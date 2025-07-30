# AI-Agentic Crypto Browser - HFT Trading System

## Overview

The AI-Agentic Crypto Browser includes a comprehensive High-Frequency Trading (HFT) system designed for cryptocurrency markets. This system combines traditional HFT capabilities with AI-powered analysis and decision-making through MCP (Model Context Protocol) tools integration.

## Architecture

### Core Components

1. **HFT Engine** (`internal/hft/`)
   - Real-time order execution engine
   - Sub-millisecond latency optimization
   - Advanced risk management
   - Performance monitoring and analytics

2. **Binance Integration** (`internal/binance/`)
   - REST API and WebSocket connections
   - Real-time market data streaming
   - Order management and execution
   - Account and portfolio tracking

3. **TradingView Integration** (`internal/tradingview/`)
   - Browser automation for signal extraction
   - Technical indicator monitoring
   - Chart pattern recognition
   - Signal processing and filtering

4. **MCP Tools Integration** (`internal/mcp/`)
   - AI-powered market analysis
   - Sentiment analysis from multiple sources
   - News impact assessment
   - Predictive analytics

5. **Strategy Engine** (`pkg/strategies/`)
   - Market making strategies
   - Arbitrage detection and execution
   - Momentum and trend-following strategies
   - Custom strategy framework

6. **Real-time Dashboard** (`web/src/components/trading/`)
   - Live trading interface
   - Portfolio management
   - Risk monitoring
   - Performance analytics

7. **Comprehensive API** (`api/`)
   - RESTful endpoints
   - WebSocket streaming
   - Real-time notifications
   - Complete CRUD operations

## Features

### High-Frequency Trading
- **Sub-millisecond Execution**: Optimized for ultra-low latency trading
- **Multi-Exchange Support**: Currently supports Binance with extensible architecture
- **Advanced Order Types**: Market, limit, stop-loss, and custom order types
- **Risk Management**: Real-time position monitoring and automatic risk controls

### AI-Powered Analysis
- **TradingView Signals**: Automated extraction of technical analysis signals
- **Sentiment Analysis**: Real-time sentiment from news, social media, and forums
- **Predictive Analytics**: AI-powered price prediction and trend analysis
- **Market Intelligence**: Comprehensive market insights and recommendations

### Trading Strategies
- **Market Making**: Professional market making with inventory management
- **Arbitrage**: Cross-exchange arbitrage detection and execution
- **Momentum**: Trend-following and momentum strategies
- **Custom Strategies**: Extensible framework for custom trading algorithms

### Risk Management
- **Position Limits**: Configurable position size and exposure limits
- **Loss Limits**: Daily, weekly, and monthly loss limits
- **Drawdown Protection**: Maximum drawdown monitoring and protection
- **Emergency Controls**: Instant stop-all and position closure capabilities

### Real-time Monitoring
- **Live Dashboard**: Comprehensive trading control center
- **Performance Analytics**: Detailed strategy and portfolio performance
- **Risk Monitoring**: Real-time risk metrics and alerts
- **Market Data**: Live market data visualization and analysis

## Quick Start

### Prerequisites

