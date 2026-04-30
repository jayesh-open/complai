"use client";

import { useRouter } from "next/navigation";
import { cn, formatINR } from "@complai/ui-components";
import { ArrowRight, Clock, AlertTriangle } from "lucide-react";
import type { FilingGridCell, QuarterFilingStatus, TDSFormType } from "../types";
import { QUARTERS } from "../types";

interface FilingStatusOverviewProps {
  cells: FilingGridCell[];
  taxYear: string;
}

const STATUS_PILL: Record<QuarterFilingStatus, { bg: string; text: string; label: string }> = {
  NOT_STARTED: { bg: "bg-[var(--bg-tertiary)]", text: "text-[var(--text-muted)]", label: "Not Started" },
  DRAFT: { bg: "bg-[var(--warning-muted)]", text: "text-[var(--warning)]", label: "Draft" },
  SUBMITTED: { bg: "bg-[var(--info-muted)]", text: "text-[var(--info)]", label: "Submitted" },
  FILED: { bg: "bg-[var(--success-muted)]", text: "text-[var(--success)]", label: "Filed" },
  REJECTED: { bg: "bg-[var(--danger-muted)]", text: "text-[var(--danger)]", label: "Rejected" },
};

function daysUntil(dueDate: string): number {
  const [dd, mm, yyyy] = dueDate.split("/").map(Number);
  const due = new Date(yyyy, mm - 1, dd);
  return Math.ceil((due.getTime() - Date.now()) / 86400000);
}

export function FilingStatusOverview({ cells, taxYear }: FilingStatusOverviewProps) {
  const router = useRouter();
  const forms: TDSFormType[] = ["138", "140", "144"];

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-xs">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            <th className="text-left py-2 px-3 text-[var(--text-muted)] font-medium">Form</th>
            {QUARTERS.map((q) => (
              <th key={q.id} className="text-center py-2 px-3 text-[var(--text-muted)] font-medium">
                <div>{q.label}</div>
                <div className="text-[9px] font-normal">{q.dateRange}</div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {forms.map((form) => (
            <tr key={form} className="border-b border-[var(--border-default)]">
              <td className="py-3 px-3">
                <span className="font-semibold text-[var(--text-primary)]">
                  {cells.find((c) => c.formType === form)?.formLabel}
                </span>
              </td>
              {QUARTERS.map((q) => {
                const cell = cells.find((c) => c.formType === form && c.quarter === q.id);
                if (!cell) return <td key={q.id} />;
                const style = STATUS_PILL[cell.status];
                const days = daysUntil(cell.dueDate);
                const overdue = days < 0 && cell.status !== "FILED";
                const urgent = days >= 0 && days <= 7 && cell.status !== "FILED";
                const canStart = cell.status !== "FILED" && cell.status !== "SUBMITTED";

                return (
                  <td key={q.id} className="py-3 px-3 text-center">
                    <div className="space-y-1.5">
                      <span className={cn("inline-block text-[9px] px-1.5 py-0.5 rounded font-semibold", style.bg, style.text)}>
                        {style.label}
                      </span>
                      {cell.entryCount > 0 && (
                        <p className="text-[10px] text-[var(--text-muted)]">
                          {cell.entryCount} entries · {formatINR(cell.totalTds)}
                        </p>
                      )}
                      <div className="flex items-center justify-center gap-0.5">
                        {overdue ? (
                          <AlertTriangle className="w-2.5 h-2.5 text-[var(--danger)]" />
                        ) : (
                          <Clock className={cn("w-2.5 h-2.5", urgent ? "text-[var(--warning)]" : "text-[var(--text-muted)]")} />
                        )}
                        <span className={cn(
                          "text-[9px]",
                          overdue ? "text-[var(--danger)] font-semibold" : urgent ? "text-[var(--warning)]" : "text-[var(--text-muted)]",
                        )}>
                          {overdue ? `${Math.abs(days)}d overdue` : `${days}d left`}
                        </span>
                      </div>
                      {canStart && (
                        <button
                          onClick={() => router.push(`/compliance/tds/file/${form}/${taxYear}/${q.id}`)}
                          className={cn(
                            "inline-flex items-center gap-1 text-[10px] font-semibold px-2 py-0.5 rounded",
                            "text-[var(--accent)] hover:bg-[var(--accent-muted)] transition-colors",
                          )}
                        >
                          {cell.status === "DRAFT" ? "Resume" : "Start Filing"}
                          <ArrowRight className="w-2.5 h-2.5" />
                        </button>
                      )}
                    </div>
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
