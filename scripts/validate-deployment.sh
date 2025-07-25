#!/bin/bash

# AI Agentic Browser Deployment Validation Script

set -e

echo "🚀 AI Agentic Browser - Deployment Validation"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if Docker is running
echo "Checking Docker environment..."
if ! docker info > /dev/null 2>&1; then
    print_status 1 "Docker is not running"
    echo "Please start Docker and try again"
    exit 1
fi
print_status 0 "Docker is running"

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_status 1 "docker-compose not found"
    echo "Please install docker-compose"
    exit 1
fi
print_status 0 "docker-compose is available"

# Check environment file
if [ ! -f ".env" ]; then
    print_warning ".env file not found, creating from example..."
    if [ -f ".env.example" ]; then
        cp .env.example .env
        print_info "Please edit .env file with your API keys before continuing"
        echo "Required variables:"
        echo "  - OPENAI_API_KEY"
        echo "  - JWT_SECRET"
        echo "  - ETHEREUM_RPC_URL (optional)"
        echo ""
        read -p "Press Enter after configuring .env file..."
    else
        print_status 1 ".env.example not found"
        exit 1
    fi
fi
print_status 0 ".env file exists"

# Validate critical environment variables
source .env
if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "your-super-secret-jwt-key-change-in-production" ]; then
    print_warning "JWT_SECRET not properly configured"
    echo "Please set a secure JWT_SECRET in .env file"
fi

if [ -z "$OPENAI_API_KEY" ] || [ "$OPENAI_API_KEY" = "your-openai-api-key" ]; then
    print_warning "OPENAI_API_KEY not configured"
    echo "AI features will not work without a valid OpenAI API key"
fi

# Build and start services
echo ""
echo "Building and starting services..."
print_info "This may take a few minutes on first run..."

# Build images
docker-compose build --parallel

if [ $? -eq 0 ]; then
    print_status 0 "Docker images built successfully"
else
    print_status 1 "Failed to build Docker images"
    exit 1
fi

# Start infrastructure services first
print_info "Starting infrastructure services..."
docker-compose up -d postgres redis jaeger prometheus grafana

# Wait for infrastructure to be ready
print_info "Waiting for infrastructure services to be ready..."
sleep 10

# Check PostgreSQL
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
        print_status 0 "PostgreSQL is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        print_status 1 "PostgreSQL failed to start"
        exit 1
    fi
    sleep 1
done

# Check Redis
for i in {1..30}; do
    if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
        print_status 0 "Redis is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        print_status 1 "Redis failed to start"
        exit 1
    fi
    sleep 1
done

# Start application services
print_info "Starting application services..."
docker-compose up -d auth-service ai-agent browser-service web3-service

# Wait for services to start
sleep 15

# Start API Gateway
print_info "Starting API Gateway..."
docker-compose up -d api-gateway

# Wait for API Gateway
sleep 10

# Start Frontend
print_info "Starting Frontend..."
docker-compose up -d frontend

# Wait for all services
sleep 20

# Function to check service health
check_service_health() {
    local service_name=$1
    local health_url=$2
    local max_attempts=30
    local attempt=1

    print_info "Checking $service_name health..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$health_url" > /dev/null 2>&1; then
            print_status 0 "$service_name is healthy"
            return 0
        fi
        
        if [ $((attempt % 5)) -eq 0 ]; then
            echo -n "."
        fi
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_status 1 "$service_name health check failed"
    return 1
}

# Check all service health endpoints
echo ""
echo "Validating service health..."

check_service_health "Auth Service" "http://localhost:8081/health"
check_service_health "AI Agent Service" "http://localhost:8082/health"
check_service_health "Browser Service" "http://localhost:8083/health"
check_service_health "Web3 Service" "http://localhost:8084/health"
check_service_health "API Gateway" "http://localhost:8080/health"
check_service_health "Frontend Application" "http://localhost:3000"

# Test API Gateway status endpoint
echo ""
echo "Testing API Gateway integration..."
if curl -s "http://localhost:8080/api/status" | jq . > /dev/null 2>&1; then
    print_status 0 "API Gateway status endpoint working"
else
    print_status 1 "API Gateway status endpoint failed"
fi

# Test user registration and authentication flow
echo ""
echo "Testing authentication flow..."

# Generate test user data
TEST_EMAIL="test-$(date +%s)@example.com"
TEST_PASSWORD="TestPassword123!"

# Test user registration
print_info "Testing user registration..."
REGISTER_RESPONSE=$(curl -s -w "%{http_code}" -X POST http://localhost:8081/auth/register \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"first_name\":\"Test\",\"last_name\":\"User\"}" \
    -o /tmp/register_response.json)

