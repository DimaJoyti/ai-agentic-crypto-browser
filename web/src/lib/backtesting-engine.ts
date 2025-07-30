import { type HistoricalCandle } from './historical-data'
import { type SignalStrategy, type TradingSignal } from './trading-signals'

export interface BacktestConfig {
  id: string
  name: string
  strategy: SignalStrategy
  symbol: string
  timeframe: string
  startTime: number
  endTime: number
  initialCapital: number
  positionSize: number // Percentage of capital per trade
  maxPositions: number
  commission: number // Percentage
  slippage: number // Percentage
  stopLoss?: number // Percentage
  takeProfit?: number // Percentage
  riskManagement: RiskManagement
}

export interface RiskManagement {
  maxDrawdown: number // Percentage
  maxDailyLoss: number // Percentage
  maxConsecutiveLosses: number
  positionSizing: 'fixed' | 'percentage' | 'kelly' | 'volatility'
  riskPerTrade: number // Percentage of capital
}

export interface BacktestTrade {
  id: string
  signal: TradingSignal
  entryTime: number
  entryPrice: number
  exitTime?: number
  exitPrice?: number
  quantity: number
  side: 'long' | 'short'
  status: 'open' | 'closed'
  pnl: number
  pnlPercent: number
  commission: number
  slippage: number
  exitReason?: 'signal' | 'stop_loss' | 'take_profit' | 'time_limit' | 'max_drawdown'
  holdingPeriod?: number
}

export interface BacktestResults {
  config: BacktestConfig
  trades: BacktestTrade[]
  performance: BacktestPerformance
  equity: EquityPoint[]
  drawdown: DrawdownPoint[]
  monthlyReturns: MonthlyReturn[]
  statistics: BacktestStatistics
  riskMetrics: BacktestRiskMetrics
}

export interface BacktestPerformance {
  totalReturn: number
  totalReturnPercent: number
  annualizedReturn: number
  sharpeRatio: number
  sortinoRatio: number
  calmarRatio: number
  maxDrawdown: number
  maxDrawdownPercent: number
  volatility: number
  winRate: number
  profitFactor: number
  avgWin: number
  avgLoss: number
  avgTrade: number
  totalTrades: number
  winningTrades: number
  losingTrades: number
  largestWin: number
  largestLoss: number
  avgHoldingPeriod: number
  totalCommission: number
  totalSlippage: number
}

export interface EquityPoint {
  timestamp: number
  equity: number
  drawdown: number
  trades: number
}

export interface DrawdownPoint {
  timestamp: number
  drawdown: number
  drawdownPercent: number
  peak: number
  valley: number
}

export interface MonthlyReturn {
  year: number
  month: number
  return: number
  returnPercent: number
  trades: number
}

export interface BacktestStatistics {
  startDate: Date
  endDate: Date
  duration: number // days
  totalCandles: number
  signalsGenerated: number
  signalsTraded: number
  avgTradesPerMonth: number
  bestMonth: MonthlyReturn
  worstMonth: MonthlyReturn
  consecutiveWins: number
  consecutiveLosses: number
  maxConsecutiveWins: number
  maxConsecutiveLosses: number
}

export interface BacktestRiskMetrics {
  var95: number // Value at Risk 95%
  var99: number // Value at Risk 99%
  cvar95: number // Conditional Value at Risk 95%
  beta: number
  alpha: number
  informationRatio: number
  treynorRatio: number
  ulcerIndex: number
  painIndex: number
  gainToPainRatio: number
}

export class BacktestingEngine {
  private static instance: BacktestingEngine
  private runningBacktests = new Map<string, BacktestResults>()
  private completedBacktests = new Map<string, BacktestResults>()

  private constructor() {}

  static getInstance(): BacktestingEngine {
    if (!BacktestingEngine.instance) {
      BacktestingEngine.instance = new BacktestingEngine()
    }
    return BacktestingEngine.instance
  }

