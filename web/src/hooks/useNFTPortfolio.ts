import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { toast } from 'sonner'

export interface NFTPortfolioState {
  portfolio: NFTPortfolio | null
  collections: NFTCollectionHolding[]
  assets: NFTAssetHolding[]
  performance: PortfolioPerformance | null
  analytics: PortfolioAnalytics | null
  isLoading: boolean
  error: string | null
  lastUpdate: number | null
}

export interface NFTPortfolio {
  id: string
  ownerAddress: Address
  name: string
  description?: string
  totalValue: number
  totalValueUSD: number
  totalAssets: number
  collections: NFTCollectionHolding[]
  assets: NFTAssetHolding[]
  performance: PortfolioPerformance
  analytics: PortfolioAnalytics
  diversification: DiversificationMetrics
  riskMetrics: PortfolioRiskMetrics
  lastUpdated: string
  createdAt: string
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

export interface AcquisitionInfo {
  date: string
  price: number
  currency: string
  currencyAddress: Address
  usdValue: number
  transactionHash: string
  marketplace: string
  method: 'purchase' | 'mint' | 'transfer' | 'airdrop' | 'unknown'
  gasUsed?: string
  gasCost?: number
}

export interface AssetValuation {
  currentValue: number
  currentValueUSD: number
  lastSalePrice?: number
  floorPrice: number
  traitFloorPrice?: number
  estimatedValue: number
  valuationMethod: 'floor' | 'last_sale' | 'trait_floor' | 'ml_estimate' | 'manual'
  confidence: number
  lastUpdated: string
}

export interface AssetPerformance {
  costBasis: number
  currentValue: number
  unrealizedPnL: number
  unrealizedPnLPercent: number
  holdingPeriod: number
  annualizedReturn: number
  bestOffer?: number
  listingPrice?: number
  priceHistory: PriceHistoryPoint[]
}

export interface PriceHistoryPoint {
  timestamp: string
  price: number
  event: 'purchase' | 'sale' | 'listing' | 'offer' | 'floor_update'
  source: string
}

export interface AssetMetadata {
  standard: 'ERC721' | 'ERC1155'
  amount: string
  frozen: boolean
  transferable: boolean
  royalties: RoyaltyInfo[]
  externalUrl?: string
  backgroundColor?: string
}

export interface RoyaltyInfo {
  recipient: Address
  percentage: number
  amount?: number
}

export interface PortfolioPerformance {
  totalCostBasis: number
  totalCurrentValue: number
  totalUnrealizedPnL: number
  totalUnrealizedPnLPercent: number
  totalRealizedPnL: number
  totalRealizedPnLPercent: number
  bestPerformer: AssetPerformanceInfo
  worstPerformer: AssetPerformanceInfo
  averageHoldingPeriod: number
  winRate: number
  profitFactor: number
  sharpeRatio: number
  maxDrawdown: number
  volatility: number
  performanceHistory: PerformanceHistoryPoint[]
}

export interface AssetPerformanceInfo {
  contractAddress: Address
  tokenId: string
  name: string
  imageUrl: string
  pnl: number
  pnlPercent: number
}

export interface PerformanceHistoryPoint {
  timestamp: string
  totalValue: number
  totalValueUSD: number
  pnl: number
  pnlPercent: number
}

export interface PortfolioAnalytics {
  topCollections: TopCollectionInfo[]
  rarityDistribution: RarityDistribution
  priceDistribution: PriceDistribution
  acquisitionAnalysis: AcquisitionAnalysis
  marketplaceDistribution: MarketplaceDistribution
  chainDistribution: ChainDistribution
  categoryDistribution: CategoryDistribution
  timeAnalysis: TimeAnalysis
}

export interface TopCollectionInfo {
  collection: CollectionInfo
  totalOwned: number
  totalValue: number
  allocation: number
  performance: number
  rank: number
}

export interface RarityDistribution {
  legendary: number
  epic: number
  rare: number
  uncommon: number
  common: number
  averageRarity: number
  rarityScore: number
}

export interface PriceDistribution {
  under1ETH: number
  between1and5ETH: number
  between5and10ETH: number
  between10and50ETH: number
  over50ETH: number
  averagePrice: number
  medianPrice: number
}

export interface AcquisitionAnalysis {
  byMethod: Record<string, number>
  byMarketplace: Record<string, number>
  byTimeframe: Record<string, number>
  averageCost: number
  totalSpent: number
  acquisitionRate: number
}

export interface MarketplaceDistribution {
  opensea: number
  looksrare: number
  x2y2: number
  blur: number
  other: number
}

export interface ChainDistribution {
  ethereum: number
  polygon: number
  arbitrum: number
  optimism: number
  other: number
}

export interface CategoryDistribution {
  art: number
  gaming: number
  collectibles: number
  utility: number
  pfp: number
  metaverse: number
  other: number
}

export interface TimeAnalysis {
  acquisitionTrend: TrendData[]
  valueTrend: TrendData[]
  performanceTrend: TrendData[]
  seasonality: SeasonalityData
}

export interface TrendData {
  period: string
  value: number
  change: number
  changePercent: number
}

export interface SeasonalityData {
  bestMonth: string
  worstMonth: string
  monthlyPerformance: Record<string, number>
  weeklyPattern: Record<string, number>
}

export interface DiversificationMetrics {
  collectionConcentration: number
  herfindahlIndex: number
  diversificationRatio: number
  correlationMatrix: CorrelationMatrix
  concentrationRisk: ConcentrationRisk
  diversificationScore: number
  recommendations: DiversificationRecommendation[]
}

export interface CorrelationMatrix {
  collections: string[]
  matrix: number[][]
  averageCorrelation: number
}

export interface ConcentrationRisk {
  topCollectionWeight: number
  top3CollectionsWeight: number
  top5CollectionsWeight: number
  riskLevel: 'low' | 'medium' | 'high' | 'critical'
  recommendations: string[]
}

export interface DiversificationRecommendation {
  type: 'reduce_concentration' | 'add_collection' | 'rebalance' | 'hedge'
  priority: 'high' | 'medium' | 'low'
  description: string
  expectedImpact: number
  implementation: string[]
}

export interface PortfolioRiskMetrics {
  overallRisk: 'low' | 'medium' | 'high' | 'critical'
  riskScore: number
  liquidityRisk: LiquidityRisk
  concentrationRisk: ConcentrationRisk
  marketRisk: MarketRisk
  volatilityMetrics: VolatilityMetrics
  valueAtRisk: ValueAtRisk
  stressTestResults: StressTestResult[]
}

export interface LiquidityRisk {
  liquidityScore: number
  averageDailyVolume: number
  liquidAssets: number
  illiquidAssets: number
  liquidationTime: number
  marketImpact: number
}

export interface MarketRisk {
  beta: number
  correlation: number
  marketExposure: number
  systematicRisk: number
  idiosyncraticRisk: number
}

export interface VolatilityMetrics {
  dailyVolatility: number
  weeklyVolatility: number
  monthlyVolatility: number
  annualizedVolatility: number
  volatilityTrend: 'increasing' | 'decreasing' | 'stable'
}

export interface ValueAtRisk {
  var95: number
  var99: number
  expectedShortfall: number
  timeHorizon: number
  confidence: number
}

export interface StressTestResult {
  scenario: string
  portfolioImpact: number
  worstCaseValue: number
  recoveryTime: number
  affectedAssets: number
}

export interface CollectionRarityMetrics {
  averageRarity: number
  rarityScore: number
  topRarityAsset: AssetRarityInfo
  rarityDistribution: Record<string, number>
}

export interface CollectionPerformance {
  totalCostBasis: number
  totalCurrentValue: number
  unrealizedPnL: number
  unrealizedPnLPercent: number
  bestAsset: AssetPerformanceInfo
  worstAsset: AssetPerformanceInfo
  floorPriceChange: number
  volumeChange: number
}

export interface UseNFTPortfolioOptions {
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseNFTPortfolioReturn {
  // State
  state: NFTPortfolioState
  
  // Portfolio Operations
  loadPortfolio: (ownerAddress: Address) => Promise<void>
  createPortfolio: (ownerAddress: Address, name: string, description?: string) => Promise<NFTPortfolio>
  updatePortfolio: () => Promise<void>
  
  // Analytics
  getPortfolioSummary: () => PortfolioSummary
  getPerformanceMetrics: () => PerformanceMetrics
  getRiskAnalysis: () => RiskAnalysis
  getDiversificationAnalysis: () => DiversificationAnalysis
  
  // Utilities
  clearError: () => void
  refresh: () => void
}

export interface PortfolioSummary {
  totalValue: number
  totalValueUSD: number
  totalAssets: number
  totalCollections: number
  totalUnrealizedPnL: number
  totalUnrealizedPnLPercent: number
  topCollection: NFTCollectionHolding | null
  bestPerformer: NFTAssetHolding | null
  worstPerformer: NFTAssetHolding | null
}

export interface PerformanceMetrics {
  totalReturn: number
  totalReturnPercent: number
  winRate: number
  sharpeRatio: number
  maxDrawdown: number
  volatility: number
  averageHoldingPeriod: number
  profitFactor: number
}

export interface RiskAnalysis {
  overallRisk: string
  riskScore: number
  liquidityRisk: number
  concentrationRisk: number
  marketRisk: number
  recommendations: string[]
}

export interface DiversificationAnalysis {
  diversificationScore: number
  concentrationRisk: number
  topCollectionWeight: number
  herfindahlIndex: number
  recommendations: DiversificationRecommendation[]
}

export const useNFTPortfolio = (
  options: UseNFTPortfolioOptions = {}
): UseNFTPortfolioReturn => {
  const {
    enableNotifications = true,
    autoRefresh = false,
    refreshInterval = 300000 // 5 minutes
  } = options

  const [state, setState] = useState<NFTPortfolioState>({
    portfolio: null,
    collections: [],
    assets: [],
    performance: null,
    analytics: null,
    isLoading: false,
    error: null,
    lastUpdate: null
  })

  // Mock portfolio data generator
  const generateMockPortfolio = useCallback(async (ownerAddress: Address): Promise<NFTPortfolio> => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1500))

