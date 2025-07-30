import { useState, useEffect, useCallback } from 'react'
import { useAccount } from 'wagmi'
import { type Address } from 'viem'
import { 
  nftMarketplaceIntegration,
  type NFTMarketplace,
  type NFTListing,
  type NFTOffer,
  type NFTTransaction,
  type NFTPortfolio,
  type MarketplaceFilters,
  ListingType,
  OfferType,
  type MarketplaceEvent
} from '@/lib/nft-marketplace'
import { toast } from 'sonner'

export interface NFTMarketplaceState {
  marketplaces: NFTMarketplace[]
  listings: NFTListing[]
  offers: NFTOffer[]
  transactions: NFTTransaction[]
  portfolio: NFTPortfolio | null
  isLoading: boolean
  isExecuting: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseNFTMarketplaceOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseNFTMarketplaceReturn {
  // State
  state: NFTMarketplaceState
  
  // Listing Operations
  getListings: (filters?: MarketplaceFilters) => Promise<NFTListing[]>
  createListing: (
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    listingType?: ListingType,
    endTime?: string
  ) => Promise<NFTListing>
  
  // Trading Operations
  purchaseNFT: (listingId: string, price?: string) => Promise<NFTTransaction>
  createOffer: (
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    offerType?: OfferType,
    expiration?: string
  ) => Promise<NFTOffer>
  
  // Portfolio Operations
  getUserPortfolio: (userAddress?: Address) => Promise<NFTPortfolio>
  refreshPortfolio: () => Promise<void>
  
  // Data Access
  getMarketplaces: () => NFTMarketplace[]
  getMarketplace: (id: string) => NFTMarketplace | null
  getTransaction: (id: string) => NFTTransaction | null
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useNFTMarketplace = (
  options: UseNFTMarketplaceOptions = {}
): UseNFTMarketplaceReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefresh = true,
    refreshInterval = 30000
  } = options

  const { address } = useAccount()

  const [state, setState] = useState<NFTMarketplaceState>({
    marketplaces: [],
    listings: [],
    offers: [],
    transactions: [],
    portfolio: null,
    isLoading: false,
    isExecuting: false,
    error: null,
    lastUpdate: null
  })

