"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, Calendar, Building2, ArrowRight, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { GovStatusPill } from "@complai/ui-components";
import { ANNUAL_ENTRIES } from "./mock-data";
import type { AnnualReturnEntry, GSTR9Status, GSTR9CStatus } from "./types";

function formatINR(amount: number): string {
  if (amount >= 10_000_000) return `₹${(amount / 10_000_000).toFixed(2)} Cr`;
  if (amount >= 100_000) return `₹${(amount / 100_000).toFixed(1)} L`;
  return `₹${amount.toLocaleString("en-IN")}`;
}

const FY_OPTIONS = ["2025-26", "2024-25", "2023-24"];

function gstr9StatusPill(status: GSTR9Status) {
  const map: Record<GSTR9Status, { status: "success" | "warning" | "danger" | "info"; label: string }> = {
    NOT_STARTED: { status: "danger", label: "Not Started" },
    AGGREGATING: { status: "info", label: "Aggregating" },
    IN_REVIEW: { status: "warning", label: "In Review" },
    FILED: { status: "success", label: "Filed" },
    ACKNOWLEDGED: { status: "success", label: "Acknowledged" },
  };
  return map[status];
}

function gstr9cStatusPill(status: GSTR9CStatus) {
  const map: Record<GSTR9CStatus, { status: "success" | "warning" | "danger" | "info"; label: string }> = {
    NOT_APPLICABLE: { status: "info", label: "N/A" },
    NOT_STARTED: { status: "danger", label: "Not Started" },
    RECONCILING: { status: "warning", label: "Reconciling" },
    CERTIFIED: { status: "info", label: "Certified" },
    FILED: { status: "success", label: "Filed" },
  };
  return map[status];
}

export default function AnnualReturnLandingPage() {
  const [fy, setFy] = useState("2025-26");
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    return ANNUAL_ENTRIES.filter((e) => {
      if (e.fy !== fy) return false;
      if (search) {
        const q = search.toLowerCase();
        return e.gstin.toLowerCase().includes(q) || e.legalName.toLowerCase().includes(q);
      }
      return true;
    });
  }, [fy, search]);

  const filed = filtered.filter((e) => e.gstr9Status === "FILED" || e.gstr9Status === "ACKNOWLEDGED").length;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">Annual Returns</h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            GSTR-9 &amp; GSTR-9C filing for all registered GSTINs
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 text-[var(--text-muted)]" />
          <select
            value={fy}
            onChange={(e) => setFy(e.target.value)}
            className={cn(
              "px-3 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]"
            )}
          >
            {FY_OPTIONS.map((f) => (
              <option key={f} value={f}>FY {f}</option>
            ))}
          </select>
        </div>
      </div>

      {/* KPI row */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
          <p className="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wide">Total GSTINs</p>
          <p className="text-2xl font-bold tabular-nums text-[var(--text-primary)] mt-1">{filtered.length}</p>
        </div>
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
          <p className="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wide">GSTR-9 Filed</p>
          <p className="text-2xl font-bold tabular-nums text-[var(--success)] mt-1">{filed}</p>
        </div>
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
          <p className="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wide">9C Required</p>
          <p className="text-2xl font-bold tabular-nums text-[var(--warning)] mt-1">
            {filtered.filter((e) => e.gstr9cRequired).length}
          </p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by GSTIN or legal name..."
          className={cn(
            "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
            "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]"
          )}
        />
      </div>

      {/* Table */}
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-[var(--border-default)] bg-[var(--bg-tertiary)]">
              <th className="px-4 py-3 text-left font-semibold text-[var(--text-muted)]">GSTIN</th>
              <th className="px-4 py-3 text-left font-semibold text-[var(--text-muted)]">Legal Name</th>
              <th className="px-4 py-3 text-right font-semibold text-[var(--text-muted)]">Turnover</th>
              <th className="px-4 py-3 text-center font-semibold text-[var(--text-muted)]">GSTR-9</th>
              <th className="px-4 py-3 text-center font-semibold text-[var(--text-muted)]">GSTR-9C</th>
              <th className="px-4 py-3 text-right font-semibold text-[var(--text-muted)]">Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((entry) => (
              <EntryRow key={entry.id} entry={entry} />
            ))}
            {filtered.length === 0 && (
              <tr>
                <td colSpan={6} className="px-4 py-12 text-center text-[var(--text-muted)]">
                  No annual returns found for FY {fy}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function EntryRow({ entry }: { entry: AnnualReturnEntry }) {
  const g9 = gstr9StatusPill(entry.gstr9Status);
  const g9c = gstr9cStatusPill(entry.gstr9cStatus);

  return (
    <tr className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors">
      <td className="px-4 py-3 font-mono font-semibold text-[var(--text-primary)]">{entry.gstin}</td>
      <td className="px-4 py-3">
        <div className="flex items-center gap-2">
          <Building2 className="w-3.5 h-3.5 text-[var(--text-muted)]" />
          <span className="text-[var(--text-secondary)]">{entry.legalName}</span>
        </div>
      </td>
      <td className="px-4 py-3 text-right tabular-nums font-medium text-[var(--text-primary)]">
        {formatINR(entry.turnover)}
      </td>
      <td className="px-4 py-3 text-center">
        <GovStatusPill system="GSTN" status={g9.status} label={g9.label} />
      </td>
      <td className="px-4 py-3 text-center">
        {entry.gstr9cRequired ? (
          <GovStatusPill system="GSTN" status={g9c.status} label={g9c.label} />
        ) : (
          <span className="text-[10px] text-[var(--text-muted)]">Not required</span>
        )}
      </td>
      <td className="px-4 py-3 text-right">
        <div className="flex items-center gap-2 justify-end">
          <Link
            href={`/compliance/gst-returns/annual/${entry.fy}/${entry.gstin}/9`}
            className={cn(
              "inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[10px] font-semibold transition-colors",
              entry.gstr9Status === "FILED" || entry.gstr9Status === "ACKNOWLEDGED"
                ? "border border-[var(--success)] text-[var(--success)] hover:bg-[var(--success-muted)]"
                : "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
            )}
          >
            {entry.gstr9Status === "FILED" || entry.gstr9Status === "ACKNOWLEDGED" ? (
              <>
                <CheckCircle2 className="w-3 h-3" />
                View
              </>
            ) : (
              <>
                GSTR-9
                <ArrowRight className="w-3 h-3" />
              </>
            )}
          </Link>
          {entry.gstr9cRequired && entry.gstr9Status !== "NOT_STARTED" && (
            <Link
              href={`/compliance/gst-returns/annual/${entry.fy}/${entry.gstin}/9c`}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[10px] font-semibold border border-[var(--border-default)] text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors"
            >
              9C
              <ArrowRight className="w-3 h-3" />
            </Link>
          )}
        </div>
      </td>
    </tr>
  );
}
