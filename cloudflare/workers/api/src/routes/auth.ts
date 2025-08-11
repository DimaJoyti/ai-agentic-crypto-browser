/**
 * Authentication routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { ValidationError, UnauthorizedError } from '../utils/errorHandler';
import { Env } from '../index';

const router = Router({ base: '/api/auth' });

// User registration
router.post('/register', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { email?: string; password?: string; name?: string };
    const { email, password, name } = body;

    // Validate input
    if (!email || !password || !name) {
      throw new ValidationError('Email, password, and name are required');
    }

    if (password.length < 8) {
      throw new ValidationError('Password must be at least 8 characters long');
    }

    // Check if user already exists
    const existingUser = await env.DB.prepare(
      'SELECT id FROM users WHERE email = ?'
    ).bind(email).first();

    if (existingUser) {
      throw new ValidationError('User with this email already exists');
    }

    // Hash password
    const hashedPassword = await hashPassword(password);

    // Create user
    const userId = crypto.randomUUID();
    await env.DB.prepare(
      'INSERT INTO users (id, email, password_hash, name, role, created_at) VALUES (?, ?, ?, ?, ?, ?)'
    ).bind(
      userId,
      email,
      hashedPassword,
      name,
      'user',
      new Date().toISOString()
    ).run();

    // Generate JWT token
    const token = await generateJWT({
      sub: userId,
      email,
      role: 'user',
      permissions: ['read', 'write'],
    }, env.JWT_SECRET);

    // Store session
    const sessionKey = `session:${userId}:${token}`;
    await env.SESSIONS.put(sessionKey, JSON.stringify({
      userId,
      email,
      createdAt: new Date().toISOString(),
    }), {
      expirationTtl: 24 * 60 * 60, // 24 hours
    });

    return new Response(JSON.stringify({
      success: true,
      token,
      user: {
        id: userId,
        email,
        name,
        role: 'user',
      },
    }), {
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// User login
router.post('/login', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { email?: string; password?: string };
    const { email, password } = body;

    // Validate input
    if (!email || !password) {
      throw new ValidationError('Email and password are required');
    }

    // Get user from database
    const user = await env.DB.prepare(
      'SELECT id, email, password_hash, name, role FROM users WHERE email = ?'
    ).bind(email).first() as { id: string; email: string; password_hash: string; name: string; role: string } | null;

    if (!user) {
      throw new UnauthorizedError('Invalid email or password');
    }

    // Verify password
    const isValidPassword = await verifyPassword(password, user.password_hash);
    if (!isValidPassword) {
      throw new UnauthorizedError('Invalid email or password');
    }

    // Generate JWT token
    const token = await generateJWT({
      sub: user.id,
      email: user.email,
      role: user.role,
      permissions: ['read', 'write'],
    }, env.JWT_SECRET);

    // Store session
    const sessionKey = `session:${user.id}:${token}`;
    await env.SESSIONS.put(sessionKey, JSON.stringify({
      userId: user.id,
      email: user.email,
      loginAt: new Date().toISOString(),
    }), {
      expirationTtl: 24 * 60 * 60, // 24 hours
    });

    return new Response(JSON.stringify({
      success: true,
      token,
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        role: user.role,
      },
    }), {
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// User logout
router.post('/logout', async (request: Request, env: Env) => {
  try {
    const authHeader = request.headers.get('Authorization');
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7);
      const user = (request as any).user;
      
      if (user) {
        // Remove session
        const sessionKey = `session:${user.id}:${token}`;
        await env.SESSIONS.delete(sessionKey);
      }
    }

    return new Response(JSON.stringify({
      success: true,
      message: 'Logged out successfully',
    }), {
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// Get current user
router.get('/me', async (request: Request, env: Env) => {
  try {
    const user = (request as any).user;
    
    if (!user) {
      throw new UnauthorizedError('User not authenticated');
    }

    // Get full user details from database
    const userDetails = await env.DB.prepare(
      'SELECT id, email, name, role, created_at FROM users WHERE id = ?'
    ).bind(user.id).first();

    return new Response(JSON.stringify({
      success: true,
      user: userDetails,
    }), {
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// Helper functions
async function hashPassword(password: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(password);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return Array.from(new Uint8Array(hash))
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}

async function verifyPassword(password: string, hash: string): Promise<boolean> {
  const hashedInput = await hashPassword(password);
  return hashedInput === hash;
}

async function generateJWT(payload: any, secret: string): Promise<string> {
  const header = {
    alg: 'HS256',
    typ: 'JWT',
  };

  const now = Math.floor(Date.now() / 1000);
  const jwtPayload = {
    ...payload,
    iat: now,
    exp: now + (24 * 60 * 60), // 24 hours
  };

  const encoder = new TextEncoder();
  const headerB64 = btoa(JSON.stringify(header)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
  const payloadB64 = btoa(JSON.stringify(jwtPayload)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');

  const data = encoder.encode(`${headerB64}.${payloadB64}`);
  const secretKey = await crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign']
  );

  const signature = await crypto.subtle.sign('HMAC', secretKey, data);
  const signatureB64 = btoa(String.fromCharCode(...new Uint8Array(signature)))
    .replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');

  return `${headerB64}.${payloadB64}.${signatureB64}`;
}

export { router as authRoutes };
