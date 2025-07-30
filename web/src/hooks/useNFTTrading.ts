import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  nftTradingSystem,
  type NFTOrder,
  type TradingTransaction,
  type BulkOperation,
  type TradingStrategy,
  OrderType,
  OrderStatus,
  type TransactionType,
  TransactionStatus,
  type BulkOperationType,
  type BulkOperationItem,
  type TradingSystemEvent,
  type OrderFilters,
  type TransactionFilters
} from '@/lib/nft-trading'
import { toast } from 'sonner'

export interface NFTTradingState {
  orders: NFTOrder[]
  transactions: TradingTransaction[]
  bulkOperations: BulkOperation[]
  strategies: TradingStrategy[]
  isLoading: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseNFTTradingOptions {
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseNFTTradingReturn {
  // State
  state: NFTTradingState
  
  // Order Operations
  createBuyOrder: (
    contractAddress: Address,
    tokenId: string,
    price: string,
    currency: Address,
    marketplace: string,
    chainId?: number
  ) => Promise<NFTOrder>
  
  createSellOrder: (
    contractAddress: Address,
    tokenId: string,
    price: string,
    currency: Address,
    marketplace: string,
    chainId?: number
  ) => Promise<NFTOrder>
  
  cancelOrder: (orderId: string) => Promise<void>
  
  // Transaction Operations
  executeBuy: (orderId: string, buyerAddress: Address) => Promise<TradingTransaction>
  
  // Auction Operations
  createAuction: (
    contractAddress: Address,
    tokenId: string,
    startingPrice: string,
    reservePrice: string,
    duration: number,
    marketplace: string,
    chainId?: number
  ) => Promise<NFTOrder>
  
  placeBid: (auctionId: string, bidAmount: string, bidderAddress: Address) => Promise<NFTOrder>
  
  // Bulk Operations
  createBulkOperation: (type: BulkOperationType, items: BulkOperationItem[]) => Promise<BulkOperation>
  
  // Data Access
  getOrders: (filters?: OrderFilters) => NFTOrder[]
  getTransactions: (filters?: TransactionFilters) => TradingTransaction[]
  getBulkOperations: () => BulkOperation[]
  getStrategies: () => TradingStrategy[]
  
  // Analytics
  getTradingStats: () => TradingStats
  getOrderStats: () => OrderStats
  
  // Utilities
  clearError: () => void
  refresh: () => void
}

export interface TradingStats {
  totalOrders: number
  activeOrders: number
  filledOrders: number
  cancelledOrders: number
  totalVolume: number
  totalTransactions: number
  successRate: number
  averageOrderValue: number
}

export interface OrderStats {
  byType: Record<OrderType, number>
  byStatus: Record<OrderStatus, number>
  byMarketplace: Record<string, number>
  recentActivity: NFTOrder[]
}

export const useNFTTrading = (
  options: UseNFTTradingOptions = {}
): UseNFTTradingReturn => {
  const {
    enableNotifications = true,
    autoRefresh = false,
    refreshInterval = 30000 // 30 seconds
  } = options

  const [state, setState] = useState<NFTTradingState>({
    orders: [],
    transactions: [],
    bulkOperations: [],
    strategies: [],
    isLoading: false,
    error: null,
    lastUpdate: null
  })

  // Handle trading system events
  const handleTradingEvent = useCallback((event: TradingSystemEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'order_created':
          toast.success('Order Created', {
            description: `${event.order?.type} order created successfully`
          })
          break
        case 'order_filled':
          toast.success('Order Filled', {
            description: `Order filled at ${event.order?.price.amount} ${event.order?.price.currency}`
          })
          break
        case 'order_cancelled':
          toast.info('Order Cancelled', {
            description: 'Order has been cancelled'
          })
          break
        case 'order_failed':
          toast.error('Order Failed', {
            description: event.error?.message || 'Failed to create order'
          })
          break
        case 'transaction_created':
          toast.info('Transaction Created', {
            description: `${event.transaction?.type} transaction initiated`
          })
          break
        case 'transaction_confirmed':
          toast.success('Transaction Confirmed', {
            description: `Transaction confirmed in block ${event.transaction?.blockNumber}`
          })
          break
        case 'transaction_failed':
          toast.error('Transaction Failed', {
            description: event.error?.message || 'Transaction failed'
          })
          break
        case 'auction_created':
          toast.success('Auction Created', {
            description: `Auction started with ${event.order?.price.amount} ${event.order?.price.currency} starting price`
          })
          break
        case 'bid_placed':
          toast.info('Bid Placed', {
            description: `Bid of ${event.order?.price.amount} ${event.order?.price.currency} placed`
          })
          break
        case 'bulk_operation_created':
          toast.info('Bulk Operation Started', {
            description: `Processing ${event.bulkOperation?.operations.length} items`
          })
          break
        case 'bulk_operation_completed':
          const operation = event.bulkOperation!
          toast.success('Bulk Operation Completed', {
            description: `${operation.progress.completed}/${operation.progress.total} items processed successfully`
          })
          break
      }
    }

    // Update state
    setState(prev => ({
      ...prev,
      orders: nftTradingSystem.getOrders(),
      transactions: nftTradingSystem.getTransactions(),
      bulkOperations: nftTradingSystem.getBulkOperations(),
      strategies: nftTradingSystem.getStrategies(),
      error: event.type.includes('failed') ? event.error?.message || 'Operation failed' : null,
      lastUpdate: Date.now()
    }))
  }, [enableNotifications])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = nftTradingSystem.addEventListener(handleTradingEvent)

