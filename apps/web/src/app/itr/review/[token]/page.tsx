"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { Loader2, AlertTriangle, CheckCircle2, Clock } from "lucide-react";
import { cn } from "@/lib/utils";
import { verifyMagicLinkToken, type TokenVerificationResult } from "../mock-token-data";
import { IncomeHeadCard } from "../../../compliance/itr/components/IncomeHeadCard";
import { TaxComputationPanel } from "../../../compliance/itr/components/TaxComputationPanel";
import { AISReconciliationPanel } from "../../../compliance/itr/components/AISReconciliationPanel";
import { RefundIndicator } from "../../../compliance/itr/components/RefundIndicator";
import { EmployeeReviewDisclosureCard } from "../../../compliance/itr/components/EmployeeReviewDisclosureCard";
import { MagicLinkExpiryBanner } from "../../../compliance/itr/components/MagicLinkExpiryBanner";
import { RequestChangesModal } from "../../../compliance/itr/components/RequestChangesModal";
import { EmployeeApprovalConfirmation } from "../../../compliance/itr/components/EmployeeApprovalConfirmation";

type PageState = "loading" | "valid" | "expired" | "used" | "approved";

export default function MagicLinkReviewPage() {
  const params = useParams();
  const token = params.token as string;
  const [pageState, setPageState] = useState<PageState>("loading");
  const [result, setResult] = useState<TokenVerificationResult | null>(null);
  const [showChangesModal, setShowChangesModal] = useState(false);
  const [arn, setArn] = useState<string | null>(null);

  useEffect(() => {
    verifyMagicLinkToken(token).then((res) => {
      setResult(res);
      setPageState(res.status);
    });
  }, [token]);

  const handleApprove = useCallback(() => {
    setPageState("loading");
    setTimeout(() => {
      setArn(`CPC/202627/${Date.now().toString().slice(-8)}`);
      setPageState("approved");
    }, 1200);
  }, []);

  const handleRequestChanges = useCallback((feedback: string) => {
    setShowChangesModal(false);
    alert(`Feedback sent to HR:\n\n${feedback}`);
  }, []);

  if (pageState === "loading") {
    return (
      <div className="flex flex-col items-center justify-center py-32 gap-4">
        <Loader2 className="w-8 h-8 text-[var(--accent)] animate-spin" />
        <p className="text-sm text-[var(--text-muted)]">Verifying your link…</p>
      </div>
    );
  }

  if (pageState === "expired") {
    return <ExpiredState expiredAt={result?.expiredAt} />;
  }

  if (pageState === "used") {
    return <UsedState outcome={result?.outcome} usedAt={result?.usedAt} />;
  }

  if (pageState === "approved" && arn) {
    return (
      <div className="space-y-6">
        <EmployeeApprovalConfirmation
          arn={arn}
          filedAt={new Date().toISOString()}
          signMethod="evc"
        />
      </div>
    );
  }

  const detail = result?.employeeDetail;
  if (!detail) return null;
  const { employee: emp, computation: comp } = detail;

  return (
    <div className="space-y-6">
      <EmployeeReviewDisclosureCard
        employeeName={emp.name}
        taxYear={emp.taxYear}
        employerName={result?.employerName ?? "Your Employer"}
        recommendedForm={emp.recommendedForm}
        regime={emp.regime}
      />

      {result?.expiresAt && <MagicLinkExpiryBanner expiresAt={result.expiresAt} />}

      <div>
        <h2 className="text-sm font-semibold text-[var(--text-primary)] mb-3">Income Summary</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {detail.incomeHeads.filter((h) => h.visible).map((head) => (
            <IncomeHeadCard key={head.head} detail={head} />
          ))}
        </div>
      </div>

      {emp.regime === "OLD" && (
        <div>
          <h2 className="text-sm font-semibold text-[var(--text-primary)] mb-3">Deductions (Old Regime)</h2>
          <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 space-y-1">
            {detail.deductions.map((d) => (
              <div key={d.section} className="flex items-center justify-between py-1">
                <span className="text-[11px] text-[var(--text-muted)]">§{d.section}: {d.label}</span>
                <span className="text-[11px] tabular-nums font-medium text-[var(--text-primary)]">
                  ₹{d.declared.toLocaleString("en-IN")}
                  {d.limit > 0 && <span className="text-[var(--text-muted)]"> / ₹{d.limit.toLocaleString("en-IN")}</span>}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <TaxComputationPanel
          computation={comp}
          deductions={emp.regime === "OLD" ? detail.deductions : undefined}
          regime={emp.regime}
        />
        <RefundIndicator amount={comp.refundOrPayable} />
      </div>

      <AISReconciliationPanel mismatches={detail.mismatches} />

      {/* Actions Footer */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 flex flex-col sm:flex-row items-center justify-between gap-4">
        <p className="text-[11px] text-[var(--text-muted)]">
          Review the details above carefully. By approving, you confirm the information is accurate.
        </p>
        <div className="flex items-center gap-3 flex-shrink-0">
          <button
            onClick={() => setShowChangesModal(true)}
            className={cn(
              "flex items-center gap-1.5 px-4 py-2.5 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)] text-[var(--text-secondary)]",
              "hover:bg-[var(--bg-tertiary)] transition-colors"
            )}
          >
            Request Changes
          </button>
          <button
            onClick={handleApprove}
            className={cn(
              "flex items-center gap-1.5 px-5 py-2.5 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors"
            )}
          >
            <CheckCircle2 className="w-3.5 h-3.5" />
            Approve &amp; E-Verify
          </button>
        </div>
      </div>

      <RequestChangesModal
        open={showChangesModal}
        onClose={() => setShowChangesModal(false)}
        onSubmit={handleRequestChanges}
        employeeName={emp.name}
      />
    </div>
  );
}

function ExpiredState({ expiredAt }: { expiredAt?: string }) {
  const formatted = expiredAt
    ? new Date(expiredAt).toLocaleDateString("en-IN", { day: "2-digit", month: "short", year: "numeric" })
    : "an unknown date";

  return (
    <div className="flex flex-col items-center justify-center py-24 text-center gap-4">
      <div className="w-14 h-14 rounded-full bg-[var(--warning-muted)] flex items-center justify-center">
        <AlertTriangle className="w-7 h-7 text-[var(--warning)]" />
      </div>
      <h1 className="text-lg font-bold text-[var(--text-primary)]">Link Expired</h1>
      <p className="text-sm text-[var(--text-muted)] max-w-md">
        This link expired on {formatted}. Contact your HR team for a new link.
      </p>
    </div>
  );
}

function UsedState({ outcome, usedAt }: { outcome?: string; usedAt?: string }) {
  const formatted = usedAt
    ? new Date(usedAt).toLocaleDateString("en-IN", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" })
    : "";

  return (
    <div className="flex flex-col items-center justify-center py-24 text-center gap-4">
      <div className="w-14 h-14 rounded-full bg-[var(--success-muted)] flex items-center justify-center">
        <CheckCircle2 className="w-7 h-7 text-[var(--success)]" />
      </div>
      <h1 className="text-lg font-bold text-[var(--text-primary)]">Link Already Used</h1>
      <p className="text-sm text-[var(--text-muted)] max-w-md">
        This link has been used. Status: <span className="font-semibold text-[var(--text-primary)]">{outcome ?? "Processed"}</span>
        {formatted && <> on {formatted}</>}.
      </p>
      <div className="flex items-center gap-1.5 text-xs text-[var(--success)] font-medium">
        <Clock className="w-3.5 h-3.5" />
        No further action required
      </div>
    </div>
  );
}
