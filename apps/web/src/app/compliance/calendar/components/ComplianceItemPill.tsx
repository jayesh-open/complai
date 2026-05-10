"use client";

import type { ComplianceEvent } from "../types";
import { ComplianceStatusIcon } from "./ComplianceStatusIcon";
import { getCategoryStyle } from "./ComplianceCategoryBadge";

export function ComplianceItemPill({ event }: { event: ComplianceEvent & { dueDate: Date } }) {
  const cat = getCategoryStyle(event.category);

  return (
    <div
      data-testid="event-pill"
      className="flex items-center gap-1 rounded px-1.5 py-0.5 max-w-full overflow-hidden"
      style={{ backgroundColor: cat.bg }}
    >
      <ComplianceStatusIcon status={event.status} size={10} />
      <span
        className="text-[10px] font-medium truncate leading-tight"
        style={{ color: cat.text }}
      >
        {event.title}
      </span>
    </div>
  );
}
