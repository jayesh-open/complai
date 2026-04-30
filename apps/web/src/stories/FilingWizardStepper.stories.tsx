import type { Meta, StoryObj } from "@storybook/react";
import { FilingWizardStepper } from "../app/compliance/tds/file/components/FilingWizardStepper";

const meta: Meta<typeof FilingWizardStepper> = {
  title: "Compliance/TDS/Filing/FilingWizardStepper",
  component: FilingWizardStepper,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FilingWizardStepper>;

export const AtPull: Story = { args: { currentStep: "pull" } };
export const AtValidate: Story = { args: { currentStep: "validate" } };
export const AtPreview: Story = { args: { currentStep: "preview" } };
export const AtSubmit: Story = { args: { currentStep: "submit" } };
export const AtAcknowledge: Story = { args: { currentStep: "acknowledge" } };
