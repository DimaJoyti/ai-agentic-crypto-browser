import { type Address, type Chain } from 'viem'
import { 
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
  polygonMumbai,
  arbitrumGoerli,
  optimismGoerli,
  goerli
} from 'viem/chains'

export interface ChainConfig extends Chain {
  logo: string
  color: string
  blockExplorer: string
  nativeCurrency: {
    name: string
    symbol: string
    decimals: number
    logo: string
  }
  rpcUrls: {
    default: { http: string[] }
    public: { http: string[] }
    infura?: { http: string[] }
    alchemy?: { http: string[] }
    quicknode?: { http: string[] }
  }
  features: string[]
  gasPrice: {
    slow: string
    standard: string
    fast: string
  }
  bridgeSupport: string[]
  dexes: string[]
  lending: string[]
  category: 'mainnet' | 'testnet' | 'l2' | 'sidechain'
  tvl?: string
  dailyVolume?: string
}

export interface ChainBalance {
  chainId: number
  chainName: string
  nativeBalance: string
  nativeBalanceUSD: string
  tokens: TokenBalance[]
  totalValueUSD: string
  lastUpdated: number
}

export interface TokenBalance {
  address: Address
  symbol: string
  name: string
  decimals: number
  balance: string
  balanceFormatted: string
  valueUSD: string
  logo?: string
  price?: number
  change24h?: number
}

export interface CrossChainAsset {
  symbol: string
  name: string
  logo: string
  totalBalance: string
  totalValueUSD: string
  chains: {
    chainId: number
    chainName: string
    balance: string
    valueUSD: string
    tokenAddress: Address
  }[]
}

export interface ChainSwitchRequest {
  chainId: number
  reason: 'user' | 'dapp' | 'auto'
  timestamp: number
  success: boolean
  error?: string
}

