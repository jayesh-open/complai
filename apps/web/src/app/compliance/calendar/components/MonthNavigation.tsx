"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";

const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

interface MonthNavigationProps {
  year: number;
  month: number;
  onPrev: () => void;
  onNext: () => void;
  onToday: () => void;
}

export function MonthNavigation({ year, month, onPrev, onNext, onToday }: MonthNavigationProps) {
  return (
    <div className="flex items-center gap-3">
      <h2 className="text-display text-[var(--text-primary)] select-none">
        {MONTH_NAMES[month]} {year}
      </h2>
      <div className="flex items-center gap-1 ml-2">
        <button
          onClick={onPrev}
          className={cn(
            "p-1.5 rounded-lg border border-[var(--border-default)]",
            "text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)]",
            "transition-colors duration-150",
          )}
          aria-label="Previous month"
        >
          <ChevronLeft className="w-4 h-4" />
        </button>
        <button
          onClick={onNext}
          className={cn(
            "p-1.5 rounded-lg border border-[var(--border-default)]",
            "text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)]",
            "transition-colors duration-150",
          )}
          aria-label="Next month"
        >
          <ChevronRight className="w-4 h-4" />
        </button>
        <button
          onClick={onToday}
          className={cn(
            "ml-1 px-3 py-1 rounded-lg border border-[var(--border-default)]",
            "text-xs font-medium text-[var(--text-secondary)]",
            "hover:bg-[var(--bg-tertiary)] transition-colors duration-150",
          )}
        >
          Today
        </button>
      </div>
    </div>
  );
}
