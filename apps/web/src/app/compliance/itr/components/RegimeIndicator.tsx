"use client";

import { cn } from "@/lib/utils";
import type { TaxRegime } from "../types";

interface RegimeIndicatorProps {
  regime: TaxRegime;
  className?: string;
}

export function RegimeIndicator({ regime, className }: RegimeIndicatorProps) {
  const isNew = regime === "NEW";
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-semibold border",
        isNew
          ? "bg-[var(--accent-muted)] text-[var(--accent)] border-[var(--accent)]"
          : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] border-[var(--border-default)]",
        className
      )}
    >
      <span className={cn("w-1.5 h-1.5 rounded-full", isNew ? "bg-[var(--accent)]" : "bg-[var(--text-muted)]")} />
      {isNew ? "New Regime (§202)" : "Old Regime"}
    </span>
  );
}
