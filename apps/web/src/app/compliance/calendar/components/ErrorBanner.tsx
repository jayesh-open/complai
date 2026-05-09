"use client";

import { useState } from "react";
import { X } from "lucide-react";

export function ErrorBanner() {
  const [dismissed, setDismissed] = useState(false);

  if (dismissed) return null;

  return (
    <div className="flex items-center gap-3 px-4 py-2.5 rounded-lg bg-[#FFFBEB] border border-[#F59E0B]/30 text-sm text-[#92400E]">
      <span className="flex-1">
        Showing sample data — live calendar service unavailable. Retrying automatically.
      </span>
      <button
        type="button"
        onClick={() => setDismissed(true)}
        className="p-1 rounded hover:bg-[#FEF3C7] transition-colors shrink-0"
        aria-label="Dismiss"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
}
