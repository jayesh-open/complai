"use client";

import { Truck, CheckCircle2, XCircle, AlertTriangle } from "lucide-react";
import { cn } from "@/lib/utils";
import type { EwbRecord } from "../types";

interface EwbKPIsProps {
  records: EwbRecord[];
}

interface StatCardProps {
  label: string;
  value: number;
  icon: React.ElementType;
  color: string;
}

function StatCard({ label, value, icon: Icon, color }: StatCardProps) {
  return (
    <div className={cn(
      "flex items-center gap-3 p-4 rounded-xl",
      "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
    )}>
      <div className={cn(
        "w-9 h-9 rounded-lg flex items-center justify-center",
        "bg-[var(--bg-tertiary)]",
      )}>
        <Icon className={cn("w-4 h-4", color)} />
      </div>
      <div>
        <div className="text-lg font-bold text-[var(--text-primary)] tabular-nums">
          {value}
        </div>
        <div className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide font-semibold">
          {label}
        </div>
      </div>
    </div>
  );
}

export function EwbKPIs({ records }: EwbKPIsProps) {
  const now = Date.now();
  const fourHours = 4 * 3600000;
  const active = records.filter((r) => r.status === "ACTIVE").length;
  const nearExpiry = records.filter(
    (r) => r.status === "ACTIVE" &&
      new Date(r.validUntil).getTime() - now > 0 &&
      new Date(r.validUntil).getTime() - now < fourHours,
  ).length;
  const cancelled = records.filter((r) => r.status === "CANCELLED").length;
  const expired = records.filter((r) => r.status === "EXPIRED").length;

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard label="Total EWBs" value={records.length} icon={Truck} color="text-[var(--accent)]" />
      <StatCard label="Active" value={active} icon={CheckCircle2} color="text-[var(--success)]" />
      <StatCard label="Nearing Expiry" value={nearExpiry} icon={AlertTriangle} color="text-[var(--warning)]" />
      <StatCard label="Cancelled" value={cancelled + expired} icon={XCircle} color="text-[var(--danger)]" />
    </div>
  );
}
