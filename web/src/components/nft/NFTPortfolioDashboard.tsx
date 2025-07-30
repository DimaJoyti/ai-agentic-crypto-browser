'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Wallet, 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Activity, 
  PieChart,
  BarChart3,
  Shield,
  Target,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  Star,
  Image as ImageIcon
} from 'lucide-react'
import { useNFTPortfolio } from '@/hooks/useNFTPortfolio'
import { type Address } from 'viem'
import { formatCurrency, formatNumber, formatPercentage } from '@/lib/utils'

interface NFTPortfolioDashboardProps {
  ownerAddress: Address
  className?: string
}

export function NFTPortfolioDashboard({ ownerAddress, className }: NFTPortfolioDashboardProps) {
  const { 
    state, 
    loadPortfolio, 
    updatePortfolio,
    getPortfolioSummary,
    getPerformanceMetrics,
    getRiskAnalysis,
    getDiversificationAnalysis,
    clearError
  } = useNFTPortfolio({ enableNotifications: true })

  const [activeTab, setActiveTab] = useState('overview')

  const portfolioSummary = getPortfolioSummary()
  const performanceMetrics = getPerformanceMetrics()
  const riskAnalysis = getRiskAnalysis()
  const diversificationAnalysis = getDiversificationAnalysis()

  useEffect(() => {
    loadPortfolio(ownerAddress)
  }, [ownerAddress, loadPortfolio])

  const handleRefresh = () => {
    updatePortfolio()
  }

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'low': return 'text-green-600 bg-green-100'
      case 'medium': return 'text-yellow-600 bg-yellow-100'
      case 'high': return 'text-orange-600 bg-orange-100'
      case 'critical': return 'text-red-600 bg-red-100'
      default: return 'text-gray-600 bg-gray-100'
    }
  }

  const getPerformanceColor = (value: number) => {
    return value >= 0 ? 'text-green-600' : 'text-red-600'
  }

  const getPerformanceIcon = (value: number) => {
    return value >= 0 ? <TrendingUp className="h-4 w-4" /> : <TrendingDown className="h-4 w-4" />
  }

  if (state.error) {
    return (
      <Alert className="border-red-200 bg-red-50">
        <AlertTriangle className="h-4 w-4 text-red-600" />
        <AlertDescription className="text-red-800">
          Failed to load portfolio: {state.error}
          <Button 
            variant="outline" 
            size="sm" 
            onClick={clearError}
            className="ml-2"
          >
            Retry
          </Button>
        </AlertDescription>
      </Alert>
    )
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">NFT Portfolio</h2>
          <p className="text-muted-foreground">
            {ownerAddress.slice(0, 6)}...{ownerAddress.slice(-4)}
          </p>
        </div>
        <div className="flex items-center gap-2">
          {state.lastUpdate && (
            <span className="text-sm text-muted-foreground">
              Last updated: {new Date(state.lastUpdate).toLocaleTimeString()}
            </span>
          )}
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={state.isLoading}
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="collections">Collections</TabsTrigger>
          <TabsTrigger value="assets">Assets</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="risk">Risk Analysis</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          {/* Portfolio Summary */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Value</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {formatCurrency(portfolioSummary.totalValue)} ETH
                </div>
                <p className="text-xs text-muted-foreground">
                  ${formatNumber(portfolioSummary.totalValueUSD)}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
                {getPerformanceIcon(portfolioSummary.totalUnrealizedPnL)}
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${getPerformanceColor(portfolioSummary.totalUnrealizedPnL)}`}>
                  {portfolioSummary.totalUnrealizedPnL >= 0 ? '+' : ''}
                  {formatCurrency(portfolioSummary.totalUnrealizedPnL)} ETH
                </div>
                <p className={`text-xs ${getPerformanceColor(portfolioSummary.totalUnrealizedPnLPercent)}`}>
                  {portfolioSummary.totalUnrealizedPnLPercent >= 0 ? '+' : ''}
                  {formatPercentage(portfolioSummary.totalUnrealizedPnLPercent)}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Assets</CardTitle>
                <ImageIcon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatNumber(portfolioSummary.totalAssets)}</div>
                <p className="text-xs text-muted-foreground">
                  {portfolioSummary.totalCollections} collections
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Risk Level</CardTitle>
                <Shield className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{riskAnalysis.riskScore}/100</div>
                <Badge className={getRiskColor(riskAnalysis.overallRisk)} variant="secondary">
                  {riskAnalysis.overallRisk.toUpperCase()}
                </Badge>
              </CardContent>
            </Card>
          </div>

          {/* Top Collections */}
          {state.analytics && (
            <Card>
              <CardHeader>
                <CardTitle>Top Collections</CardTitle>
                <CardDescription>Your largest holdings by value</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {state.analytics.topCollections.slice(0, 5).map(collection => (
                    <div key={collection.collection.contractAddress} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <img 
                          src={collection.collection.imageUrl} 
                          alt={collection.collection.name}
                          className="w-10 h-10 rounded-lg object-cover"
                        />
                        <div>
                          <div className="font-medium">{collection.collection.name}</div>
                          <div className="text-sm text-muted-foreground">
                            {collection.totalOwned} items • Floor: {formatCurrency(collection.collection.floorPrice)} ETH
                          </div>
                        </div>
                        {collection.collection.verified && (
                          <CheckCircle className="h-4 w-4 text-blue-600" />
                        )}
                      </div>
                      <div className="text-right">
                        <div className="font-medium">
                          {formatCurrency(collection.totalValue)} ETH
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {formatPercentage(collection.allocation)}
                        </div>
                        <div className={`text-sm ${getPerformanceColor(collection.performance)}`}>
                          {collection.performance >= 0 ? '+' : ''}
                          {formatPercentage(collection.performance)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Performance Overview */}
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Performance Metrics</CardTitle>
                <CardDescription>Key performance indicators</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>Win Rate</span>
                    <span className="font-medium">{formatPercentage(performanceMetrics.winRate)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Sharpe Ratio</span>
                    <span className="font-medium">{performanceMetrics.sharpeRatio.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Max Drawdown</span>
                    <span className="font-medium text-red-600">
                      -{formatPercentage(performanceMetrics.maxDrawdown)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span>Volatility</span>
                    <span className="font-medium">{formatPercentage(performanceMetrics.volatility)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Avg Holding Period</span>
                    <span className="font-medium">{performanceMetrics.averageHoldingPeriod} days</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Analysis</CardTitle>
                <CardDescription>Portfolio risk breakdown</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span>Liquidity Risk</span>
                    <div className="flex items-center gap-2">
                      <Progress value={riskAnalysis.liquidityRisk} className="w-20 h-2" />
                      <span className="text-sm font-medium">{riskAnalysis.liquidityRisk}/100</span>
                    </div>
                  </div>
                  <div className="flex justify-between items-center">
                    <span>Concentration Risk</span>
                    <div className="flex items-center gap-2">
                      <Progress value={riskAnalysis.concentrationRisk} className="w-20 h-2" />
                      <span className="text-sm font-medium">{formatPercentage(riskAnalysis.concentrationRisk)}</span>
                    </div>
                  </div>
                  <div className="flex justify-between items-center">
                    <span>Market Risk</span>
                    <div className="flex items-center gap-2">
                      <Progress value={riskAnalysis.marketRisk} className="w-20 h-2" />
                      <span className="text-sm font-medium">{riskAnalysis.marketRisk}/100</span>
                    </div>
                  </div>
                  <div className="flex justify-between items-center">
                    <span>Diversification Score</span>
                    <div className="flex items-center gap-2">
                      <Progress value={diversificationAnalysis.diversificationScore} className="w-20 h-2" />
                      <span className="text-sm font-medium">{diversificationAnalysis.diversificationScore}/100</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="collections" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Collection Holdings</CardTitle>
              <CardDescription>Detailed view of your collection holdings</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {state.collections.map(collection => (
                  <div key={`${collection.contractAddress}_${collection.chainId}`} className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <img 
                          src={collection.imageUrl} 
                          alt={collection.name}
                          className="w-12 h-12 rounded-lg object-cover"
                        />
                        <div>
                          <div className="font-medium text-lg">{collection.name}</div>
                          <div className="text-sm text-muted-foreground">
                            {collection.totalOwned} items • Floor: {formatCurrency(collection.floorPrice)} ETH
                          </div>
                        </div>
                        {collection.verified && (
                          <CheckCircle className="h-5 w-5 text-blue-600" />
                        )}
                      </div>
                      <Badge className={getRiskColor('medium')} variant="secondary">
                        {formatPercentage(collection.allocation)} allocation
                      </Badge>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                      <div>
                        <div className="text-muted-foreground">Total Value</div>
                        <div className="font-medium text-lg">
                          {formatCurrency(collection.totalValue)} ETH
                        </div>
                        <div className="text-muted-foreground">
                          ${formatNumber(collection.totalValueUSD)}
                        </div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Unrealized P&L</div>
                        <div className={`font-medium text-lg ${getPerformanceColor(collection.unrealizedPnL)}`}>
                          {collection.unrealizedPnL >= 0 ? '+' : ''}
                          {formatCurrency(collection.unrealizedPnL)} ETH
                        </div>
                        <div className={`text-sm ${getPerformanceColor(collection.unrealizedPnLPercent)}`}>
                          {collection.unrealizedPnLPercent >= 0 ? '+' : ''}
                          {formatPercentage(collection.unrealizedPnLPercent)}
                        </div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Avg Cost Basis</div>
                        <div className="font-medium text-lg">
                          {formatCurrency(collection.averageCostBasis)} ETH
                        </div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Rarity Score</div>
                        <div className="font-medium text-lg">
                          {collection.rarity.rarityScore.toFixed(1)}
                        </div>
                        <div className="flex items-center gap-1">
                          <Star className="h-3 w-3 text-yellow-500" />
                          <span className="text-sm">Top {formatPercentage(100 - collection.rarity.averageRarity)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="assets" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Individual Assets</CardTitle>
              <CardDescription>Detailed view of your NFT assets</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {state.assets.map(asset => (
                  <div key={asset.id} className="border rounded-lg p-4">
                    <div className="aspect-square mb-3">
                      <img 
                        src={asset.imageUrl} 
                        alt={asset.name}
                        className="w-full h-full object-cover rounded-lg"
                      />
                    </div>
                    <div className="space-y-2">
                      <div>
                        <div className="font-medium">{asset.name}</div>
                        <div className="text-sm text-muted-foreground">
                          {asset.collection.name} • #{asset.tokenId}
                        </div>
                      </div>
                      
                      <div className="flex justify-between text-sm">
                        <span>Current Value</span>
                        <span className="font-medium">
                          {formatCurrency(asset.valuation.currentValue)} ETH
                        </span>
                      </div>
                      
                      <div className="flex justify-between text-sm">
                        <span>Cost Basis</span>
                        <span>{formatCurrency(asset.performance.costBasis)} ETH</span>
                      </div>
                      
                      <div className="flex justify-between text-sm">
                        <span>P&L</span>
                        <span className={getPerformanceColor(asset.performance.unrealizedPnL)}>
                          {asset.performance.unrealizedPnL >= 0 ? '+' : ''}
                          {formatCurrency(asset.performance.unrealizedPnL)} ETH
                          ({asset.performance.unrealizedPnLPercent >= 0 ? '+' : ''}
                          {formatPercentage(asset.performance.unrealizedPnLPercent)})
                        </span>
                      </div>
                      
                      <div className="flex justify-between text-sm">
                        <span>Rarity</span>
                        <div className="flex items-center gap-1">
                          <Star className="h-3 w-3 text-yellow-500" />
                          <span>#{asset.rarity.rank} ({asset.rarity.tier})</span>
                        </div>
                      </div>
                      
                      <div className="flex justify-between text-sm">
                        <span>Holding Period</span>
                        <span>{asset.performance.holdingPeriod} days</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="performance" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Performance Summary</CardTitle>
                <CardDescription>Overall portfolio performance</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className={`text-3xl font-bold ${getPerformanceColor(performanceMetrics.totalReturn)}`}>
                      {performanceMetrics.totalReturn >= 0 ? '+' : ''}
                      {formatCurrency(performanceMetrics.totalReturn)} ETH
                    </div>
                    <div className={`text-lg ${getPerformanceColor(performanceMetrics.totalReturnPercent)}`}>
                      {performanceMetrics.totalReturnPercent >= 0 ? '+' : ''}
                      {formatPercentage(performanceMetrics.totalReturnPercent)}
                    </div>
                    <div className="text-sm text-muted-foreground">Total Return</div>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4 text-center">
                    <div>
                      <div className="text-2xl font-bold text-green-600">
                        {formatPercentage(performanceMetrics.winRate)}
                      </div>
                      <div className="text-sm text-muted-foreground">Win Rate</div>
                    </div>
                    <div>
                      <div className="text-2xl font-bold">
                        {performanceMetrics.profitFactor.toFixed(2)}
                      </div>
                      <div className="text-sm text-muted-foreground">Profit Factor</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Metrics</CardTitle>
                <CardDescription>Risk-adjusted performance</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>Sharpe Ratio</span>
                    <span className="font-medium">{performanceMetrics.sharpeRatio.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Max Drawdown</span>
                    <span className="font-medium text-red-600">
                      -{formatPercentage(performanceMetrics.maxDrawdown)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span>Volatility</span>
                    <span className="font-medium">{formatPercentage(performanceMetrics.volatility)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Average Holding Period</span>
                    <span className="font-medium">{performanceMetrics.averageHoldingPeriod} days</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Best and Worst Performers */}
          {portfolioSummary.bestPerformer && portfolioSummary.worstPerformer && (
            <div className="grid gap-4 md:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle className="text-green-600">Best Performer</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-3">
                    <img 
                      src={portfolioSummary.bestPerformer.imageUrl} 
                      alt={portfolioSummary.bestPerformer.name}
                      className="w-16 h-16 rounded-lg object-cover"
                    />
                    <div>
                      <div className="font-medium">{portfolioSummary.bestPerformer.name}</div>
                      <div className="text-sm text-muted-foreground">
                        #{portfolioSummary.bestPerformer.tokenId}
                      </div>
                      <div className="text-lg font-bold text-green-600">
                        +{formatCurrency(portfolioSummary.bestPerformer.performance.unrealizedPnL)} ETH
                      </div>
                      <div className="text-sm text-green-600">
                        +{formatPercentage(portfolioSummary.bestPerformer.performance.unrealizedPnLPercent)}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-red-600">Worst Performer</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center gap-3">
                    <img 
                      src={portfolioSummary.worstPerformer.imageUrl} 
                      alt={portfolioSummary.worstPerformer.name}
                      className="w-16 h-16 rounded-lg object-cover"
                    />
                    <div>
                      <div className="font-medium">{portfolioSummary.worstPerformer.name}</div>
                      <div className="text-sm text-muted-foreground">
                        #{portfolioSummary.worstPerformer.tokenId}
                      </div>
                      <div className="text-lg font-bold text-red-600">
                        {formatCurrency(portfolioSummary.worstPerformer.performance.unrealizedPnL)} ETH
                      </div>
                      <div className="text-sm text-red-600">
                        {formatPercentage(portfolioSummary.worstPerformer.performance.unrealizedPnLPercent)}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="risk" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Risk Overview</CardTitle>
                <CardDescription>Portfolio risk assessment</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-3xl font-bold">{riskAnalysis.riskScore}/100</div>
                    <Badge className={getRiskColor(riskAnalysis.overallRisk)} variant="secondary">
                      {riskAnalysis.overallRisk.toUpperCase()} RISK
                    </Badge>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span>Liquidity Risk</span>
                      <div className="flex items-center gap-2">
                        <Progress value={riskAnalysis.liquidityRisk} className="w-20 h-2" />
                        <span className="text-sm font-medium">{riskAnalysis.liquidityRisk}/100</span>
                      </div>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>Concentration Risk</span>
                      <div className="flex items-center gap-2">
                        <Progress value={riskAnalysis.concentrationRisk} className="w-20 h-2" />
                        <span className="text-sm font-medium">{formatPercentage(riskAnalysis.concentrationRisk)}</span>
                      </div>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>Market Risk</span>
                      <div className="flex items-center gap-2">
                        <Progress value={riskAnalysis.marketRisk} className="w-20 h-2" />
                        <span className="text-sm font-medium">{riskAnalysis.marketRisk}/100</span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Diversification Analysis</CardTitle>
                <CardDescription>Portfolio diversification metrics</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-3xl font-bold">{diversificationAnalysis.diversificationScore}/100</div>
                    <div className="text-sm text-muted-foreground">Diversification Score</div>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex justify-between">
                      <span>Top Collection Weight</span>
                      <span className="font-medium">{formatPercentage(diversificationAnalysis.topCollectionWeight)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Herfindahl Index</span>
                      <span className="font-medium">{diversificationAnalysis.herfindahlIndex.toFixed(3)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Concentration Risk</span>
                      <span className="font-medium">{formatPercentage(diversificationAnalysis.concentrationRisk)}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Risk Recommendations */}
          {(riskAnalysis.recommendations.length > 0 || diversificationAnalysis.recommendations.length > 0) && (
            <Card>
              <CardHeader>
                <CardTitle>Risk Management Recommendations</CardTitle>
                <CardDescription>Suggestions to improve your portfolio risk profile</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {riskAnalysis.recommendations.map((recommendation, index) => (
                    <div key={index} className="flex items-start gap-2 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                      <AlertTriangle className="h-4 w-4 text-yellow-600 mt-0.5" />
                      <span className="text-sm">{recommendation}</span>
                    </div>
                  ))}
                  {diversificationAnalysis.recommendations.map((recommendation, index) => (
                    <div key={`div_${index}`} className="flex items-start gap-2 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                      <Target className="h-4 w-4 text-blue-600 mt-0.5" />
                      <div className="text-sm">
                        <div className="font-medium">{recommendation.description}</div>
                        <div className="text-muted-foreground mt-1">
                          Expected Impact: {formatPercentage(recommendation.expectedImpact)} improvement
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
