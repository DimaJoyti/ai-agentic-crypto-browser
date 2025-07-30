import { type Address, type Hash } from 'viem'

// Missing interface definitions
export interface PortfolioAnalytics {
  totalValue: number
  totalItems: number
  averageHoldingPeriod: number
  profitLoss: number
  profitLossPercentage: number
  topPerformers: string[]
  worstPerformers: string[]
}

export interface DiversificationMetrics {
  collectionCount: number
  averageItemsPerCollection: number
  concentrationRisk: number
  diversificationScore: number
  categoryDistribution: Record<string, number>
}

export interface PortfolioRiskMetrics {
  volatilityScore: number
  liquidityRisk: number
  concentrationRisk: number
  marketRisk: number
  overallRiskScore: number
}

export interface CollectionRarityMetrics {
  averageRarity: number
  rarityDistribution: Record<string, number>
  uniqueTraits: number
  rarityScore: number
}

export interface AcquisitionInfo {
  date: number
  price: number
  currency: string
  marketplace: string
  transactionHash: Hash
}

export interface AssetValuation {
  currentValue: number
  estimatedValue: number
  lastSalePrice: number
  floorPrice: number
  currency: string
  lastUpdated: number
}

export interface AssetPerformance {
  profitLoss: number
  profitLossPercentage: number
  holdingPeriod: number
  annualizedReturn: number
}

export interface AssetMetadata {
  name: string
  description: string
  image: string
  attributes: Array<{
    trait_type: string
    value: string | number
    rarity?: number
  }>
  externalUrl?: string
}

export interface PortfolioCollection {
  contractAddress: Address
  name: string
  symbol: string
  imageUrl?: string
  chainId?: number
  items?: PortfolioNFT[]
  totalItems: number
  totalValue: number
  totalCost?: number
  unrealizedPnL?: number
  unrealizedPnLPercentage?: number
  floorPrice: number
  averagePrice?: number
  averageBuyPrice?: number
  allocation?: number
  riskScore?: number
  performance: CollectionPerformance
  rarity: CollectionRarityMetrics
  lastUpdated: number
}

export interface PortfolioNFT {
  tokenId: string
  contractAddress: Address
  collection: string
  name: string
  image: string
  imageUrl?: string
  chainId?: number
  acquiredAt?: string
  acquiredPrice?: number
  acquiredPriceUSD?: number
  currentValue?: number
  currentValueUSD?: number
  unrealizedPnL?: number
  unrealizedPnLPercentage?: number
  rarity: number
  rarityRank: number
  acquisition: AcquisitionInfo
  valuation: AssetValuation
  performance: AssetPerformance
  metadata: AssetMetadata
  isStaked: boolean
  stakingRewards?: number
}

export interface NFTTransaction {
  id: string
  type: 'buy' | 'sell' | 'transfer' | 'mint'
  tokenId: string
  contractAddress: Address
  price?: number
  currency?: string
  from?: Address
  to?: Address
  timestamp: number
  transactionHash: Hash
}

export interface NFTPortfolioManager {
  userAddress?: Address
  portfolio?: NFTPortfolio
  portfolios?: NFTPortfolio[]
  collections: PortfolioCollection[]
  assets: PortfolioNFT[]
  analytics: PortfolioAnalytics
  diversification: DiversificationMetrics
  riskMetrics: PortfolioRiskMetrics
  totalValue?: number
  totalItems?: number
  totalCollections?: number
  performance?: any
  allocation?: any
  recommendations?: any[]
  lastUpdate?: number
  lastUpdated: number
}

export interface NFTPortfolio {
  id: string
  ownerAddress: Address
  userAddress?: Address
  name: string
  description?: string
  totalValue: number
  totalValueUSD: number
  totalAssets: number
  totalItems?: number
  collections: NFTCollectionHolding[]
  topCollections?: PortfolioCollection[]
  assets: NFTAssetHolding[]
  recentActivity?: NFTTransaction[]
  performance: PortfolioPerformance
  analytics: PortfolioAnalytics
  diversification: DiversificationMetrics
  riskMetrics: PortfolioRiskMetrics
  allocation?: any
  lastUpdated: string
  createdAt: string
  updatedAt?: string
  isDefault?: boolean
  tags?: string[]
  riskScore?: number
}

