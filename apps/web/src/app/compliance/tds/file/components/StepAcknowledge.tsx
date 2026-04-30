"use client";

import { CheckCircle2, Download, ExternalLink } from "lucide-react";
import { formatINR } from "@complai/ui-components";
import Link from "next/link";
import type { TDSFilingData } from "../types";

interface StepAcknowledgeProps {
  data: TDSFilingData;
  arn: string;
}

export function StepAcknowledge({ data, arn }: StepAcknowledgeProps) {
  return (
    <div className="p-6 space-y-4 text-center">
      <div className="flex justify-center">
        <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center">
          <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
        </div>
      </div>

      <div>
        <h2 className="text-heading-lg text-[var(--text-primary)]">{data.formLabel} Filed Successfully</h2>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Your TDS return has been submitted to TRACES via Sandbox gateway.
        </p>
      </div>

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 inline-block mx-auto space-y-2 text-left">
        {[
          { label: "ARN", value: arn },
          { label: "Form", value: data.formLabel },
          { label: "TAN", value: data.tan },
          { label: "Tax Year", value: data.taxYear },
          { label: "Quarter", value: data.quarterLabel },
          { label: "Filed At", value: new Date().toLocaleString("en-IN") },
          { label: "Deductees", value: String(data.deducteeCount) },
          { label: "Total Tax", value: formatINR(data.totalTax) },
        ].map((d) => (
          <div key={d.label} className="flex gap-3 text-xs">
            <span className="text-[var(--text-muted)] min-w-[100px]">{d.label}:</span>
            <span className="text-[var(--text-primary)] font-bold font-mono">{d.value}</span>
          </div>
        ))}
      </div>

      <div className="flex items-center justify-center gap-3 pt-2">
        <button className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-muted)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors">
          <Download className="w-3.5 h-3.5" />
          Download Acknowledgement
        </button>
        <Link
          href="/compliance/tds/file"
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium text-[var(--accent)] hover:underline"
        >
          <ExternalLink className="w-3.5 h-3.5" />
          Back to Filing
        </Link>
      </div>
    </div>
  );
}
