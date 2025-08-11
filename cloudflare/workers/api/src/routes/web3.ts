/**
 * Web3 routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { ValidationError } from '../utils/errorHandler';
import { Env } from '../index';

const router = Router({ base: '/api/web3' });

// Connect wallet
router.post('/connect-wallet', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { address?: string; signature?: string; message?: string };
    const { address, signature, message } = body;
    const user = (request as any).user;

    if (!address || !signature || !message) {
      throw new ValidationError('Address, signature, and message are required');
    }

    // Store wallet connection
    await env.DB.prepare(
      'INSERT OR REPLACE INTO user_wallets (user_id, address, signature, connected_at) VALUES (?, ?, ?, ?)'
    ).bind(
      user.id,
      address,
      signature,
      new Date().toISOString()
    ).run();

    return new Response(JSON.stringify({
      success: true,
      message: 'Wallet connected successfully',
      address,
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

// Get wallet balance
router.get('/balance', async (request: Request, env: Env) => {
  try {
    const user = (request as any).user;
    const url = new URL(request.url);
    const address = url.searchParams.get('address');

    if (!address) {
      throw new ValidationError('Wallet address is required');
    }

    // Mock balance data - in production, integrate with real blockchain APIs
    const balance = {
      address,
      balances: [
        { symbol: 'ETH', balance: '1.5432', usd_value: 2456.78 },
        { symbol: 'USDC', balance: '1000.00', usd_value: 1000.00 },
        { symbol: 'BTC', balance: '0.0234', usd_value: 1234.56 },
      ],
      total_usd_value: 4691.34,
    };

    return new Response(JSON.stringify({
      success: true,
      ...balance,
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

// Send transaction
router.post('/transaction', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { to?: string; amount?: string; token?: string; gas_price?: string };
    const { to, amount, token, gas_price } = body;
    const user = (request as any).user;

    if (!to || !amount) {
      throw new ValidationError('Recipient address and amount are required');
    }

    // Mock transaction - in production, integrate with real blockchain
    const txHash = '0x' + Array.from(crypto.getRandomValues(new Uint8Array(32)))
      .map(b => b.toString(16).padStart(2, '0')).join('');

    // Store transaction record
    await env.DB.prepare(
      'INSERT INTO transactions (id, user_id, tx_hash, to_address, amount, token, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)'
    ).bind(
      crypto.randomUUID(),
      user.id,
      txHash,
      to,
      amount,
      token || 'ETH',
      'pending',
      new Date().toISOString()
    ).run();

    return new Response(JSON.stringify({
      success: true,
      tx_hash: txHash,
      status: 'pending',
      message: 'Transaction submitted successfully',
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

// Get DeFi positions
router.get('/defi/positions', async (request: Request, env: Env) => {
  try {
    const user = (request as any).user;
    const url = new URL(request.url);
    const address = url.searchParams.get('address');

    if (!address) {
      throw new ValidationError('Wallet address is required');
    }

    // Mock DeFi positions - in production, integrate with DeFi protocols
    const positions = {
      address,
      positions: [
        {
          protocol: 'Uniswap V3',
          type: 'liquidity_pool',
          pair: 'ETH/USDC',
          amount: '1000.00',
          apy: 12.5,
          rewards: '15.67',
        },
        {
          protocol: 'Compound',
          type: 'lending',
          asset: 'USDC',
          amount: '5000.00',
          apy: 4.2,
          rewards: '8.33',
        },
      ],
      total_value: 6000.00,
      total_rewards: 24.00,
    };

    return new Response(JSON.stringify({
      success: true,
      ...positions,
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

export { router as web3Routes };
