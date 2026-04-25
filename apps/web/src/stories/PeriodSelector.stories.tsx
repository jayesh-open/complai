import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";

const FY_OPTIONS = ["2024-25", "2025-26", "2026-27"];
const MONTH_OPTIONS = ["Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb", "Mar"];

function PeriodSelector({ financialYear, month, onYearChange, onMonthChange }: {
  financialYear: string; month: string;
  onYearChange: (y: string) => void; onMonthChange: (m: string) => void;
}) {
  return (
    <div className="flex items-center gap-2">
      <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">Period:</span>
      <select value={financialYear} onChange={(e) => onYearChange(e.target.value)}
        className="px-3 py-1.5 rounded-lg border text-xs font-medium bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]">
        {FY_OPTIONS.map((fy) => <option key={fy} value={fy}>FY {fy}</option>)}
      </select>
      <select value={month} onChange={(e) => onMonthChange(e.target.value)}
        className="px-3 py-1.5 rounded-lg border text-xs font-medium bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]">
        {MONTH_OPTIONS.map((m) => (
          <option key={m} value={m}>{m} {financialYear.split("-")[m === "Jan" || m === "Feb" || m === "Mar" ? 1 : 0]}</option>
        ))}
      </select>
    </div>
  );
}

const meta: Meta<typeof PeriodSelector> = {
  title: "Compliance/PeriodSelector",
  component: PeriodSelector,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof PeriodSelector>;

export const Default: Story = {
  render: () => {
    const [fy, setFy] = useState("2026-27");
    const [month, setMonth] = useState("Apr");
    return <PeriodSelector financialYear={fy} month={month} onYearChange={setFy} onMonthChange={setMonth} />;
  },
};
