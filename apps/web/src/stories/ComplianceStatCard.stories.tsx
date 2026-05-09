import type { Meta, StoryObj } from "@storybook/react";
import { ComplianceStatCard } from "@/app/compliance/calendar/components/ComplianceStatCard";

const meta: Meta<typeof ComplianceStatCard> = {
  title: "Compliance/Calendar/ComplianceStatCard",
  component: ComplianceStatCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ComplianceStatCard>;

export const Filed: Story = { args: { label: "Filed", count: 12, variant: "success" } };
export const DueSoon: Story = { args: { label: "Due in 7 days", count: 3, variant: "warning" } };
export const Upcoming: Story = { args: { label: "Upcoming this month", count: 8, variant: "neutral" } };
