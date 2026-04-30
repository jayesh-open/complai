import type { Meta, StoryObj } from "@storybook/react";
import { DTAAEvidenceBadge } from "../app/compliance/tds/components/DTAAEvidenceBadge";

const meta: Meta<typeof DTAAEvidenceBadge> = {
  title: "Compliance/TDS/DTAAEvidenceBadge",
  component: DTAAEvidenceBadge,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DTAAEvidenceBadge>;

export const AllClear: Story = {
  args: {
    form41Filed: true,
    trcAttached: true,
  },
};

export const MissingTRC: Story = {
  args: {
    form41Filed: true,
    trcAttached: false,
  },
};

export const MissingForm41: Story = {
  args: {
    form41Filed: false,
    trcAttached: true,
  },
};

export const NeitherFiled: Story = {
  args: {
    form41Filed: false,
    trcAttached: false,
  },
};
