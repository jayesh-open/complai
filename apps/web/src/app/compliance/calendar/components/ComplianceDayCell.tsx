"use client";

import type { ComplianceEvent } from "../types";
import { ComplianceItemPill } from "./ComplianceItemPill";
import { cn } from "@/lib/utils";

interface ComplianceDayCellProps {
  date: Date;
  events: (ComplianceEvent & { dueDate: Date })[];
  isToday: boolean;
  isCurrentMonth: boolean;
  onClick: () => void;
}

export function ComplianceDayCell({
  date,
  events,
  isToday,
  isCurrentMonth,
  onClick,
}: ComplianceDayCellProps) {
  const day = date.getDate();
  const visible = events.slice(0, 3);
  const overflow = events.length - 3;

  return (
    <button
      type="button"
      onClick={onClick}
      data-testid={`day-cell-${date.getDate()}`}
      className={cn(
        "flex flex-col items-start p-1.5 border-b border-r border-[var(--border-default)] text-left transition-colors duration-100",
        isCurrentMonth
          ? "bg-[var(--bg-primary)]"
          : "bg-[var(--bg-tertiary)]",
        isToday && "bg-[#FAEEDA]",
        !isToday && isCurrentMonth && events.length > 0 && "hover:bg-[var(--bg-secondary)]",
        !isToday && !isCurrentMonth && "hover:bg-[var(--bg-tertiary)]",
      )}
    >
      <div className="flex items-center gap-1 mb-0.5">
        <span
          className={cn(
            "text-xs font-medium tabular-nums w-5 h-5 flex items-center justify-center rounded",
            isToday && "text-[#633806] font-bold",
            !isToday && isCurrentMonth && "text-[var(--text-primary)]",
            !isToday && !isCurrentMonth && "text-[var(--text-disabled)]",
          )}
        >
          {day}
        </span>
        {isToday && (
          <span className="text-[9px] font-semibold text-[#633806] bg-[#F5D799] px-1.5 py-px rounded-full leading-tight" data-testid="today-indicator">
            Today
          </span>
        )}
      </div>

      <div className="flex flex-col gap-0.5 w-full min-h-0 flex-1 overflow-hidden">
        {visible.map((evt) => (
          <ComplianceItemPill key={evt.id} event={evt} />
        ))}
        {overflow > 0 && (
          <span className="text-[9px] text-[var(--text-muted)] font-medium pl-1">
            +{overflow} more
          </span>
        )}
      </div>
    </button>
  );
}
