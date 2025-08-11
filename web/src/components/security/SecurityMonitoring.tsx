'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield, 
  AlertTriangle,
  Eye,
  Activity,
  Globe,
  Clock,
  TrendingUp,
  TrendingDown,
  MapPin,
  Monitor,
  Wifi,
  Lock,
  Unlock,
  Ban,
  CheckCircle,
  XCircle,
  Zap,
  Target,
  Bell,
  Settings
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface SecurityThreat {
  id: string
  type: 'brute_force' | 'suspicious_login' | 'unusual_activity' | 'malware' | 'phishing' | 'data_breach'
  severity: 'low' | 'medium' | 'high' | 'critical'
  title: string
  description: string
  timestamp: number
  source: string
  ipAddress: string
  location: string
  status: 'active' | 'investigating' | 'resolved' | 'false_positive'
  affectedAssets: string[]
  recommendedActions: string[]
}

interface SecurityMetric {
  id: string
  name: string
  value: number
  previousValue: number
  unit: string
  trend: 'up' | 'down' | 'stable'
  status: 'good' | 'warning' | 'critical'
  description: string
}

interface GeolocationData {
  country: string
  city: string
  latitude: number
  longitude: number
  loginAttempts: number
  successfulLogins: number
  blockedAttempts: number
  riskScore: number
}

interface DeviceFingerprint {
  id: string
  deviceType: string
  browser: string
  os: string
  screenResolution: string
  timezone: string
  language: string
  firstSeen: number
  lastSeen: number
  loginCount: number
  riskScore: number
  isTrusted: boolean
  isBlocked: boolean
}

