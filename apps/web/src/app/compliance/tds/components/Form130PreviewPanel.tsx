"use client";

import { formatINR } from "@complai/ui-components";
import type { Form130Detail } from "../certificates/types";

export function Form130PreviewPanel({ data }: { data: Form130Detail }) {
  return (
    <div className="space-y-4 text-xs">
      <div className="border border-[var(--border-default)] rounded-lg p-4 space-y-3">
        <h3 className="text-body-sm font-semibold text-[var(--text-primary)]">Annexure I — Salary Breakup</h3>
        <div className="grid grid-cols-2 gap-x-8 gap-y-1">
          {data.salaryBreakup.map((b) => (
            <div key={b.component} className="flex justify-between">
              <span className="text-[var(--text-muted)]">{b.component}</span>
              <span className="font-mono text-[var(--text-primary)]">{formatINR(b.amount)}</span>
            </div>
          ))}
        </div>
        <div className="border-t border-[var(--border-default)] pt-2 flex justify-between font-medium">
          <span className="text-[var(--text-primary)]">Gross Total</span>
          <span className="font-mono text-[var(--text-primary)]">{formatINR(data.totalGross)}</span>
        </div>
      </div>

      <div className="border border-[var(--border-default)] rounded-lg p-4 space-y-3">
        <h3 className="text-body-sm font-semibold text-[var(--text-primary)]">Annexure II — Deductions & Tax</h3>
        <div className="space-y-1">
          <div className="flex justify-between">
            <span className="text-[var(--text-muted)]">Standard Deduction</span>
            <span className="font-mono text-[var(--text-primary)]">{formatINR(data.standardDeduction)}</span>
          </div>
          {data.deductions80C.map((d) => (
            <div key={d.section} className="flex justify-between">
              <span className="text-[var(--text-muted)]">{d.section} — {d.description}</span>
              <span className="font-mono text-[var(--text-primary)]">{formatINR(d.verified)}</span>
            </div>
          ))}
          <div className="border-t border-[var(--border-default)] pt-2 flex justify-between">
            <span className="text-[var(--text-muted)]">Taxable Income</span>
            <span className="font-mono font-medium text-[var(--text-primary)]">{formatINR(data.taxableIncome)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-[var(--text-muted)]">Tax Computed</span>
            <span className="font-mono font-medium text-[var(--text-primary)]">{formatINR(data.taxComputed)}</span>
          </div>
        </div>
      </div>

      <div className="border border-[var(--border-default)] rounded-lg p-4 space-y-2">
        <h3 className="text-body-sm font-semibold text-[var(--text-primary)]">Challan References</h3>
        <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
                <th className="text-left px-3 py-2 font-medium">Quarter</th>
                <th className="text-left px-3 py-2 font-medium">BSR Code</th>
                <th className="text-left px-3 py-2 font-medium">Challan No.</th>
                <th className="text-left px-3 py-2 font-medium">Date</th>
                <th className="text-right px-3 py-2 font-medium">Amount</th>
              </tr>
            </thead>
            <tbody>
              {data.tdsChallanRefs.map((c) => (
                <tr key={c.challanSerial} className="border-t border-[var(--border-default)]">
                  <td className="px-3 py-2 text-[var(--text-muted)]">{c.quarter}</td>
                  <td className="px-3 py-2 font-mono text-[var(--text-muted)]">{c.bsrCode}</td>
                  <td className="px-3 py-2 font-mono text-[var(--text-primary)]">{c.challanSerial}</td>
                  <td className="px-3 py-2 text-[var(--text-muted)]">{c.depositDate}</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-primary)]">{formatINR(c.amount)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
