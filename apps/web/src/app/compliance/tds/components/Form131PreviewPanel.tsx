"use client";

import { formatINR } from "@complai/ui-components";
import type { Form131Detail } from "../certificates/types";

export function Form131PreviewPanel({ data }: { data: Form131Detail }) {
  return (
    <div className="space-y-4 text-xs">
      <div className="border border-[var(--border-default)] rounded-lg p-4 space-y-2">
        <h3 className="text-body-sm font-semibold text-[var(--text-primary)]">Section-wise Breakdown</h3>
        <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
                <th className="text-left px-3 py-2 font-medium">Payment Code</th>
                <th className="text-left px-3 py-2 font-medium">Nature</th>
                <th className="text-right px-3 py-2 font-medium">Amount</th>
                <th className="text-right px-3 py-2 font-medium">Rate</th>
                <th className="text-right px-3 py-2 font-medium">TDS</th>
                <th className="text-right px-3 py-2 font-medium">Surcharge</th>
                <th className="text-right px-3 py-2 font-medium">Cess</th>
                <th className="text-left px-3 py-2 font-medium">Date</th>
              </tr>
            </thead>
            <tbody>
              {data.sectionBreakdown.map((s, i) => (
                <tr key={i} className="border-t border-[var(--border-default)]">
                  <td className="px-3 py-2 font-mono text-[var(--text-primary)]">{s.paymentCode}</td>
                  <td className="px-3 py-2 text-[var(--text-muted)] max-w-[160px] truncate">{s.paymentLabel}</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-primary)]">{formatINR(s.amount)}</td>
                  <td className="px-3 py-2 text-right text-[var(--text-muted)]">{s.tdsRate}%</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-primary)]">{formatINR(s.tdsAmount)}</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-muted)]">{formatINR(s.surcharge)}</td>
                  <td className="px-3 py-2 text-right font-mono text-[var(--text-muted)]">{formatINR(s.cess)}</td>
                  <td className="px-3 py-2 text-[var(--text-muted)]">{s.dateOfPayment}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="flex justify-end gap-6 pt-2 font-medium text-[var(--text-primary)]">
          <span>Total TDS: {formatINR(data.totalTds)}</span>
          <span>Total Tax: {formatINR(data.totalTax)}</span>
        </div>
      </div>

      {data.challanRefs.length > 0 && (
        <div className="border border-[var(--border-default)] rounded-lg p-4 space-y-2">
          <h3 className="text-body-sm font-semibold text-[var(--text-primary)]">Challan References</h3>
          <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
                  <th className="text-left px-3 py-2 font-medium">BSR Code</th>
                  <th className="text-left px-3 py-2 font-medium">Challan No.</th>
                  <th className="text-left px-3 py-2 font-medium">Date</th>
                  <th className="text-right px-3 py-2 font-medium">Amount</th>
                </tr>
              </thead>
              <tbody>
                {data.challanRefs.map((c) => (
                  <tr key={c.challanSerial} className="border-t border-[var(--border-default)]">
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
      )}
    </div>
  );
}
