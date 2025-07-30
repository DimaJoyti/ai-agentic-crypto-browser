import { createConfig, http, fallback } from 'wagmi'
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
  walletConnect,
  coinbaseWallet,
  safe,
  metaMask
} from 'wagmi/connectors'

// Environment variables with fallbacks
const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || ''
const infuraKey = process.env.NEXT_PUBLIC_INFURA_KEY || ''
const alchemyKey = process.env.NEXT_PUBLIC_ALCHEMY_KEY || ''
const quicknodeKey = process.env.NEXT_PUBLIC_QUICKNODE_KEY || ''

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

// Enhanced connector configuration with better error handling and features
const connectors = [
  // MetaMask with enhanced configuration
  metaMask({
    dappMetadata: {
      name: 'AI Agentic Browser',
      url: typeof window !== 'undefined' ? window.location.origin : 'https://ai-agentic-browser.com',
      iconUrl: 'https://ai-agentic-browser.com/icon.png'
    },
    extensionOnly: false,
    preferDesktop: true,
    infuraAPIKey: infuraKey,
  }),

  // Generic injected connector for other browser wallets
  injected({
    target: 'metaMask'
  }),

  // WalletConnect with enhanced configuration
  walletConnect({
    projectId,
    metadata: {
      name: 'AI Agentic Browser',
      description: 'AI-powered web browser with Web3 integration and DeFi capabilities',
      url: typeof window !== 'undefined' ? window.location.origin : 'https://ai-agentic-browser.com',
      icons: ['https://ai-agentic-browser.com/icon.png']
    },
    showQrModal: true,
    qrModalOptions: {
      themeMode: 'light',
      themeVariables: {
        '--wcm-z-index': '1000'
      }
    }
  }),

  // Coinbase Wallet with enhanced configuration
  coinbaseWallet({
    appName: 'AI Agentic Browser',
    appLogoUrl: 'https://ai-agentic-browser.com/icon.png',
    headlessMode: false
  }),

  // Safe (Gnosis Safe) connector
  safe({
    allowedDomains: [/gnosis-safe.io$/, /app.safe.global$/],
    debug: process.env.NODE_ENV === 'development'
  }),


]

export const config = createConfig({
  chains,
  connectors,
  transports,
  ssr: true,

})
