import { useState, useEffect, useCallback, useRef } from 'react'
import { 
  solanaRealtimeData, 
  PriceUpdate, 
  DeFiUpdate, 
  NFTUpdate, 
  NetworkUpdate 
} from '@/lib/solana/realtime-data'

export interface RealtimeState<T> {
  data: T[]
  isConnected: boolean
  lastUpdate: Date | null
  error: string | null
}

export interface UseSolanaRealtimeOptions {
  enabled?: boolean
  maxDataPoints?: number
  onUpdate?: (data: any) => void
  onError?: (error: any) => void
}

// Hook for real-time price updates
export function useSolanaPriceRealtime(
  symbols: string[], 
  options: UseSolanaRealtimeOptions = {}
) {
  const {
    enabled = true,
    maxDataPoints = 100,
    onUpdate,
    onError
  } = options

  const [state, setState] = useState<RealtimeState<PriceUpdate>>({
    data: [],
    isConnected: false,
    lastUpdate: null,
    error: null
  })

  const subscriptionIdRef = useRef<string | null>(null)

  const handlePriceUpdate = useCallback((update: PriceUpdate) => {
    setState(prev => {
      const newData = [...prev.data, update]
      // Keep only the latest maxDataPoints
      if (newData.length > maxDataPoints) {
        newData.splice(0, newData.length - maxDataPoints)
      }
      
      return {
        ...prev,
        data: newData,
        isConnected: true,
        lastUpdate: update.timestamp,
        error: null
      }
    })
    
    onUpdate?.(update)
  }, [maxDataPoints, onUpdate])

  const handleError = useCallback((error: any) => {
    setState(prev => ({
      ...prev,
      isConnected: false,
      error: error.message || 'Connection error'
    }))
    
    onError?.(error)
  }, [onError])

  useEffect(() => {
    if (!enabled || symbols.length === 0) return

    try {
      subscriptionIdRef.current = solanaRealtimeData.subscribeToPrices(
        symbols,
        handlePriceUpdate
      )

      // Listen for errors
      solanaRealtimeData.on('error', handleError)

      setState(prev => ({ ...prev, isConnected: true, error: null }))
    } catch (error) {
      handleError(error)
    }

    return () => {
      if (subscriptionIdRef.current) {
        solanaRealtimeData.unsubscribe(subscriptionIdRef.current)
        subscriptionIdRef.current = null
      }
      solanaRealtimeData.off('error', handleError)
    }
  }, [enabled, symbols.join(','), handlePriceUpdate, handleError])

  const getLatestPrice = useCallback((symbol: string): PriceUpdate | null => {
    const symbolData = state.data.filter(d => d.symbol === symbol)
    return symbolData.length > 0 ? symbolData[symbolData.length - 1] : null
  }, [state.data])

  const getPriceHistory = useCallback((symbol: string): PriceUpdate[] => {
    return state.data.filter(d => d.symbol === symbol)
  }, [state.data])

  return {
    ...state,
    getLatestPrice,
    getPriceHistory,
    symbols
  }
}

// Hook for real-time DeFi updates
export function useSolanaDeFiRealtime(
  protocols: string[], 
  options: UseSolanaRealtimeOptions = {}
) {
  const {
    enabled = true,
    maxDataPoints = 50,
    onUpdate,
    onError
  } = options

  const [state, setState] = useState<RealtimeState<DeFiUpdate>>({
    data: [],
    isConnected: false,
    lastUpdate: null,
    error: null
  })

  const subscriptionIdRef = useRef<string | null>(null)

  const handleDeFiUpdate = useCallback((update: DeFiUpdate) => {
    setState(prev => {
      const newData = [...prev.data, update]
      if (newData.length > maxDataPoints) {
        newData.splice(0, newData.length - maxDataPoints)
      }
      
      return {
        ...prev,
        data: newData,
        isConnected: true,
        lastUpdate: update.timestamp,
        error: null
      }
    })
    
    onUpdate?.(update)
  }, [maxDataPoints, onUpdate])

  const handleError = useCallback((error: any) => {
    setState(prev => ({
      ...prev,
      isConnected: false,
      error: error.message || 'Connection error'
    }))
    
    onError?.(error)
  }, [onError])

  useEffect(() => {
    if (!enabled || protocols.length === 0) return

    try {
      subscriptionIdRef.current = solanaRealtimeData.subscribeToDeFi(
        protocols,
        handleDeFiUpdate
      )

      solanaRealtimeData.on('error', handleError)
      setState(prev => ({ ...prev, isConnected: true, error: null }))
    } catch (error) {
      handleError(error)
    }

    return () => {
      if (subscriptionIdRef.current) {
        solanaRealtimeData.unsubscribe(subscriptionIdRef.current)
        subscriptionIdRef.current = null
      }
      solanaRealtimeData.off('error', handleError)
    }
  }, [enabled, protocols.join(','), handleDeFiUpdate, handleError])

  const getLatestProtocolData = useCallback((protocol: string): DeFiUpdate | null => {
    const protocolData = state.data.filter(d => d.protocol === protocol)
    return protocolData.length > 0 ? protocolData[protocolData.length - 1] : null
  }, [state.data])

  return {
    ...state,
    getLatestProtocolData,
    protocols
  }
}

