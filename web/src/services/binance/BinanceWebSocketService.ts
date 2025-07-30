/**
 * Binance WebSocket Service
 * Handles real-time market data streams from Binance
 */

export interface BinanceTickerStream {
  e: string; // Event type
  E: number; // Event time
  s: string; // Symbol
  c: string; // Close price
  o: string; // Open price
  h: string; // High price
  l: string; // Low price
  v: string; // Total traded base asset volume
  q: string; // Total traded quote asset volume
  P: string; // Price change percent
  p: string; // Price change
  Q: string; // Last quantity
  F: number; // First trade ID
  L: number; // Last trade ID
  C: number; // Close time
  O: number; // Open time
  n: number; // Total number of trades
}

export interface BinanceOrderBookStream {
  e: string; // Event type
  E: number; // Event time
  s: string; // Symbol
  U: number; // First update ID in event
  u: number; // Final update ID in event
  b: [string, string][]; // Bids to be updated
  a: [string, string][]; // Asks to be updated
}

export interface BinanceTradeStream {
  e: string; // Event type
  E: number; // Event time
  s: string; // Symbol
  t: number; // Trade ID
  p: string; // Price
  q: string; // Quantity
  b: number; // Buyer order ID
  a: number; // Seller order ID
  T: number; // Trade time
  m: boolean; // Is the buyer the market maker?
  M: boolean; // Ignore
}

export interface BinanceKlineStream {
  e: string; // Event type
  E: number; // Event time
  s: string; // Symbol
  k: {
    t: number; // Kline start time
    T: number; // Kline close time
    s: string; // Symbol
    i: string; // Interval
    f: number; // First trade ID
    L: number; // Last trade ID
    o: string; // Open price
    c: string; // Close price
    h: string; // High price
    l: string; // Low price
    v: string; // Base asset volume
    n: number; // Number of trades
    x: boolean; // Is this kline closed?
    q: string; // Quote asset volume
    V: string; // Taker buy base asset volume
    Q: string; // Taker buy quote asset volume
    B: string; // Ignore
  };
}

export type StreamCallback<T> = (data: T) => void;

export class BinanceWebSocketService {
  private connections: Map<string, WebSocket> = new Map();
  private callbacks: Map<string, Set<StreamCallback<any>>> = new Map();
  private reconnectAttempts: Map<string, number> = new Map();
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second

  /**
   * Subscribe to ticker stream for a symbol
   */
  subscribeTicker(symbol: string, callback: StreamCallback<BinanceTickerStream>): () => void {
    const streamName = `${symbol.toLowerCase()}@ticker`;
    return this.subscribe(streamName, callback);
  }

  /**
   * Subscribe to order book depth stream
   */
  subscribeOrderBook(symbol: string, callback: StreamCallback<BinanceOrderBookStream>): () => void {
    const streamName = `${symbol.toLowerCase()}@depth`;
    return this.subscribe(streamName, callback);
  }

  /**
   * Subscribe to trade stream
   */
  subscribeTrades(symbol: string, callback: StreamCallback<BinanceTradeStream>): () => void {
    const streamName = `${symbol.toLowerCase()}@trade`;
    return this.subscribe(streamName, callback);
  }

  /**
   * Subscribe to kline/candlestick stream
   */
  subscribeKline(symbol: string, interval: string, callback: StreamCallback<BinanceKlineStream>): () => void {
    const streamName = `${symbol.toLowerCase()}@kline_${interval}`;
    return this.subscribe(streamName, callback);
  }

  /**
   * Subscribe to multiple streams
   */
  subscribeMultiple(streams: string[], callback: StreamCallback<any>): () => void {
    const streamName = streams.join('/');
    return this.subscribe(streamName, callback);
  }

  /**
   * Generic subscribe method
   */
  private subscribe<T>(streamName: string, callback: StreamCallback<T>): () => void {
    // Add callback to the set
    if (!this.callbacks.has(streamName)) {
      this.callbacks.set(streamName, new Set());
    }
    this.callbacks.get(streamName)!.add(callback);

    // Create connection if it doesn't exist
    if (!this.connections.has(streamName)) {
      this.createConnection(streamName);
    }

    // Return unsubscribe function
    return () => {
      this.unsubscribe(streamName, callback);
    };
  }

  /**
   * Unsubscribe from a stream
   */
  private unsubscribe<T>(streamName: string, callback: StreamCallback<T>): void {
    const callbacks = this.callbacks.get(streamName);
    if (callbacks) {
      callbacks.delete(callback);
      
      // If no more callbacks, close the connection
      if (callbacks.size === 0) {
        this.closeConnection(streamName);
      }
    }
  }

  /**
   * Create WebSocket connection
   */
  private createConnection(streamName: string): void {
    const wsUrl = `wss://stream.binance.com:9443/ws/${streamName}`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log(`Connected to Binance stream: ${streamName}`);
      this.reconnectAttempts.set(streamName, 0);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const callbacks = this.callbacks.get(streamName);
        if (callbacks) {
          callbacks.forEach(callback => callback(data));
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = (event) => {
      console.log(`Disconnected from Binance stream: ${streamName}`, event.code, event.reason);
      this.connections.delete(streamName);
      
      // Attempt to reconnect if there are still callbacks
      const callbacks = this.callbacks.get(streamName);
      if (callbacks && callbacks.size > 0) {
        this.attemptReconnect(streamName);
      }
    };

    ws.onerror = (error) => {
      console.error(`WebSocket error for stream ${streamName}:`, error);
    };

    this.connections.set(streamName, ws);
  }

  /**
   * Attempt to reconnect with exponential backoff
   */
  private attemptReconnect(streamName: string): void {
    const attempts = this.reconnectAttempts.get(streamName) || 0;
    
    if (attempts < this.maxReconnectAttempts) {
      const delay = this.reconnectDelay * Math.pow(2, attempts);
      
      setTimeout(() => {
        console.log(`Attempting to reconnect to ${streamName} (attempt ${attempts + 1})`);
        this.reconnectAttempts.set(streamName, attempts + 1);
        this.createConnection(streamName);
      }, delay);
    } else {
      console.error(`Max reconnection attempts reached for stream: ${streamName}`);
      this.callbacks.delete(streamName);
      this.reconnectAttempts.delete(streamName);
    }
  }

  /**
   * Close connection
   */
  private closeConnection(streamName: string): void {
    const ws = this.connections.get(streamName);
    if (ws) {
      ws.close();
      this.connections.delete(streamName);
    }
    this.callbacks.delete(streamName);
    this.reconnectAttempts.delete(streamName);
  }

  /**
   * Close all connections
   */
  closeAll(): void {
    this.connections.forEach((ws, streamName) => {
      ws.close();
    });
    this.connections.clear();
    this.callbacks.clear();
    this.reconnectAttempts.clear();
  }

  /**
   * Get connection status
   */
  getConnectionStatus(streamName: string): string {
    const ws = this.connections.get(streamName);
    if (!ws) return 'disconnected';
    
    switch (ws.readyState) {
      case WebSocket.CONNECTING: return 'connecting';
      case WebSocket.OPEN: return 'connected';
      case WebSocket.CLOSING: return 'closing';
      case WebSocket.CLOSED: return 'closed';
      default: return 'unknown';
    }
  }

  /**
   * Get all active connections
   */
  getActiveConnections(): string[] {
    return Array.from(this.connections.keys());
  }
}

// Export singleton instance
export const binanceWebSocketService = new BinanceWebSocketService();
