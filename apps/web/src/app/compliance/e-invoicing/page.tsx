"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, Plus, Download, Upload, FileCheck2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { generateMockEInvoices } from "./mock-data";
import { EInvoiceKPIs } from "./components/EInvoiceKPIs";
import { EInvoiceTable } from "./components/EInvoiceTable";
import type { IRNStatus } from "./types";

type StatusFilter = "ALL" | IRNStatus;

const ALL_RECORDS = generateMockEInvoices();

export default function EInvoicingListPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("ALL");

  const filtered = useMemo(() => {
    return ALL_RECORDS.filter((r) => {
      if (statusFilter !== "ALL" && r.status !== statusFilter) return false;
      if (search) {
        const q = search.toLowerCase();
        return (
          r.invoiceNo.toLowerCase().includes(q) ||
          r.irn.toLowerCase().includes(q) ||
          r.buyerName.toLowerCase().includes(q) ||
          r.gstin.includes(q)
        );
      }
      return true;
    });
  }, [search, statusFilter]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            E-Invoicing
          </h1>
          <p className="text-body-sm text-[var(--text-muted)] mt-1">
            Generate and manage Invoice Reference Numbers (IRN) via IRP
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Link
            href="/compliance/e-invoicing/bulk-generate"
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors"
            )}
          >
            <Upload className="w-3.5 h-3.5" />
            Bulk Generate
          </Link>
          <Link
            href="/compliance/e-invoicing/generate"
            className={cn(
              "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)]",
              "hover:bg-[var(--accent-hover)] transition-colors"
            )}
          >
            <Plus className="w-3.5 h-3.5" />
            Generate IRN
          </Link>
        </div>
      </div>

      <EInvoiceKPIs records={ALL_RECORDS} />

      <div className="flex items-center gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search invoice no., IRN, buyer..."
            className={cn(
              "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)]",
              "focus:ring-2 focus:ring-[var(--accent-muted)]"
            )}
          />
        </div>
        <div className="flex items-center gap-1 p-1 rounded-lg bg-[var(--bg-secondary)] border border-[var(--border-default)]">
          {(["ALL", "GENERATED", "CANCELLED"] as StatusFilter[]).map((s) => (
            <button
              key={s}
              onClick={() => setStatusFilter(s)}
              className={cn(
                "px-3 py-1.5 rounded-md text-[10px] font-semibold uppercase tracking-wide transition-colors",
                statusFilter === s
                  ? "bg-[var(--accent)] text-[var(--accent-text)]"
                  : "text-[var(--text-muted)] hover:text-[var(--text-primary)]"
              )}
            >
              {s === "ALL" ? "All" : s === "GENERATED" ? "Generated" : "Cancelled"}
            </button>
          ))}
        </div>
        <button
          className={cn(
            "flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
            "transition-colors"
          )}
        >
          <Download className="w-3.5 h-3.5" />
          Export
        </button>
      </div>

      <EInvoiceTable records={filtered} />
    </div>
  );
}
