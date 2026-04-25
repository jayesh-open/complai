# Complai

Enterprise GST, TDS, and compliance SaaS platform for Indian businesses.

## Architecture overview

Complai is composed of ~28 services organized into five families:

- **Platform services (9)** -- identity, tenant, user/role, master data, document, notification, audit, workflow, rules engine. Shared infrastructure every domain module depends on.
- **Domain services (11)** -- one per compliance workflow: GST, GSTR-9, e-Invoice, E-Way Bill, TDS, ITR, vendor management, reconciliation, AP automation, billing, secretarial.
- **Gateway services (8+)** -- thin services that talk to external providers (Adaequare for GST/IRP/EWB, Sandbox.co.in for TDS/ITR/KYC, direct bank and ERP integrations). Each gateway normalizes responses into an internal contract.
- **AI/ML services (3)** -- Python services for reconciliation matching (CatBoost), LLM copilot, and OCR extraction. Phase 4+.
- **Web-facing services (4)** -- Node/NestJS BFFs for each frontend app, plus a Go reporting service for PDF/Excel/CSV generation.

## Tech stack

| Layer | Technology |
|---|---|
| Cloud | AWS (ap-south-1 primary, ap-south-2 DR) |
| Compute | EKS 1.30 + Istio 1.22 ambient mesh |
| Backend | Go 1.22 (domain + gateways) |
| AI/ML | Python 3.12 |
| Frontend | TypeScript 5.4 + Next.js 15 + React 19 + Tailwind + shadcn/ui |
| BFF | Node 20 + NestJS |
| Database | Amazon RDS PostgreSQL 16 (Multi-AZ + read replica) |
| Cache | ElastiCache Redis 7 |
| Messaging | SQS + SNS (MessageBus abstraction) |
| Search | Amazon OpenSearch Service 2 |
| Workflow | Temporal Cloud (managed) |
| Identity | Keycloak 24 (self-hosted on EKS) |
| CDN/DNS/WAF | Cloudflare |
| Email | Amazon SES |
| Observability | Last9 (OpenTelemetry-native) |
| Secrets | AWS Secrets Manager + KMS |
| CI/CD | GitHub Actions + ArgoCD |

## Prerequisites

- Node.js >= 20.18
- Go >= 1.22
- Docker and Docker Compose
- pnpm >= 9

## Quick start

```bash
# Clone the repository
git clone git@github.com:your-org/complai.git
cd complai

# Install Node/TypeScript dependencies
pnpm install

# Start local infrastructure (Postgres, Redis, LocalStack, Mailpit, Temporal dev server, Jaeger)
make dev

# Run database migrations for all services
make migrate

# Start a specific Go service (example: tenant-service)
cd services/go/tenant-service && make run

# Start the web app
cd apps/web && pnpm dev

# Start the BFF
cd services/node/web-bff-service && pnpm dev
```

### Local dev environment

All AWS service interactions use LocalStack (configured in `docker-compose.dev.yml`). Go services use `aws-sdk-go-v2` with `AWS_ENDPOINT_URL` pointed at LocalStack. Mailpit replaces SES, Temporal dev server replaces Temporal Cloud, Jaeger replaces Last9 for tracing.

## Repo layout

