"use client";

import { cn, formatINR } from "@complai/ui-components";
import { DataTable } from "@complai/ui-components";
import { AlertTriangle, ShieldX } from "lucide-react";
import { SectionPill } from "../../components/SectionPill";
import { DTAAEvidenceBadge } from "../../components/DTAAEvidenceBadge";
import type { TDSFilingData, TDSFormType, DeducteeSummaryRow, PaymentCodeDistribution } from "../types";

interface StepValidateProps {
  formType: TDSFormType;
  data: TDSFilingData;
}

type Row = Record<string, unknown>;
const asD = (row: Row) => row as unknown as DeducteeSummaryRow;
const asP = (row: Row) => row as unknown as PaymentCodeDistribution;

const deducteeColumns = [
  { key: "pan", header: "PAN", render: (row: Row) => <span className="font-mono text-[11px]">{asD(row).pan}</span> },
  { key: "name", header: "Name", render: (row: Row) => <span className="text-xs">{asD(row).name}</span> },
  { key: "entryCount", header: "Entries", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono">{asD(row).entryCount}</span> },
  { key: "grossTotal", header: "Gross", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono">{formatINR(asD(row).grossTotal)}</span> },
  { key: "totalTax", header: "Total Tax", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono font-semibold">{formatINR(asD(row).totalTax)}</span> },
];

const nrColumns = [
  ...deducteeColumns,
  {
    key: "dtaa", header: "DTAA Evidence", render: (row: Row) => (
      <DTAAEvidenceBadge form41Filed={asD(row).form41Filed ?? false} trcAttached={asD(row).trcAttached ?? false} />
    ),
  },
];

const codeColumns = [
  { key: "code", header: "Code", render: (row: Row) => <span className="font-mono text-[11px]">{asP(row).code}</span> },
  { key: "label", header: "Nature of Payment", render: (row: Row) => <span className="text-xs">{asP(row).label}</span> },
  { key: "count", header: "Count", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono">{asP(row).count}</span> },
  { key: "grossTotal", header: "Gross", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono">{formatINR(asP(row).grossTotal)}</span> },
  { key: "tdsTotal", header: "TDS", align: "right" as const, render: (row: Row) => <span className="text-xs font-mono font-semibold">{formatINR(asP(row).tdsTotal)}</span> },
];

function Form138Details({ data }: { data: TDSFilingData }) {
  return (
    <>
      <SummaryStats data={data} />
      <div className="text-xs font-semibold text-[var(--text-primary)] mt-4 mb-2">Employee Aggregation (Salary)</div>
      <DataTable columns={deducteeColumns} data={data.deducteeSummaries as unknown as Row[]} emptyMessage="No salary deductees" />
    </>
  );
}

function Form140Details({ data }: { data: TDSFilingData }) {
  return (
    <>
      <SummaryStats data={data} />
      <div className="text-xs font-semibold text-[var(--text-primary)] mt-4 mb-2">Payment Code Distribution</div>
      <DataTable columns={codeColumns} data={data.paymentCodeDistribution as unknown as Row[]} emptyMessage="No entries" />
      <div className="text-xs font-semibold text-[var(--text-primary)] mt-4 mb-2">Deductee Summary</div>
      <DataTable columns={deducteeColumns} data={data.deducteeSummaries as unknown as Row[]} emptyMessage="No deductees" />
    </>
  );
}

function Form144Details({ data }: { data: TDSFilingData }) {
  return (
    <>
      <SummaryStats data={data} showCess />
      {data.dtaaBlockers.length > 0 && (
        <div className="p-3 bg-[var(--danger-muted)] border border-[var(--danger-border)] rounded-lg space-y-1">
          <div className="flex items-center gap-1.5">
            <ShieldX className="w-4 h-4 text-[var(--danger)]" />
            <span className="text-xs font-bold text-[var(--danger)]">DTAA Evidence Missing — Submit Blocked</span>
          </div>
          {data.dtaaBlockers.map((b, i) => (
            <p key={i} className="text-[10px] text-[var(--danger)] ml-5.5">{b}</p>
          ))}
        </div>
      )}
      <div className="text-xs font-semibold text-[var(--text-primary)] mt-4 mb-2">Non-Resident Deductees</div>
      <DataTable columns={nrColumns} data={data.deducteeSummaries as unknown as Row[]} emptyMessage="No NR deductees" />
    </>
  );
}

function SummaryStats({ data, showCess }: { data: TDSFilingData; showCess?: boolean }) {
  const items = [
    { label: "Deductees", value: String(data.deducteeCount) },
    { label: "Entries", value: String(data.entries.length) },
    { label: "Gross Total", value: formatINR(data.totalGross) },
    { label: "TDS", value: formatINR(data.totalTds) },
  ];
  if (showCess) {
    items.push({ label: "Surcharge", value: formatINR(data.totalSurcharge) });
    items.push({ label: "Cess (4%)", value: formatINR(data.totalCess) });
  }
  items.push({ label: "Total Tax", value: formatINR(data.totalTax) });

  return (
    <div className="grid grid-cols-4 gap-3">
      {items.map((s) => (
        <div key={s.label} className="bg-[var(--bg-tertiary)] rounded-lg p-3">
          <p className="text-[10px] text-[var(--text-muted)] uppercase">{s.label}</p>
          <p className="text-sm font-bold font-mono text-[var(--text-primary)]">{s.value}</p>
        </div>
      ))}
    </div>
  );
}

export function StepValidate({ formType, data }: StepValidateProps) {
  return (
    <div className="p-6 space-y-4">
      <div>
        <h2 className="text-heading-lg text-[var(--text-primary)]">Step 2: Validate</h2>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Review aggregated data before generating the FVU file.
        </p>
      </div>
      {formType === "138" && <Form138Details data={data} />}
      {formType === "140" && <Form140Details data={data} />}
      {formType === "144" && <Form144Details data={data} />}
    </div>
  );
}
