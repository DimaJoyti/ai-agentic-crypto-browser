import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  TrendingUp, 
  BarChart3, 
  Target, 
  Clock,
  DollarSign,
  Percent
} from 'lucide-react';

export const PerformancePanel: React.FC = () => {
  const performanceMetrics = {
    totalPnL: 15420.50,
    dailyPnL: 1234.56,
    winRate: 68.5,
    sharpeRatio: 1.85,
    maxDrawdown: 2150.30,
    totalTrades: 1247,
    avgTradeSize: 0.05,
    avgHoldTime: '2.3 minutes'
  };

  const strategyPerformance = [
    { name: 'Market Making', pnl: 8420.30, trades: 856, winRate: 72.1 },
    { name: 'Arbitrage', pnl: 5200.15, trades: 234, winRate: 89.2 },
    { name: 'Momentum', pnl: 1800.05, trades: 157, winRate: 45.8 }
  ];

  const timeframes = [
    { period: '1H', pnl: 156.78, change: '+2.1%' },
    { period: '1D', pnl: 1234.56, change: '+8.7%' },
    { period: '1W', pnl: 4567.89, change: '+12.3%' },
    { period: '1M', pnl: 15420.50, change: '+18.9%' }
  ];

  return (
    <div className="space-y-6">
      {/* Performance Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              +${performanceMetrics.totalPnL.toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">
              Today: +${performanceMetrics.dailyPnL.toFixed(2)}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Win Rate</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{performanceMetrics.winRate}%</div>
            <Progress value={performanceMetrics.winRate} className="mt-2" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Sharpe Ratio</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{performanceMetrics.sharpeRatio}</div>
            <p className="text-xs text-muted-foreground">
              Excellent performance
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Trades</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{performanceMetrics.totalTrades.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              Avg hold: {performanceMetrics.avgHoldTime}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Performance by Timeframe */}
      <Card>
        <CardHeader>
          <CardTitle>Performance by Timeframe</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {timeframes.map((tf) => (
              <div key={tf.period} className="text-center p-4 border rounded-lg">
                <div className="text-sm text-muted-foreground mb-1">{tf.period}</div>
                <div className="text-xl font-bold text-green-600">
                  +${tf.pnl.toLocaleString()}
                </div>
                <div className="text-sm text-green-600">{tf.change}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Strategy Performance */}
      <Card>
        <CardHeader>
          <CardTitle>Strategy Performance</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {strategyPerformance.map((strategy, index) => (
              <div key={index} className="border rounded-lg p-4">
                <div className="flex items-center justify-between mb-3">
                  <h4 className="font-semibold">{strategy.name}</h4>
                  <Badge variant="outline">{strategy.trades} trades</Badge>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">P&L</p>
                    <p className="text-lg font-bold text-green-600">
                      +${strategy.pnl.toLocaleString()}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Win Rate</p>
                    <div className="flex items-center space-x-2">
                      <Progress value={strategy.winRate} className="flex-1" />
                      <span className="text-sm font-medium">{strategy.winRate}%</span>
                    </div>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Avg P&L per Trade</p>
                    <p className="text-lg font-bold">
                      ${(strategy.pnl / strategy.trades).toFixed(2)}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Risk Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Risk Metrics</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Max Drawdown</span>
                <span className="font-semibold text-red-600">
                  -${performanceMetrics.maxDrawdown.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Sharpe Ratio</span>
                <span className="font-semibold">{performanceMetrics.sharpeRatio}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Sortino Ratio</span>
                <span className="font-semibold">2.14</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Calmar Ratio</span>
                <span className="font-semibold">7.17</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Trading Statistics</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Avg Trade Size</span>
                <span className="font-semibold">{performanceMetrics.avgTradeSize} BTC</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Avg Hold Time</span>
                <span className="font-semibold">{performanceMetrics.avgHoldTime}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Largest Win</span>
                <span className="font-semibold text-green-600">+$456.78</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Largest Loss</span>
                <span className="font-semibold text-red-600">-$123.45</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
