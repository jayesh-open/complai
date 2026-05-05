import type {
  AnnualReturnEntry,
  GSTR9Data,
  GSTR9Table,
  HSNRow,
  LateITCEntry,
  FeesDemandsRow,
  GSTR9CData,
  GSTR9CMismatch,
  AuditedFinancials,
} from "./types";

export const ANNUAL_ENTRIES: AnnualReturnEntry[] = [
  { id: "ar-1", fy: "2025-26", gstin: "29AABCA1234A1Z5", legalName: "Acme Technologies Pvt. Ltd.", turnover: 85000000, gstr9Status: "IN_REVIEW", gstr9cStatus: "NOT_STARTED", gstr9cRequired: true },
  { id: "ar-2", fy: "2025-26", gstin: "27BBBCB5678B2Y6", legalName: "Globex Mfg. Ltd.", turnover: 32000000, gstr9Status: "FILED", gstr9cStatus: "NOT_APPLICABLE", gstr9cRequired: false, filedDate: "2026-12-15", arn: "ARN-9-20260001" },
  { id: "ar-3", fy: "2025-26", gstin: "33CCCDC9012C3X7", legalName: "Initech Solutions", turnover: 120000000, gstr9Status: "NOT_STARTED", gstr9cStatus: "NOT_STARTED", gstr9cRequired: true },
  { id: "ar-4", fy: "2024-25", gstin: "29AABCA1234A1Z5", legalName: "Acme Technologies Pvt. Ltd.", turnover: 72000000, gstr9Status: "ACKNOWLEDGED", gstr9cStatus: "FILED", gstr9cRequired: true, filedDate: "2025-12-28", arn: "ARN-9-20250042" },
];

function buildTables(turnover: number): GSTR9Table[] {
  const t = turnover;
  return [
    {
      tableNumber: 4, title: "Details of advances, inward and outward supplies on which tax is payable",
      rows: [
        { serial: "4A", description: "Supplies made to un-registered persons (B2C)", taxableValue: Math.round(t * 0.3), cgst: Math.round(t * 0.027), sgst: Math.round(t * 0.027), igst: 0, cess: 0, sourceReturn: "GSTR-1" },
        { serial: "4B", description: "Supplies made to registered persons (B2B)", taxableValue: Math.round(t * 0.5), cgst: Math.round(t * 0.03), sgst: Math.round(t * 0.03), igst: Math.round(t * 0.03), cess: 0, sourceReturn: "GSTR-1" },
        { serial: "4C", description: "Zero rated supply (Export) on payment of tax", taxableValue: Math.round(t * 0.08), cgst: 0, sgst: 0, igst: Math.round(t * 0.014), cess: 0, sourceReturn: "GSTR-1" },
      ],
    },
    {
      tableNumber: 5, title: "Details of outward supplies on which tax is NOT payable",
      rows: [
        { serial: "5A", description: "Zero rated supply (Export) without payment of tax", taxableValue: Math.round(t * 0.05), cgst: 0, sgst: 0, igst: 0, cess: 0, sourceReturn: "GSTR-1" },
        { serial: "5B", description: "Supply to SEZs without payment of tax", taxableValue: Math.round(t * 0.02), cgst: 0, sgst: 0, igst: 0, cess: 0, sourceReturn: "GSTR-1" },
        { serial: "5C", description: "Exempted", taxableValue: Math.round(t * 0.03), cgst: 0, sgst: 0, igst: 0, cess: 0, sourceReturn: "GSTR-1" },
      ],
    },
    {
      tableNumber: 6, title: "Details of ITC availed during the financial year",
      rows: [
        { serial: "6A", description: "Total amount of ITC availed through FORM GSTR-3B", taxableValue: 0, cgst: Math.round(t * 0.025), sgst: Math.round(t * 0.025), igst: Math.round(t * 0.04), cess: 0, sourceReturn: "GSTR-3B" },
        { serial: "6B", description: "Inward supplies (other than imports and inward from SEZ)", taxableValue: 0, cgst: Math.round(t * 0.02), sgst: Math.round(t * 0.02), igst: Math.round(t * 0.015), cess: 0, sourceReturn: "GSTR-2B" },
        { serial: "6E", description: "Import of goods", taxableValue: 0, cgst: 0, sgst: 0, igst: Math.round(t * 0.015), cess: 0, sourceReturn: "BoE" },
      ],
    },
    {
      tableNumber: 7, title: "Details of ITC reversed and ineligible ITC",
      rows: [
        { serial: "7A", description: "As per Rule 37", taxableValue: 0, cgst: Math.round(t * 0.001), sgst: Math.round(t * 0.001), igst: 0, cess: 0 },
        { serial: "7B", description: "As per Rule 39", taxableValue: 0, cgst: Math.round(t * 0.0005), sgst: Math.round(t * 0.0005), igst: 0, cess: 0 },
        { serial: "7H", description: "Other reversals", taxableValue: 0, cgst: Math.round(t * 0.0008), sgst: Math.round(t * 0.0008), igst: Math.round(t * 0.001), cess: 0 },
      ],
    },
    {
      tableNumber: 8, title: "Other ITC related information",
      rows: [
        { serial: "8A", description: "ITC as per GSTR-2B (Table 3 & 5)", taxableValue: 0, cgst: Math.round(t * 0.022), sgst: Math.round(t * 0.022), igst: Math.round(t * 0.035), cess: 0, sourceReturn: "GSTR-2B" },
      ],
    },
    {
      tableNumber: 9, title: "Details of tax paid as declared in returns filed during the FY",
      rows: [
        { serial: "9", description: "Tax paid through cash ledger", taxableValue: 0, cgst: Math.round(t * 0.008), sgst: Math.round(t * 0.008), igst: Math.round(t * 0.012), cess: 0, sourceReturn: "GSTR-3B" },
      ],
    },
  ];
}

