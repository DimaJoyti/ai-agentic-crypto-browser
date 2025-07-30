import { type Address, type Hash } from 'viem'

export interface SecurityMonitor {
  id: string
  name: string
  description: string
  enabled: boolean
  monitorType: MonitorType
  targets: MonitorTarget[]
  rules: MonitoringRule[]
  alertConfig: AlertConfiguration
  lastCheck: string
  status: MonitorStatus
  metrics: MonitorMetrics
}

export enum MonitorType {
  WALLET_MONITORING = 'wallet_monitoring',
  CONTRACT_MONITORING = 'contract_monitoring',
  TRANSACTION_MONITORING = 'transaction_monitoring',
  PROTOCOL_MONITORING = 'protocol_monitoring',
  MARKET_MONITORING = 'market_monitoring',
  NETWORK_MONITORING = 'network_monitoring'
}

export interface MonitorTarget {
  id: string
  type: 'address' | 'contract' | 'protocol' | 'token'
  value: string
  chainId?: number
  metadata?: Record<string, any>
}

export interface MonitoringRule {
  id: string
  name: string
  description: string
  condition: RuleCondition
  threshold: RuleThreshold
  severity: AlertSeverity
  enabled: boolean
  cooldown: number
  lastTriggered?: string
}

export interface RuleCondition {
  field: string
  operator: 'equals' | 'not_equals' | 'greater_than' | 'less_than' | 'contains' | 'regex' | 'change_percent'
  value: string | number
  timeWindow?: number
}

export interface RuleThreshold {
  warning: number
  critical: number
  emergency: number
}

export enum AlertSeverity {
  INFO = 'info',
  WARNING = 'warning',
  CRITICAL = 'critical',
  EMERGENCY = 'emergency'
}

export enum MonitorStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  ERROR = 'error',
  MAINTENANCE = 'maintenance'
}

export interface MonitorMetrics {
  checksPerformed: number
  alertsTriggered: number
  lastCheckDuration: number
  averageCheckDuration: number
  uptime: number
  errorRate: number
}

export interface AlertConfiguration {
  channels: AlertChannel[]
  escalation: EscalationRule[]
  suppressions: SuppressionRule[]
  formatting: AlertFormatting
}

export interface AlertChannel {
  type: 'email' | 'sms' | 'push' | 'webhook' | 'slack' | 'discord'
  config: Record<string, any>
  enabled: boolean
  severityFilter: AlertSeverity[]
}

export interface EscalationRule {
  condition: string
  delay: number
  action: string
  targets: string[]
}

export interface SuppressionRule {
  condition: string
  duration: number
  reason: string
}

export interface AlertFormatting {
  template: string
  includeMetadata: boolean
  includeRecommendations: boolean
  maxLength: number
}

export interface SecurityAlert {
  id: string
  monitorId: string
  ruleId: string
  severity: AlertSeverity
  title: string
  description: string
  timestamp: string
  source: AlertSource
  metadata: AlertMetadata
  recommendations: string[]
  status: AlertStatus
  acknowledgment?: AlertAcknowledgment
  resolution?: AlertResolution
}

export interface AlertSource {
  type: string
  identifier: string
  chainId?: number
  blockNumber?: number
  transactionHash?: Hash
}

export interface AlertMetadata {
  affectedAssets: string[]
  estimatedImpact: string
  confidence: number
  falsePositiveRate: number
  relatedAlerts: string[]
  context: Record<string, any>
}

export enum AlertStatus {
  ACTIVE = 'active',
  ACKNOWLEDGED = 'acknowledged',
  RESOLVED = 'resolved',
  SUPPRESSED = 'suppressed',
  FALSE_POSITIVE = 'false_positive'
}

export interface AlertAcknowledgment {
  acknowledgedBy: string
  acknowledgedAt: string
  notes?: string
}

export interface AlertResolution {
  resolvedBy: string
  resolvedAt: string
  resolution: string
  notes?: string
}

