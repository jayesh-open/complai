# ---------------------------------------------------------------
# VPC — Complai Network Foundation
# ---------------------------------------------------------------
# /16 CIDR split into public, private (application), and data
# (no internet egress) subnets across 3 AZs.
# ---------------------------------------------------------------

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name        = "${var.project}-${var.environment}-vpc"
    Environment = var.environment
    Project     = var.project
    ManagedBy   = "terraform"
  }
}

# ---------------------------------------------------------------
# Internet Gateway
# ---------------------------------------------------------------

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${var.project}-${var.environment}-igw"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Elastic IP for NAT Gateway
# ---------------------------------------------------------------

resource "aws_eip" "nat" {
  domain = "vpc"

  tags = {
    Name        = "${var.project}-${var.environment}-nat-eip"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# NAT Gateway (single for cost; prod may use one per AZ)
# ---------------------------------------------------------------

resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id

  tags = {
    Name        = "${var.project}-${var.environment}-nat"
    Environment = var.environment
    Project     = var.project
  }

  depends_on = [aws_internet_gateway.main]
}

# ---------------------------------------------------------------
# Public Subnets (3 AZs) — ALB, NAT Gateway
# ---------------------------------------------------------------

resource "aws_subnet" "public" {
  count = length(var.azs)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = var.azs[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name                                = "${var.project}-${var.environment}-public-${var.azs[count.index]}"
    Environment                         = var.environment
    Project                             = var.project
    "kubernetes.io/role/elb"            = "1"
    "kubernetes.io/cluster/${var.project}-${var.environment}" = "shared"
  }
}

# ---------------------------------------------------------------
# Private Subnets (3 AZs) — EKS nodes, application workloads
# ---------------------------------------------------------------

resource "aws_subnet" "private" {
  count = length(var.azs)

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 10)
  availability_zone = var.azs[count.index]

  tags = {
    Name                                = "${var.project}-${var.environment}-private-${var.azs[count.index]}"
    Environment                         = var.environment
    Project                             = var.project
    "kubernetes.io/role/internal-elb"    = "1"
    "kubernetes.io/cluster/${var.project}-${var.environment}" = "shared"
  }
}

# ---------------------------------------------------------------
# Data Subnets (3 AZs) — RDS, ElastiCache, OpenSearch (no internet)
# ---------------------------------------------------------------

resource "aws_subnet" "data" {
  count = length(var.azs)

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 20)
  availability_zone = var.azs[count.index]

  tags = {
    Name        = "${var.project}-${var.environment}-data-${var.azs[count.index]}"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Route Tables
# ---------------------------------------------------------------

# Public route table — routes to IGW
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name        = "${var.project}-${var.environment}-public-rt"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_route_table_association" "public" {
  count = length(var.azs)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Private route table — routes to NAT Gateway
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main.id
  }

  tags = {
    Name        = "${var.project}-${var.environment}-private-rt"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_route_table_association" "private" {
  count = length(var.azs)

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

# Data route table — NO internet route (isolated)
resource "aws_route_table" "data" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name        = "${var.project}-${var.environment}-data-rt"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_route_table_association" "data" {
  count = length(var.azs)

  subnet_id      = aws_subnet.data[count.index].id
  route_table_id = aws_route_table.data.id
}
