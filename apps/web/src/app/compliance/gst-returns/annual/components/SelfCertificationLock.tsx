"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import { Lock, ShieldCheck } from "lucide-react";

interface SelfCertificationLockProps {
  gstin: string;
  fy: string;
  totalMismatches: number;
  resolvedCount: number;
  onCertify: () => void;
  locked?: boolean;
  className?: string;
}

export function SelfCertificationLock({
  gstin,
  fy,
  totalMismatches,
  resolvedCount,
  onCertify,
  locked = false,
  className,
}: SelfCertificationLockProps) {
  const [typed, setTyped] = useState("");
  const confirmPhrase = `I CERTIFY GSTR-9C ${fy} ${gstin}`;
  const matches = typed.toUpperCase() === confirmPhrase.toUpperCase();

  if (locked) {
    return (
      <div className={cn("border border-[var(--success)] rounded-xl p-5 bg-[color-mix(in_srgb,var(--success)_4%,transparent)]", className)} data-testid="certification-locked">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-[var(--success)]/10 flex items-center justify-center">
            <Lock className="w-5 h-5 text-[var(--success)]" />
          </div>
          <div>
            <p className="text-sm font-semibold text-[var(--success)]">Self-Certification Complete</p>
            <p className="text-[11px] text-[var(--text-muted)]">
              GSTR-9C for {gstin} &middot; FY {fy} has been certified and locked.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={cn("border border-[var(--border-default)] rounded-xl p-5", className)} data-testid="certification-form">
      <div className="flex items-center gap-3 mb-4">
        <div className="w-10 h-10 rounded-lg bg-[var(--accent-muted)] flex items-center justify-center">
          <ShieldCheck className="w-5 h-5 text-[var(--accent)]" />
        </div>
        <div>
          <h3 className="text-sm font-semibold text-[var(--text-primary)]">Self-Certification</h3>
          <p className="text-[11px] text-[var(--text-muted)]">
            {resolvedCount}/{totalMismatches} mismatches resolved
          </p>
        </div>
      </div>

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-3 mb-4 text-[11px] text-[var(--text-secondary)] leading-relaxed">
        I hereby certify that the reconciliation statement in GSTR-9C for the financial year {fy}
        pertaining to GSTIN {gstin} is true and correct to the best of my knowledge and belief,
        and that the information provided therein is in accordance with the provisions of the
        Central Goods and Services Tax Act, 2017.
      </div>

      <div className="space-y-2">
        <label className="text-[11px] font-medium text-[var(--text-secondary)] block">
          Type &ldquo;<span className="font-mono text-[var(--accent)]">{confirmPhrase}</span>&rdquo; to certify
        </label>
        <input
          type="text"
          value={typed}
          onChange={(e) => setTyped(e.target.value)}
          placeholder={confirmPhrase}
          className={cn(
            "w-full px-3 py-2 rounded-lg border text-xs font-mono",
            "bg-[var(--bg-secondary)] text-[var(--text-primary)]",
            "placeholder:text-[var(--text-muted)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/30",
            "border-[var(--border-default)]",
          )}
        />
        <button
          onClick={onCertify}
          disabled={!matches}
          className={cn(
            "w-full px-4 py-2.5 rounded-lg text-xs font-semibold transition-colors mt-2",
            matches
              ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
              : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] cursor-not-allowed",
          )}
        >
          Certify GSTR-9C
        </button>
      </div>
    </div>
  );
}
