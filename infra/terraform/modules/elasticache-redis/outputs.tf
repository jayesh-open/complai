output "primary_endpoint" {
  description = "Configuration endpoint for the Redis cluster (use for cluster-mode clients)"
  value       = aws_elasticache_replication_group.main.configuration_endpoint_address
}

output "reader_endpoint" {
  description = "Reader endpoint for Redis read replicas"
  value       = aws_elasticache_replication_group.main.reader_endpoint_address
}

output "port" {
  description = "Redis port"
  value       = aws_elasticache_replication_group.main.port
}

output "security_group_id" {
  description = "ID of the Redis security group"
  value       = aws_security_group.redis.id
}