export interface NFTCollectionHolding {
  contractAddress: Address
  chainId: number
  name: string
  slug: string
  imageUrl: string
  verified: boolean
  totalOwned: number
  floorPrice: number
  totalValue: number
  totalValueUSD: number
  averageCostBasis: number
  unrealizedPnL: number
  unrealizedPnLPercent: number
  allocation: number
  rarity: CollectionRarityMetrics
  performance: CollectionPerformance
  assets: NFTAssetHolding[]
}

export interface NFTAssetHolding {
  id: string
  contractAddress: Address
  tokenId: string
  chainId: number
  name: string
  description?: string
  imageUrl: string
  animationUrl?: string
  attributes: NFTAttribute[]
  rarity: AssetRarityInfo
  collection: CollectionInfo
  acquisition: AcquisitionInfo
  valuation: AssetValuation
  performance: AssetPerformance
  metadata: AssetMetadata
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

export interface AssetRarityInfo {
  rank: number
  score: number
  tier: string
  percentile: number
  method: string
  rarityScore: number
}

export interface CollectionInfo {
  contractAddress: Address
  name: string
  slug: string
  imageUrl: string
  verified: boolean
  floorPrice: number
  totalSupply: number
}

export interface PortfolioPerformance {
  totalReturn: number
  totalReturnPercentage: number
  realizedPnL: number
  unrealizedPnL: number
  bestPerformer: PerformanceItem
  worstPerformer: PerformanceItem
  timeWeightedReturn: number
  sharpeRatio: number
  maxDrawdown: number
  winRate: number
  averageHoldingPeriod: number
  performanceHistory: PerformancePoint[]
}

export interface PerformanceItem {
  tokenId: string
  contractAddress: Address
  name: string
  return: number
  returnPercentage: number
}

export interface PerformancePoint {
  timestamp: string
  totalValue: number
  totalReturn: number
  returnPercentage: number
}

export interface CollectionPerformance {
  return: number
  returnPercentage: number
  volatility: number
  sharpeRatio: number
  maxDrawdown: number
  correlation: number
  beta: number
}

export interface AllocationAnalysis {
  byCollection: AllocationItem[]
  byCategory: AllocationItem[]
  byRarity: AllocationItem[]
  byChain: AllocationItem[]
  diversificationScore: number
  concentrationRisk: number
  recommendations: AllocationRecommendation[]
}

export interface AllocationItem {
  name: string
  value: number
  percentage: number
  targetPercentage?: number
  deviation?: number
}

export interface AllocationRecommendation {
  type: 'rebalance' | 'diversify' | 'concentrate' | 'hedge'
  description: string
  impact: number
  confidence: number
  actions: RecommendationAction[]
}

export interface RecommendationAction {
  action: 'buy' | 'sell' | 'hold'
  asset: string
  amount: number
  reason: string
}

export interface RiskMetrics {
  overallRisk: number
  liquidityRisk: number
  concentrationRisk: number
  marketRisk: number
  collectionRisks: CollectionRisk[]
  riskFactors: RiskFactor[]
  stressTest: StressTestResult[]
}

export interface CollectionRisk {
  contractAddress: Address
  name: string
  riskScore: number
  liquidityScore: number
  volatilityScore: number
  concentrationScore: number
  factors: string[]
}

export interface RiskFactor {
  factor: string
  impact: 'high' | 'medium' | 'low'
  probability: number
  mitigation: string
}

export interface StressTestResult {
  scenario: string
  portfolioImpact: number
  worstCaseValue: number
  recoveryTime: number
  mitigationStrategies: string[]
}

export interface PortfolioRecommendation {
  id: string
  type: 'buy' | 'sell' | 'hold' | 'rebalance' | 'diversify'
  priority: 'high' | 'medium' | 'low'
  title: string
  description: string
  expectedImpact: number
  confidence: number
  timeframe: string
  actions: RecommendationAction[]
  reasoning: string[]
  risks: string[]
}

export interface PortfolioOptimization {
  currentAllocation: AllocationItem[]
  optimizedAllocation: AllocationItem[]
  expectedImprovement: number
  riskReduction: number
  returnIncrease: number
  rebalanceActions: RebalanceAction[]
  optimizationMethod: 'mean_variance' | 'risk_parity' | 'black_litterman' | 'custom'
}

export interface RebalanceAction {
  action: 'buy' | 'sell'
  asset: string
  currentWeight: number
  targetWeight: number
  amount: number
  estimatedCost: number
}

export interface PortfolioFilters {
  collections?: Address[]
  chains?: number[]
  rarityTiers?: string[]
  priceRange?: { min: number; max: number }
  tags?: string[]
  isListed?: boolean
  sortBy?: 'value' | 'rarity' | 'performance' | 'acquired_date'
  sortOrder?: 'asc' | 'desc'
}

export class NFTPortfolioEngine {
  private static instance: NFTPortfolioEngine
  private portfolios = new Map<string, NFTPortfolioManager>()
  private eventListeners = new Set<(event: PortfolioEvent) => void>()

