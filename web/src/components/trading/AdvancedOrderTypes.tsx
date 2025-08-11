'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import { Slider } from '@/components/ui/slider'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Target, 
  TrendingUp, 
  TrendingDown,
  Zap,
  Shield,
  BarChart3,
  Settings,
  Play,
  Pause,
  Square,
  AlertTriangle,
  Info,
  Clock,
  Layers,
  Activity,
  Brain,
  Cpu,
  Timer,
  Repeat,
  Grid3X3,
  LineChart,
  PieChart,
  Scissors,
  Shuffle,
  ArrowUpDown,
  Calculator,
  Gauge,
  Crosshair,
  Radar,
  Waves,
  Sparkles,
  Rocket,
  Bot,
  Eye,
  Edit,
  Trash2,
  Copy,
  Download,
  Upload,
  RefreshCw,
  CheckCircle,
  XCircle,
  DollarSign,
  Percent,
  Plus,
  Star
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface AdvancedOrder {
  id: string
  type: 'oco' | 'trailing_stop' | 'iceberg' | 'twap' | 'vwap' | 'grid' | 'bracket' | 'stop_limit' | 'hidden' | 'post_only' | 'reduce_only'
  symbol: string
  side: 'buy' | 'sell'
  status: 'active' | 'paused' | 'completed' | 'cancelled' | 'partially_filled' | 'pending' | 'rejected'
  createdAt: number
  updatedAt: number
  parameters: Record<string, any>
  progress?: number
  pnl?: number
  filledQuantity?: number
  totalQuantity?: number
  averagePrice?: number
  estimatedCompletion?: number
  riskScore?: number
  slippage?: number
  fees?: number
  venue?: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
}

interface AlgoStrategy {
  id: string
  name: string
  description: string
  type: 'dca' | 'grid' | 'momentum' | 'mean_reversion' | 'arbitrage' | 'market_making' | 'pairs_trading' | 'statistical_arbitrage'
  category: 'conservative' | 'moderate' | 'aggressive' | 'experimental'
  isActive: boolean
  performance: number
  totalTrades: number
  winRate: number
  sharpeRatio: number
  maxDrawdown: number
  totalPnL: number
  dailyPnL: number
  parameters: Record<string, any>
  riskMetrics: {
    var95: number
    var99: number
    expectedShortfall: number
    beta: number
    alpha: number
  }
  backtestResults?: {
    period: string
    totalReturn: number
    annualizedReturn: number
    volatility: number
    maxDrawdown: number
    calmarRatio: number
  }
}

interface OrderTemplate {
  id: string
  name: string
  description: string
  type: string
  parameters: Record<string, any>
  isPublic: boolean
  createdBy: string
  usageCount: number
  rating: number
}

interface RiskParameters {
  maxPositionSize: number
  maxDailyLoss: number
  maxOrderValue: number
  allowedSymbols: string[]
  forbiddenSymbols: string[]
  maxLeverage: number
  requireConfirmation: boolean
  stopLossRequired: boolean
  takeProfitRequired: boolean
}

