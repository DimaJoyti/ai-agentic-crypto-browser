import { type Address } from 'viem'

export interface NFTCollection {
  id: string
  contractAddress: Address
  name: string
  symbol: string
  description: string
  imageUrl: string
  bannerUrl?: string
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
  listedCount: number
  listedPercentage: number
  royaltyFee: number
  createdAt: string
  verified: boolean
  featured: boolean
  trending: boolean
  category: NFTCategory
  blockchain: string
  metadata: CollectionMetadata
  stats: CollectionStats
  socialLinks: SocialLinks
}

export interface CollectionMetadata {
  creator: Address
  creatorName?: string
  website?: string
  discord?: string
  twitter?: string
  instagram?: string
  telegram?: string
  medium?: string
  github?: string
  externalUrl?: string
  wikiUrl?: string
  traits: TraitInfo[]
  rarityEnabled: boolean
  revealDate?: string
  mintPrice?: number
  maxSupply?: number
}

export interface TraitInfo {
  traitType: string
  valueCount: number
  values: TraitValue[]
  rarity: number
}

export interface TraitValue {
  value: string
  count: number
  percentage: number
  rarity: number
}

export interface CollectionStats {
  sales24h: number
  sales7d: number
  sales30d: number
  averagePrice24h: number
  averagePrice7d: number
  averagePrice30d: number
  priceChange24h: number
  priceChange7d: number
  priceChange30d: number
  volumeChange24h: number
  volumeChange7d: number
  volumeChange30d: number
  uniqueOwners: number
  ownershipDistribution: OwnershipTier[]
}

export interface OwnershipTier {
  tier: string
  count: number
  percentage: number
  minTokens: number
  maxTokens: number
}

export interface SocialLinks {
  website?: string
  discord?: string
  twitter?: string
  instagram?: string
  telegram?: string
  medium?: string
  github?: string
}

export interface NFTToken {
  id: string
  tokenId: string
  contractAddress: Address
  collectionId: string
  name: string
  description?: string
  imageUrl: string
  animationUrl?: string
  externalUrl?: string
  chainId: number
  owner: Address
  creator?: Address
  minter?: Address
  mintedAt: string
  lastSale?: NFTSale
  currentListing?: NFTListing
  rarity: RarityInfo
  traits: NFTTrait[]
  metadata: TokenMetadata
  history: NFTActivity[]
  verified: boolean
}

export interface RarityInfo {
  rank: number
  score: number
  totalSupply: number
  rarityTier: 'Common' | 'Uncommon' | 'Rare' | 'Epic' | 'Legendary' | 'Mythic'
  percentile: number
}

export interface NFTTrait {
  traitType: string
  value: string
  displayType?: string
  rarity: number
  count: number
  percentage: number
}

export interface TokenMetadata {
  standard: 'ERC721' | 'ERC1155'
  tokenURI: string
  metadataURI?: string
  contentType?: string
  fileSize?: number
  dimensions?: {
    width: number
    height: number
  }
  duration?: number
  hasAnimation: boolean
  hasAudio: boolean
  isVideo: boolean
  isImage: boolean
  is3D: boolean
}

export interface NFTSale {
  id: string
  tokenId: string
  contractAddress: Address
  seller: Address
  buyer: Address
  price: string
  priceUSD: number
  currency: string
  marketplace: string
  transactionHash: string
  blockNumber: number
  timestamp: string
  gasUsed: string
  gasFee: string
}

export interface NFTListing {
  id: string
  tokenId: string
  contractAddress: Address
  seller: Address
  price: string
  priceUSD: number
  currency: string
  marketplace: string
  listingType: 'fixed' | 'auction' | 'dutch'
  startTime: string
  endTime?: string
  reservePrice?: string
  currentBid?: string
  bidCount?: number
  status: 'active' | 'sold' | 'cancelled' | 'expired'
  createdAt: string
  updatedAt: string
}

export interface NFTActivity {
  id: string
  type: 'mint' | 'sale' | 'transfer' | 'listing' | 'bid' | 'cancel' | 'burn'
  tokenId: string
  contractAddress: Address
  from?: Address
  to?: Address
  price?: string
  priceUSD?: number
  currency?: string
  marketplace?: string
  transactionHash: string
  blockNumber: number
  timestamp: string
  gasUsed?: string
  gasFee?: string
}

