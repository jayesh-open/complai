import { cn } from '../lib/utils';
import type { ReactNode } from 'react';

interface ReconciliationSplitPaneProps {
  leftTitle: string;
  rightTitle: string;
  leftContent: ReactNode;
  rightContent: ReactNode;
  className?: string;
}

export function ReconciliationSplitPane({ leftTitle, rightTitle, leftContent, rightContent, className }: ReconciliationSplitPaneProps) {
  return (
    <div className={cn('border border-[var(--border-default)] rounded-xl overflow-hidden', className)} data-testid="recon-split-pane">
      <div className="grid grid-cols-2 divide-x divide-[var(--border-default)]">
        <div className="px-4 py-2.5 bg-[var(--bg-tertiary)]">
          <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">{leftTitle}</span>
        </div>
        <div className="px-4 py-2.5 bg-[var(--bg-tertiary)]">
          <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">{rightTitle}</span>
        </div>
      </div>
      <div className="grid grid-cols-2 divide-x divide-[var(--border-default)] min-h-[200px]">
        <div className="p-4">{leftContent}</div>
        <div className="p-4">{rightContent}</div>
      </div>
    </div>
  );
}
