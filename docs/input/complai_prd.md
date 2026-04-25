# Complai — Product Requirements Document

**Version:** 2.0
**Status:** Approved for build
**Owner:** Founding Team
**Date:** April 2026

---

## 0. Bank Open Ecosystem Architecture

Complai is one of four sibling applications in the **Bank Open** product family:

| App | Domain | Status |
|---|---|---|
| **Apex** | Procure-to-Pay (vendor master, POs, GRNs, AP invoices, AP payments) | In UAT |
| **Aura** | Order-to-Cash + AR (customer master, AR invoices, payment collection, e-Invoice, E-Way Bill, GSTR-2A/2B view, MaxITC view) | Early stage |
| **Bridge** | Contract management | Early stage |
| **Complai** | Compliance (GST filings, TDS, ITR, Secretarial, audit) | This app |

Plus an **external HRMS** for payroll and Form 16.

### Boundary rules

- **Complai does not own** vendor master (Apex), AR invoices (Aura), or contracts (Bridge). It **consumes** data from siblings via gateway services and adds compliance value on top (filings, scoring, reconciliation, audit).
- **Shared compliance modules** — e-Invoice, E-Way Bill, GSTR-2A/2B view, and MaxITC view exist in both Aura (AR operational user) and Complai (compliance officer) with real-time sync. For Phase 1, Complai builds these modules autonomously; cross-app sync is a Phase 2 concern.
- **Data flow:** Bridge → Apex (POs/contracts) → AP invoices in Apex → consumed by Complai for TDS + 2A/2B recon. Aura → AR invoices → consumed by Complai for GSTR-1. HRMS → Form 16 → consumed by Complai for ITR + 24Q.

---

## 1. Product Overview

### 1.1 What Complai is

Complai is the **compliance layer** in the Bank Open product family — an enterprise-grade compliance SaaS platform for Indian businesses. It handles a company's regulatory obligations: GST returns, e-Invoicing, e-Way Bills, ITC reconciliation, TDS/TCS, income tax, and secretarial compliance. It consumes transactional data from sibling apps (Apex for AP, Aura for AR, Bridge for contracts, HRMS for payroll) and adds compliance intelligence on top.

Complai targets mid-market and enterprise Indian businesses (₹100 Cr to ₹10,000 Cr revenue) that today either stitch together 4–8 point solutions or run ad-hoc spreadsheet-based processes with dedicated compliance teams. The platform replaces the compliance stack with one unified product.

### 1.2 Who it serves

Primary customers: mid-to-large Indian enterprises, including listed companies, PE-backed companies, manufacturing, e-commerce, BFSI adjacent businesses (NBFCs, brokers), and rapidly-growing startups crossing ₹50 Cr revenue.

Primary users within these customers:
- **Tax Managers** (filing GST, TDS, ITR)
- **Tax Analysts** (data preparation, reconciliation)
- **CFOs / Finance Controllers** (oversight dashboards)
- **Company Secretaries** (ROC filings, minutes, registers)
- **Auditors** (read access for statutory audits)
- **External CAs** (assisted filings via marketplace)
- **Integration Admins** (configuring data sync from Apex/Aura/Bridge/HRMS)

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

## 2. Product Scope — The 7 Modules

Complai organizes into seven compliance-focused modules plus a horizontal platform layer and an integration layer for consuming data from sibling Bank Open apps. Each module is independently valuable; the platform's power comes from their integration.

### 2.1 Module overview

| # | Module | Primary Users | Why it matters |
|---|---|---|---|
| 1 | **GST Compliance** | Tax Managers, Analysts | Core monthly/annual filings for every business |
| 2 | **E-Invoicing (IRN)** | Tax Analysts | Mandatory for ₹5 Cr+ turnover; syncs with Aura (Phase 2) |
| 3 | **E-Way Bill** | Logistics, Tax Analysts | Mandatory for goods movement ₹50K+; syncs with Aura (Phase 2) |
| 4 | **ITC Reconciliation + MaxITC + Vendor Compliance Scoring** | Tax Managers, CFOs | Protects ₹ lakhs-crores of ITC; vendor master sourced from Apex |
| 5 | **TDS / TCS** | Tax Analysts | Quarterly filings, Form 16 generation |
| 6 | **Income Tax (ITR)** | Employees, CAs, HR | Employer-bulk filing, CA marketplace |
| 7 | **Secretarial (Compliance Cloud)** | Company Secretaries | ROC filings, registers, minutes |

