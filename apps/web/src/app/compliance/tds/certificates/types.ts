export type CertificateStatus = "GENERATED" | "PENDING" | "ISSUED" | "REVOKED";

export type TaxRegime = "OLD" | "NEW";

export interface Form130Row {
  employeeId: string;
  employeeName: string;
  pan: string;
  grossSalary: number;
  totalTds: number;
  regime: TaxRegime;
  status: CertificateStatus;
  taxYear: string;
  departedDate?: string;
}

export interface SalaryBreakup {
  component: string;
  amount: number;
}

export interface Deduction80C {
  section: string;
  description: string;
  declared: number;
  verified: number;
}

export interface ChallanRef {
  bsrCode: string;
  challanSerial: string;
  depositDate: string;
  amount: number;
  quarter: string;
}

export interface Form130Detail {
  employeeId: string;
  employeeName: string;
  pan: string;
  designation: string;
  taxYear: string;
  regime: TaxRegime;
  employerName: string;
  employerTan: string;
  salaryBreakup: SalaryBreakup[];
  totalGross: number;
  standardDeduction: number;
  deductions80C: Deduction80C[];
  totalDeductions: number;
  taxableIncome: number;
  taxComputed: number;
  tdsChallanRefs: ChallanRef[];
  totalTdsDeducted: number;
  status: CertificateStatus;
}

export interface Form131Row {
  deducteeId: string;
  deducteeName: string;
  pan: string;
  sectionCodes: string[];
  totalAmount: number;
  totalTds: number;
  status: CertificateStatus;
  taxYear: string;
  quarter: string;
}

export interface Form131Detail {
  deducteeId: string;
  deducteeName: string;
  pan: string;
  category: string;
  taxYear: string;
  quarter: string;
  quarterLabel: string;
  deductorName: string;
  deductorTan: string;
  sectionBreakdown: SectionBreakdownRow[];
  totalAmount: number;
  totalTds: number;
  totalSurcharge: number;
  totalCess: number;
  totalTax: number;
  challanRefs: ChallanRef[];
  status: CertificateStatus;
}

export interface SectionBreakdownRow {
  paymentCode: string;
  paymentLabel: string;
  amount: number;
  tdsRate: number;
  tdsAmount: number;
  surcharge: number;
  cess: number;
  dateOfPayment: string;
}