function buildHSN(turnover: number): HSNRow[] {
  return [
    { hsn: "8471", description: "Computers and peripherals", uqc: "NOS", quantity: 2500, taxableValue: Math.round(turnover * 0.25), cgst: Math.round(turnover * 0.0225), sgst: Math.round(turnover * 0.0225), igst: Math.round(turnover * 0.01), digitTier: 4 },
    { hsn: "998314", description: "IT consulting services", uqc: "NOS", quantity: 1200, taxableValue: Math.round(turnover * 0.35), cgst: Math.round(turnover * 0.0315), sgst: Math.round(turnover * 0.0315), igst: Math.round(turnover * 0.02), digitTier: 6 },
    { hsn: "85176290", description: "Network equipment", uqc: "NOS", quantity: 800, taxableValue: Math.round(turnover * 0.15), cgst: Math.round(turnover * 0.0135), sgst: Math.round(turnover * 0.0135), igst: Math.round(turnover * 0.008), digitTier: 8 },
    { hsn: "9983", description: "Other professional services", uqc: "NOS", quantity: 450, taxableValue: Math.round(turnover * 0.12), cgst: Math.round(turnover * 0.0108), sgst: Math.round(turnover * 0.0108), igst: Math.round(turnover * 0.005), digitTier: 4 },
  ];
}

function buildLateITC(turnover: number): LateITCEntry[] {
  return [
    { table: "6H", description: "ITC reclaimed on reversal of Rule 37 (payment within 180 days)", amount: Math.round(turnover * 0.002), period: "Q2 2025-26", rule: "Rule 37 proviso — ITC reversed for non-payment can be reclaimed when payment is made within the extended period" },
    { table: "8C", description: "Difference between ITC claimed in GSTR-3B vs GSTR-2B (gap rectification)", amount: Math.round(turnover * 0.001), period: "Sep 2025", rule: "Section 16(4) — ITC missed in monthly returns can be claimed in annual return up to Sep 30 of next FY" },
    { table: "13", description: "ITC declared in current FY relating to previous FY invoices", amount: Math.round(turnover * 0.0015), period: "FY 2024-25 invoices", rule: "Section 16(4) read with GSTR-9 Table 13 — ITC on prior-year invoices claimed in current year annual" },
  ];
}

