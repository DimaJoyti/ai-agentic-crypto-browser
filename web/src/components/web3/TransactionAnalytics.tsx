'use client'

import { useState, useEffect, useMemo } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  BarChart3, 
  TrendingUp, 
  TrendingDown,
  Activity, 
  DollarSign,
  Zap,
  Clock,
  Download,
  Filter,
  Calendar,
  PieChart,
  LineChart
} from 'lucide-react'
import { useTransactionMonitor } from '@/hooks/useTransactionMonitor'
import { 
  TransactionAnalytics as Analytics,
  type AnalyticsFilters,
  type TransactionMetrics,
  type ChainMetrics,
  type TypeMetrics,
  type TimeSeriesData
} from '@/lib/transaction-analytics'
import { TransactionStatus, TransactionType } from '@/lib/transaction-monitor'
import { SUPPORTED_CHAINS } from '@/lib/chains'
import { AnalyticsCharts } from './AnalyticsCharts'
import { AnalyticsFilters as FiltersComponent } from './AnalyticsFilters'

export function TransactionAnalytics() {
  const [activeTab, setActiveTab] = useState('overview')
  const [filters, setFilters] = useState<AnalyticsFilters>(Analytics.getDefaultFilters())
  const [showFilters, setShowFilters] = useState(false)

  const { transactions } = useTransactionMonitor()

  // Filter transactions based on current filters
  const filteredTransactions = useMemo(() => {
    return Analytics.filterTransactions(transactions, filters)
  }, [transactions, filters])

  // Calculate metrics
  const metrics = useMemo(() => {
    return Analytics.calculateMetrics(filteredTransactions)
  }, [filteredTransactions])

  const chainMetrics = useMemo(() => {
    return Analytics.calculateChainMetrics(filteredTransactions)
  }, [filteredTransactions])

  const typeMetrics = useMemo(() => {
    return Analytics.calculateTypeMetrics(filteredTransactions)
  }, [filteredTransactions])

  const timeSeriesData = useMemo(() => {
    return Analytics.generateTimeSeriesData(filteredTransactions, filters.timeframe)
  }, [filteredTransactions, filters.timeframe])

  const handleExportCSV = () => {
    Analytics.exportToCSV(filteredTransactions, `transaction-analytics-${Date.now()}.csv`)
  }

  const handleExportJSON = () => {
    Analytics.exportToJSON(filteredTransactions, `transaction-analytics-${Date.now()}.json`)
  }

  const getMetricChange = (current: number, previous: number): { value: number; isPositive: boolean } => {
    if (previous === 0) return { value: 0, isPositive: true }
    const change = ((current - previous) / previous) * 100
    return { value: Math.abs(change), isPositive: change >= 0 }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <BarChart3 className="w-6 h-6" />
            Transaction Analytics
          </h2>
          <p className="text-muted-foreground">
            Comprehensive insights into your blockchain transaction activity
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowFilters(!showFilters)}
            className={showFilters ? 'bg-blue-50 border-blue-200' : ''}
          >
            <Filter className="w-4 h-4 mr-2" />
            Filters
          </Button>
          <Button variant="outline" size="sm" onClick={handleExportCSV}>
            <Download className="w-4 h-4 mr-2" />
            Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={handleExportJSON}>
            <Download className="w-4 h-4 mr-2" />
            Export JSON
          </Button>
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          exit={{ opacity: 0, height: 0 }}
        >
          <FiltersComponent filters={filters} onFiltersChange={setFilters} />
        </motion.div>
      )}

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Transactions</p>
                <p className="text-2xl font-bold">{metrics.totalTransactions.toLocaleString()}</p>
                <div className="flex items-center gap-1 mt-1">
                  <Badge variant="secondary" className="text-xs">
                    {filters.timeframe.label}
                  </Badge>
                </div>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold">{metrics.successRate.toFixed(1)}%</p>
                <div className="flex items-center gap-1 mt-1">
                  <TrendingUp className="w-3 h-3 text-green-500" />
                  <span className="text-xs text-green-600">
                    {metrics.successfulTransactions} successful
                  </span>
                </div>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Volume</p>
                <p className="text-2xl font-bold">{Analytics.formatCurrency(metrics.totalVolume)}</p>
                <div className="flex items-center gap-1 mt-1">
                  <DollarSign className="w-3 h-3 text-blue-500" />
                  <span className="text-xs text-muted-foreground">ETH</span>
                </div>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Gas Used</p>
                <p className="text-2xl font-bold">{Analytics.formatGas(metrics.averageGasUsed)}</p>
                <div className="flex items-center gap-1 mt-1">
                  <Zap className="w-3 h-3 text-yellow-500" />
                  <span className="text-xs text-muted-foreground">
                    {Analytics.formatTime(metrics.averageConfirmationTime)}
                  </span>
                </div>
              </div>
              <Zap className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Analytics Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="chains">By Chain</TabsTrigger>
          <TabsTrigger value="types">By Type</TabsTrigger>
          <TabsTrigger value="trends">Trends</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Transaction Status Distribution */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Transaction Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <AnalyticsCharts.StatusPieChart data={[
                  { name: 'Confirmed', value: metrics.successfulTransactions, color: '#10b981' },
                  { name: 'Failed', value: metrics.failedTransactions, color: '#ef4444' },
                  { name: 'Pending', value: metrics.pendingTransactions, color: '#f59e0b' }
                ]} />
              </CardContent>
            </Card>

            {/* Top Chains */}
            <Card>
              <CardHeader>
                <CardTitle>Top Chains</CardTitle>
                <CardDescription>Most active blockchain networks</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {chainMetrics.slice(0, 5).map((chain, index) => (
                    <div key={chain.chainId} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center text-sm font-medium">
                          {index + 1}
                        </div>
                        <div>
                          <p className="font-medium">{chain.chainName}</p>
                          <p className="text-sm text-muted-foreground">
                            {chain.transactionCount} transactions
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{Analytics.formatCurrency(chain.volume)}</p>
                        <p className="text-sm text-muted-foreground">
                          {chain.successRate.toFixed(1)}% success
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Transaction Volume Over Time */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <LineChart className="w-5 h-5" />
                Transaction Activity
              </CardTitle>
              <CardDescription>
                Transaction volume and count over {filters.timeframe.label.toLowerCase()}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AnalyticsCharts.TimeSeriesChart data={timeSeriesData} />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="chains" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Chain Performance</CardTitle>
              <CardDescription>
                Detailed metrics for each blockchain network
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AnalyticsCharts.ChainBarChart data={chainMetrics} />
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Chain Distribution</CardTitle>
              </CardHeader>
              <CardContent>
                <AnalyticsCharts.ChainPieChart data={chainMetrics} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Chain Details</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {chainMetrics.map(chain => (
                    <div key={chain.chainId} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <h4 className="font-medium">{chain.chainName}</h4>
                        <Badge variant="outline">
                          {chain.transactionCount} txs
                        </Badge>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <p className="text-muted-foreground">Volume</p>
                          <p className="font-medium">{Analytics.formatCurrency(chain.volume)} ETH</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Success Rate</p>
                          <p className="font-medium">{chain.successRate.toFixed(1)}%</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Avg Gas</p>
                          <p className="font-medium">{Analytics.formatGas(chain.gasUsed)}</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Avg Time</p>
                          <p className="font-medium">{Analytics.formatTime(chain.averageConfirmationTime)}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="types" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Transaction Types</CardTitle>
              <CardDescription>
                Performance metrics by transaction type
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AnalyticsCharts.TypeBarChart data={typeMetrics} />
            </CardContent>
          </Card>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Type Distribution</CardTitle>
              </CardHeader>
              <CardContent>
                <AnalyticsCharts.TypePieChart data={typeMetrics} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Type Performance</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {typeMetrics.map(type => (
                    <div key={type.type} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <h4 className="font-medium capitalize">
                          {type.type.replace('_', ' ')}
                        </h4>
                        <Badge variant="outline">
                          {type.count} txs
                        </Badge>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <p className="text-muted-foreground">Volume</p>
                          <p className="font-medium">{Analytics.formatCurrency(type.volume)} ETH</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Success Rate</p>
                          <p className="font-medium">{type.successRate.toFixed(1)}%</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Avg Gas</p>
                          <p className="font-medium">{Analytics.formatGas(type.averageGasUsed)}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="trends" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Success Rate Trend</CardTitle>
              </CardHeader>
              <CardContent>
                <AnalyticsCharts.SuccessRateTrend data={timeSeriesData} />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Gas Usage Trend</CardTitle>
              </CardHeader>
              <CardContent>
                <AnalyticsCharts.GasUsageTrend data={timeSeriesData} />
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Volume Trend</CardTitle>
              <CardDescription>
                Transaction volume over {filters.timeframe.label.toLowerCase()}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AnalyticsCharts.VolumeTrend data={timeSeriesData} />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