// Enhanced chain configurations
export const SUPPORTED_CHAINS: Record<number, ChainConfig> = {
  [mainnet.id]: {
    ...mainnet,
    logo: '/chains/ethereum.svg',
    color: '#627EEA',
    blockExplorer: 'https://etherscan.io',
    nativeCurrency: {
      ...mainnet.nativeCurrency,
      logo: '/tokens/eth.svg'
    },
    rpcUrls: {
      default: { http: ['https://eth.llamarpc.com'] },
      public: { http: ['https://eth.llamarpc.com'] },
      infura: { http: [`https://mainnet.infura.io/v3/${process.env.NEXT_PUBLIC_INFURA_KEY}`] },
      alchemy: { http: [`https://eth-mainnet.g.alchemy.com/v2/${process.env.NEXT_PUBLIC_ALCHEMY_KEY}`] }
    },
    features: ['DeFi Hub', 'NFTs', 'Staking', 'Layer 2s'],
    gasPrice: { slow: '20', standard: '25', fast: '35' },
    bridgeSupport: ['Polygon', 'Arbitrum', 'Optimism', 'Base'],
    dexes: ['Uniswap', 'SushiSwap', '1inch', 'Curve'],
    lending: ['Aave', 'Compound', 'MakerDAO'],
    category: 'mainnet',
    tvl: '$50B+',
    dailyVolume: '$2B+'
  },
  [polygon.id]: {
    ...polygon,
    logo: '/chains/polygon.svg',
    color: '#8247E5',
    blockExplorer: 'https://polygonscan.com',
    nativeCurrency: {
      ...polygon.nativeCurrency,
      logo: '/tokens/matic.svg'
    },
    rpcUrls: {
      default: { http: ['https://polygon.llamarpc.com'] },
      public: { http: ['https://polygon.llamarpc.com'] },
      infura: { http: [`https://polygon-mainnet.infura.io/v3/${process.env.NEXT_PUBLIC_INFURA_KEY}`] },
      alchemy: { http: [`https://polygon-mainnet.g.alchemy.com/v2/${process.env.NEXT_PUBLIC_ALCHEMY_KEY}`] }
    },
    features: ['Low Fees', 'Fast Transactions', 'Ethereum Compatible'],
    gasPrice: { slow: '30', standard: '35', fast: '50' },
    bridgeSupport: ['Ethereum', 'BSC', 'Avalanche'],
    dexes: ['QuickSwap', 'SushiSwap', '1inch'],
    lending: ['Aave', 'Compound'],
    category: 'sidechain',
    tvl: '$1B+',
    dailyVolume: '$100M+'
  },
  [arbitrum.id]: {
    ...arbitrum,
    logo: '/chains/arbitrum.svg',
    color: '#28A0F0',
    blockExplorer: 'https://arbiscan.io',
    nativeCurrency: {
      ...arbitrum.nativeCurrency,
      logo: '/tokens/eth.svg'
    },
    rpcUrls: {
      default: { http: ['https://arbitrum.llamarpc.com'] },
      public: { http: ['https://arbitrum.llamarpc.com'] },
      infura: { http: [`https://arbitrum-mainnet.infura.io/v3/${process.env.NEXT_PUBLIC_INFURA_KEY}`] },
      alchemy: { http: [`https://arb-mainnet.g.alchemy.com/v2/${process.env.NEXT_PUBLIC_ALCHEMY_KEY}`] }
    },
    features: ['Layer 2', 'Low Fees', 'Ethereum Security'],
    gasPrice: { slow: '0.1', standard: '0.15', fast: '0.25' },
    bridgeSupport: ['Ethereum', 'Polygon'],
    dexes: ['Uniswap V3', 'SushiSwap', 'Camelot'],
    lending: ['Aave', 'Radiant'],
    category: 'l2',
    tvl: '$2B+',
    dailyVolume: '$200M+'
  },
  [optimism.id]: {
    ...optimism,
    logo: '/chains/optimism.svg',
    color: '#FF0420',
    blockExplorer: 'https://optimistic.etherscan.io',
    nativeCurrency: {
      ...optimism.nativeCurrency,
      logo: '/tokens/eth.svg'
    },
    rpcUrls: {
      default: { http: ['https://optimism.llamarpc.com'] },
      public: { http: ['https://optimism.llamarpc.com'] },
      infura: { http: [`https://optimism-mainnet.infura.io/v3/${process.env.NEXT_PUBLIC_INFURA_KEY}`] },
      alchemy: { http: [`https://opt-mainnet.g.alchemy.com/v2/${process.env.NEXT_PUBLIC_ALCHEMY_KEY}`] }
    },
    features: ['Layer 2', 'Optimistic Rollups', 'Low Fees'],
    gasPrice: { slow: '0.001', standard: '0.002', fast: '0.005' },
    bridgeSupport: ['Ethereum', 'Base'],
    dexes: ['Uniswap V3', 'Velodrome', '1inch'],
    lending: ['Aave', 'Exactly'],
    category: 'l2',
    tvl: '$1B+',
    dailyVolume: '$50M+'
  },
  [base.id]: {
    ...base,
    logo: '/chains/base.svg',
    color: '#0052FF',
    blockExplorer: 'https://basescan.org',
    nativeCurrency: {
      ...base.nativeCurrency,
      logo: '/tokens/eth.svg'
    },
    rpcUrls: {
      default: { http: ['https://base.llamarpc.com'] },
      public: { http: ['https://base.llamarpc.com'] },
      alchemy: { http: [`https://base-mainnet.g.alchemy.com/v2/${process.env.NEXT_PUBLIC_ALCHEMY_KEY}`] }
    },
    features: ['Coinbase L2', 'Low Fees', 'Growing Ecosystem'],
    gasPrice: { slow: '0.001', standard: '0.002', fast: '0.005' },
    bridgeSupport: ['Ethereum', 'Optimism'],
    dexes: ['Uniswap V3', 'Aerodrome'],
    lending: ['Aave', 'Compound'],
    category: 'l2',
    tvl: '$500M+',
    dailyVolume: '$30M+'
  },
  [bsc.id]: {
    ...bsc,
    logo: '/chains/bsc.svg',
    color: '#F3BA2F',
    blockExplorer: 'https://bscscan.com',
    nativeCurrency: {
      ...bsc.nativeCurrency,
      logo: '/tokens/bnb.svg'
    },
    rpcUrls: {
      default: { http: ['https://bsc-dataseed.binance.org'] },
      public: { http: ['https://bsc-dataseed.binance.org'] }
    },
    features: ['Low Fees', 'Fast Transactions', 'Large Ecosystem'],
    gasPrice: { slow: '3', standard: '5', fast: '10' },
    bridgeSupport: ['Ethereum', 'Polygon', 'Avalanche'],
    dexes: ['PancakeSwap', 'Uniswap V3', '1inch'],
    lending: ['Venus', 'Alpaca'],
    category: 'sidechain',
    tvl: '$3B+',
    dailyVolume: '$300M+'
  },
  [avalanche.id]: {
    ...avalanche,
    logo: '/chains/avalanche.svg',
    color: '#E84142',
    blockExplorer: 'https://snowtrace.io',
    nativeCurrency: {
      ...avalanche.nativeCurrency,
      logo: '/tokens/avax.svg'
    },
    rpcUrls: {
      default: { http: ['https://api.avax.network/ext/bc/C/rpc'] },
      public: { http: ['https://api.avax.network/ext/bc/C/rpc'] }
    },
    features: ['Fast Finality', 'Low Fees', 'Subnets'],
    gasPrice: { slow: '25', standard: '30', fast: '40' },
    bridgeSupport: ['Ethereum', 'BSC', 'Polygon'],
    dexes: ['Trader Joe', 'Pangolin', 'SushiSwap'],
    lending: ['Aave', 'Benqi'],
    category: 'mainnet',
    tvl: '$800M+',
    dailyVolume: '$80M+'
  }
}

