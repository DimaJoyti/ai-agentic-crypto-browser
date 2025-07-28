import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  nftAnalyticsService, 
  type RarityScore,
  type FloorPriceData,
  type VolumeData,
  type MarketTrend,
  type NFTValuation,
  type CollectionAnalytics,
  type PriceHistory,
  RarityTier,
  TrendDirection 
} from '@/lib/nft-analytics'
import { toast } from 'sonner'

export interface UseNFTAnalyticsOptions {
  contractAddress?: Address
  tokenId?: string
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface NFTAnalyticsState {
  rarityScore: RarityScore | null
  floorPriceData: FloorPriceData | null
  volumeData: VolumeData | null
  marketTrend: MarketTrend | null
  nftValuation: NFTValuation | null
  collectionAnalytics: CollectionAnalytics | null
  priceHistory: PriceHistory[]
  topCollections: CollectionAnalytics[]
  trendingCollections: CollectionAnalytics[]
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useNFTAnalytics(options: UseNFTAnalyticsOptions = {}) {
  const {
    contractAddress,
    tokenId,
    autoRefresh = true,
    refreshInterval = 300000, // 5 minutes
    enableNotifications = true
  } = options

  const [state, setState] = useState<NFTAnalyticsState>({
    rarityScore: null,
    floorPriceData: null,
    volumeData: null,
    marketTrend: null,
    nftValuation: null,
    collectionAnalytics: null,
    priceHistory: [],
    topCollections: [],
    trendingCollections: [],
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load analytics data
  const loadData = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const promises: Promise<any>[] = []

      // Load collection-specific data
      if (contractAddress) {
        promises.push(
          nftAnalyticsService.getFloorPriceData(contractAddress),
          nftAnalyticsService.getVolumeData(contractAddress),
          nftAnalyticsService.getMarketTrend(contractAddress),
          nftAnalyticsService.getCollectionAnalytics(contractAddress)
        )

        // Load NFT-specific data
        if (tokenId) {
          promises.push(
            nftAnalyticsService.getRarityScore(contractAddress, tokenId),
            nftAnalyticsService.getNFTValuation(contractAddress, tokenId),
            nftAnalyticsService.getPriceHistory(contractAddress, tokenId)
          )
        }
      }

      // Load global data
      promises.push(
        nftAnalyticsService.getTopCollections(10),
        nftAnalyticsService.getTrendingCollections(10)
      )

      const results = await Promise.all(promises)
      let resultIndex = 0

      let floorPriceData = null
      let volumeData = null
      let marketTrend = null
      let collectionAnalytics = null
      let rarityScore = null
      let nftValuation = null
      let priceHistory: PriceHistory[] = []

      if (contractAddress) {
        floorPriceData = results[resultIndex++]
        volumeData = results[resultIndex++]
        marketTrend = results[resultIndex++]
        collectionAnalytics = results[resultIndex++]

        if (tokenId) {
          rarityScore = results[resultIndex++]
          nftValuation = results[resultIndex++]
          priceHistory = results[resultIndex++]
        }
      }

      const topCollections = results[resultIndex++]
      const trendingCollections = results[resultIndex++]

      setState(prev => ({
        ...prev,
        rarityScore,
        floorPriceData,
        volumeData,
        marketTrend,
        nftValuation,
        collectionAnalytics,
        priceHistory,
        topCollections,
        trendingCollections,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load analytics data'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Analytics Error', {
          description: errorMessage
        })
      }
    }
  }, [contractAddress, tokenId, enableNotifications])

  // Get rarity score for specific NFT
  const getRarityScore = useCallback(async (contractAddress: Address, tokenId: string) => {
    try {
      return await nftAnalyticsService.getRarityScore(contractAddress, tokenId)
    } catch (error) {
      console.error('Failed to get rarity score:', error)
      return null
    }
  }, [])

  // Get floor price data for collection
  const getFloorPriceData = useCallback(async (contractAddress: Address) => {
    try {
      return await nftAnalyticsService.getFloorPriceData(contractAddress)
    } catch (error) {
      console.error('Failed to get floor price data:', error)
      return null
    }
  }, [])

  // Get NFT valuation
  const getNFTValuation = useCallback(async (contractAddress: Address, tokenId: string) => {
    try {
      return await nftAnalyticsService.getNFTValuation(contractAddress, tokenId)
    } catch (error) {
      console.error('Failed to get NFT valuation:', error)
      return null
    }
  }, [])

  // Calculate analytics metrics
  const getAnalyticsMetrics = useCallback(() => {
    const { floorPriceData, volumeData, marketTrend, collectionAnalytics } = state

    if (!floorPriceData || !volumeData || !marketTrend) {
      return {
        floorPrice: '0',
        floorPriceChange: '0',
        volume24h: '0',
        sales24h: 0,
        averagePrice: '0',
        marketCap: '0',
        trendDirection: TrendDirection.STABLE,
        trendStrength: 0,
        liquidityScore: 0,
        volatility: 0
      }
    }

    return {
      floorPrice: floorPriceData.currentFloor,
      floorPriceChange: floorPriceData.changePercentage,
      volume24h: volumeData.volume24h,
      sales24h: volumeData.sales24h,
      averagePrice: volumeData.averagePrice24h,
      marketCap: collectionAnalytics?.marketCap || '0',
      trendDirection: marketTrend.trendDirection,
      trendStrength: marketTrend.trendStrength,
      liquidityScore: marketTrend.liquidityScore,
      volatility: marketTrend.volatility
    }
  }, [state])

  // Get rarity distribution
  const getRarityDistribution = useCallback(() => {
    return state.collectionAnalytics?.rarityDistribution || []
  }, [state.collectionAnalytics])

  // Get price distribution
  const getPriceDistribution = useCallback(() => {
    return state.collectionAnalytics?.priceDistribution || []
  }, [state.collectionAnalytics])

  // Get top sales
  const getTopSales = useCallback(() => {
    return state.collectionAnalytics?.topSales || []
  }, [state.collectionAnalytics])

  // Get market metrics
  const getMarketMetrics = useCallback(() => {
    return state.collectionAnalytics?.marketMetrics || {
      liquidityScore: 0,
      volatilityIndex: 0,
      trendScore: 0,
      momentumIndicator: 0
    }
  }, [state.collectionAnalytics])

  // Get predictions
  const getPredictions = useCallback(() => {
    return state.collectionAnalytics?.predictions || {
      floorPrice7d: '0',
      floorPrice30d: '0',
      confidence: 0
    }
  }, [state.collectionAnalytics])

  // Get rarity tier color
  const getRarityTierColor = useCallback((tier: RarityTier) => {
    switch (tier) {
      case RarityTier.COMMON:
        return 'bg-gray-100 text-gray-800'
      case RarityTier.UNCOMMON:
        return 'bg-green-100 text-green-800'
      case RarityTier.RARE:
        return 'bg-blue-100 text-blue-800'
      case RarityTier.EPIC:
        return 'bg-purple-100 text-purple-800'
      case RarityTier.LEGENDARY:
        return 'bg-yellow-100 text-yellow-800'
      case RarityTier.MYTHIC:
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }, [])

  // Get trend color
  const getTrendColor = useCallback((trend: TrendDirection) => {
    switch (trend) {
      case TrendDirection.UP:
        return 'text-green-600'
      case TrendDirection.DOWN:
        return 'text-red-600'
      case TrendDirection.STABLE:
        return 'text-gray-600'
      default:
        return 'text-gray-600'
    }
  }, [])

  // Format utilities
  const formatPrice = useCallback((price: string) => {
    return nftAnalyticsService.formatPrice(price)
  }, [])

  const formatPercentage = useCallback((value: string) => {
    return nftAnalyticsService.formatPercentage(value)
  }, [])

  const formatVolume = useCallback((volume: string) => {
    return nftAnalyticsService.formatVolume(volume)
  }, [])

  // Get confidence level color
  const getConfidenceColor = useCallback((confidence: number) => {
    if (confidence >= 80) return 'text-green-600'
    if (confidence >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }, [])

  // Get liquidity score color
  const getLiquidityColor = useCallback((score: number) => {
    if (score >= 80) return 'text-green-600'
    if (score >= 60) return 'text-blue-600'
    if (score >= 40) return 'text-yellow-600'
    return 'text-red-600'
  }, [])

  // Get volatility color
  const getVolatilityColor = useCallback((volatility: number) => {
    if (volatility >= 70) return 'text-red-600'
    if (volatility >= 40) return 'text-yellow-600'
    return 'text-green-600'
  }, [])

  // Auto-refresh setup
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadData, refreshInterval)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, loadData])

  // Initial load
  useEffect(() => {
    loadData()
  }, [loadData])

  return {
    // State
    ...state,

    // Actions
    loadData,
    getRarityScore,
    getFloorPriceData,
    getNFTValuation,

    // Analytics
    getAnalyticsMetrics,
    getRarityDistribution,
    getPriceDistribution,
    getTopSales,
    getMarketMetrics,
    getPredictions,

    // Utilities
    getRarityTierColor,
    getTrendColor,
    getConfidenceColor,
    getLiquidityColor,
    getVolatilityColor,
    formatPrice,
    formatPercentage,
    formatVolume,

    // Computed values
    analyticsMetrics: getAnalyticsMetrics(),
    rarityDistribution: getRarityDistribution(),
    priceDistribution: getPriceDistribution(),
    topSales: getTopSales(),
    marketMetrics: getMarketMetrics(),
    predictions: getPredictions(),

    // Quick access to trend data
    isUpTrend: state.marketTrend?.trendDirection === TrendDirection.UP,
    isDownTrend: state.marketTrend?.trendDirection === TrendDirection.DOWN,
    isStableTrend: state.marketTrend?.trendDirection === TrendDirection.STABLE,

    // Quick access to rarity data
    isRareNFT: state.rarityScore?.tier === RarityTier.RARE || 
               state.rarityScore?.tier === RarityTier.EPIC || 
               state.rarityScore?.tier === RarityTier.LEGENDARY || 
               state.rarityScore?.tier === RarityTier.MYTHIC,

    // Quick access to market health
    isHealthyMarket: state.marketTrend?.liquidityScore && state.marketTrend.liquidityScore >= 70,
    isVolatileMarket: state.marketTrend?.volatility && state.marketTrend.volatility >= 60
  }
}
