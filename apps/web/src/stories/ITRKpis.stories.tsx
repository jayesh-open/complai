import type { Meta, StoryObj } from "@storybook/react";
import { ITRKpis } from "../app/compliance/itr/components/ITRKpis";
import { ALL_EMPLOYEES } from "../app/compliance/itr/mock-data";
import type { ITREmployee } from "../app/compliance/itr/types";

const meta: Meta<typeof ITRKpis> = {
  title: "Compliance/ITR/ITRKpis",
  component: ITRKpis,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ITRKpis>;

export const Default: Story = { args: { employees: ALL_EMPLOYEES } };

export const AllFiled: Story = {
  args: {
    employees: ALL_EMPLOYEES.map((e) => ({ ...e, filingStatus: "FILED" })) as ITREmployee[],
  },
};

export const NoneStarted: Story = {
  args: {
    employees: ALL_EMPLOYEES.map((e) => ({ ...e, filingStatus: "NOT_STARTED" })) as ITREmployee[],
  },
};
