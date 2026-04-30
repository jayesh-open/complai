"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { ArrowLeft, FileText } from "lucide-react";
import { TaxYearSelector } from "../components/TaxYearSelector";
import { FilingStatusOverview } from "./components/FilingStatusOverview";
import { generateFilingGrid } from "./mock-data";

export default function TDSFilingLandingPage() {
  const [taxYear, setTaxYear] = useState("2026-27");
  const cells = useMemo(() => generateFilingGrid(), []);

  const filed = cells.filter((c) => c.status === "FILED").length;
  const pending = cells.filter((c) => c.status !== "FILED" && c.status !== "NOT_STARTED").length;

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/tds"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to TDS
      </Link>

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">TDS Quarterly Filing</h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            File Form 138 (Salary), Form 140 (Resident Non-Salary), Form 144 (Non-Resident)
          </p>
        </div>
        <TaxYearSelector value={taxYear} onChange={setTaxYear} />
      </div>

      <div className="grid grid-cols-3 gap-4">
        <StatCard icon={FileText} label="Filed" value={filed} color="text-[var(--success)]" bg="bg-[var(--success-muted)]" />
        <StatCard icon={FileText} label="In Progress" value={pending} color="text-[var(--warning)]" bg="bg-[var(--warning-muted)]" />
        <StatCard icon={FileText} label="Total Forms" value={12} color="text-[var(--text-muted)]" bg="bg-[var(--bg-tertiary)]" />
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
        <FilingStatusOverview cells={cells} taxYear={taxYear} />
      </div>
    </div>
  );
}

function StatCard({ icon: Icon, label, value, color, bg }: {
  icon: React.ElementType; label: string; value: number; color: string; bg: string;
}) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4 flex items-center gap-3">
      <div className={`w-10 h-10 rounded-lg ${bg} flex items-center justify-center`}>
        <Icon className={`w-5 h-5 ${color}`} />
      </div>
      <div>
        <p className="text-2xl font-bold text-[var(--text-primary)]">{value}</p>
        <p className="text-[10px] text-[var(--text-muted)] uppercase">{label}</p>
      </div>
    </div>
  );
}
