# BUILD_PLAN.md — Living checklist

Last updated: 2026-05-03

## Current part
Part 11 — Sibling gateway services (Aura, Bridge, HRMS) (not started).

## Recently completed
Part 10c — GSTR-9/9C backends. gstr9-service (annual return aggregation, GSTR-9C reconciliation engine, mismatch resolution, certification). gstn-gateway-service extended with 7 GSTR-9/9C endpoints (save/submit/file/status). Mock + real provider stubs. 115+ gstr9-service tests, 80+ new gstn-gateway tests. All handlers ≥80% individually.
Part 10b — AIS reconciliation engine, employee bulk filing flow, ITR-4/5/6/7 form generators + eligibility checkers. Migration 002_bulk_filing.sql, 6 bulk API endpoints, coverage targets met.
Part 10a — itr-service + itr-gateway-service backends. ITA 2025 tax computation engine, ITR-1/2/3 eligibility, TDS reconciliation, 5-head income calculators, mock Sandbox.co.in ITR APIs.
Part 9 — TDS module complete. ITA 2025, 4-digit payment codes, Form 138/140/144 filing wizards, certificates (Form 130/131), challan tracking, 3 Playwright E2E specs, all verifications green.

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

- [x] Part 8: e-Invoicing + E-Way Bill
  - [x] irp-gateway-service: Mock Adaequare IRP APIs — Generate IRN, Cancel IRN, Get by IRN/Doc (port 8098)
    - [x] Request/response mapping to IRP v1.05 schema, double-wrap gateway envelope
    - [x] Idempotency via X-Idempotency-Key header
    - [x] Unit test coverage: API 94.4%, 15 tests
  - [x] einvoice-service: E-Invoice domain logic — generate, cancel (24h window), list, summary (port 8099, einvoice_db)
    - [x] 24h cancellation enforcement: fixedClock-based test, 12h within window → 200, 25h expired → 422, exact 24h boundary → 422
    - [x] CancellationWindowOpen helper with clock interface for deterministic tests
    - [x] ValidityDaysForDistance: 0-200km→1, 201-400→2, 401-600→3, 1500→8 (shared with EWB)
    - [x] IRP gateway integration: generate → store → return signed invoice + QR
    - [x] Unit test coverage: API handlers 65.1%, 18 tests
  - [x] ewb-gateway-service: Mock Adaequare EWB APIs — Generate, Cancel, Update Vehicle, Extend, Consolidate, Get (port 8100)
    - [x] Unit test coverage: API 76.9%, 12 tests
  - [x] ewb-service: E-Way Bill domain logic — full lifecycle (port 8101, ewb_db)
    - [x] State machine: PENDING→ACTIVE→{VEHICLE_UPDATED,EXTENDED,CANCELLED,CONSOLIDATED}
    - [x] CanTransitionTo: 22 transition tests, terminal states (CANCELLED, CONSOLIDATED) block further changes
    - [x] 24h cancellation window: same fixedClock pattern as einvoice, 422 on expired
    - [x] Vehicle update: multi-update history tracked, cancelled EWBs reject updates (422)
    - [x] Extend validity: blocked on cancelled EWBs
    - [x] Consolidate: min 2 EWBs, all must be ACTIVE, atomic status update
    - [x] Validity calculation: regular (200km/day) + ODC (20km/day), 17 test cases
    - [x] Unit test coverage: API handlers 63.9%, 30 tests
  - [x] E-Invoicing UI (apps/web/src/app/compliance/e-invoicing/)
    - [x] List page: 30 mock records, KPIs (Total/Generated/Cancelled), status filter (All/Generated/Cancelled), search
    - [x] Generate page: 3-step flow (select invoice → validate payload → generate IRN), success with QR + signed JSON
    - [x] Bulk generate page: multi-select invoices, batch progress, InvoiceSelector component
    - [x] Detail page: IRN info, signed JSON viewer (collapsible), QR code, cancel button (24h window)
    - [x] Components: IRNStatusPill, EInvoiceKPIs, EInvoiceTable, SelectInvoiceStep, ValidateStep, SignedJsonViewer
    - [x] All files ≤250 lines (split into components/ subdirectories)
  - [x] E-Way Bill UI (apps/web/src/app/compliance/e-way-bill/)
    - [x] List page: 25 mock records, KPIs (Total/Active/Nearing Expiry/Cancelled), 5 status tabs, search
    - [x] Generate page: 5-step flow (select → form → confirm → generating → success)
    - [x] Detail page: 3-col layout (summary + items, vehicle history), cancel modal (type CANCEL)
    - [x] Update Vehicle page: vehicle history timeline, new vehicle form, type-to-confirm
    - [x] Extend Validity page: current/new validity comparison, additional distance with DistanceValidityCalculator
    - [x] Consolidate page: multi-select active EWBs (max 15), type CONSOLIDATE to confirm
    - [x] Components: EwbStatusPill (5 variants), DistanceValidityCalculator, VehicleUpdateTimeline, EwbKPIs, EwbTable, EwbActions, EwbDetailItems
    - [x] All files ≤250 lines
  - [x] Storybook: 3 new story files (EwbStatusPill, DistanceValidityCalculator, VehicleUpdateTimeline), 48 total stories, build clean
  - [x] Playwright E2E: 6 tests across 2 spec files
    - [x] einvoicing-lifecycle.spec.ts: generate IRN → detail with QR + signed JSON → cancel within 24h; cancelled record → no cancel button
    - [x] ewb-lifecycle.spec.ts: generate EWB → 600km = 3 days validity → detail → update vehicle; cancel within 24h; extend validity; cancelled record → no action buttons
  - [x] Postgres RLS: einvoice_db 0 cross-tenant rows, ewb_db 0 cross-tenant rows
  - [x] Goose migrations: einvoice_db (einvoices, einvoice_line_items, outbox), ewb_db (ewb, ewb_items, ewb_vehicle_updates, ewb_consolidations, ewb_consolidation_items, ewb_outbox)
  - [x] TypeScript typecheck clean (tsc --noEmit)
  - [x] go vet + go build clean on all 4 Part 8 services
  - [x] Storybook axe-core: 11 component suites, 0 a11y violations
  - [x] 6 compliance modules now in browser: GST Returns, E-Invoicing, E-Way Bill, ITC Reconciliation, Vendor Compliance, GSTR-3B

