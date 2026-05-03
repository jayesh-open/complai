import type { Meta, StoryObj } from "@storybook/react";
import { AISReconciliationPanel } from "../app/compliance/itr/components/AISReconciliationPanel";
import type { AISMismatch } from "../app/compliance/itr/types";

const meta: Meta<typeof AISReconciliationPanel> = {
  title: "Compliance/ITR/AISReconciliationPanel",
  component: AISReconciliationPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof AISReconciliationPanel>;

const mismatches: AISMismatch[] = [
  { id: "1", category: "SALARY", field: "Gross Salary", itrValue: 1850000, aisValue: 1845000, severity: "info", resolved: true, resolution: "Rounding difference" },
  { id: "2", category: "TDS", field: "TDS on Salary", itrValue: 185000, aisValue: 192000, severity: "warn", resolved: false },
  { id: "3", category: "INTEREST", field: "FD Interest (HDFC)", itrValue: 45000, aisValue: 52000, severity: "error", resolved: false },
  { id: "4", category: "DIVIDEND", field: "Dividend — Reliance", itrValue: 12000, aisValue: 12000, severity: "info", resolved: true, resolution: "Exact match" },
];

export const WithErrors: Story = { args: { mismatches } };

export const AllResolved: Story = {
  args: { mismatches: mismatches.map((m) => ({ ...m, resolved: true, resolution: "Accepted" })) },
};

export const Empty: Story = { args: { mismatches: [] } };
