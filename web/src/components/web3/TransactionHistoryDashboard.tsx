'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { 
  History, 
  Search, 
  Filter,
  Download,
  RefreshCw,
  TrendingUp,
  TrendingDown,
  Activity,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  ExternalLink,
  Copy,
  Calendar,
  BarChart3,
  PieChart,
  ArrowUpRight,
  ArrowDownLeft,
  Repeat,
  Zap
} from 'lucide-react'
import { useTransactionHistory } from '@/hooks/useTransactionHistory'
import { TransactionType, TransactionStatus, TransactionCategory } from '@/lib/transaction-history'
import { type Address } from 'viem'

interface TransactionHistoryDashboardProps {
  userAddress?: Address
  chainId?: number
}

export function TransactionHistoryDashboard({ userAddress, chainId = 1 }: TransactionHistoryDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedType, setSelectedType] = useState<TransactionType | 'all'>('all')
  const [selectedCategory, setSelectedCategory] = useState<TransactionCategory | 'all'>('all')
  const [selectedStatus, setSelectedStatus] = useState<TransactionStatus | 'all'>('all')

  const {
    transactions,
    summary,
    isLoading,
    isSearching,
    error,
    hasMore,
    totalCount,
    loadTransactions,
    loadMore,
    searchTransactions,
    clearSearch,
    exportTransactions,
    analytics,
    recentTransactions,
    failedTransactions,
    pendingTransactions,
    confirmedCount,
    pendingCount,
    failedCount,
    formatValue,
    formatGasPrice,
    getTypeIcon,
    getStatusColor,
    getCategoryColor,
    isSearchActive
  } = useTransactionHistory({
    address: userAddress,
    chainId,
    autoRefresh: true,
    enableNotifications: true
  })

  const handleSearch = async () => {
    const filters: any = {}

    if (searchQuery.trim()) {
      filters.searchQuery = searchQuery.trim()
    }

    if (selectedType !== 'all') {
      filters.type = [selectedType]
    }

    if (selectedCategory !== 'all') {
      filters.category = [selectedCategory]
    }

    if (selectedStatus !== 'all') {
      filters.status = [selectedStatus]
    }

    if (Object.keys(filters).length > 0) {
      await searchTransactions(filters)
    } else {
      await clearSearch()
    }
  }

  const handleClearSearch = async () => {
    setSearchQuery('')
    setSelectedType('all')
    setSelectedCategory('all')
    setSelectedStatus('all')
    await clearSearch()
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  const getStatusIcon = (status: TransactionStatus) => {
    switch (status) {
      case TransactionStatus.CONFIRMED:
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case TransactionStatus.PENDING:
        return <Clock className="w-4 h-4 text-yellow-500" />
      case TransactionStatus.FAILED:
        return <XCircle className="w-4 h-4 text-red-500" />
      default:
        return <AlertCircle className="w-4 h-4 text-gray-500" />
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <History className="w-6 h-6" />
            Transaction History & Search
          </h2>
          <p className="text-muted-foreground">
            Advanced transaction tracking with search, filtering, and analytics
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => exportTransactions('csv')}>
            <Download className="w-4 h-4 mr-2" />
            Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={() => loadTransactions()}>
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Transactions</p>
                <p className="text-2xl font-bold">{totalCount}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold">{analytics.successRate.toFixed(1)}%</p>
                <p className="text-xs text-muted-foreground">
                  {confirmedCount} confirmed
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Gas Cost</p>
                <p className="text-2xl font-bold">{formatValue(analytics.totalGasCost.toString())} ETH</p>
                <p className="text-xs text-muted-foreground">
                  Avg: {formatGasPrice(analytics.averageGasPrice.toString())}
                </p>
              </div>
              <Zap className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Pending</p>
                <p className="text-2xl font-bold">{pendingCount}</p>
                <p className="text-xs text-muted-foreground">
                  {failedCount} failed
                </p>
              </div>
              <Clock className="w-8 h-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Search className="w-5 h-5" />
            Search & Filter
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            <div className="md:col-span-2">
              <Label htmlFor="search">Search</Label>
              <Input
                id="search"
                placeholder="Hash, address, contract, or description..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>

            <div>
              <Label htmlFor="type">Type</Label>
              <Select value={selectedType} onValueChange={(value) => setSelectedType(value as any)}>
                <SelectTrigger>
                  <SelectValue placeholder="All types" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  {Object.values(TransactionType).map(type => (
                    <SelectItem key={type} value={type}>
                      {getTypeIcon(type)} {type.replace('_', ' ')}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label htmlFor="category">Category</Label>
              <Select value={selectedCategory} onValueChange={(value) => setSelectedCategory(value as any)}>
                <SelectTrigger>
                  <SelectValue placeholder="All categories" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Categories</SelectItem>
                  {Object.values(TransactionCategory).map(category => (
                    <SelectItem key={category} value={category}>
                      {category.replace('_', ' ')}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-end gap-2">
              <Button onClick={handleSearch} disabled={isSearching}>
                {isSearching ? (
                  <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                ) : (
                  <Search className="w-4 h-4 mr-2" />
                )}
                Search
              </Button>
              {isSearchActive && (
                <Button variant="outline" onClick={handleClearSearch}>
                  Clear
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="pending">Pending</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Recent Transactions */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Transactions</CardTitle>
              <CardDescription>
                Your latest transaction activity
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentTransactions.map((tx, index) => (
                  <motion.div
                    key={tx.hash}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-3 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-3">
                      <div className="text-2xl">{getTypeIcon(tx.type)}</div>
                      <div>
                        <div className="flex items-center gap-2">
                          <p className="font-medium">{tx.description || tx.type.replace('_', ' ')}</p>
                          <Badge className={getCategoryColor(tx.category)}>
                            {tx.category}
                          </Badge>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          {getStatusIcon(tx.status)}
                          <span>{new Date(tx.timestamp).toLocaleDateString()}</span>
                          <span>•</span>
                          <span>{tx.hash.slice(0, 10)}...</span>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">{formatValue(tx.value)} ETH</p>
                      <p className="text-sm text-muted-foreground">
                        Gas: {tx.gasCostETH} ETH
                      </p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Quick Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5 text-green-500" />
                  Most Active Type
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {Object.entries(analytics.typeDistribution)
                    .sort(([,a], [,b]) => b - a)
                    .slice(0, 3)
                    .map(([type, count]) => (
                      <div key={type} className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span>{getTypeIcon(type as TransactionType)}</span>
                          <span className="text-sm">{type.replace('_', ' ')}</span>
                        </div>
                        <span className="font-medium">{count}</span>
                      </div>
                    ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5 text-blue-500" />
                  Category Breakdown
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {Object.entries(analytics.categoryDistribution)
                    .sort(([,a], [,b]) => b - a)
                    .slice(0, 3)
                    .map(([category, count]) => (
                      <div key={category} className="flex items-center justify-between">
                        <Badge className={getCategoryColor(category as TransactionCategory)}>
                          {category}
                        </Badge>
                        <span className="font-medium">{count}</span>
                      </div>
                    ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5 text-purple-500" />
                  Gas Efficiency
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Average Gas Price</span>
                    <span className="font-medium">
                      {formatGasPrice(analytics.averageGasPrice.toString())}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Gas Cost</span>
                    <span className="font-medium">
                      {formatValue(analytics.totalGasCost.toString())} ETH
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Success Rate</span>
                    <span className="font-medium text-green-600">
                      {analytics.successRate.toFixed(1)}%
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="transactions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>All Transactions</CardTitle>
              <CardDescription>
                {isSearchActive ? `Search results (${totalCount} found)` : `Complete transaction history (${totalCount} total)`}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {transactions.map((tx, index) => (
                  <motion.div
                    key={tx.hash}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="text-2xl">{getTypeIcon(tx.type)}</div>
                        <div>
                          <div className="flex items-center gap-2">
                            <h4 className="font-medium">{tx.description || tx.type.replace('_', ' ')}</h4>
                            <Badge className={getCategoryColor(tx.category)}>
                              {tx.category}
                            </Badge>
                            {tx.protocol && (
                              <Badge variant="outline">{tx.protocol}</Badge>
                            )}
                          </div>
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            {getStatusIcon(tx.status)}
                            <span className={getStatusColor(tx.status)}>{tx.status}</span>
                            <span>•</span>
                            <span>{new Date(tx.timestamp).toLocaleString()}</span>
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-bold text-lg">{formatValue(tx.value)} ETH</p>
                        <p className="text-sm text-muted-foreground">
                          Gas: {tx.gasCostETH} ETH
                        </p>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                      <div>
                        <p className="text-muted-foreground">From</p>
                        <div className="flex items-center gap-1">
                          <p className="font-mono">{tx.from.slice(0, 10)}...</p>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-1"
                            onClick={() => copyToClipboard(tx.from)}
                          >
                            <Copy className="w-3 h-3" />
                          </Button>
                        </div>
                      </div>

                      <div>
                        <p className="text-muted-foreground">To</p>
                        <div className="flex items-center gap-1">
                          <p className="font-mono">{tx.to?.slice(0, 10) || 'Contract'}...</p>
                          {tx.to && (
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-auto p-1"
                              onClick={() => copyToClipboard(tx.to!)}
                            >
                              <Copy className="w-3 h-3" />
                            </Button>
                          )}
                        </div>
                      </div>

                      <div>
                        <p className="text-muted-foreground">Block</p>
                        <p className="font-mono">{tx.blockNumber.toString()}</p>
                      </div>

                      <div>
                        <p className="text-muted-foreground">Gas Price</p>
                        <p>{formatGasPrice(tx.gasPrice)}</p>
                      </div>
                    </div>

                    {/* Token Transfers */}
                    {tx.tokenTransfers.length > 0 && (
                      <div className="mt-3 pt-3 border-t">
                        <p className="text-sm font-medium mb-2">Token Transfers</p>
                        <div className="space-y-1">
                          {tx.tokenTransfers.map((transfer, i) => (
                            <div key={i} className="flex items-center justify-between text-sm">
                              <span>{transfer.valueFormatted} {transfer.tokenSymbol}</span>
                              <span className="text-muted-foreground">
                                {transfer.from.slice(0, 6)}...→{transfer.to.slice(0, 6)}...
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* NFT Transfers */}
                    {tx.nftTransfers.length > 0 && (
                      <div className="mt-3 pt-3 border-t">
                        <p className="text-sm font-medium mb-2">NFT Transfers</p>
                        <div className="space-y-1">
                          {tx.nftTransfers.map((transfer, i) => (
                            <div key={i} className="flex items-center justify-between text-sm">
                              <span>{transfer.tokenName} #{transfer.tokenId}</span>
                              <span className="text-muted-foreground">
                                {transfer.from.slice(0, 6)}...→{transfer.to.slice(0, 6)}...
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    <div className="flex items-center justify-between mt-3 pt-3 border-t">
                      <div className="flex items-center gap-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => copyToClipboard(tx.hash)}
                        >
                          <Copy className="w-3 h-3 mr-2" />
                          Copy Hash
                        </Button>
                        <Button variant="ghost" size="sm">
                          <ExternalLink className="w-3 h-3 mr-2" />
                          View on Explorer
                        </Button>
                      </div>
                      <div className="flex items-center gap-1 text-xs text-muted-foreground">
                        {tx.tags.map(tag => (
                          <Badge key={tag} variant="secondary" className="text-xs">
                            {tag}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </motion.div>
                ))}

                {hasMore && (
                  <div className="text-center pt-4">
                    <Button variant="outline" onClick={loadMore} disabled={isLoading}>
                      {isLoading ? (
                        <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                      ) : (
                        'Load More'
                      )}
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Transaction Types</CardTitle>
                <CardDescription>
                  Distribution of transaction types
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {Object.entries(analytics.typeDistribution)
                    .sort(([,a], [,b]) => b - a)
                    .map(([type, count]) => {
                      const percentage = analytics.totalTransactions > 0 
                        ? (count / analytics.totalTransactions) * 100 
                        : 0
                      return (
                        <div key={type} className="space-y-1">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <span>{getTypeIcon(type as TransactionType)}</span>
                              <span className="text-sm font-medium">{type.replace('_', ' ')}</span>
                            </div>
                            <span className="text-sm text-muted-foreground">
                              {count} ({percentage.toFixed(1)}%)
                            </span>
                          </div>
                          <div className="w-full bg-muted rounded-full h-2">
                            <div 
                              className="bg-primary h-2 rounded-full transition-all duration-300"
                              style={{ width: `${percentage}%` }}
                            />
                          </div>
                        </div>
                      )
                    })}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Categories</CardTitle>
                <CardDescription>
                  Transaction categories breakdown
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {Object.entries(analytics.categoryDistribution)
                    .sort(([,a], [,b]) => b - a)
                    .map(([category, count]) => {
                      const percentage = analytics.totalTransactions > 0 
                        ? (count / analytics.totalTransactions) * 100 
                        : 0
                      return (
                        <div key={category} className="space-y-1">
                          <div className="flex items-center justify-between">
                            <Badge className={getCategoryColor(category as TransactionCategory)}>
                              {category}
                            </Badge>
                            <span className="text-sm text-muted-foreground">
                              {count} ({percentage.toFixed(1)}%)
                            </span>
                          </div>
                          <div className="w-full bg-muted rounded-full h-2">
                            <div 
                              className="bg-primary h-2 rounded-full transition-all duration-300"
                              style={{ width: `${percentage}%` }}
                            />
                          </div>
                        </div>
                      )
                    })}
                </div>
              </CardContent>
            </Card>
          </div>

          {summary && (
            <Card>
              <CardHeader>
                <CardTitle>Summary Statistics</CardTitle>
                <CardDescription>
                  Overall transaction performance metrics
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{summary.totalTransactions}</p>
                    <p className="text-sm text-muted-foreground">Total Transactions</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{formatValue(summary.totalValueETH)} ETH</p>
                    <p className="text-sm text-muted-foreground">Total Value</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{formatValue(summary.totalGasCostETH)} ETH</p>
                    <p className="text-sm text-muted-foreground">Total Gas Cost</p>
                  </div>
                  <div className="text-center p-4 border rounded-lg">
                    <p className="text-2xl font-bold">{summary.successRate.toFixed(1)}%</p>
                    <p className="text-sm text-muted-foreground">Success Rate</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="pending" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="w-5 h-5 text-yellow-500" />
                  Pending Transactions
                </CardTitle>
                <CardDescription>
                  Transactions waiting for confirmation
                </CardDescription>
              </CardHeader>
              <CardContent>
                {pendingTransactions.length > 0 ? (
                  <div className="space-y-3">
                    {pendingTransactions.map(tx => (
                      <div key={tx.hash} className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">{tx.description || tx.type}</p>
                          <p className="text-sm text-muted-foreground">{tx.hash.slice(0, 20)}...</p>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatValue(tx.value)} ETH</p>
                          <p className="text-sm text-yellow-600">Pending</p>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground">No pending transactions</p>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <XCircle className="w-5 h-5 text-red-500" />
                  Failed Transactions
                </CardTitle>
                <CardDescription>
                  Transactions that failed to execute
                </CardDescription>
              </CardHeader>
              <CardContent>
                {failedTransactions.length > 0 ? (
                  <div className="space-y-3">
                    {failedTransactions.map(tx => (
                      <div key={tx.hash} className="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <p className="font-medium">{tx.description || tx.type}</p>
                          <p className="text-sm text-muted-foreground">{tx.hash.slice(0, 20)}...</p>
                          {tx.errorReason && (
                            <p className="text-xs text-red-600">{tx.errorReason}</p>
                          )}
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatValue(tx.value)} ETH</p>
                          <p className="text-sm text-red-600">Failed</p>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-4" />
                    <p className="text-muted-foreground">No failed transactions</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
