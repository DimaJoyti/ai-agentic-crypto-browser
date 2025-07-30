import React, { useState, useEffect, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { 
  DollarSign, 
  TrendingUp, 
  TrendingDown, 
  PieChart,
  BarChart3,
  RefreshCw,
  AlertTriangle,
  Target
} from 'lucide-react';

interface Position {
  symbol: string;
  size: string;
  avgPrice: string;
  currentPrice: string;
  unrealizedPnL: string;
  realizedPnL: string;
  commission: string;
  openTime: string;
  updateTime: string;
  exchange: string;
  strategyId: string;
}

interface PortfolioMetrics {
  totalValue: string;
  cashBalance: string;
  totalPnL: string;
  unrealizedPnL: string;
  realizedPnL: string;
  dayPnL: string;
  maxDrawdown: string;
  highWaterMark: string;
  openPositions: number;
  totalTrades: number;
}

interface RiskMetrics {
  var95: string;
  var99: string;
  expectedShortfall: string;
  sharpeRatio: number;
  sortinoRatio: number;
  maxDrawdown: string;
  beta: number;
  alpha: number;
  volatility: number;
}

interface PortfolioPanelProps {
  className?: string;
}

export const PortfolioPanel: React.FC<PortfolioPanelProps> = ({ className }) => {
  const [positions, setPositions] = useState<Position[]>([]);
  const [portfolioMetrics, setPortfolioMetrics] = useState<PortfolioMetrics | null>(null);
  const [riskMetrics, setRiskMetrics] = useState<RiskMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedView, setSelectedView] = useState<'positions' | 'metrics' | 'risk'>('positions');

  // Fetch portfolio data
  useEffect(() => {
    const fetchPortfolioData = async () => {
      try {
        setLoading(true);
        
        // Fetch positions
        const positionsResponse = await fetch('/api/portfolio/positions');
        if (positionsResponse.ok) {
          const positionsData = await positionsResponse.json();
          setPositions(positionsData);
        }
        
        // Fetch portfolio metrics
        const metricsResponse = await fetch('/api/portfolio/metrics');
        if (metricsResponse.ok) {
          const metricsData = await metricsResponse.json();
          setPortfolioMetrics(metricsData);
        }
        
        // Fetch risk metrics
        const riskResponse = await fetch('/api/portfolio/risk');
        if (riskResponse.ok) {
          const riskData = await riskResponse.json();
          setRiskMetrics(riskData);
        }
      } catch (error) {
        console.error('Failed to fetch portfolio data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchPortfolioData();
    
    // Set up auto-refresh
    const interval = setInterval(fetchPortfolioData, 5000);
    return () => clearInterval(interval);
  }, []);

  // Calculate portfolio allocation
  const portfolioAllocation = useMemo(() => {
    if (!positions.length || !portfolioMetrics) return [];
    
    const totalValue = parseFloat(portfolioMetrics.totalValue);
    
    return positions.map(position => {
      const positionValue = Math.abs(parseFloat(position.size)) * parseFloat(position.currentPrice);
      const allocation = (positionValue / totalValue) * 100;
      
      return {
        symbol: position.symbol,
        value: positionValue,
        allocation: allocation,
        pnl: parseFloat(position.unrealizedPnL),
        size: parseFloat(position.size)
      };
    }).sort((a, b) => b.allocation - a.allocation);
  }, [positions, portfolioMetrics]);

  const formatCurrency = (value: string | number) => {
    const num = typeof value === 'string' ? parseFloat(value) : value;
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(num);
  };

  const formatPercent = (value: number) => {
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`;
  };

  const formatPnL = (value: string | number) => {
    const num = typeof value === 'string' ? parseFloat(value) : value;
    const isPositive = num >= 0;
    return (
      <span className={isPositive ? 'text-green-600' : 'text-red-600'}>
        {formatCurrency(num)}
      </span>
    );
  };

  const getPositionIcon = (size: number) => {
    if (size > 0) return <TrendingUp className="w-4 h-4 text-green-600" />;
    if (size < 0) return <TrendingDown className="w-4 h-4 text-red-600" />;
    return <BarChart3 className="w-4 h-4 text-gray-400" />;
  };

  if (loading) {
    return (
      <Card className={className}>
        <CardContent className="flex items-center justify-center h-64">
          <RefreshCw className="w-8 h-8 animate-spin" />
        </CardContent>
      </Card>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Portfolio Summary */}
      {portfolioMetrics && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Value</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatCurrency(portfolioMetrics.totalValue)}</div>
              <p className="text-xs text-muted-foreground">
                Cash: {formatCurrency(portfolioMetrics.cashBalance)}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatPnL(portfolioMetrics.totalPnL)}</div>
              <p className="text-xs text-muted-foreground">
                Day P&L: {formatPnL(portfolioMetrics.dayPnL)}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Unrealized P&L</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatPnL(portfolioMetrics.unrealizedPnL)}</div>
              <p className="text-xs text-muted-foreground">
                Realized: {formatPnL(portfolioMetrics.realizedPnL)}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Positions</CardTitle>
              <PieChart className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{portfolioMetrics.openPositions}</div>
              <p className="text-xs text-muted-foreground">
                Total Trades: {portfolioMetrics.totalTrades}
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* View Selector */}
      <div className="flex space-x-2">
        <Button
          variant={selectedView === 'positions' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSelectedView('positions')}
        >
          Positions
        </Button>
        <Button
          variant={selectedView === 'metrics' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSelectedView('metrics')}
        >
          Metrics
        </Button>
        <Button
          variant={selectedView === 'risk' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSelectedView('risk')}
        >
          Risk
        </Button>
      </div>

      {/* Positions View */}
      {selectedView === 'positions' && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Positions Table */}
          <Card className="lg:col-span-2">
            <CardHeader>
              <CardTitle>Open Positions</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Symbol</TableHead>
                    <TableHead className="text-right">Size</TableHead>
                    <TableHead className="text-right">Avg Price</TableHead>
                    <TableHead className="text-right">Current Price</TableHead>
                    <TableHead className="text-right">P&L</TableHead>
                    <TableHead className="text-center">Exchange</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {positions.map((position, index) => {
                    const size = parseFloat(position.size);
                    const pnl = parseFloat(position.unrealizedPnL);
                    const pnlPercent = ((parseFloat(position.currentPrice) - parseFloat(position.avgPrice)) / parseFloat(position.avgPrice)) * 100;
                    
                    return (
                      <TableRow key={index}>
                        <TableCell className="font-medium">
                          <div className="flex items-center">
                            {getPositionIcon(size)}
                            <span className="ml-2">{position.symbol}</span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          {Math.abs(size).toFixed(4)}
                          <Badge variant={size > 0 ? 'default' : 'destructive'} className="ml-2 text-xs">
                            {size > 0 ? 'LONG' : 'SHORT'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          ${parseFloat(position.avgPrice).toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          ${parseFloat(position.currentPrice).toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right">
                          <div>
                            {formatPnL(pnl)}
                            <div className="text-xs text-muted-foreground">
                              {formatPercent(pnlPercent)}
                            </div>
                          </div>
                        </TableCell>
                        <TableCell className="text-center">
                          <Badge variant="outline" className="text-xs">
                            {position.exchange.toUpperCase()}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
              
              {positions.length === 0 && (
                <div className="text-center py-8 text-muted-foreground">
                  <Target className="w-12 h-12 mx-auto mb-4 opacity-50" />
                  <p>No open positions</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Portfolio Allocation */}
          <Card>
            <CardHeader>
              <CardTitle>Portfolio Allocation</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {portfolioAllocation.slice(0, 10).map((item, index) => (
                  <div key={item.symbol} className="space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">{item.symbol}</span>
                      <span className="text-sm text-muted-foreground">
                        {item.allocation.toFixed(1)}%
                      </span>
                    </div>
                    <Progress value={item.allocation} className="h-2" />
                    <div className="flex justify-between text-xs text-muted-foreground">
                      <span>{formatCurrency(item.value)}</span>
                      <span className={item.pnl >= 0 ? 'text-green-600' : 'text-red-600'}>
                        {formatCurrency(item.pnl)}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Risk Metrics View */}
      {selectedView === 'risk' && riskMetrics && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <AlertTriangle className="w-5 h-5 mr-2" />
              Risk Metrics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div>
                <h4 className="font-semibold mb-3">Value at Risk</h4>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">VaR 95%</span>
                    <span className="font-mono">{formatCurrency(riskMetrics.var95)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">VaR 99%</span>
                    <span className="font-mono">{formatCurrency(riskMetrics.var99)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Expected Shortfall</span>
                    <span className="font-mono">{formatCurrency(riskMetrics.expectedShortfall)}</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-3">Performance Ratios</h4>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Sharpe Ratio</span>
                    <span className="font-mono">{riskMetrics.sharpeRatio.toFixed(3)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Sortino Ratio</span>
                    <span className="font-mono">{riskMetrics.sortinoRatio.toFixed(3)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Alpha</span>
                    <span className="font-mono">{formatPercent(riskMetrics.alpha * 100)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Beta</span>
                    <span className="font-mono">{riskMetrics.beta.toFixed(3)}</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-3">Risk Measures</h4>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Max Drawdown</span>
                    <span className="font-mono text-red-600">{formatCurrency(riskMetrics.maxDrawdown)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Volatility</span>
                    <span className="font-mono">{formatPercent(riskMetrics.volatility * 100)}</span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};
