import { type PriceData } from './price-feed-manager'

export interface HistoricalCandle {
  timestamp: number
  open: number
  high: number
  low: number
  close: number
  volume: number
  symbol: string
  timeframe: '1m' | '5m' | '15m' | '1h' | '4h' | '1d' | '1w' | '1M'
}

export interface HistoricalDataRange {
  symbol: string
  timeframe: string
  startTime: number
  endTime: number
  candles: HistoricalCandle[]
}

export interface DataSource {
  id: string
  name: string
  type: 'api' | 'file' | 'websocket'
  url?: string
  apiKey?: string
  rateLimit: number // requests per minute
  supportedSymbols: string[]
  supportedTimeframes: string[]
  isActive: boolean
}

export interface DataRequest {
  id: string
  symbol: string
  timeframe: string
  startTime: number
  endTime: number
  status: 'pending' | 'processing' | 'completed' | 'failed'
  progress: number
  error?: string
  createdAt: number
  completedAt?: number
}

export interface DataStats {
  totalCandles: number
  symbols: number
  timeframes: number
  oldestData: number
  newestData: number
  dataGaps: DataGap[]
  storageSize: number
}

export interface DataGap {
  symbol: string
  timeframe: string
  startTime: number
  endTime: number
  expectedCandles: number
  missingCandles: number
}

export class HistoricalDataManager {
  private static instance: HistoricalDataManager
  private data = new Map<string, Map<string, HistoricalCandle[]>>() // symbol -> timeframe -> candles
  private dataSources = new Map<string, DataSource>()
  private requests = new Map<string, DataRequest>()
  private cache = new Map<string, HistoricalDataRange>()
  private maxCacheSize = 1000 // Maximum cached ranges

  private constructor() {
    this.initializeDefaultSources()
  }

  static getInstance(): HistoricalDataManager {
    if (!HistoricalDataManager.instance) {
      HistoricalDataManager.instance = new HistoricalDataManager()
    }
    return HistoricalDataManager.instance
  }

  /**
   * Get historical data for symbol and timeframe
   */
  async getHistoricalData(
    symbol: string,
    timeframe: string,
    startTime: number,
    endTime: number,
    forceRefresh = false
  ): Promise<HistoricalCandle[]> {
    const cacheKey = `${symbol}_${timeframe}_${startTime}_${endTime}`
    
    // Check cache first
    if (!forceRefresh && this.cache.has(cacheKey)) {
      const cached = this.cache.get(cacheKey)!
      return cached.candles
    }

    // Check if data exists in storage
    const existingData = this.getStoredData(symbol, timeframe, startTime, endTime)
    if (existingData.length > 0 && !forceRefresh) {
      this.updateCache(cacheKey, { symbol, timeframe, startTime, endTime, candles: existingData })
      return existingData
    }

    // Fetch new data
    return this.fetchHistoricalData(symbol, timeframe, startTime, endTime)
  }

