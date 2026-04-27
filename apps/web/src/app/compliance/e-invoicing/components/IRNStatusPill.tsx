"use client";

import { cn } from "@/lib/utils";
import type { IRNStatus } from "../types";

interface IRNStatusPillProps {
  status: IRNStatus;
  className?: string;
}

const STATUS_CONFIG: Record<
  IRNStatus,
  { label: string; dot: string; bg: string }
> = {
  GENERATED: {
    label: "Generated",
    dot: "bg-[var(--success)]",
    bg: "bg-[var(--success-muted)] border-[var(--success-border)]",
  },
  CANCELLED: {
    label: "Cancelled",
    dot: "bg-[var(--danger)]",
    bg: "bg-[var(--danger-muted)] border-[var(--danger-border)]",
  },
};

export function IRNStatusPill({ status, className }: IRNStatusPillProps) {
  const cfg = STATUS_CONFIG[status];
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 h-6 px-2.5 rounded-[6px] border",
        "text-[10px] font-semibold",
        cfg.bg,
        className
      )}
      data-testid="irn-status-pill"
    >
      <span
        className={cn("w-1.5 h-1.5 rounded-full flex-shrink-0", cfg.dot)}
      />
      <span className="text-[var(--text-primary)]">{cfg.label}</span>
    </span>
  );
}
