'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  RefreshCw,
  Activity,
  BarChart3,
  Coins,
  Image,
  ArrowUpRight,
  ArrowDownRight,
  Clock,
  Loader2
} from 'lucide-react'
import { useSolanaMarketData } from '@/hooks/useSolanaMarketData'
import { useSolanaDeFi } from '@/hooks/useSolanaDeFi'
import { useSolanaNFT } from '@/hooks/useSolanaNFT'
import { formatCurrency, formatNumber, formatPercentage } from '@/lib/utils'
import { cn } from '@/lib/utils'
import { SolanaDeFiInterface } from './SolanaDeFiInterface'
import { SolanaNFTMarketplace } from './SolanaNFTMarketplace'

interface SolanaMarketDashboardProps {
  className?: string
  autoRefresh?: boolean
  refreshInterval?: number
}

export function SolanaMarketDashboard({
  className,
  autoRefresh = true,
  refreshInterval = 30000
}: SolanaMarketDashboardProps) {
  const [mounted, setMounted] = useState(false)
  const [activeTab, setActiveTab] = useState('overview')
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  const {
    solPrice,
    marketData,
    isLoading: marketLoading,
    error: marketError,
    refresh: refreshMarket
  } = useSolanaMarketData({ autoRefresh, refreshInterval })

  const {
    protocols,
    totalTVL,
    topYields,
    isLoading: defiLoading,
    error: defiError,
    refresh: refreshDeFi
  } = useSolanaDeFi({ autoRefresh, refreshInterval })

  const {
    collections,
    totalVolume,
    topCollections,
    isLoading: nftLoading,
    error: nftError,
    refresh: refreshNFT
  } = useSolanaNFT({ autoRefresh, refreshInterval })

  const isLoading = marketLoading || defiLoading || nftLoading
  const hasError = marketError || defiError || nftError

  useEffect(() => {
    setMounted(true)
  }, [])

  const handleRefresh = async () => {
    await Promise.all([
      refreshMarket(),
      refreshDeFi(),
      refreshNFT()
    ])
    setLastUpdated(new Date())
  }

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(() => {
      handleRefresh()
    }, refreshInterval)

    return () => clearInterval(interval)
  }, [autoRefresh, refreshInterval])

  const priceChangeColor = (change: number) => {
    if (change > 0) return 'text-green-600'
    if (change < 0) return 'text-red-600'
    return 'text-gray-600'
  }

  const priceChangeIcon = (change: number) => {
    if (change > 0) return <ArrowUpRight className="h-4 w-4" />
    if (change < 0) return <ArrowDownRight className="h-4 w-4" />
    return null
  }

  // Prevent hydration mismatch by not rendering until mounted
  if (!mounted) {
    return (
      <div className="space-y-6 flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
          <p className="text-sm text-muted-foreground">Loading Solana market data...</p>
        </div>
      </div>
    )
  }

  return (
    <div className={cn('space-y-6', className)}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Solana Market</h2>
          <p className="text-muted-foreground">
            Real-time market data, DeFi protocols, and NFT collections
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Badge variant="outline" className="text-xs">
            <Clock className="h-3 w-3 mr-1" />
            Updated {lastUpdated.toLocaleTimeString()}
          </Badge>
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

      {/* Error Alert */}
      {hasError && (
        <Alert variant="destructive">
          <AlertDescription>
            Failed to load market data. Please try refreshing.
          </AlertDescription>
        </Alert>
      )}

      {/* SOL Price Overview */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Coins className="h-5 w-5 mr-2" />
            Solana (SOL)
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">Price</p>
              <div className="flex items-center space-x-2">
                <span className="text-2xl font-bold">
                  {formatCurrency(solPrice?.price || 0)}
                </span>
                {solPrice?.change24h && (
                  <div className={cn("flex items-center", priceChangeColor(solPrice.change24h))}>
                    {priceChangeIcon(solPrice.change24h)}
                    <span className="text-sm font-medium">
                      {formatPercentage(Math.abs(solPrice.change24h))}
                    </span>
                  </div>
                )}
              </div>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">Market Cap</p>
              <p className="text-xl font-semibold">
                {formatCurrency(marketData?.marketCap || 0)}
              </p>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">24h Volume</p>
              <p className="text-xl font-semibold">
                {formatCurrency(marketData?.volume24h || 0)}
              </p>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">Circulating Supply</p>
              <p className="text-xl font-semibold">
                {formatNumber(marketData?.circulatingSupply || 0, 0)} SOL
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabs for different sections */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="defi">DeFi Protocols</TabsTrigger>
          <TabsTrigger value="nft">NFT Collections</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* DeFi TVL Summary */}
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total DeFi TVL</CardTitle>
                <BarChart3 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {formatCurrency(totalTVL || 0)}
                </div>
                <p className="text-xs text-muted-foreground">
                  Across {protocols?.length || 0} protocols
                </p>
              </CardContent>
            </Card>

            {/* NFT Volume Summary */}
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">NFT 24h Volume</CardTitle>
                <Image className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {formatNumber(totalVolume || 0, 0)} SOL
                </div>
                <p className="text-xs text-muted-foreground">
                  {collections?.length || 0} active collections
                </p>
              </CardContent>
            </Card>

            {/* Network Activity */}
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Network Activity</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {formatNumber(marketData?.transactions24h || 0, 0)}
                </div>
                <p className="text-xs text-muted-foreground">
                  Transactions (24h)
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Top Performers */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Top DeFi Yields</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {topYields?.slice(0, 5).map((yield_, index) => (
                    <div key={index} className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <Badge variant="outline">{yield_.protocol}</Badge>
                        <span className="text-sm">{yield_.pool}</span>
                      </div>
                      <span className="font-semibold text-green-600">
                        {formatPercentage(yield_.apy)}
                      </span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Top NFT Collections</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {topCollections?.slice(0, 5).map((collection, index) => (
                    <div key={index} className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <span className="text-sm font-medium">{collection.name}</span>
                      </div>
                      <div className="text-right">
                        <div className="font-semibold">
                          {formatNumber(collection.floorPrice)} SOL
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {formatNumber(collection.volume24h, 0)} SOL vol
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="defi" className="space-y-4">
          <SolanaDeFiInterface />
        </TabsContent>

        <TabsContent value="nft" className="space-y-4">
          <SolanaNFTMarketplace />
        </TabsContent>
      </Tabs>

      {/* Loading Overlay */}
      {isLoading && (
        <div className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="flex items-center space-x-2">
            <Loader2 className="h-6 w-6 animate-spin" />
            <span>Loading market data...</span>
          </div>
        </div>
      )}
    </div>
  )
}
