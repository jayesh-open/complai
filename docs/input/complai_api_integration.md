# Complai — External API Integration Specification

**Version:** 1.0
**Date:** April 2026
**Posture:** Enriched APIs everywhere they exist. Operational simplicity over marginal cost optimization.

**Provider stack:**
- **GST / IRP / EWB:** Adaequare uGSP (Enriched APIs)
- **TDS / ITR / KYC / Tax Payment:** Sandbox.co.in (Quicko)
- **MCA21, Banks, ERPs:** direct integration (no aggregator)

**Companion docs:** `complai_prd.md`, `complai_architecture.md`, `complai_design_system.md`, `complai_build_prompt.md`

This is a **practical implementation spec**: every endpoint listed is a real URL with a real HTTP method that an engineer can copy into a Postman collection and test.

---

## 1. Why Enriched-Only (Adaequare)

Complai commits to **Adaequare's Enriched APIs** as the sole integration mode. This is the correct call for a regulated, multi-tenant SaaS, for these reasons:

1. **No encryption in our code.** Adaequare handles AES-256 + SEK + HMAC + GSTN public-key wrapping. This is the highest-risk code in any GST integration — eliminating it removes an entire class of incidents.
2. **No GSTN auth-token lifecycle.** Adaequare auto-generates and refreshes the per-GSTIN auth_token. We authenticate to Adaequare once (24h bearer), then make plain-JSON calls on behalf of any GSTIN. The 6-hour SEK lifecycle, OTP re-auth flows, and refresh-token race conditions are abstracted.
3. **No schema-drift fighting.** When GSTN changes a field, rotates a public key, or modifies an error code, Adaequare's wrapper updates — our service code does not.
4. **Filing-week reliability.** Auto-retry, fault tolerance, GSTN failure handling, and callbacks are built-in. Our on-call engineer debugs business logic, not crypto.
5. **Faster shipping.** A GSTR-1 filing flow ships in 2 weeks on Enriched vs 6 weeks on pass-through.
6. **Single mental model.** All Adaequare calls have the same shape: `Authorization: Bearer <gsp_token>` + `gstin` + `requestid` + body. No exceptions.

The marginal per-call cost of Enriched is the price of eliminating ~30% of our gateway code and most of our filing-day risk. This trade is paid in money rather than engineering time and incident severity.

---

## 2. Adaequare Authentication (One Layer, Plain JSON)

### 2.1 The complete handshake (one-time per environment)

```http
POST https://gsp.adaequare.com/gsp/authenticate?action=GSP&grant_type=token
Headers:
  gspappid:     <client_id>
  gspappsecret: <client_secret>
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer",
  "expires_in": 86399,
  "scope": "gsp",
  "jti": "26b1d566-45cf-4772-bb62-9bbce061ef21"
}
```

Token lifetime: **24 hours**. Authentication calls are not billed.

### 2.2 The contract for every subsequent call

```http
POST https://gsp.adaequare.com/enriched/<service>/...
Headers:
  Content-Type:  application/json
  Authorization: Bearer <gsp_access_token>
  user_name:     <taxpayer_gstn_username>     (for e-Invoice/EWB)
  username:      <taxpayer_gstn_username>     (variant for some EWB endpoints)
  password:      <taxpayer_gstn_password>
  gstin:         <15-char GSTIN>
  requestid:     <unique alphanumeric per call>
Body: <plain JSON, no encryption>
```

No SEK. No HMAC. No public-key wrap. No auth_token refresh. Adaequare handles all of that.

### 2.3 Service identity in Complai

Our `gstn-gateway`, `irp-gateway`, and `ewb-gateway` services maintain:
- **One Adaequare bearer token per environment** in Redis, refreshed proactively at 23h
- **Tenant-scoped credential vault**: each tenant's GSTIN username + password stored encrypted in AWS Secrets Manager at `complai/tenant/{tenant_id}/{gstin}/adaequare-credentials`
- **`requestid` generation**: UUID v4 per logical operation, stored in the outbox row for idempotency

---

## 3. Adaequare GSP Services & Endpoints

```
Sandbox:    https://gsp.adaequare.com/test/...
Production: https://gsp.adaequare.com/...
```

### 3.1 GSP Authentication

| Purpose | Method | URL |
|---|---|---|
| Get GSP access token | POST | `/gsp/authenticate?action=GSP&grant_type=token` |

