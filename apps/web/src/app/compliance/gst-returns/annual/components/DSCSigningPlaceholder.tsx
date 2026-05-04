"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";
import { ShieldCheck, Usb, Loader2, CheckCircle2 } from "lucide-react";

type DSCState = "idle" | "detecting" | "signing" | "signed";

interface DSCSigningPlaceholderProps {
  onSigned?: () => void;
  label?: string;
  className?: string;
}

export function DSCSigningPlaceholder({
  onSigned,
  label = "Sign with DSC",
  className,
}: DSCSigningPlaceholderProps) {
  const [state, setState] = useState<DSCState>("idle");

  function handleSign() {
    setState("detecting");
    setTimeout(() => {
      setState("signing");
      setTimeout(() => {
        setState("signed");
        onSigned?.();
      }, 1200);
    }, 1000);
  }

  if (state === "signed") {
    return (
      <div className={cn("flex items-center gap-2 px-4 py-3 rounded-xl bg-[var(--success-muted)] border border-[var(--success)]", className)}>
        <CheckCircle2 className="w-5 h-5 text-[var(--success)]" />
        <div>
          <p className="text-xs font-semibold text-[var(--success)]">Signed with DSC</p>
          <p className="text-[10px] text-[var(--text-muted)]">Digital Signature Certificate applied</p>
        </div>
      </div>
    );
  }

  return (
    <div className={cn("border border-[var(--border-default)] rounded-xl p-5", className)}>
      <div className="flex items-center gap-3 mb-3">
        <div className="w-10 h-10 rounded-lg bg-[var(--accent-muted)] flex items-center justify-center">
          <Usb className="w-5 h-5 text-[var(--accent)]" />
        </div>
        <div>
          <h3 className="text-xs font-semibold text-[var(--text-primary)]">Digital Signature (DSC)</h3>
          <p className="text-[10px] text-[var(--text-muted)]">
            Connect your USB DSC token to sign the return
          </p>
        </div>
      </div>

      {state === "idle" && (
        <button
          onClick={handleSign}
          className={cn(
            "w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg text-xs font-semibold",
            "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors"
          )}
        >
          <ShieldCheck className="w-4 h-4" />
          {label}
        </button>
      )}

      {state === "detecting" && (
        <div className="flex items-center gap-2 justify-center py-3 text-xs text-[var(--text-muted)]">
          <Loader2 className="w-4 h-4 animate-spin text-[var(--accent)]" />
          Detecting DSC token…
        </div>
      )}

      {state === "signing" && (
        <div className="flex items-center gap-2 justify-center py-3 text-xs text-[var(--accent)]">
          <Loader2 className="w-4 h-4 animate-spin" />
          Signing document…
        </div>
      )}
    </div>
  );
}
