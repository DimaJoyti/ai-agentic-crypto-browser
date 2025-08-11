'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  CreditCard, 
  Building2, 
  ArrowUpDown,
  DollarSign,
  Euro,
  PoundSterling,
  Clock,
  CheckCircle,
  AlertTriangle,
  Shield,
  Zap,
  Globe,
  TrendingUp,
  Download,
  Upload,
  Eye,
  Settings,
  Plus,
  Minus,
  RefreshCw
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface FiatCurrency {
  code: string
  name: string
  symbol: string
  icon: React.ReactNode
  supported: boolean
  buyEnabled: boolean
  sellEnabled: boolean
  minAmount: number
  maxAmount: number
  fees: {
    buy: number
    sell: number
    withdrawal: number
  }
}

interface PaymentMethod {
  id: string
  type: 'credit_card' | 'debit_card' | 'bank_transfer' | 'sepa' | 'ach' | 'wire'
  name: string
  description: string
  icon: React.ReactNode
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
  supported: boolean
  instant: boolean
}

interface FiatTransaction {
  id: string
  type: 'buy' | 'sell'
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled'
  fiatAmount: number
  fiatCurrency: string
  cryptoAmount: number
  cryptoCurrency: string
  paymentMethod: string
  fees: number
  exchangeRate: number
  createdAt: number
  completedAt?: number
  reference: string
}

interface FiatProvider {
  id: string
  name: string
  logo: string
  description: string
  rating: number
  fees: {
    buy: number
    sell: number
  }
  limits: {
    daily: number
    monthly: number
  }
  processingTime: string
  supportedCountries: string[]
  supportedCurrencies: string[]
  features: string[]
  isRecommended: boolean
}

