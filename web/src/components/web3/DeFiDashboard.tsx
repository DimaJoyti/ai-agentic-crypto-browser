'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  TrendingUp, 
  DollarSign, 
  Zap, 
  Shield, 
  Target,
  Coins,
  BarChart3,
  RefreshCw,
  ExternalLink,
  ArrowUpRight,
  ArrowDownLeft,
  Repeat,
  Lock,
  Unlock,
  Flame,
  Droplets,
  PieChart
} from 'lucide-react'
import { useDeFiProtocols } from '@/hooks/useDeFiProtocols'
import { ProtocolType } from '@/lib/defi-protocols'
import { SUPPORTED_CHAINS } from '@/lib/chains'

interface DeFiDashboardProps {
  chainId: number
  userAddress?: string
}

export function DeFiDashboard({ chainId, userAddress }: DeFiDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const {
    protocols,
    lendingPositions,
    yieldPositions,
    isLoading,
    error,
    lastUpdated,
    refreshData,
    dexProtocols,
    lendingProtocols,
    stakingProtocols,
    yieldFarmingProtocols,
    portfolioValue,
    topProtocols,
    lowRiskProtocols,
    mediumRiskProtocols,
    highRiskProtocols
  } = useDeFiProtocols({
    chainId,
    userAddress: userAddress as any,
    autoRefresh: true
  })

  const chain = SUPPORTED_CHAINS[chainId]

  const getProtocolIcon = (protocolId: string) => {
    const icons: Record<string, string> = {
      'uniswap-v3': '🦄',
      'aave-v3': '👻',
      'compound-v3': '🏛️',
      'curve': '🌀',
      'lido': '🏊',
      'pancakeswap': '🥞'
    }
    return icons[protocolId] || '🔗'
  }

  const getTypeIcon = (type: ProtocolType) => {
    switch (type) {
      case ProtocolType.DEX:
        return <Repeat className="w-4 h-4" />
      case ProtocolType.LENDING:
        return <DollarSign className="w-4 h-4" />
      case ProtocolType.STAKING:
        return <Lock className="w-4 h-4" />
      case ProtocolType.YIELD_FARMING:
        return <Target className="w-4 h-4" />
      default:
        return <Coins className="w-4 h-4" />
    }
  }

  const getRiskColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'bg-green-100 text-green-800'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800'
      case 'high':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getAuditColor = (auditStatus: string) => {
    switch (auditStatus) {
      case 'audited':
        return 'bg-green-100 text-green-800'
      case 'partially_audited':
        return 'bg-yellow-100 text-yellow-800'
      case 'unaudited':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatCurrency = (amount: number) => {
    if (amount >= 1e9) return `$${(amount / 1e9).toFixed(2)}B`
    if (amount >= 1e6) return `$${(amount / 1e6).toFixed(2)}M`
    if (amount >= 1e3) return `$${(amount / 1e3).toFixed(2)}K`
    return `$${amount.toFixed(2)}`
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <TrendingUp className="w-6 h-6" />
            DeFi Dashboard
          </h2>
          <p className="text-muted-foreground">
            Explore and interact with decentralized finance protocols on {chain?.name}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            {chain?.shortName || 'Unknown'}
          </Badge>
          <Button variant="outline" size="sm" onClick={refreshData}>
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Portfolio Overview */}
      {userAddress && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Supplied</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioValue.totalSupplied)}</p>
                </div>
                <ArrowUpRight className="w-8 h-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Borrowed</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioValue.totalBorrowed)}</p>
                </div>
                <ArrowDownLeft className="w-8 h-8 text-red-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Net Worth</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioValue.netWorth)}</p>
                </div>
                <DollarSign className="w-8 h-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Total Yield</p>
                  <p className="text-2xl font-bold">{formatCurrency(portfolioValue.totalYield)}</p>
                </div>
                <Target className="w-8 h-8 text-yellow-500" />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* DeFi Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="dex">DEX</TabsTrigger>
          <TabsTrigger value="lending">Lending</TabsTrigger>
          <TabsTrigger value="staking">Staking</TabsTrigger>
          <TabsTrigger value="yield">Yield Farming</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Top Protocols */}
          <Card>
            <CardHeader>
              <CardTitle>Top Protocols by TVL</CardTitle>
              <CardDescription>
                Most popular DeFi protocols on {chain?.name}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topProtocols.map((protocol, index) => (
                  <motion.div
                    key={protocol.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center text-2xl">
                        {getProtocolIcon(protocol.id)}
                      </div>
                      <div>
                        <h4 className="font-medium">{protocol.name}</h4>
                        <p className="text-sm text-muted-foreground">{protocol.description}</p>
                        <div className="flex items-center gap-2 mt-1">
                          {getTypeIcon(protocol.type)}
                          <span className="text-xs text-muted-foreground capitalize">
                            {protocol.type.replace('_', ' ')}
                          </span>
                          <Badge variant="secondary" className={getRiskColor(protocol.riskLevel)}>
                            {protocol.riskLevel} risk
                          </Badge>
                        </div>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <p className="font-bold text-lg">{protocol.tvl}</p>
                      <p className="text-sm text-muted-foreground">TVL</p>
                      {protocol.apy && (
                        <p className="text-sm text-green-600">{protocol.apy} APY</p>
                      )}
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Protocol Categories */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">DEX Protocols</p>
                    <p className="text-2xl font-bold">{dexProtocols.length}</p>
                  </div>
                  <Repeat className="w-8 h-8 text-blue-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Lending Protocols</p>
                    <p className="text-2xl font-bold">{lendingProtocols.length}</p>
                  </div>
                  <DollarSign className="w-8 h-8 text-green-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Staking Protocols</p>
                    <p className="text-2xl font-bold">{stakingProtocols.length}</p>
                  </div>
                  <Lock className="w-8 h-8 text-purple-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">Yield Farming</p>
                    <p className="text-2xl font-bold">{yieldFarmingProtocols.length}</p>
                  </div>
                  <Target className="w-8 h-8 text-yellow-500" />
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Risk Analysis */}
          <Card>
            <CardHeader>
              <CardTitle>Risk Distribution</CardTitle>
              <CardDescription>
                Protocol distribution by risk level
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Shield className="w-4 h-4 text-green-500" />
                    <span>Low Risk</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">{lowRiskProtocols.length} protocols</span>
                    <Progress value={(lowRiskProtocols.length / protocols.length) * 100} className="w-24 h-2" />
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Flame className="w-4 h-4 text-yellow-500" />
                    <span>Medium Risk</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">{mediumRiskProtocols.length} protocols</span>
                    <Progress value={(mediumRiskProtocols.length / protocols.length) * 100} className="w-24 h-2" />
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Zap className="w-4 h-4 text-red-500" />
                    <span>High Risk</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">{highRiskProtocols.length} protocols</span>
                    <Progress value={(highRiskProtocols.length / protocols.length) * 100} className="w-24 h-2" />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="dex" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Decentralized Exchanges</CardTitle>
              <CardDescription>
                Trade tokens on decentralized exchanges
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {dexProtocols.map((protocol) => (
                  <motion.div
                    key={protocol.id}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <span className="text-2xl">{getProtocolIcon(protocol.id)}</span>
                        <div>
                          <h4 className="font-medium">{protocol.name}</h4>
                          <p className="text-sm text-muted-foreground">{protocol.tvl} TVL</p>
                        </div>
                      </div>
                      <Button size="sm" variant="outline">
                        <ExternalLink className="w-3 h-3 mr-2" />
                        Trade
                      </Button>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex flex-wrap gap-1">
                        {protocol.features.map((feature, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {feature}
                          </Badge>
                        ))}
                      </div>
                      
                      <div className="flex items-center justify-between text-sm">
                        <Badge className={getRiskColor(protocol.riskLevel)}>
                          {protocol.riskLevel} risk
                        </Badge>
                        <Badge className={getAuditColor(protocol.auditStatus)}>
                          {protocol.auditStatus.replace('_', ' ')}
                        </Badge>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="lending" className="space-y-6">
          {/* User Positions */}
          {userAddress && lendingPositions.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Your Lending Positions</CardTitle>
                <CardDescription>
                  Current lending and borrowing positions
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {lendingPositions.map((position, index) => (
                    <div key={index} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <span className="text-2xl">👻</span>
                          <div>
                            <h4 className="font-medium">{position.protocol}</h4>
                            <p className="text-sm text-muted-foreground">{position.asset.symbol}</p>
                          </div>
                        </div>
                        <Badge variant="outline">
                          Health: {position.healthFactor}
                        </Badge>
                      </div>
                      
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <p className="text-muted-foreground">Supplied</p>
                          <p className="font-medium">{position.supplied} {position.asset.symbol}</p>
                          <p className="text-green-600">{position.supplyApy} APY</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Borrowed</p>
                          <p className="font-medium">{position.borrowed} {position.asset.symbol}</p>
                          <p className="text-red-600">{position.borrowApy} APY</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Collateral Ratio</p>
                          <p className="font-medium">{position.collateralRatio}</p>
                        </div>
                        <div>
                          <p className="text-muted-foreground">Health Factor</p>
                          <p className={`font-medium ${parseFloat(position.healthFactor) > 1.5 ? 'text-green-600' : 'text-red-600'}`}>
                            {position.healthFactor}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Available Lending Protocols */}
          <Card>
            <CardHeader>
              <CardTitle>Lending Protocols</CardTitle>
              <CardDescription>
                Lend and borrow assets across different protocols
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {lendingProtocols.map((protocol) => (
                  <motion.div
                    key={protocol.id}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <span className="text-2xl">{getProtocolIcon(protocol.id)}</span>
                        <div>
                          <h4 className="font-medium">{protocol.name}</h4>
                          <p className="text-sm text-muted-foreground">{protocol.tvl} TVL</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <Button size="sm">
                          Lend
                        </Button>
                        {protocol.apy && (
                          <p className="text-sm text-green-600 mt-1">{protocol.apy} APY</p>
                        )}
                      </div>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex flex-wrap gap-1">
                        {protocol.features.map((feature, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {feature}
                          </Badge>
                        ))}
                      </div>
                      
                      <div className="flex items-center justify-between text-sm">
                        <Badge className={getRiskColor(protocol.riskLevel)}>
                          {protocol.riskLevel} risk
                        </Badge>
                        <Badge className={getAuditColor(protocol.auditStatus)}>
                          {protocol.auditStatus.replace('_', ' ')}
                        </Badge>
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
              <CardTitle>Staking Protocols</CardTitle>
              <CardDescription>
                Stake your assets to earn rewards
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {stakingProtocols.map((protocol) => (
                  <motion.div
                    key={protocol.id}
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <span className="text-2xl">{getProtocolIcon(protocol.id)}</span>
                        <div>
                          <h4 className="font-medium">{protocol.name}</h4>
                          <p className="text-sm text-muted-foreground">{protocol.tvl} TVL</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <Button size="sm">
                          Stake
                        </Button>
                        {protocol.apy && (
                          <p className="text-sm text-green-600 mt-1">{protocol.apy} APY</p>
                        )}
                      </div>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex flex-wrap gap-1">
                        {protocol.features.map((feature, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {feature}
                          </Badge>
                        ))}
                      </div>
                      
                      <div className="flex items-center justify-between text-sm">
                        <Badge className={getRiskColor(protocol.riskLevel)}>
                          {protocol.riskLevel} risk
                        </Badge>
                        <Badge className={getAuditColor(protocol.auditStatus)}>
                          {protocol.auditStatus.replace('_', ' ')}
                        </Badge>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="yield" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Yield Farming</CardTitle>
              <CardDescription>
                Provide liquidity and earn farming rewards
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-8">
                <Target className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Yield Farming Coming Soon</h3>
                <p className="text-muted-foreground">
                  Advanced yield farming strategies and liquidity mining will be available soon
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
