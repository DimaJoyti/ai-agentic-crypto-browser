import { useState, useEffect, useCallback } from 'react'
import { tradingApi } from '@/lib/trading-api'
import { useWebSocket } from './useWebSocket'

interface HFTMetrics {
  ordersPerSecond: number
  avgLatencyMicros: number
  totalOrders: number
  successfulOrders: number
  failedOrders: number
  uptime: string
  isRunning: boolean
}

interface MarketData {
  symbol: string
  price: string
  change24h: string
  volume: string
  high24h: string
  low24h: string
  timestamp: string
}

interface TradingSignal {
  id: string
  symbol: string
  side: 'BUY' | 'SELL'
  price: string
  confidence: number
  source: string
  timestamp: string
}

interface UseTradingDashboardReturn {
  // HFT Engine
  hftMetrics: HFTMetrics | null
  hftStatus: any
  isHFTRunning: boolean
  startHFTEngine: () => Promise<void>
  stopHFTEngine: () => Promise<void>
  
  // Market Data
  marketData: MarketData[]
  
  // Trading
  orders: any[]
  positions: any[]
  signals: TradingSignal[]
  
  // Portfolio
  portfolioSummary: any
  portfolioMetrics: any
  portfolioRisk: any
  
  // Strategies
  strategies: any[]
  
  // Risk Management
  riskLimits: any[]
  riskViolations: any[]
  riskMetrics: any
  
  // System
  systemStatus: any

  // Connection
  isConnected: boolean

  // Loading states
  isLoading: boolean
  error: string | null
  
  // Actions
  refreshData: () => Promise<void>
  emergencyStop: (reason: string) => Promise<void>
}

