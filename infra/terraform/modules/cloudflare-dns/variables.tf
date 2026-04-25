variable "zone_id" {
  description = "Cloudflare zone ID for the domain"
  type        = string
}

variable "domain" {
  description = "Domain name (e.g. complai.in)"
  type        = string
}

variable "alb_dns_name" {
  description = "DNS name of the ALB to point records to"
  type        = string
}

variable "environment" {
  description = "Deployment environment (dev, sandbox, staging, prod)"
  type        = string
}
