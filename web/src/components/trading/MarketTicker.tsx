import React, { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useBinanceTicker } from '@/hooks/useBinanceMarketData';
import {
  TrendingUp,
  TrendingDown,
  Activity,
  Volume2,
  Clock,
  Zap,
  Wifi,
  WifiOff
} from 'lucide-react';

interface MarketTickerProps {
  symbol: string;
  className?: string;
}

export const MarketTicker: React.FC<MarketTickerProps> = ({ symbol, className }) => {
  const { ticker, isLoading, error, isConnected } = useBinanceTicker(symbol);
  const [priceDirection, setPriceDirection] = useState<'up' | 'down' | 'neutral'>('neutral');
  const [lastPrice, setLastPrice] = useState<string>('');

  // Track price changes for direction indicator
  useEffect(() => {
    if (ticker?.lastPrice && lastPrice) {
      const current = parseFloat(ticker.lastPrice);
      const previous = parseFloat(lastPrice);

      if (current > previous) {
        setPriceDirection('up');
      } else if (current < previous) {
        setPriceDirection('down');
      } else {
        setPriceDirection('neutral');
      }
    }

    if (ticker?.lastPrice) {
      setLastPrice(ticker.lastPrice);
    }
  }, [ticker?.lastPrice, lastPrice]);

  // Reset direction after animation
  useEffect(() => {
    if (priceDirection !== 'neutral') {
      const timer = setTimeout(() => {
        setPriceDirection('neutral');
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [priceDirection]);

  if (isLoading) {
    return (
      <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
        <CardContent className="p-4">
          <div className="animate-pulse">
            <div className="h-6 bg-gray-700 rounded mb-2"></div>
            <div className="h-8 bg-gray-700 rounded mb-2"></div>
            <div className="h-4 bg-gray-700 rounded"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
        <CardContent className="p-4">
          <div className="text-center text-red-400">
            <WifiOff className="w-6 h-6 mx-auto mb-2" />
            <p className="text-sm">Failed to load market data</p>
            <p className="text-xs text-gray-500">{error}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!ticker) return null;

  const isPositive = parseFloat(ticker.priceChangePercent) >= 0;
  const formatPrice = (price: string) => {
    return parseFloat(price).toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 8
    });
  };

  const formatVolume = (volume: string) => {
    const val = parseFloat(volume);
    if (val >= 1000000) return `${(val / 1000000).toFixed(2)}M`;
    if (val >= 1000) return `${(val / 1000).toFixed(2)}K`;
    return val.toFixed(2);
  };

  return (
    <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
      <CardContent className="p-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            {isConnected ? (
              <Wifi className="w-4 h-4 text-green-400" />
            ) : (
              <WifiOff className="w-4 h-4 text-red-400" />
            )}
            <span className="text-lg font-bold text-white">{ticker.symbol}</span>
            <Badge
              variant="outline"
              className={`text-xs ${isConnected ? 'text-green-400 border-green-400' : 'text-red-400 border-red-400'}`}
            >
              {isConnected ? 'LIVE' : 'OFFLINE'}
            </Badge>
          </div>

          <div className="flex items-center space-x-1">
            <Activity className={`w-4 h-4 ${isConnected ? 'text-green-400' : 'text-gray-400'}`} />
            <span className="text-xs text-gray-400">
              {new Date().toLocaleTimeString()}
            </span>
          </div>
        </div>

        {/* Price Display */}
        <div className="mb-4">
          <div className="flex items-center space-x-3">
            <span
              className={`text-3xl font-bold transition-all duration-300 ${
                priceDirection === 'up' ? 'text-green-400 scale-105' :
                priceDirection === 'down' ? 'text-red-400 scale-105' : 'text-white scale-100'
              }`}
            >
              ${formatPrice(ticker.lastPrice)}
            </span>
            <div className="flex items-center space-x-1">
              {isPositive ? (
                <TrendingUp className="w-5 h-5 text-green-400" />
              ) : (
                <TrendingDown className="w-5 h-5 text-red-400" />
              )}
              <Badge
                className={`${isPositive ? 'bg-green-600' : 'bg-red-600'} text-white`}
              >
                {isPositive ? '+' : ''}{parseFloat(ticker.priceChangePercent).toFixed(2)}%
              </Badge>
            </div>
          </div>

          <div className="text-sm text-gray-400 mt-1">
            {isPositive ? '+' : ''}{formatPrice(ticker.priceChange)} USDT (24h)
          </div>
        </div>

        {/* Market Stats Grid */}
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-400">24h High</span>
              <span className="text-green-400 font-mono">${formatPrice(ticker.highPrice)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">24h Low</span>
              <span className="text-red-400 font-mono">${formatPrice(ticker.lowPrice)}</span>
            </div>
          </div>

          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-400">24h Volume</span>
              <span className="text-white font-mono">{formatVolume(ticker.volume)} {ticker.symbol.replace('USDT', '')}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Quote Volume</span>
              <span className="text-gray-300 font-mono">{formatVolume(ticker.quoteVolume)} USDT</span>
            </div>
          </div>
        </div>

        {/* Price Movement Indicator */}
        <div className="mt-4 pt-3 border-t border-gray-700">
          <div className="flex items-center justify-between text-xs">
            <span className="text-gray-400">Price Movement</span>
            <div className="flex items-center space-x-2">
              <div className={`w-2 h-2 rounded-full animate-pulse ${
                priceDirection === 'up' ? 'bg-green-400' :
                priceDirection === 'down' ? 'bg-red-400' : 'bg-gray-400'
              }`} />
              <span className={`font-mono ${
                priceDirection === 'up' ? 'text-green-400' :
                priceDirection === 'down' ? 'text-red-400' : 'text-gray-400'
              }`}>
                {priceDirection === 'up' ? '↗' : priceDirection === 'down' ? '↘' : '→'}
              </span>
            </div>
          </div>
        </div>

        {/* Quick Stats */}
        <div className="mt-3 grid grid-cols-3 gap-2 text-xs">
          <div className="text-center p-2 bg-gray-800/50 rounded">
            <div className="text-gray-400">Bid</div>
            <div className="text-green-400 font-mono">${formatPrice(ticker.bidPrice)}</div>
          </div>
          <div className="text-center p-2 bg-gray-800/50 rounded">
            <div className="text-gray-400">Ask</div>
            <div className="text-red-400 font-mono">${formatPrice(ticker.askPrice)}</div>
          </div>
          <div className="text-center p-2 bg-gray-800/50 rounded">
            <div className="text-gray-400">Spread</div>
            <div className="text-white font-mono">
              {(parseFloat(ticker.askPrice) - parseFloat(ticker.bidPrice)).toFixed(2)}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
