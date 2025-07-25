#!/bin/bash

# AI Agent and Browser Service Test Script

set -e

echo "🤖 Testing AI Agent and Browser Services"
echo "========================================"

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

# Check if required environment variables are set
echo "Checking environment configuration..."

if [ ! -f ".env" ]; then
    print_status 1 ".env file not found"
    echo "Please create .env file from .env.example and configure your API keys"
    exit 1
fi

# Check for OpenAI API key
if grep -q "OPENAI_API_KEY=your-openai-api-key" .env || ! grep -q "OPENAI_API_KEY=" .env; then
    print_warning "OpenAI API key not configured in .env"
    echo "Some AI features may not work without a valid API key"
else
    print_status 0 "OpenAI API key configured"
fi

# Check if services can compile
echo ""
echo "Testing service compilation..."

# Test auth service compilation
if go build -o /tmp/test-auth ./cmd/auth-service > /dev/null 2>&1; then
    print_status 0 "Auth service compiles successfully"
    rm -f /tmp/test-auth
else
    print_status 1 "Auth service compilation failed"
    echo "Run 'go build ./cmd/auth-service' for details"
    exit 1
fi

# Test AI agent compilation
if go build -o /tmp/test-ai ./cmd/ai-agent > /dev/null 2>&1; then
    print_status 0 "AI agent service compiles successfully"
    rm -f /tmp/test-ai
else
    print_status 1 "AI agent service compilation failed"
    echo "Run 'go build ./cmd/ai-agent' for details"
    exit 1
fi

# Test browser service compilation
if go build -o /tmp/test-browser ./cmd/browser-service > /dev/null 2>&1; then
    print_status 0 "Browser service compiles successfully"
    rm -f /tmp/test-browser
else
    print_status 1 "Browser service compilation failed"
    echo "Run 'go build ./cmd/browser-service' for details"
    exit 1
fi

# Check if infrastructure is running
echo ""
echo "Checking infrastructure services..."

# Check PostgreSQL
if docker-compose ps postgres | grep -q "Up"; then
    print_status 0 "PostgreSQL is running"
elif pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    print_status 0 "PostgreSQL is running (local)"
else
    print_status 1 "PostgreSQL is not running"
    echo "Start with: docker-compose up -d postgres"
    exit 1
fi

# Check Redis
if docker-compose ps redis | grep -q "Up"; then
    print_status 0 "Redis is running"
elif redis-cli ping > /dev/null 2>&1; then
    print_status 0 "Redis is running (local)"
else
    print_status 1 "Redis is not running"
    echo "Start with: docker-compose up -d redis"
    exit 1
fi

# Function to wait for service to be ready
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1

    print_info "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            print_status 0 "$service_name is ready"
            return 0
        fi
        
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    print_status 1 "$service_name failed to start within $max_attempts seconds"
    return 1
}

# Function to test service endpoint
test_endpoint() {
    local url=$1
    local expected_status=$2
    local description=$3
    
    response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$url")
    
    if [ "$response" = "$expected_status" ]; then
        print_status 0 "$description"
        return 0
    else
        print_status 1 "$description (got HTTP $response)"
        return 1
    fi
}

# Start testing services
echo ""
echo "Starting service tests..."

# Test 1: Start Auth Service
print_info "Starting auth service..."
go run cmd/auth-service/main.go > /tmp/auth.log 2>&1 &
AUTH_PID=$!

if wait_for_service "http://localhost:8081/health" "Auth service"; then
    # Test auth endpoints
    test_endpoint "http://localhost:8081/health" "200" "Auth service health check"
    
    # Test user registration
    print_info "Testing user registration..."
    response=$(curl -s -w "%{http_code}" -X POST http://localhost:8081/auth/register \
        -H "Content-Type: application/json" \
        -d '{"email":"test@example.com","password":"password123","first_name":"Test","last_name":"User"}' \
        -o /tmp/register.json)
    
    if [ "$response" = "201" ]; then
        print_status 0 "User registration successful"
        
        # Test user login
        print_info "Testing user login..."
        login_response=$(curl -s -w "%{http_code}" -X POST http://localhost:8081/auth/login \
            -H "Content-Type: application/json" \
            -d '{"email":"test@example.com","password":"password123"}' \
            -o /tmp/login.json)
        
        if [ "$login_response" = "200" ]; then
            print_status 0 "User login successful"
            
            # Extract access token for further tests
            ACCESS_TOKEN=$(jq -r '.access_token' /tmp/login.json)
            print_info "Access token obtained for authenticated requests"
        else
            print_status 1 "User login failed"
        fi
    else
        print_status 1 "User registration failed"
    fi
