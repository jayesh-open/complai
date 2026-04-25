# Complai — Product Requirements Document

**Version:** 1.0
**Status:** Approved for build
**Owner:** Founding Team
**Date:** April 2026

---

## 1. Product Overview

### 1.1 What Complai is

Complai is an enterprise-grade compliance SaaS platform for Indian businesses. It is the single operating system for a company's regulatory obligations: GST returns, e-Invoicing, e-Way Bills, TDS/TCS, income tax, secretarial compliance, accounts payable automation, vendor management, and invoice discounting.

Complai targets mid-market and enterprise Indian businesses (₹100 Cr to ₹10,000 Cr revenue) that today either stitch together 4–8 point solutions or run ad-hoc spreadsheet-based processes with dedicated compliance teams. The platform replaces the stack with one unified product.

### 1.2 Who it serves

Primary customers: mid-to-large Indian enterprises, including listed companies, PE-backed companies, manufacturing, e-commerce, BFSI adjacent businesses (NBFCs, brokers), and rapidly-growing startups crossing ₹50 Cr revenue.

Primary users within these customers:
- **Tax Managers** (filing GST, TDS, ITR)
- **Tax Analysts** (data preparation, reconciliation)
- **AP Clerks** (invoice processing)
- **AP Approvers** (payment approvals)
- **CFOs / Finance Controllers** (oversight dashboards)
- **Company Secretaries** (ROC filings, minutes, registers)
- **Auditors** (read access for statutory audits)
- **External CAs** (assisted filings via marketplace)
- **Vendors** (external portal for invoice submission and compliance visibility)

### 1.3 Why Complai exists

Indian compliance is uniquely complex — 20+ return forms, 50+ TDS sections, state-specific rules, monthly filing deadlines, and continuously evolving regulations (GSTR-9 format changes yearly, IMS went live in Oct 2024, e-Invoice thresholds keep expanding). Existing solutions fall into three buckets:

- **Point products** (one for GST, one for TDS, one for AP) — require integration work and have no cross-module intelligence
- **ERPs with compliance modules** (SAP, Oracle) — expensive, rigid, poor UX, slow to adapt to regulatory changes
- **Legacy compliance vendors** — functional but built on outdated architecture, poor UX, limited automation

Complai is built for the post-2024 Indian regulatory era: API-first, AI-augmented, enterprise-grade from day one, with a modern interface that treats compliance as a product category rather than a back-office utility.

### 1.4 North-star metrics

- **Adoption:** 500 enterprise customers by end of Year 1; 50,000 GSTINs under management by end of Year 3
- **Coverage:** 10+ compliance workflows per customer (not just GST)
- **Reliability:** 99.99% SLO on Tier-0 services during filing peaks
- **Filing compliance:** >99.5% of customer filings submitted on time
- **NPS:** >60 among Tax Managers

---

## 2. Product Scope — The 11 Modules

Complai organizes into eleven functional modules plus a horizontal platform layer. Each module is independently valuable; the platform's power comes from their integration.

### 2.1 Module overview

| # | Module | Primary Users | Why it matters |
|---|---|---|---|
| 1 | **GST Compliance** | Tax Managers, Analysts | Core monthly/annual filings for every business |
| 2 | **E-Invoicing (IRN)** | Tax Analysts, AR teams | Mandatory for ₹5 Cr+ turnover |
| 3 | **E-Way Bill** | Logistics, Tax Analysts | Mandatory for goods movement ₹50K+ |
| 4 | **MaxITC & Vendor Compliance** | Tax Managers, CFOs | Protects ₹ lakhs-crores of input tax credit |
| 5 | **TDS / TCS** | Tax Analysts, Payroll | Quarterly filings, Form 16 generation |
| 6 | **Income Tax (ITR)** | Employees, CAs, HR | Employer-bulk filing, CA marketplace |
| 7 | **AP Automation** | AP Clerks, Approvers, CFO | Replaces paper/email invoice workflows |
| 8 | **Invoice Discounting** | CFOs, vendors | Working capital via TReDS |
| 9 | **Complai One (SMB)** | SMB founders, bookkeepers | Simplified GST+billing for small businesses |
| 10 | **Vendor Management** | Procurement, Tax, AP | Unified vendor master + portal |
| 11 | **Secretarial (Compliance Cloud)** | Company Secretaries | ROC filings, registers, minutes |

