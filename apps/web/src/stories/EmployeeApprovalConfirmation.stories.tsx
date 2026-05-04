import type { Meta, StoryObj } from "@storybook/react";
import { EmployeeApprovalConfirmation } from "../app/compliance/itr/components/EmployeeApprovalConfirmation";

const meta: Meta<typeof EmployeeApprovalConfirmation> = {
  title: "Compliance/ITR/EmployeeApprovalConfirmation",
  component: EmployeeApprovalConfirmation,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof EmployeeApprovalConfirmation>;

export const WithEVC: Story = {
  args: {
    arn: "CPC/202627/00012345",
    filedAt: "2026-07-25T14:30:00Z",
    signMethod: "evc",
  },
};

export const WithDSC: Story = {
  args: {
    arn: "CPC/202627/00067890",
    filedAt: "2026-07-25T16:00:00Z",
    signMethod: "dsc",
  },
};
