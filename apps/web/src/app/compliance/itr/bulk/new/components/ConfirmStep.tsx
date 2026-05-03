"use client";

import { cn } from "@/lib/utils";
import { Send } from "lucide-react";
import type { ITREmployee } from "../../../types";

interface ConfirmStepProps {
  selected: Set<string>;
  employees: ITREmployee[];
  batchName: string;
  autoFetchAIS: boolean;
  sendMagicLinks: boolean;
  onBack: () => void;
  onSubmit: () => void;
}

export function ConfirmStep({ selected, employees, batchName, autoFetchAIS, sendMagicLinks, onBack, onSubmit }: ConfirmStepProps) {
  const selectedEmployees = employees.filter((e) => selected.has(e.id));
  return (
    <div className="space-y-6">
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6 space-y-4">
        <h3 className="text-sm font-semibold text-[var(--text-primary)]">Batch Summary</h3>
        <div className="grid grid-cols-2 gap-4 text-xs">
          <div>
            <span className="text-[var(--text-muted)]">Batch Name:</span>
            <span className="ml-2 text-[var(--text-primary)] font-medium">{batchName}</span>
          </div>
          <div>
            <span className="text-[var(--text-muted)]">Employees:</span>
            <span className="ml-2 text-[var(--text-primary)] font-medium">{selected.size}</span>
          </div>
          <div>
            <span className="text-[var(--text-muted)]">Auto-fetch AIS:</span>
            <span className="ml-2 text-[var(--text-primary)] font-medium">{autoFetchAIS ? "Yes" : "No"}</span>
          </div>
          <div>
            <span className="text-[var(--text-muted)]">Magic Links:</span>
            <span className="ml-2 text-[var(--text-primary)] font-medium">{sendMagicLinks ? "Yes" : "No"}</span>
          </div>
        </div>
        <div className="border-t border-[var(--border-default)] pt-3">
          <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide mb-2">Selected Employees</p>
          <div className="flex flex-wrap gap-2">
            {selectedEmployees.slice(0, 10).map((emp) => (
              <span key={emp.id} className="inline-flex items-center px-2 py-1 rounded-md text-[10px] bg-[var(--bg-tertiary)] text-[var(--text-secondary)] border border-[var(--border-default)]">
                {emp.name}
              </span>
            ))}
            {selectedEmployees.length > 10 && (
              <span className="inline-flex items-center px-2 py-1 rounded-md text-[10px] text-[var(--text-muted)]">
                +{selectedEmployees.length - 10} more
              </span>
            )}
          </div>
        </div>
      </div>
      <div className="flex items-center justify-between">
        <button onClick={onBack} className="text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] font-medium">← Back</button>
        <button
          onClick={onSubmit}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
            "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors"
          )}
        >
          <Send className="w-3.5 h-3.5" />
          Create Batch & Start Filing
        </button>
      </div>
    </div>
  );
}