// Add testnet configurations for development
if (process.env.NODE_ENV === 'development') {
  SUPPORTED_CHAINS[sepolia.id] = {
    ...sepolia,
    logo: '/chains/ethereum.svg',
    color: '#627EEA',
    blockExplorer: 'https://sepolia.etherscan.io',
    nativeCurrency: { ...sepolia.nativeCurrency, logo: '/tokens/eth.svg' },
    rpcUrls: {
      default: { http: ['https://sepolia.infura.io/v3/'] },
      public: { http: ['https://sepolia.infura.io/v3/'] }
    },
    features: ['Testnet', 'Development'],
    gasPrice: { slow: '1', standard: '2', fast: '5' },
    bridgeSupport: [],
    dexes: [],
    lending: [],
    category: 'testnet'
  }
}

export class MultiChainManager {
  private static instance: MultiChainManager
  private balanceCache: Map<string, ChainBalance> = new Map()
  private priceCache: Map<string, number> = new Map()
  private switchHistory: ChainSwitchRequest[] = []
  private updateInterval: NodeJS.Timeout | null = null

  private constructor() {
    this.startBalanceUpdates()
  }

  static getInstance(): MultiChainManager {
    if (!MultiChainManager.instance) {
      MultiChainManager.instance = new MultiChainManager()
    }
    return MultiChainManager.instance
  }

  /**
   * Get all supported chains
   */
  getSupportedChains(): ChainConfig[] {
    return Object.values(SUPPORTED_CHAINS)
  }

  /**
   * Get chains by category
   */
  getChainsByCategory(category: ChainConfig['category']): ChainConfig[] {
    return Object.values(SUPPORTED_CHAINS).filter(chain => chain.category === category)
  }

  /**
   * Get chain configuration by ID
   */
  getChainConfig(chainId: number): ChainConfig | undefined {
    return SUPPORTED_CHAINS[chainId]
  }

  /**
   * Get mainnet chains only
   */
  getMainnetChains(): ChainConfig[] {
    return this.getChainsByCategory('mainnet').concat(
      this.getChainsByCategory('l2'),
      this.getChainsByCategory('sidechain')
    )
  }

  /**
   * Get testnet chains
   */
  getTestnetChains(): ChainConfig[] {
    return this.getChainsByCategory('testnet')
  }

  /**
   * Check if chain switching is supported
   */
  isChainSwitchSupported(fromChainId: number, toChainId: number): boolean {
    const fromChain = this.getChainConfig(fromChainId)
    const toChain = this.getChainConfig(toChainId)
    
    if (!fromChain || !toChain) return false
    
    // Check if chains support bridging to each other
    return fromChain.bridgeSupport.includes(toChain.name) ||
           toChain.bridgeSupport.includes(fromChain.name)
  }

  /**
   * Get recommended chains for a user
   */
  getRecommendedChains(userActivity?: { chainId: number; txCount: number }[]): ChainConfig[] {
    const mainnetChains = this.getMainnetChains()
    
    if (!userActivity || userActivity.length === 0) {
      // Default recommendations
      return [
        SUPPORTED_CHAINS[mainnet.id],
        SUPPORTED_CHAINS[polygon.id],
        SUPPORTED_CHAINS[arbitrum.id],
        SUPPORTED_CHAINS[optimism.id]
      ].filter(Boolean)
    }
    
    // Sort by user activity
    const sortedByActivity = userActivity
      .sort((a, b) => b.txCount - a.txCount)
      .map(activity => SUPPORTED_CHAINS[activity.chainId])
      .filter(Boolean)
    
    // Fill remaining slots with popular chains
    const remaining = mainnetChains.filter(
      chain => !sortedByActivity.find(active => active.id === chain.id)
    )
    
    return [...sortedByActivity, ...remaining].slice(0, 6)
  }

  /**
   * Get gas price recommendations for a chain
   */
  getGasPriceRecommendations(chainId: number): { slow: string; standard: string; fast: string } | null {
    const chain = this.getChainConfig(chainId)
    return chain?.gasPrice || null
  }

  /**
   * Start automatic balance updates
   */
  private startBalanceUpdates(): void {
    // Update balances every 30 seconds
    this.updateInterval = setInterval(() => {
      this.updateAllBalances()
    }, 30000)
  }

