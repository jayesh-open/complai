# ---------------------------------------------------------------
# Application Load Balancer — Complai Ingress
# ---------------------------------------------------------------
# ALB in public subnets, HTTPS listener with ACM cert,
# security group restricted to Cloudflare IP ranges only.
# ---------------------------------------------------------------

# ---------------------------------------------------------------
# Security Group — Cloudflare-only ingress
# ---------------------------------------------------------------

resource "aws_security_group" "alb" {
  name_prefix = "${var.project}-${var.environment}-alb-"
  description = "ALB security group — allows HTTPS from Cloudflare IPs only"
  vpc_id      = var.vpc_id

  tags = {
    Name        = "${var.project}-${var.environment}-alb-sg"
    Environment = var.environment
    Project     = var.project
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "alb_https_ingress" {
  count = length(var.cloudflare_ip_ranges)

  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = [var.cloudflare_ip_ranges[count.index]]
  security_group_id = aws_security_group.alb.id
  description       = "Allow HTTPS from Cloudflare IP range"
}

resource "aws_security_group_rule" "alb_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.alb.id
  description       = "Allow all outbound traffic"
}

# ---------------------------------------------------------------
# Application Load Balancer
# ---------------------------------------------------------------

resource "aws_lb" "main" {
  name               = "${var.project}-${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.environment == "prod"
  drop_invalid_header_fields = true

  tags = {
    Name        = "${var.project}-${var.environment}-alb"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# Target Group (for EKS Istio ingress gateway)
# ---------------------------------------------------------------

resource "aws_lb_target_group" "eks" {
  name        = "${var.project}-${var.environment}-eks-tg"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    path                = "/healthz"
    matcher             = "200"
  }

  tags = {
    Name        = "${var.project}-${var.environment}-eks-tg"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# HTTPS Listener
# ---------------------------------------------------------------

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.eks.arn
  }

  tags = {
    Name        = "${var.project}-${var.environment}-https-listener"
    Environment = var.environment
    Project     = var.project
  }
}

# ---------------------------------------------------------------
# HTTP → HTTPS Redirect
# ---------------------------------------------------------------

resource "aws_lb_listener" "http_redirect" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  tags = {
    Name        = "${var.project}-${var.environment}-http-redirect"
    Environment = var.environment
    Project     = var.project
  }
}
