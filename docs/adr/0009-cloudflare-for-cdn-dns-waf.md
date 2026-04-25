# ADR-0009: Cloudflare for CDN/DNS/WAF

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai needs CDN (static asset delivery for Next.js apps), DNS management, WAF (web application firewall), and DDoS protection in front of our AWS ALB. These are distinct concerns but tightly coupled operationally.

Options evaluated:

1. **AWS-native stack** -- CloudFront (CDN) + Route 53 (DNS) + AWS WAF -- tightly integrated with AWS but more expensive and less flexible for edge rules.
2. **Cloudflare** -- unified CDN + DNS + WAF + DDoS in one platform, global edge network, simpler operational model.
3. **Mix** -- Route 53 for DNS, Cloudflare for CDN/WAF -- adds complexity without clear benefit.

## Decision

Cloudflare (Business plan) for DNS, CDN, WAF, DDoS protection, and TLS termination for all Complai domains (`complai.in`, `app.complai.in`, `api.complai.in`).

- **CDN:** caches static assets (JS, CSS, images) at Cloudflare's global edge.
- **WAF:** OWASP Top 10 managed rules, custom rules for India-specific attack patterns, bot management.
- **DDoS protection:** included at the Business tier.
- **TLS termination:** at Cloudflare edge with origin-pull TLS to ALB.
- **Origin lockdown:** ALB security group accepts traffic only from Cloudflare IP ranges. Direct-to-origin bypass is blocked.

## Consequences

### Positive
- Global edge network reduces latency for static assets, especially for users across India's varied network infrastructure.
- Integrated WAF + DDoS protection in a single platform -- no separate configuration or billing for each.
- DNS management through Cloudflare's UI/API is straightforward and well-documented.
- Origin lockdown ensures all traffic passes through WAF/DDoS protection -- no bypass path.
- Business plan cost (~200 USD/month) is competitive with equivalent AWS-native setup.

### Negative
- Additional vendor relationship outside of AWS.
- Must keep the Cloudflare IP allowlist updated in ALB security groups (Cloudflare publishes IP ranges; we automate updates).
- Some advanced features (custom page rules, advanced bot management) require Business or Enterprise plan tiers.

### Risks
- Cloudflare outage would make the platform unreachable even if AWS infrastructure is healthy. Mitigated by: Cloudflare's track record of 99.99%+ uptime; in a prolonged outage, we can temporarily update DNS to point directly to ALB (sacrificing WAF/CDN protection).