### 3.2 GST Common Services

| Purpose | Method | URL |
|---|---|---|
| Get AATO (Annual Aggregate Turnover) per GSTIN per FY | GET | `/test/enriched/services?action=GETAATO&gstin={gstin}&finyr={fy}` |

Request Adaequare's full GST Returns Enriched API doc from `gsp_support@adaequare.com` for the GSTR-1/3B/9/9C/2A/2B/IMS endpoint URLs at onboarding.

### 3.3 e-Invoice APIs (15 endpoints)

| Purpose | Method | URL |
|---|---|---|
| Generate IRN | POST | `/enriched/ei/api/invoice` |
| Cancel IRN (within 24h) | POST | `/enriched/ei/api/invoice/cancel` |
| Get e-Invoice by IRN | GET | `/enriched/ei/api/invoice/irn?irn={irn}` |
| Get IRN by document details | GET | `/enriched/ei/api/invoice/irnbydocdetails?doctype={INV/CRN/DBN}&docnum={no}&docdate={DD/MM/YYYY}` |
| Get Taxpayer details (validate counterparty GSTIN) | GET | `/enriched/ei/api/master/gstin?gstin={gstin}` |
| Sync GSTIN details from Common Portal | GET | `/enriched/ei/api/master/syncgstin?gstin={gstin}` |
| Generate EWB by IRN (one-shot e-Invoice + EWB) | POST | `/enriched/ei/api/ewaybill` |
| Cancel EWB (e-Invoice path) | POST | `/enriched/ei/api/ewayapi` |
| Get EWB details by IRN | GET | `/enriched/ei/api/ewaybill/irn?irn={irn}` |
| Extract QR code data | POST | `/enriched/ei/others/extract/qr` |
| Extract Signed Invoice payload | POST | `/enriched/ei/others/extract/invoice` |
| Generate QR code as image (POST) | POST | `/enriched/ei/others/qr/image` |
| Get QR code image (GET) | GET | `/enriched/ei/others/qr/image` |

**Headers:** every e-Invoice call needs `user_name`, `password`, `gstin`, `requestid`, plus GSP `Authorization` bearer.

**Owned by:** `irp-gateway` service.

### 3.4 e-Way Bill APIs (24 endpoints)

EWB uses the action-parameter pattern: `?action=GENEWAYBILL`, `?action=CANEWB`, etc.

#### Generation
| Purpose | Method | URL |
|---|---|---|
| Generate single EWB | POST | `/enriched/ewb/ewayapi?action=GENEWAYBILL` |
| Generate Consolidated EWB | POST | `/enriched/ewb/ewayapi?action=GENCEWB` |
| Re-generate Consolidated EWB | POST | `/enriched/ewb/ewayapi?action=REGENTRIPSHEET` |

#### Update / management
| Purpose | Method | URL |
|---|---|---|
| Update Part-B / vehicle number | POST | `/enriched/ewb/ewayapi?action=VEHEWB` |
| Update Transporter | POST | `/enriched/ewb/ewayapi?action=UPDATETRANSPORTER` |
| Cancel EWB (within 24h) | POST | `/enriched/ewb/ewayapi?action=CANEWB` |
| Reject EWB (recipient, within 72h) | POST | `/enriched/ewb/ewayapi?action=REJEWB` |
| Extend EWB validity | POST | `/enriched/ewb/ewayapi?action=EXTENDVALIDITY` |

#### Multi-vehicle movement
| Purpose | Method | URL |
|---|---|---|
| Initiate multi-vehicle movement | POST | `/enriched/ewb/ewayapi?action=MULTIVEHMOVINT` |
| Add vehicle to group | POST | `/enriched/ewb/ewayapi?action=MULTIVEHADD` |
| Change vehicle in group | POST | `/enriched/ewb/ewayapi?action=MULTIVEHUPD` |

