import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { PermissionMatrix } from "../app/configure/users/components/PermissionMatrix";
import type { Permission } from "../app/configure/users/types";

function makePerm(resource: string, action: string): Permission {
  return { id: `${resource}_${action}`, tenant_id: "t1", resource, action, created_at: "2026-01-01T00:00:00Z" };
}

const ALL_MODULES: Record<string, string[]> = {
  gst_returns: ["view", "edit", "file"],
  gstr_9_9c: ["view", "edit", "file"],
  e_invoicing: ["view", "generate", "cancel"],
  e_way_bill: ["view", "generate", "cancel"],
  itc_reconciliation: ["view", "edit"],
  vendor_compliance: ["view", "edit"],
  tds: ["view", "calculate", "file", "issue_cert"],
  itr: ["view", "calculate", "file", "approve"],
  compliance_calendar: ["view"],
  users_roles: ["view", "manage"],
  connected_apps: ["view", "manage"],
  billing: ["view", "manage"],
};

const adminPerms: Permission[] = Object.entries(ALL_MODULES).flatMap(([mod, actions]) =>
  actions.map((a) => makePerm(mod, a)),
);

const auditorPerms: Permission[] = Object.keys(ALL_MODULES).map((mod) => makePerm(mod, "view"));

const customPerms: Permission[] = [
  makePerm("gst_returns", "view"), makePerm("gst_returns", "edit"), makePerm("gst_returns", "file"),
  makePerm("gstr_9_9c", "view"), makePerm("gstr_9_9c", "edit"),
  makePerm("e_invoicing", "view"), makePerm("e_invoicing", "generate"),
  makePerm("tds", "view"), makePerm("tds", "calculate"),
  makePerm("compliance_calendar", "view"),
];

const meta: Meta<typeof PermissionMatrix> = {
  title: "Configure/Users/PermissionMatrix",
  component: PermissionMatrix,
  tags: ["autodocs"],
  args: { onChange: fn() },
};
export default meta;
type Story = StoryObj<typeof PermissionMatrix>;

export const AdminAllGranted: Story = { args: { permissions: adminPerms, isSystemRole: true } };
export const AuditorViewOnly: Story = { args: { permissions: auditorPerms, isSystemRole: true } };
export const CustomMixed: Story = { args: { permissions: customPerms, isSystemRole: false } };
export const EmptyCustom: Story = { args: { permissions: [], isSystemRole: false } };
