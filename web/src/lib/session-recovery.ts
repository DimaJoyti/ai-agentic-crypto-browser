import { type Address } from 'viem'
import { useWalletStore } from '@/stores/walletStore'
import { 
  type WalletProvider, 
  getWalletPreference, 
  WALLET_PROVIDERS 
} from '@/lib/wallet-providers'

export interface SessionRecoveryData {
  address: Address
  chainId: number
  walletId: string
  connectedAt: number
  lastActivity: number
  sessionId: string
}

export interface RecoveryOptions {
  maxAge?: number // Maximum age in milliseconds
  requireUserConfirmation?: boolean
  autoReconnect?: boolean
  fallbackToPreferred?: boolean
}

export class SessionRecoveryManager {
  private static instance: SessionRecoveryManager
  private recoveryKey = 'wallet_session_recovery'
  private maxSessionAge = 24 * 60 * 60 * 1000 // 24 hours
  private checkInterval: NodeJS.Timeout | null = null

  private constructor() {
    this.startPeriodicCleanup()
  }

  static getInstance(): SessionRecoveryManager {
    if (!SessionRecoveryManager.instance) {
      SessionRecoveryManager.instance = new SessionRecoveryManager()
    }
    return SessionRecoveryManager.instance
  }

  /**
   * Save current session for recovery
   */
  saveSession(data: SessionRecoveryData): void {
    try {
      const sessions = this.getSavedSessions()
      const updatedSessions = {
        ...sessions,
        [data.address]: {
          ...data,
          lastActivity: Date.now()
        }
      }
      
      localStorage.setItem(this.recoveryKey, JSON.stringify(updatedSessions))
      
      // Also update the wallet store
      const walletStore = useWalletStore.getState()
      walletStore.updateSessionActivity()
      
    } catch (error) {
      console.warn('Failed to save session for recovery:', error)
    }
  }

  /**
   * Get all saved sessions
   */
  getSavedSessions(): Record<string, SessionRecoveryData> {
    try {
      const saved = localStorage.getItem(this.recoveryKey)
      return saved ? JSON.parse(saved) : {}
    } catch (error) {
      console.warn('Failed to get saved sessions:', error)
      return {}
    }
  }

  /**
   * Get the most recent valid session
   */
  getMostRecentSession(options: RecoveryOptions = {}): SessionRecoveryData | null {
    const { maxAge = this.maxSessionAge } = options
    const sessions = this.getSavedSessions()
    const now = Date.now()

    const validSessions = Object.values(sessions).filter(session => {
      const age = now - session.lastActivity
      return age <= maxAge
    })

    if (validSessions.length === 0) return null

    // Return the most recently active session
    return validSessions.sort((a, b) => b.lastActivity - a.lastActivity)[0]
  }

  /**
   * Get session for specific address
   */
  getSessionForAddress(address: Address): SessionRecoveryData | null {
    const sessions = this.getSavedSessions()
    const session = sessions[address]
    
    if (!session) return null
    
    const age = Date.now() - session.lastActivity
    if (age > this.maxSessionAge) {
      this.removeSession(address)
      return null
    }
    
    return session
  }

  /**
   * Check if there's a recoverable session
   */
  hasRecoverableSession(options: RecoveryOptions = {}): boolean {
    return this.getMostRecentSession(options) !== null
  }

  /**
   * Attempt to recover the most recent session
   */
  async recoverSession(
    connectWallet: (walletId: string) => Promise<void>,
    options: RecoveryOptions = {}
  ): Promise<boolean> {
    const session = this.getMostRecentSession(options)
    if (!session) return false

    const { requireUserConfirmation = true, autoReconnect = false } = options

    try {
      // Check if user confirmation is required
      if (requireUserConfirmation && !autoReconnect) {
        const confirmed = await this.requestUserConfirmation(session)
        if (!confirmed) return false
      }

      // Attempt to reconnect
      await connectWallet(session.walletId)
      
      // Update session activity
      this.saveSession({
        ...session,
        lastActivity: Date.now()
      })

      return true
    } catch (error) {
      console.error('Session recovery failed:', error)
      this.removeSession(session.address)
      return false
    }
  }

  /**
   * Recover session with fallback to preferred wallet
   */
  async recoverWithFallback(
    connectWallet: (walletId: string) => Promise<void>,
    options: RecoveryOptions = {}
  ): Promise<boolean> {
    // First try to recover the most recent session
    const recovered = await this.recoverSession(connectWallet, options)
    if (recovered) return true

    // If that fails and fallback is enabled, try preferred wallet
    if (options.fallbackToPreferred) {
      const preferredWallet = getWalletPreference()
      if (preferredWallet) {
        try {
          await connectWallet(preferredWallet)
          return true
        } catch (error) {
          console.warn('Fallback to preferred wallet failed:', error)
        }
      }
    }

    return false
  }

  /**
   * Request user confirmation for session recovery
   */
  private async requestUserConfirmation(session: SessionRecoveryData): Promise<boolean> {
    const wallet = WALLET_PROVIDERS.find(w => w.id === session.walletId)
    const walletName = wallet?.name || session.walletId
    const timeAgo = this.formatTimeAgo(Date.now() - session.lastActivity)
    
    return new Promise((resolve) => {
      // Create a simple confirmation dialog
      const confirmed = window.confirm(
        `Would you like to reconnect to your ${walletName} wallet?\n\n` +
        `Address: ${session.address.slice(0, 6)}...${session.address.slice(-4)}\n` +
        `Last active: ${timeAgo} ago`
      )
      resolve(confirmed)
    })
  }

