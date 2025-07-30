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
  Sprout,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Zap,
  Shield,
  AlertTriangle,
  CheckCircle,
  Clock,
  Activity,
  BarChart3,
  RefreshCw,
  Settings,
  Plus,
  Minus,
  Target,
  ExternalLink,
  Info,
  Lightbulb
} from 'lucide-react'
import { useYieldFarming, useFarmingPosition, useYieldOptimization, useFarmingAnalytics } from '@/hooks/useYieldFarming'
import { type YieldFarm, type FarmType } from '@/lib/yield-farming-integration'
import { cn } from '@/lib/utils'

export function YieldFarmingDashboard() {
  const [activeTab, setActiveTab] = useState('farms')
  const [selectedFarm, setSelectedFarm] = useState<YieldFarm | null>(null)
  const [actionType, setActionType] = useState<'stake' | 'unstake' | 'claim' | 'compound'>('stake')
  const [amount, setAmount] = useState('')
  const [filterProtocol, setFilterProtocol] = useState<string>('')
  const [filterMinAPY, setFilterMinAPY] = useState<number>(0)
  const [filterRisk, setFilterRisk] = useState<string>('')

  // Mock functions
  const useAdvancedYieldFarming = (_options?: any) => ({
    state: {
      farms: [],
      positions: [],
      protocols: [],
      analytics: {},
      isLoading: false,
      isExecuting: false,
      error: null
    },
    getFarms: async (_protocol?: string, _minAPY?: string, _risk?: string) => {},
    stake: async (_farmId: string, _amount: string) => {},
    unstake: async (_farmId: string, _amount: string) => {},
    claimRewards: async (_farmId: string) => {},
    compound: async (_farmId: string) => {},
    getOptimizationRecommendations: async () => {},
    refreshPositions: async () => {},
    clearError: () => {}
  })

  const getYieldFarms = async () => {}

  const {
    state,
    getFarms,
    stake,
    unstake,
    claimRewards,
    compound,
    getOptimizationRecommendations,
    refreshPositions,
    clearError
  } = useAdvancedYieldFarming({
    autoLoad: true,
    enableNotifications: true,
    autoRefresh: true,
    enableOptimization: true
  })

  useEffect(() => {
    getYieldFarms()
  }, [])

  const handleAction = async () => {
    if (!selectedFarm || (!amount && actionType !== 'claim' && actionType !== 'compound')) return

    try {
      switch (actionType) {
        case 'stake':
          await stake(selectedFarm.id, amount)
          break
        case 'unstake':
          await unstake(selectedFarm.id, amount)
          break
        case 'claim':
          await claimRewards(selectedFarm.id)
          break
        case 'compound':
          await compound(selectedFarm.id)
          break
      }
      setAmount('')
    } catch (error) {
      console.error('Action failed:', error)
    }
  }

  const handleFilter = async () => {
    // Apply filters to the existing farms data
    await getYieldFarms()
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

  const formatPercentage = (value: number) => {
    return `${value.toFixed(2)}%`
  }

  const getRiskLevelColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'text-green-600 bg-green-100 dark:bg-green-900'
      case 'medium':
        return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900'
      case 'high':
        return 'text-orange-600 bg-orange-100 dark:bg-orange-900'
      case 'extreme':
        return 'text-red-600 bg-red-100 dark:bg-red-900'
      default:
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900'
    }
  }

  const getProtocolLogo = (protocolId: string) => {
    const logos: Record<string, string> = {
      'pancakeswap': '🥞',
      'curve': '🌀',
      'convex': '🔺'
    }
    return logos[protocolId] || '🌾'
  }

  const getFarmTypeIcon = (farmType: FarmType) => {
    switch (farmType) {
      case 'liquidity_pool':
        return <Activity className="w-4 h-4" />
      case 'single_token':
        return <Target className="w-4 h-4" />
      case 'vault':
        return <Shield className="w-4 h-4" />
      default:
        return <Sprout className="w-4 h-4" />
    }
  }

  // Calculate portfolio analytics
  const portfolioAnalytics = {
    totalStaked: state.positions.reduce((sum: number, p: any) => sum + (p.stakedAmountUSD || 0), 0),
    totalRewards: state.positions.reduce((sum: number, p: any) => sum + (p.totalRewardsClaimedUSD || 0), 0),
    totalPendingRewards: state.positions.reduce((sum: number, p: any) => sum + (p.pendingRewardsUSD || 0), 0),
    averageAPY: state.positions.length > 0
      ? state.positions.reduce((sum: number, p: any) => sum + (p.currentAPY || 0), 0) / state.positions.length
      : 0,
    totalPositions: state.positions.length
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Yield Farming</h2>
          <p className="text-muted-foreground">
            Maximize your returns through automated yield farming strategies
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={refreshPositions}>
            <RefreshCw className="w-4 h-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" size="sm" onClick={() => getOptimizationRecommendations()}>
            <Lightbulb className="w-4 h-4 mr-2" />
            Optimize
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
                <p className="text-sm font-medium text-muted-foreground">Total Staked</p>
                <p className="text-2xl font-bold">{formatCurrency(portfolioAnalytics.totalStaked)}</p>
              </div>
              <Sprout className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Across {portfolioAnalytics.totalPositions} farms
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Rewards</p>
                <p className="text-2xl font-bold text-green-600">
                  {formatCurrency(portfolioAnalytics.totalRewards)}
                </p>
              </div>
              <TrendingUp className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Lifetime earnings
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Pending Rewards</p>
                <p className="text-2xl font-bold text-orange-600">
                  {formatCurrency(portfolioAnalytics.totalPendingRewards)}
                </p>
              </div>
              <Clock className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Ready to claim
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Average APY</p>
                <p className="text-2xl font-bold text-blue-600">
                  {formatPercentage(portfolioAnalytics.averageAPY)}
                </p>
              </div>
              <BarChart3 className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Portfolio yield
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="farms">Available Farms ({state.farms.length})</TabsTrigger>
          <TabsTrigger value="positions">My Positions ({state.positions.length})</TabsTrigger>
          <TabsTrigger value="optimization">Optimization</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="farms" className="space-y-6">
          {/* Filters */}
          <Card>
            <CardHeader>
              <CardTitle>Farm Filters</CardTitle>
              <CardDescription>Filter farms by protocol, APY, and risk level</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-4">
                <div className="space-y-2">
                  <Label>Protocol</Label>
                  <select
                    value={filterProtocol}
                    onChange={(e) => setFilterProtocol(e.target.value)}
                    className="w-full p-2 border rounded-md"
                  >
                    <option value="">All Protocols</option>
                    {state.protocols.map((protocol: any) => (
                      <option key={protocol.id} value={protocol.id}>
                        {protocol.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <Label>Min APY (%)</Label>
                  <Input
                    type="number"
                    placeholder="0"
                    value={filterMinAPY || ''}
                    onChange={(e) => setFilterMinAPY(parseFloat(e.target.value) || 0)}
                  />
                </div>

                <div className="space-y-2">
                  <Label>Max Risk</Label>
                  <select
                    value={filterRisk}
                    onChange={(e) => setFilterRisk(e.target.value)}
                    className="w-full p-2 border rounded-md"
                  >
                    <option value="">All Risk Levels</option>
                    <option value="low">Low Risk</option>
                    <option value="medium">Medium Risk</option>
                    <option value="high">High Risk</option>
                    <option value="extreme">Extreme Risk</option>
                  </select>
                </div>

                <div className="flex items-end">
                  <Button onClick={handleFilter} className="w-full">
                    Apply Filters
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Farm List */}
          {state.farms.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              <AnimatePresence>
                {state.farms.map((farm: any, index: number) => (
                  <motion.div
                    key={farm.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ delay: index * 0.05 }}
                  >
                    <Card 
                      className={cn(
                        "transition-all duration-200 cursor-pointer hover:shadow-md",
                        selectedFarm?.id === farm.id && "ring-2 ring-primary"
                      )}
                      onClick={() => setSelectedFarm(farm)}
                    >
                      <CardHeader className="pb-3">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <span className="text-xl">{getProtocolLogo(farm.protocolId)}</span>
                            <div>
                              <CardTitle className="text-lg">{farm.name}</CardTitle>
                              <CardDescription className="flex items-center gap-1">
                                {getFarmTypeIcon(farm.type)}
                                {farm.type.replace('_', ' ')}
                              </CardDescription>
                            </div>
                          </div>
                          <div className="flex flex-col items-end gap-1">
                            <Badge className={cn("text-xs", getRiskLevelColor(farm.riskLevel))}>
                              {farm.riskLevel.toUpperCase()}
                            </Badge>
                            {!farm.isActive && (
                              <Badge variant="secondary" className="text-xs">
                                Inactive
                              </Badge>
                            )}
                          </div>
                        </div>
                      </CardHeader>
                      <CardContent>
                        <div className="space-y-3">
                          <div className="flex items-center justify-between">
                            <span className="text-sm text-muted-foreground">APY</span>
                            <span className="text-lg font-bold text-green-600">
                              {formatPercentage(farm.apy)}
                            </span>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">TVL</span>
                            <span>{formatCurrency(parseFloat(farm.tvl))}</span>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Daily Rewards</span>
                            <span>{formatCurrency(parseFloat(farm.dailyRewards))}</span>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">IL Risk</span>
                            <span className={cn(
                              farm.impermanentLossRisk > 20 ? 'text-red-600' :
                              farm.impermanentLossRisk > 10 ? 'text-yellow-600' : 'text-green-600'
                            )}>
                              {farm.impermanentLossRisk}%
                            </span>
                          </div>

                          {/* Token Composition */}
                          <div className="space-y-1">
                            <span className="text-xs text-muted-foreground">Tokens</span>
                            <div className="flex flex-wrap gap-1">
                              {farm.tokens.map((token: any, idx: number) => (
                                <Badge key={idx} variant="outline" className="text-xs">
                                  {token.symbol}
                                </Badge>
                              ))}
                            </div>
                          </div>

                          {/* Reward Tokens */}
                          <div className="space-y-1">
                            <span className="text-xs text-muted-foreground">Rewards</span>
                            <div className="flex flex-wrap gap-1">
                              {farm.rewardTokens.map((reward: any, idx: number) => (
                                <Badge key={idx} variant="secondary" className="text-xs">
                                  {reward.symbol}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Sprout className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Farms Found</h3>
                <p className="text-muted-foreground mb-4">
                  {state.isLoading ? 'Loading farms...' : 'No farms match your current filters'}
                </p>
                {!state.isLoading && (
                  <Button onClick={() => getYieldFarms()}>
                    <RefreshCw className="w-4 h-4 mr-2" />
                    Refresh Farms
                  </Button>
                )}
              </CardContent>
            </Card>
          )}

          {/* Action Panel */}
          {selectedFarm && (
            <Card>
              <CardHeader>
                <CardTitle>Farm Actions - {selectedFarm.name}</CardTitle>
                <CardDescription>
                  Stake, unstake, or claim rewards from this farm
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Action Type Selection */}
                <div className="grid grid-cols-2 gap-2 md:grid-cols-4">
                  {(['stake', 'unstake', 'claim', 'compound'] as const).map((action) => (
                    <Button
                      key={action}
                      variant={actionType === action ? "default" : "outline"}
                      size="sm"
                      onClick={() => setActionType(action)}
                    >
                      {action === 'stake' && <Plus className="w-4 h-4 mr-1" />}
                      {action === 'unstake' && <Minus className="w-4 h-4 mr-1" />}
                      {action === 'claim' && <DollarSign className="w-4 h-4 mr-1" />}
                      {action === 'compound' && <Zap className="w-4 h-4 mr-1" />}
                      {action.charAt(0).toUpperCase() + action.slice(1)}
                    </Button>
                  ))}
                </div>

                {/* Amount Input */}
                {(actionType === 'stake' || actionType === 'unstake') && (
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
                        ≈ ${formatNumber(parseFloat(amount) * 100, 2)} USD
                      </p>
                    )}
                  </div>
                )}

                {/* Execute Button */}
                <Button
                  onClick={handleAction}
                  disabled={
                    (!amount && actionType !== 'claim' && actionType !== 'compound') || 
                    state.isExecuting
                  }
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
                      {actionType.charAt(0).toUpperCase() + actionType.slice(1)}
                    </>
                  )}
                </Button>

                {/* Farm Details */}
                <div className="pt-4 border-t space-y-2">
                  <h4 className="font-medium">Farm Details</h4>
                  <div className="grid gap-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Protocol:</span>
                      <span>{selectedFarm.protocolId}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Strategy:</span>
                      <span>{selectedFarm.strategy.name}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Auto-Compound:</span>
                      <span>{selectedFarm.strategy.autoCompound ? 'Yes' : 'No'}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Deposit Fee:</span>
                      <span>{selectedFarm.depositFee}%</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Performance Fee:</span>
                      <span>{selectedFarm.performanceFee}%</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="positions" className="space-y-4">
          {state.positions.length > 0 ? (
            <div className="space-y-4">
              {state.positions.map((position: any, index: number) => (
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
                            <CardTitle className="text-lg">{position.farmId}</CardTitle>
                            <CardDescription>
                              Staked: {formatCurrency(position.stakedAmountUSD)}
                            </CardDescription>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-lg font-bold text-green-600">
                            {formatPercentage(position.currentAPY)}
                          </p>
                          <p className="text-sm text-muted-foreground">Current APY</p>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="grid gap-4 md:grid-cols-3">
                        <div className="space-y-2">
                          <h4 className="font-medium text-green-600">Rewards</h4>
                          <div className="space-y-1">
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">Pending:</span>
                              <span>{formatCurrency(position.pendingRewardsUSD)}</span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">Claimed:</span>
                              <span>{formatCurrency(position.totalRewardsClaimedUSD)}</span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">Total:</span>
                              <span className="font-medium">
                                {formatCurrency(position.pendingRewardsUSD + position.totalRewardsClaimedUSD)}
                              </span>
                            </div>
                          </div>
                        </div>

                        <div className="space-y-2">
                          <h4 className="font-medium text-blue-600">Performance</h4>
                          <div className="space-y-1">
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">ROI:</span>
                              <span className={cn(
                                position.roi >= 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {formatPercentage(position.roi)}
                              </span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">IL:</span>
                              <span className={cn(
                                position.impermanentLoss <= 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {formatPercentage(position.impermanentLoss)}
                              </span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span className="text-muted-foreground">Net P&L:</span>
                              <span className={cn(
                                position.netProfitLoss >= 0 ? 'text-green-600' : 'text-red-600'
                              )}>
                                {formatCurrency(position.netProfitLoss)}
                              </span>
                            </div>
                          </div>
                        </div>

                        <div className="space-y-2">
                          <h4 className="font-medium text-purple-600">Actions</h4>
                          <div className="space-y-2">
                            <Button
                              size="sm"
                              className="w-full"
                              onClick={() => claimRewards(position.farmId)}
                              disabled={position.pendingRewardsUSD === 0}
                            >
                              <DollarSign className="w-4 h-4 mr-1" />
                              Claim Rewards
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              className="w-full"
                              onClick={() => compound(position.farmId)}
                              disabled={position.pendingRewardsUSD === 0}
                            >
                              <Zap className="w-4 h-4 mr-1" />
                              Compound
                            </Button>
                          </div>
                        </div>
                      </div>

                      {/* Position Timeline */}
                      <div className="mt-4 pt-4 border-t">
                        <div className="flex items-center justify-between text-sm text-muted-foreground">
                          <span>Deposited: {new Date(position.depositedAt).toLocaleDateString()}</span>
                          {position.lastClaimAt > 0 && (
                            <span>Last Claim: {new Date(position.lastClaimAt).toLocaleDateString()}</span>
                          )}
                          {position.autoCompound && (
                            <Badge variant="secondary" className="text-xs">
                              Auto-Compound
                            </Badge>
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
                <Sprout className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Active Positions</h3>
                <p className="text-muted-foreground mb-4">
                  Start farming to see your positions here
                </p>
                <Button onClick={() => setActiveTab('farms')}>
                  <Plus className="w-4 h-4 mr-2" />
                  Explore Farms
                </Button>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="optimization" className="space-y-4">
          <Card>
            <CardContent className="p-12 text-center">
              <Lightbulb className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">Yield Optimization</h3>
              <p className="text-muted-foreground mb-4">
                AI-powered yield optimization recommendations coming soon
              </p>
              <Button onClick={() => getOptimizationRecommendations()}>
                <Target className="w-4 h-4 mr-2" />
                Get Recommendations
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="space-y-4">
          <Card>
            <CardContent className="p-12 text-center">
              <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">Farming Analytics</h3>
              <p className="text-muted-foreground">
                Detailed analytics and performance tracking coming soon
              </p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