export function FiatOnOffRamps() {
  const [activeTab, setActiveTab] = useState('buy')
  const [selectedCurrency, setSelectedCurrency] = useState('USD')
  const [selectedCrypto, setSelectedCrypto] = useState('BTC')
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState('')
  const [amount, setAmount] = useState('')
  const [transactions, setTransactions] = useState<FiatTransaction[]>([])
  const [providers, setProviders] = useState<FiatProvider[]>([])
  const [exchangeRates, setExchangeRates] = useState<Record<string, number>>({})

  const { address, isConnected } = useAccount()

  const fiatCurrencies: FiatCurrency[] = [
    {
      code: 'USD',
      name: 'US Dollar',
      symbol: '$',
      icon: <DollarSign className="w-4 h-4" />,
      supported: true,
      buyEnabled: true,
      sellEnabled: true,
      minAmount: 10,
      maxAmount: 50000,
      fees: { buy: 1.5, sell: 1.2, withdrawal: 5 }
    },
    {
      code: 'EUR',
      name: 'Euro',
      symbol: '€',
      icon: <Euro className="w-4 h-4" />,
      supported: true,
      buyEnabled: true,
      sellEnabled: true,
      minAmount: 10,
      maxAmount: 45000,
      fees: { buy: 1.8, sell: 1.5, withdrawal: 4 }
    },
    {
      code: 'GBP',
      name: 'British Pound',
      symbol: '£',
      icon: <PoundSterling className="w-4 h-4" />,
      supported: true,
      buyEnabled: true,
      sellEnabled: true,
      minAmount: 8,
      maxAmount: 40000,
      fees: { buy: 2.0, sell: 1.8, withdrawal: 3 }
    }
  ]

  const paymentMethods: PaymentMethod[] = [
    {
      id: 'credit_card',
      type: 'credit_card',
      name: 'Credit Card',
      description: 'Instant purchase with Visa/Mastercard',
      icon: <CreditCard className="w-4 h-4" />,
      processingTime: 'Instant',
      fees: { fixed: 0, percentage: 3.5 },
      limits: { min: 10, max: 5000, daily: 10000, monthly: 50000 },
      supported: true,
      instant: true
    },
    {
      id: 'debit_card',
      type: 'debit_card',
      name: 'Debit Card',
      description: 'Direct bank account debit',
      icon: <CreditCard className="w-4 h-4" />,
      processingTime: 'Instant',
      fees: { fixed: 0, percentage: 2.5 },
      limits: { min: 10, max: 3000, daily: 5000, monthly: 25000 },
      supported: true,
      instant: true
    },
    {
      id: 'bank_transfer',
      type: 'bank_transfer',
      name: 'Bank Transfer',
      description: 'Direct bank wire transfer',
      icon: <Building2 className="w-4 h-4" />,
      processingTime: '1-3 business days',
      fees: { fixed: 15, percentage: 0.5 },
      limits: { min: 100, max: 100000, daily: 250000, monthly: 1000000 },
      supported: true,
      instant: false
    },
    {
      id: 'sepa',
      type: 'sepa',
      name: 'SEPA Transfer',
      description: 'European bank transfer',
      icon: <Building2 className="w-4 h-4" />,
      processingTime: '1-2 business days',
      fees: { fixed: 2, percentage: 0.2 },
      limits: { min: 50, max: 50000, daily: 100000, monthly: 500000 },
      supported: true,
      instant: false
    }
  ]

  useEffect(() => {
    if (!isConnected) return

    // Generate mock transaction data
    const mockTransactions: FiatTransaction[] = [
      {
        id: 'tx1',
        type: 'buy',
        status: 'completed',
        fiatAmount: 1000,
        fiatCurrency: 'USD',
        cryptoAmount: 0.022,
        cryptoCurrency: 'BTC',
        paymentMethod: 'Credit Card',
        fees: 35,
        exchangeRate: 45000,
        createdAt: Date.now() - 86400000,
        completedAt: Date.now() - 86400000 + 300000,
        reference: 'REF-001-BUY'
      },
      {
        id: 'tx2',
        type: 'sell',
        status: 'processing',
        fiatAmount: 2500,
        fiatCurrency: 'EUR',
        cryptoAmount: 1.0,
        cryptoCurrency: 'ETH',
        paymentMethod: 'Bank Transfer',
        fees: 12.5,
        exchangeRate: 2500,
        createdAt: Date.now() - 3600000,
        reference: 'REF-002-SELL'
      },
      {
        id: 'tx3',
        type: 'buy',
        status: 'pending',
        fiatAmount: 500,
        fiatCurrency: 'GBP',
        cryptoAmount: 0.2,
        cryptoCurrency: 'ETH',
        paymentMethod: 'Debit Card',
        fees: 12.5,
        exchangeRate: 2500,
        createdAt: Date.now() - 1800000,
        reference: 'REF-003-BUY'
      }
    ]

    const mockProviders: FiatProvider[] = [
      {
        id: 'moonpay',
        name: 'MoonPay',
        logo: '/providers/moonpay.png',
        description: 'Leading fiat-to-crypto gateway with global coverage',
        rating: 4.8,
        fees: { buy: 3.5, sell: 2.5 },
        limits: { daily: 10000, monthly: 50000 },
        processingTime: 'Instant',
        supportedCountries: ['US', 'EU', 'UK', 'CA', 'AU'],
        supportedCurrencies: ['USD', 'EUR', 'GBP', 'CAD', 'AUD'],
        features: ['Instant', 'KYC Verified', '24/7 Support'],
        isRecommended: true
      },
      {
        id: 'simplex',
        name: 'Simplex',
        logo: '/providers/simplex.png',
        description: 'Secure and compliant fiat infrastructure',
        rating: 4.6,
        fees: { buy: 4.0, sell: 3.0 },
        limits: { daily: 20000, monthly: 100000 },
        processingTime: 'Instant',
        supportedCountries: ['US', 'EU', 'UK'],
        supportedCurrencies: ['USD', 'EUR', 'GBP'],
        features: ['High Limits', 'Enterprise Grade', 'Fraud Protection'],
        isRecommended: false
      },
      {
        id: 'banxa',
        name: 'Banxa',
        logo: '/providers/banxa.png',
        description: 'Global payment infrastructure for digital assets',
        rating: 4.4,
        fees: { buy: 3.8, sell: 2.8 },
        limits: { daily: 15000, monthly: 75000 },
        processingTime: '5-10 minutes',
        supportedCountries: ['US', 'EU', 'UK', 'AU', 'CA'],
        supportedCurrencies: ['USD', 'EUR', 'GBP', 'AUD', 'CAD'],
        features: ['Multiple Payment Methods', 'Competitive Rates', 'Fast Processing'],
        isRecommended: false
      }
    ]

    const mockExchangeRates = {
      'BTC/USD': 45000,
      'ETH/USD': 2500,
      'BTC/EUR': 42000,
      'ETH/EUR': 2300,
      'BTC/GBP': 36000,
      'ETH/GBP': 2000
    }

    setTransactions(mockTransactions)
    setProviders(mockProviders)
    setExchangeRates(mockExchangeRates)
  }, [isConnected])

  const formatCurrency = (amount: number, currency: string) => {
    const currencyData = fiatCurrencies.find(c => c.code === currency)
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'text-green-500'
      case 'processing': return 'text-blue-500'
      case 'pending': return 'text-yellow-500'
      case 'failed': case 'cancelled': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'completed': return 'default'
      case 'processing': case 'pending': return 'secondary'
      case 'failed': case 'cancelled': return 'destructive'
      default: return 'outline'
    }
  }

  const calculateFees = (amount: number, paymentMethod: PaymentMethod) => {
    return paymentMethod.fees.fixed + (amount * paymentMethod.fees.percentage / 100)
  }

  const calculateTotal = (amount: number, fees: number, type: 'buy' | 'sell') => {
    return type === 'buy' ? amount + fees : amount - fees
  }

  const getExchangeRate = (crypto: string, fiat: string) => {
    return exchangeRates[`${crypto}/${fiat}`] || 0
  }

  const calculateCryptoAmount = (fiatAmount: number, crypto: string, fiat: string) => {
    const rate = getExchangeRate(crypto, fiat)
    return rate > 0 ? fiatAmount / rate : 0
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <CreditCard className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access fiat on/off ramp services
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
          <h2 className="text-2xl font-bold">Fiat On/Off Ramps</h2>
          <p className="text-muted-foreground">
            Buy and sell cryptocurrency with fiat currency
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
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Volume</span>
            </div>
            <div className="text-2xl font-bold">
              {formatCurrency(
                transactions.reduce((sum, tx) => sum + tx.fiatAmount, 0),
                'USD'
              )}
            </div>
            <div className="text-xs text-muted-foreground">
              This month
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <ArrowUpDown className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Transactions</span>
            </div>
            <div className="text-2xl font-bold">{transactions.length}</div>
            <div className="text-xs text-green-500">
              {transactions.filter(tx => tx.status === 'completed').length} completed
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Clock className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Avg Processing</span>
            </div>
            <div className="text-2xl font-bold">2.5 min</div>
            <div className="text-xs text-muted-foreground">
              For instant methods
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Shield className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Success Rate</span>
            </div>
            <div className="text-2xl font-bold text-green-500">99.2%</div>
            <div className="text-xs text-muted-foreground">
              Last 30 days
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="buy">Buy Crypto</TabsTrigger>
          <TabsTrigger value="sell">Sell Crypto</TabsTrigger>
          <TabsTrigger value="history">Transaction History</TabsTrigger>
          <TabsTrigger value="providers">Payment Providers</TabsTrigger>
        </TabsList>

        <TabsContent value="buy" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Buy Form */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Plus className="w-5 h-5" />
                    Buy Cryptocurrency
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Amount Input */}
                  <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="fiat-amount">You Pay</Label>
                        <div className="flex mt-1">
                          <Input
                            id="fiat-amount"
                            type="number"
                            placeholder="0.00"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            className="rounded-r-none"
                          />
                          <Select value={selectedCurrency} onValueChange={setSelectedCurrency}>
                            <SelectTrigger className="w-24 rounded-l-none border-l-0">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              {fiatCurrencies.map((currency) => (
                                <SelectItem key={currency.code} value={currency.code}>
                                  <div className="flex items-center gap-2">
                                    {currency.icon}
                                    {currency.code}
                                  </div>
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>
                      </div>

                      <div>
                        <Label htmlFor="crypto-amount">You Receive</Label>
                        <div className="flex mt-1">
                          <Input
                            id="crypto-amount"
                            type="number"
                            placeholder="0.00000000"
                            value={amount ? calculateCryptoAmount(parseFloat(amount), selectedCrypto, selectedCurrency).toFixed(8) : ''}
                            readOnly
                            className="rounded-r-none bg-muted"
                          />
                          <Select value={selectedCrypto} onValueChange={setSelectedCrypto}>
                            <SelectTrigger className="w-24 rounded-l-none border-l-0">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="BTC">BTC</SelectItem>
                              <SelectItem value="ETH">ETH</SelectItem>
                              <SelectItem value="USDT">USDT</SelectItem>
                              <SelectItem value="BNB">BNB</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      </div>
                    </div>

                    {/* Exchange Rate */}
                    <div className="flex items-center justify-between text-sm text-muted-foreground">
                      <span>Exchange Rate</span>
                      <span>
                        1 {selectedCrypto} = {formatCurrency(getExchangeRate(selectedCrypto, selectedCurrency), selectedCurrency)}
                      </span>
                    </div>
                  </div>

                  {/* Payment Method Selection */}
                  <div>
                    <Label>Payment Method</Label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                      {paymentMethods.filter(method => method.supported).map((method) => (
                        <div
                          key={method.id}
                          className={cn(
                            "p-4 border rounded-lg cursor-pointer transition-colors",
                            selectedPaymentMethod === method.id
                              ? "border-primary bg-primary/5"
                              : "border-border hover:border-primary/50"
                          )}
                          onClick={() => setSelectedPaymentMethod(method.id)}
                        >
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center gap-2">
                              {method.icon}
                              <span className="font-medium">{method.name}</span>
                            </div>
                            {method.instant && (
                              <Badge variant="outline" className="text-xs">
                                <Zap className="w-3 h-3 mr-1" />
                                Instant
                              </Badge>
                            )}
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">
                            {method.description}
                          </p>
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>Fee: {method.fees.percentage}%</span>
                            <span>{method.processingTime}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Transaction Summary */}
                  {amount && selectedPaymentMethod && (
                    <div className="p-4 bg-muted/50 rounded-lg space-y-2">
                      <h4 className="font-medium">Transaction Summary</h4>
                      <div className="space-y-1 text-sm">
                        <div className="flex justify-between">
                          <span>Amount</span>
                          <span>{formatCurrency(parseFloat(amount), selectedCurrency)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span>Fees</span>
                          <span>
                            {formatCurrency(
                              calculateFees(parseFloat(amount), paymentMethods.find(m => m.id === selectedPaymentMethod)!),
                              selectedCurrency
                            )}
                          </span>
                        </div>
                        <div className="flex justify-between font-medium border-t pt-1">
                          <span>Total</span>
                          <span>
                            {formatCurrency(
                              calculateTotal(
                                parseFloat(amount),
                                calculateFees(parseFloat(amount), paymentMethods.find(m => m.id === selectedPaymentMethod)!),
                                'buy'
                              ),
                              selectedCurrency
                            )}
                          </span>
                        </div>
                      </div>
                    </div>
                  )}

                  <Button 
                    className="w-full" 
                    size="lg"
                    disabled={!amount || !selectedPaymentMethod}
                  >
                    <Plus className="w-4 h-4 mr-2" />
                    Buy {selectedCrypto}
                  </Button>
                </CardContent>
              </Card>
            </div>

            {/* Payment Providers */}
            <div>
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Recommended Providers</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {providers.filter(p => p.isRecommended).map((provider) => (
                      <div key={provider.id} className="p-3 border rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center">
                              <span className="text-xs font-bold">{provider.name[0]}</span>
                            </div>
                            <span className="font-medium">{provider.name}</span>
                          </div>
                          <div className="flex items-center gap-1">
                            <span className="text-sm text-yellow-500">★</span>
                            <span className="text-sm">{provider.rating}</span>
                          </div>
                        </div>
                        <p className="text-xs text-muted-foreground mb-2">
                          {provider.description}
                        </p>
                        <div className="flex justify-between text-xs">
                          <span>Fee: {provider.fees.buy}%</span>
                          <span>{provider.processingTime}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="sell" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Sell Form */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Minus className="w-5 h-5" />
                    Sell Cryptocurrency
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Amount Input */}
                  <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="crypto-sell-amount">You Sell</Label>
                        <div className="flex mt-1">
                          <Input
                            id="crypto-sell-amount"
                            type="number"
                            placeholder="0.00000000"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            className="rounded-r-none"
                          />
                          <Select value={selectedCrypto} onValueChange={setSelectedCrypto}>
                            <SelectTrigger className="w-24 rounded-l-none border-l-0">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="BTC">BTC</SelectItem>
                              <SelectItem value="ETH">ETH</SelectItem>
                              <SelectItem value="USDT">USDT</SelectItem>
                              <SelectItem value="BNB">BNB</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      </div>

                      <div>
                        <Label htmlFor="fiat-receive-amount">You Receive</Label>
                        <div className="flex mt-1">
                          <Input
                            id="fiat-receive-amount"
                            type="number"
                            placeholder="0.00"
                            value={amount ? (parseFloat(amount) * getExchangeRate(selectedCrypto, selectedCurrency)).toFixed(2) : ''}
                            readOnly
                            className="rounded-r-none bg-muted"
                          />
                          <Select value={selectedCurrency} onValueChange={setSelectedCurrency}>
                            <SelectTrigger className="w-24 rounded-l-none border-l-0">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              {fiatCurrencies.map((currency) => (
                                <SelectItem key={currency.code} value={currency.code}>
                                  <div className="flex items-center gap-2">
                                    {currency.icon}
                                    {currency.code}
                                  </div>
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>
                      </div>
                    </div>

                    {/* Exchange Rate */}
                    <div className="flex items-center justify-between text-sm text-muted-foreground">
                      <span>Exchange Rate</span>
                      <span>
                        1 {selectedCrypto} = {formatCurrency(getExchangeRate(selectedCrypto, selectedCurrency), selectedCurrency)}
                      </span>
                    </div>
                  </div>

                  {/* Withdrawal Method */}
                  <div>
                    <Label>Withdrawal Method</Label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                      {paymentMethods.filter(method => method.supported && method.type !== 'credit_card').map((method) => (
                        <div
                          key={method.id}
                          className={cn(
                            "p-4 border rounded-lg cursor-pointer transition-colors",
                            selectedPaymentMethod === method.id
                              ? "border-primary bg-primary/5"
                              : "border-border hover:border-primary/50"
                          )}
                          onClick={() => setSelectedPaymentMethod(method.id)}
                        >
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center gap-2">
                              {method.icon}
                              <span className="font-medium">{method.name}</span>
                            </div>
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">
                            {method.description}
                          </p>
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>Fee: {method.fees.percentage}%</span>
                            <span>{method.processingTime}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>

                  <Button 
                    className="w-full" 
                    size="lg"
                    disabled={!amount || !selectedPaymentMethod}
                  >
                    <Minus className="w-4 h-4 mr-2" />
                    Sell {selectedCrypto}
                  </Button>
                </CardContent>
              </Card>
            </div>

            {/* Withdrawal Info */}
            <div>
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Withdrawal Information</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Alert>
                      <AlertTriangle className="h-4 w-4" />
                      <AlertDescription>
                        Withdrawals may take 1-3 business days depending on your bank and withdrawal method.
                      </AlertDescription>
                    </Alert>

                    <div className="space-y-3">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Daily Limit</span>
                        <span className="font-medium">$25,000</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Monthly Limit</span>
                        <span className="font-medium">$100,000</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Processing Time</span>
                        <span className="font-medium">1-3 business days</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Transaction History</h3>
            <Button variant="outline" size="sm">
              <Download className="w-4 h-4 mr-2" />
              Export
            </Button>
          </div>

          <div className="space-y-4">
            {transactions.map((transaction) => (
              <Card key={transaction.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className={cn(
                        "w-10 h-10 rounded-full flex items-center justify-center",
                        transaction.type === 'buy' ? "bg-green-100 text-green-600" : "bg-red-100 text-red-600"
                      )}>
                        {transaction.type === 'buy' ? <Plus className="w-5 h-5" /> : <Minus className="w-5 h-5" />}
                      </div>
                      <div>
                        <h4 className="font-bold capitalize">
                          {transaction.type} {transaction.cryptoCurrency}
                        </h4>
                        <p className="text-sm text-muted-foreground">
                          {transaction.reference}
                        </p>
                      </div>
                    </div>
                    <Badge variant={getStatusBadgeVariant(transaction.status)}>
                      {transaction.status}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Amount</div>
                      <div className="font-medium">
                        {formatCurrency(transaction.fiatAmount, transaction.fiatCurrency)}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Crypto</div>
                      <div className="font-medium">
                        {transaction.cryptoAmount} {transaction.cryptoCurrency}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Payment Method</div>
                      <div className="font-medium">{transaction.paymentMethod}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Date</div>
                      <div className="font-medium">{formatTime(transaction.createdAt)}</div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between mt-4 pt-4 border-t">
                    <div className="text-sm text-muted-foreground">
                      Fees: {formatCurrency(transaction.fees, transaction.fiatCurrency)}
                    </div>
                    <Button variant="outline" size="sm">
                      <Eye className="w-3 h-3 mr-1" />
                      View Details
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="providers" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Payment Providers</h3>
            <Button variant="outline" size="sm">
              <RefreshCw className="w-4 h-4 mr-2" />
              Refresh Rates
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {providers.map((provider) => (
              <Card key={provider.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <span className="font-bold">{provider.name[0]}</span>
                      </div>
                      <div>
                        <h4 className="font-bold">{provider.name}</h4>
                        <div className="flex items-center gap-1">
                          <span className="text-sm text-yellow-500">★</span>
                          <span className="text-sm">{provider.rating}</span>
                        </div>
                      </div>
                    </div>
                    {provider.isRecommended && (
                      <Badge variant="outline">Recommended</Badge>
                    )}
                  </div>

                  <p className="text-sm text-muted-foreground mb-4">
                    {provider.description}
                  </p>

                  <div className="space-y-3">
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <div className="text-muted-foreground">Buy Fee</div>
                        <div className="font-medium">{provider.fees.buy}%</div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Sell Fee</div>
                        <div className="font-medium">{provider.fees.sell}%</div>
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <div className="text-muted-foreground">Daily Limit</div>
                        <div className="font-medium">{formatCurrency(provider.limits.daily, 'USD')}</div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Processing</div>
                        <div className="font-medium">{provider.processingTime}</div>
                      </div>
                    </div>

                    <div>
                      <div className="text-muted-foreground text-sm mb-2">Features</div>
                      <div className="flex flex-wrap gap-1">
                        {provider.features.map((feature) => (
                          <Badge key={feature} variant="outline" className="text-xs">
                            {feature}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </div>

                  <Button className="w-full mt-4" variant="outline">
                    Select Provider
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
