import { type Address } from 'viem'
import { type PriceData } from './price-feed-manager'

export interface PortfolioPosition {
  id: string
  symbol: string
  contractAddress?: Address
  amount: number
  averageCost: number
  currentPrice: number
  marketValue: number
  unrealizedPnL: number
  unrealizedPnLPercent: number
  realizedPnL: number
  totalCost: number
  allocation: number // percentage of total portfolio
  firstPurchaseDate: number
  lastUpdateDate: number
  transactions: PortfolioTransaction[]
}

export interface PortfolioTransaction {
  id: string
  type: 'buy' | 'sell' | 'transfer_in' | 'transfer_out' | 'stake' | 'unstake' | 'reward'
  symbol: string
  amount: number
  price: number
  value: number
  fee: number
  timestamp: number
  txHash?: string
  exchange?: string
  notes?: string
}

export interface PortfolioMetrics {
  totalValue: number
  totalCost: number
  totalPnL: number
  totalPnLPercent: number
  realizedPnL: number
  unrealizedPnL: number
  dayChange: number
  dayChangePercent: number
  weekChange: number
  weekChangePercent: number
  monthChange: number
  monthChangePercent: number
  allTimeHigh: number
  allTimeLow: number
  maxDrawdown: number
  sharpeRatio: number
  volatility: number
  beta: number
  alpha: number
}

export interface AssetAllocation {
  symbol: string
  value: number
  percentage: number
  color: string
}

export interface RiskMetrics {
  portfolioRisk: 'low' | 'medium' | 'high' | 'extreme'
  riskScore: number // 0-100
  diversificationScore: number // 0-100
  concentrationRisk: number // 0-100
  volatilityRisk: number // 0-100
  correlationRisk: number // 0-100
  liquidityRisk: number // 0-100
  recommendations: string[]
}

export interface PerformanceData {
  timestamp: number
  totalValue: number
  totalPnL: number
  totalPnLPercent: number
  dayChange: number
  positions: number
}

export interface PortfolioSummary {
  metrics: PortfolioMetrics
  allocation: AssetAllocation[]
  riskMetrics: RiskMetrics
  topGainers: PortfolioPosition[]
  topLosers: PortfolioPosition[]
  recentTransactions: PortfolioTransaction[]
  performanceHistory: PerformanceData[]
}

export class PortfolioAnalytics {
  private static instance: PortfolioAnalytics
  private portfolios = new Map<string, Map<string, PortfolioPosition>>() // userId -> positions
  private transactions = new Map<string, PortfolioTransaction[]>() // userId -> transactions
  private performanceHistory = new Map<string, PerformanceData[]>() // userId -> history
  private priceCache = new Map<string, PriceData>()

  private constructor() {}

  static getInstance(): PortfolioAnalytics {
    if (!PortfolioAnalytics.instance) {
      PortfolioAnalytics.instance = new PortfolioAnalytics()
    }
    return PortfolioAnalytics.instance
  }

  /**
   * Add or update portfolio position
   */
  updatePosition(userId: string, position: Partial<PortfolioPosition> & { symbol: string }): void {
    if (!this.portfolios.has(userId)) {
      this.portfolios.set(userId, new Map())
    }

    const userPortfolio = this.portfolios.get(userId)!
    const existingPosition = userPortfolio.get(position.symbol)

    const updatedPosition: PortfolioPosition = {
      id: existingPosition?.id || `pos_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
      symbol: position.symbol,
      contractAddress: position.contractAddress,
      amount: position.amount || existingPosition?.amount || 0,
      averageCost: position.averageCost || existingPosition?.averageCost || 0,
      currentPrice: position.currentPrice || existingPosition?.currentPrice || 0,
      marketValue: 0, // Will be calculated
      unrealizedPnL: 0, // Will be calculated
      unrealizedPnLPercent: 0, // Will be calculated
      realizedPnL: position.realizedPnL || existingPosition?.realizedPnL || 0,
      totalCost: 0, // Will be calculated
      allocation: 0, // Will be calculated
      firstPurchaseDate: position.firstPurchaseDate || existingPosition?.firstPurchaseDate || Date.now(),
      lastUpdateDate: Date.now(),
      transactions: existingPosition?.transactions || []
    }

    // Calculate derived values
    updatedPosition.marketValue = updatedPosition.amount * updatedPosition.currentPrice
    updatedPosition.totalCost = updatedPosition.amount * updatedPosition.averageCost
    updatedPosition.unrealizedPnL = updatedPosition.marketValue - updatedPosition.totalCost
    updatedPosition.unrealizedPnLPercent = updatedPosition.totalCost > 0 
      ? (updatedPosition.unrealizedPnL / updatedPosition.totalCost) * 100 
      : 0

    userPortfolio.set(position.symbol, updatedPosition)
    this.updatePortfolioAllocations(userId)
  }

  /**
   * Add transaction to portfolio
   */
  addTransaction(userId: string, transaction: Omit<PortfolioTransaction, 'id'>): void {
    if (!this.transactions.has(userId)) {
      this.transactions.set(userId, [])
    }

    const fullTransaction: PortfolioTransaction = {
      ...transaction,
      id: `tx_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    }

    this.transactions.get(userId)!.push(fullTransaction)

    // Update position based on transaction
    this.processTransaction(userId, fullTransaction)
  }

