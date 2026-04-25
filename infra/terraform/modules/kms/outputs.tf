output "platform_key_arn" {
  description = "ARN of the platform KMS CMK"
  value       = aws_kms_key.platform.arn
}

output "platform_key_id" {
  description = "Key ID of the platform KMS CMK"
  value       = aws_kms_key.platform.key_id
}

output "platform_key_alias" {
  description = "Alias of the platform KMS CMK"
  value       = aws_kms_alias.platform.name
}
