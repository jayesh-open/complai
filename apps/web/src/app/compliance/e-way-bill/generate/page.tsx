"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import { ArrowLeft, Loader2, CheckCircle2, Search } from "lucide-react";
import { cn } from "@/lib/utils";
import { FilingConfirmationModal, formatINR } from "@complai/ui-components";
import { generateMockSourceInvoicesForEwb, generateMockEwbRecords, getValidityDays } from "../mock-data";
import { EwbForm } from "./components/EwbForm";
import type { SourceInvoiceForEwb, TransportMode } from "../types";

type Step = "select" | "form" | "confirm" | "generating" | "success";

const SOURCE_INVOICES = generateMockSourceInvoicesForEwb();
const SAMPLE_EWB = generateMockEwbRecords()[0];

export default function GenerateEwbPage() {
  const [step, setStep] = useState<Step>("select");
  const [selected, setSelected] = useState<SourceInvoiceForEwb | null>(null);
  const [search, setSearch] = useState("");
  const [vehicleNo, setVehicleNo] = useState("");
  const [distanceKm, setDistanceKm] = useState("");
  const [transportMode, setTransportMode] = useState<TransportMode>("Road");

  const filtered = SOURCE_INVOICES.filter((inv) => {
    if (!search) return true;
    const q = search.toLowerCase();
    return inv.invoiceNo.toLowerCase().includes(q) || inv.consigneeName.toLowerCase().includes(q);
  });

  const handleSelect = useCallback((inv: SourceInvoiceForEwb) => {
    setSelected(inv);
    setStep("form");
  }, []);

  const handleGenerate = useCallback(async () => {
    setStep("generating");
    await new Promise((r) => setTimeout(r, 2000));
    setStep("success");
  }, []);

  const dist = parseInt(distanceKm, 10) || 0;
  const days = getValidityDays(dist);

  if (step === "success") {
    return (
      <div className="space-y-6 max-w-2xl mx-auto">
        <div className="text-center space-y-4 py-8">
          <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
            <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
          </div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            E-Way Bill Generated
          </h1>
          <p className="text-body-sm text-[var(--text-muted)]">
            EWB for {selected?.invoiceNo} has been registered
          </p>
          <p className="font-mono text-sm text-[var(--text-primary)]">
            {SAMPLE_EWB.ewbNumber}
          </p>
        </div>
        <div className="flex justify-center gap-3">
          <Link
            href={`/compliance/e-way-bill/${SAMPLE_EWB.ewbNumber}`}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors",
            )}
          >
            View Details
          </Link>
          <Link
            href="/compliance/e-way-bill"
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] transition-colors",
            )}
          >
            Back to List
          </Link>
        </div>
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
        <h1 className="text-heading-lg text-[var(--text-primary)]">Generate E-Way Bill</h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Select a source invoice and specify transport details
        </p>
      </div>

      {step === "select" && (
        <div className="space-y-4">
          <div className="relative max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search invoices..."
              className={cn(
                "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
                "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
                "focus:outline-none focus:border-[var(--accent)] focus:ring-2 focus:ring-[var(--accent-muted)]",
              )}
            />
          </div>
          <div className="grid gap-3">
            {filtered.map((inv) => (
              <button
                key={inv.id}
                onClick={() => handleSelect(inv)}
                className={cn(
                  "flex items-center justify-between p-4 rounded-xl text-left",
                  "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                  "hover:border-[var(--accent)] hover:bg-[var(--bg-tertiary)] transition-colors",
                )}
              >
                <div className="space-y-1">
                  <div className="text-xs font-semibold text-[var(--text-primary)] font-mono">{inv.invoiceNo}</div>
                  <div className="text-[11px] text-[var(--text-muted)]">
                    {inv.consigneeName} &middot; <span className="font-mono">{inv.consigneeGstin}</span>
                  </div>
                </div>
                <div className="text-sm font-bold text-[var(--text-primary)] tabular-nums font-mono">
                  {formatINR(inv.totalValue)}
                </div>
              </button>
            ))}
          </div>
        </div>
      )}

      {step === "form" && selected && (
        <EwbForm
          invoice={selected}
          vehicleNo={vehicleNo}
          onVehicleChange={setVehicleNo}
          distanceKm={distanceKm}
          onDistanceChange={setDistanceKm}
          transportMode={transportMode}
          onTransportModeChange={setTransportMode}
          onBack={() => setStep("select")}
          onSubmit={() => setStep("confirm")}
        />
      )}

      {step === "confirm" && (
        <FilingConfirmationModal
          open
          onClose={() => setStep("form")}
          onConfirm={handleGenerate}
          title="Generate E-Way Bill"
          details={[
            { label: "Invoice", value: selected?.invoiceNo ?? "" },
            { label: "Vehicle", value: vehicleNo },
            { label: "Distance", value: `${distanceKm} km` },
            { label: "Transport", value: transportMode },
            { label: "Validity", value: `${days} day${days !== 1 ? "s" : ""}` },
          ]}
          warningText="This will register the EWB with the NIC portal."
          confirmWord="GENERATE"
        />
      )}

      {step === "generating" && (
        <div className="flex flex-col items-center justify-center py-20 space-y-4">
          <Loader2 className="w-8 h-8 text-[var(--accent)] animate-spin" />
          <p className="text-sm text-[var(--text-muted)]">Submitting to EWB portal...</p>
        </div>
      )}
    </div>
  );
}
