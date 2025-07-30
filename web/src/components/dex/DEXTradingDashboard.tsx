'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  ArrowUpDown,
  TrendingUp,
  TrendingDown,
  Zap,
  RefreshCw,
  Settings,
  AlertTriangle,
  CheckCircle,
  Clock,
  DollarSign,
  Activity,
  BarChart3,
  ArrowRight,
  ExternalLink,
  Info
} from 'lucide-react'
import { useDEX, useTokenSwap, useDEXComparison, useDEXAnalytics } from '@/hooks/useDEX'
import { type Token, type SwapQuote, type DEXProtocol } from '@/lib/dex-integration'
import { cn } from '@/lib/utils'

export function DEXTradingDashboard() {
  const [activeTab, setActiveTab] = useState('swap')
  const [tokenIn, setTokenIn] = useState<Token | null>(null)
  const [tokenOut, setTokenOut] = useState<Token | null>(null)
  const [amountIn, setAmountIn] = useState('')
  const [slippage, setSlippage] = useState(0.5)
  const [selectedQuote, setSelectedQuote] = useState<SwapQuote | null>(null)
  const [quotes, setQuotes] = useState<SwapQuote[]>([])

  const {
    state,
    getSwapQuote,
    executeSwap,
    clearError
  } = useDEX({
    autoLoad: true,
    enableNotifications: true,
    defaultSlippage: 0.5
  })

  const { compareRates } = useDEXComparison()
  const analytics = useDEXAnalytics()

  // Mock tokens for demo
  const mockTokens: Token[] = [
    {
      address: '0xA0b86a33E6441E6C8C7F1C7C8C7F1C7C8C7F1C7C',
      symbol: 'ETH',
      name: 'Ethereum',
      decimals: 18,
      chainId: 1,
      verified: true,
      priceUSD: 1800
    },
    {
      address: '0xA0b86a33E6441E6C8C7F1C7C8C7F1C7C8C7F1C7D',
      symbol: 'USDC',
      name: 'USD Coin',
      decimals: 6,
      chainId: 1,
      verified: true,
      priceUSD: 1
    },
    {
      address: '0xA0b86a33E6441E6C8C7F1C7C8C7F1C7C8C7F1C7E',
      symbol: 'WBTC',
      name: 'Wrapped Bitcoin',
      decimals: 8,
      chainId: 1,
      verified: true,
      priceUSD: 35000
    }
  ]

  useEffect(() => {
    if (!tokenIn && !tokenOut) {
      setTokenIn(mockTokens[0])
      setTokenOut(mockTokens[1])
    }
  }, [])

  const handleGetQuotes = async () => {
    if (!tokenIn || !tokenOut || !amountIn) return

    try {
      const newQuotes = await getSwapQuote(tokenIn, tokenOut, amountIn, slippage)
      setQuotes(newQuotes)
      if (newQuotes.length > 0) {
        setSelectedQuote(newQuotes[0])
      }
    } catch (error) {
      console.error('Failed to get quotes:', error)
    }
  }

  const handleSwapTokens = () => {
    const temp = tokenIn
    setTokenIn(tokenOut)
    setTokenOut(temp)
    setQuotes([])
    setSelectedQuote(null)
  }

  const handleExecuteSwap = async () => {
    if (!selectedQuote) return

    try {
      await executeSwap(selectedQuote)
      setQuotes([])
      setSelectedQuote(null)
      setAmountIn('')
    } catch (error) {
      console.error('Swap failed:', error)
    }
  }

  const formatNumber = (value: string | number, decimals: number = 4) => {
    return parseFloat(value.toString()).toFixed(decimals)
  }

  const formatPercentage = (value: number) => {
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  }

  const getProtocolLogo = (dexId: string) => {
    const logos: Record<string, string> = {
      'uniswap-v3': '🦄',
      'sushiswap': '🍣',
      '1inch': '1️⃣'
    }
    return logos[dexId] || '🔄'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">DEX Trading</h2>
          <p className="text-muted-foreground">
            Trade tokens across multiple decentralized exchanges
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <BarChart3 className="w-4 h-4 mr-2" />
            Analytics
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Trading Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Swaps</p>
                <p className="text-2xl font-bold">{analytics.totalSwaps}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              All time transactions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Success Rate</p>
                <p className="text-2xl font-bold text-green-600">
                  {analytics.successRate.toFixed(1)}%
                </p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Transaction success rate
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Volume</p>
                <p className="text-2xl font-bold">${formatNumber(analytics.totalVolume * 1800, 0)}</p>
              </div>
              <DollarSign className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Trading volume (USD)
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Avg Slippage</p>
                <p className="text-2xl font-bold">{analytics.averageSlippage.toFixed(2)}%</p>
              </div>
              <TrendingDown className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Average slippage
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Trading Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="swap">Swap</TabsTrigger>
          <TabsTrigger value="quotes">Compare Quotes</TabsTrigger>
          <TabsTrigger value="history">Transaction History</TabsTrigger>
        </TabsList>

        <TabsContent value="swap" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            {/* Swap Interface */}
            <Card>
              <CardHeader>
                <CardTitle>Token Swap</CardTitle>
                <CardDescription>
                  Swap tokens at the best available rates
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Token In */}
                <div className="space-y-2">
                  <Label>From</Label>
                  <div className="flex items-center gap-2">
                    <div className="flex-1">
                      <Input
                        type="number"
                        placeholder="0.0"
                        value={amountIn}
                        onChange={(e) => setAmountIn(e.target.value)}
                      />
                    </div>
                    <Button variant="outline" className="min-w-[100px]">
                      {tokenIn ? (
                        <div className="flex items-center gap-2">
                          <span>{tokenIn.symbol}</span>
                        </div>
                      ) : (
                        'Select Token'
                      )}
                    </Button>
                  </div>
                  {tokenIn && (
                    <p className="text-xs text-muted-foreground">
                      Balance: 10.5 {tokenIn.symbol} • ${formatNumber((tokenIn.priceUSD || 0) * 10.5, 2)}
                    </p>
                  )}
                </div>

                {/* Swap Button */}
                <div className="flex justify-center">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleSwapTokens}
                    className="rounded-full"
                  >
                    <ArrowUpDown className="w-4 h-4" />
                  </Button>
                </div>

                {/* Token Out */}
                <div className="space-y-2">
                  <Label>To</Label>
                  <div className="flex items-center gap-2">
                    <div className="flex-1">
                      <Input
                        type="number"
                        placeholder="0.0"
                        value={selectedQuote ? formatNumber(selectedQuote.amountOut) : ''}
                        readOnly
                      />
                    </div>
                    <Button variant="outline" className="min-w-[100px]">
                      {tokenOut ? (
                        <div className="flex items-center gap-2">
                          <span>{tokenOut.symbol}</span>
                        </div>
                      ) : (
                        'Select Token'
                      )}
                    </Button>
                  </div>
                  {tokenOut && selectedQuote && (
                    <p className="text-xs text-muted-foreground">
                      ≈ ${formatNumber(parseFloat(selectedQuote.amountOut) * (tokenOut.priceUSD || 0), 2)}
                    </p>
                  )}
                </div>

                {/* Slippage Settings */}
                <div className="space-y-2">
                  <Label>Slippage Tolerance</Label>
                  <div className="flex items-center gap-2">
                    <Input
                      type="number"
                      step="0.1"
                      min="0.1"
                      max="50"
                      value={slippage}
                      onChange={(e) => setSlippage(parseFloat(e.target.value))}
                      className="w-20"
                    />
                    <span className="text-sm text-muted-foreground">%</span>
                    <div className="flex gap-1 ml-auto">
                      {[0.1, 0.5, 1.0].map((value) => (
                        <Button
                          key={value}
                          variant={slippage === value ? "default" : "outline"}
                          size="sm"
                          onClick={() => setSlippage(value)}
                        >
                          {value}%
                        </Button>
                      ))}
                    </div>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="space-y-2">
                  <Button
                    onClick={handleGetQuotes}
                    disabled={!tokenIn || !tokenOut || !amountIn || state.isGettingQuote}
                    className="w-full"
                  >
                    {state.isGettingQuote ? (
                      <>
                        <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                        Getting Quotes...
                      </>
                    ) : (
                      'Get Quotes'
                    )}
                  </Button>

                  {selectedQuote && (
                    <Button
                      onClick={handleExecuteSwap}
                      disabled={state.isSwapping}
                      className="w-full"
                      variant="default"
                    >
                      {state.isSwapping ? (
                        <>
                          <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                          Swapping...
                        </>
                      ) : (
                        <>
                          <Zap className="w-4 h-4 mr-2" />
                          Execute Swap
                        </>
                      )}
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Quote Details */}
            <Card>
              <CardHeader>
                <CardTitle>Best Quote</CardTitle>
                <CardDescription>
                  Optimal swap route and pricing details
                </CardDescription>
              </CardHeader>
              <CardContent>
                {selectedQuote ? (
                  <div className="space-y-4">
                    {/* DEX Info */}
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-2xl">{getProtocolLogo(selectedQuote.dexId)}</span>
                        <div>
                          <p className="font-medium">{selectedQuote.dexId}</p>
                          <p className="text-sm text-muted-foreground">
                            Confidence: {selectedQuote.confidence.toFixed(0)}%
                          </p>
                        </div>
                      </div>
                      <Badge variant="default">Best Rate</Badge>
                    </div>

                    {/* Price Details */}
                    <div className="space-y-3">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Exchange Rate</span>
                        <span>1 {tokenIn?.symbol} = {formatNumber(selectedQuote.executionPrice)} {tokenOut?.symbol}</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Price Impact</span>
                        <span className={cn(
                          selectedQuote.priceImpact > 5 ? 'text-red-600' : 
                          selectedQuote.priceImpact > 2 ? 'text-yellow-600' : 'text-green-600'
                        )}>
                          {formatPercentage(selectedQuote.priceImpact)}
                        </span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Minimum Received</span>
                        <span>{formatNumber(selectedQuote.minimumReceived)} {tokenOut?.symbol}</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Gas Fee</span>
                        <span>${formatNumber(selectedQuote.fees.gasFee)}</span>
                      </div>

                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Total Fees</span>
                        <span>${formatNumber(selectedQuote.fees.totalFeeUSD)}</span>
                      </div>
                    </div>

                    {/* Route Visualization */}
                    <div className="space-y-2">
                      <Label>Swap Route</Label>
                      <div className="flex items-center gap-2 p-3 bg-muted rounded-lg">
                        <span className="font-medium">{tokenIn?.symbol}</span>
                        <ArrowRight className="w-4 h-4 text-muted-foreground" />
                        <span className="font-medium">{tokenOut?.symbol}</span>
                        <Badge variant="outline" className="ml-auto">
                          Direct
                        </Badge>
                      </div>
                    </div>

                    {/* Warnings */}
                    {selectedQuote.priceImpact > 5 && (
                      <Alert variant="destructive">
                        <AlertTriangle className="h-4 w-4" />
                        <AlertDescription>
                          High price impact ({formatPercentage(selectedQuote.priceImpact)}). 
                          Consider reducing trade size.
                        </AlertDescription>
                      </Alert>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <TrendingUp className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                    <h3 className="text-lg font-medium mb-2">No Quote Selected</h3>
                    <p className="text-muted-foreground">
                      Get quotes to see pricing details and swap routes
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Available Quotes */}
          {quotes.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Available Quotes ({quotes.length})</CardTitle>
                <CardDescription>
                  Compare rates across different DEX protocols
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {quotes.map((quote, index) => (
                    <motion.div
                      key={quote.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className={cn(
                        "flex items-center justify-between p-4 border rounded-lg cursor-pointer transition-all",
                        selectedQuote?.id === quote.id ? "border-primary bg-primary/5" : "hover:bg-muted/50"
                      )}
                      onClick={() => setSelectedQuote(quote)}
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-xl">{getProtocolLogo(quote.dexId)}</span>
                        <div>
                          <p className="font-medium">{quote.dexId}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatNumber(quote.amountOut)} {tokenOut?.symbol}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="font-medium">
                          ${formatNumber(quote.fees.totalFeeUSD)} fees
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {formatPercentage(quote.priceImpact)} impact
                        </p>
                      </div>

                      <div className="flex items-center gap-2">
                        <Badge variant={index === 0 ? "default" : "outline"}>
                          {index === 0 ? 'Best' : `#${index + 1}`}
                        </Badge>
                        {selectedQuote?.id === quote.id && (
                          <CheckCircle className="w-4 h-4 text-primary" />
                        )}
                      </div>
                    </motion.div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="quotes" className="space-y-4">
          <Card>
            <CardContent className="p-12 text-center">
              <BarChart3 className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">Quote Comparison</h3>
              <p className="text-muted-foreground">
                Advanced quote comparison features coming soon
              </p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          {state.transactions.length > 0 ? (
            <div className="space-y-3">
              {state.transactions.map((transaction, index) => (
                <motion.div
                  key={transaction.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <Card>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex items-center gap-2">
                            {transaction.status === 'confirmed' ? (
                              <CheckCircle className="w-4 h-4 text-green-500" />
                            ) : transaction.status === 'failed' ? (
                              <AlertTriangle className="w-4 h-4 text-red-500" />
                            ) : (
                              <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
                            )}
                          </div>
                          
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <p className="font-medium">
                                {formatNumber(transaction.amountIn)} {transaction.tokenIn.symbol}
                              </p>
                              <ArrowRight className="w-3 h-3 text-muted-foreground" />
                              <p className="font-medium">
                                {formatNumber(transaction.actualAmountOut || transaction.amountOut)} {transaction.tokenOut.symbol}
                              </p>
                            </div>
                            
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <span>Status: {transaction.status}</span>
                              <span>Slippage: {transaction.slippage}%</span>
                              <span>{new Date(transaction.timestamp).toLocaleTimeString()}</span>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          {transaction.hash && (
                            <Button variant="ghost" size="sm">
                              <ExternalLink className="w-4 h-4" />
                            </Button>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Transaction History</h3>
                <p className="text-muted-foreground">
                  Your swap transactions will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
