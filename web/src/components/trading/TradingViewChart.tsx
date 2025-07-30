import React, { useState } from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import TradingViewWidget from './TradingViewWidget';
import EmbeddedTradingViewChart from './EmbeddedTradingViewChart';
import {
  BarChart3,
  Maximize2,
  Settings,
  TrendingUp,
  RefreshCw,
  Fullscreen,
  RotateCcw
} from 'lucide-react';

interface TradingViewChartProps {
  symbol: string;
  onSymbolChange: (symbol: string) => void;
  className?: string;
}

export const TradingViewChart: React.FC<TradingViewChartProps> = ({
  symbol,
  onSymbolChange,
  className
}) => {
  const [theme, setTheme] = useState<'light' | 'dark'>('dark');
  const [interval, setInterval] = useState('1D');
  const [useEmbedded, setUseEmbedded] = useState(false);

  const handleIntervalChange = (newInterval: string) => {
    setInterval(newInterval);
  };

  const handleThemeToggle = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
  };

  const handleFullscreen = () => {
    const element = document.documentElement;
    if (element.requestFullscreen) {
      element.requestFullscreen();
    }
  };

  const toggleChartType = () => {
    setUseEmbedded(!useEmbedded);
  };

  const symbols = [
    'BTCUSDT', 'ETHUSDT', 'BNBUSDT', 'ADAUSDT',
    'XRPUSDT', 'SOLUSDT', 'DOTUSDT', 'LINKUSDT',
    'AVAXUSDT', 'MATICUSDT', 'ATOMUSDT', 'NEARUSDT'
  ];

  const timeframes = [
    { label: '1m', value: '1' },
    { label: '5m', value: '5' },
    { label: '15m', value: '15' },
    { label: '1h', value: '60' },
    { label: '4h', value: '240' },
    { label: '1d', value: '1D' }
  ];

  return (
    <div className={`w-full h-full ${className}`}>
      {/* Chart Controls */}
      <div className="flex items-center justify-between mb-4 p-2 bg-gray-800/30 rounded-lg">
        <div className="flex items-center space-x-3">
          <Select value={symbol} onValueChange={onSymbolChange}>
            <SelectTrigger className="w-32 bg-gray-800 border-gray-700 text-white">
              <SelectValue />
            </SelectTrigger>
            <SelectContent className="bg-gray-800 border-gray-700">
              {symbols.map((sym) => (
                <SelectItem key={sym} value={sym} className="text-white hover:bg-gray-700">
                  {sym}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <div className="flex space-x-1">
            {timeframes.map((tf) => (
              <Button
                key={tf.value}
                variant={interval === tf.value ? "default" : "outline"}
                size="sm"
                className={`text-xs px-3 py-1 ${
                  interval === tf.value
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-800 border-gray-700 text-gray-300 hover:bg-gray-700'
                }`}
                onClick={() => handleIntervalChange(tf.value)}
              >
                {tf.label}
              </Button>
            ))}
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <Badge variant="outline" className="text-green-400 border-green-400">
            <TrendingUp className="w-3 h-3 mr-1" />
            Live
          </Badge>

          <Button
            variant="outline"
            size="sm"
            className="bg-gray-800 border-gray-700 text-gray-300 hover:bg-gray-700"
            onClick={handleThemeToggle}
          >
            {theme === 'dark' ? '🌙' : '☀️'}
          </Button>

          <Button
            variant="outline"
            size="sm"
            className="bg-gray-800 border-gray-700 text-gray-300 hover:bg-gray-700"
            onClick={handleFullscreen}
          >
            <Maximize2 className="w-4 h-4" />
          </Button>

          <Button
            variant="outline"
            size="sm"
            className="bg-gray-800 border-gray-700 text-gray-300 hover:bg-gray-700"
          >
            <Settings className="w-4 h-4" />
          </Button>
        </div>
      </div>

      {/* TradingView Chart Container */}
      <div className="relative w-full h-full bg-gray-900 rounded-lg overflow-hidden">
        {useEmbedded ? (
          <EmbeddedTradingViewChart
            symbol={symbol}
            theme={theme}
            interval={interval}
            height={500}
            className="w-full h-full"
          />
        ) : (
          <TradingViewWidget
            symbol={`BINANCE:${symbol}`}
            theme={theme}
            interval={interval}
            autosize={true}
            className="w-full h-full"
          />
        )}

        {/* Chart Type Toggle */}
        <Button
          variant="ghost"
          size="sm"
          className="absolute top-2 right-2 bg-gray-800/80 text-gray-300 hover:bg-gray-700"
          onClick={toggleChartType}
        >
          {useEmbedded ? 'Widget' : 'Embedded'}
        </Button>
      </div>
    </div>
  );
};
