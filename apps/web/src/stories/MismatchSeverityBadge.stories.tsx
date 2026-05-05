import type { Meta, StoryObj } from "@storybook/react";
import { MismatchSeverityBadge } from "../app/compliance/gst-returns/annual/components/MismatchSeverityBadge";

const meta: Meta<typeof MismatchSeverityBadge> = {
  title: "Compliance/GST/MismatchSeverityBadge",
  component: MismatchSeverityBadge,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MismatchSeverityBadge>;

export const Info: Story = { args: { severity: "INFO" } };
export const Warn: Story = { args: { severity: "WARN" } };
export const Error: Story = { args: { severity: "ERROR" } };
export const AllSeverities: Story = {
  render: () => (
    <div className="flex items-center gap-3">
      <MismatchSeverityBadge severity="INFO" />
      <MismatchSeverityBadge severity="WARN" />
      <MismatchSeverityBadge severity="ERROR" />
    </div>
  ),
};