  /**
   * Stop automatic balance updates
   */
  stopBalanceUpdates(): void {
    if (this.updateInterval) {
      clearInterval(this.updateInterval)
      this.updateInterval = null
    }
  }

  /**
   * Update balances for all chains
   */
  private async updateAllBalances(): Promise<void> {
    // This would integrate with actual blockchain APIs
    // For now, this is a placeholder
    console.log('Updating balances for all chains...')
  }

  /**
   * Get cached balance for a chain
   */
  getChainBalance(chainId: number, address: Address): ChainBalance | null {
    const key = `${chainId}-${address}`
    return this.balanceCache.get(key) || null
  }

  /**
   * Get aggregated balance across all chains
   */
  getAggregatedBalance(address: Address): {
    totalValueUSD: string
    chainBalances: ChainBalance[]
    crossChainAssets: CrossChainAsset[]
  } {
    const chainBalances: ChainBalance[] = []
    let totalValueUSD = 0

    // Get balances from all supported chains
    Object.keys(SUPPORTED_CHAINS).forEach(chainIdStr => {
      const chainId = parseInt(chainIdStr)
      const balance = this.getChainBalance(chainId, address)
      if (balance) {
        chainBalances.push(balance)
        totalValueUSD += parseFloat(balance.totalValueUSD)
      }
    })

    // Aggregate cross-chain assets
    const crossChainAssets = this.aggregateCrossChainAssets(chainBalances)

    return {
      totalValueUSD: totalValueUSD.toString(),
      chainBalances,
      crossChainAssets
    }
  }

  /**
   * Aggregate assets across chains
   */
  private aggregateCrossChainAssets(chainBalances: ChainBalance[]): CrossChainAsset[] {
    const assetMap = new Map<string, CrossChainAsset>()

    chainBalances.forEach(chainBalance => {
      chainBalance.tokens.forEach(token => {
        const key = token.symbol.toLowerCase()
        
        if (!assetMap.has(key)) {
          assetMap.set(key, {
            symbol: token.symbol,
            name: token.name,
            logo: token.logo || '',
            totalBalance: '0',
            totalValueUSD: '0',
            chains: []
          })
        }

        const asset = assetMap.get(key)!
        asset.chains.push({
          chainId: chainBalance.chainId,
          chainName: chainBalance.chainName,
          balance: token.balance,
          valueUSD: token.valueUSD,
          tokenAddress: token.address
        })

        // Update totals
        const currentTotal = parseFloat(asset.totalBalance)
        const tokenBalance = parseFloat(token.balance)
        asset.totalBalance = (currentTotal + tokenBalance).toString()

        const currentValueUSD = parseFloat(asset.totalValueUSD)
        const tokenValueUSD = parseFloat(token.valueUSD)
        asset.totalValueUSD = (currentValueUSD + tokenValueUSD).toString()
      })
    })

    return Array.from(assetMap.values())
      .filter(asset => parseFloat(asset.totalValueUSD) > 0)
      .sort((a, b) => parseFloat(b.totalValueUSD) - parseFloat(a.totalValueUSD))
  }

  /**
   * Record chain switch attempt
   */
  recordChainSwitch(request: ChainSwitchRequest): void {
    this.switchHistory.push(request)
    
    // Keep only last 100 switches
    if (this.switchHistory.length > 100) {
      this.switchHistory = this.switchHistory.slice(-100)
    }
  }

  /**
   * Get chain switch statistics
   */
  getChainSwitchStats(): {
    totalSwitches: number
    successRate: number
    mostSwitchedTo: number | null
    recentSwitches: ChainSwitchRequest[]
  } {
    const totalSwitches = this.switchHistory.length
    const successfulSwitches = this.switchHistory.filter(s => s.success).length
    const successRate = totalSwitches > 0 ? (successfulSwitches / totalSwitches) * 100 : 0

    // Find most switched to chain
    const chainCounts = new Map<number, number>()
    this.switchHistory.forEach(switch_ => {
      const count = chainCounts.get(switch_.chainId) || 0
      chainCounts.set(switch_.chainId, count + 1)
    })

    const mostSwitchedTo = chainCounts.size > 0 
      ? Array.from(chainCounts.entries()).sort(([,a], [,b]) => b - a)[0][0]
      : null

    return {
      totalSwitches,
      successRate,
      mostSwitchedTo,
      recentSwitches: this.switchHistory.slice(-10)
    }
  }

  /**
   * Clear cache
   */
  clearCache(): void {
    this.balanceCache.clear()
    this.priceCache.clear()
  }

  /**
   * Cleanup resources
   */
  cleanup(): void {
    this.stopBalanceUpdates()
    this.clearCache()
  }
}

// Export singleton instance
export const multiChainManager = MultiChainManager.getInstance()
