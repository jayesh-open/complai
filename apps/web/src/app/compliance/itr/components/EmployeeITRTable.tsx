"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";
import { ITRFormRecommendationPill } from "./ITRFormRecommendationPill";
import { RegimeIndicator } from "./RegimeIndicator";
import { FilingStatusPill } from "./FilingStatusPill";
import type { ITREmployee } from "../types";

interface EmployeeITRTableProps {
  employees: ITREmployee[];
  className?: string;
}

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

export function EmployeeITRTable({ employees, className }: EmployeeITRTableProps) {
  return (
    <div className={cn("bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden", className)}>
      <table className="w-full">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Employee</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">PAN</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Form</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Regime</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">Gross Income</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Status</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">AIS</th>
          </tr>
        </thead>
        <tbody>
          {employees.length === 0 ? (
            <tr>
              <td colSpan={7} className="text-center py-12 text-[var(--text-muted)] text-sm">
                No employees found
              </td>
            </tr>
          ) : (
            employees.map((emp) => (
              <tr
                key={emp.id}
                className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors"
              >
                <td className="px-[18px] py-3">
                  <Link href={`/compliance/itr/${emp.id}`} className="hover:underline">
                    <div className="text-xs font-medium text-[var(--text-primary)]">{emp.name}</div>
                    <div className="text-[10px] text-[var(--text-muted)]">{emp.department}</div>
                  </Link>
                </td>
                <td className="px-[18px] py-3 text-xs font-mono text-[var(--text-secondary)]">{emp.pan}</td>
                <td className="px-[18px] py-3"><ITRFormRecommendationPill form={emp.recommendedForm} /></td>
                <td className="px-[18px] py-3"><RegimeIndicator regime={emp.regime} /></td>
                <td className="px-[18px] py-3 text-xs text-right tabular-nums font-semibold text-[var(--text-primary)]">
                  {formatINR(emp.grossIncome)}
                </td>
                <td className="px-[18px] py-3"><FilingStatusPill status={emp.filingStatus} /></td>
                <td className="px-[18px] py-3">
                  <span className={cn(
                    "inline-flex items-center gap-1 text-[10px] font-medium",
                    emp.aisReconciled ? "text-[var(--success)]" : "text-[var(--text-muted)]"
                  )}>
                    <span className={cn("w-1.5 h-1.5 rounded-full", emp.aisReconciled ? "bg-[var(--success)]" : "bg-[var(--border-default)]")} />
                    {emp.aisReconciled ? "Reconciled" : "Pending"}
                  </span>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
