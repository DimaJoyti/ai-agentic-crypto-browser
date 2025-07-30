'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  TrendingUp,
  TrendingDown,
  Activity,
  DollarSign,
  BarChart3,
  Bell,
  RefreshCw,
  Zap,
  AlertCircle,
  Eye,
  Settings
} from 'lucide-react'
import { PriceTicker } from './PriceTicker'
import { PriceAlerts } from './PriceAlerts'
import { usePriceFeed } from '@/hooks/usePriceFeed'
import { type PriceData } from '@/lib/price-feed-manager'

export function MarketDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedTimeframe, setSelectedTimeframe] = useState('24h')
  const [watchlist] = useState(['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'])

  const {
    state,
    getPrices,
    formatPrice,
    formatChange,
    restart,
    clearError
  } = usePriceFeed({
    symbols: watchlist,
    autoStart: true,
    enableAlerts: true
  })

  const prices = getPrices(watchlist)
  const pricesArray = Array.from(prices.values())

  // Calculate market statistics
  const marketStats = {
    totalMarketCap: pricesArray.reduce((sum, price) => sum + price.marketCap, 0),
    totalVolume: pricesArray.reduce((sum, price) => sum + price.volume24h, 0),
    gainers: pricesArray.filter(price => price.changePercent24h > 0).length,
    losers: pricesArray.filter(price => price.changePercent24h < 0).length,
    avgChange: pricesArray.length > 0 
      ? pricesArray.reduce((sum, price) => sum + price.changePercent24h, 0) / pricesArray.length 
      : 0
  }

  const topGainers = pricesArray
    .filter(price => price.changePercent24h > 0)
    .sort((a, b) => b.changePercent24h - a.changePercent24h)
    .slice(0, 3)

  const topLosers = pricesArray
    .filter(price => price.changePercent24h < 0)
    .sort((a, b) => a.changePercent24h - b.changePercent24h)
    .slice(0, 3)

  const formatMarketCap = (value: number) => {
    if (value >= 1e12) return `$${(value / 1e12).toFixed(2)}T`
    if (value >= 1e9) return `$${(value / 1e9).toFixed(2)}B`
    if (value >= 1e6) return `$${(value / 1e6).toFixed(2)}M`
    return `$${value.toFixed(2)}`
  }

  const getMarketSentiment = () => {
    if (marketStats.avgChange > 2) return { text: 'Bullish', color: 'text-green-600', icon: TrendingUp }
    if (marketStats.avgChange < -2) return { text: 'Bearish', color: 'text-red-600', icon: TrendingDown }
    return { text: 'Neutral', color: 'text-gray-600', icon: Activity }
  }

  const sentiment = getMarketSentiment()
  const SentimentIcon = sentiment.icon

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Market Dashboard</h1>
          <p className="text-muted-foreground">
            Real-time cryptocurrency market data and analytics
          </p>
        </div>
        <div className="flex items-center gap-2">
          {state.isConnected && (
            <Badge variant="default" className="text-xs">
              <Zap className="w-3 h-3 mr-1" />
              Live Data
            </Badge>
          )}
          <Button
            variant="outline"
            size="sm"
            onClick={restart}
            disabled={state.isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Market Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Market Cap</p>
                <p className="text-2xl font-bold">{formatMarketCap(marketStats.totalMarketCap)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              24h Volume: {formatMarketCap(marketStats.totalVolume)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Market Sentiment</p>
                <p className={`text-2xl font-bold ${sentiment.color}`}>{sentiment.text}</p>
              </div>
              <SentimentIcon className={`w-8 h-8 ${sentiment.color}`} />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Avg Change: {formatChange(marketStats.avgChange, true)}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Gainers</p>
                <p className="text-2xl font-bold text-green-600">{marketStats.gainers}</p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Losers: {marketStats.losers}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Alerts</p>
                <p className="text-2xl font-bold">{state.stats.activeAlerts}</p>
              </div>
              <Bell className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Data Sources: {state.stats.activeSources}
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

      {/* Top Movers */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-green-500" />
              Top Gainers
            </CardTitle>
            <CardDescription>Best performing assets in the last 24 hours</CardDescription>
          </CardHeader>
          <CardContent>
            {topGainers.length > 0 ? (
              <div className="space-y-3">
                {topGainers.map((price, index) => (
                  <motion.div
                    key={price.symbol}
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
                        <p className="font-medium">{price.symbol}</p>
                        <p className="text-sm text-muted-foreground">{formatPrice(price.price)}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-green-600 dark:text-green-400 font-medium">
                        +{formatChange(price.changePercent24h, true)}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        +{formatChange(price.change24h)}
                      </p>
                    </div>
                  </motion.div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No gainers in the current timeframe
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
            <CardDescription>Worst performing assets in the last 24 hours</CardDescription>
          </CardHeader>
          <CardContent>
            {topLosers.length > 0 ? (
              <div className="space-y-3">
                {topLosers.map((price, index) => (
                  <motion.div
                    key={price.symbol}
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
                        <p className="font-medium">{price.symbol}</p>
                        <p className="text-sm text-muted-foreground">{formatPrice(price.price)}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-red-600 dark:text-red-400 font-medium">
                        {formatChange(price.changePercent24h, true)}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {formatChange(price.change24h)}
                      </p>
                    </div>
                  </motion.div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No losers in the current timeframe
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="prices">Live Prices</TabsTrigger>
          <TabsTrigger value="alerts">Price Alerts</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Market Overview</CardTitle>
                <CardDescription>
                  Real-time market data and key metrics for your watchlist
                </CardDescription>
              </CardHeader>
              <CardContent>
                <PriceTicker 
                  symbols={watchlist}
                  showChange={true}
                  showVolume={true}
                  showMarketCap={true}
                  compact={false}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Market Statistics</CardTitle>
                <CardDescription>
                  Key market metrics and performance indicators
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{pricesArray.length}</p>
                    <p className="text-sm text-muted-foreground">Tracked Assets</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{state.stats.activeSources}</p>
                    <p className="text-sm text-muted-foreground">Data Sources</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{state.stats.subscribers}</p>
                    <p className="text-sm text-muted-foreground">Active Subscriptions</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">
                      {state.lastUpdate ? new Date(state.lastUpdate).toLocaleTimeString() : 'N/A'}
                    </p>
                    <p className="text-sm text-muted-foreground">Last Update</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="prices" className="space-y-6">
          <PriceTicker 
            symbols={watchlist}
            showChange={true}
            showVolume={true}
            showMarketCap={true}
            compact={false}
          />
        </TabsContent>

        <TabsContent value="alerts" className="space-y-6">
          <PriceAlerts symbols={watchlist} />
        </TabsContent>
      </Tabs>
    </div>
  )
}
