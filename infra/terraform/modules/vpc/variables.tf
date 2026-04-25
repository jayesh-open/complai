variable "vpc_cidr" {
  description = "CIDR block for the VPC (e.g. 10.0.0.0/16)"
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

variable "azs" {
  description = "List of availability zones to deploy into"
  type        = list(string)
}