export interface ThreatIntelligence {
  threatId: string
  threatType: ThreatType
  severity: AlertSeverity
  description: string
  indicators: ThreatIndicator[]
  attribution: ThreatAttribution
  timeline: ThreatTimeline[]
  mitigation: ThreatMitigation
  confidence: number
  lastUpdated: string
}

export enum ThreatType {
  PHISHING = 'phishing',
  MALWARE = 'malware',
  SCAM = 'scam',
  RUG_PULL = 'rug_pull',
  EXPLOIT = 'exploit',
  FLASH_LOAN_ATTACK = 'flash_loan_attack',
  GOVERNANCE_ATTACK = 'governance_attack',
  ORACLE_MANIPULATION = 'oracle_manipulation',
  BRIDGE_EXPLOIT = 'bridge_exploit',
  SOCIAL_ENGINEERING = 'social_engineering'
}

export interface ThreatIndicator {
  type: 'address' | 'domain' | 'hash' | 'pattern'
  value: string
  confidence: number
  context: string
}

export interface ThreatAttribution {
  actor: string
  group?: string
  motivation: string
  sophistication: 'low' | 'medium' | 'high'
  geography?: string
}

export interface ThreatTimeline {
  timestamp: string
  event: string
  description: string
  impact: string
}

export interface ThreatMitigation {
  immediate: string[]
  shortTerm: string[]
  longTerm: string[]
  preventive: string[]
}

export interface SecurityDashboard {
  overview: SecurityOverview
  activeAlerts: SecurityAlert[]
  recentThreats: ThreatIntelligence[]
  monitorStatus: MonitorStatusSummary
  riskMetrics: SecurityRiskMetrics
  trends: SecurityTrend[]
}

export interface SecurityOverview {
  totalMonitors: number
  activeMonitors: number
  totalAlerts: number
  criticalAlerts: number
  resolvedAlerts: number
  averageResolutionTime: number
  securityScore: number
}

export interface MonitorStatusSummary {
  active: number
  inactive: number
  error: number
  maintenance: number
}

export interface SecurityRiskMetrics {
  overallRisk: number
  threatLevel: AlertSeverity
  exposureScore: number
  vulnerabilityCount: number
  incidentRate: number
  mttr: number // Mean Time To Resolution
}

export interface SecurityTrend {
  metric: string
  timeframe: string
  values: TrendDataPoint[]
  trend: 'increasing' | 'decreasing' | 'stable'
  changePercent: number
}

export interface TrendDataPoint {
  timestamp: string
  value: number
}

export class SecurityMonitoringSystem {
  private static instance: SecurityMonitoringSystem
  private monitors = new Map<string, SecurityMonitor>()
  private alerts = new Map<string, SecurityAlert>()
  private threats = new Map<string, ThreatIntelligence>()
  private eventListeners = new Set<(event: SecurityMonitoringEvent) => void>()
  private monitoringInterval: NodeJS.Timeout | null = null

  private constructor() {
    this.initializeDefaultMonitors()
    this.startMonitoring()
  }

  static getInstance(): SecurityMonitoringSystem {
    if (!SecurityMonitoringSystem.instance) {
      SecurityMonitoringSystem.instance = new SecurityMonitoringSystem()
    }
    return SecurityMonitoringSystem.instance
  }

