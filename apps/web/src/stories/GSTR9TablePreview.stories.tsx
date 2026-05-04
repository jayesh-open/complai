import type { Meta, StoryObj } from "@storybook/react";
import { GSTR9TablePreview } from "../app/compliance/gst-returns/annual/components/GSTR9TablePreview";
import type { GSTR9Table } from "../app/compliance/gst-returns/annual/types";

const table4: GSTR9Table = {
  tableNumber: 4,
  title: "Details of advances, inward and outward supplies on which tax is payable",
  rows: [
    { serial: "4A", description: "Supplies made to un-registered persons (B2C)", taxableValue: 25500000, cgst: 2295000, sgst: 2295000, igst: 0, cess: 0, sourceReturn: "GSTR-1" },
    { serial: "4B", description: "Supplies made to registered persons (B2B)", taxableValue: 42500000, cgst: 2550000, sgst: 2550000, igst: 2550000, cess: 0, sourceReturn: "GSTR-1" },
    { serial: "4C", description: "Zero rated supply (Export) on payment of tax", taxableValue: 6800000, cgst: 0, sgst: 0, igst: 1190000, cess: 0, sourceReturn: "GSTR-1" },
  ],
};

const table7: GSTR9Table = {
  tableNumber: 7,
  title: "Details of ITC reversed and ineligible ITC",
  rows: [
    { serial: "7A", description: "As per Rule 37", taxableValue: 0, cgst: 85000, sgst: 85000, igst: 0, cess: 0 },
    { serial: "7B", description: "As per Rule 39", taxableValue: 0, cgst: 42500, sgst: 42500, igst: 0, cess: 0 },
    { serial: "7H", description: "Other reversals", taxableValue: 0, cgst: 68000, sgst: 68000, igst: 85000, cess: 0 },
  ],
};

const meta: Meta<typeof GSTR9TablePreview> = {
  title: "Compliance/GST/GSTR9TablePreview",
  component: GSTR9TablePreview,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof GSTR9TablePreview>;

export const OutwardSupplies: Story = { args: { table: table4 } };
export const ITCReversals: Story = { args: { table: table7 } };
export const SingleRow: Story = {
  args: {
    table: {
      tableNumber: 9,
      title: "Details of tax paid as declared in returns filed during the FY",
      rows: [{ serial: "9", description: "Tax paid through cash ledger", taxableValue: 0, cgst: 680000, sgst: 680000, igst: 1020000, cess: 0, sourceReturn: "GSTR-3B" }],
    },
  },
};
