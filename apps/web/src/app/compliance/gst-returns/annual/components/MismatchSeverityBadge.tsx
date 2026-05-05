"use client";

import { cn } from "@/lib/utils";
import type { MismatchSeverity } from "../types";

interface MismatchSeverityBadgeProps {
  severity: MismatchSeverity;
  className?: string;
}

const SEVERITY_STYLES: Record<MismatchSeverity, string> = {
  INFO: "bg-[color-mix(in_srgb,var(--accent)_12%,transparent)] text-[var(--accent)] border-[var(--accent)]",
  WARN: "bg-[color-mix(in_srgb,var(--warning)_12%,transparent)] text-[var(--warning)] border-[var(--warning)]",
  ERROR: "bg-[color-mix(in_srgb,var(--danger)_12%,transparent)] text-[var(--danger)] border-[var(--danger)]",
};

export function MismatchSeverityBadge({ severity, className }: MismatchSeverityBadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wide border",
        SEVERITY_STYLES[severity],
        className,
      )}
      data-testid={`severity-badge-${severity.toLowerCase()}`}
    >
      {severity}
    </span>
  );
}
