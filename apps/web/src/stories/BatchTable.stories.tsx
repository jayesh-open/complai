import type { Meta, StoryObj } from "@storybook/react";
import { BatchTable } from "../app/compliance/itr/components/BatchTable";
import { ALL_BATCHES } from "../app/compliance/itr/mock-data";

const meta: Meta<typeof BatchTable> = {
  title: "Compliance/ITR/BatchTable",
  component: BatchTable,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BatchTable>;

export const Default: Story = { args: { batches: ALL_BATCHES } };
export const Empty: Story = { args: { batches: [] } };
export const SingleBatch: Story = { args: { batches: [ALL_BATCHES[0]] } };
