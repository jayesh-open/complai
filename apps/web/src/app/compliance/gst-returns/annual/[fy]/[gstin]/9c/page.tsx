"use client";

import { useState, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { ArrowLeft, ArrowRight } from "lucide-react";

import {
  GSTR9C_STEPS,
  type GSTR9CStep,
  type GSTR9CData,
  type GSTR9CMismatch,
} from "../../../types";
import { generateGSTR9CData } from "../../../mock-data";
import { MismatchResolveModal } from "../../../components/MismatchResolveModal";
import { SelfCertificationLock } from "../../../components/SelfCertificationLock";
import {
  StepThresholdCheck,
  StepUploadFinancials,
  StepReconciliation,
  StepResolveMismatches,
  StepFileDSC,
} from "./components/WizardSteps9C";

function StepBar({ current }: { current: GSTR9CStep }) {
  const idx = GSTR9C_STEPS.findIndex((s) => s.id === current);
  return (
    <div className="flex items-center gap-1" data-testid="step-bar-9c">
      {GSTR9C_STEPS.map((s, i) => (
        <div key={s.id} className="flex items-center gap-1">
          <div className={cn(
            "w-7 h-7 rounded-full flex items-center justify-center text-[10px] font-bold border-2 transition-colors",
            i < idx && "bg-[var(--accent)] border-[var(--accent)] text-white",
            i === idx && "border-[var(--accent)] text-[var(--accent)] bg-[var(--accent-muted)]",
            i > idx && "border-[var(--border-default)] text-[var(--text-muted)]",
          )}>
            {i < idx ? "✓" : s.number}
          </div>
          {i < GSTR9C_STEPS.length - 1 && (
            <div className={cn("w-8 h-0.5 rounded-full", i < idx ? "bg-[var(--accent)]" : "bg-[var(--border-default)]")} />
          )}
        </div>
      ))}
    </div>
  );
}

export default function GSTR9CWizardPage() {
  const params = useParams<{ fy: string; gstin: string }>();
  const fy = params.fy ?? "2025-26";
  const gstin = params.gstin ?? "29AABCA1234A1Z5";

  const [step, setStep] = useState<GSTR9CStep>("threshold-check");
  const [data] = useState<GSTR9CData>(() => generateGSTR9CData(gstin, fy));
  const [mismatches, setMismatches] = useState<GSTR9CMismatch[]>(data.mismatches);
  const [financialsUploaded, setFinancialsUploaded] = useState(false);
  const [certified, setCertified] = useState(false);
  const [signed, setSigned] = useState(false);
  const [arn, setArn] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [resolveTarget, setResolveTarget] = useState<GSTR9CMismatch | null>(null);

  const stepIndex = GSTR9C_STEPS.findIndex((s) => s.id === step);
  const unresolvedErrors = mismatches.filter((m) => m.severity === "ERROR" && !m.resolved).length;
  const resolvedCount = mismatches.filter((m) => m.resolved).length;

  const handleResolve = useCallback((id: string, reason: string) => {
    setMismatches((prev) =>
      prev.map((m) => (m.id === id ? { ...m, resolved: true, resolvedReason: reason } : m)),
    );
  }, []);

  const handleNext = useCallback(() => {
    if (stepIndex < GSTR9C_STEPS.length - 1) {
      setStep(GSTR9C_STEPS[stepIndex + 1].id);
    }
  }, [stepIndex]);

  const handlePrev = useCallback(() => {
    if (stepIndex > 0) setStep(GSTR9C_STEPS[stepIndex - 1].id);
  }, [stepIndex]);

  const handleFile = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 2000));
    setArn(`AC${gstin.slice(0, 2)}09260001${Math.floor(Math.random() * 9000 + 1000)}`);
    setLoading(false);
  }, [gstin]);

  const canProceed = (): boolean => {
    switch (step) {
      case "threshold-check": return true;
      case "upload-financials": return financialsUploaded;
      case "reconciliation": return true;
      case "resolve-mismatches": return unresolvedErrors === 0;
      case "self-certification": return certified;
      case "file-dsc": return signed;
      default: return false;
    }
  };

  return (
    <div className="space-y-4" data-testid="gstr9c-wizard">
      <div className="flex items-center gap-3">
        <Link href="/compliance/gst-returns/annual" className="text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div>
          <h1 className="text-heading-xl text-[var(--text-primary)]">GSTR-9C Reconciliation</h1>
          <p className="text-body-sm text-[var(--text-muted)]">{gstin} &middot; FY {fy} &middot; {data.legalName}</p>
        </div>
      </div>

      <StepBar current={step} />

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
        {step === "threshold-check" && (
          <StepThresholdCheck turnover={data.gstr9Turnover} required={data.gstr9cRequired} />
        )}
        {step === "upload-financials" && (
          <StepUploadFinancials audited={data.audited} uploaded={financialsUploaded} onUpload={() => setFinancialsUploaded(true)} />
        )}
        {step === "reconciliation" && (
          <StepReconciliation mismatches={mismatches} />
        )}
        {step === "resolve-mismatches" && (
          <StepResolveMismatches
            mismatches={mismatches}
            unresolvedErrors={unresolvedErrors}
            onResolve={(id) => setResolveTarget(mismatches.find((m) => m.id === id) ?? null)}
          />
        )}
        {step === "self-certification" && (
          <SelfCertificationLock
            gstin={gstin}
            fy={fy}
            totalMismatches={mismatches.length}
            resolvedCount={resolvedCount}
            onCertify={() => setCertified(true)}
            locked={certified}
          />
        )}
        {step === "file-dsc" && (
          <StepFileDSC signed={signed} onSigned={() => setSigned(true)} arn={arn} loading={loading} />
        )}
      </div>

      {!arn && (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl px-4 py-3 flex items-center justify-between">
          <div />
          <div className="flex items-center gap-2">
            {stepIndex > 0 && (
              <button onClick={handlePrev} className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors">
                <ArrowLeft className="w-3.5 h-3.5" />Previous
              </button>
            )}
            <button
              onClick={step === "file-dsc" ? handleFile : handleNext}
              disabled={!canProceed() || loading}
              className={cn(
                "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors",
                step === "file-dsc"
                  ? "bg-[var(--danger)] text-white hover:opacity-90 disabled:opacity-50"
                  : "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)] disabled:opacity-50",
              )}
            >
              {step === "file-dsc" ? "File GSTR-9C" : "Next"}
              {step !== "file-dsc" && <ArrowRight className="w-3.5 h-3.5" />}
            </button>
          </div>
        </div>
      )}

      <MismatchResolveModal
        mismatch={resolveTarget}
        open={resolveTarget !== null}
        onClose={() => setResolveTarget(null)}
        onResolve={handleResolve}
      />
    </div>
  );
}
