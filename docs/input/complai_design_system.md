# Complai — Design System

**Version:** 1.0
**Status:** Approved for build
**Foundation:** Aura Design System (see `AURA_DESIGN_SYSTEM.md`)

---

## 0. Relationship to Aura

Complai's design system is built on top of Aura — the foundational design system used across our product family. Aura provides the 15 theme variants, type scale, spacing, animation tokens, and primitive components. Complai applies Aura to the specific information architecture, information density, and regulatory affordances of an enterprise compliance product.

**Inherited from Aura:**
- All 15 theme variants (5 families × 3 variants)
- Type scale, spacing system, radius tokens, animation timings
- Component primitives: cards, buttons, inputs, modals, toasts, skeletons
- Selected-state patterns (accent-tinted bg + accent text + accent border)
- Action-guidance hierarchy (pulsing CTA → highlighted card → step indicator → badge nudge → empty-state CTA → tooltip → helper text)
- Formatting helpers (DD/MM/YYYY dates, locale numbers, font-mono for IDs)

**Added and overridden for Complai:**
1. **Default theme: Light Classic** — warm cream + terracotta
2. **Information architecture** — sidebar grouping specific to the 11 compliance modules
3. **Compliance-specific components** — filing wizards, reconciliation grids, government status pills, tax form replicas, audit timelines, vendor scorecards, period selectors
4. **Regulatory affordances** — signed-action confirmations, audit-trail surfacing, data-residency indicators, downstream system status banners
5. **Default density: Compact** — compliance users scan thousands of rows
6. **Multi-tenant affordances** — tenant switcher, GSTIN/PAN scope selector, environment indicator

Where this document and Aura disagree, this document wins for Complai. Where silent, Aura applies.

---

## 1. Theme — Light Classic (Default)

Light Classic is Complai's default theme — warm off-white background, warm cream sidebar, terracotta accent. All other Aura themes (Dark, Ocean, Forest, etc.) remain available for user selection in Settings → Appearance.

### 1.1 Color tokens

```ts
export const lightClassicTheme: ThemeColors = {
  bgPrimary:    "#faf8f5",   // page background, warm off-white
  bgSecondary:  "#ffffff",   // cards, panels, elevated surfaces
  bgTertiary:   "#f5f3f0",   // inputs, table-row hover
  bgSidebar:    "#f0ebe3",   // warm cream sidebar
  bgOverlay:    "rgba(90, 70, 50, 0.25)",

  borderDefault: "#e0d8cc",
  borderLight:   "#eae3d8",
  borderFocus:   "#d97706",

  textPrimary:   "#1a1612",
  textSecondary: "#4a4239",
  textMuted:     "#6b6358",
  textDisabled:  "#c4bdb4",

  accent:        "#d97706",   // terracotta / amber-700
  accentHover:   "#b45309",
  accentMuted:   "rgba(217, 119, 6, 0.08)",
  accentBorder:  "rgba(217, 119, 6, 0.20)",
  accentText:    "#ffffff",

  success: "#059669", successMuted: "rgba(5,150,105,0.08)",   successBorder: "rgba(5,150,105,0.20)",
  danger:  "#dc2626", dangerMuted:  "rgba(220,38,38,0.07)",   dangerBorder:  "rgba(220,38,38,0.18)",
  warning: "#d97706", warningMuted: "rgba(217,119,6,0.08)",   warningBorder: "rgba(217,119,6,0.18)",
  info:    "#2563eb", infoMuted:    "rgba(37,99,235,0.07)",   infoBorder:    "rgba(37,99,235,0.18)",

  purple: "#7c3aed", purpleMuted: "rgba(124,58,237,0.07)",
  pink:   "#db2777", pinkMuted:   "rgba(219,39,119,0.07)",
  teal:   "#0d9488", tealMuted:   "rgba(13,148,136,0.07)",
  orange: "#ea580c", orangeMuted: "rgba(234,88,12,0.07)",

  shadowSm: "0 1px 3px rgba(140,120,100,0.08)",
  shadowMd: "0 4px 12px rgba(140,120,100,0.10)",
  shadowLg: "0 8px 32px rgba(140,120,100,0.12)",
  shadowAccent: "0 4px 16px rgba(217,119,6,0.18)",
};
```

### 1.2 Semantic color mapping

Strict mapping for compliance domain states:

| Domain state | Color token | Examples |
|---|---|---|
| Filed / Reconciled / Approved / Verified / Matched | `success` | "Filed", "Approved" |
| Failed / Rejected / Defective / Mismatch / Quarantined | `danger` | "Rejected", "Mismatch" |
| Pending Approval / On Hold / In Progress / Draft | `warning` | "Pending Approval", "Draft" |
| Processing / Submitted / Under Review | `info` | "Submitted", "Processing" |
| RCM / Composition / SEZ / Export-LUT | `purple` | "RCM", "Export-LUT" |
| Vendor categories / Tags | `teal` | "Strategic", "MSME-Registered" |
| Reverse-charge / Refund / Reversed | `orange` | "Refund", "Reversed" |

