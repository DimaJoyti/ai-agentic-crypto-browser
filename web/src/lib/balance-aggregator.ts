import { type Address } from 'viem'
import { 
  type ChainBalance, 
  type TokenBalance, 
  type CrossChainAsset,
  multiChainManager,
  SUPPORTED_CHAINS 
} from '@/lib/multi-chain-manager'

export interface BalanceUpdateOptions {
  forceRefresh?: boolean
  includeTestnets?: boolean
  minValueUSD?: number
  timeout?: number
}

export interface PriceData {
  symbol: string
  price: number
  change24h: number
  marketCap?: number
  volume24h?: number
  lastUpdated: number
}

export interface PortfolioSummary {
  totalValueUSD: string
  totalChange24h: string
  totalChange24hPercent: number
  chainCount: number
  tokenCount: number
  topAssets: CrossChainAsset[]
  chainDistribution: {
    chainId: number
    chainName: string
    valueUSD: string
    percentage: number
    color: string
  }[]
}

export interface BalanceHistory {
  timestamp: number
  totalValueUSD: string
  chainBalances: { chainId: number; valueUSD: string }[]
}

export class BalanceAggregator {
  private static instance: BalanceAggregator
  private balanceCache = new Map<string, ChainBalance>()
  private priceCache = new Map<string, PriceData>()
  private historyCache = new Map<string, BalanceHistory[]>()
  private updatePromises = new Map<string, Promise<ChainBalance>>()
  private lastUpdateTime = new Map<string, number>()
  
  // Cache durations
  private readonly BALANCE_CACHE_DURATION = 30000 // 30 seconds
  private readonly PRICE_CACHE_DURATION = 60000 // 1 minute
  private readonly HISTORY_CACHE_DURATION = 300000 // 5 minutes

  private constructor() {}

  static getInstance(): BalanceAggregator {
    if (!BalanceAggregator.instance) {
      BalanceAggregator.instance = new BalanceAggregator()
    }
    return BalanceAggregator.instance
  }

  /**
   * Get aggregated balances for an address across all chains
   */
  async getAggregatedBalances(
    address: Address, 
    options: BalanceUpdateOptions = {}
  ): Promise<{
    totalValueUSD: string
    chainBalances: ChainBalance[]
    crossChainAssets: CrossChainAsset[]
    lastUpdated: number
  }> {
    const { 
      forceRefresh = false, 
      includeTestnets = false, 
      minValueUSD = 0.01 
    } = options

    const chains = includeTestnets 
      ? multiChainManager.getSupportedChains()
      : multiChainManager.getMainnetChains()

    // Get balances for all chains
    const balancePromises = chains.map(chain => 
      this.getChainBalance(address, chain.id, forceRefresh)
    )

    const chainBalances = (await Promise.allSettled(balancePromises))
      .map(result => result.status === 'fulfilled' ? result.value : null)
      .filter((balance): balance is ChainBalance => 
        balance !== null && parseFloat(balance.totalValueUSD) >= minValueUSD
      )

    // Calculate total value
    const totalValueUSD = chainBalances
      .reduce((sum, balance) => sum + parseFloat(balance.totalValueUSD), 0)
      .toString()

    // Aggregate cross-chain assets
    const crossChainAssets = this.aggregateCrossChainAssets(chainBalances)

    return {
      totalValueUSD,
      chainBalances,
      crossChainAssets,
      lastUpdated: Date.now()
    }
  }

  /**
   * Get balance for a specific chain
   */
  async getChainBalance(
    address: Address, 
    chainId: number, 
    forceRefresh = false
  ): Promise<ChainBalance> {
    const cacheKey = `${address}-${chainId}`
    const lastUpdate = this.lastUpdateTime.get(cacheKey) || 0
    const now = Date.now()

    // Check cache first
    if (!forceRefresh && now - lastUpdate < this.BALANCE_CACHE_DURATION) {
      const cached = this.balanceCache.get(cacheKey)
      if (cached) return cached
    }

    // Check if update is already in progress
    const existingPromise = this.updatePromises.get(cacheKey)
    if (existingPromise) return existingPromise

    // Start new update
    const updatePromise = this.fetchChainBalance(address, chainId)
    this.updatePromises.set(cacheKey, updatePromise)

    try {
      const balance = await updatePromise
      this.balanceCache.set(cacheKey, balance)
      this.lastUpdateTime.set(cacheKey, now)
      return balance
    } finally {
      this.updatePromises.delete(cacheKey)
    }
  }

