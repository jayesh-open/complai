import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { UsersTable } from "../app/configure/users/components/UsersTable";
import { MOCK_USERS } from "../app/configure/users/mock-data";

const meta: Meta<typeof UsersTable> = {
  title: "Configure/Users/UsersTable",
  component: UsersTable,
  tags: ["autodocs"],
  args: {
    onSelectUser: fn(),
    onDeactivate: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof UsersTable>;

export const Populated: Story = { args: { users: MOCK_USERS } };
export const Empty: Story = { args: { users: [] } };
export const Loading: Story = {
  args: { users: [] },
  decorators: [
    (Story) => (
      <div>
        <p className="text-xs text-[var(--text-muted)] mb-2">
          Loading state is rendered by the page orchestrator, not the table. See the page for skeleton.
        </p>
        <Story />
      </div>
    ),
  ],
};
