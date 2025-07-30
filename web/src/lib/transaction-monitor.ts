import { createPublicClient, http, parseAbiItem, formatEther, type Hash, type TransactionReceipt } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum TransactionStatus {
  PENDING = 'pending',
  CONFIRMED = 'confirmed',
  FAILED = 'failed',
  DROPPED = 'dropped',
  REPLACED = 'replaced'
}

export enum TransactionType {
  SEND = 'send',
  RECEIVE = 'receive',
  SWAP = 'swap',
  APPROVE = 'approve',
  STAKE = 'stake',
  UNSTAKE = 'unstake',
  MINT = 'mint',
  BURN = 'burn',
  CONTRACT_INTERACTION = 'contract_interaction'
}

export interface TransactionData {
  hash: Hash
  chainId: number
  from: string
  to?: string
  value: string
  gasPrice?: string
  gasLimit?: string
  gasUsed?: string
  nonce?: number
  blockNumber?: number
  blockHash?: string
  transactionIndex?: number
  status: TransactionStatus
  type: TransactionType
  timestamp: number
  confirmations: number
  maxConfirmations: number
  receipt?: TransactionReceipt
  error?: string
  metadata?: {
    tokenSymbol?: string
    tokenAmount?: string
    contractAddress?: string
    methodName?: string
    description?: string
  }
}

export interface TransactionUpdate {
  hash: Hash
  status: TransactionStatus
  confirmations: number
  blockNumber?: number
  gasUsed?: string
  receipt?: TransactionReceipt
  error?: string
  timestamp: number
}

export type TransactionCallback = (update: TransactionUpdate) => void

export class TransactionMonitor {
  private static instance: TransactionMonitor
  private clients: Map<number, any> = new Map()
  private trackedTransactions: Map<Hash, TransactionData> = new Map()
  private callbacks: Map<Hash, TransactionCallback[]> = new Map()
  private pollingIntervals: Map<Hash, NodeJS.Timeout> = new Map()
  private websockets: Map<number, WebSocket> = new Map()

  private constructor() {
    this.initializeClients()
  }

  static getInstance(): TransactionMonitor {
    if (!TransactionMonitor.instance) {
      TransactionMonitor.instance = new TransactionMonitor()
    }
    return TransactionMonitor.instance
  }

