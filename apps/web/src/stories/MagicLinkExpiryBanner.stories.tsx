import type { Meta, StoryObj } from "@storybook/react";
import { MagicLinkExpiryBanner } from "../app/compliance/itr/components/MagicLinkExpiryBanner";

const meta: Meta<typeof MagicLinkExpiryBanner> = {
  title: "Compliance/ITR/MagicLinkExpiryBanner",
  component: MagicLinkExpiryBanner,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MagicLinkExpiryBanner>;

const inOneWeek = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString();
const inTwoHours = new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString();
const pastDate = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString();

export const Normal: Story = { args: { expiresAt: inOneWeek } };
export const Urgent: Story = { args: { expiresAt: inTwoHours } };
export const Expired: Story = { args: { expiresAt: pastDate } };
