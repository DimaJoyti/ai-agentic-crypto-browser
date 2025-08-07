'use client'

import { useState, useEffect, useCallback, useMemo } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAccount, useChainId } from 'wagmi'
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
  Wallet,
  Smartphone,
  Shield,
  Zap,
  CheckCircle,
  AlertCircle,
  Loader2,
  ExternalLink,
  Copy,
  LogOut,
  Monitor,
  HardDrive,
  Building,
  Star,
  Download,
  Info
} from 'lucide-react'
import { HardwareWalletModal } from './HardwareWalletModal'
import { useWalletConnection } from '@/hooks/useWalletConnection'
import {
  type WalletProvider,
  WALLET_PROVIDERS,
  detectInstalledWallets,
  getRecommendedWallets
} from '@/lib/wallet-providers'
import { toast } from 'sonner'

interface WalletConnectionModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (address: string, chainId: number) => void
}

const getCategoryIcon = (category: string) => {
  switch (category) {
    case 'browser':
      return <Monitor className="w-5 h-5" />
    case 'mobile':
      return <Smartphone className="w-5 h-5" />
    case 'hardware':
      return <HardDrive className="w-5 h-5" />
    case 'institutional':
      return <Building className="w-5 h-5" />
    default:
      return <Wallet className="w-5 h-5" />
  }
}

const getSecurityBadgeColor = (security: string) => {
  switch (security) {
    case 'high':
      return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
    case 'medium':
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
    case 'low':
      return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
  }
}

