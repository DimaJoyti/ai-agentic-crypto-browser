'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  TrendingUp, 
  TrendingDown, 
  Minus,
  RefreshCw,
  Bell,
  BellOff,
  Settings,
  Activity,
  Zap
} from 'lucide-react'
import { usePriceFeed, usePrices } from '@/hooks/usePriceFeed'
import { type PriceData } from '@/lib/price-feed-manager'
import { cn } from '@/lib/utils'

interface PriceTickerProps {
  symbols?: string[]
  showChange?: boolean
  showVolume?: boolean
  showMarketCap?: boolean
  compact?: boolean
  autoRefresh?: boolean
  className?: string
}

export function PriceTicker({
  symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'],
  showChange = true,
  showVolume = false,
  showMarketCap = false,
  compact = false,
  autoRefresh = true,
  className
}: PriceTickerProps) {
  const [selectedSymbol, setSelectedSymbol] = useState<string | null>(null)
  const [alertsEnabled, setAlertsEnabled] = useState(false)

  const { prices, isLoading } = usePrices(symbols)
  const { 
    state, 
    restart, 
    addAlert, 
    getAlertsForSymbol,
    formatPrice, 
    formatChange 
  } = usePriceFeed({
    symbols,
    autoStart: autoRefresh,
    enableAlerts: alertsEnabled
  })

  const getPriceChangeColor = (changePercent: number) => {
    if (changePercent > 0) return 'text-green-600 dark:text-green-400'
    if (changePercent < 0) return 'text-red-600 dark:text-red-400'
    return 'text-gray-600 dark:text-gray-400'
  }

  const getPriceChangeIcon = (changePercent: number) => {
    if (changePercent > 0) return <TrendingUp className="w-4 h-4" />
    if (changePercent < 0) return <TrendingDown className="w-4 h-4" />
    return <Minus className="w-4 h-4" />
  }

  const formatVolume = (volume: number) => {
    if (volume >= 1e9) return `$${(volume / 1e9).toFixed(2)}B`
    if (volume >= 1e6) return `$${(volume / 1e6).toFixed(2)}M`
    if (volume >= 1e3) return `$${(volume / 1e3).toFixed(2)}K`
    return `$${volume.toFixed(2)}`
  }

  const formatMarketCap = (marketCap: number) => {
    if (marketCap >= 1e12) return `$${(marketCap / 1e12).toFixed(2)}T`
    if (marketCap >= 1e9) return `$${(marketCap / 1e9).toFixed(2)}B`
    if (marketCap >= 1e6) return `$${(marketCap / 1e6).toFixed(2)}M`
    return `$${marketCap.toFixed(2)}`
  }

  const handleAddAlert = (symbol: string, price: number) => {
    // Add a simple price alert 5% above current price
    addAlert({
      symbol,
      type: 'above',
      value: price * 1.05,
      isActive: true,
      triggered: false
    })
  }

  const getSymbolAlerts = (symbol: string) => {
    return getAlertsForSymbol(symbol).filter(alert => alert.isActive && !alert.triggered)
  }

  if (isLoading && prices.size === 0) {
    return (
      <Card className={cn("w-full", className)}>
        <CardContent className="p-4">
          <div className="flex items-center justify-center space-x-2">
            <RefreshCw className="w-4 h-4 animate-spin" />
            <span className="text-sm text-muted-foreground">Loading prices...</span>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (compact) {
    return (
      <div className={cn("flex items-center space-x-4 overflow-x-auto", className)}>
        <AnimatePresence>
          {Array.from(prices.entries()).map(([symbol, price]) => (
            <motion.div
              key={symbol}
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.9 }}
              className="flex items-center space-x-2 whitespace-nowrap"
            >
              <span className="font-medium">{symbol}</span>
              <span className="text-sm">{formatPrice(price.price)}</span>
              {showChange && (
                <span className={cn("text-xs flex items-center", getPriceChangeColor(price.changePercent24h))}>
                  {getPriceChangeIcon(price.changePercent24h)}
                  {formatChange(price.changePercent24h, true)}
                </span>
              )}
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    )
  }

  return (
    <Card className={cn("w-full", className)}>
      <CardContent className="p-4">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-2">
            <Activity className="w-5 h-5" />
            <h3 className="text-lg font-semibold">Live Prices</h3>
            {state.isConnected && (
              <Badge variant="default" className="text-xs">
                <Zap className="w-3 h-3 mr-1" />
                Live
              </Badge>
            )}
          </div>
          
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setAlertsEnabled(!alertsEnabled)}
            >
              {alertsEnabled ? <Bell className="w-4 h-4" /> : <BellOff className="w-4 h-4" />}
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={restart}
              disabled={state.isLoading}
            >
              <RefreshCw className={cn("w-4 h-4", state.isLoading && "animate-spin")} />
            </Button>
          </div>
        </div>

        <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          <AnimatePresence>
            {Array.from(prices.entries()).map(([symbol, price]) => {
              const alerts = getSymbolAlerts(symbol)
              const isSelected = selectedSymbol === symbol

              return (
                <motion.div
                  key={symbol}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  whileHover={{ scale: 1.02 }}
                  className={cn(
                    "p-3 border rounded-lg cursor-pointer transition-all",
                    isSelected && "ring-2 ring-primary",
                    "hover:shadow-md"
                  )}
                  onClick={() => setSelectedSymbol(isSelected ? null : symbol)}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center space-x-2">
                      <span className="font-bold text-lg">{symbol}</span>
                      {alerts.length > 0 && (
                        <Badge variant="secondary" className="text-xs">
                          <Bell className="w-3 h-3 mr-1" />
                          {alerts.length}
                        </Badge>
                      )}
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleAddAlert(symbol, price.price)
                      }}
                      className="h-6 w-6 p-0"
                    >
                      <Bell className="w-3 h-3" />
                    </Button>
                  </div>

                  <div className="space-y-1">
                    <div className="text-xl font-bold">
                      {formatPrice(price.price)}
                    </div>

                    {showChange && (
                      <div className={cn("flex items-center space-x-1 text-sm", getPriceChangeColor(price.changePercent24h))}>
                        {getPriceChangeIcon(price.changePercent24h)}
                        <span>{formatChange(price.change24h)}</span>
                        <span>({formatChange(price.changePercent24h, true)})</span>
                      </div>
                    )}

                    {(showVolume || showMarketCap) && (
                      <div className="text-xs text-muted-foreground space-y-1">
                        {showVolume && (
                          <div>Vol: {formatVolume(price.volume24h)}</div>
                        )}
                        {showMarketCap && (
                          <div>MCap: {formatMarketCap(price.marketCap)}</div>
                        )}
                      </div>
                    )}

                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                      <span>H: {formatPrice(price.high24h)}</span>
                      <span>L: {formatPrice(price.low24h)}</span>
                    </div>
                  </div>

                  {isSelected && (
                    <motion.div
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      exit={{ opacity: 0, height: 0 }}
                      className="mt-3 pt-3 border-t space-y-2"
                    >
                      <div className="text-xs text-muted-foreground">
                        <div>Source: {price.source}</div>
                        <div>Updated: {new Date(price.timestamp).toLocaleTimeString()}</div>
                      </div>
                      
                      {alerts.length > 0 && (
                        <div className="space-y-1">
                          <div className="text-xs font-medium">Active Alerts:</div>
                          {alerts.slice(0, 2).map((alert) => (
                            <div key={alert.id} className="text-xs text-muted-foreground">
                              {alert.type} {formatPrice(alert.value)}
                            </div>
                          ))}
                        </div>
                      )}
                    </motion.div>
                  )}
                </motion.div>
              )
            })}
          </AnimatePresence>
        </div>

        {state.error && (
          <div className="mt-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <div className="flex items-center space-x-2 text-red-600 dark:text-red-400">
              <span className="text-sm">{state.error}</span>
              <Button variant="ghost" size="sm" onClick={restart}>
                Retry
              </Button>
            </div>
          </div>
        )}

        <div className="mt-4 flex items-center justify-between text-xs text-muted-foreground">
          <div className="flex items-center space-x-4">
            <span>Sources: {state.stats.activeSources}</span>
            <span>Alerts: {state.stats.activeAlerts}</span>
            <span>Subscribers: {state.stats.subscribers}</span>
          </div>
          {state.lastUpdate && (
            <span>Last update: {new Date(state.lastUpdate).toLocaleTimeString()}</span>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
