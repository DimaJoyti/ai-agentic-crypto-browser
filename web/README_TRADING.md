# Real-time Trading Dashboard

## Overview

The Real-time Trading Dashboard is a comprehensive interface for monitoring and controlling the AI-Agentic Crypto Browser's High-Frequency Trading (HFT) system. It provides real-time data visualization, trading controls, and risk management capabilities.

## Features

### 🚀 **HFT Engine Control**
- Start/Stop HFT engine
- Real-time performance metrics
- Latency monitoring
- Order execution statistics

### 📊 **Market Data**
- Live price feeds for multiple symbols
- Real-time orderbook data
- Volume and volatility indicators
- Market depth visualization

### 💼 **Portfolio Management**
- Real-time portfolio value
- Position tracking
- P&L monitoring
- Performance analytics

### ⚡ **Trading Operations**
- Order management interface
- Signal monitoring
- Strategy execution
- Trade history

### 🛡️ **Risk Management**
- Real-time risk metrics
- Position limits monitoring
- Emergency stop controls
- Risk violation alerts

### 📈 **Performance Analytics**
- Strategy performance tracking
- Sharpe ratio calculations
- Drawdown monitoring
- Win/loss statistics

## Getting Started

### Prerequisites

1. **Backend Running**: Ensure the Go HFT backend is running on `localhost:8080`
2. **Node.js 18+**: Required for the React frontend
3. **Environment Setup**: Configure environment variables

### Installation

1. **Navigate to web directory**:
   ```bash
   cd web
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Configure environment**:
   ```bash
   cp .env.local.example .env.local
   # Edit .env.local with your configuration
   ```

4. **Start development server**:
   ```bash
   npm run dev
   ```

5. **Access the dashboard**:
   - Open http://localhost:3000
   - Navigate to `/trading` for the HFT dashboard

## Configuration

### Environment Variables

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# Trading Configuration
NEXT_PUBLIC_HFT_ENABLED=true
NEXT_PUBLIC_TRADING_SYMBOLS=BTCUSDT,ETHUSDT,BNBUSDT

# Features
NEXT_PUBLIC_ENABLE_PAPER_TRADING=true
NEXT_PUBLIC_ENABLE_LIVE_TRADING=false
```

## Dashboard Components

### 1. **HFT Metrics Panel**
- Orders per second
- Average latency
- Success rate
- Engine status

### 2. **Market Data Panel**
- Real-time price feeds
- Volume indicators
- Price change percentages
- Market status

### 3. **Portfolio Panel**
- Total portfolio value
- Current positions
- Unrealized P&L
- Risk metrics

### 4. **Order Management Panel**
- Active orders
- Order history
- Quick order placement
- Order cancellation

### 5. **Strategy Panel**
- Active strategies
- Strategy performance
- Start/stop controls
- Configuration options

### 6. **Risk Management Panel**
- Risk limits
- Violation alerts
- Emergency controls
- Risk metrics

## Real-time Features

### WebSocket Integration
- Real-time data updates
- Live order status
- Market data streaming
- System notifications

### Auto-refresh
- Configurable refresh intervals
- Manual refresh controls
- Connection status monitoring
- Error handling

## API Integration

The dashboard integrates with the Go backend through:

### REST API Endpoints
- `/api/hft/*` - HFT engine control
- `/api/trading/*` - Trading operations
- `/api/portfolio/*` - Portfolio data
- `/api/strategies/*` - Strategy management
- `/api/risk/*` - Risk management

### WebSocket Streams
- `ws://localhost:8080/ws/trading` - Real-time trading data
- Market data updates
- Order status changes
- System notifications

## Development

### Project Structure
```
web/src/
├── app/
│   ├── trading/page.tsx          # Trading dashboard route
│   └── dashboard/page.tsx        # Main dashboard
├── components/trading/
│   ├── TradingDashboard.tsx      # Main dashboard component
│   ├── MarketDataPanel.tsx       # Market data display
│   ├── PortfolioPanel.tsx        # Portfolio management
│   ├── OrderManagementPanel.tsx  # Order operations
│   ├── StrategyPanel.tsx         # Strategy control
│   ├── RiskManagementPanel.tsx   # Risk monitoring
│   └── PerformancePanel.tsx      # Performance analytics
├── hooks/
│   ├── useTradingDashboard.ts    # Main trading hook
│   └── useWebSocket.ts           # WebSocket management
└── lib/
    ├── trading-api.ts            # API client
    └── utils.ts                  # Utility functions
```

### Key Hooks

#### `useTradingDashboard`
Main hook that provides:
- HFT engine control
- Real-time data management
- API integration
- Error handling

#### `useWebSocket`
WebSocket management:
- Connection handling
- Message processing
- Reconnection logic
- Error recovery

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Ensure backend is running on port 8080
   - Check WebSocket URL configuration
   - Verify CORS settings

2. **No Data Loading**
   - Check API endpoint URLs
   - Verify backend API responses
   - Check browser console for errors

3. **WebSocket Disconnections**
   - Check network connectivity
   - Verify WebSocket endpoint
   - Monitor connection status

### Debug Mode

Enable debug mode in `.env.local`:
```env
NEXT_PUBLIC_DEBUG=true
```

This will show additional logging and debug information.

## Production Deployment

### Build for Production
```bash
npm run build
npm start
```

### Environment Configuration
- Update API URLs for production
- Configure proper CORS settings
- Set up SSL/TLS for WebSocket connections
- Enable production optimizations

## Security Considerations

1. **API Security**
   - Use HTTPS in production
   - Implement proper authentication
   - Validate all user inputs

2. **WebSocket Security**
   - Use WSS (WebSocket Secure) in production
   - Implement connection authentication
   - Rate limit connections

3. **Data Protection**
   - Encrypt sensitive data
   - Implement proper access controls
   - Log security events

## Support

For issues and questions:
- Check the main project documentation
- Review the Go backend logs
- Monitor browser console for errors
- Check WebSocket connection status

## License

This project is part of the AI-Agentic Crypto Browser system and follows the same license terms.
