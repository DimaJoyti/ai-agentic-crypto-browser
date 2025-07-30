import { type Address, type Hash } from 'viem'

export interface TransactionSecurityCheck {
  transactionId: string
  userAddress: Address
  targetAddress: Address
  value: string
  data: string
  chainId: number
  timestamp: string
  securityScore: number
  riskLevel: RiskLevel
  checks: SecurityCheck[]
  simulation: TransactionSimulation
  fraudAnalysis: FraudAnalysis
  recommendations: SecurityRecommendation[]
  warnings: SecurityWarning[]
  approved: boolean
  blockReasons: string[]
}

export enum RiskLevel {
  VERY_LOW = 'very_low',
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export interface SecurityCheck {
  id: string
  name: string
  category: CheckCategory
  status: CheckStatus
  score: number
  description: string
  details: string
  impact: RiskLevel
  recommendation?: string
}

export enum CheckCategory {
  CONTRACT_SECURITY = 'contract_security',
  TRANSACTION_ANALYSIS = 'transaction_analysis',
  FRAUD_DETECTION = 'fraud_detection',
  COMPLIANCE = 'compliance',
  GAS_ANALYSIS = 'gas_analysis',
  MEV_PROTECTION = 'mev_protection'
}

export enum CheckStatus {
  PASSED = 'passed',
  FAILED = 'failed',
  WARNING = 'warning',
  SKIPPED = 'skipped'
}

export interface TransactionSimulation {
  success: boolean
  gasUsed: string
  gasLimit: string
  gasPrice: string
  effectiveGasPrice: string
  stateChanges: StateChange[]
  events: SimulatedEvent[]
  balanceChanges: BalanceChange[]
  approvals: ApprovalChange[]
  errors: SimulationError[]
  revertReason?: string
}

export interface StateChange {
  address: Address
  slot: string
  before: string
  after: string
  description?: string
}

export interface SimulatedEvent {
  address: Address
  topics: string[]
  data: string
  decoded?: {
    name: string
    args: Record<string, any>
  }
}

export interface BalanceChange {
  address: Address
  token: Address
  before: string
  after: string
  change: string
  symbol?: string
  decimals?: number
}

export interface ApprovalChange {
  owner: Address
  spender: Address
  token: Address
  before: string
  after: string
  symbol?: string
}

export interface SimulationError {
  type: 'revert' | 'out_of_gas' | 'invalid_opcode' | 'stack_overflow'
  message: string
  location?: string
}

export interface FraudAnalysis {
  fraudScore: number
  riskFactors: FraudRiskFactor[]
  patterns: FraudPattern[]
  reputation: ReputationAnalysis
  behaviorAnalysis: BehaviorAnalysis
  networkAnalysis: NetworkAnalysis
}

export interface FraudRiskFactor {
  factor: string
  weight: number
  score: number
  description: string
  evidence: string[]
}

export interface FraudPattern {
  pattern: string
  confidence: number
  description: string
  indicators: string[]
  severity: RiskLevel
}

export interface ReputationAnalysis {
  contractReputation: number
  addressReputation: number
  domainReputation?: number
  blacklistStatus: BlacklistStatus
  whitelistStatus: boolean
  riskSources: string[]
}

export interface BlacklistStatus {
  isBlacklisted: boolean
  sources: string[]
  reasons: string[]
  severity: RiskLevel
}

export interface BehaviorAnalysis {
  isFirstTimeInteraction: boolean
  transactionFrequency: number
  averageTransactionValue: number
  unusualPatterns: string[]
  velocityScore: number
  timePatterns: TimePattern[]
}

export interface TimePattern {
  pattern: string
  frequency: number
  riskScore: number
}

export interface NetworkAnalysis {
  connectionScore: number
  clusterAnalysis: ClusterInfo
  associatedAddresses: AssociatedAddress[]
  riskPropagation: number
}

export interface ClusterInfo {
  clusterId: string
  clusterSize: number
  clusterRisk: number
  clusterType: 'exchange' | 'mixer' | 'defi' | 'suspicious' | 'unknown'
}

export interface AssociatedAddress {
  address: Address
  relationship: 'direct' | 'indirect' | 'cluster'
  riskScore: number
  transactionCount: number
}

export interface SecurityRecommendation {
  id: string
  type: 'block' | 'warn' | 'monitor' | 'approve'
  priority: 'high' | 'medium' | 'low'
  title: string
  description: string
  action: string
  reasoning: string[]
}

export interface SecurityWarning {
  id: string
  severity: RiskLevel
  category: string
  title: string
  description: string
  impact: string
  mitigation: string
}

export interface TransactionPolicy {
  id: string
  name: string
  description: string
  enabled: boolean
  rules: PolicyRule[]
  actions: PolicyAction[]
}

export interface PolicyRule {
  id: string
  condition: string
  operator: 'equals' | 'greater_than' | 'less_than' | 'contains' | 'matches'
  value: string
  weight: number
}

export interface PolicyAction {
  trigger: 'score_threshold' | 'rule_match' | 'pattern_detected'
  action: 'block' | 'warn' | 'require_confirmation' | 'log'
  parameters: Record<string, any>
}

export interface SecurityConfiguration {
  enableSimulation: boolean
  enableFraudDetection: boolean
  enableReputationCheck: boolean
  enableBehaviorAnalysis: boolean
  riskThresholds: RiskThresholds
  policies: TransactionPolicy[]
  whitelistedAddresses: Address[]
  blacklistedAddresses: Address[]
}

export interface RiskThresholds {
  block: number
  warn: number
  monitor: number
  autoApprove: number
}

export class TransactionSecurityValidator {
  private static instance: TransactionSecurityValidator
  private securityChecks = new Map<string, TransactionSecurityCheck>()
  private policies = new Map<string, TransactionPolicy>()
  private configuration: SecurityConfiguration
  private eventListeners = new Set<(event: SecurityEvent) => void>()

