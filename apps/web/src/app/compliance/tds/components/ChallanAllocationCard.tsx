"use client";

import { formatINR } from "@complai/ui-components";
import type { ChallanAllocation } from "../challans/types";

interface Props {
  allocations: ChallanAllocation[];
  totalAmount: number;
  allocatedAmount: number;
}

export function ChallanAllocationCard({ allocations, totalAmount, allocatedAmount }: Props) {
  const pct = totalAmount > 0 ? Math.round((allocatedAmount / totalAmount) * 100) : 0;

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <span className="text-body-sm font-medium text-[var(--text-primary)]">Allocation</span>
        <span className="text-body-sm text-[var(--text-muted)]">{pct}% allocated</span>
      </div>
      <div className="h-2 rounded-full bg-[var(--bg-tertiary)] overflow-hidden">
        <div className="h-full rounded-full bg-[var(--accent)]" style={{ width: `${pct}%` }} />
      </div>
      <div className="flex justify-between text-[10px] text-[var(--text-muted)]">
        <span>Allocated: {formatINR(allocatedAmount)}</span>
        <span>Unallocated: {formatINR(totalAmount - allocatedAmount)}</span>
      </div>
      {allocations.length > 0 && (
        <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
          <table className="w-full text-xs">
            <thead>
              <tr className="bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
                <th className="text-left px-3 py-2 font-medium">Deductee</th>
                <th className="text-left px-3 py-2 font-medium">PAN</th>
                <th className="text-left px-3 py-2 font-medium">Section</th>
                <th className="text-right px-3 py-2 font-medium">Amount</th>
              </tr>
            </thead>
            <tbody>
              {allocations.map((a) => (
                <tr key={a.entryId} className="border-t border-[var(--border-default)]">
                  <td className="px-3 py-2 text-[var(--text-primary)]">{a.deducteeName}</td>
                  <td className="px-3 py-2 font-mono text-[var(--text-muted)]">{a.pan}</td>
                  <td className="px-3 py-2 text-[var(--text-muted)]">{a.section}</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-primary)]">{formatINR(a.amount)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
