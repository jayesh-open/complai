import type { Meta, StoryObj } from "@storybook/react";
import { FilingStatusOverview } from "../app/compliance/tds/file/components/FilingStatusOverview";
import { generateFilingGrid } from "../app/compliance/tds/file/mock-data";

const meta: Meta<typeof FilingStatusOverview> = {
  title: "Compliance/TDS/Filing/FilingStatusOverview",
  component: FilingStatusOverview,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FilingStatusOverview>;

export const Default: Story = {
  args: {
    cells: generateFilingGrid(),
    taxYear: "2026-27",
  },
};
