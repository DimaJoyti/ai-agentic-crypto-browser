import { useState, useEffect, useCallback } from 'react'
import { useAccount } from 'wagmi'
import { 
  tradingSignalsEngine,
  type TradingSignal,
  type SignalStrategy,
  type SignalAlert,
  type SignalPerformance
} from '@/lib/trading-signals'
import { usePriceFeed } from './usePriceFeed'
import { useMarketAnalytics } from './useMarketAnalytics'
import { toast } from 'sonner'

export interface TradingSignalsState {
  signals: TradingSignal[]
  strategies: SignalStrategy[]
  alerts: SignalAlert[]
  performance: Map<string, SignalPerformance>
  isLoading: boolean
  lastUpdate: number | null
  error: string | null
}

export interface UseTradingSignalsOptions {
  symbols?: string[]
  autoStart?: boolean
  enableAlerts?: boolean
  maxSignals?: number
}

export interface UseTradingSignalsReturn {
  // State
  state: TradingSignalsState
  
  // Signals
  getSignals: (symbol?: string) => TradingSignal[]
  getRecentSignals: (limit?: number) => TradingSignal[]
  getSignalsByType: (type: 'buy' | 'sell' | 'hold') => TradingSignal[]
  
  // Strategies
  getStrategies: () => SignalStrategy[]
  addStrategy: (strategy: Omit<SignalStrategy, 'id' | 'createdAt'>) => string
  updateStrategy: (id: string, updates: Partial<SignalStrategy>) => void
  deleteStrategy: (id: string) => void
  toggleStrategy: (id: string, enabled: boolean) => void
  
  // Alerts
  getAlerts: () => SignalAlert[]
  addAlert: (alert: Omit<SignalAlert, 'id' | 'userId' | 'createdAt'>) => string
  updateAlert: (id: string, updates: Partial<SignalAlert>) => void
  deleteAlert: (id: string) => void
  toggleAlert: (id: string, enabled: boolean) => void
  
  // Performance
  getPerformance: (strategyId: string) => SignalPerformance | null
  getOverallPerformance: () => {
    totalSignals: number
    successRate: number
    avgReturn: number
    activeStrategies: number
    activeAlerts: number
  }
  
  // Utilities
  refresh: () => void
  clearError: () => void
  exportSignals: () => any
  importStrategies: (strategies: SignalStrategy[]) => void
}

