import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { PeriodSelector } from "@complai/ui-components";

const meta: Meta<typeof PeriodSelector> = {
  title: "Compliance/PeriodSelector",
  component: PeriodSelector,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof PeriodSelector>;

export const Default: Story = {
  render: () => {
    const [fy, setFy] = useState("2026-27");
    const [month, setMonth] = useState("Apr");
    return <PeriodSelector financialYear={fy} month={month} onYearChange={setFy} onMonthChange={setMonth} />;
  },
};
