'use client'

import React, { useMemo, useState, useEffect } from 'react'
import { ConnectionProvider, WalletProvider } from '@solana/wallet-adapter-react'
import { WalletAdapterNetwork } from '@solana/wallet-adapter-base'
import { WalletModalProvider } from '@solana/wallet-adapter-react-ui'
import {
  PhantomWalletAdapter,
  SolflareWalletAdapter,
  LedgerWalletAdapter,
  TrezorWalletAdapter,
  TorusWalletAdapter,
  CoinbaseWalletAdapter,
  MathWalletAdapter,
  Coin98WalletAdapter,
  CloverWalletAdapter,
  SafePalWalletAdapter,
  SolongWalletAdapter,
  TokenPocketWalletAdapter,
  TrustWalletAdapter
} from '@solana/wallet-adapter-wallets'
import { clusterApiUrl } from '@solana/web3.js'

// Import wallet adapter CSS
import '@solana/wallet-adapter-react-ui/styles.css'

interface SolanaWalletProviderProps {
  children: React.ReactNode
  network?: WalletAdapterNetwork
  endpoint?: string
  autoConnect?: boolean
}

export function SolanaWalletProvider({
  children,
  network = WalletAdapterNetwork.Mainnet,
  endpoint,
  autoConnect = true
}: SolanaWalletProviderProps) {
  const [isClient, setIsClient] = useState(false)

  useEffect(() => {
    setIsClient(true)
  }, [])

  // The network can be set to 'devnet', 'testnet', or 'mainnet-beta'
  const solanaNetwork = network

  // You can also provide a custom RPC endpoint
  const rpcEndpoint = useMemo(() => {
    if (endpoint) return endpoint
    
    // Use custom RPC endpoints for better performance
    switch (solanaNetwork) {
      case WalletAdapterNetwork.Mainnet:
        return process.env.NEXT_PUBLIC_SOLANA_RPC_URL || 
               process.env.NEXT_PUBLIC_HELIUS_RPC_URL ||
               'https://api.mainnet-beta.solana.com'
      case WalletAdapterNetwork.Devnet:
        return process.env.NEXT_PUBLIC_SOLANA_DEVNET_RPC_URL || 
               clusterApiUrl(WalletAdapterNetwork.Devnet)
      case WalletAdapterNetwork.Testnet:
        return clusterApiUrl(WalletAdapterNetwork.Testnet)
      default:
        return clusterApiUrl(solanaNetwork)
    }
  }, [solanaNetwork, endpoint])

  // Configure supported wallets
  const wallets = useMemo(
    () => [
      // Most popular Solana wallets
      new PhantomWalletAdapter(),
      new SolflareWalletAdapter({ network: solanaNetwork }),

      // Note: Some wallets like BackpackWalletAdapter, GlowWalletAdapter, and SlopeWalletAdapter
      // are not available in @solana/wallet-adapter-wallets and require separate packages:
      // - @backpack-wallet/wallet-adapter for Backpack
      // - @glow-wallet/wallet-adapter for Glow
      // - @slope-finance/wallet-adapter for Slope

      // Hardware wallets
      new LedgerWalletAdapter(),
      new TrezorWalletAdapter({
        email: process.env.NEXT_PUBLIC_TREZOR_EMAIL || 'support@example.com'
      }),

      // Web wallets
      new TorusWalletAdapter(),
      new CoinbaseWalletAdapter(),

      // Mobile and other wallets
      new MathWalletAdapter(),
      new Coin98WalletAdapter(),
      new CloverWalletAdapter(),
      new SafePalWalletAdapter(),
      new SolongWalletAdapter(),
      new TokenPocketWalletAdapter(),
      new TrustWalletAdapter()
    ],
    [solanaNetwork]
  )

  // Prevent hydration mismatch by only rendering wallet components on client
  if (!isClient) {
    return <div>{children}</div>
  }

  return (
    <ConnectionProvider
      endpoint={rpcEndpoint}
      config={{
        commitment: 'confirmed',
        wsEndpoint: rpcEndpoint.replace('https://', 'wss://').replace('http://', 'ws://'),
        confirmTransactionInitialTimeout: 60000
      }}
    >
      <WalletProvider
        wallets={wallets}
        autoConnect={autoConnect}
        onError={(error) => {
          console.error('Wallet error:', error)
          // You can add custom error handling here
        }}
      >
        <WalletModalProvider>
          {children}
        </WalletModalProvider>
      </WalletProvider>
    </ConnectionProvider>
  )
}

// Hook to get current network info
export function useSolanaNetwork() {
  return {
    network: WalletAdapterNetwork.Mainnet, // This could be dynamic based on environment
    isMainnet: true,
    isDevnet: false,
    isTestnet: false
  }
}

// Hook to get RPC endpoint info
export function useSolanaRPC() {
  const network: WalletAdapterNetwork = WalletAdapterNetwork.Mainnet

  const endpoint = useMemo(() => {
    return process.env.NEXT_PUBLIC_SOLANA_RPC_URL ||
           process.env.NEXT_PUBLIC_HELIUS_RPC_URL ||
           'https://api.mainnet-beta.solana.com'
  }, [])

  return {
    endpoint,
    network,
    isCustomRPC: !!process.env.NEXT_PUBLIC_SOLANA_RPC_URL
  }
}
