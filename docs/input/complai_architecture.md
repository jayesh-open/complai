# Complai — Technical Architecture

**Version:** 1.0
**Status:** Approved for build
**Companion docs:** PRD, Design System, API Integration Spec, Build Prompt

---

## 1. Overview

Complai is an enterprise multi-tenant compliance SaaS platform for Indian businesses. This document specifies the technical architecture: services, data stores, messaging, infrastructure, scaling model, security posture, and operational practices.

The architecture is designed to:

- Support **500 enterprise customers in Year 1, scaling to 50,000 GSTINs by Year 3**
- Meet **99.99% SLO on Tier-0 services** during filing peaks (10th/11th/15th/20th of every month)
- Enforce **strict multi-tenant isolation** at the database level
- Remain **cost-efficient** in Phase 1 without compromising production-grade reliability
- **Deploy entirely in India** (`ap-south-1` / `ap-south-2`) for DPDP Act data-residency compliance

---

## 2. Architectural Principles

1. **Multi-tenancy is invariant.** Every row in every table carries `tenant_id`. Postgres Row-Level Security (RLS) enforces isolation at the database layer, independent of application code correctness. No exceptions.

2. **One primary language per layer.** Go for all backend services (domain + gateway). Python only for AI/ML. TypeScript for frontend and BFFs. No polyglot sprawl within a layer.

3. **Outbox pattern for every external call.** Business services never call external APIs directly. State changes write to an outbox table in the same transaction, and a CDC pipeline ships them to SQS for reliable delivery. Every external call is idempotent via `request_id`.

4. **Managed where it matters.** Temporal Cloud for orchestration, Amazon RDS Multi-AZ for Postgres, Last9 for observability. We run what differentiates us; we buy what doesn't.

5. **One Postgres, multiple databases, RLS everywhere.** We don't run MongoDB, Kafka, or ClickHouse in Phase 1. A single RDS Postgres instance with per-service logical databases plus a read replica for analytics handles our entire transactional and analytical surface at this scale.

6. **Observability first.** Every service emits OTel traces, metrics, and structured logs. `tenant_id`, `gstin`, `tan`, `pan` ride as baggage on every span. Nothing is silently dropped.

7. **Provider abstraction.** Internal consumers never see Adaequare or Sandbox.co.in specifics. Gateway services normalize responses into a common internal contract.

8. **Design for the 11th.** The 10th/11th/15th/20th of every month are traffic spikes of 10–20× baseline. The architecture assumes peak is the common case, not the exception.

---

## 3. Service Topology

Complai is composed of roughly **28 services** organized into 5 families.

### 3.1 Platform services (9)

Foundational services every domain module depends on.

| Service | Language | Purpose |
|---|---|---|
| `identity-service` | Go | Authentication, MFA, SSO, session management |
| `tenant-service` | Go | Tenant lifecycle, PAN/GSTIN/TAN hierarchy |
| `user-role-service` | Go | Roles, permissions, approval workflows |
| `master-data-service` | Go | Vendors, customers, employees, items, HSN, pincodes |
| `document-service` | Go | S3 + Postgres metadata for all documents; OCR trigger |
| `notification-service` | Go | Email (SES), WhatsApp, SMS, in-app |
| `audit-service` | Go | Tamper-evident audit log; search and export |
| `workflow-service` | Go | Temporal Cloud integration layer |
| `rules-engine-service` | Go | Tax determination, validation, eligibility rules |

### 3.2 Domain services (11)

One per major compliance workflow.

| Service | Language | Purpose |
|---|---|---|
| `gst-service` | Go | GSTR-1, GSTR-3B, GSTR-2A/2B/IMS lifecycle |
| `gstr9-service` | Go | Annual returns (GSTR-9, GSTR-9C) |
| `einvoice-service` | Go | IRN generation, signed JSON/QR, canonical transformer |
| `ewb-service` | Go | E-Way Bill lifecycle including multi-vehicle |
| `tds-service` | Go | TDS/TCS deductees, payments, returns, Form 16 |
| `itr-service` | Go | ITR filing for employer-bulk and CA-assisted flows |
| `vendor-service` | Go | Vendor master, onboarding, compliance scoring, portal |
| `recon-service` | Go | 5-stage match pipeline, IMS actions, bucket views |
| `ap-service` | Go | Invoice ingestion, OCR, 3-way match, approvals, payments |
| `billing-service` | Go | Complai One SMB billing, recurring invoices, payment links |
| `secretarial-service` | Go | ROC filings (MCA21), registers, minutes, resolutions |

### 3.3 Gateway services (8)

Thin, single-responsibility services that talk to external providers.

| Service | Provider | Purpose |
|---|---|---|
| `gstn-gateway` | Adaequare Enriched | GSTR returns, IMS, ledgers, challans |
| `irp-gateway` | Adaequare Enriched | e-Invoice IRN |
| `ewb-gateway` | Adaequare Enriched | E-Way Bill |
| `tds-gateway` | Sandbox.co.in | TDS calc, FVU, e-File, Form 16 |
| `itd-gateway` | Sandbox.co.in | ITR filing, prefill, AIS, 26AS |
| `kyc-gateway` | Sandbox.co.in | PAN, Aadhaar, GSTIN, bank, MCA, Udyam |
| `tax-payment-gateway` | Sandbox.co.in | TDS and income tax challans |
| `bank-gateway` | Direct (HDFC, ICICI, Axis, SBI, Razorpay, Cashfree) | Payments and collections |
| `mca-gateway` | Direct (MCA21 V3) | ROC filings |
| `erp-gateway` | Direct (SAP, Oracle, Tally, Dynamics) | Bidirectional ERP sync |

