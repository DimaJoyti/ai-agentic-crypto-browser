# Autonomous Trading API Documentation

## 🎯 Overview

This document provides comprehensive API documentation for the autonomous trading and DeFi operations system. The API enables users to create and manage trading portfolios, execute autonomous trading strategies, interact with DeFi protocols, and perform intelligent portfolio rebalancing.

## 🔗 Base URL

```
http://localhost:8084/api/v1
```

## 🔐 Authentication

All endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## 📊 Trading Engine Endpoints

### Create Portfolio

Create a new autonomous trading portfolio.

**Endpoint:** `POST /web3/trading/portfolio`

**Request Body:**
```json
{
  "name": "AI Growth Portfolio",
  "initial_balance": "25000.00",
  "risk_profile": {
    "level": "moderate",
    "max_position_size": "0.10",
    "max_daily_loss": "0.05",
    "stop_loss_percentage": "0.10",
    "take_profit_percentage": "0.20",
    "allowed_strategies": ["momentum", "mean_reversion"]
  }
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user-uuid",
  "name": "AI Growth Portfolio",
  "total_value": "25000.00",
  "available_balance": "25000.00",
  "invested_amount": "0.00",
  "total_pnl": "0.00",
  "daily_pnl": "0.00",
  "holdings": {},
  "active_positions": [],
  "trading_strategies": [],
  "risk_profile": {
    "level": "moderate",
    "max_position_size": "0.10",
    "max_daily_loss": "0.05",
    "stop_loss_percentage": "0.10",
    "take_profit_percentage": "0.20",
    "allowed_strategies": ["momentum", "mean_reversion"]
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Get Portfolio

Retrieve portfolio details with real-time data.

**Endpoint:** `GET /web3/trading/portfolio/{portfolio_id}`

**Response:**
```json
{
  "id": "portfolio-uuid",
  "user_id": "user-uuid",
  "name": "AI Growth Portfolio",
  "total_value": "27500.00",
  "available_balance": "12500.00",
  "invested_amount": "15000.00",
  "total_pnl": "2500.00",
  "daily_pnl": "150.00",
  "holdings": {
    "ETH": {
      "token_address": "0x...",
      "token_symbol": "ETH",
      "amount": "6.25",
      "average_price": "2400.00",
      "current_price": "2450.00",
      "value": "15312.50",
      "pnl": "312.50",
      "pnl_percentage": "2.08",
      "last_updated": "2024-01-15T15:45:00Z"
    }
  },
  "active_positions": [
    {
      "id": "position-uuid",
      "strategy_name": "momentum_strategy",
      "token_address": "0x...",
      "token_symbol": "ETH",
      "amount": "6.25",
      "entry_price": "2400.00",
      "current_price": "2450.00",
      "unrealized_pnl": "312.50",
      "stop_loss": "2160.00",
      "take_profit": "2880.00",
      "status": "open",
      "opened_at": "2024-01-15T14:30:00Z"
    }
  ]
}
```

### Start Trading

Start autonomous trading for a portfolio (global engine already running).

**Endpoint:** `POST /web3/trading/portfolio/{portfolio_id}/start`

**Response:**
```json
{
  "message": "Trading engine is already running globally",
  "status": "success"
}
```

### Stop Trading

Stop autonomous trading engine.

**Endpoint:** `POST /web3/trading/portfolio/{portfolio_id}/stop`

**Response:**
```json
{
  "message": "Trading stopped successfully",
  "status": "success"
}
```

### Get Active Positions

Retrieve all active positions for a portfolio.

**Endpoint:** `GET /web3/trading/positions/{portfolio_id}`

**Response:**
```json
{
  "portfolio_id": "portfolio-uuid",
  "positions": [
    {
      "id": "position-uuid-1",
      "user_id": "user-uuid",
      "strategy_name": "momentum_strategy",
      "token_address": "0x...",
      "token_symbol": "ETH",
      "amount": "6.25",
      "entry_price": "2400.00",
      "current_price": "2450.00",
      "unrealized_pnl": "312.50",
      "realized_pnl": "0.00",
      "stop_loss": "2160.00",
      "take_profit": "2880.00",
      "status": "open",
      "opened_at": "2024-01-15T14:30:00Z",
      "updated_at": "2024-01-15T15:45:00Z"
    }
  ]
}
```

### Close Position

Manually close a trading position.

**Endpoint:** `POST /web3/trading/positions/{position_id}/close`

**Request Body:**
```json
{
  "reason": "Manual close - taking profits"
}
```

**Response:**
```json
{
  "message": "Position closed successfully",
  "position_id": "position-uuid",
  "reason": "Manual close - taking profits"
}
```

## 🏦 DeFi Protocol Endpoints

### Get All Protocols

List all supported DeFi protocols.

**Endpoint:** `GET /web3/defi/protocols`

**Response:**
```json
{
  "protocols": {
    "uniswap_v3": {
      "id": "uniswap_v3",
      "name": "Uniswap V3",
      "type": "dex",
      "chain_id": 1,
      "address": "0xE592427A0AEce92De3Edee1F18E0157C05861564",
      "tvl": "5000000000.00",
      "apy": "0.15",
      "fees": "0.003",
      "risk_score": 25,
      "is_active": true,
      "pools": {
        "USDC_ETH": {
          "id": "uniswap_v3_usdc_eth",
          "name": "USDC/ETH 0.3%",
          "token_a": "USDC",
          "token_b": "ETH",
          "total_liquidity": "200000000.00",
          "apy": "0.18",
          "volume_24h": "50000000.00",
          "risk_level": "medium"
        }
      }
    }
  }
}
```

### Get Protocol Details

Get detailed information about a specific protocol.

**Endpoint:** `GET /web3/defi/protocols/{protocol_id}`

**Response:**
```json
{
  "id": "compound",
  "name": "Compound",
  "type": "lending",
  "chain_id": 1,
  "address": "0x3d9819210A31b4961b30EF54bE2aeD79B9c9Cd3B",
  "tvl": "3000000000.00",
  "apy": "0.08",
  "fees": "0.001",
  "risk_score": 15,
  "is_active": true,
  "pools": {
    "USDC": {
      "id": "compound_usdc",
      "name": "cUSDC",
      "token_a": "USDC",
      "reserve_a": "500000000.00",
      "total_liquidity": "500000000.00",
      "apy": "0.06",
      "risk_level": "low"
    }
  }
}
```

### Get Yield Opportunities

Discover the best yield farming opportunities.

**Endpoint:** `GET /web3/defi/opportunities`

**Query Parameters:**
- `min_apy` (optional): Minimum APY threshold (default: 0.01)
- `max_risk` (optional): Maximum risk level (very_low, low, medium, high, critical)

**Example:** `GET /web3/defi/opportunities?min_apy=0.05&max_risk=medium`

**Response:**
```json
{
  "opportunities": [
    {
      "protocol_id": "uniswap_v3",
      "protocol_name": "Uniswap V3",
      "pool_id": "uniswap_v3_usdc_eth",
      "pool_name": "USDC/ETH 0.3%",
      "apy": "0.18",
      "tvl": "200000000.00",
      "risk_level": "medium",
      "token_a": "USDC",
      "token_b": "ETH",
      "fees": "0.003",
      "impermanent_loss": "0.02"
    },
    {
      "protocol_id": "aave",
      "protocol_name": "Aave",
      "pool_id": "aave_eth",
      "pool_name": "aETH",
      "apy": "0.09",
      "tvl": "400000000.00",
      "risk_level": "low",
      "token_a": "ETH",
      "token_b": "",
      "fees": "0.0005",
      "impermanent_loss": "0.00"
    }
  ],
  "filters": {
    "min_apy": "0.05",
    "max_risk": "medium"
  }
}
```

## ⚖️ Portfolio Rebalancing Endpoints

### Create Rebalancing Strategy

Create a new portfolio rebalancing strategy.

**Endpoint:** `POST /web3/rebalance/strategy`

**Request Body:**
```json
{
  "portfolio_id": "portfolio-uuid",
  "name": "Balanced Growth Strategy",
  "type": "fixed",
  "target_allocations": {
    "ETH": "0.40",
    "BTC": "0.30",
    "USDC": "0.30"
  }
}
```

**Response:**
```json
{
  "id": "strategy-uuid",
  "portfolio_id": "portfolio-uuid",
  "name": "Balanced Growth Strategy",
  "type": "fixed",
  "target_allocations": {
    "ETH": "0.40",
    "BTC": "0.30",
    "USDC": "0.30"
  },
  "constraints": [],
  "trigger_conditions": [
    {
      "type": "drift",
      "threshold": "0.05",
      "condition": "deviation"
    },
    {
      "type": "time",
      "threshold": "24",
      "condition": "hours"
    }
  ],
  "is_active": true,
  "created_at": "2024-01-15T16:00:00Z"
}
```

### Get Rebalancing Strategy

Retrieve rebalancing strategy for a portfolio.

**Endpoint:** `GET /web3/rebalance/strategy/{portfolio_id}`

**Response:**
```json
{
  "portfolio_id": "portfolio-uuid",
  "message": "Rebalance strategy retrieval not implemented yet"
}
```

### Execute Rebalancing

Manually trigger portfolio rebalancing.

**Endpoint:** `POST /web3/rebalance/execute/{portfolio_id}`

**Response:**
```json
{
  "message": "Portfolio rebalanced successfully",
  "portfolio_id": "portfolio-uuid"
}
```

## 🔧 Enhanced Web3 Endpoints

### Create Enhanced Transaction

Create a transaction with advanced features like gas optimization and MEV protection.

**Endpoint:** `POST /web3/enhanced/transaction`

**Request Body:**
```json
{
  "wallet_id": "wallet-uuid",
  "to_address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
  "value": "1000000000000000000",
  "data": "0x",
  "gas_optimization": true,
  "mev_protection": true,
  "priority": "high"
}
```

**Response:**
```json
{
  "transaction_id": "tx-uuid",
  "transaction": {
    "id": "tx-uuid",
    "hash": "0x...",
    "from": "0x...",
    "to": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
    "value": "1000000000000000000",
    "gas_limit": 21000,
    "gas_price": "20000000000",
    "status": "pending"
  },
  "tx_hash": "0x...",
  "status": "pending"
}
```

## 📈 Risk Levels

The system uses the following risk levels for DeFi protocols and strategies:

- **very_low**: Minimal risk, stable protocols (0-20 risk score)
- **low**: Low risk, established protocols (21-40 risk score)
- **medium**: Moderate risk, proven protocols (41-60 risk score)
- **high**: High risk, newer protocols (61-80 risk score)
- **critical**: Very high risk, experimental protocols (81-100 risk score)

## 🎯 Trading Strategy Types

- **momentum**: Trend-following strategy based on RSI and volume
- **mean_reversion**: Contrarian strategy using Bollinger Bands
- **arbitrage**: Cross-DEX price difference exploitation

## ⚖️ Rebalancing Strategy Types

- **fixed**: Traditional percentage-based allocation
- **dynamic**: Market condition-based allocation
- **risk_parity**: Risk-weighted allocation
- **momentum**: Trend-following allocation
- **mean_revert**: Contrarian allocation strategy

## 🚨 Error Responses

All endpoints return standard HTTP status codes with JSON error responses:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

Common status codes:
- `400`: Bad Request - Invalid input parameters
- `401`: Unauthorized - Missing or invalid authentication
- `404`: Not Found - Resource not found
- `500`: Internal Server Error - Server-side error

## 📊 Rate Limits

- **Portfolio Operations**: 10 requests per minute
- **Position Management**: 20 requests per minute
- **DeFi Queries**: 30 requests per minute
- **Rebalancing**: 5 requests per minute

## 🔐 Security Considerations

- All sensitive operations require user authentication
- Portfolio access is restricted to the owner
- Position modifications are logged for audit trails
- Risk limits are enforced at multiple levels
- Emergency stop mechanisms are available for critical situations
