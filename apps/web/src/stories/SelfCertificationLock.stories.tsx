import type { Meta, StoryObj } from "@storybook/react";
import { SelfCertificationLock } from "../app/compliance/gst-returns/annual/components/SelfCertificationLock";

const meta: Meta<typeof SelfCertificationLock> = {
  title: "Compliance/GST/SelfCertificationLock",
  component: SelfCertificationLock,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof SelfCertificationLock>;

export const Unlocked: Story = {
  args: {
    gstin: "29AABCA1234A1Z5",
    fy: "2025-26",
    totalMismatches: 6,
    resolvedCount: 6,
    onCertify: () => {},
    locked: false,
  },
};
export const Locked: Story = {
  args: {
    gstin: "29AABCA1234A1Z5",
    fy: "2025-26",
    totalMismatches: 6,
    resolvedCount: 6,
    onCertify: () => {},
    locked: true,
  },
};
export const PartiallyResolved: Story = {
  args: {
    gstin: "33CCCDC9012C3X7",
    fy: "2025-26",
    totalMismatches: 6,
    resolvedCount: 4,
    onCertify: () => {},
    locked: false,
  },
};
