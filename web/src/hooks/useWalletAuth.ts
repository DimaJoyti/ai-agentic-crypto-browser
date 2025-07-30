import { useState, useEffect, useCallback, useRef } from 'react'
import { useAccount, useSignMessage } from 'wagmi'
import { type Address, type Hex } from 'viem'
import { toast } from 'sonner'
import { 
  walletSecurity, 
  type AuthenticationChallenge, 
  type SecuritySession,
  type SecurityEvent 
} from '@/lib/wallet-security'
import { useWalletStore } from '@/stores/walletStore'

export interface AuthState {
  isAuthenticated: boolean
  isAuthenticating: boolean
  session: SecuritySession | null
  challenge: AuthenticationChallenge | null
  error: string | null
  lastAuthTime: number | null
}

export interface AuthOptions {
  autoAuth?: boolean
  requireAuth?: boolean
  sessionTimeout?: number
  onAuthSuccess?: (session: SecuritySession) => void
  onAuthFailure?: (error: string) => void
  onSessionExpired?: () => void
}

export interface UseWalletAuthReturn {
  // Auth state
  authState: AuthState
  
  // Auth methods
  authenticate: () => Promise<boolean>
  logout: () => void
  refreshSession: () => boolean
  
  // Challenge management
  generateChallenge: () => AuthenticationChallenge | null
  verifySignature: (signature: Hex) => Promise<boolean>
  
  // Session management
  getActiveSessions: () => SecuritySession[]
  terminateAllSessions: () => void
  
  // Security
  getSecurityEvents: () => SecurityEvent[]
  isAddressLocked: () => boolean
  
  // Utilities
  clearError: () => void
  requiresAuth: boolean
}

