import { useState, useEffect, useCallback } from 'react'
import { 
  historicalDataManager,
  type HistoricalCandle,
  type DataStats,
  type DataRequest,
  type DataSource
} from '@/lib/historical-data'
import { 
  backtestingEngine,
  type BacktestConfig,
  type BacktestResults
} from '@/lib/backtesting-engine'
import { toast } from 'sonner'

export interface HistoricalDataState {
  candles: HistoricalCandle[]
  stats: DataStats
  requests: DataRequest[]
  sources: DataSource[]
  isLoading: boolean
  error: string | null
  lastUpdate: number | null
}

export interface BacktestState {
  results: BacktestResults[]
  currentBacktest: BacktestResults | null
  isRunning: boolean
  progress: number
  error: string | null
}

export interface UseHistoricalDataOptions {
  symbol?: string
  timeframe?: string
  autoLoad?: boolean
  cacheEnabled?: boolean
}

export interface UseHistoricalDataReturn {
  // Historical Data State
  dataState: HistoricalDataState
  
  // Data Management
  loadHistoricalData: (symbol: string, timeframe: string, startTime: number, endTime: number) => Promise<HistoricalCandle[]>
  getCandles: (symbol: string, timeframe: string, limit?: number) => HistoricalCandle[]
  getAvailableSymbols: () => string[]
  getAvailableTimeframes: (symbol: string) => string[]
  
  // Data Sources
  getDataSources: () => DataSource[]
  addDataSource: (source: DataSource) => void
  
  // Data Statistics
  getDataStats: () => DataStats
  clearOldData: (beforeTimestamp: number) => void
  
  // Import/Export
  exportData: (symbol?: string, timeframe?: string) => any
  importData: (data: any, symbol?: string, timeframe?: string) => void
  
  // Backtesting State
  backtestState: BacktestState
  
