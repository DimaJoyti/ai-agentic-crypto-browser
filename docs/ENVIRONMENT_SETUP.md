# Environment Variables Setup Guide

This guide explains how to properly configure environment variables for the AI Agentic Crypto Browser project.

## 📁 Environment Files Structure

```
ai-agentic-crypto-browser/
├── .env                    # Backend environment variables
├── .env.example           # Backend environment template
├── web/.env.local         # Frontend environment variables
└── web/.env.example       # Frontend environment template
```

## 🚀 Quick Setup

### 1. Backend Environment Setup

Copy the example file and customize:
```bash
cp .env.example .env
```

### 2. Frontend Environment Setup

```bash
cd web
cp .env.example .env.local
```

### 3. Required API Keys

For development, you can use the demo/placeholder values provided. For production, you'll need to obtain real API keys.

## 🔑 API Keys Configuration

### Core Services

#### Database & Cache
```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/agentic_browser?sslmode=disable
REDIS_URL=redis://localhost:6379
```

#### Authentication
```bash
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h
```

### Solana Integration

#### RPC Endpoints
```bash
# Development (Devnet)
SOLANA_RPC_URL=https://api.devnet.solana.com
SOLANA_WS_URL=wss://api.devnet.solana.com

# Production (Mainnet)
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com
```

#### DeFi Protocols
```bash
# Jupiter DEX Aggregator
JUPITER_API_URL=https://quote-api.jup.ag/v6

# Raydium AMM
RAYDIUM_API_URL=https://api.raydium.io/v2

# Orca DEX
ORCA_API_URL=https://api.orca.so

# Marinade Finance
MARINADE_API_URL=https://api.marinade.finance
```

#### NFT Marketplaces
```bash
# Magic Eden
MAGIC_EDEN_API_URL=https://api-mainnet.magiceden.dev/v2
MAGIC_EDEN_API_KEY=your_magic_eden_api_key_here

# Tensor
TENSOR_API_URL=https://api.tensor.trade
TENSOR_API_KEY=your_tensor_api_key_here
```

### External APIs

#### Market Data
```bash
# CoinGecko
COINGECKO_API_URL=https://api.coingecko.com/api/v3
COINGECKO_API_KEY=your_coingecko_api_key_here

# DeFiLlama
DEFILLAMA_API_URL=https://api.llama.fi
```

#### Trading
```bash
# Binance (Testnet for development)
BINANCE_API_KEY=your_binance_api_key_here
BINANCE_SECRET_KEY=your_binance_secret_key_here
BINANCE_TESTNET=true
BINANCE_BASE_URL=https://testnet.binance.vision
BINANCE_WS_URL=wss://testnet.binance.vision/ws
```

### Payment & Banking

#### Stripe (Demo keys for development)
```bash
STRIPE_PUBLISHABLE_KEY=pk_test_demo_key_1234567890abcdef
STRIPE_SECRET_KEY=sk_test_demo_secret_1234567890abcdef
STRIPE_WEBHOOK_SECRET=whsec_demo_webhook_secret_12345
```

#### Plaid Banking
```bash
PLAID_CLIENT_ID=demo_plaid_client_id
PLAID_SECRET=demo_plaid_secret
PLAID_ENV=sandbox
```

### Communication Services

#### Email (SendGrid)
```bash
SENDGRID_API_KEY=demo_sendgrid_api_key
SENDGRID_FROM_EMAIL=noreply@yourdomain.com
```

#### SMS (Twilio)
```bash
TWILIO_ACCOUNT_SID=demo_twilio_account_sid
TWILIO_AUTH_TOKEN=demo_twilio_auth_token
TWILIO_PHONE_NUMBER=+1234567890
```

## 🎛️ Feature Flags

Control which features are enabled:

