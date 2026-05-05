import type { Meta, StoryObj } from "@storybook/react";
import { ReconciliationDeltaTable } from "../app/compliance/gst-returns/annual/components/ReconciliationDeltaTable";
import type { GSTR9CMismatch } from "../app/compliance/gst-returns/annual/types";

const meta: Meta<typeof ReconciliationDeltaTable> = {
  title: "Compliance/GST/ReconciliationDeltaTable",
  component: ReconciliationDeltaTable,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ReconciliationDeltaTable>;

const mockMismatches: GSTR9CMismatch[] = [
  { id: "m-1", section: "II", category: "turnover", description: "Aggregate turnover mismatch", booksAmount: 86700000, gstr9Amount: 85000000, difference: 1700000, severity: "ERROR", resolved: false },
  { id: "m-2", section: "II", category: "turnover", description: "Unbilled revenue adjustment", booksAmount: 680000, gstr9Amount: 0, difference: 680000, severity: "WARN", resolved: false },
  { id: "m-3", section: "III", category: "tax", description: "CGST payable at 18% slab", booksAmount: 2720000, gstr9Amount: 2550000, difference: 170000, severity: "WARN", resolved: false },
  { id: "m-4", section: "IV", category: "itc", description: "ITC claimed (IGST) — import duty credit", booksAmount: 1530000, gstr9Amount: 1275000, difference: 255000, severity: "ERROR", resolved: false },
];

export const TurnoverSection: Story = {
  args: { section: "II", sectionTitle: "Turnover Reconciliation", mismatches: mockMismatches },
};
export const TaxSection: Story = {
  args: { section: "III", sectionTitle: "Tax Reconciliation", mismatches: mockMismatches },
};
export const ITCSection: Story = {
  args: { section: "IV", sectionTitle: "ITC Reconciliation", mismatches: mockMismatches },
};
export const EmptySection: Story = {
  args: { section: "III", sectionTitle: "Tax Reconciliation", mismatches: mockMismatches.filter((m) => m.section === "II") },
};
