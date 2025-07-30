import { type Address } from 'viem'

export interface MarketplaceAPI {
  id: string
  name: string
  baseUrl: string
  apiKey?: string
  rateLimit: RateLimit
  endpoints: MarketplaceEndpoints
  features: MarketplaceFeature[]
  status: APIStatus
  lastSync: string
}

export interface RateLimit {
  requestsPerSecond: number
  requestsPerMinute: number
  requestsPerHour: number
  burstLimit: number
}

export interface MarketplaceEndpoints {
  collections: string
  assets: string
  orders: string
  events: string
  stats: string
  account: string
  bundles?: string
  offers?: string
}

export enum MarketplaceFeature {
  COLLECTIONS = 'collections',
  ASSETS = 'assets',
  ORDERS = 'orders',
  EVENTS = 'events',
  STATS = 'stats',
  ACCOUNT = 'account',
  BUNDLES = 'bundles',
  OFFERS = 'offers',
  REAL_TIME = 'real_time',
  HISTORICAL = 'historical'
}

export enum APIStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  ERROR = 'error',
  RATE_LIMITED = 'rate_limited',
  MAINTENANCE = 'maintenance'
}

export interface NFTCollection {
  id: string
  slug: string
  name: string
  description: string
  imageUrl: string
  bannerImageUrl?: string
  contractAddress: Address
  chainId: number
  totalSupply: number
  ownersCount: number
  floorPrice: number
  floorPriceUSD: number
  volume24h: number
  volume7d: number
  volume30d: number
  volumeTotal: number
  marketCap: number
  averagePrice: number
  sales24h: number
  sales7d: number
  sales30d: number
  salesTotal: number
  createdAt: string
  verified: boolean
  featured: boolean
  category: string
  traits: CollectionTrait[]
  socialLinks: SocialLinks
  royalties: RoyaltyInfo
  marketplace: string
  lastUpdated: string
}

export interface CollectionTrait {
  traitType: string
  values: TraitValue[]
  count: number
}

export interface TraitValue {
  value: string
  count: number
  percentage: number
  floorPrice?: number
}

export interface SocialLinks {
  website?: string
  discord?: string
  twitter?: string
  instagram?: string
  telegram?: string
  medium?: string
}

export interface RoyaltyInfo {
  recipient: Address
  percentage: number
  enforced: boolean
}

export interface NFTAsset {
  id: string
  tokenId: string
  contractAddress: Address
  chainId: number
  name: string
  description?: string
  imageUrl: string
  animationUrl?: string
  externalUrl?: string
  attributes: NFTAttribute[]
  rarity: RarityInfo
  collection: CollectionInfo
  owner: Address
  creator?: Address
  currentPrice?: PriceInfo
  lastSale?: SaleInfo
  orders: OrderInfo[]
  transferHistory: TransferInfo[]
  marketplace: string
  lastUpdated: string
}

export interface NFTAttribute {
  traitType: string
  value: string
  displayType?: string
  rarity?: number
  count?: number
  percentage?: number
}

export interface RarityInfo {
  rank: number
  score: number
  tier: string
  percentile: number
  method: string
}

export interface CollectionInfo {
  id: string
  name: string
  slug: string
  imageUrl: string
  verified: boolean
}

export interface PriceInfo {
  amount: string
  currency: string
  currencyAddress: Address
  usdValue: number
  marketplace: string
  listingType: 'fixed' | 'auction' | 'dutch_auction'
  expirationTime?: string
}

export interface SaleInfo {
  price: string
  currency: string
  currencyAddress: Address
  usdValue: number
  buyer: Address
  seller: Address
  timestamp: string
  transactionHash: string
  marketplace: string
}

export interface OrderInfo {
  id: string
  type: 'listing' | 'offer' | 'bid'
  price: string
  currency: string
  currencyAddress: Address
  usdValue: number
  maker: Address
  taker?: Address
  startTime: string
  endTime?: string
  status: 'active' | 'filled' | 'cancelled' | 'expired'
  marketplace: string
}

export interface TransferInfo {
  from: Address
  to: Address
  timestamp: string
  transactionHash: string
  price?: string
  currency?: string
  marketplace?: string
}

