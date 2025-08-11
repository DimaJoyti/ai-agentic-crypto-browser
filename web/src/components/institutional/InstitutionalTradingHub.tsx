'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import {
  Building2,
  Key,
  Shield,
  BarChart3,
  Users,
  DollarSign,
  Clock,
  Settings,
  Download,
  Upload,
  Eye,
  EyeOff,
  Copy,
  CheckCircle,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Zap,
  Globe,
  Lock
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'
import { PrimeServices } from './PrimeServices'
import { RiskManagement } from './RiskManagement'
import { ComplianceReporting } from './ComplianceReporting'

interface APICredential {
  id: string
  name: string
  apiKey: string
  secretKey: string
  permissions: string[]
  ipWhitelist: string[]
  rateLimit: number
  isActive: boolean
  createdAt: number
  lastUsed?: number
  usage24h: number
  environment: 'sandbox' | 'production'
  webhookUrl?: string
  expiresAt?: number
}

interface OTCOrder {
  id: string
  type: 'buy' | 'sell'
  asset: string
  amount: number
  price: number
  counterparty: string
  status: 'pending' | 'negotiating' | 'confirmed' | 'settled'
  createdAt: number
  settleBy: number
}

interface CustodyAccount {
  id: string
  name: string
  type: 'hot' | 'cold' | 'multi_sig'
  assets: Array<{
    symbol: string
    balance: number
    value: number
  }>
  securityLevel: 'standard' | 'enhanced' | 'institutional'
  signatories: string[]
  requiredSignatures: number
}

