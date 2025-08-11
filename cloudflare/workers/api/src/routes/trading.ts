/**
 * Trading routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/trading' });

// Get trading pairs
router.get('/pairs', async (request: Request, env: Env) => {
  const pairs = [
    { symbol: 'BTCUSDT', price: 45500.00, change: 2.5 },
    { symbol: 'ETHUSDT', price: 2800.00, change: -1.2 },
    { symbol: 'ADAUSDT', price: 0.45, change: 5.8 },
  ];

  return new Response(JSON.stringify({ success: true, pairs }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

// Place order
router.post('/orders', async (request: Request, env: Env) => {
  const body = await request.json();
  const orderId = crypto.randomUUID();
  
  return new Response(JSON.stringify({
    success: true,
    order_id: orderId,
    status: 'pending',
  }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as tradingRoutes };
