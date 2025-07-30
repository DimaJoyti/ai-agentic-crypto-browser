'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  PieChart,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Target,
  Shield,
  Activity,
  BarChart3,
  RefreshCw,
  Download,
  Upload,
  Settings,
  AlertCircle,
  Eye,
  EyeOff
} from 'lucide-react'
import { PortfolioOverview } from './PortfolioOverview'
import { PortfolioPositions } from './PortfolioPositions'
import { usePortfolioAnalytics } from '@/hooks/usePortfolioAnalytics'
import { cn } from '@/lib/utils'

export function PortfolioAnalyticsDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [showValues, setShowValues] = useState(true)

  const {
    state,
    getAllocation,
    getMetrics,
    getRiskMetrics,
    exportPortfolio,
    refreshPortfolio,
    clearError
  } = usePortfolioAnalytics({
    autoSync: true,
    trackPriceUpdates: true
  })

  const allocation = getAllocation()
  const metrics = getMetrics()
  const riskMetrics = getRiskMetrics()

  const formatCurrency = (value: number) => {
    if (!showValues) return '****'
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }

  const formatPercent = (value: number) => {
    if (!showValues) return '**%'
    const sign = value >= 0 ? '+' : ''
    return `${sign}${value.toFixed(2)}%`
  }

  const getChangeColor = (value: number) => {
    if (value > 0) return 'text-green-600 dark:text-green-400'
    if (value < 0) return 'text-red-600 dark:text-red-400'
    return 'text-gray-600 dark:text-gray-400'
  }

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'low':
        return 'text-green-600 dark:text-green-400'
      case 'medium':
        return 'text-yellow-600 dark:text-yellow-400'
      case 'high':
        return 'text-orange-600 dark:text-orange-400'
      case 'extreme':
        return 'text-red-600 dark:text-red-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const handleExport = () => {
    const data = exportPortfolio()
    if (data) {
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `portfolio-${new Date().toISOString().split('T')[0]}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Portfolio Analytics</h1>
          <p className="text-muted-foreground">
            Comprehensive portfolio tracking and performance analysis
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowValues(!showValues)}
          >
            {showValues ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleExport}
            disabled={state.positions.length === 0}
          >
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={refreshPortfolio}
            disabled={state.isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                <p className="text-2xl font-bold">{formatCurrency(metrics?.totalValue || 0)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
            <div className={cn("mt-2 text-sm", getChangeColor(metrics?.dayChangePercent || 0))}>
              {formatPercent(metrics?.dayChangePercent || 0)} today
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total P&L</p>
                <p className={cn("text-2xl font-bold", getChangeColor(metrics?.totalPnL || 0))}>
                  {formatCurrency(metrics?.totalPnL || 0)}
                </p>
              </div>
              {(metrics?.totalPnL || 0) >= 0 ? (
                <TrendingUp className="w-8 h-8 text-green-500" />
              ) : (
                <TrendingDown className="w-8 h-8 text-red-500" />
              )}
            </div>
            <div className={cn("mt-2 text-sm", getChangeColor(metrics?.totalPnLPercent || 0))}>
              {formatPercent(metrics?.totalPnLPercent || 0)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Positions</p>
                <p className="text-2xl font-bold">{state.positions.length}</p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {allocation.length > 0 ? `Top: ${allocation[0].symbol}` : 'No positions'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Risk Level</p>
                <p className={cn("text-2xl font-bold", getRiskColor(riskMetrics?.portfolioRisk || 'low'))}>
                  {riskMetrics?.portfolioRisk.toUpperCase() || 'LOW'}
                </p>
              </div>
              <Shield className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Score: {riskMetrics?.riskScore.toFixed(0) || '0'}/100
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Asset Allocation */}
      {allocation.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <PieChart className="w-5 h-5" />
              Asset Allocation
            </CardTitle>
            <CardDescription>
              Portfolio distribution by asset value
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {allocation.slice(0, 6).map((asset, index) => (
                <motion.div
                  key={asset.symbol}
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center justify-between p-3 border rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    <div 
                      className="w-4 h-4 rounded-full"
                      style={{ backgroundColor: asset.color }}
                    />
                    <div>
                      <p className="font-medium">{asset.symbol}</p>
                      <p className="text-sm text-muted-foreground">
                        {formatCurrency(asset.value)}
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline">
                    {asset.percentage.toFixed(1)}%
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="positions">Positions</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <PortfolioOverview showValues={showValues} />
        </TabsContent>

        <TabsContent value="positions" className="space-y-6">
          <PortfolioPositions showValues={showValues} />
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            {/* Performance Chart Placeholder */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Performance Chart
                </CardTitle>
                <CardDescription>Portfolio value over time</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-64 flex items-center justify-center border-2 border-dashed border-muted rounded-lg">
                  <div className="text-center">
                    <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-2" />
                    <p className="text-muted-foreground">Performance chart coming soon</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Risk Analysis */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Risk Analysis
                </CardTitle>
                <CardDescription>Portfolio risk assessment</CardDescription>
              </CardHeader>
              <CardContent>
                {riskMetrics ? (
                  <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div className="text-center p-3 border rounded">
                        <p className="text-lg font-bold">{riskMetrics.diversificationScore.toFixed(0)}%</p>
                        <p className="text-sm text-muted-foreground">Diversification</p>
                      </div>
                      <div className="text-center p-3 border rounded">
                        <p className="text-lg font-bold">{riskMetrics.concentrationRisk.toFixed(0)}%</p>
                        <p className="text-sm text-muted-foreground">Concentration</p>
                      </div>
                      <div className="text-center p-3 border rounded">
                        <p className="text-lg font-bold">{riskMetrics.volatilityRisk.toFixed(0)}%</p>
                        <p className="text-sm text-muted-foreground">Volatility</p>
                      </div>
                      <div className="text-center p-3 border rounded">
                        <p className="text-lg font-bold">{riskMetrics.liquidityRisk.toFixed(0)}%</p>
                        <p className="text-sm text-muted-foreground">Liquidity</p>
                      </div>
                    </div>

                    {riskMetrics.recommendations.length > 0 && (
                      <div className="pt-4 border-t">
                        <h4 className="font-medium mb-2">Recommendations</h4>
                        <ul className="text-sm text-muted-foreground space-y-1">
                          {riskMetrics.recommendations.slice(0, 3).map((rec, index) => (
                            <li key={index}>• {rec}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-8 text-muted-foreground">
                    No risk data available
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="transactions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="w-5 h-5" />
                Transaction History
              </CardTitle>
              <CardDescription>
                Recent portfolio transactions and trades
              </CardDescription>
            </CardHeader>
            <CardContent>
              {state.transactions.length > 0 ? (
                <div className="space-y-3">
                  {state.transactions.slice(0, 10).map((transaction, index) => (
                    <motion.div
                      key={transaction.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.05 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className={cn("w-8 h-8 rounded-full flex items-center justify-center",
                          transaction.type === 'buy' ? 'bg-green-100 dark:bg-green-900' :
                          transaction.type === 'sell' ? 'bg-red-100 dark:bg-red-900' :
                          'bg-blue-100 dark:bg-blue-900'
                        )}>
                          <span className={cn("text-xs font-bold",
                            transaction.type === 'buy' ? 'text-green-600 dark:text-green-400' :
                            transaction.type === 'sell' ? 'text-red-600 dark:text-red-400' :
                            'text-blue-600 dark:text-blue-400'
                          )}>
                            {transaction.type === 'buy' ? 'B' : transaction.type === 'sell' ? 'S' : 'T'}
                          </span>
                        </div>
                        <div>
                          <p className="font-medium">
                            {transaction.type.toUpperCase()} {transaction.symbol}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            {new Date(transaction.timestamp).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">
                          {transaction.amount.toFixed(8)} @ {formatCurrency(transaction.price)}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {formatCurrency(transaction.value)}
                        </p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  <Activity className="w-12 h-12 mx-auto mb-4" />
                  <p>No transactions yet</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
