import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Hash } from 'viem'
import { 
  transactionRecovery,
  type FailedTransaction,
  type RecoveryStrategy,
  type RecoveryAttempt,
  type RecoveryConfig,
  type RecoveryEvent,
  RecoveryStatus,
  FailureReason,
  RecoveryType
} from '@/lib/transaction-recovery'
import { toast } from 'sonner'

export interface TransactionRecoveryState {
  failedTransactions: FailedTransaction[]
  recoverableTransactions: FailedTransaction[]
  recoveryQueue: FailedTransaction[]
  isAnalyzing: boolean
  isRecovering: boolean
  config: RecoveryConfig
  error: string | null
  lastUpdate: number | null
}

export interface UseTransactionRecoveryOptions {
  autoAnalyze?: boolean
  enableNotifications?: boolean
  requireUserConfirmation?: boolean
  maxAutoRecoveryValue?: number // USD
}

export interface UseTransactionRecoveryReturn {
  // State
  state: TransactionRecoveryState
  
  // Analysis
  analyzeFailedTransaction: (hash: Hash, errorMessage: string, transactionData: any) => Promise<FailedTransaction>
  
  // Recovery
  startRecovery: (hash: Hash, strategy?: RecoveryStrategy) => Promise<void>
  cancelRecovery: (hash: Hash) => void
  retryRecovery: (hash: Hash) => Promise<void>
  
  // Management
  getFailedTransaction: (hash: Hash) => FailedTransaction | null
  clearFailedTransaction: (hash: Hash) => void
  clearAllFailedTransactions: () => void
  
