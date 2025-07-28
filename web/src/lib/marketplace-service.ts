import { createPublicClient, http, type Address, type Hash } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum MarketplaceName {
  OPENSEA = 'OpenSea',
  BLUR = 'Blur',
  LOOKSRARE = 'LooksRare',
  X2Y2 = 'X2Y2',
  FOUNDATION = 'Foundation',
  SUPERRARE = 'SuperRare'
}

export enum OrderType {
  LISTING = 'listing',
  OFFER = 'offer',
  AUCTION = 'auction',
  DUTCH_AUCTION = 'dutch_auction'
}

export enum OrderStatus {
  ACTIVE = 'active',
  FILLED = 'filled',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired'
}

export interface MarketplaceConfig {
  name: MarketplaceName
  baseUrl: string
  apiUrl: string
  contractAddress: Address
  chainId: number
  feePercentage: number
  supportedOrderTypes: OrderType[]
  isActive: boolean
}

export interface NFTListing {
  id: string
  marketplace: MarketplaceName
  nftId: string
  contractAddress: Address
  tokenId: string
  seller: Address
  price: string
  currency: Address
  orderType: OrderType
  status: OrderStatus
  startTime: number
  endTime?: number
  createdAt: number
  updatedAt: number
  signature?: string
  fees: {
    marketplace: string
    royalty: string
    total: string
  }
}

export interface NFTOffer {
  id: string
  marketplace: MarketplaceName
  nftId: string
  contractAddress: Address
  tokenId: string
  buyer: Address
  price: string
  currency: Address
  status: OrderStatus
  expirationTime: number
  createdAt: number
  signature?: string
}

export interface MarketplaceActivity {
  id: string
  marketplace: MarketplaceName
  type: 'sale' | 'listing' | 'offer' | 'transfer' | 'mint'
  nftId: string
  contractAddress: Address
  tokenId: string
  from: Address
  to: Address
  price?: string
  currency?: Address
  timestamp: number
  transactionHash: string
  blockNumber: number
}

export interface MarketplaceStats {
  marketplace: MarketplaceName
  volume24h: string
  volume7d: string
  volume30d: string
  volumeAll: string
  sales24h: number
  sales7d: number
  sales30d: number
  averagePrice24h: string
  floorPrice: string
  activeListings: number
  uniqueTraders24h: number
}

export interface CollectionStats {
  contractAddress: Address
  chainId: number
  floorPrice: string
  volume24h: string
  volume7d: string
  volume30d: string
  sales24h: number
  averagePrice24h: string
  marketCap: string
  totalSupply: number
  ownersCount: number
  listedCount: number
  marketplaceBreakdown: {
    [key in MarketplaceName]?: {
      listings: number
      volume24h: string
      floorPrice: string
    }
  }
}

