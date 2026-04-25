import type { Meta, StoryObj } from "@storybook/react";
import { MakerCheckerApprovalCard } from "@complai/ui-components";

const meta: Meta<typeof MakerCheckerApprovalCard> = {
  title: "Compliance/MakerCheckerApprovalCard",
  component: MakerCheckerApprovalCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MakerCheckerApprovalCard>;

export const Default: Story = {
  args: {
    title: "GSTR-3B Filing — April 2026",
    subtitle: "29AABCA1234A1Z5 · Karnataka",
    submittedBy: "Priya Sharma",
    submittedAt: "25/04/2026 14:30",
    details: [
      { label: "Period", value: "April 2026" },
      { label: "Tax Payable", value: "₹12,45,678" },
      { label: "ITC Claimed", value: "₹8,92,340" },
      { label: "Net Payable", value: "₹3,53,338" },
    ],
    warnings: [
      "ITC claimed exceeds auto-populated 2B amount by ₹42,000",
      "Late filing penalty of ₹2,000 may apply after 20th",
    ],
    onApprove: (c: string) => alert(`Approved${c ? `: ${c}` : ""}`),
    onReject: (c: string) => alert(`Rejected${c ? `: ${c}` : ""}`),
    onSendBack: (c: string) => alert(`Sent back${c ? `: ${c}` : ""}`),
  },
};

export const NoWarnings: Story = {
  args: {
    title: "TDS Return — Q4 FY 2025-26",
    subtitle: "AAACR1234A · Form 26Q",
    submittedBy: "Amit Verma",
    submittedAt: "24/04/2026 10:15",
    details: [
      { label: "Quarter", value: "Q4 (Jan-Mar 2026)" },
      { label: "Deductees", value: "142" },
      { label: "Total TDS", value: "₹18,72,400" },
    ],
    onApprove: (c: string) => alert(`Approved${c ? `: ${c}` : ""}`),
    onReject: (c: string) => alert(`Rejected${c ? `: ${c}` : ""}`),
  },
};
