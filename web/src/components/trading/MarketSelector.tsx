'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Search, 
  Star, 
  TrendingUp, 
  TrendingDown, 
  Volume2,
  Filter,
  ArrowUpDown
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface TradingPair {
  symbol: string
  baseAsset: string
  quoteAsset: string
  price: number
  change24h: number
  volume24h: number
  high24h: number
  low24h: number
  isFavorite?: boolean
  category?: string
}

interface MarketSelectorProps {
  selectedPair: TradingPair
  onPairChange: (pair: TradingPair) => void
  isOpen?: boolean
  onClose?: () => void
}

type SortField = 'symbol' | 'price' | 'change24h' | 'volume24h'
type SortDirection = 'asc' | 'desc'

export function MarketSelector({ 
  selectedPair, 
  onPairChange, 
  isOpen = false, 
  onClose 
}: MarketSelectorProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [activeTab, setActiveTab] = useState('favorites')
  const [favorites, setFavorites] = useState<Set<string>>(new Set(['BTC/USDT', 'ETH/USDT', 'BNB/USDT']))
  const [sortField, setSortField] = useState<SortField>('volume24h')
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc')
  const [allPairs, setAllPairs] = useState<TradingPair[]>([])

  // Generate mock trading pairs
  useEffect(() => {
    const generateMockPairs = (): TradingPair[] => {
      const baseAssets = [
        'BTC', 'ETH', 'BNB', 'ADA', 'SOL', 'XRP', 'DOT', 'DOGE', 'AVAX', 'MATIC',
        'LINK', 'UNI', 'LTC', 'ATOM', 'FTM', 'NEAR', 'ALGO', 'VET', 'ICP', 'THETA',
        'AAVE', 'COMP', 'MKR', 'SNX', 'CRV', 'YFI', 'SUSHI', '1INCH', 'BAL', 'REN'
      ]
      
      const quoteAssets = ['USDT', 'USDC', 'BTC', 'ETH', 'BNB']
      const categories = ['DeFi', 'Layer 1', 'Layer 2', 'Meme', 'Gaming', 'NFT', 'AI']

      const pairs: TradingPair[] = []

      baseAssets.forEach(base => {
        quoteAssets.forEach(quote => {
          if (base !== quote) {
            const symbol = `${base}/${quote}`
            const basePrice = base === 'BTC' ? 45000 : 
                            base === 'ETH' ? 2500 :
                            Math.random() * 100 + 1

            const changePercent = (Math.random() - 0.5) * 20
            const change24h = basePrice * (changePercent / 100)

            pairs.push({
              symbol,
              baseAsset: base,
              quoteAsset: quote,
              price: basePrice + change24h,
              change24h: changePercent,
              volume24h: Math.random() * 1000000000,
              high24h: basePrice * 1.1,
              low24h: basePrice * 0.9,
              isFavorite: favorites.has(symbol),
              category: categories[Math.floor(Math.random() * categories.length)]
            })
          }
        })
      })

      return pairs.slice(0, 100) // Limit to 100 pairs for demo
    }

    setAllPairs(generateMockPairs())
  }, [favorites])

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDirection('desc')
    }
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

  const filteredPairs = allPairs
    .filter(pair => {
      const matchesSearch = pair.symbol.toLowerCase().includes(searchTerm.toLowerCase()) ||
                           pair.baseAsset.toLowerCase().includes(searchTerm.toLowerCase())
      
      switch (activeTab) {
        case 'favorites':
          return matchesSearch && favorites.has(pair.symbol)
        case 'usdt':
          return matchesSearch && pair.quoteAsset === 'USDT'
        case 'btc':
          return matchesSearch && pair.quoteAsset === 'BTC'
        case 'eth':
          return matchesSearch && pair.quoteAsset === 'ETH'
        case 'defi':
          return matchesSearch && pair.category === 'DeFi'
        default:
          return matchesSearch
      }
    })
    .sort((a, b) => {
      let aValue: number | string
      let bValue: number | string

      switch (sortField) {
        case 'symbol':
          aValue = a.symbol
          bValue = b.symbol
          break
        case 'price':
          aValue = a.price
          bValue = b.price
          break
        case 'change24h':
          aValue = a.change24h
          bValue = b.change24h
          break
        case 'volume24h':
          aValue = a.volume24h
          bValue = b.volume24h
          break
        default:
          return 0
      }

      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortDirection === 'asc' 
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue)
      }

      return sortDirection === 'asc' 
        ? (aValue as number) - (bValue as number)
        : (bValue as number) - (aValue as number)
    })

  const formatPrice = (price: number) => {
    if (price >= 1000) {
      return price.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 2 })
    }
    return price.toFixed(6)
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

  if (!isOpen) {
    return (
      <Card className="cursor-pointer hover:shadow-md transition-shadow" onClick={() => {}}>
        <CardContent className="p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <h3 className="text-lg font-bold">{selectedPair.symbol}</h3>
              <Badge variant="outline">{selectedPair.category || 'Spot'}</Badge>
            </div>
            <div className="text-right">
              <div className="text-lg font-mono">${formatPrice(selectedPair.price)}</div>
              <div className={cn(
                "text-sm flex items-center gap-1",
                selectedPair.change24h >= 0 ? "text-green-500" : "text-red-500"
              )}>
                {selectedPair.change24h >= 0 ? (
                  <TrendingUp className="w-3 h-3" />
                ) : (
                  <TrendingDown className="w-3 h-3" />
                )}
                {selectedPair.change24h >= 0 ? '+' : ''}{selectedPair.change24h.toFixed(2)}%
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full max-w-4xl">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Select Trading Pair</CardTitle>
          <Button variant="ghost" size="sm" onClick={onClose}>
            ×
          </Button>
        </div>
        
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input
            placeholder="Search pairs..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
      </CardHeader>

      <CardContent>
        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab} className="mb-4">
          <TabsList className="grid w-full grid-cols-6">
            <TabsTrigger value="favorites" className="flex items-center gap-1">
              <Star className="w-3 h-3" />
              Favorites
            </TabsTrigger>
            <TabsTrigger value="usdt">USDT</TabsTrigger>
            <TabsTrigger value="btc">BTC</TabsTrigger>
            <TabsTrigger value="eth">ETH</TabsTrigger>
            <TabsTrigger value="defi">DeFi</TabsTrigger>
            <TabsTrigger value="all">All</TabsTrigger>
          </TabsList>
        </Tabs>

        {/* Table Header */}
        <div className="grid grid-cols-5 gap-4 p-2 text-xs font-medium text-muted-foreground border-b">
          <button 
            className="flex items-center gap-1 hover:text-foreground transition-colors"
            onClick={() => handleSort('symbol')}
          >
            Pair
            <ArrowUpDown className="w-3 h-3" />
          </button>
          <button 
            className="flex items-center gap-1 hover:text-foreground transition-colors text-right"
            onClick={() => handleSort('price')}
          >
            Price
            <ArrowUpDown className="w-3 h-3" />
          </button>
          <button 
            className="flex items-center gap-1 hover:text-foreground transition-colors text-right"
            onClick={() => handleSort('change24h')}
          >
            24h Change
            <ArrowUpDown className="w-3 h-3" />
          </button>
          <button 
            className="flex items-center gap-1 hover:text-foreground transition-colors text-right"
            onClick={() => handleSort('volume24h')}
          >
            24h Volume
            <ArrowUpDown className="w-3 h-3" />
          </button>
          <div className="text-center">Action</div>
        </div>

        {/* Pairs List */}
        <div className="max-h-96 overflow-y-auto">
          {filteredPairs.map((pair) => (
            <div
              key={pair.symbol}
              className={cn(
                "grid grid-cols-5 gap-4 p-3 hover:bg-muted/50 cursor-pointer transition-colors border-b border-border/50",
                selectedPair.symbol === pair.symbol && "bg-muted"
              )}
              onClick={() => {
                onPairChange(pair)
                onClose?.()
              }}
            >
              {/* Pair */}
              <div className="flex items-center gap-2">
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    toggleFavorite(pair.symbol)
                  }}
                  className="text-muted-foreground hover:text-yellow-500 transition-colors"
                >
                  <Star 
                    className={cn(
                      "w-4 h-4",
                      favorites.has(pair.symbol) && "fill-yellow-500 text-yellow-500"
                    )} 
                  />
                </button>
                <div>
                  <div className="font-medium">{pair.symbol}</div>
                  <div className="text-xs text-muted-foreground">{pair.category}</div>
                </div>
              </div>

              {/* Price */}
              <div className="text-right font-mono">
                ${formatPrice(pair.price)}
              </div>

              {/* Change */}
              <div className={cn(
                "text-right flex items-center justify-end gap-1",
                pair.change24h >= 0 ? "text-green-500" : "text-red-500"
              )}>
                {pair.change24h >= 0 ? (
                  <TrendingUp className="w-3 h-3" />
                ) : (
                  <TrendingDown className="w-3 h-3" />
                )}
                {pair.change24h >= 0 ? '+' : ''}{pair.change24h.toFixed(2)}%
              </div>

              {/* Volume */}
              <div className="text-right text-muted-foreground">
                {formatVolume(pair.volume24h)}
              </div>

              {/* Action */}
              <div className="text-center">
                <Button size="sm" variant="outline">
                  Trade
                </Button>
              </div>
            </div>
          ))}
        </div>

        {filteredPairs.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            No trading pairs found
          </div>
        )}
      </CardContent>
    </Card>
  )
}
