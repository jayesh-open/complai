"use client";

import { Upload } from "lucide-react";
import { cn } from "@/lib/utils";
import { formatINR } from "@complai/ui-components";
import type { SourceInvoice } from "../../types";

interface InvoiceSelectorProps {
  invoices: SourceInvoice[];
  selected: Set<string>;
  onToggle: (id: string) => void;
  onToggleAll: () => void;
  onSubmit: () => void;
}

export function InvoiceSelector({
  invoices,
  selected,
  onToggle,
  onToggleAll,
  onSubmit,
}: InvoiceSelectorProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <label className="flex items-center gap-2 text-xs text-[var(--text-secondary)]">
          <input
            type="checkbox"
            checked={selected.size === invoices.length}
            onChange={onToggleAll}
            className="accent-[var(--accent)]"
          />
          Select all ({invoices.length})
        </label>
        <div className="flex items-center gap-3">
          <span className="text-xs text-[var(--text-muted)]">
            {selected.size} selected (max 50)
          </span>
          <button
            onClick={onSubmit}
            disabled={selected.size === 0}
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
              selected.size > 0
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
            )}
          >
            <Upload className="w-3.5 h-3.5" />
            Generate {selected.size} IRNs
          </button>
        </div>
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[var(--border-default)]">
              <th className="w-10 px-4 py-2" />
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Invoice
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Buyer
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Date
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                Value
              </th>
            </tr>
          </thead>
          <tbody>
            {invoices.map((inv) => (
              <tr
                key={inv.id}
                onClick={() => onToggle(inv.id)}
                className={cn(
                  "border-b border-[var(--border-default)] last:border-b-0",
                  "cursor-pointer hover:bg-[var(--bg-tertiary)] transition-colors",
                  selected.has(inv.id) && "bg-[var(--accent-muted)]"
                )}
              >
                <td className="px-4 py-2.5 text-center">
                  <input
                    type="checkbox"
                    checked={selected.has(inv.id)}
                    onChange={() => onToggle(inv.id)}
                    className="accent-[var(--accent)]"
                  />
                </td>
                <td className="px-4 py-2.5 text-xs font-mono text-[var(--text-primary)]">
                  {inv.invoiceNo}
                </td>
                <td className="px-4 py-2.5">
                  <div className="text-xs text-[var(--text-primary)]">
                    {inv.buyerName}
                  </div>
                  <div className="text-[10px] font-mono text-[var(--text-muted)]">
                    {inv.buyerGstin}
                  </div>
                </td>
                <td className="px-4 py-2.5 text-xs text-[var(--text-muted)]">
                  {inv.invoiceDate}
                </td>
                <td className="px-4 py-2.5 text-xs font-mono font-semibold text-[var(--text-primary)] text-right tabular-nums">
                  {formatINR(inv.totalValue)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
