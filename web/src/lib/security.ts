import { createHash, randomBytes } from 'crypto'

// Security utilities for wallet management
export class WalletSecurity {
  private static readonly STORAGE_PREFIX = 'ai_browser_wallet_'
  private static readonly ENCRYPTION_ALGORITHM = 'AES-GCM'
  private static readonly KEY_LENGTH = 32
  private static readonly IV_LENGTH = 12
  private static readonly TAG_LENGTH = 16

  /**
   * Generate a secure random key for encryption
   */
  static generateKey(): string {
    return randomBytes(this.KEY_LENGTH).toString('hex')
  }

  /**
   * Generate a secure random salt
   */
  static generateSalt(): string {
    return randomBytes(16).toString('hex')
  }

  /**
   * Derive a key from password using PBKDF2
   */
  static async deriveKey(password: string, salt: string): Promise<CryptoKey> {
    const encoder = new TextEncoder()
    const keyMaterial = await crypto.subtle.importKey(
      'raw',
      encoder.encode(password) as BufferSource,
      { name: 'PBKDF2' },
      false,
      ['deriveBits', 'deriveKey']
    )

    return crypto.subtle.deriveKey(
      {
        name: 'PBKDF2',
        salt: encoder.encode(salt) as BufferSource,
        iterations: 100000,
        hash: 'SHA-256'
      },
      keyMaterial,
      { name: 'AES-GCM', length: 256 },
      true,
      ['encrypt', 'decrypt']
    )
  }

  /**
   * Encrypt data using AES-GCM
   */
  static async encrypt(data: string, key: CryptoKey): Promise<string> {
    const encoder = new TextEncoder()
    const iv = crypto.getRandomValues(new Uint8Array(this.IV_LENGTH))
    
    const encrypted = await crypto.subtle.encrypt(
      { name: 'AES-GCM', iv },
      key,
      encoder.encode(data) as BufferSource
    )

    // Combine IV and encrypted data
    const combined = new Uint8Array(iv.length + encrypted.byteLength)
    combined.set(iv)
    combined.set(new Uint8Array(encrypted), iv.length)

    return btoa(String.fromCharCode.apply(null, Array.from(combined)))
  }

  /**
   * Decrypt data using AES-GCM
   */
  static async decrypt(encryptedData: string, key: CryptoKey): Promise<string> {
    const combined = new Uint8Array(
      Array.from(atob(encryptedData)).map(char => char.charCodeAt(0))
    )

    const iv = combined.slice(0, this.IV_LENGTH)
    const encrypted = combined.slice(this.IV_LENGTH)

    const decrypted = await crypto.subtle.decrypt(
      { name: 'AES-GCM', iv },
      key,
      encrypted
    )

    return new TextDecoder().decode(decrypted)
  }

  /**
   * Hash data using SHA-256
   */
  static hash(data: string): string {
    return createHash('sha256').update(data).digest('hex')
  }

  /**
   * Validate password strength
   */
  static validatePassword(password: string): {
    isValid: boolean
    score: number
    feedback: string[]
  } {
    const feedback: string[] = []
    let score = 0

    // Length check
    if (password.length >= 12) {
      score += 2
    } else if (password.length >= 8) {
      score += 1
    } else {
      feedback.push('Password should be at least 8 characters long')
    }

    // Character variety checks
    if (/[a-z]/.test(password)) score += 1
    else feedback.push('Include lowercase letters')

    if (/[A-Z]/.test(password)) score += 1
    else feedback.push('Include uppercase letters')

    if (/\d/.test(password)) score += 1
    else feedback.push('Include numbers')

    if (/[^a-zA-Z\d]/.test(password)) score += 1
    else feedback.push('Include special characters')

    // Common patterns check
    if (!/(.)\1{2,}/.test(password)) score += 1
    else feedback.push('Avoid repeating characters')

    return {
      isValid: score >= 5,
      score: Math.min(score, 6),
      feedback
    }
  }

  /**
   * Generate a secure session token
   */
  static generateSessionToken(): string {
    return randomBytes(32).toString('hex')
  }

  /**
   * Check if running in secure context
   */
  static isSecureContext(): boolean {
    return typeof window !== 'undefined' && window.isSecureContext
  }

  /**
   * Sanitize storage key
   */
  static sanitizeKey(key: string): string {
    return this.STORAGE_PREFIX + key.replace(/[^a-zA-Z0-9_-]/g, '_')
  }
}

// Secure storage interface
export interface SecureStorageItem {
  data: string
  timestamp: number
  expiresAt?: number
  checksum: string
}

export class SecureStorage {
  private static readonly STORAGE_KEY_PREFIX = 'secure_'

