import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Address } from 'viem'
import { 
  dexIntegration,
  type DEXProtocol,
  type Token,
  type SwapQuote,
  type SwapTransaction,
  type TradingPair,
  type DEXConfig,
  type DEXEvent
} from '@/lib/dex-integration'
import { toast } from 'sonner'

export interface DEXState {
  protocols: DEXProtocol[]
  tokens: Token[]
  pairs: TradingPair[]
  quotes: SwapQuote[]
  transactions: SwapTransaction[]
  isLoading: boolean
  isSwapping: boolean
  isGettingQuote: boolean
  config: DEXConfig
  error: string | null
  lastUpdate: number | null
}

export interface UseDEXOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefreshQuotes?: boolean
  defaultSlippage?: number
}

export interface UseDEXReturn {
  // State
  state: DEXState
  
  // Quote Management
  getSwapQuote: (tokenIn: Token, tokenOut: Token, amountIn: string, slippage?: number, dexIds?: string[]) => Promise<SwapQuote[]>
  refreshQuote: (quoteId: string) => Promise<SwapQuote | null>
  
  // Swap Execution
  executeSwap: (quote: SwapQuote) => Promise<SwapTransaction>
  getSwapTransaction: (id: string) => SwapTransaction | null
  
  // Data Access
  getProtocols: (chainId?: number) => DEXProtocol[]
  getSupportedTokens: (chainId?: number) => Token[]
  getTradingPairs: (dexId?: string) => TradingPair[]
  
  // Configuration
  updateConfig: (config: Partial<DEXConfig>) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useDEX = (
  options: UseDEXOptions = {}
): UseDEXReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefreshQuotes = true,
    defaultSlippage = 0.5
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<DEXState>({
    protocols: [],
    tokens: [],
    pairs: [],
    quotes: [],
    transactions: [],
    isLoading: false,
    isSwapping: false,
    isGettingQuote: false,
    config: dexIntegration.getConfig(),
    error: null,
    lastUpdate: null
  })

