"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import { ArrowLeft, Upload, FileSpreadsheet, CheckCircle2, XCircle, AlertTriangle } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import { DataTable } from "@complai/ui-components";
import { PAYMENT_CODE_MAP } from "../payment-codes";
import { generateMockImportRows } from "../mock-data";
import { StatChip } from "../components/StatChip";
import type { TDSImportRow } from "../types";

export default function TDSImportPage() {
  const [step, setStep] = useState<"upload" | "preview" | "done">("upload");
  const [rows, setRows] = useState<TDSImportRow[]>([]);
  const [dragOver, setDragOver] = useState(false);

  const handleFile = useCallback(() => {
    setRows(generateMockImportRows());
    setStep("preview");
  }, []);

  const validCount = rows.filter((r) => r.status === "valid").length;
  const errorCount = rows.filter((r) => r.status === "error").length;
  const warningCount = rows.filter((r) => r.status === "warning").length;

  type Row = Record<string, unknown>;
  const r = (row: Row) => row as unknown as TDSImportRow;

  const columns = [
    {
      key: "rowNumber",
      header: "#",
      render: (row: Row) => (
        <span className="font-mono text-[11px] text-[var(--text-muted)]">{r(row).rowNumber}</span>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (row: Row) => {
        const s = r(row).status;
        const Icon = s === "valid" ? CheckCircle2 : s === "error" ? XCircle : AlertTriangle;
        const color = s === "valid" ? "text-[var(--success)]" : s === "error" ? "text-[var(--danger)]" : "text-[var(--warning)]";
        return <Icon className={cn("w-4 h-4", color)} />;
      },
    },
    { key: "deducteePan", header: "PAN", render: (row: Row) => <span className="font-mono text-[11px]">{r(row).deducteePan}</span> },
    { key: "deducteeName", header: "Name", render: (row: Row) => <span className="text-xs text-[var(--text-primary)]">{r(row).deducteeName}</span> },
    {
      key: "paymentCode",
      header: "Code",
      render: (row: Row) => {
        const info = PAYMENT_CODE_MAP.get(r(row).paymentCode);
        return <span className="font-mono text-[11px]" title={info?.label}>{r(row).paymentCode}</span>;
      },
    },
    {
      key: "amount",
      header: "Amount",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs">{formatINR(r(row).amount)}</span>
      ),
    },
    {
      key: "tdsAmount",
      header: "TDS",
      align: "right" as const,
      render: (row: Row) => (
        <span className="font-mono text-xs font-semibold">{formatINR(r(row).tdsAmount)}</span>
      ),
    },
    {
      key: "errors",
      header: "Issues",
      render: (row: Row) => {
        const errs = r(row).errors;
        if (errs.length === 0) return <span className="text-[10px] text-[var(--text-muted)]">—</span>;
        return (
          <span className="text-[10px] text-[var(--danger)]">{errs.join("; ")}</span>
        );
      },
    },
  ];

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/tds"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to TDS
      </Link>

      <div>
        <h1 className="text-heading-lg text-[var(--text-primary)]">
          Import TDS Entries
        </h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Upload CSV or pull from Apex gateway
        </p>
      </div>

      {step === "upload" && (
        <div
          onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
          onDragLeave={() => setDragOver(false)}
          onDrop={(e) => { e.preventDefault(); setDragOver(false); handleFile(); }}
          className={cn(
            "rounded-xl border-2 border-dashed p-12",
            "flex flex-col items-center justify-center gap-4",
            "transition-colors",
            dragOver
              ? "border-[var(--accent)] bg-[var(--accent-muted)]"
              : "border-[var(--border-default)] bg-[var(--bg-secondary)]"
          )}
        >
          <FileSpreadsheet className="w-12 h-12 text-[var(--text-muted)]" />
          <div className="text-center">
            <p className="text-sm text-[var(--text-primary)] font-medium">
              Drag & drop a CSV file here
            </p>
            <p className="text-xs text-[var(--text-muted)] mt-1">
              Columns: PAN, Name, Payment Code, Amount, TDS Amount, Date
            </p>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={handleFile}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
                "bg-[var(--accent)] text-[var(--accent-text)]",
                "hover:bg-[var(--accent-hover)] transition-colors"
              )}
            >
              <Upload className="w-3.5 h-3.5" />
              Browse File
            </button>
            <button
              onClick={handleFile}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-medium",
                "border border-[var(--border-default)]",
                "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
                "transition-colors"
              )}
            >
              Pull from Apex
            </button>
          </div>
        </div>
      )}

      {step === "preview" && (
        <>
          <div className="flex items-center gap-4">
            <StatChip
              icon={CheckCircle2}
              color="text-[var(--success)]"
              label="Valid"
              count={validCount}
            />
            <StatChip
              icon={AlertTriangle}
              color="text-[var(--warning)]"
              label="Warnings"
              count={warningCount}
            />
            <StatChip
              icon={XCircle}
              color="text-[var(--danger)]"
              label="Errors"
              count={errorCount}
            />
            <div className="flex-1" />
            <button
              onClick={() => setStep("done")}
              disabled={validCount === 0}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold",
                "transition-colors",
                validCount > 0
                  ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                  : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
              )}
            >
              Import {validCount} entries
            </button>
            <button
              onClick={() => { setStep("upload"); setRows([]); }}
              className={cn(
                "px-4 py-2 rounded-lg text-xs font-medium",
                "border border-[var(--border-default)]",
                "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]"
              )}
            >
              Cancel
            </button>
          </div>

          <DataTable
            columns={columns}
            data={rows as unknown as Row[]}
            emptyMessage="No rows"
          />
        </>
      )}

      {step === "done" && (
        <div className="rounded-xl border border-[var(--success-border)] bg-[var(--success-muted)] p-8 text-center space-y-3">
          <CheckCircle2 className="w-12 h-12 text-[var(--success)] mx-auto" />
          <p className="text-sm font-semibold text-[var(--text-primary)]">
            {validCount} entries imported successfully
          </p>
          <p className="text-xs text-[var(--text-muted)]">
            {errorCount} rows skipped due to errors
          </p>
          <button
            onClick={() => { setStep("upload"); setRows([]); }}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]"
            )}
          >
            Import more
          </button>
        </div>
      )}
    </div>
  );
}

