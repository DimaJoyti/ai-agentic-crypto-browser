'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useAccount, useBalance, useChainId, useDisconnect, useConnectors } from 'wagmi'
import { formatEther } from 'viem'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Wallet, 
  Plus, 
  Settings, 
  Copy, 
  ExternalLink,
  Check,
  AlertCircle,
  Trash2,
  Edit,
  Star,
  StarOff,
  Shield,
  Smartphone,
  Zap
} from 'lucide-react'
import { SUPPORTED_CHAINS } from '@/lib/chains'
import { SimpleWalletModal } from './SimpleWalletModal'
import { toast } from 'sonner'

interface ConnectedWallet {
  id: string
  name: string
  address: string
  connector: string
  chainId: number
  balance: string
  isActive: boolean
  isFavorite: boolean
  lastUsed: Date
  type: 'injected' | 'walletconnect' | 'coinbase' | 'hardware'
}

interface WalletGroup {
  chainId: number
  wallets: ConnectedWallet[]
}

export function MultiWalletManager() {
  const [showConnectionModal, setShowConnectionModal] = useState(false)
  const [connectedWallets, setConnectedWallets] = useState<ConnectedWallet[]>([])
  const [activeTab, setActiveTab] = useState('all')
  const [editingWallet, setEditingWallet] = useState<string | null>(null)
  const [newWalletName, setNewWalletName] = useState('')

  const { address, isConnected, connector } = useAccount()
  const { data: balance } = useBalance({ address })
  const chainId = useChainId()
  const { disconnect } = useDisconnect()
  const connectors = useConnectors()

  // Mock data for demonstration - in real app, this would come from local storage or API
  useEffect(() => {
    if (isConnected && address) {
      const mockWallets: ConnectedWallet[] = [
        {
          id: '1',
          name: 'Main Wallet',
          address: address,
          connector: connector?.name || 'Unknown',
          chainId: chainId,
          balance: balance ? formatEther(balance.value) : '0',
          isActive: true,
          isFavorite: true,
          lastUsed: new Date(),
          type: 'injected'
        },
        {
          id: '2',
          name: 'Trading Wallet',
          address: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1',
          connector: 'MetaMask',
          chainId: 137,
          balance: '1250.45',
          isActive: false,
          isFavorite: false,
          lastUsed: new Date(Date.now() - 86400000), // 1 day ago
          type: 'injected'
        },
        {
          id: '3',
          name: 'Hardware Wallet',
          address: '0x8ba1f109551bD432803012645Hac136c22C501e',
          connector: 'Ledger',
          chainId: 1,
          balance: '5.67',
          isActive: false,
          isFavorite: true,
          lastUsed: new Date(Date.now() - 172800000), // 2 days ago
          type: 'hardware'
        }
      ]
      setConnectedWallets(mockWallets)
    }
  }, [isConnected, address, balance, chainId, connector])

  const handleAddWallet = () => {
    setShowConnectionModal(true)
  }

  const handleWalletConnectionSuccess = (address: string, chainId: number) => {
    toast.success('New wallet connected successfully!')
    setShowConnectionModal(false)
  }

  const handleToggleFavorite = (walletId: string) => {
    setConnectedWallets(prev => 
      prev.map(wallet => 
        wallet.id === walletId 
          ? { ...wallet, isFavorite: !wallet.isFavorite }
          : wallet
      )
    )
  }

  const handleRenameWallet = (walletId: string, newName: string) => {
    setConnectedWallets(prev => 
      prev.map(wallet => 
        wallet.id === walletId 
          ? { ...wallet, name: newName }
          : wallet
      )
    )
    setEditingWallet(null)
    setNewWalletName('')
    toast.success('Wallet renamed successfully!')
  }

  const handleRemoveWallet = (walletId: string) => {
    setConnectedWallets(prev => prev.filter(wallet => wallet.id !== walletId))
    toast.success('Wallet removed successfully!')
  }

  const copyAddress = (address: string) => {
    navigator.clipboard.writeText(address)
    toast.success('Address copied to clipboard')
  }

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  const getWalletIcon = (type: string) => {
    switch (type) {
      case 'hardware':
        return <Shield className="w-4 h-4" />
      case 'walletconnect':
        return <Smartphone className="w-4 h-4" />
      case 'coinbase':
        return <Shield className="w-4 h-4" />
      default:
        return <Wallet className="w-4 h-4" />
    }
  }

  const getFilteredWallets = () => {
    switch (activeTab) {
      case 'favorites':
        return connectedWallets.filter(wallet => wallet.isFavorite)
      case 'hardware':
        return connectedWallets.filter(wallet => wallet.type === 'hardware')
      case 'active':
        return connectedWallets.filter(wallet => wallet.isActive)
      default:
        return connectedWallets
    }
  }

  const groupWalletsByChain = (wallets: ConnectedWallet[]): WalletGroup[] => {
    const groups = wallets.reduce((acc, wallet) => {
      const existing = acc.find(group => group.chainId === wallet.chainId)
      if (existing) {
        existing.wallets.push(wallet)
      } else {
        acc.push({ chainId: wallet.chainId, wallets: [wallet] })
      }
      return acc
    }, [] as WalletGroup[])

    return groups.sort((a, b) => {
      // Sort by chain priority (mainnet first, then by chain ID)
      const chainA = SUPPORTED_CHAINS[a.chainId]
      const chainB = SUPPORTED_CHAINS[b.chainId]
      if (chainA?.isTestnet !== chainB?.isTestnet) {
        return chainA?.isTestnet ? 1 : -1
      }
      return a.chainId - b.chainId
    })
  }

  const filteredWallets = getFilteredWallets()
  const walletGroups = groupWalletsByChain(filteredWallets)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Multi-Wallet Manager</h2>
          <p className="text-muted-foreground">
            Manage multiple wallets across different networks
          </p>
        </div>
        <Button onClick={handleAddWallet} className="gap-2">
          <Plus className="w-4 h-4" />
          Add Wallet
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Wallets</p>
                <p className="text-2xl font-bold">{connectedWallets.length}</p>
              </div>
              <Wallet className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Active Wallets</p>
                <p className="text-2xl font-bold">
                  {connectedWallets.filter(w => w.isActive).length}
                </p>
              </div>
              <Check className="w-8 h-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Hardware Wallets</p>
                <p className="text-2xl font-bold">
                  {connectedWallets.filter(w => w.type === 'hardware').length}
                </p>
              </div>
              <Shield className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Networks</p>
                <p className="text-2xl font-bold">
                  {new Set(connectedWallets.map(w => w.chainId)).size}
                </p>
              </div>
              <Zap className="w-8 h-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Wallet List */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Connected Wallets</CardTitle>
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList>
                <TabsTrigger value="all">All</TabsTrigger>
                <TabsTrigger value="favorites">Favorites</TabsTrigger>
                <TabsTrigger value="hardware">Hardware</TabsTrigger>
                <TabsTrigger value="active">Active</TabsTrigger>
              </TabsList>
            </Tabs>
          </div>
        </CardHeader>
        <CardContent>
          {walletGroups.length === 0 ? (
            <div className="text-center py-12">
              <Wallet className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">No wallets found</h3>
              <p className="text-muted-foreground mb-6">
                {activeTab === 'all' 
                  ? 'Connect your first wallet to get started'
                  : `No wallets in the ${activeTab} category`
                }
              </p>
              <Button onClick={handleAddWallet}>
                <Plus className="w-4 h-4 mr-2" />
                Add Wallet
              </Button>
            </div>
          ) : (
            <div className="space-y-6">
              {walletGroups.map((group) => {
                const chainInfo = SUPPORTED_CHAINS[group.chainId]
                return (
                  <div key={group.chainId} className="space-y-3">
                    <div className="flex items-center gap-2">
                      <div className={`w-4 h-4 rounded-full ${chainInfo?.color || 'bg-gray-500'}`}>
                        {chainInfo?.icon}
                      </div>
                      <h3 className="font-semibold">{chainInfo?.name || `Chain ${group.chainId}`}</h3>
                      <Badge variant="secondary">{group.wallets.length} wallet{group.wallets.length !== 1 ? 's' : ''}</Badge>
                    </div>
                    
                    <div className="grid gap-3">
                      {group.wallets.map((wallet) => (
                        <motion.div
                          key={wallet.id}
                          initial={{ opacity: 0, y: 20 }}
                          animate={{ opacity: 1, y: 0 }}
                          className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                        >
                          <div className="flex items-center gap-3">
                            <div className="w-10 h-10 bg-secondary rounded-full flex items-center justify-center">
                              {getWalletIcon(wallet.type)}
                            </div>
                            <div>
                              <div className="flex items-center gap-2">
                                {editingWallet === wallet.id ? (
                                  <input
                                    type="text"
                                    value={newWalletName}
                                    onChange={(e) => setNewWalletName(e.target.value)}
                                    onBlur={() => handleRenameWallet(wallet.id, newWalletName)}
                                    onKeyDown={(e) => {
                                      if (e.key === 'Enter') {
                                        handleRenameWallet(wallet.id, newWalletName)
                                      } else if (e.key === 'Escape') {
                                        setEditingWallet(null)
                                        setNewWalletName('')
                                      }
                                    }}
                                    className="bg-background border rounded px-2 py-1 text-sm"
                                    autoFocus
                                  />
                                ) : (
                                  <span className="font-medium">{wallet.name}</span>
                                )}
                                {wallet.isActive && (
                                  <Badge variant="default" className="text-xs">Active</Badge>
                                )}
                                {wallet.isFavorite && (
                                  <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                                )}
                              </div>
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <span>{formatAddress(wallet.address)}</span>
                                <span>•</span>
                                <span>{wallet.connector}</span>
                                <span>•</span>
                                <span>{parseFloat(wallet.balance).toFixed(4)} {chainInfo?.gasToken}</span>
                              </div>
                            </div>
                          </div>
                          
                          <div className="flex items-center gap-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleToggleFavorite(wallet.id)}
                            >
                              {wallet.isFavorite ? (
                                <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                              ) : (
                                <StarOff className="w-4 h-4" />
                              )}
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => copyAddress(wallet.address)}
                            >
                              <Copy className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => {
                                setEditingWallet(wallet.id)
                                setNewWalletName(wallet.name)
                              }}
                            >
                              <Edit className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleRemoveWallet(wallet.id)}
                              className="text-red-500 hover:text-red-600"
                            >
                              <Trash2 className="w-4 h-4" />
                            </Button>
                          </div>
                        </motion.div>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </CardContent>
      </Card>

      <SimpleWalletModal
        isOpen={showConnectionModal}
        onClose={() => setShowConnectionModal(false)}
        onSuccess={handleWalletConnectionSuccess}
      />
    </div>
  )
}
