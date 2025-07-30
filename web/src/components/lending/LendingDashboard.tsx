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
  TrendingUp,
  TrendingDown,
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
  Plus,
  Minus,
  ArrowUpDown,
  ExternalLink,
  Info
} from 'lucide-react'
import { useLending, useLendingPosition, useYieldFarming, useLiquidationMonitor, useLendingAnalytics } from '@/hooks/useLending'
import { type LendingAsset, type UserPosition } from '@/lib/lending-integration'
import { cn } from '@/lib/utils'

export function LendingDashboard() {
  const [activeTab, setActiveTab] = useState('positions')
  const [selectedProtocol, setSelectedProtocol] = useState<string>('')
  const [selectedAsset, setSelectedAsset] = useState<LendingAsset | null>(null)
  const [actionType, setActionType] = useState<'supply' | 'withdraw' | 'borrow' | 'repay'>('supply')
  const [amount, setAmount] = useState('')

  const {
    state,
    supply,
    withdraw,
    borrow,
    repay,
    refreshPositions,
    clearError
  } = useLending({
    autoLoad: true,
    enableNotifications: true,
    autoRefresh: true
  })

  const { opportunities: yieldOpportunities } = useYieldFarming()
  const { atRiskPositions, liquidationOpportunities } = useLiquidationMonitor()
  const analytics = useLendingAnalytics()

  useEffect(() => {
    if (state.protocols.length > 0 && !selectedProtocol) {
      setSelectedProtocol(state.protocols[0].id)
    }
  }, [state.protocols, selectedProtocol])

  const handleAction = async () => {
    if (!selectedAsset || !amount || !selectedProtocol) return

    try {
      switch (actionType) {
        case 'supply':
          await supply(selectedProtocol, selectedAsset, amount)
          break
        case 'withdraw':
          await withdraw(selectedProtocol, selectedAsset, amount)
          break
        case 'borrow':
          await borrow(selectedProtocol, selectedAsset, amount)
          break
        case 'repay':
          await repay(selectedProtocol, selectedAsset, amount)
          break
      }
      setAmount('')
    } catch (error) {
      console.error('Action failed:', error)
    }
  }

  const formatNumber = (value: string | number, decimals: number = 2) => {
    return parseFloat(value.toString()).toFixed(decimals)
  }

  const formatPercentage = (value: number) => {
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value)
  }

  const getHealthFactorColor = (healthFactor: number) => {
    if (healthFactor < 1.2) return 'text-red-600'
    if (healthFactor < 1.5) return 'text-orange-600'
    if (healthFactor < 2.0) return 'text-yellow-600'
    return 'text-green-600'
  }

  const getRiskLevelColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'text-green-600 bg-green-100 dark:bg-green-900'
      case 'medium':
        return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900'
      case 'high':
        return 'text-orange-600 bg-orange-100 dark:bg-orange-900'
      case 'critical':
        return 'text-red-600 bg-red-100 dark:bg-red-900'
      default:
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900'
    }
  }

  const getProtocolLogo = (protocolId: string) => {
    const logos: Record<string, string> = {
      'aave-v3': '👻',
      'compound-v3': '🏛️',
      'morpho': '🦋'
    }
    return logos[protocolId] || '🏦'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Lending & Borrowing</h2>
          <p className="text-muted-foreground">
            Manage your DeFi lending positions across multiple protocols
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
                <p className="text-sm font-medium text-muted-foreground">Total Supplied</p>
                <p className="text-2xl font-bold">{formatCurrency(analytics.totalSupplied)}</p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Across {analytics.totalPositions} positions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Borrowed</p>
                <p className="text-2xl font-bold">{formatCurrency(analytics.totalBorrowed)}</p>
              </div>
              <TrendingDown className="w-8 h-8 text-red-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Current debt positions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Net Worth</p>
                <p className="text-2xl font-bold text-green-600">
                  {formatCurrency(analytics.totalNetWorth)}
                </p>
              </div>
              <DollarSign className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Portfolio value
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Health Factor</p>
                <p className={cn("text-2xl font-bold", getHealthFactorColor(analytics.averageHealthFactor))}>
                  {analytics.averageHealthFactor.toFixed(2)}
                </p>
              </div>
              <Shield className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Liquidation safety
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="positions">Positions ({state.positions.length})</TabsTrigger>
          <TabsTrigger value="markets">Markets</TabsTrigger>
          <TabsTrigger value="yield">Yield ({yieldOpportunities.length})</TabsTrigger>
          <TabsTrigger value="liquidations">Liquidations ({liquidationOpportunities.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="positions" className="space-y-6">
          {state.positions.length > 0 ? (
            <div className="space-y-4">
              {state.positions.map((position, index) => (
                <motion.div
                  key={position.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <Card>
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <span className="text-2xl">{getProtocolLogo(position.protocolId)}</span>
                          <div>
                            <CardTitle className="text-lg">{position.protocolId}</CardTitle>
                            <CardDescription>
                              Health Factor: <span className={getHealthFactorColor(position.healthFactor)}>
                                {position.healthFactor.toFixed(2)}
                              </span>
                            </CardDescription>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge className={cn("text-xs", getRiskLevelColor(position.riskLevel))}>
                            {position.riskLevel.toUpperCase()}
                          </Badge>
                          {position.isAtRisk && (
                            <Badge variant="destructive">
                              <AlertTriangle className="w-3 h-3 mr-1" />
                              At Risk
                            </Badge>
                          )}
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="grid gap-6 md:grid-cols-2">
                        {/* Supply Positions */}
                        <div className="space-y-3">
                          <h4 className="font-medium text-green-600">Supply Positions</h4>
                          {position.supplies.length > 0 ? (
                            <div className="space-y-2">
                              {position.supplies.map((supply, idx) => (
                                <div key={idx} className="flex items-center justify-between p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
                                  <div>
                                    <p className="font-medium">{supply.asset.symbol}</p>
                                    <p className="text-sm text-muted-foreground">
                                      {formatNumber(supply.amount)} • {formatCurrency(supply.amountUSD)}
                                    </p>
                                  </div>
                                  <div className="text-right">
                                    <p className="font-medium text-green-600">
                                      {formatPercentage(supply.totalAPY)}
                                    </p>
                                    <p className="text-xs text-muted-foreground">
                                      {supply.isCollateral ? 'Collateral' : 'Not Collateral'}
                                    </p>
                                  </div>
                                </div>
                              ))}
                            </div>
                          ) : (
                            <p className="text-sm text-muted-foreground">No supply positions</p>
                          )}
                        </div>

                        {/* Borrow Positions */}
                        <div className="space-y-3">
                          <h4 className="font-medium text-red-600">Borrow Positions</h4>
                          {position.borrows.length > 0 ? (
                            <div className="space-y-2">
                              {position.borrows.map((borrow, idx) => (
                                <div key={idx} className="flex items-center justify-between p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
                                  <div>
                                    <p className="font-medium">{borrow.asset.symbol}</p>
                                    <p className="text-sm text-muted-foreground">
                                      {formatNumber(borrow.amount)} • {formatCurrency(borrow.amountUSD)}
                                    </p>
                                  </div>
                                  <div className="text-right">
                                    <p className="font-medium text-red-600">
                                      {formatPercentage(borrow.netAPY)}
                                    </p>
                                    <p className="text-xs text-muted-foreground">
                                      {borrow.rateMode} rate
                                    </p>
                                  </div>
                                </div>
                              ))}
                            </div>
                          ) : (
                            <p className="text-sm text-muted-foreground">No borrow positions</p>
                          )}
                        </div>
                      </div>

                      {/* Position Summary */}
                      <div className="mt-6 pt-4 border-t">
                        <div className="grid gap-4 md:grid-cols-3">
                          <div className="text-center">
                            <p className="text-sm text-muted-foreground">LTV Ratio</p>
                            <p className="text-lg font-medium">{position.ltv.toFixed(1)}%</p>
                            <Progress value={position.ltv} className="mt-1" />
                          </div>
                          <div className="text-center">
                            <p className="text-sm text-muted-foreground">Net APY</p>
                            <p className={cn("text-lg font-medium", 
                              position.netAPY >= 0 ? "text-green-600" : "text-red-600"
                            )}>
                              {formatPercentage(position.netAPY)}
                            </p>
                          </div>
                          <div className="text-center">
                            <p className="text-sm text-muted-foreground">Available to Borrow</p>
                            <p className="text-lg font-medium">
                              {formatNumber(position.availableBorrowsETH)} ETH
                            </p>
                          </div>
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
                <TrendingUp className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Lending Positions</h3>
                <p className="text-muted-foreground mb-4">
                  Start earning yield by supplying assets to lending protocols
                </p>
                <Button onClick={() => setActiveTab('markets')}>
                  <Plus className="w-4 h-4 mr-2" />
                  Explore Markets
                </Button>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="markets" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Action Panel */}
            <Card>
              <CardHeader>
                <CardTitle>Lending Actions</CardTitle>
                <CardDescription>
                  Supply, withdraw, borrow, or repay assets
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Protocol Selection */}
                <div className="space-y-2">
                  <Label>Protocol</Label>
                  <div className="flex gap-2">
                    {state.protocols.map((protocol) => (
                      <Button
                        key={protocol.id}
                        variant={selectedProtocol === protocol.id ? "default" : "outline"}
                        size="sm"
                        onClick={() => setSelectedProtocol(protocol.id)}
                      >
                        <span className="mr-2">{getProtocolLogo(protocol.id)}</span>
                        {protocol.name}
                      </Button>
                    ))}
                  </div>
                </div>

                {/* Action Type */}
                <div className="space-y-2">
                  <Label>Action</Label>
                  <div className="grid grid-cols-2 gap-2">
                    {(['supply', 'withdraw', 'borrow', 'repay'] as const).map((action) => (
                      <Button
                        key={action}
                        variant={actionType === action ? "default" : "outline"}
                        size="sm"
                        onClick={() => setActionType(action)}
                      >
                        {action === 'supply' && <Plus className="w-4 h-4 mr-1" />}
                        {action === 'withdraw' && <Minus className="w-4 h-4 mr-1" />}
                        {action === 'borrow' && <TrendingDown className="w-4 h-4 mr-1" />}
                        {action === 'repay' && <TrendingUp className="w-4 h-4 mr-1" />}
                        {action.charAt(0).toUpperCase() + action.slice(1)}
                      </Button>
                    ))}
                  </div>
                </div>

                {/* Asset Selection */}
                <div className="space-y-2">
                  <Label>Asset</Label>
                  <div className="grid gap-2">
                    {state.assets.slice(0, 4).map((asset) => (
                      <Button
                        key={asset.address}
                        variant={selectedAsset?.address === asset.address ? "default" : "outline"}
                        className="justify-between"
                        onClick={() => setSelectedAsset(asset)}
                      >
                        <span>{asset.symbol}</span>
                        <span className="text-sm text-muted-foreground">
                          {actionType === 'supply' || actionType === 'withdraw' 
                            ? `${asset.rates.supplyAPY.toFixed(2)}% APY`
                            : `${asset.rates.variableBorrowAPY.toFixed(2)}% APY`
                          }
                        </span>
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
                  {selectedAsset && amount && (
                    <p className="text-sm text-muted-foreground">
                      ≈ ${formatNumber(parseFloat(amount) * selectedAsset.priceUSD, 2)}
                    </p>
                  )}
                </div>

                {/* Execute Button */}
                <Button
                  onClick={handleAction}
                  disabled={!selectedAsset || !amount || !selectedProtocol || state.isExecuting}
                  className="w-full"
                >
                  {state.isExecuting ? (
                    <>
                      <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                      Processing...
                    </>
                  ) : (
                    <>
                      <Zap className="w-4 h-4 mr-2" />
                      {actionType.charAt(0).toUpperCase() + actionType.slice(1)} {selectedAsset?.symbol}
                    </>
                  )}
                </Button>
              </CardContent>
            </Card>

            {/* Market Overview */}
            <Card>
              <CardHeader>
                <CardTitle>Market Overview</CardTitle>
                <CardDescription>
                  Current lending and borrowing rates
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {state.assets.slice(0, 6).map((asset, index) => (
                    <motion.div
                      key={asset.address}
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.05 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div>
                          <p className="font-medium">{asset.symbol}</p>
                          <p className="text-sm text-muted-foreground">{asset.name}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="flex items-center gap-4">
                          <div className="text-center">
                            <p className="text-sm text-green-600 font-medium">
                              {asset.rates.supplyAPY.toFixed(2)}%
                            </p>
                            <p className="text-xs text-muted-foreground">Supply</p>
                          </div>
                          <div className="text-center">
                            <p className="text-sm text-red-600 font-medium">
                              {asset.rates.variableBorrowAPY.toFixed(2)}%
                            </p>
                            <p className="text-xs text-muted-foreground">Borrow</p>
                          </div>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="yield" className="space-y-4">
          <Card>
            <CardContent className="p-12 text-center">
              <TrendingUp className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">Yield Opportunities</h3>
              <p className="text-muted-foreground">
                Advanced yield farming strategies coming soon
              </p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="liquidations" className="space-y-4">
          <Card>
            <CardContent className="p-12 text-center">
              <AlertTriangle className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">Liquidation Monitor</h3>
              <p className="text-muted-foreground">
                Liquidation monitoring and protection features coming soon
              </p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
