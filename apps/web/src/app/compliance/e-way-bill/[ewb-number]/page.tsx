"use client";

import { useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { formatINR, FilingConfirmationModal } from "@complai/ui-components";
import { cn } from "@/lib/utils";
import { generateMockEwbRecords } from "../mock-data";
import { EwbStatusPill } from "../components/EwbStatusPill";
import { DistanceValidityCalculator } from "../components/DistanceValidityCalculator";
import { VehicleUpdateTimeline } from "../components/VehicleUpdateTimeline";
import { EwbDetailItems } from "./components/EwbDetailItems";
import { EwbActions } from "./components/EwbActions";

const ALL_RECORDS = generateMockEwbRecords();

function fmtDate(iso: string): string {
  const d = new Date(iso);
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const min = String(d.getMinutes()).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()} ${hh}:${min}`;
}

function InfoRow({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex gap-4 text-xs">
      <span className="text-[var(--text-muted)] min-w-[130px]">{label}</span>
      <span className="text-[var(--text-primary)]">{children}</span>
    </div>
  );
}

export default function EwbDetailPage() {
  const { "ewb-number": ewbNumber } = useParams<{ "ewb-number": string }>();
  const [cancelOpen, setCancelOpen] = useState(false);
  const record = ALL_RECORDS.find((r) => r.ewbNumber === ewbNumber) ?? ALL_RECORDS[0];

  return (
    <div className="space-y-6">
      <div>
        <Link
          href="/compliance/e-way-bill"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" /> Back to e-way bills
        </Link>
        <div className="flex items-center gap-3">
          <h1 className="text-heading-lg text-[var(--text-primary)] font-mono">
            {record.ewbNumber}
          </h1>
          <EwbStatusPill status={record.status} />
        </div>
      </div>

      <EwbActions record={record} onCancel={() => setCancelOpen(true)} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className={cn(
            "bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 space-y-3",
          )}>
            <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-3">
              EWB Summary
            </h3>
            <InfoRow label="EWB Number">
              <span className="font-mono">{record.ewbNumber}</span>
            </InfoRow>
            <InfoRow label="Source Invoice">
              <span className="font-mono">{record.invoiceNo}</span>
            </InfoRow>
            <InfoRow label="Supplier GSTIN">
              <span className="font-mono">{record.gstin}</span>
            </InfoRow>
            <InfoRow label="Consignee">
              {record.consigneeName}{" "}
              <span className="font-mono text-[10px] text-[var(--text-muted)]">
                ({record.consigneeGstin})
              </span>
            </InfoRow>
            <InfoRow label="Transport Mode">{record.transportMode}</InfoRow>
            <InfoRow label="Vehicle">
              <span className="font-mono">{record.vehicleNo}</span>
            </InfoRow>
            <InfoRow label="Distance">{record.distanceKm} km</InfoRow>
            <InfoRow label="Valid From">{fmtDate(record.validFrom)}</InfoRow>
            <InfoRow label="Valid Until">{fmtDate(record.validUntil)}</InfoRow>
            <InfoRow label="Total Value">
              <span className="font-semibold">{formatINR(record.totalValue)}</span>
            </InfoRow>
          </div>

          <DistanceValidityCalculator distanceKm={record.distanceKm} />
          <EwbDetailItems items={record.items} />
        </div>

        <div className="space-y-6">
          {record.vehicleHistory.length > 0 && (
            <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
              <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-4">
                Vehicle History
              </h3>
              <VehicleUpdateTimeline entries={record.vehicleHistory} />
            </div>
          )}

          {record.consolidatedEwbNo && (
            <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5">
              <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-2">
                Consolidation
              </h3>
              <p className="text-xs text-[var(--text-secondary)]">
                Part of consolidated EWB:{" "}
                <span className="font-mono font-semibold">{record.consolidatedEwbNo}</span>
              </p>
            </div>
          )}
        </div>
      </div>

      <FilingConfirmationModal
        open={cancelOpen}
        onClose={() => setCancelOpen(false)}
        onConfirm={() => setCancelOpen(false)}
        title="Cancel E-Way Bill"
        details={[
          { label: "EWB Number", value: record.ewbNumber },
          { label: "Invoice", value: record.invoiceNo },
          { label: "Vehicle", value: record.vehicleNo },
        ]}
        warningText="Cancellation is irreversible. The EWB cannot be reinstated."
        confirmWord="CANCEL"
      />
    </div>
  );
}
