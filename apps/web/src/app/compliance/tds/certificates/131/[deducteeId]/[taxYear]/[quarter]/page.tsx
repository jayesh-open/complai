"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Download, Printer } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import { generateForm131Detail } from "../../../../mock-data";
import { CertificateStatusPill } from "../../../../../components/CertificateStatusPill";
import { Form131PreviewPanel } from "../../../../../components/Form131PreviewPanel";

export default function Form131DetailPage() {
  const params = useParams<{ deducteeId: string; taxYear: string; quarter: string }>();
  const deducteeId = params.deducteeId;
  const taxYear = params.taxYear ?? "2026-27";
  const quarter = params.quarter ?? "q1";

  const data = generateForm131Detail(deducteeId, taxYear, quarter);

  return (
    <div className="space-y-6">
      <Link href="/compliance/tds/certificates/131" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]">
        <ArrowLeft className="w-3.5 h-3.5" /> Back to Form 131 List
      </Link>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-heading-lg text-[var(--text-primary)]">Form 131 — {data.deducteeName}</h1>
            <div className="flex items-center gap-3 mt-1">
              <span className="font-mono text-xs text-[var(--text-muted)]">PAN: {data.pan}</span>
              <span className="text-[10px] text-[var(--text-muted)] uppercase font-medium">{data.category}</span>
              <span className="text-xs font-semibold text-[var(--text-primary)]">{data.quarterLabel}</span>
              <CertificateStatusPill status={data.status} />
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium", "border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors")}>
              <Printer className="w-3.5 h-3.5" /> Print
            </button>
            <button className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold", "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors")}>
              <Download className="w-3.5 h-3.5" /> Download PDF
            </button>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-4 mt-6 pt-4 border-t border-[var(--border-default)]">
          <MiniStat label="Tax Year" value={data.taxYear} />
          <MiniStat label="Quarter" value={data.quarterLabel} />
          <MiniStat label="Deductor" value={data.deductorName} />
          <MiniStat label="TAN" value={data.deductorTan} />
        </div>
        <div className="grid grid-cols-4 gap-4 mt-4">
          <MiniStat label="Total Amount" value={formatINR(data.totalAmount)} />
          <MiniStat label="TDS" value={formatINR(data.totalTds)} />
          <MiniStat label="Surcharge" value={formatINR(data.totalSurcharge)} />
          <MiniStat label="Total Tax" value={formatINR(data.totalTax)} />
        </div>
      </div>

      <Form131PreviewPanel data={data} />
    </div>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">{label}</div>
      <div className="text-sm font-semibold text-[var(--text-primary)] mt-0.5">{value}</div>
    </div>
  );
}
