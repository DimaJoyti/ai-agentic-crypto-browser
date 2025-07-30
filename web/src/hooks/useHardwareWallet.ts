import { useState, useEffect, useCallback, useRef } from 'react'
import { type Address, type Hex } from 'viem'
import { toast } from 'sonner'
import {
  hardwareWalletManager,
  type HardwareWalletDevice,
  type HardwareWalletAccount,
  type HardwareWalletConnection,
  type SigningRequest
} from '@/lib/hardware-wallet-manager'

export interface HardwareWalletState {
  isScanning: boolean
  isConnecting: boolean
  isSigning: boolean
  devices: HardwareWalletDevice[]
  connectedDevices: HardwareWalletDevice[]
  selectedDevice: HardwareWalletDevice | null
  accounts: HardwareWalletAccount[]
  selectedAccount: HardwareWalletAccount | null
  signingRequests: SigningRequest[]
  error: string | null
}

export interface HardwareWalletOptions {
  autoScan?: boolean
  autoConnect?: boolean
  accountCount?: number
  onDeviceConnect?: (device: HardwareWalletDevice) => void
  onDeviceDisconnect?: (device: HardwareWalletDevice) => void
  onAccountSelect?: (account: HardwareWalletAccount) => void
  onSigningComplete?: (signature: Hex) => void
  onError?: (error: Error) => void
}

export interface UseHardwareWalletReturn {
  // State
  state: HardwareWalletState
  
  // Device management
  scanDevices: () => Promise<HardwareWalletDevice[]>
  connectDevice: (deviceId: string) => Promise<void>
  disconnectDevice: (deviceId: string) => Promise<void>
  selectDevice: (device: HardwareWalletDevice) => void
  
  // Account management
  loadAccounts: (deviceId?: string, startIndex?: number, count?: number) => Promise<void>
  selectAccount: (account: HardwareWalletAccount) => void
  refreshAccounts: () => Promise<void>
  
  // Signing operations
  signTransaction: (transaction: any) => Promise<Hex>
  signMessage: (message: string) => Promise<Hex>
  signTypedData: (typedData: any) => Promise<Hex>
  
  // Utilities
  clearError: () => void
  reset: () => void
  isDeviceSupported: (deviceType: string) => boolean
  getDeviceInfo: (deviceId: string) => HardwareWalletDevice | null
}

