variable "identifier" {
  description = "Unique identifier for the RDS instance"
  type        = string
}

variable "engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16"
}

variable "instance_class" {
  description = "RDS instance class for the primary (e.g. db.r7g.2xlarge)"
  type        = string
}

variable "storage" {
  description = "Storage configuration for the RDS instance"
  type = object({
    allocated     = number
    max_allocated = number
    iops          = number
    throughput    = number
  })
  default = {
    allocated     = 500
    max_allocated = 1000
    iops          = 12000
    throughput    = 500
  }
}

variable "db_name" {
  description = "Name of the default database to create"
  type        = string
  default     = "complai"
}

variable "master_username" {
  description = "Master username for the RDS instance"
  type        = string
  default     = "complai_admin"
}

variable "multi_az" {
  description = "Enable Multi-AZ deployment for high availability"
  type        = bool
  default     = true
}

variable "backup_retention_period" {
  description = "Number of days to retain automated backups (max 35 for PITR)"
  type        = number
  default     = 35
}

variable "deletion_protection" {
  description = "Enable deletion protection on the RDS instance"
  type        = bool
  default     = true
}

variable "create_replica" {
  description = "Whether to create a read replica for analytics workloads"
  type        = bool
  default     = true
}

variable "replica_instance_class" {
  description = "Instance class for the read replica (defaults to primary instance class)"
  type        = string
  default     = ""
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "data_subnet_ids" {
  description = "List of data subnet IDs for the DB subnet group"
  type        = list(string)
}

variable "allowed_security_groups" {
  description = "List of security group IDs allowed to connect to the database"
  type        = list(string)
  default     = []
}

variable "kms_key_arn" {
  description = "ARN of the KMS key for storage encryption"
  type        = string
  default     = null
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
