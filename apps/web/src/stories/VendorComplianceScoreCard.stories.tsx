import type { Meta, StoryObj } from "@storybook/react";

interface ScoreDimension {
  label: string; score: number; maxScore: number; status: "pass" | "warn" | "fail"; note?: string;
}

function VendorComplianceScoreCard({ vendorName, gstin, state, totalScore, maxScore = 100, riskLevel, category, lastReviewed, dimensions }: {
  vendorName: string; gstin: string; state: string; totalScore: number; maxScore?: number;
  riskLevel: string; category: string; lastReviewed: string; dimensions: ScoreDimension[];
}) {
  const scoreColor = totalScore >= 80 ? "var(--success)" : totalScore >= 60 ? "var(--warning)" : "var(--danger)";
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl p-5">
      <div className="text-sm font-bold text-[var(--text-primary)]">{vendorName}</div>
      <div className="text-[11px] text-[var(--text-muted)] mt-0.5 font-mono">{gstin} · {state}</div>
      <div className="mt-4 flex items-center gap-3">
        <div className="flex gap-0.5">
          {Array.from({ length: 10 }, (_, i) => (
            <div key={i} className="w-2.5 h-2.5 rounded-full"
              style={{ backgroundColor: i < Math.round(totalScore / (maxScore / 10)) ? scoreColor : "var(--border-default)" }} />
          ))}
        </div>
        <span className="text-sm font-bold" style={{ color: scoreColor }}>{totalScore} / {maxScore}</span>
        <span className="text-[11px] text-[var(--text-muted)]">· Risk: {riskLevel}</span>
      </div>
      <div className="mt-3 bg-[var(--bg-tertiary)] border border-[var(--border-default)] rounded-lg p-3 space-y-2">
        {dimensions.map((d) => (
          <div key={d.label} className="flex items-center gap-2 text-xs">
            <span className="text-[var(--text-muted)] min-w-[140px]">{d.label}:</span>
            <span className="font-semibold text-[var(--text-primary)]">{d.score}/{d.maxScore}</span>
            <span>{d.status === "pass" ? "✓" : "⚠"}</span>
            {d.note && <span className="text-[var(--text-disabled)]">— {d.note}</span>}
          </div>
        ))}
      </div>
      <div className="mt-3 text-[11px] text-[var(--text-muted)]">
        Category: {category} · Last reviewed: {lastReviewed}
      </div>
    </div>
  );
}

const meta: Meta<typeof VendorComplianceScoreCard> = {
  title: "Compliance/VendorComplianceScoreCard",
  component: VendorComplianceScoreCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof VendorComplianceScoreCard>;

export const HighScore: Story = {
  args: {
    vendorName: "Tata Steel Limited",
    gstin: "29AABCT1234A1Z5",
    state: "Karnataka",
    totalScore: 92,
    riskLevel: "Low",
    category: "Strategic",
    lastReviewed: "20/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 20, maxScore: 20, status: "pass" },
      { label: "GSTR-2B Match Rate", score: 18, maxScore: 20, status: "pass", note: "98.2% match" },
      { label: "Payment Timeliness", score: 19, maxScore: 20, status: "pass" },
      { label: "Document Quality", score: 17, maxScore: 20, status: "pass", note: "All e-invoiced" },
      { label: "Compliance History", score: 18, maxScore: 20, status: "pass" },
    ],
  },
};

export const MediumScore: Story = {
  args: {
    vendorName: "QuickParts Industries",
    gstin: "27AADCQ5678B2Z3",
    state: "Maharashtra",
    totalScore: 64,
    riskLevel: "Medium",
    category: "Operational",
    lastReviewed: "18/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 14, maxScore: 20, status: "warn", note: "2 late filings" },
      { label: "GSTR-2B Match Rate", score: 12, maxScore: 20, status: "warn", note: "76% match" },
      { label: "Payment Timeliness", score: 16, maxScore: 20, status: "pass" },
      { label: "Document Quality", score: 10, maxScore: 20, status: "fail", note: "Manual invoices" },
      { label: "Compliance History", score: 12, maxScore: 20, status: "warn" },
    ],
  },
};

export const LowScore: Story = {
  args: {
    vendorName: "ABC Traders",
    gstin: "33AABCA9999C1Z1",
    state: "Tamil Nadu",
    totalScore: 38,
    riskLevel: "High",
    category: "Transactional",
    lastReviewed: "15/04/2026",
    dimensions: [
      { label: "GST Filing Status", score: 6, maxScore: 20, status: "fail", note: "4 months unfiled" },
      { label: "GSTR-2B Match Rate", score: 8, maxScore: 20, status: "fail", note: "42% match" },
      { label: "Payment Timeliness", score: 10, maxScore: 20, status: "warn" },
      { label: "Document Quality", score: 8, maxScore: 20, status: "fail", note: "Missing HSN codes" },
      { label: "Compliance History", score: 6, maxScore: 20, status: "fail", note: "SCN received" },
    ],
  },
};
