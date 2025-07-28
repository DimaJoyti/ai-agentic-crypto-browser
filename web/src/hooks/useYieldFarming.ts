import { useState, useEffect, useCallback } from 'react'
import { type Address, type Hash } from 'viem'
import { 
  yieldFarmingService, 
  type YieldFarm, 
  type StakingPool,
  type UserYieldPosition,
  type YieldOptimization,
  YieldStrategy,
  RiskLevel 
} from '@/lib/yield-farming'
import { toast } from 'sonner'

export interface UseYieldFarmingOptions {
  userAddress?: Address
  chainId?: number
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface YieldFarmingState {
  farms: YieldFarm[]
  stakingPools: StakingPool[]
  userPositions: UserYieldPosition[]
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useYieldFarming(options: UseYieldFarmingOptions = {}) {
  const {
    userAddress,
    chainId,
    autoRefresh = true,
    refreshInterval = 60000, // 1 minute
    enableNotifications = true
  } = options

  const [state, setState] = useState<YieldFarmingState>({
    farms: [],
    stakingPools: [],
    userPositions: [],
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load farms and pools
  const loadData = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const farms = yieldFarmingService.getAllFarms(chainId)
      const stakingPools = yieldFarmingService.getAllStakingPools(chainId)
      const userPositions = userAddress 
        ? yieldFarmingService.getUserPositions(userAddress)
        : []

      setState(prev => ({
        ...prev,
        farms,
        stakingPools,
        userPositions,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load yield farming data'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
    }
  }, [chainId, userAddress])

  // Get farms by strategy
  const getFarmsByStrategy = useCallback((strategy: YieldStrategy): YieldFarm[] => {
    return yieldFarmingService.getFarmsByStrategy(strategy, chainId)
  }, [chainId])

  // Get farms by risk level
  const getFarmsByRisk = useCallback((riskLevel: RiskLevel): YieldFarm[] => {
    return yieldFarmingService.getFarmsByRisk(riskLevel, chainId)
  }, [chainId])

  // Get top farms by APY
  const getTopFarmsByAPY = useCallback((limit: number = 10): YieldFarm[] => {
    return yieldFarmingService.getTopFarmsByAPY(limit, chainId)
  }, [chainId])

  // Stake in farm
  const stakeFarm = useCallback(async (
    farmId: string,
    amount: string
  ): Promise<Hash> => {
    if (!userAddress) {
      throw new Error('User address required for staking')
    }

    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await yieldFarmingService.stakeFarm(farmId, amount, userAddress)
      
      // Refresh data after staking
      loadData()

      if (enableNotifications) {
        const farm = state.farms.find(f => f.id === farmId)
        toast.success('Staking successful!', {
          description: `Staked ${amount} ${farm?.stakingToken.symbol} in ${farm?.name}`,
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Staking failed'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Staking failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [userAddress, loadData, enableNotifications, state.farms])

  // Unstake from farm
  const unstakeFarm = useCallback(async (
    positionId: string,
    amount: string
  ): Promise<Hash> => {
    if (!userAddress) {
      throw new Error('User address required for unstaking')
    }

    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await yieldFarmingService.unstakeFarm(positionId, amount, userAddress)
      
      // Refresh data after unstaking
      loadData()

      if (enableNotifications) {
        toast.success('Unstaking successful!', {
          description: `Unstaked ${amount} tokens`,
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unstaking failed'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Unstaking failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [userAddress, loadData, enableNotifications])

  // Claim rewards
  const claimRewards = useCallback(async (positionId: string): Promise<Hash> => {
    if (!userAddress) {
      throw new Error('User address required for claiming rewards')
    }

    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await yieldFarmingService.claimRewards(positionId, userAddress)
      
      // Refresh data after claiming
      loadData()

      if (enableNotifications) {
        toast.success('Rewards claimed!', {
          description: 'Successfully claimed pending rewards',
          action: {
            label: 'View Transaction',
            onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        })
      }

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Claiming rewards failed'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      if (enableNotifications) {
        toast.error('Claiming failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [userAddress, loadData, enableNotifications])

  // Get yield optimization suggestions
  const getYieldOptimization = useCallback((positionId: string): YieldOptimization | null => {
    if (!userAddress) return null
    return yieldFarmingService.getYieldOptimization(positionId, userAddress)
  }, [userAddress])

  // Calculate portfolio metrics
  const getPortfolioMetrics = useCallback(() => {
    const totalStaked = state.userPositions.reduce((sum, pos) => {
      return sum + parseFloat(pos.stakedAmount)
    }, 0)

    const totalValue = state.userPositions.reduce((sum, pos) => {
      return sum + parseFloat(pos.currentValue)
    }, 0)

    const totalEarned = state.userPositions.reduce((sum, pos) => {
      return sum + parseFloat(pos.totalEarned)
    }, 0)

    const totalPendingRewards = state.userPositions.reduce((sum, pos) => {
      return sum + pos.pendingRewards.reduce((rewardSum, reward) => {
        return rewardSum + (parseFloat(reward.amount) * parseFloat(reward.price))
      }, 0)
    }, 0)

    const totalClaimableRewards = state.userPositions.reduce((sum, pos) => {
      return sum + pos.claimableRewards.reduce((rewardSum, reward) => {
        return rewardSum + (parseFloat(reward.amount) * parseFloat(reward.price))
      }, 0)
    }, 0)

    const averageAPY = state.userPositions.length > 0 
      ? state.userPositions.reduce((sum, pos) => {
          const farm = state.farms.find(f => f.id === pos.farmId)
          return sum + (farm ? parseFloat(farm.apy) : 0)
        }, 0) / state.userPositions.length
      : 0

    return {
      totalStaked,
      totalValue,
      totalEarned,
      totalPendingRewards,
      totalClaimableRewards,
      averageAPY,
      totalPnl: totalValue - totalStaked,
      totalPnlPercentage: totalStaked > 0 ? ((totalValue - totalStaked) / totalStaked) * 100 : 0,
      activePositions: state.userPositions.length,
      profitablePositions: state.userPositions.filter(pos => parseFloat(pos.pnl) > 0).length
    }
  }, [state.userPositions, state.farms])

  // Get strategy distribution
  const getStrategyDistribution = useCallback(() => {
    const distribution = Object.values(YieldStrategy).map(strategy => {
      const farms = getFarmsByStrategy(strategy)
      const userPositions = state.userPositions.filter(pos => {
        const farm = state.farms.find(f => f.id === pos.farmId)
        return farm?.strategy === strategy
      })
      
      return {
        strategy,
        farmCount: farms.length,
        userPositions: userPositions.length,
        totalStaked: userPositions.reduce((sum, pos) => sum + parseFloat(pos.stakedAmount), 0)
      }
    }).filter(item => item.farmCount > 0)

    return distribution.sort((a, b) => b.farmCount - a.farmCount)
  }, [getFarmsByStrategy, state.userPositions, state.farms])

  // Get risk distribution
  const getRiskDistribution = useCallback(() => {
    const distribution = Object.values(RiskLevel).map(riskLevel => {
      const farms = getFarmsByRisk(riskLevel)
      const userPositions = state.userPositions.filter(pos => {
        const farm = state.farms.find(f => f.id === pos.farmId)
        return farm?.riskLevel === riskLevel
      })
      
      return {
        riskLevel,
        farmCount: farms.length,
        userPositions: userPositions.length,
        totalStaked: userPositions.reduce((sum, pos) => sum + parseFloat(pos.stakedAmount), 0)
      }
    }).filter(item => item.farmCount > 0)

    return distribution.sort((a, b) => b.farmCount - a.farmCount)
  }, [getFarmsByRisk, state.userPositions, state.farms])

  // Format currency
  const formatCurrency = useCallback((amount: number): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
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
    stakeFarm,
    unstakeFarm,
    claimRewards,

    // Getters
    getFarmsByStrategy,
    getFarmsByRisk,
    getTopFarmsByAPY,
    getYieldOptimization,

    // Analytics
    getPortfolioMetrics,
    getStrategyDistribution,
    getRiskDistribution,

    // Utilities
    formatCurrency,

    // Computed values
    portfolioMetrics: getPortfolioMetrics(),
    strategyDistribution: getStrategyDistribution(),
    riskDistribution: getRiskDistribution(),
    topFarms: getTopFarmsByAPY(5),

    // Quick access to strategies
    liquidityMiningFarms: getFarmsByStrategy(YieldStrategy.LIQUIDITY_MINING),
    yieldFarmingFarms: getFarmsByStrategy(YieldStrategy.YIELD_FARMING),
    stakingFarms: getFarmsByStrategy(YieldStrategy.SINGLE_STAKING),
    lendingYieldFarms: getFarmsByStrategy(YieldStrategy.LENDING_YIELD),
    autocompoundingFarms: getFarmsByStrategy(YieldStrategy.AUTOCOMPOUNDING),
    liquidStakingFarms: getFarmsByStrategy(YieldStrategy.LIQUID_STAKING),

    // Quick access to risk levels
    lowRiskFarms: getFarmsByRisk(RiskLevel.LOW),
    mediumRiskFarms: getFarmsByRisk(RiskLevel.MEDIUM),
    highRiskFarms: getFarmsByRisk(RiskLevel.HIGH),
    veryHighRiskFarms: getFarmsByRisk(RiskLevel.VERY_HIGH)
  }
}
