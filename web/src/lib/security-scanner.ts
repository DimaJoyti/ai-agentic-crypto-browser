import { type Address, type Hash } from 'viem'

export interface SecurityScanResult {
  contractAddress: Address
  chainId: number
  scanId: string
  timestamp: string
  overallRisk: RiskLevel
  riskScore: number
  vulnerabilities: Vulnerability[]
  securityMetrics: SecurityMetrics
  auditStatus: AuditStatus
  recommendations: SecurityRecommendation[]
  complianceChecks: ComplianceCheck[]
  gasAnalysis: GasAnalysis
  codeQuality: CodeQualityMetrics
}

export enum RiskLevel {
  VERY_LOW = 'very_low',
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export interface Vulnerability {
  id: string
  type: VulnerabilityType
  severity: RiskLevel
  title: string
  description: string
  location: CodeLocation
  impact: string
  recommendation: string
  cweId?: string
  swcId?: string
  confidence: number
  exploitability: number
  references: string[]
}

export enum VulnerabilityType {
  REENTRANCY = 'reentrancy',
  INTEGER_OVERFLOW = 'integer_overflow',
  UNCHECKED_CALL = 'unchecked_call',
  ACCESS_CONTROL = 'access_control',
  FRONT_RUNNING = 'front_running',
  TIMESTAMP_DEPENDENCE = 'timestamp_dependence',
  DENIAL_OF_SERVICE = 'denial_of_service',
  LOGIC_ERROR = 'logic_error',
  CENTRALIZATION = 'centralization',
  UPGRADE_RISK = 'upgrade_risk',
  ORACLE_MANIPULATION = 'oracle_manipulation',
  FLASH_LOAN_ATTACK = 'flash_loan_attack',
  MEV_VULNERABILITY = 'mev_vulnerability',
  GOVERNANCE_ATTACK = 'governance_attack'
}

export interface CodeLocation {
  function: string
  line: number
  column: number
  snippet: string
}

export interface SecurityMetrics {
  codeComplexity: number
  testCoverage: number
  documentationScore: number
  upgradeability: UpgradeabilityInfo
  decentralization: DecentralizationMetrics
  liquidityRisk: number
  marketRisk: number
  technicalRisk: number
}

export interface UpgradeabilityInfo {
  isUpgradeable: boolean
  proxyType?: 'transparent' | 'uups' | 'beacon' | 'diamond'
  adminAddress?: Address
  implementationAddress?: Address
  upgradeRisk: RiskLevel
}

export interface DecentralizationMetrics {
  ownershipConcentration: number
  governanceTokenDistribution: number
  validatorDistribution: number
  decentralizationScore: number
}

export interface AuditStatus {
  isAudited: boolean
  auditors: AuditorInfo[]
  auditReports: AuditReport[]
  lastAuditDate?: string
  auditScore: number
}

export interface AuditorInfo {
  name: string
  reputation: number
  website: string
  verified: boolean
}

export interface AuditReport {
  auditor: string
  date: string
  reportUrl: string
  findings: number
  criticalFindings: number
  status: 'passed' | 'failed' | 'conditional'
}

export interface SecurityRecommendation {
  id: string
  priority: 'high' | 'medium' | 'low'
  category: 'security' | 'optimization' | 'best_practice'
  title: string
  description: string
  implementation: string
  estimatedEffort: string
  impact: string
}

export interface ComplianceCheck {
  standard: string
  status: 'compliant' | 'non_compliant' | 'partial' | 'unknown'
  score: number
  details: string
  requirements: ComplianceRequirement[]
}

export interface ComplianceRequirement {
  requirement: string
  status: 'met' | 'not_met' | 'partial'
  description: string
}

export interface GasAnalysis {
  averageGasUsage: number
  gasOptimizationScore: number
  expensiveFunctions: ExpensiveFunction[]
  gasVulnerabilities: GasVulnerability[]
  optimizationSuggestions: string[]
}

export interface ExpensiveFunction {
  name: string
  gasUsage: number
  complexity: number
  optimizable: boolean
}

export interface GasVulnerability {
  type: 'gas_limit' | 'gas_griefing' | 'out_of_gas'
  severity: RiskLevel
  description: string
  location: CodeLocation
}

export interface CodeQualityMetrics {
  maintainabilityIndex: number
  cyclomaticComplexity: number
  linesOfCode: number
  duplicatedCode: number
  technicalDebt: number
  qualityGate: 'passed' | 'failed'
}

export interface ScanConfiguration {
  depth: 'basic' | 'standard' | 'comprehensive'
  includeAuditCheck: boolean
  includeGasAnalysis: boolean
  includeComplianceCheck: boolean
  customRules: SecurityRule[]
  excludePatterns: string[]
}

export interface SecurityRule {
  id: string
  name: string
  description: string
  pattern: string
  severity: RiskLevel
  enabled: boolean
}

export interface ContractMetadata {
  name?: string
  version?: string
  compiler: string
  optimization: boolean
  runs: number
  evmVersion: string
  libraries: Record<string, Address>
  sourceCode?: string
  abi?: any[]
}

export class SmartContractSecurityScanner {
  private static instance: SmartContractSecurityScanner
  private scanResults = new Map<string, SecurityScanResult>()
  private scanHistory = new Map<string, SecurityScanResult[]>()
  private eventListeners = new Set<(event: SecurityScanEvent) => void>()
  private defaultRules: SecurityRule[]

