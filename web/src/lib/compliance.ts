import { type Address, type Hash } from 'viem'

export interface ComplianceFramework {
  id: string
  name: string
  description: string
  jurisdiction: string
  version: string
  requirements: ComplianceRequirement[]
  reportingSchedule: ReportingSchedule
  lastUpdate: string
  status: ComplianceStatus
}

export enum ComplianceStatus {
  COMPLIANT = 'compliant',
  NON_COMPLIANT = 'non_compliant',
  PARTIAL = 'partial',
  PENDING = 'pending',
  UNKNOWN = 'unknown'
}

export interface ComplianceRequirement {
  id: string
  name: string
  description: string
  category: RequirementCategory
  mandatory: boolean
  deadline?: string
  status: ComplianceStatus
  evidence: ComplianceEvidence[]
  controls: ComplianceControl[]
  riskLevel: RiskLevel
}

export enum RequirementCategory {
  AML = 'aml',
  KYC = 'kyc',
  SANCTIONS = 'sanctions',
  TAX_REPORTING = 'tax_reporting',
  DATA_PROTECTION = 'data_protection',
  OPERATIONAL = 'operational',
  FINANCIAL = 'financial',
  TECHNICAL = 'technical'
}

export enum RiskLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export interface ComplianceEvidence {
  id: string
  type: EvidenceType
  description: string
  documentUrl?: string
  hash?: string
  timestamp: string
  verifiedBy?: string
  expiryDate?: string
}

export enum EvidenceType {
  DOCUMENT = 'document',
  TRANSACTION_LOG = 'transaction_log',
  AUDIT_REPORT = 'audit_report',
  CERTIFICATION = 'certification',
  POLICY = 'policy',
  PROCEDURE = 'procedure',
  TRAINING_RECORD = 'training_record'
}

export interface ComplianceControl {
  id: string
  name: string
  description: string
  controlType: ControlType
  automated: boolean
  frequency: string
  lastExecution?: string
  nextExecution?: string
  effectiveness: number
  status: ControlStatus
}

export enum ControlType {
  PREVENTIVE = 'preventive',
  DETECTIVE = 'detective',
  CORRECTIVE = 'corrective',
  COMPENSATING = 'compensating'
}

export enum ControlStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  FAILED = 'failed',
  PENDING = 'pending'
}

export interface ReportingSchedule {
  frequency: ReportingFrequency
  nextDue: string
  lastSubmitted?: string
  recipients: string[]
  format: ReportFormat
  automated: boolean
}

export enum ReportingFrequency {
  DAILY = 'daily',
  WEEKLY = 'weekly',
  MONTHLY = 'monthly',
  QUARTERLY = 'quarterly',
  ANNUALLY = 'annually',
  ON_DEMAND = 'on_demand'
}

export enum ReportFormat {
  PDF = 'pdf',
  CSV = 'csv',
  JSON = 'json',
  XML = 'xml',
  XLSX = 'xlsx'
}

export interface AMLCheck {
  id: string
  address: Address
  timestamp: string
  checkType: AMLCheckType
  result: AMLResult
  riskScore: number
  riskLevel: RiskLevel
  sources: AMLSource[]
  findings: AMLFinding[]
  recommendations: string[]
  expiryDate: string
}

export enum AMLCheckType {
  SANCTIONS_SCREENING = 'sanctions_screening',
  PEP_SCREENING = 'pep_screening',
  ADVERSE_MEDIA = 'adverse_media',
  TRANSACTION_MONITORING = 'transaction_monitoring',
  ENHANCED_DUE_DILIGENCE = 'enhanced_due_diligence'
}

export interface AMLResult {
  status: 'clear' | 'alert' | 'blocked'
  confidence: number
  falsePositiveRate: number
  reviewRequired: boolean
  autoApproved: boolean
}

export interface AMLSource {
  name: string
  type: 'sanctions_list' | 'pep_list' | 'adverse_media' | 'internal_database'
  lastUpdate: string
  reliability: number
}

export interface AMLFinding {
  id: string
  type: string
  description: string
  severity: RiskLevel
  source: string
  confidence: number
  details: Record<string, any>
}

