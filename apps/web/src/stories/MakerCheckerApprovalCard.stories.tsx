import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";

function MakerCheckerApprovalCard({ title, subtitle, submittedBy, submittedAt, details, warnings, onApprove, onReject, onSendBack }: {
  title: string; subtitle: string; submittedBy: string; submittedAt: string;
  details: { label: string; value: string }[]; warnings?: string[];
  onApprove: (c: string) => void; onReject: (c: string) => void; onSendBack?: (c: string) => void;
}) {
  const [comment, setComment] = useState("");
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl">
      <div className="px-5 py-3 border-b border-[var(--border-default)]">
        <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--warning)]">Approval Required</div>
      </div>
      <div className="px-5 py-4 space-y-3">
        <div>
          <div className="text-sm font-bold text-[var(--text-primary)]">{title}</div>
          <div className="text-xs text-[var(--text-muted)]">{subtitle}</div>
        </div>
        <div className="text-[11px] text-[var(--text-muted)]">Submitted by: {submittedBy} · {submittedAt}</div>
        <div className="space-y-1.5">
          {details.map((d) => (
            <div key={d.label} className="flex gap-2 text-xs">
              <span className="text-[var(--text-muted)] min-w-[80px]">{d.label}:</span>
              <span className="text-[var(--text-primary)] font-medium">{d.value}</span>
            </div>
          ))}
        </div>
        {warnings?.map((w, i) => (
          <div key={i} className="flex items-center gap-1.5 text-[11px] text-[var(--warning)]">
            <span>⚠</span> {w}
          </div>
        ))}
        <div>
          <label className="text-[11px] font-medium text-[var(--text-muted)]">Comments (optional):</label>
          <textarea value={comment} onChange={(e) => setComment(e.target.value)} rows={2}
            className="w-full mt-1 px-3 py-2 rounded-lg border text-xs bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]" />
        </div>
      </div>
      <div className="px-5 py-3 border-t border-[var(--border-default)] flex items-center gap-2 justify-between">
        <div className="flex gap-2">
          <button onClick={() => onReject(comment)}
            className="px-3 py-1.5 text-[11px] font-semibold rounded-lg bg-[var(--danger-muted)] text-[var(--danger)]">Reject</button>
          {onSendBack && (
            <button onClick={() => onSendBack(comment)}
              className="px-3 py-1.5 text-[11px] font-medium rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)]">Send Back for Edit</button>
          )}
        </div>
        <button onClick={() => onApprove(comment)}
          className="px-4 py-1.5 text-[11px] font-bold rounded-lg bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]">
          Approve &amp; Continue
        </button>
      </div>
    </div>
  );
}

const meta: Meta<typeof MakerCheckerApprovalCard> = {
  title: "Compliance/MakerCheckerApprovalCard",
  component: MakerCheckerApprovalCard,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof MakerCheckerApprovalCard>;

export const Default: Story = {
  args: {
    title: "GSTR-3B Filing — April 2026",
    subtitle: "29AABCA1234A1Z5 · Karnataka",
    submittedBy: "Priya Sharma",
    submittedAt: "25/04/2026 14:30",
    details: [
      { label: "Period", value: "April 2026" },
      { label: "Tax Payable", value: "₹12,45,678" },
      { label: "ITC Claimed", value: "₹8,92,340" },
      { label: "Net Payable", value: "₹3,53,338" },
    ],
    warnings: [
      "ITC claimed exceeds auto-populated 2B amount by ₹42,000",
      "Late filing penalty of ₹2,000 may apply after 20th",
    ],
    onApprove: (c: string) => alert(`Approved${c ? `: ${c}` : ""}`),
    onReject: (c: string) => alert(`Rejected${c ? `: ${c}` : ""}`),
    onSendBack: (c: string) => alert(`Sent back${c ? `: ${c}` : ""}`),
  },
};

export const NoWarnings: Story = {
  args: {
    title: "TDS Return — Q4 FY 2025-26",
    subtitle: "AAACR1234A · Form 26Q",
    submittedBy: "Amit Verma",
    submittedAt: "24/04/2026 10:15",
    details: [
      { label: "Quarter", value: "Q4 (Jan-Mar 2026)" },
      { label: "Deductees", value: "142" },
      { label: "Total TDS", value: "₹18,72,400" },
    ],
    onApprove: (c: string) => alert(`Approved${c ? `: ${c}` : ""}`),
    onReject: (c: string) => alert(`Rejected${c ? `: ${c}` : ""}`),
  },
};
