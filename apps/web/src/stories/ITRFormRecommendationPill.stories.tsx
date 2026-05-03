import type { Meta, StoryObj } from "@storybook/react";
import { ITRFormRecommendationPill } from "../app/compliance/itr/components/ITRFormRecommendationPill";

const meta: Meta<typeof ITRFormRecommendationPill> = {
  title: "Compliance/ITR/ITRFormRecommendationPill",
  component: ITRFormRecommendationPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ITRFormRecommendationPill>;

export const ITR1: Story = { args: { form: "ITR-1" } };
export const ITR2: Story = { args: { form: "ITR-2" } };
export const ITR3: Story = { args: { form: "ITR-3" } };
export const ITR4: Story = { args: { form: "ITR-4" } };
export const ITR5: Story = { args: { form: "ITR-5" } };
export const ITR6: Story = { args: { form: "ITR-6" } };
export const ITR7: Story = { args: { form: "ITR-7" } };
