"use client";

import { ShieldCheck, ShieldX } from "lucide-react";
import { cn } from "@complai/ui-components";

interface DTAAEvidenceBadgeProps {
  form41Filed: boolean;
  trcAttached: boolean;
  className?: string;
}

export function DTAAEvidenceBadge({ form41Filed, trcAttached, className }: DTAAEvidenceBadgeProps) {
  const allClear = form41Filed && trcAttached;
  return (
    <div className={cn("inline-flex items-center gap-2", className)}>
      <div
        className={cn(
          "inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-semibold border",
          form41Filed
            ? "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]"
            : "bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-border)]"
        )}
      >
        {form41Filed ? <ShieldCheck className="w-3 h-3" /> : <ShieldX className="w-3 h-3" />}
        Form 41
      </div>
      <div
        className={cn(
          "inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-semibold border",
          trcAttached
            ? "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]"
            : "bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-border)]"
        )}
      >
        {trcAttached ? <ShieldCheck className="w-3 h-3" /> : <ShieldX className="w-3 h-3" />}
        TRC
      </div>
      {allClear && (
        <span className="text-[10px] text-[var(--success)] font-medium">
          DTAA eligible
        </span>
      )}
    </div>
  );
}
