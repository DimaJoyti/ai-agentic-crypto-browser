import { type Hash } from 'viem'

export interface FailedTransaction {
  hash: Hash
  chainId: number
  from: string
  to?: string
  value: string
  gasLimit: string
  gasPrice: string
  maxFeePerGas?: string
  maxPriorityFeePerGas?: string
  nonce: number
  data?: string
  timestamp: number
  failureReason: FailureReason
  failureDetails: FailureDetails
  recoveryAttempts: RecoveryAttempt[]
  recoveryStatus: RecoveryStatus
  canRecover: boolean
  suggestedFix: RecoveryStrategy
  metadata?: TransactionMetadata
}

export interface FailureDetails {
  errorCode?: string
  errorMessage: string
  gasUsed?: string
  revertReason?: string
  blockNumber?: number
  blockHash?: string
  transactionIndex?: number
  logs?: any[]
  internalTransactions?: any[]
  debugTrace?: any
}

export interface RecoveryAttempt {
  id: string
  timestamp: number
  strategy: RecoveryStrategy
  newHash?: Hash
  status: 'pending' | 'success' | 'failed'
  gasAdjustment?: {
    oldGasPrice: string
    newGasPrice: string
    oldGasLimit: string
    newGasLimit: string
  }
  error?: string
  cost?: string
}

export interface RecoveryStrategy {
  type: RecoveryType
  title: string
  description: string
  confidence: number // 0-100
  estimatedCost: string
  estimatedTime: number // seconds
  riskLevel: 'low' | 'medium' | 'high'
  parameters: RecoveryParameters
  prerequisites?: string[]
  warnings?: string[]
}

export interface RecoveryParameters {
  gasPrice?: string
  gasLimit?: string
  maxFeePerGas?: string
  maxPriorityFeePerGas?: string
  nonce?: number
  data?: string
  value?: string
  deadline?: number
  slippage?: number
  replacementTx?: boolean
  speedUp?: boolean
}

export enum FailureReason {
  OUT_OF_GAS = 'out_of_gas',
  INSUFFICIENT_FUNDS = 'insufficient_funds',
  NONCE_TOO_LOW = 'nonce_too_low',
  NONCE_TOO_HIGH = 'nonce_too_high',
  GAS_PRICE_TOO_LOW = 'gas_price_too_low',
  TRANSACTION_UNDERPRICED = 'transaction_underpriced',
  REPLACEMENT_UNDERPRICED = 'replacement_underpriced',
  INTRINSIC_GAS_TOO_LOW = 'intrinsic_gas_too_low',
  EXECUTION_REVERTED = 'execution_reverted',
  INVALID_SIGNATURE = 'invalid_signature',
  NETWORK_ERROR = 'network_error',
  TIMEOUT = 'timeout',
  SLIPPAGE_EXCEEDED = 'slippage_exceeded',
  DEADLINE_EXCEEDED = 'deadline_exceeded',
  INSUFFICIENT_LIQUIDITY = 'insufficient_liquidity',
  CONTRACT_ERROR = 'contract_error',
  UNKNOWN = 'unknown'
}

export enum RecoveryType {
  INCREASE_GAS_PRICE = 'increase_gas_price',
  INCREASE_GAS_LIMIT = 'increase_gas_limit',
  ADJUST_NONCE = 'adjust_nonce',
  RETRY_TRANSACTION = 'retry_transaction',
  REPLACE_TRANSACTION = 'replace_transaction',
  SPEED_UP_TRANSACTION = 'speed_up_transaction',
  CANCEL_TRANSACTION = 'cancel_transaction',
  ADJUST_SLIPPAGE = 'adjust_slippage',
  EXTEND_DEADLINE = 'extend_deadline',
  SPLIT_TRANSACTION = 'split_transaction',
  BATCH_TRANSACTION = 'batch_transaction',
  ALTERNATIVE_ROUTE = 'alternative_route',
  MANUAL_INTERVENTION = 'manual_intervention'
}

