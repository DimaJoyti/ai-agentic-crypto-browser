import React, { useState, useEffect, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { 
  TrendingUp, 
  TrendingDown, 
  Activity, 
  Search,
  RefreshCw,
  Filter
} from 'lucide-react';

interface MarketData {
  symbol: string;
  price: string;
  change24h: string;
  volume: string;
  high24h: string;
  low24h: string;
  timestamp: string;
  bidPrice?: string;
  askPrice?: string;
  exchange?: string;
}

interface MarketDataPanelProps {
  marketData: MarketData[];
  selectedSymbol?: string;
  onSymbolChange?: (symbol: string) => void;
  className?: string;
}

interface SymbolStats {
  symbol: string;
  lastPrice: string;
  priceChange: string;
  priceChangePercent: string;
  volume: string;
  high24h: string;
  low24h: string;
  spread: string;
  lastUpdate: string;
  exchange: string;
}

export const MarketDataPanel: React.FC<MarketDataPanelProps> = ({
  marketData,
  selectedSymbol,
  onSymbolChange,
  className
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [exchangeFilter, setExchangeFilter] = useState('all');
  const [sortBy, setSortBy] = useState('volume');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Process market data into symbol statistics
  const symbolStats = useMemo(() => {
    const statsMap = new Map<string, SymbolStats>();
    
    marketData.forEach(tick => {
      const existing = statsMap.get(tick.symbol);
      
      if (!existing || new Date(tick.timestamp) > new Date(existing.lastUpdate)) {
        // Calculate spread (mock values if not provided)
        const bid = parseFloat(tick.bidPrice || (parseFloat(tick.price) * 0.999).toString());
        const ask = parseFloat(tick.askPrice || (parseFloat(tick.price) * 1.001).toString());
        const spread = ask - bid;
        const spreadPercent = ((spread / ask) * 100).toFixed(4);
        
        // Mock price change calculation (in real app, would track historical data)
        const priceChange = (Math.random() - 0.5) * parseFloat(tick.price) * 0.1;
        const priceChangePercent = ((priceChange / parseFloat(tick.price)) * 100).toFixed(2);
        
        statsMap.set(tick.symbol, {
          symbol: tick.symbol,
          lastPrice: tick.price,
          priceChange: priceChange.toFixed(2),
          priceChangePercent,
          volume: tick.volume,
          high24h: (parseFloat(tick.price) * 1.05).toFixed(2), // Mock high
          low24h: (parseFloat(tick.price) * 0.95).toFixed(2),  // Mock low
          spread: `${spread.toFixed(2)} (${spreadPercent}%)`,
          lastUpdate: tick.timestamp,
          exchange: tick.exchange || 'Unknown'
        });
      }
    });
    
    return Array.from(statsMap.values());
  }, [marketData]);

  // Filter and sort symbol stats
  const filteredStats = useMemo(() => {
    let filtered = symbolStats;
    
    // Apply search filter
    if (searchTerm) {
      filtered = filtered.filter(stat => 
        stat.symbol.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    // Apply exchange filter
    if (exchangeFilter !== 'all') {
      filtered = filtered.filter(stat => stat.exchange === exchangeFilter);
    }
    
    // Apply sorting
    filtered.sort((a, b) => {
      let aValue: number, bValue: number;
      
      switch (sortBy) {
        case 'symbol':
          return sortOrder === 'asc' 
            ? a.symbol.localeCompare(b.symbol)
            : b.symbol.localeCompare(a.symbol);
        case 'price':
          aValue = parseFloat(a.lastPrice);
          bValue = parseFloat(b.lastPrice);
          break;
        case 'change':
          aValue = parseFloat(a.priceChangePercent);
          bValue = parseFloat(b.priceChangePercent);
          break;
        case 'volume':
          aValue = parseFloat(a.volume);
          bValue = parseFloat(b.volume);
          break;
        default:
          return 0;
      }
      
      return sortOrder === 'asc' ? aValue - bValue : bValue - aValue;
    });
    
    return filtered;
  }, [symbolStats, searchTerm, exchangeFilter, sortBy, sortOrder]);

  // Get unique exchanges for filter
  const exchanges = useMemo(() => {
    const exchangeSet = new Set(marketData.map(tick => tick.exchange || 'Unknown'));
    return Array.from(exchangeSet).filter(Boolean);
  }, [marketData]);

  const formatPrice = (price: string) => {
    const value = parseFloat(price);
    return value.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 8
    });
  };

  const formatVolume = (volume: string) => {
    const value = parseFloat(volume);
    if (value >= 1e9) return `${(value / 1e9).toFixed(2)}B`;
    if (value >= 1e6) return `${(value / 1e6).toFixed(2)}M`;
    if (value >= 1e3) return `${(value / 1e3).toFixed(2)}K`;
    return value.toFixed(2);
  };

  const formatChange = (change: string, changePercent: string) => {
    const isPositive = parseFloat(change) >= 0;
    return (
      <div className={`flex items-center ${isPositive ? 'text-green-600' : 'text-red-600'}`}>
        {isPositive ? <TrendingUp className="w-4 h-4 mr-1" /> : <TrendingDown className="w-4 h-4 mr-1" />}
        <span>{isPositive ? '+' : ''}{change}</span>
        <span className="ml-1">({isPositive ? '+' : ''}{changePercent}%)</span>
      </div>
    );
  };

  const handleSort = (column: string) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('desc');
    }
  };

  return (
    <div className={`space-y-4 ${className}`}>
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center">
              <Activity className="w-5 h-5 mr-2" />
              Market Data
            </CardTitle>
            <Badge variant="outline">
              {filteredStats.length} symbols
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          {/* Filters */}
          <div className="flex flex-col sm:flex-row gap-4 mb-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  placeholder="Search symbols..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            
            <Select value={exchangeFilter} onValueChange={setExchangeFilter}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Exchange" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Exchanges</SelectItem>
                {exchanges.map(exchange => (
                  <SelectItem key={exchange} value={exchange || 'unknown'}>
                    {(exchange || 'Unknown').toUpperCase()}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            
            <Select value={sortBy} onValueChange={setSortBy}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="symbol">Symbol</SelectItem>
                <SelectItem value="price">Price</SelectItem>
                <SelectItem value="change">Change</SelectItem>
                <SelectItem value="volume">Volume</SelectItem>
              </SelectContent>
            </Select>
            
            <Button variant="outline" size="sm">
              <RefreshCw className="w-4 h-4" />
            </Button>
          </div>

          {/* Market Data Table */}
          <div className="border rounded-lg overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead 
                    className="cursor-pointer hover:bg-gray-50"
                    onClick={() => handleSort('symbol')}
                  >
                    Symbol {sortBy === 'symbol' && (sortOrder === 'asc' ? '↑' : '↓')}
                  </TableHead>
                  <TableHead 
                    className="cursor-pointer hover:bg-gray-50 text-right"
                    onClick={() => handleSort('price')}
                  >
                    Price {sortBy === 'price' && (sortOrder === 'asc' ? '↑' : '↓')}
                  </TableHead>
                  <TableHead 
                    className="cursor-pointer hover:bg-gray-50 text-right"
                    onClick={() => handleSort('change')}
                  >
                    24h Change {sortBy === 'change' && (sortOrder === 'asc' ? '↑' : '↓')}
                  </TableHead>
                  <TableHead 
                    className="cursor-pointer hover:bg-gray-50 text-right"
                    onClick={() => handleSort('volume')}
                  >
                    Volume {sortBy === 'volume' && (sortOrder === 'asc' ? '↑' : '↓')}
                  </TableHead>
                  <TableHead className="text-right">Spread</TableHead>
                  <TableHead className="text-center">Exchange</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredStats.slice(0, 50).map((stat) => (
                  <TableRow 
                    key={`${stat.symbol}-${stat.exchange}`}
                    className={`cursor-pointer hover:bg-gray-50 ${
                      selectedSymbol === stat.symbol ? 'bg-blue-50' : ''
                    }`}
                    onClick={() => onSymbolChange?.(stat.symbol)}
                  >
                    <TableCell className="font-medium">
                      <div className="flex items-center">
                        <span>{stat.symbol}</span>
                        {selectedSymbol === stat.symbol && (
                          <Badge variant="default" className="ml-2 text-xs">
                            Selected
                          </Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-right font-mono">
                      ${formatPrice(stat.lastPrice)}
                    </TableCell>
                    <TableCell className="text-right">
                      {formatChange(stat.priceChange, stat.priceChangePercent)}
                    </TableCell>
                    <TableCell className="text-right font-mono">
                      {formatVolume(stat.volume)}
                    </TableCell>
                    <TableCell className="text-right font-mono text-sm">
                      {stat.spread}
                    </TableCell>
                    <TableCell className="text-center">
                      <Badge variant="outline" className="text-xs">
                        {stat.exchange.toUpperCase()}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {filteredStats.length === 0 && (
            <div className="text-center py-8 text-muted-foreground">
              <Activity className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p>No market data available</p>
              <p className="text-sm">Check your connection and try again</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Selected Symbol Details */}
      {selectedSymbol && (
        <Card>
          <CardHeader>
            <CardTitle>{selectedSymbol} Details</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {filteredStats
                .filter(stat => stat.symbol === selectedSymbol)
                .map(stat => (
                  <React.Fragment key={stat.symbol}>
                    <div>
                      <p className="text-sm text-muted-foreground">Last Price</p>
                      <p className="text-lg font-bold">${formatPrice(stat.lastPrice)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">24h High</p>
                      <p className="text-lg font-bold">${formatPrice(stat.high24h)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">24h Low</p>
                      <p className="text-lg font-bold">${formatPrice(stat.low24h)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-muted-foreground">Volume</p>
                      <p className="text-lg font-bold">{formatVolume(stat.volume)}</p>
                    </div>
                  </React.Fragment>
                ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};
