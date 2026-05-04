import type { Meta, StoryObj } from "@storybook/react";
import { DSCSigningPlaceholder } from "../app/compliance/gst-returns/annual/components/DSCSigningPlaceholder";

const meta: Meta<typeof DSCSigningPlaceholder> = {
  title: "Compliance/GST/DSCSigningPlaceholder",
  component: DSCSigningPlaceholder,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DSCSigningPlaceholder>;

export const Default: Story = { args: {} };
export const CustomLabel: Story = { args: { label: "Sign GSTR-9C with DSC" } };
export const WithCallback: Story = {
  args: {
    label: "Sign & Submit",
    onSigned: () => console.log("DSC signed"),
  },
};
