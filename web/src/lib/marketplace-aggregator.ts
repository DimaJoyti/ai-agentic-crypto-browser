import { type Address } from 'viem'
import { 
  marketplaceAPIManager,
  type NFTCollection,
  type NFTAsset,
  type MarketplaceEvent,
  type MarketplaceStats,
  type SearchFilters
} from './marketplace-api'

export interface AggregatedMarketData {
  collections: AggregatedCollection[]
  assets: AggregatedAsset[]
  events: AggregatedEvent[]
  globalStats: GlobalMarketStats
  priceComparisons: PriceComparison[]
  marketShare: MarketShareData[]
  trends: MarketTrend[]
  lastUpdated: string
}

export interface AggregatedCollection extends NFTCollection {
  marketplaceData: MarketplaceCollectionData[]
  bestPrice: PriceInfo
  priceRange: PriceRange
  availability: AvailabilityInfo
  crossMarketplaceStats: CrossMarketplaceStats
}

export interface MarketplaceCollectionData {
  marketplace: string
  floorPrice: number
  volume24h: number
  sales24h: number
  listings: number
  lastUpdated: string
}

export interface PriceInfo {
  price: number
  marketplace: string
  currency: string
  usdValue: number
}

export interface PriceRange {
  min: number
  max: number
  average: number
  median: number
}

export interface AvailabilityInfo {
  totalListings: number
  marketplaceBreakdown: Record<string, number>
  averageListingTime: number
}

export interface CrossMarketplaceStats {
  totalVolume: number
  totalSales: number
  averagePrice: number
  priceVariance: number
  liquidityScore: number
}

export interface AggregatedAsset extends NFTAsset {
  marketplaceListings: MarketplaceListing[]
  bestListing: MarketplaceListing | null
  priceHistory: PriceHistoryPoint[]
  marketAnalysis: AssetMarketAnalysis
}

export interface MarketplaceListing {
  marketplace: string
  price: number
  currency: string
  usdValue: number
  listingType: 'fixed' | 'auction' | 'dutch_auction'
  expirationTime?: string
  seller: Address
  listingUrl: string
}

export interface PriceHistoryPoint {
  timestamp: string
  price: number
  marketplace: string
  eventType: 'sale' | 'listing' | 'offer'
}

export interface AssetMarketAnalysis {
  priceScore: number
  liquidityScore: number
  demandScore: number
  rarityPremium: number
  marketEfficiency: number
  recommendedAction: 'buy' | 'sell' | 'hold' | 'watch'
}

export interface AggregatedEvent extends MarketplaceEvent {
  crossMarketplaceData?: CrossMarketplaceEventData
  priceImpact?: PriceImpactData
  marketContext?: MarketContextData
}

export interface CrossMarketplaceEventData {
  similarEvents: MarketplaceEvent[]
  priceComparison: number[]
  volumeImpact: number
}

export interface PriceImpactData {
  priceChange: number
  volumeChange: number
  liquidityChange: number
  marketSentiment: 'bullish' | 'bearish' | 'neutral'
}

export interface MarketContextData {
  marketCondition: 'hot' | 'warm' | 'cold'
  trendDirection: 'up' | 'down' | 'sideways'
  volatility: number
  confidence: number
}

export interface GlobalMarketStats {
  totalVolume24h: number
  totalSales24h: number
  averagePrice24h: number
  totalActiveListings: number
  totalUniqueTraders24h: number
  marketCapitalization: number
  dominanceIndex: Record<string, number>
  growthMetrics: GrowthMetrics
}

export interface GrowthMetrics {
  volumeGrowth24h: number
  salesGrowth24h: number
  traderGrowth24h: number
  listingGrowth24h: number
}

export interface PriceComparison {
  collection: CollectionInfo
  marketplacePrices: Record<string, number>
  bestPrice: number
  worstPrice: number
  priceSpread: number
  arbitrageOpportunity: number
  recommendedMarketplace: string
}

export interface CollectionInfo {
  id: string
  name: string
  contractAddress: Address
  imageUrl: string
}

export interface MarketShareData {
  marketplace: string
  volumeShare: number
  salesShare: number
  listingsShare: number
  tradersShare: number
  growthRate: number
}

export interface MarketTrend {
  metric: string
  timeframe: string
  direction: 'up' | 'down' | 'stable'
  magnitude: number
  confidence: number
  dataPoints: TrendDataPoint[]
}

export interface TrendDataPoint {
  timestamp: string
  value: number
  marketplace?: string
}

export class MarketplaceDataAggregator {
  private static instance: MarketplaceDataAggregator
  private aggregatedData: AggregatedMarketData | null = null
  private lastAggregation: number = 0
  private aggregationInterval: number = 300000 // 5 minutes

