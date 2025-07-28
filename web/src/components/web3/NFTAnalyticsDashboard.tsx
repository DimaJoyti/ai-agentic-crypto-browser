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
  Crown,
  Star,
  Zap,
  Target,
  Eye,
  PieChart,
  LineChart,
  Gauge,
  Brain,
  Calculator,
  Award,
  Flame,
  Shield
} from 'lucide-react'
import { useNFTAnalytics } from '@/hooks/useNFTAnalytics'
import { RarityTier, TrendDirection } from '@/lib/nft-analytics'
import { type Address } from 'viem'

interface NFTAnalyticsDashboardProps {
  contractAddress?: Address
  tokenId?: string
}

export function NFTAnalyticsDashboard({ contractAddress, tokenId }: NFTAnalyticsDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const {
    rarityScore,
    floorPriceData,
    volumeData,
    marketTrend,
    nftValuation,
    collectionAnalytics,
    priceHistory,
    topCollections,
    trendingCollections,
    isLoading,
    loadData,
    analyticsMetrics,
    rarityDistribution,
    priceDistribution,
    topSales,
    marketMetrics,
    predictions,
    getRarityTierColor,
    getTrendColor,
    getConfidenceColor,
    getLiquidityColor,
    getVolatilityColor,
    formatPrice,
    formatPercentage,
    formatVolume,
    isUpTrend,
    isDownTrend,
    isRareNFT,
    isHealthyMarket,
    isVolatileMarket
  } = useNFTAnalytics({
    contractAddress,
    tokenId,
    autoRefresh: true,
    enableNotifications: true
  })

  const getRarityIcon = (tier: RarityTier) => {
    switch (tier) {
      case RarityTier.MYTHIC:
        return <Flame className="w-4 h-4" />
      case RarityTier.LEGENDARY:
        return <Crown className="w-4 h-4" />
      case RarityTier.EPIC:
        return <Star className="w-4 h-4" />
      case RarityTier.RARE:
        return <Zap className="w-4 h-4" />
      case RarityTier.UNCOMMON:
        return <Target className="w-4 h-4" />
      default:
        return <Shield className="w-4 h-4" />
    }
  }

  const getTrendIcon = (trend: TrendDirection) => {
    switch (trend) {
      case TrendDirection.UP:
        return <TrendingUp className="w-4 h-4" />
      case TrendDirection.DOWN:
        return <TrendingDown className="w-4 h-4" />
      default:
        return <Activity className="w-4 h-4" />
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <BarChart3 className="w-6 h-6" />
            NFT Analytics & Valuation
          </h2>
          <p className="text-muted-foreground">
            Advanced analytics, rarity scoring, and AI-powered valuation for NFT collections
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadData}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Analytics Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Floor Price</p>
                <p className="text-2xl font-bold">{formatPrice(analyticsMetrics.floorPrice)}</p>
                <p className={`text-sm ${getTrendColor(analyticsMetrics.trendDirection)}`}>
                  {formatPercentage(analyticsMetrics.floorPriceChange)} 24h
                </p>
              </div>
              <div className={getTrendColor(analyticsMetrics.trendDirection)}>
                {getTrendIcon(analyticsMetrics.trendDirection)}
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Volume 24h</p>
                <p className="text-2xl font-bold">{formatVolume(analyticsMetrics.volume24h)}</p>
                <p className="text-sm text-muted-foreground">
                  {analyticsMetrics.sales24h} sales
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
                <p className="text-sm font-medium text-muted-foreground">Market Health</p>
                <p className="text-2xl font-bold flex items-center gap-2">
                  <span className={getLiquidityColor(analyticsMetrics.liquidityScore)}>
                    {analyticsMetrics.liquidityScore}
                  </span>
                  <span className="text-sm text-muted-foreground">/100</span>
                </p>
                <p className="text-sm text-muted-foreground">
                  {isHealthyMarket ? 'Healthy' : 'Moderate'} liquidity
                </p>
              </div>
              <Gauge className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Volatility</p>
                <p className="text-2xl font-bold flex items-center gap-2">
                  <span className={getVolatilityColor(analyticsMetrics.volatility)}>
                    {analyticsMetrics.volatility}%
                  </span>
                </p>
                <p className="text-sm text-muted-foreground">
                  {isVolatileMarket ? 'High' : 'Low'} volatility
                </p>
              </div>
              <LineChart className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="rarity">Rarity</TabsTrigger>
          <TabsTrigger value="valuation">Valuation</TabsTrigger>
          <TabsTrigger value="trends">Trends</TabsTrigger>
          <TabsTrigger value="collections">Collections</TabsTrigger>
          <TabsTrigger value="predictions">Predictions</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Collection Analytics */}
          {collectionAnalytics && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Price Distribution</CardTitle>
                  <CardDescription>
                    Distribution of NFT prices in the collection
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {priceDistribution.map((item, index) => (
                      <div key={index} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium">{item.range}</span>
                          <span className="text-sm text-muted-foreground">
                            {item.count} ({item.percentage}%)
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
                  <CardTitle>Rarity Distribution</CardTitle>
                  <CardDescription>
                    Distribution of rarity tiers in the collection
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {rarityDistribution.map((item, index) => (
                      <div key={index} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            {getRarityIcon(item.tier)}
                            <span className="text-sm font-medium capitalize">{item.tier}</span>
                          </div>
                          <span className="text-sm text-muted-foreground">
                            {item.count} ({item.percentage}%)
                          </span>
                        </div>
                        <Progress value={item.percentage} className="h-2" />
                        <div className="text-xs text-muted-foreground">
                          Floor: {formatPrice(item.floorPrice)}
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {/* Top Sales */}
          <Card>
            <CardHeader>
              <CardTitle>Top Sales</CardTitle>
              <CardDescription>
                Highest value sales in the collection
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topSales.map((sale, index) => (
                  <motion.div
                    key={sale.tokenId}
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
                        <p className="text-sm font-medium">#{sale.tokenId}</p>
                        <p className="text-xs text-muted-foreground">
                          {sale.marketplace} • {new Date(sale.timestamp).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-bold text-lg">{formatPrice(sale.price)}</p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="rarity" className="space-y-6">
          {/* NFT Rarity Score */}
          {rarityScore && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Award className="w-5 h-5" />
                  Rarity Analysis
                </CardTitle>
                <CardDescription>
                  Detailed rarity breakdown for NFT #{rarityScore.tokenId}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <div className="text-3xl font-bold mb-2">{rarityScore.rarityScore.toFixed(1)}</div>
                    <div className="text-sm text-muted-foreground">Rarity Score</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold mb-2">#{rarityScore.rarityRank}</div>
                    <div className="text-sm text-muted-foreground">Rank</div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold mb-2">{rarityScore.percentile.toFixed(1)}%</div>
                    <div className="text-sm text-muted-foreground">Percentile</div>
                  </div>
                </div>

                <div className="mt-6">
                  <div className="flex items-center justify-center mb-4">
                    <Badge className={getRarityTierColor(rarityScore.tier)}>
                      {getRarityIcon(rarityScore.tier)}
                      <span className="ml-2 capitalize">{rarityScore.tier}</span>
                    </Badge>
                  </div>

                  <div className="space-y-3">
                    <h4 className="font-medium">Attribute Breakdown</h4>
                    {rarityScore.attributes.map((attr, index) => (
                      <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="text-sm font-medium">{attr.trait_type}</p>
                          <p className="text-xs text-muted-foreground">{attr.value}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-sm font-medium">{attr.percentage}%</p>
                          <p className="text-xs text-muted-foreground">
                            Score: {attr.rarityScore.toFixed(1)}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="valuation" className="space-y-6">
          {/* NFT Valuation */}
          {nftValuation && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Brain className="w-5 h-5" />
                  AI Valuation
                </CardTitle>
                <CardDescription>
                  Machine learning-powered valuation analysis
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <div className="text-center p-6 bg-muted rounded-lg mb-4">
                      <div className="text-3xl font-bold mb-2">
                        {formatPrice(nftValuation.estimatedValue)}
                      </div>
                      <div className="text-sm text-muted-foreground">Estimated Value</div>
                      <div className={`text-sm mt-1 ${getConfidenceColor(nftValuation.confidence)}`}>
                        {nftValuation.confidence}% confidence
                      </div>
                    </div>

                    <div className="space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Method</span>
                        <Badge variant="outline" className="capitalize">
                          {nftValuation.valuationMethod.replace('_', ' ')}
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Low Estimate</span>
                        <span className="font-medium">{formatPrice(nftValuation.priceRange.low)}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">High Estimate</span>
                        <span className="font-medium">{formatPrice(nftValuation.priceRange.high)}</span>
                      </div>
                    </div>
                  </div>

                  <div>
                    <h4 className="font-medium mb-3">Valuation Factors</h4>
                    <div className="space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Rarity Score</span>
                        <span className="font-medium">{nftValuation.factors.rarityScore}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Floor Price Multiplier</span>
                        <span className="font-medium">{nftValuation.factors.floorPriceMultiplier}x</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Recent Sales Weight</span>
                        <span className="font-medium">{(nftValuation.factors.recentSalesWeight * 100).toFixed(0)}%</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Market Trend Weight</span>
                        <span className="font-medium">{(nftValuation.factors.marketTrendWeight * 100).toFixed(0)}%</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-muted-foreground">Liquidity Factor</span>
                        <span className="font-medium">{(nftValuation.factors.liquidityFactor * 100).toFixed(0)}%</span>
                      </div>
                    </div>

                    <div className="mt-6">
                      <h4 className="font-medium mb-3">Comparable Sales</h4>
                      <div className="space-y-2">
                        {nftValuation.comparableSales.map((sale, index) => (
                          <div key={index} className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">#{sale.tokenId}</span>
                            <span className="font-medium">{formatPrice(sale.price)}</span>
                            <span className="text-muted-foreground">
                              {(sale.similarityScore * 100).toFixed(0)}% similar
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="trends" className="space-y-6">
          {/* Market Trends */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Market Metrics</CardTitle>
                <CardDescription>
                  Key market health indicators
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Liquidity Score</span>
                      <span className={`font-medium ${getLiquidityColor(marketMetrics.liquidityScore)}`}>
                        {marketMetrics.liquidityScore}/100
                      </span>
                    </div>
                    <Progress value={marketMetrics.liquidityScore} className="h-2" />
                  </div>

                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Volatility Index</span>
                      <span className={`font-medium ${getVolatilityColor(marketMetrics.volatilityIndex)}`}>
                        {marketMetrics.volatilityIndex}%
                      </span>
                    </div>
                    <Progress value={marketMetrics.volatilityIndex} className="h-2" />
                  </div>

                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Trend Score</span>
                      <span className="font-medium">{marketMetrics.trendScore}/100</span>
                    </div>
                    <Progress value={marketMetrics.trendScore} className="h-2" />
                  </div>

                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Momentum</span>
                      <span className="font-medium">{marketMetrics.momentumIndicator}</span>
                    </div>
                    <Progress 
                      value={Math.abs(marketMetrics.momentumIndicator)} 
                      className="h-2" 
                    />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Price History</CardTitle>
                <CardDescription>
                  Recent transaction history
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {priceHistory.map((transaction, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                      <div>
                        <p className="text-sm font-medium capitalize">{transaction.type}</p>
                        <p className="text-xs text-muted-foreground">
                          {transaction.marketplace} • {new Date(transaction.timestamp).toLocaleDateString()}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{formatPrice(transaction.price)}</p>
                        <Button variant="ghost" size="sm" className="h-auto p-0">
                          <Eye className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="collections" className="space-y-6">
          {/* Top Collections */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Top Collections</CardTitle>
                <CardDescription>
                  Highest value collections by market cap
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {topCollections.map((collection, index) => (
                    <motion.div
                      key={collection.contractAddress}
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
                          <p className="text-sm font-medium">{collection.name}</p>
                          <p className="text-xs text-muted-foreground">
                            Floor: {formatPrice(collection.floorPrice)}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{formatVolume(collection.volume24h)}</p>
                        <p className="text-xs text-muted-foreground">24h volume</p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Trending Collections</CardTitle>
                <CardDescription>
                  Collections with highest volume growth
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {trendingCollections.map((collection, index) => (
                    <motion.div
                      key={collection.contractAddress}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center">
                          <TrendingUp className="w-4 h-4 text-green-500" />
                        </div>
                        <div>
                          <p className="text-sm font-medium">{collection.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {collection.sales24h} sales
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-medium">{formatPrice(collection.averagePrice)}</p>
                        <p className="text-xs text-muted-foreground">avg price</p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="predictions" className="space-y-6">
          {/* Price Predictions */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calculator className="w-5 h-5" />
                Price Predictions
              </CardTitle>
              <CardDescription>
                AI-powered price forecasts based on market trends
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold mb-2">
                    {formatPrice(predictions.floorPrice7d)}
                  </div>
                  <div className="text-sm text-muted-foreground">7-day Prediction</div>
                  <div className="text-xs text-green-600 mt-1">
                    {formatPercentage(((parseFloat(predictions.floorPrice7d) / parseFloat(analyticsMetrics.floorPrice) - 1) * 100).toString())}
                  </div>
                </div>

                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold mb-2">
                    {formatPrice(predictions.floorPrice30d)}
                  </div>
                  <div className="text-sm text-muted-foreground">30-day Prediction</div>
                  <div className="text-xs text-green-600 mt-1">
                    {formatPercentage(((parseFloat(predictions.floorPrice30d) / parseFloat(analyticsMetrics.floorPrice) - 1) * 100).toString())}
                  </div>
                </div>

                <div className="text-center p-4 border rounded-lg">
                  <div className="text-2xl font-bold mb-2">{predictions.confidence}%</div>
                  <div className="text-sm text-muted-foreground">Confidence</div>
                  <div className={`text-xs mt-1 ${getConfidenceColor(predictions.confidence)}`}>
                    {predictions.confidence >= 80 ? 'High' : predictions.confidence >= 60 ? 'Medium' : 'Low'}
                  </div>
                </div>
              </div>

              <div className="mt-6 p-4 bg-muted rounded-lg">
                <h4 className="font-medium mb-2">Prediction Factors</h4>
                <div className="text-sm text-muted-foreground space-y-1">
                  <p>• Historical price patterns and market cycles</p>
                  <p>• Current market sentiment and trading volume</p>
                  <p>• Collection rarity distribution and holder behavior</p>
                  <p>• Broader NFT market trends and correlations</p>
                  <p>• Upcoming events and roadmap milestones</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
