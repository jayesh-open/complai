import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { PermissionCell } from "../app/configure/users/components/PermissionCell";

const meta: Meta<typeof PermissionCell> = {
  title: "Configure/Users/PermissionCell",
  component: PermissionCell,
  tags: ["autodocs"],
  args: { onChange: fn() },
};
export default meta;
type Story = StoryObj<typeof PermissionCell>;

export const GstAllGranted: Story = {
  args: { module: "gst_returns", actions: ["view", "edit", "file"], grantedActions: ["view", "edit", "file"], isSystemRole: false },
};
export const GstViewOnly: Story = {
  args: { module: "gst_returns", actions: ["view", "edit", "file"], grantedActions: ["view"], isSystemRole: false },
};
export const TdsAllGranted: Story = {
  args: { module: "tds", actions: ["view", "calculate", "file", "issue_cert"], grantedActions: ["view", "calculate", "file", "issue_cert"], isSystemRole: false },
};
export const EInvoicing: Story = {
  args: { module: "e_invoicing", actions: ["view", "generate", "cancel"], grantedActions: ["view", "generate"], isSystemRole: false },
};
export const SystemReadOnly: Story = {
  args: { module: "gst_returns", actions: ["view", "edit", "file"], grantedActions: ["view", "edit", "file"], isSystemRole: true },
};
