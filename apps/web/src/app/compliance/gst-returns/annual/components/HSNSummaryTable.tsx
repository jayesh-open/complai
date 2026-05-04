"use client";

import { cn } from "@/lib/utils";
import type { HSNRow } from "../types";

function formatINR(amount: number): string {
  if (amount === 0) return "—";
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

const TIER_STYLE: Record<number, string> = {
  4: "bg-[var(--success-muted)] text-[var(--success)]",
  6: "bg-[var(--info-muted)] text-[var(--info)]",
  8: "bg-[var(--purple-muted)] text-[var(--purple)]",
};

interface HSNSummaryTableProps {
  rows: HSNRow[];
  className?: string;
}

export function HSNSummaryTable({ rows, className }: HSNSummaryTableProps) {
  return (
    <div className={cn("bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden", className)}>
      <div className="px-4 py-3 border-b border-[var(--border-default)] flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="text-[10px] font-bold bg-[var(--accent-muted)] text-[var(--accent)] px-2 py-0.5 rounded">
            Table 18
          </span>
          <h3 className="text-xs font-semibold text-[var(--text-primary)]">HSN-wise Summary of Outward Supplies</h3>
        </div>
        <div className="flex items-center gap-2">
          {[4, 6, 8].map((tier) => (
            <span key={tier} className={cn("text-[9px] font-semibold px-1.5 py-0.5 rounded", TIER_STYLE[tier])}>
              {tier}-digit
            </span>
          ))}
        </div>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-[11px]">
          <thead>
            <tr className="border-b border-[var(--border-default)] bg-[var(--bg-tertiary)]">
              <th className="px-3 py-2 text-left font-semibold text-[var(--text-muted)]">HSN</th>
              <th className="px-3 py-2 text-left font-semibold text-[var(--text-muted)]">Description</th>
              <th className="px-3 py-2 text-center font-semibold text-[var(--text-muted)]">Tier</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">Qty</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">Taxable</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">CGST</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">SGST</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">IGST</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.hsn} className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)]">
                <td className="px-3 py-2 font-mono font-semibold text-[var(--text-primary)]">{row.hsn}</td>
                <td className="px-3 py-2 text-[var(--text-secondary)] max-w-[200px] truncate">{row.description}</td>
                <td className="px-3 py-2 text-center">
                  <span className={cn("text-[9px] font-semibold px-1.5 py-0.5 rounded", TIER_STYLE[row.digitTier])}>
                    {row.digitTier}D
                  </span>
                </td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{row.quantity.toLocaleString("en-IN")}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-primary)] font-medium">{formatINR(row.taxableValue)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.cgst)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.sgst)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.igst)}</td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr className="bg-[var(--bg-tertiary)] font-semibold">
              <td colSpan={4} className="px-3 py-2 text-[var(--text-primary)]">Total</td>
              <td className="px-3 py-2 text-right tabular-nums text-[var(--text-primary)]">{formatINR(rows.reduce((s, r) => s + r.taxableValue, 0))}</td>
              <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(rows.reduce((s, r) => s + r.cgst, 0))}</td>
              <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(rows.reduce((s, r) => s + r.sgst, 0))}</td>
              <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(rows.reduce((s, r) => s + r.igst, 0))}</td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  );
}
