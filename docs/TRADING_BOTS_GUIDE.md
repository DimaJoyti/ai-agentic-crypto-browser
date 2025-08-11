# 🤖 7 Trading Bots System - Complete Guide

## 📋 **Overview**

The AI-Agentic Crypto Browser now includes a sophisticated 7-bot trading system with distinct strategies, comprehensive risk management, and institutional-grade features.

### **🎯 Bot Specifications**

| Bot # | Strategy | Trading Pairs | Exchange | Risk Level | Capital % |
|-------|----------|---------------|----------|------------|-----------|
| 1 | **DCA Bot** | BTC/USDT, ETH/USDT | Binance | Low | 20% |
| 2 | **Grid Trading Bot** | BNB/USDT, ADA/USDT | Binance | Medium | 15% |
| 3 | **Momentum Bot** | SOL/USDT, AVAX/USDT | Coinbase | High | 10% |
| 4 | **Mean Reversion Bot** | DOT/USDT, LINK/USDT | Kraken | Medium | 15% |
| 5 | **Arbitrage Bot** | BTC/USDT (Multi-exchange) | Binance+Coinbase | Low | 20% |
| 6 | **Scalping Bot** | ETH/USDT | Binance | High | 10% |
| 7 | **Swing Trading Bot** | Multiple Altcoins | Binance | Medium | 10% |

## 🚀 **Quick Start**

### **1. Prerequisites**

```bash
# Required software
- Go 1.21+
- PostgreSQL 15+ with TimescaleDB
- Redis 7.0+
- Docker & Docker Compose

# Exchange API Keys (for each exchange you plan to use)
- Binance API Key & Secret
- Coinbase Pro API Key & Secret & Passphrase
- Kraken API Key & Secret
```

### **2. Installation**

```bash
# Clone the repository
git clone https://github.com/DimaJoyti/ai-agentic-crypto-browser.git
cd ai-agentic-crypto-browser

# Install dependencies
go mod tidy

# Run database migrations
make migrate-up

# Apply trading bots schema
psql -d ai_agentic_browser -f migrations/006_trading_bots_schema.sql
```

### **3. Configuration**

```bash
# Copy and edit trading bots configuration
cp configs/trading-bots.yaml.example configs/trading-bots.yaml

# Edit configuration with your settings
nano configs/trading-bots.yaml
```

### **4. Start the Trading Bots System**

```bash
# Start infrastructure services
docker-compose up -d postgres redis

# Start the trading bots service
go run cmd/trading-bots/main.go

# Or build and run
make build-trading-bots
./bin/trading-bots
```

## ⚙️ **Configuration Guide**

### **Trading Bots Configuration (`configs/trading-bots.yaml`)**

```yaml
# Server Configuration
server:
  host: "0.0.0.0"
  port: 8090
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s

# Global Bot Settings
trading_bots:
  max_concurrent_bots: 7
  execution_interval: 5s
  order_timeout: 30s
  retry_attempts: 3
  performance_update_interval: 1m
  health_check_interval: 30s

# Exchange Configuration
exchanges:
  binance:
    api_url: "https://api.binance.com"
    testnet_url: "https://testnet.binance.vision"
    rate_limit: 1200
    sandbox: false
    api_key: "your_binance_api_key"
    api_secret: "your_binance_api_secret"
    
  coinbase:
    api_url: "https://api.exchange.coinbase.com"
    sandbox_url: "https://api-public.sandbox.exchange.coinbase.com"
    rate_limit: 10
    sandbox: false
    api_key: "your_coinbase_api_key"
    api_secret: "your_coinbase_api_secret"
    passphrase: "your_coinbase_passphrase"
```

### **Individual Bot Configuration**

Each bot can be configured with specific parameters:

```yaml
# Example: DCA Bot Configuration
dca_bot:
  strategy_params:
    investment_amount: 100.0    # USDT per interval
    interval: "1h"              # Investment frequency
    max_deviation: 0.05         # 5% price deviation threshold
    
  risk_management:
    max_position_size: 0.20     # 20% of total capital
    stop_loss: 0.15             # 15% stop loss
    take_profit: 0.30           # 30% take profit
    max_drawdown: 0.10          # 10% max drawdown
```

## 📊 **Bot Strategies Explained**

