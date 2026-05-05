"use client";

import { cn } from "@/lib/utils";
import { CheckCircle2 } from "lucide-react";
import type { GSTR9CMismatch } from "../types";
import { MismatchSeverityBadge } from "./MismatchSeverityBadge";
import { BookVsReturnDeltaCell } from "./BookVsReturnDeltaCell";

interface MismatchRowProps {
  mismatch: GSTR9CMismatch;
  onResolve: (id: string) => void;
  className?: string;
}

export function MismatchRow({ mismatch, onResolve, className }: MismatchRowProps) {
  return (
    <div
      className={cn(
        "flex items-center gap-3 px-4 py-3 rounded-lg border transition-colors",
        mismatch.resolved
          ? "border-[var(--success)]/30 bg-[color-mix(in_srgb,var(--success)_4%,transparent)]"
          : "border-[var(--border-default)] bg-[var(--bg-secondary)] hover:bg-[var(--bg-tertiary)]",
        className,
      )}
      data-testid={`mismatch-row-${mismatch.id}`}
    >
      <MismatchSeverityBadge severity={mismatch.severity} />

      <div className="flex-1 min-w-0">
        <p className="text-xs font-medium text-[var(--text-primary)] truncate">
          Part {mismatch.section} &middot; {mismatch.description}
        </p>
        <div className="mt-1">
          <BookVsReturnDeltaCell
            booksAmount={mismatch.booksAmount}
            returnAmount={mismatch.gstr9Amount}
          />
        </div>
      </div>

      {mismatch.resolved ? (
        <div className="flex items-center gap-1.5 text-[var(--success)]">
          <CheckCircle2 className="w-4 h-4" />
          <span className="text-[10px] font-semibold">Resolved</span>
        </div>
      ) : (
        <button
          onClick={() => onResolve(mismatch.id)}
          className="shrink-0 px-3 py-1.5 rounded-lg text-[10px] font-semibold border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors"
        >
          Resolve
        </button>
      )}
    </div>
  );
}
