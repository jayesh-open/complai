"use client";

import { useRouter } from "next/navigation";
import { cn, formatINR } from "@complai/ui-components";
import { DataTable } from "@complai/ui-components";
import { IRNStatusPill } from "./IRNStatusPill";
import type { EInvoiceRecord } from "../types";

interface EInvoiceTableProps {
  records: EInvoiceRecord[];
  className?: string;
}

export function EInvoiceTable({ records, className }: EInvoiceTableProps) {
  const router = useRouter();

  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as EInvoiceRecord;

  const columns = [
    {
      key: "invoiceNo",
      header: "Invoice No.",
      render: (row: Row) => (
        <span className="font-mono text-[var(--text-primary)] font-medium">
          {r(row).invoiceNo}
        </span>
      ),
    },
    {
      key: "gstin",
      header: "GSTIN",
      render: (row: Row) => (
        <span className="font-mono text-[11px]">{r(row).gstin}</span>
      ),
    },
    {
      key: "buyerName",
      header: "Buyer",
      render: (row: Row) => (
        <div>
          <div className="text-[var(--text-primary)] text-xs">
            {r(row).buyerName}
          </div>
          <div className="font-mono text-[10px] text-[var(--text-muted)]">
            {r(row).buyerGstin}
          </div>
        </div>
      ),
    },
    {
      key: "irn",
      header: "IRN",
      render: (row: Row) => (
        <span className="font-mono text-[10px] text-[var(--text-muted)]">
          {r(row).irn.slice(0, 20)}...
        </span>
      ),
    },
    {
      key: "totalValue",
      header: "Total Value",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono">{formatINR(r(row).totalValue)}</span>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (row: Row) => (
        <IRNStatusPill status={r(row).status} />
      ),
    },
    {
      key: "generatedAt",
      header: "Generated",
      render: (row: Row) => {
        const d = new Date(r(row).generatedAt);
        const dd = String(d.getDate()).padStart(2, "0");
        const mm = String(d.getMonth() + 1).padStart(2, "0");
        const yyyy = d.getFullYear();
        const hh = String(d.getHours()).padStart(2, "0");
        const min = String(d.getMinutes()).padStart(2, "0");
        return (
          <span className="text-[11px] text-[var(--text-muted)] tabular-nums">
            {dd}/{mm}/{yyyy} {hh}:{min}
          </span>
        );
      },
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={records as unknown as Row[]}
      onRowClick={(row) => {
        const rec = r(row);
        router.push(`/compliance/e-invoicing/${rec.gstin}/${rec.irn}`);
      }}
      emptyMessage="No e-invoices found"
      className={className}
    />
  );
}
