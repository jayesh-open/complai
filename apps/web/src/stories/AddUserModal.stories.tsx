import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { AddUserModal } from "../app/configure/users/components/AddUserModal";
import { MOCK_ROLES } from "../app/configure/users/mock-data";

const meta: Meta<typeof AddUserModal> = {
  title: "Configure/Users/AddUserModal",
  component: AddUserModal,
  tags: ["autodocs"],
  args: {
    open: true,
    roles: MOCK_ROLES,
    onClose: fn(),
    onCreated: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof AddUserModal>;

export const EmptyForm: Story = {};

export const NoRoles: Story = { args: { roles: [] } };
