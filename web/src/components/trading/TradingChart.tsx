'use client'

import React, { useEffect, useRef, useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  BarChart3, 
  TrendingUp, 
  Settings, 
  Maximize2, 
  Volume2,
  Activity,
  Target,
  Layers
} from 'lucide-react'

interface TradingChartProps {
  pair: string
  height?: number
  showToolbar?: boolean
}

interface ChartData {
  time: number
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export function TradingChart({ pair, height = 400, showToolbar = true }: TradingChartProps) {
  const chartContainerRef = useRef<HTMLDivElement>(null)
  const [timeframe, setTimeframe] = useState('1h')
  const [chartType, setChartType] = useState('candlestick')
  const [indicators, setIndicators] = useState<string[]>(['volume'])
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [chartData, setChartData] = useState<ChartData[]>([])

  // Generate mock chart data
  useEffect(() => {
    const generateMockData = () => {
      const data: ChartData[] = []
      let basePrice = 2500 // Starting price for ETH
      const now = Date.now()
      const timeframeMs = getTimeframeMs(timeframe)
      
      for (let i = 100; i >= 0; i--) {
        const time = now - (i * timeframeMs)
        const volatility = 0.02 // 2% volatility
        
        const open = basePrice
        const change = (Math.random() - 0.5) * basePrice * volatility
        const close = open + change
        const high = Math.max(open, close) + Math.random() * basePrice * volatility * 0.5
        const low = Math.min(open, close) - Math.random() * basePrice * volatility * 0.5
        const volume = Math.random() * 1000 + 100
        
        data.push({
          time,
          open,
          high,
          low,
          close,
          volume
        })
        
        basePrice = close
      }
      
      return data
    }

    setChartData(generateMockData())
  }, [pair, timeframe])

  const getTimeframeMs = (tf: string) => {
    const timeframes: Record<string, number> = {
      '1m': 60 * 1000,
      '5m': 5 * 60 * 1000,
      '15m': 15 * 60 * 1000,
      '1h': 60 * 60 * 1000,
      '4h': 4 * 60 * 60 * 1000,
      '1d': 24 * 60 * 60 * 1000,
      '1w': 7 * 24 * 60 * 60 * 1000
    }
    return timeframes[tf] || timeframes['1h']
  }

  const timeframes = [
    { value: '1m', label: '1m' },
    { value: '5m', label: '5m' },
    { value: '15m', label: '15m' },
    { value: '1h', label: '1h' },
    { value: '4h', label: '4h' },
    { value: '1d', label: '1D' },
    { value: '1w', label: '1W' }
  ]

  const chartTypes = [
    { value: 'candlestick', label: 'Candles', icon: BarChart3 },
    { value: 'line', label: 'Line', icon: TrendingUp },
    { value: 'area', label: 'Area', icon: Layers }
  ]

  const availableIndicators = [
    { value: 'volume', label: 'Volume' },
    { value: 'ma', label: 'Moving Average' },
    { value: 'rsi', label: 'RSI' },
    { value: 'macd', label: 'MACD' },
    { value: 'bb', label: 'Bollinger Bands' }
  ]

  const toggleIndicator = (indicator: string) => {
    setIndicators(prev => 
      prev.includes(indicator) 
        ? prev.filter(i => i !== indicator)
        : [...prev, indicator]
    )
  }

  // Mock chart rendering - in a real app, you'd use a charting library like TradingView, Chart.js, or D3
  const renderChart = () => {
    if (chartData.length === 0) return null

    const maxPrice = Math.max(...chartData.map(d => d.high))
    const minPrice = Math.min(...chartData.map(d => d.low))
    const priceRange = maxPrice - minPrice
    const chartHeight = height - 100 // Account for toolbar and padding

    return (
      <div className="relative w-full" style={{ height: chartHeight }}>
        <svg width="100%" height="100%" className="overflow-visible">
          {/* Grid lines */}
          <defs>
            <pattern id="grid" width="50" height="40" patternUnits="userSpaceOnUse">
              <path d="M 50 0 L 0 0 0 40" fill="none" stroke="currentColor" strokeWidth="0.5" opacity="0.1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />
          
          {/* Price chart */}
          {chartType === 'candlestick' && chartData.map((candle, index) => {
            const x = (index / (chartData.length - 1)) * 100
            const bodyTop = ((maxPrice - Math.max(candle.open, candle.close)) / priceRange) * 100
            const bodyBottom = ((maxPrice - Math.min(candle.open, candle.close)) / priceRange) * 100
            const wickTop = ((maxPrice - candle.high) / priceRange) * 100
            const wickBottom = ((maxPrice - candle.low) / priceRange) * 100
            const isGreen = candle.close > candle.open
            
            return (
              <g key={index}>
                {/* Wick */}
                <line
                  x1={`${x}%`}
                  y1={`${wickTop}%`}
                  x2={`${x}%`}
                  y2={`${wickBottom}%`}
                  stroke={isGreen ? '#22c55e' : '#ef4444'}
                  strokeWidth="1"
                />
                {/* Body */}
                <rect
                  x={`${x - 0.4}%`}
                  y={`${bodyTop}%`}
                  width="0.8%"
                  height={`${bodyBottom - bodyTop}%`}
                  fill={isGreen ? '#22c55e' : '#ef4444'}
                  opacity="0.8"
                />
              </g>
            )
          })}
          
          {/* Line chart */}
          {chartType === 'line' && (
            <polyline
              fill="none"
              stroke="#3b82f6"
              strokeWidth="2"
              points={chartData.map((candle, index) => {
                const x = (index / (chartData.length - 1)) * 100
                const y = ((maxPrice - candle.close) / priceRange) * 100
                return `${x},${y}`
              }).join(' ')}
            />
          )}
          
          {/* Area chart */}
          {chartType === 'area' && (
            <polygon
              fill="url(#areaGradient)"
              stroke="#3b82f6"
              strokeWidth="2"
              points={[
                ...chartData.map((candle, index) => {
                  const x = (index / (chartData.length - 1)) * 100
                  const y = ((maxPrice - candle.close) / priceRange) * 100
                  return `${x},${y}`
                }),
                `100,100`,
                `0,100`
              ].join(' ')}
            />
          )}
          
          {/* Gradient for area chart */}
          <defs>
            <linearGradient id="areaGradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="#3b82f6" stopOpacity="0.3"/>
              <stop offset="100%" stopColor="#3b82f6" stopOpacity="0.05"/>
            </linearGradient>
          </defs>
        </svg>
        
        {/* Price labels */}
        <div className="absolute right-0 top-0 bottom-0 w-16 flex flex-col justify-between text-xs text-muted-foreground py-2">
          <span>${maxPrice.toFixed(2)}</span>
          <span>${((maxPrice + minPrice) / 2).toFixed(2)}</span>
          <span>${minPrice.toFixed(2)}</span>
        </div>
      </div>
    )
  }

  return (
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-background' : 'relative'}`}>
      <Card className="h-full">
        {showToolbar && (
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              {/* Timeframe Selector */}
              <Tabs value={timeframe} onValueChange={setTimeframe}>
                <TabsList className="h-8">
                  {timeframes.map(tf => (
                    <TabsTrigger key={tf.value} value={tf.value} className="text-xs px-2">
                      {tf.label}
                    </TabsTrigger>
                  ))}
                </TabsList>
              </Tabs>

              {/* Chart Controls */}
              <div className="flex items-center gap-2">
                {/* Chart Type */}
                <Tabs value={chartType} onValueChange={setChartType}>
                  <TabsList className="h-8">
                    {chartTypes.map(type => {
                      const Icon = type.icon
                      return (
                        <TabsTrigger key={type.value} value={type.value} className="text-xs px-2">
                          <Icon className="w-3 h-3" />
                        </TabsTrigger>
                      )
                    })}
                  </TabsList>
                </Tabs>

                {/* Indicators */}
                <div className="flex items-center gap-1">
                  {availableIndicators.map(indicator => (
                    <Button
                      key={indicator.value}
                      variant={indicators.includes(indicator.value) ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => toggleIndicator(indicator.value)}
                      className="text-xs h-8 px-2"
                    >
                      {indicator.label}
                    </Button>
                  ))}
                </div>

                {/* Settings */}
                <Button variant="outline" size="sm" className="h-8 px-2">
                  <Settings className="w-3 h-3" />
                </Button>

                {/* Fullscreen */}
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={() => setIsFullscreen(!isFullscreen)}
                  className="h-8 px-2"
                >
                  <Maximize2 className="w-3 h-3" />
                </Button>
              </div>
            </div>

            {/* Active Indicators */}
            {indicators.length > 0 && (
              <div className="flex items-center gap-2 mt-2">
                <span className="text-xs text-muted-foreground">Indicators:</span>
                {indicators.map(indicator => (
                  <Badge key={indicator} variant="secondary" className="text-xs">
                    {availableIndicators.find(i => i.value === indicator)?.label}
                  </Badge>
                ))}
              </div>
            )}
          </CardHeader>
        )}

        <CardContent className="p-0">
          <div ref={chartContainerRef} className="w-full">
            {renderChart()}
          </div>

          {/* Volume Chart */}
          {indicators.includes('volume') && (
            <div className="h-20 border-t">
              <div className="p-2">
                <div className="flex items-center gap-2 mb-2">
                  <Volume2 className="w-3 h-3 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">Volume</span>
                </div>
                <div className="h-12 flex items-end gap-px">
                  {chartData.map((candle, index) => {
                    const maxVolume = Math.max(...chartData.map(d => d.volume))
                    const height = (candle.volume / maxVolume) * 100
                    const isGreen = candle.close > candle.open
                    
                    return (
                      <div
                        key={index}
                        className={`flex-1 ${isGreen ? 'bg-green-500' : 'bg-red-500'} opacity-60`}
                        style={{ height: `${height}%` }}
                      />
                    )
                  })}
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
