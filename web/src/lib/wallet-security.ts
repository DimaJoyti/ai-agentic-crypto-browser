import { type Address, type Hex } from 'viem'
import { verifyMessage, recoverAddress } from 'viem'

export interface SecurityConfig {
  sessionTimeout: number // in milliseconds
  maxFailedAttempts: number
  lockoutDuration: number // in milliseconds
  requireSignatureAuth: boolean
  enableBiometric: boolean
  encryptionEnabled: boolean
}

export interface AuthenticationChallenge {
  id: string
  message: string
  timestamp: number
  expiresAt: number
  address: Address
  nonce: string
  domain: string
}

export interface SecuritySession {
  id: string
  address: Address
  walletId: string
  startTime: number
  lastActivity: number
  expiresAt: number
  isActive: boolean
  ipAddress?: string
  userAgent?: string
  permissions: string[]
  riskScore: number
}

export interface SecurityEvent {
  id: string
  type: 'login' | 'logout' | 'failed_auth' | 'suspicious_activity' | 'permission_change'
  address: Address
  timestamp: number
  details: Record<string, any>
  riskLevel: 'low' | 'medium' | 'high' | 'critical'
  resolved: boolean
}

export interface BiometricData {
  id: string
  type: 'fingerprint' | 'face' | 'voice'
  hash: string
  createdAt: number
  lastUsed: number
}

export class WalletSecurityManager {
  private static instance: WalletSecurityManager
  private config: SecurityConfig
  private activeSessions = new Map<string, SecuritySession>()
  private pendingChallenges = new Map<string, AuthenticationChallenge>()
  private securityEvents: SecurityEvent[] = []
  private failedAttempts = new Map<Address, { count: number; lastAttempt: number }>()
  private biometricData = new Map<Address, BiometricData[]>()
  private encryptionKey: CryptoKey | null = null

  private constructor(config: Partial<SecurityConfig> = {}) {
    this.config = {
      sessionTimeout: 24 * 60 * 60 * 1000, // 24 hours
      maxFailedAttempts: 5,
      lockoutDuration: 15 * 60 * 1000, // 15 minutes
      requireSignatureAuth: true,
      enableBiometric: false,
      encryptionEnabled: true,
      ...config
    }

    this.initializeEncryption()
    this.startSessionCleanup()
  }

  static getInstance(config?: Partial<SecurityConfig>): WalletSecurityManager {
    if (!WalletSecurityManager.instance) {
      WalletSecurityManager.instance = new WalletSecurityManager(config)
    }
    return WalletSecurityManager.instance
  }

  /**
   * Initialize encryption for secure data storage
   */
  private async initializeEncryption(): Promise<void> {
    if (!this.config.encryptionEnabled || !window.crypto?.subtle) return

    try {
      // Generate or retrieve encryption key
      const keyData = localStorage.getItem('wallet_security_key')
      if (keyData) {
        const keyBuffer = new Uint8Array(JSON.parse(keyData))
        this.encryptionKey = await window.crypto.subtle.importKey(
          'raw',
          keyBuffer,
          { name: 'AES-GCM' },
          false,
          ['encrypt', 'decrypt']
        )
      } else {
        this.encryptionKey = await window.crypto.subtle.generateKey(
          { name: 'AES-GCM', length: 256 },
          true,
          ['encrypt', 'decrypt']
        )
        
        const keyBuffer = await window.crypto.subtle.exportKey('raw', this.encryptionKey)
        localStorage.setItem('wallet_security_key', JSON.stringify(Array.from(new Uint8Array(keyBuffer))))
      }
    } catch (error) {
      console.warn('Failed to initialize encryption:', error)
      this.config.encryptionEnabled = false
    }
  }