### 1.3 Module identity colors

Each top-level module has an identity color used in its icon and module hero. Identity colors are **never** used for action buttons — actions always use the global `accent`.

| Module | Identity color |
|---|---|
| Dashboard | `accent` |
| GST | `accent` |
| E-Invoice | `info` |
| E-Way Bill | `teal` |
| TDS | `purple` |
| ITR | `pink` |
| MaxITC / Vendor | `success` |
| AP Automation | `accent` |
| Invoice Discounting | `orange` |
| Compliance Cloud | `info` |
| Complai One (SMB) | `teal` |
| Reports & Analytics | `purple` |
| Settings & Config | `textMuted` |

---

## 2. Information Architecture

### 2.1 Sidebar — seven groups

```
┌─────────────────────────────────────┐
│ [Logo] Complai                      │
│         COMPLIANCE PLATFORM         │
├─────────────────────────────────────┤
│ Dashboard                           │
│ My Tasks                  [5]       │
│ Inbox                     [12]      │
│                                     │
│ ── PROCUREMENT ─────────────────    │
│   Purchase Requests                 │
│   Quotations                        │
│   Purchase Orders                   │
│   Goods Receipts                    │
│                                     │
│ ── PAYABLES ────────────────────    │
│   Invoices                [40]      │
│   Debit Notes                       │
│   Vendors                           │
│   Payments                          │
│                                     │
│ ── COMPLIANCE ──────────────────    │
│   GST                               │
│   E-Invoicing                       │
│   E-Way Bill                        │
│   TDS                               │
│   GSTR-2A Recon                     │
│   ITR (employer)                    │
│   Secretarial                       │
│                                     │
│ ── INSIGHTS ────────────────────    │
│   Reports & Analytics               │
│   Vendor Evaluation                 │
│   Audit Trail                       │
│   CFO Dashboard                     │
│                                     │
│ ── DOCUMENTS ───────────────────    │
│   Documents                         │
│   Email Inbox                       │
│                                     │
│ ── CONFIGURE ───────────────────    │
│   Settings                          │
│   Users & Roles                     │
│   Approval Workflows                │
│   Exception Policies                │
│   Billing                           │
│   Setup Wizard                      │
├─────────────────────────────────────┤
│ [Avatar] User Name                  │
│         Role · Tenant               │
└─────────────────────────────────────┘
```

### 2.2 Sidebar rules

- **Group labels are uppercase**, 11px, font-weight 600, letter-spacing 0.04em, color `textMuted`
- **Each group is collapsible** (chevron next to label); state persisted in localStorage per user
- **Groups maintain fixed order** — never reorder dynamically based on usage
- **Permissioned visibility** — items the user has no permission for are hidden, not greyed out; empty groups hide their label too
- **Pending count badges** — use Aura's badge nudge pattern: 18px circle, `bg-accent` for normal, `bg-danger` for overdue/blocked, `bg-warning` for SLA-breach risk

### 2.3 Sidebar collapsed mode

- 64px wide; only icons visible
- Group labels hidden; group dividers shown as a thin 1px horizontal line
- Tooltip on hover shows item label + count
- Submenus (rare) open as flyout panels to the right

### 2.4 Top header

Fixed 52px height:

```
┌──────────────────────────────────────────────────────────────────────────┐
│ [Page Title H4]  [breadcrumb if deep]                                     │
│                               [Search ⌘K] [🔔 12] [Wallet ₹1.25L] [Avatar]│
└──────────────────────────────────────────────────────────────────────────┘
```

- **Left:** page title (heading-md) + breadcrumb (caption, chevron-separated)
- **Center:** optional context bar — GSTIN/PAN scope selector for multi-scope users
- **Right:** search (⌘K), notifications, wallet balance (if relevant), user avatar

### 2.5 Tenant context selectors

Three header dropdowns for users whose role spans multiple scopes:

1. **Tenant switcher** — for CAs, multi-tenant admins, internal Complai staff; shows tenant name + brand color chip
2. **PAN selector** — for users whose tenant has multiple PANs
3. **GSTIN selector** — every page that operates on a single GSTIN

### 2.6 Environment indicator (non-prod only)

24px-tall colored bar above the header:

- **Sandbox:** `info` bg + white text — "Sandbox environment — no real filings will be made"
- **Staging:** `warning` bg — "Staging environment — using masked production data"
- **DR / Failover:** `danger` bg — "Operating from DR region — some features may be slower"

