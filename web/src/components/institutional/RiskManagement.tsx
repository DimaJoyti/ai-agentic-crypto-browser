'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield, 
  AlertTriangle, 
  TrendingDown,
  TrendingUp,
  Target,
  BarChart3,
  Activity,
  Clock,
  Zap,
  Eye,
  Settings,
  Bell,
  CheckCircle,
  XCircle
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface RiskMetric {
  id: string
  name: string
  value: number
  threshold: number
  warningLevel: number
  unit: string
  status: 'normal' | 'warning' | 'critical'
  trend: 'up' | 'down' | 'stable'
  description: string
}

interface RiskLimit {
  id: string
  name: string
  type: 'position' | 'exposure' | 'concentration' | 'leverage' | 'var'
  currentValue: number
  limitValue: number
  utilizationPercent: number
  status: 'compliant' | 'warning' | 'breach'
  asset?: string
  timeframe: string
}

interface RiskAlert {
  id: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  type: 'limit_breach' | 'concentration' | 'volatility' | 'correlation' | 'liquidity'
  title: string
  description: string
  timestamp: number
  isAcknowledged: boolean
  affectedAssets: string[]
  recommendedAction: string
}

interface StressTestScenario {
  id: string
  name: string
  description: string
  type: 'market_crash' | 'liquidity_crisis' | 'correlation_breakdown' | 'custom'
  parameters: Record<string, number>
  results: {
    portfolioImpact: number
    worstAsset: string
    timeToRecover: number
    liquidityImpact: number
  }
  lastRun: number
  status: 'passed' | 'warning' | 'failed'
}

