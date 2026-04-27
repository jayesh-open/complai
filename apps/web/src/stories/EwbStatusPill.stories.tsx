import type { Meta, StoryObj } from "@storybook/react";
import { EwbStatusPill } from "../app/compliance/e-way-bill/components/EwbStatusPill";

const meta: Meta<typeof EwbStatusPill> = {
  title: "Compliance/EwbStatusPill",
  component: EwbStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof EwbStatusPill>;

export const Active: Story = { args: { status: "ACTIVE" } };
export const Expired: Story = { args: { status: "EXPIRED" } };
export const Cancelled: Story = { args: { status: "CANCELLED" } };
export const Consolidated: Story = { args: { status: "CONSOLIDATED" } };
export const NearingExpiry: Story = { args: { status: "ACTIVE", nearingExpiry: true } };

export const AllStatuses: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <EwbStatusPill status="ACTIVE" />
      <EwbStatusPill status="ACTIVE" nearingExpiry />
      <EwbStatusPill status="EXPIRED" />
      <EwbStatusPill status="CANCELLED" />
      <EwbStatusPill status="CONSOLIDATED" />
    </div>
  ),
};
