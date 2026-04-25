output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_endpoint" {
  description = "EKS cluster API endpoint"
  value       = module.eks.cluster_endpoint
}

output "eks_cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "rds_primary_endpoint" {
  description = "RDS primary instance endpoint"
  value       = module.rds.primary_endpoint
}

output "rds_replica_endpoint" {
  description = "RDS read replica endpoint"
  value       = module.rds.replica_endpoint
}

output "redis_primary_endpoint" {
  description = "Redis cluster configuration endpoint"
  value       = module.redis.primary_endpoint
}

output "opensearch_endpoint" {
  description = "OpenSearch domain endpoint"
  value       = module.opensearch.domain_endpoint
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = module.alb.alb_dns_name
}

output "s3_bucket_names" {
  description = "Map of S3 bucket names"
  value       = module.s3.bucket_names
}

output "sqs_queue_urls" {
  description = "Map of SQS queue URLs"
  value       = module.sqs_sns.queue_urls
}

output "sns_topic_arns" {
  description = "Map of SNS topic ARNs"
  value       = module.sqs_sns.topic_arns
}

output "irsa_role_arns" {
  description = "Map of service IRSA role ARNs"
  value       = module.iam_irsa.role_arns
}

output "kms_platform_key_arn" {
  description = "Platform KMS CMK ARN"
  value       = module.kms.platform_key_arn
}
