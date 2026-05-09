import type { Meta, StoryObj } from "@storybook/react";
import { ComplianceItemPill } from "@/app/compliance/calendar/components/ComplianceItemPill";
import type { ComplianceEvent } from "@/app/compliance/calendar/types";

const meta: Meta<typeof ComplianceItemPill> = {
  title: "Compliance/Calendar/ComplianceItemPill",
  component: ComplianceItemPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ComplianceItemPill>;

const base: ComplianceEvent & { dueDate: Date } = {
  id: "evt-001",
  title: "GSTR-1 (May)",
  description: "Monthly outward supplies return.",
  category: "indirect_tax",
  authority: "CBIC",
  sectionRef: "GST § 37",
  formRef: "GSTR-1",
  dueDateOffset: 2,
  penalty: "₹50/day late fee",
  status: "upcoming",
  linkedModule: "/compliance/gst",
  dueDate: new Date(),
};

export const GSTUpcoming: Story = { args: { event: base } };
export const TDSFiled: Story = {
  args: { event: { ...base, id: "evt-002", title: "TDS Deposit (April)", category: "direct_tax", status: "filed" } },
};
export const StatutoryDueSoon: Story = {
  args: { event: { ...base, id: "evt-003", title: "PF Contribution", category: "statutory", status: "due_soon" } },
};
