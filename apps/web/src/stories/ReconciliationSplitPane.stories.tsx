import type { Meta, StoryObj } from "@storybook/react";

function ReconciliationSplitPane({ leftTitle, rightTitle, leftContent, rightContent }: {
  leftTitle: string; rightTitle: string; leftContent: React.ReactNode; rightContent: React.ReactNode;
}) {
  return (
    <div className="border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div className="grid grid-cols-2 divide-x divide-[var(--border-default)]">
        <div className="px-4 py-2.5 bg-[var(--bg-tertiary)]">
          <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">{leftTitle}</span>
        </div>
        <div className="px-4 py-2.5 bg-[var(--bg-tertiary)]">
          <span className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">{rightTitle}</span>
        </div>
      </div>
      <div className="grid grid-cols-2 divide-x divide-[var(--border-default)] min-h-[200px]">
        <div className="p-4">{leftContent}</div>
        <div className="p-4">{rightContent}</div>
      </div>
    </div>
  );
}

function InvoiceRow({ no, vendor, amount, match }: { no: string; vendor: string; amount: string; match?: boolean }) {
  return (
    <div className={`flex items-center gap-3 px-3 py-2 rounded-lg text-xs ${match === true ? "bg-[color-mix(in_srgb,var(--success)_8%,transparent)]" : match === false ? "bg-[color-mix(in_srgb,var(--danger)_8%,transparent)]" : ""}`}>
      <span className="font-mono text-[var(--text-primary)] min-w-[120px]">{no}</span>
      <span className="text-[var(--text-secondary)] flex-1">{vendor}</span>
      <span className="tabular-nums font-semibold text-[var(--text-primary)]">{amount}</span>
      {match === true && <span className="text-[var(--success)] text-[10px] font-semibold">MATCHED</span>}
      {match === false && <span className="text-[var(--danger)] text-[10px] font-semibold">MISMATCH</span>}
    </div>
  );
}

const meta: Meta<typeof ReconciliationSplitPane> = {
  title: "Compliance/ReconciliationSplitPane",
  component: ReconciliationSplitPane,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ReconciliationSplitPane>;

export const GSTR1vs2B: Story = {
  render: () => (
    <ReconciliationSplitPane
      leftTitle="GSTR-1 (Sales Register)"
      rightTitle="GSTR-2B (Auto-populated)"
      leftContent={
        <div className="space-y-1">
          <InvoiceRow no="INV-2026-0042" vendor="Tata Steel Ltd" amount="₹4,52,800" match={true} />
          <InvoiceRow no="INV-2026-0041" vendor="Infosys BPO" amount="₹1,28,000" match={true} />
          <InvoiceRow no="INV-2026-0040" vendor="Reliance Jio" amount="₹86,400" match={false} />
          <InvoiceRow no="INV-2026-0039" vendor="Wipro Tech" amount="₹2,15,600" match={true} />
          <InvoiceRow no="INV-2026-0038" vendor="Mahindra Log" amount="₹3,72,000" />
        </div>
      }
      rightContent={
        <div className="space-y-1">
          <InvoiceRow no="INV-2026-0042" vendor="Tata Steel Ltd" amount="₹4,52,800" match={true} />
          <InvoiceRow no="INV-2026-0041" vendor="Infosys BPO" amount="₹1,28,000" match={true} />
          <InvoiceRow no="INV-2026-0040" vendor="Reliance Jio" amount="₹84,200" match={false} />
          <InvoiceRow no="INV-2026-0039" vendor="Wipro Tech" amount="₹2,15,600" match={true} />
          <div className="flex items-center justify-center h-10 text-xs text-[var(--text-muted)]">Not found in 2B</div>
        </div>
      }
    />
  ),
};
