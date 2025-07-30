import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Address } from 'viem'
import { 
  lendingIntegration,
  type LendingProtocol,
  type LendingAsset,
  type UserPosition,
  type LendingTransaction,
  type YieldOpportunity,
  type LiquidationOpportunity,
  type LendingConfig,
  type LendingEvent,
  type TransactionType,
  type TransactionMetadata
} from '@/lib/lending-integration'
import { toast } from 'sonner'

export interface LendingState {
  protocols: LendingProtocol[]
  assets: LendingAsset[]
  positions: UserPosition[]
  transactions: LendingTransaction[]
  yieldOpportunities: YieldOpportunity[]
  liquidationOpportunities: LiquidationOpportunity[]
  isLoading: boolean
  isExecuting: boolean
  config: LendingConfig
  error: string | null
  lastUpdate: number | null
}

export interface UseLendingOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseLendingReturn {
  // State
  state: LendingState
  
  // Position Management
  getUserPosition: (userAddress?: Address, protocolId?: string) => Promise<UserPosition[]>
  refreshPositions: () => Promise<void>
  
  // Transaction Execution
  supply: (protocolId: string, asset: LendingAsset, amount: string, useAsCollateral?: boolean) => Promise<LendingTransaction>
  withdraw: (protocolId: string, asset: LendingAsset, amount: string) => Promise<LendingTransaction>
  borrow: (protocolId: string, asset: LendingAsset, amount: string, rateMode?: 'variable' | 'stable') => Promise<LendingTransaction>
  repay: (protocolId: string, asset: LendingAsset, amount: string, rateMode?: 'variable' | 'stable') => Promise<LendingTransaction>
  
  // Yield & Opportunities
  getYieldOpportunities: (minAPY?: number, maxRisk?: number) => Promise<YieldOpportunity[]>
  getLiquidationOpportunities: (minProfit?: number) => Promise<LiquidationOpportunity[]>
  
  // Data Access
  getProtocols: (chainId?: number) => LendingProtocol[]
  getSupportedAssets: (protocolId?: string) => LendingAsset[]
  getTransaction: (id: string) => LendingTransaction | null
  
  // Configuration
  updateConfig: (config: Partial<LendingConfig>) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useLending = (
  options: UseLendingOptions = {}
): UseLendingReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefresh = true,
    refreshInterval = 30000
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<LendingState>({
    protocols: [],
    assets: [],
    positions: [],
    transactions: [],
    yieldOpportunities: [],
    liquidationOpportunities: [],
    isLoading: false,
    isExecuting: false,
    config: lendingIntegration.getConfig(),
    error: null,
    lastUpdate: null
  })

