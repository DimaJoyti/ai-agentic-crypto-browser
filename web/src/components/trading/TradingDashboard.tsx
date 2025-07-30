import React, { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Progress } from '@/components/ui/progress';
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
  CheckCircle,
  Zap,
  Target,
  Shield,
  Clock,
  Wifi,
  WifiOff,
  RefreshCw,
  Maximize2,
  Minimize2,
  Volume2,
  VolumeX,
  Bell,
  BellOff,
  Power,
  PowerOff,
  Database,
  Server,
  Cpu,
  MemoryStick,
  HardDrive,
  Network,
  Eye,
  EyeOff,
  Filter,
  Search,
  Download,
  Upload,
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  Plus,
  Minus,
  X,
  Check,
  Info,
  AlertTriangle as Warning,
  AlertCircle,
  Trash2,
  Edit,
  Copy,
  ExternalLink
} from 'lucide-react';

import { MarketDataPanel } from './MarketDataPanel';
import { OrderManagementPanel } from './OrderManagementPanel';
import { PortfolioPanel } from './PortfolioPanel';
import { StrategyPanel } from './StrategyPanel';
import { RiskManagementPanel } from './RiskManagementPanel';
import { PerformancePanel } from './PerformancePanel';
import { TradingViewChart } from './TradingViewChart';
import { MarketTicker } from './MarketTicker';
import { OrderBook } from './OrderBook';
import { useTradingDashboard } from '@/hooks/useTradingDashboard';
import { useBinanceTicker, useBinanceMultipleTickers } from '@/hooks/useBinanceMarketData';
import { binanceMarketService } from '@/services/binance/BinanceMarketService';

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
  cpuUsage: number;
  memoryUsage: number;
  networkLatency: number;
  diskUsage: number;
  activeConnections: number;
  queueDepth: number;
  throughput: number;
  errorRate: number;
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
  change24h: string;
  changePercent24h: string;
  high24h: string;
  low24h: string;
  volume24h: string;
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
  status: 'PENDING' | 'EXECUTED' | 'CANCELLED' | 'FAILED';
  executionTime?: number;
  slippage?: number;
}

interface OrderBookEntry {
  price: string;
  size: string;
  total: string;
}

interface OrderBook {
  symbol: string;
  bids: OrderBookEntry[];
  asks: OrderBookEntry[];
  lastUpdate: string;
}

interface Position {
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
}

interface SystemStatus {
  cpu: number;
  memory: number;
  disk: number;
  network: number;
  latency: number;
  uptime: string;
  connections: number;
  errors: number;
}

