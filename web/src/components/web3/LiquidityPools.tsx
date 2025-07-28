'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Droplets, 
  Plus, 
  Minus, 
  TrendingUp, 
  DollarSign,
  Target,
  BarChart3,
  ExternalLink,
  Info,
  Zap,
  Clock,
  Users
} from 'lucide-react'
import { useDeFiProtocols } from '@/hooks/useDeFiProtocols'
import { type Address } from 'viem'

interface LiquidityPoolsProps {
  chainId: number
  userAddress?: Address
}

export function LiquidityPools({ chainId, userAddress }: LiquidityPoolsProps) {
  const [activeTab, setActiveTab] = useState('pools')

  const { dexProtocols } = useDeFiProtocols({
    chainId,
    userAddress
  })

  // Mock liquidity pools data
  const liquidityPools = [
    {
      id: '1',
      protocol: 'Uniswap V3',
      token0: { symbol: 'USDC', name: 'USD Coin' },
      token1: { symbol: 'WETH', name: 'Wrapped Ether' },
      fee: 0.05,
      tvl: '$12.5M',
      volume24h: '$45M',
      apy: '8.5%',
      myLiquidity: '$2,500',
      myShare: '0.02%',
      priceRange: '$2,350 - $2,450'
    },
    {
      id: '2',
      protocol: 'Uniswap V3',
      token0: { symbol: 'USDC', name: 'USD Coin' },
      token1: { symbol: 'DAI', name: 'Dai Stablecoin' },
      fee: 0.01,
      tvl: '$8.2M',
      volume24h: '$12M',
      apy: '3.2%',
      myLiquidity: '$0',
      myShare: '0%',
      priceRange: '$0.998 - $1.002'
    },
    {
      id: '3',
      protocol: 'Curve',
      token0: { symbol: 'USDC', name: 'USD Coin' },
      token1: { symbol: 'USDT', name: 'Tether USD' },
      fee: 0.04,
      tvl: '$25.8M',
      volume24h: '$78M',
      apy: '5.8%',
      myLiquidity: '$1,200',
      myShare: '0.005%',
      priceRange: 'Stable'
    }
  ]

  const myPositions = liquidityPools.filter(pool => parseFloat(pool.myLiquidity.replace(/[$,]/g, '')) > 0)

  const totalLiquidity = myPositions.reduce((sum, pool) => {
    return sum + parseFloat(pool.myLiquidity.replace(/[$,]/g, ''))
  }, 0)

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount)
  }

  const getProtocolIcon = (protocol: string) => {
    const icons: Record<string, string> = {
      'Uniswap V3': '🦄',
      'Curve': '🌀',
      'SushiSwap': '🍣',
      'Balancer': '⚖️'
    }
    return icons[protocol] || '🔗'
  }

  const getApyColor = (apy: string) => {
    const apyValue = parseFloat(apy.replace('%', ''))
    if (apyValue >= 10) return 'text-green-600'
    if (apyValue >= 5) return 'text-yellow-600'
    return 'text-blue-600'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <Droplets className="w-6 h-6" />
          Liquidity Pools
        </h2>
        <p className="text-muted-foreground">
          Provide liquidity to earn trading fees and rewards
        </p>
      </div>

      {/* Portfolio Overview */}
      {userAddress && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Liquidity</p>
                  <p className="text-2xl font-bold">{formatCurrency(totalLiquidity)}</p>
                </div>
                <Droplets className="w-8 h-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Active Positions</p>
                  <p className="text-2xl font-bold">{myPositions.length}</p>
                </div>
                <Target className="w-8 h-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">24h Fees Earned</p>
                  <p className="text-2xl font-bold">$12.45</p>
                </div>
                <DollarSign className="w-8 h-8 text-yellow-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Avg APY</p>
                  <p className="text-2xl font-bold">6.8%</p>
                </div>
                <TrendingUp className="w-8 h-8 text-purple-500" />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="pools">All Pools</TabsTrigger>
          <TabsTrigger value="positions">My Positions</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="pools" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Available Liquidity Pools</CardTitle>
              <CardDescription>
                Discover pools with the best APY and trading volume
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {liquidityPools.map((pool, index) => (
                  <motion.div
                    key={pool.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="flex items-center gap-2">
                          <span className="text-2xl">{getProtocolIcon(pool.protocol)}</span>
                          <div>
                            <h4 className="font-medium">
                              {pool.token0.symbol}/{pool.token1.symbol}
                            </h4>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                              <span>{pool.protocol}</span>
                              <Badge variant="outline" className="text-xs">
                                {pool.fee}% fee
                              </Badge>
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="grid grid-cols-4 gap-8 text-right">
                        <div>
                          <p className="text-sm text-muted-foreground">TVL</p>
                          <p className="font-medium">{pool.tvl}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">24h Volume</p>
                          <p className="font-medium">{pool.volume24h}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">APY</p>
                          <p className={`font-medium ${getApyColor(pool.apy)}`}>{pool.apy}</p>
                        </div>
                        <div className="flex gap-2">
                          <Button size="sm" variant="outline">
                            <Plus className="w-3 h-3 mr-2" />
                            Add
                          </Button>
                          <Button size="sm" variant="outline">
                            <ExternalLink className="w-3 h-3" />
                          </Button>
                        </div>
                      </div>
                    </div>

                    <div className="mt-4 pt-4 border-t">
                      <div className="grid grid-cols-3 gap-4 text-sm">
                        <div>
                          <p className="text-muted-foreground">Price Range</p>
                          <p className="font-medium">{pool.priceRange}</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">My Liquidity</p>
                          <p className="font-medium">{pool.myLiquidity}</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">My Share</p>
                          <p className="font-medium">{pool.myShare}</p>
                        </div>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="positions" className="space-y-6">
          {myPositions.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle>My Liquidity Positions</CardTitle>
                <CardDescription>
                  Manage your active liquidity positions
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {myPositions.map((pool, index) => (
                    <motion.div
                      key={pool.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="border rounded-lg p-4"
                    >
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-4">
                          <span className="text-2xl">{getProtocolIcon(pool.protocol)}</span>
                          <div>
                            <h4 className="font-medium">
                              {pool.token0.symbol}/{pool.token1.symbol}
                            </h4>
                            <p className="text-sm text-muted-foreground">{pool.protocol}</p>
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
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div>
                          <p className="text-sm text-muted-foreground">My Liquidity</p>
                          <p className="font-bold text-lg">{pool.myLiquidity}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">Pool Share</p>
                          <p className="font-medium">{pool.myShare}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">APY</p>
                          <p className={`font-medium ${getApyColor(pool.apy)}`}>{pool.apy}</p>
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">24h Fees</p>
                          <p className="font-medium text-green-600">+$4.25</p>
                        </div>
                      </div>

                      <div className="mt-4 pt-4 border-t">
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-sm text-muted-foreground">Price Range</p>
                            <p className="font-medium">{pool.priceRange}</p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm text-muted-foreground">In Range</p>
                            <div className="flex items-center gap-2">
                              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                              <span className="text-sm font-medium">Active</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="p-8 text-center">
                <Droplets className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">No Liquidity Positions</h3>
                <p className="text-muted-foreground mb-4">
                  Start providing liquidity to earn trading fees and rewards
                </p>
                <Button>
                  <Plus className="w-4 h-4 mr-2" />
                  Add Liquidity
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
                  Pool Performance
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Volume (24h)</span>
                    <span className="font-medium">$135M</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Total Fees (24h)</span>
                    <span className="font-medium">$405K</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Active Pools</span>
                    <span className="font-medium">1,247</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Liquidity Providers</span>
                    <span className="font-medium">8,432</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5" />
                  Top Performing Pools
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {liquidityPools.slice(0, 3).map((pool, index) => (
                    <div key={pool.id} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-lg">{getProtocolIcon(pool.protocol)}</span>
                        <span className="text-sm font-medium">
                          {pool.token0.symbol}/{pool.token1.symbol}
                        </span>
                      </div>
                      <div className="text-right">
                        <p className={`text-sm font-medium ${getApyColor(pool.apy)}`}>
                          {pool.apy}
                        </p>
                        <p className="text-xs text-muted-foreground">{pool.volume24h}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Liquidity Distribution</CardTitle>
              <CardDescription>
                Distribution of liquidity across different protocols
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {dexProtocols.map((protocol, index) => (
                  <div key={protocol.id} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">{protocol.name}</span>
                      <span className="text-sm text-muted-foreground">{protocol.tvl}</span>
                    </div>
                    <Progress value={Math.random() * 100} className="h-2" />
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
