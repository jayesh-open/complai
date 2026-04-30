import type { TDSEntry, Deductee } from "../types";

export type TDSFormType = "138" | "140" | "144";

export const FORM_LABELS: Record<TDSFormType, string> = {
  "138": "Form 138 (Salary)",
  "140": "Form 140 (Non-Salary Resident)",
  "144": "Form 144 (Non-Resident)",
};

export const FORM_SECTIONS: Record<TDSFormType, string> = {
  "138": "Section 392",
  "140": "Section 393(1) + 393(3)",
  "144": "Section 393(2) + 393(3)",
};

export type FilingWizardStep = "pull" | "validate" | "preview" | "submit" | "acknowledge";

export interface FilingStepDef {
  id: FilingWizardStep;
  label: string;
  number: number;
}

export const FILING_STEPS: FilingStepDef[] = [
  { id: "pull", label: "Pull Entries", number: 1 },
  { id: "validate", label: "Validate", number: 2 },
  { id: "preview", label: "FVU Preview", number: 3 },
  { id: "submit", label: "Submit", number: 4 },
  { id: "acknowledge", label: "Acknowledgement", number: 5 },
];

export type QuarterFilingStatus = "NOT_STARTED" | "DRAFT" | "SUBMITTED" | "FILED" | "REJECTED";

export interface QuarterInfo {
  id: string;
  label: string;
  dateRange: string;
  dueDate: string;
  dueDateObj: Date;
}

export const QUARTERS: QuarterInfo[] = [
  { id: "q1", label: "Q1", dateRange: "Apr – Jun", dueDate: "31/07/2026", dueDateObj: new Date(2026, 6, 31) },
  { id: "q2", label: "Q2", dateRange: "Jul – Sep", dueDate: "31/10/2026", dueDateObj: new Date(2026, 9, 31) },
  { id: "q3", label: "Q3", dateRange: "Oct – Dec", dueDate: "31/01/2027", dueDateObj: new Date(2027, 0, 31) },
  { id: "q4", label: "Q4", dateRange: "Jan – Mar", dueDate: "31/05/2027", dueDateObj: new Date(2027, 4, 31) },
];

export interface FilingGridCell {
  formType: TDSFormType;
  formLabel: string;
  quarter: string;
  status: QuarterFilingStatus;
  dueDate: string;
  entryCount: number;
  totalTds: number;
}

export interface DeducteeSummaryRow {
  deducteeId: string;
  pan: string;
  name: string;
  category: string;
  residency: string;
  entryCount: number;
  grossTotal: number;
  tdsTotal: number;
  surchargeTotal: number;
  cessTotal: number;
  totalTax: number;
  paymentCodes: string[];
  form41Filed?: boolean;
  trcAttached?: boolean;
  countryCode?: string;
}

export interface PaymentCodeDistribution {
  code: string;
  label: string;
  count: number;
  grossTotal: number;
  tdsTotal: number;
}

export interface TDSFilingData {
  formType: TDSFormType;
  formLabel: string;
  taxYear: string;
  quarter: string;
  quarterLabel: string;
  tan: string;
  deducteeSummaries: DeducteeSummaryRow[];
  entries: TDSEntry[];
  paymentCodeDistribution: PaymentCodeDistribution[];
  totalGross: number;
  totalTds: number;
  totalSurcharge: number;
  totalCess: number;
  totalTax: number;
  deducteeCount: number;
  dtaaBlockers: string[];
}
