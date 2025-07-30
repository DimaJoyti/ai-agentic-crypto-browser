import { format, subDays, startOfDay, endOfDay, isWithinInterval } from 'date-fns'
import { TransactionData, TransactionStatus, TransactionType } from './transaction-monitor'

// Re-export types for external use
export { TransactionType }
import { SUPPORTED_CHAINS } from './chains'

export enum AnalyticsTimeframe {
  LAST_24_HOURS = '1d',
  LAST_7_DAYS = '7d',
  LAST_30_DAYS = '30d',
  LAST_90_DAYS = '90d',
  LAST_YEAR = '365d',
  ALL_TIME = 'all'
}

export interface AnalyticsTimeframeConfig {
  label: string
  days: number
  format: string
}

export interface TransactionMetrics {
  totalTransactions: number
  successfulTransactions: number
  failedTransactions: number
  pendingTransactions: number
  successRate: number
  totalVolume: number
  averageGasUsed: number
  totalGasFees: number
  averageConfirmationTime: number
}

export interface AnalyticsMetrics extends TransactionMetrics {
  // Additional analytics-specific metrics
  totalValue: number
  averageGasPrice: number
  totalGasUsed: number
}

export interface ChainMetrics {
  chainId: number
  chainName: string
  transactionCount: number
  volume: number
  gasUsed: number
  successRate: number
  averageConfirmationTime: number
}

export interface TypeMetrics {
  type: TransactionType
  count: number
  volume: number
  successRate: number
  averageGasUsed: number
}

export interface TimeSeriesData {
  date: string
  timestamp: number
  transactions: number
  volume: number
  gasUsed: number
  successRate: number
  averageConfirmationTime: number
}

export interface AnalyticsFilters {
  timeframe: AnalyticsTimeframeConfig
  chains: number[]
  types: TransactionType[]
  status: TransactionStatus[]
  minAmount?: number
  maxAmount?: number
  dateRange?: {
    start: Date
    end: Date
  }
}

export interface AnalyticsInsights {
  recommendations: AnalyticsRecommendation[]
  warnings: AnalyticsWarning[]
  opportunities: AnalyticsOpportunity[]
  trends: AnalyticsTrend[]
}

export interface AnalyticsRecommendation {
  type: 'gas_optimization' | 'timing' | 'chain_selection' | 'cost_reduction' | 'security'
  title: string
  description: string
  impact: 'low' | 'medium' | 'high'
  potentialSavings?: string
  action?: string
  priority: number
}

export interface AnalyticsWarning {
  type: 'high_fees' | 'failed_transactions' | 'security' | 'unusual_activity' | 'performance'
  title: string
  description: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  affectedTransactions?: string[]
  recommendation?: string
  timestamp: number
}

export interface AnalyticsOpportunity {
  type: 'gas_savings' | 'timing_optimization' | 'chain_migration' | 'batch_transactions' | 'defi_yield'
  title: string
  description: string
  potentialBenefit: string
  difficulty: 'easy' | 'medium' | 'hard'
  timeframe: string
  estimatedSavings: number
}

export interface AnalyticsTrend {
  metric: string
  direction: 'increasing' | 'decreasing' | 'stable'
  change: number
  period: string
  significance: 'low' | 'medium' | 'high'
  description: string
}

export interface CostAnalysis {
  totalCosts: number
  costsByChain: Record<number, number>
  costsByType: Record<TransactionType, number>
  monthlyTrends: Array<{ month: string; costs: number; transactions: number }>
  costDistribution: Array<{ range: string; count: number; percentage: number }>
  expensivePeriods: Array<{ date: string; costs: number; reason: string }>
  averageCostPerTransaction: number
  costEfficiencyScore: number
}

export interface PerformanceAnalysis {
  averageConfirmationTime: number
  confirmationTimeByChain: Record<number, number>
  confirmationTimeByType: Record<TransactionType, number>
  gasEfficiency: number
  failureRate: number
  retryRate: number
  performanceScore: number
  bottlenecks: string[]
}

export class TransactionAnalytics {
  private static readonly TIMEFRAMES: AnalyticsTimeframeConfig[] = [
    { label: '24 Hours', days: 1, format: 'HH:mm' },
    { label: '7 Days', days: 7, format: 'MMM dd' },
    { label: '30 Days', days: 30, format: 'MMM dd' },
    { label: '90 Days', days: 90, format: 'MMM dd' },
    { label: '1 Year', days: 365, format: 'MMM yyyy' }
  ]

