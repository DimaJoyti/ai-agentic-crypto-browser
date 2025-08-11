import { useState, useEffect, useCallback } from 'react'
import { useWallet, useConnection } from '@solana/wallet-adapter-react'
import { LAMPORTS_PER_SOL, PublicKey } from '@solana/web3.js'
import { TOKEN_PROGRAM_ID } from '@solana/spl-token'
import axios from 'axios'

export interface TokenBalance {
  mint: string
  symbol: string
  name: string
  balance: number
  decimals: number
  uiAmount: number
  value?: number
  logo?: string
}

export interface SolanaWalletState {
  balance: number
  tokens: TokenBalance[]
  isLoading: boolean
  error: string | null
  lastUpdated: Date | null
}

export interface UseSolanaWalletOptions {
  autoRefresh?: boolean
  refreshInterval?: number
  includeTokens?: boolean
}

export function useSolanaWallet(options: UseSolanaWalletOptions = {}) {
  const {
    autoRefresh = false,
    refreshInterval = 30000,
    includeTokens = true
  } = options

  const { connection } = useConnection()
  const { publicKey, connected } = useWallet()

  const [state, setState] = useState<SolanaWalletState>({
    balance: 0,
    tokens: [],
    isLoading: false,
    error: null,
    lastUpdated: null
  })

  const fetchWalletData = useCallback(async () => {
    if (!publicKey || !connected) {
      setState(prev => ({
        ...prev,
        balance: 0,
        tokens: [],
        isLoading: false,
        error: null
      }))
      return
    }

    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      // Fetch SOL balance
      const balance = await connection.getBalance(publicKey)
      const solBalance = balance / LAMPORTS_PER_SOL

      let tokens: TokenBalance[] = []

      if (includeTokens) {
        try {
          // Fetch token accounts
          const tokenAccounts = await connection.getParsedTokenAccountsByOwner(
            publicKey,
            { programId: TOKEN_PROGRAM_ID }
          )

          // Process token balances
          const tokenPromises = tokenAccounts.value.map(async (tokenAccount) => {
            const accountData = tokenAccount.account.data.parsed
            const mintAddress = accountData.info.mint
            const balance = accountData.info.tokenAmount.uiAmount || 0
            const decimals = accountData.info.tokenAmount.decimals

            // Skip tokens with zero balance
            if (balance === 0) return null

            try {
              // Try to get token metadata
              const tokenInfo = await getTokenInfo(mintAddress)
              return {
                mint: mintAddress,
                symbol: tokenInfo?.symbol || 'UNKNOWN',
                name: tokenInfo?.name || 'Unknown Token',
                balance,
                decimals,
                uiAmount: balance,
                value: tokenInfo?.price ? balance * tokenInfo.price : undefined,
                logo: tokenInfo?.logo
              }
            } catch (error) {
              // If metadata fetch fails, return basic info
              return {
                mint: mintAddress,
                symbol: 'UNKNOWN',
                name: 'Unknown Token',
                balance,
                decimals,
                uiAmount: balance
              }
            }
          })

          const resolvedTokens = await Promise.all(tokenPromises)
          tokens = resolvedTokens.filter((token): token is TokenBalance => token !== null)
        } catch (error) {
          console.warn('Failed to fetch token balances:', error)
        }
      }

      setState(prev => ({
        ...prev,
        balance: solBalance,
        tokens,
        isLoading: false,
        lastUpdated: new Date()
      }))
    } catch (error) {
      console.error('Failed to fetch wallet data:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to fetch wallet data'
      }))
    }
  }, [publicKey, connected, connection, includeTokens])

  const refresh = useCallback(async () => {
    await fetchWalletData()
  }, [fetchWalletData])

  // Initial fetch when wallet connects
  useEffect(() => {
    if (connected && publicKey) {
      fetchWalletData()
    }
  }, [connected, publicKey, fetchWalletData])

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh || !connected) return

    const interval = setInterval(fetchWalletData, refreshInterval)
    return () => clearInterval(interval)
  }, [autoRefresh, connected, refreshInterval, fetchWalletData])

  // Computed values
  const totalValue = state.balance * 100 + // Placeholder SOL price
    state.tokens.reduce((sum, token) => sum + (token.value || 0), 0)

  const hasTokens = state.tokens.length > 0

  return {
    ...state,
    refresh,
    totalValue,
    hasTokens,
    isConnected: connected,
    address: publicKey?.toString() || null
  }
}

// Helper function to get token information
async function getTokenInfo(mintAddress: string): Promise<{
  symbol: string
  name: string
  price?: number
  logo?: string
} | null> {
  try {
    // Try Jupiter token list first
    const response = await axios.get(`https://token.jup.ag/strict`, {
      timeout: 5000
    })
    
    const token = response.data.find((t: any) => t.address === mintAddress)
    if (token) {
      return {
        symbol: token.symbol,
        name: token.name,
        logo: token.logoURI
      }
    }

    // Fallback to Solana token list
    const solanaTokenResponse = await axios.get(
      'https://raw.githubusercontent.com/solana-labs/token-list/main/src/tokens/solana.tokenlist.json',
      { timeout: 5000 }
    )
    
    const solanaToken = solanaTokenResponse.data.tokens.find((t: any) => t.address === mintAddress)
    if (solanaToken) {
      return {
        symbol: solanaToken.symbol,
        name: solanaToken.name,
        logo: solanaToken.logoURI
      }
    }

    return null
  } catch (error) {
    console.warn(`Failed to fetch token info for ${mintAddress}:`, error)
    return null
  }
}

// Helper hook for just SOL balance
export function useSolBalance() {
  const { balance, isLoading, error, refresh } = useSolanaWallet({
    includeTokens: false,
    autoRefresh: true
  })

  return {
    balance,
    isLoading,
    error,
    refresh
  }
}

// Helper hook for token balances only
export function useSolanaTokens() {
  const { tokens, isLoading, error, refresh } = useSolanaWallet({
    includeTokens: true,
    autoRefresh: true
  })

  const sortedTokens = tokens.sort((a, b) => {
    // Sort by value if available, otherwise by balance
    if (a.value && b.value) return b.value - a.value
    return b.balance - a.balance
  })

  const totalTokenValue = tokens.reduce((sum, token) => sum + (token.value || 0), 0)

  return {
    tokens: sortedTokens,
    totalTokenValue,
    tokenCount: tokens.length,
    isLoading,
    error,
    refresh
  }
}
