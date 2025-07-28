'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Input } from '@/components/ui/input'
import { 
  Image as ImageIcon, 
  TrendingUp, 
  DollarSign, 
  Crown,
  Shield,
  Star,
  Search,
  RefreshCw,
  ExternalLink,
  Eye,
  Heart,
  BarChart3,
  PieChart,
  Activity,
  Zap,
  Filter,
  Grid3X3,
  List
} from 'lucide-react'
import { useNFTCollection } from '@/hooks/useNFTCollection'
import { NFTRarity, type NFTCollection, type NFT } from '@/lib/nft-service'
import { type Address } from 'viem'

interface NFTCollectionDashboardProps {
  userAddress?: Address
  chainId?: number
}

export function NFTCollectionDashboard({ userAddress, chainId }: NFTCollectionDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')

  const {
    collections,
    userNFTs,
    userPortfolio,
    isLoading,
    loadData,
    searchCollections,
    getCollectionsByCategory,
    portfolioMetrics,
    collectionDistribution,
    rarityDistribution,
    topCollections,
    trendingCollections,
    featuredCollections,
    collectionCategories,
    formatCurrency,
    formatETH
  } = useNFTCollection({
    userAddress,
    chainId,
    autoRefresh: true,
    enableNotifications: true
  })

  const getRarityColor = (rarity: NFTRarity) => {
    switch (rarity) {
      case NFTRarity.COMMON:
        return 'bg-gray-100 text-gray-800'
      case NFTRarity.UNCOMMON:
        return 'bg-green-100 text-green-800'
      case NFTRarity.RARE:
        return 'bg-blue-100 text-blue-800'
      case NFTRarity.EPIC:
        return 'bg-purple-100 text-purple-800'
      case NFTRarity.LEGENDARY:
        return 'bg-yellow-100 text-yellow-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getRarityIcon = (rarity: NFTRarity) => {
    switch (rarity) {
      case NFTRarity.LEGENDARY:
        return <Crown className="w-3 h-3" />
      case NFTRarity.EPIC:
        return <Star className="w-3 h-3" />
      case NFTRarity.RARE:
        return <Zap className="w-3 h-3" />
      default:
        return null
    }
  }

  const filteredCollections = searchQuery 
    ? searchCollections(searchQuery)
    : selectedCategory 
      ? getCollectionsByCategory(selectedCategory)
      : collections

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <ImageIcon className="w-6 h-6" />
            NFT Collection Manager
          </h2>
          <p className="text-muted-foreground">
            Discover, track, and manage your NFT collections across all chains
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadData}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Portfolio Overview */}
      {userAddress && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                  <p className="text-2xl font-bold">{formatETH(portfolioMetrics.totalValue)}</p>
                  <p className="text-sm text-muted-foreground">
                    {portfolioMetrics.totalNFTs} NFTs
                  </p>
                </div>
                <DollarSign className="w-8 h-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Collections</p>
                  <p className="text-2xl font-bold">{portfolioMetrics.totalCollections}</p>
                  <p className="text-sm text-muted-foreground">
                    Avg: {formatETH(portfolioMetrics.averageValue)}
                  </p>
                </div>
                <ImageIcon className="w-8 h-8 text-purple-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">P&L</p>
                  <p className={`text-2xl font-bold ${portfolioMetrics.totalGainLoss >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {portfolioMetrics.totalGainLoss >= 0 ? '+' : ''}{formatETH(portfolioMetrics.totalGainLoss)}
                  </p>
                  <p className={`text-sm ${portfolioMetrics.totalGainLoss >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {portfolioMetrics.gainLossPercentage.toFixed(2)}%
                  </p>
                </div>
                <TrendingUp className="w-8 h-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Top NFT</p>
                  <p className="text-lg font-bold">
                    {portfolioMetrics.topNFT?.metadata.name || 'None'}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {portfolioMetrics.topNFT ? formatETH(portfolioMetrics.topNFT.currentPrice || '0') : 'N/A'}
                  </p>
                </div>
                <Crown className="w-8 h-8 text-yellow-500" />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="collections">Collections</TabsTrigger>
          <TabsTrigger value="portfolio">My NFTs</TabsTrigger>
          <TabsTrigger value="trending">Trending</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Featured Collections */}
          <Card>
            <CardHeader>
              <CardTitle>Featured Collections</CardTitle>
              <CardDescription>
                Curated and verified NFT collections
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {featuredCollections.slice(0, 3).map((collection, index) => (
                  <motion.div
                    key={collection.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors cursor-pointer"
                  >
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                        <ImageIcon className="w-6 h-6" />
                      </div>
                      <div>
                        <h4 className="font-medium flex items-center gap-2">
                          {collection.name}
                          {collection.verified && <Shield className="w-4 h-4 text-blue-500" />}
                        </h4>
                        <p className="text-sm text-muted-foreground">{collection.symbol}</p>
                      </div>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Floor Price</span>
                        <span className="font-medium">{formatETH(collection.floorPrice)}</span>
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">24h Volume</span>
                        <span className="font-medium">{formatETH(collection.volume24h)}</span>
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Owners</span>
                        <span className="font-medium">{collection.ownersCount.toLocaleString()}</span>
                      </div>
                    </div>

                    <Button className="w-full mt-3" size="sm" variant="outline">
                      <Eye className="w-3 h-3 mr-2" />
                      View Collection
                    </Button>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Portfolio Distribution */}
          {userAddress && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Collection Distribution</CardTitle>
                  <CardDescription>
                    Your NFTs by collection
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {collectionDistribution.slice(0, 5).map((item, index) => (
                      <div key={index} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium">{item.name}</span>
                          <span className="text-sm text-muted-foreground">
                            {item.count} ({item.percentage.toFixed(1)}%)
                          </span>
                        </div>
                        <Progress value={item.percentage} className="h-2" />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Rarity Distribution</CardTitle>
                  <CardDescription>
                    Your NFTs by rarity level
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {rarityDistribution.map((item, index) => (
                      <div key={index} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            {getRarityIcon(item.rarity)}
                            <span className="text-sm font-medium capitalize">{item.rarity}</span>
                          </div>
                          <span className="text-sm text-muted-foreground">
                            {item.count} ({item.percentage.toFixed(1)}%)
                          </span>
                        </div>
                        <Progress value={item.percentage} className="h-2" />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="collections" className="space-y-6">
          {/* Search and Filters */}
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center gap-4">
                <div className="flex-1 relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
                  <Input
                    placeholder="Search collections..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <select
                  value={selectedCategory}
                  onChange={(e) => setSelectedCategory(e.target.value)}
                  className="px-3 py-2 border rounded-md bg-background"
                >
                  <option value="">All Categories</option>
                  {collectionCategories.map(category => (
                    <option key={category} value={category}>{category}</option>
                  ))}
                </select>
                <div className="flex items-center gap-1 border rounded-md">
                  <Button
                    variant={viewMode === 'grid' ? 'default' : 'ghost'}
                    size="sm"
                    onClick={() => setViewMode('grid')}
                  >
                    <Grid3X3 className="w-4 h-4" />
                  </Button>
                  <Button
                    variant={viewMode === 'list' ? 'default' : 'ghost'}
                    size="sm"
                    onClick={() => setViewMode('list')}
                  >
                    <List className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Collections Grid/List */}
          <Card>
            <CardHeader>
              <CardTitle>All Collections</CardTitle>
              <CardDescription>
                Discover and explore NFT collections
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className={viewMode === 'grid' 
                ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
                : 'space-y-4'
              }>
                {filteredCollections.map((collection, index) => (
                  <motion.div
                    key={collection.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className={`border rounded-lg p-4 hover:bg-accent/50 transition-colors cursor-pointer ${
                      viewMode === 'list' ? 'flex items-center gap-4' : ''
                    }`}
                  >
                    <div className={`${viewMode === 'list' ? 'flex items-center gap-4 flex-1' : ''}`}>
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center flex-shrink-0">
                        <ImageIcon className="w-6 h-6" />
                      </div>
                      
                      <div className={`${viewMode === 'list' ? 'flex-1' : 'mt-3'}`}>
                        <div className="flex items-center gap-2 mb-1">
                          <h4 className="font-medium">{collection.name}</h4>
                          {collection.verified && <Shield className="w-4 h-4 text-blue-500" />}
                          {collection.featured && <Star className="w-4 h-4 text-yellow-500" />}
                        </div>
                        <p className="text-sm text-muted-foreground mb-2">{collection.symbol}</p>
                        
                        <div className={`${viewMode === 'list' ? 'flex items-center gap-6' : 'space-y-1'}`}>
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Floor</span>
                            <span className="font-medium">{formatETH(collection.floorPrice)}</span>
                          </div>
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Volume</span>
                            <span className="font-medium">{formatETH(collection.volume24h)}</span>
                          </div>
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Items</span>
                            <span className="font-medium">{collection.totalSupply.toLocaleString()}</span>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className={`${viewMode === 'list' ? 'flex items-center gap-2' : 'mt-3 flex gap-2'}`}>
                      <Button size="sm" variant="outline">
                        <Eye className="w-3 h-3 mr-2" />
                        View
                      </Button>
                      <Button size="sm" variant="outline">
                        <ExternalLink className="w-3 h-3" />
                      </Button>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="portfolio" className="space-y-6">
          {userNFTs.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle>My NFT Collection</CardTitle>
                <CardDescription>
                  Your owned NFTs across all collections
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                  {userNFTs.map((nft, index) => (
                    <motion.div
                      key={nft.id}
                      initial={{ opacity: 0, scale: 0.95 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: index * 0.05 }}
                      className="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
                    >
                      <div className="aspect-square bg-muted flex items-center justify-center">
                        <ImageIcon className="w-12 h-12 text-muted-foreground" />
                      </div>
                      
                      <div className="p-4">
                        <div className="flex items-center justify-between mb-2">
                          <h4 className="font-medium truncate">{nft.metadata.name}</h4>
                          {nft.rarity && (
                            <Badge className={getRarityColor(nft.rarity)}>
                              {getRarityIcon(nft.rarity)}
                              <span className="ml-1 capitalize">{nft.rarity}</span>
                            </Badge>
                          )}
                        </div>
                        
                        <div className="space-y-1 text-sm">
                          <div className="flex items-center justify-between">
                            <span className="text-muted-foreground">Value</span>
                            <span className="font-medium">{formatETH(nft.currentPrice || nft.floorPrice || '0')}</span>
                          </div>
                          {nft.rarityRank && (
                            <div className="flex items-center justify-between">
                              <span className="text-muted-foreground">Rank</span>
                              <span className="font-medium">#{nft.rarityRank}</span>
                            </div>
                          )}
                          {nft.isListed && (
                            <div className="flex items-center justify-between">
                              <span className="text-muted-foreground">Listed</span>
                              <span className="font-medium text-green-600">{formatETH(nft.listingPrice || '0')}</span>
                            </div>
                          )}
                        </div>

                        <div className="flex gap-2 mt-3">
                          <Button size="sm" variant="outline" className="flex-1">
                            <Eye className="w-3 h-3 mr-2" />
                            View
                          </Button>
                          <Button size="sm" variant="outline">
                            <Heart className="w-3 h-3" />
                          </Button>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="p-8 text-center">
                <ImageIcon className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">No NFTs Found</h3>
                <p className="text-muted-foreground">
                  {userAddress 
                    ? "You don't have any NFTs in your wallet yet"
                    : "Connect your wallet to view your NFT collection"
                  }
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="trending" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Trending Collections</CardTitle>
              <CardDescription>
                Collections with the highest 24h trading volume
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {trendingCollections.map((collection, index) => (
                  <motion.div
                    key={collection.id}
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center text-sm font-bold">
                        {index + 1}
                      </div>
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                        <ImageIcon className="w-6 h-6" />
                      </div>
                      <div>
                        <h4 className="font-medium flex items-center gap-2">
                          {collection.name}
                          {collection.verified && <Shield className="w-4 h-4 text-blue-500" />}
                        </h4>
                        <p className="text-sm text-muted-foreground">
                          Floor: {formatETH(collection.floorPrice)}
                        </p>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <p className="font-bold text-green-600">{formatETH(collection.volume24h)}</p>
                      <p className="text-sm text-muted-foreground">24h Volume</p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Market Overview
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Collections</span>
                    <span className="font-medium">{collections.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Verified Collections</span>
                    <span className="font-medium">{collections.filter(c => c.verified).length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Featured Collections</span>
                    <span className="font-medium">{featuredCollections.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Categories</span>
                    <span className="font-medium">{collectionCategories.length}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Portfolio Health
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Activity className="w-4 h-4 text-green-500" />
                      <span className="text-sm text-muted-foreground">Total NFTs</span>
                    </div>
                    <span className="font-medium">{portfolioMetrics.totalNFTs}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <TrendingUp className="w-4 h-4 text-blue-500" />
                      <span className="text-sm text-muted-foreground">Portfolio Value</span>
                    </div>
                    <span className="font-medium">{formatETH(portfolioMetrics.totalValue)}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <DollarSign className="w-4 h-4 text-purple-500" />
                      <span className="text-sm text-muted-foreground">Average Value</span>
                    </div>
                    <span className="font-medium">{formatETH(portfolioMetrics.averageValue)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
