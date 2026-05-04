"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";
import {
  CheckCircle2, AlertTriangle, FileText, ArrowLeft,
} from "lucide-react";

import type { GSTR9Data, GSTR9Step } from "../../../../types";
import { GSTR9_STEPS } from "../../../../types";
import { GSTR9TablePreview } from "../../../../components/GSTR9TablePreview";
import { HSNSummaryTable } from "../../../../components/HSNSummaryTable";
import { LateITCReclaimPanel } from "../../../../components/LateITCReclaimPanel";
import { DSCSigningPlaceholder } from "../../../../components/DSCSigningPlaceholder";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

export function StepBar({ current }: { current: GSTR9Step }) {
  const idx = GSTR9_STEPS.findIndex((s) => s.id === current);
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 sticky top-0 z-10">
      <div className="flex items-center gap-1">
        {GSTR9_STEPS.map((s, i) => (
          <div key={s.id} className="flex items-center flex-1">
            <div className="flex items-center gap-2 flex-1">
              <div className={cn(
                "w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0",
                i < idx ? "bg-[var(--success)] text-white" : s.id === current ? "bg-[var(--accent)] text-[var(--accent-text)]" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
              )}>
                {i < idx ? <CheckCircle2 className="w-4 h-4" /> : s.number}
              </div>
              <span className={cn("text-xs font-medium whitespace-nowrap", s.id === current ? "text-[var(--text-primary)]" : "text-[var(--text-muted)]")}>{s.label}</span>
            </div>
            {i < GSTR9_STEPS.length - 1 && <div className={cn("h-[2px] flex-1 mx-2 rounded", i < idx ? "bg-[var(--success)]" : "bg-[var(--border-default)]")} />}
          </div>
        ))}
      </div>
    </div>
  );
}

export function StepThreshold({ data }: { data: GSTR9Data }) {
  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 1: Threshold Check</h2>
      <div className="grid grid-cols-2 gap-4">
        <div className="border border-[var(--border-default)] rounded-lg p-4">
          <p className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide">Aggregate Turnover (FY {data.fy})</p>
          <p className="text-xl font-bold tabular-nums text-[var(--text-primary)] mt-1">{formatINR(data.turnover)}</p>
        </div>
        <div className="border border-[var(--border-default)] rounded-lg p-4">
          <p className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide">GSTR-9 Applicability</p>
          <div className="flex items-center gap-2 mt-2">
            {data.thresholdExceeded ? (
              <>
                <AlertTriangle className="w-4 h-4 text-[var(--warning)]" />
                <span className="text-xs font-semibold text-[var(--warning)]">Required — turnover exceeds ₹2 Cr threshold</span>
              </>
            ) : (
              <>
                <CheckCircle2 className="w-4 h-4 text-[var(--success)]" />
                <span className="text-xs font-semibold text-[var(--success)]">Optional — turnover below ₹2 Cr threshold</span>
              </>
            )}
          </div>
        </div>
      </div>
      <div className="flex items-start gap-2 bg-[var(--info-muted)] border border-[var(--info)] rounded-lg p-3 text-xs text-[var(--info)]">
        <FileText className="w-4 h-4 flex-shrink-0 mt-0.5" />
        <p>GSTR-9 is mandatory for taxpayers with aggregate turnover exceeding ₹2 crore in the financial year. Composition dealers must file GSTR-9A instead.</p>
      </div>
    </div>
  );
}

export function StepReviewTables({ data }: { data: GSTR9Data }) {
  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 2: Review Aggregated Tables</h2>
      <p className="text-xs text-[var(--text-muted)]">Auto-populated from GSTR-1, GSTR-3B, and GSTR-2B. Review and confirm each table.</p>
      <div className="space-y-4">
        {data.tables.map((t) => (
          <GSTR9TablePreview key={t.tableNumber} table={t} />
        ))}
      </div>
    </div>
  );
}

export function StepLateITC({ data }: { data: GSTR9Data }) {
  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 3: Late ITC Reclaim</h2>
      <p className="text-xs text-[var(--text-muted)]">ITC reversed in monthly returns that is now reclaimable in the annual return.</p>
      <LateITCReclaimPanel entries={data.lateITC} />
    </div>
  );
}

