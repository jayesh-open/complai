# BUILD_PLAN.md — Living checklist

Last updated: 2026-04-25

## Current part
Part 4.5 complete (Bank Open scope correction). Next = Part 5.

## Completed
- [x] Part 0.5: Repo init, CLAUDE.md, BUILD_PLAN.md, input docs, ADR template
- [x] Part 1: Repo skeleton + shared foundation
- [x] Part 2: Identity + Tenant + User/Role services + auth
  - [x] identity-service: Keycloak OIDC login/refresh/logout, JWT with tenant_id claim, MFA enrollment, step-up auth (5-min window)
  - [x] tenant-service: CRUD + PAN/GSTIN/TAN hierarchy, per-tenant KMS CMK via LocalStack, suspend/reactivate
  - [x] user-role-service: RBAC (roles, permissions, policy check), maker-checker approval workflow (self-approval denied), role templates
  - [x] Postgres RLS enforced on all 3 databases (complai_app role, SET LOCAL app.tenant_id)
  - [x] seed-dev.sh: 3 tenants, 15 users, 15 roles, 55 permissions, 15 Keycloak users with tenant_id attribute
  - [x] Dockerfiles for all 3 services
  - [x] Unit test coverage: identity 98.4%, tenant 100%, user-role 100%
  - [x] E2E verified: auth flow, RLS isolation, maker-checker, step-up, KMS key creation
  - [x] Monorepo scaffolding (pnpm workspaces + Turborepo + go.work)
  - [x] Shared Go kernel (packages/shared-kernel-go) — 11 sub-packages, all tests pass
  - [x] Shared Node kernel (packages/shared-kernel-node) — types, Zod schemas, formatters, client, errors
  - [x] UI components scaffolding (packages/ui-components)
  - [x] Events library (packages/events) — 6 Protobuf schemas + TS types
  - [x] OpenAPI library (packages/openapi) — health-probe spec + codegen config
  - [x] Terraform scaffolding — 13 modules, 4 environments (dev/sandbox/staging/prod), 57 .tf files
  - [x] Helm charts — base chart (Istio/ESO/ArgoCD/cert-manager) + service template chart
  - [x] Docker compose dev environment — 11 services (Postgres, Redis, LocalStack, OpenSearch, Keycloak, Temporal, Mailpit, Jaeger, OTel Collector)
  - [x] CI/CD workflows — ci.yml, deploy-staging.yml, security.yml
  - [x] Dev tooling — Makefile, dev-bootstrap.sh, .env.example
  - [x] ADRs 0001-0015 written
  - [x] README.md
  - [x] Health-probe-service — boots, /health returns 200, /metrics returns Prometheus format
  - [x] All tests green (Go build, Go tests, TS typecheck, lint)
- [x] Part 3: Platform services (master-data, document, notification, audit, workflow, rules)
  - [x] master-data-service: CRUD for HSN codes, state codes, vendors, customers, items (port 8084, master_data_db)
  - [x] document-service: Upload with DEK envelope encryption (AES-256-GCM + KMS CMK), download/decrypt, OCR trigger, document lineage (port 8085, document_db)
  - [x] notification-service: Email send via SMTP (Mailpit in dev), digest consolidation (N→1 email), bounce processing with email_valid flag, user preferences (port 8086, notification_db)
  - [x] audit-service: Immutable event log, Merkle hash chain per hour bucket, integrity check with tamper detection (port 8087, audit_db)
  - [x] workflow-service: Workflow definitions/instances, human tasks, Temporal integration (NoopEngine fallback), signal-based task completion (port 8089, workflow_db)
  - [x] rules-engine-service: Tax determination (intra/inter-state CGST/SGST/IGST split), HSN validation with rate lookup, TDS applicability (6 sections), execution logging (port 8090, rules_engine_db)
  - [x] Postgres RLS enforced on all 6 new databases (complai_app role, SET LOCAL app.tenant_id)
  - [x] Outbox pattern verified: pending row → SQS delivery → published status
  - [x] Merkle chain integrity: 100 events → hash computed → tamper detected
  - [x] Notification digest: 5 notifications → 1 digest email via Mailpit
  - [x] SES bounce simulation: bounce → email_valid=false, bounce_count=1
  - [x] Unit test coverage: master-data 98.0%, document 90.1%, notification 89.9%, audit 98.9%, workflow 90.8%, rules-engine 98.3%
  - [x] Dockerfiles for all 6 services
  - [x] Goose migrations for all 6 databases

