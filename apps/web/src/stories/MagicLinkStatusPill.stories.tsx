import type { Meta, StoryObj } from "@storybook/react";
import { MagicLinkStatusPill } from "../app/compliance/itr/components/MagicLinkStatusPill";

const meta: Meta<typeof MagicLinkStatusPill> = {
  title: "Compliance/ITR/MagicLinkStatusPill",
  component: MagicLinkStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MagicLinkStatusPill>;

export const Sent: Story = { args: { status: "SENT" } };
export const Viewed: Story = { args: { status: "VIEWED" } };
export const Approved: Story = { args: { status: "APPROVED" } };
export const Expired: Story = { args: { status: "EXPIRED" } };
export const Used: Story = { args: { status: "USED" } };
