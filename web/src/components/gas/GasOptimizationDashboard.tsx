'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Fuel,
  TrendingUp,
  TrendingDown,
  Clock,
  Shield,
  Zap,
  Target,
  AlertTriangle,
  CheckCircle,
  RefreshCw,
  Settings,
  Lightbulb,
  BarChart3,
  Timer
} from 'lucide-react'
import { useChainId } from 'wagmi'
import { useGasOptimization } from '@/hooks/useGasOptimization'
import { GasPriority, type GasOptimizationSuggestion } from '@/lib/gas-optimization'
import { cn } from '@/lib/utils'

export function GasOptimizationDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [gasLimit, setGasLimit] = useState(BigInt(21000))
  const [transactionType, setTransactionType] = useState('transfer')
  const [amount, setAmount] = useState('')

  const chainId = useChainId()
  
  const {
    gasEstimates,
    gasTracker,
    suggestions,
    isLoading,
    error,
    getGasEstimates,
    generateSuggestions,
    getOptimizationAnalysis,
    applySuggestion,
    getRecommendedPriority,
    getEstimateForPriority,
    calculateSavings,
    formatGasPrice,
    getCongestionLevel,
    getMEVRisk,
    getTimeBasedRecommendation,
    isOptimalTime,
    shouldWait,
    gasTrend
  } = useGasOptimization({
    chainId: chainId || 1,
    autoUpdate: true,
    enableNotifications: true
  })

  const congestionLevel = getCongestionLevel()
  const recommendedPriority = getRecommendedPriority()
  const mevRisk = getMEVRisk(transactionType)

  useEffect(() => {
    if (chainId) {
      getGasEstimates(gasLimit)
      generateSuggestions(transactionType, amount)
    }
  }, [chainId, gasLimit, transactionType, amount, getGasEstimates, generateSuggestions])

  const getPriorityIcon = (priority: GasPriority) => {
    switch (priority) {
      case GasPriority.SLOW:
        return <Clock className="w-4 h-4 text-blue-500" />
      case GasPriority.STANDARD:
        return <Target className="w-4 h-4 text-green-500" />
      case GasPriority.FAST:
        return <Zap className="w-4 h-4 text-orange-500" />
      case GasPriority.INSTANT:
        return <TrendingUp className="w-4 h-4 text-red-500" />
      default:
        return <Target className="w-4 h-4 text-gray-500" />
    }
  }

  const getPriorityColor = (priority: GasPriority) => {
    switch (priority) {
      case GasPriority.SLOW:
        return 'border-blue-200 bg-blue-50 dark:border-blue-800 dark:bg-blue-950'
      case GasPriority.STANDARD:
        return 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950'
      case GasPriority.FAST:
        return 'border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950'
      case GasPriority.INSTANT:
        return 'border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950'
      default:
        return 'border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-950'
    }
  }

  const getSuggestionIcon = (type: GasOptimizationSuggestion['type']) => {
    switch (type) {
      case 'timing':
        return <Timer className="w-4 h-4" />
      case 'batching':
        return <BarChart3 className="w-4 h-4" />
      case 'route':
        return <Target className="w-4 h-4" />
      case 'mev':
        return <Shield className="w-4 h-4" />
      case 'layer2':
        return <Zap className="w-4 h-4" />
      default:
        return <Lightbulb className="w-4 h-4" />
    }
  }

  const formatTime = (seconds: number) => {
    if (seconds < 60) return `~${seconds}s`
    if (seconds < 3600) return `~${Math.round(seconds / 60)}m`
    return `~${Math.round(seconds / 3600)}h`
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Gas Optimization</h2>
          <p className="text-muted-foreground">
            Optimize your transaction costs with intelligent gas management
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => getGasEstimates(gasLimit)}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Network Status */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Network Status</p>
                <p className={cn("text-lg font-bold", congestionLevel.color)}>
                  {congestionLevel.level.toUpperCase()}
                </p>
              </div>
              <Fuel className="w-8 h-8 text-blue-500" />
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              {congestionLevel.description}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Current Gas</p>
                <p className="text-lg font-bold">
                  {gasTracker ? formatGasPrice(gasTracker.currentGasPrice) : 'Loading...'}
                </p>
              </div>
              {gasTrend === 'rising' ? (
                <TrendingUp className="w-8 h-8 text-red-500" />
              ) : gasTrend === 'falling' ? (
                <TrendingDown className="w-8 h-8 text-green-500" />
              ) : (
                <Target className="w-8 h-8 text-gray-500" />
              )}
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              Trend: {gasTrend}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Recommendation</p>
                <p className="text-lg font-bold">
                  {isOptimalTime ? 'TRANSACT NOW' : shouldWait ? 'WAIT' : 'PROCEED'}
                </p>
              </div>
              {isOptimalTime ? (
                <CheckCircle className="w-8 h-8 text-green-500" />
              ) : shouldWait ? (
                <Clock className="w-8 h-8 text-yellow-500" />
              ) : (
                <Target className="w-8 h-8 text-blue-500" />
              )}
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              Based on current conditions
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">MEV Risk</p>
                <p className={cn("text-lg font-bold", 
                  mevRisk === 'high' ? 'text-red-600' :
                  mevRisk === 'medium' ? 'text-yellow-600' : 'text-green-600'
                )}>
                  {mevRisk.toUpperCase()}
                </p>
              </div>
              <Shield className={cn("w-8 h-8",
                mevRisk === 'high' ? 'text-red-500' :
                mevRisk === 'medium' ? 'text-yellow-500' : 'text-green-500'
              )} />
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              Frontrunning protection
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Gas Estimates</TabsTrigger>
          <TabsTrigger value="optimization">Optimization</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Gas Estimates */}
          <Card>
            <CardHeader>
              <CardTitle>Gas Price Estimates</CardTitle>
              <CardDescription>
                Choose the right gas price for your transaction priority
              </CardDescription>
            </CardHeader>
            <CardContent>
              {gasEstimates.length > 0 ? (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                  {gasEstimates.map((estimate, index) => {
                    const isRecommended = estimate.priority === recommendedPriority
                    const savings = calculateSavings(GasPriority.FAST, estimate.priority)
                    
                    return (
                      <motion.div
                        key={estimate.priority}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                      >
                        <Card className={cn(
                          "transition-all duration-200 cursor-pointer hover:shadow-md",
                          getPriorityColor(estimate.priority),
                          isRecommended && "ring-2 ring-primary"
                        )}>
                          <CardContent className="p-4">
                            <div className="flex items-center justify-between mb-3">
                              <div className="flex items-center gap-2">
                                {getPriorityIcon(estimate.priority)}
                                <span className="font-medium capitalize">
                                  {estimate.priority}
                                </span>
                              </div>
                              {isRecommended && (
                                <Badge variant="default" className="text-xs">
                                  Recommended
                                </Badge>
                              )}
                            </div>
                            
                            <div className="space-y-2">
                              <div>
                                <p className="text-sm text-muted-foreground">Gas Price</p>
                                <p className="text-lg font-bold">
                                  {formatGasPrice(estimate.gasPrice)}
                                </p>
                              </div>
                              
                              <div>
                                <p className="text-sm text-muted-foreground">Cost</p>
                                <p className="font-medium">${estimate.cost}</p>
                              </div>
                              
                              <div>
                                <p className="text-sm text-muted-foreground">Time</p>
                                <p className="font-medium">{formatTime(estimate.estimatedTime)}</p>
                              </div>

                              {savings && savings.percentage > 0 && (
                                <div className="pt-2 border-t">
                                  <p className="text-xs text-green-600">
                                    Save {savings.percentage}% vs Fast
                                  </p>
                                </div>
                              )}
                            </div>
                          </CardContent>
                        </Card>
                      </motion.div>
                    )
                  })}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Fuel className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">
                    {isLoading ? 'Loading gas estimates...' : 'No gas estimates available'}
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="optimization" className="space-y-6">
          {/* Optimization Suggestions */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Lightbulb className="w-5 h-5" />
                Optimization Suggestions
              </CardTitle>
              <CardDescription>
                Smart recommendations to reduce your gas costs
              </CardDescription>
            </CardHeader>
            <CardContent>
              {suggestions.length > 0 ? (
                <div className="space-y-4">
                  {suggestions.map((suggestion, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-start justify-between p-4 border rounded-lg"
                    >
                      <div className="flex items-start gap-3">
                        <div className="mt-1">
                          {getSuggestionIcon(suggestion.type)}
                        </div>
                        <div>
                          <h4 className="font-medium">{suggestion.title}</h4>
                          <p className="text-sm text-muted-foreground mt-1">
                            {suggestion.description}
                          </p>
                          <div className="flex items-center gap-4 mt-2">
                            <Badge variant="outline" className="text-xs">
                              Save {suggestion.potentialSavings}
                            </Badge>
                            <Badge 
                              variant={suggestion.difficulty === 'easy' ? 'default' : 
                                     suggestion.difficulty === 'medium' ? 'secondary' : 'destructive'}
                              className="text-xs"
                            >
                              {suggestion.difficulty}
                            </Badge>
                          </div>
                        </div>
                      </div>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => applySuggestion(suggestion)}
                      >
                        Apply
                      </Button>
                    </motion.div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Lightbulb className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">
                    No optimization suggestions available
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="analytics" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Gas Analytics</CardTitle>
              <CardDescription>
                Historical gas price trends and network analysis
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-center py-12">
                <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">Analytics Coming Soon</h3>
                <p className="text-muted-foreground">
                  Advanced gas analytics and historical trends will be available here
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
