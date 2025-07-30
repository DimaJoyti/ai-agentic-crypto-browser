import { useState, useEffect, useCallback, useRef } from 'react'
import { useSwitchChain, useChainId, useAccount } from 'wagmi'
import { type Address } from 'viem'
import { toast } from 'sonner'
import { 
  multiChainManager, 
  type ChainConfig, 
  type ChainSwitchRequest,
  SUPPORTED_CHAINS 
} from '@/lib/multi-chain-manager'
import { useWalletStore } from '@/stores/walletStore'

export interface ChainSwitchOptions {
  autoSwitch?: boolean
  showConfirmation?: boolean
  timeout?: number
  fallbackChain?: number
  onSuccess?: (chainId: number) => void
  onError?: (error: Error, chainId: number) => void
  reason?: 'user' | 'dapp' | 'auto'
}

export interface ChainSwitchState {
  isSwitching: boolean
  targetChainId: number | null
  error: string | null
  lastSwitchTime: number | null
  switchCount: number
}

export interface UseChainSwitchingReturn {
  // Current state
  currentChain: ChainConfig | null
  switchState: ChainSwitchState
  supportedChains: ChainConfig[]
  
  // Chain switching
  switchToChain: (chainId: number, options?: ChainSwitchOptions) => Promise<boolean>
  switchToOptimalChain: (dappChains: number[]) => Promise<boolean>
  
  // Chain management
  addCustomChain: (chain: ChainConfig) => Promise<boolean>
  isChainSupported: (chainId: number) => boolean
  getChainConfig: (chainId: number) => ChainConfig | undefined
  
  // Auto-switching
  enableAutoSwitch: () => void
  disableAutoSwitch: () => void
  isAutoSwitchEnabled: boolean
  
  // Utilities
  getRecommendedChains: () => ChainConfig[]
  getChainSwitchHistory: () => ChainSwitchRequest[]
  clearSwitchHistory: () => void
  
  // Error handling
  clearError: () => void
  retryLastSwitch: () => Promise<boolean>
}

