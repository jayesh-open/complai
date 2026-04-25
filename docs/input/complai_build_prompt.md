# Complai — Claude Code Build Prompt

**Version:** 1.0
**Status:** Approved for build
**Target:** Claude Code (CLI agent)
**Product:** Complai — Enterprise GST + TDS + Compliance SaaS for Indian enterprises

**Input documents (must be in `/docs/input/` before Part 1):**
- `complai_prd.md` — product requirements
- `complai_architecture.md` — technical architecture
- `complai_design_system.md` — design system (inherits from Aura)
- `complai_api_integration.md` — Adaequare + Sandbox.co.in integration spec
- `AURA_DESIGN_SYSTEM.md` — foundation design system

**Invocation model:** paste each numbered part into Claude Code one at a time. After Part N finishes and tests are green, invoke Part N+1.

---

## PART 0 — Meta-prompt (paste into `CLAUDE.md`; paste at the start of every new session)

You are building **Complai**, an enterprise GST + TDS + compliance SaaS platform for Indian enterprises. Before writing any code in any session, do these four things:

1. **Read `/CLAUDE.md`** — session-durable memory (invariants, stack, current state).
2. **Read `/BUILD_PLAN.md`** — living checklist. The last unchecked item is current work.
3. **Read the input docs** listed for the current part.
4. **Announce the part + Definition of Done (DoD) back to the user** before writing any code.

### Non-negotiable operating rules (every part, every session)

- **Multi-tenancy is an invariant.** Every table has `tenant_id` + Postgres Row-Level Security (RLS). Every service method asserts it. No exceptions.
- **Enriched APIs only** for Adaequare. No SEK, no AES, no HMAC code in our repo.
- **Provider-agnostic internal contracts.** Callers of `gstn-gateway`, `tds-gateway`, etc. never see Adaequare/Sandbox specifics.
- **Design-system rules win.** Every UI component reads theme CSS variables. No hex colors in component code. Default theme = Light Classic. Default density = Compact.
- **DD/MM/YYYY dates, ₹ prefix, Indian numbering.** Use shared formatters from `packages/shared-kernel-node`.
- **Idempotency on every external call.** UUID `request_id` in the outbox row; providers de-dupe on it.
- **Outbox pattern for every state-changing external call.** Never call providers directly from business services — always via the gateway, always via the outbox.
- **Observability first.** Every service emits OpenTelemetry traces with `tenant_id` + `gstin/tan/pan` as baggage. No silently-swallowed exceptions.
- **Tests before DoD.** Every part has tests. DoD is not green until they pass.
- **No secrets in code.** AWS Secrets Manager only. Reject any PR with hardcoded credentials.
- **Conventional commits**, one commit per part minimum: `feat(module): ...`, `test: ...`, `docs: ...`.
- **No framework drift.** Every Go service uses: chi + pgx + sqlc + goose + fx + viper + zerolog + OTel + testify + gomock + shopspring/decimal + aws-sdk-go-v2 + gobreaker. Deviation requires an ADR.

### When stuck

- If a decision isn't obvious from the docs, write an ADR at `/docs/adr/NNNN-title.md`, make the call, proceed.
- If a credential isn't yet available (Adaequare, Sandbox, AWS, Cloudflare, SES, Last9), stub with a mock provider, note in `BUILD_PLAN.md` under "Credentials Needed", proceed.
- If truly blocked, stop, write the blocker to `BUILD_PLAN.md`, commit progress, tell the user.

### When you finish a part

- Run the part's tests. All must pass.
- Update `BUILD_PLAN.md` (check off items; add new items discovered).
- Update `CLAUDE.md` if architectural facts changed.
- Commit.
- Announce "DoD achieved" + one-paragraph hand-off summary for the next part.
- **Stop.** Do not start the next part without explicit user invocation.

---

## PART 0.5 — Initialize repo + memory scaffolding (~30 min, one-time)

**Goal:** Create the repo, meta-memory files, input-doc folder, ADRs folder. After this, any session can resume from `/CLAUDE.md`.

**Tasks:**

1. `git init` in `~/workspace/complai/`.
2. Create `/docs/input/` and `/docs/adr/` directories.
3. Place the 5 input documents in `/docs/input/`.
4. Create `/CLAUDE.md` with this content:

```markdown
# CLAUDE.md — Session-durable memory for Complai

## What we're building
Complai is an enterprise GST + TDS + compliance SaaS for Indian enterprises. Target: 500 Year-1, 50K GSTINs Year-3. Competitive reference: Clear / ClearTax.

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
- services/go/{name}-service — Go services (domain + gateways)
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

## Current build state
- [ ] Part 1: Repo skeleton + shared foundation
- [ ] Part 2: Identity + Tenant + User/Role services + auth
- [ ] Part 3: Platform services (master-data, document, notification, audit, workflow, rules)
- [ ] Part 4: API Gateway + BFF + Web Shell + design system components
- [ ] Part 5: Adaequare GST gateway + GSTR-1 flow
- [ ] Part 6: Sandbox KYC gateway + Vendor Compliance + Apex Sync
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
(populated as written — see /docs/adr/)
```

5. Create `/BUILD_PLAN.md`:

```markdown
# BUILD_PLAN.md — Living checklist

Last updated: (date)

## Current part
None yet. Next = Part 1.

## Completed
None yet.

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
- MCA21 test access (Part 13)
- SAP/Oracle test systems (Part 13)
- HDFC/ICICI sandbox bank APIs (Part 13)

## Notes
- Docs in /docs/input are authoritative; don't re-derive
- ADRs in /docs/adr are authoritative for past decisions
- Every part ends with: tests green + BUILD_PLAN updated + commit
```

6. Create `/docs/adr/README.md` with ADR format template.
7. Initial commit: `chore: init complai repo with input docs and memory scaffolding`.

**DoD:**
- Clean `git log` with one commit.
- `/CLAUDE.md`, `/BUILD_PLAN.md`, `/docs/input/*.md`, `/docs/adr/README.md` all present.
- `tree -L 2` shows the initial structure.

---

## PART 1 — Repo skeleton + shared foundation (~6h)

