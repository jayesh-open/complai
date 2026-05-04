import type { Meta, StoryObj } from "@storybook/react";
import { LateITCReclaimPanel } from "../app/compliance/gst-returns/annual/components/LateITCReclaimPanel";
import type { LateITCEntry } from "../app/compliance/gst-returns/annual/types";

const entries: LateITCEntry[] = [
  { table: "6H", description: "ITC reclaimed on reversal of Rule 37 (payment within 180 days)", amount: 170000, period: "Q2 2025-26", rule: "Rule 37 proviso — ITC reversed for non-payment can be reclaimed when payment is made within the extended period" },
  { table: "8C", description: "Difference between ITC claimed in GSTR-3B vs GSTR-2B (gap rectification)", amount: 85000, period: "Sep 2025", rule: "Section 16(4) — ITC missed in monthly returns can be claimed in annual return up to Sep 30 of next FY" },
  { table: "13", description: "ITC declared in current FY relating to previous FY invoices", amount: 127500, period: "FY 2024-25 invoices", rule: "Section 16(4) read with GSTR-9 Table 13 — ITC on prior-year invoices claimed in current year annual" },
];

const meta: Meta<typeof LateITCReclaimPanel> = {
  title: "Compliance/GST/LateITCReclaimPanel",
  component: LateITCReclaimPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof LateITCReclaimPanel>;

export const Default: Story = { args: { entries } };
export const SingleEntry: Story = { args: { entries: [entries[0]] } };
export const HighValue: Story = {
  args: {
    entries: entries.map((e) => ({ ...e, amount: e.amount * 10 })),
  },
};