Plus the **Platform Layer** — identity, tenancy, users/roles, master data, documents, notifications, audit, workflow, rules engine — shared across all 11 modules.

---

## 3. Module 1 — GST Compliance

### 3.1 Purpose

End-to-end management of a business's GST obligations: monthly GSTR-1 and GSTR-3B, annual GSTR-9 and GSTR-9C, quarterly CMP-08 (composition), and associated returns (4, 5, 6, 7, 8, ITC-04, DRC-03, RFD-01).

### 3.2 Sub-modules

#### GSTR-1 (Outward Supply)

**What it does:** ingests the business's sales register (from ERP, Excel, or manual upload), categorizes invoices into 11 GSTR-1 sections (B2B, B2CL, B2CS, CDNR, CDNUR, EXP, AT, ATADJ, NIL, HSN, DOCS), validates against GST rules, and files the return with DSC or EVC.

**Features:**
- Sales register ingestion (CSV, Excel, ERP push, email-based)
- Auto-categorization via rules engine (place-of-supply, registration type, invoice type)
- Validation: GSTIN format, HSN validity, state-wise place-of-supply, reverse-charge identification
- Government-form replica review screens (line-by-line match to the GSTR-1 form)
- Amendments for prior periods (Tables 9, 10, 11)
- Recipient-initiated amendments via GSTR-1A
- Filing with DSC token or Aadhaar-OTP based EVC
- ARN capture and acknowledgment archival

**How it works:**
1. Tax analyst imports sales register
2. System auto-categorizes, flags exceptions (missing HSN, invalid GSTIN, etc.)
3. Analyst resolves exceptions
4. Tax manager reviews section-by-section in government-form layout
5. Maker-checker approval if enterprise policy requires
6. File with DSC/EVC; system captures ARN and signed acknowledgment
7. Full audit trail retained

#### GSTR-3B (Summary + Tax Payment)

**What it does:** consolidates outward supply (from GSTR-1), inward supply and ITC (from GSTR-2B + IMS), computes liability, offsets via cash and credit ledgers, and files.

**Features:**
- Auto-populate from filed GSTR-1 and latest GSTR-2B/IMS
- System-computed tax liability with user-override + justification
- Cash ledger and credit ledger balance view
- Offset wizard: maximize credit utilization within rules
- Challan generation if cash shortfall
- Filing + ARN capture

#### GSTR-2A / GSTR-2B / IMS (Inward Supply Reconciliation)

**What it does:** pulls data from the government systems, reconciles against the business's purchase register, and enables actioning via IMS.

**GSTR-2A:** dynamic — updates as suppliers file. Pulled section-by-section (B2B, B2BA, CDN, CDNA, ISD, ISDA, IMPG, IMPGSEZ, TDS, TCS).

**GSTR-2B:** static monthly snapshot (generated 14th of each month).

**IMS (Invoice Management System, live Oct 2024):** recipient actions (Accept/Reject/Pending) on every supplier invoice. Actions feed into GSTR-2B regeneration.

**Features:**
- Scheduled pulls (nightly for 2A, monthly for 2B)
- 5-stage reconciliation pipeline (exact, fuzzy, AI-assisted, partial, unmatched)
- Bucket views (Matched, Mismatch, Partial, Missing-2B, Missing-PR, Duplicate)
- IMS action UI with bulk operations
- Reason codes for every mismatch
- AI-suggested matches (Module 12 integration)

#### GSTR-9 / GSTR-9C (Annual Return + Reconciliation)

**What it does:** auto-populates annual return from 12 months of GSTR-1 and GSTR-3B, reconciles against books of accounts, supports auditor sign-off on 9C.

**Features:**
- Auto-population from filed returns
- Table 8A bulk document-level pull (often 100K+ rows)
- Difference report (books vs returns)
- Auditor collaboration: CA/auditor can review, annotate, and sign 9C
- Filing with DSC

#### Other returns

- **GSTR-4** — composition dealer annual return
- **GSTR-5** — non-resident taxable person
- **GSTR-6 / 6A** — input service distributor
- **GSTR-7** — TDS deductor under GST
- **GSTR-8** — TCS by e-commerce operators
- **CMP-08** — composition quarterly
- **ITC-04** — job-work
- **GST RFD-01** — refunds (with statement uploads)
- **GST DRC-03** — voluntary payment

Each follows the same pattern: prepare → validate → submit → file → acknowledge.

### 3.3 Non-functional requirements

