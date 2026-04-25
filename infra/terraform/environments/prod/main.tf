# ---------------------------------------------------------------
# Complai — Production Environment
# ---------------------------------------------------------------
# Full spec from architecture doc: Multi-AZ everything,
# production-grade instances, full node group fleet, read
# replicas, cluster-mode Redis, 3-node OpenSearch.
# ---------------------------------------------------------------

terraform {
  required_version = ">= 1.7.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }

  # backend "s3" {
  #   bucket         = "complai-terraform-state"
  #   key            = "prod/terraform.tfstate"
  #   region         = "ap-south-1"
  #   dynamodb_table = "complai-terraform-locks"
  #   encrypt        = true
  # }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "complai"
      ManagedBy   = "terraform"
    }
  }
}

# ---------------------------------------------------------------
# Random password for OpenSearch
# ---------------------------------------------------------------

resource "random_password" "opensearch_master" {
  length  = 32
  special = true
}

# ---------------------------------------------------------------
# KMS
# ---------------------------------------------------------------

module "kms" {
  source = "../../modules/kms"

  environment = var.environment
  project     = "complai"
}

# ---------------------------------------------------------------
# VPC
# ---------------------------------------------------------------

module "vpc" {
  source = "../../modules/vpc"

  vpc_cidr    = var.vpc_cidr
  environment = var.environment
  project     = "complai"
  azs         = ["${var.region}a", "${var.region}b", "${var.region}c"]
}

# ---------------------------------------------------------------
# EKS
# ---------------------------------------------------------------

module "eks" {
  source = "../../modules/eks"

  cluster_name       = "complai-${var.environment}"
  cluster_version    = "1.30"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  environment        = var.environment

  node_groups = {
    system = {
      instance_types = ["m7i.large"]
      desired_size   = 2
      min_size       = 2
      max_size       = 4
      labels = {
        "workload-type" = "system"
      }
    }
    application = {
      instance_types = ["m7i.xlarge"]
      desired_size   = 5
      min_size       = 5
      max_size       = 30
      labels = {
        "workload-type" = "application"
      }
    }
    batch = {
      instance_types = ["c7i.2xlarge"]
      desired_size   = 0
      min_size       = 0
      max_size       = 10
      capacity_type  = "SPOT"
      labels = {
        "workload-type" = "batch"
      }
      taints = [
        {
          key    = "workload-type"
          value  = "batch"
          effect = "NO_SCHEDULE"
        }
      ]
    }
  }
}

# ---------------------------------------------------------------
# RDS PostgreSQL
# ---------------------------------------------------------------

module "rds" {
  source = "../../modules/rds-postgres"

  identifier     = "complai-${var.environment}"
  engine_version = "16"
  instance_class = "db.r7g.2xlarge"
  environment    = var.environment

  storage = {
    allocated     = 500
    max_allocated = 2000
    iops          = 12000
    throughput    = 500
  }

  multi_az                = true
  backup_retention_period = 35
  deletion_protection     = true
  create_replica          = true

  vpc_id                  = module.vpc.vpc_id
  data_subnet_ids         = module.vpc.data_subnet_ids
  allowed_security_groups = [module.eks.cluster_security_group_id]
  kms_key_arn             = module.kms.platform_key_arn
}

# ---------------------------------------------------------------
# ElastiCache Redis
# ---------------------------------------------------------------

module "redis" {
  source = "../../modules/elasticache-redis"

  cluster_id         = "complai-${var.environment}"
  node_type          = "cache.r7g.large"
  num_shards         = 3
  replicas_per_shard = 2
  environment        = var.environment

  vpc_id                  = module.vpc.vpc_id
  data_subnet_ids         = module.vpc.data_subnet_ids
  allowed_security_groups = [module.eks.cluster_security_group_id]
  kms_key_arn             = module.kms.platform_key_arn
}

# ---------------------------------------------------------------
# OpenSearch
# ---------------------------------------------------------------

module "opensearch" {
  source = "../../modules/opensearch"

  domain_name          = "complai-${var.environment}"
  instance_type        = "m7g.large.search"
  instance_count       = 3
  volume_size          = 200
  volume_iops          = 3000
  volume_throughput    = 125
  environment          = var.environment
  master_user_password = random_password.opensearch_master.result

  vpc_id                  = module.vpc.vpc_id
  data_subnet_ids         = module.vpc.data_subnet_ids
  allowed_security_groups = [module.eks.cluster_security_group_id]
  kms_key_arn             = module.kms.platform_key_arn
}

# ---------------------------------------------------------------
# S3 Buckets
# ---------------------------------------------------------------

module "s3" {
  source = "../../modules/s3"

  environment     = var.environment
  kms_key_arn     = module.kms.platform_key_arn
  vpc_id          = module.vpc.vpc_id
  route_table_ids = [] # Populated after VPC module provides route table IDs
}

# ---------------------------------------------------------------
# SQS + SNS
# ---------------------------------------------------------------

module "sqs_sns" {
  source = "../../modules/sqs-sns"

  environment = var.environment
  kms_key_arn = module.kms.platform_key_arn
}

# ---------------------------------------------------------------
# Secrets Manager
# ---------------------------------------------------------------

module "secrets" {
  source = "../../modules/secrets-manager"

  environment = var.environment
  kms_key_arn = module.kms.platform_key_arn
}

# ---------------------------------------------------------------
# ACM Certificate
# ---------------------------------------------------------------

module "acm" {
  source = "../../modules/acm"

  domain_name = var.domain_name
  zone_id     = var.cloudflare_zone_id
  environment = var.environment
}

# ---------------------------------------------------------------
# ALB
# ---------------------------------------------------------------

module "alb" {
  source = "../../modules/alb"

  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  certificate_arn   = module.acm.certificate_arn
  environment       = var.environment
}

# ---------------------------------------------------------------
# Cloudflare DNS
# ---------------------------------------------------------------

module "cloudflare_dns" {
  source = "../../modules/cloudflare-dns"

  zone_id      = var.cloudflare_zone_id
  domain       = var.domain_name
  alb_dns_name = module.alb.alb_dns_name
  environment  = var.environment
}

# ---------------------------------------------------------------
# IAM IRSA Roles
# ---------------------------------------------------------------

module "iam_irsa" {
  source = "../../modules/iam-irsa"

  oidc_provider_arn = module.eks.oidc_provider_arn
  environment       = var.environment
  kms_key_arns      = [module.kms.platform_key_arn]

  service_names = [
    "identity-service",
    "tenant-service",
    "user-role-service",
    "master-data-service",
    "document-service",
    "notification-service",
    "audit-service",
    "workflow-service",
    "rules-engine-service",
    "gst-service",
    "gstr9-service",
    "einvoice-service",
    "ewb-service",
    "tds-service",
    "itr-service",
    "vendor-service",
    "recon-service",
    "ap-service",
    "billing-service",
    "secretarial-service",
    "gstn-gateway",
    "irp-gateway",
    "ewb-gateway",
    "tds-gateway",
    "itd-gateway",
    "kyc-gateway",
    "tax-payment-gateway",
    "bank-gateway",
    "mca-gateway",
    "erp-gateway",
    "web-bff-service",
    "portal-bff-service",
    "smb-bff-service",
    "reporting-service",
    "matching-ml-service",
    "llm-copilot-service",
    "ocr-service",
  ]
}
