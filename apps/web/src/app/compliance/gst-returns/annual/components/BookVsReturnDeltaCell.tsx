"use client";

import { cn } from "@/lib/utils";

interface BookVsReturnDeltaCellProps {
  booksAmount: number;
  returnAmount: number;
  className?: string;
}

function fmtINR(amount: number): string {
  if (Math.abs(amount) >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (Math.abs(amount) >= 100_000) return `₹${(amount / 100_000).toFixed(2)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

export function BookVsReturnDeltaCell({ booksAmount, returnAmount, className }: BookVsReturnDeltaCellProps) {
  const delta = booksAmount - returnAmount;
  const sign = delta > 0 ? "+" : delta < 0 ? "-" : "";
  const isZero = Math.abs(delta) < 1;

  return (
    <div className={cn("grid grid-cols-3 gap-2 text-xs tabular-nums", className)} data-testid="delta-cell">
      <span className="text-[var(--text-primary)] font-medium">{fmtINR(booksAmount)}</span>
      <span className="text-[var(--text-primary)] font-medium">{fmtINR(returnAmount)}</span>
      <span
        className={cn(
          "font-semibold",
          isZero && "text-[var(--text-muted)]",
          delta > 0 && "text-[var(--danger)]",
          delta < 0 && "text-[var(--success)]",
        )}
      >
        {isZero ? "—" : `${sign}${fmtINR(Math.abs(delta))}`}
      </span>
    </div>
  );
}
