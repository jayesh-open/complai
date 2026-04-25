#!/bin/bash
# Complai — LocalStack initialisation script
# Creates SQS queues (with DLQs), SNS topics, S3 buckets, and test secrets.
# Executed automatically by LocalStack on startup via ready.d hook.

set -euo pipefail

REGION="ap-south-1"

echo "============================================================"
echo "  Complai LocalStack Init — Region: ${REGION}"
echo "============================================================"

# ---------------------------------------------------------------------------
# SQS — Dead-letter queues
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating SQS dead-letter queues..."

DLQ_QUEUES=(
  gov-outbound-gstn-dlq.fifo
  gov-outbound-irp-dlq
  gov-outbound-ewb-dlq
  gov-outbound-tds-dlq
  gov-outbound-itd-dlq
  gov-outbound-kyc-dlq
  notification-jobs-dlq
  ocr-jobs-dlq
)

for queue in "${DLQ_QUEUES[@]}"; do
  echo "  -> ${queue}"
  if [[ "${queue}" == *.fifo ]]; then
    awslocal sqs create-queue \
      --queue-name "${queue}" \
      --attributes FifoQueue=true,ContentBasedDeduplication=true \
      --region "${REGION}"
  else
    awslocal sqs create-queue \
      --queue-name "${queue}" \
      --region "${REGION}"
  fi
done

# ---------------------------------------------------------------------------
# SQS — Primary queues (with redrive policy pointing to DLQs)
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating SQS primary queues..."

declare -A QUEUE_MAP
QUEUE_MAP=(
  ["gov-outbound-gstn.fifo"]="gov-outbound-gstn-dlq.fifo"
  ["gov-outbound-irp"]="gov-outbound-irp-dlq"
  ["gov-outbound-ewb"]="gov-outbound-ewb-dlq"
  ["gov-outbound-tds"]="gov-outbound-tds-dlq"
  ["gov-outbound-itd"]="gov-outbound-itd-dlq"
  ["gov-outbound-kyc"]="gov-outbound-kyc-dlq"
  ["notification-jobs"]="notification-jobs-dlq"
  ["ocr-jobs"]="ocr-jobs-dlq"
)

for queue in "${!QUEUE_MAP[@]}"; do
  dlq="${QUEUE_MAP[$queue]}"
  dlq_arn=$(awslocal sqs get-queue-attributes \
    --queue-url "http://sqs.${REGION}.localhost.localstack.cloud:4566/000000000000/${dlq}" \
    --attribute-names QueueArn \
    --region "${REGION}" \
    --query 'Attributes.QueueArn' \
    --output text)

  echo "  -> ${queue} (DLQ: ${dlq})"

  if [[ "${queue}" == *.fifo ]]; then
    awslocal sqs create-queue \
      --queue-name "${queue}" \
      --attributes FifoQueue=true,ContentBasedDeduplication=true \
      --region "${REGION}"
  else
    awslocal sqs create-queue \
      --queue-name "${queue}" \
      --region "${REGION}"
  fi

  queue_url=$(awslocal sqs get-queue-url \
    --queue-name "${queue}" \
    --region "${REGION}" \
    --query 'QueueUrl' \
    --output text)

  awslocal sqs set-queue-attributes \
    --queue-url "${queue_url}" \
    --cli-input-json "{\"QueueUrl\": \"${queue_url}\", \"Attributes\": {\"RedrivePolicy\": \"{\\\"deadLetterTargetArn\\\": \\\"${dlq_arn}\\\", \\\"maxReceiveCount\\\": \\\"3\\\"}\"}}" \
    --region "${REGION}"
done

# ---------------------------------------------------------------------------
# SNS — Topics
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating SNS topics..."

SNS_TOPICS=(
  FilingCompleted
  InvoiceCreated
  VendorCreated
  MasterDataChanged
)

for topic in "${SNS_TOPICS[@]}"; do
  echo "  -> ${topic}"
  awslocal sns create-topic \
    --name "${topic}" \
    --region "${REGION}"
done

# ---------------------------------------------------------------------------
# S3 — Buckets
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating S3 buckets..."

S3_BUCKETS=(
  complai-dev-documents
  complai-dev-uploads
  complai-dev-exports
  complai-dev-backups
)

for bucket in "${S3_BUCKETS[@]}"; do
  echo "  -> ${bucket}"
  awslocal s3 mb "s3://${bucket}" \
    --region "${REGION}" || true
done

# ---------------------------------------------------------------------------
# KMS — Dev encryption key
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating KMS key..."

KMS_KEY_ID=$(awslocal kms create-key \
  --description "Complai dev encryption key" \
  --region "${REGION}" \
  --query 'KeyMetadata.KeyId' \
  --output text)

awslocal kms create-alias \
  --alias-name "alias/complai-dev" \
  --target-key-id "${KMS_KEY_ID}" \
  --region "${REGION}"

echo "  -> Key created: ${KMS_KEY_ID} (alias/complai-dev)"

# ---------------------------------------------------------------------------
# Secrets Manager — Test secrets
# ---------------------------------------------------------------------------
echo ""
echo "==> Creating Secrets Manager test secrets..."

awslocal secretsmanager create-secret \
  --name "complai/dev/db-credentials" \
  --secret-string '{"username":"complai","password":"complai_dev","host":"postgres","port":"5432"}' \
  --region "${REGION}"

awslocal secretsmanager create-secret \
  --name "complai/dev/redis-credentials" \
  --secret-string '{"host":"redis","port":"6379"}' \
  --region "${REGION}"

awslocal secretsmanager create-secret \
  --name "complai/dev/adaequare-credentials" \
  --secret-string '{"api_key":"test-adaequare-key","api_secret":"test-adaequare-secret","base_url":"https://sandbox.adaequare.com"}' \
  --region "${REGION}"

awslocal secretsmanager create-secret \
  --name "complai/dev/sandbox-credentials" \
  --secret-string '{"api_key":"test-sandbox-key","api_secret":"test-sandbox-secret","base_url":"https://api.sandbox.co.in"}' \
  --region "${REGION}"

awslocal secretsmanager create-secret \
  --name "complai/dev/keycloak-credentials" \
  --secret-string '{"admin_user":"admin","admin_password":"admin","client_secret":"dev-client-secret"}' \
  --region "${REGION}"

awslocal secretsmanager create-secret \
  --name "complai/dev/jwt-signing-key" \
  --secret-string '{"key":"dev-jwt-signing-key-do-not-use-in-production-32bytes!"}' \
  --region "${REGION}"

echo ""
echo "============================================================"
echo "  Complai LocalStack Init — Complete"
echo "============================================================"
