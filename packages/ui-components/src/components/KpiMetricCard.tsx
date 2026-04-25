import { cn } from '../lib/utils';
import type { ReactNode } from 'react';

interface KpiMetricCardProps {
  icon: ReactNode;
  iconColor?: string;
  value: string;
  label: string;
  subtitle?: string;
  trend?: { value: string; favorable: boolean };
  className?: string;
}

export function KpiMetricCard({ icon, iconColor = 'var(--accent)', value, label, subtitle, trend, className }: KpiMetricCardProps) {
  return (
    <div className={cn(
      'bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-[16px] p-4',
      'hover:border-[var(--accent)] transition-colors duration-150',
      className,
    )}>
      <div className="flex items-start justify-between mb-3">
        <div
          className="w-[30px] h-[30px] rounded-lg flex items-center justify-center border"
          style={{
            background: `color-mix(in srgb, ${iconColor} 15%, transparent)`,
            borderColor: `color-mix(in srgb, ${iconColor} 30%, transparent)`,
          }}
        >
          {icon}
        </div>
        {trend && (
          <span className={cn(
            'text-[11px] font-semibold',
            trend.favorable ? 'text-[var(--success)]' : 'text-[var(--danger)]',
          )}>
            {trend.favorable ? '↑' : '↓'} {trend.value}
          </span>
        )}
      </div>
      <div className="text-[20px] font-bold leading-tight text-[var(--text-primary)]" style={{ fontFeatureSettings: '"tnum"' }}>
        {value}
      </div>
      <div className="text-[10px] font-semibold uppercase tracking-[0.05em] text-[var(--text-muted)] mt-1">
        {label}
      </div>
      {subtitle && (
        <>
          <div className="border-t border-[var(--border-light)] my-2" />
          <div className="text-xs text-[var(--text-muted)]">{subtitle}</div>
        </>
      )}
    </div>
  );
}