export interface MarketplaceEvent {
  id: string
  type: EventType
  contractAddress: Address
  tokenId: string
  from?: Address
  to?: Address
  price?: string
  currency?: string
  currencyAddress?: Address
  usdValue?: number
  timestamp: string
  blockNumber: number
  transactionHash: string
  logIndex: number
  marketplace: string
  collection: CollectionInfo
  asset: AssetInfo
}

export enum EventType {
  SALE = 'sale',
  LISTING = 'listing',
  OFFER = 'offer',
  BID = 'bid',
  TRANSFER = 'transfer',
  MINT = 'mint',
  BURN = 'burn',
  CANCEL = 'cancel',
  ACCEPT = 'accept'
}

export interface AssetInfo {
  tokenId: string
  name: string
  imageUrl: string
  rarity?: RarityInfo
}

export interface MarketplaceStats {
  marketplace: string
  volume24h: number
  volume7d: number
  volume30d: number
  volumeTotal: number
  sales24h: number
  sales7d: number
  sales30d: number
  salesTotal: number
  averagePrice24h: number
  averagePrice7d: number
  averagePrice30d: number
  activeListings: number
  uniqueTraders24h: number
  uniqueTraders7d: number
  uniqueTraders30d: number
  topCollections: TopCollection[]
  lastUpdated: string
}

export interface TopCollection {
  collection: CollectionInfo
  volume24h: number
  volume7d: number
  sales24h: number
  sales7d: number
  floorPrice: number
  floorPriceChange24h: number
  rank: number
}

export interface APIResponse<T> {
  data: T
  success: boolean
  error?: string
  pagination?: PaginationInfo
  rateLimit?: RateLimitInfo
  marketplace: string
  timestamp: string
}

export interface PaginationInfo {
  page: number
  limit: number
  total: number
  hasNext: boolean
  hasPrevious: boolean
  nextCursor?: string
  previousCursor?: string
}

export interface RateLimitInfo {
  remaining: number
  reset: number
  limit: number
}

export interface SearchFilters {
  collections?: string[]
  priceMin?: number
  priceMax?: number
  currency?: string
  traits?: Record<string, string[]>
  rarity?: string[]
  status?: string[]
  sortBy?: 'price' | 'rarity' | 'recent' | 'oldest'
  sortOrder?: 'asc' | 'desc'
  limit?: number
  offset?: number
  cursor?: string
}

export class MarketplaceAPIManager {
  private static instance: MarketplaceAPIManager
  private apis = new Map<string, MarketplaceAPI>()
  private rateLimiters = new Map<string, RateLimiter>()
  private cache = new Map<string, CacheEntry>()
  private eventListeners = new Set<(event: MarketplaceAPIEvent) => void>()

  private constructor() {
    this.initializeMarketplaceAPIs()
  }

  static getInstance(): MarketplaceAPIManager {
    if (!MarketplaceAPIManager.instance) {
      MarketplaceAPIManager.instance = new MarketplaceAPIManager()
    }
    return MarketplaceAPIManager.instance
  }

