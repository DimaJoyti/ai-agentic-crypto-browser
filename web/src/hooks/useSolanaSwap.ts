import React, { useState, useCallback } from 'react'
import { useWallet, useConnection } from '@solana/wallet-adapter-react'
import { Transaction, VersionedTransaction } from '@solana/web3.js'
import { SolanaDeFiService } from '@/services/solana/SolanaDeFiService'

export interface SwapQuote {
  inputAmount: number
  outputAmount: number
  priceImpact: number
  route: any[]
  fees: number
  slippageBps: number
  inputMint: string
  outputMint: string
}

export interface SwapResult {
  success: boolean
  signature?: string
  error?: string
}

export interface SwapParams {
  inputMint: string
  outputMint: string
  amount: number
  slippageBps?: number
}

export interface SolanaSwapState {
  quote: SwapQuote | null
  isLoading: boolean
  error: string | null
  lastQuoteTime: Date | null
}

export function useSolanaSwap() {
  const { connection } = useConnection()
  const { publicKey, signTransaction, signAllTransactions } = useWallet()
  
  const [state, setState] = useState<SolanaSwapState>({
    quote: null,
    isLoading: false,
    error: null,
    lastQuoteTime: null
  })

  const defiService = new SolanaDeFiService()

  const getQuote = useCallback(async (params: SwapParams): Promise<SwapQuote | null> => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      const quote = await defiService.getSwapQuote({
        inputMint: params.inputMint,
        outputMint: params.outputMint,
        amount: params.amount * 1e9, // Convert to lamports/smallest unit
        slippageBps: params.slippageBps || 50
      })

      const swapQuote: SwapQuote = {
        inputAmount: quote.inputAmount / 1e9,
        outputAmount: quote.outputAmount / 1e9,
        priceImpact: quote.priceImpact,
        route: quote.route,
        fees: quote.fees,
        slippageBps: params.slippageBps || 50,
        inputMint: params.inputMint,
        outputMint: params.outputMint
      }

      setState(prev => ({
        ...prev,
        quote: swapQuote,
        isLoading: false,
        lastQuoteTime: new Date()
      }))

      return swapQuote
    } catch (error) {
      console.error('Failed to get swap quote:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to get quote'
      }))
      return null
    }
  }, [defiService])

  const executeSwap = useCallback(async (): Promise<SwapResult> => {
    if (!state.quote || !publicKey || !signTransaction) {
      return {
        success: false,
        error: 'Missing quote, wallet, or signing capability'
      }
    }

    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      // Get swap transaction from Jupiter
      const swapResponse = await fetch('https://quote-api.jup.ag/v6/swap', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          quoteResponse: {
            inputMint: state.quote.inputMint,
            outputMint: state.quote.outputMint,
            inAmount: (state.quote.inputAmount * 1e9).toString(),
            outAmount: (state.quote.outputAmount * 1e9).toString(),
            slippageBps: state.quote.slippageBps,
            routePlan: state.quote.route
          },
          userPublicKey: publicKey.toString(),
          wrapAndUnwrapSol: true,
          dynamicComputeUnitLimit: true,
          prioritizationFeeLamports: 'auto'
        })
      })

      if (!swapResponse.ok) {
        throw new Error('Failed to get swap transaction')
      }

      const { swapTransaction } = await swapResponse.json()
      
      // Deserialize the transaction
      const transactionBuf = Buffer.from(swapTransaction, 'base64')
      let transaction: Transaction | VersionedTransaction

      try {
        // Try as VersionedTransaction first
        transaction = VersionedTransaction.deserialize(transactionBuf)
      } catch {
        // Fallback to legacy Transaction
        transaction = Transaction.from(transactionBuf)
      }

      // Sign the transaction
      const signedTransaction = await signTransaction(transaction)

      // Send the transaction
      const signature = await connection.sendRawTransaction(
        signedTransaction.serialize(),
        {
          skipPreflight: false,
          preflightCommitment: 'confirmed',
          maxRetries: 3
        }
      )

      // Confirm the transaction
      const confirmation = await connection.confirmTransaction(
        signature,
        'confirmed'
      )

      if (confirmation.value.err) {
        throw new Error(`Transaction failed: ${confirmation.value.err}`)
      }

      setState(prev => ({
        ...prev,
        isLoading: false,
        quote: null // Clear quote after successful swap
      }))

      return {
        success: true,
        signature
      }
    } catch (error) {
      console.error('Swap execution failed:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Swap execution failed'
      }))

      return {
        success: false,
        error: error instanceof Error ? error.message : 'Swap execution failed'
      }
    }
  }, [state.quote, publicKey, signTransaction, connection])

  const resetQuote = useCallback(() => {
    setState(prev => ({
      ...prev,
      quote: null,
      error: null
    }))
  }, [])

  const refreshQuote = useCallback(async (): Promise<SwapQuote | null> => {
    if (!state.quote) return null

    return getQuote({
      inputMint: state.quote.inputMint,
      outputMint: state.quote.outputMint,
      amount: state.quote.inputAmount,
      slippageBps: state.quote.slippageBps
    })
  }, [state.quote, getQuote])

  // Check if quote is stale (older than 30 seconds)
  const isQuoteStale = state.lastQuoteTime && 
    (Date.now() - state.lastQuoteTime.getTime()) > 30000

  return {
    ...state,
    getQuote,
    executeSwap,
    resetQuote,
    refreshQuote,
    isQuoteStale,
    canSwap: !!state.quote && !!publicKey && !state.isLoading
  }
}

// Helper hook for token swaps with automatic quote refresh
export function useSolanaAutoSwap(refreshInterval: number = 30000) {
  const swap = useSolanaSwap()

  // Auto-refresh quote if it becomes stale
  React.useEffect(() => {
    if (!swap.quote || !swap.isQuoteStale) return

    const interval = setInterval(() => {
      if (swap.isQuoteStale) {
        swap.refreshQuote()
      }
    }, refreshInterval)

    return () => clearInterval(interval)
  }, [swap.quote, swap.isQuoteStale, swap.refreshQuote, refreshInterval])

  return swap
}

// Helper hook for price impact warnings
export function useSwapPriceImpact(quote: SwapQuote | null) {
  const getPriceImpactSeverity = (impact: number) => {
    if (impact < 0.1) return 'low'
    if (impact < 1) return 'medium'
    if (impact < 5) return 'high'
    return 'critical'
  }

  const getPriceImpactColor = (impact: number) => {
    const severity = getPriceImpactSeverity(impact)
    switch (severity) {
      case 'low': return 'text-green-600'
      case 'medium': return 'text-yellow-600'
      case 'high': return 'text-orange-600'
      case 'critical': return 'text-red-600'
      default: return 'text-gray-600'
    }
  }

  const shouldWarnUser = (impact: number) => impact > 1

  return {
    priceImpact: quote?.priceImpact || 0,
    severity: quote ? getPriceImpactSeverity(quote.priceImpact) : 'low',
    color: quote ? getPriceImpactColor(quote.priceImpact) : 'text-gray-600',
    shouldWarn: quote ? shouldWarnUser(quote.priceImpact) : false
  }
}