export enum RecoveryStatus {
  PENDING_ANALYSIS = 'pending_analysis',
  ANALYSIS_COMPLETE = 'analysis_complete',
  RECOVERY_AVAILABLE = 'recovery_available',
  RECOVERY_IN_PROGRESS = 'recovery_in_progress',
  RECOVERY_SUCCESS = 'recovery_success',
  RECOVERY_FAILED = 'recovery_failed',
  MANUAL_REQUIRED = 'manual_required',
  NOT_RECOVERABLE = 'not_recoverable'
}

export interface TransactionMetadata {
  type: 'send' | 'swap' | 'approve' | 'stake' | 'unstake' | 'mint' | 'burn' | 'other'
  dappName?: string
  tokenSymbol?: string
  tokenAmount?: string
  swapDetails?: {
    tokenIn: string
    tokenOut: string
    amountIn: string
    expectedAmountOut: string
    slippage: number
    deadline: number
  }
}

export interface RecoveryConfig {
  enableAutoRecovery: boolean
  maxRecoveryAttempts: number
  gasIncreasePercentage: number
  maxGasIncrease: number
  retryDelay: number
  maxRetryDelay: number
  timeoutThreshold: number
  enableNotifications: boolean
  requireUserConfirmation: boolean
  autoApproveThreshold: number // USD value
}

export class TransactionRecoveryEngine {
  private static instance: TransactionRecoveryEngine
  private failedTransactions = new Map<Hash, FailedTransaction>()
  private recoveryQueue = new Set<Hash>()
  private config: RecoveryConfig
  private isProcessing = false
  private processingInterval?: NodeJS.Timeout
  private eventListeners = new Set<(event: RecoveryEvent) => void>()

  private constructor() {
    this.config = {
      enableAutoRecovery: true,
      maxRecoveryAttempts: 3,
      gasIncreasePercentage: 20,
      maxGasIncrease: 200,
      retryDelay: 30000, // 30 seconds
      maxRetryDelay: 300000, // 5 minutes
      timeoutThreshold: 600000, // 10 minutes
      enableNotifications: true,
      requireUserConfirmation: true,
      autoApproveThreshold: 10 // $10 USD
    }
  }

  static getInstance(): TransactionRecoveryEngine {
    if (!TransactionRecoveryEngine.instance) {
      TransactionRecoveryEngine.instance = new TransactionRecoveryEngine()
    }
    return TransactionRecoveryEngine.instance
  }

  /**
   * Analyze failed transaction and determine recovery options
   */
  async analyzeFailedTransaction(
    hash: Hash,
    errorMessage: string,
    transactionData: any
  ): Promise<FailedTransaction> {
    const failureReason = this.classifyFailure(errorMessage)
    const failureDetails = await this.getFailureDetails(hash, errorMessage)
    const recoveryStrategies = this.generateRecoveryStrategies(failureReason, transactionData, failureDetails)
    
    const failedTx: FailedTransaction = {
      hash,
      chainId: transactionData.chainId,
      from: transactionData.from,
      to: transactionData.to,
      value: transactionData.value || '0',
      gasLimit: transactionData.gasLimit || transactionData.gas,
      gasPrice: transactionData.gasPrice || '0',
      maxFeePerGas: transactionData.maxFeePerGas,
      maxPriorityFeePerGas: transactionData.maxPriorityFeePerGas,
      nonce: transactionData.nonce,
      data: transactionData.data,
      timestamp: Date.now(),
      failureReason,
      failureDetails,
      recoveryAttempts: [],
      recoveryStatus: RecoveryStatus.ANALYSIS_COMPLETE,
      canRecover: recoveryStrategies.length > 0,
      suggestedFix: recoveryStrategies[0] || this.getManualInterventionStrategy(),
      metadata: transactionData.metadata
    }

    this.failedTransactions.set(hash, failedTx)

    // Emit analysis complete event
    this.emitEvent({
      type: 'analysis_complete',
      transaction: failedTx,
      timestamp: Date.now()
    })

    // Auto-queue for recovery if enabled
    if (this.config.enableAutoRecovery && failedTx.canRecover) {
      this.queueForRecovery(hash)
    }

    return failedTx
  }

