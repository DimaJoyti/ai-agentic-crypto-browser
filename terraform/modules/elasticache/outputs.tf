output "redis_replication_group_id" {
  description = "ID of the ElastiCache replication group"
  value       = aws_elasticache_replication_group.redis.id
}

output "redis_replication_group_arn" {
  description = "ARN of the ElastiCache replication group"
  value       = aws_elasticache_replication_group.redis.arn
}

output "redis_primary_endpoint_address" {
  description = "Address of the endpoint for the primary node in the replication group"
  value       = aws_elasticache_replication_group.redis.primary_endpoint_address
}

output "redis_configuration_endpoint_address" {
  description = "Address of the replication group configuration endpoint when cluster mode is enabled"
  value       = aws_elasticache_replication_group.redis.configuration_endpoint_address
}

output "redis_port" {
  description = "Port number on which the cache nodes accept connections"
  value       = aws_elasticache_replication_group.redis.port
}

output "redis_auth_token" {
  description = "Auth token for Redis"
  value       = var.auth_token_enabled ? random_password.auth_token[0].result : null
  sensitive   = true
}

output "redis_security_group_id" {
  description = "Security group ID for ElastiCache"
  value       = aws_security_group.elasticache.id
}

output "redis_parameter_group_name" {
  description = "Name of the parameter group"
  value       = aws_elasticache_parameter_group.redis.name
}

output "redis_secret_arn" {
  description = "ARN of the Redis auth token secret"
  value       = var.auth_token_enabled ? aws_secretsmanager_secret.redis_auth_token[0].arn : null
}

output "redis_secret_name" {
  description = "Name of the Redis auth token secret"
  value       = var.auth_token_enabled ? aws_secretsmanager_secret.redis_auth_token[0].name : null
}
