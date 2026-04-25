import { cn, formatDate } from '../lib/utils';

interface AuditEntry {
  action: string;
  actor: string;
  timestamp: Date | string;
  detail?: string;
  status?: 'success' | 'warning' | 'info' | 'danger' | 'default';
}

interface AuditTrailTimelineProps {
  entries: AuditEntry[];
  className?: string;
}

const dotColor: Record<string, string> = {
  success: 'bg-[var(--success)]',
  warning: 'bg-[var(--warning)]',
  info: 'bg-[var(--info)]',
  danger: 'bg-[var(--danger)]',
  default: 'bg-[var(--text-muted)]',
};

export function AuditTrailTimeline({ entries, className }: AuditTrailTimelineProps) {
  return (
    <div className={cn('space-y-0', className)} data-testid="audit-trail">
      {entries.map((entry, i) => (
        <div key={i} className="flex gap-3">
          <div className="flex flex-col items-center">
            <div className={cn('w-2.5 h-2.5 rounded-full mt-1.5 flex-shrink-0', dotColor[entry.status || 'default'])} />
            {i < entries.length - 1 && <div className="w-px flex-1 bg-[var(--border-default)] my-1" />}
          </div>
          <div className="pb-4 min-w-0">
            <div className="flex items-baseline justify-between gap-4">
              <span className="text-[13px] font-semibold text-[var(--text-primary)]">{entry.action}</span>
              <span className="text-[11px] text-[var(--text-muted)] whitespace-nowrap">
                {formatDate(entry.timestamp)}
              </span>
            </div>
            <div className="text-[11px] text-[var(--text-muted)] mt-0.5">by {entry.actor}</div>
            {entry.detail && <div className="text-[11px] text-[var(--text-disabled)] mt-0.5">{entry.detail}</div>}
          </div>
        </div>
      ))}
    </div>
  );
}
