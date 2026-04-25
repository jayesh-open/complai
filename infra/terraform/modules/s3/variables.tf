variable "environment" {
  description = "Deployment environment (dev, sandbox, staging, prod)"
  type        = string
}

variable "project" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "complai"
}

variable "kms_key_arn" {
  description = "ARN of the KMS key for S3 server-side encryption"
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC for the S3 Gateway VPC Endpoint"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-south-1"
}

variable "route_table_ids" {
  description = "List of route table IDs to associate with the S3 VPC Endpoint"
  type        = list(string)
}
