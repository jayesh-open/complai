import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { UserDetailPanel } from "../app/configure/users/components/UserDetailPanel";
import { MOCK_USERS, MOCK_ROLES } from "../app/configure/users/mock-data";

const meta: Meta<typeof UserDetailPanel> = {
  title: "Configure/Users/UserDetailPanel",
  component: UserDetailPanel,
  tags: ["autodocs"],
  args: {
    open: true,
    roles: MOCK_ROLES,
    onClose: fn(),
    onUpdated: fn(),
  },
  decorators: [
    (Story) => (
      <div style={{ height: "600px", position: "relative", overflow: "hidden" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof UserDetailPanel>;

export const ReadMode: Story = { args: { user: MOCK_USERS[0]! } };

export const DeactivatedUser: Story = { args: { user: MOCK_USERS[5]! } };

export const NoRole: Story = { args: { user: MOCK_USERS[8]! } };
