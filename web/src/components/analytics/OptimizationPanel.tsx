'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  Zap,
  TrendingUp,
  Target,
  AlertTriangle,
  CheckCircle,
  DollarSign,
  Settings,
  Lightbulb,
  ArrowRight,
  Star
} from 'lucide-react'

interface OptimizationOpportunity {
  id: string
  title: string
  description: string
  category: 'performance' | 'cost' | 'risk' | 'efficiency'
  priority: 'high' | 'medium' | 'low'
  impact: number
  effort: number
  estimatedSavings: string
  timeToImplement: string
  status: 'pending' | 'in_progress' | 'completed' | 'dismissed'
}

interface OptimizationMetrics {
  overallScore: number
  totalOpportunities: number
  implementedOptimizations: number
  potentialSavings: string
  performanceGain: number
  riskReduction: number
  opportunities: OptimizationOpportunity[]
}

export const OptimizationPanel: React.FC = () => {
  const [metrics, setMetrics] = useState<OptimizationMetrics | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchOptimizationData = async () => {
      try {
        const response = await fetch('/api/analytics/performance/optimization')
        if (response.ok) {
          const data = await response.json()
          setMetrics(data)
        }
      } catch (error) {
        console.error('Error fetching optimization data:', error)
        // Set mock data for demo
        setMetrics({
          overallScore: 78.5,
          totalOpportunities: 12,
          implementedOptimizations: 8,
          potentialSavings: '$45,200',
          performanceGain: 23.4,
          riskReduction: 15.8,
          opportunities: [
            {
              id: '1',
              title: 'Optimize Database Queries',
              description: 'Implement query caching and indexing for frequently accessed data',
              category: 'performance',
              priority: 'high',
              impact: 85,
              effort: 60,
              estimatedSavings: '$12,000',
              timeToImplement: '2-3 weeks',
              status: 'pending'
            },
            {
              id: '2',
              title: 'Implement Connection Pooling',
              description: 'Reduce connection overhead by implementing database connection pooling',
              category: 'performance',
              priority: 'high',
              impact: 75,
              effort: 40,
              estimatedSavings: '$8,500',
              timeToImplement: '1 week',
              status: 'in_progress'
            },
            {
              id: '3',
              title: 'Optimize Memory Usage',
              description: 'Reduce memory footprint by optimizing data structures and garbage collection',
              category: 'cost',
              priority: 'medium',
              impact: 60,
              effort: 70,
              estimatedSavings: '$15,000',
              timeToImplement: '3-4 weeks',
              status: 'pending'
            },
            {
              id: '4',
              title: 'Implement Circuit Breakers',
              description: 'Add circuit breaker pattern to prevent cascade failures',
              category: 'risk',
              priority: 'high',
              impact: 90,
              effort: 50,
              estimatedSavings: '$25,000',
              timeToImplement: '2 weeks',
              status: 'pending'
            }
          ]
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchOptimizationData()
  }, [])



  const getPriorityBadgeVariant = (priority: string) => {
    switch (priority) {
      case 'high': return 'destructive'
      case 'medium': return 'secondary'
      case 'low': return 'default'
      default: return 'outline'
    }
  }

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'performance': return <Zap className="h-4 w-4" />
      case 'cost': return <DollarSign className="h-4 w-4" />
      case 'risk': return <AlertTriangle className="h-4 w-4" />
      case 'efficiency': return <Target className="h-4 w-4" />
      default: return <Settings className="h-4 w-4" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'text-green-600'
      case 'in_progress': return 'text-blue-600'
      case 'pending': return 'text-yellow-600'
      case 'dismissed': return 'text-gray-600'
      default: return 'text-gray-600'
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading optimization insights...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Optimization Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Optimization Score</CardTitle>
            <Star className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {metrics?.overallScore.toFixed(1)}
            </div>
            <Progress value={metrics?.overallScore || 0} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-1">
              {(metrics?.overallScore || 0) >= 80 ? 'Excellent' :
               (metrics?.overallScore || 0) >= 60 ? 'Good' : 'Needs Improvement'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Opportunities</CardTitle>
            <Lightbulb className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {metrics?.totalOpportunities}
            </div>
            <p className="text-xs text-muted-foreground">
              {metrics?.implementedOptimizations} implemented
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Potential Savings</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {metrics?.potentialSavings}
            </div>
            <p className="text-xs text-muted-foreground">
              Annual cost reduction
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Performance Gain</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              +{metrics?.performanceGain.toFixed(1)}%
            </div>
            <p className="text-xs text-muted-foreground">
              Expected improvement
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Optimization Tabs */}
      <Tabs defaultValue="opportunities" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="opportunities">Opportunities</TabsTrigger>
          <TabsTrigger value="implemented">Implemented</TabsTrigger>
          <TabsTrigger value="insights">Insights</TabsTrigger>
        </TabsList>

        <TabsContent value="opportunities" className="space-y-6">
          <div className="space-y-4">
            {metrics?.opportunities
              .filter(opp => opp.status === 'pending' || opp.status === 'in_progress')
              .sort((a, b) => {
                const priorityOrder = { high: 3, medium: 2, low: 1 }
                return priorityOrder[b.priority] - priorityOrder[a.priority]
              })
              .map((opportunity) => (
                <Card key={opportunity.id} className="border-l-4 border-l-blue-500">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="space-y-1">
                        <CardTitle className="text-lg flex items-center gap-2">
                          {getCategoryIcon(opportunity.category)}
                          {opportunity.title}
                        </CardTitle>
                        <CardDescription>{opportunity.description}</CardDescription>
                      </div>
                      <div className="flex flex-col gap-2">
                        <Badge variant={getPriorityBadgeVariant(opportunity.priority)}>
                          {opportunity.priority} priority
                        </Badge>
                        <Badge variant="outline" className={getStatusColor(opportunity.status)}>
                          {opportunity.status.replace('_', ' ')}
                        </Badge>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                      <div>
                        <p className="text-sm font-medium">Impact</p>
                        <div className="flex items-center gap-2">
                          <Progress value={opportunity.impact} className="flex-1" />
                          <span className="text-sm">{opportunity.impact}%</span>
                        </div>
                      </div>
                      <div>
                        <p className="text-sm font-medium">Effort</p>
                        <div className="flex items-center gap-2">
                          <Progress value={opportunity.effort} className="flex-1" />
                          <span className="text-sm">{opportunity.effort}%</span>
                        </div>
                      </div>
                      <div>
                        <p className="text-sm font-medium">Savings</p>
                        <p className="text-sm text-green-600 font-medium">{opportunity.estimatedSavings}</p>
                      </div>
                      <div>
                        <p className="text-sm font-medium">Timeline</p>
                        <p className="text-sm text-muted-foreground">{opportunity.timeToImplement}</p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button size="sm" className="flex items-center gap-1">
                        Implement
                        <ArrowRight className="h-3 w-3" />
                      </Button>
                      <Button size="sm" variant="outline">
                        Learn More
                      </Button>
                      <Button size="sm" variant="ghost">
                        Dismiss
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
          </div>
        </TabsContent>

        <TabsContent value="implemented" className="space-y-6">
          <div className="space-y-4">
            {metrics?.opportunities
              .filter(opp => opp.status === 'completed')
              .map((opportunity) => (
                <Card key={opportunity.id} className="border-l-4 border-l-green-500">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="space-y-1">
                        <CardTitle className="text-lg flex items-center gap-2">
                          <CheckCircle className="h-5 w-5 text-green-600" />
                          {opportunity.title}
                        </CardTitle>
                        <CardDescription>{opportunity.description}</CardDescription>
                      </div>
                      <Badge variant="default" className="bg-green-100 text-green-800">
                        Completed
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                      <div>
                        <p className="text-sm font-medium">Impact Achieved</p>
                        <p className="text-sm text-green-600 font-medium">{opportunity.impact}%</p>
                      </div>
                      <div>
                        <p className="text-sm font-medium">Savings Realized</p>
                        <p className="text-sm text-green-600 font-medium">{opportunity.estimatedSavings}</p>
                      </div>
                      <div>
                        <p className="text-sm font-medium">Implementation Time</p>
                        <p className="text-sm text-muted-foreground">{opportunity.timeToImplement}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
          </div>
        </TabsContent>

        <TabsContent value="insights" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Optimization Categories</CardTitle>
                <CardDescription>Distribution of optimization opportunities</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-3">
                  {['performance', 'cost', 'risk', 'efficiency'].map((category) => {
                    const count = metrics?.opportunities.filter(opp => opp.category === category).length || 0
                    const percentage = (count / (metrics?.totalOpportunities || 1)) * 100
                    return (
                      <div key={category} className="space-y-1">
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium capitalize flex items-center gap-2">
                            {getCategoryIcon(category)}
                            {category}
                          </span>
                          <span className="text-sm">{count} opportunities</span>
                        </div>
                        <Progress value={percentage} className="h-2" />
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Implementation Progress</CardTitle>
                <CardDescription>Current optimization implementation status</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Completion Rate</span>
                    <span className="text-sm">
                      {metrics?.implementedOptimizations}/{metrics?.totalOpportunities}
                    </span>
                  </div>
                  <Progress
                    value={((metrics?.implementedOptimizations || 0) / (metrics?.totalOpportunities || 1)) * 100}
                    className="h-3"
                  />
                  <div className="grid grid-cols-2 gap-4 text-center">
                    <div>
                      <p className="text-2xl font-bold text-green-600">
                        {metrics?.implementedOptimizations}
                      </p>
                      <p className="text-xs text-muted-foreground">Completed</p>
                    </div>
                    <div>
                      <p className="text-2xl font-bold text-blue-600">
                        {(metrics?.totalOpportunities || 0) - (metrics?.implementedOptimizations || 0)}
                      </p>
                      <p className="text-xs text-muted-foreground">Remaining</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <Alert>
            <Lightbulb className="h-4 w-4" />
            <AlertDescription>
              <strong>Pro Tip:</strong> Focus on high-impact, low-effort optimizations first to maximize your ROI.
              The current recommendations could improve performance by {metrics?.performanceGain.toFixed(1)}%
              and reduce costs by {metrics?.potentialSavings} annually.
            </AlertDescription>
          </Alert>
        </TabsContent>
      </Tabs>
    </div>
  )
}
