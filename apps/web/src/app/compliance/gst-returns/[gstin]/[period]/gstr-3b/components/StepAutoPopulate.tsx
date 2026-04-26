"use client";

import { CheckCircle2, Download, Loader2, AlertTriangle } from "lucide-react";
import { cn, formatINR } from "@/lib/utils";
import type { GSTR3BData } from "./types";

interface StepAutoPopulateProps {
  data: GSTR3BData | null;
  loading: boolean;
  onPopulate: () => void;
}

export function StepAutoPopulate({ data, loading, onPopulate }: StepAutoPopulateProps) {
  return (
    <div className="p-6 space-y-4" data-testid="step-auto-populate">
      <div>
        <h2 className="text-heading-lg text-foreground">Step 1: Auto-Populate GSTR-3B</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Pull data from filed GSTR-1, GSTR-2B/IMS, and compute ITC eligibility.
        </p>
      </div>

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)] mb-2">Data Sources</div>
        {[
          { label: "GSTR-1 (filed)", desc: "Outward supplies, tax liability", status: "ready" },
          { label: "GSTR-2B (auto-drafted)", desc: "Inward supplies, ITC available", status: "ready" },
          { label: "IMS Actions", desc: "Accept/Reject decisions on supplier invoices", status: "ready" },
          { label: "Cash & Credit Ledger", desc: "Available balances for offset", status: "ready" },
        ].map((src) => (
          <div key={src.label} className="flex items-center gap-3 text-xs">
            <CheckCircle2 className="w-3.5 h-3.5 text-[var(--success)] flex-shrink-0" />
            <div className="flex-1">
              <span className="font-medium text-foreground">{src.label}</span>
              <span className="text-foreground-muted ml-2">— {src.desc}</span>
            </div>
          </div>
        ))}
      </div>

      {!data && (
        <button
          data-testid="populate-button"
          onClick={onPopulate}
          disabled={loading}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
            "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
            loading && "opacity-50 cursor-not-allowed",
          )}
        >
          {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Download className="w-4 h-4" />}
          {loading ? "Fetching data..." : "Auto-Populate from Sources"}
        </button>
      )}

      {data && (
        <div className="space-y-3">
          <div className="flex items-center gap-2 p-3 bg-[var(--success-muted)] rounded-lg">
            <CheckCircle2 className="w-5 h-5 text-[var(--success)]" />
            <span className="text-xs text-[var(--success)] font-medium">
              GSTR-3B auto-populated from all sources
            </span>
          </div>

          {data.flags.length > 0 && (
            <div className="space-y-1.5">
              {data.flags.map((flag, i) => (
                <div key={i} className="flex items-center gap-2 p-2 bg-[var(--warning-muted)] rounded-lg">
                  <AlertTriangle className="w-3.5 h-3.5 text-[var(--warning)]" />
                  <span className="text-[11px] text-[var(--warning)] font-medium">{flag}</span>
                </div>
              ))}
            </div>
          )}

          <div className="grid grid-cols-3 gap-3">
            {[
              { label: "Total Liability", value: data.totalLiability.cgst + data.totalLiability.sgst + data.totalLiability.igst, color: "text-foreground" },
              { label: "Total ITC", value: data.totalITC.cgst + data.totalITC.sgst + data.totalITC.igst, color: "text-[var(--success)]" },
              { label: "Net Payable", value: data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst, color: "text-[var(--danger)]" },
            ].map((stat) => (
              <div key={stat.label} className="bg-[var(--bg-tertiary)] rounded-lg p-3 text-center">
                <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">{stat.label}</div>
                <div className={cn("text-lg font-bold font-mono", stat.color)}>{formatINR(stat.value)}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