### 3.4 AI/ML services (3) — Phase 4

| Service | Language | Purpose |
|---|---|---|
| `matching-ml-service` | Python | CatBoost model for reconciliation Stage 3 |
| `llm-copilot-service` | Python | Vendor comms, NL queries, explanation generation |
| `ocr-service` | Python | Form 16 / 26AS / invoice extraction |

### 3.5 Web-facing services (4)

| Service | Language | Purpose |
|---|---|---|
| `web-bff-service` | Node + NestJS | BFF for `apps/web` (main product) |
| `portal-bff-service` | Node + NestJS | BFF for `apps/vendor-portal` |
| `smb-bff-service` | Node + NestJS | BFF for `apps/complai-one` |
| `reporting-service` | Go | Report execution (PDF/Excel/CSV generation) |

---

## 4. Data Architecture

### 4.1 Primary database — Amazon RDS for PostgreSQL 16

**Phase 1 setup:**

- **Instance class:** `db.r7g.2xlarge` (8 vCPU, 64 GB RAM, up to 2.5 Gbps network)
- **Deployment:** Multi-AZ for HA; synchronous replication to standby in a different AZ
- **Storage:** 500 GB gp3 with 12,000 IOPS provisioned
- **Backups:** automated daily + 35-day PITR retention
- **Read replica:** same instance class in same region, asynchronous replication — dedicated to analytics/BI workloads

**Database organization:**

One RDS instance hosts ~20 logical databases, one per owning service. Each service owns its schema and migrations; no service writes to another's database.

```
complai-primary-rds
├── identity_db
├── tenant_db
├── user_role_db
├── master_data_db
├── document_db
├── audit_db
├── gst_db
├── gstr9_db
├── einvoice_db
├── ewb_db
├── tds_db
├── itr_db
├── vendor_db
├── recon_db
├── ap_db
├── billing_db
├── secretarial_db
├── workflow_db
├── reporting_db
└── rules_engine_db
```

**Multi-tenancy pattern:**

Every table has:
- `tenant_id UUID NOT NULL`
- RLS policy: `USING (tenant_id = current_setting('app.tenant_id')::uuid)`
- Tenant-scoped indexes: `(tenant_id, ...)` rather than single-column indexes

A shared Go middleware (`packages/shared-kernel-go/tenant`) extracts `tenant_id` from the JWT, sets the Postgres session variable at transaction start, and ensures it propagates. A service method that forgets to set it gets zero rows back — RLS is a compile-time-ish safety net, not just a runtime check.

**JSONB for flexible fields:**

Where schema variability is needed (OCR extraction output, rules engine config, document metadata, user preferences), we use Postgres `JSONB` columns with GIN indexes on queried paths. This gives us Mongo-like flexibility without a second database.

### 4.2 Phase 1 analytics — Postgres read replica

All BI, reporting, and dashboard workloads read from the read replica. The primary handles only OLTP.

**Materialized views** pre-compute expensive aggregates, refreshed on schedule:

- `mv_vendor_compliance_scores` — hourly refresh
- `mv_filing_summary_by_period` — daily refresh
- `mv_itc_at_risk` — every 15 min
- `mv_cfo_dashboard_kpis` — every 15 min
- `mv_recon_bucket_counts` — on-demand + every 5 min during business hours

**Upgrade path:** when a single table crosses ~500M rows (likely `audit_events` or `transaction_stream` in Year 2), we add ClickHouse alongside via Debezium CDC for that specific workload. Until then, one Postgres does both jobs.

### 4.3 Cache — Amazon ElastiCache for Redis 7

**Setup:** cluster-mode enabled, 3 shards × 2 replicas, `cache.r7g.large` nodes.

**What we cache:**

| Purpose | Key pattern | TTL |
|---|---|---|
| Adaequare bearer token | `provider:adaequare:{env}:bearer` | 23h |
| Sandbox bearer token | `provider:sandbox:{env}:bearer` | 23h |
| Tenant-scoped provider credentials | `tenant:{id}:adaequare:{gstin}:creds` | 5 min (hot cache) |
| GSTIN KYC cache | `kyc:gstin:{gstin}` | 7 days |
| PAN KYC cache | `kyc:pan:{pan}` | 30 days |
| TAN KYC cache | `kyc:tan:{tan}` | 30 days |
| HSN validation cache | `kyc:hsn:{code}` | 90 days |
| User session | `session:{user_id}:{device_id}` | 24h |
| CSRF tokens | `csrf:{session_id}` | 1h |
| Rate-limit counters | `rl:{tenant}:{user}:{action}` | sliding window |
| OTP reference IDs (CSI flow) | `otp:{flow_id}` | 2 min |
| Filing idempotency keys | `filing:idem:{request_id}` | 24h |
| Workflow semaphores | `wf:lock:{gstin}:{period}:{form}` | 10 min |

Redis pub/sub powers real-time UI updates (bulk operation progress, filing status changes).

