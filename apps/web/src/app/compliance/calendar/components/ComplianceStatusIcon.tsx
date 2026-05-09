"use client";

import type { EventStatus } from "../types";

const STATUS_CONFIG: Record<EventStatus, { color: string; label: string }> = {
  filed: { color: "#0F6E56", label: "Filed" },
  due_soon: { color: "#BA7517", label: "Due soon" },
  upcoming: { color: "#6B7280", label: "Upcoming" },
  overdue: { color: "#DC2626", label: "Overdue" },
};

export function ComplianceStatusIcon({
  status,
  size = 14,
}: {
  status: EventStatus;
  size?: number;
}) {
  const cfg = STATUS_CONFIG[status];

  if (status === "filed") {
    return (
      <svg width={size} height={size} viewBox="0 0 16 16" fill="none" aria-label={cfg.label}>
        <circle cx="8" cy="8" r="7" fill={cfg.color} />
        <path d="M5 8.5L7 10.5L11 6" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    );
  }

  if (status === "due_soon") {
    return (
      <svg width={size} height={size} viewBox="0 0 16 16" fill="none" aria-label={cfg.label}>
        <circle cx="8" cy="8" r="7" fill={cfg.color} />
        <circle cx="8" cy="8" r="2" fill="white" />
      </svg>
    );
  }

  if (status === "overdue") {
    return (
      <svg width={size} height={size} viewBox="0 0 16 16" fill="none" aria-label={cfg.label}>
        <circle cx="8" cy="8" r="7" fill={cfg.color} />
        <path d="M8 5V9" stroke="white" strokeWidth="1.5" strokeLinecap="round" />
        <circle cx="8" cy="11" r="0.75" fill="white" />
      </svg>
    );
  }

  return (
    <svg width={size} height={size} viewBox="0 0 16 16" fill="none" aria-label={cfg.label}>
      <circle cx="8" cy="8" r="6.5" stroke={cfg.color} strokeWidth="1" fill="none" />
    </svg>
  );
}

export function StatusLegend() {
  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center gap-1.5">
        <ComplianceStatusIcon status="filed" size={12} />
        <span className="text-caption text-[var(--text-muted)]">Filed</span>
      </div>
      <div className="flex items-center gap-1.5">
        <ComplianceStatusIcon status="due_soon" size={12} />
        <span className="text-caption text-[var(--text-muted)]">Due in 7 days</span>
      </div>
      <div className="flex items-center gap-1.5">
        <ComplianceStatusIcon status="upcoming" size={12} />
        <span className="text-caption text-[var(--text-muted)]">Upcoming</span>
      </div>
    </div>
  );
}