Production has no bar.

---

## 3. Page Layout Patterns

### 3.1 Workflow List — the dominant pattern

Used for invoices, vendors, e-Invoices, EWBs, returns, TDS deductees, and every other collection. Structure:

```
┌──────────────────────────────────────────────────────────────────────────┐
│ PAGE HEADER                                                               │
│ Title (heading-md) + info icon + lineage icon                             │
│ Subtitle (caption, text-muted)                                            │
│                          [Grouping ▾] [Export ▾] [+ Primary CTA]          │
├──────────────────────────────────────────────────────────────────────────┤
│ KPI METRIC ROW (always 4 cards)                                           │
│ [Total] [Pending] [Unpaid] [Overdue]                                      │
├──────────────────────────────────────────────────────────────────────────┤
│ TAB BAR (status-based; counts mandatory)                                  │
│ [All 228] [Pending 40] [Approved 171] [Paid 59] [RCM 7] [On Hold 3]       │
├──────────────────────────────────────────────────────────────────────────┤
│ FILTER BAR (collapsible)                                                  │
│ [🔍 Search] [Vendor ▾] [Date Range ▾] [Amount ▾]   [Save view]            │
├──────────────────────────────────────────────────────────────────────────┤
│ DATA TABLE (compact density default)                                      │
├──────────────────────────────────────────────────────────────────────────┤
│ FOOTER: pagination + density toggle + total count                         │
└──────────────────────────────────────────────────────────────────────────┘
```

**Mandatory rules:**

1. Always show KPI cards above the tab bar, even if values are 0
2. Tabs are status taxonomy only — never mix actions into tabs
3. Filters in a single horizontal bar; collapse extras into "More filters" popover if >6
4. "Save view" — named filter combinations become user-personal saved tabs
5. Default density = Compact (40px row height); toggle in footer
6. **Bulk action bar** appears when rows selected: filter bar transforms with 150ms slide-down animation
7. **Empty state inside the table area** when filters return zero — not replacing the page

### 3.2 Filing Wizard

Used for: GSTR-1, GSTR-3B, GSTR-9, TDS returns, ITR, EWB (multi-leg).

```
┌──────────────────────────────────────────────────────────────────────────┐
│ WIZARD HEADER                                                             │
│ Return Type · Period · GSTIN                                              │
│ Status pill + last saved timestamp + autosave indicator                   │
├──────────────────────────────────────────────────────────────────────────┤
│ STEP INDICATOR (horizontal, sticky at top)                                │
│ ✓ Ingest ── ✓ Validate ── ● Review ── ○ Pay ── ○ File ── ○ Acknowledge   │
├──────────────────────────────────────────────────────────────────────────┤
│ STEP CONTENT                                                              │
├──────────────────────────────────────────────────────────────────────────┤
│ STICKY FOOTER (always visible)                                            │
│ [Save Draft] [Discard]              [← Previous]  [Next →] [File Now]    │
└──────────────────────────────────────────────────────────────────────────┘
```

**Mandatory rules:**

1. Steps are linear and verb-named: Ingest, Validate, Review, Pay, File, Acknowledge
2. Autosave every 10s with "Saved 12s ago" caption
3. "File" action triggers irreversible-action modal (§4.6) with red confirm button
4. Status pill reflects: Draft → Validated → Signed → Submitted → Acknowledged → Filed
5. Block navigation on unsaved changes (autosave handles most cases; this catches races)
6. Step gating — can't skip to future steps; past steps clickable as read-only review
7. **Government form replica** — where the user reviews line-by-line data matching a government form, layout mirrors the official form (regulatory trust)

### 3.3 Reconciliation Workspace

Used for GSTR-2B ↔ Purchase Register, IMS reconciliation, EWB ↔ GSTR-1.

```
┌──────────────────────────────────────────────────────────────────────────┐
│ HEADER: "Reconciliation — GSTR-2B vs Purchase Register"                   │
│ Period selector · GSTIN selector · Last run timestamp · [Run Recon]       │
├──────────────────────────────────────────────────────────────────────────┤
│ BUCKET SUMMARY (horizontal bar of clickable status counts)                │
│ ✓ Matched 1,847   ⚠ Mismatch 124   ◐ Partial 38   ○ Missing-2B 19       │
│ ○ Missing-PR 7    ✗ Duplicate 3                                          │
├──────────────────────────────────────────────────────────────────────────┤
│ FILTER + AI SUGGEST                                                       │
│ [🔍 Search] [Reason ▾] [Amount-diff ▾] [✨ AI suggestions]               │
├──────────────────────────────────────────────────────────────────────────┤
│ SPLIT-PANE TABLE                                                          │
│ ┌────────────────────────┬────────────────────────┐                       │
│ │ PURCHASE REGISTER      │ GSTR-2B (GSTN data)    │                       │
│ │ Invoice / Date / ₹     │ Invoice / Date / ₹     │                       │
│ │ Vendor / GSTIN         │ Vendor / GSTIN         │                       │
│ │ [Accept] [Link] [Skip] │                        │                       │
│ └────────────────────────┴────────────────────────┘                       │
├──────────────────────────────────────────────────────────────────────────┤
│ BULK: [Accept all matched] [Send vendor reminders] [Export]               │
└──────────────────────────────────────────────────────────────────────────┘
```

