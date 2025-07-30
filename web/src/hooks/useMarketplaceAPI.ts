import { useState, useEffect, useCallback, useRef } from 'react'
import { 
  marketplaceAPIManager,
  type NFTCollection,
  type NFTAsset,
  type MarketplaceEvent,
  type MarketplaceStats,
  type SearchFilters,
  type MarketplaceAPI,
  type MarketplaceAPIEvent,
  type APIResponse
} from '@/lib/marketplace-api'
import { toast } from 'sonner'

export interface MarketplaceAPIState {
  collections: NFTCollection[]
  assets: NFTAsset[]
  events: MarketplaceEvent[]
  stats: Record<string, MarketplaceStats>
  marketplaces: MarketplaceAPI[]
  isLoading: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseMarketplaceAPIOptions {
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
  defaultMarketplaces?: string[]
}

export interface UseMarketplaceAPIReturn {
  // State
  state: MarketplaceAPIState
  
  // Collection Operations
  getCollections: (marketplace: string, filters?: SearchFilters) => Promise<APIResponse<NFTCollection[]>>
  aggregateCollections: (marketplaces: string[], filters?: SearchFilters) => Promise<any>
  
  // Asset Operations
  getAssets: (marketplace: string, filters?: SearchFilters) => Promise<APIResponse<NFTAsset[]>>
  aggregateAssets: (marketplaces: string[], filters?: SearchFilters) => Promise<any>
  
  // Event Operations
  getEvents: (marketplace: string, filters?: SearchFilters) => Promise<APIResponse<MarketplaceEvent[]>>
  
  // Stats Operations
  getStats: (marketplace: string) => Promise<APIResponse<MarketplaceStats>>
  getAllStats: (marketplaces: string[]) => Promise<Record<string, MarketplaceStats>>
  
  // Search Operations
  searchCollections: (query: string, marketplaces?: string[]) => Promise<NFTCollection[]>
  searchAssets: (query: string, marketplaces?: string[]) => Promise<NFTAsset[]>
  
