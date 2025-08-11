/**
 * Cloudflare KV utilities for caching and session management
 * Replaces Redis functionality with KV storage
 */

export interface CacheOptions {
  expirationTtl?: number; // TTL in seconds
  metadata?: any;
}

export interface SessionData {
  userId: string;
  email: string;
  role: string;
  permissions: string[];
  createdAt: string;
  lastAccess: string;
  ipAddress?: string;
  userAgent?: string;
}

export interface CacheEntry<T = any> {
  data: T;
  timestamp: number;
  ttl?: number;
  metadata?: any;
}

/**
 * KV Cache Manager
 * Provides Redis-like functionality using Cloudflare KV
 */
export class KVCacheManager {
  constructor(private kv: KVNamespace) {}

  /**
   * Set a value in cache with optional TTL
   */
  async set<T>(key: string, value: T, options?: CacheOptions): Promise<void> {
    const cacheEntry: CacheEntry<T> = {
      data: value,
      timestamp: Date.now(),
      ttl: options?.expirationTtl,
      metadata: options?.metadata,
    };

    const kvOptions: any = {};
    if (options?.expirationTtl) {
      kvOptions.expirationTtl = options.expirationTtl;
    }
    if (options?.metadata) {
      kvOptions.metadata = options.metadata;
    }

    await this.kv.put(key, JSON.stringify(cacheEntry), kvOptions);
  }

  /**
   * Get a value from cache
   */
  async get<T>(key: string): Promise<T | null> {
    const value = await this.kv.get(key);
    if (!value) {
      return null;
    }

    try {
      const cacheEntry: CacheEntry<T> = JSON.parse(value);
      
      // Check if entry has expired (additional check beyond KV TTL)
      if (cacheEntry.ttl && cacheEntry.timestamp + (cacheEntry.ttl * 1000) < Date.now()) {
        await this.delete(key);
        return null;
      }

      return cacheEntry.data;
    } catch (error) {
      console.error('Error parsing cached value:', error);
      return null;
    }
  }

  /**
   * Delete a value from cache
   */
  async delete(key: string): Promise<void> {
    await this.kv.delete(key);
  }

  /**
   * Check if a key exists in cache
   */
  async exists(key: string): Promise<boolean> {
    const value = await this.kv.get(key);
    return value !== null;
  }

  /**
   * Increment a numeric value (atomic operation simulation)
   */
  async increment(key: string, delta: number = 1): Promise<number> {
    const current = await this.get<number>(key) || 0;
    const newValue = current + delta;
    await this.set(key, newValue);
    return newValue;
  }

  /**
   * Set multiple values at once
   */
  async mset(entries: Record<string, any>, options?: CacheOptions): Promise<void> {
    const promises = Object.entries(entries).map(([key, value]) =>
      this.set(key, value, options)
    );
    await Promise.all(promises);
  }

  /**
   * Get multiple values at once
   */
  async mget<T>(keys: string[]): Promise<(T | null)[]> {
    const promises = keys.map(key => this.get<T>(key));
    return Promise.all(promises);
  }

  /**
   * List keys with a prefix
   */
  async listKeys(prefix?: string, limit?: number): Promise<string[]> {
    const options: any = {};
    if (prefix) options.prefix = prefix;
    if (limit) options.limit = limit;

    const result = await this.kv.list(options);
    return result.keys.map(key => key.name);
  }
}

/**
 * Session Manager using KV storage
 */
export class KVSessionManager {
  private cache: KVCacheManager;
  private readonly SESSION_PREFIX = 'session:';
  private readonly USER_SESSIONS_PREFIX = 'user_sessions:';
  private readonly DEFAULT_TTL = 24 * 60 * 60; // 24 hours

  constructor(private kv: KVNamespace) {
    this.cache = new KVCacheManager(kv);
  }

  /**
   * Create a new session
   */
  async createSession(sessionId: string, sessionData: SessionData, ttl?: number): Promise<void> {
    const sessionKey = this.SESSION_PREFIX + sessionId;
    const userSessionsKey = this.USER_SESSIONS_PREFIX + sessionData.userId;

    // Store session data
    await this.cache.set(sessionKey, sessionData, {
      expirationTtl: ttl || this.DEFAULT_TTL,
    });

    // Track user sessions for cleanup
    const userSessions = await this.cache.get<string[]>(userSessionsKey) || [];
    userSessions.push(sessionId);
    await this.cache.set(userSessionsKey, userSessions, {
      expirationTtl: (ttl || this.DEFAULT_TTL) + 3600, // Slightly longer TTL
    });
  }

