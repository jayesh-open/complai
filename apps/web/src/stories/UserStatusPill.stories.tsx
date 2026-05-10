import type { Meta, StoryObj } from "@storybook/react";
import { UserStatusPill } from "../app/configure/users/components/UserStatusPill";

const meta: Meta<typeof UserStatusPill> = {
  title: "Configure/Users/UserStatusPill",
  component: UserStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof UserStatusPill>;

export const Active: Story = { args: { status: "active" } };
export const Inactive: Story = { args: { status: "inactive" } };
