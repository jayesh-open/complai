import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { TaxYearSelector } from "../app/compliance/tds/components/TaxYearSelector";

const meta: Meta<typeof TaxYearSelector> = {
  title: "Compliance/TDS/TaxYearSelector",
  component: TaxYearSelector,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof TaxYearSelector>;

export const Default: Story = {
  render: () => {
    const [value, setValue] = useState("2026-27");
    return <TaxYearSelector value={value} onChange={setValue} />;
  },
};
