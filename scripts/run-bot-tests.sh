#!/bin/bash

# Bot Testing Framework - Test Runner Script
# AI-Agentic Crypto Browser - Trading Bot Testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8090/api/v1/testing"
TIMEOUT=300
VERBOSE=false
PARALLEL=true

# Default test configuration
DEFAULT_STRATEGIES=("dca" "grid" "momentum" "mean_reversion" "arbitrage" "scalping" "portfolio_rebalancing")
DEFAULT_TEST_TYPES=("unit" "integration" "backtest")

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print usage
print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -v, --verbose           Enable verbose output"
    echo "  -s, --sequential        Run tests sequentially (default: parallel)"
    echo "  -t, --type TYPE         Test type: unit, integration, backtest, stress, performance"
    echo "  -S, --strategy STRATEGY Strategy to test: dca, grid, momentum, etc."
    echo "  -T, --timeout SECONDS   Test timeout in seconds (default: 300)"
    echo "  -u, --url URL           API base URL (default: http://localhost:8090/api/v1/testing)"
    echo "  --unit                  Run unit tests only"
    echo "  --integration           Run integration tests only"
    echo "  --backtest              Run backtests only"
    echo "  --stress                Run stress tests only"
    echo "  --performance           Run performance tests only"
    echo "  --all                   Run all test types"
    echo "  --reset-env             Reset test environment before running"
    echo "  --market-condition COND Set market condition: bull, bear, sideways, volatile"
    echo ""
    echo "Examples:"
    echo "  $0 --unit                           # Run all unit tests"
    echo "  $0 --backtest -S dca               # Run DCA backtest"
    echo "  $0 --all -S grid --sequential      # Run all Grid tests sequentially"
    echo "  $0 --stress --market-condition volatile  # Run stress tests in volatile market"
}

# Function to check if API is available
check_api() {
    print_status $BLUE "Checking API availability..."
    if curl -s -f "${API_BASE_URL}/environment" > /dev/null 2>&1; then
        print_status $GREEN "✓ API is available"
        return 0
    else
        print_status $RED "✗ API is not available at ${API_BASE_URL}"
        print_status $YELLOW "Please ensure the trading bot service is running"
        return 1
    fi
}

# Function to reset test environment
reset_environment() {
    print_status $BLUE "Resetting test environment..."
    if curl -s -X POST "${API_BASE_URL}/environment/reset" > /dev/null 2>&1; then
        print_status $GREEN "✓ Test environment reset"
    else
        print_status $YELLOW "⚠ Failed to reset test environment"
    fi
}

# Function to set market condition
set_market_condition() {
    local condition=$1
    print_status $BLUE "Setting market condition to: $condition"
    
    local response=$(curl -s -X POST "${API_BASE_URL}/environment/market-condition" \
        -H "Content-Type: application/json" \
        -d "{\"condition\": \"$condition\"}")
    
    if echo "$response" | grep -q "updated"; then
        print_status $GREEN "✓ Market condition set to $condition"
    else
        print_status $YELLOW "⚠ Failed to set market condition"
    fi
}

# Function to submit a test
submit_test() {
    local test_type=$1
    local strategy=$2
    local scenario=${3:-""}
    
    local payload="{
        \"type\": \"$test_type\",
        \"strategy\": \"$strategy\",
        \"timeout\": \"${TIMEOUT}s\",
        \"config\": {
            \"trading_pairs\": [\"BTC/USDT\"],
            \"exchange\": \"mock\",
            \"base_currency\": \"USDT\",
            \"strategy_params\": {
                \"strategy\": \"$strategy\"
            },
            \"enabled\": true
        }"
    
    if [ -n "$scenario" ]; then
        payload="${payload}, \"scenario\": $scenario"
    fi
    
    payload="${payload}}"
    
    local response=$(curl -s -X POST "${API_BASE_URL}/tests" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    local test_id=$(echo "$response" | grep -o '"test_id":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$test_id" ]; then
        echo "$test_id"
        return 0
    else
        print_status $RED "✗ Failed to submit test: $test_type/$strategy"
        if [ "$VERBOSE" = true ]; then
            echo "Response: $response"
        fi
        return 1
    fi
}

# Function to wait for test completion
wait_for_test() {
    local test_id=$1
    local max_wait=${2:-$TIMEOUT}
    local start_time=$(date +%s)
    
    while true; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        if [ $elapsed -gt $max_wait ]; then
            print_status $RED "✗ Test $test_id timed out after ${max_wait}s"
            return 1
        fi
        
        local status_response=$(curl -s "${API_BASE_URL}/tests/$test_id")
        local status=$(echo "$status_response" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        
        case "$status" in
            "completed"|"failed"|"cancelled")
                return 0
                ;;
            "running")
                if [ "$VERBOSE" = true ]; then
                    local progress=$(echo "$status_response" | grep -o '"progress":[0-9.]*' | cut -d':' -f2)
                    print_status $BLUE "Test $test_id progress: ${progress}%"
                fi
                sleep 2
                ;;
            *)
                sleep 1
                ;;
        esac
    done
}

