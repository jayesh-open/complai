"use client";

import { useState, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { FilingConfirmationModal } from "@complai/ui-components";
import { ArrowLeft, ArrowRight, Save, Trash2, Loader2 } from "lucide-react";

import { GSTR9_STEPS, type GSTR9Step, type GSTR9Data } from "../../../types";
import { generateGSTR9Data } from "../../../mock-data";
import {
  StepBar, StepThreshold, StepReviewTables, StepLateITC,
  StepHSN, StepFeesDemands, StepSubmit, StepAcknowledge,
} from "./components/WizardSteps";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

export default function GSTR9WizardPage() {
  const params = useParams<{ fy: string; gstin: string }>();
  const fy = params.fy ?? "2025-26";
  const gstin = params.gstin ?? "29AABCA1234A1Z5";

  const [step, setStep] = useState<GSTR9Step>("threshold");
  const [data] = useState<GSTR9Data>(() => generateGSTR9Data(gstin, fy));
  const [isDirty, setIsDirty] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const [loading, setLoading] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [arn, setArn] = useState<string | null>(null);
  const [signed, setSigned] = useState(false);

  const stepIndex = GSTR9_STEPS.findIndex((s) => s.id === step);
  const isLast = step === "acknowledge";

  const handleNext = useCallback(() => {
    if (step === "submit") {
      setShowConfirm(true);
      return;
    }
    if (stepIndex < GSTR9_STEPS.length - 1) {
      setStep(GSTR9_STEPS[stepIndex + 1].id);
    }
  }, [step, stepIndex]);

  const handlePrev = useCallback(() => {
    if (stepIndex > 0) setStep(GSTR9_STEPS[stepIndex - 1].id);
  }, [stepIndex]);

  const handleFile = useCallback(async () => {
    setShowConfirm(false);
    setLoading(true);
    await new Promise((r) => setTimeout(r, 2000));
    setArn(`AA${gstin.slice(0, 2)}09260001${Math.floor(Math.random() * 9000 + 1000)}`);
    setLoading(false);
    setIsDirty(false);
    setStep("acknowledge");
  }, [gstin]);

  const totalPayable = data.feesAndDemands
    .filter((f) => f.category === "demand" || f.category === "late_fee")
    .reduce((s, f) => s + f.amount, 0);

  return (
    <div className="space-y-4" data-testid="gstr9-wizard">
      <div className="flex items-center gap-3">
        <Link href="/compliance/gst-returns/annual" className="text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div>
          <h1 className="text-heading-xl text-[var(--text-primary)]">GSTR-9 Annual Return</h1>
          <p className="text-body-sm text-[var(--text-muted)]">{gstin} &middot; FY {fy} &middot; {data.legalName}</p>
        </div>
      </div>

      <StepBar current={step} />

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
        {step === "threshold" && <StepThreshold data={data} />}
        {step === "review-tables" && <StepReviewTables data={data} />}
        {step === "late-itc" && <StepLateITC data={data} />}
        {step === "hsn-summary" && <StepHSN data={data} />}
        {step === "fees-demands" && <StepFeesDemands data={data} />}
        {step === "submit" && <StepSubmit data={data} signed={signed} onSigned={() => setSigned(true)} />}
        {step === "acknowledge" && arn && <StepAcknowledge data={data} arn={arn} />}
        {step === "acknowledge" && loading && (
          <div className="flex items-center justify-center gap-2 py-12 text-sm text-[var(--text-muted)]">
            <Loader2 className="w-5 h-5 animate-spin text-[var(--accent)]" />
            Filing with GSTN...
          </div>
        )}
      </div>

      {!isLast && (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={() => { setLastSaved(new Date()); setIsDirty(false); }}
              disabled={!isDirty || loading}
              className={cn(
                "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] transition-colors",
                isDirty ? "text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]" : "text-[var(--text-muted)] cursor-not-allowed opacity-50",
              )}
            >
              <Save className="w-3.5 h-3.5" />Save Draft
            </button>
            <button
              onClick={() => setIsDirty(false)}
              disabled={!isDirty || loading}
              className={cn(
                "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors",
                isDirty ? "text-[var(--danger)] hover:bg-[var(--danger)]/5" : "text-[var(--text-muted)] cursor-not-allowed opacity-50",
              )}
            >
              <Trash2 className="w-3.5 h-3.5" />Discard
            </button>
            {lastSaved && (
              <span className="text-[10px] text-[var(--text-muted)]">
                Saved {Math.max(0, Math.round((Date.now() - lastSaved.getTime()) / 1000))}s ago
              </span>
            )}
          </div>
          <div className="flex items-center gap-2">
            {stepIndex > 0 && (
              <button
                onClick={handlePrev}
                disabled={loading}
                className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors"
              >
                <ArrowLeft className="w-3.5 h-3.5" />Previous
              </button>
            )}
            <button
              onClick={handleNext}
              disabled={loading || (step === "submit" && !signed)}
              className={cn(
                "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors",
                step === "submit"
                  ? "bg-[var(--danger)] text-white hover:opacity-90 disabled:opacity-50"
                  : "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
              )}
            >
              {step === "submit" ? "File GSTR-9" : "Next"}
              {step !== "submit" && <ArrowRight className="w-3.5 h-3.5" />}
            </button>
          </div>
        </div>
      )}

      <FilingConfirmationModal
        open={showConfirm}
        onClose={() => setShowConfirm(false)}
        onConfirm={handleFile}
        title={`File GSTR-9 — FY ${fy}`}
        details={[
          { label: "GSTIN", value: gstin },
          { label: "FY", value: fy },
          { label: "Turnover", value: formatINR(data.turnover) },
          { label: "Tax Payable", value: formatINR(totalPayable) },
        ]}
        warningText="This action is irreversible. GSTR-9 cannot be revised once filed."
        confirmWord="FILE"
      />
    </div>
  );
}
