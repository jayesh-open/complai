import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";

interface BulkJob {
  id: string; title: string; progress: number; total: number;
  status: "running" | "done" | "error"; eta?: string;
}

function BulkOperationTray({ jobs: initialJobs, onStop }: { jobs: BulkJob[]; onStop?: (id: string) => void }) {
  const [collapsed, setCollapsed] = useState(false);
  const jobs = initialJobs;
  if (jobs.length === 0) return null;
  return (
    <div className="w-80 bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl shadow-[var(--shadow-lg)]">
      <div className="px-4 py-3 flex items-center justify-between cursor-pointer border-b border-[var(--border-default)]"
        onClick={() => setCollapsed(!collapsed)}>
        <span className="text-xs font-semibold text-[var(--text-primary)]">{jobs.length} Background Jobs</span>
        <button className="text-[var(--text-muted)] text-xs">{collapsed ? "▲" : "—"}</button>
      </div>
      {!collapsed && (
        <div className="max-h-[360px] overflow-y-auto divide-y divide-[var(--border-default)]">
          {jobs.map((job) => (
            <div key={job.id} className="px-4 py-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-[var(--text-primary)]">{job.title}</span>
                {job.status === "done" && <span className="text-[10px] text-[var(--success)] font-semibold">✓ Done</span>}
                {job.status === "error" && <span className="text-[10px] text-[var(--danger)] font-semibold">✗ Error</span>}
              </div>
              {job.status === "running" && (
                <>
                  <div className="mt-2 h-1.5 bg-[var(--bg-tertiary)] rounded-full overflow-hidden">
                    <div className="h-full bg-[var(--accent)] rounded-full transition-all duration-500"
                      style={{ width: `${(job.progress / job.total) * 100}%` }} />
                  </div>
                  <div className="flex items-center justify-between mt-1">
                    <span className="text-[10px] text-[var(--text-muted)]">
                      {job.progress} / {job.total} — {Math.round((job.progress / job.total) * 100)}%
                    </span>
                    {onStop && <button onClick={() => onStop(job.id)} className="text-[10px] text-[var(--danger)] hover:underline">stop</button>}
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

const meta: Meta<typeof BulkOperationTray> = {
  title: "Compliance/BulkOperationTray",
  component: BulkOperationTray,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BulkOperationTray>;

export const Default: Story = {
  args: {
    jobs: [
      { id: "1", title: "E-Invoice Generation (Batch 42)", progress: 142, total: 228, status: "running", eta: "~3 min" },
      { id: "2", title: "GSTR-1 JSON Export", progress: 500, total: 500, status: "done" },
      { id: "3", title: "Vendor GSTIN Validation", progress: 38, total: 120, status: "running" },
    ],
    onStop: (id: string) => alert(`Stopping job ${id}`),
  },
};

export const WithError: Story = {
  args: {
    jobs: [
      { id: "1", title: "TDS Certificate Download", progress: 45, total: 200, status: "error" },
      { id: "2", title: "Bulk Invoice Upload", progress: 180, total: 300, status: "running" },
    ],
  },
};

export const AllDone: Story = {
  args: {
    jobs: [
      { id: "1", title: "E-Invoice Generation", progress: 228, total: 228, status: "done" },
      { id: "2", title: "GSTR-1 JSON Export", progress: 500, total: 500, status: "done" },
    ],
  },
};
