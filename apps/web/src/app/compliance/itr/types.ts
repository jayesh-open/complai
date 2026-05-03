export type ITRForm = "ITR-1" | "ITR-2" | "ITR-3" | "ITR-4" | "ITR-5" | "ITR-6" | "ITR-7";

export type TaxRegime = "NEW" | "OLD";

export type FilingStatus =
  | "NOT_STARTED"
  | "AIS_FETCHED"
  | "FORM_GENERATED"
  | "REVIEW_PENDING"
  | "EMPLOYEE_APPROVED"
  | "FILED"
  | "ACKNOWLEDGED"
  | "DEFECTIVE";

export type BatchStatus = "DRAFT" | "IN_PROGRESS" | "COMPLETED" | "FAILED";

export type IncomeHead =
  | "SALARY"
  | "HOUSE_PROPERTY"
  | "BUSINESS_PROFESSION"
  | "CAPITAL_GAINS"
  | "OTHER_SOURCES";

export interface ITREmployee {
  id: string;
  tenantId: string;
  pan: string;
  name: string;
  email: string;
  designation: string;
  department: string;
  taxYear: string;
  regime: TaxRegime;
  recommendedForm: ITRForm;
  filingStatus: FilingStatus;
  grossIncome: number;
  taxPayable: number;
  taxPaid: number;
  refundDue: number;
  aisReconciled: boolean;
  lastUpdated: string;
  acknowledgementNumber?: string;
  filedDate?: string;
}

export interface BulkBatch {
  id: string;
  tenantId: string;
  name: string;
  taxYear: string;
  status: BatchStatus;
  totalEmployees: number;
  filed: number;
  pending: number;
  failed: number;
  createdAt: string;
  createdBy: string;
  completedAt?: string;
}

export interface IncomeBreakdown {
  head: IncomeHead;
  amount: number;
  aisAmount?: number;
  reconciled: boolean;
}

export type AISMismatchSeverity = "info" | "warn" | "error";
export type AISMismatchCategory =
  | "SALARY"
  | "TDS"
  | "INTEREST"
  | "DIVIDEND"
  | "SECURITIES"
  | "PROPERTY";

export interface IncomeSubItem {
  label: string;
  amount: number;
  section?: string;
}

export interface IncomeHeadDetail {
  head: IncomeHead;
  gross: number;
  deductions: number;
  net: number;
  subItems: IncomeSubItem[];
  visible: boolean;
}

export interface AISMismatch {
  id: string;
  category: AISMismatchCategory;
  field: string;
  itrValue: number;
  aisValue: number;
  severity: AISMismatchSeverity;
  resolved: boolean;
  resolution?: string;
}

export interface TaxSlab {
  from: number;
  to: number | null;
  rate: number;
  tax: number;
}

export interface TaxComputation {
  totalIncome: number;
  standardDeduction: number;
  taxableIncome: number;
  slabs: TaxSlab[];
  slabTax: number;
  surchargeRate: number;
  surchargeAmount: number;
  surchargeThreshold: string;
  healthEducationCess: number;
  grossTax: number;
  rebate87A: number;
  totalLiability: number;
  tdsCredit: number;
  advanceTax: number;
  selfAssessmentTax: number;
  refundOrPayable: number;
}

export interface DeductionItem {
  section: string;
  label: string;
  declared: number;
  limit: number;
}

export interface EmployeeITRDetail {
  employee: ITREmployee;
  formReason: string;
  incomeHeads: IncomeHeadDetail[];
  deductions: DeductionItem[];
  computation: TaxComputation;
  mismatches: AISMismatch[];
  form10IEAFiled?: boolean;
  auditTrail: {
    action: string;
    actor: string;
    timestamp: string;
    detail?: string;
    status?: "success" | "warning" | "info" | "danger" | "default";
  }[];
}
