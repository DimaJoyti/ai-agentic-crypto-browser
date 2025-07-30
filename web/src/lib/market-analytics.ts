import { type PriceData } from './price-feed-manager'

export interface MarketIndicator {
  name: string
  value: number
  signal: 'bullish' | 'bearish' | 'neutral'
  strength: number // 0-100
  description: string
}

export interface TechnicalIndicators {
  rsi: number
  macd: {
    macd: number
    signal: number
    histogram: number
  }
  movingAverages: {
    sma20: number
    sma50: number
    ema12: number
    ema26: number
  }
  bollinger: {
    upper: number
    middle: number
    lower: number
  }
  stochastic: {
    k: number
    d: number
  }
}

export interface VolumeAnalysis {
  averageVolume: number
  volumeRatio: number
  volumeTrend: 'increasing' | 'decreasing' | 'stable'
  volumeProfile: {
    high: number
    medium: number
    low: number
  }
  onBalanceVolume: number
}

export interface MarketSentiment {
  overall: 'extremely_bullish' | 'bullish' | 'neutral' | 'bearish' | 'extremely_bearish'
  score: number // -100 to 100
  fearGreedIndex: number // 0-100
  socialSentiment: number // -100 to 100
  newssentiment: number // -100 to 100
  technicalSentiment: number // -100 to 100
  confidence: number // 0-100
}

export interface MarketMetrics {
  volatility: number
  correlation: Map<string, number>
  beta: number
  sharpeRatio: number
  maxDrawdown: number
  supportLevels: number[]
  resistanceLevels: number[]
}