### Backend Feature Flags
```bash
# Solana Features
ENABLE_SOLANA_DEFI=true
ENABLE_SOLANA_NFT=true
ENABLE_SOLANA_STAKING=true
ENABLE_SOLANA_SWAPS=true

# Trading Features
ENABLE_ALGORITHMIC_TRADING=true
ENABLE_SOCIAL_TRADING=true
ENABLE_COPY_TRADING=true

# Security Features
ENABLE_MFA=true
ENABLE_FRAUD_DETECTION=true
ENABLE_KYC_VERIFICATION=false  # Disable in development

# Institutional Features
ENABLE_PRIME_BROKERAGE=false   # Disable in development
ENABLE_COMPLIANCE_REPORTING=false
ENABLE_INSTITUTIONAL_API=false

# Fiat Features
ENABLE_FIAT_ONRAMP=false       # Disable in development
ENABLE_FIAT_OFFRAMP=false
ENABLE_BANK_INTEGRATION=false
```

### Frontend Feature Flags
```bash
# Add NEXT_PUBLIC_ prefix for frontend
NEXT_PUBLIC_ENABLE_SOLANA_DEFI=true
NEXT_PUBLIC_ENABLE_SOLANA_NFT=true
NEXT_PUBLIC_ENABLE_ALGORITHMIC_TRADING=true
NEXT_PUBLIC_ENABLE_SOCIAL_TRADING=true

# Development flags
NEXT_PUBLIC_DEMO_MODE=true
NEXT_PUBLIC_USE_MOCK_DATA=true
NEXT_PUBLIC_ENABLE_PAPER_TRADING=true
```

## 🔒 Security Best Practices

### 1. Never Commit Real API Keys
- Always use `.env` files (which are gitignored)
- Use demo/placeholder values in example files
- Rotate API keys regularly in production

### 2. Environment-Specific Configuration
- **Development**: Use testnet/sandbox APIs
- **Staging**: Use testnet with production-like data
- **Production**: Use mainnet with real API keys

### 3. API Key Patterns
- Use descriptive prefixes: `demo_`, `test_`, `prod_`
- Never use patterns that trigger GitHub secret detection
- Avoid `sk_live_`, `pk_live_`, `rk_live_` patterns

## 🚦 Environment Validation

The application validates required environment variables on startup:

### Backend Validation
```go
// Required variables are checked in config/config.go
requiredVars := []string{
    "DATABASE_URL",
    "REDIS_URL", 
    "JWT_SECRET",
}
```

### Frontend Validation
```typescript
// Required variables are checked in next.config.js
const requiredEnvVars = [
    'NEXT_PUBLIC_API_URL',
    'NEXT_PUBLIC_SOLANA_RPC_URL',
]
```

## 🐛 Troubleshooting

### Common Issues

1. **Missing Environment Variables**
   ```bash
   Error: Required environment variable DATABASE_URL is not set
   ```
   **Solution**: Copy from `.env.example` and set the value

2. **Invalid API URLs**
   ```bash
   Error: Failed to connect to Solana RPC
   ```
   **Solution**: Check `SOLANA_RPC_URL` is correct and accessible

3. **Frontend Can't Connect to Backend**
   ```bash
   Error: Network Error
   ```
   **Solution**: Verify `NEXT_PUBLIC_API_URL` matches backend port

### Debug Mode
Enable debug logging:
```bash
DEBUG=true
LOG_LEVEL=debug
NEXT_PUBLIC_DEBUG=true
```

## 📚 Additional Resources

- [Solana RPC Endpoints](https://docs.solana.com/cluster/rpc-endpoints)
- [Jupiter API Documentation](https://docs.jup.ag/)
- [Magic Eden API](https://api.magiceden.dev/)
- [CoinGecko API](https://www.coingecko.com/en/api)
- [Binance API](https://binance-docs.github.io/apidocs/)

## 🤝 Contributing

When adding new environment variables:

1. Add to both `.env.example` and `web/.env.example`
2. Document in this guide
3. Add validation in the application
4. Use demo/safe values in examples
5. Update the deployment documentation