#### Retrieval
| Purpose | Method | URL |
|---|---|---|
| Get EWB by EWB number | GET | `/enriched/ewb/ewayapi/GetEwayBill?ewbNo={ewbno}` |
| Get EWBs for transporter (by date) | GET | `/enriched/ewb/ewayapi/GetEwayBillsForTransporter?date={DD/MM/YYYY}` |
| Get EWBs for transporter (by GSTIN) | GET | `/enriched/ewb/ewayapi/GetEwayBillsForTransporterByGSTIN?gstin={gstin}&date={DD/MM/YYYY}` |
| Get EWBs generated by other party | GET | `/enriched/ewb/ewayapi/GetEwayBillsofOtherParty?date={DD/MM/YYYY}` |
| Get EWBs by date (own generated) | GET | `/enriched/ewb/ewayapi/GetEwayBillsByDate?date={DD/MM/YYYY}` |
| Get EWBs rejected by others | GET | `/enriched/ewb/ewayapi/GetEwayBillsRejectedByOthers?date={DD/MM/YYYY}` |
| Get EWB by document details | GET | `/enriched/ewb/ewayapi/GetEwayBillByDocNo?...` |
| Get EWBs for transporter by state | GET | `/enriched/ewb/ewayapi/GetEwayBillsForTransporterByState?date={DD/MM/YYYY}&stateCode={code}` |
| Get Trip Sheet | GET | `/enriched/ewb/ewayapi/GetTripSheet?tripSheetNo={no}` |
| Get EWB report by transporter assigned date | GET | `/enriched/ewb/ewayapi/GetEwayBillReportByTransporterAssignedDate?...` |

#### Masters
| Purpose | Method | URL |
|---|---|---|
| Get Transporter details | GET | `/enriched/ewb/master/GetTransporterDetails?trn_no={transin/gstin}` |
| Get GSTIN details (EWB context) | GET | `/enriched/ewb/master/GetGSTINDetails?gstin={gstin}` |
| Get HSN details by code | GET | `/enriched/ewb/master/GetHsnDetailsByHsnCode?hsncode={hsn}` |
| Get error list (decoder) | GET | `/enriched/ewb/master/GetErrorList` |

**Owned by:** `ewb-gateway` service.

### 3.5 GST Return Filing APIs

These follow the `/enriched/gst/...` pattern with action parameters. Confirm exact URLs at Adaequare onboarding.

#### GSTR-1 (Outward Supply)
Save sections (B2B, B2CL, B2CS, CDNR, CDNUR, EXP, AT, ATADJ, NIL, HSN, DOCS) · Get saved · Get summary · Submit (lock) · File (DSC/EVC) · Reset draft · Status/ARN · Amendments (Tables 9/10/11) · GSTR-1A get + action

#### GSTR-3B (Summary + Tax Payment)
Save summary · Get auto-fill (from 1 + 2B + IMS) · Submit · Offset liability (cash + credit ledger) · File (DSC/EVC) · Reset · Status/ARN · Payment status

#### GSTR-2A (dynamic) and GSTR-2B (static)
2A by section (B2B, B2BA, CDN, CDNA, ISD, ISDA, IMPG, IMPGSEZ, TDS, TCS) · 2A invoice-level · 2B summary · 2B invoice-level · Regenerate 2B after IMS actions

#### IMS — Invoice Management System (live since Oct 2024)
Get pending IMS records · Actions (Accept/Reject/Pending) · Bulk actions (up to 5,000/call) · Summary by counterparty · Action history · Trigger 2B regeneration

#### GSTR-9 / 9C (Annual returns)
Get auto-populated 9 · Save edits · Table 8A bulk pull (paginated) · Compute liability · File with DSC · Status · 9C: pre-filled get, save reconciliation entries, file with auditor DSC

#### Other returns
GSTR-4 (composition annual), GSTR-5 (non-resident), GSTR-6 + 6A (ISD), GSTR-7 (TDS deductor under GST), GSTR-8 (TCS by e-com), CMP-08 (composition quarterly), ITC-04 (job-work), GST RFD-01 (refund), GST DRC-03 (voluntary payment)

#### Ledgers & Challans
Cash ledger · Credit (ITC) ledger · Liability ledger · Negative liability ledger · Payment history · Challan create / get / pay / status / history

#### Public APIs (no per-GSTIN session)
- `searchTP` — search taxpayer by GSTIN
- `searchByPAN` — all GSTINs under a PAN
- `searchHSN` — HSN/SAC search
- `getReturnTrack` — filing status of any GSTIN+period
- `getReturnHistory` — public filing history

**Owned by:** `gstn-gateway` service.

---

## 4. Sandbox.co.in — TDS / ITR / KYC / Tax Payment