  private constructor() {
    this.initializeMockData()
  }

  static getInstance(): NFTPortfolioEngine {
    if (!NFTPortfolioEngine.instance) {
      NFTPortfolioEngine.instance = new NFTPortfolioEngine()
    }
    return NFTPortfolioEngine.instance
  }

  /**
   * Initialize mock portfolio data
   */
  private initializeMockData(): void {
    const mockPortfolio: NFTPortfolioManager = {
      userAddress: '0x0000000000000000000000000000000000000001',
      portfolios: [
        {
          id: 'default',
          name: 'Main Portfolio',
          description: 'Primary NFT collection',
          userAddress: '0x0000000000000000000000000000000000000001',
          collections: this.generateMockCollections(),
          totalValue: 125000,
          totalItems: 15,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: new Date().toISOString(),
          isDefault: true,
          tags: ['main', 'diversified'],
          performance: this.generateMockPerformance(),
          allocation: this.generateMockAllocation(),
          riskScore: 65
        }
      ],
      totalValue: 125000,
      totalItems: 15,
      totalCollections: 3,
      performance: this.generateMockPerformance(),
      allocation: this.generateMockAllocation(),
      riskMetrics: this.generateMockRiskMetrics(),
      recommendations: this.generateMockRecommendations(),
      lastUpdate: Date.now()
    }

    this.portfolios.set(mockPortfolio.portfolio.ownerAddress.toLowerCase(), mockPortfolio)
  }

