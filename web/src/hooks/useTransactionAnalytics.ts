import { useState, useEffect, useCallback, useMemo } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { 
  TransactionAnalytics,
  AnalyticsTimeframe,
  AnalyticsTimeframeConfig,
  type AnalyticsMetrics,
  type AnalyticsFilters,
  type AnalyticsInsights,
  type CostAnalysis,
  type PerformanceAnalysis,
  TransactionType
} from '@/lib/transaction-analytics'
import { type TransactionData, TransactionStatus } from '@/lib/transaction-monitor'
import { toast } from 'sonner'

export interface TransactionAnalyticsState {
  transactions: TransactionData[]
  filteredTransactions: TransactionData[]
  metrics: AnalyticsMetrics
  insights: AnalyticsInsights
  costAnalysis: CostAnalysis
  performanceAnalysis: PerformanceAnalysis
  isLoading: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseTransactionAnalyticsOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  defaultTimeframe?: AnalyticsTimeframe
  maxTransactions?: number
  refreshInterval?: number
}

export interface UseTransactionAnalyticsReturn {
  // State
  state: TransactionAnalyticsState
  
  // Data Management
  loadTransactions: (filters?: AnalyticsFilters) => Promise<void>
  addTransaction: (transaction: TransactionData) => void
  updateTransaction: (hash: string, updates: Partial<TransactionData>) => void
  
  // Analytics
  getMetrics: (filters?: AnalyticsFilters) => AnalyticsMetrics
  getInsights: (filters?: AnalyticsFilters) => AnalyticsInsights
  getCostAnalysis: (filters?: AnalyticsFilters) => CostAnalysis
  getPerformanceAnalysis: (filters?: AnalyticsFilters) => PerformanceAnalysis
  
  // Filtering
  applyFilters: (filters: AnalyticsFilters) => void
  clearFilters: () => void
  
  // Export
  exportToCSV: (filename?: string) => void
  exportToJSON: (filename?: string) => void
  
  // Utilities
  refresh: () => Promise<void>
  clearError: () => void
}