  /**
   * Update price data for portfolio calculations
   */
  updatePriceData(priceData: PriceData): void {
    this.priceCache.set(priceData.symbol, priceData)

    // Update all positions with new price
    for (const [userId, portfolio] of Array.from(this.portfolios)) {
      const position = portfolio.get(priceData.symbol)
      if (position) {
        this.updatePosition(userId, {
          ...position,
          currentPrice: priceData.price
        })
      }
    }

    // Update performance history
    this.updatePerformanceHistory()
  }

  /**
   * Get portfolio summary for user
   */
  getPortfolioSummary(userId: string): PortfolioSummary {
    const positions = Array.from(this.portfolios.get(userId)?.values() || [])
    const userTransactions = this.transactions.get(userId) || []
    const history = this.performanceHistory.get(userId) || []

    const metrics = this.calculateMetrics(positions, history)
    const allocation = this.calculateAllocation(positions)
    const riskMetrics = this.calculateRiskMetrics(positions)

    const topGainers = positions
      .filter(p => p.unrealizedPnLPercent > 0)
      .sort((a, b) => b.unrealizedPnLPercent - a.unrealizedPnLPercent)
      .slice(0, 5)

    const topLosers = positions
      .filter(p => p.unrealizedPnLPercent < 0)
      .sort((a, b) => a.unrealizedPnLPercent - b.unrealizedPnLPercent)
      .slice(0, 5)

    const recentTransactions = userTransactions
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, 10)

