# Real-time Data and Monitoring API Documentation

## 🔍 Overview

This document provides comprehensive API documentation for the real-time data and monitoring features of the autonomous trading system, including live market data feeds, portfolio analytics, system monitoring, and alert management.

## 🔗 Base URL

```
http://localhost:8084/api/v1
```

## 🔐 Authentication

All endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## 📊 Real-time Market Data

### Get Market Data Status

Check the status of real-time market data connections.

**Endpoint:** `GET /web3/realtime/market/status`

**Response:**
```json
{
  "status": "active",
  "connections": {
    "binance": {
      "exchange": "binance",
      "is_connected": true,
      "last_ping": "2024-01-15T10:30:00Z",
      "last_pong": "2024-01-15T10:30:05Z",
      "reconnects": 0,
      "message_count": 15420,
      "error_count": 2
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Subscribe to Market Data Stream

Subscribe to real-time market data updates for a specific symbol.

**Endpoint:** `GET /web3/realtime/market/subscribe/{symbol}`

**Example:** `GET /web3/realtime/market/subscribe/BTCUSDT`

**Response:** Server-Sent Events (SSE) stream

```
data: {"type":"connected","symbol":"BTCUSDT"}

data: {"exchange":"binance","symbol":"BTCUSDT","type":"ticker","price":"45250.00","volume":"1250.50","bid":"45249.50","ask":"45250.50","high_24h":"46100.00","low_24h":"44800.00","change_24h":"1.2","timestamp":"2024-01-15T10:30:00Z"}

data: {"exchange":"binance","symbol":"BTCUSDT","type":"trade","price":"45251.00","volume":"0.5","timestamp":"2024-01-15T10:30:01Z"}
```

**Usage Example:**
```javascript
const eventSource = new EventSource('/web3/realtime/market/subscribe/BTCUSDT');

eventSource.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Market update:', data);
};