export interface KYCProfile {
  id: string
  userAddress: Address
  status: KYCStatus
  tier: KYCTier
  verificationLevel: VerificationLevel
  documents: KYCDocument[]
  checks: KYCCheck[]
  riskAssessment: KYCRiskAssessment
  createdAt: string
  updatedAt: string
  expiryDate?: string
}

export enum KYCStatus {
  PENDING = 'pending',
  VERIFIED = 'verified',
  REJECTED = 'rejected',
  EXPIRED = 'expired',
  SUSPENDED = 'suspended'
}

export enum KYCTier {
  BASIC = 'basic',
  STANDARD = 'standard',
  ENHANCED = 'enhanced',
  INSTITUTIONAL = 'institutional'
}

export enum VerificationLevel {
  LEVEL_1 = 'level_1', // Basic identity verification
  LEVEL_2 = 'level_2', // Enhanced verification with documents
  LEVEL_3 = 'level_3', // Full verification with biometrics
  LEVEL_4 = 'level_4'  // Institutional verification
}

export interface KYCDocument {
  id: string
  type: DocumentType
  status: DocumentStatus
  uploadedAt: string
  verifiedAt?: string
  expiryDate?: string
  hash: string
  metadata: DocumentMetadata
}

export enum DocumentType {
  PASSPORT = 'passport',
  DRIVERS_LICENSE = 'drivers_license',
  NATIONAL_ID = 'national_id',
  PROOF_OF_ADDRESS = 'proof_of_address',
  BANK_STATEMENT = 'bank_statement',
  UTILITY_BILL = 'utility_bill',
  BUSINESS_REGISTRATION = 'business_registration',
  ARTICLES_OF_INCORPORATION = 'articles_of_incorporation'
}

export enum DocumentStatus {
  UPLOADED = 'uploaded',
  PROCESSING = 'processing',
  VERIFIED = 'verified',
  REJECTED = 'rejected',
  EXPIRED = 'expired'
}

export enum CheckStatus {
  PENDING = 'pending',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  FAILED = 'failed',
  EXPIRED = 'expired'
}

export interface DocumentMetadata {
  fileName: string
  fileSize: number
  mimeType: string
  extractedData?: Record<string, any>
  ocrConfidence?: number
  biometricMatch?: number
}

export interface KYCCheck {
  id: string
  type: KYCCheckType
  status: CheckStatus
  result: CheckResult
  timestamp: string
  provider: string
  confidence: number
  details: Record<string, any>
}

export enum KYCCheckType {
  IDENTITY_VERIFICATION = 'identity_verification',
  DOCUMENT_VERIFICATION = 'document_verification',
  BIOMETRIC_VERIFICATION = 'biometric_verification',
  ADDRESS_VERIFICATION = 'address_verification',
  PHONE_VERIFICATION = 'phone_verification',
  EMAIL_VERIFICATION = 'email_verification',
  LIVENESS_CHECK = 'liveness_check'
}

export interface CheckResult {
  passed: boolean
  score: number
  reasons: string[]
  recommendations: string[]
}

export interface KYCRiskAssessment {
  overallRisk: RiskLevel
  riskScore: number
  riskFactors: RiskFactor[]
  mitigatingFactors: string[]
  recommendations: string[]
  lastAssessed: string
  nextReview: string
}

export interface RiskFactor {
  factor: string
  weight: number
  score: number
  description: string
}

export interface TaxReporting {
  id: string
  userAddress: Address
  taxYear: number
  jurisdiction: string
  reportType: TaxReportType
  transactions: TaxableTransaction[]
  calculations: TaxCalculation
  status: ReportStatus
  generatedAt: string
  submittedAt?: string
  filingDeadline: string
}

export enum TaxReportType {
  CAPITAL_GAINS = 'capital_gains',
  INCOME = 'income',
  DEFI_ACTIVITY = 'defi_activity',
  NFT_TRANSACTIONS = 'nft_transactions',
  MINING_REWARDS = 'mining_rewards',
  STAKING_REWARDS = 'staking_rewards'
}

export interface TaxableTransaction {
  id: string
  hash: Hash
  timestamp: string
  type: TransactionType
  asset: string
  amount: string
  costBasis?: string
  fairMarketValue: string
  gainLoss?: string
  taxable: boolean
  category: TaxCategory
}

