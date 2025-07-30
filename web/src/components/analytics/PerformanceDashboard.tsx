'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import {
  BarChart3,
  TrendingUp,
  Activity,
  Zap,
  Target,
  AlertTriangle,
  CheckCircle,
  Clock,
  RefreshCw,
  Download,
  Settings
} from 'lucide-react'
import {
  PerformanceOverview,
  TradingAnalytics,
  SystemMonitoring,
  OptimizationPanel,
  BenchmarkComparison
} from '.'

interface PerformanceDashboardProps {
  className?: string
}

interface PerformanceMetrics {
  overallScore: number
  tradingScore: number
  systemScore: number
  portfolioScore: number
  optimizationScore: number
}

interface KeyMetrics {
  totalReturn: string
  sharpeRatio: number
  maxDrawdown: string
  successRate: number
  avgLatencyMs: number
  cpuUsage: number
  memoryUsage: number
}

interface SystemHealth {
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  errorRate: number
  uptime: string
}

interface OptimizationData {
  score: number
  opportunitiesCount: number
  potentialImprovement: number
  recommendedActions: string[]
}

export const PerformanceDashboard: React.FC<PerformanceDashboardProps> = ({ className }) => {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null)
  const [keyMetrics, setKeyMetrics] = useState<KeyMetrics | null>(null)
  const [systemHealth, setSystemHealth] = useState<SystemHealth | null>(null)
  const [optimization, setOptimization] = useState<OptimizationData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState('overview')
  const [lastUpdated, setLastUpdated] = useState<string>('')

  useEffect(() => {
    loadDashboardData()
    const interval = setInterval(loadDashboardData, 30000) // Update every 30 seconds
    return () => clearInterval(interval)
  }, [])

  const loadDashboardData = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/analytics/dashboard')
      if (!response.ok) {
        throw new Error('Failed to load analytics dashboard')
      }
      
      const data = await response.json()
      
      setMetrics(data.overview)
      setKeyMetrics(data.key_metrics)
      setSystemHealth(data.system_health)
      setOptimization(data.optimization)
      setLastUpdated(data.last_updated)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-600'
    if (score >= 70) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getScoreBadgeVariant = (score: number) => {
    if (score >= 90) return 'default'
    if (score >= 70) return 'secondary'
    return 'destructive'
  }

  const getHealthStatus = (value: number, threshold: number = 80) => {
    if (value < threshold) return { status: 'healthy', color: 'text-green-600', icon: CheckCircle }
    if (value < 90) return { status: 'warning', color: 'text-yellow-600', icon: AlertTriangle }
    return { status: 'critical', color: 'text-red-600', icon: AlertTriangle }
  }

  if (isLoading) {
    return (
      <div className={`space-y-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading performance dashboard...</p>
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
            Error loading performance dashboard: {error}
            <Button 
              variant="outline" 
              size="sm" 
              className="ml-2"
              onClick={loadDashboardData}
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
          <h1 className="text-3xl font-bold">Performance Analytics</h1>
          <p className="text-muted-foreground">
            Comprehensive performance monitoring and optimization insights
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-sm text-muted-foreground">
            Last updated: {new Date(lastUpdated).toLocaleTimeString()}
          </div>
          <Button variant="outline" onClick={loadDashboardData}>
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

      {/* Performance Scores */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Overall Score</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.overallScore || 0)}`}>
              {metrics?.overallScore || 0}
            </div>
            <Progress value={metrics?.overallScore || 0} className="mt-2" />
            <Badge variant={getScoreBadgeVariant(metrics?.overallScore || 0)} className="mt-2">
              {(metrics?.overallScore || 0) >= 90 ? 'Excellent' : 
               (metrics?.overallScore || 0) >= 70 ? 'Good' : 'Needs Improvement'}
            </Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Trading Score</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.tradingScore || 0)}`}>
              {(metrics?.tradingScore || 0).toFixed(1)}
            </div>
            <Progress value={metrics?.tradingScore || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              Sharpe: {keyMetrics?.sharpeRatio || 0}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">System Score</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.systemScore || 0)}`}>
              {(metrics?.systemScore || 0).toFixed(1)}
            </div>
            <Progress value={metrics?.systemScore || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              CPU: {systemHealth?.cpuUsage || 0}%
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Portfolio Score</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.portfolioScore || 0)}`}>
              {(metrics?.portfolioScore || 0).toFixed(1)}
            </div>
            <Progress value={metrics?.portfolioScore || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              Return: {keyMetrics?.totalReturn || '0%'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Optimization</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(optimization?.score || 0)}`}>
              {(optimization?.score || 0).toFixed(1)}
            </div>
            <Progress value={optimization?.score || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              {optimization?.opportunitiesCount || 0} opportunities
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Return</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {keyMetrics?.totalReturn || '0%'}
            </div>
            <p className="text-xs text-muted-foreground">
              Max Drawdown: {keyMetrics?.maxDrawdown || '0%'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {(keyMetrics?.successRate || 0).toFixed(1)}%
            </div>
            <p className="text-xs text-muted-foreground">
              Trading execution success
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
            <Clock className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {keyMetrics?.avgLatencyMs || 0}ms
            </div>
            <p className="text-xs text-muted-foreground">
              API response time
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">System Health</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-xs">CPU</span>
                <span className={`text-xs font-medium ${getHealthStatus(systemHealth?.cpuUsage || 0).color}`}>
                  {(systemHealth?.cpuUsage || 0).toFixed(1)}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs">Memory</span>
                <span className={`text-xs font-medium ${getHealthStatus(systemHealth?.memoryUsage || 0).color}`}>
                  {(systemHealth?.memoryUsage || 0).toFixed(1)}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs">Disk</span>
                <span className={`text-xs font-medium ${getHealthStatus(systemHealth?.diskUsage || 0).color}`}>
                  {(systemHealth?.diskUsage || 0).toFixed(1)}%
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="trading">Trading</TabsTrigger>
          <TabsTrigger value="system">System</TabsTrigger>
          <TabsTrigger value="optimization">Optimization</TabsTrigger>
          <TabsTrigger value="benchmarks">Benchmarks</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <PerformanceOverview 
            metrics={metrics}
            keyMetrics={keyMetrics}
            systemHealth={systemHealth}
            optimization={optimization}
          />
        </TabsContent>

        <TabsContent value="trading" className="space-y-6">
          <TradingAnalytics />
        </TabsContent>

        <TabsContent value="system" className="space-y-6">
          <SystemMonitoring />
        </TabsContent>

        <TabsContent value="optimization" className="space-y-6">
          <OptimizationPanel />
        </TabsContent>

        <TabsContent value="benchmarks" className="space-y-6">
          <BenchmarkComparison />
        </TabsContent>
      </Tabs>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>
            Common performance analysis and optimization actions
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Button variant="outline" className="h-20 flex-col gap-2">
              <BarChart3 className="h-6 w-6" />
              <span>Generate Report</span>
            </Button>
            <Button variant="outline" className="h-20 flex-col gap-2">
              <Target className="h-6 w-6" />
              <span>Run Optimization</span>
            </Button>
            <Button variant="outline" className="h-20 flex-col gap-2">
              <TrendingUp className="h-6 w-6" />
              <span>Compare Benchmarks</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
