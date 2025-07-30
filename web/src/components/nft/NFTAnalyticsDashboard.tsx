'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  TrendingUp, 
  TrendingDown, 
  Brain, 
  Target, 
  BarChart3,
  PieChart,
  Activity,
  Star,
  Shield,
  Zap,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  Eye,
  DollarSign
} from 'lucide-react'
import { type Address } from 'viem'
import { formatCurrency, formatNumber, formatPercentage } from '@/lib/utils'

interface NFTAnalyticsDashboardProps {
  contractAddress: Address
  tokenId: string
  chainId?: number
  className?: string
}

interface AnalyticsData {
  rarity: {
    score: number
    rank: number
    tier: string
    percentile: number
    traits: Array<{
      type: string
      value: string
      rarity: number
      impact: number
    }>
  }
  price: {
    current: number
    predicted: Array<{
      timeframe: string
      price: number
      confidence: number
    }>
    targets: Array<{
      price: number
      probability: number
      timeframe: string
    }>
    volatility: number
  }
  market: {
    sentiment: number
    trend: 'bullish' | 'bearish' | 'neutral'
    volume: number
    liquidity: number
    marketCap: number
  }
  investment: {
    rating: 'strong_buy' | 'buy' | 'hold' | 'sell' | 'strong_sell'
    score: number
    confidence: number
    risks: string[]
    opportunities: string[]
  }
  risk: {
    overall: number
    factors: Array<{
      factor: string
      impact: number
      probability: number
    }>
    var95: number
    maxDrawdown: number
  }
}

