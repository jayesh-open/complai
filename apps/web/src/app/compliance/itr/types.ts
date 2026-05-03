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