**Goal:** Monorepo structure, shared libraries, AWS Terraform scaffolding, dev environment. Nothing runs yet except a hello-world per runtime.

**Inputs:** `/CLAUDE.md`, `/BUILD_PLAN.md`, architecture §3, §12, §13.

**Tasks:**

1. **Monorepo scaffolding** — `package.json` with pnpm workspaces + Turborepo; `pnpm-workspace.yaml` listing apps/services/packages; `go.work` listing all Go services + `packages/shared-kernel-go`; `.editorconfig`, `.gitignore`, `.nvmrc` (Node 20.18).

2. **Shared Go library (`packages/shared-kernel-go`)** — tenant (JWT middleware, WithTenantContext), outbox (OutboxRow + PublishInTx), messagebus (SQS/SNS interface + impl), auth (JWT validation), errors (domain taxonomy), formatters (INR, DD/MM/YYYY, GSTIN/PAN/TAN), otel (setup with Last9 exporter), config (Viper + Secrets Manager), http (chi middleware: tenant, tracing, logging, rate limit), db (pgx pool, goose runner), testutil (testcontainers helpers).

3. **Shared Node library (`packages/shared-kernel-node`)** — TS types shared with frontend; Zod schemas (Tenant, User, Vendor, Invoice); formatters mirroring Go versions; API client factory with auth + tracing; error taxonomy mirroring Go.

4. **UI library scaffolding (`packages/ui-components`)** — exports only; primitives come in Part 4.

5. **Events library (`packages/events`)** — Protobuf schemas: `TenantCreated`, `UserCreated`, `InvoiceCreated`, `FilingSubmitted`, `FilingAcknowledged`; Canonical Invoice Schema at `events/schemas/canonical-invoice.proto`; `buf.yaml` + codegen scripts (Go + TS outputs).

6. **OpenAPI library (`packages/openapi`)** — `specs/` per-service OpenAPI 3.1 files; `oapi-codegen` config for Go types; `openapi-typescript` for TS.

7. **Terraform (`infra/terraform/`)** — AWS provider; ap-south-1 + ap-south-2; modules: vpc, eks, rds-postgres, elasticache-redis, opensearch, s3, sqs-sns, secrets-manager, kms, acm, alb, cloudflare-dns, iam-irsa; environments: dev, sandbox, staging, prod; workspaces per env; state in S3 + DynamoDB locking. Ready to plan, not yet applied.

8. **Helm charts (`infra/helm/`)** — base chart (Istio ambient + External Secrets Operator + ArgoCD + cert-manager); service chart template.

9. **Dev environment (`docker-compose.dev.yml`)** — Postgres 16 (with uuid-ossp, pgcrypto, pg_trgm), Redis 7, LocalStack (SQS, SNS, S3, Secrets Manager, KMS emulation), OpenSearch 2 (dev mode), Keycloak 24 (with dev realm import), Temporal dev server, Mailpit, OTel Collector + Jaeger. All on `complai-dev` network with healthchecks. ~4 GB RAM footprint.

10. **CI/CD (`.github/workflows/`)** — `ci.yml` (per-service matrix build, Go + pnpm caches, test gates, Docker + ECR push, Trivy scan, SBOM); `deploy-staging.yml` (ArgoCD sync); `security.yml` (Trivy + Snyk + GitLeaks + Semgrep); commitlint.

11. **Dev tooling** — `Makefile` (dev, db-up, migrate, seed, lint, test, build, clean, down); `scripts/dev-bootstrap.sh`; pre-commit hooks (lint-staged, commitlint, detect-secrets, gofmt, golangci-lint).

12. **Docs** — `README.md` with overview + local-run; write ADRs 0001-0020 as enumerated in architecture §18.

13. **Dummy service** — `services/go/health-probe-service/` — minimal chi server with `/health` and `/metrics`, demonstrates shared-kernel-go usage.

**Tests:**

- `make dev` brings up docker-compose; all containers healthy within 120s.
- `make lint` passes across all workspaces.
- `pnpm -r run typecheck` passes.
- `go work sync && go build ./...` passes.
- `terraform -chdir=infra/terraform validate` passes for all environments.
- `health-probe-service` boots, `/health` returns 200, emits OTel span visible in Jaeger.
- Unit tests for shared-kernel-go helpers pass (tenant context, outbox helper, message bus).

**DoD:**
- Repo structure matches architecture §13.
- All tests green.
- BUILD_PLAN updated.
- Committed: `feat(foundation): repo skeleton, shared libs, dev env, adrs`.

---

## PART 2 — Identity, Tenant, User/Role + auth (~7h)

**Goal:** Three core platform services built. End-to-end login via Keycloak. JWT carries `tenant_id` + role claims.

**Inputs:** Architecture §4 (multi-tenancy), §8 (identity/authz); PRD §16.1.

**Tasks:**

1. **`identity-service` (Go)** — Postgres schema with RLS: `users`, `user_credentials` (refs Keycloak), `user_sessions`, `mfa_factors`, `step_up_events`. Endpoints: `/login`, `/logout`, `/refresh`, `/mfa/enroll`, `/mfa/verify`, `/password/reset`, `/sso/callback` (OIDC for Google Workspace, Azure AD, Okta). Step-up auth: re-authenticate within 5-min window for `action=FILE_*`.

2. **`tenant-service` (Go)** — Schema: `tenants`, `tenant_pans`, `tenant_gstins`, `tenant_tans`, `tenant_settings`, `tenant_feature_flags`. Lifecycle: create (platform-admin only), suspend, reactivate, delete (soft + hard). Tenancy tier enum: `pooled | bridge | silo | on_prem`. PAN hierarchy: tenant → PANs → GSTINs → TANs.

3. **`user-role-service` (Go)** — Schema: `roles`, `permissions`, `role_permissions`, `user_roles`, `role_templates`, `approval_workflows`. Role templates: Tenant Admin, Tax Manager, Tax Analyst, AP Clerk, AP Approver, Viewer, Vendor, CA, Platform Admin. Permissions as resource-action pairs. Maker-checker: N-of-M chains per resource. Policy eval: `POST /policy/check` with `(user_id, resource, action)` → `{allow: bool, reasons: []string}`.

