# CI/CD Setup Guide

This document describes the complete CI/CD pipeline for the AI Agentic Crypto Browser project.

## 🏗️ **Pipeline Overview**

The CI/CD pipeline consists of three main workflows:

1. **Test Workflow** (`test.yml`) - Runs on every push and PR
2. **Infrastructure Workflow** (`infrastructure.yml`) - Manages Terraform and Kubernetes
3. **Deploy Workflow** (`deploy.yml`) - Handles application deployments

## 🔄 **Workflows**

### 1. Test Workflow (`test.yml`)

**Triggers:**
- Every push to any branch
- Every pull request

**Jobs:**
- **Code Quality**: Linting, formatting, security scans
- **Unit Tests**: Go tests across multiple versions (1.22, 1.23)
- **Integration Tests**: Database and Redis connectivity tests
- **E2E Tests**: Browser automation tests
- **Load Tests**: Performance testing
- **Security Tests**: Vulnerability scanning

**Artifacts:**
- Test results and coverage reports
- Security scan results
- Performance test results

### 2. Infrastructure Workflow (`infrastructure.yml`)

**Triggers:**
- Push to main (for Terraform/K8s changes)
- Pull requests (for validation)
- Manual dispatch (for specific environments)

**Jobs:**
- **Terraform Dev**: Automatic planning and applying for development
- **Terraform Staging**: Manual deployment to staging
- **Terraform Prod**: Manual deployment to production with approval
- **K8s Validation**: Validates Kubernetes manifests and Helm charts
- **Security Scan**: Infrastructure security scanning with Checkov

**Manual Deployment:**
```bash
# Trigger via GitHub UI or API
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/OWNER/REPO/actions/workflows/infrastructure.yml/dispatches \
  -d '{"ref":"main","inputs":{"environment":"staging","action":"apply"}}'
```

### 3. Deploy Workflow (`deploy.yml`)

**Triggers:**
- Successful completion of tests on main branch
- Manual dispatch for specific environments

**Jobs:**
- **Build & Push**: Docker images to ECR
- **Deploy Staging**: Automatic deployment to staging
- **Deploy Production**: Manual deployment to production
- **Rollback**: Emergency rollback capability

## 🔧 **Required Secrets**

### GitHub Repository Secrets

```bash
# AWS Credentials (Development/Staging)
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
AWS_REGION=us-east-1

# AWS Credentials (Production - separate account recommended)
AWS_ACCESS_KEY_ID_PROD=AKIA...
AWS_SECRET_ACCESS_KEY_PROD=...

# EKS Cluster Names
EKS_CLUSTER_NAME_DEV=ai-crypto-browser-dev
EKS_CLUSTER_NAME_STAGING=ai-crypto-browser-staging
EKS_CLUSTER_NAME_PROD=ai-crypto-browser-prod

# Container Registry
ECR_REGISTRY=123456789012.dkr.ecr.us-east-1.amazonaws.com

# Notifications (Optional)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
```

### Environment-Specific Secrets

Each environment (dev, staging, prod) should have its own GitHub Environment with protection rules:

**Development Environment:**
- No protection rules
- Auto-deployment on main branch

**Staging Environment:**
- Require reviewers: 1 person
- Manual deployment trigger

**Production Environment:**
- Require reviewers: 2 people
- Manual deployment trigger
- Deployment branches: main only

## 📋 **Setup Instructions**

### 1. Configure AWS Infrastructure

```bash
# 1. Set up AWS credentials with appropriate permissions
aws configure

# 2. Create S3 bucket for Terraform state (optional)
aws s3 mb s3://your-terraform-state-bucket

# 3. Update backend configuration in terraform files
# terraform/environments/*/main.tf
```

### 2. Configure GitHub Repository

```bash
# 1. Add repository secrets (see above)
# 2. Create GitHub Environments:
#    - Settings > Environments > New environment
#    - Configure protection rules for staging/prod

# 3. Enable GitHub Actions
#    - Settings > Actions > General > Allow all actions
```

