import type { Meta, StoryObj } from "@storybook/react";

type GovSystem = "GSTN" | "IRP" | "EWB" | "TRACES" | "MCA" | "OLTAS";
type GovStatus = "success" | "warning" | "danger" | "info";

function GovStatusPill({ system, status, label }: {
  system: GovSystem; status: GovStatus; label: string;
}) {
  const dotColor: Record<GovStatus, string> = {
    success: "bg-[var(--success)]",
    warning: "bg-[var(--warning)]",
    danger: "bg-[var(--danger)]",
    info: "bg-[var(--info)]",
  };
  return (
    <span className="inline-flex items-center gap-1.5 h-6 px-2.5 rounded-[6px] border text-[10px] font-semibold bg-[var(--bg-tertiary)] border-[var(--border-default)]">
      <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${dotColor[status]}`} />
      <span className="font-mono uppercase tracking-wide text-[var(--text-primary)]">{system}</span>
      <span className="text-[var(--text-muted)]">·</span>
      <span className="text-[var(--text-secondary)] normal-case">{label}</span>
    </span>
  );
}

const meta: Meta<typeof GovStatusPill> = {
  title: "Compliance/GovStatusPill",
  component: GovStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof GovStatusPill>;

export const Success: Story = { args: { system: "GSTN", status: "success", label: "Filed" } };
export const Warning: Story = { args: { system: "IRP", status: "warning", label: "Pending" } };
export const Danger: Story = { args: { system: "EWB", status: "danger", label: "Expired" } };
export const Info: Story = { args: { system: "TRACES", status: "info", label: "Processing" } };

export const AllSystems: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <GovStatusPill system="GSTN" status="success" label="Filed" />
      <GovStatusPill system="IRP" status="warning" label="Queued" />
      <GovStatusPill system="EWB" status="danger" label="Expired" />
      <GovStatusPill system="TRACES" status="info" label="Processing" />
      <GovStatusPill system="MCA" status="success" label="Acknowledged" />
      <GovStatusPill system="OLTAS" status="info" label="Submitted" />
    </div>
  ),
};
