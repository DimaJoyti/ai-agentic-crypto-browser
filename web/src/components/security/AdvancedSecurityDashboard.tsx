'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import { Progress } from '@/components/ui/progress'
import {
  Shield,
  Key,
  Smartphone,
  Lock,
  Unlock,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Eye,
  EyeOff,
  Clock,
  MapPin,
  Monitor,
  Wifi,
  Globe,
  UserCheck,
  Settings,
  Download,
  Upload,
  Activity
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'
import { MultiFactorAuth } from './MultiFactorAuth'
import { SecurityMonitoring } from './SecurityMonitoring'
import { FraudDetection } from './FraudDetection'

interface SecurityEvent {
  id: string
  type: 'login' | 'withdrawal' | 'api_access' | 'password_change' | 'suspicious'
  timestamp: number
  location: string
  device: string
  ipAddress: string
  status: 'success' | 'failed' | 'blocked'
  riskLevel: 'low' | 'medium' | 'high'
}

interface WhitelistAddress {
  id: string
  address: string
  label: string
  network: string
  addedAt: number
  isActive: boolean
  lastUsed?: number
}

interface APIKey {
  id: string
  name: string
  key: string
  permissions: string[]
  createdAt: number
  lastUsed?: number
  isActive: boolean
  ipWhitelist: string[]
}

export function AdvancedSecurityDashboard() {
  const [securityEvents, setSecurityEvents] = useState<SecurityEvent[]>([])
  const [whitelistAddresses, setWhitelistAddresses] = useState<WhitelistAddress[]>([])
  const [apiKeys, setApiKeys] = useState<APIKey[]>([])
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false)
  const [withdrawalLockEnabled, setWithdrawalLockEnabled] = useState(true)
  const [sessionTimeout, setSessionTimeout] = useState(30)
  const [securityScore, setSecurityScore] = useState(85)
  const [activeTab, setActiveTab] = useState('overview')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock security data
    const mockEvents: SecurityEvent[] = [
      {
        id: 'event1',
        type: 'login',
        timestamp: Date.now() - 3600000,
        location: 'New York, US',
        device: 'Chrome on Windows',
        ipAddress: '192.168.1.100',
        status: 'success',
        riskLevel: 'low'
      },
      {
        id: 'event2',
        type: 'withdrawal',
        timestamp: Date.now() - 7200000,
        location: 'New York, US',
        device: 'Chrome on Windows',
        ipAddress: '192.168.1.100',
        status: 'success',
        riskLevel: 'medium'
      },
      {
        id: 'event3',
        type: 'suspicious',
        timestamp: Date.now() - 86400000,
        location: 'Unknown',
        device: 'Unknown Browser',
        ipAddress: '45.123.45.67',
        status: 'blocked',
        riskLevel: 'high'
      }
    ]

    const mockWhitelist: WhitelistAddress[] = [
      {
        id: 'addr1',
        address: '0x742d35Cc6634C0532925a3b8D4C2C4e4C4C4C4C4',
        label: 'Main Wallet',
        network: 'Ethereum',
        addedAt: Date.now() - 86400000 * 7,
        isActive: true,
        lastUsed: Date.now() - 3600000
      },
      {
        id: 'addr2',
        address: '0x123d35Cc6634C0532925a3b8D4C2C4e4C4C4C4C4',
        label: 'Cold Storage',
        network: 'Ethereum',
        addedAt: Date.now() - 86400000 * 30,
        isActive: true
      }
    ]

    const mockApiKeys: APIKey[] = [
      {
        id: 'api1',
        name: 'Trading Bot',
        key: 'ak_live_1234567890abcdef',
        permissions: ['read', 'trade'],
        createdAt: Date.now() - 86400000 * 14,
        lastUsed: Date.now() - 3600000,
        isActive: true,
        ipWhitelist: ['192.168.1.100', '10.0.0.1']
      }
    ]

    setSecurityEvents(mockEvents)
    setWhitelistAddresses(mockWhitelist)
    setApiKeys(mockApiKeys)
  }, [isConnected])

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'low': return 'text-green-500'
      case 'medium': return 'text-yellow-500'
      case 'high': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskBadgeVariant = (level: string) => {
    switch (level) {
      case 'low': return 'default'
      case 'medium': return 'secondary'
      case 'high': return 'destructive'
      default: return 'outline'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success': return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'failed': return <XCircle className="w-4 h-4 text-red-500" />
      case 'blocked': return <Shield className="w-4 h-4 text-red-500" />
      default: return <AlertTriangle className="w-4 h-4 text-yellow-500" />
    }
  }

  const handleToggle2FA = () => {
    setTwoFactorEnabled(!twoFactorEnabled)
    // In real implementation, this would trigger 2FA setup/disable flow
  }

  const handleAddWhitelistAddress = () => {
    // In real implementation, this would open a modal to add new address
    console.log('Add whitelist address')
  }

  const handleCreateAPIKey = () => {
    // In real implementation, this would open API key creation modal
    console.log('Create API key')
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Your Wallet</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access advanced security features and monitoring
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Security Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Shield className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Security Score</span>
            </div>
            <div className="text-2xl font-bold mb-2">{securityScore}/100</div>
            <Progress value={securityScore} className="h-2" />
            <div className="text-xs text-muted-foreground mt-1">
              {securityScore >= 80 ? 'Excellent' : securityScore >= 60 ? 'Good' : 'Needs Improvement'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Key className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">2FA Status</span>
            </div>
            <div className="flex items-center gap-2">
              <div className={cn(
                "text-lg font-bold",
                twoFactorEnabled ? "text-green-500" : "text-red-500"
              )}>
                {twoFactorEnabled ? 'Enabled' : 'Disabled'}
              </div>
              {twoFactorEnabled ? (
                <CheckCircle className="w-4 h-4 text-green-500" />
              ) : (
                <XCircle className="w-4 h-4 text-red-500" />
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Lock className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Whitelist Addresses</span>
            </div>
            <div className="text-2xl font-bold">{whitelistAddresses.filter(a => a.isActive).length}</div>
            <div className="text-xs text-muted-foreground">
              Active addresses
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <AlertTriangle className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Security Events</span>
            </div>
            <div className="text-2xl font-bold">
              {securityEvents.filter(e => e.timestamp > Date.now() - 86400000).length}
            </div>
            <div className="text-xs text-muted-foreground">
              Last 24 hours
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Security Dashboard */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="mfa">Multi-Factor Auth</TabsTrigger>
          <TabsTrigger value="monitoring">Monitoring</TabsTrigger>
          <TabsTrigger value="fraud">Fraud Detection</TabsTrigger>
          <TabsTrigger value="whitelist">Whitelist & API</TabsTrigger>
          <TabsTrigger value="activity">Activity Log</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Security Settings */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="w-5 h-5" />
                  Security Settings
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="font-medium">Two-Factor Authentication</Label>
                    <p className="text-sm text-muted-foreground">
                      Add an extra layer of security to your account
                    </p>
                  </div>
                  <Switch checked={twoFactorEnabled} onCheckedChange={handleToggle2FA} />
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label className="font-medium">Withdrawal Lock</Label>
                    <p className="text-sm text-muted-foreground">
                      Require additional confirmation for withdrawals
                    </p>
                  </div>
                  <Switch checked={withdrawalLockEnabled} onCheckedChange={setWithdrawalLockEnabled} />
                </div>

                <div className="space-y-2">
                  <Label className="font-medium">Session Timeout (minutes)</Label>
                  <Input
                    type="number"
                    value={sessionTimeout}
                    onChange={(e) => setSessionTimeout(Number(e.target.value))}
                    min="5"
                    max="120"
                  />
                  <p className="text-xs text-muted-foreground">
                    Automatically log out after inactivity
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* Recent Security Events */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="w-5 h-5" />
                  Recent Security Events
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {securityEvents.slice(0, 5).map((event) => (
                    <div key={event.id} className="flex items-center justify-between p-3 border rounded">
                      <div className="flex items-center gap-3">
                        {getStatusIcon(event.status)}
                        <div>
                          <div className="font-medium capitalize">{event.type.replace('_', ' ')}</div>
                          <div className="text-sm text-muted-foreground">
                            {event.location} • {event.device}
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <Badge variant={getRiskBadgeVariant(event.riskLevel)}>
                          {event.riskLevel}
                        </Badge>
                        <div className="text-xs text-muted-foreground mt-1">
                          {formatTime(event.timestamp)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>



        <TabsContent value="whitelist" className="space-y-6">
          {/* Withdrawal Whitelist */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Lock className="w-5 h-5" />
                  Withdrawal Whitelist
                </CardTitle>
                <Button onClick={handleAddWhitelistAddress}>
                  Add Address
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {whitelistAddresses.map((addr) => (
                  <div key={addr.id} className="flex items-center justify-between p-4 border rounded">
                    <div>
                      <div className="font-medium">{addr.label}</div>
                      <div className="text-sm text-muted-foreground font-mono">
                        {addr.address.slice(0, 10)}...{addr.address.slice(-8)}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {addr.network} • Added {formatTime(addr.addedAt)}
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={addr.isActive ? 'default' : 'secondary'}>
                        {addr.isActive ? 'Active' : 'Inactive'}
                      </Badge>
                      <Button variant="outline" size="sm">
                        Edit
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* API Key Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Key className="w-5 h-5" />
                  API Keys
                </CardTitle>
                <Button onClick={handleCreateAPIKey}>
                  Create API Key
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {apiKeys.map((key) => (
                  <div key={key.id} className="p-4 border rounded">
                    <div className="flex items-center justify-between mb-3">
                      <div>
                        <div className="font-medium">{key.name}</div>
                        <div className="text-sm text-muted-foreground font-mono">
                          {key.key}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant={key.isActive ? 'default' : 'secondary'}>
                          {key.isActive ? 'Active' : 'Inactive'}
                        </Badge>
                        <Button variant="outline" size="sm">
                          Manage
                        </Button>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <span className="text-muted-foreground">Permissions: </span>
                        {key.permissions.join(', ')}
                      </div>
                      <div>
                        <span className="text-muted-foreground">Last Used: </span>
                        {key.lastUsed ? formatTime(key.lastUsed) : 'Never'}
                      </div>
                    </div>

                    <div className="mt-2 text-sm">
                      <span className="text-muted-foreground">IP Whitelist: </span>
                      {key.ipWhitelist.join(', ')}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="activity" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="w-5 h-5" />
                Security Activity Log
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {securityEvents.map((event) => (
                  <div key={event.id} className="flex items-start gap-4 p-4 border rounded">
                    <div className="mt-1">
                      {getStatusIcon(event.status)}
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-2">
                        <div className="font-medium capitalize">
                          {event.type.replace('_', ' ')}
                        </div>
                        <Badge variant={getRiskBadgeVariant(event.riskLevel)}>
                          {event.riskLevel} risk
                        </Badge>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm text-muted-foreground">
                        <div className="flex items-center gap-1">
                          <MapPin className="w-3 h-3" />
                          {event.location}
                        </div>
                        <div className="flex items-center gap-1">
                          <Monitor className="w-3 h-3" />
                          {event.device}
                        </div>
                        <div className="flex items-center gap-1">
                          <Globe className="w-3 h-3" />
                          {event.ipAddress}
                        </div>
                        <div className="flex items-center gap-1">
                          <Clock className="w-3 h-3" />
                          {formatTime(event.timestamp)}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="mfa">
          <MultiFactorAuth />
        </TabsContent>

        <TabsContent value="monitoring">
          <SecurityMonitoring />
        </TabsContent>

        <TabsContent value="fraud">
          <FraudDetection />
        </TabsContent>
      </Tabs>
    </div>
  )
}