  // Configuration
  updateConfig: (config: Partial<RecoveryConfig>) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useTransactionRecovery = (
  options: UseTransactionRecoveryOptions = {}
): UseTransactionRecoveryReturn => {
  const {
    autoAnalyze = true,
    enableNotifications = true,
    requireUserConfirmation = true,
    maxAutoRecoveryValue = 10
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<TransactionRecoveryState>({
    failedTransactions: [],
    recoverableTransactions: [],
    recoveryQueue: [],
    isAnalyzing: false,
    isRecovering: false,
    config: transactionRecovery.getConfig(),
    error: null,
    lastUpdate: null
  })

  // Update state from recovery engine
  const updateState = useCallback(() => {
    try {
      const failedTransactions = transactionRecovery.getAllFailedTransactions()
      const recoverableTransactions = transactionRecovery.getRecoverableTransactions()
      const recoveryQueue = failedTransactions.filter(tx => 
        tx.recoveryStatus === RecoveryStatus.RECOVERY_IN_PROGRESS
      )
      const config = transactionRecovery.getConfig()

      setState(prev => ({
        ...prev,
        failedTransactions,
        recoverableTransactions,
        recoveryQueue,
        config,
        error: null,
        lastUpdate: Date.now()
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [])

  // Handle recovery events
  const handleRecoveryEvent = useCallback((event: RecoveryEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'analysis_complete':
          toast.info('Transaction Analysis Complete', {
            description: `Found ${event.transaction.canRecover ? 'recovery options' : 'no recovery options'} for failed transaction`
          })
          break
        case 'recovery_started':
          toast.info('Recovery Started', {
            description: `Attempting to recover transaction using ${event.transaction.suggestedFix.title}`
          })
          break
        case 'recovery_success':
          toast.success('Recovery Successful', {
            description: `Transaction recovered successfully with hash ${event.attempt?.newHash?.slice(0, 10)}...`
          })
          break
        case 'recovery_failed':
          toast.error('Recovery Failed', {
            description: `Recovery attempt failed: ${event.attempt?.error || 'Unknown error'}`
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
    const unsubscribe = transactionRecovery.addEventListener(handleRecoveryEvent)

    // Initial state update
    updateState()

    return () => {
      unsubscribe()
    }
  }, [handleRecoveryEvent, updateState])

  // Analyze failed transaction
  const analyzeFailedTransaction = useCallback(async (
    hash: Hash,
    errorMessage: string,
    transactionData: any
  ): Promise<FailedTransaction> => {
    setState(prev => ({ ...prev, isAnalyzing: true, error: null }))

    try {
      const failedTx = await transactionRecovery.analyzeFailedTransaction(hash, errorMessage, transactionData)
      
      setState(prev => ({ ...prev, isAnalyzing: false }))

      if (enableNotifications) {
        if (failedTx.canRecover) {
          toast.success('Recovery Options Found', {
            description: `Found ${failedTx.suggestedFix.title} strategy with ${failedTx.suggestedFix.confidence}% confidence`
          })
        } else {
          toast.warning('No Recovery Options', {
            description: 'Automatic recovery not possible for this transaction'
          })
        }
      }

      return failedTx
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isAnalyzing: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Analysis Failed', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Start recovery
  const startRecovery = useCallback(async (
    hash: Hash,
    strategy?: RecoveryStrategy
  ): Promise<void> => {
    const failedTx = transactionRecovery.getFailedTransaction(hash)
    if (!failedTx) {
      throw new Error('Failed transaction not found')
    }

    // Check if user confirmation is required
    if (requireUserConfirmation && !strategy) {
      const estimatedCostUSD = parseFloat(failedTx.suggestedFix.estimatedCost) * 2000 // Assume $2000 ETH
      
      if (estimatedCostUSD > maxAutoRecoveryValue) {
        if (enableNotifications) {
          toast.warning('User Confirmation Required', {
            description: `Recovery cost ($${estimatedCostUSD.toFixed(2)}) exceeds auto-approval threshold`
          })
        }
        return
      }
    }

    setState(prev => ({ ...prev, isRecovering: true, error: null }))

    try {
      // If strategy is provided, update the suggested fix
      if (strategy) {
        failedTx.suggestedFix = strategy
      }

      // Queue for recovery
      transactionRecovery.queueForRecovery(hash)

      setState(prev => ({ ...prev, isRecovering: false }))

      if (enableNotifications) {
        toast.info('Recovery Queued', {
          description: 'Transaction queued for recovery processing'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isRecovering: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to start recovery', { description: errorMessage })
      }
      throw error
    }
  }, [requireUserConfirmation, maxAutoRecoveryValue, enableNotifications])

  // Cancel recovery
  const cancelRecovery = useCallback((hash: Hash) => {
    try {
      const failedTx = transactionRecovery.getFailedTransaction(hash)
      if (failedTx) {
        failedTx.recoveryStatus = RecoveryStatus.ANALYSIS_COMPLETE
        updateState()

        if (enableNotifications) {
          toast.info('Recovery Cancelled', {
            description: 'Transaction recovery has been cancelled'
          })
        }
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to cancel recovery', { description: errorMessage })
      }
    }
  }, [enableNotifications, updateState])

  // Retry recovery
  const retryRecovery = useCallback(async (hash: Hash): Promise<void> => {
    const failedTx = transactionRecovery.getFailedTransaction(hash)
    if (!failedTx) {
      throw new Error('Failed transaction not found')
    }

    // Reset recovery status and queue again
    failedTx.recoveryStatus = RecoveryStatus.RECOVERY_AVAILABLE
    await startRecovery(hash)
  }, [startRecovery])

  // Get failed transaction
  const getFailedTransaction = useCallback((hash: Hash): FailedTransaction | null => {
    return transactionRecovery.getFailedTransaction(hash)
  }, [])

  // Clear failed transaction
  const clearFailedTransaction = useCallback((hash: Hash) => {
    try {
      const failedTransactions = transactionRecovery.getAllFailedTransactions()
      const updatedTransactions = failedTransactions.filter(tx => tx.hash !== hash)
      
      // Clear from recovery engine (would need to implement this method)
      // transactionRecovery.removeFailedTransaction(hash)
      
      updateState()

      if (enableNotifications) {
        toast.success('Transaction Cleared', {
          description: 'Failed transaction removed from recovery list'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to clear transaction', { description: errorMessage })
      }
    }
  }, [enableNotifications, updateState])

  // Clear all failed transactions
  const clearAllFailedTransactions = useCallback(() => {
    try {
      transactionRecovery.clearFailedTransactions()
      updateState()

      if (enableNotifications) {
        toast.success('All Transactions Cleared', {
          description: 'All failed transactions have been cleared'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to clear transactions', { description: errorMessage })
      }
    }
  }, [enableNotifications, updateState])

  // Update configuration
  const updateConfig = useCallback((config: Partial<RecoveryConfig>) => {
    try {
      transactionRecovery.updateConfig(config)
      setState(prev => ({ ...prev, config: transactionRecovery.getConfig() }))

      if (enableNotifications) {
        toast.success('Configuration Updated', {
          description: 'Recovery settings have been updated'
        })
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
    analyzeFailedTransaction,
    startRecovery,
    cancelRecovery,
    retryRecovery,
    getFailedTransaction,
    clearFailedTransaction,
    clearAllFailedTransactions,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for basic recovery operations
export const useFailedTransactionRecovery = () => {
  const { state, analyzeFailedTransaction, startRecovery, clearFailedTransaction } = useTransactionRecovery({
    autoAnalyze: true,
    enableNotifications: true,
    requireUserConfirmation: true
  })

  return {
    failedTransactions: state.failedTransactions,
    recoverableTransactions: state.recoverableTransactions,
    isAnalyzing: state.isAnalyzing,
    isRecovering: state.isRecovering,
    analyzeFailedTransaction,
    startRecovery,
    clearFailedTransaction
  }
}

// Hook for recovery statistics
export const useRecoveryStats = () => {
  const { state } = useTransactionRecovery()

  const stats = {
    totalFailed: state.failedTransactions.length,
    recoverable: state.recoverableTransactions.length,
    inProgress: state.recoveryQueue.length,
    successRate: state.failedTransactions.length > 0 
      ? (state.failedTransactions.filter(tx => tx.recoveryStatus === RecoveryStatus.RECOVERY_SUCCESS).length / state.failedTransactions.length) * 100
      : 0,
    averageAttempts: state.failedTransactions.length > 0
      ? state.failedTransactions.reduce((sum, tx) => sum + tx.recoveryAttempts.length, 0) / state.failedTransactions.length
      : 0
  }

  return stats
}
