"use client";

import { cn } from "@/lib/utils";
import { Send, Eye, CheckCircle2, Clock, RotateCcw } from "lucide-react";

export type MagicLinkPillStatus = "SENT" | "VIEWED" | "APPROVED" | "EXPIRED" | "USED";

const STATUS_CONFIG: Record<MagicLinkPillStatus, { label: string; icon: typeof Send; style: string }> = {
  SENT: { label: "Sent", icon: Send, style: "bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-border)]" },
  VIEWED: { label: "Viewed", icon: Eye, style: "bg-[var(--accent-muted)] text-[var(--accent)] border-[var(--accent)]" },
  APPROVED: { label: "Approved", icon: CheckCircle2, style: "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]" },
  EXPIRED: { label: "Expired", icon: Clock, style: "bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning)]" },
  USED: { label: "Used", icon: RotateCcw, style: "bg-[var(--bg-tertiary)] text-[var(--text-muted)] border-[var(--border-default)]" },
};

interface MagicLinkStatusPillProps {
  status: MagicLinkPillStatus;
  className?: string;
}

export function MagicLinkStatusPill({ status, className }: MagicLinkStatusPillProps) {
  const config = STATUS_CONFIG[status];
  const Icon = config.icon;

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-semibold border",
        config.style,
        className
      )}
    >
      <Icon className="w-3 h-3" />
      {config.label}
    </span>
  );
}