4. **RLS middleware** — in `packages/shared-kernel-go/tenant/`. Every table gets RLS policies: `USING (tenant_id = current_setting('app.tenant_id')::uuid)`. `WithTenantContext` enforces `SET LOCAL app.tenant_id`.

5. **JWT shape (ADR-0004)** — OIDC standard + `tenant_id`, `active_pan`, `active_gstin`, `active_tan`, `roles[]`, `permissions[]`, `step_up_at`. Signed by Keycloak; verified at Istio + service level.

6. **Seed data** — `scripts/seed-dev.sh` creates 1 platform tenant, 2 customer tenants, 5 users per tenant across role templates; loads Keycloak realm configs.

**Tests:**

- **Integration** — register → email verified (Mailpit) → login via Keycloak → JWT has `tenant_id` → access protected endpoint.
- **RLS test (critical)** — user from tenant A cannot see ANY tenant B data, even via direct SQL. Lynchpin test — if this fails, stop everything.
- **Maker-checker** — tax-analyst creates filing intent → `pending_approval` → tax-manager approves → `approved`; analyst cannot self-approve.
- **Step-up** — valid JWT but no recent step-up attempts FILE action → 403 + step-up-required → provides MFA → proceeds.
- **Role templates** — assigning "Tax Analyst" grants expected permission set.

**DoD:** login works in browser. All tests green. Seed data loads. Committed: `feat(identity): auth, tenancy, rbac, rls`.

---

## PART 3 — Platform services: Master-Data, Document, Notification, Audit, Workflow, Rules (~8h)

**Goal:** Six platform services every domain module depends on.

**Inputs:** Architecture §3.1, §7.2 (outbox); PRD §16.

**Tasks:**

1. **`master-data-service` (Go)** — Postgres + Redis cache. Tables: `vendors`, `customers`, `employees`, `items`, `chart_of_accounts`, `hsn_codes`, `pincodes`, `state_codes`, `bank_branches`. CSV import/export. RLS. Outbox → SNS `master-data.changes` topic. ERP sync stubs.

2. **`document-service` (Go)** — Postgres (metadata) + S3 (binaries). Upload with virus scan (ClamAV sidecar), OCR queue trigger, DEK envelope encryption per tenant via KMS. Pre-signed URL retrieval (15-min TTL). Document lineage. Retention policies per doc class.

3. **`notification-service` (Go)** — Channels: email (Amazon SES), SMS (MSG91), WhatsApp Business (Meta), in-app (SSE + Redis pub/sub). Template registry with Handlebars. Per-user preferences + quiet hours. Digest bundler (09:00 IST). DPDP consent tracking.