  /**
   * Initialize marketplace APIs
   */
  private initializeMarketplaceAPIs(): void {
    // OpenSea API
    this.apis.set('opensea', {
      id: 'opensea',
      name: 'OpenSea',
      baseUrl: 'https://api.opensea.io/api/v1',
      rateLimit: {
        requestsPerSecond: 4,
        requestsPerMinute: 240,
        requestsPerHour: 14400,
        burstLimit: 10
      },
      endpoints: {
        collections: '/collections',
        assets: '/assets',
        orders: '/orders',
        events: '/events',
        stats: '/collection/{slug}/stats',
        account: '/account'
      },
      features: [
        MarketplaceFeature.COLLECTIONS,
        MarketplaceFeature.ASSETS,
        MarketplaceFeature.ORDERS,
        MarketplaceFeature.EVENTS,
        MarketplaceFeature.STATS,
        MarketplaceFeature.ACCOUNT,
        MarketplaceFeature.BUNDLES
      ],
      status: APIStatus.ACTIVE,
      lastSync: new Date().toISOString()
    })

    // LooksRare API
    this.apis.set('looksrare', {
      id: 'looksrare',
      name: 'LooksRare',
      baseUrl: 'https://api.looksrare.org/api/v1',
      rateLimit: {
        requestsPerSecond: 10,
        requestsPerMinute: 600,
        requestsPerHour: 36000,
        burstLimit: 20
      },
      endpoints: {
        collections: '/collections',
        assets: '/tokens',
        orders: '/orders',
        events: '/events',
        stats: '/collections/stats',
        account: '/accounts'
      },
      features: [
        MarketplaceFeature.COLLECTIONS,
        MarketplaceFeature.ASSETS,
        MarketplaceFeature.ORDERS,
        MarketplaceFeature.EVENTS,
        MarketplaceFeature.STATS
      ],
      status: APIStatus.ACTIVE,
      lastSync: new Date().toISOString()
    })

    // X2Y2 API
    this.apis.set('x2y2', {
      id: 'x2y2',
      name: 'X2Y2',
      baseUrl: 'https://api.x2y2.org/v1',
      rateLimit: {
        requestsPerSecond: 5,
        requestsPerMinute: 300,
        requestsPerHour: 18000,
        burstLimit: 15
      },
      endpoints: {
        collections: '/collections',
        assets: '/tokens',
        orders: '/orders',
        events: '/events',
        stats: '/stats',
        account: '/account'
      },
      features: [
        MarketplaceFeature.COLLECTIONS,
        MarketplaceFeature.ASSETS,
        MarketplaceFeature.ORDERS,
        MarketplaceFeature.EVENTS,
        MarketplaceFeature.STATS
      ],
      status: APIStatus.ACTIVE,
      lastSync: new Date().toISOString()
    })

    // Blur API
    this.apis.set('blur', {
      id: 'blur',
      name: 'Blur',
      baseUrl: 'https://api.blur.io/v1',
      rateLimit: {
        requestsPerSecond: 8,
        requestsPerMinute: 480,
        requestsPerHour: 28800,
        burstLimit: 25
      },
      endpoints: {
        collections: '/collections',
        assets: '/tokens',
        orders: '/orders',
        events: '/events',
        stats: '/collections/stats',
        account: '/portfolio'
      },
      features: [
        MarketplaceFeature.COLLECTIONS,
        MarketplaceFeature.ASSETS,
        MarketplaceFeature.ORDERS,
        MarketplaceFeature.EVENTS,
        MarketplaceFeature.STATS,
        MarketplaceFeature.REAL_TIME
      ],
      status: APIStatus.ACTIVE,
      lastSync: new Date().toISOString()
    })

    // Initialize rate limiters
    for (const api of Array.from(this.apis.values())) {
      this.rateLimiters.set(api.id, new RateLimiter(api.rateLimit))
    }
  }

