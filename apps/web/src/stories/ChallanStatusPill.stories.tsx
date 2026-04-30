import type { Meta, StoryObj } from "@storybook/react";
import { ChallanStatusPill } from "../app/compliance/tds/components/ChallanStatusPill";

const meta: Meta<typeof ChallanStatusPill> = {
  title: "Compliance/TDS/Challans/ChallanStatusPill",
  component: ChallanStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ChallanStatusPill>;

export const Pending: Story = { args: { status: "PENDING" } };
export const Cleared: Story = { args: { status: "CLEARED" } };
export const Rejected: Story = { args: { status: "REJECTED" } };
