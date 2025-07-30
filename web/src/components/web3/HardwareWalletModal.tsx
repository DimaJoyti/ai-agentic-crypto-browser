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
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
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
  RefreshCw,
  Wallet,
  Monitor,
  Copy,
  Eye,
  EyeOff
} from 'lucide-react'
import { useHardwareWallet } from '@/hooks/useHardwareWallet'
import { type HardwareWalletDevice, type HardwareWalletAccount } from '@/lib/hardware-wallet-manager'
import { toast } from 'sonner'

interface HardwareWalletModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (device: HardwareWalletDevice, account: HardwareWalletAccount) => void
}
export function HardwareWalletModal({ isOpen, onClose, onSuccess }: HardwareWalletModalProps) {
  const [activeTab, setActiveTab] = useState('devices')
  const [showAddresses, setShowAddresses] = useState(false)

  const {
    state,
    scanDevices,
    connectDevice,
    disconnectDevice,
    selectDevice,
    loadAccounts,
    selectAccount,
    clearError
  } = useHardwareWallet({
    autoScan: true,
    onDeviceConnect: (device) => {
      setActiveTab('accounts')
    },
    onAccountSelect: (account) => {
      if (state.selectedDevice) {
        onSuccess?.(state.selectedDevice, account)
        onClose()
      }
    },
    onError: (error) => {
      toast.error('Hardware wallet error', {
        description: error.message
      })
    }
  })

  const handleDeviceConnect = async (device: HardwareWalletDevice) => {
    try {
      await connectDevice(device.id)
    } catch (error) {
      console.error('Failed to connect device:', error)
    }
  }

  const handleAccountSelect = (account: HardwareWalletAccount) => {
    selectAccount(account)
  }

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard')
  }

  const getDeviceIcon = (type: string) => {
    switch (type) {
      case 'ledger':
        return <Shield className="w-8 h-8" />
      case 'trezor':
        return <Lock className="w-8 h-8" />
      case 'gridplus':
        return <Zap className="w-8 h-8" />
      default:
        return <Wallet className="w-8 h-8" />
    }
  }

  const getConnectionMethodIcon = (method: string) => {
    switch (method) {
      case 'usb':
        return <Usb className="w-4 h-4" />
      case 'bluetooth':
        return <Monitor className="w-4 h-4" />
      case 'webusb':
        return <Monitor className="w-4 h-4" />
      default:
        return <Usb className="w-4 h-4" />
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            Hardware Wallet Connection
          </DialogTitle>
          <DialogDescription>
            Connect your hardware wallet for maximum security and control
          </DialogDescription>
        </DialogHeader>

        {state.error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              {state.error}
              <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
                Dismiss
              </Button>
            </AlertDescription>
          </Alert>
        )}

        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="devices">Devices</TabsTrigger>
            <TabsTrigger value="accounts" disabled={!state.selectedDevice}>
              Accounts
            </TabsTrigger>
          </TabsList>

          <TabsContent value="devices" className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-medium">Available Devices</h3>
              <Button
                variant="outline"
                size="sm"
                onClick={scanDevices}
                disabled={state.isScanning}
              >
                <RefreshCw className={`w-4 h-4 mr-2 ${state.isScanning ? 'animate-spin' : ''}`} />
                {state.isScanning ? 'Scanning...' : 'Scan'}
              </Button>
            </div>

            <div className="space-y-3">
              <AnimatePresence>
                {state.devices.map((device, index) => (
                  <motion.div
                    key={device.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                  >
                    <Card className={`cursor-pointer transition-all hover:shadow-md ${
                      state.selectedDevice?.id === device.id ? 'ring-2 ring-primary' : ''
                    }`}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-12 h-12 bg-secondary rounded-lg flex items-center justify-center">
                              {getDeviceIcon(device.type)}
                            </div>
                            <div>
                              <h4 className="font-medium">{device.model}</h4>
                              <p className="text-sm text-muted-foreground capitalize">
                                {device.type} • {device.version}
                              </p>
                              <div className="flex items-center gap-2 mt-1">
                                {getConnectionMethodIcon(device.connectionMethod)}
                                <span className="text-xs text-muted-foreground capitalize">
                                  {device.connectionMethod}
                                </span>
                              </div>
                            </div>
                          </div>
                          <div className="flex flex-col items-end gap-2">
                            <div className="flex items-center gap-2">
                              {device.isConnected ? (
                                <Badge variant="default">Connected</Badge>
                              ) : (
                                <Badge variant="outline">Disconnected</Badge>
                              )}
                              {device.isLocked && (
                                <Badge variant="secondary">
                                  <Lock className="w-3 h-3 mr-1" />
                                  Locked
                                </Badge>
                              )}
                            </div>
                            <Button
                              size="sm"
                              onClick={() => handleDeviceConnect(device)}
                              disabled={state.isConnecting || device.isConnected}
                            >
                              {state.isConnecting && state.selectedDevice?.id === device.id ? (
                                <>
                                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                                  Connecting...
                                </>
                              ) : device.isConnected ? (
                                'Connected'
                              ) : (
                                'Connect'
                              )}
                            </Button>
                          </div>
                        </div>

                        {device.supportedApps.length > 0 && (
                          <div className="mt-3 pt-3 border-t">
                            <p className="text-xs text-muted-foreground mb-2">Supported Apps:</p>
                            <div className="flex flex-wrap gap-1">
                              {device.supportedApps.slice(0, 4).map((app) => (
                                <Badge key={app} variant="secondary" className="text-xs">
                                  {app}
                                </Badge>
                              ))}
                              {device.supportedApps.length > 4 && (
                                <Badge variant="secondary" className="text-xs">
                                  +{device.supportedApps.length - 4} more
                                </Badge>
                              )}
                            </div>
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  </motion.div>
                ))}
              </AnimatePresence>

              {state.devices.length === 0 && !state.isScanning && (
                <div className="text-center py-8">
                  <Wallet className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Hardware Wallets Found</h3>
                  <p className="text-muted-foreground mb-4">
                    Make sure your device is connected and unlocked
                  </p>
                  <Button onClick={scanDevices} disabled={state.isScanning}>
                    <RefreshCw className="w-4 h-4 mr-2" />
                    Scan for Devices
                  </Button>
                </div>
              )}
            </div>
          </TabsContent>

          <TabsContent value="accounts" className="space-y-4">
            {state.selectedDevice && (
              <>
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-medium">Select Account</h3>
                    <p className="text-sm text-muted-foreground">
                      Choose an account from your {state.selectedDevice.model}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setShowAddresses(!showAddresses)}
                    >
                      {showAddresses ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => loadAccounts(state.selectedDevice!.id)}
                    >
                      <RefreshCw className="w-4 h-4 mr-2" />
                      Refresh
                    </Button>
                  </div>
                </div>

                <div className="space-y-3">
                  <AnimatePresence>
                    {state.accounts.map((account, index) => (
                      <motion.div
                        key={account.address}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                      >
                        <Card className={`cursor-pointer transition-all hover:shadow-md ${
                          state.selectedAccount?.address === account.address ? 'ring-2 ring-primary' : ''
                        }`}>
                          <CardContent className="p-4">
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                                  <Wallet className="w-5 h-5" />
                                </div>
                                <div>
                                  <h4 className="font-medium">Account {account.index + 1}</h4>
                                  <p className="text-sm text-muted-foreground font-mono">
                                    {showAddresses ? account.address : formatAddress(account.address)}
                                  </p>
                                  <p className="text-xs text-muted-foreground">
                                    {account.path}
                                  </p>
                                </div>
                              </div>
                              <div className="flex items-center gap-2">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => copyToClipboard(account.address)}
                                >
                                  <Copy className="w-4 h-4" />
                                </Button>
                                <Button
                                  size="sm"
                                  onClick={() => handleAccountSelect(account)}
                                  variant={state.selectedAccount?.address === account.address ? "default" : "outline"}
                                >
                                  {state.selectedAccount?.address === account.address ? 'Selected' : 'Select'}
                                </Button>
                              </div>
                            </div>
                            {account.balance && (
                              <div className="mt-2 pt-2 border-t">
                                <p className="text-sm text-muted-foreground">
                                  Balance: {account.balance} ETH
                                </p>
                              </div>
                            )}
                          </CardContent>
                        </Card>
                      </motion.div>
                    ))}
                  </AnimatePresence>

                  {state.accounts.length === 0 && (
                    <div className="text-center py-8">
                      <Wallet className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                      <h3 className="text-lg font-medium mb-2">No Accounts Found</h3>
                      <p className="text-muted-foreground mb-4">
                        Make sure the Ethereum app is open on your device
                      </p>
                      <Button onClick={() => loadAccounts(state.selectedDevice!.id)}>
                        <RefreshCw className="w-4 h-4 mr-2" />
                        Load Accounts
                      </Button>
                    </div>
                  )}
                </div>
              </>
            )}
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
