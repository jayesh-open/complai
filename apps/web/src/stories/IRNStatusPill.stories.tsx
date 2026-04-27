import type { Meta, StoryObj } from "@storybook/react";
import { IRNStatusPill } from "../app/compliance/e-invoicing/components/IRNStatusPill";

const meta: Meta<typeof IRNStatusPill> = {
  title: "Compliance/IRNStatusPill",
  component: IRNStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof IRNStatusPill>;

export const Generated: Story = { args: { status: "GENERATED" } };
export const Cancelled: Story = { args: { status: "CANCELLED" } };

export const AllStatuses: Story = {
  render: () => (
    <div className="flex gap-3">
      <IRNStatusPill status="GENERATED" />
      <IRNStatusPill status="CANCELLED" />
    </div>
  ),
};