  // Update state from lending integration
  const updateState = useCallback(async () => {
    try {
      const protocols = lendingIntegration.getProtocols(chainId)
      const assets = lendingIntegration.getSupportedAssets()
      const config = lendingIntegration.getConfig()

      setState(prev => ({
        ...prev,
        protocols,
        assets,
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

  // Handle lending events
  const handleLendingEvent = useCallback((event: LendingEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'transaction_success':
          toast.success('Transaction Successful', {
            description: `${event.transaction?.type} completed successfully`
          })
          break
        case 'transaction_failed':
          toast.error('Transaction Failed', {
            description: `${event.transaction?.type} failed: ${event.error?.message}`
          })
          break
        case 'liquidation_risk':
          toast.warning('Liquidation Risk', {
            description: 'Your position is at risk of liquidation. Consider adding collateral.'
          })
          break
        case 'yield_opportunity':
          toast.info('Yield Opportunity', {
            description: `New high-yield opportunity available: ${event.opportunity?.currentAPY.toFixed(2)}% APY`
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
    const unsubscribe = lendingIntegration.addEventListener(handleLendingEvent)

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, handleLendingEvent, updateState])

  // Auto-refresh positions
  useEffect(() => {
    if (autoRefresh && address && refreshInterval > 0) {
      const interval = setInterval(() => {
        refreshPositions()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, address, refreshInterval])

  // Get user position
  const getUserPosition = useCallback(async (
    userAddress?: Address,
    protocolId?: string
  ): Promise<UserPosition[]> => {
    const targetAddress = userAddress || address
    if (!targetAddress) {
      throw new Error('User address not available')
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const positions = await lendingIntegration.getUserPosition(targetAddress, protocolId)
      
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
      await getUserPosition(address)
    }
  }, [address, getUserPosition])

  // Supply assets
  const supply = useCallback(async (
    protocolId: string,
    asset: LendingAsset,
    amount: string,
    useAsCollateral: boolean = true
  ): Promise<LendingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await lendingIntegration.executeLendingTransaction(
        protocolId,
        'supply' as TransactionType,
        asset,
        amount,
        address,
        { useAsCollateral }
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

  // Withdraw assets
  const withdraw = useCallback(async (
    protocolId: string,
    asset: LendingAsset,
    amount: string
  ): Promise<LendingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await lendingIntegration.executeLendingTransaction(
        protocolId,
        'withdraw' as TransactionType,
        asset,
        amount,
        address
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

  // Borrow assets
  const borrow = useCallback(async (
    protocolId: string,
    asset: LendingAsset,
    amount: string,
    rateMode: 'variable' | 'stable' = 'variable'
  ): Promise<LendingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await lendingIntegration.executeLendingTransaction(
        protocolId,
        'borrow' as TransactionType,
        asset,
        amount,
        address,
        { interestRateMode: rateMode }
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

  // Repay assets
  const repay = useCallback(async (
    protocolId: string,
    asset: LendingAsset,
    amount: string,
    rateMode: 'variable' | 'stable' = 'variable'
  ): Promise<LendingTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const transaction = await lendingIntegration.executeLendingTransaction(
        protocolId,
        'repay' as TransactionType,
        asset,
        amount,
        address,
        { interestRateMode: rateMode }
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

  // Get yield opportunities
  const getYieldOpportunities = useCallback(async (
    minAPY?: number,
    maxRisk?: number
  ): Promise<YieldOpportunity[]> => {
    try {
      const opportunities = await lendingIntegration.getYieldOpportunities(minAPY, maxRisk)
      
      setState(prev => ({
        ...prev,
        yieldOpportunities: opportunities
      }))

      return opportunities
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      throw error
    }
  }, [])

  // Get liquidation opportunities
  const getLiquidationOpportunities = useCallback(async (
    minProfit?: number
  ): Promise<LiquidationOpportunity[]> => {
    try {
      const opportunities = await lendingIntegration.getLiquidationOpportunities(minProfit)
      
      setState(prev => ({
        ...prev,
        liquidationOpportunities: opportunities
      }))

      return opportunities
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      throw error
    }
  }, [])

  // Get protocols
  const getProtocols = useCallback((protocolChainId?: number): LendingProtocol[] => {
    return lendingIntegration.getProtocols(protocolChainId || chainId)
  }, [chainId])

  // Get supported assets
  const getSupportedAssets = useCallback((protocolId?: string): LendingAsset[] => {
    return lendingIntegration.getSupportedAssets(protocolId)
  }, [])

  // Get transaction
  const getTransaction = useCallback((id: string): LendingTransaction | null => {
    return lendingIntegration.getTransaction(id)
  }, [])

  // Update configuration
  const updateConfig = useCallback((config: Partial<LendingConfig>) => {
    try {
      lendingIntegration.updateConfig(config)
      setState(prev => ({ ...prev, config: lendingIntegration.getConfig() }))

      if (enableNotifications) {
        toast.success('Configuration Updated', {
          description: 'Lending settings have been updated'
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
    }
  }, [updateState, address, refreshPositions])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    getUserPosition,
    refreshPositions,
    supply,
    withdraw,
    borrow,
    repay,
    getYieldOpportunities,
    getLiquidationOpportunities,
    getProtocols,
    getSupportedAssets,
    getTransaction,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for lending operations
export const useLendingPosition = (protocolId?: string) => {
  const { state, getUserPosition, supply, withdraw, borrow, repay } = useLending()
  const { address } = useAccount()

  const position = protocolId 
    ? state.positions.find(p => p.protocolId === protocolId)
    : state.positions[0]

  const refreshPosition = useCallback(async () => {
    if (address) {
      await getUserPosition(address, protocolId)
    }
  }, [address, protocolId, getUserPosition])

  return {
    position,
    supply,
    withdraw,
    borrow,
    repay,
    refreshPosition,
    isLoading: state.isLoading,
    isExecuting: state.isExecuting,
    error: state.error
  }
}

// Hook for yield farming
export const useYieldFarming = () => {
  const { getYieldOpportunities, state } = useLending()

  const findBestYield = useCallback(async (
    minAPY: number = 5,
    maxRisk: number = 50
  ) => {
    const opportunities = await getYieldOpportunities(minAPY, maxRisk)
    return opportunities.sort((a, b) => b.currentAPY - a.currentAPY)
  }, [getYieldOpportunities])

  return {
    opportunities: state.yieldOpportunities,
    findBestYield,
    isLoading: state.isLoading
  }
}

// Hook for liquidation monitoring
export const useLiquidationMonitor = () => {
  const { getLiquidationOpportunities, state } = useLending()

  const monitorPositions = useCallback(async () => {
    const atRiskPositions = state.positions.filter(p => p.isAtRisk)
    const liquidationOpps = await getLiquidationOpportunities()
    
    return {
      atRiskPositions,
      liquidationOpportunities: liquidationOpps,
      totalAtRisk: atRiskPositions.length,
      totalOpportunities: liquidationOpps.length
    }
  }, [state.positions, getLiquidationOpportunities])

  return {
    monitorPositions,
    atRiskPositions: state.positions.filter(p => p.isAtRisk),
    liquidationOpportunities: state.liquidationOpportunities,
    isLoading: state.isLoading
  }
}

// Hook for lending analytics
export const useLendingAnalytics = () => {
  const { state } = useLending()

  const analytics = {
    totalPositions: state.positions.length,
    totalSupplied: state.positions.reduce((sum, p) => sum + p.totalSupplyUSD, 0),
    totalBorrowed: state.positions.reduce((sum, p) => sum + p.totalBorrowUSD, 0),
    totalNetWorth: state.positions.reduce((sum, p) => sum + p.netWorthUSD, 0),
    averageHealthFactor: state.positions.length > 0
      ? state.positions.reduce((sum, p) => sum + p.healthFactor, 0) / state.positions.length
      : 0,
    averageAPY: state.positions.length > 0
      ? state.positions.reduce((sum, p) => sum + p.netAPY, 0) / state.positions.length
      : 0,
    riskDistribution: state.positions.reduce((acc, p) => {
      acc[p.riskLevel] = (acc[p.riskLevel] || 0) + 1
      return acc
    }, {} as Record<string, number>),
    protocolDistribution: state.positions.reduce((acc, p) => {
      acc[p.protocolId] = (acc[p.protocolId] || 0) + p.totalSupplyUSD
      return acc
    }, {} as Record<string, number>)
  }

  return analytics
}
