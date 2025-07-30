'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { AlertTriangle } from 'lucide-react'
import { 
  Search,
  TrendingUp,
  Star,
  Eye,
  DollarSign,
  Users,
  Activity,
  BarChart3,
  RefreshCw,
  Settings,
  Filter,
  Grid,
  List,
  ExternalLink,
  Heart,
  Share,
  Zap,
  Crown,
  Flame
} from 'lucide-react'
import { useNFTDiscovery, useCollectionBrowser, useNFTSearch, useNFTCategories } from '@/hooks/useNFTDiscovery'
import { type NFTCollection, NFTCategory, type CollectionSearchFilters } from '@/lib/nft-discovery'
import { cn } from '@/lib/utils'

export function NFTDiscoveryDashboard() {
  const [activeTab, setActiveTab] = useState('trending')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<NFTCategory | ''>('')
  const [sortBy, setSortBy] = useState<'volume' | 'floor_price' | 'market_cap'>('volume')
  const [priceRange, setPriceRange] = useState<{ min: string; max: string }>({ min: '', max: '' })

  const {
    state,
    searchCollections,
    clearError
  } = useNFTDiscovery({
    autoLoad: true,
    enableNotifications: true
  })

  const { categories } = useNFTCategories()

  const handleSearch = async () => {
    if (!searchQuery.trim()) return

    const filters: CollectionSearchFilters = {
      category: selectedCategory || undefined,
      sortBy,
      sortOrder: 'desc',
      floorPriceMin: priceRange.min ? parseFloat(priceRange.min) : undefined,
      floorPriceMax: priceRange.max ? parseFloat(priceRange.max) : undefined,
      limit: 20
    }

    try {
      await searchCollections(searchQuery, filters)
    } catch (error) {
      console.error('Search failed:', error)
    }
  }

  const handleClearFilters = () => {
    setSearchQuery('')
    setSelectedCategory('')
    setPriceRange({ min: '', max: '' })
    setSortBy('volume')
  }

  const formatNumber = (value: number, decimals: number = 2) => {
    if (value >= 1000000) {
      return `${(value / 1000000).toFixed(decimals)}M`
    }
    if (value >= 1000) {
      return `${(value / 1000).toFixed(decimals)}K`
    }
    return value.toFixed(decimals)
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value)
  }

  const formatETH = (value: number) => {
    return `${value.toFixed(2)} ETH`
  }

  const getCategoryIcon = (category: NFTCategory) => {
    const icons: Record<NFTCategory, string> = {
      [NFTCategory.ART]: '🎨',
      [NFTCategory.COLLECTIBLES]: '🏆',
      [NFTCategory.GAMING]: '🎮',
      [NFTCategory.METAVERSE]: '🌐',
      [NFTCategory.MUSIC]: '🎵',
      [NFTCategory.PHOTOGRAPHY]: '📸',
      [NFTCategory.SPORTS]: '⚽',
      [NFTCategory.TRADING_CARDS]: '🃏',
      [NFTCategory.UTILITY]: '🔧',
      [NFTCategory.VIRTUAL_WORLDS]: '🏰',
      [NFTCategory.DOMAIN_NAMES]: '🌍',
      [NFTCategory.MEMES]: '😂'
    }
    return icons[category] || '🖼️'
  }

  const getChangeColor = (change: number) => {
    if (change > 0) return 'text-green-600'
    if (change < 0) return 'text-red-600'
    return 'text-gray-600'
  }

  const CollectionCard = ({ collection, index }: { collection: NFTCollection; index: number }) => (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.1 }}
      className="group"
    >
      <Card className="overflow-hidden hover:shadow-lg transition-all duration-300 cursor-pointer">
        <div className="relative">
          <div className="aspect-square bg-gradient-to-br from-blue-100 to-purple-100 dark:from-blue-900 dark:to-purple-900">
            {collection.imageUrl ? (
              <img
                src={collection.imageUrl}
                alt={collection.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-6xl">
                {getCategoryIcon(collection.category)}
              </div>
            )}
          </div>
          
          {/* Overlay with actions */}
          <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-2">
            <Button size="sm" variant="secondary">
              <Eye className="w-4 h-4 mr-2" />
              View
            </Button>
            <Button size="sm" variant="secondary">
              <Heart className="w-4 h-4" />
            </Button>
            <Button size="sm" variant="secondary">
              <Share className="w-4 h-4" />
            </Button>
          </div>

          {/* Badges */}
          <div className="absolute top-2 left-2 flex gap-1">
            {collection.verified && (
              <Badge variant="default" className="text-xs">
                <Crown className="w-3 h-3 mr-1" />
                Verified
              </Badge>
            )}
            {collection.trending && (
              <Badge variant="destructive" className="text-xs">
                <Flame className="w-3 h-3 mr-1" />
                Trending
              </Badge>
            )}
            {collection.featured && (
              <Badge variant="secondary" className="text-xs">
                <Star className="w-3 h-3 mr-1" />
                Featured
              </Badge>
            )}
          </div>

          {/* Rank */}
          <div className="absolute top-2 right-2">
            <Badge variant="outline" className="text-xs font-bold">
              #{index + 1}
            </Badge>
          </div>
        </div>

        <CardContent className="p-4">
          <div className="space-y-3">
            {/* Collection Info */}
            <div>
              <h3 className="font-semibold text-lg truncate">{collection.name}</h3>
              <p className="text-sm text-muted-foreground truncate">{collection.description}</p>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div>
                <p className="text-muted-foreground">Floor Price</p>
                <p className="font-medium">{formatETH(collection.floorPrice)}</p>
                <p className="text-xs text-muted-foreground">{formatCurrency(collection.floorPriceUSD)}</p>
              </div>
              <div>
                <p className="text-muted-foreground">24h Volume</p>
                <p className="font-medium">{formatETH(collection.volume24h)}</p>
                <p className={cn("text-xs", getChangeColor(collection.stats?.volumeChange24h || 0))}>
                  {collection.stats?.volumeChange24h > 0 ? '+' : ''}{collection.stats?.volumeChange24h?.toFixed(1)}%
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">Total Supply</p>
                <p className="font-medium">{formatNumber(collection.totalSupply, 0)}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Owners</p>
                <p className="font-medium">{formatNumber(collection.ownersCount, 0)}</p>
                <p className="text-xs text-muted-foreground">
                  {((collection.ownersCount / collection.totalSupply) * 100).toFixed(1)}% unique
                </p>
              </div>
            </div>

            {/* Category */}
            <div className="flex items-center justify-between">
              <Badge variant="outline" className="text-xs">
                {getCategoryIcon(collection.category)} {collection.category.replace(/_/g, ' ')}
              </Badge>
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                <Users className="w-3 h-3" />
                {formatNumber(collection.ownersCount, 0)}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">NFT Discovery</h2>
          <p className="text-muted-foreground">
            Discover, explore, and track NFT collections across the ecosystem
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}>
            {viewMode === 'grid' ? <List className="w-4 h-4" /> : <Grid className="w-4 h-4" />}
          </Button>
          <Button variant="outline" size="sm">
            <RefreshCw className="w-4 h-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Search & Filters</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Search Bar */}
          <div className="flex gap-2">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input
                placeholder="Search collections..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>
            <Button onClick={handleSearch} disabled={state.isSearching}>
              {state.isSearching ? (
                <RefreshCw className="w-4 h-4 animate-spin" />
              ) : (
                <Search className="w-4 h-4" />
              )}
            </Button>
          </div>

          {/* Filters */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {/* Category Filter */}
            <div className="space-y-2">
              <Label>Category</Label>
              <select
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value as NFTCategory | '')}
                className="w-full p-2 border rounded-md bg-background"
              >
                <option value="">All Categories</option>
                {categories.map((category) => (
                  <option key={category.value} value={category.value}>
                    {getCategoryIcon(category.value)} {category.label} ({category.count})
                  </option>
                ))}
              </select>
            </div>

            {/* Sort By */}
            <div className="space-y-2">
              <Label>Sort By</Label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as 'volume' | 'floor_price' | 'market_cap')}
                className="w-full p-2 border rounded-md bg-background"
              >
                <option value="volume">24h Volume</option>
                <option value="floor_price">Floor Price</option>
                <option value="market_cap">Market Cap</option>
              </select>
            </div>

            {/* Price Range */}
            <div className="space-y-2">
              <Label>Min Floor Price (ETH)</Label>
              <Input
                type="number"
                placeholder="0.0"
                value={priceRange.min}
                onChange={(e) => setPriceRange(prev => ({ ...prev, min: e.target.value }))}
              />
            </div>

            <div className="space-y-2">
              <Label>Max Floor Price (ETH)</Label>
              <Input
                type="number"
                placeholder="100.0"
                value={priceRange.max}
                onChange={(e) => setPriceRange(prev => ({ ...prev, max: e.target.value }))}
              />
            </div>
          </div>

          {/* Clear Filters */}
          <div className="flex justify-end">
            <Button variant="outline" size="sm" onClick={handleClearFilters}>
              <Filter className="w-4 h-4 mr-2" />
              Clear Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="trending">
            <Flame className="w-4 h-4 mr-2" />
            Trending ({state.trendingCollections.length})
          </TabsTrigger>
          <TabsTrigger value="featured">
            <Star className="w-4 h-4 mr-2" />
            Featured ({state.featuredCollections.length})
          </TabsTrigger>
          <TabsTrigger value="search">
            <Search className="w-4 h-4 mr-2" />
            Search Results ({state.collections.length})
          </TabsTrigger>
          <TabsTrigger value="categories">
            <Grid className="w-4 h-4 mr-2" />
            Categories ({categories.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="trending" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Trending Collections</h3>
            <Badge variant="destructive">
              <TrendingUp className="w-3 h-3 mr-1" />
              Hot
            </Badge>
          </div>

          {state.isLoading ? (
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {Array.from({ length: 8 }).map((_, i) => (
                <Card key={i} className="animate-pulse">
                  <div className="aspect-square bg-muted" />
                  <CardContent className="p-4">
                    <div className="space-y-2">
                      <div className="h-4 bg-muted rounded" />
                      <div className="h-3 bg-muted rounded w-2/3" />
                      <div className="grid grid-cols-2 gap-2">
                        <div className="h-3 bg-muted rounded" />
                        <div className="h-3 bg-muted rounded" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : state.trendingCollections.length > 0 ? (
            <div className={cn(
              viewMode === 'grid' 
                ? "grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
                : "space-y-4"
            )}>
              {state.trendingCollections.map((collection, index) => (
                <CollectionCard key={collection.id} collection={collection} index={index} />
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Flame className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Trending Collections</h3>
                <p className="text-muted-foreground">
                  Trending collections will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="featured" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Featured Collections</h3>
            <Badge variant="secondary">
              <Star className="w-3 h-3 mr-1" />
              Curated
            </Badge>
          </div>

          {state.featuredCollections.length > 0 ? (
            <div className={cn(
              viewMode === 'grid' 
                ? "grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
                : "space-y-4"
            )}>
              {state.featuredCollections.map((collection, index) => (
                <CollectionCard key={collection.id} collection={collection} index={index} />
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Star className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Featured Collections</h3>
                <p className="text-muted-foreground">
                  Featured collections will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="search" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">
              Search Results
              {state.searchQuery && (
                <span className="text-muted-foreground ml-2">for "{state.searchQuery}"</span>
              )}
            </h3>
            {state.total > 0 && (
              <Badge variant="outline">
                {state.total} results
              </Badge>
            )}
          </div>

          {state.isSearching ? (
            <div className="text-center py-8">
              <RefreshCw className="w-8 h-8 mx-auto animate-spin text-muted-foreground mb-4" />
              <p className="text-muted-foreground">Searching collections...</p>
            </div>
          ) : state.collections.length > 0 ? (
            <div className={cn(
              viewMode === 'grid' 
                ? "grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
                : "space-y-4"
            )}>
              {state.collections.map((collection, index) => (
                <CollectionCard key={collection.id} collection={collection} index={index} />
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Search className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Search Results</h3>
                <p className="text-muted-foreground">
                  {state.searchQuery 
                    ? `No collections found for "${state.searchQuery}"`
                    : "Enter a search term to find collections"
                  }
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="categories" className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Browse by Category</h3>
            <Badge variant="outline">
              {categories.length} categories
            </Badge>
          </div>

          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {categories.map((category, index) => (
              <motion.div
                key={category.value}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="hover:shadow-lg transition-all duration-300 cursor-pointer">
                  <CardContent className="p-6 text-center">
                    <div className="text-4xl mb-3">{getCategoryIcon(category.value)}</div>
                    <h4 className="font-semibold mb-2">{category.label}</h4>
                    <p className="text-sm text-muted-foreground">
                      {category.count} collections
                    </p>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
