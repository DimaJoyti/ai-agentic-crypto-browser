import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  transactionHistoryService, 
  type TransactionDetails,
  type TransactionSearchFilters,
  type TransactionSearchResult,
  type TransactionSummary,
  TransactionType,
  TransactionStatus,
  TransactionCategory 
} from '@/lib/transaction-history'
import { toast } from 'sonner'

export interface UseTransactionHistoryOptions {
  address?: Address
  chainId?: number
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface TransactionHistoryState {
  transactions: TransactionDetails[]
  summary: TransactionSummary | null
  isLoading: boolean
  isSearching: boolean
  error: string | null
  hasMore: boolean
  nextCursor?: string
  totalCount: number
  lastUpdated: number
}

export function useTransactionHistory(options: UseTransactionHistoryOptions = {}) {
  const {
    address,
    chainId = 1,
    autoRefresh = true,
    refreshInterval = 30000, // 30 seconds
    enableNotifications = true
  } = options

  const [state, setState] = useState<TransactionHistoryState>({
    transactions: [],
    summary: null,
    isLoading: false,
    isSearching: false,
    error: null,
    hasMore: false,
    totalCount: 0,
    lastUpdated: 0
  })

  const [searchFilters, setSearchFilters] = useState<TransactionSearchFilters>({})

  // Load transaction history
  const loadTransactions = useCallback(async (
    loadMore: boolean = false,
    filters?: TransactionSearchFilters
  ) => {
    if (!address) return

    setState(prev => ({ 
      ...prev, 
      isLoading: !loadMore, 
      isSearching: !!filters,
      error: null 
    }))

    try {
      let result: TransactionSearchResult

      if (filters && Object.keys(filters).length > 0) {
        // Search with filters
        result = await transactionHistoryService.searchTransactions(
          { ...filters, address, chainId },
          50,
          loadMore ? state.nextCursor : undefined
        )
      } else {
        // Load regular history
        result = await transactionHistoryService.getTransactionHistory(
          address,
          chainId,
          50,
          loadMore ? state.nextCursor : undefined
        )
      }

      setState(prev => ({
        ...prev,
        transactions: loadMore ? [...prev.transactions, ...result.transactions] : result.transactions,
        hasMore: result.hasMore,
        nextCursor: result.nextCursor,
        totalCount: result.totalCount,
        isLoading: false,
        isSearching: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load transactions'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
        isSearching: false
      }))

      if (enableNotifications) {
        toast.error('Transaction Error', {
          description: errorMessage
        })
      }
    }
  }, [address, chainId, enableNotifications, state.nextCursor])

  // Load transaction summary
  const loadSummary = useCallback(async (fromDate?: Date, toDate?: Date) => {
    if (!address) return

    try {
      const summary = await transactionHistoryService.getTransactionSummary(
        address,
        chainId,
        fromDate,
        toDate
      )

      setState(prev => ({
        ...prev,
        summary
      }))
    } catch (error) {
      console.error('Failed to load transaction summary:', error)
    }
  }, [address, chainId])

  // Search transactions
  const searchTransactions = useCallback(async (filters: TransactionSearchFilters) => {
    setSearchFilters(filters)
    await loadTransactions(false, filters)
  }, [loadTransactions])

  // Clear search
  const clearSearch = useCallback(async () => {
    setSearchFilters({})
    await loadTransactions(false)
  }, [loadTransactions])

  // Load more transactions
  const loadMore = useCallback(async () => {
    if (state.hasMore && !state.isLoading) {
      await loadTransactions(true, Object.keys(searchFilters).length > 0 ? searchFilters : undefined)
    }
  }, [loadTransactions, state.hasMore, state.isLoading, searchFilters])

  // Get specific transaction
  const getTransaction = useCallback(async (hash: string) => {
    try {
      return await transactionHistoryService.getTransaction(hash as any)
    } catch (error) {
      console.error('Failed to get transaction:', error)
      return null
    }
  }, [])

  // Filter transactions by type
  const getTransactionsByType = useCallback((type: TransactionType) => {
    return state.transactions.filter(tx => tx.type === type)
  }, [state.transactions])

  // Filter transactions by category
  const getTransactionsByCategory = useCallback((category: TransactionCategory) => {
    return state.transactions.filter(tx => tx.category === category)
  }, [state.transactions])

  // Filter transactions by status
  const getTransactionsByStatus = useCallback((status: TransactionStatus) => {
    return state.transactions.filter(tx => tx.status === status)
  }, [state.transactions])

  // Get transactions by date range
  const getTransactionsByDateRange = useCallback((fromDate: Date, toDate: Date) => {
    return state.transactions.filter(tx => 
      tx.timestamp >= fromDate.getTime() && tx.timestamp <= toDate.getTime()
    )
  }, [state.transactions])

  // Get recent transactions
  const getRecentTransactions = useCallback((limit: number = 10) => {
    return state.transactions
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
  }, [state.transactions])

  // Get failed transactions
  const getFailedTransactions = useCallback(() => {
    return state.transactions.filter(tx => tx.status === TransactionStatus.FAILED)
  }, [state.transactions])

  // Get pending transactions
  const getPendingTransactions = useCallback(() => {
    return state.transactions.filter(tx => tx.status === TransactionStatus.PENDING)
  }, [state.transactions])

  // Calculate analytics
  const getAnalytics = useCallback(() => {
    const transactions = state.transactions
    const confirmedTxs = transactions.filter(tx => tx.status === TransactionStatus.CONFIRMED)
    
    const totalValue = transactions.reduce((sum, tx) => sum + parseFloat(tx.value), 0)
    const totalGasCost = transactions.reduce((sum, tx) => sum + parseFloat(tx.gasCostETH), 0)
    const averageGasPrice = transactions.length > 0 
      ? transactions.reduce((sum, tx) => sum + parseFloat(tx.gasPrice), 0) / transactions.length
      : 0

    const successRate = transactions.length > 0 
      ? (confirmedTxs.length / transactions.length) * 100 
      : 0

    const typeDistribution = transactions.reduce((acc, tx) => {
      acc[tx.type] = (acc[tx.type] || 0) + 1
      return acc
    }, {} as Record<TransactionType, number>)

    const categoryDistribution = transactions.reduce((acc, tx) => {
      acc[tx.category] = (acc[tx.category] || 0) + 1
      return acc
    }, {} as Record<TransactionCategory, number>)

    return {
      totalTransactions: transactions.length,
      confirmedTransactions: confirmedTxs.length,
      totalValue,
      totalGasCost,
      averageGasPrice,
      successRate,
      typeDistribution,
      categoryDistribution
    }
  }, [state.transactions])

  // Format utilities
  const formatValue = useCallback((value: string, decimals: number = 18) => {
    return transactionHistoryService.formatTransactionValue(value, decimals)
  }, [])

  const formatGasPrice = useCallback((gasPrice: string) => {
    return transactionHistoryService.formatGasPrice(gasPrice)
  }, [])

  const getTypeIcon = useCallback((type: TransactionType) => {
    return transactionHistoryService.getTransactionTypeIcon(type)
  }, [])

  const getStatusColor = useCallback((status: TransactionStatus) => {
    return transactionHistoryService.getTransactionStatusColor(status)
  }, [])

  const getCategoryColor = useCallback((category: TransactionCategory) => {
    return transactionHistoryService.getCategoryColor(category)
  }, [])

  // Export transactions
  const exportTransactions = useCallback((format: 'csv' | 'json' = 'csv') => {
    const transactions = state.transactions

    if (format === 'csv') {
      const headers = [
        'Hash',
        'Date',
        'Type',
        'Category',
        'From',
        'To',
        'Value (ETH)',
        'Gas Cost (ETH)',
        'Status',
        'Description'
      ]

      const csvContent = [
        headers.join(','),
        ...transactions.map(tx => [
          tx.hash,
          new Date(tx.timestamp).toISOString(),
          tx.type,
          tx.category,
          tx.from,
          tx.to || '',
          formatValue(tx.value),
          tx.gasCostETH,
          tx.status,
          `"${tx.description || ''}"`
        ].join(','))
      ].join('\n')

      const blob = new Blob([csvContent], { type: 'text/csv' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `transactions-${Date.now()}.csv`
      a.click()
      URL.revokeObjectURL(url)
    } else {
      const jsonContent = JSON.stringify(transactions, null, 2)
      const blob = new Blob([jsonContent], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `transactions-${Date.now()}.json`
      a.click()
      URL.revokeObjectURL(url)
    }

    if (enableNotifications) {
      toast.success('Export Complete', {
        description: `Transactions exported as ${format.toUpperCase()}`
      })
    }
  }, [state.transactions, formatValue, enableNotifications])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Auto-refresh setup
  useEffect(() => {
    if (autoRefresh && address) {
      const interval = setInterval(() => {
        loadTransactions(false, Object.keys(searchFilters).length > 0 ? searchFilters : undefined)
      }, refreshInterval)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, address, refreshInterval, loadTransactions, searchFilters])

  // Initial load
  useEffect(() => {
    if (address) {
      loadTransactions()
      loadSummary()
    }
  }, [address, chainId, loadTransactions, loadSummary])

  return {
    // State
    ...state,

    // Actions
    loadTransactions,
    loadSummary,
    searchTransactions,
    clearSearch,
    loadMore,
    getTransaction,
    exportTransactions,
    clearError,

    // Filters
    getTransactionsByType,
    getTransactionsByCategory,
    getTransactionsByStatus,
    getTransactionsByDateRange,
    getRecentTransactions,
    getFailedTransactions,
    getPendingTransactions,

    // Analytics
    getAnalytics,
    analytics: getAnalytics(),

    // Utilities
    formatValue,
    formatGasPrice,
    getTypeIcon,
    getStatusColor,
    getCategoryColor,

    // Search state
    searchFilters,
    isSearchActive: Object.keys(searchFilters).length > 0,

    // Quick access
    recentTransactions: getRecentTransactions(5),
    failedTransactions: getFailedTransactions(),
    pendingTransactions: getPendingTransactions(),

    // Transaction types
    sendTransactions: getTransactionsByType(TransactionType.SEND),
    receiveTransactions: getTransactionsByType(TransactionType.RECEIVE),
    swapTransactions: getTransactionsByType(TransactionType.SWAP),
    defiTransactions: getTransactionsByCategory(TransactionCategory.DEFI),
    nftTransactions: getTransactionsByCategory(TransactionCategory.NFT),

    // Status counts
    confirmedCount: getTransactionsByStatus(TransactionStatus.CONFIRMED).length,
    pendingCount: getPendingTransactions().length,
    failedCount: getFailedTransactions().length,

    // Summary data
    totalValue: state.summary?.totalValueETH || '0',
    totalGasCost: state.summary?.totalGasCostETH || '0',
    successRate: state.summary?.successRate || 0
  }
}
