"use client";

import { Truck } from "lucide-react";
import { cn } from "@/lib/utils";
import type { VehicleEntry } from "../types";

interface VehicleUpdateTimelineProps {
  entries: VehicleEntry[];
  className?: string;
}

function fmtDate(iso: string): string {
  const d = new Date(iso);
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const min = String(d.getMinutes()).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()} ${hh}:${min}`;
}

export function VehicleUpdateTimeline({ entries, className }: VehicleUpdateTimelineProps) {
  if (entries.length === 0) return null;

  return (
    <div data-testid="vehicle-timeline" className={cn("space-y-0", className)}>
      {entries.map((entry, i) => (
        <div key={`${entry.vehicleNo}-${i}`} className="flex gap-3">
          <div className="flex flex-col items-center">
            <div className={cn(
              "w-6 h-6 rounded-full flex items-center justify-center mt-0.5 flex-shrink-0",
              i === entries.length - 1
                ? "bg-[var(--accent-muted)]"
                : "bg-[var(--bg-tertiary)]",
            )}>
              <Truck className={cn(
                "w-3 h-3",
                i === entries.length - 1 ? "text-[var(--accent)]" : "text-[var(--text-muted)]",
              )} />
            </div>
            {i < entries.length - 1 && (
              <div className="w-px flex-1 bg-[var(--border-default)] my-1" />
            )}
          </div>
          <div className="pb-4 min-w-0">
            <div className="flex items-baseline gap-3">
              <span className="text-xs font-mono font-semibold text-[var(--text-primary)]">
                {entry.vehicleNo}
              </span>
              <span className="text-[10px] text-[var(--text-muted)] tabular-nums">
                {fmtDate(entry.updatedAt)}
              </span>
            </div>
            <div className="text-[11px] text-[var(--text-muted)] mt-0.5">
              by {entry.updatedBy}
            </div>
            {entry.reason && (
              <div className="text-[11px] text-[var(--text-secondary)] mt-0.5">
                {entry.reason}
              </div>
            )}
            {entry.fromPlace && (
              <div className="text-[10px] text-[var(--text-muted)] mt-0.5">
                From: {entry.fromPlace}
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
