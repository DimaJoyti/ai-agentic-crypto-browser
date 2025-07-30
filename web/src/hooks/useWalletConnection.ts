import { useState, useEffect, useCallback, useRef } from 'react'
import { useConnect, useDisconnect, useAccount, useChainId, useSwitchChain } from 'wagmi'
import { type Connector } from 'wagmi'
import { type Address } from 'viem'
import { toast } from 'sonner'
import {
  type WalletProvider,
  type WalletConnectionState,
  type WalletConnectionOptions,
  WalletConnectionError,
  WALLET_ERROR_CODES,
  getWalletErrorMessage,
  getWalletByConnector,
  saveWalletPreference,
  getWalletPreference,
  clearWalletPreference,
  WALLET_PROVIDERS
} from '@/lib/wallet-providers'
import { useWalletStore, type ConnectedWallet } from '@/stores/walletStore'
import { useSessionRecovery } from '@/lib/session-recovery'
import { useWalletPersistence } from '@/lib/wallet-persistence'

export interface UseWalletConnectionReturn {
  // Connection state
  connectionState: WalletConnectionState
  
  // Connection methods
  connectWallet: (walletId: string, options?: WalletConnectionOptions) => Promise<void>
  disconnectWallet: () => Promise<void>
  reconnectWallet: () => Promise<void>
  
  // Wallet management
  availableWallets: WalletProvider[]
  connectedWallet: WalletProvider | null
  
  // Chain management
  switchChain: (chainId: number) => Promise<void>
  
  // Error handling
  clearError: () => void
  
  // Auto-connect
  enableAutoConnect: () => void
  disableAutoConnect: () => void
  isAutoConnectEnabled: boolean
}