  /**
   * Generate mock collections
   */
  private generateMockCollections(): PortfolioCollection[] {
    return [
      {
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        name: 'Bored Ape Yacht Club',
        symbol: 'BAYC',
        imageUrl: '/nft/bayc.jpg',
        chainId: 1,
        items: this.generateMockNFTs('BAYC', 3),
        totalItems: 3,
        totalValue: 85000,
        totalCost: 65000,
        unrealizedPnL: 20000,
        unrealizedPnLPercentage: 30.8,
        floorPrice: 15.5,
        averageBuyPrice: 12.0,
        allocation: 68,
        riskScore: 45,
        performance: {
          return: 20000,
          returnPercentage: 30.8,
          volatility: 0.35,
          sharpeRatio: 0.88,
          maxDrawdown: 0.15,
          correlation: 0.75,
          beta: 1.2
        },
        rarity: {
          averageRarity: 0.15,
          rarityDistribution: { common: 0.6, rare: 0.3, legendary: 0.1 },
          uniqueTraits: 150,
          rarityScore: 85
        },
        lastUpdated: Date.now()
      },
      {
        contractAddress: '0xED5AF388653567Af2F388E6224dC7C4b3241C544',
        name: 'Azuki',
        symbol: 'AZUKI',
        imageUrl: '/nft/azuki.jpg',
        chainId: 1,
        items: this.generateMockNFTs('AZUKI', 5),
        totalItems: 5,
        totalValue: 30000,
        totalCost: 28000,
        unrealizedPnL: 2000,
        unrealizedPnLPercentage: 7.1,
        floorPrice: 8.2,
        averageBuyPrice: 7.8,
        allocation: 24,
        riskScore: 55,
        performance: {
          return: 2000,
          returnPercentage: 7.1,
          volatility: 0.42,
          sharpeRatio: 0.17,
          maxDrawdown: 0.25,
          correlation: 0.68,
          beta: 1.1
        },
        rarity: {
          averageRarity: 0.25,
          rarityDistribution: { common: 0.5, rare: 0.35, legendary: 0.15 },
          uniqueTraits: 200,
          rarityScore: 75
        },
        lastUpdated: Date.now()
      },
      {
        contractAddress: '0x60E4d786628Fea6478F785A6d7e704777c86a7c6',
        name: 'Mutant Ape Yacht Club',
        symbol: 'MAYC',
        imageUrl: '/nft/mayc.jpg',
        chainId: 1,
        items: this.generateMockNFTs('MAYC', 7),
        totalItems: 7,
        totalValue: 10000,
        totalCost: 12000,
        unrealizedPnL: -2000,
        unrealizedPnLPercentage: -16.7,
        floorPrice: 2.8,
        averageBuyPrice: 3.2,
        allocation: 8,
        riskScore: 70,
        performance: {
          return: -2000,
          returnPercentage: -16.7,
          volatility: 0.58,
          sharpeRatio: -0.29,
          maxDrawdown: 0.35,
          correlation: 0.82,
          beta: 1.4
        },
        rarity: {
          averageRarity: 0.35,
          rarityDistribution: { common: 0.4, rare: 0.4, legendary: 0.2 },
          uniqueTraits: 120,
          rarityScore: 65
        },
        lastUpdated: Date.now()
      }
    ]
  }

  /**
   * Generate mock NFTs
   */
  private generateMockNFTs(collection: string, count: number): PortfolioNFT[] {
    const nfts: PortfolioNFT[] = []

    for (let i = 1; i <= count; i++) {
      nfts.push({
        tokenId: i.toString(),
        contractAddress: '0x0000000000000000000000000000000000000000',
        collection: collection,
        name: `${collection} #${i}`,
        image: `/nft/${collection.toLowerCase()}-${i}.jpg`,
        imageUrl: `/nft/${collection.toLowerCase()}-${i}.jpg`,
        chainId: 1,
        acquiredAt: '2023-06-15T00:00:00Z',
        acquiredPrice: 10 + Math.random() * 5,
        acquiredPriceUSD: (10 + Math.random() * 5) * 1800,
        currentValue: 12 + Math.random() * 8,
        currentValueUSD: (12 + Math.random() * 8) * 1800,
        unrealizedPnL: 2 + Math.random() * 3,
        unrealizedPnLPercentage: 15 + Math.random() * 20,
        rarity: Math.random() * 100,
        rarityRank: Math.floor(Math.random() * 10000) + 1,
        acquisition: {
          date: Date.now() - Math.random() * 365 * 24 * 60 * 60 * 1000,
          price: 10 + Math.random() * 5,
          currency: 'ETH',
          marketplace: 'OpenSea',
          transactionHash: `0x${Math.random().toString(16).substring(2, 66)}` as Hash
        },
        valuation: {
          currentValue: 12 + Math.random() * 8,
          estimatedValue: 12 + Math.random() * 8,
          lastSalePrice: 10 + Math.random() * 5,
          floorPrice: 8 + Math.random() * 3,
          currency: 'ETH',
          lastUpdated: Date.now()
        },
        performance: {
          profitLoss: 2 + Math.random() * 3,
          profitLossPercentage: 15 + Math.random() * 20,
          holdingPeriod: Math.floor(Math.random() * 365),
          annualizedReturn: 0.15 + Math.random() * 0.2
        },
        metadata: {
          name: `${collection} #${i + 1}`,
          description: `A unique ${collection} NFT`,
          image: `https://example.com/${collection.toLowerCase()}/${i + 1}.png`,
          attributes: [
            {
              trait_type: 'Background',
              value: 'Blue',
              rarity: Math.random() * 100
            }
          ]
        },
        isStaked: Math.random() > 0.8,
        stakingRewards: Math.random() > 0.8 ? Math.random() * 0.1 : undefined
      })
    }

    return nfts
  }

