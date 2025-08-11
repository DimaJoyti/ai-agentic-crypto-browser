import { type Address } from 'viem'
import { useWalletStore, type ConnectedWallet, type WalletPreferences } from '@/stores/walletStore'

export interface PersistenceConfig {
  enableEncryption?: boolean
  compressionLevel?: number
  maxStorageSize?: number // in bytes
  autoCleanup?: boolean
  backupInterval?: number // in milliseconds
}

export interface StorageQuota {
  used: number
  available: number
  total: number
  percentage: number
}

export interface BackupData {
  wallets: ConnectedWallet[]
  preferences: WalletPreferences
  sessions: any[]
  analytics: any
  timestamp: number
  version: string
  checksum: string
}

export class WalletPersistenceManager {
  private static instance: WalletPersistenceManager
  private config: PersistenceConfig
  private storageKeys = {
    wallets: 'wallet_connected_wallets',
    preferences: 'wallet_preferences',
    sessions: 'wallet_sessions',
    analytics: 'wallet_analytics',
    backup: 'wallet_backup',
    metadata: 'wallet_metadata'
  }

  private constructor(config: PersistenceConfig = {}) {
    this.config = {
      enableEncryption: false,
      compressionLevel: 1,
      maxStorageSize: 5 * 1024 * 1024, // 5MB
      autoCleanup: true,
      backupInterval: 24 * 60 * 60 * 1000, // 24 hours
      ...config
    }

    if (this.config.autoCleanup) {
      this.startAutoCleanup()
    }
  }

  static getInstance(config?: PersistenceConfig): WalletPersistenceManager {
    if (!WalletPersistenceManager.instance) {
      WalletPersistenceManager.instance = new WalletPersistenceManager(config)
    }
    return WalletPersistenceManager.instance
  }

  /**
   * Save wallet data with optional encryption and compression
   */
  async saveWalletData(key: string, data: any): Promise<boolean> {
    try {
      let processedData = JSON.stringify(data)

      // Apply compression if enabled
      if (this.config.compressionLevel && this.config.compressionLevel > 0) {
        processedData = await this.compressData(processedData)
      }

      // Apply encryption if enabled
      if (this.config.enableEncryption) {
        processedData = await this.encryptData(processedData)
      }

      // Check storage quota
      const quota = await this.getStorageQuota()
      const dataSize = new Blob([processedData]).size

      if (quota.used + dataSize > this.config.maxStorageSize!) {
        console.warn('Storage quota exceeded, attempting cleanup...')
        await this.performCleanup()
        
        // Check again after cleanup
        const newQuota = await this.getStorageQuota()
        if (newQuota.used + dataSize > this.config.maxStorageSize!) {
          throw new Error('Insufficient storage space')
        }
      }

      localStorage.setItem(key, processedData)
      
      // Update metadata
      await this.updateMetadata(key, dataSize)
      
      return true
    } catch (error) {
      console.error('Failed to save wallet data:', error)
      return false
    }
  }

  /**
   * Load wallet data with decryption and decompression
   */
  async loadWalletData<T>(key: string): Promise<T | null> {
    try {
      let data = localStorage.getItem(key)
      if (!data) return null

      // Apply decryption if enabled
      if (this.config.enableEncryption) {
        data = await this.decryptData(data)
      }

      // Apply decompression if enabled
      if (this.config.compressionLevel && this.config.compressionLevel > 0) {
        data = await this.decompressData(data)
      }

      return JSON.parse(data) as T
    } catch (error) {
      console.error('Failed to load wallet data:', error)
      return null
    }
  }

  /**
   * Get current storage quota information
   */
  async getStorageQuota(): Promise<StorageQuota> {
    try {
      if ('storage' in navigator && 'estimate' in navigator.storage) {
        const estimate = await navigator.storage.estimate()
        const used = estimate.usage || 0
        const total = estimate.quota || 0
        const available = total - used
        const percentage = total > 0 ? (used / total) * 100 : 0

        return { used, available, total, percentage }
      }
    } catch (error) {
      console.warn('Failed to get storage estimate:', error)
    }

    // Fallback: estimate based on localStorage
    const localStorageSize = this.getLocalStorageSize()
    const estimatedTotal = 10 * 1024 * 1024 // 10MB estimate
    
    return {
      used: localStorageSize,
      available: estimatedTotal - localStorageSize,
      total: estimatedTotal,
      percentage: (localStorageSize / estimatedTotal) * 100
    }
  }

