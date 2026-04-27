import type { Meta, StoryObj } from "@storybook/react";
import { CancellationCountdown } from "../app/compliance/e-invoicing/components/CancellationCountdown";

const meta: Meta<typeof CancellationCountdown> = {
  title: "Compliance/CancellationCountdown",
  component: CancellationCountdown,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof CancellationCountdown>;

export const PlentyOfTime: Story = {
  args: {
    generatedAt: new Date(Date.now() - 2 * 3600 * 1000).toISOString(),
  },
};

export const Urgent: Story = {
  args: {
    generatedAt: new Date(Date.now() - 23.5 * 3600 * 1000).toISOString(),
  },
};

export const Expired: Story = {
  args: {
    generatedAt: new Date(Date.now() - 25 * 3600 * 1000).toISOString(),
  },
};