**Mandatory rules:**

1. Bucket counts are clickable filters
2. Our data on the left, government data on the right — never reverse
3. AI suggestion icon (✨) appears next to high-confidence probable matches; click to preview, one-click accept
4. Reason codes first-class — every mismatch has a "Why?" link with explanation popover
5. Bulk operations on entire bucket are supported

### 3.4 Entity Detail Page

Used for vendor, customer, employee, director.

Standard Aura detail page with these compliance tabs as defaults:

```
[Overview] [Invoices] [Compliance] [Communications] [Documents] [Audit Trail]
```

- **Overview** has a Compliance Score panel at top
- **Compliance** shows a 12-month × tax-type heatmap (green=on-time, yellow=late, red=not-filed)
- **Audit Trail** is mandatory — every entity detail page has it

### 3.5 CFO Dashboard

Single-page, no scroll at 1440px:

```
Row 1: KPI cards (4) — Net GST outflow, ITC unblocked YTD, Pending refunds, Vendor risk score
Row 2: Charts (2) — GST outflow trend (12 months), Vendor risk distribution
Row 3: Compliance health (4 mini-cards by module) — GST, TDS, E-Invoice, Secretarial
Row 4: Action items (1 wide card) — top 5 things needing CFO attention
```

Every card has drill-down: click → relevant module with filter pre-applied.

### 3.6 Settings

Two-row tabs: primary (most-used) on top, secondary (less-used) below.

**Primary:** Organization, Departments, GSTIN Management, Chart of Accounts, Dimensions, Appearance, GST Configuration, API Integrations, Notifications, Policies, Bank Accounts, Vendor Onboarding, Procurement Policy.

**Secondary:** Vendor Evaluation, Email Inbox, Messaging, Process Configuration, Security, Vendor Portal.

Setup-complete banner at top (dismissible): "Setup complete — all essential configurations are done. [View Summary] [↺ Restart]".

---

## 4. Compliance-Specific Components

### 4.1 KPI Metric Card

Compliance variant adds regulatory context lines:

```
┌─────────────────────────────────┐
│  [icon 30×30 — module color]   │
│                                 │
│  228                            │
│  TOTAL INVOICES                 │
│  ──────────                     │
│  ₹8.76 Cr                       │
└─────────────────────────────────┘
```

Trend indicators (↑ 12%) in top-right; success color for favorable, danger for unfavorable. Favorability is domain-aware: ↑ Total Invoices = neutral, ↑ Overdue = unfavorable.

### 4.2 Data Table — Compliance Density

Aura's data table is the base. Overrides:

- Default density: Compact (40px row, 8px×14px padding)
- Tabular numerics on every numeric column
- Sticky header; sticky first column on horizontal scroll
- Frozen totals row at bottom for financial tables
- Inline status badges (XS variant inside cells)
- Number cells right-aligned, `font-mono`, tabular-nums
- Date cells: 80px fixed, `DD/MM/YYYY`, em-dash for empty
- Match column at far right: Direct / Probable / No Match pill

#### Row meta-icons (12px, after primary identifier)

- ✉ envelope — received via email
- 🔁 repeat — recurring invoice
- 📎 paperclip — has attachments
- `DUP` pill (orange XS) — duplicate detected
- `2A…` pill (info XS) — auto-matched with 2A
- `NON-PO` pill (default-grey XS) — non-PO invoice

### 4.3 Status Badges

In addition to Aura's base variants:

| Variant | Token | Usage |
|---|---|---|
| `filed` | success | "Filed", "Acknowledged" |
| `pending-approval` | warning | "Pending Approval" |
| `draft` | default-grey | "Draft" |
| `quarantined` | danger (strong border) | "Quarantined" |
| `non-po` | default-grey XS | "NON-PO" |
| `dup` | warning XS | "DUP" |
| `direct` | success XS with tick | "Direct ✓2713" |
| `rcm` | purple | "RCM" |
| `sez` | purple | "SEZ" |
| `export-lut` | purple | "EXPORT-LUT" |
| `composition` | teal | "COMPOSITION" |

Always paired with text — never color-only.

### 4.4 Government Status Pill

