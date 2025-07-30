import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Address } from 'viem'
import {
  yieldFarmingIntegration,
  type YieldFarmingProtocol,
  type YieldFarm,
  type UserFarmPosition,
  type FarmingTransaction,
  type OptimizationResult,
  type FarmingConfig,
  type FarmingEvent,
  type FarmingTransactionType,
  type FarmingTransactionMetadata,
  type FarmType
} from '@/lib/yield-farming-integration'
import { toast } from 'sonner'

export interface UseYieldFarmingOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
  enableOptimization?: boolean
}

export interface YieldFarmingState {
  protocols: YieldFarmingProtocol[]
  farms: YieldFarm[]
  positions: UserFarmPosition[]
  transactions: FarmingTransaction[]
  optimizations: OptimizationResult[]
  isLoading: boolean
  isExecuting: boolean
  isOptimizing: boolean
  config: FarmingConfig
  error: string | null
  lastUpdate: number | null
}

export interface UseYieldFarmingReturn {
  // State
  state: YieldFarmingState

  // Farm Management
  getFarms: (protocolId?: string, minAPY?: number, maxRisk?: string, farmType?: FarmType) => Promise<YieldFarm[]>
  refreshFarms: () => Promise<void>

  // Position Management
  getUserPositions: (userAddress?: Address, protocolId?: string) => Promise<UserFarmPosition[]>
  refreshPositions: () => Promise<void>

  // Transaction Execution
  stake: (farmId: string, amount: string, metadata?: FarmingTransactionMetadata) => Promise<FarmingTransaction>
  unstake: (farmId: string, amount: string, metadata?: FarmingTransactionMetadata) => Promise<FarmingTransaction>
  harvest: (farmId: string, metadata?: FarmingTransactionMetadata) => Promise<FarmingTransaction>
  compound: (farmId: string, metadata?: FarmingTransactionMetadata) => Promise<FarmingTransaction>

  // Optimization
  getOptimizationRecommendations: (riskTolerance?: 'conservative' | 'moderate' | 'aggressive') => Promise<OptimizationResult>

  // Data Access
  getProtocols: (chainId?: number) => YieldFarmingProtocol[]
  getTransaction: (id: string) => FarmingTransaction | null

  // Configuration
  updateConfig: (config: Partial<FarmingConfig>) => void

  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useYieldFarming = (
  options: UseYieldFarmingOptions = {}
): UseYieldFarmingReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefresh = true,
    refreshInterval = 30000,
    enableOptimization = true
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<YieldFarmingState>({
    protocols: [],
    farms: [],
    positions: [],
    transactions: [],
    optimizations: [],
    isLoading: false,
    isExecuting: false,
    isOptimizing: false,
    config: yieldFarmingIntegration.getConfig(),
    error: null,
    lastUpdate: null
  })