  /**
   * Run backtest
   */
  async runBacktest(
    config: BacktestConfig,
    historicalData: HistoricalCandle[],
    onProgress?: (progress: number) => void
  ): Promise<BacktestResults> {
    const trades: BacktestTrade[] = []
    const equity: EquityPoint[] = []
    const signals: TradingSignal[] = []

    let currentCapital = config.initialCapital
    let currentEquity = config.initialCapital
    let openPositions: BacktestTrade[] = []
    let peak = config.initialCapital
    let maxDrawdown = 0
    let consecutiveWins = 0
    let consecutiveLosses = 0
    let maxConsecutiveWins = 0
    let maxConsecutiveLosses = 0

    // Sort data by timestamp
    const sortedData = [...historicalData].sort((a, b) => a.timestamp - b.timestamp)
    const totalCandles = sortedData.length

    for (let i = 0; i < sortedData.length; i++) {
      const candle = sortedData[i]
      const progress = (i / totalCandles) * 100

      if (onProgress) {
        onProgress(progress)
      }

      // Generate signals (simplified - in real implementation would use strategy engine)
      const signal = this.generateSignalForCandle(candle, config.strategy, sortedData.slice(0, i + 1))
      if (signal) {
        signals.push(signal)
      }

      // Process existing positions
      this.processOpenPositions(openPositions, candle, config)

      // Close completed positions
      const closedPositions = openPositions.filter(pos => pos.status === 'closed')
      for (const closedPos of closedPositions) {
        trades.push(closedPos)
        currentCapital += closedPos.pnl
        
        // Update consecutive wins/losses
        if (closedPos.pnl > 0) {
          consecutiveWins++
          consecutiveLosses = 0
          maxConsecutiveWins = Math.max(maxConsecutiveWins, consecutiveWins)
        } else {
          consecutiveLosses++
          consecutiveWins = 0
          maxConsecutiveLosses = Math.max(maxConsecutiveLosses, consecutiveLosses)
        }
      }

      // Remove closed positions
      openPositions = openPositions.filter(pos => pos.status === 'open')

      // Open new positions based on signals
      if (signal && openPositions.length < config.maxPositions) {
        const newTrade = this.createTradeFromSignal(signal, candle, config, currentCapital)
        if (newTrade) {
          openPositions.push(newTrade)
          currentCapital -= newTrade.quantity * newTrade.entryPrice
        }
      }

      // Calculate current equity
      currentEquity = currentCapital + this.calculateOpenPositionsValue(openPositions, candle)

      // Update peak and drawdown
      if (currentEquity > peak) {
        peak = currentEquity
      }
      const currentDrawdown = peak - currentEquity
      const currentDrawdownPercent = (currentDrawdown / peak) * 100
      maxDrawdown = Math.max(maxDrawdown, currentDrawdownPercent)

      // Record equity point
      equity.push({
        timestamp: candle.timestamp,
        equity: currentEquity,
        drawdown: currentDrawdownPercent,
        trades: trades.length
      })

      // Risk management - stop trading if max drawdown exceeded
      if (currentDrawdownPercent > config.riskManagement.maxDrawdown) {
        // Close all positions
        for (const pos of openPositions) {
          this.closePosition(pos, candle, 'max_drawdown')
          trades.push(pos)
        }
        break
      }
    }

    // Close any remaining open positions
    const lastCandle = sortedData[sortedData.length - 1]
    for (const pos of openPositions) {
      this.closePosition(pos, lastCandle, 'time_limit')
      trades.push(pos)
    }

    // Calculate performance metrics
    const performance = this.calculatePerformance(trades, config.initialCapital, equity)
    const drawdown = this.calculateDrawdown(equity)
    const monthlyReturns = this.calculateMonthlyReturns(trades, equity)
    const statistics = this.calculateStatistics(trades, signals, sortedData, monthlyReturns)
    const riskMetrics = this.calculateRiskMetrics(equity, trades)

    const results: BacktestResults = {
      config,
      trades,
      performance,
      equity,
      drawdown,
      monthlyReturns,
      statistics,
      riskMetrics
    }

    this.completedBacktests.set(config.id, results)
    return results
  }

  /**
   * Map candle timeframe to signal timeframe (filter out unsupported timeframes)
   */
  private mapToSignalTimeframe(candleTimeframe: string): '1m' | '5m' | '15m' | '1h' | '4h' | '1d' {
    switch (candleTimeframe) {
      case '1m':
      case '5m':
      case '15m':
      case '1h':
      case '4h':
      case '1d':
        return candleTimeframe as '1m' | '5m' | '15m' | '1h' | '4h' | '1d'
      case '1w':
        return '1d' // Map weekly to daily
      case '1M':
        return '1d' // Map monthly to daily
      default:
        return '1h' // Default fallback
    }
  }