export function RiskManagement() {
  const [riskMetrics, setRiskMetrics] = useState<RiskMetric[]>([])
  const [riskLimits, setRiskLimits] = useState<RiskLimit[]>([])
  const [riskAlerts, setRiskAlerts] = useState<RiskAlert[]>([])
  const [stressTests, setStressTests] = useState<StressTestScenario[]>([])
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    // Generate mock risk data
    const mockMetrics: RiskMetric[] = [
      {
        id: 'var',
        name: 'Value at Risk (95%)',
        value: 2.8,
        threshold: 5.0,
        warningLevel: 4.0,
        unit: '%',
        status: 'normal',
        trend: 'stable',
        description: '1-day VaR at 95% confidence level'
      },
      {
        id: 'sharpe',
        name: 'Sharpe Ratio',
        value: 1.85,
        threshold: 1.0,
        warningLevel: 1.2,
        unit: '',
        status: 'normal',
        trend: 'up',
        description: 'Risk-adjusted return metric'
      },
      {
        id: 'beta',
        name: 'Portfolio Beta',
        value: 1.23,
        threshold: 1.5,
        warningLevel: 1.3,
        unit: '',
        status: 'warning',
        trend: 'up',
        description: 'Systematic risk relative to market'
      },
      {
        id: 'concentration',
        name: 'Concentration Risk',
        value: 35.2,
        threshold: 40.0,
        warningLevel: 35.0,
        unit: '%',
        status: 'warning',
        trend: 'up',
        description: 'Largest single position as % of portfolio'
      }
    ]

    const mockLimits: RiskLimit[] = [
      {
        id: 'limit1',
        name: 'Single Asset Exposure',
        type: 'concentration',
        currentValue: 15000000,
        limitValue: 20000000,
        utilizationPercent: 75,
        status: 'warning',
        asset: 'BTC',
        timeframe: 'Real-time'
      },
      {
        id: 'limit2',
        name: 'Portfolio Leverage',
        type: 'leverage',
        currentValue: 2.1,
        limitValue: 3.0,
        utilizationPercent: 70,
        status: 'compliant',
        timeframe: 'Daily'
      },
      {
        id: 'limit3',
        name: 'Daily VaR Limit',
        type: 'var',
        currentValue: 2800000,
        limitValue: 5000000,
        utilizationPercent: 56,
        status: 'compliant',
        timeframe: 'Daily'
      }
    ]

    const mockAlerts: RiskAlert[] = [
      {
        id: 'alert1',
        severity: 'high',
        type: 'concentration',
        title: 'High Concentration Risk',
        description: 'BTC position exceeds 75% of concentration limit',
        timestamp: Date.now() - 1800000,
        isAcknowledged: false,
        affectedAssets: ['BTC'],
        recommendedAction: 'Consider reducing BTC exposure or increasing portfolio diversification'
      },
      {
        id: 'alert2',
        severity: 'medium',
        type: 'volatility',
        title: 'Increased Volatility',
        description: 'Portfolio volatility has increased by 25% over the last 24 hours',
        timestamp: Date.now() - 3600000,
        isAcknowledged: true,
        affectedAssets: ['ETH', 'SOL'],
        recommendedAction: 'Monitor positions closely and consider hedging strategies'
      }
    ]

    const mockStressTests: StressTestScenario[] = [
      {
        id: 'stress1',
        name: 'Market Crash Scenario',
        description: '30% market decline with increased correlations',
        type: 'market_crash',
        parameters: { marketDecline: -30, correlationIncrease: 0.8 },
        results: {
          portfolioImpact: -18.5,
          worstAsset: 'ETH',
          timeToRecover: 180,
          liquidityImpact: -25.0
        },
        lastRun: Date.now() - 86400000,
        status: 'warning'
      },
      {
        id: 'stress2',
        name: 'Liquidity Crisis',
        description: 'Severe liquidity constraints across all assets',
        type: 'liquidity_crisis',
        parameters: { liquidityReduction: -60, spreadIncrease: 5.0 },
        results: {
          portfolioImpact: -12.3,
          worstAsset: 'MATIC',
          timeToRecover: 90,
          liquidityImpact: -45.0
        },
        lastRun: Date.now() - 172800000,
        status: 'passed'
      }
    ]

    setRiskMetrics(mockMetrics)
    setRiskLimits(mockLimits)
    setRiskAlerts(mockAlerts)
    setStressTests(mockStressTests)
  }, [])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'normal': case 'compliant': case 'passed': return 'text-green-500'
      case 'warning': return 'text-yellow-500'
      case 'critical': case 'breach': case 'failed': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'normal': case 'compliant': case 'passed': return 'default'
      case 'warning': return 'secondary'
      case 'critical': case 'breach': case 'failed': return 'destructive'
      default: return 'outline'
    }
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'low': return 'text-blue-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-orange-500'
      case 'critical': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="w-3 h-3 text-red-500" />
      case 'down': return <TrendingDown className="w-3 h-3 text-green-500" />
      case 'stable': return <Activity className="w-3 h-3 text-muted-foreground" />
      default: return null
    }
  }

  const acknowledgeAlert = (alertId: string) => {
    setRiskAlerts(prev => prev.map(alert => 
      alert.id === alertId ? { ...alert, isAcknowledged: true } : alert
    ))
  }

  const getOverallRiskScore = () => {
    const criticalCount = riskMetrics.filter(m => m.status === 'critical').length
    const warningCount = riskMetrics.filter(m => m.status === 'warning').length
    
    if (criticalCount > 0) return { score: 'High', color: 'text-red-500' }
    if (warningCount > 1) return { score: 'Medium', color: 'text-yellow-500' }
    return { score: 'Low', color: 'text-green-500' }
  }

  const overallRisk = getOverallRiskScore()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Risk Management</h2>
          <p className="text-muted-foreground">
            Comprehensive risk monitoring and control systems
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Risk Score: <span className={overallRisk.color}>{overallRisk.score}</span>
          </Badge>
          <Button variant="outline" size="sm">
            <Settings className="w-3 h-3 mr-1" />
            Configure
          </Button>
        </div>
      </div>

      {/* Active Alerts */}
      {riskAlerts.filter(alert => !alert.isAcknowledged).length > 0 && (
        <div className="space-y-2">
          {riskAlerts.filter(alert => !alert.isAcknowledged).map((alert) => (
            <Alert key={alert.id} className={cn(
              "border-l-4",
              alert.severity === 'critical' ? "border-l-red-500" :
              alert.severity === 'high' ? "border-l-orange-500" :
              alert.severity === 'medium' ? "border-l-yellow-500" :
              "border-l-blue-500"
            )}>
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription className="flex items-center justify-between">
                <div>
                  <div className="font-medium">{alert.title}</div>
                  <div className="text-sm text-muted-foreground">{alert.description}</div>
                </div>
                <Button 
                  variant="outline" 
                  size="sm"
                  onClick={() => acknowledgeAlert(alert.id)}
                >
                  Acknowledge
                </Button>
              </AlertDescription>
            </Alert>
          ))}
        </div>
      )}

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Target className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Portfolio VaR</span>
            </div>
            <div className="text-2xl font-bold">
              {riskMetrics.find(m => m.id === 'var')?.value}%
            </div>
            <div className="text-xs text-muted-foreground">
              95% confidence, 1-day
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <BarChart3 className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Active Limits</span>
            </div>
            <div className="text-2xl font-bold">{riskLimits.length}</div>
            <div className="text-xs text-red-500">
              {riskLimits.filter(l => l.status === 'breach').length} breaches
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Bell className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Open Alerts</span>
            </div>
            <div className="text-2xl font-bold">
              {riskAlerts.filter(a => !a.isAcknowledged).length}
            </div>
            <div className="text-xs text-muted-foreground">
              {riskAlerts.filter(a => a.severity === 'high' || a.severity === 'critical').length} high priority
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Clock className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Last Stress Test</span>
            </div>
            <div className="text-2xl font-bold">
              {Math.floor((Date.now() - Math.max(...stressTests.map(s => s.lastRun))) / (24 * 60 * 60 * 1000))}d
            </div>
            <div className="text-xs text-muted-foreground">
              ago
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="limits">Risk Limits</TabsTrigger>
          <TabsTrigger value="alerts">Alerts</TabsTrigger>
          <TabsTrigger value="stress">Stress Tests</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Risk Metrics */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Key Risk Metrics
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {riskMetrics.map((metric) => (
                    <div key={metric.id} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{metric.name}</span>
                          {getTrendIcon(metric.trend)}
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="font-bold">
                            {metric.value}{metric.unit}
                          </span>
                          <Badge variant={getStatusBadgeVariant(metric.status)}>
                            {metric.status}
                          </Badge>
                        </div>
                      </div>
                      <Progress 
                        value={(metric.value / metric.threshold) * 100} 
                        className="h-2" 
                      />
                      <div className="flex justify-between text-xs text-muted-foreground">
                        <span>{metric.description}</span>
                        <span>Limit: {metric.threshold}{metric.unit}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Recent Alerts */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Bell className="w-5 h-5" />
                  Recent Alerts
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {riskAlerts.slice(0, 5).map((alert) => (
                    <div key={alert.id} className="flex items-start gap-3 p-3 border rounded">
                      <div className={cn(
                        "w-2 h-2 rounded-full mt-2",
                        alert.severity === 'critical' ? "bg-red-500" :
                        alert.severity === 'high' ? "bg-orange-500" :
                        alert.severity === 'medium' ? "bg-yellow-500" :
                        "bg-blue-500"
                      )} />
                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <div className="font-medium">{alert.title}</div>
                          <div className="flex items-center gap-1">
                            {alert.isAcknowledged ? (
                              <CheckCircle className="w-4 h-4 text-green-500" />
                            ) : (
                              <XCircle className="w-4 h-4 text-red-500" />
                            )}
                          </div>
                        </div>
                        <div className="text-sm text-muted-foreground">{alert.description}</div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {formatTime(alert.timestamp)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="limits" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Risk Limits</h3>
            <Button>
              <Target className="w-4 h-4 mr-2" />
              Add New Limit
            </Button>
          </div>

          <div className="space-y-4">
            {riskLimits.map((limit) => (
              <Card key={limit.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{limit.name}</h4>
                      <p className="text-sm text-muted-foreground capitalize">
                        {limit.type.replace('_', ' ')} • {limit.timeframe}
                      </p>
                    </div>
                    <Badge variant={getStatusBadgeVariant(limit.status)}>
                      {limit.status}
                    </Badge>
                  </div>

                  <div className="space-y-3">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Utilization</span>
                      <span className="font-medium">{limit.utilizationPercent}%</span>
                    </div>
                    <Progress value={limit.utilizationPercent} className="h-2" />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>
                        Current: {limit.type === 'leverage' ? limit.currentValue.toFixed(1) : formatCurrency(limit.currentValue)}
                      </span>
                      <span>
                        Limit: {limit.type === 'leverage' ? limit.limitValue.toFixed(1) : formatCurrency(limit.limitValue)}
                      </span>
                    </div>
                  </div>

                  {limit.asset && (
                    <div className="mt-3 text-sm text-muted-foreground">
                      Asset: {limit.asset}
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="alerts" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Risk Alerts</h3>
            <Button variant="outline">
              <Eye className="w-4 h-4 mr-2" />
              Mark All Read
            </Button>
          </div>

          <div className="space-y-4">
            {riskAlerts.map((alert) => (
              <Card key={alert.id}>
                <CardContent className="p-6">
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
                      <Badge variant="outline" className={getSeverityColor(alert.severity)}>
                        {alert.severity}
                      </Badge>
                      {alert.isAcknowledged ? (
                        <CheckCircle className="w-4 h-4 text-green-500" />
                      ) : (
                        <Button 
                          variant="outline" 
                          size="sm"
                          onClick={() => acknowledgeAlert(alert.id)}
                        >
                          Acknowledge
                        </Button>
                      )}
                    </div>
                  </div>

                  <div className="text-sm">
                    <div className="text-muted-foreground mb-2">Recommended Action:</div>
                    <div className="bg-muted/50 p-3 rounded">{alert.recommendedAction}</div>
                  </div>

                  <div className="flex items-center justify-between mt-3 text-xs text-muted-foreground">
                    <span>Affected: {alert.affectedAssets.join(', ')}</span>
                    <span>{formatTime(alert.timestamp)}</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="stress" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Stress Test Scenarios</h3>
            <Button>
              <Zap className="w-4 h-4 mr-2" />
              Run New Test
            </Button>
          </div>

          <div className="space-y-4">
            {stressTests.map((test) => (
              <Card key={test.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{test.name}</h4>
                      <p className="text-sm text-muted-foreground">{test.description}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={getStatusBadgeVariant(test.status)}>
                        {test.status}
                      </Badge>
                      <Button variant="outline" size="sm">
                        Run Test
                      </Button>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Portfolio Impact</div>
                      <div className={cn(
                        "font-bold",
                        test.results.portfolioImpact < 0 ? "text-red-500" : "text-green-500"
                      )}>
                        {test.results.portfolioImpact > 0 ? '+' : ''}{test.results.portfolioImpact}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Worst Asset</div>
                      <div className="font-medium">{test.results.worstAsset}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Recovery Time</div>
                      <div className="font-medium">{test.results.timeToRecover} days</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Liquidity Impact</div>
                      <div className="font-medium text-red-500">
                        {test.results.liquidityImpact}%
                      </div>
                    </div>
                  </div>

                  <div className="mt-3 text-xs text-muted-foreground">
                    Last run: {formatTime(test.lastRun)}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
