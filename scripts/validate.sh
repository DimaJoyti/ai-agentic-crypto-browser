#!/bin/bash

# AI Agentic Crypto Browser - Validation Script
# This script validates the infrastructure setup and configuration

set -euo pipefail

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

# Check if required tools are installed
check_tools() {
    log_info "Checking required tools..."
    
    local tools=("terraform" "kubectl" "aws" "docker" "helm" "kustomize")
    local missing_tools=()
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        else
            log_success "$tool is installed"
        fi
    done
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing tools: ${missing_tools[*]}"
        log_info "Please install the missing tools before proceeding"
        return 1
    fi
    
    log_success "All required tools are installed"
}

# Check AWS credentials and permissions
check_aws() {
    log_info "Checking AWS configuration..."
    
    # Check if AWS credentials are configured
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials not configured"
        log_info "Please run 'aws configure' to set up your credentials"
        return 1
    fi
    
    local account_id=$(aws sts get-caller-identity --query Account --output text)
    local region=$(aws configure get region)
    
    log_success "AWS credentials configured"
    log_info "Account ID: $account_id"
    log_info "Region: $region"
    
    # Check basic permissions
    if aws iam get-user &> /dev/null || aws sts get-caller-identity &> /dev/null; then
        log_success "AWS permissions verified"
    else
        log_warning "Could not verify AWS permissions"
    fi
}

# Validate Terraform configuration
validate_terraform() {
    log_info "Validating Terraform configuration..."
    
    cd terraform/environments/dev
    
    # Initialize Terraform
    if terraform init -backend=false &> /dev/null; then
        log_success "Terraform initialization successful"
    else
        log_error "Terraform initialization failed"
        cd - > /dev/null
        return 1
    fi
    
    # Validate configuration
    if terraform validate &> /dev/null; then
        log_success "Terraform configuration is valid"
    else
        log_error "Terraform configuration validation failed"
        terraform validate
        cd - > /dev/null
        return 1
    fi
    
    cd - > /dev/null
}

# Validate Kubernetes manifests
validate_kubernetes() {
    log_info "Validating Kubernetes manifests..."
    
    # Check if kustomize can build the manifests
    if kustomize build k8s/overlays/dev > /dev/null; then
        log_success "Kubernetes manifests are valid"
    else
        log_error "Kubernetes manifest validation failed"
        return 1
    fi
    
    # Check for common issues
    local manifests=$(kustomize build k8s/overlays/dev)
    
    # Check for missing image tags
    if echo "$manifests" | grep -q "image:.*:latest"; then
        log_warning "Found 'latest' image tags - consider using specific versions in production"
    fi
    
    # Check for resource limits
    if ! echo "$manifests" | grep -q "resources:"; then
        log_warning "Some containers may not have resource limits defined"
    fi
}

# Check Docker setup
check_docker() {
    log_info "Checking Docker setup..."
    
    if docker info &> /dev/null; then
        log_success "Docker is running"
    else
        log_error "Docker is not running or not accessible"
        log_info "Please start Docker and ensure your user has permission to access it"
        return 1
    fi
    
    # Check if we can pull a test image
    if docker pull hello-world &> /dev/null; then
        log_success "Docker can pull images"
        docker rmi hello-world &> /dev/null
    else
        log_warning "Docker may have issues pulling images"
    fi
}

# Check file permissions
check_permissions() {
    log_info "Checking file permissions..."
    
    if [ -x "scripts/deploy.sh" ]; then
        log_success "Deploy script is executable"
    else
        log_warning "Deploy script is not executable - fixing..."
        chmod +x scripts/deploy.sh
        log_success "Deploy script permissions fixed"
    fi
}

# Validate configuration files
validate_configs() {
    log_info "Validating configuration files..."
    
    # Check if example tfvars exists
    if [ -f "terraform/environments/dev/terraform.tfvars.example" ]; then
        log_success "Example tfvars file exists"
        
        if [ ! -f "terraform/environments/dev/terraform.tfvars" ]; then
            log_warning "terraform.tfvars not found - you may want to copy from the example"
            log_info "Run: cp terraform/environments/dev/terraform.tfvars.example terraform/environments/dev/terraform.tfvars"
        fi
    fi
    
    # Check for sensitive data in configs
    if grep -r "CHANGE_ME\|YOUR_KEY\|REPLACE_THIS" k8s/ terraform/ 2>/dev/null; then
        log_warning "Found placeholder values that need to be replaced"
    fi
}

# Main validation function
main() {
    log_info "Starting infrastructure validation..."
    echo
    
    local failed=0
    
    check_tools || failed=1
    echo
    
    check_aws || failed=1
    echo
    
    check_docker || failed=1
    echo
    
    check_permissions || failed=1
    echo
    
    validate_terraform || failed=1
    echo
    
    validate_kubernetes || failed=1
    echo
    
    validate_configs || failed=1
    echo
    
    if [ $failed -eq 0 ]; then
        log_success "All validations passed! You're ready to deploy."
        echo
        log_info "Next steps:"
        log_info "1. Review and customize terraform/environments/dev/terraform.tfvars"
        log_info "2. Update API keys in k8s/base/secrets.yaml"
        log_info "3. Run: ./scripts/deploy.sh dev"
    else
        log_error "Some validations failed. Please fix the issues before deploying."
        exit 1
    fi
}

# Show usage
usage() {
    echo "Usage: $0"
    echo ""
    echo "This script validates the infrastructure setup and configuration."
    echo "It checks for required tools, AWS credentials, and validates"
    echo "Terraform and Kubernetes configurations."
}

# Parse command line arguments
case "${1:-}" in
    -h|--help)
        usage
        exit 0
        ;;
    *)
        main
        ;;
esac
