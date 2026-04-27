"use client";

import { useState, useCallback } from "react";
import Link from "next/link";
import { ArrowLeft, Loader2, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { generateMockSourceInvoices, generateMockEInvoices } from "../mock-data";
import { QRCodeDisplay } from "../components/QRCodeDisplay";
import { SelectInvoiceStep } from "./components/SelectInvoiceStep";
import { ValidateStep } from "./components/ValidateStep";
import type { SourceInvoice } from "../types";

type Step = "select" | "validate" | "generating" | "success";

const SOURCE_INVOICES = generateMockSourceInvoices();
const MOCK_EINVOICES = generateMockEInvoices();

function InfoRow({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex gap-4 text-xs">
      <span className="text-[var(--text-muted)] min-w-[120px]">{label}</span>
      <span className="text-[var(--text-primary)]">{children}</span>
    </div>
  );
}

export default function GenerateIRNPage() {
  const [step, setStep] = useState<Step>("select");
  const [selected, setSelected] = useState<SourceInvoice | null>(null);
  const [search, setSearch] = useState("");
  const [errors, setErrors] = useState<string[]>([]);

  const filtered = SOURCE_INVOICES.filter((inv) => {
    if (!search) return true;
    const q = search.toLowerCase();
    return (
      inv.invoiceNo.toLowerCase().includes(q) ||
      inv.buyerName.toLowerCase().includes(q)
    );
  });

  const handleSelect = useCallback((inv: SourceInvoice) => {
    setSelected(inv);
    setStep("validate");
    setErrors([]);

    const validationErrors: string[] = [];
    if (!inv.gstin || inv.gstin.length !== 15) {
      validationErrors.push("Supplier GSTIN is invalid");
    }
    if (!inv.buyerGstin || inv.buyerGstin.length !== 15) {
      validationErrors.push("Buyer GSTIN is invalid");
    }
    if (inv.items.length === 0) {
      validationErrors.push("At least one line item required");
    }
    if (inv.totalValue <= 0) {
      validationErrors.push("Total value must be positive");
    }
    const itemTotal = inv.items.reduce((s, it) => s + it.totalAmount, 0);
    if (Math.abs(itemTotal - inv.taxableValue) > 1) {
      validationErrors.push("Item total does not match taxable value");
    }
    setErrors(validationErrors);
  }, []);

  const handleGenerate = useCallback(async () => {
    setStep("generating");
    await new Promise((r) => setTimeout(r, 2000));
    setStep("success");
  }, []);

  const sampleResult = MOCK_EINVOICES[0];

  if (step === "success" && sampleResult) {
    return (
      <div className="space-y-6 max-w-2xl mx-auto" data-testid="irn-success">
        <div className="text-center space-y-4 py-8">
          <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center mx-auto">
            <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
          </div>
          <h1 className="text-heading-lg text-[var(--text-primary)]">
            IRN Generated Successfully
          </h1>
          <p className="text-body-sm text-[var(--text-muted)]">
            Invoice {selected?.invoiceNo} has been registered with IRP
          </p>
        </div>

        <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-5 space-y-3">
          <InfoRow label="IRN">
            <span className="font-mono text-[10px] break-all">
              {sampleResult.irn}
            </span>
          </InfoRow>
          <InfoRow label="Ack Number">
            <span className="font-mono">{sampleResult.ackNumber}</span>
          </InfoRow>
          <InfoRow label="Ack Date">{sampleResult.ackDate}</InfoRow>
          <InfoRow label="Status">
            <span className="text-[var(--success)] font-semibold">
              Generated
            </span>
          </InfoRow>
        </div>

        <div className="flex justify-center">
          <QRCodeDisplay
            value={sampleResult.qrCodeData}
            size={180}
            label={`ACK: ${sampleResult.ackNumber}`}
          />
        </div>

        <div className="flex justify-center gap-3">
          <Link
            data-testid="view-details-link"
            href={`/compliance/e-invoicing/${sampleResult.gstin}/${sampleResult.irn}`}
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-semibold",
              "bg-[var(--accent)] text-[var(--accent-text)]",
              "hover:bg-[var(--accent-hover)] transition-colors"
            )}
          >
            View Details
          </Link>
          <Link
            href="/compliance/e-invoicing"
            className={cn(
              "px-4 py-2 rounded-lg text-xs font-medium",
              "border border-[var(--border-default)]",
              "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
              "transition-colors"
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
          href="/compliance/e-invoicing"
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to e-invoices
        </Link>
        <h1 className="text-heading-lg text-[var(--text-primary)]">
          Generate IRN
        </h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Select a source invoice to generate an Invoice Reference Number
        </p>
      </div>

      {step === "select" && (
        <SelectInvoiceStep
          invoices={filtered}
          search={search}
          onSearch={setSearch}
          onSelect={handleSelect}
        />
      )}

      {step === "validate" && selected && (
        <ValidateStep
          invoice={selected}
          errors={errors}
          onBack={() => setStep("select")}
          onGenerate={handleGenerate}
        />
      )}

      {step === "generating" && (
        <div className="flex flex-col items-center justify-center py-20 space-y-4">
          <Loader2 className="w-8 h-8 text-[var(--accent)] animate-spin" />
          <p className="text-sm text-[var(--text-muted)]">
            Submitting to IRP...
          </p>
        </div>
      )}
    </div>
  );
}
