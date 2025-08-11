'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Percent, 
  PieChart,
  BarChart3,
  Target,
  Shield,
  AlertTriangle,
  Calendar,
  Activity,
  Zap,
  ArrowUpRight,
  ArrowDownRight
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface PortfolioAsset {
  symbol: string
  name: string
  balance: number
  value: number
  price: number
  change24h: number
  allocation: number
  avgBuyPrice: number
  pnl: number
  pnlPercent: number
}

interface PerformanceMetric {
  label: string
  value: number
  change: number
  isPercentage: boolean
  period: string
}

interface RiskMetric {
  label: string
  value: number
  level: 'low' | 'medium' | 'high'
  description: string
}

interface TradeAnalytics {
  totalTrades: number
  winRate: number
  avgWin: number
  avgLoss: number
  profitFactor: number
  sharpeRatio: number
  maxDrawdown: number
}

export function AdvancedPortfolioAnalytics() {
  const [assets, setAssets] = useState<PortfolioAsset[]>([])
  const [performanceMetrics, setPerformanceMetrics] = useState<PerformanceMetric[]>([])
  const [riskMetrics, setRiskMetrics] = useState<RiskMetric[]>([])
  const [tradeAnalytics, setTradeAnalytics] = useState<TradeAnalytics | null>(null)
  const [timeframe, setTimeframe] = useState('7d')
  const [totalValue, setTotalValue] = useState(0)
  const [totalPnL, setTotalPnL] = useState(0)
  const [totalPnLPercent, setTotalPnLPercent] = useState(0)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock portfolio data
    const mockAssets: PortfolioAsset[] = [
      {
        symbol: 'ETH',
        name: 'Ethereum',
        balance: 5.2,
        value: 13000,
        price: 2500,
        change24h: 3.2,
        allocation: 52,
        avgBuyPrice: 2300,
        pnl: 1040,
        pnlPercent: 8.7
      },
      {
        symbol: 'BTC',
        name: 'Bitcoin',
        balance: 0.25,
        value: 11250,
        price: 45000,
        change24h: -1.5,
        allocation: 45,
        avgBuyPrice: 42000,
        pnl: 750,
        pnlPercent: 7.1
      },
      {
        symbol: 'USDC',
        name: 'USD Coin',
        balance: 750,
        value: 750,
        price: 1,
        change24h: 0,
        allocation: 3,
        avgBuyPrice: 1,
        pnl: 0,
        pnlPercent: 0
      }
    ]

    const mockPerformanceMetrics: PerformanceMetric[] = [
      { label: 'Total Return', value: 1790, change: 7.16, isPercentage: false, period: timeframe },
      { label: 'Return %', value: 7.16, change: 0.8, isPercentage: true, period: timeframe },
      { label: 'Volatility', value: 12.4, change: -2.1, isPercentage: true, period: timeframe },
      { label: 'Sharpe Ratio', value: 1.85, change: 0.15, isPercentage: false, period: timeframe }
    ]

    const mockRiskMetrics: RiskMetric[] = [
      {
        label: 'Portfolio Risk Score',
        value: 65,
        level: 'medium',
        description: 'Moderate risk with balanced allocation'
      },
      {
        label: 'Concentration Risk',
        value: 52,
        level: 'medium',
        description: 'ETH represents 52% of portfolio'
      },
      {
        label: 'Liquidity Risk',
        value: 15,
        level: 'low',
        description: 'High liquidity across all assets'
      }
    ]

    const mockTradeAnalytics: TradeAnalytics = {
      totalTrades: 47,
      winRate: 68.1,
      avgWin: 245.50,
      avgLoss: -128.75,
      profitFactor: 1.91,
      sharpeRatio: 1.85,
      maxDrawdown: -8.2
    }

    setAssets(mockAssets)
    setPerformanceMetrics(mockPerformanceMetrics)
    setRiskMetrics(mockRiskMetrics)
    setTradeAnalytics(mockTradeAnalytics)

    const total = mockAssets.reduce((sum, asset) => sum + asset.value, 0)
    const totalPnLValue = mockAssets.reduce((sum, asset) => sum + asset.pnl, 0)
    const totalPnLPercentValue = (totalPnLValue / (total - totalPnLValue)) * 100

    setTotalValue(total)
    setTotalPnL(totalPnLValue)
    setTotalPnLPercent(totalPnLPercentValue)
  }, [isConnected, timeframe])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 2
    }).format(amount)
  }

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'low': return 'text-green-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskBadgeVariant = (level: string) => {
    switch (level) {
      case 'low': return 'default'
      case 'medium': return 'secondary'
      case 'high': return 'destructive'
      default: return 'outline'
    }
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Your Wallet</h3>
          <p className="text-muted-foreground">
            Connect your wallet to view detailed portfolio analytics and performance metrics
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Portfolio Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Value</span>
            </div>
            <div className="text-2xl font-bold">{formatCurrency(totalValue)}</div>
            <div className={cn(
              "text-xs flex items-center gap-1",
              totalPnL >= 0 ? "text-green-500" : "text-red-500"
            )}>
              {totalPnL >= 0 ? <ArrowUpRight className="w-3 h-3" /> : <ArrowDownRight className="w-3 h-3" />}
              {totalPnL >= 0 ? '+' : ''}{formatCurrency(totalPnL)} ({totalPnLPercent >= 0 ? '+' : ''}{totalPnLPercent.toFixed(2)}%)
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">24h Change</span>
            </div>
            <div className="text-2xl font-bold text-green-500">+2.8%</div>
            <div className="text-xs text-muted-foreground">
              {formatCurrency(700)} gain
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Target className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Best Performer</span>
            </div>
            <div className="text-2xl font-bold">ETH</div>
            <div className="text-xs text-green-500">
              +8.7% total return
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Shield className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Risk Score</span>
            </div>
            <div className="text-2xl font-bold text-yellow-500">65</div>
            <div className="text-xs text-muted-foreground">
              Medium risk
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Time Frame Selector */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Portfolio Analytics</h2>
        <Select value={timeframe} onValueChange={setTimeframe}>
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1d">1 Day</SelectItem>
            <SelectItem value="7d">7 Days</SelectItem>
            <SelectItem value="30d">30 Days</SelectItem>
            <SelectItem value="90d">90 Days</SelectItem>
            <SelectItem value="1y">1 Year</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Main Analytics */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="assets">Assets</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="risk">Risk Analysis</TabsTrigger>
          <TabsTrigger value="trades">Trade Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Asset Allocation */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Asset Allocation
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {assets.map((asset) => (
                    <div key={asset.symbol} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 rounded-full bg-primary" />
                          <span className="font-medium">{asset.symbol}</span>
                        </div>
                        <div className="text-right">
                          <div className="font-medium">{asset.allocation}%</div>
                          <div className="text-xs text-muted-foreground">
                            {formatCurrency(asset.value)}
                          </div>
                        </div>
                      </div>
                      <Progress value={asset.allocation} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Performance Metrics */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Performance Metrics
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {performanceMetrics.map((metric) => (
                    <div key={metric.label} className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">{metric.label}</span>
                      <div className="text-right">
                        <div className="font-medium">
                          {metric.isPercentage ? `${metric.value.toFixed(2)}%` : 
                           metric.label === 'Total Return' ? formatCurrency(metric.value) : 
                           metric.value.toFixed(2)}
                        </div>
                        <div className={cn(
                          "text-xs flex items-center gap-1",
                          metric.change >= 0 ? "text-green-500" : "text-red-500"
                        )}>
                          {metric.change >= 0 ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
                          {metric.change >= 0 ? '+' : ''}{metric.change.toFixed(2)}
                          {metric.isPercentage ? 'pp' : '%'}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="assets" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Asset Breakdown</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {assets.map((asset) => (
                  <div key={asset.symbol} className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                          <span className="font-bold text-sm">{asset.symbol}</span>
                        </div>
                        <div>
                          <div className="font-medium">{asset.name}</div>
                          <div className="text-sm text-muted-foreground">
                            {asset.balance} {asset.symbol}
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold">{formatCurrency(asset.value)}</div>
                        <div className={cn(
                          "text-sm",
                          asset.pnl >= 0 ? "text-green-500" : "text-red-500"
                        )}>
                          {asset.pnl >= 0 ? '+' : ''}{formatCurrency(asset.pnl)} ({asset.pnlPercent >= 0 ? '+' : ''}{asset.pnlPercent.toFixed(2)}%)
                        </div>
                      </div>
                    </div>
                    
                    <div className="grid grid-cols-3 gap-4 text-sm">
                      <div>
                        <div className="text-muted-foreground">Current Price</div>
                        <div className="font-medium">{formatCurrency(asset.price)}</div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Avg Buy Price</div>
                        <div className="font-medium">{formatCurrency(asset.avgBuyPrice)}</div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">24h Change</div>
                        <div className={cn(
                          "font-medium",
                          asset.change24h >= 0 ? "text-green-500" : "text-red-500"
                        )}>
                          {asset.change24h >= 0 ? '+' : ''}{asset.change24h.toFixed(2)}%
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="performance" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle>Performance Chart</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-64 flex items-center justify-center text-muted-foreground">
                  <div className="text-center">
                    <BarChart3 className="w-12 h-12 mx-auto mb-2 opacity-50" />
                    <p>Performance chart visualization</p>
                    <p className="text-xs">Integration with charting library needed</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Key Metrics</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {performanceMetrics.map((metric) => (
                    <div key={metric.label} className="flex items-center justify-between p-3 border rounded">
                      <span className="font-medium">{metric.label}</span>
                      <div className="text-right">
                        <div className="font-bold">
                          {metric.isPercentage ? `${metric.value.toFixed(2)}%` : 
                           metric.label === 'Total Return' ? formatCurrency(metric.value) : 
                           metric.value.toFixed(2)}
                        </div>
                        <div className={cn(
                          "text-xs",
                          metric.change >= 0 ? "text-green-500" : "text-red-500"
                        )}>
                          {metric.change >= 0 ? '+' : ''}{metric.change.toFixed(2)}
                          {metric.isPercentage ? 'pp' : '%'} vs prev period
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="risk" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Risk Assessment
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {riskMetrics.map((risk) => (
                    <div key={risk.label} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="font-medium">{risk.label}</span>
                        <Badge variant={getRiskBadgeVariant(risk.level)}>
                          {risk.level}
                        </Badge>
                      </div>
                      <Progress value={risk.value} className="h-2" />
                      <p className="text-xs text-muted-foreground">{risk.description}</p>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="w-5 h-5" />
                  Risk Recommendations
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="p-3 bg-yellow-50 dark:bg-yellow-950/20 border border-yellow-200 dark:border-yellow-800 rounded">
                    <div className="font-medium text-yellow-800 dark:text-yellow-200">
                      High ETH Concentration
                    </div>
                    <p className="text-xs text-yellow-700 dark:text-yellow-300 mt-1">
                      Consider diversifying your portfolio to reduce concentration risk
                    </p>
                  </div>
                  
                  <div className="p-3 bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800 rounded">
                    <div className="font-medium text-blue-800 dark:text-blue-200">
                      Add Stablecoins
                    </div>
                    <p className="text-xs text-blue-700 dark:text-blue-300 mt-1">
                      Increase stablecoin allocation for better risk management
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="trades" className="space-y-4">
          {tradeAnalytics && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Activity className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">Total Trades</span>
                  </div>
                  <div className="text-2xl font-bold">{tradeAnalytics.totalTrades}</div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Target className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">Win Rate</span>
                  </div>
                  <div className="text-2xl font-bold text-green-500">{tradeAnalytics.winRate}%</div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <TrendingUp className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">Profit Factor</span>
                  </div>
                  <div className="text-2xl font-bold">{tradeAnalytics.profitFactor}</div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <TrendingDown className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">Max Drawdown</span>
                  </div>
                  <div className="text-2xl font-bold text-red-500">{tradeAnalytics.maxDrawdown}%</div>
                </CardContent>
              </Card>
            </div>
          )}

          <Card>
            <CardHeader>
              <CardTitle>Trading Performance</CardTitle>
            </CardHeader>
            <CardContent>
              {tradeAnalytics && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Average Win</span>
                      <span className="font-medium text-green-500">
                        {formatCurrency(tradeAnalytics.avgWin)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Average Loss</span>
                      <span className="font-medium text-red-500">
                        {formatCurrency(tradeAnalytics.avgLoss)}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Sharpe Ratio</span>
                      <span className="font-medium">{tradeAnalytics.sharpeRatio}</span>
                    </div>
                  </div>
                  
                  <div className="space-y-4">
                    <div>
                      <div className="flex justify-between mb-2">
                        <span className="text-sm text-muted-foreground">Win Rate</span>
                        <span className="text-sm font-medium">{tradeAnalytics.winRate}%</span>
                      </div>
                      <Progress value={tradeAnalytics.winRate} className="h-2" />
                    </div>
                    
                    <div>
                      <div className="flex justify-between mb-2">
                        <span className="text-sm text-muted-foreground">Profit Factor</span>
                        <span className="text-sm font-medium">{tradeAnalytics.profitFactor}</span>
                      </div>
                      <Progress value={(tradeAnalytics.profitFactor / 3) * 100} className="h-2" />
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
