# ---------------------------------------------------------------
# Cloudflare DNS + WAF — Complai Edge
# ---------------------------------------------------------------
# DNS records pointing to ALB, WAF with OWASP managed rules.
# Cloudflare sits in front of the ALB for CDN, WAF, and DDoS.
# ---------------------------------------------------------------

terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

# ---------------------------------------------------------------
# DNS Records
# ---------------------------------------------------------------

# Root domain → ALB (proxied through Cloudflare)
resource "cloudflare_record" "root" {
  zone_id = var.zone_id
  name    = var.domain
  content = var.alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied

  comment = "Complai ${var.environment} - root domain to ALB"
}

# Wildcard → ALB (proxied through Cloudflare)
resource "cloudflare_record" "wildcard" {
  zone_id = var.zone_id
  name    = "*.${var.domain}"
  content = var.alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1

  comment = "Complai ${var.environment} - wildcard to ALB"
}

# API subdomain → ALB
resource "cloudflare_record" "api" {
  zone_id = var.zone_id
  name    = "api.${var.domain}"
  content = var.alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1

  comment = "Complai ${var.environment} - API endpoint"
}

# Vendor portal subdomain → ALB
resource "cloudflare_record" "vendor_portal" {
  zone_id = var.zone_id
  name    = "vendor.${var.domain}"
  content = var.alb_dns_name
  type    = "CNAME"
  proxied = true
  ttl     = 1

  comment = "Complai ${var.environment} - vendor portal"
}

# ---------------------------------------------------------------
# WAF Ruleset — OWASP Managed Rules
# ---------------------------------------------------------------

resource "cloudflare_ruleset" "waf_managed" {
  zone_id     = var.zone_id
  name        = "Complai WAF - OWASP Managed Rules"
  description = "Deploy Cloudflare OWASP managed ruleset for Complai"
  kind        = "zone"
  phase       = "http_request_firewall_managed"

  # Cloudflare OWASP Core Ruleset
  rules {
    action = "execute"
    action_parameters {
      id      = "efb7b8c949ac4650a09736fc376e9aee"
      version = "latest"
    }
    expression  = "true"
    description = "Execute Cloudflare Managed OWASP Core Ruleset"
    enabled     = true
  }

  # Cloudflare Managed Ruleset
  rules {
    action = "execute"
    action_parameters {
      id      = "efb7b8c949ac4650a09736fc376e9aee"
      version = "latest"
    }
    expression  = "true"
    description = "Execute Cloudflare Managed Ruleset"
    enabled     = true
  }
}

# ---------------------------------------------------------------
# Rate Limiting for API endpoints
# ---------------------------------------------------------------

resource "cloudflare_ruleset" "rate_limit" {
  zone_id     = var.zone_id
  name        = "Complai Rate Limiting"
  description = "Rate limiting rules for API and auth endpoints"
  kind        = "zone"
  phase       = "http_ratelimit"

  rules {
    action = "block"
    ratelimit {
      characteristics     = ["ip.src"]
      period              = 60
      requests_per_period = 100
      mitigation_timeout  = 600
    }
    expression  = "(http.request.uri.path contains \"/api/v1/auth\")"
    description = "Rate limit authentication endpoints"
    enabled     = true
  }
}
