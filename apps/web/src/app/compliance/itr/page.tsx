"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, Plus, Download } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_EMPLOYEES, ALL_BATCHES } from "./mock-data";
import { ITRKpis } from "./components/ITRKpis";
import { EmployeeITRTable } from "./components/EmployeeITRTable";
import { BatchTable } from "./components/BatchTable";
import { TaxYearSelector } from "../tds/components/TaxYearSelector";
import type { FilingStatus, ITRForm } from "./types";

type Tab = "employees" | "batches";
type StatusFilter = "ALL" | FilingStatus;
type FormFilter = "ALL" | ITRForm;

const TABS: { key: Tab; label: string }[] = [
  { key: "employees", label: "All Employees" },
  { key: "batches", label: "Bulk Batches" },
];

export default function ITRLandingPage() {
  const [tab, setTab] = useState<Tab>("employees");
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");
  const [formFilter, setFormFilter] = useState<FormFilter>("ALL");
  const [taxYear, setTaxYear] = useState("2026-27");

  const filteredEmployees = useMemo(() => {
    return ALL_EMPLOYEES.filter((e) => {
      if (e.taxYear !== taxYear) return false;
      if (statusFilter !== "ALL" && e.filingStatus !== statusFilter) return false;
      if (formFilter !== "ALL" && e.recommendedForm !== formFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return e.pan.toLowerCase().includes(q) || e.name.toLowerCase().includes(q);
      }
      return true;
    });
  }, [search, statusFilter, formFilter, taxYear]);

  const filteredBatches = useMemo(() => {
    return ALL_BATCHES.filter((b) => b.taxYear === taxYear);
  }, [taxYear]);

  return (
    <div className="space-y-6" data-testid="itr-landing">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            Income Tax Returns
          </h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            Employee bulk ITR filing — Section 202 (ITA 2025)
          </p>
        </div>
        <div className="flex items-center gap-3">
          <TaxYearSelector value={taxYear} onChange={setTaxYear} />
          <Link
            href="/compliance/itr/bulk/new"
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)]",
              "hover:bg-[var(--accent-hover)] transition-colors"
            )}
          >
            <Plus className="w-3.5 h-3.5" />
            New Batch
          </Link>
        </div>
      </div>

      <ITRKpis employees={ALL_EMPLOYEES} />

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

        {tab === "employees" && (
          <>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
              className={cn(
                "px-3 py-2 rounded-lg text-xs",
                "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                "text-[var(--text-primary)]",
                "focus:outline-none focus:border-[var(--accent)]"
              )}
            >
              <option value="ALL">All Statuses</option>
              <option value="NOT_STARTED">Not Started</option>
              <option value="AIS_FETCHED">AIS Fetched</option>
              <option value="FORM_GENERATED">Form Ready</option>
              <option value="REVIEW_PENDING">Review Pending</option>
              <option value="EMPLOYEE_APPROVED">Approved</option>
              <option value="FILED">Filed</option>
              <option value="ACKNOWLEDGED">Acknowledged</option>
            </select>
            <select
              value={formFilter}
              onChange={(e) => setFormFilter(e.target.value as FormFilter)}
              className={cn(
                "px-3 py-2 rounded-lg text-xs",
                "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                "text-[var(--text-primary)]",
                "focus:outline-none focus:border-[var(--accent)]"
              )}
            >
              <option value="ALL">All Forms</option>
              <option value="ITR-1">ITR-1</option>
              <option value="ITR-2">ITR-2</option>
              <option value="ITR-3">ITR-3</option>
              <option value="ITR-4">ITR-4</option>
            </select>
          </>
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

      {tab === "employees" && <EmployeeITRTable employees={filteredEmployees} />}
      {tab === "batches" && <BatchTable batches={filteredBatches} />}
    </div>
  );
}