  /**
   * Generate mock performance
   */
  private generateMockPerformance(): PortfolioPerformance {
    return {
      totalReturn: 20000,
      totalReturnPercentage: 19.0,
      realizedPnL: 5000,
      unrealizedPnL: 15000,
      bestPerformer: {
        tokenId: '1',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        name: 'BAYC #1',
        return: 8000,
        returnPercentage: 45.2
      },
      worstPerformer: {
        tokenId: '3',
        contractAddress: '0x60E4d786628Fea6478F785A6d7e704777c86a7c6',
        name: 'MAYC #3',
        return: -1500,
        returnPercentage: -25.8
      },
      timeWeightedReturn: 18.5,
      sharpeRatio: 0.65,
      maxDrawdown: 0.22,
      winRate: 0.73,
      averageHoldingPeriod: 125,
      performanceHistory: this.generatePerformanceHistory()
    }
  }

  /**
   * Generate performance history
   */
  private generatePerformanceHistory(): PerformancePoint[] {
    const points: PerformancePoint[] = []
    const now = Date.now()
    const dayMs = 24 * 60 * 60 * 1000

    for (let i = 30; i >= 0; i--) {
      const timestamp = new Date(now - i * dayMs).toISOString()
      const baseValue = 105000
      const variation = (Math.random() - 0.5) * 0.1
      const totalValue = baseValue * (1 + variation)

      points.push({
        timestamp,
        totalValue,
        totalReturn: totalValue - 105000,
        returnPercentage: ((totalValue - 105000) / 105000) * 100
      })
    }

    return points
  }

  /**
   * Generate mock allocation
   */
  private generateMockAllocation(): AllocationAnalysis {
    return {
      byCollection: [
        { name: 'BAYC', value: 85000, percentage: 68, targetPercentage: 60, deviation: 8 },
        { name: 'Azuki', value: 30000, percentage: 24, targetPercentage: 25, deviation: -1 },
        { name: 'MAYC', value: 10000, percentage: 8, targetPercentage: 15, deviation: -7 }
      ],
      byCategory: [
        { name: 'Art', value: 95000, percentage: 76 },
        { name: 'Gaming', value: 20000, percentage: 16 },
        { name: 'Utility', value: 10000, percentage: 8 }
      ],
      byRarity: [
        { name: 'Common', value: 45000, percentage: 36 },
        { name: 'Rare', value: 50000, percentage: 40 },
        { name: 'Legendary', value: 30000, percentage: 24 }
      ],
      byChain: [
        { name: 'Ethereum', value: 125000, percentage: 100 }
      ],
      diversificationScore: 0.72,
      concentrationRisk: 0.68,
      recommendations: [
        {
          type: 'diversify',
          description: 'Consider reducing BAYC concentration',
          impact: 0.15,
          confidence: 0.8,
          actions: [
            { action: 'sell', asset: 'BAYC', amount: 1, reason: 'Reduce concentration risk' }
          ]
        }
      ]
    }
  }

  /**
   * Generate mock risk metrics
   */
  private generateMockRiskMetrics(): PortfolioRiskMetrics {
    return {
      volatilityScore: 65,
      liquidityRisk: 45,
      concentrationRisk: 75,
      marketRisk: 60,
      overallRiskScore: 65
    }
  }

  /**
   * Generate mock recommendations
   */
  private generateMockRecommendations(): PortfolioRecommendation[] {
    return [
      {
        id: 'rec_1',
        type: 'diversify',
        priority: 'high',
        title: 'Reduce BAYC Concentration',
        description: 'Your portfolio is heavily concentrated in BAYC (68%). Consider diversifying.',
        expectedImpact: 0.15,
        confidence: 0.85,
        timeframe: '1-2 weeks',
        actions: [
          { action: 'sell', asset: 'BAYC #3', amount: 1, reason: 'Reduce concentration' }
        ],
        reasoning: ['High concentration risk', 'Better risk-adjusted returns'],
        risks: ['Potential missed upside', 'Transaction costs']
      }
    ]
  }