  // Backtesting Operations
  runBacktest: (config: BacktestConfig) => Promise<BacktestResults>
  getBacktests: () => BacktestResults[]
  getBacktest: (id: string) => BacktestResults | null
  deleteBacktest: (id: string) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useHistoricalData = (
  options: UseHistoricalDataOptions = {}
): UseHistoricalDataReturn => {
  const {
    symbol,
    timeframe,
    autoLoad = false,
    cacheEnabled = true
  } = options

  const [dataState, setDataState] = useState<HistoricalDataState>({
    candles: [],
    stats: {
      totalCandles: 0,
      symbols: 0,
      timeframes: 0,
      oldestData: Date.now(),
      newestData: 0,
      dataGaps: [],
      storageSize: 0
    },
    requests: [],
    sources: [],
    isLoading: false,
    error: null,
    lastUpdate: null
  })

  const [backtestState, setBacktestState] = useState<BacktestState>({
    results: [],
    currentBacktest: null,
    isRunning: false,
    progress: 0,
    error: null
  })

  // Auto-load data on mount
  useEffect(() => {
    if (autoLoad && symbol && timeframe) {
      const endTime = Date.now()
      const startTime = endTime - 30 * 24 * 60 * 60 * 1000 // 30 days ago
      loadHistoricalData(symbol, timeframe, startTime, endTime)
    }
    
    updateDataState()
    updateBacktestState()
  }, [autoLoad, symbol, timeframe])

  // Update data state
  const updateDataState = useCallback(() => {
    try {
      const stats = historicalDataManager.getDataStats()
      const requests = historicalDataManager.getDataRequests()
      const sources = historicalDataManager.getDataSources()
      
      let candles: HistoricalCandle[] = []
      if (symbol && timeframe) {
        candles = getCandles(symbol, timeframe)
      }

      setDataState(prev => ({
        ...prev,
        candles,
        stats,
        requests,
        sources,
        isLoading: false,
        lastUpdate: Date.now(),
        error: null
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setDataState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))
    }
  }, [symbol, timeframe])

  // Update backtest state
  const updateBacktestState = useCallback(() => {
    try {
      const results = backtestingEngine.getCompletedBacktests()
      
      setBacktestState(prev => ({
        ...prev,
        results,
        error: null
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setBacktestState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [])

  // Load historical data
  const loadHistoricalData = useCallback(async (
    symbol: string,
    timeframe: string,
    startTime: number,
    endTime: number
  ): Promise<HistoricalCandle[]> => {
    setDataState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const candles = await historicalDataManager.getHistoricalData(
        symbol,
        timeframe,
        startTime,
        endTime,
        !cacheEnabled
      )

      updateDataState()
      toast.success(`Loaded ${candles.length} candles for ${symbol} ${timeframe}`)
      return candles
    } catch (error) {
      const errorMessage = (error as Error).message
      setDataState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))
      toast.error('Failed to load historical data', { description: errorMessage })
      throw error
    }
  }, [cacheEnabled, updateDataState])

  // Get candles
  const getCandles = useCallback((
    symbol: string,
    timeframe: string,
    limit?: number
  ): HistoricalCandle[] => {
    const endTime = Date.now()
    const startTime = limit 
      ? endTime - limit * getTimeframeMs(timeframe)
      : endTime - 365 * 24 * 60 * 60 * 1000 // 1 year default

    const result = await historicalDataManager.getHistoricalData(symbol, timeframe, startTime, endTime)
    return Array.isArray(result) ? result : []
      .then(candles => limit ? candles.slice(-limit) : candles)
      .catch(() => [])
  }, [])

  // Get available symbols
  const getAvailableSymbols = useCallback((): string[] => {
    return historicalDataManager.getAvailableSymbols()
  }, [])

  // Get available timeframes
  const getAvailableTimeframes = useCallback((symbol: string): string[] => {
    return historicalDataManager.getAvailableTimeframes(symbol)
  }, [])

  // Get data sources
  const getDataSources = useCallback((): DataSource[] => {
    return historicalDataManager.getDataSources()
  }, [])

  // Add data source
  const addDataSource = useCallback((source: DataSource) => {
    try {
      historicalDataManager.addDataSource(source)
      updateDataState()
      toast.success(`Added data source: ${source.name}`)
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to add data source', { description: errorMessage })
    }
  }, [updateDataState])

  // Get data statistics
  const getDataStats = useCallback((): DataStats => {
    return historicalDataManager.getDataStats()
  }, [])

  // Clear old data
  const clearOldData = useCallback((beforeTimestamp: number) => {
    try {
      historicalDataManager.clearOldData(beforeTimestamp)
      updateDataState()
      toast.success('Old data cleared successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to clear old data', { description: errorMessage })
    }
  }, [updateDataState])

  // Export data
  const exportData = useCallback((symbol?: string, timeframe?: string) => {
    try {
      const data = historicalDataManager.exportData(symbol, timeframe)
      
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `historical-data-${symbol || 'all'}-${timeframe || 'all'}-${new Date().toISOString().split('T')[0]}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
      
      toast.success('Data exported successfully')
      return data
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to export data', { description: errorMessage })
      return null
    }
  }, [])

  // Import data
  const importData = useCallback((data: any, symbol?: string, timeframe?: string) => {
    try {
      historicalDataManager.importData(data, symbol, timeframe)
      updateDataState()
      toast.success('Data imported successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to import data', { description: errorMessage })
    }
  }, [updateDataState])

  // Run backtest
  const runBacktest = useCallback(async (config: BacktestConfig): Promise<BacktestResults> => {
    setBacktestState(prev => ({
      ...prev,
      isRunning: true,
      progress: 0,
      error: null
    }))

    try {
      // Load historical data for backtest
      const historicalData = await historicalDataManager.getHistoricalData(
        config.symbol,
        config.timeframe,
        config.startTime,
        config.endTime
      )

      if (historicalData.length === 0) {
        throw new Error('No historical data available for the specified period')
      }

      // Run backtest with progress callback
      const results = await backtestingEngine.runBacktest(
        config,
        historicalData,
        (progress) => {
          setBacktestState(prev => ({ ...prev, progress }))
        }
      )

      setBacktestState(prev => ({
        ...prev,
        isRunning: false,
        progress: 100,
        currentBacktest: results
      }))

      updateBacktestState()
      toast.success(`Backtest completed: ${results.performance.totalReturnPercent.toFixed(2)}% return`)
      return results
    } catch (error) {
      const errorMessage = (error as Error).message
      setBacktestState(prev => ({
        ...prev,
        isRunning: false,
        progress: 0,
        error: errorMessage
      }))
      toast.error('Backtest failed', { description: errorMessage })
      throw error
    }
  }, [updateBacktestState])

  // Get backtests
  const getBacktests = useCallback((): BacktestResults[] => {
    return backtestingEngine.getCompletedBacktests()
  }, [])

  // Get backtest
  const getBacktest = useCallback((id: string): BacktestResults | null => {
    return backtestingEngine.getBacktest(id)
  }, [])

  // Delete backtest
  const deleteBacktest = useCallback((id: string) => {
    try {
      backtestingEngine.deleteBacktest(id)
      updateBacktestState()
      toast.success('Backtest deleted successfully')
    } catch (error) {
      const errorMessage = (error as Error).message
      toast.error('Failed to delete backtest', { description: errorMessage })
    }
  }, [updateBacktestState])

  // Refresh
  const refresh = useCallback(() => {
    updateDataState()
    updateBacktestState()
  }, [updateDataState, updateBacktestState])

  // Clear error
  const clearError = useCallback(() => {
    setDataState(prev => ({ ...prev, error: null }))
    setBacktestState(prev => ({ ...prev, error: null }))
  }, [])

  // Helper function to get timeframe in milliseconds
  const getTimeframeMs = (timeframe: string): number => {
    const timeframes: Record<string, number> = {
      '1m': 60 * 1000,
      '5m': 5 * 60 * 1000,
      '15m': 15 * 60 * 1000,
      '1h': 60 * 60 * 1000,
      '4h': 4 * 60 * 60 * 1000,
      '1d': 24 * 60 * 60 * 1000,
      '1w': 7 * 24 * 60 * 60 * 1000,
      '1M': 30 * 24 * 60 * 60 * 1000
    }
    return timeframes[timeframe] || 60 * 60 * 1000
  }

  return {
    dataState,
    loadHistoricalData,
    getCandles,
    getAvailableSymbols,
    getAvailableTimeframes,
    getDataSources,
    addDataSource,
    getDataStats,
    clearOldData,
    exportData,
    importData,
    backtestState,
    runBacktest,
    getBacktests,
    getBacktest,
    deleteBacktest,
    refresh,
    clearError
  }
}

// Hook for simplified historical data access
export const useCandles = (symbol: string, timeframe: string, limit = 100) => {
  const { dataState, loadHistoricalData, getCandles } = useHistoricalData({
    symbol,
    timeframe,
    autoLoad: true
  })

  return {
    candles: getCandles(symbol, timeframe, limit),
    isLoading: dataState.isLoading,
    error: dataState.error,
    reload: () => {
      const endTime = Date.now()
      const startTime = endTime - limit * getTimeframeMs(timeframe)
      return loadHistoricalData(symbol, timeframe, startTime, endTime)
    }
  }
}

// Helper function
const getTimeframeMs = (timeframe: string): number => {
  const timeframes: Record<string, number> = {
    '1m': 60 * 1000,
    '5m': 5 * 60 * 1000,
    '15m': 15 * 60 * 1000,
    '1h': 60 * 60 * 1000,
    '4h': 4 * 60 * 60 * 1000,
    '1d': 24 * 60 * 60 * 1000,
    '1w': 7 * 24 * 60 * 60 * 1000,
    '1M': 30 * 24 * 60 * 60 * 1000
  }
  return timeframes[timeframe] || 60 * 60 * 1000
}
