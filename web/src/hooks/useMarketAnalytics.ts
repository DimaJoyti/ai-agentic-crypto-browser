import { useState, useEffect, useCallback } from 'react'
import { 
  marketAnalytics,
  type TechnicalIndicators,
  type VolumeAnalysis,
  type MarketSentiment,
  type MarketMetrics,
  type MarketIndicator,
  type CandlestickData
} from '@/lib/market-analytics'
import { usePriceFeed } from './usePriceFeed'

export interface MarketAnalyticsState {
  indicators: Map<string, TechnicalIndicators>
  volumeAnalysis: Map<string, VolumeAnalysis>
  sentiment: Map<string, MarketSentiment>
  metrics: Map<string, MarketMetrics>
  marketIndicators: Map<string, MarketIndicator[]>
  candlestickData: Map<string, CandlestickData[]>
  isLoading: boolean
  lastUpdate: number | null
}

export interface UseMarketAnalyticsOptions {
  symbols?: string[]
  autoUpdate?: boolean
  updateInterval?: number
  candlestickLimit?: number
}

export interface UseMarketAnalyticsReturn {
  // State
  state: MarketAnalyticsState
  
  // Data getters
  getTechnicalIndicators: (symbol: string) => TechnicalIndicators | null
  getVolumeAnalysis: (symbol: string) => VolumeAnalysis | null
  getMarketSentiment: (symbol: string) => MarketSentiment | null
  getMarketMetrics: (symbol: string) => MarketMetrics | null
  getMarketIndicators: (symbol: string) => MarketIndicator[]
  getCandlestickData: (symbol: string, limit?: number) => CandlestickData[]
  
  // Analysis functions
  getOverallMarketSentiment: () => MarketSentiment
  getTopPerformers: (metric: 'rsi' | 'volume' | 'sentiment') => Array<{ symbol: string; value: number }>
  getMarketSummary: () => {
    totalSymbols: number
    bullishSignals: number
    bearishSignals: number
    neutralSignals: number
    avgSentiment: number
    highVolumeAssets: number
  }
  
  // Utilities
  refresh: () => void
  addSymbol: (symbol: string) => void
  removeSymbol: (symbol: string) => void
}