  /**
   * Remove session for specific address
   */
  removeSession(address: Address): void {
    try {
      const sessions = this.getSavedSessions()
      delete sessions[address]
      localStorage.setItem(this.recoveryKey, JSON.stringify(sessions))
    } catch (error) {
      console.warn('Failed to remove session:', error)
    }
  }

  /**
   * Clear all saved sessions
   */
  clearAllSessions(): void {
    try {
      localStorage.removeItem(this.recoveryKey)
    } catch (error) {
      console.warn('Failed to clear sessions:', error)
    }
  }

  /**
   * Clean up expired sessions
   */
  cleanupExpiredSessions(): void {
    const sessions = this.getSavedSessions()
    const now = Date.now()
    const validSessions: Record<string, SessionRecoveryData> = {}

    Object.entries(sessions).forEach(([address, session]) => {
      const age = now - session.lastActivity
      if (age <= this.maxSessionAge) {
        validSessions[address] = session
      }
    })

    try {
      localStorage.setItem(this.recoveryKey, JSON.stringify(validSessions))
    } catch (error) {
      console.warn('Failed to cleanup expired sessions:', error)
    }
  }

  /**
   * Start periodic cleanup of expired sessions
   */
  private startPeriodicCleanup(): void {
    // Clean up every hour
    this.checkInterval = setInterval(() => {
      this.cleanupExpiredSessions()
    }, 60 * 60 * 1000)
  }

  /**
   * Stop periodic cleanup
   */
  stopPeriodicCleanup(): void {
    if (this.checkInterval) {
      clearInterval(this.checkInterval)
      this.checkInterval = null
    }
  }

  /**
   * Format time ago string
   */
  private formatTimeAgo(milliseconds: number): string {
    const seconds = Math.floor(milliseconds / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)

    if (days > 0) return `${days} day${days > 1 ? 's' : ''}`
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''}`
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''}`
    return `${seconds} second${seconds > 1 ? 's' : ''}`
  }

  /**
   * Get session statistics
   */
  getSessionStats(): {
    totalSessions: number
    activeSessions: number
    oldestSession: number | null
    newestSession: number | null
  } {
    const sessions = Object.values(this.getSavedSessions())
    const now = Date.now()
    
    const activeSessions = sessions.filter(s => 
      (now - s.lastActivity) <= this.maxSessionAge
    )

    return {
      totalSessions: sessions.length,
      activeSessions: activeSessions.length,
      oldestSession: sessions.length > 0 
        ? Math.min(...sessions.map(s => s.lastActivity))
        : null,
      newestSession: sessions.length > 0 
        ? Math.max(...sessions.map(s => s.lastActivity))
        : null
    }
  }

  /**
   * Export session data for backup
   */
  exportSessions(): string {
    const sessions = this.getSavedSessions()
    const walletStore = useWalletStore.getState()
    
    const exportData = {
      sessions,
      preferences: walletStore.preferences,
      analytics: walletStore.analytics,
      exportedAt: Date.now(),
      version: '1.0'
    }
    
    return JSON.stringify(exportData, null, 2)
  }

  /**
   * Import session data from backup
   */
  importSessions(data: string): boolean {
    try {
      const importData = JSON.parse(data)
      
      if (importData.sessions) {
        localStorage.setItem(this.recoveryKey, JSON.stringify(importData.sessions))
      }
      
      if (importData.preferences || importData.analytics) {
        const walletStore = useWalletStore.getState()
        
        if (importData.preferences) {
          walletStore.updatePreferences(importData.preferences)
        }
        
        if (importData.analytics) {
          // Merge analytics data
          walletStore.resetAnalytics()
          // Note: You might want to implement a more sophisticated merge strategy
        }
      }
      
      return true
    } catch (error) {
      console.error('Failed to import session data:', error)
      return false
    }
  }
}

// Export singleton instance
export const sessionRecovery = SessionRecoveryManager.getInstance()

// Hook for using session recovery in React components
export const useSessionRecovery = () => {
  const walletStore = useWalletStore()
  
  const saveCurrentSession = () => {
    if (walletStore.currentWallet && walletStore.currentSession) {
      sessionRecovery.saveSession({
        address: walletStore.currentWallet.address,
        chainId: walletStore.currentWallet.chainId,
        walletId: walletStore.currentWallet.id,
        connectedAt: walletStore.currentWallet.connectedAt,
        lastActivity: Date.now(),
        sessionId: walletStore.currentSession.id
      })
    }
  }

  const hasRecoverableSession = (options?: RecoveryOptions) => {
    return sessionRecovery.hasRecoverableSession(options)
  }

  const recoverSession = async (
    connectWallet: (walletId: string) => Promise<void>,
    options?: RecoveryOptions
  ) => {
    return sessionRecovery.recoverSession(connectWallet, options)
  }

  const clearSessions = () => {
    sessionRecovery.clearAllSessions()
  }

  return {
    saveCurrentSession,
    hasRecoverableSession,
    recoverSession,
    clearSessions,
    sessionStats: sessionRecovery.getSessionStats(),
    exportSessions: sessionRecovery.exportSessions,
    importSessions: sessionRecovery.importSessions
  }
}
