import type { Meta, StoryObj } from "@storybook/react";
import { BulkOperationTray } from "@complai/ui-components";

const meta: Meta<typeof BulkOperationTray> = {
  title: "Compliance/BulkOperationTray",
  component: BulkOperationTray,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BulkOperationTray>;

export const Default: Story = {
  args: {
    jobs: [
      { id: "1", title: "E-Invoice Generation (Batch 42)", progress: 142, total: 228, status: "running", eta: "~3 min" },
      { id: "2", title: "GSTR-1 JSON Export", progress: 500, total: 500, status: "done" },
      { id: "3", title: "Vendor GSTIN Validation", progress: 38, total: 120, status: "running" },
    ],
    onStop: (id: string) => alert(`Stopping job ${id}`),
  },
};

export const WithError: Story = {
  args: {
    jobs: [
      { id: "1", title: "TDS Certificate Download", progress: 45, total: 200, status: "error" },
      { id: "2", title: "Bulk Invoice Upload", progress: 180, total: 300, status: "running" },
    ],
  },
};

export const AllDone: Story = {
  args: {
    jobs: [
      { id: "1", title: "E-Invoice Generation", progress: 228, total: 228, status: "done" },
      { id: "2", title: "GSTR-1 JSON Export", progress: 500, total: 500, status: "done" },
    ],
  },
};
