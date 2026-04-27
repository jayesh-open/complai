"use client";

import { FileCheck2, XCircle, Clock, FileText } from "lucide-react";
import { cn } from "@/lib/utils";
import type { EInvoiceRecord } from "../types";

interface EInvoiceKPIsProps {
  records: EInvoiceRecord[];
  className?: string;
}

function StatCard({
  icon: Icon,
  iconColor,
  label,
  value,
}: {
  icon: React.ElementType;
  iconColor: string;
  label: string;
  value: string | number;
}) {
  return (
    <div
      className={cn(
        "flex items-center gap-3 px-4 py-3 rounded-xl",
        "bg-[var(--bg-secondary)] border border-[var(--border-default)]"
      )}
    >
      <div
        className={cn(
          "w-9 h-9 rounded-lg flex items-center justify-center",
          iconColor
        )}
      >
        <Icon className="w-4.5 h-4.5 text-white" />
      </div>
      <div>
        <div className="text-lg font-bold text-[var(--text-primary)] tabular-nums">
          {value}
        </div>
        <div className="text-[10px] text-[var(--text-muted)] font-medium uppercase tracking-wide">
          {label}
        </div>
      </div>
    </div>
  );
}

export function EInvoiceKPIs({ records, className }: EInvoiceKPIsProps) {
  const total = records.length;
  const generated = records.filter((r) => r.status === "GENERATED").length;
  const cancelled = records.filter((r) => r.status === "CANCELLED").length;
  const last24h = records.filter((r) => {
    const diff = Date.now() - new Date(r.generatedAt).getTime();
    return diff < 24 * 3600 * 1000;
  }).length;

  return (
    <div className={cn("grid grid-cols-4 gap-4", className)}>
      <StatCard
        icon={FileText}
        iconColor="bg-[var(--accent)]"
        label="Total IRNs"
        value={total}
      />
      <StatCard
        icon={FileCheck2}
        iconColor="bg-[var(--success)]"
        label="Generated"
        value={generated}
      />
      <StatCard
        icon={XCircle}
        iconColor="bg-[var(--danger)]"
        label="Cancelled"
        value={cancelled}
      />
      <StatCard
        icon={Clock}
        iconColor="bg-[var(--info)]"
        label="Last 24h"
        value={last24h}
      />
    </div>
  );
}
