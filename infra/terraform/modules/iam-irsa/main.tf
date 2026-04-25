# ---------------------------------------------------------------
# IAM Roles for Service Accounts (IRSA) — Complai
# ---------------------------------------------------------------
# One IAM role per Go/Node service, each scoped to its specific
# Secrets Manager secrets, S3 prefix, SQS queues, and KMS keys.
# ---------------------------------------------------------------

data "aws_caller_identity" "current" {}

locals {
  # Extract the OIDC issuer URL without https:// for trust policy
  oidc_issuer = replace(var.oidc_provider_arn, "/.*oidc-provider\\//", "")
}

# ---------------------------------------------------------------
# IAM Role per Service
# ---------------------------------------------------------------

resource "aws_iam_role" "service" {
  for_each = toset(var.service_names)

  name = "${var.project}-${var.environment}-${each.value}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = var.oidc_provider_arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${local.oidc_issuer}:sub" = "system:serviceaccount:${var.environment}:${each.value}"
            "${local.oidc_issuer}:aud" = "sts.amazonaws.com"
          }
        }
      }
    ]
  })

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}-role"
    Environment = var.environment
    Project     = var.project
    Service     = each.value
  }
}

# ---------------------------------------------------------------
# Secrets Manager Access — scoped to service-specific secrets
# ---------------------------------------------------------------

resource "aws_iam_role_policy" "secrets_access" {
  for_each = toset(var.service_names)

  name = "${each.value}-secrets-access"
  role = aws_iam_role.service[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ReadSecrets"
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret",
        ]
        Resource = "arn:aws:secretsmanager:*:${data.aws_caller_identity.current.account_id}:secret:${var.environment}/complai/${each.value}/*"
      },
      {
        Sid    = "ReadSharedSecrets"
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret",
        ]
        Resource = "arn:aws:secretsmanager:*:${data.aws_caller_identity.current.account_id}:secret:${var.environment}/complai/shared/*"
      },
    ]
  })
}

# ---------------------------------------------------------------
# S3 Access — scoped to service-specific prefix
# ---------------------------------------------------------------

resource "aws_iam_role_policy" "s3_access" {
  for_each = toset(var.service_names)

  name = "${each.value}-s3-access"
  role = aws_iam_role.service[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ListBucket"
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
        ]
        Resource = "arn:aws:s3:::${var.project}-${var.environment}-*"
        Condition = {
          StringLike = {
            "s3:prefix" = ["${each.value}/*"]
          }
        }
      },
      {
        Sid    = "ReadWriteObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
        ]
        Resource = "arn:aws:s3:::${var.project}-${var.environment}-*/${each.value}/*"
      },
    ]
  })
}

# ---------------------------------------------------------------
# SQS Access — scoped to service-related queues
# ---------------------------------------------------------------

resource "aws_iam_role_policy" "sqs_access" {
  for_each = toset(var.service_names)

  name = "${each.value}-sqs-access"
  role = aws_iam_role.service[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "SQSAccess"
        Effect = "Allow"
        Action = [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:GetQueueUrl",
          "sqs:ChangeMessageVisibility",
        ]
        Resource = "arn:aws:sqs:*:${data.aws_caller_identity.current.account_id}:${var.project}-${var.environment}-*"
      },
    ]
  })
}

# ---------------------------------------------------------------
# KMS Access — decrypt using platform key
# ---------------------------------------------------------------

resource "aws_iam_role_policy" "kms_access" {
  for_each = toset(var.service_names)

  name = "${each.value}-kms-access"
  role = aws_iam_role.service[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "KMSDecrypt"
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:GenerateDataKey",
          "kms:DescribeKey",
        ]
        Resource = var.kms_key_arns
      },
    ]
  })
}

# ---------------------------------------------------------------
# SNS Publish — for event publishing
# ---------------------------------------------------------------

resource "aws_iam_role_policy" "sns_access" {
  for_each = toset(var.service_names)

  name = "${each.value}-sns-access"
  role = aws_iam_role.service[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "SNSPublish"
        Effect = "Allow"
        Action = [
          "sns:Publish",
        ]
        Resource = "arn:aws:sns:*:${data.aws_caller_identity.current.account_id}:${var.project}-${var.environment}-*"
      },
    ]
  })
}
