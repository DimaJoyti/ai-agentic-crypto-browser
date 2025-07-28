'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Download, 
  Upload, 
  Shield, 
  Key, 
  AlertTriangle,
  CheckCircle,
  Copy,
  Eye,
  EyeOff,
  FileText,
  QrCode,
  Smartphone
} from 'lucide-react'
import { WalletSecurity, SecureStorage } from '@/lib/security'
import { toast } from 'sonner'

interface BackupData {
  wallets: any[]
  settings: any
  timestamp: string
  version: string
}

export function WalletBackup() {
  const [activeTab, setActiveTab] = useState('backup')
  const [backupPassword, setBackupPassword] = useState('')
  const [restorePassword, setRestorePassword] = useState('')
  const [restoreData, setRestoreData] = useState('')
  const [showMnemonic, setShowMnemonic] = useState(false)
  const [isCreatingBackup, setIsCreatingBackup] = useState(false)
  const [isRestoring, setIsRestoring] = useState(false)

  // Mock mnemonic phrase for demonstration
  const mockMnemonic = [
    'abandon', 'ability', 'able', 'about', 'above', 'absent',
    'absorb', 'abstract', 'absurd', 'abuse', 'access', 'accident'
  ]

  const createBackup = async () => {
    if (!backupPassword) {
      toast.error('Please enter a backup password')
      return
    }

    const passwordValidation = WalletSecurity.validatePassword(backupPassword)
    if (!passwordValidation.isValid) {
      toast.error('Password does not meet security requirements')
      return
    }

    setIsCreatingBackup(true)

    try {
      // Create backup data
      const backupData: BackupData = {
        wallets: [
          {
            id: '1',
            name: 'Main Wallet',
            address: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1',
            type: 'injected',
            chainId: 1
          }
        ],
        settings: {
          autoLock: true,
          encryptLocalStorage: true
        },
        timestamp: new Date().toISOString(),
        version: '1.0.0'
      }

      // Encrypt backup data
      await SecureStorage.store('wallet_backup', backupData, backupPassword, 24 * 365) // 1 year

      // Create downloadable backup file
      const encryptedBackup = JSON.stringify({
        data: await WalletSecurity.encrypt(JSON.stringify(backupData), await WalletSecurity.deriveKey(backupPassword, WalletSecurity.generateSalt())),
        timestamp: backupData.timestamp,
        version: backupData.version
      })

      const blob = new Blob([encryptedBackup], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `wallet-backup-${Date.now()}.json`
      a.click()
      URL.revokeObjectURL(url)

      toast.success('Backup created and downloaded successfully')
    } catch (error) {
      toast.error('Failed to create backup')
    } finally {
      setIsCreatingBackup(false)
    }
  }

  const restoreFromBackup = async () => {
    if (!restorePassword || !restoreData) {
      toast.error('Please provide both password and backup data')
      return
    }

    setIsRestoring(true)

    try {
      const backupFile = JSON.parse(restoreData)
      const decryptedData = await WalletSecurity.decrypt(
        backupFile.data, 
        await WalletSecurity.deriveKey(restorePassword, WalletSecurity.generateSalt())
      )
      
      const backupData: BackupData = JSON.parse(decryptedData)
      
      // Validate backup data structure
      if (!backupData.wallets || !backupData.settings || !backupData.timestamp) {
        throw new Error('Invalid backup data structure')
      }

      // Restore wallets and settings
      localStorage.setItem('restored_wallets', JSON.stringify(backupData.wallets))
      localStorage.setItem('security_settings', JSON.stringify(backupData.settings))

      toast.success('Wallet backup restored successfully')
      setRestoreData('')
      setRestorePassword('')
    } catch (error) {
      toast.error('Failed to restore backup. Please check your password and backup data.')
    } finally {
      setIsRestoring(false)
    }
  }

  const copyMnemonic = () => {
    navigator.clipboard.writeText(mockMnemonic.join(' '))
    toast.success('Mnemonic phrase copied to clipboard')
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        setRestoreData(e.target?.result as string)
      }
      reader.readAsText(file)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <Shield className="w-6 h-6" />
          Wallet Backup & Recovery
        </h2>
        <p className="text-muted-foreground">
          Secure your wallet with encrypted backups and recovery options
        </p>
      </div>

      {/* Security Warning */}
      <Alert>
        <AlertTriangle className="h-4 w-4" />
        <AlertDescription>
          <strong>Important:</strong> Keep your backup files and recovery phrases secure. 
          Anyone with access to them can control your wallet. Never share them online or store them in unsecured locations.
        </AlertDescription>
      </Alert>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="backup">Create Backup</TabsTrigger>
          <TabsTrigger value="restore">Restore Wallet</TabsTrigger>
          <TabsTrigger value="mnemonic">Recovery Phrase</TabsTrigger>
        </TabsList>

        <TabsContent value="backup" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Download className="w-5 h-5" />
                Create Encrypted Backup
              </CardTitle>
              <CardDescription>
                Create a secure, encrypted backup of your wallet data
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Backup Password</Label>
                <Input
                  type="password"
                  value={backupPassword}
                  onChange={(e) => setBackupPassword(e.target.value)}
                  placeholder="Enter a strong password for your backup"
                />
                <p className="text-xs text-muted-foreground">
                  This password will be required to restore your backup. Make sure it's strong and memorable.
                </p>
              </div>

              <div className="bg-muted p-4 rounded-lg">
                <h4 className="font-medium mb-2">What will be backed up:</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• Connected wallet addresses</li>
                  <li>• Security settings and preferences</li>
                  <li>• Transaction history metadata</li>
                  <li>• Custom wallet names and favorites</li>
                </ul>
                <p className="text-xs mt-2 text-muted-foreground">
                  <strong>Note:</strong> Private keys are never included in backups for security reasons.
                </p>
              </div>

              <Button 
                onClick={createBackup}
                disabled={isCreatingBackup || !backupPassword}
                className="w-full"
              >
                {isCreatingBackup ? (
                  <>
                    <Download className="w-4 h-4 mr-2 animate-pulse" />
                    Creating Backup...
                  </>
                ) : (
                  <>
                    <Download className="w-4 h-4 mr-2" />
                    Create & Download Backup
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="restore" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Upload className="w-5 h-5" />
                Restore from Backup
              </CardTitle>
              <CardDescription>
                Restore your wallet from an encrypted backup file
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Backup File</Label>
                <div className="flex gap-2">
                  <Input
                    type="file"
                    accept=".json"
                    onChange={handleFileUpload}
                    className="flex-1"
                  />
                  <Button variant="outline" size="sm">
                    <Upload className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              <div className="space-y-2">
                <Label>Backup Data (JSON)</Label>
                <Textarea
                  value={restoreData}
                  onChange={(e) => setRestoreData(e.target.value)}
                  placeholder="Paste your backup data here or upload a file above"
                  rows={6}
                />
              </div>

              <div className="space-y-2">
                <Label>Backup Password</Label>
                <Input
                  type="password"
                  value={restorePassword}
                  onChange={(e) => setRestorePassword(e.target.value)}
                  placeholder="Enter the password used to create this backup"
                />
              </div>

              <Alert>
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription>
                  Restoring from backup will overwrite your current wallet settings. 
                  Make sure to backup your current data first if needed.
                </AlertDescription>
              </Alert>

              <Button 
                onClick={restoreFromBackup}
                disabled={isRestoring || !restorePassword || !restoreData}
                className="w-full"
              >
                {isRestoring ? (
                  <>
                    <Upload className="w-4 h-4 mr-2 animate-pulse" />
                    Restoring...
                  </>
                ) : (
                  <>
                    <Upload className="w-4 h-4 mr-2" />
                    Restore Wallet
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="mnemonic" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Key className="w-5 h-5" />
                Recovery Phrase
              </CardTitle>
              <CardDescription>
                Your 12-word recovery phrase for wallet restoration
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <Alert variant="destructive">
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription>
                  <strong>Critical Security Warning:</strong> Your recovery phrase gives complete access to your wallet. 
                  Never share it with anyone or store it digitally. Write it down and keep it in a secure location.
                </AlertDescription>
              </Alert>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <Label>Recovery Phrase</Label>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setShowMnemonic(!showMnemonic)}
                    >
                      {showMnemonic ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={copyMnemonic}
                      disabled={!showMnemonic}
                    >
                      <Copy className="w-4 h-4" />
                    </Button>
                  </div>
                </div>

                <div className="grid grid-cols-3 gap-2 p-4 bg-muted rounded-lg">
                  {mockMnemonic.map((word, index) => (
                    <div
                      key={index}
                      className="flex items-center gap-2 p-2 bg-background rounded border"
                    >
                      <span className="text-xs text-muted-foreground w-6">
                        {index + 1}.
                      </span>
                      <span className="font-mono">
                        {showMnemonic ? word : '••••••'}
                      </span>
                    </div>
                  ))}
                </div>

                <div className="bg-blue-50 dark:bg-blue-950 p-4 rounded-lg">
                  <h4 className="font-medium mb-2 flex items-center gap-2">
                    <FileText className="w-4 h-4" />
                    Best Practices for Recovery Phrase Storage:
                  </h4>
                  <ul className="text-sm space-y-1 text-muted-foreground">
                    <li>• Write it down on paper and store in a secure location</li>
                    <li>• Consider using a metal backup for fire/water resistance</li>
                    <li>• Store copies in multiple secure locations</li>
                    <li>• Never store it digitally or take screenshots</li>
                    <li>• Never share it with anyone, including support staff</li>
                    <li>• Test your backup by restoring a test wallet</li>
                  </ul>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
