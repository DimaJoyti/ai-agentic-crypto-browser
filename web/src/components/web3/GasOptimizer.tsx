'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Zap, 
  TrendingUp, 
  TrendingDown, 
  Minus,
  Clock, 
  DollarSign,
  AlertTriangle,
  CheckCircle,
  Lightbulb,
  RefreshCw,
  Settings,
  Target,
  Timer,
  Gauge,
  Flame,
  Snowflake,
  Activity
} from 'lucide-react'
import { useGasOptimization } from '@/hooks/useGasOptimization'
import { GasPriority, type GasEstimate } from '@/lib/gas-optimization'
import { SUPPORTED_CHAINS } from '@/lib/chains'
import { toast } from 'sonner'

interface GasOptimizerProps {
  chainId: number
  onChainChange?: (chainId: number) => void
}

export function GasOptimizer({ chainId, onChainChange }: GasOptimizerProps) {
  const [gasLimit, setGasLimit] = useState('21000')
  const [transactionType, setTransactionType] = useState('send')
  const [amount, setAmount] = useState('')
  const [selectedPriority, setSelectedPriority] = useState<GasPriority>(GasPriority.STANDARD)

  const {
    gasEstimates,
    gasTracker,
    suggestions,
    isLoading,
    error,
    lastUpdated,
    getGasEstimates,
    updateGasTracker,
    generateSuggestions,
    getRecommendedPriority,
    getEstimateForPriority,
    calculateSavings,
    formatGasPrice,
    getCongestionLevel,
    isOptimalTime,
    shouldWait,
    gasTrend,
    formattedGasPrice
  } = useGasOptimization({
    chainId,
    autoUpdate: true,
    enableNotifications: true
  })

  const chain = SUPPORTED_CHAINS[chainId]
  const congestion = getCongestionLevel()
  const recommendedPriority = getRecommendedPriority()

  useEffect(() => {
    if (gasLimit) {
      getGasEstimates(BigInt(gasLimit))
    }
  }, [gasLimit, getGasEstimates])

  useEffect(() => {
    generateSuggestions(transactionType, amount)
  }, [transactionType, amount, generateSuggestions])

  const handleRefresh = () => {
    updateGasTracker()
    if (gasLimit) {
      getGasEstimates(BigInt(gasLimit))
    }
    toast.success('Gas data refreshed')
  }

  const getPriorityIcon = (priority: GasPriority) => {
    switch (priority) {
      case GasPriority.SLOW:
        return <Snowflake className="w-4 h-4 text-blue-500" />
      case GasPriority.STANDARD:
        return <Activity className="w-4 h-4 text-green-500" />
      case GasPriority.FAST:
        return <Zap className="w-4 h-4 text-yellow-500" />
      case GasPriority.INSTANT:
        return <Flame className="w-4 h-4 text-red-500" />
    }
  }

  const getPriorityColor = (priority: GasPriority) => {
    switch (priority) {
      case GasPriority.SLOW:
        return 'border-blue-200 bg-blue-50 dark:bg-blue-950'
      case GasPriority.STANDARD:
        return 'border-green-200 bg-green-50 dark:bg-green-950'
      case GasPriority.FAST:
        return 'border-yellow-200 bg-yellow-50 dark:bg-yellow-950'
      case GasPriority.INSTANT:
        return 'border-red-200 bg-red-50 dark:bg-red-950'
    }
  }

  const getTrendIcon = () => {
    switch (gasTrend) {
      case 'rising':
        return <TrendingUp className="w-4 h-4 text-red-500" />
      case 'falling':
        return <TrendingDown className="w-4 h-4 text-green-500" />
      default:
        return <Minus className="w-4 h-4 text-gray-500" />
    }
  }

  const getSuggestionIcon = (type: string) => {
    switch (type) {
      case 'timing':
        return <Timer className="w-4 h-4" />
      case 'batching':
        return <Target className="w-4 h-4" />
      case 'route':
        return <Activity className="w-4 h-4" />
      default:
        return <Lightbulb className="w-4 h-4" />
    }
  }

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy':
        return 'bg-green-100 text-green-800'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800'
      case 'hard':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Zap className="w-6 h-6" />
            Gas Optimizer
          </h2>
          <p className="text-muted-foreground">
            Optimize transaction costs with intelligent gas price analysis
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="gap-1">
            {chain?.shortName || 'Unknown'}
          </Badge>
          <Button variant="outline" size="sm" onClick={handleRefresh}>
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Current Gas Status */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Current Gas Price</p>
                <p className="text-2xl font-bold">{formattedGasPrice}</p>
                <div className="flex items-center gap-1 mt-1">
                  {getTrendIcon()}
                  <span className="text-xs text-muted-foreground capitalize">
                    {gasTrend}
                  </span>
                </div>
              </div>
              <Gauge className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Network Congestion</p>
                <p className={`text-2xl font-bold capitalize ${congestion.color}`}>
                  {congestion.level}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  {congestion.description}
                </p>
              </div>
              <Activity className={`w-8 h-8 ${congestion.color}`} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Recommendation</p>
                <p className="text-2xl font-bold">
                  {isOptimalTime ? 'Act Now' : shouldWait ? 'Wait' : 'Proceed'}
                </p>
                <div className="flex items-center gap-1 mt-1">
                  {isOptimalTime ? (
                    <CheckCircle className="w-3 h-3 text-green-500" />
                  ) : shouldWait ? (
                    <AlertTriangle className="w-3 h-3 text-yellow-500" />
                  ) : (
                    <Clock className="w-3 h-3 text-blue-500" />
                  )}
                  <span className="text-xs text-muted-foreground">
                    {isOptimalTime ? 'Optimal time' : shouldWait ? 'High fees' : 'Normal time'}
                  </span>
                </div>
              </div>
              {isOptimalTime ? (
                <CheckCircle className="w-8 h-8 text-green-500" />
              ) : shouldWait ? (
                <AlertTriangle className="w-8 h-8 text-yellow-500" />
              ) : (
                <Clock className="w-8 h-8 text-blue-500" />
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="estimates" className="space-y-6">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="estimates">Gas Estimates</TabsTrigger>
          <TabsTrigger value="suggestions">Optimization Tips</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="estimates" className="space-y-6">
          {/* Transaction Settings */}
          <Card>
            <CardHeader>
              <CardTitle>Transaction Settings</CardTitle>
              <CardDescription>
                Configure your transaction to get accurate gas estimates
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label>Gas Limit</Label>
                  <Input
                    type="number"
                    value={gasLimit}
                    onChange={(e) => setGasLimit(e.target.value)}
                    placeholder="21000"
                  />
                </div>
                <div className="space-y-2">
                  <Label>Transaction Type</Label>
                  <select
                    value={transactionType}
                    onChange={(e) => setTransactionType(e.target.value)}
                    className="w-full p-2 border rounded-md"
                  >
                    <option value="send">Send</option>
                    <option value="swap">Swap</option>
                    <option value="approve">Approve</option>
                    <option value="stake">Stake</option>
                    <option value="multiple">Multiple</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <Label>Amount (ETH)</Label>
                  <Input
                    type="number"
                    step="0.001"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    placeholder="0.1"
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Gas Estimates */}
          <Card>
            <CardHeader>
              <CardTitle>Gas Price Options</CardTitle>
              <CardDescription>
                Choose the right gas price for your transaction timing needs
              </CardDescription>
            </CardHeader>
            <CardContent>
              {error ? (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              ) : (
                <div className="space-y-3">
                  {gasEstimates.map((estimate) => {
                    const isRecommended = estimate.priority === recommendedPriority
                    const savings = selectedPriority !== estimate.priority 
                      ? calculateSavings(selectedPriority, estimate.priority)
                      : null

                    return (
                      <motion.div
                        key={estimate.priority}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className={`border rounded-lg p-4 cursor-pointer transition-all ${
                          selectedPriority === estimate.priority
                            ? 'ring-2 ring-blue-500 ' + getPriorityColor(estimate.priority)
                            : 'hover:bg-accent/50'
                        } ${isRecommended ? 'ring-2 ring-green-500' : ''}`}
                        onClick={() => setSelectedPriority(estimate.priority)}
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            {getPriorityIcon(estimate.priority)}
                            <div>
                              <div className="flex items-center gap-2">
                                <h4 className="font-medium capitalize">{estimate.priority}</h4>
                                {isRecommended && (
                                  <Badge variant="secondary" className="text-xs">
                                    Recommended
                                  </Badge>
                                )}
                              </div>
                              <p className="text-sm text-muted-foreground">
                                ~{estimate.estimatedTime < 60 
                                  ? `${estimate.estimatedTime}s` 
                                  : `${Math.round(estimate.estimatedTime / 60)}m`
                                }
                              </p>
                            </div>
                          </div>
                          
                          <div className="text-right">
                            <p className="font-medium">{estimate.cost} ETH</p>
                            <p className="text-sm text-muted-foreground">
                              {formatGasPrice(estimate.gasPrice)}
                            </p>
                            {savings && (
                              <p className={`text-xs ${savings.percentage > 0 ? 'text-green-600' : 'text-red-600'}`}>
                                {savings.percentage > 0 ? '-' : '+'}{Math.abs(savings.percentage)}% 
                              </p>
                            )}
                          </div>
                        </div>
                        
                        <div className="mt-3">
                          <div className="flex items-center justify-between text-xs text-muted-foreground mb-1">
                            <span>Confidence</span>
                            <span>{estimate.confidence}%</span>
                          </div>
                          <Progress value={estimate.confidence} className="h-1" />
                        </div>
                      </motion.div>
                    )
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="suggestions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Optimization Suggestions</CardTitle>
              <CardDescription>
                Smart recommendations to reduce your transaction costs
              </CardDescription>
            </CardHeader>
            <CardContent>
              {suggestions.length === 0 ? (
                <div className="text-center py-8">
                  <Lightbulb className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No suggestions available</h3>
                  <p className="text-muted-foreground">
                    Configure your transaction settings to get personalized optimization tips
                  </p>
                </div>
              ) : (
                <div className="space-y-4">
                  {suggestions.map((suggestion, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="border rounded-lg p-4"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex items-start gap-3 flex-1">
                          {getSuggestionIcon(suggestion.type)}
                          <div className="flex-1">
                            <h4 className="font-medium">{suggestion.title}</h4>
                            <p className="text-sm text-muted-foreground mt-1">
                              {suggestion.description}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2 ml-4">
                          <Badge variant="outline" className="text-green-600">
                            Save {suggestion.potentialSavings}
                          </Badge>
                          <Badge 
                            variant="secondary" 
                            className={getDifficultyColor(suggestion.difficulty)}
                          >
                            {suggestion.difficulty}
                          </Badge>
                        </div>
                      </div>
                      {suggestion.action && (
                        <div className="mt-3">
                          <Button size="sm" onClick={suggestion.action}>
                            Apply Suggestion
                          </Button>
                        </div>
                      )}
                    </motion.div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Gas Optimization Settings</CardTitle>
              <CardDescription>
                Configure how gas optimization works for your transactions
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label>Auto-refresh gas prices</Label>
                    <p className="text-sm text-muted-foreground">
                      Automatically update gas prices every 15 seconds
                    </p>
                  </div>
                  <input type="checkbox" defaultChecked className="toggle" />
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label>Gas price notifications</Label>
                    <p className="text-sm text-muted-foreground">
                      Get notified when gas prices change significantly
                    </p>
                  </div>
                  <input type="checkbox" defaultChecked className="toggle" />
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label>Optimization suggestions</Label>
                    <p className="text-sm text-muted-foreground">
                      Show smart suggestions to reduce transaction costs
                    </p>
                  </div>
                  <input type="checkbox" defaultChecked className="toggle" />
                </div>
              </div>

              <div className="pt-4 border-t">
                <h4 className="font-medium mb-2">Default Gas Preferences</h4>
                <div className="space-y-3">
                  <div>
                    <Label>Preferred gas priority</Label>
                    <select className="w-full p-2 border rounded-md mt-1">
                      <option value="standard">Standard (Recommended)</option>
                      <option value="slow">Slow (Cheapest)</option>
                      <option value="fast">Fast</option>
                      <option value="instant">Instant (Most Expensive)</option>
                    </select>
                  </div>
                  
                  <div>
                    <Label>Max gas price (gwei)</Label>
                    <Input type="number" placeholder="100" className="mt-1" />
                    <p className="text-xs text-muted-foreground mt-1">
                      Transactions will be rejected if gas price exceeds this limit
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
