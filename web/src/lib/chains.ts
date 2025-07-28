import { Chain } from 'wagmi/chains'
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
  gnosis
} from 'wagmi/chains'

export interface ExtendedChainInfo {
  id: number
  name: string
  shortName: string
  nativeCurrency: {
    name: string
    symbol: string
    decimals: number
  }
  rpcUrls: any // Use any to avoid type conflicts with wagmi chains
  blockExplorers: any // Use any to avoid type conflicts with wagmi chains
  icon: string
  color: string
  gasToken: string
  avgGasPrice: string
  blockTime: string
  isTestnet: boolean
  category: 'mainnet' | 'layer2' | 'sidechain' | 'testnet'
  features: string[]
  ecosystem: string[]
  tvl?: string
  dailyTxs?: string
  status: 'healthy' | 'congested' | 'degraded' | 'maintenance'
}

export const SUPPORTED_CHAINS: Record<number, ExtendedChainInfo> = {
  // Ethereum Mainnet
  1: {
    id: 1,
    name: 'Ethereum Mainnet',
    shortName: 'Ethereum',
    nativeCurrency: mainnet.nativeCurrency,
    rpcUrls: mainnet.rpcUrls,
    blockExplorers: mainnet.blockExplorers,
    icon: '🔷',
    color: 'bg-blue-500',
    gasToken: 'ETH',
    avgGasPrice: '25 gwei',
    blockTime: '12s',
    isTestnet: false,
    category: 'mainnet',
    features: ['DeFi', 'NFTs', 'Smart Contracts', 'Staking'],
    ecosystem: ['Uniswap', 'Aave', 'Compound', 'MakerDAO'],
    tvl: '$50B+',
    dailyTxs: '1M+',
    status: 'healthy'
  },

  // Polygon
  137: {
    id: 137,
    name: 'Polygon Mainnet',
    shortName: 'Polygon',
    nativeCurrency: polygon.nativeCurrency,
    rpcUrls: polygon.rpcUrls,
    blockExplorers: polygon.blockExplorers,
    icon: '🟣',
    color: 'bg-purple-500',
    gasToken: 'MATIC',
    avgGasPrice: '30 gwei',
    blockTime: '2s',
    isTestnet: false,
    category: 'sidechain',
    features: ['Low Fees', 'Fast Transactions', 'EVM Compatible'],
    ecosystem: ['QuickSwap', 'SushiSwap', 'Aave'],
    tvl: '$1.2B+',
    dailyTxs: '3M+',
    status: 'healthy'
  },

  // Arbitrum One
  42161: {
    id: 42161,
    name: 'Arbitrum One',
    shortName: 'Arbitrum',
    nativeCurrency: arbitrum.nativeCurrency,
    rpcUrls: arbitrum.rpcUrls,
    blockExplorers: arbitrum.blockExplorers,
    icon: '🔵',
    color: 'bg-blue-400',
    gasToken: 'ETH',
    avgGasPrice: '0.1 gwei',
    blockTime: '1s',
    isTestnet: false,
    category: 'layer2',
    features: ['Optimistic Rollup', 'Low Fees', 'Ethereum Security'],
    ecosystem: ['GMX', 'Camelot', 'Radiant'],
    tvl: '$2.1B+',
    dailyTxs: '500K+',
    status: 'congested'
  },

  // Optimism
  10: {
    id: 10,
    name: 'Optimism',
    shortName: 'Optimism',
    nativeCurrency: optimism.nativeCurrency,
    rpcUrls: optimism.rpcUrls,
    blockExplorers: optimism.blockExplorers,
    icon: '🔴',
    color: 'bg-red-500',
    gasToken: 'ETH',
    avgGasPrice: '0.001 gwei',
    blockTime: '2s',
    isTestnet: false,
    category: 'layer2',
    features: ['Optimistic Rollup', 'OP Token', 'Retroactive Funding'],
    ecosystem: ['Velodrome', 'Synthetix', 'Kwenta'],
    tvl: '$800M+',
    dailyTxs: '200K+',
    status: 'healthy'
  },

  // Base
  8453: {
    id: 8453,
    name: 'Base',
    shortName: 'Base',
    nativeCurrency: base.nativeCurrency,
    rpcUrls: base.rpcUrls,
    blockExplorers: base.blockExplorers,
    icon: '🔵',
    color: 'bg-blue-600',
    gasToken: 'ETH',
    avgGasPrice: '0.01 gwei',
    blockTime: '2s',
    isTestnet: false,
    category: 'layer2',
    features: ['Coinbase Backed', 'Low Fees', 'Developer Friendly'],
    ecosystem: ['Aerodrome', 'BaseSwap', 'Moonwell'],
    tvl: '$400M+',
    dailyTxs: '150K+',
    status: 'healthy'
  },

  // Avalanche
  43114: {
    id: 43114,
    name: 'Avalanche C-Chain',
    shortName: 'Avalanche',
    nativeCurrency: avalanche.nativeCurrency,
    rpcUrls: avalanche.rpcUrls,
    blockExplorers: avalanche.blockExplorers,
    icon: '🔺',
    color: 'bg-red-600',
    gasToken: 'AVAX',
    avgGasPrice: '25 nAVAX',
    blockTime: '1s',
    isTestnet: false,
    category: 'mainnet',
    features: ['Sub-second Finality', 'Subnets', 'EVM Compatible'],
    ecosystem: ['Trader Joe', 'Pangolin', 'Benqi'],
    tvl: '$600M+',
    dailyTxs: '100K+',
    status: 'healthy'
  },

  // BNB Smart Chain
  56: {
    id: 56,
    name: 'BNB Smart Chain',
    shortName: 'BSC',
    nativeCurrency: bsc.nativeCurrency,
    rpcUrls: bsc.rpcUrls,
    blockExplorers: bsc.blockExplorers,
    icon: '🟡',
    color: 'bg-yellow-500',
    gasToken: 'BNB',
    avgGasPrice: '5 gwei',
    blockTime: '3s',
    isTestnet: false,
    category: 'sidechain',
    features: ['Low Fees', 'Fast Transactions', 'Binance Ecosystem'],
    ecosystem: ['PancakeSwap', 'Venus', 'Alpaca Finance'],
    tvl: '$3.2B+',
    dailyTxs: '2M+',
    status: 'healthy'
  },

  // Fantom
  250: {
    id: 250,
    name: 'Fantom Opera',
    shortName: 'Fantom',
    nativeCurrency: fantom.nativeCurrency,
    rpcUrls: fantom.rpcUrls,
    blockExplorers: fantom.blockExplorers,
    icon: '👻',
    color: 'bg-blue-300',
    gasToken: 'FTM',
    avgGasPrice: '20 gwei',
    blockTime: '1s',
    isTestnet: false,
    category: 'mainnet',
    features: ['DAG Technology', 'Fast Finality', 'Low Fees'],
    ecosystem: ['SpookySwap', 'Beethoven X', 'Geist Finance'],
    tvl: '$100M+',
    dailyTxs: '50K+',
    status: 'healthy'
  },

  // Gnosis Chain
  100: {
    id: 100,
    name: 'Gnosis Chain',
    shortName: 'Gnosis',
    nativeCurrency: gnosis.nativeCurrency,
    rpcUrls: gnosis.rpcUrls,
    blockExplorers: gnosis.blockExplorers,
    icon: '🟢',
    color: 'bg-green-500',
    gasToken: 'xDAI',
    avgGasPrice: '2 gwei',
    blockTime: '5s',
    isTestnet: false,
    category: 'sidechain',
    features: ['Stable Fees', 'PoS Consensus', 'Community Owned'],
    ecosystem: ['Honeyswap', 'Symmetric', 'Agave'],
    tvl: '$50M+',
    dailyTxs: '20K+',
    status: 'healthy'
  },

  // Sepolia Testnet
  11155111: {
    id: 11155111,
    name: 'Sepolia Testnet',
    shortName: 'Sepolia',
    nativeCurrency: sepolia.nativeCurrency,
    rpcUrls: sepolia.rpcUrls,
    blockExplorers: sepolia.blockExplorers,
    icon: '🧪',
    color: 'bg-yellow-500',
    gasToken: 'ETH',
    avgGasPrice: '1 gwei',
    blockTime: '12s',
    isTestnet: true,
    category: 'testnet',
    features: ['Testing', 'Development', 'Faucet Available'],
    ecosystem: ['Testnet DApps'],
    status: 'healthy'
  }
}

