'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { motion } from 'framer-motion'
import {
  Shield,
  AlertTriangle,
  CheckCircle,
  FileText,
  TrendingUp,
  Settings,
  Download,
  RefreshCw,
  Clock,
  Bell,
  Sparkles,
  ShieldCheck,
  TrendingDown
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
  const [lastUpdated, setLastUpdated] = useState<string | null>(null)

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
      setLastUpdated(new Date().toISOString())
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
    <div className={`space-y-8 ${className}`}>
      {/* Enhanced Header with Gradient Background */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-600 via-teal-600 to-cyan-700 p-8 text-white"
      >
        <div className="absolute inset-0 bg-grid-pattern opacity-10" />
        <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />

        <div className="relative flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
          <div className="space-y-2">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.2, duration: 0.6 }}
              className="flex items-center gap-3"
            >
              <div className="p-3 bg-white/20 rounded-xl backdrop-blur-sm">
                <ShieldCheck className="h-8 w-8" />
              </div>
              <div>
                <h1 className="text-4xl font-bold tracking-tight">Compliance Dashboard</h1>
                <p className="text-emerald-100 text-lg">
                  Regulatory compliance and risk management
                </p>
              </div>
            </motion.div>

            {lastUpdated && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.4, duration: 0.6 }}
                className="flex items-center gap-2 text-emerald-200"
              >
                <Clock className="h-4 w-4" />
                <span className="text-sm">Last updated: {new Date(lastUpdated).toLocaleString()}</span>
              </motion.div>
            )}
          </div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.3, duration: 0.6 }}
            className="flex items-center gap-3"
          >
            <Button
              variant="secondary"
              size="sm"
              onClick={loadComplianceData}
              disabled={isLoading}
              className="bg-white/20 hover:bg-white/30 text-white border-white/30 backdrop-blur-sm"
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Button
              variant="secondary"
              size="sm"
              className="bg-white/20 hover:bg-white/30 text-white border-white/30 backdrop-blur-sm"
            >
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
            <Button
              variant="secondary"
              className="bg-white/20 hover:bg-white/30 text-white border-white/30 backdrop-blur-sm"
            >
              <Settings className="h-4 w-4 mr-2" />
              Settings
            </Button>
          </motion.div>
        </div>
      </motion.div>

      {/* Enhanced Key Metrics */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5, duration: 0.6 }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
      >
        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/50 dark:to-teal-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-emerald-700 dark:text-emerald-300">Compliance Rate</CardTitle>
              <div className="p-2 bg-emerald-100 dark:bg-emerald-900/50 rounded-lg">
                <Shield className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className="text-3xl font-bold text-emerald-600 dark:text-emerald-400">
                  {metrics?.overallComplianceRate || 0}%
                </div>
                <Sparkles className="h-5 w-5 text-yellow-500 animate-pulse" />
              </div>
              <Progress
                value={metrics?.overallComplianceRate || 0}
                className="mt-3 h-2 bg-emerald-100 dark:bg-emerald-900/30"
              />
              <p className="text-xs text-emerald-600 dark:text-emerald-400 mt-3 font-medium">
                {metrics?.compliantFrameworks || 0} of {metrics?.totalFrameworks || 0} frameworks compliant
              </p>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-orange-50 to-red-50 dark:from-orange-950/50 dark:to-red-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-orange-700 dark:text-orange-300">Risk Level</CardTitle>
              <div className="p-2 bg-orange-100 dark:bg-orange-900/50 rounded-lg">
                {metrics?.riskLevel === 'LOW' ?
                  <TrendingDown className="h-4 w-4 text-green-600 dark:text-green-400" /> :
                  <TrendingUp className="h-4 w-4 text-orange-600 dark:text-orange-400" />
                }
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className="text-3xl font-bold text-orange-600 dark:text-orange-400">
                  {metrics?.riskLevel || 'Unknown'}
                </div>
                <div className={`w-3 h-3 rounded-full animate-pulse ${
                  metrics?.riskLevel === 'LOW' ? 'bg-green-500' :
                  metrics?.riskLevel === 'MEDIUM' ? 'bg-yellow-500' :
                  metrics?.riskLevel === 'HIGH' ? 'bg-red-500' : 'bg-gray-500'
                }`} />
              </div>
              <div className="mt-3 p-2 bg-orange-100 dark:bg-orange-900/30 rounded-lg">
                <p className="text-xs text-orange-600 dark:text-orange-400 font-medium">
                  Risk Score: {metrics?.riskScore || 0}/100
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-red-50 to-pink-50 dark:from-red-950/50 dark:to-pink-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-red-700 dark:text-red-300">Active Violations</CardTitle>
              <div className="p-2 bg-red-100 dark:bg-red-900/50 rounded-lg">
                <AlertTriangle className="h-4 w-4 text-red-600 dark:text-red-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className="text-3xl font-bold text-red-600 dark:text-red-400">
                  {metrics?.activeViolations || 0}
                </div>
                {(metrics?.activeViolations || 0) > 0 && (
                  <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
                )}
              </div>
              <div className="mt-3 p-2 bg-red-100 dark:bg-red-900/30 rounded-lg">
                <p className="text-xs text-red-600 dark:text-red-400 font-medium">
                  {metrics?.resolvedViolations || 0} resolved this month
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-950/50 dark:to-indigo-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-blue-700 dark:text-blue-300">Pending Reports</CardTitle>
              <div className="p-2 bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                <FileText className="h-4 w-4 text-blue-600 dark:text-blue-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className="text-3xl font-bold text-blue-600 dark:text-blue-400">
                  {metrics?.pendingReports || 0}
                </div>
                <Clock className="h-5 w-5 text-yellow-500 animate-pulse" />
              </div>
              <div className="mt-3 p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                <p className="text-xs text-blue-600 dark:text-blue-400 font-medium">
                  Due within 7 days
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </motion.div>

      {/* Enhanced Main Content Tabs */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.7, duration: 0.6 }}
      >
        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-5 bg-gradient-to-r from-emerald-100 to-teal-200 dark:from-emerald-800 dark:to-teal-900 p-1 rounded-xl shadow-lg">
            <TabsTrigger
              value="overview"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-emerald-600 transition-all duration-300"
            >
              Overview
            </TabsTrigger>
            <TabsTrigger
              value="risk"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-emerald-600 transition-all duration-300"
            >
              Risk Monitoring
            </TabsTrigger>
            <TabsTrigger
              value="audit"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-emerald-600 transition-all duration-300"
            >
              Audit Trail
            </TabsTrigger>
            <TabsTrigger
              value="reports"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-emerald-600 transition-all duration-300"
            >
              Reports
            </TabsTrigger>
            <TabsTrigger
              value="alerts"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-emerald-600 transition-all duration-300"
            >
              Alerts
            </TabsTrigger>
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
      </motion.div>

      {/* Enhanced Recent Activity Summary */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.9, duration: 0.6 }}
        className="grid grid-cols-1 lg:grid-cols-2 gap-6"
      >
        <Card className="glass-card border-0 shadow-xl bg-gradient-to-br from-red-50 to-orange-50 dark:from-red-950/30 dark:to-orange-950/30">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-xl">
              <div className="p-2 bg-gradient-to-br from-red-500 to-orange-600 rounded-lg">
                <AlertTriangle className="h-5 w-5 text-white" />
              </div>
              Recent Violations
            </CardTitle>
            <CardDescription className="text-base">
              Latest compliance violations requiring attention
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {violations.slice(0, 5).map((violation, index) => (
                <motion.div
                  key={violation.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1, duration: 0.3 }}
                  className="flex items-center justify-between p-4 bg-white/50 dark:bg-black/20 border border-red-200 dark:border-red-800 rounded-xl hover:shadow-md transition-all duration-300"
                >
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-red-100 dark:bg-red-900/50 rounded-lg">
                      <AlertTriangle className="h-4 w-4 text-red-600 dark:text-red-400" />
                    </div>
                    <div>
                      <p className="font-medium text-sm text-red-800 dark:text-red-200">{violation.description}</p>
                      <p className="text-xs text-red-600 dark:text-red-400">
                        {violation.framework} • {new Date(violation.timestamp).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={getSeverityColor(violation.severity)} className="shadow-sm">
                      {violation.severity}
                    </Badge>
                    {violation.resolved && (
                      <CheckCircle className="h-4 w-4 text-green-500" />
                    )}
                  </div>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="glass-card border-0 shadow-xl bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-950/30 dark:to-cyan-950/30">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-xl">
              <div className="p-2 bg-gradient-to-br from-blue-500 to-cyan-600 rounded-lg">
                <Bell className="h-5 w-5 text-white" />
              </div>
              Risk Alerts
            </CardTitle>
            <CardDescription className="text-base">
              Active risk monitoring alerts requiring review
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {alerts.slice(0, 5).map((alert, index) => (
                <motion.div
                  key={alert.id}
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1, duration: 0.3 }}
                  className="flex items-center justify-between p-4 bg-white/50 dark:bg-black/20 border border-blue-200 dark:border-blue-800 rounded-xl hover:shadow-md transition-all duration-300"
                >
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                      <TrendingUp className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                    </div>
                    <div>
                      <p className="font-medium text-sm text-blue-800 dark:text-blue-200">{alert.title}</p>
                      <p className="text-xs text-blue-600 dark:text-blue-400">
                        {alert.type} • {new Date(alert.timestamp).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={getSeverityColor(alert.severity)} className="shadow-sm">
                      {alert.severity}
                    </Badge>
                    {alert.acknowledged && (
                      <CheckCircle className="h-4 w-4 text-green-500" />
                    )}
                  </div>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}
