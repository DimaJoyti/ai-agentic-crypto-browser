'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { 
  TrendingUp,
  TrendingDown,
  DollarSign,
  PieChart,
  Shield,
  AlertTriangle,
  Target,
  Activity,
  RefreshCw,
  Plus,
  Eye,
  EyeOff
} from 'lucide-react'
import { usePortfolioAnalytics } from '@/hooks/usePortfolioAnalytics'
import { cn } from '@/lib/utils'

interface PortfolioOverviewProps {
  showValues?: boolean
  compact?: boolean
}

export function PortfolioOverview({ showValues = true, compact = false }: PortfolioOverviewProps) {
  const [hideValues, setHideValues] = useState(!showValues)

  const {
    state,
    getMetrics,
    getAllocation,
    getRiskMetrics,
    getTopGainers,
    getTopLosers,
    refreshPortfolio
  } = usePortfolioAnalytics({
    autoSync: true,
    trackPriceUpdates: true
  })

  const metrics = getMetrics()
  const allocation = getAllocation()
  const riskMetrics = getRiskMetrics()
  const topGainers = getTopGainers()
  const topLosers = getTopLosers()

  const formatCurrency = (value: number) => {
    if (hideValues) return '****'
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }

  const formatPercent = (value: number) => {
    if (hideValues) return '**%'
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

  if (state.isLoading) {
    return (
      <Card className="w-full">
        <CardContent className="p-6">
          <div className="flex items-center justify-center space-x-2">
            <RefreshCw className="w-4 h-4 animate-spin" />
            <span className="text-sm text-muted-foreground">Loading portfolio...</span>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (!metrics || state.positions.length === 0) {
    return (
      <Card className="w-full">
        <CardContent className="p-6">
          <div className="text-center">
            <PieChart className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <h3 className="text-lg font-medium mb-2">No Portfolio Data</h3>
            <p className="text-muted-foreground mb-4">
              Start by adding your first position to track your portfolio performance
            </p>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Add Position
            </Button>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (compact) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                <p className="text-xl font-bold">{formatCurrency(metrics.totalValue)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
            <div className={cn("mt-2 text-sm", getChangeColor(metrics.dayChangePercent))}>
              {formatPercent(metrics.dayChangePercent)} today
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total P&L</p>
                <p className={cn("text-xl font-bold", getChangeColor(metrics.totalPnL))}>
                  {formatCurrency(metrics.totalPnL)}
                </p>
              </div>
              {metrics.totalPnL >= 0 ? (
                <TrendingUp className="w-8 h-8 text-green-500" />
              ) : (
                <TrendingDown className="w-8 h-8 text-red-500" />
              )}
            </div>
            <div className={cn("mt-2 text-sm", getChangeColor(metrics.totalPnLPercent))}>
              {formatPercent(metrics.totalPnLPercent)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Positions</p>
                <p className="text-xl font-bold">{state.positions.length}</p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {allocation.length > 0 ? `Top: ${allocation[0].symbol}` : 'No positions'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Risk Level</p>
                <p className={cn("text-xl font-bold", getRiskColor(riskMetrics?.portfolioRisk || 'low'))}>
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
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Portfolio Overview</h2>
          <p className="text-muted-foreground">
            Track your cryptocurrency portfolio performance and analytics
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setHideValues(!hideValues)}
          >
            {hideValues ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
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

      {/* Main Metrics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Portfolio Value</p>
                <p className="text-2xl font-bold">{formatCurrency(metrics.totalValue)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
            <div className={cn("mt-2 text-sm flex items-center", getChangeColor(metrics.dayChangePercent))}>
              {metrics.dayChangePercent >= 0 ? (
                <TrendingUp className="w-4 h-4 mr-1" />
              ) : (
                <TrendingDown className="w-4 h-4 mr-1" />
              )}
              {formatCurrency(metrics.dayChange)} ({formatPercent(metrics.dayChangePercent)}) today
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total P&L</p>
                <p className={cn("text-2xl font-bold", getChangeColor(metrics.totalPnL))}>
                  {formatCurrency(metrics.totalPnL)}
                </p>
              </div>
              {metrics.totalPnL >= 0 ? (
                <TrendingUp className="w-8 h-8 text-green-500" />
              ) : (
                <TrendingDown className="w-8 h-8 text-red-500" />
              )}
            </div>
            <div className={cn("mt-2 text-sm", getChangeColor(metrics.totalPnLPercent))}>
              {formatPercent(metrics.totalPnLPercent)} return
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Cost</p>
                <p className="text-2xl font-bold">{formatCurrency(metrics.totalCost)}</p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {state.positions.length} position{state.positions.length !== 1 ? 's' : ''}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Risk Level</p>
                <p className={cn("text-2xl font-bold", getRiskColor(riskMetrics?.portfolioRisk || 'low'))}>
                  {riskMetrics?.portfolioRisk.replace('_', ' ').toUpperCase() || 'LOW'}
                </p>
              </div>
              <Shield className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Diversification: {riskMetrics?.diversificationScore.toFixed(0) || '0'}%
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Performance Metrics */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Performance Metrics</CardTitle>
            <CardDescription>Key portfolio performance indicators</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-muted-foreground">24h Change</p>
                  <p className={cn("text-lg font-semibold", getChangeColor(metrics.dayChangePercent))}>
                    {formatPercent(metrics.dayChangePercent)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">7d Change</p>
                  <p className={cn("text-lg font-semibold", getChangeColor(metrics.weekChangePercent))}>
                    {formatPercent(metrics.weekChangePercent)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">30d Change</p>
                  <p className={cn("text-lg font-semibold", getChangeColor(metrics.monthChangePercent))}>
                    {formatPercent(metrics.monthChangePercent)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Max Drawdown</p>
                  <p className="text-lg font-semibold text-red-600">
                    -{formatPercent(metrics.maxDrawdown)}
                  </p>
                </div>
              </div>

              <div className="pt-4 border-t">
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">All-Time High</span>
                    <span>{formatCurrency(metrics.allTimeHigh)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">All-Time Low</span>
                    <span>{formatCurrency(metrics.allTimeLow)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Sharpe Ratio</span>
                    <span>{metrics.sharpeRatio.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Volatility</span>
                    <span>{(metrics.volatility * 100).toFixed(1)}%</span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Risk Analysis</CardTitle>
            <CardDescription>Portfolio risk assessment and recommendations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Overall Risk</span>
                <Badge variant={
                  riskMetrics?.portfolioRisk === 'low' ? 'default' :
                  riskMetrics?.portfolioRisk === 'medium' ? 'secondary' :
                  riskMetrics?.portfolioRisk === 'high' ? 'destructive' :
                  'destructive'
                }>
                  {riskMetrics?.portfolioRisk.toUpperCase() || 'LOW'}
                </Badge>
              </div>

              <div className="space-y-3">
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Diversification</span>
                    <span>{riskMetrics?.diversificationScore.toFixed(0) || '0'}%</span>
                  </div>
                  <Progress value={riskMetrics?.diversificationScore || 0} className="h-2" />
                </div>

                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Concentration Risk</span>
                    <span>{riskMetrics?.concentrationRisk.toFixed(0) || '0'}%</span>
                  </div>
                  <Progress value={riskMetrics?.concentrationRisk || 0} className="h-2" />
                </div>

                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Volatility Risk</span>
                    <span>{riskMetrics?.volatilityRisk.toFixed(0) || '0'}%</span>
                  </div>
                  <Progress value={riskMetrics?.volatilityRisk || 0} className="h-2" />
                </div>
              </div>

              {riskMetrics?.recommendations && riskMetrics.recommendations.length > 0 && (
                <div className="pt-4 border-t">
                  <h4 className="text-sm font-medium mb-2 flex items-center">
                    <AlertTriangle className="w-4 h-4 mr-1 text-yellow-500" />
                    Recommendations
                  </h4>
                  <ul className="text-xs text-muted-foreground space-y-1">
                    {riskMetrics.recommendations.slice(0, 3).map((rec, index) => (
                      <li key={index}>• {rec}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Performers */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-green-500" />
              Top Gainers
            </CardTitle>
            <CardDescription>Best performing positions</CardDescription>
          </CardHeader>
          <CardContent>
            {topGainers.length > 0 ? (
              <div className="space-y-3">
                {topGainers.slice(0, 5).map((position, index) => (
                  <motion.div
                    key={position.symbol}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
                        <span className="text-sm font-bold text-green-600 dark:text-green-400">
                          {index + 1}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium">{position.symbol}</p>
                        <p className="text-sm text-muted-foreground">
                          {position.amount.toFixed(4)} @ {formatCurrency(position.averageCost)}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-green-600 dark:text-green-400 font-medium">
                        +{formatPercent(position.unrealizedPnLPercent)}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {formatCurrency(position.unrealizedPnL)}
                      </p>
                    </div>
                  </motion.div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No gainers in current portfolio
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingDown className="w-5 h-5 text-red-500" />
              Top Losers
            </CardTitle>
            <CardDescription>Underperforming positions</CardDescription>
          </CardHeader>
          <CardContent>
            {topLosers.length > 0 ? (
              <div className="space-y-3">
                {topLosers.slice(0, 5).map((position, index) => (
                  <motion.div
                    key={position.symbol}
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-red-100 dark:bg-red-900 rounded-full flex items-center justify-center">
                        <span className="text-sm font-bold text-red-600 dark:text-red-400">
                          {index + 1}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium">{position.symbol}</p>
                        <p className="text-sm text-muted-foreground">
                          {position.amount.toFixed(4)} @ {formatCurrency(position.averageCost)}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-red-600 dark:text-red-400 font-medium">
                        {formatPercent(position.unrealizedPnLPercent)}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {formatCurrency(position.unrealizedPnL)}
                      </p>
                    </div>
                  </motion.div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No losers in current portfolio
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
