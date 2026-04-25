output "bucket_arns" {
  description = "Map of bucket purpose to ARN"
  value = {
    for k, v in aws_s3_bucket.buckets : k => v.arn
  }
}

output "bucket_names" {
  description = "Map of bucket purpose to name"
  value = {
    for k, v in aws_s3_bucket.buckets : k => v.id
  }
}

output "vpc_endpoint_id" {
  description = "ID of the S3 VPC Gateway Endpoint"
  value       = aws_vpc_endpoint.s3.id
}
