"use client";

import { Users, FileCheck, Clock, AlertTriangle } from "lucide-react";
import { cn } from "@/lib/utils";
import type { ITREmployee } from "../types";

interface ITRKpisProps {
  employees: ITREmployee[];
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
      <div className={cn("w-9 h-9 rounded-lg flex items-center justify-center", iconColor)}>
        <Icon className="w-4.5 h-4.5 text-white" />
      </div>
      <div>
        <div className="text-lg font-bold text-[var(--text-primary)] tabular-nums">{value}</div>
        <div className="text-[10px] text-[var(--text-muted)] font-medium uppercase tracking-wide">{label}</div>
        {subtitle && <div className="text-[10px] text-[var(--text-muted)] mt-0.5">{subtitle}</div>}
      </div>
    </div>
  );
}

export function ITRKpis({ employees, className }: ITRKpisProps) {
  const total = employees.length;
  const filed = employees.filter((e) => e.filingStatus === "FILED" || e.filingStatus === "ACKNOWLEDGED").length;
  const pending = employees.filter((e) => e.filingStatus === "REVIEW_PENDING" || e.filingStatus === "EMPLOYEE_APPROVED").length;
  const notStarted = employees.filter((e) => e.filingStatus === "NOT_STARTED").length;

  return (
    <div className={cn("grid grid-cols-4 gap-4", className)}>
      <StatCard icon={Users} iconColor="bg-[var(--accent)]" label="Total Employees" value={total} subtitle="Tax Year 2026-27" />
      <StatCard icon={FileCheck} iconColor="bg-[var(--success)]" label="Filed" value={filed} subtitle={`${Math.round((filed / total) * 100)}% complete`} />
      <StatCard icon={Clock} iconColor="bg-[var(--warning)]" label="Pending Review" value={pending} subtitle="Awaiting approval" />
      <StatCard icon={AlertTriangle} iconColor="bg-[var(--danger)]" label="Not Started" value={notStarted} subtitle={notStarted > 0 ? "Action required" : "All initiated"} />
    </div>
  );
}
