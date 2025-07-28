import { useState, useEffect, useCallback } from 'react'
import { type Address, type Hash } from 'viem'
import { 
  defiProtocolService, 
  type DeFiProtocol, 
  type LiquidityPool, 
  type LendingPosition,
  type YieldPosition,
  ProtocolType 
} from '@/lib/defi-protocols'
import { toast } from 'sonner'

export interface UseDeFiProtocolsOptions {
  chainId: number
  userAddress?: Address
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface DeFiState {
  protocols: DeFiProtocol[]
  pools: LiquidityPool[]
  lendingPositions: LendingPosition[]
  yieldPositions: YieldPosition[]
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useDeFiProtocols(options: UseDeFiProtocolsOptions) {
  const {
    chainId,
    userAddress,
    autoRefresh = false,
    refreshInterval = 30000 // 30 seconds
  } = options

  const [state, setState] = useState<DeFiState>({
    protocols: [],
    pools: [],
    lendingPositions: [],
    yieldPositions: [],
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load protocols for current chain
  const loadProtocols = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const protocols = defiProtocolService.getProtocolsByChain(chainId)
      setState(prev => ({
        ...prev,
        protocols,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load protocols'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
    }
  }, [chainId])

  // Load user positions
  const loadUserPositions = useCallback(async () => {
    if (!userAddress) return

    setState(prev => ({ ...prev, isLoading: true }))

    try {
      // Load Aave positions
      const aavePositions = await defiProtocolService.getAavePositions(chainId, userAddress)
      
      setState(prev => ({
        ...prev,
        lendingPositions: aavePositions,
        isLoading: false,
        lastUpdated: Date.now()
      }))
    } catch (error) {
      console.warn('Failed to load user positions:', error)
      setState(prev => ({ ...prev, isLoading: false }))
    }
  }, [chainId, userAddress])

  // Get protocols by type
  const getProtocolsByType = useCallback((type: ProtocolType): DeFiProtocol[] => {
    return state.protocols.filter(protocol => protocol.type === type)
  }, [state.protocols])

  // Search protocols
  const searchProtocols = useCallback((query: string): DeFiProtocol[] => {
    return defiProtocolService.searchProtocols(query)
  }, [])

  // Get protocol statistics
  const getProtocolStats = useCallback(async (protocolId: string) => {
    try {
      return await defiProtocolService.getProtocolStats(protocolId, chainId)
    } catch (error) {
      console.error('Failed to get protocol stats:', error)
      throw error
    }
  }, [chainId])

  // Estimate token swap
  const estimateSwap = useCallback(async (
    tokenIn: Address,
    tokenOut: Address,
    amountIn: string,
    protocol: string = 'uniswap-v3'
  ) => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const estimate = await defiProtocolService.estimateSwap(
        chainId,
        tokenIn,
        tokenOut,
        amountIn,
        protocol
      )

      setState(prev => ({ ...prev, isLoading: false }))
      return estimate
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to estimate swap'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      toast.error('Swap estimation failed', {
        description: errorMessage
      })
      throw error
    }
  }, [chainId])

  // Execute token swap
  const executeSwap = useCallback(async (
    tokenIn: Address,
    tokenOut: Address,
    amountIn: string,
    minAmountOut: string,
    protocol: string = 'uniswap-v3'
  ): Promise<Hash> => {
    if (!userAddress) {
      throw new Error('User address required for swap execution')
    }

    setState(prev => ({ ...prev, isLoading: true }))

    try {
      const txHash = await defiProtocolService.executeSwap(
        chainId,
        tokenIn,
        tokenOut,
        amountIn,
        minAmountOut,
        userAddress,
        protocol
      )

      setState(prev => ({ ...prev, isLoading: false }))
      
      toast.success('Swap executed successfully', {
        description: `Transaction hash: ${txHash.slice(0, 10)}...`,
        action: {
          label: 'View',
          onClick: () => {
            // Open block explorer
            window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
          }
        }
      })

      return txHash
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to execute swap'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      
      toast.error('Swap execution failed', {
        description: errorMessage
      })
      throw error
    }
  }, [chainId, userAddress])

