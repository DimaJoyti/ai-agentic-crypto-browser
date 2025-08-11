#!/bin/bash

# Cloudflare Deployment Script for AI Agentic Crypto Browser
# Deploys frontend to Pages and backend to Workers

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
PROJECT_ROOT=$(pwd)
WEB_DIR="$PROJECT_ROOT/web"
WORKER_DIR="$PROJECT_ROOT/cloudflare/workers/api"

log_info "Starting Cloudflare deployment for environment: $ENVIRONMENT"

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if wrangler is installed
    if ! command -v wrangler &> /dev/null; then
        log_error "Wrangler CLI is not installed. Install it with: npm install -g wrangler"
        exit 1
    fi
    
    # Check if user is logged in
    if ! wrangler whoami &> /dev/null; then
        log_error "Please login to Cloudflare first: wrangler login"
        exit 1
    fi
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed"
        exit 1
    fi
    
    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed"
        exit 1
    fi
    
    log_success "All prerequisites met"
}

# Setup environment variables
setup_environment() {
    log_info "Setting up environment variables for $ENVIRONMENT..."
    
    case $ENVIRONMENT in
        production)
            export NODE_ENV=production
            export NEXT_PUBLIC_API_URL="https://api.your-domain.com"
            export NEXT_PUBLIC_WS_URL="wss://api.your-domain.com"
            ;;
        staging)
            export NODE_ENV=production
            export NEXT_PUBLIC_API_URL="https://api-staging.your-domain.com"
            export NEXT_PUBLIC_WS_URL="wss://api-staging.your-domain.com"
            ;;
        development)
            export NODE_ENV=development
            export NEXT_PUBLIC_API_URL="https://api-dev.your-domain.com"
            export NEXT_PUBLIC_WS_URL="wss://api-dev.your-domain.com"
            ;;
        *)
            log_error "Unknown environment: $ENVIRONMENT"
            exit 1
            ;;
    esac
    
    export NEXT_PUBLIC_CLOUDFLARE_DEPLOYMENT=true
    
    log_success "Environment variables set for $ENVIRONMENT"
}

# Build frontend
build_frontend() {
    log_info "Building frontend for Cloudflare Pages..."
    
    cd "$WEB_DIR"
    
    # Install dependencies
    log_info "Installing frontend dependencies..."
    npm ci
    
    # Build the application
    log_info "Building Next.js application..."
    npm run build:cloudflare
    
    # Verify build output
    if [ ! -d "out" ]; then
        log_error "Frontend build failed - output directory not found"
        exit 1
    fi
    
    log_success "Frontend built successfully"
    cd "$PROJECT_ROOT"
}

# Deploy frontend to Cloudflare Pages
deploy_frontend() {
    log_info "Deploying frontend to Cloudflare Pages..."
    
    cd "$WEB_DIR"
    
    # Deploy based on environment
    case $ENVIRONMENT in
        production)
            wrangler pages deploy out --project-name=ai-agentic-crypto-browser --compatibility-date=2024-01-01
            ;;
        staging)
            wrangler pages deploy out --project-name=ai-agentic-crypto-browser-staging --compatibility-date=2024-01-01
            ;;
        development)
            wrangler pages deploy out --project-name=ai-agentic-crypto-browser-dev --compatibility-date=2024-01-01
            ;;
    esac
    
    log_success "Frontend deployed to Cloudflare Pages"
    cd "$PROJECT_ROOT"
}

# Setup database
setup_database() {
    log_info "Setting up Cloudflare D1 database..."
    
    cd "$PROJECT_ROOT/cloudflare/database"
    
    # Check if database exists
    DB_NAME="ai-crypto-browser-db"
    if [ "$ENVIRONMENT" != "production" ]; then
        DB_NAME="ai-crypto-browser-db-$ENVIRONMENT"
    fi
    
    # Run migrations
    log_info "Running database migrations..."
    wrangler d1 execute $DB_NAME --file=./migrations/001_initial_schema.sql || true
    wrangler d1 execute $DB_NAME --file=./migrations/002_trading_tables.sql || true
    wrangler d1 execute $DB_NAME --file=./migrations/003_ai_analytics_tables.sql || true
    wrangler d1 execute $DB_NAME --file=./migrations/004_user_preferences.sql || true
    
    log_success "Database setup completed"
    cd "$PROJECT_ROOT"
}

# Deploy Workers
deploy_workers() {
    log_info "Deploying Cloudflare Workers..."
    
    cd "$WORKER_DIR"
    
    # Install dependencies
    log_info "Installing Worker dependencies..."
    npm ci
    
    # Deploy based on environment
    case $ENVIRONMENT in
        production)
            wrangler deploy --env production
            ;;
        staging)
            wrangler deploy --env staging
            ;;
        development)
            wrangler deploy --env development
            ;;
    esac
    
    log_success "Workers deployed successfully"
    cd "$PROJECT_ROOT"
}

# Setup KV namespaces
setup_kv() {
    log_info "Setting up KV namespaces..."
    
    # This would typically be done once during initial setup
    log_warning "KV namespaces should be created manually or via setup script"
    log_info "Run: ./cloudflare/kv/setup.sh"
}

# Verify deployment
verify_deployment() {
    log_info "Verifying deployment..."
    
    # Test frontend
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
    esac
    
    log_info "Testing frontend at: $FRONTEND_URL"
    if curl -f -s "$FRONTEND_URL" > /dev/null; then
        log_success "Frontend is accessible"
    else
        log_warning "Frontend may not be ready yet (this is normal for new deployments)"
    fi
    
    log_info "Testing API at: $API_URL/health"
    if curl -f -s "$API_URL/health" > /dev/null; then
        log_success "API is accessible"
    else
        log_warning "API may not be ready yet (this is normal for new deployments)"
    fi
}

# Cleanup function
cleanup() {
    log_info "Cleaning up temporary files..."
    # Add any cleanup logic here
    log_success "Cleanup completed"
}

# Main deployment function
main() {
    log_info "🚀 Starting Cloudflare deployment for AI Agentic Crypto Browser"
    log_info "Environment: $ENVIRONMENT"
    log_info "Project root: $PROJECT_ROOT"
    
    # Set trap for cleanup on exit
    trap cleanup EXIT
    
    # Run deployment steps
    check_prerequisites
    setup_environment
    build_frontend
    deploy_frontend
    setup_database
    deploy_workers
    verify_deployment
    
    log_success "🎉 Deployment completed successfully!"
    log_info "Frontend URL: $FRONTEND_URL"
    log_info "API URL: $API_URL"
    
    echo ""
    log_info "📋 Post-deployment checklist:"
    echo "1. Update DNS records to point to your custom domain"
    echo "2. Configure SSL certificates"
    echo "3. Set up monitoring and alerts"
    echo "4. Test all functionality"
    echo "5. Update documentation with new URLs"
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
