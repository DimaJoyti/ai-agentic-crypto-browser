'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Percent, 
  Zap,
  Shield,
  Target,
  Coins,
  ArrowRightLeft,
  Layers,
  PieChart,
  BarChart3
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface DeFiProtocol {
  id: string
  name: string
  logo: string
  tvl: number
  apy: number
  category: 'lending' | 'dex' | 'yield' | 'staking'
  riskLevel: 'low' | 'medium' | 'high'
  userDeposited?: number
  userEarned?: number
}

interface YieldOpportunity {
  id: string
  protocol: string
  pair: string
  apy: number
  tvl: number
  risk: 'low' | 'medium' | 'high'
  category: string
  minDeposit: number
}

interface StakingPool {
  id: string
  token: string
  apy: number
  totalStaked: number
  userStaked?: number
  lockPeriod: number
  rewards: string[]
}

export function EnhancedDeFiDashboard() {
  const [protocols, setProtocols] = useState<DeFiProtocol[]>([])
  const [yieldOpportunities, setYieldOpportunities] = useState<YieldOpportunity[]>([])
  const [stakingPools, setStakingPools] = useState<StakingPool[]>([])
  const [totalPortfolioValue, setTotalPortfolioValue] = useState(0)
  const [totalEarnings, setTotalEarnings] = useState(0)
  const [activeTab, setActiveTab] = useState('overview')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    // Generate mock data
    const mockProtocols: DeFiProtocol[] = [
      {
        id: 'aave',
        name: 'Aave',
        logo: '/protocols/aave.svg',
        tvl: 8500000000,
        apy: 12.5,
        category: 'lending',
        riskLevel: 'low',
        userDeposited: 5000,
        userEarned: 125
      },
      {
        id: 'uniswap',
        name: 'Uniswap V3',
        logo: '/protocols/uniswap.svg',
        tvl: 5200000000,
        apy: 18.7,
        category: 'dex',
        riskLevel: 'medium',
        userDeposited: 3000,
        userEarned: 89
      },
      {
        id: 'compound',
        name: 'Compound',
        logo: '/protocols/compound.svg',
        tvl: 3100000000,
        apy: 9.8,
        category: 'lending',
        riskLevel: 'low',
        userDeposited: 2000,
        userEarned: 45
      },
      {
        id: 'curve',
        name: 'Curve Finance',
        logo: '/protocols/curve.svg',
        tvl: 4800000000,
        apy: 15.2,
        category: 'dex',
        riskLevel: 'medium',
        userDeposited: 1500,
        userEarned: 67
      }
    ]

    const mockYieldOpportunities: YieldOpportunity[] = [
      {
        id: 'eth-usdc-pool',
        protocol: 'Uniswap V3',
        pair: 'ETH/USDC',
        apy: 24.5,
        tvl: 450000000,
        risk: 'medium',
        category: 'Liquidity Pool',
        minDeposit: 100
      },
      {
        id: 'steth-lending',
        protocol: 'Aave',
        pair: 'stETH',
        apy: 16.8,
        tvl: 890000000,
        risk: 'low',
        category: 'Lending',
        minDeposit: 50
      },
      {
        id: 'crv-staking',
        protocol: 'Curve',
        pair: 'CRV',
        apy: 32.1,
        tvl: 120000000,
        risk: 'high',
        category: 'Staking',
        minDeposit: 25
      }
    ]

    const mockStakingPools: StakingPool[] = [
      {
        id: 'eth-staking',
        token: 'ETH',
        apy: 5.2,
        totalStaked: 28000000,
        userStaked: 2.5,
        lockPeriod: 0,
        rewards: ['ETH']
      },
      {
        id: 'matic-staking',
        token: 'MATIC',
        apy: 8.7,
        totalStaked: 5600000000,
        userStaked: 1000,
        lockPeriod: 21,
        rewards: ['MATIC']
      }
    ]

    setProtocols(mockProtocols)
    setYieldOpportunities(mockYieldOpportunities)
    setStakingPools(mockStakingPools)
    
    // Calculate totals
    const totalValue = mockProtocols.reduce((sum, p) => sum + (p.userDeposited || 0), 0)
    const totalEarn = mockProtocols.reduce((sum, p) => sum + (p.userEarned || 0), 0)
    setTotalPortfolioValue(totalValue)
    setTotalEarnings(totalEarn)
  }, [])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatLargeNumber = (num: number) => {
    if (num >= 1e9) {
      return `$${(num / 1e9).toFixed(1)}B`
    }
    if (num >= 1e6) {
      return `$${(num / 1e6).toFixed(1)}M`
    }
    if (num >= 1e3) {
      return `$${(num / 1e3).toFixed(1)}K`
    }
    return `$${num.toFixed(0)}`
  }

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'low': return 'text-green-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskBadgeVariant = (risk: string) => {
    switch (risk) {
      case 'low': return 'default'
      case 'medium': return 'secondary'
      case 'high': return 'destructive'
      default: return 'outline'
    }
  }

  return (
    <div className="space-y-6">
      {/* Portfolio Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Value</span>
            </div>
            <div className="text-2xl font-bold">{formatCurrency(totalPortfolioValue)}</div>
            <div className="text-xs text-green-500 flex items-center gap-1">
              <TrendingUp className="w-3 h-3" />
              +12.5% (24h)
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Percent className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Earnings</span>
            </div>
            <div className="text-2xl font-bold text-green-500">{formatCurrency(totalEarnings)}</div>
            <div className="text-xs text-muted-foreground">
              Avg APY: 14.2%
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Shield className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Active Positions</span>
            </div>
            <div className="text-2xl font-bold">{protocols.filter(p => p.userDeposited).length}</div>
            <div className="text-xs text-muted-foreground">
              Across {protocols.length} protocols
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
              {Math.max(...yieldOpportunities.map(y => y.apy)).toFixed(1)}%
            </div>
            <div className="text-xs text-muted-foreground">
              CRV Staking
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Dashboard */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="protocols">Protocols</TabsTrigger>
          <TabsTrigger value="yield">Yield Farming</TabsTrigger>
          <TabsTrigger value="staking">Staking</TabsTrigger>
          <TabsTrigger value="bridge">Bridge</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Active Positions */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Active Positions
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {protocols.filter(p => p.userDeposited).map((protocol) => (
                    <div key={protocol.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                          <Coins className="w-4 h-4" />
                        </div>
                        <div>
                          <div className="font-medium">{protocol.name}</div>
                          <div className="text-xs text-muted-foreground capitalize">
                            {protocol.category}
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-medium">{formatCurrency(protocol.userDeposited!)}</div>
                        <div className="text-xs text-green-500">
                          +{formatCurrency(protocol.userEarned!)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Top Yield Opportunities */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5" />
                  Top Opportunities
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {yieldOpportunities.slice(0, 3).map((opportunity) => (
                    <div key={opportunity.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div>
                        <div className="font-medium">{opportunity.pair}</div>
                        <div className="text-xs text-muted-foreground">
                          {opportunity.protocol} • {opportunity.category}
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-green-500">{opportunity.apy.toFixed(1)}% APY</div>
                        <Badge variant={getRiskBadgeVariant(opportunity.risk)} className="text-xs">
                          {opportunity.risk} risk
                        </Badge>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="protocols" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {protocols.map((protocol) => (
              <Card key={protocol.id} className="hover:shadow-md transition-shadow">
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Coins className="w-5 h-5" />
                      </div>
                      <div>
                        <CardTitle className="text-lg">{protocol.name}</CardTitle>
                        <Badge variant="outline" className="text-xs capitalize">
                          {protocol.category}
                        </Badge>
                      </div>
                    </div>
                    <Badge variant={getRiskBadgeVariant(protocol.riskLevel)}>
                      {protocol.riskLevel}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between">
                      <span className="text-sm text-muted-foreground">TVL</span>
                      <span className="font-medium">{formatLargeNumber(protocol.tvl)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm text-muted-foreground">APY</span>
                      <span className="font-medium text-green-500">{protocol.apy}%</span>
                    </div>
                    {protocol.userDeposited && (
                      <>
                        <div className="flex justify-between">
                          <span className="text-sm text-muted-foreground">Your Deposit</span>
                          <span className="font-medium">{formatCurrency(protocol.userDeposited)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-sm text-muted-foreground">Earned</span>
                          <span className="font-medium text-green-500">
                            +{formatCurrency(protocol.userEarned!)}
                          </span>
                        </div>
                      </>
                    )}
                    <Button className="w-full" variant={protocol.userDeposited ? "outline" : "default"}>
                      {protocol.userDeposited ? "Manage" : "Deposit"}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="yield" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="w-5 h-5" />
                Yield Farming Opportunities
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {yieldOpportunities.map((opportunity) => (
                  <div key={opportunity.id} className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <Layers className="w-6 h-6" />
                      </div>
                      <div>
                        <div className="font-medium">{opportunity.pair}</div>
                        <div className="text-sm text-muted-foreground">
                          {opportunity.protocol} • {opportunity.category}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          TVL: {formatLargeNumber(opportunity.tvl)} • Min: {formatCurrency(opportunity.minDeposit)}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-xl font-bold text-green-500">{opportunity.apy.toFixed(1)}%</div>
                      <div className="text-sm text-muted-foreground">APY</div>
                      <Badge variant={getRiskBadgeVariant(opportunity.risk)} className="text-xs mt-1">
                        {opportunity.risk} risk
                      </Badge>
                    </div>
                    <Button>
                      Farm
                    </Button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="staking" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {stakingPools.map((pool) => (
              <Card key={pool.id}>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Shield className="w-5 h-5" />
                    {pool.token} Staking
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between">
                      <span className="text-sm text-muted-foreground">APY</span>
                      <span className="font-bold text-green-500">{pool.apy}%</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm text-muted-foreground">Total Staked</span>
                      <span className="font-medium">{formatLargeNumber(pool.totalStaked)}</span>
                    </div>
                    {pool.userStaked && (
                      <div className="flex justify-between">
                        <span className="text-sm text-muted-foreground">Your Stake</span>
                        <span className="font-medium">{pool.userStaked} {pool.token}</span>
                      </div>
                    )}
                    <div className="flex justify-between">
                      <span className="text-sm text-muted-foreground">Lock Period</span>
                      <span className="font-medium">
                        {pool.lockPeriod === 0 ? 'Flexible' : `${pool.lockPeriod} days`}
                      </span>
                    </div>
                    <Button className="w-full" variant={pool.userStaked ? "outline" : "default"}>
                      {pool.userStaked ? "Manage Stake" : "Start Staking"}
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="bridge" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ArrowRightLeft className="w-5 h-5" />
                Cross-Chain Bridge
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-center py-8 text-muted-foreground">
                <ArrowRightLeft className="w-12 h-12 mx-auto mb-4 opacity-50" />
                <h3 className="text-lg font-medium mb-2">Cross-Chain Bridge</h3>
                <p>Transfer assets between different blockchain networks</p>
                <Button className="mt-4">
                  Coming Soon
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
