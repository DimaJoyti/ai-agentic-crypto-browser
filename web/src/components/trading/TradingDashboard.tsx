import React, { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { 
  Activity, 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  BarChart3,
  Settings,
  Play,
  Pause,
  AlertTriangle,
  CheckCircle
} from 'lucide-react';

import { MarketDataPanel } from './MarketDataPanel';
import { OrderManagementPanel } from './OrderManagementPanel';
import { PortfolioPanel } from './PortfolioPanel';
import { StrategyPanel } from './StrategyPanel';
import { RiskManagementPanel } from './RiskManagementPanel';
import { PerformancePanel } from './PerformancePanel';
import { TradingViewChart } from './TradingViewChart';
import { useTradingDashboard } from '@/hooks/useTradingDashboard';

interface TradingDashboardProps {
  className?: string;
}

interface HFTMetrics {
  ordersPerSecond: number;
  avgLatencyMicros: number;
  totalOrders: number;
  successfulOrders: number;
  failedOrders: number;
  uptime: string;
  isRunning: boolean;
}

interface MarketTick {
  symbol: string;
  price: string;
  volume: string;
  bidPrice: string;
  askPrice: string;
  bidSize: string;
  askSize: string;
  timestamp: string;
  exchange: string;
  sequence: number;
}

interface TradingSignal {
  id: string;
  symbol: string;
  side: 'BUY' | 'SELL';
  orderType: string;
  quantity: string;
  price: string;
  confidence: number;
  strategyId: string;
  timestamp: string;
  metadata: Record<string, any>;
}

export const TradingDashboard: React.FC<TradingDashboardProps> = ({ className }) => {
  const [selectedSymbol, setSelectedSymbol] = useState('BTCUSDT');
  const [autoRefresh, setAutoRefresh] = useState(true);

  // Use the trading dashboard hook for all data and actions
  const {
    hftMetrics,
    hftStatus,
    isHFTRunning,
    startHFTEngine,
    stopHFTEngine,
    marketData,
    orders,
    positions,
    signals,
    portfolioSummary,
    portfolioMetrics,
    portfolioRisk,
    strategies,
    riskLimits,
    riskViolations,
    riskMetrics,
    systemStatus,
    isConnected,
    isLoading,
    error,
    refreshData,
    emergencyStop,
  } = useTradingDashboard();

  const [activeTab, setActiveTab] = useState('overview');

  const formatLatency = (micros: number) => {
    if (micros < 1000) return `${micros}μs`;
    if (micros < 1000000) return `${(micros / 1000).toFixed(1)}ms`;
    return `${(micros / 1000000).toFixed(2)}s`;
  };

  const formatPnL = (pnl: string) => {
    const value = parseFloat(pnl);
    const isPositive = value >= 0;
    return (
      <span className={isPositive ? 'text-green-600' : 'text-red-600'}>
        {isPositive ? '+' : ''}${value.toFixed(2)}
      </span>
    );
  };

  // Show loading state
  if (isLoading) {
    return (
      <div className={`space-y-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading trading dashboard...</p>
          </div>
        </div>
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div className={`space-y-6 ${className}`}>
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            Error loading dashboard: {error}
            <Button
              variant="outline"
              size="sm"
              className="ml-2"
              onClick={refreshData}
            >
              Retry
            </Button>
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">High-Frequency Trading Dashboard</h1>
          <p className="text-muted-foreground">Real-time trading operations and analytics</p>
        </div>
        
        <div className="flex items-center space-x-4">
          <Badge variant={isConnected ? "default" : "destructive"}>
            {isConnected ? (
              <>
                <CheckCircle className="w-4 h-4 mr-1" />
                Connected
              </>
            ) : (
              <>
                <AlertTriangle className="w-4 h-4 mr-1" />
                Disconnected
              </>
            )}
          </Badge>
          
          <Button
            variant="outline"
            size="sm"
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            {autoRefresh ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
            {autoRefresh ? 'Pause' : 'Resume'}
          </Button>
          
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4" />
            Settings
          </Button>
        </div>
      </div>

      {/* Status Alert */}
      {!isConnected && (
        <Alert>
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            Connection to trading system lost. Some features may not work properly.
          </AlertDescription>
        </Alert>
      )}

      {/* Key Metrics */}
      {hftMetrics && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Orders/Second</CardTitle>
              <Activity className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{hftMetrics.ordersPerSecond}</div>
              <p className="text-xs text-muted-foreground">
                {hftMetrics.isRunning ? 'Engine Running' : 'Engine Stopped'}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatLatency(hftMetrics.avgLatencyMicros)}</div>
              <p className="text-xs text-muted-foreground">
                Success Rate: {hftMetrics.successfulOrders > 0 ? ((hftMetrics.successfulOrders / hftMetrics.totalOrders) * 100).toFixed(1) : 0}%
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{portfolioSummary ? formatPnL(portfolioSummary.total_pnl || "0") : "$0.00"}</div>
              <p className="text-xs text-muted-foreground">
                Open Positions: {positions?.length || 0}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Strategies</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{strategies?.filter(s => s.enabled).length || 0}</div>
              <div className="flex space-x-2 mt-2">
                <Button
                  size="sm"
                  variant={isHFTRunning ? "destructive" : "default"}
                  onClick={isHFTRunning ? stopHFTEngine : startHFTEngine}
                >
                  {isHFTRunning ? (
                    <>
                      <Pause className="w-3 h-3 mr-1" />
                      Stop
                    </>
                  ) : (
                    <>
                      <Play className="w-3 h-3 mr-1" />
                      Start
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList className="grid w-full grid-cols-7">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="market">Market Data</TabsTrigger>
          <TabsTrigger value="orders">Orders</TabsTrigger>
          <TabsTrigger value="portfolio">Portfolio</TabsTrigger>
          <TabsTrigger value="strategies">Strategies</TabsTrigger>
          <TabsTrigger value="risk">Risk</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <TradingViewChart 
              symbol={selectedSymbol}
              onSymbolChange={setSelectedSymbol}
            />
            <MarketDataPanel 
              marketData={marketData}
              selectedSymbol={selectedSymbol}
            />
          </div>
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Recent Trading Signals</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 max-h-64 overflow-y-auto">
                  {signals.slice(0, 10).map((signal: any) => (
                    <div key={signal.id} className="flex items-center justify-between p-2 border rounded">
                      <div className="flex items-center space-x-2">
                        <Badge variant={signal.side === 'BUY' ? 'default' : 'destructive'}>
                          {signal.side}
                        </Badge>
                        <span className="font-medium">{signal.symbol}</span>
                        <span className="text-sm text-muted-foreground">
                          {signal.quantity} @ ${signal.price}
                        </span>
                      </div>
                      <div className="text-right">
                        <div className="text-sm font-medium">
                          {(signal.confidence * 100).toFixed(0)}%
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {signal.strategyId}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <PortfolioPanel />
          </div>
        </TabsContent>

        <TabsContent value="market">
          <MarketDataPanel 
            marketData={marketData}
            selectedSymbol={selectedSymbol}
            onSymbolChange={setSelectedSymbol}
          />
        </TabsContent>

        <TabsContent value="orders">
          <OrderManagementPanel />
        </TabsContent>

        <TabsContent value="portfolio">
          <PortfolioPanel />
        </TabsContent>

        <TabsContent value="strategies">
          <StrategyPanel />
        </TabsContent>

        <TabsContent value="risk">
          <RiskManagementPanel />
        </TabsContent>

        <TabsContent value="performance">
          <PerformancePanel />
        </TabsContent>
      </Tabs>
    </div>
  );
};
