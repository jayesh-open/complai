"use client";

import { Download, Loader2, Database } from "lucide-react";
import { cn } from "@complai/ui-components";
import type { TDSFilingData, TDSFormType } from "../types";
import { FORM_LABELS, FORM_SECTIONS } from "../types";

interface StepPullProps {
  formType: TDSFormType;
  taxYear: string;
  quarter: string;
  data: TDSFilingData | null;
  loading: boolean;
  onPull: () => void;
}

export function StepPull({ formType, taxYear, quarter, data, loading, onPull }: StepPullProps) {
  return (
    <div className="p-6 space-y-4">
      <div>
        <h2 className="text-heading-lg text-[var(--text-primary)]">Step 1: Pull Entries</h2>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Pull TDS entries from the system for {FORM_LABELS[formType]}.
        </p>
      </div>

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">Pull Parameters</div>
        {[
          { label: "Form", value: FORM_LABELS[formType] },
          { label: "Section", value: FORM_SECTIONS[formType] },
          { label: "Tax Year", value: taxYear },
          { label: "Quarter", value: quarter.toUpperCase() },
          { label: "Source", value: formType === "138" ? "Apex Gateway (Salary)" : "TDS Entries" },
        ].map((item) => (
          <div key={item.label} className="flex gap-3 text-xs">
            <span className="text-[var(--text-muted)] min-w-[100px]">{item.label}:</span>
            <span className="text-[var(--text-primary)] font-medium">{item.value}</span>
          </div>
        ))}
      </div>

      {!data && (
        <button
          onClick={onPull}
          disabled={loading}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
            "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]",
            loading && "opacity-50 cursor-not-allowed",
          )}
        >
          {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Download className="w-4 h-4" />}
          {loading ? "Pulling entries…" : "Pull Entries"}
        </button>
      )}

      {data && (
        <div className="flex items-center gap-3 p-3 bg-[var(--success-muted)] rounded-lg border border-[var(--success-border)]">
          <Database className="w-5 h-5 text-[var(--success)]" />
          <div>
            <p className="text-xs font-semibold text-[var(--text-primary)]">
              {data.entries.length} entries pulled ({data.deducteeCount} deductees)
            </p>
            <p className="text-[10px] text-[var(--text-muted)]">
              Ready for validation. Click Next to proceed.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