### 4.4 Binary storage — Amazon S3

**Buckets:**

- `complai-{env}-documents` — all persistent documents (invoices, signed JSONs, FVU files, Form 16 PDFs, minutes)
- `complai-{env}-uploads` — temporary staging for uploads before OCR
- `complai-{env}-exports` — generated reports, zip bundles
- `complai-{env}-backups` — application-level backups (audit log archives, etc.)

**Configuration:**

- Versioning enabled
- Server-side encryption with KMS (tenant-scoped DEKs)
- Lifecycle rules: Standard → IA after 90 days → Glacier after 2 years → Deep Archive after 5 years
- VPC endpoints so service-to-S3 traffic stays on the AWS backbone
- Bucket policies deny public access; all access via IAM roles or presigned URLs
- Object lock enabled on the documents bucket for regulatory retention (8-year for GST, 7-year for TDS)

**Pre-signed URLs** (15-minute TTL) for all direct client access. No service proxies large files.

### 4.5 Messaging — Amazon SQS + SNS

**Pattern: Single-consumer queues for most flows; SNS+SQS fan-out where multiple services consume the same event.**

**Standard queues (single-consumer):**

- `gov.outbound.gstn.fifo` — FIFO queue grouped by GSTIN; consumed by `gstn-gateway`
- `gov.outbound.irp.queue` — consumed by `irp-gateway`
- `gov.outbound.ewb.queue` — consumed by `ewb-gateway`
- `gov.outbound.tds.queue` — consumed by `tds-gateway`
- `gov.outbound.itd.queue` — consumed by `itd-gateway`
- `gov.outbound.kyc.queue` — consumed by `kyc-gateway`
- `gov.outcome.*.queue` — per-domain outcome queues consumed by the originating business service
- `ocr.jobs.queue` — consumed by `ocr-service`
- `notification.jobs.queue` — consumed by `notification-service`

**SNS topics with fan-out:**

- `FilingCompleted.topic` → SQS subscribers: `notification-service`, `gl-stream`, `audit-service`, `maxitc-orchestrator`
- `InvoiceCreated.topic` → SQS subscribers: `audit-service`, `recon-service`, `matching-ml-service`
- `VendorCreated.topic` → SQS subscribers: `audit-service`, `notification-service`, `kyc-gateway`
- `MasterDataChanged.topic` → SQS subscribers: dependent services

**Operational details:**

- Every queue has a **dead-letter queue** with a 5-retry threshold
- CloudWatch alarms on DLQ depth > 0
- **FIFO queues** (message group ID = resource key) where strict ordering matters: filing sagas per GSTIN, TDS flows per TAN
- **S3 pointer pattern** for messages > 256 KB (payload in S3, message contains bucket+key)
- **Message Bus interface** in `packages/shared-kernel-go/messaging` abstracts SQS/SNS; publishers call `bus.Publish(topic, event)` and consumers register `bus.Subscribe(queue, handler)`. No SQS primitives leak into business code.

### 4.6 Search — Amazon OpenSearch Service

**Setup:** 3-node cluster, `m7g.large.search` instances, 200 GB gp3 each.

**Indexes:**

- `audit-events-{YYYY-MM}` — 30-day hot retention; ISM policy archives older months to S3 with Iceberg metadata
- `documents` — full-text search across invoice content, vendor names, resolution text
- `compliance-search` — cross-module search ("find all filings for Acme GSTIN in 2025")

### 4.7 Long-term archive — S3 + Apache Iceberg

For the 8-year regulatory retention requirement:

- Audit events older than 30 days move from OpenSearch to S3 (Glacier tier)
- Iceberg metadata tables on S3 enable time-travel queries
- Accessible via Amazon Athena for ad-hoc regulatory queries
- Cold — not queried often, but required for compliance

---

## 5. Workflow Orchestration — Temporal Cloud

**Managed service:** Temporal Cloud hosted in AWS Mumbai (low-latency to our EKS).

**Workflow namespaces:**

- `complai-filings` — GSTR-1/3B/9 sagas, TDS quarter flows, ITR filing
- `complai-reconciliation` — 2A/2B pulls, recon runs, IMS actions
- `complai-bulk` — bulk IRN generation, bulk EWB, Form 16 generation, vendor bulk import
- `complai-ap` — invoice ingestion, 3-way match, approval chains, payment execution
- `complai-onboarding` — tenant and vendor onboarding flows

**Workflow characteristics:**

- **Determinism enforced** — workflow code only calls activities, never direct I/O
- **Activities are Go functions** exposing typed inputs/outputs
- **Human-task integration:** a workflow emits a `human_task` event; the UI picks it up, shows the task to the user; on user action, workflow resumes via signal
- **Compensations:** every workflow has saga semantics — if Step N fails, Steps 1..N-1 are compensated
- **Visibility:** Temporal's built-in workflow search indexed by tenant, GSTIN, period, form type

**Why Temporal Cloud instead of self-hosted:**

- Zero operational burden (Temporal has non-trivial ops tax — history shards, visibility backend, cluster upgrades)
- First-class Go SDK
- AWS-native hosting, low latency to our services
- Graduates to self-hosted later if scale demands, with zero code changes

---

## 6. External Integration Strategy

Complai integrates with external systems through well-defined gateway services. Each gateway speaks one external provider's protocol and exposes a uniform internal contract to the rest of the platform.

