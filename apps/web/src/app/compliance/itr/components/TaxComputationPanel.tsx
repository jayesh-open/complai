"use client";

import { cn } from "@/lib/utils";
import type { TaxComputation, DeductionItem, TaxRegime } from "../types";
import { Calculator } from "lucide-react";

function formatINR(amount: number): string {
  const abs = Math.abs(amount);
  const sign = amount < 0 ? "-" : "";
  if (abs >= 10_000_000) return `${sign}₹${(abs / 10_000_000).toFixed(2)} Cr`;
  if (abs >= 100_000) return `${sign}₹${(abs / 100_000).toFixed(1)} L`;
  return `${sign}₹${abs.toLocaleString("en-IN")}`;
}

interface TaxComputationPanelProps {
  computation: TaxComputation;
  deductions?: DeductionItem[];
  regime: TaxRegime;
  className?: string;
}

function Row({ label, value, bold, accent, section }: {
  label: string; value: string; bold?: boolean; accent?: string; section?: string;
}) {
  return (
    <div className="flex items-center justify-between py-1">
      <span className={cn("text-[11px]", bold ? "font-semibold text-[var(--text-primary)]" : "text-[var(--text-muted)]")}>
        {label}
        {section && (
          <span className="ml-1 text-[9px] font-mono px-1 py-0.5 rounded bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
            {section}
          </span>
        )}
      </span>
      <span className={cn("text-[11px] tabular-nums font-medium", accent || (bold ? "text-[var(--text-primary)]" : "text-[var(--text-secondary)]"))}>
        {value}
      </span>
    </div>
  );
}

export function TaxComputationPanel({ computation, deductions, regime, className }: TaxComputationPanelProps) {
  const c = computation;

  return (
    <div className={cn("bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4", className)}>
      <div className="flex items-center gap-2 mb-3">
        <Calculator className="w-4 h-4 text-[var(--accent)]" />
        <h3 className="text-xs font-semibold text-[var(--text-primary)]">Tax Computation</h3>
      </div>

      <div className="space-y-0.5">
        <Row label="Total Income" value={formatINR(c.totalIncome)} bold />
        <Row label="Less: Standard Deduction" value={`-${formatINR(c.standardDeduction)}`} section="§16" />

        {regime === "OLD" && deductions && deductions.length > 0 && (
          <div className="border-t border-[var(--border-default)] mt-2 pt-2">
            <p className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">
              Chapter VI-A Deductions (Old Regime)
            </p>
            {deductions.map((d) => (
              <Row
                key={d.section}
                label={`§${d.section}: ${d.label}`}
                value={`-${formatINR(d.declared)}${d.limit > 0 ? ` / ${formatINR(d.limit)}` : ""}`}
              />
            ))}
          </div>
        )}

        <div className="border-t border-[var(--border-default)] mt-2 pt-2">
          <Row label="Taxable Income" value={formatINR(c.taxableIncome)} bold />
        </div>

        <div className="border-t border-[var(--border-default)] mt-2 pt-2">
          <p className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">
            Slab Tax (§197)
          </p>
          {c.slabs.filter((s) => s.tax > 0).map((s) => (
            <Row
              key={s.from}
              label={`${formatINR(s.from)} – ${s.to ? formatINR(s.to) : "above"} @ ${s.rate}%`}
              value={formatINR(s.tax)}
            />
          ))}
          <Row label="Slab Tax" value={formatINR(c.slabTax)} bold />
        </div>

        <div className="border-t border-[var(--border-default)] mt-2 pt-2">
          <Row
            label={`Surcharge (${c.surchargeRate}%)`}
            value={formatINR(c.surchargeAmount)}
            section={c.surchargeRate > 0 ? c.surchargeThreshold : undefined}
          />
          {c.surchargeRate > 0 && (
            <div className="flex gap-1 ml-0 mb-1">
              {[10, 15, 25].map((tier) => (
                <span
                  key={tier}
                  className={cn(
                    "text-[9px] px-1.5 py-0.5 rounded-full font-semibold",
                    tier === c.surchargeRate
                      ? "bg-[var(--accent)] text-[var(--accent-text)]"
                      : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]"
                  )}
                >
                  {tier}%
                </span>
              ))}
            </div>
          )}
          <Row label="Health & Education Cess (4%)" value={formatINR(c.healthEducationCess)} />
          <Row label="Gross Tax" value={formatINR(c.grossTax)} bold />
        </div>

        {c.rebate87A > 0 && (
          <div className="border-t border-[var(--border-default)] mt-2 pt-2">
            <Row label="Less: Rebate u/s 87A (≤ ₹7L)" value={`-${formatINR(c.rebate87A)}`} accent="text-[var(--success)]" section="§87A" />
          </div>
        )}

        <div className="border-t-2 border-[var(--text-primary)] mt-2 pt-2">
          <Row label="Total Tax Liability" value={formatINR(c.totalLiability)} bold />
        </div>

        <div className="border-t border-[var(--border-default)] mt-2 pt-2">
          <Row label="TDS Credit" value={`-${formatINR(c.tdsCredit)}`} accent="text-[var(--success)]" />
          <Row label="Advance Tax Paid" value={`-${formatINR(c.advanceTax)}`} accent="text-[var(--success)]" />
          {c.selfAssessmentTax > 0 && (
            <Row label="Self-Assessment Tax" value={`-${formatINR(c.selfAssessmentTax)}`} accent="text-[var(--success)]" />
          )}
        </div>
      </div>
    </div>
  );
}