- Filing confirmation modal with type-to-confirm for any filing where tax > ₹10 lakh
- DSC + EVC both supported as signing methods
- 8-year retention of signed acknowledgments (regulatory requirement)
- Filing must complete in under 5 minutes end-to-end during peak (10th/11th of every month)
- Must handle 12,000+ concurrent filings during peak

---

## 4. Module 2 — E-Invoicing (IRN Generation)

### 4.1 Purpose

Generate Invoice Reference Numbers (IRNs) from the NIC Invoice Registration Portal for every B2B invoice and credit/debit note, as mandated for businesses with turnover above ₹5 Cr (threshold keeps reducing).

### 4.2 Features

- Single IRN generation (user-initiated from an invoice)
- Bulk IRN generation (up to 1,000 invoices per request; progress tray in UI)
- IRN retrieval by IRN, by document details, by EWB number
- Cancellation within 24-hour window (with reason codes)
- Signed JSON + Signed QR code persistence (8-year retention)
- Automatic ITR data + EWB generation from IRN (single flow)
- Master data validations: supplier/buyer GSTIN, HSN, state, country, port, currency, UQC
- QR code extraction and verification (for incoming supplier e-Invoices)
- e-Invoice schema v1.1 compliance with transformer from canonical invoice schema

### 4.3 Integration points

- Integrates with AR systems (invoice created in ERP → IRN back to ERP)
- Integrates with AP automation (incoming supplier e-Invoice QR verified)
- Feeds GSTR-1 (auto-populated from e-Invoice data per GSTN)

---

## 5. Module 3 — E-Way Bill

### 5.1 Purpose

Generate E-Way Bills for goods movement above ₹50,000, as mandated by state GST authorities.

### 5.2 Features

- Single and bulk EWB generation (up to 100 per bulk request)
- Generate EWB directly from IRN (no duplicate payload — single-flow)
- Part-A + Part-B EWB support
- Multi-vehicle movement (long-distance shipments crossing multiple states — initiate → add vehicle → change vehicle)
- Consolidated EWBs (multiple EWBs on one vehicle / one trip sheet)
- Validity extension with reason codes (in-transit, transshipment, natural calamity)
- Cancellation (24-hour window) + rejection (72-hour window for recipient)
- Retrieval by EWB number, by date, by consignor/consignee/transporter, by GSTIN
- Transporter master management
- Distance calculation (PIN to PIN) for validity
- HSN and place master validation

### 5.3 Non-functional

- EWB validity timer visible in UI (green >24h, amber 4–24h, red <4h)
- Alerts on approaching expiry
- Auto-extension workflow if enabled by enterprise

---

## 6. Module 4 — MaxITC & Vendor Compliance

### 6.1 Purpose

Protect Input Tax Credit (ITC) by ensuring every vendor is compliant and every invoice is reconciled before it impacts GSTR-3B ITC claims. For a typical mid-market enterprise, ITC represents ₹2–20 crore of working capital annually — losing eligibility because of vendor non-compliance is a real and common risk.

### 6.2 Features

#### Vendor compliance scoring

- 100-point score across 5 dimensions:
  - Filing regularity (30 points) — how consistently the vendor files GSTR-1/3B
  - IRN compliance (20 points) — does the vendor generate IRNs properly
  - Mismatch rate (20 points) — historical reconciliation gap rate
  - Payment behavior (15 points) — pays on time, responds to queries
  - Document hygiene (15 points) — has GST certificate, bank details, PAN, etc.
- Category assignment: A (≥80), B (60–79), C (40–59), D (<40)
- 12-month rolling window
- Automatic vendor-level dashboards

#### ITC orchestration

- Scheduled nightly pull of GSTR-2A for every GSTIN
- Monthly 2B snapshot capture
- IMS pull and action management
- Cross-referencing: which invoices have no matching 2B entry? Which invoices have supplier GSTR-1 mismatches?
- Automated chase emails (via LLM copilot) to delinquent vendors
- ITC-at-risk dashboard: ₹ amount, by vendor, by aging

#### Vendor portal

- External-facing portal
- Magic-link authentication (no passwords for external users)
- Vendor sees: their compliance score, invoices raised to the buyer, payment status, pending documents
- Self-service updates: bank account change, address change (triggers re-KYC)

### 6.3 Intelligence features (Phase 2+)

- AI-suggested matches for reconciliation ambiguities
- Anomaly detection: unusual mismatch patterns, sudden vendor non-compliance
- Predictive vendor scoring based on behavioral signals

