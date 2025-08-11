'use client'

import React, { useState, useEffect } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { TrendingUp, TrendingDown, Clock, Filter, Download } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Trade {
  id: string
  timestamp: number
  pair: string
  side: 'buy' | 'sell'
  price: number
  amount: number
  total: number
  fee: number
  type: 'market' | 'limit'
}

interface TradeHistoryProps {
  pair?: string
  userAddress?: string
  maxItems?: number
  showFilters?: boolean
}

export function TradeHistory({ 
  pair, 
  userAddress, 
  maxItems = 50, 
  showFilters = true 
}: TradeHistoryProps) {
  const [trades, setTrades] = useState<Trade[]>([])
  const [filter, setFilter] = useState<'all' | 'buy' | 'sell'>('all')
  const [isLoading, setIsLoading] = useState(false)

  // Generate mock trade data
  useEffect(() => {
    const generateMockTrades = (): Trade[] => {
      const mockTrades: Trade[] = []
      const now = Date.now()
      
      for (let i = 0; i < maxItems; i++) {
        const timestamp = now - (i * 60000) // 1 minute intervals
        const side = Math.random() > 0.5 ? 'buy' : 'sell'
        const price = 2500 + (Math.random() - 0.5) * 100 // ETH price around 2500
        const amount = Math.random() * 5 + 0.1
        const total = price * amount
        const fee = total * 0.001 // 0.1% fee
        
        mockTrades.push({
          id: `trade_${i}`,
          timestamp,
          pair: pair || 'ETH/USDT',
          side,
          price,
          amount,
          total,
          fee,
          type: Math.random() > 0.7 ? 'market' : 'limit'
        })
      }
      
      return mockTrades
    }

    setIsLoading(true)
    // Simulate API call delay
    setTimeout(() => {
      setTrades(generateMockTrades())
      setIsLoading(false)
    }, 500)
  }, [pair, maxItems])

  const filteredTrades = trades.filter(trade => {
    if (filter === 'all') return true
    return trade.side === filter
  })

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  const formatPrice = (price: number) => {
    return price.toFixed(6)
  }

  const formatAmount = (amount: number) => {
    return amount.toFixed(6)
  }

  const exportTrades = () => {
    const csv = [
      'Time,Pair,Side,Price,Amount,Total,Fee,Type',
      ...filteredTrades.map(trade => 
        `${new Date(trade.timestamp).toISOString()},${trade.pair},${trade.side},${trade.price},${trade.amount},${trade.total},${trade.fee},${trade.type}`
      )
    ].join('\n')
    
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `trades_${pair || 'all'}_${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  if (isLoading) {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="animate-pulse">
            <div className="h-8 bg-muted rounded" />
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Filters and Controls */}
      {showFilters && (
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Button
              variant={filter === 'all' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setFilter('all')}
            >
              All
            </Button>
            <Button
              variant={filter === 'buy' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setFilter('buy')}
              className="text-green-600 border-green-600 hover:bg-green-50"
            >
              <TrendingUp className="w-3 h-3 mr-1" />
              Buy
            </Button>
            <Button
              variant={filter === 'sell' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setFilter('sell')}
              className="text-red-600 border-red-600 hover:bg-red-50"
            >
              <TrendingDown className="w-3 h-3 mr-1" />
              Sell
            </Button>
          </div>
          
          <Button variant="outline" size="sm" onClick={exportTrades}>
            <Download className="w-3 h-3 mr-1" />
            Export
          </Button>
        </div>
      )}

      {/* Table Header */}
      <div className="grid grid-cols-6 gap-2 text-xs font-medium text-muted-foreground pb-2 border-b">
        <div className="flex items-center gap-1">
          <Clock className="w-3 h-3" />
          Time
        </div>
        <div>Side</div>
        <div className="text-right">Price</div>
        <div className="text-right">Amount</div>
        <div className="text-right">Total</div>
        <div className="text-right">Fee</div>
      </div>

      {/* Trades List */}
      <ScrollArea className="h-64">
        <div className="space-y-1">
          {filteredTrades.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <Filter className="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p>No trades found</p>
              {userAddress && (
                <p className="text-xs mt-1">Connect your wallet to see your trade history</p>
              )}
            </div>
          ) : (
            filteredTrades.map((trade) => (
              <div
                key={trade.id}
                className="grid grid-cols-6 gap-2 text-xs py-2 hover:bg-muted/50 rounded transition-colors"
              >
                {/* Time */}
                <div className="text-muted-foreground font-mono">
                  {formatTime(trade.timestamp)}
                </div>

                {/* Side */}
                <div>
                  <Badge
                    variant="outline"
                    className={cn(
                      "text-xs",
                      trade.side === 'buy' 
                        ? "text-green-600 border-green-600" 
                        : "text-red-600 border-red-600"
                    )}
                  >
                    {trade.side.toUpperCase()}
                  </Badge>
                </div>

                {/* Price */}
                <div className="text-right font-mono">
                  {formatPrice(trade.price)}
                </div>

                {/* Amount */}
                <div className="text-right font-mono">
                  {formatAmount(trade.amount)}
                </div>

                {/* Total */}
                <div className="text-right font-mono font-medium">
                  {trade.total.toFixed(2)}
                </div>

                {/* Fee */}
                <div className="text-right font-mono text-muted-foreground">
                  {trade.fee.toFixed(4)}
                </div>
              </div>
            ))
          )}
        </div>
      </ScrollArea>

      {/* Summary */}
      {filteredTrades.length > 0 && (
        <div className="pt-2 border-t">
          <div className="grid grid-cols-3 gap-4 text-xs">
            <div>
              <span className="text-muted-foreground">Total Trades: </span>
              <span className="font-medium">{filteredTrades.length}</span>
            </div>
            <div>
              <span className="text-muted-foreground">Total Volume: </span>
              <span className="font-medium">
                {filteredTrades.reduce((sum, trade) => sum + trade.total, 0).toFixed(2)} USDT
              </span>
            </div>
            <div>
              <span className="text-muted-foreground">Total Fees: </span>
              <span className="font-medium">
                {filteredTrades.reduce((sum, trade) => sum + trade.fee, 0).toFixed(4)} USDT
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
