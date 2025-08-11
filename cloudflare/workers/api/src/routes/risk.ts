/**
 * Risk management routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/risk' });

// Get risk assessment
router.get('/assessment', async (request: Request, env: Env) => {
  const assessment = {
    risk_score: 6.5,
    risk_level: 'moderate',
    factors: ['Market volatility', 'Portfolio concentration'],
    recommendations: ['Diversify holdings', 'Set stop losses'],
  };

  return new Response(JSON.stringify({ success: true, ...assessment }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as riskRoutes };