if [ "$REGISTER_RESPONSE" = "201" ]; then
    print_status 0 "User registration successful"
    
    # Test user login
    print_info "Testing user login..."
    LOGIN_RESPONSE=$(curl -s -w "%{http_code}" -X POST http://localhost:8081/auth/login \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
        -o /tmp/login_response.json)
    
    if [ "$LOGIN_RESPONSE" = "200" ]; then
        print_status 0 "User login successful"
        
        # Extract access token
        ACCESS_TOKEN=$(jq -r '.access_token' /tmp/login_response.json)
        
        if [ "$ACCESS_TOKEN" != "null" ] && [ ! -z "$ACCESS_TOKEN" ]; then
            print_status 0 "Access token obtained"
            
            # Test protected endpoint
            print_info "Testing protected endpoint..."
            PROFILE_RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8081/auth/me \
                -H "Authorization: Bearer $ACCESS_TOKEN" \
                -o /tmp/profile_response.json)
            
            if [ "$PROFILE_RESPONSE" = "200" ]; then
                print_status 0 "Protected endpoint access successful"
            else
                print_status 1 "Protected endpoint access failed"
            fi
        else
            print_status 1 "Failed to extract access token"
        fi
    else
        print_status 1 "User login failed"
    fi
else
    print_status 1 "User registration failed"
fi

# Test AI Agent (if OpenAI key is configured)
if [ ! -z "$OPENAI_API_KEY" ] && [ "$OPENAI_API_KEY" != "your-openai-api-key" ] && [ ! -z "$ACCESS_TOKEN" ]; then
    echo ""
    echo "Testing AI Agent functionality..."
    
    AI_RESPONSE=$(curl -s -w "%{http_code}" -X POST http://localhost:8082/ai/chat \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{"message":"Hello, can you help me test the system?"}' \
        -o /tmp/ai_response.json)
    
    if [ "$AI_RESPONSE" = "200" ]; then
        print_status 0 "AI Agent chat functionality working"
    else
        print_status 1 "AI Agent chat functionality failed"
    fi
else
    print_warning "Skipping AI Agent test (OpenAI API key not configured or no access token)"
fi

# Test Browser Service
if [ ! -z "$ACCESS_TOKEN" ]; then
    echo ""
    echo "Testing Browser Service functionality..."
    
    # Create browser session
    SESSION_RESPONSE=$(curl -s -w "%{http_code}" -X POST http://localhost:8083/browser/sessions \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{"session_name":"Test Session"}' \
        -o /tmp/session_response.json)
    
    if [ "$SESSION_RESPONSE" = "201" ]; then
        print_status 0 "Browser session creation working"
        
        SESSION_ID=$(jq -r '.session.id' /tmp/session_response.json)
        
        # Test navigation
        NAV_RESPONSE=$(curl -s -w "%{http_code}" -X POST http://localhost:8083/browser/navigate \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -H "X-Session-ID: $SESSION_ID" \
            -d '{"url":"https://httpbin.org/get","timeout":10}' \
            -o /tmp/nav_response.json)
        
        if [ "$NAV_RESPONSE" = "200" ]; then
            print_status 0 "Browser navigation working"
        else
            print_status 1 "Browser navigation failed"
        fi
    else
        print_status 1 "Browser session creation failed"
    fi
else
    print_warning "Skipping Browser Service test (no access token)"
fi

# Test Web3 Service
if [ ! -z "$ACCESS_TOKEN" ]; then
    echo ""
    echo "Testing Web3 Service functionality..."
    
    # Test supported chains endpoint
    CHAINS_RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8084/web3/chains \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -o /tmp/chains_response.json)
    
    if [ "$CHAINS_RESPONSE" = "200" ]; then
        print_status 0 "Web3 supported chains endpoint working"
    else
        print_status 1 "Web3 supported chains endpoint failed"
    fi
    
    # Test price endpoint
    PRICES_RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8084/web3/prices \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -o /tmp/prices_response.json)
    
    if [ "$PRICES_RESPONSE" = "200" ]; then
        print_status 0 "Web3 prices endpoint working"
    else
        print_status 1 "Web3 prices endpoint failed"
    fi
else
    print_warning "Skipping Web3 Service test (no access token)"
fi

# Clean up test files
rm -f /tmp/register_response.json /tmp/login_response.json /tmp/profile_response.json
rm -f /tmp/ai_response.json /tmp/session_response.json /tmp/nav_response.json
rm -f /tmp/chains_response.json /tmp/prices_response.json

echo ""
echo "🎉 Deployment validation complete!"
echo ""
echo "📋 Service URLs:"
echo "  Frontend:         http://localhost:3000"
echo "  API Gateway:      http://localhost:8080"
echo "  Auth Service:     http://localhost:8081"
echo "  AI Agent:         http://localhost:8082"
echo "  Browser Service:  http://localhost:8083"
echo "  Web3 Service:     http://localhost:8084"
echo ""
echo "📊 Monitoring URLs:"
echo "  Grafana:          http://localhost:3001 (admin/admin)"
echo "  Prometheus:       http://localhost:9090"
echo "  Jaeger:           http://localhost:16686"
echo ""
echo "🔧 Management Commands:"
echo "  View logs:        docker-compose logs -f"
echo "  Stop services:    docker-compose down"
echo "  Restart:          docker-compose restart"
echo "  Health check:     curl http://localhost:8080/api/status"
echo ""
echo "✨ Your AI Agentic Browser is ready to use!"
