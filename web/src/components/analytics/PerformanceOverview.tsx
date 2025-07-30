'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import {
  TrendingUp,
  TrendingDown,
  Activity,
  Zap,
  Target,
  AlertTriangle,
  CheckCircle,
  Clock,
  DollarSign
} from 'lucide-react'

interface PerformanceOverviewProps {
  metrics: any
  keyMetrics: any
  systemHealth: any
  optimization: any
}

export const PerformanceOverview: React.FC<PerformanceOverviewProps> = ({
  metrics,
  keyMetrics,
  systemHealth,
  optimization
}) => {
  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-600'
    if (score >= 70) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getScoreBadgeVariant = (score: number) => {
    if (score >= 90) return 'default'
    if (score >= 70) return 'secondary'
    return 'destructive'
  }

  const getTrendIcon = (value: string) => {
    if (value?.includes('+') || parseFloat(value) > 0) {
      return <TrendingUp className="h-4 w-4 text-green-600" />
    } else if (value?.includes('-') || parseFloat(value) < 0) {
      return <TrendingDown className="h-4 w-4 text-red-600" />
    }
    return <Activity className="h-4 w-4 text-blue-600" />
  }

  return (
    <div className="space-y-6">
      {/* Performance Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Overall Performance</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(metrics?.overallScore || 0)}`}>
              {(metrics?.overallScore || 0).toFixed(1)}
            </div>
            <Progress value={metrics?.overallScore || 0} className="mt-2" />
            <Badge variant={getScoreBadgeVariant(metrics?.overallScore || 0)} className="mt-2">
              {(metrics?.overallScore || 0) >= 90 ? 'Excellent' :
               (metrics?.overallScore || 0) >= 70 ? 'Good' : 'Needs Improvement'}
            </Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Return</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {keyMetrics?.totalReturn || '0%'}
            </div>
            <div className="flex items-center gap-1 mt-1">
              {getTrendIcon(keyMetrics?.totalReturn)}
              <p className="text-xs text-muted-foreground">
                vs last period
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">System Health</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {systemHealth?.cpuUsage < 80 ? (
                <CheckCircle className="h-8 w-8 text-green-600" />
              ) : (
                <AlertTriangle className="h-8 w-8 text-yellow-600" />
              )}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              CPU: {(systemHealth?.cpuUsage || 0).toFixed(1)}%
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Optimization Score</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${getScoreColor(optimization?.score || 0)}`}>
              {(optimization?.score || 0).toFixed(1)}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {optimization?.opportunitiesCount || 0} opportunities
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Detailed Performance Metrics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Key Performance Indicators</CardTitle>
            <CardDescription>Critical metrics for system performance</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Success Rate</span>
              <div className="flex items-center gap-2">
                <span className="text-sm">{(keyMetrics?.successRate || 0).toFixed(1)}%</span>
                <Progress value={keyMetrics?.successRate || 0} className="w-20" />
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Average Latency</span>
              <div className="flex items-center gap-2">
                <span className="text-sm">{keyMetrics?.avgLatencyMs || 0}ms</span>
                <Clock className="h-4 w-4 text-muted-foreground" />
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Sharpe Ratio</span>
              <span className="text-sm font-medium">{(keyMetrics?.sharpeRatio || 0).toFixed(2)}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Max Drawdown</span>
              <span className="text-sm text-red-600">{keyMetrics?.maxDrawdown || '0%'}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Resource Usage</CardTitle>
            <CardDescription>Current system resource utilization</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">CPU Usage</span>
                <span className="text-sm">{(systemHealth?.cpuUsage || 0).toFixed(1)}%</span>
              </div>
              <Progress value={systemHealth?.cpuUsage || 0} className="h-2" />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Memory Usage</span>
                <span className="text-sm">{(systemHealth?.memoryUsage || 0).toFixed(1)}%</span>
              </div>
              <Progress value={systemHealth?.memoryUsage || 0} className="h-2" />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Disk Usage</span>
                <span className="text-sm">{(systemHealth?.diskUsage || 0).toFixed(1)}%</span>
              </div>
              <Progress value={systemHealth?.diskUsage || 0} className="h-2" />
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Uptime</span>
              <span className="text-sm">{systemHealth?.uptime || 'N/A'}</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