  /**
   * Generate authentication challenge for wallet
   */
  generateChallenge(address: Address, domain: string = window.location.hostname): AuthenticationChallenge {
    const nonce = this.generateNonce()
    const timestamp = Date.now()
    const expiresAt = timestamp + (5 * 60 * 1000) // 5 minutes

    const challenge: AuthenticationChallenge = {
      id: `challenge_${timestamp}_${Math.random().toString(36).substr(2, 9)}`,
      message: this.createChallengeMessage(address, nonce, domain, timestamp),
      timestamp,
      expiresAt,
      address,
      nonce,
      domain
    }

    this.pendingChallenges.set(challenge.id, challenge)
    
    // Clean up expired challenges
    this.cleanupExpiredChallenges()

    return challenge
  }

  /**
   * Create challenge message for signing
   */
  private createChallengeMessage(address: Address, nonce: string, domain: string, timestamp: number): string {
    return `Welcome to ${domain}!

This request will not trigger a blockchain transaction or cost any gas fees.

Wallet address: ${address}
Nonce: ${nonce}
Issued at: ${new Date(timestamp).toISOString()}`
  }

  /**
   * Verify signature and authenticate wallet
   */
  async verifySignature(
    challengeId: string, 
    signature: Hex, 
    address: Address
  ): Promise<{ success: boolean; session?: SecuritySession; error?: string }> {
    const challenge = this.pendingChallenges.get(challengeId)
    
    if (!challenge) {
      return { success: false, error: 'Challenge not found or expired' }
    }

    if (Date.now() > challenge.expiresAt) {
      this.pendingChallenges.delete(challengeId)
      return { success: false, error: 'Challenge expired' }
    }

    if (challenge.address.toLowerCase() !== address.toLowerCase()) {
      return { success: false, error: 'Address mismatch' }
    }

    // Check if address is locked out
    if (this.isAddressLockedOut(address)) {
      return { success: false, error: 'Address temporarily locked due to failed attempts' }
    }

    try {
      // Verify the signature
      const isValid = await verifyMessage({
        address,
        message: challenge.message,
        signature
      })

      if (!isValid) {
        this.recordFailedAttempt(address)
        this.logSecurityEvent({
          type: 'failed_auth',
          address,
          details: { challengeId, reason: 'Invalid signature' },
          riskLevel: 'medium'
        })
        return { success: false, error: 'Invalid signature' }
      }

      // Clear failed attempts on successful auth
      this.failedAttempts.delete(address)
      
      // Create security session
      const session = this.createSession(address, challenge.domain)
      
      // Clean up challenge
      this.pendingChallenges.delete(challengeId)

      this.logSecurityEvent({
        type: 'login',
        address,
        details: { sessionId: session.id, domain: challenge.domain },
        riskLevel: 'low'
      })

      return { success: true, session }

    } catch (error) {
      this.recordFailedAttempt(address)
      this.logSecurityEvent({
        type: 'failed_auth',
        address,
        details: { challengeId, error: (error as Error).message },
        riskLevel: 'high'
      })
      return { success: false, error: 'Signature verification failed' }
    }
  }

  /**
   * Create a new security session
   */
  private createSession(address: Address, domain: string): SecuritySession {
    const now = Date.now()
    const session: SecuritySession = {
      id: `session_${now}_${Math.random().toString(36).substr(2, 9)}`,
      address,
      walletId: `wallet_${address}`,
      startTime: now,
      lastActivity: now,
      expiresAt: now + this.config.sessionTimeout,
      isActive: true,
      ipAddress: this.getClientIP(),
      userAgent: navigator.userAgent,
      permissions: ['read', 'write'], // Default permissions
      riskScore: this.calculateRiskScore(address)
    }

    this.activeSessions.set(session.id, session)
    return session
  }

  /**
   * Validate and refresh session
   */
  validateSession(sessionId: string): SecuritySession | null {
    const session = this.activeSessions.get(sessionId)
    
    if (!session || !session.isActive) {
      return null
    }

    const now = Date.now()
    
    if (now > session.expiresAt) {
      this.terminateSession(sessionId)
      return null
    }

    // Update last activity and extend session
    session.lastActivity = now
    session.expiresAt = now + this.config.sessionTimeout
    
    return session
  }

