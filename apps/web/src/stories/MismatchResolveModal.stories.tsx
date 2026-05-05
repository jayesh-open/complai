import type { Meta, StoryObj } from "@storybook/react";
import { MismatchResolveModal } from "../app/compliance/gst-returns/annual/components/MismatchResolveModal";
import type { GSTR9CMismatch } from "../app/compliance/gst-returns/annual/types";

const meta: Meta<typeof MismatchResolveModal> = {
  title: "Compliance/GST/MismatchResolveModal",
  component: MismatchResolveModal,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MismatchResolveModal>;

const errorMismatch: GSTR9CMismatch = {
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

const warnMismatch: GSTR9CMismatch = {
  id: "m-2",
  section: "III",
  category: "tax",
  description: "CGST payable variance at 18% slab",
  booksAmount: 2720000,
  gstr9Amount: 2550000,
  difference: 170000,
  severity: "WARN",
  resolved: false,
};

export const ErrorSeverity: Story = {
  args: { mismatch: errorMismatch, open: true, onClose: () => {}, onResolve: () => {} },
};
export const WarnSeverity: Story = {
  args: { mismatch: warnMismatch, open: true, onClose: () => {}, onResolve: () => {} },
};
export const Closed: Story = {
  args: { mismatch: errorMismatch, open: false, onClose: () => {}, onResolve: () => {} },
};
