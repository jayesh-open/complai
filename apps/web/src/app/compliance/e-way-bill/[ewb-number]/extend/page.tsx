"use client";

import { useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2, CheckCircle2, Clock } from "lucide-react";
import { cn } from "@/lib/utils";
import { FilingConfirmationModal } from "@complai/ui-components";
import { generateMockEwbRecords, getValidityDays } from "../../mock-data";
import { DistanceValidityCalculator } from "../../components/DistanceValidityCalculator";

const ALL_RECORDS = generateMockEwbRecords();

function fmtDate(iso: string): string {
  const d = new Date(iso);
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const min = String(d.getMinutes()).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()} ${hh}:${min}`;
}

export default function ExtendValidityPage() {
  const router = useRouter();
  const { "ewb-number": ewbNumber } = useParams<{ "ewb-number": string }>();
  const record = ALL_RECORDS.find((r) => r.ewbNumber === ewbNumber) ?? ALL_RECORDS[0];

  const [additionalKm, setAdditionalKm] = useState("");
  const [reason, setReason] = useState("");
  const [newVehicle, setNewVehicle] = useState("");
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);

  const addKm = parseInt(additionalKm, 10) || 0;
  const newDays = getValidityDays(record.distanceKm + addKm);
  const newValidUntil = new Date(
    new Date(record.validFrom).getTime() + newDays * 86400000,
  );

  const handleConfirm = useCallback(async () => {
    setConfirmOpen(false);
    setSubmitting(true);
    await new Promise((r) => setTimeout(r, 1500));
    setSubmitting(false);
    setDone(true);
  }, []);

  const canSubmit = addKm > 0 && reason.length > 0;

  if (done) {
    return (
      <div className="space-y-6 max-w-lg mx-auto text-center py-12">
        <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
          <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
        </div>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Validity Extended</h1>
        <p className="text-body-sm text-[var(--text-muted)]">
          New validity: {fmtDate(newValidUntil.toISOString())}
        </p>
        <button
          onClick={() => router.push(`/compliance/e-way-bill/${record.ewbNumber}`)}
          className={cn(
            "px-4 py-2 rounded-lg text-xs font-semibold",
            "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors",
          )}
        >
          Back to EWB Detail
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <Link
          href={`/compliance/e-way-bill/${record.ewbNumber}`}
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" /> Back to EWB
        </Link>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Extend Validity</h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          EWB: <span className="font-mono">{record.ewbNumber}</span>
        </p>
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 max-w-lg space-y-2">
        <div className="flex items-center gap-2 text-xs">
          <Clock className="w-4 h-4 text-[var(--text-muted)]" />
          <span className="text-[var(--text-muted)]">Current validity expires:</span>
          <span className="font-mono font-semibold text-[var(--text-primary)]">
            {fmtDate(record.validUntil)}
          </span>
        </div>
        {addKm > 0 && (
          <div className="flex items-center gap-2 text-xs">
            <Clock className="w-4 h-4 text-[var(--success)]" />
            <span className="text-[var(--text-muted)]">New validity expires:</span>
            <span className="font-mono font-semibold text-[var(--success)]">
              {fmtDate(newValidUntil.toISOString())}
            </span>
          </div>
        )}
      </div>

      <div className="space-y-4 max-w-lg">
        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Additional Distance (km)
          </label>
          <input
            type="number"
            value={additionalKm}
            onChange={(e) => setAdditionalKm(e.target.value)}
            placeholder="e.g. 200"
            min={1}
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs font-mono",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>

        {addKm > 0 && (
          <DistanceValidityCalculator distanceKm={record.distanceKm + addKm} />
        )}

        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Reason for Extension
          </label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="e.g. Consignment delayed due to weather"
            rows={2}
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs resize-none",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>

        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            New Vehicle (optional)
          </label>
          <input
            type="text"
            value={newVehicle}
            onChange={(e) => setNewVehicle(e.target.value.toUpperCase())}
            placeholder="Leave blank if same vehicle"
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs font-mono",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>

        <div className="flex items-center gap-3 pt-2">
          <Link
            href={`/compliance/e-way-bill/${record.ewbNumber}`}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
            )}
          >
            Cancel
          </Link>
          <button
            onClick={() => setConfirmOpen(true)}
            disabled={!canSubmit || submitting}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
              canSubmit && !submitting
                ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed",
            )}
          >
            {submitting ? (
              <span className="flex items-center gap-2">
                <Loader2 className="w-3.5 h-3.5 animate-spin" /> Extending...
              </span>
            ) : (
              "Extend Validity"
            )}
          </button>
        </div>
      </div>

      <FilingConfirmationModal
        open={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        onConfirm={handleConfirm}
        title="Confirm Validity Extension"
        details={[
          { label: "EWB", value: record.ewbNumber },
          { label: "Current Expiry", value: fmtDate(record.validUntil) },
          { label: "New Expiry", value: fmtDate(newValidUntil.toISOString()) },
          { label: "Additional km", value: `${addKm} km` },
          { label: "New Vehicle", value: newVehicle || "(same)" },
        ]}
        warningText="This will extend the EWB validity on the NIC portal."
        confirmWord="EXTEND"
      />
    </div>
  );
}
