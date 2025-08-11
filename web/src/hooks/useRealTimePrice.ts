import { useState, useEffect, useRef } from 'react'

interface PriceData {
  symbol: string
  price: number
  change24h: number
  changePercent24h: number
  volume24h: number
  high24h: number
  low24h: number
  timestamp: number
}

interface UseRealTimePriceReturn {
  price: number | null
  priceData: PriceData | null
  isLoading: boolean
  error: string | null
  isConnected: boolean
}

export function useRealTimePrice(symbol: string): UseRealTimePriceReturn {
  const [price, setPrice] = useState<number | null>(null)
  const [priceData, setPriceData] = useState<PriceData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 5

  const generateMockPrice = (basePrice: number) => {
    const volatility = 0.001 // 0.1% volatility
    const change = (Math.random() - 0.5) * basePrice * volatility
    return basePrice + change
  }

  const getBasePrice = (symbol: string) => {
    const basePrices: Record<string, number> = {
      'BTC/USDT': 45000,
      'ETH/USDT': 2500,
      'BNB/USDT': 300,
      'ADA/USDT': 0.5,
      'SOL/USDT': 100,
      'XRP/USDT': 0.6,
      'DOT/USDT': 7,
      'DOGE/USDT': 0.08,
      'AVAX/USDT': 35,
      'MATIC/USDT': 0.9
    }
    return basePrices[symbol] || 1
  }

  const connectWebSocket = () => {
    try {
      // In a real implementation, you would connect to a real WebSocket endpoint
      // For demo purposes, we'll simulate WebSocket behavior
      setIsLoading(true)
      setError(null)

      // Simulate connection delay
      setTimeout(() => {
        setIsConnected(true)
        setIsLoading(false)
        reconnectAttempts.current = 0

        // Generate initial price data
        const basePrice = getBasePrice(symbol)
        const initialPrice = generateMockPrice(basePrice)
        const change24h = (Math.random() - 0.5) * 20 // -10% to +10%
        
        const initialData: PriceData = {
          symbol,
          price: initialPrice,
          change24h,
          changePercent24h: change24h,
          volume24h: Math.random() * 1000000000,
          high24h: basePrice * 1.05,
          low24h: basePrice * 0.95,
          timestamp: Date.now()
        }

        setPrice(initialPrice)
        setPriceData(initialData)

        // Start price updates
        const interval = setInterval(() => {
          const newPrice = generateMockPrice(initialData.price)
          const newChange = ((newPrice - basePrice) / basePrice) * 100
          
          const updatedData: PriceData = {
            ...initialData,
            price: newPrice,
            change24h: newChange,
            changePercent24h: newChange,
            timestamp: Date.now()
          }

          setPrice(newPrice)
          setPriceData(updatedData)
        }, 1000) // Update every second

        // Store interval reference for cleanup
        wsRef.current = { close: () => clearInterval(interval) } as any
      }, 500)

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to connect')
      setIsLoading(false)
      scheduleReconnect()
    }
  }

  const scheduleReconnect = () => {
    if (reconnectAttempts.current < maxReconnectAttempts) {
      const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000)
      reconnectTimeoutRef.current = setTimeout(() => {
        reconnectAttempts.current++
        connectWebSocket()
      }, delay)
    } else {
      setError('Max reconnection attempts reached')
    }
  }

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    
    setIsConnected(false)
    reconnectAttempts.current = 0
  }

  useEffect(() => {
    if (symbol) {
      connectWebSocket()
    }

    return () => {
      disconnect()
    }
  }, [symbol])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [])

  return {
    price,
    priceData,
    isLoading,
    error,
    isConnected
  }
}
