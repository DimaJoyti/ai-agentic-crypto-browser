#!/bin/bash

# AI Agentic Browser Setup Test Script

set -e

echo "🚀 Testing AI Agentic Browser Setup"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Function to print info
print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if required tools are installed
echo "Checking required tools..."

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status 0 "Go is installed (version: $GO_VERSION)"
else
    print_status 1 "Go is not installed"
    echo "Please install Go 1.22 or later from https://golang.org/dl/"
    exit 1
fi

# Check Docker
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    print_status 0 "Docker is installed (version: $DOCKER_VERSION)"
else
    print_status 1 "Docker is not installed"
    echo "Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker Compose
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
    print_status 0 "Docker Compose is installed (version: $COMPOSE_VERSION)"
else
    print_status 1 "Docker Compose is not installed"
    echo "Please install Docker Compose"
    exit 1
fi

# Check Node.js (for frontend)
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    print_status 0 "Node.js is installed (version: $NODE_VERSION)"
else
    print_status 1 "Node.js is not installed"
    echo "Please install Node.js 18 or later from https://nodejs.org/"
fi

echo ""
echo "Checking project structure..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_status 1 "go.mod not found. Are you in the project root directory?"
    exit 1
fi

# Check key files
FILES_TO_CHECK=(
    "go.mod"
    "docker-compose.yml"
    ".env.example"
    "cmd/auth-service/main.go"
    "internal/config/config.go"
    "pkg/database/postgres.go"
    "scripts/init.sql"
)

for file in "${FILES_TO_CHECK[@]}"; do
    if [ -f "$file" ]; then
        print_status 0 "$file exists"
    else
        print_status 1 "$file is missing"
    fi
done

echo ""
echo "Testing Go module setup..."

# Test Go module
if go mod verify &> /dev/null; then
    print_status 0 "Go modules are valid"
else
    print_info "Running go mod tidy..."
    if go mod tidy; then
        print_status 0 "Go modules tidied successfully"
    else
        print_status 1 "Failed to tidy Go modules"
    fi
fi

echo ""
echo "Testing Docker setup..."

# Test Docker
if docker info &> /dev/null; then
    print_status 0 "Docker daemon is running"
else
    print_status 1 "Docker daemon is not running"
    echo "Please start Docker daemon"
    exit 1
fi

echo ""
echo "Checking environment configuration..."

# Check .env file
if [ -f ".env" ]; then
    print_status 0 ".env file exists"
    
    # Check for required environment variables
    if grep -q "JWT_SECRET=" .env && [ "$(grep JWT_SECRET= .env | cut -d'=' -f2)" != "your-super-secret-jwt-key-change-in-production" ]; then
        print_status 0 "JWT_SECRET is configured"
    else
        print_status 1 "JWT_SECRET needs to be configured in .env"
    fi
    
    if grep -q "DATABASE_URL=" .env; then
        print_status 0 "DATABASE_URL is configured"
    else
        print_status 1 "DATABASE_URL needs to be configured in .env"
    fi
else
    print_info "Creating .env file from .env.example..."
    if cp .env.example .env; then
        print_status 0 ".env file created"
        echo "Please edit .env file with your configuration"
    else
        print_status 1 "Failed to create .env file"
    fi
fi

echo ""
echo "Testing basic compilation..."

# Test if auth service compiles
if go build -o /tmp/auth-service ./cmd/auth-service &> /dev/null; then
    print_status 0 "Auth service compiles successfully"
    rm -f /tmp/auth-service
else
    print_status 1 "Auth service compilation failed"
    echo "Run 'go build ./cmd/auth-service' for details"
fi

echo ""
echo "Setup test complete!"
echo ""

# Provide next steps
echo "🎯 Next Steps:"
echo "1. Update .env file with your API keys and configuration"
echo "2. Start infrastructure: make dev-infra"
echo "3. Run services individually or use: make docker-up"
echo "4. Access the application at http://localhost:3000"
echo ""
echo "📚 Useful commands:"
echo "  make help          - Show all available commands"
echo "  make dev           - Start development environment"
echo "  make docker-up     - Start all services with Docker"
echo "  make health        - Check service health"
echo ""
echo "🔗 Monitoring URLs:"
echo "  Application:  http://localhost:3000"
echo "  API Gateway:  http://localhost:8080"
echo "  Grafana:      http://localhost:3001 (admin/admin)"
echo "  Prometheus:   http://localhost:9090"
echo "  Jaeger:       http://localhost:16686"