  /**
   * Get user portfolio
   */
  async getUserPortfolio(userAddress: Address): Promise<NFTPortfolioManager | null> {
    return this.portfolios.get(userAddress.toLowerCase()) || null
  }

  /**
   * Create portfolio
   */
  async createPortfolio(
    userAddress: Address,
    name: string,
    description?: string
  ): Promise<NFTPortfolio> {
    const portfolio: NFTPortfolio = {
      id: `portfolio_${Date.now()}`,
      name,
      description,
      userAddress,
      collections: [],
      totalValue: 0,
      totalItems: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      isDefault: false,
      tags: [],
      performance: {
        totalReturn: 0,
        totalReturnPercentage: 0,
        realizedPnL: 0,
        unrealizedPnL: 0,
        bestPerformer: { tokenId: '', contractAddress: '0x0000000000000000000000000000000000000000', name: '', return: 0, returnPercentage: 0 },
        worstPerformer: { tokenId: '', contractAddress: '0x0000000000000000000000000000000000000000', name: '', return: 0, returnPercentage: 0 },
        timeWeightedReturn: 0,
        sharpeRatio: 0,
        maxDrawdown: 0,
        winRate: 0,
        averageHoldingPeriod: 0,
        performanceHistory: []
      },
      allocation: {
        byCollection: [],
        byCategory: [],
        byRarity: [],
        byChain: [],
        diversificationScore: 0,
        concentrationRisk: 0,
        recommendations: []
      },
      riskScore: 0
    }

    // Add to user's portfolio manager
    let portfolioManager = this.portfolios.get(userAddress.toLowerCase())
    if (!portfolioManager) {
      portfolioManager = {
        userAddress,
        portfolios: [],
        totalValue: 0,
        totalItems: 0,
        totalCollections: 0,
        performance: portfolio.performance,
        allocation: portfolio.allocation,
        riskMetrics: this.generateMockRiskMetrics(),
        recommendations: [],
        lastUpdate: Date.now()
      }
      this.portfolios.set(userAddress.toLowerCase(), portfolioManager)
    }

    portfolioManager.portfolios.push(portfolio)

    // Emit event
    this.emitEvent({
      type: 'portfolio_created',
      portfolio,
      timestamp: Date.now()
    })

    return portfolio
  }

  /**
   * Optimize portfolio
   */
  async optimizePortfolio(
    userAddress: Address,
    portfolioId: string,
    method: 'mean_variance' | 'risk_parity' | 'black_litterman' | 'custom' = 'mean_variance'
  ): Promise<PortfolioOptimization> {
    // Mock optimization
    const optimization: PortfolioOptimization = {
      currentAllocation: [
        { name: 'BAYC', value: 85000, percentage: 68 },
        { name: 'Azuki', value: 30000, percentage: 24 },
        { name: 'MAYC', value: 10000, percentage: 8 }
      ],
      optimizedAllocation: [
        { name: 'BAYC', value: 75000, percentage: 60 },
        { name: 'Azuki', value: 31250, percentage: 25 },
        { name: 'MAYC', value: 18750, percentage: 15 }
      ],
      expectedImprovement: 0.12,
      riskReduction: 0.08,
      returnIncrease: 0.04,
      rebalanceActions: [
        {
          action: 'sell',
          asset: 'BAYC',
          currentWeight: 68,
          targetWeight: 60,
          amount: 1,
          estimatedCost: 500
        }
      ],
      optimizationMethod: method
    }

    return optimization
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: PortfolioEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in portfolio event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: PortfolioEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
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

export interface PortfolioEvent {
  type: 'portfolio_created' | 'portfolio_updated' | 'nft_added' | 'nft_removed' | 'performance_updated' | 'optimization_completed'
  portfolio?: NFTPortfolio
  portfolioManager?: NFTPortfolioManager
  optimization?: PortfolioOptimization
  timestamp: number
}

// Export singleton instance
export const nftPortfolioEngine = NFTPortfolioEngine.getInstance()
