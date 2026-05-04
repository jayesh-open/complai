"use client";

import { cn } from "@/lib/utils";
import { CheckCircle2, FileDown, ShieldCheck, Clock } from "lucide-react";

interface EmployeeApprovalConfirmationProps {
  arn: string;
  filedAt: string;
  signMethod: "dsc" | "evc";
  className?: string;
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function EmployeeApprovalConfirmation({
  arn,
  filedAt,
  signMethod,
  className,
}: EmployeeApprovalConfirmationProps) {
  return (
    <div className={cn("space-y-4", className)}>
      {/* Success Banner */}
      <div className="bg-[var(--success-muted)] border-2 border-[var(--success)] rounded-xl p-5">
        <div className="flex items-center gap-2 mb-3">
          <CheckCircle2 className="w-5 h-5 text-[var(--success)]" />
          <h2 className="text-sm font-bold text-[var(--success)]">
            ITR Approved & Submitted
          </h2>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          <div>
            <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">ARN</p>
            <p className="text-sm font-mono font-bold text-[var(--text-primary)] mt-0.5">{arn}</p>
          </div>
          <div>
            <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">Filed On</p>
            <p className="text-xs font-medium text-[var(--text-primary)] mt-0.5">{formatDate(filedAt)}</p>
          </div>
          <div>
            <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">e-Verification</p>
            <span className={cn(
              "inline-flex items-center gap-1 mt-0.5 text-[10px] font-semibold px-2 py-0.5 rounded-md",
              signMethod === "dsc"
                ? "bg-[var(--success-muted)] text-[var(--success)]"
                : "bg-[var(--warning-muted)] text-[var(--warning)]"
            )}>
              <ShieldCheck className="w-3 h-3" />
              {signMethod === "dsc" ? "Verified (DSC)" : "Pending EVC"}
            </span>
          </div>
        </div>
      </div>

      {/* What's Next Panel */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
        <h3 className="text-xs font-semibold text-[var(--text-primary)] mb-3">What&apos;s Next</h3>
        <div className="space-y-3">
          <div className="flex items-start gap-2">
            <FileDown className="w-4 h-4 text-[var(--accent)] flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-xs font-medium text-[var(--text-primary)]">Download ITR-V</p>
              <p className="text-[10px] text-[var(--text-muted)]">
                Your ITR-V acknowledgement is available for download.
              </p>
              <button className="mt-1 text-[11px] font-medium text-[var(--accent)] hover:underline">
                Download ITR-V (PDF)
              </button>
            </div>
          </div>
          <div className="flex items-start gap-2">
            <Clock className="w-4 h-4 text-[var(--warning)] flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-xs font-medium text-[var(--text-primary)]">30-Day Verification Window</p>
              <p className="text-[10px] text-[var(--text-muted)]">
                {signMethod === "evc"
                  ? "Complete e-Verification via Aadhaar OTP, net banking, or DSC within 30 days of filing. Unverified returns are treated as not filed."
                  : "Your return has been verified via DSC. No further action needed."
                }
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
