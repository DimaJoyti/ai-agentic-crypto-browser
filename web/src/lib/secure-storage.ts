import { type Address } from 'viem'

export interface SecureStorageConfig {
  encryptionEnabled: boolean
  compressionEnabled: boolean
  keyRotationInterval: number // in milliseconds
  maxStorageSize: number // in bytes
  autoCleanup: boolean
}

export interface StorageItem {
  id: string
  key: string
  value: string
  encrypted: boolean
  compressed: boolean
  createdAt: number
  lastAccessed: number
  expiresAt?: number
  metadata?: Record<string, any>
}

export interface StorageStats {
  totalItems: number
  totalSize: number
  encryptedItems: number
  compressedItems: number
  expiredItems: number
  oldestItem: number | null
  newestItem: number | null
}

export class SecureStorage {
  private static instance: SecureStorage
  private config: SecureStorageConfig
  private encryptionKey: CryptoKey | null = null
  private keyVersion: number = 1
  private storagePrefix = 'secure_wallet_'
  private metadataKey = 'secure_storage_metadata'

  private constructor(config: Partial<SecureStorageConfig> = {}) {
    this.config = {
      encryptionEnabled: true,
      compressionEnabled: false, // Disabled for simplicity
      keyRotationInterval: 30 * 24 * 60 * 60 * 1000, // 30 days
      maxStorageSize: 10 * 1024 * 1024, // 10MB
      autoCleanup: true,
      ...config
    }

    this.initializeEncryption()
    
    if (this.config.autoCleanup) {
      this.startAutoCleanup()
    }
  }

  static getInstance(config?: Partial<SecureStorageConfig>): SecureStorage {
    if (!SecureStorage.instance) {
      SecureStorage.instance = new SecureStorage(config)
    }
    return SecureStorage.instance
  }

  /**
   * Initialize encryption system
   */
  private async initializeEncryption(): Promise<void> {
    if (!this.config.encryptionEnabled || !window.crypto?.subtle) {
      console.warn('Encryption not available or disabled')
      return
    }

    try {
      const keyData = localStorage.getItem(`${this.storagePrefix}encryption_key`)
      const keyVersionData = localStorage.getItem(`${this.storagePrefix}key_version`)
      
      if (keyData && keyVersionData) {
        // Load existing key
        const keyBuffer = new Uint8Array(JSON.parse(keyData))
        this.keyVersion = parseInt(keyVersionData)
        
        this.encryptionKey = await window.crypto.subtle.importKey(
          'raw',
          keyBuffer,
          { name: 'AES-GCM' },
          false,
          ['encrypt', 'decrypt']
        )
      } else {
        // Generate new key
        await this.generateNewEncryptionKey()
      }

      // Check if key rotation is needed
      this.checkKeyRotation()
    } catch (error) {
      console.error('Failed to initialize encryption:', error)
      this.config.encryptionEnabled = false
    }
  }

  /**
   * Generate new encryption key
   */
  private async generateNewEncryptionKey(): Promise<void> {
    if (!window.crypto?.subtle) return

    try {
      this.encryptionKey = await window.crypto.subtle.generateKey(
        { name: 'AES-GCM', length: 256 },
        true,
        ['encrypt', 'decrypt']
      )

      const keyBuffer = await window.crypto.subtle.exportKey('raw', this.encryptionKey)
      const keyArray = Array.from(new Uint8Array(keyBuffer))
      
      localStorage.setItem(`${this.storagePrefix}encryption_key`, JSON.stringify(keyArray))
      localStorage.setItem(`${this.storagePrefix}key_version`, this.keyVersion.toString())
      localStorage.setItem(`${this.storagePrefix}key_created`, Date.now().toString())
    } catch (error) {
      console.error('Failed to generate encryption key:', error)
      this.config.encryptionEnabled = false
    }
  }

  /**
   * Check if key rotation is needed
   */
  private checkKeyRotation(): void {
    const keyCreated = localStorage.getItem(`${this.storagePrefix}key_created`)
    if (!keyCreated) return

    const keyAge = Date.now() - parseInt(keyCreated)
    if (keyAge > this.config.keyRotationInterval) {
      this.rotateEncryptionKey()
    }
  }

  /**
   * Rotate encryption key
   */
  private async rotateEncryptionKey(): Promise<void> {
    console.log('Rotating encryption key...')
    
    // Get all encrypted items
    const encryptedItems = this.getAllItems().filter(item => item.encrypted)
    
    // Generate new key
    this.keyVersion++
    await this.generateNewEncryptionKey()
    
    // Re-encrypt all items with new key
    for (const item of encryptedItems) {
      try {
        const decryptedValue = await this.decryptValue(item.value, this.keyVersion - 1)
        const encryptedValue = await this.encryptValue(decryptedValue)
        
        const updatedItem: StorageItem = {
          ...item,
          value: encryptedValue,
          lastAccessed: Date.now()
        }
        
        this.saveItem(updatedItem)
      } catch (error) {
        console.error(`Failed to re-encrypt item ${item.id}:`, error)
      }
    }
    
    console.log('Key rotation completed')
  }

