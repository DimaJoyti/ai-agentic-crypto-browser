'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle 
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { 
  Shield, 
  Usb, 
  CheckCircle, 
  AlertCircle,
  Loader2,
  Zap,
  Fingerprint,
  Lock,
  Unlock,
  RefreshCw
} from 'lucide-react'
import { toast } from 'sonner'

interface HardwareWalletModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (deviceInfo: HardwareDeviceInfo) => void
}

interface HardwareDeviceInfo {
  deviceId: string
  model: string
  manufacturer: string
  version: string
  isLocked: boolean
  appName: string
  appVersion: string
  addresses?: string[]
}

interface HardwareWalletType {
  id: string
  name: string
  manufacturer: string
  icon: React.ReactNode
  description: string
  features: string[]
  isSupported: boolean
  connectionType: 'usb' | 'bluetooth' | 'wifi'
}

const hardwareWallets: HardwareWalletType[] = [
  {
    id: 'ledger',
    name: 'Ledger',
    manufacturer: 'Ledger SAS',
    icon: <Shield className="w-8 h-8" />,
    description: 'Nano S Plus, Nano X, Stax',
    features: ['USB', 'Bluetooth', 'Multi-Chain', 'Secure Element'],
    isSupported: true,
    connectionType: 'usb'
  },
  {
    id: 'trezor',
    name: 'Trezor',
    manufacturer: 'SatoshiLabs',
    icon: <Lock className="w-8 h-8" />,
    description: 'Model T, Model One',
    features: ['USB', 'Touchscreen', 'Open Source', 'PIN Protection'],
    isSupported: true,
    connectionType: 'usb'
  },
  {
    id: 'gridplus',
    name: 'GridPlus',
    manufacturer: 'GridPlus Inc.',
    icon: <Zap className="w-8 h-8" />,
    description: 'Lattice1',
    features: ['WiFi', 'Large Screen', 'Card Reader', 'Always Online'],
    isSupported: true,
    connectionType: 'wifi'
  }
]

const connectionSteps = [
  'Detecting device...',
  'Establishing connection...',
  'Verifying device...',
  'Loading applications...',
  'Ready to use!'
]

