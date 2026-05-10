import type { Meta, StoryObj } from "@storybook/react";
import { ViewAsBanner } from "../components/layout/ViewAsBanner";
import { useViewAsStore } from "../store/view-as-store";
import { useEffect } from "react";
import type { Role } from "../app/configure/users/types";

const AUDITOR: Role = { id: "r7", tenant_id: "t1", name: "auditor", display_name: "Auditor", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" };
const AP_MGR: Role = { id: "r3", tenant_id: "t1", name: "ap_manager", display_name: "AP Manager", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" };

function WithViewAs({ role, children }: { role: Role; children: React.ReactNode }) {
  const enterViewAs = useViewAsStore((s) => s.enterViewAs);
  useEffect(() => { enterViewAs(role, []); return () => useViewAsStore.getState().exitViewAs(); }, [role, enterViewAs]);
  return <>{children}</>;
}

const meta: Meta<typeof ViewAsBanner> = {
  title: "Layout/ViewAsBanner",
  component: ViewAsBanner,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ViewAsBanner>;

export const AsAuditor: Story = {
  decorators: [(Story) => <WithViewAs role={AUDITOR}><Story /></WithViewAs>],
};
export const AsAPManager: Story = {
  decorators: [(Story) => <WithViewAs role={AP_MGR}><Story /></WithViewAs>],
};
