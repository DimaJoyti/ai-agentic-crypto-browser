export interface PriceData {
  symbol: string
  price: number
  change24h: number
  changePercent24h: number
  volume24h: number
  marketCap: number
  high24h: number
  low24h: number
  timestamp: number
  source: string
}

export interface PriceFeedConfig {
  symbols: string[]
  updateInterval: number
  enableWebSocket: boolean
  dataSources: string[]
  aggregationMethod: 'average' | 'median' | 'weighted'
  maxRetries: number
  timeout: number
}

export interface PriceAlert {
  id: string
  symbol: string
  type: 'above' | 'below' | 'change_percent'
  value: number
  isActive: boolean
  triggered: boolean
  createdAt: number
  triggeredAt?: number
  userId?: string
}

export interface DataSource {
  name: string
  isActive: boolean
  priority: number
  lastUpdate: number
  errorCount: number
  connect(): Promise<void>
  disconnect(): void
  subscribe(symbols: string[]): void
  unsubscribe(symbols: string[]): void
  getPrices(symbols: string[]): Promise<PriceData[]>
}

export class PriceFeedManager {
  private static instance: PriceFeedManager
  private config: PriceFeedConfig
  private dataSources = new Map<string, DataSource>()
  private priceCache = new Map<string, PriceData>()
  private priceHistory = new Map<string, PriceData[]>()
  private alerts = new Map<string, PriceAlert>()
  private subscribers = new Map<string, Set<(data: PriceData) => void>>()
  private updateInterval: NodeJS.Timeout | null = null
  private isRunning = false

  private constructor(config: Partial<PriceFeedConfig> = {}) {
    this.config = {
      symbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'],
      updateInterval: 5000, // 5 seconds
      enableWebSocket: true,
      dataSources: ['coingecko', 'coinmarketcap', 'binance'],
      aggregationMethod: 'weighted',
      maxRetries: 3,
      timeout: 10000,
      ...config
    }

    this.initializeDataSources()
  }

  static getInstance(config?: Partial<PriceFeedConfig>): PriceFeedManager {
    if (!PriceFeedManager.instance) {
      PriceFeedManager.instance = new PriceFeedManager(config)
    }
    return PriceFeedManager.instance
  }

  /**
   * Initialize data sources
   */
  private initializeDataSources(): void {
    // Initialize CoinGecko data source
    if (this.config.dataSources.includes('coingecko')) {
      this.dataSources.set('coingecko', new CoinGeckoDataSource())
    }

    // Initialize CoinMarketCap data source
    if (this.config.dataSources.includes('coinmarketcap')) {
      this.dataSources.set('coinmarketcap', new CoinMarketCapDataSource())
    }

    // Initialize Binance data source
    if (this.config.dataSources.includes('binance')) {
      this.dataSources.set('binance', new BinanceDataSource())
    }
  }

  /**
   * Start price feed updates
   */
  async start(): Promise<void> {
    if (this.isRunning) return

    this.isRunning = true

    try {
      // Connect to data sources
      for (const [name, source] of Array.from(this.dataSources)) {
        try {
          await source.connect()
          source.subscribe(this.config.symbols)
          console.log(`Connected to ${name} data source`)
        } catch (error) {
          console.error(`Failed to connect to ${name}:`, error)
        }
      }

      // Start periodic updates
      this.startPeriodicUpdates()

      console.log('Price feed manager started')
    } catch (error) {
      console.error('Failed to start price feed manager:', error)
      this.isRunning = false
      throw error
    }
  }

  /**
   * Stop price feed updates
   */
  stop(): void {
    if (!this.isRunning) return

    this.isRunning = false

    // Disconnect from data sources
    for (const [name, source] of Array.from(this.dataSources)) {
      try {
        source.disconnect()
        console.log(`Disconnected from ${name} data source`)
      } catch (error) {
        console.error(`Failed to disconnect from ${name}:`, error)
      }
    }

    // Stop periodic updates
    if (this.updateInterval) {
      clearInterval(this.updateInterval)
      this.updateInterval = null
    }

    console.log('Price feed manager stopped')
  }