export class MarketplaceService {
  private static instance: MarketplaceService
  private clients: Map<number, any> = new Map()
  private marketplaces: Map<MarketplaceName, MarketplaceConfig> = new Map()
  private listings: Map<string, NFTListing> = new Map()
  private offers: Map<string, NFTOffer> = new Map()
  private activities: Map<string, MarketplaceActivity[]> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeMarketplaces()
    this.initializeMockData()
  }

  static getInstance(): MarketplaceService {
    if (!MarketplaceService.instance) {
      MarketplaceService.instance = new MarketplaceService()
    }
    return MarketplaceService.instance
  }

  private initializeClients() {
    Object.values(SUPPORTED_CHAINS).forEach(chain => {
      if (!chain.isTestnet || chain.id === 11155111) {
        try {
          const client = createPublicClient({
            chain: {
              id: chain.id,
              name: chain.name,
              network: chain.shortName.toLowerCase(),
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: chain.rpcUrls
            } as any,
            transport: http()
          })
          this.clients.set(chain.id, client)
        } catch (error) {
          console.warn(`Failed to initialize marketplace client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMarketplaces() {
    // OpenSea
    this.marketplaces.set(MarketplaceName.OPENSEA, {
      name: MarketplaceName.OPENSEA,
      baseUrl: 'https://opensea.io',
      apiUrl: 'https://api.opensea.io/v2',
      contractAddress: '0x00000000000000ADc04C56Bf30aC9d3c0aAF14dC' as Address, // Seaport
      chainId: 1,
      feePercentage: 2.5,
      supportedOrderTypes: [OrderType.LISTING, OrderType.OFFER, OrderType.AUCTION],
      isActive: true
    })

    // Blur
    this.marketplaces.set(MarketplaceName.BLUR, {
      name: MarketplaceName.BLUR,
      baseUrl: 'https://blur.io',
      apiUrl: 'https://api.blur.io/v1',
      contractAddress: '0x000000000000Ad05Ccc4F10045630fb830B95127' as Address,
      chainId: 1,
      feePercentage: 0.5,
      supportedOrderTypes: [OrderType.LISTING, OrderType.OFFER],
      isActive: true
    })

    // LooksRare
    this.marketplaces.set(MarketplaceName.LOOKSRARE, {
      name: MarketplaceName.LOOKSRARE,
      baseUrl: 'https://looksrare.org',
      apiUrl: 'https://api.looksrare.org/api/v2',
      contractAddress: '0x59728544B08AB483533076417FbBB2fD0B17CE3a' as Address,
      chainId: 1,
      feePercentage: 2.0,
      supportedOrderTypes: [OrderType.LISTING, OrderType.OFFER, OrderType.AUCTION],
      isActive: true
    })

    // X2Y2
    this.marketplaces.set(MarketplaceName.X2Y2, {
      name: MarketplaceName.X2Y2,
      baseUrl: 'https://x2y2.io',
      apiUrl: 'https://api.x2y2.org/v1',
      contractAddress: '0x74312363e45DCaBA76c59ec49a7Aa8A65a67EeD3' as Address,
      chainId: 1,
      feePercentage: 0.5,
      supportedOrderTypes: [OrderType.LISTING, OrderType.OFFER, OrderType.AUCTION],
      isActive: true
    })

    // Foundation
    this.marketplaces.set(MarketplaceName.FOUNDATION, {
      name: MarketplaceName.FOUNDATION,
      baseUrl: 'https://foundation.app',
      apiUrl: 'https://api.foundation.app/v1',
      contractAddress: '0xcDA72070E455bb31C7690a170224Ce43623d0B6f' as Address,
      chainId: 1,
      feePercentage: 15.0,
      supportedOrderTypes: [OrderType.AUCTION],
      isActive: true
    })

    // SuperRare
    this.marketplaces.set(MarketplaceName.SUPERRARE, {
      name: MarketplaceName.SUPERRARE,
      baseUrl: 'https://superrare.com',
      apiUrl: 'https://api.superrare.com/v1',
      contractAddress: '0xb932a70A57673d89f4acfFBE830E8ed7f75Fb9e0' as Address,
      chainId: 1,
      feePercentage: 15.0,
      supportedOrderTypes: [OrderType.AUCTION, OrderType.LISTING],
      isActive: true
    })
  }

  private initializeMockData() {
    // Mock NFT Listings
    const mockListing1: NFTListing = {
      id: 'listing-1',
      marketplace: MarketplaceName.OPENSEA,
      nftId: 'bayc-1234',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      tokenId: '1234',
      seller: '0x1234567890123456789012345678901234567890' as Address,
      price: '18.5',
      currency: '0x0000000000000000000000000000000000000000' as Address, // ETH
      orderType: OrderType.LISTING,
      status: OrderStatus.ACTIVE,
      startTime: Date.now(),
      endTime: Date.now() + 86400000 * 7, // 7 days
      createdAt: Date.now() - 86400000,
      updatedAt: Date.now() - 86400000,
      fees: {
        marketplace: '0.4625', // 2.5% of 18.5
        royalty: '0.4625', // 2.5% royalty
        total: '0.925'
      }
    }

    const mockListing2: NFTListing = {
      id: 'listing-2',
      marketplace: MarketplaceName.BLUR,
      nftId: 'azuki-5678',
      contractAddress: '0xED5AF388653567Af2F388E6224dC7C4b3241C544' as Address,
      tokenId: '5678',
      seller: '0x0987654321098765432109876543210987654321' as Address,
      price: '10.2',
      currency: '0x0000000000000000000000000000000000000000' as Address,
      orderType: OrderType.LISTING,
      status: OrderStatus.ACTIVE,
      startTime: Date.now(),
      createdAt: Date.now() - 86400000 * 2,
      updatedAt: Date.now() - 86400000 * 2,
      fees: {
        marketplace: '0.051', // 0.5% of 10.2
        royalty: '0.51', // 5% royalty
        total: '0.561'
      }
    }

    this.listings.set(mockListing1.id, mockListing1)
    this.listings.set(mockListing2.id, mockListing2)

    // Mock NFT Offers
    const mockOffer1: NFTOffer = {
      id: 'offer-1',
      marketplace: MarketplaceName.OPENSEA,
      nftId: 'bayc-1234',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      tokenId: '1234',
      buyer: '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd' as Address,
      price: '17.8',
      currency: '0x0000000000000000000000000000000000000000' as Address,
      status: OrderStatus.ACTIVE,
      expirationTime: Date.now() + 86400000 * 3, // 3 days
      createdAt: Date.now() - 86400000
    }

    this.offers.set(mockOffer1.id, mockOffer1)

    // Mock Marketplace Activities
    const mockActivities: MarketplaceActivity[] = [
      {
        id: 'activity-1',
        marketplace: MarketplaceName.OPENSEA,
        type: 'sale',
        nftId: 'bayc-9999',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
        tokenId: '9999',
        from: '0x1111111111111111111111111111111111111111' as Address,
        to: '0x2222222222222222222222222222222222222222' as Address,
        price: '19.5',
        currency: '0x0000000000000000000000000000000000000000' as Address,
        timestamp: Date.now() - 3600000, // 1 hour ago
        transactionHash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890',
        blockNumber: 18500000
      },
      {
        id: 'activity-2',
        marketplace: MarketplaceName.BLUR,
        type: 'sale',
        nftId: 'azuki-1111',
        contractAddress: '0xED5AF388653567Af2F388E6224dC7C4b3241C544' as Address,
        tokenId: '1111',
        from: '0x3333333333333333333333333333333333333333' as Address,
        to: '0x4444444444444444444444444444444444444444' as Address,
        price: '9.8',
        currency: '0x0000000000000000000000000000000000000000' as Address,
        timestamp: Date.now() - 7200000, // 2 hours ago
        transactionHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
        blockNumber: 18499950
      }
    ]

    this.activities.set('recent', mockActivities)
  }

  // Public methods
  getAllMarketplaces(): MarketplaceConfig[] {
    return Array.from(this.marketplaces.values()).filter(mp => mp.isActive)
  }

  getMarketplace(name: MarketplaceName): MarketplaceConfig | undefined {
    return this.marketplaces.get(name)
  }

  getListings(contractAddress?: Address, tokenId?: string, marketplace?: MarketplaceName): NFTListing[] {
    let listings = Array.from(this.listings.values())
    
    if (contractAddress) {
      listings = listings.filter(listing => 
        listing.contractAddress.toLowerCase() === contractAddress.toLowerCase()
      )
    }
    
    if (tokenId) {
      listings = listings.filter(listing => listing.tokenId === tokenId)
    }
    
    if (marketplace) {
      listings = listings.filter(listing => listing.marketplace === marketplace)
    }
    
    return listings.filter(listing => listing.status === OrderStatus.ACTIVE)
  }

  getOffers(contractAddress?: Address, tokenId?: string, marketplace?: MarketplaceName): NFTOffer[] {
    let offers = Array.from(this.offers.values())
    
    if (contractAddress) {
      offers = offers.filter(offer => 
        offer.contractAddress.toLowerCase() === contractAddress.toLowerCase()
      )
    }
    
    if (tokenId) {
      offers = offers.filter(offer => offer.tokenId === tokenId)
    }
    
    if (marketplace) {
      offers = offers.filter(offer => offer.marketplace === marketplace)
    }
    
    return offers.filter(offer => 
      offer.status === OrderStatus.ACTIVE && 
      offer.expirationTime > Date.now()
    )
  }

  getRecentActivity(limit: number = 20): MarketplaceActivity[] {
    const activities = this.activities.get('recent') || []
    return activities
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
  }

  async createListing(
    contractAddress: Address,
    tokenId: string,
    price: string,
    marketplace: MarketplaceName,
    duration: number = 86400000 * 7 // 7 days default
  ): Promise<Hash> {
    const marketplaceConfig = this.getMarketplace(marketplace)
    if (!marketplaceConfig) {
      throw new Error('Marketplace not supported')
    }

    // Mock listing creation
    const listingId = `listing-${Date.now()}`
    const newListing: NFTListing = {
      id: listingId,
      marketplace,
      nftId: `${contractAddress}-${tokenId}`,
      contractAddress,
      tokenId,
      seller: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address, // Mock user
      price,
      currency: '0x0000000000000000000000000000000000000000' as Address,
      orderType: OrderType.LISTING,
      status: OrderStatus.ACTIVE,
      startTime: Date.now(),
      endTime: Date.now() + duration,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      fees: {
        marketplace: (parseFloat(price) * marketplaceConfig.feePercentage / 100).toString(),
        royalty: (parseFloat(price) * 0.025).toString(), // 2.5% royalty
        total: (parseFloat(price) * (marketplaceConfig.feePercentage + 2.5) / 100).toString()
      }
    }

    this.listings.set(listingId, newListing)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  async createOffer(
    contractAddress: Address,
    tokenId: string,
    price: string,
    marketplace: MarketplaceName,
    duration: number = 86400000 * 3 // 3 days default
  ): Promise<Hash> {
    const marketplaceConfig = this.getMarketplace(marketplace)
    if (!marketplaceConfig) {
      throw new Error('Marketplace not supported')
    }

    // Mock offer creation
    const offerId = `offer-${Date.now()}`
    const newOffer: NFTOffer = {
      id: offerId,
      marketplace,
      nftId: `${contractAddress}-${tokenId}`,
      contractAddress,
      tokenId,
      buyer: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address, // Mock user
      price,
      currency: '0x0000000000000000000000000000000000000000' as Address,
      status: OrderStatus.ACTIVE,
      expirationTime: Date.now() + duration,
      createdAt: Date.now()
    }

    this.offers.set(offerId, newOffer)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  async cancelListing(listingId: string): Promise<Hash> {
    const listing = this.listings.get(listingId)
    if (!listing) {
      throw new Error('Listing not found')
    }

    listing.status = OrderStatus.CANCELLED
    listing.updatedAt = Date.now()
    this.listings.set(listingId, listing)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  async cancelOffer(offerId: string): Promise<Hash> {
    const offer = this.offers.get(offerId)
    if (!offer) {
      throw new Error('Offer not found')
    }

    offer.status = OrderStatus.CANCELLED
    this.offers.set(offerId, offer)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  getMarketplaceStats(): MarketplaceStats[] {
    return Array.from(this.marketplaces.values()).map(marketplace => ({
      marketplace: marketplace.name,
      volume24h: this.getMockVolume(marketplace.name, '24h'),
      volume7d: this.getMockVolume(marketplace.name, '7d'),
      volume30d: this.getMockVolume(marketplace.name, '30d'),
      volumeAll: this.getMockVolume(marketplace.name, 'all'),
      sales24h: this.getMockSales(marketplace.name, '24h'),
      sales7d: this.getMockSales(marketplace.name, '7d'),
      sales30d: this.getMockSales(marketplace.name, '30d'),
      averagePrice24h: this.getMockAveragePrice(marketplace.name),
      floorPrice: this.getMockFloorPrice(marketplace.name),
      activeListings: this.getMockActiveListings(marketplace.name),
      uniqueTraders24h: this.getMockUniqueTraders(marketplace.name)
    }))
  }

  private getMockVolume(marketplace: MarketplaceName, period: string): string {
    const baseVolumes: Record<MarketplaceName, Record<string, string>> = {
      [MarketplaceName.OPENSEA]: { '24h': '15420', '7d': '98650', '30d': '425800', 'all': '12500000' },
      [MarketplaceName.BLUR]: { '24h': '8950', '7d': '62100', '30d': '285400', 'all': '3200000' },
      [MarketplaceName.LOOKSRARE]: { '24h': '2150', '7d': '14800', '30d': '68500', 'all': '850000' },
      [MarketplaceName.X2Y2]: { '24h': '1850', '7d': '12200', '30d': '55800', 'all': '720000' },
      [MarketplaceName.FOUNDATION]: { '24h': '450', '7d': '3200', '30d': '15800', 'all': '180000' },
      [MarketplaceName.SUPERRARE]: { '24h': '320', '7d': '2100', '30d': '9800', 'all': '125000' }
    }
    return baseVolumes[marketplace]?.[period] || '0'
  }

  private getMockSales(marketplace: MarketplaceName, period: string): number {
    const baseSales: Record<MarketplaceName, Record<string, number>> = {
      [MarketplaceName.OPENSEA]: { '24h': 1250, '7d': 8500, '30d': 35200 },
      [MarketplaceName.BLUR]: { '24h': 890, '7d': 6200, '30d': 28500 },
      [MarketplaceName.LOOKSRARE]: { '24h': 180, '7d': 1200, '30d': 5800 },
      [MarketplaceName.X2Y2]: { '24h': 150, '7d': 980, '30d': 4200 },
      [MarketplaceName.FOUNDATION]: { '24h': 25, '7d': 180, '30d': 850 },
      [MarketplaceName.SUPERRARE]: { '24h': 18, '7d': 125, '30d': 580 }
    }
    return baseSales[marketplace]?.[period] || 0
  }

  private getMockAveragePrice(marketplace: MarketplaceName): string {
    const averagePrices: Record<MarketplaceName, string> = {
      [MarketplaceName.OPENSEA]: '12.34',
      [MarketplaceName.BLUR]: '10.05',
      [MarketplaceName.LOOKSRARE]: '11.94',
      [MarketplaceName.X2Y2]: '12.33',
      [MarketplaceName.FOUNDATION]: '18.00',
      [MarketplaceName.SUPERRARE]: '17.78'
    }
    return averagePrices[marketplace] || '0'
  }

  private getMockFloorPrice(marketplace: MarketplaceName): string {
    const floorPrices: Record<MarketplaceName, string> = {
      [MarketplaceName.OPENSEA]: '8.50',
      [MarketplaceName.BLUR]: '8.45',
      [MarketplaceName.LOOKSRARE]: '8.52',
      [MarketplaceName.X2Y2]: '8.48',
      [MarketplaceName.FOUNDATION]: '12.00',
      [MarketplaceName.SUPERRARE]: '15.50'
    }
    return floorPrices[marketplace] || '0'
  }

  private getMockActiveListings(marketplace: MarketplaceName): number {
    const activeListings: Record<MarketplaceName, number> = {
      [MarketplaceName.OPENSEA]: 15420,
      [MarketplaceName.BLUR]: 8950,
      [MarketplaceName.LOOKSRARE]: 2150,
      [MarketplaceName.X2Y2]: 1850,
      [MarketplaceName.FOUNDATION]: 450,
      [MarketplaceName.SUPERRARE]: 320
    }
    return activeListings[marketplace] || 0
  }

  private getMockUniqueTraders(marketplace: MarketplaceName): number {
    const uniqueTraders: Record<MarketplaceName, number> = {
      [MarketplaceName.OPENSEA]: 5420,
      [MarketplaceName.BLUR]: 3950,
      [MarketplaceName.LOOKSRARE]: 850,
      [MarketplaceName.X2Y2]: 720,
      [MarketplaceName.FOUNDATION]: 180,
      [MarketplaceName.SUPERRARE]: 125
    }
    return uniqueTraders[marketplace] || 0
  }
}

// Export singleton instance
export const marketplaceService = MarketplaceService.getInstance()
