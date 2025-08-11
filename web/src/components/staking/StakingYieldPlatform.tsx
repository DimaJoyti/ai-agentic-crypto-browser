'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { 
  Shield, 
  Coins, 
  TrendingUp, 
  Clock, 
  Zap,
  Target,
  DollarSign,
  Percent,
  Calendar,
  Lock,
  Unlock,
  ArrowRight,
  Info,
  AlertCircle
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount, useBalance } from 'wagmi'

interface StakingPool {
  id: string
  token: string
  name: string
  apy: number
  totalStaked: number
  userStaked: number
  userRewards: number
  lockPeriod: number
  minStake: number
  maxStake: number
  isActive: boolean
  autoCompound: boolean
  riskLevel: 'low' | 'medium' | 'high'
  features: string[]
}

interface YieldFarm {
  id: string
  protocol: string
  pair: string
  apy: number
  tvl: number
  userDeposited: number
  userRewards: number
  multiplier: number
  endDate: number
  isActive: boolean
  riskLevel: 'low' | 'medium' | 'high'
}

interface LiquidityPool {
  id: string
  pair: string
  protocol: string
  apy: number
  tvl: number
  userLiquidity: number
  userRewards: number
  fees24h: number
  volume24h: number
  impermanentLoss: number
  riskLevel: 'low' | 'medium' | 'high'
}

