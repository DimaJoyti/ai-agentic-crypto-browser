import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  nftService, 
  type NFT, 
  type NFTCollection,
  type UserNFTPortfolio,
  type NFTValuation,
  type NFTActivity,
  NFTRarity 
} from '@/lib/nft-service'
import { toast } from 'sonner'

export interface UseNFTCollectionOptions {
  userAddress?: Address
  chainId?: number
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface NFTCollectionState {
  collections: NFTCollection[]
  userNFTs: NFT[]
  userPortfolio: UserNFTPortfolio | null
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useNFTCollection(options: UseNFTCollectionOptions = {}) {
  const {
    userAddress,
    chainId,
    autoRefresh = true,
    refreshInterval = 300000, // 5 minutes
    enableNotifications = true
  } = options

  const [state, setState] = useState<NFTCollectionState>({
    collections: [],
    userNFTs: [],
    userPortfolio: null,
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load collections and user data
  const loadData = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const collections = nftService.getAllCollections(chainId)
      const userNFTs = userAddress ? nftService.getUserNFTs(userAddress, chainId) : []
      const userPortfolio = userAddress ? nftService.getUserPortfolio(userAddress) : null

      setState(prev => ({
        ...prev,
        collections,
        userNFTs,
        userPortfolio,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load NFT data'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
    }
  }, [chainId, userAddress])

  // Get collection by ID
  const getCollection = useCallback((collectionId: string): NFTCollection | undefined => {
    return nftService.getCollection(collectionId)
  }, [])

  // Get collection by contract address
  const getCollectionByAddress = useCallback((contractAddress: Address, chainId: number): NFTCollection | undefined => {
    return nftService.getCollectionByAddress(contractAddress, chainId)
  }, [])

  // Get NFT by ID
  const getNFT = useCallback((nftId: string): NFT | undefined => {
    return nftService.getNFT(nftId)
  }, [])

  // Get NFT valuation
  const getNFTValuation = useCallback(async (nftId: string): Promise<NFTValuation | null> => {
    try {
      return await nftService.getNFTValuation(nftId)
    } catch (error) {
      console.error('Failed to get NFT valuation:', error)
      return null
    }
  }, [])

  // Get top collections
  const getTopCollections = useCallback((limit: number = 10): NFTCollection[] => {
    return nftService.getTopCollections(limit)
  }, [])

  // Get trending collections
  const getTrendingCollections = useCallback((limit: number = 10): NFTCollection[] => {
    return nftService.getTrendingCollections(limit)
  }, [])

  // Search collections
  const searchCollections = useCallback((query: string): NFTCollection[] => {
    return nftService.searchCollections(query)
  }, [])

  // Get collections by category
  const getCollectionsByCategory = useCallback((category: string): NFTCollection[] => {
    return nftService.getCollectionsByCategory(category)
  }, [])

  // Calculate portfolio metrics
  const getPortfolioMetrics = useCallback(() => {
    if (!userAddress) {
      return {
        totalValue: 0,
        totalNFTs: 0,
        totalCollections: 0,
        topNFT: null,
        averageValue: 0,
        totalGainLoss: 0,
        gainLossPercentage: 0
      }
    }

    const portfolioValue = nftService.calculatePortfolioValue(userAddress)
    const portfolio = state.userPortfolio

    return {
      totalValue: portfolioValue.totalValue,
      totalNFTs: portfolio?.totalNFTs || 0,
      totalCollections: portfolio?.totalCollections || 0,
      topNFT: portfolioValue.topNFT,
      averageValue: portfolioValue.averageValue,
      totalGainLoss: portfolioValue.totalGainLoss,
      gainLossPercentage: portfolioValue.totalValue > 0 
        ? (portfolioValue.totalGainLoss / portfolioValue.totalValue) * 100 
        : 0
    }
  }, [userAddress, state.userPortfolio])

  // Get NFTs by rarity
  const getNFTsByRarity = useCallback((rarity: NFTRarity): NFT[] => {
    return state.userNFTs.filter(nft => nft.rarity === rarity)
  }, [state.userNFTs])

  // Get listed NFTs
  const getListedNFTs = useCallback(() => {
    return state.userNFTs.filter(nft => nft.isListed)
  }, [state.userNFTs])

  // Get NFTs by collection
  const getNFTsByCollection = useCallback((collectionId: string) => {
    const collection = getCollection(collectionId)
    if (!collection) return []

    return state.userNFTs.filter(nft =>
      nft.contractAddress.toLowerCase() === collection.contractAddress.toLowerCase() &&
      nft.chainId === collection.chainId
    )
  }, [state.userNFTs, getCollection])

  // Get collection distribution
  const getCollectionDistribution = useCallback(() => {
    const distribution = new Map<string, number>()
    
    state.userNFTs.forEach(nft => {
      const collection = getCollectionByAddress(nft.contractAddress, nft.chainId)
      if (collection) {
        const current = distribution.get(collection.name) || 0
        distribution.set(collection.name, current + 1)
      }
    })

    return Array.from(distribution.entries())
      .map(([name, count]) => ({ name, count, percentage: (count / state.userNFTs.length) * 100 }))
      .sort((a, b) => b.count - a.count)
  }, [state.userNFTs, getCollectionByAddress])

  // Get rarity distribution
  const getRarityDistribution = useCallback(() => {
    const distribution = new Map<NFTRarity, number>()
    
    state.userNFTs.forEach(nft => {
      if (nft.rarity) {
        const current = distribution.get(nft.rarity) || 0
        distribution.set(nft.rarity, current + 1)
      }
    })

    return Array.from(distribution.entries())
      .map(([rarity, count]) => ({ rarity, count, percentage: (count / state.userNFTs.length) * 100 }))
      .sort((a, b) => b.count - a.count)
  }, [state.userNFTs])

  // Get recent activity
  const getRecentActivity = useCallback((limit: number = 10) => {
    return state.userPortfolio?.recentActivity.slice(0, limit) || []
  }, [state.userPortfolio])

  // Format currency
  const formatCurrency = useCallback((amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
  }, [])

  // Format ETH
  const formatETH = useCallback((amount: string | number) => {
    const value = typeof amount === 'string' ? parseFloat(amount) : amount
    return `${value.toFixed(3)} ETH`
  }, [])

  // Get collection categories
  const getCollectionCategories = useCallback(() => {
    const categories = new Set<string>()
    state.collections.forEach(collection => {
      categories.add(collection.category)
    })
    return Array.from(categories).sort()
  }, [state.collections])

  // Get featured collections
  const getFeaturedCollections = useCallback(() => {
    return state.collections.filter(collection => collection.featured)
  }, [state.collections])

  // Get verified collections
  const getVerifiedCollections = useCallback(() => {
    return state.collections.filter(collection => collection.verified)
  }, [state.collections])

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

    // Getters
    getCollection,
    getCollectionByAddress,
    getNFT,
    getNFTValuation,
    getTopCollections,
    getTrendingCollections,
    searchCollections,
    getCollectionsByCategory,
    getNFTsByRarity,
    getListedNFTs,
    getNFTsByCollection,

    // Analytics
    getPortfolioMetrics,
    getCollectionDistribution,
    getRarityDistribution,
    getRecentActivity,
    getCollectionCategories,

    // Utilities
    formatCurrency,
    formatETH,

    // Computed values
    portfolioMetrics: getPortfolioMetrics(),
    collectionDistribution: getCollectionDistribution(),
    rarityDistribution: getRarityDistribution(),
    recentActivity: getRecentActivity(),
    topCollections: getTopCollections(5),
    trendingCollections: getTrendingCollections(5),
    featuredCollections: getFeaturedCollections(),
    verifiedCollections: getVerifiedCollections(),
    collectionCategories: getCollectionCategories(),

    // Quick access to rarity groups
    commonNFTs: getNFTsByRarity(NFTRarity.COMMON),
    uncommonNFTs: getNFTsByRarity(NFTRarity.UNCOMMON),
    rareNFTs: getNFTsByRarity(NFTRarity.RARE),
    epicNFTs: getNFTsByRarity(NFTRarity.EPIC),
    legendaryNFTs: getNFTsByRarity(NFTRarity.LEGENDARY),

    // Quick access to listing status
    listedNFTs: getListedNFTs(),
    unlistedNFTs: state.userNFTs.filter(nft => !nft.isListed)
  }
}