- [x] Part 4: API Gateway + BFF + Web Shell + design system components
  - [x] web-bff-service (Node/Express): aggregates backend services, CSRF, idempotency, tenant context propagation (port 3100)
  - [x] apps/web (Next.js 15 App Router): authenticated layout with Complai shell, sidebar with 6 groups, theme provider (15 themes), command palette
  - [x] packages/ui-components: Aura primitives + compliance-specific components (StatusBadge, GovStatusPill, FilingConfirmationModal, VendorComplianceScoreCard, MakerCheckerApprovalCard, etc.)
  - [x] Design tokens: CSS variables per Aura themes.ts, Tailwind config, density modes
  - [x] Storybook 8 with a11y addon
  - [x] E2E tests (Playwright): sidebar groups, theme switching, density, command palette, dashboard
- [x] Part 4.5: Scope correction — align with Bank Open ecosystem (Apex/Aura/Bridge/HRMS)
  - [x] PRD updated: removed modules 7-10 (AP Automation, Invoice Discounting, Complai One, Vendor Management), renumbered to 7 modules, added Bank Open ecosystem + Integration Layer sections
  - [x] Architecture updated: removed ap-service + billing-service, renamed vendor-service → vendor-compliance-service, added 4 sibling gateways, added canonical schemas (Payment, Contract, Payroll)
  - [x] Design system updated: sidebar 7→6 groups (removed Procurement/Payables, added Data Sources), module colors updated to 7 modules
  - [x] Build prompt updated: relabeled Parts 6, 7, 11, 13 for Bank Open alignment, updated E2E test suite
  - [x] Sidebar component updated: 6 groups with Data Sources, compliance items renamed
  - [x] Command palette + header updated: removed procurement/payables commands, added compliance + data-sources
  - [x] Data Sources route pages created: connected-apps, sync-status, ar-invoices, ap-invoices, vendors, contracts, payroll
  - [x] E2E tests updated: 6 groups assertion, compliance group collapse test
  - [x] CLAUDE.md + BUILD_PLAN.md updated

## Bank Open ecosystem
Complai is one of four sibling apps in the Bank Open family:
- **Apex P2P** — owns vendor master, AP invoices, POs, GRNs, payments (in UAT)
- **Aura O2C** — owns customer master, AR invoices, collections (early stage)
- **Bridge** — owns contracts, obligations, renewals (early stage)
- **HRMS** — external, payroll data + Form 16

Complai consumes from siblings via 4 gateway services (apex-gateway, aura-gateway, bridge-gateway, hrms-gateway). Phase 1 uses mock data sources; Phase 2 connects to real sibling APIs.

## Out of scope (moved to siblings)
- **AP Automation** — Apex P2P owns invoice ingestion, OCR, 3-way match, approval workflows, payment file generation
- **Invoice Discounting** — Apex P2P will integrate with TReDS
- **Complai One (SMB billing)** — deferred; Aura O2C covers SMB invoicing
- **Vendor Management (CRUD)** — Apex P2P owns vendor master; Complai only does compliance scoring on synced data
- **Vendor Portal** — Apex P2P provides vendor self-service portal
- **Portal BFF + SMB BFF** — removed (no vendor-portal or complai-one apps)

## Blockers / credentials needed
- AWS account with IAM roles + ap-south-1/ap-south-2 enabled (Part 1)
- Cloudflare account (Pro or Business plan) (Part 1 / Part 14)
- GitHub organization for complai repo (Part 1)
- Adaequare sandbox credentials + GST Returns Enriched API doc (Part 5)
- Sandbox.co.in sandbox api_key + api_secret (Part 6)
- Temporal Cloud account (Part 3)
- Last9 account (Part 3 / Part 14)
- Amazon SES production access (1-3 business day request) (Part 3)
- Azure OpenAI subscription + private endpoint (Part 12)
- Apex P2P UAT API access + webhook config (Part 6 mock, Part 13 real)
- Aura O2C API access (Part 11 mock, Part 13 real)
- Bridge API access (Part 11 mock, Part 13 real)
- HRMS API access (Part 11 mock, Part 13 real)
- MCA21 test access (Part 13)

## DevOps handoff

Items that require real AWS/cloud access. The DevOps team should execute these when engaged.