export function StakingYieldPlatform() {
  const [stakingPools, setStakingPools] = useState<StakingPool[]>([])
  const [yieldFarms, setYieldFarms] = useState<YieldFarm[]>([])
  const [liquidityPools, setLiquidityPools] = useState<LiquidityPool[]>([])
  const [selectedPool, setSelectedPool] = useState<StakingPool | null>(null)
  const [stakeAmount, setStakeAmount] = useState('')
  const [activeTab, setActiveTab] = useState('staking')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    // Generate mock data
    const mockStakingPools: StakingPool[] = [
      {
        id: 'eth-staking',
        token: 'ETH',
        name: 'Ethereum 2.0 Staking',
        apy: 5.2,
        totalStaked: 28000000,
        userStaked: 2.5,
        userRewards: 0.13,
        lockPeriod: 0,
        minStake: 0.1,
        maxStake: 1000,
        isActive: true,
        autoCompound: true,
        riskLevel: 'low',
        features: ['Liquid Staking', 'Auto-compound', 'No Lock Period']
      },
      {
        id: 'matic-staking',
        token: 'MATIC',
        name: 'Polygon Staking',
        apy: 8.7,
        totalStaked: 5600000000,
        userStaked: 1000,
        userRewards: 87,
        lockPeriod: 21,
        minStake: 1,
        maxStake: 100000,
        isActive: true,
        autoCompound: false,
        riskLevel: 'medium',
        features: ['High APY', 'Validator Rewards', '21-day Unbonding']
      },
      {
        id: 'atom-staking',
        token: 'ATOM',
        name: 'Cosmos Hub Staking',
        apy: 12.4,
        totalStaked: 280000000,
        userStaked: 500,
        userRewards: 62,
        lockPeriod: 21,
        minStake: 1,
        maxStake: 50000,
        isActive: true,
        autoCompound: true,
        riskLevel: 'medium',
        features: ['High APY', 'IBC Rewards', 'Governance Rights']
      }
    ]

    const mockYieldFarms: YieldFarm[] = [
      {
        id: 'uni-eth-usdc',
        protocol: 'Uniswap V3',
        pair: 'ETH/USDC',
        apy: 24.5,
        tvl: 450000000,
        userDeposited: 5000,
        userRewards: 122.5,
        multiplier: 2.5,
        endDate: Date.now() + 30 * 24 * 60 * 60 * 1000,
        isActive: true,
        riskLevel: 'medium'
      },
      {
        id: 'curve-3pool',
        protocol: 'Curve Finance',
        pair: '3Pool',
        apy: 18.2,
        tvl: 890000000,
        userDeposited: 3000,
        userRewards: 54.6,
        multiplier: 1.8,
        endDate: Date.now() + 60 * 24 * 60 * 60 * 1000,
        isActive: true,
        riskLevel: 'low'
      }
    ]

    const mockLiquidityPools: LiquidityPool[] = [
      {
        id: 'uni-eth-usdt',
        pair: 'ETH/USDT',
        protocol: 'Uniswap V3',
        apy: 15.8,
        tvl: 320000000,
        userLiquidity: 2500,
        userRewards: 39.5,
        fees24h: 45000,
        volume24h: 15000000,
        impermanentLoss: -2.1,
        riskLevel: 'medium'
      },
      {
        id: 'curve-steth',
        pair: 'stETH/ETH',
        protocol: 'Curve Finance',
        apy: 9.4,
        tvl: 680000000,
        userLiquidity: 1800,
        userRewards: 16.92,
        fees24h: 12000,
        volume24h: 4000000,
        impermanentLoss: -0.3,
        riskLevel: 'low'
      }
    ]

    setStakingPools(mockStakingPools)
    setYieldFarms(mockYieldFarms)
    setLiquidityPools(mockLiquidityPools)
  }, [])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 2
    }).format(amount)
  }

  const formatLargeNumber = (num: number) => {
    if (num >= 1e9) {
      return `${(num / 1e9).toFixed(1)}B`
    }
    if (num >= 1e6) {
      return `${(num / 1e6).toFixed(1)}M`
    }
    if (num >= 1e3) {
      return `${(num / 1e3).toFixed(1)}K`
    }
    return num.toFixed(0)
  }

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'low': return 'text-green-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskBadgeVariant = (level: string) => {
    switch (level) {
      case 'low': return 'default'
      case 'medium': return 'secondary'
      case 'high': return 'destructive'
      default: return 'outline'
    }
  }

  const handleStake = (pool: StakingPool) => {
    setSelectedPool(pool)
    // In real implementation, this would open a staking modal
    console.log('Staking in pool:', pool.id, 'Amount:', stakeAmount)
  }

  const calculateRewards = (pool: StakingPool, amount: number) => {
    const dailyRate = pool.apy / 365 / 100
    return amount * dailyRate
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Your Wallet</h3>
          <p className="text-muted-foreground">
            Connect your wallet to start earning rewards through staking and yield farming
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Staked</span>
            </div>
            <div className="text-2xl font-bold">
              {formatCurrency(stakingPools.reduce((sum, pool) => sum + pool.userStaked * 2500, 0))}
            </div>
            <div className="text-xs text-green-500">
              Across {stakingPools.filter(p => p.userStaked > 0).length} pools
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Percent className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Rewards</span>
            </div>
            <div className="text-2xl font-bold text-green-500">
              {formatCurrency(stakingPools.reduce((sum, pool) => sum + pool.userRewards * 2500, 0))}
            </div>
            <div className="text-xs text-muted-foreground">
              Pending rewards
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Avg APY</span>
            </div>
            <div className="text-2xl font-bold">
              {(stakingPools.reduce((sum, pool) => sum + pool.apy, 0) / stakingPools.length).toFixed(1)}%
            </div>
            <div className="text-xs text-muted-foreground">
              Weighted average
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Target className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Best APY</span>
            </div>
            <div className="text-2xl font-bold text-green-500">
              {Math.max(...yieldFarms.map(f => f.apy)).toFixed(1)}%
            </div>
            <div className="text-xs text-muted-foreground">
              ETH/USDC Farm
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Platform */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="staking">Staking Pools</TabsTrigger>
          <TabsTrigger value="farming">Yield Farming</TabsTrigger>
          <TabsTrigger value="liquidity">Liquidity Mining</TabsTrigger>
          <TabsTrigger value="rewards">My Rewards</TabsTrigger>
        </TabsList>

        <TabsContent value="staking" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
            {stakingPools.map((pool) => (
              <Card key={pool.id} className="hover:shadow-md transition-shadow">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Coins className="w-5 h-5" />
                      </div>
                      <div>
                        <CardTitle className="text-lg">{pool.token}</CardTitle>
                        <p className="text-sm text-muted-foreground">{pool.name}</p>
                      </div>
                    </div>
                    <Badge variant={getRiskBadgeVariant(pool.riskLevel)}>
                      {pool.riskLevel}
                    </Badge>
                  </div>
                </CardHeader>
                
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">APY</div>
                      <div className="font-bold text-green-500 text-lg">{pool.apy}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Total Staked</div>
                      <div className="font-medium">{formatLargeNumber(pool.totalStaked)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Lock Period</div>
                      <div className="font-medium flex items-center gap-1">
                        {pool.lockPeriod === 0 ? (
                          <>
                            <Unlock className="w-3 h-3" />
                            Flexible
                          </>
                        ) : (
                          <>
                            <Lock className="w-3 h-3" />
                            {pool.lockPeriod} days
                          </>
                        )}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Min Stake</div>
                      <div className="font-medium">{pool.minStake} {pool.token}</div>
                    </div>
                  </div>

                  {pool.userStaked > 0 && (
                    <div className="p-3 bg-muted/50 rounded-lg">
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-sm text-muted-foreground">Your Stake</span>
                        <span className="font-medium">{pool.userStaked} {pool.token}</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-muted-foreground">Pending Rewards</span>
                        <span className="font-medium text-green-500">
                          {pool.userRewards.toFixed(4)} {pool.token}
                        </span>
                      </div>
                    </div>
                  )}

                  <div className="space-y-2">
                    <div className="flex flex-wrap gap-1">
                      {pool.features.map((feature) => (
                        <Badge key={feature} variant="outline" className="text-xs">
                          {feature}
                        </Badge>
                      ))}
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <Button 
                      className="flex-1" 
                      onClick={() => handleStake(pool)}
                      variant={pool.userStaked > 0 ? "outline" : "default"}
                    >
                      {pool.userStaked > 0 ? "Manage" : "Stake"}
                    </Button>
                    {pool.userRewards > 0 && (
                      <Button variant="outline" size="sm">
                        Claim
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="farming" className="space-y-4">
          <div className="space-y-4">
            {yieldFarms.map((farm) => (
              <Card key={farm.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <TrendingUp className="w-6 h-6" />
                      </div>
                      <div>
                        <h3 className="font-bold text-lg">{farm.pair}</h3>
                        <p className="text-sm text-muted-foreground">{farm.protocol}</p>
                        <div className="flex items-center gap-2 mt-1">
                          <Badge variant={getRiskBadgeVariant(farm.riskLevel)}>
                            {farm.riskLevel} risk
                          </Badge>
                          <Badge variant="outline">
                            {farm.multiplier}x multiplier
                          </Badge>
                        </div>
                      </div>
                    </div>

                    <div className="text-right">
                      <div className="text-2xl font-bold text-green-500">{farm.apy}%</div>
                      <div className="text-sm text-muted-foreground">APY</div>
                      <div className="text-xs text-muted-foreground mt-1">
                        TVL: {formatCurrency(farm.tvl)}
                      </div>
                    </div>
                  </div>

                  {farm.userDeposited > 0 && (
                    <div className="mt-4 p-3 bg-muted/50 rounded-lg">
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <div className="text-muted-foreground">Your Deposit</div>
                          <div className="font-medium">{formatCurrency(farm.userDeposited)}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Pending Rewards</div>
                          <div className="font-medium text-green-500">
                            {formatCurrency(farm.userRewards)}
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  <div className="mt-4 flex items-center justify-between">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Calendar className="w-4 h-4" />
                      Ends in {Math.ceil((farm.endDate - Date.now()) / (24 * 60 * 60 * 1000))} days
                    </div>
                    <div className="flex gap-2">
                      <Button variant={farm.userDeposited > 0 ? "outline" : "default"}>
                        {farm.userDeposited > 0 ? "Manage" : "Farm"}
                      </Button>
                      {farm.userRewards > 0 && (
                        <Button variant="outline" size="sm">
                          Harvest
                        </Button>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="liquidity" className="space-y-4">
          <div className="space-y-4">
            {liquidityPools.map((pool) => (
              <Card key={pool.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <Coins className="w-6 h-6" />
                      </div>
                      <div>
                        <h3 className="font-bold text-lg">{pool.pair}</h3>
                        <p className="text-sm text-muted-foreground">{pool.protocol}</p>
                        <Badge variant={getRiskBadgeVariant(pool.riskLevel)} className="mt-1">
                          {pool.riskLevel} risk
                        </Badge>
                      </div>
                    </div>

                    <div className="text-right">
                      <div className="text-2xl font-bold text-green-500">{pool.apy}%</div>
                      <div className="text-sm text-muted-foreground">APY</div>
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">TVL</div>
                      <div className="font-medium">{formatCurrency(pool.tvl)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">24h Volume</div>
                      <div className="font-medium">{formatCurrency(pool.volume24h)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">24h Fees</div>
                      <div className="font-medium">{formatCurrency(pool.fees24h)}</div>
                    </div>
                  </div>

                  {pool.userLiquidity > 0 && (
                    <div className="p-3 bg-muted/50 rounded-lg mb-4">
                      <div className="grid grid-cols-3 gap-4 text-sm">
                        <div>
                          <div className="text-muted-foreground">Your Liquidity</div>
                          <div className="font-medium">{formatCurrency(pool.userLiquidity)}</div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Earned Fees</div>
                          <div className="font-medium text-green-500">
                            {formatCurrency(pool.userRewards)}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">IL</div>
                          <div className={cn(
                            "font-medium",
                            pool.impermanentLoss >= 0 ? "text-green-500" : "text-red-500"
                          )}>
                            {pool.impermanentLoss >= 0 ? '+' : ''}{pool.impermanentLoss.toFixed(2)}%
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Info className="w-4 h-4" />
                      Impermanent Loss: {pool.impermanentLoss.toFixed(2)}%
                    </div>
                    <div className="flex gap-2">
                      <Button variant={pool.userLiquidity > 0 ? "outline" : "default"}>
                        {pool.userLiquidity > 0 ? "Manage" : "Add Liquidity"}
                      </Button>
                      {pool.userRewards > 0 && (
                        <Button variant="outline" size="sm">
                          Claim
                        </Button>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="rewards" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="w-5 h-5" />
                  Claimable Rewards
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {stakingPools.filter(p => p.userRewards > 0).map((pool) => (
                    <div key={pool.id} className="flex items-center justify-between p-3 border rounded">
                      <div>
                        <div className="font-medium">{pool.token} Staking</div>
                        <div className="text-sm text-muted-foreground">{pool.name}</div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-green-500">
                          {pool.userRewards.toFixed(4)} {pool.token}
                        </div>
                        <Button size="sm" className="mt-1">
                          Claim
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Target className="w-5 h-5" />
                  Auto-Compound Settings
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {stakingPools.filter(p => p.userStaked > 0).map((pool) => (
                    <div key={pool.id} className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">{pool.token}</div>
                        <div className="text-sm text-muted-foreground">
                          Auto-compound rewards
                        </div>
                      </div>
                      <Switch checked={pool.autoCompound} />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
