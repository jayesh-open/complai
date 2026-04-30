import { ALL_DEDUCTEES, ALL_ENTRIES } from "../mock-data";
import type {
  ChallanRow, ChallanDetail, ChallanAllocation,
  ChallanStatus, ReconciliationSummary,
} from "./types";

const BANKS = ["State Bank of India", "HDFC Bank", "ICICI Bank", "Axis Bank", "Punjab National Bank"];
const BRANCHES = ["Fort Branch, Mumbai", "Andheri West, Mumbai", "Nariman Point, Mumbai", "Borivali East, Mumbai", "Bandra Kurla Complex, Mumbai"];
const STATUSES: ChallanStatus[] = ["CLEARED", "CLEARED", "CLEARED", "PENDING", "REJECTED"];

function bsr(i: number): string {
  return `000${(2000 + i).toString()}`.slice(-7);
}

export function generateChallanRows(): ChallanRow[] {
  const seen = new Map<string, { amount: number; date: string; quarter: string }>();
  for (const e of ALL_ENTRIES) {
    if (e.challanNumber && e.challanDate) {
      const existing = seen.get(e.challanNumber);
      if (existing) {
        existing.amount += e.totalTax;
      } else {
        seen.set(e.challanNumber, { amount: e.totalTax, date: e.challanDate, quarter: e.quarter });
      }
    }
  }

  const rows: ChallanRow[] = [];
  let idx = 0;
  for (const [serial, info] of seen) {
    const allocated = Math.round(info.amount * (idx % 3 === 0 ? 1 : idx % 3 === 1 ? 0.7 : 0));
    rows.push({
      challanId: `chl-${String(idx + 1).padStart(4, "0")}`,
      bsrCode: bsr(idx),
      challanSerial: serial,
      depositDate: info.date,
      amount: info.amount,
      allocatedAmount: allocated,
      unallocatedAmount: info.amount - allocated,
      bankName: BANKS[idx % BANKS.length],
      status: STATUSES[idx % STATUSES.length],
      taxYear: "2026-27",
      quarter: info.quarter,
    });
    idx++;
  }
  return rows;
}

export function generateChallanDetail(challanId: string): ChallanDetail | null {
  const rows = generateChallanRows();
  const row = rows.find((r) => r.challanId === challanId);
  if (!row) return null;

  const matchingEntries = ALL_ENTRIES.filter((e) => e.challanNumber === row.challanSerial);
  const allocations: ChallanAllocation[] = matchingEntries.map((e) => {
    const ded = ALL_DEDUCTEES.find((d) => d.id === e.deducteeId);
    return {
      deducteeId: e.deducteeId,
      deducteeName: ded?.name ?? e.deducteeName,
      pan: e.deducteePan,
      section: e.section2025,
      amount: e.totalTax,
      entryId: e.id,
    };
  });

  const idx = rows.indexOf(row);
  return {
    ...row,
    branchName: BRANCHES[idx % BRANCHES.length],
    tan: "MUMA12345B",
    paymentDate: row.depositDate,
    allocations,
  };
}

export function generateReconciliationSummary(): ReconciliationSummary {
  const rows = generateChallanRows();
  return {
    totalChallans: rows.length,
    totalDeposited: rows.reduce((a, r) => a + r.amount, 0),
    totalAllocated: rows.reduce((a, r) => a + r.allocatedAmount, 0),
    totalUnallocated: rows.reduce((a, r) => a + r.unallocatedAmount, 0),
    fullyAllocated: rows.filter((r) => r.unallocatedAmount === 0).length,
    partiallyAllocated: rows.filter((r) => r.allocatedAmount > 0 && r.unallocatedAmount > 0).length,
    unallocated: rows.filter((r) => r.allocatedAmount === 0).length,
  };
}
