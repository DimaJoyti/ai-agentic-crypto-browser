import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Address } from 'viem'
import { 
  crossChainBridge,
  type BridgeProtocol,
  type BridgeRoute,
  type BridgeTransaction,
  type CrossChainPosition,
  type BridgeEvent,
  BridgeStatus
} from '@/lib/cross-chain-bridge'
import { toast } from 'sonner'

export interface CrossChainState {
  protocols: BridgeProtocol[]
  routes: BridgeRoute[]
  transactions: BridgeTransaction[]
  positions: CrossChainPosition[]
  isLoading: boolean
  isBridging: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseCrossChainOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseCrossChainReturn {
  // State
  state: CrossChainState
  
  // Bridge Operations
  getBridgeRoutes: (fromChain: number, toChain: number, tokenSymbol: string, amount: string) => Promise<BridgeRoute[]>
  executeBridge: (routeId: string, amount: string) => Promise<BridgeTransaction>
  
  // Position Management
  getCrossChainPositions: (userAddress?: Address) => Promise<CrossChainPosition[]>
  refreshPositions: () => Promise<void>
  
  // Data Access
  getProtocols: () => BridgeProtocol[]
  getTransaction: (id: string) => BridgeTransaction | null
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useCrossChain = (
  options: UseCrossChainOptions = {}
): UseCrossChainReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    autoRefresh = true,
    refreshInterval = 30000
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<CrossChainState>({
    protocols: [],
    routes: [],
    transactions: [],
    positions: [],
    isLoading: false,
    isBridging: false,
    error: null,
    lastUpdate: null
  })

  // Update state from cross-chain bridge
  const updateState = useCallback(async () => {
    try {
      const protocols = crossChainBridge.getProtocols()

      setState(prev => ({
        ...prev,
        protocols,
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
  }, [])

  // Handle bridge events
  const handleBridgeEvent = useCallback((event: BridgeEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'bridge_initiated':
          toast.success('Bridge Initiated', {
            description: `Bridge transaction started for ${event.transaction?.amount} ${event.transaction?.token.symbol}`
          })
          break
        case 'bridge_completed':
          toast.success('Bridge Completed', {
            description: `Successfully bridged ${event.transaction?.amount} ${event.transaction?.token.symbol}`
          })
          break
        case 'bridge_failed':
          toast.error('Bridge Failed', {
            description: `Bridge failed: ${event.error?.message || 'Unknown error'}`
          })
          break
        case 'position_updated':
          toast.info('Position Updated', {
            description: 'Cross-chain position has been updated'
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
    const unsubscribe = crossChainBridge.addEventListener(handleBridgeEvent)

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, handleBridgeEvent, updateState])

  // Auto-refresh data
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Get bridge routes
  const getBridgeRoutes = useCallback(async (
    fromChain: number,
    toChain: number,
    tokenSymbol: string,
    amount: string
  ): Promise<BridgeRoute[]> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const routes = await crossChainBridge.getBridgeRoutes(fromChain, toChain, tokenSymbol, amount)
      
      setState(prev => ({
        ...prev,
        isLoading: false,
        routes
      }))

      if (enableNotifications && routes.length > 0) {
        toast.success('Routes Found', {
          description: `Found ${routes.length} bridge routes`
        })
      }

      return routes
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to get routes', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications])

  // Execute bridge
  const executeBridge = useCallback(async (
    routeId: string,
    amount: string
  ): Promise<BridgeTransaction> => {
    if (!address) {
      throw new Error('Wallet not connected')
    }

    setState(prev => ({ ...prev, isBridging: true, error: null }))

    try {
      const transaction = await crossChainBridge.executeBridge(routeId, address, amount)
      
      setState(prev => ({
        ...prev,
        isBridging: false,
        transactions: [...prev.transactions, transaction]
      }))

      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isBridging: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address])

  // Get cross-chain positions
  const getCrossChainPositions = useCallback(async (
    userAddress?: Address
  ): Promise<CrossChainPosition[]> => {
    const targetAddress = userAddress || address
    if (!targetAddress) {
      throw new Error('User address not available')
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const positions = await crossChainBridge.getCrossChainPositions(targetAddress)
      
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
      await getCrossChainPositions(address)
    }
  }, [address, getCrossChainPositions])

  // Get protocols
  const getProtocols = useCallback((): BridgeProtocol[] => {
    return crossChainBridge.getProtocols()
  }, [])

  // Get transaction
  const getTransaction = useCallback((id: string): BridgeTransaction | null => {
    return crossChainBridge.getTransaction(id)
  }, [])

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
    getBridgeRoutes,
    executeBridge,
    getCrossChainPositions,
    refreshPositions,
    getProtocols,
    getTransaction,
    refresh,
    clearError
  }
}

// Simplified hook for bridge operations
export const useBridge = (fromChain?: number, toChain?: number) => {
  const { state, getBridgeRoutes, executeBridge } = useCrossChain()

  const bridge = useCallback(async (tokenSymbol: string, amount: string) => {
    if (!fromChain || !toChain) {
      throw new Error('Chains not specified')
    }

    const routes = await getBridgeRoutes(fromChain, toChain, tokenSymbol, amount)
    if (routes.length === 0) {
      throw new Error('No routes available')
    }

    const bestRoute = routes[0]
    return executeBridge(bestRoute.id, amount)
  }, [fromChain, toChain, getBridgeRoutes, executeBridge])

  const getRoutes = useCallback(async (tokenSymbol: string, amount: string) => {
    if (!fromChain || !toChain) {
      throw new Error('Chains not specified')
    }

    return getBridgeRoutes(fromChain, toChain, tokenSymbol, amount)
  }, [fromChain, toChain, getBridgeRoutes])

  return {
    bridge,
    getRoutes,
    isBridging: state.isBridging,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for cross-chain portfolio
export const useCrossChainPortfolio = () => {
  const { state, getCrossChainPositions } = useCrossChain()

  const portfolio = {
    positions: state.positions,
    totalValue: state.positions.reduce((sum, pos) => sum + pos.totalValueUSD, 0),
    totalYield: state.positions.reduce((sum, pos) => sum + pos.totalYield, 0),
    averageRisk: state.positions.length > 0
      ? state.positions.reduce((sum, pos) => sum + pos.riskScore, 0) / state.positions.length
      : 0,
    chainDistribution: state.positions.reduce((acc, pos) => {
      pos.positions.forEach(chainPos => {
        acc[chainPos.chainId] = (acc[chainPos.chainId] || 0) + chainPos.valueUSD
      })
      return acc
    }, {} as Record<number, number>),
    protocolDistribution: state.positions.reduce((acc, pos) => {
      acc[pos.protocol] = (acc[pos.protocol] || 0) + pos.totalValueUSD
      return acc
    }, {} as Record<string, number>)
  }

  return {
    portfolio,
    refreshPortfolio: getCrossChainPositions,
    isLoading: state.isLoading
  }
}

// Hook for bridge analytics
export const useBridgeAnalytics = () => {
  const { state } = useCrossChain()

  const analytics = {
    totalTransactions: state.transactions.length,
    completedTransactions: state.transactions.filter(tx => tx.status === BridgeStatus.COMPLETED).length,
    failedTransactions: state.transactions.filter(tx => tx.status === BridgeStatus.FAILED).length,
    totalVolume: state.transactions.reduce((sum, tx) => {
      if (tx.status === BridgeStatus.COMPLETED) {
        return sum + tx.amountUSD
      }
      return sum
    }, 0),
    averageFees: state.transactions.length > 0
      ? state.transactions.reduce((sum, tx) => sum + parseFloat(tx.fees.baseFee), 0) / state.transactions.length
      : 0,
    mostUsedProtocol: state.transactions.length > 0
      ? state.transactions.reduce((acc, tx) => {
          const route = state.routes.find(r => r.id === tx.routeId)
          if (route) {
            acc[route.protocol.id] = (acc[route.protocol.id] || 0) + 1
          }
          return acc
        }, {} as Record<string, number>)
      : {},
    successRate: state.transactions.length > 0
      ? (state.transactions.filter(tx => tx.status === BridgeStatus.COMPLETED).length / state.transactions.length) * 100
      : 0,
    chainPairs: state.transactions.reduce((acc, tx) => {
      const pair = `${tx.fromChain}-${tx.toChain}`
      acc[pair] = (acc[pair] || 0) + 1
      return acc
    }, {} as Record<string, number>)
  }

  return analytics
}