export enum TransactionType {
  BUY = 'buy',
  SELL = 'sell',
  TRADE = 'trade',
  TRANSFER = 'transfer',
  STAKE = 'stake',
  UNSTAKE = 'unstake',
  REWARD = 'reward',
  AIRDROP = 'airdrop',
  MINING = 'mining',
  DEFI_INTERACTION = 'defi_interaction'
}

export enum TaxCategory {
  CAPITAL_GAINS_SHORT = 'capital_gains_short',
  CAPITAL_GAINS_LONG = 'capital_gains_long',
  ORDINARY_INCOME = 'ordinary_income',
  BUSINESS_INCOME = 'business_income',
  NON_TAXABLE = 'non_taxable'
}

export interface TaxCalculation {
  totalGainLoss: string
  shortTermGains: string
  longTermGains: string
  ordinaryIncome: string
  taxableIncome: string
  estimatedTax: string
  deductions: string
  methodology: string
}

export enum ReportStatus {
  DRAFT = 'draft',
  GENERATED = 'generated',
  REVIEWED = 'reviewed',
  SUBMITTED = 'submitted',
  ACCEPTED = 'accepted',
  REJECTED = 'rejected'
}

export interface ComplianceReport {
  id: string
  type: ReportType
  framework: string
  period: ReportingPeriod
  status: ReportStatus
  data: ComplianceReportData
  generatedAt: string
  submittedAt?: string
  approvedBy?: string
  recipients: string[]
}

export enum ReportType {
  COMPLIANCE_SUMMARY = 'compliance_summary',
  AML_REPORT = 'aml_report',
  KYC_REPORT = 'kyc_report',
  TAX_REPORT = 'tax_report',
  TRANSACTION_REPORT = 'transaction_report',
  RISK_ASSESSMENT = 'risk_assessment',
  AUDIT_REPORT = 'audit_report'
}

export interface ReportingPeriod {
  startDate: string
  endDate: string
  description: string
}

export interface ComplianceReportData {
  summary: ReportSummary
  metrics: ComplianceMetrics
  findings: ComplianceFinding[]
  recommendations: string[]
  attachments: string[]
}

export interface ReportSummary {
  totalRequirements: number
  compliantRequirements: number
  nonCompliantRequirements: number
  pendingRequirements: number
  overallComplianceRate: number
  riskLevel: RiskLevel
}

export interface ComplianceMetrics {
  amlChecks: number
  kycVerifications: number
  transactionsMonitored: number
  alertsGenerated: number
  falsePositives: number
  truePositives: number
  averageProcessingTime: number
}

export interface ComplianceFinding {
  id: string
  type: string
  severity: RiskLevel
  description: string
  requirement: string
  evidence: string[]
  remediation: string
  deadline?: string
  responsible: string
}

export class ComplianceManager {
  private static instance: ComplianceManager
  private frameworks = new Map<string, ComplianceFramework>()
  private amlChecks = new Map<string, AMLCheck>()
  private kycProfiles = new Map<string, KYCProfile>()
  private taxReports = new Map<string, TaxReporting>()
  private complianceReports = new Map<string, ComplianceReport>()
  private eventListeners = new Set<(event: ComplianceEvent) => void>()

  private constructor() {
    this.initializeDefaultFrameworks()
  }

  static getInstance(): ComplianceManager {
    if (!ComplianceManager.instance) {
      ComplianceManager.instance = new ComplianceManager()
    }
    return ComplianceManager.instance
  }

