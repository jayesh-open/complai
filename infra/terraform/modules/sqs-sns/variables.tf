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
  description = "ARN of the KMS key for SQS/SNS encryption"
  type        = string
}
