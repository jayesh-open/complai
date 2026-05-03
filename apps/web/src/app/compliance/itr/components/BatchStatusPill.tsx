"use client";

import { cn } from "@/lib/utils";
import type { BatchStatus } from "../types";

const STATUS_CONFIG: Record<BatchStatus, { label: string; style: string }> = {
  DRAFT: { label: "Draft", style: "bg-[var(--bg-tertiary)] text-[var(--text-muted)] border-[var(--border-default)]" },
  IN_PROGRESS: { label: "In Progress", style: "bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning-muted)]" },
  COMPLETED: { label: "Completed", style: "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]" },
  FAILED: { label: "Failed", style: "bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-muted)]" },
};

interface BatchStatusPillProps {
  status: BatchStatus;
  className?: string;
}

export function BatchStatusPill({ status, className }: BatchStatusPillProps) {
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