  private static readonly TIMEFRAME_MAPPING: Record<AnalyticsTimeframe, AnalyticsTimeframeConfig> = {
    [AnalyticsTimeframe.LAST_24_HOURS]: { label: '24 Hours', days: 1, format: 'HH:mm' },
    [AnalyticsTimeframe.LAST_7_DAYS]: { label: '7 Days', days: 7, format: 'MMM dd' },
    [AnalyticsTimeframe.LAST_30_DAYS]: { label: '30 Days', days: 30, format: 'MMM dd' },
    [AnalyticsTimeframe.LAST_90_DAYS]: { label: '90 Days', days: 90, format: 'MMM dd' },
    [AnalyticsTimeframe.LAST_YEAR]: { label: '1 Year', days: 365, format: 'MMM yyyy' },
    [AnalyticsTimeframe.ALL_TIME]: { label: 'All Time', days: 9999, format: 'MMM yyyy' }
  }

  static getTimeframeConfig(timeframe: AnalyticsTimeframe): AnalyticsTimeframeConfig {
    return this.TIMEFRAME_MAPPING[timeframe]
  }

  static getTimeframes(): AnalyticsTimeframeConfig[] {
    return this.TIMEFRAMES
  }

  static getDefaultFilters(): AnalyticsFilters {
    return {
      timeframe: this.TIMEFRAMES[1], // 7 days
      chains: [],
      types: [],
      status: [],
      dateRange: {
        start: subDays(new Date(), 7),
        end: new Date()
      }
    }
  }

  static filterTransactions(transactions: TransactionData[], filters: AnalyticsFilters): TransactionData[] {
    return transactions.filter(tx => {
      // Time filter
      const txDate = new Date(tx.timestamp)
      const { start, end } = filters.dateRange || {
        start: subDays(new Date(), filters.timeframe.days),
        end: new Date()
      }

      if (!isWithinInterval(txDate, { start: startOfDay(start), end: endOfDay(end) })) {
        return false
      }

      // Chain filter
      if (filters.chains.length > 0 && !filters.chains.includes(tx.chainId)) {
        return false
      }

      // Type filter
      if (filters.types.length > 0 && !filters.types.includes(tx.type)) {
        return false
      }

      // Status filter
      if (filters.status.length > 0 && !filters.status.includes(tx.status)) {
        return false
      }

      // Amount filter
      const amount = parseFloat(tx.value || '0')
      if (filters.minAmount !== undefined && amount < filters.minAmount) {
        return false
      }
      if (filters.maxAmount !== undefined && amount > filters.maxAmount) {
        return false
      }

      return true
    })
  }

  static calculateMetrics(transactions: TransactionData[]): AnalyticsMetrics {
    const total = transactions.length
    const successful = transactions.filter(tx => tx.status === TransactionStatus.CONFIRMED).length
    const failed = transactions.filter(tx => tx.status === TransactionStatus.FAILED).length
    const pending = transactions.filter(tx => tx.status === TransactionStatus.PENDING).length

    const totalVolume = transactions.reduce((sum, tx) => {
      return sum + parseFloat(tx.value || '0')
    }, 0)

    const gasUsedTransactions = transactions.filter(tx => tx.gasUsed)
    const totalGasUsed = gasUsedTransactions.reduce((sum, tx) => {
      return sum + parseInt(tx.gasUsed || '0')
    }, 0)
    const averageGasUsed = gasUsedTransactions.length > 0 ? totalGasUsed / gasUsedTransactions.length : 0

    const gasPriceTransactions = transactions.filter(tx => tx.gasPrice && tx.gasUsed)
    const totalGasFees = gasPriceTransactions.reduce((sum, tx) => {
      const gasPrice = parseFloat(tx.gasPrice || '0')
      const gasUsed = parseInt(tx.gasUsed || '0')
      return sum + (gasPrice * gasUsed)
    }, 0)

    // Calculate average confirmation time for confirmed transactions
    const confirmedTransactions = transactions.filter(tx => 
      tx.status === TransactionStatus.CONFIRMED && tx.blockNumber
    )
    const averageConfirmationTime = confirmedTransactions.length > 0
      ? confirmedTransactions.reduce((sum, tx) => sum + (tx.confirmations * 15), 0) / confirmedTransactions.length // Assume 15s block time
      : 0

    return {
      totalTransactions: total,
      successfulTransactions: successful,
      failedTransactions: failed,
      pendingTransactions: pending,
      successRate: total > 0 ? (successful / total) * 100 : 0,
      totalVolume,
      averageGasUsed,
      totalGasFees,
      averageConfirmationTime,
      totalValue: totalVolume, // Use totalVolume as totalValue
      averageGasPrice: total > 0 ? totalGasFees / total : 0,
      totalGasUsed
    }
  }

