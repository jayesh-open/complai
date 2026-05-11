"use client";

import { Eye, X } from "lucide-react";
import { useViewAsStore } from "@/store/view-as-store";

export function ViewAsBanner() {
  const { isViewAs, viewAsRole, exitViewAs } = useViewAsStore();

  if (!isViewAs || !viewAsRole) return null;

  return (
    <div
      className="bg-[var(--warning)]/15 border-b border-[var(--warning)]/30 px-5 py-2 flex items-center justify-between flex-shrink-0"
      data-testid="view-as-banner"
    >
      <div className="flex items-center gap-2 text-xs text-[var(--text-primary)]">
        <Eye className="w-3.5 h-3.5 text-[var(--warning)]" />
        <span>
          Viewing as: <strong>{viewAsRole.display_name}</strong>. You see only what this role can access.
        </span>
      </div>
      <button
        onClick={exitViewAs}
        data-testid="exit-view-as"
        className="flex items-center gap-1 px-3 py-1 rounded-lg text-xs font-medium border border-[var(--warning)]/40 text-[var(--warning)] hover:bg-[var(--warning)]/10 transition-colors"
      >
        <X className="w-3 h-3" />
        Exit View As
      </button>
    </div>
  );
}