  /**
   * Get localStorage size in bytes
   */
  private getLocalStorageSize(): number {
    let total = 0
    for (const key in localStorage) {
      if (localStorage.hasOwnProperty(key)) {
        total += localStorage[key].length + key.length
      }
    }
    return total
  }

  /**
   * Perform cleanup of old and unnecessary data
   */
  async performCleanup(): Promise<void> {
    try {
      // Remove expired sessions
      const walletStore = useWalletStore.getState()
      walletStore.clearExpiredSessions()

      // Clean up old analytics data (keep last 6 months)
      const sixMonthsAgo = new Date()
      sixMonthsAgo.setMonth(sixMonthsAgo.getMonth() - 6)
      const cutoffMonth = sixMonthsAgo.toISOString().slice(0, 7)

      const analytics = walletStore.analytics
      const cleanedMonthlyStats: Record<string, any> = {}
      
      Object.entries(analytics.monthlyStats).forEach(([month, stats]) => {
        if (month >= cutoffMonth) {
          cleanedMonthlyStats[month] = stats
        }
      })

      // Update analytics with cleaned data
      walletStore.resetAnalytics()
      // Note: You'd need to implement a method to set specific analytics data

      // Remove orphaned data
      this.removeOrphanedData()

      console.log('Wallet data cleanup completed')
    } catch (error) {
      console.error('Cleanup failed:', error)
    }
  }

  /**
   * Remove orphaned data that's no longer referenced
   */
  private removeOrphanedData(): void {
    const validKeys = Object.values(this.storageKeys)
    const keysToRemove: string[] = []

    for (const key in localStorage) {
      if (key.startsWith('wallet_') && !validKeys.includes(key)) {
        keysToRemove.push(key)
      }
    }

    keysToRemove.forEach(key => {
      localStorage.removeItem(key)
    })

    if (keysToRemove.length > 0) {
      console.log(`Removed ${keysToRemove.length} orphaned storage keys`)
    }
  }

  /**
   * Create a backup of all wallet data
   */
  async createBackup(): Promise<BackupData> {
    const walletStore = useWalletStore.getState()
    
    const backupData: BackupData = {
      wallets: walletStore.connectedWallets,
      preferences: walletStore.preferences,
      sessions: walletStore.sessions,
      analytics: walletStore.analytics,
      timestamp: Date.now(),
      version: '1.0',
      checksum: ''
    }

    // Generate checksum
    backupData.checksum = await this.generateChecksum(backupData)

    // Save backup
    await this.saveWalletData(this.storageKeys.backup, backupData)

    return backupData
  }

  /**
   * Restore from backup
   */
  async restoreFromBackup(backupData: BackupData): Promise<boolean> {
    try {
      // Verify checksum
      const expectedChecksum = await this.generateChecksum({
        ...backupData,
        checksum: ''
      })

      if (backupData.checksum !== expectedChecksum) {
        throw new Error('Backup data integrity check failed')
      }

      // Restore data to store
      const walletStore = useWalletStore.getState()
      
      // Clear existing data
      walletStore.resetStore()

      // Restore wallets
      backupData.wallets.forEach(wallet => {
        walletStore.addConnectedWallet(wallet)
      })

      // Restore preferences
      walletStore.updatePreferences(backupData.preferences)

      // Note: You might want to implement methods to restore sessions and analytics

      console.log('Backup restored successfully')
      return true
    } catch (error) {
      console.error('Failed to restore backup:', error)
      return false
    }
  }

  /**
   * Export wallet data for external backup
   */
  async exportWalletData(): Promise<string> {
    const backup = await this.createBackup()
    return JSON.stringify(backup, null, 2)
  }

  /**
   * Import wallet data from external backup
   */
  async importWalletData(data: string): Promise<boolean> {
    try {
      const backupData = JSON.parse(data) as BackupData
      return await this.restoreFromBackup(backupData)
    } catch (error) {
      console.error('Failed to import wallet data:', error)
      return false
    }
  }

