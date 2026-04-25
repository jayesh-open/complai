# ---------------------------------------------------------------
# AWS Secrets Manager — Complai Secrets
# ---------------------------------------------------------------
# Stores provider credentials, DB password, and JWT signing key.
# All secrets encrypted with the platform KMS CMK.
# ---------------------------------------------------------------

# ---------------------------------------------------------------
# Adaequare API Credentials
# ---------------------------------------------------------------

resource "aws_secretsmanager_secret" "adaequare_creds" {
  name       = "${var.environment}/complai/adaequare-creds"
  kms_key_id = var.kms_key_arn

  tags = {
    Name        = "${var.environment}/complai/adaequare-creds"
    Environment = var.environment
    Project     = "complai"
    Provider    = "adaequare"
  }
}

resource "aws_secretsmanager_secret_version" "adaequare_creds" {
  secret_id = aws_secretsmanager_secret.adaequare_creds.id
  secret_string = jsonencode({
    asp_id     = "PLACEHOLDER_ASP_ID"
    asp_secret = "PLACEHOLDER_ASP_SECRET"
    auth_token = "PLACEHOLDER_AUTH_TOKEN"
    base_url   = "https://api.adaequare.com"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# ---------------------------------------------------------------
# Sandbox.co.in API Credentials
# ---------------------------------------------------------------

resource "aws_secretsmanager_secret" "sandbox_creds" {
  name       = "${var.environment}/complai/sandbox-creds"
  kms_key_id = var.kms_key_arn

  tags = {
    Name        = "${var.environment}/complai/sandbox-creds"
    Environment = var.environment
    Project     = "complai"
    Provider    = "sandbox"
  }
}

resource "aws_secretsmanager_secret_version" "sandbox_creds" {
  secret_id = aws_secretsmanager_secret.sandbox_creds.id
  secret_string = jsonencode({
    api_key    = "PLACEHOLDER_API_KEY"
    api_secret = "PLACEHOLDER_API_SECRET"
    base_url   = "https://api.sandbox.co.in"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# ---------------------------------------------------------------
# Database Master Password
# ---------------------------------------------------------------

resource "aws_secretsmanager_secret" "db_master_password" {
  name       = "${var.environment}/complai/db-master-password"
  kms_key_id = var.kms_key_arn

  tags = {
    Name        = "${var.environment}/complai/db-master-password"
    Environment = var.environment
    Project     = "complai"
  }
}

resource "aws_secretsmanager_secret_version" "db_master_password" {
  secret_id = aws_secretsmanager_secret.db_master_password.id
  secret_string = jsonencode({
    username = "complai_admin"
    password = "PLACEHOLDER_CHANGE_ME"
    host     = "PLACEHOLDER_RDS_ENDPOINT"
    port     = 5432
    dbname   = "complai"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# ---------------------------------------------------------------
# DB Password Rotation (Lambda placeholder)
# ---------------------------------------------------------------
# TODO: Implement rotation Lambda when RDS is provisioned.
# Use aws_secretsmanager_secret_rotation with a Lambda that
# calls the RDS API to rotate the master password.
#
# resource "aws_secretsmanager_secret_rotation" "db_master_password" {
#   secret_id           = aws_secretsmanager_secret.db_master_password.id
#   rotation_lambda_arn = aws_lambda_function.secret_rotation.arn
#
#   rotation_rules {
#     automatically_after_days = 30
#   }
# }

# ---------------------------------------------------------------
# JWT Signing Key
# ---------------------------------------------------------------

resource "aws_secretsmanager_secret" "jwt_signing_key" {
  name       = "${var.environment}/complai/jwt-signing-key"
  kms_key_id = var.kms_key_arn

  tags = {
    Name        = "${var.environment}/complai/jwt-signing-key"
    Environment = var.environment
    Project     = "complai"
  }
}

resource "aws_secretsmanager_secret_version" "jwt_signing_key" {
  secret_id = aws_secretsmanager_secret.jwt_signing_key.id
  secret_string = jsonencode({
    algorithm   = "RS256"
    private_key = "PLACEHOLDER_GENERATE_RSA_KEY"
    public_key  = "PLACEHOLDER_GENERATE_RSA_KEY"
    key_id      = "complai-jwt-v1"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}
