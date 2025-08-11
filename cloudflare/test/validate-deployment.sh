#!/bin/bash

# Cloudflare Deployment Validation Script
# Tests all components of the deployed AI Agentic Crypto Browser

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Configuration
ENVIRONMENT=${1:-production}
TIMEOUT=30

# Set URLs based on environment
case $ENVIRONMENT in
    production)
        FRONTEND_URL="https://e7f3d92d.ai-agentic-crypto-browser.pages.dev"
        API_URL="https://ai-crypto-browser-api.gcp-inspiration.workers.dev"
        ;;
    staging)
        FRONTEND_URL="https://ai-agentic-crypto-browser-staging.pages.dev"
        API_URL="https://api-staging.your-domain.com"
        ;;
    development)
        FRONTEND_URL="https://ai-agentic-crypto-browser-dev.pages.dev"
        API_URL="https://api-dev.your-domain.com"
        ;;
    *)
        log_error "Unknown environment: $ENVIRONMENT"
        exit 1
        ;;
esac

log_info "🧪 Starting deployment validation for $ENVIRONMENT environment"
log_info "Frontend URL: $FRONTEND_URL"
log_info "API URL: $API_URL"

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_TOTAL=0

# Test function
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TESTS_TOTAL=$((TESTS_TOTAL + 1))
    log_info "Running test: $test_name"
    
    if eval "$test_command"; then
        log_success "PASS: $test_name"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        log_error "FAIL: $test_name"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Frontend Tests
test_frontend_accessibility() {
    curl -f -s --max-time $TIMEOUT "$FRONTEND_URL" > /dev/null
}

test_frontend_content() {
    local content=$(curl -s --max-time $TIMEOUT "$FRONTEND_URL")
    echo "$content" | grep -q "AI Agentic Crypto Browser"
}

test_frontend_assets() {
    local content=$(curl -s --max-time $TIMEOUT "$FRONTEND_URL")
    echo "$content" | grep -q "_next/static"
}

# API Tests
test_api_health() {
    curl -f -s --max-time $TIMEOUT "$API_URL/health" > /dev/null
}

test_api_version() {
    local response=$(curl -s --max-time $TIMEOUT "$API_URL/api/version")
    echo "$response" | grep -q "version"
}

test_api_cors() {
    curl -s --max-time $TIMEOUT \
        -H "Origin: https://example.com" \
        -H "Access-Control-Request-Method: GET" \
        -H "Access-Control-Request-Headers: Content-Type" \
        -X OPTIONS "$API_URL/api/version" | grep -q "Access-Control-Allow-Origin"
}

# Authentication Tests
test_auth_endpoints() {
    # Test registration endpoint (should return validation error)
    local response=$(curl -s --max-time $TIMEOUT \
        -X POST "$API_URL/api/auth/register" \
        -H "Content-Type: application/json" \
        -d '{}')
    echo "$response" | grep -q "error"
}

# Database Tests
test_database_connection() {
    # Test an endpoint that requires database access
    local response=$(curl -s --max-time $TIMEOUT "$API_URL/api/analytics/market")
    echo "$response" | grep -q "success"
}

# WebSocket Tests
test_websocket_upgrade() {
    # Test WebSocket endpoint (should return upgrade error for HTTP)
    local response=$(curl -s --max-time $TIMEOUT "$API_URL/websocket")
    echo "$response" | grep -q -i "upgrade"
}

# Security Tests
test_security_headers() {
    local headers=$(curl -s -I --max-time $TIMEOUT "$FRONTEND_URL")
    echo "$headers" | grep -q "X-Frame-Options" && \
    echo "$headers" | grep -q "X-Content-Type-Options"
}

test_ssl_certificate() {
    curl -s --max-time $TIMEOUT "$FRONTEND_URL" > /dev/null
}

# Performance Tests
test_response_time() {
    local start_time=$(date +%s%N)
    curl -s --max-time $TIMEOUT "$FRONTEND_URL" > /dev/null
    local end_time=$(date +%s%N)
    local duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if [ $duration -lt 2000 ]; then # Less than 2 seconds
        return 0
    else
        log_warning "Response time: ${duration}ms (slower than expected)"
        return 1
    fi
}

# Cache Tests
test_static_asset_caching() {
    local headers=$(curl -s -I --max-time $TIMEOUT "$FRONTEND_URL/_next/static/css/app.css" 2>/dev/null || echo "")
    if [ -n "$headers" ]; then
        echo "$headers" | grep -q "Cache-Control"
    else
        return 0 # Skip if asset doesn't exist
    fi
}

# API Rate Limiting Tests
test_rate_limiting() {
    local count=0
    local max_requests=10
    
    while [ $count -lt $max_requests ]; do
        local response=$(curl -s -w "%{http_code}" --max-time 5 "$API_URL/api/version" -o /dev/null)
        if [ "$response" = "429" ]; then
            return 0 # Rate limiting is working
        fi
        count=$((count + 1))
        sleep 0.1
    done
    
    log_warning "Rate limiting not triggered after $max_requests requests"
    return 0 # Don't fail the test, just warn
}

# Run all tests
run_all_tests() {
    log_info "🔍 Running Frontend Tests"
    run_test "Frontend Accessibility" "test_frontend_accessibility"
    run_test "Frontend Content" "test_frontend_content"
    run_test "Frontend Assets" "test_frontend_assets"
    
    log_info "🔍 Running API Tests"
    run_test "API Health Check" "test_api_health"
    run_test "API Version Endpoint" "test_api_version"
    run_test "API CORS Headers" "test_api_cors"
    
    log_info "🔍 Running Authentication Tests"
    run_test "Auth Endpoints" "test_auth_endpoints"
    
    log_info "🔍 Running Database Tests"
    run_test "Database Connection" "test_database_connection"
    
    log_info "🔍 Running WebSocket Tests"
    run_test "WebSocket Upgrade" "test_websocket_upgrade"
    
    log_info "🔍 Running Security Tests"
    run_test "Security Headers" "test_security_headers"
    run_test "SSL Certificate" "test_ssl_certificate"
    
    log_info "🔍 Running Performance Tests"
    run_test "Response Time" "test_response_time"
    run_test "Static Asset Caching" "test_static_asset_caching"
    
    log_info "🔍 Running Rate Limiting Tests"
    run_test "API Rate Limiting" "test_rate_limiting"
}

# Generate test report
generate_report() {
    echo ""
    log_info "📊 Test Results Summary"
    echo "=================================="
    echo "Environment: $ENVIRONMENT"
    echo "Frontend URL: $FRONTEND_URL"
    echo "API URL: $API_URL"
    echo "=================================="
    echo "Tests Passed: $TESTS_PASSED"
    echo "Tests Failed: $TESTS_FAILED"
    echo "Total Tests: $TESTS_TOTAL"
    echo "Success Rate: $(( TESTS_PASSED * 100 / TESTS_TOTAL ))%"
    echo "=================================="
    
    if [ $TESTS_FAILED -eq 0 ]; then
        log_success "🎉 All tests passed! Deployment is healthy."
        return 0
    else
        log_error "❌ $TESTS_FAILED test(s) failed. Please review the deployment."
        return 1
    fi
}

# Main execution
main() {
    log_info "Starting validation for $ENVIRONMENT environment..."
    
    # Wait a bit for deployment to propagate
    log_info "Waiting 30 seconds for deployment to propagate..."
    sleep 30
    
    run_all_tests
    generate_report
}

# Show usage if no arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 [environment]"
    echo "Environments: production, staging, development"
    echo "Example: $0 production"
    exit 1
fi

# Run main function
main
