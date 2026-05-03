"use client";

import { cn } from "@/lib/utils";
import type { AISMismatch } from "../types";
import { AISReconciliationMismatchRow } from "./AISReconciliationMismatchRow";
import { ShieldAlert, ShieldCheck } from "lucide-react";

interface AISReconciliationPanelProps {
  mismatches: AISMismatch[];
  onResolve?: (id: string, reason: string) => void;
  className?: string;
}

export function AISReconciliationPanel({
  mismatches,
  onResolve,
  className,
}: AISReconciliationPanelProps) {
  const unresolvedErrors = mismatches.filter((m) => !m.resolved && m.severity === "error");
  const unresolvedWarnings = mismatches.filter((m) => !m.resolved && m.severity === "warn");
  const resolvedCount = mismatches.filter((m) => m.resolved).length;
  const hasBlockingErrors = unresolvedErrors.length > 0;

  return (
    <div
      className={cn(
        "bg-[var(--bg-secondary)] border rounded-xl p-4",
        hasBlockingErrors ? "border-[var(--danger)]" : "border-[var(--border-default)]",
        className
      )}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          {hasBlockingErrors ? (
            <ShieldAlert className="w-4 h-4 text-[var(--danger)]" />
          ) : (
            <ShieldCheck className="w-4 h-4 text-[var(--success)]" />
          )}
          <h3 className="text-xs font-semibold text-[var(--text-primary)]">
            AIS Reconciliation
          </h3>
        </div>
        <div className="flex items-center gap-3 text-[10px] font-medium">
          {unresolvedErrors.length > 0 && (
            <span className="text-[var(--danger)]">{unresolvedErrors.length} error{unresolvedErrors.length > 1 ? "s" : ""}</span>
          )}
          {unresolvedWarnings.length > 0 && (
            <span className="text-[var(--warning)]">{unresolvedWarnings.length} warning{unresolvedWarnings.length > 1 ? "s" : ""}</span>
          )}
          <span className="text-[var(--success)]">{resolvedCount} resolved</span>
        </div>
      </div>

      {hasBlockingErrors && (
        <div className="mb-3 px-3 py-2 rounded-lg bg-[var(--danger-muted)] border border-[var(--danger)] text-[11px] text-[var(--danger)] font-medium">
          Filing blocked — resolve all error-severity mismatches before submission
        </div>
      )}

      <div className="space-y-2">
        {mismatches.map((m) => (
          <AISReconciliationMismatchRow
            key={m.id}
            mismatch={m}
            onResolve={onResolve}
          />
        ))}
      </div>

      {mismatches.length === 0 && (
        <div className="text-center py-8 text-[var(--text-muted)] text-sm">
          No mismatches found — AIS data fully reconciled
        </div>
      )}
    </div>
  );
}
