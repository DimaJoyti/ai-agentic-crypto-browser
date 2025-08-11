/**
 * Analytics routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/analytics' });

// Get market analytics
router.get('/market', async (request: Request, env: Env) => {
  const analytics = {
    market_cap: 1200000000000,
    volume_24h: 45000000000,
    dominance: { BTC: 42.5, ETH: 18.2 },
    fear_greed_index: 65,
  };

  return new Response(JSON.stringify({ success: true, ...analytics }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as analyticsRoutes };
