import React, { useState, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useBinanceOrderBook, useBinanceTicker } from '@/hooks/useBinanceMarketData';
import {
  Settings,
  TrendingUp,
  TrendingDown,
  Volume2,
  ArrowUpDown,
  Zap,
  Wifi,
  WifiOff,
  Loader2
} from 'lucide-react';

interface OrderBookEntry {
  price: string;
  size: string;
  total: string;
  percentage: number;
}

interface OrderBookProps {
  symbol: string;
  className?: string;
  limit?: number;
}

export const OrderBook: React.FC<OrderBookProps> = ({ symbol, className, limit = 20 }) => {
  const [precision, setPrecision] = useState(2);
  const { orderBook, isLoading, error, isConnected } = useBinanceOrderBook(symbol, limit);
  const { ticker } = useBinanceTicker(symbol);

  // Process order book data for display
  const processedOrderBook = useMemo(() => {
    if (!orderBook) return { asks: [], bids: [], spread: 0, spreadPercent: 0 };

    // Calculate totals and percentages
    let bidTotal = 0;
    let askTotal = 0;

    const processedBids = orderBook.bids.map((bid, index) => {
      bidTotal += parseFloat(bid.quantity);
      return {
        price: parseFloat(bid.price).toFixed(precision),
        size: parseFloat(bid.quantity).toFixed(4),
        total: (parseFloat(bid.price) * parseFloat(bid.quantity)).toFixed(2),
        percentage: 0 // Will be calculated after
      };
    });

    const processedAsks = orderBook.asks.map((ask, index) => {
      askTotal += parseFloat(ask.quantity);
      return {
        price: parseFloat(ask.price).toFixed(precision),
        size: parseFloat(ask.quantity).toFixed(4),
        total: (parseFloat(ask.price) * parseFloat(ask.quantity)).toFixed(2),
        percentage: 0 // Will be calculated after
      };
    });

    // Calculate percentages based on max volume
    const maxVolume = Math.max(bidTotal, askTotal);

    processedBids.forEach((bid, index) => {
      const cumulative = processedBids.slice(0, index + 1).reduce((sum, b) => sum + parseFloat(b.size), 0);
      bid.percentage = (cumulative / maxVolume) * 100;
    });

    processedAsks.forEach((ask, index) => {
      const cumulative = processedAsks.slice(0, index + 1).reduce((sum, a) => sum + parseFloat(a.size), 0);
      ask.percentage = (cumulative / maxVolume) * 100;
    });

    // Calculate spread
    const bestBid = orderBook.bids[0] ? parseFloat(orderBook.bids[0].price) : 0;
    const bestAsk = orderBook.asks[0] ? parseFloat(orderBook.asks[0].price) : 0;
    const spread = bestAsk - bestBid;
    const spreadPercent = bestBid > 0 ? (spread / bestBid) * 100 : 0;

    return {
      asks: processedAsks,
      bids: processedBids,
      spread,
      spreadPercent
    };
  }, [orderBook, precision]);

  if (isLoading) {
    return (
      <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm text-white flex items-center">
            <Loader2 className="w-4 h-4 mr-2 text-blue-400 animate-spin" />
            Loading Order Book
          </CardTitle>
        </CardHeader>
        <CardContent className="p-4">
          <div className="animate-pulse space-y-2">
            {[...Array(10)].map((_, i) => (
              <div key={i} className="h-4 bg-gray-700 rounded"></div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm text-white flex items-center">
            <WifiOff className="w-4 h-4 mr-2 text-red-400" />
            Order Book Error
          </CardTitle>
        </CardHeader>
        <CardContent className="p-4">
          <div className="text-center text-red-400">
            <p className="text-sm">Failed to load order book</p>
            <p className="text-xs text-gray-500 mt-1">{error}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm text-white flex items-center">
            <ArrowUpDown className="w-4 h-4 mr-2 text-blue-400" />
            Order Book
          </CardTitle>
          <div className="flex items-center space-x-2">
            <Badge variant="outline" className="text-xs text-gray-400 border-gray-600">
              {symbol}
            </Badge>
            <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white">
              <Settings className="w-3 h-3" />
            </Button>
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="space-y-1 p-4">
        {/* Header */}
        <div className="flex justify-between text-xs text-gray-400 mb-2 px-1">
          <span>Price (USDT)</span>
          <span>Size (BTC)</span>
          <span>Total</span>
        </div>

        {/* Asks (Sell Orders) */}
        <div className="space-y-1">
          {processedOrderBook.asks.slice().reverse().map((ask, index) => (
            <div
              key={`ask-${index}`}
              className="relative flex justify-between text-xs py-1.5 px-1 hover:bg-red-900/20 rounded cursor-pointer transition-colors"
            >
              {/* Background bar */}
              <div
                className="absolute right-0 top-0 h-full bg-red-500/10 rounded"
                style={{ width: `${ask.percentage}%` }}
              />

              <span className="text-red-400 font-mono relative z-10">{ask.price}</span>
              <span className="text-gray-300 font-mono relative z-10">{ask.size}</span>
              <span className="text-gray-500 font-mono relative z-10">{ask.total}</span>
            </div>
          ))}
        </div>

        {/* Spread */}
        <div className="flex justify-center py-3 border-y border-gray-700 my-2">
          <div className="text-center">
            <div className="text-xs text-gray-400">Spread</div>
            <div className="text-sm font-mono text-white">
              ${processedOrderBook.spread.toFixed(2)} ({processedOrderBook.spreadPercent.toFixed(3)}%)
            </div>
          </div>
        </div>

        {/* Bids (Buy Orders) */}
        <div className="space-y-1">
          {processedOrderBook.bids.map((bid, index) => (
            <div
              key={`bid-${index}`}
              className="relative flex justify-between text-xs py-1.5 px-1 hover:bg-green-900/20 rounded cursor-pointer transition-colors"
            >
              {/* Background bar */}
              <div
                className="absolute right-0 top-0 h-full bg-green-500/10 rounded"
                style={{ width: `${bid.percentage}%` }}
              />

              <span className="text-green-400 font-mono relative z-10">{bid.price}</span>
              <span className="text-gray-300 font-mono relative z-10">{bid.size}</span>
              <span className="text-gray-500 font-mono relative z-10">{bid.total}</span>
            </div>
          ))}
        </div>

        {/* Market Summary */}
        {ticker && (
          <div className="mt-4 pt-3 border-t border-gray-700">
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div className="flex items-center justify-between">
                <span className="text-gray-400">24h Volume</span>
                <span className="text-white font-mono">
                  {parseFloat(ticker.volume).toLocaleString(undefined, { maximumFractionDigits: 2 })}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-400">24h High</span>
                <span className="text-green-400 font-mono">
                  {parseFloat(ticker.highPrice).toFixed(precision)}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-400">24h Low</span>
                <span className="text-red-400 font-mono">
                  {parseFloat(ticker.lowPrice).toFixed(precision)}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-400">24h Change</span>
                <span className={`font-mono ${parseFloat(ticker.priceChangePercent) >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                  {parseFloat(ticker.priceChangePercent) >= 0 ? '+' : ''}{parseFloat(ticker.priceChangePercent).toFixed(2)}%
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Quick Actions */}
        <div className="mt-4 pt-3 border-t border-gray-700">
          <div className="flex space-x-2">
            <Button 
              size="sm" 
              className="flex-1 bg-green-600 hover:bg-green-700 text-white text-xs"
            >
              <TrendingUp className="w-3 h-3 mr-1" />
              Buy
            </Button>
            <Button 
              size="sm" 
              className="flex-1 bg-red-600 hover:bg-red-700 text-white text-xs"
            >
              <TrendingDown className="w-3 h-3 mr-1" />
              Sell
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
