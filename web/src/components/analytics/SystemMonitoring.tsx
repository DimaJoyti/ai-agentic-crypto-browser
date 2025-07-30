'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  Cpu,
  MemoryStick,
  HardDrive,
  Network,
  Clock,
  AlertTriangle,
  CheckCircle,
  Database,
  Wifi
} from 'lucide-react'

interface SystemMetrics {
  cpuUsage: number
  memoryUsage: number
  diskUsage: number
  networkLatency: number
  databaseLatency: number
  apiLatency: number
  throughput: number
  errorRate: number
  uptime: string
  activeConnections: number
  queueDepth: number
  cacheHitRate: number
  gcPauseTime: number
  goroutineCount: number
  heapSize: number
  allocRate: number
}

export const SystemMonitoring: React.FC = () => {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchSystemMetrics = async () => {
      try {
        const response = await fetch('/api/analytics/performance/system')
        if (response.ok) {
          const data = await response.json()
          setMetrics(data)
        }
      } catch (error) {
        console.error('Error fetching system metrics:', error)
        // Set mock data for demo
        setMetrics({
          cpuUsage: 45.2,
          memoryUsage: 68.7,
          diskUsage: 34.1,
          networkLatency: 12,
          databaseLatency: 8,
          apiLatency: 25,
          throughput: 1547,
          errorRate: 0.3,
          uptime: '15d 8h 42m',
          activeConnections: 156,
          queueDepth: 23,
          cacheHitRate: 94.8,
          gcPauseTime: 2.1,
          goroutineCount: 847,
          heapSize: 256,
          allocRate: 12.4
        })
      } finally {
        setIsLoading(false)
      }
    }

    fetchSystemMetrics()

    // Set up real-time updates
    const interval = setInterval(fetchSystemMetrics, 5000)
    return () => clearInterval(interval)
  }, [])

  const getHealthStatus = (value: number, threshold: number = 80) => {
    if (value < threshold) return { status: 'healthy', color: 'text-green-600', icon: CheckCircle }
    if (value < 90) return { status: 'warning', color: 'text-yellow-600', icon: AlertTriangle }
    return { status: 'critical', color: 'text-red-600', icon: AlertTriangle }
  }

  const getLatencyColor = (latency: number) => {
    if (latency < 50) return 'text-green-600'
    if (latency < 100) return 'text-yellow-600'
    return 'text-red-600'
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading system metrics...</p>
          </div>
        </div>
      </div>
    )
  }

  const cpuHealth = getHealthStatus(metrics?.cpuUsage || 0)
  const memoryHealth = getHealthStatus(metrics?.memoryUsage || 0)
  const diskHealth = getHealthStatus(metrics?.diskUsage || 0)

  return (
    <div className="space-y-6">
      {/* System Health Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">CPU Usage</CardTitle>
            <Cpu className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${cpuHealth.color}`}>
              {metrics?.cpuUsage.toFixed(1)}%
            </div>
            <Progress value={metrics?.cpuUsage || 0} className="mt-2" />
            <div className="flex items-center gap-1 mt-2">
              <cpuHealth.icon className="h-3 w-3" />
              <span className="text-xs capitalize">{cpuHealth.status}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
            <MemoryStick className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${memoryHealth.color}`}>
              {metrics?.memoryUsage.toFixed(1)}%
            </div>
            <Progress value={metrics?.memoryUsage || 0} className="mt-2" />
            <div className="flex items-center gap-1 mt-2">
              <memoryHealth.icon className="h-3 w-3" />
              <span className="text-xs capitalize">{memoryHealth.status}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Disk Usage</CardTitle>
            <HardDrive className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${diskHealth.color}`}>
              {metrics?.diskUsage.toFixed(1)}%
            </div>
            <Progress value={metrics?.diskUsage || 0} className="mt-2" />
            <div className="flex items-center gap-1 mt-2">
              <diskHealth.icon className="h-3 w-3" />
              <span className="text-xs capitalize">{diskHealth.status}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">System Uptime</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {metrics?.uptime}
            </div>
            <p className="text-xs text-muted-foreground mt-2">
              Continuous operation
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Alerts */}
      {(metrics?.cpuUsage || 0) > 80 && (
        <Alert>
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            High CPU usage detected ({metrics?.cpuUsage.toFixed(1)}%). Consider scaling resources.
          </AlertDescription>
        </Alert>
      )}

      {/* Detailed Monitoring Tabs */}
      <Tabs defaultValue="resources" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="resources">Resources</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="network">Network</TabsTrigger>
          <TabsTrigger value="runtime">Runtime</TabsTrigger>
        </TabsList>

        <TabsContent value="resources" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Resource Utilization</CardTitle>
                <CardDescription>Current system resource usage</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium flex items-center gap-2">
                      <Cpu className="h-4 w-4" />
                      CPU
                    </span>
                    <span className="text-sm">{metrics?.cpuUsage.toFixed(1)}%</span>
                  </div>
                  <Progress value={metrics?.cpuUsage || 0} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium flex items-center gap-2">
                      <MemoryStick className="h-4 w-4" />
                      Memory
                    </span>
                    <span className="text-sm">{metrics?.memoryUsage.toFixed(1)}%</span>
                  </div>
                  <Progress value={metrics?.memoryUsage || 0} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium flex items-center gap-2">
                      <HardDrive className="h-4 w-4" />
                      Disk
                    </span>
                    <span className="text-sm">{metrics?.diskUsage.toFixed(1)}%</span>
                  </div>
                  <Progress value={metrics?.diskUsage || 0} className="h-2" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Cache Performance</CardTitle>
                <CardDescription>Cache efficiency metrics</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Hit Rate</span>
                  <span className="text-sm font-medium text-green-600">
                    {metrics?.cacheHitRate.toFixed(1)}%
                  </span>
                </div>
                <Progress value={metrics?.cacheHitRate || 0} className="h-2" />
                <p className="text-xs text-muted-foreground">
                  {(metrics?.cacheHitRate || 0) > 90 ? 'Excellent' :
                   (metrics?.cacheHitRate || 0) > 80 ? 'Good' : 'Poor'} cache performance
                </p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="performance" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Throughput</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.throughput.toLocaleString()}</div>
                <p className="text-xs text-muted-foreground">Requests per second</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Error Rate</CardTitle>
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${(metrics?.errorRate || 0) < 1 ? 'text-green-600' : 'text-red-600'}`}>
                  {metrics?.errorRate.toFixed(2)}%
                </div>
                <p className="text-xs text-muted-foreground">Failed requests</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Active Connections</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metrics?.activeConnections}</div>
                <p className="text-xs text-muted-foreground">Current connections</p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="network" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Network className="h-4 w-4" />
                  Network Latency
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${getLatencyColor(metrics?.networkLatency || 0)}`}>
                  {metrics?.networkLatency}ms
                </div>
                <p className="text-xs text-muted-foreground">Network response time</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-4 w-4" />
                  Database Latency
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${getLatencyColor(metrics?.databaseLatency || 0)}`}>
                  {metrics?.databaseLatency}ms
                </div>
                <p className="text-xs text-muted-foreground">Database response time</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Wifi className="h-4 w-4" />
                  API Latency
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${getLatencyColor(metrics?.apiLatency || 0)}`}>
                  {metrics?.apiLatency}ms
                </div>
                <p className="text-xs text-muted-foreground">API response time</p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="runtime" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Go Runtime Metrics</CardTitle>
                <CardDescription>Go-specific performance metrics</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Goroutines</span>
                  <span className="text-sm font-medium">{metrics?.goroutineCount.toLocaleString()}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Heap Size</span>
                  <span className="text-sm font-medium">{metrics?.heapSize}MB</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">GC Pause</span>
                  <span className="text-sm font-medium">{metrics?.gcPauseTime}ms</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Alloc Rate</span>
                  <span className="text-sm font-medium">{metrics?.allocRate}MB/s</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Queue Metrics</CardTitle>
                <CardDescription>Request queue performance</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Queue Depth</span>
                  <span className="text-sm font-medium">{metrics?.queueDepth}</span>
                </div>
                <Progress value={Math.min((metrics?.queueDepth || 0) / 100 * 100, 100)} className="h-2" />
                <p className="text-xs text-muted-foreground">
                  {(metrics?.queueDepth || 0) < 50 ? 'Normal' :
                   (metrics?.queueDepth || 0) < 100 ? 'High' : 'Critical'} queue depth
                </p>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
