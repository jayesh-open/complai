import type { Meta, StoryObj } from "@storybook/react";
import { KpiMetricCard } from "@complai/ui-components";
import { FileText, Receipt, FileCheck2, AlertTriangle } from "lucide-react";

const meta: Meta<typeof KpiMetricCard> = {
  title: "Compliance/KpiMetricCard",
  component: KpiMetricCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof KpiMetricCard>;

export const Default: Story = {
  args: {
    icon: <FileText className="w-4 h-4" style={{ color: "var(--accent)" }} />,
    value: "228",
    label: "Total Invoices",
    subtitle: "₹8.76 Cr",
    trend: { value: "12%", favorable: true },
  },
};

export const CardGrid: Story = {
  render: () => (
    <div className="flex gap-4 flex-wrap">
      <KpiMetricCard
        icon={<FileText className="w-4 h-4" style={{ color: "var(--accent)" }} />}
        iconColor="var(--accent)" value="228" label="Total Invoices"
        subtitle="₹8.76 Cr" trend={{ value: "12%", favorable: true }}
      />
      <KpiMetricCard
        icon={<Receipt className="w-4 h-4" style={{ color: "var(--info)" }} />}
        iconColor="var(--info)" value="40" label="Pending Approval"
        subtitle="₹1.52 Cr" trend={{ value: "8%", favorable: false }}
      />
      <KpiMetricCard
        icon={<FileCheck2 className="w-4 h-4" style={{ color: "var(--success)" }} />}
        iconColor="var(--success)" value="₹45.2L" label="ITC Claimed YTD"
        subtitle="Eligible: ₹52.8L" trend={{ value: "15%", favorable: true }}
      />
      <KpiMetricCard
        icon={<AlertTriangle className="w-4 h-4" style={{ color: "var(--danger)" }} />}
        iconColor="var(--danger)" value="7" label="Overdue Filings"
        subtitle="GSTR-1: 3 · TDS: 4"
      />
    </div>
  ),
};
