import axios from 'axios'
import { SolPrice, MarketData } from '@/hooks/useSolanaMarketData'

export class SolanaMarketService {
  private readonly baseURL: string
  private readonly coingeckoAPI: string
  private readonly solanaAPI: string

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    this.coingeckoAPI = 'https://api.coingecko.com/api/v3'
    this.solanaAPI = 'https://api.mainnet-beta.solana.com'
  }

  async getSolPrice(): Promise<SolPrice> {
    try {
      // Try backend API first
      const response = await axios.get(`${this.baseURL}/api/solana/price`)
      return response.data
    } catch (error) {
      console.warn('Backend API unavailable, falling back to CoinGecko')
      
      // Fallback to CoinGecko API
      const response = await axios.get(
        `${this.coingeckoAPI}/simple/price?ids=solana&vs_currencies=usd&include_24hr_change=true&include_24hr_vol=true&include_market_cap=true`
      )
      
      const data = response.data.solana
      return {
        price: data.usd,
        change24h: data.usd_24h_change || 0,
        change7d: 0, // Not available in simple API
        volume24h: data.usd_24h_vol || 0,
        marketCap: data.usd_market_cap || 0,
        lastUpdated: new Date()
      }
    }
  }

  async getMarketData(): Promise<MarketData> {
    try {
      // Try backend API first
      const response = await axios.get(`${this.baseURL}/api/solana/market-data`)
      return response.data
    } catch (error) {
      console.warn('Backend API unavailable, using mock data')
      
      // Mock data for development
      return {
        marketCap: 45000000000, // $45B
        volume24h: 2500000000,  // $2.5B
        circulatingSupply: 467000000, // 467M SOL
        totalSupply: 580000000,      // 580M SOL
        transactions24h: 45000000,   // 45M transactions
        activeAddresses24h: 1200000, // 1.2M addresses
        networkHashRate: 0,          // Not applicable for Solana
        blockTime: 400               // 400ms average block time
      }
    }
  }

  async getHistoricalPrices(days: number = 30): Promise<Array<{ timestamp: number; price: number }>> {
    try {
      const response = await axios.get(
        `${this.coingeckoAPI}/coins/solana/market_chart?vs_currency=usd&days=${days}`
      )
      
      return response.data.prices.map(([timestamp, price]: [number, number]) => ({
        timestamp,
        price
      }))
    } catch (error) {
      console.error('Failed to fetch historical prices:', error)
      throw new Error('Failed to fetch historical price data')
    }
  }

  async getNetworkStats(): Promise<{
    tps: number
    blockHeight: number
    epochInfo: {
      epoch: number
      slotIndex: number
      slotsInEpoch: number
    }
    supply: {
      total: number
      circulating: number
      nonCirculating: number
    }
  }> {
    try {
      // Try backend API first
      const response = await axios.get(`${this.baseURL}/api/solana/network-stats`)
      return response.data
    } catch (error) {
      console.warn('Backend API unavailable, using Solana RPC')
      
      // Fallback to direct Solana RPC calls
      const [blockHeight, epochInfo, supply] = await Promise.all([
        this.getSolanaRPCData('getBlockHeight'),
        this.getSolanaRPCData('getEpochInfo'),
        this.getSolanaRPCData('getSupply')
      ])

      return {
        tps: 2500, // Approximate current TPS
        blockHeight: blockHeight.result,
        epochInfo: epochInfo.result,
        supply: {
          total: supply.result.value.total / 1e9, // Convert lamports to SOL
          circulating: supply.result.value.circulating / 1e9,
          nonCirculating: supply.result.value.nonCirculating / 1e9
        }
      }
    }
  }

  private async getSolanaRPCData(method: string, params: any[] = []): Promise<any> {
    const response = await axios.post(this.solanaAPI, {
      jsonrpc: '2.0',
      id: 1,
      method,
      params
    })
    
    if (response.data.error) {
      throw new Error(`Solana RPC error: ${response.data.error.message}`)
    }
    
    return response.data
  }

  async getTokenPrices(tokenMints: string[]): Promise<Record<string, number>> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/token-prices`, {
        params: { mints: tokenMints.join(',') }
      })
      return response.data
    } catch (error) {
      console.error('Failed to fetch token prices:', error)
      // Return empty object if API fails
      return {}
    }
  }

  async getTopTokens(limit: number = 20): Promise<Array<{
    mint: string
    symbol: string
    name: string
    price: number
    change24h: number
    volume24h: number
    marketCap: number
    logo?: string
  }>> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/top-tokens`, {
        params: { limit }
      })
      return response.data
    } catch (error) {
      console.error('Failed to fetch top tokens:', error)
      return []
    }
  }

  async getGasTracker(): Promise<{
    slow: number
    standard: number
    fast: number
    instant: number
    averageBlockTime: number
    congestionLevel: 'low' | 'medium' | 'high'
  }> {
    try {
      const response = await axios.get(`${this.baseURL}/api/solana/gas-tracker`)
      return response.data
    } catch (error) {
      console.warn('Gas tracker unavailable, using defaults')
      return {
        slow: 0.000005,     // 5,000 lamports
        standard: 0.000005, // 5,000 lamports (Solana has fixed fees)
        fast: 0.000005,     // 5,000 lamports
        instant: 0.000005,  // 5,000 lamports
        averageBlockTime: 400, // 400ms
        congestionLevel: 'low'
      }
    }
  }

  // WebSocket connection for real-time price updates
  createPriceWebSocket(onPriceUpdate: (price: SolPrice) => void): WebSocket | null {
    try {
      const ws = new WebSocket('wss://api.coingecko.com/api/v3/coins/solana/tickers')
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'price_update') {
            onPriceUpdate({
              price: data.price,
              change24h: data.change24h,
              change7d: data.change7d || 0,
              volume24h: data.volume24h,
              marketCap: data.marketCap,
              lastUpdated: new Date()
            })
          }
        } catch (error) {
          console.error('Failed to parse WebSocket price data:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('Price WebSocket error:', error)
      }

      return ws
    } catch (error) {
      console.error('Failed to create price WebSocket:', error)
      return null
    }
  }
}
