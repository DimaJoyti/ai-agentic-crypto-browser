import { createPublicClient, http, type Address } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum ProtocolCategory {
  LENDING = 'lending',
  DEX = 'dex',
  YIELD_FARMING = 'yield_farming',
  LIQUID_STAKING = 'liquid_staking',
  DERIVATIVES = 'derivatives',
  INSURANCE = 'insurance',
  BRIDGE = 'bridge',
  SYNTHETIC = 'synthetic'
}

export enum RiskLevel {
  VERY_LOW = 'very_low',
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  VERY_HIGH = 'very_high'
}

export interface ProtocolMetrics {
  protocolId: string
  name: string
  category: ProtocolCategory
  chainId: number
  tvl: string
  volume24h: string
  volume7d: string
  volume30d: string
  fees24h: string
  revenue24h: string
  users24h: number
  users7d: number
  users30d: number
  marketShare: number
  dominanceIndex: number
  lastUpdated: number
}

export interface APYData {
  protocolId: string
  poolId: string
  poolName: string
  baseAPY: string
  rewardAPY: string
  totalAPY: string
  apy7d: string
  apy30d: string
  apyHistory: {
    timestamp: number
    apy: string
  }[]
  riskLevel: RiskLevel
  tvl: string
  utilization: string
  lastUpdated: number
}

export interface YieldOpportunity {
  id: string
  protocolId: string
  protocolName: string
  poolName: string
  category: ProtocolCategory
  apy: string
  tvl: string
  riskLevel: RiskLevel
  riskScore: number
  liquidityScore: number
  stabilityScore: number
  rewardTokens: {
    symbol: string
    address: Address
    allocation: number
  }[]
  requirements: {
    minDeposit: string
    lockPeriod?: number
    impermanentLoss: boolean
  }
  fees: {
    deposit: string
    withdrawal: string
    performance: string
  }
}

export interface ProtocolComparison {
  protocols: {
    protocolId: string
    name: string
    category: ProtocolCategory
    apy: string
    tvl: string
    riskLevel: RiskLevel
    fees: string
    users24h: number
    volume24h: string
  }[]
  bestAPY: string
  bestTVL: string
  bestVolume: string
  averageAPY: string
  totalTVL: string
  riskDistribution: {
    [key in RiskLevel]: number
  }
}

export interface DeFiTrends {
  tvlTrend: {
    current: string
    change24h: string
    change7d: string
    change30d: string
    trend: 'up' | 'down' | 'stable'
  }
  apyTrend: {
    average: string
    change24h: string
    change7d: string
    trend: 'up' | 'down' | 'stable'
  }
  volumeTrend: {
    total24h: string
    change24h: string
    change7d: string
    trend: 'up' | 'down' | 'stable'
  }
  userTrend: {
    active24h: number
    change24h: string
    change7d: string
    trend: 'up' | 'down' | 'stable'
  }
  topGainers: {
    protocolId: string
    name: string
    change: string
  }[]
  topLosers: {
    protocolId: string
    name: string
    change: string
  }[]
}

export interface PortfolioAnalytics {
  userAddress: Address
  totalValue: string
  totalYield24h: string
  totalYield7d: string
  totalYield30d: string
  averageAPY: string
  riskScore: number
  diversificationScore: number
  positions: {
    protocolId: string
    protocolName: string
    poolName: string
    value: string
    apy: string
    yield24h: string
    riskLevel: RiskLevel
    allocation: number
  }[]
  recommendations: {
    type: 'rebalance' | 'migrate' | 'diversify' | 'exit'
    description: string
    expectedGain: string
    riskChange: string
  }[]
}

