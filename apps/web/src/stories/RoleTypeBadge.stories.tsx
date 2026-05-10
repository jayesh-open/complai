import type { Meta, StoryObj } from "@storybook/react";
import { RoleTypeBadge } from "../app/configure/users/components/RoleTypeBadge";

const meta: Meta<typeof RoleTypeBadge> = {
  title: "Configure/Users/RoleTypeBadge",
  component: RoleTypeBadge,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof RoleTypeBadge>;

export const System: Story = { args: { isSystem: true } };
export const Custom: Story = { args: { isSystem: false } };
