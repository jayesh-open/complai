"use client";

import { IndianRupee, AlertTriangle, FileText, Clock } from "lucide-react";
import { cn, formatCompact } from "@complai/ui-components";
import type { TDSEntry, ReturnStatus } from "../types";

interface TdsKpisProps {
  entries: TDSEntry[];
  returns: ReturnStatus[];
  className?: string;
}

function StatCard({
  icon: Icon,
  iconColor,
  label,
  value,
  subtitle,
}: {
  icon: React.ElementType;
  iconColor: string;
  label: string;
  value: string | number;
  subtitle?: string;
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
        {subtitle && (
          <div className="text-[10px] text-[var(--text-muted)] mt-0.5">
            {subtitle}
          </div>
        )}
      </div>
    </div>
  );
}

export function TdsKpis({ entries, returns, className }: TdsKpisProps) {
  const totalDeducted = entries.reduce((s, e) => s + e.totalTax, 0);
  const pendingChallans = entries.filter((e) => e.status === "PENDING").length;
  const openReturns = returns.filter(
    (r) => r.status === "PENDING" || r.status === "OVERDUE"
  ).length;
  const overdueReturns = returns.filter((r) => r.status === "OVERDUE").length;

  return (
    <div className={cn("grid grid-cols-4 gap-4", className)}>
      <StatCard
        icon={IndianRupee}
        iconColor="bg-[var(--accent)]"
        label="Total Deducted"
        value={formatCompact(totalDeducted)}
        subtitle="Tax Year 2026-27"
      />
      <StatCard
        icon={Clock}
        iconColor="bg-[var(--warning)]"
        label="Pending Challans"
        value={pendingChallans}
        subtitle={pendingChallans > 0 ? "Awaiting deposit" : "All deposited"}
      />
      <StatCard
        icon={FileText}
        iconColor="bg-[var(--info)]"
        label="Open Returns"
        value={openReturns}
        subtitle="Forms 138/140/144"
      />
      <StatCard
        icon={AlertTriangle}
        iconColor="bg-[var(--danger)]"
        label="Late Filing Risk"
        value={overdueReturns}
        subtitle={overdueReturns > 0 ? "Action required" : "On track"}
      />
    </div>
  );
}
