import { type Connector } from 'wagmi'
import { type Address } from 'viem'

// Extend Window interface for wallet providers
declare global {
  interface Window {
    phantom?: {
      ethereum?: any
      solana?: any
    }
  }
}

export interface WalletProvider {
  id: string
  name: string
  icon: string
  description: string
  type: 'injected' | 'walletconnect' | 'coinbase' | 'hardware' | 'safe'
  isInstalled?: boolean
  downloadUrl?: string
  features: string[]
  supported: boolean
  priority: number
  category: 'browser' | 'mobile' | 'hardware' | 'institutional'
  security: 'high' | 'medium' | 'low'
  platforms: ('web' | 'mobile' | 'desktop')[]
}

export interface WalletConnectionState {
  isConnecting: boolean
  isConnected: boolean
  address?: Address
  chainId?: number
  connector?: Connector
  error?: string
  lastConnected?: Date
}

export interface WalletConnectionOptions {
  autoConnect?: boolean
  preferredChainId?: number
  timeout?: number
  retryAttempts?: number
  onSuccess?: (address: Address, chainId: number) => void
  onError?: (error: Error) => void
  onDisconnect?: () => void
}

// Comprehensive wallet provider definitions
export const WALLET_PROVIDERS: WalletProvider[] = [
  {
    id: 'metamask',
    name: 'MetaMask',
    icon: '/icons/metamask.svg',
    description: 'The most popular Ethereum wallet browser extension',
    type: 'injected',
    isInstalled: typeof window !== 'undefined' && window.ethereum?.isMetaMask,
    downloadUrl: 'https://metamask.io/download/',
    features: ['Browser Extension', 'Mobile App', 'Hardware Wallet Support', 'Swaps', 'Staking'],
    supported: true,
    priority: 1,
    category: 'browser',
    security: 'high',
    platforms: ['web', 'mobile']
  },
  {
    id: 'walletconnect',
    name: 'WalletConnect',
    icon: '/icons/walletconnect.svg',
    description: 'Connect to 300+ mobile wallets via QR code',
    type: 'walletconnect',
    isInstalled: typeof window !== 'undefined', // Only available client-side
    features: ['Mobile Wallets', 'QR Code', 'Cross-Platform', '300+ Wallets'],
    supported: typeof window !== 'undefined', // Only supported client-side
    priority: 2,
    category: 'mobile',
    security: 'high',
    platforms: ['web', 'mobile']
  },
  {
    id: 'coinbase',
    name: 'Coinbase Wallet',
    icon: '/icons/coinbase.svg',
    description: 'Self-custody wallet from Coinbase',
    type: 'coinbase',
    isInstalled: typeof window !== 'undefined' && window.ethereum?.isCoinbaseWallet,
    downloadUrl: 'https://www.coinbase.com/wallet',
    features: ['Browser Extension', 'Mobile App', 'DeFi Access', 'NFT Support'],
    supported: true,
    priority: 3,
    category: 'browser',
    security: 'high',
    platforms: ['web', 'mobile']
  },
  {
    id: 'ledger',
    name: 'Ledger',
    icon: '/icons/ledger.svg',
    description: 'Hardware wallet with maximum security',
    type: 'hardware',
    isInstalled: true,
    downloadUrl: 'https://www.ledger.com/',
    features: ['Hardware Security', 'Multi-Currency', 'Staking', 'DeFi'],
    supported: true,
    priority: 4,
    category: 'hardware',
    security: 'high',
    platforms: ['web', 'desktop']
  },
  {
    id: 'trezor',
    name: 'Trezor',
    icon: '/icons/trezor.svg',
    description: 'Open-source hardware wallet',
    type: 'hardware',
    isInstalled: true,
    downloadUrl: 'https://trezor.io/',
    features: ['Hardware Security', 'Open Source', 'Multi-Currency', 'Privacy'],
    supported: true,
    priority: 5,
    category: 'hardware',
    security: 'high',
    platforms: ['web', 'desktop']
  },
  {
    id: 'safe',
    name: 'Safe (Gnosis)',
    icon: '/icons/safe.svg',
    description: 'Multi-signature smart contract wallet',
    type: 'safe',
    isInstalled: true,
    features: ['Multi-Signature', 'Smart Contract', 'Team Management', 'Advanced Security'],
    supported: true,
    priority: 6,
    category: 'institutional',
    security: 'high',
    platforms: ['web']
  },
  {
    id: 'rainbow',
    name: 'Rainbow',
    icon: '/icons/rainbow.svg',
    description: 'Fun, simple, and secure Ethereum wallet',
    type: 'walletconnect',
    isInstalled: true,
    downloadUrl: 'https://rainbow.me/',
    features: ['Mobile First', 'Beautiful UI', 'DeFi', 'NFTs'],
    supported: true,
    priority: 7,
    category: 'mobile',
    security: 'high',
    platforms: ['mobile']
  },
  {
    id: 'trust',
    name: 'Trust Wallet',
    icon: '/icons/trust.svg',
    description: 'Multi-cryptocurrency wallet',
    type: 'walletconnect',
    isInstalled: true,
    downloadUrl: 'https://trustwallet.com/',
    features: ['Multi-Chain', 'DeFi', 'NFTs', 'Staking'],
    supported: true,
    priority: 8,
    category: 'mobile',
    security: 'high',
    platforms: ['mobile']
  },
  {
    id: 'phantom',
    name: 'Phantom',
    icon: '/icons/phantom.svg',
    description: 'Solana and Ethereum wallet',
    type: 'injected',
    isInstalled: typeof window !== 'undefined' && window.phantom?.ethereum,
    downloadUrl: 'https://phantom.app/',
    features: ['Multi-Chain', 'Solana Support', 'NFTs', 'Swaps'],
    supported: true,
    priority: 9,
    category: 'browser',
    security: 'high',
    platforms: ['web', 'mobile']
  },
  {
    id: 'brave',
    name: 'Brave Wallet',
    icon: '/icons/brave.svg',
    description: 'Built-in Brave browser wallet',
    type: 'injected',
    isInstalled: typeof window !== 'undefined' && window.ethereum?.isBraveWallet,
    downloadUrl: 'https://brave.com/wallet/',
    features: ['Built-in Browser', 'Privacy Focused', 'Multi-Chain', 'No Extensions'],
    supported: true,
    priority: 10,
    category: 'browser',
    security: 'high',
    platforms: ['web', 'mobile']
  }
]