    return {
      metrics,
      allocation,
      riskMetrics,
      topGainers,
      topLosers,
      recentTransactions,
      performanceHistory: history.slice(-100) // Last 100 data points
    }
  }

  /**
   * Get portfolio positions for user
   */
  getPositions(userId: string): PortfolioPosition[] {
    return Array.from(this.portfolios.get(userId)?.values() || [])
  }

  /**
   * Get portfolio transactions for user
   */
  getTransactions(userId: string, limit = 50): PortfolioTransaction[] {
    return (this.transactions.get(userId) || [])
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
  }

  /**
   * Calculate portfolio metrics
   */
  private calculateMetrics(positions: PortfolioPosition[], history: PerformanceData[]): PortfolioMetrics {
    const totalValue = positions.reduce((sum, p) => sum + p.marketValue, 0)
    const totalCost = positions.reduce((sum, p) => sum + p.totalCost, 0)
    const totalPnL = positions.reduce((sum, p) => sum + p.unrealizedPnL, 0)
    const realizedPnL = positions.reduce((sum, p) => sum + p.realizedPnL, 0)

    const totalPnLPercent = totalCost > 0 ? (totalPnL / totalCost) * 100 : 0

    // Calculate time-based changes
    const now = Date.now()
    const dayAgo = now - 24 * 60 * 60 * 1000
    const weekAgo = now - 7 * 24 * 60 * 60 * 1000
    const monthAgo = now - 30 * 24 * 60 * 60 * 1000

    const dayHistory = history.find(h => h.timestamp >= dayAgo)
    const weekHistory = history.find(h => h.timestamp >= weekAgo)
    const monthHistory = history.find(h => h.timestamp >= monthAgo)

    const dayChange = dayHistory ? totalValue - dayHistory.totalValue : 0
    const weekChange = weekHistory ? totalValue - weekHistory.totalValue : 0
    const monthChange = monthHistory ? totalValue - monthHistory.totalValue : 0

    const dayChangePercent = dayHistory && dayHistory.totalValue > 0 
      ? (dayChange / dayHistory.totalValue) * 100 : 0
    const weekChangePercent = weekHistory && weekHistory.totalValue > 0 
      ? (weekChange / weekHistory.totalValue) * 100 : 0
    const monthChangePercent = monthHistory && monthHistory.totalValue > 0 
      ? (monthChange / monthHistory.totalValue) * 100 : 0

    // Calculate all-time high/low
    const allTimeHigh = Math.max(totalValue, ...history.map(h => h.totalValue))
    const allTimeLow = history.length > 0 
      ? Math.min(totalValue, ...history.map(h => h.totalValue))
      : totalValue

    // Calculate max drawdown
    const maxDrawdown = this.calculateMaxDrawdown(history)

    // Mock advanced metrics (in real implementation, these would be calculated properly)
    const sharpeRatio = Math.random() * 2 - 0.5 // -0.5 to 1.5
    const volatility = Math.random() * 0.5 + 0.1 // 0.1 to 0.6
    const beta = Math.random() * 2 + 0.5 // 0.5 to 2.5
    const alpha = Math.random() * 0.2 - 0.1 // -0.1 to 0.1

    return {
      totalValue,
      totalCost,
      totalPnL,
      totalPnLPercent,
      realizedPnL,
      unrealizedPnL: totalPnL,
      dayChange,
      dayChangePercent,
      weekChange,
      weekChangePercent,
      monthChange,
      monthChangePercent,
      allTimeHigh,
      allTimeLow,
      maxDrawdown,
      sharpeRatio,
      volatility,
      beta,
      alpha
    }
  }

  /**
   * Calculate asset allocation
   */
  private calculateAllocation(positions: PortfolioPosition[]): AssetAllocation[] {
    const totalValue = positions.reduce((sum, p) => sum + p.marketValue, 0)
    
    const colors = [
      '#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6',
      '#06b6d4', '#84cc16', '#f97316', '#ec4899', '#6366f1'
    ]

    return positions
      .map((position, index) => ({
        symbol: position.symbol,
        value: position.marketValue,
        percentage: totalValue > 0 ? (position.marketValue / totalValue) * 100 : 0,
        color: colors[index % colors.length]
      }))
      .sort((a, b) => b.percentage - a.percentage)
  }

  /**
   * Calculate risk metrics
   */
  private calculateRiskMetrics(positions: PortfolioPosition[]): RiskMetrics {
    const totalValue = positions.reduce((sum, p) => sum + p.marketValue, 0)
    
    // Diversification score (based on number of positions and allocation spread)
    const diversificationScore = Math.min(100, positions.length * 10 + 
      (100 - Math.max(...positions.map(p => p.allocation))))

    // Concentration risk (highest single position allocation)
    const concentrationRisk = Math.max(...positions.map(p => p.allocation))

    // Mock other risk metrics
    const volatilityRisk = Math.random() * 100
    const correlationRisk = Math.random() * 100
    const liquidityRisk = Math.random() * 100

    const riskScore = (concentrationRisk + volatilityRisk + correlationRisk + liquidityRisk) / 4

    let portfolioRisk: RiskMetrics['portfolioRisk'] = 'low'
    if (riskScore > 75) portfolioRisk = 'extreme'
    else if (riskScore > 50) portfolioRisk = 'high'
    else if (riskScore > 25) portfolioRisk = 'medium'

    const recommendations: string[] = []
    if (concentrationRisk > 50) {
      recommendations.push('Consider diversifying your portfolio to reduce concentration risk')
    }
    if (positions.length < 5) {
      recommendations.push('Add more assets to improve diversification')
    }
    if (diversificationScore < 50) {
      recommendations.push('Rebalance portfolio to improve asset allocation')
    }

    return {
      portfolioRisk,
      riskScore,
      diversificationScore,
      concentrationRisk,
      volatilityRisk,
      correlationRisk,
      liquidityRisk,
      recommendations
    }
  }

  /**
   * Process transaction and update position
   */
  private processTransaction(userId: string, transaction: PortfolioTransaction): void {
    const portfolio = this.portfolios.get(userId)
    if (!portfolio) return

    const position = portfolio.get(transaction.symbol)
    
    if (transaction.type === 'buy') {
      if (position) {
        // Update average cost
        const totalCost = position.totalCost + transaction.value
        const totalAmount = position.amount + transaction.amount
        const newAverageCost = totalAmount > 0 ? totalCost / totalAmount : 0

        this.updatePosition(userId, {
          ...position,
          amount: totalAmount,
          averageCost: newAverageCost,
          transactions: [...position.transactions, transaction]
        })
      } else {
        // Create new position
        this.updatePosition(userId, {
          symbol: transaction.symbol,
          amount: transaction.amount,
          averageCost: transaction.price,
          currentPrice: transaction.price,
          firstPurchaseDate: transaction.timestamp,
          transactions: [transaction]
        })
      }
    } else if (transaction.type === 'sell' && position) {
      const newAmount = Math.max(0, position.amount - transaction.amount)
      const realizedPnL = (transaction.price - position.averageCost) * transaction.amount
      
      this.updatePosition(userId, {
        ...position,
        amount: newAmount,
        realizedPnL: position.realizedPnL + realizedPnL,
        transactions: [...position.transactions, transaction]
      })
    }
  }

  /**
   * Update portfolio allocations
   */
  private updatePortfolioAllocations(userId: string): void {
    const portfolio = this.portfolios.get(userId)
    if (!portfolio) return

    const positions = Array.from(portfolio.values())
    const totalValue = positions.reduce((sum, p) => sum + p.marketValue, 0)

    for (const position of positions) {
      position.allocation = totalValue > 0 ? (position.marketValue / totalValue) * 100 : 0
    }
  }

  /**
   * Update performance history
   */
  private updatePerformanceHistory(): void {
    const now = Date.now()
    
    for (const [userId, portfolio] of Array.from(this.portfolios)) {
      const positions = Array.from(portfolio.values())
      const totalValue = positions.reduce((sum: number, p: any) => sum + (p.marketValue || 0), 0)
      const totalPnL = positions.reduce((sum: number, p: any) => sum + (p.unrealizedPnL || 0), 0)
      const totalCost = positions.reduce((sum: number, p: any) => sum + (p.totalCost || 0), 0)
      const totalPnLPercent = totalCost > 0 ? (totalPnL / totalCost) * 100 : 0

      if (!this.performanceHistory.has(userId)) {
        this.performanceHistory.set(userId, [])
      }

      const history = this.performanceHistory.get(userId)!
      const lastEntry = history[history.length - 1]
      const dayChange = lastEntry ? totalValue - (lastEntry.totalValue || 0) : 0

      const performanceData: PerformanceData = {
        timestamp: now,
        totalValue: totalValue as number,
        totalPnL: totalPnL as number,
        totalPnLPercent,
        dayChange,
        positions: positions.length
      }

      history.push(performanceData)

      // Keep only last 1000 entries
      if (history.length > 1000) {
        history.splice(0, history.length - 1000)
      }
    }
  }

  /**
   * Calculate maximum drawdown
   */
  private calculateMaxDrawdown(history: PerformanceData[]): number {
    if (history.length < 2) return 0

    let maxDrawdown = 0
    let peak = history[0].totalValue

    for (const data of history) {
      if (data.totalValue > peak) {
        peak = data.totalValue
      }
      
      const drawdown = peak > 0 ? ((peak - data.totalValue) / peak) * 100 : 0
      maxDrawdown = Math.max(maxDrawdown, drawdown)
    }

    return maxDrawdown
  }

  /**
   * Import portfolio from external source
   */
  importPortfolio(userId: string, positions: Partial<PortfolioPosition>[]): void {
    for (const position of positions) {
      if (position.symbol) {
        this.updatePosition(userId, { ...position, symbol: position.symbol || 'UNKNOWN' })
      }
    }
  }

  /**
   * Export portfolio data
   */
  exportPortfolio(userId: string): {
    positions: PortfolioPosition[]
    transactions: PortfolioTransaction[]
    summary: PortfolioSummary
  } {
    return {
      positions: this.getPositions(userId),
      transactions: this.getTransactions(userId),
      summary: this.getPortfolioSummary(userId)
    }
  }

  /**
   * Clear portfolio data for user
   */
  clearPortfolio(userId: string): void {
    this.portfolios.delete(userId)
    this.transactions.delete(userId)
    this.performanceHistory.delete(userId)
  }
}

// Export singleton instance
export const portfolioAnalytics = PortfolioAnalytics.getInstance()
