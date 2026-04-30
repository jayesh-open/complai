"use client";

import { cn } from "@complai/ui-components";
import type { ChallanStatus } from "../challans/types";

const STYLES: Record<ChallanStatus, string> = {
  PENDING: "bg-[var(--warning-muted)] text-[var(--warning)]",
  CLEARED: "bg-[var(--success-muted)] text-[var(--success)]",
  REJECTED: "bg-[var(--danger-muted)] text-[var(--danger)]",
};

export function ChallanStatusPill({ status }: { status: ChallanStatus }) {
  return (
    <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold uppercase", STYLES[status])}>
      {status}
    </span>
  );
}