  // Get liquidity pools
  const getLiquidityPools = useCallback(async (
    tokenA: Address,
    tokenB: Address,
    protocol: string = 'uniswap-v3'
  ): Promise<LiquidityPool[]> => {
    setState(prev => ({ ...prev, isLoading: true }))

    try {
      let pools: LiquidityPool[] = []

      if (protocol === 'uniswap-v3') {
        pools = await defiProtocolService.getUniswapPools(chainId, tokenA, tokenB)
      }

      setState(prev => ({
        ...prev,
        pools,
        isLoading: false,
        lastUpdated: Date.now()
      }))

      return pools
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get liquidity pools'
      setState(prev => ({ ...prev, error: errorMessage, isLoading: false }))
      throw error
    }
  }, [chainId])

  // Calculate portfolio value
  const calculatePortfolioValue = useCallback((): {
    totalSupplied: number
    totalBorrowed: number
    netWorth: number
    totalYield: number
  } => {
    const totalSupplied = state.lendingPositions.reduce((sum, position) => {
      return sum + (parseFloat(position.supplied) * parseFloat(position.asset.price || '0'))
    }, 0)

    const totalBorrowed = state.lendingPositions.reduce((sum, position) => {
      return sum + (parseFloat(position.borrowed) * parseFloat(position.asset.price || '0'))
    }, 0)

    const totalYield = state.yieldPositions.reduce((sum, position) => {
      return sum + parseFloat(position.harvestable)
    }, 0)

    return {
      totalSupplied,
      totalBorrowed,
      netWorth: totalSupplied - totalBorrowed,
      totalYield
    }
  }, [state.lendingPositions, state.yieldPositions])

  // Get top protocols by TVL
  const getTopProtocols = useCallback((limit: number = 5): DeFiProtocol[] => {
    return state.protocols
      .sort((a, b) => {
        const aTvl = parseFloat(a.tvl.replace(/[$B,M]/g, ''))
        const bTvl = parseFloat(b.tvl.replace(/[$B,M]/g, ''))
        return bTvl - aTvl
      })
      .slice(0, limit)
  }, [state.protocols])

  // Get protocols by risk level
  const getProtocolsByRisk = useCallback((riskLevel: 'low' | 'medium' | 'high'): DeFiProtocol[] => {
    return state.protocols.filter(protocol => protocol.riskLevel === riskLevel)
  }, [state.protocols])

  // Refresh all data
  const refreshData = useCallback(async () => {
    loadProtocols()
    if (userAddress) {
      await loadUserPositions()
    }
  }, [loadProtocols, loadUserPositions, userAddress])

  // Auto-refresh setup
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(refreshData, refreshInterval)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, refreshData])

  // Initial load
  useEffect(() => {
    loadProtocols()
  }, [loadProtocols])

  // Load user positions when address changes
  useEffect(() => {
    if (userAddress) {
      loadUserPositions()
    }
  }, [userAddress, loadUserPositions])

  return {
    // State
    ...state,

    // Actions
    loadProtocols,
    loadUserPositions,
    refreshData,
    estimateSwap,
    executeSwap,
    getLiquidityPools,
    getProtocolStats,

    // Computed values
    getProtocolsByType,
    searchProtocols,
    calculatePortfolioValue,
    getTopProtocols,
    getProtocolsByRisk,

    // Utilities
    dexProtocols: getProtocolsByType(ProtocolType.DEX),
    lendingProtocols: getProtocolsByType(ProtocolType.LENDING),
    stakingProtocols: getProtocolsByType(ProtocolType.STAKING),
    yieldFarmingProtocols: getProtocolsByType(ProtocolType.YIELD_FARMING),
    portfolioValue: calculatePortfolioValue(),
    topProtocols: getTopProtocols(),
    lowRiskProtocols: getProtocolsByRisk('low'),
    mediumRiskProtocols: getProtocolsByRisk('medium'),
    highRiskProtocols: getProtocolsByRisk('high')
  }
}
