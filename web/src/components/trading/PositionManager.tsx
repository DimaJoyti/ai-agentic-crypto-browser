'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ScrollArea } from '@/components/ui/scroll-area'
import { 
  TrendingUp, 
  TrendingDown, 
  X, 
  Target, 
  Shield,
  AlertTriangle,
  DollarSign,
  Percent,
  Clock
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface Position {
  id: string
  symbol: string
  side: 'long' | 'short'
  size: number
  entryPrice: number
  currentPrice: number
  pnl: number
  pnlPercent: number
  margin: number
  leverage: number
  liquidationPrice: number
  timestamp: number
  stopLoss?: number
  takeProfit?: number
}

interface OpenOrder {
  id: string
  symbol: string
  side: 'buy' | 'sell'
  type: 'limit' | 'market' | 'stop'
  amount: number
  price?: number
  stopPrice?: number
  filled: number
  status: 'pending' | 'partial' | 'filled' | 'cancelled'
  timestamp: number
}

interface PositionManagerProps {
  userAddress?: string
}

export function PositionManager({ userAddress }: PositionManagerProps) {
  const [positions, setPositions] = useState<Position[]>([])
  const [openOrders, setOpenOrders] = useState<OpenOrder[]>([])
  const [activeTab, setActiveTab] = useState('positions')
  const [isLoading, setIsLoading] = useState(false)

  // Generate mock data
  useEffect(() => {
    if (!userAddress) return

    const generateMockPositions = (): Position[] => {
      return [
        {
          id: 'pos_1',
          symbol: 'ETH/USDT',
          side: 'long',
          size: 2.5,
          entryPrice: 2450,
          currentPrice: 2520,
          pnl: 175,
          pnlPercent: 2.86,
          margin: 1225,
          leverage: 5,
          liquidationPrice: 2205,
          timestamp: Date.now() - 3600000,
          stopLoss: 2400,
          takeProfit: 2600
        },
        {
          id: 'pos_2',
          symbol: 'BTC/USDT',
          side: 'short',
          size: 0.1,
          entryPrice: 45200,
          currentPrice: 44800,
          pnl: 40,
          pnlPercent: 0.88,
          margin: 2260,
          leverage: 2,
          liquidationPrice: 47300,
          timestamp: Date.now() - 7200000
        }
      ]
    }

    const generateMockOrders = (): OpenOrder[] => {
      return [
        {
          id: 'order_1',
          symbol: 'ETH/USDT',
          side: 'buy',
          type: 'limit',
          amount: 1.0,
          price: 2480,
          filled: 0,
          status: 'pending',
          timestamp: Date.now() - 1800000
        },
        {
          id: 'order_2',
          symbol: 'BTC/USDT',
          side: 'sell',
          type: 'stop',
          amount: 0.05,
          stopPrice: 44000,
          filled: 0,
          status: 'pending',
          timestamp: Date.now() - 900000
        }
      ]
    }

    setIsLoading(true)
    setTimeout(() => {
      setPositions(generateMockPositions())
      setOpenOrders(generateMockOrders())
      setIsLoading(false)
    }, 500)
  }, [userAddress])

  const closePosition = (positionId: string) => {
    setPositions(prev => prev.filter(p => p.id !== positionId))
    // In real implementation, this would call the API to close the position
    console.log('Closing position:', positionId)
  }

  const cancelOrder = (orderId: string) => {
    setOpenOrders(prev => prev.filter(o => o.id !== orderId))
    // In real implementation, this would call the API to cancel the order
    console.log('Cancelling order:', orderId)
  }

  const formatPrice = (price: number) => {
    return price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 6 })
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getTotalPnL = () => {
    return positions.reduce((sum, pos) => sum + pos.pnl, 0)
  }

  const getTotalMargin = () => {
    return positions.reduce((sum, pos) => sum + pos.margin, 0)
  }

  if (!userAddress) {
    return (
      <Card>
        <CardContent className="p-6 text-center">
          <Shield className="w-8 h-8 mx-auto mb-2 text-muted-foreground" />
          <p className="text-muted-foreground">Connect wallet to view positions</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-medium">Portfolio</CardTitle>
          <div className="flex items-center gap-2 text-xs">
            <span className="text-muted-foreground">Total PnL:</span>
            <span className={cn(
              "font-medium",
              getTotalPnL() >= 0 ? "text-green-500" : "text-red-500"
            )}>
              {getTotalPnL() >= 0 ? '+' : ''}${getTotalPnL().toFixed(2)}
            </span>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="positions" className="text-xs">
              Positions ({positions.length})
            </TabsTrigger>
            <TabsTrigger value="orders" className="text-xs">
              Orders ({openOrders.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="positions" className="mt-4">
            {isLoading ? (
              <div className="space-y-2">
                {Array.from({ length: 2 }).map((_, i) => (
                  <div key={i} className="animate-pulse h-16 bg-muted rounded" />
                ))}
              </div>
            ) : positions.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <Target className="w-8 h-8 mx-auto mb-2 opacity-50" />
                <p>No open positions</p>
              </div>
            ) : (
              <ScrollArea className="h-64">
                <div className="space-y-2">
                  {positions.map((position) => (
                    <div
                      key={position.id}
                      className="p-3 border rounded-lg hover:bg-muted/50 transition-colors"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-sm">{position.symbol}</span>
                          <Badge
                            variant={position.side === 'long' ? 'default' : 'secondary'}
                            className={cn(
                              "text-xs",
                              position.side === 'long' 
                                ? "bg-green-100 text-green-700 border-green-200" 
                                : "bg-red-100 text-red-700 border-red-200"
                            )}
                          >
                            {position.side.toUpperCase()} {position.leverage}x
                          </Badge>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => closePosition(position.id)}
                          className="h-6 w-6 p-0 text-muted-foreground hover:text-red-500"
                        >
                          <X className="w-3 h-3" />
                        </Button>
                      </div>

                      <div className="grid grid-cols-2 gap-4 text-xs">
                        <div>
                          <div className="text-muted-foreground">Size</div>
                          <div className="font-mono">{position.size}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Entry</div>
                          <div className="font-mono">${formatPrice(position.entryPrice)}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Current</div>
                          <div className="font-mono">${formatPrice(position.currentPrice)}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">PnL</div>
                          <div className={cn(
                            "font-mono font-medium",
                            position.pnl >= 0 ? "text-green-500" : "text-red-500"
                          )}>
                            {position.pnl >= 0 ? '+' : ''}${position.pnl.toFixed(2)}
                            <span className="ml-1">
                              ({position.pnl >= 0 ? '+' : ''}{position.pnlPercent.toFixed(2)}%)
                            </span>
                          </div>
                        </div>
                      </div>

                      {(position.stopLoss || position.takeProfit) && (
                        <div className="mt-2 pt-2 border-t grid grid-cols-2 gap-4 text-xs">
                          {position.stopLoss && (
                            <div>
                              <div className="text-muted-foreground">Stop Loss</div>
                              <div className="font-mono text-red-500">${formatPrice(position.stopLoss)}</div>
                            </div>
                          )}
                          {position.takeProfit && (
                            <div>
                              <div className="text-muted-foreground">Take Profit</div>
                              <div className="font-mono text-green-500">${formatPrice(position.takeProfit)}</div>
                            </div>
                          )}
                        </div>
                      )}

                      <div className="mt-2 pt-2 border-t flex items-center justify-between text-xs text-muted-foreground">
                        <span>Margin: ${position.margin.toFixed(2)}</span>
                        <span>Liq: ${formatPrice(position.liquidationPrice)}</span>
                        <span>{formatTime(position.timestamp)}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            )}

            {positions.length > 0 && (
              <div className="mt-4 pt-3 border-t">
                <div className="grid grid-cols-2 gap-4 text-xs">
                  <div>
                    <span className="text-muted-foreground">Total Margin: </span>
                    <span className="font-medium">${getTotalMargin().toFixed(2)}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">Total PnL: </span>
                    <span className={cn(
                      "font-medium",
                      getTotalPnL() >= 0 ? "text-green-500" : "text-red-500"
                    )}>
                      {getTotalPnL() >= 0 ? '+' : ''}${getTotalPnL().toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>
            )}
          </TabsContent>

          <TabsContent value="orders" className="mt-4">
            {isLoading ? (
              <div className="space-y-2">
                {Array.from({ length: 2 }).map((_, i) => (
                  <div key={i} className="animate-pulse h-16 bg-muted rounded" />
                ))}
              </div>
            ) : openOrders.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <Clock className="w-8 h-8 mx-auto mb-2 opacity-50" />
                <p>No open orders</p>
              </div>
            ) : (
              <ScrollArea className="h-64">
                <div className="space-y-2">
                  {openOrders.map((order) => (
                    <div
                      key={order.id}
                      className="p-3 border rounded-lg hover:bg-muted/50 transition-colors"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-sm">{order.symbol}</span>
                          <Badge
                            variant={order.side === 'buy' ? 'default' : 'secondary'}
                            className={cn(
                              "text-xs",
                              order.side === 'buy' 
                                ? "bg-green-100 text-green-700 border-green-200" 
                                : "bg-red-100 text-red-700 border-red-200"
                            )}
                          >
                            {order.side.toUpperCase()}
                          </Badge>
                          <Badge variant="outline" className="text-xs">
                            {order.type.toUpperCase()}
                          </Badge>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => cancelOrder(order.id)}
                          className="h-6 w-6 p-0 text-muted-foreground hover:text-red-500"
                        >
                          <X className="w-3 h-3" />
                        </Button>
                      </div>

                      <div className="grid grid-cols-2 gap-4 text-xs">
                        <div>
                          <div className="text-muted-foreground">Amount</div>
                          <div className="font-mono">{order.amount}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">
                            {order.type === 'stop' ? 'Stop Price' : 'Price'}
                          </div>
                          <div className="font-mono">
                            ${formatPrice(order.stopPrice || order.price || 0)}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Filled</div>
                          <div className="font-mono">
                            {order.filled} / {order.amount}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Status</div>
                          <Badge
                            variant={order.status === 'pending' ? 'secondary' : 'default'}
                            className="text-xs"
                          >
                            {order.status}
                          </Badge>
                        </div>
                      </div>

                      <div className="mt-2 pt-2 border-t text-xs text-muted-foreground">
                        <span>{formatTime(order.timestamp)}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            )}
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}