  /**
   * Initialize default security monitors
   */
  private initializeDefaultMonitors(): void {
    const defaultMonitors: SecurityMonitor[] = [
      {
        id: 'wallet_balance_monitor',
        name: 'Wallet Balance Monitor',
        description: 'Monitors wallet balance changes for suspicious activity',
        enabled: true,
        monitorType: MonitorType.WALLET_MONITORING,
        targets: [],
        rules: [
          {
            id: 'large_outflow',
            name: 'Large Balance Outflow',
            description: 'Detects large outflows from monitored wallets',
            condition: {
              field: 'balance_change',
              operator: 'less_than',
              value: -10000,
              timeWindow: 3600
            },
            threshold: {
              warning: 5000,
              critical: 10000,
              emergency: 50000
            },
            severity: AlertSeverity.CRITICAL,
            enabled: true,
            cooldown: 300
          }
        ],
        alertConfig: {
          channels: [
            {
              type: 'push',
              config: {},
              enabled: true,
              severityFilter: [AlertSeverity.CRITICAL, AlertSeverity.EMERGENCY]
            }
          ],
          escalation: [],
          suppressions: [],
          formatting: {
            template: 'Security Alert: {title}\n{description}',
            includeMetadata: true,
            includeRecommendations: true,
            maxLength: 500
          }
        },
        lastCheck: new Date().toISOString(),
        status: MonitorStatus.ACTIVE,
        metrics: {
          checksPerformed: 0,
          alertsTriggered: 0,
          lastCheckDuration: 0,
          averageCheckDuration: 0,
          uptime: 100,
          errorRate: 0
        }
      },
      {
        id: 'contract_interaction_monitor',
        name: 'Contract Interaction Monitor',
        description: 'Monitors interactions with suspicious contracts',
        enabled: true,
        monitorType: MonitorType.CONTRACT_MONITORING,
        targets: [],
        rules: [
          {
            id: 'suspicious_contract',
            name: 'Suspicious Contract Interaction',
            description: 'Detects interactions with flagged contracts',
            condition: {
              field: 'contract_risk_score',
              operator: 'greater_than',
              value: 80
            },
            threshold: {
              warning: 60,
              critical: 80,
              emergency: 95
            },
            severity: AlertSeverity.WARNING,
            enabled: true,
            cooldown: 60
          }
        ],
        alertConfig: {
          channels: [
            {
              type: 'push',
              config: {},
              enabled: true,
              severityFilter: [AlertSeverity.WARNING, AlertSeverity.CRITICAL]
            }
          ],
          escalation: [],
          suppressions: [],
          formatting: {
            template: 'Contract Alert: {title}\n{description}',
            includeMetadata: true,
            includeRecommendations: true,
            maxLength: 500
          }
        },
        lastCheck: new Date().toISOString(),
        status: MonitorStatus.ACTIVE,
        metrics: {
          checksPerformed: 0,
          alertsTriggered: 0,
          lastCheckDuration: 0,
          averageCheckDuration: 0,
          uptime: 100,
          errorRate: 0
        }
      },
      {
        id: 'transaction_anomaly_monitor',
        name: 'Transaction Anomaly Monitor',
        description: 'Detects anomalous transaction patterns',
        enabled: true,
        monitorType: MonitorType.TRANSACTION_MONITORING,
        targets: [],
        rules: [
          {
            id: 'unusual_gas_price',
            name: 'Unusual Gas Price',
            description: 'Detects transactions with unusually high gas prices',
            condition: {
              field: 'gas_price',
              operator: 'greater_than',
              value: 100000000000, // 100 gwei
              timeWindow: 300
            },
            threshold: {
              warning: 50000000000,
              critical: 100000000000,
              emergency: 200000000000
            },
            severity: AlertSeverity.WARNING,
            enabled: true,
            cooldown: 60
          }
        ],
        alertConfig: {
          channels: [
            {
              type: 'push',
              config: {},
              enabled: true,
              severityFilter: [AlertSeverity.WARNING]
            }
          ],
          escalation: [],
          suppressions: [],
          formatting: {
            template: 'Transaction Alert: {title}\n{description}',
            includeMetadata: false,
            includeRecommendations: false,
            maxLength: 300
          }
        },
        lastCheck: new Date().toISOString(),
        status: MonitorStatus.ACTIVE,
        metrics: {
          checksPerformed: 0,
          alertsTriggered: 0,
          lastCheckDuration: 0,
          averageCheckDuration: 0,
          uptime: 100,
          errorRate: 0
        }
      }
    ]

    defaultMonitors.forEach(monitor => {
      this.monitors.set(monitor.id, monitor)
    })
  }

