# AI-Powered Agentic Crypto Browser

An intelligent web browser that uses AI agents to autonomously navigate, interact with, and extract information from websites, with integrated cryptocurrency and Web3 functionality.

## 🚀 Features

### Core Capabilities
- **AI-Powered Browsing**: Autonomous web navigation using natural language commands
- **Intelligent Content Analysis**: Automatic page summarization and data extraction
- **Web3 Integration**: Cryptocurrency wallet connection and DeFi interactions
- **Microservices Architecture**: Scalable, maintainable backend services
- **Real-time Monitoring**: Comprehensive observability with OpenTelemetry

### AI Agent Features
- Natural language command processing
- Context-aware navigation and interaction
- Automated form filling and data extraction
- Learning from user behavior patterns
- Intelligent content summarization

### Cryptocurrency Integration
- Multi-wallet support (MetaMask, WalletConnect, Coinbase Wallet)
- Transaction monitoring and management
- DeFi protocol interactions (Uniswap, Aave, Compound)
- NFT marketplace browsing and transactions
- Real-time cryptocurrency price tracking

## 🏗️ Architecture

### Backend Services (Go)
- **API Gateway**: Main entry point and request routing
- **Auth Service**: User authentication and session management
- **AI Agent Service**: Natural language processing and agent logic
- **Browser Service**: Headless browser automation and page interaction
- **Web3 Service**: Cryptocurrency and blockchain interactions

### Frontend (React/Next.js)
- Modern React application with TypeScript
- TailwindCSS for styling with Shadcn/ui components
- Web3 integration with Wagmi and Viem
- Real-time communication with WebSocket support

### Infrastructure
- **Database**: PostgreSQL for structured data, Redis for caching
- **Observability**: OpenTelemetry, Prometheus, Grafana, Jaeger
- **Containerization**: Docker containers with custom networking
- **Security**: JWT authentication, rate limiting, input validation

## 🛠️ Quick Start

### Prerequisites
- Go 1.22 or later
- Node.js 18 or later
- Docker
- PostgreSQL 16+
- Redis 7+

### Environment Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ai-agentic-browser
   ```

2. **Copy environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Update environment variables**
   Edit `.env` file with your API keys and configuration:
   ```bash
   # Required API Keys
   OPENAI_API_KEY=your-openai-api-key
   ANTHROPIC_API_KEY=your-anthropic-api-key
   ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/your-project-id
   
   # JWT Secret (generate a secure random string)
   JWT_SECRET=your-super-secret-jwt-key-change-in-production
   ```

### Development Setup

1. **Start infrastructure services**
   ```bash
   # Create a Docker network
   docker network create ai-browser-network

   # Start PostgreSQL
   docker run -d --name postgres \
     --network ai-browser-network \
     -e POSTGRES_DB=ai_agentic_browser \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -p 5432:5432 \
     postgres:16

   # Start Redis
   docker run -d --name redis \
     --network ai-browser-network \
     -p 6379:6379 \
     redis:7-alpine

   # Start Jaeger (optional - for tracing)
   docker run -d --name jaeger \
     --network ai-browser-network \
     -p 16686:16686 \
     -p 14268:14268 \
     jaegertracing/all-in-one:latest

   # Start Prometheus (optional - for metrics)
   docker run -d --name prometheus \
     --network ai-browser-network \
     -p 9090:9090 \
     prom/prometheus:latest

   # Start Grafana (optional - for dashboards)
   docker run -d --name grafana \
     --network ai-browser-network \
     -p 3001:3000 \
     -e GF_SECURITY_ADMIN_PASSWORD=admin \
     grafana/grafana:latest
   ```

2. **Initialize Go modules**
   ```bash
   go mod tidy
   ```

3. **Run database migrations**
   ```bash
   # Database will be initialized automatically via Docker
   # Check logs: docker logs postgres
   ```

4. **Start backend services**
   ```bash
   # Terminal 1: Auth Service
   go run cmd/auth-service/main.go
   
   # Terminal 2: AI Agent Service
   go run cmd/ai-agent/main.go
   
   # Terminal 3: Browser Service
   go run cmd/browser-service/main.go
   
   # Terminal 4: Web3 Service
   go run cmd/web3-service/main.go
   
   # Terminal 5: API Gateway
   go run cmd/api-gateway/main.go
   ```

5. **Start frontend**
   ```bash
   cd web
   npm install
   npm run dev
   ```

### Using Docker for Complete Setup

For a complete setup with all services using Docker:

```bash
# Build and run all services
./scripts/docker-setup.sh

# View logs for specific service
docker logs -f postgres
docker logs -f redis
docker logs -f api-gateway