export const useTransactionAnalytics = (
  options: UseTransactionAnalyticsOptions = {}
): UseTransactionAnalyticsReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    defaultTimeframe = AnalyticsTimeframe.LAST_30_DAYS,
    maxTransactions = 1000,
    refreshInterval = 30000
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<TransactionAnalyticsState>({
    transactions: [],
    filteredTransactions: [],
    metrics: {
      totalTransactions: 0,
      successfulTransactions: 0,
      failedTransactions: 0,
      pendingTransactions: 0,
      successRate: 0,
      totalVolume: 0,
      averageGasUsed: 0,
      totalGasFees: 0,
      averageConfirmationTime: 0,
      totalValue: 0,
      averageGasPrice: 0,
      totalGasUsed: 0
    },
    insights: {
      recommendations: [],
      warnings: [],
      opportunities: [],
      trends: []
    },
    costAnalysis: {
      totalCosts: 0,
      costsByChain: {},
      costsByType: {
        [TransactionType.SEND]: 0,
        [TransactionType.RECEIVE]: 0,
        [TransactionType.SWAP]: 0,
        [TransactionType.APPROVE]: 0,
        [TransactionType.STAKE]: 0,
        [TransactionType.UNSTAKE]: 0,
        [TransactionType.MINT]: 0,
        [TransactionType.BURN]: 0,
        [TransactionType.CONTRACT_INTERACTION]: 0
      },
      monthlyTrends: [],
      costDistribution: [],
      expensivePeriods: [],
      averageCostPerTransaction: 0,
      costEfficiencyScore: 0
    },
    performanceAnalysis: {
      averageConfirmationTime: 0,
      confirmationTimeByChain: {},
      confirmationTimeByType: {
        [TransactionType.SEND]: 0,
        [TransactionType.RECEIVE]: 0,
        [TransactionType.SWAP]: 0,
        [TransactionType.APPROVE]: 0,
        [TransactionType.STAKE]: 0,
        [TransactionType.UNSTAKE]: 0,
        [TransactionType.MINT]: 0,
        [TransactionType.BURN]: 0,
        [TransactionType.CONTRACT_INTERACTION]: 0
      },
      gasEfficiency: 0,
      failureRate: 0,
      retryRate: 0,
      performanceScore: 0,
      bottlenecks: []
    },
    isLoading: false,
    error: null,
    lastUpdate: null
  })

  const [currentFilters, setCurrentFilters] = useState<AnalyticsFilters>({
    timeframe: TransactionAnalytics.getTimeframeConfig(defaultTimeframe),
    chains: [],
    types: [],
    status: [],
    dateRange: undefined
  })

  // Load transactions from various sources
  const loadTransactions = useCallback(async (filters?: AnalyticsFilters) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      // In a real implementation, this would fetch from:
      // - Local storage
      // - Transaction monitor
      // - Blockchain APIs
      // - Database
      
      // For now, we'll simulate loading transactions
      const mockTransactions: TransactionData[] = []
      
      // Apply filters if provided
      const activeFilters = filters || currentFilters
      let filteredTransactions = mockTransactions

      if (activeFilters.timeframe.days > 0) {
        const now = Date.now()
        const timeLimit = activeFilters.timeframe.days * 24 * 60 * 60 * 1000
        if (timeLimit > 0) {
          filteredTransactions = filteredTransactions.filter(tx => 
            tx.timestamp > now - timeLimit
          )
        }
      }

      if (activeFilters.chains.length > 0) {
        filteredTransactions = filteredTransactions.filter(tx => 
          activeFilters.chains.includes(tx.chainId)
        )
      }

      if (activeFilters.types.length > 0) {
        filteredTransactions = filteredTransactions.filter(tx => 
          activeFilters.types.includes(tx.type)
        )
      }

      if (activeFilters.status.length > 0) {
        filteredTransactions = filteredTransactions.filter(tx => 
          activeFilters.status.includes(tx.status)
        )
      }

      // Limit number of transactions
      if (filteredTransactions.length > maxTransactions) {
        filteredTransactions = filteredTransactions
          .sort((a, b) => b.timestamp - a.timestamp)
          .slice(0, maxTransactions)
      }

      // Calculate analytics
      const metrics = TransactionAnalytics.calculateMetrics(filteredTransactions)
      const insights = TransactionAnalytics.generateInsights(filteredTransactions)
      const costAnalysis = TransactionAnalytics.generateCostAnalysis(filteredTransactions)
      const performanceAnalysis = TransactionAnalytics.generatePerformanceAnalysis(filteredTransactions)

      setState(prev => ({
        ...prev,
        transactions: mockTransactions,
        filteredTransactions,
        metrics,
        insights,
        costAnalysis,
        performanceAnalysis,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      if (enableNotifications && filteredTransactions.length > 0) {
        toast.success(`Loaded ${filteredTransactions.length} transactions`)
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to load transactions', { description: errorMessage })
      }
    }
  }, [currentFilters, maxTransactions, enableNotifications])

  // Auto-load transactions
  useEffect(() => {
    if (autoLoad && address) {
      loadTransactions()
    }
  }, [autoLoad, address, loadTransactions])

  // Auto-refresh
  useEffect(() => {
    if (refreshInterval > 0 && address) {
      const interval = setInterval(() => {
        loadTransactions()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [refreshInterval, address, loadTransactions])

  // Add transaction
  const addTransaction = useCallback((transaction: TransactionData) => {
    setState(prev => {
      const newTransactions = [...prev.transactions, transaction]
      const filteredTransactions = newTransactions // Would apply current filters
      const metrics = TransactionAnalytics.calculateMetrics(filteredTransactions)
      const insights = TransactionAnalytics.generateInsights(filteredTransactions)
      const costAnalysis = TransactionAnalytics.generateCostAnalysis(filteredTransactions)
      const performanceAnalysis = TransactionAnalytics.generatePerformanceAnalysis(filteredTransactions)

      return {
        ...prev,
        transactions: newTransactions,
        filteredTransactions,
        metrics,
        insights,
        costAnalysis,
        performanceAnalysis,
        lastUpdate: Date.now()
      }
    })

    if (enableNotifications) {
      toast.info('Transaction added to analytics')
    }
  }, [enableNotifications])

  // Update transaction
  const updateTransaction = useCallback((hash: string, updates: Partial<TransactionData>) => {
    setState(prev => {
      const newTransactions = prev.transactions.map(tx => 
        tx.hash === hash ? { ...tx, ...updates } : tx
      )
      const filteredTransactions = newTransactions // Would apply current filters
      const metrics = TransactionAnalytics.calculateMetrics(filteredTransactions)
      const insights = TransactionAnalytics.generateInsights(filteredTransactions)
      const costAnalysis = TransactionAnalytics.generateCostAnalysis(filteredTransactions)
      const performanceAnalysis = TransactionAnalytics.generatePerformanceAnalysis(filteredTransactions)

      return {
        ...prev,
        transactions: newTransactions,
        filteredTransactions,
        metrics,
        insights,
        costAnalysis,
        performanceAnalysis,
        lastUpdate: Date.now()
      }
    })
  }, [])

  // Get metrics with optional filters
  const getMetrics = useCallback((filters?: AnalyticsFilters): AnalyticsMetrics => {
    const transactions = filters ? state.transactions : state.filteredTransactions
    return TransactionAnalytics.calculateMetrics(transactions)
  }, [state.transactions, state.filteredTransactions])

  // Get insights with optional filters
  const getInsights = useCallback((filters?: AnalyticsFilters): AnalyticsInsights => {
    const transactions = filters ? state.transactions : state.filteredTransactions
    return TransactionAnalytics.generateInsights(transactions)
  }, [state.transactions, state.filteredTransactions])

  // Get cost analysis with optional filters
  const getCostAnalysis = useCallback((filters?: AnalyticsFilters): CostAnalysis => {
    const transactions = filters ? state.transactions : state.filteredTransactions
    return TransactionAnalytics.generateCostAnalysis(transactions)
  }, [state.transactions, state.filteredTransactions])

  // Get performance analysis with optional filters
  const getPerformanceAnalysis = useCallback((filters?: AnalyticsFilters): PerformanceAnalysis => {
    const transactions = filters ? state.transactions : state.filteredTransactions
    return TransactionAnalytics.generatePerformanceAnalysis(transactions)
  }, [state.transactions, state.filteredTransactions])

  // Apply filters
  const applyFilters = useCallback((filters: AnalyticsFilters) => {
    setCurrentFilters(filters)
    loadTransactions(filters)
  }, [loadTransactions])

  // Clear filters
  const clearFilters = useCallback(() => {
    const defaultFilters: AnalyticsFilters = {
      timeframe: TransactionAnalytics.getTimeframeConfig(defaultTimeframe),
      chains: [],
      types: [],
      status: []
    }
    setCurrentFilters(defaultFilters)
    loadTransactions(defaultFilters)
  }, [defaultTimeframe, loadTransactions])

  // Export to CSV
  const exportToCSV = useCallback((filename?: string) => {
    try {
      TransactionAnalytics.exportToCSV(state.filteredTransactions, filename)
      
      if (enableNotifications) {
        toast.success('Transactions exported to CSV')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      
      if (enableNotifications) {
        toast.error('Failed to export CSV', { description: errorMessage })
      }
    }
  }, [state.filteredTransactions, enableNotifications])

  // Export to JSON
  const exportToJSON = useCallback((filename?: string) => {
    try {
      TransactionAnalytics.exportToJSON(state.filteredTransactions, filename)
      
      if (enableNotifications) {
        toast.success('Transactions exported to JSON')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      
      if (enableNotifications) {
        toast.error('Failed to export JSON', { description: errorMessage })
      }
    }
  }, [state.filteredTransactions, enableNotifications])

  // Refresh
  const refresh = useCallback(async () => {
    await loadTransactions()
  }, [loadTransactions])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    loadTransactions,
    addTransaction,
    updateTransaction,
    getMetrics,
    getInsights,
    getCostAnalysis,
    getPerformanceAnalysis,
    applyFilters,
    clearFilters,
    exportToCSV,
    exportToJSON,
    refresh,
    clearError
  }
}

// Simplified hook for basic analytics
export const useTransactionMetrics = (timeframe: AnalyticsTimeframe = AnalyticsTimeframe.LAST_30_DAYS) => {
  const { state, applyFilters } = useTransactionAnalytics({
    defaultTimeframe: timeframe,
    autoLoad: true
  })

  return {
    metrics: state.metrics,
    insights: state.insights,
    isLoading: state.isLoading,
    error: state.error,
    changeTimeframe: (newTimeframe: AnalyticsTimeframe) => {
      const timeframeConfig = TransactionAnalytics.getTimeframeConfig(newTimeframe)
      applyFilters({
        timeframe: timeframeConfig,
        chains: [],
        types: [],
        status: []
      })
    }
  }
}
