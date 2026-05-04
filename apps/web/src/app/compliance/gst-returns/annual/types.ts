export type GSTR9Status = "NOT_STARTED" | "AGGREGATING" | "IN_REVIEW" | "FILED" | "ACKNOWLEDGED";
export type GSTR9CStatus = "NOT_APPLICABLE" | "NOT_STARTED" | "RECONCILING" | "CERTIFIED" | "FILED";

export type GSTR9Step =
  | "threshold"
  | "review-tables"
  | "late-itc"
  | "hsn-summary"
  | "fees-demands"
  | "submit"
  | "acknowledge";

export interface GSTR9StepDef {
  id: GSTR9Step;
  label: string;
  number: number;
}

export const GSTR9_STEPS: GSTR9StepDef[] = [
  { id: "threshold", label: "Threshold", number: 1 },
  { id: "review-tables", label: "Review Tables", number: 2 },
  { id: "late-itc", label: "Late ITC", number: 3 },
  { id: "hsn-summary", label: "HSN Summary", number: 4 },
  { id: "fees-demands", label: "Fees & Demands", number: 5 },
  { id: "submit", label: "Submit", number: 6 },
  { id: "acknowledge", label: "Acknowledgement", number: 7 },
];

export interface AnnualReturnEntry {
  id: string;
  fy: string;
  gstin: string;
  legalName: string;
  turnover: number;
  gstr9Status: GSTR9Status;
  gstr9cStatus: GSTR9CStatus;
  gstr9cRequired: boolean;
  filedDate?: string;
  arn?: string;
}

export interface GSTR9TableRow {
  serial: string;
  description: string;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  cess: number;
  sourceReturn?: string;
}

export interface GSTR9Table {
  tableNumber: number;
  title: string;
  rows: GSTR9TableRow[];
}

export interface HSNRow {
  hsn: string;
  description: string;
  uqc: string;
  quantity: number;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  digitTier: 4 | 6 | 8;
}

export interface LateITCEntry {
  table: "6H" | "8C" | "13";
  description: string;
  amount: number;
  period: string;
  rule: string;
}

export interface FeesDemandsRow {
  table: number;
  description: string;
  amount: number;
  category: "late_fee" | "demand" | "refund";
}

export interface GSTR9Data {
  gstin: string;
  fy: string;
  legalName: string;
  turnover: number;
  thresholdExceeded: boolean;
  tables: GSTR9Table[];
  hsnRows: HSNRow[];
  lateITC: LateITCEntry[];
  feesAndDemands: FeesDemandsRow[];
}
