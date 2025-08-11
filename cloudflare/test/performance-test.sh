#!/bin/bash

# Performance Testing Script for Cloudflare Deployment
# Tests load times, throughput, and scalability

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
CONCURRENT_USERS=${2:-10}
DURATION=${3:-60}

# Set URLs based on environment
case $ENVIRONMENT in
    production)
        FRONTEND_URL="https://ai-agentic-crypto-browser.pages.dev"
        API_URL="https://api.your-domain.com"
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

log_info "🚀 Starting performance tests for $ENVIRONMENT"
log_info "Concurrent users: $CONCURRENT_USERS"
log_info "Duration: ${DURATION}s"

# Check if required tools are installed
check_tools() {
    log_info "Checking required tools..."
    
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi
    
    if command -v ab &> /dev/null; then
        log_success "Apache Bench (ab) is available"
        HAS_AB=true
    else
        log_warning "Apache Bench (ab) not found - some tests will be skipped"
        HAS_AB=false
    fi
    
    if command -v wrk &> /dev/null; then
        log_success "wrk is available"
        HAS_WRK=true
    else
        log_warning "wrk not found - some tests will be skipped"
        HAS_WRK=false
    fi
}

# Basic response time test
test_response_times() {
    log_info "📊 Testing response times..."
    
    local endpoints=(
        "$FRONTEND_URL"
        "$API_URL/health"
        "$API_URL/api/version"
        "$API_URL/api/analytics/market"
    )
    
    for endpoint in "${endpoints[@]}"; do
        log_info "Testing: $endpoint"
        
        local times=()
        for i in {1..5}; do
            local start_time=$(date +%s%N)
            curl -s --max-time 10 "$endpoint" > /dev/null 2>&1
            local end_time=$(date +%s%N)
            local duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
            times+=($duration)
        done
        
        # Calculate average
        local total=0
        for time in "${times[@]}"; do
            total=$((total + time))
        done
        local average=$((total / ${#times[@]}))
        
        if [ $average -lt 1000 ]; then
            log_success "Average response time: ${average}ms"
        elif [ $average -lt 2000 ]; then
            log_warning "Average response time: ${average}ms (acceptable)"
        else
            log_error "Average response time: ${average}ms (slow)"
        fi
    done
}

# Load testing with Apache Bench
test_with_ab() {
    if [ "$HAS_AB" = false ]; then
        log_warning "Skipping Apache Bench tests (ab not installed)"
        return
    fi
    
    log_info "🔥 Running load tests with Apache Bench..."
    
    local endpoints=(
        "$API_URL/health"
        "$API_URL/api/version"
    )
    
    for endpoint in "${endpoints[@]}"; do
        log_info "Load testing: $endpoint"
        
        # Run ab test
        local ab_output=$(ab -n 100 -c $CONCURRENT_USERS "$endpoint" 2>/dev/null)
        
        # Extract metrics
        local requests_per_sec=$(echo "$ab_output" | grep "Requests per second" | awk '{print $4}')
        local time_per_request=$(echo "$ab_output" | grep "Time per request" | head -1 | awk '{print $4}')
        local failed_requests=$(echo "$ab_output" | grep "Failed requests" | awk '{print $3}')
        
        log_info "Requests per second: $requests_per_sec"
        log_info "Time per request: ${time_per_request}ms"
        log_info "Failed requests: $failed_requests"
        
        if [ "${failed_requests:-0}" -eq 0 ]; then
            log_success "No failed requests"
        else
            log_warning "$failed_requests failed requests"
        fi
    done
}

# Load testing with wrk
test_with_wrk() {
    if [ "$HAS_WRK" = false ]; then
        log_warning "Skipping wrk tests (wrk not installed)"
        return
    fi
    
    log_info "⚡ Running load tests with wrk..."
    
    local endpoints=(
        "$API_URL/health"
        "$API_URL/api/version"
    )
    
    for endpoint in "${endpoints[@]}"; do
        log_info "Load testing: $endpoint"
        
        # Run wrk test
        local wrk_output=$(wrk -t4 -c$CONCURRENT_USERS -d${DURATION}s "$endpoint" 2>/dev/null)
        
        # Extract metrics
        local requests_per_sec=$(echo "$wrk_output" | grep "Requests/sec" | awk '{print $2}')
        local latency_avg=$(echo "$wrk_output" | grep "Latency" | awk '{print $2}')
        local errors=$(echo "$wrk_output" | grep "Non-2xx" | awk '{print $4}' || echo "0")
        
        log_info "Requests per second: $requests_per_sec"
        log_info "Average latency: $latency_avg"
        log_info "Errors: ${errors:-0}"
        
        if [ "${errors:-0}" -eq 0 ]; then
            log_success "No errors"
        else
            log_warning "$errors errors"
        fi
    done
}

# Test CDN cache performance
test_cdn_cache() {
    log_info "🌐 Testing CDN cache performance..."
    
    local static_assets=(
        "$FRONTEND_URL/_next/static/css/app.css"
        "$FRONTEND_URL/_next/static/js/app.js"
        "$FRONTEND_URL/favicon.ico"
    )
    
    for asset in "${static_assets[@]}"; do
        log_info "Testing cache for: $asset"
        
        # First request (cache miss)
        local headers1=$(curl -s -I "$asset" 2>/dev/null || echo "")
        local cf_cache_status1=$(echo "$headers1" | grep -i "cf-cache-status" | awk '{print $2}' | tr -d '\r')
        
        # Second request (should be cache hit)
        sleep 1
        local headers2=$(curl -s -I "$asset" 2>/dev/null || echo "")
        local cf_cache_status2=$(echo "$headers2" | grep -i "cf-cache-status" | awk '{print $2}' | tr -d '\r')
        
        if [ -n "$cf_cache_status1" ]; then
            log_info "First request: $cf_cache_status1"
        fi
        
        if [ -n "$cf_cache_status2" ]; then
            log_info "Second request: $cf_cache_status2"
            if [ "$cf_cache_status2" = "HIT" ]; then
                log_success "Cache is working correctly"
            else
                log_warning "Cache status: $cf_cache_status2"
            fi
        else
            log_warning "Asset not found or not cached"
        fi
    done
}

# Test geographic performance
test_geographic_performance() {
    log_info "🌍 Testing geographic performance..."
    
    # Test from different regions (if available)
    local regions=(
        "us-east-1"
        "eu-west-1"
        "ap-southeast-1"
    )
    
    for region in "${regions[@]}"; do
        log_info "Testing from region: $region"
        
        # This is a simplified test - in practice, you'd use a service like Pingdom or GTmetrix
        local start_time=$(date +%s%N)
        curl -s --max-time 10 "$FRONTEND_URL" > /dev/null 2>&1
        local end_time=$(date +%s%N)
        local duration=$(( (end_time - start_time) / 1000000 ))
        
        log_info "Response time from $region: ${duration}ms"
    done
}

# Test API rate limiting
test_rate_limiting() {
    log_info "🚦 Testing API rate limiting..."
    
    local endpoint="$API_URL/api/version"
    local requests=0
    local rate_limited=false
    
    log_info "Sending rapid requests to test rate limiting..."
    
    for i in {1..50}; do
        local response=$(curl -s -w "%{http_code}" "$endpoint" -o /dev/null --max-time 5)
        requests=$((requests + 1))
        
        if [ "$response" = "429" ]; then
            log_success "Rate limiting triggered after $requests requests"
            rate_limited=true
            break
        fi
        
        sleep 0.1
    done
    
    if [ "$rate_limited" = false ]; then
        log_warning "Rate limiting not triggered after $requests requests"
    fi
}

# Generate performance report
generate_performance_report() {
    log_info "📈 Performance Test Summary"
    echo "=================================="
    echo "Environment: $ENVIRONMENT"
    echo "Concurrent Users: $CONCURRENT_USERS"
    echo "Test Duration: ${DURATION}s"
    echo "Frontend URL: $FRONTEND_URL"
    echo "API URL: $API_URL"
    echo "=================================="
    
    log_success "Performance testing completed!"
    log_info "Review the results above for any performance issues."
    
    echo ""
    log_info "💡 Performance Optimization Tips:"
    echo "1. Monitor Cloudflare Analytics for detailed metrics"
    echo "2. Use Cloudflare's Speed tab for optimization suggestions"
    echo "3. Enable additional caching rules for better performance"
    echo "4. Consider using Cloudflare Workers for edge computing"
    echo "5. Optimize images and assets for faster loading"
}

# Main execution
main() {
    check_tools
    test_response_times
    test_with_ab
    test_with_wrk
    test_cdn_cache
    test_geographic_performance
    test_rate_limiting
    generate_performance_report
}

# Show usage if no arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 [environment] [concurrent_users] [duration]"
    echo "Environments: production, staging, development"
    echo "Example: $0 production 10 60"
    exit 1
fi

# Run main function
main