  private constructor() {}

  static getInstance(): MarketplaceDataAggregator {
    if (!MarketplaceDataAggregator.instance) {
      MarketplaceDataAggregator.instance = new MarketplaceDataAggregator()
    }
    return MarketplaceDataAggregator.instance
  }

  /**
   * Aggregate data from all marketplaces
   */
  async aggregateMarketData(
    marketplaces: string[] = ['opensea', 'looksrare', 'x2y2', 'blur'],
    forceRefresh: boolean = false
  ): Promise<AggregatedMarketData> {
    const now = Date.now()
    
    // Return cached data if recent and not forcing refresh
    if (!forceRefresh && this.aggregatedData && (now - this.lastAggregation) < this.aggregationInterval) {
      return this.aggregatedData
    }

    try {
      // Fetch data from all marketplaces in parallel
      const [collectionsData, statsData, eventsData] = await Promise.all([
        this.aggregateCollections(marketplaces),
        this.aggregateStats(marketplaces),
        this.aggregateEvents(marketplaces)
      ])

      // Calculate global statistics
      const globalStats = this.calculateGlobalStats(statsData)

      // Generate price comparisons
      const priceComparisons = this.generatePriceComparisons(collectionsData)

      // Calculate market share
      const marketShare = this.calculateMarketShare(statsData)

      // Analyze trends
      const trends = this.analyzeTrends(statsData, eventsData)

      this.aggregatedData = {
        collections: collectionsData,
        assets: [], // Will be populated on demand
        events: eventsData,
        globalStats,
        priceComparisons,
        marketShare,
        trends,
        lastUpdated: new Date().toISOString()
      }

      this.lastAggregation = now

      return this.aggregatedData

    } catch (error) {
      console.error('Error aggregating market data:', error)
      throw error
    }
  }

  /**
   * Aggregate collections from multiple marketplaces
   */
  private async aggregateCollections(marketplaces: string[]): Promise<AggregatedCollection[]> {
    const allCollections = new Map<string, AggregatedCollection>()

    for (const marketplace of marketplaces) {
      try {
        const response = await marketplaceAPIManager.getCollections(marketplace, { limit: 100 })
        
        if (response.success) {
          for (const collection of response.data) {
            const key = `${collection.contractAddress}_${collection.chainId}`
            
            if (allCollections.has(key)) {
              // Merge data from multiple marketplaces
              const existing = allCollections.get(key)!
              existing.marketplaceData.push({
                marketplace,
                floorPrice: collection.floorPrice,
                volume24h: collection.volume24h,
                sales24h: collection.sales24h,
                listings: 0, // Would be fetched from marketplace
                lastUpdated: collection.lastUpdated
              })

              // Update aggregated stats
              this.updateAggregatedCollectionStats(existing)
            } else {
              // Create new aggregated collection
              const aggregated: AggregatedCollection = {
                ...collection,
                marketplaceData: [{
                  marketplace,
                  floorPrice: collection.floorPrice,
                  volume24h: collection.volume24h,
                  sales24h: collection.sales24h,
                  listings: 0,
                  lastUpdated: collection.lastUpdated
                }],
                bestPrice: {
                  price: collection.floorPrice,
                  marketplace,
                  currency: 'ETH',
                  usdValue: collection.floorPriceUSD
                },
                priceRange: {
                  min: collection.floorPrice,
                  max: collection.floorPrice,
                  average: collection.floorPrice,
                  median: collection.floorPrice
                },
                availability: {
                  totalListings: 0,
                  marketplaceBreakdown: { [marketplace]: 0 },
                  averageListingTime: 0
                },
                crossMarketplaceStats: {
                  totalVolume: collection.volume24h,
                  totalSales: collection.sales24h,
                  averagePrice: collection.averagePrice,
                  priceVariance: 0,
                  liquidityScore: 0
                }
              }

              allCollections.set(key, aggregated)
            }
          }
        }
      } catch (error) {
        console.error(`Error fetching collections from ${marketplace}:`, error)
      }
    }

    return Array.from(allCollections.values())
  }