  /**
   * Start monitoring system
   */
  private startMonitoring(): void {
    if (this.monitoringInterval) {
      clearInterval(this.monitoringInterval)
    }

    this.monitoringInterval = setInterval(() => {
      this.performMonitoringChecks()
    }, 30000) // Check every 30 seconds
  }

  /**
   * Perform monitoring checks
   */
  private async performMonitoringChecks(): Promise<void> {
    for (const monitor of Array.from(this.monitors.values())) {
      if (monitor.enabled && monitor.status === MonitorStatus.ACTIVE) {
        try {
          await this.executeMonitor(monitor)
        } catch (error) {
          console.error(`Error executing monitor ${monitor.id}:`, error)
          monitor.status = MonitorStatus.ERROR
          monitor.metrics.errorRate++
        }
      }
    }
  }

  /**
   * Execute individual monitor
   */
  private async executeMonitor(monitor: SecurityMonitor): Promise<void> {
    const startTime = Date.now()

    try {
      // Mock monitoring logic - in real app, this would check actual conditions
      for (const rule of monitor.rules) {
        if (rule.enabled) {
          const shouldTrigger = await this.evaluateRule(monitor, rule)
          
          if (shouldTrigger) {
            await this.triggerAlert(monitor, rule)
          }
        }
      }

      // Update metrics
      const duration = Date.now() - startTime
      monitor.metrics.checksPerformed++
      monitor.metrics.lastCheckDuration = duration
      monitor.metrics.averageCheckDuration = 
        (monitor.metrics.averageCheckDuration * (monitor.metrics.checksPerformed - 1) + duration) / 
        monitor.metrics.checksPerformed
      monitor.lastCheck = new Date().toISOString()

    } catch (error) {
      monitor.status = MonitorStatus.ERROR
      throw error
    }
  }

  /**
   * Evaluate monitoring rule
   */
  private async evaluateRule(_monitor: SecurityMonitor, rule: MonitoringRule): Promise<boolean> {
    // Mock rule evaluation - in real app, this would check actual conditions
    
    // Check cooldown
    if (rule.lastTriggered) {
      const lastTriggered = new Date(rule.lastTriggered).getTime()
      const now = Date.now()
      if (now - lastTriggered < rule.cooldown * 1000) {
        return false
      }
    }

    // Random trigger for demo (5% chance)
    return Math.random() < 0.05
  }

  /**
   * Trigger security alert
   */
  private async triggerAlert(monitor: SecurityMonitor, rule: MonitoringRule): Promise<void> {
    const alertId = `alert_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`

    const alert: SecurityAlert = {
      id: alertId,
      monitorId: monitor.id,
      ruleId: rule.id,
      severity: rule.severity,
      title: rule.name,
      description: rule.description,
      timestamp: new Date().toISOString(),
      source: {
        type: monitor.monitorType,
        identifier: monitor.id
      },
      metadata: {
        affectedAssets: [],
        estimatedImpact: this.getEstimatedImpact(rule.severity),
        confidence: 85,
        falsePositiveRate: 0.1,
        relatedAlerts: [],
        context: {
          monitorName: monitor.name,
          ruleName: rule.name
        }
      },
      recommendations: this.generateRecommendations(rule),
      status: AlertStatus.ACTIVE
    }

    // Store alert
    this.alerts.set(alertId, alert)

    // Update rule
    rule.lastTriggered = new Date().toISOString()

    // Update monitor metrics
    monitor.metrics.alertsTriggered++

    // Send notifications
    await this.sendAlertNotifications(alert, monitor.alertConfig)

    // Emit event
    this.emitEvent({
      type: 'alert_triggered',
      alert,
      monitor,
      timestamp: Date.now()
    })
  }

  /**
   * Get estimated impact
   */
  private getEstimatedImpact(severity: AlertSeverity): string {
    switch (severity) {
      case AlertSeverity.EMERGENCY:
        return 'Critical - Immediate action required'
      case AlertSeverity.CRITICAL:
        return 'High - Significant risk to assets'
      case AlertSeverity.WARNING:
        return 'Medium - Potential security concern'
      case AlertSeverity.INFO:
        return 'Low - Informational only'
      default:
        return 'Unknown'
    }
  }

