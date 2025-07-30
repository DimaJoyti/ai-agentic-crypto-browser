import React, { useEffect, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { 
  BarChart3, 
  Maximize2, 
  Settings, 
  TrendingUp,
  RefreshCw
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
  const chartContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // In a real implementation, this would initialize TradingView widget
    // For now, we'll create a professional-looking mock chart
    if (chartContainerRef.current) {
      // Clear previous content
      chartContainerRef.current.innerHTML = '';

      // Create professional mock chart
      const mockChart = document.createElement('div');
      mockChart.style.width = '100%';
      mockChart.style.height = '100%';
      mockChart.style.background = 'linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%)';
      mockChart.style.borderRadius = '8px';
      mockChart.style.position = 'relative';
      mockChart.style.overflow = 'hidden';

      // Create grid pattern
      const gridPattern = document.createElement('div');
      gridPattern.style.position = 'absolute';
      gridPattern.style.top = '0';
      gridPattern.style.left = '0';
      gridPattern.style.width = '100%';
      gridPattern.style.height = '100%';
      gridPattern.style.backgroundImage = `
        linear-gradient(rgba(255, 255, 255, 0.1) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, 0.1) 1px, transparent 1px)
      `;
      gridPattern.style.backgroundSize = '20px 20px';
      gridPattern.style.opacity = '0.3';

      // Create price line simulation
      const priceLine = document.createElement('div');
      priceLine.style.position = 'absolute';
      priceLine.style.top = '40%';
      priceLine.style.left = '10%';
      priceLine.style.width = '80%';
      priceLine.style.height = '2px';
      priceLine.style.background = 'linear-gradient(90deg, #10b981 0%, #3b82f6 50%, #ef4444 100%)';
      priceLine.style.borderRadius = '1px';

      // Create candlestick simulation
      const candlesticks = document.createElement('div');
      candlesticks.style.position = 'absolute';
      candlesticks.style.bottom = '20%';
      candlesticks.style.left = '10%';
      candlesticks.style.width = '80%';
      candlesticks.style.height = '60%';
      candlesticks.style.display = 'flex';
      candlesticks.style.alignItems = 'end';
      candlesticks.style.justifyContent = 'space-between';

      // Add mock candlesticks
      for (let i = 0; i < 50; i++) {
        const candle = document.createElement('div');
        const height = Math.random() * 80 + 20;
        const isGreen = Math.random() > 0.5;
        candle.style.width = '2px';
        candle.style.height = `${height}%`;
        candle.style.backgroundColor = isGreen ? '#10b981' : '#ef4444';
        candle.style.opacity = '0.8';
        candlesticks.appendChild(candle);
      }

      // Create overlay info
      const overlay = document.createElement('div');
      overlay.style.position = 'absolute';
      overlay.style.top = '20px';
      overlay.style.left = '20px';
      overlay.style.color = 'white';
      overlay.style.fontSize = '14px';
      overlay.style.fontFamily = 'monospace';
      overlay.innerHTML = `
        <div style="margin-bottom: 10px; font-size: 18px; font-weight: bold; color: #3b82f6;">${symbol}</div>
        <div style="margin-bottom: 5px;">O: <span style="color: #10b981;">$43,250.50</span></div>
        <div style="margin-bottom: 5px;">H: <span style="color: #10b981;">$44,125.00</span></div>
        <div style="margin-bottom: 5px;">L: <span style="color: #ef4444;">$42,850.00</span></div>
        <div style="margin-bottom: 5px;">C: <span style="color: #10b981;">$43,251.50</span></div>
        <div style="color: #6b7280; font-size: 12px;">Volume: 1,234.56 BTC</div>
        </div>
      `;

      // Assemble the chart
      mockChart.appendChild(gridPattern);
      mockChart.appendChild(candlesticks);
      mockChart.appendChild(priceLine);
      mockChart.appendChild(overlay);

      chartContainerRef.current.appendChild(mockChart);
    }
  }, [symbol]);

  const symbols = [
    'BTCUSDT', 'ETHUSDT', 'BNBUSDT', 'ADAUSDT', 
    'XRPUSDT', 'SOLUSDT', 'DOTUSDT', 'LINKUSDT'
  ];

  const timeframes = ['1m', '5m', '15m', '1h', '4h', '1d'];

  return (
    <Card className={className}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center">
            <BarChart3 className="w-5 h-5 mr-2" />
            Price Chart
          </CardTitle>
          <div className="flex items-center space-x-2">
            <Badge variant="outline" className="flex items-center">
              <TrendingUp className="w-3 h-3 mr-1" />
              Live
            </Badge>
            <Button variant="outline" size="sm">
              <Maximize2 className="w-4 h-4" />
            </Button>
            <Button variant="outline" size="sm">
              <Settings className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {/* Chart Controls */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Select value={symbol} onValueChange={onSymbolChange}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {symbols.map((sym) => (
                  <SelectItem key={sym} value={sym}>
                    {sym}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            
            <div className="flex space-x-1">
              {timeframes.map((tf) => (
                <Button
                  key={tf}
                  variant="outline"
                  size="sm"
                  className="text-xs px-2 py-1"
                >
                  {tf}
                </Button>
              ))}
            </div>
          </div>
          
          <Button variant="outline" size="sm">
            <RefreshCw className="w-4 h-4" />
          </Button>
        </div>

        {/* Chart Container */}
        <div ref={chartContainerRef} className="w-full" />

        {/* Chart Info */}
        <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <span className="text-muted-foreground">Open:</span>
            <span className="ml-2 font-mono">$45,123.45</span>
          </div>
          <div>
            <span className="text-muted-foreground">High:</span>
            <span className="ml-2 font-mono text-green-600">$45,678.90</span>
          </div>
          <div>
            <span className="text-muted-foreground">Low:</span>
            <span className="ml-2 font-mono text-red-600">$44,890.12</span>
          </div>
          <div>
            <span className="text-muted-foreground">Volume:</span>
            <span className="ml-2 font-mono">1.23M</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