  /**
   * Initialize default compliance frameworks
   */
  private initializeDefaultFrameworks(): void {
    const frameworks: ComplianceFramework[] = [
      {
        id: 'us_bsa',
        name: 'US Bank Secrecy Act',
        description: 'US Anti-Money Laundering regulations',
        jurisdiction: 'United States',
        version: '2023.1',
        requirements: [
          {
            id: 'customer_identification',
            name: 'Customer Identification Program',
            description: 'Verify customer identity before account opening',
            category: RequirementCategory.KYC,
            mandatory: true,
            status: ComplianceStatus.COMPLIANT,
            evidence: [],
            controls: [],
            riskLevel: RiskLevel.HIGH
          },
          {
            id: 'suspicious_activity_reporting',
            name: 'Suspicious Activity Reporting',
            description: 'Report suspicious transactions to FinCEN',
            category: RequirementCategory.AML,
            mandatory: true,
            status: ComplianceStatus.COMPLIANT,
            evidence: [],
            controls: [],
            riskLevel: RiskLevel.CRITICAL
          }
        ],
        reportingSchedule: {
          frequency: ReportingFrequency.MONTHLY,
          nextDue: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
          recipients: ['compliance@example.com'],
          format: ReportFormat.PDF,
          automated: true
        },
        lastUpdate: new Date().toISOString(),
        status: ComplianceStatus.COMPLIANT
      },
      {
        id: 'eu_amld5',
        name: 'EU Anti-Money Laundering Directive 5',
        description: 'European Union AML regulations',
        jurisdiction: 'European Union',
        version: '2020.1',
        requirements: [
          {
            id: 'customer_due_diligence',
            name: 'Customer Due Diligence',
            description: 'Perform enhanced due diligence on high-risk customers',
            category: RequirementCategory.KYC,
            mandatory: true,
            status: ComplianceStatus.PARTIAL,
            evidence: [],
            controls: [],
            riskLevel: RiskLevel.HIGH
          }
        ],
        reportingSchedule: {
          frequency: ReportingFrequency.QUARTERLY,
          nextDue: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000).toISOString(),
          recipients: ['eu-compliance@example.com'],
          format: ReportFormat.XML,
          automated: false
        },
        lastUpdate: new Date().toISOString(),
        status: ComplianceStatus.PARTIAL
      }
    ]