  /**
   * Generate recommendations
   */
  private generateRecommendations(rule: MonitoringRule): string[] {
    const recommendations: string[] = []

    switch (rule.id) {
      case 'large_outflow':
        recommendations.push('Verify the transaction was authorized')
        recommendations.push('Check for unauthorized access to your wallet')
        recommendations.push('Consider moving remaining funds to a secure wallet')
        break
      case 'suspicious_contract':
        recommendations.push('Do not interact with the flagged contract')
        recommendations.push('Verify the contract address and source code')
        recommendations.push('Check community reports and audit status')
        break
      case 'unusual_gas_price':
        recommendations.push('Consider waiting for lower gas prices')
        recommendations.push('Check if the transaction is urgent')
        recommendations.push('Use gas optimization tools')
        break
      default:
        recommendations.push('Review the alert details carefully')
        recommendations.push('Take appropriate security measures')
    }

    return recommendations
  }

  /**
   * Send alert notifications
   */
  private async sendAlertNotifications(alert: SecurityAlert, config: AlertConfiguration): Promise<void> {
    for (const channel of config.channels) {
      if (channel.enabled && channel.severityFilter.includes(alert.severity)) {
        try {
          await this.sendNotification(alert, channel, config.formatting)
        } catch (error) {
          console.error(`Failed to send notification via ${channel.type}:`, error)
        }
      }
    }
  }

  /**
   * Send individual notification
   */
  private async sendNotification(
    alert: SecurityAlert,
    channel: AlertChannel,
    formatting: AlertFormatting
  ): Promise<void> {
    // Mock notification sending - in real app, this would send actual notifications
    const message = this.formatAlertMessage(alert, formatting)
    
    console.log(`Sending ${channel.type} notification:`, message)
    
    // Simulate notification delay
    await new Promise(resolve => setTimeout(resolve, 100))
  }

  /**
   * Format alert message
   */
  private formatAlertMessage(alert: SecurityAlert, formatting: AlertFormatting): string {
    let message = formatting.template
      .replace('{title}', alert.title)
      .replace('{description}', alert.description)
      .replace('{severity}', alert.severity.toUpperCase())
      .replace('{timestamp}', alert.timestamp)

    if (formatting.includeMetadata) {
      message += `\n\nImpact: ${alert.metadata.estimatedImpact}`
      message += `\nConfidence: ${alert.metadata.confidence}%`
    }

    if (formatting.includeRecommendations && alert.recommendations.length > 0) {
      message += '\n\nRecommendations:'
      alert.recommendations.forEach((rec, index) => {
        message += `\n${index + 1}. ${rec}`
      })
    }

    // Truncate if too long
    if (message.length > formatting.maxLength) {
      message = message.substring(0, formatting.maxLength - 3) + '...'
    }

    return message
  }

  /**
   * Add monitor
   */
  addMonitor(monitor: SecurityMonitor): void {
    this.monitors.set(monitor.id, monitor)
    
    this.emitEvent({
      type: 'monitor_added',
      monitor,
      timestamp: Date.now()
    })
  }

  /**
   * Remove monitor
   */
  removeMonitor(monitorId: string): boolean {
    const monitor = this.monitors.get(monitorId)
    if (monitor) {
      this.monitors.delete(monitorId)
      
      this.emitEvent({
        type: 'monitor_removed',
        monitor,
        timestamp: Date.now()
      })
      
      return true
    }
    return false
  }

  /**
   * Update monitor
   */
  updateMonitor(monitorId: string, updates: Partial<SecurityMonitor>): boolean {
    const monitor = this.monitors.get(monitorId)
    if (monitor) {
      Object.assign(monitor, updates)
      
      this.emitEvent({
        type: 'monitor_updated',
        monitor,
        timestamp: Date.now()
      })
      
      return true
    }
    return false
  }

