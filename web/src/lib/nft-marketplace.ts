import { type Address, type Hash } from 'viem'

export interface NFTMarketplace {
  id: string
  name: string
  version: string
  chainId: number
  contractAddresses: MarketplaceContracts
  supportedFeatures: MarketplaceFeature[]
  fees: MarketplaceFees
  metadata: MarketplaceMetadata
}

export interface MarketplaceContracts {
  exchange: Address
  transferManager: Address
  royaltyRegistry?: Address
  weth: Address
  protocolFeeRecipient: Address
}

export interface MarketplaceFees {
  protocolFee: number // percentage
  royaltyFee: number // percentage
  gasEstimate: {
    listing: string
    purchase: string
    bid: string
    cancel: string
  }
}

export interface MarketplaceMetadata {
  description: string
  website: string
  documentation: string
  logo: string
  volume24h: string
  volume7d: string
  volumeTotal: string
  activeListings: number
  totalUsers: number
  averagePrice: number
}

export enum MarketplaceFeature {
  FIXED_PRICE = 'fixed_price',
  AUCTION = 'auction',
  DUTCH_AUCTION = 'dutch_auction',
  BUNDLE_SALES = 'bundle_sales',
  COLLECTION_OFFERS = 'collection_offers',
  TRAIT_OFFERS = 'trait_offers',
  ROYALTIES = 'royalties',
  LAZY_MINTING = 'lazy_minting',
  BULK_OPERATIONS = 'bulk_operations'
}

export interface NFTListing {
  id: string
  marketplace: string
  tokenId: string
  contractAddress: Address
  seller: Address
  price: string
  priceUSD: number
  currency: string
  currencyAddress: Address
  listingType: ListingType
  startTime: string
  endTime?: string
  reservePrice?: string
  currentBid?: string
  bidCount: number
  status: ListingStatus
  signature?: string
  salt?: string
  createdAt: string
  updatedAt: string
  metadata: ListingMetadata
}

export interface ListingMetadata {
  isPrivate: boolean
  allowedBuyers?: Address[]
  bundleItems?: BundleItem[]
  royaltyRecipient?: Address
  royaltyAmount?: string
  protocolData?: string
}

export interface BundleItem {
  contractAddress: Address
  tokenId: string
  amount: string
}

export enum ListingType {
  FIXED_PRICE = 'fixed_price',
  AUCTION = 'auction',
  DUTCH_AUCTION = 'dutch_auction',
  BUNDLE = 'bundle'
}

export enum ListingStatus {
  ACTIVE = 'active',
  SOLD = 'sold',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired',
  INVALID = 'invalid'
}

export interface NFTOffer {
  id: string
  marketplace: string
  tokenId?: string
  contractAddress: Address
  offerer: Address
  price: string
  priceUSD: number
  currency: string
  currencyAddress: Address
  offerType: OfferType
  expiration: string
  signature: string
  salt: string
  status: OfferStatus
  createdAt: string
  updatedAt: string
  metadata: OfferMetadata
}

export interface OfferMetadata {
  isPrivate: boolean
  allowedSellers?: Address[]
  traitCriteria?: TraitCriteria[]
  collectionCriteria?: CollectionCriteria
  protocolData?: string
}

export interface TraitCriteria {
  traitType: string
  value: string
}

export interface CollectionCriteria {
  contractAddress: Address
  minFloorPrice?: string
  maxFloorPrice?: string
  traits?: TraitCriteria[]
}

export enum OfferType {
  TOKEN = 'token',
  COLLECTION = 'collection',
  TRAIT = 'trait'
}

export enum OfferStatus {
  ACTIVE = 'active',
  ACCEPTED = 'accepted',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired',
  INVALID = 'invalid'
}

export interface NFTTransaction {
  id: string
  hash?: Hash
  marketplace: string
  type: TransactionType
  tokenId: string
  contractAddress: Address
  seller: Address
  buyer: Address
  price: string
  priceUSD: number
  currency: string
  currencyAddress: Address
  royaltyAmount?: string
  protocolFee?: string
  gasUsed?: string
  gasFee?: string
  status: 'pending' | 'confirmed' | 'failed' | 'reverted'
  timestamp: string
  blockNumber?: number
  error?: string
  metadata?: TransactionMetadata
}

export interface TransactionMetadata {
  listingId?: string
  offerId?: string
  bundleItems?: BundleItem[]
  isPrivateSale?: boolean
  referrer?: Address
  protocolData?: string
}

export enum TransactionType {
  PURCHASE = 'purchase',
  SALE = 'sale',
  BID = 'bid',
  ACCEPT_BID = 'accept_bid',
  CANCEL_LISTING = 'cancel_listing',
  CANCEL_OFFER = 'cancel_offer',
  TRANSFER = 'transfer'
}

