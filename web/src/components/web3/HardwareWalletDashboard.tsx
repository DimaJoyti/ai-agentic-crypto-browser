'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Shield,
  Usb,
  Bluetooth,
  Wifi,
  CheckCircle,
  AlertCircle,
  Loader2,
  RefreshCw,
  Zap,
  Lock,
  Unlock,
  Wallet,
  Monitor,
  Smartphone,
  HardDrive,
  Eye,
  EyeOff,
  Copy,
  Settings,
  Activity,
  Clock,
  FileText,
  Download
} from 'lucide-react'
import { useHardwareWallet } from '@/hooks/useHardwareWallet'
import { type HardwareWalletDevice, type HardwareWalletAccount, type SigningRequest } from '@/lib/hardware-wallet-manager'
import { toast } from 'sonner'

export function HardwareWalletDashboard() {
  const [activeTab, setActiveTab] = useState('overview')
  const [showAddresses, setShowAddresses] = useState(false)
  const [selectedDevice, setSelectedDevice] = useState<HardwareWalletDevice | null>(null)

  const {
    state,
    scanDevices,
    connectDevice,
    disconnectDevice,
    loadAccounts,
    selectAccount,
    signMessage,
    clearError,
    reset
  } = useHardwareWallet({
    autoScan: true,
    onDeviceConnect: (device) => {
      toast.success(`${device.model} connected successfully`)
    },
    onDeviceDisconnect: (device) => {
      toast.info(`${device.model} disconnected`)
    },
    onError: (error) => {
      toast.error('Hardware wallet error', {
        description: error.message
      })
    }
  })

  const handleDeviceSelect = (device: HardwareWalletDevice) => {
    setSelectedDevice(device)
    if (!device.isConnected) {
      connectDevice(device.id)
    }
  }

  const handleTestSigning = async () => {
    if (!state.selectedAccount) {
      toast.error('Please select an account first')
      return
    }

    try {
      const message = `Test message signed at ${new Date().toISOString()}`
      await signMessage(message)
      toast.success('Message signed successfully!')
    } catch (error) {
      console.error('Signing failed:', error)
    }
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
        return <Shield className="w-8 h-8 text-blue-500" />
      case 'trezor':
        return <Lock className="w-8 h-8 text-green-500" />
      case 'gridplus':
        return <Zap className="w-8 h-8 text-purple-500" />
      default:
        return <Wallet className="w-8 h-8 text-gray-500" />
    }
  }

  const getConnectionIcon = (method: string) => {
    switch (method) {
      case 'usb':
        return <Usb className="w-4 h-4" />
      case 'bluetooth':
        return <Bluetooth className="w-4 h-4" />
      case 'wifi':
        return <Wifi className="w-4 h-4" />
      default:
        return <Monitor className="w-4 h-4" />
    }
  }

  const getStatusColor = (device: HardwareWalletDevice) => {
    if (device.isConnected) return 'text-green-600'
    if (device.isLocked) return 'text-red-600'
    return 'text-yellow-600'
  }

  const getStatusText = (device: HardwareWalletDevice) => {
    if (device.isConnected) return 'Connected'
    if (device.isLocked) return 'Locked'
    return 'Disconnected'
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Hardware Wallet Dashboard</h2>
          <p className="text-muted-foreground">
            Manage your hardware wallet devices and accounts
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={scanDevices}
            disabled={state.isScanning}
          >
            <RefreshCw className={`w-4 h-4 mr-2 ${state.isScanning ? 'animate-spin' : ''}`} />
            {state.isScanning ? 'Scanning...' : 'Scan Devices'}
          </Button>
          <Button variant="outline" size="sm" onClick={reset}>
            Reset
          </Button>
        </div>
      </div>

      {/* Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Connected Devices</p>
                <p className="text-2xl font-bold">{state.connectedDevices.length}</p>
              </div>
              <Shield className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {state.devices.length} total devices
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Available Accounts</p>
                <p className="text-2xl font-bold">{state.accounts.length}</p>
              </div>
              <Wallet className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {state.selectedAccount ? 'Account selected' : 'No account selected'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Signing Requests</p>
                <p className="text-2xl font-bold">{state.signingRequests.length}</p>
              </div>
              <FileText className="w-8 h-8 text-purple-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              {state.isSigning ? 'Signing in progress' : 'Ready to sign'}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Security Level</p>
                <p className="text-2xl font-bold text-green-600">High</p>
              </div>
              <CheckCircle className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Hardware-secured
            </div>
          </CardContent>
        </Card>
      </div>

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
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="devices">Devices</TabsTrigger>
          <TabsTrigger value="accounts">Accounts</TabsTrigger>
          <TabsTrigger value="activity">Activity</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            {/* Current Device */}
            <Card>
              <CardHeader>
                <CardTitle>Current Device</CardTitle>
                <CardDescription>Active hardware wallet device</CardDescription>
              </CardHeader>
              <CardContent>
                {state.selectedDevice ? (
                  <div className="space-y-4">
                    <div className="flex items-center gap-3">
                      {getDeviceIcon(state.selectedDevice.type)}
                      <div>
                        <h4 className="font-medium">{state.selectedDevice.model}</h4>
                        <p className="text-sm text-muted-foreground">
                          {state.selectedDevice.type} • {state.selectedDevice.version}
                        </p>
                      </div>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Status</span>
                        <span className={getStatusColor(state.selectedDevice)}>
                          {getStatusText(state.selectedDevice)}
                        </span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Connection</span>
                        <div className="flex items-center gap-1">
                          {getConnectionIcon(state.selectedDevice.connectionMethod)}
                          <span className="capitalize">{state.selectedDevice.connectionMethod}</span>
                        </div>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Last Connected</span>
                        <span>{formatTimeAgo(state.selectedDevice.lastConnected)}</span>
                      </div>
                    </div>

                    {state.selectedAccount && (
                      <div className="pt-4 border-t">
                        <h5 className="font-medium mb-2">Selected Account</h5>
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-sm font-mono">
                              {showAddresses ? state.selectedAccount.address : formatAddress(state.selectedAccount.address)}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Account {state.selectedAccount.index + 1} • {state.selectedAccount.path}
                            </p>
                          </div>
                          <div className="flex items-center gap-1">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => setShowAddresses(!showAddresses)}
                            >
                              {showAddresses ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => copyToClipboard(state.selectedAccount!.address)}
                            >
                              <Copy className="w-4 h-4" />
                            </Button>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Wallet className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                    <p className="text-muted-foreground">No device selected</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
                <CardDescription>Common hardware wallet operations</CardDescription>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button 
                  className="w-full" 
                  onClick={handleTestSigning}
                  disabled={!state.selectedAccount || state.isSigning}
                >
                  {state.isSigning ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Signing...
                    </>
                  ) : (
                    <>
                      <FileText className="w-4 h-4 mr-2" />
                      Test Message Signing
                    </>
                  )}
                </Button>
                
                <Button 
                  variant="outline" 
                  className="w-full"
                  onClick={() => state.selectedDevice && loadAccounts(state.selectedDevice.id)}
                  disabled={!state.selectedDevice}
                >
                  <RefreshCw className="w-4 h-4 mr-2" />
                  Refresh Accounts
                </Button>
                
                <Button 
                  variant="outline" 
                  className="w-full"
                  onClick={scanDevices}
                  disabled={state.isScanning}
                >
                  <Monitor className="w-4 h-4 mr-2" />
                  Scan for Devices
                </Button>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="devices" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
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
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          {getDeviceIcon(device.type)}
                          <div>
                            <h4 className="font-medium">{device.model}</h4>
                            <p className="text-sm text-muted-foreground capitalize">
                              {device.type} • {device.version}
                            </p>
                          </div>
                        </div>
                        <div className="flex flex-col items-end gap-1">
                          <Badge variant={device.isConnected ? 'default' : 'outline'}>
                            {getStatusText(device)}
                          </Badge>
                          {device.isLocked && (
                            <Badge variant="secondary" className="text-xs">
                              <Lock className="w-3 h-3 mr-1" />
                              Locked
                            </Badge>
                          )}
                        </div>
                      </div>

                      <div className="space-y-2 mb-3">
                        <div className="flex items-center gap-2 text-sm">
                          {getConnectionIcon(device.connectionMethod)}
                          <span className="text-muted-foreground capitalize">
                            {device.connectionMethod}
                          </span>
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                          <Clock className="w-4 h-4" />
                          <span className="text-muted-foreground">
                            {formatTimeAgo(device.lastConnected)}
                          </span>
                        </div>
                      </div>

                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant={device.isConnected ? "secondary" : "default"}
                          onClick={() => handleDeviceSelect(device)}
                          disabled={state.isConnecting}
                          className="flex-1"
                        >
                          {state.isConnecting && selectedDevice?.id === device.id ? (
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
                        {device.isConnected && (
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => disconnectDevice(device.id)}
                          >
                            <Unlock className="w-4 h-4" />
                          </Button>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>

          {state.devices.length === 0 && !state.isScanning && (
            <div className="text-center py-12">
              <HardDrive className="w-16 h-16 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">No Hardware Wallets Found</h3>
              <p className="text-muted-foreground mb-4">
                Connect your hardware wallet and make sure it's unlocked
              </p>
              <Button onClick={scanDevices}>
                <RefreshCw className="w-4 h-4 mr-2" />
                Scan for Devices
              </Button>
            </div>
          )}
        </TabsContent>

        <TabsContent value="accounts" className="space-y-4">
          {state.accounts.length > 0 ? (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-medium">Available Accounts</h3>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowAddresses(!showAddresses)}
                >
                  {showAddresses ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </Button>
              </div>

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
                              onClick={() => selectAccount(account)}
                              variant={state.selectedAccount?.address === account.address ? "default" : "outline"}
                            >
                              {state.selectedAccount?.address === account.address ? 'Selected' : 'Select'}
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <div className="text-center py-12">
              <Wallet className="w-16 h-16 mx-auto text-muted-foreground mb-4" />
              <h3 className="text-lg font-medium mb-2">No Accounts Available</h3>
              <p className="text-muted-foreground mb-4">
                Connect a hardware wallet to view accounts
              </p>
              {state.selectedDevice && (
                <Button onClick={() => loadAccounts(state.selectedDevice!.id)}>
                  <RefreshCw className="w-4 h-4 mr-2" />
                  Load Accounts
                </Button>
              )}
            </div>
          )}
        </TabsContent>

        <TabsContent value="activity" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Signing Activity</CardTitle>
              <CardDescription>Recent signing requests and transactions</CardDescription>
            </CardHeader>
            <CardContent>
              {state.signingRequests.length > 0 ? (
                <div className="space-y-3">
                  {state.signingRequests.slice(0, 10).map((request) => (
                    <div key={request.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center gap-3">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                          request.status === 'signed' ? 'bg-green-100 dark:bg-green-900' :
                          request.status === 'error' ? 'bg-red-100 dark:bg-red-900' :
                          'bg-yellow-100 dark:bg-yellow-900'
                        }`}>
                          {request.status === 'signed' && <CheckCircle className="w-4 h-4 text-green-600" />}
                          {request.status === 'error' && <AlertCircle className="w-4 h-4 text-red-600" />}
                          {request.status === 'pending' && <Clock className="w-4 h-4 text-yellow-600" />}
                        </div>
                        <div>
                          <p className="font-medium capitalize">{request.type}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatTimeAgo(request.createdAt)}
                          </p>
                        </div>
                      </div>
                      <Badge variant={
                        request.status === 'signed' ? 'default' :
                        request.status === 'error' ? 'destructive' :
                        'secondary'
                      }>
                        {request.status}
                      </Badge>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">No signing activity yet</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
