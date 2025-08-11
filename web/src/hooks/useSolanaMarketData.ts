import { useState, useEffect, useCallback } from 'react'
import { SolanaMarketService } from '@/services/solana/SolanaMarketService'

export interface SolPrice {
  price: number
  change24h: number
  change7d: number
  volume24h: number
  marketCap: number
  lastUpdated: Date
}

export interface MarketData {
  marketCap: number
  volume24h: number
  circulatingSupply: number
  totalSupply: number
  transactions24h: number
  activeAddresses24h: number
  networkHashRate: number
  blockTime: number
}

export interface SolanaMarketDataState {
  solPrice: SolPrice | null
  marketData: MarketData | null
  isLoading: boolean
  error: string | null
  lastUpdated: Date | null
}

export interface UseSolanaMarketDataOptions {
  autoRefresh?: boolean
  refreshInterval?: number
  enableWebSocket?: boolean
}

export function useSolanaMarketData(options: UseSolanaMarketDataOptions = {}) {
  const {
    autoRefresh = false,
    refreshInterval = 30000,
    enableWebSocket = true
  } = options

  const [state, setState] = useState<SolanaMarketDataState>({
    solPrice: null,
    marketData: null,
    isLoading: true,
    error: null,
    lastUpdated: null
  })

  const marketService = new SolanaMarketService()

  const fetchMarketData = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      const [priceData, marketStats] = await Promise.all([
        marketService.getSolPrice(),
        marketService.getMarketData()
      ])

      setState(prev => ({
        ...prev,
        solPrice: priceData,
        marketData: marketStats,
        isLoading: false,
        lastUpdated: new Date()
      }))
    } catch (error) {
      console.error('Failed to fetch Solana market data:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to fetch market data'
      }))
    }
  }, [])

  const refresh = useCallback(async () => {
    await fetchMarketData()
  }, [fetchMarketData])

  // Initial data fetch
  useEffect(() => {
    fetchMarketData()
  }, [fetchMarketData])

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(fetchMarketData, refreshInterval)
    return () => clearInterval(interval)
  }, [autoRefresh, refreshInterval, fetchMarketData])

  // WebSocket connection for real-time updates
  useEffect(() => {
    if (!enableWebSocket) return

    let ws: WebSocket | null = null

    const connectWebSocket = () => {
      try {
        // Connect to Solana price feed WebSocket
        ws = new WebSocket('wss://api.coingecko.com/api/v3/coins/solana/tickers')
        
        ws.onopen = () => {
          console.log('Connected to Solana price WebSocket')
        }

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            if (data.type === 'price_update') {
              setState(prev => ({
                ...prev,
                solPrice: prev.solPrice ? {
                  ...prev.solPrice,
                  price: data.price,
                  change24h: data.change24h,
                  lastUpdated: new Date()
                } : null
              }))
            }
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error)
          }
        }

        ws.onerror = (error) => {
          console.error('WebSocket error:', error)
        }

        ws.onclose = () => {
          console.log('WebSocket connection closed')
          // Attempt to reconnect after 5 seconds
          setTimeout(connectWebSocket, 5000)
        }
      } catch (error) {
        console.error('Failed to connect to WebSocket:', error)
      }
    }

    connectWebSocket()

    return () => {
      if (ws) {
        ws.close()
      }
    }
  }, [enableWebSocket])

  return {
    ...state,
    refresh
  }
}

// Helper hook for just SOL price
export function useSolPrice(options: UseSolanaMarketDataOptions = {}) {
  const { solPrice, isLoading, error, refresh } = useSolanaMarketData(options)
  
  return {
    price: solPrice?.price || 0,
    change24h: solPrice?.change24h || 0,
    change7d: solPrice?.change7d || 0,
    volume24h: solPrice?.volume24h || 0,
    marketCap: solPrice?.marketCap || 0,
    isLoading,
    error,
    refresh
  }
}

// Helper hook for network statistics
export function useSolanaNetworkStats(options: UseSolanaMarketDataOptions = {}) {
  const { marketData, isLoading, error, refresh } = useSolanaMarketData(options)
  
  return {
    transactions24h: marketData?.transactions24h || 0,
    activeAddresses24h: marketData?.activeAddresses24h || 0,
    networkHashRate: marketData?.networkHashRate || 0,
    blockTime: marketData?.blockTime || 0,
    isLoading,
    error,
    refresh
  }
}
