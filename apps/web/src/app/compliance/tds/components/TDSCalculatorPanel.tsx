"use client";

import { cn, formatINR } from "@complai/ui-components";
import { getPaymentCodeInfo } from "../payment-codes";
import type { TDSCalculationResult } from "../types";

interface TDSCalculatorPanelProps {
  paymentCode: string;
  amount: number;
  isNonResident: boolean;
  noPan: boolean;
  className?: string;
}

function calculateTDS(
  paymentCode: string,
  amount: number,
  isNonResident: boolean,
  noPan: boolean
): TDSCalculationResult | null {
  const info = getPaymentCodeInfo(paymentCode);
  if (!info || amount <= 0) return null;

  const thresholds: Record<string, number> = {
    "1009": 50000,
    "1023": 30000,
    "1024": 30000,
    "1026": 15000,
    "1027": 30000,
    "1028": 30000,
    "1031": 5000000,
  };

  const threshold = thresholds[paymentCode] ?? 0;
  const thresholdMet = threshold === 0 || amount >= threshold;
  let baseRate = info.baseRate;

  if (noPan) {
    baseRate = paymentCode === "1031" || paymentCode === "1035" ? 5 : 20;
  }

  const cessRate = isNonResident ? 4 : 0;
  const surchargeRate = amount > 5000000 ? 10 : 0;
  const tdsAmount = Math.round((amount * baseRate) / 100);
  const surcharge = Math.round((tdsAmount * surchargeRate) / 100);
  const cess = Math.round(((tdsAmount + surcharge) * cessRate) / 100);
  const effectiveRate =
    baseRate +
    (baseRate * surchargeRate) / 100 +
    ((baseRate + (baseRate * surchargeRate) / 100) * cessRate) / 100;

  return {
    section: info.section,
    subClause: info.subClause,
    paymentCode: info.code,
    baseRate,
    cessRate,
    surchargeRate,
    effectiveRate: Math.round(effectiveRate * 100) / 100,
    tdsAmount,
    surcharge,
    cess,
    totalTax: tdsAmount + surcharge + cess,
    thresholdMet,
    thresholdAmount: threshold,
    noPanApplied: noPan,
  };
}

export function TDSCalculatorPanel({
  paymentCode,
  amount,
  isNonResident,
  noPan,
  className,
}: TDSCalculatorPanelProps) {
  const result = calculateTDS(paymentCode, amount, isNonResident, noPan);
  const info = getPaymentCodeInfo(paymentCode);

  if (!info) {
    return (
      <div
        className={cn(
          "rounded-xl border border-[var(--border-default)] bg-[var(--bg-secondary)] p-6",
          "text-center text-sm text-[var(--text-muted)]",
          className
        )}
      >
        Select a payment code and enter amount to see TDS calculation
      </div>
    );
  }

  return (
    <div
      className={cn(
        "rounded-xl border border-[var(--border-default)] bg-[var(--bg-secondary)] p-6 space-y-4",
        className
      )}
    >
      <div className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide">
        Live Calculation Preview
      </div>

      <div className="space-y-2">
        <CalcRow
          label="Section Reference"
          value={`${info.section}${info.subClause ? `[${info.subClause}]` : ""}`}
          mono
        />
        <CalcRow label="Payment Code" value={info.code} mono />
        <CalcRow label="Description" value={info.description} />
      </div>

      {result && (
        <>
          <div className="border-t border-[var(--border-default)] pt-3 space-y-2">
            <CalcRow label="Base Rate" value={`${result.baseRate}%`} />
            {result.cessRate > 0 && (
              <CalcRow label="Cess Rate (H&E)" value={`${result.cessRate}%`} />
            )}
            {result.surchargeRate > 0 && (
              <CalcRow label="Surcharge Rate" value={`${result.surchargeRate}%`} />
            )}
            <CalcRow
              label="Effective Rate"
              value={`${result.effectiveRate}%`}
              highlight
            />
          </div>

          <div className="border-t border-[var(--border-default)] pt-3 space-y-2">
            <CalcRow label="TDS Amount" value={formatINR(result.tdsAmount)} />
            {result.surcharge > 0 && (
              <CalcRow label="Surcharge" value={formatINR(result.surcharge)} />
            )}
            {result.cess > 0 && (
              <CalcRow label="Cess" value={formatINR(result.cess)} />
            )}
            <CalcRow
              label="Total Deduction"
              value={formatINR(result.totalTax)}
              highlight
            />
          </div>

          <div className="border-t border-[var(--border-default)] pt-3">
            {result.thresholdAmount > 0 && (
              <div
                className={cn(
                  "text-[10px] font-medium px-2 py-1 rounded",
                  result.thresholdMet
                    ? "bg-[var(--success-muted)] text-[var(--success)]"
                    : "bg-[var(--warning-muted)] text-[var(--warning)]"
                )}
              >
                Threshold: {formatINR(result.thresholdAmount)} —{" "}
                {result.thresholdMet ? "Met" : "Not met (no TDS required)"}
              </div>
            )}
            {result.noPanApplied && (
              <div className="text-[10px] font-medium px-2 py-1 mt-1 rounded bg-[var(--danger-muted)] text-[var(--danger)]">
                Section 397(2) — No-PAN rate applied ({result.baseRate}%)
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}

function CalcRow({
  label,
  value,
  mono,
  highlight,
}: {
  label: string;
  value: string;
  mono?: boolean;
  highlight?: boolean;
}) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-[11px] text-[var(--text-muted)]">{label}</span>
      <span
        className={cn(
          "text-xs",
          mono && "font-mono",
          highlight
            ? "text-[var(--text-primary)] font-bold"
            : "text-[var(--text-secondary)] font-medium"
        )}
      >
        {value}
      </span>
    </div>
  );
}