export const CHAIN_CATEGORIES = {
  mainnet: 'Mainnet',
  layer2: 'Layer 2',
  sidechain: 'Sidechain',
  testnet: 'Testnet'
} as const

export const WAGMI_CHAINS: Chain[] = [
  mainnet,
  polygon,
  arbitrum,
  optimism,
  base,
  avalanche,
  bsc,
  fantom,
  gnosis,
  sepolia
]

// Helper functions
export function getChainInfo(chainId: number): ExtendedChainInfo | undefined {
  return SUPPORTED_CHAINS[chainId]
}

export function getChainsByCategory(category: keyof typeof CHAIN_CATEGORIES): ExtendedChainInfo[] {
  return Object.values(SUPPORTED_CHAINS).filter(chain => chain.category === category)
}

export function getMainnetChains(): ExtendedChainInfo[] {
  return Object.values(SUPPORTED_CHAINS).filter(chain => !chain.isTestnet)
}

export function getTestnetChains(): ExtendedChainInfo[] {
  return Object.values(SUPPORTED_CHAINS).filter(chain => chain.isTestnet)
}

export function isChainSupported(chainId: number): boolean {
  return chainId in SUPPORTED_CHAINS
}

export function getChainStatus(chainId: number): string {
  return SUPPORTED_CHAINS[chainId]?.status || 'unknown'
}

export function formatChainName(chainId: number): string {
  const chain = SUPPORTED_CHAINS[chainId]
  return chain ? chain.shortName : `Chain ${chainId}`
}
