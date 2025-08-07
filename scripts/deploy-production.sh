#!/bin/bash

# AI Agentic Browser - Production Deployment Script
# This script handles the complete production deployment process

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENVIRONMENT="${ENVIRONMENT:-production}"
NAMESPACE="${NAMESPACE:-ai-browser}"
HELM_RELEASE="${HELM_RELEASE:-ai-browser}"
DOCKER_REGISTRY="${DOCKER_REGISTRY:-ghcr.io/dimajoyti/ai-agentic-crypto-browser}"
KUBECTL_TIMEOUT="${KUBECTL_TIMEOUT:-600s}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Error handling
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log_error "Deployment failed with exit code $exit_code"
        log_info "Rolling back deployment..."
        helm rollback "$HELM_RELEASE" --namespace "$NAMESPACE" || true
    fi
    exit $exit_code
}

trap cleanup EXIT

# Validation functions
validate_prerequisites() {
    log_info "Validating prerequisites..."
    
    # Check required tools
    local required_tools=("kubectl" "helm" "docker" "aws" "jq")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "Required tool '$tool' is not installed"
            exit 1
        fi
    done
    
    # Check kubectl connection
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    # Check Helm
    if ! helm version &> /dev/null; then
        log_error "Helm is not properly configured"
        exit 1
    fi
    
    # Check Docker registry access
    if ! docker info &> /dev/null; then
        log_error "Docker is not running or accessible"
        exit 1
    fi
    
    log_success "Prerequisites validated"
}

validate_environment() {
    log_info "Validating environment configuration..."
    
    # Check required environment variables
    local required_vars=("AWS_REGION" "ENVIRONMENT")
    for var in "${required_vars[@]}"; do
        if [ -z "${!var:-}" ]; then
            log_error "Required environment variable '$var' is not set"
            exit 1
        fi
    done
    
    # Validate environment value
    if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
        log_error "Invalid environment: $ENVIRONMENT. Must be development, staging, or production"
        exit 1
    fi
    
    log_success "Environment configuration validated"
}

# Pre-deployment checks
pre_deployment_checks() {
    log_info "Running pre-deployment checks..."
    
    # Check if namespace exists
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        log_info "Creating namespace: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
        kubectl label namespace "$NAMESPACE" environment="$ENVIRONMENT"
    fi
    
    # Check cluster resources
    local node_count
    node_count=$(kubectl get nodes --no-headers | wc -l)
    if [ "$node_count" -lt 3 ]; then
        log_warning "Cluster has only $node_count nodes. Recommended minimum is 3 for production"
    fi
    
    # Check storage classes
    if ! kubectl get storageclass gp2 &> /dev/null; then
        log_warning "Storage class 'gp2' not found. This may cause PVC issues"
    fi
    
    log_success "Pre-deployment checks completed"
}

# Build and push images
build_and_push_images() {
    log_info "Building and pushing Docker images..."
    
    local image_tag="${GITHUB_SHA:-$(git rev-parse HEAD)}"
    local components=("api-gateway" "ai-agent" "auth-service" "browser-service" "web3-service" "frontend")
    
    for component in "${components[@]}"; do
        log_info "Building $component..."
        
        local image_name="$DOCKER_REGISTRY/$component:$image_tag"
        local dockerfile="$PROJECT_ROOT/deployments/docker/Dockerfile.$component"
        
        if [ ! -f "$dockerfile" ]; then
            log_warning "Dockerfile not found for $component: $dockerfile"
            continue
        fi
        
        # Build image
        docker build \
            -t "$image_name" \
            -f "$dockerfile" \
            "$PROJECT_ROOT"
        
        # Push image
        docker push "$image_name"
        
        log_success "Built and pushed $component: $image_name"
    done
}

# Deploy infrastructure
deploy_infrastructure() {
    log_info "Deploying infrastructure with Terraform..."
    
    cd "$PROJECT_ROOT/terraform"
    
    # Initialize Terraform
    terraform init -upgrade
    
    # Plan deployment
    terraform plan \
        -var="environment=$ENVIRONMENT" \
        -var="aws_region=${AWS_REGION}" \
        -out=tfplan
    
    # Apply deployment
    terraform apply tfplan
    
    # Update kubeconfig
    local cluster_name="ai-browser-$ENVIRONMENT"
    aws eks update-kubeconfig \
        --region "$AWS_REGION" \
        --name "$cluster_name"
    
    log_success "Infrastructure deployed"
}

