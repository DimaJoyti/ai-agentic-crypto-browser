import { type PriceData } from './price-feed-manager'
import { type TechnicalIndicators } from './market-analytics'

export interface TradingSignal {
  id: string
  symbol: string
  type: 'buy' | 'sell' | 'hold'
  strength: number // 0-100
  confidence: number // 0-100
  price: number
  timestamp: number
  source: 'technical' | 'sentiment' | 'volume' | 'pattern' | 'composite'
  indicators: string[]
  description: string
  targetPrice?: number
  stopLoss?: number
  timeframe: '1m' | '5m' | '15m' | '1h' | '4h' | '1d'
  expiresAt?: number
  metadata?: Record<string, any>
}

export interface SignalStrategy {
  id: string
  name: string
  description: string
  enabled: boolean
  symbols: string[]
  timeframes: string[]
  conditions: SignalCondition[]
  filters: SignalFilter[]
  riskLevel: 'low' | 'medium' | 'high'
  minConfidence: number
  maxSignalsPerDay: number
  cooldownPeriod: number // minutes
  createdAt: number
  lastTriggered?: number
}

export interface SignalCondition {
  id: string
  type: 'rsi' | 'macd' | 'ma_cross' | 'bollinger' | 'volume' | 'price' | 'pattern'
  operator: 'gt' | 'lt' | 'eq' | 'cross_above' | 'cross_below' | 'between'
  value: number | number[]
  weight: number // 0-1
  required: boolean
}

export interface SignalFilter {
  id: string
  type: 'volume_min' | 'market_cap_min' | 'price_range' | 'volatility_max' | 'time_range'
  value: number | number[] | string
  enabled: boolean
}

export interface SignalAlert {
  id: string
  signalId: string
  userId: string
  type: 'email' | 'push' | 'sms' | 'webhook' | 'in_app'
  enabled: boolean
  conditions: AlertCondition[]
  cooldownPeriod: number
  maxAlertsPerDay: number
  createdAt: number
  lastTriggered?: number
}

export interface AlertCondition {
  signalType: 'buy' | 'sell' | 'hold' | 'any'
  minStrength: number
  minConfidence: number
  symbols: string[]
  strategies: string[]
}

export interface SignalPerformance {
  strategyId: string
  totalSignals: number
  successfulSignals: number
  successRate: number
  avgReturn: number
  maxReturn: number
  minReturn: number
  avgHoldTime: number
  profitFactor: number
  sharpeRatio: number
  maxDrawdown: number
  winRate: number
  lossRate: number
  avgWin: number
  avgLoss: number
}

export class TradingSignalsEngine {
  private static instance: TradingSignalsEngine
  private signals = new Map<string, TradingSignal[]>() // symbol -> signals
  private strategies = new Map<string, SignalStrategy>()
  private alerts = new Map<string, SignalAlert[]>() // userId -> alerts
  private performance = new Map<string, SignalPerformance>()
  private priceHistory = new Map<string, PriceData[]>()
  private lastSignalTime = new Map<string, number>() // strategy -> timestamp

  private constructor() {
    this.initializeDefaultStrategies()
  }

  static getInstance(): TradingSignalsEngine {
    if (!TradingSignalsEngine.instance) {
      TradingSignalsEngine.instance = new TradingSignalsEngine()
    }
    return TradingSignalsEngine.instance
  }

  /**
   * Add price data and generate signals
   */
  addPriceData(priceData: PriceData, indicators?: TechnicalIndicators): void {
    // Store price history
    if (!this.priceHistory.has(priceData.symbol)) {
      this.priceHistory.set(priceData.symbol, [])
    }
    
    const history = this.priceHistory.get(priceData.symbol)!
    history.push(priceData)
    
    // Keep only last 1000 entries
    if (history.length > 1000) {
      history.splice(0, history.length - 1000)
    }

    // Generate signals for all enabled strategies
    this.generateSignals(priceData, indicators)
  }

  /**
   * Generate trading signals based on strategies
   */
  private generateSignals(priceData: PriceData, indicators?: TechnicalIndicators): void {
    const enabledStrategies = Array.from(this.strategies.values()).filter(s => s.enabled)

    for (const strategy of enabledStrategies) {
      if (!strategy.symbols.includes(priceData.symbol)) continue

      // Check cooldown period
      const lastTriggered = this.lastSignalTime.get(strategy.id)
      if (lastTriggered && Date.now() - lastTriggered < strategy.cooldownPeriod * 60 * 1000) {
        continue
      }

      // Check daily signal limit
      const todaySignals = this.getTodaySignals(priceData.symbol, strategy.id)
      if (todaySignals >= strategy.maxSignalsPerDay) continue

      // Evaluate strategy conditions
      const signal = this.evaluateStrategy(strategy, priceData, indicators)
      if (signal && signal.confidence >= strategy.minConfidence) {
        this.addSignal(signal)
        this.lastSignalTime.set(strategy.id, Date.now())
        
        // Trigger alerts
        this.triggerAlerts(signal)
      }
    }
  }

