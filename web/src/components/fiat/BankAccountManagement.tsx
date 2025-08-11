'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Building2, 
  CreditCard,
  CheckCircle,
  AlertTriangle,
  Clock,
  Plus,
  Trash2,
  Edit,
  Eye,
  EyeOff,
  Shield,
  Globe,
  Zap,
  Download,
  Upload,
  Settings
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface BankAccount {
  id: string
  type: 'checking' | 'savings' | 'business'
  bankName: string
  accountName: string
  accountNumber: string
  routingNumber: string
  iban?: string
  swiftCode?: string
  currency: string
  country: string
  status: 'pending' | 'verified' | 'rejected' | 'suspended'
  isDefault: boolean
  addedAt: number
  verifiedAt?: number
  lastUsed?: number
  withdrawalLimits: {
    daily: number
    monthly: number
  }
}

interface PaymentCard {
  id: string
  type: 'credit' | 'debit'
  cardNumber: string
  expiryMonth: string
  expiryYear: string
  cardholderName: string
  brand: 'visa' | 'mastercard' | 'amex' | 'discover'
  currency: string
  country: string
  status: 'active' | 'expired' | 'blocked' | 'pending'
  isDefault: boolean
  addedAt: number
  lastUsed?: number
  limits: {
    daily: number
    monthly: number
  }
}

interface WithdrawalMethod {
  id: string
  type: 'bank_transfer' | 'sepa' | 'ach' | 'wire' | 'card'
  name: string
  description: string
  processingTime: string
  fees: {
    fixed: number
    percentage: number
  }
  limits: {
    min: number
    max: number
    daily: number
    monthly: number
  }
  supportedCurrencies: string[]
  supportedCountries: string[]
  isAvailable: boolean
}

