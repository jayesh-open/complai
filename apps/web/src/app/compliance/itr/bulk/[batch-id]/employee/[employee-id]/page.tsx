"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, CheckCircle2, XCircle, Pencil, Send, Eye } from "lucide-react";
import { cn } from "@/lib/utils";
import { getEmployeeDetail } from "../../../../mock-detail-data";
import { ITRFormRecommendationPill } from "../../../../components/ITRFormRecommendationPill";
import { RegimeIndicator } from "../../../../components/RegimeIndicator";
import { FilingStatusPill } from "../../../../components/FilingStatusPill";
import { IncomeHeadCard } from "../../../../components/IncomeHeadCard";
import { TaxComputationPanel } from "../../../../components/TaxComputationPanel";
import { AISReconciliationPanel } from "../../../../components/AISReconciliationPanel";
import { RefundIndicator } from "../../../../components/RefundIndicator";
import { AuditTrailTimeline } from "@complai/ui-components";

export default function EmployeeITRDetailPage() {
  const params = useParams();
  const router = useRouter();
  const batchId = params["batch-id"] as string;
  const employeeId = params["employee-id"] as string;

  const detail = useMemo(() => getEmployeeDetail(employeeId), [employeeId]);
  const [mismatches, setMismatches] = useState(detail?.mismatches ?? []);

  const unresolvedErrors = mismatches.filter((m) => !m.resolved && m.severity === "error");
  const canSubmit = unresolvedErrors.length === 0;

  function handleResolve(id: string, reason: string) {
    setMismatches((prev) =>
      prev.map((m) => (m.id === id ? { ...m, resolved: true, resolution: reason } : m))
    );
  }

  if (!detail) {
    return (
      <div className="space-y-6">
        <Link
          href={`/compliance/itr/bulk/${batchId}`}
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to Batch
        </Link>
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-12 text-center">
          <p className="text-sm text-[var(--text-muted)]">Employee not found</p>
        </div>
      </div>
    );
  }

  const { employee: emp, computation: comp } = detail;

  return (
    <div className="space-y-6">
      <Link
        href={`/compliance/itr/bulk/${batchId}`}
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to Batch
      </Link>

      {/* Header */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-heading-lg text-[var(--text-primary)]">{emp.name}</h1>
            <p className="text-body-sm text-[var(--text-muted)] mt-1">
              PAN: <span className="font-mono font-medium">{emp.pan}</span> · Tax Year {emp.taxYear} · {emp.department}
            </p>
            <div className="flex items-center gap-2 mt-2">
              <FilingStatusPill status={emp.filingStatus} />
              <RegimeIndicator regime={emp.regime} />
              <ITRFormRecommendationPill form={emp.recommendedForm} />
              <span className="text-[10px] text-[var(--text-muted)] italic ml-1">{detail.formReason}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Income Summary — 5 head cards */}
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
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold text-[var(--text-primary)]">Deductions (Old Regime)</h2>
            {detail.form10IEAFiled !== undefined && (
              <span className={cn(
                "text-[10px] font-medium px-2 py-1 rounded-lg border",
                detail.form10IEAFiled
                  ? "bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]"
                  : "bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning)]"
              )}>
                Form 10-IEA {detail.form10IEAFiled ? "Filed" : "Not Filed"}
              </span>
            )}
          </div>
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
        <div className="space-y-4">
          <RefundIndicator amount={comp.refundOrPayable} />
          <AuditTrailTimeline entries={detail.auditTrail} />
        </div>
      </div>

      {/* AIS Reconciliation */}
      <AISReconciliationPanel mismatches={mismatches} onResolve={handleResolve} />

      {/* Actions Footer */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          {!canSubmit && (
            <span className="text-[10px] text-[var(--danger)] font-medium">
              {unresolvedErrors.length} unresolved error{unresolvedErrors.length > 1 ? "s" : ""} — resolve before approval
            </span>
          )}
        </div>
        <div className="flex items-center gap-3">
          <button className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)] text-[var(--text-secondary)]",
            "hover:bg-[var(--bg-tertiary)] transition-colors"
          )}>
            <Eye className="w-3.5 h-3.5" />
            Form Preview
          </button>
          <button className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)] text-[var(--text-secondary)]",
            "hover:bg-[var(--bg-tertiary)] transition-colors"
          )}>
            <Pencil className="w-3.5 h-3.5" />
            Edit
          </button>
          <button className={cn(
            "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--danger)] text-[var(--danger)]",
            "hover:bg-[var(--danger-muted)] transition-colors"
          )}>
            <XCircle className="w-3.5 h-3.5" />
            Reject
          </button>
          <button
            onClick={() => router.push(`/compliance/itr/summary/${employeeId}/${emp.taxYear}`)}
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)] text-[var(--info)]",
              "hover:bg-[var(--info-muted)] transition-colors"
            )}
          >
            <Send className="w-3.5 h-3.5" />
            Send Magic Link
          </button>
          <button
            disabled={!canSubmit}
            className={cn(
              "flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
              canSubmit
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
            )}
          >
            <CheckCircle2 className="w-3.5 h-3.5" />
            Approve
          </button>
        </div>
      </div>
    </div>
  );
}
