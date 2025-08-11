/**
 * AI Agentic Crypto Browser - Cloudflare Worker API Gateway
 * Handles all API requests and routes them to appropriate handlers
 */

/// <reference path="./types/cloudflare.d.ts" />

import { Router } from 'itty-router';
import { corsHeaders, handleCORS } from './utils/cors';
import { authMiddleware } from './middleware/auth';
import { rateLimitMiddleware } from './middleware/rateLimit';
import { errorHandler } from './utils/errorHandler';
import { authRoutes } from './routes/auth';
import { aiRoutes } from './routes/ai';
import { web3Routes } from './routes/web3';
import { tradingRoutes } from './routes/trading';
import { analyticsRoutes } from './routes/analytics';
import { portfolioRoutes } from './routes/portfolio';
import { riskRoutes } from './routes/risk';
import { hftRoutes } from './routes/hft';
import { mcpRoutes } from './routes/mcp';
import { solanaRoutes } from './routes/solana';

export interface Env {
  // KV Namespaces
  CACHE: KVNamespace;
  SESSIONS: KVNamespace;

  // D1 Database
  DB: D1Database;

  // R2 Storage
  STORAGE: R2Bucket;

  // Durable Objects
  WEBSOCKET_HANDLER: DurableObjectNamespace;
  
  // Environment variables
  ENVIRONMENT: string;
  API_VERSION: string;
  JWT_SECRET: string;
  OPENAI_API_KEY: string;
  ANTHROPIC_API_KEY: string;
  ETHEREUM_RPC_URL: string;
  POLYGON_RPC_URL: string;
  BINANCE_API_KEY: string;
  BINANCE_SECRET_KEY: string;
}

// Create router
const router = Router();

// Health check endpoint
router.get('/health', () => {
  return new Response(JSON.stringify({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  }), {
    headers: {
      'Content-Type': 'application/json',
      ...corsHeaders
    }
  });
});

// API version endpoint
router.get('/api/version', () => {
  return new Response(JSON.stringify({
    version: '1.0.0',
    api_version: 'v1',
    environment: 'production'
  }), {
    headers: {
      'Content-Type': 'application/json',
      ...corsHeaders
    }
  });
});

// Apply middleware
router.all('*', rateLimitMiddleware);
router.all('/api/*', authMiddleware);

// Mount route handlers
router.all('/api/auth/*', authRoutes.handle);
router.all('/api/ai/*', aiRoutes.handle);
router.all('/api/web3/*', web3Routes.handle);
router.all('/api/trading/*', tradingRoutes.handle);
router.all('/api/analytics/*', analyticsRoutes.handle);
router.all('/api/portfolio/*', portfolioRoutes.handle);
router.all('/api/risk/*', riskRoutes.handle);
router.all('/api/hft/*', hftRoutes.handle);
router.all('/api/mcp/*', mcpRoutes.handle);
router.all('/api/solana/*', solanaRoutes.handle);

// 404 handler
router.all('*', () => {
  return new Response(JSON.stringify({
    error: 'Not Found',
    message: 'The requested endpoint does not exist'
  }), {
    status: 404,
    headers: {
      'Content-Type': 'application/json',
      ...corsHeaders
    }
  });
});

// Main worker handler
export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    try {
      // Handle CORS preflight requests
      if (request.method === 'OPTIONS') {
        return handleCORS(request);
      }

      // Route the request
      return router.handle(request, env, ctx);
    } catch (error) {
      return errorHandler(error);
    }
  }
};

// Export Durable Object classes
export { WebSocketHandler } from './durable-objects/websocket';
