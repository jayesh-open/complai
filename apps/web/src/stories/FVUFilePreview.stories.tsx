import type { Meta, StoryObj } from "@storybook/react";
import { FVUFilePreview } from "../app/compliance/tds/file/components/FVUFilePreview";

const meta: Meta<typeof FVUFilePreview> = {
  title: "Compliance/TDS/Filing/FVUFilePreview",
  component: FVUFilePreview,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FVUFilePreview>;

const sampleFVU = `^FVU|FORM140|5.2|
^BH|MUMA12345B|2026-27|Q1|ORIGINAL|
^CH|1|MUMA12345B|2026-27|Q1|5|1250000|

^DD|AAACT1234A|TCS Ltd|4|500000|10000|0|0|10000|
^DD|AABCI5678B|Infosys Ltd|3|300000|30000|0|0|30000|
^DD|GHIPA5678K|Amit Jain|2|75000|3750|0|0|3750|

^TV|TOTAL|5|875000|43750|0|0|43750|
^FV|END|`;

export const Default: Story = {
  args: {
    content: sampleFVU,
    formLabel: "Form 140 (Non-Salary Resident)",
  },
};
