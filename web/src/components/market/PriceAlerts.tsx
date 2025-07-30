'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { 
  Bell,
  BellOff,
  Plus,
  Trash2,
  TrendingUp,
  TrendingDown,
  Activity,
  CheckCircle,
  AlertTriangle,
  Clock
} from 'lucide-react'
import { usePriceFeed } from '@/hooks/usePriceFeed'
import { type PriceAlert } from '@/lib/price-feed-manager'
import { toast } from 'sonner'

interface PriceAlertsProps {
  symbols?: string[]
}

export function PriceAlerts({ 
  symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL', 'MATIC', 'AVAX'] 
}: PriceAlertsProps) {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [newAlert, setNewAlert] = useState({
    symbol: '',
    type: 'above' as 'above' | 'below' | 'change_percent',
    value: ''
  })

  const {
    state,
    getAlerts,
    addAlert,
    removeAlert,
    getPrice,
    formatPrice
  } = usePriceFeed({
    symbols,
    autoStart: true,
    enableAlerts: true,
    onAlert: (alert, price) => {
      toast.success(`Alert Triggered: ${alert.symbol}`, {
        description: `${alert.symbol} is ${alert.type} ${alert.value}. Current price: ${formatPrice(price.price)}`
      })
    }
  })

  const alerts = getAlerts()
  const activeAlerts = alerts.filter(alert => alert.isActive && !alert.triggered)
  const triggeredAlerts = alerts.filter(alert => alert.triggered)

  const handleCreateAlert = () => {
    if (!newAlert.symbol || !newAlert.value) {
      toast.error('Please fill in all fields')
      return
    }

    const value = parseFloat(newAlert.value)
    if (isNaN(value) || value <= 0) {
      toast.error('Please enter a valid price value')
      return
    }

    addAlert({
      symbol: newAlert.symbol.toUpperCase(),
      type: newAlert.type,
      value,
      isActive: true,
      triggered: false
    })

    setNewAlert({ symbol: '', type: 'above', value: '' })
    setIsDialogOpen(false)
  }

  const handleRemoveAlert = (alertId: string) => {
    removeAlert(alertId)
  }

  const getAlertIcon = (type: string) => {
    switch (type) {
      case 'above':
        return <TrendingUp className="w-4 h-4" />
      case 'below':
        return <TrendingDown className="w-4 h-4" />
      case 'change_percent':
        return <Activity className="w-4 h-4" />
      default:
        return <Bell className="w-4 h-4" />
    }
  }

  const getAlertColor = (alert: PriceAlert) => {
    if (alert.triggered) return 'text-green-600 dark:text-green-400'
    if (!alert.isActive) return 'text-gray-600 dark:text-gray-400'
    return 'text-blue-600 dark:text-blue-400'
  }

  const getAlertDescription = (alert: PriceAlert) => {
    const currentPrice = getPrice(alert.symbol)
    const currentPriceText = currentPrice ? formatPrice(currentPrice.price) : 'N/A'

    switch (alert.type) {
      case 'above':
        return `Alert when ${alert.symbol} goes above ${formatPrice(alert.value)} (Current: ${currentPriceText})`
      case 'below':
        return `Alert when ${alert.symbol} goes below ${formatPrice(alert.value)} (Current: ${currentPriceText})`
      case 'change_percent':
        return `Alert when ${alert.symbol} changes by ${alert.value}% in 24h`
      default:
        return `Alert for ${alert.symbol}`
    }
  }

  const formatTimeAgo = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`
    return 'Just now'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Price Alerts</h2>
          <p className="text-muted-foreground">
            Set up alerts to be notified when prices reach your target levels
          </p>
        </div>
        
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Add Alert
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Price Alert</DialogTitle>
              <DialogDescription>
                Set up a new price alert for your selected cryptocurrency
              </DialogDescription>
            </DialogHeader>
            
            <div className="space-y-4">
              <div>
                <Label htmlFor="symbol">Cryptocurrency</Label>
                <Select value={newAlert.symbol} onValueChange={(value) => setNewAlert(prev => ({ ...prev, symbol: value }))}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a cryptocurrency" />
                  </SelectTrigger>
                  <SelectContent>
                    {symbols.map((symbol) => (
                      <SelectItem key={symbol} value={symbol}>
                        {symbol}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="type">Alert Type</Label>
                <Select value={newAlert.type} onValueChange={(value: any) => setNewAlert(prev => ({ ...prev, type: value }))}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="above">Price Above</SelectItem>
                    <SelectItem value="below">Price Below</SelectItem>
                    <SelectItem value="change_percent">24h Change %</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="value">
                  {newAlert.type === 'change_percent' ? 'Change Percentage' : 'Price Value'}
                </Label>
                <Input
                  id="value"
                  type="number"
                  step="0.01"
                  placeholder={newAlert.type === 'change_percent' ? 'e.g., 5' : 'e.g., 50000'}
                  value={newAlert.value}
                  onChange={(e) => setNewAlert(prev => ({ ...prev, value: e.target.value }))}
                />
              </div>

              {newAlert.symbol && (
                <Alert>
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    {getAlertDescription({
                      ...newAlert,
                      value: parseFloat(newAlert.value) || 0,
                      id: '',
                      isActive: true,
                      triggered: false,
                      createdAt: Date.now()
                    })}
                  </AlertDescription>
                </Alert>
              )}
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreateAlert}>
                Create Alert
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Alerts</p>
                <p className="text-2xl font-bold">{activeAlerts.length}</p>
              </div>
              <Bell className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Triggered Today</p>
                <p className="text-2xl font-bold">{triggeredAlerts.length}</p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Alerts</p>
                <p className="text-2xl font-bold">{alerts.length}</p>
              </div>
              <Activity className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Active Alerts */}
      <Card>
        <CardHeader>
          <CardTitle>Active Alerts</CardTitle>
          <CardDescription>
            Alerts that are currently monitoring price movements
          </CardDescription>
        </CardHeader>
        <CardContent>
          {activeAlerts.length > 0 ? (
            <div className="space-y-3">
              <AnimatePresence>
                {activeAlerts.map((alert) => (
                  <motion.div
                    key={alert.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center space-x-3">
                      <div className={`w-10 h-10 rounded-lg flex items-center justify-center bg-blue-100 dark:bg-blue-900 ${getAlertColor(alert)}`}>
                        {getAlertIcon(alert.type)}
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <span className="font-medium">{alert.symbol}</span>
                          <Badge variant="outline" className="text-xs">
                            {alert.type.replace('_', ' ')}
                          </Badge>
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {getAlertDescription(alert)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Created {formatTimeAgo(alert.createdAt)}
                        </p>
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleRemoveAlert(alert.id)}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <div className="text-center py-8">
              <BellOff className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">No Active Alerts</h3>
              <p className="text-muted-foreground mb-4">
                Create your first price alert to get notified when prices reach your target levels
              </p>
              <Button onClick={() => setIsDialogOpen(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Add Your First Alert
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Triggered Alerts */}
      {triggeredAlerts.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Recently Triggered</CardTitle>
            <CardDescription>
              Alerts that have been triggered recently
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <AnimatePresence>
                {triggeredAlerts.slice(0, 5).map((alert) => (
                  <motion.div
                    key={alert.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="flex items-center justify-between p-4 border rounded-lg bg-green-50 dark:bg-green-900/20"
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-10 h-10 rounded-lg flex items-center justify-center bg-green-100 dark:bg-green-900 text-green-600 dark:text-green-400">
                        <CheckCircle className="w-5 h-5" />
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <span className="font-medium">{alert.symbol}</span>
                          <Badge variant="default" className="text-xs bg-green-600">
                            Triggered
                          </Badge>
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {getAlertDescription(alert)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Triggered {alert.triggeredAt ? formatTimeAgo(alert.triggeredAt) : 'recently'}
                        </p>
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleRemoveAlert(alert.id)}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