export class DeFiAnalyticsService {
  private static instance: DeFiAnalyticsService
  private clients: Map<number, any> = new Map()
  private protocolMetrics: Map<string, ProtocolMetrics> = new Map()
  private apyData: Map<string, APYData> = new Map()
  private yieldOpportunities: Map<string, YieldOpportunity> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeMockData()
  }

  static getInstance(): DeFiAnalyticsService {
    if (!DeFiAnalyticsService.instance) {
      DeFiAnalyticsService.instance = new DeFiAnalyticsService()
    }
    return DeFiAnalyticsService.instance
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
          console.warn(`Failed to initialize DeFi analytics client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeMockData() {
    // Mock Protocol Metrics
    const uniswapMetrics: ProtocolMetrics = {
      protocolId: 'uniswap-v3',
      name: 'Uniswap V3',
      category: ProtocolCategory.DEX,
      chainId: 1,
      tvl: '4250000000',
      volume24h: '1850000000',
      volume7d: '12500000000',
      volume30d: '52000000000',
      fees24h: '5550000',
      revenue24h: '2775000',
      users24h: 125000,
      users7d: 650000,
      users30d: 2100000,
      marketShare: 65.2,
      dominanceIndex: 8.5,
      lastUpdated: Date.now()
    }

    const aaveMetrics: ProtocolMetrics = {
      protocolId: 'aave-v3',
      name: 'Aave V3',
      category: ProtocolCategory.LENDING,
      chainId: 1,
      tvl: '8500000000',
      volume24h: '450000000',
      volume7d: '3200000000',
      volume30d: '14500000000',
      fees24h: '1350000',
      revenue24h: '675000',
      users24h: 45000,
      users7d: 280000,
      users30d: 950000,
      marketShare: 42.8,
      dominanceIndex: 9.2,
      lastUpdated: Date.now()
    }

    const compoundMetrics: ProtocolMetrics = {
      protocolId: 'compound-v3',
      name: 'Compound V3',
      category: ProtocolCategory.LENDING,
      chainId: 1,
      tvl: '3200000000',
      volume24h: '180000000',
      volume7d: '1250000000',
      volume30d: '5800000000',
      fees24h: '540000',
      revenue24h: '270000',
      users24h: 28000,
      users7d: 185000,
      users30d: 620000,
      marketShare: 16.1,
      dominanceIndex: 7.8,
      lastUpdated: Date.now()
    }

    this.protocolMetrics.set('uniswap-v3', uniswapMetrics)
    this.protocolMetrics.set('aave-v3', aaveMetrics)
    this.protocolMetrics.set('compound-v3', compoundMetrics)

    // Mock APY Data
    const usdcAaveAPY: APYData = {
      protocolId: 'aave-v3',
      poolId: 'usdc-supply',
      poolName: 'USDC Supply',
      baseAPY: '3.45',
      rewardAPY: '1.25',
      totalAPY: '4.70',
      apy7d: '4.82',
      apy30d: '4.35',
      apyHistory: [
        { timestamp: Date.now() - 86400000 * 7, apy: '4.82' },
        { timestamp: Date.now() - 86400000 * 6, apy: '4.75' },
        { timestamp: Date.now() - 86400000 * 5, apy: '4.68' },
        { timestamp: Date.now() - 86400000 * 4, apy: '4.72' },
        { timestamp: Date.now() - 86400000 * 3, apy: '4.65' },
        { timestamp: Date.now() - 86400000 * 2, apy: '4.71' },
        { timestamp: Date.now() - 86400000 * 1, apy: '4.70' }
      ],
      riskLevel: RiskLevel.LOW,
      tvl: '1250000000',
      utilization: '78.5',
      lastUpdated: Date.now()
    }

    const ethUniswapAPY: APYData = {
      protocolId: 'uniswap-v3',
      poolId: 'eth-usdc-005',
      poolName: 'ETH/USDC 0.05%',
      baseAPY: '8.25',
      rewardAPY: '3.75',
      totalAPY: '12.00',
      apy7d: '11.85',
      apy30d: '12.45',
      apyHistory: [
        { timestamp: Date.now() - 86400000 * 7, apy: '11.85' },
        { timestamp: Date.now() - 86400000 * 6, apy: '12.10' },
        { timestamp: Date.now() - 86400000 * 5, apy: '11.95' },
        { timestamp: Date.now() - 86400000 * 4, apy: '12.25' },
        { timestamp: Date.now() - 86400000 * 3, apy: '12.05' },
        { timestamp: Date.now() - 86400000 * 2, apy: '11.90' },
        { timestamp: Date.now() - 86400000 * 1, apy: '12.00' }
      ],
      riskLevel: RiskLevel.MEDIUM,
      tvl: '850000000',
      utilization: '92.3',
      lastUpdated: Date.now()
    }

    this.apyData.set('aave-v3-usdc-supply', usdcAaveAPY)
    this.apyData.set('uniswap-v3-eth-usdc-005', ethUniswapAPY)

    // Mock Yield Opportunities
    const lidoStaking: YieldOpportunity = {
      id: 'lido-eth-staking',
      protocolId: 'lido',
      protocolName: 'Lido',
      poolName: 'ETH Liquid Staking',
      category: ProtocolCategory.LIQUID_STAKING,
      apy: '3.8',
      tvl: '32500000000',
      riskLevel: RiskLevel.LOW,
      riskScore: 25,
      liquidityScore: 95,
      stabilityScore: 90,
      rewardTokens: [
        {
          symbol: 'stETH',
          address: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
          allocation: 100
        }
      ],
      requirements: {
        minDeposit: '0.01',
        impermanentLoss: false
      },
      fees: {
        deposit: '0',
        withdrawal: '0',
        performance: '10'
      }
    }

    this.yieldOpportunities.set('lido-eth-staking', lidoStaking)
  }

  // Public methods
  async getProtocolMetrics(protocolId?: string): Promise<ProtocolMetrics[]> {
    if (protocolId) {
      const metrics = this.protocolMetrics.get(protocolId)
      return metrics ? [metrics] : []
    }
    return Array.from(this.protocolMetrics.values())
  }

  async getAPYData(protocolId?: string, poolId?: string): Promise<APYData[]> {
    let data = Array.from(this.apyData.values())
    
    if (protocolId) {
      data = data.filter(apy => apy.protocolId === protocolId)
    }
    
    if (poolId) {
      data = data.filter(apy => apy.poolId === poolId)
    }
    
    return data
  }

  async getYieldOpportunities(
    category?: ProtocolCategory,
    riskLevel?: RiskLevel,
    minAPY?: number
  ): Promise<YieldOpportunity[]> {
    let opportunities = Array.from(this.yieldOpportunities.values())
    
    if (category) {
      opportunities = opportunities.filter(opp => opp.category === category)
    }
    
    if (riskLevel) {
      opportunities = opportunities.filter(opp => opp.riskLevel === riskLevel)
    }
    
    if (minAPY) {
      opportunities = opportunities.filter(opp => parseFloat(opp.apy) >= minAPY)
    }
    
    return opportunities.sort((a, b) => parseFloat(b.apy) - parseFloat(a.apy))
  }

  async getProtocolComparison(category: ProtocolCategory): Promise<ProtocolComparison> {
    const protocols = Array.from(this.protocolMetrics.values())
      .filter(protocol => protocol.category === category)
      .map(protocol => {
        const apyData = Array.from(this.apyData.values())
          .find(apy => apy.protocolId === protocol.protocolId)
        
        return {
          protocolId: protocol.protocolId,
          name: protocol.name,
          category: protocol.category,
          apy: apyData?.totalAPY || '0',
          tvl: protocol.tvl,
          riskLevel: apyData?.riskLevel || RiskLevel.MEDIUM,
          fees: '0.3', // Mock fee
          users24h: protocol.users24h,
          volume24h: protocol.volume24h
        }
      })

    const apyValues = protocols.map(p => parseFloat(p.apy))
    const tvlValues = protocols.map(p => parseFloat(p.tvl))

    const riskDistribution = protocols.reduce((acc, protocol) => {
      acc[protocol.riskLevel] = (acc[protocol.riskLevel] || 0) + 1
      return acc
    }, {} as { [key in RiskLevel]: number })

    return {
      protocols,
      bestAPY: Math.max(...apyValues).toFixed(2),
      bestTVL: Math.max(...tvlValues).toLocaleString(),
      bestVolume: Math.max(...protocols.map(p => parseFloat(p.volume24h))).toLocaleString(),
      averageAPY: (apyValues.reduce((sum, apy) => sum + apy, 0) / apyValues.length).toFixed(2),
      totalTVL: tvlValues.reduce((sum, tvl) => sum + tvl, 0).toLocaleString(),
      riskDistribution
    }
  }

  async getDeFiTrends(): Promise<DeFiTrends> {
    const protocols = Array.from(this.protocolMetrics.values())
    const totalTVL = protocols.reduce((sum, p) => sum + parseFloat(p.tvl), 0)
    const totalVolume = protocols.reduce((sum, p) => sum + parseFloat(p.volume24h), 0)
    const totalUsers = protocols.reduce((sum, p) => sum + p.users24h, 0)
    
    const apyData = Array.from(this.apyData.values())
    const averageAPY = apyData.reduce((sum, apy) => sum + parseFloat(apy.totalAPY), 0) / apyData.length

    return {
      tvlTrend: {
        current: totalTVL.toLocaleString(),
        change24h: '2.5',
        change7d: '8.2',
        change30d: '15.7',
        trend: 'up'
      },
      apyTrend: {
        average: averageAPY.toFixed(2),
        change24h: '-0.15',
        change7d: '0.45',
        trend: 'stable'
      },
      volumeTrend: {
        total24h: totalVolume.toLocaleString(),
        change24h: '5.8',
        change7d: '12.3',
        trend: 'up'
      },
      userTrend: {
        active24h: totalUsers,
        change24h: '3.2',
        change7d: '7.8',
        trend: 'up'
      },
      topGainers: [
        { protocolId: 'uniswap-v3', name: 'Uniswap V3', change: '8.5' },
        { protocolId: 'aave-v3', name: 'Aave V3', change: '5.2' }
      ],
      topLosers: [
        { protocolId: 'compound-v3', name: 'Compound V3', change: '-2.1' }
      ]
    }
  }

  async getPortfolioAnalytics(userAddress: Address): Promise<PortfolioAnalytics> {
    // Mock portfolio analytics
    return {
      userAddress,
      totalValue: '125000',
      totalYield24h: '15.25',
      totalYield7d: '108.50',
      totalYield30d: '485.75',
      averageAPY: '4.85',
      riskScore: 35,
      diversificationScore: 75,
      positions: [
        {
          protocolId: 'aave-v3',
          protocolName: 'Aave V3',
          poolName: 'USDC Supply',
          value: '75000',
          apy: '4.70',
          yield24h: '9.63',
          riskLevel: RiskLevel.LOW,
          allocation: 60
        },
        {
          protocolId: 'uniswap-v3',
          protocolName: 'Uniswap V3',
          poolName: 'ETH/USDC LP',
          value: '50000',
          apy: '12.00',
          yield24h: '16.44',
          riskLevel: RiskLevel.MEDIUM,
          allocation: 40
        }
      ],
      recommendations: [
        {
          type: 'diversify',
          description: 'Consider adding liquid staking exposure for better risk-adjusted returns',
          expectedGain: '0.8%',
          riskChange: '+5'
        }
      ]
    }
  }

  // Utility methods
  formatAPY(apy: string): string {
    const value = parseFloat(apy)
    return `${value.toFixed(2)}%`
  }

  formatTVL(tvl: string): string {
    const value = parseFloat(tvl)
    if (value >= 1e9) return `$${(value / 1e9).toFixed(1)}B`
    if (value >= 1e6) return `$${(value / 1e6).toFixed(1)}M`
    if (value >= 1e3) return `$${(value / 1e3).toFixed(1)}K`
    return `$${value.toFixed(0)}`
  }

  formatVolume(volume: string): string {
    return this.formatTVL(volume)
  }

  getRiskColor(riskLevel: RiskLevel): string {
    switch (riskLevel) {
      case RiskLevel.VERY_LOW:
        return 'bg-green-100 text-green-800'
      case RiskLevel.LOW:
        return 'bg-green-100 text-green-700'
      case RiskLevel.MEDIUM:
        return 'bg-yellow-100 text-yellow-800'
      case RiskLevel.HIGH:
        return 'bg-orange-100 text-orange-800'
      case RiskLevel.VERY_HIGH:
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  getCategoryIcon(category: ProtocolCategory): string {
    switch (category) {
      case ProtocolCategory.LENDING:
        return '🏦'
      case ProtocolCategory.DEX:
        return '🔄'
      case ProtocolCategory.YIELD_FARMING:
        return '🌾'
      case ProtocolCategory.LIQUID_STAKING:
        return '💧'
      case ProtocolCategory.DERIVATIVES:
        return '📈'
      case ProtocolCategory.INSURANCE:
        return '🛡️'
      case ProtocolCategory.BRIDGE:
        return '🌉'
      case ProtocolCategory.SYNTHETIC:
        return '⚗️'
      default:
        return '🔗'
    }
  }
}

// Export singleton instance
export const defiAnalyticsService = DeFiAnalyticsService.getInstance()
