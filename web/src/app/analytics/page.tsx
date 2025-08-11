'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { DatePickerWithRange } from '@/components/ui/date-range-picker'
import { 
  BarChart, 
  Bar, 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  AreaChart,
  Area
} from 'recharts'
import { 
  TrendingUp, 
  TrendingDown, 
  Activity, 
  Users, 
  DollarSign,
  Clock,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Zap,
  Brain,
  Globe,
  Wallet
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { AuthGuard } from '@/components/auth-guard'
import { formatCurrency, formatDuration, formatPercentage } from '@/lib/utils'
import { AdvancedPortfolioAnalytics } from '@/components/analytics/AdvancedPortfolioAnalytics'

interface AnalyticsData {
  metrics: {
    total_executions: number
    success_rate: number
    average_execution_time: number
    cost_analysis: {
      total_cost: number
      cost_per_execution: number
      cost_trend: string
    }
    user_engagement: {
      active_users: number
      sessions_per_user: number
    }
    workflow_performance: Array<{
      id: string
      name: string
      executions: number
      success_rate: number
      avg_duration: number
    }>
    resource_utilization: {
      cpu_utilization: number
      memory_utilization: number
      storage_utilization: number
      concurrent_sessions: number
    }
    error_analysis: {
      total_errors: number
      error_rate: number
      error_types: Array<{
        type: string
        count: number
      }>
    }
  }
  time_series: Array<{
    timestamp: string
    values: {
      executions: number
      success_rate: number
    }
  }>
  insights: Array<{
    type: string
    title: string
    description: string
    severity: 'info' | 'warning' | 'critical'
    value: any
  }>
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8']

export default function AnalyticsPage() {
  const { user } = useAuth()
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [timeRange, setTimeRange] = useState('7d')
  const [granularity, setGranularity] = useState('day')
  const [selectedMetrics, setSelectedMetrics] = useState([
    'total_executions',
    'success_rate',
    'average_execution_time',
    'cost_analysis',
    'user_engagement',
    'workflow_performance',
    'resource_utilization',
    'error_analysis'
  ])

  useEffect(() => {
    loadAnalyticsData()
  }, [timeRange, granularity, selectedMetrics])

  const loadAnalyticsData = async () => {
    setIsLoading(true)
    try {
      // Calculate date range
      const endDate = new Date()
      const startDate = new Date()
      
      switch (timeRange) {
        case '24h':
          startDate.setHours(startDate.getHours() - 24)
          break
        case '7d':
          startDate.setDate(startDate.getDate() - 7)
          break
        case '30d':
          startDate.setDate(startDate.getDate() - 30)
          break
        case '90d':
          startDate.setDate(startDate.getDate() - 90)
          break
      }

      const response = await fetch('/api/analytics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
        body: JSON.stringify({
          start_date: startDate.toISOString(),
          end_date: endDate.toISOString(),
          granularity,
          metrics: selectedMetrics,
        }),
      })

      if (response.ok) {
        const data = await response.json()
        setAnalyticsData(data)
      } else {
        // Mock data for demonstration
        setAnalyticsData(generateMockData())
      }
    } catch (error) {
      console.error('Failed to load analytics:', error)
      setAnalyticsData(generateMockData())
    } finally {
      setIsLoading(false)
    }
  }

  const generateMockData = (): AnalyticsData => {
    const timeSeries = Array.from({ length: 7 }, (_, i) => ({
      timestamp: new Date(Date.now() - (6 - i) * 24 * 60 * 60 * 1000).toISOString(),
      values: {
        executions: Math.floor(Math.random() * 100) + 50,
        success_rate: Math.random() * 20 + 80,
      },
    }))

    return {
      metrics: {
        total_executions: 1247,
        success_rate: 94.2,
        average_execution_time: 12.5,
        cost_analysis: {
          total_cost: 156.75,
          cost_per_execution: 0.125,
          cost_trend: 'decreasing',
        },
        user_engagement: {
          active_users: 23,
          sessions_per_user: 4.2,
        },
        workflow_performance: [
          { id: '1', name: 'Data Scraping', executions: 345, success_rate: 96.5, avg_duration: 8.2 },
          { id: '2', name: 'Form Automation', executions: 234, success_rate: 92.1, avg_duration: 15.3 },
          { id: '3', name: 'Web3 Trading', executions: 189, success_rate: 98.9, avg_duration: 5.7 },
          { id: '4', name: 'Content Analysis', executions: 156, success_rate: 89.7, avg_duration: 22.1 },
        ],
        resource_utilization: {
          cpu_utilization: 65.5,
          memory_utilization: 72.3,
          storage_utilization: 45.8,
          concurrent_sessions: 15,
        },
        error_analysis: {
          total_errors: 72,
          error_rate: 5.8,
          error_types: [
            { type: 'Timeout Error', count: 28 },
            { type: 'Network Error', count: 19 },
            { type: 'Validation Error', count: 15 },
            { type: 'Authentication Error', count: 10 },
          ],
        },
      },
      time_series: timeSeries,
      insights: [
        {
          type: 'success_rate',
          title: 'High Success Rate',
          description: 'Your workflows are performing excellently with 94.2% success rate',
          severity: 'info',
          value: 94.2,
        },
        {
          type: 'cost_optimization',
          title: 'Cost Optimization',
          description: 'Consider optimizing workflows to reduce execution time and costs',
          severity: 'warning',
          value: 0.125,
        },
        {
          type: 'performance_trend',
          title: 'Increasing Usage',
          description: 'Workflow executions have increased by 15% this week',
          severity: 'info',
          value: 15,
        },
      ],
    }
  }

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
        return <XCircle className="h-4 w-4 text-red-500" />
      case 'warning':
        return <AlertTriangle className="h-4 w-4 text-yellow-500" />
      default:
        return <CheckCircle className="h-4 w-4 text-green-500" />
    }
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'border-red-500 bg-red-50'
      case 'warning':
        return 'border-yellow-500 bg-yellow-50'
      default:
        return 'border-green-500 bg-green-50'
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="spinner mb-4" />
          <p className="text-muted-foreground">Loading analytics...</p>
        </div>
      </div>
    )
  }

  if (!analyticsData) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">Failed to load analytics data</p>
          <Button onClick={loadAnalyticsData} className="mt-4">
            Retry
          </Button>
        </div>
      </div>
    )
  }

  return (
    <AuthGuard>
      <div className="min-h-screen bg-background">
        {/* Header */}
        <div className="border-b border-border bg-card">
          <div className="container mx-auto px-4 py-6">
            <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold">Analytics Dashboard</h1>
              <p className="text-muted-foreground mt-1">
                Comprehensive insights into your automation performance
              </p>
            </div>
            <div className="flex items-center gap-4">
              <Select value={timeRange} onValueChange={setTimeRange}>
                <SelectTrigger className="w-32">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="24h">Last 24h</SelectItem>
                  <SelectItem value="7d">Last 7 days</SelectItem>
                  <SelectItem value="30d">Last 30 days</SelectItem>
                  <SelectItem value="90d">Last 90 days</SelectItem>
                </SelectContent>
              </Select>
              <Select value={granularity} onValueChange={setGranularity}>
                <SelectTrigger className="w-32">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="hour">Hourly</SelectItem>
                  <SelectItem value="day">Daily</SelectItem>
                  <SelectItem value="week">Weekly</SelectItem>
                  <SelectItem value="month">Monthly</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Executions</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{analyticsData.metrics.total_executions.toLocaleString()}</div>
                <p className="text-xs text-green-600 flex items-center">
                  <TrendingUp className="h-3 w-3 mr-1" />
                  +12% from last period
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
                <CheckCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatPercentage(analyticsData.metrics.success_rate)}</div>
                <p className="text-xs text-green-600 flex items-center">
                  <TrendingUp className="h-3 w-3 mr-1" />
                  +2.1% from last period
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Avg Execution Time</CardTitle>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatDuration(analyticsData.metrics.average_execution_time)}</div>
                <p className="text-xs text-red-600 flex items-center">
                  <TrendingDown className="h-3 w-3 mr-1" />
                  -8% from last period
                </p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.4 }}
          >
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Cost</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatCurrency(analyticsData.metrics.cost_analysis.total_cost)}</div>
                <p className="text-xs text-green-600 flex items-center">
                  <TrendingDown className="h-3 w-3 mr-1" />
                  -5% from last period
                </p>
              </CardContent>
            </Card>
          </motion.div>
        </div>

        {/* Insights */}
        {analyticsData.insights.length > 0 && (
          <div className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Key Insights</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {analyticsData.insights.map((insight, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.5, delay: index * 0.1 }}
                >
                  <Card className={`border-l-4 ${getSeverityColor(insight.severity)}`}>
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        {getSeverityIcon(insight.severity)}
                        <div className="flex-1">
                          <h3 className="font-medium text-sm">{insight.title}</h3>
                          <p className="text-xs text-muted-foreground mt-1">{insight.description}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          </div>
        )}

        {/* Charts and Detailed Analytics */}
        <Tabs defaultValue="overview" className="w-full">
          <TabsList className="grid w-full grid-cols-6">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="portfolio">Portfolio</TabsTrigger>
            <TabsTrigger value="performance">Performance</TabsTrigger>
            <TabsTrigger value="costs">Costs</TabsTrigger>
            <TabsTrigger value="resources">Resources</TabsTrigger>
            <TabsTrigger value="errors">Errors</TabsTrigger>
          </TabsList>

          <TabsContent value="portfolio" className="space-y-6">
            <AdvancedPortfolioAnalytics />
          </TabsContent>

          <TabsContent value="overview" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Execution Trend</CardTitle>
                  <CardDescription>Daily workflow executions over time</CardDescription>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={300}>
                    <AreaChart data={analyticsData.time_series}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis 
                        dataKey="timestamp" 
                        tickFormatter={(value) => new Date(value).toLocaleDateString()}
                      />
                      <YAxis />
                      <Tooltip 
                        labelFormatter={(value) => new Date(value).toLocaleDateString()}
                      />
                      <Area 
                        type="monotone" 
                        dataKey="values.executions" 
                        stroke="#8884d8" 
                        fill="#8884d8" 
                        fillOpacity={0.3}
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Success Rate Trend</CardTitle>
                  <CardDescription>Workflow success rate over time</CardDescription>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={300}>
                    <LineChart data={analyticsData.time_series}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis 
                        dataKey="timestamp" 
                        tickFormatter={(value) => new Date(value).toLocaleDateString()}
                      />
                      <YAxis domain={[80, 100]} />
                      <Tooltip 
                        labelFormatter={(value) => new Date(value).toLocaleDateString()}
                        formatter={(value: any) => [`${value.toFixed(1)}%`, 'Success Rate']}
                      />
                      <Line 
                        type="monotone" 
                        dataKey="values.success_rate" 
                        stroke="#00C49F" 
                        strokeWidth={2}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="performance" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Top Performing Workflows</CardTitle>
                <CardDescription>Workflows ranked by execution count and success rate</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={400}>
                  <BarChart data={analyticsData.metrics.workflow_performance}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="executions" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="costs" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Cost Breakdown</CardTitle>
                  <CardDescription>Distribution of costs across services</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Brain className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">AI Services</span>
                      </div>
                      <span className="font-medium">{formatCurrency(75.30)}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Globe className="h-4 w-4 text-green-500" />
                        <span className="text-sm">Browser Automation</span>
                      </div>
                      <span className="font-medium">{formatCurrency(25.20)}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Wallet className="h-4 w-4 text-purple-500" />
                        <span className="text-sm">Web3 Services</span>
                      </div>
                      <span className="font-medium">{formatCurrency(15.00)}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Cost Metrics</CardTitle>
                  <CardDescription>Key cost performance indicators</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Cost per Execution</span>
                        <span className="font-medium">{formatCurrency(analyticsData.metrics.cost_analysis.cost_per_execution)}</span>
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Target: {formatCurrency(0.10)} per execution
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Monthly Projection</span>
                        <span className="font-medium">{formatCurrency(analyticsData.metrics.cost_analysis.total_cost * 4)}</span>
                      </div>
                      <div className="text-xs text-green-600">
                        15% below budget
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="resources" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Resource Utilization</CardTitle>
                  <CardDescription>Current system resource usage</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">CPU Usage</span>
                        <span className="font-medium">{analyticsData.metrics.resource_utilization.cpu_utilization}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-blue-600 h-2 rounded-full" 
                          style={{ width: `${analyticsData.metrics.resource_utilization.cpu_utilization}%` }}
                        />
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Memory Usage</span>
                        <span className="font-medium">{analyticsData.metrics.resource_utilization.memory_utilization}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-green-600 h-2 rounded-full" 
                          style={{ width: `${analyticsData.metrics.resource_utilization.memory_utilization}%` }}
                        />
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Storage Usage</span>
                        <span className="font-medium">{analyticsData.metrics.resource_utilization.storage_utilization}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-yellow-600 h-2 rounded-full" 
                          style={{ width: `${analyticsData.metrics.resource_utilization.storage_utilization}%` }}
                        />
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Active Sessions</CardTitle>
                  <CardDescription>Current concurrent browser sessions</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-center">
                    <div className="text-4xl font-bold text-primary mb-2">
                      {analyticsData.metrics.resource_utilization.concurrent_sessions}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Active browser sessions
                    </div>
                    <div className="mt-4 text-xs text-green-600">
                      Within normal limits
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="errors" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Error Distribution</CardTitle>
                  <CardDescription>Most common error types</CardDescription>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={300}>
                    <PieChart>
                      <Pie
                        data={analyticsData.metrics.error_analysis.error_types}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        label={({ type, percent }) => `${type} ${percent ? (percent * 100).toFixed(0) : 0}%`}
                        outerRadius={80}
                        fill="#8884d8"
                        dataKey="count"
                      >
                        {analyticsData.metrics.error_analysis.error_types.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Error Metrics</CardTitle>
                  <CardDescription>Error rate and recovery statistics</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Error Rate</span>
                        <span className="font-medium">{formatPercentage(analyticsData.metrics.error_analysis.error_rate)}</span>
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Target: &lt; 5%
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Total Errors</span>
                        <span className="font-medium">{analyticsData.metrics.error_analysis.total_errors}</span>
                      </div>
                      <div className="text-xs text-green-600">
                        -12% from last period
                      </div>
                    </div>
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm">Mean Time to Recovery</span>
                        <span className="font-medium">45.5 min</span>
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Industry average: 60 min
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
    </AuthGuard>
  )
}
