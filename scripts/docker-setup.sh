#!/bin/bash

# AI Agentic Browser - Docker Setup Script
# This script sets up all required infrastructure services using Docker

set -e

echo "🚀 Setting up AI Agentic Browser infrastructure..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

print_status "Creating Docker network..."
docker network create ai-browser-network 2>/dev/null || print_warning "Network ai-browser-network already exists"

print_status "Starting PostgreSQL..."
docker run -d --name postgres \
    --network ai-browser-network \
    -e POSTGRES_DB=ai_agentic_browser \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    -v postgres_data:/var/lib/postgresql/data \
    postgres:16 2>/dev/null || print_warning "PostgreSQL container already exists"

print_status "Starting Redis..."
docker run -d --name redis \
    --network ai-browser-network \
    -p 6379:6379 \
    -v redis_data:/data \
    redis:7-alpine 2>/dev/null || print_warning "Redis container already exists"

print_status "Starting Jaeger (distributed tracing)..."
docker run -d --name jaeger \
    --network ai-browser-network \
    -p 16686:16686 \
    -p 14268:14268 \
    jaegertracing/all-in-one:latest 2>/dev/null || print_warning "Jaeger container already exists"

print_status "Starting Prometheus (metrics)..."
docker run -d --name prometheus \
    --network ai-browser-network \
    -p 9090:9090 \
    prom/prometheus:latest 2>/dev/null || print_warning "Prometheus container already exists"

print_status "Starting Grafana (dashboards)..."
docker run -d --name grafana \
    --network ai-browser-network \
    -p 3001:3000 \
    -e GF_SECURITY_ADMIN_PASSWORD=admin \
    -v grafana_data:/var/lib/grafana \
    grafana/grafana:latest 2>/dev/null || print_warning "Grafana container already exists"

# Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Check service health
print_status "Checking service health..."

# Check PostgreSQL
if docker exec postgres pg_isready -U postgres &> /dev/null; then
    print_status "✅ PostgreSQL is ready"
else
    print_warning "⚠️  PostgreSQL is not ready yet"
fi

# Check Redis
if docker exec redis redis-cli ping &> /dev/null; then
    print_status "✅ Redis is ready"
else
    print_warning "⚠️  Redis is not ready yet"
fi

# Check if Jaeger is responding
if curl -s http://localhost:16686 &> /dev/null; then
    print_status "✅ Jaeger is ready"
else
    print_warning "⚠️  Jaeger is not ready yet"
fi

# Check if Prometheus is responding
if curl -s http://localhost:9090 &> /dev/null; then
    print_status "✅ Prometheus is ready"
else
    print_warning "⚠️  Prometheus is not ready yet"
fi

# Check if Grafana is responding
if curl -s http://localhost:3001 &> /dev/null; then
    print_status "✅ Grafana is ready"
else
    print_warning "⚠️  Grafana is not ready yet"
fi

echo ""
print_status "🎉 Infrastructure setup complete!"
echo ""
echo "📊 Access Points:"
echo "  - PostgreSQL: localhost:5432 (postgres/postgres)"
echo "  - Redis: localhost:6379"
echo "  - Jaeger UI: http://localhost:16686"
echo "  - Prometheus: http://localhost:9090"
echo "  - Grafana: http://localhost:3001 (admin/admin)"
echo ""
echo "🔧 Next steps:"
echo "  1. Run 'go mod tidy' to install Go dependencies"
echo "  2. Start the backend services with 'go run cmd/*/main.go'"
echo "  3. Start the frontend with 'cd web && npm install && npm run dev'"
echo ""
echo "🛑 To stop all services:"
echo "  ./scripts/docker-cleanup.sh"
