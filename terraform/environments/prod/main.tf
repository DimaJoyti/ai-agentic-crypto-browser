# Production Environment for AI Agentic Crypto Browser

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }

  # Uncomment and configure for remote state
  # backend "s3" {
  #   bucket = "your-terraform-state-bucket"
  #   key    = "ai-crypto-browser/prod/terraform.tfstate"
  #   region = "us-east-1"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# Local values
locals {
  cluster_name = "${var.project_name}-${var.environment}"
  
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# VPC Module
module "vpc" {
  source = "../../modules/vpc"

  project_name = var.project_name
  cluster_name = local.cluster_name
  vpc_cidr     = var.vpc_cidr

  public_subnet_cidrs   = var.public_subnet_cidrs
  private_subnet_cidrs  = var.private_subnet_cidrs
  database_subnet_cidrs = var.database_subnet_cidrs

  # High availability for production
  enable_nat_gateway     = true
  single_nat_gateway     = false
  one_nat_gateway_per_az = true

  tags = local.common_tags
}

# EKS Module
module "eks" {
  source = "../../modules/eks"

  cluster_name       = local.cluster_name
  vpc_id            = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  public_subnet_ids  = module.vpc.public_subnet_ids

  kubernetes_version                    = var.kubernetes_version
  cluster_endpoint_public_access_cidrs = var.cluster_endpoint_public_access_cidrs

  node_group_capacity_type   = var.node_group_capacity_type
  node_group_instance_types  = var.node_group_instance_types
  node_group_ami_type       = var.node_group_ami_type
  node_group_disk_size      = var.node_group_disk_size
  node_group_desired_size   = var.node_group_desired_size
  node_group_max_size       = var.node_group_max_size
  node_group_min_size       = var.node_group_min_size

  tags = local.common_tags
}

# RDS Module
module "rds" {
  source = "../../modules/rds"

  project_name         = var.project_name
  vpc_id              = module.vpc.vpc_id
  db_subnet_group_name = module.vpc.database_subnet_group_name

  allowed_security_groups = [
    module.eks.node_group_security_group_id,
    module.eks.cluster_security_group_id
  ]

  postgres_version        = var.postgres_version
  instance_class         = var.rds_instance_class
  allocated_storage      = var.rds_allocated_storage
  max_allocated_storage  = var.rds_max_allocated_storage
  database_name          = var.database_name
  database_username      = var.database_username
  backup_retention_period = var.backup_retention_period
  deletion_protection    = var.rds_deletion_protection

  tags = local.common_tags
}

# ElastiCache Module
module "elasticache" {
  source = "../../modules/elasticache"

  project_name        = var.project_name
  vpc_id             = module.vpc.vpc_id
  subnet_group_name  = module.vpc.elasticache_subnet_group_name

  allowed_security_groups = [
    module.eks.node_group_security_group_id,
    module.eks.cluster_security_group_id
  ]

  redis_version           = var.redis_version
  node_type              = var.elasticache_node_type
  num_cache_clusters     = var.elasticache_num_cache_clusters
  snapshot_retention_limit = var.elasticache_snapshot_retention_limit
  auth_token_enabled     = var.elasticache_auth_token_enabled

  tags = local.common_tags
}

# ECR Repositories for container images
resource "aws_ecr_repository" "api_gateway" {
  name                 = "${var.project_name}/api-gateway"
  image_tag_mutability = "IMMUTABLE"  # Production should use immutable tags

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = local.common_tags
}

resource "aws_ecr_repository" "auth_service" {
  name                 = "${var.project_name}/auth-service"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = local.common_tags
}

resource "aws_ecr_repository" "browser_service" {
  name                 = "${var.project_name}/browser-service"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = local.common_tags
}

resource "aws_ecr_repository" "web3_service" {
  name                 = "${var.project_name}/web3-service"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = local.common_tags
}

resource "aws_ecr_repository" "frontend" {
  name                 = "${var.project_name}/frontend"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = local.common_tags
}

# ECR Lifecycle Policies (more aggressive for production)
resource "aws_ecr_lifecycle_policy" "api_gateway" {
  repository = aws_ecr_repository.api_gateway.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 20 production images"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["v", "prod"]
          countType     = "imageCountMoreThan"
          countNumber   = 20
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Delete untagged images after 1 day"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 1
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# WAF for production
resource "aws_wafv2_web_acl" "main" {
  name  = "${var.project_name}-waf"
  scope = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "RateLimitRule"
    priority = 1

    override_action {
      none {}
    }

    statement {
      rate_based_statement {
        limit              = 2000
        aggregate_key_type = "IP"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "RateLimitRule"
      sampled_requests_enabled   = true
    }

    action {
      block {}
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${var.project_name}-waf"
    sampled_requests_enabled   = true
  }

  tags = local.common_tags
}
