# AI Agentic Crypto Browser - Infrastructure as Code

This repository contains the complete infrastructure setup for the AI Agentic Crypto Browser using Terraform and Kubernetes.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                           AWS Cloud                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   Public Subnet │  │   Public Subnet │  │   Public Subnet │  │
│  │       AZ-a      │  │       AZ-b      │  │       AZ-c      │  │
│  │                 │  │                 │  │                 │  │
│  │  ┌─────────────┐│  │  ┌─────────────┐│  │  ┌─────────────┐│  │
│  │  │     NAT     ││  │  │     NAT     ││  │  │     NAT     ││  │
│  │  │   Gateway   ││  │  │   Gateway   ││  │  │   Gateway   ││  │
│  │  └─────────────┘│  │  └─────────────┘│  │  └─────────────┘│  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  Private Subnet │  │  Private Subnet │  │  Private Subnet │  │
│  │       AZ-a      │  │       AZ-b      │  │       AZ-c      │  │
│  │                 │  │                 │  │                 │  │
│  │  ┌─────────────┐│  │  ┌─────────────┐│  │  ┌─────────────┐│  │
│  │  │     EKS     ││  │  │     EKS     ││  │  │     EKS     ││  │
│  │  │    Nodes    ││  │  │    Nodes    ││  │  │    Nodes    ││  │
│  │  └─────────────┘│  │  └─────────────┘│  │  └─────────────┘│  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │Database Subnet  │  │Database Subnet  │  │Database Subnet  │  │
│  │       AZ-a      │  │       AZ-b      │  │       AZ-c      │  │
│  │                 │  │                 │  │                 │  │
│  │  ┌─────────────┐│  │  ┌─────────────┐│  │  ┌─────────────┐│  │
│  │  │     RDS     ││  │  │ElastiCache  ││  │  │             ││  │
│  │  │ PostgreSQL  ││  │  │    Redis    ││  │  │             ││  │
│  │  └─────────────┘│  │  └─────────────┘│  │  └─────────────┘│  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- AWS CLI configured with appropriate permissions
- Terraform >= 1.0
- kubectl
- Docker
- Helm
- eksctl

### Easy Setup (Recommended)
```bash
# Run the setup script to install tools and configure environment
./scripts/setup.sh

# Validate your setup
./scripts/validate.sh

# Deploy to development environment
./scripts/deploy.sh dev
```

### One-Command Deployment (if already set up)
```bash
./scripts/deploy.sh dev
```

This will:
1. Deploy AWS infrastructure with Terraform
2. Create EKS cluster and managed services
3. Build and push Docker images to ECR
4. Deploy application to Kubernetes
5. Configure auto-scaling and monitoring

## 📁 Directory Structure

```
├── terraform/                 # Infrastructure as Code
│   ├── modules/               # Reusable Terraform modules
│   │   ├── vpc/              # VPC, subnets, routing
│   │   ├── eks/              # EKS cluster and node groups
│   │   ├── rds/              # PostgreSQL database
│   │   └── elasticache/      # Redis cluster
│   └── environments/         # Environment-specific configs
│       ├── dev/              # Development environment
│       ├── staging/          # Staging environment
│       └── prod/             # Production environment
├── k8s/                      # Kubernetes manifests
│   ├── base/                 # Base Kubernetes resources
│   │   ├── services/         # Service deployments
│   │   ├── configmap.yaml    # Application configuration
│   │   ├── secrets.yaml      # Secret management
│   │   └── namespace.yaml    # Namespace definitions
│   ├── overlays/             # Environment-specific overlays
│   │   ├── dev/              # Development overrides
│   │   ├── staging/          # Staging overrides
│   │   └── prod/             # Production overrides
│   └── helm/                 # Helm charts
│       └── ai-crypto-browser/ # Main application chart
├── scripts/                  # Deployment and utility scripts
│   └── deploy.sh            # Main deployment script
└── docs/                    # Documentation
    └── TERRAFORM_KUBERNETES_DEPLOYMENT.md
```

## 🛠️ Infrastructure Components

### AWS Services Used

| Service | Purpose | Configuration |
|---------|---------|---------------|
| **VPC** | Network isolation | Multi-AZ with public/private/database subnets |
| **EKS** | Kubernetes cluster | Managed control plane with auto-scaling nodes |
| **RDS** | PostgreSQL database | Multi-AZ, encrypted, automated backups |
| **ElastiCache** | Redis cluster | Multi-node, encrypted, auth enabled |
| **ECR** | Container registry | Private repositories for each service |
| **ALB** | Load balancing | Application Load Balancer for ingress |
| **IAM** | Access control | Service-specific roles and policies |
| **Secrets Manager** | Secret storage | Encrypted secret management |
| **CloudWatch** | Monitoring | Logs and metrics collection |

### Kubernetes Resources