  /**
   * Classify failure reason from error message
   */
  private classifyFailure(errorMessage: string): FailureReason {
    const message = errorMessage.toLowerCase()

    if (message.includes('out of gas') || message.includes('gas required exceeds allowance')) {
      return FailureReason.OUT_OF_GAS
    }
    if (message.includes('insufficient funds') || message.includes('insufficient balance')) {
      return FailureReason.INSUFFICIENT_FUNDS
    }
    if (message.includes('nonce too low') || message.includes('nonce has already been used')) {
      return FailureReason.NONCE_TOO_LOW
    }
    if (message.includes('nonce too high') || message.includes('nonce gap')) {
      return FailureReason.NONCE_TOO_HIGH
    }
    if (message.includes('gas price too low') || message.includes('underpriced')) {
      return FailureReason.GAS_PRICE_TOO_LOW
    }
    if (message.includes('replacement transaction underpriced')) {
      return FailureReason.REPLACEMENT_UNDERPRICED
    }
    if (message.includes('intrinsic gas too low')) {
      return FailureReason.INTRINSIC_GAS_TOO_LOW
    }
    if (message.includes('execution reverted') || message.includes('revert')) {
      return FailureReason.EXECUTION_REVERTED
    }
    if (message.includes('invalid signature')) {
      return FailureReason.INVALID_SIGNATURE
    }
    if (message.includes('network') || message.includes('connection')) {
      return FailureReason.NETWORK_ERROR
    }
    if (message.includes('timeout') || message.includes('deadline')) {
      return FailureReason.TIMEOUT
    }
    if (message.includes('slippage') || message.includes('price impact')) {
      return FailureReason.SLIPPAGE_EXCEEDED
    }
    if (message.includes('insufficient liquidity')) {
      return FailureReason.INSUFFICIENT_LIQUIDITY
    }

    return FailureReason.UNKNOWN
  }

  /**
   * Get detailed failure information
   */
  private async getFailureDetails(hash: Hash, errorMessage: string): Promise<FailureDetails> {
    try {
      // In a real implementation, this would:
      // 1. Fetch transaction receipt
      // 2. Get debug trace if available
      // 3. Parse revert reason
      // 4. Analyze gas usage
      
      return {
        errorMessage,
        errorCode: this.extractErrorCode(errorMessage),
        revertReason: this.extractRevertReason(errorMessage)
      }
    } catch (error) {
      return {
        errorMessage,
        errorCode: 'UNKNOWN'
      }
    }
  }

  /**
   * Extract error code from error message
   */
  private extractErrorCode(errorMessage: string): string | undefined {
    const codeMatch = errorMessage.match(/code[:\s]+([A-Z_]+)/i)
    return codeMatch?.[1]
  }

  /**
   * Extract revert reason from error message
   */
  private extractRevertReason(errorMessage: string): string | undefined {
    const revertMatch = errorMessage.match(/revert[:\s]+(.+)/i)
    return revertMatch?.[1]?.trim()
  }

  /**
   * Generate recovery strategies based on failure reason
   */
  private generateRecoveryStrategies(
    failureReason: FailureReason,
    transactionData: any,
    failureDetails: FailureDetails
  ): RecoveryStrategy[] {
    const strategies: RecoveryStrategy[] = []

    switch (failureReason) {
      case FailureReason.OUT_OF_GAS:
        strategies.push(this.createGasLimitIncreaseStrategy(transactionData))
        break

      case FailureReason.GAS_PRICE_TOO_LOW:
      case FailureReason.TRANSACTION_UNDERPRICED:
        strategies.push(this.createGasPriceIncreaseStrategy(transactionData))
        strategies.push(this.createSpeedUpStrategy(transactionData))
        break

      case FailureReason.REPLACEMENT_UNDERPRICED:
        strategies.push(this.createReplacementStrategy(transactionData))
        break

      case FailureReason.NONCE_TOO_LOW:
        strategies.push(this.createNonceAdjustmentStrategy(transactionData, 'increase'))
        break

      case FailureReason.NONCE_TOO_HIGH:
        strategies.push(this.createNonceAdjustmentStrategy(transactionData, 'decrease'))
        break

      case FailureReason.SLIPPAGE_EXCEEDED:
        strategies.push(this.createSlippageAdjustmentStrategy(transactionData))
        break

      case FailureReason.DEADLINE_EXCEEDED:
        strategies.push(this.createDeadlineExtensionStrategy(transactionData))
        break

      case FailureReason.INSUFFICIENT_FUNDS:
        strategies.push(this.createSplitTransactionStrategy(transactionData))
        break

      case FailureReason.NETWORK_ERROR:
      case FailureReason.TIMEOUT:
        strategies.push(this.createRetryStrategy(transactionData))
        break

      case FailureReason.EXECUTION_REVERTED:
        strategies.push(this.createContractAnalysisStrategy(transactionData, failureDetails))
        break

      default:
        strategies.push(this.createGenericRetryStrategy(transactionData))
        break
    }

    // Sort by confidence level
    return strategies.sort((a, b) => b.confidence - a.confidence)
  }

