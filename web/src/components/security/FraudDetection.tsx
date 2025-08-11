'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield, 
  AlertTriangle,
  Target,
  Brain,
  TrendingUp,
  TrendingDown,
  Activity,
  Clock,
  DollarSign,
  Users,
  Eye,
  Ban,
  CheckCircle,
  XCircle,
  Zap,
  BarChart3,
  Settings,
  Download
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface FraudAlert {
  id: string
  type: 'account_takeover' | 'money_laundering' | 'wash_trading' | 'pump_dump' | 'insider_trading' | 'identity_theft'
  severity: 'low' | 'medium' | 'high' | 'critical'
  title: string
  description: string
  timestamp: number
  userId: string
  transactionId?: string
  riskScore: number
  confidence: number
  status: 'active' | 'investigating' | 'resolved' | 'false_positive'
  indicators: string[]
  affectedAmount?: number
  currency?: string
}

interface MLModel {
  id: string
  name: string
  type: 'anomaly_detection' | 'pattern_recognition' | 'behavioral_analysis' | 'risk_scoring'
  description: string
  accuracy: number
  precision: number
  recall: number
  lastTrained: number
  isActive: boolean
  alertsGenerated: number
  falsePositiveRate: number
}

interface RiskProfile {
  userId: string
  overallRiskScore: number
  behaviorScore: number
  transactionScore: number
  networkScore: number
  timeScore: number
  lastUpdated: number
  riskFactors: Array<{
    factor: string
    score: number
    weight: number
    description: string
  }>
}

interface FraudMetric {
  id: string
  name: string
  value: number
  previousValue: number
  unit: string
  trend: 'up' | 'down' | 'stable'
  status: 'good' | 'warning' | 'critical'
  description: string
}

