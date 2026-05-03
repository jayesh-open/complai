import type { Meta, StoryObj } from "@storybook/react";
import { TaxComputationPanel } from "../app/compliance/itr/components/TaxComputationPanel";
import type { TaxComputation, DeductionItem } from "../app/compliance/itr/types";

const meta: Meta<typeof TaxComputationPanel> = {
  title: "Compliance/ITR/TaxComputationPanel",
  component: TaxComputationPanel,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof TaxComputationPanel>;

const computation: TaxComputation = {
  totalIncome: 2275000,
  standardDeduction: 75000,
  taxableIncome: 2200000,
  slabs: [
    { from: 0, to: 400000, rate: 0, tax: 0 },
    { from: 400000, to: 800000, rate: 5, tax: 20000 },
    { from: 800000, to: 1200000, rate: 10, tax: 40000 },
    { from: 1200000, to: 1600000, rate: 15, tax: 60000 },
    { from: 1600000, to: 2000000, rate: 20, tax: 80000 },
    { from: 2000000, to: 2400000, rate: 25, tax: 50000 },
    { from: 2400000, to: null, rate: 30, tax: 0 },
  ],
  slabTax: 250000,
  surchargeRate: 0,
  surchargeAmount: 0,
  surchargeThreshold: "N/A",
  healthEducationCess: 10000,
  grossTax: 260000,
  rebate87A: 0,
  totalLiability: 260000,
  tdsCredit: 182000,
  advanceTax: 22750,
  selfAssessmentTax: 0,
  refundOrPayable: 55250,
};

export const NewRegime: Story = {
  args: { computation, regime: "NEW" },
};

const deductions: DeductionItem[] = [
  { section: "80C", label: "PPF, ELSS, LIC, Tuition", declared: 150000, limit: 150000 },
  { section: "80D", label: "Health Insurance Premium", declared: 25000, limit: 25000 },
  { section: "24(b)", label: "Home Loan Interest", declared: 200000, limit: 200000 },
];

export const OldRegimeWithDeductions: Story = {
  args: { computation: { ...computation, taxableIncome: 1825000 }, deductions, regime: "OLD" },
};

export const WithRebate87A: Story = {
  args: {
    computation: {
      ...computation,
      taxableIncome: 650000,
      slabTax: 12500,
      surchargeRate: 0,
      surchargeAmount: 0,
      healthEducationCess: 500,
      grossTax: 13000,
      rebate87A: 13000,
      totalLiability: 0,
    },
    regime: "NEW",
  },
};

export const WithSurcharge: Story = {
  args: {
    computation: {
      ...computation,
      taxableIncome: 12000000,
      slabTax: 3150000,
      surchargeRate: 15,
      surchargeAmount: 472500,
      surchargeThreshold: "> ₹1 Cr",
      healthEducationCess: 144900,
      grossTax: 3767400,
      totalLiability: 3767400,
    },
    regime: "NEW",
  },
};