  /**
   * Fetch historical data from external sources
   */
  private async fetchHistoricalData(
    symbol: string,
    timeframe: string,
    startTime: number,
    endTime: number
  ): Promise<HistoricalCandle[]> {
    const requestId = `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    
    const request: DataRequest = {
      id: requestId,
      symbol,
      timeframe,
      startTime,
      endTime,
      status: 'pending',
      progress: 0,
      createdAt: Date.now()
    }

    this.requests.set(requestId, request)

    try {
      request.status = 'processing'
      
      // Find suitable data source
      const source = this.findBestDataSource(symbol, timeframe)
      if (!source) {
        throw new Error(`No data source available for ${symbol} ${timeframe}`)
      }

      // Simulate data fetching (in real implementation, would call actual APIs)
      const candles = await this.simulateDataFetch(symbol, timeframe, startTime, endTime, request)
      
      // Store the data
      this.storeData(symbol, timeframe, candles)
      
      request.status = 'completed'
      request.progress = 100
      request.completedAt = Date.now()

      // Cache the result
      const cacheKey = `${symbol}_${timeframe}_${startTime}_${endTime}`
      this.updateCache(cacheKey, { symbol, timeframe, startTime, endTime, candles })

      return candles
    } catch (error) {
      request.status = 'failed'
      request.error = (error as Error).message
      throw error
    }
  }

  /**
   * Simulate data fetching (replace with real API calls)
   */
  private async simulateDataFetch(
    symbol: string,
    timeframe: string,
    startTime: number,
    endTime: number,
    request: DataRequest
  ): Promise<HistoricalCandle[]> {
    const candles: HistoricalCandle[] = []
    const timeframeMs = this.getTimeframeMs(timeframe)
    
    let currentTime = startTime
    let basePrice = 50000 + Math.random() * 50000 // Random base price
    
    const totalCandles = Math.floor((endTime - startTime) / timeframeMs)
    let processedCandles = 0

    while (currentTime < endTime) {
      // Simulate price movement
      const volatility = 0.02 // 2% volatility
      const change = (Math.random() - 0.5) * volatility
      const open = basePrice
      const close = basePrice * (1 + change)
      const high = Math.max(open, close) * (1 + Math.random() * 0.01)
      const low = Math.min(open, close) * (1 - Math.random() * 0.01)
      const volume = Math.random() * 1000000 + 100000

      candles.push({
        timestamp: currentTime,
        open,
        high,
        low,
        close,
        volume,
        symbol,
        timeframe: timeframe as any
      })

      basePrice = close
      currentTime += timeframeMs
      processedCandles++

      // Update progress
      request.progress = Math.floor((processedCandles / totalCandles) * 100)

      // Simulate processing delay
      if (processedCandles % 100 === 0) {
        await new Promise(resolve => setTimeout(resolve, 10))
      }
    }

    return candles
  }

  /**
   * Get stored data from memory
   */
  private getStoredData(
    symbol: string,
    timeframe: string,
    startTime: number,
    endTime: number
  ): HistoricalCandle[] {
    const symbolData = this.data.get(symbol)
    if (!symbolData) return []

    const timeframeData = symbolData.get(timeframe)
    if (!timeframeData) return []

    return timeframeData.filter(candle => 
      candle.timestamp >= startTime && candle.timestamp <= endTime
    )
  }

  /**
   * Store data in memory
   */
  private storeData(symbol: string, timeframe: string, candles: HistoricalCandle[]): void {
    if (!this.data.has(symbol)) {
      this.data.set(symbol, new Map())
    }

    const symbolData = this.data.get(symbol)!
    if (!symbolData.has(timeframe)) {
      symbolData.set(timeframe, [])
    }

    const timeframeData = symbolData.get(timeframe)!
    
    // Merge new candles with existing data
    const mergedCandles = [...timeframeData, ...candles]
    
    // Remove duplicates and sort by timestamp
    const uniqueCandles = mergedCandles
      .filter((candle, index, arr) => 
        arr.findIndex(c => c.timestamp === candle.timestamp) === index
      )
      .sort((a, b) => a.timestamp - b.timestamp)

    symbolData.set(timeframe, uniqueCandles)

    // Limit data size (keep last 10000 candles per timeframe)
    if (uniqueCandles.length > 10000) {
      symbolData.set(timeframe, uniqueCandles.slice(-10000))
    }
  }

  /**
   * Find best data source for symbol and timeframe
   */
  private findBestDataSource(symbol: string, timeframe: string): DataSource | null {
    const sources = Array.from(this.dataSources.values())
      .filter(source => 
        source.isActive &&
        source.supportedSymbols.includes(symbol) &&
        source.supportedTimeframes.includes(timeframe)
      )
      .sort((a, b) => b.rateLimit - a.rateLimit) // Prefer higher rate limits

    return sources[0] || null
  }

  /**
   * Get timeframe in milliseconds
   */
  private getTimeframeMs(timeframe: string): number {
    const timeframes: Record<string, number> = {
      '1m': 60 * 1000,
      '5m': 5 * 60 * 1000,
      '15m': 15 * 60 * 1000,
      '1h': 60 * 60 * 1000,
      '4h': 4 * 60 * 60 * 1000,
      '1d': 24 * 60 * 60 * 1000,
      '1w': 7 * 24 * 60 * 60 * 1000,
      '1M': 30 * 24 * 60 * 60 * 1000
    }
    return timeframes[timeframe] || 60 * 60 * 1000
  }

  /**
   * Update cache with size limit
   */
  private updateCache(key: string, data: HistoricalDataRange): void {
    if (this.cache.size >= this.maxCacheSize) {
      // Remove oldest cache entry
      const firstKey = this.cache.keys().next().value
      if (firstKey) {
        this.cache.delete(firstKey)
      }
    }
    this.cache.set(key, data)
  }

  /**
   * Get latest candle for symbol and timeframe
   */
  getLatestCandle(symbol: string, timeframe: string): HistoricalCandle | null {
    const symbolData = this.data.get(symbol)
    if (!symbolData) return null

    const timeframeData = symbolData.get(timeframe)
    if (!timeframeData || timeframeData.length === 0) return null

    return timeframeData[timeframeData.length - 1]
  }

  /**
   * Add real-time candle data
   */
  addRealtimeCandle(candle: HistoricalCandle): void {
    this.storeData(candle.symbol, candle.timeframe, [candle])
  }

  /**
   * Convert price data to candle
   */
  priceDataToCandle(priceData: PriceData, timeframe: string): HistoricalCandle {
    return {
      timestamp: priceData.timestamp,
      open: priceData.price,
      high: priceData.high24h,
      low: priceData.low24h,
      close: priceData.price,
      volume: priceData.volume24h,
      symbol: priceData.symbol,
      timeframe: timeframe as any
    }
  }

  /**
   * Get data statistics
   */
  getDataStats(): DataStats {
    let totalCandles = 0
    const symbols = new Set<string>()
    const timeframes = new Set<string>()
    let oldestData = Date.now()
    let newestData = 0

    for (const [symbol, symbolData] of Array.from(this.data)) {
      symbols.add(symbol)
      
      for (const [timeframe, candles] of Array.from(symbolData)) {
        timeframes.add(timeframe)
        totalCandles += candles.length
        
        if (candles.length > 0) {
          oldestData = Math.min(oldestData, candles[0].timestamp)
          newestData = Math.max(newestData, candles[candles.length - 1].timestamp)
        }
      }
    }

    // Calculate approximate storage size
    const storageSize = totalCandles * 64 // Approximate bytes per candle

    return {
      totalCandles,
      symbols: symbols.size,
      timeframes: timeframes.size,
      oldestData,
      newestData,
      dataGaps: this.findDataGaps(),
      storageSize
    }
  }

  /**
   * Find gaps in historical data
   */
  private findDataGaps(): DataGap[] {
    const gaps: DataGap[] = []

    for (const [symbol, symbolData] of Array.from(this.data)) {
      for (const [timeframe, candles] of Array.from(symbolData)) {
        if (candles.length < 2) continue

        const timeframeMs = this.getTimeframeMs(timeframe)
        
        for (let i = 1; i < candles.length; i++) {
          const expectedTime = candles[i - 1].timestamp + timeframeMs
          const actualTime = candles[i].timestamp
          
          if (actualTime > expectedTime + timeframeMs) {
            const missingCandles = Math.floor((actualTime - expectedTime) / timeframeMs)
            
            gaps.push({
              symbol,
              timeframe,
              startTime: expectedTime,
              endTime: actualTime - timeframeMs,
              expectedCandles: missingCandles,
              missingCandles
            })
          }
        }
      }
    }

    return gaps
  }

  /**
   * Get available symbols
   */
  getAvailableSymbols(): string[] {
    return Array.from(this.data.keys())
  }

  /**
   * Get available timeframes for symbol
   */
  getAvailableTimeframes(symbol: string): string[] {
    const symbolData = this.data.get(symbol)
    return symbolData ? Array.from(symbolData.keys()) : []
  }

  /**
   * Get data requests
   */
  getDataRequests(): DataRequest[] {
    return Array.from(this.requests.values())
      .sort((a, b) => b.createdAt - a.createdAt)
  }

  /**
   * Clear old data
   */
  clearOldData(beforeTimestamp: number): void {
    for (const [symbol, symbolData] of Array.from(this.data)) {
      for (const [timeframe, candles] of Array.from(symbolData)) {
        const filteredCandles = candles.filter((candle: any) => candle.timestamp >= beforeTimestamp)
        symbolData.set(timeframe, filteredCandles)
      }
    }
  }

  /**
   * Export data
   */
  exportData(symbol?: string, timeframe?: string): any {
    if (symbol && timeframe) {
      const symbolData = this.data.get(symbol)
      return symbolData?.get(timeframe) || []
    } else if (symbol) {
      return Object.fromEntries(this.data.get(symbol) || new Map())
    } else {
      const result: any = {}
      for (const [sym, symbolData] of Array.from(this.data)) {
        result[sym] = Object.fromEntries(symbolData)
      }
      return result
    }
  }

  /**
   * Import data
   */
  importData(data: any, symbol?: string, timeframe?: string): void {
    if (symbol && timeframe && Array.isArray(data)) {
      this.storeData(symbol, timeframe, data)
    } else if (symbol && typeof data === 'object') {
      for (const [tf, candles] of Object.entries(data)) {
        if (Array.isArray(candles)) {
          this.storeData(symbol, tf, candles as HistoricalCandle[])
        }
      }
    } else if (typeof data === 'object') {
      for (const [sym, symbolData] of Object.entries(data)) {
        if (typeof symbolData === 'object') {
          for (const [tf, candles] of Object.entries(symbolData as any)) {
            if (Array.isArray(candles)) {
              this.storeData(sym, tf, candles as HistoricalCandle[])
            }
          }
        }
      }
    }
  }

  /**
   * Initialize default data sources
   */
  private initializeDefaultSources(): void {
    // CoinGecko
    this.dataSources.set('coingecko', {
      id: 'coingecko',
      name: 'CoinGecko',
      type: 'api',
      url: 'https://api.coingecko.com/api/v3',
      rateLimit: 50,
      supportedSymbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'],
      supportedTimeframes: ['1h', '4h', '1d'],
      isActive: true
    })

    // Binance
    this.dataSources.set('binance', {
      id: 'binance',
      name: 'Binance',
      type: 'api',
      url: 'https://api.binance.com/api/v3',
      rateLimit: 1200,
      supportedSymbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'],
      supportedTimeframes: ['1m', '5m', '15m', '1h', '4h', '1d'],
      isActive: true
    })

    // CoinMarketCap
    this.dataSources.set('coinmarketcap', {
      id: 'coinmarketcap',
      name: 'CoinMarketCap',
      type: 'api',
      url: 'https://pro-api.coinmarketcap.com/v1',
      rateLimit: 333,
      supportedSymbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
      supportedTimeframes: ['1h', '1d'],
      isActive: true
    })
  }

  /**
   * Add data source
   */
  addDataSource(source: DataSource): void {
    this.dataSources.set(source.id, source)
  }

  /**
   * Get data sources
   */
  getDataSources(): DataSource[] {
    return Array.from(this.dataSources.values())
  }
}

// Export singleton instance
export const historicalDataManager = HistoricalDataManager.getInstance()
