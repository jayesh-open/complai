"use client";

import { cn } from "@/lib/utils";
import { ReconciliationSplitPane } from "@complai/ui-components";
import type { GSTR9CMismatch } from "../types";
import { MismatchSeverityBadge } from "./MismatchSeverityBadge";

interface ReconciliationDeltaTableProps {
  section: "II" | "III" | "IV";
  sectionTitle: string;
  mismatches: GSTR9CMismatch[];
  className?: string;
}

function fmtINR(amount: number): string {
  if (Math.abs(amount) >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (Math.abs(amount) >= 100_000) return `₹${(amount / 100_000).toFixed(2)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

function DeltaRow({ m }: { m: GSTR9CMismatch }) {
  const delta = m.booksAmount - m.gstr9Amount;
  return (
    <div className="flex items-center gap-2 px-3 py-2 rounded-lg text-xs hover:bg-[var(--bg-tertiary)] transition-colors">
      <MismatchSeverityBadge severity={m.severity} />
      <span className="flex-1 text-[var(--text-secondary)] truncate">{m.description}</span>
      <span className={cn("tabular-nums font-semibold", delta > 0 ? "text-[var(--danger)]" : "text-[var(--success)]")}>
        {delta > 0 ? "+" : ""}{fmtINR(delta)}
      </span>
    </div>
  );
}

export function ReconciliationDeltaTable({ section, sectionTitle, mismatches, className }: ReconciliationDeltaTableProps) {
  const filtered = mismatches.filter((m) => m.section === section);

  if (filtered.length === 0) {
    return (
      <div className={cn("border border-[var(--border-default)] rounded-xl p-4", className)}>
        <p className="text-xs text-[var(--text-muted)]">No mismatches in Part {section}</p>
      </div>
    );
  }

  return (
    <ReconciliationSplitPane
      className={className}
      leftTitle={`Books (Audited) — Part ${section}`}
      rightTitle={`GSTR-9 Filed — Part ${section}`}
      leftContent={
        <div className="space-y-1">
          {filtered.map((m) => (
            <div key={m.id} className="flex items-center justify-between px-3 py-2 rounded-lg text-xs">
              <span className="text-[var(--text-secondary)] truncate max-w-[70%]">{m.description}</span>
              <span className="tabular-nums font-medium text-[var(--text-primary)]">{fmtINR(m.booksAmount)}</span>
            </div>
          ))}
        </div>
      }
      rightContent={
        <div className="space-y-1">
          {filtered.map((m) => (
            <DeltaRow key={m.id} m={m} />
          ))}
        </div>
      }
    />
  );
}
