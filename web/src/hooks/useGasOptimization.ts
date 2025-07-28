import { useState, useEffect, useCallback, useRef } from 'react'
import { 
  gasOptimizer, 
  type GasEstimate, 
  type GasTracker, 
  type GasOptimizationSuggestion,
  GasPriority 
} from '@/lib/gas-optimization'
import { toast } from 'sonner'

export interface UseGasOptimizationOptions {
  chainId: number
  autoUpdate?: boolean
  updateInterval?: number // milliseconds
  enableNotifications?: boolean
}

export interface GasOptimizationState {
  gasEstimates: GasEstimate[]
  gasTracker: GasTracker | null
  suggestions: GasOptimizationSuggestion[]
  isLoading: boolean
  error: string | null
  lastUpdated: number
}

export function useGasOptimization(options: UseGasOptimizationOptions) {
  const {
    chainId,
    autoUpdate = true,
    updateInterval = 15000, // 15 seconds
    enableNotifications = true
  } = options

  const [state, setState] = useState<GasOptimizationState>({
    gasEstimates: [],
    gasTracker: null,
    suggestions: [],
    isLoading: false,
    error: null,
    lastUpdated: 0
  })

  const intervalRef = useRef<NodeJS.Timeout>()
  const previousTrackerRef = useRef<GasTracker | null>(null)

  // Get gas estimates for a specific transaction
  const getGasEstimates = useCallback(async (gasLimit: bigint) => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const estimates = await gasOptimizer.getGasEstimates(chainId, gasLimit)
      setState(prev => ({
        ...prev,
        gasEstimates: estimates,
        isLoading: false,
        lastUpdated: Date.now()
      }))
      return estimates
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get gas estimates'
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      
      if (enableNotifications) {
        toast.error('Gas estimation failed', {
          description: errorMessage
        })
      }
      throw error
    }
  }, [chainId, enableNotifications])

  // Update gas tracker data
  const updateGasTracker = useCallback(() => {
    const tracker = gasOptimizer.getGasTracker(chainId)
    
    setState(prev => ({
      ...prev,
      gasTracker: tracker || null,
      lastUpdated: Date.now()
    }))

    // Check for significant gas price changes
    if (enableNotifications && tracker && previousTrackerRef.current) {
      const previous = previousTrackerRef.current
      const currentGwei = Number(tracker.currentGasPrice) / 1e9
      const previousGwei = Number(previous.currentGasPrice) / 1e9
      const change = (currentGwei - previousGwei) / previousGwei

      // Notify on significant changes (>20%)
      if (Math.abs(change) > 0.2) {
        const direction = change > 0 ? 'increased' : 'decreased'
        const percentage = Math.abs(change * 100).toFixed(0)
        
        toast.info(`Gas prices ${direction}`, {
          description: `Gas price ${direction} by ${percentage}% to ${currentGwei.toFixed(1)} gwei`,
          action: tracker.recommendedAction === 'wait' ? {
            label: 'Wait',
            onClick: () => {}
          } : tracker.recommendedAction === 'urgent' ? {
            label: 'Act Now',
            onClick: () => {}
          } : undefined
        })
      }

      // Notify on recommendation changes
      if (previous.recommendedAction !== tracker.recommendedAction) {
        const messages = {
          wait: 'Consider waiting for lower gas prices',
          proceed: 'Good time to proceed with transactions',
          urgent: 'Excellent time to transact - gas prices are low!'
        }

        toast.info('Gas recommendation updated', {
          description: messages[tracker.recommendedAction]
        })
      }
    }

    previousTrackerRef.current = tracker || null
  }, [chainId, enableNotifications])

  // Generate optimization suggestions
  const generateSuggestions = useCallback((
    transactionType: string,
    amount?: string
  ) => {
    const suggestions = gasOptimizer.generateOptimizationSuggestions(
      chainId,
      transactionType,
      amount
    )
    
    setState(prev => ({
      ...prev,
      suggestions,
      lastUpdated: Date.now()
    }))

    return suggestions
  }, [chainId])

  // Get recommended gas priority based on current conditions
  const getRecommendedPriority = useCallback((): GasPriority => {
    const tracker = state.gasTracker
    if (!tracker) return GasPriority.STANDARD

    switch (tracker.recommendedAction) {
      case 'wait':
        return GasPriority.SLOW
      case 'urgent':
        return GasPriority.FAST
      default:
        return GasPriority.STANDARD
    }
  }, [state.gasTracker])

  // Get gas estimate for specific priority
  const getEstimateForPriority = useCallback((priority: GasPriority): GasEstimate | null => {
    return state.gasEstimates.find(estimate => estimate.priority === priority) || null
  }, [state.gasEstimates])

  // Calculate potential savings between priorities
  const calculateSavings = useCallback((
    fromPriority: GasPriority,
    toPriority: GasPriority
  ): { amount: string; percentage: number } | null => {
    const fromEstimate = getEstimateForPriority(fromPriority)
    const toEstimate = getEstimateForPriority(toPriority)

    if (!fromEstimate || !toEstimate) return null

    const fromCost = parseFloat(fromEstimate.cost)
    const toCost = parseFloat(toEstimate.cost)
    const savings = fromCost - toCost
    const percentage = (savings / fromCost) * 100

    return {
      amount: savings.toFixed(6),
      percentage: Math.round(percentage)
    }
  }, [getEstimateForPriority])

  // Format gas price for display
  const formatGasPrice = useCallback((gasPrice: bigint): string => {
    const gwei = Number(gasPrice) / 1e9
    return `${gwei.toFixed(1)} gwei`
  }, [])

  // Get network congestion level
  const getCongestionLevel = useCallback((): {
    level: 'low' | 'medium' | 'high'
    color: string
    description: string
  } => {
    const tracker = state.gasTracker
    if (!tracker) {
      return {
        level: 'medium',
        color: 'text-yellow-500',
        description: 'Unknown'
      }
    }

    const colors = {
      low: 'text-green-500',
      medium: 'text-yellow-500',
      high: 'text-red-500'
    }

    const descriptions = {
      low: 'Low congestion - good time to transact',
      medium: 'Moderate congestion - normal gas prices',
      high: 'High congestion - consider waiting'
    }

    return {
      level: tracker.networkCongestion,
      color: colors[tracker.networkCongestion],
      description: descriptions[tracker.networkCongestion]
    }
  }, [state.gasTracker])

  // Start auto-update
  useEffect(() => {
    if (autoUpdate) {
      updateGasTracker()
      
      intervalRef.current = setInterval(() => {
        updateGasTracker()
      }, updateInterval)

      return () => {
        if (intervalRef.current) {
          clearInterval(intervalRef.current)
        }
      }
    }
  }, [autoUpdate, updateInterval, updateGasTracker])

  // Update tracker when chainId changes
  useEffect(() => {
    updateGasTracker()
  }, [chainId, updateGasTracker])

  return {
    // State
    ...state,

    // Actions
    getGasEstimates,
    updateGasTracker,
    generateSuggestions,

    // Computed values
    getRecommendedPriority,
    getEstimateForPriority,
    calculateSavings,
    formatGasPrice,
    getCongestionLevel,

    // Utilities
    isOptimalTime: state.gasTracker?.recommendedAction === 'urgent',
    shouldWait: state.gasTracker?.recommendedAction === 'wait',
    gasTrend: state.gasTracker?.trend || 'stable',
    currentGasPrice: state.gasTracker?.currentGasPrice,
    formattedGasPrice: state.gasTracker ? formatGasPrice(state.gasTracker.currentGasPrice) : 'Unknown'
  }
}