# Stop all services
docker stop postgres redis jaeger prometheus grafana
docker rm postgres redis jaeger prometheus grafana
docker network rm ai-browser-network
```

## 📊 Monitoring and Observability

### Access Points
- **Application**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686

### Health Checks
```bash
# Check service health
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # AI Agent Service
curl http://localhost:8083/health  # Browser Service
curl http://localhost:8084/health  # Web3 Service
```

## 🧪 Testing

### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
go test ./internal/auth/...
```

### Integration Tests
```bash
# Start test environment (PostgreSQL and Redis)
docker run -d --name test-postgres \
  -e POSTGRES_DB=ai_agentic_browser_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:16

docker run -d --name test-redis \
  -p 6380:6379 \
  redis:7-alpine

# Run integration tests
go test -tags=integration ./test/...

# Cleanup test environment
docker stop test-postgres test-redis
docker rm test-postgres test-redis
```

## 🔧 Development

### Project Structure
```
ai-agentic-browser/
├── cmd/                    # Application entrypoints
├── internal/               # Private application code
├── pkg/                    # Public packages
├── web/                    # Frontend application
├── configs/                # Configuration files
├── scripts/                # Build and utility scripts
├── test/                   # Test utilities and fixtures
└── docs/                   # Documentation
```

### Adding New Features

1. **Backend Services**: Add new services in `cmd/` and `internal/`
2. **Frontend Components**: Add React components in `web/components/`
3. **Database Changes**: Update schema in `scripts/init.sql`
4. **API Endpoints**: Define in `api/openapi/` and implement in services

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt` and `golangci-lint`
- **TypeScript/React**: Use ESLint and Prettier for formatting
- **Commits**: Use conventional commit messages

## 🚀 Deployment

### Production Deployment

1. **Build Docker images**
   ```bash
   # Build individual service images
   docker build -t ai-browser/auth-service -f cmd/auth-service/Dockerfile .
   docker build -t ai-browser/ai-agent -f cmd/ai-agent/Dockerfile .
   docker build -t ai-browser/browser-service -f cmd/browser-service/Dockerfile .
   docker build -t ai-browser/web3-service -f cmd/web3-service/Dockerfile .
   docker build -t ai-browser/api-gateway -f cmd/api-gateway/Dockerfile .
   docker build -t ai-browser/web -f web/Dockerfile ./web
   ```

2. **Deploy to Kubernetes**
   ```bash
   kubectl apply -f deployments/k8s/
   ```

3. **Configure environment variables**
   - Update production secrets
   - Configure external databases
   - Set up monitoring and alerting

### Security Considerations

- Use strong JWT secrets in production
- Enable HTTPS/TLS for all services
- Configure proper CORS origins
- Set up rate limiting and DDoS protection
- Regular security audits and dependency updates

## 📚 API Documentation

### Authentication Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh
- `POST /auth/logout` - User logout
- `GET /auth/me` - Get user profile

### AI Agent Endpoints
- `POST /ai/chat` - Send message to AI agent
- `GET /ai/conversations` - List user conversations
- `GET /ai/conversations/{id}` - Get conversation details
- `POST /ai/tasks` - Create new AI task
- `GET /ai/tasks/{id}` - Get task status

### Browser Endpoints
- `POST /browser/sessions` - Create browser session
- `GET /browser/sessions` - List user sessions
- `POST /browser/navigate` - Navigate to URL (requires X-Session-ID header)
- `POST /browser/interact` - Interact with page elements
- `POST /browser/extract` - Extract page content
- `POST /browser/screenshot` - Take page screenshot

### Web3 Endpoints
- `POST /web3/connect-wallet` - Connect cryptocurrency wallet
- `GET /web3/balance` - Get wallet balance
- `POST /web3/transaction` - Send transaction
- `GET /web3/defi/positions` - Get DeFi positions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the `docs/` directory
- **Issues**: Report bugs and feature requests on GitHub
- **Discussions**: Join our community discussions

## 🗺️ Roadmap

### Phase 1: Foundation ✅
- [x] Project setup and architecture
- [x] Authentication service
- [x] Database schema and migrations
- [x] Basic observability setup

### Phase 2: AI Agent Core ✅
- [x] AI service integration with OpenAI
- [x] Natural language processing and chat interface
- [x] Browser automation service with Chromedp
- [x] Basic agent workflows and task execution

### Phase 3: Web3 Integration
- [ ] Wallet connection infrastructure
- [ ] Transaction monitoring
- [ ] DeFi protocol integrations
- [ ] NFT marketplace support

### Phase 4: Advanced Features
- [ ] Advanced AI capabilities
- [ ] Learning and adaptation
- [ ] Performance optimization
- [ ] Security hardening

### Phase 5: Production Ready
- [ ] Comprehensive testing
- [ ] Production deployment
- [ ] Monitoring and alerting
- [ ] Documentation and tutorials