function buildFeesAndDemands(turnover: number): FeesDemandsRow[] {
  return [
    { table: 14, description: "Differential tax paid on account of GSTR-9 declarations", amount: Math.round(turnover * 0.001), category: "demand" },
    { table: 15, description: "Refund claimed from electronic cash ledger", amount: Math.round(turnover * 0.0005), category: "refund" },
    { table: 16, description: "Demands and refunds — DRC-07 demand", amount: Math.round(turnover * 0.002), category: "demand" },
    { table: 17, description: "Late fees payable and paid", amount: 10000, category: "late_fee" },
  ];
}

export function generateGSTR9Data(gstin: string, fy: string): GSTR9Data {
  const entry = ANNUAL_ENTRIES.find((e) => e.gstin === gstin && e.fy === fy);
  const turnover = entry?.turnover ?? 50000000;
  const legalName = entry?.legalName ?? "Unknown Entity";

  return {
    gstin,
    fy,
    legalName,
    turnover,
    thresholdExceeded: turnover > 20000000,
    tables: buildTables(turnover),
    hsnRows: buildHSN(turnover),
    lateITC: buildLateITC(turnover),
    feesAndDemands: buildFeesAndDemands(turnover),
  };
}

function buildMismatches(turnover: number): GSTR9CMismatch[] {
  return [
    { id: "m-1", section: "II", category: "turnover", description: "Aggregate turnover mismatch (books vs GSTR-9 Part II)", booksAmount: turnover * 1.02, gstr9Amount: turnover, difference: turnover * 0.02, severity: "ERROR", resolved: false },
    { id: "m-2", section: "II", category: "turnover", description: "Unbilled revenue adjustment not reflected in GSTR-9", booksAmount: Math.round(turnover * 0.008), gstr9Amount: 0, difference: Math.round(turnover * 0.008), severity: "WARN", resolved: false },
    { id: "m-3", section: "III", category: "tax", description: "CGST payable — rate-wise variance at 18% slab", booksAmount: Math.round(turnover * 0.032), gstr9Amount: Math.round(turnover * 0.03), difference: Math.round(turnover * 0.002), severity: "WARN", resolved: false },
    { id: "m-4", section: "III", category: "tax", description: "IGST payable — export supply tax under-reported", booksAmount: Math.round(turnover * 0.015), gstr9Amount: Math.round(turnover * 0.014), difference: Math.round(turnover * 0.001), severity: "INFO", resolved: false },
    { id: "m-5", section: "IV", category: "itc", description: "ITC claimed (CGST) — excess ITC in books vs GSTR-9", booksAmount: Math.round(turnover * 0.026), gstr9Amount: Math.round(turnover * 0.025), difference: Math.round(turnover * 0.001), severity: "INFO", resolved: false },
    { id: "m-6", section: "IV", category: "itc", description: "ITC claimed (IGST) — import duty credit discrepancy", booksAmount: Math.round(turnover * 0.018), gstr9Amount: Math.round(turnover * 0.015), difference: Math.round(turnover * 0.003), severity: "ERROR", resolved: false },
  ];
}

export function generateGSTR9CData(gstin: string, fy: string): GSTR9CData {
  const entry = ANNUAL_ENTRIES.find((e) => e.gstin === gstin && e.fy === fy);
  const turnover = entry?.turnover ?? 50000000;
  const legalName = entry?.legalName ?? "Unknown Entity";

  const audited: AuditedFinancials = {
    grossTurnover: Math.round(turnover * 1.02),
    taxableTurnover: Math.round(turnover * 0.88 * 1.02),
    taxPayable: {
      cgst: Math.round(turnover * 0.032),
      sgst: Math.round(turnover * 0.032),
      igst: Math.round(turnover * 0.015),
      cess: 0,
    },
    itcClaimed: {
      cgst: Math.round(turnover * 0.026),
      sgst: Math.round(turnover * 0.026),
      igst: Math.round(turnover * 0.018),
      cess: 0,
    },
  };

  return {
    gstin,
    fy,
    legalName,
    gstr9Turnover: turnover,
    gstr9cRequired: turnover > 50000000,
    audited,
    mismatches: buildMismatches(turnover),
  };
}
