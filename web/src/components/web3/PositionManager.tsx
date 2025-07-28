'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Target, 
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
  RefreshCw,
  Bell,
  BellOff,
  Eye,
  EyeOff,
  BarChart3,
  PieChart,
  Activity
} from 'lucide-react'
import { usePositionManager } from '@/hooks/usePositionManager'
import { PositionType, PositionStatus, type DeFiPosition } from '@/lib/position-manager'
import { type Address } from 'viem'

interface PositionManagerProps {
  userAddress?: Address
}

export function PositionManager({ userAddress }: PositionManagerProps) {
  const [activeTab, setActiveTab] = useState('overview')
  const [showAlerts, setShowAlerts] = useState(true)

  const {
    positions,
    alerts,
    summary,
    isLoading,
    error,
    lastUpdated,
    loadPositions,
    acknowledgeAlert,
    portfolioMetrics,
    topPerformers,
    worstPerformers,
    positionsRequiringAttention,
    positionTypeDistribution,
    protocolDistribution,
    formatCurrency,
    formatPercentage,
    lendingPositions,
    borrowingPositions,
    liquidityPositions,
    stakingPositions,
    criticalAlerts,
    highAlerts
  } = usePositionManager({
    userAddress,
    autoRefresh: true,
    enableNotifications: true
  })

  const getPositionTypeIcon = (type: PositionType) => {
    switch (type) {
      case PositionType.LENDING:
        return <DollarSign className="w-4 h-4" />
      case PositionType.BORROWING:
        return <TrendingDown className="w-4 h-4" />
      case PositionType.LIQUIDITY:
        return <Droplets className="w-4 h-4" />
      case PositionType.STAKING:
        return <Lock className="w-4 h-4" />
      case PositionType.YIELD_FARMING:
        return <Target className="w-4 h-4" />
      default:
        return <Activity className="w-4 h-4" />
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

  const getAlertSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'border-red-500 bg-red-50'
      case 'high':
        return 'border-orange-500 bg-orange-50'
      case 'medium':
        return 'border-yellow-500 bg-yellow-50'
      case 'low':
        return 'border-blue-500 bg-blue-50'
      default:
        return 'border-gray-500 bg-gray-50'
    }
  }

  const getPnlColor = (pnl: string) => {
    const value = parseFloat(pnl)
    return value >= 0 ? 'text-green-600' : 'text-red-600'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Target className="w-6 h-6" />
            Position Manager
          </h2>
          <p className="text-muted-foreground">
            Track and manage your DeFi positions across all protocols
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowAlerts(!showAlerts)}
          >
            {showAlerts ? <BellOff className="w-4 h-4" /> : <Bell className="w-4 h-4" />}
          </Button>
          <Button variant="outline" size="sm" onClick={loadPositions}>
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Critical Alerts */}
      {showAlerts && (criticalAlerts.length > 0 || highAlerts.length > 0) && (
        <div className="space-y-2">
          {criticalAlerts.map((alert) => (
            <Alert key={alert.id} variant="destructive">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription className="flex items-center justify-between">
                <span>{alert.title}: {alert.message}</span>
                <Button size="sm" variant="outline" onClick={() => acknowledgeAlert(alert.id)}>
                  Acknowledge
                </Button>
              </AlertDescription>
            </Alert>
          ))}
          {highAlerts.map((alert) => (
            <Alert key={alert.id} className="border-orange-500 bg-orange-50">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription className="flex items-center justify-between">
                <span>{alert.title}: {alert.message}</span>
                <Button size="sm" variant="outline" onClick={() => acknowledgeAlert(alert.id)}>
                  Acknowledge
                </Button>
              </AlertDescription>
            </Alert>
          ))}
        </div>
      )}

      {/* Portfolio Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Value</p>
                <p className="text-2xl font-bold">{formatCurrency(portfolioMetrics.totalValue)}</p>
                <p className={`text-sm ${getPnlColor(portfolioMetrics.totalPnl.toString())}`}>
                  {formatPercentage(portfolioMetrics.totalPnlPercentage)}
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
                <p className="text-sm font-medium text-muted-foreground">Total P&L</p>
                <p className={`text-2xl font-bold ${getPnlColor(portfolioMetrics.totalPnl.toString())}`}>
                  {formatCurrency(portfolioMetrics.totalPnl)}
                </p>
                <p className="text-sm text-muted-foreground">
                  {portfolioMetrics.profitablePositions} profitable
                </p>
              </div>
              {portfolioMetrics.totalPnl >= 0 ? (
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
                <p className="text-sm font-medium text-muted-foreground">Active Positions</p>
                <p className="text-2xl font-bold">{summary.totalPositions}</p>
                <p className="text-sm text-muted-foreground">
                  {portfolioMetrics.atRiskPositions} at risk
                </p>
              </div>
              <Target className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Rewards</p>
                <p className="text-2xl font-bold">{formatCurrency(portfolioMetrics.totalRewards)}</p>
                <p className="text-sm text-muted-foreground">
                  {summary.activeAlerts} alerts
                </p>
              </div>
              <Zap className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="positions">All Positions</TabsTrigger>
          <TabsTrigger value="alerts">Alerts</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Positions Requiring Attention */}
          {positionsRequiringAttention.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="w-5 h-5 text-yellow-500" />
                  Positions Requiring Attention
                </CardTitle>
                <CardDescription>
                  Positions that need immediate review or action
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {positionsRequiringAttention.map((position) => (
                    <div key={position.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        {getPositionTypeIcon(position.type)}
                        <div>
                          <h4 className="font-medium">{position.protocol}</h4>
                          <p className="text-sm text-muted-foreground capitalize">
                            {position.type.replace('_', ' ')}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <Badge className={getStatusColor(position.status)}>
                          {position.status.replace('_', ' ')}
                        </Badge>
                        <p className="text-sm text-muted-foreground mt-1">
                          {formatCurrency(position.totalValue)}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Position Type Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Position Distribution</CardTitle>
                <CardDescription>
                  Breakdown by position type
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {positionTypeDistribution.map((item) => (
                    <div key={item.type} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          {getPositionTypeIcon(item.type)}
                          <span className="text-sm font-medium capitalize">
                            {item.type.replace('_', ' ')}
                          </span>
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {item.count} ({item.percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <Progress value={item.percentage} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Protocol Distribution</CardTitle>
                <CardDescription>
                  Breakdown by protocol
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {protocolDistribution.map((item) => (
                    <div key={item.protocol} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-sm font-medium">{item.protocol}</span>
                        <span className="text-sm text-muted-foreground">
                          {item.count} ({item.percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <Progress value={item.percentage} className="h-2" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="positions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>All Positions</CardTitle>
              <CardDescription>
                Complete overview of your DeFi positions
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {positions.map((position, index) => (
                  <motion.div
                    key={position.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="border rounded-lg p-4"
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        {getPositionTypeIcon(position.type)}
                        <div>
                          <h4 className="font-medium">{position.protocol}</h4>
                          <p className="text-sm text-muted-foreground capitalize">
                            {position.type.replace('_', ' ')}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <Badge className={getStatusColor(position.status)}>
                          {position.status.replace('_', ' ')}
                        </Badge>
                        {position.apy && (
                          <p className="text-sm text-green-600 mt-1">{position.apy}</p>
                        )}
                      </div>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div>
                        <p className="text-sm text-muted-foreground">Total Value</p>
                        <p className="font-medium">{formatCurrency(position.totalValue)}</p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">P&L</p>
                        <p className={`font-medium ${getPnlColor(position.pnl)}`}>
                          {formatCurrency(position.pnl)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">P&L %</p>
                        <p className={`font-medium ${getPnlColor(position.pnl)}`}>
                          {formatPercentage(position.pnlPercentage)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground">Rewards</p>
                        <p className="font-medium">
                          {position.rewards
                            ? formatCurrency(position.rewards.reduce((sum, r) => sum + parseFloat(r.value), 0))
                            : '$0.00'
                          }
                        </p>
                      </div>
                    </div>

                    <div className="flex gap-2 mt-4">
                      <Button size="sm" variant="outline">
                        <Eye className="w-3 h-3 mr-2" />
                        View Details
                      </Button>
                      {position.type === PositionType.LIQUIDITY && (
                        <Button size="sm" variant="outline">
                          Manage LP
                        </Button>
                      )}
                      {position.type === PositionType.STAKING && position.rewards && position.rewards.length > 0 && (
                        <Button size="sm" variant="outline">
                          <Zap className="w-3 h-3 mr-2" />
                          Claim Rewards
                        </Button>
                      )}
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="alerts" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Active Alerts</CardTitle>
              <CardDescription>
                Monitor important position updates and warnings
              </CardDescription>
            </CardHeader>
            <CardContent>
              {alerts.length === 0 ? (
                <div className="text-center py-8">
                  <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No Active Alerts</h3>
                  <p className="text-muted-foreground">
                    All your positions are healthy and performing well
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  {alerts.map((alert) => (
                    <div
                      key={alert.id}
                      className={`border rounded-lg p-4 ${getAlertSeverityColor(alert.severity)}`}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium">{alert.title}</h4>
                          <p className="text-sm text-muted-foreground mt-1">
                            {alert.message}
                          </p>
                          <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                            <span>Threshold: {alert.threshold}</span>
                            <span>Current: {alert.currentValue}</span>
                            <span>
                              {new Date(alert.createdAt).toLocaleTimeString()}
                            </span>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="capitalize">
                            {alert.severity}
                          </Badge>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => acknowledgeAlert(alert.id)}
                          >
                            Acknowledge
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="w-5 h-5" />
                  Portfolio Breakdown
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Lending Positions</span>
                    <span className="font-medium">{lendingPositions.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Borrowing Positions</span>
                    <span className="font-medium">{borrowingPositions.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Liquidity Positions</span>
                    <span className="font-medium">{liquidityPositions.length}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Staking Positions</span>
                    <span className="font-medium">{stakingPositions.length}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Risk Assessment
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Shield className="w-4 h-4 text-green-500" />
                      <span className="text-sm text-muted-foreground">Healthy Positions</span>
                    </div>
                    <span className="font-medium">
                      {positions.filter(p => p.status === PositionStatus.ACTIVE).length}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <AlertTriangle className="w-4 h-4 text-yellow-500" />
                      <span className="text-sm text-muted-foreground">At Risk</span>
                    </div>
                    <span className="font-medium">{portfolioMetrics.atRiskPositions}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Zap className="w-4 h-4 text-red-500" />
                      <span className="text-sm text-muted-foreground">Liquidatable</span>
                    </div>
                    <span className="font-medium">{portfolioMetrics.liquidatablePositions}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="performance" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5 text-green-500" />
                  Top Performers
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {topPerformers.map((position) => (
                    <div key={position.id} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        {getPositionTypeIcon(position.type)}
                        <span className="text-sm font-medium">{position.protocol}</span>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-medium text-green-600">
                          {formatPercentage(position.pnlPercentage)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatCurrency(position.pnl)}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingDown className="w-5 h-5 text-red-500" />
                  Worst Performers
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {worstPerformers.map((position) => (
                    <div key={position.id} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        {getPositionTypeIcon(position.type)}
                        <span className="text-sm font-medium">{position.protocol}</span>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-medium text-red-600">
                          {formatPercentage(position.pnlPercentage)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatCurrency(position.pnl)}
                        </p>
                      </div>
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
