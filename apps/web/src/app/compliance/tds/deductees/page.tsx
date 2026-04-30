"use client";

import { useState, useMemo } from "react";
import { Search, Download } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_DEDUCTEES } from "../mock-data";
import { DeducteeTable } from "../components/DeducteeTable";
import { TaxYearSelector } from "../components/TaxYearSelector";
import type { Section2025, ResidencyStatus } from "../types";

type SectionFilter = "ALL" | Section2025;
type ResidencyFilter = "ALL" | ResidencyStatus;

export default function DeducteesListPage() {
  const [search, setSearch] = useState("");
  const [sectionFilter, setSectionFilter] = useState<SectionFilter>("ALL");
  const [residencyFilter, setResidencyFilter] = useState<ResidencyFilter>("ALL");
  const [taxYear, setTaxYear] = useState("2026-27");

  const filtered = useMemo(() => {
    return ALL_DEDUCTEES.filter((d) => {
      if (sectionFilter !== "ALL" && d.sectionPreference !== sectionFilter) return false;
      if (residencyFilter !== "ALL" && d.residency !== residencyFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return d.pan.toLowerCase().includes(q) || d.name.toLowerCase().includes(q);
      }
      return true;
    });
  }, [search, sectionFilter, residencyFilter]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">Deductees</h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            All TDS deductees for Tax Year {taxYear}
          </p>
        </div>
        <TaxYearSelector value={taxYear} onChange={setTaxYear} />
      </div>

      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by PAN or name..."
            className={cn(
              "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)]",
              "focus:ring-2 focus:ring-[var(--accent-muted)]"
            )}
          />
        </div>
        <select
          value={sectionFilter}
          onChange={(e) => setSectionFilter(e.target.value as SectionFilter)}
          className={cn(
            "px-3 py-2 rounded-lg text-xs",
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "text-[var(--text-primary)]",
            "focus:outline-none focus:border-[var(--accent)]"
          )}
        >
          <option value="ALL">All Sections</option>
          <option value="392">392 — Salary</option>
          <option value="393(1)">393(1) — Resident</option>
          <option value="393(2)">393(2) — Non-Resident</option>
          <option value="393(3)">393(3) — TCS</option>
        </select>
        <select
          value={residencyFilter}
          onChange={(e) => setResidencyFilter(e.target.value as ResidencyFilter)}
          className={cn(
            "px-3 py-2 rounded-lg text-xs",
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "text-[var(--text-primary)]",
            "focus:outline-none focus:border-[var(--accent)]"
          )}
        >
          <option value="ALL">Any Residency</option>
          <option value="RESIDENT">Resident</option>
          <option value="NON_RESIDENT">Non-Resident</option>
          <option value="ANY_PERSON">Any Person</option>
        </select>
        <button
          className={cn(
            "flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
            "transition-colors"
          )}
        >
          <Download className="w-3.5 h-3.5" />
          Export
        </button>
      </div>

      <DeducteeTable deductees={filtered} />
    </div>
  );
}