  // Update state from marketplace integration
  const updateState = useCallback(async () => {
    try {
      const marketplaces = nftMarketplaceIntegration.getMarketplaces()

      setState(prev => ({
        ...prev,
        marketplaces,
        error: null,
        lastUpdate: Date.now()
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [])

  // Handle marketplace events
  const handleMarketplaceEvent = useCallback((event: MarketplaceEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'listing_created':
          toast.success('Listing Created', {
            description: `Successfully listed NFT for ${event.listing?.price} ${event.listing?.currency}`
          })
          break
        case 'purchase_completed':
          toast.success('Purchase Completed', {
            description: `Successfully purchased NFT for ${event.transaction?.price} ${event.transaction?.currency}`
          })
          break
        case 'purchase_failed':
          toast.error('Purchase Failed', {
            description: `Purchase failed: ${event.error?.message || 'Unknown error'}`
          })
          break
        case 'offer_created':
          toast.success('Offer Created', {
            description: `Offer submitted for ${event.offer?.price} ${event.offer?.currency}`
          })
          break
        case 'offer_accepted':
          toast.success('Offer Accepted', {
            description: `Your offer has been accepted!`
          })
          break
        case 'portfolio_updated':
          toast.info('Portfolio Updated', {
            description: 'Your NFT portfolio has been updated'
          })
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = nftMarketplaceIntegration.addEventListener(handleMarketplaceEvent)

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, handleMarketplaceEvent, updateState])

  // Auto-refresh data
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Get listings
  const getListings = useCallback(async (
    filters?: MarketplaceFilters
  ): Promise<NFTListing[]> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const listings = await nftMarketplaceIntegration.getListings(filters)
      
      setState(prev => ({
        ...prev,
        isLoading: false,
        listings
      }))

      return listings
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get listings', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Create listing
  const createListing = useCallback(async (
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    listingType: ListingType = ListingType.FIXED_PRICE,
    endTime?: string
  ): Promise<NFTListing> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const listing = await nftMarketplaceIntegration.createListing(
        tokenId,
        contractAddress,
        price,
        currency,
        marketplace,
        address,
        listingType,
        endTime
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        listings: [...prev.listings, listing]
      }))

      return listing
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Purchase NFT
  const purchaseNFT = useCallback(async (
    listingId: string,
    price?: string
  ): Promise<NFTTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await nftMarketplaceIntegration.purchaseNFT(
        listingId,
        address,
        price
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Create offer
  const createOffer = useCallback(async (
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    offerType: OfferType = OfferType.TOKEN,
    expiration?: string
  ): Promise<NFTOffer> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const offer = await nftMarketplaceIntegration.createOffer(
        tokenId,
        contractAddress,
        price,
        currency,
        marketplace,
        address,
        offerType,
        expiration
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        offers: [...prev.offers, offer]
      }))

      return offer
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Get user portfolio
  const getUserPortfolio = useCallback(async (
    userAddress?: Address
  ): Promise<NFTPortfolio> => {
    const targetAddress = userAddress || address
    if (!targetAddress) {
      throw new Error('User address not available')
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const portfolio = await nftMarketplaceIntegration.getUserPortfolio(targetAddress)
      
      setState(prev => ({
        ...prev,
        isLoading: false,
        portfolio
      }))

      return portfolio
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get portfolio', { description: errorMessage })
      }
      throw error
    }
  }, [address, enableNotifications])

  // Refresh portfolio
  const refreshPortfolio = useCallback(async () => {
    if (address) {
      await getUserPortfolio(address)
    }
  }, [address, getUserPortfolio])

  // Get marketplaces
  const getMarketplaces = useCallback((): NFTMarketplace[] => {
    return nftMarketplaceIntegration.getMarketplaces()
  }, [])

  // Get marketplace by ID
  const getMarketplace = useCallback((id: string): NFTMarketplace | null => {
    return nftMarketplaceIntegration.getMarketplace(id)
  }, [])

  // Get transaction
  const getTransaction = useCallback((id: string): NFTTransaction | null => {
    return nftMarketplaceIntegration.getTransaction(id)
  }, [])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
    if (address) {
      refreshPortfolio()
    }
  }, [updateState, address, refreshPortfolio])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    getListings,
    createListing,
    purchaseNFT,
    createOffer,
    getUserPortfolio,
    refreshPortfolio,
    getMarketplaces,
    getMarketplace,
    getTransaction,
    refresh,
    clearError
  }
}

// Simplified hook for NFT trading
export const useNFTTrading = () => {
  const { state, purchaseNFT, createOffer } = useNFTMarketplace()

  const trade = useCallback(async (
    listingId: string,
    type: 'buy' | 'offer',
    price?: string
  ) => {
    if (type === 'buy') {
      return purchaseNFT(listingId, price)
    } else {
      // For offers, we'd need additional parameters
      throw new Error('Offer creation requires more parameters')
    }
  }, [purchaseNFT])

  return {
    trade,
    isExecuting: state.isExecuting,
    transactions: state.transactions,
    error: state.error
  }
}

// Hook for NFT portfolio management
export const useNFTPortfolio = () => {
  const { state, getUserPortfolio, refreshPortfolio } = useNFTMarketplace()

  const analytics = state.portfolio ? {
    totalValue: state.portfolio.totalValue,
    totalItems: state.portfolio.totalItems,
    totalCollections: state.portfolio.collections.length,
    totalUnrealizedPnL: state.portfolio.collections.reduce((sum, col) => sum + col.unrealizedPnL, 0),
    averageUnrealizedPnL: state.portfolio.collections.length > 0
      ? state.portfolio.collections.reduce((sum, col) => sum + col.unrealizedPnLPercentage, 0) / state.portfolio.collections.length
      : 0,
    topCollection: state.portfolio.collections.sort((a, b) => b.totalValue - a.totalValue)[0],
    floorValueChange24h: state.portfolio.floorValueChange24h
  } : null

  return {
    portfolio: state.portfolio,
    analytics,
    getUserPortfolio,
    refreshPortfolio,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for marketplace analytics
export const useMarketplaceAnalytics = () => {
  const { state } = useNFTMarketplace()

  const analytics = {
    totalMarketplaces: state.marketplaces.length,
    totalListings: state.listings.length,
    totalOffers: state.offers.length,
    totalTransactions: state.transactions.length,
    totalVolume: state.transactions.reduce((sum, tx) => {
      if (tx.status === 'confirmed') {
        return sum + tx.priceUSD
      }
      return sum
    }, 0),
    averagePrice: state.transactions.length > 0
      ? state.transactions.reduce((sum, tx) => sum + tx.priceUSD, 0) / state.transactions.length
      : 0,
    marketplaceDistribution: state.transactions.reduce((acc, tx) => {
      acc[tx.marketplace] = (acc[tx.marketplace] || 0) + 1
      return acc
    }, {} as Record<string, number>),
    successRate: state.transactions.length > 0
      ? (state.transactions.filter(tx => tx.status === 'confirmed').length / state.transactions.length) * 100
      : 0
  }

  return analytics
}
