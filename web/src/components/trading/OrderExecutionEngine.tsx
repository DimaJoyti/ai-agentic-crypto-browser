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
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Zap,
  Activity,
  Clock,
  Target,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Shield,
  Settings,
  Play,
  Pause,
  Square,
  RefreshCw,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Eye,
  Edit,
  Trash2,
  Copy,
  Download,
  Gauge,
  Cpu,
  Wifi,
  WifiOff,
  Server,
  Database,
  Globe,
  MapPin,
  Timer,
  Layers,
  Grid3X3,
  Crosshair,
  Radar,
  Bot,
  Brain,
  Sparkles
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface ExecutionVenue {
  id: string
  name: string
  type: 'cex' | 'dex' | 'dark_pool' | 'otc'
  latency: number
  fees: number
  liquidity: number
  reliability: number
  isConnected: boolean
  lastPing: number
  supportedPairs: string[]
  minOrderSize: number
  maxOrderSize: number
}

interface ExecutionAlgorithm {
  id: string
  name: string
  description: string
  type: 'twap' | 'vwap' | 'implementation_shortfall' | 'arrival_price' | 'pov' | 'iceberg' | 'sniper'
  parameters: Record<string, any>
  isActive: boolean
  performance: {
    avgSlippage: number
    fillRate: number
    marketImpact: number
    executionTime: number
  }
}

interface OrderExecution {
  id: string
  orderId: string
  symbol: string
  side: 'buy' | 'sell'
  quantity: number
  algorithm: string
  venue: string
  status: 'pending' | 'executing' | 'completed' | 'failed' | 'cancelled'
  progress: number
  filledQuantity: number
  avgPrice: number
  slippage: number
  marketImpact: number
  startTime: number
  endTime?: number
  executions: Array<{
    venue: string
    quantity: number
    price: number
    timestamp: number
    fees: number
  }>
}

interface SmartOrderRouter {
  isEnabled: boolean
  routingStrategy: 'best_price' | 'best_execution' | 'minimal_impact' | 'fastest_fill'
  venueWeights: Record<string, number>
  latencyThreshold: number
  liquidityThreshold: number
  maxVenueCount: number
}

