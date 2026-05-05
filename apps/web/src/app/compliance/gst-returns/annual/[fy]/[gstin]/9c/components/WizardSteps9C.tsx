"use client";

import { Info, Upload, Loader2 } from "lucide-react";
import type { GSTR9CMismatch, AuditedFinancials } from "../../../../types";
import { DSCSigningPlaceholder } from "../../../../components/DSCSigningPlaceholder";
import { MismatchRow } from "../../../../components/MismatchRow";
import { ReconciliationDeltaTable } from "../../../../components/ReconciliationDeltaTable";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

export function StepThresholdCheck({ turnover, required }: { turnover: number; required: boolean }) {
  return (
    <div className="space-y-3" data-testid="step-threshold-check">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Threshold Check</h2>
      <div className="flex items-center gap-4">
        <div className="bg-[var(--bg-tertiary)] rounded-lg px-4 py-3">
          <p className="text-[10px] text-[var(--text-muted)]">Turnover (from GSTR-9)</p>
          <p className="text-lg font-bold text-[var(--text-primary)] tabular-nums">{formatINR(turnover)}</p>
        </div>
        <div className="bg-[var(--bg-tertiary)] rounded-lg px-4 py-3">
          <p className="text-[10px] text-[var(--text-muted)]">9C Threshold</p>
          <p className="text-lg font-bold text-[var(--text-primary)]">₹5 Cr</p>
        </div>
      </div>
      {required ? (
        <div className="flex items-start gap-2 bg-[color-mix(in_srgb,var(--warning)_8%,transparent)] border border-[var(--warning)]/30 rounded-lg px-4 py-3">
          <Info className="w-4 h-4 text-[var(--warning)] mt-0.5 shrink-0" />
          <p className="text-xs text-[var(--text-secondary)]">Turnover exceeds ₹5 Cr — GSTR-9C reconciliation statement is <strong>mandatory</strong>. Proceed to upload audited financials.</p>
        </div>
      ) : (
        <div className="flex items-start gap-2 bg-[color-mix(in_srgb,var(--success)_8%,transparent)] border border-[var(--success)]/30 rounded-lg px-4 py-3">
          <Info className="w-4 h-4 text-[var(--success)] mt-0.5 shrink-0" />
          <p className="text-xs text-[var(--text-secondary)]">Turnover is ≤ ₹5 Cr — GSTR-9C is <strong>optional</strong>. You may skip or proceed voluntarily.</p>
        </div>
      )}
    </div>
  );
}

export function StepUploadFinancials({ audited, uploaded, onUpload }: { audited: AuditedFinancials; uploaded: boolean; onUpload: () => void }) {
  return (
    <div className="space-y-3" data-testid="step-upload-financials">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Upload Audited Financials</h2>
      {!uploaded ? (
        <div className="border-2 border-dashed border-[var(--border-default)] rounded-xl p-8 flex flex-col items-center gap-3">
          <Upload className="w-8 h-8 text-[var(--text-muted)]" />
          <p className="text-xs text-[var(--text-muted)]">Upload CSV/JSON or enter manually</p>
          <button onClick={onUpload} className="px-4 py-2 rounded-lg text-xs font-semibold bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors">
            Use Mock Financials
          </button>
        </div>
      ) : (
        <div className="space-y-2">
          <p className="text-xs text-[var(--success)] font-medium">Audited financials loaded</p>
          <div className="grid grid-cols-2 gap-3">
            <FinancialCard label="Gross Turnover (Audited)" value={audited.grossTurnover} />
            <FinancialCard label="Taxable Turnover (Audited)" value={audited.taxableTurnover} />
            <FinancialCard label="Tax Payable (Total)" value={audited.taxPayable.cgst + audited.taxPayable.sgst + audited.taxPayable.igst + audited.taxPayable.cess} />
            <FinancialCard label="ITC Claimed (Total)" value={audited.itcClaimed.cgst + audited.itcClaimed.sgst + audited.itcClaimed.igst + audited.itcClaimed.cess} />
          </div>
        </div>
      )}
    </div>
  );
}

function FinancialCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-[var(--bg-tertiary)] rounded-lg px-4 py-3">
      <p className="text-[10px] text-[var(--text-muted)]">{label}</p>
      <p className="text-sm font-bold text-[var(--text-primary)] tabular-nums">{formatINR(value)}</p>
    </div>
  );
}

export function StepReconciliation({ mismatches }: { mismatches: GSTR9CMismatch[] }) {
  return (
    <div className="space-y-4" data-testid="step-reconciliation">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Reconciliation by Section</h2>
      <ReconciliationDeltaTable section="II" sectionTitle="Turnover Reconciliation" mismatches={mismatches} />
      <ReconciliationDeltaTable section="III" sectionTitle="Tax Reconciliation" mismatches={mismatches} />
      <ReconciliationDeltaTable section="IV" sectionTitle="ITC Reconciliation" mismatches={mismatches} />
    </div>
  );
}

export function StepResolveMismatches({ mismatches, unresolvedErrors, onResolve }: { mismatches: GSTR9CMismatch[]; unresolvedErrors: number; onResolve: (id: string) => void }) {
  return (
    <div className="space-y-3" data-testid="step-resolve-mismatches">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-[var(--text-primary)]">Resolve Mismatches</h2>
        {unresolvedErrors > 0 && (
          <span className="text-[10px] font-bold text-[var(--danger)] bg-[color-mix(in_srgb,var(--danger)_10%,transparent)] px-2 py-1 rounded">
            {unresolvedErrors} ERROR{unresolvedErrors > 1 ? "S" : ""} unresolved — blocks proceed
          </span>
        )}
      </div>
      <div className="space-y-2">
        {mismatches.map((m) => (
          <MismatchRow key={m.id} mismatch={m} onResolve={onResolve} />
        ))}
      </div>
    </div>
  );
}

export function StepFileDSC({ signed, onSigned, arn, loading }: { signed: boolean; onSigned: () => void; arn: string | null; loading: boolean }) {
  if (arn) {
    return (
      <div className="text-center py-8 space-y-3" data-testid="step-arn-9c">
        <div className="inline-flex items-center justify-center w-14 h-14 rounded-full bg-[var(--success)]/10 mb-2">
          <svg className="w-7 h-7 text-[var(--success)]" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" /></svg>
        </div>
        <h2 className="text-sm font-semibold text-[var(--success)]">GSTR-9C Filed Successfully</h2>
        <p className="text-xs text-[var(--text-muted)]">ARN: <span className="font-mono font-bold text-[var(--text-primary)]">{arn}</span></p>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center gap-2 py-12 text-sm text-[var(--text-muted)]">
        <Loader2 className="w-5 h-5 animate-spin text-[var(--accent)]" />
        Filing GSTR-9C with GSTN...
      </div>
    );
  }

  return (
    <div className="space-y-4" data-testid="step-file-dsc">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">File with DSC</h2>
      <p className="text-xs text-[var(--text-muted)]">GSTR-9C requires mandatory DSC signing. Connect your USB token to proceed.</p>
      <DSCSigningPlaceholder label="Sign GSTR-9C with DSC" onSigned={onSigned} />
    </div>
  );
}
