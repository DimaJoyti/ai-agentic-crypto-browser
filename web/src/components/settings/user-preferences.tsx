'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { 
  User, 
  Bell, 
  Shield, 
  Palette, 
  Globe, 
  Zap, 
  TrendingUp,
  DollarSign,
  Clock,
  Volume2,
  Smartphone,
  Mail,
  MessageSquare,
  Eye,
  Lock,
  Key,
  Fingerprint,
  AlertTriangle,
  Save,
  RotateCcw
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface UserPreferences {
  // Profile
  displayName: string
  email: string
  timezone: string
  language: string
  currency: string
  
  // Notifications
  emailNotifications: boolean
  pushNotifications: boolean
  tradingAlerts: boolean
  priceAlerts: boolean
  newsAlerts: boolean
  securityAlerts: boolean
  
  // Trading
  defaultOrderType: string
  confirmOrders: boolean
  showAdvancedTools: boolean
  autoRefreshInterval: number
  riskWarnings: boolean
  
  // Display
  theme: string
  compactMode: boolean
  showAnimations: boolean
  highContrast: boolean
  fontSize: string
  
  // Privacy & Security
  twoFactorAuth: boolean
  sessionTimeout: number
  dataSharing: boolean
  analyticsTracking: boolean
  
  // Advanced
  apiAccess: boolean
  webhookUrl: string
  customCss: string
}

const defaultPreferences: UserPreferences = {
  displayName: 'John Doe',
  email: 'john@example.com',
  timezone: 'UTC',
  language: 'en',
  currency: 'USD',
  emailNotifications: true,
  pushNotifications: true,
  tradingAlerts: true,
  priceAlerts: true,
  newsAlerts: false,
  securityAlerts: true,
  defaultOrderType: 'limit',
  confirmOrders: true,
  showAdvancedTools: false,
  autoRefreshInterval: 5,
  riskWarnings: true,
  theme: 'system',
  compactMode: false,
  showAnimations: true,
  highContrast: false,
  fontSize: 'medium',
  twoFactorAuth: false,
  sessionTimeout: 30,
  dataSharing: false,
  analyticsTracking: true,
  apiAccess: false,
  webhookUrl: '',
  customCss: ''
}

