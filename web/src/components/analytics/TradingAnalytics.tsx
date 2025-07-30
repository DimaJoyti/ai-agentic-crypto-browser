'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Target,
  Activity,
  BarChart3,
  DollarSign,
  Percent,
  Clock,
  Zap
} from 'lucide-react'

interface TradingMetrics {
  totalTrades: number
  successfulTrades: number
  successRate: number
  totalPnL: string
  sharpeRatio: number
  sortinoRatio: number
  maxDrawdown: string
  winRate: number
  profitFactor: number
  averageWin: string
  averageLoss: string
  largestWin: string
  largestLoss: string
  volatility: number
  beta: number
  alpha: number
}

export const TradingAnalytics: React.FC = () => {
  const [metrics, setMetrics] = useState<TradingMetrics | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchTradingMetrics = async () => {
      try {
        const response = await fetch('/api/analytics/performance/trading')
        if (response.ok) {
          const data = await response.json()
          setMetrics(data)
        }
      } catch (error) {
        console.error('Error fetching trading metrics:', error)
        // Set mock data for demo
        setMetrics({
          totalTrades: 1247,
          successfulTrades: 1156,
          successRate: 92.7,
          totalPnL: '+$125,430',
          sharpeRatio: 1.85,
          sortinoRatio: 2.34,
          maxDrawdown: '-8.2%',
          winRate: 68.5,
          profitFactor: 1.42,
          averageWin: '+$245',
          averageLoss: '-$156',
          largestWin: '+$2,450',
          largestLoss: '-$890',
          volatility: 15.3,
          beta: 0.85,
          alpha: 2.1
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchTradingMetrics()
  }, [])

  const getScoreColor = (score: number, threshold: number = 70) => {
    if (score >= threshold) return 'text-green-600'
    if (score >= threshold * 0.7) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getScoreBadgeVariant = (score: number, threshold: number = 70) => {
    if (score >= threshold) return 'default'
    if (score >= threshold * 0.7) return 'secondary'
    return 'destructive'
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading trading analytics...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Trading Performance Summary */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Trades</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{metrics?.totalTrades.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              {metrics?.successfulTrades} successful
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.successRate || 0, 90)}`}>
              {metrics?.successRate.toFixed(1)}%
            </div>
            <Progress value={metrics?.successRate || 0} className="mt-2" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total P&L</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${metrics?.totalPnL.includes('+') ? 'text-green-600' : 'text-red-600'}`}>
              {metrics?.totalPnL}
            </div>
            <p className="text-xs text-muted-foreground">
              Realized gains/losses
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Sharpe Ratio</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.sharpeRatio || 0, 1.5)}`}>
              {metrics?.sharpeRatio.toFixed(2)}
            </div>
            <Badge variant={getScoreBadgeVariant(metrics?.sharpeRatio || 0, 1.5)} className="mt-2">
              {(metrics?.sharpeRatio || 0) >= 1.5 ? 'Excellent' :
               (metrics?.sharpeRatio || 0) >= 1.0 ? 'Good' : 'Poor'}
            </Badge>
          </CardContent>
        </Card>
      </div>

      {/* Detailed Analytics Tabs */}
      <Tabs defaultValue="performance" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="risk">Risk Metrics</TabsTrigger>
          <TabsTrigger value="trades">Trade Analysis</TabsTrigger>
          <TabsTrigger value="ratios">Financial Ratios</TabsTrigger>
        </TabsList>

        <TabsContent value="performance" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Win/Loss Analysis</CardTitle>
                <CardDescription>Trading success breakdown</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Win Rate</span>
                  <div className="flex items-center gap-2">
                    <span className="text-sm">{metrics?.winRate.toFixed(1)}%</span>
                    <Progress value={metrics?.winRate || 0} className="w-20" />
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Profit Factor</span>
                  <span className={`text-sm font-medium ${getScoreColor(metrics?.profitFactor || 0, 1.2)}`}>
                    {metrics?.profitFactor.toFixed(2)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Average Win</span>
                  <span className="text-sm text-green-600">{metrics?.averageWin}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Average Loss</span>
                  <span className="text-sm text-red-600">{metrics?.averageLoss}</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Extreme Values</CardTitle>
                <CardDescription>Best and worst trades</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Largest Win</span>
                  <span className="text-sm text-green-600 font-medium">{metrics?.largestWin}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Largest Loss</span>
                  <span className="text-sm text-red-600 font-medium">{metrics?.largestLoss}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Max Drawdown</span>
                  <span className="text-sm text-red-600 font-medium">{metrics?.maxDrawdown}</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="risk" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Risk-Adjusted Returns</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Sharpe Ratio</span>
                  <span className={`text-sm font-medium ${getScoreColor(metrics?.sharpeRatio || 0, 1.5)}`}>
                    {metrics?.sharpeRatio.toFixed(2)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Sortino Ratio</span>
                  <span className={`text-sm font-medium ${getScoreColor(metrics?.sortinoRatio || 0, 2.0)}`}>
                    {metrics?.sortinoRatio.toFixed(2)}
                  </span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Market Exposure</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Beta</span>
                  <span className="text-sm font-medium">{metrics?.beta.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Alpha</span>
                  <span className={`text-sm font-medium ${(metrics?.alpha || 0) > 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {(metrics?.alpha || 0).toFixed(1)}%
                  </span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Volatility</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Annualized</span>
                  <span className="text-sm font-medium">{metrics?.volatility.toFixed(1)}%</span>
                </div>
                <Progress value={metrics?.volatility || 0} className="mt-2" />
                <p className="text-xs text-muted-foreground">
                  {(metrics?.volatility || 0) < 20 ? 'Low' : (metrics?.volatility || 0) < 30 ? 'Moderate' : 'High'} volatility
                </p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="trades" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Trade Distribution</CardTitle>
              <CardDescription>Analysis of trading patterns</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-8">
                <div className="space-y-4">
                  <h4 className="text-sm font-medium">Successful Trades</h4>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Count</span>
                      <span className="text-sm font-medium text-green-600">{metrics?.successfulTrades}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Percentage</span>
                      <span className="text-sm font-medium text-green-600">{metrics?.successRate.toFixed(1)}%</span>
                    </div>
                  </div>
                </div>
                <div className="space-y-4">
                  <h4 className="text-sm font-medium">Failed Trades</h4>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Count</span>
                      <span className="text-sm font-medium text-red-600">
                        {(metrics?.totalTrades || 0) - (metrics?.successfulTrades || 0)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Percentage</span>
                      <span className="text-sm font-medium text-red-600">
                        {(100 - (metrics?.successRate || 0)).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="ratios" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Percent className="h-4 w-4" />
                  Efficiency Ratios
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-sm">Profit Factor</span>
                  <span className="font-medium">{metrics?.profitFactor.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Win Rate</span>
                  <span className="font-medium">{metrics?.winRate.toFixed(1)}%</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="h-4 w-4" />
                  Risk Ratios
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-sm">Sharpe</span>
                  <span className="font-medium">{metrics?.sharpeRatio.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Sortino</span>
                  <span className="font-medium">{metrics?.sortinoRatio.toFixed(2)}</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="h-4 w-4" />
                  Market Ratios
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-sm">Beta</span>
                  <span className="font-medium">{metrics?.beta.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm">Alpha</span>
                  <span className="font-medium">{metrics?.alpha.toFixed(1)}%</span>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
