# BUILD_PLAN.md — Living checklist

Last updated: 2026-04-25

## Current part
None yet. Next = Part 1.

## Completed
- [x] Part 0.5: Repo init, CLAUDE.md, BUILD_PLAN.md, input docs, ADR template

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
- [ ] Configure SES inbound email for per-tenant AP inboxes (`*@inbox.complai.in`)
- [ ] Set up SES → Lambda → S3 → SQS pipeline for email ingestion

### Observability (Part 3 / Part 14)
- [ ] Create Last9 account and obtain OTLP endpoint
- [ ] Configure OTel Collector DaemonSet to export to Last9
- [ ] Set up Last9 dashboards: per-service RED, per-provider, per-tenant, SLO
- [ ] Configure Last9 alerting → PagerDuty + Slack

### Workflow orchestration (Part 3)
- [ ] Create Temporal Cloud account (AWS Mumbai region)
- [ ] Provision namespaces: complai-filings, complai-reconciliation, complai-bulk, complai-ap, complai-onboarding
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
- [ ] HDFC ENet / ICICI iConnect: obtain sandbox API access (Part 13)
- [ ] SAP / Oracle NetSuite: obtain test system access (Part 13)

## Notes
- Docs in /docs/input are authoritative; don't re-derive
- ADRs in /docs/adr are authoritative for past decisions
- Every part ends with: tests green + BUILD_PLAN updated + commit
- All dev uses LocalStack — same code runs on real AWS via env var swap
- Terraform files are scaffolding for DevOps team — never run locally