### Part 9 final verification (2026-04-30) — COMPLETE

  - [x] Docker: 10/10 containers healthy
  - [x] Goose migrations: tds_db at version 3 (tds_entries, deductees, tds_filings, tds_aggregates + RLS)
  - [x] GOWORK=off: 23/23 services PASS
  - [x] Coverage: tds-service (api 91.4%, domain 95.6%, filing 84.5%, gateway 92.2%, store 24.1%), tds-gateway-service (api 100%, provider 100%)
  - [x] Performance: 1000 TDS calculations in 269µs (benchmark + assertion test)
  - [x] ITA 1961 rejection: TestRejectITA1961Sections — 8 legacy section codes rejected
  - [x] DTAA enforcement: 12/12 non-resident + DTAA tests pass
  - [x] Tax Year terminology: zero "Assessment Year" references in codebase
  - [x] Postgres RLS: 0 cross-tenant rows on all 4 TDS tables (deductees, tds_entries, tds_filings, tds_aggregates)
  - [x] tsc --noEmit: clean
  - [x] go vet + go build: 23/23 clean
  - [x] Storybook: 32 story bundles built clean
  - [x] Playwright E2E (3 tests, 2 spec files):
    - [x] tds-calculation.spec.ts: contractor ₹50K → section 393(1), code 1024, 2% rate, ₹1,000 TDS, save entry
    - [x] tds-form144-filing.spec.ts: DTAA blocks Form 144 submit; Form 138 salary wizard reaches ARN acknowledgement
  - [x] 7 compliance modules now in browser: GST Returns, E-Invoicing, E-Way Bill, ITC Reconciliation, Vendor Compliance, GSTR-3B, TDS

