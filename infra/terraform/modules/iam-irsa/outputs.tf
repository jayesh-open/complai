output "role_arns" {
  description = "Map of service name to IAM role ARN"
  value = {
    for k, v in aws_iam_role.service : k => v.arn
  }
}

output "role_names" {
  description = "Map of service name to IAM role name"
  value = {
    for k, v in aws_iam_role.service : k => v.name
  }
}