4. **`audit-service` (Go)** — OpenSearch (90d hot) + S3 archive. Every state-change event writes here. Merkle-chain integrity (hourly hash of previous hour → next hour's head; tamper-evident). Search API. Signed PDF export for regulatory submission.

5. **`workflow-service` (Go + Temporal Cloud client)** — Namespaces per environment. Workflow templates: `gstr1-filing-saga`, `gstr3b-filing-saga`, `tds-quarterly-saga`, `reconciliation-job`, `bulk-irn-generation`, `form16-bulk`, `vendor-onboarding`, `itr-saga`. Human-task integration.

6. **`rules-engine-service` (Go)** — Rules as versioned JSON per tenant. Categories: tax-determination, HSN-applicability, TDS-applicability, validation-rules, eligibility-rules. Evaluation: given canonical invoice → `{applicable_taxes[], tds_sections[], itc_eligibility, warnings[]}`. Seed with FY 2026-27 Indian tax rules.

7. **Outbox publisher** — for Phase 1, a Go outbox-poller goroutine in each service ships rows to SQS directly (simpler than Debezium + Kafka Connect). Poll every 500ms. ADR-0011 covers this.

**Tests:**

- **Document lifecycle** — upload PDF → virus-scanned → OCR queued → retrievable via pre-signed URL → KMS-encrypted at rest (round-trip).
- **Audit chain** — tamper with one event → Merkle check detects mismatch.
- **Outbox delivery** — kill SQS consumer → state change still commits + outbox queues → consumer recovers → downstream receives.
- **Notification digest** — 5 notifications within window → single digest email.
- **Rules engine** — sample B2B intra-state invoice → correct CGST/SGST split + HSN validation + TDS applicability.
- **Workflow** — sample saga with human_task interrupt → resumes after mock approval.

**DoD:** all six services running, ≥80% line coverage each, integration tests pass. Committed: `feat(platform): master-data, document, notification, audit, workflow, rules`.

> **Credential checkpoint:** Temporal Cloud + Last9 + Amazon SES production access by end of Part 3. Use local Temporal + Mailpit as fallback; note in BUILD_PLAN.

---

## PART 4 — API Gateway, BFF, Web Shell + Design System (~9h)

**Goal:** UI shell is live. User logs in, sees sidebar with 6 groups, navigates placeholder pages, toggles theme, sees Complai design system in Storybook.

**Inputs:** `complai_design_system.md` (whole); `AURA_DESIGN_SYSTEM.md`; architecture §3.5, §11.

**Tasks:**

1. **Cloudflare + AWS ALB + Istio** — Terraform apply: ALB + Istio ingress; Istio JWT validation (Keycloak public keys); Cloudflare in front of ALB (DNS + WAF + TLS); ALB security group restricts to Cloudflare IP ranges.

2. **`web-bff-service` (Node/NestJS)** — aggregates backend services; REST (GraphQL deferred); session mgmt, CSRF tokens, idempotency; per-tenant context propagation.

3. **`apps/web` (Next.js 15 App Router)** — authenticated layout with Complai shell; sidebar with 6 groups (Dashboard / My Tasks / Inbox → Compliance → Insights → Data Sources → Documents → Configure), collapsible, badge nudges, environment indicator, tenant/PAN/GSTIN selectors; theme provider (all 15 themes, Light Classic default); command palette (⌘K); routes rendering placeholders including `/data-sources/*` for Bank Open sibling sync views.

4. **`packages/ui-components` — Complai primitives** — Aura primitives (Button, Input, Select, Checkbox, Radio, Switch, Textarea, DatePicker, NumberInput, Card, Modal, Drawer, Toast, Tooltip, Popover, Skeleton, EmptyState, Tab, Breadcrumb); compliance-specific per design system §5 (StatusBadge, GovStatusPill, PeriodSelector, AuditTrailTimeline, KpiMetricCard, DataTable, FilingConfirmationModal, VendorComplianceScoreCard (skeleton), BulkOperationTray, MakerCheckerApprovalCard (skeleton), ReconciliationSplitPane (skeleton)); Storybook 8 with a11y addon.

5. **Design tokens** — CSS variables per Aura `themes.ts` verbatim; Tailwind config with `app-*` and `foreground-*` utilities; density modes via `data-density` attribute.

6. **Authenticated flows** — login → Keycloak → redirect with JWT → dashboard; logout → clear session → redirect; "Environment: Sandbox" banner on non-prod; tenant switcher functional.

**Tests:**

- **E2E (Playwright)** — new user logs in → dashboard with 6 sidebar groups → switches to Ocean theme → persists across refresh → logs out.
- **Accessibility** — every Storybook story passes axe-core (0 violations).
- **Visual regression (Chromatic)** — baselines for all primitives in 3 themes.
- **Density** — DataTable renders at correct row heights in all three modes.
- **Command palette** — ⌘K opens, fuzzy search filters, Enter navigates.
- **GovStatusPill** — renders correctly for GSTN/IRP/EWB/TRACES/MCA in success/warning/danger.

**DoD:** login → dashboard E2E works; sidebar 6 groups per design §2.1; Light Classic default, all themes switchable; Storybook published. Committed: `feat(ui): shell, design system, bff, cloudflare + alb + istio`.

---

## PART 5 — Adaequare GST gateway + GSTR-1 filing wizard (~11h)

**Goal:** First vertical slice of real compliance. Tenant uploads sales register, reviews derived GSTR-1, files it (Adaequare sandbox).

**Inputs:** `complai_api_integration.md` §3 (Adaequare GST); PRD §4; design system §4.2 (Filing Wizard).

**Tasks:**

1. **`gstn-gateway-service` (Go)** — talks to Adaequare Enriched GST APIs; Adaequare bearer cache in Redis (refresh at 23h); tenant GST credential vault at AWS Secrets Manager `complai/tenant/{tenant_id}/{gstin}/adaequare-credentials`; internal contract `POST /v1/gateway/adaequare/gstr1/{action}` with `X-Tenant-Id` + `Idempotency-Key`; actions: save, get, submit, file, status, reset, amendment, 1A get, 1A action; SQS consumer on `gov.outbound.gstn.queue`; emits to SNS `gov.outcome.gstn.topic`; circuit breaker (gobreaker).

2. **`gst-service` (Go)** — Schema: `gstr1_drafts`, `gstr1_sections`, `gstr1_invoices`, `gstr1_filings`, `gstr1_amendments`. Orchestrates 11 sections: b2b, b2cl, b2cs, cdnr, cdnur, exp, at, atadj, nil, hsn, docs. Sales register ingestion: CSV/Excel → CIS transform → section categorization. Rules engine integration. Temporal workflow `GSTR1FilingSaga`.

3. **UI: GSTR-1 Filing Wizard** (design system §4.2) — Steps: Ingest → Validate → Review → Pay (N/A) → File → Acknowledge; sticky step indicator; autosave every 10s; section review tables as government-form replicas; filing confirmation modal with type-to-confirm for tax > ₹10L; DSC/EVC selector; live progress view; GovStatusPill throughout; block navigation on unsaved changes.

4. **Sales register ingestion** — CSV/Excel uploader with column mapper (paste-from-Excel supported); named mapping templates per tenant; duplicate detection; preview + confirm before persist.

5. **Public GSTIN validation** — via Sandbox `/kyc/gstin` (preferred) or Adaequare `searchTP` (fallback); bulk endpoint up to 5K GSTINs per call.

**Tests:**

- **Unit** — rules engine categorizes 50 sample invoices into correct sections; tax computation accurate to ₹0.01 (shopspring/decimal).
- **Integration** — 1,000-invoice CSV ingested → sections auto-populated → validation passes → summary matches hand-calculated totals.
- **Gateway (Adaequare sandbox)** — save → get → reset → save → submit → file with EVC.
- **E2E (Playwright)** — login → Compliance → GST → period → upload CSV → review wizard → file → ARN on Acknowledge.
- **Idempotency** — same save `request_id` twice → same response, no duplicate.
- **Failure recovery** — kill gstn-gateway mid-filing → restart → SQS re-delivers → saga resumes → success.
- **Visual regression** — wizard at each step matches design-system spec.

**DoD:** file GSTR-1 end-to-end in sandbox; ARN received + shown; audit trail complete. Committed: `feat(gst): gstr-1 filing wizard + adaequare gateway`.

> **Credential checkpoint:** Adaequare sandbox credentials + GST Returns Enriched API doc. Mock gateway if not available.

---

## PART 6 — Sandbox KYC gateway + Vendor Compliance + Apex Sync (~8h)

**Goal:** Vendor compliance scoring on vendor master synced from Apex P2P. KYC validation via Sandbox.co.in. No vendor CRUD in Complai — Apex owns the vendor master.

**Inputs:** `complai_api_integration.md` §4.4; PRD §4 (Module 4); design system §5.8.

**Tasks:**

1. **`kyc-gateway-service` (Go)** — Sandbox KYC API; auth: api_key + api_secret → bearer cached in Redis; internal contract `POST /v1/gateway/sandbox/kyc/{action}`; actions: pan-verify, aadhaar-otp, aadhaar-verify, gstin-verify, tan-verify, bank-account-verify, mca-company, mca-director, udyam-verify, digilocker-session; response caching (GSTIN 7d, PAN 30d, TAN 30d, HSN 90d, bank 90d).

2. **`apex-gateway-service` (Go)** — Consumes vendor master, AP invoices, payments, POs, GRNs from Apex P2P. Phase 1: mock data source with realistic sample data. Internal contract `POST /v1/gateway/apex/{resource}/{action}`. Webhook receiver for real-time sync (Phase 2). Publishes `VendorSynced.topic` and `InvoiceSynced.topic` via SNS.

3. **`vendor-compliance-service` (Go)** — Read-only vendor store synced from Apex (no vendor CRUD). Compliance scoring 0-100: filing regularity (30) + IRN compliance (20) + mismatch rate (20) + payment behavior (15) + document hygiene (15). Category A/B/C/D. KYC enrichment via Sandbox. Compliance check via Adaequare `getReturnTrack`. MaxITC orchestration stub (completed in Part 12). Temporal workflow: `vendor-compliance-sync`.

4. **UI: Vendor Compliance** — Vendor compliance list (Workflow List pattern): KPI cards (total vendors, high-risk count, avg score), filters by category/score/state; vendor compliance detail with tabs (Compliance Score, Filing History heatmap, ITC Impact, Documents, Audit Trail); Data Sources > Imported Vendors view.

**Tests:**

- **Sandbox integration** — PAN verify cached first call → second call within TTL hits cache.
- **Apex sync** — mock Apex publishes 100 vendors → vendor-compliance-service receives → stores read-only copies.
- **Score accuracy** — 100% on-time + no mismatches → ≥90; late filings → <70.
- **No CRUD** — POST/PUT/DELETE vendor endpoints return 405; only Apex sync creates vendors.
- **Data Sources UI** — Imported Vendors page shows synced vendors with sync timestamp.

**DoD:** vendor compliance scoring works on Apex-synced vendors; KYC enrichment functional. Committed: `feat(vendor): compliance scoring, kyc gateway, apex sync`.

> **Credential checkpoint:** Sandbox.co.in sandbox api_key/secret. Mock pattern if missing. Apex gateway uses mock data source in Phase 1.

---

## PART 7 — Reconciliation Engine + GSTR-3B + GSTR-2B/IMS (~11h)

**Goal:** Full ITC cycle. Pull 2A/2B/IMS via Adaequare → reconcile vs purchase register (sourced from Apex AP invoices) → IMS actions → auto-populate GSTR-3B → file.

**Inputs:** PRD §4.2.5, §4.2.3; API spec §3.5 (IMS); design system §4.3.

**Tasks:**

1. **`recon-service` (Go)** — Consumes AP invoice register from Apex via `apex-gateway-service` as "purchase register" (Complai does not own AP invoices). Pull GSTR-2A section-by-section on schedule; GSTR-2B monthly; IMS + changes. 5-stage match pipeline: Exact → Fuzzy (Levenshtein on invoice_no, ±2d date, ±₹1 amount) → AI (Phase 12 stub) → Partial → Unmatched. Bucket persistence per period per GSTIN. IMS action sync via outbox → gstn-gateway → poll 2B regeneration → re-match.

2. **`gstr3b-service` (Go)** — Auto-populate from GSTR-1 (Tables 1-6), GSTR-2B/IMS (Tables 4A-D), ledgers. User override with justification. Liability offset (cash + credit ledger). Temporal: Prepare → Validate → Offset → Sign → File → Ack.

3. **UI: Reconciliation Workspace** (design §4.3) — split-pane (PR left, 2B/IMS right); clickable bucket counts; AI suggestion icon (stub); bulk "Accept all matched"; IMS action toolbar per row; reason-code popovers.

4. **UI: GSTR-3B Wizard** — same wizard pattern; auto-populated with "From GSTR-1" / "From 2B" / "Computed" badges; offset screen with live ledger balances; type-to-confirm for all 3B > ₹10L.

5. **Scheduled jobs** (Temporal cron) — monthly 2A + 2B + IMS pull on 14th night; weekly incremental 2A; on-demand re-run.

**Tests:**

- **Match accuracy** — 1,000 PR + 1,000 2B with known match set → >95% precision, >90% recall.
- **Fuzzy** — "INV-001" vs "INV001" → matched with confidence flag.
- **IMS round-trip** — accept in UI → action to sandbox → 2B regenerates → recon refreshes.
- **GSTR-3B auto-fill** — fully-reconciled period → all fields correct.
- **E2E** — Jan 2026: GSTR-1 filed → 2B pulled → recon → 3B auto-fill → offset → file → ARN.
- **Performance** — 10K × 10K recon < 5 min.

**DoD:** full ITC cycle in sandbox; workspace matches design. Committed: `feat(gst): reconciliation, gstr-3b, ims`.

---

## PART 8 — e-Invoicing + E-Way Bill (~9h)

**Goal:** Generate IRNs and EWBs; bulk ops; cancel flows. All 15 e-Invoice + 24 EWB endpoints wired.

**Inputs:** API spec §3.3, §3.4; PRD §5, §6.

**Tasks:**

1. **`irp-gateway-service` (Go)** — all 15 Adaequare e-Invoice endpoints; bulk IRN with tray progress.

2. **`einvoice-service` (Go)** — Transformer CIS → e-Invoice schema v1.1 (12 blocks / ~150 fields); validations; IRN + signed JSON + signed QR persisted; IRN-in-PDF composer; cancellation (24h window, reasons 1-4).

3. **`ewb-gateway-service` (Go)** — all 24 EWB endpoints; multi-vehicle movement (MULTIVEHMOVINT/ADD/UPD).

4. **`ewb-service` (Go)** — EWB lifecycle: generate → update Part-B → extend → cancel; consolidated EWB; integration with e-Invoice (`genewbfromirn`); transporter management.

5. **UI: e-Invoicing + EWB consoles** — e-Invoice list (Workflow List); bulk IRN wizard; cancel flow with reason + type-to-confirm; EWB list with validity timer (green >24h, amber 4-24h, red <4h/expired); one-click "Generate EWB" on invoice detail; multi-vehicle UI; QR viewer.

**Tests:**

- Single IRN + signed JSON + QR; retrieval works.
- Bulk IRN: 1,000 invoices → tray progress → all have IRNs.
- Cancel within 24h succeeds; after 24h rejects.
- EWB from IRN: no duplicate payload.
- Part-B update, extend validity, multi-vehicle flows.
- Transformer: 20 canonical invoices → produce e-Invoice JSON validating against NIC v1.1 schema.

**DoD:** generate + view + cancel IRN/EWB works; bulk IRN 1K works. Committed: `feat(compliance): e-invoicing, e-way-bill`.

---

## PART 9 — Sandbox TDS gateway + TDS module (~11h)

**Goal:** Full TDS cycle. Deductee master → payments → challan → TXT → CSI → FVU → e-file → Form 16 bulk.

**Inputs:** API spec §4.2 (10-step job flow); PRD §8.

**Tasks:**

1. **`tds-gateway-service` (Go)** — wraps all Sandbox TDS APIs; job-based polling via Temporal (every 10s first min, 30s next 5min, 2min thereafter); OTP flow for CSI download.

2. **`tds-service` (Go)** — Schema: `deductors`, `deductees`, `tds_payments`, `tds_challans`, `tds_quarters`, `tds_returns`, `form16_batches`. TDS calc per transaction via Sandbox calculator. 206AB/206CCA check at payment time. Quarterly returns: 24Q (salary), 26Q (non-salary resident), 27Q (non-resident), 27EQ (TCS). 10-step flow as Temporal workflow. Form 16 bulk: generate → zip → email via SES → download tracker.

3. **UI: TDS module** — TDS Dashboard; Deductees list with 206AB results; Payment recording (manual + bulk); Return wizard (Filing Wizard): Prepare → Generate TXT → Download CSI (OTP) → Generate FVU → E-File → E-Verify → Acknowledge; Form 16 bulk generator.

4. **Tax payment UI** — TDS challan via Sandbox tax-payment; CIN capture + ledger link.

**Tests:**

- Calculator: 50 transactions across sections → correct deduction.
- 206AB: mix of specified/non → correct higher-rate detection.
- Full quarter (Sandbox sandbox): 100 deductees + 500 payments + 10 challans → TXT valid → CSI with OTP sim → FVU → E-File → E-Verify → RRR.
- Form 16 bulk: 500 employees → zip + individual downloads.
- Correction: simulated defective → correction flow → resubmit.

**DoD:** full 24Q cycle in sandbox; Form 16 bulk works. Committed: `feat(tds): full cycle via sandbox`.

---

## PART 10 — Sandbox ITR + GSTR-9/9C (~10h)

**Goal:** Annual filings. ITR for employees (bulk Form 16 → 3-min filing). GSTR-9 + 9C for enterprises.

**Inputs:** API spec §4.3, §3.5; PRD §9.

**Tasks:**

1. **`itd-gateway-service` (Go)** — all Sandbox ITR endpoints (tax-payer registration, prefill, validate, submit, ITR-V, status, e-verify); OCR endpoints for Form 16/26AS/AIS.

2. **`itr-service` (Go)** — Employer-side: bulk Form 16 → email + magic link → 3-min filing flow. CA-assisted: CA dashboard, bulk pre-fill review, submission. Forms ITR-1 through ITR-7. Regime comparator (old vs new).

3. **`gstr9-service` (Go)** — Annual recon against books; Table 8A bulk pull (paginated); GSTR-9C auditor sign-off; differences report.

4. **UI: ITR module** — Employer console ("Send ITR invites" after Form 16); employee mobile-first flow (magic link → pre-filled → review → consent → e-verify → done); CA console.

5. **UI: GSTR-9/9C** — Annual return wizard; side-by-side returns vs books; auditor collaboration for 9C.

**Tests:**

- Employee ITR: simulated employee + Form 16 + savings → <3 min → ITR-V.
- Employer bulk: 500 invited → 300 file → tracker accurate.
- GSTR-9 prep: full year of 1/3B → auto-populated annual matches computed.
- Table 8A: 50K line items paginated.

**DoD:** both filings in sandbox. Committed: `feat(itr-gstr9): annual filings`.

---

## PART 11 — Sibling Gateway Services (Aura, Bridge, HRMS) (~8h)

**Goal:** Wire remaining Bank Open sibling gateways. Complai consumes AR invoices from Aura, contracts from Bridge, payroll from HRMS. Phase 1: mock data sources with realistic sample data.

**Inputs:** PRD §11 (Integration Layer); architecture §10A.

**Tasks:**

1. **`aura-gateway-service` (Go)** — Consumes customer master and AR invoices from Aura O2C. Publishes back filed-IRN-status and EWB status to Aura. Phase 1: mock data source. Internal contract `POST /v1/gateway/aura/{resource}/{action}`. Publishes `InvoiceSynced.topic` (AR). Data Sources > Imported AR Invoices view.

2. **`bridge-gateway-service` (Go)** — Consumes contracts from Bridge for TDS section determination and secretarial obligations. Phase 1: mock data source. Internal contract `POST /v1/gateway/bridge/{resource}/{action}`. Canonical Contract Schema mapping.

3. **`hrms-gateway-service` (Go)** — Consumes payroll data and Form 16 from external HRMS. Phase 1: mock data source. Internal contract `POST /v1/gateway/hrms/{resource}/{action}`. Canonical Payroll Schema mapping. Used by TDS (24Q salary) and ITR (employer-side bulk filing).

4. **UI: Data Sources module** — Connected Apps page (status cards for Apex, Aura, Bridge, HRMS); Sync Status dashboard (last sync time, record counts, errors); Imported AR Invoices, Imported Contracts, Imported Payroll Data list pages.

5. **Canonical schema adapters** — Implement Canonical Payment Schema, Contract Schema, Payroll Schema transforms from each sibling's format.

**Tests:**

- **Aura sync** — mock Aura publishes 500 AR invoices → synced → available for GSTR-1.
- **Bridge sync** — mock Bridge publishes 50 contracts → synced → TDS section determination works.
- **HRMS sync** — mock HRMS publishes 200 payroll records → synced → available for 24Q.
- **Data Sources UI** — Connected Apps shows 4 sibling apps with status; Sync Status shows timestamps and counts.
- **Canonical transforms** — each gateway maps raw sibling data to canonical schema correctly.

**DoD:** all 4 sibling gateways operational with mock data; Data Sources UI functional. Committed: `feat(gateways): aura, bridge, hrms sibling sync`.

---

## PART 12 — AI Layer + MaxITC + LLM Vendor Comms (~8h)

**Goal:** ML-driven differentiators. Better match rates; automated vendor chase; recommendations.

**Inputs:** Architecture §3.4; PRD §7.

**Tasks:**

1. **`matching-ml-service` (Python)** — CatBoost model for invoice-match confidence; features: vendor, amount, date delta, invoice-number edit distance, HSN similarity, seasonal patterns; training: synthetic initial data + real tenant data (with consent) post-rollout; FastAPI serving called by recon-service for Stage 3.

2. **`llm-copilot-service` (Python, Azure OpenAI via private endpoint)** — vendor comms: chase emails for missing invoices, IRN mismatches; natural-language search; return-prep assistant; safety: PII redaction pre-LLM, prompt-injection defenses, per-tenant usage caps.

3. **`maxitc-orchestrator-service` (Go)** — Daily cross-module intelligence; identify ITC at risk, vendors to chase, deadlines; action items → CFO dashboard + My Tasks.

**Tests:**

- Match lift: ML adds 5-10pp recall vs fuzzy-only baseline on held-out set.
- PII redaction: PAN/Aadhaar/bank strings → redacted before LLM call.
- Vendor chase: generates polite English/Hindi email with invoice numbers.
- Usage caps: per-tenant LLM spend capped, exceeding → 429 with error.

**DoD:** ML match integrated; copilot live; MaxITC runs daily. Committed: `feat(ai): ml + copilot + maxitc`.

---

## PART 13 — Real Bank Open Sibling Sync + GL-Stream + Compliance Cloud (~10h)

**Goal:** Replace mock sibling gateways with real Apex/Aura/Bridge/HRMS integrations; Compliance Cloud (secretarial); GL stream back to siblings.

**Inputs:** PRD §11 (Integration Layer); architecture §10A.

**Tasks:**

1. **Real Apex sync (`apex-gateway-service`)** — replace mock data source with real Apex P2P API/webhook integration; vendor master bidirectional compliance status; AP invoice real-time sync; PO + GRN sync for recon context.

2. **Real Aura sync (`aura-gateway-service`)** — replace mock with real Aura O2C API; AR invoice sync; publish back filed-IRN-status + EWB status to Aura for their invoice lifecycle.

3. **Real Bridge sync (`bridge-gateway-service`)** — replace mock with real Bridge API; contract sync for TDS section determination + secretarial obligation tracking.

4. **Real HRMS sync (`hrms-gateway-service`)** — replace mock with real HRMS API; payroll + Form 16 sync for 24Q + ITR.

5. **`gl-stream-service` (Go)** — real-time GL posting from compliance actions (tax paid, ITC claimed, TDS deposited); double-entry integrity; push journals back to Apex/Aura via sibling gateways.

6. **`secretarial-service` (Go) — Compliance Cloud** — entity registry (companies, LLPs, directors, DINs); filing calendar (AOC-4, MGT-7, DIR-3 KYC, ADT-1, CHG-1); MCA21 V3 integration (direct); document mgmt; compliance health score per entity; consumes company structure from Bridge.

**Tests:**

- Apex real sync: vendor master changes in Apex → reflected in Complai within 60s.
- Aura round-trip: AR invoice synced → IRN generated → status pushed back to Aura.
- Bridge: contract synced → correct TDS section derived → secretarial obligations created.
- HRMS: payroll synced → 24Q pre-populated correctly.
- ROC: AOC-4 → XML → MCA21 test → SRN received.
- GL integrity: daily imbalance = 0 across synthetic week.

**DoD:** all 4 sibling gateways on real APIs; secretarial works; GL sync operational. Committed: `feat(integrations): real sibling sync, mca21, gl-stream`.

> **Credential checkpoint:** Apex/Aura/Bridge API access (UAT), HRMS API, MCA21 test. Stub if missing.

---

## PART 14 — Reporting, Observability, Production Hardening (~9h)

**Goal:** Production readiness. Measured, secured, DR-rehearsed.

**Inputs:** Architecture §9, §10, §17.

**Tasks:**

1. **Reporting** — CFO Dashboard per design §4.5; pre-canned reports (monthly GST summary, TDS summary, aging, vendor scores, exceptions); custom report builder (Postgres + materialized views); scheduled email delivery (SES); PDF + Excel + CSV exports.

2. **Last9 observability hardening** — all services shipping metrics/logs/traces via OTel collector; dashboards (per-service RED, per-provider, per-tenant usage); alerting (filing success rate, gateway errors, auth failures, RLS violations) → PagerDuty + Slack; SLO dashboards per tier.

3. **Security hardening** — Secrets Manager auto-rotation every 90d; CI scanning (Trivy, Snyk, GitLeaks, Semgrep); OWASP ASVS L2 checklist green; DPDP (consent mgmt, data export/delete, PII access audit); Cloudflare WAF tuned; AWS GuardDuty + Security Hub active; SOC 2 evidence collection begins.

4. **Performance hardening** — k6 load tests (filing peak — 10K concurrent users, 1K IRNs/sec, 500 recon runs); Postgres indexes reviewed; Redis cache hit rates >80%; Temporal worker scaling.

5. **DR / BC** — full AWS infra applied to `ap-south-2`; RDS cross-region replica; S3 cross-region replication; EKS warm (1 node per nodegroup); Cloudflare + DNS failover; RTO 60min, RPO 5min for T0; scripted DR drill; monthly backup restore test.

6. **Documentation** — runbooks per service; OpenAPI per service (published); user guides (Admin, Tax Manager, AP Clerk, Vendor, CA); developer onboarding ("how to add a new compliance module").

**Tests:**

- Load: 1,000 IRNs/sec for 30 min — P95 < 300ms, error rate < 0.1%.
- DR drill: failover → GSTR-1 save/retrieve works → flip back → no data loss.
- SLO: 30 days green on synthetic probes for T0 services.
- Security: axe-core a11y pass, ZAP baseline clean, Trivy no High/Critical.

**DoD:** tests green; DR drill successful; runbooks written; prod-readiness sign-off. BUILD_PLAN complete. Committed: `feat(prod): reporting, last9, hardening, dr`.

---

## CONSOLIDATED E2E TEST SUITE (runs after Part 14)

Complete platform acceptance test. Runs against fresh staging. ~4 hours end-to-end.

**Setup:** fresh Postgres, empty Keycloak, clean SQS, seeded master data.

### 1. Tenant Lifecycle
Platform admin creates "Acme Manufacturing" (3 PANs, 15 GSTINs, 5 TANs). Acme admin adds 10 users across role templates. Permission matrix validated. Adaequare + Sandbox creds connected. 3-of-5 approval workflow for filings > ₹50L.

### 2. Vendor Compliance at Scale
5,000 vendors synced from Apex (mock data source). KYC enrichment runs. Compliance scores computed. 30 high-risk flagged. No vendor CRUD — all vendor data read-only from Apex.

### 3. Full GST Cycle (April 2026)
3,000 AP invoices synced from Apex via apex-gateway-service. 800 IRNs + 400 EWBs. GSTR-1 prep → filed → ARN. 2B pulled → recon against Apex purchase register → 95% matched → IMS actions on 120 → 2B regenerated → GSTR-3B auto-fill → offset → filed → ARN. Audit trail complete. CFO dashboard accurate.

### 4. Full TDS Cycle (Q4 FY 2025-26)
500 salaried (24Q), 2,000 non-salary (26Q), 100 non-resident (27Q). TDS calc + 206AB (30 specified-persons caught). Challans deposited + CIN captured. TXT → CSI (OTP) → FVU → E-File → E-Verify → RRR. Q4: Form 16 bulk for 500 employees.

### 5. ITR Filing (AY 2026-27)
500 employees invited. 300 complete in a week. Mix of ITR-1 + ITR-2. 50 complex cases → CA marketplace → filed → fees collected. All ITR-Vs e-verified.

### 6. Bank Open Sibling Sync
All 4 sibling gateways operational. Apex: 5,000 vendors + 3,000 AP invoices synced. Aura: 2,000 AR invoices synced, IRN status published back. Bridge: 200 contracts synced, TDS sections derived. HRMS: 500 payroll records synced, Form 16 data available for 24Q. Data Sources UI shows all connected apps with sync status. GL journals posted back to siblings.

### 7. Annual Filings + Compliance Cloud
GSTR-9 + 9C for one GSTIN → filed. AOC-4 + MGT-7 + DIR-3 KYC for 1 company + 5 directors via MCA21.

### 8. AI / MaxITC
MaxITC identifies ₹2.5L ITC at risk from 12 vendors. LLM generates chase emails (English + Hindi) → sent → replies tracked. ML matching lifts recon rate from 88% (fuzzy) to 95%.

### 9. Reliability / Performance
Kill gstn-gateway mid-filing → restarts → outbox picks up → success. DR failover drill. Load: 1,000 IRN/sec for 15min → P95 < 300ms. 10K concurrent CFO dashboard → P95 < 2s.

### 10. Security + Compliance
OWASP ZAP + Burp → 0 High/Critical. DPDP: data export ZIP; deletion scrubs PII. RLS: cross-tenant SQL returns 0 rows. No secrets in git history.

### 11. Accessibility
axe-core on 30 pages → 0 violations. Keyboard-only: login, GSTR-1 wizard, recon workspace, TDS wizard. Screen reader on filing modal → irreversibility announced.

### 12. Visual Fidelity
Chromatic regression: all pages match baseline in Light Classic, Dark Classic, Ocean. Reference UI screenshots reproduced at pixel fidelity.

### Final acceptance
Test suite green in CI, tagged `v1.0-acceptance`, production go-live approved.

---

## HOW TO RUN THIS

### First session (one-time)
```bash
mkdir -p ~/workspace/complai && cd ~/workspace/complai
claude

# Paste:
"I'm building Complai. The build prompt is at /path/to/complai_build_prompt.md. Read it. Execute Part 0.5 (repo init) now. Then stop and wait for me."
```

### Every subsequent session
```bash
cd ~/workspace/complai
claude

# Paste:
"Read /CLAUDE.md and /BUILD_PLAN.md. Announce current state and next part. Execute it. Update BUILD_PLAN + CLAUDE.md as you go. Commit and summarize when done. If blocked, stop and tell me."
```

### If Claude seems to have "forgotten" something
```
"Re-read /CLAUDE.md, /BUILD_PLAN.md, and any ADRs in /docs/adr/. Also re-read /docs/input/{relevant}.md for the current part. Then continue."
```

### Pacing

- 14 parts × ~9h avg = ~130h of Claude Code execution.
- Calendar: 4-5 weeks at 5-6 active hours/day.
- Parallel engineers (if available): code review + integration testing during Claude Code execution.

### Credential checkpoints

Parts 1, 3, 5, 6, 9, 12, 13, 14 each gate on specific credentials. Claude Code builds with mock gateways if missing — no blocking. Real-integration testing flipped on when creds arrive.

---

## FINAL NOTES

**What this prompt optimizes for:**

- **Resumability** — every part starts by re-reading CLAUDE.md; session boundaries don't lose progress.
- **Verifiable progress** — every part has a hard DoD + passing tests.
- **Invariant preservation** — non-negotiables stated in every part, asserted in tests.
- **Independent value** — Part 5 = working GSTR-1; Part 7 = working ITC engine; Part 9 = working TDS; each is demo-able.

**What you provide as we go:**

- Part 1: AWS account, GitHub org, Cloudflare account.
- Part 3: Temporal Cloud, Last9, Amazon SES production access.
- Part 5: Adaequare sandbox credentials + GST Returns Enriched API doc.
- Part 6: Sandbox.co.in sandbox api_key/secret.
- Part 12: Azure OpenAI subscription + private endpoint.
- Part 13: SAP test system + HDFC/ICICI sandbox + MCA21 test access.
- Pre-prod: production credentials across all providers, real DSC + EVC for final tests.

**Start with Part 0.5, then Part 1.**