Specialized badge showing record state vs. an external government system:

```
┌────────────────────────────────────┐
│  ● GSTN  ·  Submitted              │
└────────────────────────────────────┘
```

- 24px tall, `rounded-[6px]`, padding `2px 10px`
- 6px dot, color = status (success/warning/danger/info)
- System abbreviation in `font-mono` uppercase: `GSTN`, `IRP`, `EWB`, `TRACES`, `MCA`, `OLTAS`
- Middle-dot separator before status text
- Hover tooltip: last status update timestamp, latency from submission, retry count

#### System-status banner (page-level)

When a government system is degraded:

```
┌──────────────────────────────────────────────────────────────────────────┐
│ ⚠ GSTN portal experiencing delays. Filings are being queued and will be  │
│   submitted automatically. Last successful call: 14:32. [View status]    │
└──────────────────────────────────────────────────────────────────────────┘
```

Background `warningMuted`, border-bottom `warningBorder`. Auto-dismissed when system recovers.

### 4.5 Audit Trail Timeline

Mandatory on every entity-detail page:

```
│ ●  Filed                                          12/04/2026, 14:32      │
│ │  by Priya Mehta · ARN AA0904260012345                                  │
│ │                                                                        │
│ ●  Submitted to GSTN                              12/04/2026, 14:30      │
│ │  by Priya Mehta · System call took 1.8s                                │
│ │                                                                        │
│ ●  Reviewed                                       12/04/2026, 11:15      │
│ │  by Rohan Shah · No changes                                            │
│ │                                                                        │
│ ●  Edited (Table 4 ITC)                           11/04/2026, 16:42      │
│ │  by Priya Mehta · 3 line items modified [view diff]                    │
│ ●  Created                                        11/04/2026, 09:00      │
│    by System (auto from Sales Register)                                  │
```

- Vertical line in `borderDefault`; dots colored per state
- Each entry: action verb (heading-sm) + actor + timestamp + context line
- "View diff" links open modal with before/after comparison
- Filterable by user, action, date range
- Exportable to signed PDF

### 4.6 Filing Confirmation Modal (irreversible action)

For every irreversible regulatory action (file return, generate IRN, file ROC form, submit FVU):

```
┌──────────────────────────────────────────────────────────────────────────┐
│ ⚠  File GSTR-3B for April 2026                                            │
├──────────────────────────────────────────────────────────────────────────┤
│ You are about to file:                                                    │
│   Period:       April 2026 (FY 2026-27)                                   │
│   GSTIN:        29AABCA1234A1Z5 (Karnataka)                               │
│   Tax payable:  ₹12,45,678 — paid via challan CIN 0123…                   │
│                                                                           │
│ This action is irreversible. Once filed, you cannot revise this return.   │
│                                                                           │
│ Sign with:     ( ) DSC token (USB)    ( ) EVC OTP                         │
│                                                                           │
│ Type "FILE" to confirm: [           ]                                     │
├──────────────────────────────────────────────────────────────────────────┤
│                                          [Cancel]  [Confirm & File]      │
│                                                    ↑ red, disabled until │
│                                                      type-to-confirm OK  │
└──────────────────────────────────────────────────────────────────────────┘
```

**Rules:**

- Confirm button is `danger` color, not accent — regulatory gravity
- Type-to-confirm required for: filings with tax > ₹10 lakh, ROC filings, ITR filings, deletion of master data
- DSC/EVC selector mandatory
- On confirm, modal transforms to live progress view: "Validating → Signing → Submitting to GSTN → Awaiting ARN → Done" — don't close it

### 4.7 Period Selector

```
┌──────────────────────────────────────────┐
│  Period:  [FY 2026-27 ▾]  [Apr 2026 ▾]   │
└──────────────────────────────────────────┘
```

- Two coupled dropdowns: financial year + month/quarter
- Quick presets: "Current Period", "Previous Period", "Last 3 months", "FYTD"
- Annual returns: collapses to single FY dropdown
- Ad-hoc reports: extend with "Custom range" → date-range picker

### 4.8 Vendor Compliance Score Card

```
┌──────────────────────────────────────────────────────────────────────────┐
│  CWC Services [OPC] Pvt Ltd                                               │
│  GSTIN 29ABCDE1234F1Z5  ·  Karnataka                                      │
│                                                                           │
│  Compliance Score                                                         │
│  ●●●●●○○○○○  62 / 100  ·  Risk: Moderate                                  │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │  Filing regularity:    8/10  ✓                                      │  │
│  │  IRN compliance:       9/10  ✓                                      │  │
│  │  Mismatch rate:        4/10  ⚠ — 12% of invoices have mismatches    │  │
│  │  Payment behaviour:    7/10  ✓                                      │  │
│  │  Document hygiene:     6/10  ⚠ — Missing GST cert                   │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│  Category: B (Occasional delay)  ·  Last reviewed: 03/04/2026             │
└──────────────────────────────────────────────────────────────────────────┘
```