fi

# Test 2: Start AI Agent Service
print_info "Starting AI agent service..."
go run cmd/ai-agent/main.go > /tmp/ai.log 2>&1 &
AI_PID=$!

if wait_for_service "http://localhost:8082/health" "AI agent service"; then
    test_endpoint "http://localhost:8082/health" "200" "AI agent service health check"
    
    if [ ! -z "$ACCESS_TOKEN" ]; then
        print_info "Testing AI chat endpoint..."
        chat_response=$(curl -s -w "%{http_code}" -X POST http://localhost:8082/ai/chat \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -d '{"message":"Hello, can you help me navigate to google.com?"}' \
            -o /tmp/chat.json)
        
        if [ "$chat_response" = "200" ]; then
            print_status 0 "AI chat endpoint working"
        else
            print_status 1 "AI chat endpoint failed (HTTP $chat_response)"
        fi
    else
        print_warning "Skipping AI chat test (no access token)"
    fi
fi

# Test 3: Start Browser Service
print_info "Starting browser service..."
go run cmd/browser-service/main.go > /tmp/browser.log 2>&1 &
BROWSER_PID=$!

if wait_for_service "http://localhost:8083/health" "Browser service"; then
    test_endpoint "http://localhost:8083/health" "200" "Browser service health check"
    
    if [ ! -z "$ACCESS_TOKEN" ]; then
        print_info "Testing browser session creation..."
        session_response=$(curl -s -w "%{http_code}" -X POST http://localhost:8083/browser/sessions \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -d '{"session_name":"Test Session"}' \
            -o /tmp/session.json)
        
        if [ "$session_response" = "201" ]; then
            print_status 0 "Browser session creation working"
            
            # Extract session ID for navigation test
            SESSION_ID=$(jq -r '.session.id' /tmp/session.json)
            
            print_info "Testing browser navigation..."
            nav_response=$(curl -s -w "%{http_code}" -X POST http://localhost:8083/browser/navigate \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $ACCESS_TOKEN" \
                -H "X-Session-ID: $SESSION_ID" \
                -d '{"url":"https://httpbin.org/get","timeout":10}' \
                -o /tmp/navigate.json)
            
            if [ "$nav_response" = "200" ]; then
                print_status 0 "Browser navigation working"
            else
                print_status 1 "Browser navigation failed (HTTP $nav_response)"
            fi
        else
            print_status 1 "Browser session creation failed (HTTP $session_response)"
        fi
    else
        print_warning "Skipping browser tests (no access token)"
    fi
fi

# Cleanup
echo ""
echo "Cleaning up test processes..."

if [ ! -z "$AUTH_PID" ]; then
    kill $AUTH_PID > /dev/null 2>&1 || true
    print_info "Auth service stopped"
fi

if [ ! -z "$AI_PID" ]; then
    kill $AI_PID > /dev/null 2>&1 || true
    print_info "AI agent service stopped"
fi

if [ ! -z "$BROWSER_PID" ]; then
    kill $BROWSER_PID > /dev/null 2>&1 || true
    print_info "Browser service stopped"
fi

# Clean up temporary files
rm -f /tmp/auth.log /tmp/ai.log /tmp/browser.log
rm -f /tmp/response.json /tmp/register.json /tmp/login.json /tmp/chat.json /tmp/session.json /tmp/navigate.json

echo ""
echo "🎉 Service testing complete!"
echo ""
echo "📋 Next Steps:"
echo "1. Start all services: make docker-up"
echo "2. Access the application at http://localhost:3000"
echo "3. Monitor services with: make health"
echo "4. View logs with: docker-compose logs -f"
echo ""
echo "🔗 Service URLs:"
echo "  Auth Service:     http://localhost:8081"
echo "  AI Agent:         http://localhost:8082"
echo "  Browser Service:  http://localhost:8083"
echo "  Web3 Service:     http://localhost:8084"
echo "  API Gateway:      http://localhost:8080"