eventSource.onerror = function(event) {
  console.error('Connection error:', event);
};
```

## 📈 Portfolio Analytics

### Get Portfolio Analytics

Retrieve comprehensive analytics for a portfolio.

**Endpoint:** `GET /web3/analytics/portfolio/{portfolio_id}`

**Response:**
```json
{
  "portfolio_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "user-uuid",
  "name": "AI Growth Portfolio",
  "total_value": "27500.00",
  "total_pnl": "2500.00",
  "total_pnl_percent": "10.00",
  "daily_pnl": "150.00",
  "weekly_pnl": "800.00",
  "monthly_pnl": "2500.00",
  "max_drawdown": "-5.2",
  "sharpe_ratio": "1.85",
  "sortino_ratio": "2.12",
  "volatility": "0.15",
  "beta": "1.05",
  "alpha": "0.02",
  "holdings": [
    {
      "symbol": "ETH",
      "token_address": "0x...",
      "amount": "6.25",
      "average_price": "2400.00",
      "current_price": "2450.00",
      "value": "15312.50",
      "pnl": "312.50",
      "pnl_percent": "2.08",
      "weight": "55.68",
      "day_change": "1.5",
      "week_change": "8.2",
      "month_change": "15.7",
      "last_updated": "2024-01-15T10:30:00Z"
    }
  ],
  "positions": [
    {
      "id": "position-uuid",
      "symbol": "ETH",
      "strategy": "momentum",
      "side": "long",
      "amount": "6.25",
      "entry_price": "2400.00",
      "current_price": "2450.00",
      "unrealized_pnl": "312.50",
      "realized_pnl": "0.00",
      "stop_loss": "2160.00",
      "take_profit": "2880.00",
      "duration": "24h0m0s",
      "status": "open",
      "opened_at": "2024-01-14T10:30:00Z",
      "last_updated": "2024-01-15T10:30:00Z"
    }
  ],
  "performance": {
    "daily": [
      {
        "timestamp": "2024-01-15T00:00:00Z",
        "value": "27500.00",
        "pnl": "150.00",
        "pnl_percent": "0.55",
        "drawdown": "0.00",
        "volume": "5000.00",
        "trades": 3
      }
    ],
    "weekly": [...],
    "monthly": [...]
  },
  "risk_metrics": {
    "var_95": "-0.025",
    "var_99": "-0.045",
    "cvar_95": "-0.032",
    "max_drawdown": "-0.052",
    "downside_risk": "0.08",
    "upside_capture": "1.05",
    "downside_capture": "0.95",
    "correlation": "0.75",
    "risk_score": "35.0",
    "risk_grade": "B"
  },
  "last_updated": "2024-01-15T10:30:00Z"
}
```

### Get Portfolio Performance

Retrieve detailed performance metrics for a portfolio.

**Endpoint:** `GET /web3/analytics/portfolio/{portfolio_id}/performance`

**Response:**
```json
{
  "portfolio_id": "550e8400-e29b-41d4-a716-446655440000",
  "performance": {
    "daily": [...],
    "weekly": [...],
    "monthly": [...]
  },
  "risk_metrics": {
    "var_95": "-0.025",
    "var_99": "-0.045",
    "max_drawdown": "-0.052",
    "risk_score": "35.0",
    "risk_grade": "B"
  },
  "sharpe_ratio": "1.85",
  "sortino_ratio": "2.12",
  "max_drawdown": "-5.2",
  "volatility": "0.15"
}
```

### Compare Portfolios

Compare performance metrics across multiple portfolios.

**Endpoint:** `GET /web3/analytics/portfolio/compare`

**Query Parameters:**
- `portfolio_ids` (required): Comma-separated list of portfolio UUIDs

**Example:** `GET /web3/analytics/portfolio/compare?portfolio_ids=uuid1,uuid2,uuid3`

**Response:**
```json
{
  "550e8400-e29b-41d4-a716-446655440000": {
    "name": "AI Growth Portfolio",
    "total_value": "27500.00",
    "total_pnl_percent": "10.00",
    "sharpe_ratio": "1.85",
    "max_drawdown": "-5.2",
    "risk_grade": "B"
  },
  "660f9511-f3ac-52e5-b827-557766551111": {
    "name": "Conservative Portfolio",
    "total_value": "15200.00",
    "total_pnl_percent": "5.20",
    "sharpe_ratio": "1.42",
    "max_drawdown": "-2.8",
    "risk_grade": "A"
  }
}
```

## 🖥️ System Monitoring

### Get System Health

Check overall system health status.

**Endpoint:** `GET /web3/monitoring/health`

**Response:**
```json
{
  "status": "healthy",
  "score": 92.5,
  "components": {
    "cpu": "healthy",
    "memory": "healthy",
    "application": "healthy",
    "database": "healthy",
    "websocket": "healthy"
  },
  "issues": [],
  "last_check": "2024-01-15T10:30:00Z"
}
```

### Get System Metrics

Retrieve comprehensive system performance metrics.

**Endpoint:** `GET /web3/monitoring/metrics`

**Response:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "cpu": {
    "usage_percent": 45.2,
    "load_average_1m": 1.2,
    "load_average_5m": 1.1,
    "load_average_15m": 1.0,
    "cores": 8,
    "goroutines": 150,
    "cgo_calls": 1250
  },
  "memory": {
    "total_bytes": 8589934592,
    "used_bytes": 3221225472,
    "free_bytes": 5368709120,
    "usage_percent": 37.5,
    "heap_bytes": 134217728,
    "stack_bytes": 8388608,
    "gc_pause_ms": 2.5,
    "gc_count": 45
  },
  "disk": {
    "total_bytes": 107374182400,
    "used_bytes": 32212254720,
    "free_bytes": 75161927680,
    "usage_percent": 30.0,
    "read_ops": 1000,
    "write_ops": 500,
    "read_bytes": 1048576,
    "write_bytes": 524288
  },
  "network": {
    "bytes_received": 104857600,
    "bytes_sent": 52428800,
    "packets_received": 10000,
    "packets_sent": 5000,
    "connections": 100,
    "active_sockets": 50,
    "drop_rate": 0.01
  },
  "application": {
    "request_count": 10000,
    "error_count": 50,
    "error_rate": 0.5,
    "avg_response_time": "100ms",
    "p95_response_time": "200ms",
    "p99_response_time": "500ms",
    "active_users": 250,
    "throughput_rps": 100.0,
    "cache_hit_rate": 0.85,
    "queue_length": 10
  },
  "trading": {
    "active_portfolios": 50,
    "active_positions": 150,
    "total_trades": 1000,
    "successful_trades": 850,
    "failed_trades": 150,
    "trade_success_rate": 85.0,
    "avg_execution_time": "500ms",
    "total_volume": "1000000.00",
    "total_pnl": "50000.00",
    "risk_alerts": 5,
    "strategy_performance": {
      "momentum": 85.5,
      "mean_reversion": 78.2,
      "arbitrage": 92.1
    }
  },
  "database": {
    "connections_active": 20,
    "connections_idle": 30,
    "connections_total": 50,
    "queries_per_second": 100.0,
    "avg_query_time": "10ms",
    "slow_queries": 2,
    "deadlock_count": 0,
    "cache_hit_ratio": 0.95,
    "replication_lag": "100ms"
  },
  "websocket": {
    "total_connections": 100,
    "active_connections": 85,
    "messages_received": 10000,
    "messages_sent": 8000,
    "connection_errors": 5,
    "reconnect_count": 10,
    "avg_latency": 50.0,
    "data_throughput": 1048576
  },
  "health": {
    "status": "healthy",
    "score": 92.5,
    "components": {
      "cpu": "healthy",
      "memory": "healthy",
      "application": "healthy"
    },
    "issues": [],
    "last_check": "2024-01-15T10:30:00Z"
  }
}
```

### Get System Status

Get combined health and alert status.

**Endpoint:** `GET /web3/monitoring/status`

