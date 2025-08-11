# 🚀 Cloudflare Deployment Guide

This guide covers deploying the AI Agentic Crypto Browser to Cloudflare's edge infrastructure using Pages, Workers, D1, and KV.

## 📋 Prerequisites

1. **Cloudflare Account**: Sign up at [cloudflare.com](https://cloudflare.com)
2. **Domain**: Add your domain to Cloudflare (optional but recommended)
3. **Wrangler CLI**: Install globally with `npm install -g wrangler`
4. **Node.js**: Version 18 or higher
5. **Git**: For version control and CI/CD

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cloudflare    │    │   Cloudflare    │    │   Cloudflare    │
│     Pages       │    │    Workers      │    │       D1        │
│   (Frontend)    │────│   (Backend)     │────│   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   Cloudflare    │              │
         └──────────────│       KV        │──────────────┘
                        │   (Caching)     │
                        └─────────────────┘
```

## 🚀 Quick Start

### 1. Login to Cloudflare

```bash
wrangler login
```

### 2. Set up the database

```bash
cd cloudflare/database
chmod +x setup.sh
./setup.sh
```

### 3. Set up KV namespaces

```bash
cd cloudflare/kv
chmod +x setup.sh
./setup.sh
```

### 4. Configure Workers

Update `cloudflare/workers/api/wrangler.toml` with your:
- Database ID
- KV namespace IDs
- Environment variables

### 5. Deploy everything

```bash
cd cloudflare
chmod +x deploy.sh
./deploy.sh production
```

## 📁 Project Structure

```
cloudflare/
├── database/           # D1 database setup
│   ├── schema.sql     # Complete database schema
│   ├── migrations/    # Migration files
│   ├── seeds/         # Sample data
│   └── setup.sh       # Database setup script
├── workers/           # Cloudflare Workers
│   └── api/           # Main API worker
│       ├── src/       # Worker source code
│       ├── wrangler.toml
│       └── package.json
├── kv/                # KV namespace setup
│   └── setup.sh       # KV setup script
├── security/          # Security configuration
│   ├── waf-rules.json # WAF rules
│   ├── page-rules.json # Page rules
│   └── setup.sh       # Security setup
├── dns/               # DNS configuration
│   └── dns-records.json
├── deploy.sh          # Main deployment script
└── README.md          # This file
```

## 🔧 Configuration

### Environment Variables

Set these secrets in your GitHub repository or Cloudflare dashboard:

```bash
# Required
CLOUDFLARE_API_TOKEN=your-api-token
CLOUDFLARE_ACCOUNT_ID=your-account-id

# Optional
NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=your-project-id
OPENAI_API_KEY=your-openai-key
ANTHROPIC_API_KEY=your-anthropic-key
JWT_SECRET=your-jwt-secret
```

### Wrangler Configuration

Update `cloudflare/workers/api/wrangler.toml`:

```toml
name = "ai-crypto-browser-api"
main = "src/index.ts"
compatibility_date = "2024-01-01"

# Add your database ID
[[d1_databases]]
binding = "DB"
database_name = "ai-crypto-browser-db"
database_id = "your-database-id"

# Add your KV namespace IDs
[[kv_namespaces]]
binding = "CACHE"
id = "your-cache-namespace-id"
preview_id = "your-cache-preview-id"

[[kv_namespaces]]
binding = "SESSIONS"
id = "your-sessions-namespace-id"
preview_id = "your-sessions-preview-id"
```

## 🌍 Environments

### Production
- **Frontend**: `https://ai-agentic-crypto-browser.pages.dev`
- **API**: `https://api.your-domain.com`
- **Branch**: `main`

### Staging
- **Frontend**: `https://ai-agentic-crypto-browser-staging.pages.dev`
- **API**: `https://api-staging.your-domain.com`
- **Branch**: `staging`

### Development
- **Frontend**: `https://ai-agentic-crypto-browser-dev.pages.dev`
- **API**: `https://api-dev.your-domain.com`
- **Branch**: `dev`

## 🔄 CI/CD Pipeline

The GitHub Actions workflow automatically:

1. **Tests** code on pull requests
2. **Deploys** database migrations
3. **Deploys** Workers (API backend)
4. **Deploys** Pages (frontend)
5. **Verifies** deployment health

### Manual Deployment

```bash
# Deploy to production
./cloudflare/deploy.sh production

# Deploy to staging
./cloudflare/deploy.sh staging

# Deploy to development
./cloudflare/deploy.sh development
```

## 🔒 Security Features

- **WAF Protection**: SQL injection, XSS, and other attack prevention
- **DDoS Protection**: Automatic mitigation of DDoS attacks
- **Rate Limiting**: API endpoint protection
- **SSL/TLS**: Full strict encryption
- **Bot Management**: Automated bot detection and mitigation

## 📊 Monitoring

### Health Checks

```bash
# Frontend
curl https://your-domain.com

# API
curl https://api.your-domain.com/health

# Database
wrangler d1 execute ai-crypto-browser-db --command="SELECT COUNT(*) FROM users"
```

### Logs

```bash
# Worker logs
wrangler tail ai-crypto-browser-api

# Pages logs
wrangler pages deployment tail ai-agentic-crypto-browser
```

## 🛠️ Development

### Local Development

```bash
# Start frontend
cd web
npm run dev

# Start Worker locally
cd cloudflare/workers/api
wrangler dev
```

### Testing

```bash
# Frontend tests
cd web
npm test

# Worker tests
cd cloudflare/workers/api
npm test
```

## 📈 Performance Optimization

- **Edge Caching**: Static assets cached globally
- **KV Storage**: Fast key-value storage for sessions
- **D1 Database**: SQLite at the edge
- **Minification**: Automatic CSS/JS/HTML minification
- **Compression**: Brotli and Gzip compression

## 🔧 Troubleshooting

### Common Issues

1. **Build Failures**
   ```bash
   # Clear cache and rebuild
   cd web
   rm -rf .next out node_modules
   npm install
   npm run build:cloudflare
   ```

2. **Database Connection Issues**
   ```bash
   # Check database status
   wrangler d1 info ai-crypto-browser-db
   ```

3. **Worker Deployment Issues**
   ```bash
   # Check Worker status
   wrangler status
   ```

### Getting Help

- **Cloudflare Docs**: [developers.cloudflare.com](https://developers.cloudflare.com)
- **Discord**: [Cloudflare Developers Discord](https://discord.gg/cloudflaredev)
- **GitHub Issues**: Create an issue in this repository

## 📚 Additional Resources

- [Cloudflare Pages Documentation](https://developers.cloudflare.com/pages/)
- [Cloudflare Workers Documentation](https://developers.cloudflare.com/workers/)
- [Cloudflare D1 Documentation](https://developers.cloudflare.com/d1/)
- [Cloudflare KV Documentation](https://developers.cloudflare.com/workers/runtime-apis/kv/)

## 🎉 Success!

Once deployed, your AI Agentic Crypto Browser will be running on Cloudflare's global edge network with:

- ⚡ **Ultra-fast performance** with edge caching
- 🔒 **Enterprise-grade security** with WAF and DDoS protection
- 🌍 **Global availability** with 200+ data centers
- 📈 **Automatic scaling** based on demand
- 💰 **Cost-effective** with generous free tiers
