"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import {
  ArrowLeft, CheckCircle2, XCircle, Loader2,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { generateMockSourceInvoices } from "../mock-data";
import { InvoiceSelector } from "./components/InvoiceSelector";

type BatchStatus = "idle" | "running" | "done";

interface BatchResult {
  invoiceNo: string;
  success: boolean;
  irn?: string;
  error?: string;
}

const SOURCE_INVOICES = generateMockSourceInvoices();

export default function BulkGeneratePage() {
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [batchStatus, setBatchStatus] = useState<BatchStatus>("idle");
  const [results, setResults] = useState<BatchResult[]>([]);
  const [progress, setProgress] = useState(0);

  const toggleSelect = useCallback((id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else if (next.size < 50) next.add(id);
      return next;
    });
  }, []);

  const toggleAll = useCallback(() => {
    if (selected.size === SOURCE_INVOICES.length) {
      setSelected(new Set());
    } else {
      setSelected(
        new Set(SOURCE_INVOICES.slice(0, 50).map((inv) => inv.id))
      );
    }
  }, [selected.size]);

  const handleBulkGenerate = useCallback(async () => {
    const invoices = SOURCE_INVOICES.filter((inv) => selected.has(inv.id));
    setBatchStatus("running");
    setResults([]);
    setProgress(0);

    const batchResults: BatchResult[] = [];
    for (let i = 0; i < invoices.length; i++) {
      await new Promise((r) => setTimeout(r, 300 + Math.random() * 400));
      const success = Math.random() > 0.1;
      batchResults.push({
        invoiceNo: invoices[i].invoiceNo,
        success,
        irn: success ? `irn-${invoices[i].id}-${Date.now()}` : undefined,
        error: success ? undefined : "IRP validation failed: duplicate doc number",
      });
      setResults([...batchResults]);
      setProgress(Math.round(((i + 1) / invoices.length) * 100));
    }

    setBatchStatus("done");
  }, [selected]);

  const successCount = results.filter((r) => r.success).length;
  const failCount = results.filter((r) => !r.success).length;

  return (
    <div className="space-y-6">
      <div>
        <Link
          href="/compliance/e-invoicing"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to e-invoices
        </Link>
        <h1 className="text-heading-lg text-[var(--text-primary)]">
          Bulk Generate IRN
        </h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Select up to 50 invoices to generate IRNs in batch
        </p>
      </div>

      {batchStatus === "idle" && (
        <InvoiceSelector
          invoices={SOURCE_INVOICES}
          selected={selected}
          onToggle={toggleSelect}
          onToggleAll={toggleAll}
          onSubmit={handleBulkGenerate}
        />
      )}

      {(batchStatus === "running" || batchStatus === "done") && (
        <div className="space-y-4">
          <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide">
                Progress
              </h3>
              <span className="text-xs text-[var(--text-muted)] tabular-nums">
                {results.length} / {selected.size}
              </span>
            </div>
            <div className="w-full h-2 bg-[var(--bg-tertiary)] rounded-full overflow-hidden">
              <div
                className="h-full bg-[var(--accent)] rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
            {batchStatus === "done" && (
              <div className="flex items-center gap-4 mt-3">
                <span className="flex items-center gap-1 text-xs text-[var(--success)]">
                  <CheckCircle2 className="w-3.5 h-3.5" />
                  {successCount} succeeded
                </span>
                {failCount > 0 && (
                  <span className="flex items-center gap-1 text-xs text-[var(--danger)]">
                    <XCircle className="w-3.5 h-3.5" />
                    {failCount} failed
                  </span>
                )}
              </div>
            )}
            {batchStatus === "running" && (
              <div className="flex items-center gap-2 mt-3">
                <Loader2 className="w-3.5 h-3.5 text-[var(--accent)] animate-spin" />
                <span className="text-xs text-[var(--text-muted)]">
                  Processing...
                </span>
              </div>
            )}
          </div>

          <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="border-b border-[var(--border-default)]">
                  <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                    Invoice
                  </th>
                  <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                    Status
                  </th>
                  <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                    IRN / Error
                  </th>
                </tr>
              </thead>
              <tbody>
                {results.map((r) => (
                  <tr
                    key={r.invoiceNo}
                    className="border-b border-[var(--border-default)] last:border-b-0"
                  >
                    <td className="px-5 py-2.5 text-xs font-mono text-[var(--text-primary)]">
                      {r.invoiceNo}
                    </td>
                    <td className="px-5 py-2.5">
                      {r.success ? (
                        <span className="flex items-center gap-1 text-xs text-[var(--success)]">
                          <CheckCircle2 className="w-3.5 h-3.5" /> Success
                        </span>
                      ) : (
                        <span className="flex items-center gap-1 text-xs text-[var(--danger)]">
                          <XCircle className="w-3.5 h-3.5" /> Failed
                        </span>
                      )}
                    </td>
                    <td className="px-5 py-2.5 text-[10px] font-mono text-[var(--text-muted)]">
                      {r.success ? r.irn : r.error}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {batchStatus === "done" && (
            <div className="flex gap-3">
              <Link
                href="/compliance/e-invoicing"
                className={cn(
                  "px-4 py-2 rounded-lg text-xs font-semibold",
                  "bg-[var(--accent)] text-[var(--accent-text)]",
                  "hover:bg-[var(--accent-hover)] transition-colors"
                )}
              >
                View All IRNs
              </Link>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
