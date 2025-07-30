import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Hash } from 'viem'
import { 
  transactionPool,
  type QueuedTransaction,
  type QueueStats,
  type QueueConfig,
  type QueueEvent,
  type TransactionMetadata,
  TransactionPriority,
  QueueStatus
} from '@/lib/transaction-pool'
import { toast } from 'sonner'

export interface TransactionPoolState {
  transactions: QueuedTransaction[]
  queuedTransactions: QueuedTransaction[]
  pendingTransactions: QueuedTransaction[]
  submittedTransactions: QueuedTransaction[]
  confirmedTransactions: QueuedTransaction[]
  failedTransactions: QueuedTransaction[]
  stats: QueueStats
  config: QueueConfig
  isProcessing: boolean
  error: string | null
}

export interface UseTransactionPoolOptions {
  autoStart?: boolean
  enableNotifications?: boolean
  maxDisplayTransactions?: number
  filterByAddress?: boolean
}

export interface UseTransactionPoolReturn {
  // State
  state: TransactionPoolState
  
  // Queue Management
  addTransaction: (transaction: Omit<QueuedTransaction, 'id' | 'status' | 'retryCount' | 'createdAt'>) => string
  removeTransaction: (id: string) => boolean
  cancelTransaction: (id: string) => boolean
  updatePriority: (id: string, priority: TransactionPriority) => boolean
  
  // Transaction Access
  getTransaction: (id: string) => QueuedTransaction | null
  getTransactionsByStatus: (status: QueueStatus) => QueuedTransaction[]
  getTransactionsByAddress: (address: string) => QueuedTransaction[]
  
  // Queue Control
  startProcessing: () => void
  stopProcessing: () => void
  clearCompleted: () => void
  
