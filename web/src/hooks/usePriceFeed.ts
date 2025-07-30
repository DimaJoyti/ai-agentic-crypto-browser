import { useState, useEffect, useCallback, useRef } from 'react'
import { toast } from 'sonner'
import { 
  priceFeedManager, 
  type PriceData, 
  type PriceAlert,
  type PriceFeedConfig 
} from '@/lib/price-feed-manager'

export interface PriceFeedState {
  prices: Map<string, PriceData>
  isConnected: boolean
  isLoading: boolean
  lastUpdate: number | null
  error: string | null
  stats: {
    activeSources: number
    cachedPrices: number
    activeAlerts: number
    subscribers: number
    isRunning: boolean
  }
}

export interface UsePriceFeedOptions {
  symbols?: string[]
  autoStart?: boolean
  enableAlerts?: boolean
  onPriceUpdate?: (symbol: string, price: PriceData) => void
  onAlert?: (alert: PriceAlert, price: PriceData) => void
  onError?: (error: Error) => void
}

export interface UsePriceFeedReturn {
  // State
  state: PriceFeedState
  
  // Price data
  getPrice: (symbol: string) => PriceData | null
  getPrices: (symbols: string[]) => Map<string, PriceData>
  getPriceHistory: (symbol: string, limit?: number) => PriceData[]
  
  // Connection management
  start: () => Promise<void>
  stop: () => void
  restart: () => Promise<void>
  
  // Subscriptions
  subscribe: (symbols: string | string[]) => () => void
  
  // Alerts
  addAlert: (alert: Omit<PriceAlert, 'id' | 'createdAt'>) => string
  removeAlert: (alertId: string) => boolean
  getAlerts: () => PriceAlert[]
  getAlertsForSymbol: (symbol: string) => PriceAlert[]
  
  // Utilities
  clearError: () => void
  formatPrice: (price: number, currency?: string) => string
  formatChange: (change: number, isPercent?: boolean) => string
  calculatePriceChange: (current: number, previous: number) => { change: number; changePercent: number }
}

export const usePriceFeed = (options: UsePriceFeedOptions = {}): UsePriceFeedReturn => {
  const [state, setState] = useState<PriceFeedState>({
    prices: new Map(),
    isConnected: false,
    isLoading: false,
    lastUpdate: null,
    error: null,
    stats: {
      activeSources: 0,
      cachedPrices: 0,
      activeAlerts: 0,
      subscribers: 0,
      isRunning: false
    }
  })

  const subscriptionRef = useRef<(() => void) | null>(null)
  const alertListenerRef = useRef<((event: CustomEvent) => void) | null>(null)

  const {
    symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
    autoStart = true,
    enableAlerts = true,
    onPriceUpdate,
    onAlert,
    onError
  } = options

  // Initialize price feed
  useEffect(() => {
    if (autoStart) {
      start()
    }

    return () => {
      stop()
    }
  }, [autoStart])

  // Set up alert listener
  useEffect(() => {
    if (enableAlerts) {
      const handleAlert = (event: CustomEvent) => {
        const { alert, price } = event.detail
        onAlert?.(alert, price)
        
        toast.info(`Price Alert: ${alert.symbol}`, {
          description: `${alert.symbol} is ${alert.type} ${alert.value}`
        })
      }

      alertListenerRef.current = handleAlert
      window.addEventListener('priceAlert', handleAlert as EventListener)

      return () => {
        if (alertListenerRef.current) {
          window.removeEventListener('priceAlert', alertListenerRef.current as EventListener)
        }
      }
    }
  }, [enableAlerts, onAlert])

  // Update stats periodically
  useEffect(() => {
    const updateStats = () => {
      const stats = priceFeedManager.getStats()
      setState(prev => ({ ...prev, stats }))
    }

    const interval = setInterval(updateStats, 1000)
    updateStats() // Initial update

    return () => clearInterval(interval)
  }, [])

  // Start price feed
  const start = useCallback(async (): Promise<void> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      await priceFeedManager.start()
      
      // Subscribe to price updates
      const unsubscribe = subscribe(symbols)
      subscriptionRef.current = unsubscribe

      setState(prev => ({
        ...prev,
        isConnected: true,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      toast.success('Price feed connected')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isConnected: false,
        isLoading: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Failed to connect price feed', {
        description: errorMessage
      })
    }
  }, [symbols, onError])

  // Stop price feed
  const stop = useCallback((): void => {
    priceFeedManager.stop()
    
    if (subscriptionRef.current) {
      subscriptionRef.current()
      subscriptionRef.current = null
    }

    setState(prev => ({
      ...prev,
      isConnected: false,
      isLoading: false
    }))

    toast.info('Price feed disconnected')
  }, [])

  // Restart price feed
  const restart = useCallback(async (): Promise<void> => {
    stop()
    await new Promise(resolve => setTimeout(resolve, 1000)) // Wait a bit
    await start()
  }, [start, stop])

  // Subscribe to price updates
  const subscribe = useCallback((symbolsToSubscribe: string | string[]): (() => void) => {
    const symbolList = Array.isArray(symbolsToSubscribe) ? symbolsToSubscribe : [symbolsToSubscribe]
    
    const unsubscribe = priceFeedManager.subscribe(symbolList, (priceData) => {
      setState(prev => {
        const newPrices = new Map(prev.prices)
        newPrices.set(priceData.symbol, priceData)
        
        return {
          ...prev,
          prices: newPrices,
          lastUpdate: Date.now()
        }
      })

      onPriceUpdate?.(priceData.symbol, priceData)
    })

    return unsubscribe
  }, [onPriceUpdate])

  // Get price for a symbol
  const getPrice = useCallback((symbol: string): PriceData | null => {
    return priceFeedManager.getPrice(symbol)
  }, [])

  // Get prices for multiple symbols
  const getPrices = useCallback((symbolsToGet: string[]): Map<string, PriceData> => {
    return priceFeedManager.getPrices(symbolsToGet)
  }, [])

  // Get price history
  const getPriceHistory = useCallback((symbol: string, limit = 100): PriceData[] => {
    return priceFeedManager.getPriceHistory(symbol, limit)
  }, [])

  // Add price alert
  const addAlert = useCallback((alert: Omit<PriceAlert, 'id' | 'createdAt'>): string => {
    const alertId = priceFeedManager.addAlert(alert)
    toast.success('Price alert added', {
      description: `Alert for ${alert.symbol} ${alert.type} ${alert.value}`
    })
    return alertId
  }, [])

  // Remove price alert
  const removeAlert = useCallback((alertId: string): boolean => {
    const removed = priceFeedManager.removeAlert(alertId)
    if (removed) {
      toast.success('Price alert removed')
    }
    return removed
  }, [])

  // Get all alerts
  const getAlerts = useCallback((): PriceAlert[] => {
    return priceFeedManager.getAlerts()
  }, [])

  // Get alerts for symbol
  const getAlertsForSymbol = useCallback((symbol: string): PriceAlert[] => {
    return priceFeedManager.getAlertsForSymbol(symbol)
  }, [])

  // Clear error
  const clearError = useCallback((): void => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Format price
  const formatPrice = useCallback((price: number, currency = 'USD'): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
      minimumFractionDigits: price < 1 ? 6 : 2,
      maximumFractionDigits: price < 1 ? 6 : 2
    }).format(price)
  }, [])

  // Format change
  const formatChange = useCallback((change: number, isPercent = false): string => {
    const sign = change >= 0 ? '+' : ''
    const suffix = isPercent ? '%' : ''
    
    if (isPercent) {
      return `${sign}${change.toFixed(2)}${suffix}`
    } else {
      return `${sign}${formatPrice(Math.abs(change)).replace('$', '')}${suffix}`
    }
  }, [formatPrice])

  // Calculate price change
  const calculatePriceChange = useCallback((current: number, previous: number) => {
    const change = current - previous
    const changePercent = previous !== 0 ? (change / previous) * 100 : 0
    
    return { change, changePercent }
  }, [])

  return {
    state,
    getPrice,
    getPrices,
    getPriceHistory,
    start,
    stop,
    restart,
    subscribe,
    addAlert,
    removeAlert,
    getAlerts,
    getAlertsForSymbol,
    clearError,
    formatPrice,
    formatChange,
    calculatePriceChange
  }
}

