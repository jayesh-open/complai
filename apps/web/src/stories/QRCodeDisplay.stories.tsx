import type { Meta, StoryObj } from "@storybook/react";
import { QRCodeDisplay } from "../app/compliance/e-invoicing/components/QRCodeDisplay";

const meta: Meta<typeof QRCodeDisplay> = {
  title: "Compliance/QRCodeDisplay",
  component: QRCodeDisplay,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof QRCodeDisplay>;

export const Default: Story = {
  args: {
    value: "upi://pay?irn=abc123def456&ack=100000000001",
    size: 160,
    label: "ACK: 100000000001",
  },
};

export const Small: Story = {
  args: {
    value: "upi://pay?irn=abc123def456&ack=100000000001",
    size: 100,
  },
};

export const Large: Story = {
  args: {
    value: "upi://pay?irn=abc123def456&ack=100000000001",
    size: 240,
    label: "Scan to verify e-Invoice",
  },
};
