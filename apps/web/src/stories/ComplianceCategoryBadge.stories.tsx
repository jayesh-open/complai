import type { Meta, StoryObj } from "@storybook/react";
import { ComplianceCategoryBadge } from "@/app/compliance/calendar/components/ComplianceCategoryBadge";

const meta: Meta<typeof ComplianceCategoryBadge> = {
  title: "Compliance/Calendar/ComplianceCategoryBadge",
  component: ComplianceCategoryBadge,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ComplianceCategoryBadge>;

export const DirectTax: Story = { args: { category: "direct_tax", active: true } };
export const IndirectTax: Story = { args: { category: "indirect_tax", active: true } };
export const Statutory: Story = { args: { category: "statutory", active: true } };
export const Inactive: Story = { args: { category: "direct_tax", active: false } };
export const Small: Story = { args: { category: "indirect_tax", active: true, size: "small" } };