**Response:**
```json
{
  "health": {
    "status": "healthy",
    "score": 92.5,
    "components": {
      "cpu": "healthy",
      "memory": "healthy",
      "application": "healthy"
    },
    "issues": [],
    "last_check": "2024-01-15T10:30:00Z"
  },
  "alerts": [
    {
      "id": "cpu_usage_1642248600",
      "type": "system",
      "severity": "warning",
      "title": "High CPU Usage",
      "description": "CPU usage is 85.2%, exceeding threshold of 80.0%",
      "metric": "cpu_usage",
      "value": 85.2,
      "threshold": 80.0,
      "timestamp": "2024-01-15T10:30:00Z",
      "resolved": false
    }
  ],
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🚨 Alert Management

### Get Alerts

Retrieve recent alerts with optional limit.

**Endpoint:** `GET /web3/alerts`

**Query Parameters:**
- `limit` (optional): Maximum number of alerts to return (default: 50)

**Example:** `GET /web3/alerts?limit=20`

**Response:**
```json
{
  "alerts": [
    {
      "id": "alert-uuid-1",
      "rule_id": "high_cpu_usage",
      "title": "High CPU Usage",
      "message": "CPU usage exceeds threshold: 85.2% > 80.0%",
      "severity": "warning",
      "metric": "cpu_usage_percent",
      "value": "85.2",
      "threshold": "80.0",
      "timestamp": "2024-01-15T10:30:00Z",
      "resolved": false,
      "channels": ["email", "slack"],
      "metadata": {}
    }
  ],
  "count": 1,
  "limit": 20
}
```

### Get Active Alerts

Retrieve only unresolved alerts.

**Endpoint:** `GET /web3/alerts/active`

**Response:**
```json
{
  "alerts": [
    {
      "id": "alert-uuid-1",
      "rule_id": "high_cpu_usage",
      "title": "High CPU Usage",
      "message": "CPU usage exceeds threshold: 85.2% > 80.0%",
      "severity": "warning",
      "metric": "cpu_usage_percent",
      "value": "85.2",
      "threshold": "80.0",
      "timestamp": "2024-01-15T10:30:00Z",
      "resolved": false,
      "channels": ["email", "slack"]
    }
  ],
  "count": 1
}
```

### Resolve Alert

Mark an alert as resolved.

**Endpoint:** `POST /web3/alerts/{alert_id}/resolve`

**Response:**
```json
{
  "message": "Alert resolved successfully",
  "alert_id": "alert-uuid-1"
}
```

### Subscribe to Alert Stream

Subscribe to real-time alert notifications.

**Endpoint:** `GET /web3/alerts/subscribe/{topic}`

**Topics:**
- `all` - All alerts
- `severity_critical` - Critical alerts only
- `severity_warning` - Warning alerts only
- `metric_cpu_usage` - CPU usage alerts only

**Example:** `GET /web3/alerts/subscribe/severity_critical`

**Response:** Server-Sent Events (SSE) stream

```
data: {"type":"connected","topic":"severity_critical"}

data: {"id":"alert-uuid-2","rule_id":"high_error_rate","title":"High Error Rate","message":"Error rate exceeds threshold: 6.5% > 5.0%","severity":"critical","metric":"error_rate_percent","value":"6.5","threshold":"5.0","timestamp":"2024-01-15T10:35:00Z","resolved":false,"channels":["email","slack","webhook"]}
```

## 📊 Performance Metrics

### Response Times
- **Market Data Status**: <50ms average
- **Portfolio Analytics**: <200ms average
- **System Metrics**: <100ms average
- **Alert Operations**: <50ms average

### Real-time Capabilities
- **Market Data Updates**: <100ms latency from exchange
- **Alert Notifications**: <1 second from trigger to delivery
- **System Monitoring**: 30-second collection intervals
- **Portfolio Analytics**: 5-minute update intervals

### Scalability
- **Concurrent Connections**: 1000+ WebSocket connections
- **Market Data Throughput**: 10,000+ messages/second
- **Alert Processing**: 1000+ alerts/minute
- **Analytics Queries**: 100+ concurrent requests

## 🔧 Configuration

### Market Data Configuration
```json
{
  "exchanges": [
    {
      "name": "binance",
      "ws_url": "wss://stream.binance.com:9443/ws",
      "symbols": ["BTCUSDT", "ETHUSDT", "ADAUSDT"],
      "channels": ["ticker", "trade"],
      "enabled": true
    }
  ],
  "reconnect_delay": "5s",
  "ping_interval": "30s",
  "max_reconnects": 10,
  "buffer_size": 1000,
  "enable_heartbeat": true
}
```

### Monitoring Configuration
```json
{
  "collection_interval": "30s",
  "retention_period": "24h",
  "alert_thresholds": {
    "cpu_threshold": 80.0,
    "memory_threshold": 85.0,
    "error_rate_threshold": 5.0
  },
  "enable_profiling": true,
  "enable_tracing": true
}
```

### Alert Configuration
```json
{
  "max_history_size": 1000,
  "default_cooldown": "5m",
  "enable_email": true,
  "enable_webhook": true,
  "enable_slack": true,
  "enable_push_notifications": true
}
```

## 🚨 Error Handling

All endpoints return standard HTTP status codes with JSON error responses:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

Common status codes:
- `400`: Bad Request - Invalid parameters
- `401`: Unauthorized - Missing or invalid authentication
- `404`: Not Found - Resource not found
- `500`: Internal Server Error - Server-side error
- `503`: Service Unavailable - Service temporarily unavailable
