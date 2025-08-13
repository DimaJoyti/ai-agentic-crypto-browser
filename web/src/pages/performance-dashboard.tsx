import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { 
  TrendingUp, 
  TrendingDown,
  DollarSign, 
  Target, 
  Award,
  BarChart3,
  PieChart,
  Activity,
  Trophy,
  Zap,
  Shield
} from 'lucide-react';

interface PerformanceSummary {
  total_trades: number;
  profitable_trades: number;
  total_pnl: number;
  total_performance_fees: number;
  win_rate: number;
  average_return: number;
  max_drawdown: number;
  sharpe_ratio: number;
  current_high_water_mark: number;
}

interface TradeRecord {
  id: string;
  symbol: string;
  side: string;
  quantity: number;
  entry_price: number;
  exit_price: number;
  pnl: number;
  performance_fee: number;
  entry_timestamp: string;
  exit_timestamp: string;
  status: string;
}

const PerformanceDashboard = () => {
  const [summary, setSummary] = useState<PerformanceSummary | null>(null);
  const [trades, setTrades] = useState<TradeRecord[]>([]);
  const [analytics, setAnalytics] = useState<any>(null);
  const [billing, setBilling] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchPerformanceData();
  }, []);

  const fetchPerformanceData = async () => {
    try {
      const [summaryResponse, tradesResponse, analyticsResponse, billingResponse] = await Promise.all([
        fetch('/api/performance/summary'),
        fetch('/api/performance/trades?limit=10'),
        fetch('/api/performance/analytics'),
        fetch('/api/performance/fees/billing')
      ]);

      const summaryData = await summaryResponse.json();
      const tradesData = await tradesResponse.json();
      const analyticsData = await analyticsResponse.json();
      const billingData = await billingResponse.json();

      setSummary(summaryData);
      setTrades(tradesData.trades || []);
      setAnalytics(analyticsData);
      setBilling(billingData);
    } catch (error) {
      console.error('Error fetching performance data:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const formatPercentage = (rate: number) => {
    return `${(rate * 100).toFixed(1)}%`;
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/3"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Performance Dashboard</h1>
          <p className="text-gray-600">Track your trading performance and fees</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline">
            <BarChart3 className="w-4 h-4 mr-2" />
            Export Report
          </Button>
          <Button>
            <Target className="w-4 h-4 mr-2" />
            Set Goals
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      {summary && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(summary.total_pnl)}
              </div>
              <p className="text-xs text-gray-600">
                +{formatPercentage(summary.average_return)} average return
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Performance Fees</CardTitle>
              <DollarSign className="h-4 w-4 text-blue-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">
                {formatCurrency(summary.total_performance_fees)}
              </div>
              <p className="text-xs text-gray-600">
                {summary.total_pnl > 0
                  ? `${((summary.total_performance_fees / summary.total_pnl) * 100).toFixed(1)}% of profits`
                  : 'No profits yet'
                }
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Win Rate</CardTitle>
              <Target className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-purple-600">
                {formatPercentage(summary.win_rate)}
              </div>
              <p className="text-xs text-gray-600">
                {summary.profitable_trades} of {summary.total_trades} trades
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Sharpe Ratio</CardTitle>
              <Award className="h-4 w-4 text-orange-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-600">
                {summary.sharpe_ratio.toFixed(2)}
              </div>
              <p className="text-xs text-gray-600">
                Risk-adjusted returns
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs defaultValue="overview" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="trades">Trades</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="billing">Billing</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* High Water Mark */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Shield className="w-5 h-5 mr-2 text-blue-600" />
                  High Water Mark
                </CardTitle>
                <CardDescription>
                  Performance fees only charged above this level
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-3xl font-bold text-blue-600">
                    {summary && formatCurrency(summary.current_high_water_mark)}
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Current Portfolio</span>
                      <span>{summary && formatCurrency(summary.current_high_water_mark)}</span>
                    </div>
                    <Progress value={100} className="h-2" />
                    <p className="text-xs text-gray-600">
                      At all-time high • No recovery needed
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Performance Metrics */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Activity className="w-5 h-5 mr-2 text-green-600" />
                  Performance Metrics
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex justify-between items-center">
                    <span className="text-sm">Max Drawdown</span>
                    <span className="font-medium text-red-600">
                      {summary && formatPercentage(summary.max_drawdown)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">Sharpe Ratio</span>
                    <span className="font-medium text-green-600">
                      {summary && summary.sharpe_ratio.toFixed(2)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">Total Trades</span>
                    <span className="font-medium">
                      {summary && summary.total_trades}
                    </span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">Average Return</span>
                    <span className="font-medium text-green-600">
                      {summary && formatPercentage(summary.average_return)}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Recent Trades */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Trades</CardTitle>
              <CardDescription>Your latest trading activity</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {trades.slice(0, 5).map((trade) => (
                  <div key={trade.id} className="flex justify-between items-center p-4 border rounded-lg">
                    <div className="flex items-center space-x-4">
                      <div className={`p-2 rounded-full ${
                        trade.pnl > 0 ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
                      }`}>
                        {trade.pnl > 0 ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
                      </div>
                      <div>
                        <div className="font-medium">{trade.symbol}</div>
                        <div className="text-sm text-gray-600">
                          {trade.side.toUpperCase()} • {trade.quantity} units
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className={`font-bold ${
                        trade.pnl > 0 ? 'text-green-600' : 'text-red-600'
                      }`}>
                        {formatCurrency(trade.pnl)}
                      </div>
                      <div className="text-sm text-gray-600">
                        Fee: {formatCurrency(trade.performance_fee)}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Trades Tab */}
        <TabsContent value="trades" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Trade History</CardTitle>
              <CardDescription>Complete record of your trades and performance fees</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {trades.map((trade) => (
                  <div key={trade.id} className="grid grid-cols-6 gap-4 p-4 border rounded-lg">
                    <div>
                      <div className="font-medium">{trade.symbol}</div>
                      <div className="text-sm text-gray-600">{trade.side.toUpperCase()}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Entry</div>
                      <div className="font-medium">{formatCurrency(trade.entry_price)}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Exit</div>
                      <div className="font-medium">{formatCurrency(trade.exit_price)}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">P&L</div>
                      <div className={`font-medium ${
                        trade.pnl > 0 ? 'text-green-600' : 'text-red-600'
                      }`}>
                        {formatCurrency(trade.pnl)}
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Performance Fee</div>
                      <div className="font-medium">{formatCurrency(trade.performance_fee)}</div>
                    </div>
                    <div>
                      <Badge variant={trade.status === 'completed' ? 'default' : 'secondary'}>
                        {trade.status}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Analytics Tab */}
        <TabsContent value="analytics" className="space-y-6">
          {analytics && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Performance Metrics</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {analytics.performance_metrics && Object.entries(analytics.performance_metrics).map(([key, value]) => (
                      <div key={key} className="flex justify-between">
                        <span className="text-sm capitalize">{key.replace(/_/g, ' ')}</span>
                        <span className="font-medium">{String(value)}</span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Fee Breakdown</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {analytics.fee_breakdown && Object.entries(analytics.fee_breakdown).map(([key, value]) => (
                      <div key={key} className="flex justify-between">
                        <span className="text-sm capitalize">{key.replace(/_/g, ' ')}</span>
                        <span className="font-medium">{String(value)}</span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        {/* Billing Tab */}
        <TabsContent value="billing" className="space-y-6">
          {billing && (
            <>
              <Card>
                <CardHeader>
                  <CardTitle>Current Month Billing</CardTitle>
                  <CardDescription>Performance fees for {billing.current_month.period}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div>
                      <div className="text-sm text-gray-600">Trades</div>
                      <div className="text-2xl font-bold">{billing.current_month.trades}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">P&L</div>
                      <div className="text-2xl font-bold text-green-600">
                        {formatCurrency(billing.current_month.total_pnl)}
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Performance Fees</div>
                      <div className="text-2xl font-bold text-blue-600">
                        {formatCurrency(billing.current_month.performance_fees)}
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-600">Status</div>
                      <Badge variant="outline" className="mt-1">
                        {billing.current_month.status}
                      </Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Billing History</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {billing.billing_history.map((bill: any, index: number) => (
                      <div key={index} className="flex justify-between items-center p-4 border rounded-lg">
                        <div>
                          <div className="font-medium">{bill.period}</div>
                          <div className="text-sm text-gray-600">
                            {bill.trades} trades • {bill.profitable_trades} profitable
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold">{formatCurrency(bill.performance_fees)}</div>
                          <Badge variant={bill.status === 'paid' ? 'default' : 'secondary'}>
                            {bill.status}
                          </Badge>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default PerformanceDashboard;
