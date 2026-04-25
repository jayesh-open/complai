"use client";
import { useState } from 'react';
import { cn } from '../lib/utils';

interface ApprovalCardProps {
  title: string;
  subtitle: string;
  submittedBy: string;
  submittedAt: string;
  details: { label: string; value: string }[];
  warnings?: string[];
  onApprove: (comment: string) => void;
  onReject: (comment: string) => void;
  onSendBack?: (comment: string) => void;
  className?: string;
}

export function MakerCheckerApprovalCard({
  title, subtitle, submittedBy, submittedAt,
  details, warnings, onApprove, onReject, onSendBack, className,
}: ApprovalCardProps) {
  const [comment, setComment] = useState('');

  return (
    <div className={cn('bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl', className)} data-testid="approval-card">
      <div className="px-5 py-3 border-b border-[var(--border-default)]">
        <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--warning)]">Approval Required</div>
      </div>
      <div className="px-5 py-4 space-y-3">
        <div>
          <div className="text-sm font-bold text-[var(--text-primary)]">{title}</div>
          <div className="text-xs text-[var(--text-muted)]">{subtitle}</div>
        </div>
        <div className="text-[11px] text-[var(--text-muted)]">
          Submitted by: {submittedBy} &middot; {submittedAt}
        </div>
        <div className="space-y-1.5">
          {details.map((d) => (
            <div key={d.label} className="flex gap-2 text-xs">
              <span className="text-[var(--text-muted)] min-w-[80px]">{d.label}:</span>
              <span className="text-[var(--text-primary)] font-medium">{d.value}</span>
            </div>
          ))}
        </div>
        {warnings?.map((w, i) => (
          <div key={i} className="flex items-center gap-1.5 text-[11px] text-[var(--warning)]">
            <span>{'⚠'}</span> {w}
          </div>
        ))}
        <div>
          <label htmlFor="approval-comment" className="text-[11px] font-medium text-[var(--text-muted)]">Comments (optional):</label>
          <textarea
            id="approval-comment"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            rows={2}
            className={cn(
              'w-full mt-1 px-3 py-2 rounded-lg border text-xs',
              'bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)]',
              'focus:outline-none focus:border-[var(--accent)]',
            )}
          />
        </div>
      </div>
      <div className="px-5 py-3 border-t border-[var(--border-default)] flex items-center gap-2 justify-between">
        <div className="flex gap-2">
          <button
            onClick={() => onReject(comment)}
            className="px-3 py-1.5 text-[11px] font-semibold rounded-lg bg-[var(--danger-muted)] text-[var(--danger)]"
          >
            Reject
          </button>
          {onSendBack && (
            <button
              onClick={() => onSendBack(comment)}
              className="px-3 py-1.5 text-[11px] font-medium rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)]"
            >
              Send Back for Edit
            </button>
          )}
        </div>
        <button
          onClick={() => onApprove(comment)}
          className="px-4 py-1.5 text-[11px] font-bold rounded-lg bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]"
        >
          Approve &amp; Continue
        </button>
      </div>
    </div>
  );
}
