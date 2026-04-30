"use client";

import { useState, useMemo } from "react";
import { Search, Download, FileText } from "lucide-react";
import Link from "next/link";
import { cn, formatINR } from "@complai/ui-components";
import { DataTable } from "@complai/ui-components";
import { generateForm131Rows } from "../mock-data";
import { CertificateStatusPill } from "../../components/CertificateStatusPill";
import { TaxYearSelector } from "../../components/TaxYearSelector";
import type { CertificateStatus } from "../types";

type StatusFilter = "ALL" | CertificateStatus;
type QuarterFilter = "ALL" | "Q1" | "Q2" | "Q3" | "Q4";

export default function Form131ListPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");
  const [quarterFilter, setQuarterFilter] = useState<QuarterFilter>("ALL");
  const [taxYear, setTaxYear] = useState("2026-27");

  const rows = useMemo(() => generateForm131Rows(), []);
  const filtered = useMemo(() => {
    return rows.filter((r) => {
      if (statusFilter !== "ALL" && r.status !== statusFilter) return false;
      if (quarterFilter !== "ALL" && r.quarter !== quarterFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return r.pan.toLowerCase().includes(q) || r.deducteeName.toLowerCase().includes(q);
      }
      return true;
    });
  }, [rows, search, statusFilter, quarterFilter]);

  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as (typeof rows)[0];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">Form 131 Certificates</h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">Quarterly non-salary TDS certificates (ITA 2025 §393)</p>
        </div>
        <div className="flex items-center gap-3">
          <TaxYearSelector value={taxYear} onChange={setTaxYear} />
          <button className={cn("flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-semibold", "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors")}>
            <FileText className="w-3.5 h-3.5" /> Bulk Generate
          </button>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
          <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search by PAN or name..."
            className={cn("w-full pl-9 pr-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]", "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]")} />
        </div>
        <select value={quarterFilter} onChange={(e) => setQuarterFilter(e.target.value as QuarterFilter)}
          className={cn("px-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]")}>
          <option value="ALL">All Quarters</option>
          <option value="Q1">Q1</option>
          <option value="Q2">Q2</option>
          <option value="Q3">Q3</option>
          <option value="Q4">Q4</option>
        </select>
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
          className={cn("px-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]")}>
          <option value="ALL">All Statuses</option>
          <option value="ISSUED">Issued</option>
          <option value="PENDING">Pending</option>
          <option value="GENERATED">Generated</option>
        </select>
        <button className={cn("flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium", "border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors")}>
          <Download className="w-3.5 h-3.5" /> Export
        </button>
      </div>

      <DataTable
        data={filtered as unknown as Row[]}
        columns={[
          { header: "Deductee", key: "deducteeName", render: (row: Row) => (
            <Link href={`/compliance/tds/certificates/131/${r(row).deducteeId}/${taxYear}/${r(row).quarter.toLowerCase()}`} className="text-[var(--accent)] hover:underline font-medium">{r(row).deducteeName}</Link>
          )},
          { header: "PAN", key: "pan", render: (row: Row) => <span className="font-mono text-[var(--text-muted)]">{r(row).pan}</span> },
          { header: "Quarter", key: "quarter", render: (row: Row) => <span className="font-semibold text-xs">{r(row).quarter}</span> },
          { header: "Sections", key: "sectionCodes", render: (row: Row) => (
            <div className="flex gap-1">{r(row).sectionCodes.slice(0, 3).map((c: string) => <span key={c} className="px-1.5 py-0.5 rounded text-[9px] font-mono bg-[var(--bg-tertiary)] text-[var(--text-muted)]">{c}</span>)}</div>
          )},
          { header: "Amount", key: "totalAmount", render: (row: Row) => <span className="font-mono">{formatINR(r(row).totalAmount)}</span> },
          { header: "TDS", key: "totalTds", render: (row: Row) => <span className="font-mono">{formatINR(r(row).totalTds)}</span> },
          { header: "Status", key: "status", render: (row: Row) => <CertificateStatusPill status={r(row).status} /> },
        ]}
      />
    </div>
  );
}