export function HardwareWalletModal({ isOpen, onClose, onSuccess }: HardwareWalletModalProps) {
  const [selectedWallet, setSelectedWallet] = useState<string | null>(null)
  const [isConnecting, setIsConnecting] = useState(false)
  const [connectionStep, setConnectionStep] = useState(0)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [deviceInfo, setDeviceInfo] = useState<HardwareDeviceInfo | null>(null)
  const [isDeviceReady, setIsDeviceReady] = useState(false)

  const handleWalletSelect = async (walletId: string) => {
    setSelectedWallet(walletId)
    setIsConnecting(true)
    setConnectionError(null)
    setConnectionStep(0)

    try {
      // Simulate hardware wallet connection process
      for (let i = 0; i < connectionSteps.length; i++) {
        setConnectionStep(i)
        await new Promise(resolve => setTimeout(resolve, 1000))
      }

      // Simulate device info retrieval
      const mockDeviceInfo: HardwareDeviceInfo = {
        deviceId: `${walletId}-${Date.now()}`,
        model: hardwareWallets.find(w => w.id === walletId)?.name || 'Unknown',
        manufacturer: hardwareWallets.find(w => w.id === walletId)?.manufacturer || 'Unknown',
        version: '2.1.0',
        isLocked: false,
        appName: 'Ethereum',
        appVersion: '1.9.0',
        addresses: []
      }

      setDeviceInfo(mockDeviceInfo)
      setIsDeviceReady(true)
      setIsConnecting(false)
      
      toast.success(`${mockDeviceInfo.model} connected successfully!`)
      onSuccess?.(mockDeviceInfo)
      
    } catch (error) {
      console.error('Hardware wallet connection error:', error)
      setConnectionError('Failed to connect to hardware wallet. Please check your device and try again.')
      setIsConnecting(false)
    }
  }

  const handleRetry = () => {
    if (selectedWallet) {
      handleWalletSelect(selectedWallet)
    }
  }

  const handleClose = () => {
    setSelectedWallet(null)
    setIsConnecting(false)
    setConnectionStep(0)
    setConnectionError(null)
    setDeviceInfo(null)
    setIsDeviceReady(false)
    onClose()
  }

  const getConnectionIcon = () => {
    if (connectionError) return <AlertCircle className="w-8 h-8 text-red-500" />
    if (isDeviceReady) return <CheckCircle className="w-8 h-8 text-green-500" />
    if (isConnecting) return <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
    return <Usb className="w-8 h-8 text-gray-500" />
  }

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            Connect Hardware Wallet
          </DialogTitle>
          <DialogDescription>
            Connect your hardware wallet for maximum security
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {!selectedWallet ? (
            // Wallet Selection
            <div className="space-y-4">
              <div className="text-sm text-muted-foreground">
                Choose your hardware wallet type:
              </div>
              
              <div className="grid gap-3">
                {hardwareWallets.map((wallet) => (
                  <Card 
                    key={wallet.id}
                    className={`cursor-pointer transition-all hover:shadow-md ${
                      !wallet.isSupported ? 'opacity-50 cursor-not-allowed' : ''
                    }`}
                    onClick={() => wallet.isSupported && handleWalletSelect(wallet.id)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-12 h-12 bg-secondary rounded-lg flex items-center justify-center">
                            {wallet.icon}
                          </div>
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <h3 className="font-medium">{wallet.name}</h3>
                              {!wallet.isSupported && (
                                <Badge variant="outline" className="text-xs">
                                  Coming Soon
                                </Badge>
                              )}
                            </div>
                            <p className="text-sm text-muted-foreground">
                              {wallet.description}
                            </p>
                            <div className="flex gap-1 mt-1">
                              {wallet.features.slice(0, 2).map((feature) => (
                                <Badge key={feature} variant="secondary" className="text-xs">
                                  {feature}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <Badge variant="outline" className="text-xs">
                            {wallet.connectionType.toUpperCase()}
                          </Badge>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          ) : (
            // Connection Process
            <div className="space-y-6">
              <div className="text-center">
                <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                  {getConnectionIcon()}
                </div>
                
                {isConnecting && (
                  <div className="space-y-4">
                    <h3 className="text-lg font-semibold">
                      Connecting to {hardwareWallets.find(w => w.id === selectedWallet)?.name}
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      {connectionSteps[connectionStep]}
                    </p>
                    <Progress value={(connectionStep + 1) / connectionSteps.length * 100} className="w-full" />
                  </div>
                )}

                {connectionError && (
                  <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-red-600">Connection Failed</h3>
                    <Alert variant="destructive">
                      <AlertCircle className="h-4 w-4" />
                      <AlertDescription>{connectionError}</AlertDescription>
                    </Alert>
                    <div className="flex gap-2">
                      <Button onClick={handleRetry} variant="outline">
                        <RefreshCw className="w-4 h-4 mr-2" />
                        Retry
                      </Button>
                      <Button onClick={handleClose} variant="ghost">
                        Cancel
                      </Button>
                    </div>
                  </div>
                )}

                {isDeviceReady && deviceInfo && (
                  <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-green-600">Device Connected!</h3>
                    <Card>
                      <CardContent className="p-4">
                        <div className="space-y-2 text-sm">
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">Model:</span>
                            <span className="font-medium">{deviceInfo.model}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">Version:</span>
                            <span className="font-medium">{deviceInfo.version}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">App:</span>
                            <span className="font-medium">{deviceInfo.appName} {deviceInfo.appVersion}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">Status:</span>
                            <div className="flex items-center gap-1">
                              {deviceInfo.isLocked ? (
                                <>
                                  <Lock className="w-3 h-3 text-red-500" />
                                  <span className="text-red-500">Locked</span>
                                </>
                              ) : (
                                <>
                                  <Unlock className="w-3 h-3 text-green-500" />
                                  <span className="text-green-500">Unlocked</span>
                                </>
                              )}
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                    <Button onClick={handleClose} className="w-full">
                      Continue with {deviceInfo.model}
                    </Button>
                  </div>
                )}
              </div>
            </div>
          )}

          {!selectedWallet && (
            <div className="text-center">
              <p className="text-sm text-muted-foreground">
                Make sure your hardware wallet is connected and unlocked before proceeding.
              </p>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