  /**
   * Terminate a session
   */
  terminateSession(sessionId: string): void {
    const session = this.activeSessions.get(sessionId)
    if (session) {
      session.isActive = false
      this.activeSessions.delete(sessionId)
      
      this.logSecurityEvent({
        type: 'logout',
        address: session.address,
        details: { sessionId, duration: Date.now() - session.startTime },
        riskLevel: 'low'
      })
    }
  }

  /**
   * Terminate all sessions for an address
   */
  terminateAllSessions(address: Address): void {
    const sessionsToTerminate: string[] = []
    
    this.activeSessions.forEach((session, sessionId) => {
      if (session.address.toLowerCase() === address.toLowerCase()) {
        sessionsToTerminate.push(sessionId)
      }
    })

    sessionsToTerminate.forEach(sessionId => this.terminateSession(sessionId))
  }

  /**
   * Get active sessions for an address
   */
  getActiveSessions(address: Address): SecuritySession[] {
    return Array.from(this.activeSessions.values())
      .filter(session => 
        session.address.toLowerCase() === address.toLowerCase() && 
        session.isActive
      )
  }

  /**
   * Check if address is locked out
   */
  private isAddressLockedOut(address: Address): boolean {
    const attempts = this.failedAttempts.get(address)
    if (!attempts) return false

    const now = Date.now()
    const timeSinceLastAttempt = now - attempts.lastAttempt

    if (timeSinceLastAttempt > this.config.lockoutDuration) {
      this.failedAttempts.delete(address)
      return false
    }

    return attempts.count >= this.config.maxFailedAttempts
  }

  /**
   * Record failed authentication attempt
   */
  private recordFailedAttempt(address: Address): void {
    const now = Date.now()
    const existing = this.failedAttempts.get(address)
    
    if (existing && (now - existing.lastAttempt) < this.config.lockoutDuration) {
      existing.count++
      existing.lastAttempt = now
    } else {
      this.failedAttempts.set(address, { count: 1, lastAttempt: now })
    }
  }

  /**
   * Calculate risk score for address
   */
  private calculateRiskScore(address: Address): number {
    let score = 0

    // Check failed attempts
    const attempts = this.failedAttempts.get(address)
    if (attempts) {
      score += attempts.count * 10
    }

    // Check recent security events
    const recentEvents = this.securityEvents
      .filter(event => 
        event.address.toLowerCase() === address.toLowerCase() &&
        Date.now() - event.timestamp < 24 * 60 * 60 * 1000 // Last 24 hours
      )

    score += recentEvents.length * 5

    // Check for high-risk events
    const highRiskEvents = recentEvents.filter(event => 
      event.riskLevel === 'high' || event.riskLevel === 'critical'
    )
    score += highRiskEvents.length * 20

    return Math.min(score, 100) // Cap at 100
  }

