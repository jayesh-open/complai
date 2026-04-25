# ---------------------------------------------------------------
# ElastiCache Redis 7 — Complai Cache Layer
# ---------------------------------------------------------------
# Cluster-mode enabled with configurable shards and replicas.
# Deployed in data subnets (no internet egress).
# ---------------------------------------------------------------

# ---------------------------------------------------------------
# Subnet Group
# ---------------------------------------------------------------

resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.cluster_id}-subnet-group"
  subnet_ids = var.data_subnet_ids

  tags = {
    Name        = "${var.cluster_id}-subnet-group"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Parameter Group
# ---------------------------------------------------------------

resource "aws_elasticache_parameter_group" "main" {
  name   = "${var.cluster_id}-params"
  family = "redis7"

  parameter {
    name  = "maxmemory-policy"
    value = "volatile-lru"
  }

  parameter {
    name  = "notify-keyspace-events"
    value = "Ex"
  }

  tags = {
    Name        = "${var.cluster_id}-params"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Security Group
# ---------------------------------------------------------------

resource "aws_security_group" "redis" {
  name_prefix = "${var.cluster_id}-redis-"
  description = "Security group for Complai ElastiCache Redis"
  vpc_id      = var.vpc_id

  tags = {
    Name        = "${var.cluster_id}-redis-sg"
    Environment = var.environment
    Project     = var.project
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "redis_ingress" {
  count = length(var.allowed_security_groups)

  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = var.allowed_security_groups[count.index]
  security_group_id        = aws_security_group.redis.id
  description              = "Allow Redis from application security group"
}

# ---------------------------------------------------------------
# Replication Group (cluster-mode enabled)
# ---------------------------------------------------------------

resource "aws_elasticache_replication_group" "main" {
  replication_group_id = var.cluster_id
  description          = "Complai Redis cluster for ${var.environment}"

  engine         = "redis"
  engine_version = var.engine_version
  node_type      = var.node_type
  port           = 6379

  num_node_groups         = var.num_shards
  replicas_per_node_group = var.replicas_per_shard

  automatic_failover_enabled = var.num_shards > 1 || var.replicas_per_shard > 0
  multi_az_enabled           = var.replicas_per_shard > 0

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [aws_security_group.redis.id]
  parameter_group_name = aws_elasticache_parameter_group.main.name

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  kms_key_id                 = var.kms_key_arn

  snapshot_retention_limit = var.environment == "prod" ? 7 : 1
  snapshot_window          = "03:00-05:00"
  maintenance_window       = "sun:05:00-sun:06:00"

  auto_minor_version_upgrade = true

  tags = {
    Name        = var.cluster_id
    Environment = var.environment
    Project     = var.project
  }
}