  static calculateChainMetrics(transactions: TransactionData[]): ChainMetrics[] {
    const chainGroups = transactions.reduce((groups, tx) => {
      if (!groups[tx.chainId]) {
        groups[tx.chainId] = []
      }
      groups[tx.chainId].push(tx)
      return groups
    }, {} as Record<number, TransactionData[]>)

    return Object.entries(chainGroups).map(([chainId, txs]) => {
      const chain = SUPPORTED_CHAINS[parseInt(chainId)]
      const metrics = this.calculateMetrics(txs)
      
      return {
        chainId: parseInt(chainId),
        chainName: chain?.name || `Chain ${chainId}`,
        transactionCount: txs.length,
        volume: metrics.totalVolume,
        gasUsed: metrics.averageGasUsed,
        successRate: metrics.successRate,
        averageConfirmationTime: metrics.averageConfirmationTime
      }
    }).sort((a, b) => b.transactionCount - a.transactionCount)
  }

  static calculateTypeMetrics(transactions: TransactionData[]): TypeMetrics[] {
    const typeGroups = transactions.reduce((groups, tx) => {
      if (!groups[tx.type]) {
        groups[tx.type] = []
      }
      groups[tx.type].push(tx)
      return groups
    }, {} as Record<TransactionType, TransactionData[]>)

    return Object.entries(typeGroups).map(([type, txs]) => {
      const metrics = this.calculateMetrics(txs)
      
      return {
        type: type as TransactionType,
        count: txs.length,
        volume: metrics.totalVolume,
        successRate: metrics.successRate,
        averageGasUsed: metrics.averageGasUsed
      }
    }).sort((a, b) => b.count - a.count)
  }

  static generateTimeSeriesData(
    transactions: TransactionData[],
    timeframe: AnalyticsTimeframe
  ): TimeSeriesData[] {
    const config = this.getTimeframeConfig(timeframe)
    const now = new Date()
    const start = subDays(now, config.days)
    const data: TimeSeriesData[] = []

    // Generate time buckets based on timeframe
    const bucketSize = config.days <= 1 ? 'hour' : 'day'
    const buckets = this.generateTimeBuckets(start, now, bucketSize)

    buckets.forEach(bucket => {
      const bucketTransactions = transactions.filter(tx => {
        const txDate = new Date(tx.timestamp)
        return isWithinInterval(txDate, bucket.interval)
      })

      const metrics = this.calculateMetrics(bucketTransactions)

      data.push({
        date: format(bucket.date, config.format),
        timestamp: bucket.date.getTime(),
        transactions: bucketTransactions.length,
        volume: metrics.totalVolume,
        gasUsed: metrics.averageGasUsed,
        successRate: metrics.successRate,
        averageConfirmationTime: metrics.averageConfirmationTime
      })
    })

    return data
  }

  private static generateTimeBuckets(start: Date, end: Date, bucketSize: 'hour' | 'day') {
    const buckets = []
    const current = new Date(start)

    while (current <= end) {
      const bucketStart = new Date(current)
      const bucketEnd = new Date(current)

      if (bucketSize === 'hour') {
        bucketEnd.setHours(bucketEnd.getHours() + 1)
      } else {
        bucketEnd.setDate(bucketEnd.getDate() + 1)
      }

      buckets.push({
        date: new Date(bucketStart),
        interval: { start: bucketStart, end: bucketEnd }
      })

      if (bucketSize === 'hour') {
        current.setHours(current.getHours() + 1)
      } else {
        current.setDate(current.getDate() + 1)
      }
    }

    return buckets
  }