Single subscription, single auth model, four product families.

### 4.1 Sandbox auth

```http
POST https://api.sandbox.co.in/authenticate
Headers:
  x-api-key:     <your_api_key>
  x-api-secret:  <your_api_secret>
  x-api-version: 1.0.0
```

Response carries `access_token`. Subsequent calls:
```http
Authorization: <access_token>
x-api-key:     <your_api_key>
x-api-version: <api_version>
```

Token lifetime ~24h. Test host: `developer.sandbox.co.in`. Production: `api.sandbox.co.in`.

Many heavy Sandbox operations (FVU generation, e-File, Form 16 bulk, ITR-V fetch) are **job-based**: start → `job_id` → poll → fetch result.

### 4.2 TDS APIs — for `tds-gateway` / `tds-service`

#### Calculator endpoints
| Purpose | Method | URL |
|---|---|---|
| Calculate TDS on salary | POST | `/tds/calculator/salary` |
| Calculate TDS on non-salary (any section) | POST | `/tds/calculator/non-salary` |
| Compliance check u/s 206AB / 206CCA | POST | `/tds/compliance/206ab` |
| TAN search | POST | `/tds/search/tan?tan={tan}` |

#### Business APIs (Sandbox stores state)
| Purpose | Method | URL |
|---|---|---|
| Create / update Deductor | POST/PUT | `/tds/business/deductors` |
| Create / update Deductee | POST/PUT | `/tds/business/deductees` |
| Record Payment to deductee | POST | `/tds/business/payments` |
| Record Challan (post-deposit) | POST | `/tds/business/challans` |
| Salary master records (24Q Annexure-II) | POST | `/tds/business/salary-details` |
| Business analytics (pre-return summary) | GET | `/tds/business/analytics` |

#### Compliance — return generation & filing
| Purpose | Method | URL |
|---|---|---|
| Generate TXT (RPU output) for 24Q/26Q/27Q/27EQ | POST | `/tds/reports/txt` |
| Download CSI from OLTAS (job, OTP-driven) | POST | `/tds/csi/download` |
| Generate FVU + Form 27A (job) | POST | `/tds/compliance/fvu` |
| Fetch FVU job status + zip | GET | `/tds/compliance/fvu/jobs` |
| E-File TDS Return (job) | POST | `/tds/compliance/e-file` |
| E-File job status | GET | `/tds/compliance/e-file/jobs/{job_id}` |
| E-Verify TDS Return (Form 27A PDF) | POST | `/tds/compliance/e-verify` |
| E-Verify job status | GET | `/tds/compliance/e-verify/jobs/{job_id}` |
| List TDS return jobs | GET | `/tds/compliance/jobs?tan=&financial_year=&quarter=&form=` |

#### Form 16 / 16A
| Purpose | Method | URL |
|---|---|---|
| Generate Form 16 (bulk) | POST | `/tds/form16/generate` |
| Generate Form 16A (bulk) | POST | `/tds/form16a/generate` |
| Download certificates | GET | `/tds/form16/jobs/{job_id}/download` |

**Canonical happy-path flow** (orchestrated by `TdsQuarterlyFilingSaga` in Temporal):

```
1. Build deductee + payment + challan masters (Business APIs)
2. POST /tds/reports/txt → receive TXT blob
3. POST /tds/csi/download → triggers OTP to TAN mobile; user enters OTP; receive CSI file
4. POST /tds/compliance/fvu (TXT + CSI) → job_id → poll → download FVU + Form 27A zip
5. POST /tds/compliance/e-file (FVU zip) → job_id → poll → return_receipt_number (RRR)
6. POST /tds/compliance/e-verify (RRR + Form 27A PDF) → job_id → poll → success
7. POST /tds/form16/generate (after Q4 24Q / each Q for 26Q)
```

OTP is the only human-in-the-loop interruption. Handled via a `human_task` event emitted to the UI.

### 4.3 Income Tax (ITR) APIs — for `itd-gateway` / `itr-service`

Sandbox holds the ERI (e-Return Intermediary) licence; we call clean APIs.

