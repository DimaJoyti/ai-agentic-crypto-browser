/**
 * Rate limiting middleware for Cloudflare Worker
 */

import { RateLimitError } from '../utils/errorHandler';
import { Env } from '../index';

interface RateLimitConfig {
  windowMs: number; // Time window in milliseconds
  maxRequests: number; // Maximum requests per window
  keyGenerator?: (request: Request) => string;
}

const DEFAULT_CONFIG: RateLimitConfig = {
  windowMs: 60 * 1000, // 1 minute
  maxRequests: 100, // 100 requests per minute
};

// Different rate limits for different endpoints
const ENDPOINT_CONFIGS: Record<string, RateLimitConfig> = {
  '/api/auth/login': {
    windowMs: 15 * 60 * 1000, // 15 minutes
    maxRequests: 5, // 5 login attempts per 15 minutes
  },
  '/api/auth/register': {
    windowMs: 60 * 60 * 1000, // 1 hour
    maxRequests: 3, // 3 registrations per hour
  },
  '/api/ai/': {
    windowMs: 60 * 1000, // 1 minute
    maxRequests: 20, // 20 AI requests per minute
  },
  '/api/trading/': {
    windowMs: 60 * 1000, // 1 minute
    maxRequests: 50, // 50 trading requests per minute
  },
  '/api/web3/': {
    windowMs: 60 * 1000, // 1 minute
    maxRequests: 30, // 30 Web3 requests per minute
  },
};

export async function rateLimitMiddleware(
  request: Request,
  env: Env,
  ctx: ExecutionContext
): Promise<Response | void> {
  const url = new URL(request.url);
  const pathname = url.pathname;

  // Skip rate limiting for health checks
  if (pathname === '/health' || pathname === '/api/version') {
    return;
  }

  // Get rate limit configuration for this endpoint
  const config = getConfigForEndpoint(pathname);
  
  // Generate rate limit key
  const key = generateRateLimitKey(request, config);
  
  // Check current rate limit status
  const currentCount = await getCurrentCount(env.CACHE, key, config.windowMs);
  
  if (currentCount >= config.maxRequests) {
    // Rate limit exceeded
    const resetTime = await getResetTime(env.CACHE, key);
    
    throw new RateLimitError(`Rate limit exceeded. Try again in ${Math.ceil((resetTime - Date.now()) / 1000)} seconds`);
  }

  // Increment counter
  await incrementCounter(env.CACHE, key, config.windowMs);
  
  // Add rate limit headers to response (will be added by the response handler)
  const remaining = config.maxRequests - currentCount - 1;
  const resetTime = Date.now() + config.windowMs;
  
  // Store rate limit info in request context for response headers
  (request as any).rateLimitInfo = {
    limit: config.maxRequests,
    remaining: Math.max(0, remaining),
    reset: Math.floor(resetTime / 1000),
  };
}

function getConfigForEndpoint(pathname: string): RateLimitConfig {
  // Find matching endpoint configuration
  for (const [endpoint, config] of Object.entries(ENDPOINT_CONFIGS)) {
    if (pathname.startsWith(endpoint)) {
      return config;
    }
  }
  
  return DEFAULT_CONFIG;
}

function generateRateLimitKey(request: Request, config: RateLimitConfig): string {
  if (config.keyGenerator) {
    return config.keyGenerator(request);
  }
  
  // Default key generation based on IP address
  const clientIP = request.headers.get('CF-Connecting-IP') || 
                   request.headers.get('X-Forwarded-For') || 
                   'unknown';
  
  const url = new URL(request.url);
  const endpoint = url.pathname.split('/').slice(0, 3).join('/'); // e.g., /api/auth
  
  return `rate_limit:${clientIP}:${endpoint}`;
}

async function getCurrentCount(cache: KVNamespace, key: string, windowMs: number): Promise<number> {
  const data = await cache.get(key);
  
  if (!data) {
    return 0;
  }
  
  try {
    const parsed = JSON.parse(data);
    const now = Date.now();
    
    // Check if the window has expired
    if (now - parsed.firstRequest > windowMs) {
      // Window expired, reset count
      await cache.delete(key);
      return 0;
    }
    
    return parsed.count || 0;
  } catch (error) {
    console.error('Error parsing rate limit data:', error);
    return 0;
  }
}

async function incrementCounter(cache: KVNamespace, key: string, windowMs: number): Promise<void> {
  const now = Date.now();
  const data = await cache.get(key);
  
  let rateLimitData = {
    count: 1,
    firstRequest: now,
  };
  
  if (data) {
    try {
      const parsed = JSON.parse(data);
      
      // Check if window has expired
      if (now - parsed.firstRequest <= windowMs) {
        rateLimitData = {
          count: (parsed.count || 0) + 1,
          firstRequest: parsed.firstRequest,
        };
      }
    } catch (error) {
      console.error('Error parsing existing rate limit data:', error);
    }
  }
  
  // Store with TTL slightly longer than window to handle clock skew
  const ttl = Math.ceil(windowMs / 1000) + 60; // Add 60 seconds buffer
  
  await cache.put(key, JSON.stringify(rateLimitData), {
    expirationTtl: ttl,
  });
}

async function getResetTime(cache: KVNamespace, key: string): Promise<number> {
  const data = await cache.get(key);
  
  if (!data) {
    return Date.now();
  }
  
  try {
    const parsed = JSON.parse(data);
    return parsed.firstRequest + DEFAULT_CONFIG.windowMs;
  } catch (error) {
    return Date.now();
  }
}
