"use client";

import { useState, useMemo } from "react";
import { Search, Download } from "lucide-react";
import Link from "next/link";
import { cn, formatINR } from "@complai/ui-components";
import { DataTable } from "@complai/ui-components";
import { generateChallanRows, generateReconciliationSummary } from "./mock-data";
import { ChallanStatusPill } from "../components/ChallanStatusPill";
import { TaxYearSelector } from "../components/TaxYearSelector";
import type { ChallanStatus } from "./types";

type StatusFilter = "ALL" | ChallanStatus;
type QuarterFilter = "ALL" | "Q1" | "Q2" | "Q3" | "Q4";

export default function ChallansListPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");
  const [quarterFilter, setQuarterFilter] = useState<QuarterFilter>("ALL");
  const [taxYear, setTaxYear] = useState("2026-27");

  const rows = useMemo(() => generateChallanRows(), []);
  const recon = useMemo(() => generateReconciliationSummary(), []);

  const filtered = useMemo(() => {
    return rows.filter((r) => {
      if (statusFilter !== "ALL" && r.status !== statusFilter) return false;
      if (quarterFilter !== "ALL" && r.quarter !== quarterFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return r.challanSerial.toLowerCase().includes(q) || r.bsrCode.toLowerCase().includes(q);
      }
      return true;
    });
  }, [rows, search, statusFilter, quarterFilter]);

  const allocatedPct = recon.totalDeposited > 0 ? Math.round((recon.totalAllocated / recon.totalDeposited) * 100) : 0;

  type Row = Record<string, unknown>;
  const r2 = (row: Row) => row as unknown as (typeof rows)[0];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">Challan Tracking</h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">TDS challan deposits and allocation reconciliation</p>
        </div>
        <TaxYearSelector value={taxYear} onChange={setTaxYear} />
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
        <div className="flex items-center justify-between mb-3">
          <span className="text-body-sm font-medium text-[var(--text-primary)]">Reconciliation</span>
          <span className="text-body-sm text-[var(--text-muted)]">{allocatedPct}% allocated</span>
        </div>
        <div className="h-3 rounded-full bg-[var(--bg-tertiary)] overflow-hidden mb-3">
          <div className="h-full rounded-full bg-[var(--accent)]" style={{ width: `${allocatedPct}%` }} />
        </div>
        <div className="grid grid-cols-4 gap-4">
          <ReconStat label="Total Deposited" value={formatINR(recon.totalDeposited)} />
          <ReconStat label="Allocated" value={formatINR(recon.totalAllocated)} accent />
          <ReconStat label="Unallocated" value={formatINR(recon.totalUnallocated)} warn={recon.totalUnallocated > 0} />
          <ReconStat label="Challans" value={`${recon.fullyAllocated} full · ${recon.partiallyAllocated} partial · ${recon.unallocated} none`} small />
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
          <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search by challan no. or BSR code..."
            className={cn("w-full pl-9 pr-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]", "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]")} />
        </div>
        <select value={quarterFilter} onChange={(e) => setQuarterFilter(e.target.value as QuarterFilter)}
          className={cn("px-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]")}>
          <option value="ALL">All Quarters</option>
          <option value="Q1">Q1</option><option value="Q2">Q2</option><option value="Q3">Q3</option><option value="Q4">Q4</option>
        </select>
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
          className={cn("px-3 py-2 rounded-lg text-xs", "bg-[var(--bg-secondary)] border border-[var(--border-default)]", "text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]")}>
          <option value="ALL">All Statuses</option>
          <option value="CLEARED">Cleared</option><option value="PENDING">Pending</option><option value="REJECTED">Rejected</option>
        </select>
        <button className={cn("flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium", "border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors")}>
          <Download className="w-3.5 h-3.5" /> Export
        </button>
      </div>

      <DataTable
        data={filtered as unknown as Row[]}
        columns={[
          { header: "Challan No.", key: "challanSerial", render: (row: Row) => (
            <Link href={`/compliance/tds/challans/${r2(row).challanId}`} className="text-[var(--accent)] hover:underline font-mono font-medium">{r2(row).challanSerial}</Link>
          )},
          { header: "BSR Code", key: "bsrCode", render: (row: Row) => <span className="font-mono text-[var(--text-muted)]">{r2(row).bsrCode}</span> },
          { header: "Date", key: "depositDate" },
          { header: "Quarter", key: "quarter" },
          { header: "Amount", key: "amount", render: (row: Row) => <span className="font-mono">{formatINR(r2(row).amount)}</span> },
          { header: "Allocated", key: "allocatedAmount", render: (row: Row) => <span className="font-mono text-[var(--success)]">{formatINR(r2(row).allocatedAmount)}</span> },
          { header: "Unallocated", key: "unallocatedAmount", render: (row: Row) => (
            <span className={cn("font-mono", r2(row).unallocatedAmount > 0 ? "text-[var(--warning)]" : "text-[var(--text-muted)]")}>{formatINR(r2(row).unallocatedAmount)}</span>
          )},
          { header: "Bank", key: "bankName" },
          { header: "Status", key: "status", render: (row: Row) => <ChallanStatusPill status={r2(row).status} /> },
        ]}
      />
    </div>
  );
}

function ReconStat({ label, value, accent, warn, small }: { label: string; value: string; accent?: boolean; warn?: boolean; small?: boolean }) {
  return (
    <div>
      <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">{label}</div>
      <div className={cn("mt-0.5 font-semibold tabular-nums", small ? "text-xs" : "text-sm", accent && "text-[var(--accent)]", warn && "text-[var(--warning)]", !accent && !warn && "text-[var(--text-primary)]")}>{value}</div>
    </div>
  );
}
