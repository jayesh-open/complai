"use client";

import { CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { STEPS, type WizardStep } from "./types";

interface StepIndicatorProps {
  currentStep: WizardStep;
}

export function StepIndicator({ currentStep }: StepIndicatorProps) {
  const stepIndex = STEPS.findIndex((s) => s.id === currentStep);

  return (
    <div className="bg-app-card border border-app-border rounded-card p-4 sticky top-0 z-10">
      <div className="flex items-center gap-1">
        {STEPS.map((step, i) => {
          const isActive = step.id === currentStep;
          const isCompleted = i < stepIndex;
          return (
            <div key={step.id} className="flex items-center flex-1">
              <div className="flex items-center gap-2 flex-1">
                <div
                  className={cn(
                    "w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0",
                    isCompleted
                      ? "bg-[var(--success)] text-white"
                      : isActive
                        ? "bg-[var(--accent)] text-[var(--accent-text)]"
                        : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
                  )}
                >
                  {isCompleted ? <CheckCircle2 className="w-4 h-4" /> : step.number}
                </div>
                <span
                  className={cn(
                    "text-xs font-medium whitespace-nowrap",
                    isActive ? "text-foreground" : "text-foreground-muted",
                  )}
                >
                  {step.label}
                </span>
              </div>
              {i < STEPS.length - 1 && (
                <div
                  className={cn(
                    "h-[2px] flex-1 mx-2 rounded",
                    i < stepIndex ? "bg-[var(--success)]" : "bg-[var(--border-default)]",
                  )}
                />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
