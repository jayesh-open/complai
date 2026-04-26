"use client";

import { AlertTriangle, FileSpreadsheet, Loader2 } from "lucide-react";
import { cn, formatINR } from "@/lib/utils";
import type { GSTR3BData } from "./types";

interface StepFileProps {
  data: GSTR3BData;
  signMethod: "dsc" | "evc";
  loading: boolean;
  onOpenConfirmModal: () => void;
}

export function StepFile({ data, signMethod, loading, onOpenConfirmModal }: StepFileProps) {
  const totalPayable = data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst;

  return (
    <div className="p-6 space-y-4" data-testid="step-file">
      <div>
        <h2 className="text-heading-lg text-foreground">Step 5: File GSTR-3B</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Review final summary and file your return on the GST portal.
        </p>
      </div>

      {/* Final Summary */}
      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-3">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">Filing Summary</div>
        {[
          { label: "GSTIN", value: data.gstin },
          { label: "Period", value: data.periodLabel },
          { label: "Total Liability", value: formatINR(data.totalLiability.cgst + data.totalLiability.sgst + data.totalLiability.igst) },
          { label: "Total ITC Used", value: formatINR(data.totalITC.cgst + data.totalITC.sgst + data.totalITC.igst) },
          { label: "Signing Method", value: signMethod === "dsc" ? "Digital Signature" : "Electronic Verification" },
        ].map((item) => (
          <div key={item.label} className="flex gap-3 text-xs">
            <span className="text-foreground-muted min-w-[120px]">{item.label}:</span>
            <span className="text-foreground font-medium">{item.value}</span>
          </div>
        ))}
      </div>

      {/* Tax Amount Prominent Display */}
      <div className="bg-[var(--danger)]/5 border border-[var(--danger)]/20 rounded-lg p-6 text-center">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--danger)] mb-2">Cash Payment Due</div>
        <div className="text-3xl font-bold font-mono text-[var(--danger)]" data-testid="total-payable">
          {formatINR(totalPayable)}
        </div>
        <div className="text-xs text-foreground-muted mt-2">
          This amount will be debited from your cash ledger upon filing.
        </div>
      </div>

      {/* Warning */}
      <div className="flex items-start gap-2 p-3 bg-[var(--warning-muted)] rounded-lg">
        <AlertTriangle className="w-4 h-4 text-[var(--warning)] flex-shrink-0 mt-0.5" />
        <div className="text-xs text-[var(--warning)]">
          <span className="font-bold">Irreversible action.</span> Once filed, GSTR-3B cannot be revised.
          Ensure all values are correct before proceeding.
        </div>
      </div>

      <button
        data-testid="file-button"
        onClick={onOpenConfirmModal}
        disabled={loading}
        className={cn(
          "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
          "bg-[var(--danger)] text-white hover:opacity-90",
          loading && "opacity-50 cursor-not-allowed",
        )}
      >
        {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <FileSpreadsheet className="w-4 h-4" />}
        {loading ? "Filing..." : "File GSTR-3B"}
      </button>
    </div>
  );
}
