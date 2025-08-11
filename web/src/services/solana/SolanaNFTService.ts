import axios from 'axios'
import { NFTCollection, NFTMarketStats, NFTTrend } from '@/hooks/useSolanaNFT'

export class SolanaNFTService {
  private readonly baseURL: string
  private readonly magicEdenAPI: string
  private readonly tensorAPI: string

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    this.magicEdenAPI = 'https://api-mainnet.magiceden.dev/v2'
    this.tensorAPI = 'https://api.tensor.trade/graphql'
  }

  async getCollections(options: {
    marketplace?: NFTCollection['marketplace']
    minFloorPrice?: number
    maxFloorPrice?: number
    sortBy?: 'volume' | 'floor_price' | 'market_cap' | 'sales'
    limit?: number
  } = {}): Promise<NFTCollection[]> {
    try {
      // Try backend API first
      const response = await axios.get(`${this.baseURL}/api/solana/nft/collections`, {
        params: options
      })
      return response.data
    } catch (error) {
      console.warn('Backend API unavailable, falling back to Magic Eden')
      
      // Fallback to Magic Eden API
      const response = await axios.get(`${this.magicEdenAPI}/collections`, {
        params: {
          offset: 0,
          limit: options.limit || 50
        }
      })
      
      return response.data.map((collection: any) => ({
        id: collection.symbol,
        name: collection.name,
        symbol: collection.symbol,
        description: collection.description || '',
        image: collection.image,
        floorPrice: collection.floorPrice / 1e9 || 0, // Convert lamports to SOL
        floorPriceChange24h: 0, // Not available in basic API
        volume24h: collection.volumeAll / 1e9 || 0,
        volumeChange24h: 0,
        marketCap: (collection.floorPrice * collection.totalItems) / 1e9 || 0,
        supply: collection.totalItems || 0,
        owners: collection.uniqueHolders || 0,
        listedCount: collection.listedCount || 0,
        listedPercentage: collection.totalItems > 0 
          ? (collection.listedCount / collection.totalItems) * 100 
          : 0,
        averagePrice24h: 0,
        sales24h: 0,
        website: collection.website,
        twitter: collection.twitter,
        discord: collection.discord,
        verified: collection.verified || false,
        marketplace: 'magic-eden'
      }))
    }
  }

  async getCollection(collectionId: string): Promise<NFTCollection> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/collections/${collectionId}`)
      return response.data
    } catch (error) {
      throw new Error(`Failed to fetch collection ${collectionId}`)
    }
  }

  async getMarketStats(): Promise<NFTMarketStats> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/market-stats`)
      return response.data
    } catch (error) {
      console.warn('NFT market stats API unavailable, calculating from collections')
      
      const collections = await this.getCollections({ limit: 100 })
      
      const totalVolume24h = collections.reduce((sum, collection) => sum + collection.volume24h, 0)
      const totalSales24h = collections.reduce((sum, collection) => sum + collection.sales24h, 0)
      const averageFloorPrice = collections.length > 0
        ? collections.reduce((sum, collection) => sum + collection.floorPrice, 0) / collections.length
        : 0
      const marketCapTotal = collections.reduce((sum, collection) => sum + collection.marketCap, 0)
      
      const topByVolume = collections.sort((a, b) => b.volume24h - a.volume24h)[0]
      const topByFloor = collections.sort((a, b) => b.floorPrice - a.floorPrice)[0]

      return {
        totalVolume24h,
        totalSales24h,
        totalCollections: collections.length,
        averageFloorPrice,
        topCollectionByVolume: topByVolume?.name || '',
        topCollectionByFloor: topByFloor?.name || '',
        marketCapTotal
      }
    }
  }

  async getTrends(): Promise<NFTTrend[]> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/trends`)
      return response.data
    } catch (error) {
      console.warn('NFT trends API unavailable, using mock data')
      
      // Mock trends for development
      return [
        {
          collection: 'Mad Lads',
          metric: 'volume',
          change24h: 150.5,
          changePercentage: 25.3,
          rank: 1
        },
        {
          collection: 'Okay Bears',
          metric: 'floor_price',
          change24h: 12.8,
          changePercentage: 18.7,
          rank: 2
        },
        {
          collection: 'DeGods',
          metric: 'sales',
          change24h: 89,
          changePercentage: 15.2,
          rank: 3
        },
        {
          collection: 'Solana Monkey Business',
          metric: 'volume',
          change24h: -45.2,
          changePercentage: -12.8,
          rank: 4
        },
        {
          collection: 'Famous Fox Federation',
          metric: 'floor_price',
          change24h: 8.3,
          changePercentage: 11.4,
          rank: 5
        }
      ]
    }
  }

  async getCollectionActivity(collectionId: string, limit: number = 20): Promise<Array<{
    signature: string
    type: 'sale' | 'listing' | 'bid' | 'cancel'
    price: number
    buyer?: string
    seller?: string
    timestamp: Date
    marketplace: string
  }>> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/collections/${collectionId}/activity`, {
        params: { limit }
      })
      return response.data
    } catch (error) {
      console.error(`Failed to fetch activity for collection ${collectionId}:`, error)
      return []
    }
  }

  async getFloorPriceHistory(collectionId: string, days: number = 30): Promise<Array<{
    timestamp: number
    floorPrice: number
    volume: number
    sales: number
  }>> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/collections/${collectionId}/history`, {
        params: { days }
      })
      return response.data
    } catch (error) {
      console.error(`Failed to fetch floor price history for ${collectionId}:`, error)
      return []
    }
  }

  async searchCollections(query: string, limit: number = 20): Promise<NFTCollection[]> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/search`, {
        params: { q: query, limit }
      })
      return response.data
    } catch (error) {
      console.error('Failed to search collections:', error)
      return []
    }
  }

  async getTopCollectionsByCategory(category: string, limit: number = 10): Promise<NFTCollection[]> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/categories/${category}`, {
        params: { limit }
      })
      return response.data
    } catch (error) {
      console.error(`Failed to fetch collections for category ${category}:`, error)
      return []
    }
  }

  async getMarketplaceStats(marketplace: NFTCollection['marketplace']): Promise<{
    totalVolume24h: number
    totalSales24h: number
    averagePrice: number
    topCollection: string
    marketShare: number
  }> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/nft/marketplaces/${marketplace}/stats`)
      return response.data
    } catch (error) {
      console.warn(`Marketplace stats for ${marketplace} unavailable`)
      return {
        totalVolume24h: 0,
        totalSales24h: 0,
        averagePrice: 0,
        topCollection: '',
        marketShare: 0
      }
    }
  }

  // Real-time price updates via WebSocket
  createCollectionWebSocket(
    collectionId: string, 
    onUpdate: (data: { floorPrice: number; volume24h: number; sales24h: number }) => void
  ): WebSocket | null {
    try {
      const ws = new WebSocket(`wss://api.magiceden.dev/ws/collections/${collectionId}`)
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'collection_update') {
            onUpdate({
              floorPrice: data.floorPrice / 1e9,
              volume24h: data.volume24h / 1e9,
              sales24h: data.sales24h
            })
          }
        } catch (error) {
          console.error('Failed to parse WebSocket NFT data:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('NFT WebSocket error:', error)
      }

      return ws
    } catch (error) {
      console.error('Failed to create NFT WebSocket:', error)
      return null
    }
  }
}
