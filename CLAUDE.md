# CLAUDE.md — Session-durable memory for Complai

## What we're building
Complai is the compliance layer in the Bank Open product family — an enterprise GST + TDS + compliance SaaS for Indian enterprises. Target: 500 Year-1, 50K GSTINs Year-3. Competitive reference: Clear / ClearTax.

## Bank Open ecosystem
Complai is one of four sibling apps plus an external HRMS:
- **Apex** — Procure-to-Pay (P2P). Owns vendor master, AP invoices, POs, GRNs, payments. In UAT.
- **Aura** — Order-to-Cash (O2C). Owns customer master, AR invoices, collections. Early stage.
- **Bridge** — Contract management. Owns contracts, obligations, renewals. Early stage.
- **Complai** — Compliance only. Consumes data from siblings via gateways.
- **HRMS** — External. Payroll data, Form 16.

**Boundary rule:** Complai does NOT own vendor master (Apex), AP invoices (Apex), AR invoices (Aura), or contracts (Bridge). It consumes from siblings via gateways and adds compliance value on top.

**7 compliance modules:** GST Returns, E-Invoicing, E-Way Bill, ITC Reconciliation + MaxITC + Vendor Compliance Scoring, TDS/TCS, ITR, Secretarial.

**Shared compliance modules (Phase 2):** e-Invoice, E-Way Bill, GSTR-2A/2B view, MaxITC view exist in BOTH Aura and Complai with real-time sync. Phase 1 = standalone.

## Authoritative input docs (read on every session start)
- /docs/input/complai_prd.md
- /docs/input/complai_architecture.md
- /docs/input/complai_design_system.md
- /docs/input/complai_api_integration.md
- /docs/input/AURA_DESIGN_SYSTEM.md

## Non-negotiable invariants
1. Multi-tenant: tenant_id on every row, Postgres RLS, asserted in every service
2. Enriched APIs only (Adaequare). No SEK/crypto in our repo.
3. Two providers: Adaequare (GST/IRP/EWB) + Sandbox.co.in (TDS/ITR/KYC/Tax-Payment)
4. Cloud: AWS, ap-south-1 (Mumbai) primary, ap-south-2 (Hyderabad) DR
5. Backend: Go 1.22 everywhere (domain + gateways); Python 3.12 for AI only
6. Frontend: TypeScript 5.4 + Next.js 15
7. OLTP + Phase 1 analytics: Postgres RDS (Multi-AZ + read replica)
8. No MongoDB. No Kafka. No ClickHouse in Phase 1.
9. Messaging: SQS + SNS, abstracted behind MessageBus interface
10. Design system: Light Classic default, Compact density, DD/MM/YYYY, ₹ + crore/lakh
11. Outbox pattern for every external call
12. Idempotency via request_id UUIDs
13. OTel traces with tenant_id + gstin/tan/pan baggage → Last9

## Stack (pinned)
- Cloud: AWS (ap-south-1 primary, ap-south-2 DR)
- Compute: EKS 1.30 + Istio 1.22 ambient
- Backend: Go 1.22 (domain + gateways), Python 3.12 (AI only)
- Frontend: TypeScript 5.4 + Next.js 15 + React 19 + Tailwind + shadcn/ui
- OLTP: Amazon RDS PostgreSQL 16 Multi-AZ + read replica
- Cache: ElastiCache Redis 7
- Messaging: SQS + SNS (MessageBus abstraction)
- Search: Amazon OpenSearch Service 2
- Storage: S3
- Secrets: AWS Secrets Manager + KMS
- Identity: Keycloak 24 on EKS
- Workflow: Temporal Cloud (managed)
- Email: Amazon SES
- CDN/DNS/WAF: Cloudflare
- Observability: Last9 (OpenTelemetry-native)
- CI/CD: GitHub Actions + ArgoCD

