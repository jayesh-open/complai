"use client";

import { useMemo } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { ArrowLeft, RefreshCw, Send } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_BATCHES, ALL_EMPLOYEES } from "../../mock-data";
import { BatchStatusPill } from "../../components/BatchStatusPill";
import { BulkBatchProgressBar } from "../../components/BulkBatchProgressBar";
import { EmployeeITRTable } from "../../components/EmployeeITRTable";

export default function BatchDetailPage() {
  const params = useParams();
  const batchId = params["batch-id"] as string;

  const batch = useMemo(() => ALL_BATCHES.find((b) => b.id === batchId), [batchId]);
  const batchEmployees = useMemo(
    () => ALL_EMPLOYEES.slice(0, batch?.totalEmployees ?? 5),
    [batch]
  );

  if (!batch) {
    return (
      <div className="space-y-6">
        <Link
          href="/compliance/itr"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to ITR
        </Link>
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-12 text-center">
          <p className="text-sm text-[var(--text-muted)]">Batch not found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/itr"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to ITR
      </Link>

      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-heading-lg text-[var(--text-primary)]">{batch.name}</h1>
            <BatchStatusPill status={batch.status} />
          </div>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            Tax Year {batch.taxYear} · Created {batch.createdAt} by {batch.createdBy}
          </p>
        </div>
        <div className="flex items-center gap-3">
          {batch.status === "IN_PROGRESS" && (
            <button
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-medium",
                "border border-[var(--border-default)]",
                "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
              )}
            >
              <RefreshCw className="w-3.5 h-3.5" />
              Retry Failed
            </button>
          )}
          {batch.status === "DRAFT" && (
            <button
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
                "bg-[var(--accent)] text-[var(--accent-text)]",
                "hover:bg-[var(--accent-hover)] transition-colors"
              )}
            >
              <Send className="w-3.5 h-3.5" />
              Start Batch
            </button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-4 gap-4">
        <StatCard label="Total" value={batch.totalEmployees} />
        <StatCard label="Filed" value={batch.filed} color="text-[var(--success)]" />
        <StatCard label="Pending" value={batch.pending} color="text-[var(--warning)]" />
        <StatCard label="Failed" value={batch.failed} color="text-[var(--danger)]" />
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
        <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide mb-3">
          Overall Progress
        </p>
        <BulkBatchProgressBar
          filed={batch.filed}
          pending={batch.pending}
          failed={batch.failed}
          total={batch.totalEmployees}
        />
      </div>

      <div>
        <h2 className="text-sm font-semibold text-[var(--text-primary)] mb-3">
          Employees in Batch
        </h2>
        <EmployeeITRTable employees={batchEmployees} batchId={batchId} />
      </div>
    </div>
  );
}

function StatCard({ label, value, color }: { label: string; value: number; color?: string }) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 text-center">
      <p className={cn("text-2xl font-bold tabular-nums", color || "text-[var(--text-primary)]")}>{value}</p>
      <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide mt-1">{label}</p>
    </div>
  );
}