### 6.1 Two primary API providers

- **Adaequare uGSP (Enriched APIs)** — GST, e-Invoice, E-Way Bill
- **Sandbox.co.in (Quicko)** — TDS, ITR, KYC, Tax Payment

The API Integration Specification document details the full endpoint list, auth flow, and headers for each.

### 6.2 Direct integrations (no aggregator)

- **MCA21 V3** — ROC filings (no aggregator provides end-to-end coverage)
- **Banks** — HDFC ENet, ICICI iConnect, Axis CIB, SBI ePay, Razorpay, Cashfree, TReDS platforms (RXIL, Invoicemart, M1)
- **ERPs** — SAP (RFC/BAPI + OData), Oracle (SuiteTalk), Tally (ODBC + XML), Dynamics (Dataverse API)

### 6.3 Gateway service common shape

Every gateway service implements:

```
POST /v1/gateway/{provider}/{action}
Headers:
  X-Tenant-Id:        <uuid>
  X-Idempotency-Key:  <uuid>
  X-Trace-Id:         <otel-trace>
Body:                 <plain JSON, gateway-specific schema>

Response:
  200/202: { data, meta: { request_id, latency_ms, provider_status } }
  4xx/5xx: { error: { code, message, details, retry_after_ms } }
```

Internal callers never see provider-specific fields, URLs, or error codes.

---

## 7. Infrastructure — AWS

### 7.1 Regions

- **Primary:** `ap-south-1` (Mumbai)
- **DR:** `ap-south-2` (Hyderabad) — warm standby
- All data stays within India (DPDP Act compliance)

### 7.2 Compute — Amazon EKS

**Cluster:** `complai-prod-{region}`, Kubernetes 1.30

**Node groups:**

| Name | Instance type | Min–Max | Purpose |
|---|---|---|---|
| `system` | `m7i.large` | 2–4 | Istio, ArgoCD, cert-manager, monitoring agents |
| `application` | `m7i.xlarge` | 5–30 | All Complai services (autoscaled) |
| `batch` | `c7i.2xlarge` | 0–10 | Recon batch jobs, OCR, bulk ops (scales to zero) |
| `ai` | `g5.xlarge` | 0–3 | Phase 4+ — ML inference (scales to zero) |

**Autoscaling:**

- Cluster Autoscaler for node-level
- HPA (Horizontal Pod Autoscaler) for pod-level — scale on CPU + custom metrics (SQS queue depth, request rate)
- KEDA for event-driven autoscaling (scale based on SQS backlog)

**Istio Ambient Mesh 1.22:**

- Service-to-service mTLS (all internal traffic encrypted)
- L7 routing, retries, circuit breakers
- Traffic policies (per-tenant rate limits, canary deployments)

### 7.3 Networking

**VPC:** custom VPC per environment, `/16` CIDR

**Subnets (3 AZs):**

- **Public subnets** (3) — ALB, NAT Gateway
- **Private subnets** (3) — EKS nodes, Temporal workers
- **Data subnets** (3) — RDS, ElastiCache, OpenSearch (no internet egress)

**VPC endpoints** for AWS services (no NAT charges, better latency):

- S3 Gateway endpoint
- Interface endpoints: SQS, SNS, Secrets Manager, STS, ECR, KMS, SES, Systems Manager

**Egress:**

- NAT Gateway with **Elastic IPs** — Adaequare and Sandbox.co.in require IP whitelisting; we register these EIPs at onboarding
- Cloudflare proxies inbound traffic; ALB accepts traffic only from Cloudflare IP ranges

### 7.4 CDN + DNS + WAF — Cloudflare

- **DNS:** Cloudflare (all Complai domains — `complai.in`, `app.complai.in`, `api.complai.in`, etc.)
- **CDN:** Cloudflare — caches static assets (JS, CSS, images) at edge
- **WAF:** Cloudflare WAF (OWASP Top 10 rules, custom rules, bot management)
- **DDoS protection:** included
- **TLS termination:** at Cloudflare edge with origin-pull TLS to ALB
- **Origin lockdown:** ALB security group allows only Cloudflare IP ranges (prevents origin-direct bypass)

### 7.5 Load balancing

**AWS Application Load Balancer (ALB)** in public subnet:

- Terminates TLS from Cloudflare origin-pull
- Routes by hostname to appropriate Ingress in EKS
- Ingress controller: AWS Load Balancer Controller (integrates ALB with K8s)

### 7.6 Email — Amazon SES

- **Configured domains:** `complai.in`, `notifications.complai.in`, tenant-specific sender subdomains
- **DKIM + SPF + DMARC** set up for all sending domains
- **Configuration sets** for tracking (opens, clicks, bounces, complaints)
- **Event destinations:** SNS topic → notification-service for bounce handling

**Per-tenant inbound email:**

- Customers get addresses like `ap-acme@inbox.complai.in` for AP invoice ingestion
- SES receives → writes email to S3 → Lambda triggers → enqueues to `ap.ingestion.queue`

### 7.7 Secrets — AWS Secrets Manager + KMS

**Secrets Manager** holds:

- Provider credentials: `complai/adaequare/{env}/{field}`, `complai/sandbox/{env}/{field}`
- Tenant-scoped provider credentials: `complai/tenant/{tenant_id}/adaequare/{gstin}/{field}`
- Database master passwords (rotated automatically via Lambda)
- JWT signing keys