export const useHardwareWallet = (
  options: HardwareWalletOptions = {}
): UseHardwareWalletReturn => {
  const [state, setState] = useState<HardwareWalletState>({
    isScanning: false,
    isConnecting: false,
    isSigning: false,
    devices: [],
    connectedDevices: [],
    selectedDevice: null,
    accounts: [],
    selectedAccount: null,
    signingRequests: [],
    error: null
  })

  const scanTimeoutRef = useRef<NodeJS.Timeout>()
  const {
    autoScan = false,
    autoConnect = false,
    accountCount = 5,
    onDeviceConnect,
    onDeviceDisconnect,
    onAccountSelect,
    onSigningComplete,
    onError
  } = options

  // Auto-scan on mount
  useEffect(() => {
    if (autoScan) {
      scanDevices()
    }
  }, [autoScan])

  // Auto-connect to first available device
  useEffect(() => {
    if (autoConnect && state.devices.length > 0 && state.connectedDevices.length === 0) {
      const firstDevice = state.devices[0]
      connectDevice(firstDevice.id)
    }
  }, [autoConnect, state.devices, state.connectedDevices])

  // Update connected devices
  useEffect(() => {
    const updateConnectedDevices = () => {
      const connected = hardwareWalletManager.getConnectedDevices()
      setState(prev => ({ ...prev, connectedDevices: connected }))
    }

    const interval = setInterval(updateConnectedDevices, 1000)
    return () => clearInterval(interval)
  }, [])

  // Scan for hardware wallet devices
  const scanDevices = useCallback(async (): Promise<HardwareWalletDevice[]> => {
    setState(prev => ({ ...prev, isScanning: true, error: null }))

    try {
      const devices = await hardwareWalletManager.scanForDevices()
      
      setState(prev => ({
        ...prev,
        devices,
        isScanning: false
      }))

      if (devices.length === 0) {
        toast.info('No hardware wallets found', {
          description: 'Make sure your device is connected and unlocked'
        })
      } else {
        toast.success(`Found ${devices.length} hardware wallet${devices.length > 1 ? 's' : ''}`)
      }

      return devices
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isScanning: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Device scan failed', {
        description: errorMessage
      })
      return []
    }
  }, [onError])

  // Connect to a hardware wallet device
  const connectDevice = useCallback(async (deviceId: string): Promise<void> => {
    setState(prev => ({ ...prev, isConnecting: true, error: null }))

    try {
      const connection = await hardwareWalletManager.connectDevice(deviceId)
      
      setState(prev => ({
        ...prev,
        isConnecting: false,
        selectedDevice: connection.device
      }))

      onDeviceConnect?.(connection.device)
      toast.success(`Connected to ${connection.device.model}`)

      // Auto-load accounts
      await loadAccounts(deviceId)

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isConnecting: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Device connection failed', {
        description: errorMessage
      })
    }
  }, [onDeviceConnect, onError])

  // Disconnect from a hardware wallet device
  const disconnectDevice = useCallback(async (deviceId: string): Promise<void> => {
    try {
      await hardwareWalletManager.disconnectDevice(deviceId)
      
      setState(prev => ({
        ...prev,
        selectedDevice: prev.selectedDevice?.id === deviceId ? null : prev.selectedDevice,
        accounts: prev.selectedDevice?.id === deviceId ? [] : prev.accounts,
        selectedAccount: prev.selectedDevice?.id === deviceId ? null : prev.selectedAccount
      }))

      const device = state.devices.find(d => d.id === deviceId)
      if (device) {
        onDeviceDisconnect?.(device)
        toast.success(`Disconnected from ${device.model}`)
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      onError?.(error as Error)
      toast.error('Disconnection failed', {
        description: errorMessage
      })
    }
  }, [state.devices, onDeviceDisconnect, onError])

  // Select a device
  const selectDevice = useCallback((device: HardwareWalletDevice) => {
    setState(prev => ({ ...prev, selectedDevice: device }))
  }, [])

  // Load accounts from hardware wallet
  const loadAccounts = useCallback(async (
    deviceId?: string,
    startIndex = 0,
    count = accountCount
  ): Promise<void> => {
    const targetDeviceId = deviceId || state.selectedDevice?.id
    if (!targetDeviceId) {
      throw new Error('No device selected')
    }

    try {
      const accounts = await hardwareWalletManager.getAccounts(
        targetDeviceId,
        startIndex,
        count
      )

      setState(prev => ({ ...prev, accounts }))

      if (accounts.length > 0 && !state.selectedAccount) {
        selectAccount(accounts[0])
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))
      onError?.(error as Error)
      toast.error('Failed to load accounts', {
        description: errorMessage
      })
    }
  }, [state.selectedDevice, state.selectedAccount, accountCount, onError])

  // Select an account
  const selectAccount = useCallback((account: HardwareWalletAccount) => {
    setState(prev => ({ ...prev, selectedAccount: account }))
    onAccountSelect?.(account)
  }, [onAccountSelect])

  // Refresh accounts
  const refreshAccounts = useCallback(async (): Promise<void> => {
    if (state.selectedDevice) {
      await loadAccounts(state.selectedDevice.id)
    }
  }, [state.selectedDevice, loadAccounts])

  // Sign transaction
  const signTransaction = useCallback(async (transaction: any): Promise<Hex> => {
    if (!state.selectedDevice || !state.selectedAccount) {
      throw new Error('No device or account selected')
    }

    setState(prev => ({ ...prev, isSigning: true, error: null }))

    try {
      const signature = await hardwareWalletManager.signTransaction(
        state.selectedDevice.id,
        state.selectedAccount,
        transaction
      )

      setState(prev => ({ ...prev, isSigning: false }))
      onSigningComplete?.(signature)
      toast.success('Transaction signed successfully')

      return signature
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSigning: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Transaction signing failed', {
        description: errorMessage
      })
      throw error
    }
  }, [state.selectedDevice, state.selectedAccount, onSigningComplete, onError])

  // Sign message
  const signMessage = useCallback(async (message: string): Promise<Hex> => {
    if (!state.selectedDevice || !state.selectedAccount) {
      throw new Error('No device or account selected')
    }

    setState(prev => ({ ...prev, isSigning: true, error: null }))

    try {
      const signature = await hardwareWalletManager.signMessage(
        state.selectedDevice.id,
        state.selectedAccount,
        message
      )

      setState(prev => ({ ...prev, isSigning: false }))
      onSigningComplete?.(signature)
      toast.success('Message signed successfully')

      return signature
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSigning: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Message signing failed', {
        description: errorMessage
      })
      throw error
    }
  }, [state.selectedDevice, state.selectedAccount, onSigningComplete, onError])

  // Sign typed data
  const signTypedData = useCallback(async (typedData: any): Promise<Hex> => {
    if (!state.selectedDevice || !state.selectedAccount) {
      throw new Error('No device or account selected')
    }

    setState(prev => ({ ...prev, isSigning: true, error: null }))

    try {
      // For now, treat typed data like a message
      // In a real implementation, you'd use the appropriate typed data signing method
      const signature = await hardwareWalletManager.signMessage(
        state.selectedDevice.id,
        state.selectedAccount,
        JSON.stringify(typedData)
      )

      setState(prev => ({ ...prev, isSigning: false }))
      onSigningComplete?.(signature)
      toast.success('Typed data signed successfully')

      return signature
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isSigning: false,
        error: errorMessage
      }))

      onError?.(error as Error)
      toast.error('Typed data signing failed', {
        description: errorMessage
      })
      throw error
    }
  }, [state.selectedDevice, state.selectedAccount, onSigningComplete, onError])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Reset state
  const reset = useCallback(() => {
    setState({
      isScanning: false,
      isConnecting: false,
      isSigning: false,
      devices: [],
      connectedDevices: [],
      selectedDevice: null,
      accounts: [],
      selectedAccount: null,
      signingRequests: [],
      error: null
    })
  }, [])

  // Check if device type is supported
  const isDeviceSupported = useCallback((deviceType: string): boolean => {
    const supportedTypes = ['ledger', 'trezor', 'gridplus']
    return supportedTypes.includes(deviceType.toLowerCase())
  }, [])

  // Get device info
  const getDeviceInfo = useCallback((deviceId: string): HardwareWalletDevice | null => {
    return state.devices.find(device => device.id === deviceId) || null
  }, [state.devices])

  // Update signing requests
  useEffect(() => {
    const updateSigningRequests = () => {
      const requests = hardwareWalletManager.getSigningRequests()
      setState(prev => ({ ...prev, signingRequests: requests }))
    }

    const interval = setInterval(updateSigningRequests, 1000)
    return () => clearInterval(interval)
  }, [])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (scanTimeoutRef.current) {
        clearTimeout(scanTimeoutRef.current)
      }
    }
  }, [])

  return {
    state,
    scanDevices,
    connectDevice,
    disconnectDevice,
    selectDevice,
    loadAccounts,
    selectAccount,
    refreshAccounts,
    signTransaction,
    signMessage,
    signTypedData,
    clearError,
    reset,
    isDeviceSupported,
    getDeviceInfo
  }
}