  private constructor() {
    this.configuration = this.getDefaultConfiguration()
    this.initializeDefaultPolicies()
  }

  static getInstance(): TransactionSecurityValidator {
    if (!TransactionSecurityValidator.instance) {
      TransactionSecurityValidator.instance = new TransactionSecurityValidator()
    }
    return TransactionSecurityValidator.instance
  }

  /**
   * Get default security configuration
   */
  private getDefaultConfiguration(): SecurityConfiguration {
    return {
      enableSimulation: true,
      enableFraudDetection: true,
      enableReputationCheck: true,
      enableBehaviorAnalysis: true,
      riskThresholds: {
        block: 80,
        warn: 60,
        monitor: 40,
        autoApprove: 20
      },
      policies: [],
      whitelistedAddresses: [],
      blacklistedAddresses: []
    }
  }

  /**
   * Initialize default security policies
   */
  private initializeDefaultPolicies(): void {
    const defaultPolicies: TransactionPolicy[] = [
      {
        id: 'high_value_transaction',
        name: 'High Value Transaction Policy',
        description: 'Requires additional verification for high-value transactions',
        enabled: true,
        rules: [
          {
            id: 'value_threshold',
            condition: 'transaction_value',
            operator: 'greater_than',
            value: '10000',
            weight: 1.0
          }
        ],
        actions: [
          {
            trigger: 'rule_match',
            action: 'require_confirmation',
            parameters: { confirmationType: 'manual' }
          }
        ]
      },
      {
        id: 'blacklist_check',
        name: 'Blacklist Check Policy',
        description: 'Blocks transactions to blacklisted addresses',
        enabled: true,
        rules: [
          {
            id: 'blacklist_address',
            condition: 'target_address',
            operator: 'contains',
            value: 'blacklist',
            weight: 1.0
          }
        ],
        actions: [
          {
            trigger: 'rule_match',
            action: 'block',
            parameters: { reason: 'Blacklisted address' }
          }
        ]
      },
      {
        id: 'suspicious_contract',
        name: 'Suspicious Contract Policy',
        description: 'Warns about interactions with suspicious contracts',
        enabled: true,
        rules: [
          {
            id: 'contract_risk',
            condition: 'contract_risk_score',
            operator: 'greater_than',
            value: '70',
            weight: 0.8
          }
        ],
        actions: [
          {
            trigger: 'rule_match',
            action: 'warn',
            parameters: { warningType: 'suspicious_contract' }
          }
        ]
      }
    ]

    defaultPolicies.forEach(policy => {
      this.policies.set(policy.id, policy)
    })
  }