export const useWalletAuth = (options: AuthOptions = {}): UseWalletAuthReturn => {
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    isAuthenticating: false,
    session: null,
    challenge: null,
    error: null,
    lastAuthTime: null
  })

  const sessionCheckInterval = useRef<NodeJS.Timeout>()
  const { 
    autoAuth = false, 
    requireAuth = true,
    onAuthSuccess,
    onAuthFailure,
    onSessionExpired
  } = options

  // Wagmi hooks
  const { address, isConnected } = useAccount()
  const { signMessage, isPending: isSigningPending, error: signError } = useSignMessage()
  
  // Store hooks
  const walletStore = useWalletStore()

  // Check if authentication is required
  const requiresAuth = requireAuth && isConnected && !!address

  // Auto-authenticate on wallet connection
  useEffect(() => {
    if (autoAuth && isConnected && address && !authState.isAuthenticated) {
      authenticate()
    }
  }, [autoAuth, isConnected, address, authState.isAuthenticated])

  // Handle wallet disconnection
  useEffect(() => {
    if (!isConnected || !address) {
      logout()
    }
  }, [isConnected, address])

  // Start session monitoring
  useEffect(() => {
    if (authState.session) {
      startSessionMonitoring()
    } else {
      stopSessionMonitoring()
    }

    return () => stopSessionMonitoring()
  }, [authState.session])

  // Handle signing errors
  useEffect(() => {
    if (signError) {
      const errorMessage = getSignErrorMessage(signError)
      setAuthState(prev => ({
        ...prev,
        error: errorMessage,
        isAuthenticating: false
      }))
      onAuthFailure?.(errorMessage)
    }
  }, [signError, onAuthFailure])

  // Generate authentication challenge
  const generateChallenge = useCallback((): AuthenticationChallenge | null => {
    if (!address) {
      setAuthState(prev => ({ ...prev, error: 'No wallet connected' }))
      return null
    }

    try {
      const challenge = walletSecurity.generateChallenge(address)
      setAuthState(prev => ({ ...prev, challenge, error: null }))
      return challenge
    } catch (error) {
      const errorMessage = (error as Error).message
      setAuthState(prev => ({ ...prev, error: errorMessage }))
      return null
    }
  }, [address])

  // Main authentication function
  const authenticate = useCallback(async (): Promise<boolean> => {
    if (!address) {
      const error = 'No wallet connected'
      setAuthState(prev => ({ ...prev, error }))
      onAuthFailure?.(error)
      return false
    }

    setAuthState(prev => ({ 
      ...prev, 
      isAuthenticating: true, 
      error: null 
    }))

    try {
      // Generate challenge
      const challenge = generateChallenge()
      if (!challenge) {
        throw new Error('Failed to generate authentication challenge')
      }

      // Request signature
      const signature = (await signMessage({
        message: challenge.message
      }) as unknown) as Hex

      // Verify signature
      const result = await walletSecurity.verifySignature(
        challenge.id, 
        signature, 
        address
      )

      if (result.success && result.session) {
        setAuthState(prev => ({
          ...prev,
          isAuthenticated: true,
          isAuthenticating: false,
          session: result.session!,
          challenge: null,
          error: null,
          lastAuthTime: Date.now()
        }))

        // Update wallet store
        walletStore.updateSessionActivity()

        onAuthSuccess?.(result.session)
        toast.success('Wallet authenticated successfully')
        return true
      } else {
        throw new Error(result.error || 'Authentication failed')
      }

    } catch (error) {
      const errorMessage = (error as Error).message
      setAuthState(prev => ({
        ...prev,
        isAuthenticating: false,
        error: errorMessage,
        challenge: null
      }))

      onAuthFailure?.(errorMessage)
      toast.error('Authentication failed', {
        description: errorMessage
      })
      return false
    }
  }, [address, signMessage, generateChallenge, onAuthSuccess, onAuthFailure, walletStore])

  // Verify signature for existing challenge
  const verifySignature = useCallback(async (signature: Hex): Promise<boolean> => {
    if (!authState.challenge || !address) {
      return false
    }

    try {
      const result = await walletSecurity.verifySignature(
        authState.challenge.id,
        signature,
        address
      )

      if (result.success && result.session) {
        setAuthState(prev => ({
          ...prev,
          isAuthenticated: true,
          isAuthenticating: false,
          session: result.session!,
          challenge: null,
          error: null,
          lastAuthTime: Date.now()
        }))

        onAuthSuccess?.(result.session)
        return true
      } else {
        setAuthState(prev => ({
          ...prev,
          error: result.error || 'Signature verification failed'
        }))
        return false
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setAuthState(prev => ({ ...prev, error: errorMessage }))
      return false
    }
  }, [authState.challenge, address, onAuthSuccess])

  // Logout function
  const logout = useCallback(() => {
    if (authState.session) {
      walletSecurity.terminateSession(authState.session.id)
    }

    setAuthState({
      isAuthenticated: false,
      isAuthenticating: false,
      session: null,
      challenge: null,
      error: null,
      lastAuthTime: null
    })

    stopSessionMonitoring()
    toast.success('Logged out successfully')
  }, [authState.session])

  // Refresh session
  const refreshSession = useCallback((): boolean => {
    if (!authState.session) return false

    const refreshedSession = walletSecurity.validateSession(authState.session.id)
    
    if (refreshedSession) {
      setAuthState(prev => ({
        ...prev,
        session: refreshedSession
      }))
      return true
    } else {
      // Session expired
      setAuthState(prev => ({
        ...prev,
        isAuthenticated: false,
        session: null
      }))
      onSessionExpired?.()
      toast.warning('Session expired. Please authenticate again.')
      return false
    }
  }, [authState.session, onSessionExpired])

  // Get active sessions
  const getActiveSessions = useCallback((): SecuritySession[] => {
    if (!address) return []
    return walletSecurity.getActiveSessions(address)
  }, [address])

  // Terminate all sessions
  const terminateAllSessions = useCallback(() => {
    if (!address) return
    
    walletSecurity.terminateAllSessions(address)
    setAuthState(prev => ({
      ...prev,
      isAuthenticated: false,
      session: null
    }))
    
    toast.success('All sessions terminated')
  }, [address])

  // Get security events
  const getSecurityEvents = useCallback((): SecurityEvent[] => {
    if (!address) return []
    return walletSecurity.getSecurityEvents(address)
  }, [address])

  // Check if address is locked
  const isAddressLocked = useCallback((): boolean => {
    if (!address) return false
    // This would need to be implemented in the security manager
    return false
  }, [address])

  // Clear error
  const clearError = useCallback(() => {
    setAuthState(prev => ({ ...prev, error: null }))
  }, [])

  // Start session monitoring
  const startSessionMonitoring = useCallback(() => {
    if (sessionCheckInterval.current) return

    sessionCheckInterval.current = setInterval(() => {
      if (!refreshSession()) {
        stopSessionMonitoring()
      }
    }, 60000) // Check every minute
  }, [refreshSession])

  // Stop session monitoring
  const stopSessionMonitoring = useCallback(() => {
    if (sessionCheckInterval.current) {
      clearInterval(sessionCheckInterval.current)
      sessionCheckInterval.current = undefined
    }
  }, [])

  return {
    authState: {
      ...authState,
      isAuthenticating: authState.isAuthenticating || isSigningPending
    },
    authenticate,
    logout,
    refreshSession,
    generateChallenge,
    verifySignature,
    getActiveSessions,
    terminateAllSessions,
    getSecurityEvents,
    isAddressLocked,
    clearError,
    requiresAuth
  }
}

// Helper function to get user-friendly error messages
function getSignErrorMessage(error: Error): string {
  const message = error.message.toLowerCase()
  
  if (message.includes('user rejected') || message.includes('user denied')) {
    return 'Signature request was rejected by user'
  }
  
  if (message.includes('not connected')) {
    return 'Wallet not connected. Please connect your wallet first.'
  }
  
  if (message.includes('network')) {
    return 'Network error occurred. Please check your connection.'
  }
  
  if (message.includes('timeout')) {
    return 'Signature request timed out. Please try again.'
  }
  
  return error.message || 'Unknown error occurred during signing'
}
