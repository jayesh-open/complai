"use client";

import { cn } from "@/lib/utils";

interface ComplianceStatCardProps {
  label: string;
  count: number;
  variant: "success" | "warning" | "neutral";
}

const VARIANT_STYLES = {
  success: {
    bg: "bg-[#ECFDF5]",
    text: "text-[#0F6E56]",
    countText: "text-[#065F46]",
    border: "border-[#A7F3D0]",
  },
  warning: {
    bg: "bg-[#FFFBEB]",
    text: "text-[#BA7517]",
    countText: "text-[#92400E]",
    border: "border-[#FDE68A]",
  },
  neutral: {
    bg: "bg-[var(--bg-secondary)]",
    text: "text-[var(--text-muted)]",
    countText: "text-[var(--text-primary)]",
    border: "border-[var(--border-default)]",
  },
};

export function ComplianceStatCard({ label, count, variant }: ComplianceStatCardProps) {
  const s = VARIANT_STYLES[variant];

  return (
    <div className={cn("flex items-center gap-3 rounded-xl border px-4 py-3", s.bg, s.border)}>
      <span className={cn("text-[22px] font-bold tabular-nums", s.countText)}>{count}</span>
      <span className={cn("text-body-sm", s.text)}>{label}</span>
    </div>
  );
}
