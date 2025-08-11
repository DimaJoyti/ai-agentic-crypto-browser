/**
 * CORS utilities for Cloudflare Worker
 */

export const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With, X-Session-ID',
  'Access-Control-Max-Age': '86400',
};

export function handleCORS(request: Request): Response {
  // Handle CORS preflight requests
  if (request.method === 'OPTIONS') {
    return new Response(null, {
      status: 200,
      headers: corsHeaders,
    });
  }

  return new Response('Method not allowed', {
    status: 405,
    headers: corsHeaders,
  });
}

export function addCORSHeaders(response: Response): Response {
  // Add CORS headers to existing response
  const newHeaders = new Headers(response.headers);
  Object.entries(corsHeaders).forEach(([key, value]) => {
    newHeaders.set(key, value);
  });

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: newHeaders,
  });
}
