# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI-Powered Agentic Browser - An intelligent web browser that uses AI agents to autonomously navigate, interact with, and extract information from websites, with integrated cryptocurrency and Web3 functionality.

## Development Commands

### Core Development
- `make dev` - Start infrastructure services (postgres, redis, jaeger, prometheus, grafana) and print service startup commands
- `make dev-infra` - Start only infrastructure services
- `make build` - Build all Go services to bin/ directory
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report (creates coverage.html)
- `make lint` - Run golangci-lint (install with `make install-tools`)
- `make format` - Format Go code with gofmt and goimports
- `make deps` - Download and tidy Go dependencies (`go mod download && go mod tidy`)
- `make setup` - Full environment setup (install tools, deps, create .env)

**Note**: This project uses a Go workspace (`go.work`) with a single module containing all services.

### Individual Services (Default Ports)
- `go run cmd/auth-service/main.go` - Start auth service (port 8081)
- `go run cmd/ai-agent/main.go` - Start AI agent service (port 8082)  
- `go run cmd/browser-service/main.go` - Start browser service (port 8083)
- `go run cmd/web3-service/main.go` - Start Web3 service (port 8084)
- `go run cmd/api-gateway/main.go` - Start API gateway (port 8080)

### Make Shortcuts for Services
- `make run-auth` - Start auth service
- `make run-ai` - Start AI agent service
- `make run-browser` - Start browser service
- `make run-web3` - Start Web3 service
- `make run-gateway` - Start API gateway

### Frontend Development
- `cd web && npm install` - Install frontend dependencies
- `cd web && npm run dev` - Start Next.js development server (port 3000)
- `cd web && npm run build` - Build frontend for production
- `cd web && npm run lint` - Run ESLint
- `cd web && npm run type-check` - Run TypeScript type checking

### Docker Operations
- `make docker-up` - Start all services with Docker Compose
- `make docker-down` - Stop all Docker services
- `make docker-logs` - View all service logs
- `make docker-build` - Build Docker images
- `make health` - Check health of all services (requires curl and jq)
- `docker-compose up -d` - Direct docker-compose command

## Architecture

### Backend Services (Go)
**Microservices architecture with the following services:**

- **API Gateway** (`cmd/api-gateway/`) - Main entry point, request routing, middleware
- **Auth Service** (`cmd/auth-service/`) - JWT authentication, user management
- **AI Agent Service** (`cmd/ai-agent/`) - OpenAI/Anthropic integration, natural language processing
- **Browser Service** (`cmd/browser-service/`) - Headless Chrome automation via chromedp
- **Web3 Service** (`cmd/web3-service/`) - Blockchain interactions, wallet connectivity

**Core packages:**
- `internal/config/` - Environment-based configuration management
- `pkg/database/` - PostgreSQL and Redis database utilities  
- `pkg/middleware/` - HTTP middleware (JWT auth, rate limiting, CORS, logging, tracing)
- `pkg/observability/` - OpenTelemetry tracing and structured logging

### Frontend (Next.js/React)
- **Stack**: Next.js 14, TypeScript, TailwindCSS, Radix UI components
- **State Management**: Zustand for global state, React Query for server state
- **Web3**: Wagmi + Viem for blockchain interactions
- **Location**: `web/` directory with standard Next.js app router structure

### Database Schema
- **PostgreSQL**: User accounts, sessions, tasks, browser sessions, wallets
- **Redis**: Caching, session storage, rate limiting
- **Migrations**: Auto-applied via `scripts/init.sql` on container startup

## Key Implementation Details

### Authentication Flow
JWT-based authentication with refresh tokens. Protected routes use JWT middleware to extract user context. All services except auth-service require valid Authorization header. Sessions stored in Redis with configurable expiry.

### AI Agent Integration  
Supports OpenAI and Anthropic providers via AI_MODEL_PROVIDER env var. The AI service integrates with browser service for task execution. Task types include navigate, extract, interact. Some endpoints are implemented as TODO placeholders.

### Browser Automation
Uses chromedp for headless Chrome automation. Session-based architecture - users create browser sessions via API. Supports element interaction, content extraction, screenshots. Configured for Docker with disabled GPU and sandbox.

### Web3 Integration
Multi-chain RPC configuration (Ethereum, Polygon, Arbitrum, Optimism). Frontend uses Wagmi/Viem for wallet connections. Backend Web3 service handles blockchain interactions.

### Observability
OpenTelemetry integration with Jaeger tracing and structured logging. Health endpoints on all services. Prometheus metrics and Grafana dashboards configured via docker-compose.

## Environment Setup

Essential environment variables (copy from `.env.example`):
```bash
# Required for all services
DATABASE_URL=postgres://postgres:postgres@localhost:5432/agentic_browser?sslmode=disable
JWT_SECRET=your-super-secure-jwt-secret  # Change in production!

# Required for AI features
OPENAI_API_KEY=sk-your-openai-key
# OR
ANTHROPIC_API_KEY=your-anthropic-key

# Optional Web3 features
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/your-project-id
```

**Setup**: Run `make setup` to install tools, dependencies, and create .env from example.

## Testing

- **Unit Tests**: `go test ./...` for all Go packages
- **Test with Coverage**: `make test-coverage` (creates coverage.html report)  
- **Integration Tests**: Not yet implemented (no docker-compose.test.yml found)
- **Individual Service Tests**: `make test-auth`, `make test-ai`, `make test-browser` (automated test runners)
- **Setup Validation**: `./scripts/test-setup.sh` - Validates environment and dependencies

## Service Dependencies

Services have these startup dependencies:
- All services depend on PostgreSQL and Redis
- API Gateway coordinates with all other services
- AI Agent requires valid OPENAI_API_KEY or ANTHROPIC_API_KEY in environment
- Web3 Service requires RPC URLs for blockchain access
- Browser Service needs Chrome/Chromium runtime in container

## Common Development Patterns

### Adding New AI Task Types
1. Define task type constant in `internal/ai/models.go`
2. Implement execution logic in `internal/ai/service.go`
3. Add case to task execution switch statement

### Adding New Browser Actions
1. Define action type in `internal/browser/models.go`
2. Implement action handler in `internal/browser/service.go`
3. Add case to action execution switch

### Adding New API Endpoints
1. Define in appropriate service main.go
2. Add middleware for auth/validation as needed
3. Follow existing error handling patterns

## Code Conventions

- **Go**: Standard Go conventions, use gofmt/goimports
- **Error Handling**: Always return detailed error context
- **Logging**: Use structured logging with context
- **HTTP**: JSON APIs with proper status codes
- **Database**: Use transactions for multi-table operations
- **Security**: Validate all inputs, use prepared statements