"use client";

import { DataTable, formatINR } from "@complai/ui-components";
import { FilingStatusPill } from "./EntryStatusPill";
import type { ReturnStatus } from "../types";

interface FilingStatusTableProps {
  returns: ReturnStatus[];
  className?: string;
}

export function FilingStatusTable({ returns, className }: FilingStatusTableProps) {
  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as ReturnStatus;

  const columns = [
    {
      key: "formLabel",
      header: "Form",
      render: (row: Row) => (
        <span className="text-[var(--text-primary)] text-xs font-medium">
          {r(row).formLabel}
        </span>
      ),
    },
    {
      key: "taxYear",
      header: "Tax Year",
      render: (row: Row) => (
        <span className="font-mono text-[11px]">{r(row).taxYear}</span>
      ),
    },
    {
      key: "quarter",
      header: "Quarter",
      render: (row: Row) => (
        <span className="font-mono text-[11px]">{r(row).quarter}</span>
      ),
    },
    {
      key: "dueDate",
      header: "Due Date",
      render: (row: Row) => (
        <span className="text-[11px] text-[var(--text-muted)] tabular-nums">
          {r(row).dueDate}
        </span>
      ),
    },
    {
      key: "entryCount",
      header: "Entries",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs">{r(row).entryCount}</span>
      ),
    },
    {
      key: "totalTds",
      header: "Total TDS",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs">
          {r(row).totalTds > 0 ? formatINR(r(row).totalTds) : "—"}
        </span>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (row: Row) => <FilingStatusPill status={r(row).status} />,
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={returns as unknown as Row[]}
      emptyMessage="No returns found"
      className={className}
    />
  );
}
