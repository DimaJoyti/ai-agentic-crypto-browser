'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  ArrowUpDown, 
  Zap, 
  Settings, 
  Info,
  AlertTriangle,
  CheckCircle,
  RefreshCw,
  ExternalLink,
  TrendingUp,
  Clock,
  DollarSign
} from 'lucide-react'
import { useDeFiProtocols } from '@/hooks/useDeFiProtocols'
import { useGasOptimization } from '@/hooks/useGasOptimization'
import { type Address } from 'viem'
import { toast } from 'sonner'

interface Token {
  address: Address
  symbol: string
  name: string
  decimals: number
  logoURI?: string
  balance?: string
  price?: string
}

interface TokenSwapProps {
  chainId: number
  userAddress?: Address
}

export function TokenSwap({ chainId, userAddress }: TokenSwapProps) {
  const [fromToken, setFromToken] = useState<Token | null>(null)
  const [toToken, setToToken] = useState<Token | null>(null)
  const [fromAmount, setFromAmount] = useState('')
  const [toAmount, setToAmount] = useState('')
  const [slippage, setSlippage] = useState('0.5')
  const [deadline, setDeadline] = useState('20')
  const [isEstimating, setIsEstimating] = useState(false)
  const [swapEstimate, setSwapEstimate] = useState<any>(null)

  const { estimateSwap, executeSwap, dexProtocols } = useDeFiProtocols({
    chainId,
    userAddress
  })

  const { getGasEstimates, getRecommendedPriority } = useGasOptimization({
    chainId
  })

  // Mock token list
  const tokens: Token[] = [
    {
      address: '0xA0b86a33E6441b8435b662f0E2d0c2837b0b3c0' as Address,
      symbol: 'USDC',
      name: 'USD Coin',
      decimals: 6,
      balance: '1000.00',
      price: '1.00'
    },
    {
      address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' as Address,
      symbol: 'WETH',
      name: 'Wrapped Ether',
      decimals: 18,
      balance: '2.5',
      price: '2400.00'
    },
    {
      address: '0x6B175474E89094C44Da98b954EedeAC495271d0F' as Address,
      symbol: 'DAI',
      name: 'Dai Stablecoin',
      decimals: 18,
      balance: '500.00',
      price: '1.00'
    }
  ]

  const handleEstimateSwap = async () => {
    if (!fromToken || !toToken || !fromAmount) return

    setIsEstimating(true)
    try {
      const estimate = await estimateSwap(
        fromToken.address,
        toToken.address,
        fromAmount
      )
      setSwapEstimate(estimate)
      setToAmount(estimate.amountOut)
    } catch (error) {
      console.error('Failed to estimate swap:', error)
    } finally {
      setIsEstimating(false)
    }
  }

  const handleSwap = async () => {
    if (!fromToken || !toToken || !fromAmount || !toAmount || !userAddress) {
      toast.error('Please fill in all required fields')
      return
    }

    try {
      const minAmountOut = (parseFloat(toAmount) * (1 - parseFloat(slippage) / 100)).toString()
      
      const txHash = await executeSwap(
        fromToken.address,
        toToken.address,
        fromAmount,
        minAmountOut
      )

      toast.success('Swap executed successfully!', {
        description: `Swapped ${fromAmount} ${fromToken.symbol} for ${toAmount} ${toToken.symbol}`,
        action: {
          label: 'View Transaction',
          onClick: () => window.open(`https://etherscan.io/tx/${txHash}`, '_blank')
        }
      })

      // Reset form
      setFromAmount('')
      setToAmount('')
      setSwapEstimate(null)
    } catch (error) {
      console.error('Swap failed:', error)
    }
  }

  const handleFlipTokens = () => {
    const tempToken = fromToken
    setFromToken(toToken)
    setToToken(tempToken)
    
    const tempAmount = fromAmount
    setFromAmount(toAmount)
    setToAmount(tempAmount)
  }

  const calculateUsdValue = (amount: string, token: Token | null) => {
    if (!amount || !token?.price) return '$0.00'
    const value = parseFloat(amount) * parseFloat(token.price)
    return `$${value.toFixed(2)}`
  }

  useEffect(() => {
    if (fromToken && toToken && fromAmount) {
      const debounceTimer = setTimeout(() => {
        handleEstimateSwap()
      }, 500)
      return () => clearTimeout(debounceTimer)
    }
  }, [fromToken, toToken, fromAmount])

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <ArrowUpDown className="w-6 h-6" />
          Token Swap
        </h2>
        <p className="text-muted-foreground">
          Swap tokens across multiple DEX protocols with best price routing
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Swap Interface */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Swap Tokens</CardTitle>
              <CardDescription>
                Exchange tokens at the best available rates
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* From Token */}
              <div className="space-y-2">
                <Label>From</Label>
                <div className="border rounded-lg p-4 space-y-3">
                  <div className="flex items-center justify-between">
                    <Select
                      value={fromToken?.symbol || ''}
                      onValueChange={(symbol) => {
                        const token = tokens.find(t => t.symbol === symbol)
                        setFromToken(token || null)
                      }}
                    >
                      <SelectTrigger className="w-32">
                        <SelectValue placeholder="Token" />
                      </SelectTrigger>
                      <SelectContent>
                        {tokens.map((token) => (
                          <SelectItem key={token.symbol} value={token.symbol}>
                            <div className="flex items-center gap-2">
                              <span>{token.symbol}</span>
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    
                    <Input
                      type="number"
                      placeholder="0.0"
                      value={fromAmount}
                      onChange={(e) => setFromAmount(e.target.value)}
                      className="text-right text-lg font-medium border-0 bg-transparent"
                    />
                  </div>
                  
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span>{calculateUsdValue(fromAmount, fromToken)}</span>
                    {fromToken?.balance && (
                      <span>Balance: {fromToken.balance} {fromToken.symbol}</span>
                    )}
                  </div>
                </div>
              </div>

              {/* Flip Button */}
              <div className="flex justify-center">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleFlipTokens}
                  className="rounded-full w-10 h-10 p-0"
                >
                  <ArrowUpDown className="w-4 h-4" />
                </Button>
              </div>

              {/* To Token */}
              <div className="space-y-2">
                <Label>To</Label>
                <div className="border rounded-lg p-4 space-y-3">
                  <div className="flex items-center justify-between">
                    <Select
                      value={toToken?.symbol || ''}
                      onValueChange={(symbol) => {
                        const token = tokens.find(t => t.symbol === symbol)
                        setToToken(token || null)
                      }}
                    >
                      <SelectTrigger className="w-32">
                        <SelectValue placeholder="Token" />
                      </SelectTrigger>
                      <SelectContent>
                        {tokens.map((token) => (
                          <SelectItem key={token.symbol} value={token.symbol}>
                            <div className="flex items-center gap-2">
                              <span>{token.symbol}</span>
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    
                    <div className="text-right">
                      <p className="text-lg font-medium">
                        {isEstimating ? (
                          <RefreshCw className="w-4 h-4 animate-spin" />
                        ) : (
                          toAmount || '0.0'
                        )}
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span>{calculateUsdValue(toAmount, toToken)}</span>
                    {toToken?.balance && (
                      <span>Balance: {toToken.balance} {toToken.symbol}</span>
                    )}
                  </div>
                </div>
              </div>

              {/* Swap Details */}
              {swapEstimate && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="border rounded-lg p-4 space-y-3"
                >
                  <h4 className="font-medium">Swap Details</h4>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Price Impact</span>
                      <span className={parseFloat(swapEstimate.priceImpact) > 3 ? 'text-red-600' : 'text-green-600'}>
                        {swapEstimate.priceImpact}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Estimated Gas</span>
                      <span>{swapEstimate.gasEstimate}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Route</span>
                      <span className="text-right">
                        {fromToken?.symbol} → {toToken?.symbol}
                      </span>
                    </div>
                  </div>
                </motion.div>
              )}

              {/* Settings */}
              <div className="border rounded-lg p-4 space-y-4">
                <h4 className="font-medium flex items-center gap-2">
                  <Settings className="w-4 h-4" />
                  Settings
                </h4>
                
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label>Slippage Tolerance</Label>
                    <div className="flex gap-2">
                      {['0.1', '0.5', '1.0'].map((value) => (
                        <Button
                          key={value}
                          variant={slippage === value ? 'default' : 'outline'}
                          size="sm"
                          onClick={() => setSlippage(value)}
                        >
                          {value}%
                        </Button>
                      ))}
                      <Input
                        type="number"
                        value={slippage}
                        onChange={(e) => setSlippage(e.target.value)}
                        className="w-20"
                        placeholder="Custom"
                      />
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <Label>Transaction Deadline</Label>
                    <div className="flex items-center gap-2">
                      <Input
                        type="number"
                        value={deadline}
                        onChange={(e) => setDeadline(e.target.value)}
                        className="w-20"
                      />
                      <span className="text-sm text-muted-foreground">minutes</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Swap Button */}
              <Button
                onClick={handleSwap}
                disabled={!fromToken || !toToken || !fromAmount || !userAddress || isEstimating}
                className="w-full"
                size="lg"
              >
                {!userAddress ? (
                  'Connect Wallet'
                ) : !fromToken || !toToken ? (
                  'Select Tokens'
                ) : !fromAmount ? (
                  'Enter Amount'
                ) : isEstimating ? (
                  <>
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                    Estimating...
                  </>
                ) : (
                  <>
                    <Zap className="w-4 h-4 mr-2" />
                    Swap {fromToken.symbol} for {toToken.symbol}
                  </>
                )}
              </Button>

              {/* Warnings */}
              {swapEstimate && parseFloat(swapEstimate.priceImpact) > 3 && (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    High price impact ({swapEstimate.priceImpact}). You may lose a significant portion of your funds.
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Available DEXs */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Available DEXs</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {dexProtocols.map((protocol) => (
                  <div key={protocol.id} className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">🦄</span>
                      <span className="text-sm font-medium">{protocol.name}</span>
                    </div>
                    <Badge variant="outline" className="text-xs">
                      {protocol.tvl}
                    </Badge>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Market Info */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Market Info</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">24h Volume</span>
                  <span className="text-sm font-medium">$125M</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">24h Fees</span>
                  <span className="text-sm font-medium">$375K</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Active Traders</span>
                  <span className="text-sm font-medium">15.4K</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Tips */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Info className="w-4 h-4" />
                Tips
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3 text-sm">
                <div className="flex items-start gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500 mt-0.5" />
                  <span>Use 0.5% slippage for most trades</span>
                </div>
                <div className="flex items-start gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500 mt-0.5" />
                  <span>Check price impact before swapping</span>
                </div>
                <div className="flex items-start gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500 mt-0.5" />
                  <span>Consider gas costs for small trades</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
