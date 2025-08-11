/**
 * Solana routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/solana' });

// Get Solana balance
router.get('/balance', async (request: Request, env: Env) => {
  const url = new URL(request.url);
  const address = url.searchParams.get('address');
  
  const balance = {
    address,
    sol_balance: 12.5,
    tokens: [
      { mint: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', symbol: 'USDC', balance: 1000.0 },
    ],
  };

  return new Response(JSON.stringify({ success: true, ...balance }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as solanaRoutes };
