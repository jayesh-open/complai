# ---------------------------------------------------------------
# AWS KMS — Complai Platform Encryption Keys
# ---------------------------------------------------------------
# Platform CMK used for S3, SQS, SNS, Secrets Manager, RDS,
# ElastiCache, and OpenSearch encryption. Per-tenant CMKs are
# created dynamically by tenant-service at runtime.
# ---------------------------------------------------------------

data "aws_caller_identity" "current" {}

# ---------------------------------------------------------------
# Platform CMK
# ---------------------------------------------------------------

resource "aws_kms_key" "platform" {
  description             = "Complai ${var.environment} platform encryption key"
  deletion_window_in_days = 30
  enable_key_rotation     = true
  multi_region            = false

  policy = jsonencode({
    Version = "2012-10-17"
    Id      = "complai-platform-key-policy"
    Statement = [
      {
        Sid    = "EnableRootAccountAccess"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action   = "kms:*"
        Resource = "*"
      },
      {
        Sid    = "AllowKeyAdministration"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action = [
          "kms:Create*",
          "kms:Describe*",
          "kms:Enable*",
          "kms:List*",
          "kms:Put*",
          "kms:Update*",
          "kms:Revoke*",
          "kms:Disable*",
          "kms:Get*",
          "kms:Delete*",
          "kms:TagResource",
          "kms:UntagResource",
          "kms:ScheduleKeyDeletion",
          "kms:CancelKeyDeletion",
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowServiceUsage"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action = [
          "kms:Encrypt",
          "kms:Decrypt",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:DescribeKey",
        ]
        Resource = "*"
      },
      {
        Sid    = "AllowAWSServicesGrant"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action = [
          "kms:CreateGrant",
          "kms:ListGrants",
          "kms:RevokeGrant",
        ]
        Resource = "*"
        Condition = {
          Bool = {
            "kms:GrantIsForAWSResource" = "true"
          }
        }
      },
    ]
  })

  tags = {
    Name        = "${var.project}-${var.environment}-platform-cmk"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Platform CMK Alias
# ---------------------------------------------------------------

resource "aws_kms_alias" "platform" {
  name          = "alias/${var.project}-${var.environment}-platform"
  target_key_id = aws_kms_key.platform.key_id
}

# ---------------------------------------------------------------
# Per-Tenant CMKs — Created dynamically
# ---------------------------------------------------------------
# Per-tenant CMKs are NOT created in Terraform. They are
# provisioned at runtime by the tenant-service when a new
# tenant is onboarded. The tenant-service uses aws-sdk-go-v2
# to create a CMK scoped to the tenant's data.
#
# Key naming convention:
#   alias/complai-{environment}-tenant-{tenant_id}
#
# This ensures tenant data isolation at the encryption layer.
# ---------------------------------------------------------------
