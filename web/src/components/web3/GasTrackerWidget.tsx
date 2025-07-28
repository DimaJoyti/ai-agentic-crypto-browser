'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  Zap, 
  TrendingUp, 
  TrendingDown, 
  Minus,
  ChevronDown,
  ChevronUp,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  Clock,
  Gauge
} from 'lucide-react'
import { useGasOptimization } from '@/hooks/useGasOptimization'
import { SUPPORTED_CHAINS } from '@/lib/chains'

interface GasTrackerWidgetProps {
  chainId: number
  compact?: boolean
}

export function GasTrackerWidget({ chainId, compact = false }: GasTrackerWidgetProps) {
  const [isExpanded, setIsExpanded] = useState(false)

  const {
    gasTracker,
    isLoading,
    updateGasTracker,
    formatGasPrice,
    getCongestionLevel,
    isOptimalTime,
    shouldWait,
    gasTrend,
    formattedGasPrice
  } = useGasOptimization({
    chainId,
    autoUpdate: true,
    enableNotifications: false
  })

  const chain = SUPPORTED_CHAINS[chainId]
  const congestion = getCongestionLevel()

  const getTrendIcon = () => {
    switch (gasTrend) {
      case 'rising':
        return <TrendingUp className="w-3 h-3 text-red-500" />
      case 'falling':
        return <TrendingDown className="w-3 h-3 text-green-500" />
      default:
        return <Minus className="w-3 h-3 text-gray-500" />
    }
  }

  const getStatusIcon = () => {
    if (isOptimalTime) {
      return <CheckCircle className="w-4 h-4 text-green-500" />
    } else if (shouldWait) {
      return <AlertTriangle className="w-4 h-4 text-yellow-500" />
    } else {
      return <Clock className="w-4 h-4 text-blue-500" />
    }
  }

  const getStatusText = () => {
    if (isOptimalTime) return 'Optimal'
    if (shouldWait) return 'Wait'
    return 'Normal'
  }

  const getStatusColor = () => {
    if (isOptimalTime) return 'text-green-600'
    if (shouldWait) return 'text-yellow-600'
    return 'text-blue-600'
  }

  if (compact) {
    return (
      <Card className="w-full">
        <CardContent className="p-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Zap className="w-4 h-4 text-blue-500" />
              <div>
                <p className="text-sm font-medium">{formattedGasPrice}</p>
                <div className="flex items-center gap-1">
                  {getTrendIcon()}
                  <span className="text-xs text-muted-foreground capitalize">
                    {gasTrend}
                  </span>
                </div>
              </div>
            </div>
            
            <div className="flex items-center gap-2">
              <Badge 
                variant="outline" 
                className={`text-xs ${congestion.color}`}
              >
                {congestion.level}
              </Badge>
              {getStatusIcon()}
            </div>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="w-full">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg flex items-center gap-2">
            <Gauge className="w-5 h-5" />
            Gas Tracker
          </CardTitle>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="text-xs">
              {chain?.shortName || 'Unknown'}
            </Badge>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsExpanded(!isExpanded)}
            >
              {isExpanded ? (
                <ChevronUp className="w-4 h-4" />
              ) : (
                <ChevronDown className="w-4 h-4" />
              )}
            </Button>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Current Gas Price */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">Current Price</p>
            <p className="text-xl font-bold">{formattedGasPrice}</p>
          </div>
          <div className="flex items-center gap-1">
            {getTrendIcon()}
            <span className="text-sm text-muted-foreground capitalize">
              {gasTrend}
            </span>
          </div>
        </div>

        {/* Network Status */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">Network</p>
            <p className={`font-medium capitalize ${congestion.color}`}>
              {congestion.level} congestion
            </p>
          </div>
          <div className="flex items-center gap-2">
            {getStatusIcon()}
            <span className={`text-sm font-medium ${getStatusColor()}`}>
              {getStatusText()}
            </span>
          </div>
        </div>

        <AnimatePresence>
          {isExpanded && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="space-y-4 border-t pt-4"
            >
              {/* Detailed Information */}
              <div className="space-y-3">
                <div>
                  <p className="text-sm text-muted-foreground">Recommendation</p>
                  <p className="text-sm">
                    {isOptimalTime 
                      ? 'Excellent time to transact - gas prices are low!'
                      : shouldWait 
                      ? 'Consider waiting for lower gas prices'
                      : 'Normal time to proceed with transactions'
                    }
                  </p>
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">Network Health</p>
                  <p className="text-sm">{congestion.description}</p>
                </div>

                {gasTracker && (
                  <div>
                    <p className="text-sm text-muted-foreground">Next Update</p>
                    <p className="text-sm">{gasTracker.nextUpdateIn}s</p>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={updateGasTracker}
                  disabled={isLoading}
                  className="flex-1"
                >
                  <RefreshCw className={`w-3 h-3 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
                
                {isOptimalTime && (
                  <Button size="sm" className="flex-1">
                    <Zap className="w-3 h-3 mr-2" />
                    Transact Now
                  </Button>
                )}
              </div>

              {/* Quick Gas Estimates */}
              <div className="space-y-2">
                <p className="text-sm font-medium">Quick Estimates</p>
                <div className="grid grid-cols-2 gap-2 text-xs">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Slow:</span>
                    <span>{gasTracker ? formatGasPrice(gasTracker.currentGasPrice * BigInt(80) / BigInt(100)) : '-'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Fast:</span>
                    <span>{gasTracker ? formatGasPrice(gasTracker.currentGasPrice * BigInt(120) / BigInt(100)) : '-'}</span>
                  </div>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </CardContent>
    </Card>
  )
}
