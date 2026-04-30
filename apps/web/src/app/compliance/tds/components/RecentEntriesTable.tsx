"use client";

import { DataTable, formatINR } from "@complai/ui-components";
import { SectionPill } from "./SectionPill";
import { EntryStatusPill } from "./EntryStatusPill";
import type { TDSEntry } from "../types";

interface RecentEntriesTableProps {
  entries: TDSEntry[];
  className?: string;
}

export function RecentEntriesTable({ entries, className }: RecentEntriesTableProps) {
  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as TDSEntry;

  const columns = [
    {
      key: "deducteeName",
      header: "Deductee",
      render: (row: Row) => (
        <div>
          <div className="text-[var(--text-primary)] text-xs font-medium">
            {r(row).deducteeName}
          </div>
          <div className="font-mono text-[10px] text-[var(--text-muted)]">
            {r(row).deducteePan}
          </div>
        </div>
      ),
    },
    {
      key: "section2025",
      header: "Section",
      render: (row: Row) => <SectionPill section={r(row).section2025} />,
    },
    {
      key: "paymentCode",
      header: "Code",
      render: (row: Row) => (
        <span className="font-mono text-[11px] text-[var(--text-primary)]">
          {r(row).paymentCode}
        </span>
      ),
    },
    {
      key: "grossAmount",
      header: "Amount",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs">{formatINR(r(row).grossAmount)}</span>
      ),
    },
    {
      key: "totalTax",
      header: "TDS",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs font-semibold">{formatINR(r(row).totalTax)}</span>
      ),
    },
    {
      key: "transactionDate",
      header: "Date",
      render: (row: Row) => {
        const d = new Date(r(row).transactionDate);
        const dd = String(d.getDate()).padStart(2, "0");
        const mm = String(d.getMonth() + 1).padStart(2, "0");
        return (
          <span className="text-[11px] text-[var(--text-muted)] tabular-nums">
            {dd}/{mm}/{d.getFullYear()}
          </span>
        );
      },
    },
    {
      key: "status",
      header: "Status",
      render: (row: Row) => <EntryStatusPill status={r(row).status} />,
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={entries as unknown as Row[]}
      emptyMessage="No entries found"
      className={className}
    />
  );
}