export function NFTAnalyticsDashboard({ 
  contractAddress, 
  tokenId, 
  chainId = 1, 
  className 
}: NFTAnalyticsDashboardProps) {
  const [activeTab, setActiveTab] = useState('overview')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null)

  // Mock analytics data
  const mockAnalytics: AnalyticsData = {
    rarity: {
      score: 344.5,
      rank: 1,
      tier: 'Legendary',
      percentile: 99.99,
      traits: [
        { type: 'Background', value: 'Gold', rarity: 0.5, impact: 0.3 },
        { type: 'Eyes', value: 'Laser Eyes', rarity: 0.1, impact: 0.5 },
        { type: 'Mouth', value: 'Grin', rarity: 2.5, impact: 0.1 },
        { type: 'Hat', value: 'Crown', rarity: 0.2, impact: 0.4 }
      ]
    },
    price: {
      current: 25.5,
      predicted: [
        { timeframe: '7d', price: 27.2, confidence: 0.78 },
        { timeframe: '30d', price: 31.5, confidence: 0.65 },
        { timeframe: '90d', price: 38.2, confidence: 0.52 }
      ],
      targets: [
        { price: 30.0, probability: 0.65, timeframe: '30d' },
        { price: 35.0, probability: 0.45, timeframe: '90d' },
        { price: 42.0, probability: 0.25, timeframe: '180d' }
      ],
      volatility: 0.45
    },
    market: {
      sentiment: 0.72,
      trend: 'bullish',
      volume: 1250.5,
      liquidity: 0.85,
      marketCap: 180000
    },
    investment: {
      rating: 'buy',
      score: 78,
      confidence: 0.82,
      risks: ['Market volatility', 'Liquidity constraints', 'Regulatory uncertainty'],
      opportunities: ['Rarity premium expansion', 'Utility development', 'Institutional adoption']
    },
    risk: {
      overall: 65,
      factors: [
        { factor: 'Market Risk', impact: 0.3, probability: 0.7 },
        { factor: 'Liquidity Risk', impact: 0.25, probability: 0.4 },
        { factor: 'Collection Risk', impact: 0.2, probability: 0.3 }
      ],
      var95: 0.15,
      maxDrawdown: 0.35
    }
  }

  useEffect(() => {
    loadAnalytics()
  }, [contractAddress, tokenId, chainId])

  const loadAnalytics = async () => {
    setIsLoading(true)
    setError(null)

    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000))
      setAnalytics(mockAnalytics)
    } catch (error) {
      setError((error as Error).message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleRefresh = () => {
    loadAnalytics()
  }

  const getRatingColor = (rating: string) => {
    switch (rating) {
      case 'strong_buy': return 'text-green-700 bg-green-100'
      case 'buy': return 'text-green-600 bg-green-50'
      case 'hold': return 'text-yellow-600 bg-yellow-50'
      case 'sell': return 'text-red-600 bg-red-50'
      case 'strong_sell': return 'text-red-700 bg-red-100'
      default: return 'text-gray-600 bg-gray-50'
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'bullish': return <TrendingUp className="h-4 w-4 text-green-600" />
      case 'bearish': return <TrendingDown className="h-4 w-4 text-red-600" />
      default: return <Activity className="h-4 w-4 text-gray-600" />
    }
  }

  const getRiskColor = (risk: number) => {
    if (risk < 30) return 'text-green-600'
    if (risk < 60) return 'text-yellow-600'
    if (risk < 80) return 'text-orange-600'
    return 'text-red-600'
  }

  if (error) {
    return (
      <Alert className="border-red-200 bg-red-50">
        <AlertTriangle className="h-4 w-4 text-red-600" />
        <AlertDescription className="text-red-800">
          Failed to load analytics: {error}
          <Button 
            variant="outline" 
            size="sm" 
            onClick={handleRefresh}
            className="ml-2"
          >
            Retry
          </Button>
        </AlertDescription>
      </Alert>
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
          <Brain className="h-8 w-8 animate-pulse mx-auto mb-2 text-blue-600" />
          <p className="text-sm text-muted-foreground">Analyzing NFT...</p>
        </div>
      </div>
    )
  }

  if (!analytics) {
    return (
      <div className="text-center p-8 text-muted-foreground">
        No analytics data available
      </div>
    )
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">NFT Analytics</h2>
          <p className="text-muted-foreground">
            {contractAddress.slice(0, 6)}...{contractAddress.slice(-4)} #{tokenId}
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleRefresh}
          disabled={isLoading}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="rarity">Rarity Analysis</TabsTrigger>
          <TabsTrigger value="price">Price Prediction</TabsTrigger>
          <TabsTrigger value="market">Market Analysis</TabsTrigger>
          <TabsTrigger value="investment">Investment</TabsTrigger>
          <TabsTrigger value="risk">Risk Assessment</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          {/* Key Metrics */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Rarity Score</CardTitle>
                <Star className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{analytics.rarity.score}</div>
                <p className="text-xs text-muted-foreground">
                  Rank #{analytics.rarity.rank} ({analytics.rarity.tier})
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Current Price</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatCurrency(analytics.price.current)} ETH</div>
                <p className="text-xs text-muted-foreground">
                  Volatility: {formatPercentage(analytics.price.volatility)}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Market Sentiment</CardTitle>
                {getTrendIcon(analytics.market.trend)}
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatPercentage(analytics.market.sentiment)}</div>
                <p className="text-xs text-muted-foreground capitalize">
                  {analytics.market.trend} trend
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Investment Rating</CardTitle>
                <Target className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{analytics.investment.score}/100</div>
                <Badge className={getRatingColor(analytics.investment.rating)} variant="secondary">
                  {analytics.investment.rating.replace('_', ' ').toUpperCase()}
                </Badge>
              </CardContent>
            </Card>
          </div>

          {/* Quick Insights */}
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Price Predictions</CardTitle>
                <CardDescription>AI-powered price forecasts</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {analytics.price.predicted.map(prediction => (
                    <div key={prediction.timeframe} className="flex items-center justify-between">
                      <span className="text-sm">{prediction.timeframe}</span>
                      <div className="text-right">
                        <div className="font-medium">{formatCurrency(prediction.price)} ETH</div>
                        <div className="text-xs text-muted-foreground">
                          {formatPercentage(prediction.confidence)} confidence
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Assessment</CardTitle>
                <CardDescription>Portfolio risk analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span>Overall Risk</span>
                    <div className="flex items-center gap-2">
                      <Progress value={analytics.risk.overall} className="w-20 h-2" />
                      <span className={`text-sm font-medium ${getRiskColor(analytics.risk.overall)}`}>
                        {analytics.risk.overall}/100
                      </span>
                    </div>
                  </div>
                  {analytics.risk.factors.slice(0, 3).map(factor => (
                    <div key={factor.factor} className="flex items-center justify-between text-sm">
                      <span>{factor.factor}</span>
                      <span className="font-medium">{formatPercentage(factor.impact)}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="rarity" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Rarity Overview</CardTitle>
                <CardDescription>Comprehensive rarity analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-4xl font-bold text-yellow-600">{analytics.rarity.score}</div>
                    <div className="text-lg font-medium">{analytics.rarity.tier}</div>
                    <div className="text-sm text-muted-foreground">
                      Rank #{analytics.rarity.rank} • Top {formatPercentage(100 - analytics.rarity.percentile)}
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Rarity Percentile</span>
                      <span className="font-medium">{formatPercentage(analytics.rarity.percentile)}</span>
                    </div>
                    <Progress value={analytics.rarity.percentile} className="h-2" />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Trait Analysis</CardTitle>
                <CardDescription>Individual trait rarity breakdown</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {analytics.rarity.traits.map(trait => (
                    <div key={`${trait.type}_${trait.value}`} className="space-y-1">
                      <div className="flex justify-between text-sm">
                        <span className="font-medium">{trait.type}</span>
                        <span className="text-muted-foreground">{trait.value}</span>
                      </div>
                      <div className="flex justify-between text-xs">
                        <span>Rarity: {formatPercentage(trait.rarity)}</span>
                        <span>Impact: {formatPercentage(trait.impact)}</span>
                      </div>
                      <Progress value={100 - (trait.rarity * 100)} className="h-1" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="price" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Price Predictions</CardTitle>
                <CardDescription>AI-powered price forecasts</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-3xl font-bold">{formatCurrency(analytics.price.current)} ETH</div>
                    <div className="text-sm text-muted-foreground">Current Price</div>
                  </div>
                  
                  <div className="space-y-3">
                    {analytics.price.predicted.map(prediction => (
                      <div key={prediction.timeframe} className="p-3 border rounded-lg">
                        <div className="flex justify-between items-center mb-2">
                          <span className="font-medium">{prediction.timeframe}</span>
                          <span className="text-lg font-bold">
                            {formatCurrency(prediction.price)} ETH
                          </span>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span>Confidence</span>
                          <span>{formatPercentage(prediction.confidence)}</span>
                        </div>
                        <Progress value={prediction.confidence * 100} className="h-1 mt-1" />
                      </div>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Price Targets</CardTitle>
                <CardDescription>Potential price levels</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {analytics.price.targets.map((target, index) => (
                    <div key={index} className="p-3 border rounded-lg">
                      <div className="flex justify-between items-center mb-2">
                        <span className="font-medium">{formatCurrency(target.price)} ETH</span>
                        <Badge variant="outline">{target.timeframe}</Badge>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>Probability</span>
                        <span>{formatPercentage(target.probability)}</span>
                      </div>
                      <Progress value={target.probability * 100} className="h-1 mt-1" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="market" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader>
                <CardTitle>Market Sentiment</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-center">
                  <div className="text-3xl font-bold text-green-600">
                    {formatPercentage(analytics.market.sentiment)}
                  </div>
                  <div className="flex items-center justify-center gap-2 mt-2">
                    {getTrendIcon(analytics.market.trend)}
                    <span className="capitalize">{analytics.market.trend}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Volume</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-center">
                  <div className="text-3xl font-bold">
                    {formatCurrency(analytics.market.volume)} ETH
                  </div>
                  <div className="text-sm text-muted-foreground mt-2">
                    24h Volume
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Liquidity</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-center">
                  <div className="text-3xl font-bold">
                    {formatPercentage(analytics.market.liquidity)}
                  </div>
                  <div className="text-sm text-muted-foreground mt-2">
                    Liquidity Score
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="investment" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Investment Rating</CardTitle>
                <CardDescription>AI-powered investment analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-4xl font-bold">{analytics.investment.score}/100</div>
                    <Badge className={getRatingColor(analytics.investment.rating)} variant="secondary">
                      {analytics.investment.rating.replace('_', ' ').toUpperCase()}
                    </Badge>
                    <div className="text-sm text-muted-foreground mt-2">
                      Confidence: {formatPercentage(analytics.investment.confidence)}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Key Factors</CardTitle>
                <CardDescription>Investment considerations</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium text-green-600 mb-2">Opportunities</h4>
                    <ul className="space-y-1">
                      {analytics.investment.opportunities.map((opportunity, index) => (
                        <li key={index} className="text-sm flex items-start gap-2">
                          <CheckCircle className="h-3 w-3 text-green-600 mt-0.5 flex-shrink-0" />
                          {opportunity}
                        </li>
                      ))}
                    </ul>
                  </div>
                  
                  <div>
                    <h4 className="font-medium text-red-600 mb-2">Risks</h4>
                    <ul className="space-y-1">
                      {analytics.investment.risks.map((risk, index) => (
                        <li key={index} className="text-sm flex items-start gap-2">
                          <AlertTriangle className="h-3 w-3 text-red-600 mt-0.5 flex-shrink-0" />
                          {risk}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="risk" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Risk Overview</CardTitle>
                <CardDescription>Comprehensive risk assessment</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className={`text-4xl font-bold ${getRiskColor(analytics.risk.overall)}`}>
                      {analytics.risk.overall}/100
                    </div>
                    <div className="text-sm text-muted-foreground">Overall Risk Score</div>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex justify-between text-sm">
                      <span>Value at Risk (95%)</span>
                      <span className="font-medium">{formatPercentage(analytics.risk.var95)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Max Drawdown</span>
                      <span className="font-medium">{formatPercentage(analytics.risk.maxDrawdown)}</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Risk Factors</CardTitle>
                <CardDescription>Individual risk components</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {analytics.risk.factors.map(factor => (
                    <div key={factor.factor} className="space-y-1">
                      <div className="flex justify-between text-sm">
                        <span className="font-medium">{factor.factor}</span>
                        <span>{formatPercentage(factor.impact)}</span>
                      </div>
                      <div className="flex justify-between text-xs text-muted-foreground">
                        <span>Probability: {formatPercentage(factor.probability)}</span>
                      </div>
                      <Progress value={factor.impact * 100} className="h-1" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
