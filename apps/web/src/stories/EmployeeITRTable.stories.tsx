import type { Meta, StoryObj } from "@storybook/react";
import { EmployeeITRTable } from "../app/compliance/itr/components/EmployeeITRTable";
import { ALL_EMPLOYEES } from "../app/compliance/itr/mock-data";

const meta: Meta<typeof EmployeeITRTable> = {
  title: "Compliance/ITR/EmployeeITRTable",
  component: EmployeeITRTable,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof EmployeeITRTable>;

export const Default: Story = { args: { employees: ALL_EMPLOYEES.slice(0, 10) } };
export const Empty: Story = { args: { employees: [] } };
export const SingleEmployee: Story = { args: { employees: ALL_EMPLOYEES.slice(0, 1) } };
