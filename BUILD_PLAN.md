# BUILD_PLAN.md — Living checklist

Last updated: 2026-04-26

## Current part
Part 7 complete (Reconciliation engine + GSTR-3B + GSTR-2B/IMS). Next = Part 8.

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

- [x] Part 5: Adaequare GST gateway + GSTR-1 flow
  - [x] gstn-gateway-service: Mock Adaequare Enriched GST APIs — Save, Submit, File GSTR-1 (port 8091)
  - [x] aura-gateway-service: Stubbed AR invoice data — 100 mock invoices (B2B intra/inter, B2CL, B2CS, export, NIL-rated, CRN/DBN) with type mix (port 8092)
  - [x] gst-service: GSTR-1 domain logic — 11-section categorization, validation, filing lifecycle (ingest→validate→review→approve→file→acknowledge), maker-checker, outbox pattern (port 8093)
  - [x] Categorizer engine: B2B, B2CL, B2CS, CDNR, CDNUR, EXP, NIL, HSN, DOCS, AT, ATADJ — 14 tests, 100% coverage
  - [x] Filing wizard UI: 6-step wizard (Next.js), step indicator, section review table (7 columns), filing confirmation modal (type-to-confirm, DSC/EVC selector)
  - [x] GST Returns landing page with cards for GSTR-1, GSTR-3B, GSTR-2B, GSTR-9
  - [x] Postgres RLS enforced on gst_db (sales_register, gstr1_filings, gstr1_sections, validation_errors, outbox)
  - [x] Goose migrations for gst_db (5 tables, RLS policies, unique constraints)
  - [x] Unit test coverage: gst-service handlers 84.3%, categorizer 100%, aura-gateway 100%, gstn-gateway 100%
  - [x] TypeScript typecheck clean (tsc --noEmit)
  - [x] go.work updated with 3 new services
  - [x] Dockerfiles for all 3 services
  - [x] Step-up auth gate: StepUpVerifier interface, File handler blocks without valid step-up (403 step_up_required), 2 tests PASS
  - [x] Maker-checker enforcement: CreatedBy tracked on filing, self-approval denied (403 self_approval_denied), different user can approve, 2 tests PASS
  - [x] Playwright E2E: Full wizard lifecycle (Ingest→Validate→Review→Approve→File with Cancel+type FILE+EVC→Acknowledge with ARN), 1 test PASS

- [x] Part 6: Sandbox KYC gateway + Vendor Compliance + Apex Sync
  - [x] kyc-gateway-service: Mock Sandbox.co.in KYC APIs — PAN, GSTIN, TAN, Bank verification (port 8094)
    - [x] PAN: format validation (10-char alphanumeric), entity type detection (Company/Individual/HUF/Trust/Firm/AOP/BOI/Government/AJP/Local Authority)
    - [x] GSTIN: 15-char validation, state code lookup (37 states), PAN extraction, status + registration type
    - [x] TAN: 10-char validation, deductor name generation
    - [x] Bank: IFSC validation (11-char format), 6-bank mapping (SBI/HDFC/ICICI/Axis/Kotak/PNB), account verification
    - [x] Unit test coverage: API 83.2%, provider 100% (43 tests total)
  - [x] apex-gateway-service: Mock Apex P2P client — 50 diverse vendor profiles + AP invoices (port 8095)
    - [x] 50 vendors across 4 compliance tiers: 10 exemplary (Cat A), 15 good (Cat B), 15 average (Cat C), 10 poor (Cat D)
    - [x] Indian enterprise names (Tata Steel, Infosys, Reliance, etc.), valid GSTIN format, 10 states, 5 categories
    - [x] 270+ AP invoices (3-9 per vendor), Oct 2025 - Mar 2026, 4 GST rates, CGST/SGST vs IGST
    - [x] Deterministic mock data (no randomness)
    - [x] Unit test coverage: API 86.7%, provider 96.3% (35 tests total)
  - [x] vendor-compliance-service: 100-point compliance scoring engine (port 8096, vendor_compliance_db)
    - [x] 5-dimension scorer: Filing Regularity (30pts), IRN Compliance (20pts), Mismatch Rate (20pts), Payment Behavior (15pts), Document Hygiene (15pts)
    - [x] Categories: A (≥90), B (60-89), C (40-59), D (<40). Risk levels: Low (≥80), Medium (60-79), High (40-59), Critical (<40)
    - [x] Sync endpoint: POST /sync → fetches 50 vendors from apex-gateway, scores all, persists to vendor_compliance_db
    - [x] Read-only API: GET /vendors, GET /vendors/{id}/score, GET /summary, GET /sync/status (no vendor CRUD)
    - [x] Postgres RLS enforced: vendor_snapshots, compliance_scores, sync_status tables
    - [x] Goose migration with CHECK constraints on score ranges and category/risk values
    - [x] Unit test coverage: API 82.7%, scorer 95.7% (33 tests total)
  - [x] Vendor compliance UI page (apps/web/src/app/compliance/vendor-compliance/page.tsx)
    - [x] List view: 50 mock vendors, KPI summary cards, category filter (A/B/C/D/All), search, vendor table
    - [x] Detail view: VendorComplianceScoreCard (10-dot bar + 5-dimension breakdown), vendor details, category badge
    - [x] Sync from Apex button with loading state
    - [x] data-testid attributes for Playwright testing
  - [x] Playwright E2E test (apps/web/e2e/vendor-compliance.spec.ts)
    - [x] List view renders with 50 vendors, KPI summary visible
    - [x] Category filter (A→10, D→10, All→50)
    - [x] Search filter ("Tata"→1 row)
    - [x] Click vendor → detail view with VendorComplianceScoreCard (5 dimensions visible)
    - [x] Back button returns to list, sync button interaction
  - [x] Infrastructure: vendor_compliance_db added to postgres-init.sh, go.work updated with 3 services, Makefile MIGRATE_ORDER updated
  - [x] Storybook: VendorComplianceScoreCard story exists, Storybook builds clean
  - [x] TypeScript typecheck clean (tsc --noEmit)
  - [x] go vet + go build clean on all 17 Go modules
  - [x] Dockerfiles for all 3 new services
  - [x] RLS regression verified: cross-tenant query returns 0 rows
  - [x] Read-only API verified: PUT/DELETE return 405
  - [x] End-to-end sync verified: 50 vendors synced + scored (CatA=6, CatB=27, CatC=11, CatD=6, Avg=69)

