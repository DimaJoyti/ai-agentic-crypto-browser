import { useState, useEffect, useCallback } from 'react';
import { 
  binanceTradingService, 
  BinanceAccount, 
  BinanceOrder, 
  PlaceOrderRequest,
  CancelOrderRequest 
} from '@/services/binance/BinanceTradingService';

/**
 * Hook for Binance account information
 */
export function useBinanceAccount() {
  const [account, setAccount] = useState<BinanceAccount | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  const fetchAccount = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      // Test connectivity first
      const connected = await binanceTradingService.testConnectivity();
      setIsConnected(connected);
      
      if (connected) {
        const accountData = await binanceTradingService.getAccount();
        setAccount(accountData);
      }
      
      setIsLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch account data');
      setIsLoading(false);
      setIsConnected(false);
    }
  }, []);

  useEffect(() => {
    fetchAccount();
  }, [fetchAccount]);

  return { 
    account, 
    isLoading, 
    error, 
    isConnected, 
    refetch: fetchAccount 
  };
}

/**
 * Hook for managing orders
 */
export function useBinanceOrders(symbol?: string) {
  const [orders, setOrders] = useState<BinanceOrder[]>([]);
  const [openOrders, setOpenOrders] = useState<BinanceOrder[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrders = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const [allOrders, openOrdersData] = await Promise.all([
        symbol ? binanceTradingService.getAllOrders(symbol, 100) : Promise.resolve([]),
        binanceTradingService.getOpenOrders(symbol)
      ]);
      
      setOrders(allOrders);
      setOpenOrders(openOrdersData);
      setIsLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch orders');
      setIsLoading(false);
    }
  }, [symbol]);

  const placeOrder = useCallback(async (orderRequest: PlaceOrderRequest): Promise<BinanceOrder> => {
    try {
      setError(null);
      const newOrder = await binanceTradingService.placeOrder(orderRequest);
      
      // Refresh orders after placing
      await fetchOrders();
      
      return newOrder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to place order';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [fetchOrders]);

  const cancelOrder = useCallback(async (cancelRequest: CancelOrderRequest): Promise<BinanceOrder> => {
    try {
      setError(null);
      const cancelledOrder = await binanceTradingService.cancelOrder(cancelRequest);
      
      // Refresh orders after cancelling
      await fetchOrders();
      
      return cancelledOrder;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to cancel order';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [fetchOrders]);

  const cancelAllOrders = useCallback(async (orderSymbol: string): Promise<BinanceOrder[]> => {
    try {
      setError(null);
      const cancelledOrders = await binanceTradingService.cancelAllOrders(orderSymbol);
      
      // Refresh orders after cancelling all
      await fetchOrders();
      
      return cancelledOrders;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to cancel all orders';
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  }, [fetchOrders]);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  return {
    orders,
    openOrders,
    isLoading,
    error,
    placeOrder,
    cancelOrder,
    cancelAllOrders,
    refetch: fetchOrders
  };
}

/**
 * Hook for trading operations with real-time updates
 */
export function useBinanceTrading(symbol: string) {
  const { account, isConnected: accountConnected } = useBinanceAccount();
  const { 
    orders, 
    openOrders, 
    isLoading: ordersLoading, 
    error: ordersError,
    placeOrder,
    cancelOrder,
    cancelAllOrders,
    refetch: refetchOrders
  } = useBinanceOrders(symbol);

  const [isPlacingOrder, setIsPlacingOrder] = useState(false);
  const [isCancellingOrder, setIsCancellingOrder] = useState(false);

  // Get balance for the base and quote assets
  const getBalance = useCallback((asset: string) => {
    if (!account) return { free: '0', locked: '0', total: '0' };
    
    const balance = account.balances.find(b => b.asset === asset);
    if (!balance) return { free: '0', locked: '0', total: '0' };
    
    const total = (parseFloat(balance.free) + parseFloat(balance.locked)).toString();
    return { ...balance, total };
  }, [account]);

  // Place order with loading state
  const placeOrderWithLoading = useCallback(async (orderRequest: PlaceOrderRequest) => {
    setIsPlacingOrder(true);
    try {
      const result = await placeOrder(orderRequest);
      return result;
    } finally {
      setIsPlacingOrder(false);
    }
  }, [placeOrder]);

  // Cancel order with loading state
  const cancelOrderWithLoading = useCallback(async (cancelRequest: CancelOrderRequest) => {
    setIsCancellingOrder(true);
    try {
      const result = await cancelOrder(cancelRequest);
      return result;
    } finally {
      setIsCancellingOrder(false);
    }
  }, [cancelOrder]);

  // Get trading fees
  const [tradingFees, setTradingFees] = useState<any>(null);
  
  useEffect(() => {
    const fetchTradingFees = async () => {
      try {
        const fees = await binanceTradingService.getTradingFees(symbol);
        setTradingFees(fees);
      } catch (err) {
        console.error('Failed to fetch trading fees:', err);
      }
    };

    if (accountConnected) {
      fetchTradingFees();
    }
  }, [symbol, accountConnected]);

  // Calculate available balance for trading
  const getAvailableBalance = useCallback((asset: string, price?: number) => {
    const balance = getBalance(asset);
    const available = parseFloat(balance.free);
    
    if (asset === 'USDT' && price) {
      // For USDT, calculate how much of the base asset can be bought
      return {
        asset,
        available: available.toString(),
        canBuy: (available / price).toString(),
        usdValue: available.toString()
      };
    }
    
    return {
      asset,
      available: available.toString(),
      canSell: available.toString(),
      usdValue: price ? (available * price).toString() : '0'
    };
  }, [getBalance]);

  return {
    // Account data
    account,
    isConnected: accountConnected,
    
    // Orders data
    orders,
    openOrders,
    ordersLoading,
    ordersError,
    
    // Trading operations
    placeOrder: placeOrderWithLoading,
    cancelOrder: cancelOrderWithLoading,
    cancelAllOrders,
    isPlacingOrder,
    isCancellingOrder,
    
    // Utility functions
    getBalance,
    getAvailableBalance,
    tradingFees,
    refetchOrders
  };
}

/**
 * Hook for order book trading (quick buy/sell at market prices)
 */
export function useOrderBookTrading(symbol: string) {
  const { placeOrder, isPlacingOrder } = useBinanceTrading(symbol);

  const quickBuy = useCallback(async (quantity: string, price?: string) => {
    const orderRequest: PlaceOrderRequest = {
      symbol,
      side: 'BUY',
      type: price ? 'LIMIT' : 'MARKET',
      quantity,
      ...(price && { price, timeInForce: 'GTC' })
    };

    return placeOrder(orderRequest);
  }, [symbol, placeOrder]);

  const quickSell = useCallback(async (quantity: string, price?: string) => {
    const orderRequest: PlaceOrderRequest = {
      symbol,
      side: 'SELL',
      type: price ? 'LIMIT' : 'MARKET',
      quantity,
      ...(price && { price, timeInForce: 'GTC' })
    };

    return placeOrder(orderRequest);
  }, [symbol, placeOrder]);

  return {
    quickBuy,
    quickSell,
    isPlacingOrder
  };
}
