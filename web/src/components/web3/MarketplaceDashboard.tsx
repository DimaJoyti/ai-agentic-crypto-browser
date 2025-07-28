'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Store, 
  TrendingUp, 
  DollarSign, 
  Activity,
  RefreshCw,
  ExternalLink,
  Plus,
  Minus,
  Clock,
  Users,
  BarChart3,
  PieChart,
  Zap,
  Eye,
  ShoppingCart,
  Tag,
  Gavel,
  Timer
} from 'lucide-react'
import { useMarketplace } from '@/hooks/useMarketplace'
import { MarketplaceName, OrderType, OrderStatus } from '@/lib/marketplace-service'
import { type Address } from 'viem'

interface MarketplaceDashboardProps {
  userAddress?: Address
}

export function MarketplaceDashboard({ userAddress }: MarketplaceDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const {
    marketplaces,
    listings,
    offers,
    recentActivity,
    marketplaceStats,
    isLoading,
    loadData,
    createListing,
    createOffer,
    cancelListing,
    cancelOffer,
    marketplaceAnalytics,
    marketplaceDistribution,
    userListings,
    userOffers,
    formatETH,
    getMarketplaceIcon
  } = useMarketplace({
    userAddress,
    autoRefresh: true,
    enableNotifications: true
  })

  const getOrderTypeIcon = (orderType: OrderType) => {
    switch (orderType) {
      case OrderType.LISTING:
        return <Tag className="w-4 h-4" />
      case OrderType.OFFER:
        return <DollarSign className="w-4 h-4" />
      case OrderType.AUCTION:
        return <Gavel className="w-4 h-4" />
      case OrderType.DUTCH_AUCTION:
        return <Timer className="w-4 h-4" />
      default:
        return <Store className="w-4 h-4" />
    }
  }

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.ACTIVE:
        return 'bg-green-100 text-green-800'
      case OrderStatus.FILLED:
        return 'bg-blue-100 text-blue-800'
      case OrderStatus.CANCELLED:
        return 'bg-gray-100 text-gray-800'
      case OrderStatus.EXPIRED:
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatTimeRemaining = (endTime: number) => {
    const remaining = endTime - Date.now()
    if (remaining <= 0) return 'Expired'
    
    const days = Math.floor(remaining / (1000 * 60 * 60 * 24))
    const hours = Math.floor((remaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    
    if (days > 0) return `${days}d ${hours}h`
    if (hours > 0) return `${hours}h`
    return '< 1h'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Store className="w-6 h-6" />
            NFT Marketplace Hub
          </h2>
          <p className="text-muted-foreground">
            Trade NFTs across multiple marketplaces with real-time data and analytics
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadData}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Marketplace Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Volume 24h</p>
                <p className="text-2xl font-bold">{formatETH(marketplaceAnalytics.totalVolume24h)}</p>
                <p className="text-sm text-muted-foreground">
                  {marketplaceAnalytics.totalSales24h} sales
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
                <p className="text-sm font-medium text-muted-foreground">Active Listings</p>
                <p className="text-2xl font-bold">{marketplaceAnalytics.totalListings}</p>
                <p className="text-sm text-muted-foreground">
                  {marketplaceAnalytics.totalOffers} offers
                </p>
              </div>
              <ShoppingCart className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Average Price</p>
                <p className="text-2xl font-bold">{formatETH(marketplaceAnalytics.averagePrice)}</p>
                <p className="text-sm text-muted-foreground">
                  24h average
                </p>
              </div>
              <BarChart3 className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Top Marketplace</p>
                <p className="text-2xl font-bold flex items-center gap-2">
                  {getMarketplaceIcon(marketplaceAnalytics.topMarketplace)}
                  {marketplaceAnalytics.topMarketplace}
                </p>
                <p className="text-sm text-muted-foreground">
                  By volume
                </p>
              </div>
              <Store className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="marketplaces">Marketplaces</TabsTrigger>
          <TabsTrigger value="listings">Listings</TabsTrigger>
          <TabsTrigger value="offers">Offers</TabsTrigger>
          <TabsTrigger value="activity">Activity</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Marketplace Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Marketplace Distribution</CardTitle>
                <CardDescription>
                  Active listings by marketplace
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {marketplaceDistribution.map((item, index) => (
                    <div key={index} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-lg">{getMarketplaceIcon(item.marketplace)}</span>
                          <span className="text-sm font-medium">{item.marketplace}</span>
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

            <Card>
              <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
                <CardDescription>
                  Latest marketplace transactions
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {recentActivity.slice(0, 5).map((activity, index) => (
                    <motion.div
                      key={activity.id}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center">
                          <Activity className="w-4 h-4" />
                        </div>
                        <div>
                          <p className="text-sm font-medium capitalize">{activity.type}</p>
                          <p className="text-xs text-muted-foreground">
                            {getMarketplaceIcon(activity.marketplace)} {activity.marketplace} • #{activity.tokenId}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        {activity.price && (
                          <p className="text-sm font-medium">{formatETH(activity.price)}</p>
                        )}
                        <p className="text-xs text-muted-foreground">
                          {new Date(activity.timestamp).toLocaleTimeString()}
                        </p>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* User's Active Orders */}
          {userAddress && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>My Listings</CardTitle>
                  <CardDescription>
                    Your active NFT listings
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {userListings.length > 0 ? (
                    <div className="space-y-3">
                      {userListings.slice(0, 3).map((listing, index) => (
                        <motion.div
                          key={listing.id}
                          initial={{ opacity: 0, x: 20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                          className="flex items-center justify-between p-3 border rounded-lg"
                        >
                          <div>
                            <p className="text-sm font-medium">#{listing.tokenId}</p>
                            <p className="text-xs text-muted-foreground">
                              {getMarketplaceIcon(listing.marketplace)} {listing.marketplace}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm font-medium">{formatETH(listing.price)}</p>
                            <p className="text-xs text-muted-foreground">
                              {listing.endTime ? formatTimeRemaining(listing.endTime) : 'No expiry'}
                            </p>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <Tag className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                      <p className="text-sm text-muted-foreground">No active listings</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>My Offers</CardTitle>
                  <CardDescription>
                    Your active NFT offers
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {userOffers.length > 0 ? (
                    <div className="space-y-3">
                      {userOffers.slice(0, 3).map((offer, index) => (
                        <motion.div
                          key={offer.id}
                          initial={{ opacity: 0, x: 20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.1 }}
                          className="flex items-center justify-between p-3 border rounded-lg"
                        >
                          <div>
                            <p className="text-sm font-medium">#{offer.tokenId}</p>
                            <p className="text-xs text-muted-foreground">
                              {getMarketplaceIcon(offer.marketplace)} {offer.marketplace}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm font-medium">{formatETH(offer.price)}</p>
                            <p className="text-xs text-muted-foreground">
                              {formatTimeRemaining(offer.expirationTime)}
                            </p>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <DollarSign className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                      <p className="text-sm text-muted-foreground">No active offers</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="marketplaces" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Supported Marketplaces</CardTitle>
              <CardDescription>
                NFT marketplaces integrated with the platform
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {marketplaces.map((marketplace, index) => (
                  <motion.div
                    key={marketplace.name}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center text-2xl">
                        {getMarketplaceIcon(marketplace.name)}
                      </div>
                      <div>
                        <h4 className="font-medium">{marketplace.name}</h4>
                        <p className="text-sm text-muted-foreground">
                          {marketplace.feePercentage}% fee
                        </p>
                      </div>
                    </div>

                    <div className="space-y-2 mb-4">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Supported Orders</span>
                        <span className="font-medium">{marketplace.supportedOrderTypes.length}</span>
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {marketplace.supportedOrderTypes.map(orderType => (
                          <Badge key={orderType} variant="outline" className="text-xs">
                            {getOrderTypeIcon(orderType)}
                            <span className="ml-1 capitalize">{orderType.replace('_', ' ')}</span>
                          </Badge>
                        ))}
                      </div>
                    </div>

                    <div className="flex gap-2">
                      <Button size="sm" variant="outline" className="flex-1">
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

        <TabsContent value="listings" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Active Listings</CardTitle>
              <CardDescription>
                Current NFT listings across all marketplaces
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {listings.map((listing, index) => (
                  <motion.div
                    key={listing.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                        {getOrderTypeIcon(listing.orderType)}
                      </div>
                      <div>
                        <h4 className="font-medium">NFT #{listing.tokenId}</h4>
                        <p className="text-sm text-muted-foreground">
                          {getMarketplaceIcon(listing.marketplace)} {listing.marketplace}
                        </p>
                        <div className="flex items-center gap-2 mt-1">
                          <Badge className={getStatusColor(listing.status)}>
                            {listing.status}
                          </Badge>
                          <Badge variant="outline" className="text-xs">
                            {listing.orderType.replace('_', ' ')}
                          </Badge>
                        </div>
                      </div>
                    </div>

                    <div className="text-right">
                      <p className="font-bold text-lg">{formatETH(listing.price)}</p>
                      <p className="text-sm text-muted-foreground">
                        {listing.endTime ? formatTimeRemaining(listing.endTime) : 'No expiry'}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Fee: {formatETH(listing.fees.total)}
                      </p>
                    </div>

                    <div className="flex gap-2">
                      <Button size="sm" variant="outline">
                        <Eye className="w-3 h-3 mr-2" />
                        View
                      </Button>
                      {userAddress && listing.seller.toLowerCase() === userAddress.toLowerCase() && (
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => cancelListing(listing.id)}
                        >
                          <Minus className="w-3 h-3 mr-2" />
                          Cancel
                        </Button>
                      )}
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="offers" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Active Offers</CardTitle>
              <CardDescription>
                Current NFT offers across all marketplaces
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {offers.map((offer, index) => (
                  <motion.div
                    key={offer.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                        <DollarSign className="w-6 h-6" />
                      </div>
                      <div>
                        <h4 className="font-medium">NFT #{offer.tokenId}</h4>
                        <p className="text-sm text-muted-foreground">
                          {getMarketplaceIcon(offer.marketplace)} {offer.marketplace}
                        </p>
                        <Badge className={getStatusColor(offer.status)}>
                          {offer.status}
                        </Badge>
                      </div>
                    </div>

                    <div className="text-right">
                      <p className="font-bold text-lg">{formatETH(offer.price)}</p>
                      <p className="text-sm text-muted-foreground">
                        Expires: {formatTimeRemaining(offer.expirationTime)}
                      </p>
                    </div>

                    <div className="flex gap-2">
                      <Button size="sm" variant="outline">
                        <Eye className="w-3 h-3 mr-2" />
                        View
                      </Button>
                      {userAddress && offer.buyer.toLowerCase() === userAddress.toLowerCase() && (
                        <Button 
                          size="sm" 
                          variant="outline"
                          onClick={() => cancelOffer(offer.id)}
                        >
                          <Minus className="w-3 h-3 mr-2" />
                          Cancel
                        </Button>
                      )}
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="activity" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Marketplace Activity</CardTitle>
              <CardDescription>
                Recent transactions across all marketplaces
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentActivity.map((activity, index) => (
                  <motion.div
                    key={activity.id}
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                        <Activity className="w-6 h-6" />
                      </div>
                      <div>
                        <h4 className="font-medium capitalize">{activity.type}</h4>
                        <p className="text-sm text-muted-foreground">
                          {getMarketplaceIcon(activity.marketplace)} {activity.marketplace} • NFT #{activity.tokenId}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {new Date(activity.timestamp).toLocaleString()}
                        </p>
                      </div>
                    </div>

                    <div className="text-right">
                      {activity.price && (
                        <p className="font-bold text-lg">{formatETH(activity.price)}</p>
                      )}
                      <Button size="sm" variant="outline">
                        <ExternalLink className="w-3 h-3 mr-2" />
                        View Tx
                      </Button>
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
                  Marketplace Stats
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {marketplaceStats.map((stats, index) => (
                    <div key={stats.marketplace} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-lg">{getMarketplaceIcon(stats.marketplace)}</span>
                          <span className="text-sm font-medium">{stats.marketplace}</span>
                        </div>
                        <span className="text-sm font-medium">{formatETH(stats.volume24h)}</span>
                      </div>
                      <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground">
                        <span>Sales: {stats.sales24h}</span>
                        <span>Avg: {formatETH(stats.averagePrice24h)}</span>
                        <span>Floor: {formatETH(stats.floorPrice)}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Platform Overview
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Marketplaces</span>
                    <span className="font-medium">{marketplaces.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Active Listings</span>
                    <span className="font-medium">{listings.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Active Offers</span>
                    <span className="font-medium">{offers.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">24h Volume</span>
                    <span className="font-medium">{formatETH(marketplaceAnalytics.totalVolume24h)}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">24h Sales</span>
                    <span className="font-medium">{marketplaceAnalytics.totalSales24h}</span>
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
