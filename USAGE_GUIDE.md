# AI Agentic Browser - Complete Usage Guide

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 18+ (for local development)
- Go 1.22+ (for local development)

### 1. Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd ai-agentic-browser

# Copy environment configuration
cp .env.example .env

# Edit .env with your API keys
nano .env
```

### 2. Required Environment Variables

```bash
# Essential Configuration
JWT_SECRET=your-super-secure-jwt-secret-here
OPENAI_API_KEY=sk-your-openai-api-key-here

# Optional Web3 Configuration
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/your-project-id
POLYGON_RPC_URL=https://polygon-mainnet.infura.io/v3/your-project-id
ARBITRUM_RPC_URL=https://arbitrum-mainnet.infura.io/v3/your-project-id
OPTIMISM_RPC_URL=https://optimism-mainnet.infura.io/v3/your-project-id

# Optional Frontend Configuration
WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id
```

### 3. Start the Application

```bash
# Option 1: Full deployment with validation
./scripts/validate-deployment.sh

# Option 2: Quick start
docker-compose up -d

# Option 3: Development mode
make dev-infra  # Start infrastructure only
# Then run services individually:
make run-auth
make run-ai
make run-browser
make run-web3
make run-gateway
cd web && npm run dev
```

## 🎯 **Core Features & Usage**

### **1. AI-Powered Web Automation**

#### Chat with AI Agent
```bash
# Via API
curl -X POST http://localhost:8082/ai/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "message": "Navigate to google.com and search for AI news"
  }'
```

#### Create Automation Tasks
```bash
# Navigate to a website
curl -X POST http://localhost:8082/ai/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "task_type": "navigate",
    "description": "Navigate to example.com",
    "input_data": {
      "url": "https://example.com",
      "wait_for_selector": "body"
    }
  }'

# Extract content from a page
curl -X POST http://localhost:8082/ai/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "task_type": "extract",
    "description": "Extract all links from the page",
    "input_data": {
      "url": "https://example.com",
      "data_type": "links"
    }
  }'
```

### **2. Browser Automation**

#### Create Browser Session
```bash
curl -X POST http://localhost:8083/browser/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "session_name": "My Automation Session"
  }'
```

#### Navigate and Interact
```bash
# Navigate to URL
curl -X POST http://localhost:8083/browser/navigate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Session-ID: $SESSION_ID" \
  -d '{
    "url": "https://example.com",
    "wait_for_selector": "body",
    "timeout": 10
  }'

# Interact with page elements
curl -X POST http://localhost:8083/browser/interact \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Session-ID: $SESSION_ID" \
  -d '{
    "actions": [
      {
        "type": "click",
        "selector": "#search-button"
      },
      {
        "type": "type",
        "selector": "#search-input",
        "value": "AI automation"
      }
    ],
    "screenshot": true
  }'
```

#### Extract Content
```bash
curl -X POST http://localhost:8083/browser/extract \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Session-ID: $SESSION_ID" \
  -d '{
    "data_type": "text",
    "selectors": ["h1", "p", ".content"]
  }'
```

### **3. Web3 & Cryptocurrency**

#### Connect Wallet
```bash
curl -X POST http://localhost:8084/web3/connect-wallet \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
    "chain_id": 1,
    "wallet_type": "metamask"
  }'
```

#### Get Wallet Balance
```bash
curl -X GET http://localhost:8084/web3/balance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "wallet_id": "$WALLET_ID"
  }'
```

#### Get Cryptocurrency Prices
```bash
curl -X GET http://localhost:8084/web3/prices \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

#### DeFi Protocol Interaction
```bash
curl -X POST http://localhost:8084/web3/defi/interact \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "wallet_id": "$WALLET_ID",
    "protocol": "uniswap",
    "action": "add_liquidity",
    "token_address": "0xA0b86a33E6441b8C4505E2E8E3C3C5C8E6441b8C",
    "amount": "1000000000000000000"
  }'
```

## 🖥️ **Frontend Usage**

### **1. Authentication**
1. Visit http://localhost:3000
2. Click "Get Started" or "Login"
3. Register a new account or login with existing credentials
4. Access the dashboard

### **2. AI Agent Interface**
1. Navigate to the AI Chat section
2. Type natural language commands:
   - "Navigate to google.com and take a screenshot"
   - "Extract all links from this page"
   - "Fill out the contact form with my details"
   - "Search for AI news and summarize the results"

### **3. Browser Automation**
1. Create a new browser session
2. Use the visual interface to:
   - Navigate to websites
   - Interact with page elements
   - Extract content
   - Take screenshots
   - Monitor automation progress