export function UserPreferences() {
  const [preferences, setPreferences] = useState<UserPreferences>(defaultPreferences)
  const [hasChanges, setHasChanges] = useState(false)
  const [saving, setSaving] = useState(false)

  const updatePreference = <K extends keyof UserPreferences>(
    key: K,
    value: UserPreferences[K]
  ) => {
    setPreferences(prev => ({ ...prev, [key]: value }))
    setHasChanges(true)
  }

  const savePreferences = async () => {
    setSaving(true)
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    setSaving(false)
    setHasChanges(false)
  }

  const resetToDefaults = () => {
    setPreferences(defaultPreferences)
    setHasChanges(true)
  }

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">User Preferences</h1>
          <p className="text-muted-foreground">Customize your trading experience</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="outline" onClick={resetToDefaults}>
            <RotateCcw className="h-4 w-4 mr-2" />
            Reset
          </Button>
          <Button 
            onClick={savePreferences} 
            disabled={!hasChanges || saving}
            className="min-w-[100px]"
          >
            {saving ? (
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
              >
                <Save className="h-4 w-4 mr-2" />
              </motion.div>
            ) : (
              <Save className="h-4 w-4 mr-2" />
            )}
            {saving ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </div>

      {hasChanges && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-yellow-50 border border-yellow-200 rounded-lg p-4"
        >
          <div className="flex items-center space-x-2">
            <AlertTriangle className="h-4 w-4 text-yellow-600" />
            <span className="text-sm text-yellow-800">You have unsaved changes</span>
          </div>
        </motion.div>
      )}

      <Tabs defaultValue="profile" className="w-full">
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="profile">Profile</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
          <TabsTrigger value="trading">Trading</TabsTrigger>
          <TabsTrigger value="display">Display</TabsTrigger>
          <TabsTrigger value="security">Security</TabsTrigger>
          <TabsTrigger value="advanced">Advanced</TabsTrigger>
        </TabsList>

        {/* Profile Tab */}
        <TabsContent value="profile" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <User className="h-5 w-5" />
                <span>Profile Information</span>
              </CardTitle>
              <CardDescription>
                Manage your personal information and regional settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="displayName">Display Name</Label>
                  <Input
                    id="displayName"
                    value={preferences.displayName}
                    onChange={(e) => updatePreference('displayName', e.target.value)}
                  />
                </div>
                <div>
                  <Label htmlFor="email">Email Address</Label>
                  <Input
                    id="email"
                    type="email"
                    value={preferences.email}
                    onChange={(e) => updatePreference('email', e.target.value)}
                  />
                </div>
              </div>
              
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="timezone">Timezone</Label>
                  <Select value={preferences.timezone} onValueChange={(value) => updatePreference('timezone', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="UTC">UTC</SelectItem>
                      <SelectItem value="EST">Eastern Time</SelectItem>
                      <SelectItem value="PST">Pacific Time</SelectItem>
                      <SelectItem value="GMT">Greenwich Mean Time</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="language">Language</Label>
                  <Select value={preferences.language} onValueChange={(value) => updatePreference('language', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="en">English</SelectItem>
                      <SelectItem value="es">Spanish</SelectItem>
                      <SelectItem value="fr">French</SelectItem>
                      <SelectItem value="de">German</SelectItem>
                      <SelectItem value="ja">Japanese</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="currency">Default Currency</Label>
                  <Select value={preferences.currency} onValueChange={(value) => updatePreference('currency', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="USD">USD ($)</SelectItem>
                      <SelectItem value="EUR">EUR (€)</SelectItem>
                      <SelectItem value="GBP">GBP (£)</SelectItem>
                      <SelectItem value="JPY">JPY (¥)</SelectItem>
                      <SelectItem value="BTC">BTC (₿)</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notifications Tab */}
        <TabsContent value="notifications" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Bell className="h-5 w-5" />
                <span>Notification Preferences</span>
              </CardTitle>
              <CardDescription>
                Choose how and when you want to be notified
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <h4 className="font-medium flex items-center space-x-2">
                  <Mail className="h-4 w-4" />
                  <span>Delivery Methods</span>
                </h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="emailNotifications">Email Notifications</Label>
                      <p className="text-sm text-muted-foreground">Receive notifications via email</p>
                    </div>
                    <Switch
                      id="emailNotifications"
                      checked={preferences.emailNotifications}
                      onCheckedChange={(checked) => updatePreference('emailNotifications', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="pushNotifications">Push Notifications</Label>
                      <p className="text-sm text-muted-foreground">Browser and mobile notifications</p>
                    </div>
                    <Switch
                      id="pushNotifications"
                      checked={preferences.pushNotifications}
                      onCheckedChange={(checked) => updatePreference('pushNotifications', checked)}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <h4 className="font-medium flex items-center space-x-2">
                  <TrendingUp className="h-4 w-4" />
                  <span>Trading Alerts</span>
                </h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="tradingAlerts">Order Execution</Label>
                      <p className="text-sm text-muted-foreground">When orders are filled or cancelled</p>
                    </div>
                    <Switch
                      id="tradingAlerts"
                      checked={preferences.tradingAlerts}
                      onCheckedChange={(checked) => updatePreference('tradingAlerts', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="priceAlerts">Price Alerts</Label>
                      <p className="text-sm text-muted-foreground">When price targets are reached</p>
                    </div>
                    <Switch
                      id="priceAlerts"
                      checked={preferences.priceAlerts}
                      onCheckedChange={(checked) => updatePreference('priceAlerts', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="newsAlerts">Market News</Label>
                      <p className="text-sm text-muted-foreground">Important market updates</p>
                    </div>
                    <Switch
                      id="newsAlerts"
                      checked={preferences.newsAlerts}
                      onCheckedChange={(checked) => updatePreference('newsAlerts', checked)}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <h4 className="font-medium flex items-center space-x-2">
                  <Shield className="h-4 w-4" />
                  <span>Security Alerts</span>
                </h4>
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="securityAlerts">Security Events</Label>
                    <p className="text-sm text-muted-foreground">Login attempts and security changes</p>
                  </div>
                  <Switch
                    id="securityAlerts"
                    checked={preferences.securityAlerts}
                    onCheckedChange={(checked) => updatePreference('securityAlerts', checked)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Trading Tab */}
        <TabsContent value="trading" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <TrendingUp className="h-5 w-5" />
                <span>Trading Preferences</span>
              </CardTitle>
              <CardDescription>
                Configure your trading interface and behavior
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="defaultOrderType">Default Order Type</Label>
                  <Select value={preferences.defaultOrderType} onValueChange={(value) => updatePreference('defaultOrderType', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="market">Market Order</SelectItem>
                      <SelectItem value="limit">Limit Order</SelectItem>
                      <SelectItem value="stop">Stop Order</SelectItem>
                      <SelectItem value="stop-limit">Stop-Limit Order</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="autoRefreshInterval">Auto Refresh (seconds)</Label>
                  <Select 
                    value={preferences.autoRefreshInterval.toString()} 
                    onValueChange={(value) => updatePreference('autoRefreshInterval', parseInt(value))}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1">1 second</SelectItem>
                      <SelectItem value="5">5 seconds</SelectItem>
                      <SelectItem value="10">10 seconds</SelectItem>
                      <SelectItem value="30">30 seconds</SelectItem>
                      <SelectItem value="60">1 minute</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="confirmOrders">Order Confirmation</Label>
                    <p className="text-sm text-muted-foreground">Require confirmation before placing orders</p>
                  </div>
                  <Switch
                    id="confirmOrders"
                    checked={preferences.confirmOrders}
                    onCheckedChange={(checked) => updatePreference('confirmOrders', checked)}
                  />
                </div>
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="showAdvancedTools">Advanced Tools</Label>
                    <p className="text-sm text-muted-foreground">Show advanced trading features</p>
                  </div>
                  <Switch
                    id="showAdvancedTools"
                    checked={preferences.showAdvancedTools}
                    onCheckedChange={(checked) => updatePreference('showAdvancedTools', checked)}
                  />
                </div>
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="riskWarnings">Risk Warnings</Label>
                    <p className="text-sm text-muted-foreground">Show risk warnings for large orders</p>
                  </div>
                  <Switch
                    id="riskWarnings"
                    checked={preferences.riskWarnings}
                    onCheckedChange={(checked) => updatePreference('riskWarnings', checked)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Display Tab */}
        <TabsContent value="display" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Palette className="h-5 w-5" />
                <span>Display Settings</span>
              </CardTitle>
              <CardDescription>
                Customize the appearance and accessibility
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="theme">Theme</Label>
                  <Select value={preferences.theme} onValueChange={(value) => updatePreference('theme', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="light">Light</SelectItem>
                      <SelectItem value="dark">Dark</SelectItem>
                      <SelectItem value="system">System</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="fontSize">Font Size</Label>
                  <Select value={preferences.fontSize} onValueChange={(value) => updatePreference('fontSize', value)}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="small">Small</SelectItem>
                      <SelectItem value="medium">Medium</SelectItem>
                      <SelectItem value="large">Large</SelectItem>
                      <SelectItem value="extra-large">Extra Large</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="compactMode">Compact Mode</Label>
                    <p className="text-sm text-muted-foreground">Reduce spacing and padding</p>
                  </div>
                  <Switch
                    id="compactMode"
                    checked={preferences.compactMode}
                    onCheckedChange={(checked) => updatePreference('compactMode', checked)}
                  />
                </div>
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="showAnimations">Animations</Label>
                    <p className="text-sm text-muted-foreground">Enable smooth transitions and effects</p>
                  </div>
                  <Switch
                    id="showAnimations"
                    checked={preferences.showAnimations}
                    onCheckedChange={(checked) => updatePreference('showAnimations', checked)}
                  />
                </div>
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="highContrast">High Contrast</Label>
                    <p className="text-sm text-muted-foreground">Increase color contrast for better visibility</p>
                  </div>
                  <Switch
                    id="highContrast"
                    checked={preferences.highContrast}
                    onCheckedChange={(checked) => updatePreference('highContrast', checked)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Tab */}
        <TabsContent value="security" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Shield className="h-5 w-5" />
                <span>Security & Privacy</span>
              </CardTitle>
              <CardDescription>
                Manage your account security and privacy settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <h4 className="font-medium flex items-center space-x-2">
                  <Lock className="h-4 w-4" />
                  <span>Authentication</span>
                </h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="twoFactorAuth">Two-Factor Authentication</Label>
                      <p className="text-sm text-muted-foreground">Add an extra layer of security</p>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="twoFactorAuth"
                        checked={preferences.twoFactorAuth}
                        onCheckedChange={(checked) => updatePreference('twoFactorAuth', checked)}
                      />
                      {preferences.twoFactorAuth && (
                        <Badge variant="secondary">Enabled</Badge>
                      )}
                    </div>
                  </div>
                  <div>
                    <Label htmlFor="sessionTimeout">Session Timeout (minutes)</Label>
                    <Select 
                      value={preferences.sessionTimeout.toString()} 
                      onValueChange={(value) => updatePreference('sessionTimeout', parseInt(value))}
                    >
                      <SelectTrigger className="mt-1">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="15">15 minutes</SelectItem>
                        <SelectItem value="30">30 minutes</SelectItem>
                        <SelectItem value="60">1 hour</SelectItem>
                        <SelectItem value="240">4 hours</SelectItem>
                        <SelectItem value="480">8 hours</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <h4 className="font-medium flex items-center space-x-2">
                  <Eye className="h-4 w-4" />
                  <span>Privacy</span>
                </h4>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="dataSharing">Data Sharing</Label>
                      <p className="text-sm text-muted-foreground">Share anonymized data for platform improvement</p>
                    </div>
                    <Switch
                      id="dataSharing"
                      checked={preferences.dataSharing}
                      onCheckedChange={(checked) => updatePreference('dataSharing', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="analyticsTracking">Analytics Tracking</Label>
                      <p className="text-sm text-muted-foreground">Allow usage analytics for better experience</p>
                    </div>
                    <Switch
                      id="analyticsTracking"
                      checked={preferences.analyticsTracking}
                      onCheckedChange={(checked) => updatePreference('analyticsTracking', checked)}
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Advanced Tab */}
        <TabsContent value="advanced" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Zap className="h-5 w-5" />
                <span>Advanced Settings</span>
              </CardTitle>
              <CardDescription>
                Advanced configuration for power users
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="apiAccess">API Access</Label>
                    <p className="text-sm text-muted-foreground">Enable programmatic access to your account</p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Switch
                      id="apiAccess"
                      checked={preferences.apiAccess}
                      onCheckedChange={(checked) => updatePreference('apiAccess', checked)}
                    />
                    {preferences.apiAccess && (
                      <Badge variant="outline">Beta</Badge>
                    )}
                  </div>
                </div>

                {preferences.apiAccess && (
                  <div>
                    <Label htmlFor="webhookUrl">Webhook URL</Label>
                    <Input
                      id="webhookUrl"
                      placeholder="https://your-webhook-url.com"
                      value={preferences.webhookUrl}
                      onChange={(e) => updatePreference('webhookUrl', e.target.value)}
                      className="mt-1"
                    />
                    <p className="text-sm text-muted-foreground mt-1">
                      Receive real-time notifications via webhook
                    </p>
                  </div>
                )}

                <div>
                  <Label htmlFor="customCss">Custom CSS</Label>
                  <textarea
                    id="customCss"
                    placeholder="/* Add your custom styles here */"
                    value={preferences.customCss}
                    onChange={(e) => updatePreference('customCss', e.target.value)}
                    className="w-full h-32 mt-1 p-3 border border-border rounded-md bg-background font-mono text-sm"
                  />
                  <p className="text-sm text-muted-foreground mt-1">
                    Customize the appearance with your own CSS
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
