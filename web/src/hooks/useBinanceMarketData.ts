import { useState, useEffect, useCallback, useRef } from 'react';
import { 
  binanceMarketService, 
  BinanceTicker, 
  BinanceOrderBook, 
  BinanceTrade,
  BinanceKline 
} from '@/services/binance/BinanceMarketService';
import { 
  binanceWebSocketService, 
  BinanceTickerStream, 
  BinanceOrderBookStream,
  BinanceTradeStream 
} from '@/services/binance/BinanceWebSocketService';

/**
 * Hook for real-time ticker data
 */
export function useBinanceTicker(symbol: string) {
  const [ticker, setTicker] = useState<BinanceTicker | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    let unsubscribe: (() => void) | null = null;

    const initializeTicker = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get initial ticker data
        const initialTicker = await binanceMarketService.getTicker(symbol) as BinanceTicker;
        setTicker(initialTicker);

        // Subscribe to real-time updates
        unsubscribe = binanceWebSocketService.subscribeTicker(symbol, (data: BinanceTickerStream) => {
          setTicker(prev => ({
            ...prev!,
            price: data.c,
            priceChange: data.p,
            priceChangePercent: data.P,
            lastPrice: data.c,
            openPrice: data.o,
            highPrice: data.h,
            lowPrice: data.l,
            volume: data.v,
            quoteVolume: data.q,
            closeTime: data.C,
            openTime: data.O,
            count: data.n
          }));
          setIsConnected(true);
        });

        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch ticker data');
        setIsLoading(false);
      }
    };

    initializeTicker();

    return () => {
      if (unsubscribe) {
        unsubscribe();
      }
      setIsConnected(false);
    };
  }, [symbol]);

  return { ticker, isLoading, error, isConnected };
}

/**
 * Hook for real-time order book data
 */
export function useBinanceOrderBook(symbol: string, limit: number = 20) {
  const [orderBook, setOrderBook] = useState<BinanceOrderBook | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    let unsubscribe: (() => void) | null = null;

    const initializeOrderBook = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get initial order book
        const initialOrderBook = await binanceMarketService.getOrderBook(symbol, limit);
        setOrderBook(initialOrderBook);

        // Subscribe to real-time updates
        unsubscribe = binanceWebSocketService.subscribeOrderBook(symbol, (data: BinanceOrderBookStream) => {
          setOrderBook(prev => {
            if (!prev) return prev;

            const newOrderBook = { ...prev };

            // Update bids
            data.b.forEach(([price, quantity]) => {
              const index = newOrderBook.bids.findIndex(bid => bid.price === price);
              if (parseFloat(quantity) === 0) {
                // Remove if quantity is 0
                if (index !== -1) {
                  newOrderBook.bids.splice(index, 1);
                }
              } else {
                // Update or add
                if (index !== -1) {
                  newOrderBook.bids[index].quantity = quantity;
                } else {
                  newOrderBook.bids.push({ price, quantity });
                  newOrderBook.bids.sort((a, b) => parseFloat(b.price) - parseFloat(a.price));
                }
              }
            });

            // Update asks
            data.a.forEach(([price, quantity]) => {
              const index = newOrderBook.asks.findIndex(ask => ask.price === price);
              if (parseFloat(quantity) === 0) {
                // Remove if quantity is 0
                if (index !== -1) {
                  newOrderBook.asks.splice(index, 1);
                }
              } else {
                // Update or add
                if (index !== -1) {
                  newOrderBook.asks[index].quantity = quantity;
                } else {
                  newOrderBook.asks.push({ price, quantity });
                  newOrderBook.asks.sort((a, b) => parseFloat(a.price) - parseFloat(b.price));
                }
              }
            });

            // Keep only the requested limit
            newOrderBook.bids = newOrderBook.bids.slice(0, limit);
            newOrderBook.asks = newOrderBook.asks.slice(0, limit);

            return newOrderBook;
          });
          setIsConnected(true);
        });

        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch order book');
        setIsLoading(false);
      }
    };

    initializeOrderBook();

    return () => {
      if (unsubscribe) {
        unsubscribe();
      }
      setIsConnected(false);
    };
  }, [symbol, limit]);

  return { orderBook, isLoading, error, isConnected };
}

/**
 * Hook for recent trades
 */
export function useBinanceTrades(symbol: string, limit: number = 100) {
  const [trades, setTrades] = useState<BinanceTrade[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    let unsubscribe: (() => void) | null = null;

    const initializeTrades = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get initial trades
        const initialTrades = await binanceMarketService.getRecentTrades(symbol, limit);
        setTrades(initialTrades.reverse()); // Most recent first

        // Subscribe to real-time trade updates
        unsubscribe = binanceWebSocketService.subscribeTrades(symbol, (data: BinanceTradeStream) => {
          setTrades(prev => {
            const newTrade: BinanceTrade = {
              id: data.t,
              price: data.p,
              qty: data.q,
              quoteQty: (parseFloat(data.p) * parseFloat(data.q)).toString(),
              time: data.T,
              isBuyerMaker: data.m,
              isBestMatch: true
            };

            const newTrades = [newTrade, ...prev];
            return newTrades.slice(0, limit); // Keep only the limit
          });
          setIsConnected(true);
        });

        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch trades');
        setIsLoading(false);
      }
    };

    initializeTrades();

    return () => {
      if (unsubscribe) {
        unsubscribe();
      }
      setIsConnected(false);
    };
  }, [symbol, limit]);

  return { trades, isLoading, error, isConnected };
}

/**
 * Hook for multiple tickers
 */
export function useBinanceMultipleTickers(symbols: string[]) {
  const [tickers, setTickers] = useState<Record<string, BinanceTicker>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshTickers = useCallback(async () => {
    try {
      setError(null);
      const allTickers = await binanceMarketService.getTicker() as BinanceTicker[];
      
      const filteredTickers = allTickers
        .filter(ticker => symbols.includes(ticker.symbol))
        .reduce((acc, ticker) => {
          acc[ticker.symbol] = ticker;
          return acc;
        }, {} as Record<string, BinanceTicker>);

      setTickers(filteredTickers);
      setIsLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tickers');
      setIsLoading(false);
    }
  }, [symbols]);

  useEffect(() => {
    refreshTickers();
    
    // Refresh every 5 seconds
    const interval = setInterval(refreshTickers, 5000);
    
    return () => clearInterval(interval);
  }, [refreshTickers]);

  return { tickers, isLoading, error, refreshTickers };
}

/**
 * Hook for candlestick/kline data
 */
export function useBinanceKlines(symbol: string, interval: string, limit: number = 500) {
  const [klines, setKlines] = useState<BinanceKline[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchKlines = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const data = await binanceMarketService.getKlines(symbol, interval, limit);
      setKlines(data);
      setIsLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch klines');
      setIsLoading(false);
    }
  }, [symbol, interval, limit]);

  useEffect(() => {
    fetchKlines();
  }, [fetchKlines]);

  return { klines, isLoading, error, refetch: fetchKlines };
}
