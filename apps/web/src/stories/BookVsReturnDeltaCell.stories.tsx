import type { Meta, StoryObj } from "@storybook/react";
import { BookVsReturnDeltaCell } from "../app/compliance/gst-returns/annual/components/BookVsReturnDeltaCell";

const meta: Meta<typeof BookVsReturnDeltaCell> = {
  title: "Compliance/GST/BookVsReturnDeltaCell",
  component: BookVsReturnDeltaCell,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof BookVsReturnDeltaCell>;

export const PositiveDelta: Story = {
  args: { booksAmount: 85000000, returnAmount: 83000000 },
};
export const NegativeDelta: Story = {
  args: { booksAmount: 45000000, returnAmount: 46500000 },
};
export const ZeroDelta: Story = {
  args: { booksAmount: 72000000, returnAmount: 72000000 },
};
export const SmallAmounts: Story = {
  args: { booksAmount: 15200, returnAmount: 14800 },
};
