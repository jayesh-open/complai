"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";
import { BatchStatusPill } from "./BatchStatusPill";
import { BulkBatchProgressBar } from "./BulkBatchProgressBar";
import type { BulkBatch } from "../types";

interface BatchTableProps {
  batches: BulkBatch[];
  className?: string;
}

export function BatchTable({ batches, className }: BatchTableProps) {
  return (
    <div className={cn("bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden", className)}>
      <table className="w-full">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Batch</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Status</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Progress</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">Created</th>
            <th className="px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">By</th>
          </tr>
        </thead>
        <tbody>
          {batches.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center py-12 text-[var(--text-muted)] text-sm">
                No batches created yet
              </td>
            </tr>
          ) : (
            batches.map((batch) => (
              <tr
                key={batch.id}
                className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors"
              >
                <td className="px-[18px] py-3">
                  <Link href={`/compliance/itr/bulk/${batch.id}`} className="hover:underline">
                    <div className="text-xs font-medium text-[var(--text-primary)]">{batch.name}</div>
                    <div className="text-[10px] text-[var(--text-muted)]">{batch.totalEmployees} employees</div>
                  </Link>
                </td>
                <td className="px-[18px] py-3"><BatchStatusPill status={batch.status} /></td>
                <td className="px-[18px] py-3 min-w-[160px]">
                  <BulkBatchProgressBar
                    filed={batch.filed}
                    pending={batch.pending}
                    failed={batch.failed}
                    total={batch.totalEmployees}
                  />
                </td>
                <td className="px-[18px] py-3 text-xs text-[var(--text-secondary)]">{batch.createdAt}</td>
                <td className="px-[18px] py-3 text-xs text-[var(--text-secondary)]">{batch.createdBy}</td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
