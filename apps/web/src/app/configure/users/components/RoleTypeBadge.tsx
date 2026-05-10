"use client";

import { cn } from "@/lib/utils";

interface RoleTypeBadgeProps {
  isSystem: boolean;
}

export function RoleTypeBadge({ isSystem }: RoleTypeBadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded-full",
        isSystem
          ? "bg-[var(--success)]/15 text-[var(--success)]"
          : "bg-[var(--text-muted)]/10 text-[var(--text-muted)]",
      )}
    >
      {isSystem ? "System" : "Custom"}
    </span>
  );
}
