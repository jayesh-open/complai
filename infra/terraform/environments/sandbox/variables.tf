variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "sandbox"
}

variable "region" {
  description = "AWS region for deployment"
  type        = string
  default     = "ap-south-1"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.3.0.0/16"
}

variable "domain_name" {
  description = "Base domain name for the application"
  type        = string
  default     = "complai.in"
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for DNS management"
  type        = string
}
