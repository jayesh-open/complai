import type { Meta, StoryObj } from "@storybook/react";
import { IncomeHeadCard } from "../app/compliance/itr/components/IncomeHeadCard";
import type { IncomeHeadDetail } from "../app/compliance/itr/types";

const meta: Meta<typeof IncomeHeadCard> = {
  title: "Compliance/ITR/IncomeHeadCard",
  component: IncomeHeadCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof IncomeHeadCard>;

const salary: IncomeHeadDetail = {
  head: "SALARY",
  gross: 1850000,
  deductions: 148000,
  net: 1702000,
  subItems: [
    { label: "Basic Pay", amount: 925000 },
    { label: "HRA", amount: 370000 },
    { label: "Special Allowance", amount: 370000 },
    { label: "Other Allowances", amount: 185000 },
    { label: "HRA Exemption", amount: -88800, section: "§392" },
    { label: "LTA Exemption", amount: -59200 },
  ],
  visible: true,
};

export const Salary: Story = { args: { detail: salary } };

export const HouseProperty: Story = {
  args: {
    detail: {
      head: "HOUSE_PROPERTY",
      gross: 240000,
      deductions: 72000,
      net: 168000,
      subItems: [
        { label: "Self-Occupied — NAV", amount: 0 },
        { label: "Let-Out — Gross Rental", amount: 240000 },
        { label: "Municipal Tax", amount: -24000 },
        { label: "Std Deduction (30%)", amount: -48000 },
      ],
      visible: true,
    },
  },
};

export const CapitalGains: Story = {
  args: {
    detail: {
      head: "CAPITAL_GAINS",
      gross: 320000,
      deductions: 0,
      net: 320000,
      subItems: [
        { label: "STCG (Listed Equity)", amount: 120000 },
        { label: "LTCG (§112A — up to ₹1.25L exempt)", amount: 150000, section: "§112A" },
        { label: "Schedule VDA (Crypto)", amount: 50000 },
      ],
      visible: true,
    },
  },
};

export const BusinessProfession: Story = {
  args: {
    detail: {
      head: "BUSINESS_PROFESSION",
      gross: 500000,
      deductions: 160000,
      net: 340000,
      subItems: [
        { label: "Gross Receipts", amount: 500000 },
        { label: "Presumptive (§44AD)", amount: -160000, section: "§44AD" },
      ],
      visible: true,
    },
  },
};

export const OtherSources: Story = {
  args: {
    detail: {
      head: "OTHER_SOURCES",
      gross: 85000,
      deductions: 0,
      net: 85000,
      subItems: [
        { label: "Savings Interest", amount: 35000 },
        { label: "FD Interest", amount: 28000 },
        { label: "Dividend Income", amount: 15000 },
        { label: "Gift Income", amount: 7000 },
      ],
      visible: true,
    },
  },
};