| Purpose | Method | URL |
|---|---|---|
| Add a tax-payer by PAN | POST | `/itd/eri/tax-payers` |
| List tax-payers | GET | `/itd/eri/tax-payers` |
| Get pre-fill data | GET | `/itd/eri/tax-payers/{pan}/prefill?assessment_year={ay}` |
| Compute tax (old vs new regime) | POST | `/itd/calculator/tax` |
| Compute capital gains (domestic / foreign / crypto) | POST | `/itd/calculator/capital-gains` |
| Validate ITR JSON (pre-submit) | POST | `/itd/eri/tax-payers/{pan}/itrs/{ay}/validate` |
| Submit ITR | POST | `/itd/eri/tax-payers/{pan}/itrs/{ay}/submit` |
| Get ITR-V (JSON/PDF/XML) | GET | `/itd/eri/tax-payers/{pan}/itrs/{ay}/itr-v?format={json/pdf/xml}` |
| Get ITR status | GET | `/itd/eri/tax-payers/{pan}/itrs/{ay}/status` |
| Initiate e-Verify | POST | `/itd/eri/tax-payers/{pan}/itrs/{ay}/e-verify` |
| E-Verify status | GET | `/itd/eri/tax-payers/{pan}/itrs/{ay}/e-verify/status` |
| Get refund status | GET | `/itd/eri/tax-payers/{pan}/itrs/{ay}/refund` |
| OCR Form 16 | POST | `/itd/ocr/form-16` |
| OCR Form 26AS / AIS | POST | `/itd/ocr/26as` |
| Fetch AIS | GET | `/itd/eri/tax-payers/{pan}/ais?financial_year={fy}` |
| Fetch 26AS | GET | `/itd/eri/tax-payers/{pan}/26as?financial_year={fy}` |

Forms supported: ITR-1 (Sahaj) through ITR-7. The submit/validate endpoints handle form-specific schema via the JSON body.

### 4.4 KYC APIs — for `kyc-gateway` / vendor onboarding / customer KYC / employee onboarding

Unified KYC layer. One integration, all Indian identity/entity verifications.

#### Individual KYC
| Purpose | Method | URL |
|---|---|---|
| PAN verify + tax-status | GET | `/kyc/pan?pan={pan}` |
| PAN-Aadhaar link status | GET | `/kyc/pan/aadhaar-link?pan={pan}` |
| Aadhaar OTP send | POST | `/kyc/aadhaar/okyc/otp` |
| Aadhaar OTP verify | POST | `/kyc/aadhaar/okyc/verify` |
| Aadhaar XML (offline e-KYC) | POST | `/kyc/aadhaar/xml` |
| DigiLocker — initiate session | POST | `/kyc/digilocker/sessions` |
| DigiLocker — fetch document | GET | `/kyc/digilocker/sessions/{session_id}/documents/{type}` |
| Bank account verify (penny-drop) | POST | `/kyc/bank-account` |
| Bank account verify (reverse penny-drop) | POST | `/kyc/bank-account/reverse-penny-drop` |
| UPI VPA verify | GET | `/kyc/upi?vpa={vpa}` |
| Driving license verify | GET | `/kyc/driving-license?dl_number=&dob=` |
| Voter ID verify | GET | `/kyc/voter-id?epic_number=` |
| Passport verify | POST | `/kyc/passport` |

#### Business KYB
| Purpose | Method | URL |
|---|---|---|
| GSTIN verify | GET | `/kyc/gstin?gstin={gstin}` |
| GSTIN by PAN | GET | `/kyc/gstin/by-pan?pan={pan}` |
| TAN verify | GET | `/kyc/tan?tan={tan}` |
| Business PAN verify | GET | `/kyc/pan/business?pan={pan}` |
| MCA Company Master Data (CIN) | GET | `/kyc/mca/companies/{cin}` |
| MCA Director Master Data (DIN) | GET | `/kyc/mca/directors/{din}` |
| MCA LLP Master Data | GET | `/kyc/mca/llps/{llpin}` |
| Udyam (MSME) verify | GET | `/kyc/udyam?udyam_no={no}` |
| Entitylocker — fetch verified business docs | POST | `/kyc/entitylocker/sessions` |

**Routing decision:** for GSTIN validation in vendor-compliance flows, use Sandbox `/kyc/gstin` (cleaner contract, single API key). Adaequare's `searchTP` is the fallback if Sandbox is unavailable.

### 4.5 Tax Payment APIs

