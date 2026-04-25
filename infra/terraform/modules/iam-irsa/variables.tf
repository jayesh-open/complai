variable "oidc_provider_arn" {
  description = "ARN of the EKS OIDC provider for IRSA trust"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, sandbox, staging, prod)"
  type        = string
}

variable "project" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "complai"
}

variable "service_names" {
  description = "List of service names to create IRSA roles for"
  type        = list(string)
}

variable "kms_key_arns" {
  description = "List of KMS key ARNs the services are allowed to use"
  type        = list(string)
  default     = []
}