### 3. Initial Infrastructure Deployment

```bash
# 1. Deploy development infrastructure manually first
cd terraform/environments/dev
terraform init
terraform plan
terraform apply

# 2. Configure kubectl
aws eks update-kubeconfig --region us-east-1 --name ai-crypto-browser-dev

# 3. Install required operators
kubectl apply -f https://github.com/external-secrets/external-secrets/releases/latest/download/bundle.yaml
```

### 4. Test the Pipeline

```bash
# 1. Create a test branch
git checkout -b test-ci-cd

# 2. Make a small change
echo "# Test" >> README.md

# 3. Push and create PR
git add .
git commit -m "Test CI/CD pipeline"
git push origin test-ci-cd

# 4. Create PR and verify all checks pass
```

## 🚀 **Deployment Process**

### Development Deployment
1. Push to main branch
2. Tests run automatically
3. If tests pass, infrastructure is updated (if changed)
4. Application is deployed automatically

### Staging Deployment
1. Trigger infrastructure workflow manually for staging
2. Trigger deploy workflow manually for staging
3. Verify deployment in staging environment

### Production Deployment
1. Ensure staging deployment is successful
2. Trigger infrastructure workflow manually for production
3. Wait for manual approval
4. Trigger deploy workflow manually for production
5. Monitor deployment and verify health

## 🔍 **Monitoring and Troubleshooting**

### Pipeline Monitoring
- **GitHub Actions**: Monitor workflow runs in GitHub UI
- **Slack Notifications**: Receive alerts for failures (if configured)
- **AWS CloudWatch**: Monitor infrastructure and application logs

### Common Issues

#### 1. Test Failures
```bash
# Check test logs in GitHub Actions
# Run tests locally:
go test ./... -v
docker-compose -f docker-compose.test.yml up -d
go test -tags=integration ./test/integration/...
```

#### 2. Infrastructure Deployment Failures
```bash
# Check Terraform logs in GitHub Actions
# Run locally:
cd terraform/environments/dev
terraform plan
terraform apply
```

#### 3. Application Deployment Failures
```bash
# Check Kubernetes logs:
kubectl get pods -n ai-crypto-browser
kubectl logs -f deployment/api-gateway -n ai-crypto-browser
kubectl describe pod <pod-name> -n ai-crypto-browser
```

#### 4. Image Build Failures
```bash
# Check Docker build logs in GitHub Actions
# Test locally:
docker build -f cmd/api-gateway/Dockerfile .
```

### Rollback Procedures

#### Application Rollback
```bash
# Via GitHub Actions (recommended)
# Trigger rollback workflow

# Manual rollback
helm rollback ai-agentic-browser -n ai-crypto-browser
```

#### Infrastructure Rollback
```bash
# Revert Terraform changes
cd terraform/environments/prod
git checkout HEAD~1 -- .
terraform plan
terraform apply
```

## 📊 **Pipeline Metrics**

### Performance Targets
- **Test Suite**: < 10 minutes
- **Build & Push**: < 15 minutes
- **Deployment**: < 20 minutes
- **Total Pipeline**: < 45 minutes

### Success Rates
- **Tests**: > 95% pass rate
- **Deployments**: > 98% success rate
- **Infrastructure**: > 99% success rate

## 🔒 **Security Considerations**

### Secrets Management
- Use GitHub Secrets for sensitive data
- Rotate credentials regularly
- Use least privilege access
- Separate production credentials

### Infrastructure Security
- Enable AWS CloudTrail
- Use IAM roles instead of access keys where possible
- Enable VPC Flow Logs
- Regular security scanning with Checkov

### Application Security
- Container image scanning
- Dependency vulnerability scanning
- SAST/DAST security testing
- Regular security updates

## 📚 **Additional Resources**

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [AWS EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
