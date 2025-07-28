'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { useAccount, useBalance, useChainId, useSwitchChain, useDisconnect } from 'wagmi'
import { formatEther } from 'viem'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Progress } from '@/components/ui/progress'
import { 
  Wallet, 
  TrendingUp, 
  TrendingDown, 
  Copy, 
  ExternalLink, 
  Settings,
  RefreshCw,
  Shield,
  Zap,
  DollarSign,
  BarChart3,
  Eye,
  EyeOff
} from 'lucide-react'
import { WalletConnectionModal } from './WalletConnectionModal'
import { ChainSwitcher } from './ChainSwitcher'
import { toast } from 'sonner'

interface TokenBalance {
  symbol: string
  name: string
  balance: string
  value: number
  change24h: number
  icon?: string
}

interface WalletStats {
  totalValue: number
  change24h: number
  transactionCount: number
  gasSpent: number
}

export function WalletManager() {
  const [showConnectionModal, setShowConnectionModal] = useState(false)
  const [showBalances, setShowBalances] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [walletStats, setWalletStats] = useState<WalletStats>({
    totalValue: 0,
    change24h: 0,
    transactionCount: 0,
    gasSpent: 0
  })
  const [tokenBalances, setTokenBalances] = useState<TokenBalance[]>([])

  const { address, isConnected, connector } = useAccount()
  const { data: balance, refetch: refetchBalance } = useBalance({ address })
  const chainId = useChainId()
  const { switchChain } = useSwitchChain()
  const { disconnect } = useDisconnect()

  // Mock data for demonstration - in real app, fetch from API
  useEffect(() => {
    if (isConnected && address) {
      // Simulate fetching wallet stats and token balances
      setWalletStats({
        totalValue: 12450.67,
        change24h: 5.23,
        transactionCount: 142,
        gasSpent: 0.234
      })

      setTokenBalances([
        {
          symbol: 'ETH',
          name: 'Ethereum',
          balance: balance ? formatEther(balance.value) : '0',
          value: 2456.78,
          change24h: 3.45
        },
        {
          symbol: 'USDC',
          name: 'USD Coin',
          balance: '5000.00',
          value: 5000.00,
          change24h: 0.01
        },
        {
          symbol: 'UNI',
          name: 'Uniswap',
          balance: '125.50',
          value: 892.35,
          change24h: -2.15
        }
      ])
    }
  }, [isConnected, address, balance])

  const handleRefresh = async () => {
    setIsRefreshing(true)
    try {
      await refetchBalance()
      // Refresh other data
      toast.success('Wallet data refreshed')
    } catch (error) {
      toast.error('Failed to refresh wallet data')
    } finally {
      setIsRefreshing(false)
    }
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

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(value)
  }

  const getChainName = (chainId: number) => {
    const chains: Record<number, string> = {
      1: 'Ethereum',
      137: 'Polygon',
      42161: 'Arbitrum',
      10: 'Optimism',
      11155111: 'Sepolia'
    }
    return chains[chainId] || `Chain ${chainId}`
  }

  if (!isConnected) {
    return (
      <div className="space-y-6">
        <Card>
          <CardContent className="p-8 text-center">
            <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
              <Wallet className="w-8 h-8" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Connect Your Wallet</h3>
            <p className="text-muted-foreground mb-6">
              Connect your cryptocurrency wallet to access Web3 features, manage your portfolio, and interact with DeFi protocols.
            </p>
            <Button onClick={() => setShowConnectionModal(true)} className="w-full">
              <Wallet className="w-4 h-4 mr-2" />
              Connect Wallet
            </Button>
          </CardContent>
        </Card>

        <WalletConnectionModal
          isOpen={showConnectionModal}
          onClose={() => setShowConnectionModal(false)}
          onSuccess={() => setShowConnectionModal(false)}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Wallet Header */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                <Wallet className="w-5 h-5 text-primary" />
              </div>
              <div>
                <CardTitle className="text-lg">
                  {connector?.name || 'Connected Wallet'}
                </CardTitle>
                <CardDescription className="flex items-center gap-2">
                  {formatAddress(address!)}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={copyAddress}
                    className="h-auto p-1"
                  >
                    <Copy className="w-3 h-3" />
                  </Button>
                </CardDescription>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <ChainSwitcher />
              <Button
                variant="outline"
                size="sm"
                onClick={handleRefresh}
                disabled={isRefreshing}
              >
                <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => disconnect()}
              >
                Disconnect
              </Button>
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* Portfolio Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Value</p>
                <p className="text-2xl font-bold">
                  {showBalances ? formatCurrency(walletStats.totalValue) : '••••••'}
                </p>
              </div>
              <DollarSign className="w-8 h-8 text-green-500" />
            </div>
            <div className="flex items-center gap-1 mt-2">
              {walletStats.change24h >= 0 ? (
                <TrendingUp className="w-4 h-4 text-green-500" />
              ) : (
                <TrendingDown className="w-4 h-4 text-red-500" />
              )}
              <span className={`text-sm ${walletStats.change24h >= 0 ? 'text-green-500' : 'text-red-500'}`}>
                {walletStats.change24h >= 0 ? '+' : ''}{walletStats.change24h}%
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Transactions</p>
                <p className="text-2xl font-bold">{walletStats.transactionCount}</p>
              </div>
              <BarChart3 className="w-8 h-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Gas Spent</p>
                <p className="text-2xl font-bold">{walletStats.gasSpent} ETH</p>
              </div>
              <Zap className="w-8 h-8 text-yellow-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Network</p>
                <p className="text-lg font-semibold">{getChainName(chainId)}</p>
              </div>
              <Shield className="w-8 h-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Wallet Tabs */}
      <Tabs defaultValue="tokens" className="space-y-4">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="tokens">Tokens</TabsTrigger>
            <TabsTrigger value="nfts">NFTs</TabsTrigger>
            <TabsTrigger value="defi">DeFi</TabsTrigger>
            <TabsTrigger value="history">History</TabsTrigger>
          </TabsList>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowBalances(!showBalances)}
          >
            {showBalances ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
          </Button>
        </div>

        <TabsContent value="tokens" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Token Balances</CardTitle>
              <CardDescription>Your cryptocurrency holdings</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {tokenBalances.map((token, index) => (
                  <motion.div
                    key={token.symbol}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center justify-between p-3 rounded-lg border"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-secondary rounded-full flex items-center justify-center">
                        <span className="font-semibold text-sm">{token.symbol}</span>
                      </div>
                      <div>
                        <p className="font-medium">{token.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {showBalances ? `${parseFloat(token.balance).toFixed(4)} ${token.symbol}` : '••••••'}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">
                        {showBalances ? formatCurrency(token.value) : '••••••'}
                      </p>
                      <div className="flex items-center gap-1">
                        {token.change24h >= 0 ? (
                          <TrendingUp className="w-3 h-3 text-green-500" />
                        ) : (
                          <TrendingDown className="w-3 h-3 text-red-500" />
                        )}
                        <span className={`text-xs ${token.change24h >= 0 ? 'text-green-500' : 'text-red-500'}`}>
                          {token.change24h >= 0 ? '+' : ''}{token.change24h}%
                        </span>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="nfts">
          <Card>
            <CardContent className="p-8 text-center">
              <p className="text-muted-foreground">NFT collection coming soon...</p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="defi">
          <Card>
            <CardContent className="p-8 text-center">
              <p className="text-muted-foreground">DeFi positions coming soon...</p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="history">
          <Card>
            <CardContent className="p-8 text-center">
              <p className="text-muted-foreground">Transaction history coming soon...</p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <WalletConnectionModal
        isOpen={showConnectionModal}
        onClose={() => setShowConnectionModal(false)}
        onSuccess={() => setShowConnectionModal(false)}
      />
    </div>
  )
}