  /**
   * Generate signal for candle (simplified)
   */
  private generateSignalForCandle(
    candle: HistoricalCandle,
    strategy: SignalStrategy,
    historicalData: HistoricalCandle[]
  ): TradingSignal | null {
    // Simplified signal generation - in real implementation would use full strategy engine
    if (historicalData.length < 20) return null

    const prices = historicalData.map(c => c.close)
    const sma20 = this.calculateSMA(prices, 20)
    const rsi = this.calculateRSI(prices, 14)

    // Map candle timeframe to signal timeframe
    const signalTimeframe = this.mapToSignalTimeframe(candle.timeframe)

    // Simple RSI strategy
    if (rsi < 30 && candle.close > sma20) {
      return {
        id: `signal_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
        symbol: candle.symbol,
        type: 'buy',
        strength: 80,
        confidence: 75,
        price: candle.close,
        timestamp: candle.timestamp,
        source: 'technical',
        indicators: ['rsi', 'sma'],
        description: 'RSI oversold with price above SMA20',
        timeframe: signalTimeframe
      }
    } else if (rsi > 70 && candle.close < sma20) {
      return {
        id: `signal_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
        symbol: candle.symbol,
        type: 'sell',
        strength: 80,
        confidence: 75,
        price: candle.close,
        timestamp: candle.timestamp,
        source: 'technical',
        indicators: ['rsi', 'sma'],
        description: 'RSI overbought with price below SMA20',
        timeframe: signalTimeframe
      }
    }

    return null
  }

