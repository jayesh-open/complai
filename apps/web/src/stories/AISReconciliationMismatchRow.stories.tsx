import type { Meta, StoryObj } from "@storybook/react";
import { AISReconciliationMismatchRow } from "../app/compliance/itr/components/AISReconciliationMismatchRow";

const meta: Meta<typeof AISReconciliationMismatchRow> = {
  title: "Compliance/ITR/AISReconciliationMismatchRow",
  component: AISReconciliationMismatchRow,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof AISReconciliationMismatchRow>;

export const ErrorUnresolved: Story = {
  args: {
    mismatch: { id: "1", category: "INTEREST", field: "FD Interest (HDFC)", itrValue: 45000, aisValue: 52000, severity: "error", resolved: false },
  },
};

export const WarningUnresolved: Story = {
  args: {
    mismatch: { id: "2", category: "TDS", field: "TDS on Salary", itrValue: 185000, aisValue: 192000, severity: "warn", resolved: false },
  },
};

export const InfoResolved: Story = {
  args: {
    mismatch: { id: "3", category: "SALARY", field: "Gross Salary", itrValue: 1850000, aisValue: 1845000, severity: "info", resolved: true, resolution: "Rounding difference — accepted" },
  },
};
