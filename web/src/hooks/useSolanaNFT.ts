import { useState, useEffect, useCallback } from 'react'
import { SolanaNFTService } from '@/services/solana/SolanaNFTService'

export interface NFTCollection {
  id: string
  name: string
  symbol: string
  description: string
  image: string
  floorPrice: number
  floorPriceChange24h: number
  volume24h: number
  volumeChange24h: number
  marketCap: number
  supply: number
  owners: number
  listedCount: number
  listedPercentage: number
  averagePrice24h: number
  sales24h: number
  website?: string
  twitter?: string
  discord?: string
  verified: boolean
  marketplace: 'magic-eden' | 'tensor' | 'solanart' | 'exchange-art'
}

export interface NFTMarketStats {
  totalVolume24h: number
  totalSales24h: number
  totalCollections: number
  averageFloorPrice: number
  topCollectionByVolume: string
  topCollectionByFloor: string
  marketCapTotal: number
}

export interface NFTTrend {
  collection: string
  metric: 'volume' | 'floor_price' | 'sales'
  change24h: number
  changePercentage: number
  rank: number
}

export interface SolanaNFTState {
  collections: NFTCollection[]
  marketStats: NFTMarketStats | null
  trends: NFTTrend[]
  isLoading: boolean
  error: string | null
  lastUpdated: Date | null
}

export interface UseSolanaNFTOptions {
  autoRefresh?: boolean
  refreshInterval?: number
  marketplace?: NFTCollection['marketplace']
  minFloorPrice?: number
  maxFloorPrice?: number
  sortBy?: 'volume' | 'floor_price' | 'market_cap' | 'sales'
  limit?: number
}

export function useSolanaNFT(options: UseSolanaNFTOptions = {}) {
  const {
    autoRefresh = false,
    refreshInterval = 120000, // 2 minutes for NFT data
    marketplace,
    minFloorPrice = 0,
    maxFloorPrice,
    sortBy = 'volume',
    limit = 50
  } = options

  const [state, setState] = useState<SolanaNFTState>({
    collections: [],
    marketStats: null,
    trends: [],
    isLoading: true,
    error: null,
    lastUpdated: null
  })

  const nftService = new SolanaNFTService()

  const fetchNFTData = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      const [collectionsData, statsData, trendsData] = await Promise.all([
        nftService.getCollections({
          marketplace,
          minFloorPrice,
          maxFloorPrice,
          sortBy,
          limit
        }),
        nftService.getMarketStats(),
        nftService.getTrends()
      ])

      setState(prev => ({
        ...prev,
        collections: collectionsData,
        marketStats: statsData,
        trends: trendsData,
        isLoading: false,
        lastUpdated: new Date()
      }))
    } catch (error) {
      console.error('Failed to fetch Solana NFT data:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to fetch NFT data'
      }))
    }
  }, [marketplace, minFloorPrice, maxFloorPrice, sortBy, limit])

  const refresh = useCallback(async () => {
    await fetchNFTData()
  }, [fetchNFTData])

  // Initial data fetch
  useEffect(() => {
    fetchNFTData()
  }, [fetchNFTData])

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(fetchNFTData, refreshInterval)
    return () => clearInterval(interval)
  }, [autoRefresh, refreshInterval, fetchNFTData])

  // Computed values
  const totalVolume = state.marketStats?.totalVolume24h || 0
  const topCollections = state.collections
    .sort((a, b) => {
      switch (sortBy) {
        case 'volume':
          return b.volume24h - a.volume24h
        case 'floor_price':
          return b.floorPrice - a.floorPrice
        case 'market_cap':
          return b.marketCap - a.marketCap
        case 'sales':
          return b.sales24h - a.sales24h
        default:
          return b.volume24h - a.volume24h
      }
    })
    .slice(0, 10)

  const collectionsByMarketplace = state.collections.reduce((acc, collection) => {
    if (!acc[collection.marketplace]) {
      acc[collection.marketplace] = []
    }
    acc[collection.marketplace].push(collection)
    return acc
  }, {} as Record<NFTCollection['marketplace'], NFTCollection[]>)

  const trendingUp = state.trends
    .filter(trend => trend.changePercentage > 0)
    .sort((a, b) => b.changePercentage - a.changePercentage)
    .slice(0, 5)

  const trendingDown = state.trends
    .filter(trend => trend.changePercentage < 0)
    .sort((a, b) => a.changePercentage - b.changePercentage)
    .slice(0, 5)

  return {
    ...state,
    refresh,
    totalVolume,
    topCollections,
    collectionsByMarketplace,
    trendingUp,
    trendingDown
  }
}

// Helper hook for specific collection data
export function useNFTCollection(collectionId: string) {
  const [collection, setCollection] = useState<NFTCollection | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const nftService = new SolanaNFTService()

  const fetchCollection = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const collectionData = await nftService.getCollection(collectionId)
      setCollection(collectionData)
    } catch (error) {
      console.error(`Failed to fetch collection ${collectionId}:`, error)
      setError(error instanceof Error ? error.message : 'Failed to fetch collection')
    } finally {
      setIsLoading(false)
    }
  }, [collectionId])

  useEffect(() => {
    if (collectionId) {
      fetchCollection()
    }
  }, [collectionId, fetchCollection])

  return {
    collection,
    isLoading,
    error,
    refresh: fetchCollection
  }
}

// Helper hook for marketplace-specific data
export function useNFTMarketplace(marketplace: NFTCollection['marketplace']) {
  const { collectionsByMarketplace, isLoading, error, refresh } = useSolanaNFT({
    marketplace,
    autoRefresh: true,
    refreshInterval: 60000 // 1 minute
  })

  const marketplaceCollections = collectionsByMarketplace[marketplace] || []
  const marketplaceVolume = marketplaceCollections.reduce(
    (sum, collection) => sum + collection.volume24h, 
    0
  )
  const marketplaceSales = marketplaceCollections.reduce(
    (sum, collection) => sum + collection.sales24h, 
    0
  )

  return {
    collections: marketplaceCollections,
    totalVolume: marketplaceVolume,
    totalSales: marketplaceSales,
    collectionCount: marketplaceCollections.length,
    isLoading,
    error,
    refresh
  }
}

// Helper hook for NFT trends and analytics
export function useNFTTrends() {
  const { trends, isLoading, error, refresh } = useSolanaNFT({
    autoRefresh: true,
    refreshInterval: 300000 // 5 minutes
  })

  const volumeTrends = trends.filter(trend => trend.metric === 'volume')
  const floorPriceTrends = trends.filter(trend => trend.metric === 'floor_price')
  const salesTrends = trends.filter(trend => trend.metric === 'sales')

  const biggestGainers = trends
    .filter(trend => trend.changePercentage > 0)
    .sort((a, b) => b.changePercentage - a.changePercentage)
    .slice(0, 10)

  const biggestLosers = trends
    .filter(trend => trend.changePercentage < 0)
    .sort((a, b) => a.changePercentage - b.changePercentage)
    .slice(0, 10)

  return {
    trends,
    volumeTrends,
    floorPriceTrends,
    salesTrends,
    biggestGainers,
    biggestLosers,
    isLoading,
    error,
    refresh
  }
}
