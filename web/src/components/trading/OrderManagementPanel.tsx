import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { 
  ShoppingCart, 
  X, 
  Clock, 
  CheckCircle, 
  AlertCircle,
  RefreshCw
} from 'lucide-react';

interface Order {
  id: string;
  symbol: string;
  side: 'BUY' | 'SELL';
  type: string;
  quantity: string;
  price: string;
  filledQty: string;
  status: string;
  createdAt: string;
  exchange: string;
  strategyId: string;
}

export const OrderManagementPanel: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock data for demonstration
    const mockOrders: Order[] = [
      {
        id: '1',
        symbol: 'BTCUSDT',
        side: 'BUY',
        type: 'LIMIT',
        quantity: '0.1',
        price: '45000.00',
        filledQty: '0.05',
        status: 'PARTIAL_FILL',
        createdAt: new Date().toISOString(),
        exchange: 'binance',
        strategyId: 'market_making_1'
      },
      {
        id: '2',
        symbol: 'ETHUSDT',
        side: 'SELL',
        type: 'MARKET',
        quantity: '1.0',
        price: '3000.00',
        filledQty: '1.0',
        status: 'FILLED',
        createdAt: new Date(Date.now() - 60000).toISOString(),
        exchange: 'binance',
        strategyId: 'arbitrage_1'
      }
    ];
    
    setOrders(mockOrders);
    setLoading(false);
  }, []);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'FILLED':
        return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'PARTIAL_FILL':
        return <Clock className="w-4 h-4 text-yellow-600" />;
      case 'CANCELED':
        return <X className="w-4 h-4 text-red-600" />;
      case 'REJECTED':
        return <AlertCircle className="w-4 h-4 text-red-600" />;
      default:
        return <Clock className="w-4 h-4 text-blue-600" />;
    }
  };

  const getStatusBadge = (status: string) => {
    const variant = status === 'FILLED' ? 'default' : 
                   status === 'PARTIAL_FILL' ? 'secondary' : 'destructive';
    return <Badge variant={variant}>{status}</Badge>;
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center">
            <ShoppingCart className="w-5 h-5 mr-2" />
            Order Management
          </CardTitle>
          <Button variant="outline" size="sm">
            <RefreshCw className="w-4 h-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Order ID</TableHead>
              <TableHead>Symbol</TableHead>
              <TableHead>Side</TableHead>
              <TableHead>Type</TableHead>
              <TableHead className="text-right">Quantity</TableHead>
              <TableHead className="text-right">Price</TableHead>
              <TableHead className="text-right">Filled</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Exchange</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {orders.map((order) => (
              <TableRow key={order.id}>
                <TableCell className="font-mono text-sm">{order.id}</TableCell>
                <TableCell className="font-medium">{order.symbol}</TableCell>
                <TableCell>
                  <Badge variant={order.side === 'BUY' ? 'default' : 'destructive'}>
                    {order.side}
                  </Badge>
                </TableCell>
                <TableCell>{order.type}</TableCell>
                <TableCell className="text-right font-mono">{order.quantity}</TableCell>
                <TableCell className="text-right font-mono">${order.price}</TableCell>
                <TableCell className="text-right font-mono">{order.filledQty}</TableCell>
                <TableCell>
                  <div className="flex items-center">
                    {getStatusIcon(order.status)}
                    <span className="ml-2">{getStatusBadge(order.status)}</span>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">{order.exchange.toUpperCase()}</Badge>
                </TableCell>
                <TableCell>
                  {order.status === 'NEW' || order.status === 'PARTIAL_FILL' ? (
                    <Button variant="outline" size="sm">
                      Cancel
                    </Button>
                  ) : (
                    <span className="text-muted-foreground text-sm">-</span>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