  /**
   * Start automatic cleanup process
   */
  private startAutoCleanup(): void {
    setInterval(() => {
      this.performCleanup()
    }, this.config.backupInterval!)
  }

  /**
   * Update metadata for storage tracking
   */
  private async updateMetadata(key: string, size: number): Promise<void> {
    try {
      const metadata = await this.loadWalletData<Record<string, any>>(this.storageKeys.metadata) || {}
      metadata[key] = {
        size,
        lastUpdated: Date.now()
      }
      await this.saveWalletData(this.storageKeys.metadata, metadata)
    } catch (error) {
      console.warn('Failed to update metadata:', error)
    }
  }

  /**
   * Generate checksum for data integrity
   */
  private async generateChecksum(data: any): Promise<string> {
    const jsonString = JSON.stringify(data)
    const encoder = new TextEncoder()
    const dataBuffer = encoder.encode(jsonString)
    
    if ('crypto' in window && 'subtle' in window.crypto) {
      const hashBuffer = await window.crypto.subtle.digest('SHA-256', dataBuffer as BufferSource)
      const hashArray = Array.from(new Uint8Array(hashBuffer))
      return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
    }
    
    // Fallback: simple hash
    let hash = 0
    for (let i = 0; i < jsonString.length; i++) {
      const char = jsonString.charCodeAt(i)
      hash = ((hash << 5) - hash) + char
      hash = hash & hash // Convert to 32-bit integer
    }
    return hash.toString(16)
  }

  /**
   * Simple compression (placeholder - in real implementation, use a proper compression library)
   */
  private async compressData(data: string): Promise<string> {
    // This is a placeholder. In a real implementation, you'd use a compression library
    // like pako or lz-string
    return data
  }

  /**
   * Simple decompression (placeholder)
   */
  private async decompressData(data: string): Promise<string> {
    // This is a placeholder. In a real implementation, you'd use a compression library
    return data
  }

  /**
   * Simple encryption (placeholder - in real implementation, use proper encryption)
   */
  private async encryptData(data: string): Promise<string> {
    // This is a placeholder. In a real implementation, you'd use proper encryption
    // with user-provided passwords or device-specific keys
    return btoa(data)
  }

  /**
   * Simple decryption (placeholder)
   */
  private async decryptData(data: string): Promise<string> {
    // This is a placeholder. In a real implementation, you'd use proper decryption
    try {
      return atob(data)
    } catch {
      return data // Fallback for unencrypted data
    }
  }

  /**
   * Clear all wallet data
   */
  clearAllData(): void {
    Object.values(this.storageKeys).forEach(key => {
      localStorage.removeItem(key)
    })
    
    // Also clear the main store
    const walletStore = useWalletStore.getState()
    walletStore.resetStore()
  }

  /**
   * Get storage statistics
   */
  async getStorageStats(): Promise<{
    quota: StorageQuota
    itemCount: number
    totalSize: number
    itemSizes: Record<string, number>
  }> {
    const quota = await this.getStorageQuota()
    const metadata = await this.loadWalletData<Record<string, any>>(this.storageKeys.metadata) || {}
    
    const itemSizes: Record<string, number> = {}
    let totalSize = 0
    let itemCount = 0

    Object.entries(metadata).forEach(([key, data]) => {
      if (data && typeof data === 'object' && 'size' in data) {
        itemSizes[key] = data.size
        totalSize += data.size
        itemCount++
      }
    })

    return {
      quota,
      itemCount,
      totalSize,
      itemSizes
    }
  }
}

// Export singleton instance
export const walletPersistence = WalletPersistenceManager.getInstance()

// Hook for using persistence in React components
export const useWalletPersistence = () => {
  const createBackup = () => walletPersistence.createBackup()
  const exportData = () => walletPersistence.exportWalletData()
  const importData = (data: string) => walletPersistence.importWalletData(data)
  const clearData = () => walletPersistence.clearAllData()
  const getStats = () => walletPersistence.getStorageStats()
  const getQuota = () => walletPersistence.getStorageQuota()

  return {
    createBackup,
    exportData,
    importData,
    clearData,
    getStats,
    getQuota
  }
}
