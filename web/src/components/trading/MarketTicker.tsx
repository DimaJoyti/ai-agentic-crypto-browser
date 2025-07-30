import React, { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  TrendingUp, 
  TrendingDown,
  Activity,
  Volume2,
  Clock,
  Zap
} from 'lucide-react';

interface MarketData {
  symbol: string;
  price: string;
  change24h: string;
  changePercent24h: string;
  high24h: string;
  low24h: string;
  volume24h: string;
  lastUpdate: string;
}

interface MarketTickerProps {
  symbol: string;
  className?: string;
}

export const MarketTicker: React.FC<MarketTickerProps> = ({ symbol, className }) => {
  const [isLive, setIsLive] = useState(true);
  const [lastPrice, setLastPrice] = useState('43,251.50');
  const [priceDirection, setPriceDirection] = useState<'up' | 'down' | 'neutral'>('up');

  // Mock market data
  const marketData: MarketData = {
    symbol: symbol,
    price: '43,251.50',
    change24h: '+987.50',
    changePercent24h: '+2.34',
    high24h: '44,125.00',
    low24h: '42,850.00',
    volume24h: '1,234.56',
    lastUpdate: new Date().toLocaleTimeString()
  };

  // Simulate price updates
  useEffect(() => {
    if (!isLive) return;

    const interval = setInterval(() => {
      const change = (Math.random() - 0.5) * 10;
      const newPrice = (parseFloat(lastPrice.replace(/,/g, '')) + change).toFixed(2);
      const formattedPrice = parseFloat(newPrice).toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      });
      
      setPriceDirection(change > 0 ? 'up' : change < 0 ? 'down' : 'neutral');
      setLastPrice(formattedPrice);
    }, 2000);

    return () => clearInterval(interval);
  }, [isLive, lastPrice]);

  const isPositive = marketData.changePercent24h.startsWith('+');

  return (
    <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
      <CardContent className="p-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Activity className="w-4 h-4 text-green-400" />
            <span className="text-lg font-bold text-white">{symbol}</span>
            <Badge 
              variant="outline" 
              className={`text-xs ${isLive ? 'text-green-400 border-green-400' : 'text-gray-400 border-gray-400'}`}
            >
              {isLive ? 'LIVE' : 'PAUSED'}
            </Badge>
          </div>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsLive(!isLive)}
            className="text-gray-400 hover:text-white"
          >
            {isLive ? <Zap className="w-4 h-4" /> : <Clock className="w-4 h-4" />}
          </Button>
        </div>

        {/* Price Display */}
        <div className="mb-4">
          <div className="flex items-center space-x-3">
            <span 
              className={`text-3xl font-bold transition-colors duration-300 ${
                priceDirection === 'up' ? 'text-green-400' : 
                priceDirection === 'down' ? 'text-red-400' : 'text-white'
              }`}
            >
              ${lastPrice}
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
                {marketData.changePercent24h}%
              </Badge>
            </div>
          </div>
          
          <div className="text-sm text-gray-400 mt-1">
            {isPositive ? '+' : ''}{marketData.change24h} USDT (24h)
          </div>
        </div>

        {/* Market Stats Grid */}
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-400">24h High</span>
              <span className="text-green-400 font-mono">${marketData.high24h}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">24h Low</span>
              <span className="text-red-400 font-mono">${marketData.low24h}</span>
            </div>
          </div>
          
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-400">24h Volume</span>
              <span className="text-white font-mono">{marketData.volume24h} BTC</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Last Update</span>
              <span className="text-gray-300 font-mono text-xs">{marketData.lastUpdate}</span>
            </div>
          </div>
        </div>

        {/* Price Movement Indicator */}
        <div className="mt-4 pt-3 border-t border-gray-700">
          <div className="flex items-center justify-between text-xs">
            <span className="text-gray-400">Price Movement</span>
            <div className="flex items-center space-x-2">
              <div className={`w-2 h-2 rounded-full ${
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
            <div className="text-green-400 font-mono">43,250.50</div>
          </div>
          <div className="text-center p-2 bg-gray-800/50 rounded">
            <div className="text-gray-400">Ask</div>
            <div className="text-red-400 font-mono">43,251.00</div>
          </div>
          <div className="text-center p-2 bg-gray-800/50 rounded">
            <div className="text-gray-400">Spread</div>
            <div className="text-white font-mono">0.50</div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
