import type { Meta, StoryObj } from "@storybook/react";
import { BulkBatchProgressBar } from "../app/compliance/itr/components/BulkBatchProgressBar";

const meta: Meta<typeof BulkBatchProgressBar> = {
  title: "Compliance/ITR/BulkBatchProgressBar",
  component: BulkBatchProgressBar,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BulkBatchProgressBar>;

export const InProgress: Story = { args: { filed: 8, pending: 4, failed: 1, total: 13 } };
export const AllFiled: Story = { args: { filed: 15, pending: 0, failed: 0, total: 15 } };
export const WithFailures: Story = { args: { filed: 10, pending: 0, failed: 5, total: 15 } };
export const JustStarted: Story = { args: { filed: 1, pending: 14, failed: 0, total: 15 } };
export const Empty: Story = { args: { filed: 0, pending: 0, failed: 0, total: 0 } };
