'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import {
  Shield,
  AlertTriangle,
  CheckCircle,
  FileText,
  TrendingUp,
  Settings,
  Download,
  RefreshCw
} from 'lucide-react'
import {
  ComplianceOverview,
  RiskMonitoring,
  AuditTrail,
  ComplianceReports,
  AlertManagement
} from '.'

interface ComplianceDashboardProps {
  className?: string
}

interface ComplianceMetrics {
  overallComplianceRate: number
  riskLevel: string
  totalFrameworks: number
  compliantFrameworks: number
  activeViolations: number
  pendingReports: number
  complianceScore: number
  riskScore: number
  violationsThisMonth: number
  resolvedViolations: number
}

interface ComplianceFramework {
  id: string
  name: string
  description: string
  jurisdiction: string
  status: string
  lastUpdate: string
  complianceRate: number
}

interface ComplianceViolation {
  id: string
  type: string
  severity: string
  description: string
  timestamp: string
  resolved: boolean
  framework: string
}

interface RiskAlert {
  id: string
  type: string
  severity: string
  title: string
  description: string
  timestamp: string
  acknowledged: boolean
  resolved: boolean
}

export const ComplianceDashboard: React.FC<ComplianceDashboardProps> = ({ className }) => {
  const [metrics, setMetrics] = useState<ComplianceMetrics | null>(null)
  const [frameworks, setFrameworks] = useState<ComplianceFramework[]>([])
  const [violations, setViolations] = useState<ComplianceViolation[]>([])
  const [alerts, setAlerts] = useState<RiskAlert[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    loadComplianceData()
  }, [])

  const loadComplianceData = async () => {
    setIsLoading(true)
    setError(null)

    try {
      // Load compliance dashboard data
      const dashboardResponse = await fetch('/api/compliance/dashboard')
      if (!dashboardResponse.ok) {
        throw new Error('Failed to load compliance dashboard')
      }
      const dashboardData = await dashboardResponse.json()

      // Load frameworks
      const frameworksResponse = await fetch('/api/compliance/frameworks')
      if (!frameworksResponse.ok) {
        throw new Error('Failed to load compliance frameworks')
      }
      const frameworksData = await frameworksResponse.json()

      // Load violations
      const violationsResponse = await fetch('/api/compliance/violations?limit=10')
      if (!violationsResponse.ok) {
        throw new Error('Failed to load compliance violations')
      }
      const violationsData = await violationsResponse.json()

      // Load risk alerts
      const alertsResponse = await fetch('/api/risk/alerts?limit=10')
      if (!alertsResponse.ok) {
        throw new Error('Failed to load risk alerts')
      }
      const alertsData = await alertsResponse.json()

      setMetrics(dashboardData.metrics)
      setFrameworks(frameworksData.frameworks || [])
      setViolations(violationsData.violations || [])
      setAlerts(alertsData.alerts || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred')
    } finally {
      setIsLoading(false)
    }
  }



  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical':
        return 'destructive'
      case 'high':
        return 'destructive'
      case 'medium':
        return 'default'
      case 'low':
        return 'secondary'
      default:
        return 'outline'
    }
  }

  if (isLoading) {
    return (
      <div className={`space-y-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading compliance dashboard...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className={`space-y-6 ${className}`}>
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            Error loading compliance dashboard: {error}
            <Button 
              variant="outline" 
              size="sm" 
              className="ml-2"
              onClick={loadComplianceData}
            >
              Retry
            </Button>
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Compliance & Risk Management</h1>
          <p className="text-muted-foreground">
            Monitor regulatory compliance and risk metrics
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Button variant="outline" onClick={loadComplianceData}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Export Report
          </Button>
          <Button variant="outline">
            <Settings className="h-4 w-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Compliance Rate</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.overallComplianceRate || 0}%</div>
            <Progress value={metrics?.overallComplianceRate || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              {metrics?.compliantFrameworks || 0} of {metrics?.totalFrameworks || 0} frameworks compliant
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Risk Level</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.riskLevel || 'Unknown'}</div>
            <div className="flex items-center gap-2 mt-2">
              <div className={`w-2 h-2 rounded-full ${
                metrics?.riskLevel === 'LOW' ? 'bg-green-500' :
                metrics?.riskLevel === 'MEDIUM' ? 'bg-yellow-500' :
                metrics?.riskLevel === 'HIGH' ? 'bg-red-500' : 'bg-gray-500'
              }`} />
              <span className="text-xs text-muted-foreground">
                Risk Score: {metrics?.riskScore || 0}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Violations</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.activeViolations || 0}</div>
            <p className="text-xs text-muted-foreground">
              {metrics?.resolvedViolations || 0} resolved this month
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending Reports</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.pendingReports || 0}</div>
            <p className="text-xs text-muted-foreground">
              Due within 7 days
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="risk">Risk Monitoring</TabsTrigger>
          <TabsTrigger value="audit">Audit Trail</TabsTrigger>
          <TabsTrigger value="reports">Reports</TabsTrigger>
          <TabsTrigger value="alerts">Alerts</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <ComplianceOverview 
            frameworks={frameworks}
            violations={violations}
            alerts={alerts}
            metrics={metrics}
          />
        </TabsContent>

        <TabsContent value="risk" className="space-y-6">
          <RiskMonitoring />
        </TabsContent>

        <TabsContent value="audit" className="space-y-6">
          <AuditTrail />
        </TabsContent>

        <TabsContent value="reports" className="space-y-6">
          <ComplianceReports />
        </TabsContent>

        <TabsContent value="alerts" className="space-y-6">
          <AlertManagement />
        </TabsContent>
      </Tabs>

      {/* Recent Activity Summary */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Recent Violations</CardTitle>
            <CardDescription>Latest compliance violations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {violations.slice(0, 5).map((violation) => (
                <div key={violation.id} className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    <AlertTriangle className="h-4 w-4 text-orange-500" />
                    <div>
                      <p className="font-medium text-sm">{violation.description}</p>
                      <p className="text-xs text-muted-foreground">
                        {violation.framework} • {new Date(violation.timestamp).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={getSeverityColor(violation.severity)}>
                      {violation.severity}
                    </Badge>
                    {violation.resolved && (
                      <CheckCircle className="h-4 w-4 text-green-500" />
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Risk Alerts</CardTitle>
            <CardDescription>Active risk monitoring alerts</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {alerts.slice(0, 5).map((alert) => (
                <div key={alert.id} className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    <TrendingUp className="h-4 w-4 text-blue-500" />
                    <div>
                      <p className="font-medium text-sm">{alert.title}</p>
                      <p className="text-xs text-muted-foreground">
                        {alert.type} • {new Date(alert.timestamp).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={getSeverityColor(alert.severity)}>
                      {alert.severity}
                    </Badge>
                    {alert.acknowledged && (
                      <CheckCircle className="h-4 w-4 text-green-500" />
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
