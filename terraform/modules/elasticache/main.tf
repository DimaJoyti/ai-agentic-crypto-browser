# ElastiCache Module for AI Agentic Crypto Browser

terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
}

# Security Group for ElastiCache
resource "aws_security_group" "elasticache" {
  name_prefix = "${var.project_name}-elasticache-sg"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = var.allowed_security_groups
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-elasticache-sg"
  })
}

# ElastiCache Parameter Group
resource "aws_elasticache_parameter_group" "redis" {
  family = "redis7.x"
  name   = "${var.project_name}-redis-params"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }

  parameter {
    name  = "timeout"
    value = "300"
  }

  parameter {
    name  = "tcp-keepalive"
    value = "300"
  }

  tags = var.tags
}

# ElastiCache Replication Group (Redis Cluster)
resource "aws_elasticache_replication_group" "redis" {
  replication_group_id       = "${var.project_name}-redis"
  description                = "Redis cluster for ${var.project_name}"

  # Node configuration
  node_type               = var.node_type
  port                    = 6379
  parameter_group_name    = aws_elasticache_parameter_group.redis.name

  # Cluster configuration
  num_cache_clusters      = var.num_cache_clusters
  
  # Engine configuration
  engine_version          = var.redis_version
  
  # Network configuration
  subnet_group_name       = var.subnet_group_name
  security_group_ids      = [aws_security_group.elasticache.id]

  # Backup configuration
  snapshot_retention_limit = var.snapshot_retention_limit
  snapshot_window         = var.snapshot_window
  
  # Maintenance
  maintenance_window      = var.maintenance_window
  auto_minor_version_upgrade = true

  # Security
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  auth_token                 = var.auth_token_enabled ? random_password.auth_token[0].result : null

  # Logging
  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.redis_slow.name
    destination_type = "cloudwatch-logs"
    log_format       = "text"
    log_type         = "slow-log"
  }

  tags = var.tags

  lifecycle {
    ignore_changes = [auth_token]
  }
}

# Random auth token for Redis
resource "random_password" "auth_token" {
  count   = var.auth_token_enabled ? 1 : 0
  length  = 32
  special = false
}

# CloudWatch Log Group for Redis slow logs
resource "aws_cloudwatch_log_group" "redis_slow" {
  name              = "/aws/elasticache/redis/${var.project_name}/slow-log"
  retention_in_days = 7

  tags = var.tags
}

# Store Redis auth token in AWS Secrets Manager (if enabled)
resource "aws_secretsmanager_secret" "redis_auth_token" {
  count                   = var.auth_token_enabled ? 1 : 0
  name                    = "${var.project_name}/elasticache/redis/auth-token"
  description             = "Redis auth token"
  recovery_window_in_days = 7

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "redis_auth_token" {
  count     = var.auth_token_enabled ? 1 : 0
  secret_id = aws_secretsmanager_secret.redis_auth_token[0].id
  secret_string = jsonencode({
    auth_token = random_password.auth_token[0].result
    endpoint   = aws_elasticache_replication_group.redis.configuration_endpoint_address != "" ? aws_elasticache_replication_group.redis.configuration_endpoint_address : aws_elasticache_replication_group.redis.primary_endpoint_address
    port       = aws_elasticache_replication_group.redis.port
  })
}

# CloudWatch Alarms for monitoring
resource "aws_cloudwatch_metric_alarm" "redis_cpu" {
  alarm_name          = "${var.project_name}-redis-cpu-utilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ElastiCache"
  period              = "300"
  statistic           = "Average"
  threshold           = "80"
  alarm_description   = "This metric monitors redis cpu utilization"
  alarm_actions       = var.alarm_actions

  dimensions = {
    CacheClusterId = aws_elasticache_replication_group.redis.id
  }

  tags = var.tags
}

resource "aws_cloudwatch_metric_alarm" "redis_memory" {
  alarm_name          = "${var.project_name}-redis-memory-utilization"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseMemoryUsagePercentage"
  namespace           = "AWS/ElastiCache"
  period              = "300"
  statistic           = "Average"
  threshold           = "80"
  alarm_description   = "This metric monitors redis memory utilization"
  alarm_actions       = var.alarm_actions

  dimensions = {
    CacheClusterId = aws_elasticache_replication_group.redis.id
  }

  tags = var.tags
}