- 10-dot bar; ≥80 success, 60–79 warning, <60 danger
- Category letter from MaxITC (A/B/C/D)
- Click → detailed scorecard page

### 4.9 Tax Rule Banner

For regulatory changes:

```
┌──────────────────────────────────────────────────────────────────────────┐
│ ✨ New: IMS-based ITC auto-population now available for FY 2024-25 GSTR-9 │
│    [Learn what's new]  [Don't show again]                                 │
└──────────────────────────────────────────────────────────────────────────┘
```

- `accentMuted` background, `accentBorder`
- Sparkle icon (not info — info implies neutrality, sparkle implies new feature)
- Dismissible per user
- Linked to "What's New" article

### 4.10 Bulk Operation Progress Tray

Floating bottom-right:

```
┌──────────────────────────────────────────────┐
│ 3 Background Jobs                       [—]  │
├──────────────────────────────────────────────┤
│ Bulk IRN — 1,247 / 5,000   ████░░░░  25%     │
│   Started 12:30 · ETA 4 min          [stop]  │
├──────────────────────────────────────────────┤
│ Form 16 generation — Done ✓                  │
│   234 / 234 employees · [Download zip]       │
├──────────────────────────────────────────────┤
│ GSTR-2A pull (FY 2025-26) — 8 of 12 months   │
└──────────────────────────────────────────────┘
```

- 320px wide, max-height 480px
- Collapsible; count badge when collapsed
- Toast on each job completion
- Persists across page navigation

### 4.11 Maker-Checker Approval Card

```
┌──────────────────────────────────────────────────────────────────────────┐
│ APPROVAL REQUIRED                                                         │
│ Invoice INV-2026-0847 · ₹2,95,000                                         │
│ Submitted by: Priya Mehta (Tax Analyst) · 2 hours ago                     │
│                                                                           │
│ Vendor:  Raj Catering Services                                            │
│ GSTIN:   29ABCDE5678G1Z9                                                  │
│ ITC:     ₹45,000 (eligible)                                               │
│                                                                           │
│ ⚠ Vendor compliance score: 58 (Moderate risk)                             │
│                                                                           │
│ Comments (optional): [                                                  ] │
│                                                                           │
│   [Reject] [Send Back for Edit]              [Approve & Continue]        │
└──────────────────────────────────────────────────────────────────────────┘
```

- Always surfaces maker's name + timestamp
- Always surfaces risk signals (compliance score, deviations, missing data)
- Reject = danger color, Send Back = ghost, Approve = primary accent
- Comment optional on approve, mandatory on reject/send-back
- N-of-M sequential chains shown as vertical step indicator on the right

---

## 5. Interaction Patterns

### 5.1 Keyboard shortcuts

Global:

| Shortcut | Action |
|---|---|
| `⌘K` / `Ctrl+K` | Command palette / search |
| `g d` | Dashboard |
| `g i` | Invoices |
| `g g` | GST |
| `g t` | TDS |
| `n i` | New Invoice |
| `n p` | New Purchase Order |
| `/` | Focus filter search |
| `j` / `k` | Move selection down / up |
| `x` | Toggle selection of current row |
| `⌘A` | Select all visible rows |
| `Esc` | Close modal, clear selection, exit edit |
| `?` | Show shortcut cheat sheet |

`⌘K` hint visible in header search input. `?` hint in footer.

### 5.2 Command palette (⌘K)

Categories:

- **Navigate** — jump to any module/page
- **Create** — quick-create vendor, invoice, etc.
- **Recent** — last 10 entities viewed
- **Filings** — pending filings across modules
- **Reports** — saved + standard reports
- **Help** — docs, contact KAM

Each command: icon + name + shortcut + category. Fuzzy search.

### 5.3 Inline editing in tables

For sales register prep, journal entries:

- Click cell → edit in place
- Tab / Shift+Tab → next / previous
- Enter saves + moves down; Esc cancels
- Modified cells: yellow left border
- "Save changes" sticky footer with change count
- Validation on save, not per-keystroke

### 5.4 Drag-and-drop uploads

- Dashed `borderDefault`, `bgTertiary`, centered icon + label
- On drag-over: border `accentBorder`, bg `accentMuted`, label "Drop to upload"
- On drop: instant upload with per-file progress
- For invoices: OCR confidence shown per field

### 5.5 Excel paste

- "Paste from Excel" button beside upload
- Modal with textarea; tab-delimited paste → structured grid
- Column mapper: source → canonical field
- Save mappings as per-tenant templates