| Resource Type | Purpose | Count |
|---------------|---------|-------|
| **Deployments** | Application workloads | 5 (API Gateway, Auth, Browser, Web3, Frontend) |
| **Services** | Service discovery | 5 (one per deployment) |
| **ConfigMaps** | Configuration | 1 (shared application config) |
| **Secrets** | Sensitive data | 3 (DB, Redis, JWT/API keys) |
| **HPA** | Auto-scaling | 5 (one per deployment) |
| **ServiceAccounts** | Pod identity | 5 (one per service) |

## 🔧 Configuration

### Environment Variables
Key configuration options in `k8s/base/configmap.yaml`:

```yaml
# Application
APP_ENV: "development"
LOG_LEVEL: "info"

# Services
AUTH_SERVICE_URL: "http://auth-service:8081"
BROWSER_SERVICE_URL: "http://browser-service:8082"
WEB3_SERVICE_URL: "http://web3-service:8083"

# AI Configuration
OLLAMA_MODEL: "qwen3"
OLLAMA_HOST: "http://ollama-service:11434"

# Database
DB_MAX_CONNECTIONS: "25"
DB_CONNECTION_MAX_LIFETIME: "300s"

# Security
CORS_ALLOWED_ORIGINS: "*"
RATE_LIMIT_REQUESTS: "100"
```

### Secrets Management
Secrets are managed through AWS Secrets Manager and External Secrets Operator:

1. **Database credentials** - Auto-generated by Terraform
2. **Redis auth token** - Auto-generated by Terraform  
3. **JWT secrets** - Manually configured
4. **API keys** - Update with your actual keys

## 📊 Monitoring & Observability

### Health Checks
- **Liveness probes**: Restart unhealthy containers
- **Readiness probes**: Remove unhealthy pods from load balancing
- **Startup probes**: Handle slow-starting containers

### Metrics
- Prometheus metrics exposed on port 9090
- Custom application metrics
- Infrastructure metrics via CloudWatch

### Logging
- Structured JSON logging
- Centralized log collection
- Request tracing and correlation

## 🔒 Security Features

### Network Security
- Private subnets for application workloads
- Security groups with least privilege
- Network policies for pod-to-pod communication

### Container Security
- Non-root user execution
- Read-only root filesystem
- Dropped Linux capabilities
- Resource limits and requests

### Data Security
- Encryption at rest (RDS, ElastiCache, EBS)
- Encryption in transit (TLS/SSL)
- Secrets stored in AWS Secrets Manager
- IAM roles for service authentication

## 📈 Scaling & Performance

### Auto-scaling
- **Horizontal Pod Autoscaler**: Scale pods based on CPU/memory
- **Cluster Autoscaler**: Scale nodes based on pod requirements
- **Vertical Pod Autoscaler**: Adjust resource requests automatically

### Performance Optimization
- **Resource requests/limits**: Prevent resource contention
- **Affinity rules**: Optimize pod placement
- **Pod disruption budgets**: Maintain availability during updates

## 🚀 Deployment Strategies

### Blue-Green Deployment
```bash
# Deploy new version
kubectl set image deployment/api-gateway api-gateway=new-image:tag

# Monitor rollout
kubectl rollout status deployment/api-gateway

# Rollback if needed
kubectl rollout undo deployment/api-gateway
```

### Canary Deployment
```bash
# Scale up new version gradually
kubectl patch hpa api-gateway -p '{"spec":{"minReplicas":4}}'
```

## 🔧 Maintenance

### Regular Tasks
- **Update dependencies**: Keep Terraform providers and Kubernetes versions current
- **Rotate secrets**: Regular rotation of API keys and certificates
- **Backup verification**: Test database and configuration backups
- **Security scanning**: Regular vulnerability scans of images and infrastructure

### Monitoring Alerts
Set up alerts for:
- High CPU/memory usage
- Pod restart loops
- Database connection issues
- External API failures
- Security events

## 🆘 Troubleshooting

### Common Issues

#### Pod Startup Failures
```bash
kubectl describe pod <pod-name> -n ai-crypto-browser
kubectl logs <pod-name> -n ai-crypto-browser
```

#### Service Discovery Issues
```bash
kubectl get endpoints -n ai-crypto-browser
kubectl run debug --image=busybox -it --rm -- nslookup auth-service
```

#### Database Connection Problems
```bash
kubectl get secret database-secret -o yaml
aws rds describe-db-instances --db-instance-identifier ai-crypto-browser-postgres
```

### Useful Commands
```bash
# Get all resources
kubectl get all -n ai-crypto-browser

# Check resource usage
kubectl top pods -n ai-crypto-browser

# Scale deployment
kubectl scale deployment api-gateway --replicas=3

# Update image
kubectl set image deployment/api-gateway api-gateway=new-image:tag

# Restart deployment
kubectl rollout restart deployment/api-gateway
```

## 🧹 Cleanup

### Destroy Everything
```bash
# Delete Kubernetes resources
kubectl delete -k k8s/overlays/dev

# Destroy AWS infrastructure
cd terraform/environments/dev
terraform destroy
```

## 📚 Additional Resources

- [Terraform Documentation](https://www.terraform.io/docs)
- [Kubernetes Documentation](https://kubernetes.io/docs)
- [AWS EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Helm Documentation](https://helm.sh/docs/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