export function useTradingDashboard(): UseTradingDashboardReturn {
  // State
  const [hftMetrics, setHftMetrics] = useState<HFTMetrics | null>(null)
  const [hftStatus, setHftStatus] = useState<any>(null)
  const [isHFTRunning, setIsHFTRunning] = useState(false)
  const [marketData, setMarketData] = useState<MarketData[]>([])
  const [orders, setOrders] = useState<any[]>([])
  const [positions, setPositions] = useState<any[]>([])
  const [signals, setSignals] = useState<TradingSignal[]>([])
  const [portfolioSummary, setPortfolioSummary] = useState<any>(null)
  const [portfolioMetrics, setPortfolioMetrics] = useState<any>(null)
  const [portfolioRisk, setPortfolioRisk] = useState<any>(null)
  const [strategies, setStrategies] = useState<any[]>([])
  const [riskLimits, setRiskLimits] = useState<any[]>([])
  const [riskViolations, setRiskViolations] = useState<any[]>([])
  const [riskMetrics, setRiskMetrics] = useState<any>(null)
  const [systemStatus, setSystemStatus] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // WebSocket for real-time updates
  const { isConnected, lastMessage } = useWebSocket('ws://localhost:8080/ws/trading')

  // Handle WebSocket messages
  useEffect(() => {
    if (lastMessage) {
      try {
        const { type, data } = lastMessage
        
        switch (type) {
          case 'hft_metrics':
            setHftMetrics(data)
            break
          case 'market_data':
            setMarketData(prev => {
              const updated = [...prev]
              const index = updated.findIndex(item => item.symbol === data.symbol)
              if (index >= 0) {
                updated[index] = data
              } else {
                updated.push(data)
              }
              return updated
            })
            break
          case 'new_order':
            setOrders(prev => [data, ...prev.slice(0, 99)]) // Keep last 100 orders
            break
          case 'order_update':
            setOrders(prev => prev.map(order => 
              order.id === data.id ? { ...order, ...data } : order
            ))
            break
          case 'new_signal':
            setSignals(prev => [data, ...prev.slice(0, 49)]) // Keep last 50 signals
            break
          case 'portfolio_update':
            setPortfolioSummary(data)
            break
          case 'strategy_update':
            setStrategies(prev => prev.map(strategy => 
              strategy.id === data.id ? { ...strategy, ...data } : strategy
            ))
            break
          case 'risk_violation':
            setRiskViolations(prev => [data, ...prev])
            break
          default:
            console.log('Unknown WebSocket message type:', type)
        }
      } catch (error) {
        console.error('Error processing WebSocket message:', error)
      }
    }
  }, [lastMessage])

  // API calls
  const fetchHFTData = useCallback(async () => {
    try {
      const [status, metrics] = await Promise.all([
        tradingApi.getHFTStatus(),
        tradingApi.getHFTMetrics(),
      ])
      setHftStatus(status)
      setHftMetrics(metrics as HFTMetrics)
      setIsHFTRunning((status as any)?.isRunning || false)
    } catch (error) {
      console.error('Error fetching HFT data:', error)
    }
  }, [])

  const fetchTradingData = useCallback(async () => {
    try {
      const [ordersData, positionsData, signalsData] = await Promise.all([
        tradingApi.getOrders(50),
        tradingApi.getPositions(),
        tradingApi.getSignals(),
      ])
      setOrders((ordersData as any)?.orders || [])
      setPositions((positionsData as any)?.positions || [])
      setSignals((signalsData as any)?.signals || [])
    } catch (error) {
      console.error('Error fetching trading data:', error)
    }
  }, [])

  const fetchPortfolioData = useCallback(async () => {
    try {
      const [summary, metrics, risk] = await Promise.all([
        tradingApi.getPortfolioSummary(),
        tradingApi.getPortfolioMetrics(),
        tradingApi.getPortfolioRisk(),
      ])
      setPortfolioSummary(summary)
      setPortfolioMetrics(metrics)
      setPortfolioRisk(risk)
    } catch (error) {
      console.error('Error fetching portfolio data:', error)
    }
  }, [])

  const fetchStrategiesData = useCallback(async () => {
    try {
      const strategiesData = await tradingApi.getStrategies()
      setStrategies((strategiesData as any)?.strategies || [])
    } catch (error) {
      console.error('Error fetching strategies data:', error)
    }
  }, [])

  const fetchRiskData = useCallback(async () => {
    try {
      const [limits, violations, metrics] = await Promise.all([
        tradingApi.getRiskLimits(),
        tradingApi.getRiskViolations(),
        tradingApi.getRiskMetrics(),
      ])
      setRiskLimits((limits as any)?.limits || [])
      setRiskViolations((violations as any)?.violations || [])
      setRiskMetrics(metrics)
    } catch (error) {
      console.error('Error fetching risk data:', error)
    }
  }, [])

  const fetchSystemData = useCallback(async () => {
    try {
      const status = await tradingApi.getSystemStatus()
      setSystemStatus(status)
    } catch (error) {
      console.error('Error fetching system data:', error)
    }
  }, [])

  // Refresh all data
  const refreshData = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    
    try {
      await Promise.all([
        fetchHFTData(),
        fetchTradingData(),
        fetchPortfolioData(),
        fetchStrategiesData(),
        fetchRiskData(),
        fetchSystemData(),
      ])
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Unknown error occurred')
    } finally {
      setIsLoading(false)
    }
  }, [fetchHFTData, fetchTradingData, fetchPortfolioData, fetchStrategiesData, fetchRiskData, fetchSystemData])

  // HFT Engine controls
  const startHFTEngine = useCallback(async () => {
    try {
      await tradingApi.startHFTEngine()
      setIsHFTRunning(true)
      await fetchHFTData()
    } catch (error) {
      console.error('Error starting HFT engine:', error)
      setError('Failed to start HFT engine')
    }
  }, [fetchHFTData])

  const stopHFTEngine = useCallback(async () => {
    try {
      await tradingApi.stopHFTEngine()
      setIsHFTRunning(false)
      await fetchHFTData()
    } catch (error) {
      console.error('Error stopping HFT engine:', error)
      setError('Failed to stop HFT engine')
    }
  }, [fetchHFTData])

  // Emergency stop
  const emergencyStop = useCallback(async (reason: string) => {
    try {
      await tradingApi.emergencyStop(reason, {
        stop_trading: true,
        cancel_orders: true,
        close_positions: true,
      })
      setIsHFTRunning(false)
      await refreshData()
    } catch (error) {
      console.error('Error executing emergency stop:', error)
      setError('Failed to execute emergency stop')
    }
  }, [refreshData])

  // Initial data load
  useEffect(() => {
    refreshData()
  }, [refreshData])

  // Auto-refresh data every 30 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      if (!isLoading) {
        refreshData()
      }
    }, 30000)

    return () => clearInterval(interval)
  }, [refreshData, isLoading])

  return {
    // HFT Engine
    hftMetrics,
    hftStatus,
    isHFTRunning,
    startHFTEngine,
    stopHFTEngine,
    
    // Market Data
    marketData,
    
    // Trading
    orders,
    positions,
    signals,
    
    // Portfolio
    portfolioSummary,
    portfolioMetrics,
    portfolioRisk,
    
    // Strategies
    strategies,
    
    // Risk Management
    riskLimits,
    riskViolations,
    riskMetrics,
    
    // System
    systemStatus,

    // Connection
    isConnected,

    // Loading states
    isLoading,
    error,
    
    // Actions
    refreshData,
    emergencyStop,
  }
}
