'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  TrendingUp,
  TrendingDown,
  Target,
  Bell,
  BellOff,
  Zap,
  Activity,
  BarChart3,
  Settings,
  RefreshCw,
  Download,
  Plus,
  AlertCircle,
  CheckCircle,
  Clock,
  Filter
} from 'lucide-react'
import { useTradingSignals } from '@/hooks/useTradingSignals'
import { type TradingSignal } from '@/lib/trading-signals'
import { cn } from '@/lib/utils'

export function TradingSignalsDashboard() {
  const [activeTab, setActiveTab] = useState('signals')
  const [selectedSignalType, setSelectedSignalType] = useState<'all' | 'buy' | 'sell' | 'hold'>('all')
  const [selectedTimeframe, setSelectedTimeframe] = useState<'all' | '1h' | '4h' | '1d'>('all')

  const {
    state,
    getSignals,
    getSignalsByType,
    getOverallPerformance,
    refresh,
    exportSignals,
    clearError
  } = useTradingSignals({
    symbols: ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'],
    autoStart: true,
    enableAlerts: true
  })

  const performance = getOverallPerformance()
  
  const filteredSignals = state.signals.filter(signal => {
    if (selectedSignalType !== 'all' && signal.type !== selectedSignalType) return false
    if (selectedTimeframe !== 'all' && signal.timeframe !== selectedTimeframe) return false
    return true
  })

  const getSignalColor = (type: string) => {
    switch (type) {
      case 'buy':
        return 'text-green-600 dark:text-green-400'
      case 'sell':
        return 'text-red-600 dark:text-red-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const getSignalIcon = (type: string) => {
    switch (type) {
      case 'buy':
        return <TrendingUp className="w-4 h-4" />
      case 'sell':
        return <TrendingDown className="w-4 h-4" />
      default:
        return <Activity className="w-4 h-4" />
    }
  }

  const getSignalBadge = (signal: TradingSignal) => {
    const variant = signal.type === 'buy' ? 'default' : 
                   signal.type === 'sell' ? 'destructive' : 'secondary'
    
    return (
      <Badge variant={variant} className="flex items-center gap-1">
        {getSignalIcon(signal.type)}
        {signal.type.toUpperCase()}
      </Badge>
    )
  }

  const getStrengthColor = (strength: number) => {
    if (strength >= 80) return 'text-green-600 dark:text-green-400'
    if (strength >= 60) return 'text-yellow-600 dark:text-yellow-400'
    return 'text-red-600 dark:text-red-400'
  }

  const formatTimeAgo = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days}d ago`
    if (hours > 0) return `${hours}h ago`
    if (minutes > 0) return `${minutes}m ago`
    return 'Just now'
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 6
    }).format(price)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Trading Signals</h1>
          <p className="text-muted-foreground">
            AI-powered trading signals and automated alerts
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={exportSignals}
          >
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={refresh}
            disabled={state.isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Signals</p>
                <p className="text-2xl font-bold">{performance.totalSignals}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Last 24 hours
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {performance.successRate.toFixed(1)}%
                </p>
              </div>
              <Target className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Average performance
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Strategies</p>
                <p className="text-2xl font-bold">{performance.activeStrategies}</p>
              </div>
              <BarChart3 className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Currently running
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Alerts</p>
                <p className="text-2xl font-bold">{performance.activeAlerts}</p>
              </div>
              <Bell className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Monitoring signals
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Signal Type Summary */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Buy Signals</p>
                <p className="text-xl font-bold text-green-600">
                  {getSignalsByType('buy').length}
                </p>
              </div>
              <TrendingUp className="w-6 h-6 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Sell Signals</p>
                <p className="text-xl font-bold text-red-600">
                  {getSignalsByType('sell').length}
                </p>
              </div>
              <TrendingDown className="w-6 h-6 text-red-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Hold Signals</p>
                <p className="text-xl font-bold text-gray-600">
                  {getSignalsByType('hold').length}
                </p>
              </div>
              <Activity className="w-6 h-6 text-gray-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="signals">Signals</TabsTrigger>
          <TabsTrigger value="strategies">Strategies</TabsTrigger>
          <TabsTrigger value="alerts">Alerts</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>

        <TabsContent value="signals" className="space-y-6">
          {/* Filters */}
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Filter className="w-4 h-4" />
              <span className="text-sm font-medium">Filters:</span>
            </div>
            
            <select
              value={selectedSignalType}
              onChange={(e) => setSelectedSignalType(e.target.value as any)}
              className="px-3 py-1 border rounded-md text-sm"
            >
              <option value="all">All Types</option>
              <option value="buy">Buy Signals</option>
              <option value="sell">Sell Signals</option>
              <option value="hold">Hold Signals</option>
            </select>

            <select
              value={selectedTimeframe}
              onChange={(e) => setSelectedTimeframe(e.target.value as any)}
              className="px-3 py-1 border rounded-md text-sm"
            >
              <option value="all">All Timeframes</option>
              <option value="1h">1 Hour</option>
              <option value="4h">4 Hours</option>
              <option value="1d">1 Day</option>
            </select>
          </div>

          {/* Signals List */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Trading Signals</CardTitle>
              <CardDescription>
                Latest AI-generated trading signals with confidence scores
              </CardDescription>
            </CardHeader>
            <CardContent>
              {filteredSignals.length > 0 ? (
                <div className="space-y-3">
                  <AnimatePresence>
                    {filteredSignals.slice(0, 20).map((signal, index) => (
                      <motion.div
                        key={signal.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        transition={{ delay: index * 0.05 }}
                        className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50"
                      >
                        <div className="flex items-center gap-4">
                          <div className="w-12 h-12 bg-secondary rounded-lg flex items-center justify-center">
                            <span className="font-bold text-lg">{signal.symbol.slice(0, 2)}</span>
                          </div>
                          
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <h4 className="font-medium">{signal.symbol}</h4>
                              {getSignalBadge(signal)}
                              <Badge variant="outline" className="text-xs">
                                {signal.timeframe}
                              </Badge>
                            </div>
                            <p className="text-sm text-muted-foreground">
                              {signal.description}
                            </p>
                            <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
                              <span>Price: {formatPrice(signal.price)}</span>
                              {signal.targetPrice && (
                                <span>Target: {formatPrice(signal.targetPrice)}</span>
                              )}
                              {signal.stopLoss && (
                                <span>Stop: {formatPrice(signal.stopLoss)}</span>
                              )}
                            </div>
                          </div>
                        </div>

                        <div className="text-right">
                          <div className="flex items-center gap-2 mb-1">
                            <span className="text-sm font-medium">Strength:</span>
                            <span className={cn("text-sm font-bold", getStrengthColor(signal.strength))}>
                              {signal.strength.toFixed(0)}%
                            </span>
                          </div>
                          <div className="flex items-center gap-2 mb-1">
                            <span className="text-sm font-medium">Confidence:</span>
                            <span className="text-sm font-bold">
                              {signal.confidence.toFixed(0)}%
                            </span>
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {formatTimeAgo(signal.timestamp)}
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </AnimatePresence>
                </div>
              ) : (
                <div className="text-center py-12">
                  <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Signals Found</h3>
                  <p className="text-muted-foreground">
                    {selectedSignalType !== 'all' || selectedTimeframe !== 'all' 
                      ? 'No signals match your current filters'
                      : 'Waiting for trading signals to be generated'
                    }
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="strategies" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span>Trading Strategies</span>
                <Button size="sm">
                  <Plus className="w-4 h-4 mr-2" />
                  Add Strategy
                </Button>
              </CardTitle>
              <CardDescription>
                Manage your automated trading signal strategies
              </CardDescription>
            </CardHeader>
            <CardContent>
              {state.strategies.length > 0 ? (
                <div className="space-y-4">
                  {state.strategies.map((strategy, index) => (
                    <motion.div
                      key={strategy.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-4 border rounded-lg"
                    >
                      <div className="flex items-center gap-4">
                        <div className={cn("w-3 h-3 rounded-full", 
                          strategy.enabled ? 'bg-green-500' : 'bg-gray-400'
                        )} />
                        <div>
                          <h4 className="font-medium">{strategy.name}</h4>
                          <p className="text-sm text-muted-foreground">{strategy.description}</p>
                          <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
                            <span>Risk: {strategy.riskLevel}</span>
                            <span>Min Confidence: {strategy.minConfidence}%</span>
                            <span>Symbols: {strategy.symbols.length}</span>
                          </div>
                        </div>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <Badge variant={strategy.enabled ? 'default' : 'secondary'}>
                          {strategy.enabled ? 'Active' : 'Inactive'}
                        </Badge>
                        <Button variant="ghost" size="sm">
                          <Settings className="w-4 h-4" />
                        </Button>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Strategies</h3>
                  <p className="text-muted-foreground mb-4">
                    Create your first trading strategy to start generating signals
                  </p>
                  <Button>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Strategy
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="alerts" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span>Signal Alerts</span>
                <Button size="sm">
                  <Plus className="w-4 h-4 mr-2" />
                  Add Alert
                </Button>
              </CardTitle>
              <CardDescription>
                Configure notifications for trading signals
              </CardDescription>
            </CardHeader>
            <CardContent>
              {state.alerts.length > 0 ? (
                <div className="space-y-4">
                  {state.alerts.map((alert, index) => (
                    <motion.div
                      key={alert.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-4 border rounded-lg"
                    >
                      <div className="flex items-center gap-4">
                        <div className={cn("w-3 h-3 rounded-full", 
                          alert.enabled ? 'bg-green-500' : 'bg-gray-400'
                        )} />
                        <div>
                          <h4 className="font-medium">Signal Alert</h4>
                          <p className="text-sm text-muted-foreground">
                            {alert.type} notifications for trading signals
                          </p>
                          <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
                            <span>Type: {alert.type}</span>
                            <span>Cooldown: {alert.cooldownPeriod}m</span>
                            <span>Max/day: {alert.maxAlertsPerDay}</span>
                          </div>
                        </div>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <Badge variant={alert.enabled ? 'default' : 'secondary'}>
                          {alert.enabled ? 'Active' : 'Inactive'}
                        </Badge>
                        <Button variant="ghost" size="sm">
                          <Settings className="w-4 h-4" />
                        </Button>
                      </div>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <Bell className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Alerts</h3>
                  <p className="text-muted-foreground mb-4">
                    Set up alerts to be notified when trading signals are generated
                  </p>
                  <Button>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Alert
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="performance" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Strategy Performance</CardTitle>
              <CardDescription>
                Track the performance of your trading strategies
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold text-green-600">
                    {performance.successRate.toFixed(1)}%
                  </p>
                  <p className="text-sm text-muted-foreground">Success Rate</p>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">
                    {performance.avgReturn > 0 ? '+' : ''}{performance.avgReturn.toFixed(2)}%
                  </p>
                  <p className="text-sm text-muted-foreground">Avg Return</p>
                </div>
                <div className="text-center p-4 border rounded-lg">
                  <p className="text-2xl font-bold">{performance.totalSignals}</p>
                  <p className="text-sm text-muted-foreground">Total Signals</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