  // Update state from DEX integration
  const updateState = useCallback(() => {
    try {
      const protocols = dexIntegration.getProtocols(chainId)
      const tokens = dexIntegration.getSupportedTokens(chainId || 1)
      const pairs = dexIntegration.getTradingPairs()
      const config = dexIntegration.getConfig()

      setState(prev => ({
        ...prev,
        protocols,
        tokens,
        pairs,
        config,
        error: null,
        lastUpdate: Date.now()
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [chainId])

  // Handle DEX events
  const handleDEXEvent = useCallback((event: DEXEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'swap_success':
          toast.success('Swap Successful', {
            description: `Swapped ${event.transaction?.tokenIn.symbol} for ${event.transaction?.tokenOut.symbol}`
          })
          break
        case 'swap_failed':
          toast.error('Swap Failed', {
            description: `Swap failed: ${event.error?.message || 'Unknown error'}`
          })
          break
        case 'quote_updated':
          toast.info('Quote Updated', {
            description: 'Swap quote has been refreshed with latest prices'
          })
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = dexIntegration.addEventListener(handleDEXEvent)

    // Update configuration
    dexIntegration.updateConfig({
      defaultSlippage,
      autoRefreshQuotes
    })

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, defaultSlippage, autoRefreshQuotes, handleDEXEvent, updateState])

  // Get swap quote
  const getSwapQuote = useCallback(async (
    tokenIn: Token,
    tokenOut: Token,
    amountIn: string,
    slippage: number = defaultSlippage,
    dexIds?: string[]
  ): Promise<SwapQuote[]> => {
    setState(prev => ({ ...prev, isGettingQuote: true, error: null }))

    try {
      const quotes = await dexIntegration.getSwapQuote(tokenIn, tokenOut, amountIn, slippage, dexIds)
      
      setState(prev => ({
        ...prev,
        isGettingQuote: false,
        quotes: [...prev.quotes, ...quotes]
      }))

      if (enableNotifications && quotes.length > 0) {
        const bestQuote = quotes[0]
        toast.success('Quotes Retrieved', {
          description: `Best rate: ${parseFloat(bestQuote.amountOut).toFixed(4)} ${tokenOut.symbol} from ${bestQuote.dexId}`
        })
      }

      return quotes
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isGettingQuote: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get quotes', { description: errorMessage })
      }
      throw error
    }
  }, [defaultSlippage, enableNotifications])

  // Refresh quote
  const refreshQuote = useCallback(async (quoteId: string): Promise<SwapQuote | null> => {
    const existingQuote = dexIntegration.getQuote(quoteId)
    if (!existingQuote) {
      return null
    }

    try {
      const newQuotes = await dexIntegration.getSwapQuote(
        existingQuote.tokenIn,
        existingQuote.tokenOut,
        existingQuote.amountIn,
        existingQuote.slippage,
        [existingQuote.dexId]
      )

      const refreshedQuote = newQuotes[0]
      if (refreshedQuote) {
        setState(prev => ({
          ...prev,
          quotes: prev.quotes.map(q => q.id === quoteId ? refreshedQuote : q)
        }))

        if (enableNotifications) {
          toast.info('Quote Refreshed', {
            description: `Updated rate: ${parseFloat(refreshedQuote.amountOut).toFixed(4)} ${refreshedQuote.tokenOut.symbol}`
          })
        }
      }

      return refreshedQuote
    } catch (error) {
      console.error('Failed to refresh quote:', error)
      return null
    }
  }, [enableNotifications])

  // Execute swap
  const executeSwap = useCallback(async (quote: SwapQuote): Promise<SwapTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isSwapping: true, error: null }))

    try {
      const transaction = await dexIntegration.executeSwap(quote, address)
      
      setState(prev => ({
        ...prev,
        isSwapping: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSwapping: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Get swap transaction
  const getSwapTransaction = useCallback((id: string): SwapTransaction | null => {
    return dexIntegration.getSwapTransaction(id)
  }, [])

  // Get protocols
  const getProtocols = useCallback((protocolChainId?: number): DEXProtocol[] => {
    return dexIntegration.getProtocols(protocolChainId || chainId)
  }, [chainId])

  // Get supported tokens
  const getSupportedTokens = useCallback((tokenChainId?: number): Token[] => {
    return dexIntegration.getSupportedTokens(tokenChainId || chainId || 1)
  }, [chainId])

  // Get trading pairs
  const getTradingPairs = useCallback((dexId?: string): TradingPair[] => {
    return dexIntegration.getTradingPairs(dexId)
  }, [])

  // Update configuration
  const updateConfig = useCallback((config: Partial<DEXConfig>) => {
    try {
      dexIntegration.updateConfig(config)
      setState(prev => ({ ...prev, config: dexIntegration.getConfig() }))

      if (enableNotifications) {
        toast.success('Configuration Updated', {
          description: 'DEX settings have been updated'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to update configuration', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
  }, [updateState])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    getSwapQuote,
    refreshQuote,
    executeSwap,
    getSwapTransaction,
    getProtocols,
    getSupportedTokens,
    getTradingPairs,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for token swapping
export const useTokenSwap = (tokenIn?: Token, tokenOut?: Token) => {
  const { getSwapQuote, executeSwap, state } = useDEX()

  const swap = useCallback(async (amountIn: string, slippage?: number) => {
    if (!tokenIn || !tokenOut) {
      throw new Error('Tokens not specified')
    }

    const quotes = await getSwapQuote(tokenIn, tokenOut, amountIn, slippage)
    if (quotes.length === 0) {
      throw new Error('No quotes available')
    }

    const bestQuote = quotes[0]
    return executeSwap(bestQuote)
  }, [tokenIn, tokenOut, getSwapQuote, executeSwap])

  const getQuote = useCallback(async (amountIn: string, slippage?: number) => {
    if (!tokenIn || !tokenOut) {
      throw new Error('Tokens not specified')
    }

    return getSwapQuote(tokenIn, tokenOut, amountIn, slippage)
  }, [tokenIn, tokenOut, getSwapQuote])

  return {
    swap,
    getQuote,
    isSwapping: state.isSwapping,
    isGettingQuote: state.isGettingQuote,
    error: state.error
  }
}

// Hook for DEX comparison
export const useDEXComparison = () => {
  const { getSwapQuote, state } = useDEX()

  const compareRates = useCallback(async (
    tokenIn: Token,
    tokenOut: Token,
    amountIn: string,
    slippage?: number
  ) => {
    const quotes = await getSwapQuote(tokenIn, tokenOut, amountIn, slippage)
    
    // Sort by amount out (best rate first)
    const sortedQuotes = quotes.sort((a, b) => parseFloat(b.amountOut) - parseFloat(a.amountOut))
    
    // Calculate savings compared to worst rate
    const worstRate = sortedQuotes[sortedQuotes.length - 1]
    const bestRate = sortedQuotes[0]
    
    const savings = worstRate && bestRate 
      ? ((parseFloat(bestRate.amountOut) - parseFloat(worstRate.amountOut)) / parseFloat(worstRate.amountOut)) * 100
      : 0

    return {
      quotes: sortedQuotes,
      bestQuote: bestRate,
      worstQuote: worstRate,
      savingsPercentage: savings,
      totalQuotes: quotes.length
    }
  }, [getSwapQuote])

  return {
    compareRates,
    protocols: state.protocols,
    isLoading: state.isGettingQuote
  }
}

// Hook for DEX analytics
export const useDEXAnalytics = () => {
  const { state } = useDEX()

  const analytics = {
    totalSwaps: state.transactions.length,
    successfulSwaps: state.transactions.filter(tx => tx.status === 'confirmed').length,
    failedSwaps: state.transactions.filter(tx => tx.status === 'failed').length,
    totalVolume: state.transactions.reduce((sum, tx) => {
      if (tx.status === 'confirmed') {
        return sum + parseFloat(tx.amountIn)
      }
      return sum
    }, 0),
    averageSlippage: state.transactions.length > 0
      ? state.transactions.reduce((sum, tx) => sum + tx.slippage, 0) / state.transactions.length
      : 0,
    mostUsedDEX: state.transactions.length > 0
      ? state.transactions.reduce((acc, tx) => {
          const quote = state.quotes.find(q => q.id === tx.quoteId)
          if (quote) {
            acc[quote.dexId] = (acc[quote.dexId] || 0) + 1
          }
          return acc
        }, {} as Record<string, number>)
      : {},
    successRate: state.transactions.length > 0
      ? (state.transactions.filter(tx => tx.status === 'confirmed').length / state.transactions.length) * 100
      : 0
  }

  return analytics
}
