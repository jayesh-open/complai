"use client";

import { cn } from "@/lib/utils";
import type { IncomeHeadDetail } from "../types";
import {
  Briefcase,
  Building2,
  TrendingUp,
  Wallet,
  Coins,
} from "lucide-react";

const HEAD_CONFIG: Record<
  string,
  { label: string; icon: typeof Briefcase; color: string }
> = {
  SALARY: { label: "Salary", icon: Briefcase, color: "var(--accent)" },
  HOUSE_PROPERTY: { label: "House Property", icon: Building2, color: "var(--info)" },
  CAPITAL_GAINS: { label: "Capital Gains", icon: TrendingUp, color: "var(--purple)" },
  BUSINESS_PROFESSION: { label: "Business / Profession", icon: Wallet, color: "var(--orange)" },
  OTHER_SOURCES: { label: "Other Sources", icon: Coins, color: "var(--warning)" },
};

function formatINR(amount: number): string {
  const abs = Math.abs(amount);
  const sign = amount < 0 ? "-" : "";
  if (abs >= 10_000_000) return `${sign}₹${(abs / 10_000_000).toFixed(2)} Cr`;
  if (abs >= 100_000) return `${sign}₹${(abs / 100_000).toFixed(1)} L`;
  return `${sign}₹${abs.toLocaleString("en-IN")}`;
}

interface IncomeHeadCardProps {
  detail: IncomeHeadDetail;
  className?: string;
}

export function IncomeHeadCard({ detail, className }: IncomeHeadCardProps) {
  if (!detail.visible) return null;

  const config = HEAD_CONFIG[detail.head];
  if (!config) return null;
  const Icon = config.icon;

  return (
    <div
      className={cn(
        "bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-4",
        className
      )}
    >
      <div className="flex items-center gap-2 mb-3">
        <div
          className="w-7 h-7 rounded-lg flex items-center justify-center"
          style={{ backgroundColor: `color-mix(in srgb, ${config.color} 15%, transparent)` }}
        >
          <Icon className="w-3.5 h-3.5" style={{ color: config.color }} />
        </div>
        <h3 className="text-xs font-semibold text-[var(--text-primary)]">
          {config.label}
        </h3>
      </div>

      <div className="grid grid-cols-3 gap-3 mb-3">
        <div>
          <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide">Gross</p>
          <p className="text-sm font-bold tabular-nums text-[var(--text-primary)]">{formatINR(detail.gross)}</p>
        </div>
        <div>
          <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide">Deductions</p>
          <p className="text-sm font-bold tabular-nums text-[var(--danger)]">{formatINR(-detail.deductions)}</p>
        </div>
        <div>
          <p className="text-[10px] text-[var(--text-muted)] uppercase font-semibold tracking-wide">Net</p>
          <p className="text-sm font-bold tabular-nums text-[var(--text-primary)]">{formatINR(detail.net)}</p>
        </div>
      </div>

      {detail.subItems.length > 0 && (
        <div className="border-t border-[var(--border-default)] pt-2 space-y-1">
          {detail.subItems.map((item) => (
            <div key={item.label} className="flex items-center justify-between text-[11px]">
              <span className="text-[var(--text-muted)] flex items-center gap-1">
                {item.label}
                {item.section && (
                  <span className="text-[9px] font-mono px-1 py-0.5 rounded bg-[var(--bg-tertiary)] text-[var(--text-muted)]">
                    {item.section}
                  </span>
                )}
              </span>
              <span
                className={cn(
                  "tabular-nums font-medium",
                  item.amount < 0 ? "text-[var(--danger)]" : "text-[var(--text-primary)]"
                )}
              >
                {formatINR(item.amount)}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
