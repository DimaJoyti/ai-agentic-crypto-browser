import { createPublicClient, http, formatGwei, parseGwei, type Hash } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum GasPriority {
  SLOW = 'slow',
  STANDARD = 'standard',
  FAST = 'fast',
  INSTANT = 'instant'
}

export interface GasEstimate {
  priority: GasPriority
  gasPrice: bigint
  maxFeePerGas?: bigint
  maxPriorityFeePerGas?: bigint
  estimatedTime: number // seconds
  cost: string // in ETH
  confidence: number // 0-100
}

export interface GasOptimizationSuggestion {
  type: 'timing' | 'batching' | 'route' | 'token'
  title: string
  description: string
  potentialSavings: string
  difficulty: 'easy' | 'medium' | 'hard'
  action?: () => void
}

export interface TransactionBatch {
  id: string
  transactions: any[]
  estimatedGasSavings: string
  totalGasCost: string
  status: 'pending' | 'ready' | 'executing' | 'completed' | 'failed'
  createdAt: number
}

export interface GasTracker {
  chainId: number
  currentGasPrice: bigint
  trend: 'rising' | 'falling' | 'stable'
  networkCongestion: 'low' | 'medium' | 'high'
  recommendedAction: 'wait' | 'proceed' | 'urgent'
  nextUpdateIn: number
}

export class GasOptimizer {
  private static instance: GasOptimizer
  private clients: Map<number, any> = new Map()
  private gasHistory: Map<number, bigint[]> = new Map()
  private trackers: Map<number, GasTracker> = new Map()
  private updateIntervals: Map<number, NodeJS.Timeout> = new Map()

  private constructor() {
    this.initializeClients()
    this.startGasTracking()
  }

  static getInstance(): GasOptimizer {
    if (!GasOptimizer.instance) {
      GasOptimizer.instance = new GasOptimizer()
    }
    return GasOptimizer.instance
  }

