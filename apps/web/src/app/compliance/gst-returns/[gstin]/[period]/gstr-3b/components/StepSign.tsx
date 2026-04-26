"use client";

import { ShieldCheck, KeyRound } from "lucide-react";
import { cn, formatINR } from "@/lib/utils";
import type { GSTR3BData } from "./types";

interface StepSignProps {
  data: GSTR3BData;
  signMethod: "dsc" | "evc";
  onSignMethodChange: (method: "dsc" | "evc") => void;
}

export function StepSign({ data, signMethod, onSignMethodChange }: StepSignProps) {
  return (
    <div className="p-6 space-y-4" data-testid="step-sign">
      <div>
        <h2 className="text-heading-lg text-foreground">Step 4: Select Signing Method</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Choose how to authenticate your GSTR-3B filing.
        </p>
      </div>

      {/* Organization Details */}
      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)] mb-2">Filing Details</div>
        {[
          { label: "GSTIN", value: data.gstin },
          { label: "Return Type", value: "GSTR-3B" },
          { label: "Period", value: data.periodLabel },
          { label: "Net Payable", value: formatINR(data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst) },
          { label: "Legal Name", value: "Complai Technologies Pvt Ltd" },
          { label: "Authorized Signatory", value: "Jayesh Haridasan" },
        ].map((item) => (
          <div key={item.label} className="flex gap-3 text-xs">
            <span className="text-foreground-muted min-w-[140px]">{item.label}:</span>
            <span className="text-foreground font-medium font-mono">{item.value}</span>
          </div>
        ))}
      </div>

      {/* Sign Method Selection */}
      <div className="grid grid-cols-2 gap-4">
        <button
          data-testid="sign-dsc"
          onClick={() => onSignMethodChange("dsc")}
          className={cn(
            "p-4 rounded-lg border-2 text-left transition-colors",
            signMethod === "dsc"
              ? "border-[var(--accent)] bg-[var(--accent-muted)]"
              : "border-[var(--border-default)] hover:border-[var(--accent)]/50",
          )}
        >
          <div className="flex items-center gap-3 mb-2">
            <div className={cn(
              "w-10 h-10 rounded-lg flex items-center justify-center",
              signMethod === "dsc" ? "bg-[var(--accent)] text-white" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
            )}>
              <KeyRound className="w-5 h-5" />
            </div>
            <div>
              <div className="text-sm font-bold text-foreground">DSC (Digital Signature)</div>
              <div className="text-[10px] text-foreground-muted">Class 2/3 Digital Certificate</div>
            </div>
          </div>
          <p className="text-[11px] text-foreground-muted">
            Sign using your registered Digital Signature Certificate. Required for companies and LLPs.
          </p>
        </button>

        <button
          data-testid="sign-evc"
          onClick={() => onSignMethodChange("evc")}
          className={cn(
            "p-4 rounded-lg border-2 text-left transition-colors",
            signMethod === "evc"
              ? "border-[var(--accent)] bg-[var(--accent-muted)]"
              : "border-[var(--border-default)] hover:border-[var(--accent)]/50",
          )}
        >
          <div className="flex items-center gap-3 mb-2">
            <div className={cn(
              "w-10 h-10 rounded-lg flex items-center justify-center",
              signMethod === "evc" ? "bg-[var(--accent)] text-white" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
            )}>
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <div className="text-sm font-bold text-foreground">EVC (Electronic Verification)</div>
              <div className="text-[10px] text-foreground-muted">OTP to registered mobile/email</div>
            </div>
          </div>
          <p className="text-[11px] text-foreground-muted">
            Verify via OTP sent to registered mobile number or email. Available for proprietors and partners.
          </p>
        </button>
      </div>

      <div className="flex items-center gap-2 p-3 bg-[var(--accent-muted)] rounded-lg">
        <ShieldCheck className="w-4 h-4 text-[var(--accent)]" />
        <span className="text-xs text-[var(--accent)] font-medium">
          Selected: {signMethod === "dsc" ? "Digital Signature Certificate" : "Electronic Verification Code"}
        </span>
      </div>
    </div>
  );
}
