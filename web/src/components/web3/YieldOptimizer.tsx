'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Target, 
  TrendingUp, 
  ArrowRight,
  Lightbulb,
  AlertTriangle,
  CheckCircle,
  RefreshCw,
  Zap,
  DollarSign,
  BarChart3,
  Shield,
  Clock
} from 'lucide-react'
import { useYieldFarming } from '@/hooks/useYieldFarming'
import { RiskLevel, type YieldOptimization } from '@/lib/yield-farming'
import { type Address } from 'viem'

interface YieldOptimizerProps {
  userAddress?: Address
  chainId?: number
}

export function YieldOptimizer({ userAddress, chainId }: YieldOptimizerProps) {
  const [selectedPosition, setSelectedPosition] = useState<string | null>(null)
  const [optimization, setOptimization] = useState<YieldOptimization | null>(null)

  const {
    state
  } = useYieldFarming({
    autoRefresh: true
  })

  // Extract data from state
  const userPositions = state.positions || []
  const farms = state.farms || []

  // Mock functions
  const getYieldOptimization = async (): Promise<YieldOptimization> => ({
    recommendation: 'stay',
    currentPosition: {
      id: '',
      farmId: '',
      stakedAmount: '0',
      pendingRewards: '0',
      apy: 0
    },
    suggestedActions: [],
    projectedReturns: {
      daily: 0,
      weekly: 0,
      monthly: 0,
      yearly: 0
    },
    riskAssessment: {
      level: 'low',
      factors: []
    }
  })
  const formatCurrency = (amount: number) => new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)

  const handleOptimizePosition = async (positionId: string) => {
    setSelectedPosition(positionId)
    const opt = await getYieldOptimization()
    setOptimization(opt)
  }

  const getRecommendationColor = (recommendation: string) => {
    switch (recommendation) {
      case 'migrate':
        return 'bg-green-100 text-green-800 border-green-200'
      case 'diversify':
        return 'bg-blue-100 text-blue-800 border-blue-200'
      case 'stay':
        return 'bg-gray-100 text-gray-800 border-gray-200'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  const getRecommendationIcon = (recommendation: string) => {
    switch (recommendation) {
      case 'migrate':
        return <ArrowRight className="w-4 h-4" />
      case 'diversify':
        return <BarChart3 className="w-4 h-4" />
      case 'stay':
        return <CheckCircle className="w-4 h-4" />
      default:
        return <Target className="w-4 h-4" />
    }
  }

  const getRiskColor = (riskLevel: RiskLevel) => {
    switch (riskLevel) {
      case RiskLevel.LOW:
        return 'bg-green-100 text-green-800'
      case RiskLevel.MEDIUM:
        return 'bg-yellow-100 text-yellow-800'
      case RiskLevel.HIGH:
        return 'bg-orange-100 text-orange-800'
      case RiskLevel.VERY_HIGH:
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const calculateYieldProjection = (apy: string, amount: string, days: number = 365) => {
    const apyValue = parseFloat(apy.replace('%', '')) / 100
    const principal = parseFloat(amount)
    const dailyRate = apyValue / 365
    const projection = principal * dailyRate * days
    return projection
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <Lightbulb className="w-6 h-6" />
          Yield Optimizer
        </h2>
        <p className="text-muted-foreground">
          Analyze and optimize your yield farming positions for maximum returns
        </p>
      </div>

      {userPositions.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <Target className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">No Positions to Optimize</h3>
            <p className="text-muted-foreground">
              Start yield farming to get personalized optimization recommendations
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Position Selection */}
          <Card>
            <CardHeader>
              <CardTitle>Select Position to Optimize</CardTitle>
              <CardDescription>
                Choose a position to analyze for optimization opportunities
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {userPositions.map((position: any) => {
                  const farm = farms.find((f: any) => f.id === position.farmId)
                  if (!farm) return null

                  const isSelected = selectedPosition === position.id
                  const yearlyProjection = calculateYieldProjection(String(typeof farm.apy === 'number' ? farm.apy : parseFloat(String(farm.apy))), parseFloat(String(position.stakedAmount)))

                  return (
                    <motion.div
                      key={position.id}
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className={`border rounded-lg p-4 cursor-pointer transition-all ${
                        isSelected 
                          ? 'ring-2 ring-blue-500 bg-blue-50 dark:bg-blue-950' 
                          : 'hover:bg-accent/50'
                      }`}
                      onClick={() => handleOptimizePosition(position.id)}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <h4 className="font-medium">{farm.name}</h4>
                          <p className="text-sm text-muted-foreground">{farm.protocolId}</p>
                          <div className="flex items-center gap-2 mt-1">
                            <Badge variant="outline" className="text-xs">
                              {farm.apy} APY
                            </Badge>
                            <Badge className={getRiskColor(farm.riskLevel as any)}>
                              {farm.riskLevel}
                            </Badge>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="font-medium">{formatCurrency(parseFloat(position.stakedAmount))}</p>
                          <p className="text-sm text-muted-foreground">Staked</p>
                          <p className="text-xs text-green-600">
                            +{formatCurrency(yearlyProjection)}/year
                          </p>
                        </div>
                      </div>
                    </motion.div>
                  )
                })}
              </div>
            </CardContent>
          </Card>

          {/* Optimization Results */}
          <Card>
            <CardHeader>
              <CardTitle>Optimization Analysis</CardTitle>
              <CardDescription>
                Recommendations to maximize your yield
              </CardDescription>
            </CardHeader>
            <CardContent>
              {!optimization ? (
                <div className="text-center py-8">
                  <BarChart3 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">Select a Position</h3>
                  <p className="text-muted-foreground">
                    Choose a position from the left to see optimization recommendations
                  </p>
                </div>
              ) : (
                <div className="space-y-6">
                  {/* Current Position */}
                  <div>
                    <h4 className="font-medium mb-3">Current Position</h4>
                    <div className="border rounded-lg p-4 bg-muted/50">
                      <div className="flex items-center justify-between">
                        <div>
                          <h5 className="font-medium">{optimization.currentFarm.name}</h5>
                          <p className="text-sm text-muted-foreground">{optimization.currentFarm.protocol}</p>
                        </div>
                        <div className="text-right">
                          <p className="font-bold text-lg">{optimization.currentFarm.apy}</p>
                          <p className="text-sm text-muted-foreground">Current APY</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Recommendation */}
                  <div>
                    <h4 className="font-medium mb-3">Recommendation</h4>
                    <Alert className={getRecommendationColor(optimization.recommendation)}>
                      <div className="flex items-center gap-2">
                        {getRecommendationIcon(optimization.recommendation)}
                        <AlertDescription>
                          <strong className="capitalize">{optimization.recommendation}:</strong> {optimization.reason}
                        </AlertDescription>
                      </div>
                    </Alert>
                  </div>

                  {/* Potential Gains */}
                  {optimization.recommendation !== 'stay' && (
                    <div>
                      <h4 className="font-medium mb-3">Potential Impact</h4>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="border rounded-lg p-3">
                          <div className="flex items-center gap-2 mb-1">
                            <TrendingUp className="w-4 h-4 text-green-500" />
                            <span className="text-sm font-medium">Potential Gains</span>
                          </div>
                          <p className="text-lg font-bold text-green-600">{optimization.potentialGains}</p>
                        </div>
                        <div className="border rounded-lg p-3">
                          <div className="flex items-center gap-2 mb-1">
                            <DollarSign className="w-4 h-4 text-blue-500" />
                            <span className="text-sm font-medium">Migration Cost</span>
                          </div>
                          <p className="text-lg font-bold">{optimization.migrationCost}</p>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Suggested Farms */}
                  {optimization.suggestedFarms.length > 0 && (
                    <div>
                      <h4 className="font-medium mb-3">Suggested Alternatives</h4>
                      <div className="space-y-3">
                        {optimization.suggestedFarms.map((farm, index) => {
                          const currentPosition = userPositions.find((p: any) => p.id === selectedPosition)
                          const currentProjection = currentPosition 
                            ? calculateYieldProjection(optimization.currentFarm.apy, currentPosition.stakedAmount)
                            : 0
                          const newProjection = currentPosition 
                            ? calculateYieldProjection(farm.apy, currentPosition.stakedAmount)
                            : 0
                          const additionalYield = newProjection - currentProjection

                          return (
                            <motion.div
                              key={farm.id}
                              initial={{ opacity: 0, x: 20 }}
                              animate={{ opacity: 1, x: 0 }}
                              transition={{ delay: index * 0.1 }}
                              className="border rounded-lg p-4 hover:bg-accent/50 transition-colors"
                            >
                              <div className="flex items-center justify-between">
                                <div>
                                  <h5 className="font-medium">{farm.name}</h5>
                                  <p className="text-sm text-muted-foreground">{farm.protocol}</p>
                                  <div className="flex items-center gap-2 mt-1">
                                    <Badge className={getRiskColor(farm.riskLevel)}>
                                      {farm.riskLevel}
                                    </Badge>
                                    <Badge variant="outline" className="text-xs">
                                      {farm.tvl} TVL
                                    </Badge>
                                  </div>
                                </div>
                                <div className="text-right">
                                  <p className="font-bold text-lg text-green-600">{farm.apy}</p>
                                  <p className="text-sm text-muted-foreground">APY</p>
                                  <p className="text-xs text-green-600">
                                    +{formatCurrency(additionalYield)}/year
                                  </p>
                                </div>
                              </div>

                              <div className="mt-3 pt-3 border-t">
                                <div className="flex items-center justify-between text-sm">
                                  <span className="text-muted-foreground">Min Stake:</span>
                                  <span>{farm.minimumStake} {farm.stakingToken.symbol}</span>
                                </div>
                                {farm.lockPeriod && (
                                  <div className="flex items-center justify-between text-sm mt-1">
                                    <span className="text-muted-foreground">Lock Period:</span>
                                    <span>{Math.floor(farm.lockPeriod / 86400)} days</span>
                                  </div>
                                )}
                              </div>

                              <Button className="w-full mt-3" size="sm">
                                <ArrowRight className="w-3 h-3 mr-2" />
                                Migrate to {farm.protocol}
                              </Button>
                            </motion.div>
                          )
                        })}
                      </div>
                    </div>
                  )}

                  {/* Risk Analysis */}
                  <div>
                    <h4 className="font-medium mb-3">Risk Analysis</h4>
                    <div className="space-y-3">
                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <Shield className="w-4 h-4" />
                          <span className="text-sm">Current Risk Level</span>
                        </div>
                        <Badge className={getRiskColor(optimization.currentFarm.riskLevel)}>
                          {optimization.currentFarm.riskLevel}
                        </Badge>
                      </div>

                      {optimization.suggestedFarms.length > 0 && (
                        <div className="flex items-center justify-between p-3 border rounded-lg">
                          <div className="flex items-center gap-2">
                            <AlertTriangle className="w-4 h-4" />
                            <span className="text-sm">Suggested Risk Range</span>
                          </div>
                          <div className="flex gap-1">
                            {Array.from(new Set(optimization.suggestedFarms.map(f => f.riskLevel))).map(risk => (
                              <Badge key={risk} className={getRiskColor(risk)}>
                                {risk}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      )}

                      <div className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex items-center gap-2">
                          <Clock className="w-4 h-4" />
                          <span className="text-sm">Time Horizon</span>
                        </div>
                        <span className="text-sm font-medium">Long-term (1+ year)</span>
                      </div>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-2">
                    <Button className="flex-1">
                      <RefreshCw className="w-4 h-4 mr-2" />
                      Refresh Analysis
                    </Button>
                    {optimization.recommendation === 'migrate' && (
                      <Button variant="outline" className="flex-1">
                        <Zap className="w-4 h-4 mr-2" />
                        Auto-Migrate
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
