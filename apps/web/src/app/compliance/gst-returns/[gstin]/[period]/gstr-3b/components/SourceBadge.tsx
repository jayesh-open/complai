"use client";

import { cn } from "@/lib/utils";
import type { DataSource } from "./types";

const BADGE_CONFIG: Record<DataSource, { label: string; className: string }> = {
  gstr1: { label: "From GSTR-1", className: "bg-[var(--success-muted)] text-[var(--success)]" },
  gstr2b: { label: "From 2B", className: "bg-[var(--accent-muted)] text-[var(--accent)]" },
  computed: { label: "Computed", className: "bg-[var(--warning-muted)] text-[var(--warning)]" },
  override: { label: "Override", className: "bg-[var(--danger)]/10 text-[var(--danger)]" },
};

export function SourceBadge({ source }: { source: DataSource }) {
  const config = BADGE_CONFIG[source];
  return (
    <span className={cn("px-1.5 py-0.5 rounded text-[9px] font-bold uppercase whitespace-nowrap", config.className)}>
      {config.label}
    </span>
  );
}