```
complai/
├── apps/
│   ├── web/                        # Next.js 15 -- main product
│   ├── vendor-portal/              # Next.js 15 -- external vendor app
│   ├── complai-one/                # Next.js 15 + PWA -- SMB billing
│   └── admin/                      # Internal admin console
├── services/
│   ├── go/
│   │   ├── identity-service/       # Auth, MFA, SSO, sessions
│   │   ├── tenant-service/         # Tenant lifecycle, PAN/GSTIN/TAN hierarchy
│   │   ├── user-role-service/      # Roles, permissions, approval workflows
│   │   ├── master-data-service/    # Vendors, customers, employees, HSN, items
│   │   ├── document-service/       # S3 + Postgres metadata, OCR trigger
│   │   ├── notification-service/   # Email (SES), WhatsApp, SMS, in-app
│   │   ├── audit-service/          # Tamper-evident audit log
│   │   ├── workflow-service/       # Temporal Cloud integration layer
│   │   ├── rules-engine-service/   # Tax determination, validation rules
│   │   ├── gst-service/            # GSTR-1, GSTR-3B, GSTR-2A/2B/IMS
│   │   ├── gstr9-service/          # Annual returns (GSTR-9, GSTR-9C)
│   │   ├── einvoice-service/       # IRN generation, signed JSON/QR
│   │   ├── ewb-service/            # E-Way Bill lifecycle
│   │   ├── tds-service/            # TDS/TCS deductees, payments, returns
│   │   ├── itr-service/            # ITR filing flows
│   │   ├── vendor-service/         # Vendor master, compliance scoring
│   │   ├── recon-service/          # 5-stage match pipeline, IMS actions
│   │   ├── ap-service/             # Invoice ingestion, 3-way match, approvals
│   │   ├── billing-service/        # Complai One SMB billing
│   │   ├── secretarial-service/    # ROC filings, registers, minutes
│   │   ├── reporting-service/      # PDF/Excel/CSV report generation
│   │   ├── gstn-gateway/           # Adaequare -- GSTR returns, IMS, ledgers
│   │   ├── irp-gateway/            # Adaequare -- e-Invoice IRN
│   │   ├── ewb-gateway/            # Adaequare -- E-Way Bill
│   │   ├── tds-gateway/            # Sandbox -- TDS calc, FVU, e-File
│   │   ├── itd-gateway/            # Sandbox -- ITR filing, prefill, AIS
│   │   ├── kyc-gateway/            # Sandbox -- PAN, Aadhaar, GSTIN, bank
│   │   ├── tax-payment-gateway/    # Sandbox -- TDS and income tax challans
│   │   ├── bank-gateway/           # Direct -- payments and collections
│   │   ├── mca-gateway/            # Direct -- MCA21 ROC filings
│   │   └── erp-gateway/            # Direct -- SAP, Oracle, Tally, Dynamics
│   ├── python/
│   │   ├── matching-ml-service/    # CatBoost model for recon Stage 3
│   │   ├── llm-copilot-service/    # NL queries, explanation generation
│   │   └── ocr-service/            # Form 16 / invoice extraction
│   └── node/
│       ├── web-bff-service/        # BFF for apps/web
│       ├── portal-bff-service/     # BFF for apps/vendor-portal
│       └── smb-bff-service/        # BFF for apps/complai-one
├── packages/
│   ├── shared-kernel-go/           # Go: tenant ctx, outbox, OTel, RLS, message bus
│   ├── shared-kernel-node/         # TS: formatters, Zod schemas, types
│   ├── ui-components/              # Complai component library
│   └── events/                     # Protobuf schemas, JSON schemas
├── infra/
│   ├── terraform/                  # AWS infra as code
│   ├── helm/                       # K8s Helm charts
│   └── argocd/                     # GitOps manifests
├── docs/
│   ├── input/                      # PRD, architecture, design, API spec
│   ├── adr/                        # Architecture decision records
│   ├── runbooks/                   # Per-service operational runbooks
│   └── api/                        # OpenAPI specs per service
├── scripts/                        # Dev tooling
├── go.work                         # Go workspace root
├── pnpm-workspace.yaml             # pnpm workspace config
├── turbo.json                      # Turborepo config
└── docker-compose.dev.yml          # Local dev infrastructure
```

### Standard Go service layout

```
services/go/{name}-service/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── api/                        # HTTP handlers (chi)
│   ├── app/                        # Application services (use cases)
│   ├── domain/                     # Domain models + interfaces
│   ├── infra/
│   │   ├── postgres/               # Repositories (sqlc-generated)
│   │   ├── redis/
│   │   ├── sqs/
│   │   └── temporal/               # Workflow + activity registrations
│   └── config/
├── migrations/                     # Goose-managed SQL migrations
├── api/                            # OpenAPI spec
├── test/                           # Integration tests (testcontainers)
├── Dockerfile
├── go.mod
└── Makefile
```

## Development workflow

1. Create a feature branch from `main`.
2. Make changes. Run `make lint` and `make test` locally.
3. Open a pull request. CI runs:
   - Affected-service detection (path filters).
   - Per-service build, lint, unit tests, integration tests.
   - Trivy container scan, Snyk dependency scan, GitLeaks secret scan.
   - Helm chart render + diff.
4. Code review and approval.
5. Merge to `main`.
6. ArgoCD syncs staging automatically. Production requires manual approval.

## Key references

- [Architecture Decision Records](docs/adr/)
- [Product Requirements Document](docs/input/complai_prd.md)
- [Technical Architecture](docs/input/complai_architecture.md)
- [Design System](docs/input/complai_design_system.md)
- [API Integration Spec](docs/input/complai_api_integration.md)
