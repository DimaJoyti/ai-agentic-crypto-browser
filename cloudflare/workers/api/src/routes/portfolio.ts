/**
 * Portfolio routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/portfolio' });

// Get portfolio
router.get('/', async (request: Request, env: Env) => {
  const user = (request as any).user;
  
  const portfolio = {
    total_value: 10000.00,
    holdings: [
      { symbol: 'BTC', amount: 0.25, value: 11375.00 },
      { symbol: 'ETH', amount: 3.2, value: 8960.00 },
    ],
    performance: { daily: 2.5, weekly: -1.2, monthly: 15.8 },
  };

  return new Response(JSON.stringify({ success: true, ...portfolio }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as portfolioRoutes };