  /**
   * Subscribe to price updates for specific symbols
   */
  subscribe(symbols: string | string[], callback: (data: PriceData) => void): () => void {
    const symbolList = Array.isArray(symbols) ? symbols : [symbols]

    for (const symbol of symbolList) {
      if (!this.subscribers.has(symbol)) {
        this.subscribers.set(symbol, new Set())
      }
      this.subscribers.get(symbol)!.add(callback)
    }

    // Return unsubscribe function
    return () => {
      for (const symbol of symbolList) {
        const callbacks = this.subscribers.get(symbol)
        if (callbacks) {
          callbacks.delete(callback)
          if (callbacks.size === 0) {
            this.subscribers.delete(symbol)
          }
        }
      }
    }
  }

  /**
   * Get current price for a symbol
   */
  getPrice(symbol: string): PriceData | null {
    return this.priceCache.get(symbol.toUpperCase()) || null
  }

  /**
   * Get current prices for multiple symbols
   */
  getPrices(symbols: string[]): Map<string, PriceData> {
    const prices = new Map<string, PriceData>()
    for (const symbol of symbols) {
      const price = this.getPrice(symbol)
      if (price) {
        prices.set(symbol, price)
      }
    }
    return prices
  }

  /**
   * Get price history for a symbol
   */
  getPriceHistory(symbol: string, limit = 100): PriceData[] {
    const history = this.priceHistory.get(symbol.toUpperCase()) || []
    return history.slice(-limit)
  }

  /**
   * Add price alert
   */
  addAlert(alert: Omit<PriceAlert, 'id' | 'createdAt'>): string {
    const id = `alert_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`
    const fullAlert: PriceAlert = {
      ...alert,
      id,
      createdAt: Date.now()
    }

    this.alerts.set(id, fullAlert)
    return id
  }

  /**
   * Remove price alert
   */
  removeAlert(alertId: string): boolean {
    return this.alerts.delete(alertId)
  }

  /**
   * Get all alerts
   */
  getAlerts(): PriceAlert[] {
    return Array.from(this.alerts.values())
  }

  /**
   * Get alerts for a specific symbol
   */
  getAlertsForSymbol(symbol: string): PriceAlert[] {
    return Array.from(this.alerts.values()).filter(alert => 
      alert.symbol.toUpperCase() === symbol.toUpperCase()
    )
  }

  /**
   * Start periodic price updates
   */
  private startPeriodicUpdates(): void {
    this.updateInterval = setInterval(async () => {
      await this.updatePrices()
    }, this.config.updateInterval)

    // Initial update
    this.updatePrices()
  }

  /**
   * Update prices from all data sources
   */
  private async updatePrices(): Promise<void> {
    const allPrices = new Map<string, PriceData[]>()

    // Collect prices from all sources
    for (const [sourceName, source] of Array.from(this.dataSources)) {
      if (!source.isActive) continue

      try {
        const prices = await source.getPrices(this.config.symbols)
        
        for (const price of prices) {
          if (!allPrices.has(price.symbol)) {
            allPrices.set(price.symbol, [])
          }
          allPrices.get(price.symbol)!.push(price)
        }
      } catch (error) {
        console.error(`Failed to get prices from ${sourceName}:`, error)
        source.errorCount++
      }
    }

    // Aggregate prices and update cache
    for (const [symbol, prices] of Array.from(allPrices)) {
      if (prices.length === 0) continue

      const aggregatedPrice = this.aggregatePrices(prices)
      const previousPrice = this.priceCache.get(symbol)

      // Update cache
      this.priceCache.set(symbol, aggregatedPrice)

      // Update history
      this.updatePriceHistory(symbol, aggregatedPrice)

      // Check alerts
      this.checkAlerts(symbol, aggregatedPrice, previousPrice)

      // Notify subscribers
      this.notifySubscribers(symbol, aggregatedPrice)
    }
  }

