import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  nftDiscoveryEngine,
  type NFTCollection,
  type NFTToken,
  type NFTCategory,
  type CollectionSearchFilters,
  type NFTSearchFilters,
  type SearchResult,
  type NFTDiscoveryEvent
} from '@/lib/nft-discovery'
import { toast } from 'sonner'

export interface NFTDiscoveryState {
  collections: NFTCollection[]
  tokens: NFTToken[]
  trendingCollections: NFTCollection[]
  featuredCollections: NFTCollection[]
  categories: { value: NFTCategory; label: string; count: number }[]
  isLoading: boolean
  isSearching: boolean
  error: string | null
  searchQuery: string
  lastUpdate: number | null
  hasMore: boolean
  total: number
}

export interface UseNFTDiscoveryOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseNFTDiscoveryReturn {
  // State
  state: NFTDiscoveryState
  
  // Collection Operations
  searchCollections: (query?: string, filters?: CollectionSearchFilters) => Promise<SearchResult<NFTCollection>>
  getTrendingCollections: (limit?: number) => Promise<NFTCollection[]>
  getFeaturedCollections: (limit?: number) => Promise<NFTCollection[]>
  getCollection: (id: string) => Promise<NFTCollection | null>
  getCollectionByContract: (contractAddress: Address) => Promise<NFTCollection | null>
  
  // Token Operations
  searchTokens: (query?: string, filters?: NFTSearchFilters) => Promise<SearchResult<NFTToken>>
  
  // Data Access
  getCategories: () => { value: NFTCategory; label: string; count: number }[]
  
