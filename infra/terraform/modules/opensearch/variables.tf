variable "domain_name" {
  description = "Name of the OpenSearch domain"
  type        = string
}

variable "engine_version" {
  description = "OpenSearch engine version (e.g. 2.11)"
  type        = string
  default     = "2.11"
}

variable "instance_type" {
  description = "OpenSearch instance type (e.g. m7g.large.search)"
  type        = string
}

variable "instance_count" {
  description = "Number of data nodes in the OpenSearch cluster"
  type        = number
  default     = 3
}

variable "volume_size" {
  description = "EBS volume size in GB per data node"
  type        = number
  default     = 200
}

variable "volume_iops" {
  description = "Provisioned IOPS for gp3 EBS volumes"
  type        = number
  default     = 3000
}

variable "volume_throughput" {
  description = "Provisioned throughput in MiB/s for gp3 EBS volumes"
  type        = number
  default     = 125
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "data_subnet_ids" {
  description = "List of data subnet IDs for VPC deployment"
  type        = list(string)
}

variable "allowed_security_groups" {
  description = "List of security group IDs allowed to connect to OpenSearch"
  type        = list(string)
  default     = []
}

variable "kms_key_arn" {
  description = "ARN of the KMS key for encryption at rest"
  type        = string
  default     = null
}

variable "master_user_name" {
  description = "Master user name for OpenSearch fine-grained access control"
  type        = string
  default     = "complai_admin"
}

variable "master_user_password" {
  description = "Master user password for OpenSearch (use Secrets Manager in practice)"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "AWS region for the OpenSearch domain"
  type        = string
  default     = "ap-south-1"
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
