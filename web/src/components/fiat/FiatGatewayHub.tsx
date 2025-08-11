'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { 
  CreditCard, 
  Building2, 
  Shield,
  ArrowUpDown,
  User,
  Settings,
  TrendingUp,
  Globe,
  Zap,
  CheckCircle
} from 'lucide-react'
import { useAccount } from 'wagmi'
import { FiatOnOffRamps } from './FiatOnOffRamps'
import { KYCVerification } from './KYCVerification'
import { BankAccountManagement } from './BankAccountManagement'

export function FiatGatewayHub() {
  const [activeTab, setActiveTab] = useState('ramps')
  const { address, isConnected } = useAccount()

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <CreditCard className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access fiat gateway services
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
          <h2 className="text-2xl font-bold">Fiat Gateway</h2>
          <p className="text-muted-foreground">
            Complete fiat-to-crypto infrastructure with global payment support
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            KYC Verified
          </Badge>
          <Badge variant="outline">
            <Globe className="w-3 h-3 mr-1" />
            Global Coverage
          </Badge>
          <Badge variant="outline">
            <CheckCircle className="w-3 h-3 mr-1" />
            Bank Grade Security
          </Badge>
        </div>
      </div>

      {/* Quick Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <ArrowUpDown className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Supported Methods</span>
            </div>
            <div className="text-2xl font-bold">12+</div>
            <div className="text-xs text-muted-foreground">
              Payment options
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Globe className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Global Reach</span>
            </div>
            <div className="text-2xl font-bold">180+</div>
            <div className="text-xs text-muted-foreground">
              Countries supported
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Zap className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Processing</span>
            </div>
            <div className="text-2xl font-bold">Instant</div>
            <div className="text-xs text-muted-foreground">
              For most methods
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Success Rate</span>
            </div>
            <div className="text-2xl font-bold text-green-500">99.2%</div>
            <div className="text-xs text-muted-foreground">
              Transaction success
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="ramps" className="flex items-center gap-2">
            <ArrowUpDown className="w-4 h-4" />
            Buy/Sell
          </TabsTrigger>
          <TabsTrigger value="kyc" className="flex items-center gap-2">
            <User className="w-4 h-4" />
            Verification
          </TabsTrigger>
          <TabsTrigger value="payments" className="flex items-center gap-2">
            <Building2 className="w-4 h-4" />
            Payment Methods
          </TabsTrigger>
          <TabsTrigger value="settings" className="flex items-center gap-2">
            <Settings className="w-4 h-4" />
            Settings
          </TabsTrigger>
        </TabsList>

        <TabsContent value="ramps">
          <FiatOnOffRamps />
        </TabsContent>

        <TabsContent value="kyc">
          <KYCVerification />
        </TabsContent>

        <TabsContent value="payments">
          <BankAccountManagement />
        </TabsContent>

        <TabsContent value="settings" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="w-5 h-5" />
                Fiat Gateway Settings
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {/* Notification Preferences */}
                <div>
                  <h4 className="font-medium mb-3">Notification Preferences</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Transaction Confirmations</div>
                        <div className="text-sm text-muted-foreground">
                          Get notified when transactions are completed
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">KYC Status Updates</div>
                        <div className="text-sm text-muted-foreground">
                          Receive updates on verification status
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Security Alerts</div>
                        <div className="text-sm text-muted-foreground">
                          Important security notifications
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                  </div>
                </div>

                {/* Default Settings */}
                <div>
                  <h4 className="font-medium mb-3">Default Settings</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Default Currency</div>
                        <div className="text-sm text-muted-foreground">
                          Primary fiat currency for transactions
                        </div>
                      </div>
                      <Badge variant="outline">USD</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Default Payment Method</div>
                        <div className="text-sm text-muted-foreground">
                          Preferred payment method for purchases
                        </div>
                      </div>
                      <Badge variant="outline">Credit Card</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Auto-Convert</div>
                        <div className="text-sm text-muted-foreground">
                          Automatically convert to preferred crypto
                        </div>
                      </div>
                      <Badge variant="outline">Disabled</Badge>
                    </div>
                  </div>
                </div>

                {/* Security Settings */}
                <div>
                  <h4 className="font-medium mb-3">Security Settings</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Transaction Limits</div>
                        <div className="text-sm text-muted-foreground">
                          Daily and monthly transaction limits
                        </div>
                      </div>
                      <Badge variant="outline">Level 2</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Two-Factor Authentication</div>
                        <div className="text-sm text-muted-foreground">
                          Additional security for fiat transactions
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Withdrawal Whitelist</div>
                        <div className="text-sm text-muted-foreground">
                          Only allow withdrawals to verified accounts
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                  </div>
                </div>

                {/* Provider Preferences */}
                <div>
                  <h4 className="font-medium mb-3">Provider Preferences</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Preferred Provider</div>
                        <div className="text-sm text-muted-foreground">
                          Default payment provider for transactions
                        </div>
                      </div>
                      <Badge variant="outline">MoonPay</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Auto-Select Best Rate</div>
                        <div className="text-sm text-muted-foreground">
                          Automatically choose provider with best rates
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Fee Optimization</div>
                        <div className="text-sm text-muted-foreground">
                          Minimize fees by selecting optimal methods
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                  </div>
                </div>

                {/* Compliance Settings */}
                <div>
                  <h4 className="font-medium mb-3">Compliance & Reporting</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Tax Reporting</div>
                        <div className="text-sm text-muted-foreground">
                          Generate tax reports for fiat transactions
                        </div>
                      </div>
                      <Badge variant="outline">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Transaction History Export</div>
                        <div className="text-sm text-muted-foreground">
                          Export transaction data for accounting
                        </div>
                      </div>
                      <Badge variant="outline">Available</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="font-medium">Regulatory Compliance</div>
                        <div className="text-sm text-muted-foreground">
                          Automatic compliance with local regulations
                        </div>
                      </div>
                      <Badge variant="outline">Active</Badge>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Supported Regions */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Globe className="w-5 h-5" />
                Supported Regions
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <h4 className="font-medium mb-2">North America</h4>
                  <div className="space-y-1 text-sm">
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>United States</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Canada</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Mexico</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Europe</h4>
                  <div className="space-y-1 text-sm">
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>European Union</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>United Kingdom</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Switzerland</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Asia Pacific</h4>
                  <div className="space-y-1 text-sm">
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Australia</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Singapore</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle className="w-3 h-3 text-green-500" />
                      <span>Japan</span>
                    </div>
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
