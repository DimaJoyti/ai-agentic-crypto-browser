import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  defiAnalyticsService, 
  type ProtocolMetrics,
  type APYData,
  type YieldOpportunity,
  type ProtocolComparison,
  type DeFiTrends,
  type PortfolioAnalytics,
  ProtocolCategory,
  RiskLevel 
} from '@/lib/defi-analytics'
import { toast } from 'sonner'

export interface UseDeFiAnalyticsOptions {
  userAddress?: Address
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface DeFiAnalyticsState {
  protocolMetrics: ProtocolMetrics[]
  apyData: APYData[]
  yieldOpportunities: YieldOpportunity[]
  defiTrends: DeFiTrends | null
  portfolioAnalytics: PortfolioAnalytics | null
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useDeFiAnalytics(options: UseDeFiAnalyticsOptions = {}) {
  const {
    userAddress,
    autoRefresh = true,
    refreshInterval = 300000, // 5 minutes
    enableNotifications = true
  } = options

  const [state, setState] = useState<DeFiAnalyticsState>({
    protocolMetrics: [],
    apyData: [],
    yieldOpportunities: [],
    defiTrends: null,
    portfolioAnalytics: null,
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load analytics data
  const loadData = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const [
        protocolMetrics,
        apyData,
        yieldOpportunities,
        defiTrends,
        portfolioAnalytics
      ] = await Promise.all([
        defiAnalyticsService.getProtocolMetrics(),
        defiAnalyticsService.getAPYData(),
        defiAnalyticsService.getYieldOpportunities(),
        defiAnalyticsService.getDeFiTrends(),
        userAddress ? defiAnalyticsService.getPortfolioAnalytics(userAddress) : null
      ])

      setState(prev => ({
        ...prev,
        protocolMetrics,
        apyData,
        yieldOpportunities,
        defiTrends,
        portfolioAnalytics,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load DeFi analytics'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))

      if (enableNotifications) {
        toast.error('Analytics Error', {
          description: errorMessage
        })
      }
    }
  }, [userAddress, enableNotifications])

  // Get protocol metrics by category
  const getProtocolsByCategory = useCallback((category: ProtocolCategory) => {
    return state.protocolMetrics.filter(protocol => protocol.category === category)
  }, [state.protocolMetrics])

  // Get APY data for specific protocol
  const getProtocolAPY = useCallback((protocolId: string) => {
    return state.apyData.filter(apy => apy.protocolId === protocolId)
  }, [state.apyData])

  // Get yield opportunities by criteria
  const getFilteredOpportunities = useCallback((
    category?: ProtocolCategory,
    riskLevel?: RiskLevel,
    minAPY?: number
  ) => {
    let opportunities = state.yieldOpportunities

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
  }, [state.yieldOpportunities])

  // Get protocol comparison
  const getProtocolComparison = useCallback(async (category: ProtocolCategory) => {
    try {
      return await defiAnalyticsService.getProtocolComparison(category)
    } catch (error) {
      console.error('Failed to get protocol comparison:', error)
      return null
    }
  }, [])

  // Calculate analytics metrics
  const getAnalyticsMetrics = useCallback(() => {
    const { protocolMetrics, apyData, defiTrends } = state

    const totalTVL = protocolMetrics.reduce((sum, protocol) => {
      return sum + parseFloat(protocol.tvl)
    }, 0)

    const totalVolume24h = protocolMetrics.reduce((sum, protocol) => {
      return sum + parseFloat(protocol.volume24h)
    }, 0)

    const averageAPY = apyData.length > 0 
      ? apyData.reduce((sum, apy) => sum + parseFloat(apy.totalAPY), 0) / apyData.length
      : 0

    const highestAPY = apyData.length > 0 
      ? Math.max(...apyData.map(apy => parseFloat(apy.totalAPY)))
      : 0

    const totalUsers24h = protocolMetrics.reduce((sum, protocol) => {
      return sum + protocol.users24h
    }, 0)

    return {
      totalTVL,
      totalVolume24h,
      averageAPY,
      highestAPY,
      totalUsers24h,
      protocolCount: protocolMetrics.length,
      opportunityCount: state.yieldOpportunities.length,
      tvlTrend: defiTrends?.tvlTrend.trend || 'stable',
      apyTrend: defiTrends?.apyTrend.trend || 'stable'
    }
  }, [state])

  // Get category distribution
  const getCategoryDistribution = useCallback(() => {
    const distribution = new Map<ProtocolCategory, number>()
    
    state.protocolMetrics.forEach(protocol => {
      const current = distribution.get(protocol.category) || 0
      distribution.set(protocol.category, current + parseFloat(protocol.tvl))
    })

    const total = Array.from(distribution.values()).reduce((sum, value) => sum + value, 0)

    return Array.from(distribution.entries())
      .map(([category, tvl]) => ({ 
        category, 
        tvl, 
        percentage: total > 0 ? (tvl / total) * 100 : 0 
      }))
      .sort((a, b) => b.tvl - a.tvl)
  }, [state.protocolMetrics])

  // Get risk distribution
  const getRiskDistribution = useCallback(() => {
    const distribution = new Map<RiskLevel, number>()
    
    state.yieldOpportunities.forEach(opportunity => {
      const current = distribution.get(opportunity.riskLevel) || 0
      distribution.set(opportunity.riskLevel, current + 1)
    })

    const total = state.yieldOpportunities.length

    return Array.from(distribution.entries())
      .map(([riskLevel, count]) => ({ 
        riskLevel, 
        count, 
        percentage: total > 0 ? (count / total) * 100 : 0 
      }))
      .sort((a, b) => b.count - a.count)
  }, [state.yieldOpportunities])

  // Get top protocols by TVL
  const getTopProtocolsByTVL = useCallback((limit: number = 5) => {
    return state.protocolMetrics
      .sort((a, b) => parseFloat(b.tvl) - parseFloat(a.tvl))
      .slice(0, limit)
  }, [state.protocolMetrics])

  // Get top protocols by volume
  const getTopProtocolsByVolume = useCallback((limit: number = 5) => {
    return state.protocolMetrics
      .sort((a, b) => parseFloat(b.volume24h) - parseFloat(a.volume24h))
      .slice(0, limit)
  }, [state.protocolMetrics])

  // Get best yield opportunities
  const getBestYieldOpportunities = useCallback((limit: number = 5) => {
    return state.yieldOpportunities
      .sort((a, b) => parseFloat(b.apy) - parseFloat(a.apy))
      .slice(0, limit)
  }, [state.yieldOpportunities])

  // Get portfolio metrics
  const getPortfolioMetrics = useCallback(() => {
    if (!state.portfolioAnalytics) {
      return {
        totalValue: 0,
        totalYield24h: 0,
        averageAPY: 0,
        riskScore: 0,
        diversificationScore: 0,
        positionCount: 0
      }
    }

    const { portfolioAnalytics } = state
    return {
      totalValue: parseFloat(portfolioAnalytics.totalValue),
      totalYield24h: parseFloat(portfolioAnalytics.totalYield24h),
      averageAPY: parseFloat(portfolioAnalytics.averageAPY),
      riskScore: portfolioAnalytics.riskScore,
      diversificationScore: portfolioAnalytics.diversificationScore,
      positionCount: portfolioAnalytics.positions.length
    }
  }, [state.portfolioAnalytics])

  // Format utilities
  const formatAPY = useCallback((apy: string) => {
    return defiAnalyticsService.formatAPY(apy)
  }, [])

  const formatTVL = useCallback((tvl: string) => {
    return defiAnalyticsService.formatTVL(tvl)
  }, [])

  const formatVolume = useCallback((volume: string) => {
    return defiAnalyticsService.formatVolume(volume)
  }, [])

  const getRiskColor = useCallback((riskLevel: RiskLevel) => {
    return defiAnalyticsService.getRiskColor(riskLevel)
  }, [])

  const getCategoryIcon = useCallback((category: ProtocolCategory) => {
    return defiAnalyticsService.getCategoryIcon(category)
  }, [])

  // Get trend color
  const getTrendColor = useCallback((trend: 'up' | 'down' | 'stable') => {
    switch (trend) {
      case 'up':
        return 'text-green-600'
      case 'down':
        return 'text-red-600'
      case 'stable':
        return 'text-gray-600'
      default:
        return 'text-gray-600'
    }
  }, [])

  // Auto-refresh setup
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadData, refreshInterval)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, loadData])

  // Initial load
  useEffect(() => {
    loadData()
  }, [loadData])

  return {
    // State
    ...state,

    // Actions
    loadData,
    getProtocolComparison,

    // Getters
    getProtocolsByCategory,
    getProtocolAPY,
    getFilteredOpportunities,
    getTopProtocolsByTVL,
    getTopProtocolsByVolume,
    getBestYieldOpportunities,

    // Analytics
    getAnalyticsMetrics,
    getCategoryDistribution,
    getRiskDistribution,
    getPortfolioMetrics,

    // Utilities
    formatAPY,
    formatTVL,
    formatVolume,
    getRiskColor,
    getCategoryIcon,
    getTrendColor,

    // Computed values
    analyticsMetrics: getAnalyticsMetrics(),
    categoryDistribution: getCategoryDistribution(),
    riskDistribution: getRiskDistribution(),
    portfolioMetrics: getPortfolioMetrics(),
    topProtocolsByTVL: getTopProtocolsByTVL(),
    topProtocolsByVolume: getTopProtocolsByVolume(),
    bestYieldOpportunities: getBestYieldOpportunities(),

    // Quick access to categories
    lendingProtocols: getProtocolsByCategory(ProtocolCategory.LENDING),
    dexProtocols: getProtocolsByCategory(ProtocolCategory.DEX),
    yieldFarmingProtocols: getProtocolsByCategory(ProtocolCategory.YIELD_FARMING),
    liquidStakingProtocols: getProtocolsByCategory(ProtocolCategory.LIQUID_STAKING),

    // Quick access to risk levels
    lowRiskOpportunities: getFilteredOpportunities(undefined, RiskLevel.LOW),
    mediumRiskOpportunities: getFilteredOpportunities(undefined, RiskLevel.MEDIUM),
    highRiskOpportunities: getFilteredOpportunities(undefined, RiskLevel.HIGH),

    // Quick access to trends
    isBullishTrend: state.defiTrends?.tvlTrend.trend === 'up' && state.defiTrends?.volumeTrend.trend === 'up',
    isBearishTrend: state.defiTrends?.tvlTrend.trend === 'down' && state.defiTrends?.volumeTrend.trend === 'down'
  }
}
