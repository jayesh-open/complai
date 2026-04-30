export type Section2025 = "392" | "393(1)" | "393(2)" | "393(3)";

export type ResidencyStatus = "RESIDENT" | "NON_RESIDENT" | "ANY_PERSON";

export type DeducteeCategory =
  | "INDIVIDUAL"
  | "HUF"
  | "COMPANY"
  | "FIRM"
  | "TRUST"
  | "AOP"
  | "LOCAL_AUTHORITY"
  | "GOVERNMENT";

export type EntryStatus = "PENDING" | "DEPOSITED" | "FILED" | "REVISED";

export type FilingStatus = "NOT_DUE" | "PENDING" | "FILED" | "OVERDUE" | "REVISED";

export interface PaymentCodeInfo {
  code: string;
  label: string;
  section: Section2025;
  subClause: string;
  baseRate: number;
  description: string;
}

export interface Deductee {
  id: string;
  tenantId: string;
  pan: string;
  name: string;
  category: DeducteeCategory;
  residency: ResidencyStatus;
  sectionPreference: Section2025;
  countryCode?: string;
  dtaaApplicable: boolean;
  form41Filed: boolean;
  trcAttached: boolean;
  totalDeductedYTD: number;
  lastTransactionDate: string;
  address?: string;
  email?: string;
  phone?: string;
}

export interface TDSEntry {
  id: string;
  tenantId: string;
  deducteeId: string;
  deducteePan: string;
  deducteeName: string;
  section2025: Section2025;
  paymentCode: string;
  subClause: string;
  taxYear: string;
  quarter: string;
  transactionDate: string;
  paymentDate?: string;
  grossAmount: number;
  baseRate: number;
  cessRate: number;
  surchargeRate: number;
  effectiveRate: number;
  tdsAmount: number;
  surcharge: number;
  cess: number;
  totalTax: number;
  challanNumber?: string;
  challanDate?: string;
  status: EntryStatus;
  noPanDeduction: boolean;
  lowerCertApplied: boolean;
  dtaaCountryCode?: string;
  natureOfPayment: string;
  invoiceNumber?: string;
}

export interface ReturnStatus {
  formType: "138" | "140" | "144";
  formLabel: string;
  taxYear: string;
  quarter: string;
  status: FilingStatus;
  dueDate: string;
  filedDate?: string;
  tokenNumber?: string;
  acknowledgementNumber?: string;
  entryCount: number;
  totalTds: number;
}

export interface TDSImportRow {
  rowNumber: number;
  deducteePan: string;
  deducteeName: string;
  paymentCode: string;
  amount: number;
  tdsAmount: number;
  transactionDate: string;
  status: "valid" | "error" | "warning";
  errors: string[];
}

export interface TDSCalculationResult {
  section: Section2025;
  subClause: string;
  paymentCode: string;
  baseRate: number;
  cessRate: number;
  surchargeRate: number;
  effectiveRate: number;
  tdsAmount: number;
  surcharge: number;
  cess: number;
  totalTax: number;
  thresholdMet: boolean;
  thresholdAmount: number;
  noPanApplied: boolean;
}
