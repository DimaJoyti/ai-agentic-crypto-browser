import { useState, useEffect, useCallback } from 'react'
import { useAccount } from 'wagmi'
import { 
  portfolioAnalytics,
  type PortfolioPosition,
  type PortfolioTransaction,
  type PortfolioSummary,
  type PortfolioMetrics,
  type AssetAllocation,
  type RiskMetrics
} from '@/lib/portfolio-analytics'
import { usePriceFeed } from './usePriceFeed'
import { toast } from 'sonner'

export interface PortfolioState {
  summary: PortfolioSummary | null
  positions: PortfolioPosition[]
  transactions: PortfolioTransaction[]
  isLoading: boolean
  lastUpdate: number | null
  error: string | null
}

export interface UsePortfolioAnalyticsOptions {
  autoSync?: boolean
  syncInterval?: number
  trackPriceUpdates?: boolean
  symbols?: string[]
}

export interface UsePortfolioAnalyticsReturn {
  // State
  state: PortfolioState
  
  // Portfolio management
  addPosition: (position: Partial<PortfolioPosition> & { symbol: string }) => void
  updatePosition: (symbol: string, updates: Partial<PortfolioPosition>) => void
  removePosition: (symbol: string) => void
  
  // Transaction management
  addTransaction: (transaction: Omit<PortfolioTransaction, 'id'>) => void
  getTransactionHistory: (symbol?: string) => PortfolioTransaction[]
  
  // Analytics
  getMetrics: () => PortfolioMetrics | null
  getAllocation: () => AssetAllocation[]
  getRiskMetrics: () => RiskMetrics | null
  getTopGainers: () => PortfolioPosition[]
  getTopLosers: () => PortfolioPosition[]
  
  // Portfolio operations
  importPortfolio: (positions: Partial<PortfolioPosition>[]) => void
  exportPortfolio: () => any
  clearPortfolio: () => void
  refreshPortfolio: () => void
  
  // Utilities
  calculatePositionValue: (symbol: string, amount: number) => number
  getPositionBySymbol: (symbol: string) => PortfolioPosition | null
  getTotalValue: () => number
  getTotalPnL: () => number
  clearError: () => void
}

