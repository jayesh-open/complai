# ADR-0012: Keycloak self-hosted for identity

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai serves enterprise customers who require SSO integration (Google Workspace, Azure AD, Okta), multi-factor authentication, and tenant-isolated identity management. We need an identity provider that supports:

- Per-tenant SSO configuration (each enterprise customer may use a different IdP).
- SAML 2.0 and OIDC protocols.
- MFA (TOTP, SMS, email).
- User federation and social login.
- Custom authentication flows (step-up auth for filing operations).

Options evaluated:

1. **Amazon Cognito** -- managed, AWS-native, but per-user-pool model becomes awkward at >100 tenants. Limited customization of auth flows.
2. **Auth0** -- excellent developer experience, but expensive at enterprise scale and SaaS dependency for a critical path.
3. **Keycloak (self-hosted)** -- per-realm model maps cleanly to multi-tenant architecture. Full control over auth flows. Open source.

## Decision

Self-hosted Keycloak 24 on EKS with a Postgres backend on RDS (separate `keycloak_db` logical database).

- **3-node cluster** behind an internal ALB for high availability.
- **Per-tenant realms** for isolated SSO configuration -- each enterprise customer gets their own Keycloak realm with their own IdP connections, branding, and MFA policies.
- **Protocols:** SAML 2.0 and OIDC for enterprise SSO; OIDC for all internal flows.
- **MFA:** TOTP (authenticator apps), SMS (via SES), and email as second factors.
- **Step-up authentication:** filing operations require re-authentication within a 5-minute window.
- **JWT claims** carry `tenant_id`, `active_pan`, `active_gstin`, roles, and permissions.

## Consequences

### Positive
- Per-realm model maps naturally to multi-tenant enterprise SSO -- each tenant's IdP configuration is fully isolated.
- Full control over authentication flows, token claims, and session management.
- Supports SAML 2.0 + OIDC + social login out of the box, covering all enterprise SSO requirements.
- Open source with active community and predictable release cadence.
- No per-user licensing costs.

### Negative
- Operational overhead: EKS deployment, health monitoring, version upgrades, and configuration backup require a few engineer-hours per month.
- Must manage Keycloak's Postgres database alongside the application databases (though it runs on the same RDS instance).
- Keycloak's admin UI and configuration model have a learning curve.

### Risks
- Keycloak vulnerability requiring emergency patching. Mitigated by: WAF in front of the auth endpoints, security mailing list subscription for early notification, and a tested upgrade runbook.
- Realm count scaling (1000+ tenants) may require Keycloak performance tuning. Mitigated by: realm metadata is small; the primary bottleneck is authentication throughput, which scales with EKS pod count.