| Purpose | Method | URL |
|---|---|---|
| Pay TDS challan (ITNS-281) | POST | `/tax-payment/tds` |
| Pay self-assessment / advance tax (ITNS-280) | POST | `/tax-payment/income-tax` |
| Pay STT | POST | `/tax-payment/stt` |
| Get challan status | GET | `/tax-payment/challans/{cin}` |
| Download paid challan (PDF) | GET | `/tax-payment/challans/{cin}/download` |

---

## 5. Service-to-Provider Mapping

| Complai service | Provider(s) | Purpose |
|---|---|---|
| `gstn-gateway` | Adaequare Enriched (GST Returns) | GSTR-1/3B/9/9C/2A/2B/IMS, ledgers, challan, refunds |
| `irp-gateway` | Adaequare Enriched (e-Invoice) | IRN generation, cancel, lookup, QR, signed JSON |
| `ewb-gateway` | Adaequare Enriched (EWB) | EWB generation, multi-vehicle, transporter, masters |
| `tds-gateway` | Sandbox.co.in (TDS) | 24Q/26Q/27Q/27EQ end-to-end |
| `itd-gateway` | Sandbox.co.in (Income Tax) | ITR-1 to ITR-7, prefill, AIS, 26AS, e-verify |
| `kyc-gateway` | Sandbox.co.in (KYC) | PAN, Aadhaar, GSTIN, bank, MCA, DigiLocker |
| `tax-payment-gateway` | Sandbox.co.in (Tax Payment) | TDS challan, advance tax, self-assessment |
| `mca-gateway` | Direct MCA21 (no aggregator) | ROC form filing (AOC-4, MGT-7, DIR-3 KYC) |
| `bank-gateway` | Direct bank APIs + Razorpay/Cashfree + TReDS | Payments, collections, discounting |
| `erp-gateway` | Direct per-ERP | SAP, Oracle, Tally, Dynamics 365 |

Major architectural simplifications from the two-provider strategy:
- **No TRACES RPA service** — Sandbox abstracts it internally
- **No in-process encryption library** — Adaequare handles it
- **No per-GSTIN SEK lifecycle management in our code** — Adaequare handles it
- **KYC unified** — single Sandbox integration replaces 5 separate vendor integrations

---

## 6. PRD Module → Provider Mapping

### Module 1 — GST Compliance Platform
Adaequare Enriched (GST Returns): GSTR-1, 3B, 9, 9C, 2A, 2B, IMS, ledgers, challan, refunds, public APIs. One provider, one gateway service.

### Module 2 — E-Invoicing
Adaequare Enriched (e-Invoice): all 15 endpoints (§3.3).

### Module 3 — E-Way Bill
Adaequare Enriched (EWB): all 24 endpoints (§3.4).

### Module 4 — MaxITC + Vendor Compliance
- Adaequare Enriched: GSTR-2A/2B/IMS pulls (heavy, scheduled)
- Sandbox KYC: bulk GSTIN validation, PAN-to-GSTIN lookup
- Adaequare public APIs: `searchTP`, `getReturnTrack` as fallback

### Module 5 — TDS / TCS
Sandbox TDS: end-to-end (§4.2). Sandbox Tax Payment for challan deposit. Sandbox KYC for PAN/TAN verification.

### Module 6 — ITR Filing
Sandbox Income Tax: all endpoints (§4.3). Sandbox KYC for PAN/Aadhaar.

### Module 7 — AP Automation
Sandbox KYC: vendor GSTIN/PAN/bank at submission. Adaequare GST: `GSTR-2B` for ITC eligibility. Direct banks: payment execution.

### Module 8 — Invoice Discounting
Sandbox KYC: vendor onboarding. Adaequare: `getReturnTrack` as alt-data. Bank APIs + TReDS: settlement.

### Module 9 — ClearOne (SMB billing)
Adaequare e-Invoice + EWB + GSTR-1. Sandbox KYC for customer GSTIN. Razorpay/Cashfree for collection.

### Module 10 — Vendor Management
Sandbox KYC as single integration point. Adaequare `getReturnTrack` for compliance scoring.

### Module 11 — Compliance Cloud (Secretarial)
Sandbox KYC: MCA reads. Direct MCA21: form filing (no aggregator covers this).

---

## 7. Implementation Standards

### 7.1 Uniform gateway contract

Every gateway service exposes:

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