  /**
   * Acknowledge alert
   */
  acknowledgeAlert(alertId: string, acknowledgedBy: string, notes?: string): boolean {
    const alert = this.alerts.get(alertId)
    if (alert && alert.status === AlertStatus.ACTIVE) {
      alert.status = AlertStatus.ACKNOWLEDGED
      alert.acknowledgment = {
        acknowledgedBy,
        acknowledgedAt: new Date().toISOString(),
        notes
      }

      this.emitEvent({
        type: 'alert_acknowledged',
        alert,
        timestamp: Date.now()
      })

      return true
    }
    return false
  }

  /**
   * Resolve alert
   */
  resolveAlert(alertId: string, resolvedBy: string, resolution: string, notes?: string): boolean {
    const alert = this.alerts.get(alertId)
    if (alert && (alert.status === AlertStatus.ACTIVE || alert.status === AlertStatus.ACKNOWLEDGED)) {
      alert.status = AlertStatus.RESOLVED
      alert.resolution = {
        resolvedBy,
        resolvedAt: new Date().toISOString(),
        resolution,
        notes
      }

      this.emitEvent({
        type: 'alert_resolved',
        alert,
        timestamp: Date.now()
      })

      return true
    }
    return false
  }

  /**
   * Get security dashboard
   */
  getSecurityDashboard(): SecurityDashboard {
    const activeAlerts = Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.ACTIVE)
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())

    const criticalAlerts = activeAlerts.filter(alert => 
      alert.severity === AlertSeverity.CRITICAL || alert.severity === AlertSeverity.EMERGENCY
    )

    const resolvedAlerts = Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.RESOLVED)

    const monitors = Array.from(this.monitors.values())
    const activeMonitors = monitors.filter(m => m.status === MonitorStatus.ACTIVE)

    return {
      overview: {
        totalMonitors: monitors.length,
        activeMonitors: activeMonitors.length,
        totalAlerts: this.alerts.size,
        criticalAlerts: criticalAlerts.length,
        resolvedAlerts: resolvedAlerts.length,
        averageResolutionTime: this.calculateAverageResolutionTime(),
        securityScore: this.calculateSecurityScore()
      },
      activeAlerts: activeAlerts.slice(0, 10),
      recentThreats: Array.from(this.threats.values()).slice(0, 5),
      monitorStatus: {
        active: monitors.filter(m => m.status === MonitorStatus.ACTIVE).length,
        inactive: monitors.filter(m => m.status === MonitorStatus.INACTIVE).length,
        error: monitors.filter(m => m.status === MonitorStatus.ERROR).length,
        maintenance: monitors.filter(m => m.status === MonitorStatus.MAINTENANCE).length
      },
      riskMetrics: {
        overallRisk: this.calculateOverallRisk(),
        threatLevel: this.calculateThreatLevel(),
        exposureScore: this.calculateExposureScore(),
        vulnerabilityCount: this.calculateVulnerabilityCount(),
        incidentRate: this.calculateIncidentRate(),
        mttr: this.calculateAverageResolutionTime()
      },
      trends: this.generateSecurityTrends()
    }
  }

  /**
   * Calculate average resolution time
   */
  private calculateAverageResolutionTime(): number {
    const resolvedAlerts = Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.RESOLVED && alert.resolution)

    if (resolvedAlerts.length === 0) return 0

    const totalTime = resolvedAlerts.reduce((sum, alert) => {
      const created = new Date(alert.timestamp).getTime()
      const resolved = new Date(alert.resolution!.resolvedAt).getTime()
      return sum + (resolved - created)
    }, 0)

    return totalTime / resolvedAlerts.length / 1000 / 60 // Convert to minutes
  }

  /**
   * Calculate security score
   */
  private calculateSecurityScore(): number {
    const monitors = Array.from(this.monitors.values())
    const activeMonitors = monitors.filter(m => m.status === MonitorStatus.ACTIVE)
    const activeAlerts = Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.ACTIVE)

    let score = 100

    // Deduct for inactive monitors
    score -= (monitors.length - activeMonitors.length) * 10

    // Deduct for active alerts
    score -= activeAlerts.length * 5

    // Deduct more for critical alerts
    const criticalAlerts = activeAlerts.filter(alert => 
      alert.severity === AlertSeverity.CRITICAL || alert.severity === AlertSeverity.EMERGENCY
    )
    score -= criticalAlerts.length * 15

    return Math.max(score, 0)
  }

  /**
   * Calculate overall risk
   */
  private calculateOverallRisk(): number {
    const securityScore = this.calculateSecurityScore()
    return 100 - securityScore
  }

  /**
   * Calculate threat level
   */
  private calculateThreatLevel(): AlertSeverity {
    const activeAlerts = Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.ACTIVE)

    if (activeAlerts.some(alert => alert.severity === AlertSeverity.EMERGENCY)) {
      return AlertSeverity.EMERGENCY
    }
    if (activeAlerts.some(alert => alert.severity === AlertSeverity.CRITICAL)) {
      return AlertSeverity.CRITICAL
    }
    if (activeAlerts.some(alert => alert.severity === AlertSeverity.WARNING)) {
      return AlertSeverity.WARNING
    }
    return AlertSeverity.INFO
  }

  /**
   * Calculate exposure score
   */
  private calculateExposureScore(): number {
    // Mock calculation
    return Math.floor(Math.random() * 100)
  }

  /**
   * Calculate vulnerability count
   */
  private calculateVulnerabilityCount(): number {
    return Array.from(this.alerts.values())
      .filter(alert => alert.status === AlertStatus.ACTIVE)
      .length
  }

  /**
   * Calculate incident rate
   */
  private calculateIncidentRate(): number {
    const now = Date.now()
    const oneDayAgo = now - 24 * 60 * 60 * 1000

    const recentAlerts = Array.from(this.alerts.values())
      .filter(alert => new Date(alert.timestamp).getTime() > oneDayAgo)

    return recentAlerts.length
  }

  /**
   * Generate security trends
   */
  private generateSecurityTrends(): SecurityTrend[] {
    // Mock trend data
    return [
      {
        metric: 'Alert Volume',
        timeframe: '24h',
        values: [
          { timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(), value: 5 },
          { timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(), value: 8 },
          { timestamp: new Date().toISOString(), value: 12 }
        ],
        trend: 'increasing',
        changePercent: 140
      },
      {
        metric: 'Security Score',
        timeframe: '7d',
        values: [
          { timestamp: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(), value: 85 },
          { timestamp: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(), value: 82 },
          { timestamp: new Date().toISOString(), value: 78 }
        ],
        trend: 'decreasing',
        changePercent: -8.2
      }
    ]
  }

  /**
   * Get monitors
   */
  getMonitors(): SecurityMonitor[] {
    return Array.from(this.monitors.values())
  }

  /**
   * Get alerts
   */
  getAlerts(status?: AlertStatus): SecurityAlert[] {
    const alerts = Array.from(this.alerts.values())
    return status ? alerts.filter(alert => alert.status === status) : alerts
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: SecurityMonitoringEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in security monitoring event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: SecurityMonitoringEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Stop monitoring
   */
  stopMonitoring(): void {
    if (this.monitoringInterval) {
      clearInterval(this.monitoringInterval)
      this.monitoringInterval = null
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.monitors.clear()
    this.alerts.clear()
    this.threats.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.stopMonitoring()
    this.clear()
    this.eventListeners.clear()
  }
}

export interface SecurityMonitoringEvent {
  type: 'monitor_added' | 'monitor_removed' | 'monitor_updated' | 'alert_triggered' | 'alert_acknowledged' | 'alert_resolved' | 'threat_detected'
  monitor?: SecurityMonitor
  alert?: SecurityAlert
  threat?: ThreatIntelligence
  timestamp: number
}

// Export singleton instance
export const securityMonitoringSystem = SecurityMonitoringSystem.getInstance()
