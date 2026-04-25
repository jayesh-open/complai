import type { Meta, StoryObj } from "@storybook/react";
import { StatusBadge } from "@complai/ui-components";

const meta: Meta<typeof StatusBadge> = {
  title: "Compliance/StatusBadge",
  component: StatusBadge,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof StatusBadge>;

export const AllVariants: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <StatusBadge variant="success">Filed</StatusBadge>
      <StatusBadge variant="danger">Rejected</StatusBadge>
      <StatusBadge variant="warning">Pending</StatusBadge>
      <StatusBadge variant="info">Processing</StatusBadge>
      <StatusBadge variant="purple">RCM</StatusBadge>
      <StatusBadge variant="teal">MSME</StatusBadge>
      <StatusBadge variant="default">Draft</StatusBadge>
    </div>
  ),
};

export const XSSize: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <StatusBadge variant="success" size="xs">Direct ✓</StatusBadge>
      <StatusBadge variant="warning" size="xs">DUP</StatusBadge>
      <StatusBadge variant="info" size="xs">2A Match</StatusBadge>
      <StatusBadge variant="default" size="xs">NON-PO</StatusBadge>
    </div>
  ),
};
