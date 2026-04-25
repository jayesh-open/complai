# ---------------------------------------------------------------
# S3 Backend Configuration for Complai Terraform State
# ---------------------------------------------------------------
# This backend block is commented out for scaffolding purposes.
# Before using, ensure the S3 bucket and DynamoDB table exist.
# Bootstrap them manually or via a separate "bootstrap" Terraform
# config before uncommenting this block.
#
# terraform {
#   backend "s3" {
#     bucket         = "complai-terraform-state"
#     key            = "complai/terraform.tfstate"
#     region         = "ap-south-1"
#     dynamodb_table = "complai-terraform-locks"
#     encrypt        = true
#   }
# }