Plus the **Platform Layer** — identity, tenancy, users/roles, master data, documents, notifications, audit, workflow, rules engine — shared across all 7 modules.

Plus the **Integration Layer** — gateway services consuming data from Apex (AP), Aura (AR), Bridge (contracts), and HRMS (payroll). See §8.

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

## 6. Module 4 — ITC Reconciliation + MaxITC + Vendor Compliance Scoring

### 6.1 Purpose

Protect Input Tax Credit (ITC) by ensuring every vendor is compliant and every invoice is reconciled before it impacts GSTR-3B ITC claims. For a typical mid-market enterprise, ITC represents ₹2–20 crore of working capital annually — losing eligibility because of vendor non-compliance is a real and common risk.

**Important:** Complai does not own the vendor master — Apex does. Complai consumes the vendor master from Apex via the `apex-gateway-service` (read-only sync) and adds compliance scoring + MaxITC orchestration on top. No vendor CRUD in Complai.

### 6.2 Features

#### Vendor compliance scoring (on Apex-sourced vendor master)

- 100-point score across 5 dimensions:
  - Filing regularity (30 points) — how consistently the vendor files GSTR-1/3B
  - IRN compliance (20 points) — does the vendor generate IRNs properly
  - Mismatch rate (20 points) — historical reconciliation gap rate
  - Payment behavior (15 points) — pays on time, responds to queries (data from Apex)
  - Document hygiene (15 points) — has GST certificate, bank details, PAN, etc. (data from Apex)
- Category assignment: A (≥80), B (60–79), C (40–59), D (<40)
- 12-month rolling window
- Automatic vendor-level compliance dashboards

#### ITC orchestration

- Scheduled nightly pull of GSTR-2A for every GSTIN
- Monthly 2B snapshot capture
- IMS pull and action management
- Cross-referencing: which invoices have no matching 2B entry? Which invoices have supplier GSTR-1 mismatches?
- Automated chase emails (via LLM copilot) to delinquent vendors
- ITC-at-risk dashboard: ₹ amount, by vendor, by aging
- Will sync ITC view with Aura (Phase 2)

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

## 9. Module 7 — Compliance Cloud (Secretarial)

### 9.1 Purpose

Company Secretary / legal compliance for corporates and LLPs — ROC filings, board/shareholder resolutions, statutory registers, director KYC.

### 9.2 Features

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

### 9.3 Notes

- MCA21 integration is direct (no aggregator covers this fully)
- Supports companies with 1 to 500+ subsidiaries

---

## 10. Platform Layer

Shared services underneath all 7 modules.

### 10.1 Identity & Access

- Keycloak-based identity
- SSO via SAML 2.0 and OIDC (Google Workspace, Azure AD, Okta)
- MFA (TOTP, SMS, email)
- Step-up authentication for filing operations (re-auth within 5-minute window)
- Session management with device tracking

### 10.2 Tenant & Entity Management

- Multi-tenant with data residency tags
- Tenancy tiers: Pooled, Bridge, Silo, On-Premise
- Tenant hierarchy: Tenant → PANs → GSTINs, TANs, CINs
- Per-tenant branding (logo, colors, email sender)
- Per-tenant feature flags

### 10.3 Users & Roles

- Role templates:
  - Tenant Admin
  - Tax Manager
  - Tax Analyst
  - Viewer
  - CA (external)
  - Integration Admin (configures Apex/Aura/Bridge/HRMS sync)
  - Internal Admin (Complai staff)
- Permission model: resource-action pairs (e.g., `gst.gstr1.file`, `vendor-compliance.score.view`)
- Maker-checker workflows (N-of-M approver chains, per-resource configurable)
- Fine-grained permission overrides

### 10.4 Master Data

- Vendors, customers, employees, items, chart of accounts, HSN codes, state codes, pincodes, bank branches
- Tenant-scoped with RLS enforcement
- Change-data-capture for real-time sync
- ERP sync (bidirectional)

### 10.5 Document Management

