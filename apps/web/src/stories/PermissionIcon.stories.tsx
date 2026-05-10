import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { PermissionIcon } from "../app/configure/users/components/PermissionIcon";

const meta: Meta<typeof PermissionIcon> = {
  title: "Configure/Users/PermissionIcon",
  component: PermissionIcon,
  tags: ["autodocs"],
  args: { onClick: fn() },
};
export default meta;
type Story = StoryObj<typeof PermissionIcon>;

export const ViewGranted: Story = { args: { action: "view", granted: true, disabled: false } };
export const ViewNotGranted: Story = { args: { action: "view", granted: false, disabled: false } };
export const EditGranted: Story = { args: { action: "edit", granted: true, disabled: false } };
export const FileGranted: Story = { args: { action: "file", granted: true, disabled: false } };
export const GenerateGranted: Story = { args: { action: "generate", granted: true, disabled: false } };
export const CancelGranted: Story = { args: { action: "cancel", granted: true, disabled: false } };
export const CalculateGranted: Story = { args: { action: "calculate", granted: true, disabled: false } };
export const ApproveGranted: Story = { args: { action: "approve", granted: true, disabled: false } };
export const IssueCertGranted: Story = { args: { action: "issue_cert", granted: true, disabled: false } };
export const ManageGranted: Story = { args: { action: "manage", granted: true, disabled: false } };
export const SystemDisabled: Story = { args: { action: "view", granted: true, disabled: true } };
