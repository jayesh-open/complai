import type { Meta, StoryObj } from "@storybook/react";
import { HSNSummaryTable } from "../app/compliance/gst-returns/annual/components/HSNSummaryTable";
import type { HSNRow } from "../app/compliance/gst-returns/annual/types";

const rows: HSNRow[] = [
  { hsn: "8471", description: "Computers and peripherals", uqc: "NOS", quantity: 2500, taxableValue: 21250000, cgst: 1912500, sgst: 1912500, igst: 850000, digitTier: 4 },
  { hsn: "998314", description: "IT consulting services", uqc: "NOS", quantity: 1200, taxableValue: 29750000, cgst: 2677500, sgst: 2677500, igst: 1700000, digitTier: 6 },
  { hsn: "85176290", description: "Network equipment", uqc: "NOS", quantity: 800, taxableValue: 12750000, cgst: 1147500, sgst: 1147500, igst: 680000, digitTier: 8 },
  { hsn: "9983", description: "Other professional services", uqc: "NOS", quantity: 450, taxableValue: 10200000, cgst: 918000, sgst: 918000, igst: 425000, digitTier: 4 },
];

const meta: Meta<typeof HSNSummaryTable> = {
  title: "Compliance/GST/HSNSummaryTable",
  component: HSNSummaryTable,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof HSNSummaryTable>;

export const Default: Story = { args: { rows } };
export const SingleItem: Story = { args: { rows: [rows[0]] } };
export const AllTiers: Story = { args: { rows } };