  /**
   * Evaluate strategy conditions and generate signal
   */
  private evaluateStrategy(
    strategy: SignalStrategy, 
    priceData: PriceData, 
    indicators?: TechnicalIndicators
  ): TradingSignal | null {
    if (!indicators) return null

    let totalScore = 0
    let totalWeight = 0
    let requiredMet = true
    const triggeredIndicators: string[] = []

    // Evaluate each condition
    for (const condition of strategy.conditions) {
      const result = this.evaluateCondition(condition, priceData, indicators)
      
      if (result.triggered) {
        totalScore += result.score * condition.weight
        triggeredIndicators.push(condition.type)
      } else if (condition.required) {
        requiredMet = false
        break
      }
      
      totalWeight += condition.weight
    }

    if (!requiredMet || totalWeight === 0) return null

    const confidence = (totalScore / totalWeight) * 100
    if (confidence < strategy.minConfidence) return null

    // Determine signal type and strength
    const signalType = this.determineSignalType(strategy, triggeredIndicators, indicators)
    const strength = Math.min(100, confidence * 1.2) // Boost strength slightly

    // Apply filters
    if (!this.passesFilters(strategy.filters, priceData)) return null

    return {
      id: `signal_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      symbol: priceData.symbol,
      type: signalType,
      strength,
      confidence,
      price: priceData.price,
      timestamp: Date.now(),
      source: 'technical',
      indicators: triggeredIndicators,
      description: this.generateSignalDescription(signalType, triggeredIndicators, strength),
      timeframe: '1h', // Default timeframe
      targetPrice: this.calculateTargetPrice(signalType, priceData.price, strength),
      stopLoss: this.calculateStopLoss(signalType, priceData.price, strength)
    }
  }

  /**
   * Evaluate individual condition
   */
  private evaluateCondition(
    condition: SignalCondition,
    priceData: PriceData,
    indicators: TechnicalIndicators
  ): { triggered: boolean; score: number } {
    let value: number
    let triggered = false
    let score = 0

    switch (condition.type) {
      case 'rsi':
        value = indicators.rsi
        triggered = this.checkOperator(value, condition.operator, condition.value)
        score = triggered ? Math.abs(value - 50) / 50 : 0 // Distance from neutral
        break

      case 'macd':
        value = indicators.macd.histogram
        triggered = this.checkOperator(value, condition.operator, condition.value)
        score = triggered ? Math.min(Math.abs(value) * 1000, 1) : 0
        break

      case 'ma_cross':
        const sma20 = indicators.movingAverages.sma20
        const sma50 = indicators.movingAverages.sma50
        triggered = condition.operator === 'cross_above' 
          ? sma20 > sma50 && priceData.price > sma20
          : sma20 < sma50 && priceData.price < sma20
        score = triggered ? Math.abs(sma20 - sma50) / Math.max(sma20, sma50) : 0
        break

      case 'bollinger':
        const { upper, lower } = indicators.bollinger
        const position = (priceData.price - lower) / (upper - lower)
        triggered = this.checkOperator(position, condition.operator, condition.value)
        score = triggered ? Math.abs(position - 0.5) * 2 : 0
        break

      case 'volume':
        value = priceData.volume24h
        triggered = this.checkOperator(value, condition.operator, condition.value)
        score = triggered ? 0.8 : 0 // Volume signals get moderate score
        break

      case 'price':
        value = priceData.price
        triggered = this.checkOperator(value, condition.operator, condition.value)
        score = triggered ? 0.7 : 0
        break

      default:
        return { triggered: false, score: 0 }
    }

    return { triggered, score }
  }

  /**
   * Check operator condition
   */
  private checkOperator(value: number, operator: string, target: number | number[]): boolean {
    switch (operator) {
      case 'gt':
        return value > (target as number)
      case 'lt':
        return value < (target as number)
      case 'eq':
        return Math.abs(value - (target as number)) < 0.01
      case 'between':
        const [min, max] = target as number[]
        return value >= min && value <= max
      default:
        return false
    }
  }

  /**
   * Determine signal type based on conditions
   */
  private determineSignalType(
    strategy: SignalStrategy,
    indicators: string[],
    technicalIndicators: TechnicalIndicators
  ): 'buy' | 'sell' | 'hold' {
    let buyScore = 0
    let sellScore = 0

    // RSI signals
    if (indicators.includes('rsi')) {
      if (technicalIndicators.rsi < 30) buyScore += 0.3
      if (technicalIndicators.rsi > 70) sellScore += 0.3
    }

    // MACD signals
    if (indicators.includes('macd')) {
      if (technicalIndicators.macd.histogram > 0) buyScore += 0.2
      else sellScore += 0.2
    }

    // Moving average signals
    if (indicators.includes('ma_cross')) {
      if (technicalIndicators.movingAverages.sma20 > technicalIndicators.movingAverages.sma50) {
        buyScore += 0.25
      } else {
        sellScore += 0.25
      }
    }

    // Volume confirmation
    if (indicators.includes('volume')) {
      buyScore += 0.1
      sellScore += 0.1
    }

    if (buyScore > sellScore && buyScore > 0.4) return 'buy'
    if (sellScore > buyScore && sellScore > 0.4) return 'sell'
    return 'hold'
  }

  /**
   * Check if signal passes filters
   */
  private passesFilters(filters: SignalFilter[], priceData: PriceData): boolean {
    for (const filter of filters) {
      if (!filter.enabled) continue

      switch (filter.type) {
        case 'volume_min':
          if (priceData.volume24h < (filter.value as number)) return false
          break
        case 'market_cap_min':
          if (priceData.marketCap < (filter.value as number)) return false
          break
        case 'price_range':
          const [minPrice, maxPrice] = filter.value as number[]
          if (priceData.price < minPrice || priceData.price > maxPrice) return false
          break
      }
    }
    return true
  }

  /**
   * Calculate target price
   */
  private calculateTargetPrice(type: 'buy' | 'sell' | 'hold', price: number, strength: number): number {
    const multiplier = strength / 100 * 0.1 // Max 10% target
    
    if (type === 'buy') {
      return price * (1 + multiplier)
    } else if (type === 'sell') {
      return price * (1 - multiplier)
    }
    return price
  }

  /**
   * Calculate stop loss
   */
  private calculateStopLoss(type: 'buy' | 'sell' | 'hold', price: number, strength: number): number {
    const multiplier = (100 - strength) / 100 * 0.05 // Max 5% stop loss
    
    if (type === 'buy') {
      return price * (1 - multiplier)
    } else if (type === 'sell') {
      return price * (1 + multiplier)
    }
    return price
  }

  /**
   * Generate signal description
   */
  private generateSignalDescription(type: string, indicators: string[], strength: number): string {
    const action = type.toUpperCase()
    const strengthText = strength > 80 ? 'Strong' : strength > 60 ? 'Moderate' : 'Weak'
    const indicatorText = indicators.join(', ').toUpperCase()
    
    return `${strengthText} ${action} signal based on ${indicatorText} indicators`
  }

  /**
   * Add signal to storage
   */
  private addSignal(signal: TradingSignal): void {
    if (!this.signals.has(signal.symbol)) {
      this.signals.set(signal.symbol, [])
    }

    const symbolSignals = this.signals.get(signal.symbol)!
    symbolSignals.push(signal)

    // Keep only last 100 signals per symbol
    if (symbolSignals.length > 100) {
      symbolSignals.splice(0, symbolSignals.length - 100)
    }

    // Emit signal event
    const event = new CustomEvent('tradingSignal', { detail: signal })
    window.dispatchEvent(event)
  }

  /**
   * Trigger alerts for signal
   */
  private triggerAlerts(signal: TradingSignal): void {
    for (const [userId, userAlerts] of Array.from(this.alerts)) {
      for (const alert of userAlerts) {
        if (!alert.enabled) continue

        // Check alert conditions
        const shouldTrigger = alert.conditions.some((condition: any) =>
          this.matchesAlertCondition(condition, signal)
        )

        if (shouldTrigger) {
          this.sendAlert(alert, signal)
        }
      }
    }
  }

  /**
   * Check if signal matches alert condition
   */
  private matchesAlertCondition(condition: AlertCondition, signal: TradingSignal): boolean {
    // Check signal type
    if (condition.signalType !== 'any' && condition.signalType !== signal.type) {
      return false
    }

    // Check strength and confidence
    if (signal.strength < condition.minStrength || signal.confidence < condition.minConfidence) {
      return false
    }

    // Check symbols
    if (condition.symbols.length > 0 && !condition.symbols.includes(signal.symbol)) {
      return false
    }

    return true
  }

  /**
   * Send alert notification
   */
  private sendAlert(alert: SignalAlert, signal: TradingSignal): void {
    // Check cooldown
    if (alert.lastTriggered && Date.now() - alert.lastTriggered < alert.cooldownPeriod * 60 * 1000) {
      return
    }

    // Emit alert event
    const event = new CustomEvent('signalAlert', {
      detail: { alert, signal }
    })
    window.dispatchEvent(event)

    // Update last triggered time
    alert.lastTriggered = Date.now()
  }

  /**
   * Get signals for symbol
   */
  getSignals(symbol: string, limit = 50): TradingSignal[] {
    const signals = this.signals.get(symbol) || []
    return signals.slice(-limit).reverse()
  }

  /**
   * Get all recent signals
   */
  getAllRecentSignals(limit = 100): TradingSignal[] {
    const allSignals: TradingSignal[] = []
    
    for (const signals of Array.from(this.signals.values())) {
      allSignals.push(...signals)
    }

    return allSignals
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
  }

  /**
   * Get today's signal count
   */
  private getTodaySignals(symbol: string, strategyId: string): number {
    const today = new Date().toDateString()
    const signals = this.signals.get(symbol) || []
    
    return signals.filter(signal => 
      new Date(signal.timestamp).toDateString() === today &&
      signal.metadata?.strategyId === strategyId
    ).length
  }

  /**
   * Initialize default strategies
   */
  private initializeDefaultStrategies(): void {
    // RSI Oversold/Overbought Strategy
    this.strategies.set('rsi_strategy', {
      id: 'rsi_strategy',
      name: 'RSI Oversold/Overbought',
      description: 'Buy when RSI < 30, sell when RSI > 70',
      enabled: true,
      symbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
      timeframes: ['1h', '4h'],
      conditions: [
        {
          id: 'rsi_oversold',
          type: 'rsi',
          operator: 'lt',
          value: 30,
          weight: 0.6,
          required: false
        },
        {
          id: 'rsi_overbought',
          type: 'rsi',
          operator: 'gt',
          value: 70,
          weight: 0.6,
          required: false
        }
      ],
      filters: [
        {
          id: 'min_volume',
          type: 'volume_min',
          value: 1000000,
          enabled: true
        }
      ],
      riskLevel: 'medium',
      minConfidence: 60,
      maxSignalsPerDay: 5,
      cooldownPeriod: 60,
      createdAt: Date.now()
    })

    // MACD Crossover Strategy
    this.strategies.set('macd_strategy', {
      id: 'macd_strategy',
      name: 'MACD Crossover',
      description: 'Buy/sell signals based on MACD line crossovers',
      enabled: true,
      symbols: ['BTC', 'ETH', 'BNB'],
      timeframes: ['4h', '1d'],
      conditions: [
        {
          id: 'macd_bullish',
          type: 'macd',
          operator: 'gt',
          value: 0,
          weight: 0.7,
          required: false
        }
      ],
      filters: [],
      riskLevel: 'medium',
      minConfidence: 65,
      maxSignalsPerDay: 3,
      cooldownPeriod: 120,
      createdAt: Date.now()
    })
  }

  /**
   * Add custom strategy
   */
  addStrategy(strategy: Omit<SignalStrategy, 'id' | 'createdAt'>): string {
    const id = `strategy_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    const fullStrategy: SignalStrategy = {
      ...strategy,
      id,
      createdAt: Date.now()
    }
    
    this.strategies.set(id, fullStrategy)
    return id
  }

  /**
   * Get all strategies
   */
  getStrategies(): SignalStrategy[] {
    return Array.from(this.strategies.values())
  }

  /**
   * Add alert
   */
  addAlert(userId: string, alert: Omit<SignalAlert, 'id' | 'createdAt'>): string {
    const id = `alert_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    const fullAlert: SignalAlert = {
      ...alert,
      id,
      createdAt: Date.now()
    }

    if (!this.alerts.has(userId)) {
      this.alerts.set(userId, [])
    }

    this.alerts.get(userId)!.push(fullAlert)
    return id
  }

  /**
   * Get user alerts
   */
  getUserAlerts(userId: string): SignalAlert[] {
    return this.alerts.get(userId) || []
  }
}

// Export singleton instance
export const tradingSignalsEngine = TradingSignalsEngine.getInstance()
