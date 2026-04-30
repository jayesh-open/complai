import type { Meta, StoryObj } from "@storybook/react";
import { Form130PreviewPanel } from "../app/compliance/tds/components/Form130PreviewPanel";
import { generateForm130Detail } from "../app/compliance/tds/certificates/mock-data";

const meta: Meta<typeof Form130PreviewPanel> = {
  title: "Compliance/TDS/Certificates/Form130PreviewPanel",
  component: Form130PreviewPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof Form130PreviewPanel>;

export const NewRegime: Story = {
  args: { data: generateForm130Detail("ded-0003", "2026-27") },
};

export const OldRegime: Story = {
  args: { data: generateForm130Detail("ded-0004", "2026-27") },
};