- Invoice PDFs, signed JSONs, FVU files, challans, Form 16 PDFs, resolutions, minutes
- S3-backed with metadata in Postgres (JSONB for flexible fields)
- Tenant-scoped DEK (Data Encryption Key) envelope encryption via AWS KMS
- Retention policies per document class (invoice: 8 years, signed return: 8 years, TDS: 7 years)
- Virus scanning on upload
- Pre-signed URL access (15-min TTL)
- Document lineage graph

### 10.6 Notifications

- Email (Amazon SES)
- WhatsApp Business (Meta)
- SMS (MSG91)
- In-app (server-sent events + WebSocket)
- Per-user channel preferences + quiet hours
- Daily digest bundling (09:00 IST)
- Template registry with Handlebars
- DPDP consent tracking + unsubscribe

### 10.7 Audit Trail

- Every state-change event across every service captured
- Tamper-evident (hourly Merkle-chain hashing)
- Search by tenant, user, resource, action, date
- Export to PDF (signed) for regulatory submission
- 30-day hot retention (OpenSearch); long-term archive to S3 with Iceberg metadata

### 10.8 Workflow Orchestration

- Temporal Cloud for all long-running workflows (filing sagas, reconciliation, bulk IRN, Form 16 generation, vendor onboarding)
- Human-task integration (approval, OTP entry, DSC selection)
- Retry and compensation policies per workflow

### 10.9 Rules Engine

- Rules as versioned JSON per tenant
- Categories: tax determination (CGST/SGST/IGST split), HSN applicability, TDS applicability, validation rules, eligibility rules (ITC, RCM, composition)
- Seeded with FY 2026-27 Indian tax rules
- Hot-reload (rule changes propagate without service restart)

### 10.10 AI & Intelligence Layer (Phase 12+)

- ML-based invoice matching (CatBoost) for reconciliation Stage 3
- LLM copilot (Azure OpenAI or Bedrock) for vendor communication drafting, natural-language dashboard queries, return-prep explanation
- PII redaction before LLM calls
- Per-tenant usage caps

---

## 11. Integration Layer

Complai consumes data from the Bank Open sibling applications via dedicated gateway services. Each gateway normalizes sibling-app data into canonical schemas before it enters Complai's core.

### Data flow diagram

```
┌──────────┐    POs, contracts    ┌──────────┐
│  Bridge   │ ──────────────────▶ │   Apex   │
│(contracts)│                     │  (P2P)   │
└──────────┘                     └────┬─────┘
                                      │ AP invoices, vendor master,
                                      │ payments, POs, GRNs
                                      ▼
                              ┌───────────────┐
                              │    Complai     │
                              │ (compliance)   │
                              └───────┬───────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                   │
        ┌───────────▼──┐   ┌────────▼────────┐  ┌──────▼──────┐
        │  TDS + 2A/2B │   │  GSTR-1 (from   │  │ ITR + 24Q   │
        │  recon (from  │   │  Aura AR inv.)  │  │ (from HRMS  │
        │  Apex AP inv.)│   │                 │  │  Form 16)   │
        └──────────────┘   └─────────────────┘  └─────────────┘

┌──────────┐    AR invoices,      ┌───────────────┐
│   Aura   │ ──────────────────▶ │    Complai     │
│  (O2C)   │    customer master   │ (compliance)   │
└──────────┘                     └───────────────┘
                                      │
                              filed-IRN-status, EWB status
                              published back to Aura (Phase 2)

┌──────────┐    Contracts         ┌───────────────┐
│  Bridge   │ ──────────────────▶ │    Complai     │
│(contracts)│                     │ (compliance)   │
└──────────┘                     └───────────────┘
        used for: TDS section determination, secretarial obligations

┌──────────┐    Payroll,          ┌───────────────┐
│   HRMS   │    Form 16 ────────▶ │    Complai     │
│(external)│                     │ (compliance)   │
└──────────┘                     └───────────────┘
        used for: ITR filing, 24Q salary TDS
```

### Gateway services

| Gateway | Source App | Consumes | Publishes back (Phase 2) |
|---|---|---|---|
| `apex-gateway-service` | Apex | Vendor master, AP invoices, payments, POs, GRNs | — |
| `aura-gateway-service` | Aura | Customer master, AR invoices | Filed-IRN-status, EWB status |
| `bridge-gateway-service` | Bridge | Contracts | — |
| `hrms-gateway-service` | HRMS | Payroll data, Form 16 | — |

