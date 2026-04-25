# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the Complai platform.

## Format

Each ADR follows this template:

```markdown
# ADR-NNNN: Title

**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-XXXX
**Date:** YYYY-MM-DD
**Deciders:** [list]

## Context

What is the issue that we're seeing that is motivating this decision?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

What becomes easier or more difficult to do because of this change?

### Positive
- ...

### Negative
- ...

### Risks
- ...
```

## Naming convention

`NNNN-short-kebab-title.md` -- e.g., `0001-multi-tenancy-via-postgres-rls.md`

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [0001](0001-multi-tenancy-via-postgres-rls.md) | Multi-tenancy via Postgres RLS | Accepted |
| [0002](0002-adaequare-enriched-apis-only.md) | Adaequare Enriched APIs only (no pass-through/SEK) | Accepted |
| [0003](0003-two-provider-api-strategy.md) | Two-provider API strategy (Adaequare + Sandbox.co.in) | Accepted |
| [0004](0004-go-as-primary-backend-language.md) | Go as primary backend language | Accepted |
| [0005](0005-aws-ap-south-1-primary.md) | AWS as cloud provider, ap-south-1 primary | Accepted |
| [0006](0006-postgres-only-for-phase-1.md) | Postgres-only for Phase 1 (OLTP + analytics) | Accepted |
| [0007](0007-sqs-sns-over-kafka.md) | SQS/SNS over Kafka for Phase 1 messaging | Accepted |
| [0008](0008-temporal-cloud-for-workflows.md) | Temporal Cloud (managed) for workflow orchestration | Accepted |
| [0009](0009-cloudflare-for-cdn-dns-waf.md) | Cloudflare for CDN/DNS/WAF | Accepted |
| [0010](0010-amazon-ses-for-email.md) | Amazon SES for email | Accepted |
| [0011](0011-last9-for-observability.md) | Last9 for observability | Accepted |
| [0012](0012-keycloak-self-hosted-for-identity.md) | Keycloak self-hosted for identity | Accepted |
| [0013](0013-outbox-pattern-via-polling.md) | Outbox pattern via polling (not Debezium) in Phase 1 | Accepted |
| [0014](0014-canonical-invoice-schema.md) | Canonical Invoice Schema as lingua franca | Accepted |
| [0015](0015-monorepo-with-go-workspaces-pnpm-turborepo.md) | Monorepo with Go workspaces + pnpm + Turborepo | Accepted |