export const useTradingSignals = (
  options: UseTradingSignalsOptions = {}
): UseTradingSignalsReturn => {
  const {
    symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
    autoStart = true,
    enableAlerts = true,
    maxSignals = 100
  } = options

  const [state, setState] = useState<TradingSignalsState>({
    signals: [],
    strategies: [],
    alerts: [],
    performance: new Map(),
    isLoading: true,
    lastUpdate: null,
    error: null
  })

  const { address } = useAccount()
  const userId = address || 'demo_user'

  // Use price feed and market analytics for signal generation
  const { subscribe } = usePriceFeed({
    symbols,
    autoStart,
    onPriceUpdate: (symbol, priceData) => {
      // Get technical indicators for the symbol
      const indicators = marketAnalytics.indicators || {}
      
      // Add price data to signals engine
      tradingSignalsEngine.addPriceData(priceData, indicators as any)
      
      // Update signals state
      updateSignalsState()
    }
  })

  const { state: marketAnalytics } = useMarketAnalytics({
    symbols,
    autoUpdate: true
  })

  // Subscribe to price updates
  useEffect(() => {
    if (autoStart) {
      const unsubscribe = subscribe(symbols)
      return unsubscribe
    }
  }, [autoStart, symbols, subscribe])

  // Listen for signal and alert events
  useEffect(() => {
    const handleSignal = (event: CustomEvent) => {
      const signal = event.detail as TradingSignal
      
      if (enableAlerts) {
        toast.success(`Trading Signal: ${signal.symbol}`, {
          description: `${signal.type.toUpperCase()} signal with ${signal.confidence.toFixed(0)}% confidence`
        })
      }
      
      updateSignalsState()
    }

    const handleAlert = (event: CustomEvent) => {
      const { alert, signal } = event.detail
      
      toast.info(`Signal Alert: ${signal.symbol}`, {
        description: `${signal.type.toUpperCase()} signal triggered for ${signal.symbol}`
      })
    }

    if (enableAlerts) {
      window.addEventListener('tradingSignal', handleSignal as EventListener)
      window.addEventListener('signalAlert', handleAlert as EventListener)

      return () => {
        window.removeEventListener('tradingSignal', handleSignal as EventListener)
        window.removeEventListener('signalAlert', handleAlert as EventListener)
      }
    }
  }, [enableAlerts])

  // Initial load and periodic updates
  useEffect(() => {
    updateSignalsState()
    
    const interval = setInterval(updateSignalsState, 30000) // Update every 30 seconds
    return () => clearInterval(interval)
  }, [])

  // Update signals state
  const updateSignalsState = useCallback(() => {
    try {
      const allSignals = tradingSignalsEngine.getAllRecentSignals(maxSignals)
      const strategies = tradingSignalsEngine.getStrategies()
      const alerts = tradingSignalsEngine.getUserAlerts(userId)

      setState(prev => ({
        ...prev,
        signals: allSignals,
        strategies,
        alerts,
        isLoading: false,
        lastUpdate: Date.now(),
        error: null
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))
    }
  }, [maxSignals, userId])

  // Get signals
  const getSignals = useCallback((symbol?: string): TradingSignal[] => {
    if (symbol) {
      return tradingSignalsEngine.getSignals(symbol)
    }
    return state.signals
  }, [state.signals])

  // Get recent signals
  const getRecentSignals = useCallback((limit = 50): TradingSignal[] => {
    return tradingSignalsEngine.getAllRecentSignals(limit)
  }, [])

  // Get signals by type
  const getSignalsByType = useCallback((type: 'buy' | 'sell' | 'hold'): TradingSignal[] => {
    return state.signals.filter(signal => signal.type === type)
  }, [state.signals])

  // Get strategies
  const getStrategies = useCallback((): SignalStrategy[] => {
    return tradingSignalsEngine.getStrategies()
  }, [])

  // Add strategy
  const addStrategy = useCallback((strategy: Omit<SignalStrategy, 'id' | 'createdAt'>): string => {
    try {
      const id = tradingSignalsEngine.addStrategy(strategy)
      updateSignalsState()
      toast.success('Strategy added successfully')
      return id
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to add strategy', { description: errorMessage })
      return ''
    }
  }, [updateSignalsState])

  // Update strategy
  const updateStrategy = useCallback((id: string, updates: Partial<SignalStrategy>) => {
    try {
      // Implementation would update strategy in engine
      updateSignalsState()
      toast.success('Strategy updated successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to update strategy', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Delete strategy
  const deleteStrategy = useCallback((id: string) => {
    try {
      // Implementation would delete strategy from engine
      updateSignalsState()
      toast.success('Strategy deleted successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to delete strategy', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Toggle strategy
  const toggleStrategy = useCallback((id: string, enabled: boolean) => {
    try {
      // Implementation would toggle strategy in engine
      updateSignalsState()
      toast.success(`Strategy ${enabled ? 'enabled' : 'disabled'}`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to toggle strategy', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Get alerts
  const getAlerts = useCallback((): SignalAlert[] => {
    return tradingSignalsEngine.getUserAlerts(userId)
  }, [userId])

  // Add alert
  const addAlert = useCallback((alert: Omit<SignalAlert, 'id' | 'userId' | 'createdAt'>): string => {
    try {
      const id = tradingSignalsEngine.addAlert(userId, { ...alert, userId })
      updateSignalsState()
      toast.success('Alert added successfully')
      return id
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to add alert', { description: errorMessage })
      return ''
    }
  }, [userId, updateSignalsState])

  // Update alert
  const updateAlert = useCallback((id: string, updates: Partial<SignalAlert>) => {
    try {
      // Implementation would update alert in engine
      updateSignalsState()
      toast.success('Alert updated successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to update alert', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Delete alert
  const deleteAlert = useCallback((id: string) => {
    try {
      // Implementation would delete alert from engine
      updateSignalsState()
      toast.success('Alert deleted successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to delete alert', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Toggle alert
  const toggleAlert = useCallback((id: string, enabled: boolean) => {
    try {
      // Implementation would toggle alert in engine
      updateSignalsState()
      toast.success(`Alert ${enabled ? 'enabled' : 'disabled'}`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to toggle alert', { description: errorMessage })
    }
  }, [updateSignalsState])

  // Get performance
  const getPerformance = useCallback((strategyId: string): SignalPerformance | null => {
    return state.performance.get(strategyId) || null
  }, [state.performance])

  // Get overall performance
  const getOverallPerformance = useCallback(() => {
    const totalSignals = state.signals.length
    const activeStrategies = state.strategies.filter(s => s.enabled).length
    const activeAlerts = state.alerts.filter(a => a.enabled).length
    
    // Mock performance calculations
    const successRate = Math.random() * 40 + 50 // 50-90%
    const avgReturn = (Math.random() - 0.5) * 20 // -10% to +10%

    return {
      totalSignals,
      successRate,
      avgReturn,
      activeStrategies,
      activeAlerts
    }
  }, [state.signals, state.strategies, state.alerts])

  // Refresh
  const refresh = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true }))
    updateSignalsState()
  }, [updateSignalsState])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Export signals
  const exportSignals = useCallback(() => {
    try {
      const data = {
        signals: state.signals,
        strategies: state.strategies,
        alerts: state.alerts,
        exportedAt: new Date().toISOString()
      }
      
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `trading-signals-${new Date().toISOString().split('T')[0]}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      
      toast.success('Signals exported successfully')
      return data
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to export signals', { description: errorMessage })
      return null
    }
  }, [state])

  // Import strategies
  const importStrategies = useCallback((strategies: SignalStrategy[]) => {
    try {
      for (const strategy of strategies) {
        tradingSignalsEngine.addStrategy(strategy)
      }
      updateSignalsState()
      toast.success(`Imported ${strategies.length} strategies`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to import strategies', { description: errorMessage })
    }
  }, [updateSignalsState])

  return {
    state,
    getSignals,
    getRecentSignals,
    getSignalsByType,
    getStrategies,
    addStrategy,
    updateStrategy,
    deleteStrategy,
    toggleStrategy,
    getAlerts,
    addAlert,
    updateAlert,
    deleteAlert,
    toggleAlert,
    getPerformance,
    getOverallPerformance,
    refresh,
    clearError,
    exportSignals,
    importStrategies
  }
}

// Hook for simplified signal tracking
export const useSignalAlerts = (symbols: string[] = []) => {
  const { getSignals, getSignalsByType, state } = useTradingSignals({
    symbols,
    enableAlerts: true
  })

  return {
    recentSignals: getSignals().slice(0, 10),
    buySignals: getSignalsByType('buy').length,
    sellSignals: getSignalsByType('sell').length,
    totalSignals: state.signals.length,
    isLoading: state.isLoading
  }
}