---

## 7. Module 5 — TDS / TCS

### 7.1 Purpose

End-to-end TDS and TCS management: deductee master, payment recording, challan management, quarterly return filing (24Q/26Q/27Q/27EQ), Form 16 and 16A bulk generation, correction statements.

### 7.2 Features

#### Transaction-level TDS

- TDS calculation across all sections (194C, 194J, 194Q, 194I, 194BA, 194O, and every other active section)
- 206AB/206CCA specified-person compliance check at payment time
- TDS computation on salary (per 24Q Annexure-II)

#### Deductee management

- Deductee master with PAN validation, 206AB status, higher-rate applicability
- Bulk import from payroll / AP
- Status tracking (active, inactive, deceased)

#### Challan and payment

- Challan generation (ITNS-281)
- Online payment via authorized banks
- CIN capture and reconciliation
- Payment history per deductor

#### Quarterly returns

- 24Q (salary), 26Q (non-salary resident), 27Q (non-resident), 27EQ (TCS)
- Generate TXT file (per NSDL RPU schema)
- Download CSI (Challan Status Inquiry) from portal
- Generate FVU (File Validation Utility output) + Form 27A
- E-file the return
- E-verify with Form 27A upload
- Receipt (RRR) capture

#### Form 16 / 16A generation

- Bulk Form 16 generation post-Q4 (24Q) — Part A from TRACES + Part B system-generated
- Bulk Form 16A generation quarterly for 26Q/27Q/27EQ
- Email delivery to employees/deductees
- Download zip + individual PDFs

#### Correction statements

- Support for defective returns
- Correction type classification (C1/C2/C3/C4/C5/Y)
- Resubmission and status tracking

### 7.3 Non-functional

- Full quarter processing time <30 min (500 deductees, 2,000 payments)
- OTP-based CSI download flow (human-in-the-loop at CSI step only)
- 7-year retention of filed returns

---

## 8. Module 6 — Income Tax (ITR)

### 8.1 Purpose

Two distinct workflows:
1. **Employer-initiated bulk ITR filing** for employees (post-Form 16 dispatch, employees file in under 3 minutes via a magic-link)
2. **CA-assisted marketplace** where complex cases are routed to empaneled CAs who file on behalf of the employee (Complai takes a referral + transaction fee)

### 8.2 Features

#### Employer-bulk flow

- Post-Q4 Form 16 bulk dispatch via SES
- Email includes a magic link to a pre-filled ITR flow
- Employee clicks link → sees auto-filled personal, income, tax details → reviews → consents → e-verifies → done
- Supports ITR-1 (Sahaj) and ITR-2 for most cases

#### CA-assisted marketplace

- Empaneled CA network
- Complex ITR routing (capital gains, foreign income, crypto, multi-source)
- CA dashboard: client list, batch filing, fee collection
- Complai collects platform fee per filing

#### Full ITR forms supported

ITR-1 through ITR-7, including:
- Pre-fill from AIS, 26AS, Form 16
- Regime comparator (old vs new) with recommendation
- Capital gains computation (domestic, foreign, crypto)
- Validation against ITD schema
- E-file via ERI channel
- ITR-V generation + e-verification (Aadhaar OTP, EVC, NetBanking, DSC)
- Refund status tracking

### 8.3 Non-functional

- Employee flow completion target: under 3 minutes
- Mobile-responsive (most employees file from mobile)
- Support Hindi + English (Phase 2+)

---

## 9. Module 7 — AP Automation

### 9.1 Purpose

Automate the accounts payable workflow: vendor invoice receipt → OCR extraction → 3-way match (PO + GRN + Invoice) → approval workflow → payment execution.

### 9.2 Features

#### Invoice ingestion

- Per-tenant dedicated email inbox (e.g., `ap@acme.complai.com`)
- Email attachment auto-extraction
- Drag-drop upload
- ERP push (SAP, Oracle, Tally, Dynamics, NetSuite)
- Vendor portal submission
- Mobile capture (photograph → upload)

#### OCR and extraction

- High-accuracy OCR trained on Indian invoice formats
- Field-level confidence scoring (user reviews low-confidence fields)
- GST detail extraction (GSTIN, HSN, tax amounts)
- Multi-page invoice support
- Auto-linking to existing vendor master

#### Matching and approval

