"use client";

import { cn } from "@/lib/utils";
import { DistanceValidityCalculator } from "../../components/DistanceValidityCalculator";
import type { TransportMode, SourceInvoiceForEwb } from "../../types";

interface EwbFormProps {
  invoice: SourceInvoiceForEwb;
  vehicleNo: string;
  onVehicleChange: (v: string) => void;
  distanceKm: string;
  onDistanceChange: (v: string) => void;
  transportMode: TransportMode;
  onTransportModeChange: (m: TransportMode) => void;
  onBack: () => void;
  onSubmit: () => void;
}

const MODES: TransportMode[] = ["Road", "Rail", "Air", "Ship"];

export function EwbForm({
  invoice,
  vehicleNo,
  onVehicleChange,
  distanceKm,
  onDistanceChange,
  transportMode,
  onTransportModeChange,
  onBack,
  onSubmit,
}: EwbFormProps) {
  const dist = parseInt(distanceKm, 10) || 0;
  const canSubmit = vehicleNo.length >= 6 && dist > 0;

  return (
    <div className="space-y-5 max-w-lg">
      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 space-y-2">
        <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide mb-2">
          Source Invoice
        </h3>
        <div className="text-xs">
          <span className="font-mono font-semibold text-[var(--text-primary)]">
            {invoice.invoiceNo}
          </span>
          <span className="text-[var(--text-muted)]"> &middot; {invoice.consigneeName}</span>
        </div>
      </div>

      <div className="space-y-4">
        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Transport Mode
          </label>
          <div className="flex gap-2">
            {MODES.map((m) => (
              <button
                key={m}
                onClick={() => onTransportModeChange(m)}
                className={cn(
                  "px-3 py-1.5 rounded-md text-xs font-medium transition-colors",
                  transportMode === m
                    ? "bg-[var(--accent-muted)] text-[var(--accent)] border border-[var(--accent)]"
                    : "bg-[var(--bg-tertiary)] text-[var(--text-secondary)] border border-[var(--border-default)] hover:bg-[var(--bg-secondary)]",
                )}
              >
                {m}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)] block mb-1.5">
            Vehicle Number
          </label>
          <input
            type="text"
            value={vehicleNo}
            onChange={(e) => onVehicleChange(e.target.value.toUpperCase())}
            placeholder="e.g. KA01AB1234"
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
            Distance (km)
          </label>
          <input
            type="number"
            value={distanceKm}
            onChange={(e) => onDistanceChange(e.target.value)}
            placeholder="e.g. 450"
            min={1}
            className={cn(
              "w-full px-3 py-2 rounded-lg text-xs font-mono",
              "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
              "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
              "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
            )}
          />
        </div>

        {dist > 0 && <DistanceValidityCalculator distanceKm={dist} />}
      </div>

      <div className="flex items-center gap-3 pt-2">
        <button
          onClick={onBack}
          className={cn(
            "px-4 py-2 rounded-lg text-xs font-medium",
            "border border-[var(--border-default)]",
            "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
          )}
        >
          Back
        </button>
        <button
          onClick={onSubmit}
          disabled={!canSubmit}
          className={cn(
            "px-4 py-2 rounded-lg text-xs font-semibold transition-colors",
            canSubmit
              ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
              : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed",
          )}
        >
          Review &amp; Confirm
        </button>
      </div>
    </div>
  );
}
