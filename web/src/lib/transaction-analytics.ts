import { format, subDays, startOfDay, endOfDay, isWithinInterval, parseISO } from 'date-fns'
import { TransactionData, TransactionStatus, TransactionType } from './transaction-monitor'
import { SUPPORTED_CHAINS } from './chains'

export interface AnalyticsTimeframe {
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
  timeframe: AnalyticsTimeframe
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

export class TransactionAnalytics {
  private static readonly TIMEFRAMES: AnalyticsTimeframe[] = [
    { label: '24 Hours', days: 1, format: 'HH:mm' },
    { label: '7 Days', days: 7, format: 'MMM dd' },
    { label: '30 Days', days: 30, format: 'MMM dd' },
    { label: '90 Days', days: 90, format: 'MMM dd' },
    { label: '1 Year', days: 365, format: 'MMM yyyy' }
  ]

  static getTimeframes(): AnalyticsTimeframe[] {
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

  static calculateMetrics(transactions: TransactionData[]): TransactionMetrics {
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
      averageConfirmationTime
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
    const now = new Date()
    const start = subDays(now, timeframe.days)
    const data: TimeSeriesData[] = []

    // Generate time buckets based on timeframe
    const bucketSize = timeframe.days <= 1 ? 'hour' : 'day'
    const buckets = this.generateTimeBuckets(start, now, bucketSize)

    buckets.forEach(bucket => {
      const bucketTransactions = transactions.filter(tx => {
        const txDate = new Date(tx.timestamp)
        return isWithinInterval(txDate, bucket.interval)
      })

      const metrics = this.calculateMetrics(bucketTransactions)

      data.push({
        date: format(bucket.date, timeframe.format),
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
}
