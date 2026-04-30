"use client";

import { cn } from "@complai/ui-components";
import type { EntryStatus, FilingStatus } from "../types";

const ENTRY_STYLES: Record<EntryStatus, string> = {
  PENDING: "bg-[var(--warning-muted)] text-[var(--warning)]",
  DEPOSITED: "bg-[var(--info-muted)] text-[var(--info)]",
  FILED: "bg-[var(--success-muted)] text-[var(--success)]",
  REVISED: "bg-[var(--purple-muted)] text-[var(--purple)]",
};

const FILING_STYLES: Record<FilingStatus, string> = {
  NOT_DUE: "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
  PENDING: "bg-[var(--warning-muted)] text-[var(--warning)]",
  FILED: "bg-[var(--success-muted)] text-[var(--success)]",
  OVERDUE: "bg-[var(--danger-muted)] text-[var(--danger)]",
  REVISED: "bg-[var(--purple-muted)] text-[var(--purple)]",
};

export function EntryStatusPill({ status }: { status: EntryStatus }) {
  return (
    <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold uppercase", ENTRY_STYLES[status])}>
      {status}
    </span>
  );
}

export function FilingStatusPill({ status }: { status: FilingStatus }) {
  const label = status === "NOT_DUE" ? "Not Due" : status;
  return (
    <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold uppercase", FILING_STYLES[status])}>
      {label}
    </span>
  );
}
