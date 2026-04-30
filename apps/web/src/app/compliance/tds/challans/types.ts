export type ChallanStatus = "PENDING" | "CLEARED" | "REJECTED";

export interface ChallanRow {
  challanId: string;
  bsrCode: string;
  challanSerial: string;
  depositDate: string;
  amount: number;
  allocatedAmount: number;
  unallocatedAmount: number;
  bankName: string;
  status: ChallanStatus;
  taxYear: string;
  quarter: string;
}

export interface ChallanAllocation {
  deducteeId: string;
  deducteeName: string;
  pan: string;
  section: string;
  amount: number;
  entryId: string;
}

export interface ChallanDetail {
  challanId: string;
  bsrCode: string;
  challanSerial: string;
  depositDate: string;
  amount: number;
  allocatedAmount: number;
  unallocatedAmount: number;
  bankName: string;
  branchName: string;
  status: ChallanStatus;
  taxYear: string;
  quarter: string;
  tan: string;
  paymentDate: string;
  allocations: ChallanAllocation[];
}

export interface ReconciliationSummary {
  totalChallans: number;
  totalDeposited: number;
  totalAllocated: number;
  totalUnallocated: number;
  fullyAllocated: number;
  partiallyAllocated: number;
  unallocated: number;
}
