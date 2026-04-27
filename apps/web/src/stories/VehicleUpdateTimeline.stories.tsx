import type { Meta, StoryObj } from "@storybook/react";
import { VehicleUpdateTimeline } from "../app/compliance/e-way-bill/components/VehicleUpdateTimeline";

const meta: Meta<typeof VehicleUpdateTimeline> = {
  title: "Compliance/VehicleUpdateTimeline",
  component: VehicleUpdateTimeline,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof VehicleUpdateTimeline>;

export const SingleEntry: Story = {
  args: {
    entries: [
      {
        vehicleNo: "KA01AB1234",
        updatedAt: new Date(Date.now() - 3600000 * 10).toISOString(),
        updatedBy: "finance@acme.com",
      },
    ],
  },
};

export const MultipleUpdates: Story = {
  args: {
    entries: [
      {
        vehicleNo: "KA01AB1234",
        updatedAt: new Date(Date.now() - 3600000 * 10).toISOString(),
        updatedBy: "finance@acme.com",
      },
      {
        vehicleNo: "MH02CD5678",
        updatedAt: new Date(Date.now() - 3600000 * 6).toISOString(),
        updatedBy: "logistics@acme.com",
        reason: "Vehicle breakdown — replaced",
        fromPlace: "Pune",
      },
      {
        vehicleNo: "DL04GH3456",
        updatedAt: new Date(Date.now() - 3600000 * 1).toISOString(),
        updatedBy: "logistics@acme.com",
        reason: "Transshipment at Delhi hub",
        fromPlace: "Delhi",
      },
    ],
  },
};
