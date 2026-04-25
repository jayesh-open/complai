import type { Meta, StoryObj } from "@storybook/react";
import { GovStatusPill } from "@complai/ui-components";

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
