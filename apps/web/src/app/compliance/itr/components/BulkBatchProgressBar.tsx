"use client";

import { cn } from "@/lib/utils";

interface BulkBatchProgressBarProps {
  filed: number;
  pending: number;
  failed: number;
  total: number;
  className?: string;
}

export function BulkBatchProgressBar({ filed, pending, failed, total, className }: BulkBatchProgressBarProps) {
  const filedPct = total > 0 ? (filed / total) * 100 : 0;
  const failedPct = total > 0 ? (failed / total) * 100 : 0;

  return (
    <div className={cn("w-full", className)}>
      <div className="h-2 bg-[var(--bg-tertiary)] rounded-full overflow-hidden flex">
        {filedPct > 0 && (
          <div
            className="h-full bg-[var(--success)] transition-all duration-300"
            style={{ width: `${filedPct}%` }}
          />
        )}
        {failedPct > 0 && (
          <div
            className="h-full bg-[var(--danger)] transition-all duration-300"
            style={{ width: `${failedPct}%` }}
          />
        )}
      </div>
      <div className="flex items-center gap-3 mt-1.5">
        <span className="text-[10px] text-[var(--success)] font-medium">{filed} filed</span>
        {pending > 0 && <span className="text-[10px] text-[var(--text-muted)] font-medium">{pending} pending</span>}
        {failed > 0 && <span className="text-[10px] text-[var(--danger)] font-medium">{failed} failed</span>}
      </div>
    </div>
  );
}