  /**
   * Aggregate prices from multiple sources
   */
  private aggregatePrices(prices: PriceData[]): PriceData {
    if (prices.length === 1) return prices[0]

    const symbol = prices[0].symbol
    const timestamp = Date.now()

    switch (this.config.aggregationMethod) {
      case 'average':
        return {
          symbol,
          price: prices.reduce((sum, p) => sum + p.price, 0) / prices.length,
          change24h: prices.reduce((sum, p) => sum + p.change24h, 0) / prices.length,
          changePercent24h: prices.reduce((sum, p) => sum + p.changePercent24h, 0) / prices.length,
          volume24h: prices.reduce((sum, p) => sum + p.volume24h, 0) / prices.length,
          marketCap: prices.reduce((sum, p) => sum + p.marketCap, 0) / prices.length,
          high24h: Math.max(...prices.map(p => p.high24h)),
          low24h: Math.min(...prices.map(p => p.low24h)),
          timestamp,
          source: 'aggregated'
        }

      case 'median':
        const sortedPrices = prices.map(p => p.price).sort((a, b) => a - b)
        const medianPrice = sortedPrices[Math.floor(sortedPrices.length / 2)]
        const medianPriceData = prices.find(p => p.price === medianPrice) || prices[0]
        return { ...medianPriceData, timestamp, source: 'aggregated' }

      case 'weighted':
        // Weight by source priority and reliability
        let totalWeight = 0
        let weightedPrice = 0
        let weightedChange24h = 0
        let weightedChangePercent24h = 0

        for (const price of prices) {
          const source = this.dataSources.get(price.source)
          const weight = source ? (source.priority * (1 - source.errorCount * 0.1)) : 1
          
          totalWeight += weight
          weightedPrice += price.price * weight
          weightedChange24h += price.change24h * weight
          weightedChangePercent24h += price.changePercent24h * weight
        }

        return {
          symbol,
          price: weightedPrice / totalWeight,
          change24h: weightedChange24h / totalWeight,
          changePercent24h: weightedChangePercent24h / totalWeight,
          volume24h: prices.reduce((sum, p) => sum + p.volume24h, 0) / prices.length,
          marketCap: prices.reduce((sum, p) => sum + p.marketCap, 0) / prices.length,
          high24h: Math.max(...prices.map(p => p.high24h)),
          low24h: Math.min(...prices.map(p => p.low24h)),
          timestamp,
          source: 'aggregated'
        }

      default:
        return prices[0]
    }
  }

  /**
   * Update price history
   */
  private updatePriceHistory(symbol: string, price: PriceData): void {
    if (!this.priceHistory.has(symbol)) {
      this.priceHistory.set(symbol, [])
    }

    const history = this.priceHistory.get(symbol)!
    history.push(price)

    // Keep only last 1000 entries
    if (history.length > 1000) {
      history.splice(0, history.length - 1000)
    }
  }

  /**
   * Check and trigger price alerts
   */
  private checkAlerts(symbol: string, currentPrice: PriceData, _previousPrice?: PriceData): void {
    const symbolAlerts = this.getAlertsForSymbol(symbol).filter(alert => 
      alert.isActive && !alert.triggered
    )

    for (const alert of symbolAlerts) {
      let shouldTrigger = false

      switch (alert.type) {
        case 'above':
          shouldTrigger = currentPrice.price >= alert.value
          break
        case 'below':
          shouldTrigger = currentPrice.price <= alert.value
          break
        case 'change_percent':
          shouldTrigger = Math.abs(currentPrice.changePercent24h) >= alert.value
          break
      }

      if (shouldTrigger) {
        alert.triggered = true
        alert.triggeredAt = Date.now()
        this.triggerAlert(alert, currentPrice)
      }
    }
  }

  /**
   * Trigger alert notification
   */
  private triggerAlert(alert: PriceAlert, price: PriceData): void {
    // Emit alert event
    const event = new CustomEvent('priceAlert', {
      detail: { alert, price }
    })
    window.dispatchEvent(event)

    console.log(`Price alert triggered: ${alert.symbol} ${alert.type} ${alert.value}`)
  }