## Go framework stack (consistent across all services)
- HTTP: go-chi/chi/v5
- DB: jackc/pgx/v5 + sqlc + goose
- DI: uber-go/fx
- Config: spf13/viper
- Logging: rs/zerolog
- Validation: go-playground/validator/v10
- Tracing: OpenTelemetry Go SDK
- Testing: testify + testcontainers-go
- Mocking: gomock
- Money: shopspring/decimal
- JWT: golang-jwt/jwt/v5
- Temporal: temporal.io/sdk
- AWS: aws-sdk-go-v2
- Circuit breaker: sony/gobreaker
- Linting: golangci-lint

## Repo layout
- apps/web — Next.js main product
- services/go/{name}-service — Go services (domain + gateways, including sibling gateways)
- services/python/{name}-service — Python AI services
- services/node/{name}-bff-service — TypeScript BFFs
- packages/shared-kernel-go — Go shared libs (tenant, outbox, messagebus)
- packages/shared-kernel-node — TS shared libs
- packages/ui-components — Complai component library
- packages/events — Protobuf schemas
- packages/openapi — OpenAPI specs
- infra/terraform — AWS infra
- infra/helm — K8s charts
- docs/adr — ADRs
- docs/input — authoritative docs

## Local dev environment
- All AWS service interactions use LocalStack (via docker-compose.dev.yml)
- Go services use aws-sdk-go-v2 with AWS_ENDPOINT_URL pointed at LocalStack
- AWS region is ap-south-1 everywhere — awslocal and aws-sdk-go-v2 must use this region
- Mailpit replaces Amazon SES for email
- Temporal dev server replaces Temporal Cloud
- Jaeger replaces Last9 for tracing
- LocalStack KMS replaces real KMS
- Terraform files generated as scaffolding only — not executed locally
- No AWS CLI or Terraform installed on dev machine
- All 11 service databases auto-provisioned via `scripts/postgres-init.sh`; apply migrations with `make migrate-all`

## Current build state
- [x] Part 0.5: Repo init + memory scaffolding
- [x] Part 1: Repo skeleton + shared foundation
- [x] Part 2: Identity + Tenant + User/Role services + auth
- [x] Part 3: Platform services (master-data, document, notification, audit, workflow, rules)
- [x] Part 4: API Gateway + BFF + Web Shell + design system components
- [x] Part 4.5: Scope correction — align with Bank Open ecosystem
- [x] Part 5: Adaequare GST gateway + GSTR-1 flow
- [x] Part 6: Sandbox KYC gateway + Vendor Compliance + Apex Sync
- [ ] Part 7: Reconciliation engine + GSTR-3B + GSTR-2B/IMS (AP register from Apex)
- [ ] Part 8: e-Invoicing + E-Way Bill
- [ ] Part 9: Sandbox TDS gateway + TDS module
- [ ] Part 10: Sandbox ITR + GSTR-9/9C
- [ ] Part 11: Sibling gateway services (Aura, Bridge, HRMS)
- [ ] Part 12: AI layer + MaxITC
- [ ] Part 13: Real Bank Open sibling sync + GL-Stream + Compliance Cloud
- [ ] Part 14: Reporting + observability + production hardening

## Credentials / blockers needed
(populated as encountered)

## Key ADRs
- ADR-0001: Multi-tenancy via Postgres RLS
- ADR-0002: Adaequare Enriched APIs only (no pass-through)
- ADR-0003: Two-provider API strategy (Adaequare + Sandbox.co.in)
- ADR-0004: Go as primary backend language
- ADR-0005: AWS as cloud provider, ap-south-1 primary
- ADR-0006: Postgres-only for Phase 1 (OLTP + analytics)
- ADR-0007: SQS/SNS over Kafka for Phase 1 messaging
- ADR-0008: Temporal Cloud (managed) for workflows
- ADR-0009: Cloudflare for CDN/DNS/WAF
- ADR-0010: Amazon SES for email
- ADR-0011: Last9 for observability
- ADR-0012: Keycloak self-hosted for identity
- ADR-0013: Outbox pattern via polling (not Debezium) in Phase 1
- ADR-0014: Canonical Invoice Schema as lingua franca
- ADR-0015: Monorepo with Go workspaces + pnpm + Turborepo