Internal callers never see Adaequare or Sandbox specifics — they get a normalized response. The gateway translates.

### 7.2 Idempotency

- **Adaequare**: `requestid: <UUID>` on every call. Adaequare de-dupes ~24h.
- **Sandbox**: `transaction_id: <UUID>` on every call. Sandbox de-dupes.
- **Our outbox**: each row has `request_id` with a `(tenant_id, request_id)` UNIQUE constraint. Replay-safe.

### 7.3 Retry policy

| Class | Action | Max attempts |
|---|---|---|
| HTTP 401 (token expired) | Refresh provider bearer, retry | 1 |
| HTTP 429 (rate limited) | Honor `Retry-After`, exp backoff | 5 |
| HTTP 5xx | Exp backoff (1s, 2s, 4s, 8s, 16s) with jitter | 5 |
| Provider transient | Backoff | 5 |
| Provider 4xx validation | Surface to user, no retry | 0 |
| Network timeout | Retry once after 5s | 2 |

Retries live in the gateway. Callers see one final outcome.

### 7.4 Job-based operations

Many Sandbox operations are async (FVU generation, e-File, e-Verify, Form 16 bulk, ITR-V fetch). Pattern:

```
1. Caller invokes gateway action (e.g., "generate FVU")
2. Gateway calls Sandbox → receives job_id
3. Gateway returns 202 Accepted + operation_id
4. Gateway starts a Temporal workflow that polls Sandbox status
5. On completion, gateway emits an SNS event to caller's channel
6. Caller's service updates state and notifies user
```

Polling cadence: every 10s for first minute, 30s for next 5 min, 2 min thereafter. Max wait 24h.

### 7.5 OTP-driven flows

Some flows require human OTP entry (CSI download, Aadhaar OKYC). Pattern:

```
1. User clicks "Download CSI" in UI
2. Backend calls Sandbox `/tds/csi/download/initiate` → triggers OTP, returns reference_id
3. UI shows OTP input modal (2-min window)
4. Backend calls Sandbox `/tds/csi/download/verify` with OTP + reference_id → receives CSI file
5. CSI persisted in document-service; FVU job triggered
```

Reference ID in Redis with 2-min TTL. If OTP expires, user re-initiates.

### 7.6 Master data caching

All KYC results are cached in Redis:

| Entity | TTL | Invalidation |
|---|---|---|
| GSTIN details | 7 days | Manual refresh, vendor edits |
| PAN details | 30 days | Manual refresh |
| TAN details | 30 days | Manual refresh |
| HSN code | 90 days | Manual refresh |
| Bank account verification | 90 days | Bank change |
| MCA company data | 7 days | New filings |

Cache key: `kyc:{type}:{value}`.

### 7.7 Audit logging

Every external call emits OTel span with:
- `tenant_id`, `gstin/tan/pan`, `provider`, `action`, `request_id`, `attempt`, `latency_ms`, `http_status`, `provider_status_code`, `outcome`, `cost_units`

Never log: passwords, OTPs, encrypted payloads, full PAN/Aadhaar (mask to last 4), bank account numbers, DSC content.

Aggregated in Last9. Per-provider dashboards.

### 7.8 Cost tracking

- **Adaequare**: per-call. Tag each call with `cost_unit=1` (or weight) for monthly reconciliation against invoice.
- **Sandbox**: subscription. Track call counts per product family (TDS/ITR/KYC/Tax-Payment) for capacity planning.
- Per-tenant cost accumulation for internal chargeback.

---

## 8. Onboarding Procedures

### 8.1 Adaequare onboarding (one-time, per environment)

1. Sign Adaequare contract; choose tier.
2. Receive sandbox `client_id` + `client_secret`.
3. Whitelist our NAT egress IPs (prod + DR + staging).
4. Test authentication: `POST /gsp/authenticate` → receive bearer.
5. Test e-Invoice generation in sandbox with NIC test GSTINs.
6. Test EWB generation in sandbox.
7. Submit test summary to Adaequare → receive production credentials.
8. Store secrets in AWS Secrets Manager: `complai/adaequare/{env}/client_id`, `complai/adaequare/{env}/client_secret`.
9. Configure billing.

### 8.2 Sandbox.co.in onboarding

