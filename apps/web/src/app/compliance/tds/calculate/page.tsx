"use client";

import { useState } from "react";
import Link from "next/link";
import { ArrowLeft, Save } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_DEDUCTEES } from "../mock-data";
import { DeducteeSearchCombobox } from "../components/DeducteeSearchCombobox";
import { PaymentCodePicker } from "../components/PaymentCodePicker";
import { TDSCalculatorPanel } from "../components/TDSCalculatorPanel";
import type { Deductee } from "../types";

export default function TDSCalculatePage() {
  const [deductee, setDeductee] = useState<Deductee | null>(null);
  const [paymentCode, setPaymentCode] = useState("");
  const [amount, setAmount] = useState<number>(0);
  const [txDate, setTxDate] = useState("2026-04-30");
  const [noPan, setNoPan] = useState(false);
  const [confirmWord, setConfirmWord] = useState("");
  const [saved, setSaved] = useState(false);

  const isNonResident = deductee?.residency === "NON_RESIDENT";
  const canSave = deductee && paymentCode && amount > 0 && confirmWord === "CONFIRM";

  function handleSave() {
    if (!canSave) return;
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  }

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
          Calculate TDS
        </h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          One-off TDS calculation — Sections 392 / 393 (ITA 2025)
        </p>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="space-y-4">
          <FieldGroup label="Deductee (search by PAN)">
            <DeducteeSearchCombobox
              deductees={ALL_DEDUCTEES}
              value={deductee}
              onChange={setDeductee}
            />
          </FieldGroup>

          <FieldGroup label="Payment Code / Transaction Type">
            <PaymentCodePicker
              value={paymentCode}
              onChange={setPaymentCode}
              sectionFilter={deductee?.sectionPreference}
            />
          </FieldGroup>

          <FieldGroup label="Gross Amount (₹)">
            <input
              type="number"
              min={0}
              value={amount || ""}
              onChange={(e) => setAmount(Number(e.target.value))}
              placeholder="Enter amount..."
              className={cn(
                "w-full px-3 py-2 rounded-lg text-xs",
                "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
                "focus:outline-none focus:border-[var(--accent)]",
                "focus:ring-2 focus:ring-[var(--accent-muted)]",
                "[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
              )}
            />
          </FieldGroup>

          <FieldGroup label="Transaction Date">
            <input
              type="date"
              value={txDate}
              onChange={(e) => setTxDate(e.target.value)}
              className={cn(
                "w-full px-3 py-2 rounded-lg text-xs",
                "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                "text-[var(--text-primary)]",
                "focus:outline-none focus:border-[var(--accent)]",
                "focus:ring-2 focus:ring-[var(--accent-muted)]"
              )}
            />
          </FieldGroup>

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={noPan}
              onChange={(e) => setNoPan(e.target.checked)}
              className="rounded border-[var(--border-default)] text-[var(--accent)]"
            />
            <span className="text-xs text-[var(--text-secondary)]">
              No PAN available — apply Section 397(2) rate
            </span>
          </label>

          <div className="border-t border-[var(--border-default)] pt-4 space-y-3">
            <FieldGroup label='Type "CONFIRM" to save this entry'>
              <input
                type="text"
                value={confirmWord}
                onChange={(e) => setConfirmWord(e.target.value.toUpperCase())}
                placeholder="CONFIRM"
                className={cn(
                  "w-full px-3 py-2 rounded-lg text-xs font-mono",
                  "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
                  "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
                  "focus:outline-none focus:border-[var(--accent)]",
                  "focus:ring-2 focus:ring-[var(--accent-muted)]"
                )}
              />
            </FieldGroup>
            <button
              onClick={handleSave}
              disabled={!canSave}
              className={cn(
                "w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-xs font-semibold",
                "transition-colors",
                canSave
                  ? "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)]"
                  : "bg-[var(--bg-tertiary)] text-[var(--text-disabled)] cursor-not-allowed"
              )}
            >
              <Save className="w-3.5 h-3.5" />
              Save Entry
            </button>
            {saved && (
              <div className="text-xs text-[var(--success)] font-medium text-center">
                Entry saved successfully
              </div>
            )}
          </div>
        </div>

        <TDSCalculatorPanel
          paymentCode={paymentCode}
          amount={amount}
          isNonResident={isNonResident}
          noPan={noPan}
        />
      </div>
    </div>
  );
}

function FieldGroup({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <label className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide">
        {label}
      </label>
      {children}
    </div>
  );
}