  /**
   * Store data securely
   */
  async setItem(
    key: string, 
    value: any, 
    options: {
      encrypt?: boolean
      compress?: boolean
      expiresIn?: number // milliseconds
      metadata?: Record<string, any>
    } = {}
  ): Promise<boolean> {
    try {
      const {
        encrypt = this.config.encryptionEnabled,
        compress = this.config.compressionEnabled,
        expiresIn,
        metadata
      } = options

      // Check storage quota
      if (!(await this.checkStorageQuota())) {
        throw new Error('Storage quota exceeded')
      }

      let processedValue = JSON.stringify(value)

      // Compress if enabled
      if (compress) {
        processedValue = await this.compressValue(processedValue)
      }

      // Encrypt if enabled
      if (encrypt && this.encryptionKey) {
        processedValue = await this.encryptValue(processedValue)
      }

      const now = Date.now()
      const item: StorageItem = {
        id: this.generateId(),
        key,
        value: processedValue,
        encrypted: encrypt && !!this.encryptionKey,
        compressed: compress,
        createdAt: now,
        lastAccessed: now,
        expiresAt: expiresIn ? now + expiresIn : undefined,
        metadata
      }

      this.saveItem(item)
      this.updateMetadata()
      
      return true
    } catch (error) {
      console.error('Failed to store item:', error)
      return false
    }
  }

  /**
   * Retrieve data securely
   */
  async getItem<T = any>(key: string): Promise<T | null> {
    try {
      const item = this.loadItem(key)
      if (!item) return null

      // Check expiration
      if (item.expiresAt && Date.now() > item.expiresAt) {
        this.removeItem(key)
        return null
      }

      // Update last accessed
      item.lastAccessed = Date.now()
      this.saveItem(item)

      let value = item.value

      // Decrypt if encrypted
      if (item.encrypted) {
        value = await this.decryptValue(value)
      }

      // Decompress if compressed
      if (item.compressed) {
        value = await this.decompressValue(value)
      }

      return JSON.parse(value) as T
    } catch (error) {
      console.error('Failed to retrieve item:', error)
      return null
    }
  }

  /**
   * Remove item
   */
  removeItem(key: string): boolean {
    try {
      const storageKey = this.getStorageKey(key)
      localStorage.removeItem(storageKey)
      this.updateMetadata()
      return true
    } catch (error) {
      console.error('Failed to remove item:', error)
      return false
    }
  }

  /**
   * Check if item exists
   */
  hasItem(key: string): boolean {
    const storageKey = this.getStorageKey(key)
    return localStorage.getItem(storageKey) !== null
  }