  /**
   * Create gas limit increase strategy
   */
  private createGasLimitIncreaseStrategy(transactionData: any): RecoveryStrategy {
    const currentGasLimit = BigInt(transactionData.gasLimit || transactionData.gas)
    const newGasLimit = (currentGasLimit * BigInt(150) / BigInt(100)).toString() // 50% increase

    return {
      type: RecoveryType.INCREASE_GAS_LIMIT,
      title: 'Increase Gas Limit',
      description: 'Transaction failed due to insufficient gas. Increasing gas limit by 50%.',
      confidence: 85,
      estimatedCost: this.calculateGasCostIncrease(transactionData, newGasLimit),
      estimatedTime: 60,
      riskLevel: 'low',
      parameters: {
        gasLimit: newGasLimit
      },
      warnings: ['Higher gas limit will increase transaction cost']
    }
  }

  /**
   * Create gas price increase strategy
   */
  private createGasPriceIncreaseStrategy(transactionData: any): RecoveryStrategy {
    const currentGasPrice = BigInt(transactionData.gasPrice || '0')
    const increasePercentage = this.config.gasIncreasePercentage
    const newGasPrice = (currentGasPrice * BigInt(100 + increasePercentage) / BigInt(100)).toString()

    return {
      type: RecoveryType.INCREASE_GAS_PRICE,
      title: 'Increase Gas Price',
      description: `Transaction underpriced. Increasing gas price by ${increasePercentage}%.`,
      confidence: 90,
      estimatedCost: this.calculateGasCostIncrease(transactionData, undefined, newGasPrice),
      estimatedTime: 30,
      riskLevel: 'low',
      parameters: {
        gasPrice: newGasPrice
      }
    }
  }

  /**
   * Create speed up strategy
   */
  private createSpeedUpStrategy(transactionData: any): RecoveryStrategy {
    return {
      type: RecoveryType.SPEED_UP_TRANSACTION,
      title: 'Speed Up Transaction',
      description: 'Send a replacement transaction with higher gas price to speed up confirmation.',
      confidence: 95,
      estimatedCost: this.calculateSpeedUpCost(transactionData),
      estimatedTime: 15,
      riskLevel: 'low',
      parameters: {
        speedUp: true,
        replacementTx: true
      }
    }
  }

  /**
   * Create replacement strategy
   */
  private createReplacementStrategy(transactionData: any): RecoveryStrategy {
    return {
      type: RecoveryType.REPLACE_TRANSACTION,
      title: 'Replace Transaction',
      description: 'Replace the stuck transaction with a new one using higher gas price.',
      confidence: 88,
      estimatedCost: this.calculateReplacementCost(transactionData),
      estimatedTime: 45,
      riskLevel: 'medium',
      parameters: {
        replacementTx: true
      },
      warnings: ['Original transaction will be cancelled']
    }
  }

