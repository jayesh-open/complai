"use client";

import { useState, useCallback, useEffect, useRef } from "react";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { formatINR } from "@/lib/utils";
import { FilingConfirmationModal } from "@complai/ui-components";

import { STEPS, type WizardStep, type GSTR3BData, generateMockGSTR3B } from "./components/types";
import { StepIndicator } from "./components/StepIndicator";
import { WizardFooter } from "./components/WizardFooter";
import { StepAutoPopulate } from "./components/StepAutoPopulate";
import { StepReview } from "./components/StepReview";
import { StepPay } from "./components/StepPay";
import { StepSign } from "./components/StepSign";
import { StepFile } from "./components/StepFile";
import { StepAcknowledge } from "./components/StepAcknowledge";

const AUTO_SAVE_INTERVAL = 10_000;

export default function GSTR3BWizard() {
  const params = useParams<{ gstin: string; period: string }>();
  const gstin = params.gstin ?? "29AABCA1234A1Z5";
  const period = params.period ?? "2026-04";

  const [currentStep, setCurrentStep] = useState<WizardStep>("auto-populate");
  const [data, setData] = useState<GSTR3BData | null>(null);
  const [loading, setLoading] = useState(false);
  const [isDirty, setIsDirty] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [signMethod, setSignMethod] = useState<"dsc" | "evc">("evc");
  const [arn, setArn] = useState<string | null>(null);

  const autoSaveRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (!isDirty) return;
    autoSaveRef.current = setInterval(() => {
      setLastSaved(new Date());
      setIsDirty(false);
    }, AUTO_SAVE_INTERVAL);
    return () => {
      if (autoSaveRef.current) clearInterval(autoSaveRef.current);
    };
  }, [isDirty]);

  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
      }
    };
    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => window.removeEventListener("beforeunload", handleBeforeUnload);
  }, [isDirty]);

  const stepIndex = STEPS.findIndex((s) => s.id === currentStep);

  const handlePopulate = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1500));
    const mockData = generateMockGSTR3B(gstin, period);
    setData(mockData);
    setIsDirty(true);
    setLoading(false);
    setCurrentStep("review");
  }, [gstin, period]);

  const handleOverride = useCallback((_section: string, _index: number, _field: string, _value: number, _reason: string) => {
    setIsDirty(true);
  }, []);

  const handleSaveDraft = useCallback(() => {
    setLastSaved(new Date());
    setIsDirty(false);
  }, []);

  const handleDiscard = useCallback(() => {
    if (data) {
      setData(generateMockGSTR3B(gstin, period));
      setIsDirty(false);
    }
  }, [data, gstin, period]);

  const handlePrevious = useCallback(() => {
    if (stepIndex > 0) {
      setCurrentStep(STEPS[stepIndex - 1].id);
    }
  }, [stepIndex]);

  const handleNext = useCallback(() => {
    if (currentStep === "file") {
      setShowConfirmModal(true);
      return;
    }
    if (stepIndex < STEPS.length - 1) {
      setCurrentStep(STEPS[stepIndex + 1].id);
    }
  }, [currentStep, stepIndex]);

  const handleFile = useCallback(async () => {
    setShowConfirmModal(false);
    setLoading(true);
    await new Promise((r) => setTimeout(r, 2000));
    setArn(`AA290420260003${Math.floor(Math.random() * 9000 + 1000)}`);
    setLoading(false);
    setIsDirty(false);
    setCurrentStep("acknowledge");
  }, []);

  const totalPayable = data
    ? data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst
    : 0;

  return (
    <div className="space-y-4" data-testid="gstr3b-wizard">
      {/* Header */}
      <div className="flex items-center gap-3">
        <Link
          href="/compliance/gst"
          className="text-foreground-muted hover:text-foreground transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div>
          <h1 className="text-heading-xl text-foreground">GSTR-3B Filing</h1>
          <p className="text-body-sm text-foreground-muted">
            {gstin} &middot; {period}
          </p>
        </div>
      </div>

      <StepIndicator currentStep={currentStep} />

      {/* Step Content */}
      <div className="bg-app-card border border-app-border rounded-card">
        {currentStep === "auto-populate" && (
          <StepAutoPopulate data={data} loading={loading} onPopulate={handlePopulate} />
        )}

        {currentStep === "review" && data && (
          <StepReview data={data} onOverride={handleOverride} />
        )}

        {currentStep === "pay" && data && (
          <StepPay data={data} />
        )}

        {currentStep === "sign" && data && (
          <StepSign data={data} signMethod={signMethod} onSignMethodChange={setSignMethod} />
        )}

        {currentStep === "file" && data && (
          <StepFile data={data} signMethod={signMethod} loading={loading} onOpenConfirmModal={() => setShowConfirmModal(true)} />
        )}

        {currentStep === "acknowledge" && data && arn && (
          <StepAcknowledge data={data} arn={arn} />
        )}
      </div>

      <WizardFooter
        currentStep={currentStep}
        lastSaved={lastSaved}
        isDirty={isDirty}
        loading={loading}
        onSaveDraft={handleSaveDraft}
        onDiscard={handleDiscard}
        onPrevious={handlePrevious}
        onNext={handleNext}
      />

      {data && (
        <FilingConfirmationModal
          open={showConfirmModal}
          onClose={() => setShowConfirmModal(false)}
          onConfirm={handleFile}
          title={`File GSTR-3B — ${data.periodLabel}`}
          details={[
            { label: "GSTIN", value: data.gstin },
            { label: "Period", value: data.periodLabel },
            { label: "Net Payable", value: formatINR(totalPayable) },
            { label: "Sign Method", value: signMethod === "dsc" ? "DSC" : "EVC" },
          ]}
          warningText="This action is irreversible. GSTR-3B cannot be revised once filed. Tax payment will be debited from your cash ledger."
          confirmWord="FILE"
          signMethod={signMethod}
          onSignMethodChange={setSignMethod}
        />
      )}
    </div>
  );
}