# Function to get test result
get_test_result() {
    local test_id=$1
    
    local result_response=$(curl -s "${API_BASE_URL}/tests/$test_id/result")
    local passed=$(echo "$result_response" | grep -o '"passed":[^,]*' | cut -d':' -f2)
    local strategy=$(echo "$result_response" | grep -o '"strategy":"[^"]*"' | cut -d'"' -f4)
    local test_type=$(echo "$result_response" | grep -o '"test_type":"[^"]*"' | cut -d'"' -f4)
    local duration=$(echo "$result_response" | grep -o '"duration":"[^"]*"' | cut -d'"' -f4)
    
    if [ "$passed" = "true" ]; then
        print_status $GREEN "✓ $test_type/$strategy PASSED ($duration)"
        
        if [ "$VERBOSE" = true ]; then
            local metrics=$(echo "$result_response" | grep -o '"metrics":{[^}]*}')
            if [ -n "$metrics" ]; then
                echo "  Metrics: $metrics"
            fi
        fi
        return 0
    else
        print_status $RED "✗ $test_type/$strategy FAILED ($duration)"
        
        local failure_reasons=$(echo "$result_response" | grep -o '"failure_reasons":\[[^\]]*\]')
        if [ -n "$failure_reasons" ]; then
            echo "  Reasons: $failure_reasons"
        fi
        return 1
    fi
}

# Function to run a single test
run_single_test() {
    local test_type=$1
    local strategy=$2
    local scenario=${3:-""}
    
    print_status $BLUE "Running $test_type test for $strategy strategy..."
    
    local test_id=$(submit_test "$test_type" "$strategy" "$scenario")
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    if [ "$VERBOSE" = true ]; then
        print_status $BLUE "Test ID: $test_id"
    fi
    
    wait_for_test "$test_id"
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    get_test_result "$test_id"
    return $?
}

# Function to run test suite
run_test_suite() {
    local test_types=("$@")
    local strategies=("${DEFAULT_STRATEGIES[@]}")
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    print_status $BLUE "Running test suite..."
    print_status $BLUE "Test types: ${test_types[*]}"
    print_status $BLUE "Strategies: ${strategies[*]}"
    print_status $BLUE "Parallel execution: $PARALLEL"
    
    # Calculate total tests
    for test_type in "${test_types[@]}"; do
        for strategy in "${strategies[@]}"; do
            ((total_tests++))
        done
    done
    
    print_status $BLUE "Total tests to run: $total_tests"
    echo ""
    
    # Run tests
    local test_pids=()
    
    for test_type in "${test_types[@]}"; do
        for strategy in "${strategies[@]}"; do
            if [ "$PARALLEL" = true ]; then
                # Run in background
                (
                    if run_single_test "$test_type" "$strategy"; then
                        exit 0
                    else
                        exit 1
                    fi
                ) &
                test_pids+=($!)
            else
                # Run sequentially
                if run_single_test "$test_type" "$strategy"; then
                    ((passed_tests++))
                else
                    ((failed_tests++))
                fi
            fi
        done
    done
    
    # Wait for parallel tests to complete
    if [ "$PARALLEL" = true ]; then
        for pid in "${test_pids[@]}"; do
            if wait $pid; then
                ((passed_tests++))
            else
                ((failed_tests++))
            fi
        done
    fi
    
    # Print summary
    echo ""
    print_status $BLUE "Test Summary:"
    print_status $GREEN "  Passed: $passed_tests"
    print_status $RED "  Failed: $failed_tests"
    print_status $BLUE "  Total:  $total_tests"
    
    local success_rate=$((passed_tests * 100 / total_tests))
    print_status $BLUE "  Success Rate: ${success_rate}%"
    
    if [ $failed_tests -eq 0 ]; then
        print_status $GREEN "🎉 All tests passed!"
        return 0
    else
        print_status $RED "❌ Some tests failed"
        return 1
    fi
}

# Parse command line arguments
TEST_TYPES=()
STRATEGIES=()
MARKET_CONDITION=""
RESET_ENV=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -s|--sequential)
            PARALLEL=false
            shift
            ;;
        -t|--type)
            TEST_TYPES+=("$2")
            shift 2
            ;;
        -S|--strategy)
            STRATEGIES+=("$2")
            shift 2
            ;;
        -T|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -u|--url)
            API_BASE_URL="$2"
            shift 2
            ;;
        --unit)
            TEST_TYPES+=("unit")
            shift
            ;;
        --integration)
            TEST_TYPES+=("integration")
            shift
            ;;
        --backtest)
            TEST_TYPES+=("backtest")
            shift
            ;;
        --stress)
            TEST_TYPES+=("stress")
            shift
            ;;
        --performance)
            TEST_TYPES+=("performance")
            shift
            ;;
        --all)
            TEST_TYPES=("${DEFAULT_TEST_TYPES[@]}")
            shift
            ;;
        --reset-env)
            RESET_ENV=true
            shift
            ;;
        --market-condition)
            MARKET_CONDITION="$2"
            shift 2
            ;;
        *)
            print_status $RED "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Set defaults if not specified
if [ ${#TEST_TYPES[@]} -eq 0 ]; then
    TEST_TYPES=("unit")
fi

if [ ${#STRATEGIES[@]} -eq 0 ]; then
    STRATEGIES=("${DEFAULT_STRATEGIES[@]}")
fi

# Main execution
main() {
    print_status $BLUE "🧪 Bot Testing Framework"
    print_status $BLUE "========================"
    echo ""
    
    # Check API availability
    if ! check_api; then
        exit 1
    fi
    
    # Reset environment if requested
    if [ "$RESET_ENV" = true ]; then
        reset_environment
    fi
    
    # Set market condition if specified
    if [ -n "$MARKET_CONDITION" ]; then
        set_market_condition "$MARKET_CONDITION"
    fi
    
    echo ""
    
    # Run tests
    if run_test_suite "${TEST_TYPES[@]}"; then
        exit 0
    else
        exit 1
    fi
}

# Run main function
main
