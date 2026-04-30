"use client";

import { cn } from "@complai/ui-components";
import { Clock, AlertTriangle } from "lucide-react";
import type { QuarterFilingStatus } from "../types";
import { QUARTERS } from "../types";

interface QuarterPickerProps {
  value: string;
  onChange: (quarter: string) => void;
  statuses?: Record<string, QuarterFilingStatus>;
}

const STATUS_STYLES: Record<QuarterFilingStatus, { bg: string; text: string; label: string }> = {
  NOT_STARTED: { bg: "bg-[var(--bg-tertiary)]", text: "text-[var(--text-muted)]", label: "Not Started" },
  DRAFT: { bg: "bg-[var(--warning-muted)]", text: "text-[var(--warning)]", label: "Draft" },
  SUBMITTED: { bg: "bg-[var(--info-muted)]", text: "text-[var(--info)]", label: "Submitted" },
  FILED: { bg: "bg-[var(--success-muted)]", text: "text-[var(--success)]", label: "Filed" },
  REJECTED: { bg: "bg-[var(--danger-muted)]", text: "text-[var(--danger)]", label: "Rejected" },
};

function daysUntil(date: Date): number {
  const now = new Date();
  return Math.ceil((date.getTime() - now.getTime()) / 86400000);
}

export function QuarterPicker({ value, onChange, statuses }: QuarterPickerProps) {
  return (
    <div className="grid grid-cols-4 gap-3">
      {QUARTERS.map((q) => {
        const selected = value === q.id;
        const status = statuses?.[q.id] ?? "NOT_STARTED";
        const style = STATUS_STYLES[status];
        const days = daysUntil(q.dueDateObj);
        const overdue = days < 0;
        const urgent = days >= 0 && days <= 7;

        return (
          <button
            key={q.id}
            onClick={() => onChange(q.id)}
            className={cn(
              "rounded-lg border p-3 text-left transition-all",
              selected
                ? "border-[var(--accent)] ring-2 ring-[var(--accent)]/30 bg-[var(--accent-muted)]"
                : "border-[var(--border-default)] bg-[var(--bg-secondary)] hover:border-[var(--border-hover)]",
            )}
          >
            <div className="flex items-center justify-between mb-1">
              <span className="text-sm font-bold text-[var(--text-primary)]">{q.label}</span>
              <span className={cn("text-[9px] px-1.5 py-0.5 rounded font-semibold", style.bg, style.text)}>
                {style.label}
              </span>
            </div>
            <p className="text-[10px] text-[var(--text-muted)]">{q.dateRange}</p>
            <div className="flex items-center gap-1 mt-2">
              {overdue ? (
                <AlertTriangle className="w-3 h-3 text-[var(--danger)]" />
              ) : (
                <Clock className={cn("w-3 h-3", urgent ? "text-[var(--warning)]" : "text-[var(--text-muted)]")} />
              )}
              <span className={cn(
                "text-[10px] font-medium",
                overdue ? "text-[var(--danger)]" : urgent ? "text-[var(--warning)]" : "text-[var(--text-muted)]",
              )}>
                {overdue ? `Overdue by ${Math.abs(days)}d` : `Due: ${q.dueDate} (${days}d)`}
              </span>
            </div>
          </button>
        );
      })}
    </div>
  );
}
