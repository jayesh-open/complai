"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import { MessageSquare } from "lucide-react";

interface RequestChangesModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (feedback: string) => void;
  employeeName: string;
}

export function RequestChangesModal({
  open,
  onClose,
  onSubmit,
  employeeName,
}: RequestChangesModalProps) {
  const [feedback, setFeedback] = useState("");

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-[var(--bg-overlay)] backdrop-blur-sm" onClick={onClose} />
      <div
        className={cn(
          "relative bg-[var(--bg-secondary)] border border-[var(--border-default)]",
          "rounded-2xl shadow-[var(--shadow-lg)] w-full max-w-lg mx-4",
          "animate-in fade-in zoom-in-95 duration-150"
        )}
        role="dialog"
        aria-modal="true"
        aria-labelledby="request-changes-title"
      >
        <div className="px-6 py-4 border-b border-[var(--border-default)] flex items-center gap-2">
          <MessageSquare className="w-4 h-4 text-[var(--accent)]" />
          <h2 id="request-changes-title" className="text-sm font-bold text-[var(--text-primary)]">
            Request Changes
          </h2>
        </div>
        <div className="px-6 py-5 space-y-4">
          <p className="text-[13px] text-[var(--text-secondary)]">
            Describe what needs to be corrected in your ITR. Your HR team ({employeeName}&apos;s employer) will review and update.
          </p>
          <textarea
            value={feedback}
            onChange={(e) => setFeedback(e.target.value)}
            placeholder="e.g., My HRA exemption amount seems incorrect — I have rent receipts showing ₹15,000/month, not ₹10,000…"
            rows={5}
            className={cn(
              "w-full px-3 py-2.5 rounded-lg border text-sm",
              "bg-[var(--bg-tertiary)] border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
              "resize-none"
            )}
            autoFocus
          />
          <p className="text-[10px] text-[var(--text-muted)]">
            Your feedback will be sent to your employer&apos;s compliance team.
          </p>
        </div>
        <div className="px-6 py-4 border-t border-[var(--border-default)] flex justify-end gap-3">
          <button
            onClick={onClose}
            className="px-4 py-2 text-xs font-medium rounded-lg border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={() => { if (feedback.trim()) { onSubmit(feedback); setFeedback(""); } }}
            disabled={!feedback.trim()}
            className={cn(
              "px-4 py-2 text-xs font-semibold rounded-lg transition-colors",
              feedback.trim()
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
            )}
          >
            Submit Feedback
          </button>
        </div>
      </div>
    </div>
  );
}