  /**
   * Create nonce adjustment strategy
   */
  private createNonceAdjustmentStrategy(transactionData: any, direction: 'increase' | 'decrease'): RecoveryStrategy {
    const currentNonce = transactionData.nonce
    const newNonce = direction === 'increase' ? currentNonce + 1 : Math.max(0, currentNonce - 1)

    return {
      type: RecoveryType.ADJUST_NONCE,
      title: `${direction === 'increase' ? 'Increase' : 'Decrease'} Nonce`,
      description: `Nonce issue detected. ${direction === 'increase' ? 'Increasing' : 'Decreasing'} nonce to ${newNonce}.`,
      confidence: 75,
      estimatedCost: '0',
      estimatedTime: 30,
      riskLevel: 'medium',
      parameters: {
        nonce: newNonce
      },
      warnings: ['Nonce adjustment may affect other pending transactions']
    }
  }

  /**
   * Create slippage adjustment strategy
   */
  private createSlippageAdjustmentStrategy(transactionData: any): RecoveryStrategy {
    const currentSlippage = transactionData.metadata?.swapDetails?.slippage || 0.5
    const newSlippage = Math.min(currentSlippage * 2, 5) // Double slippage, max 5%

    return {
      type: RecoveryType.ADJUST_SLIPPAGE,
      title: 'Increase Slippage Tolerance',
      description: `Slippage exceeded. Increasing tolerance from ${currentSlippage}% to ${newSlippage}%.`,
      confidence: 80,
      estimatedCost: '0',
      estimatedTime: 30,
      riskLevel: 'medium',
      parameters: {
        slippage: newSlippage
      },
      warnings: ['Higher slippage may result in less favorable exchange rate']
    }
  }

  /**
   * Create deadline extension strategy
   */
  private createDeadlineExtensionStrategy(transactionData: any): RecoveryStrategy {
    const newDeadline = Date.now() + 600000 // 10 minutes from now

    return {
      type: RecoveryType.EXTEND_DEADLINE,
      title: 'Extend Deadline',
      description: 'Transaction deadline exceeded. Extending deadline by 10 minutes.',
      confidence: 85,
      estimatedCost: '0',
      estimatedTime: 30,
      riskLevel: 'low',
      parameters: {
        deadline: newDeadline
      }
    }
  }

  /**
   * Create split transaction strategy
   */
  private createSplitTransactionStrategy(transactionData: any): RecoveryStrategy {
    return {
      type: RecoveryType.SPLIT_TRANSACTION,
      title: 'Split Transaction',
      description: 'Insufficient funds for full amount. Split into smaller transactions.',
      confidence: 70,
      estimatedCost: this.calculateSplitCost(transactionData),
      estimatedTime: 120,
      riskLevel: 'medium',
      parameters: {
        value: (BigInt(transactionData.value || '0') / BigInt(2)).toString()
      },
      warnings: ['Will require multiple transactions to complete']
    }
  }

  /**
   * Create retry strategy
   */
  private createRetryStrategy(transactionData: any): RecoveryStrategy {
    return {
      type: RecoveryType.RETRY_TRANSACTION,
      title: 'Retry Transaction',
      description: 'Network error detected. Retry transaction with same parameters.',
      confidence: 60,
      estimatedCost: '0',
      estimatedTime: 60,
      riskLevel: 'low',
      parameters: {}
    }
  }

  /**
   * Create contract analysis strategy
   */
  private createContractAnalysisStrategy(transactionData: any, failureDetails: FailureDetails): RecoveryStrategy {
    return {
      type: RecoveryType.MANUAL_INTERVENTION,
      title: 'Contract Analysis Required',
      description: `Contract execution reverted: ${failureDetails.revertReason || 'Unknown reason'}`,
      confidence: 30,
      estimatedCost: '0',
      estimatedTime: 0,
      riskLevel: 'high',
      parameters: {},
      prerequisites: ['Manual analysis of contract state and parameters required'],
      warnings: ['May require contract interaction changes']
    }
  }

  /**
   * Create generic retry strategy
   */
  private createGenericRetryStrategy(transactionData: any): RecoveryStrategy {
    return {
      type: RecoveryType.RETRY_TRANSACTION,
      title: 'Generic Retry',
      description: 'Unknown failure reason. Retry with slightly higher gas price.',
      confidence: 40,
      estimatedCost: this.calculateGasCostIncrease(transactionData),
      estimatedTime: 60,
      riskLevel: 'medium',
      parameters: {
        gasPrice: (BigInt(transactionData.gasPrice || '0') * BigInt(110) / BigInt(100)).toString()
      }
    }
  }

