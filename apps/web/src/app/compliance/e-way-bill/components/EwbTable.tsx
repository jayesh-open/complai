"use client";

import { useRouter } from "next/navigation";
import { DataTable } from "@complai/ui-components";
import { EwbStatusPill } from "./EwbStatusPill";
import type { EwbRecord } from "../types";

interface EwbTableProps {
  records: EwbRecord[];
  className?: string;
}

function isNearExpiry(r: EwbRecord): boolean {
  if (r.status !== "ACTIVE") return false;
  const remaining = new Date(r.validUntil).getTime() - Date.now();
  return remaining > 0 && remaining < 4 * 3600000;
}

function fmtDate(iso: string): string {
  const d = new Date(iso);
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const min = String(d.getMinutes()).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()} ${hh}:${min}`;
}

export function EwbTable({ records, className }: EwbTableProps) {
  const router = useRouter();

  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as EwbRecord;

  const columns = [
    {
      key: "ewbNumber",
      header: "EWB Number",
      render: (row: Row) => (
        <span className="font-mono text-[var(--text-primary)] font-medium text-xs">
          {r(row).ewbNumber}
        </span>
      ),
    },
    {
      key: "invoiceNo",
      header: "Source Invoice",
      render: (row: Row) => (
        <span className="font-mono text-[11px] text-[var(--text-secondary)]">
          {r(row).invoiceNo}
        </span>
      ),
    },
    {
      key: "vehicleNo",
      header: "Vehicle",
      render: (row: Row) => (
        <span className="font-mono text-xs text-[var(--text-primary)]">
          {r(row).vehicleNo}
        </span>
      ),
    },
    {
      key: "distanceKm",
      header: "Distance",
      align: "right" as const,
      render: (row: Row) => (
        <span className="text-xs tabular-nums">{r(row).distanceKm} km</span>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (row: Row) => (
        <EwbStatusPill status={r(row).status} nearingExpiry={isNearExpiry(r(row))} />
      ),
    },
    {
      key: "validUntil",
      header: "Valid Until",
      render: (row: Row) => (
        <span className="text-[11px] text-[var(--text-muted)] tabular-nums">
          {fmtDate(r(row).validUntil)}
        </span>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={records as unknown as Row[]}
      onRowClick={(row) => {
        router.push(`/compliance/e-way-bill/${r(row).ewbNumber}`);
      }}
      emptyMessage="No e-way bills found"
      className={className}
    />
  );
}
