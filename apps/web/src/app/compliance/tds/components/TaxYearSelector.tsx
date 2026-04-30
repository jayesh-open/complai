"use client";

import { cn } from "@complai/ui-components";

interface TaxYearSelectorProps {
  value: string;
  onChange: (year: string) => void;
  className?: string;
}

const TAX_YEARS = ["2026-27", "2027-28", "2028-29"];

export function TaxYearSelector({ value, onChange, className }: TaxYearSelectorProps) {
  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={cn(
        "px-3 py-2 rounded-lg text-xs font-medium",
        "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
        "text-[var(--text-primary)]",
        "focus:outline-none focus:border-[var(--accent)]",
        "focus:ring-2 focus:ring-[var(--accent-muted)]",
        className
      )}
    >
      {TAX_YEARS.map((ty) => (
        <option key={ty} value={ty}>
          Tax Year {ty}
        </option>
      ))}
    </select>
  );
}
