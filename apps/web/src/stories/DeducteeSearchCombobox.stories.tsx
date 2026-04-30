import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { DeducteeSearchCombobox } from "../app/compliance/tds/components/DeducteeSearchCombobox";
import { generateMockDeductees } from "../app/compliance/tds/mock-data";
import type { Deductee } from "../app/compliance/tds/types";

const mockDeductees = generateMockDeductees();

const meta: Meta<typeof DeducteeSearchCombobox> = {
  title: "Compliance/TDS/DeducteeSearchCombobox",
  component: DeducteeSearchCombobox,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DeducteeSearchCombobox>;

export const Default: Story = {
  render: () => {
    const [value, setValue] = useState<Deductee | null>(null);
    return (
      <div className="max-w-md">
        <DeducteeSearchCombobox
          deductees={mockDeductees}
          value={value}
          onChange={setValue}
        />
      </div>
    );
  },
};

export const WithSelection: Story = {
  render: () => {
    const [value, setValue] = useState<Deductee | null>(mockDeductees[0]);
    return (
      <div className="max-w-md">
        <DeducteeSearchCombobox
          deductees={mockDeductees}
          value={value}
          onChange={setValue}
        />
      </div>
    );
  },
};
