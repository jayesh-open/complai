import type { Meta, StoryObj } from "@storybook/react";
import { TDSCalculatorPanel } from "../app/compliance/tds/components/TDSCalculatorPanel";

const meta: Meta<typeof TDSCalculatorPanel> = {
  title: "Compliance/TDS/TDSCalculatorPanel",
  component: TDSCalculatorPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof TDSCalculatorPanel>;

export const ContractorPayment: Story = {
  args: {
    paymentCode: "1024",
    amount: 500000,
    isNonResident: false,
    noPan: false,
  },
};

export const ProfessionalFees: Story = {
  args: {
    paymentCode: "1027",
    amount: 300000,
    isNonResident: false,
    noPan: false,
  },
};

export const NonResidentWithCess: Story = {
  args: {
    paymentCode: "1057",
    amount: 1000000,
    isNonResident: true,
    noPan: false,
  },
};

export const NoPanRate: Story = {
  args: {
    paymentCode: "1024",
    amount: 200000,
    isNonResident: false,
    noPan: true,
  },
};

export const EmptyState: Story = {
  args: {
    paymentCode: "",
    amount: 0,
    isNonResident: false,
    noPan: false,
  },
};
