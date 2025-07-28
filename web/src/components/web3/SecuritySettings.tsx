'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield, 
  Lock, 
  Key, 
  Eye, 
  EyeOff, 
  AlertTriangle,
  CheckCircle,
  Clock,
  Fingerprint,
  Smartphone,
  Wifi,
  WifiOff,
  RefreshCw,
  Download,
  Trash2,
  Settings
} from 'lucide-react'
import { WalletSecurity, SecurityLogger, SecurityEventType, SecureStorage } from '@/lib/security'
import { toast } from 'sonner'

interface SecuritySettings {
  autoLock: boolean
  autoLockTimeout: number // minutes
  requirePasswordForTransactions: boolean
  enableBiometric: boolean
  enableTwoFactor: boolean
  sessionTimeout: number // hours
  encryptLocalStorage: boolean
  allowRemoteConnections: boolean
  logSecurityEvents: boolean
}

interface PasswordStrength {
  score: number
  feedback: string[]
  isValid: boolean
}

export function SecuritySettings() {
  const [settings, setSettings] = useState<SecuritySettings>({
    autoLock: true,
    autoLockTimeout: 15,
    requirePasswordForTransactions: true,
    enableBiometric: false,
    enableTwoFactor: false,
    sessionTimeout: 24,
    encryptLocalStorage: true,
    allowRemoteConnections: false,
    logSecurityEvents: true
  })

  const [passwords, setPasswords] = useState({
    current: '',
    new: '',
    confirm: ''
  })

  const [showPasswords, setShowPasswords] = useState({
    current: false,
    new: false,
    confirm: false
  })

  const [passwordStrength, setPasswordStrength] = useState<PasswordStrength>({
    score: 0,
    feedback: [],
    isValid: false
  })

  const [isChangingPassword, setIsChangingPassword] = useState(false)
  const [securityEvents, setSecurityEvents] = useState<any[]>([])
  const [activeTab, setActiveTab] = useState('general')

  useEffect(() => {
    // Load security events
    SecurityLogger.loadEvents()
    setSecurityEvents(SecurityLogger.getEvents(20))

    // Load settings from secure storage if available
    loadSecuritySettings()
  }, [])

  useEffect(() => {
    // Update password strength when new password changes
    if (passwords.new) {
      const strength = WalletSecurity.validatePassword(passwords.new)
      setPasswordStrength(strength)
    } else {
      setPasswordStrength({ score: 0, feedback: [], isValid: false })
    }
  }, [passwords.new])

  const loadSecuritySettings = async () => {
    try {
      // In a real app, this would load from secure storage
      // For demo, we'll use localStorage
      const stored = localStorage.getItem('security_settings')
      if (stored) {
        setSettings(JSON.parse(stored))
      }
    } catch (error) {
      console.warn('Failed to load security settings:', error)
    }
  }

  const saveSecuritySettings = async (newSettings: SecuritySettings) => {
    try {
      setSettings(newSettings)
      localStorage.setItem('security_settings', JSON.stringify(newSettings))
      
      SecurityLogger.logEvent(
        SecurityEventType.PASSWORD_CHANGED,
        { settingsChanged: true },
        'medium'
      )
      
      toast.success('Security settings updated successfully')
    } catch (error) {
      toast.error('Failed to save security settings')
    }
  }

  const handleSettingChange = (key: keyof SecuritySettings, value: any) => {
    const newSettings = { ...settings, [key]: value }
    saveSecuritySettings(newSettings)
  }

  const handlePasswordChange = async () => {
    if (!passwords.current || !passwords.new || !passwords.confirm) {
      toast.error('Please fill in all password fields')
      return
    }

    if (passwords.new !== passwords.confirm) {
      toast.error('New passwords do not match')
      return
    }

    if (!passwordStrength.isValid) {
      toast.error('Password does not meet security requirements')
      return
    }

    setIsChangingPassword(true)

    try {
      // In a real app, this would verify the current password and update it
      await new Promise(resolve => setTimeout(resolve, 1000)) // Simulate API call

      SecurityLogger.logEvent(
        SecurityEventType.PASSWORD_CHANGED,
        { timestamp: Date.now() },
        'high'
      )

      setPasswords({ current: '', new: '', confirm: '' })
      toast.success('Password changed successfully')
    } catch (error) {
      toast.error('Failed to change password')
    } finally {
      setIsChangingPassword(false)
    }
  }

  const clearSecurityEvents = () => {
    SecurityLogger.clearEvents()
    setSecurityEvents([])
    toast.success('Security events cleared')
  }

  const exportSecurityData = () => {
    const data = {
      settings,
      events: securityEvents,
      timestamp: new Date().toISOString()
    }
    
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `wallet-security-${Date.now()}.json`
    a.click()
    URL.revokeObjectURL(url)
    
    toast.success('Security data exported')
  }

  const getPasswordStrengthColor = (score: number) => {
    if (score <= 2) return 'bg-red-500'
    if (score <= 4) return 'bg-yellow-500'
    return 'bg-green-500'
  }

  const getPasswordStrengthText = (score: number) => {
    if (score <= 2) return 'Weak'
    if (score <= 4) return 'Medium'
    return 'Strong'
  }

  const getEventSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'bg-red-500'
      case 'high': return 'bg-orange-500'
      case 'medium': return 'bg-yellow-500'
      default: return 'bg-blue-500'
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <Shield className="w-6 h-6" />
            Security Settings
          </h2>
          <p className="text-muted-foreground">
            Manage your wallet security and privacy settings
          </p>
        </div>
        <Button onClick={exportSecurityData} variant="outline" className="gap-2">
          <Download className="w-4 h-4" />
          Export Data
        </Button>
      </div>

      {/* Security Status */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
                <Shield className="w-6 h-6 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <h3 className="font-semibold">Security Status</h3>
                <p className="text-sm text-muted-foreground">Your wallet is secure</p>
              </div>
            </div>
            <Badge variant="secondary" className="bg-green-100 text-green-800">
              Protected
            </Badge>
          </div>
        </CardContent>
      </Card>

      {/* Security Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="password">Password</TabsTrigger>
          <TabsTrigger value="advanced">Advanced</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
        </TabsList>

        <TabsContent value="general" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>General Security</CardTitle>
              <CardDescription>
                Basic security settings for your wallet
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label>Auto-lock wallet</Label>
                  <p className="text-sm text-muted-foreground">
                    Automatically lock wallet after inactivity
                  </p>
                </div>
                <Switch
                  checked={settings.autoLock}
                  onCheckedChange={(checked) => handleSettingChange('autoLock', checked)}
                />
              </div>

              {settings.autoLock && (
                <div className="space-y-2">
                  <Label>Auto-lock timeout (minutes)</Label>
                  <Input
                    type="number"
                    value={settings.autoLockTimeout}
                    onChange={(e) => handleSettingChange('autoLockTimeout', parseInt(e.target.value))}
                    min="1"
                    max="60"
                  />
                </div>
              )}

              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label>Require password for transactions</Label>
                  <p className="text-sm text-muted-foreground">
                    Ask for password before signing transactions
                  </p>
                </div>
                <Switch
                  checked={settings.requirePasswordForTransactions}
                  onCheckedChange={(checked) => handleSettingChange('requirePasswordForTransactions', checked)}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label>Encrypt local storage</Label>
                  <p className="text-sm text-muted-foreground">
                    Encrypt sensitive data stored locally
                  </p>
                </div>
                <Switch
                  checked={settings.encryptLocalStorage}
                  onCheckedChange={(checked) => handleSettingChange('encryptLocalStorage', checked)}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label>Log security events</Label>
                  <p className="text-sm text-muted-foreground">
                    Keep a log of security-related activities
                  </p>
                </div>
                <Switch
                  checked={settings.logSecurityEvents}
                  onCheckedChange={(checked) => handleSettingChange('logSecurityEvents', checked)}
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="password" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Change Password</CardTitle>
              <CardDescription>
                Update your wallet password for enhanced security
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Current Password</Label>
                <div className="relative">
                  <Input
                    type={showPasswords.current ? 'text' : 'password'}
                    value={passwords.current}
                    onChange={(e) => setPasswords(prev => ({ ...prev, current: e.target.value }))}
                    placeholder="Enter current password"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-2 top-1/2 transform -translate-y-1/2"
                    onClick={() => setShowPasswords(prev => ({ ...prev, current: !prev.current }))}
                  >
                    {showPasswords.current ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </Button>
                </div>
              </div>

              <div className="space-y-2">
                <Label>New Password</Label>
                <div className="relative">
                  <Input
                    type={showPasswords.new ? 'text' : 'password'}
                    value={passwords.new}
                    onChange={(e) => setPasswords(prev => ({ ...prev, new: e.target.value }))}
                    placeholder="Enter new password"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-2 top-1/2 transform -translate-y-1/2"
                    onClick={() => setShowPasswords(prev => ({ ...prev, new: !prev.new }))}
                  >
                    {showPasswords.new ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </Button>
                </div>
                
                {passwords.new && (
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-sm">
                      <span>Password strength:</span>
                      <span className={`font-medium ${passwordStrength.score <= 2 ? 'text-red-500' : passwordStrength.score <= 4 ? 'text-yellow-500' : 'text-green-500'}`}>
                        {getPasswordStrengthText(passwordStrength.score)}
                      </span>
                    </div>
                    <Progress 
                      value={(passwordStrength.score / 6) * 100} 
                      className="h-2"
                    />
                    {passwordStrength.feedback.length > 0 && (
                      <ul className="text-xs text-muted-foreground space-y-1">
                        {passwordStrength.feedback.map((item, index) => (
                          <li key={index}>• {item}</li>
                        ))}
                      </ul>
                    )}
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <Label>Confirm New Password</Label>
                <div className="relative">
                  <Input
                    type={showPasswords.confirm ? 'text' : 'password'}
                    value={passwords.confirm}
                    onChange={(e) => setPasswords(prev => ({ ...prev, confirm: e.target.value }))}
                    placeholder="Confirm new password"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-2 top-1/2 transform -translate-y-1/2"
                    onClick={() => setShowPasswords(prev => ({ ...prev, confirm: !prev.confirm }))}
                  >
                    {showPasswords.confirm ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </Button>
                </div>
              </div>

              {passwords.new && passwords.confirm && passwords.new !== passwords.confirm && (
                <Alert variant="destructive">
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    Passwords do not match
                  </AlertDescription>
                </Alert>
              )}

              <Button 
                onClick={handlePasswordChange}
                disabled={isChangingPassword || !passwordStrength.isValid || passwords.new !== passwords.confirm}
                className="w-full"
              >
                {isChangingPassword ? (
                  <>
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                    Changing Password...
                  </>
                ) : (
                  <>
                    <Key className="w-4 h-4 mr-2" />
                    Change Password
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="advanced" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Advanced Security</CardTitle>
              <CardDescription>
                Advanced security features and experimental options
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="flex items-center gap-2">
                    <Fingerprint className="w-4 h-4" />
                    Biometric authentication
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    Use fingerprint or face recognition (if supported)
                  </p>
                </div>
                <Switch
                  checked={settings.enableBiometric}
                  onCheckedChange={(checked) => handleSettingChange('enableBiometric', checked)}
                  disabled={!('webauthn' in window)}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="flex items-center gap-2">
                    <Smartphone className="w-4 h-4" />
                    Two-factor authentication
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    Require additional verification for sensitive operations
                  </p>
                </div>
                <Switch
                  checked={settings.enableTwoFactor}
                  onCheckedChange={(checked) => handleSettingChange('enableTwoFactor', checked)}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="flex items-center gap-2">
                    <Wifi className="w-4 h-4" />
                    Allow remote connections
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    Allow WalletConnect and other remote wallet connections
                  </p>
                </div>
                <Switch
                  checked={settings.allowRemoteConnections}
                  onCheckedChange={(checked) => handleSettingChange('allowRemoteConnections', checked)}
                />
              </div>

              <div className="space-y-2">
                <Label>Session timeout (hours)</Label>
                <Input
                  type="number"
                  value={settings.sessionTimeout}
                  onChange={(e) => handleSettingChange('sessionTimeout', parseInt(e.target.value))}
                  min="1"
                  max="168"
                />
                <p className="text-xs text-muted-foreground">
                  Automatically log out after this period of inactivity
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="events" className="space-y-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Security Events</CardTitle>
                  <CardDescription>
                    Recent security-related activities and alerts
                  </CardDescription>
                </div>
                <Button onClick={clearSecurityEvents} variant="outline" size="sm">
                  <Trash2 className="w-4 h-4 mr-2" />
                  Clear Events
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {securityEvents.length === 0 ? (
                <div className="text-center py-8">
                  <Shield className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No security events</h3>
                  <p className="text-muted-foreground">
                    Security events will appear here when they occur
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  {securityEvents.map((event, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className={`w-2 h-2 rounded-full ${getEventSeverityColor(event.severity)}`} />
                        <div>
                          <p className="font-medium">{event.type.replace(/_/g, ' ')}</p>
                          <p className="text-sm text-muted-foreground">
                            {new Date(event.timestamp).toLocaleString()}
                          </p>
                        </div>
                      </div>
                      <Badge variant="outline" className="text-xs">
                        {event.severity}
                      </Badge>
                    </motion.div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
