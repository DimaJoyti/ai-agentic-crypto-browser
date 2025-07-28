'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useConnect, useAccount, useDisconnect, useChainId, useSwitchChain } from 'wagmi'
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
  LogOut
} from 'lucide-react'
import { HardwareWalletModal } from './HardwareWalletModal'
import { toast } from 'sonner'

interface WalletConnectionModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: (address: string, chainId: number) => void
}

interface WalletOption {
  id: string
  name: string
  icon: React.ReactNode
  description: string
  type: 'injected' | 'walletconnect' | 'coinbase' | 'hardware'
  isInstalled?: boolean
  downloadUrl?: string
  features: string[]
}

const walletOptions: WalletOption[] = [
  {
    id: 'metamask',
    name: 'MetaMask',
    icon: <Wallet className="w-8 h-8" />,
    description: 'Connect using MetaMask browser extension',
    type: 'injected',
    isInstalled: typeof window !== 'undefined' && window.ethereum?.isMetaMask,
    downloadUrl: 'https://metamask.io/download/',
    features: ['Browser Extension', 'Mobile App', 'Hardware Wallet Support']
  },
  {
    id: 'walletconnect',
    name: 'WalletConnect',
    icon: <Smartphone className="w-8 h-8" />,
    description: 'Connect using WalletConnect protocol',
    type: 'walletconnect',
    isInstalled: true,
    features: ['Mobile Wallets', 'QR Code', 'Cross-Platform']
  },
  {
    id: 'coinbase',
    name: 'Coinbase Wallet',
    icon: <Shield className="w-8 h-8" />,
    description: 'Connect using Coinbase Wallet',
    type: 'coinbase',
    isInstalled: typeof window !== 'undefined' && window.ethereum?.isCoinbaseWallet,
    downloadUrl: 'https://wallet.coinbase.com/',
    features: ['Self-Custody', 'DeFi Access', 'NFT Support']
  },
  {
    id: 'hardware',
    name: 'Hardware Wallet',
    icon: <Zap className="w-8 h-8" />,
    description: 'Connect Ledger, Trezor, or GridPlus',
    type: 'hardware',
    isInstalled: true,
    features: ['Maximum Security', 'Offline Storage', 'Multi-Chain']
  }
]

export function WalletConnectionModal({ isOpen, onClose, onSuccess }: WalletConnectionModalProps) {
  const [selectedWallet, setSelectedWallet] = useState<string | null>(null)
  const [isConnecting, setIsConnecting] = useState(false)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [showHardwareModal, setShowHardwareModal] = useState(false)
  
  const { connect, connectors, isPending, error } = useConnect()
  const { address, isConnected } = useAccount()
  const { disconnect } = useDisconnect()
  const chainId = useChainId()
  const { switchChain } = useSwitchChain()

  useEffect(() => {
    if (isConnected && address) {
      onSuccess?.(address, chainId)
      onClose()
      toast.success('Wallet connected successfully!')
    }
  }, [isConnected, address, chainId, onSuccess, onClose])

  useEffect(() => {
    if (error) {
      setConnectionError(error.message)
      setIsConnecting(false)
    }
  }, [error])

  const handleWalletConnect = async (walletId: string) => {
    if (walletId === 'hardware') {
      setShowHardwareModal(true)
      return
    }

    setSelectedWallet(walletId)
    setIsConnecting(true)
    setConnectionError(null)

    try {
      const connector = connectors.find(c =>
        c.id.toLowerCase().includes(walletId) ||
        c.name.toLowerCase().includes(walletId)
      )

      if (!connector) {
        throw new Error(`Connector for ${walletId} not found`)
      }

      await connect({ connector })
    } catch (err) {
      console.error('Wallet connection error:', err)
      setConnectionError(err instanceof Error ? err.message : 'Failed to connect wallet')
      setIsConnecting(false)
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

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Wallet className="w-5 h-5" />
            {isConnected ? 'Wallet Connected' : 'Connect Wallet'}
          </DialogTitle>
          <DialogDescription>
            {isConnected 
              ? 'Your wallet is successfully connected to the AI Agentic Browser'
              : 'Choose a wallet to connect to the AI Agentic Browser'
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
                        <p className="font-medium">Connected</p>
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
                        onClick={() => disconnect()}
                      >
                        <LogOut className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <div className="flex gap-2">
                <Button onClick={onClose} className="flex-1">
                  Continue
                </Button>
              </div>
            </motion.div>
          ) : (
            <>
              {connectionError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{connectionError}</AlertDescription>
                </Alert>
              )}

              <div className="grid gap-3">
                <AnimatePresence>
                  {walletOptions.map((wallet, index) => (
                    <motion.div
                      key={wallet.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      <Card 
                        className={`cursor-pointer transition-all hover:shadow-md ${
                          selectedWallet === wallet.id ? 'ring-2 ring-primary' : ''
                        } ${!wallet.isInstalled ? 'opacity-60' : ''}`}
                        onClick={() => {
                          if (wallet.isInstalled) {
                            handleWalletConnect(wallet.id)
                          } else if (wallet.downloadUrl) {
                            handleInstallWallet(wallet.downloadUrl)
                          }
                        }}
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
                                  {!wallet.isInstalled && (
                                    <Badge variant="outline" className="text-xs">
                                      Install
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
                            <div className="flex items-center gap-2">
                              {isConnecting && selectedWallet === wallet.id ? (
                                <Loader2 className="w-5 h-5 animate-spin" />
                              ) : !wallet.isInstalled ? (
                                <ExternalLink className="w-5 h-5" />
                              ) : (
                                <div className="w-5 h-5" />
                              )}
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    </motion.div>
                  ))}
                </AnimatePresence>
              </div>

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