### Infrastructure provisioning (Part 1+)
- [ ] Run `terraform init` and `terraform plan` for all environments in `infra/terraform/`
- [ ] Apply Terraform for dev/sandbox environment first (`infra/terraform/environments/dev/`)
- [ ] Provision VPC, subnets, NAT Gateway, VPC endpoints in ap-south-1
- [ ] Provision EKS cluster (`complai-dev-ap-south-1`) with node groups: system, application, batch
- [ ] Provision RDS PostgreSQL 16 Multi-AZ + read replica
- [ ] Provision ElastiCache Redis 7 cluster-mode
- [ ] Provision OpenSearch 2 cluster (3-node dev)
- [ ] Provision S3 buckets (documents, uploads, exports, backups) with lifecycle rules
- [ ] Provision SQS queues + SNS topics per architecture §4.5
- [ ] Provision Secrets Manager secrets + KMS CMKs
- [ ] Configure IAM roles + IRSA for EKS service accounts
- [ ] Set up ECR repositories for all service images

### Networking & CDN (Part 1 / Part 4)
- [ ] Configure Cloudflare DNS for `complai.in`, `app.complai.in`, `api.complai.in`
- [ ] Set up Cloudflare WAF rules (OWASP Top 10, bot management)
- [ ] Provision ALB in public subnet, restrict to Cloudflare IP ranges
- [ ] Configure TLS: Cloudflare edge termination + origin-pull to ALB
- [ ] Set up NAT Gateway Elastic IPs and whitelist with Adaequare/Sandbox
- [ ] Install Istio ambient mesh 1.22 on EKS

### Email (Part 3)
- [ ] Request Amazon SES production access (1-3 business day approval)
- [ ] Configure SES sending domains: `complai.in`, `notifications.complai.in`
- [ ] Set up DKIM + SPF + DMARC for all sending domains
- [ ] Configure SES inbound email for notification replies (`*@replies.complai.in`)

### Observability (Part 3 / Part 14)
- [ ] Create Last9 account and obtain OTLP endpoint
- [ ] Configure OTel Collector DaemonSet to export to Last9
- [ ] Set up Last9 dashboards: per-service RED, per-provider, per-tenant, SLO
- [ ] Configure Last9 alerting → PagerDuty + Slack

### Workflow orchestration (Part 3)
- [ ] Create Temporal Cloud account (AWS Mumbai region)
- [ ] Provision namespaces: complai-filings, complai-reconciliation, complai-bulk, complai-sync, complai-onboarding
- [ ] Configure Temporal Cloud mTLS certificates

### CI/CD (Part 1)
- [ ] Set up GitHub organization and repository
- [ ] Configure GitHub Actions runners with AWS access
- [ ] Set up ArgoCD on EKS for GitOps deployment
- [ ] Configure ECR push permissions for CI
- [ ] Set up Trivy, Snyk, GitLeaks, Semgrep in CI pipeline

### DR setup (Part 14)
- [ ] Apply Terraform for ap-south-2 (Hyderabad) warm standby
- [ ] Provision EKS cluster in DR region (1 node per group)
- [ ] Set up RDS cross-region read replica (promotable)
- [ ] Enable S3 cross-region replication
- [ ] Configure Cloudflare / Route 53 health-check failover
- [ ] Run DR drill: failover → verify filings work → failback

### Security hardening (Part 14)
- [ ] Enable AWS GuardDuty + Security Hub
- [ ] Configure Secrets Manager auto-rotation (90-day cycle)
- [ ] Set up VPC Flow Logs → S3
- [ ] Run OWASP ASVS L2 checklist
- [ ] Enable AWS CloudTrail for Secrets Manager access audit

### Provider onboarding
- [ ] Adaequare: sign contract, receive sandbox credentials, whitelist NAT EIPs, test e-Invoice + EWB (Part 5)
- [ ] Sandbox.co.in: create account, subscribe to TDS + IT + KYC + Tax Payment, receive sandbox API key (Part 6)
- [ ] MCA21 V3: obtain test access for ROC form filing (Part 13)

### Sibling app onboarding (Bank Open ecosystem)
- [ ] Apex P2P: obtain UAT API credentials, configure webhook endpoints for vendor + AP invoice sync (Part 13)
- [ ] Aura O2C: obtain API credentials, configure webhook endpoints for AR invoice sync + IRN/EWB status pushback (Part 13)
- [ ] Bridge: obtain API credentials, configure webhook for contract sync (Part 13)
- [ ] HRMS: obtain API credentials, configure webhook for payroll + Form 16 sync (Part 13)

## Notes
- Docs in /docs/input are authoritative; don't re-derive
- ADRs in /docs/adr are authoritative for past decisions
- Every part ends with: tests green + BUILD_PLAN updated + commit
- All dev uses LocalStack — same code runs on real AWS via env var swap
- Terraform files are scaffolding for DevOps team — never run locally
