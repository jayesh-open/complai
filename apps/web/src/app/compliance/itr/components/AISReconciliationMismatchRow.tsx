"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import type { AISMismatch } from "../types";
import { AlertTriangle, CheckCircle2, Info, XCircle } from "lucide-react";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

const SEVERITY_CONFIG = {
  info: { icon: Info, color: "var(--info)", bg: "bg-[var(--info-muted)]", label: "Info" },
  warn: { icon: AlertTriangle, color: "var(--warning)", bg: "bg-[var(--warning-muted)]", label: "Warning" },
  error: { icon: XCircle, color: "var(--danger)", bg: "bg-[var(--danger-muted)]", label: "Error" },
};

interface AISReconciliationMismatchRowProps {
  mismatch: AISMismatch;
  onResolve?: (id: string, reason: string) => void;
  className?: string;
}

export function AISReconciliationMismatchRow({
  mismatch,
  onResolve,
  className,
}: AISReconciliationMismatchRowProps) {
  const [showResolve, setShowResolve] = useState(false);
  const [reason, setReason] = useState("");
  const config = SEVERITY_CONFIG[mismatch.severity];
  const Icon = config.icon;
  const diff = mismatch.itrValue - mismatch.aisValue;

  return (
    <div
      className={cn(
        "border border-[var(--border-default)] rounded-lg p-3",
        mismatch.resolved && "opacity-60",
        className
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-start gap-2 min-w-0">
          <div
            className={cn("w-5 h-5 rounded flex items-center justify-center flex-shrink-0 mt-0.5", config.bg)}
          >
            <Icon className="w-3 h-3" style={{ color: config.color }} />
          </div>
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <span className="text-xs font-semibold text-[var(--text-primary)]">
                {mismatch.field}
              </span>
              <span
                className={cn(
                  "text-[9px] font-semibold uppercase px-1.5 py-0.5 rounded",
                  config.bg
                )}
                style={{ color: config.color }}
              >
                {config.label}
              </span>
              <span className="text-[10px] font-mono text-[var(--text-muted)]">
                {mismatch.category}
              </span>
            </div>
            <div className="flex items-center gap-3 mt-1 text-[11px]">
              <span className="text-[var(--text-muted)]">
                ITR: <span className="font-medium text-[var(--text-primary)]">{formatINR(mismatch.itrValue)}</span>
              </span>
              <span className="text-[var(--text-muted)]">
                AIS: <span className="font-medium text-[var(--text-primary)]">{formatINR(mismatch.aisValue)}</span>
              </span>
              <span className={cn("font-medium", diff > 0 ? "text-[var(--success)]" : "text-[var(--danger)]")}>
                {diff > 0 ? "+" : ""}{formatINR(diff)}
              </span>
            </div>
          </div>
        </div>

        {mismatch.resolved ? (
          <div className="flex items-center gap-1 text-[10px] text-[var(--success)] font-medium flex-shrink-0">
            <CheckCircle2 className="w-3.5 h-3.5" />
            Resolved
          </div>
        ) : (
          <button
            onClick={() => setShowResolve(!showResolve)}
            className="text-[10px] font-medium px-2 py-1 rounded border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors flex-shrink-0"
          >
            Resolve
          </button>
        )}
      </div>

      {mismatch.resolved && mismatch.resolution && (
        <p className="text-[10px] text-[var(--text-muted)] mt-1.5 ml-7 italic">{mismatch.resolution}</p>
      )}

      {showResolve && !mismatch.resolved && (
        <div className="mt-2 ml-7 flex items-center gap-2">
          <input
            type="text"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Resolution reason…"
            className={cn(
              "flex-1 px-2 py-1.5 text-[11px] rounded-lg border",
              "bg-[var(--bg-tertiary)] border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)]"
            )}
          />
          <button
            onClick={() => { if (reason.trim()) { onResolve?.(mismatch.id, reason); setShowResolve(false); setReason(""); } }}
            disabled={!reason.trim()}
            className={cn(
              "px-3 py-1.5 text-[10px] font-semibold rounded-lg transition-colors",
              reason.trim()
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
            )}
          >
            Confirm
          </button>
        </div>
      )}
    </div>
  );
}
