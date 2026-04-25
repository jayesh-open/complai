import type { Meta, StoryObj } from "@storybook/react";
import { FileText, Receipt, FileCheck2, AlertTriangle } from "lucide-react";

function KpiMetricCard({ icon, iconColor = "var(--accent)", value, label, subtitle, trend }: {
  icon: React.ReactNode; iconColor?: string; value: string; label: string;
  subtitle?: string; trend?: { value: string; favorable: boolean };
}) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-[16px] p-4 hover:border-[var(--accent)] transition-colors w-[240px]">
      <div className="flex items-start justify-between mb-3">
        <div className="w-[30px] h-[30px] rounded-lg flex items-center justify-center border"
          style={{
            background: `color-mix(in srgb, ${iconColor} 15%, transparent)`,
            borderColor: `color-mix(in srgb, ${iconColor} 30%, transparent)`,
          }}>
          {icon}
        </div>
        {trend && (
          <span className={`text-[11px] font-semibold ${trend.favorable ? "text-[var(--success)]" : "text-[var(--danger)]"}`}>
            {trend.favorable ? "↑" : "↓"} {trend.value}
          </span>
        )}
      </div>
      <div className="text-[20px] font-bold leading-tight text-[var(--text-primary)]" style={{ fontFeatureSettings: '"tnum"' }}>
        {value}
      </div>
      <div className="text-[10px] font-semibold uppercase tracking-[0.05em] text-[var(--text-muted)] mt-1">{label}</div>
      {subtitle && (
        <>
          <div className="border-t border-[var(--border-light)] my-2" />
          <div className="text-xs text-[var(--text-muted)]">{subtitle}</div>
        </>
      )}
    </div>
  );
}

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