### Canonical schemas

- **Canonical Invoice Schema** (existing) — covers both AP and AR invoices
- **Canonical Payment Schema** (new) — AP payments from Apex, TDS challans
- **Canonical Contract Schema** (new) — contract clauses for TDS section determination
- **Canonical Payroll Schema** (new) — salary components for 24Q and Form 16

For Phase 1, gateway services use mock data sources (since sibling apps don't expose APIs yet). Real integration switches on in Part 13.

---

## 12. Non-Functional Requirements

### 12.1 Scale targets

| Year | Enterprises | GSTINs | Peak IRN/min | Peak filings/min | Peak recon ops/sec |
|---|---|---|---|---|---|
| Year 1 | 500 | 5,000 | 2,000 | 500 | 1,000 |
| Year 2 | 2,000 | 20,000 | 6,000 | 1,500 | 3,000 |
| Year 3 | 5,000 | 50,000 | 12,000 | 3,000 | 6,000 |

### 12.2 Service Level Objectives

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

### 12.3 Security & Compliance

- **DPDP Act** (Digital Personal Data Protection Act, India) compliant from day one
- **SOC 2 Type II** within 12 months of production
- **ISO 27001** within 18 months
- **Row-level tenant isolation** (Postgres RLS)
- **Tenant-scoped encryption** (per-tenant DEKs via AWS KMS)
- **Data residency:** all data stored in India (`ap-south-1` primary, `ap-south-2` DR)
- **No PII in logs, traces, or LLM prompts** without explicit redaction
- **DSC and EVC** support for regulatory signing
- **Audit-trail immutability** via Merkle-chain hashing

### 12.4 Accessibility

- WCAG 2.2 Level AA compliance
- Keyboard-only completion for every critical workflow (login, filing, reconciliation, approvals)
- Screen reader support for financial tables with proper ARIA labels
- Color-plus-text for every status indicator (never color-only)

### 12.5 Localization

- **Day 1:** English, Hindi
- **Phase 2:** Tamil, Marathi, Gujarati, Kannada, Telugu, Bengali
- DD/MM/YYYY date format, ₹ currency prefix, Indian numbering (crore/lakh)

### 12.6 Performance

- Page load (authenticated pages): <2s P95
- Dashboard render: <3s P95 including data fetch
- Data tables: render 10,000 rows without virtualization lag
- Filing flow: complete in <5 min end-to-end during peak

---

## 13. Tech Stack Summary

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

## 14. Go-to-Market & Build Phasing

### Phase 1 (Weeks 1–8): Core MVP
Modules: Platform layer + GST (1) + e-Invoicing (2) + EWB (3) + ITC Recon + Vendor Compliance Scoring (4) — limited to core reconciliation; vendor master stubbed from Apex

### Phase 2 (Weeks 9–14): TDS + ITR + Secretarial
Modules: TDS (5) + ITR (6) + Compliance Cloud (7)

### Phase 3 (Weeks 15–20): Integration gateways + AI layer
Apex/Aura/Bridge/HRMS gateway services (mock → real); AI intelligence; MaxITC orchestration

### Phase 4 (Weeks 21–24): Production hardening + Bank Open sync
Real sibling app integration; e-Invoice/EWB/2A-2B/MaxITC sync with Aura; observability; security; DR; go-live readiness

---

## 15. Success Criteria — End of Phase 1

- 5 pilot enterprise customers onboarded
- 50 GSTINs filing GSTR-1 and GSTR-3B through Complai
- 10,000+ IRNs generated in production
- Zero P0 security incidents
- >99% filing success rate across all pilot customers
- Platform runs stably at target SLOs

---

## 16. Open Questions & Future Considerations

- **Multi-IRP support:** NIC1 default; add IRIS/Cygnet as failover in Phase 2
- **International expansion:** GCC (Gulf) VAT compliance explored in Year 2+
- **Embedded finance:** direct lending against invoices (vs TReDS routing) in Year 2+
- **Vertical-specific modules:** pharma (CDSCO compliance), logistics (FASTag integration) in Year 2+

---

**End of PRD v2.0. Approved for build. Reflects Bank Open ecosystem scope correction (April 2026).**