  /**
   * Update aggregated collection statistics
   */
  private updateAggregatedCollectionStats(collection: AggregatedCollection): void {
    const marketplaceData = collection.marketplaceData

    // Find best price
    const bestPriceData = marketplaceData.reduce((best, current) => 
      current.floorPrice < best.floorPrice ? current : best
    )
    
    collection.bestPrice = {
      price: bestPriceData.floorPrice,
      marketplace: bestPriceData.marketplace,
      currency: 'ETH',
      usdValue: bestPriceData.floorPrice * 1800 // Mock ETH price
    }

    // Calculate price range
    const prices = marketplaceData.map(d => d.floorPrice)
    collection.priceRange = {
      min: Math.min(...prices),
      max: Math.max(...prices),
      average: prices.reduce((sum, price) => sum + price, 0) / prices.length,
      median: this.calculateMedian(prices)
    }

    // Update cross-marketplace stats
    collection.crossMarketplaceStats = {
      totalVolume: marketplaceData.reduce((sum, d) => sum + d.volume24h, 0),
      totalSales: marketplaceData.reduce((sum, d) => sum + d.sales24h, 0),
      averagePrice: collection.priceRange.average,
      priceVariance: this.calculateVariance(prices),
      liquidityScore: this.calculateLiquidityScore(marketplaceData)
    }
  }

  /**
   * Aggregate statistics from multiple marketplaces
   */
  private async aggregateStats(marketplaces: string[]): Promise<Record<string, MarketplaceStats>> {
    const stats: Record<string, MarketplaceStats> = {}

    for (const marketplace of marketplaces) {
      try {
        const response = await marketplaceAPIManager.getStats(marketplace)
        if (response.success) {
          stats[marketplace] = response.data
        }
      } catch (error) {
        console.error(`Error fetching stats from ${marketplace}:`, error)
      }
    }

    return stats
  }

  /**
   * Aggregate events from multiple marketplaces
   */
  private async aggregateEvents(marketplaces: string[]): Promise<AggregatedEvent[]> {
    const allEvents: AggregatedEvent[] = []

    for (const marketplace of marketplaces) {
      try {
        const response = await marketplaceAPIManager.getEvents(marketplace, { limit: 100 })
        
        if (response.success) {
          const aggregatedEvents = response.data.map(event => ({
            ...event,
            marketContext: this.analyzeMarketContext(event),
            priceImpact: this.analyzePriceImpact(event)
          }))

          allEvents.push(...aggregatedEvents)
        }
      } catch (error) {
        console.error(`Error fetching events from ${marketplace}:`, error)
      }
    }

    // Sort by timestamp (most recent first)
    return allEvents.sort((a, b) => 
      new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    )
  }

  /**
   * Calculate global market statistics
   */
  private calculateGlobalStats(statsData: Record<string, MarketplaceStats>): GlobalMarketStats {
    const stats = Object.values(statsData)

    const totalVolume24h = stats.reduce((sum, s) => sum + s.volume24h, 0)
    const totalSales24h = stats.reduce((sum, s) => sum + s.sales24h, 0)
    const totalActiveListings = stats.reduce((sum, s) => sum + s.activeListings, 0)
    const totalUniqueTraders24h = stats.reduce((sum, s) => sum + s.uniqueTraders24h, 0)

    // Calculate dominance index
    const dominanceIndex: Record<string, number> = {}
    Object.entries(statsData).forEach(([marketplace, data]) => {
      dominanceIndex[marketplace] = totalVolume24h > 0 ? (data.volume24h / totalVolume24h) * 100 : 0
    })

    return {
      totalVolume24h,
      totalSales24h,
      averagePrice24h: totalSales24h > 0 ? totalVolume24h / totalSales24h : 0,
      totalActiveListings,
      totalUniqueTraders24h,
      marketCapitalization: totalVolume24h * 365, // Rough estimate
      dominanceIndex,
      growthMetrics: {
        volumeGrowth24h: 5.2, // Mock data
        salesGrowth24h: 3.8,
        traderGrowth24h: 2.1,
        listingGrowth24h: 1.5
      }
    }
  }

  /**
   * Generate price comparisons across marketplaces
   */
  private generatePriceComparisons(collections: AggregatedCollection[]): PriceComparison[] {
    return collections
      .filter(collection => collection.marketplaceData.length > 1)
      .map(collection => {
        const marketplacePrices: Record<string, number> = {}
        collection.marketplaceData.forEach(data => {
          marketplacePrices[data.marketplace] = data.floorPrice
        })

        const prices = Object.values(marketplacePrices)
        const bestPrice = Math.min(...prices)
        const worstPrice = Math.max(...prices)
        const priceSpread = ((worstPrice - bestPrice) / bestPrice) * 100

        const bestMarketplace = Object.entries(marketplacePrices)
          .find(([, price]) => price === bestPrice)?.[0] || ''

        return {
          collection: {
            id: collection.id,
            name: collection.name,
            contractAddress: collection.contractAddress,
            imageUrl: collection.imageUrl
          },
          marketplacePrices,
          bestPrice,
          worstPrice,
          priceSpread,
          arbitrageOpportunity: priceSpread > 5 ? priceSpread : 0,
          recommendedMarketplace: bestMarketplace
        }
      })
      .filter(comparison => comparison.arbitrageOpportunity > 0)
      .sort((a, b) => b.arbitrageOpportunity - a.arbitrageOpportunity)
  }

