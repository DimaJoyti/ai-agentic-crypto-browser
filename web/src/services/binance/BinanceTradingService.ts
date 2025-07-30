/**
 * Binance Trading Service (Backend Integration)
 * Handles secure trading operations through backend API
 * NEVER expose API keys in frontend code
 */

export interface BinanceAccount {
  makerCommission: number;
  takerCommission: number;
  buyerCommission: number;
  sellerCommission: number;
  canTrade: boolean;
  canWithdraw: boolean;
  canDeposit: boolean;
  updateTime: number;
  accountType: string;
  balances: BinanceBalance[];
  permissions: string[];
}

export interface BinanceBalance {
  asset: string;
  free: string;
  locked: string;
}

export interface BinanceOrder {
  symbol: string;
  orderId: number;
  orderListId: number;
  clientOrderId: string;
  price: string;
  origQty: string;
  executedQty: string;
  cummulativeQuoteQty: string;
  status: string;
  timeInForce: string;
  type: string;
  side: string;
  stopPrice: string;
  icebergQty: string;
  time: number;
  updateTime: number;
  isWorking: boolean;
  origQuoteOrderQty: string;
}

export interface PlaceOrderRequest {
  symbol: string;
  side: 'BUY' | 'SELL';
  type: 'MARKET' | 'LIMIT' | 'STOP_LOSS' | 'STOP_LOSS_LIMIT' | 'TAKE_PROFIT' | 'TAKE_PROFIT_LIMIT';
  quantity?: string;
  quoteOrderQty?: string;
  price?: string;
  stopPrice?: string;
  timeInForce?: 'GTC' | 'IOC' | 'FOK';
  newClientOrderId?: string;
  icebergQty?: string;
  newOrderRespType?: 'ACK' | 'RESULT' | 'FULL';
}

export interface CancelOrderRequest {
  symbol: string;
  orderId?: number;
  origClientOrderId?: string;
  newClientOrderId?: string;
}

export class BinanceTradingService {
  private readonly baseUrl: string;

  constructor(baseUrl: string = '/api/binance') {
    this.baseUrl = baseUrl;
  }

  /**
   * Get account information
   */
  async getAccount(): Promise<BinanceAccount> {
    const response = await fetch(`${this.baseUrl}/account`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies for authentication
    });

    if (!response.ok) {
      throw new Error(`Failed to get account info: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get all open orders
   */
  async getOpenOrders(symbol?: string): Promise<BinanceOrder[]> {
    const url = symbol 
      ? `${this.baseUrl}/openOrders?symbol=${symbol}`
      : `${this.baseUrl}/openOrders`;

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to get open orders: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get all orders for a symbol
   */
  async getAllOrders(symbol: string, limit?: number): Promise<BinanceOrder[]> {
    let url = `${this.baseUrl}/allOrders?symbol=${symbol}`;
    if (limit) url += `&limit=${limit}`;

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to get orders: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Place a new order
   */
  async placeOrder(orderRequest: PlaceOrderRequest): Promise<BinanceOrder> {
    const response = await fetch(`${this.baseUrl}/order`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(orderRequest),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(`Failed to place order: ${error.message || response.statusText}`);
    }

    return response.json();
  }

  /**
   * Cancel an order
   */
  async cancelOrder(cancelRequest: CancelOrderRequest): Promise<BinanceOrder> {
    const response = await fetch(`${this.baseUrl}/order`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(cancelRequest),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(`Failed to cancel order: ${error.message || response.statusText}`);
    }

    return response.json();
  }

  /**
   * Cancel all open orders for a symbol
   */
  async cancelAllOrders(symbol: string): Promise<BinanceOrder[]> {
    const response = await fetch(`${this.baseUrl}/openOrders`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({ symbol }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(`Failed to cancel all orders: ${error.message || response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get order status
   */
  async getOrder(symbol: string, orderId: number): Promise<BinanceOrder> {
    const response = await fetch(`${this.baseUrl}/order?symbol=${symbol}&orderId=${orderId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to get order: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get trading fees
   */
  async getTradingFees(symbol?: string): Promise<any> {
    const url = symbol 
      ? `${this.baseUrl}/tradeFee?symbol=${symbol}`
      : `${this.baseUrl}/tradeFee`;

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to get trading fees: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Test connectivity to the trading API
   */
  async testConnectivity(): Promise<boolean> {
    try {
      const response = await fetch(`${this.baseUrl}/ping`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      return response.ok;
    } catch (error) {
      console.error('Trading API connectivity test failed:', error);
      return false;
    }
  }

  /**
   * Get server time
   */
  async getServerTime(): Promise<{ serverTime: number }> {
    const response = await fetch(`${this.baseUrl}/time`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Failed to get server time: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Helper method to format order for display
   */
  formatOrder(order: BinanceOrder) {
    return {
      id: order.orderId,
      symbol: order.symbol,
      side: order.side,
      type: order.type,
      quantity: parseFloat(order.origQty),
      executedQuantity: parseFloat(order.executedQty),
      price: parseFloat(order.price),
      status: order.status,
      time: new Date(order.time).toLocaleString(),
      updateTime: new Date(order.updateTime).toLocaleString(),
    };
  }

  /**
   * Helper method to calculate order value
   */
  calculateOrderValue(order: BinanceOrder): number {
    return parseFloat(order.executedQty) * parseFloat(order.price);
  }
}

// Export singleton instance
export const binanceTradingService = new BinanceTradingService();