  /**
   * Notify subscribers of price updates
   */
  private notifySubscribers(symbol: string, price: PriceData): void {
    const callbacks = this.subscribers.get(symbol)
    if (callbacks) {
      callbacks.forEach(callback => {
        try {
          callback(price)
        } catch (error) {
          console.error('Error in price update callback:', error)
        }
      })
    }
  }

  /**
   * Get statistics
   */
  getStats(): {
    activeSources: number
    cachedPrices: number
    activeAlerts: number
    subscribers: number
    isRunning: boolean
  } {
    return {
      activeSources: Array.from(this.dataSources.values()).filter(s => s.isActive).length,
      cachedPrices: this.priceCache.size,
      activeAlerts: Array.from(this.alerts.values()).filter(a => a.isActive && !a.triggered).length,
      subscribers: this.subscribers.size,
      isRunning: this.isRunning
    }
  }
}

// Placeholder data source implementations
class CoinGeckoDataSource implements DataSource {
  name = 'coingecko'
  isActive = false
  priority = 3
  lastUpdate = 0
  errorCount = 0

  async connect(): Promise<void> {
    this.isActive = true
  }

  disconnect(): void {
    this.isActive = false
  }

  subscribe(symbols: string[]): void {
    // WebSocket subscription would go here
  }

  unsubscribe(symbols: string[]): void {
    // WebSocket unsubscription would go here
  }

  async getPrices(symbols: string[]): Promise<PriceData[]> {
    // Mock implementation - in real app, call CoinGecko API
    return symbols.map(symbol => ({
      symbol: symbol.toUpperCase(),
      price: Math.random() * 50000,
      change24h: (Math.random() - 0.5) * 1000,
      changePercent24h: (Math.random() - 0.5) * 10,
      volume24h: Math.random() * 1000000000,
      marketCap: Math.random() * 100000000000,
      high24h: Math.random() * 55000,
      low24h: Math.random() * 45000,
      timestamp: Date.now(),
      source: this.name
    }))
  }
}

class CoinMarketCapDataSource implements DataSource {
  name = 'coinmarketcap'
  isActive = false
  priority = 2
  lastUpdate = 0
  errorCount = 0

  async connect(): Promise<void> {
    this.isActive = true
  }

  disconnect(): void {
    this.isActive = false
  }

  subscribe(symbols: string[]): void {
    // WebSocket subscription would go here
  }

  unsubscribe(symbols: string[]): void {
    // WebSocket unsubscription would go here
  }

  async getPrices(symbols: string[]): Promise<PriceData[]> {
    // Mock implementation - in real app, call CoinMarketCap API
    return symbols.map(symbol => ({
      symbol: symbol.toUpperCase(),
      price: Math.random() * 50000,
      change24h: (Math.random() - 0.5) * 1000,
      changePercent24h: (Math.random() - 0.5) * 10,
      volume24h: Math.random() * 1000000000,
      marketCap: Math.random() * 100000000000,
      high24h: Math.random() * 55000,
      low24h: Math.random() * 45000,
      timestamp: Date.now(),
      source: this.name
    }))
  }
}

class BinanceDataSource implements DataSource {
  name = 'binance'
  isActive = false
  priority = 1
  lastUpdate = 0
  errorCount = 0

  async connect(): Promise<void> {
    this.isActive = true
  }

  disconnect(): void {
    this.isActive = false
  }

  subscribe(symbols: string[]): void {
    // WebSocket subscription would go here
  }

  unsubscribe(symbols: string[]): void {
    // WebSocket unsubscription would go here
  }

  async getPrices(symbols: string[]): Promise<PriceData[]> {
    // Mock implementation - in real app, call Binance API
    return symbols.map(symbol => ({
      symbol: symbol.toUpperCase(),
      price: Math.random() * 50000,
      change24h: (Math.random() - 0.5) * 1000,
      changePercent24h: (Math.random() - 0.5) * 10,
      volume24h: Math.random() * 1000000000,
      marketCap: Math.random() * 100000000000,
      high24h: Math.random() * 55000,
      low24h: Math.random() * 45000,
      timestamp: Date.now(),
      source: this.name
    }))
  }
}

// Export singleton instance
export const priceFeedManager = PriceFeedManager.getInstance()
