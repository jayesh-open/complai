"use client";

import { useState, useMemo } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, User, MapPin, Globe, Mail, Phone } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import { ALL_DEDUCTEES, ALL_ENTRIES } from "../../mock-data";
import { SectionPill } from "../../components/SectionPill";
import { DTAAEvidenceBadge } from "../../components/DTAAEvidenceBadge";
import { RecentEntriesTable } from "../../components/RecentEntriesTable";

type DetailTab = "transactions" | "certificates" | "compliance";

const TABS: { key: DetailTab; label: string }[] = [
  { key: "transactions", label: "Transactions" },
  { key: "certificates", label: "Certificates (Form 130/131)" },
  { key: "compliance", label: "Compliance Score" },
];

export default function DeducteeDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [tab, setTab] = useState<DetailTab>("transactions");

  const deductee = ALL_DEDUCTEES.find((d) => d.id === id);
  const entries = useMemo(
    () => ALL_ENTRIES.filter((e) => e.deducteeId === id),
    [id]
  );

  if (!deductee) {
    return (
      <div className="space-y-4">
        <Link
          href="/compliance/tds"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to TDS
        </Link>
        <div className="text-center py-12 text-[var(--text-muted)]">
          Deductee not found
        </div>
      </div>
    );
  }

  const totalDeducted = entries.reduce((s, e) => s + e.totalTax, 0);
  const totalGross = entries.reduce((s, e) => s + e.grossAmount, 0);
  const isNonResident = deductee.residency === "NON_RESIDENT";

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/tds"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to TDS
      </Link>

      <div
        className={cn(
          "rounded-xl border border-[var(--border-default)]",
          "bg-[var(--bg-secondary)] p-6"
        )}
      >
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-4">
            <div
              className={cn(
                "w-12 h-12 rounded-lg flex items-center justify-center",
                "bg-[var(--accent-muted)]"
              )}
            >
              <User className="w-6 h-6 text-[var(--accent)]" />
            </div>
            <div>
              <h1 className="text-heading-lg text-[var(--text-primary)]">
                {deductee.name}
              </h1>
              <div className="flex items-center gap-3 mt-1">
                <span className="font-mono text-xs text-[var(--text-muted)]">
                  PAN: {deductee.pan}
                </span>
                <SectionPill section={deductee.sectionPreference} />
                <span className="text-[10px] text-[var(--text-muted)] uppercase font-medium">
                  {deductee.category}
                </span>
                {isNonResident && (
                  <span className="text-[10px] text-[var(--purple)] font-semibold">
                    NON-RESIDENT
                  </span>
                )}
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className="text-lg font-bold text-[var(--text-primary)] tabular-nums">
              {formatINR(totalDeducted)}
            </div>
            <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">
              Total TDS deducted
            </div>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-4 mt-6 pt-4 border-t border-[var(--border-default)]">
          <InfoItem
            icon={MapPin}
            label="Address"
            value={deductee.address ?? "—"}
          />
          {isNonResident && (
            <InfoItem
              icon={Globe}
              label="Country"
              value={deductee.countryCode ?? "—"}
            />
          )}
          <InfoItem
            icon={Mail}
            label="Email"
            value={deductee.email ?? "—"}
          />
          <InfoItem
            icon={Phone}
            label="Phone"
            value={deductee.phone ?? "—"}
          />
        </div>

        {isNonResident && (
          <div className="mt-4 pt-4 border-t border-[var(--border-default)]">
            <DTAAEvidenceBadge
              form41Filed={deductee.form41Filed}
              trcAttached={deductee.trcAttached}
            />
          </div>
        )}

        <div className="grid grid-cols-3 gap-4 mt-4 pt-4 border-t border-[var(--border-default)]">
          <MiniStat label="Gross Payments" value={formatINR(totalGross)} />
          <MiniStat label="Total TDS" value={formatINR(totalDeducted)} />
          <MiniStat label="Transactions" value={String(entries.length)} />
        </div>
      </div>

      <div className="flex items-center gap-1 p-1 rounded-lg bg-[var(--bg-secondary)] border border-[var(--border-default)] w-fit">
        {TABS.map((t) => (
          <button
            key={t.key}
            onClick={() => setTab(t.key)}
            className={cn(
              "px-3 py-1.5 rounded-md text-[10px] font-semibold uppercase tracking-wide transition-colors",
              tab === t.key
                ? "bg-[var(--accent)] text-[var(--accent-text)]"
                : "text-[var(--text-muted)] hover:text-[var(--text-primary)]"
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      {tab === "transactions" && <RecentEntriesTable entries={entries} />}
      {tab === "certificates" && (
        <div className="text-center py-12 text-[var(--text-muted)] text-sm">
          Certificate generation available in Part 9e
        </div>
      )}
      {tab === "compliance" && (
        <div className="text-center py-12 text-[var(--text-muted)] text-sm">
          Compliance scoring coming soon
        </div>
      )}
    </div>
  );
}

function InfoItem({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ElementType;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-start gap-2">
      <Icon className="w-3.5 h-3.5 text-[var(--text-muted)] mt-0.5 shrink-0" />
      <div>
        <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">
          {label}
        </div>
        <div className="text-xs text-[var(--text-primary)] mt-0.5 break-all">
          {value}
        </div>
      </div>
    </div>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="text-sm font-bold text-[var(--text-primary)] tabular-nums">
        {value}
      </div>
      <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">
        {label}
      </div>
    </div>
  );
}
