"use client";

import { Search } from "lucide-react";
import { cn } from "@/lib/utils";
import { formatINR } from "@complai/ui-components";
import type { SourceInvoice } from "../../types";

interface SelectInvoiceStepProps {
  invoices: SourceInvoice[];
  search: string;
  onSearch: (v: string) => void;
  onSelect: (inv: SourceInvoice) => void;
}

export function SelectInvoiceStep({
  invoices,
  search,
  onSearch,
  onSelect,
}: SelectInvoiceStepProps) {
  return (
    <div className="space-y-4">
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
        <input
          type="text"
          value={search}
          onChange={(e) => onSearch(e.target.value)}
          placeholder="Search invoices..."
          className={cn(
            "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
            "focus:outline-none focus:border-[var(--accent)]",
            "focus:ring-2 focus:ring-[var(--accent-muted)]"
          )}
        />
      </div>
      <div className="grid gap-3">
        {invoices.map((inv) => (
          <button
            key={inv.id}
            data-testid={`invoice-card-${inv.id}`}
            onClick={() => onSelect(inv)}
            className={cn(
              "flex items-center justify-between p-4 rounded-xl text-left",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "hover:border-[var(--accent)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors"
            )}
          >
            <div className="space-y-1">
              <div className="text-xs font-semibold text-[var(--text-primary)] font-mono">
                {inv.invoiceNo}
              </div>
              <div className="text-[11px] text-[var(--text-muted)]">
                {inv.buyerName} &middot;{" "}
                <span className="font-mono">{inv.buyerGstin}</span>
              </div>
              <div className="text-[10px] text-[var(--text-muted)]">
                {inv.invoiceDate} &middot; {inv.items.length} item
                {inv.items.length > 1 ? "s" : ""}
              </div>
            </div>
            <div className="text-right">
              <div className="text-sm font-bold text-[var(--text-primary)] tabular-nums font-mono">
                {formatINR(inv.totalValue)}
              </div>
              <div className="text-[10px] text-[var(--text-muted)]">
                Taxable: {formatINR(inv.taxableValue)}
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