  private initializeClients() {
    // Initialize viem clients for each supported chain
    Object.values(SUPPORTED_CHAINS).forEach(chain => {
      if (!chain.isTestnet || chain.id === 11155111) { // Include Sepolia for testing
        try {
          const client = createPublicClient({
            chain: {
              id: chain.id,
              name: chain.name,
              network: chain.shortName.toLowerCase(),
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: chain.rpcUrls
            } as any,
            transport: http()
          })
          this.clients.set(chain.id, client)
        } catch (error) {
          console.warn(`Failed to initialize client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  /**
   * Start tracking a transaction
   */
  async trackTransaction(
    hash: Hash,
    chainId: number,
    type: TransactionType = TransactionType.SEND,
    metadata?: TransactionData['metadata']
  ): Promise<void> {
    const client = this.clients.get(chainId)
    if (!client) {
      throw new Error(`No client available for chain ${chainId}`)
    }

    try {
      // Get initial transaction data
      const tx = await client.getTransaction({ hash })
      
      const transactionData: TransactionData = {
        hash,
        chainId,
        from: tx.from,
        to: tx.to,
        value: formatEther(tx.value),
        gasPrice: tx.gasPrice ? formatEther(tx.gasPrice) : undefined,
        gasLimit: tx.gas?.toString(),
        nonce: tx.nonce,
        status: TransactionStatus.PENDING,
        type,
        timestamp: Date.now(),
        confirmations: 0,
        maxConfirmations: this.getMaxConfirmations(chainId),
        metadata
      }

      this.trackedTransactions.set(hash, transactionData)
      this.startPolling(hash, chainId)

      // Try to get receipt immediately in case it's already mined
      this.checkTransactionStatus(hash, chainId)
    } catch (error) {
      console.error(`Failed to track transaction ${hash}:`, error)
      throw error
    }
  }

  /**
   * Add callback for transaction updates
   */
  onTransactionUpdate(hash: Hash, callback: TransactionCallback): () => void {
    if (!this.callbacks.has(hash)) {
      this.callbacks.set(hash, [])
    }
    this.callbacks.get(hash)!.push(callback)

    // Return unsubscribe function
    return () => {
      const callbacks = this.callbacks.get(hash)
      if (callbacks) {
        const index = callbacks.indexOf(callback)
        if (index > -1) {
          callbacks.splice(index, 1)
        }
      }
    }
  }

  /**
   * Get transaction data
   */
  getTransaction(hash: Hash): TransactionData | undefined {
    return this.trackedTransactions.get(hash)
  }

  /**
   * Get all tracked transactions
   */
  getAllTransactions(): TransactionData[] {
    return Array.from(this.trackedTransactions.values())
  }

  /**
   * Stop tracking a transaction
   */
  stopTracking(hash: Hash): void {
    const interval = this.pollingIntervals.get(hash)
    if (interval) {
      clearInterval(interval)
      this.pollingIntervals.delete(hash)
    }
    this.trackedTransactions.delete(hash)
    this.callbacks.delete(hash)
  }

  private startPolling(hash: Hash, chainId: number): void {
    const interval = setInterval(() => {
      this.checkTransactionStatus(hash, chainId)
    }, 3000) // Poll every 3 seconds

    this.pollingIntervals.set(hash, interval)

    // Stop polling after 30 minutes if transaction is still pending
    setTimeout(() => {
      const tx = this.trackedTransactions.get(hash)
      if (tx && tx.status === TransactionStatus.PENDING) {
        this.updateTransactionStatus(hash, {
          hash,
          status: TransactionStatus.DROPPED,
          confirmations: 0,
          timestamp: Date.now(),
          error: 'Transaction dropped after 30 minutes'
        })
        this.stopTracking(hash)
      }
    }, 30 * 60 * 1000)
  }

  private async checkTransactionStatus(hash: Hash, chainId: number): Promise<void> {
    const client = this.clients.get(chainId)
    const tx = this.trackedTransactions.get(hash)
    
    if (!client || !tx) return

    try {
      // Try to get transaction receipt
      const receipt = await client.getTransactionReceipt({ hash })
      
      if (receipt) {
        // Transaction is mined
        const currentBlock = await client.getBlockNumber()
        const confirmations = Number(currentBlock - receipt.blockNumber)
        
        const status = receipt.status === 'success' 
          ? TransactionStatus.CONFIRMED 
          : TransactionStatus.FAILED

        const update: TransactionUpdate = {
          hash,
          status,
          confirmations,
          blockNumber: Number(receipt.blockNumber),
          gasUsed: receipt.gasUsed.toString(),
          receipt,
          timestamp: Date.now()
        }

        this.updateTransactionStatus(hash, update)

        // Stop polling if we have enough confirmations or transaction failed
        if (confirmations >= tx.maxConfirmations || status === TransactionStatus.FAILED) {
          this.stopTracking(hash)
        }
      }
    } catch (error) {
      // Transaction might not be mined yet, continue polling
      console.debug(`Transaction ${hash} not yet mined:`, error)
    }
  }

  private updateTransactionStatus(hash: Hash, update: TransactionUpdate): void {
    const tx = this.trackedTransactions.get(hash)
    if (!tx) return

    // Update transaction data
    const updatedTx: TransactionData = {
      ...tx,
      status: update.status,
      confirmations: update.confirmations,
      blockNumber: update.blockNumber,
      gasUsed: update.gasUsed,
      receipt: update.receipt,
      error: update.error
    }

    this.trackedTransactions.set(hash, updatedTx)

    // Notify callbacks
    const callbacks = this.callbacks.get(hash)
    if (callbacks) {
      callbacks.forEach(callback => {
        try {
          callback(update)
        } catch (error) {
          console.error('Error in transaction callback:', error)
        }
      })
    }
  }

  private getMaxConfirmations(chainId: number): number {
    // Different chains have different confirmation requirements
    switch (chainId) {
      case 1: // Ethereum
        return 12
      case 137: // Polygon
        return 20
      case 42161: // Arbitrum
        return 1
      case 10: // Optimism
        return 1
      case 8453: // Base
        return 1
      default:
        return 6
    }
  }

  /**
   * Initialize WebSocket connection for real-time updates
   */
  private initializeWebSocket(chainId: number): void {
    const chain = SUPPORTED_CHAINS[chainId]
    if (!chain) return

    // For now, skip WebSocket initialization as wsRpcUrl is not available
    // In a real implementation, you would need to add wsRpcUrl to the chain config
    return

    // WebSocket functionality disabled - requires wsRpcUrl in chain config
    /*
    try {
      const ws = new WebSocket(chain.wsRpcUrl)

      ws.onopen = () => {
        console.log(`WebSocket connected for ${chain.name}`)

        // Subscribe to pending transactions
        ws.send(JSON.stringify({
          id: 1,
          method: 'eth_subscribe',
          params: ['newPendingTransactions']
        }))

        // Subscribe to new blocks
        ws.send(JSON.stringify({
          id: 2,
          method: 'eth_subscribe',
          params: ['newHeads']
        }))
      }

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.handleWebSocketMessage(chainId, data)
        } catch (error) {
          console.error('Error parsing WebSocket message:', error)
        }
      }

      ws.onclose = () => {
        console.log(`WebSocket disconnected for ${chain.name}`)
        this.websockets.delete(chainId)

        // Attempt to reconnect after 5 seconds
        setTimeout(() => {
          this.initializeWebSocket(chainId)
        }, 5000)
      }

      ws.onerror = (error) => {
        console.error(`WebSocket error for ${chain.name}:`, error)
      }

      this.websockets.set(chainId, ws)
    } catch (error) {
      console.error(`Failed to initialize WebSocket for ${chain.name}:`, error)
    }
    */
  }

  /**
   * Handle WebSocket messages
   */
  private handleWebSocketMessage(chainId: number, data: any): void {
    if (data.method === 'eth_subscription') {
      const { subscription, result } = data.params

      // Handle new block
      if (result.number) {
        this.handleNewBlock(chainId, result)
      }

      // Handle new pending transaction
      if (typeof result === 'string') {
        this.handleNewPendingTransaction(chainId, result)
      }
    }
  }

  /**
   * Handle new block for confirmation updates
   */
  private handleNewBlock(chainId: number, block: any): void {
    // Update confirmations for all pending transactions on this chain
    for (const [hash, transaction] of Array.from(this.trackedTransactions.entries())) {
      if (transaction.chainId === chainId &&
          transaction.status === TransactionStatus.CONFIRMED &&
          transaction.blockNumber) {
        const confirmations = Number(block.number) - transaction.blockNumber + 1

        if (confirmations !== transaction.confirmations) {
          this.updateTransactionStatus(hash, {
            hash,
            status: transaction.status,
            confirmations,
            timestamp: Date.now()
          })
        }
      }
    }
  }

  /**
   * Handle new pending transaction
   */
  private handleNewPendingTransaction(chainId: number, txHash: string): void {
    const transaction = this.trackedTransactions.get(txHash as Hash)
    if (transaction && transaction.status === TransactionStatus.PENDING) {
      console.log(`Tracked transaction ${txHash} appeared in mempool`)
      // Transaction appeared in mempool, could trigger notification
    }
  }

  /**
   * Start real-time monitoring for a chain
   */
  startRealtimeMonitoring(chainId: number): void {
    if (!this.websockets.has(chainId)) {
      this.initializeWebSocket(chainId)
    }
  }

  /**
   * Stop real-time monitoring for a chain
   */
  stopRealtimeMonitoring(chainId: number): void {
    const ws = this.websockets.get(chainId)
    if (ws) {
      ws.close()
      this.websockets.delete(chainId)
    }
  }

  /**
   * Estimate transaction type based on transaction data
   */
  static estimateTransactionType(tx: any): TransactionType {
    if (!tx.to) return TransactionType.CONTRACT_INTERACTION
    if (tx.value && BigInt(tx.value) > BigInt(0)) return TransactionType.SEND
    if (tx.input && tx.input !== '0x') {
      // Try to decode common method signatures
      const methodId = tx.input.slice(0, 10)
      switch (methodId) {
        case '0xa9059cbb': // transfer(address,uint256)
          return TransactionType.SEND
        case '0x095ea7b3': // approve(address,uint256)
          return TransactionType.APPROVE
        case '0x38ed1739': // swapExactTokensForTokens
        case '0x7ff36ab5': // swapExactETHForTokens
          return TransactionType.SWAP
        default:
          return TransactionType.CONTRACT_INTERACTION
      }
    }
    return TransactionType.SEND
  }

  /**
   * Get transactions by status
   */
  getTransactionsByStatus(status: TransactionStatus): TransactionData[] {
    return Array.from(this.trackedTransactions.values())
      .filter(tx => tx.status === status)
      .sort((a, b) => b.timestamp - a.timestamp)
  }

  /**
   * Get transactions by address
   */
  getTransactionsByAddress(address: string): TransactionData[] {
    const normalizedAddress = address.toLowerCase()
    return Array.from(this.trackedTransactions.values())
      .filter(tx =>
        tx.from.toLowerCase() === normalizedAddress ||
        tx.to?.toLowerCase() === normalizedAddress
      )
      .sort((a, b) => b.timestamp - a.timestamp)
  }

  /**
   * Get transactions by chain
   */
  getTransactionsByChain(chainId: number): TransactionData[] {
    return Array.from(this.trackedTransactions.values())
      .filter(tx => tx.chainId === chainId)
      .sort((a, b) => b.timestamp - a.timestamp)
  }

  /**
   * Get transaction statistics
   */
  getTransactionStats(): {
    total: number
    pending: number
    confirmed: number
    failed: number
    dropped: number
    avgConfirmationTime: number
  } {
    const transactions = Array.from(this.trackedTransactions.values())
    const total = transactions.length
    const pending = transactions.filter(tx => tx.status === TransactionStatus.PENDING).length
    const confirmed = transactions.filter(tx => tx.status === TransactionStatus.CONFIRMED).length
    const failed = transactions.filter(tx => tx.status === TransactionStatus.FAILED).length
    const dropped = transactions.filter(tx => tx.status === TransactionStatus.DROPPED).length

    // Calculate average confirmation time for confirmed transactions
    const confirmedTxs = transactions.filter(tx =>
      tx.status === TransactionStatus.CONFIRMED && tx.receipt
    )
    const avgConfirmationTime = confirmedTxs.length > 0
      ? confirmedTxs.reduce((sum, tx) => {
          // Estimate confirmation time based on block time
          const blockTime = this.getAverageBlockTime(tx.chainId)
          return sum + (tx.confirmations * blockTime)
        }, 0) / confirmedTxs.length
      : 0

    return {
      total,
      pending,
      confirmed,
      failed,
      dropped,
      avgConfirmationTime
    }
  }

  /**
   * Get average block time for a chain
   */
  private getAverageBlockTime(chainId: number): number {
    switch (chainId) {
      case 1: // Ethereum
        return 12
      case 137: // Polygon
        return 2
      case 42161: // Arbitrum
        return 0.25
      case 10: // Optimism
        return 2
      case 8453: // Base
        return 2
      case 56: // BSC
        return 3
      default:
        return 12
    }
  }

  /**
   * Retry a failed transaction
   */
  async retryTransaction(hash: Hash, newGasPrice?: string): Promise<Hash> {
    const tx = this.trackedTransactions.get(hash)
    if (!tx || tx.status !== TransactionStatus.FAILED) {
      throw new Error('Transaction not found or not failed')
    }

    // This would implement transaction retry logic
    // For now, just return the original hash
    console.log(`Retrying transaction ${hash} with new gas price: ${newGasPrice}`)
    return hash
  }

  /**
   * Cancel a pending transaction
   */
  async cancelTransaction(hash: Hash): Promise<Hash> {
    const tx = this.trackedTransactions.get(hash)
    if (!tx || tx.status !== TransactionStatus.PENDING) {
      throw new Error('Transaction not found or not pending')
    }

    // This would implement transaction cancellation logic
    console.log(`Cancelling transaction ${hash}`)
    return hash
  }

  /**
   * Speed up a pending transaction
   */
  async speedUpTransaction(hash: Hash, newGasPrice: string): Promise<Hash> {
    const tx = this.trackedTransactions.get(hash)
    if (!tx || tx.status !== TransactionStatus.PENDING) {
      throw new Error('Transaction not found or not pending')
    }

    // This would implement transaction speed up logic
    console.log(`Speeding up transaction ${hash} with gas price: ${newGasPrice}`)
    return hash
  }

  /**
   * Clean up old transactions
   */
  cleanupOldTransactions(maxAge: number = 24 * 60 * 60 * 1000): void {
    const cutoff = Date.now() - maxAge
    const toDelete: Hash[] = []

    for (const [hash, tx] of Array.from(this.trackedTransactions.entries())) {
      if (tx.timestamp < cutoff && tx.status !== TransactionStatus.PENDING) {
        toDelete.push(hash)
      }
    }

    toDelete.forEach(hash => {
      this.stopTracking(hash)
    })

    console.log(`Cleaned up ${toDelete.length} old transactions`)
  }

  /**
   * Clean up resources
   */
  destroy(): void {
    // Clear all intervals
    this.pollingIntervals.forEach(interval => clearInterval(interval))
    this.pollingIntervals.clear()

    // Close WebSocket connections
    this.websockets.forEach(ws => ws.close())
    this.websockets.clear()

    // Clear data
    this.trackedTransactions.clear()
    this.callbacks.clear()
  }
}

// Export singleton instance
export const transactionMonitor = TransactionMonitor.getInstance()
