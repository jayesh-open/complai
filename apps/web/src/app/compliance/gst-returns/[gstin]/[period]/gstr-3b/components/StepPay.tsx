"use client";

import { ArrowRight, Wallet, CreditCard } from "lucide-react";
import { cn, formatINR } from "@/lib/utils";
import type { GSTR3BData } from "./types";

interface StepPayProps {
  data: GSTR3BData;
}

export function StepPay({ data }: StepPayProps) {
  return (
    <div className="p-6 space-y-4" data-testid="step-pay">
      <div>
        <h2 className="text-heading-lg text-foreground">Step 3: Payment & Offset</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Review ITC utilisation and cash payment required.
        </p>
      </div>

      {/* Ledger Balances */}
      <div className="grid grid-cols-2 gap-4">
        <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
          <div className="bg-[var(--bg-tertiary)] px-4 py-2.5 border-b border-[var(--border-default)] flex items-center gap-2">
            <CreditCard className="w-4 h-4 text-[var(--accent)]" />
            <span className="text-xs font-bold text-foreground">Credit Ledger (ITC Available)</span>
          </div>
          <div className="p-4 space-y-2">
            {[
              { label: "CGST", value: data.ledger.creditCgst },
              { label: "SGST", value: data.ledger.creditSgst },
              { label: "IGST", value: data.ledger.creditIgst },
            ].map((item) => (
              <div key={item.label} className="flex justify-between text-xs">
                <span className="text-foreground-muted">{item.label}</span>
                <span className="font-mono font-medium text-[var(--success)]">{formatINR(item.value)}</span>
              </div>
            ))}
            <div className="pt-2 border-t border-[var(--border-default)] flex justify-between text-xs">
              <span className="font-bold text-foreground">Total Credit</span>
              <span className="font-mono font-bold text-[var(--success)]">
                {formatINR(data.ledger.creditCgst + data.ledger.creditSgst + data.ledger.creditIgst)}
              </span>
            </div>
          </div>
        </div>

        <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
          <div className="bg-[var(--bg-tertiary)] px-4 py-2.5 border-b border-[var(--border-default)] flex items-center gap-2">
            <Wallet className="w-4 h-4 text-[var(--accent)]" />
            <span className="text-xs font-bold text-foreground">Cash Ledger (Available Balance)</span>
          </div>
          <div className="p-4 space-y-2">
            {[
              { label: "CGST", value: data.ledger.cashCgst },
              { label: "SGST", value: data.ledger.cashSgst },
              { label: "IGST", value: data.ledger.cashIgst },
            ].map((item) => (
              <div key={item.label} className="flex justify-between text-xs">
                <span className="text-foreground-muted">{item.label}</span>
                <span className="font-mono font-medium text-foreground">{formatINR(item.value)}</span>
              </div>
            ))}
            <div className="pt-2 border-t border-[var(--border-default)] flex justify-between text-xs">
              <span className="font-bold text-foreground">Total Cash</span>
              <span className="font-mono font-bold text-foreground">
                {formatINR(data.ledger.cashCgst + data.ledger.cashSgst + data.ledger.cashIgst)}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Offset Wizard */}
      <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
        <div className="bg-[var(--bg-tertiary)] px-4 py-2.5 border-b border-[var(--border-default)]">
          <span className="text-xs font-bold text-foreground">ITC Offset Order (as per Section 49A)</span>
        </div>
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-[var(--border-default)] bg-[var(--bg-secondary)]">
              <th className="text-left px-4 py-2 text-[var(--text-muted)] font-semibold">Tax Head</th>
              <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Liability</th>
              <th className="text-center px-4 py-2 text-[var(--text-muted)] font-semibold" />
              <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Credit Used</th>
              <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Cash Required</th>
              <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Remaining</th>
            </tr>
          </thead>
          <tbody>
            {data.offsets.map((row) => (
              <tr key={row.head} className="border-b border-[var(--border-default)] last:border-0">
                <td className="px-4 py-2.5 font-semibold text-foreground">{row.head}</td>
                <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(row.liability)}</td>
                <td className="px-4 py-2.5 text-center">
                  <ArrowRight className="w-3 h-3 text-[var(--text-muted)] inline" />
                </td>
                <td className="px-4 py-2.5 text-right font-mono text-[var(--success)]">{formatINR(row.creditUsed)}</td>
                <td className="px-4 py-2.5 text-right font-mono text-[var(--danger)]">{formatINR(row.cashUsed)}</td>
                <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(row.remaining)}</td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr className="bg-[var(--bg-tertiary)] border-t-2 border-[var(--accent)]">
              <td className="px-4 py-2.5 font-bold text-foreground">Total</td>
              <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                {formatINR(data.offsets.reduce((s, r) => s + r.liability, 0))}
              </td>
              <td />
              <td className="px-4 py-2.5 text-right font-mono font-bold text-[var(--success)]">
                {formatINR(data.offsets.reduce((s, r) => s + r.creditUsed, 0))}
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold text-[var(--danger)]">
                {formatINR(data.offsets.reduce((s, r) => s + r.cashUsed, 0))}
              </td>
              <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                {formatINR(data.offsets.reduce((s, r) => s + r.remaining, 0))}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      {/* Interest & Late Fee */}
      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="flex justify-between text-xs">
            <span className="text-foreground-muted">Interest (if any)</span>
            <span className="font-mono text-foreground">{formatINR(data.interestLateFee.interest)}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="text-foreground-muted">Late Fee (if any)</span>
            <span className="font-mono text-foreground">{formatINR(data.interestLateFee.lateFee)}</span>
          </div>
        </div>
      </div>

      {/* Net Payable Summary */}
      <div className="bg-[var(--danger)]/5 border border-[var(--danger)]/20 rounded-lg p-4 text-center">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--danger)] mb-1">Total Cash Payment Required</div>
        <div className="text-2xl font-bold font-mono text-[var(--danger)]">
          {formatINR(data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst)}
        </div>
        <div className="text-[10px] text-foreground-muted mt-1">
          CGST: {formatINR(data.netPayable.cgst)} + SGST: {formatINR(data.netPayable.sgst)} + IGST: {formatINR(data.netPayable.igst)}
        </div>
      </div>
    </div>
  );
}
