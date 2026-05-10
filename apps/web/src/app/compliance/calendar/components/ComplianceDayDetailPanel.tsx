"use client";

import { useEffect, useRef } from "react";
import { X } from "lucide-react";
import type { ComplianceEvent } from "../types";
import { ComplianceItemDetail } from "./ComplianceItemDetail";
import { cn } from "@/lib/utils";

const DAY_NAMES = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

interface ComplianceDayDetailPanelProps {
  date: Date;
  events: (ComplianceEvent & { dueDate: Date })[];
  open: boolean;
  onClose: () => void;
}

export function ComplianceDayDetailPanel({
  date,
  events,
  open,
  onClose,
}: ComplianceDayDetailPanelProps) {
  const panelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  const dateStr = `${MONTH_NAMES[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}`;
  const dayName = DAY_NAMES[date.getDay()];

  return (
    <>
      {/* Backdrop */}
      <div
        className={cn(
          "fixed inset-0 z-40 transition-opacity duration-200",
          open ? "bg-black/40 pointer-events-auto" : "bg-transparent pointer-events-none",
        )}
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Panel */}
      <div
        ref={panelRef}
        data-testid="day-detail-panel"
        className={cn(
          "fixed top-0 right-0 z-50 h-screen w-[480px] max-w-[90vw] bg-[var(--bg-primary)]",
          "border-l border-[var(--border-default)] shadow-xl flex flex-col",
          "transition-transform duration-200 ease-out",
          open ? "translate-x-0" : "translate-x-full",
        )}
      >
        {/* Header */}
        <div className="flex items-start justify-between p-5 border-b border-[var(--border-default)]">
          <div>
            <h3 className="text-heading-lg text-[var(--text-primary)]">
              {dateStr} &middot; {dayName}
            </h3>
            <p className="text-body-sm text-[var(--text-muted)] mt-0.5">
              {events.length} {events.length === 1 ? "obligation" : "obligations"}
            </p>
          </div>
          <button
            onClick={onClose}
            data-testid="close-panel"
            className="p-1.5 rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-colors"
            aria-label="Close panel"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto p-5 space-y-3">
          {events.length === 0 && (
            <p className="text-body text-[var(--text-muted)] py-8 text-center">
              No compliance obligations on this date.
            </p>
          )}
          {events.map((evt) => (
            <ComplianceItemDetail key={evt.id} event={evt} />
          ))}
        </div>

        {/* Footer */}
        <div className="border-t border-[var(--border-default)] px-5 py-3">
          <p className="text-caption text-[var(--text-disabled)]">
            Calendar shows obligations applicable to your tenant. To configure, see Settings.
          </p>
        </div>
      </div>
    </>
  );
}
