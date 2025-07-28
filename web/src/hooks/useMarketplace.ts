import { useState, useEffect, useCallback } from 'react'
import { type Address, type Hash } from 'viem'
import { 
  marketplaceService, 
  type MarketplaceConfig,
  type NFTListing,
  type NFTOffer,
  type MarketplaceActivity,
  type MarketplaceStats,
  MarketplaceName,
  OrderType,
  OrderStatus 
} from '@/lib/marketplace-service'
import { toast } from 'sonner'

export interface UseMarketplaceOptions {
  userAddress?: Address
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface MarketplaceState {
  marketplaces: MarketplaceConfig[]
  listings: NFTListing[]
  offers: NFTOffer[]
  recentActivity: MarketplaceActivity[]
  marketplaceStats: MarketplaceStats[]
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useMarketplace(options: UseMarketplaceOptions = {}) {
  const {
    userAddress,
    autoRefresh = true,
    refreshInterval = 60000, // 1 minute
    enableNotifications = true
  } = options

  const [state, setState] = useState<MarketplaceState>({
    marketplaces: [],
    listings: [],
    offers: [],
    recentActivity: [],
    marketplaceStats: [],
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load marketplace data
  const loadData = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const marketplaces = marketplaceService.getAllMarketplaces()
      const listings = marketplaceService.getListings()
      const offers = marketplaceService.getOffers()
      const recentActivity = marketplaceService.getRecentActivity()
      const marketplaceStats = marketplaceService.getMarketplaceStats()

      setState(prev => ({
        ...prev,
        marketplaces,
        listings,
        offers,
        recentActivity,
        marketplaceStats,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load marketplace data'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
    }
  }, [])

  // Get marketplace by name
  const getMarketplace = useCallback((name: MarketplaceName) => {
    return marketplaceService.getMarketplace(name)
  }, [])

  // Get listings for specific NFT
  const getListingsForNFT = useCallback((contractAddress: Address, tokenId: string) => {
    return marketplaceService.getListings(contractAddress, tokenId)
  }, [])

  // Get offers for specific NFT
  const getOffersForNFT = useCallback((contractAddress: Address, tokenId: string) => {
    return marketplaceService.getOffers(contractAddress, tokenId)
  }, [])

  // Get listings by marketplace
  const getListingsByMarketplace = useCallback((marketplace: MarketplaceName) => {
    return marketplaceService.getListings(undefined, undefined, marketplace)
  }, [])

  // Create listing
  const createListing = useCallback(async (
    contractAddress: Address,
    tokenId: string,
    price: string,
    marketplace: MarketplaceName,
    duration?: number
  ): Promise<Hash> => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await marketplaceService.createListing(
        contractAddress,
        tokenId,
        price,
        marketplace,
        duration
      )

      // Refresh data after creating listing
      loadData()

      if (enableNotifications) {
        const marketplaceConfig = getMarketplace(marketplace)
        toast.success('Listing created successfully!', {
          description: `Listed NFT #${tokenId} for ${price} ETH on ${marketplaceConfig?.name}`,
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create listing'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Listing failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [loadData, enableNotifications, getMarketplace])

  // Create offer
  const createOffer = useCallback(async (
    contractAddress: Address,
    tokenId: string,
    price: string,
    marketplace: MarketplaceName,
    duration?: number
  ): Promise<Hash> => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await marketplaceService.createOffer(
        contractAddress,
        tokenId,
        price,
        marketplace,
        duration
      )

      // Refresh data after creating offer
      loadData()

      if (enableNotifications) {
        const marketplaceConfig = getMarketplace(marketplace)
        toast.success('Offer created successfully!', {
          description: `Offered ${price} ETH for NFT #${tokenId} on ${marketplaceConfig?.name}`,
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create offer'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Offer failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [loadData, enableNotifications, getMarketplace])

  // Cancel listing
  const cancelListing = useCallback(async (listingId: string): Promise<Hash> => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await marketplaceService.cancelListing(listingId)

      // Refresh data after canceling listing
      loadData()

      if (enableNotifications) {
        toast.success('Listing cancelled successfully!', {
          description: 'Your NFT listing has been cancelled',
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to cancel listing'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Cancellation failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [loadData, enableNotifications])

  // Cancel offer
  const cancelOffer = useCallback(async (offerId: string): Promise<Hash> => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await marketplaceService.cancelOffer(offerId)

      // Refresh data after canceling offer
      loadData()

      if (enableNotifications) {
        toast.success('Offer cancelled successfully!', {
          description: 'Your NFT offer has been cancelled',
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to cancel offer'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Cancellation failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [loadData, enableNotifications])

  // Get marketplace analytics
  const getMarketplaceAnalytics = useCallback(() => {
    const totalVolume24h = state.marketplaceStats.reduce((sum, stats) => {
      return sum + parseFloat(stats.volume24h)
    }, 0)

    const totalSales24h = state.marketplaceStats.reduce((sum, stats) => {
      return sum + stats.sales24h
    }, 0)

    const averagePrice = totalSales24h > 0 ? totalVolume24h / totalSales24h : 0

    const topMarketplace = state.marketplaceStats.reduce((top, current) => {
      return parseFloat(current.volume24h) > parseFloat(top.volume24h) ? current : top
    }, state.marketplaceStats[0])

    return {
      totalVolume24h,
      totalSales24h,
      averagePrice,
      topMarketplace: topMarketplace?.marketplace,
      activeMarketplaces: state.marketplaces.length,
      totalListings: state.listings.length,
      totalOffers: state.offers.length
    }
  }, [state.marketplaceStats, state.marketplaces, state.listings, state.offers])

  // Get user's listings
  const getUserListings = useCallback(() => {
    if (!userAddress) return []
    return state.listings.filter(listing => 
      listing.seller.toLowerCase() === userAddress.toLowerCase() &&
      listing.status === OrderStatus.ACTIVE
    )
  }, [state.listings, userAddress])

  // Get user's offers
  const getUserOffers = useCallback(() => {
    if (!userAddress) return []
    return state.offers.filter(offer => 
      offer.buyer.toLowerCase() === userAddress.toLowerCase() &&
      offer.status === OrderStatus.ACTIVE
    )
  }, [state.offers, userAddress])

  // Get marketplace distribution
  const getMarketplaceDistribution = useCallback(() => {
    const distribution = new Map<MarketplaceName, number>()
    
    state.listings.forEach(listing => {
      const current = distribution.get(listing.marketplace) || 0
      distribution.set(listing.marketplace, current + 1)
    })

    return Array.from(distribution.entries())
      .map(([marketplace, count]) => ({ 
        marketplace, 
        count, 
        percentage: state.listings.length > 0 ? (count / state.listings.length) * 100 : 0 
      }))
      .sort((a, b) => b.count - a.count)
  }, [state.listings])

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

  // Get marketplace icon
  const getMarketplaceIcon = useCallback((marketplace: MarketplaceName) => {
    const icons = {
      [MarketplaceName.OPENSEA]: '🌊',
      [MarketplaceName.BLUR]: '💨',
      [MarketplaceName.LOOKSRARE]: '👀',
      [MarketplaceName.X2Y2]: '❌',
      [MarketplaceName.FOUNDATION]: '🏛️',
      [MarketplaceName.SUPERRARE]: '💎'
    }
    return icons[marketplace] || '🏪'
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
    createListing,
    createOffer,
    cancelListing,
    cancelOffer,

    // Getters
    getMarketplace,
    getListingsForNFT,
    getOffersForNFT,
    getListingsByMarketplace,
    getUserListings,
    getUserOffers,

    // Analytics
    getMarketplaceAnalytics,
    getMarketplaceDistribution,

    // Utilities
    formatCurrency,
    formatETH,
    getMarketplaceIcon,

    // Computed values
    marketplaceAnalytics: getMarketplaceAnalytics(),
    marketplaceDistribution: getMarketplaceDistribution(),
    userListings: getUserListings(),
    userOffers: getUserOffers(),

    // Quick access to marketplace types
    openSeaListings: getListingsByMarketplace(MarketplaceName.OPENSEA),
    blurListings: getListingsByMarketplace(MarketplaceName.BLUR),
    looksRareListings: getListingsByMarketplace(MarketplaceName.LOOKSRARE),
    x2y2Listings: getListingsByMarketplace(MarketplaceName.X2Y2),

    // Quick access to order types
    activeListings: state.listings.filter(l => l.status === OrderStatus.ACTIVE),
    activeOffers: state.offers.filter(o => o.status === OrderStatus.ACTIVE),
    expiredOffers: state.offers.filter(o => o.expirationTime < Date.now())
  }
}
