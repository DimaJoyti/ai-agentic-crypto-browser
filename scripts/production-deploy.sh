#!/bin/bash

# AI Agentic Browser - Production Deployment Script

set -e

echo "🚀 AI Agentic Browser - Production Deployment"
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

# Check deployment type
DEPLOYMENT_TYPE=${1:-docker}
ENVIRONMENT=${2:-production}

print_info "Deployment Type: $DEPLOYMENT_TYPE"
print_info "Environment: $ENVIRONMENT"

# Validate environment
if [ "$ENVIRONMENT" != "production" ] && [ "$ENVIRONMENT" != "staging" ]; then
    print_status 1 "Invalid environment. Use 'production' or 'staging'"
    exit 1
fi

# Check prerequisites
echo ""
echo "Checking prerequisites..."

# Check Docker
if ! command -v docker &> /dev/null; then
    print_status 1 "Docker not found"
    exit 1
fi
print_status 0 "Docker is available"

# Check docker-compose for Docker deployment
if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
    if ! command -v docker-compose &> /dev/null; then
        print_status 1 "docker-compose not found"
        exit 1
    fi
    print_status 0 "docker-compose is available"
fi

# Check kubectl for Kubernetes deployment
if [ "$DEPLOYMENT_TYPE" = "k8s" ]; then
    if ! command -v kubectl &> /dev/null; then
        print_status 1 "kubectl not found"
        exit 1
    fi
    print_status 0 "kubectl is available"
    
    # Check cluster connection
    if ! kubectl cluster-info &> /dev/null; then
        print_status 1 "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    print_status 0 "Kubernetes cluster is accessible"
fi

# Environment configuration
echo ""
echo "Configuring environment..."

ENV_FILE=".env.${ENVIRONMENT}"
if [ ! -f "$ENV_FILE" ]; then
    print_warning "$ENV_FILE not found, creating from template..."
    
    cat > "$ENV_FILE" << EOF
# Production Environment Configuration
ENVIRONMENT=$ENVIRONMENT

# Database
POSTGRES_PASSWORD=your-secure-postgres-password-here

# JWT
JWT_SECRET=your-super-secure-jwt-secret-here

# AI Services
OPENAI_API_KEY=your-openai-api-key-here
ANTHROPIC_API_KEY=your-anthropic-api-key-here
AI_MODEL=gpt-4-turbo-preview
MAX_TOKENS=2000
TEMPERATURE=0.7

# Web3 RPC URLs
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/your-project-id
POLYGON_RPC_URL=https://polygon-mainnet.infura.io/v3/your-project-id
ARBITRUM_RPC_URL=https://arbitrum-mainnet.infura.io/v3/your-project-id
OPTIMISM_RPC_URL=https://optimism-mainnet.infura.io/v3/your-project-id

# External APIs
COINGECKO_API_KEY=your-coingecko-api-key
ALCHEMY_API_KEY=your-alchemy-api-key
WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id

# Security
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
BCRYPT_COST=12

# Rate Limiting
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=60

# Browser Service
CHROME_HEADLESS=true
CHROME_NO_SANDBOX=true
MAX_CONCURRENT_SESSIONS=10

# Frontend
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_WS_URL=wss://api.yourdomain.com
NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id

# Monitoring
GRAFANA_PASSWORD=your-grafana-password-here
EOF

    print_warning "Please edit $ENV_FILE with your actual configuration values"
    read -p "Press Enter after configuring the environment file..."
fi

# Load environment variables
source "$ENV_FILE"
print_status 0 "Environment configuration loaded"

# Validate critical environment variables
echo ""
echo "Validating configuration..."

REQUIRED_VARS=(
    "POSTGRES_PASSWORD"
    "JWT_SECRET"
    "OPENAI_API_KEY"
)

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ] || [ "${!var}" = "your-${var,,}-here" ]; then
        print_status 1 "$var is not properly configured"
        exit 1
    fi
done
print_status 0 "Critical environment variables are configured"

# Build images
echo ""
echo "Building application images..."

