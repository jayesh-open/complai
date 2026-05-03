import type { Meta, StoryObj } from "@storybook/react";
import { RefundIndicator } from "../app/compliance/itr/components/RefundIndicator";

const meta: Meta<typeof RefundIndicator> = {
  title: "Compliance/ITR/RefundIndicator",
  component: RefundIndicator,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof RefundIndicator>;

export const Refund: Story = { args: { amount: -45000 } };
export const Payable: Story = { args: { amount: 85000 } };
export const Zero: Story = { args: { amount: 0 } };
export const LargeRefund: Story = { args: { amount: -325000 } };
