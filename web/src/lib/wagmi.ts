import { createConfig, http, fallback, createStorage } from 'wagmi'
import {
  mainnet,
  polygon,
  arbitrum,
  optimism,
  sepolia,
  base,
  avalanche,
  bsc,
  fantom,
  gnosis,
  polygonMumbai,
  arbitrumGoerli,
  optimismGoerli,
  goerli
} from 'wagmi/chains'
import {
  injected,
  coinbaseWallet,
  safe,
  metaMask
} from 'wagmi/connectors'
import { createWalletConnectConnector, isWalletConnectAvailable } from './walletconnect-client'

// Environment variables with fallbacks
const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || ''
const infuraKey = process.env.NEXT_PUBLIC_INFURA_KEY || ''
const alchemyKey = process.env.NEXT_PUBLIC_ALCHEMY_KEY || ''
const quicknodeKey = process.env.NEXT_PUBLIC_QUICKNODE_KEY || ''

// Validate required environment variables
if (!projectId && process.env.NODE_ENV === 'production') {
  console.warn('NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID is not set. WalletConnect features may not work properly.')
}

// Enhanced chain configuration with multiple RPC providers for redundancy
const chains = [
  mainnet,
  polygon,
  arbitrum,
  optimism,
  base,
  avalanche,
  bsc,
  fantom,
  gnosis,
  sepolia,
  // Testnets for development
  ...(process.env.NODE_ENV === 'development' ? [
    polygonMumbai,
    arbitrumGoerli,
    optimismGoerli,
    goerli
  ] : [])
] as const

// Create enhanced transports with fallback providers for reliability
const createTransportWithFallback = (chainId: number) => {
  const transports = []

  // Primary providers
  if (infuraKey) {
    transports.push(http(`https://mainnet.infura.io/v3/${infuraKey}`))
  }
  if (alchemyKey) {
    transports.push(http(`https://eth-mainnet.g.alchemy.com/v2/${alchemyKey}`))
  }
  if (quicknodeKey) {
    transports.push(http(`https://api.quicknode.com/${quicknodeKey}`))
  }

  // Fallback to public RPC
  transports.push(http())

  return fallback(transports)
}

const transports = Object.fromEntries(
  chains.map(chain => [chain.id, createTransportWithFallback(chain.id)])
) as Record<number, any>

// Server-safe connectors (without WalletConnect to prevent SSR issues)
const getServerConnectors = () => [
  // MetaMask with enhanced configuration
  metaMask({
    dappMetadata: {
      name: 'AI Agentic Browser',
      url: process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000',
      iconUrl: `${process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'}/icons/icon-192x192.svg`
    },
    extensionOnly: false,
    preferDesktop: true,
    infuraAPIKey: infuraKey,
  }),

  // Generic injected connector for other browser wallets
  injected({
    target: 'metaMask'
  }),

  // Coinbase Wallet with enhanced configuration
  coinbaseWallet({
    appName: 'AI Agentic Browser',
    appLogoUrl: `${process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'}/icons/icon-192x192.svg`,
    headlessMode: false
  }),

  // Safe (Gnosis Safe) connector
  safe({
    allowedDomains: [/gnosis-safe.io$/, /app.safe.global$/],
    debug: process.env.NODE_ENV === 'development'
  }),
]

// Client-side connectors (with WalletConnect)
const getClientConnectors = () => [
  // MetaMask with enhanced configuration
  metaMask({
    dappMetadata: {
      name: 'AI Agentic Browser',
      url: window.location.origin,
      iconUrl: `${window.location.origin}/icons/icon-192x192.svg`
    },
    extensionOnly: false,
    preferDesktop: true,
    infuraAPIKey: infuraKey,
  }),

  // Generic injected connector for other browser wallets
  injected({
    target: 'metaMask'
  }),

  // WalletConnect with enhanced configuration (only if project ID is available and client-side)
  ...(projectId && isWalletConnectAvailable() ? [createWalletConnectConnector(projectId)] : []),

  // Coinbase Wallet with enhanced configuration
  coinbaseWallet({
    appName: 'AI Agentic Browser',
    appLogoUrl: `${window.location.origin}/icons/icon-192x192.svg`,
    headlessMode: false
  }),

  // Safe (Gnosis Safe) connector
  safe({
    allowedDomains: [/gnosis-safe.io$/, /app.safe.global$/],
    debug: process.env.NODE_ENV === 'development'
  }),
]

// Server-safe config (without WalletConnect)
export const serverConfig = createConfig({
  chains,
  connectors: getServerConnectors(),
  transports,
  ssr: true,
})

// Client-side config (with WalletConnect)
export const clientConfig = createConfig({
  chains,
  connectors: getClientConnectors(),
  transports,
  ssr: true,
  storage: typeof window !== 'undefined' ? createStorage({ storage: window.localStorage }) : undefined,
})

// Default export for backward compatibility
export const config = typeof window !== 'undefined' ? clientConfig : serverConfig
