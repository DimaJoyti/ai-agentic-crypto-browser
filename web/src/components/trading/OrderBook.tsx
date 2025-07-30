import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Settings, 
  TrendingUp, 
  TrendingDown,
  Volume2,
  ArrowUpDown,
  Zap
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
}

export const OrderBook: React.FC<OrderBookProps> = ({ symbol, className }) => {
  const [precision, setPrecision] = useState(2);
  const [grouping, setGrouping] = useState(0.01);

  // Mock order book data
  const mockAsks: OrderBookEntry[] = [
    { price: '43,254.00', size: '0.5678', total: '24,550.00', percentage: 15 },
    { price: '43,253.25', size: '1.8765', total: '81,200.00', percentage: 45 },
    { price: '43,252.50', size: '0.9876', total: '42,700.00', percentage: 25 },
    { price: '43,251.75', size: '1.5432', total: '66,700.00', percentage: 35 },
    { price: '43,251.00', size: '0.4321', total: '18,700.00', percentage: 12 },
  ];

  const mockBids: OrderBookEntry[] = [
    { price: '43,250.50', size: '0.5432', total: '23,500.00', percentage: 18 },
    { price: '43,249.75', size: '1.2345', total: '53,400.00', percentage: 38 },
    { price: '43,249.00', size: '0.8765', total: '37,900.00', percentage: 28 },
    { price: '43,248.25', size: '2.1234', total: '91,800.00', percentage: 55 },
    { price: '43,247.50', size: '0.6789', total: '29,350.00', percentage: 22 },
  ];

  const spread = 0.50;
  const spreadPercent = 0.001;

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
          {mockAsks.reverse().map((ask, index) => (
            <div 
              key={index} 
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
              ${spread.toFixed(2)} ({(spreadPercent * 100).toFixed(3)}%)
            </div>
          </div>
        </div>
        
        {/* Bids (Buy Orders) */}
        <div className="space-y-1">
          {mockBids.map((bid, index) => (
            <div 
              key={index} 
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
        <div className="mt-4 pt-3 border-t border-gray-700">
          <div className="grid grid-cols-2 gap-3 text-xs">
            <div className="flex items-center justify-between">
              <span className="text-gray-400">24h Volume</span>
              <span className="text-white font-mono">1,234.56</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-400">24h High</span>
              <span className="text-green-400 font-mono">44,125.00</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-400">24h Low</span>
              <span className="text-red-400 font-mono">42,850.00</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-400">24h Change</span>
              <span className="text-green-400 font-mono">+2.34%</span>
            </div>
          </div>
        </div>

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
