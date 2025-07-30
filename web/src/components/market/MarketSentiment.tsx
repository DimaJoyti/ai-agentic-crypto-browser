'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { 
  TrendingUp,
  TrendingDown,
  Activity,
  Brain,
  MessageSquare,
  Newspaper,
  BarChart3,
  Thermometer,
  Eye,
  Users,
  Zap
} from 'lucide-react'
import { useMarketAnalytics } from '@/hooks/useMarketAnalytics'
import { type MarketSentiment as MarketSentimentType } from '@/lib/market-analytics'
import { cn } from '@/lib/utils'

interface MarketSentimentProps {
  symbols?: string[]
  showOverall?: boolean
  showIndividual?: boolean
}

export function MarketSentiment({ 
  symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
  showOverall = true,
  showIndividual = true
}: MarketSentimentProps) {
  const [selectedSymbol, setSelectedSymbol] = useState<string | null>(null)

  const {
    state,
    getOverallMarketSentiment,
    getMarketSummary
  } = useMarketAnalytics({
    symbols,
    autoUpdate: true
  })

  const overallSentiment = getOverallMarketSentiment()
  const marketSummary = getMarketSummary()

  const getSentimentColor = (sentiment: string) => {
    switch (sentiment) {
      case 'extremely_bullish':
        return 'text-green-700 dark:text-green-300'
      case 'bullish':
        return 'text-green-600 dark:text-green-400'
      case 'neutral':
        return 'text-gray-600 dark:text-gray-400'
      case 'bearish':
        return 'text-red-600 dark:text-red-400'
      case 'extremely_bearish':
        return 'text-red-700 dark:text-red-300'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const getSentimentBgColor = (sentiment: string) => {
    switch (sentiment) {
      case 'extremely_bullish':
        return 'bg-green-100 dark:bg-green-900'
      case 'bullish':
        return 'bg-green-50 dark:bg-green-900/50'
      case 'neutral':
        return 'bg-gray-50 dark:bg-gray-900/50'
      case 'bearish':
        return 'bg-red-50 dark:bg-red-900/50'
      case 'extremely_bearish':
        return 'bg-red-100 dark:bg-red-900'
      default:
        return 'bg-gray-50 dark:bg-gray-900/50'
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

  const getFearGreedLevel = (index: number) => {
    if (index >= 75) return { level: 'Extreme Greed', color: 'text-red-600' }
    if (index >= 55) return { level: 'Greed', color: 'text-orange-600' }
    if (index >= 45) return { level: 'Neutral', color: 'text-gray-600' }
    if (index >= 25) return { level: 'Fear', color: 'text-yellow-600' }
    return { level: 'Extreme Fear', color: 'text-green-600' }
  }

  const fearGreedLevel = getFearGreedLevel(overallSentiment.fearGreedIndex)

  return (
    <div className="space-y-6">
      {showOverall && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold">Market Sentiment Overview</h3>
          
          {/* Overall Sentiment Cards */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Overall Sentiment</p>
                    <p className={cn("text-xl font-bold", getSentimentColor(overallSentiment.overall))}>
                      {formatSentimentText(overallSentiment.overall)}
                    </p>
                  </div>
                  <div className={cn("w-12 h-12 rounded-lg flex items-center justify-center", 
                    getSentimentBgColor(overallSentiment.overall)
                  )}>
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
                    <p className="text-sm font-medium text-muted-foreground">Fear & Greed</p>
                    <p className={cn("text-xl font-bold", fearGreedLevel.color)}>
                      {Math.round(overallSentiment.fearGreedIndex)}
                    </p>
                  </div>
                  <Thermometer className="w-8 h-8 text-purple-500" />
                </div>
                <div className="mt-2 text-sm text-muted-foreground">
                  {fearGreedLevel.level}
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
                    <p className="text-sm font-medium text-muted-foreground">Confidence</p>
                    <p className="text-xl font-bold">{Math.round(overallSentiment.confidence)}%</p>
                  </div>
                  <Brain className="w-8 h-8 text-blue-500" />
                </div>
                <div className="mt-2 text-sm text-muted-foreground">
                  Data reliability
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sentiment Breakdown */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="w-5 h-5" />
                Sentiment Breakdown
              </CardTitle>
              <CardDescription>
                Detailed analysis of market sentiment components
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-6 md:grid-cols-3">
                {/* Technical Sentiment */}
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <Activity className="w-4 h-4" />
                    <span className="font-medium">Technical</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Score</span>
                      <span className={cn(
                        overallSentiment.technicalSentiment > 0 ? 'text-green-600' : 'text-red-600'
                      )}>
                        {overallSentiment.technicalSentiment.toFixed(1)}
                      </span>
                    </div>
                    <Progress 
                      value={((overallSentiment.technicalSentiment + 100) / 200) * 100} 
                      className="h-2" 
                    />
                  </div>
                </div>

                {/* Social Sentiment */}
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <Users className="w-4 h-4" />
                    <span className="font-medium">Social</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Score</span>
                      <span className={cn(
                        overallSentiment.socialSentiment > 0 ? 'text-green-600' : 'text-red-600'
                      )}>
                        {overallSentiment.socialSentiment.toFixed(1)}
                      </span>
                    </div>
                    <Progress 
                      value={((overallSentiment.socialSentiment + 100) / 200) * 100} 
                      className="h-2" 
                    />
                  </div>
                </div>

                {/* News Sentiment */}
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <Newspaper className="w-4 h-4" />
                    <span className="font-medium">News</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Score</span>
                      <span className={cn(
                        overallSentiment.newssentiment > 0 ? 'text-green-600' : 'text-red-600'
                      )}>
                        {overallSentiment.newssentiment.toFixed(1)}
                      </span>
                    </div>
                    <Progress 
                      value={((overallSentiment.newssentiment + 100) / 200) * 100} 
                      className="h-2" 
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {showIndividual && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold">Individual Asset Sentiment</h3>
          
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {symbols.map((symbol, index) => {
              const sentiment = state.sentiment.get(symbol)
              if (!sentiment) return null

              const isSelected = selectedSymbol === symbol

              return (
                <motion.div
                  key={symbol}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <Card 
                    className={cn(
                      "cursor-pointer transition-all hover:shadow-md",
                      isSelected && "ring-2 ring-primary"
                    )}
                    onClick={() => setSelectedSymbol(isSelected ? null : symbol)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center",
                            getSentimentBgColor(sentiment.overall)
                          )}>
                            <span className={getSentimentColor(sentiment.overall)}>
                              {getSentimentIcon(sentiment.overall)}
                            </span>
                          </div>
                          <div>
                            <h4 className="font-medium">{symbol}</h4>
                            <p className="text-sm text-muted-foreground">
                              {formatSentimentText(sentiment.overall)}
                            </p>
                          </div>
                        </div>
                        <Badge 
                          variant={sentiment.score > 0 ? 'default' : sentiment.score < 0 ? 'destructive' : 'secondary'}
                        >
                          {sentiment.score > 0 ? '+' : ''}{sentiment.score.toFixed(0)}
                        </Badge>
                      </div>

                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Fear & Greed</span>
                          <span>{Math.round(sentiment.fearGreedIndex)}</span>
                        </div>
                        <Progress value={sentiment.fearGreedIndex} className="h-2" />
                      </div>

                      {isSelected && (
                        <motion.div
                          initial={{ opacity: 0, height: 0 }}
                          animate={{ opacity: 1, height: 'auto' }}
                          exit={{ opacity: 0, height: 0 }}
                          className="mt-4 pt-4 border-t space-y-3"
                        >
                          <div className="grid grid-cols-3 gap-3 text-center">
                            <div>
                              <p className="text-xs text-muted-foreground">Technical</p>
                              <p className={cn("text-sm font-medium",
                                sentiment.technicalSentiment > 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {sentiment.technicalSentiment.toFixed(0)}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-muted-foreground">Social</p>
                              <p className={cn("text-sm font-medium",
                                sentiment.socialSentiment > 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {sentiment.socialSentiment.toFixed(0)}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-muted-foreground">News</p>
                              <p className={cn("text-sm font-medium",
                                sentiment.newssentiment > 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {sentiment.newssentiment.toFixed(0)}
                              </p>
                            </div>
                          </div>
                          
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>Confidence: {Math.round(sentiment.confidence)}%</span>
                            <span>Updated: {new Date().toLocaleTimeString()}</span>
                          </div>
                        </motion.div>
                      )}
                    </CardContent>
                  </Card>
                </motion.div>
              )
            })}
          </div>
        </div>
      )}

      {/* Market Summary */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Eye className="w-5 h-5" />
            Market Summary
          </CardTitle>
          <CardDescription>
            Key market sentiment metrics and signals
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <div className="text-center p-4 border rounded-lg">
              <p className="text-2xl font-bold">{marketSummary.totalSymbols}</p>
              <p className="text-sm text-muted-foreground">Tracked Assets</p>
            </div>
            <div className="text-center p-4 border rounded-lg">
              <p className="text-2xl font-bold text-green-600">{marketSummary.bullishSignals}</p>
              <p className="text-sm text-muted-foreground">Bullish Signals</p>
            </div>
            <div className="text-center p-4 border rounded-lg">
              <p className="text-2xl font-bold text-red-600">{marketSummary.bearishSignals}</p>
              <p className="text-sm text-muted-foreground">Bearish Signals</p>
            </div>
            <div className="text-center p-4 border rounded-lg">
              <p className="text-2xl font-bold">{marketSummary.highVolumeAssets}</p>
              <p className="text-sm text-muted-foreground">High Volume</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
