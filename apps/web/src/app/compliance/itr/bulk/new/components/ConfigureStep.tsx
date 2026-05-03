"use client";

import { cn } from "@/lib/utils";
import { FileText } from "lucide-react";

interface ConfigureStepProps {
  batchName: string;
  setBatchName: (v: string) => void;
  autoFetchAIS: boolean;
  setAutoFetchAIS: (v: boolean) => void;
  sendMagicLinks: boolean;
  setSendMagicLinks: (v: boolean) => void;
  selectedCount: number;
  onBack: () => void;
  onNext: () => void;
}

export function ConfigureStep({ batchName, setBatchName, autoFetchAIS, setAutoFetchAIS, sendMagicLinks, setSendMagicLinks, selectedCount, onBack, onNext }: ConfigureStepProps) {
  return (
    <div className="space-y-6">
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6 space-y-5">
        <div>
          <label className="block text-xs font-semibold text-[var(--text-primary)] mb-2">Batch Name</label>
          <input
            type="text"
            value={batchName}
            onChange={(e) => setBatchName(e.target.value)}
            placeholder="e.g., Tax Year 2026-27 — Engineering Team"
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs",
              "bg-[var(--bg-primary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]"
            )}
          />
        </div>
        <div className="space-y-3">
          <label className="block text-xs font-semibold text-[var(--text-primary)]">Options</label>
          <label className="flex items-center gap-3 cursor-pointer">
            <input type="checkbox" checked={autoFetchAIS} onChange={(e) => setAutoFetchAIS(e.target.checked)} className="rounded border-[var(--border-default)]" />
            <div>
              <p className="text-xs text-[var(--text-primary)]">Auto-fetch AIS from IT portal</p>
              <p className="text-[10px] text-[var(--text-muted)]">Pull Annual Information Statement for each employee via Sandbox API</p>
            </div>
          </label>
          <label className="flex items-center gap-3 cursor-pointer">
            <input type="checkbox" checked={sendMagicLinks} onChange={(e) => setSendMagicLinks(e.target.checked)} className="rounded border-[var(--border-default)]" />
            <div>
              <p className="text-xs text-[var(--text-primary)]">Send review magic links to employees</p>
              <p className="text-[10px] text-[var(--text-muted)]">Each employee receives a secure link to review and approve their ITR before filing</p>
            </div>
          </label>
        </div>
      </div>
      <div className="flex items-center justify-between">
        <button onClick={onBack} className="text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] font-medium">← Back</button>
        <button
          disabled={!batchName.trim()}
          onClick={onNext}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
            batchName.trim() ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] cursor-not-allowed"
          )}
        >
          <FileText className="w-3.5 h-3.5" />
          Review
        </button>
      </div>
    </div>
  );
}
