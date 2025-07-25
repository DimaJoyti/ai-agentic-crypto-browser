'use client'

import { useState, useEffect } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Shield, 
  Smartphone, 
  Mail, 
  Key, 
  Copy, 
  Check, 
  AlertTriangle,
  Download,
  Eye,
  EyeOff
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'

interface MFASetupData {
  secret: string
  qr_code_url: string
  backup_codes: string[]
  method: 'totp' | 'sms' | 'email'
}

interface MFASetupProps {
  onComplete: () => void
  onCancel: () => void
}

export default function MFASetup({ onComplete, onCancel }: MFASetupProps) {
  const { user } = useAuth()
  const [selectedMethod, setSelectedMethod] = useState<'totp' | 'sms' | 'email'>('totp')
  const [setupData, setSetupData] = useState<MFASetupData | null>(null)
  const [verificationCode, setVerificationCode] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [step, setStep] = useState<'method' | 'setup' | 'verify' | 'backup'>('method')
  const [copiedSecret, setCopiedSecret] = useState(false)
  const [copiedBackup, setCopiedBackup] = useState<number | null>(null)
  const [showBackupCodes, setShowBackupCodes] = useState(false)

  const setupMFA = async (method: 'totp' | 'sms' | 'email') => {
    setIsLoading(true)
    setError('')

    try {
      const response = await fetch('/api/auth/mfa/setup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
        body: JSON.stringify({ method }),
      })

      if (response.ok) {
        const data = await response.json()
        setSetupData(data)
        setStep('setup')
      } else {
        const errorData = await response.json()
        setError(errorData.error || 'Failed to setup MFA')
      }
    } catch (err) {
      setError('Network error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const verifyMFA = async () => {
    if (!verificationCode.trim()) {
      setError('Please enter the verification code')
      return
    }

    setIsLoading(true)
    setError('')

    try {
      const response = await fetch('/api/auth/mfa/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
        body: JSON.stringify({
          code: verificationCode,
          secret: setupData?.secret,
        }),
      })

      if (response.ok) {
        const data = await response.json()
        if (data.verified) {
          setStep('backup')
        } else {
          setError('Invalid verification code')
        }
      } else {
        const errorData = await response.json()
        setError(errorData.error || 'Verification failed')
      }
    } catch (err) {
      setError('Network error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const copyToClipboard = async (text: string, type: 'secret' | number) => {
    try {
      await navigator.clipboard.writeText(text)
      if (type === 'secret') {
        setCopiedSecret(true)
        setTimeout(() => setCopiedSecret(false), 2000)
      } else {
        setCopiedBackup(type)
        setTimeout(() => setCopiedBackup(null), 2000)
      }
    } catch (err) {
      console.error('Failed to copy to clipboard:', err)
    }
  }

  const downloadBackupCodes = () => {
    if (!setupData?.backup_codes) return

    const content = `AI Agentic Browser - MFA Backup Codes
Generated: ${new Date().toLocaleString()}
User: ${user?.email}

IMPORTANT: Store these codes in a safe place. Each code can only be used once.

${setupData.backup_codes.map((code, index) => `${index + 1}. ${code}`).join('\n')}

Instructions:
- Use these codes if you lose access to your authenticator device
- Each code can only be used once
- Generate new codes if you use all of them
- Keep these codes secure and private`

    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `mfa-backup-codes-${Date.now()}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  const renderMethodSelection = () => (
    <div className="space-y-4">
      <div className="text-center mb-6">
        <Shield className="h-12 w-12 text-primary mx-auto mb-4" />
        <h2 className="text-2xl font-bold">Enable Two-Factor Authentication</h2>
        <p className="text-muted-foreground">
          Add an extra layer of security to your account
        </p>
      </div>

      <Tabs value={selectedMethod} onValueChange={(value) => setSelectedMethod(value as any)}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="totp" className="flex items-center gap-2">
            <Smartphone className="h-4 w-4" />
            Authenticator App
          </TabsTrigger>
          <TabsTrigger value="sms" className="flex items-center gap-2">
            <Smartphone className="h-4 w-4" />
            SMS
          </TabsTrigger>
          <TabsTrigger value="email" className="flex items-center gap-2">
            <Mail className="h-4 w-4" />
            Email
          </TabsTrigger>
        </TabsList>

        <TabsContent value="totp" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Smartphone className="h-5 w-5" />
                Authenticator App
              </CardTitle>
              <CardDescription>
                Use an authenticator app like Google Authenticator, Authy, or 1Password
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-green-500" />
                  Most secure option
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-green-500" />
                  Works offline
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-green-500" />
                  No dependency on phone/email service
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="sms" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Smartphone className="h-5 w-5" />
                SMS Verification
              </CardTitle>
              <CardDescription>
                Receive verification codes via text message
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-green-500" />
                  Easy to use
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <AlertTriangle className="h-4 w-4 text-yellow-500" />
                  Requires phone service
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <AlertTriangle className="h-4 w-4 text-yellow-500" />
                  Less secure than authenticator apps
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="email" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Mail className="h-5 w-5" />
                Email Verification
              </CardTitle>
              <CardDescription>
                Receive verification codes via email
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-green-500" />
                  Always accessible
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <AlertTriangle className="h-4 w-4 text-yellow-500" />
                  Requires email access
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <AlertTriangle className="h-4 w-4 text-yellow-500" />
                  Less secure than authenticator apps
                </div>
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

      <div className="flex gap-3">
        <Button 
          onClick={() => setupMFA(selectedMethod)} 
          disabled={isLoading}
          className="flex-1"
        >
          {isLoading ? 'Setting up...' : 'Continue'}
        </Button>
        <Button variant="outline" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </div>
  )

  const renderTOTPSetup = () => (
    <div className="space-y-6">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-2">Scan QR Code</h2>
        <p className="text-muted-foreground">
          Use your authenticator app to scan this QR code
        </p>
      </div>

      <div className="flex justify-center">
        <div className="bg-white p-4 rounded-lg">
          <QRCodeSVG value={setupData?.qr_code_url || ''} size={200} />
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm">Manual Entry</CardTitle>
          <CardDescription>
            If you can't scan the QR code, enter this secret manually
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2">
            <Input 
              value={setupData?.secret || ''} 
              readOnly 
              className="font-mono text-sm"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => copyToClipboard(setupData?.secret || '', 'secret')}
            >
              {copiedSecret ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
            </Button>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-3">
        <Label htmlFor="verification-code">Enter verification code from your app</Label>
        <Input
          id="verification-code"
          value={verificationCode}
          onChange={(e) => setVerificationCode(e.target.value)}
          placeholder="000000"
          maxLength={6}
          className="text-center text-lg tracking-widest"
        />
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="flex gap-3">
        <Button onClick={verifyMFA} disabled={isLoading || !verificationCode} className="flex-1">
          {isLoading ? 'Verifying...' : 'Verify & Enable'}
        </Button>
        <Button variant="outline" onClick={() => setStep('method')}>
          Back
        </Button>
      </div>
    </div>
  )

  const renderBackupCodes = () => (
    <div className="space-y-6">
      <div className="text-center">
        <Key className="h-12 w-12 text-primary mx-auto mb-4" />
        <h2 className="text-2xl font-bold mb-2">Save Your Backup Codes</h2>
        <p className="text-muted-foreground">
          Store these codes safely. You'll need them if you lose access to your authenticator.
        </p>
      </div>

      <Alert>
        <AlertTriangle className="h-4 w-4" />
        <AlertDescription>
          <strong>Important:</strong> Each backup code can only be used once. 
          Store them in a secure location like a password manager.
        </AlertDescription>
      </Alert>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-sm">Backup Codes</CardTitle>
            <CardDescription>
              Use these codes if you lose access to your authenticator
            </CardDescription>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowBackupCodes(!showBackupCodes)}
            >
              {showBackupCodes ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={downloadBackupCodes}
            >
              <Download className="h-4 w-4" />
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {showBackupCodes ? (
            <div className="grid grid-cols-2 gap-2">
              {setupData?.backup_codes.map((code, index) => (
                <div key={index} className="flex items-center gap-2 p-2 bg-muted rounded">
                  <span className="font-mono text-sm flex-1">{code}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => copyToClipboard(code, index)}
                  >
                    {copiedBackup === index ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
                  </Button>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <Eye className="h-8 w-8 mx-auto mb-2" />
              <p>Click the eye icon to reveal backup codes</p>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="flex gap-3">
        <Button onClick={onComplete} className="flex-1">
          Complete Setup
        </Button>
      </div>
    </div>
  )

  return (
    <div className="max-w-md mx-auto">
      {step === 'method' && renderMethodSelection()}
      {step === 'setup' && selectedMethod === 'totp' && renderTOTPSetup()}
      {step === 'backup' && renderBackupCodes()}
    </div>
  )
}
