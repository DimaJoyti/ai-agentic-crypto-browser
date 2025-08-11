import { EventEmitter } from 'events'

export interface RealtimeDataSubscription {
  id: string
  type: 'price' | 'defi' | 'nft' | 'network'
  params: Record<string, any>
  callback: (data: any) => void
}

export interface PriceUpdate {
  symbol: string
  price: number
  change24h: number
  volume24h: number
  timestamp: Date
}

export interface DeFiUpdate {
  protocol: string
  tvl: number
  apy?: number
  volume24h: number
  timestamp: Date
}

export interface NFTUpdate {
  collection: string
  floorPrice: number
  volume24h: number
  sales24h: number
  timestamp: Date
}

export interface NetworkUpdate {
  tps: number
  blockHeight: number
  avgBlockTime: number
  timestamp: Date
}

export class SolanaRealtimeDataService extends EventEmitter {
  private subscriptions: Map<string, RealtimeDataSubscription> = new Map()
  private websockets: Map<string, WebSocket> = new Map()
  private reconnectAttempts: Map<string, number> = new Map()
  private maxReconnectAttempts = 5
  private reconnectDelay = 5000

  constructor() {
    super()
    this.setupGlobalErrorHandling()
  }

  // Subscribe to real-time price updates
  subscribeToPrices(symbols: string[], callback: (data: PriceUpdate) => void): string {
    const id = this.generateId()
    const subscription: RealtimeDataSubscription = {
      id,
      type: 'price',
      params: { symbols },
      callback
    }

    this.subscriptions.set(id, subscription)
    this.connectPriceWebSocket(symbols, callback)
    
    return id
  }

  // Subscribe to DeFi protocol updates
  subscribeToDeFi(protocols: string[], callback: (data: DeFiUpdate) => void): string {
    const id = this.generateId()
    const subscription: RealtimeDataSubscription = {
      id,
      type: 'defi',
      params: { protocols },
      callback
    }

    this.subscriptions.set(id, subscription)
    this.connectDeFiWebSocket(protocols, callback)
    
    return id
  }

  // Subscribe to NFT collection updates
  subscribeToNFTs(collections: string[], callback: (data: NFTUpdate) => void): string {
    const id = this.generateId()
    const subscription: RealtimeDataSubscription = {
      id,
      type: 'nft',
      params: { collections },
      callback
    }

    this.subscriptions.set(id, subscription)
    this.connectNFTWebSocket(collections, callback)
    
    return id
  }

  // Subscribe to network statistics
  subscribeToNetwork(callback: (data: NetworkUpdate) => void): string {
    const id = this.generateId()
    const subscription: RealtimeDataSubscription = {
      id,
      type: 'network',
      params: {},
      callback
    }

    this.subscriptions.set(id, subscription)
    this.connectNetworkWebSocket(callback)
    
    return id
  }

  // Unsubscribe from updates
  unsubscribe(subscriptionId: string): void {
    const subscription = this.subscriptions.get(subscriptionId)
    if (!subscription) return

    this.subscriptions.delete(subscriptionId)
    
    // Close WebSocket if no more subscriptions of this type
    const hasOtherSubscriptions = Array.from(this.subscriptions.values())
      .some(sub => sub.type === subscription.type)
    
    if (!hasOtherSubscriptions) {
      const wsKey = `${subscription.type}-ws`
      const ws = this.websockets.get(wsKey)
      if (ws) {
        ws.close()
        this.websockets.delete(wsKey)
      }
    }
  }

  // Unsubscribe from all updates
  unsubscribeAll(): void {
    this.subscriptions.clear()
    this.websockets.forEach(ws => ws.close())
    this.websockets.clear()
    this.reconnectAttempts.clear()
  }