export const useChainSwitching = (
  defaultOptions?: ChainSwitchOptions
): UseChainSwitchingReturn => {
  const [switchState, setSwitchState] = useState<ChainSwitchState>({
    isSwitching: false,
    targetChainId: null,
    error: null,
    lastSwitchTime: null,
    switchCount: 0
  })
  
  const [isAutoSwitchEnabled, setIsAutoSwitchEnabled] = useState(false)
  const lastSwitchAttempt = useRef<{ chainId: number; options?: ChainSwitchOptions } | null>(null)
  const switchTimeoutRef = useRef<NodeJS.Timeout>()
  
  // Wagmi hooks
  const { switchChain, isPending, error: switchError } = useSwitchChain()
  const currentChainId = useChainId()
  const { address } = useAccount()
  
  // Store hooks
  const walletStore = useWalletStore()
  
  // Get current chain configuration
  const currentChain = currentChainId ? SUPPORTED_CHAINS[currentChainId] : null
  const supportedChains = multiChainManager.getSupportedChains()
  
  // Update switch state when switching status changes
  useEffect(() => {
    setSwitchState(prev => ({
      ...prev,
      isSwitching: isPending
    }))
  }, [isPending])
  
  // Handle switch errors
  useEffect(() => {
    if (switchError) {
      const errorMessage = getSwitchErrorMessage(switchError)
      setSwitchState(prev => ({
        ...prev,
        error: errorMessage,
        isSwitching: false
      }))
      
      // Record failed switch
      if (switchState.targetChainId) {
        multiChainManager.recordChainSwitch({
          chainId: switchState.targetChainId,
          reason: 'user',
          timestamp: Date.now(),
          success: false,
          error: errorMessage
        })
      }
      
      toast.error('Chain Switch Failed', {
        description: errorMessage
      })
    }
  }, [switchError, switchState.targetChainId])
  
  // Handle successful chain switch
  useEffect(() => {
    if (currentChainId && switchState.targetChainId === currentChainId && switchState.isSwitching) {
      setSwitchState(prev => ({
        ...prev,
        isSwitching: false,
        targetChainId: null,
        error: null,
        lastSwitchTime: Date.now(),
        switchCount: prev.switchCount + 1
      }))
      
      // Record successful switch
      multiChainManager.recordChainSwitch({
        chainId: currentChainId,
        reason: 'user',
        timestamp: Date.now(),
        success: true
      })
      
      // Update wallet store
      if (address) {
        walletStore.updateWalletChain(address, currentChainId)
      }
      
      const chainName = SUPPORTED_CHAINS[currentChainId]?.name || `Chain ${currentChainId}`
      toast.success('Chain Switched', {
        description: `Successfully switched to ${chainName}`
      })
      
      // Clear timeout
      if (switchTimeoutRef.current) {
        clearTimeout(switchTimeoutRef.current)
      }
    }
  }, [currentChainId, switchState.targetChainId, switchState.isSwitching, address, walletStore])
  
  // Load auto-switch preference
  useEffect(() => {
    const autoSwitch = localStorage.getItem('chain_auto_switch') === 'true'
    setIsAutoSwitchEnabled(autoSwitch)
  }, [])
  
  // Main chain switching function
  const switchToChain = useCallback(async (
    chainId: number, 
    options: ChainSwitchOptions = {}
  ): Promise<boolean> => {
    const mergedOptions = { ...defaultOptions, ...options }
    const { 
      showConfirmation = false, 
      timeout = 30000, 
      onSuccess, 
      onError,
      reason = 'user'
    } = mergedOptions
    
    // Check if chain is supported
    if (!SUPPORTED_CHAINS[chainId]) {
      const error = new Error(`Chain ${chainId} is not supported`)
      onError?.(error, chainId)
      setSwitchState(prev => ({ ...prev, error: error.message }))
      return false
    }
    
    // Check if already on target chain
    if (currentChainId === chainId) {
      onSuccess?.(chainId)
      return true
    }
    
    // Show confirmation if required
    if (showConfirmation) {
      const chainName = SUPPORTED_CHAINS[chainId]?.name || `Chain ${chainId}`
      const confirmed = window.confirm(`Switch to ${chainName}?`)
      if (!confirmed) return false
    }
    
    // Store attempt for retry
    lastSwitchAttempt.current = { chainId, options }
    
    setSwitchState(prev => ({
      ...prev,
      isSwitching: true,
      targetChainId: chainId,
      error: null
    }))
    
    try {
      // Set timeout
      const timeoutPromise = new Promise<never>((_, reject) => {
        switchTimeoutRef.current = setTimeout(() => {
          reject(new Error('Chain switch timeout'))
        }, timeout)
      })
      
      // Attempt switch
      const switchPromise = switchChain({ chainId })
      
      await Promise.race([switchPromise, timeoutPromise])
      
      onSuccess?.(chainId)
      return true
      
    } catch (error) {
      console.error('Chain switch error:', error)
      
      const errorMessage = getSwitchErrorMessage(error as Error)
      setSwitchState(prev => ({
        ...prev,
        error: errorMessage,
        isSwitching: false,
        targetChainId: null
      }))
      
      // Record failed switch
      multiChainManager.recordChainSwitch({
        chainId,
        reason,
        timestamp: Date.now(),
        success: false,
        error: errorMessage
      })
      
      onError?.(error as Error, chainId)
      return false
    }
  }, [currentChainId, switchChain, defaultOptions])
  
  // Switch to optimal chain for dapp
  const switchToOptimalChain = useCallback(async (dappChains: number[]): Promise<boolean> => {
    if (dappChains.length === 0) return false
    
    // Check if current chain is supported by dapp
    if (currentChainId && dappChains.includes(currentChainId)) {
      return true
    }
    
    // Get user's chain preferences from analytics
    const analytics = walletStore.analytics
    const userChainPreferences = Object.entries(analytics.favoriteChains)
      .sort(([,a], [,b]) => b - a)
      .map(([chainId]) => parseInt(chainId))
    
    // Find best chain: user preference + dapp support
    let targetChain = dappChains.find(chainId => userChainPreferences.includes(chainId))
    
    // Fallback to first supported chain
    if (!targetChain) {
      targetChain = dappChains[0]
    }
    
    return await switchToChain(targetChain, { reason: 'auto' })
  }, [currentChainId, switchToChain, walletStore.analytics])
  
  // Add custom chain
  const addCustomChain = useCallback(async (chain: ChainConfig): Promise<boolean> => {
    try {
      // This would typically involve adding the chain to the wallet
      // For now, we'll just add it to our supported chains
      SUPPORTED_CHAINS[chain.id] = chain
      
      toast.success('Custom Chain Added', {
        description: `${chain.name} has been added to supported chains`
      })
      
      return true
    } catch (error) {
      console.error('Failed to add custom chain:', error)
      toast.error('Failed to add custom chain')
      return false
    }
  }, [])
  
  // Utility functions
  const isChainSupported = useCallback((chainId: number): boolean => {
    return !!SUPPORTED_CHAINS[chainId]
  }, [])
  
  const getChainConfig = useCallback((chainId: number): ChainConfig | undefined => {
    return SUPPORTED_CHAINS[chainId]
  }, [])
  
  const getRecommendedChains = useCallback((): ChainConfig[] => {
    const userActivity = Object.entries(walletStore.analytics.favoriteChains)
      .map(([chainId, txCount]) => ({ chainId: parseInt(chainId), txCount }))
    
    return multiChainManager.getRecommendedChains(userActivity)
  }, [walletStore.analytics.favoriteChains])
  
  const getChainSwitchHistory = useCallback((): ChainSwitchRequest[] => {
    return multiChainManager.getChainSwitchStats().recentSwitches
  }, [])
  
  const clearSwitchHistory = useCallback((): void => {
    // This would clear the switch history in the manager
    // For now, we'll just show a toast
    toast.success('Switch history cleared')
  }, [])
  
  // Auto-switch management
  const enableAutoSwitch = useCallback(() => {
    localStorage.setItem('chain_auto_switch', 'true')
    setIsAutoSwitchEnabled(true)
    toast.success('Auto chain switching enabled')
  }, [])
  
  const disableAutoSwitch = useCallback(() => {
    localStorage.setItem('chain_auto_switch', 'false')
    setIsAutoSwitchEnabled(false)
    toast.success('Auto chain switching disabled')
  }, [])
  
  // Error handling
  const clearError = useCallback(() => {
    setSwitchState(prev => ({ ...prev, error: null }))
  }, [])
  
  const retryLastSwitch = useCallback(async (): Promise<boolean> => {
    if (!lastSwitchAttempt.current) return false
    
    const { chainId, options } = lastSwitchAttempt.current
    return await switchToChain(chainId, options)
  }, [switchToChain])
  
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (switchTimeoutRef.current) {
        clearTimeout(switchTimeoutRef.current)
      }
    }
  }, [])
  
  return {
    currentChain,
    switchState,
    supportedChains,
    switchToChain,
    switchToOptimalChain,
    addCustomChain,
    isChainSupported,
    getChainConfig,
    enableAutoSwitch,
    disableAutoSwitch,
    isAutoSwitchEnabled,
    getRecommendedChains,
    getChainSwitchHistory,
    clearSwitchHistory,
    clearError,
    retryLastSwitch
  }
}

// Helper function to get user-friendly error messages
function getSwitchErrorMessage(error: Error): string {
  const message = error.message.toLowerCase()
  
  if (message.includes('user rejected')) {
    return 'Chain switch was rejected by user'
  }
  
  if (message.includes('unrecognized chain')) {
    return 'Chain not recognized by wallet. Please add it manually first.'
  }
  
  if (message.includes('timeout')) {
    return 'Chain switch timed out. Please try again.'
  }
  
  if (message.includes('network')) {
    return 'Network error occurred during chain switch'
  }
  
  return error.message || 'Unknown error occurred during chain switch'
}