export const useMarketAnalytics = (
  options: UseMarketAnalyticsOptions = {}
): UseMarketAnalyticsReturn => {
  const {
    symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'],
    autoUpdate = true,
    updateInterval = 30000, // 30 seconds
    candlestickLimit = 100
  } = options

  const [state, setState] = useState<MarketAnalyticsState>({
    indicators: new Map(),
    volumeAnalysis: new Map(),
    sentiment: new Map(),
    metrics: new Map(),
    marketIndicators: new Map(),
    candlestickData: new Map(),
    isLoading: true,
    lastUpdate: null
  })

  const [trackedSymbols, setTrackedSymbols] = useState<string[]>(symbols)

  // Use price feed to get real-time data
  const { subscribe } = usePriceFeed({
    symbols: trackedSymbols,
    autoStart: autoUpdate,
    onPriceUpdate: (symbol, priceData) => {
      // Add price data to analytics engine
      marketAnalytics.addPriceData(symbol, priceData)
      
      // Update state with new analytics data
      updateAnalyticsData()
    }
  })

  // Update analytics data from the engine
  const updateAnalyticsData = useCallback(() => {
    const newState: Partial<MarketAnalyticsState> = {
      indicators: new Map(),
      volumeAnalysis: new Map(),
      sentiment: new Map(),
      metrics: new Map(),
      marketIndicators: new Map(),
      candlestickData: new Map()
    }

    for (const symbol of trackedSymbols) {
      // Get technical indicators
      const indicators = marketAnalytics.getTechnicalIndicators(symbol)
      if (indicators) {
        newState.indicators!.set(symbol, indicators)
      }

      // Get volume analysis
      const volume = marketAnalytics.getVolumeAnalysis(symbol)
      if (volume) {
        newState.volumeAnalysis!.set(symbol, volume)
      }

      // Get market sentiment
      const sentiment = marketAnalytics.getMarketSentiment(symbol)
      if (sentiment) {
        newState.sentiment!.set(symbol, sentiment)
      }

      // Get market metrics
      const metrics = marketAnalytics.getMarketMetrics(symbol)
      if (metrics) {
        newState.metrics!.set(symbol, metrics)
      }

      // Get market indicators
      const marketIndicators = marketAnalytics.getMarketIndicators(symbol)
      if (marketIndicators.length > 0) {
        newState.marketIndicators!.set(symbol, marketIndicators)
      }

      // Get candlestick data
      const candlesticks = marketAnalytics.getCandlestickData(symbol, candlestickLimit)
      if (candlesticks.length > 0) {
        newState.candlestickData!.set(symbol, candlesticks)
      }
    }

    setState(prev => ({
      ...prev,
      ...newState,
      isLoading: false,
      lastUpdate: Date.now()
    }))
  }, [trackedSymbols, candlestickLimit])

  // Subscribe to price updates
  useEffect(() => {
    if (autoUpdate) {
      const unsubscribe = subscribe(trackedSymbols)
      return unsubscribe
    }
  }, [trackedSymbols, autoUpdate, subscribe])

  // Periodic update
  useEffect(() => {
    if (autoUpdate) {
      const interval = setInterval(updateAnalyticsData, updateInterval)
      return () => clearInterval(interval)
    }
  }, [autoUpdate, updateInterval, updateAnalyticsData])

  // Initial data load
  useEffect(() => {
    updateAnalyticsData()
  }, [updateAnalyticsData])

  // Data getter functions
  const getTechnicalIndicators = useCallback((symbol: string): TechnicalIndicators | null => {
    return marketAnalytics.getTechnicalIndicators(symbol)
  }, [])

  const getVolumeAnalysis = useCallback((symbol: string): VolumeAnalysis | null => {
    return marketAnalytics.getVolumeAnalysis(symbol)
  }, [])

  const getMarketSentiment = useCallback((symbol: string): MarketSentiment | null => {
    return marketAnalytics.getMarketSentiment(symbol)
  }, [])

  const getMarketMetrics = useCallback((symbol: string): MarketMetrics | null => {
    return marketAnalytics.getMarketMetrics(symbol)
  }, [])

  const getMarketIndicators = useCallback((symbol: string): MarketIndicator[] => {
    return marketAnalytics.getMarketIndicators(symbol)
  }, [])

  const getCandlestickData = useCallback((symbol: string, limit = candlestickLimit): CandlestickData[] => {
    return marketAnalytics.getCandlestickData(symbol, limit)
  }, [candlestickLimit])

  // Analysis functions
  const getOverallMarketSentiment = useCallback((): MarketSentiment => {
    const sentiments = Array.from(state.sentiment.values())
    
    if (sentiments.length === 0) {
      return {
        overall: 'neutral',
        score: 0,
        fearGreedIndex: 50,
        socialSentiment: 0,
        newssentiment: 0,
        technicalSentiment: 0,
        confidence: 0
      }
    }

    const avgScore = sentiments.reduce((sum, s) => sum + s.score, 0) / sentiments.length
    const avgFearGreed = sentiments.reduce((sum, s) => sum + s.fearGreedIndex, 0) / sentiments.length
    const avgSocial = sentiments.reduce((sum, s) => sum + s.socialSentiment, 0) / sentiments.length
    const avgNews = sentiments.reduce((sum, s) => sum + s.newssentiment, 0) / sentiments.length
    const avgTechnical = sentiments.reduce((sum, s) => sum + s.technicalSentiment, 0) / sentiments.length
    const avgConfidence = sentiments.reduce((sum, s) => sum + s.confidence, 0) / sentiments.length

    let overall: MarketSentiment['overall'] = 'neutral'
    if (avgScore > 40) overall = 'extremely_bullish'
    else if (avgScore > 15) overall = 'bullish'
    else if (avgScore < -40) overall = 'extremely_bearish'
    else if (avgScore < -15) overall = 'bearish'

    return {
      overall,
      score: avgScore,
      fearGreedIndex: avgFearGreed,
      socialSentiment: avgSocial,
      newssentiment: avgNews,
      technicalSentiment: avgTechnical,
      confidence: avgConfidence
    }
  }, [state.sentiment])

  const getTopPerformers = useCallback((metric: 'rsi' | 'volume' | 'sentiment') => {
    const results: Array<{ symbol: string; value: number }> = []

    for (const symbol of trackedSymbols) {
      let value = 0

      switch (metric) {
        case 'rsi':
          const indicators = state.indicators.get(symbol)
          if (indicators) value = indicators.rsi
          break
        case 'volume':
          const volume = state.volumeAnalysis.get(symbol)
          if (volume) value = volume.volumeRatio
          break
        case 'sentiment':
          const sentiment = state.sentiment.get(symbol)
          if (sentiment) value = sentiment.score
          break
      }

      if (value !== 0) {
        results.push({ symbol, value })
      }
    }

    return results.sort((a, b) => b.value - a.value).slice(0, 5)
  }, [trackedSymbols, state.indicators, state.volumeAnalysis, state.sentiment])

  const getMarketSummary = useCallback(() => {
    const indicators = Array.from(state.marketIndicators.values()).flat()
    
    const bullishSignals = indicators.filter(i => i.signal === 'bullish').length
    const bearishSignals = indicators.filter(i => i.signal === 'bearish').length
    const neutralSignals = indicators.filter(i => i.signal === 'neutral').length

    const sentiments = Array.from(state.sentiment.values())
    const avgSentiment = sentiments.length > 0 
      ? sentiments.reduce((sum, s) => sum + s.score, 0) / sentiments.length 
      : 0

    const volumes = Array.from(state.volumeAnalysis.values())
    const highVolumeAssets = volumes.filter(v => v.volumeRatio > 1.5).length

    return {
      totalSymbols: trackedSymbols.length,
      bullishSignals,
      bearishSignals,
      neutralSignals,
      avgSentiment,
      highVolumeAssets
    }
  }, [state.marketIndicators, state.sentiment, state.volumeAnalysis, trackedSymbols])

  // Utility functions
  const refresh = useCallback(() => {
    setState(prev => ({ ...prev, isLoading: true }))
    updateAnalyticsData()
  }, [updateAnalyticsData])

  const addSymbol = useCallback((symbol: string) => {
    if (!trackedSymbols.includes(symbol.toUpperCase())) {
      setTrackedSymbols(prev => [...prev, symbol.toUpperCase()])
    }
  }, [trackedSymbols])

  const removeSymbol = useCallback((symbol: string) => {
    setTrackedSymbols(prev => prev.filter(s => s !== symbol.toUpperCase()))
  }, [])

  return {
    state,
    getTechnicalIndicators,
    getVolumeAnalysis,
    getMarketSentiment,
    getMarketMetrics,
    getMarketIndicators,
    getCandlestickData,
    getOverallMarketSentiment,
    getTopPerformers,
    getMarketSummary,
    refresh,
    addSymbol,
    removeSymbol
  }
}

// Hook for single symbol analytics
export const useSymbolAnalytics = (symbol: string) => {
  const { 
    getTechnicalIndicators,
    getVolumeAnalysis,
    getMarketSentiment,
    getMarketMetrics,
    getMarketIndicators,
    getCandlestickData,
    state
  } = useMarketAnalytics({
    symbols: [symbol],
    autoUpdate: true
  })

  return {
    indicators: getTechnicalIndicators(symbol),
    volumeAnalysis: getVolumeAnalysis(symbol),
    sentiment: getMarketSentiment(symbol),
    metrics: getMarketMetrics(symbol),
    marketIndicators: getMarketIndicators(symbol),
    candlestickData: getCandlestickData(symbol),
    isLoading: state.isLoading,
    lastUpdate: state.lastUpdate
  }
}
