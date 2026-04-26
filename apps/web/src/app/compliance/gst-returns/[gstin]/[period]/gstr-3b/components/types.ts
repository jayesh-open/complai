export type WizardStep = "auto-populate" | "review" | "pay" | "sign" | "file" | "acknowledge";

export interface StepDef {
  id: WizardStep;
  label: string;
  number: number;
}

export const STEPS: StepDef[] = [
  { id: "auto-populate", label: "Auto-Populate", number: 1 },
  { id: "review", label: "Review", number: 2 },
  { id: "pay", label: "Pay", number: 3 },
  { id: "sign", label: "Sign", number: 4 },
  { id: "file", label: "File Return", number: 5 },
  { id: "acknowledge", label: "Acknowledgement", number: 6 },
];

export type DataSource = "gstr1" | "gstr2b" | "computed" | "override";

export interface TaxRow {
  description: string;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  cess: number;
  source: DataSource;
  overrideReason?: string;
}

export interface ITCRow {
  description: string;
  cgst: number;
  sgst: number;
  igst: number;
  cess: number;
  source: DataSource;
}

export interface LedgerBalance {
  cashCgst: number;
  cashSgst: number;
  cashIgst: number;
  creditCgst: number;
  creditSgst: number;
  creditIgst: number;
}

export interface OffsetEntry {
  head: string;
  liability: number;
  creditUsed: number;
  cashUsed: number;
  remaining: number;
}

export interface GSTR3BData {
  gstin: string;
  period: string;
  periodLabel: string;
  outwardSupplies: TaxRow[];
  inwardSupplies: TaxRow[];
  itcAvailed: ITCRow[];
  itcReversed: ITCRow[];
  itcNet: ITCRow[];
  interestLateFee: { interest: number; lateFee: number };
  ledger: LedgerBalance;
  offsets: OffsetEntry[];
  totalLiability: { cgst: number; sgst: number; igst: number };
  totalITC: { cgst: number; sgst: number; igst: number };
  netPayable: { cgst: number; sgst: number; igst: number };
  flags: string[];
}

export function generateMockGSTR3B(gstin: string, period: string): GSTR3BData {
  const periodLabel = `${period.split("-")[1] === "04" ? "April" : period.split("-")[1]} ${period.split("-")[0]}`;

  return {
    gstin,
    period,
    periodLabel: "April 2026",
    outwardSupplies: [
      { description: "Outward taxable supplies (other than zero rated, nil rated and exempted)", taxableValue: 7000000, cgst: 270000, sgst: 270000, igst: 720000, cess: 0, source: "gstr1" },
      { description: "Outward taxable supplies (zero rated)", taxableValue: 2500000, cgst: 0, sgst: 0, igst: 0, cess: 0, source: "gstr1" },
      { description: "Other outward supplies (nil rated, exempted)", taxableValue: 200000, cgst: 0, sgst: 0, igst: 0, cess: 0, source: "gstr1" },
      { description: "Inward supplies (liable to reverse charge)", taxableValue: 350000, cgst: 31500, sgst: 31500, igst: 0, cess: 0, source: "gstr2b" },
      { description: "Non-GST outward supplies", taxableValue: 0, cgst: 0, sgst: 0, igst: 0, cess: 0, source: "computed" },
    ],
    inwardSupplies: [
      { description: "Inter-state supplies from registered persons", taxableValue: 3200000, cgst: 0, sgst: 0, igst: 576000, cess: 0, source: "gstr2b" },
      { description: "Inter-state supplies from unregistered persons", taxableValue: 150000, cgst: 0, sgst: 0, igst: 27000, cess: 0, source: "computed" },
    ],
    itcAvailed: [
      { description: "Import of goods", cgst: 0, sgst: 0, igst: 180000, cess: 0, source: "computed" },
      { description: "Import of services", cgst: 0, sgst: 0, igst: 45000, cess: 0, source: "computed" },
      { description: "Inward supplies liable to reverse charge", cgst: 31500, sgst: 31500, igst: 0, cess: 0, source: "gstr2b" },
      { description: "Inward supplies from ISD", cgst: 12000, sgst: 12000, igst: 8000, cess: 0, source: "computed" },
      { description: "All other ITC", cgst: 185000, sgst: 185000, igst: 420000, cess: 0, source: "gstr2b" },
    ],
    itcReversed: [
      { description: "As per rules 42 & 43 of CGST Rules", cgst: 5000, sgst: 5000, igst: 0, cess: 0, source: "computed" },
      { description: "Others", cgst: 0, sgst: 0, igst: 0, cess: 0, source: "computed" },
    ],
    itcNet: [
      { description: "Net ITC available", cgst: 223500, sgst: 223500, igst: 653000, cess: 0, source: "computed" },
    ],
    interestLateFee: { interest: 0, lateFee: 0 },
    ledger: {
      cashCgst: 150000, cashSgst: 150000, cashIgst: 200000,
      creditCgst: 223500, creditSgst: 223500, creditIgst: 653000,
    },
    offsets: [
      { head: "IGST", liability: 720000, creditUsed: 653000, cashUsed: 67000, remaining: 0 },
      { head: "CGST", liability: 301500, creditUsed: 223500, cashUsed: 78000, remaining: 0 },
      { head: "SGST", liability: 301500, creditUsed: 223500, cashUsed: 78000, remaining: 0 },
    ],
    totalLiability: { cgst: 301500, sgst: 301500, igst: 720000 },
    totalITC: { cgst: 223500, sgst: 223500, igst: 653000 },
    netPayable: { cgst: 78000, sgst: 78000, igst: 67000 },
    flags: ["2 pending IMS actions", "3 RCM invoices included"],
  };
}
