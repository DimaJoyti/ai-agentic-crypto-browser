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
  Settings,
  Sparkles,
  Gauge,
  TrendingDown
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
    <div className={`space-y-8 ${className}`}>
      {/* Enhanced Header with Gradient Background */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-blue-600 via-purple-600 to-indigo-700 p-8 text-white"
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
                <Gauge className="h-8 w-8" />
              </div>
              <div>
                <h1 className="text-4xl font-bold tracking-tight">Performance Dashboard</h1>
                <p className="text-blue-100 text-lg">
                  Real-time system performance and analytics
                </p>
              </div>
            </motion.div>

            {lastUpdated && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.4, duration: 0.6 }}
                className="flex items-center gap-2 text-blue-200"
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
              onClick={loadDashboardData}
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

      {/* Enhanced Performance Scores */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5, duration: 0.6 }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6"
      >
        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/50 dark:to-teal-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-emerald-700 dark:text-emerald-300">Overall Score</CardTitle>
              <div className="p-2 bg-emerald-100 dark:bg-emerald-900/50 rounded-lg">
                <BarChart3 className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className={`text-3xl font-bold ${getScoreColor(metrics?.overallScore || 0)}`}>
                  {metrics?.overallScore || 0}
                </div>
                <Sparkles className="h-5 w-5 text-yellow-500 animate-pulse" />
              </div>
              <Progress
                value={metrics?.overallScore || 0}
                className="mt-3 h-2 bg-emerald-100 dark:bg-emerald-900/30"
              />
              <Badge
                variant={getScoreBadgeVariant(metrics?.overallScore || 0)}
                className="mt-3 bg-emerald-100 text-emerald-800 dark:bg-emerald-900/50 dark:text-emerald-200"
              >
                {(metrics?.overallScore || 0) >= 90 ? 'Excellent' :
                 (metrics?.overallScore || 0) >= 70 ? 'Good' : 'Needs Improvement'}
              </Badge>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-950/50 dark:to-indigo-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-blue-700 dark:text-blue-300">Trading Score</CardTitle>
              <div className="p-2 bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                <TrendingUp className="h-4 w-4 text-blue-600 dark:text-blue-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className={`text-3xl font-bold ${getScoreColor(metrics?.tradingScore || 0)}`}>
                  {(metrics?.tradingScore || 0).toFixed(1)}
                </div>
                {(metrics?.tradingScore || 0) > 75 ?
                  <TrendingUp className="h-5 w-5 text-green-500" /> :
                  <TrendingDown className="h-5 w-5 text-red-500" />
                }
              </div>
              <Progress
                value={metrics?.tradingScore || 0}
                className="mt-3 h-2 bg-blue-100 dark:bg-blue-900/30"
              />
              <p className="text-xs text-blue-600 dark:text-blue-400 mt-3 font-medium">
                Sharpe Ratio: {keyMetrics?.sharpeRatio || 0}
              </p>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-950/50 dark:to-pink-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-purple-700 dark:text-purple-300">System Score</CardTitle>
              <div className="p-2 bg-purple-100 dark:bg-purple-900/50 rounded-lg">
                <Activity className="h-4 w-4 text-purple-600 dark:text-purple-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className={`text-3xl font-bold ${getScoreColor(metrics?.systemScore || 0)}`}>
                  {(metrics?.systemScore || 0).toFixed(1)}
                </div>
                <div className="flex items-center gap-1">
                  <div className={`w-2 h-2 rounded-full ${
                    (systemHealth?.cpuUsage || 0) < 70 ? 'bg-green-500' :
                    (systemHealth?.cpuUsage || 0) < 85 ? 'bg-yellow-500' : 'bg-red-500'
                  } animate-pulse`} />
                </div>
              </div>
              <Progress
                value={metrics?.systemScore || 0}
                className="mt-3 h-2 bg-purple-100 dark:bg-purple-900/30"
              />
              <p className="text-xs text-purple-600 dark:text-purple-400 mt-3 font-medium">
                CPU Usage: {systemHealth?.cpuUsage || 0}%
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
              <CardTitle className="text-sm font-medium text-orange-700 dark:text-orange-300">Portfolio Score</CardTitle>
              <div className="p-2 bg-orange-100 dark:bg-orange-900/50 rounded-lg">
                <Target className="h-4 w-4 text-orange-600 dark:text-orange-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className={`text-3xl font-bold ${getScoreColor(metrics?.portfolioScore || 0)}`}>
                  {(metrics?.portfolioScore || 0).toFixed(1)}
                </div>
                <div className="text-sm font-semibold text-green-600 dark:text-green-400">
                  {keyMetrics?.totalReturn || '0%'}
                </div>
              </div>
              <Progress
                value={metrics?.portfolioScore || 0}
                className="mt-3 h-2 bg-orange-100 dark:bg-orange-900/30"
              />
              <p className="text-xs text-orange-600 dark:text-orange-400 mt-3 font-medium">
                Total Return: {keyMetrics?.totalReturn || '0%'}
              </p>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.02, y: -5 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="glass-card border-0 shadow-xl hover:shadow-2xl transition-all duration-300 bg-gradient-to-br from-yellow-50 to-amber-50 dark:from-yellow-950/50 dark:to-amber-950/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-yellow-700 dark:text-yellow-300">Optimization</CardTitle>
              <div className="p-2 bg-yellow-100 dark:bg-yellow-900/50 rounded-lg">
                <Zap className="h-4 w-4 text-yellow-600 dark:text-yellow-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-3">
                <div className={`text-3xl font-bold ${getScoreColor(optimization?.score || 0)}`}>
                  {(optimization?.score || 0).toFixed(1)}
                </div>
                <div className="flex items-center gap-1">
                  <Sparkles className="h-4 w-4 text-yellow-500 animate-pulse" />
                </div>
              </div>
              <Progress
                value={optimization?.score || 0}
                className="mt-3 h-2 bg-yellow-100 dark:bg-yellow-900/30"
              />
              <p className="text-xs text-yellow-600 dark:text-yellow-400 mt-3 font-medium">
                {optimization?.opportunitiesCount || 0} opportunities
              </p>
            </CardContent>
          </Card>
        </motion.div>
      </motion.div>

      {/* Enhanced Key Metrics */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.7, duration: 0.6 }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
      >
        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-950/30 dark:to-emerald-950/30">
            <div className="absolute top-0 right-0 w-20 h-20 bg-green-500/10 rounded-full -translate-y-10 translate-x-10" />
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-green-700 dark:text-green-300">Total Return</CardTitle>
              <div className="p-2 bg-green-100 dark:bg-green-900/50 rounded-lg">
                <TrendingUp className="h-4 w-4 text-green-600 dark:text-green-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">
                {keyMetrics?.totalReturn || '0%'}
              </div>
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-red-500 rounded-full" />
                <p className="text-xs text-green-600 dark:text-green-400 font-medium">
                  Max Drawdown: {keyMetrics?.maxDrawdown || '0%'}
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-950/30 dark:to-cyan-950/30">
            <div className="absolute top-0 right-0 w-20 h-20 bg-blue-500/10 rounded-full -translate-y-10 translate-x-10" />
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-blue-700 dark:text-blue-300">Success Rate</CardTitle>
              <div className="p-2 bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                <CheckCircle className="h-4 w-4 text-blue-600 dark:text-blue-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">
                {(keyMetrics?.successRate || 0).toFixed(1)}%
              </div>
              <div className="flex items-center gap-2">
                <div className={`w-2 h-2 rounded-full ${
                  (keyMetrics?.successRate || 0) > 80 ? 'bg-green-500' :
                  (keyMetrics?.successRate || 0) > 60 ? 'bg-yellow-500' : 'bg-red-500'
                } animate-pulse`} />
                <p className="text-xs text-blue-600 dark:text-blue-400 font-medium">
                  Trading execution success
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-purple-50 to-violet-50 dark:from-purple-950/30 dark:to-violet-950/30">
            <div className="absolute top-0 right-0 w-20 h-20 bg-purple-500/10 rounded-full -translate-y-10 translate-x-10" />
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-purple-700 dark:text-purple-300">Avg Latency</CardTitle>
              <div className="p-2 bg-purple-100 dark:bg-purple-900/50 rounded-lg">
                <Clock className="h-4 w-4 text-purple-600 dark:text-purple-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-purple-600 dark:text-purple-400 mb-2">
                {keyMetrics?.avgLatencyMs || 0}ms
              </div>
              <div className="flex items-center gap-2">
                <div className={`w-2 h-2 rounded-full ${
                  (keyMetrics?.avgLatencyMs || 0) < 100 ? 'bg-green-500' :
                  (keyMetrics?.avgLatencyMs || 0) < 300 ? 'bg-yellow-500' : 'bg-red-500'
                } animate-pulse`} />
                <p className="text-xs text-purple-600 dark:text-purple-400 font-medium">
                  API response time
                </p>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          whileHover={{ scale: 1.05 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Card className="relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-indigo-50 to-blue-50 dark:from-indigo-950/30 dark:to-blue-950/30">
            <div className="absolute top-0 right-0 w-20 h-20 bg-indigo-500/10 rounded-full -translate-y-10 translate-x-10" />
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-indigo-700 dark:text-indigo-300">System Health</CardTitle>
              <div className="p-2 bg-indigo-100 dark:bg-indigo-900/50 rounded-lg">
                <Activity className="h-4 w-4 text-indigo-600 dark:text-indigo-400" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-indigo-600 dark:text-indigo-400">CPU</span>
                  <div className="flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${getHealthStatus(systemHealth?.cpuUsage || 0).color.includes('green') ? 'bg-green-500' : getHealthStatus(systemHealth?.cpuUsage || 0).color.includes('yellow') ? 'bg-yellow-500' : 'bg-red-500'}`} />
                    <span className={`text-xs font-bold ${getHealthStatus(systemHealth?.cpuUsage || 0).color}`}>
                      {(systemHealth?.cpuUsage || 0).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-indigo-600 dark:text-indigo-400">Memory</span>
                  <div className="flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${getHealthStatus(systemHealth?.memoryUsage || 0).color.includes('green') ? 'bg-green-500' : getHealthStatus(systemHealth?.memoryUsage || 0).color.includes('yellow') ? 'bg-yellow-500' : 'bg-red-500'}`} />
                    <span className={`text-xs font-bold ${getHealthStatus(systemHealth?.memoryUsage || 0).color}`}>
                      {(systemHealth?.memoryUsage || 0).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-indigo-600 dark:text-indigo-400">Disk</span>
                  <div className="flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${getHealthStatus(systemHealth?.diskUsage || 0).color.includes('green') ? 'bg-green-500' : getHealthStatus(systemHealth?.diskUsage || 0).color.includes('yellow') ? 'bg-yellow-500' : 'bg-red-500'}`} />
                    <span className={`text-xs font-bold ${getHealthStatus(systemHealth?.diskUsage || 0).color}`}>
                      {(systemHealth?.diskUsage || 0).toFixed(1)}%
                    </span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </motion.div>

      {/* Enhanced Main Content Tabs */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.9, duration: 0.6 }}
      >
        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-5 bg-gradient-to-r from-slate-100 to-slate-200 dark:from-slate-800 dark:to-slate-900 p-1 rounded-xl shadow-lg">
            <TabsTrigger
              value="overview"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-blue-600 transition-all duration-300"
            >
              Overview
            </TabsTrigger>
            <TabsTrigger
              value="trading"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-blue-600 transition-all duration-300"
            >
              Trading
            </TabsTrigger>
            <TabsTrigger
              value="system"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-blue-600 transition-all duration-300"
            >
              System
            </TabsTrigger>
            <TabsTrigger
              value="optimization"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-blue-600 transition-all duration-300"
            >
              Optimization
            </TabsTrigger>
            <TabsTrigger
              value="benchmarks"
              className="data-[state=active]:bg-white data-[state=active]:shadow-md data-[state=active]:text-blue-600 transition-all duration-300"
            >
              Benchmarks
            </TabsTrigger>
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
      </motion.div>

      {/* Enhanced Quick Actions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 1.1, duration: 0.6 }}
      >
        <Card className="glass-card border-0 shadow-xl bg-gradient-to-br from-slate-50 to-gray-50 dark:from-slate-900/50 dark:to-gray-900/50">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-xl">
              <div className="p-2 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg">
                <Sparkles className="h-5 w-5 text-white" />
              </div>
              Quick Actions
            </CardTitle>
            <CardDescription className="text-base">
              Common performance analysis and optimization actions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <motion.div
                whileHover={{ scale: 1.05, y: -5 }}
                whileTap={{ scale: 0.95 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Button
                  variant="outline"
                  className="h-24 flex-col gap-3 bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-950/30 dark:to-indigo-950/30 border-blue-200 dark:border-blue-800 hover:shadow-lg transition-all duration-300"
                >
                  <div className="p-2 bg-blue-100 dark:bg-blue-900/50 rounded-lg">
                    <BarChart3 className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                  </div>
                  <span className="font-medium text-blue-700 dark:text-blue-300">Generate Report</span>
                </Button>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.05, y: -5 }}
                whileTap={{ scale: 0.95 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Button
                  variant="outline"
                  className="h-24 flex-col gap-3 bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-950/30 dark:to-emerald-950/30 border-green-200 dark:border-green-800 hover:shadow-lg transition-all duration-300"
                >
                  <div className="p-2 bg-green-100 dark:bg-green-900/50 rounded-lg">
                    <Target className="h-6 w-6 text-green-600 dark:text-green-400" />
                  </div>
                  <span className="font-medium text-green-700 dark:text-green-300">Run Optimization</span>
                </Button>
              </motion.div>

              <motion.div
                whileHover={{ scale: 1.05, y: -5 }}
                whileTap={{ scale: 0.95 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Button
                  variant="outline"
                  className="h-24 flex-col gap-3 bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-950/30 dark:to-pink-950/30 border-purple-200 dark:border-purple-800 hover:shadow-lg transition-all duration-300"
                >
                  <div className="p-2 bg-purple-100 dark:bg-purple-900/50 rounded-lg">
                    <TrendingUp className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                  </div>
                  <span className="font-medium text-purple-700 dark:text-purple-300">Compare Benchmarks</span>
                </Button>
              </motion.div>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}
