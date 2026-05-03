"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { ArrowLeft, FileCheck2, ShieldCheck, Copy, ExternalLink } from "lucide-react";
import { cn } from "@/lib/utils";
import { getEmployeeDetail } from "../../../mock-detail-data";
import { ITRFormRecommendationPill } from "../../../components/ITRFormRecommendationPill";
import { RegimeIndicator } from "../../../components/RegimeIndicator";
import { FilingStatusPill } from "../../../components/FilingStatusPill";
import { IncomeHeadCard } from "../../../components/IncomeHeadCard";
import { TaxComputationPanel } from "../../../components/TaxComputationPanel";
import { RefundIndicator } from "../../../components/RefundIndicator";
import { FilingConfirmationModal } from "@complai/ui-components";

type SubmissionState = "ready" | "confirming" | "filing" | "filed";

export default function ITRSummaryPage() {
  const params = useParams();
  const employeeId = params["employee-id"] as string;
  const taxYear = params["tax-year"] as string;

  const detail = useMemo(() => getEmployeeDetail(employeeId), [employeeId]);
  const [state, setState] = useState<SubmissionState>("ready");
  const [signMethod, setSignMethod] = useState<"dsc" | "evc">("evc");
  const [arn, setArn] = useState<string | null>(null);

  if (!detail) {
    return (
      <div className="space-y-6">
        <Link
          href="/compliance/itr"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to ITR
        </Link>
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-12 text-center">
          <p className="text-sm text-[var(--text-muted)]">Employee not found</p>
        </div>
      </div>
    );
  }

  const { employee: emp, computation: comp } = detail;
  const confirmWord = `FILE ${emp.recommendedForm} ${taxYear}`;

  function handleFile() {
    setState("filing");
    setTimeout(() => {
      const generatedArn = `CPC/${taxYear.replace("-", "")}/${Date.now().toString().slice(-8)}`;
      setArn(generatedArn);
      setState("filed");
    }, 1500);
  }

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/itr"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to ITR
      </Link>

      {/* Header */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-heading-lg text-[var(--text-primary)]">
              ITR Summary — {emp.name}
            </h1>
            <p className="text-body-sm text-[var(--text-muted)] mt-1">
              PAN: <span className="font-mono font-medium">{emp.pan}</span> · Tax Year {taxYear}
            </p>
            <div className="flex items-center gap-2 mt-2">
              <FilingStatusPill status={state === "filed" ? "FILED" : emp.filingStatus} />
              <RegimeIndicator regime={emp.regime} />
              <ITRFormRecommendationPill form={emp.recommendedForm} />
            </div>
          </div>
        </div>
      </div>

      {/* Post-submission status */}
      {state === "filed" && arn && (
        <div className="bg-[var(--success-muted)] border-2 border-[var(--success)] rounded-xl p-5 space-y-3">
          <div className="flex items-center gap-2">
            <FileCheck2 className="w-5 h-5 text-[var(--success)]" />
            <h2 className="text-sm font-bold text-[var(--success)]">ITR Filed Successfully</h2>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">ARN</p>
              <div className="flex items-center gap-1.5 mt-0.5">
                <span className="text-sm font-mono font-bold text-[var(--text-primary)]">{arn}</span>
                <button
                  onClick={() => navigator.clipboard.writeText(arn)}
                  className="text-[var(--text-muted)] hover:text-[var(--text-primary)]"
                >
                  <Copy className="w-3 h-3" />
                </button>
              </div>
            </div>
            <div>
              <p className="text-[10px] uppercase font-semibold tracking-wide text-[var(--text-muted)]">ITR-V</p>
              <button className="flex items-center gap-1 mt-0.5 text-xs font-medium text-[var(--accent)] hover:underline">
                <ExternalLink className="w-3 h-3" />
                Download ITR-V
              </button>
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
      )}

      {/* Income Summary */}
      <div>
        <h2 className="text-sm font-semibold text-[var(--text-primary)] mb-3">Income Summary</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {detail.incomeHeads.filter((h) => h.visible).map((head) => (
            <IncomeHeadCard key={head.head} detail={head} />
          ))}
        </div>
      </div>

      {/* Deductions — Old Regime only */}
      {emp.regime === "OLD" && (
        <div>
          <h2 className="text-sm font-semibold text-[var(--text-primary)] mb-3">
            Deductions (Old Regime)
          </h2>
          <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 space-y-1">
            {detail.deductions.map((d) => (
              <div key={d.section} className="flex items-center justify-between py-1">
                <span className="text-[11px] text-[var(--text-muted)]">
                  §{d.section}: {d.label}
                </span>
                <span className="text-[11px] tabular-nums font-medium text-[var(--text-primary)]">
                  ₹{d.declared.toLocaleString("en-IN")}
                  {d.limit > 0 && (
                    <span className="text-[var(--text-muted)]"> / ₹{d.limit.toLocaleString("en-IN")}</span>
                  )}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Tax Computation + Refund */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <TaxComputationPanel
          computation={comp}
          deductions={emp.regime === "OLD" ? detail.deductions : undefined}
          regime={emp.regime}
        />
        <RefundIndicator amount={comp.refundOrPayable} />
      </div>

      {/* File button */}
      {state !== "filed" && (
        <div className="flex justify-end">
          <button
            onClick={() => setState("confirming")}
            disabled={state === "filing"}
            className={cn(
              "flex items-center gap-2 px-6 py-3 rounded-lg text-sm font-semibold transition-colors",
              state === "filing"
                ? "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
                : "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
            )}
          >
            <FileCheck2 className="w-4 h-4" />
            {state === "filing" ? "Filing…" : `File ${emp.recommendedForm}`}
          </button>
        </div>
      )}

      <FilingConfirmationModal
        open={state === "confirming"}
        onClose={() => setState("ready")}
        onConfirm={handleFile}
        title={`File ${emp.recommendedForm} for Tax Year ${taxYear}`}
        details={[
          { label: "Employee", value: emp.name },
          { label: "PAN", value: emp.pan },
          { label: "Form", value: emp.recommendedForm },
          { label: "Tax Year", value: taxYear },
          { label: "Regime", value: emp.regime === "NEW" ? "New Regime (§202)" : "Old Regime" },
          { label: "Tax Liability", value: `₹${comp.totalLiability.toLocaleString("en-IN")}` },
        ]}
        warningText="This will submit the ITR to the Income Tax Department. This action is irreversible."
        confirmWord={confirmWord}
        signMethod={signMethod}
        onSignMethodChange={setSignMethod}
      />
    </div>
  );
}