**KMS** holds:

- One CMK per tenant (for tenant-scoped DEK envelope encryption)
- Platform CMKs for application-level encryption
- All at-rest encryption (RDS, S3, ElastiCache, OpenSearch) uses KMS

**IRSA (IAM Roles for Service Accounts):** every pod gets a narrow IAM role scoped to the specific secrets and resources it needs. No shared credentials, no baked-in keys.

### 7.8 Observability — Last9

**Last9** is the observability platform. All three signals (metrics, logs, traces) ship to Last9 via the OpenTelemetry Collector.

**OTel Collector deployment:**

- DaemonSet on every node
- Receives OTLP from every pod
- Enriches with K8s metadata
- Exports to Last9 endpoint

**What we track:**

- **Metrics:** per-service RED (rate/errors/duration), SLI dashboards per tier, business metrics (filings/min, IRNs/min, recon-match-rate)
- **Logs:** structured JSON, tenant_id and correlation IDs on every line, PII redacted
- **Traces:** full request traces across services; tenant_id and resource IDs as baggage

**Dashboards:**

- Per-service: RED metrics, error rate, dependency health
- Per-provider: Adaequare success rate/latency, Sandbox success rate/latency
- Per-tenant: usage, errors, filing success rate
- SLO: Tier-0/1/2 burn rate tracking

**Alerting:**

- Last9 alert rules → PagerDuty for on-call rotation
- Slack notifications for non-critical alerts
- Alert on: filing success rate <99%, gateway error rate >1%, auth failure spike, RLS policy violation (any), DLQ depth >0

### 7.9 Identity — Keycloak on EKS

**Self-hosted Keycloak 24** on EKS:

- 3-node cluster behind an internal ALB
- Postgres backend on RDS (separate database, `keycloak_db`)
- Per-tenant realms for isolated SSO configuration
- SAML 2.0 and OIDC for enterprise SSO (Google Workspace, Azure AD, Okta)
- MFA: TOTP, SMS via SES, email

**Why self-hosted vs Cognito:**

Keycloak's per-realm model fits multi-tenant enterprise SSO cleanly. Cognito's pools are awkward at tenant scale (>100 tenants). The operational cost of Keycloak is modest (a handful of engineer-hours/month) and pays back in flexibility.

### 7.10 CI/CD — GitHub Actions + ArgoCD

**GitHub Actions** for:

- Per-service build matrix (Go, Python, Node)
- Test gates (unit + integration + security scans)
- Docker image build + push to ECR
- Trivy vulnerability scan
- SBOM (Software Bill of Materials) generation
- Helm chart render + diff

**ArgoCD** for:

