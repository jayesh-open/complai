output "primary_endpoint" {
  description = "Connection endpoint for the primary RDS instance"
  value       = aws_db_instance.primary.endpoint
}

output "replica_endpoint" {
  description = "Connection endpoint for the read replica (empty if no replica)"
  value       = var.create_replica ? aws_db_instance.replica[0].endpoint : ""
}

output "port" {
  description = "Database port"
  value       = aws_db_instance.primary.port
}

output "db_name" {
  description = "Name of the default database"
  value       = aws_db_instance.primary.db_name
}

output "security_group_id" {
  description = "ID of the RDS security group"
  value       = aws_security_group.rds.id
}
