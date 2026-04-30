"use client";

import { useRouter } from "next/navigation";
import { DataTable, formatINR } from "@complai/ui-components";
import { SectionPill } from "./SectionPill";
import type { Deductee } from "../types";

interface DeducteeTableProps {
  deductees: Deductee[];
  className?: string;
}

export function DeducteeTable({ deductees, className }: DeducteeTableProps) {
  const router = useRouter();
  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as Deductee;

  const columns = [
    {
      key: "pan",
      header: "PAN",
      render: (row: Row) => (
        <span className="font-mono text-[var(--text-primary)] font-medium text-[11px]">
          {r(row).pan}
        </span>
      ),
    },
    {
      key: "name",
      header: "Name",
      render: (row: Row) => (
        <div>
          <div className="text-[var(--text-primary)] text-xs font-medium">
            {r(row).name}
          </div>
          <div className="text-[10px] text-[var(--text-muted)]">
            {r(row).category} · {r(row).residency === "NON_RESIDENT" ? "Non-Resident" : "Resident"}
          </div>
        </div>
      ),
    },
    {
      key: "sectionPreference",
      header: "Section",
      render: (row: Row) => <SectionPill section={r(row).sectionPreference} />,
    },
    {
      key: "totalDeductedYTD",
      header: "Deducted YTD",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs">{formatINR(r(row).totalDeductedYTD)}</span>
      ),
    },
    {
      key: "lastTransactionDate",
      header: "Last Transaction",
      render: (row: Row) => (
        <span className="text-[11px] text-[var(--text-muted)] tabular-nums">
          {r(row).lastTransactionDate}
        </span>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={deductees as unknown as Row[]}
      onRowClick={(row) => router.push(`/compliance/tds/deductees/${r(row).id}`)}
      emptyMessage="No deductees found"
      className={className}
    />
  );
}
