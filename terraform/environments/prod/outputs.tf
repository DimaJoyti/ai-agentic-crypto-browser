# Outputs for Production Environment

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

# EKS Outputs
output "cluster_id" {
  description = "EKS cluster ID"
  value       = module.eks.cluster_id
}

output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = module.eks.cluster_arn
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster"
  value       = module.eks.cluster_security_group_id
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks.cluster_certificate_authority_data
  sensitive   = true
}

output "cluster_version" {
  description = "The Kubernetes version for the EKS cluster"
  value       = module.eks.cluster_version
}

# RDS Outputs
output "db_instance_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.db_instance_endpoint
}

output "db_instance_port" {
  description = "RDS instance port"
  value       = module.rds.db_instance_port
}

output "db_instance_name" {
  description = "RDS instance database name"
  value       = module.rds.db_instance_name
}

output "db_secret_arn" {
  description = "ARN of the database password secret"
  value       = module.rds.db_secret_arn
}

# ElastiCache Outputs
output "redis_primary_endpoint_address" {
  description = "Address of the endpoint for the primary node in the replication group"
  value       = module.elasticache.redis_primary_endpoint_address
}

output "redis_port" {
  description = "Port number on which the cache nodes accept connections"
  value       = module.elasticache.redis_port
}

output "redis_secret_arn" {
  description = "ARN of the Redis auth token secret"
  value       = module.elasticache.redis_secret_arn
}

# ECR Outputs
output "ecr_repositories" {
  description = "ECR repository URLs"
  value = {
    api_gateway     = aws_ecr_repository.api_gateway.repository_url
    auth_service    = aws_ecr_repository.auth_service.repository_url
    browser_service = aws_ecr_repository.browser_service.repository_url
    web3_service    = aws_ecr_repository.web3_service.repository_url
    frontend        = aws_ecr_repository.frontend.repository_url
  }
}

# WAF Outputs
output "waf_web_acl_arn" {
  description = "ARN of the WAF Web ACL"
  value       = aws_wafv2_web_acl.main.arn
}

# Kubectl Configuration Command
output "kubectl_config_command" {
  description = "Command to configure kubectl"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${local.cluster_name}"
}