export function FraudDetection() {
  const [fraudAlerts, setFraudAlerts] = useState<FraudAlert[]>([])
  const [mlModels, setMlModels] = useState<MLModel[]>([])
  const [riskProfiles, setRiskProfiles] = useState<RiskProfile[]>([])
  const [metrics, setMetrics] = useState<FraudMetric[]>([])

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock fraud detection data
    const mockFraudAlerts: FraudAlert[] = [
      {
        id: 'fraud1',
        type: 'wash_trading',
        severity: 'high',
        title: 'Wash Trading Pattern Detected',
        description: 'User appears to be trading with themselves across multiple accounts',
        timestamp: Date.now() - 1800000,
        userId: 'user123',
        transactionId: 'tx789',
        riskScore: 85,
        confidence: 92,
        status: 'investigating',
        indicators: ['Circular trading pattern', 'Same IP addresses', 'Timing correlation'],
        affectedAmount: 50000,
        currency: 'USD'
      },
      {
        id: 'fraud2',
        type: 'account_takeover',
        severity: 'critical',
        title: 'Potential Account Takeover',
        description: 'Unusual login pattern followed by large withdrawal attempt',
        timestamp: Date.now() - 3600000,
        userId: 'user456',
        riskScore: 95,
        confidence: 88,
        status: 'active',
        indicators: ['New device login', 'Geographic anomaly', 'Behavioral change'],
        affectedAmount: 25000,
        currency: 'USD'
      },
      {
        id: 'fraud3',
        type: 'money_laundering',
        severity: 'medium',
        title: 'Suspicious Transaction Chain',
        description: 'Complex transaction pattern with multiple intermediaries',
        timestamp: Date.now() - 7200000,
        userId: 'user789',
        transactionId: 'tx456',
        riskScore: 72,
        confidence: 76,
        status: 'resolved',
        indicators: ['Multiple small transactions', 'Rapid succession', 'Cross-border transfers'],
        affectedAmount: 15000,
        currency: 'USD'
      }
    ]

    const mockMlModels: MLModel[] = [
      {
        id: 'model1',
        name: 'Anomaly Detection Engine',
        type: 'anomaly_detection',
        description: 'Detects unusual patterns in user behavior and transactions',
        accuracy: 94.2,
        precision: 91.8,
        recall: 89.5,
        lastTrained: Date.now() - 86400000 * 7,
        isActive: true,
        alertsGenerated: 247,
        falsePositiveRate: 8.2
      },
      {
        id: 'model2',
        name: 'Behavioral Analysis Model',
        type: 'behavioral_analysis',
        description: 'Analyzes user behavior patterns to identify suspicious activities',
        accuracy: 89.7,
        precision: 87.3,
        recall: 92.1,
        lastTrained: Date.now() - 86400000 * 3,
        isActive: true,
        alertsGenerated: 156,
        falsePositiveRate: 12.7
      },
      {
        id: 'model3',
        name: 'Risk Scoring Algorithm',
        type: 'risk_scoring',
        description: 'Calculates comprehensive risk scores for users and transactions',
        accuracy: 92.5,
        precision: 90.1,
        recall: 88.9,
        lastTrained: Date.now() - 86400000 * 1,
        isActive: true,
        alertsGenerated: 89,
        falsePositiveRate: 9.9
      }
    ]

    const mockRiskProfiles: RiskProfile[] = [
      {
        userId: 'user123',
        overallRiskScore: 85,
        behaviorScore: 78,
        transactionScore: 92,
        networkScore: 65,
        timeScore: 88,
        lastUpdated: Date.now() - 3600000,
        riskFactors: [
          { factor: 'High transaction frequency', score: 85, weight: 0.3, description: 'Unusually high number of transactions' },
          { factor: 'Geographic inconsistency', score: 72, weight: 0.2, description: 'Logins from multiple countries' },
          { factor: 'Behavioral anomaly', score: 90, weight: 0.25, description: 'Deviation from normal patterns' },
          { factor: 'Network analysis', score: 65, weight: 0.25, description: 'Connected to suspicious accounts' }
        ]
      },
      {
        userId: 'user456',
        overallRiskScore: 95,
        behaviorScore: 98,
        transactionScore: 89,
        networkScore: 92,
        timeScore: 96,
        lastUpdated: Date.now() - 1800000,
        riskFactors: [
          { factor: 'Account takeover indicators', score: 98, weight: 0.4, description: 'Strong signs of compromised account' },
          { factor: 'Unusual withdrawal pattern', score: 89, weight: 0.3, description: 'Large withdrawals after login' },
          { factor: 'Device fingerprint mismatch', score: 92, weight: 0.2, description: 'New device characteristics' },
          { factor: 'Time-based anomalies', score: 96, weight: 0.1, description: 'Activity outside normal hours' }
        ]
      }
    ]

    const mockMetrics: FraudMetric[] = [
      {
        id: 'metric1',
        name: 'Fraud Detection Rate',
        value: 94.2,
        previousValue: 91.8,
        unit: '%',
        trend: 'up',
        status: 'good',
        description: 'Percentage of fraudulent activities detected'
      },
      {
        id: 'metric2',
        name: 'False Positive Rate',
        value: 8.5,
        previousValue: 12.3,
        unit: '%',
        trend: 'down',
        status: 'good',
        description: 'Percentage of legitimate activities flagged as fraud'
      },
      {
        id: 'metric3',
        name: 'Average Response Time',
        value: 2.3,
        previousValue: 3.1,
        unit: 'minutes',
        trend: 'down',
        status: 'good',
        description: 'Time to detect and alert on fraudulent activity'
      },
      {
        id: 'metric4',
        name: 'Prevented Losses',
        value: 2.4,
        previousValue: 1.8,
        unit: 'M USD',
        trend: 'up',
        status: 'good',
        description: 'Total amount of prevented fraudulent transactions'
      }
    ]

    setFraudAlerts(mockFraudAlerts)
    setMlModels(mockMlModels)
    setRiskProfiles(mockRiskProfiles)
    setMetrics(mockMetrics)
  }, [isConnected])

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const formatCurrency = (amount: number, currency: string = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'text-red-500'
      case 'high': return 'text-orange-500'
      case 'medium': return 'text-yellow-500'
      case 'low': return 'text-blue-500'
      default: return 'text-muted-foreground'
    }
  }

  const getSeverityBadgeVariant = (severity: string) => {
    switch (severity) {
      case 'critical': case 'high': return 'destructive'
      case 'medium': return 'secondary'
      case 'low': return 'outline'
      default: return 'outline'
    }
  }

  const getRiskScoreColor = (score: number) => {
    if (score >= 80) return 'text-red-500'
    if (score >= 60) return 'text-orange-500'
    if (score >= 40) return 'text-yellow-500'
    return 'text-green-500'
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="w-3 h-3" />
      case 'down': return <TrendingDown className="w-3 h-3" />
      case 'stable': return <Activity className="w-3 h-3" />
      default: return null
    }
  }

  const getMetricStatusColor = (status: string) => {
    switch (status) {
      case 'good': return 'text-green-500'
      case 'warning': return 'text-yellow-500'
      case 'critical': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access fraud detection systems
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Fraud Detection</h2>
          <p className="text-muted-foreground">
            AI-powered fraud detection and prevention system
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Brain className="w-3 h-3 mr-1" />
            AI Powered
          </Badge>
          <Button variant="outline" size="sm">
            <Settings className="w-3 h-3 mr-1" />
            Configure
          </Button>
        </div>
      </div>

      {/* Fraud Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {metrics.map((metric) => (
          <Card key={metric.id}>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Target className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">{metric.name}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="text-2xl font-bold">{metric.value}</div>
                <div className={cn("flex items-center gap-1", getMetricStatusColor(metric.status))}>
                  {getTrendIcon(metric.trend)}
                  <span className="text-xs">{metric.unit}</span>
                </div>
              </div>
              <div className="text-xs text-muted-foreground mt-1">
                Previous: {metric.previousValue} {metric.unit}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Active Fraud Alerts */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="w-5 h-5" />
              Active Fraud Alerts
            </CardTitle>
            <Button variant="outline" size="sm">
              <Download className="w-4 h-4 mr-2" />
              Export Report
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {fraudAlerts.filter(alert => alert.status === 'active' || alert.status === 'investigating').map((alert) => (
              <div key={alert.id} className="p-4 border rounded-lg">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-start gap-3">
                    <div className={cn(
                      "w-3 h-3 rounded-full mt-1",
                      alert.severity === 'critical' ? "bg-red-500" :
                      alert.severity === 'high' ? "bg-orange-500" :
                      alert.severity === 'medium' ? "bg-yellow-500" :
                      "bg-blue-500"
                    )} />
                    <div>
                      <h4 className="font-bold">{alert.title}</h4>
                      <p className="text-sm text-muted-foreground">{alert.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={getSeverityBadgeVariant(alert.severity)}>
                      {alert.severity}
                    </Badge>
                    <Badge variant="outline">
                      Risk: {alert.riskScore}%
                    </Badge>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                  <div>
                    <div className="text-muted-foreground">User ID</div>
                    <div className="font-mono">{alert.userId}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Confidence</div>
                    <div className="font-medium">{alert.confidence}%</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Amount</div>
                    <div className="font-medium">
                      {alert.affectedAmount ? formatCurrency(alert.affectedAmount, alert.currency) : 'N/A'}
                    </div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Detected</div>
                    <div className="font-medium">{formatTime(alert.timestamp)}</div>
                  </div>
                </div>

                <div className="space-y-2 mb-4">
                  <div className="text-sm font-medium">Risk Indicators:</div>
                  <div className="flex flex-wrap gap-2">
                    {alert.indicators.map((indicator, index) => (
                      <Badge key={index} variant="outline" className="text-xs">
                        {indicator}
                      </Badge>
                    ))}
                  </div>
                </div>

                <div className="flex gap-2">
                  <Button size="sm">
                    <Eye className="w-3 h-3 mr-1" />
                    Investigate
                  </Button>
                  <Button variant="outline" size="sm">
                    <Ban className="w-3 h-3 mr-1" />
                    Block User
                  </Button>
                  <Button variant="outline" size="sm">
                    <CheckCircle className="w-3 h-3 mr-1" />
                    Mark Resolved
                  </Button>
                  <Button variant="outline" size="sm">
                    <XCircle className="w-3 h-3 mr-1" />
                    False Positive
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* ML Models Performance */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Brain className="w-5 h-5" />
            ML Models Performance
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {mlModels.map((model) => (
              <div key={model.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div>
                    <h4 className="font-bold">{model.name}</h4>
                    <p className="text-sm text-muted-foreground">{model.description}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    {model.isActive ? (
                      <Badge variant="default">
                        <Zap className="w-3 h-3 mr-1" />
                        Active
                      </Badge>
                    ) : (
                      <Badge variant="outline">Inactive</Badge>
                    )}
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm">
                  <div>
                    <div className="text-muted-foreground">Accuracy</div>
                    <div className="font-medium text-green-500">{model.accuracy}%</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Precision</div>
                    <div className="font-medium">{model.precision}%</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Recall</div>
                    <div className="font-medium">{model.recall}%</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">False Positive</div>
                    <div className="font-medium text-yellow-500">{model.falsePositiveRate}%</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Alerts Generated</div>
                    <div className="font-medium">{model.alertsGenerated}</div>
                  </div>
                </div>

                <div className="mt-3 text-xs text-muted-foreground">
                  Last trained: {formatTime(model.lastTrained)}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* High-Risk User Profiles */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="w-5 h-5" />
            High-Risk User Profiles
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {riskProfiles.filter(profile => profile.overallRiskScore >= 70).map((profile) => (
              <div key={profile.userId} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h4 className="font-bold">User ID: {profile.userId}</h4>
                    <p className="text-sm text-muted-foreground">
                      Last updated: {formatTime(profile.lastUpdated)}
                    </p>
                  </div>
                  <div className="text-right">
                    <div className={cn("text-2xl font-bold", getRiskScoreColor(profile.overallRiskScore))}>
                      {profile.overallRiskScore}%
                    </div>
                    <div className="text-sm text-muted-foreground">Risk Score</div>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                  <div>
                    <div className="text-muted-foreground">Behavior</div>
                    <div className={cn("font-medium", getRiskScoreColor(profile.behaviorScore))}>
                      {profile.behaviorScore}%
                    </div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Transactions</div>
                    <div className={cn("font-medium", getRiskScoreColor(profile.transactionScore))}>
                      {profile.transactionScore}%
                    </div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Network</div>
                    <div className={cn("font-medium", getRiskScoreColor(profile.networkScore))}>
                      {profile.networkScore}%
                    </div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Timing</div>
                    <div className={cn("font-medium", getRiskScoreColor(profile.timeScore))}>
                      {profile.timeScore}%
                    </div>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="text-sm font-medium">Risk Factors:</div>
                  <div className="space-y-2">
                    {profile.riskFactors.map((factor, index) => (
                      <div key={index} className="flex items-center justify-between">
                        <div className="flex-1">
                          <div className="text-sm font-medium">{factor.factor}</div>
                          <div className="text-xs text-muted-foreground">{factor.description}</div>
                        </div>
                        <div className="flex items-center gap-2">
                          <div className="w-20">
                            <Progress value={factor.score} className="h-2" />
                          </div>
                          <div className={cn("text-sm font-medium w-12", getRiskScoreColor(factor.score))}>
                            {factor.score}%
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
