"use client";

import Link from "next/link";
import { Truck, Clock, XCircle, Layers } from "lucide-react";
import { cn } from "@/lib/utils";
import type { EwbRecord } from "../../types";

interface EwbActionsProps {
  record: EwbRecord;
  onCancel: () => void;
}

function canCancel(record: EwbRecord): boolean {
  if (record.status !== "ACTIVE") return false;
  const elapsed = Date.now() - new Date(record.generatedAt).getTime();
  return elapsed < 24 * 3600000;
}

function canExtend(record: EwbRecord): boolean {
  if (record.status !== "ACTIVE") return false;
  const total = new Date(record.validUntil).getTime() - new Date(record.validFrom).getTime();
  const elapsed = Date.now() - new Date(record.validFrom).getTime();
  return elapsed > total * 0.5;
}

export function EwbActions({ record, onCancel }: EwbActionsProps) {
  const base = `/compliance/e-way-bill/${record.ewbNumber}`;

  return (
    <div className="flex flex-wrap items-center gap-2">
      {record.status === "ACTIVE" && (
        <Link
          href={`${base}/update-vehicle`}
          className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
          )}
        >
          <Truck className="w-3.5 h-3.5" /> Update Vehicle
        </Link>
      )}
      {canExtend(record) && (
        <Link
          href={`${base}/extend`}
          className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
          )}
        >
          <Clock className="w-3.5 h-3.5" /> Extend Validity
        </Link>
      )}
      {canCancel(record) && (
        <button
          onClick={onCancel}
          className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--danger-border)]",
            "text-[var(--danger)] hover:bg-[var(--danger-muted)] transition-colors",
          )}
        >
          <XCircle className="w-3.5 h-3.5" /> Cancel EWB
        </button>
      )}
      {record.status === "ACTIVE" && !record.consolidatedEwbNo && (
        <Link
          href="/compliance/e-way-bill/consolidate"
          className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
          )}
        >
          <Layers className="w-3.5 h-3.5" /> Add to Consolidation
        </Link>
      )}
    </div>
  );
}