export function WalletConnectionModal({ isOpen, onClose, onSuccess }: WalletConnectionModalProps) {
  const [selectedWallet, setSelectedWallet] = useState<string | null>(null)
  const [showHardwareModal, setShowHardwareModal] = useState(false)
  const [activeTab, setActiveTab] = useState('recommended')

  const { address, isConnected } = useAccount()
  const chainId = useChainId()

  // Use our enhanced wallet connection hook without onSuccess callback to prevent infinite loops
  const {
    connectionState,
    connectWallet,
    disconnectWallet,
    availableWallets,
    connectedWallet,
    clearError,
    enableAutoConnect,
    disableAutoConnect,
    isAutoConnectEnabled
  } = useWalletConnection()

  // Memoize wallet categories to prevent infinite re-renders
  const installedWallets = useMemo(() => detectInstalledWallets(), [])
  const recommendedWallets = useMemo(() => getRecommendedWallets().slice(0, 6), [])
  const browserWallets = useMemo(() => availableWallets.filter(w => w.category === 'browser'), [availableWallets])
  const mobileWallets = useMemo(() => availableWallets.filter(w => w.category === 'mobile'), [availableWallets])
  const hardwareWallets = useMemo(() => availableWallets.filter(w => w.category === 'hardware'), [availableWallets])
  const institutionalWallets = useMemo(() => availableWallets.filter(w => w.category === 'institutional'), [availableWallets])

  useEffect(() => {
    if (isConnected && address) {
      onSuccess?.(address, chainId)
      onClose()
    }
  }, [isConnected, address, chainId]) // Remove onSuccess and onClose from dependencies to prevent infinite loops

  const handleWalletConnect = async (wallet: WalletProvider) => {
    if (wallet.category === 'hardware') {
      setShowHardwareModal(true)
      return
    }

    if (!wallet.isInstalled && wallet.downloadUrl) {
      handleInstallWallet(wallet.downloadUrl)
      return
    }

    setSelectedWallet(wallet.id)
    clearError()

    try {
      await connectWallet(wallet.id)
    } catch (error) {
      console.error('Wallet connection error:', error)
      // Error handling is done in the hook
    }
  }

  const handleHardwareWalletSuccess = (deviceInfo: any) => {
    setShowHardwareModal(false)
    toast.success(`Hardware wallet ${deviceInfo.model} connected successfully!`)
    // In a real implementation, you would handle the hardware wallet connection here
    onSuccess?.('hardware-wallet-address', 1) // Mock address and chain ID
    onClose()
  }

  const handleInstallWallet = (downloadUrl: string) => {
    window.open(downloadUrl, '_blank')
  }

  const copyAddress = () => {
    if (address) {
      navigator.clipboard.writeText(address)
      toast.success('Address copied to clipboard')
    }
  }

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  const handleDisconnect = async () => {
    try {
      await disconnectWallet()
      onClose()
    } catch (error) {
      console.error('Disconnect error:', error)
    }
  }

  const renderWalletCard = (wallet: WalletProvider) => (
    <motion.div
      key={wallet.id}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
    >
      <Card
        className={`cursor-pointer transition-all hover:shadow-md border-2 ${
          selectedWallet === wallet.id ? 'border-primary' : 'border-transparent'
        } ${!wallet.isInstalled ? 'opacity-60' : ''}`}
        onClick={() => handleWalletConnect(wallet)}
      >
        <CardContent className="p-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-secondary rounded-lg flex items-center justify-center">
                {getCategoryIcon(wallet.category)}
              </div>
              <div>
                <h4 className="font-medium">{wallet.name}</h4>
                <p className="text-sm text-muted-foreground">{wallet.description}</p>
              </div>
            </div>
            <div className="flex flex-col items-end gap-1">
              <Badge className={getSecurityBadgeColor(wallet.security)}>
                {wallet.security}
              </Badge>
              {!wallet.isInstalled && (
                <Badge variant="outline" className="text-xs">
                  <Download className="w-3 h-3 mr-1" />
                  Install
                </Badge>
              )}
            </div>
          </div>

          <div className="flex flex-wrap gap-1 mb-2">
            {wallet.features.slice(0, 3).map((feature, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                {feature}
              </Badge>
            ))}
            {wallet.features.length > 3 && (
              <Badge variant="secondary" className="text-xs">
                +{wallet.features.length - 3} more
              </Badge>
            )}
          </div>

          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span className="capitalize">{wallet.category}</span>
            {wallet.isInstalled && (
              <div className="flex items-center gap-1 text-green-600">
                <CheckCircle className="w-3 h-3" />
                Installed
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </motion.div>
  )

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Wallet className="w-5 h-5" />
            {isConnected ? 'Wallet Connected' : 'Connect Your Wallet'}
          </DialogTitle>
          <DialogDescription>
            {isConnected
              ? `Connected to ${connectedWallet?.name || 'wallet'} on ${chainId === 1 ? 'Ethereum' : 'network'}`
              : 'Choose from our supported wallets to get started with Web3 features'
            }
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {isConnected ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="space-y-4"
            >
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
                        <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
                      </div>
                      <div>
                        <p className="font-medium">{connectedWallet?.name || 'Wallet'} Connected</p>
                        <p className="text-sm text-muted-foreground">
                          {formatAddress(address!)}
                        </p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={copyAddress}
                      >
                        <Copy className="w-4 h-4" />
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={handleDisconnect}
                      >
                        <LogOut className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>

                  {connectedWallet && (
                    <div className="mt-3 pt-3 border-t">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Security Level</span>
                        <Badge className={getSecurityBadgeColor(connectedWallet.security)}>
                          {connectedWallet.security}
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between text-sm mt-1">
                        <span className="text-muted-foreground">Auto-connect</span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={isAutoConnectEnabled ? disableAutoConnect : enableAutoConnect}
                        >
                          {isAutoConnectEnabled ? 'Enabled' : 'Disabled'}
                        </Button>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>

              <div className="flex gap-2">
                <Button onClick={onClose} className="flex-1">
                  Continue to App
                </Button>
              </div>
            </motion.div>
          ) : (
            <>
              {connectionState.error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    {connectionState.error}
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={clearError}
                      className="ml-2"
                    >
                      Dismiss
                    </Button>
                  </AlertDescription>
                </Alert>
              )}

              {connectionState.isConnecting && (
                <Alert>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <AlertDescription>
                    Connecting to {selectedWallet}... Please check your wallet.
                  </AlertDescription>
                </Alert>
              )}

              <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                <TabsList className="grid w-full grid-cols-4">
                  <TabsTrigger value="recommended" className="text-xs">
                    <Star className="w-3 h-3 mr-1" />
                    Recommended
                  </TabsTrigger>
                  <TabsTrigger value="browser" className="text-xs">
                    <Monitor className="w-3 h-3 mr-1" />
                    Browser
                  </TabsTrigger>
                  <TabsTrigger value="mobile" className="text-xs">
                    <Smartphone className="w-3 h-3 mr-1" />
                    Mobile
                  </TabsTrigger>
                  <TabsTrigger value="hardware" className="text-xs">
                    <HardDrive className="w-3 h-3 mr-1" />
                    Hardware
                  </TabsTrigger>
                </TabsList>

                <TabsContent value="recommended" className="space-y-3 mt-4">
                  <div className="grid gap-3">
                    <AnimatePresence>
                      {recommendedWallets.map((wallet) => renderWalletCard(wallet))}
                    </AnimatePresence>
                  </div>
                </TabsContent>

                <TabsContent value="browser" className="space-y-3 mt-4">
                  <div className="grid gap-3">
                    <AnimatePresence>
                      {browserWallets.map((wallet) => renderWalletCard(wallet))}
                    </AnimatePresence>
                  </div>
                </TabsContent>

                <TabsContent value="mobile" className="space-y-3 mt-4">
                  <div className="grid gap-3">
                    <AnimatePresence>
                      {mobileWallets.map((wallet) => renderWalletCard(wallet))}
                    </AnimatePresence>
                  </div>
                </TabsContent>

                <TabsContent value="hardware" className="space-y-3 mt-4">
                  <div className="grid gap-3">
                    <AnimatePresence>
                      {hardwareWallets.map((wallet) => renderWalletCard(wallet))}
                    </AnimatePresence>
                  </div>
                </TabsContent>
              </Tabs>

              <div className="text-center">
                <p className="text-sm text-muted-foreground">
                  By connecting a wallet, you agree to our{' '}
                  <a href="/terms" className="text-primary hover:underline">
                    Terms of Service
                  </a>{' '}
                  and{' '}
                  <a href="/privacy" className="text-primary hover:underline">
                    Privacy Policy
                  </a>
                </p>
              </div>
            </>
          )}
        </div>
      </DialogContent>

      <HardwareWalletModal
        isOpen={showHardwareModal}
        onClose={() => setShowHardwareModal(false)}
        onSuccess={handleHardwareWalletSuccess}
      />
    </Dialog>
  )
}