// Hook for single symbol price tracking
export const usePrice = (symbol: string) => {
  const [price, setPrice] = useState<PriceData | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const { subscribe, getPrice } = usePriceFeed({
    symbols: [symbol],
    autoStart: true
  })

  useEffect(() => {
    // Get initial price
    const initialPrice = getPrice(symbol)
    if (initialPrice) {
      setPrice(initialPrice)
      setIsLoading(false)
    }

    // Subscribe to updates
    const unsubscribe = subscribe(symbol)

    return unsubscribe
  }, [symbol, subscribe, getPrice])

  useEffect(() => {
    const unsubscribe = subscribe(symbol)
    
    // Set up price update handler
    const handlePriceUpdate = (priceData: PriceData) => {
      if (priceData.symbol === symbol.toUpperCase()) {
        setPrice(priceData)
        setIsLoading(false)
      }
    }

    // Subscribe to price updates
    const feedUnsubscribe = priceFeedManager.subscribe(symbol, handlePriceUpdate)

    return () => {
      unsubscribe()
      feedUnsubscribe()
    }
  }, [symbol, subscribe])

  return {
    price,
    isLoading,
    symbol: symbol.toUpperCase()
  }
}

// Hook for multiple symbols price tracking
export const usePrices = (symbols: string[]) => {
  const [prices, setPrices] = useState<Map<string, PriceData>>(new Map())
  const [isLoading, setIsLoading] = useState(true)

  const { subscribe, getPrices } = usePriceFeed({
    symbols,
    autoStart: true
  })

  useEffect(() => {
    // Get initial prices
    const initialPrices = getPrices(symbols)
    if (initialPrices.size > 0) {
      setPrices(initialPrices)
      setIsLoading(false)
    }

    // Subscribe to updates
    const unsubscribe = subscribe(symbols)

    return unsubscribe
  }, [symbols, subscribe, getPrices])

  useEffect(() => {
    const unsubscribes: (() => void)[] = []

    for (const symbol of symbols) {
      const unsubscribe = priceFeedManager.subscribe(symbol, (priceData) => {
        setPrices(prev => {
          const newPrices = new Map(prev)
          newPrices.set(priceData.symbol, priceData)
          return newPrices
        })
        setIsLoading(false)
      })
      unsubscribes.push(unsubscribe)
    }

    return () => {
      unsubscribes.forEach(unsub => unsub())
    }
  }, [symbols])

  return {
    prices,
    isLoading,
    symbols: symbols.map(s => s.toUpperCase())
  }
}
