"use client";

import { useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { FilingConfirmationModal } from "@complai/ui-components";
import { generateMockEwbRecords } from "../../mock-data";
import { VehicleUpdateTimeline } from "../../components/VehicleUpdateTimeline";

const ALL_RECORDS = generateMockEwbRecords();

export default function UpdateVehiclePage() {
  const router = useRouter();
  const { "ewb-number": ewbNumber } = useParams<{ "ewb-number": string }>();
  const record = ALL_RECORDS.find((r) => r.ewbNumber === ewbNumber) ?? ALL_RECORDS[0];

  const [newVehicle, setNewVehicle] = useState("");
  const [fromPlace, setFromPlace] = useState("");
  const [reason, setReason] = useState("");
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);

  const handleConfirm = useCallback(async () => {
    setConfirmOpen(false);
    setSubmitting(true);
    await new Promise((r) => setTimeout(r, 1500));
    setSubmitting(false);
    setDone(true);
  }, []);

  const canSubmit = newVehicle.length >= 6 && reason.length > 0;

  if (done) {
    return (
      <div className="space-y-6 max-w-lg mx-auto text-center py-12">
        <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
          <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
        </div>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Vehicle Updated</h1>
        <p className="text-body-sm text-[var(--text-muted)]">
          Vehicle changed to <span className="font-mono font-semibold">{newVehicle}</span>
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
        <h1 className="text-heading-lg text-[var(--text-primary)]">Update Vehicle</h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          EWB: <span className="font-mono">{record.ewbNumber}</span>
        </p>
      </div>

      {record.vehicleHistory.length > 0 && (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 max-w-lg">
          <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-4">
            Vehicle History
          </h3>
          <VehicleUpdateTimeline entries={record.vehicleHistory} />
        </div>
      )}

      <div className="space-y-4 max-w-lg">
        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            New Vehicle Number
          </label>
          <input
            type="text"
            value={newVehicle}
            onChange={(e) => setNewVehicle(e.target.value.toUpperCase())}
            placeholder="e.g. MH12AB3456"
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs font-mono",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>
        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Current Location
          </label>
          <input
            type="text"
            value={fromPlace}
            onChange={(e) => setFromPlace(e.target.value)}
            placeholder="e.g. Pune"
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>
        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Reason for Change
          </label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="e.g. Vehicle breakdown"
            rows={2}
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs resize-none",
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
                <Loader2 className="w-3.5 h-3.5 animate-spin" /> Updating...
              </span>
            ) : (
              "Update Vehicle"
            )}
          </button>
        </div>
      </div>

      <FilingConfirmationModal
        open={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        onConfirm={handleConfirm}
        title="Confirm Vehicle Update"
        details={[
          { label: "EWB", value: record.ewbNumber },
          { label: "Current Vehicle", value: record.vehicleNo },
          { label: "New Vehicle", value: newVehicle },
          { label: "Location", value: fromPlace || "—" },
          { label: "Reason", value: reason },
        ]}
        warningText="This will update the vehicle details on the NIC portal."
        confirmWord="UPDATE"
      />
    </div>
  );
}