1. **Go 1.21+**: Install from [golang.org](https://golang.org/downloads/)
2. **Node.js 18+**: For the web dashboard
3. **Binance Account**: For API access (testnet recommended for development)
4. **Chrome/Chromium**: For TradingView integration

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-org/ai-agentic-crypto-browser.git
   cd ai-agentic-crypto-browser
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

4. **Build the application**:
   ```bash
   make build-hft
   ```

### Configuration

#### Binance API Setup

1. Create a Binance account and enable API access
2. Generate API key and secret (use testnet for development)
3. Configure in `.env`:
   ```env
   BINANCE_API_KEY=your_api_key_here
   BINANCE_SECRET_KEY=your_secret_key_here
   BINANCE_TESTNET=true
   ```

#### HFT Engine Configuration

```env
HFT_MAX_ORDERS_PER_SECOND=1000
HFT_LATENCY_TARGET_MICROS=1000
HFT_MAX_POSITION_SIZE=1.0
HFT_RISK_LIMIT_PERCENT=0.02
HFT_ENABLE_MARKET_MAKING=true
HFT_ENABLE_ARBITRAGE=true
```

#### Strategy Configuration

```env
STRATEGY_MAX_STRATEGIES=50
STRATEGY_EXECUTION_INTERVAL=100ms
STRATEGY_ENABLE_RISK_MANAGEMENT=true
STRATEGY_MAX_POSITION_SIZE=1.0
STRATEGY_MAX_DAILY_LOSS=5000.0
```

### Running the System

1. **Start the HFT system**:
   ```bash
   make run-hft
   ```

2. **Access the dashboard**:
   - Open http://localhost:8080 in your browser
   - WebSocket endpoint: ws://localhost:8080/ws/trading

3. **Monitor logs**:
   ```bash
   tail -f logs/hft-system.log
   ```

## API Endpoints

### HFT Engine
- `POST /api/hft/start` - Start the HFT engine
- `POST /api/hft/stop` - Stop the HFT engine
- `GET /api/hft/status` - Get engine status
- `GET /api/hft/metrics` - Get performance metrics

### Trading
- `GET /api/trading/orders` - List orders
- `POST /api/trading/orders` - Create order
- `GET /api/trading/positions` - List positions
- `GET /api/trading/signals` - Get trading signals

### Portfolio
- `GET /api/portfolio/summary` - Portfolio summary
- `GET /api/portfolio/positions` - Current positions
- `GET /api/portfolio/metrics` - Performance metrics
- `GET /api/portfolio/risk` - Risk metrics

### Strategies
- `GET /api/strategies` - List strategies
- `POST /api/strategies` - Create strategy
- `PUT /api/strategies/{id}` - Update strategy
- `POST /api/strategies/{id}/start` - Start strategy

### Risk Management
- `GET /api/risk/limits` - Risk limits
- `GET /api/risk/violations` - Risk violations
- `POST /api/risk/emergency-stop` - Emergency stop

## Development

### Running Tests
```bash
make test
make test-coverage
```

### Code Quality
```bash
make lint
make format
```

### Development Server
```bash
make dev-server  # With hot reload
```

### Debugging
```bash
# Enable debug mode
export DEBUG=true
export LOG_LEVEL=debug
make run-hft
```

## Production Deployment

### Docker
```bash
make docker-build
make docker-run
```

### Kubernetes
```bash
kubectl apply -f deployments/k8s/
```

### Monitoring
- Prometheus metrics: http://localhost:9090
- Grafana dashboard: http://localhost:3000
- Jaeger tracing: http://localhost:16686

## Security Considerations

1. **API Keys**: Store securely, use environment variables
2. **Network Security**: Use VPN for production deployments
3. **Access Control**: Implement proper authentication and authorization
4. **Audit Logging**: Enable comprehensive audit trails
5. **Risk Limits**: Always configure appropriate risk limits

## Performance Optimization

1. **Latency**: Optimize network connectivity to exchanges
2. **Hardware**: Use high-performance servers with low-latency networking
3. **Monitoring**: Continuously monitor and optimize performance
4. **Caching**: Implement intelligent caching strategies

## Troubleshooting

### Common Issues

1. **Connection Errors**: Check API keys and network connectivity
2. **High Latency**: Optimize network configuration and server location
3. **Order Rejections**: Verify account permissions and balance
4. **Strategy Errors**: Check strategy configuration and parameters

### Logs and Monitoring

- Application logs: `logs/hft-system.log`
- Error logs: `logs/errors.log`
- Performance metrics: Available via API and dashboard
- Health checks: `GET /health`

## Support

For support and questions:
- Documentation: [docs/](../docs/)
- Issues: GitHub Issues
- Community: Discord/Telegram (if available)

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Disclaimer

This software is for educational and research purposes. Trading cryptocurrencies involves substantial risk of loss. Use at your own risk and never trade with money you cannot afford to lose.
