import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react";
import { ComplianceDayDetailPanel } from "@/app/compliance/calendar/components/ComplianceDayDetailPanel";
import type { ComplianceEvent } from "@/app/compliance/calendar/types";

const meta: Meta<typeof ComplianceDayDetailPanel> = {
  title: "Compliance/Calendar/ComplianceDayDetailPanel",
  component: ComplianceDayDetailPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ComplianceDayDetailPanel>;

const sampleEvents: (ComplianceEvent & { dueDate: Date })[] = [
  {
    id: "evt-001",
    title: "GSTR-1 (May)",
    description: "Monthly outward supplies return for invoices issued in May 2026.",
    category: "indirect_tax",
    authority: "CBIC",
    sectionRef: "GST § 37",
    formRef: "GSTR-1",
    dueDateOffset: 0,
    penalty: "₹50/day late fee, ₹20/day for nil returns",
    status: "due_soon",
    linkedModule: "/compliance/gst",
    dueDate: new Date(),
  },
  {
    id: "evt-002",
    title: "TDS Deposit (Monthly)",
    description: "Monthly deposit of tax deducted at source to the government.",
    category: "direct_tax",
    authority: "CBDT",
    sectionRef: "ITA 2025 § 392",
    formRef: "Challan 281",
    dueDateOffset: 0,
    penalty: "1.5% per month interest on delayed payment",
    status: "filed",
    linkedModule: "/compliance/tds",
    dueDate: new Date(),
  },
  {
    id: "evt-003",
    title: "PF Contribution",
    description: "Monthly provident fund contribution deposit. 12% of basic + DA.",
    category: "statutory",
    authority: "EPFO",
    formRef: "ECR",
    dueDateOffset: 0,
    penalty: "Damages @ 5-25% of arrears",
    status: "upcoming",
    linkedModule: "/data-sources/payroll",
    dueDate: new Date(),
  },
];

export const WithEvents: Story = {
  render: () => {
    const [open, setOpen] = useState(true);
    return (
      <ComplianceDayDetailPanel
        date={new Date()}
        events={sampleEvents}
        open={open}
        onClose={() => setOpen(false)}
      />
    );
  },
};

export const Empty: Story = {
  render: () => {
    const [open, setOpen] = useState(true);
    return (
      <ComplianceDayDetailPanel
        date={new Date()}
        events={[]}
        open={open}
        onClose={() => setOpen(false)}
      />
    );
  },
};