export function SecurityMonitoring() {
  const [threats, setThreats] = useState<SecurityThreat[]>([])
  const [metrics, setMetrics] = useState<SecurityMetric[]>([])
  const [geoData, setGeoData] = useState<GeolocationData[]>([])
  const [devices, setDevices] = useState<DeviceFingerprint[]>([])
  const [activeTab, setActiveTab] = useState('overview')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock security monitoring data
    const mockThreats: SecurityThreat[] = [
      {
        id: 'threat1',
        type: 'brute_force',
        severity: 'high',
        title: 'Brute Force Attack Detected',
        description: 'Multiple failed login attempts from suspicious IP addresses',
        timestamp: Date.now() - 1800000,
        source: 'Login Monitor',
        ipAddress: '203.0.113.1',
        location: 'Unknown Location',
        status: 'investigating',
        affectedAssets: ['Login System'],
        recommendedActions: ['Block IP address', 'Enable rate limiting', 'Notify user']
      },
      {
        id: 'threat2',
        type: 'suspicious_login',
        severity: 'medium',
        title: 'Login from New Location',
        description: 'User logged in from a previously unseen geographic location',
        timestamp: Date.now() - 3600000,
        source: 'Geo Monitor',
        ipAddress: '198.51.100.1',
        location: 'Moscow, Russia',
        status: 'resolved',
        affectedAssets: ['User Account'],
        recommendedActions: ['Verify with user', 'Enable MFA']
      },
      {
        id: 'threat3',
        type: 'unusual_activity',
        severity: 'medium',
        title: 'Unusual Trading Pattern',
        description: 'Large volume trades outside normal user behavior',
        timestamp: Date.now() - 7200000,
        source: 'Behavior Analytics',
        ipAddress: '192.168.1.100',
        location: 'New York, US',
        status: 'false_positive',
        affectedAssets: ['Trading System'],
        recommendedActions: ['Monitor closely', 'Request verification']
      }
    ]

    const mockMetrics: SecurityMetric[] = [
      {
        id: 'metric1',
        name: 'Failed Login Attempts',
        value: 23,
        previousValue: 18,
        unit: 'attempts/hour',
        trend: 'up',
        status: 'warning',
        description: 'Number of failed login attempts in the last hour'
      },
      {
        id: 'metric2',
        name: 'Blocked IPs',
        value: 156,
        previousValue: 142,
        unit: 'IPs',
        trend: 'up',
        status: 'good',
        description: 'Total number of blocked IP addresses'
      },
      {
        id: 'metric3',
        name: 'MFA Success Rate',
        value: 98.5,
        previousValue: 97.8,
        unit: '%',
        trend: 'up',
        status: 'good',
        description: 'Percentage of successful MFA verifications'
      },
      {
        id: 'metric4',
        name: 'Suspicious Activities',
        value: 7,
        previousValue: 12,
        unit: 'events',
        trend: 'down',
        status: 'good',
        description: 'Number of flagged suspicious activities'
      }
    ]

    const mockGeoData: GeolocationData[] = [
      {
        country: 'United States',
        city: 'New York',
        latitude: 40.7128,
        longitude: -74.0060,
        loginAttempts: 1247,
        successfulLogins: 1198,
        blockedAttempts: 49,
        riskScore: 15
      },
      {
        country: 'United Kingdom',
        city: 'London',
        latitude: 51.5074,
        longitude: -0.1278,
        loginAttempts: 342,
        successfulLogins: 338,
        blockedAttempts: 4,
        riskScore: 8
      },
      {
        country: 'Germany',
        city: 'Berlin',
        latitude: 52.5200,
        longitude: 13.4050,
        loginAttempts: 156,
        successfulLogins: 152,
        blockedAttempts: 4,
        riskScore: 12
      },
      {
        country: 'Russia',
        city: 'Moscow',
        latitude: 55.7558,
        longitude: 37.6176,
        loginAttempts: 89,
        successfulLogins: 23,
        blockedAttempts: 66,
        riskScore: 85
      }
    ]

    const mockDevices: DeviceFingerprint[] = [
      {
        id: 'device1',
        deviceType: 'Desktop',
        browser: 'Chrome 118.0',
        os: 'Windows 11',
        screenResolution: '1920x1080',
        timezone: 'America/New_York',
        language: 'en-US',
        firstSeen: Date.now() - 86400000 * 30,
        lastSeen: Date.now() - 300000,
        loginCount: 247,
        riskScore: 5,
        isTrusted: true,
        isBlocked: false
      },
      {
        id: 'device2',
        deviceType: 'Mobile',
        browser: 'Safari 17.0',
        os: 'iOS 17.1',
        screenResolution: '390x844',
        timezone: 'America/New_York',
        language: 'en-US',
        firstSeen: Date.now() - 86400000 * 60,
        lastSeen: Date.now() - 3600000,
        loginCount: 89,
        riskScore: 8,
        isTrusted: true,
        isBlocked: false
      },
      {
        id: 'device3',
        deviceType: 'Desktop',
        browser: 'Firefox 119.0',
        os: 'Linux',
        screenResolution: '1366x768',
        timezone: 'Europe/Moscow',
        language: 'ru-RU',
        firstSeen: Date.now() - 86400000 * 2,
        lastSeen: Date.now() - 86400000,
        loginCount: 3,
        riskScore: 75,
        isTrusted: false,
        isBlocked: true
      }
    ]

    setThreats(mockThreats)
    setMetrics(mockMetrics)
    setGeoData(mockGeoData)
    setDevices(mockDevices)
  }, [isConnected])

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'text-red-500'
      case 'high': return 'text-orange-500'
      case 'medium': return 'text-yellow-500'
      case 'low': return 'text-blue-500'
      default: return 'text-muted-foreground'
    }
  }

  const getSeverityBadgeVariant = (severity: string) => {
    switch (severity) {
      case 'critical': case 'high': return 'destructive'
      case 'medium': return 'secondary'
      case 'low': return 'outline'
      default: return 'outline'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'resolved': case 'false_positive': return 'text-green-500'
      case 'investigating': return 'text-yellow-500'
      case 'active': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getMetricStatusColor = (status: string) => {
    switch (status) {
      case 'good': return 'text-green-500'
      case 'warning': return 'text-yellow-500'
      case 'critical': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="w-3 h-3" />
      case 'down': return <TrendingDown className="w-3 h-3" />
      case 'stable': return <Activity className="w-3 h-3" />
      default: return null
    }
  }

  const getRiskScoreColor = (score: number) => {
    if (score >= 70) return 'text-red-500'
    if (score >= 40) return 'text-yellow-500'
    return 'text-green-500'
  }

  const getOverallSecurityScore = () => {
    const activeThreats = threats.filter(t => t.status === 'active').length
    const criticalThreats = threats.filter(t => t.severity === 'critical').length
    const avgRiskScore = geoData.reduce((sum, geo) => sum + geo.riskScore, 0) / geoData.length
    
    let score = 100
    score -= activeThreats * 10
    score -= criticalThreats * 20
    score -= avgRiskScore * 0.3
    
    return Math.max(score, 0)
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access security monitoring
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Security Monitoring</h2>
          <p className="text-muted-foreground">
            Real-time security threat detection and monitoring
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Security Score: {getOverallSecurityScore().toFixed(0)}%
          </Badge>
          <Button variant="outline" size="sm">
            <Settings className="w-3 h-3 mr-1" />
            Configure
          </Button>
        </div>
      </div>

      {/* Security Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {metrics.map((metric) => (
          <Card key={metric.id}>
            <CardContent className="p-4">
              <div className="flex items-center gap-2 mb-2">
                <Activity className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">{metric.name}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="text-2xl font-bold">{metric.value}</div>
                <div className={cn("flex items-center gap-1", getMetricStatusColor(metric.status))}>
                  {getTrendIcon(metric.trend)}
                  <span className="text-xs">{metric.unit}</span>
                </div>
              </div>
              <div className="text-xs text-muted-foreground mt-1">
                Previous: {metric.previousValue} {metric.unit}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="threats">Threats</TabsTrigger>
          <TabsTrigger value="geography">Geography</TabsTrigger>
          <TabsTrigger value="devices">Devices</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Security Score */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Security Score
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-muted-foreground">Overall Security</span>
                    <span className="font-bold text-2xl">{getOverallSecurityScore().toFixed(0)}%</span>
                  </div>
                  <Progress value={getOverallSecurityScore()} className="h-3" />
                  
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span>Active Threats</span>
                      <span className="text-red-500">{threats.filter(t => t.status === 'active').length}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Blocked IPs</span>
                      <span className="text-green-500">{metrics.find(m => m.name === 'Blocked IPs')?.value || 0}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>MFA Success Rate</span>
                      <span className="text-green-500">{metrics.find(m => m.name === 'MFA Success Rate')?.value || 0}%</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Recent Threats */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="w-5 h-5" />
                  Recent Threats
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {threats.slice(0, 5).map((threat) => (
                    <div key={threat.id} className="flex items-start gap-3 p-3 border rounded">
                      <div className={cn(
                        "w-2 h-2 rounded-full mt-2",
                        threat.severity === 'critical' ? "bg-red-500" :
                        threat.severity === 'high' ? "bg-orange-500" :
                        threat.severity === 'medium' ? "bg-yellow-500" :
                        "bg-blue-500"
                      )} />
                      <div className="flex-1">
                        <div className="font-medium">{threat.title}</div>
                        <div className="text-sm text-muted-foreground">{threat.description}</div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {formatTime(threat.timestamp)} • {threat.location}
                        </div>
                      </div>
                      <Badge variant={getSeverityBadgeVariant(threat.severity)}>
                        {threat.severity}
                      </Badge>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="threats" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Security Threats</h3>
            <Button variant="outline" size="sm">
              <Bell className="w-4 h-4 mr-2" />
              Configure Alerts
            </Button>
          </div>

          <div className="space-y-4">
            {threats.map((threat) => (
              <Card key={threat.id}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-start gap-3">
                      <div className={cn(
                        "w-3 h-3 rounded-full mt-1",
                        threat.severity === 'critical' ? "bg-red-500" :
                        threat.severity === 'high' ? "bg-orange-500" :
                        threat.severity === 'medium' ? "bg-yellow-500" :
                        "bg-blue-500"
                      )} />
                      <div>
                        <h4 className="font-bold">{threat.title}</h4>
                        <p className="text-sm text-muted-foreground">{threat.description}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={getSeverityBadgeVariant(threat.severity)}>
                        {threat.severity}
                      </Badge>
                      <Badge variant="outline" className={getStatusColor(threat.status)}>
                        {threat.status.replace('_', ' ')}
                      </Badge>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Source</div>
                      <div className="font-medium">{threat.source}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">IP Address</div>
                      <div className="font-medium font-mono">{threat.ipAddress}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Location</div>
                      <div className="font-medium">{threat.location}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Detected</div>
                      <div className="font-medium">{formatTime(threat.timestamp)}</div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="text-sm font-medium">Recommended Actions:</div>
                    <ul className="text-sm text-muted-foreground space-y-1">
                      {threat.recommendedActions.map((action, index) => (
                        <li key={index} className="flex items-center gap-2">
                          <div className="w-1 h-1 bg-muted-foreground rounded-full" />
                          {action}
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="flex gap-2 mt-4">
                    {threat.status === 'active' && (
                      <>
                        <Button size="sm">Investigate</Button>
                        <Button variant="outline" size="sm">Block IP</Button>
                        <Button variant="outline" size="sm">Mark False Positive</Button>
                      </>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="geography" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Geographic Analysis</h3>
            <Button variant="outline" size="sm">
              <Globe className="w-4 h-4 mr-2" />
              View Map
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {geoData.map((geo, index) => (
              <Card key={index}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{geo.city}, {geo.country}</h4>
                      <p className="text-sm text-muted-foreground">
                        {geo.latitude.toFixed(4)}, {geo.longitude.toFixed(4)}
                      </p>
                    </div>
                    <div className="text-right">
                      <div className={cn("font-bold", getRiskScoreColor(geo.riskScore))}>
                        Risk: {geo.riskScore}%
                      </div>
                    </div>
                  </div>

                  <div className="space-y-3 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Total Attempts</span>
                      <span className="font-medium">{geo.loginAttempts}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Successful</span>
                      <span className="font-medium text-green-500">{geo.successfulLogins}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Blocked</span>
                      <span className="font-medium text-red-500">{geo.blockedAttempts}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Success Rate</span>
                      <span className="font-medium">
                        {((geo.successfulLogins / geo.loginAttempts) * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>

                  <div className="mt-4">
                    <Progress 
                      value={(geo.successfulLogins / geo.loginAttempts) * 100} 
                      className="h-2" 
                    />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="devices" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Device Fingerprints</h3>
            <Button variant="outline" size="sm">
              <Monitor className="w-4 h-4 mr-2" />
              Device Analytics
            </Button>
          </div>

          <div className="space-y-4">
            {devices.map((device) => (
              <Card key={device.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Monitor className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{device.deviceType} - {device.browser}</h4>
                        <p className="text-sm text-muted-foreground">
                          {device.os} • {device.screenResolution}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {device.isTrusted ? (
                        <Badge variant="default">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Trusted
                        </Badge>
                      ) : (
                        <Badge variant="destructive">
                          <XCircle className="w-3 h-3 mr-1" />
                          Untrusted
                        </Badge>
                      )}
                      {device.isBlocked && (
                        <Badge variant="destructive">
                          <Ban className="w-3 h-3 mr-1" />
                          Blocked
                        </Badge>
                      )}
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Risk Score</div>
                      <div className={cn("font-medium", getRiskScoreColor(device.riskScore))}>
                        {device.riskScore}%
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Login Count</div>
                      <div className="font-medium">{device.loginCount}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">First Seen</div>
                      <div className="font-medium">{formatTime(device.firstSeen)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Last Seen</div>
                      <div className="font-medium">{formatTime(device.lastSeen)}</div>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
                    <div>
                      <div className="text-muted-foreground">Timezone</div>
                      <div className="font-medium">{device.timezone}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Language</div>
                      <div className="font-medium">{device.language}</div>
                    </div>
                    <div className="flex gap-2">
                      {!device.isTrusted && (
                        <Button variant="outline" size="sm">
                          <CheckCircle className="w-3 h-3 mr-1" />
                          Trust
                        </Button>
                      )}
                      {!device.isBlocked ? (
                        <Button variant="outline" size="sm">
                          <Ban className="w-3 h-3 mr-1" />
                          Block
                        </Button>
                      ) : (
                        <Button variant="outline" size="sm">
                          <Unlock className="w-3 h-3 mr-1" />
                          Unblock
                        </Button>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