  /**
   * Get session data
   */
  async getSession(sessionId: string): Promise<SessionData | null> {
    const sessionKey = this.SESSION_PREFIX + sessionId;
    return this.cache.get<SessionData>(sessionKey);
  }

  /**
   * Update session data
   */
  async updateSession(sessionId: string, updates: Partial<SessionData>, ttl?: number): Promise<void> {
    const sessionKey = this.SESSION_PREFIX + sessionId;
    const existingSession = await this.getSession(sessionId);
    
    if (!existingSession) {
      throw new Error('Session not found');
    }

    const updatedSession: SessionData = {
      ...existingSession,
      ...updates,
      lastAccess: new Date().toISOString(),
    };

    await this.cache.set(sessionKey, updatedSession, {
      expirationTtl: ttl || this.DEFAULT_TTL,
    });
  }

  /**
   * Delete a session
   */
  async deleteSession(sessionId: string): Promise<void> {
    const sessionKey = this.SESSION_PREFIX + sessionId;
    const session = await this.getSession(sessionId);
    
    if (session) {
      // Remove from user sessions list
      const userSessionsKey = this.USER_SESSIONS_PREFIX + session.userId;
      const userSessions = await this.cache.get<string[]>(userSessionsKey) || [];
      const updatedSessions = userSessions.filter(id => id !== sessionId);
      
      if (updatedSessions.length > 0) {
        await this.cache.set(userSessionsKey, updatedSessions);
      } else {
        await this.cache.delete(userSessionsKey);
      }
    }

    await this.cache.delete(sessionKey);
  }

  /**
   * Delete all sessions for a user
   */
  async deleteUserSessions(userId: string): Promise<void> {
    const userSessionsKey = this.USER_SESSIONS_PREFIX + userId;
    const userSessions = await this.cache.get<string[]>(userSessionsKey) || [];

    // Delete all user sessions
    const deletePromises = userSessions.map(sessionId => 
      this.cache.delete(this.SESSION_PREFIX + sessionId)
    );
    await Promise.all(deletePromises);

    // Delete user sessions list
    await this.cache.delete(userSessionsKey);
  }

  /**
   * Get all active sessions for a user
   */
  async getUserSessions(userId: string): Promise<SessionData[]> {
    const userSessionsKey = this.USER_SESSIONS_PREFIX + userId;
    const sessionIds = await this.cache.get<string[]>(userSessionsKey) || [];

    const sessions = await Promise.all(
      sessionIds.map(sessionId => this.getSession(sessionId))
    );

    return sessions.filter(session => session !== null) as SessionData[];
  }

  /**
   * Validate and refresh session
   */
  async validateSession(sessionId: string): Promise<SessionData | null> {
    const session = await this.getSession(sessionId);
    
    if (!session) {
      return null;
    }

    // Update last access time
    await this.updateSession(sessionId, {
      lastAccess: new Date().toISOString(),
    });

    return session;
  }
}

/**
 * Rate Limiter using KV storage
 */
export class KVRateLimiter {
  private cache: KVCacheManager;

  constructor(private kv: KVNamespace) {
    this.cache = new KVCacheManager(kv);
  }

  /**
   * Check if request is within rate limit
   */
  async isAllowed(key: string, limit: number, windowMs: number): Promise<{
    allowed: boolean;
    remaining: number;
    resetTime: number;
  }> {
    const now = Date.now();
    const windowStart = now - windowMs;
    
    // Get current request count
    const requestData = await this.cache.get<{
      count: number;
      firstRequest: number;
    }>(key);

    let count = 0;
    let firstRequest = now;

    if (requestData) {
      // Check if window has expired
      if (requestData.firstRequest > windowStart) {
        count = requestData.count;
        firstRequest = requestData.firstRequest;
      }
    }

    const allowed = count < limit;
    const remaining = Math.max(0, limit - count - (allowed ? 1 : 0));
    const resetTime = firstRequest + windowMs;

    if (allowed) {
      // Increment counter
      await this.cache.set(key, {
        count: count + 1,
        firstRequest,
      }, {
        expirationTtl: Math.ceil(windowMs / 1000) + 60, // Add buffer
      });
    }

    return {
      allowed,
      remaining,
      resetTime,
    };
  }
}
