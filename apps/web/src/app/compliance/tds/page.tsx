"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, Calculator, Upload, Download } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_DEDUCTEES, ALL_ENTRIES, ALL_RETURNS } from "./mock-data";
import { TdsKpis } from "./components/TdsKpis";
import { DeducteeTable } from "./components/DeducteeTable";
import { RecentEntriesTable } from "./components/RecentEntriesTable";
import { FilingStatusTable } from "./components/FilingStatusTable";
import { TaxYearSelector } from "./components/TaxYearSelector";
import type { Section2025, ResidencyStatus } from "./types";

type Tab = "deductees" | "entries" | "filings";
type SectionFilter = "ALL" | Section2025;
type ResidencyFilter = "ALL" | ResidencyStatus;

const TABS: { key: Tab; label: string }[] = [
  { key: "deductees", label: "All Deductees" },
  { key: "entries", label: "Recent Entries" },
  { key: "filings", label: "Filing Status" },
];

export default function TdsLandingPage() {
  const [tab, setTab] = useState<Tab>("deductees");
  const [search, setSearch] = useState("");
  const [sectionFilter, setSectionFilter] = useState<SectionFilter>("ALL");
  const [residencyFilter, setResidencyFilter] = useState<ResidencyFilter>("ALL");
  const [taxYear, setTaxYear] = useState("2026-27");

  const filteredDeductees = useMemo(() => {
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

  const filteredEntries = useMemo(() => {
    return ALL_ENTRIES.filter((e) => {
      if (e.taxYear !== taxYear) return false;
      if (sectionFilter !== "ALL" && e.section2025 !== sectionFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return e.deducteePan.toLowerCase().includes(q) || e.deducteeName.toLowerCase().includes(q);
      }
      return true;
    }).slice(0, 30);
  }, [search, sectionFilter, taxYear]);

  const filteredReturns = useMemo(() => {
    return ALL_RETURNS.filter((r) => r.taxYear === taxYear);
  }, [taxYear]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            TDS / TCS
          </h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            Tax Deducted at Source — Sections 392, 393 (ITA 2025)
          </p>
        </div>
        <div className="flex items-center gap-3">
          <TaxYearSelector value={taxYear} onChange={setTaxYear} />
          <Link
            href="/compliance/tds/import"
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors"
            )}
          >
            <Upload className="w-3.5 h-3.5" />
            Import
          </Link>
          <Link
            href="/compliance/tds/calculate"
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)]",
              "hover:bg-[var(--accent-hover)] transition-colors"
            )}
          >
            <Calculator className="w-3.5 h-3.5" />
            Calculate TDS
          </Link>
        </div>
      </div>

      <TdsKpis entries={ALL_ENTRIES} returns={ALL_RETURNS} />

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

        <div className="flex items-center gap-1 p-1 rounded-lg bg-[var(--bg-secondary)] border border-[var(--border-default)]">
          {TABS.map((t) => (
            <button
              key={t.key}
              onClick={() => setTab(t.key)}
              className={cn(
                "px-3 py-1.5 rounded-md text-[10px] font-semibold uppercase tracking-wide transition-colors",
                tab === t.key
                  ? "bg-[var(--accent)] text-[var(--accent-text)]"
                  : "text-[var(--text-muted)] hover:text-[var(--text-primary)]"
              )}
            >
              {t.label}
            </button>
          ))}
        </div>

        {tab !== "filings" && (
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
        )}

        {tab === "deductees" && (
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
        )}

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

      {tab === "deductees" && <DeducteeTable deductees={filteredDeductees} />}
      {tab === "entries" && <RecentEntriesTable entries={filteredEntries} />}
      {tab === "filings" && <FilingStatusTable returns={filteredReturns} />}
    </div>
  );
}