  /**
   * Store encrypted data with integrity check
   */
  static async store(
    key: string, 
    data: any, 
    password: string, 
    expirationHours?: number
  ): Promise<void> {
    if (!WalletSecurity.isSecureContext()) {
      throw new Error('Secure storage requires HTTPS context')
    }

    const salt = WalletSecurity.generateSalt()
    const derivedKey = await WalletSecurity.deriveKey(password, salt)
    
    const serializedData = JSON.stringify(data)
    const encryptedData = await WalletSecurity.encrypt(serializedData, derivedKey)
    
    const storageItem: SecureStorageItem = {
      data: encryptedData,
      timestamp: Date.now(),
      expiresAt: expirationHours ? Date.now() + (expirationHours * 60 * 60 * 1000) : undefined,
      checksum: WalletSecurity.hash(serializedData)
    }

    const storageKey = WalletSecurity.sanitizeKey(key)
    localStorage.setItem(storageKey, JSON.stringify({ ...storageItem, salt }))
  }

  /**
   * Retrieve and decrypt data with integrity verification
   */
  static async retrieve(key: string, password: string): Promise<any> {
    if (!WalletSecurity.isSecureContext()) {
      throw new Error('Secure storage requires HTTPS context')
    }

    const storageKey = WalletSecurity.sanitizeKey(key)
    const stored = localStorage.getItem(storageKey)
    
    if (!stored) {
      throw new Error('Data not found')
    }

    const { data, timestamp, expiresAt, checksum, salt } = JSON.parse(stored)

    // Check expiration
    if (expiresAt && Date.now() > expiresAt) {
      this.remove(key)
      throw new Error('Data has expired')
    }

    const derivedKey = await WalletSecurity.deriveKey(password, salt)
    const decryptedData = await WalletSecurity.decrypt(data, derivedKey)
    
    // Verify integrity
    const computedChecksum = WalletSecurity.hash(decryptedData)
    if (computedChecksum !== checksum) {
      throw new Error('Data integrity check failed')
    }

    return JSON.parse(decryptedData)
  }

  /**
   * Remove data from storage
   */
  static remove(key: string): void {
    const storageKey = WalletSecurity.sanitizeKey(key)
    localStorage.removeItem(storageKey)
  }

  /**
   * Check if data exists and is valid
   */
  static exists(key: string): boolean {
    const storageKey = WalletSecurity.sanitizeKey(key)
    const stored = localStorage.getItem(storageKey)
    
    if (!stored) return false

    try {
      const { expiresAt } = JSON.parse(stored)
      return !expiresAt || Date.now() <= expiresAt
    } catch {
      return false
    }
  }

  /**
   * Clear all secure storage
   */
  static clearAll(): void {
    const keys = Object.keys(localStorage).filter(key => 
      key.startsWith(WalletSecurity.sanitizeKey(''))
    )
    keys.forEach(key => localStorage.removeItem(key))
  }
}

// Security event types
export enum SecurityEventType {
  WALLET_CONNECTED = 'wallet_connected',
  WALLET_DISCONNECTED = 'wallet_disconnected',
  TRANSACTION_SIGNED = 'transaction_signed',
  SECURITY_BREACH_ATTEMPT = 'security_breach_attempt',
  PASSWORD_CHANGED = 'password_changed',
  SESSION_EXPIRED = 'session_expired',
  SUSPICIOUS_ACTIVITY = 'suspicious_activity'
}

export interface SecurityEvent {
  type: SecurityEventType
  timestamp: number
  details: Record<string, any>
  severity: 'low' | 'medium' | 'high' | 'critical'
  userAgent?: string
  ipAddress?: string
}

export class SecurityLogger {
  private static events: SecurityEvent[] = []
  private static readonly MAX_EVENTS = 1000

  /**
   * Log a security event
   */
  static logEvent(
    type: SecurityEventType,
    details: Record<string, any>,
    severity: SecurityEvent['severity'] = 'medium'
  ): void {
    const event: SecurityEvent = {
      type,
      timestamp: Date.now(),
      details,
      severity,
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : undefined
    }

    this.events.unshift(event)
    
    // Keep only recent events
    if (this.events.length > this.MAX_EVENTS) {
      this.events = this.events.slice(0, this.MAX_EVENTS)
    }

    // Store in secure storage for persistence
    try {
      localStorage.setItem('security_events', JSON.stringify(this.events.slice(0, 100)))
    } catch (error) {
      console.warn('Failed to persist security events:', error)
    }

    // Alert on critical events
    if (severity === 'critical') {
      console.error('Critical security event:', event)
    }
  }

  /**
   * Get recent security events
   */
  static getEvents(limit = 50): SecurityEvent[] {
    return this.events.slice(0, limit)
  }

  /**
   * Get events by type
   */
  static getEventsByType(type: SecurityEventType): SecurityEvent[] {
    return this.events.filter(event => event.type === type)
  }

  /**
   * Clear security events
   */
  static clearEvents(): void {
    this.events = []
    localStorage.removeItem('security_events')
  }

  /**
   * Load events from storage
   */
  static loadEvents(): void {
    try {
      const stored = localStorage.getItem('security_events')
      if (stored) {
        this.events = JSON.parse(stored)
      }
    } catch (error) {
      console.warn('Failed to load security events:', error)
    }
  }
}