    const mockAssets: NFTAssetHolding[] = [
      {
        id: 'asset_1',
        contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
        tokenId: '1',
        chainId: 1,
        name: 'Bored Ape #1',
        description: 'A unique Bored Ape NFT',
        imageUrl: '/nft/bayc-1.jpg',
        attributes: [
          { traitType: 'Background', value: 'Blue', rarity: 10, count: 1000, percentage: 10 },
          { traitType: 'Eyes', value: 'Laser Eyes', rarity: 1, count: 100, percentage: 1 }
        ],
        rarity: {
          rank: 1,
          score: 344.5,
          tier: 'Legendary',
          percentile: 99.99,
          method: 'statistical',
          rarityScore: 344.5
        },
        collection: {
          contractAddress: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
          name: 'Bored Ape Yacht Club',
          slug: 'bored-ape-yacht-club',
          imageUrl: '/nft/bayc.jpg',
          verified: true,
          floorPrice: 15.5,
          totalSupply: 10000
        },
        acquisition: {
          date: '2023-01-15T10:30:00Z',
          price: 12.5,
          currency: 'ETH',
          currencyAddress: '0x0000000000000000000000000000000000000000',
          usdValue: 22500,
          transactionHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
          marketplace: 'opensea',
          method: 'purchase',
          gasUsed: '150000',
          gasCost: 0.003
        },
        valuation: {
          currentValue: 25.0,
          currentValueUSD: 45000,
          lastSalePrice: 24.5,
          floorPrice: 15.5,
          traitFloorPrice: 35.0,
          estimatedValue: 28.0,
          valuationMethod: 'trait_floor',
          confidence: 85,
          lastUpdated: new Date().toISOString()
        },
        performance: {
          costBasis: 12.5,
          currentValue: 25.0,
          unrealizedPnL: 12.5,
          unrealizedPnLPercent: 100,
          holdingPeriod: 365,
          annualizedReturn: 100,
          priceHistory: [
            { timestamp: '2023-01-15T10:30:00Z', price: 12.5, event: 'purchase', source: 'opensea' },
            { timestamp: new Date().toISOString(), price: 25.0, event: 'floor_update', source: 'opensea' }
          ]
        },
        metadata: {
          standard: 'ERC721',
          amount: '1',
          frozen: false,
          transferable: true,
          royalties: [{ recipient: '0x0000000000000000000000000000000000000001', percentage: 2.5 }]
        },
        lastUpdated: new Date().toISOString()
      }
    ]

