"use client";

import { cn } from "@/lib/utils";
import type { ITRForm } from "../types";

const FORM_STYLES: Record<ITRForm, string> = {
  "ITR-1": "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]",
  "ITR-2": "bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-border)]",
  "ITR-3": "bg-[var(--purple-muted)] text-[var(--purple)] border-[var(--purple-muted)]",
  "ITR-4": "bg-[var(--orange-muted)] text-[var(--orange)] border-[var(--orange-muted)]",
  "ITR-5": "bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning-muted)]",
  "ITR-6": "bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-muted)]",
  "ITR-7": "bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-muted)]",
};

interface ITRFormRecommendationPillProps {
  form: ITRForm;
  className?: string;
}

export function ITRFormRecommendationPill({ form, className }: ITRFormRecommendationPillProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border",
        FORM_STYLES[form],
        className
      )}
    >
      {form}
    </span>
  );
}
