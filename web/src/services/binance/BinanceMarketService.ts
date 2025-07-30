/**
 * Binance Market Data Service
 * Handles public market data from Binance API (no API keys required)
 */

export interface BinanceTicker {
  symbol: string;
  price: string;
  priceChange: string;
  priceChangePercent: string;
  weightedAvgPrice: string;
  prevClosePrice: string;
  lastPrice: string;
  lastQty: string;
  bidPrice: string;
  bidQty: string;
  askPrice: string;
  askQty: string;
  openPrice: string;
  highPrice: string;
  lowPrice: string;
  volume: string;
  quoteVolume: string;
  openTime: number;
  closeTime: number;
  firstId: number;
  lastId: number;
  count: number;
}

export interface BinanceOrderBookEntry {
  price: string;
  quantity: string;
}

export interface BinanceOrderBook {
  lastUpdateId: number;
  bids: BinanceOrderBookEntry[];
  asks: BinanceOrderBookEntry[];
}

export interface BinanceKline {
  openTime: number;
  open: string;
  high: string;
  low: string;
  close: string;
  volume: string;
  closeTime: number;
  quoteAssetVolume: string;
  numberOfTrades: number;
  takerBuyBaseAssetVolume: string;
  takerBuyQuoteAssetVolume: string;
}

export interface BinanceTrade {
  id: number;
  price: string;
  qty: string;
  quoteQty: string;
  time: number;
  isBuyerMaker: boolean;
  isBestMatch: boolean;
}

export class BinanceMarketService {
  private readonly baseUrl = 'https://api.binance.com/api/v3';
  private readonly wsBaseUrl = 'wss://stream.binance.com:9443/ws';

  /**
   * Get 24hr ticker price change statistics
   */
  async getTicker(symbol?: string): Promise<BinanceTicker | BinanceTicker[]> {
    const url = symbol 
      ? `${this.baseUrl}/ticker/24hr?symbol=${symbol}`
      : `${this.baseUrl}/ticker/24hr`;
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch ticker: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Get current average price for a symbol
   */
  async getAvgPrice(symbol: string): Promise<{ mins: number; price: string }> {
    const response = await fetch(`${this.baseUrl}/avgPrice?symbol=${symbol}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch average price: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Get order book (market depth)
   */
  async getOrderBook(symbol: string, limit: number = 100): Promise<BinanceOrderBook> {
    const response = await fetch(`${this.baseUrl}/depth?symbol=${symbol}&limit=${limit}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch order book: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Get recent trades
   */
  async getRecentTrades(symbol: string, limit: number = 500): Promise<BinanceTrade[]> {
    const response = await fetch(`${this.baseUrl}/trades?symbol=${symbol}&limit=${limit}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch recent trades: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Get kline/candlestick data
   */
  async getKlines(
    symbol: string, 
    interval: string, 
    limit: number = 500,
    startTime?: number,
    endTime?: number
  ): Promise<BinanceKline[]> {
    let url = `${this.baseUrl}/klines?symbol=${symbol}&interval=${interval}&limit=${limit}`;
    
    if (startTime) url += `&startTime=${startTime}`;
    if (endTime) url += `&endTime=${endTime}`;
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch klines: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    // Transform array format to object format
    return data.map((kline: any[]) => ({
      openTime: kline[0],
      open: kline[1],
      high: kline[2],
      low: kline[3],
      close: kline[4],
      volume: kline[5],
      closeTime: kline[6],
      quoteAssetVolume: kline[7],
      numberOfTrades: kline[8],
      takerBuyBaseAssetVolume: kline[9],
      takerBuyQuoteAssetVolume: kline[10]
    }));
  }

  /**
   * Get exchange information
   */
  async getExchangeInfo(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/exchangeInfo`);
    if (!response.ok) {
      throw new Error(`Failed to fetch exchange info: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Get server time
   */
  async getServerTime(): Promise<{ serverTime: number }> {
    const response = await fetch(`${this.baseUrl}/time`);
    if (!response.ok) {
      throw new Error(`Failed to fetch server time: ${response.statusText}`);
    }
    
    return response.json();
  }

  /**
   * Format symbol for Binance API (remove special characters)
   */
  formatSymbol(symbol: string): string {
    return symbol.replace(/[^a-zA-Z0-9]/g, '').toUpperCase();
  }

  /**
   * Get popular trading pairs
   */
  getPopularPairs(): string[] {
    return [
      'BTCUSDT',
      'ETHUSDT',
      'BNBUSDT',
      'ADAUSDT',
      'XRPUSDT',
      'SOLUSDT',
      'DOTUSDT',
      'LINKUSDT',
      'AVAXUSDT',
      'MATICUSDT',
      'ATOMUSDT',
      'NEARUSDT',
      'UNIUSDT',
      'LTCUSDT',
      'BCHUSDT'
    ];
  }
}

// Export singleton instance
export const binanceMarketService = new BinanceMarketService();