if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
    print_info "Building Docker images..."
    docker-compose -f deployments/docker-compose.prod.yml build --parallel
    
    if [ $? -eq 0 ]; then
        print_status 0 "Docker images built successfully"
    else
        print_status 1 "Failed to build Docker images"
        exit 1
    fi
elif [ "$DEPLOYMENT_TYPE" = "k8s" ]; then
    print_info "Building and pushing images for Kubernetes..."
    
    # Build and tag images
    REGISTRY=${DOCKER_REGISTRY:-"your-registry.com"}
    VERSION=${VERSION:-"latest"}
    
    services=("auth-service" "ai-agent" "browser-service" "web3-service" "api-gateway")
    
    for service in "${services[@]}"; do
        print_info "Building $service..."
        docker build -t "$REGISTRY/ai-agentic-browser/$service:$VERSION" -f "cmd/$service/Dockerfile" .
        
        if [ $? -eq 0 ]; then
            print_status 0 "$service image built"
            
            # Push to registry
            print_info "Pushing $service to registry..."
            docker push "$REGISTRY/ai-agentic-browser/$service:$VERSION"
            
            if [ $? -eq 0 ]; then
                print_status 0 "$service image pushed"
            else
                print_status 1 "Failed to push $service image"
                exit 1
            fi
        else
            print_status 1 "Failed to build $service image"
            exit 1
        fi
    done
    
    # Build frontend
    print_info "Building frontend..."
    docker build -t "$REGISTRY/ai-agentic-browser/frontend:$VERSION" \
        --build-arg NEXT_PUBLIC_API_URL="$NEXT_PUBLIC_API_URL" \
        --build-arg NEXT_PUBLIC_WS_URL="$NEXT_PUBLIC_WS_URL" \
        --build-arg NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID="$NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID" \
        web/
    
    if [ $? -eq 0 ]; then
        print_status 0 "Frontend image built"
        docker push "$REGISTRY/ai-agentic-browser/frontend:$VERSION"
        print_status 0 "Frontend image pushed"
    else
        print_status 1 "Failed to build frontend image"
        exit 1
    fi
fi

# Deploy application
echo ""
echo "Deploying application..."

if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
    print_info "Deploying with Docker Compose..."
    
    # Stop existing services
    docker-compose -f deployments/docker-compose.prod.yml down
    
    # Start services
    docker-compose -f deployments/docker-compose.prod.yml up -d
    
    if [ $? -eq 0 ]; then
        print_status 0 "Application deployed successfully"
    else
        print_status 1 "Deployment failed"
        exit 1
    fi
    
elif [ "$DEPLOYMENT_TYPE" = "k8s" ]; then
    print_info "Deploying to Kubernetes..."
    
    # Apply namespace and configs
    kubectl apply -f deployments/k8s/namespace.yaml
    
    # Update secrets with actual values
    kubectl create secret generic app-secrets \
        --from-literal=POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
        --from-literal=JWT_SECRET="$JWT_SECRET" \
        --from-literal=OPENAI_API_KEY="$OPENAI_API_KEY" \
        --from-literal=ANTHROPIC_API_KEY="$ANTHROPIC_API_KEY" \
        --from-literal=ETHEREUM_RPC_URL="$ETHEREUM_RPC_URL" \
        --from-literal=POLYGON_RPC_URL="$POLYGON_RPC_URL" \
        --from-literal=ARBITRUM_RPC_URL="$ARBITRUM_RPC_URL" \
        --from-literal=OPTIMISM_RPC_URL="$OPTIMISM_RPC_URL" \
        --from-literal=COINGECKO_API_KEY="$COINGECKO_API_KEY" \
        --from-literal=ALCHEMY_API_KEY="$ALCHEMY_API_KEY" \
        --from-literal=WALLETCONNECT_PROJECT_ID="$WALLETCONNECT_PROJECT_ID" \
        --from-literal=GRAFANA_PASSWORD="$GRAFANA_PASSWORD" \
        --namespace=agentic-browser \
        --dry-run=client -o yaml | kubectl apply -f -
    
    # Deploy infrastructure
    kubectl apply -f deployments/k8s/postgres.yaml
    kubectl apply -f deployments/k8s/redis.yaml
    
    # Wait for infrastructure
    print_info "Waiting for infrastructure to be ready..."
    kubectl wait --for=condition=ready pod -l app=postgres -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=redis -n agentic-browser --timeout=300s
    
    # Deploy services
    kubectl apply -f deployments/k8s/ai-agent.yaml
    kubectl apply -f deployments/k8s/browser-service.yaml
    kubectl apply -f deployments/k8s/web3-service.yaml
    kubectl apply -f deployments/k8s/auth-service.yaml
    kubectl apply -f deployments/k8s/api-gateway.yaml
    kubectl apply -f deployments/k8s/frontend.yaml
    
    # Wait for services
    print_info "Waiting for services to be ready..."
    kubectl wait --for=condition=ready pod -l app=ai-agent -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=browser-service -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=web3-service -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=auth-service -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=api-gateway -n agentic-browser --timeout=300s
    kubectl wait --for=condition=ready pod -l app=frontend -n agentic-browser --timeout=300s
    
    print_status 0 "Kubernetes deployment completed"
