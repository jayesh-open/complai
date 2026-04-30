"use client";

import { useState, useCallback } from "react";
import { ArrowLeft, ArrowRight, Save, Trash2 } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { cn, formatINR } from "@complai/ui-components";
import { FilingConfirmationModal } from "@complai/ui-components";

import type { TDSFormType, FilingWizardStep, TDSFilingData } from "../../../types";
import { FILING_STEPS, FORM_LABELS } from "../../../types";
import { generateFilingData } from "../../../mock-data";
import { FilingWizardStepper } from "../../../components/FilingWizardStepper";
import { StepPull } from "../../../components/StepPull";
import { StepValidate } from "../../../components/StepValidate";
import { StepPreview } from "../../../components/StepPreview";
import { StepSubmit } from "../../../components/StepSubmit";
import { StepAcknowledge } from "../../../components/StepAcknowledge";

export default function TDSFilingWizardPage() {
  const params = useParams<{ formType: string; taxYear: string; quarter: string }>();
  const formType = (params.formType ?? "140") as TDSFormType;
  const taxYear = params.taxYear ?? "2026-27";
  const quarter = params.quarter ?? "q1";

  const [currentStep, setCurrentStep] = useState<FilingWizardStep>("pull");
  const [data, setData] = useState<TDSFilingData | null>(null);
  const [loading, setLoading] = useState(false);
  const [isDirty, setIsDirty] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [arn, setArn] = useState<string | null>(null);

  const stepIndex = FILING_STEPS.findIndex((s) => s.id === currentStep);

  const handlePull = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1200));
    setData(generateFilingData(formType, taxYear, quarter));
    setIsDirty(true);
    setLoading(false);
    setCurrentStep("validate");
  }, [formType, taxYear, quarter]);

  const handlePrevious = useCallback(() => {
    if (stepIndex > 0) setCurrentStep(FILING_STEPS[stepIndex - 1].id);
  }, [stepIndex]);

  const handleNext = useCallback(() => {
    if (currentStep === "submit") { setShowConfirmModal(true); return; }
    if (stepIndex < FILING_STEPS.length - 1) setCurrentStep(FILING_STEPS[stepIndex + 1].id);
  }, [currentStep, stepIndex]);

  const handleFile = useCallback(async () => {
    setShowConfirmModal(false);
    setLoading(true);
    await new Promise((r) => setTimeout(r, 2000));
    setArn(`TDS${formType}${quarter.toUpperCase()}${taxYear.replace("-", "")}${Math.floor(Math.random() * 9000 + 1000)}`);
    setLoading(false);
    setIsDirty(false);
    setCurrentStep("acknowledge");
  }, [formType, quarter, taxYear]);

  const handleSaveDraft = useCallback(() => { setLastSaved(new Date()); setIsDirty(false); }, []);
  const handleDiscard = useCallback(() => { setData(null); setCurrentStep("pull"); setIsDirty(false); }, []);

  const hasBlockers = data ? data.dtaaBlockers.length > 0 : false;
  const confirmWord = `FILE ${quarter.toUpperCase()} ${taxYear}`;
  const isFirst = currentStep === "pull";
  const isSubmit = currentStep === "submit";
  const isLast = currentStep === "acknowledge";
  const savedAgo = lastSaved ? `Saved ${Math.max(0, Math.round((Date.now() - lastSaved.getTime()) / 1000))}s ago` : null;

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <Link href="/compliance/tds/file" className="text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div>
          <h1 className="text-heading-xl text-[var(--text-primary)]">{FORM_LABELS[formType]} Filing</h1>
          <p className="text-body-sm text-[var(--text-muted)]">Tax Year {taxYear} · {quarter.toUpperCase()}</p>
        </div>
      </div>

      <FilingWizardStepper currentStep={currentStep} />

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl">
        {currentStep === "pull" && <StepPull formType={formType} taxYear={taxYear} quarter={quarter} data={data} loading={loading} onPull={handlePull} />}
        {currentStep === "validate" && data && <StepValidate formType={formType} data={data} />}
        {currentStep === "preview" && data && <StepPreview data={data} />}
        {currentStep === "submit" && data && <StepSubmit data={data} loading={loading} hasBlockers={hasBlockers} onOpenConfirmModal={() => setShowConfirmModal(true)} />}
        {currentStep === "acknowledge" && data && arn && <StepAcknowledge data={data} arn={arn} />}
      </div>

      {!isLast && (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button onClick={handleSaveDraft} disabled={!isDirty || loading} className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] transition-colors", isDirty ? "text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]" : "text-[var(--text-muted)] cursor-not-allowed opacity-50")}>
              <Save className="w-3.5 h-3.5" /> Save Draft
            </button>
            <button onClick={handleDiscard} disabled={!isDirty || loading} className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition-colors", isDirty ? "text-[var(--danger)] hover:bg-[var(--danger)]/5" : "text-[var(--text-muted)] cursor-not-allowed opacity-50")}>
              <Trash2 className="w-3.5 h-3.5" /> Discard
            </button>
            {savedAgo && <span className="text-[10px] text-[var(--text-muted)]">{savedAgo}</span>}
          </div>
          <div className="flex items-center gap-2">
            {!isFirst && (
              <button onClick={handlePrevious} disabled={loading} className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors">
                <ArrowLeft className="w-3.5 h-3.5" /> Previous
              </button>
            )}
            {data && (
              <button onClick={handleNext} disabled={loading || (isSubmit && hasBlockers)} className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors", isSubmit ? "bg-[var(--danger)] text-white hover:opacity-90" : "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]", (loading || (isSubmit && hasBlockers)) && "opacity-50 cursor-not-allowed")}>
                {isSubmit ? "File Now" : "Next"} {!isSubmit && <ArrowRight className="w-3.5 h-3.5" />}
              </button>
            )}
          </div>
        </div>
      )}

      {data && (
        <FilingConfirmationModal
          open={showConfirmModal}
          onClose={() => setShowConfirmModal(false)}
          onConfirm={handleFile}
          title={`File ${FORM_LABELS[formType]} — ${quarter.toUpperCase()} ${taxYear}`}
          details={[
            { label: "Form", value: FORM_LABELS[formType] },
            { label: "TAN", value: data.tan },
            { label: "Tax Year", value: taxYear },
            { label: "Quarter", value: data.quarterLabel },
            { label: "Total Tax", value: formatINR(data.totalTax) },
          ]}
          warningText={`This action is irreversible. ${FORM_LABELS[formType]} for ${data.quarterLabel} cannot be revised once filed.`}
          confirmWord={confirmWord}
        />
      )}
    </div>
  );
}
