"use client";

import Link from "next/link";
import { FileSpreadsheet, ArrowRight, Clock, CheckCircle2 } from "lucide-react";
import { cn, formatINR } from "@/lib/utils";

const GST_RETURNS = [
  {
    id: "gstr1",
    label: "GSTR-1",
    description: "Outward supplies (sales)",
    href: "/compliance/gst/gstr1",
    frequency: "Monthly",
    dueDate: "11th of next month",
    status: "pending" as ReturnStatus,
  },
  {
    id: "gstr3b",
    label: "GSTR-3B",
    description: "Summary return with tax payment",
    href: "/compliance/gst/gstr3b",
    frequency: "Monthly",
    dueDate: "20th of next month",
    status: "upcoming" as ReturnStatus,
  },
  {
    id: "gstr2b",
    label: "GSTR-2B",
    description: "Auto-drafted ITC statement",
    href: "/compliance/gst/gstr2b",
    frequency: "Monthly (auto)",
    dueDate: "14th of next month",
    status: "upcoming" as ReturnStatus,
  },
  {
    id: "gstr9",
    label: "GSTR-9",
    description: "Annual return",
    href: "/compliance/gst/gstr9",
    frequency: "Annual",
    dueDate: "31st December",
    status: "upcoming" as ReturnStatus,
  },
];

type ReturnStatus = "pending" | "upcoming" | "filed";

const STATUS_STYLES: Record<ReturnStatus, string> = {
  pending: "bg-[var(--warning-muted)] text-[var(--warning)]",
  upcoming: "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
  filed: "bg-[var(--success-muted)] text-[var(--success)]",
};

const STATUS_LABEL: Record<ReturnStatus, string> = {
  pending: "Action Required",
  upcoming: "Upcoming",
  filed: "Filed",
};

export default function GSTReturnsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-heading-xl text-foreground">GST Returns</h1>
        <p className="text-body-sm text-foreground-muted mt-1">
          File and manage GST returns for all registered GSTINs
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {GST_RETURNS.map((ret) => (
          <Link
            key={ret.id}
            href={ret.href}
            className={cn(
              "bg-app-card border border-app-border rounded-card p-5",
              "hover:border-app-accent transition-colors duration-150 group",
            )}
          >
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-[var(--accent-muted)] flex items-center justify-center">
                  <FileSpreadsheet className="w-5 h-5 text-[var(--accent)]" />
                </div>
                <div>
                  <div className="text-sm font-bold text-foreground">{ret.label}</div>
                  <div className="text-xs text-foreground-muted">{ret.description}</div>
                </div>
              </div>
              <ArrowRight className="w-4 h-4 text-foreground-muted group-hover:text-[var(--accent)] transition-colors" />
            </div>
            <div className="flex items-center gap-4 text-xs text-foreground-muted">
              <span className="flex items-center gap-1">
                <Clock className="w-3 h-3" />
                {ret.frequency}
              </span>
              <span>Due: {ret.dueDate}</span>
            </div>
            <div className="mt-3">
              <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold uppercase", STATUS_STYLES[ret.status])}>
                {STATUS_LABEL[ret.status]}
              </span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}