  /**
   * Calculate market share for each marketplace
   */
  private calculateMarketShare(statsData: Record<string, MarketplaceStats>): MarketShareData[] {
    const totalVolume = Object.values(statsData).reduce((sum, s) => sum + s.volume24h, 0)
    const totalSales = Object.values(statsData).reduce((sum, s) => sum + s.sales24h, 0)
    const totalListings = Object.values(statsData).reduce((sum, s) => sum + s.activeListings, 0)
    const totalTraders = Object.values(statsData).reduce((sum, s) => sum + s.uniqueTraders24h, 0)

    return Object.entries(statsData).map(([marketplace, data]) => ({
      marketplace,
      volumeShare: totalVolume > 0 ? (data.volume24h / totalVolume) * 100 : 0,
      salesShare: totalSales > 0 ? (data.sales24h / totalSales) * 100 : 0,
      listingsShare: totalListings > 0 ? (data.activeListings / totalListings) * 100 : 0,
      tradersShare: totalTraders > 0 ? (data.uniqueTraders24h / totalTraders) * 100 : 0,
      growthRate: Math.random() * 20 - 10 // Mock growth rate
    }))
  }

  /**
   * Analyze market trends
   */
  private analyzeTrends(
    statsData: Record<string, MarketplaceStats>,
    eventsData: AggregatedEvent[]
  ): MarketTrend[] {
    const trends: MarketTrend[] = []

    // Volume trend
    const volumeData = Object.entries(statsData).map(([marketplace, data]) => ({
      timestamp: new Date().toISOString(),
      value: data.volume24h,
      marketplace
    }))

    trends.push({
      metric: 'Volume',
      timeframe: '24h',
      direction: 'up',
      magnitude: 5.2,
      confidence: 0.85,
      dataPoints: volumeData
    })

    // Sales trend
    const salesData = Object.entries(statsData).map(([marketplace, data]) => ({
      timestamp: new Date().toISOString(),
      value: data.sales24h,
      marketplace
    }))

    trends.push({
      metric: 'Sales',
      timeframe: '24h',
      direction: 'up',
      magnitude: 3.8,
      confidence: 0.78,
      dataPoints: salesData
    })

    return trends
  }

  /**
   * Utility methods
   */
  private calculateMedian(numbers: number[]): number {
    const sorted = [...numbers].sort((a, b) => a - b)
    const mid = Math.floor(sorted.length / 2)
    return sorted.length % 2 === 0 
      ? (sorted[mid - 1] + sorted[mid]) / 2 
      : sorted[mid]
  }

  private calculateVariance(numbers: number[]): number {
    const mean = numbers.reduce((sum, num) => sum + num, 0) / numbers.length
    const squaredDiffs = numbers.map(num => Math.pow(num - mean, 2))
    return squaredDiffs.reduce((sum, diff) => sum + diff, 0) / numbers.length
  }

  private calculateLiquidityScore(marketplaceData: MarketplaceCollectionData[]): number {
    // Simple liquidity score based on volume and sales
    const totalVolume = marketplaceData.reduce((sum, d) => sum + d.volume24h, 0)
    const totalSales = marketplaceData.reduce((sum, d) => sum + d.sales24h, 0)
    
    return Math.min(100, (totalVolume / 1000) + (totalSales * 2))
  }

  private analyzeMarketContext(event: MarketplaceEvent): MarketContextData {
    // Mock market context analysis
    return {
      marketCondition: 'warm',
      trendDirection: 'up',
      volatility: Math.random() * 0.5,
      confidence: 0.75
    }
  }

  private analyzePriceImpact(event: MarketplaceEvent): PriceImpactData {
    // Mock price impact analysis
    return {
      priceChange: (Math.random() - 0.5) * 10,
      volumeChange: (Math.random() - 0.5) * 20,
      liquidityChange: (Math.random() - 0.5) * 15,
      marketSentiment: Math.random() > 0.5 ? 'bullish' : 'bearish'
    }
  }

  /**
   * Get cached aggregated data
   */
  getCachedData(): AggregatedMarketData | null {
    return this.aggregatedData
  }

  /**
   * Clear cached data
   */
  clearCache(): void {
    this.aggregatedData = null
    this.lastAggregation = 0
  }
}

// Export singleton instance
export const marketplaceDataAggregator = MarketplaceDataAggregator.getInstance()
