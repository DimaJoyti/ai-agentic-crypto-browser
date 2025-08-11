import { EventEmitter } from 'events'

export interface MarketTicker {
  symbol: string
  price: number
  change24h: number
  changePercent24h: number
  volume24h: number
  high24h: number
  low24h: number
  timestamp: number
}

export interface OrderBookEntry {
  price: number
  amount: number
  total: number
  count?: number
}

export interface OrderBookData {
  symbol: string
  bids: OrderBookEntry[]
  asks: OrderBookEntry[]
  lastUpdate: number
}

export interface Trade {
  id: string
  symbol: string
  price: number
  amount: number
  side: 'buy' | 'sell'
  timestamp: number
}

export interface CandlestickData {
  timestamp: number
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export interface MarketDataSubscription {
  type: 'ticker' | 'orderbook' | 'trades' | 'kline'
  symbol: string
  interval?: string
  callback: (data: any) => void
}

class MarketDataService extends EventEmitter {
  private ws: WebSocket | null = null
  private subscriptions: Map<string, MarketDataSubscription> = new Map()
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private isConnecting = false
  private heartbeatInterval: NodeJS.Timeout | null = null
  private marketData: Map<string, MarketTicker> = new Map()
  private orderBooks: Map<string, OrderBookData> = new Map()

  constructor() {
    super()
    this.initializeMarketData()
  }

  private initializeMarketData() {
    // Initialize with mock data for popular trading pairs
    const symbols = [
      'BTC/USDT', 'ETH/USDT', 'BNB/USDT', 'ADA/USDT', 'SOL/USDT',
      'XRP/USDT', 'DOT/USDT', 'DOGE/USDT', 'AVAX/USDT', 'MATIC/USDT'
    ]

    symbols.forEach(symbol => {
      const basePrice = this.getBasePrice(symbol)
      const changePercent = (Math.random() - 0.5) * 20 // -10% to +10%
      
      this.marketData.set(symbol, {
        symbol,
        price: basePrice * (1 + changePercent / 100),
        change24h: basePrice * (changePercent / 100),
        changePercent24h: changePercent,
        volume24h: Math.random() * 1000000000,
        high24h: basePrice * 1.1,
        low24h: basePrice * 0.9,
        timestamp: Date.now()
      })

      // Initialize order book
      this.orderBooks.set(symbol, this.generateOrderBook(symbol, basePrice))
    })

    // Start price simulation
    this.startPriceSimulation()
  }

  private getBasePrice(symbol: string): number {
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

  private generateOrderBook(symbol: string, basePrice: number): OrderBookData {
    const bids: OrderBookEntry[] = []
    const asks: OrderBookEntry[] = []

    // Generate bids (buy orders) - decreasing prices
    for (let i = 0; i < 20; i++) {
      const price = basePrice - (i * 0.1) - Math.random() * 0.1
      const amount = Math.random() * 10 + 0.1
      const total = price * amount
      bids.push({ price, amount, total })
    }

    // Generate asks (sell orders) - increasing prices
    for (let i = 0; i < 20; i++) {
      const price = basePrice + (i * 0.1) + Math.random() * 0.1
      const amount = Math.random() * 10 + 0.1
      const total = price * amount
      asks.push({ price, amount, total })
    }

    return {
      symbol,
      bids,
      asks,
      lastUpdate: Date.now()
    }
  }

  private startPriceSimulation() {
    setInterval(() => {
      this.marketData.forEach((ticker, symbol) => {
        // Simulate price movement
        const volatility = 0.001 // 0.1% volatility
        const change = (Math.random() - 0.5) * ticker.price * volatility
        const newPrice = Math.max(0.01, ticker.price + change)
        
        const updatedTicker: MarketTicker = {
          ...ticker,
          price: newPrice,
          timestamp: Date.now()
        }

        this.marketData.set(symbol, updatedTicker)
        
        // Update order book
        const orderBook = this.generateOrderBook(symbol, newPrice)
        this.orderBooks.set(symbol, orderBook)

        // Emit updates to subscribers
        this.emit('ticker', updatedTicker)
        this.emit('orderbook', orderBook)
      })
    }, 1000) // Update every second
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      if (this.isConnecting) {
        return
      }

      this.isConnecting = true

      try {
        // In a real implementation, this would connect to actual exchange WebSocket
        // For demo, we'll simulate the connection
        setTimeout(() => {
          this.isConnecting = false
          this.reconnectAttempts = 0
          this.emit('connected')
          this.startHeartbeat()
          resolve()
        }, 1000)

      } catch (error) {
        this.isConnecting = false
        this.handleConnectionError(error)
        reject(error)
      }
    })
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }

