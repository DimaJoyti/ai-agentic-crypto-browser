#!/bin/bash

# AI-Agentic Crypto Browser - Cloudflare Pages Deployment Script
# This script handles deployment to Cloudflare Pages with proper cache management

set -e  # Exit on any error

echo "🚀 Starting Cloudflare Pages deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "web/package.json" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

# Set environment variables for production
export NODE_ENV=production
export NEXT_PUBLIC_CLOUDFLARE_DEPLOYMENT=true
export NEXT_PUBLIC_API_BASE_URL=https://ai-crypto-browser-api.gcp-inspiration.workers.dev
export NEXT_PUBLIC_WS_URL=wss://ai-crypto-browser-api.gcp-inspiration.workers.dev
export NEXT_PUBLIC_CHAIN_ID=1

# Handle npm cache issues in CI environments
if [ -n "${CI:-}" ] || [ -n "${GITHUB_ACTIONS:-}" ] || [ -n "${RUNNER_OS:-}" ]; then
    print_status "Detected CI environment, configuring npm cache..."
    
    # Create npm cache directory
    NPM_CACHE_DIR="/tmp/.npm-cache"
    mkdir -p "$NPM_CACHE_DIR"
    export NPM_CONFIG_CACHE="$NPM_CACHE_DIR"
    
    # Set additional npm config for CI
    npm config set cache "$NPM_CACHE_DIR" --global
    npm config set prefer-offline true --global
    npm config set audit false --global
    npm config set fund false --global
    npm config set progress false --global
    
    print_status "NPM cache configured at $NPM_CACHE_DIR"
else
    print_status "Local environment detected, using default npm cache"
fi

# Change to web directory
cd web

print_status "Installing dependencies..."

# Clear npm cache if it's corrupted
npm cache clean --force 2>/dev/null || true

# Install dependencies with CI-friendly options
if [ -n "${CI:-}" ]; then
    # CI environment - use npm ci
    npm ci \
        --prefer-offline \
        --no-audit \
        --no-fund \
        --progress=false \
        --loglevel=error
else
    # Local environment - use npm install
    npm install \
        --prefer-offline \
        --no-audit \
        --no-fund \
        --loglevel=error
fi

if [ $? -eq 0 ]; then
    print_success "Dependencies installed successfully"
else
    print_error "Failed to install dependencies"
    exit 1
fi

print_status "Building application for production..."

# Build the application
npm run build:cloudflare

if [ $? -eq 0 ]; then
    print_success "Build completed successfully"
else
    print_error "Build failed"
    exit 1
fi

# Check if out directory exists
if [ ! -d "out" ]; then
    print_error "Build output directory 'out' not found"
    exit 1
fi

FILE_COUNT=$(find out -type f | wc -l)
print_status "Build output directory contains $FILE_COUNT files"

# Deploy to Cloudflare Pages
print_status "Deploying to Cloudflare Pages..."

# Check if wrangler is available
if ! command -v wrangler &> /dev/null; then
    print_status "Installing Wrangler CLI..."
    if [ -n "${CI:-}" ]; then
        npm install -g wrangler@latest
    else
        npx wrangler --version || npm install -g wrangler@latest
    fi
fi

# Deploy using wrangler
if [ -n "${CI:-}" ]; then
    # CI environment - use environment variables for auth
    wrangler pages deploy out \
        --project-name=ai-agentic-crypto-browser \
        --compatibility-date=2024-01-01
else
    # Local environment - interactive auth
    wrangler pages deploy out \
        --project-name=ai-agentic-crypto-browser \
        --compatibility-date=2024-01-01
fi

if [ $? -eq 0 ]; then
    print_success "🎉 Deployment completed successfully!"
    print_success "🌐 Your application is live at:"
    print_success "   Production: https://01becb04.ai-agentic-crypto-browser.pages.dev"
    print_success "   Development: https://dev.ai-agentic-crypto-browser.pages.dev"
else
    print_error "Deployment failed"
    exit 1
fi

# Cleanup temporary files in CI
if [ -n "${CI:-}" ] && [ -n "${NPM_CACHE_DIR:-}" ]; then
    print_status "Cleaning up CI cache..."
    rm -rf "$NPM_CACHE_DIR" 2>/dev/null || true
fi

print_success "✨ Deployment process completed!"

# Optional: Run post-deployment tests
if [ "$1" = "--test" ]; then
    print_status "Running post-deployment tests..."
    
    # Test if the site is accessible
    SITE_URL="https://01becb04.ai-agentic-crypto-browser.pages.dev"
    
    if command -v curl &> /dev/null; then
        if curl -f -s "$SITE_URL" > /dev/null; then
            print_success "✅ Site is accessible at $SITE_URL"
        else
            print_warning "⚠️  Site might not be ready yet. Please check manually."
        fi
    else
        print_warning "curl not available, skipping connectivity test"
    fi
fi

echo ""
echo "🚀 Deployment Summary:"
echo "   Platform: Cloudflare Pages"
echo "   Status: ✅ Deployed"
echo "   URL: https://01becb04.ai-agentic-crypto-browser.pages.dev"
echo "   Build: Static export (Next.js)"
echo "   Files: $FILE_COUNT static files"
echo "   Cache: Optimized for global CDN"
echo ""
echo "🎯 Next Steps:"
echo "   1. Test the application functionality"
echo "   2. Verify API connectivity"
echo "   3. Check mobile responsiveness"
echo "   4. Monitor performance metrics"
echo ""
