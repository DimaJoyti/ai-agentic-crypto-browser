'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  TrendingUp, 
  TrendingDown, 
  Activity, 
  DollarSign, 
  Users, 
  ShoppingCart,
  BarChart3,
  PieChart,
  RefreshCw,
  ExternalLink,
  AlertTriangle
} from 'lucide-react'
import { useMarketplaceAPI, useMarketplaceStats } from '@/hooks/useMarketplaceAPI'
import { marketplaceDataAggregator, type AggregatedMarketData } from '@/lib/marketplace-aggregator'
import { formatCurrency, formatNumber, formatPercentage } from '@/lib/utils'

interface MarketplaceComparisonProps {
  marketplaces?: string[]
  className?: string
}

export function MarketplaceComparison({ 
  marketplaces = ['opensea', 'looksrare', 'x2y2', 'blur'],
  className 
}: MarketplaceComparisonProps) {
  const [aggregatedData, setAggregatedData] = useState<AggregatedMarketData | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastUpdate, setLastUpdate] = useState<string | null>(null)

  const { stats, loadStats } = useMarketplaceStats()

  // Load marketplace data
  const loadMarketplaceData = async (forceRefresh = false) => {
    setIsLoading(true)
    setError(null)

    try {
      // Load individual marketplace stats
      await loadStats(marketplaces)

      // Aggregate data from all marketplaces
      const data = await marketplaceDataAggregator.aggregateMarketData(marketplaces, forceRefresh)
      setAggregatedData(data)
      setLastUpdate(new Date().toLocaleString())
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadMarketplaceData()
  }, [marketplaces])

  const handleRefresh = () => {
    loadMarketplaceData(true)
  }

  if (error) {
    return (
      <Alert className="border-red-200 bg-red-50">
        <AlertTriangle className="h-4 w-4 text-red-600" />
        <AlertDescription className="text-red-800">
          Failed to load marketplace data: {error}
        </AlertDescription>
      </Alert>
    )
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Marketplace Comparison</h2>
          <p className="text-muted-foreground">
            Compare performance across major NFT marketplaces
          </p>
        </div>
        <div className="flex items-center gap-2">
          {lastUpdate && (
            <span className="text-sm text-muted-foreground">
              Last updated: {lastUpdate}
            </span>
          )}
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isLoading}
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="volume">Volume</TabsTrigger>
          <TabsTrigger value="market-share">Market Share</TabsTrigger>
          <TabsTrigger value="price-comparison">Price Comparison</TabsTrigger>
          <TabsTrigger value="trends">Trends</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          {/* Global Stats */}
          {aggregatedData && (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Volume (24h)</CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatCurrency(aggregatedData.globalStats.totalVolume24h)}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    <TrendingUp className="inline h-3 w-3 mr-1" />
                    +{formatPercentage(aggregatedData.globalStats.growthMetrics.volumeGrowth24h)} from yesterday
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Sales (24h)</CardTitle>
                  <ShoppingCart className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatNumber(aggregatedData.globalStats.totalSales24h)}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    <TrendingUp className="inline h-3 w-3 mr-1" />
                    +{formatPercentage(aggregatedData.globalStats.growthMetrics.salesGrowth24h)} from yesterday
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Active Listings</CardTitle>
                  <Activity className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatNumber(aggregatedData.globalStats.totalActiveListings)}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    <TrendingUp className="inline h-3 w-3 mr-1" />
                    +{formatPercentage(aggregatedData.globalStats.growthMetrics.listingGrowth24h)} from yesterday
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Unique Traders (24h)</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatNumber(aggregatedData.globalStats.totalUniqueTraders24h)}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    <TrendingUp className="inline h-3 w-3 mr-1" />
                    +{formatPercentage(aggregatedData.globalStats.growthMetrics.traderGrowth24h)} from yesterday
                  </p>
                </CardContent>
              </Card>
            </div>
          )}

          {/* Marketplace Comparison Table */}
          <Card>
            <CardHeader>
              <CardTitle>Marketplace Performance</CardTitle>
              <CardDescription>
                Compare key metrics across marketplaces
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {marketplaces.map(marketplace => {
                  const marketplaceStats = stats[marketplace]
                  const dominance = aggregatedData?.globalStats.dominanceIndex[marketplace] || 0

                  if (!marketplaceStats) {
                    return (
                      <div key={marketplace} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-gray-200 rounded-full animate-pulse" />
                          <div>
                            <div className="font-medium capitalize">{marketplace}</div>
                            <div className="text-sm text-muted-foreground">Loading...</div>
                          </div>
                        </div>
                      </div>
                    )
                  }

                  return (
                    <div key={marketplace} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-sm font-bold">
                          {marketplace.charAt(0).toUpperCase()}
                        </div>
                        <div>
                          <div className="font-medium capitalize">{marketplace}</div>
                          <div className="text-sm text-muted-foreground">
                            {formatPercentage(dominance)} market share
                          </div>
                        </div>
                      </div>

                      <div className="grid grid-cols-4 gap-8 text-right">
                        <div>
                          <div className="text-sm font-medium">
                            {formatCurrency(marketplaceStats.volume24h)}
                          </div>
                          <div className="text-xs text-muted-foreground">Volume 24h</div>
                        </div>
                        <div>
                          <div className="text-sm font-medium">
                            {formatNumber(marketplaceStats.sales24h)}
                          </div>
                          <div className="text-xs text-muted-foreground">Sales 24h</div>
                        </div>
                        <div>
                          <div className="text-sm font-medium">
                            {formatCurrency(marketplaceStats.averagePrice24h)}
                          </div>
                          <div className="text-xs text-muted-foreground">Avg Price</div>
                        </div>
                        <div>
                          <div className="text-sm font-medium">
                            {formatNumber(marketplaceStats.uniqueTraders24h)}
                          </div>
                          <div className="text-xs text-muted-foreground">Traders</div>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="volume" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Volume Comparison</CardTitle>
              <CardDescription>
                24-hour trading volume across marketplaces
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {marketplaces.map(marketplace => {
                  const marketplaceStats = stats[marketplace]
                  const totalVolume = aggregatedData?.globalStats.totalVolume24h || 1
                  const percentage = marketplaceStats ? (marketplaceStats.volume24h / totalVolume) * 100 : 0

                  return (
                    <div key={marketplace} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="font-medium capitalize">{marketplace}</span>
                        <span className="text-sm text-muted-foreground">
                          {formatCurrency(marketplaceStats?.volume24h || 0)}
                        </span>
                      </div>
                      <Progress value={percentage} className="h-2" />
                      <div className="text-xs text-muted-foreground text-right">
                        {formatPercentage(percentage)} of total volume
                      </div>
                    </div>
                  )
                })}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="market-share" className="space-y-4">
          {aggregatedData && (
            <div className="grid gap-4 md:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle>Volume Market Share</CardTitle>
                  <CardDescription>
                    Distribution of trading volume
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {aggregatedData.marketShare.map(data => (
                      <div key={data.marketplace} className="flex items-center justify-between">
                        <span className="capitalize">{data.marketplace}</span>
                        <div className="flex items-center gap-2">
                          <Progress value={data.volumeShare} className="w-20 h-2" />
                          <span className="text-sm font-medium w-12 text-right">
                            {formatPercentage(data.volumeShare)}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Sales Market Share</CardTitle>
                  <CardDescription>
                    Distribution of transaction count
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {aggregatedData.marketShare.map(data => (
                      <div key={data.marketplace} className="flex items-center justify-between">
                        <span className="capitalize">{data.marketplace}</span>
                        <div className="flex items-center gap-2">
                          <Progress value={data.salesShare} className="w-20 h-2" />
                          <span className="text-sm font-medium w-12 text-right">
                            {formatPercentage(data.salesShare)}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="price-comparison" className="space-y-4">
          {aggregatedData && aggregatedData.priceComparisons.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Price Arbitrage Opportunities</CardTitle>
                <CardDescription>
                  Collections with significant price differences across marketplaces
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {aggregatedData.priceComparisons.slice(0, 10).map(comparison => (
                    <div key={comparison.collection.id} className="p-4 border rounded-lg">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <img 
                            src={comparison.collection.imageUrl} 
                            alt={comparison.collection.name}
                            className="w-10 h-10 rounded-lg object-cover"
                          />
                          <div>
                            <div className="font-medium">{comparison.collection.name}</div>
                            <div className="text-sm text-muted-foreground">
                              {formatPercentage(comparison.priceSpread)} price spread
                            </div>
                          </div>
                        </div>
                        <Badge variant={comparison.arbitrageOpportunity > 10 ? "destructive" : "secondary"}>
                          {formatPercentage(comparison.arbitrageOpportunity)} opportunity
                        </Badge>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        {Object.entries(comparison.marketplacePrices).map(([marketplace, price]) => (
                          <div key={marketplace} className="text-center">
                            <div className="font-medium capitalize">{marketplace}</div>
                            <div className={`text-lg ${price === comparison.bestPrice ? 'text-green-600 font-bold' : ''}`}>
                              {formatCurrency(price)} ETH
                            </div>
                            {marketplace === comparison.recommendedMarketplace && (
                              <Badge variant="outline" className="text-xs mt-1">Best Price</Badge>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="trends" className="space-y-4">
          {aggregatedData && (
            <div className="grid gap-4 md:grid-cols-2">
              {aggregatedData.trends.map(trend => (
                <Card key={trend.metric}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      {trend.direction === 'up' ? (
                        <TrendingUp className="h-4 w-4 text-green-600" />
                      ) : trend.direction === 'down' ? (
                        <TrendingDown className="h-4 w-4 text-red-600" />
                      ) : (
                        <BarChart3 className="h-4 w-4 text-gray-600" />
                      )}
                      {trend.metric} Trend
                    </CardTitle>
                    <CardDescription>
                      {trend.timeframe} • {formatPercentage(trend.confidence)} confidence
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold mb-2">
                      {trend.direction === 'up' ? '+' : trend.direction === 'down' ? '-' : ''}
                      {formatPercentage(trend.magnitude)}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {trend.direction === 'up' ? 'Increasing' : trend.direction === 'down' ? 'Decreasing' : 'Stable'} over {trend.timeframe}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
