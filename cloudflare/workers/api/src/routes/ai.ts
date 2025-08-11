/**
 * AI routes for Cloudflare Worker
 */

import { Router } from 'itty-router';
import { corsHeaders } from '../utils/cors';
import { ValidationError } from '../utils/errorHandler';
import { Env } from '../index';

const router = Router({ base: '/api/ai' });

// Chat with AI
router.post('/chat', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { message?: string; conversation_id?: string };
    const { message, conversation_id } = body;
    const user = (request as any).user;

    if (!message) {
      throw new ValidationError('Message is required');
    }

    // Create conversation if not exists
    let conversationId = conversation_id;
    if (!conversationId) {
      conversationId = crypto.randomUUID();
      await env.DB.prepare(
        'INSERT INTO conversations (id, user_id, title, created_at) VALUES (?, ?, ?, ?)'
      ).bind(
        conversationId,
        user.id,
        message.substring(0, 50) + '...',
        new Date().toISOString()
      ).run();
    }

    // Store user message
    const userMessageId = crypto.randomUUID();
    await env.DB.prepare(
      'INSERT INTO messages (id, conversation_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)'
    ).bind(
      userMessageId,
      conversationId,
      'user',
      message,
      new Date().toISOString()
    ).run();

    // Call AI service (OpenAI or Anthropic)
    const aiResponse = await callAIService(message, env);

    // Store AI response
    const aiMessageId = crypto.randomUUID();
    await env.DB.prepare(
      'INSERT INTO messages (id, conversation_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)'
    ).bind(
      aiMessageId,
      conversationId,
      'assistant',
      aiResponse,
      new Date().toISOString()
    ).run();

    return new Response(JSON.stringify({
      success: true,
      conversation_id: conversationId,
      message: aiResponse,
      timestamp: new Date().toISOString(),
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

// Get conversations
router.get('/conversations', async (request: Request, env: Env) => {
  try {
    const user = (request as any).user;
    
    const conversations = await env.DB.prepare(
      'SELECT id, title, created_at FROM conversations WHERE user_id = ? ORDER BY created_at DESC LIMIT 50'
    ).bind(user.id).all();

    return new Response(JSON.stringify({
      success: true,
      conversations: conversations.results,
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

// Get conversation messages
router.get('/conversations/:id', async (request: Request & { params?: { id: string } }, env: Env) => {
  try {
    const user = (request as any).user;
    const conversationId = request.params?.id;

    if (!conversationId) {
      throw new ValidationError('Conversation ID is required');
    }

    // Verify conversation belongs to user
    const conversation = await env.DB.prepare(
      'SELECT id FROM conversations WHERE id = ? AND user_id = ?'
    ).bind(conversationId, user.id).first();

    if (!conversation) {
      throw new ValidationError('Conversation not found');
    }

    // Get messages
    const messages = await env.DB.prepare(
      'SELECT id, role, content, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC'
    ).bind(conversationId).all();

    return new Response(JSON.stringify({
      success: true,
      conversation_id: conversationId,
      messages: messages.results,
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

// Analyze crypto data
router.post('/analyze', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { symbol?: string; timeframe?: string; analysis_type?: string };
    const { symbol, timeframe, analysis_type } = body;

    if (!symbol) {
      throw new ValidationError('Symbol is required');
    }

    // Get cached analysis if available
    const cacheKey = `analysis:${symbol}:${timeframe || '1h'}:${analysis_type || 'technical'}`;
    const cached = await env.CACHE.get(cacheKey);
    
    if (cached) {
      return new Response(cached, {
        headers: {
          'Content-Type': 'application/json',
          'X-Cache': 'HIT',
          ...corsHeaders,
        },
      });
    }

    // Perform analysis (simplified)
    const analysis = await performCryptoAnalysis(symbol, timeframe || '1h', analysis_type || 'technical', env);

    // Cache result for 5 minutes
    await env.CACHE.put(cacheKey, JSON.stringify(analysis), {
      expirationTtl: 300,
    });

    return new Response(JSON.stringify(analysis), {
      headers: {
        'Content-Type': 'application/json',
        'X-Cache': 'MISS',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// Predict price
router.post('/predict/price', async (request: Request, env: Env) => {
  try {
    const body = await request.json() as { symbol?: string; timeframe?: string };
    const { symbol, timeframe } = body;

    if (!symbol) {
      throw new ValidationError('Symbol is required');
    }

    // Get cached prediction if available
    const cacheKey = `prediction:${symbol}:${timeframe || '1h'}`;
    const cached = await env.CACHE.get(cacheKey);
    
    if (cached) {
      return new Response(cached, {
        headers: {
          'Content-Type': 'application/json',
          'X-Cache': 'HIT',
          ...corsHeaders,
        },
      });
    }

    // Generate prediction (simplified)
    const prediction = await generatePricePrediction(symbol, timeframe || '1h', env);

    // Cache result for 10 minutes
    await env.CACHE.put(cacheKey, JSON.stringify(prediction), {
      expirationTtl: 600,
    });

    return new Response(JSON.stringify(prediction), {
      headers: {
        'Content-Type': 'application/json',
        'X-Cache': 'MISS',
        ...corsHeaders,
      },
    });
  } catch (error) {
    throw error;
  }
});

// Helper functions
async function callAIService(message: string, env: Env): Promise<string> {
  try {
    // Use OpenAI API
    const response = await fetch('https://api.openai.com/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${env.OPENAI_API_KEY}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        model: 'gpt-3.5-turbo',
        messages: [
          {
            role: 'system',
            content: 'You are an AI assistant specialized in cryptocurrency trading and analysis. Provide helpful, accurate, and actionable insights.',
          },
          {
            role: 'user',
            content: message,
          },
        ],
        max_tokens: 500,
        temperature: 0.7,
      }),
    });

    const data = await response.json() as any;
    return data.choices?.[0]?.message?.content || 'Sorry, I could not process your request.';
  } catch (error) {
    console.error('AI service error:', error);
    return 'Sorry, the AI service is currently unavailable.';
  }
}

async function performCryptoAnalysis(symbol: string, timeframe: string, analysisType: string, env: Env): Promise<any> {
  // Simplified analysis - in production, this would integrate with real trading APIs
  return {
    success: true,
    symbol,
    timeframe: timeframe || '1h',
    analysis_type: analysisType || 'technical',
    timestamp: new Date().toISOString(),
    analysis: {
      trend: 'bullish',
      confidence: 0.75,
      support_levels: [45000, 44500, 44000],
      resistance_levels: [46000, 46500, 47000],
      indicators: {
        rsi: 65.5,
        macd: 'bullish_crossover',
        moving_averages: {
          sma_20: 45200,
          sma_50: 44800,
          ema_12: 45300,
        },
      },
      recommendation: 'HOLD',
    },
  };
}

async function generatePricePrediction(symbol: string, timeframe: string, env: Env): Promise<any> {
  // Simplified prediction - in production, this would use ML models
  const currentPrice = 45500; // Mock current price
  const prediction = currentPrice * (1 + (Math.random() - 0.5) * 0.1); // ±5% random prediction

  return {
    success: true,
    symbol,
    timeframe: timeframe || '1h',
    timestamp: new Date().toISOString(),
    prediction: {
      current_price: currentPrice,
      predicted_price: Math.round(prediction * 100) / 100,
      confidence: 0.68,
      change_percent: ((prediction - currentPrice) / currentPrice * 100).toFixed(2),
      factors: [
        'Technical indicators suggest continued momentum',
        'Market sentiment remains positive',
        'Volume analysis indicates strong support',
      ],
    },
  };
}

export { router as aiRoutes };
