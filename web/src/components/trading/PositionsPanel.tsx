import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  TrendingUp, 
  TrendingDown,
  X,
  Settings,
  Target,
  Shield,
  AlertTriangle,
  DollarSign,
  Percent,
  Activity
} from 'lucide-react';

interface Position {
  id: string;
  symbol: string;
  side: 'LONG' | 'SHORT';
  size: string;
  entryPrice: string;
  markPrice: string;
  pnl: string;
  pnlPercent: string;
  margin: string;
  leverage: string;
  liquidationPrice: string;
  marginRatio: number;
  timestamp: string;
}

interface PositionsPanelProps {
  className?: string;
}

export const PositionsPanel: React.FC<PositionsPanelProps> = ({ className }) => {
  const [selectedPosition, setSelectedPosition] = useState<string | null>(null);

  // Mock positions data
  const positions: Position[] = [
    {
      id: '1',
      symbol: 'BTCUSDT',
      side: 'LONG',
      size: '0.5432',
      entryPrice: '42,850.00',
      markPrice: '43,250.50',
      pnl: '+217.54',
      pnlPercent: '+0.51',
      margin: '4,325.00',
      leverage: '10x',
      liquidationPrice: '38,565.00',
      marginRatio: 15.2,
      timestamp: '2024-01-15 14:30:25'
    },
    {
      id: '2',
      symbol: 'ETHUSDT',
      side: 'SHORT',
      size: '2.1234',
      entryPrice: '2,650.00',
      markPrice: '2,635.75',
      pnl: '+30.23',
      pnlPercent: '+0.54',
      margin: '1,127.50',
      leverage: '5x',
      liquidationPrice: '2,915.50',
      marginRatio: 8.7,
      timestamp: '2024-01-15 13:45:12'
    },
    {
      id: '3',
      symbol: 'ADAUSDT',
      side: 'LONG',
      size: '1000.00',
      entryPrice: '0.4850',
      markPrice: '0.4823',
      pnl: '-27.00',
      pnlPercent: '-0.56',
      margin: '97.00',
      leverage: '5x',
      liquidationPrice: '0.4365',
      marginRatio: 22.1,
      timestamp: '2024-01-15 12:15:08'
    }
  ];

  const totalPnL = positions.reduce((sum, pos) => {
    return sum + parseFloat(pos.pnl.replace(/[+$,]/g, ''));
  }, 0);

  const totalMargin = positions.reduce((sum, pos) => {
    return sum + parseFloat(pos.margin.replace(/[,$]/g, ''));
  }, 0);

  const handleClosePosition = (positionId: string) => {
    console.log('Closing position:', positionId);
    // Implementation for closing position
  };

  const getMarginRiskColor = (ratio: number) => {
    if (ratio < 10) return 'text-green-400';
    if (ratio < 20) return 'text-yellow-400';
    return 'text-red-400';
  };

  const getMarginRiskBg = (ratio: number) => {
    if (ratio < 10) return 'bg-green-500';
    if (ratio < 20) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  return (
    <Card className={`bg-gray-900/50 border-gray-800 ${className}`}>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-white flex items-center">
            <Target className="w-5 h-5 mr-2 text-blue-400" />
            Open Positions ({positions.length})
          </CardTitle>
          <div className="flex items-center space-x-2">
            <Badge 
              variant="outline" 
              className={`text-xs ${totalPnL >= 0 ? 'text-green-400 border-green-400' : 'text-red-400 border-red-400'}`}
            >
              Total P&L: {totalPnL >= 0 ? '+' : ''}${totalPnL.toFixed(2)}
            </Badge>
            <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white">
              <Settings className="w-3 h-3" />
            </Button>
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="space-y-3">
        {/* Summary Stats */}
        <div className="grid grid-cols-3 gap-3 mb-4">
          <div className="text-center p-3 bg-gray-800/50 rounded-lg">
            <div className="text-xs text-gray-400 mb-1">Total Margin</div>
            <div className="text-sm font-bold text-white">${totalMargin.toLocaleString()}</div>
          </div>
          <div className="text-center p-3 bg-gray-800/50 rounded-lg">
            <div className="text-xs text-gray-400 mb-1">Unrealized P&L</div>
            <div className={`text-sm font-bold ${totalPnL >= 0 ? 'text-green-400' : 'text-red-400'}`}>
              {totalPnL >= 0 ? '+' : ''}${totalPnL.toFixed(2)}
            </div>
          </div>
          <div className="text-center p-3 bg-gray-800/50 rounded-lg">
            <div className="text-xs text-gray-400 mb-1">Avg Risk</div>
            <div className="text-sm font-bold text-yellow-400">
              {(positions.reduce((sum, pos) => sum + pos.marginRatio, 0) / positions.length).toFixed(1)}%
            </div>
          </div>
        </div>

        {/* Positions List */}
        <div className="space-y-2">
          {positions.map((position) => (
            <div 
              key={position.id}
              className={`p-3 bg-gray-800/30 rounded-lg border transition-all duration-200 cursor-pointer ${
                selectedPosition === position.id 
                  ? 'border-blue-500 bg-blue-900/20' 
                  : 'border-gray-700 hover:border-gray-600'
              }`}
              onClick={() => setSelectedPosition(selectedPosition === position.id ? null : position.id)}
            >
              {/* Position Header */}
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center space-x-3">
                  <Badge 
                    variant={position.side === 'LONG' ? 'default' : 'destructive'} 
                    className="text-xs"
                  >
                    {position.side}
                  </Badge>
                  <span className="font-bold text-white">{position.symbol}</span>
                  <span className="text-sm text-gray-400">{position.leverage}</span>
                </div>
                
                <div className="flex items-center space-x-2">
                  <span className={`text-sm font-bold ${
                    position.pnl.startsWith('+') ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {position.pnl}
                  </span>
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    className="text-red-400 hover:text-red-300 hover:bg-red-900/20"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleClosePosition(position.id);
                    }}
                  >
                    <X className="w-3 h-3" />
                  </Button>
                </div>
              </div>

              {/* Position Details */}
              <div className="grid grid-cols-4 gap-3 text-xs">
                <div>
                  <div className="text-gray-400">Size</div>
                  <div className="text-white font-mono">{position.size}</div>
                </div>
                <div>
                  <div className="text-gray-400">Entry</div>
                  <div className="text-white font-mono">${position.entryPrice}</div>
                </div>
                <div>
                  <div className="text-gray-400">Mark</div>
                  <div className="text-white font-mono">${position.markPrice}</div>
                </div>
                <div>
                  <div className="text-gray-400">P&L %</div>
                  <div className={`font-mono ${
                    position.pnlPercent.startsWith('+') ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {position.pnlPercent}%
                  </div>
                </div>
              </div>

              {/* Risk Indicator */}
              <div className="mt-2">
                <div className="flex items-center justify-between text-xs mb-1">
                  <span className="text-gray-400">Margin Risk</span>
                  <span className={getMarginRiskColor(position.marginRatio)}>
                    {position.marginRatio}%
                  </span>
                </div>
                <Progress 
                  value={position.marginRatio} 
                  className="h-1"
                  // @ts-ignore
                  indicatorClassName={getMarginRiskBg(position.marginRatio)}
                />
              </div>

              {/* Expanded Details */}
              {selectedPosition === position.id && (
                <div className="mt-3 pt-3 border-t border-gray-700 space-y-2">
                  <div className="grid grid-cols-2 gap-3 text-xs">
                    <div className="flex justify-between">
                      <span className="text-gray-400">Margin:</span>
                      <span className="text-white font-mono">${position.margin}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">Liq. Price:</span>
                      <span className="text-red-400 font-mono">${position.liquidationPrice}</span>
                    </div>
                  </div>
                  <div className="text-xs text-gray-500">
                    Opened: {position.timestamp}
                  </div>
                  
                  {/* Quick Actions */}
                  <div className="flex space-x-2 mt-2">
                    <Button size="sm" variant="outline" className="flex-1 text-xs">
                      <Shield className="w-3 h-3 mr-1" />
                      Add Margin
                    </Button>
                    <Button size="sm" variant="outline" className="flex-1 text-xs">
                      <Target className="w-3 h-3 mr-1" />
                      Set TP/SL
                    </Button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        {positions.length === 0 && (
          <div className="text-center py-8 text-gray-400">
            <Activity className="w-8 h-8 mx-auto mb-2 opacity-50" />
            <p>No open positions</p>
            <p className="text-xs mt-1">Your active positions will appear here</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
