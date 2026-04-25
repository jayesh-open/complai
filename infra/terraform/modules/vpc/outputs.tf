output "vpc_id" {
  description = "ID of the created VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of public subnets (one per AZ)"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of private subnets (one per AZ) for EKS and application workloads"
  value       = aws_subnet.private[*].id
}

output "data_subnet_ids" {
  description = "IDs of data subnets (one per AZ) for RDS, ElastiCache, OpenSearch — no internet egress"
  value       = aws_subnet.data[*].id
}

output "nat_gateway_ips" {
  description = "Elastic IP addresses assigned to NAT Gateways"
  value       = [aws_eip.nat.public_ip]
}
