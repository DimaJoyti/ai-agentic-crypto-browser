import { createPublicClient, http, parseAbi, formatUnits, parseUnits, type Address, type Hash } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum ProtocolType {
  DEX = 'dex',
  LENDING = 'lending',
  STAKING = 'staking',
  YIELD_FARMING = 'yield_farming',
  DERIVATIVES = 'derivatives',
  INSURANCE = 'insurance'
}

export interface DeFiProtocol {
  id: string
  name: string
  type: ProtocolType
  description: string
  website: string
  logo: string
  tvl: string
  apy?: string
  supportedChains: number[]
  contracts: Record<number, ProtocolContracts>
  features: string[]
  riskLevel: 'low' | 'medium' | 'high'
  auditStatus: 'audited' | 'unaudited' | 'partially_audited'
}

export interface ProtocolContracts {
  router?: Address
  factory?: Address
  pool?: Address
  token?: Address
  staking?: Address
  rewards?: Address
  governance?: Address
}

export interface TokenInfo {
  address: Address
  symbol: string
  name: string
  decimals: number
  logoURI?: string
  price?: string
  balance?: string
}

export interface LiquidityPool {
  address: Address
  token0: TokenInfo
  token1: TokenInfo
  fee: number
  liquidity: string
  volume24h: string
  apy: string
  protocol: string
}

export interface LendingPosition {
  protocol: string
  asset: TokenInfo
  supplied: string
  borrowed: string
  supplyApy: string
  borrowApy: string
  healthFactor: string
  collateralRatio: string
}

export interface YieldPosition {
  protocol: string
  pool: LiquidityPool
  stakedAmount: string
  rewards: TokenInfo[]
  apy: string
  lockPeriod?: number
  harvestable: string
}