export function AdvancedOrderTypes() {
  const [activeOrders, setActiveOrders] = useState<AdvancedOrder[]>([])
  const [algoStrategies, setAlgoStrategies] = useState<AlgoStrategy[]>([])
  const [orderTemplates, setOrderTemplates] = useState<OrderTemplate[]>([])
  const [riskParameters, setRiskParameters] = useState<RiskParameters | null>(null)
  const [selectedOrderType, setSelectedOrderType] = useState<string>('oco')
  const [orderParams, setOrderParams] = useState<Record<string, any>>({})
  const [activeTab, setActiveTab] = useState('orders')
  const [selectedStrategy, setSelectedStrategy] = useState<string | null>(null)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate comprehensive mock data
    const mockOrders: AdvancedOrder[] = [
      {
        id: 'order1',
        type: 'oco',
        symbol: 'ETH/USDT',
        side: 'sell',
        status: 'active',
        createdAt: Date.now() - 3600000,
        updatedAt: Date.now() - 1800000,
        priority: 'high',
        parameters: {
          amount: 2.5,
          takeProfitPrice: 2600,
          stopLossPrice: 2400,
          trailingStop: false
        },
        progress: 0,
        pnl: 125.50,
        filledQuantity: 0,
        totalQuantity: 2.5,
        riskScore: 25,
        venue: 'Binance',
        fees: 3.75
      },
      {
        id: 'order2',
        type: 'trailing_stop',
        symbol: 'BTC/USDT',
        side: 'sell',
        status: 'active',
        createdAt: Date.now() - 7200000,
        updatedAt: Date.now() - 900000,
        priority: 'medium',
        parameters: {
          amount: 0.1,
          trailAmount: 500,
          currentStopPrice: 44500,
          trailPercent: 2.5,
          activationPrice: 45000
        },
        progress: 0,
        pnl: -45.20,
        filledQuantity: 0,
        totalQuantity: 0.1,
        riskScore: 35,
        venue: 'Coinbase',
        fees: 4.50
      },
      {
        id: 'order3',
        type: 'iceberg',
        symbol: 'ETH/USDT',
        side: 'buy',
        status: 'partially_filled',
        createdAt: Date.now() - 1800000,
        updatedAt: Date.now() - 300000,
        priority: 'medium',
        parameters: {
          totalAmount: 10,
          visibleAmount: 1,
          price: 2480,
          randomizeSize: true,
          minSize: 0.5,
          maxSize: 1.5,
          priceImprovement: 0.1
        },
        progress: 35,
        pnl: 89.75,
        filledQuantity: 3.5,
        totalQuantity: 10,
        averagePrice: 2475.30,
        riskScore: 15,
        venue: 'Kraken',
        fees: 12.40
      },
      {
        id: 'order4',
        type: 'twap',
        symbol: 'BTC/USDT',
        side: 'buy',
        status: 'active',
        createdAt: Date.now() - 900000,
        updatedAt: Date.now() - 60000,
        priority: 'low',
        parameters: {
          totalAmount: 0.5,
          duration: 3600,
          intervals: 12,
          priceLimit: 45500,
          adaptiveSize: true,
          marketImpactLimit: 0.1
        },
        progress: 60,
        pnl: 0,
        filledQuantity: 0.3,
        totalQuantity: 0.5,
        averagePrice: 45125.80,
        estimatedCompletion: Date.now() + 1440000,
        riskScore: 20,
        venue: 'Binance',
        fees: 6.75
      },
      {
        id: 'order5',
        type: 'vwap',
        symbol: 'ETH/USDT',
        side: 'sell',
        status: 'active',
        createdAt: Date.now() - 600000,
        updatedAt: Date.now() - 120000,
        priority: 'high',
        parameters: {
          totalAmount: 5,
          duration: 1800,
          volumeParticipation: 0.15,
          priceLimit: 2450,
          aggressiveness: 'medium'
        },
        progress: 25,
        pnl: 45.80,
        filledQuantity: 1.25,
        totalQuantity: 5,
        averagePrice: 2487.60,
        estimatedCompletion: Date.now() + 1200000,
        riskScore: 18,
        venue: 'Binance',
        fees: 9.25
      }
    ]

    const mockStrategies: AlgoStrategy[] = [
      {
        id: 'strategy1',
        name: 'ETH DCA Strategy',
        description: 'Dollar Cost Averaging for ETH with weekly purchases and dynamic sizing',
        type: 'dca',
        category: 'conservative',
        isActive: true,
        performance: 15.2,
        totalTrades: 24,
        winRate: 75,
        sharpeRatio: 1.45,
        maxDrawdown: 8.5,
        totalPnL: 1250.75,
        dailyPnL: 12.50,
        parameters: {
          symbol: 'ETH/USDT',
          baseAmount: 100,
          interval: 'weekly',
          priceDeviation: 5,
          stopLoss: 15,
          takeProfit: 25,
          dynamicSizing: true,
          volatilityAdjustment: true
        },
        riskMetrics: {
          var95: 2.5,
          var99: 4.2,
          expectedShortfall: 5.1,
          beta: 0.85,
          alpha: 0.12
        },
        backtestResults: {
          period: '1Y',
          totalReturn: 18.5,
          annualizedReturn: 18.5,
          volatility: 12.8,
          maxDrawdown: 8.5,
          calmarRatio: 2.17
        }
      },
      {
        id: 'strategy2',
        name: 'BTC Grid Bot Pro',
        description: 'Advanced grid trading with dynamic range adjustment and ML optimization',
        type: 'grid',
        category: 'moderate',
        isActive: true,
        performance: 8.7,
        totalTrades: 156,
        winRate: 68,
        sharpeRatio: 0.95,
        maxDrawdown: 12.3,
        totalPnL: 2875.40,
        dailyPnL: 23.80,
        parameters: {
          symbol: 'BTC/USDT',
          lowerPrice: 40000,
          upperPrice: 50000,
          gridLevels: 20,
          investment: 5000,
          rebalanceThreshold: 5,
          dynamicRange: true,
          mlOptimization: true,
          riskManagement: 'adaptive'
        },
        riskMetrics: {
          var95: 3.8,
          var99: 6.2,
          expectedShortfall: 7.5,
          beta: 1.15,
          alpha: 0.08
        }
      }
    ]

    setActiveOrders(mockOrders)
    setAlgoStrategies(mockStrategies)
  }, [isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-500'
      case 'partially_filled': return 'text-blue-500'
      case 'paused': return 'text-yellow-500'
      case 'completed': return 'text-gray-500'
      case 'cancelled': case 'rejected': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'active': return 'default'
      case 'partially_filled': return 'secondary'
      case 'paused': return 'outline'
      case 'completed': return 'outline'
      case 'cancelled': case 'rejected': return 'destructive'
      default: return 'outline'
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'urgent': return 'text-red-500'
      case 'high': return 'text-orange-500'
      case 'medium': return 'text-yellow-500'
      case 'low': return 'text-green-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskScoreColor = (score: number) => {
    if (score <= 20) return 'text-green-500'
    if (score <= 40) return 'text-yellow-500'
    if (score <= 60) return 'text-orange-500'
    return 'text-red-500'
  }

  const getOrderTypeIcon = (type: string) => {
    switch (type) {
      case 'oco': return <Target className="w-4 h-4" />
      case 'trailing_stop': return <TrendingUp className="w-4 h-4" />
      case 'iceberg': return <Layers className="w-4 h-4" />
      case 'twap': return <Clock className="w-4 h-4" />
      case 'vwap': return <BarChart3 className="w-4 h-4" />
      case 'grid': return <Grid3X3 className="w-4 h-4" />
      case 'bracket': return <Scissors className="w-4 h-4" />
      case 'stop_limit': return <Shield className="w-4 h-4" />
      case 'hidden': return <Eye className="w-4 h-4" />
      case 'post_only': return <Crosshair className="w-4 h-4" />
      case 'reduce_only': return <TrendingDown className="w-4 h-4" />
      default: return <Activity className="w-4 h-4" />
    }
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Brain className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access advanced order types and algorithmic trading strategies
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Advanced Order Types</h2>
          <p className="text-muted-foreground">
            Professional trading tools with algorithmic execution strategies
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Bot className="w-3 h-3 mr-1" />
            {activeOrders.filter(o => o.status === 'active').length} Active
          </Badge>
          <Badge variant="outline">
            <Zap className="w-3 h-3 mr-1" />
            {algoStrategies.filter(s => s.isActive).length} Strategies
          </Badge>
        </div>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="orders">Active Orders</TabsTrigger>
          <TabsTrigger value="create">Create Order</TabsTrigger>
          <TabsTrigger value="strategies">Algo Strategies</TabsTrigger>
          <TabsTrigger value="templates">Templates</TabsTrigger>
          <TabsTrigger value="settings">Risk Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="create" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Rocket className="w-5 h-5" />
                Create Advanced Order
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {/* Order Type Selection */}
                <div>
                  <Label className="text-base font-medium">Order Type</Label>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mt-2">
                    {[
                      { type: 'oco', name: 'OCO', description: 'One-Cancels-Other', icon: Target },
                      { type: 'trailing_stop', name: 'Trailing Stop', description: 'Dynamic stop loss', icon: TrendingUp },
                      { type: 'iceberg', name: 'Iceberg', description: 'Hidden large orders', icon: Layers },
                      { type: 'twap', name: 'TWAP', description: 'Time-weighted average', icon: Clock },
                      { type: 'vwap', name: 'VWAP', description: 'Volume-weighted average', icon: BarChart3 },
                      { type: 'grid', name: 'Grid', description: 'Grid trading bot', icon: Grid3X3 },
                      { type: 'bracket', name: 'Bracket', description: 'Stop + limit combo', icon: Scissors },
                      { type: 'stop_limit', name: 'Stop Limit', description: 'Conditional limit order', icon: Shield }
                    ].map((orderType) => {
                      const Icon = orderType.icon
                      return (
                        <div
                          key={orderType.type}
                          className={cn(
                            "p-3 border rounded-lg cursor-pointer transition-colors",
                            selectedOrderType === orderType.type ? "border-primary bg-primary/5" : "border-muted hover:border-primary/50"
                          )}
                          onClick={() => setSelectedOrderType(orderType.type)}
                        >
                          <Icon className="w-5 h-5 mb-2" />
                          <div className="font-medium text-sm">{orderType.name}</div>
                          <div className="text-xs text-muted-foreground">{orderType.description}</div>
                        </div>
                      )
                    })}
                  </div>
                </div>

                {/* Basic Parameters */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <Label htmlFor="symbol">Trading Pair</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Select pair" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="BTC/USDT">BTC/USDT</SelectItem>
                        <SelectItem value="ETH/USDT">ETH/USDT</SelectItem>
                        <SelectItem value="BNB/USDT">BNB/USDT</SelectItem>
                        <SelectItem value="ADA/USDT">ADA/USDT</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <Label htmlFor="side">Side</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Buy/Sell" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="buy">Buy</SelectItem>
                        <SelectItem value="sell">Sell</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <Label htmlFor="quantity">Quantity</Label>
                    <Input
                      id="quantity"
                      type="number"
                      placeholder="0.00"
                      step="0.001"
                    />
                  </div>
                </div>

                {/* Order-Specific Parameters */}
                {selectedOrderType === 'oco' && (
                  <div className="space-y-4 p-4 border rounded-lg bg-muted/20">
                    <h4 className="font-medium">OCO Parameters</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="takeProfitPrice">Take Profit Price</Label>
                        <Input
                          id="takeProfitPrice"
                          type="number"
                          placeholder="0.00"
                          step="0.01"
                        />
                      </div>
                      <div>
                        <Label htmlFor="stopLossPrice">Stop Loss Price</Label>
                        <Input
                          id="stopLossPrice"
                          type="number"
                          placeholder="0.00"
                          step="0.01"
                        />
                      </div>
                      <div className="flex items-center space-x-2">
                        <Switch id="trailingStop" />
                        <Label htmlFor="trailingStop">Enable Trailing Stop</Label>
                      </div>
                    </div>
                  </div>
                )}

                {selectedOrderType === 'trailing_stop' && (
                  <div className="space-y-4 p-4 border rounded-lg bg-muted/20">
                    <h4 className="font-medium">Trailing Stop Parameters</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="trailAmount">Trail Amount (USD)</Label>
                        <Input
                          id="trailAmount"
                          type="number"
                          placeholder="100"
                          step="1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="trailPercent">Trail Percentage (%)</Label>
                        <Input
                          id="trailPercent"
                          type="number"
                          placeholder="2.5"
                          step="0.1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="activationPrice">Activation Price</Label>
                        <Input
                          id="activationPrice"
                          type="number"
                          placeholder="0.00"
                          step="0.01"
                        />
                      </div>
                    </div>
                  </div>
                )}

                {selectedOrderType === 'iceberg' && (
                  <div className="space-y-4 p-4 border rounded-lg bg-muted/20">
                    <h4 className="font-medium">Iceberg Parameters</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="visibleAmount">Visible Amount</Label>
                        <Input
                          id="visibleAmount"
                          type="number"
                          placeholder="1.0"
                          step="0.1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="priceImprovement">Price Improvement (%)</Label>
                        <Input
                          id="priceImprovement"
                          type="number"
                          placeholder="0.1"
                          step="0.01"
                        />
                      </div>
                      <div className="flex items-center space-x-2">
                        <Switch id="randomizeSize" />
                        <Label htmlFor="randomizeSize">Randomize Size</Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Switch id="timeRandomization" />
                        <Label htmlFor="timeRandomization">Time Randomization</Label>
                      </div>
                    </div>
                  </div>
                )}

                {selectedOrderType === 'twap' && (
                  <div className="space-y-4 p-4 border rounded-lg bg-muted/20">
                    <h4 className="font-medium">TWAP Parameters</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="duration">Duration (minutes)</Label>
                        <Input
                          id="duration"
                          type="number"
                          placeholder="60"
                          step="1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="intervals">Number of Intervals</Label>
                        <Input
                          id="intervals"
                          type="number"
                          placeholder="12"
                          step="1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="priceLimit">Price Limit</Label>
                        <Input
                          id="priceLimit"
                          type="number"
                          placeholder="0.00"
                          step="0.01"
                        />
                      </div>
                      <div className="flex items-center space-x-2">
                        <Switch id="adaptiveSize" />
                        <Label htmlFor="adaptiveSize">Adaptive Sizing</Label>
                      </div>
                    </div>
                  </div>
                )}

                {selectedOrderType === 'vwap' && (
                  <div className="space-y-4 p-4 border rounded-lg bg-muted/20">
                    <h4 className="font-medium">VWAP Parameters</h4>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="duration">Duration (minutes)</Label>
                        <Input
                          id="duration"
                          type="number"
                          placeholder="30"
                          step="1"
                        />
                      </div>
                      <div>
                        <Label htmlFor="volumeParticipation">Volume Participation (%)</Label>
                        <Input
                          id="volumeParticipation"
                          type="number"
                          placeholder="15"
                          step="1"
                          max="50"
                        />
                      </div>
                      <div>
                        <Label htmlFor="aggressiveness">Aggressiveness</Label>
                        <Select>
                          <SelectTrigger>
                            <SelectValue placeholder="Select level" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="low">Low</SelectItem>
                            <SelectItem value="medium">Medium</SelectItem>
                            <SelectItem value="high">High</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Switch id="marketImpactLimit" />
                        <Label htmlFor="marketImpactLimit">Market Impact Limit</Label>
                      </div>
                    </div>
                  </div>
                )}

                {/* Risk Management */}
                <div className="space-y-4 p-4 border rounded-lg bg-yellow-50 dark:bg-yellow-950/20">
                  <h4 className="font-medium flex items-center gap-2">
                    <Shield className="w-4 h-4" />
                    Risk Management
                  </h4>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <Label htmlFor="maxSlippage">Max Slippage (%)</Label>
                      <Input
                        id="maxSlippage"
                        type="number"
                        placeholder="0.5"
                        step="0.1"
                      />
                    </div>
                    <div>
                      <Label htmlFor="priority">Priority</Label>
                      <Select>
                        <SelectTrigger>
                          <SelectValue placeholder="Select priority" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="low">Low</SelectItem>
                          <SelectItem value="medium">Medium</SelectItem>
                          <SelectItem value="high">High</SelectItem>
                          <SelectItem value="urgent">Urgent</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div>
                      <Label htmlFor="venue">Preferred Venue</Label>
                      <Select>
                        <SelectTrigger>
                          <SelectValue placeholder="Auto-select" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="auto">Auto-select</SelectItem>
                          <SelectItem value="binance">Binance</SelectItem>
                          <SelectItem value="coinbase">Coinbase</SelectItem>
                          <SelectItem value="kraken">Kraken</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Switch id="requireConfirmation" />
                    <Label htmlFor="requireConfirmation">Require confirmation before execution</Label>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex gap-3">
                  <Button className="flex-1">
                    <Rocket className="w-4 h-4 mr-2" />
                    Create Order
                  </Button>
                  <Button variant="outline">
                    <Copy className="w-4 h-4 mr-2" />
                    Save as Template
                  </Button>
                  <Button variant="outline">
                    <Calculator className="w-4 h-4 mr-2" />
                    Simulate
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="orders" className="space-y-4">
          <div className="grid grid-cols-1 gap-4">
            {activeOrders.map((order) => (
              <Card key={order.id}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        {getOrderTypeIcon(order.type)}
                      </div>
                      <div>
                        <h4 className="font-bold">{order.type.toUpperCase().replace('_', ' ')} Order</h4>
                        <p className="text-sm text-muted-foreground">{order.symbol} • {order.side.toUpperCase()}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={getStatusBadgeVariant(order.status)}>
                        {order.status.replace('_', ' ')}
                      </Badge>
                      <Badge variant="outline" className={getPriorityColor(order.priority)}>
                        {order.priority}
                      </Badge>
                    </div>
                  </div>

                  {/* Progress Bar */}
                  {order.progress !== undefined && (
                    <div className="mb-4">
                      <div className="flex justify-between text-sm mb-1">
                        <span>Progress</span>
                        <span>{order.progress}%</span>
                      </div>
                      <Progress value={order.progress} className="h-2" />
                    </div>
                  )}

                  {/* Order Details */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Total Quantity</div>
                      <div className="font-medium">{order.totalQuantity}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Filled</div>
                      <div className="font-medium">{order.filledQuantity || 0}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">P&L</div>
                      <div className={cn(
                        "font-medium",
                        order.pnl && order.pnl >= 0 ? "text-green-500" : "text-red-500"
                      )}>
                        {order.pnl ? formatCurrency(order.pnl) : '--'}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Risk Score</div>
                      <div className={cn("font-medium", getRiskScoreColor(order.riskScore || 0))}>
                        {order.riskScore || 0}/100
                      </div>
                    </div>
                  </div>

                  {/* Order Parameters */}
                  <div className="space-y-2 mb-4">
                    <div className="text-sm font-medium">Parameters</div>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-2 text-xs">
                      {Object.entries(order.parameters).map(([key, value]) => (
                        <div key={key} className="flex justify-between p-2 bg-muted rounded">
                          <span className="text-muted-foreground">{key.replace(/([A-Z])/g, ' $1').toLowerCase()}</span>
                          <span className="font-medium">{typeof value === 'number' ? value.toLocaleString() : value}</span>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm">
                      <Edit className="w-3 h-3 mr-1" />
                      Modify
                    </Button>
                    <Button variant="outline" size="sm">
                      <Pause className="w-3 h-3 mr-1" />
                      Pause
                    </Button>
                    <Button variant="outline" size="sm">
                      <Square className="w-3 h-3 mr-1" />
                      Cancel
                    </Button>
                    <Button variant="outline" size="sm">
                      <Eye className="w-3 h-3 mr-1" />
                      Details
                    </Button>
                  </div>

                  {/* Footer */}
                  <div className="flex items-center justify-between mt-4 pt-4 border-t text-xs text-muted-foreground">
                    <span>Created: {formatTime(order.createdAt)}</span>
                    <span>Venue: {order.venue}</span>
                    <span>Fees: {order.fees ? formatCurrency(order.fees) : '--'}</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {activeOrders.length === 0 && (
            <div className="text-center py-12">
              <Bot className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="text-lg font-medium mb-2">No Active Orders</h3>
              <p className="text-muted-foreground mb-4">
                Create your first advanced order to get started
              </p>
              <Button onClick={() => setActiveTab('create')}>
                <Plus className="w-4 h-4 mr-2" />
                Create Order
              </Button>
            </div>
          )}
        </TabsContent>

        <TabsContent value="strategies" className="space-y-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium">Algorithmic Strategies</h3>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Create Strategy
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {algoStrategies.map((strategy) => (
              <Card key={strategy.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Brain className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{strategy.name}</h4>
                        <p className="text-sm text-muted-foreground">{strategy.type.toUpperCase()}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={strategy.category === 'conservative' ? 'default' :
                                   strategy.category === 'moderate' ? 'secondary' : 'destructive'}>
                        {strategy.category}
                      </Badge>
                      {strategy.isActive ? (
                        <Badge variant="default">
                          <Play className="w-3 h-3 mr-1" />
                          Active
                        </Badge>
                      ) : (
                        <Badge variant="outline">
                          <Pause className="w-3 h-3 mr-1" />
                          Paused
                        </Badge>
                      )}
                    </div>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  <p className="text-sm text-muted-foreground">{strategy.description}</p>

                  {/* Performance Metrics */}
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Performance</div>
                      <div className={cn(
                        "font-bold text-lg",
                        strategy.performance >= 0 ? "text-green-500" : "text-red-500"
                      )}>
                        {strategy.performance >= 0 ? '+' : ''}{strategy.performance.toFixed(1)}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Total P&L</div>
                      <div className={cn(
                        "font-bold text-lg",
                        strategy.totalPnL >= 0 ? "text-green-500" : "text-red-500"
                      )}>
                        {formatCurrency(strategy.totalPnL)}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Win Rate</div>
                      <div className="font-medium">{strategy.winRate}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Sharpe Ratio</div>
                      <div className="font-medium">{strategy.sharpeRatio.toFixed(2)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Max Drawdown</div>
                      <div className="font-medium text-red-500">{strategy.maxDrawdown.toFixed(1)}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Total Trades</div>
                      <div className="font-medium">{strategy.totalTrades}</div>
                    </div>
                  </div>

                  {/* Risk Metrics */}
                  <div className="space-y-2">
                    <h5 className="font-medium text-sm">Risk Metrics</h5>
                    <div className="grid grid-cols-2 gap-2 text-xs">
                      <div className="flex justify-between p-2 bg-muted rounded">
                        <span>VaR 95%</span>
                        <span className="font-medium">{strategy.riskMetrics.var95.toFixed(1)}%</span>
                      </div>
                      <div className="flex justify-between p-2 bg-muted rounded">
                        <span>Beta</span>
                        <span className="font-medium">{strategy.riskMetrics.beta.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between p-2 bg-muted rounded">
                        <span>Alpha</span>
                        <span className="font-medium">{strategy.riskMetrics.alpha.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between p-2 bg-muted rounded">
                        <span>Expected Shortfall</span>
                        <span className="font-medium">{strategy.riskMetrics.expectedShortfall.toFixed(1)}%</span>
                      </div>
                    </div>
                  </div>

                  {/* Backtest Results */}
                  {strategy.backtestResults && (
                    <div className="space-y-2">
                      <h5 className="font-medium text-sm">Backtest Results ({strategy.backtestResults.period})</h5>
                      <div className="grid grid-cols-2 gap-2 text-xs">
                        <div className="flex justify-between p-2 bg-muted rounded">
                          <span>Total Return</span>
                          <span className="font-medium text-green-500">
                            {strategy.backtestResults.totalReturn.toFixed(1)}%
                          </span>
                        </div>
                        <div className="flex justify-between p-2 bg-muted rounded">
                          <span>Volatility</span>
                          <span className="font-medium">{strategy.backtestResults.volatility.toFixed(1)}%</span>
                        </div>
                        <div className="flex justify-between p-2 bg-muted rounded">
                          <span>Calmar Ratio</span>
                          <span className="font-medium">{strategy.backtestResults.calmarRatio.toFixed(2)}</span>
                        </div>
                        <div className="flex justify-between p-2 bg-muted rounded">
                          <span>Max DD</span>
                          <span className="font-medium text-red-500">
                            {strategy.backtestResults.maxDrawdown.toFixed(1)}%
                          </span>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Strategy Parameters */}
                  <div className="space-y-2">
                    <h5 className="font-medium text-sm">Parameters</h5>
                    <div className="space-y-1">
                      {Object.entries(strategy.parameters).slice(0, 4).map(([key, value]) => (
                        <div key={key} className="flex justify-between text-xs">
                          <span className="text-muted-foreground">
                            {key.replace(/([A-Z])/g, ' $1').toLowerCase()}
                          </span>
                          <span className="font-medium">
                            {typeof value === 'number' ? value.toLocaleString() : value}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    <Button
                      variant={strategy.isActive ? "outline" : "default"}
                      size="sm"
                      className="flex-1"
                    >
                      {strategy.isActive ? (
                        <>
                          <Pause className="w-3 h-3 mr-1" />
                          Pause
                        </>
                      ) : (
                        <>
                          <Play className="w-3 h-3 mr-1" />
                          Start
                        </>
                      )}
                    </Button>
                    <Button variant="outline" size="sm">
                      <Settings className="w-3 h-3 mr-1" />
                      Configure
                    </Button>
                    <Button variant="outline" size="sm">
                      <BarChart3 className="w-3 h-3 mr-1" />
                      Analytics
                    </Button>
                  </div>

                  {/* Daily P&L */}
                  <div className="text-xs text-muted-foreground">
                    Daily P&L: <span className={cn(
                      "font-medium",
                      strategy.dailyPnL >= 0 ? "text-green-500" : "text-red-500"
                    )}>
                      {strategy.dailyPnL >= 0 ? '+' : ''}{formatCurrency(strategy.dailyPnL)}
                    </span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Strategy Templates */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Sparkles className="w-5 h-5" />
                Strategy Templates
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {[
                  {
                    name: 'DCA Pro',
                    description: 'Advanced dollar-cost averaging with volatility adjustment',
                    type: 'dca',
                    difficulty: 'Beginner',
                    expectedReturn: '12-18%',
                    risk: 'Low'
                  },
                  {
                    name: 'Grid Master',
                    description: 'Dynamic grid trading with ML-powered range detection',
                    type: 'grid',
                    difficulty: 'Intermediate',
                    expectedReturn: '15-25%',
                    risk: 'Medium'
                  },
                  {
                    name: 'Momentum Scalper',
                    description: 'High-frequency momentum trading with advanced signals',
                    type: 'momentum',
                    difficulty: 'Advanced',
                    expectedReturn: '20-35%',
                    risk: 'High'
                  }
                ].map((template, index) => (
                  <div key={index} className="p-4 border rounded-lg hover:border-primary/50 cursor-pointer transition-colors">
                    <h4 className="font-medium">{template.name}</h4>
                    <p className="text-sm text-muted-foreground mb-3">{template.description}</p>
                    <div className="space-y-2 text-xs">
                      <div className="flex justify-between">
                        <span>Difficulty:</span>
                        <Badge variant="outline" className="text-xs">{template.difficulty}</Badge>
                      </div>
                      <div className="flex justify-between">
                        <span>Expected Return:</span>
                        <span className="font-medium text-green-500">{template.expectedReturn}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Risk Level:</span>
                        <span className={cn(
                          "font-medium",
                          template.risk === 'Low' ? 'text-green-500' :
                          template.risk === 'Medium' ? 'text-yellow-500' : 'text-red-500'
                        )}>{template.risk}</span>
                      </div>
                    </div>
                    <Button variant="outline" size="sm" className="w-full mt-3">
                      Use Template
                    </Button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="templates" className="space-y-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium">Order Templates</h3>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Create Template
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {orderTemplates.map((template) => (
              <Card key={template.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <h4 className="font-bold">{template.name}</h4>
                    <div className="flex items-center gap-1">
                      <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                      <span className="text-xs">{template.rating}</span>
                    </div>
                  </div>
                  <p className="text-sm text-muted-foreground">{template.description}</p>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-sm">
                      <span>Type:</span>
                      <Badge variant="outline">{template.type.toUpperCase()}</Badge>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Usage:</span>
                      <span className="font-medium">{template.usageCount.toLocaleString()}</span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Created by:</span>
                      <span className="font-medium">{template.createdBy}</span>
                    </div>

                    <div className="space-y-1">
                      <div className="text-sm font-medium">Parameters:</div>
                      {Object.entries(template.parameters).slice(0, 3).map(([key, value]) => (
                        <div key={key} className="flex justify-between text-xs">
                          <span className="text-muted-foreground">
                            {key.replace(/([A-Z])/g, ' $1').toLowerCase()}
                          </span>
                          <span className="font-medium">
                            {typeof value === 'boolean' ? (value ? 'Yes' : 'No') : value}
                          </span>
                        </div>
                      ))}
                    </div>

                    <div className="flex gap-2">
                      <Button size="sm" className="flex-1">
                        <Copy className="w-3 h-3 mr-1" />
                        Use
                      </Button>
                      <Button variant="outline" size="sm">
                        <Eye className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="settings" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="w-5 h-5" />
                Risk Management Settings
              </CardTitle>
            </CardHeader>
            <CardContent>
              {riskParameters && (
                <div className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="maxPositionSize">Max Position Size (USD)</Label>
                      <Input
                        id="maxPositionSize"
                        type="number"
                        value={riskParameters.maxPositionSize}
                        onChange={(e) => setRiskParameters(prev => prev ?
                          { ...prev, maxPositionSize: parseFloat(e.target.value) } : null
                        )}
                      />
                    </div>
                    <div>
                      <Label htmlFor="maxDailyLoss">Max Daily Loss (USD)</Label>
                      <Input
                        id="maxDailyLoss"
                        type="number"
                        value={riskParameters.maxDailyLoss}
                        onChange={(e) => setRiskParameters(prev => prev ?
                          { ...prev, maxDailyLoss: parseFloat(e.target.value) } : null
                        )}
                      />
                    </div>
                    <div>
                      <Label htmlFor="maxOrderValue">Max Order Value (USD)</Label>
                      <Input
                        id="maxOrderValue"
                        type="number"
                        value={riskParameters.maxOrderValue}
                        onChange={(e) => setRiskParameters(prev => prev ?
                          { ...prev, maxOrderValue: parseFloat(e.target.value) } : null
                        )}
                      />
                    </div>
                    <div>
                      <Label htmlFor="maxLeverage">Max Leverage</Label>
                      <Input
                        id="maxLeverage"
                        type="number"
                        value={riskParameters.maxLeverage}
                        onChange={(e) => setRiskParameters(prev => prev ?
                          { ...prev, maxLeverage: parseFloat(e.target.value) } : null
                        )}
                      />
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="requireConfirmation"
                        checked={riskParameters.requireConfirmation}
                        onCheckedChange={(checked) => setRiskParameters(prev => prev ?
                          { ...prev, requireConfirmation: checked } : null
                        )}
                      />
                      <Label htmlFor="requireConfirmation">Require confirmation for all orders</Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="stopLossRequired"
                        checked={riskParameters.stopLossRequired}
                        onCheckedChange={(checked) => setRiskParameters(prev => prev ?
                          { ...prev, stopLossRequired: checked } : null
                        )}
                      />
                      <Label htmlFor="stopLossRequired">Require stop loss for all positions</Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="takeProfitRequired"
                        checked={riskParameters.takeProfitRequired}
                        onCheckedChange={(checked) => setRiskParameters(prev => prev ?
                          { ...prev, takeProfitRequired: checked } : null
                        )}
                      />
                      <Label htmlFor="takeProfitRequired">Require take profit for all positions</Label>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div>
                      <Label>Allowed Trading Pairs</Label>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {riskParameters.allowedSymbols.map((symbol) => (
                          <Badge key={symbol} variant="default">
                            {symbol}
                            <button className="ml-1 hover:text-red-500">×</button>
                          </Badge>
                        ))}
                        <Button variant="outline" size="sm">
                          <Plus className="w-3 h-3 mr-1" />
                          Add Pair
                        </Button>
                      </div>
                    </div>

                    <div>
                      <Label>Forbidden Trading Pairs</Label>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {riskParameters.forbiddenSymbols.map((symbol) => (
                          <Badge key={symbol} variant="destructive">
                            {symbol}
                            <button className="ml-1 hover:text-white">×</button>
                          </Badge>
                        ))}
                        <Button variant="outline" size="sm">
                          <Plus className="w-3 h-3 mr-1" />
                          Add Pair
                        </Button>
                      </div>
                    </div>
                  </div>

                  <div className="flex gap-3">
                    <Button>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      Save Settings
                    </Button>
                    <Button variant="outline">
                      <RefreshCw className="w-4 h-4 mr-2" />
                      Reset to Defaults
                    </Button>
                    <Button variant="outline">
                      <Download className="w-4 h-4 mr-2" />
                      Export Config
                    </Button>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
