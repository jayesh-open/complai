import { cn } from '../lib/utils';

type GovSystem = 'GSTN' | 'IRP' | 'EWB' | 'TRACES' | 'MCA' | 'OLTAS';
type GovStatus = 'success' | 'warning' | 'danger' | 'info';

interface GovStatusPillProps {
  system: GovSystem;
  status: GovStatus;
  label: string;
  className?: string;
}

const dotColor: Record<GovStatus, string> = {
  success: 'bg-[var(--success)]',
  warning: 'bg-[var(--warning)]',
  danger: 'bg-[var(--danger)]',
  info: 'bg-[var(--info)]',
};

export function GovStatusPill({ system, status, label, className }: GovStatusPillProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 h-6 px-2.5 rounded-[6px] border text-[10px] font-semibold',
        'bg-[var(--bg-tertiary)] border-[var(--border-default)]',
        className,
      )}
    >
      <span className={cn('w-1.5 h-1.5 rounded-full flex-shrink-0', dotColor[status])} />
      <span className="font-mono uppercase tracking-wide text-[var(--text-primary)]">{system}</span>
      <span className="text-[var(--text-muted)]">&middot;</span>
      <span className="text-[var(--text-secondary)] normal-case">{label}</span>
    </span>
  );
}