fi

# Health checks
echo ""
echo "Performing health checks..."

sleep 30  # Wait for services to fully start

if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
    # Check Docker services
    services=("postgres" "redis" "auth-service" "ai-agent" "browser-service" "web3-service" "api-gateway" "frontend")
    
    for service in "${services[@]}"; do
        if docker-compose -f deployments/docker-compose.prod.yml ps "$service" | grep -q "Up"; then
            print_status 0 "$service is running"
        else
            print_status 1 "$service is not running"
        fi
    done
    
    # Test API endpoints
    if curl -s http://localhost:8080/health > /dev/null; then
        print_status 0 "API Gateway health check passed"
    else
        print_status 1 "API Gateway health check failed"
    fi
    
    if curl -s http://localhost:3000 > /dev/null; then
        print_status 0 "Frontend health check passed"
    else
        print_status 1 "Frontend health check failed"
    fi
    
elif [ "$DEPLOYMENT_TYPE" = "k8s" ]; then
    # Check Kubernetes pods
    if kubectl get pods -n agentic-browser | grep -q "Running"; then
        print_status 0 "Kubernetes pods are running"
        kubectl get pods -n agentic-browser
    else
        print_status 1 "Some Kubernetes pods are not running"
        kubectl get pods -n agentic-browser
    fi
fi

# Final summary
echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📋 Deployment Summary:"
echo "  Environment: $ENVIRONMENT"
echo "  Deployment Type: $DEPLOYMENT_TYPE"
echo "  Timestamp: $(date)"
echo ""

if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
    echo "🔗 Service URLs:"
    echo "  Frontend:         http://localhost:3000"
    echo "  API Gateway:      http://localhost:8080"
    echo "  Grafana:          http://localhost:3001 (admin/$GRAFANA_PASSWORD)"
    echo "  Prometheus:       http://localhost:9090"
    echo "  Jaeger:           http://localhost:16686"
    echo ""
    echo "🔧 Management Commands:"
    echo "  View logs:        docker-compose -f deployments/docker-compose.prod.yml logs -f"
    echo "  Stop services:    docker-compose -f deployments/docker-compose.prod.yml down"
    echo "  Restart:          docker-compose -f deployments/docker-compose.prod.yml restart"
    echo "  Scale service:    docker-compose -f deployments/docker-compose.prod.yml up -d --scale ai-agent=5"
elif [ "$DEPLOYMENT_TYPE" = "k8s" ]; then
    echo "🔗 Kubernetes Resources:"
    echo "  Namespace:        agentic-browser"
    echo "  Services:         kubectl get svc -n agentic-browser"
    echo "  Pods:             kubectl get pods -n agentic-browser"
    echo ""
    echo "🔧 Management Commands:"
    echo "  View logs:        kubectl logs -f deployment/ai-agent -n agentic-browser"
    echo "  Scale service:    kubectl scale deployment ai-agent --replicas=5 -n agentic-browser"
    echo "  Port forward:     kubectl port-forward svc/frontend-service 3000:3000 -n agentic-browser"
    echo "  Delete:           kubectl delete namespace agentic-browser"
fi

echo ""
echo "✨ Your AI Agentic Browser is now running in production!"
