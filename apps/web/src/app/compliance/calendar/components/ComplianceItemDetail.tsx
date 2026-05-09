"use client";

import Link from "next/link";
import { ExternalLink, Bell, CheckCircle2 } from "lucide-react";
import type { ComplianceEvent } from "../types";
import { ComplianceStatusIcon } from "./ComplianceStatusIcon";
import { ComplianceCategoryBadge } from "./ComplianceCategoryBadge";

function StatusPill({ status }: { status: ComplianceEvent["status"] }) {
  const config: Record<string, { bg: string; text: string; label: string }> = {
    filed: { bg: "#ECFDF5", text: "#065F46", label: "Filed" },
    due_soon: { bg: "#FFFBEB", text: "#92400E", label: "Due Soon" },
    upcoming: { bg: "#F3F4F6", text: "#374151", label: "Upcoming" },
    overdue: { bg: "#FEF2F2", text: "#991B1B", label: "Overdue" },
  };
  const c = config[status];
  return (
    <span
      className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium"
      style={{ backgroundColor: c.bg, color: c.text }}
    >
      <ComplianceStatusIcon status={status} size={10} />
      {c.label}
    </span>
  );
}

interface ComplianceItemDetailProps {
  event: ComplianceEvent & { dueDate: Date };
}

export function ComplianceItemDetail({ event }: ComplianceItemDetailProps) {
  return (
    <div className="border border-[var(--border-default)] rounded-xl p-4 space-y-3">
      <div className="flex items-start justify-between gap-2">
        <h4 className="text-heading-sm text-[var(--text-primary)]">{event.title}</h4>
        <StatusPill status={event.status} />
      </div>

      <ComplianceCategoryBadge category={event.category} size="small" />

      <div className="grid grid-cols-2 gap-x-4 gap-y-2 text-body-sm">
        <div>
          <span className="text-[var(--text-muted)]">Authority</span>
          <p className="text-[var(--text-primary)] font-medium">{event.authority}</p>
        </div>
        <div>
          <span className="text-[var(--text-muted)]">Section / Form</span>
          <p className="text-[var(--text-primary)] font-medium">
            {event.sectionRef || event.formRef || "—"}
          </p>
        </div>
        {event.penalty && (
          <div className="col-span-2">
            <span className="text-[var(--text-muted)]">Penalty if missed</span>
            <p className="text-[#DC2626] font-medium text-xs">{event.penalty}</p>
          </div>
        )}
      </div>

      <p className="text-body-sm text-[var(--text-secondary)]">{event.description}</p>

      <div className="flex items-center gap-2 pt-1">
        {event.linkedModule && (
          <Link
            href={event.linkedModule}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-[var(--accent)] text-[var(--accent-text)] text-xs font-medium hover:opacity-90 transition-opacity"
          >
            Open in module
            <ExternalLink className="w-3 h-3" />
          </Link>
        )}
        <button
          type="button"
          className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--border-default)] text-xs font-medium text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
        >
          <CheckCircle2 className="w-3 h-3" />
          Mark as filed
        </button>
        <button
          type="button"
          className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[var(--border-default)] text-xs font-medium text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
        >
          <Bell className="w-3 h-3" />
          Set reminder
        </button>
      </div>
    </div>
  );
}