export class DeFiProtocolService {
  private static instance: DeFiProtocolService
  private clients: Map<number, any> = new Map()
  private protocols: Map<string, DeFiProtocol> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeProtocols()
  }

  static getInstance(): DeFiProtocolService {
    if (!DeFiProtocolService.instance) {
      DeFiProtocolService.instance = new DeFiProtocolService()
    }
    return DeFiProtocolService.instance
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
          console.warn(`Failed to initialize DeFi client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeProtocols() {
    // Uniswap V3
    this.protocols.set('uniswap-v3', {
      id: 'uniswap-v3',
      name: 'Uniswap V3',
      type: ProtocolType.DEX,
      description: 'The most popular decentralized exchange with concentrated liquidity',
      website: 'https://uniswap.org',
      logo: '🦄',
      tvl: '$4.2B',
      supportedChains: [1, 137, 42161, 10, 8453],
      contracts: {
        1: {
          router: '0xE592427A0AEce92De3Edee1F18E0157C05861564' as Address,
          factory: '0x1F98431c8aD98523631AE4a59f267346ea31F984' as Address
        },
        137: {
          router: '0xE592427A0AEce92De3Edee1F18E0157C05861564' as Address,
          factory: '0x1F98431c8aD98523631AE4a59f267346ea31F984' as Address
        }
      },
      features: ['Swapping', 'Liquidity Provision', 'Concentrated Liquidity'],
      riskLevel: 'low',
      auditStatus: 'audited'
    })

    // Aave V3
    this.protocols.set('aave-v3', {
      id: 'aave-v3',
      name: 'Aave V3',
      type: ProtocolType.LENDING,
      description: 'Leading decentralized lending and borrowing protocol',
      website: 'https://aave.com',
      logo: '👻',
      tvl: '$6.8B',
      apy: '3.2%',
      supportedChains: [1, 137, 42161, 10, 43114],
      contracts: {
        1: {
          pool: '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2' as Address
        },
        137: {
          pool: '0x794a61358D6845594F94dc1DB02A252b5b4814aD' as Address
        }
      },
      features: ['Lending', 'Borrowing', 'Flash Loans', 'Rate Switching'],
      riskLevel: 'low',
      auditStatus: 'audited'
    })

    // Compound V3
    this.protocols.set('compound-v3', {
      id: 'compound-v3',
      name: 'Compound V3',
      type: ProtocolType.LENDING,
      description: 'Algorithmic money market protocol with isolated markets',
      website: 'https://compound.finance',
      logo: '🏛️',
      tvl: '$1.8B',
      apy: '2.8%',
      supportedChains: [1, 137, 42161, 8453],
      contracts: {
        1: {
          pool: '0xc3d688B66703497DAA19211EEdff47f25384cdc3' as Address
        }
      },
      features: ['Lending', 'Borrowing', 'Isolated Markets'],
      riskLevel: 'low',
      auditStatus: 'audited'
    })

    // Curve Finance
    this.protocols.set('curve', {
      id: 'curve',
      name: 'Curve Finance',
      type: ProtocolType.DEX,
      description: 'Decentralized exchange optimized for stablecoin trading',
      website: 'https://curve.fi',
      logo: '🌀',
      tvl: '$2.1B',
      supportedChains: [1, 137, 42161, 10, 43114],
      contracts: {
        1: {
          factory: '0xF18056Bbd320E96A48e3Fbf8bC061322531aac99' as Address
        }
      },
      features: ['Stablecoin Swaps', 'Liquidity Mining', 'Governance'],
      riskLevel: 'low',
      auditStatus: 'audited'
    })

    // Lido
    this.protocols.set('lido', {
      id: 'lido',
      name: 'Lido',
      type: ProtocolType.STAKING,
      description: 'Liquid staking solution for Ethereum and other PoS chains',
      website: 'https://lido.fi',
      logo: '🏊',
      tvl: '$14.2B',
      apy: '3.8%',
      supportedChains: [1, 137],
      contracts: {
        1: {
          staking: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
          token: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address
        }
      },
      features: ['Liquid Staking', 'stETH Token', 'Rewards'],
      riskLevel: 'medium',
      auditStatus: 'audited'
    })

    // PancakeSwap
    this.protocols.set('pancakeswap', {
      id: 'pancakeswap',
      name: 'PancakeSwap',
      type: ProtocolType.DEX,
      description: 'Leading DEX on BNB Smart Chain with yield farming',
      website: 'https://pancakeswap.finance',
      logo: '🥞',
      tvl: '$1.4B',
      supportedChains: [56],
      contracts: {
        56: {
          router: '0x10ED43C718714eb63d5aA57B78B54704E256024E' as Address,
          factory: '0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73' as Address
        }
      },
      features: ['Swapping', 'Yield Farming', 'Lottery', 'NFTs'],
      riskLevel: 'medium',
      auditStatus: 'audited'
    })
  }

  // Get all protocols
  getAllProtocols(): DeFiProtocol[] {
    return Array.from(this.protocols.values())
  }

  // Get protocols by type
  getProtocolsByType(type: ProtocolType): DeFiProtocol[] {
    return this.getAllProtocols().filter(protocol => protocol.type === type)
  }

  // Get protocols by chain
  getProtocolsByChain(chainId: number): DeFiProtocol[] {
    return this.getAllProtocols().filter(protocol => 
      protocol.supportedChains.includes(chainId)
    )
  }

  // Get specific protocol
  getProtocol(protocolId: string): DeFiProtocol | undefined {
    return this.protocols.get(protocolId)
  }

  // Uniswap V3 Integration
  async getUniswapPools(chainId: number, tokenA: Address, tokenB: Address): Promise<LiquidityPool[]> {
    const client = this.clients.get(chainId)
    const protocol = this.protocols.get('uniswap-v3')
    
    if (!client || !protocol) {
      throw new Error('Uniswap not supported on this chain')
    }

    // Mock implementation - in real app, would query Uniswap subgraph
    return [
      {
        address: '0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640' as Address,
        token0: {
          address: tokenA,
          symbol: 'USDC',
          name: 'USD Coin',
          decimals: 6,
          price: '1.00'
        },
        token1: {
          address: tokenB,
          symbol: 'WETH',
          name: 'Wrapped Ether',
          decimals: 18,
          price: '2400.00'
        },
        fee: 500, // 0.05%
        liquidity: '12500000',
        volume24h: '45000000',
        apy: '8.5%',
        protocol: 'uniswap-v3'
      }
    ]
  }

  // Aave V3 Integration
  async getAavePositions(chainId: number, userAddress: Address): Promise<LendingPosition[]> {
    const client = this.clients.get(chainId)
    const protocol = this.protocols.get('aave-v3')
    
    if (!client || !protocol) {
      throw new Error('Aave not supported on this chain')
    }

    // Mock implementation - in real app, would query Aave contracts
    return [
      {
        protocol: 'aave-v3',
        asset: {
          address: '0xA0b86a33E6441b8435b662f0E2d0c2837b0b3c0' as Address,
          symbol: 'USDC',
          name: 'USD Coin',
          decimals: 6,
          price: '1.00'
        },
        supplied: '10000',
        borrowed: '5000',
        supplyApy: '3.2%',
        borrowApy: '4.8%',
        healthFactor: '2.1',
        collateralRatio: '75%'
      }
    ]
  }

  // Generic token swap estimation
  async estimateSwap(
    chainId: number,
    tokenIn: Address,
    tokenOut: Address,
    amountIn: string,
    protocol: string = 'uniswap-v3'
  ): Promise<{
    amountOut: string
    priceImpact: string
    gasEstimate: string
    route: string[]
  }> {
    const client = this.clients.get(chainId)
    
    if (!client) {
      throw new Error('Chain not supported')
    }

    // Mock implementation - in real app, would use actual DEX aggregator
    const amountInNum = parseFloat(amountIn)
    const mockRate = 2400 // ETH/USDC rate
    const amountOut = (amountInNum * mockRate * 0.997).toString() // 0.3% fee
    
    return {
      amountOut,
      priceImpact: '0.15%',
      gasEstimate: '150000',
      route: [tokenIn, tokenOut]
    }
  }

  // Execute token swap
  async executeSwap(
    chainId: number,
    tokenIn: Address,
    tokenOut: Address,
    amountIn: string,
    minAmountOut: string,
    userAddress: Address,
    protocol: string = 'uniswap-v3'
  ): Promise<Hash> {
    const client = this.clients.get(chainId)
    const protocolData = this.protocols.get(protocol)
    
    if (!client || !protocolData) {
      throw new Error('Protocol not supported on this chain')
    }

    // Mock implementation - in real app, would execute actual swap
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    
    return mockTxHash as Hash
  }

  // Get protocol statistics
  async getProtocolStats(protocolId: string, chainId: number): Promise<{
    tvl: string
    volume24h: string
    fees24h: string
    users24h: number
    apy?: string
  }> {
    const protocol = this.protocols.get(protocolId)
    
    if (!protocol) {
      throw new Error('Protocol not found')
    }

    // Mock implementation - in real app, would fetch from protocol APIs/subgraphs
    return {
      tvl: protocol.tvl,
      volume24h: '$125M',
      fees24h: '$375K',
      users24h: 15420,
      apy: protocol.apy
    }
  }

  // Search protocols
  searchProtocols(query: string): DeFiProtocol[] {
    const lowercaseQuery = query.toLowerCase()
    return this.getAllProtocols().filter(protocol =>
      protocol.name.toLowerCase().includes(lowercaseQuery) ||
      protocol.description.toLowerCase().includes(lowercaseQuery) ||
      protocol.features.some(feature => feature.toLowerCase().includes(lowercaseQuery))
    )
  }
}

// Export singleton instance
export const defiProtocolService = DeFiProtocolService.getInstance()
