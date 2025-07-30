import { type Hash } from 'viem'

export interface QueuedTransaction {
  id: string
  hash?: Hash
  chainId: number
  from: string
  to?: string
  value: string
  data?: string
  gasLimit: string
  gasPrice?: string
  maxFeePerGas?: string
  maxPriorityFeePerGas?: string
  nonce?: number
  priority: TransactionPriority
  status: QueueStatus
  retryCount: number
  maxRetries: number
  createdAt: number
  scheduledAt?: number
  submittedAt?: number
  confirmedAt?: number
  failedAt?: number
  lastRetryAt?: number
  dependencies?: string[] // Transaction IDs this depends on
  metadata?: TransactionMetadata
  error?: string
}

export interface TransactionMetadata {
  type: 'send' | 'swap' | 'approve' | 'stake' | 'unstake' | 'mint' | 'burn' | 'other'
  description?: string
  dappName?: string
  tokenSymbol?: string
  tokenAmount?: string
  usdValue?: number
  deadline?: number
  slippage?: number
  userInitiated: boolean
}

export enum TransactionPriority {
  LOW = 0,
  NORMAL = 1,
  HIGH = 2,
  URGENT = 3
}

export enum QueueStatus {
  QUEUED = 'queued',
  PENDING = 'pending',
  SUBMITTED = 'submitted',
  CONFIRMED = 'confirmed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired'
}

export interface QueueConfig {
  maxQueueSize: number
  maxConcurrentTransactions: number
  defaultRetryCount: number
  retryDelayMs: number
  maxRetryDelayMs: number
  transactionTimeoutMs: number
  nonceGapTolerance: number
  priorityBoostThreshold: number
  enableAutomaticRetry: boolean
  enableNonceManagement: boolean
  enableGasOptimization: boolean
}

export interface QueueStats {
  totalTransactions: number
  queuedTransactions: number
  pendingTransactions: number
  submittedTransactions: number
  confirmedTransactions: number
  failedTransactions: number
  cancelledTransactions: number
  avgConfirmationTime: number
  avgRetryCount: number
  successRate: number
}

export interface NonceManager {
  chainId: number
  address: string
  currentNonce: number
  pendingNonce: number
  lastUpdate: number
  gaps: number[]
}

export class TransactionPool {
  private static instance: TransactionPool
  private queue = new Map<string, QueuedTransaction>()
  private nonceManagers = new Map<string, NonceManager>()
  private processingQueue = new Set<string>()
  private config: QueueConfig
  private isProcessing = false
  private processingInterval?: NodeJS.Timeout
  private eventListeners = new Set<(event: QueueEvent) => void>()

  private constructor() {
    this.config = {
      maxQueueSize: 100,
      maxConcurrentTransactions: 5,
      defaultRetryCount: 3,
      retryDelayMs: 5000,
      maxRetryDelayMs: 60000,
      transactionTimeoutMs: 300000, // 5 minutes
      nonceGapTolerance: 5,
      priorityBoostThreshold: 60000, // 1 minute
      enableAutomaticRetry: true,
      enableNonceManagement: true,
      enableGasOptimization: true
    }
  }

  static getInstance(): TransactionPool {
    if (!TransactionPool.instance) {
      TransactionPool.instance = new TransactionPool()
    }
    return TransactionPool.instance
  }

  /**
   * Add transaction to queue
   */
  addTransaction(transaction: Omit<QueuedTransaction, 'id' | 'status' | 'retryCount' | 'createdAt'>): string {
    // Check queue size limit
    if (this.queue.size >= this.config.maxQueueSize) {
      throw new Error('Transaction queue is full')
    }

    const id = `tx_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`
    const queuedTransaction: QueuedTransaction = {
      ...transaction,
      id,
      status: QueueStatus.QUEUED,
      retryCount: 0,
      maxRetries: transaction.maxRetries || this.config.defaultRetryCount,
      createdAt: Date.now()
    }

    // Assign nonce if nonce management is enabled
    if (this.config.enableNonceManagement && !queuedTransaction.nonce) {
      queuedTransaction.nonce = this.getNextNonce(queuedTransaction.chainId, queuedTransaction.from)
    }

    this.queue.set(id, queuedTransaction)

    // Emit event
    this.emitEvent({
      type: 'transaction_added',
      transaction: queuedTransaction,
      timestamp: Date.now()
    })

    // Start processing if not already running
    if (!this.isProcessing) {
      this.startProcessing()
    }

    return id
  }

