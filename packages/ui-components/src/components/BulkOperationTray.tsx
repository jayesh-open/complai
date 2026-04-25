"use client";
import { useState } from 'react';
import { cn } from '../lib/utils';

interface BulkJob {
  id: string;
  title: string;
  progress: number;
  total: number;
  status: 'running' | 'done' | 'error';
  startedAt?: string;
  eta?: string;
}

interface BulkOperationTrayProps {
  jobs: BulkJob[];
  onStop?: (jobId: string) => void;
  className?: string;
}

export function BulkOperationTray({ jobs, onStop, className }: BulkOperationTrayProps) {
  const [collapsed, setCollapsed] = useState(false);

  if (jobs.length === 0) return null;

  return (
    <div className={cn(
      'fixed bottom-4 right-4 z-40 w-80',
      'bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl shadow-[var(--shadow-lg)]',
      className,
    )} data-testid="bulk-tray">
      <div
        className="px-4 py-3 flex items-center justify-between cursor-pointer border-b border-[var(--border-default)]"
        onClick={() => setCollapsed(!collapsed)}
      >
        <span className="text-xs font-semibold text-[var(--text-primary)]">{jobs.length} Background Jobs</span>
        <button className="text-[var(--text-muted)] text-xs">{collapsed ? '▲' : '—'}</button>
      </div>
      {!collapsed && (
        <div className="max-h-[360px] overflow-y-auto divide-y divide-[var(--border-default)]">
          {jobs.map((job) => (
            <div key={job.id} className="px-4 py-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-[var(--text-primary)]">{job.title}</span>
                {job.status === 'done' && <span className="text-[10px] text-[var(--success)] font-semibold">{'✓'} Done</span>}
              </div>
              {job.status === 'running' && (
                <>
                  <div className="mt-2 h-1.5 bg-[var(--bg-tertiary)] rounded-full overflow-hidden">
                    <div
                      className="h-full bg-[var(--accent)] rounded-full transition-all duration-500"
                      style={{ width: `${(job.progress / job.total) * 100}%` }}
                    />
                  </div>
                  <div className="flex items-center justify-between mt-1">
                    <span className="text-[10px] text-[var(--text-muted)]">
                      {job.progress} / {job.total} &mdash; {Math.round((job.progress / job.total) * 100)}%
                    </span>
                    {onStop && (
                      <button
                        onClick={() => onStop(job.id)}
                        className="text-[10px] text-[var(--danger)] hover:underline"
                      >
                        stop
                      </button>
                    )}
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
