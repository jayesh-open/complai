"use client";

import { useState, useCallback, useMemo } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { ALL_EMPLOYEES } from "../../mock-data";
import { SelectStep } from "./components/SelectStep";
import { ConfigureStep } from "./components/ConfigureStep";
import { ConfirmStep } from "./components/ConfirmStep";

type WizardStep = "select" | "configure" | "confirm";

const STEPS: { id: WizardStep; number: number; label: string }[] = [
  { id: "select", number: 1, label: "Select Employees" },
  { id: "configure", number: 2, label: "Configure Batch" },
  { id: "confirm", number: 3, label: "Confirm & Submit" },
];

export default function NewBulkBatchPage() {
  const router = useRouter();
  const [step, setStep] = useState<WizardStep>("select");
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [batchName, setBatchName] = useState("");
  const [autoFetchAIS, setAutoFetchAIS] = useState(true);
  const [sendMagicLinks, setSendMagicLinks] = useState(true);

  const eligibleEmployees = useMemo(
    () => ALL_EMPLOYEES.filter((e) => e.filingStatus === "NOT_STARTED" || e.filingStatus === "AIS_FETCHED"),
    []
  );

  const toggleSelect = useCallback((id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }, []);

  const toggleAll = useCallback(() => {
    if (selected.size === eligibleEmployees.length) {
      setSelected(new Set());
    } else {
      setSelected(new Set(eligibleEmployees.map((e) => e.id)));
    }
  }, [selected.size, eligibleEmployees]);

  const handleSubmit = useCallback(() => {
    router.push("/compliance/itr/bulk/batch-003");
  }, [router]);

  const stepIndex = STEPS.findIndex((s) => s.id === step);

  return (
    <div className="space-y-6">
      <Link
        href="/compliance/itr"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to ITR
      </Link>

      <div>
        <h1 className="text-heading-lg text-[var(--text-primary)]">Create Bulk Filing Batch</h1>
        <p className="text-body-sm text-[var(--text-muted)] mt-1">
          Select employees, configure options, and initiate ITR generation
        </p>
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4">
        <div className="flex items-center gap-1">
          {STEPS.map((s, i) => {
            const isActive = s.id === step;
            const isCompleted = i < stepIndex;
            return (
              <div key={s.id} className="flex items-center flex-1">
                <div className="flex items-center gap-2 flex-1">
                  <div className={cn(
                    "w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0",
                    isCompleted ? "bg-[var(--success)] text-white" : isActive ? "bg-[var(--accent)] text-[var(--accent-text)]" : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]"
                  )}>
                    {isCompleted ? <CheckCircle2 className="w-4 h-4" /> : s.number}
                  </div>
                  <span className={cn("text-xs font-medium whitespace-nowrap", isActive ? "text-[var(--text-primary)]" : "text-[var(--text-muted)]")}>{s.label}</span>
                </div>
                {i < STEPS.length - 1 && (
                  <div className={cn("h-[2px] flex-1 mx-2 rounded", i < stepIndex ? "bg-[var(--success)]" : "bg-[var(--border-default)]")} />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {step === "select" && (
        <SelectStep
          employees={eligibleEmployees}
          selected={selected}
          onToggle={toggleSelect}
          onToggleAll={toggleAll}
          onNext={() => setStep("configure")}
        />
      )}

      {step === "configure" && (
        <ConfigureStep
          batchName={batchName}
          setBatchName={setBatchName}
          autoFetchAIS={autoFetchAIS}
          setAutoFetchAIS={setAutoFetchAIS}
          sendMagicLinks={sendMagicLinks}
          setSendMagicLinks={setSendMagicLinks}
          selectedCount={selected.size}
          onBack={() => setStep("select")}
          onNext={() => setStep("confirm")}
        />
      )}

      {step === "confirm" && (
        <ConfirmStep
          selected={selected}
          employees={eligibleEmployees}
          batchName={batchName}
          autoFetchAIS={autoFetchAIS}
          sendMagicLinks={sendMagicLinks}
          onBack={() => setStep("configure")}
          onSubmit={handleSubmit}
        />
      )}
    </div>
  );
}
