"use client";
import { cn } from '../lib/utils';

interface PeriodSelectorProps {
  financialYear: string;
  month: string;
  onYearChange: (year: string) => void;
  onMonthChange: (month: string) => void;
  className?: string;
}

const FY_OPTIONS = ['2024-25', '2025-26', '2026-27'];
const MONTH_OPTIONS = ['Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan', 'Feb', 'Mar'];

export function PeriodSelector({ financialYear, month, onYearChange, onMonthChange, className }: PeriodSelectorProps) {
  return (
    <div className={cn('flex items-center gap-2', className)} data-testid="period-selector">
      <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">Period:</span>
      <select
        value={financialYear}
        onChange={(e) => onYearChange(e.target.value)}
        aria-label="Financial year"
        className={cn(
          'px-3 py-1.5 rounded-lg border text-xs font-medium',
          'bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)]',
          'focus:outline-none focus:border-[var(--accent)]',
        )}
      >
        {FY_OPTIONS.map((fy) => (
          <option key={fy} value={fy}>FY {fy}</option>
        ))}
      </select>
      <select
        value={month}
        onChange={(e) => onMonthChange(e.target.value)}
        aria-label="Month"
        className={cn(
          'px-3 py-1.5 rounded-lg border text-xs font-medium',
          'bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)]',
          'focus:outline-none focus:border-[var(--accent)]',
        )}
      >
        {MONTH_OPTIONS.map((m) => (
          <option key={m} value={m}>{m} {financialYear.split('-')[m === 'Jan' || m === 'Feb' || m === 'Mar' ? 1 : 0]}</option>
        ))}
      </select>
    </div>
  );
}
