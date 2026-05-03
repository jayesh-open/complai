import type { Meta, StoryObj } from "@storybook/react";
import { FilingStatusPill } from "../app/compliance/itr/components/FilingStatusPill";

const meta: Meta<typeof FilingStatusPill> = {
  title: "Compliance/ITR/FilingStatusPill",
  component: FilingStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FilingStatusPill>;

export const NotStarted: Story = { args: { status: "NOT_STARTED" } };
export const AISFetched: Story = { args: { status: "AIS_FETCHED" } };
export const FormGenerated: Story = { args: { status: "FORM_GENERATED" } };
export const ReviewPending: Story = { args: { status: "REVIEW_PENDING" } };
export const EmployeeApproved: Story = { args: { status: "EMPLOYEE_APPROVED" } };
export const Filed: Story = { args: { status: "FILED" } };
export const Acknowledged: Story = { args: { status: "ACKNOWLEDGED" } };
export const Defective: Story = { args: { status: "DEFECTIVE" } };
