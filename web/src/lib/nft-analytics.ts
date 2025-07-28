import { createPublicClient, http, type Address } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum RarityTier {
  COMMON = 'common',
  UNCOMMON = 'uncommon', 
  RARE = 'rare',
  EPIC = 'epic',
  LEGENDARY = 'legendary',
  MYTHIC = 'mythic'
}

export enum TrendDirection {
  UP = 'up',
  DOWN = 'down',
  STABLE = 'stable'
}

export interface RarityScore {
  tokenId: string
  contractAddress: Address
  rarityScore: number
  rarityRank: number
  totalSupply: number
  percentile: number
  tier: RarityTier
  attributes: AttributeRarity[]
}

export interface AttributeRarity {
  trait_type: string
  value: string
  count: number
  percentage: number
  rarityScore: number
}

export interface PriceHistory {
  timestamp: number
  price: string
  marketplace: string
  transactionHash: string
  from: Address
  to: Address
  type: 'sale' | 'mint' | 'transfer'
}

export interface FloorPriceData {
  contractAddress: Address
  chainId: number
  currentFloor: string
  previousFloor: string
  change24h: string
  changePercentage: string
  trend: TrendDirection
  marketplace: string
  lastUpdated: number
  history: {
    timestamp: number
    price: string
    marketplace: string
  }[]
}

export interface VolumeData {
  contractAddress: Address
  chainId: number
  volume1h: string
  volume24h: string
  volume7d: string
  volume30d: string
  volumeAll: string
  sales1h: number
  sales24h: number
  sales7d: number
  sales30d: number
  salesAll: number
  averagePrice24h: string
  medianPrice24h: string
  lastUpdated: number
}

export interface MarketTrend {
  contractAddress: Address
  chainId: number
  trendDirection: TrendDirection
  trendStrength: number // 0-100
  momentum: number // -100 to 100
  volatility: number // 0-100
  liquidityScore: number // 0-100
  marketCapChange24h: string
  ownershipDistribution: {
    uniqueOwners: number
    topHolderPercentage: number
    whaleCount: number
    averageHolding: number
  }
  sentiment: {
    score: number // -100 to 100
    volume: number
    mentions: number
    positiveRatio: number
  }
}

export interface NFTValuation {
  tokenId: string
  contractAddress: Address
  estimatedValue: string
  confidence: number // 0-100
  valuationMethod: 'comparable_sales' | 'rarity_based' | 'ml_model' | 'hybrid'
  factors: {
    rarityScore: number
    floorPriceMultiplier: number
    recentSalesWeight: number
    marketTrendWeight: number
    liquidityFactor: number
  }
  comparableSales: {
    tokenId: string
    price: string
    timestamp: number
    similarityScore: number
  }[]
  priceRange: {
    low: string
    high: string
    mostLikely: string
  }
  lastUpdated: number
}

export interface CollectionAnalytics {
  contractAddress: Address
  chainId: number
  name: string
  totalSupply: number
  uniqueOwners: number
  floorPrice: string
  marketCap: string
  volume24h: string
  sales24h: number
  averagePrice: string
  medianPrice: string
  priceDistribution: {
    range: string
    count: number
    percentage: number
  }[]
  rarityDistribution: {
    tier: RarityTier
    count: number
    percentage: number
    floorPrice: string
  }[]
  topSales: {
    tokenId: string
    price: string
    timestamp: number
    marketplace: string
  }[]
  marketMetrics: {
    liquidityScore: number
    volatilityIndex: number
    trendScore: number
    momentumIndicator: number
  }
  predictions: {
    floorPrice7d: string
    floorPrice30d: string
    confidence: number
  }
}

