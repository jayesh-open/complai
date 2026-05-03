"use client";

import { cn } from "@/lib/utils";
import { ArrowDownCircle, ArrowUpCircle } from "lucide-react";

function formatINR(amount: number): string {
  const abs = Math.abs(amount);
  if (abs >= 10_000_000) return `₹${(abs / 10_000_000).toFixed(2)} Cr`;
  if (abs >= 100_000) return `₹${(abs / 100_000).toFixed(1)} L`;
  return `₹${abs.toLocaleString("en-IN")}`;
}

interface RefundIndicatorProps {
  amount: number;
  className?: string;
}

export function RefundIndicator({ amount, className }: RefundIndicatorProps) {
  const isRefund = amount <= 0;
  const displayAmount = Math.abs(amount);

  return (
    <div
      className={cn(
        "flex items-center justify-between rounded-xl p-5 border-2",
        isRefund
          ? "bg-[var(--success-muted)] border-[var(--success)]"
          : "bg-[var(--danger-muted)] border-[var(--danger)]",
        className
      )}
    >
      <div className="flex items-center gap-3">
        {isRefund ? (
          <ArrowDownCircle className="w-8 h-8 text-[var(--success)]" />
        ) : (
          <ArrowUpCircle className="w-8 h-8 text-[var(--danger)]" />
        )}
        <div>
          <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">
            {isRefund ? "Refund Due" : "Balance Tax Payable"}
          </p>
          <p
            className={cn(
              "text-2xl font-bold tabular-nums",
              isRefund ? "text-[var(--success)]" : "text-[var(--danger)]"
            )}
          >
            {formatINR(displayAmount)}
          </p>
        </div>
      </div>
      <span
        className={cn(
          "text-xs font-semibold px-3 py-1.5 rounded-lg",
          isRefund
            ? "bg-[var(--success)] text-white"
            : "bg-[var(--danger)] text-white"
        )}
      >
        {isRefund ? "REFUND" : "PAYABLE"}
      </span>
    </div>
  );
}
