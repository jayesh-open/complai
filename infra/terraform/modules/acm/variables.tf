variable "domain_name" {
  description = "Base domain name for the ACM certificate (e.g. complai.in)"
  type        = string
}

variable "zone_id" {
  description = "Cloudflare zone ID for DNS validation records"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, sandbox, staging, prod)"
  type        = string
  default     = "prod"
}

variable "project" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "complai"
}