export class NFTAnalyticsService {
  private static instance: NFTAnalyticsService
  private clients: Map<number, any> = new Map()
  private rarityScores: Map<string, RarityScore> = new Map()
  private floorPrices: Map<string, FloorPriceData> = new Map()
  private volumeData: Map<string, VolumeData> = new Map()
  private marketTrends: Map<string, MarketTrend> = new Map()
  private valuations: Map<string, NFTValuation> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeMockData()
  }

  static getInstance(): NFTAnalyticsService {
    if (!NFTAnalyticsService.instance) {
      NFTAnalyticsService.instance = new NFTAnalyticsService()
    }
    return NFTAnalyticsService.instance
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
          console.warn(`Failed to initialize analytics client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMockData() {
    // Mock Rarity Scores
    const mockRarity1: RarityScore = {
      tokenId: '1234',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      rarityScore: 2847.5,
      rarityRank: 245,
      totalSupply: 10000,
      percentile: 97.55,
      tier: RarityTier.RARE,
      attributes: [
        {
          trait_type: 'Background',
          value: 'Blue',
          count: 1520,
          percentage: 15.2,
          rarityScore: 6.58
        },
        {
          trait_type: 'Eyes',
          value: 'Laser Eyes',
          count: 210,
          percentage: 2.1,
          rarityScore: 47.62
        },
        {
          trait_type: 'Hat',
          value: 'Crown',
          count: 120,
          percentage: 1.2,
          rarityScore: 83.33
        }
      ]
    }

    this.rarityScores.set('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D-1234', mockRarity1)

    // Mock Floor Price Data
    const mockFloorPrice: FloorPriceData = {
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      chainId: 1,
      currentFloor: '15.5',
      previousFloor: '14.8',
      change24h: '0.7',
      changePercentage: '4.73',
      trend: TrendDirection.UP,
      marketplace: 'OpenSea',
      lastUpdated: Date.now(),
      history: [
        { timestamp: Date.now() - 86400000, price: '14.8', marketplace: 'OpenSea' },
        { timestamp: Date.now() - 172800000, price: '15.2', marketplace: 'Blur' },
        { timestamp: Date.now() - 259200000, price: '14.5', marketplace: 'OpenSea' }
      ]
    }

    this.floorPrices.set('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D', mockFloorPrice)

    // Mock Volume Data
    const mockVolume: VolumeData = {
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      chainId: 1,
      volume1h: '125.5',
      volume24h: '2850.2',
      volume7d: '18420.8',
      volume30d: '85600.5',
      volumeAll: '1250000.0',
      sales1h: 8,
      sales24h: 185,
      sales7d: 1250,
      sales30d: 5800,
      salesAll: 125000,
      averagePrice24h: '15.41',
      medianPrice24h: '14.85',
      lastUpdated: Date.now()
    }

    this.volumeData.set('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D', mockVolume)

    // Mock Market Trend
    const mockTrend: MarketTrend = {
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      chainId: 1,
      trendDirection: TrendDirection.UP,
      trendStrength: 75,
      momentum: 25,
      volatility: 45,
      liquidityScore: 85,
      marketCapChange24h: '4.73',
      ownershipDistribution: {
        uniqueOwners: 5432,
        topHolderPercentage: 2.5,
        whaleCount: 125,
        averageHolding: 1.84
      },
      sentiment: {
        score: 65,
        volume: 1250,
        mentions: 850,
        positiveRatio: 0.72
      }
    }

    this.marketTrends.set('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D', mockTrend)

    // Mock NFT Valuation
    const mockValuation: NFTValuation = {
      tokenId: '1234',
      contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address,
      estimatedValue: '18.75',
      confidence: 85,
      valuationMethod: 'hybrid',
      factors: {
        rarityScore: 2847.5,
        floorPriceMultiplier: 1.21,
        recentSalesWeight: 0.35,
        marketTrendWeight: 0.25,
        liquidityFactor: 0.85
      },
      comparableSales: [
        {
          tokenId: '2156',
          price: '19.2',
          timestamp: Date.now() - 86400000,
          similarityScore: 0.92
        },
        {
          tokenId: '7834',
          price: '17.8',
          timestamp: Date.now() - 172800000,
          similarityScore: 0.88
        }
      ],
      priceRange: {
        low: '16.5',
        high: '21.2',
        mostLikely: '18.75'
      },
      lastUpdated: Date.now()
    }

    this.valuations.set('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D-1234', mockValuation)
  }

  // Public methods
  async getRarityScore(contractAddress: Address, tokenId: string): Promise<RarityScore | null> {
    const key = `${contractAddress.toLowerCase()}-${tokenId}`
    return this.rarityScores.get(key) || null
  }

  async getFloorPriceData(contractAddress: Address): Promise<FloorPriceData | null> {
    return this.floorPrices.get(contractAddress.toLowerCase()) || null
  }

  async getVolumeData(contractAddress: Address): Promise<VolumeData | null> {
    return this.volumeData.get(contractAddress.toLowerCase()) || null
  }

  async getMarketTrend(contractAddress: Address): Promise<MarketTrend | null> {
    return this.marketTrends.get(contractAddress.toLowerCase()) || null
  }

  async getNFTValuation(contractAddress: Address, tokenId: string): Promise<NFTValuation | null> {
    const key = `${contractAddress.toLowerCase()}-${tokenId}`
    return this.valuations.get(key) || null
  }

  async getCollectionAnalytics(contractAddress: Address): Promise<CollectionAnalytics | null> {
    const floorData = await this.getFloorPriceData(contractAddress)
    const volumeData = await this.getVolumeData(contractAddress)
    const trendData = await this.getMarketTrend(contractAddress)

    if (!floorData || !volumeData || !trendData) return null

    // Mock collection analytics
    return {
      contractAddress,
      chainId: 1,
      name: 'Bored Ape Yacht Club',
      totalSupply: 10000,
      uniqueOwners: 5432,
      floorPrice: floorData.currentFloor,
      marketCap: (parseFloat(floorData.currentFloor) * 10000).toString(),
      volume24h: volumeData.volume24h,
      sales24h: volumeData.sales24h,
      averagePrice: volumeData.averagePrice24h,
      medianPrice: volumeData.medianPrice24h,
      priceDistribution: [
        { range: '0-10 ETH', count: 2500, percentage: 25 },
        { range: '10-20 ETH', count: 4500, percentage: 45 },
        { range: '20-50 ETH', count: 2500, percentage: 25 },
        { range: '50+ ETH', count: 500, percentage: 5 }
      ],
      rarityDistribution: [
        { tier: RarityTier.COMMON, count: 5000, percentage: 50, floorPrice: '15.5' },
        { tier: RarityTier.UNCOMMON, count: 2500, percentage: 25, floorPrice: '18.2' },
        { tier: RarityTier.RARE, count: 1500, percentage: 15, floorPrice: '25.8' },
        { tier: RarityTier.EPIC, count: 750, percentage: 7.5, floorPrice: '45.5' },
        { tier: RarityTier.LEGENDARY, count: 225, percentage: 2.25, floorPrice: '85.2' },
        { tier: RarityTier.MYTHIC, count: 25, percentage: 0.25, floorPrice: '250.0' }
      ],
      topSales: [
        { tokenId: '8817', price: '347.0', timestamp: Date.now() - 86400000, marketplace: 'OpenSea' },
        { tokenId: '3749', price: '285.5', timestamp: Date.now() - 172800000, marketplace: 'Blur' },
        { tokenId: '1420', price: '195.8', timestamp: Date.now() - 259200000, marketplace: 'OpenSea' }
      ],
      marketMetrics: {
        liquidityScore: trendData.liquidityScore,
        volatilityIndex: trendData.volatility,
        trendScore: trendData.trendStrength,
        momentumIndicator: trendData.momentum
      },
      predictions: {
        floorPrice7d: (parseFloat(floorData.currentFloor) * 1.05).toFixed(2),
        floorPrice30d: (parseFloat(floorData.currentFloor) * 1.12).toFixed(2),
        confidence: 75
      }
    }
  }

  calculateRarityScore(attributes: AttributeRarity[]): number {
    return attributes.reduce((total, attr) => total + attr.rarityScore, 0)
  }

  getRarityTier(rarityScore: number, totalSupply: number): RarityTier {
    const percentile = (rarityScore / totalSupply) * 100
    
    if (percentile >= 99.5) return RarityTier.MYTHIC
    if (percentile >= 97.5) return RarityTier.LEGENDARY
    if (percentile >= 92.5) return RarityTier.EPIC
    if (percentile >= 85) return RarityTier.RARE
    if (percentile >= 70) return RarityTier.UNCOMMON
    return RarityTier.COMMON
  }

  async getPriceHistory(contractAddress: Address, tokenId: string): Promise<PriceHistory[]> {
    // Mock price history
    return [
      {
        timestamp: Date.now() - 86400000 * 30,
        price: '16.8',
        marketplace: 'OpenSea',
        transactionHash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890',
        from: '0x1234567890123456789012345678901234567890' as Address,
        to: '0x0987654321098765432109876543210987654321' as Address,
        type: 'sale'
      },
      {
        timestamp: Date.now() - 86400000 * 60,
        price: '15.2',
        marketplace: 'LooksRare',
        transactionHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
        from: '0x2345678901234567890123456789012345678901' as Address,
        to: '0x1234567890123456789012345678901234567890' as Address,
        type: 'sale'
      }
    ]
  }

  async getTopCollections(limit: number = 10): Promise<CollectionAnalytics[]> {
    // Mock top collections
    const mockCollections = [
      await this.getCollectionAnalytics('0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as Address)
    ]
    
    return mockCollections.filter(Boolean).slice(0, limit) as CollectionAnalytics[]
  }

  async getTrendingCollections(limit: number = 10): Promise<CollectionAnalytics[]> {
    // Mock trending collections based on volume change
    return this.getTopCollections(limit)
  }

  formatPrice(price: string): string {
    const value = parseFloat(price)
    if (value >= 1000) return `${(value / 1000).toFixed(1)}K ETH`
    if (value >= 1) return `${value.toFixed(2)} ETH`
    return `${value.toFixed(4)} ETH`
  }

  formatPercentage(value: string): string {
    const num = parseFloat(value)
    const sign = num >= 0 ? '+' : ''
    return `${sign}${num.toFixed(2)}%`
  }

  formatVolume(volume: string): string {
    const value = parseFloat(volume)
    if (value >= 1000000) return `${(value / 1000000).toFixed(1)}M ETH`
    if (value >= 1000) return `${(value / 1000).toFixed(1)}K ETH`
    return `${value.toFixed(1)} ETH`
  }
}

// Export singleton instance
export const nftAnalyticsService = NFTAnalyticsService.getInstance()
