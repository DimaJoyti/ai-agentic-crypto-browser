#!/bin/bash

# AI Agentic Browser Environment Setup Script

echo "🔧 Setting up environment variables for AI Agentic Browser..."

# Database Configuration
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/agentic_browser?sslmode=disable"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="agentic_browser"
export DB_USER="postgres"
export DB_PASSWORD="postgres"

# Redis Configuration
export REDIS_URL="redis://localhost:6379"
export REDIS_HOST="localhost"
export REDIS_PORT="6379"

# JWT and Security
export JWT_SECRET="your-super-secret-jwt-key-change-in-production"
export API_SECRET="your-api-secret-key"
export ENCRYPTION_KEY="your-32-char-encryption-key-here"

# Server Configuration
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"
export API_PORT="8080"

# AI Configuration
export OPENAI_API_KEY="your-openai-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"

# Web3 Configuration
export WEB3_RPC_URL="https://mainnet.infura.io/v3/your-project-id"
export ETHEREUM_PRIVATE_KEY="your-ethereum-private-key"

# External APIs
export COINMARKETCAP_API_KEY="your-coinmarketcap-api-key"
export COINGECKO_API_KEY="your-coingecko-api-key"

# Monitoring
export JAEGER_ENDPOINT="http://localhost:14268/api/traces"
export PROMETHEUS_ENDPOINT="http://localhost:9090"

# Environment
export ENVIRONMENT="development"
export LOG_LEVEL="info"

echo "✅ Environment variables set!"
echo ""
echo "📋 Configuration Summary:"
echo "  Database: postgres://postgres:***@localhost:5432/agentic_browser"
echo "  Redis: redis://localhost:6379"
echo "  API Server: http://localhost:8080"
echo "  Jaeger UI: http://localhost:16686"
echo "  Prometheus: http://localhost:9090"
echo "  Grafana: http://localhost:3001"
echo ""
echo "🚀 You can now start the services:"
echo "  go run cmd/api-gateway/main.go"
echo "  go run cmd/ai-agent/main.go"
echo "  go run cmd/auth-service/main.go"
echo ""
echo "💡 To use these variables in your current shell:"
echo "  source scripts/setup-env.sh"