### **1. DCA (Dollar Cost Averaging) Bot**
- **Strategy**: Invests fixed amounts at regular intervals
- **Best For**: Long-term accumulation, reducing volatility impact
- **Risk Level**: Low
- **Timeframe**: 1 hour intervals
- **Key Parameters**: Investment amount, interval, price deviation threshold

### **2. Grid Trading Bot**
- **Strategy**: Places buy/sell orders at predetermined price levels
- **Best For**: Sideways markets, capturing small price movements
- **Risk Level**: Medium
- **Timeframe**: Continuous
- **Key Parameters**: Grid levels, spacing, upper/lower bounds

### **3. Momentum Bot**
- **Strategy**: Follows price trends and momentum indicators
- **Best For**: Trending markets, breakout scenarios
- **Risk Level**: High
- **Timeframe**: 1-4 hours
- **Key Parameters**: RSI thresholds, momentum period, volume threshold

### **4. Mean Reversion Bot**
- **Strategy**: Buys oversold and sells overbought conditions
- **Best For**: Range-bound markets, temporary price dislocations
- **Risk Level**: Medium
- **Timeframe**: 4 hours
- **Key Parameters**: Bollinger Bands, RSI levels, standard deviation

### **5. Arbitrage Bot**
- **Strategy**: Exploits price differences across exchanges
- **Best For**: Market inefficiencies, risk-free profits
- **Risk Level**: Low
- **Timeframe**: Real-time
- **Key Parameters**: Minimum profit threshold, execution time, slippage

### **6. Scalping Bot**
- **Strategy**: High-frequency trading for small profits
- **Best For**: High liquidity pairs, volatile markets
- **Risk Level**: High
- **Timeframe**: 1 minute
- **Key Parameters**: Profit target, holding time, spread threshold

### **7. Swing Trading Bot**
- **Strategy**: Captures medium-term price swings
- **Best For**: Trending markets, multi-day positions
- **Risk Level**: Medium
- **Timeframe**: 4 hours to daily
- **Key Parameters**: MACD settings, moving averages, RSI

## 🔧 **API Documentation**

### **Bot Management Endpoints**

```bash
# List all trading bots
GET /api/v1/trading-bots

# Create a new trading bot
POST /api/v1/trading-bots
{
  "name": "My DCA Bot",
  "strategy": "dca",
  "trading_pairs": ["BTC/USDT"],
  "exchange": "binance",
  "strategy_params": {
    "investment_amount": 100,
    "interval": "1h"
  },
  "capital": {
    "initial_balance": 10000,
    "allocation_percentage": 0.20
  }
}

# Get specific bot details
GET /api/v1/trading-bots/{botId}

# Start a trading bot
POST /api/v1/trading-bots/{botId}/start

# Stop a trading bot
POST /api/v1/trading-bots/{botId}/stop

# Get bot performance metrics
GET /api/v1/trading-bots/{botId}/performance

# Get bot trade history
GET /api/v1/trading-bots/{botId}/trades
```

### **Strategy Management Endpoints**

```bash
# List available strategies
GET /api/v1/trading-strategies

# Get strategy details
GET /api/v1/trading-strategies/{strategyId}

# Get strategy performance
GET /api/v1/trading-strategies/{strategyId}/performance
```

### **Bulk Operations**

```bash
# Start all bots
POST /api/v1/trading-bots/start-all

# Stop all bots
POST /api/v1/trading-bots/stop-all

# Get all bots performance
GET /api/v1/trading-bots/performance
```

## 📈 **Monitoring & Analytics**

### **Health Checks**

```bash
# System health
curl http://localhost:8090/health

# Detailed health with metrics
curl http://localhost:8090/api/v1/health
```

### **Metrics Endpoints**

```bash
# Prometheus metrics
curl http://localhost:8090/metrics

# Bot-specific metrics
curl http://localhost:8090/api/v1/trading-bots/{botId}/metrics
```

### **Performance Monitoring**

The system provides comprehensive performance tracking:

- **Real-time P&L**: Live profit/loss tracking
- **Win Rate**: Percentage of profitable trades
- **Sharpe Ratio**: Risk-adjusted returns
- **Maximum Drawdown**: Largest peak-to-trough decline
- **Trade Statistics**: Count, average duration, best/worst trades

### **Alerting**

