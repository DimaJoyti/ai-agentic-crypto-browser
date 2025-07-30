'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  TrendingUp,
  TrendingDown,
  Target,
  BarChart3,
  PieChart,
  Activity,
  Award,
  Minus
} from 'lucide-react'

interface BenchmarkData {
  name: string
  value: number
  benchmark: number
  difference: number
  percentageDiff: number
  trend: 'up' | 'down' | 'neutral'
  category: 'return' | 'risk' | 'efficiency' | 'volume'
}

interface BenchmarkComparison {
  portfolioName: string
  benchmarkName: string
  timeframe: string
  overallRating: number
  metrics: BenchmarkData[]
  summary: {
    outperforming: number
    underperforming: number
    neutral: number
  }
}

export const BenchmarkComparison: React.FC = () => {
  const [comparison, setComparison] = useState<BenchmarkComparison | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchBenchmarkData = async () => {
      try {
        const response = await fetch('/api/analytics/performance/benchmark')
        if (response.ok) {
          const data = await response.json()
          setComparison(data)
        }
      } catch (error) {
        console.error('Error fetching benchmark data:', error)
        // Set mock data for demo
        setComparison({
          portfolioName: 'Crypto Trading Portfolio',
          benchmarkName: 'Market Index (BTC/ETH Weighted)',
          timeframe: 'Last 30 Days',
          overallRating: 78.5,
          metrics: [
            {
              name: 'Total Return',
              value: 15.8,
              benchmark: 12.3,
              difference: 3.5,
              percentageDiff: 28.5,
              trend: 'up',
              category: 'return'
            },
            {
              name: 'Sharpe Ratio',
              value: 1.85,
              benchmark: 1.42,
              difference: 0.43,
              percentageDiff: 30.3,
              trend: 'up',
              category: 'risk'
            },
            {
              name: 'Max Drawdown',
              value: -8.2,
              benchmark: -12.1,
              difference: 3.9,
              percentageDiff: 32.2,
              trend: 'up',
              category: 'risk'
            },
            {
              name: 'Volatility',
              value: 18.5,
              benchmark: 22.3,
              difference: -3.8,
              percentageDiff: -17.0,
              trend: 'up',
              category: 'risk'
            },
            {
              name: 'Win Rate',
              value: 68.5,
              benchmark: 58.2,
              difference: 10.3,
              percentageDiff: 17.7,
              trend: 'up',
              category: 'efficiency'
            },
            {
              name: 'Average Trade Size',
              value: 2450,
              benchmark: 1890,
              difference: 560,
              percentageDiff: 29.6,
              trend: 'up',
              category: 'volume'
            },
            {
              name: 'Sortino Ratio',
              value: 2.34,
              benchmark: 1.89,
              difference: 0.45,
              percentageDiff: 23.8,
              trend: 'up',
              category: 'risk'
            },
            {
              name: 'Information Ratio',
              value: 0.85,
              benchmark: 0.65,
              difference: 0.20,
              percentageDiff: 30.8,
              trend: 'up',
              category: 'efficiency'
            }
          ],
          summary: {
            outperforming: 7,
            underperforming: 1,
            neutral: 0
          }
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchBenchmarkData()
  }, [])

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="h-4 w-4 text-green-600" />
      case 'down': return <TrendingDown className="h-4 w-4 text-red-600" />
      default: return <Minus className="h-4 w-4 text-gray-600" />
    }
  }

  const getPerformanceColor = (percentageDiff: number) => {
    if (percentageDiff > 10) return 'text-green-600'
    if (percentageDiff > 0) return 'text-green-500'
    if (percentageDiff > -10) return 'text-red-500'
    return 'text-red-600'
  }

  const getPerformanceBadge = (percentageDiff: number) => {
    if (percentageDiff > 10) return { variant: 'default' as const, text: 'Excellent' }
    if (percentageDiff > 0) return { variant: 'secondary' as const, text: 'Good' }
    if (percentageDiff > -10) return { variant: 'outline' as const, text: 'Below' }
    return { variant: 'destructive' as const, text: 'Poor' }
  }

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'return': return <TrendingUp className="h-4 w-4" />
      case 'risk': return <Target className="h-4 w-4" />
      case 'efficiency': return <Activity className="h-4 w-4" />
      case 'volume': return <BarChart3 className="h-4 w-4" />
      default: return <PieChart className="h-4 w-4" />
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading benchmark comparison...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Benchmark Overview */}
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Award className="h-5 w-5" />
                Benchmark Comparison
              </CardTitle>
              <CardDescription>
                {comparison?.portfolioName} vs {comparison?.benchmarkName}
              </CardDescription>
            </div>
            <div className="text-right">
              <div className="text-2xl font-bold text-blue-600">
                {comparison?.overallRating.toFixed(1)}
              </div>
              <p className="text-xs text-muted-foreground">Overall Score</p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">
                {comparison?.summary.outperforming}
              </div>
              <p className="text-sm text-muted-foreground">Outperforming</p>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-red-600">
                {comparison?.summary.underperforming}
              </div>
              <p className="text-sm text-muted-foreground">Underperforming</p>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-600">
                {comparison?.summary.neutral}
              </div>
              <p className="text-sm text-muted-foreground">Neutral</p>
            </div>
          </div>
          <div className="mt-4">
            <Progress
              value={(comparison?.summary.outperforming || 0) / ((comparison?.summary.outperforming || 0) + (comparison?.summary.underperforming || 0) + (comparison?.summary.neutral || 0)) * 100}
              className="h-2"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Performance vs benchmark ({comparison?.timeframe})
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Detailed Metrics */}
      <Tabs defaultValue="all" className="w-full">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="all">All Metrics</TabsTrigger>
          <TabsTrigger value="return">Returns</TabsTrigger>
          <TabsTrigger value="risk">Risk</TabsTrigger>
          <TabsTrigger value="efficiency">Efficiency</TabsTrigger>
          <TabsTrigger value="volume">Volume</TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-4">
          <div className="grid gap-4">
            {comparison?.metrics.map((metric, index) => (
              <Card key={index}>
                <CardContent className="pt-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      {getCategoryIcon(metric.category)}
                      <div>
                        <h4 className="font-medium">{metric.name}</h4>
                        <p className="text-sm text-muted-foreground capitalize">
                          {metric.category} metric
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-4">
                      <div className="text-right">
                        <p className="text-sm text-muted-foreground">Your Portfolio</p>
                        <p className="font-medium">
                          {metric.name.includes('Rate') || metric.name.includes('Return') || metric.name.includes('Drawdown') || metric.name.includes('Volatility')
                            ? `${metric.value}%`
                            : metric.value.toLocaleString()}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-muted-foreground">Benchmark</p>
                        <p className="font-medium">
                          {metric.name.includes('Rate') || metric.name.includes('Return') || metric.name.includes('Drawdown') || metric.name.includes('Volatility')
                            ? `${metric.benchmark}%`
                            : metric.benchmark.toLocaleString()}
                        </p>
                      </div>
                      <div className="text-right">
                        <div className="flex items-center gap-1">
                          {getTrendIcon(metric.trend)}
                          <span className={`font-medium ${getPerformanceColor(metric.percentageDiff)}`}>
                            {metric.percentageDiff > 0 ? '+' : ''}{metric.percentageDiff.toFixed(1)}%
                          </span>
                        </div>
                        <Badge variant={getPerformanceBadge(metric.percentageDiff).variant} className="mt-1">
                          {getPerformanceBadge(metric.percentageDiff).text}
                        </Badge>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {['return', 'risk', 'efficiency', 'volume'].map((category) => (
          <TabsContent key={category} value={category} className="space-y-4">
            <div className="grid gap-4">
              {comparison?.metrics
                .filter(metric => metric.category === category)
                .map((metric, index) => (
                  <Card key={index}>
                    <CardContent className="pt-6">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          {getCategoryIcon(metric.category)}
                          <div>
                            <h4 className="font-medium">{metric.name}</h4>
                            <p className="text-sm text-muted-foreground">
                              Difference: {metric.difference > 0 ? '+' : ''}{metric.difference}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-4">
                          <div className="text-right">
                            <p className="text-sm text-muted-foreground">Portfolio</p>
                            <p className="font-medium">
                              {metric.name.includes('Rate') || metric.name.includes('Return') || metric.name.includes('Drawdown') || metric.name.includes('Volatility')
                                ? `${metric.value}%`
                                : metric.value.toLocaleString()}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm text-muted-foreground">Benchmark</p>
                            <p className="font-medium">
                              {metric.name.includes('Rate') || metric.name.includes('Return') || metric.name.includes('Drawdown') || metric.name.includes('Volatility')
                                ? `${metric.benchmark}%`
                                : metric.benchmark.toLocaleString()}
                            </p>
                          </div>
                          <div className="text-right">
                            <div className="flex items-center gap-1">
                              {getTrendIcon(metric.trend)}
                              <span className={`font-medium ${getPerformanceColor(metric.percentageDiff)}`}>
                                {metric.percentageDiff > 0 ? '+' : ''}{metric.percentageDiff.toFixed(1)}%
                              </span>
                            </div>
                            <Badge variant={getPerformanceBadge(metric.percentageDiff).variant} className="mt-1">
                              {getPerformanceBadge(metric.percentageDiff).text}
                            </Badge>
                          </div>
                        </div>
                      </div>
                      <div className="mt-4">
                        <div className="flex items-center justify-between text-xs text-muted-foreground mb-1">
                          <span>Benchmark</span>
                          <span>Your Performance</span>
                        </div>
                        <div className="relative">
                          <Progress value={50} className="h-2" />
                          <div
                            className="absolute top-0 h-2 bg-blue-600 rounded-full"
                            style={{
                              width: `${Math.min(Math.max((metric.value / (metric.benchmark * 2)) * 100, 0), 100)}%`
                            }}
                          />
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
            </div>
          </TabsContent>
        ))}
      </Tabs>
    </div>
  )
}
