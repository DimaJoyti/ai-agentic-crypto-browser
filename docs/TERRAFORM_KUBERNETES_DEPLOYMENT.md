# Terraform & Kubernetes Deployment Guide

This guide covers deploying the AI Agentic Crypto Browser using Terraform for infrastructure and Kubernetes for container orchestration.

## Architecture Overview

The deployment consists of:

### Infrastructure (Terraform)
- **VPC**: Multi-AZ setup with public, private, and database subnets
- **EKS Cluster**: Managed Kubernetes cluster with auto-scaling node groups
- **RDS PostgreSQL**: Managed database with encryption and backups
- **ElastiCache Redis**: Managed Redis cluster for caching and sessions
- **ECR Repositories**: Container image repositories for each service
- **Security Groups**: Network security with least privilege access
- **IAM Roles**: Service-specific permissions following AWS best practices

### Application (Kubernetes)
- **API Gateway**: Load balancer and routing service
- **Auth Service**: Authentication and authorization
- **Browser Service**: Web scraping and AI-powered browsing
- **Web3 Service**: Blockchain and cryptocurrency interactions
- **Frontend**: Next.js web application
- **Auto-scaling**: Horizontal Pod Autoscaler for all services
- **Monitoring**: Prometheus metrics and health checks

## Prerequisites

### Required Tools
```bash
# Install Terraform
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Install eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Install Helm
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update && sudo apt-get install helm

# Install Kustomize
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
sudo mv kustomize /usr/local/bin/
```

### AWS Configuration
```bash
# Configure AWS credentials
aws configure

# Verify access
aws sts get-caller-identity
```

## Quick Start

### 1. Deploy Everything (Automated)
```bash
# Deploy to development environment
./scripts/deploy.sh dev

# Deploy to staging environment
./scripts/deploy.sh staging

# Deploy to production environment
./scripts/deploy.sh prod
```

### 2. Manual Step-by-Step Deployment

#### Step 1: Deploy Infrastructure
```bash
cd terraform/environments/dev
terraform init
terraform plan
terraform apply
```

#### Step 2: Configure kubectl
```bash
aws eks update-kubeconfig --region us-east-1 --name ai-crypto-browser-dev
kubectl cluster-info
```

#### Step 3: Build and Push Images
```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.us-east-1.amazonaws.com

# Build and push each service
docker build -f cmd/api-gateway/Dockerfile -t $(aws ecr describe-repositories --repository-names ai-crypto-browser/api-gateway --query 'repositories[0].repositoryUri' --output text):dev-latest .
docker push $(aws ecr describe-repositories --repository-names ai-crypto-browser/api-gateway --query 'repositories[0].repositoryUri' --output text):dev-latest

# Repeat for other services...
```

#### Step 4: Install Operators
```bash
# External Secrets Operator
kubectl create namespace external-secrets-system
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets-system

# AWS Load Balancer Controller
eksctl create iamserviceaccount \
    --cluster=ai-crypto-browser-dev \
    --namespace=kube-system \
    --name=aws-load-balancer-controller \
    --attach-policy-arn=arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess \
    --approve

helm repo add eks https://aws.github.io/eks-charts
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
    -n kube-system \
    --set clusterName=ai-crypto-browser-dev \
    --set serviceAccount.create=false \
    --set serviceAccount.name=aws-load-balancer-controller
```

#### Step 5: Deploy Application
```bash
kubectl apply -k k8s/overlays/dev
kubectl wait --for=condition=available --timeout=600s deployment --all -n ai-crypto-browser
```

## Configuration

### Environment Variables
Update the ConfigMap in `k8s/base/configmap.yaml` for application configuration.

### Secrets Management
Secrets are managed through AWS Secrets Manager and External Secrets Operator:

1. **Database credentials**: Automatically created by Terraform
2. **Redis auth token**: Automatically created by Terraform
3. **JWT secrets**: Manually created in `k8s/base/secrets.yaml`
4. **API keys**: Update base64 encoded values in `k8s/base/secrets.yaml`