1. Self-serve account creation at https://developer.sandbox.co.in.
2. Subscribe to TDS + IT + KYC + Tax Payment products.
3. Receive sandbox `api_key` + `api_secret`.
4. Test in sandbox with dummy data.
5. Move to production tier.
6. Store secrets in AWS Secrets Manager: `complai/sandbox/{env}/api_key`, `complai/sandbox/{env}/api_secret`.
7. Configure webhook endpoints for job completion.

### 8.3 Per-tenant per-GSTIN onboarding (Adaequare GST flows)

1. Tenant logs into GST portal → My Profile → Manage API Access → enable for 30 days.
2. In Complai, tenant enters GSTIN + GST portal username + password.
3. We store in AWS Secrets Manager: `complai/tenant/{tenant_id}/{gstin}/adaequare-credentials`.
4. Test call (GET ledger) via Adaequare → confirm success.
5. Daily liveness check; alert on failure (needs re-enablement).

### 8.4 Per-tenant per-TAN onboarding (Sandbox TDS flows)

1. Tenant enters TAN in Complai.
2. Call `/kyc/tan` to validate.
3. Tenant uploads TRACES login credentials (Sandbox uses for Form 16 Part-A pulls).
4. Test by fetching latest CSI for the TAN.

### 8.5 No tenant-side onboarding for KYC

KYC APIs run on Complai's shared Sandbox subscription. Tenants never see Sandbox credentials.

---

## 9. Security & Compliance

- **Both providers issue SOC 2 / ISO 27001 reports**. Request annually.
- **Data residency**: Adaequare and Sandbox both host in India — meets DPDP Act requirements.
- **PII handling**: Sandbox PII flows through our `kyc-gateway` and is logged with masking. Aadhaar numbers never persisted — only last 4 + verification result.
- **Secret rotation**: AWS Secrets Manager auto-rotation every 90 days for Adaequare client_secret, Sandbox api_secret, and tenant GST passwords (encrypted with tenant-scoped KMS DEK).
- **No secrets in code or logs.** GitLeaks in CI + AWS CloudTrail logging of all Secrets Manager access.

---

## 10. Cost Model

### Adaequare (per-call)

| Workload (per active GSTIN per month) | Calls |
|---|---|
| GSTR-1 prep + file | ~30 |
| GSTR-3B prep + file | ~15 |
| GSTR-2B + IMS pull | ~50 |
| Annual return work (amortized) | ~20 |
| EWB generation (500/mo avg, 2 calls/EWB) | ~1,000 |
| e-Invoice generation (2,000/mo avg) | ~2,000 |
| Vendor compliance pulls (fallback only) | ~500 |
| **Per active GSTIN per month** | **~3,615** |

ASP tier (200,000 calls for ₹50K/mo) covers ~55 active GSTINs. Year-3 at 50K GSTINs needs custom enterprise contract.

**Year-3 budget**: ~₹2.5–4 cr/year Adaequare.

### Sandbox.co.in (subscription)

Tiered subscriptions bundled across TDS + IT + KYC + Tax Payment. Estimate ~₹15–40 lakh/year at our scale. Main driver: KYC volume.

### Total Year-3 compliance external APIs

- Adaequare: ~₹2.5–4 cr/year
- Sandbox: ~₹40 lakh/year
- Direct integrations (MCA, banks): ~₹50 lakh/year
- **Total**: ~₹3.5–5 cr/year (fixed COGS line)

---

## 11. References

| Provider | Resource | URL |
|---|---|---|
| Adaequare | API Details portal | https://ugsp.adaequare.com/api-details |
| Adaequare | IMS APIs | https://ugsp.adaequare.com/ims-api |
| Adaequare | Support | gsp_support@adaequare.com / 040-66033572 |
| Adaequare | Sales | ugspcontactform@adaequare.com / +91 7207773924 |
| Sandbox | Developer portal | https://developer.sandbox.co.in |
| Sandbox | API reference | https://developer.sandbox.co.in/api-reference |
| Sandbox | Status page | https://status.api.sandbox.co.in |
| Sandbox | Postman collection | https://www.postman.com/in-co-sandbox/sandbox-api |
| Sandbox | Changelog | https://developer.sandbox.co.in/changelog |
| GSTN | GSP ecosystem | https://www.gstn.org.in/gsp-ecosystem |
| NIC e-Invoice | Sandbox | https://einv-apisandbox.nic.in |

---

**End of Specification.**
