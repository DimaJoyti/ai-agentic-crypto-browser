'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  BarChart3,
  TrendingUp,
  TrendingDown,
  Activity,
  Brain,
  Target,
  Zap,
  RefreshCw,
  Settings,
  AlertCircle,
  Eye,
  Volume2,
  Thermometer
} from 'lucide-react'
import { TechnicalIndicators } from './TechnicalIndicators'
import { MarketSentiment } from './MarketSentiment'
import { useMarketAnalytics } from '@/hooks/useMarketAnalytics'
import { cn } from '@/lib/utils'

export function MarketAnalyticsDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedSymbol, setSelectedSymbol] = useState('BTC')
  const [watchlist] = useState(['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'])

  const {
    state,
    getOverallMarketSentiment,
    getTopPerformers,
    getMarketSummary,
    refresh
  } = useMarketAnalytics({
    symbols: watchlist,
    autoUpdate: true
  })

  const overallSentiment = getOverallMarketSentiment()
  const marketSummary = getMarketSummary()
  const topRSI = getTopPerformers('rsi')
  const topVolume = getTopPerformers('volume')
  const topSentiment = getTopPerformers('sentiment')

  const getSentimentColor = (sentiment: string) => {
    switch (sentiment) {
      case 'extremely_bullish':
      case 'bullish':
        return 'text-green-600 dark:text-green-400'
      case 'extremely_bearish':
      case 'bearish':
        return 'text-red-600 dark:text-red-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const getSentimentIcon = (sentiment: string) => {
    switch (sentiment) {
      case 'extremely_bullish':
      case 'bullish':
        return <TrendingUp className="w-5 h-5" />
      case 'extremely_bearish':
      case 'bearish':
        return <TrendingDown className="w-5 h-5" />
      default:
        return <Activity className="w-5 h-5" />
    }
  }

  const formatSentimentText = (sentiment: string) => {
    return sentiment.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Market Analytics</h1>
          <p className="text-muted-foreground">
            Advanced technical analysis and market sentiment tracking
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={refresh}
            disabled={state.isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Market Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Market Sentiment</p>
                <p className={cn("text-xl font-bold", getSentimentColor(overallSentiment.overall))}>
                  {formatSentimentText(overallSentiment.overall)}
                </p>
              </div>
              <div className={cn("w-12 h-12 rounded-lg flex items-center justify-center bg-blue-100 dark:bg-blue-900")}>
                <span className={getSentimentColor(overallSentiment.overall)}>
                  {getSentimentIcon(overallSentiment.overall)}
                </span>
              </div>
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Score: {overallSentiment.score.toFixed(1)}/100
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Fear & Greed Index</p>
                <p className="text-xl font-bold">{Math.round(overallSentiment.fearGreedIndex)}</p>
              </div>
              <Thermometer className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {overallSentiment.fearGreedIndex >= 75 ? 'Extreme Greed' :
               overallSentiment.fearGreedIndex >= 55 ? 'Greed' :
               overallSentiment.fearGreedIndex >= 45 ? 'Neutral' :
               overallSentiment.fearGreedIndex >= 25 ? 'Fear' : 'Extreme Fear'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Bullish Signals</p>
                <p className="text-xl font-bold text-green-600">{marketSummary.bullishSignals}</p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              vs {marketSummary.bearishSignals} bearish
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">High Volume Assets</p>
                <p className="text-xl font-bold">{marketSummary.highVolumeAssets}</p>
              </div>
              <Volume2 className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Above average volume
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Performers */}
      <div className="grid gap-6 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target className="w-5 h-5 text-blue-500" />
              Top RSI Signals
            </CardTitle>
            <CardDescription>Assets with notable RSI levels</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topRSI.slice(0, 5).map((item, index) => (
                <motion.div
                  key={item.symbol}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center text-xs font-bold">
                      {index + 1}
                    </span>
                    <span className="font-medium">{item.symbol}</span>
                  </div>
                  <Badge variant={item.value > 70 ? 'destructive' : item.value < 30 ? 'default' : 'secondary'}>
                    {item.value.toFixed(1)}
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Volume2 className="w-5 h-5 text-orange-500" />
              High Volume Assets
            </CardTitle>
            <CardDescription>Assets with elevated trading volume</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topVolume.slice(0, 5).map((item, index) => (
                <motion.div
                  key={item.symbol}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 bg-orange-100 dark:bg-orange-900 rounded-full flex items-center justify-center text-xs font-bold">
                      {index + 1}
                    </span>
                    <span className="font-medium">{item.symbol}</span>
                  </div>
                  <Badge variant={item.value > 1.5 ? 'default' : 'secondary'}>
                    {(item.value * 100).toFixed(0)}%
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Brain className="w-5 h-5 text-purple-500" />
              Sentiment Leaders
            </CardTitle>
            <CardDescription>Assets with strongest sentiment</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topSentiment.slice(0, 5).map((item, index) => (
                <motion.div
                  key={item.symbol}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <div className="flex items-center gap-2">
                    <span className="w-6 h-6 bg-purple-100 dark:bg-purple-900 rounded-full flex items-center justify-center text-xs font-bold">
                      {index + 1}
                    </span>
                    <span className="font-medium">{item.symbol}</span>
                  </div>
                  <Badge variant={item.value > 0 ? 'default' : 'destructive'}>
                    {item.value > 0 ? '+' : ''}{item.value.toFixed(0)}
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Analytics Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="technical">Technical</TabsTrigger>
          <TabsTrigger value="sentiment">Sentiment</TabsTrigger>
          <TabsTrigger value="volume">Volume</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Symbol Selector */}
            <Card>
              <CardHeader>
                <CardTitle>Asset Analysis</CardTitle>
                <CardDescription>Select an asset for detailed technical analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-4 gap-2">
                  {watchlist.map((symbol) => (
                    <Button
                      key={symbol}
                      variant={selectedSymbol === symbol ? "default" : "outline"}
                      size="sm"
                      onClick={() => setSelectedSymbol(symbol)}
                    >
                      {symbol}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Quick Stats */}
            <Card>
              <CardHeader>
                <CardTitle>Market Statistics</CardTitle>
                <CardDescription>Real-time market analytics overview</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4">
                  <div className="text-center p-3 border rounded">
                    <p className="text-lg font-bold">{watchlist.length}</p>
                    <p className="text-sm text-muted-foreground">Tracked Assets</p>
                  </div>
                  <div className="text-center p-3 border rounded">
                    <p className="text-lg font-bold">{state.indicators.size}</p>
                    <p className="text-sm text-muted-foreground">With Indicators</p>
                  </div>
                  <div className="text-center p-3 border rounded">
                    <p className="text-lg font-bold">{Math.round(overallSentiment.confidence)}%</p>
                    <p className="text-sm text-muted-foreground">Confidence</p>
                  </div>
                  <div className="text-center p-3 border rounded">
                    <p className="text-lg font-bold">
                      {state.lastUpdate ? new Date(state.lastUpdate).toLocaleTimeString() : 'N/A'}
                    </p>
                    <p className="text-sm text-muted-foreground">Last Update</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Selected Asset Technical Indicators */}
          <TechnicalIndicators symbol={selectedSymbol} />
        </TabsContent>

        <TabsContent value="technical" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Technical Analysis</h3>
            <div className="flex items-center gap-2">
              {watchlist.map((symbol) => (
                <Button
                  key={symbol}
                  variant={selectedSymbol === symbol ? "default" : "outline"}
                  size="sm"
                  onClick={() => setSelectedSymbol(symbol)}
                >
                  {symbol}
                </Button>
              ))}
            </div>
          </div>
          <TechnicalIndicators symbol={selectedSymbol} />
        </TabsContent>

        <TabsContent value="sentiment" className="space-y-6">
          <MarketSentiment symbols={watchlist} />
        </TabsContent>

        <TabsContent value="volume" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Volume2 className="w-5 h-5" />
                Volume Analysis
              </CardTitle>
              <CardDescription>
                Trading volume analysis and trends
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {watchlist.map((symbol) => {
                  const volumeAnalysis = state.volumeAnalysis.get(symbol)
                  if (!volumeAnalysis) return null

                  return (
                    <Card key={symbol}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between mb-3">
                          <h4 className="font-medium">{symbol}</h4>
                          <Badge variant={
                            volumeAnalysis.volumeRatio > 1.5 ? 'default' :
                            volumeAnalysis.volumeRatio < 0.5 ? 'destructive' :
                            'secondary'
                          }>
                            {(volumeAnalysis.volumeRatio * 100).toFixed(0)}%
                          </Badge>
                        </div>
                        
                        <div className="space-y-2 text-sm">
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">Trend</span>
                            <span className={cn(
                              volumeAnalysis.volumeTrend === 'increasing' ? 'text-green-600' :
                              volumeAnalysis.volumeTrend === 'decreasing' ? 'text-red-600' :
                              'text-gray-600'
                            )}>
                              {volumeAnalysis.volumeTrend}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">Avg Volume</span>
                            <span>{(volumeAnalysis.averageVolume / 1e6).toFixed(1)}M</span>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  )
                })}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
