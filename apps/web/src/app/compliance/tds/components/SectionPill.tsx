"use client";

import { cn } from "@complai/ui-components";
import type { Section2025 } from "../types";

const SECTION_STYLES: Record<Section2025, string> = {
  "392": "bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-border)]",
  "393(1)": "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]",
  "393(2)": "bg-[var(--purple-muted)] text-[var(--purple)] border-[var(--purple-muted)]",
  "393(3)": "bg-[var(--orange-muted)] text-[var(--orange)] border-[var(--orange-muted)]",
};

const SECTION_NAMES: Record<Section2025, string> = {
  "392": "392 Salary",
  "393(1)": "393(1)",
  "393(2)": "393(2) NR",
  "393(3)": "393(3) TCS",
};

interface SectionPillProps {
  section: Section2025;
  className?: string;
}

export function SectionPill({ section, className }: SectionPillProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border",
        SECTION_STYLES[section],
        className
      )}
    >
      {SECTION_NAMES[section]}
    </span>
  );
}
