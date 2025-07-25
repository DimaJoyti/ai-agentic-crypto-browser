# AI-Powered Agentic Browser

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
- **Containerization**: Docker and Docker Compose
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
   docker-compose up -d postgres redis jaeger prometheus grafana
   ```

2. **Initialize Go modules**
   ```bash
   go mod tidy
   ```

3. **Run database migrations**
   ```bash
   # Database will be initialized automatically via Docker
   # Check logs: docker-compose logs postgres
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

### Using Docker Compose

For a complete setup with all services:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
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
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags=integration ./test/...
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
   docker-compose -f docker-compose.prod.yml build
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
