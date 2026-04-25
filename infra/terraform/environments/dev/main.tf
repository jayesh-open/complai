# ---------------------------------------------------------------
# Complai — Dev Environment
# ---------------------------------------------------------------
# Minimal footprint: smaller instances, single-AZ where possible,
# lower node counts. Sufficient for development and integration
# testing but not for load testing or production traffic.
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
  #   key            = "dev/terraform.tfstate"
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
  length  = 24
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
      instance_types = ["t3.large"]
      desired_size   = 1
      min_size       = 1
      max_size       = 2
      labels = {
        "workload-type" = "system"
      }
    }
    application = {
      instance_types = ["t3.xlarge"]
      desired_size   = 2
      min_size       = 1
      max_size       = 5
      labels = {
        "workload-type" = "application"
      }
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
  instance_class = "db.t4g.medium"
  environment    = var.environment

  storage = {
    allocated     = 100
    max_allocated = 200
    iops          = 3000
    throughput    = 125
  }

  multi_az                = false
  backup_retention_period = 7
  deletion_protection     = false
  create_replica          = false

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
  node_type          = "cache.t4g.micro"
  num_shards         = 1
  replicas_per_shard = 0
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
  instance_type        = "t3.small.search"
  instance_count       = 1
  volume_size          = 50
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
  domain       = "dev.${var.domain_name}"
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
    "einvoice-service",
    "ewb-service",
    "tds-service",
    "itr-service",
    "vendor-service",
    "recon-service",
    "ap-service",
    "billing-service",
    "gstn-gateway",
    "irp-gateway",
    "ewb-gateway",
    "tds-gateway",
    "itd-gateway",
    "kyc-gateway",
    "tax-payment-gateway",
    "web-bff-service",
    "portal-bff-service",
    "smb-bff-service",
    "reporting-service",
  ]
}