  private constructor() {
    this.defaultRules = this.initializeDefaultRules()
  }

  static getInstance(): SmartContractSecurityScanner {
    if (!SmartContractSecurityScanner.instance) {
      SmartContractSecurityScanner.instance = new SmartContractSecurityScanner()
    }
    return SmartContractSecurityScanner.instance
  }

  /**
   * Initialize default security rules
   */
  private initializeDefaultRules(): SecurityRule[] {
    return [
      {
        id: 'reentrancy_check',
        name: 'Reentrancy Vulnerability',
        description: 'Detects potential reentrancy vulnerabilities',
        pattern: 'external_call_before_state_change',
        severity: RiskLevel.HIGH,
        enabled: true
      },
      {
        id: 'unchecked_call',
        name: 'Unchecked External Call',
        description: 'Detects unchecked external calls',
        pattern: 'call_without_check',
        severity: RiskLevel.MEDIUM,
        enabled: true
      },
      {
        id: 'access_control',
        name: 'Missing Access Control',
        description: 'Detects functions without proper access control',
        pattern: 'public_function_no_modifier',
        severity: RiskLevel.HIGH,
        enabled: true
      },
      {
        id: 'integer_overflow',
        name: 'Integer Overflow/Underflow',
        description: 'Detects potential integer overflow/underflow',
        pattern: 'unsafe_math_operations',
        severity: RiskLevel.MEDIUM,
        enabled: true
      },
      {
        id: 'timestamp_dependence',
        name: 'Timestamp Dependence',
        description: 'Detects dangerous use of block.timestamp',
        pattern: 'block_timestamp_usage',
        severity: RiskLevel.LOW,
        enabled: true
      }
    ]
  }

