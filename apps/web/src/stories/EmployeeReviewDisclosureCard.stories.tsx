import type { Meta, StoryObj } from "@storybook/react";
import { EmployeeReviewDisclosureCard } from "../app/compliance/itr/components/EmployeeReviewDisclosureCard";

const meta: Meta<typeof EmployeeReviewDisclosureCard> = {
  title: "Compliance/ITR/EmployeeReviewDisclosureCard",
  component: EmployeeReviewDisclosureCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof EmployeeReviewDisclosureCard>;

export const NewRegime: Story = {
  args: {
    employeeName: "Rajesh Kumar",
    taxYear: "2026-27",
    employerName: "Acme Technologies Pvt. Ltd.",
    recommendedForm: "ITR-1",
    regime: "NEW",
  },
};

export const OldRegime: Story = {
  args: {
    employeeName: "Priya Sharma",
    taxYear: "2026-27",
    employerName: "Acme Technologies Pvt. Ltd.",
    recommendedForm: "ITR-2",
    regime: "OLD",
  },
};

export const BusinessIncome: Story = {
  args: {
    employeeName: "Mohit Verma",
    taxYear: "2026-27",
    employerName: "Acme Technologies Pvt. Ltd.",
    recommendedForm: "ITR-3",
    regime: "NEW",
  },
};