- [x] Part 7: Reconciliation engine + GSTR-3B + GSTR-2B/IMS (AP register from Apex)
  - [x] recon-service: 5-stage ITC reconciliation engine (port 8097, recon_db)
    - [x] Matcher pipeline: exact → fuzzy (Levenshtein) → amount → partial → residual
    - [x] 6 match types: Direct, Probable, Partial, Missing2B, MissingPR, Duplicate
    - [x] IMS round-trip: send ACCEPT/REJECT action to GSTN gateway, local action store, fetch IMS state
    - [x] Gateway integration: Apex (AP invoices, double-wrapped response), GSTN (GSTR-2B, IMS)
    - [x] Precision/recall: 1000×1000 synthetic benchmark — 100% precision, >95% recall
    - [x] Fuzzy match: Levenshtein distance ≤2, amount within 1% tolerance, 4 test scenarios
    - [x] Performance: 10K×10K reconciliation under 5 minutes (2.8s actual)
  - [x] gst-service: GSTR-3B auto-fill pipeline
    - [x] 3-source aggregation: GSTR-1 summary (outward supply), GSTR-2B (inward supply), IMS (ITC filtering)
    - [x] Computed fields: outward supply rows, inward supply, eligible ITC (IMS-filtered), gross liability, net liability
    - [x] IMS-aware ITC: rejected invoices excluded, RCM excluded, pending IMS flagged
    - [x] Filing lifecycle: populate → review → approve (maker-checker) → save → submit → file → acknowledge
    - [x] GSTR3BSummary, GSTR3BApprove (self-approval denied), GSTR3BFile (step-up, DSC/EVC) handlers
  - [x] ITC Reconciliation Workspace UI
    - [x] Bucket summary (matched/mismatch/missing2B/missingPR/duplicate)
    - [x] Reconciliation table with search, bucket filter, AI suggestions toggle
    - [x] Run reconciliation button, export button
  - [x] GSTR-3B Filing Wizard UI
    - [x] 6-step wizard: Auto-Populate → Review → Pay → Sign → File → Acknowledge
    - [x] Step indicator with completed/active/pending states
    - [x] Review tabs: Tax Liability (Tables 1-6), ITC (Tables 4A-D) with SourceBadge
    - [x] Filing confirmation modal: DSC/EVC selector, type-to-confirm, ARN receipt
  - [x] Postgres RLS enforced on recon_db (recon_runs, recon_matches, ims_actions)
  - [x] Goose migrations for recon_db (3 tables, RLS policies)
  - [x] Unit test coverage: recon-service handlers 88.5%, matcher 91.2%; gst-service handlers 81.1%, categorizer 100%
  - [x] Playwright E2E: Recon Workspace (6 tests), GSTR-3B Wizard full lifecycle (3 tests)
  - [x] TypeScript typecheck clean (tsc --noEmit)
  - [x] go vet + go build clean on all affected modules
  - [x] Storybook: build clean, 16/16 story tests pass, axe-core 0 a11y violations
  - [x] StepIndicator: data-testid added for E2E

### Deferred Part 5 hardening tests (→ Part 14)
- [ ] Idempotency E2E: duplicate ingest with same GSTIN+period returns same filing_id, no double-count
- [ ] Failure recovery: GSTN gateway returns 5xx mid-file → filing status = 'failed', retryable
- [ ] Filing Confirmation Modal Playwright: test DSC path, wrong confirm word rejected, modal a11y audit
- [ ] Full audit pipeline: filing lifecycle events → audit-service → Merkle chain verification
- [ ] Outbox delivery: outbox row created on file → SQS delivery → published status
- [ ] RLS isolation E2E: tenant A cannot read tenant B's filings via API
- [ ] Concurrent filing: two users file same period simultaneously → only one succeeds
- [ ] Wire wizard to real backend APIs (requires BFF proxy gst-service route)

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
- All 11 service databases auto-provisioned via `scripts/postgres-init.sh`; migrations applied via `make migrate-all` (dependency-ordered, stops on failure)
