import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { QuarterPicker } from "../app/compliance/tds/file/components/QuarterPicker";

const meta: Meta<typeof QuarterPicker> = {
  title: "Compliance/TDS/Filing/QuarterPicker",
  component: QuarterPicker,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof QuarterPicker>;

export const Default: Story = {
  render: () => {
    const [value, setValue] = useState("q1");
    return <QuarterPicker value={value} onChange={setValue} />;
  },
};

export const WithStatuses: Story = {
  render: () => {
    const [value, setValue] = useState("q2");
    return (
      <QuarterPicker
        value={value}
        onChange={setValue}
        statuses={{ q1: "FILED", q2: "DRAFT", q3: "NOT_STARTED", q4: "NOT_STARTED" }}
      />
    );
  },
};