  /**
   * Get all keys
   */
  getAllKeys(): string[] {
    const keys: string[] = []
    const prefix = this.storagePrefix
    
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(prefix) && !key.includes('encryption_key') && !key.includes('metadata')) {
        keys.push(key.substring(prefix.length))
      }
    }
    
    return keys
  }

  /**
   * Get all items
   */
  getAllItems(): StorageItem[] {
    const items: StorageItem[] = []
    const keys = this.getAllKeys()
    
    for (const key of keys) {
      const item = this.loadItem(key)
      if (item) {
        items.push(item)
      }
    }
    
    return items
  }

  /**
   * Clear all data
   */
  clear(): void {
    const keys = this.getAllKeys()
    keys.forEach(key => this.removeItem(key))
    
    // Also remove encryption key and metadata
    localStorage.removeItem(`${this.storagePrefix}encryption_key`)
    localStorage.removeItem(`${this.storagePrefix}key_version`)
    localStorage.removeItem(`${this.storagePrefix}key_created`)
    localStorage.removeItem(this.metadataKey)
  }

  /**
   * Get storage statistics
   */
  getStats(): StorageStats {
    const items = this.getAllItems()
    const now = Date.now()
    
    let totalSize = 0
    let encryptedItems = 0
    let compressedItems = 0
    let expiredItems = 0
    let oldestItem: number | null = null
    let newestItem: number | null = null

    for (const item of items) {
      totalSize += item.value.length
      
      if (item.encrypted) encryptedItems++
      if (item.compressed) compressedItems++
      if (item.expiresAt && now > item.expiresAt) expiredItems++
      
      if (oldestItem === null || item.createdAt < oldestItem) {
        oldestItem = item.createdAt
      }
      
      if (newestItem === null || item.createdAt > newestItem) {
        newestItem = item.createdAt
      }
    }

    return {
      totalItems: items.length,
      totalSize,
      encryptedItems,
      compressedItems,
      expiredItems,
      oldestItem,
      newestItem
    }
  }

  /**
   * Cleanup expired items
   */
  cleanup(): number {
    const items = this.getAllItems()
    const now = Date.now()
    let removedCount = 0

    for (const item of items) {
      if (item.expiresAt && now > item.expiresAt) {
        this.removeItem(item.key)
        removedCount++
      }
    }

    return removedCount
  }

  /**
   * Private helper methods
   */
  private generateId(): string {
    return `${Date.now()}_${Math.random().toString(36).substring(2, 11)}`
  }

  private getStorageKey(key: string): string {
    return `${this.storagePrefix}${key}`
  }

  private saveItem(item: StorageItem): void {
    const storageKey = this.getStorageKey(item.key)
    localStorage.setItem(storageKey, JSON.stringify(item))
  }

  private loadItem(key: string): StorageItem | null {
    try {
      const storageKey = this.getStorageKey(key)
      const data = localStorage.getItem(storageKey)
      return data ? JSON.parse(data) : null
    } catch (error) {
      console.error('Failed to load item:', error)
      return null
    }
  }

  private async encryptValue(value: string): Promise<string> {
    if (!this.encryptionKey) return value

    try {
      const encoder = new TextEncoder()
      const data = encoder.encode(value)
      const iv = window.crypto.getRandomValues(new Uint8Array(12))
      
      const encrypted = await window.crypto.subtle.encrypt(
        { name: 'AES-GCM', iv },
        this.encryptionKey,
        data as BufferSource
      )

      const encryptedArray = new Uint8Array(encrypted)
      const result = new Uint8Array(iv.length + encryptedArray.length)
      result.set(iv)
      result.set(encryptedArray, iv.length)

      return btoa(String.fromCharCode(...Array.from(result)))
    } catch (error) {
      console.error('Encryption failed:', error)
      return value
    }
  }

  private async decryptValue(encryptedValue: string, keyVersion?: number): Promise<string> {
    if (!this.encryptionKey) return encryptedValue

    try {
      const data = new Uint8Array(
        atob(encryptedValue).split('').map(char => char.charCodeAt(0))
      )
      
      const iv = data.slice(0, 12)
      const encrypted = data.slice(12)

      const decrypted = await window.crypto.subtle.decrypt(
        { name: 'AES-GCM', iv },
        this.encryptionKey,
        encrypted
      )

      const decoder = new TextDecoder()
      return decoder.decode(decrypted)
    } catch (error) {
      console.error('Decryption failed:', error)
      return encryptedValue
    }
  }

  private async compressValue(value: string): Promise<string> {
    // Placeholder for compression
    // In a real implementation, you'd use a compression library like pako
    return value
  }

  private async decompressValue(value: string): Promise<string> {
    // Placeholder for decompression
    return value
  }

  private async checkStorageQuota(): Promise<boolean> {
    try {
      const stats = this.getStats()
      return stats.totalSize < this.config.maxStorageSize
    } catch (error) {
      console.error('Failed to check storage quota:', error)
      return true // Allow operation if check fails
    }
  }

  private updateMetadata(): void {
    try {
      const metadata = {
        lastUpdated: Date.now(),
        keyVersion: this.keyVersion,
        itemCount: this.getAllKeys().length
      }
      localStorage.setItem(this.metadataKey, JSON.stringify(metadata))
    } catch (error) {
      console.error('Failed to update metadata:', error)
    }
  }

  private startAutoCleanup(): void {
    // Run cleanup every hour
    setInterval(() => {
      const removed = this.cleanup()
      if (removed > 0) {
        console.log(`Auto-cleanup removed ${removed} expired items`)
      }
    }, 60 * 60 * 1000)
  }
}

// Export singleton instance
export const secureStorage = SecureStorage.getInstance()

// Hook for using secure storage in React components
export const useSecureStorage = () => {
  const setSecureItem = async (key: string, value: any, options?: any) =>
    secureStorage.setItem(key, value, options)

  const getSecureItem = async <T = any>(key: string) =>
    secureStorage.getItem<T>(key)

  const removeSecureItem = (key: string) =>
    secureStorage.removeItem(key)

  const hasSecureItem = (key: string) =>
    secureStorage.hasItem(key)

  const clearSecureStorage = () =>
    secureStorage.clear()

  const getSecureStorageStats = () =>
    secureStorage.getStats()

  const cleanupSecureStorage = () =>
    secureStorage.cleanup()

  return {
    setSecureItem,
    getSecureItem,
    removeSecureItem,
    hasSecureItem,
    clearSecureStorage,
    getSecureStorageStats,
    cleanupSecureStorage
  }
}

// Utility functions for address-specific secure storage
export const createAddressStorage = (address: Address) => {
  const prefix = `addr_${address.toLowerCase()}_`

  return {
    setItem: (key: string, value: any, options?: any) =>
      secureStorage.setItem(`${prefix}${key}`, value, options),

    getItem: <T = any>(key: string) =>
      secureStorage.getItem<T>(`${prefix}${key}`),

    removeItem: (key: string) =>
      secureStorage.removeItem(`${prefix}${key}`),

    hasItem: (key: string) =>
      secureStorage.hasItem(`${prefix}${key}`),

    clearAll: () => {
      const keys = secureStorage.getAllKeys()
      keys.filter(key => key.startsWith(prefix))
           .forEach(key => secureStorage.removeItem(key))
    }
  }
}
