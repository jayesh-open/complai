import type { Meta, StoryObj } from "@storybook/react";
import { AuditTrailTimeline } from "@complai/ui-components";

const meta: Meta<typeof AuditTrailTimeline> = {
  title: "Compliance/AuditTrailTimeline",
  component: AuditTrailTimeline,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof AuditTrailTimeline>;

export const Default: Story = {
  render: () => (
    <div className="max-w-md">
      <AuditTrailTimeline entries={[
        { action: "GSTR-1 Filed", actor: "Priya Sharma", timestamp: new Date(2026, 3, 25, 14, 32), status: "success", detail: "Filed via DSC token for April 2026" },
        { action: "Approved by Checker", actor: "Rajesh Kumar (CFO)", timestamp: new Date(2026, 3, 25, 11, 15), status: "success" },
        { action: "Submitted for Approval", actor: "Priya Sharma", timestamp: new Date(2026, 3, 24, 16, 45), status: "info", detail: "48 invoices, ₹12.4 Cr total taxable value" },
        { action: "Validation Warning", actor: "System", timestamp: new Date(2026, 3, 24, 16, 44), status: "warning", detail: "3 invoices have HSN mismatch with purchase register" },
        { action: "Draft Created", actor: "Priya Sharma", timestamp: new Date(2026, 3, 22, 9, 30), status: "default" },
        { action: "Period Opened", actor: "System", timestamp: new Date(2026, 3, 1), status: "info", detail: "April 2026 filing period auto-opened" },
      ]} />
    </div>
  ),
};
