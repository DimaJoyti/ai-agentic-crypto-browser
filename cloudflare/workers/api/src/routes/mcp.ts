/**
 * MCP (Model Context Protocol) routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { Env } from '../index';

const router = Router({ base: '/api/mcp' });

// Get MCP status
router.get('/status', async (request: Request, env: Env) => {
  const status = {
    connected_tools: 12,
    active_sessions: 5,
    last_sync: new Date().toISOString(),
  };

  return new Response(JSON.stringify({ success: true, ...status }), {
    headers: { 'Content-Type': 'application/json', ...corsHeaders },
  });
});

export { router as mcpRoutes };
