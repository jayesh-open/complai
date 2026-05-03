"use client";

import { cn } from "@/lib/utils";
import type { FilingStatus } from "../types";

const STATUS_CONFIG: Record<FilingStatus, { label: string; style: string }> = {
  NOT_STARTED: { label: "Not Started", style: "bg-[var(--bg-tertiary)] text-[var(--text-muted)] border-[var(--border-default)]" },
  AIS_FETCHED: { label: "AIS Fetched", style: "bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-border)]" },
  FORM_GENERATED: { label: "Form Ready", style: "bg-[var(--purple-muted)] text-[var(--purple)] border-[var(--purple-muted)]" },
  REVIEW_PENDING: { label: "Review Pending", style: "bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning-muted)]" },
  EMPLOYEE_APPROVED: { label: "Approved", style: "bg-[var(--accent-muted)] text-[var(--accent)] border-[var(--accent)]" },
  FILED: { label: "Filed", style: "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]" },
  ACKNOWLEDGED: { label: "Acknowledged", style: "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]" },
  DEFECTIVE: { label: "Defective", style: "bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-muted)]" },
};

interface FilingStatusPillProps {
  status: FilingStatus;
  className?: string;
}

export function FilingStatusPill({ status, className }: FilingStatusPillProps) {
  const config = STATUS_CONFIG[status];
  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border",
        config.style,
        className
      )}
    >
      {config.label}
    </span>
  );
}
