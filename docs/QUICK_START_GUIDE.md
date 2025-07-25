# AI Agentic Crypto Browser - Quick Start Guide

This guide will get you up and running with the AI Agentic Crypto Browser in under 30 minutes.

## 🚀 Prerequisites

### System Requirements
- **OS**: Linux (Ubuntu/Debian), macOS, or Windows with WSL2
- **RAM**: Minimum 8GB, recommended 16GB+
- **Storage**: At least 20GB free space
- **Network**: Stable internet connection

### Required Accounts
- **AWS Account** with administrative permissions
- **API Keys** (optional but recommended):
  - Infura or Alchemy for Ethereum RPC
  - CoinGecko API key
  - CoinMarketCap API key

## 📋 Step-by-Step Setup

### Step 1: Clone and Setup
```bash
# Clone the repository
git clone https://github.com/your-org/ai-agentic-crypto-browser.git
cd ai-agentic-crypto-browser

# Run the interactive setup script
./scripts/setup.sh
```

The setup script will:
- Install required tools (Terraform, kubectl, AWS CLI, Docker, Helm, etc.)
- Configure AWS credentials
- Set up Terraform configuration files
- Guide you through secrets configuration

### Step 2: Configure Your Environment

#### AWS Credentials
```bash
# Option 1: Use AWS CLI
aws configure

# Option 2: Set environment variables
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
```

#### Update Configuration
```bash
# Copy and customize Terraform variables
cp terraform/environments/dev/terraform.tfvars.example terraform/environments/dev/terraform.tfvars

# Edit the file to match your preferences
nano terraform/environments/dev/terraform.tfvars
```

#### Update Secrets (Optional)
```bash
# Update API keys in Kubernetes secrets
nano k8s/base/secrets.yaml

# Encode your API keys to base64
echo -n 'your-api-key' | base64
```

### Step 3: Validate Setup
```bash
# Run validation to check everything is configured correctly
./scripts/validate.sh
```

This will check:
- ✅ Required tools are installed
- ✅ AWS credentials are configured
- ✅ Terraform configuration is valid
- ✅ Kubernetes manifests are valid
- ✅ Docker is working

### Step 4: Deploy Infrastructure
```bash
# Deploy to development environment
./scripts/deploy.sh dev
```

This single command will:
1. **Deploy AWS Infrastructure** (5-10 minutes)
   - VPC with public/private subnets
   - EKS cluster with auto-scaling nodes
   - RDS PostgreSQL database
   - ElastiCache Redis cluster
   - ECR repositories

2. **Build and Push Images** (3-5 minutes)
   - Build Docker images for all services
   - Push to ECR repositories

3. **Deploy Application** (2-3 minutes)
   - Deploy all microservices to Kubernetes
   - Configure auto-scaling and monitoring
   - Set up load balancers

### Step 5: Access Your Application
```bash
# Get application URLs
kubectl get services -n ai-crypto-browser

# Check deployment status
kubectl get pods -n ai-crypto-browser

# View logs
kubectl logs -f deployment/api-gateway -n ai-crypto-browser
```

## 🔧 Configuration Options

### Environment Variables
Key settings in `k8s/base/configmap.yaml`:

```yaml
# AI Model Configuration
OLLAMA_MODEL: "qwen3"
OLLAMA_HOST: "http://ollama-service:11434"

# Database Configuration
DB_MAX_CONNECTIONS: "25"
DB_CONNECTION_MAX_LIFETIME: "300s"

# Security Configuration
RATE_LIMIT_REQUESTS: "100"
CORS_ALLOWED_ORIGINS: "*"
```

### Scaling Configuration
Adjust in the HorizontalPodAutoscaler resources:

```yaml
spec:
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## 🔍 Monitoring and Troubleshooting

### Check Application Health
```bash
# View all resources
kubectl get all -n ai-crypto-browser

# Check pod status
kubectl describe pod <pod-name> -n ai-crypto-browser

# View logs
kubectl logs <pod-name> -n ai-crypto-browser

# Check events
kubectl get events -n ai-crypto-browser --sort-by='.lastTimestamp'
```

### Common Issues

#### 1. Pod Startup Issues
```bash
# Check pod details
kubectl describe pod <pod-name> -n ai-crypto-browser

# Common causes:
# - Image pull errors (check ECR permissions)
# - Resource limits (check node capacity)
# - Configuration errors (check ConfigMaps/Secrets)
```

#### 2. Database Connection Issues
```bash
# Check database secret
kubectl get secret database-secret -n ai-crypto-browser -o yaml

# Verify RDS instance
aws rds describe-db-instances --db-instance-identifier ai-crypto-browser-postgres
```

#### 3. Service Discovery Issues
```bash
# Test internal connectivity
kubectl run debug --image=busybox -it --rm -- nslookup auth-service.ai-crypto-browser.svc.cluster.local
```

## 🧹 Cleanup

### Remove Everything
```bash
# Delete Kubernetes resources
kubectl delete -k k8s/overlays/dev

# Destroy AWS infrastructure
cd terraform/environments/dev
terraform destroy
```

## 📚 Next Steps

### Development
- **Local Development**: Set up local development environment
- **CI/CD Pipeline**: Configure automated deployments
- **Custom Models**: Add your own AI models

### Production
- **Security Hardening**: Implement additional security measures
- **Monitoring**: Set up Prometheus and Grafana
- **Backup Strategy**: Configure database backups
- **SSL/TLS**: Set up HTTPS with cert-manager

### Scaling
- **Multi-Region**: Deploy across multiple AWS regions
- **Performance Tuning**: Optimize resource allocation
- **Cost Optimization**: Implement cost monitoring and optimization

## 🆘 Getting Help

### Documentation
- [Full Deployment Guide](TERRAFORM_KUBERNETES_DEPLOYMENT.md)
- [AI Provider Setup](AI_PROVIDERS.md)
- [Ollama/LM Studio Setup](OLLAMA_LMSTUDIO_SETUP.md)

### Support
- **GitHub Issues**: Report bugs and feature requests
- **Discussions**: Community support and questions
- **Documentation**: Comprehensive guides and examples

### Useful Commands
```bash
# Scale a service
kubectl scale deployment api-gateway --replicas=3 -n ai-crypto-browser

# Update an image
kubectl set image deployment/api-gateway api-gateway=new-image:tag -n ai-crypto-browser

# Restart a service
kubectl rollout restart deployment/api-gateway -n ai-crypto-browser

# Port forward for local access
kubectl port-forward service/api-gateway 8080:80 -n ai-crypto-browser
```

## 🎉 Success!

You now have a fully functional AI Agentic Crypto Browser running on AWS with:
- ✅ Scalable microservices architecture
- ✅ AI-powered browsing with qwen3 model
- ✅ Web3 and cryptocurrency integration
- ✅ Production-ready infrastructure
- ✅ Monitoring and observability
- ✅ Auto-scaling capabilities

Happy browsing! 🚀
