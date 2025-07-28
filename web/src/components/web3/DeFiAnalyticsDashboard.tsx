'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  BarChart3, 
  TrendingUp, 
  TrendingDown,
  DollarSign, 
  Activity,
  RefreshCw,
  Users,
  Target,
  Shield,
  Zap,
  PieChart,
  LineChart,
  Award,
  AlertTriangle,
  CheckCircle,
  ArrowUpRight,
  ArrowDownRight,
  Minus
} from 'lucide-react'
import { useDeFiAnalytics } from '@/hooks/useDeFiAnalytics'
import { ProtocolCategory, RiskLevel } from '@/lib/defi-analytics'
import { type Address } from 'viem'

interface DeFiAnalyticsDashboardProps {
  userAddress?: Address
}

export function DeFiAnalyticsDashboard({ userAddress }: DeFiAnalyticsDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const {
    protocolMetrics,
    apyData,
    yieldOpportunities,
    defiTrends,
    portfolioAnalytics,
    isLoading,
    loadData,
    analyticsMetrics,
    categoryDistribution,
    riskDistribution,
    portfolioMetrics,
    topProtocolsByTVL,
    topProtocolsByVolume,
    bestYieldOpportunities,
    formatAPY,
    formatTVL,
    formatVolume,
    getRiskColor,
    getCategoryIcon,
    getTrendColor,
    isBullishTrend,
    isBearishTrend
  } = useDeFiAnalytics({
    userAddress,
    autoRefresh: true,
    enableNotifications: true
  })

  const getTrendIcon = (trend: 'up' | 'down' | 'stable') => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="w-4 h-4" />
      case 'down':
        return <TrendingDown className="w-4 h-4" />
      default:
        return <Minus className="w-4 h-4" />
    }
  }

  const getChangeIcon = (change: string) => {
    const value = parseFloat(change)
    if (value > 0) return <ArrowUpRight className="w-3 h-3 text-green-500" />
    if (value < 0) return <ArrowDownRight className="w-3 h-3 text-red-500" />
    return <Minus className="w-3 h-3 text-gray-500" />
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <BarChart3 className="w-6 h-6" />
            DeFi Analytics & APY Tracking
          </h2>
          <p className="text-muted-foreground">
            Real-time DeFi protocol analytics, APY tracking, and yield optimization insights
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadData}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Market Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total TVL</p>
                <p className="text-2xl font-bold">{formatTVL(analyticsMetrics.totalTVL.toString())}</p>
                <p className={`text-sm flex items-center gap-1 ${getTrendColor(analyticsMetrics.tvlTrend)}`}>
                  {getTrendIcon(analyticsMetrics.tvlTrend)}
                  {defiTrends?.tvlTrend.change24h}% 24h
                </p>
              </div>
              <DollarSign className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">24h Volume</p>
                <p className="text-2xl font-bold">{formatVolume(analyticsMetrics.totalVolume24h.toString())}</p>
                <p className={`text-sm flex items-center gap-1 ${getTrendColor(defiTrends?.volumeTrend.trend || 'stable')}`}>
                  {getTrendIcon(defiTrends?.volumeTrend.trend || 'stable')}
                  {defiTrends?.volumeTrend.change24h}% 24h
                </p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Average APY</p>
                <p className="text-2xl font-bold">{formatAPY(analyticsMetrics.averageAPY.toString())}</p>
                <p className={`text-sm flex items-center gap-1 ${getTrendColor(analyticsMetrics.apyTrend)}`}>
                  {getTrendIcon(analyticsMetrics.apyTrend)}
                  {defiTrends?.apyTrend.change24h}% 24h
                </p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Users</p>
                <p className="text-2xl font-bold">{analyticsMetrics.totalUsers24h.toLocaleString()}</p>
                <p className={`text-sm flex items-center gap-1 ${getTrendColor(defiTrends?.userTrend.trend || 'stable')}`}>
                  {getTrendIcon(defiTrends?.userTrend.trend || 'stable')}
                  {defiTrends?.userTrend.change24h}% 24h
                </p>
              </div>
              <Users className="w-8 h-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Market Sentiment */}
      {(isBullishTrend || isBearishTrend) && (
        <Card>
          <CardContent className="p-4">
            <div className={`flex items-center gap-3 p-3 rounded-lg ${
              isBullishTrend ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
            }`}>
              {isBullishTrend ? (
                <CheckCircle className="w-5 h-5 text-green-600" />
              ) : (
                <AlertTriangle className="w-5 h-5 text-red-600" />
              )}
              <div>
                <p className={`font-medium ${isBullishTrend ? 'text-green-800' : 'text-red-800'}`}>
                  {isBullishTrend ? 'Bullish Market Sentiment' : 'Bearish Market Sentiment'}
                </p>
                <p className={`text-sm ${isBullishTrend ? 'text-green-600' : 'text-red-600'}`}>
                  {isBullishTrend 
                    ? 'TVL and volume are both trending upward, indicating strong market confidence'
                    : 'TVL and volume are both declining, suggesting market uncertainty'
                  }
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="protocols">Protocols</TabsTrigger>
          <TabsTrigger value="yields">Yields</TabsTrigger>
          <TabsTrigger value="trends">Trends</TabsTrigger>
          <TabsTrigger value="portfolio">Portfolio</TabsTrigger>
          <TabsTrigger value="opportunities">Opportunities</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Category Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>TVL by Category</CardTitle>
                <CardDescription>
                  Total value locked distribution across DeFi categories
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {categoryDistribution.map((item, index) => (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-lg">{getCategoryIcon(item.category)}</span>
                          <span className="text-sm font-medium capitalize">{item.category.replace('_', ' ')}</span>
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {formatTVL(item.tvl.toString())} ({item.percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <Progress value={item.percentage} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Distribution</CardTitle>
                <CardDescription>
                  Yield opportunities by risk level
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {riskDistribution.map((item, index) => (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <Shield className="w-4 h-4" />
                          <span className="text-sm font-medium capitalize">{item.riskLevel.replace('_', ' ')}</span>
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {item.count} ({item.percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <Progress value={item.percentage} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Top Protocols */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Top Protocols by TVL</CardTitle>
                <CardDescription>
                  Protocols with highest total value locked
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {topProtocolsByTVL.map((protocol, index) => (
                    <motion.div
                      key={protocol.protocolId}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center text-sm font-bold">
                          {index + 1}
                        </div>
                        <div>
                          <p className="text-sm font-medium">{protocol.name}</p>
                          <p className="text-xs text-muted-foreground flex items-center gap-1">
                            {getCategoryIcon(protocol.category)}
                            {protocol.category.replace('_', ' ')}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-bold">{formatTVL(protocol.tvl)}</p>
                        <p className="text-xs text-muted-foreground">
                          {protocol.marketShare.toFixed(1)}% share
                        </p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Top Protocols by Volume</CardTitle>
                <CardDescription>
                  Protocols with highest 24h trading volume
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {topProtocolsByVolume.map((protocol, index) => (
                    <motion.div
                      key={protocol.protocolId}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center text-sm font-bold">
                          {index + 1}
                        </div>
                        <div>
                          <p className="text-sm font-medium">{protocol.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {protocol.users24h.toLocaleString()} users
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-bold">{formatVolume(protocol.volume24h)}</p>
                        <p className="text-xs text-muted-foreground">24h volume</p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="protocols" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Protocol Analytics</CardTitle>
              <CardDescription>
                Comprehensive metrics for all tracked DeFi protocols
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {protocolMetrics.map((protocol, index) => (
                  <motion.div
                    key={protocol.protocolId}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex items-center gap-3">
                        <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center text-2xl">
                          {getCategoryIcon(protocol.category)}
                        </div>
                        <div>
                          <h4 className="font-medium">{protocol.name}</h4>
                          <p className="text-sm text-muted-foreground capitalize">
                            {protocol.category.replace('_', ' ')}
                          </p>
                        </div>
                      </div>
                      <Badge variant="outline">
                        Rank #{Math.floor(protocol.dominanceIndex)}
                      </Badge>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div>
                        <p className="text-xs text-muted-foreground">TVL</p>
                        <p className="font-bold">{formatTVL(protocol.tvl)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">24h Volume</p>
                        <p className="font-bold">{formatVolume(protocol.volume24h)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">24h Fees</p>
                        <p className="font-bold">{formatTVL(protocol.fees24h)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Users 24h</p>
                        <p className="font-bold">{protocol.users24h.toLocaleString()}</p>
                      </div>
                    </div>

                    <div className="mt-4 pt-4 border-t">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Market Share</span>
                        <span className="font-medium">{protocol.marketShare.toFixed(1)}%</span>
                      </div>
                      <Progress value={protocol.marketShare} className="h-2 mt-2" />
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="yields" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>APY Tracking</CardTitle>
              <CardDescription>
                Real-time APY data across all protocols and pools
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {apyData.map((apy, index) => (
                  <motion.div
                    key={`${apy.protocolId}-${apy.poolId}`}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="border rounded-lg p-4"
                  >
                    <div className="flex items-center justify-between mb-4">
                      <div>
                        <h4 className="font-medium">{apy.poolName}</h4>
                        <p className="text-sm text-muted-foreground">{apy.protocolId}</p>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge className={getRiskColor(apy.riskLevel)}>
                          {apy.riskLevel.replace('_', ' ')}
                        </Badge>
                        <div className="text-right">
                          <p className="font-bold text-lg">{formatAPY(apy.totalAPY)}</p>
                          <p className="text-xs text-muted-foreground">Total APY</p>
                        </div>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div>
                        <p className="text-xs text-muted-foreground">Base APY</p>
                        <p className="font-medium">{formatAPY(apy.baseAPY)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Reward APY</p>
                        <p className="font-medium">{formatAPY(apy.rewardAPY)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">TVL</p>
                        <p className="font-medium">{formatTVL(apy.tvl)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Utilization</p>
                        <p className="font-medium">{apy.utilization}%</p>
                      </div>
                    </div>

                    <div className="mt-4 pt-4 border-t">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">7d APY</span>
                        <span className="flex items-center gap-1">
                          {getChangeIcon((parseFloat(apy.totalAPY) - parseFloat(apy.apy7d)).toString())}
                          {formatAPY(apy.apy7d)}
                        </span>
                      </div>
                      <div className="flex items-center justify-between text-sm mt-1">
                        <span className="text-muted-foreground">30d APY</span>
                        <span className="flex items-center gap-1">
                          {getChangeIcon((parseFloat(apy.totalAPY) - parseFloat(apy.apy30d)).toString())}
                          {formatAPY(apy.apy30d)}
                        </span>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="trends" className="space-y-6">
          {defiTrends && (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Market Trends</CardTitle>
                    <CardDescription>
                      Key DeFi market indicators and trends
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <DollarSign className="w-4 h-4" />
                          <span className="text-sm font-medium">TVL Trend</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-sm ${getTrendColor(defiTrends.tvlTrend.trend)}`}>
                            {defiTrends.tvlTrend.change24h}%
                          </span>
                          {getTrendIcon(defiTrends.tvlTrend.trend)}
                        </div>
                      </div>

                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <Target className="w-4 h-4" />
                          <span className="text-sm font-medium">APY Trend</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-sm ${getTrendColor(defiTrends.apyTrend.trend)}`}>
                            {defiTrends.apyTrend.change24h}%
                          </span>
                          {getTrendIcon(defiTrends.apyTrend.trend)}
                        </div>
                      </div>

                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <Activity className="w-4 h-4" />
                          <span className="text-sm font-medium">Volume Trend</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-sm ${getTrendColor(defiTrends.volumeTrend.trend)}`}>
                            {defiTrends.volumeTrend.change24h}%
                          </span>
                          {getTrendIcon(defiTrends.volumeTrend.trend)}
                        </div>
                      </div>

                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <Users className="w-4 h-4" />
                          <span className="text-sm font-medium">User Trend</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-sm ${getTrendColor(defiTrends.userTrend.trend)}`}>
                            {defiTrends.userTrend.change24h}%
                          </span>
                          {getTrendIcon(defiTrends.userTrend.trend)}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Top Movers</CardTitle>
                    <CardDescription>
                      Biggest gainers and losers in the DeFi space
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div>
                        <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
                          <TrendingUp className="w-4 h-4 text-green-500" />
                          Top Gainers
                        </h4>
                        <div className="space-y-2">
                          {defiTrends.topGainers.map((gainer, index) => (
                            <div key={index} className="flex items-center justify-between p-2 bg-green-50 rounded-lg">
                              <span className="text-sm font-medium">{gainer.name}</span>
                              <span className="text-sm text-green-600">+{gainer.change}%</span>
                            </div>
                          ))}
                        </div>
                      </div>

                      <div>
                        <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
                          <TrendingDown className="w-4 h-4 text-red-500" />
                          Top Losers
                        </h4>
                        <div className="space-y-2">
                          {defiTrends.topLosers.map((loser, index) => (
                            <div key={index} className="flex items-center justify-between p-2 bg-red-50 rounded-lg">
                              <span className="text-sm font-medium">{loser.name}</span>
                              <span className="text-sm text-red-600">{loser.change}%</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </>
          )}
        </TabsContent>

        <TabsContent value="portfolio" className="space-y-6">
          {portfolioAnalytics ? (
            <>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <Card>
                  <CardContent className="p-6">
                    <div className="text-center">
                      <p className="text-sm font-medium text-muted-foreground">Portfolio Value</p>
                      <p className="text-2xl font-bold">{formatTVL(portfolioMetrics.totalValue.toString())}</p>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="text-center">
                      <p className="text-sm font-medium text-muted-foreground">24h Yield</p>
                      <p className="text-2xl font-bold text-green-600">{formatTVL(portfolioMetrics.totalYield24h.toString())}</p>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="text-center">
                      <p className="text-sm font-medium text-muted-foreground">Average APY</p>
                      <p className="text-2xl font-bold">{formatAPY(portfolioMetrics.averageAPY.toString())}</p>
                    </div>
                  </CardContent>
                </Card>

                <Card>
                  <CardContent className="p-6">
                    <div className="text-center">
                      <p className="text-sm font-medium text-muted-foreground">Risk Score</p>
                      <p className="text-2xl font-bold">{portfolioMetrics.riskScore}/100</p>
                    </div>
                  </CardContent>
                </Card>
              </div>

              <Card>
                <CardHeader>
                  <CardTitle>Portfolio Positions</CardTitle>
                  <CardDescription>
                    Your active DeFi positions and performance
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {portfolioAnalytics.positions.map((position, index) => (
                      <motion.div
                        key={`${position.protocolId}-${position.poolName}`}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className="border rounded-lg p-4"
                      >
                        <div className="flex items-center justify-between mb-4">
                          <div>
                            <h4 className="font-medium">{position.poolName}</h4>
                            <p className="text-sm text-muted-foreground">{position.protocolName}</p>
                          </div>
                          <div className="text-right">
                            <p className="font-bold">{formatTVL(position.value)}</p>
                            <p className="text-sm text-muted-foreground">{position.allocation}% allocation</p>
                          </div>
                        </div>

                        <div className="grid grid-cols-3 gap-4">
                          <div>
                            <p className="text-xs text-muted-foreground">APY</p>
                            <p className="font-medium">{formatAPY(position.apy)}</p>
                          </div>
                          <div>
                            <p className="text-xs text-muted-foreground">24h Yield</p>
                            <p className="font-medium text-green-600">{formatTVL(position.yield24h)}</p>
                          </div>
                          <div>
                            <p className="text-xs text-muted-foreground">Risk Level</p>
                            <Badge className={getRiskColor(position.riskLevel)}>
                              {position.riskLevel.replace('_', ' ')}
                            </Badge>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              {portfolioAnalytics.recommendations.length > 0 && (
                <Card>
                  <CardHeader>
                    <CardTitle>Optimization Recommendations</CardTitle>
                    <CardDescription>
                      AI-powered suggestions to improve your portfolio
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      {portfolioAnalytics.recommendations.map((rec, index) => (
                        <div key={index} className="p-3 border rounded-lg">
                          <div className="flex items-center gap-2 mb-2">
                            <Award className="w-4 h-4 text-blue-500" />
                            <span className="font-medium capitalize">{rec.type}</span>
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">{rec.description}</p>
                          <div className="flex items-center gap-4 text-xs">
                            <span className="text-green-600">Expected Gain: {rec.expectedGain}</span>
                            <span className="text-orange-600">Risk Change: {rec.riskChange}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              )}
            </>
          ) : (
            <Card>
              <CardContent className="p-8 text-center">
                <PieChart className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">No Portfolio Data</h3>
                <p className="text-muted-foreground">
                  {userAddress 
                    ? "You don't have any active DeFi positions yet"
                    : "Connect your wallet to view portfolio analytics"
                  }
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="opportunities" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Best Yield Opportunities</CardTitle>
              <CardDescription>
                Top yield opportunities ranked by APY and risk-adjusted returns
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {bestYieldOpportunities.map((opportunity, index) => (
                  <motion.div
                    key={opportunity.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center text-sm font-bold">
                          {index + 1}
                        </div>
                        <div>
                          <h4 className="font-medium">{opportunity.poolName}</h4>
                          <p className="text-sm text-muted-foreground">{opportunity.protocolName}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-lg text-green-600">{formatAPY(opportunity.apy)}</p>
                        <Badge className={getRiskColor(opportunity.riskLevel)}>
                          {opportunity.riskLevel.replace('_', ' ')}
                        </Badge>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                      <div>
                        <p className="text-xs text-muted-foreground">TVL</p>
                        <p className="font-medium">{formatTVL(opportunity.tvl)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Min Deposit</p>
                        <p className="font-medium">{opportunity.requirements.minDeposit} ETH</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Liquidity Score</p>
                        <p className="font-medium">{opportunity.liquidityScore}/100</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Stability Score</p>
                        <p className="font-medium">{opportunity.stabilityScore}/100</p>
                      </div>
                    </div>

                    <div className="flex items-center justify-between pt-4 border-t">
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-muted-foreground">Fees:</span>
                        <span className="text-sm">Performance {opportunity.fees.performance}%</span>
                      </div>
                      <Button size="sm">
                        <Zap className="w-3 h-3 mr-2" />
                        Start Earning
                      </Button>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
