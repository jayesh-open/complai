# ---------------------------------------------------------------
# Amazon OpenSearch Service 2 — Complai Search
# ---------------------------------------------------------------
# VPC deployment in data subnets. Used for audit log search,
# vendor search, and full-text queries across compliance data.
# ---------------------------------------------------------------

# ---------------------------------------------------------------
# Security Group
# ---------------------------------------------------------------

resource "aws_security_group" "opensearch" {
  name_prefix = "${var.domain_name}-os-"
  description = "Security group for Complai OpenSearch domain"
  vpc_id      = var.vpc_id

  tags = {
    Name        = "${var.domain_name}-os-sg"
    Environment = var.environment
    Project     = var.project
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "opensearch_ingress" {
  count = length(var.allowed_security_groups)

  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  source_security_group_id = var.allowed_security_groups[count.index]
  security_group_id        = aws_security_group.opensearch.id
  description              = "Allow HTTPS from application security group"
}

# ---------------------------------------------------------------
# Service-Linked Role (create if not exists)
# ---------------------------------------------------------------

data "aws_iam_policy_document" "opensearch_access" {
  statement {
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions   = ["es:*"]
    resources = ["arn:aws:es:${var.region}:${data.aws_caller_identity.current.account_id}:domain/${var.domain_name}/*"]
  }
}

data "aws_caller_identity" "current" {}

# ---------------------------------------------------------------
# OpenSearch Domain
# ---------------------------------------------------------------

resource "aws_opensearch_domain" "main" {
  domain_name    = var.domain_name
  engine_version = "OpenSearch_${var.engine_version}"

  cluster_config {
    instance_type          = var.instance_type
    instance_count         = var.instance_count
    zone_awareness_enabled = var.instance_count > 1

    dynamic "zone_awareness_config" {
      for_each = var.instance_count > 1 ? [1] : []
      content {
        availability_zone_count = min(var.instance_count, 3)
      }
    }
  }

  ebs_options {
    ebs_enabled = true
    volume_type = "gp3"
    volume_size = var.volume_size
    iops        = var.volume_iops
    throughput  = var.volume_throughput
  }

  vpc_options {
    subnet_ids         = var.instance_count > 1 ? slice(var.data_subnet_ids, 0, min(var.instance_count, 3)) : [var.data_subnet_ids[0]]
    security_group_ids = [aws_security_group.opensearch.id]
  }

  encrypt_at_rest {
    enabled    = true
    kms_key_id = var.kms_key_arn
  }

  node_to_node_encryption {
    enabled = true
  }

  domain_endpoint_options {
    enforce_https       = true
    tls_security_policy = "Policy-Min-TLS-1-2-PF-2023-10"
  }

  advanced_security_options {
    enabled                        = true
    internal_user_database_enabled = true

    master_user_options {
      master_user_name     = var.master_user_name
      master_user_password = var.master_user_password
    }
  }

  access_policies = data.aws_iam_policy_document.opensearch_access.json

  tags = {
    Name        = var.domain_name
    Environment = var.environment
    Project     = var.project
  }
}
