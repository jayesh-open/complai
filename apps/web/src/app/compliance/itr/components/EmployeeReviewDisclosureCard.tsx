"use client";

import { cn } from "@/lib/utils";
import { FileText, Building2, Info } from "lucide-react";
import type { ITRForm, TaxRegime } from "../types";
import { ITRFormRecommendationPill } from "./ITRFormRecommendationPill";
import { RegimeIndicator } from "./RegimeIndicator";

interface EmployeeReviewDisclosureCardProps {
  employeeName: string;
  taxYear: string;
  employerName: string;
  recommendedForm: ITRForm;
  regime: TaxRegime;
  className?: string;
}

export function EmployeeReviewDisclosureCard({
  employeeName,
  taxYear,
  employerName,
  recommendedForm,
  regime,
  className,
}: EmployeeReviewDisclosureCardProps) {
  return (
    <div
      className={cn(
        "bg-[var(--accent-muted)] border border-[var(--accent)] rounded-xl p-5",
        className
      )}
    >
      <div className="flex items-start gap-3">
        <div className="w-8 h-8 rounded-lg bg-[var(--accent)] flex items-center justify-center flex-shrink-0">
          <FileText className="w-4 h-4 text-[var(--accent-text)]" />
        </div>
        <div className="space-y-2 flex-1">
          <h2 className="text-sm font-bold text-[var(--text-primary)]">
            You are reviewing your ITR for Tax Year {taxYear}
          </h2>
          <div className="flex items-center gap-2 text-xs text-[var(--text-secondary)]">
            <Building2 className="w-3.5 h-3.5" />
            <span>Submitted by <span className="font-semibold">{employerName}</span></span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-[11px] text-[var(--text-muted)]">Form recommended:</span>
            <ITRFormRecommendationPill form={recommendedForm} />
            <RegimeIndicator regime={regime} />
          </div>
          <div className="flex items-start gap-1.5 mt-1 text-[11px] text-[var(--text-muted)]">
            <Info className="w-3 h-3 flex-shrink-0 mt-0.5" />
            <span>
              Based on data from your Form 130 (employer) and AIS (Form 168 from Income Tax Dept)
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