    // Group assets by collection
    const collectionsMap = new Map<string, NFTAssetHolding[]>()
    mockAssets.forEach(asset => {
      const key = `${asset.contractAddress}_${asset.chainId}`
      if (!collectionsMap.has(key)) {
        collectionsMap.set(key, [])
      }
      collectionsMap.get(key)!.push(asset)
    })

    const collections: NFTCollectionHolding[] = Array.from(collectionsMap.entries()).map(([key, assets]) => {
      const firstAsset = assets[0]
      const totalValue = assets.reduce((sum, asset) => sum + asset.valuation.currentValue, 0)
      const totalCostBasis = assets.reduce((sum, asset) => sum + asset.performance.costBasis, 0)

      return {
        contractAddress: firstAsset.contractAddress,
        chainId: firstAsset.chainId,
        name: firstAsset.collection.name,
        slug: firstAsset.collection.slug,
        imageUrl: firstAsset.collection.imageUrl,
        verified: firstAsset.collection.verified,
        totalOwned: assets.length,
        floorPrice: firstAsset.collection.floorPrice,
        totalValue,
        totalValueUSD: totalValue * 1800,
        averageCostBasis: totalCostBasis / assets.length,
        unrealizedPnL: totalValue - totalCostBasis,
        unrealizedPnLPercent: totalCostBasis > 0 ? ((totalValue - totalCostBasis) / totalCostBasis) * 100 : 0,
        allocation: 100, // Will be calculated
        rarity: {
          averageRarity: 75,
          rarityScore: 344.5,
          topRarityAsset: firstAsset.rarity,
          rarityDistribution: { 'Legendary': 1 }
        },
        performance: {
          totalCostBasis,
          totalCurrentValue: totalValue,
          unrealizedPnL: totalValue - totalCostBasis,
          unrealizedPnLPercent: totalCostBasis > 0 ? ((totalValue - totalCostBasis) / totalCostBasis) * 100 : 0,
          bestAsset: {
            contractAddress: firstAsset.contractAddress,
            tokenId: firstAsset.tokenId,
            name: firstAsset.name,
            imageUrl: firstAsset.imageUrl,
            pnl: firstAsset.performance.unrealizedPnL,
            pnlPercent: firstAsset.performance.unrealizedPnLPercent
          },
          worstAsset: {
            contractAddress: firstAsset.contractAddress,
            tokenId: firstAsset.tokenId,
            name: firstAsset.name,
            imageUrl: firstAsset.imageUrl,
            pnl: firstAsset.performance.unrealizedPnL,
            pnlPercent: firstAsset.performance.unrealizedPnLPercent
          },
          floorPriceChange: 5.2,
          volumeChange: 12.8
        },
        assets
      }
    })