### Part 10a — ITR backends (2026-05-03) — COMPLETE

  - [x] itr-service: ITR domain service (port 8100, itr_db)
    - [x] Domain models: Taxpayer, ITRFiling, IncomeEntry, Deduction, TaxComputation, TDSCredit, AISReconciliation
    - [x] ITA 2025 enforcement: IsOldSectionRef() detects 11 legacy section codes, ITA2025Equivalent() maps to new sections
    - [x] Tax computation engine: ComputeTax() for both regimes, ITA 2025 slab rates, Section 87A rebate, surcharge, 4% HEC
    - [x] New regime slabs: 0-4L=0%, 4-8L=5%, 8-12L=10%, 12-16L=15%, 16-20L=20%, 20-24L=25%, 24L+=30%, std deduction ₹75K
    - [x] Old regime slabs: 0-2.5L=0%, 2.5-5L=5%, 5-10L=20%, 10-15L=30%, std deduction ₹50K, 80C/80D deductions respected
    - [x] 5-head income calculators: Salary (Section 392), House Property, Capital Gains (LTCG/STCG/VDA/112A), Business (incl. Section 44BBC), Other Sources
    - [x] ITR-1/2/3 form eligibility checkers with ITA 2025 rules (ITR-1: ≤50L, ≤2 HP, ≤1.25L LTCG 112A)
    - [x] TDS credit reconciliation: ReconcileTDS() matching AIS/Form 168 entries against claims by TAN+section
    - [x] Form 10-IEA enforcement for old regime opt-out
    - [x] 17 API endpoints: taxpayers CRUD, filings CRUD, compute-tax, income entries, deductions, TDS credits, reconcile-tds, eligibility checks
    - [x] Postgres RLS migration: 7 tables with tenant_id isolation, proper indexes
    - [x] Unit test coverage: domain 96.1%, api 81.6%, store 40.6%
  - [x] itr-gateway-service: Mock Sandbox.co.in ITR APIs (port 8101)
    - [x] 6 endpoints: PAN-Aadhaar link check, AIS/Form 168 fetch, ITR submission, ITR-V generation, e-verification, refund status
    - [x] Mock provider: PAN validation, deterministic AIS data (2 TDS entries), stateful filing lifecycle (submit→ITRV→e-verify)
    - [x] Unit test coverage: api 85.0%, provider 100%
  - [x] go.work updated with 2 new services (25 total)
  - [x] Makefile MIGRATE_ORDER updated with itr-service
  - [x] GOWORK=off: 25/25 services build clean
  - [x] go vet: clean across both services

### Part 10b — AIS reconciliation + bulk filing + ITR-4/5/6/7 (2026-05-03) — COMPLETE

  - [x] 10b-1: AIS reconciliation engine
    - [x] ReconcileAIS() compares Form 130/TDS data against AIS (Form 168)
    - [x] 6 mismatch categories: salary, TDS, interest, dividend, securities, property
    - [x] 3 severity levels: INFO/WARN/ERROR with configurable submission blocking
    - [x] API endpoint: POST /reconcile-ais
    - [x] 20 unit tests covering all categories and severity counting
  - [x] 10b-2: Employee bulk filing flow
    - [x] BulkFilingBatch/BulkFilingEmployee/MagicLinkToken domain types
    - [x] Status tracking: PENDING→PROCESSING→COMPLETED/FAILED (batches), PENDING_REVIEW→APPROVED/REJECTED/SUBMITTED/E_VERIFIED/MISMATCH (employees)
    - [x] ProcessEmployeeForBulkFiling(): per-employee salary computation + tax calculator + AIS reconciliation
    - [x] DetermineFormType(): business→ITR-3, capgains→ITR-2, else ITR-1 eligibility check
    - [x] Magic link tokens: 32-byte crypto/rand, 7-day TTL, 1000-employee batch limit
    - [x] Migration 002_bulk_filing.sql: 3 new tables (bulk_filing_batches, bulk_filing_employees, magic_link_tokens) + RLS, severity column on ais_reconciliations, bulk_batch_id FK on itr_filings
    - [x] 6 API endpoints: create/get/list batches, add/list employees, process batch
    - [x] PgStore + MockStore implementations for all 8 new repository methods
  - [x] 10b-3: ITR-4/5/6/7 form generators
    - [x] ITR-4 (Sugam): presumptive taxation Section 44AD/44ADA/44AE, individuals/HUFs/firms, ≤₹50L, LTCG ≤₹1.25L
    - [x] ITR-5: firms, LLPs, AOPs, BOIs with partner details
    - [x] ITR-6: companies with Schedule MAT, buyback loss, deemed dividend Section 2(22)(f)
    - [x] ITR-7: trusts/charities under Section 139(4A)/(4B)/(4C)/(4D), anonymous donations
    - [x] All 4 eligibility checkers with comprehensive test coverage (11+7+4+8 = 30 test cases)
    - [x] API endpoints: /eligibility/itr4 through /eligibility/itr7
    - [x] CreateFiling now accepts ITR-1 through ITR-7 form types
  - [x] Coverage: domain 96.6%, api 82.5%, store 43.9%
  - [x] GOWORK=off: 25/25 services build clean
  - [x] Migration 002_bulk_filing.sql verified via make migrate-all

