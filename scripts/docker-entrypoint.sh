#!/bin/sh
set -e

# Docker entrypoint script for AI Agentic Browser services
# This script handles service selection and initialization

SERVICE=${1:-api-gateway}
PORT=${PORT:-8080}

echo "Starting AI Agentic Browser - $SERVICE"
echo "Port: $PORT"
echo "Environment: ${ENVIRONMENT:-development}"

# Wait for dependencies
wait_for_service() {
    local host=$1
    local port=$2
    local service_name=$3
    
    echo "Waiting for $service_name at $host:$port..."
    
    while ! nc -z "$host" "$port"; do
        echo "Waiting for $service_name to be ready..."
        sleep 2
    done
    
    echo "$service_name is ready!"
}

# Wait for database if DATABASE_URL is set
if [ -n "$DATABASE_URL" ]; then
    DB_HOST=$(echo "$DATABASE_URL" | sed -n 's/.*@\([^:]*\):.*/\1/p')
    DB_PORT=$(echo "$DATABASE_URL" | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')
    
    if [ -n "$DB_HOST" ] && [ -n "$DB_PORT" ]; then
        wait_for_service "$DB_HOST" "$DB_PORT" "PostgreSQL"
    fi
fi

# Wait for Redis if REDIS_URL is set
if [ -n "$REDIS_URL" ]; then
    REDIS_HOST=$(echo "$REDIS_URL" | sed -n 's/.*\/\/\([^:]*\):.*/\1/p')
    REDIS_PORT=$(echo "$REDIS_URL" | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')
    
    if [ -n "$REDIS_HOST" ] && [ -n "$REDIS_PORT" ]; then
        wait_for_service "$REDIS_HOST" "$REDIS_PORT" "Redis"
    fi
fi

# Run database migrations if this is the API Gateway
if [ "$SERVICE" = "api-gateway" ] && [ -n "$DATABASE_URL" ]; then
    echo "Running database migrations..."
    if [ -f "./migrate" ]; then
        ./migrate -path ./migrations -database "$DATABASE_URL" up
    else
        echo "Migration tool not found, skipping migrations"
    fi
fi

# Set up logging directory
mkdir -p /app/logs
touch /app/logs/${SERVICE}.log

# Start the appropriate service
case "$SERVICE" in
    "api-gateway")
        echo "Starting API Gateway..."
        exec ./bin/api-gateway
        ;;
    "ai-agent")
        echo "Starting AI Agent..."
        exec ./bin/ai-agent
        ;;
    "browser-service")
        echo "Starting Browser Service..."
        exec ./bin/browser-service
        ;;
    "web3-service")
        echo "Starting Web3 Service..."
        exec ./bin/web3-service
        ;;
    "auth-service")
        echo "Starting Auth Service..."
        exec ./bin/auth-service
        ;;
    *)
        echo "Unknown service: $SERVICE"
        echo "Available services: api-gateway, ai-agent, browser-service, web3-service, auth-service"
        exit 1
        ;;
esac