  // Utilities
  getMarketplaces: () => MarketplaceAPI[]
  clearCache: () => void
  refresh: () => void
  clearError: () => void
}

export const useMarketplaceAPI = (
  options: UseMarketplaceAPIOptions = {}
): UseMarketplaceAPIReturn => {
  const {
    enableNotifications = true,
    autoRefresh = false,
    refreshInterval = 300000, // 5 minutes
    defaultMarketplaces = ['opensea', 'looksrare', 'x2y2', 'blur']
  } = options

  const [state, setState] = useState<MarketplaceAPIState>({
    collections: [],
    assets: [],
    events: [],
    stats: {},
    marketplaces: [],
    isLoading: false,
    error: null,
    lastUpdate: null
  })

  const refreshTimeoutRef = useRef<NodeJS.Timeout>()

  // Handle marketplace API events
  const handleMarketplaceAPIEvent = useCallback((event: MarketplaceAPIEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'api_error':
          toast.error('Marketplace API Error', {
            description: `${event.marketplace}: ${event.error?.message || 'Unknown error'}`
          })
          break
        case 'rate_limit_hit':
          toast.warning('Rate Limit Hit', {
            description: `${event.marketplace} API rate limit reached. Please wait.`
          })
          break
        case 'api_success':
          // Only show success for important operations
          if (event.endpoint?.includes('collections') || event.endpoint?.includes('stats')) {
            toast.success('Data Updated', {
              description: `${event.marketplace} data refreshed successfully`
            })
          }
          break
      }
    }

    // Update state based on event
    setState(prev => ({
      ...prev,
      error: event.type === 'api_error' ? event.error?.message || 'API Error' : null,
      lastUpdate: Date.now()
    }))
  }, [enableNotifications])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = marketplaceAPIManager.addEventListener(handleMarketplaceAPIEvent)

    // Load initial data
    const loadInitialData = async () => {
      setState(prev => ({ ...prev, isLoading: true }))
      
      try {
        const marketplaces = marketplaceAPIManager.getMarketplaces()
        setState(prev => ({ 
          ...prev, 
          marketplaces,
          isLoading: false,
          lastUpdate: Date.now()
        }))
      } catch (error) {
        setState(prev => ({ 
          ...prev, 
          error: (error as Error).message,
          isLoading: false
        }))
      }
    }

    loadInitialData()

    return () => {
      unsubscribe()
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current)
      }
    }
  }, [handleMarketplaceAPIEvent])

  // Auto-refresh
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const scheduleRefresh = () => {
        refreshTimeoutRef.current = setTimeout(() => {
          refresh()
          scheduleRefresh()
        }, refreshInterval)
      }

      scheduleRefresh()

      return () => {
        if (refreshTimeoutRef.current) {
          clearTimeout(refreshTimeoutRef.current)
        }
      }
    }
  }, [autoRefresh, refreshInterval])

  // Get collections
  const getCollections = useCallback(async (
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<NFTCollection[]>> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await marketplaceAPIManager.getCollections(marketplace, filters)
      
      if (response.success) {
        setState(prev => ({
          ...prev,
          collections: response.data,
          isLoading: false,
          lastUpdate: Date.now()
        }))
      } else {
        setState(prev => ({
          ...prev,
          error: response.error || 'Failed to fetch collections',
          isLoading: false
        }))
      }

      return response
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Aggregate collections from multiple marketplaces
  const aggregateCollections = useCallback(async (
    marketplaces: string[],
    filters?: SearchFilters
  ) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await marketplaceAPIManager.aggregateCollections(marketplaces, filters)
      
      setState(prev => ({
        ...prev,
        collections: response.data,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      return response
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Get assets
  const getAssets = useCallback(async (
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<NFTAsset[]>> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await marketplaceAPIManager.getAssets(marketplace, filters)
      
      if (response.success) {
        setState(prev => ({
          ...prev,
          assets: response.data,
          isLoading: false,
          lastUpdate: Date.now()
        }))
      } else {
        setState(prev => ({
          ...prev,
          error: response.error || 'Failed to fetch assets',
          isLoading: false
        }))
      }

      return response
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Aggregate assets from multiple marketplaces
  const aggregateAssets = useCallback(async (
    marketplaces: string[],
    filters?: SearchFilters
  ) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      // Aggregate assets from multiple marketplaces
      const responses = await Promise.allSettled(
        marketplaces.map(marketplace => marketplaceAPIManager.getAssets(marketplace, filters))
      )

      const successful = responses
        .filter((result): result is PromiseFulfilledResult<APIResponse<NFTAsset[]>> => 
          result.status === 'fulfilled' && result.value.success
        )
        .map(result => result.value)

      const allAssets = successful.flatMap(response => response.data)
      
      setState(prev => ({
        ...prev,
        assets: allAssets,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      return {
        data: allAssets,
        sources: successful.map(r => r.marketplace),
        totalSources: marketplaces.length,
        successfulSources: successful.length
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Get events
  const getEvents = useCallback(async (
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<MarketplaceEvent[]>> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await marketplaceAPIManager.getEvents(marketplace, filters)
      
      if (response.success) {
        setState(prev => ({
          ...prev,
          events: response.data,
          isLoading: false,
          lastUpdate: Date.now()
        }))
      } else {
        setState(prev => ({
          ...prev,
          error: response.error || 'Failed to fetch events',
          isLoading: false
        }))
      }

      return response
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Get stats
  const getStats = useCallback(async (
    marketplace: string
  ): Promise<APIResponse<MarketplaceStats>> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await marketplaceAPIManager.getStats(marketplace)
      
      if (response.success) {
        setState(prev => ({
          ...prev,
          stats: { ...prev.stats, [marketplace]: response.data },
          isLoading: false,
          lastUpdate: Date.now()
        }))
      } else {
        setState(prev => ({
          ...prev,
          error: response.error || 'Failed to fetch stats',
          isLoading: false
        }))
      }

      return response
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Get all stats
  const getAllStats = useCallback(async (
    marketplaces: string[]
  ): Promise<Record<string, MarketplaceStats>> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const responses = await Promise.allSettled(
        marketplaces.map(marketplace => marketplaceAPIManager.getStats(marketplace))
      )

      const stats: Record<string, MarketplaceStats> = {}
      
      responses.forEach((result, index) => {
        if (result.status === 'fulfilled' && result.value.success) {
          stats[marketplaces[index]] = result.value.data
        }
      })

      setState(prev => ({
        ...prev,
        stats: { ...prev.stats, ...stats },
        isLoading: false,
        lastUpdate: Date.now()
      }))

      return stats
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Search collections
  const searchCollections = useCallback(async (
    query: string,
    marketplaces: string[] = defaultMarketplaces
  ): Promise<NFTCollection[]> => {
    const filters: SearchFilters = {
      limit: 20,
      sortBy: 'recent'
    }

    const response = await aggregateCollections(marketplaces, filters)
    
    // Filter results by query
    return response.data.filter(collection =>
      collection.name.toLowerCase().includes(query.toLowerCase()) ||
      collection.description.toLowerCase().includes(query.toLowerCase()) ||
      collection.slug.toLowerCase().includes(query.toLowerCase())
    )
  }, [aggregateCollections, defaultMarketplaces])

  // Search assets
  const searchAssets = useCallback(async (
    query: string,
    marketplaces: string[] = defaultMarketplaces
  ): Promise<NFTAsset[]> => {
    const filters: SearchFilters = {
      limit: 20,
      sortBy: 'recent'
    }

    const response = await aggregateAssets(marketplaces, filters)
    
    // Filter results by query
    return response.data.filter(asset =>
      asset.name.toLowerCase().includes(query.toLowerCase()) ||
      asset.description?.toLowerCase().includes(query.toLowerCase()) ||
      asset.collection.name.toLowerCase().includes(query.toLowerCase())
    )
  }, [aggregateAssets, defaultMarketplaces])

  // Get marketplaces
  const getMarketplaces = useCallback((): MarketplaceAPI[] => {
    return marketplaceAPIManager.getMarketplaces()
  }, [])

  // Clear cache
  const clearCache = useCallback(() => {
    marketplaceAPIManager.clearCache()
    toast.success('Cache Cleared', {
      description: 'Marketplace API cache has been cleared'
    })
  }, [])

  // Refresh data
  const refresh = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      // Refresh stats for default marketplaces
      await getAllStats(defaultMarketplaces)
      
      setState(prev => ({ ...prev, isLoading: false, lastUpdate: Date.now() }))
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: (error as Error).message,
        isLoading: false
      }))
    }
  }, [getAllStats, defaultMarketplaces])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    getCollections,
    aggregateCollections,
    getAssets,
    aggregateAssets,
    getEvents,
    getStats,
    getAllStats,
    searchCollections,
    searchAssets,
    getMarketplaces,
    clearCache,
    refresh,
    clearError
  }
}

// Simplified hook for collection data
export const useMarketplaceCollections = (marketplaces?: string[]) => {
  const { getCollections, aggregateCollections, state } = useMarketplaceAPI()

  const loadCollections = useCallback(async (filters?: SearchFilters) => {
    if (marketplaces && marketplaces.length > 1) {
      return aggregateCollections(marketplaces, filters)
    } else if (marketplaces && marketplaces.length === 1) {
      return getCollections(marketplaces[0], filters)
    } else {
      return aggregateCollections(['opensea', 'looksrare'], filters)
    }
  }, [getCollections, aggregateCollections, marketplaces])

  return {
    collections: state.collections,
    loadCollections,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for marketplace statistics
export const useMarketplaceStats = () => {
  const { getAllStats, state } = useMarketplaceAPI()

  const loadStats = useCallback(async (marketplaces: string[]) => {
    return getAllStats(marketplaces)
  }, [getAllStats])

  const getMarketplaceComparison = useCallback(() => {
    const stats = Object.entries(state.stats)
    
    return stats.map(([marketplace, data]) => ({
      marketplace,
      volume24h: data.volume24h,
      sales24h: data.sales24h,
      averagePrice24h: data.averagePrice24h,
      activeListings: data.activeListings,
      uniqueTraders24h: data.uniqueTraders24h
    }))
  }, [state.stats])

  return {
    stats: state.stats,
    loadStats,
    getMarketplaceComparison,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for real-time marketplace events
export const useMarketplaceEvents = () => {
  const { getEvents, state } = useMarketplaceAPI()

  const loadEvents = useCallback(async (
    marketplace: string,
    filters?: SearchFilters
  ) => {
    return getEvents(marketplace, filters)
  }, [getEvents])

  const getRecentSales = useCallback(() => {
    return state.events
      .filter(event => event.type === 'sale')
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .slice(0, 10)
  }, [state.events])

  const getRecentListings = useCallback(() => {
    return state.events
      .filter(event => event.type === 'listing')
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .slice(0, 10)
  }, [state.events])

  return {
    events: state.events,
    loadEvents,
    getRecentSales,
    getRecentListings,
    isLoading: state.isLoading,
    error: state.error
  }
}