  /**
   * Remove transaction from queue
   */
  removeTransaction(id: string): boolean {
    const transaction = this.queue.get(id)
    if (!transaction) return false

    // Can only remove queued or failed transactions
    if (![QueueStatus.QUEUED, QueueStatus.FAILED, QueueStatus.CANCELLED].includes(transaction.status)) {
      throw new Error('Cannot remove transaction that is being processed')
    }

    this.queue.delete(id)
    this.processingQueue.delete(id)

    // Emit event
    this.emitEvent({
      type: 'transaction_removed',
      transaction,
      timestamp: Date.now()
    })

    return true
  }

  /**
   * Cancel transaction
   */
  cancelTransaction(id: string): boolean {
    const transaction = this.queue.get(id)
    if (!transaction) return false

    // Can only cancel queued or pending transactions
    if (![QueueStatus.QUEUED, QueueStatus.PENDING].includes(transaction.status)) {
      return false
    }

    transaction.status = QueueStatus.CANCELLED
    this.processingQueue.delete(id)

    // Emit event
    this.emitEvent({
      type: 'transaction_cancelled',
      transaction,
      timestamp: Date.now()
    })

    return true
  }

  /**
   * Update transaction priority
   */
  updatePriority(id: string, priority: TransactionPriority): boolean {
    const transaction = this.queue.get(id)
    if (!transaction) return false

    // Can only update priority for queued transactions
    if (transaction.status !== QueueStatus.QUEUED) {
      return false
    }

    transaction.priority = priority

    // Emit event
    this.emitEvent({
      type: 'priority_updated',
      transaction,
      timestamp: Date.now()
    })

    return true
  }

  /**
   * Get transaction by ID
   */
  getTransaction(id: string): QueuedTransaction | null {
    return this.queue.get(id) || null
  }

  /**
   * Get transactions by status
   */
  getTransactionsByStatus(status: QueueStatus): QueuedTransaction[] {
    return Array.from(this.queue.values())
      .filter(tx => tx.status === status)
      .sort((a, b) => this.compareTransactionPriority(a, b))
  }

  /**
   * Get transactions by address
   */
  getTransactionsByAddress(address: string): QueuedTransaction[] {
    return Array.from(this.queue.values())
      .filter(tx => tx.from.toLowerCase() === address.toLowerCase())
      .sort((a, b) => b.createdAt - a.createdAt)
  }

  /**
   * Get queue statistics
   */
  getStats(): QueueStats {
    const transactions = Array.from(this.queue.values())
    const total = transactions.length
    
    const statusCounts = transactions.reduce((acc, tx) => {
      acc[tx.status] = (acc[tx.status] || 0) + 1
      return acc
    }, {} as Record<QueueStatus, number>)

    const confirmedTxs = transactions.filter(tx => tx.status === QueueStatus.CONFIRMED)
    const avgConfirmationTime = confirmedTxs.length > 0
      ? confirmedTxs.reduce((sum, tx) => sum + (tx.confirmedAt! - tx.createdAt), 0) / confirmedTxs.length
      : 0

    const avgRetryCount = transactions.length > 0
      ? transactions.reduce((sum, tx) => sum + tx.retryCount, 0) / transactions.length
      : 0

    const successRate = total > 0 
      ? ((statusCounts[QueueStatus.CONFIRMED] || 0) / total) * 100 
      : 0

    return {
      totalTransactions: total,
      queuedTransactions: statusCounts[QueueStatus.QUEUED] || 0,
      pendingTransactions: statusCounts[QueueStatus.PENDING] || 0,
      submittedTransactions: statusCounts[QueueStatus.SUBMITTED] || 0,
      confirmedTransactions: statusCounts[QueueStatus.CONFIRMED] || 0,
      failedTransactions: statusCounts[QueueStatus.FAILED] || 0,
      cancelledTransactions: statusCounts[QueueStatus.CANCELLED] || 0,
      avgConfirmationTime,
      avgRetryCount,
      successRate
    }
  }

  /**
   * Start processing queue
   */
  private startProcessing(): void {
    if (this.isProcessing) return

    this.isProcessing = true
    this.processingInterval = setInterval(() => {
      this.processQueue()
    }, 1000) // Process every second

    console.log('Transaction pool processing started')
  }

  /**
   * Stop processing queue
   */
  stopProcessing(): void {
    if (!this.isProcessing) return

    this.isProcessing = false
    if (this.processingInterval) {
      clearInterval(this.processingInterval)
      this.processingInterval = undefined
    }

    console.log('Transaction pool processing stopped')
  }

