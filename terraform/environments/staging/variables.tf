# Variables for Staging Environment

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "ai-crypto-browser"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "staging"
}

# VPC Configuration
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.1.1.0/24", "10.1.2.0/24", "10.1.3.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.1.11.0/24", "10.1.12.0/24", "10.1.13.0/24"]
}

variable "database_subnet_cidrs" {
  description = "CIDR blocks for database subnets"
  type        = list(string)
  default     = ["10.1.21.0/24", "10.1.22.0/24", "10.1.23.0/24"]
}

# EKS Configuration
variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "cluster_endpoint_public_access_cidrs" {
  description = "List of CIDR blocks that can access the Amazon EKS public API server endpoint"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "node_group_capacity_type" {
  description = "Type of capacity associated with the EKS Node Group"
  type        = string
  default     = "ON_DEMAND"
}

variable "node_group_instance_types" {
  description = "List of instance types associated with the EKS Node Group"
  type        = list(string)
  default     = ["t3.large"]
}

variable "node_group_ami_type" {
  description = "Type of Amazon Machine Image (AMI) associated with the EKS Node Group"
  type        = string
  default     = "AL2_x86_64"
}

variable "node_group_disk_size" {
  description = "Disk size in GiB for worker nodes"
  type        = number
  default     = 30
}

variable "node_group_desired_size" {
  description = "Desired number of worker nodes"
  type        = number
  default     = 3
}

variable "node_group_max_size" {
  description = "Maximum number of worker nodes"
  type        = number
  default     = 6
}

variable "node_group_min_size" {
  description = "Minimum number of worker nodes"
  type        = number
  default     = 2
}

# RDS Configuration
variable "postgres_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "15.4"
}

variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.small"
}

variable "rds_allocated_storage" {
  description = "Initial allocated storage in GB"
  type        = number
  default     = 50
}

variable "rds_max_allocated_storage" {
  description = "Maximum allocated storage in GB for autoscaling"
  type        = number
  default     = 200
}

variable "database_name" {
  description = "Name of the database"
  type        = string
  default     = "ai_crypto_browser"
}

variable "database_username" {
  description = "Username for the database"
  type        = string
  default     = "postgres"
}

variable "backup_retention_period" {
  description = "Backup retention period in days"
  type        = number
  default     = 14
}

variable "rds_deletion_protection" {
  description = "Enable deletion protection for RDS"
  type        = bool
  default     = true
}

# ElastiCache Configuration
variable "redis_version" {
  description = "Redis version"
  type        = string
  default     = "7.0"
}

variable "elasticache_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.small"
}

variable "elasticache_num_cache_clusters" {
  description = "Number of cache clusters"
  type        = number
  default     = 2
}

variable "elasticache_snapshot_retention_limit" {
  description = "Number of days to retain snapshots"
  type        = number
  default     = 7
}

variable "elasticache_auth_token_enabled" {
  description = "Enable auth token for Redis"
  type        = bool
  default     = true
}