- 3-way match: PO + GRN (Goods Receipt Note) + Invoice
- Configurable tolerance (amount delta, quantity delta, date delta)
- Approval workflows (N-of-M approver chains, per-tenant configurable)
- Exception policies: auto-approve rules (under ₹X, trusted vendor, etc.)
- Maker-checker separation enforced

#### Payment execution

- Payment file generation (ISO 20022 pain.001 standard + HDFC/ICICI/SBI/Axis specific formats)
- Direct bank API submission where supported
- Payment reconciliation (CIN/UTR capture, marking paid)
- Vendor payment advice dispatch

### 9.3 Integrations

- Bank APIs (Phase 13)
- ERP bidirectional sync
- e-Invoice QR verification for incoming invoices (feeds vendor compliance)

---

## 10. Module 8 — Invoice Discounting

### 10.1 Purpose

Help vendors access working capital by discounting unpaid-but-approved invoices through TReDS (Trade Receivables Discounting System) — RXIL, Invoicemart, M1exchange.

### 10.2 Features

- Vendor invoice eligibility check (approved in buyer's system, above minimum amount, buyer is eligible anchor)
- Submission to one or multiple TReDS platforms
- Bid management (vendors see bids from various financiers)
- Acceptance flow with buyer concurrence
- Settlement tracking
- Dashboard per vendor and per buyer

### 10.3 Business model

Complai takes a platform fee from successful discounting transactions. Creates stickiness for vendor usage of the platform.

---

## 11. Module 9 — Complai One (SMB Billing)

### 11.1 Purpose

Simplified billing and GST compliance app for SMBs (under ₹50 Cr revenue). Mobile-first, 5-field invoice creation, integrated GST filing, payment collection.

### 11.2 Target users

Small manufacturers, retailers, service providers, consultants, freelancers. Users who today use Tally, Zoho Books, or QuickBooks — but want lower cost and integrated GST.

### 11.3 Features

- Super-simple invoice creation (5 fields: customer, item, amount, date, GST rate)
- Mobile PWA (React Native app in Phase 4+)
- Automatic IRN generation (if above e-Invoice threshold)
- Automatic EWB if shipment
- Customer management
- Item master
- Recurring invoices
- Payment links (via Razorpay, Cashfree)
- Basic GSTR-1 filing for under-threshold taxpayers
- Simple reports (monthly sales, GST payable, outstanding)

### 11.4 Monetization

- Free tier: 10 invoices/month
- Pro tier: ₹499/month — unlimited invoices, IRN, EWB, GSTR-1 filing
- Business tier: ₹1,499/month — includes payment collection, recurring invoices, team access

---

## 12. Module 10 — Vendor Management

### 12.1 Purpose

Unified vendor master, KYB (Know Your Business) workflow, and compliance hub for all vendor-facing processes.

### 12.2 Features

- Vendor CRUD with multi-contact, multi-address, multi-bank-account support
- Onboarding workflow: Basic info → KYC (PAN, GSTIN, bank, Udyam for MSME, MCA for company data) → Compliance check → Score → Approval (maker-checker for high-risk)
- Bulk import (Excel/CSV, up to 10,000 rows, with column mapping and error reporting)
- Change history with approval trails
- Integration with AP (invoice-to-vendor linkage)
- Integration with MaxITC (compliance scoring)
- Vendor portal (external, magic-link)

### 12.3 KYC coverage

- PAN verification (individual + business)
- GSTIN verification with 12-month filing history
- TAN verification
- Bank account verification (penny-drop + reverse penny-drop via IMPS)
- Aadhaar OTP (where required for director KYC)
- MCA company/LLP/director master
- Udyam (MSME) verification
- DigiLocker for document fetch

---

## 13. Module 11 — Compliance Cloud (Secretarial)

### 13.1 Purpose

Company Secretary / legal compliance for corporates and LLPs — ROC filings, board/shareholder resolutions, statutory registers, director KYC.

### 13.2 Features

- Entity registry (companies, LLPs, directors, DINs)
- Filing calendar with statutory deadlines
- Form preparation and filing via MCA21 V3:
  - AOC-4 (annual financial statement)
  - MGT-7 (annual return)
  - DIR-3 KYC (director KYC)
  - ADT-1 (auditor appointment)
  - CHG-1 (charge creation)
  - INC-22 (registered office change)
  - Many others
- Statutory registers (members, directors, charges)
- Resolution library (templates + past resolutions)
- Minutes management
- Compliance health score per entity
- Auditor / CS collaboration

### 13.3 Notes

- MCA21 integration is direct (no aggregator covers this fully)
- Supports companies with 1 to 500+ subsidiaries

---

## 14. Platform Layer

Shared services underneath all 11 modules.

### 14.1 Identity & Access

- Keycloak-based identity
- SSO via SAML 2.0 and OIDC (Google Workspace, Azure AD, Okta)
- MFA (TOTP, SMS, email)
- Step-up authentication for filing operations (re-auth within 5-minute window)
- Session management with device tracking

### 14.2 Tenant & Entity Management

- Multi-tenant with data residency tags
- Tenancy tiers: Pooled, Bridge, Silo, On-Premise
- Tenant hierarchy: Tenant → PANs → GSTINs, TANs, CINs
- Per-tenant branding (logo, colors, email sender)
- Per-tenant feature flags

### 14.3 Users & Roles

- Role templates:
  - Tenant Admin
  - Tax Manager
  - Tax Analyst
  - AP Clerk
  - AP Approver
  - Viewer
  - Vendor (external)
  - CA (external)
  - Internal Admin (Complai staff)
- Permission model: resource-action pairs (e.g., `gst.gstr1.file`, `vendor.create`)
- Maker-checker workflows (N-of-M approver chains, per-resource configurable)
- Fine-grained permission overrides

### 14.4 Master Data

- Vendors, customers, employees, items, chart of accounts, HSN codes, state codes, pincodes, bank branches
- Tenant-scoped with RLS enforcement
- Change-data-capture for real-time sync
- ERP sync (bidirectional)

### 14.5 Document Management

- Invoice PDFs, signed JSONs, FVU files, challans, Form 16 PDFs, resolutions, minutes
- S3-backed with metadata in Postgres (JSONB for flexible fields)
- Tenant-scoped DEK (Data Encryption Key) envelope encryption via AWS KMS
- Retention policies per document class (invoice: 8 years, signed return: 8 years, TDS: 7 years)
- Virus scanning on upload
- Pre-signed URL access (15-min TTL)
- Document lineage graph

### 14.6 Notifications

- Email (Amazon SES)
- WhatsApp Business (Meta)
- SMS (MSG91)
- In-app (server-sent events + WebSocket)
- Per-user channel preferences + quiet hours
- Daily digest bundling (09:00 IST)
- Template registry with Handlebars
- DPDP consent tracking + unsubscribe

### 14.7 Audit Trail

- Every state-change event across every service captured
- Tamper-evident (hourly Merkle-chain hashing)
- Search by tenant, user, resource, action, date
- Export to PDF (signed) for regulatory submission
- 30-day hot retention (OpenSearch); long-term archive to S3 with Iceberg metadata

### 14.8 Workflow Orchestration

- Temporal Cloud for all long-running workflows (filing sagas, reconciliation, bulk IRN, Form 16 generation, vendor onboarding)
- Human-task integration (approval, OTP entry, DSC selection)
- Retry and compensation policies per workflow

### 14.9 Rules Engine

- Rules as versioned JSON per tenant
- Categories: tax determination (CGST/SGST/IGST split), HSN applicability, TDS applicability, validation rules, eligibility rules (ITC, RCM, composition)
- Seeded with FY 2026-27 Indian tax rules
- Hot-reload (rule changes propagate without service restart)

### 14.10 AI & Intelligence Layer (Phase 12+)

- ML-based invoice matching (CatBoost) for reconciliation Stage 3
- LLM copilot (Azure OpenAI or Bedrock) for vendor communication drafting, natural-language dashboard queries, return-prep explanation
- PII redaction before LLM calls
- Per-tenant usage caps

---

## 15. Non-Functional Requirements

### 15.1 Scale targets

| Year | Enterprises | GSTINs | Peak IRN/min | Peak filings/min | Peak recon ops/sec |
|---|---|---|---|---|---|
| Year 1 | 500 | 5,000 | 2,000 | 500 | 1,000 |
| Year 2 | 2,000 | 20,000 | 6,000 | 1,500 | 3,000 |
| Year 3 | 5,000 | 50,000 | 12,000 | 3,000 | 6,000 |

### 15.2 Service Level Objectives

**Tier-0 services** (direct customer impact during filing peaks): 99.99% availability, <300ms P95 latency

- IRN generation
- EWB generation
- GSTR-1/3B filing
- TDS e-filing
- Authentication
- Core data APIs

**Tier-1 services** (important but not filing-critical): 99.9% availability, <1s P95

- Reconciliation
- Document upload
- Notifications
- Reports

**Tier-2 services** (background/analytics): 99.5% availability, <3s P95

- AI/ML inference
- Audit search
- Dashboard aggregations

### 15.3 Security & Compliance

- **DPDP Act** (Digital Personal Data Protection Act, India) compliant from day one
- **SOC 2 Type II** within 12 months of production
- **ISO 27001** within 18 months
- **Row-level tenant isolation** (Postgres RLS)
- **Tenant-scoped encryption** (per-tenant DEKs via AWS KMS)
- **Data residency:** all data stored in India (`ap-south-1` primary, `ap-south-2` DR)
- **No PII in logs, traces, or LLM prompts** without explicit redaction
- **DSC and EVC** support for regulatory signing
- **Audit-trail immutability** via Merkle-chain hashing

### 15.4 Accessibility

- WCAG 2.2 Level AA compliance
- Keyboard-only completion for every critical workflow (login, filing, reconciliation, approvals)
- Screen reader support for financial tables with proper ARIA labels
- Color-plus-text for every status indicator (never color-only)

### 15.5 Localization

- **Day 1:** English, Hindi
- **Phase 2:** Tamil, Marathi, Gujarati, Kannada, Telugu, Bengali
- DD/MM/YYYY date format, ₹ currency prefix, Indian numbering (crore/lakh)

### 15.6 Performance

- Page load (authenticated pages): <2s P95
- Dashboard render: <3s P95 including data fetch
- Data tables: render 10,000 rows without virtualization lag
- Filing flow: complete in <5 min end-to-end during peak

---

## 16. Tech Stack Summary

- **Cloud:** AWS, primary `ap-south-1` (Mumbai), DR `ap-south-2` (Hyderabad)
- **Compute:** Amazon EKS + Istio ambient mesh
- **Backend:** Go 1.22 (all domain and gateway services); Python 3.12 (AI/ML only)
- **Frontend:** TypeScript 5.4 + Next.js 15 + React 19 + Tailwind + shadcn/ui
- **Primary database:** Amazon RDS for PostgreSQL 16 (OLTP + Phase 1 analytics via read replica)
- **Cache:** Amazon ElastiCache for Redis 7
- **Binary storage:** Amazon S3
- **Messaging:** Amazon SQS + SNS (fan-out where needed)
- **Search:** Amazon OpenSearch Service 2
- **Identity:** Keycloak on EKS
- **Workflow:** Temporal Cloud (managed)
- **CDN + DNS + WAF:** Cloudflare
- **Email:** Amazon SES
- **Observability:** Last9
- **Secrets:** AWS Secrets Manager + KMS
- **CI/CD:** GitHub Actions + ArgoCD

---

## 17. Go-to-Market & Build Phasing

### Phase 1 (Weeks 1–8): Core MVP
Modules: Platform layer + GST (1) + e-Invoicing (2) + EWB (3) + Vendor Compliance (4) — limited to core reconciliation

### Phase 2 (Weeks 9–14): TDS + ITR + AP
Modules: TDS (5) + ITR (6) + AP Automation (7)

### Phase 3 (Weeks 15–20): SMB + Vendor Mgmt + Discounting + Secretarial
Modules: Complai One (9) + Vendor Management (10) + Invoice Discounting (8) + Compliance Cloud (11)

### Phase 4 (Weeks 21–24): AI layer + Production hardening
AI intelligence, MaxITC orchestration, observability, security, DR, go-live readiness

---

## 18. Success Criteria — End of Phase 1

- 5 pilot enterprise customers onboarded
- 50 GSTINs filing GSTR-1 and GSTR-3B through Complai
- 10,000+ IRNs generated in production
- Zero P0 security incidents
- >99% filing success rate across all pilot customers
- Platform runs stably at target SLOs

---

## 19. Open Questions & Future Considerations

- **Multi-IRP support:** NIC1 default; add IRIS/Cygnet as failover in Phase 2
- **International expansion:** GCC (Gulf) VAT compliance explored in Year 2+
- **Embedded finance:** direct lending against invoices (vs TReDS routing) in Year 2+
- **Vertical-specific modules:** pharma (CDSCO compliance), logistics (FASTag integration) in Year 2+

---

**End of PRD v1.0. Approved for build.**