### **4. Web3 Dashboard**
1. Connect your cryptocurrency wallet (MetaMask, WalletConnect, etc.)
2. View wallet balances across multiple chains
3. Monitor transaction history
4. Interact with DeFi protocols
5. Track cryptocurrency prices

## 🔧 **Development & Customization**

### **Adding New AI Tasks**

1. Define task type in `internal/ai/models.go`:
```go
const (
    TaskTypeCustom TaskType = "custom_task"
)
```

2. Implement task execution in `internal/ai/service.go`:
```go
func (s *Service) executeCustomTask(ctx context.Context, task *Task) (map[string]interface{}, error) {
    // Your custom task logic here
    return outputData, nil
}
```

3. Add task handler in the switch statement

### **Adding New Browser Actions**

1. Define action type in `internal/browser/models.go`:
```go
const (
    ActionTypeCustom ActionType = "custom_action"
)
```

2. Implement action in `internal/browser/service.go`:
```go
case ActionTypeCustom:
    return s.executeCustomAction(ctx, action)
```

### **Adding New Web3 Protocols**

1. Add protocol support in `internal/web3/service.go`:
```go
case "new_protocol":
    txHash, position, err = s.simulateNewProtocolInteraction(ctx, wallet, req)
```

2. Implement protocol-specific logic

## 📊 **Monitoring & Debugging**

### **Service Health Monitoring**
```bash
# Check all services
curl http://localhost:8080/api/status | jq

# Individual service health
curl http://localhost:8081/health  # Auth
curl http://localhost:8082/health  # AI Agent
curl http://localhost:8083/health  # Browser
curl http://localhost:8084/health  # Web3
```

### **Logs and Debugging**
```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f auth-service
docker-compose logs -f ai-agent
docker-compose logs -f browser-service

# View database logs
docker-compose logs -f postgres

# Monitor resource usage
docker stats
```

### **Observability Stack**
- **Grafana**: http://localhost:3001 (admin/admin)
  - Pre-configured dashboards for all services
  - Real-time metrics and alerts
  
- **Prometheus**: http://localhost:9090
  - Metrics collection and querying
  - Service discovery and monitoring
  
- **Jaeger**: http://localhost:16686
  - Distributed tracing
  - Request flow visualization

## 🚨 **Troubleshooting**

### **Common Issues**

1. **Services not starting**
   ```bash
   # Check Docker resources
   docker system df
   docker system prune
   
   # Restart services
   docker-compose down
   docker-compose up -d
   ```

2. **Database connection issues**
   ```bash
   # Check PostgreSQL
   docker-compose exec postgres pg_isready -U postgres
   
   # Reset database
   docker-compose down postgres
   docker volume rm ai-agentic-browser_postgres_data
   docker-compose up -d postgres
   ```

3. **AI Agent not responding**
   - Verify OPENAI_API_KEY in .env
   - Check API quota and billing
   - Review ai-agent service logs

4. **Browser automation failing**
   - Ensure Chrome/Chromium dependencies
   - Check browser-service logs
   - Verify target website accessibility

5. **Web3 features not working**
   - Configure RPC URLs in .env
   - Check network connectivity
   - Verify wallet connection

### **Performance Optimization**

1. **Resource Allocation**
   ```bash
   # Increase Docker memory limit
   # Adjust in Docker Desktop settings
   
   # Scale services
   docker-compose up -d --scale ai-agent=2
   ```

2. **Database Optimization**
   ```sql
   -- Connect to PostgreSQL
   docker-compose exec postgres psql -U postgres -d agentic_browser
   
   -- Check performance
   SELECT * FROM pg_stat_activity;
   
   -- Optimize queries
   EXPLAIN ANALYZE SELECT * FROM users;
   ```

## 🔐 **Security Considerations**

### **Production Deployment**
1. Change default passwords and secrets
2. Enable HTTPS/TLS for all services
3. Configure proper CORS origins
4. Set up rate limiting and DDoS protection
5. Regular security audits and updates
6. Implement proper backup strategies

### **API Security**
- All endpoints require JWT authentication
- Rate limiting prevents abuse
- Input validation on all requests
- Secure headers and CORS configuration

## 📈 **Scaling & Production**

### **Horizontal Scaling**
```bash
# Scale specific services
docker-compose up -d --scale ai-agent=3 --scale browser-service=2

# Load balancer configuration
# Add nginx or traefik for load balancing
```

### **Kubernetes Deployment**
```bash
# Deploy to Kubernetes
kubectl apply -f deployments/k8s/

# Monitor deployment
kubectl get pods -n agentic-browser
kubectl logs -f deployment/ai-agent -n agentic-browser
```

This comprehensive guide covers all aspects of using, developing, and deploying the AI Agentic Browser. For additional support, refer to the API documentation at http://localhost:8080/api/docs when the system is running.
