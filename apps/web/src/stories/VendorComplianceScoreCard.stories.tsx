import type { Meta, StoryObj } from "@storybook/react";
import { VendorComplianceScoreCard } from "@complai/ui-components";

const meta: Meta<typeof VendorComplianceScoreCard> = {
  title: "Compliance/VendorComplianceScoreCard",
  component: VendorComplianceScoreCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof VendorComplianceScoreCard>;

export const HighScore: Story = {
  args: {
    vendorName: "Tata Steel Limited",
    gstin: "29AABCT1234A1Z5",
    state: "Karnataka",
    totalScore: 92,
    riskLevel: "Low",
    category: "Strategic",
    lastReviewed: "20/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 20, maxScore: 20, status: "pass" },
      { label: "GSTR-2B Match Rate", score: 18, maxScore: 20, status: "pass", note: "98.2% match" },
      { label: "Payment Timeliness", score: 19, maxScore: 20, status: "pass" },
      { label: "Document Quality", score: 17, maxScore: 20, status: "pass", note: "All e-invoiced" },
      { label: "Compliance History", score: 18, maxScore: 20, status: "pass" },
    ],
  },
};

export const MediumScore: Story = {
  args: {
    vendorName: "QuickParts Industries",
    gstin: "27AADCQ5678B2Z3",
    state: "Maharashtra",
    totalScore: 64,
    riskLevel: "Medium",
    category: "Operational",
    lastReviewed: "18/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 14, maxScore: 20, status: "warn", note: "2 late filings" },
      { label: "GSTR-2B Match Rate", score: 12, maxScore: 20, status: "warn", note: "76% match" },
      { label: "Payment Timeliness", score: 16, maxScore: 20, status: "pass" },
      { label: "Document Quality", score: 10, maxScore: 20, status: "fail", note: "Manual invoices" },
      { label: "Compliance History", score: 12, maxScore: 20, status: "warn" },
    ],
  },
};

export const LowScore: Story = {
  args: {
    vendorName: "ABC Traders",
    gstin: "33AABCA9999C1Z1",
    state: "Tamil Nadu",
    totalScore: 38,
    riskLevel: "High",
    category: "Transactional",
    lastReviewed: "15/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 6, maxScore: 20, status: "fail", note: "4 months unfiled" },
      { label: "GSTR-2B Match Rate", score: 8, maxScore: 20, status: "fail", note: "42% match" },
      { label: "Payment Timeliness", score: 10, maxScore: 20, status: "warn" },
      { label: "Document Quality", score: 8, maxScore: 20, status: "fail", note: "Missing HSN codes" },
      { label: "Compliance History", score: 6, maxScore: 20, status: "fail", note: "SCN received" },
    ],
  },
};
