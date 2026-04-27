"use client";

import { cn } from "@/lib/utils";
import type { EwbStatus } from "../types";

interface EwbStatusPillProps {
  status: EwbStatus;
  nearingExpiry?: boolean;
  className?: string;
}

const config: Record<EwbStatus, { label: string; dot: string; text: string }> = {
  ACTIVE: {
    label: "Active",
    dot: "bg-[var(--success)]",
    text: "text-[var(--success)]",
  },
  EXPIRED: {
    label: "Expired",
    dot: "bg-[var(--text-muted)]",
    text: "text-[var(--text-muted)]",
  },
  CANCELLED: {
    label: "Cancelled",
    dot: "bg-[var(--danger)]",
    text: "text-[var(--danger)]",
  },
  CONSOLIDATED: {
    label: "Consolidated",
    dot: "bg-[var(--info)]",
    text: "text-[var(--info)]",
  },
};

export function EwbStatusPill({ status, nearingExpiry, className }: EwbStatusPillProps) {
  const c = nearingExpiry && status === "ACTIVE"
    ? { label: "Expiring Soon", dot: "bg-[var(--warning)]", text: "text-[var(--warning)]" }
    : config[status];

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 h-5 px-2 rounded-md text-[10px] font-semibold",
        "bg-[var(--bg-tertiary)] border border-[var(--border-default)]",
        c.text,
        className,
      )}
    >
      <span className={cn("w-1.5 h-1.5 rounded-full flex-shrink-0", c.dot)} />
      {c.label}
    </span>
  );
}
