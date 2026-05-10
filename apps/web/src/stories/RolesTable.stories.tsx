import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { RolesTable } from "../app/configure/users/components/RolesTable";
import { MOCK_ROLES } from "../app/configure/users/mock-data";

const meta: Meta<typeof RolesTable> = {
  title: "Configure/Users/RolesTable",
  component: RolesTable,
  tags: ["autodocs"],
  args: {
    onDelete: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof RolesTable>;

export const AllRoles: Story = { args: { roles: MOCK_ROLES } };
export const SystemOnly: Story = { args: { roles: MOCK_ROLES.filter((r) => r.is_system) } };
export const Empty: Story = { args: { roles: [] } };