- GitOps deployment (every environment's desired state is a Git repository)
- Automatic sync of staging; manual approval for production
- Progressive rollouts (canary, blue/green)
- Rollback via Git revert

---

## 8. Multi-Tenancy Model

### 8.1 Tenancy tiers

Complai supports four tenancy tiers, selected per customer at onboarding:

| Tier | Database | EKS | Network | Price point |
|---|---|---|---|---|
| **Pooled** | Shared RDS, RLS | Shared cluster, namespaced | Shared VPC | Standard SMB/mid-market |
| **Bridge** | Shared RDS, per-tenant schema | Shared cluster, dedicated namespace | Shared VPC with per-tenant NetworkPolicy | Standard enterprise |
| **Silo** | Dedicated RDS | Dedicated EKS cluster | Dedicated VPC | Large enterprise / BFSI |
| **On-Premise** | Customer infrastructure | Customer infrastructure | Customer network | Strategic accounts (Year 2+) |

Phase 1 focuses on Pooled + Bridge tiers.

### 8.2 Tenant isolation mechanisms

1. **Database:** Postgres RLS on every table
2. **API:** JWT claims carry `tenant_id`; Istio policy rejects requests without it
3. **Cache:** Redis keys namespaced by `tenant_id`
4. **S3:** Objects prefixed with `tenant_id`; IAM policy restricts cross-prefix access
5. **KMS:** Per-tenant CMKs for DEK envelope encryption
6. **Temporal:** Workflow IDs include `tenant_id`; Temporal search is tenant-scoped
7. **Observability:** `tenant_id` as a label on metrics, logs, traces (enables per-tenant views without cross-tenant leaks)

### 8.3 Tenant hierarchy

```
Tenant (e.g., "Acme Manufacturing Pvt Ltd")
├── PAN 1 ("Acme Mfg")
│   ├── GSTIN 1 (Karnataka)
│   ├── GSTIN 2 (Maharashtra)
│   └── TAN 1
├── PAN 2 ("Acme Services")
│   ├── GSTIN 3 (Tamil Nadu)
│   └── TAN 2
└── CIN 1 (company master)
    └── DINs (directors)
```

A user's JWT carries `active_pan` and `active_gstin` as hints; the user can switch via the header dropdowns in the UI.

---

## 9. Outbox + Event Delivery Pattern

Every state change that needs to trigger an external call or cross-service event uses the outbox pattern.

### 9.1 Flow

1. Business service starts a Postgres transaction
2. Service writes domain changes (e.g., update `gstr1_filings.status = 'submitting'`)
3. Service writes an outbox row: `outbox (id, aggregate_type, aggregate_id, event_type, payload, target_queue)`
4. Transaction commits — both writes atomic
5. **Outbox publisher** (sidecar in each service) polls outbox every 500ms, publishes to SQS/SNS, marks outbox row as `published`
6. Target gateway service consumes from SQS, makes external call, publishes outcome event
7. Originating service consumes outcome event, updates state, emits downstream events

### 9.2 Why not CDC (Debezium)?

At our Phase 1 scale, a polling-based outbox publisher is simpler to operate than a Debezium cluster with Kafka Connect. Polling every 500ms is fine for our throughput (low latency compared to external API latency anyway). If we outgrow this, we switch to Debezium → Kinesis in Phase 2+ without changing the outbox schema.

### 9.3 Idempotency

- Every outbox row has `request_id` (UUID)
- External providers dedupe on this
- If the gateway crashes after making the external call but before marking the outbox row, retry sends the same `request_id` — provider returns the same response, no double-submission

---

## 10. Canonical Invoice Schema

A single schema for every invoice across the platform, regardless of source (ERP, email, OCR, user-entered, e-Invoice inbound). Each integration transforms into this schema before it enters the platform's core.

### 10.1 Schema summary

```
Invoice {
  tenant_id, pan, gstin
  id (ULID), document_number, document_date
  supply_type, document_type, reverse_charge
  supplier { gstin, name, address, state_code }
  buyer { gstin, name, address, state_code }
  line_items[] {
    item_id, description, hsn, unit, quantity,
    unit_price, discount, taxable_value,
    cgst { rate, amount }, sgst, igst, cess
  }
  totals { taxable, cgst, sgst, igst, cess, round_off, grand_total }
  payment { mode, bank_details }
  references { po_number, grn_number, contract_number, irn, ewb_number }
  metadata { source_system, source_document_id, created_by, created_at, tags[] }
}
```

Defined in Protobuf (generates Go structs) and JSON Schema (used by the rules engine and frontend validation). Lives in `packages/events/schemas/`.

---

## 11. Security Architecture

### 11.1 Network security

- VPC with private subnets for everything sensitive
- Security groups restrictive by default (deny all, allow explicit)
- Cloudflare WAF in front of ALB
- Origin lockdown (ALB accepts Cloudflare IPs only)
- VPC Flow Logs captured to S3 for audit
- AWS GuardDuty + Security Hub for continuous threat detection

### 11.2 Data protection

- **Encryption at rest:** RDS, S3, ElastiCache, OpenSearch — all via KMS
- **Encryption in transit:** TLS 1.3 external; mTLS via Istio internal
- **Tenant-scoped DEKs:** sensitive fields (bank accounts, PAN full values, credentials) encrypted with per-tenant DEK wrapped by tenant's KMS CMK
- **PII handling:** masking in logs (last 4 of PAN/Aadhaar/bank only); redaction before LLM calls

### 11.3 Access control

- **IAM + IRSA** for AWS resource access (no baked-in keys)
- **RBAC in Kubernetes** for intra-cluster access
- **JWT + fine-grained permissions** at application layer
- **Maker-checker** for high-impact actions (filing, high-value payments, vendor approvals)
- **Step-up authentication** for filing operations (re-auth within 5 min)

### 11.4 Audit trail

- Every state-change event across every service written to `audit_events` in Postgres + OpenSearch
- Hourly Merkle-chain hashing: each hour's log hashes concatenated with previous hour's hash, published to a write-once S3 bucket
- Tamper detection: if an audit event is altered, the chain breaks at the next hourly hash
- Exportable to signed PDF for regulatory submission

### 11.5 Vulnerability management

- CI scans: Trivy (containers), Snyk (dependencies), GitLeaks (secrets)
- Runtime: GuardDuty
- Weekly patching cadence for non-critical; immediate for critical
- Quarterly pen tests (Phase 2+)

### 11.6 Compliance posture

- **DPDP Act (Day 1):** consent flows, data-export endpoint, data-delete endpoint, PII access audit
- **SOC 2 Type II:** evidence collection automated via Drata; audit window starts Month 6, report Month 12
- **ISO 27001:** pursued alongside SOC 2; certification target Month 18

---

## 12. Reliability & DR

### 12.1 Availability targets

- Tier-0 services: 99.99% (~52 min/year downtime budget)
- Tier-1 services: 99.9% (~8.7 hr/year)
- Tier-2 services: 99.5% (~1.8 days/year)

### 12.2 Redundancy

- EKS nodes across 3 AZs
- RDS Multi-AZ (automatic failover in <60s)
- ElastiCache Multi-AZ
- OpenSearch multi-AZ
- Temporal Cloud handles its own redundancy
- Stateless services: 3+ replicas each
- Stateful state held only in RDS, ElastiCache, S3, OpenSearch

### 12.3 Backup

- **RDS:** automated daily backups + 35-day PITR + weekly snapshots retained 1 year
- **S3:** versioning + cross-region replication to DR region
- **ElastiCache:** nightly snapshots (for data that's cache-primary, like session data — rare)
- **Keycloak config:** exported weekly to S3

### 12.4 DR (Disaster Recovery)

**Pattern:** warm standby in `ap-south-2`

- EKS cluster in DR region, kept minimal (1 node per group)
- RDS cross-region read replica (promotable in DR event)
- S3 cross-region replication (active)
- Route 53 (or Cloudflare) health-check-based failover
- Cloudflare can route traffic to DR origin in minutes if primary region fails

**Targets:**

- RTO (Recovery Time Objective): 60 minutes for Tier-0
- RPO (Recovery Point Objective): 5 minutes for Tier-0

**DR drill cadence:** quarterly full failover exercise to DR region; verify full-path filings work.

### 12.5 Incident response

- On-call rotation (PagerDuty)
- Runbook per service in `/docs/runbooks/`
- Post-mortem template with blameless-culture norms
- SLO burn alerts: fast-burn (2% budget in 1hr) pages immediately; slow-burn (10% in 3 days) opens ticket

---

## 13. Performance & Scale Profile

### 13.1 Load characteristics

- **Baseline:** 500 RPS across all services, 50 IRNs/min, 10 filings/hour
- **Peak (filing days):** 5,000 RPS, 12,000 IRNs/min, 3,000 filings/min
- **Batch windows:** nightly 2A/2B pulls for every GSTIN (5,000 parallel flows in Year 1, 50,000 in Year 3)

### 13.2 Scaling strategy

- **Services:** HPA + KEDA (event-driven scaling on SQS backlog)
- **Database:** read replica for analytics; connection pooling via PgBouncer; future Aurora upgrade when write IOPS become the bottleneck
- **Cache:** ElastiCache cluster mode (adding shards rebalances online)
- **Storage:** S3 is effectively infinite

### 13.3 Performance budgets per service

| Service | P50 | P95 | P99 |
|---|---|---|---|
| Authentication | 50ms | 150ms | 300ms |
| IRN generation | 150ms | 300ms | 500ms |
| GSTR-1 save | 500ms | 2s | 5s |
| GSTR-1 file | 3s | 10s | 20s (GSTN-bound) |
| Reconciliation query | 300ms | 1s | 2s |
| CFO dashboard load | 800ms | 2s | 4s |
| Document upload (50MB) | 2s | 5s | 10s |

---

## 14. Repository Layout

Monorepo using Go workspaces + pnpm + Turborepo for a unified build.

```
complai/
├── apps/
│   ├── web/                     # Next.js 15 main product
│   ├── vendor-portal/           # Next.js 15 vendor portal
│   ├── complai-one/             # Next.js 15 + PWA for SMB
│   └── admin/                   # Internal admin console
├── services/
│   ├── go/
│   │   ├── identity-service/
│   │   ├── tenant-service/
│   │   ├── user-role-service/
│   │   ├── master-data-service/
│   │   ├── document-service/
│   │   ├── notification-service/
│   │   ├── audit-service/
│   │   ├── workflow-service/
│   │   ├── rules-engine-service/
│   │   ├── gst-service/
│   │   ├── gstr9-service/
│   │   ├── einvoice-service/
│   │   ├── ewb-service/
│   │   ├── tds-service/
│   │   ├── itr-service/
│   │   ├── vendor-service/
│   │   ├── recon-service/
│   │   ├── ap-service/
│   │   ├── billing-service/
│   │   ├── secretarial-service/
│   │   ├── reporting-service/
│   │   ├── gstn-gateway/
│   │   ├── irp-gateway/
│   │   ├── ewb-gateway/
│   │   ├── tds-gateway/
│   │   ├── itd-gateway/
│   │   ├── kyc-gateway/
│   │   ├── tax-payment-gateway/
│   │   ├── bank-gateway/
│   │   ├── mca-gateway/
│   │   └── erp-gateway/
│   ├── python/
│   │   ├── matching-ml-service/
│   │   ├── llm-copilot-service/
│   │   └── ocr-service/
│   └── node/
│       ├── web-bff-service/
│       ├── portal-bff-service/
│       └── smb-bff-service/
├── packages/
│   ├── shared-kernel-go/        # Go: tenant ctx, outbox, OTel, RLS helper, message bus
│   ├── shared-kernel-node/      # TS: formatters, Zod schemas, types
│   ├── ui-components/           # Complai component library
│   └── events/                  # Avro/Proto schemas, JSON schemas
├── infra/
│   ├── terraform/               # AWS infra as code
│   ├── helm/                    # K8s Helm charts
│   └── argocd/                  # GitOps manifests
├── docs/
│   ├── input/                   # PRD, architecture, design, API spec, build prompt
│   ├── adr/                     # Architecture decision records
│   ├── runbooks/                # Per-service operational runbooks
│   └── api/                     # OpenAPI specs per service
└── scripts/                     # Dev tooling
```

### 14.1 Standard Go service layout

```
services/go/{name}-service/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── api/                     # HTTP handlers (chi)
│   ├── app/                     # Application services (use cases)
│   ├── domain/                  # Domain models + interfaces
│   ├── infra/
│   │   ├── postgres/            # Repositories (sqlc-generated)
│   │   ├── redis/
│   │   ├── sqs/
│   │   └── temporal/            # Workflow + activity registrations
│   └── config/
├── migrations/                  # Goose-managed SQL migrations
├── api/                         # OpenAPI spec
├── test/                        # Integration tests (testcontainers)
├── Dockerfile
├── go.mod
└── Makefile
```

---

## 15. Tech Stack Reference

| Layer | Technology | Version |
|---|---|---|
| Cloud | AWS | `ap-south-1` + `ap-south-2` |
| Orchestration | Amazon EKS | 1.30 |
| Service mesh | Istio Ambient | 1.22 |
| Backend | Go | 1.22 |
| AI/ML | Python | 3.12 |
| Frontend | TypeScript | 5.4 |
| Frontend framework | Next.js / React | 15 / 19 |
| Styling | Tailwind + shadcn/ui | — |
| BFF | Node + NestJS | 20 / 10 |
| Database | Amazon RDS PostgreSQL | 16 |
| Cache | Amazon ElastiCache Redis | 7 |
| Object store | Amazon S3 | — |
| Messaging | Amazon SQS + SNS | — |
| Search | Amazon OpenSearch Service | 2 |
| Workflow | Temporal Cloud | Managed |
| Identity | Keycloak | 24 |
| CDN + DNS + WAF | Cloudflare | — |
| Email | Amazon SES | — |
| Observability | Last9 | Managed |
| Secrets | AWS Secrets Manager + KMS | — |
| CI | GitHub Actions | — |
| CD | ArgoCD | 2.11 |
| Terraform | AWS provider | 5.x |
| Container registry | Amazon ECR | — |

### 15.1 Go library stack

| Concern | Library |
|---|---|
| HTTP router | `go-chi/chi` v5 |
| DB driver | `jackc/pgx` v5 |
| DB queries | `sqlc` (SQL → typed Go) |
| Migrations | `pressly/goose` |
| Validation | `go-playground/validator` v10 |
| Config | `spf13/viper` + env vars |
| Logging | `rs/zerolog` |
| Dependency injection | `uber-go/fx` |
| Observability | OpenTelemetry Go SDK |
| Testing | `stretchr/testify` + `testcontainers/testcontainers-go` |
| Mocking | handwritten fakes + `gomock` where needed |
| Money arithmetic | `shopspring/decimal` |
| JWT | `golang-jwt/jwt` v5 |
| AWS SDK | `aws-sdk-go-v2` |
| Temporal | `go.temporal.io/sdk` |

---

## 16. Cost Model — Phase 1

Estimated monthly AWS + SaaS spend:

| Line item | USD/month |
|---|---|
| EKS cluster + ~10 application nodes | $1,800 |
| RDS Multi-AZ + read replica (`r7g.2xlarge`) | $1,500 |
| ElastiCache Redis cluster | $500 |
| OpenSearch 3-node | $700 |
| S3 storage + requests | $50 |
| SQS + SNS | $30 |
| NAT Gateway + data transfer | $200 |
| Secrets Manager + KMS | $100 |
| ALB + VPC endpoints | $50 |
| Cloudflare (Business plan) | $200 |
| Amazon SES (multi-million emails) | $100 |
| Temporal Cloud (Phase 1 tier) | $300 |
| Last9 (Phase 1 tier) | $500 |
| GuardDuty + Security Hub | $150 |
| **Primary region subtotal** | **~$6,180** |
| DR region (warm standby) | $1,500 |
| **Total Phase 1 monthly** | **~$7,680** |

Approximately **₹6.5 lakh/month** Phase 1. Scales roughly linearly with customer count for the first year. Year-3 estimate at 50,000 GSTINs: ~$40,000/month infra + bigger commercial negotiation on Adaequare/Sandbox/Last9.

---

## 17. Team Topology (reference)

A team suitable for this architecture at the end of Phase 1:

- **Founding engineers (2–3)** — full-stack, architect-level
- **Backend engineers (3–5)** — Go-proficient
- **Frontend engineers (2–3)** — TypeScript + design-system fluent
- **DevOps / SRE (1–2)** — AWS, EKS, Terraform, security
- **QA engineers (1–2)** — automation-first, domain-aware
- **ML engineer (1)** — Phase 4+
- **Product + design (2)** — PM + Product Designer

Total: ~12–15 engineering-adjacent people for the production build. Claude Code compresses the first-six-months of this journey into ~6 weeks.

---

## 18. Architecture Decision Records (summary)

Detailed ADRs live in `/docs/adr/`. Summary:

| ADR | Decision |
|---|---|
| ADR-0001 | Multi-tenancy via Postgres RLS |
| ADR-0002 | Adaequare Enriched APIs only (no pass-through) |
| ADR-0003 | Two-provider API strategy (Adaequare + Sandbox.co.in) |
| ADR-0004 | Go as primary backend language |
| ADR-0005 | AWS as cloud provider, `ap-south-1` primary |
| ADR-0006 | Postgres-only for Phase 1 (OLTP + analytics) |
| ADR-0007 | SQS/SNS over Kafka for Phase 1 messaging |
| ADR-0008 | Temporal Cloud (managed) |
| ADR-0009 | Cloudflare for CDN/DNS/WAF |
| ADR-0010 | Amazon SES for email |
| ADR-0011 | Last9 for observability |
| ADR-0012 | Keycloak self-hosted for identity |
| ADR-0013 | Outbox pattern via polling (not Debezium) in Phase 1 |
| ADR-0014 | Canonical Invoice Schema as lingua franca |
| ADR-0015 | Monorepo with Go workspaces + pnpm + Turborepo |

---

**End of Architecture v1.0. Approved for build.**