  static exportToCSV(transactions: TransactionData[], filename?: string): void {
    const headers = [
      'Hash',
      'Chain',
      'Type',
      'Status',
      'Amount',
      'Gas Used',
      'Gas Price',
      'Confirmations',
      'Block Number',
      'Timestamp',
      'Date'
    ]

    const rows = transactions.map(tx => {
      const chain = SUPPORTED_CHAINS[tx.chainId]
      return [
        tx.hash,
        chain?.shortName || tx.chainId.toString(),
        tx.type,
        tx.status,
        tx.value || '0',
        tx.gasUsed || '',
        tx.gasPrice || '',
        tx.confirmations.toString(),
        tx.blockNumber?.toString() || '',
        tx.timestamp.toString(),
        new Date(tx.timestamp).toISOString()
      ]
    })

    const csvContent = [headers, ...rows]
      .map(row => row.map(field => `"${field}"`).join(','))
      .join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    
    link.setAttribute('href', url)
    link.setAttribute('download', filename || `transactions-${Date.now()}.csv`)
    link.style.visibility = 'hidden'
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  static exportToJSON(transactions: TransactionData[], filename?: string): void {
    const jsonContent = JSON.stringify(transactions, null, 2)
    const blob = new Blob([jsonContent], { type: 'application/json;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    
    link.setAttribute('href', url)
    link.setAttribute('download', filename || `transactions-${Date.now()}.json`)
    link.style.visibility = 'hidden'
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  static formatCurrency(amount: number, decimals = 4): string {
    if (amount === 0) return '0'
    if (amount < 0.0001) return '< 0.0001'
    return amount.toFixed(decimals)
  }

  static formatGas(gas: number): string {
    if (gas === 0) return '0'
    if (gas < 1000) return gas.toString()
    if (gas < 1000000) return `${(gas / 1000).toFixed(1)}K`
    return `${(gas / 1000000).toFixed(1)}M`
  }

  static formatTime(seconds: number): string {
    if (seconds < 60) return `${Math.round(seconds)}s`
    if (seconds < 3600) return `${Math.round(seconds / 60)}m`
    return `${Math.round(seconds / 3600)}h`
  }

  static generateInsights(transactions: TransactionData[]): AnalyticsInsights {
    const insights: AnalyticsInsights = {
      recommendations: [],
      warnings: [],
      opportunities: [],
      trends: []
    }

    const metrics = this.calculateMetrics(transactions)
    const costAnalysis = this.generateCostAnalysis(transactions)
    // const performanceAnalysis = this.generatePerformanceAnalysis(transactions)

    // Generate recommendations
    if (metrics.successRate < 0.95) {
      insights.recommendations.push({
        type: 'gas_optimization',
        title: 'Improve Transaction Success Rate',
        description: 'Your transaction success rate is below 95%. Consider optimizing gas settings.',
        impact: 'high',
        potentialSavings: '5-15% cost reduction',
        action: 'Review and adjust gas price settings',
        priority: 1
      })
    }

    if (metrics.averageGasUsed > 100000) {
      insights.recommendations.push({
        type: 'gas_optimization',
        title: 'Optimize Gas Consumption',
        description: 'Your average gas usage is high. Review contract interactions to reduce consumption.',
        impact: 'medium',
        potentialSavings: '10-30% gas savings',
        action: 'Optimize contract calls and batch transactions',
        priority: 2
      })
    }

    if (costAnalysis.averageCostPerTransaction > 50) {
      insights.recommendations.push({
        type: 'chain_selection',
        title: 'Consider Layer 2 Solutions',
        description: 'High transaction costs detected. Layer 2 solutions could reduce fees significantly.',
        impact: 'high',
        potentialSavings: '80-95% fee reduction',
        action: 'Migrate to Arbitrum, Optimism, or Polygon',
        priority: 1
      })
    }

    // Generate warnings
    if (metrics.failedTransactions > metrics.totalTransactions * 0.1) {
      insights.warnings.push({
        type: 'failed_transactions',
        title: 'High Failure Rate Detected',
        description: 'More than 10% of your transactions are failing.',
        severity: 'high',
        recommendation: 'Check gas settings and network conditions before submitting transactions',
        timestamp: Date.now()
      })
    }

    if (costAnalysis.costEfficiencyScore < 60) {
      insights.warnings.push({
        type: 'high_fees',
        title: 'Poor Cost Efficiency',
        description: 'Your transaction costs are higher than optimal.',
        severity: 'medium',
        recommendation: 'Review gas optimization strategies and consider alternative chains',
        timestamp: Date.now()
      })
    }

    // Generate opportunities
    const layer1Transactions = transactions.filter(tx => tx.chainId === 1)
    if (layer1Transactions.length > transactions.length * 0.8) {
      insights.opportunities.push({
        type: 'chain_migration',
        title: 'Layer 2 Migration Opportunity',
        description: 'Most transactions are on Ethereum mainnet. Significant savings available on Layer 2.',
        potentialBenefit: '80-95% cost reduction',
        difficulty: 'medium',
        timeframe: '1-2 weeks',
        estimatedSavings: costAnalysis.totalCosts * 0.85
      })
    }

    // Generate trends
    const recentTransactions = transactions.slice(-20)
    const olderTransactions = transactions.slice(-40, -20)

    if (recentTransactions.length > 0 && olderTransactions.length > 0) {
      const recentAvgGas = recentTransactions.reduce((sum, tx) => sum + (Number(tx.gasUsed) || 0), 0) / recentTransactions.length
      const olderAvgGas = olderTransactions.reduce((sum, tx) => sum + (Number(tx.gasUsed) || 0), 0) / olderTransactions.length
      const gasChange = ((recentAvgGas - olderAvgGas) / olderAvgGas) * 100

      insights.trends.push({
        metric: 'Gas Usage',
        direction: gasChange > 5 ? 'increasing' : gasChange < -5 ? 'decreasing' : 'stable',
        change: Math.abs(gasChange),
        period: 'Last 20 transactions',
        significance: Math.abs(gasChange) > 20 ? 'high' : Math.abs(gasChange) > 10 ? 'medium' : 'low',
        description: `Gas usage has ${gasChange > 0 ? 'increased' : gasChange < 0 ? 'decreased' : 'remained stable'} by ${Math.abs(gasChange).toFixed(1)}%`
      })
    }

    return insights
  }

  static generateCostAnalysis(transactions: TransactionData[]): CostAnalysis {
    const totalCosts = transactions.reduce((sum, tx) => {
      const gasUsed = Number(tx.gasUsed) || 0
      const gasPrice = Number(tx.gasPrice) || 0
      return sum + (gasUsed * gasPrice / 1e18) // Convert to ETH
    }, 0)

    const costsByChain = transactions.reduce((acc, tx) => {
      const gasUsed = Number(tx.gasUsed) || 0
      const gasPrice = Number(tx.gasPrice) || 0
      const cost = gasUsed * gasPrice / 1e18
      acc[tx.chainId] = (acc[tx.chainId] || 0) + cost
      return acc
    }, {} as Record<number, number>)

    const costsByType = transactions.reduce((acc, tx) => {
      const gasUsed = Number(tx.gasUsed) || 0
      const gasPrice = Number(tx.gasPrice) || 0
      const cost = gasUsed * gasPrice / 1e18
      acc[tx.type] = (acc[tx.type] || 0) + cost
      return acc
    }, {} as Record<TransactionType, number>)

    const averageCostPerTransaction = transactions.length > 0 ? totalCosts / transactions.length : 0
    const costEfficiencyScore = Math.max(0, 100 - (averageCostPerTransaction * 1000)) // Simplified scoring

    return {
      totalCosts,
      costsByChain,
      costsByType,
      monthlyTrends: [],
      costDistribution: [],
      expensivePeriods: [],
      averageCostPerTransaction,
      costEfficiencyScore
    }
  }

  static generatePerformanceAnalysis(transactions: TransactionData[]): PerformanceAnalysis {
    const confirmedTxs = transactions.filter(tx => tx.status === TransactionStatus.CONFIRMED)

    const averageConfirmationTime = confirmedTxs.length > 0
      ? confirmedTxs.reduce((sum, tx) => sum + (tx.confirmations || 0), 0) / confirmedTxs.length
      : 0

    const failureRate = transactions.length > 0
      ? (transactions.filter(tx => tx.status === TransactionStatus.FAILED).length / transactions.length) * 100
      : 0

    const performanceScore = Math.max(0, 100 - failureRate - (averageConfirmationTime / 60))

    return {
      averageConfirmationTime,
      confirmationTimeByChain: {} as Record<number, number>,
      confirmationTimeByType: {} as Record<TransactionType, number>,
      gasEfficiency: 85,
      failureRate,
      retryRate: 0,
      performanceScore,
      bottlenecks: []
    }
  }
}
