"use client";

import { ArrowLeft, ArrowRight, Save, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";
import type { WizardStep } from "./types";

interface WizardFooterProps {
  currentStep: WizardStep;
  lastSaved: Date | null;
  isDirty: boolean;
  loading: boolean;
  onSaveDraft: () => void;
  onDiscard: () => void;
  onPrevious: () => void;
  onNext: () => void;
}

export function WizardFooter({
  currentStep,
  lastSaved,
  isDirty,
  loading,
  onSaveDraft,
  onDiscard,
  onPrevious,
  onNext,
}: WizardFooterProps) {
  const isFirst = currentStep === "auto-populate";
  const isFiling = currentStep === "file";
  const isLast = currentStep === "acknowledge";

  const savedAgo = lastSaved
    ? `Saved ${Math.max(0, Math.round((Date.now() - lastSaved.getTime()) / 1000))}s ago`
    : null;

  if (isLast) return null;

  return (
    <div className="bg-app-card border border-app-border rounded-card px-4 py-3 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <button
          data-testid="save-draft-button"
          onClick={onSaveDraft}
          disabled={!isDirty || loading}
          className={cn(
            "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] transition-colors",
            isDirty
              ? "text-foreground hover:bg-[var(--bg-tertiary)]"
              : "text-[var(--text-muted)] cursor-not-allowed opacity-50",
          )}
        >
          <Save className="w-3.5 h-3.5" />
          Save Draft
        </button>
        <button
          data-testid="discard-button"
          onClick={onDiscard}
          disabled={!isDirty || loading}
          className={cn(
            "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors",
            isDirty
              ? "text-[var(--danger)] hover:bg-[var(--danger)]/5"
              : "text-[var(--text-muted)] cursor-not-allowed opacity-50",
          )}
        >
          <Trash2 className="w-3.5 h-3.5" />
          Discard
        </button>
        {savedAgo && (
          <span className="text-[10px] text-[var(--text-muted)]">{savedAgo}</span>
        )}
      </div>
      <div className="flex items-center gap-2">
        {!isFirst && (
          <button
            data-testid="previous-button"
            onClick={onPrevious}
            disabled={loading}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-foreground-muted hover:text-foreground hover:bg-[var(--bg-tertiary)] transition-colors"
          >
            <ArrowLeft className="w-3.5 h-3.5" />
            Previous
          </button>
        )}
        <button
          data-testid="next-button"
          onClick={onNext}
          disabled={loading}
          className={cn(
            "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors",
            isFiling
              ? "bg-[var(--danger)] text-white hover:opacity-90"
              : "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
          )}
        >
          {isFiling ? "File Now" : "Next"}
          {!isFiling && <ArrowRight className="w-3.5 h-3.5" />}
        </button>
      </div>
    </div>
  );
}
