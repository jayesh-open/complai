import type { Meta, StoryObj } from "@storybook/react";
import { MismatchRow } from "../app/compliance/gst-returns/annual/components/MismatchRow";
import type { GSTR9CMismatch } from "../app/compliance/gst-returns/annual/types";

const meta: Meta<typeof MismatchRow> = {
  title: "Compliance/GST/MismatchRow",
  component: MismatchRow,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MismatchRow>;

const baseMismatch: GSTR9CMismatch = {
  id: "m-1",
  section: "II",
  category: "turnover",
  description: "Aggregate turnover mismatch (books vs GSTR-9 Part II)",
  booksAmount: 86700000,
  gstr9Amount: 85000000,
  difference: 1700000,
  severity: "ERROR",
  resolved: false,
};

export const UnresolvedError: Story = {
  args: { mismatch: baseMismatch, onResolve: () => {} },
};
export const UnresolvedWarn: Story = {
  args: { mismatch: { ...baseMismatch, id: "m-2", severity: "WARN", difference: 5000 }, onResolve: () => {} },
};
export const ResolvedInfo: Story = {
  args: {
    mismatch: { ...baseMismatch, id: "m-3", severity: "INFO", resolved: true, resolvedReason: "Timing difference" },
    onResolve: () => {},
  },
};