  /**
   * Create trade from signal
   */
  private createTradeFromSignal(
    signal: TradingSignal,
    candle: HistoricalCandle,
    config: BacktestConfig,
    availableCapital: number
  ): BacktestTrade | null {
    const positionValue = availableCapital * (config.positionSize / 100)
    const entryPrice = candle.close * (1 + config.slippage / 100) // Apply slippage
    const quantity = positionValue / entryPrice
    const commission = positionValue * (config.commission / 100)

    if (positionValue < entryPrice) return null // Not enough capital

    return {
      id: `trade_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
      signal,
      entryTime: candle.timestamp,
      entryPrice,
      quantity,
      side: signal.type === 'buy' ? 'long' : 'short',
      status: 'open',
      pnl: 0,
      pnlPercent: 0,
      commission,
      slippage: positionValue * (config.slippage / 100)
    }
  }

  /**
   * Process open positions
   */
  private processOpenPositions(
    positions: BacktestTrade[],
    candle: HistoricalCandle,
    config: BacktestConfig
  ): void {
    for (const position of positions) {
      if (position.status === 'closed') continue

      const currentPrice = candle.close
      const unrealizedPnl = this.calculateUnrealizedPnL(position, currentPrice)
      const unrealizedPnlPercent = (unrealizedPnl / (position.quantity * position.entryPrice)) * 100

      // Check stop loss
      if (config.stopLoss && unrealizedPnlPercent <= -config.stopLoss) {
        this.closePosition(position, candle, 'stop_loss')
        continue
      }

      // Check take profit
      if (config.takeProfit && unrealizedPnlPercent >= config.takeProfit) {
        this.closePosition(position, candle, 'take_profit')
        continue
      }

      // Check for exit signal (simplified)
      if (position.side === 'long' && currentPrice < position.entryPrice * 0.95) {
        this.closePosition(position, candle, 'signal')
      } else if (position.side === 'short' && currentPrice > position.entryPrice * 1.05) {
        this.closePosition(position, candle, 'signal')
      }
    }
  }

  /**
   * Close position
   */
  private closePosition(
    position: BacktestTrade,
    candle: HistoricalCandle,
    reason: BacktestTrade['exitReason']
  ): void {
    const exitPrice = candle.close * (1 - (position.side === 'long' ? 1 : -1) * 0.001) // Apply slippage
    const exitCommission = position.quantity * exitPrice * 0.001 // Exit commission

    position.exitTime = candle.timestamp
    position.exitPrice = exitPrice
    position.status = 'closed'
    position.exitReason = reason
    position.holdingPeriod = candle.timestamp - position.entryTime

    // Calculate P&L
    if (position.side === 'long') {
      position.pnl = (exitPrice - position.entryPrice) * position.quantity - position.commission - exitCommission
    } else {
      position.pnl = (position.entryPrice - exitPrice) * position.quantity - position.commission - exitCommission
    }

    position.pnlPercent = (position.pnl / (position.quantity * position.entryPrice)) * 100
  }

  /**
   * Calculate unrealized P&L
   */
  private calculateUnrealizedPnL(position: BacktestTrade, currentPrice: number): number {
    if (position.side === 'long') {
      return (currentPrice - position.entryPrice) * position.quantity
    } else {
      return (position.entryPrice - currentPrice) * position.quantity
    }
  }

  /**
   * Calculate open positions value
   */
  private calculateOpenPositionsValue(positions: BacktestTrade[], candle: HistoricalCandle): number {
    return positions.reduce((total, pos) => {
      if (pos.status === 'open') {
        return total + this.calculateUnrealizedPnL(pos, candle.close)
      }
      return total
    }, 0)
  }

  /**
   * Calculate performance metrics
   */
  private calculatePerformance(
    trades: BacktestTrade[],
    initialCapital: number,
    equity: EquityPoint[]
  ): BacktestPerformance {
    const finalEquity = equity[equity.length - 1]?.equity || initialCapital
    const totalReturn = finalEquity - initialCapital
    const totalReturnPercent = (totalReturn / initialCapital) * 100

    const winningTrades = trades.filter(t => t.pnl > 0)
    const losingTrades = trades.filter(t => t.pnl < 0)
    const winRate = trades.length > 0 ? (winningTrades.length / trades.length) * 100 : 0

    const avgWin = winningTrades.length > 0 
      ? winningTrades.reduce((sum, t) => sum + t.pnl, 0) / winningTrades.length 
      : 0
    const avgLoss = losingTrades.length > 0 
      ? Math.abs(losingTrades.reduce((sum, t) => sum + t.pnl, 0) / losingTrades.length)
      : 0

    const profitFactor = avgLoss > 0 ? avgWin / avgLoss : 0
    const avgTrade = trades.length > 0 ? trades.reduce((sum, t) => sum + t.pnl, 0) / trades.length : 0

    // Calculate Sharpe ratio (simplified)
    const returns = equity.map((point, i) => {
      if (i === 0) return 0
      return (point.equity - equity[i - 1].equity) / equity[i - 1].equity
    }).slice(1)

    const avgReturn = returns.reduce((sum, r) => sum + r, 0) / returns.length
    const volatility = Math.sqrt(returns.reduce((sum, r) => sum + Math.pow(r - avgReturn, 2), 0) / returns.length)
    const sharpeRatio = volatility > 0 ? avgReturn / volatility : 0

    const maxDrawdown = Math.max(...equity.map(e => e.drawdown))

    return {
      totalReturn,
      totalReturnPercent,
      annualizedReturn: totalReturnPercent, // Simplified
      sharpeRatio,
      sortinoRatio: sharpeRatio, // Simplified
      calmarRatio: totalReturnPercent / (maxDrawdown || 1),
      maxDrawdown: maxDrawdown * initialCapital / 100,
      maxDrawdownPercent: maxDrawdown,
      volatility: volatility * 100,
      winRate,
      profitFactor,
      avgWin,
      avgLoss,
      avgTrade,
      totalTrades: trades.length,
      winningTrades: winningTrades.length,
      losingTrades: losingTrades.length,
      largestWin: Math.max(...trades.map(t => t.pnl), 0),
      largestLoss: Math.min(...trades.map(t => t.pnl), 0),
      avgHoldingPeriod: trades.length > 0 
        ? trades.reduce((sum, t) => sum + (t.holdingPeriod || 0), 0) / trades.length 
        : 0,
      totalCommission: trades.reduce((sum, t) => sum + t.commission, 0),
      totalSlippage: trades.reduce((sum, t) => sum + t.slippage, 0)
    }
  }

  /**
   * Calculate drawdown points
   */
  private calculateDrawdown(equity: EquityPoint[]): DrawdownPoint[] {
    const drawdown: DrawdownPoint[] = []
    let peak = equity[0]?.equity || 0

    for (const point of equity) {
      if (point.equity > peak) {
        peak = point.equity
      }

      const currentDrawdown = peak - point.equity
      const currentDrawdownPercent = peak > 0 ? (currentDrawdown / peak) * 100 : 0

      drawdown.push({
        timestamp: point.timestamp,
        drawdown: currentDrawdown,
        drawdownPercent: currentDrawdownPercent,
        peak,
        valley: point.equity
      })
    }

    return drawdown
  }

  /**
   * Calculate monthly returns
   */
  private calculateMonthlyReturns(trades: BacktestTrade[], equity: EquityPoint[]): MonthlyReturn[] {
    const monthlyReturns: MonthlyReturn[] = []
    const monthlyData = new Map<string, { trades: BacktestTrade[], startEquity: number, endEquity: number }>()

    // Group trades by month
    for (const trade of trades) {
      if (!trade.exitTime) continue
      
      const date = new Date(trade.exitTime)
      const key = `${date.getFullYear()}-${date.getMonth()}`
      
      if (!monthlyData.has(key)) {
        monthlyData.set(key, { trades: [], startEquity: 0, endEquity: 0 })
      }
      
      monthlyData.get(key)!.trades.push(trade)
    }

    // Calculate returns for each month
    for (const [key, data] of Array.from(monthlyData.entries())) {
      const [year, month] = key.split('-').map(Number)
      const monthReturn = data.trades.reduce((sum: number, trade: BacktestTrade) => sum + trade.pnl, 0)
      const monthReturnPercent = data.startEquity > 0 ? (monthReturn / data.startEquity) * 100 : 0

      monthlyReturns.push({
        year,
        month,
        return: monthReturn,
        returnPercent: monthReturnPercent,
        trades: data.trades.length
      })
    }

    return monthlyReturns.sort((a, b) => a.year - b.year || a.month - b.month)
  }

  /**
   * Calculate statistics
   */
  private calculateStatistics(
    trades: BacktestTrade[],
    signals: TradingSignal[],
    candles: HistoricalCandle[],
    monthlyReturns: MonthlyReturn[]
  ): BacktestStatistics {
    const startDate = new Date(candles[0].timestamp)
    const endDate = new Date(candles[candles.length - 1].timestamp)
    const duration = (endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24)

    const bestMonth = monthlyReturns.reduce((best, month) => 
      month.returnPercent > best.returnPercent ? month : best, 
      monthlyReturns[0] || { year: 0, month: 0, return: 0, returnPercent: -Infinity, trades: 0 }
    )

    const worstMonth = monthlyReturns.reduce((worst, month) => 
      month.returnPercent < worst.returnPercent ? month : worst,
      monthlyReturns[0] || { year: 0, month: 0, return: 0, returnPercent: Infinity, trades: 0 }
    )

    return {
      startDate,
      endDate,
      duration,
      totalCandles: candles.length,
      signalsGenerated: signals.length,
      signalsTraded: trades.length,
      avgTradesPerMonth: monthlyReturns.length > 0 ? trades.length / monthlyReturns.length : 0,
      bestMonth,
      worstMonth,
      consecutiveWins: 0, // Would need to calculate
      consecutiveLosses: 0, // Would need to calculate
      maxConsecutiveWins: 0, // Would need to calculate
      maxConsecutiveLosses: 0 // Would need to calculate
    }
  }

  /**
   * Calculate risk metrics
   */
  private calculateRiskMetrics(equity: EquityPoint[], trades: BacktestTrade[]): BacktestRiskMetrics {
    // Simplified risk metrics calculation
    const returns = equity.map((point, i) => {
      if (i === 0) return 0
      return (point.equity - equity[i - 1].equity) / equity[i - 1].equity
    }).slice(1)

    returns.sort((a, b) => a - b)
    const var95 = returns[Math.floor(returns.length * 0.05)] || 0
    const var99 = returns[Math.floor(returns.length * 0.01)] || 0

    return {
      var95: var95 * 100,
      var99: var99 * 100,
      cvar95: var95 * 100, // Simplified
      beta: 1, // Would need market data
      alpha: 0, // Would need market data
      informationRatio: 0, // Would need benchmark
      treynorRatio: 0, // Would need beta
      ulcerIndex: 0, // Would need to calculate
      painIndex: 0, // Would need to calculate
      gainToPainRatio: 0 // Would need to calculate
    }
  }

  /**
   * Helper: Calculate SMA
   */
  private calculateSMA(prices: number[], period: number): number {
    if (prices.length < period) return prices[prices.length - 1] || 0
    const recentPrices = prices.slice(-period)
    return recentPrices.reduce((sum, price) => sum + price, 0) / period
  }

  /**
   * Helper: Calculate RSI
   */
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

  /**
   * Get completed backtests
   */
  getCompletedBacktests(): BacktestResults[] {
    return Array.from(this.completedBacktests.values())
  }

  /**
   * Get backtest by ID
   */
  getBacktest(id: string): BacktestResults | null {
    return this.completedBacktests.get(id) || null
  }

  /**
   * Delete backtest
   */
  deleteBacktest(id: string): void {
    this.completedBacktests.delete(id)
  }
}

// Export singleton instance
export const backtestingEngine = BacktestingEngine.getInstance()