Configure alerts for:
- Bot errors or failures
- Performance thresholds
- Risk limit breaches
- Exchange connectivity issues

## 🔒 **Security & Risk Management**

### **API Key Security**

```bash
# Store API keys securely (encrypted in database)
# Never commit API keys to version control
# Use environment variables for sensitive data

export BINANCE_API_KEY="your_key"
export BINANCE_API_SECRET="your_secret"
```

### **Risk Controls**

- **Position Sizing**: Maximum position size per bot
- **Stop Loss**: Automatic loss cutting
- **Take Profit**: Profit taking levels
- **Drawdown Limits**: Maximum portfolio decline
- **Correlation Limits**: Prevent over-concentration

### **Portfolio Risk Management**

```yaml
portfolio_risk:
  max_total_exposure: 0.80      # 80% max portfolio exposure
  correlation_limit: 0.70       # 70% max correlation between bots
  var_limit: 0.05               # 5% Value at Risk limit
```

## 🧪 **Testing & Backtesting**

### **Paper Trading Mode**

```bash
# Enable paper trading for testing
export PAPER_TRADING=true

# Start bots in paper trading mode
go run cmd/trading-bots/main.go --paper-trading
```

### **Backtesting**

```bash
# Run backtest for specific strategy
curl -X POST http://localhost:8090/api/v1/backtests \
  -H "Content-Type: application/json" \
  -d '{
    "strategy": "dca",
    "start_date": "2023-01-01",
    "end_date": "2023-12-31",
    "initial_balance": 10000
  }'
```

## 🚀 **Deployment**

### **Production Deployment**

```bash
# Build production binary
make build-trading-bots-prod

# Deploy with Docker
docker build -t trading-bots:latest -f cmd/trading-bots/Dockerfile .
docker run -d --name trading-bots \
  -p 8090:8090 \
  -e CONFIG_FILE=/app/configs/production.yaml \
  trading-bots:latest

# Deploy with Kubernetes
kubectl apply -f deployments/k8s/trading-bots/
```

### **Environment Variables**

```bash
# Required environment variables
export CONFIG_FILE="/path/to/config.yaml"
export DATABASE_URL="postgres://user:pass@localhost/db"
export REDIS_URL="redis://localhost:6379"

# Exchange API keys
export BINANCE_API_KEY="your_key"
export BINANCE_API_SECRET="your_secret"
export COINBASE_API_KEY="your_key"
export COINBASE_API_SECRET="your_secret"
export COINBASE_PASSPHRASE="your_passphrase"
```

## 🔧 **Troubleshooting**

### **Common Issues**

1. **Bot Not Starting**
   ```bash
   # Check logs
   tail -f logs/trading-bots.log
   
   # Verify configuration
   go run cmd/trading-bots/main.go --validate-config
   ```

2. **Exchange Connection Issues**
   ```bash
   # Test exchange connectivity
   curl -X GET http://localhost:8090/api/v1/exchanges/test
   
   # Check API key permissions
   curl -X GET http://localhost:8090/api/v1/exchanges/binance/account
   ```

3. **Database Issues**
   ```bash
   # Check database connection
   psql -d ai_agentic_browser -c "SELECT 1;"
   
   # Run migrations
   make migrate-up
   ```

### **Performance Optimization**

- **Database Indexing**: Ensure proper indexes on trade tables
- **Connection Pooling**: Configure optimal database connections
- **Rate Limiting**: Respect exchange rate limits
- **Caching**: Use Redis for market data caching

## 📚 **Additional Resources**

- [API Reference](API_REFERENCE.md)
- [Strategy Development Guide](STRATEGY_DEVELOPMENT.md)
- [Risk Management Best Practices](RISK_MANAGEMENT.md)
- [Exchange Integration Guide](EXCHANGE_INTEGRATION.md)
- [Monitoring Setup](MONITORING_SETUP.md)

## 🤝 **Support**

- **Documentation**: [docs/](../docs/)
- **Issues**: [GitHub Issues](https://github.com/DimaJoyti/ai-agentic-crypto-browser/issues)
- **Discussions**: [GitHub Discussions](https://github.com/DimaJoyti/ai-agentic-crypto-browser/discussions)

---

**⚡ Ready to deploy 7 sophisticated trading bots with institutional-grade features!**