    const totalValue = mockAssets.reduce((sum, asset) => sum + asset.valuation.currentValue, 0)
    const totalCostBasis = mockAssets.reduce((sum, asset) => sum + asset.performance.costBasis, 0)

    const portfolio: NFTPortfolio = {
      id: `portfolio_${ownerAddress.toLowerCase()}`,
      ownerAddress,
      name: 'Main Portfolio',
      description: 'Primary NFT collection',
      totalValue,
      totalValueUSD: totalValue * 1800,
      totalAssets: mockAssets.length,
      collections,
      assets: mockAssets,
      performance: {
        totalCostBasis,
        totalCurrentValue: totalValue,
        totalUnrealizedPnL: totalValue - totalCostBasis,
        totalUnrealizedPnLPercent: totalCostBasis > 0 ? ((totalValue - totalCostBasis) / totalCostBasis) * 100 : 0,
        totalRealizedPnL: 0,
        totalRealizedPnLPercent: 0,
        bestPerformer: {
          contractAddress: mockAssets[0].contractAddress,
          tokenId: mockAssets[0].tokenId,
          name: mockAssets[0].name,
          imageUrl: mockAssets[0].imageUrl,
          pnl: mockAssets[0].performance.unrealizedPnL,
          pnlPercent: mockAssets[0].performance.unrealizedPnLPercent
        },
        worstPerformer: {
          contractAddress: mockAssets[0].contractAddress,
          tokenId: mockAssets[0].tokenId,
          name: mockAssets[0].name,
          imageUrl: mockAssets[0].imageUrl,
          pnl: mockAssets[0].performance.unrealizedPnL,
          pnlPercent: mockAssets[0].performance.unrealizedPnLPercent
        },
        averageHoldingPeriod: 365,
        winRate: 75,
        profitFactor: 2.1,
        sharpeRatio: 0.85,
        maxDrawdown: 15.2,
        volatility: 42.3,
        performanceHistory: [
          {
            timestamp: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
            totalValue: totalValue * 0.9,
            totalValueUSD: totalValue * 0.9 * 1800,
            pnl: (totalValue * 0.9) - totalCostBasis,
            pnlPercent: (((totalValue * 0.9) - totalCostBasis) / totalCostBasis) * 100
          },
          {
            timestamp: new Date().toISOString(),
            totalValue,
            totalValueUSD: totalValue * 1800,
            pnl: totalValue - totalCostBasis,
            pnlPercent: ((totalValue - totalCostBasis) / totalCostBasis) * 100
          }
        ]
      },
      analytics: {
        topCollections: collections.map((collection, index) => ({
          collection: {
            contractAddress: collection.contractAddress,
            name: collection.name,
            slug: collection.slug,
            imageUrl: collection.imageUrl,
            verified: collection.verified,
            floorPrice: collection.floorPrice,
            totalSupply: 10000
          },
          totalOwned: collection.totalOwned,
          totalValue: collection.totalValue,
          allocation: (collection.totalValue / totalValue) * 100,
          performance: collection.unrealizedPnLPercent,
          rank: index + 1
        })),
        rarityDistribution: {
          legendary: 1,
          epic: 0,
          rare: 0,
          uncommon: 0,
          common: 0,
          averageRarity: 99.99,
          rarityScore: 344.5
        },
        priceDistribution: {
          under1ETH: 0,
          between1and5ETH: 0,
          between5and10ETH: 0,
          between10and50ETH: 1,
          over50ETH: 0,
          averagePrice: 25.0,
          medianPrice: 25.0
        },
        acquisitionAnalysis: {
          byMethod: { 'purchase': 1 },
          byMarketplace: { 'opensea': 1 },
          byTimeframe: { '2023': 1 },
          averageCost: 12.5,
          totalSpent: 12.5,
          acquisitionRate: 1 / 365
        },
        marketplaceDistribution: { opensea: 100, looksrare: 0, x2y2: 0, blur: 0, other: 0 },
        chainDistribution: { ethereum: 100, polygon: 0, arbitrum: 0, optimism: 0, other: 0 },
        categoryDistribution: { art: 100, gaming: 0, collectibles: 0, utility: 0, pfp: 0, metaverse: 0, other: 0 },
        timeAnalysis: {
          acquisitionTrend: [],
          valueTrend: [],
          performanceTrend: [],
          seasonality: {
            bestMonth: 'January',
            worstMonth: 'June',
            monthlyPerformance: {},
            weeklyPattern: {}
          }
        }
      },
      diversification: {
        collectionConcentration: 100,
        herfindahlIndex: 1.0,
        diversificationRatio: 1.0,
        correlationMatrix: {
          collections: ['BAYC'],
          matrix: [[1.0]],
          averageCorrelation: 1.0
        },
        concentrationRisk: {
          topCollectionWeight: 100,
          top3CollectionsWeight: 100,
          top5CollectionsWeight: 100,
          riskLevel: 'critical',
          recommendations: ['Diversify across multiple collections']
        },
        diversificationScore: 0,
        recommendations: [
          {
            type: 'add_collection',
            priority: 'high',
            description: 'Add assets from different collections to reduce concentration risk',
            expectedImpact: 25,
            implementation: ['Consider purchasing assets from Azuki or CryptoPunks']
          }
        ]
      },
      riskMetrics: {
        overallRisk: 'high',
        riskScore: 75,
        liquidityRisk: {
          liquidityScore: 85,
          averageDailyVolume: 150,
          liquidAssets: 1,
          illiquidAssets: 0,
          liquidationTime: 3,
          marketImpact: 2.5
        },
        concentrationRisk: {
          topCollectionWeight: 100,
          top3CollectionsWeight: 100,
          top5CollectionsWeight: 100,
          riskLevel: 'critical',
          recommendations: ['Diversify holdings across multiple collections']
        },
        marketRisk: {
          beta: 1.2,
          correlation: 0.85,
          marketExposure: 100,
          systematicRisk: 70,
          idiosyncraticRisk: 30
        },
        volatilityMetrics: {
          dailyVolatility: 8.5,
          weeklyVolatility: 18.2,
          monthlyVolatility: 35.7,
          annualizedVolatility: 125.3,
          volatilityTrend: 'stable'
        },
        valueAtRisk: {
          var95: totalValue * 0.15,
          var99: totalValue * 0.25,
          expectedShortfall: totalValue * 0.18,
          timeHorizon: 30,
          confidence: 95
        },
        stressTestResults: [
          {
            scenario: 'Market Crash (-50%)',
            portfolioImpact: -50,
            worstCaseValue: totalValue * 0.5,
            recoveryTime: 180,
            affectedAssets: 1
          }
        ]
      },
      lastUpdated: new Date().toISOString(),
      createdAt: new Date().toISOString()
    }

