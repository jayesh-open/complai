import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { PaymentCodePicker } from "../app/compliance/tds/components/PaymentCodePicker";

const meta: Meta<typeof PaymentCodePicker> = {
  title: "Compliance/TDS/PaymentCodePicker",
  component: PaymentCodePicker,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof PaymentCodePicker>;

export const Default: Story = {
  render: () => {
    const [value, setValue] = useState("");
    return <PaymentCodePicker value={value} onChange={setValue} />;
  },
};

export const FilteredBySection: Story = {
  render: () => {
    const [value, setValue] = useState("");
    return <PaymentCodePicker value={value} onChange={setValue} sectionFilter="393(1)" />;
  },
};

export const Preselected: Story = {
  args: {
    value: "1027",
    onChange: () => {},
  },
};
