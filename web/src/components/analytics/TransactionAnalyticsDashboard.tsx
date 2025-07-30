'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  BarChart3,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Activity,
  Clock,
  Zap,
  Shield,
  Download,
  RefreshCw,
  Filter,
  AlertTriangle,
  CheckCircle,
  Target,
  Lightbulb,
  ArrowUp,
  ArrowDown
} from 'lucide-react'
import { useTransactionAnalytics } from '@/hooks/useTransactionAnalytics'
import { TransactionAnalytics, AnalyticsTimeframe, TransactionType } from '@/lib/transaction-analytics'
import { cn } from '@/lib/utils'

export function TransactionAnalyticsDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedTimeframe, setSelectedTimeframe] = useState(AnalyticsTimeframe.LAST_30_DAYS)

  const {
    state,
    applyFilters,
    exportToCSV,
    exportToJSON,
    refresh,
    clearError
  } = useTransactionAnalytics({
    autoLoad: true,
    enableNotifications: true,
    defaultTimeframe: selectedTimeframe
  })

  const handleTimeframeChange = (timeframe: AnalyticsTimeframe) => {
    setSelectedTimeframe(timeframe)
    const timeframeConfig = TransactionAnalytics.getTimeframeConfig(timeframe)
    applyFilters({
      timeframe: timeframeConfig,
      chains: [],
      types: [],
      status: []
    })
  }

  const getTimeframeName = (timeframe: AnalyticsTimeframe) => {
    const names = {
      [AnalyticsTimeframe.LAST_24_HOURS]: '24 Hours',
      [AnalyticsTimeframe.LAST_7_DAYS]: '7 Days',
      [AnalyticsTimeframe.LAST_30_DAYS]: '30 Days',
      [AnalyticsTimeframe.LAST_90_DAYS]: '90 Days',
      [AnalyticsTimeframe.LAST_YEAR]: '1 Year',
      [AnalyticsTimeframe.ALL_TIME]: 'All Time'
    }
    return names[timeframe] || '30 Days'
  }

  const getInsightIcon = (type: string) => {
    switch (type) {
      case 'gas_optimization':
        return <Zap className="w-4 h-4" />
      case 'timing':
        return <Clock className="w-4 h-4" />
      case 'chain_selection':
        return <Target className="w-4 h-4" />
      case 'cost_reduction':
        return <DollarSign className="w-4 h-4" />
      case 'security':
        return <Shield className="w-4 h-4" />
      default:
        return <Lightbulb className="w-4 h-4" />
    }
  }

  const getWarningIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
        return <AlertTriangle className="w-4 h-4 text-red-500" />
      case 'high':
        return <AlertTriangle className="w-4 h-4 text-orange-500" />
      case 'medium':
        return <AlertTriangle className="w-4 h-4 text-yellow-500" />
      default:
        return <AlertTriangle className="w-4 h-4 text-blue-500" />
    }
  }

  const getTrendIcon = (direction: string) => {
    switch (direction) {
      case 'increasing':
        return <TrendingUp className="w-4 h-4 text-green-500" />
      case 'decreasing':
        return <TrendingDown className="w-4 h-4 text-red-500" />
      default:
        return <Activity className="w-4 h-4 text-gray-500" />
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Transaction Analytics</h2>
          <p className="text-muted-foreground">
            Comprehensive analysis of your transaction history and performance
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => exportToCSV()}
          >
            <Download className="w-4 h-4 mr-2" />
            Export CSV
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => exportToJSON()}
          >
            <Download className="w-4 h-4 mr-2" />
            Export JSON
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={refresh}
            disabled={state.isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Timeframe Selector */}
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium">Timeframe:</span>
        {Object.values(AnalyticsTimeframe).map((timeframe) => (
          <Button
            key={timeframe}
            variant={selectedTimeframe === timeframe ? "default" : "outline"}
            size="sm"
            onClick={() => handleTimeframeChange(timeframe)}
          >
            {getTimeframeName(timeframe)}
          </Button>
        ))}
      </div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Transactions</p>
                <p className="text-2xl font-bold">{state.metrics.totalTransactions}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {state.filteredTransactions.length} in current period
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                <p className="text-2xl font-bold">{state.metrics.totalValue.toFixed(4)} ETH</p>
              </div>
              <DollarSign className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Across all transactions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {(state.metrics.successRate * 100).toFixed(1)}%
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Transaction success rate
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Gas Cost</p>
                <p className="text-2xl font-bold">{state.metrics.averageGasPrice.toFixed(2)} gwei</p>
              </div>
              <Zap className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Average gas price paid
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="insights">Insights</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="costs">Costs</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Transaction Distribution */}
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Transaction Types</CardTitle>
                <CardDescription>
                  Distribution of transaction types
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {Object.entries(TransactionType).map(([key, type]) => {
                    const count = state.filteredTransactions.filter(tx => tx.type === type).length
                    const percentage = state.filteredTransactions.length > 0 
                      ? (count / state.filteredTransactions.length) * 100 
                      : 0

                    return (
                      <div key={type} className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 rounded-full bg-blue-500" />
                          <span className="text-sm font-medium capitalize">{type.replace('_', ' ')}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-sm text-muted-foreground">{count}</span>
                          <span className="text-sm font-medium">{percentage.toFixed(1)}%</span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Performance Score</CardTitle>
                <CardDescription>
                  Overall transaction performance
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-4xl font-bold text-green-600">
                      {state.performanceAnalysis.performanceScore.toFixed(0)}
                    </div>
                    <div className="text-sm text-muted-foreground">Performance Score</div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Gas Efficiency</span>
                      <span>{state.performanceAnalysis.gasEfficiency.toFixed(1)}%</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Failure Rate</span>
                      <span>{state.performanceAnalysis.failureRate.toFixed(1)}%</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Avg Confirmation</span>
                      <span>{state.performanceAnalysis.averageConfirmationTime.toFixed(1)}s</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="insights" className="space-y-6">
          {/* Recommendations */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Lightbulb className="w-5 h-5" />
                Recommendations
              </CardTitle>
              <CardDescription>
                AI-powered suggestions to optimize your transactions
              </CardDescription>
            </CardHeader>
            <CardContent>
              {state.insights.recommendations.length > 0 ? (
                <div className="space-y-4">
                  {state.insights.recommendations.map((rec, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-start gap-3 p-4 border rounded-lg"
                    >
                      <div className="mt-1">
                        {getInsightIcon(rec.type)}
                      </div>
                      <div className="flex-1">
                        <h4 className="font-medium">{rec.title}</h4>
                        <p className="text-sm text-muted-foreground mt-1">
                          {rec.description}
                        </p>
                        <div className="flex items-center gap-2 mt-2">
                          <Badge variant={rec.impact === 'high' ? 'default' : 'secondary'}>
                            {rec.impact} impact
                          </Badge>
                          {rec.potentialSavings && (
                            <Badge variant="outline">
                              Save {rec.potentialSavings}
                            </Badge>
                          )}
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Lightbulb className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">No recommendations available</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Warnings */}
          {state.insights.warnings.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="w-5 h-5" />
                  Warnings
                </CardTitle>
                <CardDescription>
                  Issues that require your attention
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {state.insights.warnings.map((warning, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-start gap-3 p-4 border rounded-lg"
                    >
                      <div className="mt-1">
                        {getWarningIcon(warning.severity)}
                      </div>
                      <div className="flex-1">
                        <h4 className="font-medium">{warning.title}</h4>
                        <p className="text-sm text-muted-foreground mt-1">
                          {warning.description}
                        </p>
                        {warning.recommendation && (
                          <p className="text-sm text-blue-600 mt-2">
                            💡 {warning.recommendation}
                          </p>
                        )}
                        <Badge variant="destructive" className="mt-2">
                          {warning.severity} severity
                        </Badge>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Trends */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="w-5 h-5" />
                Trends
              </CardTitle>
              <CardDescription>
                Transaction patterns and trends
              </CardDescription>
            </CardHeader>
            <CardContent>
              {state.insights.trends.length > 0 ? (
                <div className="space-y-4">
                  {state.insights.trends.map((trend, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-4 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        {getTrendIcon(trend.direction)}
                        <div>
                          <h4 className="font-medium">{trend.metric}</h4>
                          <p className="text-sm text-muted-foreground">
                            {trend.description}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="flex items-center gap-1">
                          {trend.direction === 'increasing' ? (
                            <ArrowUp className="w-4 h-4 text-green-500" />
                          ) : trend.direction === 'decreasing' ? (
                            <ArrowDown className="w-4 h-4 text-red-500" />
                          ) : null}
                          <span className="font-medium">{trend.change.toFixed(1)}%</span>
                        </div>
                        <Badge variant="outline" className="text-xs">
                          {trend.significance} significance
                        </Badge>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <TrendingUp className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">No trends detected</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="performance" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Performance Analysis</CardTitle>
              <CardDescription>
                Detailed performance metrics and analysis
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-12">
                <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">Performance Charts Coming Soon</h3>
                <p className="text-muted-foreground">
                  Detailed performance charts and visualizations will be available here
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="costs" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Cost Analysis</CardTitle>
              <CardDescription>
                Transaction cost breakdown and optimization opportunities
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium mb-2">Total Costs</h4>
                    <p className="text-2xl font-bold">${state.costAnalysis.totalCosts.toFixed(2)}</p>
                  </div>
                  <div>
                    <h4 className="font-medium mb-2">Average per Transaction</h4>
                    <p className="text-xl font-semibold">${state.costAnalysis.averageCostPerTransaction.toFixed(2)}</p>
                  </div>
                  <div>
                    <h4 className="font-medium mb-2">Cost Efficiency Score</h4>
                    <p className="text-xl font-semibold">{state.costAnalysis.costEfficiencyScore.toFixed(0)}/100</p>
                  </div>
                </div>
                <div className="space-y-4">
                  <h4 className="font-medium">Costs by Chain</h4>
                  {Object.entries(state.costAnalysis.costsByChain).map(([chainId, cost]) => (
                    <div key={chainId} className="flex justify-between">
                      <span>Chain {chainId}</span>
                      <span>${cost.toFixed(2)}</span>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
