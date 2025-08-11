'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { 
  TrendingUp, 
  TrendingDown, 
  Wallet, 
  Calculator,
  AlertTriangle,
  Info,
  Zap,
  Target
} from 'lucide-react'
import { useAccount, useBalance } from 'wagmi'
import { cn } from '@/lib/utils'

interface TradingPair {
  symbol: string
  baseAsset: string
  quoteAsset: string
  price: number
  change24h: number
  volume24h: number
  high24h: number
  low24h: number
}

interface OrderFormProps {
  pair: TradingPair
  orderType: 'market' | 'limit' | 'stop'
  side: 'buy' | 'sell'
  onOrderTypeChange: (type: 'market' | 'limit' | 'stop') => void
  onSideChange: (side: 'buy' | 'sell') => void
}

export function OrderForm({ 
  pair, 
  orderType, 
  side, 
  onOrderTypeChange, 
  onSideChange 
}: OrderFormProps) {
  const [amount, setAmount] = useState('')
  const [price, setPrice] = useState('')
  const [stopPrice, setStopPrice] = useState('')
  const [total, setTotal] = useState('')
  const [percentage, setPercentage] = useState([0])
  const [isAdvanced, setIsAdvanced] = useState(false)
  const [timeInForce, setTimeInForce] = useState('GTC')
  const [reduceOnly, setReduceOnly] = useState(false)
  const [postOnly, setPostOnly] = useState(false)

  const { address, isConnected } = useAccount()
  const { data: baseBalance } = useBalance({
    address,
    token: pair.baseAsset === 'ETH' ? undefined : '0x...' // Mock token address
  })
  const { data: quoteBalance } = useBalance({
    address,
    token: pair.quoteAsset === 'USDT' ? '0x...' : undefined // Mock USDT address
  })

  // Calculate total when amount or price changes
  useEffect(() => {
    if (amount && (orderType === 'market' ? pair.price : price)) {
      const calculatedTotal = parseFloat(amount) * (orderType === 'market' ? pair.price : parseFloat(price))
      setTotal(calculatedTotal.toFixed(6))
    }
  }, [amount, price, pair.price, orderType])

  // Update price when pair changes for limit orders
  useEffect(() => {
    if (orderType === 'limit' && !price) {
      setPrice(pair.price.toFixed(6))
    }
  }, [pair.price, orderType, price])

  const handlePercentageChange = (value: number[]) => {
    setPercentage(value)
    const percent = value[0]
    
    if (side === 'buy' && quoteBalance) {
      const availableBalance = parseFloat(quoteBalance.formatted)
      const orderPrice = orderType === 'market' ? pair.price : parseFloat(price) || pair.price
      const maxAmount = (availableBalance * percent / 100) / orderPrice
      setAmount(maxAmount.toFixed(6))
    } else if (side === 'sell' && baseBalance) {
      const availableBalance = parseFloat(baseBalance.formatted)
      const orderAmount = availableBalance * percent / 100
      setAmount(orderAmount.toFixed(6))
    }
  }

  const getAvailableBalance = () => {
    if (side === 'buy') {
      return quoteBalance ? parseFloat(quoteBalance.formatted) : 0
    } else {
      return baseBalance ? parseFloat(baseBalance.formatted) : 0
    }
  }

  const getBalanceSymbol = () => {
    return side === 'buy' ? pair.quoteAsset : pair.baseAsset
  }

  const calculateFees = () => {
    const totalValue = parseFloat(total) || 0
    const feeRate = 0.001 // 0.1% fee
    return totalValue * feeRate
  }

  const handleSubmitOrder = () => {
    if (!isConnected) {
      // Handle wallet connection
      return
    }

    const orderData = {
      pair: pair.symbol,
      side,
      type: orderType,
      amount: parseFloat(amount),
      price: orderType === 'market' ? undefined : parseFloat(price),
      stopPrice: orderType === 'stop' ? parseFloat(stopPrice) : undefined,
      timeInForce,
      reduceOnly,
      postOnly
    }

    console.log('Submitting order:', orderData)
    // Here you would call your trading API
  }

  const isFormValid = () => {
    if (!amount || parseFloat(amount) <= 0) return false
    if (orderType === 'limit' && (!price || parseFloat(price) <= 0)) return false
    if (orderType === 'stop' && (!stopPrice || parseFloat(stopPrice) <= 0)) return false
    return true
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-medium">Place Order</CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsAdvanced(!isAdvanced)}
          >
            <Calculator className="w-4 h-4 mr-1" />
            {isAdvanced ? 'Simple' : 'Advanced'}
          </Button>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Buy/Sell Tabs */}
        <Tabs value={side} onValueChange={(value) => onSideChange(value as 'buy' | 'sell')}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="buy" className="text-green-600 data-[state=active]:bg-green-100 data-[state=active]:text-green-700">
              <TrendingUp className="w-4 h-4 mr-1" />
              Buy
            </TabsTrigger>
            <TabsTrigger value="sell" className="text-red-600 data-[state=active]:bg-red-100 data-[state=active]:text-red-700">
              <TrendingDown className="w-4 h-4 mr-1" />
              Sell
            </TabsTrigger>
          </TabsList>
        </Tabs>

        {/* Order Type */}
        <div className="space-y-2">
          <Label className="text-xs font-medium">Order Type</Label>
          <Select value={orderType} onValueChange={(value) => onOrderTypeChange(value as any)}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="market">
                <div className="flex items-center gap-2">
                  <Zap className="w-4 h-4" />
                  Market
                </div>
              </SelectItem>
              <SelectItem value="limit">
                <div className="flex items-center gap-2">
                  <Target className="w-4 h-4" />
                  Limit
                </div>
              </SelectItem>
              <SelectItem value="stop">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="w-4 h-4" />
                  Stop
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Price Input (for limit and stop orders) */}
        {orderType !== 'market' && (
          <div className="space-y-2">
            <Label className="text-xs font-medium">
              {orderType === 'stop' ? 'Stop Price' : 'Price'} ({pair.quoteAsset})
            </Label>
            <div className="relative">
              <Input
                type="number"
                placeholder="0.00"
                value={orderType === 'stop' ? stopPrice : price}
                onChange={(e) => orderType === 'stop' ? setStopPrice(e.target.value) : setPrice(e.target.value)}
                className="pr-16"
              />
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-muted-foreground">
                {pair.quoteAsset}
              </div>
            </div>
            {orderType === 'limit' && (
              <div className="text-xs text-muted-foreground">
                Market Price: {pair.price.toFixed(6)} {pair.quoteAsset}
              </div>
            )}
          </div>
        )}

        {/* Amount Input */}
        <div className="space-y-2">
          <Label className="text-xs font-medium">Amount ({pair.baseAsset})</Label>
          <div className="relative">
            <Input
              type="number"
              placeholder="0.00"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="pr-16"
            />
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-muted-foreground">
              {pair.baseAsset}
            </div>
          </div>
        </div>

        {/* Percentage Slider */}
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <Label className="text-xs font-medium">Amount</Label>
            <span className="text-xs text-muted-foreground">{percentage[0]}%</span>
          </div>
          <Slider
            value={percentage}
            onValueChange={handlePercentageChange}
            max={100}
            step={25}
            className="w-full"
          />
          <div className="grid grid-cols-4 gap-2">
            {[25, 50, 75, 100].map((percent) => (
              <Button
                key={percent}
                variant="outline"
                size="sm"
                onClick={() => handlePercentageChange([percent])}
                className="text-xs"
              >
                {percent}%
              </Button>
            ))}
          </div>
        </div>

        {/* Total */}
        <div className="space-y-2">
          <Label className="text-xs font-medium">Total ({pair.quoteAsset})</Label>
          <div className="relative">
            <Input
              type="number"
              placeholder="0.00"
              value={total}
              onChange={(e) => setTotal(e.target.value)}
              className="pr-16"
            />
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-muted-foreground">
              {pair.quoteAsset}
            </div>
          </div>
        </div>

        {/* Advanced Options */}
        {isAdvanced && (
          <div className="space-y-4 pt-4 border-t">
            <div className="space-y-2">
              <Label className="text-xs font-medium">Time in Force</Label>
              <Select value={timeInForce} onValueChange={setTimeInForce}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="GTC">Good Till Canceled</SelectItem>
                  <SelectItem value="IOC">Immediate or Cancel</SelectItem>
                  <SelectItem value="FOK">Fill or Kill</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center justify-between">
              <Label className="text-xs font-medium">Reduce Only</Label>
              <Switch checked={reduceOnly} onCheckedChange={setReduceOnly} />
            </div>

            <div className="flex items-center justify-between">
              <Label className="text-xs font-medium">Post Only</Label>
              <Switch checked={postOnly} onCheckedChange={setPostOnly} />
            </div>
          </div>
        )}

        <Separator />

        {/* Balance and Fees */}
        <div className="space-y-2 text-xs">
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Available:</span>
            <span className="flex items-center gap-1">
              <Wallet className="w-3 h-3" />
              {getAvailableBalance().toFixed(6)} {getBalanceSymbol()}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Est. Fee:</span>
            <span>{calculateFees().toFixed(6)} {pair.quoteAsset}</span>
          </div>
        </div>

        {/* Submit Button */}
        <Button
          className={cn(
            "w-full",
            side === 'buy' ? "bg-green-600 hover:bg-green-700" : "bg-red-600 hover:bg-red-700"
          )}
          disabled={!isConnected || !isFormValid()}
          onClick={handleSubmitOrder}
        >
          {!isConnected ? (
            'Connect Wallet'
          ) : (
            `${side === 'buy' ? 'Buy' : 'Sell'} ${pair.baseAsset}`
          )}
        </Button>

        {/* Order Summary */}
        {isFormValid() && (
          <div className="p-3 bg-muted/50 rounded-lg space-y-1 text-xs">
            <div className="flex justify-between">
              <span>Order Type:</span>
              <span className="capitalize">{orderType}</span>
            </div>
            <div className="flex justify-between">
              <span>Side:</span>
              <span className={cn(
                "capitalize font-medium",
                side === 'buy' ? "text-green-600" : "text-red-600"
              )}>
                {side}
              </span>
            </div>
            <div className="flex justify-between">
              <span>Amount:</span>
              <span>{amount} {pair.baseAsset}</span>
            </div>
            {orderType !== 'market' && (
              <div className="flex justify-between">
                <span>Price:</span>
                <span>{price} {pair.quoteAsset}</span>
              </div>
            )}
            <div className="flex justify-between font-medium">
              <span>Total:</span>
              <span>{total} {pair.quoteAsset}</span>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