### Scaling Configuration
Adjust scaling parameters in the HorizontalPodAutoscaler resources:
- `minReplicas`: Minimum number of pods
- `maxReplicas`: Maximum number of pods
- `targetCPUUtilizationPercentage`: CPU threshold for scaling

## Monitoring and Observability

### Health Checks
All services include:
- **Liveness probes**: Restart unhealthy containers
- **Readiness probes**: Remove unhealthy pods from load balancing

### Metrics
Prometheus metrics are exposed on port 9090 for all services.

### Logs
View logs using kubectl:
```bash
# View all pods
kubectl get pods -n ai-crypto-browser

# View logs for a specific service
kubectl logs -f deployment/api-gateway -n ai-crypto-browser

# View logs for all services
kubectl logs -f -l app.kubernetes.io/part-of=ai-crypto-browser -n ai-crypto-browser
```

## Troubleshooting

### Common Issues

#### 1. Image Pull Errors
```bash
# Check ECR repositories
aws ecr describe-repositories

# Verify image exists
aws ecr describe-images --repository-name ai-crypto-browser/api-gateway
```

#### 2. Pod Startup Issues
```bash
# Check pod status
kubectl describe pod <pod-name> -n ai-crypto-browser

# Check events
kubectl get events -n ai-crypto-browser --sort-by='.lastTimestamp'
```

#### 3. Service Discovery Issues
```bash
# Check services
kubectl get services -n ai-crypto-browser

# Test connectivity
kubectl run debug --image=busybox -it --rm -- nslookup auth-service.ai-crypto-browser.svc.cluster.local
```

#### 4. Database Connection Issues
```bash
# Check database secret
kubectl get secret database-secret -n ai-crypto-browser -o yaml

# Verify RDS instance
aws rds describe-db-instances --db-instance-identifier ai-crypto-browser-postgres
```

### Useful Commands

```bash
# Get all resources
kubectl get all -n ai-crypto-browser

# Check resource usage
kubectl top pods -n ai-crypto-browser
kubectl top nodes

# Scale a deployment
kubectl scale deployment api-gateway --replicas=3 -n ai-crypto-browser

# Update an image
kubectl set image deployment/api-gateway api-gateway=new-image:tag -n ai-crypto-browser

# Restart a deployment
kubectl rollout restart deployment/api-gateway -n ai-crypto-browser
```

## Security Considerations

### Network Security
- All services run in private subnets
- Security groups follow least privilege principle
- Database and Redis are not publicly accessible

### Container Security
- Non-root user execution
- Read-only root filesystem
- Dropped capabilities
- Resource limits enforced

### Secrets Management
- Secrets stored in AWS Secrets Manager
- Automatic rotation supported
- Encrypted at rest and in transit

## Cost Optimization

### Development Environment
- Single replica deployments
- Smaller instance types (t3.micro for RDS, cache.t3.micro for Redis)
- Single NAT Gateway for cost optimization (~$45/month savings)
- Reduced resource requests and limits

### Staging Environment
- Production-like setup for testing
- One NAT Gateway per AZ for high availability
- Medium instance types (t3.small for RDS, cache.t3.small for Redis)

### Production Environment
- Multi-AZ deployments for high availability
- One NAT Gateway per AZ for redundancy
- Larger instance types for better performance
- Auto-scaling based on demand
- Reserved instances for cost savings

## Cleanup

### Destroy Infrastructure
```bash
cd terraform/environments/dev
terraform destroy
```

### Delete Kubernetes Resources
```bash
kubectl delete -k k8s/overlays/dev
```

## Next Steps

1. **Set up CI/CD pipeline** for automated deployments
2. **Configure monitoring** with Prometheus and Grafana
3. **Set up log aggregation** with ELK stack or CloudWatch
4. **Implement backup strategies** for databases
5. **Configure SSL/TLS** with cert-manager and Let's Encrypt
