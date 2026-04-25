# ---------------------------------------------------------------
# RDS PostgreSQL 16 — Complai Primary Database
# ---------------------------------------------------------------
# Multi-AZ deployment with read replica for analytics workloads.
# All OLTP goes to primary; BI/reporting uses the replica.
# ---------------------------------------------------------------

# ---------------------------------------------------------------
# Subnet Group (data subnets — no internet egress)
# ---------------------------------------------------------------

resource "aws_db_subnet_group" "main" {
  name       = "${var.identifier}-subnet-group"
  subnet_ids = var.data_subnet_ids

  tags = {
    Name        = "${var.identifier}-subnet-group"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Parameter Group
# ---------------------------------------------------------------

resource "aws_db_parameter_group" "main" {
  name   = "${var.identifier}-params"
  family = "postgres16"

  parameter {
    name  = "shared_preload_libraries"
    value = "pg_stat_statements,pgaudit"
  }

  parameter {
    name  = "log_statement"
    value = "ddl"
  }

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }

  parameter {
    name  = "idle_in_transaction_session_timeout"
    value = "60000"
  }

  parameter {
    name         = "rds.force_ssl"
    value        = "1"
    apply_method = "pending-reboot"
  }

  tags = {
    Name        = "${var.identifier}-params"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Security Group
# ---------------------------------------------------------------

resource "aws_security_group" "rds" {
  name_prefix = "${var.identifier}-rds-"
  description = "Security group for Complai RDS PostgreSQL"
  vpc_id      = var.vpc_id

  tags = {
    Name        = "${var.identifier}-rds-sg"
    Environment = var.environment
    Project     = var.project
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "rds_ingress" {
  count = length(var.allowed_security_groups)

  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = var.allowed_security_groups[count.index]
  security_group_id        = aws_security_group.rds.id
  description              = "Allow PostgreSQL from application security group"
}

# ---------------------------------------------------------------
# Primary Instance
# ---------------------------------------------------------------

resource "aws_db_instance" "primary" {
  identifier = var.identifier

  engine         = "postgres"
  engine_version = var.engine_version
  instance_class = var.instance_class

  allocated_storage     = var.storage.allocated
  max_allocated_storage = var.storage.max_allocated
  storage_type          = "gp3"
  iops                  = var.storage.iops
  storage_throughput    = var.storage.throughput
  storage_encrypted     = true
  kms_key_id            = var.kms_key_arn

  db_name  = var.db_name
  username = var.master_username
  port     = 5432

  manage_master_user_password = true

  multi_az               = var.multi_az
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.main.name

  backup_retention_period   = var.backup_retention_period
  backup_window             = "03:00-04:00"
  maintenance_window        = "sun:04:30-sun:05:30"
  copy_tags_to_snapshot     = true
  deletion_protection       = var.deletion_protection
  skip_final_snapshot       = var.environment == "dev" ? true : false
  final_snapshot_identifier = var.environment == "dev" ? null : "${var.identifier}-final"

  performance_insights_enabled    = true
  monitoring_interval             = 60
  monitoring_role_arn             = aws_iam_role.rds_monitoring.arn
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  auto_minor_version_upgrade = true

  tags = {
    Name        = var.identifier
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Read Replica (analytics / BI workloads)
# ---------------------------------------------------------------

resource "aws_db_instance" "replica" {
  count = var.create_replica ? 1 : 0

  identifier = "${var.identifier}-replica"

  replicate_source_db = aws_db_instance.primary.identifier
  instance_class      = var.replica_instance_class != "" ? var.replica_instance_class : var.instance_class

  storage_encrypted = true
  kms_key_id        = var.kms_key_arn

  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.main.name

  performance_insights_enabled    = true
  monitoring_interval             = 60
  monitoring_role_arn             = aws_iam_role.rds_monitoring.arn
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  auto_minor_version_upgrade = true
  skip_final_snapshot        = true

  tags = {
    Name        = "${var.identifier}-replica"
    Environment = var.environment
    Project     = var.project
    Role        = "analytics-replica"
  }
}

# ---------------------------------------------------------------
# Enhanced Monitoring IAM Role
# ---------------------------------------------------------------

resource "aws_iam_role" "rds_monitoring" {
  name = "${var.identifier}-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.identifier}-rds-monitoring-role"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
  role       = aws_iam_role.rds_monitoring.name
}
