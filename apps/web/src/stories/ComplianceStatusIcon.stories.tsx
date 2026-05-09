import type { Meta, StoryObj } from "@storybook/react";
import { ComplianceStatusIcon, StatusLegend } from "@/app/compliance/calendar/components/ComplianceStatusIcon";

const meta: Meta<typeof ComplianceStatusIcon> = {
  title: "Compliance/Calendar/ComplianceStatusIcon",
  component: ComplianceStatusIcon,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ComplianceStatusIcon>;

export const Filed: Story = { args: { status: "filed", size: 20 } };
export const DueSoon: Story = { args: { status: "due_soon", size: 20 } };
export const Upcoming: Story = { args: { status: "upcoming", size: 20 } };
export const Overdue: Story = { args: { status: "overdue", size: 20 } };

export const Legend: StoryObj<typeof StatusLegend> = {
  render: () => <StatusLegend />,
};