export function BankAccountManagement() {
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([])
  const [paymentCards, setPaymentCards] = useState<PaymentCard[]>([])
  const [withdrawalMethods, setWithdrawalMethods] = useState<WithdrawalMethod[]>([])
  const [showAddAccount, setShowAddAccount] = useState(false)
  const [showAddCard, setShowAddCard] = useState(false)
  const [showAccountNumbers, setShowAccountNumbers] = useState<Record<string, boolean>>({})

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock bank account data
    const mockBankAccounts: BankAccount[] = [
      {
        id: 'bank1',
        type: 'checking',
        bankName: 'Chase Bank',
        accountName: 'John Doe',
        accountNumber: '****1234',
        routingNumber: '021000021',
        currency: 'USD',
        country: 'US',
        status: 'verified',
        isDefault: true,
        addedAt: Date.now() - 86400000 * 30,
        verifiedAt: Date.now() - 86400000 * 28,
        lastUsed: Date.now() - 86400000 * 2,
        withdrawalLimits: { daily: 25000, monthly: 100000 }
      },
      {
        id: 'bank2',
        type: 'savings',
        bankName: 'Bank of America',
        accountName: 'John Doe',
        accountNumber: '****5678',
        routingNumber: '026009593',
        currency: 'USD',
        country: 'US',
        status: 'pending',
        isDefault: false,
        addedAt: Date.now() - 86400000 * 2,
        withdrawalLimits: { daily: 10000, monthly: 50000 }
      },
      {
        id: 'bank3',
        type: 'checking',
        bankName: 'HSBC UK',
        accountName: 'John Doe',
        accountNumber: 'GB29****1234',
        routingNumber: '',
        iban: 'GB29NWBK60161331926819',
        swiftCode: 'HBUKGB4B',
        currency: 'GBP',
        country: 'UK',
        status: 'verified',
        isDefault: false,
        addedAt: Date.now() - 86400000 * 15,
        verifiedAt: Date.now() - 86400000 * 13,
        withdrawalLimits: { daily: 20000, monthly: 80000 }
      }
    ]

    const mockPaymentCards: PaymentCard[] = [
      {
        id: 'card1',
        type: 'credit',
        cardNumber: '****1234',
        expiryMonth: '12',
        expiryYear: '2026',
        cardholderName: 'John Doe',
        brand: 'visa',
        currency: 'USD',
        country: 'US',
        status: 'active',
        isDefault: true,
        addedAt: Date.now() - 86400000 * 60,
        lastUsed: Date.now() - 86400000,
        limits: { daily: 5000, monthly: 25000 }
      },
      {
        id: 'card2',
        type: 'debit',
        cardNumber: '****5678',
        expiryMonth: '08',
        expiryYear: '2025',
        cardholderName: 'John Doe',
        brand: 'mastercard',
        currency: 'EUR',
        country: 'DE',
        status: 'active',
        isDefault: false,
        addedAt: Date.now() - 86400000 * 45,
        lastUsed: Date.now() - 86400000 * 5,
        limits: { daily: 3000, monthly: 15000 }
      }
    ]

    const mockWithdrawalMethods: WithdrawalMethod[] = [
      {
        id: 'bank_transfer',
        type: 'bank_transfer',
        name: 'Bank Transfer',
        description: 'Direct transfer to your bank account',
        processingTime: '1-3 business days',
        fees: { fixed: 15, percentage: 0.5 },
        limits: { min: 100, max: 100000, daily: 250000, monthly: 1000000 },
        supportedCurrencies: ['USD', 'EUR', 'GBP'],
        supportedCountries: ['US', 'EU', 'UK'],
        isAvailable: true
      },
      {
        id: 'sepa',
        type: 'sepa',
        name: 'SEPA Transfer',
        description: 'European bank transfer',
        processingTime: '1-2 business days',
        fees: { fixed: 2, percentage: 0.2 },
        limits: { min: 50, max: 50000, daily: 100000, monthly: 500000 },
        supportedCurrencies: ['EUR'],
        supportedCountries: ['EU'],
        isAvailable: true
      },
      {
        id: 'ach',
        type: 'ach',
        name: 'ACH Transfer',
        description: 'US domestic bank transfer',
        processingTime: '2-3 business days',
        fees: { fixed: 5, percentage: 0.1 },
        limits: { min: 25, max: 25000, daily: 50000, monthly: 200000 },
        supportedCurrencies: ['USD'],
        supportedCountries: ['US'],
        isAvailable: true
      },
      {
        id: 'wire',
        type: 'wire',
        name: 'Wire Transfer',
        description: 'International wire transfer',
        processingTime: '1-5 business days',
        fees: { fixed: 50, percentage: 1.0 },
        limits: { min: 1000, max: 500000, daily: 1000000, monthly: 5000000 },
        supportedCurrencies: ['USD', 'EUR', 'GBP', 'JPY'],
        supportedCountries: ['Global'],
        isAvailable: true
      }
    ]

    setBankAccounts(mockBankAccounts)
    setPaymentCards(mockPaymentCards)
    setWithdrawalMethods(mockWithdrawalMethods)
  }, [isConnected])

  const formatCurrency = (amount: number, currency: string = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'verified': case 'active': return 'text-green-500'
      case 'pending': return 'text-yellow-500'
      case 'rejected': case 'blocked': case 'expired': case 'suspended': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'verified': case 'active': return 'default'
      case 'pending': return 'secondary'
      case 'rejected': case 'blocked': case 'expired': case 'suspended': return 'destructive'
      default: return 'outline'
    }
  }

  const getBrandIcon = (brand: string) => {
    // In a real implementation, you would use actual brand icons
    return <CreditCard className="w-4 h-4" />
  }

  const toggleAccountVisibility = (accountId: string) => {
    setShowAccountNumbers(prev => ({ ...prev, [accountId]: !prev[accountId] }))
  }

  const setDefaultAccount = (accountId: string) => {
    setBankAccounts(prev => prev.map(account => ({
      ...account,
      isDefault: account.id === accountId
    })))
  }

  const setDefaultCard = (cardId: string) => {
    setPaymentCards(prev => prev.map(card => ({
      ...card,
      isDefault: card.id === cardId
    })))
  }

  const removeAccount = (accountId: string) => {
    setBankAccounts(prev => prev.filter(account => account.id !== accountId))
  }

  const removeCard = (cardId: string) => {
    setPaymentCards(prev => prev.filter(card => card.id !== cardId))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Building2 className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to manage bank accounts and payment methods
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
          <h2 className="text-2xl font-bold">Payment Methods</h2>
          <p className="text-muted-foreground">
            Manage your bank accounts and payment cards for fiat transactions
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Bank Grade Security
          </Badge>
          <Badge variant="outline">
            <Globe className="w-3 h-3 mr-1" />
            Global Support
          </Badge>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Building2 className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Bank Accounts</span>
            </div>
            <div className="text-2xl font-bold">{bankAccounts.length}</div>
            <div className="text-xs text-green-500">
              {bankAccounts.filter(a => a.status === 'verified').length} verified
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <CreditCard className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Payment Cards</span>
            </div>
            <div className="text-2xl font-bold">{paymentCards.length}</div>
            <div className="text-xs text-green-500">
              {paymentCards.filter(c => c.status === 'active').length} active
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Zap className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Daily Limit</span>
            </div>
            <div className="text-2xl font-bold">
              {formatCurrency(
                Math.max(...bankAccounts.map(a => a.withdrawalLimits.daily))
              )}
            </div>
            <div className="text-xs text-muted-foreground">
              Maximum available
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Clock className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Processing</span>
            </div>
            <div className="text-2xl font-bold">1-3</div>
            <div className="text-xs text-muted-foreground">
              Business days
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Bank Accounts Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Building2 className="w-5 h-5" />
              Bank Accounts
            </CardTitle>
            <Button onClick={() => setShowAddAccount(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Add Bank Account
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {bankAccounts.map((account) => (
              <div key={account.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                      <Building2 className="w-5 h-5" />
                    </div>
                    <div>
                      <h4 className="font-bold">{account.bankName}</h4>
                      <p className="text-sm text-muted-foreground">
                        {account.accountName} • {account.type} • {account.currency}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {account.isDefault && (
                      <Badge variant="outline">Default</Badge>
                    )}
                    <Badge variant={getStatusBadgeVariant(account.status)}>
                      {account.status}
                    </Badge>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                  <div>
                    <div className="text-muted-foreground">Account Number</div>
                    <div className="flex items-center gap-2">
                      <span className="font-mono">
                        {showAccountNumbers[account.id] ? account.accountNumber : '****' + account.accountNumber.slice(-4)}
                      </span>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => toggleAccountVisibility(account.id)}
                      >
                        {showAccountNumbers[account.id] ? <EyeOff className="w-3 h-3" /> : <Eye className="w-3 h-3" />}
                      </Button>
                    </div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Routing/SWIFT</div>
                    <div className="font-mono">{account.swiftCode || account.routingNumber}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Daily Limit</div>
                    <div className="font-medium">{formatCurrency(account.withdrawalLimits.daily, account.currency)}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Added</div>
                    <div className="font-medium">{formatTime(account.addedAt)}</div>
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="text-sm text-muted-foreground">
                    {account.lastUsed && `Last used: ${formatTime(account.lastUsed)}`}
                  </div>
                  <div className="flex gap-2">
                    {!account.isDefault && account.status === 'verified' && (
                      <Button 
                        variant="outline" 
                        size="sm"
                        onClick={() => setDefaultAccount(account.id)}
                      >
                        Set Default
                      </Button>
                    )}
                    <Button variant="outline" size="sm">
                      <Edit className="w-3 h-3 mr-1" />
                      Edit
                    </Button>
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => removeAccount(account.id)}
                    >
                      <Trash2 className="w-3 h-3 mr-1" />
                      Remove
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Payment Cards Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <CreditCard className="w-5 h-5" />
              Payment Cards
            </CardTitle>
            <Button onClick={() => setShowAddCard(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Add Payment Card
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {paymentCards.map((card) => (
              <div key={card.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                      {getBrandIcon(card.brand)}
                    </div>
                    <div>
                      <h4 className="font-bold capitalize">
                        {card.brand} {card.type} Card
                      </h4>
                      <p className="text-sm text-muted-foreground">
                        {card.cardholderName} • {card.currency}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {card.isDefault && (
                      <Badge variant="outline">Default</Badge>
                    )}
                    <Badge variant={getStatusBadgeVariant(card.status)}>
                      {card.status}
                    </Badge>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                  <div>
                    <div className="text-muted-foreground">Card Number</div>
                    <div className="font-mono">{card.cardNumber}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Expires</div>
                    <div className="font-mono">{card.expiryMonth}/{card.expiryYear}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Daily Limit</div>
                    <div className="font-medium">{formatCurrency(card.limits.daily, card.currency)}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Added</div>
                    <div className="font-medium">{formatTime(card.addedAt)}</div>
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="text-sm text-muted-foreground">
                    {card.lastUsed && `Last used: ${formatTime(card.lastUsed)}`}
                  </div>
                  <div className="flex gap-2">
                    {!card.isDefault && card.status === 'active' && (
                      <Button 
                        variant="outline" 
                        size="sm"
                        onClick={() => setDefaultCard(card.id)}
                      >
                        Set Default
                      </Button>
                    )}
                    <Button variant="outline" size="sm">
                      <Edit className="w-3 h-3 mr-1" />
                      Edit
                    </Button>
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => removeCard(card.id)}
                    >
                      <Trash2 className="w-3 h-3 mr-1" />
                      Remove
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Withdrawal Methods */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Download className="w-5 h-5" />
            Withdrawal Methods
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {withdrawalMethods.map((method) => (
              <div key={method.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <Building2 className="w-4 h-4" />
                    <span className="font-medium">{method.name}</span>
                  </div>
                  {method.isAvailable && (
                    <Badge variant="outline">Available</Badge>
                  )}
                </div>

                <p className="text-sm text-muted-foreground mb-3">
                  {method.description}
                </p>

                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Processing Time</span>
                    <span>{method.processingTime}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Fees</span>
                    <span>{method.fees.percentage}% + {formatCurrency(method.fees.fixed)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Daily Limit</span>
                    <span>{formatCurrency(method.limits.daily)}</span>
                  </div>
                </div>

                <div className="mt-3">
                  <div className="text-xs text-muted-foreground mb-1">Supported Currencies</div>
                  <div className="flex flex-wrap gap-1">
                    {method.supportedCurrencies.map((currency) => (
                      <Badge key={currency} variant="outline" className="text-xs">
                        {currency}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Security Notice */}
      <Alert>
        <Shield className="h-4 w-4" />
        <AlertDescription>
          All payment methods are secured with bank-grade encryption and undergo strict verification processes. 
          Your financial information is never stored on our servers and is processed through certified payment processors.
        </AlertDescription>
      </Alert>
    </div>
  )
}