export function InstitutionalTradingHub() {
  const [apiCredentials, setApiCredentials] = useState<APICredential[]>([])
  const [otcOrders, setOtcOrders] = useState<OTCOrder[]>([])
  const [custodyAccounts, setCustodyAccounts] = useState<CustodyAccount[]>([])
  const [activeTab, setActiveTab] = useState('overview')
  const [showSecrets, setShowSecrets] = useState<Record<string, boolean>>({})

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock institutional data
    const mockApiCredentials: APICredential[] = [
      {
        id: 'api1',
        name: 'Trading Bot Production',
        apiKey: 'demo_key_1234567890abcdef1234567890abcdef',
        secretKey: 'demo_secret_abcdef1234567890abcdef1234567890',
        permissions: ['read', 'trade', 'withdraw'],
        ipWhitelist: ['192.168.1.100', '10.0.0.1'],
        rateLimit: 1000,
        isActive: true,
        createdAt: Date.now() - 86400000 * 30,
        lastUsed: Date.now() - 3600000,
        usage24h: 847,
        environment: 'production',
        webhookUrl: 'https://api.institution.com/webhooks/trading',
        expiresAt: Date.now() + 365 * 24 * 60 * 60 * 1000
      },
      {
        id: 'api2',
        name: 'Market Data Feed',
        apiKey: 'demo_key_fedcba0987654321fedcba0987654321',
        secretKey: 'demo_secret_0987654321fedcba0987654321fedcba',
        permissions: ['read'],
        ipWhitelist: ['203.0.113.1', '203.0.113.2'],
        rateLimit: 5000,
        isActive: true,
        createdAt: Date.now() - 86400000 * 7,
        lastUsed: Date.now() - 300000,
        usage24h: 4523,
        environment: 'production',
        webhookUrl: 'https://api.institution.com/webhooks/data'
      }
    ]

    const mockOtcOrders: OTCOrder[] = [
      {
        id: 'otc1',
        type: 'buy',
        asset: 'BTC',
        amount: 50,
        price: 45000,
        counterparty: 'Institution ABC',
        status: 'negotiating',
        createdAt: Date.now() - 3600000,
        settleBy: Date.now() + 86400000
      },
      {
        id: 'otc2',
        type: 'sell',
        asset: 'ETH',
        amount: 1000,
        price: 2520,
        counterparty: 'Fund XYZ',
        status: 'confirmed',
        createdAt: Date.now() - 7200000,
        settleBy: Date.now() + 43200000
      }
    ]

    const mockCustodyAccounts: CustodyAccount[] = [
      {
        id: 'custody1',
        name: 'Main Trading Account',
        type: 'hot',
        assets: [
          { symbol: 'BTC', balance: 125.5, value: 5647500 },
          { symbol: 'ETH', balance: 2840.2, value: 7100500 },
          { symbol: 'USDT', balance: 1500000, value: 1500000 }
        ],
        securityLevel: 'enhanced',
        signatories: ['0x123...abc', '0x456...def'],
        requiredSignatures: 1
      },
      {
        id: 'custody2',
        name: 'Cold Storage Vault',
        type: 'cold',
        assets: [
          { symbol: 'BTC', balance: 500.0, value: 22500000 },
          { symbol: 'ETH', balance: 8000.0, value: 20000000 }
        ],
        securityLevel: 'institutional',
        signatories: ['0x123...abc', '0x456...def', '0x789...ghi'],
        requiredSignatures: 2
      }
    ]

    setApiCredentials(mockApiCredentials)
    setOtcOrders(mockOtcOrders)
    setCustodyAccounts(mockCustodyAccounts)
  }, [isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const toggleSecretVisibility = (id: string) => {
    setShowSecrets(prev => ({ ...prev, [id]: !prev[id] }))
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    // In real implementation, show toast notification
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': case 'confirmed': case 'settled': return 'text-green-500'
      case 'pending': case 'negotiating': return 'text-yellow-500'
      case 'inactive': case 'cancelled': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'active': case 'confirmed': case 'settled': return 'default'
      case 'pending': case 'negotiating': return 'secondary'
      case 'inactive': case 'cancelled': return 'destructive'
      default: return 'outline'
    }
  }

  const getTotalCustodyValue = () => {
    return custodyAccounts.reduce((total, account) => 
      total + account.assets.reduce((sum, asset) => sum + asset.value, 0), 0
    )
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Building2 className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Institutional Access Required</h3>
          <p className="text-muted-foreground">
            Connect your institutional wallet to access advanced trading features
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Institutional Trading</h2>
          <p className="text-muted-foreground">
            Enterprise-grade trading tools and custody solutions
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Institutional Grade
          </Badge>
          <Badge variant="outline">
            <Lock className="w-3 h-3 mr-1" />
            SOC 2 Compliant
          </Badge>
        </div>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total AUM</span>
            </div>
            <div className="text-2xl font-bold">{formatCurrency(getTotalCustodyValue())}</div>
            <div className="text-xs text-green-500">
              Across {custodyAccounts.length} accounts
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Key className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">API Keys</span>
            </div>
            <div className="text-2xl font-bold">{apiCredentials.filter(a => a.isActive).length}</div>
            <div className="text-xs text-muted-foreground">
              Active credentials
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <BarChart3 className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">OTC Volume</span>
            </div>
            <div className="text-2xl font-bold">
              {formatCurrency(otcOrders.reduce((sum, order) => sum + (order.amount * order.price), 0))}
            </div>
            <div className="text-xs text-muted-foreground">
              This month
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">API Usage</span>
            </div>
            <div className="text-2xl font-bold">
              {apiCredentials.reduce((sum, api) => sum + api.usage24h, 0).toLocaleString()}
            </div>
            <div className="text-xs text-muted-foreground">
              Requests (24h)
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-7">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="api">API Management</TabsTrigger>
          <TabsTrigger value="otc">OTC Desk</TabsTrigger>
          <TabsTrigger value="custody">Custody</TabsTrigger>
          <TabsTrigger value="prime">Prime Services</TabsTrigger>
          <TabsTrigger value="risk">Risk Management</TabsTrigger>
          <TabsTrigger value="compliance">Compliance</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Recent Activity */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="w-5 h-5" />
                  Recent Activity
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 border rounded">
                    <div className="flex items-center gap-3">
                      <Key className="w-4 h-4 text-blue-500" />
                      <div>
                        <div className="font-medium">API Key Used</div>
                        <div className="text-sm text-muted-foreground">Trading Bot Production</div>
                      </div>
                    </div>
                    <div className="text-xs text-muted-foreground">2 min ago</div>
                  </div>
                  
                  <div className="flex items-center justify-between p-3 border rounded">
                    <div className="flex items-center gap-3">
                      <BarChart3 className="w-4 h-4 text-green-500" />
                      <div>
                        <div className="font-medium">OTC Order Confirmed</div>
                        <div className="text-sm text-muted-foreground">1000 ETH @ $2,520</div>
                      </div>
                    </div>
                    <div className="text-xs text-muted-foreground">1 hour ago</div>
                  </div>
                  
                  <div className="flex items-center justify-between p-3 border rounded">
                    <div className="flex items-center gap-3">
                      <Shield className="w-4 h-4 text-orange-500" />
                      <div>
                        <div className="font-medium">Custody Transfer</div>
                        <div className="text-sm text-muted-foreground">50 BTC to Cold Storage</div>
                      </div>
                    </div>
                    <div className="text-xs text-muted-foreground">3 hours ago</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Performance Metrics */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5" />
                  Performance Metrics
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">API Uptime</span>
                    <span className="font-bold text-green-500">99.98%</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Avg Response Time</span>
                    <span className="font-bold">12ms</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">OTC Success Rate</span>
                    <span className="font-bold text-green-500">94.2%</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Settlement Time</span>
                    <span className="font-bold">2.3 hours</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="api" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">API Credentials</h3>
            <Button>
              <Key className="w-4 h-4 mr-2" />
              Create New API Key
            </Button>
          </div>

          <div className="space-y-4">
            {apiCredentials.map((credential) => (
              <Card key={credential.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{credential.name}</h4>
                      <p className="text-sm text-muted-foreground">
                        Created {formatTime(credential.createdAt)}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={getStatusBadgeVariant(credential.isActive ? 'active' : 'inactive')}>
                        {credential.isActive ? 'Active' : 'Inactive'}
                      </Badge>
                      <Button variant="outline" size="sm">
                        <Settings className="w-3 h-3 mr-1" />
                        Manage
                      </Button>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                    <div>
                      <Label className="text-xs text-muted-foreground">API Key</Label>
                      <div className="flex items-center gap-2 mt-1">
                        <code className="flex-1 p-2 bg-muted rounded text-xs font-mono">
                          {showSecrets[credential.id] ? credential.apiKey : '••••••••••••••••••••••••••••••••'}
                        </code>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => toggleSecretVisibility(credential.id)}
                        >
                          {showSecrets[credential.id] ? <EyeOff className="w-3 h-3" /> : <Eye className="w-3 h-3" />}
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => copyToClipboard(credential.apiKey)}
                        >
                          <Copy className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>

                    <div>
                      <Label className="text-xs text-muted-foreground">Secret Key</Label>
                      <div className="flex items-center gap-2 mt-1">
                        <code className="flex-1 p-2 bg-muted rounded text-xs font-mono">
                          {showSecrets[credential.id] ? credential.secretKey : '••••••••••••••••••••••••••••••••'}
                        </code>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => copyToClipboard(credential.secretKey)}
                        >
                          <Copy className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Permissions</div>
                      <div className="font-medium">{credential.permissions.join(', ')}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Rate Limit</div>
                      <div className="font-medium">{credential.rateLimit}/min</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">24h Usage</div>
                      <div className="font-medium">{credential.usage24h.toLocaleString()}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Last Used</div>
                      <div className="font-medium">
                        {credential.lastUsed ? formatTime(credential.lastUsed) : 'Never'}
                      </div>
                    </div>
                  </div>

                  <div className="mt-4">
                    <div className="text-sm text-muted-foreground mb-2">IP Whitelist</div>
                    <div className="flex flex-wrap gap-2">
                      {credential.ipWhitelist.map((ip) => (
                        <Badge key={ip} variant="outline" className="text-xs">
                          {ip}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="otc" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">OTC Trading Desk</h3>
            <Button>
              <BarChart3 className="w-4 h-4 mr-2" />
              Request Quote
            </Button>
          </div>

          <div className="space-y-4">
            {otcOrders.map((order) => (
              <Card key={order.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className={cn(
                        "w-10 h-10 rounded-full flex items-center justify-center",
                        order.type === 'buy' ? "bg-green-100 text-green-600" : "bg-red-100 text-red-600"
                      )}>
                        {order.type === 'buy' ? <TrendingUp className="w-5 h-5" /> : <TrendingDown className="w-5 h-5" />}
                      </div>
                      <div>
                        <h4 className="font-bold">
                          {order.type.toUpperCase()} {order.amount} {order.asset}
                        </h4>
                        <p className="text-sm text-muted-foreground">
                          @ {formatCurrency(order.price)} per {order.asset}
                        </p>
                      </div>
                    </div>
                    <Badge variant={getStatusBadgeVariant(order.status)}>
                      {order.status}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Counterparty</div>
                      <div className="font-medium">{order.counterparty}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Total Value</div>
                      <div className="font-medium">{formatCurrency(order.amount * order.price)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Created</div>
                      <div className="font-medium">{formatTime(order.createdAt)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Settle By</div>
                      <div className="font-medium">{formatTime(order.settleBy)}</div>
                    </div>
                  </div>

                  <div className="mt-4 flex gap-2">
                    {order.status === 'negotiating' && (
                      <>
                        <Button size="sm" variant="outline">
                          Counter Offer
                        </Button>
                        <Button size="sm">
                          Accept Terms
                        </Button>
                      </>
                    )}
                    {order.status === 'confirmed' && (
                      <Button size="sm">
                        Initiate Settlement
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="custody" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Custody Accounts</h3>
            <Button>
              <Shield className="w-4 h-4 mr-2" />
              Create Account
            </Button>
          </div>

          <div className="space-y-4">
            {custodyAccounts.map((account) => (
              <Card key={account.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className={cn(
                        "w-10 h-10 rounded-full flex items-center justify-center",
                        account.type === 'hot' ? "bg-orange-100 text-orange-600" :
                        account.type === 'cold' ? "bg-blue-100 text-blue-600" :
                        "bg-purple-100 text-purple-600"
                      )}>
                        <Shield className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{account.name}</h4>
                        <p className="text-sm text-muted-foreground capitalize">
                          {account.type.replace('_', ' ')} • {account.securityLevel}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-bold">
                        {formatCurrency(account.assets.reduce((sum, asset) => sum + asset.value, 0))}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {account.requiredSignatures}/{account.signatories.length} signatures required
                      </div>
                    </div>
                  </div>

                  <div className="space-y-3">
                    {account.assets.map((asset) => (
                      <div key={asset.symbol} className="flex items-center justify-between p-3 bg-muted/50 rounded">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                            <span className="text-xs font-bold">{asset.symbol}</span>
                          </div>
                          <div>
                            <div className="font-medium">{asset.balance} {asset.symbol}</div>
                            <div className="text-sm text-muted-foreground">
                              {formatCurrency(asset.value)}
                            </div>
                          </div>
                        </div>
                        <Button variant="outline" size="sm">
                          Transfer
                        </Button>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="prime">
          <PrimeServices />
        </TabsContent>

        <TabsContent value="risk">
          <RiskManagement />
        </TabsContent>

        <TabsContent value="compliance">
          <ComplianceReporting />
        </TabsContent>
      </Tabs>
    </div>
  )
}
