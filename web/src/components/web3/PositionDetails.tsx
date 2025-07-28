'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  ArrowLeft,
  TrendingUp, 
  TrendingDown,
  DollarSign, 
  AlertTriangle,
  CheckCircle,
  Clock,
  Zap,
  Shield,
  Droplets,
  Lock,
  Target,
  Activity,
  ExternalLink,
  RefreshCw,
  Plus,
  Minus,
  Gift,
  BarChart3
} from 'lucide-react'
import { 
  type DeFiPosition, 
  type LendingPosition,
  type BorrowingPosition,
  type LiquidityPosition,
  type StakingPosition,
  PositionType,
  PositionStatus 
} from '@/lib/position-manager'

interface PositionDetailsProps {
  position: DeFiPosition
  onBack: () => void
}

export function PositionDetails({ position, onBack }: PositionDetailsProps) {
  const [activeTab, setActiveTab] = useState('overview')

  const getPositionTypeIcon = (type: PositionType) => {
    switch (type) {
      case PositionType.LENDING:
        return <DollarSign className="w-5 h-5" />
      case PositionType.BORROWING:
        return <TrendingDown className="w-5 h-5" />
      case PositionType.LIQUIDITY:
        return <Droplets className="w-5 h-5" />
      case PositionType.STAKING:
        return <Lock className="w-5 h-5" />
      case PositionType.YIELD_FARMING:
        return <Target className="w-5 h-5" />
      default:
        return <Activity className="w-5 h-5" />
    }
  }

  const getStatusColor = (status: PositionStatus) => {
    switch (status) {
      case PositionStatus.ACTIVE:
        return 'bg-green-100 text-green-800'
      case PositionStatus.AT_RISK:
        return 'bg-yellow-100 text-yellow-800'
      case PositionStatus.LIQUIDATABLE:
        return 'bg-red-100 text-red-800'
      case PositionStatus.INACTIVE:
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-blue-100 text-blue-800'
    }
  }

  const getPnlColor = (pnl: string) => {
    const value = parseFloat(pnl)
    return value >= 0 ? 'text-green-600' : 'text-red-600'
  }

  const formatCurrency = (amount: string | number): string => {
    const value = typeof amount === 'string' ? parseFloat(amount) : amount
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value)
  }

  const formatPercentage = (percentage: string | number): string => {
    const value = typeof percentage === 'string' ? parseFloat(percentage) : percentage
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  const renderLendingDetails = (position: LendingPosition) => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Supplied</p>
                <p className="text-lg font-bold">{position.supplied} {position.asset.symbol}</p>
                <p className="text-sm text-muted-foreground">{formatCurrency(position.asset.value)}</p>
              </div>
              <DollarSign className="w-6 h-6 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Earned</p>
                <p className="text-lg font-bold text-green-600">{position.earned} {position.asset.symbol}</p>
                <p className="text-sm text-green-600">{position.supplyApy} APY</p>
              </div>
              <TrendingUp className="w-6 h-6 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Utilization</p>
                <p className="text-lg font-bold">{position.utilizationRate}</p>
                <p className="text-sm text-muted-foreground">Collateral: {position.collateralFactor}</p>
              </div>
              <Activity className="w-6 h-6 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex gap-2">
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Supply More
        </Button>
        <Button variant="outline">
          <Minus className="w-4 h-4 mr-2" />
          Withdraw
        </Button>
        <Button variant="outline">
          <Gift className="w-4 h-4 mr-2" />
          Claim Rewards
        </Button>
      </div>
    </div>
  )

  const renderBorrowingDetails = (position: BorrowingPosition) => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Borrowed</p>
                <p className="text-lg font-bold">{position.borrowed} {position.asset.symbol}</p>
                <p className="text-sm text-muted-foreground">{formatCurrency(position.asset.value)}</p>
              </div>
              <TrendingDown className="w-6 h-6 text-red-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Current Debt</p>
                <p className="text-lg font-bold text-red-600">{position.debt} {position.asset.symbol}</p>
                <p className="text-sm text-red-600">{position.borrowApy} APY</p>
              </div>
              <AlertTriangle className="w-6 h-6 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Health Factor</p>
                <p className={`text-lg font-bold ${parseFloat(position.healthFactor) > 1.5 ? 'text-green-600' : 'text-red-600'}`}>
                  {position.healthFactor}
                </p>
                <p className="text-sm text-muted-foreground">Threshold: {position.liquidationThreshold}</p>
              </div>
              <Shield className="w-6 h-6 text-green-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Collateral */}
      <Card>
        <CardHeader>
          <CardTitle>Collateral</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {position.collateral.map((collateral, index) => (
              <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                <div>
                  <h4 className="font-medium">{collateral.symbol}</h4>
                  <p className="text-sm text-muted-foreground">{collateral.name}</p>
                </div>
                <div className="text-right">
                  <p className="font-medium">{collateral.amount} {collateral.symbol}</p>
                  <p className="text-sm text-muted-foreground">{formatCurrency(collateral.value)}</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <div className="flex gap-2">
        <Button variant="outline">
          <Plus className="w-4 h-4 mr-2" />
          Add Collateral
        </Button>
        <Button variant="outline">
          <Minus className="w-4 h-4 mr-2" />
          Repay Debt
        </Button>
        <Button>
          <RefreshCw className="w-4 h-4 mr-2" />
          Manage Position
        </Button>
      </div>
    </div>
  )

  const renderLiquidityDetails = (position: LiquidityPosition) => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Liquidity</p>
                <p className="text-lg font-bold">{formatCurrency(position.liquidity)}</p>
                <p className="text-sm text-muted-foreground">Pool Share</p>
              </div>
              <Droplets className="w-6 h-6 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">24h Fees</p>
                <p className="text-lg font-bold text-green-600">{formatCurrency(position.fees24h)}</p>
                <p className="text-sm text-green-600">{position.apy} APY</p>
              </div>
              <TrendingUp className="w-6 h-6 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">IL</p>
                <p className="text-lg font-bold text-red-600">{position.impermanentLoss}</p>
                <p className={`text-sm ${position.inRange ? 'text-green-600' : 'text-red-600'}`}>
                  {position.inRange ? 'In Range' : 'Out of Range'}
                </p>
              </div>
              <Activity className="w-6 h-6 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Pool Composition */}
      <Card>
        <CardHeader>
          <CardTitle>Pool Composition</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 border rounded-lg">
              <div>
                <h4 className="font-medium">{position.pool.token0.symbol}</h4>
                <p className="text-sm text-muted-foreground">{position.pool.token0.name}</p>
              </div>
              <div className="text-right">
                <p className="font-medium">{position.pool.token0.amount} {position.pool.token0.symbol}</p>
                <p className="text-sm text-muted-foreground">{formatCurrency(position.pool.token0.value)}</p>
              </div>
            </div>
            
            <div className="flex items-center justify-between p-3 border rounded-lg">
              <div>
                <h4 className="font-medium">{position.pool.token1.symbol}</h4>
                <p className="text-sm text-muted-foreground">{position.pool.token1.name}</p>
              </div>
              <div className="text-right">
                <p className="font-medium">{position.pool.token1.amount} {position.pool.token1.symbol}</p>
                <p className="text-sm text-muted-foreground">{formatCurrency(position.pool.token1.value)}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Price Range */}
      <Card>
        <CardHeader>
          <CardTitle>Price Range</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Min Price</span>
              <span className="font-medium">{formatCurrency(position.pool.priceRange.min)}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Current Price</span>
              <span className="font-medium">{formatCurrency(position.pool.priceRange.current)}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Max Price</span>
              <span className="font-medium">{formatCurrency(position.pool.priceRange.max)}</span>
            </div>
            
            <div className="mt-4">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-muted-foreground">Position Range</span>
                <span className={`text-sm ${position.inRange ? 'text-green-600' : 'text-red-600'}`}>
                  {position.inRange ? 'Active' : 'Inactive'}
                </span>
              </div>
              <Progress 
                value={position.inRange ? 75 : 25} 
                className="h-2" 
              />
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="flex gap-2">
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Add Liquidity
        </Button>
        <Button variant="outline">
          <Minus className="w-4 h-4 mr-2" />
          Remove Liquidity
        </Button>
        <Button variant="outline">
          <Gift className="w-4 h-4 mr-2" />
          Collect Fees
        </Button>
      </div>
    </div>
  )

  const renderStakingDetails = (position: StakingPosition) => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Staked</p>
                <p className="text-lg font-bold">{position.staked} {position.asset.symbol}</p>
                <p className="text-sm text-muted-foreground">{formatCurrency(position.asset.value)}</p>
              </div>
              <Lock className="w-6 h-6 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Rewards</p>
                <p className="text-lg font-bold text-green-600">
                  {position.rewards.reduce((sum, r) => sum + parseFloat(r.amount), 0).toFixed(4)} {position.asset.symbol}
                </p>
                <p className="text-sm text-green-600">{position.apy} APY</p>
              </div>
              <Zap className="w-6 h-6 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Risk Level</p>
                <p className="text-lg font-bold">{position.slashingRisk}</p>
                {position.unlockDate && (
                  <p className="text-sm text-muted-foreground">
                    Unlocks: {new Date(position.unlockDate).toLocaleDateString()}
                  </p>
                )}
              </div>
              <Shield className="w-6 h-6 text-green-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex gap-2">
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Stake More
        </Button>
        <Button variant="outline">
          <Gift className="w-4 h-4 mr-2" />
          Claim Rewards
        </Button>
        {!position.lockPeriod && (
          <Button variant="outline">
            <Minus className="w-4 h-4 mr-2" />
            Unstake
          </Button>
        )}
      </div>
    </div>
  )

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="outline" size="sm" onClick={onBack}>
            <ArrowLeft className="w-4 h-4" />
          </Button>
          <div>
            <div className="flex items-center gap-3">
              {getPositionTypeIcon(position.type)}
              <h2 className="text-2xl font-bold">{position.protocol}</h2>
              <Badge className={getStatusColor(position.status)}>
                {position.status.replace('_', ' ')}
              </Badge>
            </div>
            <p className="text-muted-foreground capitalize">
              {position.type.replace('_', ' ')} Position
            </p>
          </div>
        </div>
        <Button variant="outline" size="sm">
          <ExternalLink className="w-4 h-4 mr-2" />
          View on Explorer
        </Button>
      </div>

      {/* Position Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                <p className="text-2xl font-bold">{formatCurrency(position.totalValue)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">P&L</p>
                <p className={`text-2xl font-bold ${getPnlColor(position.pnl)}`}>
                  {formatCurrency(position.pnl)}
                </p>
                <p className={`text-sm ${getPnlColor(position.pnl)}`}>
                  {formatPercentage(position.pnlPercentage)}
                </p>
              </div>
              {parseFloat(position.pnl) >= 0 ? (
                <TrendingUp className="w-8 h-8 text-green-500" />
              ) : (
                <TrendingDown className="w-8 h-8 text-red-500" />
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">APY</p>
                <p className="text-2xl font-bold">{position.apy || 'N/A'}</p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Duration</p>
                <p className="text-2xl font-bold">
                  {Math.floor((Date.now() - position.createdAt) / (1000 * 60 * 60 * 24))}d
                </p>
              </div>
              <Clock className="w-8 h-8 text-gray-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Position Details */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="history">History</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview">
          {position.type === PositionType.LENDING && renderLendingDetails(position as LendingPosition)}
          {position.type === PositionType.BORROWING && renderBorrowingDetails(position as BorrowingPosition)}
          {position.type === PositionType.LIQUIDITY && renderLiquidityDetails(position as LiquidityPosition)}
          {position.type === PositionType.STAKING && renderStakingDetails(position as StakingPosition)}
        </TabsContent>

        <TabsContent value="history">
          <Card>
            <CardHeader>
              <CardTitle>Transaction History</CardTitle>
              <CardDescription>
                Recent transactions for this position
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-8">
                <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Transaction History</h3>
                <p className="text-muted-foreground">
                  Transaction history will be available soon
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics">
          <Card>
            <CardHeader>
              <CardTitle>Position Analytics</CardTitle>
              <CardDescription>
                Detailed performance metrics and insights
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-8">
                <BarChart3 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Advanced Analytics</h3>
                <p className="text-muted-foreground">
                  Detailed analytics and performance charts coming soon
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