  /**
   * Get manual intervention strategy
   */
  private getManualInterventionStrategy(): RecoveryStrategy {
    return {
      type: RecoveryType.MANUAL_INTERVENTION,
      title: 'Manual Intervention Required',
      description: 'Automatic recovery not possible. Manual analysis required.',
      confidence: 0,
      estimatedCost: '0',
      estimatedTime: 0,
      riskLevel: 'high',
      parameters: {},
      prerequisites: ['Manual analysis and intervention required']
    }
  }

  /**
   * Calculate gas cost increase
   */
  private calculateGasCostIncrease(
    transactionData: any,
    newGasLimit?: string,
    newGasPrice?: string
  ): string {
    const gasLimit = BigInt(newGasLimit || transactionData.gasLimit || transactionData.gas)
    const gasPrice = BigInt(newGasPrice || transactionData.gasPrice || '0')
    const currentCost = BigInt(transactionData.gasLimit || transactionData.gas) * BigInt(transactionData.gasPrice || '0')
    const newCost = gasLimit * gasPrice
    const increase = newCost - currentCost
    
    return (Number(increase) / 1e18).toFixed(6) // Convert to ETH
  }

  /**
   * Calculate speed up cost
   */
  private calculateSpeedUpCost(transactionData: any): string {
    const increasePercentage = this.config.gasIncreasePercentage
    const currentCost = BigInt(transactionData.gasLimit || transactionData.gas) * BigInt(transactionData.gasPrice || '0')
    const increase = currentCost * BigInt(increasePercentage) / BigInt(100)
    
    return (Number(increase) / 1e18).toFixed(6)
  }

  /**
   * Calculate replacement cost
   */
  private calculateReplacementCost(transactionData: any): string {
    return this.calculateSpeedUpCost(transactionData)
  }

  /**
   * Calculate split cost
   */
  private calculateSplitCost(transactionData: any): string {
    const baseCost = BigInt(transactionData.gasLimit || transactionData.gas) * BigInt(transactionData.gasPrice || '0')
    const splitCost = baseCost * BigInt(2) // Assume 2 transactions
    
    return (Number(splitCost) / 1e18).toFixed(6)
  }

  /**
   * Queue transaction for recovery
   */
  queueForRecovery(hash: Hash): void {
    this.recoveryQueue.add(hash)
    
    if (!this.isProcessing) {
      this.startRecoveryProcessing()
    }
  }

  /**
   * Start recovery processing
   */
  private startRecoveryProcessing(): void {
    if (this.isProcessing) return

    this.isProcessing = true
    this.processingInterval = setInterval(() => {
      this.processRecoveryQueue()
    }, this.config.retryDelay)
  }

  /**
   * Process recovery queue
   */
  private async processRecoveryQueue(): Promise<void> {
    if (this.recoveryQueue.size === 0) {
      this.stopRecoveryProcessing()
      return
    }

    for (const hash of Array.from(this.recoveryQueue)) {
      const failedTx = this.failedTransactions.get(hash)
      if (!failedTx) {
        this.recoveryQueue.delete(hash)
        continue
      }

      if (failedTx.recoveryAttempts.length >= this.config.maxRecoveryAttempts) {
        failedTx.recoveryStatus = RecoveryStatus.RECOVERY_FAILED
        this.recoveryQueue.delete(hash)
        continue
      }

      await this.attemptRecovery(failedTx)
    }
  }

