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

    // This would connect to a WebSocket endpoint for real-time updates
    // For now, we'll use polling, but this is where WebSocket logic would go
    console.log(`WebSocket support for ${chain.name} would be initialized here`)
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