  /**
   * Get collections from marketplace
   */
  async getCollections(
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<NFTCollection[]>> {
    const api = this.apis.get(marketplace)
    if (!api) {
      throw new Error(`Marketplace API not found: ${marketplace}`)
    }

    const cacheKey = `collections_${marketplace}_${JSON.stringify(filters)}`
    const cached = this.getFromCache(cacheKey)
    if (cached) {
      return cached as APIResponse<NFTCollection[]>
    }

    try {
      await this.checkRateLimit(marketplace)

      // Mock API call - in real implementation, this would make actual HTTP requests
      const collections = await this.mockGetCollections(marketplace, filters)

      const response: APIResponse<NFTCollection[]> = {
        data: collections,
        success: true,
        marketplace,
        timestamp: new Date().toISOString(),
        pagination: {
          page: 1,
          limit: filters?.limit || 20,
          total: collections.length,
          hasNext: false,
          hasPrevious: false
        }
      }

      this.setCache(cacheKey, response, 300) // Cache for 5 minutes
      return response

    } catch (error) {
      this.emitEvent({
        type: 'api_error',
        marketplace,
        error: error as Error,
        timestamp: Date.now()
      })

      return {
        data: [],
        success: false,
        error: (error as Error).message,
        marketplace,
        timestamp: new Date().toISOString()
      }
    }
  }

  /**
   * Get assets from marketplace
   */
  async getAssets(
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<NFTAsset[]>> {
    const api = this.apis.get(marketplace)
    if (!api) {
      throw new Error(`Marketplace API not found: ${marketplace}`)
    }

    const cacheKey = `assets_${marketplace}_${JSON.stringify(filters)}`
    const cached = this.getFromCache(cacheKey)
    if (cached) {
      return cached as APIResponse<NFTAsset[]>
    }

    try {
      await this.checkRateLimit(marketplace)

      // Mock API call
      const assets = await this.mockGetAssets(marketplace, filters)

      const response: APIResponse<NFTAsset[]> = {
        data: assets,
        success: true,
        marketplace,
        timestamp: new Date().toISOString(),
        pagination: {
          page: 1,
          limit: filters?.limit || 20,
          total: assets.length,
          hasNext: false,
          hasPrevious: false
        }
      }

      this.setCache(cacheKey, response, 180) // Cache for 3 minutes
      return response

    } catch (error) {
      this.emitEvent({
        type: 'api_error',
        marketplace,
        error: error as Error,
        timestamp: Date.now()
      })

      return {
        data: [],
        success: false,
        error: (error as Error).message,
        marketplace,
        timestamp: new Date().toISOString()
      }
    }
  }

  /**
   * Get marketplace events
   */
  async getEvents(
    marketplace: string,
    filters?: SearchFilters
  ): Promise<APIResponse<MarketplaceEvent[]>> {
    const api = this.apis.get(marketplace)
    if (!api) {
      throw new Error(`Marketplace API not found: ${marketplace}`)
    }

    try {
      await this.checkRateLimit(marketplace)

      // Mock API call
      const events = await this.mockGetEvents(marketplace, filters)

      const response: APIResponse<MarketplaceEvent[]> = {
        data: events,
        success: true,
        marketplace,
        timestamp: new Date().toISOString(),
        pagination: {
          page: 1,
          limit: filters?.limit || 50,
          total: events.length,
          hasNext: false,
          hasPrevious: false
        }
      }

      return response

    } catch (error) {
      this.emitEvent({
        type: 'api_error',
        marketplace,
        error: error as Error,
        timestamp: Date.now()
      })

      return {
        data: [],
        success: false,
        error: (error as Error).message,
        marketplace,
        timestamp: new Date().toISOString()
      }
    }
  }

  /**
   * Get marketplace statistics
   */
  async getStats(marketplace: string): Promise<APIResponse<MarketplaceStats>> {
    const api = this.apis.get(marketplace)
    if (!api) {
      throw new Error(`Marketplace API not found: ${marketplace}`)
    }

    const cacheKey = `stats_${marketplace}`
    const cached = this.getFromCache(cacheKey)
    if (cached) {
      return cached as APIResponse<MarketplaceStats>
    }

    try {
      await this.checkRateLimit(marketplace)

      // Mock API call
      const stats = await this.mockGetStats(marketplace)

      const response: APIResponse<MarketplaceStats> = {
        data: stats,
        success: true,
        marketplace,
        timestamp: new Date().toISOString()
      }

      this.setCache(cacheKey, response, 600) // Cache for 10 minutes
      return response

    } catch (error) {
      this.emitEvent({
        type: 'api_error',
        marketplace,
        error: error as Error,
        timestamp: Date.now()
      })

      return {
        data: {} as MarketplaceStats,
        success: false,
        error: (error as Error).message,
        marketplace,
        timestamp: new Date().toISOString()
      }
    }
  }

  /**
   * Aggregate data from multiple marketplaces
   */
  async aggregateCollections(
    marketplaces: string[],
    filters?: SearchFilters
  ): Promise<AggregatedResponse<NFTCollection[]>> {
    const responses = await Promise.allSettled(
      marketplaces.map(marketplace => this.getCollections(marketplace, filters))
    )

    const successful = responses
      .filter((result): result is PromiseFulfilledResult<APIResponse<NFTCollection[]>> => 
        result.status === 'fulfilled' && result.value.success
      )
      .map(result => result.value)

    const failed = responses
      .filter((result): result is PromiseRejectedResult => result.status === 'rejected')
      .map(result => result.reason)

    const allCollections = successful.flatMap(response => response.data)
    const uniqueCollections = this.deduplicateCollections(allCollections)

    return {
      data: uniqueCollections,
      sources: successful.map(r => r.marketplace),
      errors: failed,
      totalSources: marketplaces.length,
      successfulSources: successful.length,
      timestamp: new Date().toISOString()
    }
  }

  /**
   * Mock API implementations (replace with real API calls)
   */
  private async mockGetCollections(marketplace: string, _filters?: SearchFilters): Promise<NFTCollection[]> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 100 + Math.random() * 200))

