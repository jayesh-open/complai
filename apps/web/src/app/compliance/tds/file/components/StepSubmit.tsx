"use client";

import { AlertTriangle, FileSpreadsheet, Loader2 } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import type { TDSFilingData } from "../types";

interface StepSubmitProps {
  data: TDSFilingData;
  loading: boolean;
  hasBlockers: boolean;
  onOpenConfirmModal: () => void;
}

export function StepSubmit({ data, loading, hasBlockers, onOpenConfirmModal }: StepSubmitProps) {
  return (
    <div className="p-6 space-y-4">
      <div>
        <h2 className="text-heading-lg text-[var(--text-primary)]">Step 4: Submit to TRACES</h2>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Review final summary and submit {data.formLabel} via Sandbox gateway.
        </p>
      </div>

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-3">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">Filing Summary</div>
        {[
          { label: "Form", value: data.formLabel },
          { label: "TAN", value: data.tan },
          { label: "Tax Year", value: data.taxYear },
          { label: "Quarter", value: data.quarterLabel },
          { label: "Deductees", value: String(data.deducteeCount) },
          { label: "Total TDS", value: formatINR(data.totalTds) },
          { label: "Surcharge", value: formatINR(data.totalSurcharge) },
          { label: "Cess", value: formatINR(data.totalCess) },
        ].map((item) => (
          <div key={item.label} className="flex gap-3 text-xs">
            <span className="text-[var(--text-muted)] min-w-[100px]">{item.label}:</span>
            <span className="text-[var(--text-primary)] font-medium">{item.value}</span>
          </div>
        ))}
      </div>

      <div className="bg-[var(--danger)]/5 border border-[var(--danger)]/20 rounded-lg p-6 text-center">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--danger)] mb-2">Total Tax Liability</div>
        <div className="text-3xl font-bold font-mono text-[var(--danger)]">{formatINR(data.totalTax)}</div>
      </div>

      <div className="flex items-start gap-2 p-3 bg-[var(--warning-muted)] rounded-lg">
        <AlertTriangle className="w-4 h-4 text-[var(--warning)] flex-shrink-0 mt-0.5" />
        <div className="text-xs text-[var(--warning)]">
          <span className="font-bold">Irreversible action.</span> Once filed, {data.formLabel} for {data.quarterLabel} cannot be revised.
          Ensure all deductee details and amounts are correct.
        </div>
      </div>

      <button
        onClick={onOpenConfirmModal}
        disabled={loading || hasBlockers}
        className={cn(
          "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
          hasBlockers
            ? "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
            : "bg-[var(--danger)] text-white hover:opacity-90",
          loading && "opacity-50 cursor-not-allowed",
        )}
      >
        {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <FileSpreadsheet className="w-4 h-4" />}
        {loading ? "Submitting…" : hasBlockers ? "Blocked — Fix DTAA evidence" : `File ${data.formLabel}`}
      </button>
    </div>
  );
}
