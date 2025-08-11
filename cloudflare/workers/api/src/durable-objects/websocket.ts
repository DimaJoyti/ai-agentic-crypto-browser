/**
 * WebSocket Durable Object for real-time communication
 */

/// <reference path="../types/cloudflare.d.ts" />

import { Env } from '../index';

export class WebSocketHandler {
  private sessions: Map<string, WebSocket> = new Map();
  private state: DurableObjectState;
  private env: Env;

  constructor(state: DurableObjectState, env: Env) {
    this.state = state;
    this.env = env;
  }

  async fetch(request: Request): Promise<Response> {
    const url = new URL(request.url);
    
    if (url.pathname === '/websocket') {
      // Handle WebSocket upgrade
      const upgradeHeader = request.headers.get('Upgrade');
      if (upgradeHeader !== 'websocket') {
        return new Response('Expected Upgrade: websocket', { status: 426 });
      }

      const webSocketPair = new WebSocketPair();
      const client = webSocketPair[0];
      const server = webSocketPair[1];

      // Accept the WebSocket connection
      server.accept();

      // Generate session ID
      const sessionId = crypto.randomUUID();
      this.sessions.set(sessionId, server);

      // Handle WebSocket events
      server.addEventListener('message', (event: MessageEvent) => {
        this.handleMessage(sessionId, event.data);
      });

      server.addEventListener('close', () => {
        this.sessions.delete(sessionId);
      });

      // Send welcome message
      server.send(JSON.stringify({
        type: 'welcome',
        sessionId,
        timestamp: new Date().toISOString(),
      }));

      return new Response(null, {
        status: 101,
        webSocket: client,
      } as ResponseInit & { webSocket: WebSocket });
    }

    return new Response('Not found', { status: 404 });
  }

  private handleMessage(sessionId: string, data: any) {
    try {
      const message = JSON.parse(data);
      
      switch (message.type) {
        case 'subscribe':
          this.handleSubscribe(sessionId, message);
          break;
        case 'unsubscribe':
          this.handleUnsubscribe(sessionId, message);
          break;
        case 'ping':
          this.handlePing(sessionId);
          break;
        default:
          console.log('Unknown message type:', message.type);
      }
    } catch (error) {
      console.error('Error handling WebSocket message:', error);
    }
  }

  private handleSubscribe(sessionId: string, message: any) {
    const socket = this.sessions.get(sessionId);
    if (!socket) return;

    // Subscribe to channels (prices, trades, etc.)
    socket.send(JSON.stringify({
      type: 'subscribed',
      channel: message.channel,
      timestamp: new Date().toISOString(),
    }));

    // Start sending mock data for demo
    this.startMockDataStream(sessionId, message.channel);
  }

  private handleUnsubscribe(sessionId: string, message: any) {
    const socket = this.sessions.get(sessionId);
    if (!socket) return;

    socket.send(JSON.stringify({
      type: 'unsubscribed',
      channel: message.channel,
      timestamp: new Date().toISOString(),
    }));
  }

  private handlePing(sessionId: string) {
    const socket = this.sessions.get(sessionId);
    if (!socket) return;

    socket.send(JSON.stringify({
      type: 'pong',
      timestamp: new Date().toISOString(),
    }));
  }

  private startMockDataStream(sessionId: string, channel: string) {
    const socket = this.sessions.get(sessionId);
    if (!socket) return;

    // Send mock price updates every 5 seconds
    const interval = setInterval(() => {
      if (!this.sessions.has(sessionId)) {
        clearInterval(interval);
        return;
      }

      const mockPrice = 45000 + (Math.random() - 0.5) * 1000;
      
      socket.send(JSON.stringify({
        type: 'price_update',
        channel,
        data: {
          symbol: 'BTCUSDT',
          price: mockPrice.toFixed(2),
          change: ((Math.random() - 0.5) * 5).toFixed(2),
          timestamp: new Date().toISOString(),
        },
      }));
    }, 5000);
  }

  // Broadcast message to all connected sessions
  broadcast(message: any) {
    const messageStr = JSON.stringify(message);
    for (const [sessionId, socket] of this.sessions) {
      try {
        socket.send(messageStr);
      } catch (error) {
        console.error(`Error sending to session ${sessionId}:`, error);
        this.sessions.delete(sessionId);
      }
    }
  }
}