  /**
   * Fetch balance from blockchain (placeholder implementation)
   */
  private async fetchChainBalance(address: Address, chainId: number): Promise<ChainBalance> {
    const chain = SUPPORTED_CHAINS[chainId]
    if (!chain) {
      throw new Error(`Unsupported chain: ${chainId}`)
    }

    // This is a placeholder implementation
    // In a real app, you would use actual blockchain APIs like:
    // - Alchemy, Infura, or QuickNode for RPC calls
    // - Moralis, Covalent, or similar for multi-chain data
    // - DEX APIs for token prices

    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 100 + Math.random() * 500))

    // Mock data for demonstration
    const mockTokens: TokenBalance[] = [
      {
        address: '0x0000000000000000000000000000000000000000' as Address,
        symbol: chain.nativeCurrency.symbol,
        name: chain.nativeCurrency.name,
        decimals: chain.nativeCurrency.decimals,
        balance: (Math.random() * 10).toFixed(6),
        balanceFormatted: '',
        valueUSD: (Math.random() * 1000).toFixed(2),
        logo: chain.nativeCurrency.logo,
        price: Math.random() * 3000,
        change24h: (Math.random() - 0.5) * 20
      }
    ]

    // Add some random ERC-20 tokens for mainnet chains
    if (chain.category !== 'testnet') {
      const commonTokens = [
        { symbol: 'USDC', name: 'USD Coin', decimals: 6 },
        { symbol: 'USDT', name: 'Tether USD', decimals: 6 },
        { symbol: 'WETH', name: 'Wrapped Ether', decimals: 18 }
      ]

      commonTokens.forEach(token => {
        if (Math.random() > 0.5) { // 50% chance to have each token
          const balance = Math.random() * 1000
          mockTokens.push({
            address: `0x${Math.random().toString(16).slice(2, 42)}` as Address,
            symbol: token.symbol,
            name: token.name,
            decimals: token.decimals,
            balance: balance.toFixed(token.decimals),
            balanceFormatted: balance.toFixed(2),
            valueUSD: (balance * (token.symbol.includes('USD') ? 1 : Math.random() * 3000)).toFixed(2),
            logo: `/tokens/${token.symbol.toLowerCase()}.svg`,
            price: token.symbol.includes('USD') ? 1 : Math.random() * 3000,
            change24h: (Math.random() - 0.5) * 20
          })
        }
      })
    }

    // Format balances
    mockTokens.forEach(token => {
      if (!token.balanceFormatted) {
        const balance = parseFloat(token.balance)
        token.balanceFormatted = balance < 0.01 ? balance.toExponential(2) : balance.toFixed(4)
      }
    })

    const totalValueUSD = mockTokens
      .reduce((sum, token) => sum + parseFloat(token.valueUSD), 0)
      .toFixed(2)

    return {
      chainId,
      chainName: chain.name,
      nativeBalance: mockTokens[0].balance,
      nativeBalanceUSD: mockTokens[0].valueUSD,
      tokens: mockTokens,
      totalValueUSD,
      lastUpdated: Date.now()
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
        
        // Check if this chain already exists for this asset
        const existingChain = asset.chains.find(c => c.chainId === chainBalance.chainId)
        if (existingChain) {
          // Update existing chain data
          const currentBalance = parseFloat(existingChain.balance)
          const tokenBalance = parseFloat(token.balance)
          existingChain.balance = (currentBalance + tokenBalance).toString()
          
          const currentValueUSD = parseFloat(existingChain.valueUSD)
          const tokenValueUSD = parseFloat(token.valueUSD)
          existingChain.valueUSD = (currentValueUSD + tokenValueUSD).toString()
        } else {
          // Add new chain
          asset.chains.push({
            chainId: chainBalance.chainId,
            chainName: chainBalance.chainName,
            balance: token.balance,
            valueUSD: token.valueUSD,
            tokenAddress: token.address
          })
        }

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
   * Get portfolio summary
   */
  async getPortfolioSummary(address: Address): Promise<PortfolioSummary> {
    const { totalValueUSD, chainBalances, crossChainAssets } = await this.getAggregatedBalances(address)

    // Calculate 24h change (mock data for now)
    const totalChange24h = (parseFloat(totalValueUSD) * (Math.random() - 0.5) * 0.1).toFixed(2)
    const totalChange24hPercent = parseFloat(totalValueUSD) > 0 
      ? (parseFloat(totalChange24h) / parseFloat(totalValueUSD)) * 100 
      : 0

    // Chain distribution
    const chainDistribution = chainBalances.map(balance => {
      const chain = SUPPORTED_CHAINS[balance.chainId]
      return {
        chainId: balance.chainId,
        chainName: balance.chainName,
        valueUSD: balance.totalValueUSD,
        percentage: parseFloat(totalValueUSD) > 0 
          ? (parseFloat(balance.totalValueUSD) / parseFloat(totalValueUSD)) * 100 
          : 0,
        color: chain?.color || '#666666'
      }
    }).sort((a, b) => b.percentage - a.percentage)

    return {
      totalValueUSD,
      totalChange24h,
      totalChange24hPercent,
      chainCount: chainBalances.length,
      tokenCount: crossChainAssets.length,
      topAssets: crossChainAssets.slice(0, 10),
      chainDistribution
    }
  }

  /**
   * Get price data for tokens
   */
  async getTokenPrices(symbols: string[]): Promise<Map<string, PriceData>> {
    const prices = new Map<string, PriceData>()
    const now = Date.now()

    for (const symbol of symbols) {
      const cached = this.priceCache.get(symbol.toLowerCase())
      if (cached && now - cached.lastUpdated < this.PRICE_CACHE_DURATION) {
        prices.set(symbol, cached)
        continue
      }

      // Mock price data (in real implementation, fetch from CoinGecko, CoinMarketCap, etc.)
      const priceData: PriceData = {
        symbol,
        price: Math.random() * 3000,
        change24h: (Math.random() - 0.5) * 20,
        marketCap: Math.random() * 1000000000,
        volume24h: Math.random() * 100000000,
        lastUpdated: now
      }

      this.priceCache.set(symbol.toLowerCase(), priceData)
      prices.set(symbol, priceData)
    }

    return prices
  }

  /**
   * Get balance history for an address
   */
  async getBalanceHistory(
    address: Address, 
    days = 30
  ): Promise<BalanceHistory[]> {
    const cacheKey = `${address}-${days}`
    const cached = this.historyCache.get(cacheKey)
    const now = Date.now()

    if (cached && cached.length > 0) {
      const lastEntry = cached[cached.length - 1]
      if (now - lastEntry.timestamp < this.HISTORY_CACHE_DURATION) {
        return cached
      }
    }

    // Generate mock history data
    const history: BalanceHistory[] = []
    const startTime = now - (days * 24 * 60 * 60 * 1000)
    
    for (let i = 0; i < days; i++) {
      const timestamp = startTime + (i * 24 * 60 * 60 * 1000)
      const baseValue = 1000 + Math.sin(i / 10) * 500 + Math.random() * 200
      
      history.push({
        timestamp,
        totalValueUSD: baseValue.toFixed(2),
        chainBalances: [
          { chainId: 1, valueUSD: (baseValue * 0.6).toFixed(2) },
          { chainId: 137, valueUSD: (baseValue * 0.25).toFixed(2) },
          { chainId: 42161, valueUSD: (baseValue * 0.15).toFixed(2) }
        ]
      })
    }

    this.historyCache.set(cacheKey, history)
    return history
  }

  /**
   * Clear all caches
   */
  clearCache(): void {
    this.balanceCache.clear()
    this.priceCache.clear()
    this.historyCache.clear()
    this.lastUpdateTime.clear()
  }

  /**
   * Clear cache for specific address
   */
  clearAddressCache(address: Address): void {
    const keysToDelete: string[] = []
    
    this.balanceCache.forEach((_, key) => {
      if (key.startsWith(address)) {
        keysToDelete.push(key)
      }
    })
    
    keysToDelete.forEach(key => {
      this.balanceCache.delete(key)
      this.lastUpdateTime.delete(key)
    })
    
    this.historyCache.delete(address)
  }

  /**
   * Get cache statistics
   */
  getCacheStats(): {
    balanceCacheSize: number
    priceCacheSize: number
    historyCacheSize: number
    totalMemoryUsage: string
  } {
    const balanceCacheSize = this.balanceCache.size
    const priceCacheSize = this.priceCache.size
    const historyCacheSize = this.historyCache.size
    
    // Rough memory usage estimation
    const avgBalanceSize = 1000 // bytes
    const avgPriceSize = 200 // bytes
    const avgHistorySize = 5000 // bytes
    
    const totalBytes = 
      (balanceCacheSize * avgBalanceSize) +
      (priceCacheSize * avgPriceSize) +
      (historyCacheSize * avgHistorySize)
    
    const totalMemoryUsage = totalBytes > 1024 * 1024 
      ? `${(totalBytes / (1024 * 1024)).toFixed(2)} MB`
      : `${(totalBytes / 1024).toFixed(2)} KB`

    return {
      balanceCacheSize,
      priceCacheSize,
      historyCacheSize,
      totalMemoryUsage
    }
  }
}

// Export singleton instance
export const balanceAggregator = BalanceAggregator.getInstance()

// Hook for using balance aggregator in React components
export const useBalanceAggregator = () => {
  const getAggregatedBalances = (address: Address, options?: BalanceUpdateOptions) =>
    balanceAggregator.getAggregatedBalances(address, options)

  const getChainBalance = (address: Address, chainId: number, forceRefresh?: boolean) =>
    balanceAggregator.getChainBalance(address, chainId, forceRefresh)

  const getPortfolioSummary = (address: Address) =>
    balanceAggregator.getPortfolioSummary(address)

  const getTokenPrices = (symbols: string[]) =>
    balanceAggregator.getTokenPrices(symbols)

  const getBalanceHistory = (address: Address, days?: number) =>
    balanceAggregator.getBalanceHistory(address, days)

  const clearCache = () => balanceAggregator.clearCache()
  const clearAddressCache = (address: Address) => balanceAggregator.clearAddressCache(address)
  const getCacheStats = () => balanceAggregator.getCacheStats()

  return {
    getAggregatedBalances,
    getChainBalance,
    getPortfolioSummary,
    getTokenPrices,
    getBalanceHistory,
    clearCache,
    clearAddressCache,
    getCacheStats
  }
}
