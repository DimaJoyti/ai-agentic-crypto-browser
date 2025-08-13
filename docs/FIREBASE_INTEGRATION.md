# Firebase MCP Tools Integration

This document describes the Firebase integration for the AI Agentic Crypto Browser project using MCP (Model Context Protocol) tools.

## Overview

The Firebase integration provides comprehensive backend services for the crypto trading platform, including:

- **Authentication**: User management and secure authentication
- **Firestore**: NoSQL database for trading data, portfolios, and analytics
- **Realtime Database**: Live market data and real-time updates
- **Cloud Storage**: File storage for reports and user documents
- **Cloud Functions**: Serverless backend logic
- **Analytics**: User behavior and system performance tracking

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Go Backend    │    │   Firebase      │
│   (React/Next)  │◄──►│   API Server    │◄──►│   Services      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   MCP Firebase  │
                       │   Client        │
                       └─────────────────┘
```

## Setup Instructions

### 1. Prerequisites

- Firebase CLI installed: `npm install -g firebase-tools`
- Go 1.21+ with Firebase SDK dependencies
- Firebase project created at [Firebase Console](https://console.firebase.google.com/)

### 2. Quick Setup

Run the setup script:

```bash
./scripts/setup-firebase.sh
```

This script will:
- Initialize Firebase configuration
- Create security rules
- Set up database indexes
- Configure environment variables

### 3. Manual Setup

#### Create Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create a new project: "ai-agentic-crypto-browser"
3. Enable the following services:
   - Authentication
   - Firestore Database
   - Realtime Database
   - Cloud Storage
   - Cloud Functions
   - Analytics

#### Download Service Account Key

1. Go to Project Settings → Service Accounts
2. Generate new private key
3. Save as `configs/firebase-service-account.json`

#### Configure Environment Variables

Add to your `.env` file:

```bash
FIREBASE_PROJECT_ID=ai-agentic-crypto-browser
FIREBASE_DATABASE_URL=https://ai-agentic-crypto-browser-default-rtdb.firebaseio.com
FIREBASE_STORAGE_BUCKET=ai-agentic-crypto-browser.appspot.com
FIREBASE_CREDENTIALS_PATH=./configs/firebase-service-account.json
```

## API Endpoints

### Authentication

#### Create User
```http
POST /api/firebase/auth/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "display_name": "John Doe"
}
```

#### Get User
```http
GET /api/firebase/auth/users/{uid}
```

#### Verify Token
```http
POST /api/firebase/auth/verify
Content-Type: application/json

{
  "id_token": "firebase_id_token"
}
```

### Firestore Database

#### Create Document
```http
POST /api/firebase/firestore/{collection}
Content-Type: application/json

{
  "symbol": "BTCUSDT",
  "price": 45000,
  "timestamp": "2024-01-15T10:00:00Z"
}
```

#### Get Document
```http
GET /api/firebase/firestore/{collection}/{documentId}
```

#### Query Collection
```http
GET /api/firebase/firestore/{collection}?limit=10
```

#### Update Document
```http
PUT /api/firebase/firestore/{collection}/{documentId}
Content-Type: application/json

{
  "price": 46000,
  "updated_at": "2024-01-15T11:00:00Z"
}
```

#### Delete Document
```http
DELETE /api/firebase/firestore/{collection}/{documentId}
```

#### Batch Operations
```http
POST /api/firebase/firestore/batch
Content-Type: application/json

{
  "operations": [
    {
      "operation": "create",
      "collection": "trading_signals",
      "document_id": "signal_1",
      "data": {
        "symbol": "ETHUSDT",
        "signal": "BUY",
        "confidence": 0.85
      }
    },
    {
      "operation": "update",
      "collection": "market_data",
      "document_id": "btc_data",
      "data": {
        "price": 45500
      }
    }
  ]
}
```

### Realtime Database

#### Set Data
```http
POST /api/firebase/realtime/live_prices/BTCUSDT
Content-Type: application/json

{
  "price": 45000,
  "volume": 1234.56,
  "timestamp": 1642248000
}
```

#### Get Data
```http
GET /api/firebase/realtime/live_prices/BTCUSDT
```

#### Update Data
```http
PATCH /api/firebase/realtime/user_sessions/{userId}
Content-Type: application/json

