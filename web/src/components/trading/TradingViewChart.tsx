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
    // For now, we'll create a placeholder chart
    if (chartContainerRef.current) {
      // Clear previous content
      chartContainerRef.current.innerHTML = '';
      
      // Create mock chart
      const mockChart = document.createElement('div');
      mockChart.style.width = '100%';
      mockChart.style.height = '400px';
      mockChart.style.background = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)';
      mockChart.style.borderRadius = '8px';
      mockChart.style.display = 'flex';
      mockChart.style.alignItems = 'center';
      mockChart.style.justifyContent = 'center';
      mockChart.style.color = 'white';
      mockChart.style.fontSize = '18px';
      mockChart.style.fontWeight = 'bold';
      mockChart.innerHTML = `
        <div style="text-align: center;">
          <div style="font-size: 24px; margin-bottom: 10px;">${symbol}</div>
          <div style="font-size: 16px; opacity: 0.8;">TradingView Chart</div>
          <div style="font-size: 14px; opacity: 0.6; margin-top: 10px;">Real-time price data and technical analysis</div>
        </div>
      `;
      
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
