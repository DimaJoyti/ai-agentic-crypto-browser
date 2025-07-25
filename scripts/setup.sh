#!/bin/bash

# AI Agentic Crypto Browser - Setup Script
# This script helps set up the development environment

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

# Setup Terraform configuration
setup_terraform() {
    log_info "Setting up Terraform configuration..."
    
    local tfvars_file="terraform/environments/dev/terraform.tfvars"
    local example_file="terraform/environments/dev/terraform.tfvars.example"
    
    if [ ! -f "$tfvars_file" ]; then
        if [ -f "$example_file" ]; then
            cp "$example_file" "$tfvars_file"
            log_success "Created terraform.tfvars from example"
            log_warning "Please review and customize $tfvars_file"
        else
            log_error "Example tfvars file not found"
            return 1
        fi
    else
        log_info "terraform.tfvars already exists"
    fi
}

# Setup AWS credentials
setup_aws() {
    log_info "Checking AWS configuration..."
    
    if ! aws sts get-caller-identity &> /dev/null; then
        log_warning "AWS credentials not configured"
        log_info "Please configure AWS credentials:"
        echo
        echo "Option 1: Use AWS CLI"
        echo "  aws configure"
        echo
        echo "Option 2: Set environment variables"
        echo "  export AWS_ACCESS_KEY_ID=your_access_key"
        echo "  export AWS_SECRET_ACCESS_KEY=your_secret_key"
        echo "  export AWS_DEFAULT_REGION=us-east-1"
        echo
        read -p "Press Enter after configuring AWS credentials..."
        
        if aws sts get-caller-identity &> /dev/null; then
            log_success "AWS credentials configured successfully"
        else
            log_error "AWS credentials still not working"
            return 1
        fi
    else
        log_success "AWS credentials already configured"
    fi
}

# Setup secrets
setup_secrets() {
    log_info "Setting up secrets configuration..."
    
    local secrets_file="k8s/base/secrets.yaml"
    
    log_warning "Please update the following secrets in $secrets_file:"
    echo
    echo "1. JWT_SECRET - Generate a secure random string"
    echo "2. ETHEREUM_RPC_URL - Your Infura or Alchemy endpoint"
    echo "3. POLYGON_RPC_URL - Your Polygon RPC endpoint"
    echo "4. COINGECKO_API_KEY - Your CoinGecko API key"
    echo "5. COINMARKETCAP_API_KEY - Your CoinMarketCap API key"
    echo
    echo "All values should be base64 encoded. Example:"
    echo "  echo -n 'your-secret-value' | base64"
    echo
    
    read -p "Press Enter after updating the secrets..."
}

# Install required tools (Ubuntu/Debian)
install_tools_ubuntu() {
    log_info "Installing required tools for Ubuntu/Debian..."
    
    # Update package list
    sudo apt-get update
    
    # Install basic tools
    sudo apt-get install -y curl wget unzip
    
    # Install Terraform
    if ! command -v terraform &> /dev/null; then
        log_info "Installing Terraform..."
        curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
        sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
        sudo apt-get update && sudo apt-get install terraform
        log_success "Terraform installed"
    fi
    
    # Install kubectl
    if ! command -v kubectl &> /dev/null; then
        log_info "Installing kubectl..."
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
        rm kubectl
        log_success "kubectl installed"
    fi
    
    # Install AWS CLI
    if ! command -v aws &> /dev/null; then
        log_info "Installing AWS CLI..."
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
        unzip awscliv2.zip
        sudo ./aws/install
        rm -rf aws awscliv2.zip
        log_success "AWS CLI installed"
    fi
    
    # Install Docker
    if ! command -v docker &> /dev/null; then
        log_info "Installing Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker $USER
        rm get-docker.sh
        log_success "Docker installed"
        log_warning "Please log out and back in for Docker permissions to take effect"
    fi
    
    # Install Helm
    if ! command -v helm &> /dev/null; then
        log_info "Installing Helm..."
        curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
        sudo apt-get update && sudo apt-get install helm
        log_success "Helm installed"
    fi
    
    # Install Kustomize
    if ! command -v kustomize &> /dev/null; then
        log_info "Installing Kustomize..."
        curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
        sudo mv kustomize /usr/local/bin/
        log_success "Kustomize installed"
    fi
    
    # Install eksctl
    if ! command -v eksctl &> /dev/null; then
        log_info "Installing eksctl..."
        curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
        sudo mv /tmp/eksctl /usr/local/bin
        log_success "eksctl installed"
    fi
}

# Install required tools (macOS)
install_tools_macos() {
    log_info "Installing required tools for macOS..."
    
    # Check if Homebrew is installed
    if ! command -v brew &> /dev/null; then
        log_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
    
    # Install tools via Homebrew
    local tools=("terraform" "kubectl" "awscli" "docker" "helm" "kustomize" "eksctl")
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_info "Installing $tool..."
            brew install "$tool"
            log_success "$tool installed"
        fi
    done
}

# Detect OS and install tools
install_tools() {
    log_info "Detecting operating system..."
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            install_tools_ubuntu
        else
            log_warning "Unsupported Linux distribution. Please install tools manually."
            log_info "Required tools: terraform, kubectl, aws, docker, helm, kustomize, eksctl"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        install_tools_macos
    else
        log_warning "Unsupported operating system. Please install tools manually."
        log_info "Required tools: terraform, kubectl, aws, docker, helm, kustomize, eksctl"
    fi
}

# Main setup function
main() {
    log_info "Starting AI Agentic Crypto Browser setup..."
    echo
    
    # Ask user what they want to do
    echo "What would you like to set up?"
    echo "1. Install required tools"
    echo "2. Configure AWS credentials"
    echo "3. Set up Terraform configuration"
    echo "4. Set up secrets"
    echo "5. All of the above"
    echo "6. Validate setup"
    echo
    read -p "Enter your choice (1-6): " choice
    
    case $choice in
        1)
            install_tools
            ;;
        2)
            setup_aws
            ;;
        3)
            setup_terraform
            ;;
        4)
            setup_secrets
            ;;
        5)
            install_tools
            echo
            setup_aws
            echo
            setup_terraform
            echo
            setup_secrets
            ;;
        6)
            if [ -x "scripts/validate.sh" ]; then
                ./scripts/validate.sh
            else
                log_error "Validation script not found or not executable"
            fi
            ;;
        *)
            log_error "Invalid choice"
            exit 1
            ;;
    esac
    
    echo
    log_success "Setup completed!"
    log_info "Next steps:"
    log_info "1. Run './scripts/validate.sh' to check your setup"
    log_info "2. Run './scripts/deploy.sh dev' to deploy the infrastructure"
}

# Show usage
usage() {
    echo "Usage: $0"
    echo ""
    echo "This script helps set up the development environment for"
    echo "the AI Agentic Crypto Browser project."
    echo ""
    echo "It can install required tools, configure AWS credentials,"
    echo "and set up configuration files."
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
