# ---------------------------------------------------------------
# ACM Certificate — Complai TLS
# ---------------------------------------------------------------
# Wildcard certificate for *.complai.in with DNS validation.
# ---------------------------------------------------------------

resource "aws_acm_certificate" "main" {
  domain_name               = "*.${var.domain_name}"
  subject_alternative_names = [var.domain_name]
  validation_method         = "DNS"

  tags = {
    Name        = "${var.domain_name}-wildcard"
    Environment = var.environment
    Project     = var.project
  }

  lifecycle {
    create_before_destroy = true
  }
}

# ---------------------------------------------------------------
# DNS Validation Records
# ---------------------------------------------------------------
# These records must be created in the DNS provider (Cloudflare).
# If using Cloudflare provider, create the records there.
# If validating via Route53, use aws_route53_record instead.
# ---------------------------------------------------------------

resource "aws_acm_certificate_validation" "main" {
  certificate_arn = aws_acm_certificate.main.arn

  # Validation will wait until DNS records are resolvable.
  # Ensure the Cloudflare DNS module creates the validation records.
}
