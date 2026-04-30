"use client";

import type { TDSFilingData } from "../types";
import { FVUFilePreview } from "./FVUFilePreview";
import { generateFVUText } from "../mock-data";

interface StepPreviewProps {
  data: TDSFilingData;
}

export function StepPreview({ data }: StepPreviewProps) {
  const fvuText = generateFVUText(data);

  return (
    <div className="p-6 space-y-4">
      <div>
        <h2 className="text-heading-lg text-[var(--text-primary)]">Step 3: FVU Preview</h2>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Review the generated File Validation Utility output before submission.
        </p>
      </div>

      <FVUFilePreview content={fvuText} formLabel={data.formLabel} />

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
        <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">File Summary</div>
        {[
          { label: "Form", value: data.formLabel },
          { label: "TAN", value: data.tan },
          { label: "Tax Year", value: data.taxYear },
          { label: "Quarter", value: data.quarterLabel },
          { label: "Deductees", value: String(data.deducteeCount) },
          { label: "Records", value: String(data.entries.length) },
        ].map((item) => (
          <div key={item.label} className="flex gap-3 text-xs">
            <span className="text-[var(--text-muted)] min-w-[100px]">{item.label}:</span>
            <span className="text-[var(--text-primary)] font-medium">{item.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