    return [
      {
        id: `${marketplace}_collection_1`,
        slug: 'bored-ape-yacht-club',
        name: 'Bored Ape Yacht Club',
        description: 'A collection of 10,000 unique Bored Ape NFTs',
        imageUrl: '/nft/bayc.jpg',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        chainId: 1,
        totalSupply: 10000,
        ownersCount: 6400,
        floorPrice: 15.5,
        floorPriceUSD: 27900,
        volume24h: 450.2,
        volume7d: 2100.8,
        volume30d: 8500.5,
        volumeTotal: 850000,
        marketCap: 155000000,
        averagePrice: 18.2,
        sales24h: 25,
        sales7d: 180,
        sales30d: 750,
        salesTotal: 125000,
        createdAt: '2021-04-23T00:00:00Z',
        verified: true,
        featured: true,
        category: 'Art',
        traits: [
          {
            traitType: 'Background',
            values: [
              { value: 'Blue', count: 1000, percentage: 10, floorPrice: 15.5 },
              { value: 'Gold', count: 100, percentage: 1, floorPrice: 85.0 }
            ],
            count: 10
          }
        ],
        socialLinks: {
          website: 'https://boredapeyachtclub.com',
          discord: 'https://discord.gg/3P5K3dzgdB',
          twitter: 'https://twitter.com/BoredApeYC'
        },
        royalties: {
          recipient: '0x0000000000000000000000000000000000000001',
          percentage: 2.5,
          enforced: true
        },
        marketplace,
        lastUpdated: new Date().toISOString()
      }
    ]
  }

  private async mockGetAssets(marketplace: string, _filters?: SearchFilters): Promise<NFTAsset[]> {
    await new Promise(resolve => setTimeout(resolve, 150 + Math.random() * 250))

    return [
      {
        id: `${marketplace}_asset_1`,
        tokenId: '1',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        chainId: 1,
        name: 'Bored Ape #1',
        description: 'A unique Bored Ape NFT',
        imageUrl: '/nft/bayc-1.jpg',
        attributes: [
          {
            traitType: 'Background',
            value: 'Blue',
            rarity: 10,
            count: 1000,
            percentage: 10
          }
        ],
        rarity: {
          rank: 1,
          score: 344.5,
          tier: 'Legendary',
          percentile: 99.99,
          method: 'statistical'
        },
        collection: {
          id: 'bored-ape-yacht-club',
          name: 'Bored Ape Yacht Club',
          slug: 'bored-ape-yacht-club',
          imageUrl: '/nft/bayc.jpg',
          verified: true
        },
        owner: '0x0000000000000000000000000000000000000001',
        currentPrice: {
          amount: '25.5',
          currency: 'ETH',
          currencyAddress: '0x0000000000000000000000000000000000000000',
          usdValue: 45900,
          marketplace,
          listingType: 'fixed'
        },
        orders: [],
        transferHistory: [],
        marketplace,
        lastUpdated: new Date().toISOString()
      }
    ]
  }

  private async mockGetEvents(marketplace: string, _filters?: SearchFilters): Promise<MarketplaceEvent[]> {
    await new Promise(resolve => setTimeout(resolve, 100 + Math.random() * 150))

    return [
      {
        id: `${marketplace}_event_1`,
        type: EventType.SALE,
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        tokenId: '1',
        from: '0x0000000000000000000000000000000000000001',
        to: '0x0000000000000000000000000000000000000002',
        price: '25.5',
        currency: 'ETH',
        currencyAddress: '0x0000000000000000000000000000000000000000',
        usdValue: 45900,
        timestamp: new Date().toISOString(),
        blockNumber: 18500000,
        transactionHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
        logIndex: 0,
        marketplace,
        collection: {
          id: 'bored-ape-yacht-club',
          name: 'Bored Ape Yacht Club',
          slug: 'bored-ape-yacht-club',
          imageUrl: '/nft/bayc.jpg',
          verified: true
        },
        asset: {
          tokenId: '1',
          name: 'Bored Ape #1',
          imageUrl: '/nft/bayc-1.jpg'
        }
      }
    ]
  }

  private async mockGetStats(marketplace: string): Promise<MarketplaceStats> {
    await new Promise(resolve => setTimeout(resolve, 200 + Math.random() * 300))

    return {
      marketplace,
      volume24h: 15000000,
      volume7d: 105000000,
      volume30d: 420000000,
      volumeTotal: 35000000000,
      sales24h: 2500,
      sales7d: 17500,
      sales30d: 75000,
      salesTotal: 12500000,
      averagePrice24h: 6000,
      averagePrice7d: 6000,
      averagePrice30d: 5600,
      activeListings: 2500000,
      uniqueTraders24h: 1200,
      uniqueTraders7d: 8400,
      uniqueTraders30d: 35000,
      topCollections: [
        {
          collection: {
            id: 'bored-ape-yacht-club',
            name: 'Bored Ape Yacht Club',
            slug: 'bored-ape-yacht-club',
            imageUrl: '/nft/bayc.jpg',
            verified: true
          },
          volume24h: 450200,
          volume7d: 2100800,
          sales24h: 25,
          sales7d: 180,
          floorPrice: 15.5,
          floorPriceChange24h: 2.5,
          rank: 1
        }
      ],
      lastUpdated: new Date().toISOString()
    }
  }

  /**
   * Deduplicate collections by contract address
   */
  private deduplicateCollections(collections: NFTCollection[]): NFTCollection[] {
    const seen = new Set<string>()
    return collections.filter(collection => {
      const key = `${collection.contractAddress}_${collection.chainId}`
      if (seen.has(key)) {
        return false
      }
      seen.add(key)
      return true
    })
  }

  /**
   * Rate limiting
   */
  private async checkRateLimit(marketplace: string): Promise<void> {
    const rateLimiter = this.rateLimiters.get(marketplace)
    if (rateLimiter) {
      await rateLimiter.checkLimit()
    }
  }

  /**
   * Cache management
   */
  private getFromCache(key: string): any {
    const entry = this.cache.get(key)
    if (entry && entry.expiresAt > Date.now()) {
      return entry.data
    }
    this.cache.delete(key)
    return null
  }

  private setCache(key: string, data: any, ttlSeconds: number): void {
    this.cache.set(key, {
      data,
      expiresAt: Date.now() + ttlSeconds * 1000
    })
  }

  /**
   * Get marketplace APIs
   */
  getMarketplaces(): MarketplaceAPI[] {
    return Array.from(this.apis.values())
  }

  /**
   * Get marketplace API
   */
  getMarketplace(id: string): MarketplaceAPI | null {
    return this.apis.get(id) || null
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: MarketplaceAPIEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in marketplace API event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: MarketplaceAPIEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear cache
   */
  clearCache(): void {
    this.cache.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clearCache()
    this.eventListeners.clear()
  }
}

interface CacheEntry {
  data: any
  expiresAt: number
}

interface AggregatedResponse<T> {
  data: T
  sources: string[]
  errors: any[]
  totalSources: number
  successfulSources: number
  timestamp: string
}

class RateLimiter {
  private requests: number[] = []
  private rateLimit: RateLimit

  constructor(rateLimit: RateLimit) {
    this.rateLimit = rateLimit
  }

  async checkLimit(): Promise<void> {
    const now = Date.now()
    
    // Clean old requests
    this.requests = this.requests.filter(time => now - time < 1000)
    
    // Check rate limit
    if (this.requests.length >= this.rateLimit.requestsPerSecond) {
      const oldestRequest = Math.min(...this.requests)
      const waitTime = 1000 - (now - oldestRequest)
      await new Promise(resolve => setTimeout(resolve, waitTime))
    }
    
    this.requests.push(now)
  }
}

export interface MarketplaceAPIEvent {
  type: 'api_success' | 'api_error' | 'rate_limit_hit' | 'cache_hit' | 'cache_miss'
  marketplace: string
  endpoint?: string
  error?: Error
  timestamp: number
}

// Export singleton instance
export const marketplaceAPIManager = MarketplaceAPIManager.getInstance()
