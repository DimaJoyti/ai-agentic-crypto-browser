'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { 
  ArrowRightLeft,
  TrendingUp,
  DollarSign,
  Shield,
  AlertTriangle,
  CheckCircle,
  Clock,
  Zap,
  Activity,
  BarChart3,
  RefreshCw,
  Settings,
  ArrowRight,
  ExternalLink,
  Info,
  Globe,
  Link
} from 'lucide-react'
import { useCrossChain, useBridge, useCrossChainPortfolio, useBridgeAnalytics } from '@/hooks/useCrossChain'
import { type BridgeRoute, type BridgeStatus } from '@/lib/cross-chain-bridge'
import { cn } from '@/lib/utils'

export function CrossChainDashboard() {
  const [activeTab, setActiveTab] = useState('bridge')
  const [fromChain, setFromChain] = useState<number>(1) // Ethereum
  const [toChain, setToChain] = useState<number>(137) // Polygon
  const [selectedToken, setSelectedToken] = useState('USDC')
  const [amount, setAmount] = useState('')
  const [selectedRoute, setSelectedRoute] = useState<BridgeRoute | null>(null)

  const {
    state,
    getBridgeRoutes,
    executeBridge,
    refreshPositions,
    clearError
  } = useCrossChain({
    autoLoad: true,
    enableNotifications: true,
    autoRefresh: true
  })

  const { portfolio } = useCrossChainPortfolio()
  const analytics = useBridgeAnalytics()

  const chains = [
    { id: 1, name: 'Ethereum', symbol: 'ETH', logo: '⟠' },
    { id: 137, name: 'Polygon', symbol: 'MATIC', logo: '⬟' },
    { id: 56, name: 'BSC', symbol: 'BNB', logo: '🟡' },
    { id: 43114, name: 'Avalanche', symbol: 'AVAX', logo: '🔺' },
    { id: 42161, name: 'Arbitrum', symbol: 'ARB', logo: '🔵' },
    { id: 10, name: 'Optimism', symbol: 'OP', logo: '🔴' }
  ]

  const tokens = ['USDC', 'USDT', 'ETH', 'WBTC']

  const handleGetRoutes = async () => {
    if (!amount) return

    try {
      const routes = await getBridgeRoutes(fromChain, toChain, selectedToken, amount)
      if (routes.length > 0) {
        setSelectedRoute(routes[0])
      }
    } catch (error) {
      console.error('Failed to get routes:', error)
    }
  }

  const handleSwapChains = () => {
    const temp = fromChain
    setFromChain(toChain)
    setToChain(temp)
    setSelectedRoute(null)
  }

  const handleExecuteBridge = async () => {
    if (!selectedRoute) return

    try {
      await executeBridge(selectedRoute.id, amount)
      setAmount('')
      setSelectedRoute(null)
    } catch (error) {
      console.error('Bridge failed:', error)
    }
  }

  const formatNumber = (value: string | number, decimals: number = 2) => {
    return parseFloat(value.toString()).toFixed(decimals)
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value)
  }

  const getChainName = (chainId: number) => {
    return chains.find(c => c.id === chainId)?.name || `Chain ${chainId}`
  }

  const getChainLogo = (chainId: number) => {
    return chains.find(c => c.id === chainId)?.logo || '🔗'
  }

  const getStatusColor = (status: BridgeStatus) => {
    switch (status) {
      case 'completed':
        return 'text-green-600 bg-green-100 dark:bg-green-900'
      case 'pending':
      case 'confirmed_source':
      case 'in_transit':
        return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900'
      case 'failed':
        return 'text-red-600 bg-red-100 dark:bg-red-900'
      default:
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900'
    }
  }

  const getProtocolLogo = (protocolId: string) => {
    const logos: Record<string, string> = {
      'stargate': '⭐',
      'hop': '🐰',
      'across': '🌉'
    }
    return logos[protocolId] || '🌉'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Cross-Chain Bridge</h2>
          <p className="text-muted-foreground">
            Bridge assets across multiple blockchains securely
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={refreshPositions}>
            <RefreshCw className="w-4 h-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" size="sm">
            <BarChart3 className="w-4 h-4 mr-2" />
            Analytics
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Portfolio Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Portfolio</p>
                <p className="text-2xl font-bold">{formatCurrency(portfolio.totalValue)}</p>
              </div>
              <Globe className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Across {Object.keys(portfolio.chainDistribution).length} chains
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Bridge Volume</p>
                <p className="text-2xl font-bold">{formatCurrency(analytics.totalVolume)}</p>
              </div>
              <ArrowRightLeft className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {analytics.totalTransactions} transactions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {analytics.successRate.toFixed(1)}%
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Bridge reliability
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Fees</p>
                <p className="text-2xl font-bold">${analytics.averageFees.toFixed(3)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Per transaction
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="bridge">Bridge</TabsTrigger>
          <TabsTrigger value="portfolio">Portfolio ({portfolio.positions.length})</TabsTrigger>
          <TabsTrigger value="history">History ({state.transactions.length})</TabsTrigger>
          <TabsTrigger value="protocols">Protocols ({state.protocols.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="bridge" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Bridge Interface */}
            <Card>
              <CardHeader>
                <CardTitle>Bridge Assets</CardTitle>
                <CardDescription>
                  Transfer tokens across different blockchains
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* From Chain */}
                <div className="space-y-2">
                  <Label>From Chain</Label>
                  <div className="flex gap-2 flex-wrap">
                    {chains.map((chain) => (
                      <Button
                        key={chain.id}
                        variant={fromChain === chain.id ? "default" : "outline"}
                        size="sm"
                        onClick={() => setFromChain(chain.id)}
                      >
                        <span className="mr-2">{chain.logo}</span>
                        {chain.name}
                      </Button>
                    ))}
                  </div>
                </div>

                {/* Swap Button */}
                <div className="flex justify-center">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleSwapChains}
                    className="rounded-full"
                  >
                    <ArrowRightLeft className="w-4 h-4" />
                  </Button>
                </div>

                {/* To Chain */}
                <div className="space-y-2">
                  <Label>To Chain</Label>
                  <div className="flex gap-2 flex-wrap">
                    {chains.map((chain) => (
                      <Button
                        key={chain.id}
                        variant={toChain === chain.id ? "default" : "outline"}
                        size="sm"
                        onClick={() => setToChain(chain.id)}
                        disabled={chain.id === fromChain}
                      >
                        <span className="mr-2">{chain.logo}</span>
                        {chain.name}
                      </Button>
                    ))}
                  </div>
                </div>

                {/* Token Selection */}
                <div className="space-y-2">
                  <Label>Token</Label>
                  <div className="flex gap-2">
                    {tokens.map((token) => (
                      <Button
                        key={token}
                        variant={selectedToken === token ? "default" : "outline"}
                        size="sm"
                        onClick={() => setSelectedToken(token)}
                      >
                        {token}
                      </Button>
                    ))}
                  </div>
                </div>

                {/* Amount Input */}
                <div className="space-y-2">
                  <Label>Amount</Label>
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder="0.0"
                      value={amount}
                      onChange={(e) => setAmount(e.target.value)}
                    />
                    <Button variant="outline" size="sm">
                      Max
                    </Button>
                  </div>
                  {amount && (
                    <p className="text-sm text-muted-foreground">
                      ≈ ${formatNumber(parseFloat(amount) * 1800, 2)} {/* Mock USD value */}
                    </p>
                  )}
                </div>

                {/* Action Buttons */}
                <div className="space-y-2">
                  <Button
                    onClick={handleGetRoutes}
                    disabled={!amount || fromChain === toChain || state.isLoading}
                    className="w-full"
                  >
                    {state.isLoading ? (
                      <>
                        <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                        Finding Routes...
                      </>
                    ) : (
                      'Get Bridge Routes'
                    )}
                  </Button>

                  {selectedRoute && (
                    <Button
                      onClick={handleExecuteBridge}
                      disabled={state.isBridging}
                      className="w-full"
                      variant="default"
                    >
                      {state.isBridging ? (
                        <>
                          <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                          Bridging...
                        </>
                      ) : (
                        <>
                          <Zap className="w-4 h-4 mr-2" />
                          Execute Bridge
                        </>
                      )}
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Route Details */}
            <Card>
              <CardHeader>
                <CardTitle>Best Route</CardTitle>
                <CardDescription>
                  Optimal bridge route and pricing details
                </CardDescription>
              </CardHeader>
              <CardContent>
                {selectedRoute ? (
                  <div className="space-y-4">
                    {/* Protocol Info */}
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-2xl">{getProtocolLogo(selectedRoute.protocol.id)}</span>
                        <div>
                          <p className="font-medium">{selectedRoute.protocol.name}</p>
                          <p className="text-sm text-muted-foreground">
                            Confidence: {selectedRoute.confidence.toFixed(0)}%
                          </p>
                        </div>
                      </div>
                      <Badge variant="default">Best Route</Badge>
                    </div>

                    {/* Route Details */}
                    <div className="space-y-3">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Estimated Time</span>
                        <span>{Math.floor(selectedRoute.estimatedTime / 60)} minutes</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Bridge Fee</span>
                        <span>{selectedRoute.fees.percentageFee}% + ${selectedRoute.fees.baseFee}</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Gas Estimate</span>
                        <span>${formatNumber(parseFloat(selectedRoute.fees.gasEstimate) * 0.00002, 3)}</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Risk Level</span>
                        <span className={cn(
                          selectedRoute.riskLevel === 'low' ? 'text-green-600' : 
                          selectedRoute.riskLevel === 'medium' ? 'text-yellow-600' : 'text-red-600'
                        )}>
                          {selectedRoute.riskLevel.toUpperCase()}
                        </span>
                      </div>
                    </div>

                    {/* Route Visualization */}
                    <div className="space-y-2">
                      <Label>Bridge Route</Label>
                      <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                        <span className="font-medium">{getChainLogo(selectedRoute.fromChain)} {getChainName(selectedRoute.fromChain)}</span>
                        <ArrowRight className="w-4 h-4 text-muted-foreground" />
                        <span className="font-medium">{getChainLogo(selectedRoute.toChain)} {getChainName(selectedRoute.toChain)}</span>
                        <Badge variant="outline" className="ml-auto">
                          {selectedRoute.protocol.name}
                        </Badge>
                      </div>
                    </div>

                    {/* Steps */}
                    <div className="space-y-2">
                      <Label>Bridge Steps</Label>
                      <div className="space-y-2">
                        {selectedRoute.steps.map((step, index) => (
                          <div key={index} className="flex items-center gap-3 p-2 bg-muted/50 rounded">
                            <div className="w-6 h-6 rounded-full bg-primary text-primary-foreground text-xs flex items-center justify-center">
                              {step.order}
                            </div>
                            <div className="flex-1">
                              <p className="text-sm font-medium">{step.action}</p>
                              <p className="text-xs text-muted-foreground">{step.description}</p>
                            </div>
                            <div className="text-xs text-muted-foreground">
                              ~{Math.floor(step.estimatedTime / 60)}m
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Link className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                    <h3 className="text-lg font-medium mb-2">No Route Selected</h3>
                    <p className="text-muted-foreground">
                      Get routes to see bridge details and pricing
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Available Routes */}
          {state.routes.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Available Routes ({state.routes.length})</CardTitle>
                <CardDescription>
                  Compare bridge routes across different protocols
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {state.routes.map((route, index) => (
                    <motion.div
                      key={route.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className={cn(
                        "flex items-center justify-between p-4 border rounded-lg cursor-pointer transition-all",
                        selectedRoute?.id === route.id ? "border-primary bg-primary/5" : "hover:bg-muted/50"
                      )}
                      onClick={() => setSelectedRoute(route)}
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-xl">{getProtocolLogo(route.protocol.id)}</span>
                        <div>
                          <p className="font-medium">{route.protocol.name}</p>
                          <p className="text-sm text-muted-foreground">
                            {Math.floor(route.estimatedTime / 60)} min • {route.fees.percentageFee}% fee
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="font-medium">
                          {route.confidence.toFixed(0)}% confidence
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {route.riskLevel} risk
                        </p>
                      </div>

                      <div className="flex items-center gap-2">
                        <Badge variant={index === 0 ? "default" : "outline"}>
                          {index === 0 ? 'Best' : `#${index + 1}`}
                        </Badge>
                        {selectedRoute?.id === route.id && (
                          <CheckCircle className="w-4 h-4 text-primary" />
                        )}
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="portfolio" className="space-y-4">
          {portfolio.positions.length > 0 ? (
            <div className="space-y-4">
              {portfolio.positions.map((position, index) => (
                <motion.div
                  key={position.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <Card>
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <div>
                          <CardTitle className="text-lg">{position.protocol.toUpperCase()}</CardTitle>
                          <CardDescription>
                            Total Value: {formatCurrency(position.totalValueUSD)}
                          </CardDescription>
                        </div>
                        <div className="text-right">
                          <p className="text-lg font-bold text-green-600">
                            {position.totalYield.toFixed(2)}% Yield
                          </p>
                          <p className="text-sm text-muted-foreground">
                            Risk: {position.riskScore}/100
                          </p>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-3">
                        {position.positions.map((chainPos, idx) => (
                          <div key={idx} className="flex items-center justify-between p-3 bg-muted rounded-lg">
                            <div className="flex items-center gap-2">
                              <span className="text-lg">{getChainLogo(chainPos.chainId)}</span>
                              <div>
                                <p className="font-medium">{getChainName(chainPos.chainId)}</p>
                                <p className="text-sm text-muted-foreground">
                                  {chainPos.type} • {chainPos.tokens.length} tokens
                                </p>
                              </div>
                            </div>
                            <div className="text-right">
                              <p className="font-medium">{formatCurrency(chainPos.valueUSD)}</p>
                              <p className="text-sm text-green-600">{chainPos.apy.toFixed(2)}% APY</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Globe className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Cross-Chain Positions</h3>
                <p className="text-muted-foreground">
                  Your cross-chain DeFi positions will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          {state.transactions.length > 0 ? (
            <div className="space-y-3">
              {state.transactions.map((transaction, index) => (
                <motion.div
                  key={transaction.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <Card>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex items-center gap-2">
                            <Badge className={cn("text-xs", getStatusColor(transaction.status))}>
                              {transaction.status.toUpperCase()}
                            </Badge>
                          </div>
                          
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <span className="text-lg">{getChainLogo(transaction.fromChain)}</span>
                              <ArrowRight className="w-3 h-3 text-muted-foreground" />
                              <span className="text-lg">{getChainLogo(transaction.toChain)}</span>
                              <p className="font-medium">
                                {formatNumber(transaction.amount)} {transaction.token.symbol}
                              </p>
                            </div>
                            
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <span>{getChainName(transaction.fromChain)} → {getChainName(transaction.toChain)}</span>
                              <span>{formatCurrency(transaction.amountUSD)}</span>
                              <span>{new Date(transaction.createdAt).toLocaleTimeString()}</span>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          {transaction.fromTxHash && (
                            <Button variant="ghost" size="sm">
                              <ExternalLink className="w-4 h-4" />
                            </Button>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Bridge History</h3>
                <p className="text-muted-foreground">
                  Your bridge transactions will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="protocols" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {state.protocols.map((protocol, index) => (
              <motion.div
                key={protocol.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <Card>
                  <CardHeader>
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">{getProtocolLogo(protocol.id)}</span>
                      <div>
                        <CardTitle className="text-lg">{protocol.name}</CardTitle>
                        <CardDescription>{protocol.metadata.description}</CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Supported Chains</span>
                        <span>{protocol.supportedChains.length}</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Success Rate</span>
                        <span className="text-green-600">{protocol.metadata.successRate}%</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Avg Time</span>
                        <span>{Math.floor(protocol.metadata.averageTime / 60)} min</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">TVL</span>
                        <span>${formatNumber(parseFloat(protocol.security.tvlLocked) / 1000000, 0)}M</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
