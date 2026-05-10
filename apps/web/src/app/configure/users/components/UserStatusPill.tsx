"use client";

import { cn } from "@/lib/utils";

interface UserStatusPillProps {
  status: "active" | "inactive";
}

export function UserStatusPill({ status }: UserStatusPillProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded-full",
        status === "active"
          ? "bg-[var(--success)]/15 text-[var(--success)]"
          : "bg-[var(--text-muted)]/10 text-[var(--text-muted)]",
      )}
    >
      <span
        className={cn(
          "w-1.5 h-1.5 rounded-full",
          status === "active" ? "bg-[var(--success)]" : "bg-[var(--text-muted)]",
        )}
      />
      {status === "active" ? "Active" : "Inactive"}
    </span>
  );
}
