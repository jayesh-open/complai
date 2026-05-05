"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import { X } from "lucide-react";
import type { GSTR9CMismatch } from "../types";
import { MismatchSeverityBadge } from "./MismatchSeverityBadge";

interface MismatchResolveModalProps {
  mismatch: GSTR9CMismatch | null;
  open: boolean;
  onClose: () => void;
  onResolve: (id: string, reason: string) => void;
}

export function MismatchResolveModal({ mismatch, open, onClose, onResolve }: MismatchResolveModalProps) {
  const [reason, setReason] = useState("");

  if (!open || !mismatch) return null;

  const canSubmit = reason.trim().length >= 10;

  function handleSubmit() {
    if (!mismatch || !canSubmit) return;
    onResolve(mismatch.id, reason.trim());
    setReason("");
    onClose();
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" data-testid="resolve-modal">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className="relative bg-[var(--bg-primary)] border border-[var(--border-default)] rounded-xl shadow-lg w-full max-w-md mx-4 p-5">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">Resolve Mismatch</h3>
            <p className="text-[11px] text-[var(--text-muted)] mt-0.5">
              Part {mismatch.section} &middot; {mismatch.category}
            </p>
          </div>
          <button onClick={onClose} className="p-1 rounded hover:bg-[var(--bg-tertiary)] transition-colors">
            <X className="w-4 h-4 text-[var(--text-muted)]" />
          </button>
        </div>

        <div className="space-y-3">
          <div className="flex items-center gap-2">
            <MismatchSeverityBadge severity={mismatch.severity} />
            <span className="text-xs text-[var(--text-secondary)]">{mismatch.description}</span>
          </div>

          <div className="grid grid-cols-3 gap-2 text-[11px] text-[var(--text-muted)] border-t border-[var(--border-default)] pt-3">
            <div>
              <p className="font-medium text-[var(--text-secondary)]">Books</p>
              <p className="tabular-nums">₹{mismatch.booksAmount.toLocaleString("en-IN")}</p>
            </div>
            <div>
              <p className="font-medium text-[var(--text-secondary)]">GSTR-9</p>
              <p className="tabular-nums">₹{mismatch.gstr9Amount.toLocaleString("en-IN")}</p>
            </div>
            <div>
              <p className="font-medium text-[var(--text-secondary)]">Delta</p>
              <p className="tabular-nums text-[var(--danger)]">₹{mismatch.difference.toLocaleString("en-IN")}</p>
            </div>
          </div>

          <div>
            <label className="text-[11px] font-medium text-[var(--text-secondary)] block mb-1">
              Resolution reason (min 10 chars)
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              rows={3}
              placeholder="Explain why this mismatch is acceptable or how it was corrected..."
              className={cn(
                "w-full px-3 py-2 rounded-lg border text-xs text-[var(--text-primary)] bg-[var(--bg-secondary)]",
                "placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/30",
                "border-[var(--border-default)]",
              )}
            />
            <p className="text-[10px] text-[var(--text-muted)] mt-0.5">{reason.length}/10 minimum</p>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <button
              onClick={onClose}
              className="flex-1 px-3 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSubmit}
              disabled={!canSubmit}
              className={cn(
                "flex-1 px-3 py-2 rounded-lg text-xs font-semibold transition-colors",
                canSubmit
                  ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                  : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] cursor-not-allowed",
              )}
            >
              {mismatch.severity === "ERROR" ? "Resolve Error" : "Acknowledge"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
