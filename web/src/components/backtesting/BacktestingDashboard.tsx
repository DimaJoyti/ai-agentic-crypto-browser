'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Play,
  Pause,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Target,
  Shield,
  Clock,
  DollarSign,
  Activity,
  Settings,
  Download,
  Upload,
  Plus,
  Trash2,
  AlertCircle,
  CheckCircle
} from 'lucide-react'
import { useHistoricalData } from '@/hooks/useHistoricalData'
import { useTradingSignals } from '@/hooks/useTradingSignals'
import { type BacktestResults } from '@/lib/backtesting-engine'
import { cn } from '@/lib/utils'

export function BacktestingDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedBacktest, setSelectedBacktest] = useState<BacktestResults | null>(null)

  const {
    dataState,
    backtestState,
    runBacktest,
    getBacktests,
    deleteBacktest,
    clearError
  } = useHistoricalData()

  const { getStrategies } = useTradingSignals()

  const backtests = getBacktests()
  const strategies = getStrategies()

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }

  const formatPercent = (value: number) => {
    const sign = value >= 0 ? '+' : ''
    return `${sign}${value.toFixed(2)}%`
  }

  const getPerformanceColor = (value: number) => {
    if (value > 0) return 'text-green-600 dark:text-green-400'
    if (value < 0) return 'text-red-600 dark:text-red-400'
    return 'text-gray-600 dark:text-gray-400'
  }

  const formatDuration = (ms: number) => {
    const days = Math.floor(ms / (1000 * 60 * 60 * 24))
    const hours = Math.floor((ms % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    
    if (days > 0) return `${days}d ${hours}h`
    if (hours > 0) return `${hours}h`
    return `${Math.floor(ms / (1000 * 60))}m`
  }

  const handleRunBacktest = async () => {
    // This would open a dialog to configure backtest parameters
    // For now, we'll run a simple example backtest
    const config = {
      id: `backtest_${Date.now()}`,
      name: 'Example Backtest',
      strategy: strategies[0] || {
        id: 'example',
        name: 'Example Strategy',
        description: 'Example RSI strategy',
        enabled: true,
        symbols: ['BTC'],
        timeframes: ['1h'],
        conditions: [],
        filters: [],
        riskLevel: 'medium' as const,
        minConfidence: 70,
        maxSignalsPerDay: 5,
        cooldownPeriod: 60,
        createdAt: Date.now()
      },
      symbol: 'BTC',
      timeframe: '1h',
      startTime: Date.now() - 30 * 24 * 60 * 60 * 1000, // 30 days ago
      endTime: Date.now(),
      initialCapital: 10000,
      positionSize: 10, // 10% per trade
      maxPositions: 3,
      commission: 0.1, // 0.1%
      slippage: 0.05, // 0.05%
      stopLoss: 5, // 5%
      takeProfit: 10, // 10%
      riskManagement: {
        maxDrawdown: 20,
        maxDailyLoss: 5,
        maxConsecutiveLosses: 3,
        positionSizing: 'percentage' as const,
        riskPerTrade: 2
      }
    }

    try {
      await runBacktest(config)
    } catch (error) {
      console.error('Backtest failed:', error)
    }
  }

  const handleDeleteBacktest = (id: string) => {
    if (window.confirm('Are you sure you want to delete this backtest?')) {
      deleteBacktest(id)
      if (selectedBacktest?.config.id === id) {
        setSelectedBacktest(null)
      }
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Backtesting</h1>
          <p className="text-muted-foreground">
            Test your trading strategies against historical data
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={backtestState.isRunning}
          >
            <Upload className="w-4 h-4 mr-2" />
            Import
          </Button>
          <Button
            onClick={handleRunBacktest}
            disabled={backtestState.isRunning || strategies.length === 0}
          >
            {backtestState.isRunning ? (
              <>
                <Pause className="w-4 h-4 mr-2" />
                Running...
              </>
            ) : (
              <>
                <Play className="w-4 h-4 mr-2" />
                Run Backtest
              </>
            )}
          </Button>
        </div>
      </div>

      {/* Running Backtest Progress */}
      {backtestState.isRunning && (
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="font-medium">Running Backtest</h3>
                <p className="text-sm text-muted-foreground">
                  Testing strategy against historical data...
                </p>
              </div>
              <div className="text-right">
                <p className="text-2xl font-bold">{backtestState.progress.toFixed(0)}%</p>
                <p className="text-sm text-muted-foreground">Complete</p>
              </div>
            </div>
            <Progress value={backtestState.progress} className="h-2" />
          </CardContent>
        </Card>
      )}

      {/* Error Alert */}
      {(dataState.error || backtestState.error) && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            {dataState.error || backtestState.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Backtests</p>
                <p className="text-2xl font-bold">{backtests.length}</p>
              </div>
              <BarChart3 className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Historical tests run
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Best Return</p>
                <p className="text-2xl font-bold text-green-600">
                  {backtests.length > 0 
                    ? formatPercent(Math.max(...backtests.map(b => b.performance.totalReturnPercent)))
                    : '0%'
                  }
                </p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Top performing strategy
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Win Rate</p>
                <p className="text-2xl font-bold">
                  {backtests.length > 0 
                    ? (backtests.reduce((sum, b) => sum + b.performance.winRate, 0) / backtests.length).toFixed(1) + '%'
                    : '0%'
                  }
                </p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Average success rate
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Data Points</p>
                <p className="text-2xl font-bold">{dataState.stats.totalCandles.toLocaleString()}</p>
              </div>
              <Activity className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Historical candles
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="results">Results</TabsTrigger>
          <TabsTrigger value="analysis">Analysis</TabsTrigger>
          <TabsTrigger value="data">Data</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Recent Backtests */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Backtests</CardTitle>
              <CardDescription>
                Latest strategy backtesting results
              </CardDescription>
            </CardHeader>
            <CardContent>
              {backtests.length > 0 ? (
                <div className="space-y-4">
                  {backtests.slice(0, 5).map((backtest, index) => (
                    <motion.div
                      key={backtest.config.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className={cn(
                        "flex items-center justify-between p-4 border rounded-lg cursor-pointer hover:bg-muted/50",
                        selectedBacktest?.config.id === backtest.config.id && "ring-2 ring-primary"
                      )}
                      onClick={() => setSelectedBacktest(backtest)}
                    >
                      <div className="flex items-center gap-4">
                        <div className="w-12 h-12 bg-secondary rounded-lg flex items-center justify-center">
                          <BarChart3 className="w-6 h-6" />
                        </div>
                        <div>
                          <h4 className="font-medium">{backtest.config.name}</h4>
                          <p className="text-sm text-muted-foreground">
                            {backtest.config.symbol} • {backtest.config.timeframe} • 
                            {backtest.performance.totalTrades} trades
                          </p>
                          <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
                            <span>Duration: {Math.floor(backtest.statistics.duration)}d</span>
                            <span>Win Rate: {backtest.performance.winRate.toFixed(1)}%</span>
                            <span>Max DD: {backtest.performance.maxDrawdownPercent.toFixed(1)}%</span>
                          </div>
                        </div>
                      </div>

                      <div className="text-right">
                        <div className={cn("text-lg font-bold", getPerformanceColor(backtest.performance.totalReturnPercent))}>
                          {formatPercent(backtest.performance.totalReturnPercent)}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {formatCurrency(backtest.performance.totalReturn)}
                        </div>
                        <div className="flex items-center gap-1 mt-2">
                          <Badge variant={backtest.performance.sharpeRatio > 1 ? 'default' : 'secondary'}>
                            Sharpe: {backtest.performance.sharpeRatio.toFixed(2)}
                          </Badge>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDeleteBacktest(backtest.config.id)
                            }}
                          >
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Backtests Yet</h3>
                  <p className="text-muted-foreground mb-4">
                    Run your first backtest to analyze strategy performance
                  </p>
                  <Button onClick={handleRunBacktest} disabled={strategies.length === 0}>
                    <Play className="w-4 h-4 mr-2" />
                    Run First Backtest
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Strategy Performance Summary */}
          {strategies.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Strategy Performance</CardTitle>
                <CardDescription>
                  Performance summary across all strategies
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-3">
                  {strategies.slice(0, 3).map((strategy, index) => {
                    const strategyBacktests = backtests.filter(b => b.config.strategy.id === strategy.id)
                    const avgReturn = strategyBacktests.length > 0
                      ? strategyBacktests.reduce((sum, b) => sum + b.performance.totalReturnPercent, 0) / strategyBacktests.length
                      : 0

                    return (
                      <motion.div
                        key={strategy.id}
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: index * 0.1 }}
                        className="p-4 border rounded-lg"
                      >
                        <div className="flex items-center justify-between mb-2">
                          <h4 className="font-medium">{strategy.name}</h4>
                          <Badge variant={strategy.enabled ? 'default' : 'secondary'}>
                            {strategy.enabled ? 'Active' : 'Inactive'}
                          </Badge>
                        </div>
                        <p className="text-sm text-muted-foreground mb-3">
                          {strategy.description}
                        </p>
                        <div className="grid grid-cols-2 gap-2 text-sm">
                          <div>
                            <span className="text-muted-foreground">Backtests:</span>
                            <p className="font-medium">{strategyBacktests.length}</p>
                          </div>
                          <div>
                            <span className="text-muted-foreground">Avg Return:</span>
                            <p className={cn("font-medium", getPerformanceColor(avgReturn))}>
                              {formatPercent(avgReturn)}
                            </p>
                          </div>
                        </div>
                      </motion.div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="results" className="space-y-6">
          {selectedBacktest ? (
            <div className="space-y-6">
              {/* Backtest Header */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>{selectedBacktest.config.name}</CardTitle>
                      <CardDescription>
                        {selectedBacktest.config.symbol} • {selectedBacktest.config.timeframe} • 
                        {new Date(selectedBacktest.statistics.startDate).toLocaleDateString()} - 
                        {new Date(selectedBacktest.statistics.endDate).toLocaleDateString()}
                      </CardDescription>
                    </div>
                    <Button variant="outline" size="sm">
                      <Download className="w-4 h-4 mr-2" />
                      Export
                    </Button>
                  </div>
                </CardHeader>
              </Card>

              {/* Performance Metrics */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground">Total Return</p>
                        <p className={cn("text-2xl font-bold", getPerformanceColor(selectedBacktest.performance.totalReturnPercent))}>
                          {formatPercent(selectedBacktest.performance.totalReturnPercent)}
                        </p>
                      </div>
                      <DollarSign className="w-8 h-8 text-blue-500" />
                    </div>
                    <div className="mt-2 text-sm text-muted-foreground">
                      {formatCurrency(selectedBacktest.performance.totalReturn)}
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground">Win Rate</p>
                        <p className="text-2xl font-bold text-green-600">
                          {selectedBacktest.performance.winRate.toFixed(1)}%
                        </p>
                      </div>
                      <Target className="w-8 h-8 text-green-500" />
                    </div>
                    <div className="mt-2 text-sm text-muted-foreground">
                      {selectedBacktest.performance.winningTrades}/{selectedBacktest.performance.totalTrades} trades
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground">Sharpe Ratio</p>
                        <p className="text-2xl font-bold">
                          {selectedBacktest.performance.sharpeRatio.toFixed(2)}
                        </p>
                      </div>
                      <BarChart3 className="w-8 h-8 text-purple-500" />
                    </div>
                    <div className="mt-2 text-sm text-muted-foreground">
                      Risk-adjusted return
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground">Max Drawdown</p>
                        <p className="text-2xl font-bold text-red-600">
                          {selectedBacktest.performance.maxDrawdownPercent.toFixed(1)}%
                        </p>
                      </div>
                      <TrendingDown className="w-8 h-8 text-red-500" />
                    </div>
                    <div className="mt-2 text-sm text-muted-foreground">
                      {formatCurrency(selectedBacktest.performance.maxDrawdown)}
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Detailed Performance */}
              <div className="grid gap-6 md:grid-cols-2">
                <Card>
                  <CardHeader>
                    <CardTitle>Trade Statistics</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <p className="text-sm text-muted-foreground">Total Trades</p>
                          <p className="text-lg font-semibold">{selectedBacktest.performance.totalTrades}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Profit Factor</p>
                          <p className="text-lg font-semibold">{selectedBacktest.performance.profitFactor.toFixed(2)}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Avg Win</p>
                          <p className="text-lg font-semibold text-green-600">
                            {formatCurrency(selectedBacktest.performance.avgWin)}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Avg Loss</p>
                          <p className="text-lg font-semibold text-red-600">
                            {formatCurrency(selectedBacktest.performance.avgLoss)}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Largest Win</p>
                          <p className="text-lg font-semibold text-green-600">
                            {formatCurrency(selectedBacktest.performance.largestWin)}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Largest Loss</p>
                          <p className="text-lg font-semibold text-red-600">
                            {formatCurrency(selectedBacktest.performance.largestLoss)}
                          </p>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Risk Metrics</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <p className="text-sm text-muted-foreground">Volatility</p>
                          <p className="text-lg font-semibold">{selectedBacktest.performance.volatility.toFixed(2)}%</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Calmar Ratio</p>
                          <p className="text-lg font-semibold">{selectedBacktest.performance.calmarRatio.toFixed(2)}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">VaR 95%</p>
                          <p className="text-lg font-semibold">{selectedBacktest.riskMetrics.var95.toFixed(2)}%</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">VaR 99%</p>
                          <p className="text-lg font-semibold">{selectedBacktest.riskMetrics.var99.toFixed(2)}%</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Avg Hold Time</p>
                          <p className="text-lg font-semibold">
                            {formatDuration(selectedBacktest.performance.avgHoldingPeriod)}
                          </p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Total Fees</p>
                          <p className="text-lg font-semibold">
                            {formatCurrency(selectedBacktest.performance.totalCommission + selectedBacktest.performance.totalSlippage)}
                          </p>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Backtest Selected</h3>
                <p className="text-muted-foreground">
                  Select a backtest from the overview to view detailed results
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="analysis" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Performance Analysis</CardTitle>
              <CardDescription>
                Detailed analysis and optimization insights
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-12">
                <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">Analysis Coming Soon</h3>
                <p className="text-muted-foreground">
                  Advanced performance analysis and optimization tools will be available here
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="data" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Historical Data</CardTitle>
              <CardDescription>
                Manage historical data for backtesting
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">{dataState.stats.totalCandles.toLocaleString()}</p>
                  <p className="text-sm text-muted-foreground">Total Candles</p>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">{dataState.stats.symbols}</p>
                  <p className="text-sm text-muted-foreground">Symbols</p>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">{dataState.stats.timeframes}</p>
                  <p className="text-sm text-muted-foreground">Timeframes</p>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">{(dataState.stats.storageSize / 1024 / 1024).toFixed(1)}MB</p>
                  <p className="text-sm text-muted-foreground">Storage Used</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