{
  "last_activity": 1642248000,
  "status": "active"
}
```

#### Delete Data
```http
DELETE /api/firebase/realtime/alerts/{alertId}
```

### System Status

#### Get Firebase Status
```http
GET /api/firebase/status
```

Response:
```json
{
  "running": true,
  "project_id": "ai-agentic-crypto-browser",
  "services": {
    "auth": true,
    "firestore": true,
    "realtime_db": true
  },
  "emulators_enabled": false
}
```

## Data Models

### Trading Signal
```json
{
  "id": "signal_123",
  "symbol": "BTCUSDT",
  "signal": "BUY",
  "confidence": 0.85,
  "price": 45000,
  "target_price": 47000,
  "stop_loss": 43000,
  "timestamp": "2024-01-15T10:00:00Z",
  "strategy_id": "strategy_456",
  "user_id": "user_789",
  "metadata": {
    "indicators": ["RSI", "MACD"],
    "timeframe": "1h"
  }
}
```

### User Portfolio
```json
{
  "id": "portfolio_123",
  "user_id": "user_789",
  "name": "Main Portfolio",
  "total_value": 50000,
  "cash_balance": 10000,
  "positions": [
    {
      "symbol": "BTCUSDT",
      "quantity": 1.5,
      "avg_price": 44000,
      "current_value": 67500,
      "pnl": 1500
    }
  ],
  "performance": {
    "total_return": 0.15,
    "daily_return": 0.02,
    "max_drawdown": 0.08
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Market Data
```json
{
  "id": "market_btc_123",
  "symbol": "BTCUSDT",
  "exchange": "binance",
  "price": 45000,
  "volume": 1234567.89,
  "high_24h": 46000,
  "low_24h": 44000,
  "change_24h": 1000,
  "change_percent_24h": 2.27,
  "timestamp": "2024-01-15T10:00:00Z",
  "indicators": {
    "rsi": 65.5,
    "macd": 0.5,
    "bollinger_upper": 46500,
    "bollinger_lower": 43500
  }
}
```

## Security Rules

### Firestore Rules

The Firestore security rules ensure:
- Users can only access their own portfolios and preferences
- Trading signals are readable by authenticated users
- Admin-only access to system logs and audit trails
- Role-based access control using custom claims

### Realtime Database Rules

The Realtime Database rules provide:
- Public read access to live market prices
- User-specific access to orders and sessions
- Admin-only write access to system status

### Authentication

Custom claims are used for role-based access:
- `admin`: Full system access
- `trader`: Can execute trades and manage strategies
- `analyst`: Can view analytics and create reports
- `viewer`: Read-only access to market data

## Development

### Local Development with Emulators

Start Firebase emulators:
```bash
firebase emulators:start
```

This starts:
- Auth Emulator: http://localhost:9099
- Firestore Emulator: http://localhost:8080
- Realtime Database Emulator: http://localhost:9000
- Storage Emulator: http://localhost:9199
- Functions Emulator: http://localhost:5001
- Emulator UI: http://localhost:4000

### Environment Configuration

Set environment variables for emulator mode:
```bash
export FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
export FIRESTORE_EMULATOR_HOST=localhost:8080
export FIREBASE_DATABASE_EMULATOR_HOST=localhost:9000
export FIREBASE_STORAGE_EMULATOR_HOST=localhost:9199
```

## Deployment

### Deploy Rules and Indexes
```bash
firebase deploy --only firestore:rules,firestore:indexes,database:rules,storage:rules
```

### Deploy Functions
```bash
firebase deploy --only functions
```

### Deploy Everything
```bash
firebase deploy
```

## Monitoring and Analytics

### Performance Monitoring

Firebase Performance Monitoring tracks:
- API response times
- Database query performance
- Trading execution latency
- Market data processing speed

### Analytics Events

Custom events tracked:
- `user_login`: User authentication
- `trade_executed`: Trading activity
- `strategy_created`: Strategy management
- `portfolio_viewed`: User engagement
- `alert_triggered`: System notifications

### Error Reporting

Firebase Crashlytics captures:
- Application crashes
- Non-fatal errors
- Performance issues
- Custom error logs

## Best Practices

### Data Structure
- Use subcollections for hierarchical data
- Implement proper indexing for queries
- Use batch operations for multiple writes
- Implement data validation rules

### Security
- Always use authentication
- Implement proper security rules
- Use custom claims for role-based access
- Regularly audit access patterns

### Performance
- Use real-time listeners sparingly
- Implement proper pagination
- Cache frequently accessed data
- Monitor query performance

### Cost Optimization
- Implement data retention policies
- Use appropriate storage classes
- Monitor usage and billing
- Optimize query patterns

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify service account key
   - Check project ID configuration
   - Ensure proper permissions

2. **Permission Denied**
   - Review security rules
   - Check user authentication
   - Verify custom claims

3. **Performance Issues**
   - Check database indexes
   - Review query patterns
   - Monitor resource usage

### Debug Mode

Enable debug logging:
```bash
export FIREBASE_DEBUG=true
export GOOGLE_CLOUD_PROJECT=ai-agentic-crypto-browser
```

## Support

For issues and questions:
- Check Firebase documentation
- Review security rules
- Monitor Firebase console
- Check application logs
