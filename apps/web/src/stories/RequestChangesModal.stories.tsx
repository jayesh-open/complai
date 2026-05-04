import type { Meta, StoryObj } from "@storybook/react";
import { RequestChangesModal } from "../app/compliance/itr/components/RequestChangesModal";

const meta: Meta<typeof RequestChangesModal> = {
  title: "Compliance/ITR/RequestChangesModal",
  component: RequestChangesModal,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof RequestChangesModal>;

export const Open: Story = {
  args: {
    open: true,
    onClose: () => {},
    onSubmit: () => {},
    employeeName: "Rajesh Kumar",
  },
};

export const Closed: Story = {
  args: {
    open: false,
    onClose: () => {},
    onSubmit: () => {},
    employeeName: "Rajesh Kumar",
  },
};