---

## 6. Content & Microcopy

### 6.1 Voice

- Direct and operational
- Use domain language (GSTR-1, ITC, RCM) — don't dumb down; provide info-tooltips for junior users
- Numbers first: "5 returns due this week" before "You have some pending work"
- Action-led button labels: "File Return" not "Submit", "Generate IRN" not "Create"
- Acknowledge the user: "Successfully filed GSTR-3B" — not "Done"

### 6.2 Formatting

- **Dates:** `DD/MM/YYYY` (e.g., 11/04/2026)
- **Datetime:** `DD/MM/YYYY, h:mm AM/PM`
- **Currency:** `₹` prefix + Indian numbering (crore/lakh in summaries, full in tables)
- **Percentages:** 1 decimal max
- **GSTIN:** 15 chars, `font-mono`
- **PAN:** 10 chars, `font-mono`
- **Invoice numbers:** `font-mono`, preserve exact user casing

### 6.3 Empty states

| Context | Title | Description | CTA |
|---|---|---|---|
| No invoices | "No invoices yet" | "Upload your first invoice or set up an ERP connection." | "+ New Invoice" |
| No vendors | "No vendors yet" | "Add vendors manually or invite them via the vendor portal." | "+ Add Vendor" |
| No recon runs | "Run your first reconciliation" | "Match your purchase register with GSTR-2B." | "Run Reconciliation" |
| No saved views | "Save your first view" | "Apply filters and save them for quick access." | (hint, no CTA) |
| Filtered to zero | "No results match your filters" | "Try widening the date range or removing filters." | "Clear filters" |

Icon at 48px, 50% opacity.

### 6.4 Error messages

- Plain English first, technical detail collapsed
- First line: "Couldn't connect to GSTN. Your data is safe — we'll retry automatically."
- Below (collapsed): "Error code: GSP-503-TIMEOUT · Request ID: req_abc123"
- Always actionable: tell the user what to do
- "Invoice INV-001 has an invalid HSN code (123). HSN codes must be 4, 6, or 8 digits. [Fix in row]"
- No jargon / stack traces for end-users; reserve for dev console

---

## 7. Accessibility

- Tabular data uses proper `<thead>`, `<tbody>`, `<th scope="col">`
- Currency `aria-label`: `aria-label="Two lakh ninety-five thousand rupees"` for ₹2,95,000
- Status badges use color + text + SVG icon (✓, ⚠, ✗) for colorblind users
- Filing confirmation modals use `role="alertdialog"` not `role="dialog"` (announces irreversibility)
- Every wizard step completable with keyboard alone — tested in CI
- Focus moves to modal heading on open; returns to trigger on close

---

## 8. Responsive Strategy

Desktop-first. Compliance work happens on large screens.

- **Tablet (≥768px):** full functionality, sidebar collapsed by default; horizontal-scroll tables with sticky first column
- **Mobile (<768px):**
  - Read-only default; approvals work fully via My Tasks
  - List pages collapse to card view
  - Filing wizards, recon workspace, bulk ops: "Best experienced on desktop" with deep-link back
  - Complai One (SMB app) is the optimized mobile surface

Don't compromise desktop density to chase mobile parity.

---

## 9. Layout Grid & Density

### 9.1 Page container

- Max width 1600px at Desktop XL
- Page padding 28px
- Content gutter 16px between sibling cards

### 9.2 Density modes

| Mode | Row height | Card padding | Default for |
|---|---|---|---|
| Compact | 40px | 16px | Default (financial tables, lists) |
| Comfortable | 52px | 22px | Detail pages, config screens |
| Spacious | 64px | 28px | Dashboards, executive views |

Persisted per user, per page-type.

### 9.3 Card grids

- **KPI row:** 4 cols Desktop XL, 2 tablet, 1 mobile
- **Chart row:** 2 cols Desktop, 1 below
- **Detail-page panels:** 12-column grid; common splits 8/4 (main + side), 6/6 (side-by-side)

---

## 10. Iconography

Lucide React. Sizes:

- Sidebar: 16px
- Buttons: 14px (small), 16px (standard)
- Table action icons: 14px
- Empty state hero: 48px
- KPI card icons: 16px inside 30×30 colored square

### 10.1 Domain icon mappings

| Concept | Icon |
|---|---|
| Invoice | `FileText` |
| GST return | `FileSpreadsheet` |
| TDS | `Receipt` |
| E-Invoice / IRN | `FileCheck2` |
| E-Way Bill | `Truck` |
| Vendor | `Building2` |
| Customer | `Users` |
| Reconciliation | `GitCompareArrows` |
| Match (success) | `Check` |
| Mismatch | `AlertTriangle` |
| Approval | `ShieldCheck` |
| Filing / Submit | `Send` |
| ITC | `Wallet` |
| Refund | `RotateCcw` |
| Workflow | `Workflow` |
| Audit trail | `History` |
| Compliance score | `Gauge` |
| Document | `Paperclip` |
| Settings | `Settings` |

