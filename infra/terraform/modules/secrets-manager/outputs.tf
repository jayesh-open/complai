output "secret_arns" {
  description = "Map of secret name to ARN"
  value = {
    adaequare_creds    = aws_secretsmanager_secret.adaequare_creds.arn
    sandbox_creds      = aws_secretsmanager_secret.sandbox_creds.arn
    db_master_password = aws_secretsmanager_secret.db_master_password.arn
    jwt_signing_key    = aws_secretsmanager_secret.jwt_signing_key.arn
  }
}