  /**
   * Log security event
   */
  private logSecurityEvent(event: Omit<SecurityEvent, 'id' | 'timestamp' | 'resolved'>): void {
    const securityEvent: SecurityEvent = {
      id: `event_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      timestamp: Date.now(),
      resolved: false,
      ...event
    }

    this.securityEvents.push(securityEvent)

    // Keep only last 1000 events
    if (this.securityEvents.length > 1000) {
      this.securityEvents = this.securityEvents.slice(-1000)
    }

    // Auto-resolve low-risk events after 24 hours
    if (event.riskLevel === 'low') {
      setTimeout(() => {
        securityEvent.resolved = true
      }, 24 * 60 * 60 * 1000)
    }
  }

  /**
   * Get security events for address
   */
  getSecurityEvents(address: Address, limit = 50): SecurityEvent[] {
    return this.securityEvents
      .filter(event => event.address.toLowerCase() === address.toLowerCase())
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit)
  }

  /**
   * Generate cryptographically secure nonce
   */
  private generateNonce(): string {
    const array = new Uint8Array(16)
    window.crypto.getRandomValues(array)
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('')
  }

  /**
   * Get client IP (placeholder - in real implementation, this would come from server)
   */
  private getClientIP(): string {
    // In a real implementation, this would be provided by the server
    return 'unknown'
  }

  /**
   * Clean up expired challenges
   */
  private cleanupExpiredChallenges(): void {
    const now = Date.now()
    const expiredChallenges: string[] = []

    this.pendingChallenges.forEach((challenge, id) => {
      if (now > challenge.expiresAt) {
        expiredChallenges.push(id)
      }
    })

    expiredChallenges.forEach(id => this.pendingChallenges.delete(id))
  }

  /**
   * Start periodic session cleanup
   */
  private startSessionCleanup(): void {
    setInterval(() => {
      const now = Date.now()
      const expiredSessions: string[] = []

      this.activeSessions.forEach((session, sessionId) => {
        if (now > session.expiresAt || !session.isActive) {
          expiredSessions.push(sessionId)
        }
      })

      expiredSessions.forEach(sessionId => this.terminateSession(sessionId))
      this.cleanupExpiredChallenges()
    }, 60000) // Check every minute
  }

  /**
   * Encrypt sensitive data
   */
  async encryptData(data: string): Promise<string> {
    if (!this.config.encryptionEnabled || !this.encryptionKey) {
      return data
    }

    try {
      const encoder = new TextEncoder()
      const dataBuffer = encoder.encode(data)
      const iv = window.crypto.getRandomValues(new Uint8Array(12))
      
      const encryptedBuffer = await window.crypto.subtle.encrypt(
        { name: 'AES-GCM', iv },
        this.encryptionKey,
        dataBuffer as BufferSource
      )

      const encryptedArray = new Uint8Array(encryptedBuffer)
      const result = new Uint8Array(iv.length + encryptedArray.length)
      result.set(iv)
      result.set(encryptedArray, iv.length)

      return btoa(String.fromCharCode(...Array.from(result)))
    } catch (error) {
      console.warn('Encryption failed:', error)
      return data
    }
  }

  /**
   * Decrypt sensitive data
   */
  async decryptData(encryptedData: string): Promise<string> {
    if (!this.config.encryptionEnabled || !this.encryptionKey) {
      return encryptedData
    }

    try {
      const dataArray = new Uint8Array(
        atob(encryptedData).split('').map(char => char.charCodeAt(0))
      )
      
      const iv = dataArray.slice(0, 12)
      const encrypted = dataArray.slice(12)

      const decryptedBuffer = await window.crypto.subtle.decrypt(
        { name: 'AES-GCM', iv },
        this.encryptionKey,
        encrypted
      )

      const decoder = new TextDecoder()
      return decoder.decode(decryptedBuffer)
    } catch (error) {
      console.warn('Decryption failed:', error)
      return encryptedData
    }
  }

  /**
   * Get security statistics
   */
  getSecurityStats(): {
    activeSessions: number
    pendingChallenges: number
    totalEvents: number
    highRiskEvents: number
    lockedAddresses: number
  } {
    const now = Date.now()
    const highRiskEvents = this.securityEvents.filter(event => 
      event.riskLevel === 'high' || event.riskLevel === 'critical'
    ).length

    const lockedAddresses = Array.from(this.failedAttempts.entries())
      .filter(([_, attempts]) => 
        attempts.count >= this.config.maxFailedAttempts &&
        (now - attempts.lastAttempt) < this.config.lockoutDuration
      ).length

    return {
      activeSessions: this.activeSessions.size,
      pendingChallenges: this.pendingChallenges.size,
      totalEvents: this.securityEvents.length,
      highRiskEvents,
      lockedAddresses
    }
  }

  /**
   * Update security configuration
   */
  updateConfig(newConfig: Partial<SecurityConfig>): void {
    this.config = { ...this.config, ...newConfig }
  }

  /**
   * Clear all security data (for testing/reset)
   */
  clearAllData(): void {
    this.activeSessions.clear()
    this.pendingChallenges.clear()
    this.securityEvents.length = 0
    this.failedAttempts.clear()
    this.biometricData.clear()
  }
}

// Export singleton instance
export const walletSecurity = WalletSecurityManager.getInstance()
