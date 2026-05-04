"use client";

import { cn } from "@/lib/utils";
import type { LateITCEntry } from "../types";
import { AlertCircle, Info } from "lucide-react";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

const TABLE_LABELS: Record<string, string> = {
  "6H": "Table 6H — ITC reclaimed (Rule 37 reversal)",
  "8C": "Table 8C — ITC gap rectification (GSTR-3B vs GSTR-2B)",
  "13": "Table 13 — ITC on prior-year invoices",
};

interface LateITCReclaimPanelProps {
  entries: LateITCEntry[];
  className?: string;
}

export function LateITCReclaimPanel({ entries, className }: LateITCReclaimPanelProps) {
  const totalReclaim = entries.reduce((sum, e) => sum + e.amount, 0);

  return (
    <div className={cn("bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4", className)}>
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-[var(--warning)]" />
          <h3 className="text-xs font-semibold text-[var(--text-primary)]">Late ITC Reclaim</h3>
        </div>
        <span className="text-xs font-bold tabular-nums text-[var(--success)]">
          +{formatINR(totalReclaim)}
        </span>
      </div>

      <p className="text-[10px] text-[var(--text-muted)] mb-3">
        ITC that was reversed in monthly returns but is now reclaimable in the annual return under specific provisions.
      </p>

      <div className="space-y-3">
        {entries.map((entry) => (
          <div
            key={entry.table}
            className="border border-[var(--border-default)] rounded-lg p-3"
          >
            <div className="flex items-center justify-between mb-1.5">
              <span className="text-[10px] font-semibold text-[var(--accent)]">
                {TABLE_LABELS[entry.table] ?? entry.table}
              </span>
              <span className="text-xs font-bold tabular-nums text-[var(--text-primary)]">
                {formatINR(entry.amount)}
              </span>
            </div>
            <p className="text-[11px] text-[var(--text-secondary)] mb-1.5">
              {entry.description}
            </p>
            <div className="flex items-center gap-1.5 text-[10px] text-[var(--text-muted)]">
              <span className="font-medium">Period: {entry.period}</span>
            </div>
            <div className="flex items-start gap-1.5 mt-2 text-[10px] text-[var(--info)] bg-[var(--info-muted)] rounded-md px-2 py-1.5">
              <Info className="w-3 h-3 flex-shrink-0 mt-0.5" />
              <span>{entry.rule}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