  /**
   * Process transaction queue
   */
  private async processQueue(): Promise<void> {
    try {
      // Get queued transactions sorted by priority
      const queuedTransactions = this.getTransactionsByStatus(QueueStatus.QUEUED)
      
      // Check for expired transactions
      this.checkExpiredTransactions()
      
      // Boost priority for old transactions
      this.boostOldTransactionPriority()
      
      // Process transactions up to concurrent limit
      const availableSlots = this.config.maxConcurrentTransactions - this.processingQueue.size
      const transactionsToProcess = queuedTransactions.slice(0, availableSlots)

      for (const transaction of transactionsToProcess) {
        if (this.canProcessTransaction(transaction)) {
          this.processTransaction(transaction)
        }
      }

      // Retry failed transactions
      if (this.config.enableAutomaticRetry) {
        this.retryFailedTransactions()
      }

    } catch (error) {
      console.error('Error processing transaction queue:', error)
    }
  }

  /**
   * Check if transaction can be processed
   */
  private canProcessTransaction(transaction: QueuedTransaction): boolean {
    // Check dependencies
    if (transaction.dependencies) {
      for (const depId of transaction.dependencies) {
        const dep = this.queue.get(depId)
        if (!dep || dep.status !== QueueStatus.CONFIRMED) {
          return false
        }
      }
    }

    // Check nonce order if nonce management is enabled
    if (this.config.enableNonceManagement && transaction.nonce !== undefined) {
      const nonceManager = this.getNonceManager(transaction.chainId, transaction.from)
      if (transaction.nonce > nonceManager.currentNonce + this.config.nonceGapTolerance) {
        return false
      }
    }

    return true
  }

  /**
   * Process individual transaction
   */
  private async processTransaction(transaction: QueuedTransaction): Promise<void> {
    this.processingQueue.add(transaction.id)
    transaction.status = QueueStatus.PENDING
    transaction.submittedAt = Date.now()

    try {
      // Emit event
      this.emitEvent({
        type: 'transaction_processing',
        transaction,
        timestamp: Date.now()
      })

      // Submit transaction (mock implementation)
      const hash = await this.submitTransaction(transaction)
      
      transaction.hash = hash
      transaction.status = QueueStatus.SUBMITTED

      // Emit event
      this.emitEvent({
        type: 'transaction_submitted',
        transaction,
        timestamp: Date.now()
      })

      // Update nonce manager
      if (this.config.enableNonceManagement && transaction.nonce !== undefined) {
        this.updateNonceManager(transaction.chainId, transaction.from, transaction.nonce)
      }

    } catch (error) {
      transaction.status = QueueStatus.FAILED
      transaction.error = (error as Error).message
      transaction.failedAt = Date.now()

      // Emit event
      this.emitEvent({
        type: 'transaction_failed',
        transaction,
        timestamp: Date.now()
      })

      console.error(`Transaction ${transaction.id} failed:`, error)
    } finally {
      this.processingQueue.delete(transaction.id)
    }
  }

