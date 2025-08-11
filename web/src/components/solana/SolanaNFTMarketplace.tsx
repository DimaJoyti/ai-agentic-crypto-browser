'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  Image as ImageIcon,
  TrendingUp,
  TrendingDown,
  Search,
  Grid3X3,
  List,
  ExternalLink,
  RefreshCw,
  Eye,
  ShoppingCart,
  BarChart3,
  Loader2
} from 'lucide-react'
import { useSolanaNFT } from '@/hooks/useSolanaNFT'
import { formatNumber, formatPercentage, cn } from '@/lib/utils'

interface SolanaNFTMarketplaceProps {
  className?: string
}

export function SolanaNFTMarketplace({ className }: SolanaNFTMarketplaceProps) {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [searchQuery, setSearchQuery] = useState('')
  const [sortBy, setSortBy] = useState<'volume' | 'floor_price' | 'market_cap' | 'sales'>('volume')
  const [selectedMarketplace, setSelectedMarketplace] = useState<'all' | 'magic-eden' | 'tensor'>('all')
  const [priceRange] = useState({ min: '', max: '' })

  const {
    collections,
    marketStats,
    isLoading,
    error,
    refresh
  } = useSolanaNFT({
    autoRefresh: true,
    refreshInterval: 120000,
    sortBy,
    limit: 50
  })

  const filteredCollections = collections.filter(collection => {
    const matchesSearch = searchQuery === '' || 
      collection.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      collection.symbol.toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesMarketplace = selectedMarketplace === 'all' || 
      collection.marketplace === selectedMarketplace
    
    const matchesPriceRange = 
      (priceRange.min === '' || collection.floorPrice >= parseFloat(priceRange.min)) &&
      (priceRange.max === '' || collection.floorPrice <= parseFloat(priceRange.max))
    
    return matchesSearch && matchesMarketplace && matchesPriceRange
  })

  const handleRefresh = () => {
    refresh()
  }

  const getChangeColor = (change: number) => {
    if (change > 0) return 'text-green-600'
    if (change < 0) return 'text-red-600'
    return 'text-gray-600'
  }

  const getChangeIcon = (change: number) => {
    if (change > 0) return <TrendingUp className="h-3 w-3" />
    if (change < 0) return <TrendingDown className="h-3 w-3" />
    return null
  }

  const CollectionCard = ({ collection, index }: { collection: any; index: number }) => (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow">
      <div className="aspect-square relative bg-muted">
        <img
          src={collection.image || '/placeholder-nft.png'}
          alt={collection.name}
          className="w-full h-full object-cover"
          onError={(e) => {
            (e.target as HTMLImageElement).src = '/placeholder-nft.png'
          }}
        />
        <div className="absolute top-2 left-2">
          <Badge variant="secondary">#{index + 1}</Badge>
        </div>
        <div className="absolute top-2 right-2">
          <Badge variant="outline" className="bg-background/80">
            {collection.marketplace}
          </Badge>
        </div>
      </div>
      <CardContent className="p-4 space-y-3">
        <div>
          <h3 className="font-semibold truncate">{collection.name}</h3>
          <p className="text-sm text-muted-foreground">{collection.symbol}</p>
        </div>
        
        <div className="grid grid-cols-2 gap-2 text-sm">
          <div>
            <p className="text-muted-foreground">Floor Price</p>
            <div className="flex items-center space-x-1">
              <span className="font-semibold">
                {formatNumber(collection.floorPrice)} SOL
              </span>
              {collection.floorPriceChange24h !== 0 && (
                <div className={cn("flex items-center", getChangeColor(collection.floorPriceChange24h))}>
                  {getChangeIcon(collection.floorPriceChange24h)}
                  <span className="text-xs">
                    {formatPercentage(Math.abs(collection.floorPriceChange24h))}
                  </span>
                </div>
              )}
            </div>
          </div>
          
          <div>
            <p className="text-muted-foreground">24h Volume</p>
            <p className="font-semibold">
              {formatNumber(collection.volume24h, 0)} SOL
            </p>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground">
          <div>
            <p>Supply</p>
            <p className="font-medium">{formatNumber(collection.supply, 0)}</p>
          </div>
          <div>
            <p>Owners</p>
            <p className="font-medium">{formatNumber(collection.owners, 0)}</p>
          </div>
          <div>
            <p>Listed</p>
            <p className="font-medium">{formatPercentage(collection.listedPercentage)}</p>
          </div>
        </div>

        <div className="flex space-x-2">
          <Button variant="outline" size="sm" className="flex-1">
            <Eye className="h-3 w-3 mr-1" />
            View
          </Button>
          <Button variant="outline" size="sm" className="flex-1">
            <ExternalLink className="h-3 w-3 mr-1" />
            Trade
          </Button>
        </div>
      </CardContent>
    </Card>
  )

  const CollectionListItem = ({ collection, index }: { collection: any; index: number }) => (
    <Card className="p-4">
      <div className="flex items-center space-x-4">
        <div className="flex items-center space-x-3">
          <Badge variant="secondary">#{index + 1}</Badge>
          <img
            src={collection.image || '/placeholder-nft.png'}
            alt={collection.name}
            className="w-12 h-12 rounded-md object-cover"
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder-nft.png'
            }}
          />
          <div>
            <h3 className="font-semibold">{collection.name}</h3>
            <p className="text-sm text-muted-foreground">{collection.symbol}</p>
          </div>
        </div>

        <div className="flex-1 grid grid-cols-5 gap-4 text-sm">
          <div>
            <p className="text-muted-foreground">Floor Price</p>
            <div className="flex items-center space-x-1">
              <span className="font-semibold">
                {formatNumber(collection.floorPrice)} SOL
              </span>
              {collection.floorPriceChange24h !== 0 && (
                <div className={cn("flex items-center", getChangeColor(collection.floorPriceChange24h))}>
                  {getChangeIcon(collection.floorPriceChange24h)}
                  <span className="text-xs">
                    {formatPercentage(Math.abs(collection.floorPriceChange24h))}
                  </span>
                </div>
              )}
            </div>
          </div>

          <div>
            <p className="text-muted-foreground">24h Volume</p>
            <p className="font-semibold">{formatNumber(collection.volume24h, 0)} SOL</p>
          </div>

          <div>
            <p className="text-muted-foreground">Market Cap</p>
            <p className="font-semibold">{formatNumber(collection.marketCap, 0)} SOL</p>
          </div>

          <div>
            <p className="text-muted-foreground">Supply</p>
            <p className="font-semibold">{formatNumber(collection.supply, 0)}</p>
          </div>

          <div>
            <p className="text-muted-foreground">Owners</p>
            <p className="font-semibold">{formatNumber(collection.owners, 0)}</p>
          </div>
        </div>

        <div className="flex space-x-2">
          <Button variant="outline" size="sm">
            <Eye className="h-3 w-3 mr-1" />
            View
          </Button>
          <Button variant="outline" size="sm">
            <ExternalLink className="h-3 w-3 mr-1" />
            Trade
          </Button>
        </div>
      </div>
    </Card>
  )

  return (
    <div className={cn('space-y-6', className)}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">NFT Marketplace</h1>
          <p className="text-muted-foreground">
            Discover and trade NFT collections on Solana
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isLoading}
          >
            <RefreshCw className={cn("h-4 w-4 mr-2", isLoading && "animate-spin")} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Market Stats */}
      {marketStats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">24h Volume</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatNumber(marketStats.totalVolume24h, 0)} SOL
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">24h Sales</CardTitle>
              <ShoppingCart className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatNumber(marketStats.totalSales24h, 0)}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Collections</CardTitle>
              <ImageIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatNumber(marketStats.totalCollections, 0)}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Avg Floor</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatNumber(marketStats.averageFloorPrice)} SOL
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap gap-4">
            <div className="flex-1 min-w-64">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search collections..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            <Select value={sortBy} onValueChange={(value: any) => setSortBy(value)}>
              <SelectTrigger className="w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="volume">Volume</SelectItem>
                <SelectItem value="floor_price">Floor Price</SelectItem>
                <SelectItem value="market_cap">Market Cap</SelectItem>
                <SelectItem value="sales">Sales</SelectItem>
              </SelectContent>
            </Select>

            <Select value={selectedMarketplace} onValueChange={(value: any) => setSelectedMarketplace(value)}>
              <SelectTrigger className="w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Markets</SelectItem>
                <SelectItem value="magic-eden">Magic Eden</SelectItem>
                <SelectItem value="tensor">Tensor</SelectItem>
              </SelectContent>
            </Select>

            <div className="flex items-center space-x-2">
              <Button
                variant={viewMode === 'grid' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setViewMode('grid')}
              >
                <Grid3X3 className="h-4 w-4" />
              </Button>
              <Button
                variant={viewMode === 'list' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setViewMode('list')}
              >
                <List className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Error Display */}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Collections Grid/List */}
      <div className="space-y-4">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin" />
            <span className="ml-2">Loading collections...</span>
          </div>
        ) : (
          <>
            {viewMode === 'grid' ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {filteredCollections.map((collection, index) => (
                  <CollectionCard key={collection.id} collection={collection} index={index} />
                ))}
              </div>
            ) : (
              <div className="space-y-4">
                {filteredCollections.map((collection, index) => (
                  <CollectionListItem key={collection.id} collection={collection} index={index} />
                ))}
              </div>
            )}

            {filteredCollections.length === 0 && !isLoading && (
              <div className="text-center py-12">
                <ImageIcon className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No collections found</h3>
                <p className="text-muted-foreground">
                  Try adjusting your search or filter criteria
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
