import type { Meta, StoryObj } from "@storybook/react";
import { ViewAsDropdown } from "../components/layout/ViewAsDropdown";

const meta: Meta<typeof ViewAsDropdown> = {
  title: "Layout/ViewAsDropdown",
  component: ViewAsDropdown,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ViewAsDropdown>;

export const Default: Story = {};
