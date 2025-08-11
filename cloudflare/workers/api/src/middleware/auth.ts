/**
 * Authentication middleware for Cloudflare Worker
 */

import { UnauthorizedError } from '../utils/errorHandler';
import { Env } from '../index';

export interface User {
  id: string;
  email: string;
  role: string;
  permissions: string[];
}

export interface AuthenticatedRequest extends Request {
  user?: User;
}

// Public endpoints that don't require authentication
const PUBLIC_ENDPOINTS = [
  '/api/auth/login',
  '/api/auth/register',
  '/api/auth/refresh',
  '/api/health',
  '/api/version',
];

export async function authMiddleware(
  request: Request,
  env: Env,
  _ctx: ExecutionContext
): Promise<Response | void> {
  const url = new URL(request.url);
  const pathname = url.pathname;

  // Skip authentication for public endpoints
  if (PUBLIC_ENDPOINTS.some(endpoint => pathname.startsWith(endpoint))) {
    return;
  }

  // Get authorization header
  const authHeader = request.headers.get('Authorization');
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    throw new UnauthorizedError('Missing or invalid authorization header');
  }

  const token = authHeader.substring(7); // Remove 'Bearer ' prefix

  try {
    // Verify JWT token
    const user = await verifyJWT(token, env.JWT_SECRET);
    
    // Add user to request context
    (request as AuthenticatedRequest).user = user;
    
    // Check if user session is valid in KV
    const sessionKey = `session:${user.id}:${token}`;
    const session = await env.SESSIONS.get(sessionKey);
    
    if (!session) {
      throw new UnauthorizedError('Session expired or invalid');
    }

    // Update session expiry
    await env.SESSIONS.put(sessionKey, JSON.stringify({
      userId: user.id,
      lastAccess: new Date().toISOString(),
    }), {
      expirationTtl: 24 * 60 * 60, // 24 hours
    });

  } catch (error) {
    console.error('Authentication error:', error);
    throw new UnauthorizedError('Invalid or expired token');
  }
}

async function verifyJWT(token: string, secret: string): Promise<User> {
  try {
    // Import the Web Crypto API
    const encoder = new TextEncoder();
    const secretKey = await crypto.subtle.importKey(
      'raw',
      encoder.encode(secret),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['verify']
    );

    // Split the JWT token
    const [headerB64, payloadB64, signatureB64] = token.split('.');
    
    if (!headerB64 || !payloadB64 || !signatureB64) {
      throw new Error('Invalid JWT format');
    }

    // Verify signature
    const data = encoder.encode(`${headerB64}.${payloadB64}`);
    const signature = base64UrlDecode(signatureB64);
    
    const isValid = await crypto.subtle.verify(
      'HMAC',
      secretKey,
      signature as BufferSource,
      data
    );

    if (!isValid) {
      throw new Error('Invalid JWT signature');
    }

    // Decode payload
    const payload = JSON.parse(atob(payloadB64.replace(/-/g, '+').replace(/_/g, '/')));
    
    // Check expiration
    if (payload.exp && payload.exp < Date.now() / 1000) {
      throw new Error('JWT token expired');
    }

    return {
      id: payload.sub,
      email: payload.email,
      role: payload.role || 'user',
      permissions: payload.permissions || [],
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    throw new Error(`JWT verification failed: ${errorMessage}`);
  }
}

function base64UrlDecode(str: string): Uint8Array {
  // Add padding if needed
  str += '='.repeat((4 - str.length % 4) % 4);
  // Replace URL-safe characters
  str = str.replace(/-/g, '+').replace(/_/g, '/');
  // Decode base64
  const binary = atob(str);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}
