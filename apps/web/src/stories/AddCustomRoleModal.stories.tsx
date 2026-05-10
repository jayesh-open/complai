import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { AddCustomRoleModal } from "../app/configure/users/components/AddCustomRoleModal";
import { MOCK_ROLES } from "../app/configure/users/mock-data";

const meta: Meta<typeof AddCustomRoleModal> = {
  title: "Configure/Users/AddCustomRoleModal",
  component: AddCustomRoleModal,
  tags: ["autodocs"],
  args: {
    onClose: fn(),
    onCreated: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof AddCustomRoleModal>;

export const Open: Story = { args: { open: true, roles: MOCK_ROLES } };
export const Closed: Story = { args: { open: false, roles: MOCK_ROLES } };
