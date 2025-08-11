'use client'

import React, { useState, useEffect } from 'react'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { TrendingUp, TrendingDown, Activity, Star } from 'lucide-react'
import { cn } from '@/lib/utils'

interface TickerData {
  symbol: string
  price: number
  change24h: number
  changePercent24h: number
  volume24h: number
  high24h: number
  low24h: number
  isFavorite?: boolean
}

interface PriceTickerBarProps {
  maxItems?: number
  autoScroll?: boolean
  scrollSpeed?: number
}

export function PriceTickerBar({ 
  maxItems = 20, 
  autoScroll = true, 
  scrollSpeed = 50 
}: PriceTickerBarProps) {
  const [tickers, setTickers] = useState<TickerData[]>([])
  const [favorites, setFavorites] = useState<Set<string>>(new Set(['BTC/USDT', 'ETH/USDT']))
  const [isScrolling, setIsScrolling] = useState(autoScroll)

  // Generate mock ticker data
  useEffect(() => {
    const generateMockTickers = (): TickerData[] => {
      const symbols = [
        'BTC/USDT', 'ETH/USDT', 'BNB/USDT', 'ADA/USDT', 'SOL/USDT',
        'XRP/USDT', 'DOT/USDT', 'DOGE/USDT', 'AVAX/USDT', 'MATIC/USDT',
        'LINK/USDT', 'UNI/USDT', 'LTC/USDT', 'ATOM/USDT', 'FTM/USDT',
        'NEAR/USDT', 'ALGO/USDT', 'VET/USDT', 'ICP/USDT', 'THETA/USDT'
      ]

      return symbols.slice(0, maxItems).map(symbol => {
        const basePrice = symbol.includes('BTC') ? 45000 : 
                         symbol.includes('ETH') ? 2500 :
                         Math.random() * 100 + 1

        const changePercent = (Math.random() - 0.5) * 20 // -10% to +10%
        const change24h = basePrice * (changePercent / 100)
        
        return {
          symbol,
          price: basePrice + change24h,
          change24h,
          changePercent24h: changePercent,
          volume24h: Math.random() * 1000000000,
          high24h: basePrice * 1.1,
          low24h: basePrice * 0.9,
          isFavorite: favorites.has(symbol)
        }
      })
    }

    const updateTickers = () => {
      setTickers(generateMockTickers())
    }

    updateTickers()
    const interval = setInterval(updateTickers, 2000) // Update every 2 seconds

    return () => clearInterval(interval)
  }, [maxItems, favorites])

  const formatPrice = (price: number, symbol: string) => {
    if (symbol.includes('BTC')) {
      return `$${price.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}`
    }
    if (symbol.includes('ETH')) {
      return `$${price.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 2 })}`
    }
    return `$${price.toFixed(4)}`
  }

  const formatVolume = (volume: number) => {
    if (volume >= 1e9) {
      return `${(volume / 1e9).toFixed(1)}B`
    }
    if (volume >= 1e6) {
      return `${(volume / 1e6).toFixed(1)}M`
    }
    if (volume >= 1e3) {
      return `${(volume / 1e3).toFixed(1)}K`
    }
    return volume.toFixed(0)
  }

  const toggleFavorite = (symbol: string) => {
    setFavorites(prev => {
      const newFavorites = new Set(prev)
      if (newFavorites.has(symbol)) {
        newFavorites.delete(symbol)
      } else {
        newFavorites.add(symbol)
      }
      return newFavorites
    })
  }

  const sortedTickers = [...tickers].sort((a, b) => {
    // Favorites first, then by volume
    if (a.isFavorite && !b.isFavorite) return -1
    if (!a.isFavorite && b.isFavorite) return 1
    return b.volume24h - a.volume24h
  })

  return (
    <div className="bg-background border-b">
      <div className="relative overflow-hidden">
        <div 
          className={cn(
            "flex gap-4 py-3 px-4",
            isScrolling && "animate-scroll"
          )}
          style={{
            animationDuration: `${scrollSpeed}s`,
            animationIterationCount: 'infinite',
            animationTimingFunction: 'linear'
          }}
          onMouseEnter={() => setIsScrolling(false)}
          onMouseLeave={() => setIsScrolling(autoScroll)}
        >
          {sortedTickers.map((ticker, index) => {
            const isPositive = ticker.changePercent24h >= 0
            
            return (
              <div
                key={`${ticker.symbol}-${index}`}
                className="flex items-center gap-3 min-w-fit cursor-pointer hover:bg-muted/50 rounded-lg px-3 py-2 transition-colors"
                onClick={() => {
                  // Handle ticker click - could navigate to trading pair
                  console.log('Selected ticker:', ticker.symbol)
                }}
              >
                {/* Favorite Star */}
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    toggleFavorite(ticker.symbol)
                  }}
                  className="text-muted-foreground hover:text-yellow-500 transition-colors"
                >
                  <Star 
                    className={cn(
                      "w-4 h-4",
                      ticker.isFavorite && "fill-yellow-500 text-yellow-500"
                    )} 
                  />
                </button>

                {/* Symbol */}
                <div className="font-medium text-sm min-w-[80px]">
                  {ticker.symbol}
                </div>

                {/* Price */}
                <div className="font-mono text-sm min-w-[100px]">
                  {formatPrice(ticker.price, ticker.symbol)}
                </div>

                {/* Change */}
                <div className={cn(
                  "flex items-center gap-1 text-sm font-medium min-w-[80px]",
                  isPositive ? "text-green-500" : "text-red-500"
                )}>
                  {isPositive ? (
                    <TrendingUp className="w-3 h-3" />
                  ) : (
                    <TrendingDown className="w-3 h-3" />
                  )}
                  {isPositive ? '+' : ''}{ticker.changePercent24h.toFixed(2)}%
                </div>

                {/* Volume */}
                <div className="text-xs text-muted-foreground min-w-[60px]">
                  Vol: {formatVolume(ticker.volume24h)}
                </div>

                {/* Separator */}
                {index < sortedTickers.length - 1 && (
                  <div className="w-px h-6 bg-border" />
                )}
              </div>
            )
          })}
        </div>

        {/* Gradient overlays for scroll effect */}
        {isScrolling && (
          <>
            <div className="absolute left-0 top-0 bottom-0 w-8 bg-gradient-to-r from-background to-transparent pointer-events-none" />
            <div className="absolute right-0 top-0 bottom-0 w-8 bg-gradient-to-l from-background to-transparent pointer-events-none" />
          </>
        )}
      </div>

      {/* Market Status Indicator */}
      <div className="flex items-center justify-between px-4 py-2 bg-muted/30 text-xs">
        <div className="flex items-center gap-4">
          <Badge variant="default" className="text-xs">
            <Activity className="w-3 h-3 mr-1" />
            Market Open
          </Badge>
          <span className="text-muted-foreground">
            24h Volume: ${(tickers.reduce((sum, t) => sum + t.volume24h, 0) / 1e9).toFixed(2)}B
          </span>
        </div>
        <div className="flex items-center gap-4 text-muted-foreground">
          <span>Fear & Greed Index: 65 (Greed)</span>
          <span>Global Market Cap: $2.1T</span>
        </div>
      </div>
    </div>
  )
}

// Add CSS for scroll animation
const scrollAnimation = `
  @keyframes scroll {
    0% {
      transform: translateX(0);
    }
    100% {
      transform: translateX(-50%);
    }
  }
  
  .animate-scroll {
    animation: scroll linear infinite;
  }
`

// Inject CSS if not already present
if (typeof document !== 'undefined' && !document.getElementById('ticker-scroll-styles')) {
  const style = document.createElement('style')
  style.id = 'ticker-scroll-styles'
  style.textContent = scrollAnimation
  document.head.appendChild(style)
}
