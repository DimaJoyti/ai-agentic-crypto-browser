'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Wallet,
  TrendingUp,
  Coins,
  Image,
  Activity,
  Shield,
  Zap,
  Globe,
  Settings,
  BarChart3,
  FileText,
  DollarSign
} from 'lucide-react'
import { WalletManager } from '@/components/web3/WalletManager'
import { MultiWalletManager } from '@/components/web3/MultiWalletManager'
import { SecuritySettings } from '@/components/web3/SecuritySettings'
import { WalletBackup } from '@/components/web3/WalletBackup'
import { TransactionTracker } from '@/components/web3/TransactionTracker'
import { TransactionNotifications } from '@/components/web3/TransactionNotifications'
import { TransactionAnalytics } from '@/components/web3/TransactionAnalytics'
import { GasOptimizer } from '@/components/web3/GasOptimizer'
import { GasTrackerWidget } from '@/components/web3/GasTrackerWidget'
import { TransactionBatcher } from '@/components/web3/TransactionBatcher'
import { DeFiDashboard } from '@/components/web3/DeFiDashboard'
import { TokenSwap } from '@/components/web3/TokenSwap'
import { LiquidityPools } from '@/components/web3/LiquidityPools'
import { PositionManager } from '@/components/web3/PositionManager'
import { YieldFarmingDashboard } from '@/components/web3/YieldFarmingDashboard'
import { YieldOptimizer } from '@/components/web3/YieldOptimizer'
import { NFTCollectionDashboard } from '@/components/web3/NFTCollectionDashboard'
import { MarketplaceDashboard } from '@/components/web3/MarketplaceDashboard'
import { NFTAnalyticsDashboard } from '@/components/web3/NFTAnalyticsDashboard'
import { DeFiAnalyticsDashboard } from '@/components/web3/DeFiAnalyticsDashboard'
import { NFTMintingDashboard } from '@/components/web3/NFTMintingDashboard'
import { TransactionHistoryDashboard } from '@/components/web3/TransactionHistoryDashboard'

const features = [
  {
    icon: <Wallet className="w-8 h-8" />,
    title: 'Multi-Wallet Support',
    description: 'Connect MetaMask, WalletConnect, Coinbase Wallet, and hardware wallets',
    color: 'bg-blue-500'
  },
  {
    icon: <Globe className="w-8 h-8" />,
    title: 'Multi-Chain',
    description: 'Support for Ethereum, Polygon, Arbitrum, Optimism, and more',
    color: 'bg-green-500'
  },
  {
    icon: <TrendingUp className="w-8 h-8" />,
    title: 'DeFi Integration',
    description: 'Interact with Uniswap, Aave, Compound, and other protocols',
    color: 'bg-purple-500'
  },
  {
    icon: <Image className="w-8 h-8" />,
    title: 'NFT Support',
    description: 'Browse, trade, and manage your NFT collections',
    color: 'bg-pink-500'
  },
  {
    icon: <Activity className="w-8 h-8" />,
    title: 'Real-time Monitoring',
    description: 'Track transactions, prices, and portfolio performance',
    color: 'bg-orange-500'
  },
  {
    icon: <Shield className="w-8 h-8" />,
    title: 'Security First',
    description: 'Advanced security features and hardware wallet support',
    color: 'bg-red-500'
  }
]