export const TradingDashboard: React.FC<TradingDashboardProps> = ({ className }) => {
  const [selectedSymbol, setSelectedSymbol] = useState('BTCUSDT');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [soundEnabled, setSoundEnabled] = useState(true);
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [layoutMode, setLayoutMode] = useState<'grid' | 'split' | 'focus'>('grid');
  const [selectedTimeframe, setSelectedTimeframe] = useState('1m');
  const [orderType, setOrderType] = useState<'MARKET' | 'LIMIT' | 'STOP'>('LIMIT');
  const [orderSide, setOrderSide] = useState<'BUY' | 'SELL'>('BUY');
  const [orderQuantity, setOrderQuantity] = useState('');
  const [orderPrice, setOrderPrice] = useState('');
  const [leverage, setLeverage] = useState(1);
  const [showOrderBook, setShowOrderBook] = useState(true);
  const [showTrades, setShowTrades] = useState(true);
  const [showPositions, setShowPositions] = useState(true);

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

  // Get real market data
  const { ticker: currentTicker } = useBinanceTicker(selectedSymbol);
  const popularSymbols = binanceMarketService.getPopularPairs();
  const { tickers: allTickers } = useBinanceMultipleTickers(popularSymbols);

  const mockPositions: Position[] = [
    {
      symbol: 'BTCUSDT',
      side: 'LONG',
      size: '0.5432',
      entryPrice: '42,850.00',
      markPrice: '43,250.50',
      pnl: '+217.54',
      pnlPercent: '+0.51',
      margin: '4,325.00',
      leverage: '10x',
      liquidationPrice: '38,565.00'
    },
    {
      symbol: 'ETHUSDT',
      side: 'SHORT',
      size: '2.1234',
      entryPrice: '2,650.00',
      markPrice: '2,635.75',
      pnl: '+30.23',
      pnlPercent: '+0.54',
      margin: '1,127.50',
      leverage: '5x',
      liquidationPrice: '2,915.50'
    }
  ];

  const mockSystemStatus: SystemStatus = {
    cpu: 45,
    memory: 67,
    disk: 23,
    network: 89,
    latency: 0.8,
    uptime: '99.98%',
    connections: 1247,
    errors: 3
  };

  const formatLatency = (micros: number) => {
    if (micros < 1000) return `${micros}μs`;
    if (micros < 1000000) return `${(micros / 1000).toFixed(1)}ms`;
    return `${(micros / 1000000).toFixed(2)}s`;
  };

  const formatPnL = (pnl: string) => {
    const value = parseFloat(pnl.replace(/[+$,]/g, ''));
    const isPositive = value >= 0;
    return (
      <span className={isPositive ? 'text-green-400' : 'text-red-400'}>
        {isPositive ? '+' : ''}${value.toFixed(2)}
      </span>
    );
  };

  const formatPrice = (price: string) => {
    return parseFloat(price).toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    });
  };

  const formatVolume = (volume: string) => {
    const val = parseFloat(volume);
    if (val >= 1000000) return `${(val / 1000000).toFixed(2)}M`;
    if (val >= 1000) return `${(val / 1000).toFixed(2)}K`;
    return val.toFixed(4);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'CONNECTED': return 'text-green-400';
      case 'DISCONNECTED': return 'text-red-400';
      case 'CONNECTING': return 'text-yellow-400';
      default: return 'text-gray-400';
    }
  };

  const handlePlaceOrder = () => {
    // Implementation for placing orders
    console.log('Placing order:', {
      symbol: selectedSymbol,
      side: orderSide,
      type: orderType,
      quantity: orderQuantity,
      price: orderPrice,
      leverage
    });
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
    <div className={`min-h-screen bg-gray-950 text-white ${className}`}>
      {/* Professional Trading Header */}
      <div className="border-b border-gray-800 bg-gray-900/50 backdrop-blur-sm">
        <div className="flex items-center justify-between px-6 py-3">
          <div className="flex items-center space-x-6">
            <div className="flex items-center space-x-2">
              <Zap className="w-6 h-6 text-blue-400" />
              <h1 className="text-xl font-bold text-white">HFT Pro</h1>
            </div>

            <div className="flex items-center space-x-4">
              <Select value={selectedSymbol} onValueChange={setSelectedSymbol}>
                <SelectTrigger className="w-32 bg-gray-800 border-gray-700">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="BTCUSDT">BTC/USDT</SelectItem>
                  <SelectItem value="ETHUSDT">ETH/USDT</SelectItem>
                  <SelectItem value="ADAUSDT">ADA/USDT</SelectItem>
                  <SelectItem value="SOLUSDT">SOL/USDT</SelectItem>
                </SelectContent>
              </Select>

              <Select value={selectedTimeframe} onValueChange={setSelectedTimeframe}>
                <SelectTrigger className="w-20 bg-gray-800 border-gray-700">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1s">1s</SelectItem>
                  <SelectItem value="1m">1m</SelectItem>
                  <SelectItem value="5m">5m</SelectItem>
                  <SelectItem value="15m">15m</SelectItem>
                  <SelectItem value="1h">1h</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* System Status Indicators */}
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-1">
                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-400' : 'bg-red-400'}`} />
                <span className="text-xs text-gray-400">
                  {isConnected ? 'LIVE' : 'OFFLINE'}
                </span>
              </div>

              <div className="flex items-center space-x-1">
                <Network className="w-3 h-3 text-gray-400" />
                <span className="text-xs text-gray-400">
                  {mockSystemStatus.latency}ms
                </span>
              </div>

              <div className="flex items-center space-x-1">
                <Server className="w-3 h-3 text-gray-400" />
                <span className="text-xs text-gray-400">
                  {mockSystemStatus.uptime}
                </span>
              </div>
            </div>

            {/* Control Buttons */}
            <div className="flex items-center space-x-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSoundEnabled(!soundEnabled)}
                className="text-gray-400 hover:text-white"
              >
                {soundEnabled ? <Volume2 className="w-4 h-4" /> : <VolumeX className="w-4 h-4" />}
              </Button>

              <Button
                variant="ghost"
                size="sm"
                onClick={() => setNotificationsEnabled(!notificationsEnabled)}
                className="text-gray-400 hover:text-white"
              >
                {notificationsEnabled ? <Bell className="w-4 h-4" /> : <BellOff className="w-4 h-4" />}
              </Button>

              <Button
                variant="ghost"
                size="sm"
                onClick={() => setAutoRefresh(!autoRefresh)}
                className="text-gray-400 hover:text-white"
              >
                {autoRefresh ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
              </Button>

              <Button
                variant="ghost"
                size="sm"
                className="text-gray-400 hover:text-white"
              >
                <Settings className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Professional Trading Layout */}
      <div className="flex h-[calc(100vh-80px)]">
        {/* Left Sidebar - Order Entry & Account */}
        <div className="w-80 border-r border-gray-800 bg-gray-900/30 p-4 space-y-4">
          {/* Order Entry Panel */}
          <Card className="bg-gray-900/50 border-gray-800">
            <CardHeader className="pb-3">
              <CardTitle className="text-lg text-white flex items-center">
                <Target className="w-5 h-5 mr-2 text-blue-400" />
                Quick Order
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Order Type Selector */}
              <div className="flex space-x-1 bg-gray-800 rounded-lg p-1">
                {(['MARKET', 'LIMIT', 'STOP'] as const).map((type) => (
                  <Button
                    key={type}
                    variant={orderType === type ? "default" : "ghost"}
                    size="sm"
                    className={`flex-1 ${orderType === type ? 'bg-blue-600' : 'text-gray-400'}`}
                    onClick={() => setOrderType(type)}
                  >
                    {type}
                  </Button>
                ))}
              </div>

              {/* Buy/Sell Selector */}
              <div className="flex space-x-2">
                <Button
                  variant={orderSide === 'BUY' ? "default" : "outline"}
                  className={`flex-1 ${orderSide === 'BUY' ? 'bg-green-600 hover:bg-green-700' : 'border-green-600 text-green-400'}`}
                  onClick={() => setOrderSide('BUY')}
                >
                  BUY
                </Button>
                <Button
                  variant={orderSide === 'SELL' ? "default" : "outline"}
                  className={`flex-1 ${orderSide === 'SELL' ? 'bg-red-600 hover:bg-red-700' : 'border-red-600 text-red-400'}`}
                  onClick={() => setOrderSide('SELL')}
                >
                  SELL
                </Button>
              </div>

              {/* Quantity Input */}
              <div>
                <label className="text-xs text-gray-400 mb-1 block">Quantity</label>
                <Input
                  type="number"
                  placeholder="0.00"
                  value={orderQuantity}
                  onChange={(e) => setOrderQuantity(e.target.value)}
                  className="bg-gray-800 border-gray-700 text-white"
                />
              </div>

              {/* Price Input (for LIMIT orders) */}
              {orderType !== 'MARKET' && (
                <div>
                  <label className="text-xs text-gray-400 mb-1 block">Price</label>
                  <Input
                    type="number"
                    placeholder="0.00"
                    value={orderPrice}
                    onChange={(e) => setOrderPrice(e.target.value)}
                    className="bg-gray-800 border-gray-700 text-white"
                  />
                </div>
              )}

              {/* Leverage Selector */}
              <div>
                <label className="text-xs text-gray-400 mb-2 block">Leverage: {leverage}x</label>
                <div className="flex space-x-1">
                  {[1, 2, 5, 10, 20, 50].map((lev) => (
                    <Button
                      key={lev}
                      variant={leverage === lev ? "default" : "ghost"}
                      size="sm"
                      className={`flex-1 text-xs ${leverage === lev ? 'bg-blue-600' : 'text-gray-400'}`}
                      onClick={() => setLeverage(lev)}
                    >
                      {lev}x
                    </Button>
                  ))}
                </div>
              </div>

              {/* Place Order Button */}
              <Button
                className={`w-full ${orderSide === 'BUY' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'}`}
                onClick={handlePlaceOrder}
              >
                {orderSide} {selectedSymbol}
              </Button>
            </CardContent>
          </Card>

          {/* Account Summary */}
          <Card className="bg-gray-900/50 border-gray-800">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm text-white">Account</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Balance:</span>
                <span className="text-white">$125,430.50</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Available:</span>
                <span className="text-green-400">$98,234.20</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Margin Used:</span>
                <span className="text-yellow-400">$27,196.30</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">P&L Today:</span>
                <span className="text-green-400">+$2,847.65</span>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Center - Chart and Market Data */}
        <div className="flex-1 flex flex-col">
          {/* Chart Area */}
          <div className="flex-1 p-4">
            <Card className="h-full bg-gray-900/30 border-gray-800">
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-white flex items-center">
                    <BarChart3 className="w-5 h-5 mr-2 text-blue-400" />
                    {selectedSymbol} Chart
                  </CardTitle>
                  <div className="flex items-center space-x-2">
                    <Button variant="ghost" size="sm" className="text-gray-400">
                      <Maximize2 className="w-4 h-4" />
                    </Button>
                    <Button variant="ghost" size="sm" className="text-gray-400">
                      <Settings className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="h-[calc(100%-80px)]">
                <TradingViewChart
                  symbol={selectedSymbol}
                  onSymbolChange={setSelectedSymbol}
                />
              </CardContent>
            </Card>
          </div>

          {/* Bottom Panel - Orders and Positions */}
          <div className="h-64 border-t border-gray-800 bg-gray-900/30">
            <Tabs defaultValue="positions" className="h-full">
              <div className="flex items-center justify-between px-4 py-2 border-b border-gray-800">
                <TabsList className="bg-gray-800">
                  <TabsTrigger value="positions" className="text-xs">Positions</TabsTrigger>
                  <TabsTrigger value="orders" className="text-xs">Orders</TabsTrigger>
                  <TabsTrigger value="history" className="text-xs">History</TabsTrigger>
                </TabsList>
                <div className="flex items-center space-x-2">
                  <Button variant="ghost" size="sm" className="text-gray-400">
                    <RefreshCw className="w-3 h-3" />
                  </Button>
                  <Button variant="ghost" size="sm" className="text-gray-400">
                    <Filter className="w-3 h-3" />
                  </Button>
                </div>
              </div>

              <TabsContent value="positions" className="p-4 h-[calc(100%-60px)] overflow-auto">
                <div className="space-y-2">
                  {mockPositions.map((position, index) => (
                    <div key={index} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <div className="flex items-center space-x-4">
                        <Badge variant={position.side === 'LONG' ? 'default' : 'destructive'} className="text-xs">
                          {position.side}
                        </Badge>
                        <span className="font-medium text-white">{position.symbol}</span>
                        <span className="text-sm text-gray-400">{position.size}</span>
                      </div>
                      <div className="flex items-center space-x-6 text-sm">
                        <div>
                          <div className="text-gray-400">Entry</div>
                          <div className="text-white">${position.entryPrice}</div>
                        </div>
                        <div>
                          <div className="text-gray-400">Mark</div>
                          <div className="text-white">${position.markPrice}</div>
                        </div>
                        <div>
                          <div className="text-gray-400">PnL</div>
                          <div className={position.pnl.startsWith('+') ? 'text-green-400' : 'text-red-400'}>
                            {position.pnl} ({position.pnlPercent})
                          </div>
                        </div>
                        <div>
                          <div className="text-gray-400">Margin</div>
                          <div className="text-white">${position.margin}</div>
                        </div>
                        <Button variant="ghost" size="sm" className="text-red-400 hover:text-red-300">
                          <X className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </TabsContent>

              <TabsContent value="orders" className="p-4 h-[calc(100%-60px)] overflow-auto">
                <div className="text-center text-gray-400 py-8">
                  No active orders
                </div>
              </TabsContent>

              <TabsContent value="history" className="p-4 h-[calc(100%-60px)] overflow-auto">
                <div className="text-center text-gray-400 py-8">
                  Trade history will appear here
                </div>
              </TabsContent>
            </Tabs>
          </div>
        </div>

        {/* Right Sidebar - Order Book & Market Data */}
        <div className="w-80 border-l border-gray-800 bg-gray-900/30 p-4 space-y-4">
          {/* Real Market Ticker */}
          <MarketTicker symbol={selectedSymbol} />

          {/* Real Order Book */}
          <OrderBook symbol={selectedSymbol} limit={15} className="flex-1" />

          {/* System Status */}
          <Card className="bg-gray-900/50 border-gray-800">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm text-white flex items-center">
                <Server className="w-4 h-4 mr-2 text-blue-400" />
                System Status
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span className="text-gray-400">CPU Usage</span>
                  <span className="text-white">{mockSystemStatus.cpu}%</span>
                </div>
                <Progress value={mockSystemStatus.cpu} className="h-1" />
              </div>

              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span className="text-gray-400">Memory</span>
                  <span className="text-white">{mockSystemStatus.memory}%</span>
                </div>
                <Progress value={mockSystemStatus.memory} className="h-1" />
              </div>

              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span className="text-gray-400">Network</span>
                  <span className="text-white">{mockSystemStatus.network}%</span>
                </div>
                <Progress value={mockSystemStatus.network} className="h-1" />
              </div>

              <div className="grid grid-cols-2 gap-2 text-xs pt-2">
                <div>
                  <div className="text-gray-400">Latency</div>
                  <div className="text-green-400">{mockSystemStatus.latency}ms</div>
                </div>
                <div>
                  <div className="text-gray-400">Uptime</div>
                  <div className="text-white">{mockSystemStatus.uptime}</div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* HFT Engine Status */}
          <Card className="bg-gray-900/50 border-gray-800">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm text-white flex items-center">
                <Zap className="w-4 h-4 mr-2 text-yellow-400" />
                HFT Engine
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-400">Status</span>
                <Badge variant={isHFTRunning ? "default" : "destructive"} className="text-xs">
                  {isHFTRunning ? 'RUNNING' : 'STOPPED'}
                </Badge>
              </div>

              <div className="grid grid-cols-2 gap-2 text-xs">
                <div>
                  <div className="text-gray-400">Orders/sec</div>
                  <div className="text-white">{hftMetrics?.ordersPerSecond || 0}</div>
                </div>
                <div>
                  <div className="text-gray-400">Latency</div>
                  <div className="text-green-400">{hftMetrics ? formatLatency(hftMetrics.avgLatencyMicros) : '0μs'}</div>
                </div>
                <div>
                  <div className="text-gray-400">Success Rate</div>
                  <div className="text-green-400">
                    {hftMetrics && hftMetrics.totalOrders > 0
                      ? ((hftMetrics.successfulOrders / hftMetrics.totalOrders) * 100).toFixed(1)
                      : 0}%
                  </div>
                </div>
                <div>
                  <div className="text-gray-400">Total Orders</div>
                  <div className="text-white">{hftMetrics?.totalOrders || 0}</div>
                </div>
              </div>

              <Button
                className={`w-full ${isHFTRunning ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'}`}
                onClick={isHFTRunning ? stopHFTEngine : startHFTEngine}
              >
                {isHFTRunning ? (
                  <>
                    <PowerOff className="w-3 h-3 mr-2" />
                    Stop Engine
                  </>
                ) : (
                  <>
                    <Power className="w-3 h-3 mr-2" />
                    Start Engine
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};
