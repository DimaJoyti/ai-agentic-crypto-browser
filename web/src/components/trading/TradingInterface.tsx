'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { TrendingUp, TrendingDown, Activity, DollarSign, BarChart3, Zap } from 'lucide-react'
import { OrderBook } from './OrderBook'
import { TradingChart } from './TradingChart'
import { TradeHistory } from './TradeHistory'
import { MarketSelector } from './MarketSelector'
import { PriceTickerBar } from './PriceTickerBar'
import { OrderForm } from './OrderForm'
import { PositionManager } from './PositionManager'
import { useAccount } from 'wagmi'
import { useRealTimePrice } from '@/hooks/useRealTimePrice'
import { useWebSocket } from '@/hooks/useWebSocket'

interface TradingPair {
  symbol: string
  baseAsset: string
  quoteAsset: string
  price: number
  change24h: number
  volume24h: number
  high24h: number
  low24h: number
}

interface TradingInterfaceProps {
  initialPair?: string
  embedded?: boolean
}

export function TradingInterface({ initialPair = 'ETH/USDT', embedded = false }: TradingInterfaceProps) {
  const [mounted, setMounted] = useState(false)
  const [selectedPair, setSelectedPair] = useState<TradingPair>({
    symbol: initialPair,
    baseAsset: initialPair.split('/')[0],
    quoteAsset: initialPair.split('/')[1],
    price: 0,
    change24h: 0,
    volume24h: 0,
    high24h: 0,
    low24h: 0
  })

  const [activeTab, setActiveTab] = useState('spot')
  const [orderType, setOrderType] = useState<'market' | 'limit' | 'stop'>('limit')
  const [side, setSide] = useState<'buy' | 'sell'>('buy')

  const { address, isConnected } = useAccount()
  const { price, isLoading: priceLoading } = useRealTimePrice(selectedPair.symbol)
  const { isConnected: wsConnected, lastMessage } = useWebSocket('wss://api.example.com/ws')

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (price) {
      setSelectedPair(prev => ({ ...prev, price }))
    }
  }, [price])

  const handlePairChange = (pair: TradingPair) => {
    setSelectedPair(pair)
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 6
    }).format(price)
  }

  const formatChange = (change: number) => {
    const isPositive = change >= 0
    return (
      <span className={`flex items-center gap-1 ${isPositive ? 'text-green-500' : 'text-red-500'}`}>
        {isPositive ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
        {Math.abs(change).toFixed(2)}%
      </span>
    )
  }

  // Prevent hydration mismatch by not rendering until mounted
  if (!mounted) {
    return (
      <div className="min-h-[600px] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
          <p className="text-sm text-muted-foreground">Loading trading interface...</p>
        </div>
      </div>
    )
  }

  const containerClass = embedded
    ? "space-y-4"
    : "min-h-screen bg-background"

  const contentClass = embedded
    ? "space-y-4"
    : "container mx-auto p-4 space-y-4"

  return (
    <div className={containerClass}>
      {/* Price Ticker Bar - only show when not embedded */}
      {!embedded && <PriceTickerBar />}

      <div className={contentClass}>
        {/* Market Selector - only show when not embedded */}
        {!embedded && (
          <MarketSelector
            selectedPair={selectedPair}
            onPairChange={handlePairChange}
          />
        )}

        {/* Main Trading Interface */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
          {/* Left Column - Order Book */}
          <div className="lg:col-span-1">
            <Card className="h-[600px]">
              <CardHeader className="pb-3">
                <CardTitle className="text-sm font-medium">Order Book</CardTitle>
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>Price ({selectedPair.quoteAsset})</span>
                  <span>Amount ({selectedPair.baseAsset})</span>
                  <span>Total</span>
                </div>
              </CardHeader>
              <CardContent className="p-0">
                <OrderBook symbol={selectedPair.symbol} />
              </CardContent>
            </Card>
          </div>

          {/* Center Column - Chart */}
          <div className="lg:col-span-2">
            <Card className="h-[600px]">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <h2 className="text-xl font-bold">{selectedPair.symbol}</h2>
                    <div className="flex items-center gap-2">
                      <span className="text-2xl font-mono">
                        {formatPrice(selectedPair.price)}
                      </span>
                      {formatChange(selectedPair.change24h)}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={wsConnected ? 'default' : 'destructive'}>
                      <Activity className="w-3 h-3 mr-1" />
                      {wsConnected ? 'Live' : 'Disconnected'}
                    </Badge>
                  </div>
                </div>
                <div className="grid grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">24h High: </span>
                    <span className="font-mono">{formatPrice(selectedPair.high24h)}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">24h Low: </span>
                    <span className="font-mono">{formatPrice(selectedPair.low24h)}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">24h Volume: </span>
                    <span className="font-mono">{selectedPair.volume24h.toLocaleString()}</span>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="p-0">
                <TradingChart 
                  pair={selectedPair.symbol}
                  height={480}
                />
              </CardContent>
            </Card>
          </div>

          {/* Right Column - Order Form */}
          <div className="lg:col-span-1 space-y-4">
            {/* Trading Tabs */}
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="spot">Spot</TabsTrigger>
                <TabsTrigger value="margin">Margin</TabsTrigger>
                <TabsTrigger value="futures">Futures</TabsTrigger>
              </TabsList>
              
              <TabsContent value="spot" className="space-y-4">
                <OrderForm
                  pair={selectedPair}
                  orderType={orderType}
                  side={side}
                  onOrderTypeChange={setOrderType}
                  onSideChange={setSide}
                />
              </TabsContent>
              
              <TabsContent value="margin" className="space-y-4">
                <Card>
                  <CardContent className="p-4">
                    <div className="text-center text-muted-foreground">
                      <BarChart3 className="w-8 h-8 mx-auto mb-2" />
                      <p>Margin trading coming soon</p>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
              
              <TabsContent value="futures" className="space-y-4">
                <Card>
                  <CardContent className="p-4">
                    <div className="text-center text-muted-foreground">
                      <Zap className="w-8 h-8 mx-auto mb-2" />
                      <p>Futures trading coming soon</p>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>

            {/* Position Manager */}
            {isConnected && (
              <PositionManager userAddress={address} />
            )}
          </div>
        </div>

        {/* Bottom Section - Trade History and Open Orders */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-sm font-medium">Recent Trades</CardTitle>
            </CardHeader>
            <CardContent>
              <TradeHistory pair={selectedPair.symbol} />
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle className="text-sm font-medium">Open Orders</CardTitle>
            </CardHeader>
            <CardContent>
              {isConnected ? (
                <div className="text-center text-muted-foreground py-8">
                  <DollarSign className="w-8 h-8 mx-auto mb-2" />
                  <p>No open orders</p>
                </div>
              ) : (
                <div className="text-center text-muted-foreground py-8">
                  <p>Connect wallet to view orders</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