  /**
   * Attempt recovery for failed transaction
   */
  private async attemptRecovery(failedTx: FailedTransaction): Promise<void> {
    const strategy = failedTx.suggestedFix
    const attemptId = `recovery_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

    const attempt: RecoveryAttempt = {
      id: attemptId,
      timestamp: Date.now(),
      strategy,
      status: 'pending'
    }

    failedTx.recoveryAttempts.push(attempt)
    failedTx.recoveryStatus = RecoveryStatus.RECOVERY_IN_PROGRESS

    try {
      // Emit recovery started event
      this.emitEvent({
        type: 'recovery_started',
        transaction: failedTx,
        attempt,
        timestamp: Date.now()
      })

      // Execute recovery strategy (mock implementation)
      const result = await this.executeRecoveryStrategy(failedTx, strategy)
      
      attempt.status = 'success'
      attempt.newHash = result.hash
      attempt.cost = result.cost
      
      failedTx.recoveryStatus = RecoveryStatus.RECOVERY_SUCCESS
      this.recoveryQueue.delete(failedTx.hash)

      // Emit recovery success event
      this.emitEvent({
        type: 'recovery_success',
        transaction: failedTx,
        attempt,
        timestamp: Date.now()
      })

    } catch (error) {
      attempt.status = 'failed'
      attempt.error = (error as Error).message

      // Emit recovery failed event
      this.emitEvent({
        type: 'recovery_failed',
        transaction: failedTx,
        attempt,
        timestamp: Date.now()
      })

      // Try next strategy if available
      const nextStrategy = this.getNextRecoveryStrategy(failedTx)
      if (nextStrategy) {
        failedTx.suggestedFix = nextStrategy
        failedTx.recoveryStatus = RecoveryStatus.RECOVERY_AVAILABLE
      } else {
        failedTx.recoveryStatus = RecoveryStatus.RECOVERY_FAILED
        this.recoveryQueue.delete(failedTx.hash)
      }
    }
  }

  /**
   * Execute recovery strategy (mock implementation)
   */
  private async executeRecoveryStrategy(
    failedTx: FailedTransaction,
    strategy: RecoveryStrategy
  ): Promise<{ hash: Hash; cost: string }> {
    // Simulate recovery execution
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Simulate 80% success rate
    if (Math.random() < 0.8) {
      return {
        hash: `0x${Math.random().toString(16).substr(2, 64)}` as Hash,
        cost: strategy.estimatedCost
      }
    } else {
      throw new Error('Recovery strategy failed')
    }
  }

  /**
   * Get next recovery strategy
   */
  private getNextRecoveryStrategy(failedTx: FailedTransaction): RecoveryStrategy | null {
    // In a real implementation, this would select the next best strategy
    // based on the failure reason and previous attempts
    return null
  }

  /**
   * Stop recovery processing
   */
  private stopRecoveryProcessing(): void {
    if (!this.isProcessing) return

    this.isProcessing = false
    if (this.processingInterval) {
      clearInterval(this.processingInterval)
      this.processingInterval = undefined
    }
  }

  /**
   * Get failed transaction
   */
  getFailedTransaction(hash: Hash): FailedTransaction | null {
    return this.failedTransactions.get(hash) || null
  }

  /**
   * Get all failed transactions
   */
  getAllFailedTransactions(): FailedTransaction[] {
    return Array.from(this.failedTransactions.values())
  }

  /**
   * Get recoverable transactions
   */
  getRecoverableTransactions(): FailedTransaction[] {
    return this.getAllFailedTransactions().filter(tx => tx.canRecover)
  }

  /**
   * Update configuration
   */
  updateConfig(config: Partial<RecoveryConfig>): void {
    this.config = { ...this.config, ...config }
  }

  /**
   * Get configuration
   */
  getConfig(): RecoveryConfig {
    return { ...this.config }
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: RecoveryEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in recovery event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: RecoveryEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear failed transactions
   */
  clearFailedTransactions(): void {
    this.failedTransactions.clear()
    this.recoveryQueue.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.stopRecoveryProcessing()
    this.failedTransactions.clear()
    this.recoveryQueue.clear()
    this.eventListeners.clear()
  }
}

export interface RecoveryEvent {
  type: 'analysis_complete' | 'recovery_started' | 'recovery_success' | 'recovery_failed'
  transaction: FailedTransaction
  attempt?: RecoveryAttempt
  timestamp: number
}

// Export singleton instance
export const transactionRecovery = TransactionRecoveryEngine.getInstance()
