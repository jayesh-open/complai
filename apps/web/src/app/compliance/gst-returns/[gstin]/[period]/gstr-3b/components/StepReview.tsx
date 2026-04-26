"use client";

import { useState } from "react";
import { cn, formatINR } from "@/lib/utils";
import { SourceBadge } from "./SourceBadge";
import type { GSTR3BData, TaxRow, ITCRow, DataSource } from "./types";

interface StepReviewProps {
  data: GSTR3BData;
  onOverride: (section: string, index: number, field: string, value: number, reason: string) => void;
}

function TaxTable({ title, tableNum, rows, showTaxable }: {
  title: string;
  tableNum: string;
  rows: TaxRow[];
  showTaxable?: boolean;
}) {
  return (
    <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
      <div className="bg-[var(--bg-tertiary)] px-4 py-2.5 border-b border-[var(--border-default)]">
        <span className="text-xs font-bold text-foreground">{tableNum}. {title}</span>
      </div>
      <table className="w-full text-xs">
        <thead>
          <tr className="border-b border-[var(--border-default)] bg-[var(--bg-secondary)]">
            <th className="text-left px-4 py-2 text-[var(--text-muted)] font-semibold">Description</th>
            {showTaxable !== false && <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Taxable Value</th>}
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">CGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">SGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">IGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Cess</th>
            <th className="text-center px-4 py-2 text-[var(--text-muted)] font-semibold">Source</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-[var(--border-default)] last:border-0 hover:bg-[var(--bg-tertiary)] transition-colors">
              <td className="px-4 py-2 text-foreground max-w-[280px]">{row.description}</td>
              {showTaxable !== false && <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.taxableValue)}</td>}
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.cgst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.sgst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.igst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.cess)}</td>
              <td className="px-4 py-2 text-center"><SourceBadge source={row.source} /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ITCTable({ title, tableNum, rows }: {
  title: string;
  tableNum: string;
  rows: ITCRow[];
}) {
  return (
    <div className="border border-[var(--border-default)] rounded-lg overflow-hidden">
      <div className="bg-[var(--bg-tertiary)] px-4 py-2.5 border-b border-[var(--border-default)]">
        <span className="text-xs font-bold text-foreground">{tableNum}. {title}</span>
      </div>
      <table className="w-full text-xs">
        <thead>
          <tr className="border-b border-[var(--border-default)] bg-[var(--bg-secondary)]">
            <th className="text-left px-4 py-2 text-[var(--text-muted)] font-semibold">Description</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">CGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">SGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">IGST</th>
            <th className="text-right px-4 py-2 text-[var(--text-muted)] font-semibold">Cess</th>
            <th className="text-center px-4 py-2 text-[var(--text-muted)] font-semibold">Source</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-[var(--border-default)] last:border-0 hover:bg-[var(--bg-tertiary)] transition-colors">
              <td className="px-4 py-2 text-foreground max-w-[280px]">{row.description}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.cgst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.sgst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.igst)}</td>
              <td className="px-4 py-2 text-right font-mono text-foreground">{formatINR(row.cess)}</td>
              <td className="px-4 py-2 text-center"><SourceBadge source={row.source} /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export function StepReview({ data, onOverride }: StepReviewProps) {
  const [activeTab, setActiveTab] = useState<"liability" | "itc">("liability");

  const tabs = [
    { id: "liability" as const, label: "Tax Liability (Tables 1-6)" },
    { id: "itc" as const, label: "ITC (Tables 4A-D)" },
  ];

  return (
    <div className="p-6 space-y-4" data-testid="step-review">
      <div>
        <h2 className="text-heading-lg text-foreground">Step 2: Review GSTR-3B</h2>
        <p className="text-body-sm text-foreground-muted mt-1">
          Review auto-populated values. Click any cell to override with a reason.
        </p>
      </div>

      <div className="flex gap-1 bg-[var(--bg-tertiary)] p-1 rounded-lg w-fit">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            data-testid={`tab-${tab.id}`}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              "px-4 py-1.5 rounded-md text-xs font-medium transition-colors",
              activeTab === tab.id
                ? "bg-app-card text-foreground shadow-sm"
                : "text-foreground-muted hover:text-foreground",
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === "liability" && (
        <div className="space-y-4">
          <TaxTable
            title="Outward and reverse charge inward supplies"
            tableNum="Table 3.1"
            rows={data.outwardSupplies}
          />
          <TaxTable
            title="Inter-state supplies to unregistered persons, composition dealers, and UIN holders"
            tableNum="Table 3.2"
            rows={data.inwardSupplies}
          />
        </div>
      )}

      {activeTab === "itc" && (
        <div className="space-y-4">
          <ITCTable
            title="Eligible ITC"
            tableNum="Table 4(A)"
            rows={data.itcAvailed}
          />
          <ITCTable
            title="ITC Reversed"
            tableNum="Table 4(B)"
            rows={data.itcReversed}
          />
          <ITCTable
            title="Net ITC Available (4A - 4B)"
            tableNum="Table 4(C)"
            rows={data.itcNet}
          />
        </div>
      )}

      <div className="bg-[var(--bg-tertiary)] rounded-lg p-4">
        <div className="grid grid-cols-4 gap-4 text-center">
          <div>
            <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">Head</div>
          </div>
          {(["CGST", "SGST", "IGST"] as const).map((head) => (
            <div key={head}>
              <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">{head}</div>
            </div>
          ))}
          <div className="text-left text-xs font-medium text-foreground">Liability</div>
          <div className="text-xs font-mono text-foreground">{formatINR(data.totalLiability.cgst)}</div>
          <div className="text-xs font-mono text-foreground">{formatINR(data.totalLiability.sgst)}</div>
          <div className="text-xs font-mono text-foreground">{formatINR(data.totalLiability.igst)}</div>
          <div className="text-left text-xs font-medium text-[var(--success)]">ITC</div>
          <div className="text-xs font-mono text-[var(--success)]">{formatINR(data.totalITC.cgst)}</div>
          <div className="text-xs font-mono text-[var(--success)]">{formatINR(data.totalITC.sgst)}</div>
          <div className="text-xs font-mono text-[var(--success)]">{formatINR(data.totalITC.igst)}</div>
          <div className="text-left text-xs font-bold text-[var(--danger)]">Net Payable</div>
          <div className="text-xs font-mono font-bold text-[var(--danger)]">{formatINR(data.netPayable.cgst)}</div>
          <div className="text-xs font-mono font-bold text-[var(--danger)]">{formatINR(data.netPayable.sgst)}</div>
          <div className="text-xs font-mono font-bold text-[var(--danger)]">{formatINR(data.netPayable.igst)}</div>
        </div>
      </div>
    </div>
  );
}
