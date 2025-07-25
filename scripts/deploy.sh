#!/bin/bash

# AI Agentic Crypto Browser - Deployment Script
# This script deploys the application to Kubernetes using Terraform and Kustomize

set -euo pipefail

# Configuration
ENVIRONMENT=${1:-dev}
AWS_REGION=${AWS_REGION:-us-east-1}
PROJECT_NAME="ai-crypto-browser"

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

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if required tools are installed
    local tools=("terraform" "kubectl" "aws" "docker" "kustomize")
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "$tool is not installed. Please install it first."
            exit 1
        fi
    done
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials not configured. Please run 'aws configure' first."
        exit 1
    fi
    
    log_success "All prerequisites met"
}

# Deploy infrastructure with Terraform
deploy_infrastructure() {
    log_info "Deploying infrastructure with Terraform..."
    
    cd "terraform/environments/$ENVIRONMENT"
    
    # Initialize Terraform
    terraform init
    
    # Plan the deployment
    terraform plan -out=tfplan
    
    # Apply the deployment
    terraform apply tfplan
    
    # Get outputs
    CLUSTER_NAME=$(terraform output -raw cluster_id)
    ECR_REPOSITORIES=$(terraform output -json ecr_repositories)
    
    log_success "Infrastructure deployed successfully"
    
    cd - > /dev/null
}

# Configure kubectl
configure_kubectl() {
    log_info "Configuring kubectl..."
    
    aws eks update-kubeconfig --region "$AWS_REGION" --name "$PROJECT_NAME-$ENVIRONMENT"
    
    # Verify connection
    if kubectl cluster-info &> /dev/null; then
        log_success "kubectl configured successfully"
    else
        log_error "Failed to configure kubectl"
        exit 1
    fi
}

# Build and push Docker images
build_and_push_images() {
    log_info "Building and pushing Docker images..."
    
    # Get ECR login token
    aws ecr get-login-password --region "$AWS_REGION" | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com"
    
    # Build and push each service
    local services=("api-gateway" "auth-service" "browser-service" "web3-service" "frontend")
    
    for service in "${services[@]}"; do
        log_info "Building $service..."
        
        # Get ECR repository URL
        ECR_REPO=$(aws ecr describe-repositories --repository-names "$PROJECT_NAME/$service" --region "$AWS_REGION" --query 'repositories[0].repositoryUri' --output text)
        
        # Build image
        if [ "$service" = "frontend" ]; then
            docker build -f "web/Dockerfile" -t "$ECR_REPO:$ENVIRONMENT-latest" web/
        else
            docker build -f "cmd/$service/Dockerfile" -t "$ECR_REPO:$ENVIRONMENT-latest" .
        fi
        
        # Push image
        docker push "$ECR_REPO:$ENVIRONMENT-latest"
        
        log_success "$service image built and pushed"
    done
}

# Install required Kubernetes operators
install_operators() {
    log_info "Installing required Kubernetes operators..."
    
    # Install External Secrets Operator
    if ! kubectl get namespace external-secrets-system &> /dev/null; then
        log_info "Installing External Secrets Operator..."
        kubectl create namespace external-secrets-system
        helm repo add external-secrets https://charts.external-secrets.io
        helm repo update
        helm install external-secrets external-secrets/external-secrets -n external-secrets-system
        
        # Wait for operator to be ready
        kubectl wait --for=condition=available --timeout=300s deployment/external-secrets -n external-secrets-system
        log_success "External Secrets Operator installed"
    else
        log_info "External Secrets Operator already installed"
    fi
    
    # Install AWS Load Balancer Controller
    if ! kubectl get deployment aws-load-balancer-controller -n kube-system &> /dev/null; then
        log_info "Installing AWS Load Balancer Controller..."
        
        # Create IAM role for AWS Load Balancer Controller
        eksctl create iamserviceaccount \
            --cluster="$PROJECT_NAME-$ENVIRONMENT" \
            --namespace=kube-system \
            --name=aws-load-balancer-controller \
            --role-name "AmazonEKSLoadBalancerControllerRole-$PROJECT_NAME-$ENVIRONMENT" \
            --attach-policy-arn=arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess \
            --approve
        
        # Install the controller
        helm repo add eks https://aws.github.io/eks-charts
        helm repo update
        helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
            -n kube-system \
            --set clusterName="$PROJECT_NAME-$ENVIRONMENT" \
            --set serviceAccount.create=false \
            --set serviceAccount.name=aws-load-balancer-controller
        
        log_success "AWS Load Balancer Controller installed"
    else
        log_info "AWS Load Balancer Controller already installed"
    fi
}

# Deploy application to Kubernetes
deploy_application() {
    log_info "Deploying application to Kubernetes..."
    
    # Apply the Kubernetes manifests using Kustomize
    kubectl apply -k "k8s/overlays/$ENVIRONMENT"
    
    # Wait for deployments to be ready
    log_info "Waiting for deployments to be ready..."
    kubectl wait --for=condition=available --timeout=600s deployment --all -n ai-crypto-browser
    
    log_success "Application deployed successfully"
}

# Get application URLs
get_application_urls() {
    log_info "Getting application URLs..."
    
    # Wait for load balancers to be ready
    sleep 30
    
    # Get API Gateway URL
    API_GATEWAY_URL=$(kubectl get service api-gateway -n ai-crypto-browser -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    if [ -n "$API_GATEWAY_URL" ]; then
        log_success "API Gateway URL: http://$API_GATEWAY_URL"
    fi
    
    # Get Frontend URL
    FRONTEND_URL=$(kubectl get service frontend -n ai-crypto-browser -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    if [ -n "$FRONTEND_URL" ]; then
        log_success "Frontend URL: http://$FRONTEND_URL"
    fi
}

# Main deployment function
main() {
    log_info "Starting deployment for environment: $ENVIRONMENT"
    
    check_prerequisites
    deploy_infrastructure
    configure_kubectl
    build_and_push_images
    install_operators
    deploy_application
    get_application_urls
    
    log_success "Deployment completed successfully!"
    log_info "You can monitor your application with:"
    log_info "  kubectl get pods -n ai-crypto-browser"
    log_info "  kubectl logs -f deployment/api-gateway -n ai-crypto-browser"
}

# Show usage
usage() {
    echo "Usage: $0 [environment]"
    echo "  environment: dev, staging, or prod (default: dev)"
    echo ""
    echo "Examples:"
    echo "  $0 dev      # Deploy to development environment"
    echo "  $0 staging  # Deploy to staging environment"
    echo "  $0 prod     # Deploy to production environment"
}

# Parse command line arguments
case "${1:-}" in
    -h|--help)
        usage
        exit 0
        ;;
    dev|staging|prod)
        main
        ;;
    "")
        main
        ;;
    *)
        log_error "Invalid environment: $1"
        usage
        exit 1
        ;;
esac