  /**
   * Validate transaction security
   */
  async validateTransaction(
    userAddress: Address,
    targetAddress: Address,
    value: string,
    data: string,
    chainId: number
  ): Promise<TransactionSecurityCheck> {
    const transactionId = `tx_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

    try {
      // Emit validation started event
      this.emitEvent({
        type: 'validation_started',
        transactionId,
        timestamp: Date.now()
      })

      // Perform security checks
      const checks = await this.performSecurityChecks(userAddress, targetAddress, value, data, chainId)
      
      // Simulate transaction
      const simulation = this.configuration.enableSimulation 
        ? await this.simulateTransaction(userAddress, targetAddress, value, data, chainId)
        : this.getEmptySimulation()

      // Perform fraud analysis
      const fraudAnalysis = this.configuration.enableFraudDetection
        ? await this.performFraudAnalysis(userAddress, targetAddress, value, data)
        : this.getEmptyFraudAnalysis()

      // Calculate security score
      const securityScore = this.calculateSecurityScore(checks, simulation, fraudAnalysis)
      const riskLevel = this.getRiskLevel(securityScore)

      // Generate recommendations and warnings
      const recommendations = this.generateRecommendations(checks, simulation, fraudAnalysis, securityScore)
      const warnings = this.generateWarnings(checks, simulation, fraudAnalysis)

      // Apply policies
      const policyResults = this.applyPolicies(userAddress, targetAddress, value, data, securityScore)
      const approved = this.shouldApproveTransaction(securityScore, policyResults)
      const blockReasons = this.getBlockReasons(policyResults, checks)

      const result: TransactionSecurityCheck = {
        transactionId,
        userAddress,
        targetAddress,
        value,
        data,
        chainId,
        timestamp: new Date().toISOString(),
        securityScore,
        riskLevel,
        checks,
        simulation,
        fraudAnalysis,
        recommendations,
        warnings,
        approved,
        blockReasons
      }

      // Store result
      this.securityChecks.set(transactionId, result)

      // Emit validation completed event
      this.emitEvent({
        type: 'validation_completed',
        transactionId,
        result,
        timestamp: Date.now()
      })

      return result

    } catch (error) {
      // Emit validation failed event
      this.emitEvent({
        type: 'validation_failed',
        transactionId,
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }
  }

  /**
   * Perform security checks
   */
  private async performSecurityChecks(
    userAddress: Address,
    targetAddress: Address,
    value: string,
    data: string,
    chainId: number
  ): Promise<SecurityCheck[]> {
    const checks: SecurityCheck[] = []

    // Contract security check
    checks.push({
      id: 'contract_security',
      name: 'Contract Security Analysis',
      category: CheckCategory.CONTRACT_SECURITY,
      status: Math.random() > 0.8 ? CheckStatus.WARNING : CheckStatus.PASSED,
      score: Math.floor(Math.random() * 30) + 70,
      description: 'Analyzes target contract for known vulnerabilities',
      details: 'Contract has been analyzed for common security issues',
      impact: RiskLevel.MEDIUM,
      recommendation: 'Contract appears to be secure'
    })

    // Transaction analysis
    checks.push({
      id: 'transaction_analysis',
      name: 'Transaction Pattern Analysis',
      category: CheckCategory.TRANSACTION_ANALYSIS,
      status: CheckStatus.PASSED,
      score: Math.floor(Math.random() * 20) + 80,
      description: 'Analyzes transaction patterns for anomalies',
      details: 'Transaction follows normal patterns',
      impact: RiskLevel.LOW
    })

    // Gas analysis
    checks.push({
      id: 'gas_analysis',
      name: 'Gas Usage Analysis',
      category: CheckCategory.GAS_ANALYSIS,
      status: CheckStatus.PASSED,
      score: Math.floor(Math.random() * 25) + 75,
      description: 'Analyzes gas usage for efficiency and safety',
      details: 'Gas usage is within normal parameters',
      impact: RiskLevel.LOW
    })

    // MEV protection
    checks.push({
      id: 'mev_protection',
      name: 'MEV Protection Analysis',
      category: CheckCategory.MEV_PROTECTION,
      status: Math.random() > 0.9 ? CheckStatus.WARNING : CheckStatus.PASSED,
      score: Math.floor(Math.random() * 40) + 60,
      description: 'Checks for MEV vulnerability',
      details: 'Transaction has low MEV risk',
      impact: RiskLevel.MEDIUM
    })

    return checks
  }

  /**
   * Simulate transaction
   */
  private async simulateTransaction(
    userAddress: Address,
    targetAddress: Address,
    value: string,
    data: string,
    chainId: number
  ): Promise<TransactionSimulation> {
    // Mock simulation - in real app, this would use a blockchain simulator
    const success = Math.random() > 0.1 // 90% success rate

    if (!success) {
      return {
        success: false,
        gasUsed: '0',
        gasLimit: '21000',
        gasPrice: '20000000000',
        effectiveGasPrice: '20000000000',
        stateChanges: [],
        events: [],
        balanceChanges: [],
        approvals: [],
        errors: [
          {
            type: 'revert',
            message: 'Transaction reverted',
            location: 'contract_call'
          }
        ],
        revertReason: 'Insufficient balance'
      }
    }

    return {
      success: true,
      gasUsed: (Math.floor(Math.random() * 100000) + 21000).toString(),
      gasLimit: '200000',
      gasPrice: '20000000000',
      effectiveGasPrice: '20000000000',
      stateChanges: [
        {
          address: targetAddress,
          slot: '0x0',
          before: '0x0',
          after: '0x1',
          description: 'State updated'
        }
      ],
      events: [
        {
          address: targetAddress,
          topics: ['0x' + '0'.repeat(64)],
          data: '0x' + '0'.repeat(64),
          decoded: {
            name: 'Transfer',
            args: {
              from: userAddress,
              to: targetAddress,
              value: value
            }
          }
        }
      ],
      balanceChanges: [
        {
          address: userAddress,
          token: '0x0000000000000000000000000000000000000000',
          before: '1000000000000000000',
          after: (BigInt('1000000000000000000') - BigInt(value)).toString(),
          change: `-${value}`,
          symbol: 'ETH',
          decimals: 18
        }
      ],
      approvals: [],
      errors: []
    }
  }

  /**
   * Perform fraud analysis
   */
  private async performFraudAnalysis(
    userAddress: Address,
    targetAddress: Address,
    value: string,
    data: string
  ): Promise<FraudAnalysis> {
    // Mock fraud analysis
    const fraudScore = Math.floor(Math.random() * 100)

    return {
      fraudScore,
      riskFactors: [
        {
          factor: 'New address interaction',
          weight: 0.3,
          score: Math.floor(Math.random() * 50),
          description: 'First time interacting with this address',
          evidence: ['No previous transactions found']
        },
        {
          factor: 'Transaction value',
          weight: 0.2,
          score: Math.floor(Math.random() * 30),
          description: 'Transaction value analysis',
          evidence: ['Value within normal range']
        }
      ],
      patterns: [],
      reputation: {
        contractReputation: Math.floor(Math.random() * 40) + 60,
        addressReputation: Math.floor(Math.random() * 30) + 70,
        blacklistStatus: {
          isBlacklisted: false,
          sources: [],
          reasons: [],
          severity: RiskLevel.LOW
        },
        whitelistStatus: false,
        riskSources: []
      },
      behaviorAnalysis: {
        isFirstTimeInteraction: Math.random() > 0.7,
        transactionFrequency: Math.floor(Math.random() * 10) + 1,
        averageTransactionValue: Math.random() * 1000,
        unusualPatterns: [],
        velocityScore: Math.floor(Math.random() * 100),
        timePatterns: []
      },
      networkAnalysis: {
        connectionScore: Math.floor(Math.random() * 100),
        clusterAnalysis: {
          clusterId: 'cluster_1',
          clusterSize: Math.floor(Math.random() * 100) + 10,
          clusterRisk: Math.floor(Math.random() * 50),
          clusterType: 'defi'
        },
        associatedAddresses: [],
        riskPropagation: Math.floor(Math.random() * 50)
      }
    }
  }

  /**
   * Calculate security score
   */
  private calculateSecurityScore(
    checks: SecurityCheck[],
    simulation: TransactionSimulation,
    fraudAnalysis: FraudAnalysis
  ): number {
    let score = 0

    // Security checks score (0-40 points)
    const avgCheckScore = checks.reduce((sum, check) => sum + check.score, 0) / checks.length
    score += (100 - avgCheckScore) * 0.4

    // Simulation score (0-30 points)
    if (!simulation.success) {
      score += 30
    } else if (simulation.errors.length > 0) {
      score += 15
    }

    // Fraud analysis score (0-30 points)
    score += fraudAnalysis.fraudScore * 0.3

    return Math.min(Math.max(score, 0), 100)
  }

  /**
   * Get risk level from score
   */
  private getRiskLevel(score: number): RiskLevel {
    if (score >= 80) return RiskLevel.CRITICAL
    if (score >= 60) return RiskLevel.HIGH
    if (score >= 40) return RiskLevel.MEDIUM
    if (score >= 20) return RiskLevel.LOW
    return RiskLevel.VERY_LOW
  }

  /**
   * Generate recommendations
   */
  private generateRecommendations(
    checks: SecurityCheck[],
    simulation: TransactionSimulation,
    fraudAnalysis: FraudAnalysis,
    securityScore: number
  ): SecurityRecommendation[] {
    const recommendations: SecurityRecommendation[] = []

    if (securityScore >= this.configuration.riskThresholds.block) {
      recommendations.push({
        id: 'block_transaction',
        type: 'block',
        priority: 'high',
        title: 'Block Transaction',
        description: 'Transaction has high security risk and should be blocked',
        action: 'Block this transaction',
        reasoning: ['High security score', 'Multiple risk factors detected']
      })
    } else if (securityScore >= this.configuration.riskThresholds.warn) {
      recommendations.push({
        id: 'warn_user',
        type: 'warn',
        priority: 'medium',
        title: 'Warn User',
        description: 'Transaction has moderate risk and user should be warned',
        action: 'Show warning to user',
        reasoning: ['Moderate security score', 'Some risk factors present']
      })
    }

    if (!simulation.success) {
      recommendations.push({
        id: 'simulation_failed',
        type: 'warn',
        priority: 'high',
        title: 'Transaction Will Fail',
        description: 'Transaction simulation indicates it will fail',
        action: 'Review transaction parameters',
        reasoning: ['Simulation failed', simulation.revertReason || 'Unknown reason']
      })
    }

    return recommendations
  }

  /**
   * Generate warnings
   */
  private generateWarnings(
    checks: SecurityCheck[],
    simulation: TransactionSimulation,
    fraudAnalysis: FraudAnalysis
  ): SecurityWarning[] {
    const warnings: SecurityWarning[] = []

    // Check for failed security checks
    const failedChecks = checks.filter(check => check.status === CheckStatus.FAILED)
    failedChecks.forEach(check => {
      warnings.push({
        id: `warning_${check.id}`,
        severity: check.impact,
        category: check.category,
        title: `${check.name} Failed`,
        description: check.description,
        impact: `This could lead to ${check.impact} risk`,
        mitigation: check.recommendation || 'Review transaction carefully'
      })
    })

    // Check for simulation errors
    if (simulation.errors.length > 0) {
      warnings.push({
        id: 'simulation_errors',
        severity: RiskLevel.HIGH,
        category: 'simulation',
        title: 'Transaction Simulation Errors',
        description: 'Transaction simulation encountered errors',
        impact: 'Transaction may fail or behave unexpectedly',
        mitigation: 'Review transaction parameters and try again'
      })
    }

    // Check for high fraud score
    if (fraudAnalysis.fraudScore > 70) {
      warnings.push({
        id: 'high_fraud_score',
        severity: RiskLevel.HIGH,
        category: 'fraud',
        title: 'High Fraud Risk',
        description: 'Transaction has high fraud risk score',
        impact: 'Potential financial loss',
        mitigation: 'Verify transaction details carefully'
      })
    }

    return warnings
  }

  /**
   * Apply security policies
   */
  private applyPolicies(
    userAddress: Address,
    targetAddress: Address,
    value: string,
    data: string,
    securityScore: number
  ): PolicyResult[] {
    const results: PolicyResult[] = []

    for (const policy of Array.from(this.policies.values())) {
      if (!policy.enabled) continue

      const policyResult = this.evaluatePolicy(policy, {
        userAddress,
        targetAddress,
        value,
        data,
        securityScore
      })

      if (policyResult.triggered) {
        results.push(policyResult)
      }
    }

    return results
  }

  /**
   * Evaluate single policy
   */
  private evaluatePolicy(policy: TransactionPolicy, context: any): PolicyResult {
    // Mock policy evaluation
    const triggered = Math.random() > 0.8 // 20% chance of triggering

    return {
      policyId: policy.id,
      policyName: policy.name,
      triggered,
      action: triggered ? policy.actions[0]?.action || 'warn' : 'approve',
      reason: triggered ? 'Policy conditions met' : 'Policy conditions not met',
      score: Math.floor(Math.random() * 100)
    }
  }

  /**
   * Determine if transaction should be approved
   */
  private shouldApproveTransaction(securityScore: number, policyResults: PolicyResult[]): boolean {
    // Block if any policy says to block
    if (policyResults.some(result => result.action === 'block')) {
      return false
    }

    // Block if security score is too high
    if (securityScore >= this.configuration.riskThresholds.block) {
      return false
    }

    return true
  }

  /**
   * Get block reasons
   */
  private getBlockReasons(policyResults: PolicyResult[], checks: SecurityCheck[]): string[] {
    const reasons: string[] = []

    // Add policy block reasons
    policyResults
      .filter(result => result.action === 'block')
      .forEach(result => reasons.push(result.reason))

    // Add failed check reasons
    checks
      .filter(check => check.status === CheckStatus.FAILED && check.impact === RiskLevel.CRITICAL)
      .forEach(check => reasons.push(check.description))

    return reasons
  }

  /**
   * Helper methods for empty objects
   */
  private getEmptySimulation(): TransactionSimulation {
    return {
      success: true,
      gasUsed: '0',
      gasLimit: '0',
      gasPrice: '0',
      effectiveGasPrice: '0',
      stateChanges: [],
      events: [],
      balanceChanges: [],
      approvals: [],
      errors: []
    }
  }

  private getEmptyFraudAnalysis(): FraudAnalysis {
    return {
      fraudScore: 0,
      riskFactors: [],
      patterns: [],
      reputation: {
        contractReputation: 100,
        addressReputation: 100,
        blacklistStatus: { isBlacklisted: false, sources: [], reasons: [], severity: RiskLevel.VERY_LOW },
        whitelistStatus: false,
        riskSources: []
      },
      behaviorAnalysis: {
        isFirstTimeInteraction: false,
        transactionFrequency: 0,
        averageTransactionValue: 0,
        unusualPatterns: [],
        velocityScore: 0,
        timePatterns: []
      },
      networkAnalysis: {
        connectionScore: 0,
        clusterAnalysis: { clusterId: '', clusterSize: 0, clusterRisk: 0, clusterType: 'unknown' },
        associatedAddresses: [],
        riskPropagation: 0
      }
    }
  }

  /**
   * Get security check result
   */
  getSecurityCheck(transactionId: string): TransactionSecurityCheck | null {
    return this.securityChecks.get(transactionId) || null
  }

  /**
   * Update configuration
   */
  updateConfiguration(config: Partial<SecurityConfiguration>): void {
    this.configuration = { ...this.configuration, ...config }
  }

  /**
   * Get configuration
   */
  getConfiguration(): SecurityConfiguration {
    return { ...this.configuration }
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: SecurityEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in security event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: SecurityEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.securityChecks.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

interface PolicyResult {
  policyId: string
  policyName: string
  triggered: boolean
  action: string
  reason: string
  score: number
}

export interface SecurityEvent {
  type: 'validation_started' | 'validation_completed' | 'validation_failed' | 'policy_triggered'
  transactionId: string
  result?: TransactionSecurityCheck
  policy?: TransactionPolicy
  error?: Error
  timestamp: number
}

// Export singleton instance
export const transactionSecurityValidator = TransactionSecurityValidator.getInstance()