// Wallet detection utilities
export const detectInstalledWallets = (): WalletProvider[] => {
  if (typeof window === 'undefined') return []
  
  return WALLET_PROVIDERS.filter(wallet => {
    switch (wallet.id) {
      case 'metamask':
        return window.ethereum?.isMetaMask
      case 'coinbase':
        return window.ethereum?.isCoinbaseWallet
      case 'phantom':
        return window.phantom?.ethereum || window.phantom?.solana
      case 'brave':
        return window.ethereum?.isBraveWallet
      case 'walletconnect':
      case 'ledger':
      case 'trezor':
      case 'safe':
      case 'rainbow':
      case 'trust':
        return true // These don't require installation detection
      default:
        return false
    }
  })
}

export const getRecommendedWallets = (category?: string): WalletProvider[] => {
  let wallets = WALLET_PROVIDERS.filter(w => w.supported)
  
  if (category) {
    wallets = wallets.filter(w => w.category === category)
  }
  
  return wallets.sort((a, b) => a.priority - b.priority)
}

export const getWalletByConnector = (connector: Connector): WalletProvider | undefined => {
  const connectorId = connector.id.toLowerCase()
  const connectorName = connector.name.toLowerCase()
  
  return WALLET_PROVIDERS.find(wallet => 
    wallet.id === connectorId || 
    wallet.name.toLowerCase() === connectorName ||
    connectorName.includes(wallet.id)
  )
}

// Wallet connection error handling
export class WalletConnectionError extends Error {
  constructor(
    message: string,
    public code: string,
    public walletId?: string,
    public originalError?: Error
  ) {
    super(message)
    this.name = 'WalletConnectionError'
  }
}

export const WALLET_ERROR_CODES = {
  USER_REJECTED: 'USER_REJECTED',
  WALLET_NOT_FOUND: 'WALLET_NOT_FOUND',
  WALLET_NOT_INSTALLED: 'WALLET_NOT_INSTALLED',
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
  UNKNOWN: 'UNKNOWN'
} as const

export const getWalletErrorMessage = (error: Error, walletId?: string): string => {
  const message = error.message.toLowerCase()
  
  if (message.includes('user rejected') || message.includes('user denied')) {
    return 'Connection was rejected by the user'
  }
  
  if (message.includes('not found') || message.includes('not installed')) {
    const wallet = WALLET_PROVIDERS.find(w => w.id === walletId)
    return `${wallet?.name || 'Wallet'} is not installed. Please install it first.`
  }
  
  if (message.includes('network') || message.includes('rpc')) {
    return 'Network connection error. Please check your internet connection.'
  }
  
  if (message.includes('timeout')) {
    return 'Connection timeout. Please try again.'
  }
  
  return error.message || 'An unknown error occurred while connecting to the wallet'
}

// Storage utilities for wallet preferences
export const WALLET_STORAGE_KEYS = {
  LAST_CONNECTED: 'wallet_last_connected',
  PREFERRED_WALLET: 'wallet_preferred',
  AUTO_CONNECT: 'wallet_auto_connect',
  CONNECTION_HISTORY: 'wallet_connection_history'
} as const

export const saveWalletPreference = (walletId: string): void => {
  try {
    localStorage.setItem(WALLET_STORAGE_KEYS.PREFERRED_WALLET, walletId)
    localStorage.setItem(WALLET_STORAGE_KEYS.LAST_CONNECTED, Date.now().toString())
  } catch (error) {
    console.warn('Failed to save wallet preference:', error)
  }
}

export const getWalletPreference = (): string | null => {
  try {
    return localStorage.getItem(WALLET_STORAGE_KEYS.PREFERRED_WALLET)
  } catch (error) {
    console.warn('Failed to get wallet preference:', error)
    return null
  }
}

export const clearWalletPreference = (): void => {
  try {
    localStorage.removeItem(WALLET_STORAGE_KEYS.PREFERRED_WALLET)
    localStorage.removeItem(WALLET_STORAGE_KEYS.LAST_CONNECTED)
  } catch (error) {
    console.warn('Failed to clear wallet preference:', error)
  }
}