export interface NFTPortfolio {
  userAddress: Address
  totalValue: number
  totalItems: number
  collections: PortfolioCollection[]
  recentActivity: NFTTransaction[]
  topGainers: PortfolioItem[]
  topLosers: PortfolioItem[]
  floorValueChange24h: number
  lastUpdate: number
}

export interface PortfolioCollection {
  contractAddress: Address
  name: string
  symbol: string
  imageUrl: string
  itemCount: number
  floorPrice: number
  totalValue: number
  averageBuyPrice: number
  unrealizedPnL: number
  unrealizedPnLPercentage: number
  items: PortfolioItem[]
}

export interface PortfolioItem {
  tokenId: string
  contractAddress: Address
  name: string
  imageUrl: string
  buyPrice?: number
  buyDate?: string
  currentValue: number
  lastSalePrice?: number
  unrealizedPnL: number
  unrealizedPnLPercentage: number
  rarity?: {
    rank: number
    score: number
    tier: string
  }
}

export interface MarketplaceFilters {
  marketplace?: string
  priceMin?: number
  priceMax?: number
  currency?: string
  listingType?: ListingType
  status?: ListingStatus
  collections?: Address[]
  traits?: Record<string, string[]>
  sortBy?: 'price' | 'recent' | 'ending_soon' | 'rarity'
  sortOrder?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export class NFTMarketplaceIntegration {
  private static instance: NFTMarketplaceIntegration
  private marketplaces = new Map<string, NFTMarketplace>()
  private listings = new Map<string, NFTListing>()
  private offers = new Map<string, NFTOffer>()
  private transactions = new Map<string, NFTTransaction>()
  private portfolios = new Map<string, NFTPortfolio>()
  private eventListeners = new Set<(event: MarketplaceEvent) => void>()

  private constructor() {
    this.initializeMarketplaces()
  }

  static getInstance(): NFTMarketplaceIntegration {
    if (!NFTMarketplaceIntegration.instance) {
      NFTMarketplaceIntegration.instance = new NFTMarketplaceIntegration()
    }
    return NFTMarketplaceIntegration.instance
  }

  /**
   * Initialize supported marketplaces
   */
  private initializeMarketplaces(): void {
    // OpenSea
    this.marketplaces.set('opensea', {
      id: 'opensea',
      name: 'OpenSea',
      version: '1.5.0',
      chainId: 1,
      contractAddresses: {
        exchange: '0x00000000000000ADc04C56Bf30aC9d3c0aAF14dC',
        transferManager: '0x0000000000000000000000000000000000000000',
        weth: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
        protocolFeeRecipient: '0x0000a26b00c1F0DF003000390027140000fAa719'
      },
      supportedFeatures: [
        MarketplaceFeature.FIXED_PRICE,
        MarketplaceFeature.AUCTION,
        MarketplaceFeature.DUTCH_AUCTION,
        MarketplaceFeature.BUNDLE_SALES,
        MarketplaceFeature.COLLECTION_OFFERS,
        MarketplaceFeature.TRAIT_OFFERS,
        MarketplaceFeature.ROYALTIES
      ],
      fees: {
        protocolFee: 2.5,
        royaltyFee: 10,
        gasEstimate: {
          listing: '150000',
          purchase: '200000',
          bid: '120000',
          cancel: '80000'
        }
      },
      metadata: {
        description: 'The largest NFT marketplace',
        website: 'https://opensea.io',
        documentation: 'https://docs.opensea.io',
        logo: '/logos/opensea.svg',
        volume24h: '15000000',
        volume7d: '105000000',
        volumeTotal: '35000000000',
        activeListings: 2500000,
        totalUsers: 1800000,
        averagePrice: 0.85
      }
    })

    // LooksRare
    this.marketplaces.set('looksrare', {
      id: 'looksrare',
      name: 'LooksRare',
      version: '1.0.0',
      chainId: 1,
      contractAddresses: {
        exchange: '0x59728544B08AB483533076417FbBB2fD0B17CE3a',
        transferManager: '0xf42aa99F011A1fA7CDA90E5E98b277E306BcA83e',
        weth: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
        protocolFeeRecipient: '0x5924A28caAF1cc016617874a2f0C3710d881f3c1'
      },
      supportedFeatures: [
        MarketplaceFeature.FIXED_PRICE,
        MarketplaceFeature.AUCTION,
        MarketplaceFeature.COLLECTION_OFFERS,
        MarketplaceFeature.ROYALTIES
      ],
      fees: {
        protocolFee: 2.0,
        royaltyFee: 10,
        gasEstimate: {
          listing: '140000',
          purchase: '180000',
          bid: '110000',
          cancel: '70000'
        }
      },
      metadata: {
        description: 'Community-first NFT marketplace',
        website: 'https://looksrare.org',
        documentation: 'https://docs.looksrare.org',
        logo: '/logos/looksrare.svg',
        volume24h: '2500000',
        volume7d: '17500000',
        volumeTotal: '8500000000',
        activeListings: 450000,
        totalUsers: 320000,
        averagePrice: 1.2
      }
    })

    // X2Y2
    this.marketplaces.set('x2y2', {
      id: 'x2y2',
      name: 'X2Y2',
      version: '1.0.0',
      chainId: 1,
      contractAddresses: {
        exchange: '0x74312363e45DCaBA76c59ec49a7Aa8A65a67EeD3',
        transferManager: '0x0000000000000000000000000000000000000000',
        weth: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
        protocolFeeRecipient: '0xd823C605807cC5E6Bd6fC0d7e58F1677d0c2e3a8'
      },
      supportedFeatures: [
        MarketplaceFeature.FIXED_PRICE,
        MarketplaceFeature.AUCTION,
        MarketplaceFeature.BUNDLE_SALES,
        MarketplaceFeature.ROYALTIES
      ],
      fees: {
        protocolFee: 0.5,
        royaltyFee: 10,
        gasEstimate: {
          listing: '130000',
          purchase: '170000',
          bid: '100000',
          cancel: '60000'
        }
      },
      metadata: {
        description: 'Low-fee NFT marketplace',
        website: 'https://x2y2.io',
        documentation: 'https://docs.x2y2.io',
        logo: '/logos/x2y2.svg',
        volume24h: '1800000',
        volume7d: '12600000',
        volumeTotal: '4200000000',
        activeListings: 320000,
        totalUsers: 180000,
        averagePrice: 0.95
      }
    })
  }

  /**
   * Get marketplace listings
   */
  async getListings(_filters?: MarketplaceFilters): Promise<NFTListing[]> {
    // Mock implementation - in real app, this would query actual marketplace APIs
    const mockListings: NFTListing[] = [
      {
        id: 'listing_1',
        marketplace: 'opensea',
        tokenId: '1',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        seller: '0x0000000000000000000000000000000000000001',
        price: '15.5',
        priceUSD: 27900,
        currency: 'ETH',
        currencyAddress: '0x0000000000000000000000000000000000000000',
        listingType: ListingType.FIXED_PRICE,
        startTime: new Date().toISOString(),
        bidCount: 0,
        status: ListingStatus.ACTIVE,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        metadata: {
          isPrivate: false
        }
      }
    ]

    return mockListings
  }

  /**
   * Create listing
   */
  async createListing(
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    seller: Address,
    listingType: ListingType = ListingType.FIXED_PRICE,
    endTime?: string
  ): Promise<NFTListing> {
    const listing: NFTListing = {
      id: `listing_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
      marketplace,
      tokenId,
      contractAddress,
      seller,
      price,
      priceUSD: parseFloat(price) * 1800, // Mock USD conversion
      currency,
      currencyAddress: currency === 'ETH' ? '0x0000000000000000000000000000000000000000' : '0x0000000000000000000000000000000000000000',
      listingType,
      startTime: new Date().toISOString(),
      endTime,
      bidCount: 0,
      status: ListingStatus.ACTIVE,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      metadata: {
        isPrivate: false
      }
    }

    this.listings.set(listing.id, listing)

    // Emit event
    this.emitEvent({
      type: 'listing_created',
      listing,
      timestamp: Date.now()
    })

    return listing
  }

  /**
   * Purchase NFT
   */
  async purchaseNFT(
    listingId: string,
    buyer: Address,
    price?: string
  ): Promise<NFTTransaction> {
    const listing = this.listings.get(listingId)
    if (!listing) {
      throw new Error(`Listing not found: ${listingId}`)
    }

    if (listing.status !== ListingStatus.ACTIVE) {
      throw new Error(`Listing is not active: ${listing.status}`)
    }

    const transaction: NFTTransaction = {
      id: `tx_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
      marketplace: listing.marketplace,
      type: TransactionType.PURCHASE,
      tokenId: listing.tokenId,
      contractAddress: listing.contractAddress,
      seller: listing.seller,
      buyer,
      price: price || listing.price,
      priceUSD: parseFloat(price || listing.price) * 1800,
      currency: listing.currency,
      currencyAddress: listing.currencyAddress,
      status: 'pending',
      timestamp: new Date().toISOString(),
      metadata: {
        listingId
      }
    }

    this.transactions.set(transaction.id, transaction)

    try {
      // Execute the purchase (mock implementation)
      const result = await this.performPurchaseTransaction(transaction)
      
      transaction.status = 'confirmed'
      transaction.hash = result.hash
      transaction.blockNumber = result.blockNumber
      transaction.gasUsed = result.gasUsed

      // Update listing status
      listing.status = ListingStatus.SOLD
      listing.updatedAt = new Date().toISOString()

      // Emit success event
      this.emitEvent({
        type: 'purchase_completed',
        transaction,
        listing,
        timestamp: Date.now()
      })

    } catch (error) {
      transaction.status = 'failed'
      transaction.error = (error as Error).message

      // Emit failure event
      this.emitEvent({
        type: 'purchase_failed',
        transaction,
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }

    return transaction
  }

  /**
   * Perform purchase transaction (mock implementation)
   */
  private async performPurchaseTransaction(_transaction: NFTTransaction): Promise<{
    hash: Hash
    blockNumber: number
    gasUsed: string
  }> {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 2000 + Math.random() * 3000))

    // Simulate 95% success rate
    if (Math.random() < 0.95) {
      return {
        hash: `0x${Math.random().toString(16).substring(2, 66)}` as Hash,
        blockNumber: Math.floor(Math.random() * 1000000) + 18000000,
        gasUsed: '180000'
      }
    } else {
      throw new Error('Purchase failed: Insufficient funds')
    }
  }

  /**
   * Create offer
   */
  async createOffer(
    tokenId: string,
    contractAddress: Address,
    price: string,
    currency: string,
    marketplace: string,
    offerer: Address,
    offerType: OfferType = OfferType.TOKEN,
    expiration?: string
  ): Promise<NFTOffer> {
    const offer: NFTOffer = {
      id: `offer_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      marketplace,
      tokenId: offerType === OfferType.TOKEN ? tokenId : undefined,
      contractAddress,
      offerer,
      price,
      priceUSD: parseFloat(price) * 1800,
      currency,
      currencyAddress: currency === 'ETH' ? '0x0000000000000000000000000000000000000000' : '0x0000000000000000000000000000000000000000',
      offerType,
      expiration: expiration || new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
      signature: `0x${Math.random().toString(16).substr(2, 130)}`,
      salt: Math.random().toString(),
      status: OfferStatus.ACTIVE,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      metadata: {
        isPrivate: false
      }
    }

    this.offers.set(offer.id, offer)

    // Emit event
    this.emitEvent({
      type: 'offer_created',
      offer,
      timestamp: Date.now()
    })

    return offer
  }

  /**
   * Get user portfolio
   */
  async getUserPortfolio(userAddress: Address): Promise<NFTPortfolio> {
    // Mock implementation - in real app, this would aggregate user's NFTs
    const mockPortfolio: NFTPortfolio = {
      userAddress,
      totalValue: 125000,
      totalItems: 15,
      collections: [
        {
          contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
          name: 'Bored Ape Yacht Club',
          symbol: 'BAYC',
          imageUrl: '/nft/bayc.jpg',
          itemCount: 3,
          floorPrice: 15.5,
          totalValue: 85000,
          averageBuyPrice: 12.0,
          unrealizedPnL: 15000,
          unrealizedPnLPercentage: 21.4,
          items: [
            {
              tokenId: '1234',
              contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
              name: 'Bored Ape #1234',
              imageUrl: '/nft/bayc-1234.jpg',
              buyPrice: 12.0,
              buyDate: '2023-01-15',
              currentValue: 28500,
              unrealizedPnL: 6900,
              unrealizedPnLPercentage: 31.9
            }
          ]
        }
      ],
      recentActivity: [],
      topGainers: [],
      topLosers: [],
      floorValueChange24h: 5.2,
      lastUpdate: Date.now()
    }

    this.portfolios.set(userAddress.toLowerCase(), mockPortfolio)
    return mockPortfolio
  }

  /**
   * Get marketplaces
   */
  getMarketplaces(): NFTMarketplace[] {
    return Array.from(this.marketplaces.values())
  }

  /**
   * Get marketplace by ID
   */
  getMarketplace(id: string): NFTMarketplace | null {
    return this.marketplaces.get(id) || null
  }

  /**
   * Get transaction
   */
  getTransaction(id: string): NFTTransaction | null {
    return this.transactions.get(id) || null
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: MarketplaceEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in marketplace event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: MarketplaceEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.listings.clear()
    this.offers.clear()
    this.transactions.clear()
    this.portfolios.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface MarketplaceEvent {
  type: 'listing_created' | 'listing_updated' | 'purchase_completed' | 'purchase_failed' | 'offer_created' | 'offer_accepted' | 'portfolio_updated'
  listing?: NFTListing
  offer?: NFTOffer
  transaction?: NFTTransaction
  portfolio?: NFTPortfolio
  error?: Error
  timestamp: number
}

// Export singleton instance
export const nftMarketplaceIntegration = NFTMarketplaceIntegration.getInstance()
