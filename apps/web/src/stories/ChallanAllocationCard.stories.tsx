import type { Meta, StoryObj } from "@storybook/react";
import { ChallanAllocationCard } from "../app/compliance/tds/components/ChallanAllocationCard";

const meta: Meta<typeof ChallanAllocationCard> = {
  title: "Compliance/TDS/Challans/ChallanAllocationCard",
  component: ChallanAllocationCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ChallanAllocationCard>;

export const FullyAllocated: Story = {
  args: {
    totalAmount: 150000,
    allocatedAmount: 150000,
    allocations: [
      { deducteeId: "ded-0001", deducteeName: "TCS Ltd", pan: "AAACT1234A", section: "393(1)", amount: 100000, entryId: "e1" },
      { deducteeId: "ded-0002", deducteeName: "Infosys Ltd", pan: "AABCI5678B", section: "393(1)", amount: 50000, entryId: "e2" },
    ],
  },
};

export const PartiallyAllocated: Story = {
  args: {
    totalAmount: 200000,
    allocatedAmount: 120000,
    allocations: [
      { deducteeId: "ded-0001", deducteeName: "TCS Ltd", pan: "AAACT1234A", section: "393(1)", amount: 120000, entryId: "e1" },
    ],
  },
};

export const NoAllocations: Story = {
  args: {
    totalAmount: 100000,
    allocatedAmount: 0,
    allocations: [],
  },
};