    this.subscriptions.clear()
    this.emit('disconnected')
  }

  private handleConnectionError(error: any) {
    console.error('WebSocket connection error:', error)
    this.emit('error', error)

    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      setTimeout(() => {
        this.reconnectAttempts++
        this.connect()
      }, this.reconnectDelay * Math.pow(2, this.reconnectAttempts))
    }
  }

  private startHeartbeat() {
    this.heartbeatInterval = setInterval(() => {
      // In real implementation, send ping to server
      this.emit('heartbeat')
    }, 30000) // Every 30 seconds
  }

  subscribe(subscription: MarketDataSubscription): string {
    const id = `${subscription.type}_${subscription.symbol}_${Date.now()}`
    this.subscriptions.set(id, subscription)

    // Set up event listener
    this.on(subscription.type, (data: any) => {
      if (data.symbol === subscription.symbol) {
        subscription.callback(data)
      }
    })

    // Send initial data if available
    if (subscription.type === 'ticker') {
      const ticker = this.marketData.get(subscription.symbol)
      if (ticker) {
        subscription.callback(ticker)
      }
    } else if (subscription.type === 'orderbook') {
      const orderBook = this.orderBooks.get(subscription.symbol)
      if (orderBook) {
        subscription.callback(orderBook)
      }
    }

    return id
  }

  unsubscribe(subscriptionId: string) {
    const subscription = this.subscriptions.get(subscriptionId)
    if (subscription) {
      this.removeAllListeners(subscription.type)
      this.subscriptions.delete(subscriptionId)
    }
  }

  getTicker(symbol: string): MarketTicker | null {
    return this.marketData.get(symbol) || null
  }

  getOrderBook(symbol: string): OrderBookData | null {
    return this.orderBooks.get(symbol) || null
  }

  getAllTickers(): MarketTicker[] {
    return Array.from(this.marketData.values())
  }

  getTopMovers(limit: number = 10): MarketTicker[] {
    return Array.from(this.marketData.values())
      .sort((a, b) => Math.abs(b.changePercent24h) - Math.abs(a.changePercent24h))
      .slice(0, limit)
  }

  getTopGainers(limit: number = 10): MarketTicker[] {
    return Array.from(this.marketData.values())
      .filter(ticker => ticker.changePercent24h > 0)
      .sort((a, b) => b.changePercent24h - a.changePercent24h)
      .slice(0, limit)
  }

  getTopLosers(limit: number = 10): MarketTicker[] {
    return Array.from(this.marketData.values())
      .filter(ticker => ticker.changePercent24h < 0)
      .sort((a, b) => a.changePercent24h - b.changePercent24h)
      .slice(0, limit)
  }

  getCandlestickData(symbol: string, interval: string, limit: number = 100): CandlestickData[] {
    // Generate mock candlestick data
    const data: CandlestickData[] = []
    const ticker = this.marketData.get(symbol)
    if (!ticker) return data

    const intervalMs = this.getIntervalMs(interval)
    const now = Date.now()
    let currentPrice = ticker.price

    for (let i = limit - 1; i >= 0; i--) {
      const timestamp = now - (i * intervalMs)
      const volatility = 0.02 // 2% volatility
      
      const open = currentPrice
      const change = (Math.random() - 0.5) * currentPrice * volatility
      const close = open + change
      const high = Math.max(open, close) + Math.random() * currentPrice * volatility * 0.5
      const low = Math.min(open, close) - Math.random() * currentPrice * volatility * 0.5
      const volume = Math.random() * 1000 + 100

      data.push({
        timestamp,
        open,
        high,
        low,
        close,
        volume
      })

      currentPrice = close
    }

    return data
  }

  private getIntervalMs(interval: string): number {
    const intervals: Record<string, number> = {
      '1m': 60 * 1000,
      '5m': 5 * 60 * 1000,
      '15m': 15 * 60 * 1000,
      '1h': 60 * 60 * 1000,
      '4h': 4 * 60 * 60 * 1000,
      '1d': 24 * 60 * 60 * 1000,
      '1w': 7 * 24 * 60 * 60 * 1000
    }
    return intervals[interval] || intervals['1h']
  }

  // Market statistics
  getMarketStats() {
    const tickers = Array.from(this.marketData.values())
    const totalVolume = tickers.reduce((sum, ticker) => sum + ticker.volume24h, 0)
    const gainers = tickers.filter(t => t.changePercent24h > 0).length
    const losers = tickers.filter(t => t.changePercent24h < 0).length
    const avgChange = tickers.reduce((sum, t) => sum + t.changePercent24h, 0) / tickers.length

    return {
      totalPairs: tickers.length,
      totalVolume24h: totalVolume,
      gainers,
      losers,
      avgChange24h: avgChange,
      lastUpdate: Date.now()
    }
  }
}

// Singleton instance
export const marketDataService = new MarketDataService()

// React hook for easy integration
export function useMarketData() {
  return {
    connect: () => marketDataService.connect(),
    disconnect: () => marketDataService.disconnect(),
    subscribe: (subscription: MarketDataSubscription) => marketDataService.subscribe(subscription),
    unsubscribe: (id: string) => marketDataService.unsubscribe(id),
    getTicker: (symbol: string) => marketDataService.getTicker(symbol),
    getOrderBook: (symbol: string) => marketDataService.getOrderBook(symbol),
    getAllTickers: () => marketDataService.getAllTickers(),
    getTopMovers: (limit?: number) => marketDataService.getTopMovers(limit),
    getTopGainers: (limit?: number) => marketDataService.getTopGainers(limit),
    getTopLosers: (limit?: number) => marketDataService.getTopLosers(limit),
    getCandlestickData: (symbol: string, interval: string, limit?: number) => 
      marketDataService.getCandlestickData(symbol, interval, limit),
    getMarketStats: () => marketDataService.getMarketStats()
  }
}
