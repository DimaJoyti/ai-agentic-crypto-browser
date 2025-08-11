'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import {
  Shield,
  Smartphone,
  Key,
  QrCode,
  CheckCircle,
  AlertTriangle,
  Clock,
  Copy,
  Download,
  RefreshCw,
  Mail,
  MessageSquare,
  Fingerprint,
  Eye,
  EyeOff,
  Settings,
  Plus,
  Trash2,
  Monitor
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface MFAMethod {
  id: string
  type: 'totp' | 'sms' | 'email' | 'hardware' | 'biometric' | 'backup_codes'
  name: string
  description: string
  isEnabled: boolean
  isPrimary: boolean
  setupAt?: number
  lastUsed?: number
  status: 'active' | 'pending' | 'disabled' | 'expired'
  metadata?: {
    phoneNumber?: string
    email?: string
    deviceName?: string
    appName?: string
  }
}

interface BackupCode {
  id: string
  code: string
  isUsed: boolean
  usedAt?: number
}

interface SecuritySession {
  id: string
  deviceName: string
  location: string
  ipAddress: string
  userAgent: string
  createdAt: number
  lastActivity: number
  isCurrent: boolean
  mfaVerified: boolean
}

export function MultiFactorAuth() {
  const [mfaMethods, setMfaMethods] = useState<MFAMethod[]>([])
  const [backupCodes, setBackupCodes] = useState<BackupCode[]>([])
  const [sessions, setSessions] = useState<SecuritySession[]>([])
  const [showSetupTotp, setShowSetupTotp] = useState(false)
  const [showBackupCodes, setShowBackupCodes] = useState(false)
  const [totpSecret, setTotpSecret] = useState('')
  const [verificationCode, setVerificationCode] = useState('')
  const [setupStep, setSetupStep] = useState(1)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock MFA data
    const mockMfaMethods: MFAMethod[] = [
      {
        id: 'totp1',
        type: 'totp',
        name: 'Authenticator App',
        description: 'Google Authenticator, Authy, or similar TOTP app',
        isEnabled: true,
        isPrimary: true,
        setupAt: Date.now() - 86400000 * 30,
        lastUsed: Date.now() - 3600000,
        status: 'active',
        metadata: { appName: 'Google Authenticator' }
      },
      {
        id: 'sms1',
        type: 'sms',
        name: 'SMS Verification',
        description: 'Text message to your phone number',
        isEnabled: true,
        isPrimary: false,
        setupAt: Date.now() - 86400000 * 60,
        lastUsed: Date.now() - 86400000 * 2,
        status: 'active',
        metadata: { phoneNumber: '+1 (555) ***-1234' }
      },
      {
        id: 'email1',
        type: 'email',
        name: 'Email Verification',
        description: 'Verification code sent to your email',
        isEnabled: false,
        isPrimary: false,
        status: 'disabled',
        metadata: { email: 'user@*****.com' }
      },
      {
        id: 'hardware1',
        type: 'hardware',
        name: 'Hardware Security Key',
        description: 'YubiKey or other FIDO2 compatible device',
        isEnabled: false,
        isPrimary: false,
        status: 'pending'
      }
    ]

    const mockBackupCodes: BackupCode[] = [
      { id: 'bc1', code: '1234-5678', isUsed: false },
      { id: 'bc2', code: '2345-6789', isUsed: true, usedAt: Date.now() - 86400000 * 5 },
      { id: 'bc3', code: '3456-7890', isUsed: false },
      { id: 'bc4', code: '4567-8901', isUsed: false },
      { id: 'bc5', code: '5678-9012', isUsed: false },
      { id: 'bc6', code: '6789-0123', isUsed: false },
      { id: 'bc7', code: '7890-1234', isUsed: false },
      { id: 'bc8', code: '8901-2345', isUsed: false }
    ]

    const mockSessions: SecuritySession[] = [
      {
        id: 'session1',
        deviceName: 'Chrome on Windows',
        location: 'New York, US',
        ipAddress: '192.168.1.100',
        userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        createdAt: Date.now() - 3600000,
        lastActivity: Date.now() - 300000,
        isCurrent: true,
        mfaVerified: true
      },
      {
        id: 'session2',
        deviceName: 'Safari on iPhone',
        location: 'New York, US',
        ipAddress: '192.168.1.101',
        userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)',
        createdAt: Date.now() - 86400000,
        lastActivity: Date.now() - 7200000,
        isCurrent: false,
        mfaVerified: true
      },
      {
        id: 'session3',
        deviceName: 'Chrome on Android',
        location: 'London, UK',
        ipAddress: '203.0.113.1',
        userAgent: 'Mozilla/5.0 (Linux; Android 11; SM-G991B)',
        createdAt: Date.now() - 86400000 * 3,
        lastActivity: Date.now() - 86400000 * 2,
        isCurrent: false,
        mfaVerified: false
      }
    ]

    setMfaMethods(mockMfaMethods)
    setBackupCodes(mockBackupCodes)
    setSessions(mockSessions)
    setTotpSecret('JBSWY3DPEHPK3PXP') // Mock TOTP secret
  }, [isConnected])

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-500'
      case 'pending': return 'text-yellow-500'
      case 'disabled': case 'expired': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'active': return 'default'
      case 'pending': return 'secondary'
      case 'disabled': case 'expired': return 'destructive'
      default: return 'outline'
    }
  }

  const getMethodIcon = (type: string) => {
    switch (type) {
      case 'totp': return <Smartphone className="w-4 h-4" />
      case 'sms': return <MessageSquare className="w-4 h-4" />
      case 'email': return <Mail className="w-4 h-4" />
      case 'hardware': return <Key className="w-4 h-4" />
      case 'biometric': return <Fingerprint className="w-4 h-4" />
      case 'backup_codes': return <Shield className="w-4 h-4" />
      default: return <Shield className="w-4 h-4" />
    }
  }

  const toggleMfaMethod = (methodId: string) => {
    setMfaMethods(prev => prev.map(method => 
      method.id === methodId 
        ? { ...method, isEnabled: !method.isEnabled, status: method.isEnabled ? 'disabled' : 'active' }
        : method
    ))
  }

  const setPrimaryMethod = (methodId: string) => {
    setMfaMethods(prev => prev.map(method => ({
      ...method,
      isPrimary: method.id === methodId
    })))
  }

  const generateBackupCodes = () => {
    const newCodes: BackupCode[] = Array.from({ length: 8 }, (_, i) => ({
      id: `bc_new_${i}`,
      code: `${Math.floor(1000 + Math.random() * 9000)}-${Math.floor(1000 + Math.random() * 9000)}`,
      isUsed: false
    }))
    setBackupCodes(newCodes)
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    // In real implementation, show toast notification
  }

  const terminateSession = (sessionId: string) => {
    setSessions(prev => prev.filter(session => session.id !== sessionId))
  }

  const getMfaStrength = () => {
    const enabledMethods = mfaMethods.filter(m => m.isEnabled).length
    const hasHardware = mfaMethods.some(m => m.type === 'hardware' && m.isEnabled)
    const hasBackup = backupCodes.some(c => !c.isUsed)
    
    let strength = enabledMethods * 25
    if (hasHardware) strength += 15
    if (hasBackup) strength += 10
    
    return Math.min(strength, 100)
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access multi-factor authentication settings
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
          <h2 className="text-2xl font-bold">Multi-Factor Authentication</h2>
          <p className="text-muted-foreground">
            Secure your account with multiple layers of protection
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            MFA Strength: {getMfaStrength()}%
          </Badge>
        </div>
      </div>

      {/* MFA Strength Indicator */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            Security Strength
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Overall MFA Strength</span>
              <span className="font-bold">{getMfaStrength()}%</span>
            </div>
            <Progress value={getMfaStrength()} className="h-3" />
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
              <div className="flex items-center gap-2">
                <CheckCircle className={cn(
                  "w-4 h-4",
                  mfaMethods.filter(m => m.isEnabled).length >= 2 ? "text-green-500" : "text-muted-foreground"
                )} />
                <span>Multiple Methods</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle className={cn(
                  "w-4 h-4",
                  mfaMethods.some(m => m.type === 'hardware' && m.isEnabled) ? "text-green-500" : "text-muted-foreground"
                )} />
                <span>Hardware Security</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle className={cn(
                  "w-4 h-4",
                  backupCodes.some(c => !c.isUsed) ? "text-green-500" : "text-muted-foreground"
                )} />
                <span>Backup Codes</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* MFA Methods */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Key className="w-5 h-5" />
              Authentication Methods
            </CardTitle>
            <Button onClick={() => setShowSetupTotp(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Add Method
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {mfaMethods.map((method) => (
              <div key={method.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                      {getMethodIcon(method.type)}
                    </div>
                    <div>
                      <h4 className="font-bold">{method.name}</h4>
                      <p className="text-sm text-muted-foreground">{method.description}</p>
                      {method.metadata && (
                        <p className="text-xs text-muted-foreground mt-1">
                          {method.metadata.phoneNumber || method.metadata.email || method.metadata.appName || method.metadata.deviceName}
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {method.isPrimary && (
                      <Badge variant="outline">Primary</Badge>
                    )}
                    <Badge variant={getStatusBadgeVariant(method.status)}>
                      {method.status}
                    </Badge>
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="text-sm text-muted-foreground">
                    {method.setupAt && `Added: ${formatTime(method.setupAt)}`}
                    {method.lastUsed && ` • Last used: ${formatTime(method.lastUsed)}`}
                  </div>
                  <div className="flex gap-2">
                    {method.isEnabled && !method.isPrimary && (
                      <Button 
                        variant="outline" 
                        size="sm"
                        onClick={() => setPrimaryMethod(method.id)}
                      >
                        Set Primary
                      </Button>
                    )}
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => toggleMfaMethod(method.id)}
                    >
                      {method.isEnabled ? 'Disable' : 'Enable'}
                    </Button>
                    {method.status === 'disabled' && (
                      <Button variant="outline" size="sm">
                        <Trash2 className="w-3 h-3 mr-1" />
                        Remove
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Backup Codes */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Key className="w-5 h-5" />
              Backup Recovery Codes
            </CardTitle>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setShowBackupCodes(!showBackupCodes)}>
                {showBackupCodes ? <EyeOff className="w-4 h-4 mr-2" /> : <Eye className="w-4 h-4 mr-2" />}
                {showBackupCodes ? 'Hide' : 'Show'} Codes
              </Button>
              <Button onClick={generateBackupCodes}>
                <RefreshCw className="w-4 h-4 mr-2" />
                Generate New
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Alert className="mb-4">
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>
              Store these backup codes in a safe place. Each code can only be used once to access your account if you lose access to your primary MFA methods.
            </AlertDescription>
          </Alert>

          {showBackupCodes && (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
              {backupCodes.map((code) => (
                <div
                  key={code.id}
                  className={cn(
                    "p-3 border rounded-lg font-mono text-center cursor-pointer transition-colors",
                    code.isUsed 
                      ? "bg-muted text-muted-foreground line-through" 
                      : "hover:bg-muted/50"
                  )}
                  onClick={() => !code.isUsed && copyToClipboard(code.code)}
                >
                  {code.code}
                  {code.isUsed && (
                    <div className="text-xs mt-1">
                      Used {code.usedAt && formatTime(code.usedAt)}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}

          <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
            <span>
              {backupCodes.filter(c => !c.isUsed).length} of {backupCodes.length} codes remaining
            </span>
            <Button variant="outline" size="sm" onClick={() => copyToClipboard(backupCodes.filter(c => !c.isUsed).map(c => c.code).join('\n'))}>
              <Copy className="w-3 h-3 mr-1" />
              Copy All
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Active Sessions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Monitor className="w-5 h-5" />
            Active Sessions
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {sessions.map((session) => (
              <div key={session.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                      <Monitor className="w-5 h-5" />
                    </div>
                    <div>
                      <h4 className="font-bold">{session.deviceName}</h4>
                      <p className="text-sm text-muted-foreground">
                        {session.location} • {session.ipAddress}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {session.isCurrent && (
                      <Badge variant="outline">Current Session</Badge>
                    )}
                    {session.mfaVerified ? (
                      <Badge variant="default">
                        <CheckCircle className="w-3 h-3 mr-1" />
                        MFA Verified
                      </Badge>
                    ) : (
                      <Badge variant="destructive">
                        <AlertTriangle className="w-3 h-3 mr-1" />
                        MFA Required
                      </Badge>
                    )}
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <div className="text-sm text-muted-foreground">
                    Created: {formatTime(session.createdAt)} • Last activity: {formatTime(session.lastActivity)}
                  </div>
                  {!session.isCurrent && (
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => terminateSession(session.id)}
                    >
                      Terminate
                    </Button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
