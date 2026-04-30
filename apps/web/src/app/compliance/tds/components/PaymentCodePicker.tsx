"use client";

import { cn } from "@complai/ui-components";
import { PAYMENT_CODES } from "../payment-codes";
import type { Section2025 } from "../types";

interface PaymentCodePickerProps {
  value: string;
  onChange: (code: string) => void;
  sectionFilter?: Section2025;
  className?: string;
}

export function PaymentCodePicker({
  value,
  onChange,
  sectionFilter,
  className,
}: PaymentCodePickerProps) {
  const codes = sectionFilter
    ? PAYMENT_CODES.filter((pc) => pc.section === sectionFilter)
    : PAYMENT_CODES;

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={cn(
        "w-full px-3 py-2 rounded-lg text-xs",
        "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
        "text-[var(--text-primary)]",
        "focus:outline-none focus:border-[var(--accent)]",
        "focus:ring-2 focus:ring-[var(--accent-muted)]",
        className
      )}
    >
      <option value="">Select payment code...</option>
      {codes.map((pc) => (
        <option key={pc.code} value={pc.code}>
          {pc.code} — {pc.label} ({pc.section}) @ {pc.baseRate}%
        </option>
      ))}
    </select>
  );
}
