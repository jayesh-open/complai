variable "environment" {
  description = "Deployment environment (dev, sandbox, staging, prod)"
  type        = string
}

variable "project" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "complai"
}
