# AI Agentic Browser Makefile

.PHONY: help build test clean dev docker-up docker-down deps lint format

# Default environment for infrastructure
ENV ?= dev

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  deps         - Download Go dependencies"
	@echo "  build        - Build all services"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run linters"
	@echo "  format       - Format code"
	@echo "  dev          - Start development environment"
	@echo "  docker-up    - Start all services with Docker"
	@echo "  docker-down  - Stop all Docker services"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "Infrastructure:"
	@echo "  setup        - Set up development environment"
	@echo "  validate     - Validate infrastructure configuration"
	@echo "  deploy       - Deploy infrastructure and application"
	@echo "  destroy      - Destroy infrastructure"
	@echo "  status       - Show deployment status"
	@echo "  logs         - View application logs"
	@echo "  scale        - Scale services (SERVICE=name REPLICAS=n)"

# Go variables
GO_VERSION := 1.22
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
SERVICES := auth-service ai-agent browser-service web3-service api-gateway

# Development
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o bin/$$service ./cmd/$$service; \
	done

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

format:
	@echo "Formatting Go code..."
	gofmt -s -w $(GO_FILES)
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w $(GO_FILES); \
	else \
		echo "goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Development environment
dev-infra:
	@echo "Starting infrastructure services..."
	./scripts/docker-setup.sh

dev-stop-infra:
	@echo "Stopping infrastructure services..."
	./scripts/docker-cleanup.sh

dev: dev-infra
	@echo "Starting development environment..."
	@echo "Infrastructure services started. You can now run individual services:"
	@echo "  go run cmd/auth-service/main.go"
	@echo "  go run cmd/ai-agent/main.go"
	@echo "  go run cmd/browser-service/main.go"
	@echo "  go run cmd/web3-service/main.go"
	@echo "  go run cmd/api-gateway/main.go"

# Docker operations
docker-up:
	@echo "Starting all services with Docker..."
	./scripts/docker-setup.sh

docker-down:
	@echo "Stopping all Docker services..."
	./scripts/docker-cleanup.sh

docker-logs:
	@echo "Showing Docker logs..."
	docker logs -f postgres &
	docker logs -f redis &
	docker logs -f jaeger &
	docker logs -f prometheus &
	docker logs -f grafana &
	wait

docker-build:
	@echo "Building Docker images..."
	docker build -t ai-browser/auth-service -f cmd/auth-service/Dockerfile .
	docker build -t ai-browser/ai-agent -f cmd/ai-agent/Dockerfile .
	docker build -t ai-browser/browser-service -f cmd/browser-service/Dockerfile .
	docker build -t ai-browser/web3-service -f cmd/web3-service/Dockerfile .
	docker build -t ai-browser/api-gateway -f cmd/api-gateway/Dockerfile .
	docker build -t ai-browser/web -f web/Dockerfile ./web

# Database operations
db-migrate:
	@echo "Running database migrations..."
	@echo "Migrations are automatically applied via Docker init script"

db-reset:
	@echo "Resetting database..."
	docker stop postgres || true
	docker rm postgres || true
	docker volume rm postgres_data || true
	docker run -d --name postgres \
		--network ai-browser-network \
		-e POSTGRES_DB=ai_agentic_browser \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-p 5432:5432 \
		-v postgres_data:/var/lib/postgresql/data \
		postgres:16

# Frontend operations
frontend-install:
	@echo "Installing frontend dependencies..."
	cd web && npm install

frontend-dev:
	@echo "Starting frontend development server..."
	cd web && npm run dev

frontend-build:
	@echo "Building frontend for production..."
	cd web && npm run build

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	docker system prune -f

# Health checks
health:
	@echo "Checking service health..."
	@echo "Auth Service:"
	@curl -s http://localhost:8081/health | jq . || echo "  ❌ Not responding"
	@echo "AI Agent Service:"
	@curl -s http://localhost:8082/health | jq . || echo "  ❌ Not responding"
	@echo "Browser Service:"
	@curl -s http://localhost:8083/health | jq . || echo "  ❌ Not responding"
	@echo "Web3 Service:"
	@curl -s http://localhost:8084/health | jq . || echo "  ❌ Not responding"
	@echo "API Gateway:"
	@curl -s http://localhost:8080/health | jq . || echo "  ❌ Not responding"

# Test individual services
test-auth:
	@echo "Testing auth service..."
	go run cmd/auth-service/main.go &
	@sleep 2
	@curl -X POST http://localhost:8081/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq .
	@pkill -f "auth-service"