export function StepHSN({ data }: { data: GSTR9Data }) {
  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 4: HSN Summary</h2>
      <p className="text-xs text-[var(--text-muted)]">HSN-wise breakup of outward supplies (Table 18). Verify digit-tier compliance based on turnover.</p>
      <HSNSummaryTable rows={data.hsnRows} />
    </div>
  );
}

export function StepFeesDemands({ data }: { data: GSTR9Data }) {
  const categories = { demand: "Demand", refund: "Refund", late_fee: "Late Fee" } as const;
  const catStyle = { demand: "text-[var(--danger)]", refund: "text-[var(--success)]", late_fee: "text-[var(--warning)]" } as const;

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 5: Fees, Demands &amp; Refunds</h2>
      <p className="text-xs text-[var(--text-muted)]">Tables 14–17: additional tax, refunds, demands, and late fees.</p>
      <div className="border border-[var(--border-default)] rounded-xl overflow-hidden">
        <table className="w-full text-xs">
          <thead>
            <tr className="bg-[var(--bg-tertiary)] border-b border-[var(--border-default)]">
              <th className="px-4 py-2 text-left font-semibold text-[var(--text-muted)]">Table</th>
              <th className="px-4 py-2 text-left font-semibold text-[var(--text-muted)]">Description</th>
              <th className="px-4 py-2 text-center font-semibold text-[var(--text-muted)]">Category</th>
              <th className="px-4 py-2 text-right font-semibold text-[var(--text-muted)]">Amount</th>
            </tr>
          </thead>
          <tbody>
            {data.feesAndDemands.map((row) => (
              <tr key={row.table} className="border-b border-[var(--border-default)] last:border-b-0">
                <td className="px-4 py-2 font-mono font-semibold text-[var(--accent)]">{row.table}</td>
                <td className="px-4 py-2 text-[var(--text-secondary)]">{row.description}</td>
                <td className="px-4 py-2 text-center">
                  <span className={cn("text-[10px] font-semibold", catStyle[row.category])}>{categories[row.category]}</span>
                </td>
                <td className="px-4 py-2 text-right tabular-nums font-medium text-[var(--text-primary)]">{formatINR(row.amount)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function StepSubmit({ data, signed, onSigned }: { data: GSTR9Data; signed: boolean; onSigned: () => void }) {
  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-[var(--text-primary)]">Step 6: Sign &amp; Submit</h2>
      <div className="grid grid-cols-2 gap-4">
        <div className="border border-[var(--border-default)] rounded-lg p-4">
          <p className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide">Turnover</p>
          <p className="text-lg font-bold tabular-nums text-[var(--text-primary)]">{formatINR(data.turnover)}</p>
        </div>
        <div className="border border-[var(--border-default)] rounded-lg p-4">
          <p className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide">Tables Reviewed</p>
          <p className="text-lg font-bold tabular-nums text-[var(--success)]">{data.tables.length}</p>
        </div>
      </div>
      <DSCSigningPlaceholder onSigned={onSigned} label="Sign GSTR-9 with DSC" />
      {signed && (
        <div className="flex items-start gap-2 bg-[var(--info-muted)] border border-[var(--info)] rounded-lg p-3 text-xs text-[var(--info)]">
          <FileText className="w-4 h-4 flex-shrink-0 mt-0.5" />
          <p>DSC applied. Click &ldquo;File GSTR-9&rdquo; to submit to GSTN. This action is irreversible.</p>
        </div>
      )}
    </div>
  );
}

export function StepAcknowledge({ data, arn }: { data: GSTR9Data; arn: string }) {
  return (
    <div className="space-y-4 text-center py-4">
      <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
        <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
      </div>
      <h2 className="text-lg font-bold text-[var(--text-primary)]">GSTR-9 Filed Successfully</h2>
      <p className="text-xs text-[var(--text-muted)]">{data.gstin} &middot; FY {data.fy} &middot; {data.legalName}</p>
      <div className="inline-flex flex-col items-center gap-1 bg-[var(--bg-tertiary)] border border-[var(--border-default)] rounded-lg px-6 py-3">
        <span className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide">ARN</span>
        <span className="text-sm font-mono font-bold text-[var(--accent)]">{arn}</span>
      </div>
      <div className="pt-4">
        <Link
          href="/compliance/gst-returns/annual"
          className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to Annual Returns
        </Link>
      </div>
    </div>
  );
}
