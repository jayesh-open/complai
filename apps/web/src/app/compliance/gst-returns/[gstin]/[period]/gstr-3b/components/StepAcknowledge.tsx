"use client";

import { CheckCircle2, Download, ExternalLink } from "lucide-react";
import { formatINR } from "@/lib/utils";
import Link from "next/link";
import type { GSTR3BData } from "./types";

interface StepAcknowledgeProps {
  data: GSTR3BData;
  arn: string;
}

export function StepAcknowledge({ data, arn }: StepAcknowledgeProps) {
  return (
    <div className="p-6 space-y-4 text-center" data-testid="step-acknowledge">
      <div className="flex justify-center">
        <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center">
          <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
        </div>
      </div>

      <div>
        <h2 className="text-heading-lg text-foreground">GSTR-3B Filed Successfully</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Your return has been filed and tax payment processed.
        </p>
      </div>

      <div data-testid="filing-receipt" className="bg-[var(--bg-tertiary)] rounded-lg p-4 inline-block mx-auto space-y-2 text-left">
        {[
          { label: "ARN", value: arn },
          { label: "GSTIN", value: data.gstin },
          { label: "Return Type", value: "GSTR-3B" },
          { label: "Period", value: data.periodLabel },
          { label: "Filed At", value: new Date().toLocaleString("en-IN") },
          { label: "Tax Paid", value: formatINR(data.netPayable.cgst + data.netPayable.sgst + data.netPayable.igst) },
          { label: "ITC Utilised", value: formatINR(data.totalITC.cgst + data.totalITC.sgst + data.totalITC.igst) },
        ].map((d) => (
          <div key={d.label} className="flex gap-3 text-xs">
            <span className="text-foreground-muted min-w-[100px]">{d.label}:</span>
            <span className="text-foreground font-bold font-mono">{d.value}</span>
          </div>
        ))}
      </div>

      <div className="flex items-center justify-center gap-3 pt-2">
        <button
          data-testid="download-ack-button"
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-foreground-muted hover:text-foreground hover:bg-[var(--bg-tertiary)] transition-colors"
        >
          <Download className="w-3.5 h-3.5" />
          Download Acknowledgement PDF
        </button>
        <Link
          href="/compliance/gst"
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium text-[var(--accent)] hover:underline"
        >
          <ExternalLink className="w-3.5 h-3.5" />
          View Audit Trail
        </Link>
      </div>

      <div className="pt-2">
        <Link
          href="/compliance/gst"
          className="text-xs text-[var(--accent)] font-medium hover:underline"
        >
          Back to GST Returns
        </Link>
      </div>
    </div>
  );
}