test-ai:
	@echo "Testing AI agent service..."
	@echo "Note: Requires valid OpenAI API key in .env"
	go run cmd/ai-agent/main.go &
	@sleep 2
	@echo "AI Agent service started on port 8082"
	@pkill -f "ai-agent"

test-browser:
	@echo "Testing browser service..."
	go run cmd/browser-service/main.go &
	@sleep 2
	@echo "Browser service started on port 8083"
	@pkill -f "browser-service"

# Run individual services for development
run-auth:
	@echo "Starting auth service..."
	go run cmd/auth-service/main.go

run-ai:
	@echo "Starting AI agent service..."
	go run cmd/ai-agent/main.go

run-browser:
	@echo "Starting browser service..."
	go run cmd/browser-service/main.go

run-web3:
	@echo "Starting web3 service..."
	go run cmd/web3-service/main.go

run-gateway:
	@echo "Starting API gateway..."
	go run cmd/api-gateway/main.go

# Run the new HFT trading system
run-hft:
	@echo "Starting HFT Trading System..."
	go run cmd/main.go

# Build the HFT trading system
build-hft:
	@echo "Building HFT Trading System..."
	go build -o bin/ai-agentic-browser cmd/main.go

# Security
security-scan:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Tools installation
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Environment setup
setup: install-tools deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please update .env with your API keys and configuration"; \
	fi
	@echo "Setup complete! Run 'make dev' to start development environment"

# Production
prod-build:
	@echo "Building for production..."
	docker build -t ai-browser/auth-service:prod -f cmd/auth-service/Dockerfile .
	docker build -t ai-browser/ai-agent:prod -f cmd/ai-agent/Dockerfile .
	docker build -t ai-browser/browser-service:prod -f cmd/browser-service/Dockerfile .
	docker build -t ai-browser/web3-service:prod -f cmd/web3-service/Dockerfile .
	docker build -t ai-browser/api-gateway:prod -f cmd/api-gateway/Dockerfile .
	docker build -t ai-browser/web:prod -f web/Dockerfile ./web

prod-deploy:
	@echo "Deploying to production..."
	@echo "This would deploy to your production environment"
	@echo "Implement your deployment strategy here"

# Infrastructure Management
setup: ## Set up development environment
	@echo "Setting up development environment..."
	@chmod +x scripts/*.sh
	@./scripts/setup.sh

validate: ## Validate infrastructure configuration
	@echo "Validating configuration..."
	@./scripts/validate.sh

deploy: ## Deploy infrastructure and application
	@echo "Deploying to $(ENV) environment..."
	@./scripts/deploy.sh $(ENV)

destroy: ## Destroy infrastructure
	@echo "Destroying $(ENV) environment..."
	@read -p "Are you sure you want to destroy the $(ENV) environment? [y/N] " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		kubectl delete -k k8s/overlays/$(ENV) || true; \
		cd terraform/environments/$(ENV) && terraform destroy; \
	else \
		echo "Cancelled."; \
	fi

# Kubernetes operations
status: ## Show deployment status
	@echo "Deployment Status:"
	@kubectl get all -n ai-crypto-browser

logs: ## View application logs
	@echo "Application Logs:"
	@kubectl logs -f -l app.kubernetes.io/part-of=ai-crypto-browser -n ai-crypto-browser --tail=100

scale: ## Scale services (usage: make scale SERVICE=api-gateway REPLICAS=3)
	@if [ -z "$(SERVICE)" ] || [ -z "$(REPLICAS)" ]; then \
		echo "Usage: make scale SERVICE=<service-name> REPLICAS=<number>"; \
		echo "Available services: api-gateway, auth-service, browser-service, web3-service, frontend"; \
	else \
		echo "Scaling $(SERVICE) to $(REPLICAS) replicas..."; \
		kubectl scale deployment/$(SERVICE) --replicas=$(REPLICAS) -n ai-crypto-browser; \
	fi

restart: ## Restart a service (usage: make restart SERVICE=api-gateway)
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make restart SERVICE=<service-name>"; \
		echo "Available services: api-gateway, auth-service, browser-service, web3-service, frontend"; \
	else \
		echo "Restarting $(SERVICE)..."; \
		kubectl rollout restart deployment/$(SERVICE) -n ai-crypto-browser; \
	fi

# Quick commands
quick-deploy: validate deploy ## Quick deploy after validation

quick-status: ## Quick status check
	@kubectl get pods -n ai-crypto-browser
	@echo ""
	@kubectl get services -n ai-crypto-browser