  // Search Management
  setSearchQuery: (query: string) => void
  clearSearch: () => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useNFTDiscovery = (
  options: UseNFTDiscoveryOptions = {}
): UseNFTDiscoveryReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefresh = false,
    refreshInterval = 60000
  } = options

  const [state, setState] = useState<NFTDiscoveryState>({
    collections: [],
    tokens: [],
    trendingCollections: [],
    featuredCollections: [],
    categories: [],
    isLoading: false,
    isSearching: false,
    error: null,
    searchQuery: '',
    lastUpdate: null,
    hasMore: false,
    total: 0
  })

  // Update state from NFT discovery engine
  const updateState = useCallback(async () => {
    try {
      const categories = nftDiscoveryEngine.getCategories()

      setState(prev => ({
        ...prev,
        categories,
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

  // Handle NFT discovery events
  const handleNFTDiscoveryEvent = useCallback((event: NFTDiscoveryEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'collection_updated':
          toast.info('Collection Updated', {
            description: `${event.collection?.name} has been updated`
          })
          break
        case 'trending_updated':
          toast.info('Trending Updated', {
            description: 'Trending collections have been updated'
          })
          break
        case 'featured_updated':
          toast.info('Featured Updated', {
            description: 'Featured collections have been updated'
          })
          break
        case 'search_completed':
          if (event.searchQuery) {
            toast.success('Search Completed', {
              description: `Found results for "${event.searchQuery}"`
            })
          }
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = nftDiscoveryEngine.addEventListener(handleNFTDiscoveryEvent)

    // Initial state update
    if (autoLoad) {
      updateState()
      loadInitialData()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, handleNFTDiscoveryEvent, updateState])

  // Auto-refresh data
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Load initial data
  const loadInitialData = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const [trending, featured] = await Promise.all([
        nftDiscoveryEngine.getTrendingCollections(10),
        nftDiscoveryEngine.getFeaturedCollections(10)
      ])

      setState(prev => ({
        ...prev,
        isLoading: false,
        trendingCollections: trending,
        featuredCollections: featured
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to load initial data', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Search collections
  const searchCollections = useCallback(async (
    query?: string,
    filters?: CollectionSearchFilters
  ): Promise<SearchResult<NFTCollection>> => {
    setState(prev => ({ ...prev, isSearching: true, error: null }))

    try {
      const result = await nftDiscoveryEngine.searchCollections(query, filters)
      
      setState(prev => ({
        ...prev,
        isSearching: false,
        collections: result.items,
        hasMore: result.hasMore,
        total: result.total,
        searchQuery: query || ''
      }))

      if (enableNotifications && query) {
        toast.success('Search Completed', {
          description: `Found ${result.total} collections for "${query}"`
        })
      }

      return result
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSearching: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Search failed', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Get trending collections
  const getTrendingCollections = useCallback(async (
    limit: number = 10
  ): Promise<NFTCollection[]> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const trending = await nftDiscoveryEngine.getTrendingCollections(limit)
      
      setState(prev => ({
        ...prev,
        isLoading: false,
        trendingCollections: trending
      }))

      return trending
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get trending collections', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Get featured collections
  const getFeaturedCollections = useCallback(async (
    limit: number = 10
  ): Promise<NFTCollection[]> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const featured = await nftDiscoveryEngine.getFeaturedCollections(limit)
      
      setState(prev => ({
        ...prev,
        isLoading: false,
        featuredCollections: featured
      }))

      return featured
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get featured collections', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Get collection by ID
  const getCollection = useCallback(async (
    id: string
  ): Promise<NFTCollection | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const collection = await nftDiscoveryEngine.getCollection(id)
      
      setState(prev => ({
        ...prev,
        isLoading: false
      }))

      return collection
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get collection', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Get collection by contract address
  const getCollectionByContract = useCallback(async (
    contractAddress: Address
  ): Promise<NFTCollection | null> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const collection = await nftDiscoveryEngine.getCollectionByContract(contractAddress)
      
      setState(prev => ({
        ...prev,
        isLoading: false
      }))

      return collection
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get collection', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Search tokens
  const searchTokens = useCallback(async (
    query?: string,
    filters?: NFTSearchFilters
  ): Promise<SearchResult<NFTToken>> => {
    setState(prev => ({ ...prev, isSearching: true, error: null }))

    try {
      const result = await nftDiscoveryEngine.searchTokens(query, filters)
      
      setState(prev => ({
        ...prev,
        isSearching: false,
        tokens: result.items,
        hasMore: result.hasMore,
        total: result.total,
        searchQuery: query || ''
      }))

      return result
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSearching: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Token search failed', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Get categories
  const getCategories = useCallback(() => {
    return nftDiscoveryEngine.getCategories()
  }, [])

  // Set search query
  const setSearchQuery = useCallback((query: string) => {
    setState(prev => ({ ...prev, searchQuery: query }))
  }, [])

  // Clear search
  const clearSearch = useCallback(() => {
    setState(prev => ({
      ...prev,
      searchQuery: '',
      collections: [],
      tokens: [],
      hasMore: false,
      total: 0
    }))
  }, [])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
    loadInitialData()
  }, [updateState, loadInitialData])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    searchCollections,
    getTrendingCollections,
    getFeaturedCollections,
    getCollection,
    getCollectionByContract,
    searchTokens,
    getCategories,
    setSearchQuery,
    clearSearch,
    refresh,
    clearError
  }
}

// Simplified hook for collection browsing
export const useCollectionBrowser = () => {
  const { state, searchCollections, getTrendingCollections, getFeaturedCollections } = useNFTDiscovery()

  const browse = useCallback(async (filters?: CollectionSearchFilters) => {
    return searchCollections(undefined, filters)
  }, [searchCollections])

  return {
    collections: state.collections,
    trendingCollections: state.trendingCollections,
    featuredCollections: state.featuredCollections,
    browse,
    getTrending: getTrendingCollections,
    getFeatured: getFeaturedCollections,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for NFT search
export const useNFTSearch = () => {
  const { state, searchCollections, searchTokens, setSearchQuery, clearSearch } = useNFTDiscovery()

  const search = useCallback(async (
    query: string,
    type: 'collections' | 'tokens' = 'collections',
    filters?: CollectionSearchFilters | NFTSearchFilters
  ) => {
    setSearchQuery(query)
    
    if (type === 'collections') {
      return searchCollections(query, filters as CollectionSearchFilters)
    } else {
      return searchTokens(query, filters as NFTSearchFilters)
    }
  }, [searchCollections, searchTokens, setSearchQuery])

  return {
    search,
    clearSearch,
    searchQuery: state.searchQuery,
    collections: state.collections,
    tokens: state.tokens,
    isSearching: state.isSearching,
    hasMore: state.hasMore,
    total: state.total,
    error: state.error
  }
}

// Hook for NFT categories
export const useNFTCategories = () => {
  const { getCategories } = useNFTDiscovery()

  const categories = getCategories()

  return {
    categories,
    getCategoryCollections: (category: NFTCategory) => {
      // This would filter collections by category
      return []
    }
  }
}
