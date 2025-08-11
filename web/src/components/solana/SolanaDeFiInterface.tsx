'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  ArrowUpDown,
  Zap,
  TrendingUp,
  RefreshCw,
  Settings,
  Info,
  AlertTriangle,
  Loader2,
  ExternalLink
} from 'lucide-react'
import { useSolanaDeFi } from '@/hooks/useSolanaDeFi'
import { useSolanaSwap } from '@/hooks/useSolanaSwap'
import { useSolanaWallet } from '@/hooks/useSolanaWallet'
import { formatCurrency, formatNumber, formatPercentage, cn } from '@/lib/utils'
import { toast } from 'sonner'

interface SolanaDeFiInterfaceProps {
  className?: string
}

export function SolanaDeFiInterface({ className }: SolanaDeFiInterfaceProps) {
  const [activeTab, setActiveTab] = useState('swap')
  const [fromToken, setFromToken] = useState('SOL')
  const [toToken, setToToken] = useState('USDC')
  const [fromAmount, setFromAmount] = useState('')
  const [slippage, setSlippage] = useState(0.5)
  const [showAdvanced, setShowAdvanced] = useState(false)

  const { isConnected, balance, tokens } = useSolanaWallet()
  const { protocols, topYields } = useSolanaDeFi({
    autoRefresh: true,
    refreshInterval: 60000
  })

  const {
    quote,
    isLoading: swapLoading,
    error: swapError,
    getQuote,
    executeSwap,
    resetQuote
  } = useSolanaSwap()

  // Available tokens for swapping
  const availableTokens = [
    { symbol: 'SOL', name: 'Solana', mint: 'So11111111111111111111111111111111111111112' },
    { symbol: 'USDC', name: 'USD Coin', mint: 'EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v' },
    { symbol: 'USDT', name: 'Tether', mint: 'Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB' },
    { symbol: 'JUP', name: 'Jupiter', mint: 'JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN' },
    { symbol: 'RAY', name: 'Raydium', mint: '4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R' },
    { symbol: 'ORCA', name: 'Orca', mint: 'orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE' }
  ]

  const handleGetQuote = async () => {
    if (!fromAmount || !fromToken || !toToken) {
      toast.error('Please fill in all fields')
      return
    }

    const fromTokenInfo = availableTokens.find(t => t.symbol === fromToken)
    const toTokenInfo = availableTokens.find(t => t.symbol === toToken)

    if (!fromTokenInfo || !toTokenInfo) {
      toast.error('Invalid token selection')
      return
    }

    try {
      await getQuote({
        inputMint: fromTokenInfo.mint,
        outputMint: toTokenInfo.mint,
        amount: parseFloat(fromAmount),
        slippageBps: slippage * 100
      })
    } catch (error) {
      toast.error('Failed to get quote')
    }
  }

  const handleSwap = async () => {
    if (!quote) {
      toast.error('No quote available')
      return
    }

    try {
      const result = await executeSwap()
      if (result.success) {
        toast.success('Swap executed successfully!')
        setFromAmount('')
        resetQuote()
      } else {
        toast.error(result.error || 'Swap failed')
      }
    } catch (error) {
      toast.error('Failed to execute swap')
    }
  }

  const handleSwapTokens = () => {
    const temp = fromToken
    setFromToken(toToken)
    setToToken(temp)
    resetQuote()
  }

  const getTokenBalance = (symbol: string) => {
    if (symbol === 'SOL') return balance
    const token = tokens.find(t => t.symbol === symbol)
    return token?.balance || 0
  }

  const canSwap = isConnected && fromAmount && quote && !swapLoading

  return (
    <div className={cn('space-y-6', className)}>
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="swap">Token Swap</TabsTrigger>
          <TabsTrigger value="yield">Yield Farming</TabsTrigger>
          <TabsTrigger value="protocols">Protocols</TabsTrigger>
        </TabsList>

        {/* Token Swap Tab */}
        <TabsContent value="swap" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <ArrowUpDown className="h-5 w-5 mr-2" />
                Token Swap
              </CardTitle>
              <CardDescription>
                Swap tokens using Jupiter aggregator for best prices
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {!isConnected && (
                <Alert>
                  <Info className="h-4 w-4" />
                  <AlertDescription>
                    Connect your wallet to start swapping tokens
                  </AlertDescription>
                </Alert>
              )}

              {/* From Token */}
              <div className="space-y-2">
                <Label>From</Label>
                <div className="flex space-x-2">
                  <Select value={fromToken} onValueChange={setFromToken}>
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {availableTokens.map(token => (
                        <SelectItem key={token.symbol} value={token.symbol}>
                          {token.symbol}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Input
                    type="number"
                    placeholder="0.0"
                    value={fromAmount}
                    onChange={(e) => setFromAmount(e.target.value)}
                    className="flex-1"
                  />
                </div>
                <div className="flex justify-between text-sm text-muted-foreground">
                  <span>Balance: {formatNumber(getTokenBalance(fromToken), 4)} {fromToken}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setFromAmount(getTokenBalance(fromToken).toString())}
                  >
                    Max
                  </Button>
                </div>
              </div>

              {/* Swap Button */}
              <div className="flex justify-center">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleSwapTokens}
                  className="rounded-full"
                >
                  <ArrowUpDown className="h-4 w-4" />
                </Button>
              </div>

              {/* To Token */}
              <div className="space-y-2">
                <Label>To</Label>
                <div className="flex space-x-2">
                  <Select value={toToken} onValueChange={setToToken}>
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {availableTokens.map(token => (
                        <SelectItem key={token.symbol} value={token.symbol}>
                          {token.symbol}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Input
                    type="number"
                    placeholder="0.0"
                    value={quote ? formatNumber(quote.outputAmount, 6) : ''}
                    readOnly
                    className="flex-1 bg-muted"
                  />
                </div>
                <div className="text-sm text-muted-foreground">
                  Balance: {formatNumber(getTokenBalance(toToken), 4)} {toToken}
                </div>
              </div>

              {/* Quote Information */}
              {quote && (
                <Card className="bg-muted/50">
                  <CardContent className="pt-4 space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Rate:</span>
                      <span>1 {fromToken} = {formatNumber(quote.outputAmount / quote.inputAmount, 6)} {toToken}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Price Impact:</span>
                      <span className={cn(
                        quote.priceImpact > 5 ? 'text-red-600' : 
                        quote.priceImpact > 1 ? 'text-yellow-600' : 'text-green-600'
                      )}>
                        {formatPercentage(quote.priceImpact)}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Slippage:</span>
                      <span>{formatPercentage(slippage)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Fees:</span>
                      <span>{formatCurrency(quote.fees / 1e9)} SOL</span>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Advanced Settings */}
              <div className="space-y-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setShowAdvanced(!showAdvanced)}
                  className="w-full"
                >
                  <Settings className="h-4 w-4 mr-2" />
                  Advanced Settings
                </Button>
                
                {showAdvanced && (
                  <div className="space-y-2 p-3 border rounded-md">
                    <Label>Slippage Tolerance (%)</Label>
                    <div className="flex space-x-2">
                      {[0.1, 0.5, 1.0, 3.0].map(value => (
                        <Button
                          key={value}
                          variant={slippage === value ? "default" : "outline"}
                          size="sm"
                          onClick={() => setSlippage(value)}
                        >
                          {value}%
                        </Button>
                      ))}
                      <Input
                        type="number"
                        placeholder="Custom"
                        value={slippage}
                        onChange={(e) => setSlippage(parseFloat(e.target.value) || 0.5)}
                        className="w-20"
                        step="0.1"
                        min="0.1"
                        max="50"
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-2">
                <Button
                  onClick={handleGetQuote}
                  disabled={!fromAmount || !isConnected || swapLoading}
                  className="flex-1"
                  variant="outline"
                >
                  {swapLoading ? (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  ) : (
                    <RefreshCw className="h-4 w-4 mr-2" />
                  )}
                  Get Quote
                </Button>
                <Button
                  onClick={handleSwap}
                  disabled={!canSwap}
                  className="flex-1"
                >
                  {swapLoading ? (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  ) : (
                    <Zap className="h-4 w-4 mr-2" />
                  )}
                  Swap
                </Button>
              </div>

              {swapError && (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>{swapError}</AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Yield Farming Tab */}
        <TabsContent value="yield" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <TrendingUp className="h-5 w-5 mr-2" />
                Yield Farming Opportunities
              </CardTitle>
              <CardDescription>
                Discover high-yield farming opportunities across Solana DeFi
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topYields?.slice(0, 5).map((yield_, index) => (
                  <Card key={index} className="p-4">
                    <div className="flex items-center justify-between">
                      <div className="space-y-1">
                        <div className="flex items-center space-x-2">
                          <Badge variant="outline">{yield_.protocol}</Badge>
                          <span className="font-medium">{yield_.pool}</span>
                        </div>
                        <div className="text-sm text-muted-foreground">
                          Tokens: {yield_.tokens.join(', ')}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          Min Deposit: {formatNumber(yield_.minimumDeposit)} {yield_.tokens[0]}
                        </div>
                      </div>
                      <div className="text-right space-y-1">
                        <div className="text-2xl font-bold text-green-600">
                          {formatPercentage(yield_.apy)}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          APY
                        </div>
                        <Badge variant={
                          yield_.risk === 'low' ? 'default' :
                          yield_.risk === 'medium' ? 'secondary' : 'destructive'
                        }>
                          {yield_.risk} risk
                        </Badge>
                      </div>
                    </div>
                  </Card>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Protocols Tab */}
        <TabsContent value="protocols" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {protocols?.slice(0, 9).map((protocol, index) => (
              <Card key={index}>
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-lg">{protocol.name}</CardTitle>
                    <Badge variant="outline">{protocol.category}</Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>TVL:</span>
                      <span className="font-semibold">
                        {formatCurrency(Number(protocol.tvl))}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>24h Volume:</span>
                      <span>{formatCurrency(Number(protocol.volume24h))}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>24h Change:</span>
                      <span className={cn(
                        protocol.tvlChange24h > 0 ? 'text-green-600' : 'text-red-600'
                      )}>
                        {formatPercentage(protocol.tvlChange24h)}
                      </span>
                    </div>
                  </div>
                  <Button variant="outline" size="sm" className="w-full">
                    <ExternalLink className="h-4 w-4 mr-2" />
                    Visit Protocol
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
