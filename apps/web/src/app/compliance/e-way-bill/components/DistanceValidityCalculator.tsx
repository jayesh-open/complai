"use client";

import { MapPin, Clock } from "lucide-react";
import { cn } from "@/lib/utils";

interface DistanceValidityCalculatorProps {
  distanceKm: number;
  isODC?: boolean;
  className?: string;
}

function calcDays(km: number, isODC: boolean): number {
  if (km <= 0) return 0;
  const perDay = isODC ? 20 : 200;
  return km <= perDay ? 1 : Math.ceil(km / perDay);
}

export function DistanceValidityCalculator({
  distanceKm,
  isODC = false,
  className,
}: DistanceValidityCalculatorProps) {
  const days = calcDays(distanceKm, isODC);

  if (distanceKm <= 0) return null;

  return (
    <div
      data-testid="distance-validity-calc"
      className={cn(
        "flex items-center gap-3 px-4 py-3 rounded-lg",
        "bg-[var(--info-muted)] border border-[var(--info-border)]",
        className,
      )}
    >
      <MapPin className="w-4 h-4 text-[var(--info)] flex-shrink-0" />
      <div className="flex items-center gap-2 text-xs">
        <span className="font-mono font-semibold text-[var(--text-primary)]">
          {distanceKm} km
        </span>
        <span className="text-[var(--text-muted)]">&rarr;</span>
        <Clock className="w-3.5 h-3.5 text-[var(--info)]" />
        <span className="font-semibold text-[var(--text-primary)]">
          {days} day{days !== 1 ? "s" : ""}
        </span>
        <span className="text-[var(--text-muted)]">
          ({isODC ? "20 km/day ODC" : "200 km/day"})
        </span>
      </div>
    </div>
  );
}