  // Update state from yield farming integration
  const updateState = useCallback(async () => {
    try {
      const protocols = yieldFarmingIntegration.getProtocols(chainId)
      const config = yieldFarmingIntegration.getConfig()

      setState(prev => ({
        ...prev,
        protocols,
        config,
        error: null,
        lastUpdate: Date.now()
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [chainId])

  // Handle farming events
  const handleFarmingEvent = useCallback((event: FarmingEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'farming_transaction_success':
          toast.success('Transaction Successful', {
            description: `${event.transaction?.type} completed successfully`
          })
          break
        case 'farming_transaction_failed':
          toast.error('Transaction Failed', {
            description: `${event.transaction?.type} failed: ${event.error?.message}`
          })
          break
        case 'position_updated':
          toast.info('Position Updated', {
            description: 'Your farming position has been updated'
          })
          break
        case 'optimization_complete':
          toast.success('Optimization Complete', {
            description: `Found ${event.optimization?.expectedImprovement.toFixed(1)}% improvement opportunity`
          })
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = yieldFarmingIntegration.addEventListener(handleFarmingEvent)

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, handleFarmingEvent, updateState])

  // Auto-refresh data
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Get farms
  const getFarms = useCallback(async (
    protocolId?: string,
    minAPY?: number,
    maxRisk?: string,
    farmType?: FarmType
  ): Promise<YieldFarm[]> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const farms = await yieldFarmingIntegration.getYieldFarms(protocolId, minAPY, maxRisk, farmType)

      setState(prev => ({
        ...prev,
        isLoading: false,
        farms
      }))

      return farms
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get farms', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Refresh farms
  const refreshFarms = useCallback(async () => {
    await getFarms()
  }, [getFarms])

  // Get user positions
  const getUserPositions = useCallback(async (
    userAddress?: Address,
    protocolId?: string
  ): Promise<UserFarmPosition[]> => {
    const targetAddress = userAddress || address
    if (!targetAddress) {
      throw new Error('User address not available')
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const positions = await yieldFarmingIntegration.getUserPositions(targetAddress, protocolId)

      setState(prev => ({
        ...prev,
        isLoading: false,
        positions
      }))

      return positions
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get positions', { description: errorMessage })
      }
      throw error
    }
  }, [address, enableNotifications])

  // Refresh positions
  const refreshPositions = useCallback(async () => {
    if (address) {
      await getUserPositions(address)
    }
  }, [address, getUserPositions])

  // Stake tokens
  const stake = useCallback(async (
    farmId: string,
    amount: string,
    metadata?: FarmingTransactionMetadata
  ): Promise<FarmingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await yieldFarmingIntegration.executeFarmingTransaction(
        farmId,
        'stake' as FarmingTransactionType,
        amount,
        address,
        metadata
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Unstake tokens
  const unstake = useCallback(async (
    farmId: string,
    amount: string,
    metadata?: FarmingTransactionMetadata
  ): Promise<FarmingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await yieldFarmingIntegration.executeFarmingTransaction(
        farmId,
        'unstake' as FarmingTransactionType,
        amount,
        address,
        metadata
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Harvest rewards
  const harvest = useCallback(async (
    farmId: string,
    metadata?: FarmingTransactionMetadata
  ): Promise<FarmingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await yieldFarmingIntegration.executeFarmingTransaction(
        farmId,
        'claim_rewards' as FarmingTransactionType,
        '0', // No amount for harvest
        address,
        metadata
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Compound rewards
  const compound = useCallback(async (
    farmId: string,
    metadata?: FarmingTransactionMetadata
  ): Promise<FarmingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await yieldFarmingIntegration.executeFarmingTransaction(
        farmId,
        'compound' as FarmingTransactionType,
        '0', // No amount for compound
        address,
        metadata
      )

      setState(prev => ({
        ...prev,
        isExecuting: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Get optimization recommendations
  const getOptimizationRecommendations = useCallback(async (
    riskTolerance: 'conservative' | 'moderate' | 'aggressive' = 'moderate'
  ): Promise<OptimizationResult> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    if (!enableOptimization) {
      throw new Error('Optimization is disabled')
    }

    setState(prev => ({ ...prev, isOptimizing: true, error: null }))

    try {
      const optimization = await yieldFarmingIntegration.getOptimizationRecommendations(address, riskTolerance)

      setState(prev => ({
        ...prev,
        isOptimizing: false,
        optimizations: [...prev.optimizations, optimization]
      }))

      if (enableNotifications) {
        toast.success('Optimization Complete', {
          description: `Found ${optimization.expectedImprovement.toFixed(1)}% improvement opportunity`
        })
      }

      return optimization
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isOptimizing: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Optimization failed', { description: errorMessage })
      }
      throw error
    }
  }, [address, enableOptimization, enableNotifications])

  // Get protocols
  const getProtocols = useCallback((protocolChainId?: number): YieldFarmingProtocol[] => {
    return yieldFarmingIntegration.getProtocols(protocolChainId || chainId)
  }, [chainId])

  // Get transaction
  const getTransaction = useCallback((id: string): FarmingTransaction | null => {
    return yieldFarmingIntegration.getTransaction(id)
  }, [])

  // Update configuration
  const updateConfig = useCallback((config: Partial<FarmingConfig>) => {
    try {
      yieldFarmingIntegration.updateConfig(config)
      setState(prev => ({ ...prev, config: yieldFarmingIntegration.getConfig() }))

      if (enableNotifications) {
        toast.success('Configuration Updated', {
          description: 'Yield farming settings have been updated'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to update configuration', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
    if (address) {
      refreshPositions()
      refreshFarms()
    }
  }, [updateState, address, refreshPositions, refreshFarms])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    getFarms,
    refreshFarms,
    getUserPositions,
    refreshPositions,
    stake,
    unstake,
    harvest,
    compound,
    getOptimizationRecommendations,
    getProtocols,
    getTransaction,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for farming operations
export const useFarmingPosition = (farmId?: string) => {
  const { state, stake, unstake, harvest, compound } = useYieldFarming()
  const { address } = useAccount()

  const position = farmId
    ? state.positions.find(p => p.farmId === farmId)
    : state.positions[0]

  const farm = farmId
    ? state.farms.find(f => f.id === farmId)
    : null

  return {
    position,
    farm,
    stake: (amount: string, metadata?: FarmingTransactionMetadata) =>
      farmId ? stake(farmId, amount, metadata) : Promise.reject(new Error('No farm selected')),
    unstake: (amount: string, metadata?: FarmingTransactionMetadata) =>
      farmId ? unstake(farmId, amount, metadata) : Promise.reject(new Error('No farm selected')),
    harvest: (metadata?: FarmingTransactionMetadata) =>
      farmId ? harvest(farmId, metadata) : Promise.reject(new Error('No farm selected')),
    compound: (metadata?: FarmingTransactionMetadata) =>
      farmId ? compound(farmId, metadata) : Promise.reject(new Error('No farm selected')),
    isExecuting: state.isExecuting,
    error: state.error
  }
}

// Hook for yield optimization
export const useYieldOptimization = () => {
  const { getOptimizationRecommendations, state } = useYieldFarming()

  const optimize = useCallback(async (riskTolerance?: 'conservative' | 'moderate' | 'aggressive') => {
    return getOptimizationRecommendations(riskTolerance)
  }, [getOptimizationRecommendations])

  return {
    optimize,
    optimizations: state.optimizations,
    isOptimizing: state.isOptimizing,
    error: state.error
  }
}

// Hook for farming analytics
export const useFarmingAnalytics = () => {
  const { state } = useYieldFarming()

  const analytics = {
    totalPositions: state.positions.length,
    totalStaked: state.positions.reduce((sum, p) => sum + p.stakedAmountUSD, 0),
    totalRewards: state.positions.reduce((sum, p) => sum + p.totalRewardsClaimedUSD, 0),
    averageAPY: state.positions.length > 0
      ? state.positions.reduce((sum, p) => sum + p.currentAPY, 0) / state.positions.length
      : 0,
    totalROI: state.positions.reduce((sum, p) => sum + p.roi, 0),
    protocolDistribution: state.positions.reduce((acc, p) => {
      acc[p.protocolId] = (acc[p.protocolId] || 0) + p.stakedAmountUSD
      return acc
    }, {} as Record<string, number>),
    farmTypeDistribution: state.farms.reduce((acc, f) => {
      acc[f.type] = (acc[f.type] || 0) + 1
      return acc
    }, {} as Record<string, number>),
    riskDistribution: state.farms.reduce((acc, f) => {
      acc[f.riskLevel] = (acc[f.riskLevel] || 0) + 1
      return acc
    }, {} as Record<string, number>)
  }

  return analytics
}
