"use client";

import { useState, useCallback, useMemo } from "react";
import Link from "next/link";
import { ArrowLeft, Layers, Loader2, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { FilingConfirmationModal, formatINR } from "@complai/ui-components";
import { generateMockEwbRecords } from "../mock-data";

const ALL_RECORDS = generateMockEwbRecords();
const MAX_CONSOLIDATION = 15;

export default function ConsolidateEwbPage() {
  const activeEwbs = useMemo(
    () => ALL_RECORDS.filter((r) => r.status === "ACTIVE" && !r.consolidatedEwbNo),
    [],
  );

  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);

  const toggleSelect = useCallback((id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else if (next.size < MAX_CONSOLIDATION) next.add(id);
      return next;
    });
  }, []);

  const toggleAll = useCallback(() => {
    if (selected.size === Math.min(activeEwbs.length, MAX_CONSOLIDATION)) {
      setSelected(new Set());
    } else {
      setSelected(new Set(activeEwbs.slice(0, MAX_CONSOLIDATION).map((r) => r.id)));
    }
  }, [selected.size, activeEwbs]);

  const handleConfirm = useCallback(async () => {
    setConfirmOpen(false);
    setSubmitting(true);
    await new Promise((r) => setTimeout(r, 2000));
    setSubmitting(false);
    setDone(true);
  }, []);

  const cewbNumber = `CEWB-${Date.now().toString().slice(-6)}`;

  if (done) {
    return (
      <div className="space-y-6 max-w-lg mx-auto text-center py-12">
        <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
          <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
        </div>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Consolidated EWB Created</h1>
        <p className="text-body-sm text-[var(--text-muted)]">
          {selected.size} EWBs consolidated into transit document
        </p>
        <p className="font-mono text-sm font-semibold text-[var(--text-primary)]">{cewbNumber}</p>
        <Link
          href="/compliance/e-way-bill"
          className={cn(
            "inline-flex px-4 py-2 rounded-lg text-xs font-semibold",
            "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors",
          )}
        >
          Back to List
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <Link
          href="/compliance/e-way-bill"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" /> Back to e-way bills
        </Link>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Consolidate E-Way Bills</h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Select active EWBs to consolidate into a single transit document (max {MAX_CONSOLIDATION})
        </p>
      </div>

      <div className="flex items-center justify-between">
        <label className="flex items-center gap-2 text-xs text-[var(--text-secondary)]">
          <input
            type="checkbox"
            checked={selected.size === Math.min(activeEwbs.length, MAX_CONSOLIDATION)}
            onChange={toggleAll}
            className="accent-[var(--accent)]"
          />
          Select all ({Math.min(activeEwbs.length, MAX_CONSOLIDATION)})
        </label>
        <div className="flex items-center gap-3">
          <span className="text-xs text-[var(--text-muted)]">
            {selected.size} selected (max {MAX_CONSOLIDATION})
          </span>
          <button
            onClick={() => setConfirmOpen(true)}
            disabled={selected.size < 2 || submitting}
            className={cn(
              "flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
              selected.size >= 2 && !submitting
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed",
            )}
          >
            {submitting ? (
              <><Loader2 className="w-3.5 h-3.5 animate-spin" /> Consolidating...</>
            ) : (
              <><Layers className="w-3.5 h-3.5" /> Consolidate {selected.size}</>
            )}
          </button>
        </div>
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[var(--border-default)]">
              <th className="w-10 px-4 py-2" />
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                EWB Number
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Invoice
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Vehicle
              </th>
              <th className="px-4 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                Value
              </th>
            </tr>
          </thead>
          <tbody>
            {activeEwbs.map((ewb) => (
              <tr
                key={ewb.id}
                onClick={() => toggleSelect(ewb.id)}
                className={cn(
                  "border-b border-[var(--border-default)] last:border-b-0",
                  "cursor-pointer hover:bg-[var(--bg-tertiary)] transition-colors",
                  selected.has(ewb.id) && "bg-[var(--accent-muted)]",
                )}
              >
                <td className="px-4 py-2.5 text-center">
                  <input
                    type="checkbox"
                    checked={selected.has(ewb.id)}
                    onChange={() => toggleSelect(ewb.id)}
                    className="accent-[var(--accent)]"
                  />
                </td>
                <td className="px-4 py-2.5 text-xs font-mono text-[var(--text-primary)]">
                  {ewb.ewbNumber}
                </td>
                <td className="px-4 py-2.5 text-xs font-mono text-[var(--text-muted)]">
                  {ewb.invoiceNo}
                </td>
                <td className="px-4 py-2.5 text-xs font-mono text-[var(--text-primary)]">
                  {ewb.vehicleNo}
                </td>
                <td className="px-4 py-2.5 text-xs font-mono font-semibold text-[var(--text-primary)] text-right tabular-nums">
                  {formatINR(ewb.totalValue)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <FilingConfirmationModal
        open={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        onConfirm={handleConfirm}
        title="Confirm Consolidation"
        details={[
          { label: "EWBs", value: `${selected.size} selected` },
          { label: "Transit Doc", value: cewbNumber },
        ]}
        warningText="This will create a consolidated EWB on the NIC portal."
        confirmWord="CONSOLIDATE"
      />
    </div>
  );
}
