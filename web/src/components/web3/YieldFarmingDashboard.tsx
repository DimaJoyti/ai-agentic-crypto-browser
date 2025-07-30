'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Target, 
  TrendingUp, 
  DollarSign, 
  Zap,
  Shield,
  Lock,
  Droplets,
  RefreshCw,
  Plus,
  Minus,
  Gift,
  BarChart3,
  PieChart,
  Activity,
  AlertTriangle,
  CheckCircle,
  Clock,
  Flame
} from 'lucide-react'
import { useYieldFarming } from '@/hooks/useYieldFarming'
import { YieldStrategy, RiskLevel, type YieldFarm } from '@/lib/yield-farming'
import { type Address } from 'viem'

interface YieldFarmingDashboardProps {
  userAddress?: Address
  chainId?: number
}

export function YieldFarmingDashboard({ userAddress, chainId }: YieldFarmingDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const {
    state,
    getFarms,
    getUserPositions,
    stake,
    unstake,
    harvest,
    refresh
  } = useYieldFarming({
    autoRefresh: true,
    enableNotifications: true
  })

  // Extract data from state
  const farms = state.farms || []
  const userPositions = state.positions || []

  // Mock data for missing properties
  const stakingPools = farms.filter(f => f.type === 'single_token') || []
  const topFarms = farms.slice(0, 5) || []

  // Mock distribution data
  const strategyDistribution = [
    { name: 'Liquidity Mining', value: 45, farmCount: 12 },
    { name: 'Single Staking', value: 30, farmCount: 8 },
    { name: 'Yield Farming', value: 25, farmCount: 6 }
  ]

  const riskDistribution = [
    { name: 'Low Risk', value: 40, farmCount: 10 },
    { name: 'Medium Risk', value: 45, farmCount: 12 },
    { name: 'High Risk', value: 15, farmCount: 4 }
  ]

  // Mock functions for missing functionality
  const stakeFarm = async (farmId: string, amount: string) => {
    return stake(farmId, amount)
  }

  const unstakeFarm = async (farmId: string, amount: string) => {
    return unstake(farmId, amount)
  }

  const claimRewards = async (farmId: string) => {
    return harvest(farmId)
  }

  const loadData = async () => {
    refresh()
  }

  const isLoading = false

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount)
  }

  const portfolioMetrics = {
    totalValue: 0,
    totalRewards: 0,
    totalStaked: 0,
    totalEarned: 0,
    totalPendingRewards: 0,
    totalClaimableRewards: 0,
    totalPnlPercentage: 0,
    averageAPY: 0,
    profitablePositions: 0,
    activePositions: 0,
    apy: 0
  }

  const getStrategyIcon = (strategy: any) => {
    const strategyType = strategy?.type || strategy?.name || 'default'
    switch (strategyType) {
      case 'SINGLE_STAKING':
      case 'single_staking':
        return <Lock className="w-4 h-4" />
      case 'LIQUIDITY_MINING':
      case 'liquidity_mining':
        return <Droplets className="w-4 h-4" />
      case 'YIELD_FARMING':
      case 'yield_farming':
        return <Target className="w-4 h-4" />
      case 'LIQUID_STAKING':
      case 'liquid_staking':
        return <Activity className="w-4 h-4" />
      case 'LENDING_YIELD':
      case 'lending_yield':
        return <DollarSign className="w-4 h-4" />
      case 'AUTOCOMPOUNDING':
      case 'autocompounding':
        return <RefreshCw className="w-4 h-4" />
      default:
        return <Target className="w-4 h-4" />
    }
  }

  const getRiskColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'bg-green-100 text-green-800'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800'
      case 'high':
        return 'bg-orange-100 text-orange-800'
      case 'extreme':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getAPYColor = (apy: number | string) => {
    const apyValue = typeof apy === 'number' ? apy : parseFloat(apy.toString().replace('%', ''))
    if (apyValue >= 20) return 'text-green-600'
    if (apyValue >= 10) return 'text-blue-600'
    if (apyValue >= 5) return 'text-yellow-600'
    return 'text-gray-600'
  }

  const handleStake = async (farmId: string, amount: string) => {
    try {
      await stakeFarm(farmId, amount)
    } catch (error) {
      console.error('Staking failed:', error)
    }
  }

  const handleUnstake = async (positionId: string, amount: string) => {
    try {
      await unstakeFarm(positionId, amount)
    } catch (error) {
      console.error('Unstaking failed:', error)
    }
  }

  const handleClaimRewards = async (positionId: string) => {
    try {
      await claimRewards(positionId)
    } catch (error) {
      console.error('Claiming rewards failed:', error)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Target className="w-6 h-6" />
            Yield Farming & Staking
          </h2>
          <p className="text-muted-foreground">
            Maximize your yields across DeFi protocols with optimized strategies
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadData}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Portfolio Overview */}
      {userAddress && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Staked</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioMetrics.totalStaked)}</p>
                  <p className="text-sm text-muted-foreground">
                    {portfolioMetrics.activePositions} positions
                  </p>
                </div>
                <DollarSign className="w-8 h-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Earned</p>
                  <p className="text-2xl font-bold text-green-600">
                    {formatCurrency(portfolioMetrics.totalEarned)}
                  </p>
                  <p className="text-sm text-green-600">
                    {portfolioMetrics.totalPnlPercentage.toFixed(2)}% gain
                  </p>
                </div>
                <TrendingUp className="w-8 h-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Pending Rewards</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioMetrics.totalPendingRewards)}</p>
                  <p className="text-sm text-muted-foreground">
                    {formatCurrency(portfolioMetrics.totalClaimableRewards)} claimable
                  </p>
                </div>
                <Zap className="w-8 h-8 text-yellow-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Average APY</p>
                  <p className="text-2xl font-bold">{portfolioMetrics.averageAPY.toFixed(1)}%</p>
                  <p className="text-sm text-muted-foreground">
                    {portfolioMetrics.profitablePositions} profitable
                  </p>
                </div>
                <Target className="w-8 h-8 text-purple-500" />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="farms">Yield Farms</TabsTrigger>
          <TabsTrigger value="staking">Staking Pools</TabsTrigger>
          <TabsTrigger value="positions">My Positions</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Top Yield Farms */}
          <Card>
            <CardHeader>
              <CardTitle>Top Yield Opportunities</CardTitle>
              <CardDescription>
                Highest APY farms and staking pools available
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topFarms.map((farm, index) => (
                  <motion.div
                    key={farm.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center">
                        {getStrategyIcon(farm.strategy)}
                      </div>
                      <div>
                        <h4 className="font-medium">{farm.name}</h4>
                        <p className="text-sm text-muted-foreground">{farm.protocolId}</p>
                        <div className="flex items-center gap-2 mt-1">
                          <Badge variant="outline" className="text-xs">
                            {farm.strategy.toString().replace('_', ' ')}
                          </Badge>
                          <Badge className={getRiskColor(farm.riskLevel)}>
                            {farm.riskLevel} risk
                          </Badge>
                        </div>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <p className={`font-bold text-lg ${getAPYColor(farm.apy)}`}>{farm.apy}</p>
                      <p className="text-sm text-muted-foreground">APY</p>
                      <p className="text-sm text-muted-foreground">{farm.tvl} TVL</p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Strategy Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Strategy Distribution</CardTitle>
                <CardDescription>
                  Available yield strategies
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {strategyDistribution.map((item) => (
                    <div key={item.name} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          {getStrategyIcon(item.name)}
                          <span className="text-sm font-medium capitalize">
                            {item.name.replace('_', ' ')}
                          </span>
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {item.farmCount} farms
                        </span>
                      </div>
                      <Progress value={(item.farmCount / farms.length) * 100} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Distribution</CardTitle>
                <CardDescription>
                  Risk levels across available farms
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {riskDistribution.map((item) => (
                    <div key={item.name} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <Shield className="w-4 h-4" />
                          <span className="text-sm font-medium capitalize">
                            {item.name}
                          </span>
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {item.farmCount} farms
                        </span>
                      </div>
                      <Progress value={(item.farmCount / farms.length) * 100} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="farms" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Available Yield Farms</CardTitle>
              <CardDescription>
                Discover and stake in high-yield farming opportunities
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {farms.map((farm) => (
                  <motion.div
                    key={farm.id}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        {getStrategyIcon(farm.strategy)}
                        <div>
                          <h4 className="font-medium">{farm.name}</h4>
                          <p className="text-sm text-muted-foreground">{farm.protocolId}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className={`font-bold text-lg ${getAPYColor(farm.apy)}`}>{farm.apy}</p>
                        <p className="text-xs text-muted-foreground">APY</p>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">TVL</span>
                        <span className="font-medium">{farm.tvl}</span>
                      </div>
                      
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Min Stake</span>
                        <span className="font-medium">1.0 ETH</span>
                      </div>

                      {farm.lockupPeriod && farm.lockupPeriod > 0 && (
                        <div className="flex items-center justify-between text-sm">
                          <span className="text-muted-foreground">Lock Period</span>
                          <span className="font-medium">{Math.floor(farm.lockupPeriod / 86400)} days</span>
                        </div>
                      )}

                      <div className="flex items-center justify-between">
                        <Badge className={getRiskColor(farm.riskLevel)}>
                          {farm.riskLevel} risk
                        </Badge>
                        <Button 
                          size="sm" 
                          onClick={() => handleStake(farm.id, '1.0')}
                          disabled={!userAddress}
                        >
                          <Plus className="w-3 h-3 mr-2" />
                          Stake
                        </Button>
                      </div>

                      <div className="pt-2 border-t">
                        <p className="text-xs text-muted-foreground">High yield farming opportunity</p>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="staking" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Staking Pools</CardTitle>
              <CardDescription>
                Stake tokens to earn rewards and secure networks
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {stakingPools.map((pool) => (
                  <motion.div
                    key={pool.id}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <Lock className="w-6 h-6" />
                        <div>
                          <h4 className="font-medium">{pool.name}</h4>
                          <p className="text-sm text-muted-foreground">{pool.protocolId}</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className={`font-bold text-lg ${getAPYColor(pool.apy)}`}>{pool.apy}</p>
                        <p className="text-xs text-muted-foreground">APY</p>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">TVL</span>
                        <span className="font-medium">{pool.tvl}</span>
                      </div>
                      
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Min Stake</span>
                        <span className="font-medium">32 ETH</span>
                      </div>

                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Validators</span>
                        <span className="font-medium">1,234</span>
                      </div>

                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <Badge className={getRiskColor(pool.riskLevel)}>
                            {pool.riskLevel} risk
                          </Badge>
                          <Badge variant="outline" className="text-xs">
                            Slashing Risk: 0.1%
                          </Badge>
                        </div>
                        <Button 
                          size="sm" 
                          onClick={() => handleStake(pool.id, '32.0')}
                          disabled={!userAddress}
                        >
                          <Plus className="w-3 h-3 mr-2" />
                          Stake
                        </Button>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="positions" className="space-y-6">
          {userPositions.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle>My Yield Positions</CardTitle>
                <CardDescription>
                  Manage your active yield farming and staking positions
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {userPositions.map((position, index) => {
                    const farm = farms.find(f => f.id === position.farmId)
                    if (!farm) return null

                    return (
                      <motion.div
                        key={position.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className="border rounded-lg p-4"
                      >
                        <div className="flex items-center justify-between mb-4">
                          <div className="flex items-center gap-3">
                            {getStrategyIcon(farm.strategy)}
                            <div>
                              <h4 className="font-medium">{farm.name}</h4>
                              <p className="text-sm text-muted-foreground">{farm.protocolId}</p>
                            </div>
                          </div>
                          
                          <div className="flex gap-2">
                            <Button size="sm" variant="outline">
                              <Plus className="w-3 h-3 mr-2" />
                              Add
                            </Button>
                            <Button size="sm" variant="outline">
                              <Minus className="w-3 h-3 mr-2" />
                              Remove
                            </Button>
                            {position.pendingRewards && Array.isArray(position.pendingRewards) && position.pendingRewards.length > 0 && (
                              <Button
                                size="sm"
                                onClick={() => handleClaimRewards(position.id)}
                              >
                                <Gift className="w-3 h-3 mr-2" />
                                Claim
                              </Button>
                            )}
                          </div>
                        </div>

                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                          <div>
                            <p className="text-sm text-muted-foreground">Staked</p>
                            <p className="font-medium">{position.stakedAmount} ETH</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Current Value</p>
                            <p className="font-medium">{formatCurrency(5420.50)}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">Total Earned</p>
                            <p className="font-medium text-green-600">{formatCurrency(142.30)}</p>
                          </div>
                          <div>
                            <p className="text-sm text-muted-foreground">P&L</p>
                            <p className="font-medium text-green-600">
                              {formatCurrency(42.30)} (+2.8%)
                            </p>
                          </div>
                        </div>

                        {position.pendingRewards.length > 0 && (
                          <div className="mt-4 pt-4 border-t">
                            <h5 className="text-sm font-medium mb-2">Pending Rewards</h5>
                            <div className="flex flex-wrap gap-2">
                              {Array.isArray(position.pendingRewards) ? position.pendingRewards.map((reward: any, idx: number) => (
                                <Badge key={idx} variant="outline">
                                  {reward.amount} {reward.symbol}
                                </Badge>
                              )) : null}
                            </div>
                          </div>
                        )}
                      </motion.div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="p-8 text-center">
                <Target className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">No Active Positions</h3>
                <p className="text-muted-foreground mb-4">
                  Start yield farming to earn rewards on your crypto assets
                </p>
                <Button onClick={() => setActiveTab('farms')}>
                  <Plus className="w-4 h-4 mr-2" />
                  Explore Farms
                </Button>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Yield Performance
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Farms</span>
                    <span className="font-medium">{farms.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Staking Pools</span>
                    <span className="font-medium">{stakingPools.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Highest APY</span>
                    <span className="font-medium text-green-600">
                      {farms.length > 0 ? Math.max(...farms.map(f => typeof f.apy === 'number' ? f.apy : parseFloat(String(f.apy)))).toFixed(1) : '0.0'}%
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Average APY</span>
                    <span className="font-medium">
                      {farms.length > 0 ? (farms.reduce((sum, f) => sum + (typeof f.apy === 'number' ? f.apy : parseFloat(String(f.apy))), 0) / farms.length).toFixed(1) : '0.0'}%
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Portfolio Health
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-4 h-4 text-green-500" />
                      <span className="text-sm text-muted-foreground">Active Positions</span>
                    </div>
                    <span className="font-medium">{portfolioMetrics.activePositions}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <TrendingUp className="w-4 h-4 text-green-500" />
                      <span className="text-sm text-muted-foreground">Profitable Positions</span>
                    </div>
                    <span className="font-medium">{portfolioMetrics.profitablePositions}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Zap className="w-4 h-4 text-yellow-500" />
                      <span className="text-sm text-muted-foreground">Claimable Rewards</span>
                    </div>
                    <span className="font-medium">{formatCurrency(portfolioMetrics.totalClaimableRewards)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