export const useWalletConnection = (
  defaultOptions?: WalletConnectionOptions
): UseWalletConnectionReturn => {
  const [connectionState, setConnectionState] = useState<WalletConnectionState>({
    isConnecting: false,
    isConnected: false
  })

  const connectionTimeoutRef = useRef<NodeJS.Timeout>()
  const retryCountRef = useRef(0)

  // Wagmi hooks
  const { connect, connectors, isPending, error: connectError } = useConnect()
  const { disconnect } = useDisconnect()
  const { address, isConnected, connector } = useAccount()
  const chainId = useChainId()
  const { switchChain: wagmiSwitchChain } = useSwitchChain()

  // Store hooks
  const walletStore = useWalletStore()
  const sessionRecovery = useSessionRecovery()
  const persistence = useWalletPersistence()

  // Get auto-connect setting from store
  const isAutoConnectEnabled = walletStore.preferences.autoConnect
  
  // Update connection state when account changes
  useEffect(() => {
    const newState = {
      isConnected,
      address: address as Address,
      chainId,
      connector,
      isConnecting: isPending
    }

    setConnectionState(prev => ({ ...prev, ...newState }))
    walletStore.setConnectionState(newState)

    // Update current wallet in store
    if (isConnected && address && connector) {
      const walletProvider = getWalletByConnector(connector)
      if (walletProvider) {
        const connectedWallet: ConnectedWallet = {
          id: `${walletProvider.id}_${address}`,
          address: address as Address,
          chainId: chainId || 1,
          connector: connector.name,
          provider: walletProvider,
          connectedAt: Date.now(),
          lastUsed: Date.now(),
          isActive: true
        }

        walletStore.setCurrentWallet(connectedWallet)
        walletStore.addConnectedWallet(connectedWallet)
        walletStore.startSession(connectedWallet)

        // Save session for recovery
        sessionRecovery.saveCurrentSession()

        // Update analytics
        walletStore.incrementConnections(walletProvider.id)
      }
    } else if (!isConnected) {
      walletStore.setCurrentWallet(null)
      walletStore.endSession()
    }
  }, [isConnected, address, chainId, connector, isPending, walletStore, sessionRecovery])
  
  // Handle connection errors
  useEffect(() => {
    if (connectError) {
      const errorMessage = getWalletErrorMessage(connectError)
      setConnectionState(prev => ({
        ...prev,
        error: errorMessage,
        isConnecting: false
      }))
      
      toast.error('Wallet Connection Failed', {
        description: errorMessage
      })
    }
  }, [connectError])
  
  // Auto-connect and session recovery on mount
  useEffect(() => {
    const initializeConnection = async () => {
      if (isConnected) return // Already connected

      // First try session recovery
      if (sessionRecovery.hasRecoverableSession()) {
        try {
          const recovered = await sessionRecovery.recoverSession(
            async (walletId) => {
              await connectWallet(walletId, { autoConnect: true })
            },
            {
              requireUserConfirmation: false,
              autoReconnect: isAutoConnectEnabled
            }
          )

          if (recovered) {
            toast.success('Session recovered successfully')
            return
          }
        } catch (error) {
          console.warn('Session recovery failed:', error)
        }
      }

      // Fallback to preferred wallet auto-connect
      if (isAutoConnectEnabled && walletStore.preferences.preferredWallet) {
        try {
          await connectWallet(walletStore.preferences.preferredWallet, { autoConnect: true })
        } catch (error) {
          console.warn('Auto-connect failed:', error)
        }
      }
    }

    initializeConnection()
  }, []) // Only run on mount
  
  // Get available wallets
  const availableWallets = WALLET_PROVIDERS.filter(wallet => {
    // Check if wallet is supported and has a corresponding connector
    const hasConnector = connectors.some(connector =>
      connector.id.toLowerCase().includes(wallet.id) ||
      connector.name.toLowerCase().includes(wallet.name.toLowerCase())
    )
    return wallet.supported && hasConnector
  })

  // Get currently connected wallet from store
  const connectedWallet = walletStore.currentWallet?.provider || null
  
  // Connect wallet function
  const connectWallet = useCallback(async (
    walletId: string, 
    options: WalletConnectionOptions = {}
  ) => {
    const mergedOptions = { ...defaultOptions, ...options }
    const { timeout = 30000, retryAttempts = 3 } = mergedOptions
    
    setConnectionState(prev => ({
      ...prev,
      isConnecting: true,
      error: undefined
    }))
    
    // Clear any existing timeout
    if (connectionTimeoutRef.current) {
      clearTimeout(connectionTimeoutRef.current)
    }
    
    try {
      // Find the appropriate connector
      const targetConnector = connectors.find(c => 
        c.id.toLowerCase().includes(walletId) ||
        c.name.toLowerCase().includes(walletId)
      )
      
      if (!targetConnector) {
        throw new WalletConnectionError(
          `Connector for ${walletId} not found`,
          WALLET_ERROR_CODES.WALLET_NOT_FOUND,
          walletId
        )
      }
      
      // Set connection timeout
      const timeoutPromise = new Promise<never>((_, reject) => {
        connectionTimeoutRef.current = setTimeout(() => {
          reject(new WalletConnectionError(
            'Connection timeout',
            WALLET_ERROR_CODES.TIMEOUT,
            walletId
          ))
        }, timeout)
      })
      
      // Attempt connection
      const connectPromise = connect({ connector: targetConnector })
      
      await Promise.race([connectPromise, timeoutPromise])
      
      // Clear timeout on success
      if (connectionTimeoutRef.current) {
        clearTimeout(connectionTimeoutRef.current)
      }
      
      // Save wallet preference
      saveWalletPreference(walletId)
      walletStore.setPreferredWallet(walletId)

      // Reset retry count
      retryCountRef.current = 0

      // Call success callback
      if (mergedOptions.onSuccess && address && chainId) {
        mergedOptions.onSuccess(address, chainId)
      }

      // Create backup after successful connection
      try {
        await persistence.createBackup()
      } catch (error) {
        console.warn('Failed to create backup after connection:', error)
      }

      toast.success('Wallet Connected', {
        description: `Successfully connected to ${targetConnector.name}`
      })
      
    } catch (error) {
      console.error('Wallet connection error:', error)
      
      const walletError = error instanceof WalletConnectionError 
        ? error 
        : new WalletConnectionError(
            getWalletErrorMessage(error as Error, walletId),
            WALLET_ERROR_CODES.UNKNOWN,
            walletId,
            error as Error
          )
      
      // Retry logic
      if (retryCountRef.current < retryAttempts && !options.autoConnect) {
        retryCountRef.current++
        toast.warning(`Connection failed, retrying... (${retryCountRef.current}/${retryAttempts})`)
        
        setTimeout(() => {
          connectWallet(walletId, options)
        }, 1000 * retryCountRef.current)
        return
      }
      
      setConnectionState(prev => ({
        ...prev,
        error: walletError.message,
        isConnecting: false
      }))
      
      // Call error callback
      if (mergedOptions.onError) {
        mergedOptions.onError(walletError)
      }
      
      throw walletError
    }
  }, [connectors, connect, address, chainId, defaultOptions])
  
  // Disconnect wallet function
  const disconnectWallet = useCallback(async () => {
    try {
      await disconnect()
      clearWalletPreference()

      // Clear store state
      walletStore.setCurrentWallet(null)
      walletStore.endSession()

      // Clear session recovery data
      if (walletStore.currentWallet) {
        sessionRecovery.clearSessions()
      }

      setConnectionState(prev => ({
        ...prev,
        error: undefined
      }))

      toast.success('Wallet Disconnected')

      // Call disconnect callback
      if (defaultOptions?.onDisconnect) {
        defaultOptions.onDisconnect()
      }

    } catch (error) {
      console.error('Disconnect error:', error)
      toast.error('Failed to disconnect wallet')
    }
  }, [disconnect, defaultOptions, walletStore, sessionRecovery])
  
  // Reconnect wallet function
  const reconnectWallet = useCallback(async () => {
    const preferredWallet = getWalletPreference()
    if (preferredWallet) {
      await connectWallet(preferredWallet)
    }
  }, [connectWallet])
  
  // Switch chain function
  const switchChain = useCallback(async (targetChainId: number) => {
    try {
      await wagmiSwitchChain({ chainId: targetChainId })
      toast.success('Network switched successfully')
    } catch (error) {
      console.error('Chain switch error:', error)
      toast.error('Failed to switch network')
      throw error
    }
  }, [wagmiSwitchChain])
  
  // Clear error function
  const clearError = useCallback(() => {
    setConnectionState(prev => ({
      ...prev,
      error: undefined
    }))
  }, [])
  
  // Auto-connect management
  const enableAutoConnect = useCallback(() => {
    walletStore.setAutoConnect(true)
  }, [walletStore])

  const disableAutoConnect = useCallback(() => {
    walletStore.setAutoConnect(false)
  }, [walletStore])
  
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (connectionTimeoutRef.current) {
        clearTimeout(connectionTimeoutRef.current)
      }
    }
  }, [])
  
  return {
    connectionState,
    connectWallet,
    disconnectWallet,
    reconnectWallet,
    availableWallets,
    connectedWallet,
    switchChain,
    clearError,
    enableAutoConnect,
    disableAutoConnect,
    isAutoConnectEnabled
  }
}
