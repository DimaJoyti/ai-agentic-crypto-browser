import { useState, useEffect, useCallback } from 'react'
import { SolanaDeFiService } from '@/services/solana/SolanaDeFiService'

export interface DeFiProtocol {
  id: string
  name: string
  tvl: number
  tvlChange24h: number
  category: 'dex' | 'lending' | 'yield' | 'derivatives' | 'insurance'
  apy?: number
  volume24h: number
  users24h: number
  logo: string
  website: string
  description: string
}

export interface YieldOpportunity {
  protocol: string
  pool: string
  apy: number
  tvl: number
  risk: 'low' | 'medium' | 'high'
  tokens: string[]
  minimumDeposit: number
  lockPeriod?: number
}

export interface DeFiStats {
  totalTVL: number
  totalVolume24h: number
  totalProtocols: number
  avgAPY: number
  topProtocolByTVL: string
  topProtocolByVolume: string
}

export interface SolanaDeFiState {
  protocols: DeFiProtocol[]
  yields: YieldOpportunity[]
  stats: DeFiStats | null
  isLoading: boolean
  error: string | null
  lastUpdated: Date | null
}

export interface UseSolanaDeFiOptions {
  autoRefresh?: boolean
  refreshInterval?: number
  category?: DeFiProtocol['category']
  minTVL?: number
  minAPY?: number
}

export function useSolanaDeFi(options: UseSolanaDeFiOptions = {}) {
  const {
    autoRefresh = false,
    refreshInterval = 60000, // 1 minute for DeFi data
    category,
    minTVL = 0,
    minAPY = 0
  } = options

  const [state, setState] = useState<SolanaDeFiState>({
    protocols: [],
    yields: [],
    stats: null,
    isLoading: true,
    error: null,
    lastUpdated: null
  })

  const defiService = new SolanaDeFiService()

  const fetchDeFiData = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }))

      const [protocolsData, yieldsData, statsData] = await Promise.all([
        defiService.getProtocols({ category, minTVL }),
        defiService.getYieldOpportunities({ minAPY }),
        defiService.getDeFiStats()
      ])

      setState(prev => ({
        ...prev,
        protocols: protocolsData,
        yields: yieldsData,
        stats: statsData,
        isLoading: false,
        lastUpdated: new Date()
      }))
    } catch (error) {
      console.error('Failed to fetch Solana DeFi data:', error)
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Failed to fetch DeFi data'
      }))
    }
  }, [category, minTVL, minAPY])

  const refresh = useCallback(async () => {
    await fetchDeFiData()
  }, [fetchDeFiData])

  // Initial data fetch
  useEffect(() => {
    fetchDeFiData()
  }, [fetchDeFiData])

  // Auto-refresh effect
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(fetchDeFiData, refreshInterval)
    return () => clearInterval(interval)
  }, [autoRefresh, refreshInterval, fetchDeFiData])

  // Computed values
  const totalTVL = state.stats?.totalTVL || 0
  const topYields = state.yields
    .sort((a, b) => b.apy - a.apy)
    .slice(0, 10)
  
  const protocolsByCategory = state.protocols.reduce((acc, protocol) => {
    if (!acc[protocol.category]) {
      acc[protocol.category] = []
    }
    acc[protocol.category].push(protocol)
    return acc
  }, {} as Record<DeFiProtocol['category'], DeFiProtocol[]>)

  const topProtocolsByTVL = state.protocols
    .sort((a, b) => b.tvl - a.tvl)
    .slice(0, 10)

  return {
    ...state,
    refresh,
    totalTVL,
    topYields,
    protocolsByCategory,
    topProtocolsByTVL
  }
}

// Helper hook for specific protocol data
export function useDeFiProtocol(protocolId: string) {
  const [protocol, setProtocol] = useState<DeFiProtocol | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const defiService = new SolanaDeFiService()

  const fetchProtocol = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const protocolData = await defiService.getProtocol(protocolId)
      setProtocol(protocolData)
    } catch (error) {
      console.error(`Failed to fetch protocol ${protocolId}:`, error)
      setError(error instanceof Error ? error.message : 'Failed to fetch protocol')
    } finally {
      setIsLoading(false)
    }
  }, [protocolId])

  useEffect(() => {
    if (protocolId) {
      fetchProtocol()
    }
  }, [protocolId, fetchProtocol])

  return {
    protocol,
    isLoading,
    error,
    refresh: fetchProtocol
  }
}

// Helper hook for yield farming opportunities
export function useYieldFarming(options: { minAPY?: number; maxRisk?: YieldOpportunity['risk'] } = {}) {
  const { yields, isLoading, error, refresh } = useSolanaDeFi({
    autoRefresh: true,
    refreshInterval: 300000, // 5 minutes
    minAPY: options.minAPY
  })

  const filteredYields = yields.filter(yield_ => {
    if (options.maxRisk) {
      const riskLevels = { low: 1, medium: 2, high: 3 }
      return riskLevels[yield_.risk] <= riskLevels[options.maxRisk]
    }
    return true
  })

  const yieldsByRisk = filteredYields.reduce((acc, yield_) => {
    if (!acc[yield_.risk]) {
      acc[yield_.risk] = []
    }
    acc[yield_.risk].push(yield_)
    return acc
  }, {} as Record<YieldOpportunity['risk'], YieldOpportunity[]>)

  return {
    yields: filteredYields,
    yieldsByRisk,
    isLoading,
    error,
    refresh
  }
}

// Helper hook for DEX data
export function useSolanaDEX() {
  const { protocolsByCategory, isLoading, error, refresh } = useSolanaDeFi({
    category: 'dex',
    autoRefresh: true,
    refreshInterval: 30000 // 30 seconds for DEX data
  })

  const dexProtocols = protocolsByCategory.dex || []
  
  const totalDEXVolume = dexProtocols.reduce((sum, protocol) => sum + protocol.volume24h, 0)
  const totalDEXTVL = dexProtocols.reduce((sum, protocol) => sum + protocol.tvl, 0)

  return {
    dexProtocols,
    totalDEXVolume,
    totalDEXTVL,
    isLoading,
    error,
    refresh
  }
}