### ITA 2025 refactor (2026-04-29) — RESOLVED

**Scoping decision:** Complai supports Income Tax Act 2025 only (effective 1 Apr 2026). ITA 1961 is out of scope.

**Refactor completed (2026-04-29).** All tds-service and tds-gateway-service code now references ITA 2025:

| ITA 1961 | ITA 2025 | Description |
|---|---|---|
| Section 192 | Section 392 | Salary TDS |
| Sections 194C/194I/194J/194Q | Section 393(1) | Resident non-salary TDS |
| Section 195 | Section 393(2) | Non-resident TDS |
| — | Section 393(3) | TCS (not yet implemented) |
| Form 24Q | Form 138 | Salary TDS return |
| Form 26Q | Form 140 | Non-salary TDS return |
| Form 27Q | Form 144 | Non-resident TDS return |
| Assessment Year | Tax Year | Terminology change |
| 3-char section codes | 4-digit payment codes (1001-1092) | Primary discriminator |

**Payment code system:** 4-digit codes (1001-1092) replace old section numbers as the primary dispatch key. Calculator, FVU generator, aggregates, and store all key on PaymentCode. Section is metadata only.

**Key ITA 2025 rules implemented:**
- Section 397(2) no-PAN rate: 20% default, 5% exception for codes 1031/1035
- 4% Health & Education Cess on non-resident TDS (Section 393(2))
- Per-payment rent threshold: ₹50,000/month (not annual aggregate)
- Standard deduction ₹75,000 for salary (Section 392)
- Tax Year = Financial Year (no Assessment Year concept)

**Files refactored:** models.go, calculator.go, filing.go, form138.go (was form24q.go), form140.go (was form26q.go), form144.go (was form27q.go), saga.go, handlers.go, router.go, store.go, mock.go, sandbox.go, migration 003_ita2025.sql, all test files.

**Gateway service:** tds-gateway-service routes, handlers, provider interface, mock provider, and domain models all updated. Routes: /form140/file, /form138/file, /form144/file.

**Test results (GOWORK=off):**
- tds-service: api 91.4%, domain 95.6%, filing 84.5%, gateway 92.2%, store 24.1%
- tds-gateway-service: api 100%, provider 100%

**Remaining:** rules-engine-service TDS applicability rules still reference ITA 1961 section numbers (deferred to Part 14 hardening).

### Deferred hardening (→ Part 14)

#### Go Module Verification Under GOWORK=off (recurring issue, Parts 5 and 7)
- [ ] Add `make verify-go-modules` target that runs `GOWORK=off go test ./...` across `packages/shared-kernel-go` and every directory in `services/go/`. This catches go.sum drift that workspace mode (GOWORK=on, default) silently masks. Without this, tests pass on dev machines but fail in Docker containers (which run as standalone modules).
- [ ] Make this part of `make verify` (broader verification target).
- [ ] Add as pre-commit hook in `.githooks/pre-commit`.
- [ ] CI: run on every PR.
- Discovered: Parts 5 and 7 — verification reported "all green" while multiple services had failing tests under GOWORK=off after the workspace cache was bypassed.

#### Gateway Response Envelope Cleanup (discovered in Part 6)
- [ ] Standardize gateway response envelope across all gateway services (apex, aura, bridge, hrms, gstn, kyc).
- [ ] Eliminate double-wrap from httputil.JSON() over GatewayResponse{Data, Meta}. Currently produces `{"data": {"data": ..., "meta": ...}}` which every consumer must handle. Either remove outer wrap, smart-wrap (detect data field), or standardize the doubled shape and document it as the contract.

