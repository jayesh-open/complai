import type { Meta, StoryObj } from "@storybook/react";
import { RegimeIndicator } from "../app/compliance/itr/components/RegimeIndicator";

const meta: Meta<typeof RegimeIndicator> = {
  title: "Compliance/ITR/RegimeIndicator",
  component: RegimeIndicator,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof RegimeIndicator>;

export const NewRegime: Story = { args: { regime: "NEW" } };
export const OldRegime: Story = { args: { regime: "OLD" } };
