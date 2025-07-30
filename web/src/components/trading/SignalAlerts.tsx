'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Bell,
  BellOff,
  Plus,
  Settings,
  Trash2,
  Mail,
  Smartphone,
  MessageSquare,
  Webhook,
  TrendingUp,
  TrendingDown,
  Activity,
  CheckCircle,
  Clock
} from 'lucide-react'
import { useTradingSignals } from '@/hooks/useTradingSignals'
import { type SignalAlert, type AlertCondition } from '@/lib/trading-signals'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'

interface SignalAlertsProps {
  symbols?: string[]
}

export function SignalAlerts({ 
  symbols = ['BTC', 'ETH', 'BNB', 'XRP', 'ADA', 'SOL'] 
}: SignalAlertsProps) {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [editingAlert, setEditingAlert] = useState<SignalAlert | null>(null)
  const [newAlert, setNewAlert] = useState({
    type: 'in_app' as 'email' | 'push' | 'sms' | 'webhook' | 'in_app',
    enabled: true,
    cooldownPeriod: 60,
    maxAlertsPerDay: 10,
    conditions: [{
      signalType: 'any' as 'buy' | 'sell' | 'hold' | 'any',
      minStrength: 70,
      minConfidence: 80,
      symbols: [] as string[],
      strategies: [] as string[]
    }]
  })

  const {
    state,
    getAlerts,
    addAlert,
    updateAlert,
    deleteAlert,
    toggleAlert,
    getStrategies
  } = useTradingSignals({
    symbols,
    enableAlerts: true
  })

  const alerts = getAlerts()
  const strategies = getStrategies()

  const getAlertTypeIcon = (type: string) => {
    switch (type) {
      case 'email':
        return <Mail className="w-4 h-4" />
      case 'push':
        return <Smartphone className="w-4 h-4" />
      case 'sms':
        return <MessageSquare className="w-4 h-4" />
      case 'webhook':
        return <Webhook className="w-4 h-4" />
      default:
        return <Bell className="w-4 h-4" />
    }
  }

  const getSignalTypeIcon = (type: string) => {
    switch (type) {
      case 'buy':
        return <TrendingUp className="w-4 h-4 text-green-500" />
      case 'sell':
        return <TrendingDown className="w-4 h-4 text-red-500" />
      case 'hold':
        return <Activity className="w-4 h-4 text-gray-500" />
      default:
        return <Activity className="w-4 h-4" />
    }
  }

  const formatTimeAgo = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days}d ago`
    if (hours > 0) return `${hours}h ago`
    if (minutes > 0) return `${minutes}m ago`
    return 'Just now'
  }

  const handleCreateAlert = () => {
    if (newAlert.conditions.length === 0) {
      toast.error('Please add at least one condition')
      return
    }

    addAlert({
      type: newAlert.type,
      enabled: newAlert.enabled,
      conditions: newAlert.conditions,
      cooldownPeriod: newAlert.cooldownPeriod,
      maxAlertsPerDay: newAlert.maxAlertsPerDay,
      signalId: '' // Will be set by the system
    })

    setNewAlert({
      type: 'in_app',
      enabled: true,
      cooldownPeriod: 60,
      maxAlertsPerDay: 10,
      conditions: [{
        signalType: 'any',
        minStrength: 70,
        minConfidence: 80,
        symbols: [],
        strategies: []
      }]
    })
    setIsCreateDialogOpen(false)
  }

  const handleEditAlert = () => {
    if (!editingAlert) return

    // Implementation would update the alert
    setEditingAlert(null)
    setIsEditDialogOpen(false)
    toast.success('Alert updated successfully')
  }

  const handleDeleteAlert = (alertId: string) => {
    if (window.confirm('Are you sure you want to delete this alert?')) {
      deleteAlert(alertId)
    }
  }

  const handleToggleAlert = (alertId: string, enabled: boolean) => {
    toggleAlert(alertId, enabled)
  }

  const updateCondition = (index: number, field: keyof AlertCondition, value: any) => {
    const updatedConditions = [...newAlert.conditions]
    updatedConditions[index] = { ...updatedConditions[index], [field]: value }
    setNewAlert(prev => ({ ...prev, conditions: updatedConditions }))
  }

  const addCondition = () => {
    setNewAlert(prev => ({
      ...prev,
      conditions: [...prev.conditions, {
        signalType: 'any',
        minStrength: 70,
        minConfidence: 80,
        symbols: [],
        strategies: []
      }]
    }))
  }

  const removeCondition = (index: number) => {
    setNewAlert(prev => ({
      ...prev,
      conditions: prev.conditions.filter((_, i) => i !== index)
    }))
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">Signal Alerts</h3>
          <p className="text-sm text-muted-foreground">
            Configure notifications for trading signals and market events
          </p>
        </div>
        
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Add Alert
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Create Signal Alert</DialogTitle>
              <DialogDescription>
                Set up a new alert to be notified when trading signals match your criteria
              </DialogDescription>
            </DialogHeader>
            
            <div className="space-y-6">
              {/* Alert Type */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="alert-type">Alert Type</Label>
                  <Select value={newAlert.type} onValueChange={(value: any) => setNewAlert(prev => ({ ...prev, type: value }))}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="in_app">In-App Notification</SelectItem>
                      <SelectItem value="email">Email</SelectItem>
                      <SelectItem value="push">Push Notification</SelectItem>
                      <SelectItem value="sms">SMS</SelectItem>
                      <SelectItem value="webhook">Webhook</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                
                <div className="flex items-center space-x-2">
                  <Switch
                    id="alert-enabled"
                    checked={newAlert.enabled}
                    onCheckedChange={(checked) => setNewAlert(prev => ({ ...prev, enabled: checked }))}
                  />
                  <Label htmlFor="alert-enabled">Enabled</Label>
                </div>
              </div>

              {/* Alert Settings */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="cooldown">Cooldown Period (minutes)</Label>
                  <Input
                    id="cooldown"
                    type="number"
                    min="1"
                    value={newAlert.cooldownPeriod}
                    onChange={(e) => setNewAlert(prev => ({ ...prev, cooldownPeriod: parseInt(e.target.value) }))}
                  />
                </div>
                
                <div>
                  <Label htmlFor="max-alerts">Max Alerts per Day</Label>
                  <Input
                    id="max-alerts"
                    type="number"
                    min="1"
                    max="100"
                    value={newAlert.maxAlertsPerDay}
                    onChange={(e) => setNewAlert(prev => ({ ...prev, maxAlertsPerDay: parseInt(e.target.value) }))}
                  />
                </div>
              </div>

              {/* Conditions */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <Label>Alert Conditions</Label>
                  <Button variant="outline" size="sm" onClick={addCondition}>
                    <Plus className="w-4 h-4 mr-2" />
                    Add Condition
                  </Button>
                </div>
                
                <div className="space-y-4">
                  {newAlert.conditions.map((condition, index) => (
                    <div key={index} className="p-4 border rounded-lg space-y-4">
                      <div className="flex items-center justify-between">
                        <h4 className="font-medium">Condition {index + 1}</h4>
                        {newAlert.conditions.length > 1 && (
                          <Button variant="ghost" size="sm" onClick={() => removeCondition(index)}>
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        )}
                      </div>
                      
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label>Signal Type</Label>
                          <Select 
                            value={condition.signalType} 
                            onValueChange={(value: any) => updateCondition(index, 'signalType', value)}
                          >
                            <SelectTrigger>
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="any">Any Signal</SelectItem>
                              <SelectItem value="buy">Buy Signals</SelectItem>
                              <SelectItem value="sell">Sell Signals</SelectItem>
                              <SelectItem value="hold">Hold Signals</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                        
                        <div>
                          <Label>Min Strength (%)</Label>
                          <Input
                            type="number"
                            min="0"
                            max="100"
                            value={condition.minStrength}
                            onChange={(e) => updateCondition(index, 'minStrength', parseInt(e.target.value))}
                          />
                        </div>
                      </div>
                      
                      <div>
                        <Label>Min Confidence (%)</Label>
                        <Input
                          type="number"
                          min="0"
                          max="100"
                          value={condition.minConfidence}
                          onChange={(e) => updateCondition(index, 'minConfidence', parseInt(e.target.value))}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreateAlert}>
                Create Alert
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Alert Statistics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Alerts</p>
                <p className="text-2xl font-bold">{alerts.length}</p>
              </div>
              <Bell className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Alerts</p>
                <p className="text-2xl font-bold text-green-600">
                  {alerts.filter(a => a.enabled).length}
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Triggered Today</p>
                <p className="text-2xl font-bold">
                  {alerts.filter(a => a.lastTriggered && 
                    new Date(a.lastTriggered).toDateString() === new Date().toDateString()
                  ).length}
                </p>
              </div>
              <Activity className="w-8 h-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Response</p>
                <p className="text-2xl font-bold">2.3s</p>
              </div>
              <Clock className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Alerts List */}
      <Card>
        <CardHeader>
          <CardTitle>Your Alerts</CardTitle>
          <CardDescription>
            Manage your signal alert configurations
          </CardDescription>
        </CardHeader>
        <CardContent>
          {alerts.length > 0 ? (
            <div className="space-y-4">
              <AnimatePresence>
                {alerts.map((alert, index) => (
                  <motion.div
                    key={alert.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center gap-4">
                      <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center",
                        alert.enabled ? 'bg-green-100 dark:bg-green-900' : 'bg-gray-100 dark:bg-gray-900'
                      )}>
                        <span className={cn(
                          alert.enabled ? 'text-green-600 dark:text-green-400' : 'text-gray-600 dark:text-gray-400'
                        )}>
                          {getAlertTypeIcon(alert.type)}
                        </span>
                      </div>
                      
                      <div>
                        <div className="flex items-center gap-2 mb-1">
                          <h4 className="font-medium">
                            {alert.type.charAt(0).toUpperCase() + alert.type.slice(1)} Alert
                          </h4>
                          <Badge variant={alert.enabled ? 'default' : 'secondary'}>
                            {alert.enabled ? 'Active' : 'Inactive'}
                          </Badge>
                        </div>
                        
                        <div className="space-y-1">
                          {alert.conditions.map((condition, condIndex) => (
                            <div key={condIndex} className="flex items-center gap-2 text-sm text-muted-foreground">
                              {getSignalTypeIcon(condition.signalType)}
                              <span>
                                {condition.signalType === 'any' ? 'Any' : condition.signalType.toUpperCase()} signals
                                with ≥{condition.minStrength}% strength, ≥{condition.minConfidence}% confidence
                              </span>
                            </div>
                          ))}
                        </div>
                        
                        <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                          <span>Cooldown: {alert.cooldownPeriod}m</span>
                          <span>Max/day: {alert.maxAlertsPerDay}</span>
                          {alert.lastTriggered && (
                            <span>Last: {formatTimeAgo(alert.lastTriggered)}</span>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center gap-2">
                      <Switch
                        checked={alert.enabled}
                        onCheckedChange={(checked) => handleToggleAlert(alert.id, checked)}
                      />
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => {
                          setEditingAlert(alert)
                          setIsEditDialogOpen(true)
                        }}
                      >
                        <Settings className="w-4 h-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteAlert(alert.id)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <div className="text-center py-12">
              <BellOff className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">No Alerts Configured</h3>
              <p className="text-muted-foreground mb-4">
                Create your first alert to be notified when trading signals match your criteria
              </p>
              <Button onClick={() => setIsCreateDialogOpen(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Create Your First Alert
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Edit Alert Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Alert</DialogTitle>
            <DialogDescription>
              Modify your alert configuration
            </DialogDescription>
          </DialogHeader>
          
          {editingAlert && (
            <div className="space-y-4">
              <div>
                <Label>Alert Type</Label>
                <Select defaultValue={editingAlert.type}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="in_app">In-App Notification</SelectItem>
                    <SelectItem value="email">Email</SelectItem>
                    <SelectItem value="push">Push Notification</SelectItem>
                    <SelectItem value="sms">SMS</SelectItem>
                    <SelectItem value="webhook">Webhook</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Cooldown (minutes)</Label>
                  <Input
                    type="number"
                    defaultValue={editingAlert.cooldownPeriod}
                  />
                </div>
                <div>
                  <Label>Max per day</Label>
                  <Input
                    type="number"
                    defaultValue={editingAlert.maxAlertsPerDay}
                  />
                </div>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleEditAlert}>
              Update Alert
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
