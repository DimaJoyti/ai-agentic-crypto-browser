#!/bin/bash

# AI Agentic Browser - Docker Cleanup Script
# This script stops and removes all Docker containers and networks

set -e

echo "🧹 Cleaning up AI Agentic Browser infrastructure..."

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

# Stop and remove containers
containers=("postgres" "redis" "jaeger" "prometheus" "grafana")

for container in "${containers[@]}"; do
    if docker ps -q -f name="$container" | grep -q .; then
        print_status "Stopping $container..."
        docker stop "$container" 2>/dev/null || print_warning "Failed to stop $container"
    fi
    
    if docker ps -aq -f name="$container" | grep -q .; then
        print_status "Removing $container..."
        docker rm "$container" 2>/dev/null || print_warning "Failed to remove $container"
    fi
done

# Remove network
print_status "Removing Docker network..."
docker network rm ai-browser-network 2>/dev/null || print_warning "Network ai-browser-network does not exist or is in use"

# Optional: Remove volumes (uncomment if you want to delete data)
# print_warning "Removing Docker volumes (this will delete all data)..."
# docker volume rm postgres_data redis_data grafana_data 2>/dev/null || print_warning "Some volumes may not exist"

print_status "🎉 Cleanup complete!"
echo ""
echo "💡 Note: Docker volumes are preserved to keep your data."
echo "   To remove volumes and delete all data, run:"
echo "   docker volume rm postgres_data redis_data grafana_data"