# Deploy application
deploy_application() {
    log_info "Deploying application with Helm..."
    
    local image_tag="${GITHUB_SHA:-$(git rev-parse HEAD)}"
    local values_file="$PROJECT_ROOT/deployments/helm/ai-agentic-browser/values-$ENVIRONMENT.yaml"
    
    # Use default values if environment-specific file doesn't exist
    if [ ! -f "$values_file" ]; then
        values_file="$PROJECT_ROOT/deployments/helm/ai-agentic-browser/values.yaml"
    fi
    
    # Add Helm repositories
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
    helm repo update
    
    # Deploy with Helm
    helm upgrade --install "$HELM_RELEASE" \
        "$PROJECT_ROOT/deployments/helm/ai-agentic-browser" \
        --namespace "$NAMESPACE" \
        --values "$values_file" \
        --set global.environment="$ENVIRONMENT" \
        --set image.tag="$image_tag" \
        --set image.registry="$DOCKER_REGISTRY" \
        --wait \
        --timeout="$KUBECTL_TIMEOUT"
    
    log_success "Application deployed"
}

# Post-deployment verification
post_deployment_verification() {
    log_info "Running post-deployment verification..."
    
    # Wait for pods to be ready
    log_info "Waiting for pods to be ready..."
    kubectl wait --for=condition=ready pod \
        -l app.kubernetes.io/instance="$HELM_RELEASE" \
        -n "$NAMESPACE" \
        --timeout="$KUBECTL_TIMEOUT"
    
    # Check service endpoints
    log_info "Checking service endpoints..."
    local services
    services=$(kubectl get services -n "$NAMESPACE" -o json | jq -r '.items[].metadata.name')
    
    for service in $services; do
        local endpoint
        endpoint=$(kubectl get endpoints "$service" -n "$NAMESPACE" -o json | jq -r '.subsets[0].addresses[0].ip // "none"')
        if [ "$endpoint" = "none" ]; then
            log_warning "Service $service has no endpoints"
        else
            log_success "Service $service is ready"
        fi
    done
    
    # Run health checks
    log_info "Running health checks..."
    local api_gateway_url
    api_gateway_url=$(kubectl get service api-gateway -n "$NAMESPACE" -o json | jq -r '.status.loadBalancer.ingress[0].hostname // "localhost"')
    
    # Wait for load balancer to be ready
    sleep 30
    
    # Health check with retry
    local max_retries=10
    local retry_count=0
    
    while [ $retry_count -lt $max_retries ]; do
        if curl -f -s "http://$api_gateway_url/health" > /dev/null; then
            log_success "Health check passed"
            break
        else
            log_info "Health check failed, retrying... ($((retry_count + 1))/$max_retries)"
            sleep 10
            ((retry_count++))
        fi
    done
    
    if [ $retry_count -eq $max_retries ]; then
        log_error "Health check failed after $max_retries attempts"
        exit 1
    fi
    
    log_success "Post-deployment verification completed"
}

# Generate deployment report
generate_deployment_report() {
    log_info "Generating deployment report..."
    
    local report_file="deployment-report-$(date +%Y%m%d-%H%M%S).json"
    
    cat > "$report_file" << EOF
{
  "deployment": {
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "environment": "$ENVIRONMENT",
    "namespace": "$NAMESPACE",
    "helm_release": "$HELM_RELEASE",
    "image_tag": "${GITHUB_SHA:-$(git rev-parse HEAD)}",
    "deployed_by": "${USER:-unknown}",
    "git_commit": "$(git rev-parse HEAD)",
    "git_branch": "$(git rev-parse --abbrev-ref HEAD)"
  },
  "services": $(kubectl get services -n "$NAMESPACE" -o json | jq '.items | map({name: .metadata.name, type: .spec.type, ports: .spec.ports})'),
  "pods": $(kubectl get pods -n "$NAMESPACE" -o json | jq '.items | map({name: .metadata.name, status: .status.phase, ready: .status.conditions | map(select(.type == "Ready")) | .[0].status})')
}
EOF
    
    log_success "Deployment report generated: $report_file"
}

# Main deployment function
main() {
    log_info "Starting production deployment for AI Agentic Browser"
    log_info "Environment: $ENVIRONMENT"
    log_info "Namespace: $NAMESPACE"
    log_info "Helm Release: $HELM_RELEASE"
    
    validate_prerequisites
    validate_environment
    pre_deployment_checks
    
    if [ "${SKIP_BUILD:-false}" != "true" ]; then
        build_and_push_images
    fi
    
    if [ "${SKIP_INFRASTRUCTURE:-false}" != "true" ]; then
        deploy_infrastructure
    fi
    
    deploy_application
    post_deployment_verification
    generate_deployment_report
    
    log_success "Production deployment completed successfully!"
    log_info "Access the application at: https://ai-browser-$ENVIRONMENT.com"
    log_info "Monitoring dashboard: https://grafana-$ENVIRONMENT.ai-browser.com"
}

# Script execution
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