    frameworks.forEach(framework => {
      this.frameworks.set(framework.id, framework)
    })
  }

  /**
   * Perform AML check
   */
  async performAMLCheck(
    address: Address,
    checkType: AMLCheckType = AMLCheckType.SANCTIONS_SCREENING
  ): Promise<AMLCheck> {
    const checkId = `aml_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`

    try {
      // Mock AML check - in real app, this would query actual AML databases
      const riskScore = Math.floor(Math.random() * 100)
      const riskLevel = this.getRiskLevel(riskScore)
      
      const findings: AMLFinding[] = []
      
      // Simulate potential findings
      if (riskScore > 70) {
        findings.push({
          id: 'finding_1',
          type: 'sanctions_match',
          description: 'Potential sanctions list match',
          severity: RiskLevel.HIGH,
          source: 'OFAC SDN List',
          confidence: 85,
          details: {
            matchType: 'partial',
            similarity: 0.85
          }
        })
      }

      const amlCheck: AMLCheck = {
        id: checkId,
        address,
        timestamp: new Date().toISOString(),
        checkType,
        result: {
          status: riskScore > 80 ? 'blocked' : riskScore > 50 ? 'alert' : 'clear',
          confidence: Math.floor(Math.random() * 30) + 70,
          falsePositiveRate: 0.05,
          reviewRequired: riskScore > 50,
          autoApproved: riskScore <= 30
        },
        riskScore,
        riskLevel,
        sources: [
          {
            name: 'OFAC SDN List',
            type: 'sanctions_list',
            lastUpdate: new Date().toISOString(),
            reliability: 95
          },
          {
            name: 'EU Sanctions List',
            type: 'sanctions_list',
            lastUpdate: new Date().toISOString(),
            reliability: 90
          }
        ],
        findings,
        recommendations: this.generateAMLRecommendations(riskLevel, findings),
        expiryDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString()
      }

      this.amlChecks.set(checkId, amlCheck)

      // Emit event
      this.emitEvent({
        type: 'aml_check_completed',
        amlCheck,
        timestamp: Date.now()
      })

      return amlCheck

    } catch (error) {
      // Emit error event
      this.emitEvent({
        type: 'aml_check_failed',
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }
  }

  /**
   * Create KYC profile
   */
  async createKYCProfile(
    userAddress: Address,
    tier: KYCTier = KYCTier.BASIC
  ): Promise<KYCProfile> {
    const profileId = `kyc_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`

    const profile: KYCProfile = {
      id: profileId,
      userAddress,
      status: KYCStatus.PENDING,
      tier,
      verificationLevel: VerificationLevel.LEVEL_1,
      documents: [],
      checks: [],
      riskAssessment: {
        overallRisk: RiskLevel.MEDIUM,
        riskScore: 50,
        riskFactors: [],
        mitigatingFactors: [],
        recommendations: [],
        lastAssessed: new Date().toISOString(),
        nextReview: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString()
      },
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }

    this.kycProfiles.set(profileId, profile)

    // Emit event
    this.emitEvent({
      type: 'kyc_profile_created',
      kycProfile: profile,
      timestamp: Date.now()
    })

    return profile
  }

  /**
   * Generate tax report
   */
  async generateTaxReport(
    userAddress: Address,
    taxYear: number,
    jurisdiction: string = 'US',
    reportType: TaxReportType = TaxReportType.CAPITAL_GAINS
  ): Promise<TaxReporting> {
    const reportId = `tax_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`

    // Mock tax calculation
    const transactions: TaxableTransaction[] = [
      {
        id: 'tx_1',
        hash: '0x1234567890abcdef' as Hash,
        timestamp: new Date().toISOString(),
        type: TransactionType.SELL,
        asset: 'BTC',
        amount: '1.0',
        costBasis: '30000',
        fairMarketValue: '45000',
        gainLoss: '15000',
        taxable: true,
        category: TaxCategory.CAPITAL_GAINS_LONG
      }
    ]

    const calculations: TaxCalculation = {
      totalGainLoss: '15000',
      shortTermGains: '0',
      longTermGains: '15000',
      ordinaryIncome: '0',
      taxableIncome: '15000',
      estimatedTax: '3000',
      deductions: '0',
      methodology: 'FIFO'
    }

    const taxReport: TaxReporting = {
      id: reportId,
      userAddress,
      taxYear,
      jurisdiction,
      reportType,
      transactions,
      calculations,
      status: ReportStatus.GENERATED,
      generatedAt: new Date().toISOString(),
      filingDeadline: new Date(taxYear + 1, 3, 15).toISOString() // April 15th
    }

    this.taxReports.set(reportId, taxReport)

    // Emit event
    this.emitEvent({
      type: 'tax_report_generated',
      taxReport,
      timestamp: Date.now()
    })

    return taxReport
  }

  /**
   * Generate compliance report
   */
  async generateComplianceReport(
    frameworkId: string,
    reportType: ReportType,
    period: ReportingPeriod
  ): Promise<ComplianceReport> {
    const reportId = `report_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`

    const framework = this.frameworks.get(frameworkId)
    if (!framework) {
      throw new Error(`Framework not found: ${frameworkId}`)
    }

    // Calculate compliance metrics
    const totalRequirements = framework.requirements.length
    const compliantRequirements = framework.requirements.filter(r => r.status === ComplianceStatus.COMPLIANT).length
    const nonCompliantRequirements = framework.requirements.filter(r => r.status === ComplianceStatus.NON_COMPLIANT).length
    const pendingRequirements = framework.requirements.filter(r => r.status === ComplianceStatus.PENDING).length

    const report: ComplianceReport = {
      id: reportId,
      type: reportType,
      framework: frameworkId,
      period,
      status: ReportStatus.GENERATED,
      data: {
        summary: {
          totalRequirements,
          compliantRequirements,
          nonCompliantRequirements,
          pendingRequirements,
          overallComplianceRate: (compliantRequirements / totalRequirements) * 100,
          riskLevel: this.calculateOverallRiskLevel(framework.requirements)
        },
        metrics: {
          amlChecks: this.amlChecks.size,
          kycVerifications: this.kycProfiles.size,
          transactionsMonitored: 1000, // Mock data
          alertsGenerated: 25,
          falsePositives: 5,
          truePositives: 20,
          averageProcessingTime: 120 // seconds
        },
        findings: this.generateComplianceFindings(framework.requirements),
        recommendations: this.generateComplianceRecommendations(framework.requirements),
        attachments: []
      },
      generatedAt: new Date().toISOString(),
      recipients: framework.reportingSchedule.recipients
    }

    this.complianceReports.set(reportId, report)

    // Emit event
    this.emitEvent({
      type: 'compliance_report_generated',
      complianceReport: report,
      timestamp: Date.now()
    })

    return report
  }

  /**
   * Helper methods
   */
  private getRiskLevel(score: number): RiskLevel {
    if (score >= 80) return RiskLevel.CRITICAL
    if (score >= 60) return RiskLevel.HIGH
    if (score >= 40) return RiskLevel.MEDIUM
    return RiskLevel.LOW
  }

  private generateAMLRecommendations(riskLevel: RiskLevel, findings: AMLFinding[]): string[] {
    const recommendations: string[] = []

    if (riskLevel === RiskLevel.CRITICAL) {
      recommendations.push('Block transaction immediately')
      recommendations.push('Escalate to compliance team')
      recommendations.push('File suspicious activity report')
    } else if (riskLevel === RiskLevel.HIGH) {
      recommendations.push('Require enhanced due diligence')
      recommendations.push('Manual review required')
      recommendations.push('Consider transaction limits')
    } else if (riskLevel === RiskLevel.MEDIUM) {
      recommendations.push('Monitor ongoing activity')
      recommendations.push('Periodic review recommended')
    }

    return recommendations
  }

  private calculateOverallRiskLevel(requirements: ComplianceRequirement[]): RiskLevel {
    const riskScores = requirements.map(r => {
      switch (r.riskLevel) {
        case RiskLevel.CRITICAL: return 4
        case RiskLevel.HIGH: return 3
        case RiskLevel.MEDIUM: return 2
        case RiskLevel.LOW: return 1
        default: return 0
      }
    })

    const averageScore = riskScores.reduce((sum: number, score) => sum + score, 0) / riskScores.length

    if (averageScore >= 3.5) return RiskLevel.CRITICAL
    if (averageScore >= 2.5) return RiskLevel.HIGH
    if (averageScore >= 1.5) return RiskLevel.MEDIUM
    return RiskLevel.LOW
  }

  private generateComplianceFindings(requirements: ComplianceRequirement[]): ComplianceFinding[] {
    return requirements
      .filter(r => r.status === ComplianceStatus.NON_COMPLIANT)
      .map(r => ({
        id: `finding_${r.id}`,
        type: 'non_compliance',
        severity: r.riskLevel,
        description: `Requirement "${r.name}" is not compliant`,
        requirement: r.id,
        evidence: [],
        remediation: 'Implement required controls and procedures',
        responsible: 'Compliance Team'
      }))
  }

  private generateComplianceRecommendations(requirements: ComplianceRequirement[]): string[] {
    const recommendations: string[] = []

    const nonCompliant = requirements.filter(r => r.status === ComplianceStatus.NON_COMPLIANT)
    if (nonCompliant.length > 0) {
      recommendations.push(`Address ${nonCompliant.length} non-compliant requirements`)
    }

    const pending = requirements.filter(r => r.status === ComplianceStatus.PENDING)
    if (pending.length > 0) {
      recommendations.push(`Complete assessment of ${pending.length} pending requirements`)
    }

    recommendations.push('Implement automated compliance monitoring')
    recommendations.push('Regular compliance training for staff')
    recommendations.push('Quarterly compliance reviews')

    return recommendations
  }

  /**
   * Get compliance dashboard
   */
  getComplianceDashboard(): ComplianceDashboard {
    const frameworks = Array.from(this.frameworks.values())
    const amlChecks = Array.from(this.amlChecks.values())
    const kycProfiles = Array.from(this.kycProfiles.values())

    return {
      overview: {
        totalFrameworks: frameworks.length,
        compliantFrameworks: frameworks.filter(f => f.status === ComplianceStatus.COMPLIANT).length,
        totalRequirements: frameworks.reduce((sum, f) => sum + f.requirements.length, 0),
        compliantRequirements: frameworks.reduce((sum, f) => 
          sum + f.requirements.filter(r => r.status === ComplianceStatus.COMPLIANT).length, 0),
        overallComplianceRate: this.calculateOverallComplianceRate(),
        riskLevel: this.calculateOverallRiskLevel(
          frameworks.flatMap(f => f.requirements)
        )
      },
      amlMetrics: {
        totalChecks: amlChecks.length,
        clearResults: amlChecks.filter(c => c.result.status === 'clear').length,
        alertResults: amlChecks.filter(c => c.result.status === 'alert').length,
        blockedResults: amlChecks.filter(c => c.result.status === 'blocked').length,
        averageRiskScore: amlChecks.length > 0 
          ? amlChecks.reduce((sum, c) => sum + c.riskScore, 0) / amlChecks.length 
          : 0
      },
      kycMetrics: {
        totalProfiles: kycProfiles.length,
        verifiedProfiles: kycProfiles.filter(p => p.status === KYCStatus.VERIFIED).length,
        pendingProfiles: kycProfiles.filter(p => p.status === KYCStatus.PENDING).length,
        rejectedProfiles: kycProfiles.filter(p => p.status === KYCStatus.REJECTED).length,
        averageVerificationTime: 24 // hours
      },
      upcomingDeadlines: this.getUpcomingDeadlines(),
      recentReports: Array.from(this.complianceReports.values())
        .sort((a, b) => new Date(b.generatedAt).getTime() - new Date(a.generatedAt).getTime())
        .slice(0, 5)
    }
  }

  private calculateOverallComplianceRate(): number {
    const frameworks = Array.from(this.frameworks.values())
    const totalRequirements = frameworks.reduce((sum, f) => sum + f.requirements.length, 0)
    const compliantRequirements = frameworks.reduce((sum, f) => 
      sum + f.requirements.filter(r => r.status === ComplianceStatus.COMPLIANT).length, 0)
    
    return totalRequirements > 0 ? (compliantRequirements / totalRequirements) * 100 : 0
  }

  private getUpcomingDeadlines(): ComplianceDeadline[] {
    const frameworks = Array.from(this.frameworks.values())
    const deadlines: ComplianceDeadline[] = []

    frameworks.forEach(framework => {
      if (framework.reportingSchedule.nextDue) {
        deadlines.push({
          id: `deadline_${framework.id}`,
          type: 'reporting',
          description: `${framework.name} Report Due`,
          dueDate: framework.reportingSchedule.nextDue,
          priority: 'high',
          framework: framework.id
        })
      }

      framework.requirements.forEach(requirement => {
        if (requirement.deadline) {
          deadlines.push({
            id: `deadline_${requirement.id}`,
            type: 'requirement',
            description: requirement.name,
            dueDate: requirement.deadline,
            priority: requirement.riskLevel === RiskLevel.CRITICAL ? 'high' : 'medium',
            framework: framework.id
          })
        }
      })
    })

    return deadlines.sort((a, b) => new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime())
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: ComplianceEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in compliance event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: ComplianceEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.frameworks.clear()
    this.amlChecks.clear()
    this.kycProfiles.clear()
    this.taxReports.clear()
    this.complianceReports.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface ComplianceDashboard {
  overview: ComplianceOverview
  amlMetrics: AMLMetrics
  kycMetrics: KYCMetrics
  upcomingDeadlines: ComplianceDeadline[]
  recentReports: ComplianceReport[]
}

export interface ComplianceOverview {
  totalFrameworks: number
  compliantFrameworks: number
  totalRequirements: number
  compliantRequirements: number
  overallComplianceRate: number
  riskLevel: RiskLevel
}

export interface AMLMetrics {
  totalChecks: number
  clearResults: number
  alertResults: number
  blockedResults: number
  averageRiskScore: number
}

export interface KYCMetrics {
  totalProfiles: number
  verifiedProfiles: number
  pendingProfiles: number
  rejectedProfiles: number
  averageVerificationTime: number
}

export interface ComplianceDeadline {
  id: string
  type: 'reporting' | 'requirement' | 'review'
  description: string
  dueDate: string
  priority: 'high' | 'medium' | 'low'
  framework: string
}

export interface ComplianceEvent {
  type: 'aml_check_completed' | 'aml_check_failed' | 'kyc_profile_created' | 'kyc_profile_updated' | 'tax_report_generated' | 'compliance_report_generated' | 'framework_updated'
  amlCheck?: AMLCheck
  kycProfile?: KYCProfile
  taxReport?: TaxReporting
  complianceReport?: ComplianceReport
  framework?: ComplianceFramework
  error?: Error
  timestamp: number
}

// Export singleton instance
export const complianceManager = ComplianceManager.getInstance()
