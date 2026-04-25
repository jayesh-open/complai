output "domain_endpoint" {
  description = "Endpoint URL for the OpenSearch domain"
  value       = aws_opensearch_domain.main.endpoint
}

output "domain_arn" {
  description = "ARN of the OpenSearch domain"
  value       = aws_opensearch_domain.main.arn
}

output "security_group_id" {
  description = "ID of the OpenSearch security group"
  value       = aws_security_group.opensearch.id
}