    return portfolio
  }, [])

  // Load portfolio
  const loadPortfolio = useCallback(async (ownerAddress: Address) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const portfolio = await generateMockPortfolio(ownerAddress)
      
      setState(prev => ({
        ...prev,
        portfolio,
        collections: portfolio.collections,
        assets: portfolio.assets,
        performance: portfolio.performance,
        analytics: portfolio.analytics,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      if (enableNotifications) {
        toast.success('Portfolio Loaded', {
          description: `Loaded ${portfolio.totalAssets} assets from ${portfolio.collections.length} collections`
        })
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Portfolio Load Failed', {
          description: errorMessage
        })
      }
    }
  }, [generateMockPortfolio, enableNotifications])

  // Create portfolio
  const createPortfolio = useCallback(async (
    ownerAddress: Address,
    name: string,
    description?: string
  ): Promise<NFTPortfolio> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const portfolio = await generateMockPortfolio(ownerAddress)
      portfolio.name = name
      portfolio.description = description

      setState(prev => ({
        ...prev,
        portfolio,
        collections: portfolio.collections,
        assets: portfolio.assets,
        performance: portfolio.performance,
        analytics: portfolio.analytics,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      if (enableNotifications) {
        toast.success('Portfolio Created', {
          description: `Created portfolio "${name}" with ${portfolio.totalAssets} assets`
        })
      }

      return portfolio

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Portfolio Creation Failed', {
          description: errorMessage
        })
      }

      throw error
    }
  }, [generateMockPortfolio, enableNotifications])

  // Update portfolio
  const updatePortfolio = useCallback(async () => {
    if (!state.portfolio) return

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const updatedPortfolio = await generateMockPortfolio(state.portfolio.ownerAddress)
      updatedPortfolio.name = state.portfolio.name
      updatedPortfolio.description = state.portfolio.description

      setState(prev => ({
        ...prev,
        portfolio: updatedPortfolio,
        collections: updatedPortfolio.collections,
        assets: updatedPortfolio.assets,
        performance: updatedPortfolio.performance,
        analytics: updatedPortfolio.analytics,
        isLoading: false,
        lastUpdate: Date.now()
      }))

      if (enableNotifications) {
        toast.success('Portfolio Updated', {
          description: 'Portfolio data has been refreshed'
        })
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Portfolio Update Failed', {
          description: errorMessage
        })
      }
    }
  }, [state.portfolio, generateMockPortfolio, enableNotifications])

  // Get portfolio summary
  const getPortfolioSummary = useCallback((): PortfolioSummary => {
    if (!state.portfolio) {
      return {
        totalValue: 0,
        totalValueUSD: 0,
        totalAssets: 0,
        totalCollections: 0,
        totalUnrealizedPnL: 0,
        totalUnrealizedPnLPercent: 0,
        topCollection: null,
        bestPerformer: null,
        worstPerformer: null
      }
    }

    const topCollection = state.collections.length > 0 
      ? state.collections.reduce((top, collection) => 
          collection.totalValue > top.totalValue ? collection : top
        )
      : null

    const bestPerformer = state.assets.length > 0
      ? state.assets.reduce((best, asset) =>
          asset.performance.unrealizedPnLPercent > best.performance.unrealizedPnLPercent ? asset : best
        )
      : null

    const worstPerformer = state.assets.length > 0
      ? state.assets.reduce((worst, asset) =>
          asset.performance.unrealizedPnLPercent < worst.performance.unrealizedPnLPercent ? asset : worst
        )
      : null

    return {
      totalValue: state.portfolio.totalValue,
      totalValueUSD: state.portfolio.totalValueUSD,
      totalAssets: state.portfolio.totalAssets,
      totalCollections: state.portfolio.collections.length,
      totalUnrealizedPnL: state.portfolio.performance.totalUnrealizedPnL,
      totalUnrealizedPnLPercent: state.portfolio.performance.totalUnrealizedPnLPercent,
      topCollection,
      bestPerformer,
      worstPerformer
    }
  }, [state.portfolio, state.collections, state.assets])

  // Get performance metrics
  const getPerformanceMetrics = useCallback((): PerformanceMetrics => {
    if (!state.performance) {
      return {
        totalReturn: 0,
        totalReturnPercent: 0,
        winRate: 0,
        sharpeRatio: 0,
        maxDrawdown: 0,
        volatility: 0,
        averageHoldingPeriod: 0,
        profitFactor: 0
      }
    }

    return {
      totalReturn: state.performance.totalUnrealizedPnL,
      totalReturnPercent: state.performance.totalUnrealizedPnLPercent,
      winRate: state.performance.winRate,
      sharpeRatio: state.performance.sharpeRatio,
      maxDrawdown: state.performance.maxDrawdown,
      volatility: state.performance.volatility,
      averageHoldingPeriod: state.performance.averageHoldingPeriod,
      profitFactor: state.performance.profitFactor
    }
  }, [state.performance])

  // Get risk analysis
  const getRiskAnalysis = useCallback((): RiskAnalysis => {
    if (!state.portfolio) {
      return {
        overallRisk: 'low',
        riskScore: 0,
        liquidityRisk: 0,
        concentrationRisk: 0,
        marketRisk: 0,
        recommendations: []
      }
    }

    return {
      overallRisk: state.portfolio.riskMetrics.overallRisk,
      riskScore: state.portfolio.riskMetrics.riskScore,
      liquidityRisk: state.portfolio.riskMetrics.liquidityRisk.liquidityScore,
      concentrationRisk: state.portfolio.riskMetrics.concentrationRisk.topCollectionWeight,
      marketRisk: state.portfolio.riskMetrics.marketRisk.marketExposure,
      recommendations: state.portfolio.riskMetrics.concentrationRisk.recommendations
    }
  }, [state.portfolio])

  // Get diversification analysis
  const getDiversificationAnalysis = useCallback((): DiversificationAnalysis => {
    if (!state.portfolio) {
      return {
        diversificationScore: 0,
        concentrationRisk: 0,
        topCollectionWeight: 0,
        herfindahlIndex: 0,
        recommendations: []
      }
    }

    return {
      diversificationScore: state.portfolio.diversification.diversificationScore,
      concentrationRisk: state.portfolio.diversification.collectionConcentration,
      topCollectionWeight: state.portfolio.diversification.concentrationRisk.topCollectionWeight,
      herfindahlIndex: state.portfolio.diversification.herfindahlIndex,
      recommendations: state.portfolio.diversification.recommendations
    }
  }, [state.portfolio])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Refresh
  const refresh = useCallback(() => {
    if (state.portfolio) {
      updatePortfolio()
    }
  }, [state.portfolio, updatePortfolio])

  // Auto-refresh
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0 && state.portfolio) {
      const interval = setInterval(() => {
        updatePortfolio()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, state.portfolio, updatePortfolio])

  return {
    state,
    loadPortfolio,
    createPortfolio,
    updatePortfolio,
    getPortfolioSummary,
    getPerformanceMetrics,
    getRiskAnalysis,
    getDiversificationAnalysis,
    clearError,
    refresh
  }
}
