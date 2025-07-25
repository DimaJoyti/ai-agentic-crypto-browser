'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Shield, 
  ShieldCheck, 
  ShieldAlert, 
  Key, 
  Smartphone, 
  Globe, 
  Clock, 
  AlertTriangle,
  CheckCircle,
  XCircle,
  Monitor,
  MapPin,
  Trash2,
  Settings
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { formatDistanceToNow } from 'date-fns'

interface SecuritySession {
  id: string
  device: string
  ip_address: string
  location: string
  last_seen: string
  is_current: boolean
}

interface SecurityEvent {
  id: string
  type: 'login' | 'logout' | 'password_change' | 'mfa_enable' | 'mfa_disable' | 'suspicious_activity'
  description: string
  ip_address: string
  location: string
  timestamp: string
  success: boolean
}

interface SecuritySettings {
  mfa_enabled: boolean
  mfa_methods: string[]
  session_timeout: number
  password_last_changed: string
  login_notifications: boolean
  suspicious_activity_alerts: boolean
}

export default function SecurityDashboard() {
  const { user } = useAuth()
  const [sessions, setSessions] = useState<SecuritySession[]>([])
  const [securityEvents, setSecurityEvents] = useState<SecurityEvent[]>([])
  const [securitySettings, setSecuritySettings] = useState<SecuritySettings | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadSecurityData()
  }, [])

  const loadSecurityData = async () => {
    setIsLoading(true)
    try {
      // Load sessions
      const sessionsResponse = await fetch('/api/auth/sessions', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      })

      if (sessionsResponse.ok) {
        const sessionsData = await sessionsResponse.json()
        setSessions(sessionsData.sessions || [])
      }

      // Load security events
      const eventsResponse = await fetch('/api/security/events', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      })

      if (eventsResponse.ok) {
        const eventsData = await eventsResponse.json()
        setSecurityEvents(eventsData.events || [])
      }

      // Load security settings
      const settingsResponse = await fetch('/api/security/settings', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      })

      if (settingsResponse.ok) {
        const settingsData = await settingsResponse.json()
        setSecuritySettings(settingsData)
      } else {
        // Mock data for demonstration
        setSecuritySettings({
          mfa_enabled: user?.mfa_enabled || false,
          mfa_methods: user?.mfa_enabled ? ['totp'] : [],
          session_timeout: 24,
          password_last_changed: '2024-01-15T10:30:00Z',
          login_notifications: true,
          suspicious_activity_alerts: true,
        })

        setSessions([
          {
            id: 'session-1',
            device: 'Chrome on Windows 11',
            ip_address: '192.168.1.100',
            location: 'New York, US',
            last_seen: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
            is_current: true,
          },
          {
            id: 'session-2',
            device: 'Safari on iPhone 15',
            ip_address: '192.168.1.101',
            location: 'New York, US',
            last_seen: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
            is_current: false,
          },
          {
            id: 'session-3',
            device: 'Firefox on macOS',
            ip_address: '10.0.0.50',
            location: 'San Francisco, US',
            last_seen: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
            is_current: false,
          },
        ])

        setSecurityEvents([
          {
            id: 'event-1',
            type: 'login',
            description: 'Successful login',
            ip_address: '192.168.1.100',
            location: 'New York, US',
            timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
            success: true,
          },
          {
            id: 'event-2',
            type: 'password_change',
            description: 'Password changed successfully',
            ip_address: '192.168.1.100',
            location: 'New York, US',
            timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 7).toISOString(),
            success: true,
          },
          {
            id: 'event-3',
            type: 'login',
            description: 'Failed login attempt',
            ip_address: '203.0.113.1',
            location: 'Unknown',
            timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2).toISOString(),
            success: false,
          },
        ])
      }
    } catch (err) {
      setError('Failed to load security data')
    } finally {
      setIsLoading(false)
    }
  }

  const revokeSession = async (sessionId: string) => {
    try {
      const response = await fetch(`/api/auth/sessions/${sessionId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
      })

      if (response.ok) {
        setSessions(sessions.filter(s => s.id !== sessionId))
      } else {
        setError('Failed to revoke session')
      }
    } catch (err) {
      setError('Network error occurred')
    }
  }

  const getSecurityScore = () => {
    let score = 0
    let maxScore = 100

    // MFA enabled (+30 points)
    if (securitySettings?.mfa_enabled) {
      score += 30
    }

    // Recent password change (+20 points)
    if (securitySettings?.password_last_changed) {
      const lastChanged = new Date(securitySettings.password_last_changed)
      const daysSinceChange = (Date.now() - lastChanged.getTime()) / (1000 * 60 * 60 * 24)
      if (daysSinceChange < 90) {
        score += 20
      } else if (daysSinceChange < 180) {
        score += 10
      }
    }

    // Limited active sessions (+20 points)
    const activeSessions = sessions.filter(s => s.is_current || 
      new Date(s.last_seen).getTime() > Date.now() - 1000 * 60 * 60 * 24)
    if (activeSessions.length <= 3) {
      score += 20
    } else if (activeSessions.length <= 5) {
      score += 10
    }

    // No recent failed logins (+15 points)
    const recentFailedLogins = securityEvents.filter(e => 
      e.type === 'login' && !e.success && 
      new Date(e.timestamp).getTime() > Date.now() - 1000 * 60 * 60 * 24 * 7)
    if (recentFailedLogins.length === 0) {
      score += 15
    } else if (recentFailedLogins.length <= 2) {
      score += 7
    }

    // Notifications enabled (+15 points)
    if (securitySettings?.login_notifications && securitySettings?.suspicious_activity_alerts) {
      score += 15
    } else if (securitySettings?.login_notifications || securitySettings?.suspicious_activity_alerts) {
      score += 7
    }

    return Math.min(score, maxScore)
  }

  const getSecurityScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600'
    if (score >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getSecurityScoreIcon = (score: number) => {
    if (score >= 80) return <ShieldCheck className="h-5 w-5 text-green-600" />
    if (score >= 60) return <Shield className="h-5 w-5 text-yellow-600" />
    return <ShieldAlert className="h-5 w-5 text-red-600" />
  }

  const getEventIcon = (event: SecurityEvent) => {
    if (!event.success) {
      return <XCircle className="h-4 w-4 text-red-500" />
    }

    switch (event.type) {
      case 'login':
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case 'password_change':
        return <Key className="h-4 w-4 text-blue-500" />
      case 'mfa_enable':
      case 'mfa_disable':
        return <Smartphone className="h-4 w-4 text-purple-500" />
      default:
        return <AlertTriangle className="h-4 w-4 text-yellow-500" />
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
          <div className="spinner mb-4" />
          <p className="text-muted-foreground">Loading security dashboard...</p>
        </div>
      </div>
    )
  }

  const securityScore = getSecurityScore()

  return (
    <div className="space-y-6">
      {/* Security Score Overview */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            {getSecurityScoreIcon(securityScore)}
            Security Score
          </CardTitle>
          <CardDescription>
            Your account security rating based on current settings and activity
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <div className={`text-3xl font-bold ${getSecurityScoreColor(securityScore)}`}>
                {securityScore}/100
              </div>
              <p className="text-sm text-muted-foreground">
                {securityScore >= 80 ? 'Excellent' : securityScore >= 60 ? 'Good' : 'Needs Improvement'}
              </p>
            </div>
            <div className="text-right">
              <div className="w-32 bg-gray-200 rounded-full h-2">
                <div 
                  className={`h-2 rounded-full ${
                    securityScore >= 80 ? 'bg-green-600' : 
                    securityScore >= 60 ? 'bg-yellow-600' : 'bg-red-600'
                  }`}
                  style={{ width: `${securityScore}%` }}
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Security Recommendations */}
      {securityScore < 80 && (
        <Alert>
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            <strong>Security Recommendations:</strong>
            <ul className="mt-2 space-y-1">
              {!securitySettings?.mfa_enabled && (
                <li>• Enable two-factor authentication for better security</li>
              )}
              {securitySettings?.password_last_changed && 
               new Date(securitySettings.password_last_changed).getTime() < Date.now() - 1000 * 60 * 60 * 24 * 90 && (
                <li>• Consider changing your password (last changed {formatDistanceToNow(new Date(securitySettings.password_last_changed))} ago)</li>
              )}
              {sessions.length > 5 && (
                <li>• Review and revoke unnecessary active sessions</li>
              )}
            </ul>
          </AlertDescription>
        </Alert>
      )}

      <Tabs defaultValue="sessions" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="sessions">Active Sessions</TabsTrigger>
          <TabsTrigger value="activity">Security Activity</TabsTrigger>
          <TabsTrigger value="settings">Security Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="sessions" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Monitor className="h-5 w-5" />
                Active Sessions
              </CardTitle>
              <CardDescription>
                Manage devices and locations where you're signed in
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {sessions.map((session) => (
                  <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <Monitor className="h-5 w-5 text-muted-foreground" />
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{session.device}</span>
                          {session.is_current && (
                            <Badge variant="secondary" className="text-xs">Current</Badge>
                          )}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          <div className="flex items-center gap-1">
                            <Globe className="h-3 w-3" />
                            {session.ip_address}
                          </div>
                          <div className="flex items-center gap-1">
                            <MapPin className="h-3 w-3" />
                            {session.location}
                          </div>
                          <div className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            Last seen {formatDistanceToNow(new Date(session.last_seen))} ago
                          </div>
                        </div>
                      </div>
                    </div>
                    {!session.is_current && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => revokeSession(session.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="activity" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Recent Security Activity</CardTitle>
              <CardDescription>
                Recent login attempts and security-related actions
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {securityEvents.map((event) => (
                  <div key={event.id} className="flex items-start gap-3 p-3 border rounded-lg">
                    {getEventIcon(event)}
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span className="font-medium">{event.description}</span>
                        <span className="text-xs text-muted-foreground">
                          {formatDistanceToNow(new Date(event.timestamp))} ago
                        </span>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {event.ip_address} • {event.location}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Security Settings
              </CardTitle>
              <CardDescription>
                Configure your account security preferences
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Two-Factor Authentication</h3>
                  <p className="text-sm text-muted-foreground">
                    Add an extra layer of security to your account
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  {securitySettings?.mfa_enabled ? (
                    <Badge variant="secondary" className="text-green-600">Enabled</Badge>
                  ) : (
                    <Badge variant="outline">Disabled</Badge>
                  )}
                  <Button variant="outline" size="sm">
                    {securitySettings?.mfa_enabled ? 'Manage' : 'Enable'}
                  </Button>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Login Notifications</h3>
                  <p className="text-sm text-muted-foreground">
                    Get notified of new sign-ins to your account
                  </p>
                </div>
                <Badge variant={securitySettings?.login_notifications ? "secondary" : "outline"}>
                  {securitySettings?.login_notifications ? 'Enabled' : 'Disabled'}
                </Badge>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Suspicious Activity Alerts</h3>
                  <p className="text-sm text-muted-foreground">
                    Get alerted about unusual account activity
                  </p>
                </div>
                <Badge variant={securitySettings?.suspicious_activity_alerts ? "secondary" : "outline"}>
                  {securitySettings?.suspicious_activity_alerts ? 'Enabled' : 'Disabled'}
                </Badge>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Session Timeout</h3>
                  <p className="text-sm text-muted-foreground">
                    Automatically sign out after {securitySettings?.session_timeout} hours of inactivity
                  </p>
                </div>
                <Button variant="outline" size="sm">
                  Configure
                </Button>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">Password</h3>
                  <p className="text-sm text-muted-foreground">
                    Last changed {securitySettings?.password_last_changed ? 
                      formatDistanceToNow(new Date(securitySettings.password_last_changed)) + ' ago' : 
                      'never'}
                  </p>
                </div>
                <Button variant="outline" size="sm">
                  Change Password
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </div>
  )
}