export default function Web3Page() {
  const [activeTab, setActiveTab] = useState('wallet')

  // Mock connected address - in real app, this would come from wallet context
  const connectedAddress = undefined // Will be replaced with actual wallet connection

  return (
    <div className="container mx-auto px-4 py-8 space-y-8 relative">
      {/* Transaction Notifications */}
      <TransactionNotifications />
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center space-y-4"
      >
        <div className="flex items-center justify-center gap-3">
          <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
            <Zap className="w-6 h-6 text-primary" />
          </div>
          <h1 className="text-4xl font-bold">Web3 Integration</h1>
        </div>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Connect your wallet and explore the decentralized web with our comprehensive Web3 features
        </p>
      </motion.div>

      {/* Features Grid */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        {features.map((feature, index) => (
          <motion.div
            key={feature.title}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 + index * 0.1 }}
          >
            <Card className="h-full hover:shadow-lg transition-shadow">
              <CardContent className="p-6">
                <div className="flex items-start gap-4">
                  <div className={`w-12 h-12 ${feature.color} rounded-lg flex items-center justify-center text-white`}>
                    {feature.icon}
                  </div>
                  <div className="flex-1">
                    <h3 className="font-semibold mb-2">{feature.title}</h3>
                    <p className="text-sm text-muted-foreground">{feature.description}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </motion.div>

      {/* Main Content */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
      >
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-9">
            <TabsTrigger value="wallet" className="flex items-center gap-2">
              <Wallet className="w-4 h-4" />
              Wallet
            </TabsTrigger>
            <TabsTrigger value="multi-wallet" className="flex items-center gap-2">
              <Settings className="w-4 h-4" />
              Multi-Wallet
            </TabsTrigger>
            <TabsTrigger value="security" className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              Security
            </TabsTrigger>
            <TabsTrigger value="transactions" className="flex items-center gap-2">
              <Activity className="w-4 h-4" />
              Transactions
            </TabsTrigger>
            <TabsTrigger value="analytics" className="flex items-center gap-2">
              <BarChart3 className="w-4 h-4" />
              Analytics
            </TabsTrigger>
            <TabsTrigger value="gas" className="flex items-center gap-2">
              <Zap className="w-4 h-4" />
              Gas
            </TabsTrigger>
            <TabsTrigger value="defi" className="flex items-center gap-2">
              <TrendingUp className="w-4 h-4" />
              DeFi
            </TabsTrigger>
            <TabsTrigger value="nfts" className="flex items-center gap-2">
              <Image className="w-4 h-4" />
              NFTs
            </TabsTrigger>
            <TabsTrigger value="reports" className="flex items-center gap-2">
              <FileText className="w-4 h-4" />
              Reports
            </TabsTrigger>
          </TabsList>

          <TabsContent value="wallet" className="space-y-6">
            <WalletManager />
          </TabsContent>

          <TabsContent value="multi-wallet" className="space-y-6">
            <MultiWalletManager />
          </TabsContent>

          <TabsContent value="security" className="space-y-6">
            <Tabs defaultValue="settings" className="space-y-4">
              <TabsList>
                <TabsTrigger value="settings">Security Settings</TabsTrigger>
                <TabsTrigger value="backup">Backup & Recovery</TabsTrigger>
              </TabsList>
              <TabsContent value="settings">
                <SecuritySettings />
              </TabsContent>
              <TabsContent value="backup">
                <WalletBackup />
              </TabsContent>
            </Tabs>
          </TabsContent>

          <TabsContent value="transactions" className="space-y-6">
            <TransactionTracker />
          </TabsContent>

          <TabsContent value="analytics" className="space-y-6">
            <TransactionAnalytics />
          </TabsContent>

          <TabsContent value="gas" className="space-y-6">
            <Tabs defaultValue="optimizer" className="space-y-4">
              <TabsList>
                <TabsTrigger value="optimizer">Gas Optimizer</TabsTrigger>
                <TabsTrigger value="batcher">Transaction Batcher</TabsTrigger>
                <TabsTrigger value="tracker">Gas Tracker</TabsTrigger>
              </TabsList>
              <TabsContent value="optimizer">
                <GasOptimizer chainId={1} />
              </TabsContent>
              <TabsContent value="batcher">
                <TransactionBatcher />
              </TabsContent>
              <TabsContent value="tracker">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                  <GasTrackerWidget chainId={1} />
                  <GasTrackerWidget chainId={137} />
                  <GasTrackerWidget chainId={42161} />
                  <GasTrackerWidget chainId={10} />
                </div>
              </TabsContent>
            </Tabs>
          </TabsContent>

          <TabsContent value="defi" className="space-y-6">
            <Tabs defaultValue="dashboard" className="space-y-4">
              <TabsList>
                <TabsTrigger value="dashboard">Dashboard</TabsTrigger>
                <TabsTrigger value="positions">Positions</TabsTrigger>
                <TabsTrigger value="yield">Yield Farming</TabsTrigger>
                <TabsTrigger value="optimizer">Optimizer</TabsTrigger>
                <TabsTrigger value="defi-analytics">DeFi Analytics</TabsTrigger>
                <TabsTrigger value="nfts">NFT Collection</TabsTrigger>
                <TabsTrigger value="marketplace">Marketplace</TabsTrigger>
                <TabsTrigger value="minting">NFT Minting</TabsTrigger>
                <TabsTrigger value="analytics">NFT Analytics</TabsTrigger>
                <TabsTrigger value="history">Transaction History</TabsTrigger>
                <TabsTrigger value="swap">Token Swap</TabsTrigger>
                <TabsTrigger value="liquidity">Liquidity Pools</TabsTrigger>
                <TabsTrigger value="lending">Lending</TabsTrigger>
              </TabsList>
              <TabsContent value="dashboard">
                <DeFiDashboard chainId={1} userAddress={connectedAddress} />
              </TabsContent>
              <TabsContent value="positions">
                <PositionManager userAddress={connectedAddress as any} />
              </TabsContent>
              <TabsContent value="yield">
                <YieldFarmingDashboard userAddress={connectedAddress as any} chainId={1} />
              </TabsContent>
              <TabsContent value="optimizer">
                <YieldOptimizer userAddress={connectedAddress as any} chainId={1} />
              </TabsContent>
              <TabsContent value="defi-analytics">
                <DeFiAnalyticsDashboard userAddress={connectedAddress as any} />
              </TabsContent>
              <TabsContent value="nfts">
                <NFTCollectionDashboard userAddress={connectedAddress as any} chainId={1} />
              </TabsContent>
              <TabsContent value="marketplace">
                <MarketplaceDashboard userAddress={connectedAddress as any} />
              </TabsContent>
              <TabsContent value="minting">
                <NFTMintingDashboard
                  userAddress={connectedAddress as any}
                  chainId={1}
                />
              </TabsContent>
              <TabsContent value="analytics">
                <NFTAnalyticsDashboard
                  contractAddress={'0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' as any}
                  tokenId="1234"
                />
              </TabsContent>
              <TabsContent value="history">
                <TransactionHistoryDashboard
                  userAddress={connectedAddress as any}
                  chainId={1}
                />
              </TabsContent>
              <TabsContent value="swap">
                <TokenSwap chainId={1} userAddress={connectedAddress as any} />
              </TabsContent>
              <TabsContent value="liquidity">
                <LiquidityPools chainId={1} userAddress={connectedAddress as any} />
              </TabsContent>
              <TabsContent value="lending">
                <div className="text-center py-8">
                  <DollarSign className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-semibold mb-2">Lending Integration</h3>
                  <p className="text-muted-foreground">
                    Advanced lending and borrowing features coming soon
                  </p>
                </div>
              </TabsContent>
            </Tabs>
          </TabsContent>

          <TabsContent value="nfts" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Image className="w-5 h-5" />
                  NFT Collections
                </CardTitle>
                <CardDescription>
                  Browse and manage your NFT collections
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center py-12">
                  <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                    <Image className="w-8 h-8" />
                  </div>
                  <h3 className="text-lg font-semibold mb-2">NFT Marketplace Integration</h3>
                  <p className="text-muted-foreground mb-6">
                    Connect your wallet to view and manage your NFT collections
                  </p>
                  <Button disabled>
                    Coming Soon
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analytics" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="w-5 h-5" />
                  Portfolio Analytics
                </CardTitle>
                <CardDescription>
                  Track your portfolio performance and analytics
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-center py-12">
                  <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                    <Activity className="w-8 h-8" />
                  </div>
                  <h3 className="text-lg font-semibold mb-2">Advanced Analytics</h3>
                  <p className="text-muted-foreground mb-6">
                    Comprehensive portfolio tracking and performance analytics
                  </p>
                  <Button disabled>
                    Coming Soon
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </motion.div>

      {/* Call to Action */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <Card className="bg-gradient-to-r from-primary/10 to-secondary/10">
          <CardContent className="p-8 text-center">
            <h2 className="text-2xl font-bold mb-4">Ready to Explore Web3?</h2>
            <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
              Connect your wallet to start using our comprehensive Web3 features. 
              Manage your portfolio, interact with DeFi protocols, and explore the decentralized web.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" className="gap-2">
                <Wallet className="w-5 h-5" />
                Get Started
              </Button>
              <Button variant="outline" size="lg">
                Learn More
              </Button>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}