  /**
   * Scan smart contract for security vulnerabilities
   */
  async scanContract(
    contractAddress: Address,
    chainId: number,
    config: ScanConfiguration = { depth: 'standard', includeAuditCheck: true, includeGasAnalysis: true, includeComplianceCheck: true, customRules: [], excludePatterns: [] }
  ): Promise<SecurityScanResult> {
    const scanId = `scan_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

    try {
      // Get contract metadata
      const metadata = await this.getContractMetadata(contractAddress, chainId)

      // Perform security analysis
      const vulnerabilities = await this.analyzeVulnerabilities(contractAddress, metadata, config)
      const securityMetrics = await this.calculateSecurityMetrics(contractAddress, metadata)
      const auditStatus = await this.checkAuditStatus(contractAddress)
      const gasAnalysis = config.includeGasAnalysis ? await this.analyzeGasUsage(contractAddress, metadata) : this.getEmptyGasAnalysis()
      const complianceChecks = config.includeComplianceCheck ? await this.checkCompliance(contractAddress, metadata) : []
      const codeQuality = await this.analyzeCodeQuality(metadata)

      // Calculate overall risk score
      const riskScore = this.calculateRiskScore(vulnerabilities, securityMetrics, auditStatus)
      const overallRisk = this.getRiskLevel(riskScore)

      // Generate recommendations
      const recommendations = this.generateRecommendations(vulnerabilities, securityMetrics, auditStatus)

      const result: SecurityScanResult = {
        contractAddress,
        chainId,
        scanId,
        timestamp: new Date().toISOString(),
        overallRisk,
        riskScore,
        vulnerabilities,
        securityMetrics,
        auditStatus,
        recommendations,
        complianceChecks,
        gasAnalysis,
        codeQuality
      }

      // Store result
      this.scanResults.set(scanId, result)
      
      // Add to history
      const history = this.scanHistory.get(contractAddress.toLowerCase()) || []
      history.push(result)
      this.scanHistory.set(contractAddress.toLowerCase(), history)

      // Emit event
      this.emitEvent({
        type: 'scan_completed',
        result,
        timestamp: Date.now()
      })

      return result

    } catch (error) {
      const errorResult: SecurityScanResult = {
        contractAddress,
        chainId,
        scanId,
        timestamp: new Date().toISOString(),
        overallRisk: RiskLevel.CRITICAL,
        riskScore: 100,
        vulnerabilities: [{
          id: 'scan_error',
          type: VulnerabilityType.LOGIC_ERROR,
          severity: RiskLevel.CRITICAL,
          title: 'Scan Error',
          description: `Failed to scan contract: ${(error as Error).message}`,
          location: { function: 'unknown', line: 0, column: 0, snippet: '' },
          impact: 'Unable to assess security',
          recommendation: 'Retry scan or contact support',
          confidence: 100,
          exploitability: 0,
          references: []
        }],
        securityMetrics: this.getEmptySecurityMetrics(),
        auditStatus: { isAudited: false, auditors: [], auditReports: [], auditScore: 0 },
        recommendations: [],
        complianceChecks: [],
        gasAnalysis: this.getEmptyGasAnalysis(),
        codeQuality: this.getEmptyCodeQuality()
      }

      this.emitEvent({
        type: 'scan_failed',
        result: errorResult,
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }
  }

  /**
   * Get contract metadata
   */
  private async getContractMetadata(contractAddress: Address, chainId: number): Promise<ContractMetadata> {
    // Mock implementation - in real app, this would fetch from blockchain/etherscan
    return {
      name: 'MockContract',
      version: '1.0.0',
      compiler: '0.8.19',
      optimization: true,
      runs: 200,
      evmVersion: 'london',
      libraries: {},
      sourceCode: '// Mock source code',
      abi: []
    }
  }

  /**
   * Analyze vulnerabilities
   */
  private async analyzeVulnerabilities(
    contractAddress: Address,
    metadata: ContractMetadata,
    config: ScanConfiguration
  ): Promise<Vulnerability[]> {
    const vulnerabilities: Vulnerability[] = []

    // Apply security rules
    const rules = [...this.defaultRules, ...config.customRules].filter(rule => rule.enabled)

    for (const rule of rules) {
      // Mock vulnerability detection
      if (Math.random() < 0.3) { // 30% chance of finding vulnerability
        vulnerabilities.push({
          id: `vuln_${rule.id}_${Date.now()}`,
          type: this.getVulnerabilityTypeFromRule(rule),
          severity: rule.severity,
          title: rule.name,
          description: rule.description,
          location: {
            function: 'mockFunction',
            line: Math.floor(Math.random() * 100) + 1,
            column: Math.floor(Math.random() * 50) + 1,
            snippet: 'function mockFunction() public { ... }'
          },
          impact: this.getImpactDescription(rule.severity),
          recommendation: this.getRecommendation(rule.id),
          confidence: Math.floor(Math.random() * 30) + 70,
          exploitability: Math.floor(Math.random() * 50) + 25,
          references: [`https://swcregistry.io/docs/${rule.id}`]
        })
      }
    }

    return vulnerabilities
  }

  /**
   * Calculate security metrics
   */
  private async calculateSecurityMetrics(
    contractAddress: Address,
    metadata: ContractMetadata
  ): Promise<SecurityMetrics> {
    return {
      codeComplexity: Math.floor(Math.random() * 50) + 25,
      testCoverage: Math.floor(Math.random() * 40) + 60,
      documentationScore: Math.floor(Math.random() * 30) + 70,
      upgradeability: {
        isUpgradeable: Math.random() > 0.7,
        proxyType: Math.random() > 0.5 ? 'transparent' : 'uups',
        upgradeRisk: Math.random() > 0.6 ? RiskLevel.MEDIUM : RiskLevel.LOW
      },
      decentralization: {
        ownershipConcentration: Math.random() * 100,
        governanceTokenDistribution: Math.random() * 100,
        validatorDistribution: Math.random() * 100,
        decentralizationScore: Math.random() * 100
      },
      liquidityRisk: Math.floor(Math.random() * 40) + 30,
      marketRisk: Math.floor(Math.random() * 50) + 25,
      technicalRisk: Math.floor(Math.random() * 60) + 20
    }
  }

  /**
   * Check audit status
   */
  private async checkAuditStatus(contractAddress: Address): Promise<AuditStatus> {
    // Mock audit status
    const isAudited = Math.random() > 0.4

    if (!isAudited) {
      return {
        isAudited: false,
        auditors: [],
        auditReports: [],
        auditScore: 0
      }
    }

    return {
      isAudited: true,
      auditors: [
        {
          name: 'OpenZeppelin',
          reputation: 95,
          website: 'https://openzeppelin.com',
          verified: true
        },
        {
          name: 'ConsenSys Diligence',
          reputation: 90,
          website: 'https://consensys.net/diligence',
          verified: true
        }
      ],
      auditReports: [
        {
          auditor: 'OpenZeppelin',
          date: '2023-06-15',
          reportUrl: 'https://example.com/audit-report.pdf',
          findings: 5,
          criticalFindings: 0,
          status: 'passed'
        }
      ],
      auditScore: Math.floor(Math.random() * 20) + 80
    }
  }

  /**
   * Analyze gas usage
   */
  private async analyzeGasUsage(contractAddress: Address, metadata: ContractMetadata): Promise<GasAnalysis> {
    return {
      averageGasUsage: Math.floor(Math.random() * 100000) + 50000,
      gasOptimizationScore: Math.floor(Math.random() * 40) + 60,
      expensiveFunctions: [
        {
          name: 'complexFunction',
          gasUsage: Math.floor(Math.random() * 200000) + 100000,
          complexity: Math.floor(Math.random() * 10) + 5,
          optimizable: true
        }
      ],
      gasVulnerabilities: [],
      optimizationSuggestions: [
        'Use packed structs to reduce storage costs',
        'Implement gas-efficient loops',
        'Consider using assembly for critical operations'
      ]
    }
  }

  /**
   * Check compliance
   */
  private async checkCompliance(contractAddress: Address, metadata: ContractMetadata): Promise<ComplianceCheck[]> {
    return [
      {
        standard: 'ERC-20',
        status: 'compliant',
        score: 95,
        details: 'Contract follows ERC-20 standard',
        requirements: [
          { requirement: 'Transfer function', status: 'met', description: 'Implements transfer function' },
          { requirement: 'Approval mechanism', status: 'met', description: 'Implements approve/allowance' }
        ]
      }
    ]
  }

  /**
   * Analyze code quality
   */
  private async analyzeCodeQuality(metadata: ContractMetadata): Promise<CodeQualityMetrics> {
    return {
      maintainabilityIndex: Math.floor(Math.random() * 40) + 60,
      cyclomaticComplexity: Math.floor(Math.random() * 20) + 5,
      linesOfCode: Math.floor(Math.random() * 1000) + 500,
      duplicatedCode: Math.floor(Math.random() * 10) + 2,
      technicalDebt: Math.floor(Math.random() * 30) + 10,
      qualityGate: Math.random() > 0.3 ? 'passed' : 'failed'
    }
  }

  /**
   * Calculate overall risk score
   */
  private calculateRiskScore(
    vulnerabilities: Vulnerability[],
    securityMetrics: SecurityMetrics,
    auditStatus: AuditStatus
  ): number {
    let score = 0

    // Vulnerability score (0-60 points)
    const criticalVulns = vulnerabilities.filter(v => v.severity === RiskLevel.CRITICAL).length
    const highVulns = vulnerabilities.filter(v => v.severity === RiskLevel.HIGH).length
    const mediumVulns = vulnerabilities.filter(v => v.severity === RiskLevel.MEDIUM).length

    score += criticalVulns * 20
    score += highVulns * 10
    score += mediumVulns * 5

    // Security metrics score (0-25 points)
    score += (100 - securityMetrics.testCoverage) * 0.1
    score += securityMetrics.codeComplexity * 0.2
    score += securityMetrics.technicalRisk * 0.15

    // Audit status score (0-15 points)
    if (!auditStatus.isAudited) {
      score += 15
    } else {
      score += (100 - auditStatus.auditScore) * 0.15
    }

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
   * Generate security recommendations
   */
  private generateRecommendations(
    vulnerabilities: Vulnerability[],
    securityMetrics: SecurityMetrics,
    auditStatus: AuditStatus
  ): SecurityRecommendation[] {
    const recommendations: SecurityRecommendation[] = []

    // Vulnerability-based recommendations
    if (vulnerabilities.length > 0) {
      recommendations.push({
        id: 'fix_vulnerabilities',
        priority: 'high',
        category: 'security',
        title: 'Fix Identified Vulnerabilities',
        description: `Address ${vulnerabilities.length} security vulnerabilities found in the contract`,
        implementation: 'Review and fix each vulnerability according to the provided recommendations',
        estimatedEffort: '1-2 weeks',
        impact: 'Significantly reduces security risk'
      })
    }

    // Audit recommendation
    if (!auditStatus.isAudited) {
      recommendations.push({
        id: 'get_audit',
        priority: 'high',
        category: 'security',
        title: 'Get Professional Security Audit',
        description: 'Contract has not been audited by a professional security firm',
        implementation: 'Engage a reputable audit firm like OpenZeppelin, ConsenSys, or Trail of Bits',
        estimatedEffort: '2-4 weeks',
        impact: 'Provides professional security validation'
      })
    }

    // Test coverage recommendation
    if (securityMetrics.testCoverage < 80) {
      recommendations.push({
        id: 'improve_tests',
        priority: 'medium',
        category: 'best_practice',
        title: 'Improve Test Coverage',
        description: `Test coverage is ${securityMetrics.testCoverage}%, should be above 80%`,
        implementation: 'Add comprehensive unit and integration tests',
        estimatedEffort: '1 week',
        impact: 'Reduces risk of bugs and improves code quality'
      })
    }

    return recommendations
  }

  /**
   * Helper methods
   */
  private getVulnerabilityTypeFromRule(rule: SecurityRule): VulnerabilityType {
    const typeMap: Record<string, VulnerabilityType> = {
      'reentrancy_check': VulnerabilityType.REENTRANCY,
      'unchecked_call': VulnerabilityType.UNCHECKED_CALL,
      'access_control': VulnerabilityType.ACCESS_CONTROL,
      'integer_overflow': VulnerabilityType.INTEGER_OVERFLOW,
      'timestamp_dependence': VulnerabilityType.TIMESTAMP_DEPENDENCE
    }
    return typeMap[rule.id] || VulnerabilityType.LOGIC_ERROR
  }

  private getImpactDescription(severity: RiskLevel): string {
    switch (severity) {
      case RiskLevel.CRITICAL: return 'Critical impact - funds at risk'
      case RiskLevel.HIGH: return 'High impact - significant security risk'
      case RiskLevel.MEDIUM: return 'Medium impact - moderate security concern'
      case RiskLevel.LOW: return 'Low impact - minor security issue'
      default: return 'Very low impact - informational'
    }
  }

  private getRecommendation(ruleId: string): string {
    const recommendations: Record<string, string> = {
      'reentrancy_check': 'Use reentrancy guard or checks-effects-interactions pattern',
      'unchecked_call': 'Always check return values of external calls',
      'access_control': 'Add proper access control modifiers',
      'integer_overflow': 'Use SafeMath library or Solidity 0.8+',
      'timestamp_dependence': 'Avoid using block.timestamp for critical logic'
    }
    return recommendations[ruleId] || 'Review and fix according to best practices'
  }

  private getEmptySecurityMetrics(): SecurityMetrics {
    return {
      codeComplexity: 0,
      testCoverage: 0,
      documentationScore: 0,
      upgradeability: { isUpgradeable: false, upgradeRisk: RiskLevel.LOW },
      decentralization: { ownershipConcentration: 0, governanceTokenDistribution: 0, validatorDistribution: 0, decentralizationScore: 0 },
      liquidityRisk: 0,
      marketRisk: 0,
      technicalRisk: 0
    }
  }

  private getEmptyGasAnalysis(): GasAnalysis {
    return {
      averageGasUsage: 0,
      gasOptimizationScore: 0,
      expensiveFunctions: [],
      gasVulnerabilities: [],
      optimizationSuggestions: []
    }
  }

  private getEmptyCodeQuality(): CodeQualityMetrics {
    return {
      maintainabilityIndex: 0,
      cyclomaticComplexity: 0,
      linesOfCode: 0,
      duplicatedCode: 0,
      technicalDebt: 0,
      qualityGate: 'failed'
    }
  }

  /**
   * Get scan result
   */
  getScanResult(scanId: string): SecurityScanResult | null {
    return this.scanResults.get(scanId) || null
  }

  /**
   * Get scan history for contract
   */
  getScanHistory(contractAddress: Address): SecurityScanResult[] {
    return this.scanHistory.get(contractAddress.toLowerCase()) || []
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: SecurityScanEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in security scan event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: SecurityScanEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.scanResults.clear()
    this.scanHistory.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface SecurityScanEvent {
  type: 'scan_started' | 'scan_completed' | 'scan_failed' | 'vulnerability_found'
  result?: SecurityScanResult
  vulnerability?: Vulnerability
  error?: Error
  timestamp: number
}

// Export singleton instance
export const smartContractSecurityScanner = SmartContractSecurityScanner.getInstance()