  private connectPriceWebSocket(symbols: string[], callback: (data: PriceUpdate) => void): void {
    const wsKey = 'price-ws'
    
    try {
      // Use CoinGecko WebSocket for price updates
      const ws = new WebSocket('wss://api.coingecko.com/api/v3/coins/solana/tickers')
      
      ws.onopen = () => {
        console.log('Connected to price WebSocket')
        this.reconnectAttempts.set(wsKey, 0)
        
        // Subscribe to symbols
        ws.send(JSON.stringify({
          method: 'subscribe',
          params: symbols.map(symbol => `${symbol.toLowerCase()}@ticker`)
        }))
      }

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.stream && data.data) {
            const symbol = data.stream.split('@')[0].toUpperCase()
            const priceUpdate: PriceUpdate = {
              symbol,
              price: parseFloat(data.data.c),
              change24h: parseFloat(data.data.P),
              volume24h: parseFloat(data.data.v),
              timestamp: new Date()
            }
            callback(priceUpdate)
            this.emit('priceUpdate', priceUpdate)
          }
        } catch (error) {
          console.error('Failed to parse price WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('Price WebSocket error:', error)
        this.emit('error', { type: 'price', error })
      }

      ws.onclose = () => {
        console.log('Price WebSocket closed')
        this.handleReconnect(wsKey, () => this.connectPriceWebSocket(symbols, callback))
      }

      this.websockets.set(wsKey, ws)
    } catch (error) {
      console.error('Failed to connect to price WebSocket:', error)
    }
  }

  private connectDeFiWebSocket(protocols: string[], callback: (data: DeFiUpdate) => void): void {
    const wsKey = 'defi-ws'
    
    // Mock DeFi WebSocket connection (replace with actual DeFi data source)
    const interval = setInterval(() => {
      protocols.forEach(protocol => {
        const mockUpdate: DeFiUpdate = {
          protocol,
          tvl: Math.random() * 1000000000,
          apy: Math.random() * 50,
          volume24h: Math.random() * 100000000,
          timestamp: new Date()
        }
        callback(mockUpdate)
        this.emit('defiUpdate', mockUpdate)
      })
    }, 30000) // Update every 30 seconds

    // Store interval as a mock WebSocket
    this.websockets.set(wsKey, { close: () => clearInterval(interval) } as any)
  }

  private connectNFTWebSocket(collections: string[], callback: (data: NFTUpdate) => void): void {
    const wsKey = 'nft-ws'
    
    try {
      // Use Magic Eden WebSocket for NFT updates
      const ws = new WebSocket('wss://api.magiceden.dev/ws')
      
      ws.onopen = () => {
        console.log('Connected to NFT WebSocket')
        this.reconnectAttempts.set(wsKey, 0)
        
        // Subscribe to collections
        collections.forEach(collection => {
          ws.send(JSON.stringify({
            method: 'subscribe',
            params: [`collection.${collection}`]
          }))
        })
      }

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'collection_update') {
            const nftUpdate: NFTUpdate = {
              collection: data.collection,
              floorPrice: data.floorPrice / 1e9, // Convert lamports to SOL
              volume24h: data.volume24h / 1e9,
              sales24h: data.sales24h,
              timestamp: new Date()
            }
            callback(nftUpdate)
            this.emit('nftUpdate', nftUpdate)
          }
        } catch (error) {
          console.error('Failed to parse NFT WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('NFT WebSocket error:', error)
        this.emit('error', { type: 'nft', error })
      }

      ws.onclose = () => {
        console.log('NFT WebSocket closed')
        this.handleReconnect(wsKey, () => this.connectNFTWebSocket(collections, callback))
      }

      this.websockets.set(wsKey, ws)
    } catch (error) {
      console.error('Failed to connect to NFT WebSocket:', error)
    }
  }

  private connectNetworkWebSocket(callback: (data: NetworkUpdate) => void): void {
    const wsKey = 'network-ws'
    
    // Mock network updates (replace with actual Solana RPC WebSocket)
    const interval = setInterval(() => {
      const networkUpdate: NetworkUpdate = {
        tps: Math.floor(Math.random() * 3000) + 2000,
        blockHeight: Math.floor(Date.now() / 400), // Approximate block height
        avgBlockTime: 400 + Math.random() * 100,
        timestamp: new Date()
      }
      callback(networkUpdate)
      this.emit('networkUpdate', networkUpdate)
    }, 5000) // Update every 5 seconds

    this.websockets.set(wsKey, { close: () => clearInterval(interval) } as any)
  }

  private handleReconnect(wsKey: string, reconnectFn: () => void): void {
    const attempts = this.reconnectAttempts.get(wsKey) || 0
    
    if (attempts < this.maxReconnectAttempts) {
      this.reconnectAttempts.set(wsKey, attempts + 1)
      
      setTimeout(() => {
        console.log(`Attempting to reconnect ${wsKey} (attempt ${attempts + 1})`)
        reconnectFn()
      }, this.reconnectDelay * Math.pow(2, attempts)) // Exponential backoff
    } else {
      console.error(`Max reconnection attempts reached for ${wsKey}`)
      this.emit('maxReconnectAttemptsReached', { wsKey })
    }
  }

  private setupGlobalErrorHandling(): void {
    this.on('error', (error) => {
      console.error('Realtime data service error:', error)
    })
  }

  private generateId(): string {
    return `sub_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  // Get connection status
  getConnectionStatus(): Record<string, boolean> {
    const status: Record<string, boolean> = {}
    
    this.websockets.forEach((ws, key) => {
      status[key] = ws.readyState === WebSocket.OPEN
    })
    
    return status
  }

  // Get active subscriptions count
  getActiveSubscriptionsCount(): number {
    return this.subscriptions.size
  }
}

// Singleton instance
export const solanaRealtimeData = new SolanaRealtimeDataService()