  /**
   * Submit transaction to network (mock implementation)
   */
  private async submitTransaction(transaction: QueuedTransaction): Promise<Hash> {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))
    
    // Simulate 90% success rate
    if (Math.random() < 0.9) {
      return `0x${Math.random().toString(16).substr(2, 64)}` as Hash
    } else {
      throw new Error('Transaction failed to submit')
    }
  }

  /**
   * Retry failed transactions
   */
  private retryFailedTransactions(): void {
    const failedTransactions = this.getTransactionsByStatus(QueueStatus.FAILED)
    
    for (const transaction of failedTransactions) {
      if (this.shouldRetryTransaction(transaction)) {
        this.retryTransaction(transaction)
      }
    }
  }

  /**
   * Check if transaction should be retried
   */
  private shouldRetryTransaction(transaction: QueuedTransaction): boolean {
    if (transaction.retryCount >= transaction.maxRetries) {
      return false
    }

    const timeSinceLastRetry = Date.now() - (transaction.lastRetryAt || transaction.failedAt!)
    const retryDelay = Math.min(
      this.config.retryDelayMs * Math.pow(2, transaction.retryCount),
      this.config.maxRetryDelayMs
    )

    return timeSinceLastRetry >= retryDelay
  }

  /**
   * Retry transaction
   */
  private retryTransaction(transaction: QueuedTransaction): void {
    transaction.status = QueueStatus.QUEUED
    transaction.retryCount++
    transaction.lastRetryAt = Date.now()
    transaction.error = undefined

    // Emit event
    this.emitEvent({
      type: 'transaction_retry',
      transaction,
      timestamp: Date.now()
    })
  }

  /**
   * Check for expired transactions
   */
  private checkExpiredTransactions(): void {
    const now = Date.now()
    
    for (const transaction of Array.from(this.queue.values())) {
      if (transaction.status === QueueStatus.PENDING || transaction.status === QueueStatus.SUBMITTED) {
        const age = now - transaction.createdAt
        if (age > this.config.transactionTimeoutMs) {
          transaction.status = QueueStatus.EXPIRED
          this.processingQueue.delete(transaction.id)

          // Emit event
          this.emitEvent({
            type: 'transaction_expired',
            transaction,
            timestamp: Date.now()
          })
        }
      }
    }
  }

  /**
   * Boost priority for old transactions
   */
  private boostOldTransactionPriority(): void {
    const now = Date.now()
    
    for (const transaction of Array.from(this.queue.values())) {
      if (transaction.status === QueueStatus.QUEUED) {
        const age = now - transaction.createdAt
        if (age > this.config.priorityBoostThreshold && transaction.priority < TransactionPriority.HIGH) {
          transaction.priority = Math.min(transaction.priority + 1, TransactionPriority.URGENT)
          
          // Emit event
          this.emitEvent({
            type: 'priority_boosted',
            transaction,
            timestamp: Date.now()
          })
        }
      }
    }
  }

  /**
   * Compare transaction priority for sorting
   */
  private compareTransactionPriority(a: QueuedTransaction, b: QueuedTransaction): number {
    // First by priority (higher first)
    if (a.priority !== b.priority) {
      return b.priority - a.priority
    }
    
    // Then by creation time (older first)
    return a.createdAt - b.createdAt
  }

  /**
   * Get next nonce for address
   */
  private getNextNonce(chainId: number, address: string): number {
    const key = `${chainId}_${address.toLowerCase()}`
    const manager = this.nonceManagers.get(key)
    
    if (!manager) {
      // Initialize nonce manager
      const newManager: NonceManager = {
        chainId,
        address: address.toLowerCase(),
        currentNonce: 0, // Would fetch from network
        pendingNonce: 0,
        lastUpdate: Date.now(),
        gaps: []
      }
      this.nonceManagers.set(key, newManager)
      return 0
    }

    return manager.pendingNonce++
  }

  /**
   * Get nonce manager
   */
  private getNonceManager(chainId: number, address: string): NonceManager {
    const key = `${chainId}_${address.toLowerCase()}`
    let manager = this.nonceManagers.get(key)
    
    if (!manager) {
      manager = {
        chainId,
        address: address.toLowerCase(),
        currentNonce: 0,
        pendingNonce: 0,
        lastUpdate: Date.now(),
        gaps: []
      }
      this.nonceManagers.set(key, manager)
    }
    
    return manager
  }

  /**
   * Update nonce manager
   */
  private updateNonceManager(chainId: number, address: string, nonce: number): void {
    const manager = this.getNonceManager(chainId, address)
    
    if (nonce === manager.currentNonce) {
      manager.currentNonce++
      manager.lastUpdate = Date.now()
      
      // Fill any gaps
      while (manager.gaps.includes(manager.currentNonce)) {
        manager.gaps = manager.gaps.filter(n => n !== manager.currentNonce)
        manager.currentNonce++
      }
    } else if (nonce > manager.currentNonce) {
      // Add gap
      for (let i = manager.currentNonce; i < nonce; i++) {
        if (!manager.gaps.includes(i)) {
          manager.gaps.push(i)
        }
      }
      manager.currentNonce = nonce + 1
      manager.lastUpdate = Date.now()
    }
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: QueueEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in queue event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: QueueEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Update configuration
   */
  updateConfig(config: Partial<QueueConfig>): void {
    this.config = { ...this.config, ...config }
  }

  /**
   * Get configuration
   */
  getConfig(): QueueConfig {
    return { ...this.config }
  }

  /**
   * Clear completed transactions
   */
  clearCompleted(): void {
    const toDelete: string[] = []
    
    for (const [id, transaction] of Array.from(this.queue.entries())) {
      if ([QueueStatus.CONFIRMED, QueueStatus.FAILED, QueueStatus.CANCELLED, QueueStatus.EXPIRED].includes(transaction.status)) {
        toDelete.push(id)
      }
    }
    
    toDelete.forEach(id => this.queue.delete(id))
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.stopProcessing()
    this.queue.clear()
    this.nonceManagers.clear()
    this.processingQueue.clear()
    this.eventListeners.clear()
  }
}

export interface QueueEvent {
  type: 'transaction_added' | 'transaction_removed' | 'transaction_cancelled' | 'priority_updated' | 
        'transaction_processing' | 'transaction_submitted' | 'transaction_failed' | 'transaction_expired' |
        'transaction_retry' | 'priority_boosted'
  transaction: QueuedTransaction
  timestamp: number
}

// Export singleton instance
export const transactionPool = TransactionPool.getInstance()
