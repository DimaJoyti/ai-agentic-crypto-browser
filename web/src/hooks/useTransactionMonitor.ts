import { useState, useEffect, useCallback, useRef } from 'react'
import { type Hash } from 'viem'
import { 
  transactionMonitor, 
  type TransactionData, 
  type TransactionUpdate, 
  type TransactionCallback,
  TransactionStatus,
  TransactionType 
} from '@/lib/transaction-monitor'
import { toast } from 'sonner'

export interface UseTransactionMonitorOptions {
  showNotifications?: boolean
  autoRemoveCompleted?: boolean
  autoRemoveAfter?: number // milliseconds
}

export interface TransactionState {
  transactions: Map<Hash, TransactionData>
  isLoading: boolean
  error: string | null
}

export function useTransactionMonitor(options: UseTransactionMonitorOptions = {}) {
  const {
    showNotifications = true,
    autoRemoveCompleted = true,
    autoRemoveAfter = 30000 // 30 seconds
  } = options

  const [state, setState] = useState<TransactionState>({
    transactions: new Map(),
    isLoading: false,
    error: null
  })

  const callbacksRef = useRef<Map<Hash, () => void>>(new Map())
  const timeoutsRef = useRef<Map<Hash, NodeJS.Timeout>>(new Map())

  // Update state with current transactions
  const updateTransactions = useCallback(() => {
    const allTransactions = transactionMonitor.getAllTransactions()
    const transactionMap = new Map<Hash, TransactionData>()
    
    allTransactions.forEach(tx => {
      transactionMap.set(tx.hash, tx)
    })

    setState(prev => ({
      ...prev,
      transactions: transactionMap
    }))
  }, [])

  // Track a new transaction
  const trackTransaction = useCallback(async (
    hash: Hash,
    chainId: number,
    type: TransactionType = TransactionType.SEND,
    metadata?: TransactionData['metadata']
  ) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      await transactionMonitor.trackTransaction(hash, chainId, type, metadata)
      
      // Set up callback for this transaction
      const callback: TransactionCallback = (update: TransactionUpdate) => {
        updateTransactions()
        
        if (showNotifications) {
          handleTransactionNotification(update, metadata)
        }

        // Auto-remove completed transactions
        if (autoRemoveCompleted && 
            (update.status === TransactionStatus.CONFIRMED || 
             update.status === TransactionStatus.FAILED ||
             update.status === TransactionStatus.DROPPED)) {
          
          const timeout = setTimeout(() => {
            stopTracking(hash)
          }, autoRemoveAfter)
          
          timeoutsRef.current.set(hash, timeout)
        }
      }

      const unsubscribe = transactionMonitor.onTransactionUpdate(hash, callback)
      callbacksRef.current.set(hash, unsubscribe)

      updateTransactions()

      if (showNotifications) {
        toast.info(`Transaction submitted: ${formatHash(hash)}`, {
          description: `Tracking transaction on ${getChainName(chainId)}`,
          action: {
            label: 'View',
            onClick: () => openBlockExplorer(hash, chainId)
          }
        })
      }

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to track transaction'
      setState(prev => ({ ...prev, error: errorMessage }))
      
      if (showNotifications) {
        toast.error('Failed to track transaction', {
          description: errorMessage
        })
      }
    } finally {
      setState(prev => ({ ...prev, isLoading: false }))
    }
  }, [showNotifications, autoRemoveCompleted, autoRemoveAfter, updateTransactions])

  // Stop tracking a transaction
  const stopTracking = useCallback((hash: Hash) => {
    // Clear timeout
    const timeout = timeoutsRef.current.get(hash)
    if (timeout) {
      clearTimeout(timeout)
      timeoutsRef.current.delete(hash)
    }

    // Unsubscribe from updates
    const unsubscribe = callbacksRef.current.get(hash)
    if (unsubscribe) {
      unsubscribe()
      callbacksRef.current.delete(hash)
    }

    // Stop monitoring
    transactionMonitor.stopTracking(hash)
    updateTransactions()
  }, [updateTransactions])

  // Get transaction by hash
  const getTransaction = useCallback((hash: Hash): TransactionData | undefined => {
    return state.transactions.get(hash)
  }, [state.transactions])

  // Get transactions by status
  const getTransactionsByStatus = useCallback((status: TransactionStatus): TransactionData[] => {
    return Array.from(state.transactions.values()).filter(tx => tx.status === status)
  }, [state.transactions])

  // Get transactions by chain
  const getTransactionsByChain = useCallback((chainId: number): TransactionData[] => {
    return Array.from(state.transactions.values()).filter(tx => tx.chainId === chainId)
  }, [state.transactions])

  // Clear all transactions
  const clearAllTransactions = useCallback(() => {
    // Clear all timeouts
    timeoutsRef.current.forEach(timeout => clearTimeout(timeout))
    timeoutsRef.current.clear()

    // Unsubscribe from all callbacks
    callbacksRef.current.forEach(unsubscribe => unsubscribe())
    callbacksRef.current.clear()

    // Clear all tracked transactions
    Array.from(state.transactions.keys()).forEach(hash => {
      transactionMonitor.stopTracking(hash)
    })

    updateTransactions()
  }, [state.transactions, updateTransactions])

  // Handle transaction notifications
  const handleTransactionNotification = useCallback((
    update: TransactionUpdate, 
    metadata?: TransactionData['metadata']
  ) => {
    const hash = formatHash(update.hash)
    const description = metadata?.description || `Transaction ${hash}`

    switch (update.status) {
      case TransactionStatus.CONFIRMED:
        toast.success('Transaction confirmed!', {
          description: `${description} has been confirmed with ${update.confirmations} confirmations`,
          action: {
            label: 'View',
            onClick: () => openBlockExplorer(update.hash, getTransactionChainId(update.hash))
          }
        })
        break

      case TransactionStatus.FAILED:
        toast.error('Transaction failed', {
          description: `${description} failed: ${update.error || 'Unknown error'}`,
          action: {
            label: 'View',
            onClick: () => openBlockExplorer(update.hash, getTransactionChainId(update.hash))
          }
        })
        break

      case TransactionStatus.DROPPED:
        toast.warning('Transaction dropped', {
          description: `${description} was dropped from the mempool`
        })
        break
    }
  }, [])

  // Helper function to get chain ID for a transaction
  const getTransactionChainId = useCallback((hash: Hash): number => {
    const tx = state.transactions.get(hash)
    return tx?.chainId || 1
  }, [state.transactions])

  // Initialize and cleanup
  useEffect(() => {
    updateTransactions()

    return () => {
      // Cleanup on unmount
      timeoutsRef.current.forEach(timeout => clearTimeout(timeout))
      callbacksRef.current.forEach(unsubscribe => unsubscribe())
    }
  }, [updateTransactions])

  // Start real-time monitoring for a chain
  const startRealtimeMonitoring = useCallback((chainId: number) => {
    transactionMonitor.startRealtimeMonitoring(chainId)
    if (showNotifications) {
      toast.info(`Real-time monitoring started for ${getChainName(chainId)}`)
    }
  }, [showNotifications])

  // Stop real-time monitoring for a chain
  const stopRealtimeMonitoring = useCallback((chainId: number) => {
    transactionMonitor.stopRealtimeMonitoring(chainId)
    if (showNotifications) {
      toast.info(`Real-time monitoring stopped for ${getChainName(chainId)}`)
    }
  }, [showNotifications])

  // Get transaction statistics
  const getStats = useCallback(() => {
    return transactionMonitor.getTransactionStats()
  }, [])

  // Retry a failed transaction
  const retryTransaction = useCallback(async (hash: Hash, newGasPrice?: string) => {
    try {
      setState(prev => ({ ...prev, isLoading: true }))
      const newHash = await transactionMonitor.retryTransaction(hash, newGasPrice)

      if (showNotifications) {
        toast.info('Transaction retry submitted', {
          description: `New transaction: ${formatHash(newHash)}`
        })
      }

      return newHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to retry transaction'
      setState(prev => ({ ...prev, error: errorMessage }))

      if (showNotifications) {
        toast.error('Failed to retry transaction', { description: errorMessage })
      }
      throw error
    } finally {
      setState(prev => ({ ...prev, isLoading: false }))
    }
  }, [showNotifications])

  // Cancel a pending transaction
  const cancelTransaction = useCallback(async (hash: Hash) => {
    try {
      setState(prev => ({ ...prev, isLoading: true }))
      const cancelHash = await transactionMonitor.cancelTransaction(hash)

      if (showNotifications) {
        toast.info('Cancellation transaction submitted', {
          description: `Cancel transaction: ${formatHash(cancelHash)}`
        })
      }

      return cancelHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to cancel transaction'
      setState(prev => ({ ...prev, error: errorMessage }))

      if (showNotifications) {
        toast.error('Failed to cancel transaction', { description: errorMessage })
      }
      throw error
    } finally {
      setState(prev => ({ ...prev, isLoading: false }))
    }
  }, [showNotifications])

  // Speed up a pending transaction
  const speedUpTransaction = useCallback(async (hash: Hash, newGasPrice: string) => {
    try {
      setState(prev => ({ ...prev, isLoading: true }))
      const speedUpHash = await transactionMonitor.speedUpTransaction(hash, newGasPrice)

      if (showNotifications) {
        toast.info('Speed up transaction submitted', {
          description: `Speed up transaction: ${formatHash(speedUpHash)}`
        })
      }

      return speedUpHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to speed up transaction'
      setState(prev => ({ ...prev, error: errorMessage }))

      if (showNotifications) {
        toast.error('Failed to speed up transaction', { description: errorMessage })
      }
      throw error
    } finally {
      setState(prev => ({ ...prev, isLoading: false }))
    }
  }, [showNotifications])

  return {
    // State
    transactions: Array.from(state.transactions.values()),
    transactionMap: state.transactions,
    isLoading: state.isLoading,
    error: state.error,

    // Actions
    trackTransaction,
    stopTracking,
    clearAllTransactions,
    retryTransaction,
    cancelTransaction,
    speedUpTransaction,

    // Real-time monitoring
    startRealtimeMonitoring,
    stopRealtimeMonitoring,

    // Getters
    getTransaction,
    getTransactionsByStatus,
    getTransactionsByChain,
    getStats,

    // Computed values
    pendingTransactions: getTransactionsByStatus(TransactionStatus.PENDING),
    confirmedTransactions: getTransactionsByStatus(TransactionStatus.CONFIRMED),
    failedTransactions: getTransactionsByStatus(TransactionStatus.FAILED),
    totalTransactions: state.transactions.size
  }
}

// Helper functions
function formatHash(hash: Hash): string {
  return `${hash.slice(0, 6)}...${hash.slice(-4)}`
}

function getChainName(chainId: number): string {
  const chains: Record<number, string> = {
    1: 'Ethereum',
    137: 'Polygon',
    42161: 'Arbitrum',
    10: 'Optimism',
    8453: 'Base',
    11155111: 'Sepolia'
  }
  return chains[chainId] || `Chain ${chainId}`
}

function openBlockExplorer(hash: Hash, chainId: number): void {
  const explorers: Record<number, string> = {
    1: 'https://etherscan.io/tx/',
    137: 'https://polygonscan.com/tx/',
    42161: 'https://arbiscan.io/tx/',
    10: 'https://optimistic.etherscan.io/tx/',
    8453: 'https://basescan.org/tx/',
    11155111: 'https://sepolia.etherscan.io/tx/'
  }
  
  const baseUrl = explorers[chainId] || 'https://etherscan.io/tx/'
  window.open(`${baseUrl}${hash}`, '_blank')
}
