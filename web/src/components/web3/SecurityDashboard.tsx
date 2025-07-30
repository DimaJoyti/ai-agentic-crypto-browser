'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAccount } from 'wagmi'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield,
  ShieldCheck,
  ShieldAlert,
  ShieldX,
  Key,
  Lock,
  Unlock,
  Eye,
  EyeOff,
  Clock,
  Activity,
  AlertTriangle,
  CheckCircle,
  XCircle,
  RefreshCw,
  Trash2,
  Download,
  Settings,
  Fingerprint,
  Smartphone,
  Monitor,
  Globe
} from 'lucide-react'
import { useWalletAuth } from '@/hooks/useWalletAuth'
import { useSecureStorage } from '@/lib/secure-storage'
import { walletSecurity } from '@/lib/wallet-security'
import { toast } from 'sonner'

export function SecurityDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [securityStats, setSecurityStats] = useState<any>(null)
  const [storageStats, setStorageStats] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(false)

  const { address } = useAccount()
  const {
    authState,
    authenticate,
    logout,
    getActiveSessions,
    terminateAllSessions,
    getSecurityEvents,
    clearError
  } = useWalletAuth({ requireAuth: false })

  const {
    getSecureStorageStats,
    cleanupSecureStorage,
    clearSecureStorage
  } = useSecureStorage()

  // Load security data
  useEffect(() => {
    loadSecurityData()
  }, [address])

  const loadSecurityData = async () => {
    setIsLoading(true)
    try {
      // Get security statistics
      const stats = walletSecurity.getSecurityStats()
      setSecurityStats(stats)

      // Get storage statistics
      const storage = getSecureStorageStats()
      setStorageStats(storage)
    } catch (error) {
      console.error('Failed to load security data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleAuthenticate = async () => {
    const success = await authenticate()
    if (success) {
      loadSecurityData()
    }
  }

  const handleLogout = () => {
    logout()
    loadSecurityData()
  }

  const handleTerminateAllSessions = () => {
    terminateAllSessions()
    loadSecurityData()
  }

  const handleCleanupStorage = () => {
    const removed = cleanupSecureStorage()
    toast.success(`Cleaned up ${removed} expired items`)
    loadSecurityData()
  }

  const handleClearStorage = () => {
    if (window.confirm('Are you sure you want to clear all secure storage? This action cannot be undone.')) {
      clearSecureStorage()
      toast.success('Secure storage cleared')
      loadSecurityData()
    }
  }

  const formatTimeAgo = (timestamp: number) => {
    const now = Date.now()
    const diff = now - timestamp
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`
    if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`
    if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`
    return 'Just now'
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const getSecurityLevel = () => {
    if (!authState.isAuthenticated) return { level: 'Low', color: 'text-red-600', icon: ShieldX }
    if (authState.session && authState.session.riskScore < 20) return { level: 'High', color: 'text-green-600', icon: ShieldCheck }
    if (authState.session && authState.session.riskScore < 50) return { level: 'Medium', color: 'text-yellow-600', icon: ShieldAlert }
    return { level: 'Low', color: 'text-red-600', icon: ShieldX }
  }

  const securityLevel = getSecurityLevel()
  const SecurityIcon = securityLevel.icon

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Security Dashboard</h2>
          <p className="text-muted-foreground">
            Monitor and manage your wallet security settings
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={loadSecurityData}
            disabled={isLoading}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Security Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Security Level</p>
                <p className={`text-2xl font-bold ${securityLevel.color}`}>{securityLevel.level}</p>
              </div>
              <SecurityIcon className={`w-8 h-8 ${securityLevel.color}`} />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {authState.isAuthenticated ? 'Wallet authenticated' : 'Authentication required'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Active Sessions</p>
                <p className="text-2xl font-bold">{securityStats?.activeSessions || 0}</p>
              </div>
              <Activity className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Across all devices
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Security Events</p>
                <p className="text-2xl font-bold">{securityStats?.totalEvents || 0}</p>
              </div>
              <AlertTriangle className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {securityStats?.highRiskEvents || 0} high-risk events
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Secure Storage</p>
                <p className="text-2xl font-bold">{storageStats?.totalItems || 0}</p>
              </div>
              <Lock className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {storageStats ? formatBytes(storageStats.totalSize) : '0 Bytes'} used
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Authentication Status */}
      {authState.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {authState.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="sessions">Sessions</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
          <TabsTrigger value="storage">Storage</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            {/* Authentication Status */}
            <Card>
              <CardHeader>
                <CardTitle>Authentication Status</CardTitle>
                <CardDescription>Current wallet authentication state</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                      authState.isAuthenticated ? 'bg-green-100 dark:bg-green-900' : 'bg-red-100 dark:bg-red-900'
                    }`}>
                      {authState.isAuthenticated ? (
                        <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
                      ) : (
                        <XCircle className="w-5 h-5 text-red-600 dark:text-red-400" />
                      )}
                    </div>
                    <div>
                      <p className="font-medium">
                        {authState.isAuthenticated ? 'Authenticated' : 'Not Authenticated'}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {authState.lastAuthTime 
                          ? `Last auth: ${formatTimeAgo(authState.lastAuthTime)}`
                          : 'Never authenticated'
                        }
                      </p>
                    </div>
                  </div>
                  {authState.isAuthenticated ? (
                    <Button variant="outline" onClick={handleLogout}>
                      <Unlock className="w-4 h-4 mr-2" />
                      Logout
                    </Button>
                  ) : (
                    <Button onClick={handleAuthenticate} disabled={authState.isAuthenticating}>
                      <Lock className="w-4 h-4 mr-2" />
                      {authState.isAuthenticating ? 'Authenticating...' : 'Authenticate'}
                    </Button>
                  )}
                </div>

                {authState.session && (
                  <div className="pt-4 border-t space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Session ID</span>
                      <span className="font-mono">{authState.session.id.slice(0, 8)}...</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Risk Score</span>
                      <Badge variant={authState.session.riskScore < 20 ? 'default' : 
                                   authState.session.riskScore < 50 ? 'secondary' : 'destructive'}>
                        {authState.session.riskScore}/100
                      </Badge>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Expires</span>
                      <span>{formatTimeAgo(authState.session.expiresAt)}</span>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Security Settings */}
            <Card>
              <CardHeader>
                <CardTitle>Security Settings</CardTitle>
                <CardDescription>Configure your security preferences</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">Signature Authentication</h4>
                    <p className="text-sm text-muted-foreground">
                      Require signature verification for sensitive actions
                    </p>
                  </div>
                  <Badge variant="default">Enabled</Badge>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">Secure Storage</h4>
                    <p className="text-sm text-muted-foreground">
                      Encrypt sensitive data in local storage
                    </p>
                  </div>
                  <Badge variant="default">Enabled</Badge>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">Session Timeout</h4>
                    <p className="text-sm text-muted-foreground">
                      Automatic logout after inactivity
                    </p>
                  </div>
                  <Badge variant="outline">24 hours</Badge>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">Failed Attempt Lockout</h4>
                    <p className="text-sm text-muted-foreground">
                      Lock account after failed attempts
                    </p>
                  </div>
                  <Badge variant="outline">5 attempts</Badge>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="sessions" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Active Sessions</CardTitle>
              <CardDescription>Manage your active wallet sessions</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {address && getActiveSessions().map((session) => (
                  <div key={session.id} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                        <Monitor className="w-5 h-5" />
                      </div>
                      <div>
                        <p className="font-medium">Session {session.id.slice(0, 8)}</p>
                        <p className="text-sm text-muted-foreground">
                          Started {formatTimeAgo(session.startTime)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Risk Score: {session.riskScore}/100
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={session.isActive ? 'default' : 'secondary'}>
                        {session.isActive ? 'Active' : 'Inactive'}
                      </Badge>
                      {session.id === authState.session?.id && (
                        <Badge variant="outline">Current</Badge>
                      )}
                    </div>
                  </div>
                ))}

                {(!address || getActiveSessions().length === 0) && (
                  <div className="text-center py-8 text-muted-foreground">
                    No active sessions
                  </div>
                )}

                {address && getActiveSessions().length > 0 && (
                  <div className="pt-4 border-t">
                    <Button variant="destructive" onClick={handleTerminateAllSessions}>
                      <XCircle className="w-4 h-4 mr-2" />
                      Terminate All Sessions
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="events" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Security Events</CardTitle>
              <CardDescription>Recent security-related activities</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {address && getSecurityEvents().slice(0, 10).map((event) => (
                  <div key={event.id} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                        event.riskLevel === 'low' ? 'bg-green-100 dark:bg-green-900' :
                        event.riskLevel === 'medium' ? 'bg-yellow-100 dark:bg-yellow-900' :
                        event.riskLevel === 'high' ? 'bg-orange-100 dark:bg-orange-900' :
                        'bg-red-100 dark:bg-red-900'
                      }`}>
                        {event.type === 'login' && <CheckCircle className="w-4 h-4 text-green-600" />}
                        {event.type === 'logout' && <XCircle className="w-4 h-4 text-gray-600" />}
                        {event.type === 'failed_auth' && <AlertTriangle className="w-4 h-4 text-red-600" />}
                        {event.type === 'suspicious_activity' && <Shield className="w-4 h-4 text-orange-600" />}
                      </div>
                      <div>
                        <p className="font-medium capitalize">{event.type.replace('_', ' ')}</p>
                        <p className="text-sm text-muted-foreground">
                          {formatTimeAgo(event.timestamp)}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={
                        event.riskLevel === 'low' ? 'default' :
                        event.riskLevel === 'medium' ? 'secondary' :
                        event.riskLevel === 'high' ? 'destructive' :
                        'destructive'
                      }>
                        {event.riskLevel}
                      </Badge>
                      {event.resolved && (
                        <Badge variant="outline">Resolved</Badge>
                      )}
                    </div>
                  </div>
                ))}

                {(!address || getSecurityEvents().length === 0) && (
                  <div className="text-center py-8 text-muted-foreground">
                    No security events recorded
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="storage" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Secure Storage</CardTitle>
              <CardDescription>Manage encrypted local storage</CardDescription>
            </CardHeader>
            <CardContent>
              {storageStats && (
                <div className="space-y-6">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <h4 className="font-medium mb-2">Storage Statistics</h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Total Items</span>
                          <span>{storageStats.totalItems}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Total Size</span>
                          <span>{formatBytes(storageStats.totalSize)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Encrypted Items</span>
                          <span>{storageStats.encryptedItems}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Expired Items</span>
                          <span>{storageStats.expiredItems}</span>
                        </div>
                      </div>
                    </div>

                    <div>
                      <h4 className="font-medium mb-2">Storage Actions</h4>
                      <div className="space-y-2">
                        <Button variant="outline" size="sm" onClick={handleCleanupStorage} className="w-full">
                          <Trash2 className="w-4 h-4 mr-2" />
                          Cleanup Expired Items
                        </Button>
                        <Button variant="destructive" size="sm" onClick={handleClearStorage} className="w-full">
                          <XCircle className="w-4 h-4 mr-2" />
                          Clear All Storage
                        </Button>
                      </div>
                    </div>
                  </div>

                  {storageStats.oldestItem && (
                    <div className="pt-4 border-t">
                      <div className="grid gap-2 md:grid-cols-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Oldest Item</span>
                          <span>{formatTimeAgo(storageStats.oldestItem)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Newest Item</span>
                          <span>{formatTimeAgo(storageStats.newestItem)}</span>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
