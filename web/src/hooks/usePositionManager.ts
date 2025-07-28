import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  positionManager, 
  type DeFiPosition, 
  type PositionAlert, 
  type PositionSummary,
  PositionType,
  PositionStatus 
} from '@/lib/position-manager'
import { toast } from 'sonner'

export interface UsePositionManagerOptions {
  userAddress?: Address
  autoRefresh?: boolean
  refreshInterval?: number
  enableNotifications?: boolean
}

export interface PositionManagerState {
  positions: DeFiPosition[]
  alerts: PositionAlert[]
  summary: PositionSummary
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function usePositionManager(options: UsePositionManagerOptions = {}) {
  const {
    userAddress,
    autoRefresh = true,
    refreshInterval = 30000, // 30 seconds
    enableNotifications = true
  } = options

  const [state, setState] = useState<PositionManagerState>({
    positions: [],
    alerts: [],
    summary: {
      totalPositions: 0,
      totalValue: '0',
      totalPnl: '0',
      totalPnlPercentage: '0',
      totalRewards: '0',
      activeAlerts: 0,
      positionsByType: {} as Record<PositionType, number>,
      positionsByProtocol: {},
      positionsByChain: {}
    },
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  // Load all positions
  const loadPositions = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const positions = positionManager.getAllPositions(userAddress)
      const alerts = positionManager.getActiveAlerts()
      const summary = positionManager.getPositionSummary(userAddress)

      setState(prev => ({
        ...prev,
        positions,
        alerts,
        summary,
        isLoading: false,
        lastUpdated: Date.now()
      }))

      // Check for new critical alerts
      if (enableNotifications) {
        const criticalAlerts = alerts.filter(alert => 
          alert.severity === 'critical' && 
          Date.now() - alert.createdAt < refreshInterval
        )

        criticalAlerts.forEach(alert => {
          toast.error(alert.title, {
            description: alert.message,
            action: {
              label: 'View Position',
              onClick: () => {
                // Navigate to position details
              }
            }
          })
        })
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to load positions'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
    }
  }, [userAddress, enableNotifications, refreshInterval])

  // Get positions by type
  const getPositionsByType = useCallback((type: PositionType): DeFiPosition[] => {
    return positionManager.getPositionsByType(type, userAddress)
  }, [userAddress])

  // Get positions by protocol
  const getPositionsByProtocol = useCallback((protocol: string): DeFiPosition[] => {
    return positionManager.getPositionsByProtocol(protocol, userAddress)
  }, [userAddress])

  // Get positions by chain
  const getPositionsByChain = useCallback((chainId: number): DeFiPosition[] => {
    return positionManager.getPositionsByChain(chainId, userAddress)
  }, [userAddress])

  // Get positions by status
  const getPositionsByStatus = useCallback((status: PositionStatus): DeFiPosition[] => {
    return state.positions.filter(position => position.status === status)
  }, [state.positions])

  // Get single position
  const getPosition = useCallback((id: string): DeFiPosition | undefined => {
    return positionManager.getPosition(id)
  }, [])

  // Acknowledge alert
  const acknowledgeAlert = useCallback((alertId: string) => {
    positionManager.acknowledgeAlert(alertId)
    loadPositions() // Refresh to update alerts
    
    if (enableNotifications) {
      toast.success('Alert acknowledged')
    }
  }, [loadPositions, enableNotifications])

  // Get alerts for specific position
  const getPositionAlerts = useCallback((positionId: string): PositionAlert[] => {
    return positionManager.getActiveAlerts(positionId)
  }, [])

  // Calculate portfolio metrics
  const getPortfolioMetrics = useCallback(() => {
    const { summary } = state
    const totalValue = parseFloat(summary.totalValue)
    const totalPnl = parseFloat(summary.totalPnl)
    const totalRewards = parseFloat(summary.totalRewards)

    return {
      totalValue,
      totalPnl,
      totalPnlPercentage: parseFloat(summary.totalPnlPercentage),
      totalRewards,
      netWorth: totalValue + totalPnl,
      profitablePositions: state.positions.filter(p => parseFloat(p.pnl) > 0).length,
      losingPositions: state.positions.filter(p => parseFloat(p.pnl) < 0).length,
      atRiskPositions: state.positions.filter(p => p.status === PositionStatus.AT_RISK).length,
      liquidatablePositions: state.positions.filter(p => p.status === PositionStatus.LIQUIDATABLE).length
    }
  }, [state])

  // Get top performing positions
  const getTopPerformers = useCallback((limit: number = 5): DeFiPosition[] => {
    return [...state.positions]
      .sort((a, b) => parseFloat(b.pnlPercentage) - parseFloat(a.pnlPercentage))
      .slice(0, limit)
  }, [state.positions])

  // Get worst performing positions
  const getWorstPerformers = useCallback((limit: number = 5): DeFiPosition[] => {
    return [...state.positions]
      .sort((a, b) => parseFloat(a.pnlPercentage) - parseFloat(b.pnlPercentage))
      .slice(0, limit)
  }, [state.positions])

  // Get positions requiring attention
  const getPositionsRequiringAttention = useCallback((): DeFiPosition[] => {
    return state.positions.filter(position => 
      position.status === PositionStatus.AT_RISK || 
      position.status === PositionStatus.LIQUIDATABLE ||
      state.alerts.some(alert => alert.positionId === position.id && alert.severity === 'high')
    )
  }, [state.positions, state.alerts])

  // Format currency values
  const formatCurrency = useCallback((amount: string | number): string => {
    const value = typeof amount === 'string' ? parseFloat(amount) : amount
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }, [])

  // Format percentage values
  const formatPercentage = useCallback((percentage: string | number): string => {
    const value = typeof percentage === 'string' ? parseFloat(percentage) : percentage
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }, [])

  // Get position type distribution
  const getPositionTypeDistribution = useCallback(() => {
    const distribution = Object.entries(state.summary.positionsByType).map(([type, count]) => ({
      type: type as PositionType,
      count,
      percentage: (count / state.summary.totalPositions) * 100
    }))
    
    return distribution.sort((a, b) => b.count - a.count)
  }, [state.summary])

  // Get protocol distribution
  const getProtocolDistribution = useCallback(() => {
    const distribution = Object.entries(state.summary.positionsByProtocol).map(([protocol, count]) => ({
      protocol,
      count,
      percentage: (count / state.summary.totalPositions) * 100
    }))
    
    return distribution.sort((a, b) => b.count - a.count)
  }, [state.summary])

  // Get chain distribution
  const getChainDistribution = useCallback(() => {
    const distribution = Object.entries(state.summary.positionsByChain).map(([chainId, count]) => ({
      chainId: parseInt(chainId),
      count,
      percentage: (count / state.summary.totalPositions) * 100
    }))
    
    return distribution.sort((a, b) => b.count - a.count)
  }, [state.summary])

  // Auto-refresh setup
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadPositions, refreshInterval)
      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, loadPositions])

  // Initial load
  useEffect(() => {
    loadPositions()
  }, [loadPositions])

  return {
    // State
    ...state,

    // Actions
    loadPositions,
    acknowledgeAlert,

    // Getters
    getPosition,
    getPositionsByType,
    getPositionsByProtocol,
    getPositionsByChain,
    getPositionsByStatus,
    getPositionAlerts,
    getTopPerformers,
    getWorstPerformers,
    getPositionsRequiringAttention,

    // Analytics
    getPortfolioMetrics,
    getPositionTypeDistribution,
    getProtocolDistribution,
    getChainDistribution,

    // Utilities
    formatCurrency,
    formatPercentage,

    // Computed values
    portfolioMetrics: getPortfolioMetrics(),
    topPerformers: getTopPerformers(),
    worstPerformers: getWorstPerformers(),
    positionsRequiringAttention: getPositionsRequiringAttention(),
    positionTypeDistribution: getPositionTypeDistribution(),
    protocolDistribution: getProtocolDistribution(),
    chainDistribution: getChainDistribution(),

    // Quick access to position types
    lendingPositions: getPositionsByType(PositionType.LENDING),
    borrowingPositions: getPositionsByType(PositionType.BORROWING),
    liquidityPositions: getPositionsByType(PositionType.LIQUIDITY),
    stakingPositions: getPositionsByType(PositionType.STAKING),
    yieldFarmingPositions: getPositionsByType(PositionType.YIELD_FARMING),

    // Quick access to position statuses
    activePositions: getPositionsByStatus(PositionStatus.ACTIVE),
    atRiskPositions: getPositionsByStatus(PositionStatus.AT_RISK),
    liquidatablePositions: getPositionsByStatus(PositionStatus.LIQUIDATABLE),

    // Alert helpers
    criticalAlerts: state.alerts.filter(a => a.severity === 'critical'),
    highAlerts: state.alerts.filter(a => a.severity === 'high'),
    mediumAlerts: state.alerts.filter(a => a.severity === 'medium'),
    lowAlerts: state.alerts.filter(a => a.severity === 'low')
  }
}
