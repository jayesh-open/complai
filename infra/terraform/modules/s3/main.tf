# ---------------------------------------------------------------
# S3 Buckets — Complai Document & Data Storage
# ---------------------------------------------------------------
# Four buckets: documents, uploads, exports, backups.
# All encrypted with KMS, versioned, public access blocked.
# Lifecycle rules transition documents through storage tiers.
# ---------------------------------------------------------------

locals {
  buckets = {
    documents = {
      name = "${var.project}-${var.environment}-documents"
    }
    uploads = {
      name = "${var.project}-${var.environment}-uploads"
    }
    exports = {
      name = "${var.project}-${var.environment}-exports"
    }
    backups = {
      name = "${var.project}-${var.environment}-backups"
    }
  }
}

# ---------------------------------------------------------------
# S3 Buckets
# ---------------------------------------------------------------

resource "aws_s3_bucket" "buckets" {
  for_each = local.buckets

  bucket = each.value.name

  tags = {
    Name        = each.value.name
    Environment = var.environment
    Project     = var.project
    Purpose     = each.key
  }
}

# ---------------------------------------------------------------
# Block Public Access (all buckets)
# ---------------------------------------------------------------

resource "aws_s3_bucket_public_access_block" "buckets" {
  for_each = local.buckets

  bucket = aws_s3_bucket.buckets[each.key].id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# ---------------------------------------------------------------
# Versioning (all buckets)
# ---------------------------------------------------------------

resource "aws_s3_bucket_versioning" "buckets" {
  for_each = local.buckets

  bucket = aws_s3_bucket.buckets[each.key].id

  versioning_configuration {
    status = "Enabled"
  }
}

# ---------------------------------------------------------------
# Server-Side Encryption — KMS (all buckets)
# ---------------------------------------------------------------

resource "aws_s3_bucket_server_side_encryption_configuration" "buckets" {
  for_each = local.buckets

  bucket = aws_s3_bucket.buckets[each.key].id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = var.kms_key_arn
    }
    bucket_key_enabled = true
  }
}

# ---------------------------------------------------------------
# Lifecycle Rules — Documents bucket
# Standard -> IA (90d) -> Glacier (2yr) -> Deep Archive (5yr)
# ---------------------------------------------------------------

resource "aws_s3_bucket_lifecycle_configuration" "documents" {
  bucket = aws_s3_bucket.buckets["documents"].id

  rule {
    id     = "document-tiering"
    status = "Enabled"

    transition {
      days          = 90
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 730
      storage_class = "GLACIER"
    }

    transition {
      days          = 1825
      storage_class = "DEEP_ARCHIVE"
    }
  }

  rule {
    id     = "cleanup-incomplete-multipart"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}

# ---------------------------------------------------------------
# Lifecycle Rules — Uploads bucket (temporary)
# ---------------------------------------------------------------

resource "aws_s3_bucket_lifecycle_configuration" "uploads" {
  bucket = aws_s3_bucket.buckets["uploads"].id

  rule {
    id     = "expire-stale-uploads"
    status = "Enabled"

    expiration {
      days = 30
    }
  }

  rule {
    id     = "cleanup-incomplete-multipart"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 3
    }
  }
}

# ---------------------------------------------------------------
# Lifecycle Rules — Exports bucket
# ---------------------------------------------------------------

resource "aws_s3_bucket_lifecycle_configuration" "exports" {
  bucket = aws_s3_bucket.buckets["exports"].id

  rule {
    id     = "expire-old-exports"
    status = "Enabled"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    expiration {
      days = 90
    }
  }
}

# ---------------------------------------------------------------
# Lifecycle Rules — Backups bucket
# ---------------------------------------------------------------

resource "aws_s3_bucket_lifecycle_configuration" "backups" {
  bucket = aws_s3_bucket.buckets["backups"].id

  rule {
    id     = "backup-tiering"
    status = "Enabled"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 90
      storage_class = "GLACIER"
    }

    transition {
      days          = 365
      storage_class = "DEEP_ARCHIVE"
    }
  }
}

# ---------------------------------------------------------------
# VPC Endpoint for S3 (Gateway type — free, no NAT needed)
# ---------------------------------------------------------------

resource "aws_vpc_endpoint" "s3" {
  vpc_id       = var.vpc_id
  service_name = "com.amazonaws.${var.region}.s3"

  vpc_endpoint_type = "Gateway"
  route_table_ids   = var.route_table_ids

  tags = {
    Name        = "${var.project}-${var.environment}-s3-vpce"
    Environment = var.environment
    Project     = var.project
  }
}