// Hook for real-time NFT updates
export function useSolanaNFTRealtime(
  collections: string[], 
  options: UseSolanaRealtimeOptions = {}
) {
  const {
    enabled = true,
    maxDataPoints = 50,
    onUpdate,
    onError
  } = options

  const [state, setState] = useState<RealtimeState<NFTUpdate>>({
    data: [],
    isConnected: false,
    lastUpdate: null,
    error: null
  })

  const subscriptionIdRef = useRef<string | null>(null)

  const handleNFTUpdate = useCallback((update: NFTUpdate) => {
    setState(prev => {
      const newData = [...prev.data, update]
      if (newData.length > maxDataPoints) {
        newData.splice(0, newData.length - maxDataPoints)
      }
      
      return {
        ...prev,
        data: newData,
        isConnected: true,
        lastUpdate: update.timestamp,
        error: null
      }
    })
    
    onUpdate?.(update)
  }, [maxDataPoints, onUpdate])

  const handleError = useCallback((error: any) => {
    setState(prev => ({
      ...prev,
      isConnected: false,
      error: error.message || 'Connection error'
    }))
    
    onError?.(error)
  }, [onError])

  useEffect(() => {
    if (!enabled || collections.length === 0) return

    try {
      subscriptionIdRef.current = solanaRealtimeData.subscribeToNFTs(
        collections,
        handleNFTUpdate
      )

      solanaRealtimeData.on('error', handleError)
      setState(prev => ({ ...prev, isConnected: true, error: null }))
    } catch (error) {
      handleError(error)
    }

    return () => {
      if (subscriptionIdRef.current) {
        solanaRealtimeData.unsubscribe(subscriptionIdRef.current)
        subscriptionIdRef.current = null
      }
      solanaRealtimeData.off('error', handleError)
    }
  }, [enabled, collections.join(','), handleNFTUpdate, handleError])

  const getLatestCollectionData = useCallback((collection: string): NFTUpdate | null => {
    const collectionData = state.data.filter(d => d.collection === collection)
    return collectionData.length > 0 ? collectionData[collectionData.length - 1] : null
  }, [state.data])

  return {
    ...state,
    getLatestCollectionData,
    collections
  }
}

// Hook for real-time network updates
export function useSolanaNetworkRealtime(options: UseSolanaRealtimeOptions = {}) {
  const {
    enabled = true,
    maxDataPoints = 20,
    onUpdate,
    onError
  } = options

  const [state, setState] = useState<RealtimeState<NetworkUpdate>>({
    data: [],
    isConnected: false,
    lastUpdate: null,
    error: null
  })

  const subscriptionIdRef = useRef<string | null>(null)

  const handleNetworkUpdate = useCallback((update: NetworkUpdate) => {
    setState(prev => {
      const newData = [...prev.data, update]
      if (newData.length > maxDataPoints) {
        newData.splice(0, newData.length - maxDataPoints)
      }
      
      return {
        ...prev,
        data: newData,
        isConnected: true,
        lastUpdate: update.timestamp,
        error: null
      }
    })
    
    onUpdate?.(update)
  }, [maxDataPoints, onUpdate])

  const handleError = useCallback((error: any) => {
    setState(prev => ({
      ...prev,
      isConnected: false,
      error: error.message || 'Connection error'
    }))
    
    onError?.(error)
  }, [onError])

  useEffect(() => {
    if (!enabled) return

    try {
      subscriptionIdRef.current = solanaRealtimeData.subscribeToNetwork(handleNetworkUpdate)
      solanaRealtimeData.on('error', handleError)
      setState(prev => ({ ...prev, isConnected: true, error: null }))
    } catch (error) {
      handleError(error)
    }

    return () => {
      if (subscriptionIdRef.current) {
        solanaRealtimeData.unsubscribe(subscriptionIdRef.current)
        subscriptionIdRef.current = null
      }
      solanaRealtimeData.off('error', handleError)
    }
  }, [enabled, handleNetworkUpdate, handleError])

  const getLatestNetworkData = useCallback((): NetworkUpdate | null => {
    return state.data.length > 0 ? state.data[state.data.length - 1] : null
  }, [state.data])

  const getAverageTPS = useCallback((): number => {
    if (state.data.length === 0) return 0
    const sum = state.data.reduce((acc, update) => acc + update.tps, 0)
    return sum / state.data.length
  }, [state.data])

  return {
    ...state,
    getLatestNetworkData,
    getAverageTPS
  }
}
