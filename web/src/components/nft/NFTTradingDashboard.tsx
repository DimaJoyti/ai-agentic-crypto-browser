'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  ShoppingCart, 
  Tag, 
  Gavel, 
  TrendingUp, 
  Clock, 
  DollarSign,
  Activity,
  Users,
  BarChart3,
  RefreshCw,
  Plus,
  X,
  CheckCircle,
  AlertCircle,
  Loader2
} from 'lucide-react'
import { useNFTTrading, useNFTOrders, useNFTAuctions, useBulkNFTOperations } from '@/hooks/useNFTTrading'
import { OrderType, OrderStatus, TransactionStatus, BulkOperationStatus } from '@/lib/nft-trading'
import { formatCurrency, formatNumber, formatPercentage } from '@/lib/utils'

interface NFTTradingDashboardProps {
  className?: string
}

export function NFTTradingDashboard({ className }: NFTTradingDashboardProps) {
  const { state, getTradingStats, getOrderStats, refresh } = useNFTTrading()
  const { orders } = useNFTOrders()
  const { auctions } = useNFTAuctions()
  const { bulkOperations } = useBulkNFTOperations()

  const [activeTab, setActiveTab] = useState('overview')

  const tradingStats = getTradingStats()
  const orderStats = getOrderStats()

  const handleRefresh = () => {
    refresh()
  }

  const getStatusColor = (status: OrderStatus | TransactionStatus | BulkOperationStatus) => {
    switch (status) {
      case 'active':
      case 'pending':
      case 'running':
        return 'bg-blue-100 text-blue-800'
      case 'filled':
      case 'confirmed':
      case 'finalized':
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'cancelled':
      case 'failed':
        return 'bg-red-100 text-red-800'
      case 'expired':
        return 'bg-gray-100 text-gray-800'
      case 'partial':
        return 'bg-yellow-100 text-yellow-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusIcon = (status: OrderStatus | TransactionStatus | BulkOperationStatus) => {
    switch (status) {
      case 'active':
      case 'pending':
      case 'running':
        return <Loader2 className="h-3 w-3 animate-spin" />
      case 'filled':
      case 'confirmed':
      case 'finalized':
      case 'completed':
        return <CheckCircle className="h-3 w-3" />
      case 'cancelled':
      case 'failed':
        return <X className="h-3 w-3" />
      case 'expired':
        return <Clock className="h-3 w-3" />
      case 'partial':
        return <AlertCircle className="h-3 w-3" />
      default:
        return <Activity className="h-3 w-3" />
    }
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">NFT Trading Dashboard</h2>
          <p className="text-muted-foreground">
            Manage your NFT orders, auctions, and trading activities
          </p>
        </div>
        <div className="flex items-center gap-2">
          {state.lastUpdate && (
            <span className="text-sm text-muted-foreground">
              Last updated: {new Date(state.lastUpdate).toLocaleTimeString()}
            </span>
          )}
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={state.isLoading}
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert className="border-red-200 bg-red-50">
          <AlertCircle className="h-4 w-4 text-red-600" />
          <AlertDescription className="text-red-800">
            {state.error}
          </AlertDescription>
        </Alert>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="orders">Orders</TabsTrigger>
          <TabsTrigger value="auctions">Auctions</TabsTrigger>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
          <TabsTrigger value="bulk-ops">Bulk Operations</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          {/* Trading Stats */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Orders</CardTitle>
                <ShoppingCart className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatNumber(tradingStats.totalOrders)}</div>
                <p className="text-xs text-muted-foreground">
                  {tradingStats.activeOrders} active
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Volume</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatCurrency(tradingStats.totalVolume)} ETH</div>
                <p className="text-xs text-muted-foreground">
                  Avg: {formatCurrency(tradingStats.averageOrderValue)} ETH
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatPercentage(tradingStats.successRate)}</div>
                <p className="text-xs text-muted-foreground">
                  {tradingStats.totalTransactions} transactions
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Active Auctions</CardTitle>
                <Gavel className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{auctions.length}</div>
                <p className="text-xs text-muted-foreground">
                  Live auctions
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Order Distribution */}
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Order Types</CardTitle>
                <CardDescription>Distribution by order type</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {Object.entries(orderStats.byType).map(([type, count]) => (
                    <div key={type} className="flex items-center justify-between">
                      <span className="capitalize">{type.replace('_', ' ')}</span>
                      <div className="flex items-center gap-2">
                        <Progress 
                          value={tradingStats.totalOrders > 0 ? (count / tradingStats.totalOrders) * 100 : 0} 
                          className="w-20 h-2" 
                        />
                        <span className="text-sm font-medium w-8 text-right">{count}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Order Status</CardTitle>
                <CardDescription>Distribution by status</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {Object.entries(orderStats.byStatus).map(([status, count]) => (
                    <div key={status} className="flex items-center justify-between">
                      <span className="capitalize">{status}</span>
                      <div className="flex items-center gap-2">
                        <Progress 
                          value={tradingStats.totalOrders > 0 ? (count / tradingStats.totalOrders) * 100 : 0} 
                          className="w-20 h-2" 
                        />
                        <span className="text-sm font-medium w-8 text-right">{count}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Recent Activity */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Activity</CardTitle>
              <CardDescription>Latest order updates</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {orderStats.recentActivity.slice(0, 5).map(order => (
                  <div key={order.id} className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center text-white text-xs font-bold">
                        {order.type === OrderType.LISTING ? 'S' : order.type === OrderType.OFFER ? 'B' : 'A'}
                      </div>
                      <div>
                        <div className="font-medium">
                          {order.type === OrderType.LISTING ? 'Sell Order' : 
                           order.type === OrderType.OFFER ? 'Buy Order' : 
                           'Auction'}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          Token #{order.tokenId} • {order.marketplace}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-medium">{formatCurrency(parseFloat(order.price.amount))} ETH</div>
                      <Badge className={getStatusColor(order.status)} variant="secondary">
                        {getStatusIcon(order.status)}
                        <span className="ml-1 capitalize">{order.status}</span>
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="orders" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>All Orders</CardTitle>
              <CardDescription>Manage your buy and sell orders</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {orders.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No orders found. Create your first order to get started.
                  </div>
                ) : (
                  orders.map(order => (
                    <div key={order.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center text-white font-bold">
                          {order.type === OrderType.LISTING ? 'S' : 'B'}
                        </div>
                        <div>
                          <div className="font-medium">
                            {order.type === OrderType.LISTING ? 'Sell Order' : 'Buy Order'}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Token #{order.tokenId} • {order.marketplace}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Created: {new Date(order.createdAt).toLocaleDateString()}
                          </div>
                        </div>
                      </div>

                      <div className="text-right">
                        <div className="font-medium text-lg">
                          {formatCurrency(parseFloat(order.price.amount))} ETH
                        </div>
                        <div className="text-sm text-muted-foreground">
                          ${formatNumber(order.price.usdValue)}
                        </div>
                        <Badge className={getStatusColor(order.status)} variant="secondary">
                          {getStatusIcon(order.status)}
                          <span className="ml-1 capitalize">{order.status}</span>
                        </Badge>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="auctions" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Active Auctions</CardTitle>
              <CardDescription>Monitor and participate in NFT auctions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {auctions.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No active auctions found.
                  </div>
                ) : (
                  auctions.map(auction => (
                    <div key={auction.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 bg-gradient-to-br from-orange-500 to-red-600 rounded-lg flex items-center justify-center text-white">
                          <Gavel className="h-5 w-5" />
                        </div>
                        <div>
                          <div className="font-medium">Auction</div>
                          <div className="text-sm text-muted-foreground">
                            Token #{auction.tokenId} • {auction.marketplace}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Ends: {new Date(auction.endTime!).toLocaleString()}
                          </div>
                        </div>
                      </div>

                      <div className="text-right">
                        <div className="font-medium text-lg">
                          {formatCurrency(parseFloat(auction.price.amount))} ETH
                        </div>
                        <div className="text-sm text-muted-foreground">
                          Starting price
                        </div>
                        <Badge className={getStatusColor(auction.status)} variant="secondary">
                          {getStatusIcon(auction.status)}
                          <span className="ml-1 capitalize">{auction.status}</span>
                        </Badge>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="transactions" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Transaction History</CardTitle>
              <CardDescription>View all trading transactions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {state.transactions.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No transactions found.
                  </div>
                ) : (
                  state.transactions.map(transaction => (
                    <div key={transaction.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 bg-gradient-to-br from-green-500 to-blue-600 rounded-lg flex items-center justify-center text-white">
                          <Activity className="h-5 w-5" />
                        </div>
                        <div>
                          <div className="font-medium capitalize">{transaction.type} Transaction</div>
                          <div className="text-sm text-muted-foreground">
                            {transaction.hash ? `${transaction.hash.slice(0, 10)}...` : 'Pending'}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {new Date(transaction.createdAt).toLocaleString()}
                          </div>
                        </div>
                      </div>

                      <div className="text-right">
                        <div className="font-medium">
                          {formatCurrency(parseFloat(transaction.value))} ETH
                        </div>
                        <div className="text-sm text-muted-foreground">
                          Gas: {transaction.gasUsed || transaction.gasLimit}
                        </div>
                        <Badge className={getStatusColor(transaction.status)} variant="secondary">
                          {getStatusIcon(transaction.status)}
                          <span className="ml-1 capitalize">{transaction.status}</span>
                        </Badge>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="bulk-ops" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Bulk Operations</CardTitle>
              <CardDescription>Monitor bulk trading operations</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {bulkOperations.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No bulk operations found.
                  </div>
                ) : (
                  bulkOperations.map(operation => (
                    <div key={operation.id} className="p-4 border rounded-lg">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg flex items-center justify-center text-white">
                            <BarChart3 className="h-5 w-5" />
                          </div>
                          <div>
                            <div className="font-medium capitalize">
                              {operation.type.replace('_', ' ')}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {operation.operations.length} items
                            </div>
                          </div>
                        </div>
                        <Badge className={getStatusColor(operation.status)} variant="secondary">
                          {getStatusIcon(operation.status)}
                          <span className="ml-1 capitalize">{operation.status}</span>
                        </Badge>
                      </div>

                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>Progress</span>
                          <span>{operation.progress.completed}/{operation.progress.total}</span>
                        </div>
                        <Progress value={operation.progress.percentage} className="h-2" />
                        {operation.progress.failed > 0 && (
                          <div className="text-sm text-red-600">
                            {operation.progress.failed} failed
                          </div>
                        )}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
