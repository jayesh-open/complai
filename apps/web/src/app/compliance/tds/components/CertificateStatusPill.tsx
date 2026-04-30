"use client";

import { cn } from "@complai/ui-components";
import type { CertificateStatus } from "../certificates/types";

const STYLES: Record<CertificateStatus, string> = {
  GENERATED: "bg-[var(--info-muted)] text-[var(--info)]",
  PENDING: "bg-[var(--warning-muted)] text-[var(--warning)]",
  ISSUED: "bg-[var(--success-muted)] text-[var(--success)]",
  REVOKED: "bg-[var(--danger-muted)] text-[var(--danger)]",
};

export function CertificateStatusPill({ status }: { status: CertificateStatus }) {
  return (
    <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold uppercase", STYLES[status])}>
      {status}
    </span>
  );
}
