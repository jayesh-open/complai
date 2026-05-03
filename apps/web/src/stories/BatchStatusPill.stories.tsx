import type { Meta, StoryObj } from "@storybook/react";
import { BatchStatusPill } from "../app/compliance/itr/components/BatchStatusPill";

const meta: Meta<typeof BatchStatusPill> = {
  title: "Compliance/ITR/BatchStatusPill",
  component: BatchStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BatchStatusPill>;

export const Draft: Story = { args: { status: "DRAFT" } };
export const InProgress: Story = { args: { status: "IN_PROGRESS" } };
export const Completed: Story = { args: { status: "COMPLETED" } };
export const Failed: Story = { args: { status: "FAILED" } };
