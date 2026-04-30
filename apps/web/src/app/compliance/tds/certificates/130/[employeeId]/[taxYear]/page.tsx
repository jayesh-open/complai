"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Download, Printer } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import { generateForm130Detail } from "../../../mock-data";
import { CertificateStatusPill } from "../../../../components/CertificateStatusPill";
import { Form130PreviewPanel } from "../../../../components/Form130PreviewPanel";

export default function Form130DetailPage() {
  const params = useParams<{ employeeId: string; taxYear: string }>();
  const employeeId = params.employeeId;
  const taxYear = params.taxYear ?? "2026-27";

  const data = generateForm130Detail(employeeId, taxYear);

  return (
    <div className="space-y-6">
      <Link href="/compliance/tds/certificates/130" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]">
        <ArrowLeft className="w-3.5 h-3.5" /> Back to Form 130 List
      </Link>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-heading-lg text-[var(--text-primary)]">Form 130 — {data.employeeName}</h1>
            <div className="flex items-center gap-3 mt-1">
              <span className="font-mono text-xs text-[var(--text-muted)]">PAN: {data.pan}</span>
              <span className="text-xs text-[var(--text-muted)]">{data.designation}</span>
              <span className={cn("px-2 py-0.5 rounded-md text-[10px] font-semibold", data.regime === "NEW" ? "bg-[var(--info-muted)] text-[var(--info)]" : "bg-[var(--warning-muted)] text-[var(--warning)]")}>{data.regime} Regime</span>
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
          <MiniStat label="Employer" value={data.employerName} />
          <MiniStat label="TAN" value={data.employerTan} />
          <MiniStat label="Total TDS" value={formatINR(data.totalTdsDeducted)} />
        </div>
      </div>

      <Form130PreviewPanel data={data} />
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
