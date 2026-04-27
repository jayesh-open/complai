import type { Meta, StoryObj } from "@storybook/react";
import { DistanceValidityCalculator } from "../app/compliance/e-way-bill/components/DistanceValidityCalculator";

const meta: Meta<typeof DistanceValidityCalculator> = {
  title: "Compliance/DistanceValidityCalculator",
  component: DistanceValidityCalculator,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DistanceValidityCalculator>;

export const Short: Story = { args: { distanceKm: 150 } };
export const Medium: Story = { args: { distanceKm: 450 } };
export const Long: Story = { args: { distanceKm: 1200 } };
export const ODC: Story = { args: { distanceKm: 100, isODC: true } };

export const AllVariants: Story = {
  render: () => (
    <div className="space-y-3 max-w-md">
      <DistanceValidityCalculator distanceKm={100} />
      <DistanceValidityCalculator distanceKm={450} />
      <DistanceValidityCalculator distanceKm={1200} />
      <DistanceValidityCalculator distanceKm={60} isODC />
    </div>
  ),
};
