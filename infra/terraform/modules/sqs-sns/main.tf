# ---------------------------------------------------------------
# SQS Queues + SNS Topics — Complai Messaging Layer
# ---------------------------------------------------------------
# Outbox-driven message bus: SQS for work queues, SNS for
# fan-out events. Every queue has a DLQ with 5-retry redrive.
# ---------------------------------------------------------------

locals {
  # Standard queues (non-FIFO)
  standard_queues = toset([
    "gov-outbound-irp",
    "gov-outbound-ewb",
    "gov-outbound-tds",
    "gov-outbound-itd",
    "gov-outbound-kyc",
    "notification-jobs",
    "ocr-jobs",
  ])

  # FIFO queues (ordering required)
  fifo_queues = toset([
    "gov-outbound-gstn",
  ])

  # SNS topics for domain events
  sns_topics = toset([
    "FilingCompleted",
    "InvoiceCreated",
    "VendorCreated",
    "MasterDataChanged",
  ])
}

# ---------------------------------------------------------------
# Standard Queues + DLQs
# ---------------------------------------------------------------

resource "aws_sqs_queue" "standard_dlq" {
  for_each = local.standard_queues

  name                       = "${var.project}-${var.environment}-${each.value}-dlq"
  message_retention_seconds  = 1209600 # 14 days
  kms_master_key_id          = var.kms_key_arn
  kms_data_key_reuse_period_seconds = 300

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}-dlq"
    Environment = var.environment
    Project     = var.project
    QueueType   = "dlq"
  }
}

resource "aws_sqs_queue" "standard" {
  for_each = local.standard_queues

  name                       = "${var.project}-${var.environment}-${each.value}"
  visibility_timeout_seconds = 300
  message_retention_seconds  = 345600 # 4 days
  receive_wait_time_seconds  = 20     # long-polling
  kms_master_key_id          = var.kms_key_arn
  kms_data_key_reuse_period_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.standard_dlq[each.value].arn
    maxReceiveCount     = 5
  })

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}"
    Environment = var.environment
    Project     = var.project
    QueueType   = "standard"
  }
}

# ---------------------------------------------------------------
# FIFO Queues + DLQs
# ---------------------------------------------------------------

resource "aws_sqs_queue" "fifo_dlq" {
  for_each = local.fifo_queues

  name                        = "${var.project}-${var.environment}-${each.value}-dlq.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  message_retention_seconds   = 1209600 # 14 days
  kms_master_key_id           = var.kms_key_arn
  kms_data_key_reuse_period_seconds = 300

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}-dlq"
    Environment = var.environment
    Project     = var.project
    QueueType   = "fifo-dlq"
  }
}

resource "aws_sqs_queue" "fifo" {
  for_each = local.fifo_queues

  name                        = "${var.project}-${var.environment}-${each.value}.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
  visibility_timeout_seconds  = 300
  message_retention_seconds   = 345600 # 4 days
  receive_wait_time_seconds   = 20     # long-polling
  kms_master_key_id           = var.kms_key_arn
  kms_data_key_reuse_period_seconds = 300

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.fifo_dlq[each.value].arn
    maxReceiveCount     = 5
  })

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}"
    Environment = var.environment
    Project     = var.project
    QueueType   = "fifo"
  }
}

# ---------------------------------------------------------------
# SNS Topics
# ---------------------------------------------------------------

resource "aws_sns_topic" "topics" {
  for_each = local.sns_topics

  name              = "${var.project}-${var.environment}-${each.value}"
  kms_master_key_id = var.kms_key_arn

  tags = {
    Name        = "${var.project}-${var.environment}-${each.value}"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# SNS → SQS Subscription Policies
# ---------------------------------------------------------------
# Allow SNS topics to send messages to SQS queues.
# Specific subscriptions are created per-environment as needed.
# This policy allows any SNS topic in the account to write to
# the standard queues (scoped by condition).
# ---------------------------------------------------------------

data "aws_caller_identity" "current" {}

resource "aws_sqs_queue_policy" "standard_allow_sns" {
  for_each = local.standard_queues

  queue_url = aws_sqs_queue.standard[each.value].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowSNSPublish"
        Effect    = "Allow"
        Principal = { Service = "sns.amazonaws.com" }
        Action    = "sqs:SendMessage"
        Resource  = aws_sqs_queue.standard[each.value].arn
        Condition = {
          ArnLike = {
            "aws:SourceArn" = "arn:aws:sns:*:${data.aws_caller_identity.current.account_id}:${var.project}-${var.environment}-*"
          }
        }
      }
    ]
  })
}