    // Load initial data
    const loadInitialData = () => {
      setState(prev => ({
        ...prev,
        orders: nftTradingSystem.getOrders(),
        transactions: nftTradingSystem.getTransactions(),
        bulkOperations: nftTradingSystem.getBulkOperations(),
        strategies: nftTradingSystem.getStrategies(),
        lastUpdate: Date.now()
      }))
    }

    loadInitialData()

    return () => {
      unsubscribe()
    }
  }, [handleTradingEvent])

  // Auto-refresh
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Create buy order
  const createBuyOrder = useCallback(async (
    contractAddress: Address,
    tokenId: string,
    price: string,
    currency: Address,
    marketplace: string,
    chainId: number = 1
  ): Promise<NFTOrder> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const order = await nftTradingSystem.createBuyOrder(
        contractAddress,
        tokenId,
        price,
        currency,
        marketplace,
        chainId
      )
      
      setState(prev => ({ ...prev, isLoading: false }))
      return order
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Create sell order
  const createSellOrder = useCallback(async (
    contractAddress: Address,
    tokenId: string,
    price: string,
    currency: Address,
    marketplace: string,
    chainId: number = 1
  ): Promise<NFTOrder> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const order = await nftTradingSystem.createSellOrder(
        contractAddress,
        tokenId,
        price,
        currency,
        marketplace,
        chainId
      )
      
      setState(prev => ({ ...prev, isLoading: false }))
      return order
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Cancel order
  const cancelOrder = useCallback(async (orderId: string): Promise<void> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      await nftTradingSystem.cancelOrder(orderId)
      setState(prev => ({ ...prev, isLoading: false }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Execute buy
  const executeBuy = useCallback(async (
    orderId: string,
    buyerAddress: Address
  ): Promise<TradingTransaction> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const transaction = await nftTradingSystem.executeBuy(orderId, buyerAddress)
      setState(prev => ({ ...prev, isLoading: false }))
      return transaction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Create auction
  const createAuction = useCallback(async (
    contractAddress: Address,
    tokenId: string,
    startingPrice: string,
    reservePrice: string,
    duration: number,
    marketplace: string,
    chainId: number = 1
  ): Promise<NFTOrder> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const auction = await nftTradingSystem.createAuction(
        contractAddress,
        tokenId,
        startingPrice,
        reservePrice,
        duration,
        marketplace,
        chainId
      )
      
      setState(prev => ({ ...prev, isLoading: false }))
      return auction
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Place bid
  const placeBid = useCallback(async (
    auctionId: string,
    bidAmount: string,
    bidderAddress: Address
  ): Promise<NFTOrder> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const bid = await nftTradingSystem.placeBid(auctionId, bidAmount, bidderAddress)
      setState(prev => ({ ...prev, isLoading: false }))
      return bid
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Create bulk operation
  const createBulkOperation = useCallback(async (
    type: BulkOperationType,
    items: BulkOperationItem[]
  ): Promise<BulkOperation> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const operation = await nftTradingSystem.createBulkOperation(type, items)
      setState(prev => ({ ...prev, isLoading: false }))
      return operation
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false
      }))
      throw error
    }
  }, [])

  // Get orders
  const getOrders = useCallback((filters?: OrderFilters): NFTOrder[] => {
    return nftTradingSystem.getOrders(filters)
  }, [])

  // Get transactions
  const getTransactions = useCallback((filters?: TransactionFilters): TradingTransaction[] => {
    return nftTradingSystem.getTransactions(filters)
  }, [])

  // Get bulk operations
  const getBulkOperations = useCallback((): BulkOperation[] => {
    return nftTradingSystem.getBulkOperations()
  }, [])

  // Get strategies
  const getStrategies = useCallback((): TradingStrategy[] => {
    return nftTradingSystem.getStrategies()
  }, [])

  // Get trading stats
  const getTradingStats = useCallback((): TradingStats => {
    const orders = state.orders
    const transactions = state.transactions

    const activeOrders = orders.filter(o => o.status === OrderStatus.ACTIVE).length
    const filledOrders = orders.filter(o => o.status === OrderStatus.FILLED).length
    const cancelledOrders = orders.filter(o => o.status === OrderStatus.CANCELLED).length

    const totalVolume = orders
      .filter(o => o.status === OrderStatus.FILLED)
      .reduce((sum, o) => sum + parseFloat(o.price.amount), 0)

    const successfulTransactions = transactions.filter(t => 
      t.status === TransactionStatus.CONFIRMED || t.status === TransactionStatus.FINALIZED
    ).length

    return {
      totalOrders: orders.length,
      activeOrders,
      filledOrders,
      cancelledOrders,
      totalVolume,
      totalTransactions: transactions.length,
      successRate: transactions.length > 0 ? (successfulTransactions / transactions.length) * 100 : 0,
      averageOrderValue: filledOrders > 0 ? totalVolume / filledOrders : 0
    }
  }, [state.orders, state.transactions])

  // Get order stats
  const getOrderStats = useCallback((): OrderStats => {
    const orders = state.orders

    const byType = orders.reduce((acc, order) => {
      acc[order.type] = (acc[order.type] || 0) + 1
      return acc
    }, {} as Record<OrderType, number>)

    const byStatus = orders.reduce((acc, order) => {
      acc[order.status] = (acc[order.status] || 0) + 1
      return acc
    }, {} as Record<OrderStatus, number>)

    const byMarketplace = orders.reduce((acc, order) => {
      acc[order.marketplace] = (acc[order.marketplace] || 0) + 1
      return acc
    }, {} as Record<string, number>)

    const recentActivity = orders
      .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
      .slice(0, 10)

    return {
      byType,
      byStatus,
      byMarketplace,
      recentActivity
    }
  }, [state.orders])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Refresh data
  const refresh = useCallback(() => {
    setState(prev => ({
      ...prev,
      orders: nftTradingSystem.getOrders(),
      transactions: nftTradingSystem.getTransactions(),
      bulkOperations: nftTradingSystem.getBulkOperations(),
      strategies: nftTradingSystem.getStrategies(),
      lastUpdate: Date.now()
    }))
  }, [])

  return {
    state,
    createBuyOrder,
    createSellOrder,
    cancelOrder,
    executeBuy,
    createAuction,
    placeBid,
    createBulkOperation,
    getOrders,
    getTransactions,
    getBulkOperations,
    getStrategies,
    getTradingStats,
    getOrderStats,
    clearError,
    refresh
  }
}

// Simplified hook for order management
export const useNFTOrders = () => {
  const { state, createBuyOrder, createSellOrder, cancelOrder, getOrders } = useNFTTrading()

  return {
    orders: state.orders,
    createBuyOrder,
    createSellOrder,
    cancelOrder,
    getOrders,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for auction management
export const useNFTAuctions = () => {
  const { state, createAuction, placeBid, getOrders } = useNFTTrading()

  const getAuctions = useCallback(() => {
    return getOrders({ type: OrderType.AUCTION })
  }, [getOrders])

  const getBids = useCallback((auctionId?: string) => {
    const bids = getOrders({ type: OrderType.BID })
    return auctionId 
      ? bids.filter(bid => bid.metadata.notes?.includes(auctionId))
      : bids
  }, [getOrders])

  return {
    auctions: getAuctions(),
    createAuction,
    placeBid,
    getBids,
    isLoading: state.isLoading,
    error: state.error
  }
}

// Hook for bulk operations
export const useBulkNFTOperations = () => {
  const { state, createBulkOperation, getBulkOperations } = useNFTTrading()

  return {
    bulkOperations: state.bulkOperations,
    createBulkOperation,
    getBulkOperations,
    isLoading: state.isLoading,
    error: state.error
  }
}
