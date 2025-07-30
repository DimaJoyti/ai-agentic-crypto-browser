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
  Minus,
  Activity,
  BarChart3,
  Target,
  Zap,
  RefreshCw,
  Info
} from 'lucide-react'
import { useSymbolAnalytics } from '@/hooks/useMarketAnalytics'
import { type MarketIndicator } from '@/lib/market-analytics'
import { cn } from '@/lib/utils'

interface TechnicalIndicatorsProps {
  symbol: string
  compact?: boolean
}

export function TechnicalIndicators({ symbol, compact = false }: TechnicalIndicatorsProps) {
  const [selectedIndicator, setSelectedIndicator] = useState<string | null>(null)

  const {
    indicators,
    marketIndicators,
    isLoading,
    lastUpdate
  } = useSymbolAnalytics(symbol)

  const getSignalColor = (signal: string) => {
    switch (signal) {
      case 'bullish':
        return 'text-green-600 dark:text-green-400'
      case 'bearish':
        return 'text-red-600 dark:text-red-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const getSignalIcon = (signal: string) => {
    switch (signal) {
      case 'bullish':
        return <TrendingUp className="w-4 h-4" />
      case 'bearish':
        return <TrendingDown className="w-4 h-4" />
      default:
        return <Minus className="w-4 h-4" />
    }
  }

  const getSignalBadge = (signal: string) => {
    switch (signal) {
      case 'bullish':
        return <Badge variant="default" className="bg-green-600">Bullish</Badge>
      case 'bearish':
        return <Badge variant="destructive">Bearish</Badge>
      default:
        return <Badge variant="secondary">Neutral</Badge>
    }
  }

  const getRSILevel = (rsi: number) => {
    if (rsi >= 70) return { level: 'Overbought', color: 'text-red-600' }
    if (rsi <= 30) return { level: 'Oversold', color: 'text-green-600' }
    return { level: 'Neutral', color: 'text-gray-600' }
  }

  const formatValue = (value: number, decimals = 2) => {
    return value.toFixed(decimals)
  }

  if (isLoading) {
    return (
      <Card className="w-full">
        <CardContent className="p-6">
          <div className="flex items-center justify-center space-x-2">
            <RefreshCw className="w-4 h-4 animate-spin" />
            <span className="text-sm text-muted-foreground">Loading indicators...</span>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (!indicators || !marketIndicators) {
    return (
      <Card className="w-full">
        <CardContent className="p-6">
          <div className="text-center text-muted-foreground">
            <Activity className="w-8 h-8 mx-auto mb-2" />
            <p>No technical indicators available</p>
            <p className="text-sm">Collecting data for {symbol}...</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (compact) {
    return (
      <div className="space-y-3">
        {marketIndicators.slice(0, 3).map((indicator, index) => (
          <motion.div
            key={indicator.name}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.1 }}
            className="flex items-center justify-between p-3 border rounded-lg"
          >
            <div className="flex items-center space-x-3">
              <div className={cn("w-8 h-8 rounded-full flex items-center justify-center", 
                indicator.signal === 'bullish' ? 'bg-green-100 dark:bg-green-900' :
                indicator.signal === 'bearish' ? 'bg-red-100 dark:bg-red-900' :
                'bg-gray-100 dark:bg-gray-900'
              )}>
                <span className={cn("text-sm font-bold", getSignalColor(indicator.signal))}>
                  {getSignalIcon(indicator.signal)}
                </span>
              </div>
              <div>
                <p className="font-medium">{indicator.name}</p>
                <p className="text-sm text-muted-foreground">
                  {formatValue(indicator.value)}
                </p>
              </div>
            </div>
            {getSignalBadge(indicator.signal)}
          </motion.div>
        ))}
      </div>
    )
  }

  const rsiLevel = getRSILevel(indicators.rsi)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">Technical Indicators</h3>
          <p className="text-sm text-muted-foreground">
            Real-time technical analysis for {symbol}
          </p>
        </div>
        {lastUpdate && (
          <div className="text-xs text-muted-foreground">
            Updated: {new Date(lastUpdate).toLocaleTimeString()}
          </div>
        )}
      </div>

      {/* Market Indicators Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {marketIndicators.map((indicator, index) => (
          <motion.div
            key={indicator.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Card className={cn(
              "cursor-pointer transition-all hover:shadow-md",
              selectedIndicator === indicator.name && "ring-2 ring-primary"
            )}>
              <CardContent 
                className="p-4"
                onClick={() => setSelectedIndicator(
                  selectedIndicator === indicator.name ? null : indicator.name
                )}
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-2">
                    <div className={cn("w-8 h-8 rounded-lg flex items-center justify-center",
                      indicator.signal === 'bullish' ? 'bg-green-100 dark:bg-green-900' :
                      indicator.signal === 'bearish' ? 'bg-red-100 dark:bg-red-900' :
                      'bg-gray-100 dark:bg-gray-900'
                    )}>
                      {getSignalIcon(indicator.signal)}
                    </div>
                    <h4 className="font-medium">{indicator.name}</h4>
                  </div>
                  {getSignalBadge(indicator.signal)}
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Value</span>
                    <span className="font-mono text-sm">{formatValue(indicator.value)}</span>
                  </div>
                  
                  <div className="space-y-1">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-muted-foreground">Strength</span>
                      <span className="text-sm">{Math.round(indicator.strength)}%</span>
                    </div>
                    <Progress value={indicator.strength} className="h-2" />
                  </div>
                </div>

                {selectedIndicator === indicator.name && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="mt-3 pt-3 border-t"
                  >
                    <p className="text-xs text-muted-foreground">
                      {indicator.description}
                    </p>
                  </motion.div>
                )}
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>

      {/* Detailed Indicators */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* RSI */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target className="w-5 h-5" />
              RSI (14)
            </CardTitle>
            <CardDescription>Relative Strength Index</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-2xl font-bold">{formatValue(indicators.rsi)}</span>
                <Badge variant={rsiLevel.level === 'Neutral' ? 'secondary' : 'default'}>
                  {rsiLevel.level}
                </Badge>
              </div>
              
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Oversold</span>
                  <span>Overbought</span>
                </div>
                <div className="relative">
                  <Progress value={(indicators.rsi / 100) * 100} className="h-3" />
                  <div className="absolute top-0 left-[30%] w-0.5 h-3 bg-green-500" />
                  <div className="absolute top-0 left-[70%] w-0.5 h-3 bg-red-500" />
                </div>
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span>30</span>
                  <span>70</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* MACD */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="w-5 h-5" />
              MACD
            </CardTitle>
            <CardDescription>Moving Average Convergence Divergence</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="grid grid-cols-3 gap-4 text-center">
                <div>
                  <p className="text-xs text-muted-foreground">MACD</p>
                  <p className="font-mono text-sm">{formatValue(indicators.macd.macd, 4)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Signal</p>
                  <p className="font-mono text-sm">{formatValue(indicators.macd.signal, 4)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Histogram</p>
                  <p className={cn("font-mono text-sm", 
                    indicators.macd.histogram > 0 ? 'text-green-600' : 'text-red-600'
                  )}>
                    {formatValue(indicators.macd.histogram, 4)}
                  </p>
                </div>
              </div>
              
              <div className="text-center">
                {indicators.macd.macd > indicators.macd.signal ? (
                  <Badge variant="default" className="bg-green-600">
                    <TrendingUp className="w-3 h-3 mr-1" />
                    Bullish Crossover
                  </Badge>
                ) : (
                  <Badge variant="destructive">
                    <TrendingDown className="w-3 h-3 mr-1" />
                    Bearish Crossover
                  </Badge>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Moving Averages */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="w-5 h-5" />
              Moving Averages
            </CardTitle>
            <CardDescription>Simple and Exponential Moving Averages</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-muted-foreground">SMA 20</p>
                  <p className="font-mono text-sm">{formatValue(indicators.movingAverages.sma20)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">SMA 50</p>
                  <p className="font-mono text-sm">{formatValue(indicators.movingAverages.sma50)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">EMA 12</p>
                  <p className="font-mono text-sm">{formatValue(indicators.movingAverages.ema12)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">EMA 26</p>
                  <p className="font-mono text-sm">{formatValue(indicators.movingAverages.ema26)}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Bollinger Bands */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Zap className="w-5 h-5" />
              Bollinger Bands
            </CardTitle>
            <CardDescription>Price volatility and trend analysis</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="grid grid-cols-3 gap-4 text-center">
                <div>
                  <p className="text-xs text-muted-foreground">Upper</p>
                  <p className="font-mono text-sm">{formatValue(indicators.bollinger.upper)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Middle</p>
                  <p className="font-mono text-sm">{formatValue(indicators.bollinger.middle)}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Lower</p>
                  <p className="font-mono text-sm">{formatValue(indicators.bollinger.lower)}</p>
                </div>
              </div>
              
              <div className="relative h-2 bg-gray-200 dark:bg-gray-700 rounded">
                <div 
                  className="absolute top-0 left-0 h-full bg-blue-500 rounded"
                  style={{ 
                    width: '100%',
                    background: 'linear-gradient(to right, #ef4444, #3b82f6, #ef4444)'
                  }}
                />
              </div>
              
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>Lower Band</span>
                <span>Upper Band</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