  // Configuration
  updateConfig: (config: Partial<QueueConfig>) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useTransactionPool = (
  options: UseTransactionPoolOptions = {}
): UseTransactionPoolReturn => {
  const {
    autoStart = true,
    enableNotifications = true,
    maxDisplayTransactions = 50,
    filterByAddress = true
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<TransactionPoolState>({
    transactions: [],
    queuedTransactions: [],
    pendingTransactions: [],
    submittedTransactions: [],
    confirmedTransactions: [],
    failedTransactions: [],
    stats: {
      totalTransactions: 0,
      queuedTransactions: 0,
      pendingTransactions: 0,
      submittedTransactions: 0,
      confirmedTransactions: 0,
      failedTransactions: 0,
      cancelledTransactions: 0,
      avgConfirmationTime: 0,
      avgRetryCount: 0,
      successRate: 0
    },
    config: transactionPool.getConfig(),
    isProcessing: false,
    error: null
  })

  // Update state from transaction pool
  const updateState = useCallback(() => {
    try {
      let allTransactions = Array.from(transactionPool['queue'].values())
      
      // Filter by address if enabled and address is available
      if (filterByAddress && address) {
        allTransactions = allTransactions.filter(tx => 
          tx.from.toLowerCase() === address.toLowerCase()
        )
      }

      // Limit number of transactions displayed
      const transactions = allTransactions
        .sort((a, b) => b.createdAt - a.createdAt)
        .slice(0, maxDisplayTransactions)

      const queuedTransactions = transactions.filter(tx => tx.status === QueueStatus.QUEUED)
      const pendingTransactions = transactions.filter(tx => tx.status === QueueStatus.PENDING)
      const submittedTransactions = transactions.filter(tx => tx.status === QueueStatus.SUBMITTED)
      const confirmedTransactions = transactions.filter(tx => tx.status === QueueStatus.CONFIRMED)
      const failedTransactions = transactions.filter(tx => tx.status === QueueStatus.FAILED)

      const stats = transactionPool.getStats()
      const config = transactionPool.getConfig()

      setState(prev => ({
        ...prev,
        transactions,
        queuedTransactions,
        pendingTransactions,
        submittedTransactions,
        confirmedTransactions,
        failedTransactions,
        stats,
        config,
        error: null
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [address, filterByAddress, maxDisplayTransactions])

  // Handle queue events
  const handleQueueEvent = useCallback((event: QueueEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'transaction_added':
          toast.info('Transaction Added to Queue', {
            description: `Transaction ${event.transaction.id.slice(0, 8)}... queued for processing`
          })
          break
        case 'transaction_submitted':
          toast.success('Transaction Submitted', {
            description: `Transaction ${event.transaction.hash?.slice(0, 10)}... submitted to network`
          })
          break
        case 'transaction_failed':
          toast.error('Transaction Failed', {
            description: `Transaction ${event.transaction.id.slice(0, 8)}... failed: ${event.transaction.error}`
          })
          break
        case 'transaction_retry':
          toast.warning('Transaction Retry', {
            description: `Retrying transaction ${event.transaction.id.slice(0, 8)}... (attempt ${event.transaction.retryCount})`
          })
          break
        case 'transaction_cancelled':
          toast.info('Transaction Cancelled', {
            description: `Transaction ${event.transaction.id.slice(0, 8)}... was cancelled`
          })
          break
        case 'priority_boosted':
          toast.info('Priority Boosted', {
            description: `Transaction ${event.transaction.id.slice(0, 8)}... priority increased due to age`
          })
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = transactionPool.addEventListener(handleQueueEvent)

    // Start processing if auto-start is enabled
    if (autoStart) {
      startProcessing()
    }

    // Initial state update
    updateState()

    return () => {
      unsubscribe()
    }
  }, [autoStart, handleQueueEvent, updateState])

  // Add transaction to queue
  const addTransaction = useCallback((
    transaction: Omit<QueuedTransaction, 'id' | 'status' | 'retryCount' | 'createdAt'>
  ): string => {
    try {
      // Set chain ID if not provided
      const txWithChain = {
        ...transaction,
        chainId: transaction.chainId || chainId || 1
      }

      const id = transactionPool.addTransaction(txWithChain)
      return id
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to add transaction', { description: errorMessage })
      }
      throw error
    }
  }, [chainId, enableNotifications])

  // Remove transaction from queue
  const removeTransaction = useCallback((id: string): boolean => {
    try {
      const success = transactionPool.removeTransaction(id)
      if (success && enableNotifications) {
        toast.success('Transaction removed from queue')
      }
      return success
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to remove transaction', { description: errorMessage })
      }
      return false
    }
  }, [enableNotifications])

  // Cancel transaction
  const cancelTransaction = useCallback((id: string): boolean => {
    try {
      const success = transactionPool.cancelTransaction(id)
      if (!success && enableNotifications) {
        toast.warning('Cannot cancel transaction', {
          description: 'Transaction may already be processed or confirmed'
        })
      }
      return success
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to cancel transaction', { description: errorMessage })
      }
      return false
    }
  }, [enableNotifications])

  // Update transaction priority
  const updatePriority = useCallback((id: string, priority: TransactionPriority): boolean => {
    try {
      const success = transactionPool.updatePriority(id, priority)
      if (success && enableNotifications) {
        const priorityNames = ['Low', 'Normal', 'High', 'Urgent']
        toast.success('Priority Updated', {
          description: `Transaction priority set to ${priorityNames[priority]}`
        })
      }
      return success
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to update priority', { description: errorMessage })
      }
      return false
    }
  }, [enableNotifications])

  // Get transaction by ID
  const getTransaction = useCallback((id: string): QueuedTransaction | null => {
    return transactionPool.getTransaction(id)
  }, [])

  // Get transactions by status
  const getTransactionsByStatus = useCallback((status: QueueStatus): QueuedTransaction[] => {
    return transactionPool.getTransactionsByStatus(status)
  }, [])

  // Get transactions by address
  const getTransactionsByAddress = useCallback((address: string): QueuedTransaction[] => {
    return transactionPool.getTransactionsByAddress(address)
  }, [])

  // Start processing
  const startProcessing = useCallback(() => {
    try {
      transactionPool['startProcessing']()
      setState(prev => ({ ...prev, isProcessing: true }))
      
      if (enableNotifications) {
        toast.success('Transaction processing started')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to start processing', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Stop processing
  const stopProcessing = useCallback(() => {
    try {
      transactionPool.stopProcessing()
      setState(prev => ({ ...prev, isProcessing: false }))
      
      if (enableNotifications) {
        toast.info('Transaction processing stopped')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to stop processing', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Clear completed transactions
  const clearCompleted = useCallback(() => {
    try {
      transactionPool.clearCompleted()
      updateState()
      
      if (enableNotifications) {
        toast.success('Completed transactions cleared')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to clear completed transactions', { description: errorMessage })
      }
    }
  }, [enableNotifications, updateState])

  // Update configuration
  const updateConfig = useCallback((config: Partial<QueueConfig>) => {
    try {
      transactionPool.updateConfig(config)
      setState(prev => ({ ...prev, config: transactionPool.getConfig() }))
      
      if (enableNotifications) {
        toast.success('Configuration updated')
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (enableNotifications) {
        toast.error('Failed to update configuration', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
  }, [updateState])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    addTransaction,
    removeTransaction,
    cancelTransaction,
    updatePriority,
    getTransaction,
    getTransactionsByStatus,
    getTransactionsByAddress,
    startProcessing,
    stopProcessing,
    clearCompleted,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for basic queue operations
export const useTransactionQueue = () => {
  const { state, addTransaction, cancelTransaction, updatePriority } = useTransactionPool({
    autoStart: true,
    enableNotifications: true,
    filterByAddress: true
  })

  return {
    queuedTransactions: state.queuedTransactions,
    pendingTransactions: state.pendingTransactions,
    stats: state.stats,
    addTransaction,
    cancelTransaction,
    updatePriority,
    isProcessing: state.isProcessing
  }
}
