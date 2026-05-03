"use client";

import { cn } from "@/lib/utils";
import { Users } from "lucide-react";
import { ITRFormRecommendationPill } from "../../../components/ITRFormRecommendationPill";
import { RegimeIndicator } from "../../../components/RegimeIndicator";
import type { ITREmployee } from "../../../types";

interface SelectStepProps {
  employees: ITREmployee[];
  selected: Set<string>;
  onToggle: (id: string) => void;
  onToggleAll: () => void;
  onNext: () => void;
}

export function SelectStep({ employees, selected, onToggle, onToggleAll, onNext }: SelectStepProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-xs text-[var(--text-muted)]">{selected.size} of {employees.length} employees selected</p>
        <button onClick={onToggleAll} className="text-xs text-[var(--accent)] hover:underline font-medium">
          {selected.size === employees.length ? "Deselect All" : "Select All"}
        </button>
      </div>
      <div className="bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden max-h-[480px] overflow-y-auto">
        <table className="w-full">
          <thead className="sticky top-0 bg-[var(--bg-secondary)] z-10">
            <tr className="border-b border-[var(--border-default)]">
              <th className="px-[18px] py-[10px] w-10"></th>
              <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Employee</th>
              <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">PAN</th>
              <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Form</th>
              <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Regime</th>
            </tr>
          </thead>
          <tbody>
            {employees.map((emp) => (
              <tr key={emp.id} className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] cursor-pointer" onClick={() => onToggle(emp.id)}>
                <td className="px-[18px] py-3">
                  <input type="checkbox" checked={selected.has(emp.id)} onChange={() => onToggle(emp.id)} className="rounded border-[var(--border-default)]" />
                </td>
                <td className="px-[18px] py-3">
                  <div className="text-xs font-medium text-[var(--text-primary)]">{emp.name}</div>
                  <div className="text-[10px] text-[var(--text-muted)]">{emp.department}</div>
                </td>
                <td className="px-[18px] py-3 text-xs font-mono text-[var(--text-secondary)]">{emp.pan}</td>
                <td className="px-[18px] py-3"><ITRFormRecommendationPill form={emp.recommendedForm} /></td>
                <td className="px-[18px] py-3"><RegimeIndicator regime={emp.regime} /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="flex justify-end">
        <button
          disabled={selected.size === 0}
          onClick={onNext}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
            selected.size > 0 ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] cursor-not-allowed"
          )}
        >
          <Users className="w-3.5 h-3.5" />
          Continue with {selected.size} employees
        </button>
      </div>
    </div>
  );
}
