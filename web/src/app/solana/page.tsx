'use client'

import React, { useState, useEffect } from 'react'
import dynamic from 'next/dynamic'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import {
  Zap,
  Shield,
  TrendingUp,
  Coins,
  Info,
  ExternalLink,
  CheckCircle,
  Loader2
} from 'lucide-react'

// Dynamically import wallet components to prevent SSR issues
const SolanaWalletProvider = dynamic(
  () => import('@/components/solana/SolanaWalletProvider').then(mod => ({ default: mod.SolanaWalletProvider })),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center p-8">
        <Loader2 className="h-6 w-6 animate-spin" />
        <span className="ml-2">Loading wallet...</span>
      </div>
    )
  }
)

const SolanaMarketDashboard = dynamic(
  () => import('@/components/solana/SolanaMarketDashboard').then(mod => ({ default: mod.SolanaMarketDashboard })),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center p-8">
        <Loader2 className="h-6 w-6 animate-spin" />
        <span className="ml-2">Loading market data...</span>
      </div>
    )
  }
)

const SolanaWalletConnect = dynamic(
  () => import('@/components/solana/SolanaWalletConnect').then(mod => ({ default: mod.SolanaWalletConnect })),
  {
    ssr: false,
    loading: () => (
      <Card className="max-w-md mx-auto">
        <CardContent className="p-6">
          <div className="text-center">
            <Loader2 className="h-6 w-6 animate-spin mx-auto mb-2" />
            <p className="text-sm text-muted-foreground">Loading wallet connection...</p>
          </div>
        </CardContent>
      </Card>
    )
  }
)

export default function SolanaPage() {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  // Prevent hydration mismatch by not rendering until mounted
  if (!mounted) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
          <p className="text-sm text-muted-foreground">Loading Solana trading hub...</p>
        </div>
      </div>
    )
  }

  return (
    <SolanaWalletProvider autoConnect={true}>
      <div className="min-h-screen bg-background">
        <div className="container mx-auto px-4 py-8 space-y-8">
          {/* Header Section */}
          <div className="text-center space-y-4">
            <div className="flex items-center justify-center space-x-2">
              <Coins className="h-8 w-8 text-purple-600" />
              <h1 className="text-4xl font-bold tracking-tight">
                Solana Trading Hub
              </h1>
            </div>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Trade tokens, explore DeFi protocols, and discover NFT collections on the fastest blockchain
            </p>
            <div className="flex items-center justify-center space-x-4">
              <Badge variant="outline" className="text-green-600 border-green-600">
                <CheckCircle className="h-3 w-3 mr-1" />
                Mainnet
              </Badge>
              <Badge variant="outline" className="text-blue-600 border-blue-600">
                <Zap className="h-3 w-3 mr-1" />
                65,000+ TPS
              </Badge>
              <Badge variant="outline" className="text-purple-600 border-purple-600">
                <TrendingUp className="h-3 w-3 mr-1" />
                $0.00025 Fees
              </Badge>
            </div>
          </div>

          {/* Features Overview */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Zap className="h-5 w-5 mr-2 text-yellow-600" />
                  Lightning Fast
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Execute trades in milliseconds with Solana's high-performance blockchain
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Shield className="h-5 w-5 mr-2 text-green-600" />
                  Secure & Reliable
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Built with industry-leading security practices and wallet integrations
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <TrendingUp className="h-5 w-5 mr-2 text-blue-600" />
                  Best Prices
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Access the best prices across all major Solana DEXs and protocols
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Wallet Connection */}
          <div className="max-w-md mx-auto">
            <SolanaWalletConnect 
              showBalance={true} 
              showTokens={true}
            />
          </div>

          {/* Important Notice */}
          <Alert>
            <Info className="h-4 w-4" />
            <AlertDescription>
              This is a demo interface. Always verify transactions and use at your own risk. 
              Make sure you understand the risks involved in DeFi trading.
            </AlertDescription>
          </Alert>

          {/* Main Dashboard */}
          <SolanaMarketDashboard 
            autoRefresh={true}
            refreshInterval={30000}
          />

          {/* Footer Information */}
          <div className="border-t pt-8">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div>
                <h3 className="font-semibold mb-3">Supported Wallets</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Phantom</li>
                  <li>• Solflare</li>
                  <li>• Backpack</li>
                  <li>• Glow</li>
                  <li>• Ledger & Trezor</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">DeFi Protocols</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Jupiter (DEX Aggregator)</li>
                  <li>• Raydium (AMM)</li>
                  <li>• Orca (DEX)</li>
                  <li>• Marinade (Liquid Staking)</li>
                  <li>• Kamino (Lending)</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">NFT Marketplaces</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• Magic Eden</li>
                  <li>• Tensor</li>
                  <li>• Solanart</li>
                  <li>• Exchange Art</li>
                  <li>• OpenSea (Solana)</li>
                </ul>
              </div>

              <div>
                <h3 className="font-semibold mb-3">Resources</h3>
                <ul className="space-y-2 text-sm">
                  <li>
                    <a 
                      href="https://docs.solana.com" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-foreground flex items-center"
                    >
                      Solana Docs
                      <ExternalLink className="h-3 w-3 ml-1" />
                    </a>
                  </li>
                  <li>
                    <a 
                      href="https://explorer.solana.com" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-foreground flex items-center"
                    >
                      Solana Explorer
                      <ExternalLink className="h-3 w-3 ml-1" />
                    </a>
                  </li>
                  <li>
                    <a 
                      href="https://status.solana.com" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-foreground flex items-center"
                    >
                      Network Status
                      <ExternalLink className="h-3 w-3 ml-1" />
                    </a>
                  </li>
                  <li>
                    <a 
                      href="https://solana.com/ecosystem" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-foreground flex items-center"
                    >
                      Ecosystem
                      <ExternalLink className="h-3 w-3 ml-1" />
                    </a>
                  </li>
                </ul>
              </div>
            </div>

            <div className="mt-8 pt-8 border-t text-center text-sm text-muted-foreground">
              <p>
                Built with ❤️ for the Solana ecosystem. 
                Always DYOR (Do Your Own Research) before trading.
              </p>
            </div>
          </div>
        </div>
      </div>
    </SolanaWalletProvider>
  )
}