export interface CandlestickData {
  timestamp: number
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export class MarketAnalytics {
  private static instance: MarketAnalytics
  private priceHistory = new Map<string, PriceData[]>()
  private candlestickData = new Map<string, CandlestickData[]>()
  private indicators = new Map<string, TechnicalIndicators>()
  private volumeAnalysis = new Map<string, VolumeAnalysis>()
  private sentimentData = new Map<string, MarketSentiment>()
  private marketMetrics = new Map<string, MarketMetrics>()

  private constructor() {}

  static getInstance(): MarketAnalytics {
    if (!MarketAnalytics.instance) {
      MarketAnalytics.instance = new MarketAnalytics()
    }
    return MarketAnalytics.instance
  }

  /**
   * Add price data for analysis
   */
  addPriceData(symbol: string, priceData: PriceData): void {
    if (!this.priceHistory.has(symbol)) {
      this.priceHistory.set(symbol, [])
    }

    const history = this.priceHistory.get(symbol)!
    history.push(priceData)

    // Keep only last 1000 entries
    if (history.length > 1000) {
      history.splice(0, history.length - 1000)
    }

    // Update candlestick data
    this.updateCandlestickData(symbol, priceData)

    // Recalculate indicators if we have enough data
    if (history.length >= 50) {
      this.calculateTechnicalIndicators(symbol)
      this.calculateVolumeAnalysis(symbol)
      this.calculateMarketSentiment(symbol)
      this.calculateMarketMetrics(symbol)
    }
  }

  /**
   * Get technical indicators for a symbol
   */
  getTechnicalIndicators(symbol: string): TechnicalIndicators | null {
    return this.indicators.get(symbol) || null
  }

  /**
   * Get volume analysis for a symbol
   */
  getVolumeAnalysis(symbol: string): VolumeAnalysis | null {
    return this.volumeAnalysis.get(symbol) || null
  }

  /**
   * Get market sentiment for a symbol
   */
  getMarketSentiment(symbol: string): MarketSentiment | null {
    return this.sentimentData.get(symbol) || null
  }

  /**
   * Get market metrics for a symbol
   */
  getMarketMetrics(symbol: string): MarketMetrics | null {
    return this.marketMetrics.get(symbol) || null
  }

  /**
   * Get candlestick data for charting
   */
  getCandlestickData(symbol: string, limit = 100): CandlestickData[] {
    const data = this.candlestickData.get(symbol) || []
    return data.slice(-limit)
  }

  /**
   * Get market indicators summary
   */
  getMarketIndicators(symbol: string): MarketIndicator[] {
    const indicators = this.getTechnicalIndicators(symbol)
    const sentiment = this.getMarketSentiment(symbol)
    const volume = this.getVolumeAnalysis(symbol)

    if (!indicators || !sentiment || !volume) return []

    const marketIndicators: MarketIndicator[] = []

    // RSI Indicator
    marketIndicators.push({
      name: 'RSI',
      value: indicators.rsi,
      signal: indicators.rsi > 70 ? 'bearish' : indicators.rsi < 30 ? 'bullish' : 'neutral',
      strength: Math.abs(indicators.rsi - 50) * 2,
      description: `RSI at ${indicators.rsi.toFixed(2)} - ${indicators.rsi > 70 ? 'Overbought' : indicators.rsi < 30 ? 'Oversold' : 'Neutral'}`
    })

    // MACD Indicator
    const macdSignal = indicators.macd.macd > indicators.macd.signal ? 'bullish' : 'bearish'
    marketIndicators.push({
      name: 'MACD',
      value: indicators.macd.histogram,
      signal: macdSignal,
      strength: Math.abs(indicators.macd.histogram) * 100,
      description: `MACD ${macdSignal} crossover with histogram at ${indicators.macd.histogram.toFixed(4)}`
    })

    // Moving Average Indicator
    const price = this.priceHistory.get(symbol)?.slice(-1)[0]?.price || 0
    const maSignal = price > indicators.movingAverages.sma20 ? 'bullish' : 'bearish'
    marketIndicators.push({
      name: 'Moving Average',
      value: indicators.movingAverages.sma20,
      signal: maSignal,
      strength: Math.abs((price - indicators.movingAverages.sma20) / indicators.movingAverages.sma20) * 100,
      description: `Price ${maSignal === 'bullish' ? 'above' : 'below'} 20-day SMA`
    })

    // Volume Indicator
    marketIndicators.push({
      name: 'Volume',
      value: volume.volumeRatio,
      signal: volume.volumeRatio > 1.5 ? 'bullish' : volume.volumeRatio < 0.5 ? 'bearish' : 'neutral',
      strength: Math.abs(volume.volumeRatio - 1) * 100,
      description: `Volume ${volume.volumeTrend} - ${(volume.volumeRatio * 100).toFixed(0)}% of average`
    })

    // Sentiment Indicator
    marketIndicators.push({
      name: 'Market Sentiment',
      value: sentiment.score,
      signal: sentiment.score > 20 ? 'bullish' : sentiment.score < -20 ? 'bearish' : 'neutral',
      strength: Math.abs(sentiment.score),
      description: `Overall sentiment: ${sentiment.overall.replace('_', ' ')}`
    })

    return marketIndicators
  }

  /**
   * Update candlestick data from price data
   */
  private updateCandlestickData(symbol: string, priceData: PriceData): void {
    if (!this.candlestickData.has(symbol)) {
      this.candlestickData.set(symbol, [])
    }

    const candlesticks = this.candlestickData.get(symbol)!
    const timestamp = Math.floor(priceData.timestamp / (5 * 60 * 1000)) * (5 * 60 * 1000) // 5-minute intervals

    // Find existing candlestick for this timestamp
    let candlestick = candlesticks.find(c => c.timestamp === timestamp)

    if (!candlestick) {
      // Create new candlestick
      candlestick = {
        timestamp,
        open: priceData.price,
        high: priceData.price,
        low: priceData.price,
        close: priceData.price,
        volume: priceData.volume24h
      }
      candlesticks.push(candlestick)
    } else {
      // Update existing candlestick
      candlestick.high = Math.max(candlestick.high, priceData.price)
      candlestick.low = Math.min(candlestick.low, priceData.price)
      candlestick.close = priceData.price
      candlestick.volume = priceData.volume24h
    }

    // Keep only last 500 candlesticks
    if (candlesticks.length > 500) {
      candlesticks.splice(0, candlesticks.length - 500)
    }
  }

  /**
   * Calculate technical indicators
   */
  private calculateTechnicalIndicators(symbol: string): void {
    const history = this.priceHistory.get(symbol)
    if (!history || history.length < 50) return

    const prices = history.map(h => h.price)
    const volumes = history.map(h => h.volume24h)

    // Calculate RSI
    const rsi = this.calculateRSI(prices, 14)

    // Calculate MACD
    const macd = this.calculateMACD(prices)

    // Calculate Moving Averages
    const sma20 = this.calculateSMA(prices, 20)
    const sma50 = this.calculateSMA(prices, 50)
    const ema12 = this.calculateEMA(prices, 12)
    const ema26 = this.calculateEMA(prices, 26)

    // Calculate Bollinger Bands
    const bollinger = this.calculateBollingerBands(prices, 20, 2)

    // Calculate Stochastic
    const highs = history.map(h => h.high24h)
    const lows = history.map(h => h.low24h)
    const stochastic = this.calculateStochastic(prices, highs, lows, 14)

    this.indicators.set(symbol, {
      rsi,
      macd,
      movingAverages: { sma20, sma50, ema12, ema26 },
      bollinger,
      stochastic
    })
  }

  /**
   * Calculate volume analysis
   */
  private calculateVolumeAnalysis(symbol: string): void {
    const history = this.priceHistory.get(symbol)
    if (!history || history.length < 20) return

    const volumes = history.map(h => h.volume24h)
    const recentVolumes = volumes.slice(-20)
    const currentVolume = volumes[volumes.length - 1]

    const averageVolume = recentVolumes.reduce((sum, v) => sum + v, 0) / recentVolumes.length
    const volumeRatio = currentVolume / averageVolume

    // Determine volume trend
    const oldAverage = volumes.slice(-40, -20).reduce((sum, v) => sum + v, 0) / 20
    let volumeTrend: 'increasing' | 'decreasing' | 'stable' = 'stable'
    
    if (averageVolume > oldAverage * 1.1) volumeTrend = 'increasing'
    else if (averageVolume < oldAverage * 0.9) volumeTrend = 'decreasing'

    // Calculate volume profile
    const sortedVolumes = [...recentVolumes].sort((a, b) => b - a)
    const volumeProfile = {
      high: sortedVolumes.slice(0, 7).reduce((sum, v) => sum + v, 0) / 7,
      medium: sortedVolumes.slice(7, 14).reduce((sum, v) => sum + v, 0) / 7,
      low: sortedVolumes.slice(14).reduce((sum, v) => sum + v, 0) / 6
    }

    // Calculate On-Balance Volume (simplified)
    const onBalanceVolume = this.calculateOBV(history)

    this.volumeAnalysis.set(symbol, {
      averageVolume,
      volumeRatio,
      volumeTrend,
      volumeProfile,
      onBalanceVolume
    })
  }

  /**
   * Calculate market sentiment
   */
  private calculateMarketSentiment(symbol: string): void {
    const indicators = this.getTechnicalIndicators(symbol)
    const volume = this.getVolumeAnalysis(symbol)
    const history = this.priceHistory.get(symbol)

    if (!indicators || !volume || !history) return

    // Technical sentiment based on indicators
    let technicalSentiment = 0
    
    // RSI contribution
    if (indicators.rsi > 70) technicalSentiment -= 20
    else if (indicators.rsi < 30) technicalSentiment += 20
    else technicalSentiment += (50 - indicators.rsi) * 0.4

    // MACD contribution
    if (indicators.macd.macd > indicators.macd.signal) technicalSentiment += 15
    else technicalSentiment -= 15

    // Moving average contribution
    const currentPrice = history[history.length - 1].price
    if (currentPrice > indicators.movingAverages.sma20) technicalSentiment += 10
    if (currentPrice > indicators.movingAverages.sma50) technicalSentiment += 10

    // Volume contribution
    if (volume.volumeRatio > 1.5) technicalSentiment += 10
    else if (volume.volumeRatio < 0.5) technicalSentiment -= 10

    // Price momentum contribution
    const priceChange = history[history.length - 1].changePercent24h
    technicalSentiment += Math.min(Math.max(priceChange, -30), 30)

    // Normalize to -100 to 100
    technicalSentiment = Math.min(Math.max(technicalSentiment, -100), 100)

    // Mock social and news sentiment (in real app, would come from APIs)
    const socialSentiment = (Math.random() - 0.5) * 100
    const newssentiment = (Math.random() - 0.5) * 80

    // Calculate overall sentiment
    const overallScore = (technicalSentiment * 0.4 + socialSentiment * 0.3 + newssentiment * 0.3)
    
    let overall: MarketSentiment['overall'] = 'neutral'
    if (overallScore > 40) overall = 'extremely_bullish'
    else if (overallScore > 15) overall = 'bullish'
    else if (overallScore < -40) overall = 'extremely_bearish'
    else if (overallScore < -15) overall = 'bearish'

    // Mock fear & greed index
    const fearGreedIndex = Math.min(Math.max(50 + overallScore * 0.5, 0), 100)

    this.sentimentData.set(symbol, {
      overall,
      score: overallScore,
      fearGreedIndex,
      socialSentiment,
      newssentiment,
      technicalSentiment,
      confidence: Math.min(history.length * 2, 100)
    })
  }

  /**
   * Calculate market metrics
   */
  private calculateMarketMetrics(symbol: string): void {
    const history = this.priceHistory.get(symbol)
    if (!history || history.length < 30) return

    const prices = history.map(h => h.price)
    const returns = this.calculateReturns(prices)

    // Calculate volatility (standard deviation of returns)
    const volatility = this.calculateStandardDeviation(returns)

    // Calculate support and resistance levels
    const supportLevels = this.findSupportLevels(prices)
    const resistanceLevels = this.findResistanceLevels(prices)

    // Mock correlation, beta, Sharpe ratio, and max drawdown
    const correlation = new Map([
      ['BTC', Math.random() * 0.8 + 0.1],
      ['ETH', Math.random() * 0.7 + 0.2],
      ['SPY', Math.random() * 0.4 + 0.1]
    ])

    const beta = Math.random() * 2 + 0.5
    const sharpeRatio = (Math.random() - 0.5) * 4
    const maxDrawdown = Math.random() * 0.5

    this.marketMetrics.set(symbol, {
      volatility,
      correlation,
      beta,
      sharpeRatio,
      maxDrawdown,
      supportLevels,
      resistanceLevels
    })
  }

  // Technical indicator calculation methods
  private calculateRSI(prices: number[], period: number): number {
    if (prices.length < period + 1) return 50

    const changes = prices.slice(1).map((price, i) => price - prices[i])
    const gains = changes.map(change => change > 0 ? change : 0)
    const losses = changes.map(change => change < 0 ? -change : 0)

    const avgGain = gains.slice(-period).reduce((sum, gain) => sum + gain, 0) / period
    const avgLoss = losses.slice(-period).reduce((sum, loss) => sum + loss, 0) / period

    if (avgLoss === 0) return 100
    const rs = avgGain / avgLoss
    return 100 - (100 / (1 + rs))
  }

  private calculateMACD(prices: number[]): { macd: number; signal: number; histogram: number } {
    const ema12 = this.calculateEMA(prices, 12)
    const ema26 = this.calculateEMA(prices, 26)
    const macd = ema12 - ema26

    // For simplicity, using a mock signal line
    const signal = macd * 0.9
    const histogram = macd - signal

    return { macd, signal, histogram }
  }

  private calculateSMA(prices: number[], period: number): number {
    if (prices.length < period) return prices[prices.length - 1]
    const recentPrices = prices.slice(-period)
    return recentPrices.reduce((sum, price) => sum + price, 0) / period
  }

  private calculateEMA(prices: number[], period: number): number {
    if (prices.length < period) return prices[prices.length - 1]
    
    const multiplier = 2 / (period + 1)
    let ema = prices[0]
    
    for (let i = 1; i < prices.length; i++) {
      ema = (prices[i] * multiplier) + (ema * (1 - multiplier))
    }
    
    return ema
  }

  private calculateBollingerBands(prices: number[], period: number, stdDev: number): { upper: number; middle: number; lower: number } {
    const sma = this.calculateSMA(prices, period)
    const recentPrices = prices.slice(-period)
    const variance = recentPrices.reduce((sum, price) => sum + Math.pow(price - sma, 2), 0) / period
    const standardDeviation = Math.sqrt(variance)

    return {
      upper: sma + (standardDeviation * stdDev),
      middle: sma,
      lower: sma - (standardDeviation * stdDev)
    }
  }

  private calculateStochastic(prices: number[], highs: number[], lows: number[], period: number): { k: number; d: number } {
    if (prices.length < period) return { k: 50, d: 50 }

    const recentPrices = prices.slice(-period)
    const recentHighs = highs.slice(-period)
    const recentLows = lows.slice(-period)

    const currentPrice = prices[prices.length - 1]
    const highestHigh = Math.max(...recentHighs)
    const lowestLow = Math.min(...recentLows)

    const k = ((currentPrice - lowestLow) / (highestHigh - lowestLow)) * 100
    const d = k * 0.9 // Simplified D calculation

    return { k, d }
  }

  private calculateOBV(history: PriceData[]): number {
    let obv = 0
    for (let i = 1; i < history.length; i++) {
      if (history[i].price > history[i - 1].price) {
        obv += history[i].volume24h
      } else if (history[i].price < history[i - 1].price) {
        obv -= history[i].volume24h
      }
    }
    return obv
  }

  private calculateReturns(prices: number[]): number[] {
    const returns: number[] = []
    for (let i = 1; i < prices.length; i++) {
      returns.push((prices[i] - prices[i - 1]) / prices[i - 1])
    }
    return returns
  }

  private calculateStandardDeviation(values: number[]): number {
    const mean = values.reduce((sum, value) => sum + value, 0) / values.length
    const variance = values.reduce((sum, value) => sum + Math.pow(value - mean, 2), 0) / values.length
    return Math.sqrt(variance)
  }

  private findSupportLevels(prices: number[]): number[] {
    const levels: number[] = []
    const recentPrices = prices.slice(-50)
    
    // Find local minima as support levels
    for (let i = 2; i < recentPrices.length - 2; i++) {
      if (recentPrices[i] < recentPrices[i - 1] && 
          recentPrices[i] < recentPrices[i + 1] &&
          recentPrices[i] < recentPrices[i - 2] && 
          recentPrices[i] < recentPrices[i + 2]) {
        levels.push(recentPrices[i])
      }
    }
    
    return levels.slice(-3) // Return last 3 support levels
  }

  private findResistanceLevels(prices: number[]): number[] {
    const levels: number[] = []
    const recentPrices = prices.slice(-50)
    
    // Find local maxima as resistance levels
    for (let i = 2; i < recentPrices.length - 2; i++) {
      if (recentPrices[i] > recentPrices[i - 1] && 
          recentPrices[i] > recentPrices[i + 1] &&
          recentPrices[i] > recentPrices[i - 2] && 
          recentPrices[i] > recentPrices[i + 2]) {
        levels.push(recentPrices[i])
      }
    }
    
    return levels.slice(-3) // Return last 3 resistance levels
  }
}

// Export singleton instance
export const marketAnalytics = MarketAnalytics.getInstance()