export const usePortfolioAnalytics = (
  options: UsePortfolioAnalyticsOptions = {}
): UsePortfolioAnalyticsReturn => {
  const {
    autoSync = true,
    syncInterval = 30000, // 30 seconds
    trackPriceUpdates = true,
    symbols = []
  } = options

  const [state, setState] = useState<PortfolioState>({
    summary: null,
    positions: [],
    transactions: [],
    isLoading: true,
    lastUpdate: null,
    error: null
  })

  const { address } = useAccount()
  const userId = address || 'demo_user'

  // Track price updates for portfolio symbols
  const portfolioSymbols = state.positions.map(p => p.symbol)
  const allSymbols = Array.from(new Set([...symbols, ...portfolioSymbols]))

  const { subscribe } = usePriceFeed({
    symbols: allSymbols,
    autoStart: trackPriceUpdates,
    onPriceUpdate: (symbol, priceData) => {
      if (trackPriceUpdates) {
        portfolioAnalytics.updatePriceData(priceData)
        refreshPortfolio()
      }
    }
  })

  // Subscribe to price updates
  useEffect(() => {
    if (trackPriceUpdates && allSymbols.length > 0) {
      const unsubscribe = subscribe(allSymbols)
      return unsubscribe
    }
  }, [trackPriceUpdates, allSymbols, subscribe])

  // Auto-sync portfolio data
  useEffect(() => {
    if (autoSync) {
      const interval = setInterval(refreshPortfolio, syncInterval)
      return () => clearInterval(interval)
    }
  }, [autoSync, syncInterval])

  // Initial load
  useEffect(() => {
    refreshPortfolio()
  }, [userId])

  // Refresh portfolio data
  const refreshPortfolio = useCallback(() => {
    try {
      const summary = portfolioAnalytics.getPortfolioSummary(userId)
      const positions = portfolioAnalytics.getPositions(userId)
      const transactions = portfolioAnalytics.getTransactions(userId, 50)

      setState(prev => ({
        ...prev,
        summary,
        positions,
        transactions,
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
      console.error('Failed to refresh portfolio:', error)
    }
  }, [userId])

  // Add position
  const addPosition = useCallback((position: Partial<PortfolioPosition> & { symbol: string }) => {
    try {
      portfolioAnalytics.updatePosition(userId, position)
      refreshPortfolio()
      toast.success(`Added ${position.symbol} to portfolio`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to add position', { description: errorMessage })
    }
  }, [userId, refreshPortfolio])

  // Update position
  const updatePosition = useCallback((symbol: string, updates: Partial<PortfolioPosition>) => {
    try {
      const existingPosition = state.positions.find(p => p.symbol === symbol)
      if (!existingPosition) {
        throw new Error(`Position ${symbol} not found`)
      }

      portfolioAnalytics.updatePosition(userId, { ...existingPosition, ...updates })
      refreshPortfolio()
      toast.success(`Updated ${symbol} position`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to update position', { description: errorMessage })
    }
  }, [userId, state.positions, refreshPortfolio])

  // Remove position
  const removePosition = useCallback((symbol: string) => {
    try {
      portfolioAnalytics.updatePosition(userId, { symbol, amount: 0 })
      refreshPortfolio()
      toast.success(`Removed ${symbol} from portfolio`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to remove position', { description: errorMessage })
    }
  }, [userId, refreshPortfolio])

  // Add transaction
  const addTransaction = useCallback((transaction: Omit<PortfolioTransaction, 'id'>) => {
    try {
      portfolioAnalytics.addTransaction(userId, transaction)
      refreshPortfolio()
      toast.success(`Added ${transaction.type} transaction for ${transaction.symbol}`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to add transaction', { description: errorMessage })
    }
  }, [userId, refreshPortfolio])

  // Get transaction history
  const getTransactionHistory = useCallback((symbol?: string): PortfolioTransaction[] => {
    const allTransactions = portfolioAnalytics.getTransactions(userId)
    return symbol 
      ? allTransactions.filter(tx => tx.symbol === symbol)
      : allTransactions
  }, [userId])

  // Get metrics
  const getMetrics = useCallback((): PortfolioMetrics | null => {
    return state.summary?.metrics || null
  }, [state.summary])

  // Get allocation
  const getAllocation = useCallback((): AssetAllocation[] => {
    return state.summary?.allocation || []
  }, [state.summary])

  // Get risk metrics
  const getRiskMetrics = useCallback((): RiskMetrics | null => {
    return state.summary?.riskMetrics || null
  }, [state.summary])

  // Get top gainers
  const getTopGainers = useCallback((): PortfolioPosition[] => {
    return state.summary?.topGainers || []
  }, [state.summary])

  // Get top losers
  const getTopLosers = useCallback((): PortfolioPosition[] => {
    return state.summary?.topLosers || []
  }, [state.summary])

  // Import portfolio
  const importPortfolio = useCallback((positions: Partial<PortfolioPosition>[]) => {
    try {
      portfolioAnalytics.importPortfolio(userId, positions)
      refreshPortfolio()
      toast.success(`Imported ${positions.length} positions`)
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to import portfolio', { description: errorMessage })
    }
  }, [userId, refreshPortfolio])

  // Export portfolio
  const exportPortfolio = useCallback(() => {
    try {
      const data = portfolioAnalytics.exportPortfolio(userId)
      toast.success('Portfolio exported successfully')
      return data
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to export portfolio', { description: errorMessage })
      return null
    }
  }, [userId])

  // Clear portfolio
  const clearPortfolio = useCallback(() => {
    try {
      portfolioAnalytics.clearPortfolio(userId)
      refreshPortfolio()
      toast.success('Portfolio cleared')
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      toast.error('Failed to clear portfolio', { description: errorMessage })
    }
  }, [userId, refreshPortfolio])

  // Calculate position value
  const calculatePositionValue = useCallback((symbol: string, amount: number): number => {
    const position = state.positions.find(p => p.symbol === symbol)
    return position ? amount * position.currentPrice : 0
  }, [state.positions])

  // Get position by symbol
  const getPositionBySymbol = useCallback((symbol: string): PortfolioPosition | null => {
    return state.positions.find(p => p.symbol === symbol) || null
  }, [state.positions])

  // Get total value
  const getTotalValue = useCallback((): number => {
    return state.summary?.metrics.totalValue || 0
  }, [state.summary])

  // Get total P&L
  const getTotalPnL = useCallback((): number => {
    return state.summary?.metrics.totalPnL || 0
  }, [state.summary])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    addPosition,
    updatePosition,
    removePosition,
    addTransaction,
    getTransactionHistory,
    getMetrics,
    getAllocation,
    getRiskMetrics,
    getTopGainers,
    getTopLosers,
    importPortfolio,
    exportPortfolio,
    clearPortfolio,
    refreshPortfolio,
    calculatePositionValue,
    getPositionBySymbol,
    getTotalValue,
    getTotalPnL,
    clearError
  }
}

// Hook for simplified portfolio tracking
export const usePortfolioValue = () => {
  const { state, getTotalValue, getTotalPnL } = usePortfolioAnalytics({
    autoSync: true,
    trackPriceUpdates: true
  })

  return {
    totalValue: getTotalValue(),
    totalPnL: getTotalPnL(),
    totalPnLPercent: state.summary?.metrics.totalPnLPercent || 0,
    dayChange: state.summary?.metrics.dayChange || 0,
    dayChangePercent: state.summary?.metrics.dayChangePercent || 0,
    isLoading: state.isLoading,
    positionCount: state.positions.length
  }
}