export enum NFTCategory {
  ART = 'art',
  COLLECTIBLES = 'collectibles',
  GAMING = 'gaming',
  METAVERSE = 'metaverse',
  MUSIC = 'music',
  PHOTOGRAPHY = 'photography',
  SPORTS = 'sports',
  TRADING_CARDS = 'trading_cards',
  UTILITY = 'utility',
  VIRTUAL_WORLDS = 'virtual_worlds',
  DOMAIN_NAMES = 'domain_names',
  MEMES = 'memes'
}

export interface NFTSearchFilters {
  category?: NFTCategory
  priceMin?: number
  priceMax?: number
  chainId?: number
  verified?: boolean
  hasListings?: boolean
  traits?: Record<string, string[]>
  rarityMin?: number
  rarityMax?: number
  sortBy?: 'price' | 'rarity' | 'recent' | 'oldest' | 'name'
  sortOrder?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export interface CollectionSearchFilters {
  category?: NFTCategory
  chainId?: number
  verified?: boolean
  trending?: boolean
  featured?: boolean
  volumeMin?: number
  volumeMax?: number
  floorPriceMin?: number
  floorPriceMax?: number
  sortBy?: 'volume' | 'floor_price' | 'market_cap' | 'created' | 'name'
  sortOrder?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export interface SearchResult<T> {
  items: T[]
  total: number
  hasMore: boolean
  nextCursor?: string
}

export class NFTDiscoveryEngine {
  private static instance: NFTDiscoveryEngine
  private collections = new Map<string, NFTCollection>()
  private tokens = new Map<string, NFTToken>()
  private trendingCollections: string[] = []
  private featuredCollections: string[] = []
  private eventListeners = new Set<(event: NFTDiscoveryEvent) => void>()

  private constructor() {
    this.initializeMockData()
  }

  static getInstance(): NFTDiscoveryEngine {
    if (!NFTDiscoveryEngine.instance) {
      NFTDiscoveryEngine.instance = new NFTDiscoveryEngine()
    }
    return NFTDiscoveryEngine.instance
  }

  /**
   * Initialize mock data for demonstration
   */
  private initializeMockData(): void {
    // Mock collections
    const mockCollections: Partial<NFTCollection>[] = [
      {
        id: 'bored-ape-yacht-club',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        name: 'Bored Ape Yacht Club',
        symbol: 'BAYC',
        description: 'A collection of 10,000 unique Bored Ape NFTs',
        imageUrl: '/nft/bayc.jpg',
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
        verified: true,
        featured: true,
        trending: true,
        category: NFTCategory.ART
      },
      {
        id: 'cryptopunks',
        contractAddress: '0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB',
        name: 'CryptoPunks',
        symbol: 'PUNK',
        description: 'One of the earliest NFT projects on Ethereum',
        imageUrl: '/nft/cryptopunks.jpg',
        chainId: 1,
        totalSupply: 10000,
        ownersCount: 3500,
        floorPrice: 65.8,
        floorPriceUSD: 118440,
        volume24h: 1200.5,
        volume7d: 5500.2,
        volume30d: 22000.8,
        volumeTotal: 2500000,
        marketCap: 658000000,
        verified: true,
        featured: true,
        trending: false,
        category: NFTCategory.COLLECTIBLES
      },
      {
        id: 'azuki',
        contractAddress: '0xED5AF388653567Af2F388E6224dC7C4b3241C544',
        name: 'Azuki',
        symbol: 'AZUKI',
        description: 'A brand for the metaverse built by Chiru Labs',
        imageUrl: '/nft/azuki.jpg',
        chainId: 1,
        totalSupply: 10000,
        ownersCount: 5200,
        floorPrice: 8.2,
        floorPriceUSD: 14760,
        volume24h: 280.5,
        volume7d: 1200.8,
        volume30d: 4800.2,
        volumeTotal: 420000,
        marketCap: 82000000,
        verified: true,
        featured: false,
        trending: true,
        category: NFTCategory.ART
      }
    ]

    mockCollections.forEach((collection, index) => {
      const fullCollection: NFTCollection = {
        id: collection.id || `collection_${index}`,
        contractAddress: collection.contractAddress || '0x0000000000000000000000000000000000000000',
        name: collection.name || `Collection ${index}`,
        symbol: collection.symbol || `COL${index}`,
        description: collection.description || '',
        imageUrl: collection.imageUrl || '',
        chainId: collection.chainId || 1,
        totalSupply: collection.totalSupply || 10000,
        ownersCount: collection.ownersCount || 5000,
        floorPrice: collection.floorPrice || 1.0,
        floorPriceUSD: collection.floorPriceUSD || 1800,
        volume24h: collection.volume24h || 100,
        volume7d: collection.volume7d || 700,
        volume30d: collection.volume30d || 3000,
        volumeTotal: collection.volumeTotal || 100000,
        marketCap: collection.marketCap || 10000000,
        averagePrice: (collection.floorPrice || 1.0) * 1.5,
        listedCount: Math.floor((collection.totalSupply || 10000) * 0.15),
        listedPercentage: 15,
        royaltyFee: 2.5,
        createdAt: '2021-04-01T00:00:00Z',
        verified: collection.verified || false,
        featured: collection.featured || false,
        trending: collection.trending || false,
        category: collection.category || NFTCategory.ART,
        blockchain: 'ethereum',
        metadata: this.generateMockMetadata(),
        stats: this.generateMockStats(),
        socialLinks: this.generateMockSocialLinks()
      }

      this.collections.set(fullCollection.id, fullCollection)

      if (fullCollection.trending) {
        this.trendingCollections.push(fullCollection.id)
      }
      if (fullCollection.featured) {
        this.featuredCollections.push(fullCollection.id)
      }
    })
  }

  /**
   * Generate mock metadata
   */
  private generateMockMetadata(): CollectionMetadata {
    return {
      creator: '0x0000000000000000000000000000000000000000',
      creatorName: 'Anonymous Creator',
      website: 'https://example.com',
      discord: 'https://discord.gg/example',
      twitter: 'https://twitter.com/example',
      traits: [
        {
          traitType: 'Background',
          valueCount: 10,
          values: [
            { value: 'Blue', count: 1000, percentage: 10, rarity: 10 },
            { value: 'Red', count: 500, percentage: 5, rarity: 20 }
          ],
          rarity: 15
        }
      ],
      rarityEnabled: true,
      mintPrice: 0.08,
      maxSupply: 10000
    }
  }

  /**
   * Generate mock stats
   */
  private generateMockStats(): CollectionStats {
    return {
      sales24h: 50,
      sales7d: 350,
      sales30d: 1500,
      averagePrice24h: 12.5,
      averagePrice7d: 11.8,
      averagePrice30d: 10.2,
      priceChange24h: 5.2,
      priceChange7d: -2.1,
      priceChange30d: 15.8,
      volumeChange24h: 8.5,
      volumeChange7d: -5.2,
      volumeChange30d: 25.2,
      uniqueOwners: 6400,
      ownershipDistribution: [
        { tier: 'Whales (10+)', count: 50, percentage: 0.8, minTokens: 10, maxTokens: 100 },
        { tier: 'Large (5-9)', count: 200, percentage: 3.1, minTokens: 5, maxTokens: 9 },
        { tier: 'Medium (2-4)', count: 800, percentage: 12.5, minTokens: 2, maxTokens: 4 },
        { tier: 'Small (1)', count: 5350, percentage: 83.6, minTokens: 1, maxTokens: 1 }
      ]
    }
  }

  /**
   * Generate mock social links
   */
  private generateMockSocialLinks(): SocialLinks {
    return {
      website: 'https://example.com',
      discord: 'https://discord.gg/example',
      twitter: 'https://twitter.com/example'
    }
  }

  /**
   * Search collections
   */
  async searchCollections(
    query?: string,
    filters?: CollectionSearchFilters
  ): Promise<SearchResult<NFTCollection>> {
    let collections = Array.from(this.collections.values())

    // Apply text search
    if (query) {
      const searchTerm = query.toLowerCase()
      collections = collections.filter(collection =>
        collection.name.toLowerCase().includes(searchTerm) ||
        collection.symbol.toLowerCase().includes(searchTerm) ||
        collection.description.toLowerCase().includes(searchTerm)
      )
    }

    // Apply filters
    if (filters) {
      if (filters.category) {
        collections = collections.filter(c => c.category === filters.category)
      }
      if (filters.chainId) {
        collections = collections.filter(c => c.chainId === filters.chainId)
      }
      if (filters.verified !== undefined) {
        collections = collections.filter(c => c.verified === filters.verified)
      }
      if (filters.trending !== undefined) {
        collections = collections.filter(c => c.trending === filters.trending)
      }
      if (filters.featured !== undefined) {
        collections = collections.filter(c => c.featured === filters.featured)
      }
      if (filters.volumeMin !== undefined) {
        collections = collections.filter(c => c.volume24h >= filters.volumeMin!)
      }
      if (filters.volumeMax !== undefined) {
        collections = collections.filter(c => c.volume24h <= filters.volumeMax!)
      }
      if (filters.floorPriceMin !== undefined) {
        collections = collections.filter(c => c.floorPrice >= filters.floorPriceMin!)
      }
      if (filters.floorPriceMax !== undefined) {
        collections = collections.filter(c => c.floorPrice <= filters.floorPriceMax!)
      }
    }

    // Apply sorting
    if (filters?.sortBy) {
      collections.sort((a, b) => {
        let aValue: number, bValue: number

        switch (filters.sortBy) {
          case 'volume':
            aValue = a.volume24h
            bValue = b.volume24h
            break
          case 'floor_price':
            aValue = a.floorPrice
            bValue = b.floorPrice
            break
          case 'market_cap':
            aValue = a.marketCap
            bValue = b.marketCap
            break
          case 'created':
            aValue = new Date(a.createdAt).getTime()
            bValue = new Date(b.createdAt).getTime()
            break
          case 'name':
            return filters.sortOrder === 'desc' 
              ? b.name.localeCompare(a.name)
              : a.name.localeCompare(b.name)
          default:
            aValue = a.volume24h
            bValue = b.volume24h
        }

        return filters.sortOrder === 'desc' ? bValue - aValue : aValue - bValue
      })
    }

    // Apply pagination
    const limit = filters?.limit || 20
    const offset = filters?.offset || 0
    const total = collections.length
    const items = collections.slice(offset, offset + limit)

    return {
      items,
      total,
      hasMore: offset + limit < total,
      nextCursor: offset + limit < total ? (offset + limit).toString() : undefined
    }
  }

  /**
   * Get trending collections
   */
  async getTrendingCollections(limit: number = 10): Promise<NFTCollection[]> {
    const trending = this.trendingCollections
      .map(id => this.collections.get(id))
      .filter(Boolean) as NFTCollection[]

    return trending.slice(0, limit)
  }

  /**
   * Get featured collections
   */
  async getFeaturedCollections(limit: number = 10): Promise<NFTCollection[]> {
    const featured = this.featuredCollections
      .map(id => this.collections.get(id))
      .filter(Boolean) as NFTCollection[]

    return featured.slice(0, limit)
  }

  /**
   * Get collection by ID
   */
  async getCollection(id: string): Promise<NFTCollection | null> {
    return this.collections.get(id) || null
  }

  /**
   * Get collection by contract address
   */
  async getCollectionByContract(contractAddress: Address): Promise<NFTCollection | null> {
    for (const collection of Array.from(this.collections.values())) {
      if (collection.contractAddress.toLowerCase() === contractAddress.toLowerCase()) {
        return collection
      }
    }
    return null
  }

  /**
   * Search NFT tokens
   */
  async searchTokens(
    _query?: string,
    _filters?: NFTSearchFilters
  ): Promise<SearchResult<NFTToken>> {
    // Mock implementation - in real app, this would query actual NFT data
    const mockTokens: NFTToken[] = []

    return {
      items: mockTokens,
      total: 0,
      hasMore: false
    }
  }

  /**
   * Get categories
   */
  getCategories(): { value: NFTCategory; label: string; count: number }[] {
    const categories = Object.values(NFTCategory).map(category => {
      const count = Array.from(this.collections.values())
        .filter(c => c.category === category).length

      return {
        value: category,
        label: category.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()),
        count
      }
    })

    return categories.filter(c => c.count > 0)
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: NFTDiscoveryEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in NFT discovery event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: NFTDiscoveryEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.collections.clear()
    this.tokens.clear()
    this.trendingCollections = []
    this.featuredCollections = []
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface NFTDiscoveryEvent {
  type: 'collection_updated' | 'trending_updated' | 'featured_updated' | 'search_completed'
  collection?: NFTCollection
  collections?: NFTCollection[]
  searchQuery?: string
  timestamp: number
}

// Export singleton instance
export const nftDiscoveryEngine = NFTDiscoveryEngine.getInstance()
