"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, Plus, Download } from "lucide-react";
import { cn } from "@/lib/utils";
import { generateMockEwbRecords } from "./mock-data";
import { EwbKPIs } from "./components/EwbKPIs";
import { EwbTable } from "./components/EwbTable";
import type { EwbStatus } from "./types";

const ALL_RECORDS = generateMockEwbRecords();

const STATUS_TABS: { label: string; value: EwbStatus | "ALL" }[] = [
  { label: "All", value: "ALL" },
  { label: "Active", value: "ACTIVE" },
  { label: "Expired", value: "EXPIRED" },
  { label: "Cancelled", value: "CANCELLED" },
  { label: "Consolidated", value: "CONSOLIDATED" },
];

export default function EWayBillListPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<EwbStatus | "ALL">("ALL");

  const filtered = useMemo(() => {
    let list = ALL_RECORDS;
    if (statusFilter !== "ALL") {
      list = list.filter((r) => r.status === statusFilter);
    }
    if (search) {
      const q = search.toLowerCase();
      list = list.filter(
        (r) =>
          r.ewbNumber.toLowerCase().includes(q) ||
          r.vehicleNo.toLowerCase().includes(q),
      );
    }
    return list;
  }, [search, statusFilter]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            E-Way Bills
          </h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            Manage e-way bills for goods transport
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Link
            href="/compliance/e-way-bill/consolidate"
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors",
            )}
          >
            Consolidate
          </Link>
          <button
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors",
            )}
          >
            <Download className="w-3.5 h-3.5" /> Export
          </button>
          <Link
            href="/compliance/e-way-bill/generate"
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)]",
              "hover:bg-[var(--accent-hover)] transition-colors",
            )}
          >
            <Plus className="w-3.5 h-3.5" /> Generate EWB
          </Link>
        </div>
      </div>

      <EwbKPIs records={ALL_RECORDS} />

      <div className="flex items-center gap-4">
        <div className="relative max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search EWB number or vehicle..."
            className={cn(
              "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)]",
              "focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>
        <div className="flex gap-1">
          {STATUS_TABS.map((tab) => (
            <button
              key={tab.value}
              onClick={() => setStatusFilter(tab.value)}
              className={cn(
                "px-3 py-1.5 rounded-md text-xs font-medium transition-colors",
                statusFilter === tab.value
                  ? "bg-[var(--accent-muted)] text-[var(--accent)]"
                  : "text-[var(--text-muted)] hover:text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <EwbTable records={filtered} />
    </div>
  );
}
