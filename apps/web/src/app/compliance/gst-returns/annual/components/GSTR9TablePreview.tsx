"use client";

import { cn } from "@/lib/utils";
import type { GSTR9Table } from "../types";
import { FileText } from "lucide-react";

function formatINR(amount: number): string {
  if (amount === 0) return "—";
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

interface GSTR9TablePreviewProps {
  table: GSTR9Table;
  className?: string;
}

export function GSTR9TablePreview({ table, className }: GSTR9TablePreviewProps) {
  return (
    <div className={cn("bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden", className)}>
      <div className="px-4 py-3 border-b border-[var(--border-default)] flex items-center gap-2">
        <span className="text-[10px] font-bold bg-[var(--accent-muted)] text-[var(--accent)] px-2 py-0.5 rounded">
          Table {table.tableNumber}
        </span>
        <h3 className="text-xs font-semibold text-[var(--text-primary)] truncate">{table.title}</h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-[11px]">
          <thead>
            <tr className="border-b border-[var(--border-default)] bg-[var(--bg-tertiary)]">
              <th className="px-3 py-2 text-left font-semibold text-[var(--text-muted)]">Sr.</th>
              <th className="px-3 py-2 text-left font-semibold text-[var(--text-muted)]">Description</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">Taxable</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">CGST</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">SGST</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">IGST</th>
              <th className="px-3 py-2 text-right font-semibold text-[var(--text-muted)]">Cess</th>
              <th className="px-3 py-2 text-left font-semibold text-[var(--text-muted)]">Source</th>
            </tr>
          </thead>
          <tbody>
            {table.rows.map((row) => (
              <tr key={row.serial} className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)]">
                <td className="px-3 py-2 font-mono font-semibold text-[var(--accent)]">{row.serial}</td>
                <td className="px-3 py-2 text-[var(--text-primary)] max-w-[240px] truncate">{row.description}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-primary)]">{formatINR(row.taxableValue)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.cgst)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.sgst)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-secondary)]">{formatINR(row.igst)}</td>
                <td className="px-3 py-2 text-right tabular-nums text-[var(--text-muted)]">{formatINR(row.cess)}</td>
                <td className="px-3 py-2">
                  {row.sourceReturn && (
                    <span className="inline-flex items-center gap-1 text-[9px] font-medium px-1.5 py-0.5 rounded bg-[var(--bg-tertiary)] text-[var(--text-muted)]" title={`Sourced from ${row.sourceReturn}`}>
                      <FileText className="w-2.5 h-2.5" />
                      {row.sourceReturn}
                    </span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
