import type { Meta, StoryObj } from "@storybook/react";
import { Form131PreviewPanel } from "../app/compliance/tds/components/Form131PreviewPanel";
import { generateForm131Detail } from "../app/compliance/tds/certificates/mock-data";

const meta: Meta<typeof Form131PreviewPanel> = {
  title: "Compliance/TDS/Certificates/Form131PreviewPanel",
  component: Form131PreviewPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof Form131PreviewPanel>;

export const Default: Story = {
  args: { data: generateForm131Detail("ded-0001", "2026-27", "Q1") },
};
