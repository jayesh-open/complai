import type { Meta, StoryObj } from "@storybook/react";

type AuditStatus = "success" | "warning" | "info" | "danger" | "default";

interface AuditEntry {
  action: string; actor: string; timestamp: string; detail?: string; status?: AuditStatus;
}

const dotColor: Record<AuditStatus, string> = {
  success: "bg-[var(--success)]", warning: "bg-[var(--warning)]",
  info: "bg-[var(--info)]", danger: "bg-[var(--danger)]", default: "bg-[var(--text-muted)]",
};

function AuditTrailTimeline({ entries }: { entries: AuditEntry[] }) {
  return (
    <div className="space-y-0">
      {entries.map((entry, i) => (
        <div key={i} className="flex gap-3">
          <div className="flex flex-col items-center">
            <div className={`w-2.5 h-2.5 rounded-full mt-1.5 flex-shrink-0 ${dotColor[entry.status || "default"]}`} />
            {i < entries.length - 1 && <div className="w-px flex-1 bg-[var(--border-default)] my-1" />}
          </div>
          <div className="pb-4 min-w-0">
            <div className="flex items-baseline justify-between gap-4">
              <span className="text-[13px] font-semibold text-[var(--text-primary)]">{entry.action}</span>
              <span className="text-[11px] text-[var(--text-muted)] whitespace-nowrap">{entry.timestamp}</span>
            </div>
            <div className="text-[11px] text-[var(--text-muted)] mt-0.5">by {entry.actor}</div>
            {entry.detail && <div className="text-[11px] text-[var(--text-disabled)] mt-0.5">{entry.detail}</div>}
          </div>
        </div>
      ))}
    </div>
  );
}

const meta: Meta<typeof AuditTrailTimeline> = {
  title: "Compliance/AuditTrailTimeline",
  component: AuditTrailTimeline,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof AuditTrailTimeline>;

export const Default: Story = {
  render: () => (
    <div className="max-w-md">
      <AuditTrailTimeline entries={[
        { action: "GSTR-1 Filed", actor: "Priya Sharma", timestamp: "25/04/2026 14:32", status: "success", detail: "Filed via DSC token for April 2026" },
        { action: "Approved by Checker", actor: "Rajesh Kumar (CFO)", timestamp: "25/04/2026 11:15", status: "success" },
        { action: "Submitted for Approval", actor: "Priya Sharma", timestamp: "24/04/2026 16:45", status: "info", detail: "48 invoices, ₹12.4 Cr total taxable value" },
        { action: "Validation Warning", actor: "System", timestamp: "24/04/2026 16:44", status: "warning", detail: "3 invoices have HSN mismatch with purchase register" },
        { action: "Draft Created", actor: "Priya Sharma", timestamp: "22/04/2026 09:30", status: "default" },
        { action: "Period Opened", actor: "System", timestamp: "01/04/2026 00:00", status: "info", detail: "April 2026 filing period auto-opened" },
      ]} />
    </div>
  ),
};
