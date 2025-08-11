import axios from 'axios'
import { DeFiProtocol, YieldOpportunity, DeFiStats } from '@/hooks/useSolanaDeFi'

export class SolanaDeFiService {
  private readonly baseURL: string
  private readonly defillama: string
  private readonly jupiterAPI: string

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    this.defillama = 'https://api.llama.fi'
    this.jupiterAPI = 'https://quote-api.jup.ag/v6'
  }

  async getProtocols(options: {
    category?: DeFiProtocol['category']
    minTVL?: number
  } = {}): Promise<DeFiProtocol[]> {
    try {
      // Try backend API first
      const response = await axios.get(`${this.baseURL}/api/solana/defi/protocols`, {
        params: options
      })
      return response.data
    } catch (error) {
      console.warn('Backend API unavailable, falling back to DeFiLlama')
      
      // Fallback to DeFiLlama API
      const response = await axios.get(`${this.defillama}/protocols`)
      const allProtocols = response.data
      
      // Filter for Solana protocols
      const solanaProtocols = allProtocols
        .filter((protocol: any) => 
          protocol.chains?.includes('Solana') && 
          protocol.tvl >= (options.minTVL || 0)
        )
        .map((protocol: any) => ({
          id: protocol.slug,
          name: protocol.name,
          tvl: protocol.tvl,
          tvlChange24h: protocol.change_1d || 0,
          category: this.mapCategory(protocol.category),
          volume24h: protocol.volume24h || 0,
          users24h: 0, // Not available in DeFiLlama
          logo: protocol.logo,
          website: protocol.url,
          description: protocol.description || ''
        }))
        .slice(0, 50) // Limit to top 50

      return solanaProtocols
    }
  }

  async getProtocol(protocolId: string): Promise<DeFiProtocol> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/defi/protocols/${protocolId}`)
      return response.data
    } catch (error) {
      throw new Error(`Failed to fetch protocol ${protocolId}`)
    }
  }

  async getYieldOpportunities(options: {
    minAPY?: number
    maxRisk?: YieldOpportunity['risk']
  } = {}): Promise<YieldOpportunity[]> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/defi/yields`, {
        params: options
      })
      return response.data
    } catch (error) {
      console.warn('Yield API unavailable, using mock data')
      
      // Mock yield opportunities for development
      const mockYields: YieldOpportunity[] = [
        {
          protocol: 'Marinade',
          pool: 'mSOL Staking',
          apy: 7.2,
          tvl: 1200000000,
          risk: 'low' as const,
          tokens: ['SOL', 'mSOL'],
          minimumDeposit: 0.01
        },
        {
          protocol: 'Raydium',
          pool: 'SOL-USDC LP',
          apy: 12.5,
          tvl: 450000000,
          risk: 'medium' as const,
          tokens: ['SOL', 'USDC'],
          minimumDeposit: 10
        },
        {
          protocol: 'Orca',
          pool: 'SOL-mSOL Whirlpool',
          apy: 9.8,
          tvl: 320000000,
          risk: 'low' as const,
          tokens: ['SOL', 'mSOL'],
          minimumDeposit: 0.1
        },
        {
          protocol: 'Jupiter',
          pool: 'JUP-SOL LP',
          apy: 25.3,
          tvl: 180000000,
          risk: 'high' as const,
          tokens: ['JUP', 'SOL'],
          minimumDeposit: 5
        },
        {
          protocol: 'Kamino',
          pool: 'USDC Lending',
          apy: 8.7,
          tvl: 890000000,
          risk: 'low' as const,
          tokens: ['USDC'],
          minimumDeposit: 100
        }
      ]

      return mockYields.filter(yield_ =>
        yield_.apy >= (options.minAPY || 0)
      )
    }
  }

  async getDeFiStats(): Promise<DeFiStats> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/defi/stats`)
      return response.data
    } catch (error) {
      console.warn('DeFi stats API unavailable, calculating from protocols')
      
      const protocols = await this.getProtocols()
      const yields = await this.getYieldOpportunities()
      
      const totalTVL = protocols.reduce((sum, protocol) => sum + protocol.tvl, 0)
      const totalVolume24h = protocols.reduce((sum, protocol) => sum + protocol.volume24h, 0)
      const avgAPY = yields.length > 0 
        ? yields.reduce((sum, yield_) => sum + yield_.apy, 0) / yields.length 
        : 0
      
      const topProtocolByTVL = protocols.sort((a, b) => b.tvl - a.tvl)[0]?.name || ''
      const topProtocolByVolume = protocols.sort((a, b) => b.volume24h - a.volume24h)[0]?.name || ''

      return {
        totalTVL,
        totalVolume24h,
        totalProtocols: protocols.length,
        avgAPY,
        topProtocolByTVL,
        topProtocolByVolume
      }
    }
  }

  async getSwapQuote(params: {
    inputMint: string
    outputMint: string
    amount: number
    slippageBps?: number
  }): Promise<{
    inputAmount: number
    outputAmount: number
    priceImpact: number
    route: any[]
    fees: number
  }> {
    try {
      const response = await axios.get(`${this.jupiterAPI}/quote`, {
        params: {
          inputMint: params.inputMint,
          outputMint: params.outputMint,
          amount: params.amount,
          slippageBps: params.slippageBps || 50
        }
      })
      
      const quote = response.data
      return {
        inputAmount: parseInt(quote.inAmount),
        outputAmount: parseInt(quote.outAmount),
        priceImpact: parseFloat(quote.priceImpactPct || '0'),
        route: quote.routePlan || [],
        fees: parseInt(quote.platformFee?.amount || '0')
      }
    } catch (error) {
      console.error('Failed to get swap quote:', error)
      throw new Error('Failed to get swap quote')
    }
  }

  private mapCategory(category: string): DeFiProtocol['category'] {
    const categoryMap: Record<string, DeFiProtocol['category']> = {
      'Dexes': 'dex',
      'Lending': 'lending',
      'Yield Farming': 'yield',
      'Derivatives': 'derivatives',
      'Insurance': 'insurance'
    }
    
    return categoryMap[category] || 'dex'
  }
}