  private initializeClients() {
    Object.values(SUPPORTED_CHAINS).forEach(chain => {
      if (!chain.isTestnet || chain.id === 11155111) {
        try {
          const client = createPublicClient({
            chain: {
              id: chain.id,
              name: chain.name,
              network: chain.shortName.toLowerCase(),
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: chain.rpcUrls
            } as any,
            transport: http()
          })
          this.clients.set(chain.id, client)
        } catch (error) {
          console.warn(`Failed to initialize gas client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private startGasTracking() {
    this.clients.forEach((client, chainId) => {
      this.updateGasPrice(chainId)
      
      // Update gas prices every 15 seconds
      const interval = setInterval(() => {
        this.updateGasPrice(chainId)
      }, 15000)
      
      this.updateIntervals.set(chainId, interval)
    })
  }

  private async updateGasPrice(chainId: number) {
    const client = this.clients.get(chainId)
    if (!client) return

    try {
      const gasPrice = await client.getGasPrice()
      
      // Update history
      if (!this.gasHistory.has(chainId)) {
        this.gasHistory.set(chainId, [])
      }
      
      const history = this.gasHistory.get(chainId)!
      history.push(gasPrice)
      
      // Keep only last 20 readings (5 minutes of history)
      if (history.length > 20) {
        history.shift()
      }

      // Calculate trend
      const trend = this.calculateTrend(history)
      const congestion = this.calculateCongestion(gasPrice, chainId)
      const recommendation = this.getRecommendation(trend, congestion)

      // Update tracker
      this.trackers.set(chainId, {
        chainId,
        currentGasPrice: gasPrice,
        trend,
        networkCongestion: congestion,
        recommendedAction: recommendation,
        nextUpdateIn: 15
      })

    } catch (error) {
      console.warn(`Failed to update gas price for chain ${chainId}:`, error)
    }
  }

  private calculateTrend(history: bigint[]): 'rising' | 'falling' | 'stable' {
    if (history.length < 3) return 'stable'

    const recent = history.slice(-3)
    const older = history.slice(-6, -3)

    if (older.length === 0) return 'stable'

    const recentAvg = recent.reduce((sum, price) => sum + price, BigInt(0)) / BigInt(recent.length)
    const olderAvg = older.reduce((sum, price) => sum + price, BigInt(0)) / BigInt(older.length)

    const change = Number(recentAvg - olderAvg) / Number(olderAvg)

    if (change > 0.1) return 'rising'
    if (change < -0.1) return 'falling'
    return 'stable'
  }

  private calculateCongestion(gasPrice: bigint, chainId: number): 'low' | 'medium' | 'high' {
    // Base gas prices for different chains (in gwei)
    const basePrices: Record<number, number> = {
      1: 20,    // Ethereum
      137: 30,  // Polygon
      42161: 0.1, // Arbitrum
      10: 0.001,  // Optimism
      8453: 0.001, // Base
      43114: 25,  // Avalanche
      56: 5,      // BSC
      250: 20,    // Fantom
      100: 2,     // Gnosis
      11155111: 20 // Sepolia
    }

    const basePrice = basePrices[chainId] || 20
    const currentGwei = Number(formatGwei(gasPrice))
    const ratio = currentGwei / basePrice

    if (ratio > 3) return 'high'
    if (ratio > 1.5) return 'medium'
    return 'low'
  }

  private getRecommendation(
    trend: 'rising' | 'falling' | 'stable',
    congestion: 'low' | 'medium' | 'high'
  ): 'wait' | 'proceed' | 'urgent' {
    if (congestion === 'high' && trend === 'rising') return 'wait'
    if (congestion === 'low' && trend === 'falling') return 'urgent'
    if (congestion === 'medium' && trend === 'stable') return 'proceed'
    if (trend === 'falling') return 'proceed'
    return 'wait'
  }

  async getGasEstimates(chainId: number, gasLimit: bigint): Promise<GasEstimate[]> {
    const client = this.clients.get(chainId)
    if (!client) {
      throw new Error(`No client available for chain ${chainId}`)
    }

    try {
      const gasPrice = await client.getGasPrice()
      const chain = SUPPORTED_CHAINS[chainId]
      
      // Check if chain supports EIP-1559
      const supportsEIP1559 = chainId === 1 || chainId === 137 || chainId === 42161 || chainId === 10 || chainId === 8453

      if (supportsEIP1559) {
        return this.getEIP1559Estimates(gasPrice, gasLimit, chainId)
      } else {
        return this.getLegacyEstimates(gasPrice, gasLimit, chainId)
      }
    } catch (error) {
      console.error(`Failed to get gas estimates for chain ${chainId}:`, error)
      throw error
    }
  }

  private getEIP1559Estimates(baseGasPrice: bigint, gasLimit: bigint, chainId: number): GasEstimate[] {
    const baseFee = baseGasPrice
    
    const estimates: GasEstimate[] = [
      {
        priority: GasPriority.SLOW,
        gasPrice: baseFee,
        maxFeePerGas: baseFee + parseGwei('1'),
        maxPriorityFeePerGas: parseGwei('1'),
        estimatedTime: 300, // 5 minutes
        cost: this.calculateCost(baseFee + parseGwei('1'), gasLimit),
        confidence: 70
      },
      {
        priority: GasPriority.STANDARD,
        gasPrice: baseFee + parseGwei('2'),
        maxFeePerGas: baseFee + parseGwei('3'),
        maxPriorityFeePerGas: parseGwei('2'),
        estimatedTime: 120, // 2 minutes
        cost: this.calculateCost(baseFee + parseGwei('3'), gasLimit),
        confidence: 85
      },
      {
        priority: GasPriority.FAST,
        gasPrice: baseFee + parseGwei('5'),
        maxFeePerGas: baseFee + parseGwei('7'),
        maxPriorityFeePerGas: parseGwei('5'),
        estimatedTime: 60, // 1 minute
        cost: this.calculateCost(baseFee + parseGwei('7'), gasLimit),
        confidence: 95
      },
      {
        priority: GasPriority.INSTANT,
        gasPrice: baseFee + parseGwei('10'),
        maxFeePerGas: baseFee + parseGwei('15'),
        maxPriorityFeePerGas: parseGwei('10'),
        estimatedTime: 15, // 15 seconds
        cost: this.calculateCost(baseFee + parseGwei('15'), gasLimit),
        confidence: 99
      }
    ]

    return estimates
  }

  private getLegacyEstimates(gasPrice: bigint, gasLimit: bigint, chainId: number): GasEstimate[] {
    const estimates: GasEstimate[] = [
      {
        priority: GasPriority.SLOW,
        gasPrice: gasPrice * BigInt(80) / BigInt(100), // 80% of current
        estimatedTime: 300,
        cost: this.calculateCost(gasPrice * BigInt(80) / BigInt(100), gasLimit),
        confidence: 70
      },
      {
        priority: GasPriority.STANDARD,
        gasPrice: gasPrice,
        estimatedTime: 120,
        cost: this.calculateCost(gasPrice, gasLimit),
        confidence: 85
      },
      {
        priority: GasPriority.FAST,
        gasPrice: gasPrice * BigInt(120) / BigInt(100), // 120% of current
        estimatedTime: 60,
        cost: this.calculateCost(gasPrice * BigInt(120) / BigInt(100), gasLimit),
        confidence: 95
      },
      {
        priority: GasPriority.INSTANT,
        gasPrice: gasPrice * BigInt(150) / BigInt(100), // 150% of current
        estimatedTime: 15,
        cost: this.calculateCost(gasPrice * BigInt(150) / BigInt(100), gasLimit),
        confidence: 99
      }
    ]

    return estimates
  }

  private calculateCost(gasPrice: bigint, gasLimit: bigint): string {
    const totalCost = gasPrice * gasLimit
    const ethCost = Number(totalCost) / 1e18
    return ethCost.toFixed(6)
  }

  getGasTracker(chainId: number): GasTracker | undefined {
    return this.trackers.get(chainId)
  }

  getAllGasTrackers(): GasTracker[] {
    return Array.from(this.trackers.values())
  }

  generateOptimizationSuggestions(
    chainId: number,
    transactionType: string,
    amount?: string
  ): GasOptimizationSuggestion[] {
    const tracker = this.trackers.get(chainId)
    const suggestions: GasOptimizationSuggestion[] = []

    if (!tracker) return suggestions

    // Timing suggestions
    if (tracker.networkCongestion === 'high' && tracker.trend === 'rising') {
      suggestions.push({
        type: 'timing',
        title: 'Wait for Lower Gas Prices',
        description: 'Network congestion is high and gas prices are rising. Consider waiting 30-60 minutes.',
        potentialSavings: '20-40%',
        difficulty: 'easy'
      })
    }

    if (tracker.trend === 'falling' && tracker.networkCongestion === 'medium') {
      suggestions.push({
        type: 'timing',
        title: 'Good Time to Transact',
        description: 'Gas prices are falling and network congestion is moderate. Good time to proceed.',
        potentialSavings: '10-20%',
        difficulty: 'easy'
      })
    }

    // Batching suggestions
    if (transactionType === 'multiple') {
      suggestions.push({
        type: 'batching',
        title: 'Batch Multiple Transactions',
        description: 'Combine multiple transactions into a single batch to save on gas costs.',
        potentialSavings: '30-50%',
        difficulty: 'medium'
      })
    }

    // Route optimization for swaps
    if (transactionType === 'swap' && chainId === 1) {
      suggestions.push({
        type: 'route',
        title: 'Optimize Swap Route',
        description: 'Use aggregators like 1inch or Paraswap to find the most gas-efficient route.',
        potentialSavings: '15-25%',
        difficulty: 'easy'
      })
    }

    // Layer 2 suggestions
    if (chainId === 1 && amount && parseFloat(amount) < 1000) {
      suggestions.push({
        type: 'route',
        title: 'Consider Layer 2',
        description: 'For smaller amounts, consider using Arbitrum, Optimism, or Polygon for lower fees.',
        potentialSavings: '80-95%',
        difficulty: 'medium'
      })
    }

    return suggestions
  }

  destroy(): void {
    this.updateIntervals.forEach(interval => clearInterval(interval))
    this.updateIntervals.clear()
    this.trackers.clear()
    this.gasHistory.clear()
  }
}

// Export singleton instance
export const gasOptimizer = GasOptimizer.getInstance()