#### Part 8 coverage uplift (discovered in Part 8e verification, corrected in Part 8.5)
Per-package breakdown (only `api` packages have tests; gateway, store, domain = 0%):
- einvoice-service: `internal/api` 65.1%, `internal/gateway` 0%, `internal/store` 0%, `internal/domain` 0%
- ewb-service: `internal/api` 63.9%, `internal/gateway` 0%, `internal/store` 0%, `internal/domain` 0%
- ewb-gateway-service: `internal/api` 76.9%, `internal/gateway` 0%, `internal/store` 0%
- [ ] Add unit tests for gateway packages (HTTP client mocking, retry logic, error mapping)
- [ ] Add unit tests for store packages (sqlc-generated query verification with testcontainers)
- [ ] Add unit tests for domain packages (business rule validation, state machine transitions)
- [ ] Bring aggregate per-service coverage to ≥80%
- Critical compliance paths (24h cancel, distance validity, state machine) are tested and green in handler tests — the gap is breadth across non-api packages.

#### Coverage reporting hygiene (added Part 8.5)
- When verifying Go test coverage, always run `go test -cover ./...` and report the per-package breakdown.
- Never present a single package's coverage (e.g. `internal/api`) as representative of the full service.
- Aggregate service coverage = weighted average across all packages with source files.
- If only one package has tests, state that explicitly — do not imply the number represents the service.

#### Database hygiene (completed Part 8.5)
- [x] Removed orphan databases from postgres-init.sh: `ap_db`, `billing_db`, `vendor_db` (remnants of pre-Part 4.5 scope, never had migrations or services)
- [x] Dropped orphan databases from running Postgres container
- [x] Reorganized DATABASES array by build part sequence with comments
- Forward-provisioned databases kept: `tds_db` (Part 9), `gstr9_db`/`itr_db` (Part 10), `reporting_db` (Part 14), `secretarial_db` (future)

#### Part 5 hardening tests
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

### ITA 2025 coverage gaps (Part 14)
- [ ] rules-engine-service: update TDS applicability rules from ITA 1961 section numbers to ITA 2025 (392/393 + payment codes)
- [ ] Payment codes not yet implemented: VDA/crypto (1066), online gaming (1067), e-commerce operator (1064), partner remuneration (1059), lottery/crossword (1052-1056)
- [ ] Section 393(3) TCS: full implementation (codes 1070-1092), separate Form 141/143
- [ ] tds-service store package coverage uplift from 24.1% to ≥80% (testcontainers-go)
- [ ] DTAA rate override validation: verify treaty rates against ITA 2025 Schedule IV mappings

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

### Part 14 Hardening — gstn-gateway-service legacy handler coverage

Discovered May 3, 2026 during 10c-4 verification. The "all handlers ≥80%" coverage standard was introduced in 10c-3.5 and applied to gstr9-service. When we ran per-handler coverage on gstn-gateway-service for the first time during 10c-4, we found 12 pre-existing handlers below 80%, including 9 at 0% coverage:

- GSTR2BGet, GSTR2AGet — 0% (Part 7)
- IMSGet, IMSAction, IMSBulkAction — 0% (Part 7)
- GSTR3BSave, GSTR3BSubmit, GSTR3BFile — 0% (Part 7)
- GSTR1SummaryHandler — 0% (Part 5/7)
- Authenticate — 57.1% (Part 5)
- GSTR1Get, GSTR1Status — 76.5% each (Part 5)

These pre-date the per-handler standard. Aggregate coverage was passing in their respective parts because new handlers in those parts had high coverage and dragged the average up. The 7 new GSTR9/9C handlers added in 10c-4 are at 94.1% each — clean.

Action in Part 14: add per-handler tests for the 12 handlers above. Target ≥80% on each.

Note: these handlers have been exercised by Playwright E2E tests in Parts 5, 7, and 8 final-verification commits, which provides some real-world confidence even without unit-level coverage. Risk profile: not "production broken" but "future regressions silent."
### Storybook story policy clarification (May 3, 2026, 10d-1)

After 10d-1, we clarified the rule for which components must have Storybook stories:

- **Reusable domain components** (in `apps/web/src/app/compliance/<module>/components/`): MUST have stories — these contribute to the design system and need a11y/visual regression coverage.
- **Route-internal step components** (in `apps/web/src/app/compliance/<module>/<route>/components/`, e.g., wizard step components like SelectStep / ConfigureStep / ConfirmStep): stories are OPTIONAL — these are page-internal scaffolding tightly coupled to route state, not design system contributions.

Rationale: stories enforce a design system standard. Route-internal scaffolding doesn't benefit from that standard and adds friction without adding value.

Going forward, Claude Code should default to creating stories for components in the `components/` subdirectory adjacent to module root, but may skip stories for components in route-specific subdirectories.