export function OrderExecutionEngine() {
  const [venues, setVenues] = useState<ExecutionVenue[]>([])
  const [algorithms, setAlgorithms] = useState<ExecutionAlgorithm[]>([])
  const [executions, setExecutions] = useState<OrderExecution[]>([])
  const [smartRouter, setSmartRouter] = useState<SmartOrderRouter | null>(null)
  const [activeTab, setActiveTab] = useState('executions')
  const [selectedExecution, setSelectedExecution] = useState<string | null>(null)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock execution venues
    const mockVenues: ExecutionVenue[] = [
      {
        id: 'binance',
        name: 'Binance',
        type: 'cex',
        latency: 12,
        fees: 0.1,
        liquidity: 95,
        reliability: 99.9,
        isConnected: true,
        lastPing: Date.now() - 1000,
        supportedPairs: ['BTC/USDT', 'ETH/USDT', 'BNB/USDT'],
        minOrderSize: 10,
        maxOrderSize: 1000000
      },
      {
        id: 'coinbase',
        name: 'Coinbase Pro',
        type: 'cex',
        latency: 18,
        fees: 0.15,
        liquidity: 88,
        reliability: 99.8,
        isConnected: true,
        lastPing: Date.now() - 2000,
        supportedPairs: ['BTC/USDT', 'ETH/USDT'],
        minOrderSize: 5,
        maxOrderSize: 500000
      },
      {
        id: 'uniswap',
        name: 'Uniswap V3',
        type: 'dex',
        latency: 45,
        fees: 0.3,
        liquidity: 75,
        reliability: 98.5,
        isConnected: true,
        lastPing: Date.now() - 3000,
        supportedPairs: ['ETH/USDT', 'WBTC/ETH'],
        minOrderSize: 1,
        maxOrderSize: 100000
      },
      {
        id: 'darkpool1',
        name: 'Institutional Dark Pool',
        type: 'dark_pool',
        latency: 25,
        fees: 0.05,
        liquidity: 60,
        reliability: 99.5,
        isConnected: true,
        lastPing: Date.now() - 1500,
        supportedPairs: ['BTC/USDT', 'ETH/USDT'],
        minOrderSize: 1000,
        maxOrderSize: 10000000
      }
    ]

    const mockAlgorithms: ExecutionAlgorithm[] = [
      {
        id: 'twap',
        name: 'Time-Weighted Average Price',
        description: 'Executes orders evenly over a specified time period',
        type: 'twap',
        parameters: {
          duration: 3600,
          intervals: 12,
          randomization: 0.1
        },
        isActive: true,
        performance: {
          avgSlippage: 0.08,
          fillRate: 98.5,
          marketImpact: 0.12,
          executionTime: 3580
        }
      },
      {
        id: 'vwap',
        name: 'Volume-Weighted Average Price',
        description: 'Matches historical volume patterns for execution',
        type: 'vwap',
        parameters: {
          lookbackPeriod: 20,
          volumeParticipation: 0.15,
          aggressiveness: 'medium'
        },
        isActive: true,
        performance: {
          avgSlippage: 0.06,
          fillRate: 97.8,
          marketImpact: 0.09,
          executionTime: 2840
        }
      },
      {
        id: 'iceberg',
        name: 'Iceberg Algorithm',
        description: 'Hides large orders by showing small visible quantities',
        type: 'iceberg',
        parameters: {
          visibleSize: 0.1,
          randomization: true,
          priceImprovement: 0.05
        },
        isActive: true,
        performance: {
          avgSlippage: 0.04,
          fillRate: 96.2,
          marketImpact: 0.06,
          executionTime: 4200
        }
      },
      {
        id: 'sniper',
        name: 'Sniper Algorithm',
        description: 'Aggressive execution for immediate fills',
        type: 'sniper',
        parameters: {
          maxSlippage: 0.2,
          timeLimit: 300,
          venueCount: 3
        },
        isActive: true,
        performance: {
          avgSlippage: 0.15,
          fillRate: 99.8,
          marketImpact: 0.18,
          executionTime: 45
        }
      }
    ]

    const mockExecutions: OrderExecution[] = [
      {
        id: 'exec1',
        orderId: 'order123',
        symbol: 'BTC/USDT',
        side: 'buy',
        quantity: 2.5,
        algorithm: 'vwap',
        venue: 'binance',
        status: 'executing',
        progress: 65,
        filledQuantity: 1.625,
        avgPrice: 44850.25,
        slippage: 0.08,
        marketImpact: 0.12,
        startTime: Date.now() - 1800000,
        executions: [
          { venue: 'binance', quantity: 0.8, price: 44820.50, timestamp: Date.now() - 1700000, fees: 3.58 },
          { venue: 'coinbase', quantity: 0.5, price: 44865.75, timestamp: Date.now() - 1500000, fees: 3.36 },
          { venue: 'binance', quantity: 0.325, price: 44890.10, timestamp: Date.now() - 1200000, fees: 1.46 }
        ]
      },
      {
        id: 'exec2',
        orderId: 'order124',
        symbol: 'ETH/USDT',
        side: 'sell',
        quantity: 10,
        algorithm: 'twap',
        venue: 'multi',
        status: 'completed',
        progress: 100,
        filledQuantity: 10,
        avgPrice: 2485.60,
        slippage: 0.05,
        marketImpact: 0.08,
        startTime: Date.now() - 3600000,
        endTime: Date.now() - 600000,
        executions: [
          { venue: 'binance', quantity: 4.2, price: 2487.30, timestamp: Date.now() - 3000000, fees: 10.45 },
          { venue: 'coinbase', quantity: 3.1, price: 2484.85, timestamp: Date.now() - 2400000, fees: 11.56 },
          { venue: 'uniswap', quantity: 2.7, price: 2484.10, timestamp: Date.now() - 1800000, fees: 20.15 }
        ]
      }
    ]

    const mockSmartRouter: SmartOrderRouter = {
      isEnabled: true,
      routingStrategy: 'best_execution',
      venueWeights: {
        binance: 0.4,
        coinbase: 0.3,
        uniswap: 0.2,
        darkpool1: 0.1
      },
      latencyThreshold: 50,
      liquidityThreshold: 70,
      maxVenueCount: 3
    }

    setVenues(mockVenues)
    setAlgorithms(mockAlgorithms)
    setExecutions(mockExecutions)
    setSmartRouter(mockSmartRouter)
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
      case 'executing': return 'text-blue-500'
      case 'completed': return 'text-green-500'
      case 'failed': case 'cancelled': return 'text-red-500'
      case 'pending': return 'text-yellow-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'executing': return 'secondary'
      case 'completed': return 'default'
      case 'failed': case 'cancelled': return 'destructive'
      case 'pending': return 'outline'
      default: return 'outline'
    }
  }

  const getVenueTypeIcon = (type: string) => {
    switch (type) {
      case 'cex': return <Server className="w-4 h-4" />
      case 'dex': return <Globe className="w-4 h-4" />
      case 'dark_pool': return <Eye className="w-4 h-4" />
      case 'otc': return <Shield className="w-4 h-4" />
      default: return <Database className="w-4 h-4" />
    }
  }

  const getLatencyColor = (latency: number) => {
    if (latency <= 20) return 'text-green-500'
    if (latency <= 50) return 'text-yellow-500'
    return 'text-red-500'
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Zap className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access the order execution engine
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
          <h2 className="text-2xl font-bold">Order Execution Engine</h2>
          <p className="text-muted-foreground">
            Advanced order execution with smart routing and algorithmic strategies
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Zap className="w-3 h-3 mr-1" />
            {executions.filter(e => e.status === 'executing').length} Active
          </Badge>
          <Badge variant="outline">
            <Server className="w-3 h-3 mr-1" />
            {venues.filter(v => v.isConnected).length}/{venues.length} Connected
          </Badge>
        </div>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="executions">Active Executions</TabsTrigger>
          <TabsTrigger value="venues">Execution Venues</TabsTrigger>
          <TabsTrigger value="algorithms">Algorithms</TabsTrigger>
          <TabsTrigger value="routing">Smart Routing</TabsTrigger>
        </TabsList>

        <TabsContent value="executions" className="space-y-4">
          <div className="grid grid-cols-1 gap-4">
            {executions.map((execution) => (
              <Card key={execution.id}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Activity className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{execution.symbol} • {execution.side.toUpperCase()}</h4>
                        <p className="text-sm text-muted-foreground">
                          {execution.quantity} @ {execution.algorithm.toUpperCase()}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={getStatusBadgeVariant(execution.status)}>
                        {execution.status}
                      </Badge>
                      <Badge variant="outline">
                        {execution.venue === 'multi' ? 'Multi-Venue' : execution.venue}
                      </Badge>
                    </div>
                  </div>

                  {/* Progress */}
                  {execution.status === 'executing' && (
                    <div className="mb-4">
                      <div className="flex justify-between text-sm mb-1">
                        <span>Execution Progress</span>
                        <span>{execution.progress}%</span>
                      </div>
                      <Progress value={execution.progress} className="h-2" />
                    </div>
                  )}

                  {/* Execution Metrics */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Filled Quantity</div>
                      <div className="font-medium">{execution.filledQuantity.toFixed(3)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Avg Price</div>
                      <div className="font-medium">{formatCurrency(execution.avgPrice)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Slippage</div>
                      <div className={cn(
                        "font-medium",
                        execution.slippage <= 0.1 ? "text-green-500" : 
                        execution.slippage <= 0.2 ? "text-yellow-500" : "text-red-500"
                      )}>
                        {(execution.slippage * 100).toFixed(2)}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Market Impact</div>
                      <div className={cn(
                        "font-medium",
                        execution.marketImpact <= 0.1 ? "text-green-500" : 
                        execution.marketImpact <= 0.2 ? "text-yellow-500" : "text-red-500"
                      )}>
                        {(execution.marketImpact * 100).toFixed(2)}%
                      </div>
                    </div>
                  </div>

                  {/* Execution Details */}
                  <div className="space-y-2 mb-4">
                    <div className="text-sm font-medium">Execution Breakdown</div>
                    <div className="space-y-1">
                      {execution.executions.map((exec, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-muted rounded text-xs">
                          <div className="flex items-center gap-2">
                            {getVenueTypeIcon('cex')}
                            <span className="font-medium">{exec.venue}</span>
                          </div>
                          <div className="flex items-center gap-4">
                            <span>{exec.quantity} @ {formatCurrency(exec.price)}</span>
                            <span className="text-muted-foreground">
                              {formatTime(exec.timestamp)}
                            </span>
                            <span className="text-muted-foreground">
                              Fee: {formatCurrency(exec.fees)}
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    {execution.status === 'executing' && (
                      <>
                        <Button variant="outline" size="sm">
                          <Pause className="w-3 h-3 mr-1" />
                          Pause
                        </Button>
                        <Button variant="outline" size="sm">
                          <Square className="w-3 h-3 mr-1" />
                          Cancel
                        </Button>
                      </>
                    )}
                    <Button variant="outline" size="sm">
                      <Eye className="w-3 h-3 mr-1" />
                      Details
                    </Button>
                    <Button variant="outline" size="sm">
                      <Download className="w-3 h-3 mr-1" />
                      Export
                    </Button>
                  </div>

                  {/* Footer */}
                  <div className="flex items-center justify-between mt-4 pt-4 border-t text-xs text-muted-foreground">
                    <span>Started: {formatTime(execution.startTime)}</span>
                    {execution.endTime && (
                      <span>Completed: {formatTime(execution.endTime)}</span>
                    )}
                    <span>Order ID: {execution.orderId}</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {executions.length === 0 && (
            <div className="text-center py-12">
              <Activity className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="text-lg font-medium mb-2">No Active Executions</h3>
              <p className="text-muted-foreground">
                Order executions will appear here when you place advanced orders
              </p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="venues" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {venues.map((venue) => (
              <Card key={venue.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        {getVenueTypeIcon(venue.type)}
                      </div>
                      <div>
                        <h4 className="font-bold">{venue.name}</h4>
                        <p className="text-sm text-muted-foreground capitalize">{venue.type.replace('_', ' ')}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-1">
                      {venue.isConnected ? (
                        <Wifi className="w-4 h-4 text-green-500" />
                      ) : (
                        <WifiOff className="w-4 h-4 text-red-500" />
                      )}
                    </div>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  {/* Performance Metrics */}
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Latency</div>
                      <div className={cn("font-medium", getLatencyColor(venue.latency))}>
                        {venue.latency}ms
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Fees</div>
                      <div className="font-medium">{venue.fees}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Liquidity</div>
                      <div className="font-medium">{venue.liquidity}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Reliability</div>
                      <div className="font-medium text-green-500">{venue.reliability}%</div>
                    </div>
                  </div>

                  {/* Order Limits */}
                  <div className="space-y-2">
                    <div className="text-sm font-medium">Order Limits</div>
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">Min:</span>
                      <span className="font-medium">{formatCurrency(venue.minOrderSize)}</span>
                    </div>
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">Max:</span>
                      <span className="font-medium">{formatCurrency(venue.maxOrderSize)}</span>
                    </div>
                  </div>

                  {/* Supported Pairs */}
                  <div className="space-y-2">
                    <div className="text-sm font-medium">Supported Pairs</div>
                    <div className="flex flex-wrap gap-1">
                      {venue.supportedPairs.slice(0, 3).map((pair) => (
                        <Badge key={pair} variant="outline" className="text-xs">
                          {pair}
                        </Badge>
                      ))}
                      {venue.supportedPairs.length > 3 && (
                        <Badge variant="outline" className="text-xs">
                          +{venue.supportedPairs.length - 3} more
                        </Badge>
                      )}
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2">
                    <Button 
                      variant={venue.isConnected ? "outline" : "default"} 
                      size="sm"
                      className="flex-1"
                    >
                      {venue.isConnected ? (
                        <>
                          <WifiOff className="w-3 h-3 mr-1" />
                          Disconnect
                        </>
                      ) : (
                        <>
                          <Wifi className="w-3 h-3 mr-1" />
                          Connect
                        </>
                      )}
                    </Button>
                    <Button variant="outline" size="sm">
                      <Settings className="w-3 h-3" />
                    </Button>
                  </div>

                  {/* Connection Status */}
                  <div className="text-xs text-muted-foreground">
                    Last ping: {Math.floor((Date.now() - venue.lastPing) / 1000)}s ago
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="algorithms" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {algorithms.map((algorithm) => (
              <Card key={algorithm.id}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="font-bold">{algorithm.name}</h4>
                      <p className="text-sm text-muted-foreground">{algorithm.description}</p>
                    </div>
                    <Badge variant={algorithm.isActive ? 'default' : 'outline'}>
                      {algorithm.isActive ? 'Active' : 'Inactive'}
                    </Badge>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  {/* Performance Metrics */}
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Avg Slippage</div>
                      <div className="font-medium text-green-500">
                        {(algorithm.performance.avgSlippage * 100).toFixed(2)}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Fill Rate</div>
                      <div className="font-medium">{algorithm.performance.fillRate.toFixed(1)}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Market Impact</div>
                      <div className="font-medium">
                        {(algorithm.performance.marketImpact * 100).toFixed(2)}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Avg Execution</div>
                      <div className="font-medium">{algorithm.performance.executionTime}s</div>
                    </div>
                  </div>

                  {/* Parameters */}
                  <div className="space-y-2">
                    <div className="text-sm font-medium">Parameters</div>
                    <div className="space-y-1">
                      {Object.entries(algorithm.parameters).map(([key, value]) => (
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
                      variant={algorithm.isActive ? "outline" : "default"} 
                      size="sm"
                      className="flex-1"
                    >
                      {algorithm.isActive ? (
                        <>
                          <Pause className="w-3 h-3 mr-1" />
                          Disable
                        </>
                      ) : (
                        <>
                          <Play className="w-3 h-3 mr-1" />
                          Enable
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
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="routing" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Brain className="w-5 h-5" />
                Smart Order Router
              </CardTitle>
            </CardHeader>
            <CardContent>
              {smartRouter && (
                <div className="space-y-6">
                  <div className="flex items-center space-x-2">
                    <Switch 
                      id="enableRouter" 
                      checked={smartRouter.isEnabled}
                      onCheckedChange={(checked) => setSmartRouter(prev => prev ? 
                        { ...prev, isEnabled: checked } : null
                      )}
                    />
                    <Label htmlFor="enableRouter">Enable Smart Order Routing</Label>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="routingStrategy">Routing Strategy</Label>
                      <Select 
                        value={smartRouter.routingStrategy}
                        onValueChange={(value) => setSmartRouter(prev => prev ? 
                          { ...prev, routingStrategy: value as any } : null
                        )}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="best_price">Best Price</SelectItem>
                          <SelectItem value="best_execution">Best Execution</SelectItem>
                          <SelectItem value="minimal_impact">Minimal Impact</SelectItem>
                          <SelectItem value="fastest_fill">Fastest Fill</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label htmlFor="maxVenues">Max Venues</Label>
                      <Input
                        id="maxVenues"
                        type="number"
                        value={smartRouter.maxVenueCount}
                        onChange={(e) => setSmartRouter(prev => prev ? 
                          { ...prev, maxVenueCount: parseInt(e.target.value) } : null
                        )}
                        min="1"
                        max="10"
                      />
                    </div>

                    <div>
                      <Label htmlFor="latencyThreshold">Latency Threshold (ms)</Label>
                      <Input
                        id="latencyThreshold"
                        type="number"
                        value={smartRouter.latencyThreshold}
                        onChange={(e) => setSmartRouter(prev => prev ? 
                          { ...prev, latencyThreshold: parseInt(e.target.value) } : null
                        )}
                      />
                    </div>

                    <div>
                      <Label htmlFor="liquidityThreshold">Liquidity Threshold (%)</Label>
                      <Input
                        id="liquidityThreshold"
                        type="number"
                        value={smartRouter.liquidityThreshold}
                        onChange={(e) => setSmartRouter(prev => prev ? 
                          { ...prev, liquidityThreshold: parseInt(e.target.value) } : null
                        )}
                        min="0"
                        max="100"
                      />
                    </div>
                  </div>

                  {/* Venue Weights */}
                  <div className="space-y-4">
                    <h4 className="font-medium">Venue Weights</h4>
                    {Object.entries(smartRouter.venueWeights).map(([venue, weight]) => (
                      <div key={venue} className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="capitalize">{venue}</span>
                          <span>{(weight * 100).toFixed(0)}%</span>
                        </div>
                        <Progress value={weight * 100} className="h-2" />
                      </div>
                    ))}
                  </div>

                  <div className="flex gap-3">
                    <Button>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      Save Configuration
                    </Button>
                    <Button variant="outline">
                      <RefreshCw className="w-4 h-4 mr-2" />
                      Reset to Defaults
                    </Button>
                    <Button variant="outline">
                      <BarChart3 className="w-4 h-4 mr-2" />
                      View Analytics
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
