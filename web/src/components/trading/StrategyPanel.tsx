import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Progress } from '@/components/ui/progress';
import { 
  Brain, 
  TrendingUp, 
  Settings, 
  Play, 
  Pause,
  BarChart3
} from 'lucide-react';

export const StrategyPanel: React.FC = () => {
  const strategies = [
    {
      id: 'market_making_1',
      name: 'Market Making Strategy',
      type: 'MARKET_MAKING',
      enabled: true,
      performance: { pnl: '+$1,234.56', winRate: 78.5 },
      symbols: ['BTCUSDT', 'ETHUSDT'],
      status: 'RUNNING'
    },
    {
      id: 'arbitrage_1',
      name: 'Cross-Exchange Arbitrage',
      type: 'ARBITRAGE',
      enabled: true,
      performance: { pnl: '+$892.34', winRate: 85.2 },
      symbols: ['BTCUSDT'],
      status: 'RUNNING'
    },
    {
      id: 'momentum_1',
      name: 'Momentum Strategy',
      type: 'MOMENTUM',
      enabled: false,
      performance: { pnl: '-$123.45', winRate: 45.8 },
      symbols: ['ETHUSDT', 'BNBUSDT'],
      status: 'STOPPED'
    }
  ];

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Strategies</CardTitle>
            <Brain className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">2</div>
            <p className="text-xs text-muted-foreground">
              1 stopped
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Strategy P&L</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">+$2,126.90</div>
            <p className="text-xs text-muted-foreground">
              Today: +$456.78
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Win Rate</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">69.8%</div>
            <p className="text-xs text-muted-foreground">
              Across all strategies
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Strategy Management</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {strategies.map((strategy) => (
              <div key={strategy.id} className="border rounded-lg p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-3">
                    <Switch checked={strategy.enabled} />
                    <div>
                      <h4 className="font-semibold">{strategy.name}</h4>
                      <p className="text-sm text-muted-foreground">{strategy.type}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge variant={strategy.status === 'RUNNING' ? 'default' : 'secondary'}>
                      {strategy.status}
                    </Badge>
                    <Button variant="outline" size="sm">
                      <Settings className="w-4 h-4" />
                    </Button>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">P&L</p>
                    <p className={`font-semibold ${strategy.performance.pnl.startsWith('+') ? 'text-green-600' : 'text-red-600'}`}>
                      {strategy.performance.pnl}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Win Rate</p>
                    <div className="flex items-center space-x-2">
                      <Progress value={strategy.performance.winRate} className="flex-1" />
                      <span className="text-sm font-medium">{strategy.performance.winRate}%</span>
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Symbols</p>
                    <div className="flex flex-wrap gap-1">
                      {strategy.symbols.map((symbol) => (
                        <Badge key={symbol} variant="outline" className="text-xs">
                          {symbol}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
