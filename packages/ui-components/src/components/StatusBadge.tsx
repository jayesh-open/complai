import { cn } from '../lib/utils';

type BadgeVariant = 'success' | 'danger' | 'warning' | 'info' | 'purple' | 'teal' | 'default';
type BadgeSize = 'default' | 'xs';

interface StatusBadgeProps {
  variant?: BadgeVariant;
  size?: BadgeSize;
  children: React.ReactNode;
  className?: string;
}

const variantStyles: Record<BadgeVariant, string> = {
  success: 'bg-[var(--success-muted)] text-[var(--success)] border-[var(--success-border)]',
  danger: 'bg-[var(--danger-muted)] text-[var(--danger)] border-[var(--danger-border)]',
  warning: 'bg-[var(--warning-muted)] text-[var(--warning)] border-[var(--warning-border)]',
  info: 'bg-[var(--info-muted)] text-[var(--info)] border-[var(--info-border)]',
  purple: 'bg-[var(--purple-muted)] text-[var(--purple)] border-[rgba(124,58,237,0.20)]',
  teal: 'bg-[var(--teal-muted)] text-[var(--teal)] border-[rgba(13,148,136,0.20)]',
  default: 'bg-[rgba(107,112,128,0.1)] text-[var(--text-muted)] border-[rgba(107,112,128,0.15)]',
};

export function StatusBadge({ variant = 'default', size = 'default', children, className }: StatusBadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center border font-semibold uppercase tracking-wide',
        size === 'xs' ? 'px-1.5 py-px text-[9px] rounded-[4px]' : 'px-2.5 py-0.5 text-[10px] rounded-[6px]',
        variantStyles[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}
