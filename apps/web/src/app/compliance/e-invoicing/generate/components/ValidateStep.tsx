"use client";

import { CheckCircle2, AlertTriangle, FileCheck2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { formatINR } from "@complai/ui-components";
import type { SourceInvoice } from "../../types";

interface ValidateStepProps {
  invoice: SourceInvoice;
  errors: string[];
  onBack: () => void;
  onGenerate: () => void;
}

export function ValidateStep({
  invoice,
  errors,
  onBack,
  onGenerate,
}: ValidateStepProps) {
  const isValid = errors.length === 0;
  return (
    <div className="space-y-6 max-w-2xl">
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 space-y-3">
        <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-3">
          Payload Summary
        </h3>
        <InfoRow label="Invoice">{invoice.invoiceNo}</InfoRow>
        <InfoRow label="Date">{invoice.invoiceDate}</InfoRow>
        <InfoRow label="Supplier GSTIN">
          <span className="font-mono">{invoice.gstin}</span>
        </InfoRow>
        <InfoRow label="Buyer">
          {invoice.buyerName}{" "}
          <span className="font-mono text-[10px] text-[var(--text-muted)]">
            ({invoice.buyerGstin})
          </span>
        </InfoRow>
        <InfoRow label="Items">{invoice.items.length}</InfoRow>
        <InfoRow label="Taxable Value">
          {formatINR(invoice.taxableValue)}
        </InfoRow>
        <InfoRow label="Total Value">
          <span className="font-semibold">{formatINR(invoice.totalValue)}</span>
        </InfoRow>
      </div>

      {isValid ? (
        <div className="flex items-center gap-2 px-4 py-3 rounded-lg bg-[var(--success-muted)] border border-[var(--success-border)]">
          <CheckCircle2 className="w-4 h-4 text-[var(--success)]" />
          <span className="text-xs font-medium text-[var(--success)]">
            Payload validated — ready to submit
          </span>
        </div>
      ) : (
        <div className="space-y-2">
          {errors.map((err) => (
            <div
              key={err}
              className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--danger-muted)] border border-[var(--danger-border)]"
            >
              <AlertTriangle className="w-4 h-4 text-[var(--danger)]" />
              <span className="text-xs text-[var(--danger)]">{err}</span>
            </div>
          ))}
        </div>
      )}

      <div className="flex items-center gap-3">
        <button
          onClick={onBack}
          className={cn(
            "px-4 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
            "transition-colors"
          )}
        >
          Back
        </button>
        <button
          data-testid="generate-irn-button"
          onClick={onGenerate}
          disabled={!isValid}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
            isValid
              ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
              : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
          )}
        >
          <FileCheck2 className="w-3.5 h-3.5" />
          Generate IRN
        </button>
      </div>
    </div>
  );
}

function InfoRow({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex gap-4 text-xs">
      <span className="text-[var(--text-muted)] min-w-[120px]">{label}</span>
      <span className="text-[var(--text-primary)]">{children}</span>
    </div>
  );
}
