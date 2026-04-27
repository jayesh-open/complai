"use client";

import { useState, useEffect } from "react";
import { Clock, AlertTriangle } from "lucide-react";
import { cn } from "@/lib/utils";

interface CancellationCountdownProps {
  generatedAt: string;
  className?: string;
}

function getTimeRemaining(generatedAt: string): {
  expired: boolean;
  hours: number;
  minutes: number;
  seconds: number;
  totalMs: number;
} {
  const deadline = new Date(generatedAt).getTime() + 24 * 3600 * 1000;
  const totalMs = deadline - Date.now();
  if (totalMs <= 0) {
    return { expired: true, hours: 0, minutes: 0, seconds: 0, totalMs: 0 };
  }
  const hours = Math.floor(totalMs / 3600000);
  const minutes = Math.floor((totalMs % 3600000) / 60000);
  const seconds = Math.floor((totalMs % 60000) / 1000);
  return { expired: false, hours, minutes, seconds, totalMs };
}

export function CancellationCountdown({
  generatedAt,
  className,
}: CancellationCountdownProps) {
  const [remaining, setRemaining] = useState(() =>
    getTimeRemaining(generatedAt)
  );

  useEffect(() => {
    const timer = setInterval(() => {
      setRemaining(getTimeRemaining(generatedAt));
    }, 1000);
    return () => clearInterval(timer);
  }, [generatedAt]);

  if (remaining.expired) {
    return (
      <div
        className={cn(
          "flex items-center gap-2 px-3 py-2 rounded-lg",
          "bg-[var(--danger-muted)] border border-[var(--danger-border)]",
          className
        )}
      >
        <AlertTriangle className="w-4 h-4 text-[var(--danger)]" />
        <span className="text-xs font-medium text-[var(--danger)]">
          Cancellation window expired
        </span>
      </div>
    );
  }

  const isUrgent = remaining.totalMs < 3600000;

  return (
    <div
      className={cn(
        "flex items-center gap-2 px-3 py-2 rounded-lg border",
        isUrgent
          ? "bg-[var(--warning-muted)] border-[var(--warning-border)]"
          : "bg-[var(--bg-tertiary)] border-[var(--border-default)]",
        className
      )}
    >
      <Clock
        className={cn(
          "w-4 h-4",
          isUrgent ? "text-[var(--warning)]" : "text-[var(--text-muted)]"
        )}
      />
      <span className="text-xs text-[var(--text-secondary)]">
        Cancel within:
      </span>
      <span
        className={cn(
          "text-xs font-mono font-semibold tabular-nums",
          isUrgent
            ? "text-[var(--warning)]"
            : "text-[var(--text-primary)]"
        )}
      >
        {String(remaining.hours).padStart(2, "0")}:
        {String(remaining.minutes).padStart(2, "0")}:
        {String(remaining.seconds).padStart(2, "0")}
      </span>
    </div>
  );
}
