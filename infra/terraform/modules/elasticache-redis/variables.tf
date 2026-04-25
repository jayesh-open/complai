variable "cluster_id" {
  description = "Identifier for the Redis replication group"
  type        = string
}

variable "engine_version" {
  description = "Redis engine version"
  type        = string
  default     = "7.1"
}

variable "node_type" {
  description = "ElastiCache node type (e.g. cache.r7g.large)"
  type        = string
}

variable "num_shards" {
  description = "Number of shards (node groups) in cluster-mode"
  type        = number
  default     = 3
}

variable "replicas_per_shard" {
  description = "Number of read replicas per shard"
  type        = number
  default     = 2
}

variable "vpc_id" {
  description = "ID of the VPC"
  type        = string
}

variable "data_subnet_ids" {
  description = "List of data subnet IDs for the cache subnet group"
  type        = list(string)
}

variable "allowed_security_groups" {
  description = "List of security group IDs allowed to connect to Redis"
  type        = list(string)
  default     = []
}

variable "kms_key_arn" {
  description = "ARN of the KMS key for at-rest encryption"
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