Never mix icon libraries.

---

## 11. URL Conventions

Predictable URL structure:

```
/dashboard
/inbox
/tasks

/procurement/purchase-requests
/procurement/purchase-requests/{id}
/procurement/purchase-orders/{id}/edit

/payables/invoices
/payables/invoices/{id}
/payables/invoices/{id}/audit

/payables/vendors/{id}
/payables/vendors/{id}/communications

/compliance/gst
/compliance/gst/{gstin}/{period}/gstr-1
/compliance/gst/{gstin}/{period}/gstr-3b
/compliance/gst/{gstin}/{period}/recon

/compliance/e-invoice
/compliance/e-invoice/generate

/compliance/tds/{tan}/{quarter}/24q
/compliance/itr/employer

/insights/reports/{report-id}
/insights/audit-trail

/configure/settings
/configure/users
```

- Breadcrumb auto-derived from URL
- Deep-linking everywhere: every filtered list, every wizard step
- Browser back works as expected; wizard "Previous" = `router.back()` when safe

---

## 12. Top 12 Rules for Complai

In priority order:

1. **Default theme is Light Classic** — cream + terracotta. All 15 themes user-switchable.
2. **Sidebar groups in fixed order:** Dashboard / My Tasks / Inbox → Procurement → Payables → Compliance → Insights → Documents → Configure.
3. **Workflow List is the dominant template** — KPIs → Tabs → Filters → Table → Footer.
4. **Every regulatory action = Filing Confirmation Modal** with type-to-confirm, DSC/EVC, red confirm button.
5. **Government Status Pills** show system + state + dot color. Everywhere external systems are involved.
6. **Reconciliation workspaces are split-pane** — our data left, government right. Bucket counts always clickable.
7. **Audit Trail mandatory** on every entity-detail page.
8. **Period + GSTIN selectors pervasive** — header on every period/GSTIN-scoped page.
9. **Default density = Compact.** Compliance data is dense.
10. **Status badges = text + icon + color, never color-only.**
11. **DD/MM/YYYY, ₹ prefix, Indian numbering in summaries, full in tables.** `font-mono` for IDs.
12. **One pulsing CTA per screen max** — guide to the next regulatory action.

---

## 13. Build Order

### Phase D-0 (Week 1–2): Theme + primitives
Port Aura's `themes.ts`, `theme-provider.tsx`, Tailwind config, chart-theme verbatim. Build 13 primitives in Storybook. Set Light Classic as default.

### Phase D-1 (Week 3–4): Shell + navigation
App shell with 7-group sidebar, top header, breadcrumbs, environment indicator, command palette (⌘K), tenant/PAN/GSTIN selectors. Page templates: Workflow List, Detail, Dashboard, Wizard, Reconciliation Workspace, Settings.

### Phase D-2 (Week 5–8): Compliance components
KPI Metric Card (compliance), Data Table with density modes, all status badges, Government Status Pill, Period Selector, Audit Trail Timeline, Filing Confirmation Modal, Reconciliation split-pane, Vendor Compliance Score Card, Bulk Operation Tray, Maker-Checker Approval Card.

### Phase D-3 (Week 9–10): First end-to-end UI
Invoices Workflow List page. Settings page with two-row tabs. GSTR-1 filing wizard end-to-end as the reference flow.

After D-3, every other module reuses these patterns.

---

## 14. Tech Stack

| Layer | Technology |
|---|---|
| Framework | Next.js 15, App Router, TypeScript 5.4 strict |
| Styling | Tailwind + shadcn/ui customized |
| State | Zustand (theme, tenant ctx, user prefs) |
| Data fetching | TanStack Query |
| Tables | TanStack Table v8 |
| Forms | React Hook Form + Zod (shared with backend) |
| Charts | Recharts |
| Icons | Lucide React |
| Fonts | System + JetBrains Mono (font-mono) |
| Animation | Framer Motion (sparingly) |
| Date | date-fns + date-fns-tz (IST default) |
| i18n | next-intl (Hindi Phase 2) |
| Realtime | SSE (job tray) + WebSocket (collab recon) |
| Component dev | Storybook 8 |
| Testing | Vitest + Playwright + axe-core |
| Visual regression | Chromatic |

---

**End of Design System v1.0. Approved for build.**

Hand this to Claude Code alongside `AURA_DESIGN_SYSTEM.md`. Aura is foundation; Complai is application to the compliance domain. Where this doc is silent, Aura applies. Where this doc speaks, it wins.
