import type { Meta, StoryObj } from "@storybook/react";
import { MonthNavigation } from "@/app/compliance/calendar/components/MonthNavigation";

const meta: Meta<typeof MonthNavigation> = {
  title: "Compliance/Calendar/MonthNavigation",
  component: MonthNavigation,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MonthNavigation>;

export const Default: Story = {
  args: {
    year: 2026,
    month: 4,
    onPrev: () => {},
    onNext: () => {},
    onToday: () => {},
  },
};

export const January: Story = {
  args: {
    year: 2026,
    month: 0,
    onPrev: () => {},
    onNext: () => {},
    onToday: () => {},
  },
};
